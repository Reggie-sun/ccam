import { parse } from "yaml";

/** Splits "---\nyaml\n---\nbody". Shared by the app and the build scripts. */
export function splitFrontMatter<T = Record<string, unknown>>(source: string): { data: T; body: string } {
  const m = /^---\n([\s\S]*?)\n---\n?/.exec(source);
  if (!m) throw new Error("missing front matter");
  return { data: parse(m[1]) as T, body: source.slice(m[0].length) };
}
