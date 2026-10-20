---
title: 모든 계정의 백그라운드 세션
description: ccam agents는 모든 계정의 Claude Code 백그라운드 에이전트와 열린 세션을 표 하나로 보여 주고, ccam attach는 ID로 세션 하나를 엽니다.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code는 백그라운드 세션(`claude --bg`, 에이전트 보기)을 계정마다 따로 관리합니다. `claude agents`는 자신이 실행 중인 계정의 세션만 보여 줍니다. `ccam agents`는 모든 계정에 한꺼번에 물어 세션을 함께 나열하고, `ccam attach <id>`는 백그라운드 세션을 그 세션을 시작한 계정으로 엽니다.

## 사용법

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State`는 백그라운드 세션의 상태(`working`, `blocked`, `done`, `failed`, `stopped`)와 그 프로세스가 하는 일(`busy`, `waiting`, `idle`)을 나타냅니다. 열려 있는 대화형 세션은 `open`으로 표시됩니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--all` | 완료, 실패, 중지된 백그라운드 세션도 나열합니다 | 꺼짐 |
| `--json` | 세션을 JSON으로 출력하며, 각 세션에 소속 계정(`account`)이 붙습니다 | 꺼짐 |

`ccam attach <id> [claude args]`는 ID 뒤의 인자를 `claude attach`에 넘깁니다. `--dry-run`을 주면 명령을 출력합니다.

## 동작 방식

1. ccam은 모든 계정에 대해 그 계정의 `CLAUDE_CONFIG_DIR`로 `claude agents --json`을 실행합니다(`--all`을 주면 함께 넘깁니다). 모든 계정을 동시에 실행하며, 계정마다 제한 시간은 15초입니다. 이 명령은 계정 자체의 세션 목록을 읽을 뿐 Claude Code의 백그라운드 서비스를 시작하지 않습니다.
2. 각 항목에 계정 이름을 붙이고, 시작 시각 기준으로 최신순으로 정렬합니다. 응답하지 않는 계정은 표준 오류에 이름을 알리고 목록에서 빼며, 이때 ccam은 종료 코드 1로 끝나므로 스크립트가 실패를 빈 목록으로 착각하지 않습니다.
3. `ccam attach`는 각 계정 폴더에서 Claude Code가 백그라운드 세션을 보관하는 `jobs/<id>`를 찾고, `ccam <name>`이 Claude Code를 시작하는 것과 같은 방식으로 그 계정으로 `claude attach <id>`를 실행합니다.

## 안전장치와 예외 상황

- 백그라운드 세션 ID는 한 계정에 속합니다. 다른 계정으로 실행한 `claude attach`는 그 세션을 찾지 못하며, `ccam attach`가 바로 이 문제를 해결합니다.
- 두 계정에 ID 앞부분이 같은 세션이 있을 수 있습니다. 그러면 `ccam attach`가 두 세션을 모두 보여 주고 ID를 더 길게 입력하라고 합니다.
- ccam은 공유 MCP 서버를 `--mcp-config=<file>` 인자 하나로 Claude Code에 넘기며, `claude agents`와 `claude attach`도 이 형식을 받습니다. 백그라운드 세션은 다시 시작할 때도 그 파일을 계속 씁니다.
- API 키 계정의 백그라운드 서비스는 자신을 시작한 셸에서 키를 가져옵니다. 이런 계정의 첫 백그라운드 세션은 ccam으로 시작해(`ccam gw --bg …`) 서비스가 키를 갖도록 하세요.

## 관련 문서

- [Claude Code 시작하기](/features/launch)
- [스크립트용 계정 고르기](/features/pick)
- [대시보드](/features/dashboard)
