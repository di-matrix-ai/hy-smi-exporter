# Release process

This project uses semantic version tags and GoReleaser.

1. Complete `release-materials/01-release-metadata.md` and remove all
   placeholders from public release documents.
2. Update `CHANGELOG.md` and move the relevant entries from `Unreleased` into a
   dated version section.
3. Run `make release-check`, `make fmt`, `make test`, `make vet`, and
   `make build` on a machine with Go 1.23 or newer. Also run
   `goreleaser release --snapshot --clean` to validate packaging without
   publishing a release.
4. Record the commit, build environment and CI run in `BUILD_INFO.md`; complete
   the dependency compliance review in `THIRD_PARTY_NOTICES.md`.
5. Complete the real-device report, compatibility matrix, metric semantics,
   collaboration/publication boundaries and all pre-artifact RC checks under
   `release-materials/`. Keep the repository private at this stage.
6. Commit the release changes and create an annotated tag, for example:

   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin main v1.0.0
   ```

7. The tag workflow creates a draft GitHub Release with Linux `amd64`, Linux
   `arm64` and checksum assets. It must not publish the Release automatically.
8. Verify the draft assets, checksums, archive contents and installation
   commands. Record the evidence in the RC acceptance and release approval
   forms.
9. After both parties approve the final draft, publish the GitHub Release,
   make the approved public repository/history available, and publish the
   external announcement.
10. Confirm or register the default port after the repository is public, then
    submit the Prometheus exporter directory PR with a DCO sign-off.

Do not publish a release containing real device serial numbers, hostnames,
secrets, or unreviewed command output in fixtures or logs.
