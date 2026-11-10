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

## Release compliance checklist

- [ ] `go.mod` and `go.sum` match the release commit.
- [ ] CI has completed `go mod verify`.
- [ ] License identifiers above have been checked against the exact downloaded
      module versions.
- [ ] No dependency introduces a license incompatible with Apache-2.0
      distribution.
- [ ] The release archive includes `LICENSE`, `NOTICE` and this file.
- [ ] The final Go module dependency inventory is archived with the release evidence.
- [ ] Dingmao's legal or open-source compliance owner has approved publication.

Compliance owner: `TBD`  
Review date: `TBD`  
Result: `TBD - approved / approved with conditions / rejected`

The licenses shown here are based on upstream project metadata. Before the
first public release, the compliance owner must verify them from the exact
module source used by CI and retain the dependency inventory or equivalent evidence.
