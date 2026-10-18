export const SITE = "https://ccam.dev";
export const REPO = "https://github.com/jhyunwoo/ccam";

/** Landing pages: English lives at "/", the others at "/<locale>/". */
export function homePath(locale: string): string {
  return locale === "en" ? "/" : `/${locale}/`;
}

export function pagePath(locale: string, section: string, slug?: string): string {
  return slug ? `/${locale}/${section}/${slug}/` : `/${locale}/${section}/`;
}

export function abs(path: string): string {
  return SITE + path;
}
