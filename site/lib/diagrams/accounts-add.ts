import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam add work --shortcut ccwork`: the steps addProfile takes, from the
// name check to the shortcut link (internal/cli/profiles.go).
export default function accountsAdd(): string {
  const x = 40, w = 320, cx = x + w / 2;
  const rows = [
    { y: 10, label: "ccam add work", sub: "config.ValidateName()", kind: "proc" as const },
    { y: 70, label: "profiles/work/", sub: "profile.Add() · 0700 · slot 1", kind: "dir" as const },
    { y: 130, label: "ln -s groups/default/settings.json", sub: "share.Build() · Apply()", kind: "file" as const },
    { y: 190, label: "~/.config/ccam/config.toml", sub: "fsutil.AtomicWrite()", kind: "file" as const },
    { y: 250, label: "claude auth login", sub: "runLogin() · CLAUDE_CONFIG_DIR", kind: "proc" as const },
    { y: 310, label: "profiles/work/.claude.json", sub: "TopLevelValue(oauthAccount)", kind: "file" as const },
    { y: 370, label: "ln -s ccam ~/.local/bin/ccwork", sub: "shim.Install()", kind: "file" as const },
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
    reveal(rows.length * step, label(cx, 442, "email = 'work@example.com'", { anchor: "middle" }) +
      label(cx, 460, "slot = 1 · shortcut = 'ccwork'", { anchor: "middle", muted: true })),
  );
  return figure("accounts-add", 400, 470, 11, parts);
}
