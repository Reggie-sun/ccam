import { getLocale, locales } from "../i18n/locales";
import { getPages } from "../lib/content";
import { format, getMessages, hasMessages } from "../lib/i18n";
import { homePath, pagePath } from "../lib/paths";
import { softwareLd } from "../lib/jsonld";
import { InstallCommand } from "./InstallCommand";
import { JsonLd } from "./JsonLd";
import { Shell } from "./Shell";
import { UsageBoard } from "./UsageBoard";

export function landingAlternates(): Record<string, string> {
  return Object.fromEntries(locales.filter((l) => hasMessages(l.code)).map((l) => [l.code, homePath(l.code)]));
}

export function Landing({ locale }: { locale: string }) {
  const l = getLocale(locale)!;
  const m = getMessages(locale);
  const guides = getPages(locale, "guides");
  const add = <code dir="ltr">ccam add</code>;
  const [before, after] = m.landing.thenAdd.split("{command}");
  return (
    <Shell locale={locale} alternates={landingAlternates()}>
      <div className="wrap">
        <section className="hero">
          <div>
            <h1>{m.landing.headline}</h1>
            <p className="lede">{m.landing.lede}</p>
            <p className="install-label">{m.landing.installLabel}</p>
            <InstallCommand copy={m.copy.copy} copied={m.copy.copied} />
            <p className="then">
              {before}
              {add}
              {after} <a href={pagePath(locale, "docs", "quick-start")}>{m.landing.readDocs}</a>
            </p>
          </div>
          <UsageBoard locale={l} m={m} />
        </section>

        <section className="section" aria-labelledby="flow">
          <h2 id="flow">{m.landing.flowTitle}</h2>
          <dl className="flow">
            {m.landing.flow.map((f) => (
              <div key={f.command}>
                <dt dir="ltr">
                  <code>{f.command}</code>
                  {"also" in f && f.also ? <code>{f.also}</code> : null}
                </dt>
                <dd>{f.body}</dd>
              </div>
            ))}
          </dl>
        </section>

        {guides.length > 0 && (
          <section className="section" aria-labelledby="guides">
            <h2 id="guides">{m.landing.guidesTitle}</h2>
            <ul className="guide-list">
              {guides.map((g) => (
                <li key={g.slug}>
                  <a href={pagePath(locale, "guides", g.slug)}>{g.title}</a>
                  <p>{g.description}</p>
                </li>
              ))}
            </ul>
          </section>
        )}

        <section className="section" aria-labelledby="facts">
          <h2 id="facts">{m.landing.factsTitle}</h2>
          <div className="facts">
            <p>{m.facts.what}</p>
            <p>{m.facts.how}</p>
            <p>{m.facts.install}</p>
            <p>{m.facts.official}</p>
          </div>
        </section>
      </div>
      <JsonLd data={softwareLd(m.site.description, l.lang)} />
    </Shell>
  );
}

export function landingTitle(locale: string): { title: string; description: string } {
  const m = getMessages(locale);
  return { title: m.site.title, description: m.site.description };
}
