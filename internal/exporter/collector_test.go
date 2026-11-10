package exporter

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/di-matrix-ai/hy-smi-exporter/internal/hysmi"
)

func TestCollectorMetricContract(t *testing.T) {
	stub := writeExporterStub(t, `[
  {
    "device_index": 0,
    "hcu": "HCU",
    "hcu_util_pct": 25,
		"fan_pwm_pct": 40,
		"sclk_mhz": 1500,
    "ras_err_count": 2,
    "ecc_correctable": 3,
    "ecc_uncorrectable": 4,
    "vram_total_bytes": 68719476736,
    "vram_used_bytes": 536870912
  }
]`)
	families := gather(t, hysmi.New(stub, time.Second))

	assertMetricValue(t, families, "hy_smi_up", 1)
	assertMetricValue(t, families, "hy_smi_devices", 1)
	assertMetricValue(t, families, "hy_smi_hcu_utilization_ratio", 0.25)
	assertMetricValue(t, families, "hy_smi_fan_speed_ratio", 0.40)
	assertMetricValue(t, families, "hy_smi_clock_hertz", 1_500_000_000)
	assertMetricValue(t, families, "hy_smi_ras_errors_total", 2)

	for _, oldName := range []string{
		"hy_smi_devices_total",
		"hy_smi_hcu_utilization_percent",
		"hy_smi_fan_speed_percent",
	} {
		if _, ok := families[oldName]; ok {
			t.Errorf("obsolete metric %s is still exposed", oldName)
		}
	}

	for name, family := range families {
		for _, metric := range family.Metric {
			for _, label := range metric.Label {
				if label.GetName() == "instance" {
					t.Errorf("metric %s contains exporter-generated instance label", name)
				}
			}
		}
	}
}

func TestCollectorReportsSuccessfulZeroDeviceQuery(t *testing.T) {
	stub := writeExporterStub(t, `{}`)
	families := gather(t, hysmi.New(stub, time.Second))
	assertMetricValue(t, families, "hy_smi_up", 1)
	assertMetricValue(t, families, "hy_smi_devices", 0)
}

func TestCollectorReportsBackendFailure(t *testing.T) {
	families := gather(t, hysmi.New(filepath.Join(t.TempDir(), "missing-hy-smi"), time.Second))
	assertMetricValue(t, families, "hy_smi_up", 0)
	if _, ok := families["hy_smi_devices"]; ok {
		t.Error("device count must not be emitted when the backend query fails")
	}
}

func gather(t *testing.T, backend *hysmi.Backend) map[string]*dto.MetricFamily {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	registry := prometheus.NewRegistry()
	registry.MustRegister(NewCollector(backend, logger))
	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	byName := make(map[string]*dto.MetricFamily, len(families))
	for _, family := range families {
		byName[family.GetName()] = family
	}
	return byName
}

func assertMetricValue(t *testing.T, families map[string]*dto.MetricFamily, name string, want float64) {
	t.Helper()
	family, ok := families[name]
	if !ok {
		t.Fatalf("metric %s was not exposed", name)
	}
	if len(family.Metric) == 0 {
		t.Fatalf("metric %s has no samples", name)
	}
	metric := family.Metric[0]
	var got float64
	switch family.GetType() {
	case dto.MetricType_COUNTER:
		got = metric.GetCounter().GetValue()
	default:
		got = metric.GetGauge().GetValue()
	}
	if got != want {
		t.Errorf("metric %s: want %v, got %v", name, want, got)
	}
}

func writeExporterStub(t *testing.T, output string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "hy-smi")
	body := "#!/bin/sh\n/bin/cat <<'EOF'\n" + output + "\nEOF\n"
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	return path
}
