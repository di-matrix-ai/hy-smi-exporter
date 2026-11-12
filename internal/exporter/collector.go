// Package exporter wires the hy-smi backend behind a single
// prometheus.Collector. Each scrape runs the configured read-only hy-smi
// queries and merges their results.
// Metrics are emitted with the `hy_smi_` prefix so this exporter
// coexists cleanly with node_exporter / dcgm-exporter / amdgpu_exporter.
package exporter

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/di-matrix-ai/hy-smi-exporter/internal/hysmi"
)

// Collector implements prometheus.Collector. Each scrape calls
// Backend.Snapshot once and emits one metric series per device.
type Collector struct {
	smi                   *hysmi.Backend
	logger                *slog.Logger
	descs                 collectorDescs
	collectOneSecondStats bool
}

// NewCollector builds a Collector over the given backend.
func NewCollector(smi *hysmi.Backend, logger *slog.Logger) *Collector {
	return &Collector{
		smi:                   smi,
		logger:                logger,
		descs:                 newCollectorDescs(),
		collectOneSecondStats: smi.CollectOneSecondStats(),
	}
}

// Describe implements prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.descs.describe(ch)
}

// Collect implements prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	snap, err := c.smi.Snapshot(context.Background(), c.logger)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.descs.up, prometheus.GaugeValue, 0)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.descs.up, prometheus.GaugeValue, 1)

	ch <- prometheus.MustNewConstMetric(
		c.descs.devices, prometheus.GaugeValue, float64(len(snap.Devices)),
	)

	if len(snap.Devices) == 0 {
		return
	}

	for _, d := range snap.Devices {
		c.emitDevice(ch, d, snap.DriverVersion)
	}
}

// emitDevice emits every metric for a single HCU device. Nil/zero
// values are still emitted — hy-smi reports 0 for absent sensors, not
// missing, and ops dashboards want a consistent series inventory.
func (c *Collector) emitDevice(ch chan<- prometheus.Metric, d hysmi.Device, driverVersion string) {
	gpu := strconv.Itoa(d.Index)
	labels := []string{gpu}

	// hy_smi_info — identity gauge, always 1.
	ch <- prometheus.MustNewConstMetric(
		c.descs.info, prometheus.GaugeValue, 1.0,
		gpu,
		d.HCU, d.DeviceID, d.UniqueID, d.Serial,
		d.ProductName, d.PartNumber, d.VBIOSVersion, driverVersion,
		d.Bus, d.Health, d.Mode,
	)

	// Temperature — emit every sensor we saw a non-zero value on;
	// emit a single temp_celsius gauge too for cards that only
	// populate "Temperature (C)".
	emitTemp := func(sensor string, v float64) {
		ch <- prometheus.MustNewConstMetric(
			c.descs.temperature, prometheus.GaugeValue, v,
			gpu, sensor,
		)
	}
	if d.EdgeTempC > 0 {
		emitTemp("edge", d.EdgeTempC)
	} else if d.TempCelsius > 0 {
		emitTemp("edge", d.TempCelsius)
	}
	if d.JunctionTempC > 0 {
		emitTemp("junction", d.JunctionTempC)
	}
	if d.MemTempC > 0 {
		emitTemp("mem", d.MemTempC)
	}
	if d.HotspotTempC > 0 {
		emitTemp("hotspot", d.HotspotTempC)
	}

	// Power.
	ch <- prometheus.MustNewConstMetric(c.descs.power, prometheus.GaugeValue, d.PowerWatts, labels...)
	if d.PowerCapWatts > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.powerCap, prometheus.GaugeValue, d.PowerCapWatts, labels...)
	}
	if d.MaxPowerWatts > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.powerMax, prometheus.GaugeValue, d.MaxPowerWatts, labels...)
	}
	if d.VoltageVolts > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.voltage, prometheus.GaugeValue, d.VoltageVolts, gpu, "core")
	}

	// Utilization values are based on one-second statistical queries on some
	// devices. Do not expose zero-valued placeholders when collection is off.
	if c.collectOneSecondStats {
		ch <- prometheus.MustNewConstMetric(c.descs.hcuUtil, prometheus.GaugeValue, d.HCUUtilPct/100, labels...)
		ch <- prometheus.MustNewConstMetric(c.descs.cuUtil, prometheus.GaugeValue, d.CUUtilPct/100, labels...)
		if d.WaveUtilPct > 0 {
			ch <- prometheus.MustNewConstMetric(c.descs.waveUtil, prometheus.GaugeValue, d.WaveUtilPct/100, labels...)
		}
		if d.DecUtilPct > 0 {
			ch <- prometheus.MustNewConstMetric(c.descs.decUtil, prometheus.GaugeValue, d.DecUtilPct/100, labels...)
		}
		if d.EncUtilPct > 0 {
			ch <- prometheus.MustNewConstMetric(c.descs.encUtil, prometheus.GaugeValue, d.EncUtilPct/100, labels...)
		}
		if d.MemUtilPct > 0 {
			ch <- prometheus.MustNewConstMetric(c.descs.memUtil, prometheus.GaugeValue, d.MemUtilPct/100, labels...)
		}
		if d.VRAMUtilPct > 0 {
			ch <- prometheus.MustNewConstMetric(c.descs.vramUtil, prometheus.GaugeValue, d.VRAMUtilPct/100, labels...)
		}
	}

	// Memory (bytes).
	ch <- prometheus.MustNewConstMetric(c.descs.vramTotal, prometheus.GaugeValue, float64(d.VRAMTotalBytes), labels...)
	ch <- prometheus.MustNewConstMetric(c.descs.vramUsed, prometheus.GaugeValue, float64(d.VRAMUsedBytes), labels...)
	if d.VRAMAvailBytes > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.vramAvail, prometheus.GaugeValue, float64(d.VRAMAvailBytes), labels...)
	}
	if d.GTTTotalBytes > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.gttTotal, prometheus.GaugeValue, float64(d.GTTTotalBytes), labels...)
	}
	if d.GTTUsedBytes > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.gttUsed, prometheus.GaugeValue, float64(d.GTTUsedBytes), labels...)
	}

	// Clocks (Hz). The vendor reports MHz, so normalize to the Prometheus base
	// unit when emitting samples.
	emitClock := func(domain string, v float64) {
		if v <= 0 {
			return
		}
		ch <- prometheus.MustNewConstMetric(c.descs.clock, prometheus.GaugeValue, v*1_000_000, gpu, domain)
	}
	emitClock("sclk", d.SCLKMHz)
	emitClock("mclk", d.MCLKMHz)
	emitClock("socclk", d.SOCCLKMHz)
	if d.MaxSCLKMHz > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.clockMax, prometheus.GaugeValue, d.MaxSCLKMHz*1_000_000, gpu, "sclk")
	}
	if d.MinSCLKMHz > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.clockMin, prometheus.GaugeValue, d.MinSCLKMHz*1_000_000, gpu, "sclk")
	}

	// PCIe.
	if d.PCIELinkSpeedGTs > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.pcieSpeed, prometheus.GaugeValue, d.PCIELinkSpeedGTs, labels...)
	}
	if d.PCIELinkWidth > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.pcieWidth, prometheus.GaugeValue, float64(d.PCIELinkWidth), labels...)
	}
	if c.collectOneSecondStats && d.PCIEBandwidthMBs > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.pcieBw, prometheus.GaugeValue, d.PCIEBandwidthMBs, labels...)
	}

	// Interconnect.
	if c.collectOneSecondStats && d.XHCLBwMBs > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.xhclBw, prometheus.GaugeValue, d.XHCLBwMBs, labels...)
	}
	if c.collectOneSecondStats && d.DFBwMBs > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.dfBw, prometheus.GaugeValue, d.DFBwMBs, labels...)
	}
	if c.collectOneSecondStats && d.UMCBwMBs > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.umcBw, prometheus.GaugeValue, d.UMCBwMBs, labels...)
	}
	if d.XHCLStatus != "" {
		ch <- prometheus.MustNewConstMetric(c.descs.xhclStatus, prometheus.GaugeValue, 1.0, gpu, d.XHCLStatus)
	}

	// Health / RAS / ECC.
	if d.Health != "" {
		ch <- prometheus.MustNewConstMetric(c.descs.health, prometheus.GaugeValue, 1.0, gpu, d.Health)
	}
	// Counters must be emitted on every scrape, even when 0, so
	// Prometheus doesn't interpret a missing sample as a reset.
	ch <- prometheus.MustNewConstMetric(c.descs.rasErrors, prometheus.CounterValue, float64(d.RASErrCount), labels...)
	ch <- prometheus.MustNewConstMetric(c.descs.eccCorrectable, prometheus.CounterValue, float64(d.ECCCorrectable), labels...)
	ch <- prometheus.MustNewConstMetric(c.descs.eccUncorrectable, prometheus.CounterValue, float64(d.ECCUncorrectable), labels...)
	ch <- prometheus.MustNewConstMetric(c.descs.badPagesRetired, prometheus.CounterValue, float64(d.BadPagesRetired), labels...)
	if d.BadPagesPending > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.badPagesPending, prometheus.GaugeValue, float64(d.BadPagesPending), labels...)
	}

	// Fan.
	if d.FanRPM > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.fanRPM, prometheus.GaugeValue, d.FanRPM, gpu, "0")
	}
	if d.FanPWMPct > 0 {
		ch <- prometheus.MustNewConstMetric(c.descs.fanRatio, prometheus.GaugeValue, d.FanPWMPct/100, gpu, "0")
	}

	// Performance level as a {level=...}=1 info gauge.
	if d.PerfLevel != "" {
		ch <- prometheus.MustNewConstMetric(c.descs.perfLevel, prometheus.GaugeValue, 1.0, gpu, d.PerfLevel)
	}
	if d.BoostLevel != "" {
		ch <- prometheus.MustNewConstMetric(c.descs.boostLevel, prometheus.GaugeValue, 1.0, gpu, d.BoostLevel)
	}
}
