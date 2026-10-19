import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccwork --resume`: the shell finds a symlink in ~/.local/bin, the kernel
// runs the ccam binary it points to, and ccam picks the account from the
// name it was run as (argv[0]) with shim.Lookup (internal/app, internal/shim).
export default function shortcutsLink(): string {
  const x = 40, w = 280, cx = x + w / 2, h = 44;
  const rows = [
    { y: 10, label: "ccwork --resume", sub: "$PATH", kind: "proc" as const },
    { y: 70, label: "~/.local/bin/ccwork", sub: "-> ~/.local/bin/ccam", kind: "file" as const },
    { y: 130, label: "ccam", sub: "argv[0] = ccwork", kind: "proc" as const },
    { y: 190, label: "~/.config/ccam/config.toml", sub: "shim.Lookup() · shortcut = \"ccwork\"", kind: "file" as const },
    { y: 250, label: "launcher.LaunchWith()", sub: "work · --resume", kind: "proc" as const },
    { y: 310, label: "syscall.Exec(claude)", sub: "argv[1:]", kind: "proc" as const },
  ];
  const step = 1.1;
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, rows[i - 1].y + h], [cx, r.y]], { dashed: i === 2 }));
  });
  rows.forEach((r) => parts.push(box(x, r.y, w, h, r.label, { kind: r.kind, sub: r.sub })));
  rows.forEach((r, i) => {
    parts.push(glow(x, r.y, w, h, i * step));
    if (i) parts.push(packet([cx, rows[i - 1].y + h], [cx, r.y], i * step - 0.6));
  });
  parts.push(
    reveal(rows.length * step, label(cx, 384, "CLAUDE_CONFIG_DIR=…/profiles/work", { anchor: "middle" }) +
      label(cx, 402, "claude --mcp-config=…/mcp.json --resume", { anchor: "middle", muted: true })),
  );
  return figure("shortcuts-link", 360, 412, 10, parts);
}
