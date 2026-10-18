// `font` "brand" uses the subset Latin web font; every other stack is system
// fonts, so no multi-megabyte CJK, Arabic or Indic font is ever downloaded.
export type Font = "brand" | "latin-ext" | "cyrillic" | "ja" | "ko" | "zh-hans" | "zh-hant" | "arabic" | "urdu" | "indic" | "thai";

export type Locale = {
  code: string; // URL segment
  lang: string; // BCP 47 for <html lang> and hreflang
  og: string; // og:locale
  dir: "ltr" | "rtl";
  phase: 1 | 2;
  font: Font;
};

export const locales: Locale[] = [
  { code: "en", lang: "en", og: "en_US", dir: "ltr", phase: 1, font: "brand" },
  { code: "zh-hans", lang: "zh-Hans", og: "zh_CN", dir: "ltr", phase: 1, font: "zh-hans" },
  { code: "zh-hant", lang: "zh-Hant", og: "zh_TW", dir: "ltr", phase: 1, font: "zh-hant" },
  { code: "hi", lang: "hi", og: "hi_IN", dir: "ltr", phase: 1, font: "indic" },
  { code: "es", lang: "es", og: "es_ES", dir: "ltr", phase: 1, font: "brand" },
  { code: "ar", lang: "ar", og: "ar_AR", dir: "rtl", phase: 1, font: "arabic" },
  { code: "fr", lang: "fr", og: "fr_FR", dir: "ltr", phase: 1, font: "brand" },
  { code: "bn", lang: "bn", og: "bn_BD", dir: "ltr", phase: 1, font: "indic" },
  { code: "pt-br", lang: "pt-BR", og: "pt_BR", dir: "ltr", phase: 1, font: "brand" },
  { code: "ru", lang: "ru", og: "ru_RU", dir: "ltr", phase: 1, font: "cyrillic" },
  { code: "id", lang: "id", og: "id_ID", dir: "ltr", phase: 1, font: "brand" },
  { code: "ur", lang: "ur", og: "ur_PK", dir: "rtl", phase: 1, font: "urdu" },
  { code: "de", lang: "de", og: "de_DE", dir: "ltr", phase: 1, font: "brand" },
  { code: "ja", lang: "ja", og: "ja_JP", dir: "ltr", phase: 1, font: "ja" },
  { code: "ko", lang: "ko", og: "ko_KR", dir: "ltr", phase: 1, font: "ko" },
  { code: "vi", lang: "vi", og: "vi_VN", dir: "ltr", phase: 1, font: "latin-ext" },
  { code: "tr", lang: "tr", og: "tr_TR", dir: "ltr", phase: 1, font: "latin-ext" },
  { code: "it", lang: "it", og: "it_IT", dir: "ltr", phase: 1, font: "brand" },
  { code: "mr", lang: "mr", og: "mr_IN", dir: "ltr", phase: 2, font: "indic" },
  { code: "te", lang: "te", og: "te_IN", dir: "ltr", phase: 2, font: "indic" },
  { code: "ta", lang: "ta", og: "ta_IN", dir: "ltr", phase: 2, font: "indic" },
  { code: "gu", lang: "gu", og: "gu_IN", dir: "ltr", phase: 2, font: "indic" },
  { code: "kn", lang: "kn", og: "kn_IN", dir: "ltr", phase: 2, font: "indic" },
  { code: "fa", lang: "fa", og: "fa_IR", dir: "rtl", phase: 2, font: "arabic" },
  { code: "sw", lang: "sw", og: "sw_KE", dir: "ltr", phase: 2, font: "brand" },
  { code: "fil", lang: "fil", og: "fil_PH", dir: "ltr", phase: 2, font: "latin-ext" },
  { code: "th", lang: "th", og: "th_TH", dir: "ltr", phase: 2, font: "thai" },
  { code: "pl", lang: "pl", og: "pl_PL", dir: "ltr", phase: 2, font: "latin-ext" },
  { code: "uk", lang: "uk", og: "uk_UA", dir: "ltr", phase: 2, font: "cyrillic" },
  { code: "ha", lang: "ha", og: "ha_NG", dir: "ltr", phase: 2, font: "latin-ext" },
];

export const defaultLocale = locales[0];

export function getLocale(code: string): Locale | undefined {
  return locales.find((l) => l.code === code);
}

/** The language's own name, e.g. "한국어" for ko. */
export function nativeName(l: Locale): string {
  const name = new Intl.DisplayNames([l.lang], { type: "language" }).of(l.lang) ?? l.code;
  return name.charAt(0).toLocaleUpperCase(l.lang) + name.slice(1);
}

/** Aliases that redirect to a canonical locale code (written to _redirects). */
export const localeAliases: Record<string, string> = {
  "zh": "zh-hans", "zh-cn": "zh-hans", "zh-sg": "zh-hans", "zh-tw": "zh-hant", "zh-hk": "zh-hant",
  "pt": "pt-br", "en-us": "en", "en-gb": "en", "es-mx": "es", "fr-fr": "fr", "de-de": "de",
  "ja-jp": "ja", "ko-kr": "ko", "tl": "fil", "in": "id",
};
