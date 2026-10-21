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

// The commands an LLM needs first, in the order people use them; the rest
// are in the command reference.
const essentials = ["add", "name", "ccam", "list", "usage", "next", "fallback", "pick", "group", "folder", "doctor", "update"];

/** How to use ccam at a glance: install, then the essential commands with
 *  the locale's descriptions, so a reader need not open any page first. */
function quickUsage(code: string): string {
  const m = messages(code);
  const lines = [
    "- `curl -fsSL https://ccam.dev/install.sh | sh` (Linux, macOS) or `irm https://ccam.dev/install.ps1 | iex` (Windows PowerShell): install ccam; Claude Code must be installed first",
    ...essentials.map((key) => {
      const c = commands.find((c) => c.key === key)!;
      return `- \`${c.usage}\`: ${m.commands[key]}`;
    }),
  ];
  return [
    `Documentation in one file: ${abs("/llms-full.txt")}`,
    "Every page listed below is plain Markdown; its HTML version is the same address without `.md`.",
    `Main commands (all of them: ${abs(`/${code}/docs/commands.md`)}):`,
    lines.join("\n"),
  ].join("\n\n");
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
    "# ccam (Claude Code Account Manager)",
    `> ${m.facts.what}`,
    [m.facts.how, m.facts.install + " `curl -fsSL https://ccam.dev/install.sh | sh`", m.facts.official].map((f) => `- ${f}`).join("\n"),
    `Website: ${abs(homePath(l.code))}\nSource code: https://github.com/jhyunwoo/ccam`,
    quickUsage(l.code),
    ...groups.filter(Boolean),
  ].join("\n\n");
  // llms.txt, and the same text as the home page's Markdown version.
  const dir = l.code === "en" ? OUT : join(OUT, l.code);
  await mkdir(dir, { recursive: true });
  await writeFile(join(dir, "llms.txt"), text + "\n");
  await writeFile(join(dir, "index.md"), text + "\n");
}

const full = pages.filter((p) => p.locale === "en").map(toMarkdown).join("\n---\n\n");
await writeFile(join(OUT, "llms-full.txt"), full);
console.log(`llms: ${pages.length} Markdown pages, ${builtLocales().length} llms.txt`);
