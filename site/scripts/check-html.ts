// HTML validity (html-validate) and JSON-LD syntax on every page.
import { readFileSync } from "node:fs";
import { join, relative } from "node:path";
import { HtmlValidate } from "html-validate";
import { fail, OUT, walk } from "./lib.ts";

const validator = new HtmlValidate({
  extends: ["html-validate:recommended"],
  rules: {
    "no-inline-style": "off", // the CSP hashes the few style attributes
    "attribute-boolean-style": "off",
    "attribute-empty-style": "off",
    "no-trailing-whitespace": "off",
    "void-style": "off",
    "long-title": "off",
    "attr-case": "off", // React writes hrefLang, charSet
    "no-raw-characters": "off",
    "prefer-native-element": "off",
    "no-redundant-role": "off",
    "valid-id": ["error", { relaxed: true }],
  },
});
// Search results cut titles and descriptions by width; CJK and other
// full-width characters take about two Latin character widths.
const width = (s: string) => [...s].reduce((n, c) => n + (/[ᄀ-ᅟ⺀-꓏가-힣豈-﫿︰-﹏＀-｠]/.test(c) ? 2 : 1), 0);
const TITLE_MAX = 60;
const DESCRIPTION = [70, 160];
const decode = (s: string) => s.replace(/&quot;/g, '"').replace(/&#x27;|&#39;/g, "'").replace(/&lt;/g, "<").replace(/&gt;/g, ">").replace(/&amp;/g, "&");

const errors: string[] = [];
for (const f of (await walk(OUT)).filter((f) => f.endsWith(".html"))) {
  const html = readFileSync(f, "utf8");
  const page = relative(OUT, f);
  const report = await validator.validateString(html, page);
  for (const r of report.results) for (const m of r.messages) errors.push(`${r.filePath}:${m.line} ${m.ruleId}: ${m.message}`);

  if (!/name="robots" content="noindex/.test(html)) {
    const title = decode(/<title>([^<]*)<\/title>/.exec(html)?.[1] ?? "");
    const description = decode(/<meta name="description" content="([^"]*)"/.exec(html)?.[1] ?? "");
    if (width(title) > TITLE_MAX) errors.push(`${page}: title is ${width(title)} wide (max ${TITLE_MAX}): ${title}`);
    if (width(description) < DESCRIPTION[0] || width(description) > DESCRIPTION[1]) errors.push(`${page}: description is ${width(description)} wide (${DESCRIPTION.join("-")}): ${description}`);
  }

  for (const [, json] of html.matchAll(/<script type="application\/ld\+json">([\s\S]*?)<\/script>/g)) {
    let data;
    try {
      data = JSON.parse(json);
    } catch (e) {
      errors.push(`${page}: invalid JSON-LD (${(e as Error).message})`);
      continue;
    }
    for (const item of [data].flat()) {
      if (!item["@context"] || !item["@type"]) errors.push(`${page}: JSON-LD without @context/@type`);
      if (item["@type"] === "BreadcrumbList") {
        const urls = item.itemListElement.map((i: { item: string }) => i.item);
        if (new Set(urls).size !== urls.length) errors.push(`${page}: breadcrumb repeats a URL`);
      }
      if (item["@type"] === "TechArticle" || item["@type"] === "BlogPosting") {
        for (const k of ["author", "publisher", "datePublished", "dateModified"]) if (!item[k]) errors.push(`${page}: ${item["@type"]} without ${k}`);
      }
      if (item["@type"] === "FAQPage") {
        // Identifiers such as CLAUDE_CONFIG_DIR must survive the Markdown-to-text step.
        const md = readFileSync(join("content", page.split("/")[0], "docs", "faq.md"), "utf8").replace(/```[\s\S]*?```/g, "");
        const text = item.mainEntity.map((q: { name: string; acceptedAnswer: { text: string } }) => `${q.name}\n${q.acceptedAnswer.text}`).join("\n");
        for (const [, code] of md.matchAll(/`([^`]+)`/g)) if (!text.includes(code)) errors.push(`${page}: FAQ answer text lost "${code}"`);
      }
    }
  }
}
fail(errors, "html");
