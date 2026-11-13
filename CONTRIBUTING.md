# Contributing to hy-smi-exporter

Thank you for helping improve hy-smi-exporter. Contributions are welcome for
parsers, metrics, documentation, deployment examples, tests, and dashboards.

## Before opening an issue

Please search existing issues first. For parser or metric problems, include:

- exporter version and commit;
- Go, Linux distribution, kernel, driver, and `hy-smi` versions;
- the command-line options used;
- relevant logs and redacted `hy-smi --json`, text fallback, and affected
  subcommand samples.

Do not include serial numbers, hostnames, PCI addresses, or other sensitive
values unless they have been reviewed and redacted.

## Development

Requirements are Go 1.23 or newer and a working Go toolchain. The test suite
does not require a real Hygon device; it uses fixtures and stub commands.

```bash
make fmt
make test
make vet
make build
```

Changes should include focused tests when behavior changes. Keep metric names,
types, units, and label cardinality stable unless the change is intentional and
documented in `CHANGELOG.md`.

## Pull requests

Pull requests should explain the user-visible change and its compatibility
impact. Keep unrelated refactors out of the same change. All CI checks must pass,
and new metrics should be documented in `README.md` and `metric.md`.

By submitting a contribution, you agree that it may be distributed under the
Apache License 2.0 included in this repository.
