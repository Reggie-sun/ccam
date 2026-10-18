// Tells Bing, Yandex and other IndexNow engines about pages updated in the
// last 7 days. Needs INDEXNOW_KEY; the key file is published by finalize.
import { abs } from "../lib/paths.ts";
import { routes } from "./lib.ts";

const key = process.env.INDEXNOW_KEY;
if (!key) {
  console.log("indexnow: INDEXNOW_KEY not set, skipped");
  process.exit(0);
}
const since = new Date(Date.now() - 7 * 864e5).toISOString().slice(0, 10);
const urlList = routes().filter((r) => r.updated >= since).map((r) => abs(r.path));
if (!urlList.length) {
  console.log("indexnow: nothing changed");
  process.exit(0);
}
const res = await fetch("https://api.indexnow.org/indexnow", {
  method: "POST",
  headers: { "Content-Type": "application/json; charset=utf-8" },
  body: JSON.stringify({ host: "ccam.dev", key, keyLocation: abs(`/${key}.txt`), urlList }),
});
console.log(`indexnow: ${urlList.length} URLs, HTTP ${res.status}`);
if (!res.ok && res.status !== 202) process.exitCode = 1;
