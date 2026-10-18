// Serves dist/client like Workers Static Assets: directory index.html,
// trailing-slash redirects, 404.html, _redirects, and the "/*" headers from
// _headers (so tests run under the real CSP). Usage: node scripts/serve.ts [port]
import { createServer } from "node:http";
import { readFile, stat } from "node:fs/promises";
import { extname, join, normalize } from "node:path";

const root = "dist/client";
const port = Number(process.argv[2] ?? 4321);
const types: Record<string, string> = {
  ".html": "text/html; charset=utf-8", ".js": "text/javascript", ".css": "text/css", ".svg": "image/svg+xml",
  ".png": "image/png", ".woff2": "font/woff2", ".xml": "application/xml", ".txt": "text/plain; charset=utf-8",
  ".md": "text/markdown; charset=utf-8", ".json": "application/json", ".avif": "image/avif", ".webp": "image/webp",
};

const headers: Record<string, string> = {};
try {
  const block = (await readFile(join(root, "_headers"), "utf8")).split("\n\n")[0].split("\n").slice(1);
  for (const line of block) {
    const i = line.indexOf(":");
    if (i > 0) headers[line.slice(0, i).trim()] = line.slice(i + 1).trim();
  }
} catch {}
const redirects: [RegExp, string][] = [];
try {
  for (const line of (await readFile(join(root, "_redirects"), "utf8")).split("\n").filter(Boolean)) {
    const [from, to] = line.split(" ");
    redirects.push([new RegExp("^" + from.replace("*", "(.*)") + "$"), to]);
  }
} catch {}

async function file(path: string) {
  try {
    const s = await stat(path);
    return s.isFile() ? path : null;
  } catch {
    return null;
  }
}

createServer(async (req, res) => {
  const url = new URL(req.url ?? "/", "http://x");
  const p = decodeURIComponent(url.pathname);
  for (const [re, to] of redirects) {
    const m = re.exec(p);
    if (m) {
      res.writeHead(301, { Location: to.replace(":splat", m[1] ?? "") }).end();
      return;
    }
  }
  const safe = normalize(p).replace(/^(\.\.[/\\])+/, "");
  let target = await file(join(root, safe));
  if (!target && !p.endsWith("/") && (await file(join(root, safe, "index.html")))) {
    res.writeHead(307, { Location: p + "/" }).end();
    return;
  }
  if (!target) target = await file(join(root, safe, "index.html"));
  const status = target ? 200 : 404;
  target ??= join(root, "404.html");
  res.writeHead(status, { ...headers, "Content-Type": types[extname(target)] ?? "application/octet-stream" });
  res.end(await readFile(target));
}).listen(port, () => console.log(`http://localhost:${port}`));
