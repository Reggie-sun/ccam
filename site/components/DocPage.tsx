import { getLocale } from "../i18n/locales";
import { faqEntries, getPage, getPages, pageLocales, renderMarkdown, type Page, type Section } from "../lib/content";
import { format, formatDate, getMessages, type Messages } from "../lib/i18n";
import { abs, homePath, pagePath, REPO } from "../lib/paths";
import { articleLd, AUTHOR, CLAUDE_CODE_VERSION, faqLd, VERSION } from "../lib/jsonld";
import { ogPath } from "../lib/seo";
import commands from "../data/commands.json";
import { ChevronIcon } from "./Icons";
import { JsonLd } from "./JsonLd";
import { Shell } from "./Shell";

export function docAlternates(section: Section, slug: string): Record<string, string> {
  return Object.fromEntries(pageLocales(section, slug).map((l) => [l.code, pagePath(l.code, section, slug)]));
}

/** Docs and guides read as one list; comparisons and the blog keep their own groups. */
function docGroups(locale: string, m: Messages): { title: string; pages: Page[] }[] {
  const t = m.docs.sections as Record<Section, string>;
  return [
    { title: t.docs, pages: [...getPages(locale, "docs"), ...getPages(locale, "guides")] },
    { title: t.compare, pages: getPages(locale, "compare") },
    { title: t.blog, pages: getPages(locale, "blog") },
  ].filter((g) => g.pages.length);
}

function DocList({ page, groups, Title }: { page: Page; groups: ReturnType<typeof docGroups>; Title: "h2" | "p" }) {
  return groups.map((g) => (
    <div key={g.title}>
      <Title className="group-title">{g.title}</Title>
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
  let { html, headings } = renderMarkdown(page, m.copy);
  html = html.replace("<!-- commands -->", commandsTable(m));
  const groups = docGroups(page.locale, m);
  const group = groups.find((g) => g.pages.some((p) => p.section === page.section && p.slug === page.slug))!;
  const siblings = group.pages;
  const i = siblings.findIndex((p) => p.section === page.section && p.slug === page.slug);
  const prev = siblings[i - 1];
  const next = siblings[i + 1];
  const path = pagePath(page.locale, page.section, page.slug);
  const active = page.slug === "faq" ? "faq" : page.section === "docs" || page.section === "guides" ? "docs" : undefined;
  const sectionName = group.title;
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
      menu={<DocList page={page} groups={groups} Title="p" />}
    >
      <div className="progress" aria-hidden="true" />
      <div className={`wrap docs${toc.length > 1 ? " has-toc" : ""}`}>
        <aside className="sidebar">
          <nav aria-label={m.nav.docs}>
            <DocList page={page} groups={groups} Title="h2" />
          </nav>
        </aside>
        <article className="article">
          <p className="crumb">{sectionName}</p>
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
            // A section's first page stands for the section; it is not listed twice.
            ...(i === 0 ? [] : [{ name: sectionName, url: abs(pagePath(page.locale, siblings[0].section, siblings[0].slug)) }]),
            { name: page.title, url: abs(path) },
          ],
        })}
      />
      {page.slug === "faq" && <JsonLd data={faqLd(faqEntries(page), l.lang)} />}
    </Shell>
  );
}
