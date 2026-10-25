---
title: ccam 설치
metaTitle: "Linux, macOS, Windows에 ccam 설치하기"
description: "Linux, macOS, Windows에서 명령 하나, .deb 패키지 또는 소스로 ccam을 설치하고 ccam update로 업데이트합니다. Claude Code를 먼저 설치해야 합니다."
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "a7f9dac24f79"
---

Linux나 macOS 터미널에서 다음 명령을 실행하세요.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

스크립트는 시스템에 맞는 릴리스를 내려받아 SHA-256 체크섬을 확인한 뒤 `ccam`을 `~/.local/bin`에 설치합니다. 나중에 `ccwork`처럼 직접 고른 단축 명령을 ccam이 만드는 곳도 이 폴더입니다. 이 폴더가 `PATH`에 없으면 스크립트가 `~/.bashrc`나 `~/.zshrc`에 추가할 줄을 알려 줍니다.

Windows에서는 PowerShell에서 다음 명령을 실행하세요.

```powershell
irm https://ccam.dev/install.ps1 | iex
```

이 스크립트는 `ccam_windows_<arch>.zip`을 내려받아 SHA-256 체크섬을 확인하고, Claude Code 설치 프로그램이 `claude.exe`를 두는 폴더인 `%USERPROFILE%\.local\bin`에 `ccam.exe`를 넣은 뒤, 이 폴더가 사용자 `PATH`에 없으면 추가합니다. 이때 `%USERPROFILE%\bin` 같은 항목은 그대로 둡니다. 설치가 끝나면 새 터미널을 여세요. 스크립트를 다시 실행하면 ccam을 교체하며, 이전 바이너리로 Claude Code 세션이 아직 실행 중이어도 됩니다(그 바이너리는 `ccam.exe.old`로, 더 이전 것이 아직 쓰이고 있으면 `ccam.exe.old-<number>`로 이름이 바뀝니다). 하드 링크로 된 단축 명령도 새로 만듭니다: `OK 2 shortcuts now run the new ccam`.

ccam은 Claude Code를 대신 실행해 주는 도구이므로, 먼저 [Claude Code](https://docs.claude.com/en/docs/claude-code/overview)를 설치하고 `claude --version`이 제대로 실행되는지 확인하세요.

## 지원 시스템

- Linux(x86-64, ARM64): Ubuntu, Debian, Fedora, Arch, Alpine 등 여러 배포판. 정적 바이너리라 별도 의존성이 없습니다.
- macOS(Apple 실리콘, Intel).
- Windows 10, 11(x86-64, ARM64): PowerShell, Windows 터미널, Git Bash에서 실행합니다. 폴더 공유는 모든 사용자가 쓸 수 있습니다. 파일 공유(`settings.json`, `CLAUDE.md`, `history.jsonl`)에는 개발자 모드(설정 > 시스템 > 개발자용)가 필요합니다. Windows에서는 이 모드이거나 관리자일 때만 파일 링크를 만들 수 있기 때문입니다. [Windows에서 공유](/features/share#on-windows)를 참고하세요. 개발자 모드가 없으면 처음 `ccam add`를 실행할 때 MCP 서버와 확장 기능만 공유하는 `default` 그룹을 만들고, 개발자 모드를 켠 뒤 나머지를 추가하는 방법을 알려 줍니다.
- WSL에서는 일반 Linux 시스템과 똑같이 동작합니다.

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

릴리스 태그(`@latest` 또는 `@vX.Y.Z`)로 만든 빌드는 그 버전을 표시하므로 `ccam update`가 동작합니다. `@main`이나 태그 이후의 로컬 체크아웃으로 만든 빌드는 `dev`로 표시됩니다.

**특정 버전 설치.** `CCAM_VERSION`을 지정합니다.

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## 셸 자동완성

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

자동완성은 계정 이름을 포함해 `ccam`의 모든 명령과 옵션을 지원합니다.

## Claude Code 플러그인

ccam 플러그인은 Claude Code가 직접 ccam을 쓰도록 가르칩니다. 어느 계정에 여유가 있는지 확인하고, `ccam pick`으로 여러 계정에서 병렬 `claude -p` 워커를 시작하고, 스킬이나 MCP 서버가 보이지 않는 이유를 `ccam ext`로 설명합니다. Claude Code 안에서 설치하세요.

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

스킬 하나, 명령 `/ccam:usage`, `/ccam:ext`, `/ccam:switch`, 그리고 훅 두 개가 추가됩니다. 세션이 사용량 한도에 걸려 멈추면 터미널 알림이 오고(Linux와 macOS에서 [`auto_switch`](/features/next#auto-switch)를 켜 두었다면 ccam이 여유 있는 다른 계정으로 세션을 옮깁니다), `/ccam:switch [account]`는 대화를 다른 계정으로 넘깁니다([자세히](/features/next#inside-claude-code)). 훅은 `ccam`을 실행하므로, Claude Code가 보는 `PATH`에 ccam이 있어야 합니다. 스킬은 로그인하거나 파일을 옮기는 명령을 직접 실행하지 않고, 대신 제안만 합니다.

## 설치 확인

```sh
ccam version
ccam doctor
```

`ccam doctor`는 Claude Code를 찾을 수 있는지, 모든 계정 폴더와 링크, 로그인이 정상인지 확인합니다. 다음으로 [빠른 시작](/docs/quick-start)을 읽어 보세요.

## 업데이트

```sh
ccam update --check
ccam update
```

`ccam update --check`는 새 릴리스가 있는지 알려 주기만 합니다. `ccam update`는 시스템에 맞는 릴리스를 내려받아 릴리스의 `checksums.txt`로 SHA-256 체크섬을 확인한 뒤 `ccam` 파일을 한 번에 교체합니다. 체크섬이 맞지 않으면 아무것도 바꾸지 않습니다. `ccwork` 같은 단축 명령은 ccam을 가리키는 링크라서 곧바로 새 버전으로 실행됩니다.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

대시보드(터미널에서 실행한 `ccam`)도 열릴 때마다 새 릴리스를 찾고, 새 릴리스가 있으면 업데이트를 제안합니다. `U`를 누르고 확인하면 됩니다. 다른 명령은 새 릴리스를 찾지 않으며, ccam이 묻지 않고 업데이트하는 일은 없습니다. 이전 버전을 포함해 특정 릴리스를 설치하려면 `--version`을 지정합니다. `ccam --dry-run update`는 아무것도 바꾸지 않고, 무엇을 내려받아 어떤 파일을 교체할지만 보여 줍니다.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

설치 스크립트를 다시 실행해도 `~/.local/bin`의 ccam이 업데이트됩니다. 릴리스 태그로 `go install`해서 설치했다면 `ccam update`도 동작하고, 같은 명령을 다시 실행해도 됩니다.

**Windows.** `ccam update`는 Windows용 `.zip`을 내려받습니다. Windows에서는 실행 중인 프로그램을 교체할 수 없으므로, 기존 바이너리의 이름을 `ccam.exe.old`로 바꿔 두고 다음 업데이트 때 지웁니다. 이전의 `ccam.exe.old`로 실행 중인 세션이 아직 있으면 이름은 `ccam.exe.old-<time>`이 됩니다. 하드 링크로 된 단축 명령은 바로 이어서 새로 만듭니다.

**Debian, Ubuntu 패키지.** apt로 설치한 `/usr/bin/ccam`은 `ccam update`가 교체하지 않고, 대신 아래 방법을 알려 줍니다. [릴리스 페이지](https://github.com/jhyunwoo/ccam/releases)에서 새 `ccam_<version>_linux_amd64.deb`를 내려받아 기존 패키지 위에 설치하세요.

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## 제거

ccam이 아직 설치되어 있을 때 계정 단축 명령부터 지웁니다. `ccam list`의 Shortcut 열에서 단축 명령을 확인하고, 하나씩 `ccam shortcut <account> --remove`로 지우거나 `~/.local/bin`에서 해당 링크를 직접 지우세요. 그다음 ccam 명령을 지우고, 계정이 더 이상 필요 없다면 데이터 폴더도 지웁니다.

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows에서는 PowerShell에서 다음을 실행합니다.

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam`을 지우면 그 안에 저장된 로그인 정보와 대화 내역도 함께 지워집니다. [가져온](/docs/commands) `~/.claude` 폴더는 원래 위치에 그대로 남지만, 그룹에 속해 있는 동안에는 공유하는 파일이 `~/.local/share/ccam/groups/`에 있습니다. 먼저 그 계정에 `ccam remove <name>`(또는 `ccam group leave`)을 실행하세요. 그러면 공유하던 모든 항목의 사본을 받고 그룹에서 빠집니다.
