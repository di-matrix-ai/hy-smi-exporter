package config

import "testing"

func TestDefaultDisablesOneSecondStats(t *testing.T) {
	if Default().CollectOneSecondStats {
		t.Fatal("one-second statistics must be disabled by default")
	}
}

func TestFromEnvEnablesOneSecondStats(t *testing.T) {
	t.Setenv("HY_SMI_EXPORTER_COLLECT_1S_STATS", "true")
	if !FromEnv().CollectOneSecondStats {
		t.Fatal("expected HY_SMI_EXPORTER_COLLECT_1S_STATS=true to enable collection")
	}
}

func TestFromEnvIgnoresInvalidOneSecondStatsValue(t *testing.T) {
	t.Setenv("HY_SMI_EXPORTER_COLLECT_1S_STATS", "not-a-bool")
	if FromEnv().CollectOneSecondStats {
		t.Fatal("invalid boolean value must preserve the default")
	}
}
