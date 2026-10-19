---
title: 명령어 레퍼런스
description: ccam add, ccam <name>부터 ccam usage, ccam next, ccam doctor까지 ccam의 모든 명령과 옵션을 정리했습니다.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 8ac012f7f7cf
---

명령 앞에 `--dry-run`을 붙이면 실제로 실행하지 않고 무엇을 할지만 보여 줍니다(예: `ccam --dry-run work`). Claude Code를 시작할 때(`ccam <name>`, `ccam run`, `ccam next`)와 `share`, `migrate`, `shortcut`, `update`, `doctor --fix`에서 쓸 수 있습니다. 무언가를 바꾸는 다른 명령은 실행하지 않고 이 옵션을 거부합니다.

<!-- commands -->

## Claude Code 시작하기

`ccam <name|number> [arguments]`는 해당 계정으로 Claude Code를 시작합니다. ccam은 이름 뒤의 인자를 해석하지 않고 같은 순서 그대로 `claude`에 넘기며, 표준 입력, 종료 코드, 시그널도 `claude`를 직접 실행할 때와 똑같이 동작합니다. 이름 뒤의 `--`는 넣어도 되고 생략해도 됩니다.

`ccam -p "summarise"`처럼 계정 이름 없이 인자만 넘겨 `ccam`을 실행하면 기본 계정을 사용합니다. 기본 계정은 `ccam use <name>`으로 지정합니다.

`ccwork` 같은 단축 명령도 같은 방식으로 해당 계정을 시작합니다. `ccwork --resume`은 `ccam work --resume`과 똑같이 동작합니다. 자세한 내용은 [단축 명령](/docs/commands#shortcuts)을 참고하세요.

## 계정 추가와 삭제

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

폴더를 지정하지 않고 `ccam import`를 실행하면 `~/.claude`와 `~/.claude-*`를 찾아 폴더마다 하나씩 물어봅니다. 가져온 폴더는 원래 위치에서 그대로 등록되며, ccam은 로그인 정보를 절대 옮기거나 복사하지 않습니다.

`ccam rename`, `ccam remove`, `ccam share`는 어떤 계정으로 Claude Code가 실행 중이면 그 계정을 바꾸지 않고 거부합니다. 먼저 세션을 종료하세요.

`ccam list --json`은 모든 계정의 슬롯 번호, 단축 명령, 이메일, 폴더를 출력합니다. `share` 객체는 각 항목(`mcp`, `settings`, `memory`, `extensions`, `history`)을 그룹 이름이나 `local`에 대응시킵니다.

## 단축 명령

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

단축 명령은 직접 이름을 붙여 한 계정으로 Claude Code를 시작하는 명령입니다. `ccam work`와 마찬가지로 모든 인자를 그대로 Claude Code에 넘기므로 `ccwork --resume`이나 `ccwork -p "hi"`도 예상대로 동작합니다. 단축 명령은 `~/.local/bin`에 있는 ccam 링크입니다.

단축 명령은 필수가 아닙니다. `ccam add` 설정 마법사가 만들지 여부를 묻고, `ccam add <name> --shortcut <command>`는 묻지 않고 바로 만듭니다. `ccam list`는 계정마다 단축 명령을 보여 주며, 없으면 `-`로 표시합니다.

이름에는 a-z, 0-9, `.`, `-`, `_`만 쓸 수 있습니다. ccam 명령 이름, 다른 계정의 단축 명령, `PATH`에 있는 다른 프로그램의 이름은 쓸 수 없습니다.

`ccam rename`은 단축 명령 이름을 그대로 두고, `ccam remove`는 계정과 함께 단축 명령도 삭제합니다. `ccam shims install`은 모든 계정의 링크를 다시 만듭니다. ccam을 다른 위치로 옮긴 뒤에 쓰면 됩니다.

## 공유

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

항목마다 `on`(`default` 그룹), `off`, 또는 그룹 이름을 지정할 수 있습니다. 자세한 내용은 [공유](/docs/sharing)를 참고하세요.

## 기존 설정 복사하기

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`는 Claude Code 폴더에 있는 `settings.json`, `CLAUDE.md`, 확장 기능, 사용자 MCP 서버를 공유 그룹으로 복사합니다. 그러면 그 항목을 공유하는 모든 계정이 같은 내용을 쓰게 됩니다. 원본 폴더는 읽기만 합니다. `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` 옵션을 주면 지정한 항목만 복사하고, 옵션이 없으면 대화 내역을 제외한 모든 항목을 복사합니다. `-g` 또는 `--group`으로 그룹을 고르고, `-y` 또는 `--yes`를 주면 묻지 않고 복사하며 내용이 다른 파일은 그룹의 사본을 유지합니다. 자세한 내용은 [기존 설정 복사하기](/docs/sharing#copy-an-existing-setup)를 참고하세요.

## MCP 서버

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>`는 그 계정에만 있는 서버를 계정의 그룹으로 옮깁니다. 그룹에 같은 이름으로 정의가 다른 서버가 이미 있으면 건너뜁니다. `--force`를 붙이면 그룹의 서버를 바꾸며, ccam이 먼저 기존 것을 백업합니다.

## 사용량 확인과 계정 전환

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## 유지 관리

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
```

`ccam update`는 ccam을 교체하기 전에 릴리스의 SHA-256 체크섬을 확인하고, `--version v0.2.0`을 주면 특정 릴리스를 설치합니다. `.deb` 패키지로 설치한 ccam은 apt로 업데이트합니다. 자세한 내용은 [업데이트](/docs/install#update)를 참고하세요.

## 파일

| 경로 | 내용 |
| --- | --- |
| `~/.config/ccam/config.toml` | 계정 목록과 각 계정의 공유 항목 |
| `~/.local/share/ccam/profiles/<name>/` | 계정별 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 공유된 설정, MCP 서버, 확장 기능, 대화 내역 |
| `~/.local/state/ccam/` | 사용량 캐시와 마지막 사용 시각 |
| `~/.local/bin/<shortcut>` | `ccwork` 같은 단축 명령(ccam 링크) |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` 변수로 이 폴더들의 위치를 바꿀 수 있습니다.
