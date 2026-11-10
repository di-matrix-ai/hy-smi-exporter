# Architecture

## Runtime flow

1. Prometheus requests the configured metrics endpoint.
2. The collector invokes the configured `hy-smi` or `hcu-smi` executable with
   a timeout.
3. The backend first attempts `-a --json` and falls back to `-a` text output
   when the JSON command is unavailable or cannot be parsed.
4. If the main query lacks memory data, the backend performs the supported
   read-only memory queries.
5. Parsed device snapshots are converted to `hy_smi_*` metrics and returned to
   Prometheus.

```text
Prometheus
    |
    v
HTTP /metrics -> Collector -> Backend -> hy-smi / hcu-smi -> driver and device
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
| Deployment examples | `examples/` | Prometheus and systemd examples; container files are future adaptation drafts |

## Security boundary

The exporter invokes read-only query commands and does not load drivers,
update firmware, reset devices, inject RAS faults or change device settings.
The vendor tool and its runtime libraries are not embedded in the project or
container image. Their installation and licensing remain the responsibility of
the deployment environment.

## Compatibility boundary

Parser unit tests protect known output shapes, but they do not replace
real-device verification. A hardware, driver, firmware, DTK or vendor-tool
upgrade may change fields, units, types or enumerations. Only combinations in
`docs/compatibility.md` marked as verified are part of the published support
statement.
