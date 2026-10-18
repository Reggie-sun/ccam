/** Three stacked gauges: the accounts and their limits. */
export function BrandMark() {
  return (
    <svg viewBox="0 0 20 20" aria-hidden="true">
      <rect x="2" y="3" width="16" height="3" fill="var(--rule)" />
      <rect x="2" y="3" width="11" height="3" fill="var(--amber)" />
      <rect x="2" y="8.5" width="16" height="3" fill="var(--rule)" />
      <rect x="2" y="8.5" width="4" height="3" fill="var(--amber)" />
      <rect x="2" y="14" width="16" height="3" fill="var(--rule)" />
      <rect x="2" y="14" width="16" height="3" fill="var(--ink)" />
    </svg>
  );
}

export function ThemeIcon() {
  return (
    <svg viewBox="0 0 20 20" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth="1.6">
      <circle cx="10" cy="10" r="4.2" />
      <path d="M10 1.5v2.2M10 16.3v2.2M1.5 10h2.2M16.3 10h2.2M4 4l1.5 1.5M14.5 14.5L16 16M4 16l1.5-1.5M14.5 5.5L16 4" />
    </svg>
  );
}
