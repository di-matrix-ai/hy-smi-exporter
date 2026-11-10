# Release process

This project uses semantic version tags and GoReleaser.

1. Complete `release-materials/01-release-metadata.md` and remove all
   placeholders from public release documents.
2. Update `CHANGELOG.md` and move the relevant entries from `Unreleased` into a
   dated version section.
3. Run `make release-check`, `make fmt`, `make test`, `make vet`, and
   `make build` on a machine with
   Go 1.23 or newer.
4. Record the commit, build environment and CI run in `BUILD_INFO.md`; complete
   the dependency compliance review in `THIRD_PARTY_NOTICES.md`.
5. Complete the real-device report, compatibility matrix, metric semantics
   confirmation, RC acceptance record and publication authorization under
   `release-materials/`.
6. Commit the release changes and create an annotated tag, for example:

   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin main v0.1.0
   ```

7. The tag workflow publishes GitHub Release archives for Linux `amd64` and
   `arm64` and checksums. The first release does not publish a container image.
8. Verify the GitHub Release, checksums and the installation commands in
   both README files before announcing the release.
9. Both parties complete `release-materials/07-release-approval.md`. Only then
   publish the external announcement and submit the Prometheus directory PR.

Do not publish a release containing real device serial numbers, hostnames,
secrets, or unreviewed command output in fixtures or logs.
