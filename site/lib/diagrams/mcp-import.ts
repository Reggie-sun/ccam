import { box, figure, glow, line, packet, reveal } from "../motion.ts";

// `ccam mcp import work pg`: the account's own server goes into its group's
// mcp.json (a differing group server only with --force, after a backup),
// then the account's entry is removed from .claude.json, also after a backup.
export default function mcpImport(): string {
  const h = 44;
  const L = { x: 10, w: 180 }, R = { x: 220, w: 170 };
  const lc = L.x + L.w / 2, rc = R.x + R.w / 2;
  const cmd = { x: 40, y: 10, w: 320 };
  const yA = 90, yB = 180, yC = 270;
  const parts: string[] = [
    line([[lc, cmd.y + h], [lc, yA]]),
    line([[L.x + L.w, yA + 22], [R.x, yA + 22]]),
    line([[rc, yA + h], [rc, yB]], { dashed: true }),
    line([[lc, yA + h], [lc, yB]]),
    line([[lc, yB + h], [lc, yC]]),
    box(cmd.x, cmd.y, cmd.w, h, "ccam mcp import work pg", { sub: "mcp.LocalServers()" }),
    box(L.x, yA, L.w, h, "work/.claude.json", { kind: "file", sub: "mcpServers.pg" }),
    box(R.x, yA, R.w, h, "default/mcp.json", { kind: "file", sub: "mcp.SameDef()" }),
    box(R.x, yB, R.w, h, ".ccam-backup/<time>/", { kind: "dir", sub: "mcp.json · --force" }),
    box(L.x, yB, L.w, h, "mcp/work/backups/", { kind: "dir", sub: ".claude.json.<time>" }),
    box(L.x, yC, L.w, h, "mcp.RemoveLocal()", { sub: "lock.DirLock()" }),
    glow(cmd.x, cmd.y, cmd.w, h, 0),
    packet([lc, cmd.y + h], [lc, yA], 0.6),
    glow(L.x, yA, L.w, h, 1.2),
    packet([L.x + L.w, yA + 22], [R.x, yA + 22], 1.8),
    glow(R.x, yA, R.w, h, 2.4),
    packet([rc, yA + h], [rc, yB], 3.0),
    glow(R.x, yB, R.w, h, 3.6),
    packet([lc, yA + h], [lc, yB], 5.0),
    glow(L.x, yB, L.w, h, 5.6),
    packet([lc, yB + h], [lc, yC], 6.2),
    glow(L.x, yC, L.w, h, 6.8),
    reveal(
      4.0,
      line([[rc, yB + h], [rc, yC]]) + box(R.x, yC, R.w, h, "default/mcp.json", { kind: "file", sub: "mcpServers.pg" }),
    ),
    glow(R.x, yC, R.w, h, 4.0),
  ];
  return figure("mcp-import", 400, 324, 9.8, parts);
}
