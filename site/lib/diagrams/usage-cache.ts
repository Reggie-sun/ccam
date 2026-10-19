import { box, figure, glow, label, line, packet } from "../motion.ts";

// The checks Fetcher.Fetch makes for one account, in order. Each check
// either ends the lookup (right column) or passes on to the next one.
export default function usageCache(): string {
  const h = 44, sx = 10, sw = 230, scx = sx + sw / 2, ex = 290, ew = 100;
  const rows = [
    { y: 10, label: "Fetch(a, opt)", sub: "internal/usage" },
    { y: 70, label: "opt.Source", sub: "== \"cache\"", exit: "fallback()", edge: "true" },
    { y: 130, label: "<name>.json", sub: "< CacheTTL (60 s)", kind: "file" as const, exit: "return", exitSub: "source: live", edge: "true" },
    { y: 190, label: "readCredentials()", sub: ".credentials.json", kind: "file" as const, exit: "logged_out", exitSub: "errNoOAuth", edge: "err" },
    { y: 250, label: "cred.expired()", sub: "!opt.Refresh", exit: "fallback()", exitSub: "expired", edge: "true" },
    { y: 310, label: "GET /api/oauth/usage", sub: "3 s · 1 MB", kind: "net" as const, exit: "fallback()", exitSub: "error", edge: "!= 200" },
    { y: 370, label: "writeCache()", sub: "state/usage/<name>.json", kind: "file" as const },
  ];
  const step = 1.1;
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[scx, rows[i - 1].y + h], [scx, r.y]]));
    if (r.exit) parts.push(line([[sx + sw, r.y + h / 2], [ex, r.y + h / 2]]));
  });
  rows.forEach((r) => {
    parts.push(box(sx, r.y, sw, h, r.label, { kind: r.kind ?? "proc", sub: r.sub }));
    if (r.exit) {
      parts.push(box(ex, r.y, ew, h, r.exit, { kind: "proc", sub: r.exitSub }));
      parts.push(label((sx + sw + ex) / 2, r.y + h / 2 - 5, r.edge!, { anchor: "middle", muted: true }));
    }
  });
  rows.forEach((r, i) => {
    parts.push(glow(sx, r.y, sw, h, i * step));
    if (i) parts.push(packet([scx, rows[i - 1].y + h], [scx, r.y], i * step - 0.6));
  });
  return figure("usage-cache", 400, 424, (rows.length - 1) * step + 3.2, parts);
}
