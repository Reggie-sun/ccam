import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// The flag strategy at launch (internal/launcher applyMCP, internal/mcp):
// one stat of the group's mcp.json, one argument inserted before the
// user's, a cheap check for servers added inside a session, then exec.
export default function mcpFlag(): string {
  const x = 40, w = 320, cx = x + w / 2, h = 44, step = 1.2;
  const rows = [
    { y: 10, label: "ccam work -p hi", sub: "launcher.applyMCP()", kind: "proc" as const },
    { y: 74, label: "groups/default/mcp.json", sub: "mcp.FlagArgs() · os.Stat()", kind: "file" as const },
    { y: 138, label: "Plan.InsertArgs()", sub: "--mcp-config=<path>", kind: "proc" as const },
    { y: 202, label: "profiles/work/.claude.json", sub: "mcp.NewLocal() · ModTime()", kind: "file" as const },
    { y: 266, label: "syscall.Exec(claude)", sub: "CLAUDE_CONFIG_DIR", kind: "proc" as const },
  ];
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, rows[i - 1].y + h], [cx, r.y]]));
  });
  rows.forEach((r) => parts.push(box(x, r.y, w, h, r.label, { kind: r.kind, sub: r.sub })));
  rows.forEach((r, i) => {
    parts.push(glow(x, r.y, w, h, i * step));
    if (i) parts.push(packet([cx, rows[i - 1].y + h], [cx, r.y], i * step - 0.6));
  });
  parts.push(
    reveal(rows.length * step, label(cx, 340, "claude --mcp-config=…/mcp.json -p hi", { anchor: "middle" }) +
      label(cx, 358, "~/.local/state/ccam/mcp/work/seen", { anchor: "middle", muted: true })),
  );
  return figure("mcp-flag", 400, 368, 9, parts);
}
