import { box, figure, glow, line, packet } from "../motion.ts";

// Claude Code pipes its status line JSON to `ccam statusline`, which saves
// the rate limits for `ccam usage` and prints the line Claude Code shows.
export default function usageStatusline(): string {
  const h = 44, w = 280, x = 60, cx = 200;
  const lx = 10, rx = 205, cw = 185, lcx = lx + cw / 2, rcx = rx + cw / 2;
  const y = { claude: 10, sl: 80, parse: 150, mid: 230, end: 310 };
  const parts: string[] = [
    line([[cx, y.claude + h], [cx, y.sl]]),
    line([[cx, y.sl + h], [cx, y.parse]]),
    line([[cx, y.parse + h], [cx, 212], [lcx, 212], [lcx, y.mid]]),
    line([[cx, y.parse + h], [cx, 212], [rcx, 212], [rcx, y.mid]]),
    line([[lcx, y.mid + h], [lcx, y.end]], { dashed: true }),
    line([[rcx, y.mid + h], [rcx, y.end]]),
    box(x, y.claude, w, h, "claude", { sub: "settings.json statusLine" }),
    box(x, y.sl, w, h, "ccam statusline", { sub: "CCAM_PROFILE=work" }),
    box(x, y.parse, w, h, "ParseStatusLine()", { sub: "model · rate_limits" }),
    box(lx, y.mid, cw, h, "work.statusline.json", { kind: "file", sub: "SaveStatusLine()" }),
    box(rx, y.mid, cw, h, "--wrap <cmd>", { kind: "group", sub: "sh -c · 2 s" }),
    box(lx, y.end, cw, h, "ccam usage", { sub: "fallback()" }),
    box(rx, y.end, cw, h, "[work] · Opus", { kind: "file", sub: "5h 24% · 7d 41%" }),
  ];
  parts.push(
    glow(x, y.claude, w, h, 0),
    packet([cx, y.claude + h], [cx, y.sl], 0.5, { label: "stdin" }),
    glow(x, y.sl, w, h, 1.1),
    packet([cx, y.sl + h], [cx, y.parse], 1.6),
    glow(x, y.parse, w, h, 2.2),
    packet([lcx, 212], [lcx, y.mid], 2.8),
    packet([rcx, 212], [rcx, y.mid], 2.8),
    glow(lx, y.mid, cw, h, 3.3),
    glow(rx, y.mid, cw, h, 3.3),
    packet([rcx, y.mid + h], [rcx, y.end], 3.9, { label: "stdout" }),
    glow(rx, y.end, cw, h, 4.4),
    packet([lcx, y.mid + h], [lcx, y.end], 5.0),
    glow(lx, y.end, cw, h, 5.5),
  );
  return figure("usage-statusline", 400, 364, 8.5, parts);
}
