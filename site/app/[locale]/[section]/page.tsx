import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { JsonLd } from "../../../components/JsonLd";
import { Shell } from "../../../components/Shell";
import { getLocale, locales } from "../../../i18n/locales";
import { getPages, sections, type Section } from "../../../lib/content";
import { getMessages, hasMessages } from "../../../lib/i18n";
import { webPageLd } from "../../../lib/jsonld";
import { abs, pagePath } from "../../../lib/paths";
import { ogPath, pageTitle, seo } from "../../../lib/seo";

// A list of every page in one section, such as /en/guides/.
type Props = { params: Promise<{ locale: string; section: string }> };

const alternates = (section: Section) =>
  Object.fromEntries(locales.filter((l) => hasMessages(l.code) && getPages(l.code, section).length).map((l) => [l.code, pagePath(l.code, section)]));

export const dynamicParams = false;
export function generateStaticParams() {
  return locales
    .filter((l) => hasMessages(l.code))
    .flatMap((l) => sections.filter((s) => getPages(l.code, s).length).map((section) => ({ locale: l.code, section })));
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale, section } = await params;
  const m = getMessages(locale);
  const s = section as Section;
  return seo({ locale, title: pageTitle(locale, m.docs.sections[s]), description: m.docs.hubs[s], path: pagePath(locale, s), alternates: alternates(s), og: ogPath(locale, `hub-${s}`) });
}

export default async function SectionHub({ params }: Props) {
  const { locale, section } = await params;
  const s = section as Section;
  const pages = getPages(locale, s);
  if (!pages.length) notFound();
  const m = getMessages(locale);
  return (
    <Shell locale={locale} alternates={alternates(s)} active={s === "docs" || s === "guides" || s === "features" ? "docs" : undefined}>
      <div className="wrap page">
        <article className="article prose">
          <h1>{m.docs.sections[s]}</h1>
          <p>{m.docs.hubs[s]}</p>
          <ul>
            {pages.map((p) => (
              <li key={p.slug}>
                <a href={pagePath(locale, s, p.slug)}>{p.title}</a>: {p.description}
              </li>
            ))}
          </ul>
        </article>
      </div>
      <JsonLd
        data={webPageLd({
          title: m.docs.sections[s],
          description: m.docs.hubs[s],
          url: abs(pagePath(locale, s)),
          updated: pages.map((p) => p.updated).sort().at(-1)!,
          lang: getLocale(locale)!.lang,
        })}
      />
    </Shell>
  );
}
