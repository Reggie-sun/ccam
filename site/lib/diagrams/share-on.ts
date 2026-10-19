import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam share side --history on` when the group already has content: the
// account's projects/ is merged into the group, moved to the backup, and
// replaced by a link (internal/share: Build, then Plan.Apply).
export default function shareOn(): string {
  const h = 44;
  const L = { x: 10, w: 170 }, R = { x: 220, w: 170 };
  const lc = L.x + L.w / 2, rc = R.x + R.w / 2;
  const cmd = { x: 40, y: 10, w: 320 };
  const yA = 90, yB = 180, yC = 270;
  const parts: string[] = [
    line([[lc, cmd.y + h], [lc, yA]]),
    line([[rc, cmd.y + h], [rc, yA]]),
    line([[L.x + L.w, yA + 22], [R.x, yA + 22]]),
    line([[lc, yA + h], [lc, yB]]),
    box(cmd.x, cmd.y, cmd.w, h, "ccam share side --history on", { sub: "share.Build() · Plan.Apply()" }),
    box(L.x, yA, L.w, h, "profiles/side/", { kind: "dir", sub: "projects/" }),
    box(R.x, yA, R.w, h, "groups/default/", { kind: "dir", sub: "projects/" }),
    box(L.x, yB, L.w, h, ".ccam-backup/<time>/", { kind: "dir", sub: "projects/" }),
    glow(cmd.x, cmd.y, cmd.w, h, 0),
    packet([lc, cmd.y + h], [lc, yA], 0.6),
    packet([rc, cmd.y + h], [rc, yA], 0.6),
    glow(L.x, yA, L.w, h, 1.2),
    packet([L.x + L.w, yA + 22], [R.x, yA + 22], 1.8),
    glow(R.x, yA, R.w, h, 2.4),
    reveal(1.8, label(200, 156, "mergeDir()", { anchor: "middle", muted: true })),
    packet([lc, yA + h], [lc, yB], 3.0),
    glow(L.x, yB, L.w, h, 3.6),
    reveal(
      4.8,
      box(L.x, yC, L.w, h, "side/projects", { kind: "file", sub: "os.Symlink()" }) +
        line([[L.x + L.w, yC + 22], [rc, yC + 22], [rc, yA + h]], { dashed: true }),
    ),
    glow(L.x, yC, L.w, h, 5.2),
  ];
  return figure("share-on", 400, 324, 8.4, parts);
}
