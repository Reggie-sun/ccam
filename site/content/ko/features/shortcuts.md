---
title: 단축 명령
description: ccam shortcut으로 계정에 ccwork 같은 전용 명령을 붙이고, ccam shims install로 ~/.local/bin의 링크를 다시 만듭니다.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

단축 명령은 직접 이름을 붙여 한 계정으로 Claude Code를 시작하는 명령입니다. `ccwork --resume`은 `ccam work --resume`과 똑같이 동작합니다. 단축 명령은 `~/.local/bin`에 있는 ccam 링크이므로 bash, zsh, 스크립트 등 `PATH`에서 프로그램을 실행하는 곳이라면 어디서나 동작합니다. 단축 명령은 필수가 아니니, 자주 시작하는 계정에만 만들면 됩니다.

## 사용법

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam은 변경할 때마다 결과를 알려 줍니다.

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

계정마다 단축 명령은 최대 하나이며, `~/.config/ccam/config.toml`에 `shortcut`으로 저장됩니다. `ccam add` 설정 마법사는 `Create a shortcut command for <name>? [y/N]`를 묻고, 이어서 `Shortcut name [cc<name>]:`를 묻습니다. 마법사 밖에서는 `--shortcut`을 줄 때만 `ccam add`가 단축 명령을 만듭니다. `ccam list`는 계정마다 단축 명령을 보여 주며, 없으면 `-`로 표시합니다.

단축 명령은 `ccam <name>`이 이름 뒤의 인자를 넘기는 것과 같은 방식으로 모든 인자를 Claude Code에 넘깁니다. `--dry-run`과 `--help`도 예외가 아니므로, 실행을 미리 보려면 `ccam --dry-run work`를 쓰세요.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--remove` | 계정의 단축 명령과 링크를 삭제합니다 | 꺼짐 |
| `-h`, `--help` | 명령 도움말 | |
| `--dry-run` | `ccam shortcut`에서 지원합니다. 이름을 검사하고, 저장하거나 링크하지 않은 채 `Would set the shortcut of work to ccwork.` 또는 `Would remove ccwork, the shortcut of work.`를 출력합니다 | 꺼짐 |

`ccam shims install`은 인자를 받지 않으며, `--dry-run`을 주면 `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``를 출력하고 거부합니다.

이름 규칙은 다음 순서로 검사합니다.

- `a-z`, `0-9`, `.`, `-`, `_`로 이루어진 1~64자이며, 문자나 숫자로 시작해야 합니다.
- `ccam`, `claude`, 또는 `run`이나 `usage` 같은 ccam 명령 이름은 쓸 수 없습니다.
- 다른 계정의 단축 명령은 쓸 수 없습니다.
- `~/.local/bin`에 있는 파일 중 ccam 링크가 아닌 파일의 이름은 쓸 수 없습니다.
- `PATH`에 있는 다른 실행 파일의 이름은 쓸 수 없습니다.

명령 이름과 `--remove`를 함께 주면 `give a command or --remove, not both` 오류를 내고 멈춥니다. 계정에 이미 있는 이름을 다시 지정하면 `work: ccwork (unchanged)`를 출력하고, 링크가 없으면 다시 만듭니다.

## 동작 방식

이름을 지정하면 `checkShortcut`이 실행되고, 이 함수가 bin 폴더와 ccam 실행 파일의 실제 경로를 넘겨 `shim.Check`를 호출합니다.

<!-- motion:shortcuts-check -->

이름은 검사를 순서대로 통과해 내려갑니다. 각 검사 옆에는 그 검사에서 거부되는 이름이 있습니다. `Work`는 패턴에 맞지 않고, `run`은 ccam 명령이며, `ccside`는 다른 계정의 것이고, `ccfile`은 `~/.local/bin`에 이미 있는 파일이며, `git`은 `/usr/bin/git`을 가리게 됩니다. 모든 검사를 통과한 이름만 저장되고 링크됩니다.

1. `shim.Check`는 이름을 `^[a-z0-9][a-z0-9._-]{0,63}$`, `config.Reserved`, `cfg.FindShortcut` 순서로 확인합니다. `~/.local/bin/<name>`에는 `os.Lstat`을 호출하고, `PATH`의 나머지 폴더마다 `os.Stat`을 호출해 ccam 링크가 아닌 실행 파일이 있으면 거부합니다.
2. ccam은 `fsutil.AtomicWrite`(임시 파일, `fsync`, 이름 바꾸기)로 `config.toml`을 저장합니다.
3. `shim.Remove`는 이전 단축 명령의 링크를 삭제하지만, ccam 링크일 때만 삭제합니다.
4. `shim.Install`은 설정된 단축 명령 중 링크가 없는 것을 모두 만듭니다. `os.Executable`과 `filepath.EvalSymlinks`로 찾은 ccam 실행 파일을 가리키도록 `os.Symlink`를 호출합니다. 이 함수는 아무것도 삭제하지 않습니다.

`os.Readlink`가 ccam 경로를 반환하거나, 링크와 ccam 경로가 같은 파일로 해석되면 그 링크는 ccam의 것으로 간주합니다. 모든 삭제와 덮어쓰기 검사가 이 확인을 거치므로, ccam은 자신이 만들지 않은 파일을 건드리지 않습니다.

단축 명령을 실행하면 ccam에게 필요한 것은 입력한 이름뿐입니다.

<!-- motion:shortcuts-link -->

셸이 `PATH`에서 `ccwork`를 찾고, 커널이 링크를 따라가 ccam 실행 파일을 시작하며, `argv[0]`은 여전히 `ccwork`입니다. `app.Main`에서 `ccam`이 아닌 이름이면 `config.toml`을 한 번 읽고 `shim.Lookup`을 호출합니다. 이 함수는 `shortcut`이 `filepath.Base(argv[0])`과 같은 계정을 찾습니다. 그다음 `launcher.LaunchWith`가 이미 읽은 설정으로 [`ccam <name>`](/features/launch)과 같은 단계를 밟고, `argv[1:]`을 손대지 않고 `claude`에 넘깁니다.

셸 별칭이나 래퍼 스크립트 대신 심볼릭 링크를 쓰는 이유는 다음과 같습니다. 심볼릭 링크는 모든 셸과 스크립트에서 동작하고, 셸을 추가로 시작하지 않으며(약 1~2 ms 절약), 실행할 때 읽는 파일도 여전히 하나입니다. `scripts/perf-check.sh`는 `hyperfine`으로 단축 명령을 측정해 `ccam <name>`과 같은 10 ms(p95) 한도와 비교합니다.

### Windows에서

단축 명령은 `%USERPROFILE%\.local\bin`의 `ccwork.exe`이며, `shim.Name`이 조회 전에 대소문자에 관계없이 `.exe`를 떼어 내므로 `ccwork`, `ccwork.exe`, `CCWORK.EXE` 모두 `work`를 시작합니다. ccam은 먼저 심볼릭 링크를 만드는데, Windows에서는 개발자 모드이거나 관리자일 때만 허용됩니다. 그 권한이 없으면 하드 링크를 만듭니다. 하드 링크는 같은 `ccam.exe` 파일의 두 번째 이름이며, 권한이 필요 없습니다. 하드 링크는 `ccam.exe`와 같은 파일일 때 ccam의 것으로 인정됩니다. `ccam update`는 `ccam.exe` 자리에 새 파일을 두므로, 이어서 하드 링크로 된 단축 명령을 모두 `shim.Relink`로 새로 만들고 `✓ shortcut ccwork renewed`를 출력합니다. 이름 검사는 다른 `PATH` 폴더에서 `PATHEXT`의 각 확장자(`.exe`, `.cmd` 등)를 붙인 `<name>`을 찾습니다.

## 안전장치와 예외 상황

- **어느 계정의 것도 아닌 링크.** `config.toml`의 `shortcut`을 직접 고친 경우처럼 `~/.local/bin`에 있는 ccam 링크의 이름을 쓰는 계정이 없으면, ccam은 기본 계정으로 대신 실행하지 않습니다. `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` 오류를 내고 멈춥니다. `ccam shortcut`으로 이름을 다시 지정하거나 링크를 삭제하세요.
- **이미 있는 파일.** ccam은 자신이 만들지 않은 파일을 절대 바꾸지 않습니다. `ccam shims install`은 `<path> already exists and is not a ccam link, so it was left alone`를 출력하고 멈춥니다.
- **ccam 실행 파일을 옮긴 경우.** 이전 링크는 더 이상 없는 경로를 가리키므로, ccam은 이 링크를 자신의 것으로 보지 않습니다. 링크를 삭제한 뒤 `ccam shims install`을 실행하세요.
- **이름 변경과 삭제.** `ccam rename`은 단축 명령 이름을 그대로 두며, 링크도 계속 동작합니다. 링크는 ccam을 가리키고 ccam이 `config.toml`에서 이름을 찾기 때문입니다. `ccam remove`는 계정의 링크를 삭제합니다. `ccam shortcut <account> <new>`는 새 링크를 만들기 전에 이전 링크를 삭제합니다.
- **PATH.** ccam이 링크를 만들 때 `~/.local/bin`이 `PATH`에 없으면 `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:`와 함께 추가할 `export PATH=...` 줄을 출력합니다. Windows에서는 대신 이 폴더를 사용자 `PATH`에 추가하는 PowerShell 명령 줄을 출력합니다.
- **나중에 확인하기.** `ccam doctor`는 `missing shortcuts in /home/me/.local/bin: ccwork`를 보고하고, 폴더가 `PATH`에 없으면 경고합니다. `ccam doctor --fix`는 없는 링크를 만듭니다.

## 관련 문서

- [Claude Code 시작하기](/features/launch)
- [계정 관리](/features/accounts)
- [명령어 레퍼런스](/docs/commands#shortcuts)
- [ccam 설치](/docs/install)
