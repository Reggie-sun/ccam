import { box, figure, glow, line, packet, reveal } from "../motion.ts";

// shim.Check, in order, for `ccam shortcut work ccwork`: each check on the
// left, and beside it a name that this check refuses. Only a name that
// passes them all is saved and linked.
export default function shortcutsCheck(): string {
  const x = 10, w = 250, cx = x + w / 2, h = 44;
  const rx = 290, rw = 100;
  const rows = [
    { y: 10, label: "ccam shortcut work ccwork", sub: "checkShortcut()" },
    { y: 66, label: "shortcutRe", sub: "^[a-z0-9][a-z0-9._-]{0,63}$", bad: "Work" },
    { y: 122, label: "config.Reserved", sub: "ccam · claude · run · usage …", bad: "run" },
    { y: 178, label: "cfg.FindShortcut()", sub: "profiles[].shortcut", bad: "ccside" },
    { y: 234, label: "os.Lstat()", sub: "~/.local/bin/<name>", bad: "ccfile" },
    { y: 290, label: "$PATH", sub: "os.Stat() · mode & 0o111", bad: "git" },
    { y: 346, label: "config.toml · os.Symlink()", sub: "shim.Install()" },
  ];
  const step = 1.1;
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, rows[i - 1].y + h], [cx, r.y]]));
    if (r.bad) parts.push(line([[x + w, r.y + h / 2], [rx, r.y + h / 2]], { dashed: true }));
  });
  rows.forEach((r, i) => parts.push(box(x, r.y, w, h, r.label, { kind: i === rows.length - 1 ? "file" : "proc", sub: r.sub })));
  rows.forEach((r, i) => {
    parts.push(glow(x, r.y, w, h, i * step));
    if (i) parts.push(packet([cx, rows[i - 1].y + h], [cx, r.y], i * step - 0.6));
    if (r.bad) parts.push(reveal(i * step, box(rx, r.y, rw, h, r.bad, { kind: "net", sub: "return err" })));
  });
  return figure("shortcuts-check", 400, 400, 10, parts);
}
