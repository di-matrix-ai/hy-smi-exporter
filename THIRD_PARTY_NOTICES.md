# Third-Party Dependency Notices

`hy-smi-exporter` is licensed under Apache License 2.0. The Go dependencies
below are declared in `go.mod`. This inventory is a release compliance record,
not a replacement for the license text distributed by each upstream project.

| Module | Version | Usage | Upstream license |
| --- | --- | --- | --- |
| `github.com/prometheus/client_golang` | `v1.20.5` | Direct dependency; Prometheus registry and HTTP exposition | Apache-2.0 |
| `github.com/beorn7/perks` | `v1.0.1` | Indirect dependency | MIT |
| `github.com/cespare/xxhash/v2` | `v2.3.0` | Indirect dependency | MIT |
| `github.com/klauspost/compress` | `v1.17.9` | Indirect dependency | Apache-2.0 |
| `github.com/munnerz/goautoneg` | `v0.0.0-20191010083416-a7dc8b61c822` | Indirect dependency | BSD-3-Clause |
| `github.com/prometheus/client_model` | `v0.6.1` | Indirect dependency | Apache-2.0 |
| `github.com/prometheus/common` | `v0.55.0` | Indirect dependency | Apache-2.0 |
| `github.com/prometheus/procfs` | `v0.15.1` | Indirect dependency | Apache-2.0 |
| `golang.org/x/sys` | `v0.22.0` | Indirect dependency | BSD-3-Clause |
| `google.golang.org/protobuf` | `v1.34.2` | Indirect dependency | BSD-3-Clause |

## Release verification

The release workflow verifies the module checksum database with `go mod
verify`. The release archive includes `LICENSE`, `NOTICE` and this file. The
licenses shown here are based on the upstream project metadata for the exact
versions declared by `go.mod` and `go.sum`.

Internal release approval and dependency-review evidence are retained by
Shanghai Ding Mao Information Technology Co., Ltd. and are not included in
this public notice file.
