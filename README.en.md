# hy-smi-exporter

Prometheus exporter for Hygon DCU / HCU accelerators. It runs the read-only
`hy-smi` command, parses JSON when available, falls back to text output for
older drivers, and exposes per-device metrics under the `hy_smi_` prefix.

## Quick start

Build and test with Go 1.23 or newer:

```bash
make build
make test
```

Run it on a host where `hy-smi` or `hcu-smi` is available:

```bash
./bin/hy-smi-exporter \
  --web.listen-address=:9584 \
  --hy-smi.path=/usr/bin/hy-smi
```

Metrics are available at <http://localhost:9584/metrics> and the liveness
endpoint is <http://localhost:9584/healthz>.

`hy_smi_up` reports whether the most recent vendor-tool collection succeeded.
Prometheus supplies the target `instance` label; the exporter does not create
its own `instance` label.

Tagged releases publish Linux `amd64` and `arm64` archives. The first release
does not publish or support a container image.

Driver compatibility and installation details are documented in
[`docs/compatibility.md`](docs/compatibility.md) and
[`docs/installation.md`](docs/installation.md). The compatibility matrix
distinguishes vendor driver support from combinations verified with this
exporter.

Architecture, support boundaries, dependency notices, build traceability and
release approval forms are documented in [`ARCHITECTURE.md`](ARCHITECTURE.md),
[`SUPPORT.md`](SUPPORT.md),
[`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md),
[`BUILD_INFO.md`](BUILD_INFO.md) and
[`release-materials/`](release-materials/00-%E5%8F%91%E5%B8%83%E6%9D%90%E6%96%99%E6%80%BB%E7%9B%AE%E5%BD%95.md).
Planned community work is tracked in [`ROADMAP.md`](ROADMAP.md).

## Requirements

- Go 1.23+ to build.
- A driver, firmware, DTK, and `hy-smi` / `hcu-smi` combination listed in the
  verified compatibility matrix.
- A host installation of `hy-smi` or `hcu-smi`; container deployment is outside
  the compatibility scope of `v0.1.0`.

The exporter invokes only read-only commands. It never uses `--set*`,
`--rasinject`, `--load*driver`, or `--reset*` options.

## Deployment examples

See the [`examples/`](examples/) directory for systemd, environment variables,
and Prometheus discovery examples. Container deployment files are not included
because container scenarios are outside the `v0.1.0` support scope.
The Chinese README contains the complete option and metric reference:
[`README.md`](README.md).

## Contributing and security

Read [`CONTRIBUTING.md`](CONTRIBUTING.md) before opening a pull request. Report
vulnerabilities privately using [`SECURITY.md`](SECURITY.md), not a public issue.

## License

Apache License 2.0. See [`LICENSE`](LICENSE).
