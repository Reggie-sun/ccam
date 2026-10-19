---
title: MCP 서버, 설정, 대화 내역 공유
description: ccam이 로그인 정보는 계정마다 분리해 둔 채 MCP 서버, settings.json, CLAUDE.md, 에이전트, 스킬, 대화 내역을 Claude Code 계정 간에 공유하는 방법을 설명합니다.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

로그인 정보는 계정마다 따로 유지됩니다. 그 밖의 항목은 하나씩 골라 공유할 수 있습니다. ccam은 계정의 파일이나 폴더를 공유 그룹 폴더에 링크하므로, 한 계정에서 바꾼 내용이 같은 그룹의 모든 계정에 반영됩니다.

## 공유할 수 있는 항목

| 항목 | 옵션 | 공유되는 내용 |
| --- | --- | --- |
| MCP 서버 | `--mcp` | 그룹의 `mcp.json`(시작할 때 Claude Code에 전달) |
| 설정 | `--settings` | `settings.json`: 권한, 훅, 모델, 상태 표시줄 |
| 전역 메모리 | `--memory` | `CLAUDE.md` |
| 확장 기능 | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| 대화 내역 | `--history` | `projects/`와 `history.jsonl` |

새 계정은 대화 내역을 제외한 모든 항목을 공유합니다.

로그인 파일인 `.credentials.json`과 `.claude.json`은 절대 공유되지 않습니다. ccam은 정해진 목록에 있는 경로만 링크하는데, 이 두 파일은 그 목록에 없습니다.

## 그룹

`on`은 `default` 그룹을 뜻합니다. 그룹 이름을 직접 지정하면 일부 계정끼리만 항목을 공유할 수 있습니다.

```sh
ccam share work --history team
ccam share side --history team
```

이제 `work`와 `side`는 서로 대화를 공유하지만, 다른 계정과는 공유하지 않습니다.

## 공유 켜기

계정에서 어떤 항목의 공유를 시작하면 ccam은 다음 순서로 진행합니다.

1. 옮기거나 복사하거나 링크할 파일을 모두 보여 주는 계획을 출력하고, 시작하기 전에 확인을 받습니다.
2. 그룹이 아직 비어 있으면 계정의 파일을 그룹 폴더로 옮기고, 그렇지 않으면 그룹에 없는 파일만 추가합니다.
3. 같은 파일인데 내용이 다르면 그룹의 버전을 유지할지, 내 버전을 쓸지, JSON을 병합할지 묻습니다.
4. 원본은 계정 폴더 안의 `.ccam-backup/<time>/`에 보관합니다.
5. 원본이 있던 자리에 그룹을 가리키는 링크를 둡니다.

대화 파일 이름은 세션 ID를 따르므로, 여러 계정의 대화 내역을 합쳐도 충돌하지 않습니다.

## 공유 끄기

```sh
ccam share work --history off
```

링크가 제거되고, 계정은 그룹의 현재 내용을 복사한 자체 사본을 갖게 됩니다. 빈 상태로 새로 시작하려면 `--fresh`를 붙이세요. 그룹의 내용은 다른 계정을 위해 그대로 남습니다.

## MCP 서버

Claude Code는 MCP 서버를 로그인 정보와 함께 `.claude.json`에 저장하므로, 이 파일은 공유할 수 없습니다. ccam은 공유 서버를 그룹의 `mcp.json`에 두고, `--mcp-config=<group>/mcp.json` 옵션으로 Claude Code를 시작합니다. 직접 지정한 `--mcp-config` 옵션도 그 위에 그대로 적용됩니다.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

세션 안에서 `claude mcp add`로 추가한 서버는 해당 계정에만 저장됩니다. 다음에 그 계정을 시작하면 ccam이 서버 이름과 함께, 그 서버를 그룹으로 옮기는 명령을 알려 줍니다.

```sh
ccam mcp import work gh
```

OAuth로 로그인하는 MCP 서버는 토큰을 계정의 로그인 정보와 함께 저장하므로, 계정마다 한 번씩 인증해야 합니다.

서버를 `claude mcp list`와 `/mcp`에서 일반 사용자 서버로 보고 싶다면 `~/.config/ccam/config.toml`에 `mcp_strategy = "sync"`를 설정하세요. 그러면 ccam이 시작 전에 각 계정의 `.claude.json`에 서버를 기록합니다. 이때 Claude Code의 파일 잠금을 그대로 사용하고 백업도 남깁니다.

## 같은 폴더에서 두 계정 쓰기

대화 내역을 공유하는 계정들은 같은 파일에 기록합니다. 같은 그룹의 다른 계정이 이미 실행 중인 폴더에서 Claude Code를 하나 더 시작하면, ccam이 경고를 보여 주고 계속할지 묻습니다.

## 문제가 생겼을 때

```sh
ccam doctor
ccam doctor --fix
```

예를 들어 백업을 직접 복원해서 공유 경로가 일반 파일로 바뀌었다면, ccam은 Claude Code를 시작하기 전에 멈추고 `ccam doctor --fix`를 실행하라고 안내합니다. 이 명령은 그 파일을 그룹에 병합하고 백업을 남깁니다.
