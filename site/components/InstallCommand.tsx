export const INSTALL = "curl -fsSL https://ccam.dev/install.sh | sh";

export function InstallCommand({ copy, copied, command = INSTALL }: { copy: string; copied: string; command?: string }) {
  return (
    <div className="cmd" dir="ltr">
      <code>{command}</code>
      <button type="button" className="copy" data-copy data-copied={copied} hidden>
        {copy}
      </button>
    </div>
  );
}
