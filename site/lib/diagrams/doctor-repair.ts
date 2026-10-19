import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam doctor --fix` on a real settings.json where the share link belongs:
// share.Build with the same settings on both sides keeps the group's copy,
// moves the account's file into .ccam-backup/<time>/ and links the path.
export default function doctorRepair(): string {
  const x = 30, w = 300, h = 44, cx = x + w / 2, mx = 356, step = 1.2;
  const rows = [
    { y: 10, label: "profiles/work/settings.json", sub: "share.Check() · NotLinked", kind: "file" as const, mark: "✗" },
    { y: 74, label: "share.Build(dir, …, s, s)", sub: "planAdopt()", kind: "proc" as const },
    { y: 138, label: "SameContent(mine, group)", sub: "false → KeepGroup", kind: "proc" as const },
    { y: 202, label: ".ccam-backup/<time>/settings.json", sub: "fsutil.Move()", kind: "file" as const },
    { y: 266, label: "settings.json → groups/default/", sub: "os.Symlink()", kind: "file" as const, mark: "✓" },
  ];
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, rows[i - 1].y + h], [cx, r.y]]));
  });
  rows.forEach((r, i) => {
    parts.push(box(x, r.y, w, h, r.label, { kind: r.kind, sub: r.sub }));
    parts.push(glow(x, r.y, w, h, i * step));
    if (i) parts.push(packet([cx, rows[i - 1].y + h], [cx, r.y], i * step - 0.6));
    if (r.mark) parts.push(reveal(i * step + 0.3, label(mx, r.y + h / 2 + 4.5, r.mark, { anchor: "middle" })));
  });
  const d = rows.length * step;
  parts.push(reveal(d, label(cx, 340, "ccam doctor → exit 0", { anchor: "middle", muted: true })));
  return figure("doctor-repair", 380, 350, d + 3, parts);
}
