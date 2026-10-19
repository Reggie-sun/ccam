// www.ccam.dev: a separate Worker (deploy/www.jsonc) that sends every request
// to the same path on ccam.dev, so the site's pages keep skipping the Worker.
export default {
  fetch(req: Request): Response {
    const url = new URL(req.url);
    url.protocol = "https:";
    url.hostname = "ccam.dev";
    url.port = "";
    return Response.redirect(url.toString(), 301);
  },
};
