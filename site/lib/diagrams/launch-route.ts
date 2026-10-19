import { box, figure, glow, line, packet } from "../motion.ts";

// How internal/app routes a command line: a shortcut name in argv[0], an
// account name in argv[1], or a ccam subcommand, which alone reaches cobra.
export default function launchRoute(): string {
  const w = 120, h = 44, xs = [10, 140, 270];
  const ys = [10, 80, 150];
  const cols = [
    [
      { label: "ccwork -p hi", sub: "argv[0]" },
      { label: "shim.Lookup()", sub: "config.toml" },
      { label: "LaunchWith()", sub: "launcher" },
    ],
    [
      { label: "ccam work -p hi", sub: "argv[1]" },
      { label: "app.parse()", sub: "routeLaunch" },
      { label: "Launch()", sub: "launcher" },
    ],
    [
      { label: "ccam usage", sub: "argv[1]" },
      { label: "app.parse()", sub: "config.Reserved" },
      { label: "cli.Execute()", sub: "cobra" },
    ],
  ];
  const execY = 230, execX = 10, execW = 250;
  const runY = 230;
  const parts: string[] = [];
  const step = 1.1, colGap = 3.3;
  cols.forEach((col, c) => {
    const cx = xs[c] + w / 2;
    for (let r = 1; r < col.length; r++) parts.push(line([[cx, ys[r - 1] + h], [cx, ys[r]]]));
  });
  // The two launch columns meet at exec; the subcommand column runs in-process.
  parts.push(line([[xs[0] + w / 2, ys[2] + h], [xs[0] + w / 2, execY]]));
  parts.push(line([[xs[1] + w / 2, ys[2] + h], [xs[1] + w / 2, execY]]));
  parts.push(line([[xs[2] + w / 2, ys[2] + h], [xs[2] + w / 2, runY]]));
  cols.forEach((col, c) =>
    col.forEach((b, r) => parts.push(box(xs[c], ys[r], w, h, b.label, { kind: "proc", sub: b.sub }))),
  );
  parts.push(box(execX, execY, execW, h, "syscall.Exec(claude)", { kind: "proc", sub: "CLAUDE_CONFIG_DIR" }));
  parts.push(box(xs[2], runY, w, h, "RunE()", { kind: "proc", sub: "cobra.Command" }));
  cols.forEach((col, c) => {
    const t0 = c * colGap, cx = xs[c] + w / 2;
    col.forEach((_, r) => {
      parts.push(glow(xs[c], ys[r], w, h, t0 + r * step));
      if (r) parts.push(packet([cx, ys[r - 1] + h], [cx, ys[r]], t0 + r * step - 0.6));
    });
    const endY = c === 2 ? runY : execY;
    parts.push(packet([cx, ys[2] + h], [cx, endY], t0 + 3 * step - 0.6));
    parts.push(c === 2 ? glow(xs[2], runY, w, h, t0 + 3 * step) : glow(execX, execY, execW, h, t0 + 3 * step));
  });
  return figure("launch-route", 400, 284, 13, parts);
}
