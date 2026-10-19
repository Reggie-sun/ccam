import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// MCP flag strategy: launcher.applyMCP asks mcp.FlagArgs for one
// "--mcp-config=<group>/mcp.json" argument and Plan.InsertArgs puts it right
// after argv[0], before the user's own arguments.
export default function launchMcp(): string {
  const h = 40;
  const fileY = 10, beforeY = 90, afterY = 180;
  const slotX = 92, slotW = 180, cx = slotX + slotW / 2;
  const tokens = (y: number, withFlag: boolean) => [
    box(20, y, 64, h, "claude"),
    ...(withFlag ? [box(slotX, y, slotW, h, "--mcp-config=…/mcp.json", { kind: "file" })] : []),
    box(280, y, 44, h, "-p"),
    box(332, y, 44, h, "hi"),
  ];
  const parts: string[] = [
    line([[cx, fileY + 44], [cx, afterY]]),
    box(70, fileY, 224, 44, "groups/default/mcp.json", { kind: "file", sub: "mcp.FlagArgs() · os.Stat" }),
    label(20, beforeY - 8, "Plan.Argv", { muted: true }),
    ...tokens(beforeY, false),
    label(cx + 8, afterY - 22, "InsertArgs()", { muted: true }),
    label(20, afterY + h + 18, "syscall.Exec(Plan.Bin, Plan.Argv, Plan.Env)", { muted: true }),
  ];
  parts.push(glow(20, beforeY, 356, h, 0));
  parts.push(glow(70, fileY, 224, 44, 1.2));
  parts.push(packet([cx, fileY + 44], [cx, beforeY + h / 2], 1.8));
  parts.push(packet([cx, beforeY + h / 2], [cx, afterY], 2.8));
  parts.push(reveal(3.8, tokens(afterY, true).join("")));
  parts.push(glow(slotX, afterY, slotW, h, 4.6));
  return figure("launch-mcp", 400, 250, 8, parts);
}
