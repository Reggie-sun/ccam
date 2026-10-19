// The document around every page: <html lang dir>, head extras, header with
// the JS-free language picker, the (hidden) language banner and the footer.
import { getLocale, locales, nativeName } from "../i18n/locales";
import { format, getMessages, hasMessages } from "../lib/i18n";
import { homePath, pagePath, REPO } from "../lib/paths";
import { websiteLd } from "../lib/jsonld";
import { BrandMark, GlobeIcon, MenuIcon, ThemeIcon } from "./Icons";
import { JsonLd } from "./JsonLd";

// Applies the saved theme before first paint; hashed into the CSP at build.
export const THEME_INIT = `try{var t=localStorage.getItem("ccam-theme");if(t)document.documentElement.dataset.theme=t}catch(e){}`;

const SPECULATION = JSON.stringify({ prerender: [{ where: { href_matches: "/*" }, eagerness: "moderate" }] });

type Props = {
  locale: string;
  /** This page in every locale it exists in (code → path). */
  alternates: Record<string, string>;
  active?: "docs" | "faq" | "changelog";
  /** Extra links for the phone menu: the doc list on doc pages. */
  menu?: React.ReactNode;
  children: React.ReactNode;
};

export function Shell({ locale, alternates, active, menu, children }: Props) {
  const l = getLocale(locale)!;
  const m = getMessages(locale);
  const built = locales.filter((x) => hasMessages(x.code));
  return (
    <html lang={l.lang} dir={l.dir} data-font={l.font}>
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME_INIT }} />
        {l.font === "brand" && <link rel="preload" href="/fonts/newsreader.woff2" as="font" type="font/woff2" crossOrigin="" />}
        <link rel="icon" href="/favicon.svg" type="image/svg+xml" />
        <meta name="theme-color" content="#f5f5f5" media="(prefers-color-scheme: light)" />
        <meta name="theme-color" content="#0c0a09" media="(prefers-color-scheme: dark)" />
        <script type="speculationrules" dangerouslySetInnerHTML={{ __html: SPECULATION }} />
      </head>
      <body>
        <a className="skip" href="#main">{m.nav.skip}</a>
        <div className="lang-banner" id="lang-banner" hidden>
          <div className="wrap">
            <span data-banner-text />
            <a data-banner-link href="/" />
            <button type="button" data-banner-dismiss>{m.banner.dismiss}</button>
          </div>
        </div>
        <header className="site-header">
          <div className="wrap">
            <a className="brand" href={homePath(locale)}>
              <BrandMark />
              ccam
            </a>
            <nav className="site-nav" aria-label={m.nav.menu}>
              <a href={pagePath(locale, "docs", "quick-start")} aria-current={active === "docs" ? "page" : undefined}>{m.nav.docs}</a>
              <a href={pagePath(locale, "docs", "faq")} aria-current={active === "faq" ? "page" : undefined}>{m.nav.faq}</a>
              <a href={REPO}>{m.nav.github}</a>
            </nav>
            <div className="actions">
              <details className="picker">
                <summary aria-label={m.nav.language}>
                  <GlobeIcon />
                  {nativeName(l)}
                </summary>
                <ul id="lang-list" className="popover">
                  {built.map((x) => {
                    const href = alternates[x.code] ?? homePath(x.code);
                    const xm = getMessages(x.code);
                    const name = nativeName(x);
                    return (
                      <li key={x.code}>
                        <a
                          href={href}
                          hrefLang={x.lang}
                          lang={x.lang}
                          data-dir={x.dir}
                          aria-current={x.code === locale ? "true" : undefined}
                          data-same={alternates[x.code] ? "" : undefined}
                          data-suggest={format(x, xm.banner.suggest, { language: name })}
                          data-open={format(x, xm.banner.open, { language: name })}
                        >
                          {name}
                        </a>
                      </li>
                    );
                  })}
                </ul>
              </details>
              <button type="button" className="icon-button theme-toggle" data-theme-toggle aria-label={m.nav.theme} hidden>
                <ThemeIcon />
              </button>
              <details className="menu">
                <summary className="icon-button" aria-label={m.nav.menu}>
                  <MenuIcon />
                </summary>
                <div className="popover menu-panel">
                  <a href={pagePath(locale, "docs", "quick-start")} aria-current={active === "docs" ? "page" : undefined}>{m.nav.docs}</a>
                  <a href={pagePath(locale, "docs", "faq")} aria-current={active === "faq" ? "page" : undefined}>{m.nav.faq}</a>
                  <a href={REPO}>{m.nav.github}</a>
                  {menu && <div className="menu-docs">{menu}</div>}
                </div>
              </details>
            </div>
          </div>
        </header>
        <main id="main">{children}</main>
        <footer className="site-footer">
          <div className="wrap">
            <a className="brand" href={homePath(locale)}>ccam</a>
            <p>{m.footer.unofficial}</p>
            <p>{m.footer.terms}</p>
            <ul className="footer-links">
              <li><a href={REPO}>{m.footer.license}</a></li>
              <li><a href={`/${locale}/languages/`}>{m.footer.languages}</a></li>
            </ul>
          </div>
        </footer>
        <JsonLd data={websiteLd(l.lang)} />
      </body>
    </html>
  );
}
