// The only JavaScript on ccam.dev. Pages work without it; with it they get the
// copy buttons, the theme toggle, the language banner, search and the motion layer.
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

  // Search: Ctrl K, ⌘K or / opens the dialog; the locale's index loads on first
  // use. Every word must appear in a section's page title, heading or text.
  const dialog = d.querySelector("dialog.search");
  const opener = d.querySelector("[data-search-open]");
  if (dialog && opener && dialog.showModal && window.fetch) {
    opener.hidden = false;
    if (/Mac|iPhone|iPad/.test(navigator.platform)) opener.querySelector("[data-search-key]").textContent = "⌘K";
    const input = dialog.querySelector("input");
    const list = dialog.querySelector("[role=listbox]");
    const status = dialog.querySelector("[role=status]");
    const norm = (s) => s.normalize("NFKC").toLowerCase();
    let entries, loading, hits = [], active = -1;
    const load = () => (loading ??= fetch(dialog.dataset.searchIndex)
      .then((r) => r.json())
      .then(({ pages }) => {
        entries = pages.flatMap((p, order) => p.sections.map((s) => ({
          url: p.url + (s.id ? `#${s.id}` : ""), page: p.title, group: p.group, heading: s.heading, text: s.text, order,
          t: norm(p.title), h: norm(s.heading), x: norm(s.text),
        })));
        run();
      }, () => (loading = undefined))); // offline: try again on the next open

    const score = (e, terms) => {
      let s = 0;
      for (const t of terms) {
        const h = e.h.includes(t), title = e.t.includes(t), x = e.x.includes(t);
        if (!h && !title && !x) return 0;
        // A page title counts fully for the page's opening section only, so one page does not fill the list.
        s += (h ? 10 : 0) + (title ? (e.heading ? 2 : 10) : 0) + (x ? 1 : 0);
      }
      return s;
    };
    // Text with each search word in <mark>; lowercasing can change a string's
    // length in a few scripts, and then the text is shown plain.
    const marked = (el, str, terms) => {
      const low = str.toLowerCase();
      const on = new Array(str.length).fill(false);
      if (low.length === str.length) for (const t of terms) for (let i = low.indexOf(t); i >= 0; i = low.indexOf(t, i + 1)) on.fill(true, i, i + t.length);
      for (let i = 0; i < str.length;) {
        let j = i;
        while (j < str.length && on[j] === on[i]) j++;
        const part = str.slice(i, j);
        el.append(on[i] ? Object.assign(d.createElement("mark"), { textContent: part }) : part);
        i = j;
      }
      return el;
    };
    const span = (cls, text = "") => Object.assign(d.createElement("span"), { className: cls, textContent: text });
    const snippet = (e, terms) => {
      const at = e.x.length === e.text.length ? Math.min(...terms.map((t) => e.x.indexOf(t)).filter((i) => i >= 0)) : 0;
      // Start a little before the first match, at a word break when one is near (CJK and Thai have none).
      let start = Number.isFinite(at) && at > 50 ? at - 40 : 0;
      const space = e.text.lastIndexOf(" ", start);
      if (start && space > start - 20) start = space + 1;
      return (start ? "…" : "") + e.text.slice(start, start + 160);
    };
    const select = (i) => {
      active = i;
      [...list.children].forEach((li, k) => li.setAttribute("aria-selected", String(k === i)));
      if (i >= 0) {
        input.setAttribute("aria-activedescendant", list.children[i].id);
        list.children[i].scrollIntoView({ block: "nearest" });
      } else input.removeAttribute("aria-activedescendant");
    };
    const run = () => {
      const terms = norm(input.value).split(/\s+/).filter(Boolean);
      hits = [];
      if (entries && terms.length) {
        const perPage = new Map();
        const ranked = entries.map((e) => [score(e, terms), e]).filter(([s]) => s > 0).sort((a, b) => b[0] - a[0] || a[1].order - b[1].order);
        for (const [, e] of ranked) {
          const n = perPage.get(e.page) ?? 0;
          if (n >= 3) continue;
          perPage.set(e.page, n + 1);
          hits.push(e);
          if (hits.length === 12) break;
        }
      }
      list.replaceChildren(...hits.map((e, i) => {
        const li = Object.assign(d.createElement("li"), { id: `search-hit-${i}` });
        li.setAttribute("role", "option");
        li.append(
          span("hit-path", e.heading ? `${e.group} › ${e.page}` : e.group),
          marked(span("hit-title"), e.heading || e.page, terms),
        );
        if (e.text) li.append(marked(span("hit-text"), snippet(e, terms), terms));
        li.addEventListener("click", () => go(e));
        return li;
      }));
      input.setAttribute("aria-expanded", String(hits.length > 0));
      select(hits.length ? 0 : -1);
      status.textContent = entries && terms.length && !hits.length ? status.dataset.empty : "";
    };
    const go = (e) => {
      dialog.close();
      location.href = e.url;
    };
    const open = () => {
      if (dialog.open) return;
      dialog.showModal();
      input.select();
      load();
    };

    opener.addEventListener("click", open);
    opener.addEventListener("pointerenter", load, { once: true });
    dialog.querySelector("[data-search-close]").addEventListener("click", () => dialog.close());
    dialog.addEventListener("click", (e) => e.target === dialog && dialog.close()); // the backdrop
    input.addEventListener("input", run);
    input.addEventListener("keydown", (e) => {
      if (e.key === "ArrowDown" || e.key === "ArrowUp") {
        e.preventDefault();
        if (hits.length) select((active + (e.key === "ArrowDown" ? 1 : hits.length - 1)) % hits.length);
      } else if (e.key === "Enter" && hits[active]) {
        e.preventDefault();
        go(hits[active]);
      } else if (e.key === "Escape") {
        e.preventDefault(); // a search field would only clear its text
        dialog.close();
      }
    });
    d.addEventListener("keydown", (e) => {
      const typing = e.target.closest?.("input, textarea, select, [contenteditable]");
      if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || (e.key === "/" && !typing)) {
        e.preventDefault();
        open();
      }
    });
  }

  // Everything below waits for an idle moment so it never lengthens the first task.
  const idle = window.requestIdleCallback || ((fn) => setTimeout(fn, 200));
  idle(() => {
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

    // The demo replays as a terminal session each time it scrolls into view.
    const demo = d.querySelector("[data-demo]");
    if (demo && IO) {
      const seq = [...demo.querySelectorAll("[data-seq]")];
      const typed = [...demo.querySelectorAll("[data-type]")];
      const texts = new Map(typed.map((t) => [t, t.textContent]));
      const pcts = [...demo.querySelectorAll(".pct")];
      const pick = demo.querySelector("[data-pick]");
      const steps = [...d.querySelectorAll("[data-demo-steps] li")];
      const step = (i) => steps.forEach((s, k) => s.classList.toggle("on", k === i));
      let fmt; // created on first play: Intl setup is a long task on slow phones
      let timers = [];
      const at = (ms, fn) => timers.push(setTimeout(fn, ms));
      const type = (el, start) => {
        const text = texts.get(el);
        for (let i = 0; i <= text.length; i++) at(start + i * 45, () => (el.textContent = text.slice(0, i)));
        return start + text.length * 45;
      };
      const count = (p, start) => {
        const v = Number(p.dataset.value);
        for (let k = 1; k <= 24; k++) at(start + k * 50, () => (p.textContent = fmt.format((v * (1 - (1 - k / 24) ** 3)) / 100)));
      };
      const reset = () => {
        timers.forEach(clearTimeout);
        timers = [];
        demo.classList.add("playing");
        seq.forEach((e) => e.classList.remove("shown"));
        pick.classList.remove("on");
        steps[0]?.parentElement.classList.add("live");
        step(-1);
        typed.forEach((t) => (t.textContent = ""));
        pcts.forEach((p) => (p.textContent = ""));
      };
      const play = () => {
        reset();
        fmt ??= new Intl.NumberFormat(d.documentElement.lang, { style: "percent" });
        let t = 300;
        for (const el of seq) {
          at(t, () => {
            el.classList.add("shown");
            if (el.dataset.step) step(Number(el.dataset.step));
          });
          for (const tp of el.querySelectorAll("[data-type]")) t = type(tp, t + 250);
          if (el.classList.contains("rows")) pcts.forEach((p) => count(p, t + 300));
          t += Number(el.dataset.seq);
        }
        at(t - 800, () => pick.classList.add("on"));
      };
      reset();
      let seen = false;
      new IO(([e]) => {
        if (e.intersectionRatio >= 0.3 && !seen) { seen = true; play(); }
        else if (!e.isIntersecting && seen) { seen = false; reset(); }
      }, { threshold: [0, 0.3] }).observe(demo);
    }
  });
})();
