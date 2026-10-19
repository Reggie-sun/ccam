import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam work -p hi`, from the command line to exec: the steps the launch
// path takes before Claude Code replaces ccam (internal/app, internal/launcher).
export default function launch(): string {
  const x = 40, w = 280, cx = x + w / 2;
  const rows = [
    { y: 10, label: "ccam work -p hi", sub: "os.Args", kind: "proc" as const },
    { y: 70, label: "app.parse()", sub: "internal/app", kind: "proc" as const },
    { y: 130, label: "~/.config/ccam/config.toml", sub: "config.Load()", kind: "file" as const },
    { y: 190, label: "launcher.Prepare()", sub: "ResolveClaude() · claude-path", kind: "proc" as const },
    { y: 250, label: "profiles/work/", sub: "share.Check() · lstat", kind: "dir" as const },
    { y: 310, label: "ccam/profile-work/<pid>", sub: "session.Register()", kind: "file" as const },
    { y: 370, label: "syscall.Exec(claude)", sub: "Plan.Exec() · MarkUsed()", kind: "proc" as const },
  ];
  const h = 44, step = 1.1;
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, rows[i - 1].y + h], [cx, r.y]]));
  });
  rows.forEach((r) => parts.push(box(x, r.y, w, h, r.label, { kind: r.kind, sub: r.sub })));
  rows.forEach((r, i) => {
    parts.push(glow(x, r.y, w, h, i * step));
    if (i) parts.push(packet([cx, rows[i - 1].y + h], [cx, r.y], i * step - 0.6));
  });
  parts.push(
    reveal(rows.length * step, label(cx, 444, "CLAUDE_CONFIG_DIR=…/profiles/work", { anchor: "middle" }) +
      label(cx, 462, "CCAM_PROFILE=work", { anchor: "middle" }) +
      label(cx, 480, "claude --mcp-config=…/mcp.json -p hi", { anchor: "middle", muted: true })),
  );
  return figure("launch", 360, 490, 11, parts);
}
