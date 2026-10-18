// HTML validity (html-validate) and JSON-LD syntax on every page.
import { readFileSync } from "node:fs";
import { relative } from "node:path";
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
const errors: string[] = [];
for (const f of (await walk(OUT)).filter((f) => f.endsWith(".html"))) {
  const html = readFileSync(f, "utf8");
  const report = await validator.validateString(html, relative(OUT, f));
  for (const r of report.results) for (const m of r.messages) errors.push(`${r.filePath}:${m.line} ${m.ruleId}: ${m.message}`);
  for (const [, json] of html.matchAll(/<script type="application\/ld\+json">([\s\S]*?)<\/script>/g)) {
    try {
      const data = JSON.parse(json);
      for (const item of [data].flat()) if (!item["@context"] || !item["@type"]) errors.push(`${relative(OUT, f)}: JSON-LD without @context/@type`);
    } catch (e) {
      errors.push(`${relative(OUT, f)}: invalid JSON-LD (${(e as Error).message})`);
    }
  }
}
fail(errors, "html");
