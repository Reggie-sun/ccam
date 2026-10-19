import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// How buildEntry (internal/migrate) sorts one source path into a plan state,
// and the choices a Conflict offers (share.NewConflict).
export default function migrateStates(): string {
  const cx0 = 10, cw = 210, sx = 262, sw = 128, h = 40, gap = 58, step = 1.1;
  const mid = cx0 + cw / 2;
  const rows = [
    { check: "EvalSymlinks(src)", sub: "err", state: "Missing" },
    { check: "EvalSymlinks(dst) == src", state: "Linked" },
    { check: "!Exists(dst)", state: "New" },
    { check: "entry.Dir", sub: "scan()", state: "Merge", stateSub: "Files == 0 → Same" },
    { check: "SameContent(src, dst)", state: "Same" },
  ];
  const parts: string[] = [];
  rows.forEach((r, i) => {
    const y = 8 + i * gap, cy = y + h / 2, d = i * step;
    if (i) {
      parts.push(line([[mid, y - gap + h], [mid, y]]));
      parts.push(label(mid + 6, y - 5, "false", { muted: true }));
      parts.push(packet([mid, y - gap + h], [mid, y], d - 0.5));
    }
    parts.push(line([[cx0 + cw, cy], [sx, cy]]));
    parts.push(label((cx0 + cw + sx) / 2, cy - 5, "true", { anchor: "middle", muted: true }));
    parts.push(box(cx0, y, cw, h, r.check, { sub: r.sub }));
    parts.push(box(sx, y, sw, h, r.state, { kind: "file", sub: r.stateSub }));
    parts.push(glow(cx0, y, cw, h, d));
    parts.push(glow(sx, y, sw, h, d + 0.4));
  });
  // The false branch of SameContent: a conflict and its choices.
  const y = 8 + rows.length * gap, d = rows.length * step;
  parts.push(line([[mid, y - gap + h], [mid, y]]));
  parts.push(label(mid + 6, y - 5, "false", { muted: true }));
  parts.push(packet([mid, y - gap + h], [mid, y], d - 0.5));
  parts.push(box(cx0, y, cw, h, "Conflict", { kind: "file", sub: "share.NewConflict()" }));
  parts.push(glow(cx0, y, cw, h, d));
  const cy = y + h + 30, bw = 92;
  const choices = ["KeepGroup", "UseMine", "MergeJSON", "MergeLines"];
  parts.push(line([[mid, y + h], [mid, cy - 12]], { arrow: false }));
  parts.push(line([[4 + bw / 2, cy - 12], [4 + 3 * (bw + 8) + bw / 2, cy - 12]], { arrow: false }));
  choices.forEach((c, i) => {
    const x = 4 + i * (bw + 8);
    parts.push(line([[x + bw / 2, cy - 12], [x + bw / 2, cy]]));
    parts.push(reveal(d + 0.8 + i * 0.3, box(x, cy, bw, h, c, { kind: "proc" })));
  });
  return figure("migrate-states", 400, cy + h + 8, d + 4, parts);
}
