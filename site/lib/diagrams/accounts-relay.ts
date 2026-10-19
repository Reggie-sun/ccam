import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// Signing in over SSH: the browser on the user's computer lands on
// localhost:<port>/callback, which does not load there; the pasted address
// is requested on the server, where claude listens (internal/cli/login.go).
export default function accountsRelay(): string {
  const h = 44;
  const L = { x: 10, w: 180 }, R = { x: 210, w: 180 };
  const lc = L.x + L.w / 2, rc = R.x + R.w / 2;
  const r1 = 10, r2 = 70, l1 = 140, r3 = 220, r4 = 280, r5 = 340, r6 = 400;
  const parts: string[] = [
    line([[rc, r1 + h], [rc, r2]]),
    line([[R.x, r2 + 22], [lc, r2 + 22], [lc, l1]]),
    line([[L.x + L.w, l1 + 22], [200, l1 + 22], [200, r3 + 22], [R.x, r3 + 22]]),
    line([[rc, r3 + h], [rc, r4]]),
    line([[rc, r4 + h], [rc, r5]]),
    line([[rc, r5 + h], [rc, r6]]),
    // claude auth login is the process that listens on the callback port.
    line([[R.x + R.w, r2 + 22], [396, r2 + 22], [396, r5 + 22], [R.x + R.w, r5 + 22]], { arrow: false, dashed: true }),
    box(R.x, r1, R.w, h, "ccam login work", { sub: "runLogin()" }),
    box(R.x, r2, R.w, h, "claude auth login", { sub: "StdinPipe()" }),
    box(L.x, l1, L.w, h, "localhost:<port>", { kind: "net", sub: "/callback?code=…&state=…" }),
    box(R.x, r3, R.w, h, "relayLogin()", { sub: "callbackURL()" }),
    box(R.x, r4, R.w, h, "GET 127.0.0.1:<port>", { kind: "net", sub: "deliverCallback() · 10s" }),
    box(R.x, r5, R.w, h, ":<port>/callback", { kind: "net", sub: "claude" }),
    box(R.x, r6, R.w, h, ".credentials.json", { kind: "file", sub: "profiles/work/" }),
  ];
  const step = 1.1;
  parts.push(
    glow(R.x, r1, R.w, h, 0),
    packet([rc, r1 + h], [rc, r2], step - 0.6),
    glow(R.x, r2, R.w, h, step),
    packet([R.x, r2 + 22], [lc, r2 + 22], 2 * step - 0.6),
    glow(L.x, l1, L.w, h, 2 * step),
    reveal(3 * step, label(lc, l1 + h + 20, "ECONNREFUSED", { anchor: "middle" })),
    packet([200, l1 + 22], [200, r3 + 22], 4 * step - 0.6),
    label(192, r3 - 14, "stdin", { anchor: "end", muted: true }),
    glow(R.x, r3, R.w, h, 4 * step),
    packet([rc, r3 + h], [rc, r4], 5 * step - 0.6),
    glow(R.x, r4, R.w, h, 5 * step),
    packet([rc, r4 + h], [rc, r5], 6 * step - 0.6),
    glow(R.x, r5, R.w, h, 6 * step),
    packet([rc, r5 + h], [rc, r6], 7 * step - 0.6),
    glow(R.x, r6, R.w, h, 7 * step),
  );
  return figure("accounts-relay", 400, 454, 11, parts);
}
