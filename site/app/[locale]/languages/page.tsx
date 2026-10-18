import type { Metadata } from "next";
import { Shell } from "../../../components/Shell";
import { locales, nativeName } from "../../../i18n/locales";
import { format, getMessages, hasMessages } from "../../../lib/i18n";
import { homePath } from "../../../lib/paths";
import { ogPath, seo } from "../../../lib/seo";
import coverage from "../../../data/coverage.json";

type Props = { params: Promise<{ locale: string }> };

const built = () => locales.filter((l) => hasMessages(l.code));
const alternates = () => Object.fromEntries(built().map((l) => [l.code, `/${l.code}/languages/`]));

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const m = getMessages(locale);
  return seo({ locale, title: `${m.languages.title} | ccam`, description: format(locale, m.languages.description, { count: built().length }), path: `/${locale}/languages/`, alternates: alternates(), og: ogPath(locale, "languages") });
}

export default async function Languages({ params }: Props) {
  const { locale } = await params;
  const m = getMessages(locale);
  const pct = (v: number) => new Intl.NumberFormat(locale, { style: "percent", maximumFractionDigits: 0 }).format(v);
  return (
    <Shell locale={locale} alternates={alternates()}>
      <div className="wrap docs">
        <div />
        <article className="article prose">
          <h1>{m.languages.title}</h1>
          <p>{format(locale, m.languages.description, { count: built().length })}</p>
          {coverage.high > 0 && (
            <>
              <p>{format(locale, m.languages.coverage, { low: pct(coverage.low), high: pct(coverage.high) })}</p>
              <p>{m.languages.method}</p>
            </>
          )}
          <ul>
            {built().map((l) => (
              <li key={l.code}>
                <a href={homePath(l.code)} hrefLang={l.lang} lang={l.lang}>{nativeName(l)}</a>
              </li>
            ))}
          </ul>
        </article>
      </div>
    </Shell>
  );
}
