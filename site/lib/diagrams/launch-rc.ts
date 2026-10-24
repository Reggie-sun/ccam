import { box, figure, glow, label, line, packet } from "../motion.ts";

// `ccam work rc`: claude remote-control takes no --mcp-config, so
// launcher.applyMCP leaves the flag out and remoteControlMCP offers to copy
// the group's servers into .claude.json while that copy is missing or stale.
export default function launchRc(): string {
  const x = 60, w = 280, cx = x + w / 2, h = 44, hAsk = 72;
  const yCmd = 10, yFlag = 74, ySynced = 138, yAsk = 202, yBranch = 312, yExec = 390;
  const split = yAsk + hAsk + 18, merge = yBranch + h + 16;
  const left = { x: 10, w: 180 }, right = { x: 210, w: 180 };
  const lcx = left.x + left.w / 2, rcx = right.x + right.w / 2;
  const parts: string[] = [
    line([[cx, yCmd + h], [cx, yFlag]]),
    line([[cx, yFlag + h], [cx, ySynced]]),
    line([[cx, ySynced + h], [cx, yAsk]]),
    line([[cx, yAsk + hAsk], [cx, split], [lcx, split], [lcx, yBranch]]),
    line([[cx, split], [rcx, split], [rcx, yBranch]]),
    line([[lcx, yBranch + h], [lcx, merge], [cx, merge], [cx, yExec]]),
    line([[rcx, yBranch + h], [rcx, merge], [cx, merge]], { arrow: false }),
    box(x, yCmd, w, h, "ccam work rc", { sub: "claudeVerb(Argv[1:])" }),
    box(x, yFlag, w, h, "mcp.FlagArgs()", { sub: "flag = nil · no --mcp-config" }),
    box(x, ySynced, w, h, "mcp.Synced()", { kind: "file", sub: "sync-stamp vs groups/default/mcp.json" }),
    // The question as ccam prints it, on three lines.
    box(x, yAsk, w, hAsk, "confirm()", {
      sub: ["Copy them into .claude.json", "(a backup is kept; servers of", "the account's own stay)? [y/N]"],
    }),
    label(lcx - 8, split - 6, "y", { anchor: "end" }),
    label(rcx + 8, split - 6, "N"),
    box(left.x, yBranch, left.w, h, "mcp.Sync()", { kind: "file", sub: ".claude.json · backups/" }),
    box(right.x, yBranch, right.w, h, "os.Stderr", { sub: "mcp_strategy = \"sync\"" }),
    box(x, yExec, w, h, "syscall.Exec(claude rc)", { sub: "CLAUDE_CONFIG_DIR" }),
  ];
  parts.push(
    glow(x, yCmd, w, h, 0),
    packet([cx, yCmd + h], [cx, yFlag], 0.5),
    glow(x, yFlag, w, h, 1.1),
    packet([cx, yFlag + h], [cx, ySynced], 1.6),
    glow(x, ySynced, w, h, 2.2),
    packet([cx, ySynced + h], [cx, yAsk], 2.7),
    glow(x, yAsk, w, hAsk, 3.3),
    packet([cx, split], [lcx, yBranch], 3.9),
    glow(left.x, yBranch, left.w, h, 4.5),
    packet([lcx, merge], [cx, yExec], 5.0),
    glow(x, yExec, w, h, 5.6),
  );
  return figure("launch-rc", 400, yExec + h + 10, 8.8, parts);
}
