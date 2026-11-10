# Security Policy

## Supported versions

Security fixes are made for the latest released version. Users should upgrade
to the newest release before reporting a vulnerability.

## Reporting a vulnerability

Please do not open a public issue for an unpatched vulnerability. Use GitHub's
private security advisory workflow for this repository, or contact Dingmao at
`esteric_shen@di-matrix.com` while the private reporting channel is being
configured.

Include the affected version or commit, deployment mode, reproduction steps,
and the impact you observed. Redact credentials, hostnames, serial numbers,
and other sensitive infrastructure details.

We will acknowledge a report as soon as practical, investigate it, and
coordinate a fix and disclosure timeline with the reporter. Please allow time
for a fix before making details public.

## Scope notes

The exporter invokes only read-only `hy-smi` commands. A report involving a
command that can modify device state is especially important. The exporter
should never receive secrets through metric labels or command output; please
report any such exposure immediately.
