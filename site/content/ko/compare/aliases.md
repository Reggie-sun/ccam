---
title: ccam과 셸 별칭 비교
description: Claude Code 계정 여러 개를 쓸 때 CLAUDE_CONFIG_DIR을 쓰는 셸 별칭과 ccam을 나란히 비교합니다.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e9902676e0ed
---

계정마다 셸 별칭을 하나씩 두는 것이 Claude Code 계정 여러 개를 쓰는 가장 간단한 방법입니다. ccam도 내부에서는 같은 일을 하며, 여기에 공유, 계정 전환, 사용량 확인을 더합니다.

| | 셸 별칭 | ccam |
| --- | --- | --- |
| 로그인 분리 | 예 | 예 |
| 계정별 설정 작업 | `~/.bashrc` 수정, 로그인 | `ccam add` |
| 스크립트와 다른 프로그램에서 동작 | 아니요, 별칭은 셸에서만 동작 | 예, `claude2`는 실제 명령 |
| 설정, `CLAUDE.md`, 확장 기능 공유 | 손으로 복사하거나 링크 | 항목별로, 백업과 함께 |
| MCP 서버 공유 | 계정마다 손으로 | `ccam mcp add` 한 번 |
| 다른 계정에서 대화 이어 가기 | 아니요 | `ccs next --continue` |
| 모든 계정의 사용량 | 아니요 | `ccs usage` |
| 여유가 가장 많은 계정 고르기 | 아니요 | `ccs next --auto` |
| 상태줄에 계정 이름 표시 | 아니요 | `ccam statusline install` |
| 깨진 링크 수리 | 손으로 | `ccam doctor --fix` |
| 추가 소프트웨어 | 없음 | 정적 바이너리 하나 |

계정이 두 개이고 서로 아무것도 공유하지 않으려면 별칭을 고르세요. 여러 계정을 하나의 환경처럼 쓰고 싶다면 ccam을 고르세요. 별칭으로 시작했다가 나중에 ccam으로 옮겨도 됩니다. `ccam import`가 기존 `~/.claude-*` 폴더를 제자리에서 등록합니다.

[컴퓨터 한 대에서 계정 두 개를 쓰는 방법](/guides/multiple-claude-code-accounts)에서 두 방법을 단계별로 설명합니다.
