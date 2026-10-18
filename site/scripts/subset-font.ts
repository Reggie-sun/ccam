// Subsets the Latin variable font to what the Latin-script locales need
// (Basic Latin, Latin-1, common punctuation), keeping it under 25 KB.
// Locales needing Latin Extended (vi, tr, pl, ha) use the system font.
import { readFile, writeFile, mkdir } from "node:fs/promises";
import subsetFont from "subset-font";

const src = "node_modules/@fontsource-variable/atkinson-hyperlegible-next/files/atkinson-hyperlegible-next-latin-wght-normal.woff2";
const out = "public/fonts/atkinson-next.woff2";

const ranges: Array<[number, number]> = [
  [0x20, 0x7e], [0xa0, 0xff], [0x2013, 0x2014], [0x2018, 0x201d], [0x2022, 0x2022], [0x2026, 0x2026], [0x20ac, 0x20ac],
];
let text = "";
for (const [from, to] of ranges) for (let c = from; c <= to; c++) text += String.fromCodePoint(c);

const font = await subsetFont(await readFile(src), text, { targetFormat: "woff2", variationAxes: { wght: { min: 400, max: 700 } } });
await mkdir("public/fonts", { recursive: true });
await writeFile(out, font);
console.log(`${out}: ${font.length} bytes`);
