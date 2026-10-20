export const INSTALL = "curl -fsSL https://ccam.dev/install.sh | sh";
export const INSTALL_WINDOWS = "irm https://ccam.dev/install.ps1 | iex";

export function InstallCommand({ copy, copied, command = INSTALL }: { copy: string; copied: string; command?: string }) {
  return (
    <div className="cmd" dir="ltr">
      <code tabIndex={0}>{command}</code>
      <button type="button" className="copy pill pill-primary" data-copy data-copied={copied} hidden>
        {copy}
      </button>
    </div>
  );
}
