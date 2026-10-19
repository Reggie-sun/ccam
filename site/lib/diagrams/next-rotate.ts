import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// next.Rotate: start at the account used last, walk forward in slot order
// (wrapping round), skip accounts with no login, stop at the first one left.
export default function nextRotate(): string {
  const h = 44, x = 40, w = 260, cx = x + w / 2, loopX = 340;
  const rows = [
    { y: 10, label: "1 personal", sub: "last-used/personal 09:12" },
    { y: 70, label: "2 work", sub: "last-used/work 14:05" },
    { y: 130, label: "3 side", sub: "creds.Has() false" },
    { y: 190, label: "4 lab", sub: "creds.Has() true" },
  ];
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, rows[i - 1].y + h], [cx, r.y]]));
  });
  const last = rows[rows.length - 1];
  parts.push(line([[x + w, last.y + h / 2], [loopX, last.y + h / 2], [loopX, rows[0].y + h / 2], [x + w, rows[0].y + h / 2]], { dashed: true }));
  parts.push(label(loopX + 8, (rows[0].y + last.y + h) / 2 + 4, "% len", { muted: true }));
  rows.forEach((r) => parts.push(box(x, r.y, w, h, r.label, { sub: r.sub })));

  parts.push(glow(x, rows[1].y, w, h, 0));
  parts.push(packet([cx, rows[1].y + h], [cx, rows[2].y], 1.0));
  parts.push(glow(x, rows[2].y, w, h, 1.6));
  parts.push(reveal(1.6, label(x - 10, rows[2].y + h / 2 + 4, "✗", { anchor: "end" })));
  parts.push(packet([cx, rows[2].y + h], [cx, rows[3].y], 2.6));
  parts.push(glow(x, rows[3].y, w, h, 3.2));
  parts.push(reveal(3.2, label(x - 10, rows[3].y + h / 2 + 4, "✓", { anchor: "end" })));
  parts.push(line([[cx, last.y + h], [cx, 270]]));
  parts.push(packet([cx, last.y + h], [cx, 270], 4.0));
  parts.push(box(x, 270, w, h, "LaunchWith(lab)", { sub: "syscall.Exec(claude)" }));
  parts.push(glow(x, 270, w, h, 4.6));
  parts.push(reveal(4.6, label(cx, 338, "→ work → lab", { anchor: "middle" })));
  return figure("next-rotate", 400, 348, 7.8, parts);
}
