# Support Policy

## Community support scope

The project accepts public Issues and Pull Requests for exporter code,
documentation, parsing compatibility, metrics and deployment examples. The
latest released version is the primary supported community version.

Community maintenance does not constitute a commercial service-level
agreement or 7x24 support commitment. Production support, response time and
on-site services require a separate agreement with the responsible party.

## Responsibility boundary

- Dingmao owns the public repository, release process, code review and
  community issue triage.
- Hygon participates as technical collaborator for device, driver,
  `hy-smi` / `hcu-smi` compatibility and metric semantics.
- Hardware warranty, driver distribution, firmware, DTK and vendor tool
  licensing remain outside this repository.
- `v0.1.0` supports direct host execution and systemd only. Docker, Docker
  Compose, Kubernetes and container images are explicitly outside the first
  release compatibility scope.
- Prometheus directory inclusion is a community listing and is not an
  official Prometheus certification or support endorsement.

## Issue routing

| Issue type | Preferred channel |
| --- | --- |
| Exporter bug or feature request | Public GitHub Issue |
| Device or driver compatibility question | GitHub Issue, then Hygon technical escalation if required |
| Security vulnerability | Private channel defined in `SECURITY.md` |
| Commercial support request | Dingmao business/support contact, to be published separately |

Before publishing logs or command output, remove serial numbers, unique IDs,
hostnames, IP addresses, PCIe addresses, customer names and other sensitive
infrastructure information.
