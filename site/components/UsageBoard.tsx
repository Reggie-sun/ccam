import type { Messages } from "../lib/i18n";
import { format } from "../lib/i18n";
import type { Locale } from "../i18n/locales";

// Sample accounts, as `ccs usage` would print them.
const rows = [
  { slot: 1, name: "personal", five: 72, week: 41, resets: 100 },
  { slot: 2, name: "work", five: 8, week: 19, resets: 252 },
  { slot: 3, name: "side", five: 100, week: 68, resets: 22 },
];

function duration(locale: Locale, minutes: number): string {
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  const unit = (value: number, u: "hour" | "minute") =>
    new Intl.NumberFormat(locale.lang, { style: "unit", unit: u, unitDisplay: "short" }).format(value);
  return h ? `${unit(h, "hour")} ${unit(m, "minute")}` : unit(m, "minute");
}

// SVG rather than a style attribute: the CSP allows no inline styles.
function Bar({ value, label }: { value: number; label: string }) {
  return (
    <svg className={`bar${value >= 100 ? " full" : ""}`} viewBox="0 0 100 10" preserveAspectRatio="none" role="img" aria-label={label}>
      <rect className="fill" width={Math.min(value, 100)} height="10" />
    </svg>
  );
}

export function UsageBoard({ locale, m }: { locale: Locale; m: Messages }) {
  const pct = (v: number) => new Intl.NumberFormat(locale.lang, { style: "percent" }).format(v / 100);
  return (
    <figure className="board">
      <figcaption>{m.landing.boardLabel}</figcaption>
      {rows.map((r) => (
        <div className="board-row" key={r.name}>
          <span className="slot" dir="ltr">{r.slot}</span>
          <span className="name" dir="ltr">{r.name}</span>
          <div className="meter">
            <span>{m.landing.board.fiveHour}</span>
            <Bar value={r.five} label={`${m.landing.board.fiveHour} ${pct(r.five)}`} />
            <span className="pct">{pct(r.five)}</span>
          </div>
          <div className="meter">
            <span>{m.landing.board.weekly}</span>
            <Bar value={r.week} label={`${m.landing.board.weekly} ${pct(r.week)}`} />
            <span className="pct">{pct(r.week)}</span>
          </div>
          <span className="resets">{format(locale, m.landing.board.resets, { time: duration(locale, r.resets) })}</span>
        </div>
      ))}
      <footer>{format(locale, m.landing.board.mostRoom, { name: "work" })}</footer>
    </figure>
  );
}
