import type { Metadata } from "next";
import { JsonLd } from "../../../components/JsonLd";
import { Shell } from "../../../components/Shell";
import { getLocale, locales } from "../../../i18n/locales";
import { getMessages, hasMessages } from "../../../lib/i18n";
import { webPageLd } from "../../../lib/jsonld";
import { abs, REPO } from "../../../lib/paths";
import { ogPath, pageTitle, seo } from "../../../lib/seo";
import releases from "../../../data/changelog.json";

type Props = { params: Promise<{ locale: string }> };
type Release = { tag: string; date: string; url: string; notes: string };

const alternates = () => Object.fromEntries(locales.filter((l) => hasMessages(l.code)).map((l) => [l.code, `/${l.code}/changelog/`]));

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const m = getMessages(locale);
  return seo({ locale, title: pageTitle(locale, m.changelog.title), description: m.changelog.description, path: `/${locale}/changelog/`, alternates: alternates(), og: ogPath(locale, "changelog") });
}

export default async function Changelog({ params }: Props) {
  const { locale } = await params;
  const m = getMessages(locale);
  const list = releases as Release[];
  return (
    <Shell locale={locale} alternates={alternates()} active="changelog">
      <div className="wrap page">
        <article className="article">
          <h1>{m.changelog.title}</h1>
          <p className="meta">{m.changelog.description}</p>
          {list.length === 0 ? (
            <p>
              {m.changelog.empty} <a href={REPO}>GitHub</a>
            </p>
          ) : (
            list.map((r) => (
              <section key={r.tag} className="prose">
                <h2 id={r.tag}>
                  <a href={r.url} dir="ltr">{r.tag}</a>
                </h2>
                <p className="meta">{r.date}</p>
                <div dir="ltr" lang="en" dangerouslySetInnerHTML={{ __html: r.notes }} />
              </section>
            ))
          )}
        </article>
      </div>
      {list.length > 0 && (
        <JsonLd data={webPageLd({ title: m.changelog.title, description: m.changelog.description, url: abs(`/${locale}/changelog/`), updated: list[0].date, lang: getLocale(locale)!.lang })} />
      )}
    </Shell>
  );
}
