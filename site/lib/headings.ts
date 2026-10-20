// Heading anchors, shared by the pages (lib/content.ts) and the search index
// (scripts/search.ts) so a search result's #fragment always exists.
import { Marked, type Tokens } from "marked";

export function slugify(text: string): string {
  return text.toLowerCase().normalize("NFKC").replace(/<[^>]+>/g, "").replace(/[^\p{L}\p{N}]+/gu, "-").replace(/^-|-$/g, "");
}

/** Anchor of each top-level heading, in order. */
export function headingIds(markdown: string): string[] {
  return new Marked().lexer(markdown).filter((t) => t.type === "heading").map((t) => slugify((t as Tokens.Heading).text.replace(/`/g, "")));
}
