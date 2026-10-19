import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// MCP sign-ins after `ccam migrate`: the mcpOAuth entries of the source login
// go into each account's own .credentials.json (internal/migrate/logins.go),
// never into the group. Accounts that cannot take them now are skipped.
export default function migrateSignins(): string {
  const parts: string[] = [];
  const sx = 20, sw = 360, cx = sx + sw / 2;
  parts.push(box(sx, 8, sw, 44, "~/.claude/.credentials.json", { kind: "file", sub: "mcpOAuth" }));
  parts.push(box(sx, 72, sw, 44, "linear|sha256(type,url,headers)[:16]", { sub: "buildLogins()" }));
  parts.push(line([[cx, 52], [cx, 72]]));
  parts.push(glow(sx, 8, sw, 44, 0));
  parts.push(packet([cx, 52], [cx, 72], 0.4));
  parts.push(glow(sx, 72, sw, 44, 1));

  const ax = 52, aw = 200, h = 40, trunk = 30, top = 140, gap = 52;
  const rows = [
    { dir: "profiles/personal/", tag: "fsutil.AtomicWrite", ok: true },
    { dir: "profiles/team/", tag: "fsutil.AtomicWrite", ok: true },
    { dir: "profiles/work/", tag: "session.Active" },
    { dir: "profiles/side/", tag: "creds.ErrNone" },
    { dir: "~/.claude-old/", tag: "acct.Imported" },
  ];
  const last = top + (rows.length - 1) * gap + h / 2;
  parts.push(line([[trunk, 116], [trunk, last]], { arrow: false }));
  rows.forEach((r, i) => {
    const y = top + i * gap, cy = y + h / 2, d = 2 + i * 1.1;
    parts.push(line([[trunk, cy], [ax, cy]], { dashed: !r.ok }));
    parts.push(box(ax, y, aw, h, r.dir, { kind: "dir", sub: ".credentials.json" }));
    parts.push(glow(ax, y, aw, h, d));
    if (r.ok) parts.push(packet([trunk, cy], [ax, cy], d - 0.3));
    parts.push(reveal(d + 0.3, label(ax + aw + 10, cy + 4, (r.ok ? "✓ " : "! ") + r.tag, { muted: true })));
  });
  const d = 2 + rows.length * 1.1;
  return figure("migrate-signins", 400, top + rows.length * gap, d + 3, parts);
}
