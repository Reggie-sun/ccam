import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { DocPage, docAlternates } from "../../../../components/DocPage";
import { allPages, getPage, type Section } from "../../../../lib/content";
import { hasMessages } from "../../../../lib/i18n";
import { pagePath } from "../../../../lib/paths";
import { ogPath, pageTitle, seo } from "../../../../lib/seo";

type Props = { params: Promise<{ locale: string; section: string; slug: string }> };

export const dynamicParams = false;
export function generateStaticParams() {
  return allPages()
    .filter((p) => hasMessages(p.locale))
    .map((p) => ({ locale: p.locale, section: p.section, slug: p.slug }));
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale, section, slug } = await params;
  const page = getPage(locale, section as Section, slug);
  if (!page) return {};
  return seo({
    locale,
    title: pageTitle(locale, page.metaTitle ?? page.title),
    description: page.description,
    path: pagePath(locale, section, slug),
    alternates: docAlternates(section as Section, slug),
    og: ogPath(locale, `${section}-${slug}`),
  });
}

export default async function ContentPage({ params }: Props) {
  const { locale, section, slug } = await params;
  const page = getPage(locale, section as Section, slug);
  if (!page) notFound();
  return <DocPage page={page} />;
}
