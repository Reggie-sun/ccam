import type { Messages } from "../lib/i18n";
import { format } from "../lib/i18n";
import type { Locale } from "../i18n/locales";

// Sample accounts, as `ccam usage` would print them.
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

/** Rendered in its final state; site.js replays it as a terminal session when it scrolls into view. */
export function UsageBoard({ locale, m }: { locale: Locale; m: Messages }) {
  const pct = (v: number) => new Intl.NumberFormat(locale.lang, { style: "percent" }).format(v / 100);
  const meter = (label: string, value: number) => (
    <div className="meter">
      <span>{label}</span>
      <Bar value={value} label={`${label} ${pct(value)}`} />
      <span className="pct" data-value={value}>{pct(value)}</span>
    </div>
  );
  return (
    <figure className="stage reveal" data-demo>
      <div className="term">
        <p className="line" dir="ltr">
          <span className="ps">$</span>
          <span data-type>ccam usage</span>
        </p>
        <div className="rows">
          {rows.map((r) => (
            <div className={`row later${r.name === "work" ? " on" : ""}`} key={r.name} data-pick={r.name === "work" ? "" : undefined}>
              <span className="slot" dir="ltr">{r.slot}</span>
              <span className="name" dir="ltr">{r.name}</span>
              {meter(m.landing.board.fiveHour, r.five)}
              {meter(m.landing.board.weekly, r.week)}
              <span className="resets">{format(locale, m.landing.board.resets, { time: duration(locale, r.resets) })}</span>
            </div>
          ))}
        </div>
        <p className="line note later">{format(locale, m.landing.board.mostRoom, { name: "work" })}</p>
        <p className="line later" dir="ltr">
          <span className="ps">$</span>
          <span data-type>ccam next --auto</span>
          <span className="caret" aria-hidden="true" />
        </p>
      </div>
      <figcaption>{m.landing.boardLabel}</figcaption>
    </figure>
  );
}
