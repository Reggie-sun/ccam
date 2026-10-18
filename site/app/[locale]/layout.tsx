import { locales } from "../../i18n/locales";
import { hasMessages } from "../../lib/i18n";

// Only locales with a messages file are built; anything else is a 404.
export const dynamicParams = false;
export function generateStaticParams() {
  return locales.filter((l) => hasMessages(l.code)).map((l) => ({ locale: l.code }));
}

export default function LocaleLayout({ children }: { children: React.ReactNode }) {
  return children;
}
