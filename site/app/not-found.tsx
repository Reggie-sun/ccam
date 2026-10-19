import type { Metadata } from "next";
import { Shell } from "../components/Shell";
import { getMessages } from "../lib/i18n";
import { pagePath } from "../lib/paths";

export const metadata: Metadata = { title: "Page not found | ccam", robots: { index: false } };

export default function NotFound() {
  const m = getMessages("en");
  return (
    <Shell locale="en" alternates={{}}>
      <div className="wrap page">
        <article className="article">
          <h1>{m.notFound.title}</h1>
          <p>
            {m.notFound.body} <a href={pagePath("en", "docs", "quick-start")}>{m.landing.readDocs}</a>
          </p>
        </article>
      </div>
    </Shell>
  );
}
