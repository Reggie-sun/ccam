// OG images per page and locale, rendered with resvg because satori cannot
// shape Arabic, Urdu or Indic text. Noto fonts are cached in node_modules/.cache.
import { existsSync, readFileSync } from "node:fs";
import { mkdir, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { Resvg } from "@resvg/resvg-js";
import { getLocale, type Font } from "../i18n/locales.ts";
import { builtLocales, messages, OUT, routes, sourcePages } from "./lib.ts";

const FONT_DIR = "node_modules/.cache/og-fonts";
const RAW = "https://github.com/google/fonts/raw/main/ofl/";
const fontFiles: Record<string, string> = {
  display: "newsreader/Newsreader%5Bopsz,wght%5D.ttf",
  latin: "notosans/NotoSans%5Bwdth,wght%5D.ttf",
  ko: "notosanskr/NotoSansKR%5Bwght%5D.ttf",
  ja: "notosansjp/NotoSansJP%5Bwght%5D.ttf",
  "zh-hans": "notosanssc/NotoSansSC%5Bwght%5D.ttf",
  "zh-hant": "notosanstc/NotoSansTC%5Bwght%5D.ttf",
  arabic: "notosansarabic/NotoSansArabic%5Bwdth,wght%5D.ttf",
  urdu: "notonastaliqurdu/NotoNastaliqUrdu%5Bwght%5D.ttf",
  devanagari: "notosansdevanagari/NotoSansDevanagari%5Bwdth,wght%5D.ttf",
  bengali: "notosansbengali/NotoSansBengali%5Bwdth,wght%5D.ttf",
  tamil: "notosanstamil/NotoSansTamil%5Bwdth,wght%5D.ttf",
  telugu: "notosanstelugu/NotoSansTelugu%5Bwdth,wght%5D.ttf",
  gujarati: "notosansgujarati/NotoSansGujarati%5Bwdth,wght%5D.ttf",
  kannada: "notosanskannada/NotoSansKannada%5Bwdth,wght%5D.ttf",
  thai: "notosansthai/NotoSansThai%5Bwdth,wght%5D.ttf",
};
const indicScript: Record<string, string> = { hi: "devanagari", mr: "devanagari", bn: "bengali", ta: "tamil", te: "telugu", gu: "gujarati", kn: "kannada" };

async function font(key: string): Promise<string> {
  const file = join(FONT_DIR, `${key}.ttf`);
  if (!existsSync(file)) {
    await mkdir(FONT_DIR, { recursive: true });
    const res = await fetch(RAW + fontFiles[key]);
    if (!res.ok) throw new Error(`font ${key}: HTTP ${res.status}`);
    await writeFile(file, Buffer.from(await res.arrayBuffer()));
  }
  return file;
}

function scriptFont(code: string, f: Font): string {
  if (f === "brand") return "display";
  if (f === "latin-ext" || f === "cyrillic") return "latin";
  if (f === "indic") return indicScript[code];
  return f;
}

const esc = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

/** Rough line breaking: SVG text does not wrap. Spaceless scripts break per character. */
function wrap(text: string, maxEm: number, wide: boolean): string[] {
  const unit = (ch: string) => (wide ? 1 : /[\u0000-ɏ]/.test(ch) ? 0.56 : 0.62);
  const width = (s: string) => [...s].reduce((w, c) => w + unit(c), 0);
  const tokens = wide ? [...text] : text.split(/(\s+)/);
  const lines: string[] = [];
  let line = "";
  for (const t of tokens) {
    if (line && width(line + t) > maxEm) {
      lines.push(line.trim());
      line = t.trimStart();
    } else line += t;
  }
  if (line.trim()) lines.push(line.trim());
  return lines.slice(0, 4);
}

function svg(opts: { title: string; kicker: string; family: string; rtl: boolean; wide: boolean }): string {
  const size = opts.wide ? 60 : 66;
  const lines = wrap(opts.title, 1000 / size, opts.wide); // 1000 px of text width
  const x = opts.rtl ? 1100 : 100;
  const anchor = opts.rtl ? "end" : "start";
  const tspans = lines.map((l, i) => `<tspan x="${x}" dy="${i ? size * 1.25 : 0}">${esc(l)}</tspan>`).join("");
  const top = 300 - ((lines.length - 1) * size * 1.25) / 2;
  return `<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="630" viewBox="0 0 1200 630">
  <rect width="1200" height="630" fill="#f5f5f5"/>
  <text x="${x}" y="120" text-anchor="${anchor}" font-family="Newsreader" font-size="56" fill="#0c0a09">ccam</text>
  <text x="${x}" y="${top}" text-anchor="${anchor}" direction="${opts.rtl ? "rtl" : "ltr"}" font-family="${opts.family}" font-size="${size}" fill="#0c0a09">${tspans}</text>
  <text x="${x}" y="545" text-anchor="${anchor}" direction="${opts.rtl ? "rtl" : "ltr"}" font-family="${opts.family}" font-size="30" fill="#4e4e4e">${esc(opts.kicker)}</text>
</svg>`;
}

const familyName: Record<string, string> = {
  display: "Newsreader", latin: "Noto Sans", ko: "Noto Sans KR", ja: "Noto Sans JP", "zh-hans": "Noto Sans SC",
  "zh-hant": "Noto Sans TC", arabic: "Noto Sans Arabic", urdu: "Noto Nastaliq Urdu", devanagari: "Noto Sans Devanagari",
  bengali: "Noto Sans Bengali", tamil: "Noto Sans Tamil", telugu: "Noto Sans Telugu", gujarati: "Noto Sans Gujarati",
  kannada: "Noto Sans Kannada", thai: "Noto Sans Thai",
};

const pages = sourcePages();
let count = 0;
for (const l of builtLocales()) {
  const m = messages(l.code);
  const key = scriptFont(l.code, l.font);
  const fonts = [await font("display"), await font(key)];
  const wide = ["ko", "ja", "zh-hans", "zh-hant"].includes(key);
  const titles: Record<string, string> = { home: m.landing.headline, changelog: m.changelog.title, languages: m.languages.title };
  for (const p of pages.filter((p) => p.locale === l.code)) titles[`${p.section}-${p.slug}`] = p.title;
  for (const r of routes().filter((r) => r.locale === l.code)) {
    const image = new Resvg(svg({ title: titles[r.key], kicker: "ccam.dev", family: familyName[key], rtl: getLocale(l.code)!.dir === "rtl", wide }), {
      font: { fontFiles: fonts, loadSystemFonts: false, defaultFontFamily: familyName[key] },
    }).render().asPng();
    await mkdir(join(OUT, "og", l.code), { recursive: true });
    await writeFile(join(OUT, "og", l.code, `${r.key}.png`), image);
    count++;
  }
}

const mark = readFileSync("public/favicon.svg", "utf8").replace(/<style>.*<\/style>/, "").replace(/class="g"/g, 'fill="#d97757"');
const png = (size: number, background?: string) => new Resvg(mark, { fitTo: { mode: "width", value: size }, background }).render().asPng();
await writeFile(join(OUT, "icon-512.png"), png(512, "#f5f5f5"));
// iOS ignores SVG icons and fills transparency with black.
await writeFile(join(OUT, "apple-touch-icon.png"), png(180, "#f5f5f5"));
// /favicon.ico for clients that request it without reading the page: an ICO
// directory holding 32 and 48 px PNGs.
const sizes = [32, 48].map((size) => ({ size, data: png(size) }));
const ico = Buffer.alloc(6 + 16 * sizes.length);
ico.writeUInt16LE(1, 2);
ico.writeUInt16LE(sizes.length, 4);
let offset = ico.length;
sizes.forEach(({ size, data }, i) => {
  const e = 6 + 16 * i;
  ico.writeUInt8(size, e);
  ico.writeUInt8(size, e + 1);
  ico.writeUInt16LE(1, e + 4);
  ico.writeUInt16LE(32, e + 6);
  ico.writeUInt32LE(data.length, e + 8);
  ico.writeUInt32LE(offset, e + 12);
  offset += data.length;
});
await writeFile(join(OUT, "favicon.ico"), Buffer.concat([ico, ...sizes.map((s) => s.data)]));
console.log(`og: ${count} images`);
