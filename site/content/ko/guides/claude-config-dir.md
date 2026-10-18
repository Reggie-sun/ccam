---
title: CLAUDE_CONFIG_DIR 알아보기
description: Claude Code의 CLAUDE_CONFIG_DIR이 하는 일, 폴더 내용, macOS의 로그인 위치, 계정끼리 공유할 수 있는 파일.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR`은 Claude Code가 `~/.claude` 대신 쓸 폴더를 지정합니다. 로그인, 설정, MCP 서버, 플러그인, 대화 내역처럼 계정마다 다른 것은 모두 이 폴더에 있습니다. 이 변수가 서로 다른 폴더를 가리키게 하면 계정이 분리됩니다.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## 폴더 안에 있는 것

Claude Code 2.1.302에서 확인한 내용입니다.

| 경로 | 내용 | 공유 가능 |
| --- | --- | --- |
| `.credentials.json` | 로그인 토큰(Linux), MCP OAuth 토큰 | 아니요 |
| `.claude.json` | 계정 정보, 사용자 MCP 서버, 프로젝트별 상태 | 아니요(MCP 서버는 따로) |
| `settings.json` | 권한, 훅, 모델, 상태줄, 플러그인 | 예 |
| `CLAUDE.md` | 전역 지침 | 예 |
| `agents/`, `commands/`, `skills/`, `plugins/` | 확장 기능 | 예 |
| `projects/` | `.jsonl` 파일로 된 대화, 프로젝트별 메모리 | 예 |
| `history.jsonl` | 프롬프트 기록 | 예 |
| `sessions/`, `shell-snapshots/`, `cache/` 등 | 실행 중 상태 | 아니요 |

## 많이들 놓치는 두 가지

**`CLAUDE_CONFIG_DIR`이 없으면 `.claude.json`은 `~/.claude` 안에 있지 않습니다.** 이때 Claude Code는 이 파일을 폴더 옆인 `~/.claude.json`에 둡니다. `CLAUDE_CONFIG_DIR=~/.claude`를 쓰기 시작하면 Claude Code는 대신 `~/.claude/.claude.json`을 찾고, 기존 계정을 찾지 못합니다. `ccam import`는 두 위치를 링크해 이 문제를 처리합니다.

**macOS에서는 로그인이 폴더가 아니라 키체인에 있습니다.** Claude Code는 로그인을 `Claude Code-credentials`라는 이름으로 저장하며, `CLAUDE_CONFIG_DIR`이 설정되어 있으면 그 뒤에 폴더 경로의 짧은 해시를 붙입니다. 따라서 폴더가 다르면 키체인 항목도 따로 생깁니다. 키체인이 쓰기를 거부하면 Claude Code는 폴더 안의 `.credentials.json`을 대신 씁니다.

## 폴더끼리 공유하기

심볼릭 링크를 쓸 수 있습니다. Claude Code는 링크된 `settings.json`이나 `projects/` 폴더를 통해 읽고 쓰며, 링크를 그대로 유지합니다. 예외는 MCP 서버입니다. MCP 서버는 로그인과 같은 `.claude.json`에 들어 있으므로, 대신 `--mcp-config`로 전달하세요.

ccam은 바로 이 규칙들을 바탕으로 동작합니다. ccam이 폴더를 링크하는 방법은 [공유](/docs/sharing)를, 전체 그림은 [여러 계정 가이드](/guides/multiple-claude-code-accounts)를 참고하세요.
