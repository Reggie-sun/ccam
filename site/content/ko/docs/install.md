---
title: ccam 설치
description: Linux나 macOS에서 명령 하나로, 또는 .deb 패키지나 소스로 ccam을 설치합니다. Claude Code가 먼저 설치되어 있어야 합니다.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Linux나 macOS 터미널에서 다음 명령을 실행하세요.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

스크립트는 시스템에 맞는 릴리스를 내려받아 SHA-256 체크섬을 확인한 뒤 `ccam`을 `~/.local/bin`에 설치합니다. 나중에 `ccam add`가 `claude1`, `claude2`… 단축 명령을 만드는 곳도 이 폴더입니다. 이 폴더가 `PATH`에 없으면 스크립트가 `~/.bashrc`나 `~/.zshrc`에 추가할 줄을 알려 줍니다.

ccam은 Claude Code를 대신 실행해 주는 도구이므로, 먼저 [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)를 설치하고 `claude --version`이 제대로 실행되는지 확인하세요.

## 지원 시스템

- Linux(x86-64, ARM64): Ubuntu, Debian, Fedora, Arch, Alpine 등 여러 배포판. 정적 바이너리라 별도 의존성이 없습니다.
- macOS(Apple 실리콘, Intel).
- Windows는 직접 지원하지 않습니다. WSL에서는 일반 Linux 시스템과 똑같이 동작합니다.

## 다른 설치 방법

**Debian, Ubuntu 패키지.** [릴리스 페이지](https://github.com/jhyunwoo/ccam/releases)에서 `ccam_<version>_linux_amd64.deb`를 내려받아 설치합니다.

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

패키지로 설치하면 `/usr/bin/ccam`에 설치됩니다.

**소스에서 설치.** Go 1.26 이상이 필요합니다.

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**특정 버전 설치.** `CCAM_VERSION`을 지정합니다.

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## 셸 자동완성

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

자동완성은 계정 이름을 포함해 `ccam`의 모든 명령과 옵션을 지원합니다.

## 설치 확인

```sh
ccam version
ccam doctor
```

`ccam doctor`는 Claude Code를 찾을 수 있는지, 모든 계정 폴더와 링크, 로그인이 정상인지 확인합니다. 다음으로 [빠른 시작](/docs/quick-start)을 읽어 보세요.

## 제거

ccam 명령과 계정 단축 명령을 지우고, 계정이 더 이상 필요 없다면 데이터 폴더도 지웁니다.

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam`을 지우면 그 안에 저장된 로그인 정보와 대화 내역도 함께 지워집니다. [가져온](/docs/commands) `~/.claude` 폴더는 원래 위치에 그대로 남습니다.
