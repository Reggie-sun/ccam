---
title: 명령어 레퍼런스
metaTitle: "ccam 명령어 레퍼런스"
description: ccam add, ccam <name>부터 ccam usage, ccam next, ccam doctor까지 ccam의 모든 명령과 옵션을 정리했습니다.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: "116c5b34ac9a"
---

명령 앞에 `--dry-run`을 붙이면 실제로 실행하지 않고 무엇을 할지만 보여 줍니다(예: `ccam --dry-run work`). Claude Code를 시작할 때(`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`)와 `group`, `migrate`, `shortcut`, `update`, `doctor --fix`에서 쓸 수 있습니다. 무언가를 바꾸는 다른 명령은 실행하지 않고 이 옵션을 거부합니다.

<!-- commands -->

## Claude Code 시작하기

`ccam <name|number> [arguments]`는 해당 계정으로 Claude Code를 시작합니다. ccam은 이름 뒤의 인자를 해석하지 않고 같은 순서 그대로 `claude`에 넘기며, 표준 입력, 종료 코드, 시그널도 `claude`를 직접 실행할 때와 똑같이 동작합니다. 이름 뒤의 `--`는 넣어도 되고 생략해도 됩니다.

`ccam -p "summarise"`처럼 계정 이름 없이 인자만 넘겨 `ccam`을 실행하면 현재 폴더에 지정된 계정을, 없으면 기본 계정을 사용합니다. 각각 `ccam folder <name> [folder]`와 `ccam use <name>`으로 지정합니다.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

계정 이름 앞의 `--fallback`은 `-p`/`--print` 실행용입니다. 계정이 사용량 한도에 도달하면 ccam은 로그인된 계정 중 여유가 가장 많은 계정으로 같은 명령을 다시 실행하고, 끝까지 실행된 쪽의 출력만 표시합니다. [폴더별 계정](/features/folder)과 [한도에서 다시 실행](/features/next#run-again-at-the-limit)을 참고하세요.

`ccwork` 같은 단축 명령도 같은 방식으로 해당 계정을 시작합니다. `ccwork --resume`은 `ccam work --resume`과 똑같이 동작합니다. 자세한 내용은 [단축 명령](/docs/commands#shortcuts)을 참고하세요.

## 계정 추가와 삭제

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

폴더를 지정하지 않고 `ccam import`를 실행하면 `~/.claude`와 `~/.claude-*`를 찾아 폴더마다 하나씩 물어봅니다. 가져온 폴더는 원래 위치에서 그대로 등록되며, ccam은 로그인 정보를 절대 옮기거나 복사하지 않습니다.

`ccam add --api`는 Claude 로그인 대신 API 키를 쓰는 계정을 추가합니다. ccam은 실행할 때마다 키를 읽으며 절대 저장하지 않습니다. [API 키 계정](/features/api)을 참고하세요.

`ccam rename`, `ccam remove`, 그리고 `ccam group` 명령들은 어떤 계정으로 Claude Code가 실행 중이면 그 계정을 바꾸지 않고 거부합니다. 먼저 세션을 종료하세요.

`ccam list --json`은 모든 계정의 슬롯 번호, 단축 명령, 이메일, 폴더를 출력합니다. `groups` 배열은 계정이 속한 그룹을 나열하고, `share` 객체는 각 항목(`mcp`, `settings`, `memory`, `extensions`, `history`)을 그 항목을 받는 그룹이나 `local`에 대응시킵니다.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

그룹은 속한 계정끼리 항목을 공유합니다. 한 계정이 여러 그룹에 속할 수 있지만, 각 항목은 한 그룹에서만 받습니다. 자세한 내용은 [그룹으로 공유하기](/features/share)를 참고하세요.

## 기존 설정 복사하기

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`는 Claude Code 폴더에 있는 `settings.json`, `CLAUDE.md`, 확장 기능, 사용자 MCP 서버를 공유 그룹으로 복사합니다. 그러면 그 항목을 공유하는 모든 계정이 같은 내용을 쓰게 됩니다. 원본 폴더는 읽기만 합니다. `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` 옵션을 주면 지정한 항목만 복사합니다. 옵션이 없으면 터미널에서는 체크리스트가 나타나 직접 고를 수 있으며, 대화 내역을 제외한 모든 항목이 선택되어 있습니다. 터미널이 아니면 대화 내역을 제외한 모든 항목을 복사합니다. `-g` 또는 `--group`으로 그룹을 고르고, `-y` 또는 `--yes`를 주면 묻지 않고 복사하며 내용이 다른 파일은 그룹의 사본을 유지합니다. 또한 MCP 서버에 저장된 로그인 정보를 각 계정 자체의 로그인 파일로 복사합니다. 자세한 내용은 [기존 설정 복사하기](/docs/sharing#copy-an-existing-setup)를 참고하세요.

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
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback`, `ccam pick`은 같은 방식으로 고릅니다. `auto_max_usage`(`config.toml`, 기본값 100) 미만인 계정을 먼저, `auto_strategy`(`headroom`, `expiring` 또는 `balance`)에 따라 고르며, `--model`의 주간 구간도 따지고, 100%인 계정은 절대 고르지 않고, 시작한 계정과 같은 조직 안에서만 고릅니다(더 넓게 찾으려면 `--cross-org` 또는 `cross_org = true`). 고를 수 있는 계정이 없으면 종료 코드 75로 끝나며, `--wait`를 주면 첫 번째 계정이 초기화될 때까지 기다립니다. `ccam next --continue`는 대화 내역을 공유하지 않는 계정 사이에서도 대화를 넘깁니다. 이때는 Claude Code가 대화를 새 계정으로 복사합니다. Claude Code 안에서는 플러그인의 `/ccam:switch [account]`가 대화를 넘기며, `ccam hook`은 플러그인의 훅이 실행하는 명령입니다. [다음 계정으로 전환](/features/next), [스크립트용 계정 고르기](/features/pick), [모든 계정의 백그라운드 세션](/features/agents)을 참고하세요.

## 유지 관리

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update`는 ccam을 교체하기 전에 릴리스의 SHA-256 체크섬을 확인하고, `--version v0.2.0`을 주면 특정 릴리스를 설치합니다. `.deb` 패키지로 설치한 ccam은 apt로 업데이트합니다. 자세한 내용은 [업데이트](/docs/install#update)를 참고하세요.

## 파일

| 경로 | 내용 |
| --- | --- |
| `~/.config/ccam/config.toml` | 계정, 그룹, 설정(`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | 계정별 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 공유된 설정, MCP 서버, 확장 기능, 대화 내역 |
| `~/.local/state/ccam/` | 사용량 캐시, 마지막 사용 시각, `/ccam:switch`와 한도 훅이 남기는 메모(`handoff/`) |
| `~/.local/bin/<shortcut>` | `ccwork` 같은 단축 명령(ccam 링크) |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` 변수로 이 폴더들의 위치를 바꿀 수 있습니다. Windows에서는 `~`가 `%USERPROFILE%`이고, 단축 명령은 `<shortcut>.exe`입니다.
