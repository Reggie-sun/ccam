type Kind = "mint" | "peach" | "lavender" | "sky" | "rose";

/** DESIGN.md's atmospheric gradient orbs: decoration only, never behind a control. */
export function Orbs({ kinds, parallax }: { kinds: Kind[]; parallax?: boolean }) {
  return (
    <div className="orbs" aria-hidden="true">
      <div className="orbs-drift" data-parallax={parallax ? "" : undefined}>
        {kinds.map((k) => (
          <span key={k} className={`orb ${k}`} />
        ))}
      </div>
    </div>
  );
}
