---
title: 확장 살펴보기
description: ccam ext는 계정마다 불러오는 스킬, 명령, 플러그인, MCP 서버, 훅과 각각의 출처, 그리고 망가진 부분을 보여 줍니다.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

각 칸은 그 계정에서 그 항목의 출처를 나타냅니다. 공유 그룹(`default`), 계정 자체 폴더에 있으면 `account`, 계정에 없으면 `-`입니다. `!`는 문제가 있는 항목을 표시하며, 문제는 표 아래에 설명합니다.

## 읽는 곳

| 종류 | 위치 | 표시되는 출처 |
| --- | --- | --- |
| 스킬 | `skills/` 안에서 `SKILL.md`가 있는 각 폴더 | `skills/`가 링크하는 그룹, 또는 `account` |
| 명령, 에이전트 | `commands/`와 `agents/` 아래의 `.md` 파일, 하위 폴더는 `folder/name` 형태 | 그룹, 또는 `account` |
| 플러그인 | `plugins/installed_plugins.json` | 그룹, 또는 `account` |
| MCP 서버 | 그룹의 `mcp.json`(MCP를 공유할 때)과 계정의 `.claude.json` | 그룹, 또는 `account` |
| 훅 | `settings.json`의 `hooks`, 이벤트별 개수 | 그룹, 또는 `account` |
| 이 프로젝트 | 현재 폴더의 `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json`, `.mcp.json` | 모든 계정에 대해 한 번만 표시 |

## 알려 주는 문제

- **깨진 링크**: `skills/`, `commands/`, `agents/` 안에서 대상이 사라진 링크, 또는 없는 위치를 가리키는 공유 폴더(`skills`, `commands`, `agents`, `plugins`)입니다(`ccam doctor --fix`로 복구).
- **`SKILL.md`가 없는 스킬 폴더**: Claude Code가 건너뜁니다.
- **두 곳에서 다르게 정의된 MCP 서버**: 그룹의 `mcp.json`과 계정의 `.claude.json`, 또는 그중 하나와 프로젝트의 `.mcp.json`입니다. 두 곳의 설정이 같으면(`mcp_strategy = "sync"`가 이렇게 씁니다) 문제가 아닙니다. 계정 자체 서버를 그룹으로 옮기려면 `ccam mcp import <account> <server>`를 쓰세요.
- **올바른 JSON이 아닌 `installed_plugins.json`.**

## 안전장치와 예외 상황

- `ccam ext`는 아무것도 바꾸지 않고 네트워크 요청도 하지 않습니다. 서버 이름과 훅 개수만 읽으며, MCP 정의에 있는 `env`나 헤더의 값은 절대 읽지 않습니다.
- 디스크에 있는 그대로를 보여 줍니다. 설치했지만 설정에서 꺼 둔 플러그인이나 관리형 설정이 막은 서버도 표시됩니다.
- 프로젝트 폴더에서 실행하면 그 프로젝트의 `.claude` 폴더와 `.mcp.json`도 포함됩니다.

## 관련 문서

- [계정 간 공유](/features/share)
- [MCP 서버 공유](/features/mcp)
- [점검과 복구](/features/doctor)
