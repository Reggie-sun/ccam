import type { Metadata } from "next";
import { Landing, landingAlternates, landingTitle } from "../components/Landing";
import { ogPath, seo } from "../lib/seo";

// "/" is the English landing page and the x-default for every language.
export function generateMetadata(): Metadata {
  return { ...seo({ locale: "en", ...landingTitle("en"), path: "/", alternates: landingAlternates(), og: ogPath("en", "home"), type: "website" }), title: { absolute: landingTitle("en").title } };
}

export default function Home() {
  return <Landing locale="en" />;
}
