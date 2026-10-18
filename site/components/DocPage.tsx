import { getLocale } from "../i18n/locales";
import { faqEntries, getPages, pageLocales, renderMarkdown, sections, type Page, type Section } from "../lib/content";
import { format, formatDate, getMessages, type Messages } from "../lib/i18n";
import { abs, homePath, pagePath, REPO } from "../lib/paths";
import { articleLd, CLAUDE_CODE_VERSION, faqLd, VERSION } from "../lib/jsonld";
import commands from "../data/commands.json";
import { JsonLd } from "./JsonLd";
import { Shell } from "./Shell";

export function docAlternates(section: Section, slug: string): Record<string, string> {
  return Object.fromEntries(pageLocales(section, slug).map((l) => [l.code, pagePath(l.code, section, slug)]));
}

function commandsTable(m: Messages): string {
  const esc = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  const rows = (commands as { key: string; usage: string }[])
    .map((c) => `<tr><td><code dir="ltr">${esc(c.usage)}</code></td><td>${esc((m.commands as Record<string, string>)[c.key] ?? "")}</td></tr>`)
    .join("");
  return `<table tabindex="0"><thead><tr><th scope="col">${esc(m.docs.commandsTable.command)}</th><th scope="col">${esc(m.docs.commandsTable.does)}</th></tr></thead><tbody>${rows}</tbody></table>`;
}

export function DocPage({ page }: { page: Page }) {
  const l = getLocale(page.locale)!;
  const m = getMessages(page.locale);
  let { html } = renderMarkdown(page, m.copy);
  html = html.replace("<!-- commands -->", commandsTable(m));
  const siblings = getPages(page.locale, page.section);
  const i = siblings.findIndex((p) => p.slug === page.slug);
  const prev = siblings[i - 1];
  const next = siblings[i + 1];
  const path = pagePath(page.locale, page.section, page.slug);
  const active = page.slug === "faq" ? "faq" : page.section === "docs" ? "docs" : page.section === "guides" ? "guides" : undefined;
  const sectionName = (m.docs.sections as Record<Section, string>)[page.section];
  return (
    <Shell locale={page.locale} alternates={docAlternates(page.section, page.slug)} active={active}>
      <div className="wrap docs">
        <aside className="sidebar">
          <details open>
            <summary>{m.nav.menu}</summary>
            <nav aria-label={m.nav.docs}>
              {sections.map((s) => {
                const list = getPages(page.locale, s);
                if (!list.length) return null;
                return (
                  <div key={s}>
                    <h2>{(m.docs.sections as Record<Section, string>)[s]}</h2>
                    <ul>
                      {list.map((p) => (
                        <li key={p.slug}>
                          <a href={pagePath(p.locale, s, p.slug)} aria-current={p.slug === page.slug && s === page.section ? "page" : undefined}>{p.title}</a>
                        </li>
                      ))}
                    </ul>
                  </div>
                );
              })}
            </nav>
          </details>
        </aside>
        <article className="article">
          <h1>{page.title}</h1>
          <p className="meta">
            {format(l, m.docs.updated, { date: formatDate(l, page.updated) })}
            <br />
            {format(l, m.docs.versions, { ccam: page.ccam ?? VERSION, claude: page.claude ?? CLAUDE_CODE_VERSION })}
          </p>
          {page.translation === "machine" && (
            <p className="notice">
              {m.docs.machine} <a href={`${REPO}/edit/main/${page.file}`}>{m.docs.edit}</a>
            </p>
          )}
          <div className="prose" dangerouslySetInnerHTML={{ __html: html }} />
          <nav className="pager" aria-label={`${m.docs.prev}, ${m.docs.next}`}>
            {prev && (
              <a href={pagePath(page.locale, page.section, prev.slug)}>
                <small>{m.docs.prev}</small>
                {prev.title}
              </a>
            )}
            {next && (
              <a className="next" href={pagePath(page.locale, page.section, next.slug)}>
                <small>{m.docs.next}</small>
                {next.title}
              </a>
            )}
          </nav>
          <p className="related">
            <a href={`${REPO}/edit/main/${page.file}`}>{m.docs.edit}</a>
          </p>
        </article>
      </div>
      <JsonLd
        data={articleLd({
          title: page.title,
          description: page.description,
          url: abs(path),
          updated: page.updated,
          lang: l.lang,
          crumbs: [
            { name: "ccam", url: abs(homePath(page.locale)) },
            { name: sectionName, url: abs(pagePath(page.locale, page.section, siblings[0].slug)) },
            { name: page.title, url: abs(path) },
          ],
        })}
      />
      {page.slug === "faq" && <JsonLd data={faqLd(faqEntries(page), l.lang)} />}
    </Shell>
  );
}
