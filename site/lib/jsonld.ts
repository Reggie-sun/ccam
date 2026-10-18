import { abs, REPO, SITE } from "./paths";

export const VERSION = "0.1.0"; // shown on the site; bump with releases
export const CLAUDE_CODE_VERSION = "2.1.302";

export function websiteLd(lang: string) {
  return [
    { "@context": "https://schema.org", "@type": "WebSite", name: "ccam", url: SITE + "/", inLanguage: lang },
    {
      "@context": "https://schema.org",
      "@type": "Organization",
      name: "ccam",
      url: SITE + "/",
      logo: abs("/icon-512.png"),
      sameAs: [REPO],
    },
  ];
}

export function softwareLd(description: string, lang: string) {
  return {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "ccam",
    description,
    inLanguage: lang,
    applicationCategory: "DeveloperApplication",
    operatingSystem: "Linux, macOS",
    softwareVersion: VERSION,
    downloadUrl: `${REPO}/releases/latest`,
    license: `${REPO}/blob/main/LICENSE`,
    offers: { "@type": "Offer", price: "0", priceCurrency: "USD" },
    url: SITE + "/",
  };
}

export function articleLd(opts: { title: string; description: string; url: string; updated: string; lang: string; crumbs: { name: string; url: string }[] }) {
  return [
    {
      "@context": "https://schema.org",
      "@type": "TechArticle",
      headline: opts.title,
      description: opts.description,
      inLanguage: opts.lang,
      dateModified: opts.updated,
      url: opts.url,
      isPartOf: { "@type": "WebSite", name: "ccam", url: SITE + "/" },
    },
    {
      "@context": "https://schema.org",
      "@type": "BreadcrumbList",
      itemListElement: opts.crumbs.map((c, i) => ({ "@type": "ListItem", position: i + 1, name: c.name, item: c.url })),
    },
  ];
}

export function faqLd(entries: { question: string; answer: string }[], lang: string) {
  return {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    inLanguage: lang,
    mainEntity: entries.map((e) => ({ "@type": "Question", name: e.question, acceptedAnswer: { "@type": "Answer", text: e.answer } })),
  };
}

/** Serializes JSON-LD safely for a <script> element. */
export function ldJson(data: unknown): string {
  return JSON.stringify(data).replace(/</g, "\\u003c");
}
