#!/bin/sh
# Compares benchmarks of the working tree against a base ref with benchstat
# and fails when any benchmark is more than 20% slower (p < 0.05).
# Usage: scripts/bench-compare.sh origin/main
set -eu
base=${1:?usage: bench-compare.sh <base-ref>}
pkgs="./internal/launcher/ ./internal/profile/ ./internal/config/"

tmp=$(mktemp -d)
trap 'git worktree remove --force "$tmp/base" 2>/dev/null; rm -rf "$tmp"' EXIT
git worktree add --detach "$tmp/base" "$base" >/dev/null

# The base may lack a package or benchmark added since, so its failures are
# tolerated, but a base with no benchmarks at all would compare nothing.
(cd "$tmp/base" && go test -run '^$' -bench . -count 8 $pkgs 2>/dev/null) > "$tmp/old.txt" || true
grep -q '^Benchmark' "$tmp/old.txt" || { echo "no benchmarks ran on $base" >&2; exit 1; }
go test -run '^$' -bench . -count 8 $pkgs > "$tmp/new.txt"
benchstat "$tmp/old.txt" "$tmp/new.txt" | tee "$tmp/diff.txt"

python3 - "$tmp/diff.txt" <<'PY'
import re, sys
# benchstat prints one table per package and unit; only time tables count
# (in B/s tables a higher number is better).
blocks = open(sys.argv[1]).read().split("\n\n")
bad = [l for b in blocks if "sec/op" in b for l in b.splitlines()
       if (m := re.search(r"([+-]\d+\.\d+)% \(p=([\d.]+)", l)) and float(m[1]) > 20 and float(m[2]) < 0.05]
if bad:
    print("\nslower by more than 20%:\n" + "\n".join(bad))
    sys.exit(1)
print("\nno benchmark regressed by more than 20%")
PY
