import assert from "node:assert/strict";
import { test } from "node:test";
import worker, { installScript, latest } from "../worker/index.ts";
import www from "../worker/www.ts";

test("install.sh comes from the latest release, else main", async () => {
  const seen: string[] = [];
  const fake = (async (url: string) => {
    seen.push(url);
    return url.includes("releases") ? new Response("", { status: 404 }) : new Response("#!/bin/sh\n");
  }) as typeof fetch;
  const res = await installScript(fake);
  assert.equal(res.status, 200);
  assert.equal(res.headers.get("Content-Type"), "text/x-shellscript; charset=utf-8");
  assert.equal(await res.text(), "#!/bin/sh\n");
  assert.equal(seen.length, 2);
});

test("install.ps1 comes from the latest release as plain text", async () => {
  const seen: string[] = [];
  const fake = (async (url: string) => {
    seen.push(url);
    return new Response("# Installs ccam\n");
  }) as typeof fetch;
  const res = await installScript(fake, "install.ps1");
  assert.equal(res.status, 200);
  assert.equal(res.headers.get("Content-Type"), "text/plain; charset=utf-8");
  assert.deepEqual(seen, ["https://github.com/jhyunwoo/ccam/releases/latest/download/install.ps1"]);
});

test("install.sh failure is a 502, not an HTML page piped into sh", async () => {
  const res = await installScript((async () => new Response("", { status: 500 })) as typeof fetch);
  assert.equal(res.status, 502);
});

test("latest redirects", () => {
  assert.equal(latest("/latest").headers.get("Location"), "https://github.com/jhyunwoo/ccam/releases/latest");
  assert.equal(latest("/latest/ccam_linux_amd64.tar.gz").headers.get("Location"), "https://github.com/jhyunwoo/ccam/releases/latest/download/ccam_linux_amd64.tar.gz");
  assert.equal(latest("/latest/../../x").status, 404);
});

test("other paths go to static assets", async () => {
  const res = await worker.fetch(new Request("https://ccam.dev/en/"), { ASSETS: { fetch: async () => new Response("asset") } });
  assert.equal(await res.text(), "asset");
});

test("www redirects to the same path on ccam.dev", () => {
  const res = www.fetch(new Request("http://www.ccam.dev/ko/docs/faq/?q=1"));
  assert.equal(res.status, 301);
  assert.equal(res.headers.get("Location"), "https://ccam.dev/ko/docs/faq/?q=1");
});
