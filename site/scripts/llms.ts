// llms.txt (llmstxt.org), llms-full.txt and a plain Markdown copy of every page.
import { readFileSync } from "node:fs";
import { mkdir, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import { nav } from "../lib/nav.ts";
import { abs, homePath } from "../lib/paths.ts";
import { builtLocales, messages, OUT, sourcePages, type SourcePage } from "./lib.ts";

const commands: { key: string; usage: string }[] = JSON.parse(readFileSync("data/commands.json", "utf8"));

function commandTable(code: string): string {
  const m = messages(code);
  const rows = commands.map((c) => `| \`${c.usage}\` | ${m.commands[c.key] ?? ""} |`);
  return [`| ${m.docs.commandsTable.command} | ${m.docs.commandsTable.does} |`, "| --- | --- |", ...rows].join("\n");
}

/** Absolute links, so the Markdown works outside the site. */
function toMarkdown(p: SourcePage): string {
  const body = p.body
    .replace("<!-- commands -->", commandTable(p.locale))
    .replace(/<!-- motion:[a-z0-9-]+ -->\n*/g, "") // diagrams are pictures; the text explains them
    .replace(/\]\(\/(docs|guides|features|compare|blog)\/([a-z0-9-]+)\/?(#[^)]*)?\)/g, (_, s, slug, hash = "") => `](${abs(`/${p.locale}/${s}/${slug}.md`)}${hash})`);
  return `# ${p.title}\n\n> ${p.description}\n\nSource: ${abs(p.path)} · Updated ${p.updated}\n\n${body.trim()}\n`;
}

const pages = sourcePages();
for (const p of pages) {
  const file = join(OUT, p.locale, p.section, `${p.slug}.md`);
  await mkdir(dirname(file), { recursive: true });
  await writeFile(file, toMarkdown(p));
}

for (const l of builtLocales()) {
  const m = messages(l.code);
  const list = pages.filter((p) => p.locale === l.code);
  // Same groups and order as the docs sidebar.
  const groups = nav.map((g) => {
    const items = g.pages
      .map((id) => list.find((p) => `${p.section}/${p.slug}` === id))
      .filter((p) => !!p)
      .map((p) => `- [${p.title}](${abs(`/${l.code}/${p.section}/${p.slug}.md`)}): ${p.description}`);
    return items.length ? `## ${m.docs.groups[g.key]}\n\n${items.join("\n")}` : "";
  });
  const text = [
    "# ccam",
    `> ${m.facts.what}`,
    [m.facts.how, m.facts.install + " `curl -fsSL https://ccam.dev/install.sh | sh`", m.facts.official].map((f) => `- ${f}`).join("\n"),
    `Website: ${abs(homePath(l.code))}\nSource code: https://github.com/jhyunwoo/ccam`,
    ...groups.filter(Boolean),
  ].join("\n\n");
  const file = l.code === "en" ? join(OUT, "llms.txt") : join(OUT, l.code, "llms.txt");
  await mkdir(dirname(file), { recursive: true });
  await writeFile(file, text + "\n");
}

const full = pages.filter((p) => p.locale === "en").map(toMarkdown).join("\n---\n\n");
await writeFile(join(OUT, "llms-full.txt"), full);
console.log(`llms: ${pages.length} Markdown pages, ${builtLocales().length} llms.txt`);
