# Architecture

## Runtime flow

1. Prometheus requests the configured metrics endpoint.
2. The collector invokes the configured `hy-smi` executable with
   a timeout.
3. The backend first attempts the lightweight `hy-smi --json` summary and
   falls back to `hy-smi` text output when JSON is unavailable or cannot be
   parsed.
4. Missing Unique ID, health and memory fields are enriched through the
   read-only `--showuniqueid`, `--healthcheck` and memory queries.
5. Detailed one-second utilization and bandwidth statistics are disabled by
   default. When enabled, the main query uses `-a --json` with `-a` text
   fallback.
6. Parsed device snapshots are converted to `hy_smi_*` metrics and returned to
   Prometheus. Collection failures emit `hy_smi_up 0` instead of being reported
   as a healthy node with zero devices.

```text
Prometheus
    |
    v
HTTP /metrics -> Collector -> Backend -> hy-smi -> driver and device
                     |          |
                     |          +-> JSON parser / text fallback
                     +-> Prometheus metric descriptors
```

## Main components

| Component | Location | Responsibility |
| --- | --- | --- |
| Program entry | `cmd/hy-smi-exporter/` | Flags, HTTP server, logging and shutdown |
| Configuration | `internal/config/` | Flags and environment defaults |
| Vendor command backend | `internal/hysmi/` | Command execution, timeout and parsing |
| Prometheus collector | `internal/exporter/` | Metric descriptions and emission |
| Deployment examples | `examples/` | Prometheus and systemd examples |

## Security boundary

The exporter invokes read-only query commands and does not load drivers,
update firmware, reset devices, inject RAS faults or change device settings.
The vendor tool and its runtime libraries are not embedded in the project.
Their installation and licensing remain the responsibility of the deployment
environment.

## Compatibility boundary

Parser unit tests protect known output shapes, but they do not replace
real-device verification. The exporter supports devices through the common
`hy-smi` command interface rather than hardware-specific access.
A device is in scope when its vendor driver supports it and the command output
remains compatible with the current parser. `docs/compatibility.md` separately
records interface-based support and combinations verified on real devices.
