// Subsets the two Latin fonts to what the Latin-script locales need
// (Basic Latin, Latin-1, common punctuation). Locales needing Latin Extended
// (vi, tr, pl, ha) and every other script use the system font.
import { readFile, writeFile, mkdir } from "node:fs/promises";
import subsetFont from "subset-font";

const ranges: Array<[number, number]> = [
  [0x20, 0x7e], [0xa0, 0xff], [0x2013, 0x2014], [0x2018, 0x201d], [0x2022, 0x2022], [0x2026, 0x2026], [0x20ac, 0x20ac], [0x2192, 0x2192],
];
let text = "";
for (const [from, to] of ranges) for (let c = from; c <= to; c++) text += String.fromCodePoint(c);

const fonts = [
  // Body and UI: Inter 400–600 (DESIGN.md body, title and button weights).
  { src: "inter/files/inter-latin-wght-normal.woff2", out: "inter.woff2", axes: { wght: { min: 400, max: 600 } } },
  // Display: Newsreader pinned to weight 300 at a display optical size, standing in for Waldenburg Light.
  { src: "newsreader/files/newsreader-latin-standard-normal.woff2", out: "newsreader.woff2", axes: { wght: 300, opsz: 60 } },
];

await mkdir("public/fonts", { recursive: true });
for (const f of fonts) {
  const font = await subsetFont(await readFile(`node_modules/@fontsource-variable/${f.src}`), text, { targetFormat: "woff2", variationAxes: f.axes });
  await writeFile(`public/fonts/${f.out}`, font);
  console.log(`public/fonts/${f.out}: ${font.length} bytes`);
}
