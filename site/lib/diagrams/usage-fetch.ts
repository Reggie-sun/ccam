import { box, figure, glow, line, packet } from "../motion.ts";

// `ccam usage`: one goroutine per account over one HTTP client; an account
// that does not answer in 3 s falls back to its newest saved numbers
// (internal/usage Fetcher.FetchAll, Fetch, fallback).
export default function usageFetch(): string {
  const h = 44, colW = 120, cols = [10, 140, 270], cx = cols.map((x) => x + colW / 2);
  const top = { x: 60, y: 10, w: 280 }, tcx = 200;
  const yAcct = 84, yNet = 158, yRes = 232, yFile = 306, yTable = 380;
  const names = ["personal", "work", "side"];
  const parts: string[] = [];

  // Lines first, so boxes sit on top of them.
  cx.forEach((x) => parts.push(line([[tcx, top.y + h], [tcx, 66], [x, 66], [x, yAcct]])));
  cx.forEach((x) => parts.push(line([[x, yAcct + h], [x, yNet]])));
  cx.forEach((x) => parts.push(line([[x, yNet + h], [x, yRes]])));
  parts.push(line([[cx[2], yRes + h], [cx[2], yFile]]));
  parts.push(line([[cx[0], yRes + h], [cx[0], yTable]]));
  parts.push(line([[cx[1], yRes + h], [cx[1], yTable]]));
  parts.push(line([[cx[2], yFile + h], [cx[2], yTable]]));

  parts.push(box(top.x, top.y, top.w, h, "ccam usage", { sub: "Fetcher.FetchAll()" }));
  names.forEach((n, i) => parts.push(box(cols[i], yAcct, colW, h, n, { sub: "wg.Go()" })));
  parts.push(box(10, yNet, 380, h, "api.anthropic.com", { kind: "net", sub: "GET /api/oauth/usage · 2 s dial · 3 s" }));
  parts.push(box(cols[0], yRes, colW, h, "200", { sub: "source: live" }));
  parts.push(box(cols[1], yRes, colW, h, "200", { sub: "source: live" }));
  parts.push(box(cols[2], yRes, colW, h, "timeout", { sub: "> 3 s" }));
  parts.push(box(cols[2], yFile, colW, h, "side*.json", { kind: "file", sub: "fallback()" }));
  parts.push(box(10, yTable, 380, h, "usage.Table()", { sub: "usage.Best()" }));

  // Timeline.
  parts.push(glow(top.x, top.y, top.w, h, 0));
  cx.forEach((x, i) => {
    parts.push(packet([x, 66], [x, yAcct], 0.6));
    parts.push(glow(cols[i], yAcct, colW, h, 1.1));
    parts.push(packet([x, yAcct + h], [x, yNet], 1.7));
  });
  parts.push(glow(10, yNet, 380, h, 2.3));
  [0, 1].forEach((i) => {
    parts.push(packet([cx[i], yNet + h], [cx[i], yRes], 2.9));
    parts.push(glow(cols[i], yRes, colW, h, 3.4));
  });
  parts.push(glow(cols[2], yRes, colW, h, 4.4));
  parts.push(packet([cx[2], yRes + h], [cx[2], yFile], 4.9));
  parts.push(glow(cols[2], yFile, colW, h, 5.4));
  [0, 1].forEach((i) => parts.push(packet([cx[i], yRes + h], [cx[i], yTable], 6.0)));
  parts.push(packet([cx[2], yFile + h], [cx[2], yTable], 6.0));
  parts.push(glow(10, yTable, 380, h, 6.6));
  return figure("usage-fetch", 400, 434, 9.6, parts);
}
