import { box, figure, glow, label, line, packet } from "../motion.ts";

// Before downloading, `ccam update` checks who owns the binary: a system
// folder or dpkg means a package manager (managedError), an unwritable
// folder means the install script; anything else goes on to Fetch.
export default function updateManaged(): string {
  const h = 44, sx = 10, sw = 230, scx = sx + sw / 2, ex = 280, ew = 110;
  const y = { exe: 10, dir: 70, dpkg: 130, wr: 190, fetch: 250 };
  const parts: string[] = [
    line([[scx, y.exe + h], [scx, y.dir]]),
    line([[scx, y.dir + h], [scx, y.dpkg]]),
    line([[scx, y.dpkg + h], [scx, y.wr]]),
    line([[scx, y.wr + h], [scx, y.fetch]]),
    line([[sx + sw, y.dir + h / 2], [ex, y.dir + h / 2]]),
    line([[sx + sw, y.dpkg + h / 2], [ex, y.dpkg + h / 2]]),
    line([[sx + sw, y.wr + h / 2], [ex, y.wr + h / 2]]),
    label((sx + sw + ex) / 2, y.dir + h / 2 - 5, "true", { anchor: "middle", muted: true }),
    label((sx + sw + ex) / 2, y.dpkg + h / 2 - 5, "exit 0", { anchor: "middle", muted: true }),
    label((sx + sw + ex) / 2, y.wr + h / 2 - 5, "err", { anchor: "middle", muted: true }),
    box(sx, y.exe, sw, h, "os.Executable()", { sub: "filepath.EvalSymlinks()" }),
    box(sx, y.dir, sw, h, "filepath.Dir(exe)", { sub: "/usr/bin /bin /usr/sbin /sbin" }),
    box(sx, y.dpkg, sw, h, "dpkg -S <exe>", { sub: "PackageOwned()" }),
    box(sx, y.wr, sw, h, "writable(dir)", { kind: "dir", sub: "os.CreateTemp()" }),
    box(sx, y.fetch, sw, h, "Fetch(tag)", { kind: "net", sub: "Install()" }),
    box(ex, y.dir, ew, y.dpkg + h - y.dir, "managedError()", { sub: "apt install *.deb" }),
    box(ex, y.wr, ew, h, "install.sh", { sub: "~/.local/bin" }),
  ];
  const t = 1.1;
  parts.push(
    glow(sx, y.exe, sw, h, 0),
    packet([scx, y.exe + h], [scx, y.dir], 0.5, { label: "/usr/bin/ccam" }),
    glow(sx, y.dir, sw, h, t),
    packet([sx + sw, y.dir + h / 2], [ex, y.dir + h / 2], t + 0.5),
    glow(ex, y.dir, ew, y.dpkg + h - y.dir, 2 * t),
    packet([scx, y.exe + h], [scx, y.dir], 3 * t, { label: "~/.local/bin/ccam" }),
    glow(sx, y.dir, sw, h, 3 * t + 0.5),
    packet([scx, y.dir + h], [scx, y.dpkg], 4 * t),
    glow(sx, y.dpkg, sw, h, 4 * t + 0.5),
    packet([scx, y.dpkg + h], [scx, y.wr], 5 * t),
    glow(sx, y.wr, sw, h, 5 * t + 0.5),
    packet([scx, y.wr + h], [scx, y.fetch], 6 * t),
    glow(sx, y.fetch, sw, h, 6 * t + 0.5),
  );
  return figure("update-managed", 400, 304, 6 * t + 3.4, parts);
}
