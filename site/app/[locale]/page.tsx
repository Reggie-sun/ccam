import type { Metadata } from "next";
import { Landing, landingAlternates, landingTitle } from "../../components/Landing";
import { homePath } from "../../lib/paths";
import { ogPath, seo } from "../../lib/seo";

type Props = { params: Promise<{ locale: string }> };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params;
  const t = landingTitle(locale);
  return { ...seo({ locale, ...t, path: homePath(locale), alternates: landingAlternates(), og: ogPath(locale, "home"), type: "website" }), title: { absolute: t.title } };
}

export default async function LocaleHome({ params }: Props) {
  const { locale } = await params;
  return <Landing locale={locale} />;
}
