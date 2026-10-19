#!/bin/sh
# Installs ccam into ~/.local/bin on Linux or macOS.
#
#   curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
#
# Environment:
#   CCAM_VERSION      release tag such as v1.0.0 (default: latest)
#   CCAM_INSTALL_DIR  where to install (default: ~/.local/bin)
#   CCAM_BASE_URL     release download base (default: GitHub releases)
#   CCAM_TARBALL      install from a local tar.gz instead of downloading
set -eu

repo="jhyunwoo/ccam"
dir="${CCAM_INSTALL_DIR:-$HOME/.local/bin}"

say() { printf '%s\n' "$*"; }
die() { printf 'ccam install failed: %s\n' "$*" >&2; exit 1; }

case "$(uname -s)" in
  Linux) os=linux ;;
  Darwin) os=darwin ;;
  *) die "only Linux and macOS are supported" ;;
esac
case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  *) die "unsupported CPU: $(uname -m)" ;;
esac

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
asset="ccam_${os}_${arch}.tar.gz"

fetch() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$1" -o "$2"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$1" -O "$2"
  else
    die "curl or wget is required"
  fi
}

if [ -n "${CCAM_TARBALL:-}" ]; then
  cp "$CCAM_TARBALL" "$tmp/$asset"
else
  if [ -n "${CCAM_VERSION:-}" ]; then
    base="${CCAM_BASE_URL:-https://github.com/$repo/releases/download}/$CCAM_VERSION"
  else
    base="${CCAM_BASE_URL:-https://github.com/$repo/releases/latest/download}"
  fi
  say "→ downloading $base/$asset"
  fetch "$base/$asset" "$tmp/$asset" || die "download failed"
  fetch "$base/checksums.txt" "$tmp/checksums.txt" || die "could not download checksums.txt"
  grep " $asset\$" "$tmp/checksums.txt" > "$tmp/sum.txt" || die "$asset is not in checksums.txt"
  if command -v sha256sum >/dev/null 2>&1; then
    (cd "$tmp" && sha256sum -c --quiet sum.txt) || die "checksum mismatch"
  else # macOS
    (cd "$tmp" && shasum -a 256 -c -s sum.txt) || die "checksum mismatch"
  fi
fi

tar -xzf "$tmp/$asset" -C "$tmp" ccam || die "could not extract the archive"
mkdir -p "$dir"
install -m 0755 "$tmp/ccam" "$dir/ccam"
# Older versions linked a ccs alias to ccam; drop it, but never a foreign ccs.
if [ -L "$dir/ccs" ] && [ "$(readlink "$dir/ccs")" = ccam ]; then
  rm -f "$dir/ccs"
fi
say "✓ $("$dir/ccam" version) → $dir/ccam"

case ":$PATH:" in
  *":$dir:"*) ;;
  *)
    say ""
    say "Note: $dir is not in PATH. Add it to ~/.bashrc or ~/.zshrc:"
    say "  export PATH=\"$dir:\$PATH\""
    ;;
esac
if ! command -v claude >/dev/null 2>&1; then
  say ""
  say "Note: Claude Code (claude) is not installed. Install it first: https://docs.claude.com/claude-code"
fi
say ""
say "Next: ccam add  (add your first account)"
