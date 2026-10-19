import AxeBuilder from "@axe-core/playwright";
import { expect, test, type Page } from "@playwright/test";
import { readdirSync, existsSync } from "node:fs";
import { join } from "node:path";

// Every built page, found on disk so new locales are covered automatically.
function pages(dir = "dist/client", prefix = ""): string[] {
  const out: string[] = [];
  for (const e of readdirSync(dir, { withFileTypes: true })) {
    if (!e.isDirectory() || ["og", "assets", "fonts", "sitemaps", "_next", "404"].includes(e.name)) continue;
    const path = `${prefix}/${e.name}`;
    if (existsSync(join(dir, e.name, "index.html"))) out.push(`${path}/`);
    out.push(...pages(join(dir, e.name), path));
  }
  return prefix ? out : ["/", ...out];
}

const sample = pages().filter((p, i, all) => !p.startsWith("/en/") || i < all.length); // all pages

async function errorsOn(page: Page, path: string): Promise<string[]> {
  const errors: string[] = [];
  page.on("console", (m) => m.type() === "error" && errors.push(m.text()));
  page.on("pageerror", (e) => errors.push(e.message));
  const res = await page.goto(path);
  expect(res?.status(), path).toBe(200);
  await page.waitForLoadState("load");
  return errors;
}

for (const path of sample) {
  test(`${path} loads under the CSP, has no a11y violations, no sideways scroll`, async ({ page }) => {
    expect(await errorsOn(page, path)).toEqual([]);
    const overflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);
    expect(overflow, "horizontal overflow in px").toBeLessThanOrEqual(0);
    const axe = await new AxeBuilder({ page }).withTags(["wcag2a", "wcag2aa", "wcag21aa"]).analyze();
    const serious = axe.violations.filter((v) => v.impact === "serious" || v.impact === "critical");
    expect(serious.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(" ")).join(", ")}`)).toEqual([]);
  });
}

test("no sideways scroll at 320 px", async ({ page }) => {
  await page.setViewportSize({ width: 320, height: 700 });
  for (const path of ["/", "/en/docs/commands/", "/en/docs/sharing/"]) {
    await page.goto(path);
    expect(await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth), path).toBeLessThanOrEqual(0);
  }
});

test("copy button copies the install command", async ({ page, context, browserName }) => {
  test.skip(browserName !== "chromium");
  await context.grantPermissions(["clipboard-read", "clipboard-write"]);
  await page.goto("/");
  const button = page.locator(".cmd .copy");
  await expect(button).toBeVisible();
  await button.click();
  await expect(button).toHaveText("Copied");
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe("curl -fsSL https://ccam.dev/install.sh | sh");
});

test("theme toggle switches and remembers the theme", async ({ page }) => {
  await page.goto("/");
  const toggle = page.locator("[data-theme-toggle]");
  const bg = () => page.evaluate(() => getComputedStyle(document.body).backgroundColor);
  const before = await bg();
  await toggle.click();
  const after = await bg();
  expect(after).not.toBe(before);
  await page.reload();
  expect(await bg()).toBe(after);
  await page.locator("[data-theme-toggle]").click();
  expect(await bg()).toBe(before);
});

test("keyboard: skip link, language picker, nav", async ({ page }) => {
  await page.goto("/en/docs/install/");
  await page.keyboard.press("Tab");
  await expect(page.locator(".skip")).toBeFocused();
  await page.keyboard.press("Enter");
  await expect(page).toHaveURL(/#main$/);
  const summary = page.locator(".picker > summary");
  await summary.focus();
  await page.keyboard.press("Enter");
  await expect(page.locator(".picker")).toHaveAttribute("open", "");
  await expect(page.locator("#lang-list a").first()).toBeVisible();
  for (const name of ["Docs", "Guides", "FAQ"]) {
    const href = await page.locator(".site-nav a", { hasText: name }).getAttribute("href");
    expect((await page.request.get(href!)).status(), name).toBe(200);
  }
});

test("language banner suggests the reader's language when the page exists in it", async ({ browser }) => {
  const other = pages().find((p) => /^\/(ko|ja|de|es)\/$/.test(p));
  test.skip(!other, "only English is built");
  const lang = other!.slice(1, -1);
  const ctx = await browser.newContext({ locale: lang });
  const page = await ctx.newPage();
  await page.goto("/");
  await expect(page.locator("#lang-banner")).toBeVisible();
  await expect(page.locator("[data-banner-link]")).toHaveAttribute("href", other!);
  await page.locator("[data-banner-dismiss]").click();
  await page.reload();
  await expect(page.locator("#lang-banner")).toBeHidden();
  await ctx.close();
});

test("404 page", async ({ page }) => {
  const res = await page.goto("/en/docs/nope/");
  expect(res?.status()).toBe(404);
  await expect(page.locator("h1")).toHaveText("Page not found");
});

test("mobile menu opens, lists the sections and closes with Escape", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/en/docs/install/");
  await page.locator(".menu > summary").click();
  await expect(page.locator(".menu-panel a").first()).toBeVisible();
  await page.keyboard.press("Escape");
  await expect(page.locator(".menu")).not.toHaveAttribute("open", "");
});

test("on this page links to every section heading", async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto("/en/docs/commands/");
  const links = page.locator(".toc-rail a");
  await expect(links.first()).toBeVisible();
  for (const href of await links.evaluateAll((as) => as.map((a) => a.getAttribute("href")!))) {
    await expect(page.locator(href)).toHaveCount(1);
  }
});
