import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam doctor`: the checks of doctor.Run in order, each with the mark it
// prints (✓ OK, · info, ! warning, ✗ error). The middle block runs once per
// account. Any unfixed ✗ makes the exit code 1.
export default function doctorChecks(): string {
  const x = 40, w = 280, h = 40, gap = 50, mx = 352, cx = x + w / 2, step = 1;
  const rows = [
    { label: "claude --version", sub: "checkClaude()", mark: "✓" },
    { label: "$ANTHROPIC_API_KEY", sub: "checkAuthEnv()", mark: "!" },
    { label: "profiles/work/", sub: "os.Stat · 0700", mark: "✓", acct: true },
    { label: ".credentials.json", sub: "lstat · 0600 · expiresAt", mark: "✓", acct: true },
    { label: "share.Check()", sub: "lstat · readlink", mark: "✗", acct: true },
    { label: "os.ReadDir()", sub: "share.Known", mark: "·", acct: true },
    { label: "plugins/*.json", sub: "checkPlugins()", mark: "✓" },
    { label: "mcp.json", sub: "mcp.Validate()", mark: "✓" },
    { label: "/api/oauth/usage", sub: "checkUsageFormat()", mark: "✓" },
    { label: "~/.local/bin/ccwork", sub: "checkShims()", mark: "✓" },
  ];
  const top = 8;
  const y = (i: number) => top + i * gap + (i >= 2 ? 10 : 0) + (i >= 6 ? 10 : 0);
  const parts: string[] = [];
  // The per-account block.
  parts.push(box(x - 16, y(2) - 8, w + 32, y(5) + h + 8 - (y(2) - 8), "", { kind: "group" }));
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, y(i - 1) + h], [cx, y(i)]]));
  });
  rows.forEach((r, i) => {
    const d = i * step;
    parts.push(box(x, y(i), w, h, r.label, { kind: r.label.endsWith("/") ? "dir" : r.label.includes(".json") ? "file" : "proc", sub: r.sub }));
    parts.push(glow(x, y(i), w, h, d));
    if (i) parts.push(packet([cx, y(i - 1) + h], [cx, y(i)], d - 0.5));
    parts.push(reveal(d + 0.3, label(mx, y(i) + h / 2 + 4.5, r.mark, { anchor: "middle" })));
  });
  const end = y(rows.length - 1) + h, d = rows.length * step;
  parts.push(reveal(d, label(cx, end + 24, "exit 1", { anchor: "middle" })));
  return figure("doctor-checks", 380, end + 34, d + 3, parts);
}
