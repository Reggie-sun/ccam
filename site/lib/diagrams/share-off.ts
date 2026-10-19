import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam share side --history off`: the link to the group is removed and
// the group's current content is copied back (OpUnlink, then OpCopyBack).
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
    line([[rc, yA + h], [rc, yC + 22], [L.x + L.w, yC + 22]]),
    box(cmd.x, cmd.y, cmd.w, h, "ccam share side --history off", { sub: "share.Build() · Plan.Apply()" }),
    box(L.x, yA, L.w, h, "side/projects", { kind: "file", sub: "os.Readlink()" }),
    box(R.x, yA, R.w, h, "groups/default/", { kind: "dir", sub: "projects/" }),
    box(L.x, yB, L.w, h, "os.Remove()", { sub: "OpUnlink" }),
    box(L.x, yC, L.w, h, "profiles/side/", { kind: "dir", sub: "projects/" }),
    glow(cmd.x, cmd.y, cmd.w, h, 0),
    packet([lc, cmd.y + h], [lc, yA], 0.6),
    glow(L.x, yA, L.w, h, 1.2),
    packet([lc, yA + h], [lc, yB], 1.8),
    glow(L.x, yB, L.w, h, 2.4),
    glow(R.x, yA, R.w, h, 3.6),
    packet([rc, yA + h], [rc, yC + 22], 3.6),
    packet([rc, yC + 22], [L.x + L.w, yC + 22], 4.8),
    reveal(4.2, label(rc + 8, yB + 26, "CopyTree()", { muted: true })),
    glow(L.x, yC, L.w, h, 6.0),
  ];
  return figure("share-off", 400, 324, 9, parts);
}
