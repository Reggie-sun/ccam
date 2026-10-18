// How much of the world population speaks at least one supported language,
// from Unicode CLDR (cldr-core territoryInfo). Two estimates
// (docs/WEBSITE_PLAN.md 5.2):
//  low:  per country, only the most-spoken supported language counts
//  high: per country, all supported languages add up, capped at 100%
// Writes data/coverage.json for the built locales and prints the plan's
// phase 1 and phase 1+2 numbers. `--strict` fails below 80% (low, all locales).
import { readFileSync, writeFileSync } from "node:fs";
import { locales, type Locale } from "../i18n/locales.ts";
import { builtLocales } from "./lib.ts";

type Territory = { _population: string; languagePopulation?: Record<string, { _populationPercent: string }> };
const info: Record<string, Territory> = JSON.parse(
  readFileSync("node_modules/cldr-core/supplemental/territoryInfo.json", "utf8"),
).supplemental.territoryInfo;

/** CLDR language codes a locale covers. Chinese is split by script. */
function cldrCodes(l: Locale): string[] {
  if (l.code === "zh-hans") return ["zh", "zh_Hans"];
  if (l.code === "zh-hant") return ["zh_Hant"];
  if (l.code === "pt-br") return ["pt"];
  return [l.code];
}

export function coverage(list: Locale[]): { low: number; high: number } {
  const codes = new Set(list.flatMap(cldrCodes));
  let world = 0, low = 0, high = 0;
  for (const t of Object.values(info)) {
    const pop = Number(t._population);
    if (!pop) continue;
    world += pop;
    const shares = Object.entries(t.languagePopulation ?? {})
      .filter(([lang]) => codes.has(lang))
      .map(([, v]) => Number(v._populationPercent) / 100);
    if (!shares.length) continue;
    low += pop * Math.max(...shares);
    high += pop * Math.min(1, shares.reduce((a, b) => a + b, 0));
  }
  return { low: low / world, high: high / world };
}

const pct = (v: number) => `${(v * 100).toFixed(1)}%`;
const built = coverage(builtLocales());
const phase1 = coverage(locales.filter((l) => l.phase === 1));
const all = coverage(locales);
writeFileSync("data/coverage.json", JSON.stringify({ low: +built.low.toFixed(4), high: +built.high.toFixed(4), locales: builtLocales().map((l) => l.code) }, null, 2) + "\n");
console.log(`coverage: built ${pct(built.low)}–${pct(built.high)}, phase 1 ${pct(phase1.low)}–${pct(phase1.high)}, all 30 ${pct(all.low)}–${pct(all.high)}`);
if (process.argv.includes("--strict") && all.low < 0.8) {
  console.error("coverage: the conservative estimate for all locales is below 80%");
  process.exitCode = 1;
}
