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

// The story: Claude Code on "side" hits its 5-hour limit, ccam usage shows the room left,
// and ccam next continues the same conversation on "work". Numbers match the docs' sample output.
const ASK = "refactor the parser and add tests";

/** Rendered in its final state; site.js replays it line by line (data-seq: pause after, ms) when it scrolls into view. */
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
        <p className="line" dir="ltr" data-seq="250" data-step="0">
          <span className="ps">$</span>
          <span data-type>claude3</span>
        </p>
        <div className="cc" dir="ltr" lang="en">
          <p className="cc-head" data-seq="250">
            <span className="cc-mark" aria-hidden="true">✻</span> Claude Code <span className="cc-acct">side</span>
          </p>
          <p className="cc-prompt" data-seq="700">
            <span className="ps">&gt;</span> <span data-type>{ASK}</span>
          </p>
          <p className="cc-limit" data-seq="1000" data-step="1">5-hour limit reached · resets in 22m</p>
        </div>
        <p className="line" dir="ltr" data-seq="250" data-step="2">
          <span className="ps">$</span>
          <span data-type>ccam usage</span>
        </p>
        <div className="rows" data-seq="1900">
          {rows.map((r) => (
            <div className={`row${r.name === "work" ? " on" : ""}`} key={r.name} data-pick={r.name === "work" ? "" : undefined}>
              <span className="slot" dir="ltr">{r.slot}</span>
              <span className="name" dir="ltr">{r.name}</span>
              {meter(m.landing.board.fiveHour, r.five)}
              {meter(m.landing.board.weekly, r.week)}
              <span className="resets">{format(locale, m.landing.board.resets, { time: duration(locale, r.resets) })}</span>
            </div>
          ))}
        </div>
        <p className="line note" data-seq="700">{format(locale, m.landing.board.mostRoom, { name: "work" })}</p>
        <p className="line" dir="ltr" data-seq="400" data-step="3">
          <span className="ps">$</span>
          <span data-type>ccam next --auto --continue</span>
        </p>
        <div className="cc" dir="ltr" lang="en">
          <p className="cc-head" data-seq="350">
            <span className="cc-mark" aria-hidden="true">✻</span> Claude Code <span className="cc-acct">work</span>
          </p>
          <p className="cc-prompt" data-seq="500">
            <span className="ps">&gt;</span> {ASK}
          </p>
          <p className="cc-status" data-seq="0">
            [work] · Opus · 5h 8% · 7d 19% <span className="caret" aria-hidden="true" />
          </p>
        </div>
      </div>
      <figcaption>{m.landing.boardLabel}</figcaption>
    </figure>
  );
}
