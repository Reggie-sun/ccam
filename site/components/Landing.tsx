import { Fragment } from "react";
import { getLocale, locales } from "../i18n/locales";
import { getPages } from "../lib/content";
import { getMessages, hasMessages } from "../lib/i18n";
import { homePath, pagePath } from "../lib/paths";
import { softwareLd } from "../lib/jsonld";
import { ArrowIcon } from "./Icons";
import { INSTALL_WINDOWS, InstallCommand } from "./InstallCommand";
import { JsonLd } from "./JsonLd";
import { Shell } from "./Shell";
import { UsageBoard } from "./UsageBoard";

// The feature page behind each "How it works" step, in the order of landing.flow.
const flowFeatures = ["accounts", "launch", "share", "usage", "next", "doctor"];

export function landingAlternates(): Record<string, string> {
  return Object.fromEntries(locales.filter((l) => hasMessages(l.code)).map((l) => [l.code, homePath(l.code)]));
}

export function Landing({ locale }: { locale: string }) {
  const l = getLocale(locale)!;
  const m = getMessages(locale);
  const guides = getPages(locale, "guides");
  const add = <code dir="ltr">ccam add</code>;
  const [before, after] = m.landing.thenAdd.split("{command}");
  // Each word rises out of its own clip; languages without spaces rise as one line.
  // "Claude Code" rises as one word so it never breaks across lines.
  const words = m.landing.headline.replace("Claude Code", "Claude\u00a0Code").split(" ");
  return (
    <Shell locale={locale} alternates={landingAlternates()}>
      <section className="hero">
        <div className="wrap">
          <h1 className="display">
            {words.map((w, i) => (
              <Fragment key={i}>
                {i > 0 && " "}
                <span className="w">
                  <span>{w}</span>
                </span>
              </Fragment>
            ))}
          </h1>
          <p className="lede focus-in d1">{m.landing.lede}</p>
          <div className="focus-in d2">
            <p className="install-label">
              <a href={pagePath(locale, "docs", "install")}>{m.landing.installLabel}</a>
            </p>
            <InstallCommand copy={m.copy.copy} copied={m.copy.copied} />
            <p className="install-label">{m.landing.installWindowsLabel}</p>
            <InstallCommand copy={m.copy.copy} copied={m.copy.copied} command={INSTALL_WINDOWS} />
          </div>
          <p className="then focus-in d3">
            <span>
              {before}
              {add}
              {after}
            </span>
            <a className="pill pill-outline" href={pagePath(locale, "docs", "quick-start")}>{m.landing.readDocs}</a>
          </p>
        </div>
      </section>

      <section className="section demo-section" aria-labelledby="demo">
        <div className="wrap demo-split">
          <UsageBoard locale={l} m={m} />
          <div className="demo-text">
            <h2 id="demo" className="section-title display">{m.landing.demo.title}</h2>
            <ol className="demo-steps" data-demo-steps>
              {m.landing.demo.steps.map((s, i) => (
                <li key={i}>
                  <span className="step-n" dir="ltr" aria-hidden="true">{String(i + 1).padStart(2, "0")}</span>
                  <div>
                    <p>{s.body}</p>
                    {s.command && <code dir="ltr">{s.command}</code>}
                  </div>
                </li>
              ))}
            </ol>
          </div>
        </div>
      </section>

      <section className="section" aria-labelledby="flow">
        <div className="wrap">
          <h2 id="flow" className="section-title display">{m.landing.flowTitle}</h2>
          <div className="flow" data-flow>
            <div className="flow-visual" aria-hidden="true">
              {m.landing.flow.map((f, i) => (
                <div className={`fv${i === 0 ? " on" : ""}`} key={f.command} dir="ltr">
                  <code>{f.command}</code>
                  {"also" in f && f.also ? <code>{f.also}</code> : null}
                </div>
              ))}
              <div className="rail">
                {m.landing.flow.map((f) => (
                  <span key={f.command} />
                ))}
              </div>
            </div>
            <ol className="steps">
              {m.landing.flow.map((f, i) => (
                <li className={`step${i === 0 ? " on" : ""}`} key={f.command}>
                  <span className="step-n" dir="ltr" aria-hidden="true">{String(i + 1).padStart(2, "0")}</span>
                  <h3 className="step-cmd" dir="ltr">
                    <a href={pagePath(locale, "features", flowFeatures[i])}>
                      <code>{f.command}</code>
                      {"also" in f && f.also ? <> <code>{f.also}</code></> : null}
                    </a>
                  </h3>
                  <p>{f.body}</p>
                </li>
              ))}
            </ol>
          </div>
        </div>
      </section>

      {guides.length > 0 && (
        <section className="section" aria-labelledby="guides">
          <div className="wrap">
            <h2 id="guides" className="section-title display">{m.landing.guidesTitle}</h2>
            <ul className="cards">
              {guides.map((g) => (
                <li key={g.slug}>
                  <a className="card reveal" href={pagePath(locale, "guides", g.slug)}>
                    <strong>{g.title}</strong>
                    <span>{g.description}</span>
                    <span className="go">
                      <ArrowIcon />
                    </span>
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </section>
      )}

      <section className="section" aria-labelledby="facts">
        <div className="wrap">
          <h2 id="facts" className="section-title display">{m.landing.factsTitle}</h2>
          <div className="facts">
            <p className="big">{m.facts.what}</p>
            <div className="small">
              <p>{m.facts.how}</p>
              <p>{m.facts.install}</p>
              <p>{m.facts.official}</p>
            </div>
          </div>
        </div>
      </section>
      <JsonLd data={softwareLd(m.site.description, l.lang)} />
    </Shell>
  );
}

export function landingTitle(locale: string): { title: string; description: string } {
  const m = getMessages(locale);
  return { title: m.site.title, description: m.site.description };
}
