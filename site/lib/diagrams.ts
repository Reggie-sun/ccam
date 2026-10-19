// Motion diagrams by name: lib/diagrams/<name>.ts, each exporting a function
// that returns the figure's HTML. A page places one with <!-- motion:<name> -->.
const files = import.meta.glob<{ default: () => string }>("./diagrams/*.ts", { eager: true });

export const diagrams: Record<string, () => string> = Object.fromEntries(
  Object.entries(files).map(([path, mod]) => [/([^/]+)\.ts$/.exec(path)![1], mod.default]),
);

export const MOTION = /<!-- motion:([a-z0-9-]+) -->/g;

/** Replaces every motion marker in html with its diagram; an unknown name fails the build. */
export function placeDiagrams(html: string, file: string): string {
  return html.replace(MOTION, (_, name: string) => {
    const draw = diagrams[name];
    if (!draw) throw new Error(`${file}: no diagram "${name}" in lib/diagrams/`);
    return draw();
  });
}
