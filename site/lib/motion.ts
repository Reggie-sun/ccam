// Motion diagrams for the feature pages: inline SVG drawn from a few pieces.
// Their CSS ("Motion diagrams" in app/globals.css) moves only transform and
// clip-path, never the opacity of text, and stops under
// prefers-reduced-motion, which leaves the whole diagram standing still.
//
// Every piece of a diagram shares one cycle (--t on the figure, seconds) and
// starts --d seconds into it, so a flow plays in order and repeats. Values
// travel in style attributes, which the CSP hashes at build.
//
// Labels are commands, paths and names, which are never translated; the page
// text around a diagram explains it, so one diagram serves every language.
// Keep diagrams at most 400 units wide and flowing downwards: a phone shows
// them at about full size.

const esc = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
const css = (vars: Record<string, string | number | undefined>) =>
  Object.entries(vars)
    .filter(([, v]) => v !== undefined)
    .map(([k, v]) => `--${k}:${v}`)
    .join(";");

export type Kind = "proc" | "file" | "dir" | "net" | "group";

/**
 * A labelled box: a process, a file, a folder, a network endpoint or a group
 * of them. sub is a smaller line under the label, or several lines.
 */
export function box(x: number, y: number, w: number, h: number, label: string, o: { kind?: Kind; sub?: string | string[] } = {}): string {
  const kind = o.kind ?? "proc";
  const subs = o.sub ? [o.sub].flat() : [];
  const cy = subs.length ? y + h / 2 - 5 - (subs.length - 1) * 6.5 : y + h / 2 + 4.5;
  return (
    `<g class="mo-box mo-${kind}"><rect x="${x}" y="${y}" width="${w}" height="${h}" rx="${kind === "group" ? 14 : 8}"/>` +
    `<text x="${x + w / 2}" y="${cy}" text-anchor="middle">${esc(label)}</text>` +
    subs.map((s, i) => `<text class="mo-sub" x="${x + w / 2}" y="${cy + 15 + i * 13}" text-anchor="middle">${esc(s)}</text>`).join("") +
    `</g>`
  );
}

/** A plain label, for example beside a line. */
export function label(x: number, y: number, text: string, o: { anchor?: "start" | "middle" | "end"; muted?: boolean } = {}): string {
  return `<text class="mo-label${o.muted ? " mo-sub" : ""}" x="${x}" y="${y}" text-anchor="${o.anchor ?? "start"}">${esc(text)}</text>`;
}

/** A connector through the given points, with an arrowhead at the end unless arrow is false. */
export function line(pts: [number, number][], o: { arrow?: boolean; dashed?: boolean } = {}): string {
  const d = pts.map(([x, y], i) => `${i ? "L" : "M"}${x} ${y}`).join(" ");
  let head = "";
  if (o.arrow !== false && pts.length > 1) {
    const [x1, y1] = pts[pts.length - 2];
    const [x2, y2] = pts[pts.length - 1];
    const a = Math.atan2(y2 - y1, x2 - x1);
    const p = (da: number) => `${(x2 - 7 * Math.cos(a + da)).toFixed(1)} ${(y2 - 7 * Math.sin(a + da)).toFixed(1)}`;
    head = `<path class="mo-head" d="M${p(0.5)} L${x2} ${y2} L${p(-0.5)}"/>`;
  }
  return `<path class="mo-line${o.dashed ? " mo-dashed" : ""}" d="${d}"/>${head}`;
}

/**
 * A dot travelling in a straight line from one point to another, starting
 * at second d of the cycle and taking a quarter of it. Chain several with
 * later d values to follow a bent path.
 */
export function packet(from: [number, number], to: [number, number], d: number, o: { label?: string } = {}): string {
  const style = css({ x0: `${from[0]}px`, y0: `${from[1]}px`, x1: `${to[0]}px`, y1: `${to[1]}px`, d: `${d}s` });
  const tag = o.label ? `<text class="mo-tag" x="0" y="-9" text-anchor="middle">${esc(o.label)}</text>` : "";
  return `<g class="mo-pkt" style="${style}"><circle r="5"/>${tag}</g>`;
}

/** A bar under or beside a box that lights up at second d, for a quarter of the cycle. */
export function glow(x: number, y: number, w: number, h: number, d: number): string {
  return `<rect class="mo-glow" x="${x}" y="${y}" width="${w}" height="${h}" rx="8" style="${css({ d: `${d}s` })}"/>`;
}

/** Content that is wiped in from the start side at second d and stays until the cycle ends. */
export function reveal(d: number, content: string): string {
  return `<g class="mo-reveal" style="${css({ d: `${d}s` })}">${content}</g>`;
}

/**
 * The figure: an SVG of the given size whose pieces share a cycle of t
 * seconds. name only identifies it in the HTML (data-motion), for tests.
 */
export function figure(name: string, w: number, h: number, t: number, parts: string[]): string {
  if (w > 400) throw new Error(`motion ${name}: ${w} units is too wide for a phone (max 400)`);
  return (
    `<figure class="mo" dir="ltr" data-motion="${esc(name)}" style="${css({ t: `${t}s` })}">` +
    `<svg viewBox="0 0 ${w} ${h}" width="${w}" height="${h}" aria-hidden="true" focusable="false">${parts.join("")}</svg></figure>`
  );
}
