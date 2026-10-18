// The only JavaScript on ccam.dev: copy buttons, the theme toggle and the
// "also available in your language" banner. Pages work without it.
(() => {
  const d = document;
  const store = {
    get: (k) => { try { return localStorage.getItem(k); } catch { return null; } },
    set: (k, v) => { try { localStorage.setItem(k, v); } catch {} },
  };

  for (const b of d.querySelectorAll("[data-copy]")) {
    if (!navigator.clipboard) continue;
    b.hidden = false;
    const label = b.textContent;
    b.addEventListener("click", async () => {
      const code = b.parentElement.querySelector("code");
      await navigator.clipboard.writeText(code.textContent.trim());
      b.textContent = b.dataset.copied;
      setTimeout(() => (b.textContent = label), 1600);
    });
  }

  const toggle = d.querySelector("[data-theme-toggle]");
  if (toggle) {
    toggle.hidden = false;
    toggle.addEventListener("click", () => {
      const root = d.documentElement;
      const dark = root.dataset.theme ? root.dataset.theme === "dark" : matchMedia("(prefers-color-scheme: dark)").matches;
      root.dataset.theme = dark ? "light" : "dark";
      store.set("ccam-theme", root.dataset.theme);
    });
  }

  const banner = d.getElementById("lang-banner");
  if (banner && !store.get("ccam-lang-dismissed")) {
    const current = d.documentElement.lang.toLowerCase();
    const links = [...d.querySelectorAll("#lang-list a[data-same]")];
    const pick = (navigator.languages || [])
      .map((lang) => lang.toLowerCase())
      .map((lang) => links.find((a) => a.hreflang.toLowerCase() === lang) ||
        links.find((a) => a.hreflang.toLowerCase().split("-")[0] === lang.split("-")[0]))
      .find(Boolean);
    if (pick && pick.hreflang.toLowerCase() !== current) {
      banner.querySelector("[data-banner-text]").textContent = pick.dataset.suggest;
      const a = banner.querySelector("[data-banner-link]");
      a.href = pick.getAttribute("href");
      a.lang = pick.lang;
      a.textContent = pick.dataset.open;
      banner.lang = pick.lang;
      banner.dir = pick.dataset.dir; // an English suggestion on an Arabic page reads left to right
      banner.hidden = false;
      banner.querySelector("[data-banner-dismiss]").addEventListener("click", () => {
        banner.hidden = true;
        store.set("ccam-lang-dismissed", "1");
      });
    }
  }
})();
