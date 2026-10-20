import { getLocale } from "../i18n/locales";
import { faqEntries, getPage, navGroups, pageLocales, renderMarkdown, type Page, type Section } from "../lib/content";
import { format, formatDate, getMessages, type Messages } from "../lib/i18n";
import { abs, homePath, pagePath, REPO } from "../lib/paths";
import { articleLd, AUTHOR, CLAUDE_CODE_VERSION, faqLd, VERSION } from "../lib/jsonld";
import { placeDiagrams } from "../lib/diagrams";
import { ogPath } from "../lib/seo";
import commands from "../data/commands.json";
import { ChevronIcon } from "./Icons";
import { JsonLd } from "./JsonLd";
import { Shell } from "./Shell";

export function docAlternates(section: Section, slug: string): Record<string, string> {
  return Object.fromEntries(pageLocales(section, slug).map((l) => [l.code, pagePath(l.code, section, slug)]));
}

/** The sidebar groups (lib/nav.ts) with their names in this locale. */
function docGroups(locale: string, m: Messages): { title: string; pages: Page[] }[] {
  return navGroups(locale).map((g) => ({ title: m.docs.groups[g.key], pages: g.pages }));
}

function DocList({ page, groups }: { page: Page; groups: ReturnType<typeof docGroups> }) {
  return groups.map((g) => (
    <div key={g.title}>
      <p className="group-title">{g.title}</p>
      <ul>
        {g.pages.map((p) => (
          <li key={`${p.section}/${p.slug}`}>
            <a href={pagePath(p.locale, p.section, p.slug)} aria-current={p.slug === page.slug && p.section === page.section ? "page" : undefined}>{p.title}</a>
          </li>
        ))}
      </ul>
    </div>
  ));
}

// The feature page that documents each command in full.
const commandFeature: Record<string, string> = {
  ccam: "launch", name: "launch", use: "launch", add: "accounts", list: "accounts", login: "accounts", logout: "accounts", rename: "accounts",
  remove: "accounts", import: "accounts", completion: "accounts", doctor: "doctor", mcp: "mcp", migrate: "migrate", next: "next", share: "share",
  shims: "shortcuts", shortcut: "shortcuts", statusline: "usage", usage: "usage", update: "update",
};

function commandsTable(locale: string, m: Messages): string {
  const esc = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  const rows = (commands as { key: string; usage: string }[])
    .map((c) => {
      const slug = commandFeature[c.key];
      const code = `<code dir="ltr">${esc(c.usage)}</code>`;
      const cell = slug ? `<a href="${pagePath(getPage(locale, "features", slug) ? locale : "en", "features", slug)}">${code}</a>` : code;
      return `<tr><td>${cell}</td><td>${esc((m.commands as Record<string, string>)[c.key] ?? "")}</td></tr>`;
    })
    .join("");
  return `<table tabindex="0"><thead><tr><th scope="col">${esc(m.docs.commandsTable.command)}</th><th scope="col">${esc(m.docs.commandsTable.does)}</th></tr></thead><tbody>${rows}</tbody></table>`;
}

export function DocPage({ page }: { page: Page }) {
  const l = getLocale(page.locale)!;
  const m = getMessages(page.locale);
  let { html, headings } = renderMarkdown(page, m.copy);
  html = placeDiagrams(html.replace("<!-- commands -->", commandsTable(page.locale, m)), page.file);
  const groups = docGroups(page.locale, m);
  // Previous and next follow the sidebar across groups, as one reading path.
  const ordered = groups.flatMap((g) => g.pages);
  const i = ordered.findIndex((p) => p.section === page.section && p.slug === page.slug);
  const prev = ordered[i - 1];
  const next = ordered[i + 1];
  const path = pagePath(page.locale, page.section, page.slug);
  const active = page.slug === "faq" ? "faq" : page.section === "docs" || page.section === "guides" || page.section === "features" ? "docs" : undefined;
  // The crumb names the page's own section (guides are listed with docs in the sidebar).
  const sectionName = (m.docs.sections as Record<Section, string>)[page.section];
  const sectionPath = pagePath(page.locale, page.section);
  // The author's name is a link, wherever the language puts it in the sentence.
  const byline = format(l, m.docs.byline, { author: "\u0000" }).split("\u0000");
  const published = getPage("en", page.section, page.slug)?.published ?? page.updated;
  const toc = headings.filter((h) => h.depth === 2);
  const tocList = (
    <ol>
      {toc.map((h) => (
        <li key={h.id}>
          <a href={`#${h.id}`}>{h.text}</a>
        </li>
      ))}
    </ol>
  );
  return (
    <Shell
      locale={page.locale}
      alternates={docAlternates(page.section, page.slug)}
      active={active}
      menu={<DocList page={page} groups={groups} />}
    >
      <div className="progress" aria-hidden="true" />
      <div className={`wrap docs${toc.length > 1 ? " has-toc" : ""}`}>
        <aside className="sidebar">
          <nav aria-label={m.nav.docs}>
            <DocList page={page} groups={groups} />
          </nav>
        </aside>
        <article className="article">
          <p className="crumb">
            <a href={sectionPath}>{sectionName}</a>
          </p>
          <h1>{page.title}</h1>
          <p className="meta">
            {byline[0]}
            <a href={AUTHOR.url} dir="ltr">
              {AUTHOR.name}
            </a>
            {byline[1]}
            <br />
            {format(l, m.docs.updated, { date: formatDate(l, page.updated) })}
            <br />
            {format(l, m.docs.versions, { ccam: page.ccam ?? VERSION, claude: page.claude ?? CLAUDE_CODE_VERSION })}
          </p>
          {page.translation === "machine" && (
            <p className="notice">
              {m.docs.machine} <a href={`${REPO}/edit/main/${page.file}`}>{m.docs.edit}</a>
            </p>
          )}
          {toc.length > 1 && (
            <details className="toc toc-inline">
              <summary>
                {m.docs.onThisPage}
                <ChevronIcon />
              </summary>
              <div className="toc-body">{tocList}</div>
            </details>
          )}
          <div className="prose" dangerouslySetInnerHTML={{ __html: html }} />
          <nav className="pager" aria-label={`${m.docs.prev}, ${m.docs.next}`}>
            {prev && (
              <a href={pagePath(page.locale, prev.section, prev.slug)}>
                <small>{m.docs.prev}</small>
                {prev.title}
              </a>
            )}
            {next && (
              <a className="next" href={pagePath(page.locale, next.section, next.slug)}>
                <small>{m.docs.next}</small>
                {next.title}
              </a>
            )}
          </nav>
          <p className="related">
            <a href={`${REPO}/edit/main/${page.file}`}>{m.docs.edit}</a>
          </p>
        </article>
        {toc.length > 1 && (
          <div className="toc toc-rail">
            <p className="toc-title">{m.docs.onThisPage}</p>
            <nav aria-label={m.docs.onThisPage}>{tocList}</nav>
          </div>
        )}
      </div>
      <JsonLd
        data={articleLd({
          type: page.section === "blog" ? "BlogPosting" : "TechArticle",
          title: page.title,
          description: page.description,
          url: abs(path),
          published,
          updated: page.updated,
          lang: l.lang,
          image: abs(ogPath(page.locale, `${page.section}-${page.slug}`)),
          crumbs: [
            { name: "ccam", url: abs(homePath(page.locale)) },
            { name: sectionName, url: abs(sectionPath) },
            { name: page.title, url: abs(path) },
          ],
        })}
      />
      {page.slug === "faq" && <JsonLd data={faqLd(faqEntries(page), l.lang)} />}
    </Shell>
  );
}
