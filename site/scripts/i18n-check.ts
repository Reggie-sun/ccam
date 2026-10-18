// Translation checks (docs/WEBSITE_PLAN.md 5.4, 11):
// - messages/<locale>.json has exactly the keys, array lengths and ICU
//   placeholders of messages/en.json
// - each translated page keeps every code block, internal link and the
//   commands marker of its English source, and the glossary terms
// - front matter is complete; sourceHash must match the current English file
//   (stale translations fail with --strict, warn otherwise)
// Options: --locale <code>   check one locale
//          --strict          stale translations are errors
//          --hash <file>     print the sourceHash of an English page
import { createHash } from "node:crypto";
import { existsSync, readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { parse as parseYaml } from "yaml";
import { parse, type MessageFormatElement } from "@formatjs/icu-messageformat-parser";
import { locales } from "../i18n/locales.ts";
import { splitFrontMatter } from "../lib/frontmatter.ts";
import { fail } from "./lib.ts";

export const sourceHash = (file: string) => createHash("sha256").update(readFileSync(file, "utf8")).digest("hex").slice(0, 12);

const args = process.argv.slice(2);
if (args[0] === "--hash") {
  console.log(sourceHash(args[1]));
  process.exit(0);
}
const only = args.includes("--locale") ? args[args.indexOf("--locale") + 1] : undefined;
const strict = args.includes("--strict");
const glossary: string[] = parseYaml(readFileSync("i18n/glossary.yaml", "utf8")).keep;

const errors: string[] = [];
const warnings: string[] = [];
// Argument names from the real ICU parser, so plural branch text is not
// mistaken for a placeholder.
function placeholders(s: string): string {
  const names = new Set<string>();
  const walk = (els: MessageFormatElement[]) => {
    for (const el of els) {
      if ("value" in el && typeof el.value === "string" && el.type !== 0 && el.type !== 7) names.add(el.value);
      if ("options" in el) for (const o of Object.values(el.options)) walk((o as { value: MessageFormatElement[] }).value);
    }
  };
  try {
    walk(parse(s));
  } catch {
    return "(unparseable)";
  }
  return [...names].sort().join(",");
}

function compare(en: unknown, tr: unknown, path: string, code: string): void {
  if (typeof en === "string") {
    if (typeof tr !== "string") return void errors.push(`${code}: ${path} should be a string`);
    if (/^landing\.flow\[\d+\]\.(command|also)$/.test(path) && tr !== en) errors.push(`${code}: ${path} is a command and must stay "${en}"`);
    if (placeholders(en) !== placeholders(tr)) errors.push(`${code}: ${path} placeholders {${placeholders(en)}} ≠ {${placeholders(tr)}}`);
    if (!tr.trim()) errors.push(`${code}: ${path} is empty`);
    return;
  }
  if (Array.isArray(en)) {
    if (!Array.isArray(tr) || tr.length !== en.length) return void errors.push(`${code}: ${path} should have ${en.length} items`);
    en.forEach((v, i) => compare(v, tr[i], `${path}[${i}]`, code));
    return;
  }
  if (en && typeof en === "object") {
    if (!tr || typeof tr !== "object") return void errors.push(`${code}: ${path} should be an object`);
    for (const k of Object.keys(en)) compare((en as any)[k], (tr as any)[k], path ? `${path}.${k}` : k, code);
    for (const k of Object.keys(tr)) if (!(k in (en as object))) errors.push(`${code}: unknown key ${path ? path + "." : ""}${k}`);
  }
}

const codeBlocks = (md: string) => [...md.matchAll(/```[\s\S]*?```/g)].map((m) => m[0]);
const inlineCode = (md: string) => [...md.replace(/```[\s\S]*?```/g, "").matchAll(/`([^`]+)`/g)].map((m) => m[1]).sort();
const links = (md: string) => [...md.matchAll(/\]\(([^)]+)\)/g)].map((m) => m[1]).sort();

const en = JSON.parse(readFileSync("messages/en.json", "utf8"));
for (const l of locales.filter((l) => l.code !== "en" && (!only || l.code === only))) {
  const mfile = `messages/${l.code}.json`;
  if (existsSync(mfile)) compare(en, JSON.parse(readFileSync(mfile, "utf8")), "", l.code);
  const dir = join("content", l.code);
  if (!existsSync(dir)) continue;
  if (!existsSync(mfile)) errors.push(`${l.code}: content without messages/${l.code}.json`);
  for (const section of readdirSync(dir)) {
    for (const f of readdirSync(join(dir, section))) {
      const file = join(dir, section, f);
      const srcFile = join("content/en", section, f);
      if (!existsSync(srcFile)) {
        errors.push(`${file}: no English source`);
        continue;
      }
      const src = splitFrontMatter<any>(readFileSync(srcFile, "utf8"));
      const tr = splitFrontMatter<any>(readFileSync(file, "utf8"));
      for (const k of ["title", "description", "updated", "order", "translation", "sourceHash"]) if (tr.data[k] === undefined) errors.push(`${file}: front matter lacks ${k}`);
      if (tr.data.order !== src.data.order) errors.push(`${file}: order ${tr.data.order} ≠ ${src.data.order}`);
      if (codeBlocks(src.body).join("\n") !== codeBlocks(tr.body).join("\n")) errors.push(`${file}: code blocks differ from the English source`);
      if (inlineCode(src.body).join("\n") !== inlineCode(tr.body).join("\n")) errors.push(`${file}: inline code differs (${inlineCode(src.body).length} vs ${inlineCode(tr.body).length})`);
      if (links(src.body).join(" ") !== links(tr.body).join(" ")) errors.push(`${file}: links differ: ${links(src.body).join(" ")} ≠ ${links(tr.body).join(" ")}`);
      const levels = (md: string) => [...md.replace(/```[\s\S]*?```/g, "").matchAll(/^(#{1,6}) /gm)].map((m) => m[1].length).join(",");
      if (levels(src.body) !== levels(tr.body)) errors.push(`${file}: headings differ from the English source (anchors are matched by position)`);
      if (src.body.includes("<!-- commands -->") !== tr.body.includes("<!-- commands -->")) errors.push(`${file}: <!-- commands --> marker missing`);
      const visible = (s: string) => s.replace(/```[\s\S]*?```/g, "").replace(/`[^`]+`/g, "");
      for (const term of glossary) {
        if (visible(src.body + src.data.title).includes(term) && !(tr.body + tr.data.title).includes(term)) errors.push(`${file}: glossary term "${term}" missing`);
      }
      const ratio = visible(tr.body).length / Math.max(1, visible(src.body).length);
      if (ratio < 0.3 || ratio > 3.2) warnings.push(`${file}: length ratio ${ratio.toFixed(2)}`);
      const want = sourceHash(srcFile);
      if (tr.data.sourceHash !== want) (strict ? errors : warnings).push(`${file}: stale (sourceHash ${tr.data.sourceHash}, English is ${want})`);
    }
  }
}
for (const w of warnings) console.warn("warning: " + w);
fail(errors, "i18n");
