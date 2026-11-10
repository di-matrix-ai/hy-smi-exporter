package hysmi

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestParseMapForm covers the map-of-cards shape emitted by the
// common hy-smi --json mode: a "system" block plus "card_0" ...
// "card_N" entries with human-readable keys like
// "Temperature (Sensor edge) (C)" and unit-bearing strings like
// "157.0W".
func TestParseMapForm(t *testing.T) {
	raw := loadFixture(t, "two_cards.json")
	snap, err := parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if snap.DriverVersion != "6.2.22" {
		t.Errorf("driver version: want 6.2.22, got %q", snap.DriverVersion)
	}
	if snap.HostName != "dcu-node-01" {
		t.Errorf("hostname: want dcu-node-01, got %q", snap.HostName)
	}
	if len(snap.Devices) != 2 {
		t.Fatalf("devices: want 2, got %d", len(snap.Devices))
	}
	if snap.Devices[0].Index != 0 || snap.Devices[1].Index != 1 {
		t.Errorf("indices: want 0,1; got %d,%d", snap.Devices[0].Index, snap.Devices[1].Index)
	}

	d0 := snap.Devices[0]
	if d0.ProductName != "Hygon C-3000" {
		t.Errorf("product: got %q", d0.ProductName)
	}
	if d0.TempCelsius != 40.0 {
		t.Errorf("temp: want 40.0, got %v", d0.TempCelsius)
	}
	if d0.EdgeTempC != 40.0 {
		t.Errorf("edge temp: want 40.0, got %v", d0.EdgeTempC)
	}
	if d0.JunctionTempC != 45.0 {
		t.Errorf("junction temp: want 45.0, got %v", d0.JunctionTempC)
	}
	if d0.PowerWatts != 157.0 {
		t.Errorf("power: want 157.0, got %v", d0.PowerWatts)
	}
	if d0.PowerCapWatts != 800.0 {
		t.Errorf("powercap: want 800.0, got %v", d0.PowerCapWatts)
	}
	if d0.VoltageVolts != 0.75 {
		t.Errorf("voltage: want 0.75, got %v", d0.VoltageVolts)
	}
	if d0.VRAMTotalBytes != 65536*1024*1024 {
		t.Errorf("vram total: want %d, got %d", 65536*1024*1024, d0.VRAMTotalBytes)
	}
	if d0.VRAMUsedBytes != 512*1024*1024 {
		t.Errorf("vram used: want %d, got %d", 512*1024*1024, d0.VRAMUsedBytes)
	}
	if d0.SCLKMHz != 1500.0 {
		t.Errorf("sclk: want 1500.0, got %v", d0.SCLKMHz)
	}
	if d0.PCIELinkWidth != 16 {
		t.Errorf("pcie width: want 16, got %d", d0.PCIELinkWidth)
	}
	if d0.PCIELinkSpeedGTs != 32.0 {
		t.Errorf("pcie speed: want 32.0, got %v", d0.PCIELinkSpeedGTs)
	}
	if d0.Health != "Healthy" {
		t.Errorf("health: got %q", d0.Health)
	}
	if d0.FanRPM != 3000 {
		t.Errorf("fan rpm: want 3000, got %v", d0.FanRPM)
	}
	if d0.FanPWMPct != 40 {
		t.Errorf("fan pwm pct: want 40, got %v", d0.FanPWMPct)
	}
}

// TestParseArrayForm covers the array-of-rows shape. The first row
// is the system block (no Device Index); the remaining rows are
// devices.
func TestParseArrayForm(t *testing.T) {
	raw := loadFixture(t, "array_form.json")
	snap, err := parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(snap.Devices) != 1 {
		t.Fatalf("devices: want 1, got %d", len(snap.Devices))
	}
	d := snap.Devices[0]
	if d.HCUUtilPct != 100.0 {
		t.Errorf("hcu util: want 100.0, got %v", d.HCUUtilPct)
	}
	if d.VRAMUtilPct != 80.0 {
		t.Errorf("vram util: want 80.0, got %v", d.VRAMUtilPct)
	}
	if d.VRAMUsedBytes != 52428*1024*1024 {
		t.Errorf("vram used: want %d, got %d", 52428*1024*1024, d.VRAMUsedBytes)
	}
	if d.RASErrCount != 2 {
		t.Errorf("ras count: want 2, got %d", d.RASErrCount)
	}
	if d.ECCCorrectable != 2 {
		t.Errorf("ecc correctable: want 2, got %d", d.ECCCorrectable)
	}
	if d.BadPagesRetired != 1 {
		t.Errorf("bad pages retired: want 1, got %d", d.BadPagesRetired)
	}
	if d.BadPagesPending != 1 {
		t.Errorf("bad pages pending: want 1, got %d", d.BadPagesPending)
	}
}

// TestParseInvalidJSON checks that a malformed document returns an
// error rather than panicking.
func TestParseInvalidJSON(t *testing.T) {
	_, err := parse([]byte("{not json"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestParseEmpty ensures an empty object yields a zero-device
// snapshot, not a crash.
func TestParseEmpty(t *testing.T) {
	snap, err := parse([]byte("{}"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(snap.Devices) != 0 {
		t.Errorf("expected 0 devices, got %d", len(snap.Devices))
	}
}

func TestParseMapCardKeyAddsDeviceIndex(t *testing.T) {
	raw := []byte(`{"card1":{"Device ID":"0x1001"},"card0":{"Device ID":"0x1000"}}`)
	snap, err := parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(snap.Devices) != 2 {
		t.Fatalf("devices: want 2, got %d", len(snap.Devices))
	}
	if snap.Devices[0].Index != 0 || snap.Devices[1].Index != 1 {
		t.Fatalf("indices: want 0,1; got %d,%d", snap.Devices[0].Index, snap.Devices[1].Index)
	}
}

func TestParseNormalizedByteFields(t *testing.T) {
	raw := []byte(`[{"device_index":0,"vram_total_bytes":68719476736,"vram_used_bytes":536870912}]`)
	snap, err := parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(snap.Devices) != 1 {
		t.Fatalf("devices: want 1, got %d", len(snap.Devices))
	}
	if got := snap.Devices[0].VRAMTotalBytes; got != 68719476736 {
		t.Fatalf("vram total bytes: want 68719476736, got %d", got)
	}
	if got := snap.Devices[0].VRAMUsedBytes; got != 536870912 {
		t.Fatalf("vram used bytes: want 536870912, got %d", got)
	}
}

func TestParseText(t *testing.T) {
	raw := []byte("\x1b[31mWarning\x1b[0m\n" +
		"Driver Version: 6.3.28-V1.3.0b\n" +
		"HCU[1]\t\t: Device ID: 0x6370\n" +
		"HCU[1]\t\t: Temperature (Sensor edge) (C): 51.0\n" +
		"HCU[1]\t\t: Temperature (Sensor core) (C): 49.0\n" +
		"HCU[1]\t\t: mclk clock level: 0 (1800Mhz)\n" +
		"HCU[1]\t\t: sclk clock level: 1 (600Mhz)\n" +
		"HCU[1]\t\t: socclk clock level: 0 (1079Mhz)\n" +
		"HCU[1]\t\t: pcie clock level: 2 (32.0GT/s, x16 1100Mhz)\n" +
		"HCU[1]\t\t: Performance Level: auto\n" +
		"HCU[1]\t\t: Max Graphics Package Power (W): 400.0\n" +
		"HCU[1]\t\t: Average Graphics Package Power (W): 53.0\n" +
		"HCU[1]\t\t: HCU use (%): 0.0\n" +
		"HCU[1]\t\t: HCU memory use (%): 0\n" +
		"HCU[1]\t\t: Serial Number: TEST-SERIAL-0001\n" +
		"HCU[1]\t\t: Voltage (mV): 743\n" +
		"HCU[1]\t\t: PCI Bus: 0000:97:00.0 --> OAM ID: 0\n" +
		"HCU[1]\t\t: Card Series:\t\t BW100\n")
	snap, err := parseText(raw)
	if err != nil {
		t.Fatalf("parseText: %v", err)
	}
	if snap.DriverVersion != "6.3.28-V1.3.0b" || len(snap.Devices) != 1 {
		t.Fatalf("snapshot: driver=%q devices=%d", snap.DriverVersion, len(snap.Devices))
	}
	d := snap.Devices[0]
	if d.Index != 1 || d.DeviceID != "0x6370" || d.ProductName != "BW100" {
		t.Errorf("identity: %+v", d)
	}
	if d.EdgeTempC != 51 || d.TempCelsius != 49 || d.SCLKMHz != 600 || d.MCLKMHz != 1800 || d.SOCCLKMHz != 1079 {
		t.Errorf("thermal/clocks: %+v", d)
	}
	if d.PCIELinkSpeedGTs != 32 || d.PCIELinkWidth != 16 {
		t.Errorf("pcie: speed=%v width=%d", d.PCIELinkSpeedGTs, d.PCIELinkWidth)
	}
	if d.PowerWatts != 53 || d.MaxPowerWatts != 400 || d.VoltageVolts != 0.743 {
		t.Errorf("power/voltage: power=%v max=%v voltage=%v", d.PowerWatts, d.MaxPowerWatts, d.VoltageVolts)
	}
}

func TestBackendWithTextFallback(t *testing.T) {
	stub := writeStubScript(t, "#!/bin/sh\n/bin/cat <<'EOF'\n"+
		"Driver Version: 6.3.28-V1.3.0b\n"+
		"HCU[0] : Device ID: 0x6370\n"+
		"HCU[0] : Temperature (Sensor edge) (C): 50.0\n"+
		"HCU[0] : Average Graphics Package Power (W): 20.0\n"+
		"EOF\n")
	b := New(stub, 0)
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	snap, err := b.Snapshot(context.Background(), logger)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(snap.Devices) != 1 || snap.Devices[0].PowerWatts != 20 {
		t.Fatalf("fallback snapshot: %+v", snap)
	}
}

// TestBackendBinaryNotFound exercises the code path where hy-smi
// is not on $PATH. We temporarily shrink $PATH to /dev/null so
// exec.LookPath fails. The snapshot must be empty and the missing
// warning latch must engage without spamming.
func TestBackendBinaryNotFound(t *testing.T) {
	t.Setenv("PATH", "/dev/null")
	b := New("", 0)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	snap, err := b.Snapshot(context.Background(), logger)
	if err == nil {
		t.Fatal("expected an error when binary is missing")
	}
	if len(snap.Devices) != 0 {
		t.Fatalf("expected 0 devices when binary missing, got %d", len(snap.Devices))
	}
	if !b.missingLogged.Load() {
		t.Error("missingLogged flag should be set after first miss")
	}
}

func TestBackendBinaryNotFoundConcurrent(t *testing.T) {
	t.Setenv("PATH", "/dev/null")
	b := New("", 0)
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = b.Snapshot(context.Background(), logger)
		}()
	}
	wg.Wait()

	if !b.missingLogged.Load() {
		t.Error("missingLogged flag should be set after concurrent misses")
	}
}

// TestBackendWithStubBinary wires a fake hy-smi shell script that
// prints the two-card fixture and verifies the Backend returns the
// parsed snapshot.
func TestBackendWithStubBinary(t *testing.T) {
	stub := writeStub(t, "two_cards.json")
	b := New(stub, 0)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	snap, err := b.Snapshot(context.Background(), logger)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(snap.Devices) != 2 {
		t.Fatalf("devices: want 2, got %d", len(snap.Devices))
	}
	if snap.DriverVersion != "6.2.22" {
		t.Errorf("driver: got %q", snap.DriverVersion)
	}
}

// TestBackendWithFailingStub verifies a non-zero exit from hy-smi
// is returned to the caller and produces an empty snapshot.
func TestBackendWithFailingStub(t *testing.T) {
	stub := writeStubScript(t, "#!/bin/sh\necho boom >&2\nexit 1\n")
	b := New(stub, 0)
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	snap, err := b.Snapshot(context.Background(), logger)
	if err == nil {
		t.Fatal("expected an error from a failing hy-smi command")
	}
	if len(snap.Devices) != 0 {
		t.Fatalf("expected 0 devices on non-zero exit, got %d", len(snap.Devices))
	}
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func writeStub(t *testing.T, fixture string) string {
	t.Helper()
	data := loadFixture(t, fixture)
	return writeStubScript(t, "#!/bin/sh\n/bin/cat <<'EOF'\n"+string(data)+"\nEOF\n")
}

func writeStubScript(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "hy-smi")
	if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	return p
}
