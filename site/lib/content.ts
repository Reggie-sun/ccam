// A page exists in a locale only when its file does: untranslated pages are
// not generated and never fall back to English (docs/WEBSITE_PLAN.md 5.3).
import { Marked, type Tokens } from "marked";
import { splitFrontMatter } from "./frontmatter";
import { locales, type Locale } from "../i18n/locales";
import { pagePath } from "./paths";

export type Section = "docs" | "guides" | "compare" | "blog";
export const sections: Section[] = ["docs", "guides", "compare", "blog"];

export type FrontMatter = {
  title: string;
  description: string;
  updated: string; // YYYY-MM-DD
  order: number;
  ccam?: string;
  claude?: string;
  translation?: "source" | "machine" | "reviewed";
  sourceHash?: string; // hash of the English source this translation was made from
};

export type Page = FrontMatter & {
  locale: string;
  section: Section;
  slug: string;
  markdown: string; // body without front matter
  file: string; // repository-relative path, for edit links
};

const raw = import.meta.glob<string>("../content/*/*/*.md", { query: "?raw", import: "default", eager: true });

const all: Page[] = Object.entries(raw).map(([path, source]) => {
  const [, locale, section, file] = /content\/([^/]+)\/([^/]+)\/([^/]+)\.md$/.exec(path)!;
  const { data, body } = splitFrontMatter<FrontMatter>(source);
  return { ...data, updated: String(data.updated), locale, section: section as Section, slug: file, markdown: body, file: `site/content/${locale}/${section}/${file}.md` };
});

export function getPages(locale: string, section?: Section): Page[] {
  return all
    .filter((p) => p.locale === locale && (!section || p.section === section))
    .sort((a, b) => a.order - b.order || a.slug.localeCompare(b.slug));
}

export function getPage(locale: string, section: Section, slug: string): Page | undefined {
  return all.find((p) => p.locale === locale && p.section === section && p.slug === slug);
}

/** Locales in which a page exists, in locale-list order. */
export function pageLocales(section: Section, slug: string): Locale[] {
  return locales.filter((l) => getPage(l.code, section, slug));
}

export function allPages(): Page[] {
  return all;
}

export type Heading = { id: string; text: string; depth: number };

function slugify(text: string): string {
  return text.toLowerCase().normalize("NFKC").replace(/<[^>]+>/g, "").replace(/[^\p{L}\p{N}]+/gu, "-").replace(/^-|-$/g, "");
}

/** Internal links go to the same locale when the page exists there, else to
 *  English. Code stays left-to-right inside right-to-left pages. */
export function renderMarkdown(page: Pick<Page, "locale" | "markdown">, labels: { copy: string; copied: string }): { html: string; headings: Heading[] } {
  const headings: Heading[] = [];
  const marked = new Marked({
    gfm: true,
    renderer: {
      heading({ tokens, depth }: Tokens.Heading) {
        const text = this.parser.parseInline(tokens);
        const id = slugify(text);
        headings.push({ id, text: text.replace(/<[^>]+>/g, ""), depth });
        return `<h${depth} id="${id}">${text}</h${depth}>\n`;
      },
      code({ text, lang }: Tokens.Code) {
        const escaped = text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
        const copy = lang === "sh" || lang === "bash" || lang === "console";
        return `<div class="code"><pre dir="ltr" tabindex="0"><code${lang ? ` data-lang="${lang}"` : ""}>${escaped}</code></pre>${copy ? `<button type="button" class="copy" data-copy data-copied="${labels.copied}" hidden>${labels.copy}</button>` : ""}</div>\n`;
      },
      codespan({ text }: Tokens.Codespan) {
        // marked passes inline code unescaped: "ccs <name>" must not become a tag.
        const escaped = text.replace(/&(?![a-z]+;|#\d+;)/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
        return `<code dir="ltr">${escaped}</code>`;
      },
      link({ href, title, tokens }: Tokens.Link) {
        const text = this.parser.parseInline(tokens);
        let target = href;
        const internal = /^\/(docs|guides|compare|blog)\/([a-z0-9-]+)\/?(#.*)?$/.exec(href);
        if (internal) {
          const [, section, slug, hash = ""] = internal;
          const locale = getPage(page.locale, section as Section, slug) ? page.locale : "en";
          target = pagePath(locale, section, slug) + hash;
        }
        const ext = /^https?:/.test(target);
        return `<a href="${target}"${title ? ` title="${title}"` : ""}${ext ? ' rel="noopener"' : ""}>${text}</a>`;
      },
    },
  });
  const html = (marked.parse(page.markdown, { async: false }) as string).replace(/<th>/g, '<th scope="col">').replace(/<table>/g, '<table tabindex="0">');
  return { html, headings };
}

/** Question/answer pairs from a FAQ page: each "## question" and the text below it. */
export function faqEntries(page: Page): { question: string; answer: string }[] {
  const out: { question: string; answer: string }[] = [];
  const parts = page.markdown.split(/^## /m).slice(1);
  for (const part of parts) {
    const [question, ...rest] = part.split("\n");
    const answer = rest.join("\n").replace(/```[\s\S]*?```/g, "").replace(/[`*_>#]/g, "").replace(/\[([^\]]+)\]\([^)]+\)/g, "$1").replace(/\s+/g, " ").trim();
    out.push({ question: question.trim(), answer });
  }
  return out;
}
