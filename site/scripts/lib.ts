// Helpers shared by the build and check scripts (plain Node, no Vite).
import { readdirSync, readFileSync, existsSync } from "node:fs";
import { readdir } from "node:fs/promises";
import { join } from "node:path";
import { locales, type Locale } from "../i18n/locales.ts";
import { splitFrontMatter } from "../lib/frontmatter.ts";
import { homePath, pagePath } from "../lib/paths.ts";

export const OUT = "dist/client";

export type SourcePage = {
  locale: string;
  section: string;
  slug: string;
  title: string;
  description: string;
  updated: string;
  order: number;
  translation?: string;
  sourceHash?: string;
  body: string;
  path: string; // URL path
  file: string; // content file
};

export function builtLocales(): Locale[] {
  return locales.filter((l) => existsSync(`messages/${l.code}.json`));
}

export function messages(code: string): any {
  return JSON.parse(readFileSync(`messages/${code}.json`, "utf8"));
}

export function sourcePages(): SourcePage[] {
  const out: SourcePage[] = [];
  for (const locale of readdirSync("content")) {
    for (const section of readdirSync(join("content", locale))) {
      for (const f of readdirSync(join("content", locale, section))) {
        if (!f.endsWith(".md")) continue;
        const file = join("content", locale, section, f);
        const { data, body } = splitFrontMatter<any>(readFileSync(file, "utf8"));
        const slug = f.slice(0, -3);
        out.push({ ...data, updated: String(data.updated), locale, section, slug, body, path: pagePath(locale, section, slug), file });
      }
    }
  }
  return out.filter((p) => existsSync(`messages/${p.locale}.json`)).sort((a, b) => a.order - b.order);
}

export type Route = { locale: string; path: string; key: string; updated: string; alternates: Record<string, string> };

/** Every page the site should contain, with its alternates (code → path). */
export function routes(): Route[] {
  const built = builtLocales();
  const pages = sourcePages();
  const latest = pages.map((p) => p.updated).sort().at(-1) ?? "2026-01-01";
  const all = (fn: (code: string) => string) => Object.fromEntries(built.map((l) => [l.code, fn(l.code)]));
  const out: Route[] = [];
  for (const l of built) {
    out.push({ locale: l.code, path: homePath(l.code), key: "home", updated: latest, alternates: all(homePath) });
    out.push({ locale: l.code, path: `/${l.code}/changelog/`, key: "changelog", updated: latest, alternates: all((c) => `/${c}/changelog/`) });
    out.push({ locale: l.code, path: `/${l.code}/languages/`, key: "languages", updated: latest, alternates: all((c) => `/${c}/languages/`) });
  }
  for (const p of pages) {
    const alternates = Object.fromEntries(
      pages.filter((q) => q.section === p.section && q.slug === p.slug).map((q) => [q.locale, q.path]),
    );
    out.push({ locale: p.locale, path: p.path, key: `${p.section}-${p.slug}`, updated: p.updated, alternates });
  }
  return out;
}

export function htmlFile(path: string): string {
  return join(OUT, path, "index.html");
}

export async function walk(dir: string): Promise<string[]> {
  const out: string[] = [];
  for (const e of await readdir(dir, { withFileTypes: true })) {
    const p = join(dir, e.name);
    if (e.isDirectory()) out.push(...(await walk(p)));
    else out.push(p);
  }
  return out;
}

export function fail(errors: string[], what: string): void {
  if (errors.length) {
    console.error(`${what}: ${errors.length} problem(s)`);
    for (const e of errors.slice(0, 50)) console.error("  " + e);
    process.exitCode = 1;
  } else {
    console.log(`${what}: ok`);
  }
}
