import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// A file both sides have with different content: the candidates come from
// share.NewConflict, the choice from the prompt or --prefer, and every
// choice but KeepGroup saves the group copy before resolve() rewrites it.
export default function shareConflict(): string {
  const h = 44;
  const top = 10, cmp = 90, ask = 170, pick = 250, out = 340;
  const mid = { x: 80, w: 240 }, cx = 200;
  const picks = [
    { x: 10, label: "--prefer group", sub: "KeepGroup" },
    { x: 140, label: "--prefer mine", sub: "UseMine" },
    { x: 270, label: "--prefer merge", sub: "MergeJSON" },
  ];
  const pw = 120;
  const parts: string[] = [
    line([[100, top + h], [100, top + h + 18], [160, top + h + 18], [160, cmp]]),
    line([[300, top + h], [300, top + h + 18], [240, top + h + 18], [240, cmp]]),
    line([[cx, cmp + h], [cx, ask]]),
    ...picks.map((p) => line([[cx, ask + h], [cx, ask + h + 18], [p.x + pw / 2, ask + h + 18], [p.x + pw / 2, pick]])),
    line([[330, pick + h], [330, out]]),
    box(10, top, 180, h, "profiles/side/", { kind: "file", sub: "settings.json" }),
    box(210, top, 180, h, "groups/default/", { kind: "file", sub: "settings.json" }),
    box(mid.x, cmp, mid.w, h, "fsutil.SameContent()", { sub: "false" }),
    box(mid.x, ask, mid.w, h, "askConflict()", { sub: "Conflict.Candidates" }),
    ...picks.map((p) => box(p.x, pick, pw, h, p.label, { sub: p.sub })),
    glow(10, top, 180, h, 0),
    glow(210, top, 180, h, 0),
    packet([160, top + h + 18], [160, cmp], 0.6),
    packet([240, top + h + 18], [240, cmp], 0.6),
    glow(mid.x, cmp, mid.w, h, 1.4),
    packet([cx, cmp + h], [cx, ask], 2.0),
    glow(mid.x, ask, mid.w, h, 2.8),
    packet([330, ask + h + 18], [330, pick], 3.6),
    glow(270, pick, pw, h, 4.4),
    reveal(
      5.4,
      box(150, out, 240, h, "groups/default/settings.json", { kind: "file", sub: "MergeJSONBytes()" }) +
        box(10, out, 124, h, ".ccam-backup/", { kind: "dir", sub: "<time>/group/" }) +
        line([[150, out + 22], [134, out + 22]]),
    ),
    packet([330, pick + h], [330, out], 4.8),
    glow(10, out, 124, h, 5.8),
    glow(150, out, 240, h, 6.8),
  ];
  return figure("share-conflict", 400, 394, 9.8, parts);
}
