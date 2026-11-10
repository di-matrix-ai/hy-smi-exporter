# Roadmap

This roadmap describes intended community work and is not a commercial
delivery commitment. Scope and dates may change based on real-device access,
driver changes and maintainer capacity.

## v0.1.x - first public release stabilization

- Complete the first verified device, driver, firmware, DTK and vendor-tool
  compatibility matrix.
- Stabilize JSON and text parsing against redacted real-device samples.
- Validate Linux amd64 and arm64 binaries and the multi-architecture image.
- Publish installation, metric semantics and known limitations.
- Complete Prometheus exporter directory submission after the first release.

## Later candidates

- Expand verified device and driver combinations through test evidence.
- Reduce container privileges where the vendor runtime permits it.
- Add a Prometheus Operator `ServiceMonitor` example.
- Evaluate a Grafana dashboard and Helm chart based on community demand.
- Evaluate process-level and virtual-device metrics only after vendor interface,
  security and performance review.

Compatibility work must include a redacted fixture, parser test, metric update,
real-device result and compatibility-matrix change in the same release cycle.
