---
title: CLAUDE_CONFIG_DIR 이해하기
description: Claude Code에서 CLAUDE_CONFIG_DIR이 하는 일, 폴더에 저장되는 내용, macOS에서 로그인 정보가 보관되는 위치, 계정 간에 공유할 수 있는 파일을 설명합니다.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: ef9317030e42
---

`CLAUDE_CONFIG_DIR`은 Claude Code가 `~/.claude` 대신 사용할 폴더를 지정합니다. 로그인 정보, 설정, MCP 서버, 플러그인, 대화 내역처럼 계정별 데이터가 모두 이 폴더에 들어갑니다. 서로 다른 폴더를 가리키게 하면 계정을 분리할 수 있습니다.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR은 공식적으로 지원되나요?

예. Claude Code는 [환경 변수](https://code.claude.com/docs/en/env-vars) 목록에서 이 변수를 설정 폴더를 바꾸는 방법으로 소개하며 "여러 계정을 나란히 실행할 때 유용하다"고 설명합니다. [인증 가이드](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)에서도 여러 계정으로 로그인할 때 이 변수를 사용합니다. 셸이나 사용자 설정에서 지정할 수 있지만, 프로젝트 설정에서는 지정할 수 없습니다.

## 폴더에는 무엇이 들어 있나요?

Claude Code 2.1.302에서 확인한 내용입니다.

| 경로 | 내용 | 공유 가능 |
| --- | --- | --- |
| `.credentials.json` | 로그인 토큰(Linux, Windows), MCP OAuth 토큰 | 아니요 |
| `.claude.json` | 계정 정보, 사용자 MCP 서버, 프로젝트별 상태 | 아니요(MCP 서버는 따로) |
| `settings.json` | 권한, 훅, 모델, 상태 표시줄, 플러그인 | 예 |
| `CLAUDE.md` | 전역 지침 | 예 |
| `agents/`, `commands/`, `skills/`, `plugins/` | 확장 기능 | 예 |
| `projects/` | `.jsonl` 파일로 저장된 대화, 프로젝트별 메모리 | 예 |
| `history.jsonl` | 프롬프트 기록 | 예 |
| `sessions/`, `shell-snapshots/`, `cache/` 등 | 실행 상태 | 아니요 |

## 자주 헷갈리는 점은 무엇인가요?

**`CLAUDE_CONFIG_DIR`이 없으면 `.claude.json`은 `~/.claude` 안에 있지 않습니다.** 이때 Claude Code는 이 파일을 폴더 옆인 `~/.claude.json`에 둡니다. `CLAUDE_CONFIG_DIR=~/.claude`로 지정하기 시작하면 Claude Code는 `~/.claude/.claude.json`을 찾게 되고, 기존 계정을 찾지 못합니다. `ccam import`는 두 파일을 링크해 이 문제를 해결합니다.

**macOS에서는 로그인 정보가 폴더가 아니라 키체인에 있습니다.** Claude Code는 `Claude Code-credentials`라는 이름으로 저장하며, `CLAUDE_CONFIG_DIR`이 설정되어 있으면 이름 뒤에 폴더 경로의 짧은 해시가 붙습니다. 그래서 폴더가 다르면 키체인 항목도 따로 생깁니다. 키체인에 쓰기가 거부되면 Claude Code는 [문서](https://code.claude.com/docs/en/authentication#credential-management)에 설명된 대로 폴더 안의 `.credentials.json`에 저장합니다.

**모든 로그인이 폴더에 저장되는 것은 아닙니다.** API 키 없이 Claude Console로 로그인한 정보는 폴더 밖에 저장되므로, 폴더를 나눠도 이런 로그인 두 개를 구분할 수 없습니다. claude.ai 로그인과 API 키는 폴더별로 저장됩니다.

## Windows에서는 폴더가 어디에 있나요?

기본값은 `%USERPROFILE%\.claude`이며, 로그인 정보는 그 안의 `.credentials.json`에 있습니다. 추가 계정마다 `CLAUDE_CONFIG_DIR`을 다른 폴더로 지정하세요. Linux, macOS와 같은 방식입니다.

## 폴더끼리 파일을 공유할 수 있나요?

심볼릭 링크를 쓸 수 있습니다. Claude Code는 링크된 `settings.json`이나 `projects/` 폴더를 통해 읽고 쓰며, 링크를 그대로 유지합니다. 예외는 MCP 서버입니다. 로그인 정보와 같은 `.claude.json`에 들어 있으므로, 대신 `--mcp-config`로 전달하세요.

ccam은 바로 이 규칙을 바탕으로 만들어졌습니다. 폴더를 링크하는 방식은 [공유](/docs/sharing)를, 전체적인 그림은 [여러 계정 가이드](/guides/multiple-claude-code-accounts)를 참고하세요.
