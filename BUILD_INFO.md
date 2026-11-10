# Build and Release Traceability

This file is the release build record for `hy-smi-exporter`. The release owner
must replace every `TBD` value for a tagged release. Development branches may
retain this template.

## Release identity

| Item | Value |
| --- | --- |
| Project | `hy-smi-exporter` |
| Version / tag | `TBD` |
| Git repository | `TBD` |
| Git commit SHA | `TBD` |
| Source archive SHA-256 | `TBD` |
| Build date (UTC) | `TBD` |
| Release owner | `TBD` |

## Build environment

| Item | Value |
| --- | --- |
| Go version | `TBD` |
| GoReleaser version | `TBD` |
| CI provider and run URL | `TBD` |
| Runner operating system | `TBD` |
| Build command | `goreleaser release --clean` |
| CGO | `CGO_ENABLED=0` |

## Release artifacts

| Artifact | Architecture | SHA-256 / digest | Verification result |
| --- | --- | --- | --- |
| Linux archive | `amd64` | `TBD` | `TBD` |
| Linux archive | `arm64` | `TBD` | `TBD` |
| GHCR image | multi-arch | 首版不发布 | 不适用 |
| SBOM | SPDX / OCI | `TBD` | `TBD` |
| Source archive | source | `TBD` | `TBD` |

## Verification evidence

| Check | Evidence | Result |
| --- | --- | --- |
| `go mod verify` | CI run URL / log | `TBD` |
| Formatting | CI run URL / log | `TBD` |
| `go test -race ./...` | CI run URL / log | `TBD` |
| `go vet ./...` | CI run URL / log | `TBD` |
| Linux amd64 build | CI run URL / log | `TBD` |
| Linux arm64 build | CI run URL / log | `TBD` |
| Container multi-arch build | 首版不执行 | 不适用 |
| Real-device acceptance | `release-materials/05-rc-acceptance.md` | `TBD` |

Release artifacts must be generated from the commit recorded above. Do not
manually replace a binary after checksums or image digests have been recorded.
