---
title: MCP 서버, 설정, 대화 내역 공유
description: ccam 그룹이 로그인 정보는 계정마다 분리해 둔 채 MCP 서버, settings.json, CLAUDE.md, 스킬, 대화 내역을 Claude Code 계정 간에 공유하는 방법을 설명합니다.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "27e0a0e5abbd"
---

로그인 정보는 계정마다 따로 유지됩니다. 그 밖의 항목은 **그룹**으로 공유할 수 있습니다. 그룹에는 공유할 항목과 속한 계정이 정해져 있고, ccam은 그 계정들의 파일을 그룹 폴더에 링크합니다. 그래서 한 계정에서 바꾼 내용이 그룹의 모든 계정에 반영됩니다.

## 공유할 수 있는 항목

| 항목 | 공유되는 내용 |
| --- | --- |
| `mcp` | MCP 서버: 그룹의 `mcp.json`(시작할 때 Claude Code에 전달) |
| `settings` | `settings.json`: 권한, 훅, 모델, 상태 표시줄 |
| `memory` | 전역 `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | 대화 내역: `projects/`와 `history.jsonl` |

처음 추가한 계정이 `default` 그룹을 만듭니다. 이 그룹은 대화 내역을 제외한 모든 항목을 공유하며, 이후에 추가하는 계정도 이 그룹에 들어갑니다.

로그인 파일인 `.credentials.json`과 `.claude.json`은 절대 공유되지 않습니다. ccam은 정해진 목록에 있는 경로만 링크하는데, 이 두 파일은 그 목록에 없습니다.

이 페이지에서는 일상적인 작업을 다룹니다. 모든 옵션은 기능 페이지에서 설명합니다. [그룹으로 공유하기](/features/share), [MCP 서버 공유](/features/mcp), [기존 설정 복사하기](/features/migrate)를 참고하세요.

## 그룹

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

한 계정이 여러 그룹에 속할 수 있지만, 각 항목은 한 그룹에서만 받습니다. 그래서 `work`와 `side`끼리만 대화 내역을 공유하고 다른 계정과는 공유하지 않으려면, 위처럼 `default`에서는 대화 내역을 빼고 두 계정만의 그룹을 따로 만드세요. `ccam share work`는 계정이 무엇을 어느 그룹으로 공유하는지 보여 줍니다. [대시보드](/features/dashboard)의 Groups 탭에서도 같은 일을 할 수 있습니다.

## 계정이 공유를 시작할 때

계정이 어떤 항목을 공유하기 시작하면(그룹에 들어가거나, 계정의 그룹에 항목이 추가되면) ccam은 다음 순서로 진행합니다.

1. 계정마다 옮기거나 복사하거나 링크할 파일을 모두 보여 주는 계획을 출력하고, 시작하기 전에 확인을 받습니다.
2. 그룹이 아직 비어 있으면 계정의 파일을 그룹 폴더로 옮기고, 그렇지 않으면 그룹에 없는 파일만 추가합니다.
3. 같은 파일인데 내용이 다르면 그룹의 버전을 유지할지, 내 버전을 쓸지, JSON을 병합할지 묻습니다.
4. 원본은 계정 폴더 안의 `.ccam-backup/<time>/`에 보관합니다.
5. 원본이 있던 자리에 그룹을 가리키는 링크를 둡니다.

대화 파일 이름은 세션 ID를 따르므로, 여러 계정의 대화 내역을 합쳐도 충돌하지 않습니다.

Claude Code가 그 계정으로 실행 중일 때는 Claude Code를 어떻게 시작했든(백그라운드 세션 포함) ccam이 계정을 바꾸지 않습니다. 세션이 옮기려는 파일에 계속 쓰기 때문입니다. 먼저 세션을 종료하세요. ccam은 저장하기 직전에 한 번 더 확인하므로, 계획이 확인을 기다리는 동안 시작한 세션도 잡아냅니다.

## 계정이 공유를 멈출 때

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

링크가 제거되고, 계정은 그룹의 현재 내용을 복사한 자체 사본을 갖게 됩니다. 그룹의 MCP 서버는 계정 자체의 서버와 함께 계정의 `.claude.json`에 들어갑니다. 빈 상태로 새로 시작하려면 `--fresh`를 붙이세요. 이때 계획에 더 이상 불러오지 않을 MCP 서버가 표시됩니다. `~/.claude`처럼 가져온 폴더는 항상 사본을 받습니다. 그룹의 내용은 다른 계정을 위해 그대로 남습니다.

## 기존 설정 복사하기

ccam을 쓰기 전에 Claude Code를 써 왔다면 설정이 `~/.claude`에 있습니다. `ccam migrate`는 이 설정을 그룹으로 복사하므로, 그 항목을 공유하는 모든 계정이 같은 설정을 쓰게 됩니다.

```sh
ccam migrate
```

터미널에서는 ccam이 먼저 항목 체크리스트를 보여 줍니다. 대화 내역을 제외한 모든 항목이 선택되어 있으며, Space로 항목 선택을 바꾸고 Enter를 누르면 계획을 볼 수 있습니다. 폴더에 해당 내용이 없는 항목에는 `not found`가 표시됩니다.

복사하는 항목은 다음과 같습니다.

- `settings.json`과 전역 `CLAUDE.md`
- 확장 기능 `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/`
- 사용자 MCP 서버(`~/.claude.json`의 `mcpServers` 항목). 그룹의 `mcp.json`으로 복사됩니다.
- 대화 내역(`projects/`와 `history.jsonl`). `--history`를 붙였을 때만 복사됩니다.

원본 폴더는 읽기만 합니다. 그 안의 어떤 것도 옮기거나 링크하거나 바꾸지 않으므로, `claude`를 그냥 실행해도 이전과 똑같이 동작합니다. 로그인 파일인 `.credentials.json`과 `.claude.json`은 그룹으로 절대 복사하지 않습니다.

OAuth로 로그인하는 MCP 서버는 토큰을 서버 목록이 아니라 로그인 정보에 보관합니다. ccam은 이 로그인 정보를 복사하지 않습니다. 복사한 로그인은 같은 갱신 토큰을 갖고 있어서, 처음 갱신하는 순간 `~/.claude`를 포함한 다른 모든 쪽이 그 서버에서 로그아웃되기 때문입니다. 계획과 요약에는 `/mcp`로 직접 서버에 로그인해야 하는 계정이 표시됩니다.

ccam은 무엇이든 복사하기 전에 계획을 출력하고 확인을 받습니다.

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Group default: personal
To share them:
  ccam group join default work
```

마지막에는 그룹의 계정을 보여 줍니다. 복사한 항목 중 그룹이 아직 공유하지 않는 항목이 있거나 그룹에 들어 있지 않은 계정이 있으면, 이를 바꾸는 `ccam group` 명령도 함께 출력합니다.

그룹에 내용이 다른 같은 파일이 이미 있으면, 공유를 켤 때와 마찬가지로 그룹의 사본을 유지할지 내 사본을 쓸지 묻습니다. `settings.json` 같은 JSON 파일은 병합할 수도 있고, `history.jsonl`은 줄 단위로 병합할 수 있습니다. 그룹에 있던 이전 사본은 그룹 폴더 안의 `.ccam-backup/<time>/`에 보관됩니다.

| 옵션 | 효과 |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | 체크리스트 없이 지정한 항목만 복사합니다. 옵션이 없고 터미널이 아니면 대화 내역을 제외한 모든 항목을 복사합니다 |
| `-g`, `--group <group>` | `default` 대신 이 그룹으로 복사합니다 |
| `-y`, `--yes` | 묻지 않고 복사합니다. 내용이 다른 파일은 그룹의 사본을 유지합니다 |

`ccam --dry-run migrate`처럼 명령 앞에 `--dry-run`을 붙이면 아무것도 복사하지 않고 계획만 보여 줍니다. 다른 폴더에서 복사하려면 폴더를 지정하세요(예: `ccam migrate ~/.claude-work --group work`).

## MCP 서버

Claude Code는 MCP 서버를 로그인 정보와 함께 `.claude.json`에 저장하므로, 이 파일은 공유할 수 없습니다. ccam은 공유 서버를 그룹의 `mcp.json`에 두고, `--mcp-config=<group>/mcp.json` 옵션으로 Claude Code를 시작합니다. 직접 지정한 `--mcp-config` 옵션도 그 위에 그대로 적용됩니다.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

세션 안에서 `claude mcp add`로 추가한 서버는 해당 계정에만, 그리고 그때 있던 폴더에만 저장됩니다(Claude Code의 local 범위이며, `-s user`를 주면 모든 폴더에서 그 계정의 서버가 됩니다). 다음에 그 계정을 시작하면 ccam이 서버 이름과 그 폴더, 그리고 그 서버를 그룹으로 옮기는 명령을 알려 줍니다.

```sh
ccam mcp import work gh
```

그룹에 같은 이름의 다른 서버가 이미 있으면 ccam은 그 서버를 건너뛰고 그 사실을 알려 줍니다. `--force`를 붙이면 그룹의 서버를 바꾸고, 이전 `mcp.json`은 그룹의 `.ccam-backup/` 폴더에 보관합니다.

OAuth로 로그인하는 MCP 서버는 토큰을 계정의 로그인 정보와 함께 저장하므로, 계정마다 한 번씩 `/mcp`로 인증해야 합니다. `ccam ext`와 `ccam doctor`는 다른 계정은 로그인했지만 이 계정은 로그인하지 않은 서버를 알려 줍니다.

서버를 `claude mcp list`와 `/mcp`에서 일반 사용자 서버로 보고 싶다면 `~/.config/ccam/config.toml`의 맨 위, 첫 `[[group]]`, `[[profile]]`, `[[folder]]`보다 위에 `mcp_strategy = "sync"`를 설정하세요. 그러면 ccam이 시작 전에 각 계정의 `.claude.json`에 서버를 기록합니다. 이때 Claude Code의 파일 잠금을 그대로 사용하고 백업도 남깁니다. 계정에 같은 이름의 자체 서버가 있으면 그 서버는 계정 것으로 남고, ccam이 그 사실을 알려 줍니다.

## 같은 폴더에서 두 계정 쓰기

대화 내역을 공유하는 계정들은 같은 파일에 기록합니다. 같은 그룹의 다른 계정이 이미 실행 중인 폴더에서 Claude Code를 하나 더 시작하면, ccam이 경고를 보여 주고 계속할지 묻습니다.

## 문제가 생겼을 때

```sh
ccam doctor
ccam doctor --fix
```

예를 들어 백업을 직접 복원해서 공유 경로가 일반 파일로 바뀌었다면, ccam은 Claude Code를 시작하기 전에 멈추고 `ccam doctor --fix`를 실행하라고 안내합니다. 이 명령은 그 경로를 다시 그룹에 링크합니다. 그룹에 사본이 없으면 그 파일이 그룹의 사본이 되고, 있으면 그룹은 자기 사본을 유지하고 내 파일은 백업으로 가며, `--fix`가 그 백업을 알려 줍니다. 중간에 멈춘 그룹 변경도 `ccam doctor --fix`가 마무리합니다.
