---
title: Updating ccam
description: ccam update replaces ccam with the latest GitHub release after a SHA-256 check; --check only reports and --version installs a given release tag.
updated: 2026-10-21
order: 10
translation: source
---

`ccam update` downloads a ccam release from GitHub, checks it against the release's `checksums.txt` and replaces the running binary in one rename. Use it when you installed ccam with the install script or from an archive. The [dashboard](/features/dashboard) also asks GitHub for the latest release each time it opens and offers it with `U`, which runs the same steps after asking. No other command contacts GitHub, so they work the same offline.

## Usage

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

When no newer release exists, ccam prints `ccam 0.4.0 is up to date.` and changes nothing. With `--version` and the version you already run, it prints `ccam 0.4.0 is already installed.`

## Options

| Option | What it does | Default |
| --- | --- | --- |
| `--check` | Only report whether a newer release exists | off |
| `--version <tag>` | Install this release tag, such as `v0.3.0`, newer or older than the current one | latest |
| `--dry-run` | Honoured: looks up the latest tag (unless `--version` names one), runs the same checks on the binary, and prints the download, checksum and replace paths without downloading | off |

`--check` and `--version` cannot be used together. The tag must look like `vMAJOR.MINOR.PATCH`, with an optional pre-release or build suffix; the leading `v` may be left out. Anything else fails, for example `--version takes a release tag such as v0.3.0, not "1.2"`.

The environment variable `CCAM_RELEASES_URL` replaces `https://github.com/jhyunwoo/ccam/releases`, for a mirror or for tests. It must be laid out the same way: `<url>/latest` redirects to `<url>/tag/<tag>`, and assets are at `<url>/download/<tag>/<asset>`.

## How it works

The code is in `internal/update`, which only `ccam update` imports. The launch path never calls it.

1. `cli.updateCmd` finds the running binary with `os.Executable` and resolves symlinks with `filepath.EvalSymlinks`. That real path is what gets replaced.
2. Without `--version`, `Updater.Latest` requests `<releases>/latest` without following the redirect, and reads the tag from the `Location` header (`…/tag/v0.5.0`). This uses the releases page instead of the GitHub API, so the API's limit of 60 requests an hour does not apply. The lookup has 10 seconds.
3. `parseVersion` checks the tag against a strict pattern. A tag ends up in download URLs and in a printed shell command, so nothing with `/`, `%`, spaces or shell characters gets through. `compareVersions` compares it with the running version; a pre-release sorts before its release.
4. `managed()` and `writable()` decide whether ccam may replace the binary at all (see below).
5. `Updater.Fetch` downloads `ccam_<os>_<arch>.tar.gz` (on Windows `ccam_windows_<arch>.zip`) and `checksums.txt` at the same time, with 2 minutes for both. Each download is capped at 256 MB.
6. The archive's SHA-256 must match its line in `checksums.txt` before anything else happens. A missing line or a different checksum stops the update with `nothing was changed`.
7. `extract()` reads the tar and takes the first regular file named `ccam` (`extractZip` takes `ccam.exe` from the zip). Directories, links and other entries are ignored, and the binary is capped at 256 MB too.
8. `Updater.Install` writes the new binary with `fsutil.AtomicWrite`: a temporary file `.ccam.tmp*` next to the binary, with the old file's permission bits, flushed to disk and renamed over the old one.

The whole pipeline, from tag lookup to the rename:

<!-- motion:update-pipeline -->

The archive and `checksums.txt` are downloaded side by side, and nothing is written until the checksum matches. The rename is atomic, so a crash or a full disk leaves either the old binary or the new one, never half of one. Shortcut commands such as `~/.local/bin/ccwork` are symlinks to the ccam binary, so they run the new version without any change.

On Windows a running program's file cannot be replaced or deleted, only renamed. `Updater.Install` writes the new binary to `ccam.exe.new`, renames the running `ccam.exe` to `ccam.exe.old`, and renames the new file to `ccam.exe`. If the last step fails, the old binary is renamed back. The next update deletes `ccam.exe.old`. Shortcuts that are hard links instead of symbolic links ([Shortcut commands](/features/shortcuts)) are renewed right after, so they also run the new version.

### Package-managed installs

A ccam installed from the `.deb` belongs to apt, and replacing it behind apt's back would leave the package database wrong. Before downloading, ccam checks who owns the binary:

<!-- motion:update-managed -->

1. A binary in `/usr/bin`, `/bin`, `/usr/sbin` or `/sbin` counts as package-managed.
2. Elsewhere, if `dpkg` is installed and `dpkg -S <path>` succeeds, it also counts as package-managed.
3. Otherwise ccam creates and removes a temporary `.ccam-update-*` file next to the binary to make sure it can write there.

For a package-managed binary ccam stops and prints the apt steps:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

On systems other than Linux, only the install script line is printed: `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex` on Windows. When the folder is not writable, ccam suggests rerunning the install script, which installs ccam in `~/.local/bin`.

## Safety and edge cases

- A development build reports its version as `dev` and has no release to compare with, so `ccam update` and `ccam update --check` refuse it. `ccam update --version vX.Y.Z` installs a release over it.
- A failed download, a missing asset for your platform, or a checksum mismatch leaves the installed binary as it was.
- A release that does not exist fails with `there is no release <tag>`. When GitHub refuses with HTTP 403 or 429, ccam says it is probably a rate limit and asks you to try again later.
- `--check` only reads. It does not check whether the binary is package-managed or writable.
- `--version` with an older tag is how you roll back.

## Related

- [Install](/docs/install)
- [Commands](/docs/commands)
- [FAQ](/docs/faq)
