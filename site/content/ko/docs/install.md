---
title: ccam 설치
description: Linux나 macOS에서 명령 하나로, .deb 패키지로, 또는 소스에서 ccam을 설치합니다. Claude Code를 먼저 설치해야 합니다.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Linux나 macOS의 터미널에서 다음을 실행합니다.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

스크립트는 시스템에 맞는 릴리스를 내려받아 SHA-256 체크섬을 확인하고, `~/.local/bin`에 명령 두 개를 넣습니다. `ccam`과 그 짧은 별칭인 `ccs`입니다. 그 폴더가 `PATH`에 없으면 `~/.bashrc`나 `~/.zshrc`에 추가할 줄을 출력합니다.

ccam이 Claude Code를 대신 실행하므로, [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)를 먼저 설치하고 `claude --version`이 동작하는지 확인하세요.

## 지원하는 시스템

- x86-64와 ARM64의 Linux: Ubuntu, Debian, Fedora, Arch, Alpine 및 그 밖의 배포판. 바이너리는 정적이며 의존성이 없습니다.
- Apple 실리콘과 Intel의 macOS.
- Windows는 직접 지원하지 않습니다. WSL 안에서는 다른 Linux 시스템에서와 똑같이 ccam이 동작합니다.

## 다른 설치 방법

**Debian과 Ubuntu 패키지.** [릴리스 페이지](https://github.com/jhyunwoo/ccam/releases)에서 `ccam_<version>_linux_amd64.deb`를 내려받아 설치합니다.

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

패키지는 `/usr/bin/ccam`과 `/usr/bin/ccs` 링크를 설치합니다.

**소스에서.** Go 1.26 이상이 필요합니다.

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**특정 버전.** `CCAM_VERSION`을 지정합니다.

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## 셸 자동완성

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

자동완성은 `ccam`과 `ccs` 모두에 적용되며 계정 이름도 완성합니다.

## 설치 확인

```sh
ccam version
ccam doctor
```

`ccam doctor`는 Claude Code를 찾을 수 있는지, 그리고 모든 계정 폴더, 링크, 로그인이 정상인지 확인합니다. 이어서 [빠른 시작](/docs/quick-start)을 보세요.

## 삭제

명령 두 개와 계정 단축 명령을 지우고, 계정이 더 이상 필요 없다면 데이터 폴더도 지웁니다.

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam`을 지우면 그 안에 저장된 로그인과 대화 내역도 함께 지워집니다. [가져온](/docs/commands) `~/.claude` 폴더는 제자리에 그대로 남습니다.
