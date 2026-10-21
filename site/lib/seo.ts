import type { Metadata } from "next";
import { getLocale, locales, type Locale } from "../i18n/locales";
import { getMessages } from "./i18n";
import { abs, SITE } from "./paths";
import { TITLE_MAX, width } from "./width";

type SeoInput = {
  locale: string;
  title: string;
  description: string;
  path: string; // this page's path
  /** Paths of this page in every locale it exists in (code → path). */
  alternates: Record<string, string>;
  og: string; // OG image path
  type?: "website" | "article";
};

/** Canonical, hreflang (both directions, only existing locales plus
 *  x-default) and Open Graph for one page. */
export function seo({ locale, title, description, path, alternates, og, type = "article" }: SeoInput): Metadata {
  const languages: Record<string, string> = {};
  for (const l of locales) {
    if (alternates[l.code]) languages[l.lang] = abs(alternates[l.code]);
  }
  if (alternates.en) languages["x-default"] = abs(alternates.en);
  const current = getLocale(locale)!;
  return {
    metadataBase: new URL(SITE),
    title,
    description,
    alternates: { canonical: abs(path), languages, types: markdownAlternate(path) },
    openGraph: {
      type,
      title,
      description,
      url: abs(path),
      siteName: "ccam",
      locale: current.og,
      alternateLocale: locales.filter((l) => alternates[l.code] && l.code !== locale).map((l) => l.og),
      images: [{ url: abs(og), width: 1200, height: 630, alt: title }],
    },
    twitter: { card: "summary_large_image", title, description, images: [abs(og)] },
  };
}

/** Docs pages also exist as plain Markdown next to the HTML (for LLMs), and
 *  each home page as index.md, the locale's llms.txt (scripts/llms.ts). */
function markdownAlternate(path: string): Record<string, string> | undefined {
  const m = /^\/([^/]+)\/(docs|guides|features|compare|blog)\/([^/]+)\/$/.exec(path);
  if (m) return { "text/markdown": abs(`/${m[1]}/${m[2]}/${m[3]}.md`) };
  return /^\/([a-z-]+\/)?$/.test(path) ? { "text/markdown": abs(`${path}index.md`) } : undefined;
}

/** "<page> | ccam for Claude Code", or "<page> | ccam" when the page title
 *  already names Claude Code or the longer form would be cut off in results. */
export function pageTitle(locale: string, title: string): string {
  const long = `${title} | ${getMessages(locale).site.titleSuffix}`;
  return title.includes("Claude Code") || width(long) > TITLE_MAX ? `${title} | ccam` : long;
}

export function ogPath(locale: Locale | string, key: string): string {
  const code = typeof locale === "string" ? locale : locale.code;
  return `/og/${code}/${key}.png`;
}
