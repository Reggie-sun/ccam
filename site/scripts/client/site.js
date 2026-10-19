// The only JavaScript on ccam.dev. Pages work without it; with it they get the
// copy buttons, the theme toggle, the language banner and the motion layer.
(() => {
  const d = document;
  const store = {
    get: (k) => { try { return localStorage.getItem(k); } catch { return null; } },
    set: (k, v) => { try { localStorage.setItem(k, v); } catch {} },
  };
  const still = matchMedia("(prefers-reduced-motion: reduce)").matches;
  const IO = window.IntersectionObserver;

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

  // Header popovers: one open at a time, closed by Escape or a click elsewhere.
  const pops = [...d.querySelectorAll("details.picker, details.menu")];
  for (const p of pops) p.addEventListener("toggle", () => p.open && pops.forEach((o) => o !== p && (o.open = false)));
  d.addEventListener("click", (e) => pops.forEach((p) => p.open && !p.contains(e.target) && (p.open = false)));
  d.addEventListener("keydown", (e) => {
    if (e.key !== "Escape") return;
    for (const p of [...pops, ...d.querySelectorAll(".sidebar details, .toc-inline")]) {
      if (p.open) { p.open = false; p.querySelector("summary").focus(); }
    }
  });
  for (const a of d.querySelectorAll(".toc-inline a, .sidebar a")) a.addEventListener("click", () => (a.closest("details").open = false));

  // On this page: mark the section being read.
  const tocLinks = [...d.querySelectorAll(".toc a")];
  const heads = [...new Set(tocLinks.map((a) => decodeURIComponent(a.hash.slice(1))))].map((id) => d.getElementById(id)).filter(Boolean);
  if (heads.length) {
    let current, queued = false;
    const mark = () => {
      queued = false;
      const line = innerHeight * 0.3;
      const head = heads.filter((h) => h.getBoundingClientRect().top < line).pop();
      if (head === current) return;
      current = head;
      for (const a of tocLinks) a.removeAttribute("aria-current");
      if (head) for (const a of tocLinks) if (a.hash === `#${head.id}`) a.setAttribute("aria-current", "location");
    };
    addEventListener("scroll", () => { if (!queued) { queued = true; requestAnimationFrame(mark); } }, { passive: true });
    mark();
  }

  // How it works: the step in the middle of the screen drives the sticky panel.
  const flow = d.querySelector("[data-flow]");
  if (flow && IO) {
    const steps = [...flow.querySelectorAll(".step")];
    const views = [...flow.querySelectorAll(".fv")];
    const io = new IO((entries) => {
      for (const e of entries) {
        if (!e.isIntersecting) continue;
        const i = steps.indexOf(e.target);
        flow.dataset.active = i;
        steps.forEach((s, k) => s.classList.toggle("on", k === i));
        views.forEach((v, k) => v.classList.toggle("on", k === i));
      }
    }, { rootMargin: "-45% 0px -45% 0px" });
    steps.forEach((s) => io.observe(s));
  }

  if (still) return;

  // Hero orbs lean toward the pointer.
  const drift = d.querySelector("[data-parallax]");
  if (drift && matchMedia("(pointer: fine)").matches) {
    let x = 0, y = 0, queued = false;
    addEventListener("pointermove", (e) => {
      x = e.clientX / innerWidth - 0.5;
      y = e.clientY / innerHeight - 0.5;
      if (queued) return;
      queued = true;
      requestAnimationFrame(() => {
        queued = false;
        drift.style.setProperty("--px", x.toFixed(3));
        drift.style.setProperty("--py", y.toFixed(3));
      });
    }, { passive: true });
  }

  // Guide cards: a soft light follows the pointer.
  for (const c of d.querySelectorAll("[data-spot]")) {
    c.addEventListener("pointermove", (e) => {
      const r = c.getBoundingClientRect();
      c.style.setProperty("--mx", `${e.clientX - r.left}px`);
      c.style.setProperty("--my", `${e.clientY - r.top}px`);
    }, { passive: true });
  }

  // The usage demo replays as a terminal session each time it scrolls into view.
  const demo = d.querySelector("[data-demo]");
  if (demo && IO) {
    const typed = [...demo.querySelectorAll("[data-type]")];
    const texts = typed.map((t) => t.textContent);
    const later = [...demo.querySelectorAll(".later")];
    const rows = [...demo.querySelectorAll(".row")];
    const pcts = [...demo.querySelectorAll(".pct")];
    const pick = demo.querySelector("[data-pick]");
    const note = demo.querySelector(".note");
    const last = later[later.length - 1];
    const fmt = new Intl.NumberFormat(d.documentElement.lang, { style: "percent" });
    let timers = [];
    const at = (ms, fn) => timers.push(setTimeout(fn, ms));
    const type = (el, text, start) => {
      for (let i = 0; i <= text.length; i++) at(start + i * 55, () => (el.textContent = text.slice(0, i)));
      return start + text.length * 55;
    };
    const reset = () => {
      timers.forEach(clearTimeout);
      timers = [];
      demo.classList.add("playing");
      later.forEach((e) => e.classList.remove("shown"));
      pick.classList.remove("on");
      typed.forEach((t) => (t.textContent = ""));
      pcts.forEach((p) => (p.textContent = fmt.format(0)));
    };
    const play = () => {
      reset();
      let t = type(typed[0], texts[0], 350) + 300;
      at(t, () => rows.forEach((r) => r.classList.add("shown")));
      for (const p of pcts) {
        const v = Number(p.dataset.value);
        for (let k = 1; k <= 24; k++) at(t + 300 + k * 50, () => (p.textContent = fmt.format((v * (1 - (1 - k / 24) ** 3)) / 100)));
      }
      t += 1900;
      at(t, () => note.classList.add("shown"));
      t += 800;
      at(t, () => last.classList.add("shown"));
      t = type(typed[1], texts[1], t + 400);
      at(t + 450, () => pick.classList.add("on"));
    };
    reset();
    let seen = false;
    new IO(([e]) => {
      if (e.intersectionRatio >= 0.4 && !seen) { seen = true; play(); }
      else if (!e.isIntersecting && seen) { seen = false; reset(); }
    }, { threshold: [0, 0.4] }).observe(demo);
  }
})();
