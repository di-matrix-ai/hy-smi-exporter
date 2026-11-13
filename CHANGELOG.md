# Changelog

All notable changes to this project will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and releases use [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Community contribution, security, and release documentation.
- Automated test, vet, build, and tagged release workflows.
- Draft publication of versioned Linux artifacts with manual approval before
  the GitHub Release becomes public.
- Optional `--hy-smi.collect-1s-stats` collection for detailed utilization and
  bandwidth statistics, disabled by default for predictable multi-card scrape
  latency.

### Fixed

- Populate K100_AI Unique ID and health labels through dedicated read-only
  `--showuniqueid` and `--healthcheck` enrichment queries.
- Remove unsupported alternate command-name references; `hy-smi` is the only
  vendor command used by the exporter.

## [1.0.0] - 2026-11-13

### Added

- Initial Hygon DCU / HCU Prometheus exporter.
- JSON parsing with text-output fallback for `hy-smi`.
- Prometheus metrics, systemd unit and Prometheus configuration examples.
