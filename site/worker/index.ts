// Runs only for the paths listed in run_worker_first (deploy/wrangler.jsonc);
// every page is served straight from static assets.
const REPO = "jhyunwoo/ccam";
const RELEASES = `https://github.com/${REPO}/releases`;

type Env = { ASSETS: { fetch(req: Request): Promise<Response> } };

export async function installScript(fetcher: typeof fetch = fetch): Promise<Response> {
  // The script attached to the latest release; main's copy until a release exists.
  for (const url of [`${RELEASES}/latest/download/install.sh`, `https://raw.githubusercontent.com/${REPO}/main/install.sh`]) {
    const res = await fetcher(url, { cf: { cacheTtl: 3600, cacheEverything: true } } as RequestInit);
    if (res.ok) {
      return new Response(res.body, {
        headers: {
          "Content-Type": "text/x-shellscript; charset=utf-8",
          "Cache-Control": "public, max-age=300, s-maxage=3600",
          "X-Content-Type-Options": "nosniff",
        },
      });
    }
  }
  return new Response("install.sh is temporarily unavailable; see https://github.com/jhyunwoo/ccam\n", { status: 502 });
}

export function latest(path: string): Response {
  const asset = path.replace(/^\/latest\/?/, "");
  if (asset && !/^[\w.-]+$/.test(asset)) return new Response("Not found\n", { status: 404 });
  return Response.redirect(asset ? `${RELEASES}/latest/download/${asset}` : `${RELEASES}/latest`, 302);
}

export default {
  async fetch(req: Request, env: Env): Promise<Response> {
    const { pathname } = new URL(req.url);
    if (pathname === "/install.sh") return installScript();
    if (pathname === "/latest" || pathname.startsWith("/latest/")) return latest(pathname);
    return env.ASSETS.fetch(req);
  },
};
