import { box, figure, glow, line, packet } from "../motion.ts";

// mcp_strategy = "sync": before exec, mcp.Sync merges the group's servers
// into the account's .claude.json under Claude Code's own lock directory,
// after saving a backup.
export default function mcpSync(): string {
  const x = 40, w = 320, cx = x + w / 2, h = 44, step = 1.2;
  const rows = [
    { y: 10, label: "ccam work", sub: "mcp_strategy = \"sync\"", kind: "proc" as const },
    { y: 74, label: "groups/default/mcp.json", sub: "ModTime() vs sync-stamp", kind: "file" as const },
    { y: 138, label: ".claude.json.lock/", sub: "lock.DirLock() · 2 s", kind: "dir" as const },
    { y: 202, label: "mcp/work/backups/", sub: ".claude.json.<time>", kind: "dir" as const },
    { y: 266, label: "profiles/work/.claude.json", sub: "mcpServers · AtomicWrite()", kind: "file" as const },
    { y: 330, label: "syscall.Exec(claude)", sub: "CLAUDE_CONFIG_DIR", kind: "proc" as const },
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
  return figure("mcp-sync", 400, 384, 9, parts);
}
