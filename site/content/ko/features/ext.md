---
title: 확장 살펴보기
description: ccam ext는 계정마다 불러오는 스킬, 명령, 플러그인, MCP 서버, 훅과 각각의 출처, 그리고 망가진 부분을 보여 줍니다.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: "45740e4efbef"
---

`ccam ext`는 "이 스킬(또는 MCP 서버)이 왜 한 계정에서는 되고 다른 계정에서는 안 될까?"라는 질문에 답합니다. 계정마다 Claude Code가 불러오는 확장과 각각의 출처(공유 그룹, 계정 자체 폴더, 지금 있는 프로젝트)를 나열합니다. 표 아래에는 확장이 예상대로 로드되지 못하게 막는 문제를 나열합니다. 읽기만 합니다.

## 사용법

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

각 칸은 그 계정에서 그 항목의 출처를 나타냅니다. 공유 그룹(`default`), 계정 자체 폴더에 있으면 `account`, 계정이 이 폴더에서만 가진 MCP 서버는 `local`, 계정에 없으면 `-`입니다. 설치했지만 켜지 않은 플러그인은 출처 뒤에 `(off)`가 붙습니다. `!`는 문제가 있는 항목을 표시하며, 문제는 표 아래에 설명합니다. 일반 `claude`에만 있는 서버처럼 어느 계정에도 속하지 않는 문제는 이름 없이 맨 끝에 나옵니다. `ccam ext --json`은 같은 보고서(`accounts`, 프로젝트, `problems`)를 출력하며, Claude Code가 네트워크로 연결하는 MCP 서버에는 `remote`가 붙습니다. 계정이 없으면 `accounts: []`인 보고서를 출력합니다.

## 읽는 곳

| 종류 | 위치 | 표시되는 출처 |
| --- | --- | --- |
| 스킬 | `skills/` 안에서 `SKILL.md`가 있는 각 폴더 | `skills/`가 링크하는 그룹, 또는 `account` |
| 명령, 에이전트 | `commands/`와 `agents/` 아래의 `.md` 파일, 하위 폴더는 `folder/name` 형태 | 그룹, 또는 `account` |
| 플러그인 | `plugins/installed_plugins.json`. 계정과 프로젝트 설정의 `enabledPlugins`로 켜고 끕니다. 한 프로젝트에만 설치한 플러그인은 그 프로젝트에서만 나타납니다 | 그룹, 또는 `account`. 켜지 않았으면 `(off)` |
| MCP 서버 | 그룹의 `mcp.json`(MCP를 공유할 때), 계정의 `.claude.json`(계정 자체 서버와 이 폴더에서만 가진 서버), 그리고 계정 플러그인의 서버(`plugin:<plugin>:<server>` 형태) | 그룹, `account` 또는 `local` |
| 훅 | `settings.json`의 `hooks`, 이벤트별 개수 | 그룹, 또는 `account` |
| 이 프로젝트 | 현재 폴더의 `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json`, `.mcp.json` | 한 번만 표시. 그 폴더를 신뢰하지 않은 계정에는 `.mcp.json` 서버에 대한 문제를 표시 |

## 알려 주는 문제

- **깨진 링크**: `skills/`, `commands/`, `agents/` 안에서 대상이 사라진 링크, 또는 없는 위치를 가리키는 공유 폴더(`skills`, `commands`, `agents`, `plugins`)입니다(`ccam doctor --fix`로 복구하며, 없는 그룹 폴더도 다시 만듭니다).
- **`SKILL.md`가 없는 스킬 폴더**: Claude Code가 건너뜁니다.
- **두 곳에서 다르게 정의된 MCP 서버**: 그룹의 `mcp.json`과 계정의 `.claude.json`, 또는 그중 하나와 프로젝트의 `.mcp.json`입니다. 두 곳의 설정이 같으면(`mcp_strategy = "sync"`가 이렇게 씁니다) 문제가 아닙니다. 계정 자체 서버를 그룹으로 옮기려면 `ccam mcp import <account> <server>`를 쓰세요.
- **올바른 JSON이 아닌 `installed_plugins.json`**, `enabledPlugins`에 이름이 없는 플러그인, 폴더가 없는 플러그인.
- **한 계정만 불러오는 서버**: 이 계정의 `.claude.json`에만 있어 그룹의 다른 계정은 불러오지 않는 서버, 그리고 이 폴더에 대한 계정의 local 범위에 있는 서버(`-s user`를 주지 않은 `claude mcp add`가 이렇게 추가합니다)입니다. 각각 그 서버를 공유하는 `ccam mcp import` 명령과 함께 표시됩니다. 일반 `claude`의 `~/.claude.json`에 있지만 어느 계정도 불러오지 않는 서버는 `ccam migrate --mcp` 명령과 함께 표시됩니다.
- **승인을 기다리는 프로젝트 서버**: 계정이 아직 그 폴더를 신뢰하지 않았으므로 Claude Code가 한 번 묻고, 그때까지 프로젝트의 `.mcp.json` 서버는 기다립니다. 문제 설명에 승인할 곳이 나옵니다. 설정 그룹의 계정들 또는 이 계정의 `settings.json`에 있는 `"enabledMcpjsonServers"`입니다. git 저장소의 메인 워크트리에 준 신뢰는 다른 워크트리에서도 인정됩니다.
- **로그인과 커넥터**: 다른 계정은 로그인했지만 이 계정에서는 로그인이 필요한 MCP 서버(로그인은 한 계정에만 유효하니 `/mcp`로 로그인하세요), 그리고 다른 계정은 썼지만 이 계정은 쓰지 않은 claude.ai 커넥터입니다.

## 안전장치와 예외 상황

- `ccam ext`는 아무것도 바꾸지 않고 네트워크 요청도 하지 않습니다. 서버 이름과 훅 개수만 읽으며, MCP 정의에 있는 `env`나 헤더의 값은 절대 읽지 않습니다. 로그인 확인을 위해 각 계정 로그인 정보에 있는 MCP 로그인의 이름만 읽으며(macOS에서는 키체인에서), 토큰은 절대 읽지 않고, 다른 계정도 불러오는 서버를 불러오는 계정에 대해서만 읽습니다.
- 디스크에 있는 그대로를 보여 줍니다. 관리형 설정이 막은 서버도 표시되며, 설정에서 꺼 둔 플러그인은 `(off)`로 표시됩니다.
- 프로젝트 폴더에서 실행하면 그 프로젝트의 `.claude` 폴더와 `.mcp.json`도 포함됩니다.

## 관련 문서

- [계정 간 공유](/features/share)
- [MCP 서버 공유](/features/mcp)
- [점검과 복구](/features/doctor)
