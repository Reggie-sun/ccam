#!/bin/sh
# Installs ccam (and the ccs alias) into ~/.local/bin.
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
die() { printf 'ccam 설치 실패: %s\n' "$*" >&2; exit 1; }

[ "$(uname -s)" = Linux ] || die "리눅스만 지원합니다"
case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  *) die "지원하지 않는 CPU: $(uname -m)" ;;
esac

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
asset="ccam_linux_${arch}.tar.gz"

fetch() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$1" -o "$2"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$1" -O "$2"
  else
    die "curl 이나 wget 이 필요합니다"
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
  say "→ $base/$asset 내려받는 중"
  fetch "$base/$asset" "$tmp/$asset" || die "내려받기 실패"
  fetch "$base/checksums.txt" "$tmp/checksums.txt" || die "checksums.txt 내려받기 실패"
  (cd "$tmp" && grep " $asset\$" checksums.txt | sha256sum -c --quiet -) || die "체크섬이 맞지 않습니다"
fi

tar -xzf "$tmp/$asset" -C "$tmp" ccam || die "압축을 풀 수 없습니다"
mkdir -p "$dir"
install -m 0755 "$tmp/ccam" "$dir/ccam"
ln -sf ccam "$dir/ccs"
say "✓ $("$dir/ccam" version) → $dir/ccam (짧게: ccs)"

case ":$PATH:" in
  *":$dir:"*) ;;
  *)
    say ""
    say "참고: $dir 가 PATH 에 없습니다. ~/.bashrc 나 ~/.zshrc 에 추가하세요:"
    say "  export PATH=\"$dir:\$PATH\""
    ;;
esac
if ! command -v claude >/dev/null 2>&1; then
  say ""
  say "참고: Claude Code(claude)가 없습니다. 먼저 설치하세요: https://docs.claude.com/claude-code"
fi
say ""
say "다음: ccam add  (첫 계정 추가)"
