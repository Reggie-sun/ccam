import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam migrate`, item by item: each share.Layout path of the source folder
// is copied into the group folder (internal/migrate). The MCP servers come
// out of ~/.claude.json; history only with --history; the login never.
export default function migrateItems(): string {
  const lx = 6, rx = 222, w = 172, h = 40, gap = 52, top = 34, step = 1.1;
  const rows = [
    { from: "~/.claude.json", sub: "--mcp", to: "mcp.json", toSub: "mcpServers" },
    { from: "settings.json", sub: "--settings", to: "settings.json" },
    { from: "CLAUDE.md", sub: "--memory", to: "CLAUDE.md" },
    { from: "agents/ skills/ …", sub: "--extensions", to: "agents/ skills/ …" },
    { from: "projects/", sub: "--history", to: "projects/", opt: true },
    { from: "history.jsonl", sub: "--history", to: "history.jsonl", opt: true },
  ];
  const parts: string[] = [
    label(lx + w / 2, 18, "~/.claude/", { anchor: "middle" }),
    label(rx + w / 2, 18, "groups/default/", { anchor: "middle" }),
  ];
  rows.forEach((r, i) => {
    const y = top + i * gap, cy = y + h / 2;
    parts.push(line([[lx + w + 4, cy], [rx - 4, cy]], { dashed: r.opt }));
    parts.push(box(lx, y, w, h, r.from, { kind: r.from.endsWith("/") || r.from.includes("/ ") ? "dir" : "file", sub: r.sub }));
    parts.push(box(rx, y, w, h, r.to, { kind: r.to.endsWith("/") || r.to.includes("/ ") ? "dir" : "file", sub: r.toSub }));
    parts.push(glow(lx, y, w, h, i * step));
    parts.push(packet([lx + w + 4, cy], [rx - 4, cy], i * step + 0.2));
    parts.push(glow(rx, y, w, h, i * step + 0.6));
  });
  const y = top + rows.length * gap, d = rows.length * step;
  parts.push(box(lx, y, w, h, ".credentials.json", { kind: "file", sub: "share.Never" }));
  parts.push(glow(lx, y, w, h, d));
  parts.push(reveal(d + 0.4, label(rx + w / 2, y + h / 2 + 4.5, "✗", { anchor: "middle" })));
  return figure("migrate-items", 400, y + h + 8, d + 3.4, parts);
}
