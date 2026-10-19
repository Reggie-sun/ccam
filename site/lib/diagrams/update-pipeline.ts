import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam update`: find the latest tag from a redirect, download the archive
// and checksums.txt together, check SHA-256, take the ccam file out of the
// tar and rename it over the running binary (internal/update).
export default function updatePipeline(): string {
  const h = 44, x = 40, w = 320, cx = 200;
  const lx = 10, rx = 205, cw = 185, lcx = lx + cw / 2, rcx = rx + cw / 2;
  const y = { cmd: 10, latest: 70, tag: 130, dl: 208, sum: 286, tar: 346, inst: 406 };
  const parts: string[] = [
    line([[cx, y.cmd + h], [cx, y.latest]]),
    line([[cx, y.latest + h], [cx, y.tag]]),
    line([[cx, y.tag + h], [cx, 190], [lcx, 190], [lcx, y.dl]]),
    line([[cx, y.tag + h], [cx, 190], [rcx, 190], [rcx, y.dl]]),
    line([[lcx, y.dl + h], [lcx, 268], [cx, 268], [cx, y.sum]]),
    line([[rcx, y.dl + h], [rcx, 268], [cx, 268], [cx, y.sum]]),
    line([[cx, y.sum + h], [cx, y.tar]]),
    line([[cx, y.tar + h], [cx, y.inst]]),
    box(x, y.cmd, w, h, "ccam update", { sub: "cli.Version 0.4.0" }),
    box(x, y.latest, w, h, "GET releases/latest", { kind: "net", sub: "302 Location: /tag/v0.5.0" }),
    box(x, y.tag, w, h, "parseVersion(tag)", { sub: "versionRe · compareVersions()" }),
    box(lx, y.dl, cw, h, "ccam_linux_amd64.tar.gz", { kind: "net", sub: "≤ 256 MB" }),
    box(rx, y.dl, cw, h, "checksums.txt", { kind: "net", sub: "≤ 256 MB" }),
    box(x, y.sum, w, h, "sha256.Sum256(archive)", { sub: "== checksumFor(asset)" }),
    box(x, y.tar, w, h, "extract()", { sub: "tar.TypeReg · ccam" }),
    box(x, y.inst, w, h, "fsutil.AtomicWrite(exe)", { kind: "file", sub: ".ccam.tmp* → rename() · same mode" }),
  ];
  const t = 1.1;
  parts.push(
    glow(x, y.cmd, w, h, 0),
    packet([cx, y.cmd + h], [cx, y.latest], 0.5),
    glow(x, y.latest, w, h, t),
    packet([cx, y.latest + h], [cx, y.tag], t + 0.5),
    glow(x, y.tag, w, h, 2 * t),
    packet([lcx, 190], [lcx, y.dl], 2 * t + 0.5),
    packet([rcx, 190], [rcx, y.dl], 2 * t + 0.5),
    glow(lx, y.dl, cw, h, 3 * t),
    glow(rx, y.dl, cw, h, 3 * t),
    packet([cx, 268], [cx, y.sum], 3 * t + 0.5),
    glow(x, y.sum, w, h, 4 * t),
    packet([cx, y.sum + h], [cx, y.tar], 4 * t + 0.5),
    glow(x, y.tar, w, h, 5 * t),
    packet([cx, y.tar + h], [cx, y.inst], 5 * t + 0.5),
    glow(x, y.inst, w, h, 6 * t),
    reveal(6 * t, label(cx, 474, "~/.local/bin/ccwork → ccam", { anchor: "middle", muted: true })),
  );
  return figure("update-pipeline", 400, 484, 6 * t + 3.2, parts);
}
