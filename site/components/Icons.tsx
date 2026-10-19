const line = { fill: "none", stroke: "currentColor", strokeWidth: 1.6, strokeLinecap: "round", strokeLinejoin: "round" } as const;

/** Sun and moon; CSS shows the one for the current theme. */
export function ThemeIcon() {
  return (
    <>
      <svg className="sun" viewBox="0 0 20 20" aria-hidden="true" {...line}>
        <circle cx="10" cy="10" r="4.2" />
        <path d="M10 1.5v2.2M10 16.3v2.2M1.5 10h2.2M16.3 10h2.2M4 4l1.5 1.5M14.5 14.5L16 16M4 16l1.5-1.5M14.5 5.5L16 4" />
      </svg>
      <svg className="moon" viewBox="0 0 20 20" aria-hidden="true" {...line}>
        <path d="M16.5 12.4A7 7 0 0 1 7.6 3.5a7 7 0 1 0 8.9 8.9z" />
      </svg>
    </>
  );
}

export function MenuIcon() {
  return (
    <svg viewBox="0 0 20 20" aria-hidden="true" {...line}>
      <path d="M3 7h14M3 13h14" />
    </svg>
  );
}

export function GlobeIcon() {
  return (
    <svg viewBox="0 0 20 20" aria-hidden="true" {...line}>
      <circle cx="10" cy="10" r="7.5" />
      <path d="M2.5 10h15M10 2.5c2 2.2 3 4.7 3 7.5s-1 5.3-3 7.5c-2-2.2-3-4.7-3-7.5s1-5.3 3-7.5z" />
    </svg>
  );
}

export function ChevronIcon() {
  return (
    <svg className="chev" viewBox="0 0 20 20" aria-hidden="true" {...line}>
      <path d="M5.5 8l4.5 4.5L14.5 8" />
    </svg>
  );
}

export function ArrowIcon() {
  return (
    <svg viewBox="0 0 20 20" aria-hidden="true" {...line}>
      <path d="M4 10h12M11 5l5 5-5 5" />
    </svg>
  );
}
