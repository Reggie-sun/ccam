// UI strings, server-side only: nothing here reaches the browser, so adding
// languages never grows the client bundle.
import { IntlMessageFormat } from "intl-messageformat";
import type en from "../messages/en.json";
import { getLocale, type Locale } from "../i18n/locales";

export type Messages = typeof en;

const files = import.meta.glob<{ default: Messages }>("../messages/*.json", { eager: true });

export function hasMessages(code: string): boolean {
  return `../messages/${code}.json` in files;
}

export function getMessages(code: string): Messages {
  const file = files[`../messages/${code}.json`];
  if (!file) throw new Error(`messages/${code}.json is missing`);
  return file.default;
}

/** Formats an ICU MessageFormat string such as "{count, plural, one {#} other {#}}". */
export function format(locale: Locale | string, message: string, values: Record<string, string | number> = {}): string {
  const lang = typeof locale === "string" ? (getLocale(locale)?.lang ?? locale) : locale.lang;
  return String(new IntlMessageFormat(message, lang).format(values));
}

export function formatDate(locale: Locale, isoDate: string): string {
  return new Intl.DateTimeFormat(locale.lang, { dateStyle: "long", timeZone: "UTC" }).format(new Date(`${isoDate}T00:00:00Z`));
}
