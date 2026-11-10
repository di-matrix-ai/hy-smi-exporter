# Build and Release Traceability

Official `hy-smi-exporter` binaries are built by the GitHub Actions release
workflow from an annotated version tag. Release identity and evidence are
resolved from immutable Git and GitHub records instead of being copied into
this source file.

## Build definition

| Item | Value |
| --- | --- |
| Repository | `https://github.com/di-matrix-ai/hy-smi-exporter` |
| Workflow | `.github/workflows/release.yml` |
| Build tool | GoReleaser v2 |
| Build command | `goreleaser release --clean` |
| Supported artifacts | Linux `amd64` and Linux `arm64` archives |
| CGO | `CGO_ENABLED=0` |
| Checksums | `checksums.txt` attached to the same GitHub Release |
| Release state | Created as draft; published manually after approval |

## Verification procedure

For a release tag such as `v0.1.0`:

1. Resolve the source commit with `git rev-list -n 1 v0.1.0`.
2. Open the workflow run for that tag and verify dependency, format, race,
   vet, native build and Linux cross-build checks.
3. Download every draft Release asset and verify it against `checksums.txt`.
4. Confirm the archive contains the license, notices, README files, metric,
   installation and compatibility documentation.
5. Record the tag, commit, workflow URL, checksums and verification result in
   `release-materials/05-rc-acceptance.md` and
   `release-materials/07-release-approval.md`.
6. Publish the draft only after both records are complete and approved.

Release artifacts must not be manually replaced after checksum verification.
