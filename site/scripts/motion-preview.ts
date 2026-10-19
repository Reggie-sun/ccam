// Renders motion diagrams to PNG for review, since headless Chromium on some
// machines (WSL) draws no frames. resvg plays no CSS animation, so each piece
// is drawn where it ends: dots at their destination, highlights on.
//   node scripts/motion-preview.ts [name…] [--dark] [--out dir]
import { readdirSync, readFileSync } from "node:fs";
import { mkdir, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { Resvg } from "@resvg/resvg-js";

const args = process.argv.slice(2);
const dark = args.includes("--dark");
const outAt = args.indexOf("--out");
const out = outAt >= 0 ? args[outAt + 1] : "node_modules/.cache/motion-preview";
const names = args.filter((a, i) => !a.startsWith("--") && i !== outAt + 1);
const all = readdirSync("lib/diagrams").filter((f) => f.endsWith(".ts")).map((f) => f.slice(0, -3));

// The tokens globals.css gives these pieces, light and dark.
const t = dark
  ? { card: "#1c1917", strong: "#292524", ink: "#fafaf9", muted: "#a8a29e", line: "#44403c" }
  : { card: "#ffffff", strong: "#f0efed", ink: "#0c0a09", muted: "#6b655e", line: "#d6d3d1" };
const style = `text{font-family:monospace;font-size:12px;fill:${t.ink}}.mo-sub{font-size:10.5px;fill:${t.muted}}
.mo-box rect{fill:${t.strong};stroke:${t.line}}.mo-file rect,.mo-dir rect,.mo-net rect{fill:${t.card}}.mo-dir rect{stroke-dasharray:5 3}
.mo-group rect{fill:none;stroke-dasharray:3 4}.mo-line,.mo-head{fill:none;stroke:${t.line};stroke-width:1.5}.mo-dashed{stroke-dasharray:4 4}
.mo-pkt circle{fill:${t.ink}}.mo-tag{font-size:10px}.mo-glow{fill:none;stroke:${t.ink};stroke-width:2}`;

await mkdir(out, { recursive: true });
for (const name of names.length ? names : all) {
  const html: string = (await import(join(process.cwd(), "lib/diagrams", `${name}.ts`))).default();
  let svg = /<svg[\s\S]*<\/svg>/.exec(html)![0].replace(" aria-hidden=\"true\" focusable=\"false\"", ' xmlns="http://www.w3.org/2000/svg"');
  // Dots at their destination.
  svg = svg.replace(/<g class="mo-pkt" style="([^"]*)">/g, (_, s: string) => {
    const v = Object.fromEntries(s.split(";").map((kv) => kv.replace(/^--/, "").split(":")));
    return `<g class="mo-pkt" transform="translate(${parseFloat(v.x1)} ${parseFloat(v.y1)})">`;
  });
  const [, w, h] = /viewBox="0 0 (\d+) (\d+)"/.exec(svg)!;
  svg = svg.replace(/<svg([^>]*)>/, `<svg$1><style>${style}</style><rect width="${w}" height="${h}" fill="${t.card}"/>`);
  const png = new Resvg(svg, { fitTo: { mode: "zoom", value: 2 }, font: { loadSystemFonts: true } }).render().asPng();
  const file = join(out, `${name}${dark ? "-dark" : ""}.png`);
  await writeFile(file, png);
  console.log(file);
}
