#!/bin/sh
# Runs the Go tests on the Windows host of a WSL machine, without Go on
# Windows: the test binaries, ccam.exe and the fake claude are cross-compiled
# here, the repository is copied to %TEMP% (tests read testdata/ relative to
# their package) and each test binary runs there. The e2e tests take the
# prebuilt binaries through CCAM_E2E_BIN, CCAM_E2E_FAKE, CCAM_E2E_ROOT and
# CCAM_E2E_RELEASES.
#
# Usage: scripts/windows-test.sh [package-dir ...] [-- test flags]
#   scripts/windows-test.sh
#   scripts/windows-test.sh internal/share internal/e2e -- -test.run Folder -test.v
set -u
command -v cmd.exe >/dev/null 2>&1 || { echo "needs WSL with Windows interop" >&2; exit 1; }
repo=$(cd "$(dirname "$0")/.." && pwd)
# WSL interop sometimes answers with nothing; wslpath -u "" would then give
# ".", and the copy would land inside the repository.
wtemp=$(cmd.exe /c "echo %TEMP%" 2>/dev/null | tr -d '\r')
[ -n "$wtemp" ] || wtemp=$(cmd.exe /c "echo %TEMP%" 2>/dev/null | tr -d '\r')
case "$wtemp" in
[A-Za-z]:\\*) ;;
*) echo "could not read %TEMP% from Windows (got '$wtemp')" >&2; exit 1 ;;
esac
out=$(wslpath -u "$wtemp")/ccam-wt
src=$(wslpath -u "$wtemp")/ccam-src
wout="$wtemp\\ccam-wt"
case "$out/" in
"$repo"/*) echo "refusing to copy the repository into itself ($out)" >&2; exit 1 ;;
esac

pkgs=""
while [ $# -gt 0 ] && [ "$1" != "--" ]; do pkgs="$pkgs $1"; shift; done
[ "${1:-}" = "--" ] && shift
cd "$repo"
[ -z "$pkgs" ] && pkgs=$(go list -f '{{if .TestGoFiles}}{{.Dir}}{{end}}' ./... | sed "s#$repo/##")

mkdir -p "$out/releases"
rsync -a --delete --exclude site --exclude .git ./ "$src/"
GOOS=windows go build -o "$out/ccam.exe" ./cmd/ccam || exit 1
GOOS=windows go build -o "$out/fakeclaude.exe" ./testdata/fakeclaude || exit 1
for v in 0.2.0 0.3.0; do
  GOOS=windows go build -trimpath -ldflags "-s -w -X github.com/jhyunwoo/ccam/internal/cli.Version=$v" \
    -o "$out/releases/ccam-$v.exe" ./cmd/ccam || exit 1
done

fail=0
for p in $pkgs; do
  name=$(echo "$p" | tr / _)
  GOOS=windows go test -c -o "$out/$name.test.exe" "./$p" || { fail=1; continue; }
  (
    cd "$src/$p" &&
      CCAM_E2E_BIN="$wout\\ccam.exe" CCAM_E2E_FAKE="$wout\\fakeclaude.exe" \
      CCAM_E2E_ROOT="$(wslpath -w "$src")" CCAM_E2E_RELEASES="$wout\\releases" \
      WSLENV=CCAM_E2E_BIN:CCAM_E2E_FAKE:CCAM_E2E_ROOT:CCAM_E2E_RELEASES \
      "$out/$name.test.exe" "$@" 2>&1 | tr -d '\r'
  ) > "$out/$name.log"
  if grep -q '^PASS$' "$out/$name.log" && ! grep -q '^FAIL' "$out/$name.log"; then
    echo "ok   $p"
  else
    echo "FAIL $p (log: $out/$name.log)"
    fail=1
  fi
done
exit $fail
