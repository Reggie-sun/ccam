import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam group leave team side`: one OpCopyBack step (copyInto). The group's
// current content is copied into side/.ccam-backup/incomplete-projects
// first; only then is the link removed (unlink, which checks it still
// points into the group) and the copy renamed into its place.
export default function shareOff(): string {
  const h = 44;
  const L = { x: 10, w: 170 }, R = { x: 220, w: 170 };
  const lc = L.x + L.w / 2, rc = R.x + R.w / 2;
  const cmd = { x: 40, y: 10, w: 320 };
  const yA = 90, yB = 180, yC = 270;
  const parts: string[] = [
    line([[lc, cmd.y + h], [lc, yA]]),
    line([[L.x + L.w, yA + 22], [R.x, yA + 22]], { dashed: true }),
    line([[lc, yA + h], [lc, yB]]),
    line([[lc, yB + h], [lc, yC]]),
    line([[rc, yA + h], [rc, yB]]),
    line([[rc, yB + h], [rc, yC + 22], [L.x + L.w, yC + 22]]),
    box(cmd.x, cmd.y, cmd.w, h, "ccam group leave team side", { sub: "group.Build() · Change.Apply()" }),
    box(L.x, yA, L.w, h, "side/projects", { kind: "file", sub: "os.Readlink()" }),
    box(R.x, yA, R.w, h, "groups/team/", { kind: "dir", sub: "projects/" }),
    box(R.x, yB, R.w, h, "side/.ccam-backup/", { kind: "dir", sub: "incomplete-projects/" }),
    box(L.x, yB, L.w, h, "os.Remove()", { sub: "unlink()" }),
    box(L.x, yC, L.w, h, "profiles/side/", { kind: "dir", sub: "projects/" }),
    glow(cmd.x, cmd.y, cmd.w, h, 0),
    packet([lc, cmd.y + h], [lc, yA], 0.6),
    glow(L.x, yA, L.w, h, 1.2),
    glow(R.x, yA, R.w, h, 1.8),
    packet([rc, yA + h], [rc, yB], 1.8),
    reveal(1.8, label(rc + 8, yA + h + 27, "CopyTree()", { muted: true })),
    glow(R.x, yB, R.w, h, 2.4),
    packet([lc, yA + h], [lc, yB], 3.0),
    glow(L.x, yB, L.w, h, 3.6),
    packet([rc, yB + h], [rc, yC + 22], 4.2),
    reveal(4.2, label(rc + 8, yB + h + 27, "os.Rename()", { muted: true })),
    packet([rc, yC + 22], [L.x + L.w, yC + 22], 5.4),
    glow(L.x, yC, L.w, h, 6.0),
  ];
  return figure("share-off", 400, 324, 9, parts);
}
