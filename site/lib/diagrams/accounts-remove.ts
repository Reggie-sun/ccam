import { box, figure, glow, line, packet } from "../motion.ts";

// `ccam remove work`: what goes where, in the order removeCmd does it
// (internal/cli/profiles.go, internal/profile). The login is revoked with
// `claude auth logout` before the folder goes to the trash.
export default function accountsRemove(): string {
  const x = 40, w = 320, cx = x + w / 2;
  const rows = [
    { y: 10, label: "ccam remove work", sub: "refuseRunning() · confirm()", kind: "proc" as const },
    { y: 70, label: "groups/default/plugins/", sub: "RewritePluginPaths()", kind: "dir" as const },
    { y: 130, label: "claude auth logout", sub: "logOut() · CLAUDE_CONFIG_DIR", kind: "proc" as const },
    { y: 190, label: "profiles/work/", sub: "profile.Trash()", kind: "dir" as const },
    { y: 250, label: "trash/work-20261019-193621/", sub: "fsutil.Move()", kind: "dir" as const },
    { y: 310, label: "~/.config/ccam/config.toml", sub: "cfg.Remove() · AtomicWrite()", kind: "file" as const },
    { y: 370, label: "last-used/work · usage/work.json", sub: "nameState() · os.RemoveAll()", kind: "file" as const },
    { y: 430, label: "~/.local/bin/ccwork", sub: "shim.Remove()", kind: "file" as const },
  ];
  const h = 44, step = 1.1;
  const parts: string[] = [];
  rows.forEach((r, i) => {
    if (i) parts.push(line([[cx, rows[i - 1].y + h], [cx, r.y]]));
  });
  rows.forEach((r) => parts.push(box(x, r.y, w, h, r.label, { kind: r.kind, sub: r.sub })));
  rows.forEach((r, i) => {
    parts.push(glow(x, r.y, w, h, i * step));
    if (i) parts.push(packet([cx, rows[i - 1].y + h], [cx, r.y], i * step - 0.6, i === 4 ? { label: "mv" } : {}));
  });
  return figure("accounts-remove", 400, 484, 11, parts);
}
