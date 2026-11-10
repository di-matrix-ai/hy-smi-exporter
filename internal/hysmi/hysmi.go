// Package hysmi is the shell-out backend that runs `hy-smi` (or
// `hcu-smi`) and parses its JSON or human-readable output into typed
// Go structs.
//
// Best-effort: failures (missing binary, timeout, parse errors) are
// logged at warn level and the backend silently contributes nothing
// to that scrape. Callers must tolerate an empty device list.
//
// Only read-only hy-smi options are invoked. No --set*/--rasinject/
// --load*driver write operations are ever issued by this package.
package hysmi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Backend wraps the hy-smi / hcu-smi binary.
type Backend struct {
	binary  string
	timeout time.Duration

	// missingLogged is set after the first "binary not found" warning
	// so the log doesn't fill up if hy-smi never appears.
	missingLogged bool
}

// New returns a Backend bound to the given binary path. binary may
// be an absolute path or a name to look up via $PATH. An empty value
// defaults to "hy-smi".
func New(binary string, timeout time.Duration) *Backend {
	if binary == "" {
		binary = "hy-smi"
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Backend{binary: binary, timeout: timeout}
}

// Device is the per-card payload returned by hy-smi --json.
// Field tags map the human-readable keys hy-smi emits in JSON mode;
// the parser also tolerates lower-case / underscore variants.
type Device struct {
	// Identity
	Index         int    `json:"device_index"`
	HCU           string `json:"hcu"`
	DeviceID      string `json:"device_id"`
	UniqueID      string `json:"unique_id"`
	Serial        string `json:"serial"`
	ProductName   string `json:"product_name"`
	PartNumber    string `json:"part_number"`
	VBIOSVersion  string `json:"vbios_version"`
	DriverVersion string `json:"driver_version"`
	Bus           string `json:"bus"`

	// Thermal / power
	TempCelsius   float64 `json:"temp_celsius"`
	EdgeTempC     float64 `json:"edge_temp_c"`
	JunctionTempC float64 `json:"junction_temp_c"`
	MemTempC      float64 `json:"mem_temp_c"`
	HotspotTempC  float64 `json:"hotspot_temp_c"`
	PowerWatts    float64 `json:"power_watts"`
	PowerCapWatts float64 `json:"power_cap_watts"`
	MaxPowerWatts float64 `json:"max_power_watts"`
	VoltageVolts  float64 `json:"voltage_volts"`

	// Utilization (percent, 0-100)
	HCUUtilPct  float64 `json:"hcu_util_pct"`
	CUUtilPct   float64 `json:"cu_util_pct"`
	WaveUtilPct float64 `json:"wave_util_pct"`
	DecUtilPct  float64 `json:"dec_util_pct"`
	EncUtilPct  float64 `json:"enc_util_pct"`
	MemUtilPct  float64 `json:"mem_util_pct"`
	VRAMUtilPct float64 `json:"vram_util_pct"`

	// Memory (bytes; hy-smi exposes MiB but we normalise)
	VRAMTotalBytes uint64 `json:"vram_total_bytes"`
	VRAMUsedBytes  uint64 `json:"vram_used_bytes"`
	VRAMAvailBytes uint64 `json:"vram_available_bytes"`
	GTTTotalBytes  uint64 `json:"gtt_total_bytes"`
	GTTUsedBytes   uint64 `json:"gtt_used_bytes"`

	// Clocks (MHz)
	SCLKMHz    float64 `json:"sclk_mhz"`
	MCLKMHz    float64 `json:"mclk_mhz"`
	SOCCLKMHz  float64 `json:"socclk_mhz"`
	MaxSCLKMHz float64 `json:"max_sclk_mhz"`
	MinSCLKMHz float64 `json:"min_sclk_mhz"`
	BoostLevel string  `json:"boost_level"`
	PerfLevel  string  `json:"perf_level"`
	PerfState  string  `json:"perf_state"`
	FDRReason  string  `json:"fdr_reason"`

	// PCIe
	PCIELinkSpeedGTs float64 `json:"pcie_link_speed_gts"`
	PCIELinkWidth    int     `json:"pcie_link_width"`
	PCIEBandwidthMBs float64 `json:"pcie_bandwidth_mbs"`

	// Interconnect
	XHCLStatus string  `json:"xhcl_status"`
	XHCLBwMBs  float64 `json:"xhcl_bw_mbs"`
	DFBwMBs    float64 `json:"df_bw_mbs"`
	UMCBwMBs   float64 `json:"umc_bw_mbs"`

	// Health
	Health           string `json:"health"`
	Mode             string `json:"mode"`
	RASErrCount      uint64 `json:"ras_err_count"`
	ECCCorrectable   uint64 `json:"ecc_correctable"`
	ECCUncorrectable uint64 `json:"ecc_uncorrectable"`
	BadPagesRetired  uint64 `json:"bad_pages_retired"`
	BadPagesPending  uint64 `json:"bad_pages_pending"`

	// Fan
	FanRPM    float64 `json:"fan_rpm"`
	FanPWMPct float64 `json:"fan_pwm_pct"`
}

// Snapshot is the parsed hy-smi output. Devices is keyed by device
// index (0..N-1) in device-index order. DriverVersion / HostName
// come from the system block hy-smi emits alongside per-card data.
type Snapshot struct {
	Devices       []Device
	DriverVersion string
	HostName      string
}

// Snapshot runs the read-only hy-smi query and returns the parsed
// snapshot. Newer hy-smi versions reject --json with --showallinfo,
// so a failed JSON parse falls back to parsing the text emitted by
// `hy-smi -a`. On any failure the function returns an empty Snapshot
// and a nil error — the caller simply sees zero devices.
func (b *Backend) Snapshot(ctx context.Context, logger *slog.Logger) Snapshot {
	if _, err := exec.LookPath(b.binary); err != nil {
		if !b.missingLogged {
			logger.Warn("hy-smi binary not found; backend disabled for this scrape",
				"binary", b.binary)
			b.missingLogged = true
		}
		return Snapshot{}
	}
	b.missingLogged = false

	cctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	// Only read-only flags. Never --set*/--rasinject/--load*driver.
	out, err := b.run(cctx, "-a", "--json")
	if err != nil {
		logger.Debug("hy-smi -a --json failed; trying text output", "err", err)
	} else if snap, parseErr := parse(out); parseErr == nil && len(snap.Devices) > 0 {
		b.enrichMemory(cctx, logger, &snap)
		return snap
	} else if parseErr != nil {
		logger.Debug("hy-smi JSON output unavailable; trying text output", "err", parseErr)
	}

	textOut, err := b.run(cctx, "-a")
	if err != nil {
		logCommandError(logger, "hy-smi -a", err)
		return Snapshot{}
	}
	snap, err := parseText(textOut)
	if err != nil {
		logger.Warn("hy-smi output parse error", "err", err)
		return Snapshot{}
	}
	b.enrichMemory(cctx, logger, &snap)
	return snap
}

func (b *Backend) run(ctx context.Context, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, b.binary, args...).Output() //nolint:gosec // binary path is operator-controlled
}

func logCommandError(logger *slog.Logger, command string, err error) {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		logger.Warn(command+" exited non-zero",
			"code", ee.ExitCode(),
			"stderr", trim(string(ee.Stderr)),
		)
		return
	}
	logger.Warn(command+" failed", "err", err)
}

func (b *Backend) enrichMemory(ctx context.Context, logger *slog.Logger, snap *Snapshot) {
	if !needsMemory(*snap) {
		return
	}
	queries := [][]string{
		{"--showmeminfo", "all", "--json"},
		{"--showmemavailable", "--json"},
	}
	for _, args := range queries {
		out, err := b.run(ctx, args...)
		if err != nil {
			logger.Debug("hy-smi memory query failed", "args", args, "err", err)
			continue
		}
		memory, err := parse(out)
		if err != nil {
			logger.Debug("hy-smi memory JSON parse failed", "args", args, "err", err)
			continue
		}
		mergeSnapshot(snap, memory)
	}
}

func needsMemory(snap Snapshot) bool {
	for _, d := range snap.Devices {
		if d.VRAMTotalBytes == 0 || d.VRAMUsedBytes == 0 || d.VRAMAvailBytes == 0 || d.GTTTotalBytes == 0 {
			return true
		}
	}
	return false
}

func mergeSnapshot(dst *Snapshot, src Snapshot) {
	byIndex := make(map[int]int, len(dst.Devices))
	for i := range dst.Devices {
		byIndex[dst.Devices[i].Index] = i
	}
	for _, source := range src.Devices {
		targetIndex, ok := byIndex[source.Index]
		if !ok {
			dst.Devices = append(dst.Devices, source)
			byIndex[source.Index] = len(dst.Devices) - 1
			continue
		}
		mergeDevice(&dst.Devices[targetIndex], &source)
	}
	sortDevices(dst.Devices)
}

func mergeDevice(dst, src *Device) {
	if dst.VRAMTotalBytes == 0 {
		dst.VRAMTotalBytes = src.VRAMTotalBytes
	}
	if dst.VRAMUsedBytes == 0 {
		dst.VRAMUsedBytes = src.VRAMUsedBytes
	}
	if dst.VRAMAvailBytes == 0 {
		dst.VRAMAvailBytes = src.VRAMAvailBytes
	}
	if dst.GTTTotalBytes == 0 {
		dst.GTTTotalBytes = src.GTTTotalBytes
	}
	if dst.GTTUsedBytes == 0 {
		dst.GTTUsedBytes = src.GTTUsedBytes
	}
}

// parse accepts the shapes hy-smi emits with --json:
//
//	{ "system": {...}, "card_0": {...}, "card_1": {...} }      // map form
//	[ {...}, {...} ]                                            // array form
//	{ "devices": [ {...}, {...} ] }                            // nested form
//
// Field names hy-smi actually emits are human-readable strings like
// "Temperature (Sensor edge) (C)" or "Average Graphics Package Power (W)".
// We normalise those into the Device struct via normaliseKey + a
// field table; unknown keys are silently dropped.
func parse(raw []byte) (Snapshot, error) {
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return Snapshot{}, fmt.Errorf("decode json: %w", err)
	}

	snap := Snapshot{}
	rows := extractRows(doc)
	for _, row := range rows {
		if isSystemBlock(row) {
			snap.DriverVersion = pickString(row, []string{"driver_version"}, []string{"Driver Version"})
			snap.HostName = pickString(row, []string{"hostname"}, []string{"Host Name"})
			continue
		}
		snap.Devices = append(snap.Devices, parseDevice(row))
	}
	// Sort by device index to keep scrape output stable.
	sortDevices(snap.Devices)
	return snap, nil
}

// extractRows normalises the three input shapes into a flat slice
// of row maps. Returns nil if the document isn't an object/array.
func extractRows(doc any) []map[string]any {
	switch v := doc.(type) {
	case map[string]any:
		// Nested "devices": [ ... ] form.
		if devs, ok := v["devices"].([]any); ok {
			out := make([]map[string]any, 0, len(devs))
			for _, d := range devs {
				if m, ok := d.(map[string]any); ok {
					out = append(out, m)
				}
			}
			return out
		}
		// Map-of-cards form: keep system/card_N entries as rows. Some
		// hy-smi JSON commands identify cards only through the map key,
		// so preserve card0/card_0 as Device Index for the parser.
		out := make([]map[string]any, 0, len(v))
		for key, val := range v {
			if m, ok := val.(map[string]any); ok {
				if idx, ok := cardIndex(key); ok && !hasDeviceIndex(m) {
					copyMap := make(map[string]any, len(m)+1)
					for k, value := range m {
						copyMap[k] = value
					}
					copyMap["device_index"] = idx
					m = copyMap
				}
				out = append(out, m)
			}
		}
		return out
	case []any:
		out := make([]map[string]any, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	default:
		return nil
	}
}

func cardIndex(key string) (int, bool) {
	key = strings.TrimSpace(strings.ToLower(key))
	if !strings.HasPrefix(key, "card") {
		return 0, false
	}
	key = strings.TrimPrefix(key, "card")
	key = strings.TrimPrefix(key, "_")
	if key == "" {
		return 0, false
	}
	idx, err := strconv.Atoi(key)
	return idx, err == nil && idx >= 0
}

func hasDeviceIndex(row map[string]any) bool {
	for _, key := range []string{"device_index", "Device Index"} {
		if _, ok := row[key]; ok {
			return true
		}
	}
	return false
}

func isSystemBlock(row map[string]any) bool {
	// hy-smi's system block carries Driver Version / Host Name but no
	// per-card fields. Treat any row lacking a device index as system.
	for _, key := range []string{"device_index", "Device Index", "Device ID", "HCU"} {
		if _, ok := row[key]; ok {
			return false
		}
	}
	return true
}

// parseDevice walks the row and fills in whichever Device fields
// match. Unknown keys are silently dropped. Numeric fields that
// come in as strings ("157.0W", "40.0C") are stripped and parsed.
func parseDevice(row map[string]any) Device {
	d := Device{
		Index:         int(pickInt64(row, []string{"device_index"}, []string{"Device Index"})),
		HCU:           pickString(row, []string{"hcu"}, []string{"HCU"}),
		DeviceID:      pickString(row, []string{"device_id"}, []string{"Device ID"}),
		UniqueID:      pickString(row, []string{"unique_id"}, []string{"Unique ID"}),
		Serial:        pickString(row, []string{"serial"}, []string{"Serial Number"}),
		ProductName:   pickString(row, []string{"product_name"}, []string{"Product Name", "SKU", "Market Name", "Card Series"}),
		PartNumber:    pickString(row, []string{"part_number"}, []string{"Part Number"}),
		VBIOSVersion:  pickString(row, []string{"vbios_version"}, []string{"VBIOS Version"}),
		DriverVersion: pickString(row, []string{"driver_version"}, []string{"Driver Version"}),
		Bus:           pickString(row, []string{"bus"}, []string{"Bus", "PCI Bus"}),
		BoostLevel:    pickString(row, []string{"boost_level"}, []string{"Boost Level"}),
		PerfLevel:     pickString(row, []string{"perf_level"}, []string{"Performance Level"}),
		PerfState:     pickString(row, []string{"perf_state"}, []string{"Performance State"}),
		FDRReason:     pickString(row, []string{"fdr_reason"}, []string{"Frequency Change Reason"}),
		XHCLStatus:    pickString(row, []string{"xhcl_status"}, []string{"xhcl Status", "XHCL Status"}),
		Health:        pickString(row, []string{"health"}, []string{"Health", "Health Status"}),
		Mode:          pickString(row, []string{"mode"}, []string{"Mode"}),
	}

	d.TempCelsius = pickFloatUnit(row, []string{"temp_celsius"}, []string{"Temperature (C)", "Temperature (Sensor core) (C)"})
	d.EdgeTempC = pickFloatUnit(row, []string{"edge_temp_c"}, []string{"Temperature (Sensor edge) (C)"})
	d.JunctionTempC = pickFloatUnit(row, []string{"junction_temp_c"}, []string{"Temperature (Sensor junction) (C)"})
	d.MemTempC = pickFloatUnit(row, []string{"mem_temp_c"}, []string{"Temperature (Sensor mem) (C)"})
	d.HotspotTempC = pickFloatUnit(row, []string{"hotspot_temp_c"}, []string{"Temperature (Sensor hotspot) (C)"})

	d.PowerWatts = pickFloatUnit(row, []string{"power_watts"}, []string{"Average Graphics Package Power (W)", "Power (W)", "AvgPwr"})
	d.PowerCapWatts = pickFloatUnit(row, []string{"power_cap_watts"}, []string{"Power Cap (W)", "PwrCap"})
	d.MaxPowerWatts = pickFloatUnit(row, []string{"max_power_watts"}, []string{"Max Power (W)", "Max Package Power (W)", "Max Graphics Package Power (W)"})
	d.VoltageVolts = pickFloatUnit(row, []string{"voltage_volts"}, []string{"Voltage (V)"})
	if d.VoltageVolts == 0 {
		d.VoltageVolts = pickFloatUnit(row, []string{"Voltage (mV)"}) / 1000
	}

	d.HCUUtilPct = pickFloatUnit(row, []string{"hcu_util_pct"}, []string{"HCU Utilization (%)", "HCU%", "HCU use (%)"})
	d.CUUtilPct = pickFloatUnit(row, []string{"cu_util_pct"}, []string{"CU Utilization (%)", "CU%"})
	d.WaveUtilPct = pickFloatUnit(row, []string{"wave_util_pct"}, []string{"Wave Utilization (%)"})
	d.DecUtilPct = pickFloatUnit(row, []string{"dec_util_pct"}, []string{"Dec Utilization (%)", "Dec%"})
	d.EncUtilPct = pickFloatUnit(row, []string{"enc_util_pct"}, []string{"Enc Utilization (%)", "Enc%"})
	d.MemUtilPct = pickFloatUnit(row, []string{"mem_util_pct"}, []string{"Memory Utilization (%)", "HCU memory use (%)"})
	d.VRAMUtilPct = pickFloatUnit(row, []string{"vram_util_pct"}, []string{"VRAM Utilization (%)", "VRAM%"})

	// hy-smi exposes memory in MiB; normalise to bytes.
	d.VRAMTotalBytes = pickMemBytes(row, []string{"vram_total_bytes"}, []string{"Total VRAM (MiB)", "Memory Total (MiB)", "vram Total Memory (MiB)", "vis_vram Total Memory (MiB)"})
	d.VRAMUsedBytes = pickMemBytes(row, []string{"vram_used_bytes"}, []string{"Used VRAM (MiB)", "Memory Used (MiB)", "vram Total Used Memory (MiB)", "vis_vram Total Used Memory (MiB)"})
	d.VRAMAvailBytes = pickMemBytes(row, []string{"vram_available_bytes"}, []string{"Available VRAM (MiB)", "Memory Available (MiB)", "Available memory size (MiB)"})
	d.GTTTotalBytes = pickMemBytes(row, []string{"gtt_total_bytes"}, []string{"Total GTT (MiB)", "gtt Total Memory (MiB)"})
	d.GTTUsedBytes = pickMemBytes(row, []string{"gtt_used_bytes"}, []string{"Used GTT (MiB)", "gtt Total Used Memory (MiB)"})

	d.SCLKMHz = pickFloatUnit(row, []string{"sclk_mhz"}, []string{"SCLK (MHz)", "GFXCLK (MHz)", "sclk clock speed"})
	d.MCLKMHz = pickFloatUnit(row, []string{"mclk_mhz"}, []string{"MCLK (MHz)", "MEMCLK (MHz)", "mclk clock speed"})
	d.SOCCLKMHz = pickFloatUnit(row, []string{"socclk_mhz"}, []string{"SOCCLK (MHz)", "socclk clock speed"})
	d.MaxSCLKMHz = pickFloatUnit(row, []string{"max_sclk_mhz"}, []string{"Max SCLK (MHz)"})
	d.MinSCLKMHz = pickFloatUnit(row, []string{"min_sclk_mhz"}, []string{"Min SCLK (MHz)"})

	d.PCIELinkSpeedGTs = pickFloatUnit(row, []string{"pcie_link_speed_gts"}, []string{"PCIe Link Speed (GT/s)", "pcie clock speed"})
	d.PCIELinkWidth = int(pickInt64(row, []string{"pcie_link_width"}, []string{"PCIe Link Width"}))
	d.PCIEBandwidthMBs = pickFloatUnit(row, []string{"pcie_bandwidth_mbs"}, []string{"PCIe Bandwidth (MB/s)", "Estimated maximum PCIe bandwidth over the last second (MB/s)"})
	if d.PCIELinkWidth == 0 {
		d.PCIELinkWidth = pcieWidth(pickString(row, []string{"pcie clock speed"}))
	}

	d.XHCLBwMBs = pickFloatUnit(row, []string{"xhcl_bw_mbs"}, []string{"xhcl Bandwidth (MB/s)"})
	d.DFBwMBs = pickFloatUnit(row, []string{"df_bw_mbs"}, []string{"DF Bandwidth (MB/s)"})
	d.UMCBwMBs = pickFloatUnit(row, []string{"umc_bw_mbs"}, []string{"UMC Bandwidth (MB/s)"})

	d.RASErrCount = uint64(pickInt64(row, []string{"ras_err_count"}, []string{"RAS Error Count"}))
	d.ECCCorrectable = uint64(pickInt64(row, []string{"ecc_correctable"}, []string{"ECC Correctable Errors"}))
	d.ECCUncorrectable = uint64(pickInt64(row, []string{"ecc_uncorrectable"}, []string{"ECC Uncorrectable Errors"}))
	d.BadPagesRetired = uint64(pickInt64(row, []string{"bad_pages_retired"}, []string{"Retired Pages"}))
	d.BadPagesPending = uint64(pickInt64(row, []string{"bad_pages_pending"}, []string{"Pending Pages"}))

	d.FanRPM = pickFloatUnit(row, []string{"fan_rpm"}, []string{"Fan Speed (RPM)"})
	d.FanPWMPct = pickFloatUnit(row, []string{"fan_pwm_pct"}, []string{"Fan PWM (%)"})

	return d
}

// pickString walks the row through any of the given key paths and
// returns the first non-empty string it finds. "N/A" is treated as
// empty.
func pickString(row map[string]any, paths ...[]string) string {
	for _, path := range paths {
		v := lookup(row, path)
		s := stringFrom(v)
		if s != "" {
			return s
		}
	}
	return ""
}

// pickFloatUnit extracts a float from a value that may be numeric or
// a "157.0W" / "40.0C" / "62.5%" style string. Non-numeric values are
// returned as 0.
func pickFloatUnit(row map[string]any, paths ...[]string) float64 {
	for _, path := range paths {
		v := lookup(row, path)
		if f, ok := toFloat(v); ok {
			return f
		}
	}
	return 0
}

// pickInt64 is pickFloatUnit but int64. Strings are parsed via
// strconv after stripping units.
func pickInt64(row map[string]any, paths ...[]string) int64 {
	for _, path := range paths {
		v := lookup(row, path)
		if f, ok := toFloat(v); ok {
			return int64(f)
		}
	}
	return 0
}

// pickMemBytes reads a memory quantity in MiB and returns bytes.
// Falls back to a raw byte count if the value is already an integer.
func pickMemBytes(row map[string]any, paths ...[]string) uint64 {
	for _, path := range paths {
		v := lookup(row, path)
		if f, ok := toFloat(v); ok {
			return uint64(f * 1024 * 1024)
		}
	}
	return 0
}

// walk navigates a nested map by keys. Missing keys return nil.
func walk(row map[string]any, path []string) any {
	var cur any = row
	for _, key := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[key]
	}
	return cur
}

// lookup supports both nested paths and the historical call style in
// this package where a path contains a list of alternative top-level
// keys, for example []string{"Product Name", "SKU"}.
func lookup(row map[string]any, path []string) any {
	if v := walk(row, path); v != nil {
		return v
	}
	if len(path) <= 1 {
		return nil
	}
	for _, candidate := range path {
		if v, ok := row[candidate]; ok {
			return v
		}
		want := normaliseKey(candidate)
		for key, v := range row {
			if normaliseKey(key) == want {
				return v
			}
		}
	}
	return nil
}

func normaliseKey(key string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(key) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// stringFrom normalises any scalar to a trimmed string. "N/A" / "" → "".
func stringFrom(v any) string {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		if s == "" || strings.EqualFold(s, "n/a") || strings.EqualFold(s, "na") {
			return ""
		}
		return s
	case float64:
		return fmt.Sprintf("%g", t)
	case bool:
		return fmt.Sprintf("%v", t)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	}
}

// toFloat parses a value that may be numeric, a unit string ("157.0W"),
// a percentage ("62.5%"), or a fraction. Returns ok=false for
// anything non-numeric.
func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case string:
		s := strings.TrimSpace(t)
		s = strings.TrimSuffix(s, "W")
		s = strings.TrimSuffix(s, "C")
		s = strings.TrimSuffix(s, "V")
		s = strings.TrimSuffix(s, "MHz")
		s = strings.TrimSuffix(s, "RPM")
		s = strings.TrimSuffix(s, "MB/s")
		s = strings.TrimSuffix(s, "GT/s")
		s = strings.TrimSuffix(s, "%")
		s = strings.TrimSpace(s)
		if s == "" || strings.EqualFold(s, "n/a") || strings.EqualFold(s, "na") {
			return 0, false
		}
		var f float64
		_, err := fmt.Sscanf(s, "%f", &f)
		if err != nil {
			return 0, false
		}
		return f, true
	case nil:
		return 0, false
	default:
		return 0, false
	}
}

var (
	ansiRE   = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
	numberRE = regexp.MustCompile(`[+-]?(?:\d+(?:\.\d*)?|\.\d+)`)
)

// parseText handles the human-readable output emitted by hy-smi
// versions that reject --json with --showallinfo.
func parseText(raw []byte) (Snapshot, error) {
	text := ansiRE.ReplaceAllString(string(raw), "")
	devices := make(map[int]Device)
	snap := Snapshot{}

	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Driver Version:") {
			snap.DriverVersion = strings.TrimSpace(strings.TrimPrefix(line, "Driver Version:"))
			continue
		}
		idx, record, ok := parseHCURecord(line)
		if !ok {
			continue
		}
		d := devices[idx]
		d.Index = idx
		if d.HCU == "" {
			d.HCU = "HCU"
		}
		key, value, ok := splitRecord(record)
		if !ok {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "device id":
			d.DeviceID = strings.TrimSpace(value)
		case "vbios version":
			d.VBIOSVersion = strings.TrimSpace(value)
		case "temperature (sensor edge) (c)":
			d.EdgeTempC = firstFloat(value)
		case "temperature (sensor junction) (c)":
			d.JunctionTempC = firstFloat(value)
		case "temperature (sensor mem) (c)":
			d.MemTempC = firstFloat(value)
		case "temperature (sensor core) (c)":
			d.TempCelsius = firstFloat(value)
		case "mclk clock level":
			d.MCLKMHz = clockMHz(value)
		case "sclk clock level":
			d.SCLKMHz = clockMHz(value)
		case "socclk clock level":
			d.SOCCLKMHz = clockMHz(value)
		case "pcie clock level":
			d.PCIELinkSpeedGTs, d.PCIELinkWidth = parsePCIe(value)
		case "performance level":
			d.PerfLevel = strings.TrimSpace(value)
		case "max graphics package power (w)":
			d.MaxPowerWatts = firstFloat(value)
		case "average graphics package power (w)":
			d.PowerWatts = firstFloat(value)
		case "hcu use (%)":
			d.HCUUtilPct = firstFloat(value)
		case "hcu memory use (%)":
			d.MemUtilPct = firstFloat(value)
		case "serial number":
			d.Serial = strings.TrimSpace(value)
		case "voltage (mv)":
			d.VoltageVolts = firstFloat(value) / 1000
		case "pci bus":
			d.Bus = strings.TrimSpace(strings.SplitN(value, "-->", 2)[0])
		case "card series":
			d.ProductName = strings.TrimSpace(value)
		}
		devices[idx] = d
	}

	if len(devices) == 0 {
		return Snapshot{}, fmt.Errorf("no HCU records found in text output")
	}
	for _, d := range devices {
		if d.EdgeTempC == 0 {
			d.EdgeTempC = d.TempCelsius
		}
		if d.TempCelsius == 0 {
			d.TempCelsius = d.EdgeTempC
		}
		snap.Devices = append(snap.Devices, d)
	}
	sortDevices(snap.Devices)
	return snap, nil
}

func parseHCURecord(line string) (int, string, bool) {
	if !strings.HasPrefix(line, "HCU[") {
		return 0, "", false
	}
	end := strings.IndexByte(line, ']')
	if end < 5 {
		return 0, "", false
	}
	idx, err := strconv.Atoi(line[4:end])
	if err != nil || idx < 0 {
		return 0, "", false
	}
	record := strings.TrimSpace(line[end+1:])
	record = strings.TrimSpace(strings.TrimPrefix(record, ":"))
	return idx, record, true
}

func splitRecord(record string) (string, string, bool) {
	sep := strings.IndexByte(record, ':')
	if sep < 1 {
		return "", "", false
	}
	return strings.TrimSpace(record[:sep]), strings.TrimSpace(record[sep+1:]), true
}

func firstFloat(value string) float64 {
	match := numberRE.FindString(value)
	if match == "" {
		return 0
	}
	f, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0
	}
	return f
}

func clockMHz(value string) float64 {
	start := strings.LastIndex(value, "(")
	if start >= 0 {
		value = value[start+1:]
	}
	return firstFloat(value)
}

func parsePCIe(value string) (float64, int) {
	if start := strings.IndexByte(value, '('); start >= 0 {
		value = value[start+1:]
	}
	return firstFloat(value), pcieWidth(value)
}

func pcieWidth(value string) int {
	value = strings.ToLower(value)
	marker := strings.Index(value, "x")
	if marker < 0 {
		return 0
	}
	return int(firstFloat(value[marker+1:]))
}

func sortDevices(devs []Device) {
	// Insertion sort — device lists are tiny (≤16 cards) and the
	// function is called once per scrape.
	for i := 1; i < len(devs); i++ {
		j := i
		for j > 0 && devs[j].Index < devs[j-1].Index {
			devs[j], devs[j-1] = devs[j-1], devs[j]
			j--
		}
	}
}

func trim(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 200 {
		return s[:200]
	}
	return s
}
