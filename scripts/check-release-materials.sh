#!/usr/bin/env sh
set -eu

failures=0

required_files="
LICENSE
NOTICE
README.md
README.en.md
CHANGELOG.md
SECURITY.md
SUPPORT.md
ARCHITECTURE.md
ROADMAP.md
THIRD_PARTY_NOTICES.md
BUILD_INFO.md
docs/installation.md
docs/compatibility.md
docs/RELEASE_NOTES_v1.0.0.md
docs/COMMUNITY_ANNOUNCEMENT_v1.0.0.md
metric.md
release-materials/00-发布材料总目录.md
release-materials/07-release-approval.md
"

for file in $required_files; do
  if [ ! -f "$file" ]; then
    echo "ERROR missing required release file: $file" >&2
    failures=$((failures + 1))
  fi
done

public_files="
NOTICE
MAINTAINERS.md
.github/CODEOWNERS
README.md
README.en.md
SECURITY.md
THIRD_PARTY_NOTICES.md
BUILD_INFO.md
docs/compatibility.md
docs/RELEASE_NOTES_v1.0.0.md
docs/COMMUNITY_ANNOUNCEMENT_v1.0.0.md
"

placeholder_pattern='TBD|<YYYY-MM-DD>|<company-org>|<company-legal-name>|<GitHub account>|<姓名 / GitHub 账号>|<安全反馈渠道>|<待填写>|<Release asset URL>|COMPANY_GITHUB_MAINTAINER|待填写|legal entity name to be confirmed'

for file in $public_files; do
  [ -f "$file" ] || continue
  if grep -nE "$placeholder_pattern" "$file"; then
    echo "ERROR release placeholder found in $file" >&2
    failures=$((failures + 1))
  fi
done

module_path=$(awk '$1 == "module" { print $2; exit }' go.mod)
if [ -z "${EXPECTED_MODULE:-}" ]; then
  echo "ERROR EXPECTED_MODULE is required, for example github.com/company/hy-smi-exporter" >&2
  failures=$((failures + 1))
elif [ "$module_path" != "$EXPECTED_MODULE" ]; then
  echo "ERROR go.mod module is '$module_path', expected '$EXPECTED_MODULE'" >&2
  failures=$((failures + 1))
fi

if command -v go >/dev/null 2>&1; then
  go mod verify
else
  echo "WARN Go is not installed; dependency verification was skipped" >&2
fi

if [ "$failures" -ne 0 ]; then
  echo "Release material check failed with $failures issue(s)." >&2
  exit 1
fi

echo "Release material check passed."
