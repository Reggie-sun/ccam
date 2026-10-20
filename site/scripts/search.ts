// The search index: one JSON file per locale at /search/<locale>.json, read by
// scripts/client/site.js the first time the search dialog opens. Pages are in
// sidebar order and split at each heading, so a result links to its section.
import { readFileSync } from "node:fs";
import { mkdir, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { gzipSync } from "node:zlib";
import { Marked, type Token } from "marked";
import { headingIds, slugify } from "../lib/headings.ts";
import { nav } from "../lib/nav.ts";
import { builtLocales, messages, OUT, sourcePages, type SourcePage } from "./lib.ts";

type Section = { id: string; heading: string; text: string };
type Entry = { url: string; title: string; group: string; sections: Section[] };

const commands: { key: string; usage: string }[] = JSON.parse(readFileSync("data/commands.json", "utf8"));

/** Readable text of a block: code as written, Markdown syntax removed. */
function plain(token: Token): string {
  if (token.type === "code") return token.text;
  if (token.type === "table") return [token.header, ...token.rows].flat().map((cell: { text: string }) => inline(cell.text)).join(" ");
  return inline(token.raw.replace(/^\s*(?:[-*+]|\d+\.)\s+/gm, "").replace(/^\s*>\s?/gm, ""));
}

function inline(md: string): string {
  return md.replace(/<!--[\s\S]*?-->/g, " ").replace(/!?\[([^\]]*)\]\([^)]*\)/g, "$1").replace(/`|\*\*/g, "");
}

function sections(p: SourcePage, english: SourcePage | undefined): Section[] {
  const m = messages(p.locale);
  const body = p.body.replace("<!-- commands -->", commands.map((c) => `${c.usage} ${m.commands[c.key] ?? ""}`).join("\n"));
  const ids = english ? headingIds(english.body) : [];
  const out: Section[] = [{ id: "", heading: "", text: "" }];
  let n = 0;
  for (const token of new Marked().lexer(body)) {
    if (token.type === "heading") {
      const heading = token.text.replace(/`/g, "");
      out.push({ id: ids[n++] ?? slugify(heading), heading, text: "" });
    } else {
      out[out.length - 1].text += " " + plain(token);
    }
  }
  for (const s of out) s.text = s.text.replace(/\s+/g, " ").trim();
  return out.filter((s) => s.heading || s.text);
}

const pages = sourcePages();
await mkdir(join(OUT, "search"), { recursive: true });
let largest = 0;
for (const l of builtLocales()) {
  const m = messages(l.code);
  const entries: Entry[] = [];
  for (const g of nav) {
    for (const id of g.pages) {
      const p = pages.find((q) => q.locale === l.code && `${q.section}/${q.slug}` === id);
      if (!p) continue;
      const english = pages.find((q) => q.locale === "en" && q.section === p.section && q.slug === p.slug);
      entries.push({ url: p.path, title: p.title, group: m.docs.groups[g.key], sections: sections(p, english) });
    }
  }
  const json = JSON.stringify({ pages: entries });
  largest = Math.max(largest, gzipSync(json).length);
  await writeFile(join(OUT, "search", `${l.code}.json`), json);
}
console.log(`search: ${builtLocales().length} indexes, largest ${largest} B gzip`);
