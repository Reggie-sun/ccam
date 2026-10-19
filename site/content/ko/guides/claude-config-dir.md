---
title: CLAUDE_CONFIG_DIR 이해하기
description: Claude Code에서 CLAUDE_CONFIG_DIR이 하는 일, 폴더에 저장되는 내용, macOS에서 로그인 정보가 보관되는 위치, 계정 간에 공유할 수 있는 파일을 설명합니다.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR`은 Claude Code가 `~/.claude` 대신 사용할 폴더를 지정합니다. 로그인 정보, 설정, MCP 서버, 플러그인, 대화 내역 등 계정에 딸린 모든 것이 이 폴더에 저장됩니다. 이 변수를 서로 다른 폴더로 지정하면 계정을 분리할 수 있습니다.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## 폴더 구성

Claude Code 2.1.302 기준입니다.

| 경로 | 내용 | 공유 가능 |
| --- | --- | --- |
| `.credentials.json` | 로그인 토큰(Linux), MCP OAuth 토큰 | 아니요 |
| `.claude.json` | 계정 정보, 사용자 MCP 서버, 프로젝트별 상태 | 아니요(MCP 서버는 별도로 공유) |
| `settings.json` | 권한, 훅, 모델, 상태 표시줄, 플러그인 | 예 |
| `CLAUDE.md` | 전역 지침 | 예 |
| `agents/`, `commands/`, `skills/`, `plugins/` | 확장 기능 | 예 |
| `projects/` | 대화(`.jsonl` 파일), 프로젝트별 메모리 | 예 |
| `history.jsonl` | 프롬프트 기록 | 예 |
| `sessions/`, `shell-snapshots/`, `cache/` 등 | 런타임 상태 | 아니요 |

## 자주 헷갈리는 두 가지

**`CLAUDE_CONFIG_DIR`을 설정하지 않으면 `.claude.json`은 `~/.claude` 안에 있지 않습니다.** 이 경우 Claude Code는 이 파일을 폴더와 같은 위치의 `~/.claude.json`에 둡니다. 그런데 `CLAUDE_CONFIG_DIR=~/.claude`로 설정해 쓰기 시작하면 Claude Code는 `~/.claude/.claude.json`을 찾기 때문에 기존 계정을 찾지 못합니다. `ccam import`는 두 파일을 링크로 연결해 이 문제를 해결합니다.

**macOS에서는 로그인 정보가 폴더가 아니라 키체인에 저장됩니다.** Claude Code는 `Claude Code-credentials`라는 이름으로 저장하고, `CLAUDE_CONFIG_DIR`이 설정되어 있으면 이름 뒤에 폴더 경로의 짧은 해시를 붙입니다. 따라서 폴더마다 키체인 항목이 따로 생깁니다. 키체인에 쓰지 못하면 Claude Code는 대신 폴더 안의 `.credentials.json`을 사용합니다.

## 폴더 간 공유

심볼릭 링크를 사용할 수 있습니다. Claude Code는 링크된 `settings.json`이나 `projects/` 폴더를 그대로 읽고 쓰며, 링크를 깨뜨리지 않습니다. 단, MCP 서버는 예외입니다. MCP 서버는 로그인 정보와 같은 `.claude.json`에 저장되므로, 대신 `--mcp-config`로 전달하세요.

ccam은 바로 이 규칙을 바탕으로 만들어졌습니다. ccam이 폴더를 링크하는 방식은 [공유](/docs/sharing)를, 전체적인 내용은 [여러 계정 가이드](/guides/multiple-claude-code-accounts)를 참고하세요.
