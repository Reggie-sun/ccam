#!/bin/sh
# Measures the launch path with hyperfine and fails when the p95 is above
# the PLAN.md targets: `ccs <name>` 10 ms, `ccs list` 30 ms, and peak memory
# above 15 MB (memory needs hyperfine 2, which records it).
set -eu

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o "$tmp/ccam" ./cmd/ccam

export CCAM_HOME="$tmp/home"
mkdir -p "$CCAM_HOME/config"
printf 'version = 1\nclaude_bin = "/bin/true"\n' > "$CCAM_HOME/config/config.toml"
"$tmp/ccam" add work --no-login >/dev/null 2>&1
"$tmp/ccam" add personal --no-login >/dev/null 2>&1
ccs="$CCAM_HOME/bin/ccs"

hyperfine -N --warmup 20 --runs 300 --export-json "$tmp/result.json" \
  "$ccs --dry-run work" \
  "$ccs work" \
  "$ccs list"

python3 - "$tmp/result.json" <<'PY'
import json, sys
limits = {"--dry-run work": 10, " work": 10, " list": 30}
failed = False
for r in json.load(open(sys.argv[1]))["results"]:
    if "times" in r:  # hyperfine 1.x
        times, mem = r["times"], []
    else:  # hyperfine 2.x
        times = [m["time_wall_clock"]["value"] for m in r["measurements"]]
        mem = [m["memory_peak_resident"]["value"] for m in r["measurements"] if "memory_peak_resident" in m]
    times = sorted(t * 1000 for t in times)
    p95 = times[int(len(times) * 0.95)]
    limit = next(v for k, v in limits.items() if r["command"].endswith(k))
    ok = p95 <= limit
    line = f"p95 {p95:6.2f} ms (limit {limit} ms)"
    if mem:
        peak = max(mem) / 2**20
        ok &= peak <= 15
        line += f", peak {peak:4.1f} MB (limit 15 MB)"
    failed |= not ok
    print(f"{'ok' if ok else 'FAIL':4} {line}  {r['command']}")
sys.exit(1 if failed else 0)
PY
