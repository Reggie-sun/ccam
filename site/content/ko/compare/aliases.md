---
title: ccam과 셸 별칭 비교
description: 여러 Claude Code 계정을 쓸 때, CLAUDE_CONFIG_DIR 기반 셸 별칭과 ccam을 나란히 비교합니다.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: b9fbcb88f111
---

여러 Claude Code 계정을 쓰는 가장 간단한 방법은 계정마다 셸 별칭을 하나씩 만드는 것입니다. ccam도 내부적으로는 같은 방식을 쓰며, 여기에 공유, 계정 전환, 사용량 확인 기능을 더합니다.

## 별칭으로 설정하기

추가 계정마다 `~/.bashrc` 또는 `~/.zshrc`에 한 줄을 넣고, 새 터미널을 연 다음 별칭을 한 번 실행해 로그인합니다.

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude`는 그대로 첫 번째 계정이고, `claude-work`는 두 번째 계정을 시작합니다. [Claude Code 문서에서 설명하는 방법](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)이 바로 이것입니다.

## ccam으로 설정하기

ccam을 설치하고 계정을 추가합니다. 여기서는 계정 전용 단축 명령도 함께 만듭니다.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work`와 `ccwork`로 계정을 시작할 수 있고, 둘 다 스크립트에서도 동작합니다.

## 나란히 비교

| | 셸 별칭 | ccam |
| --- | --- | --- |
| 로그인 분리 | 예 | 예 |
| 계정별 설정 방법 | `~/.bashrc` 수정 후 로그인 | `ccam add` |
| 스크립트와 다른 프로그램에서 사용 | 아니요(별칭은 셸에서만 동작) | 예(`ccam work`와 `ccwork` 같은 단축 명령은 실제 명령) |
| 설정, `CLAUDE.md`, 확장 기능 공유 | 직접 복사하거나 링크 | 항목별로 공유, 백업 포함 |
| MCP 서버 공유 | 계정마다 직접 설정 | `ccam mcp add` 한 번 |
| 다른 계정에서 대화 이어 가기 | 아니요 | `ccam next --continue` |
| 모든 계정의 사용량 확인 | 아니요 | `ccam usage` |
| 여유가 가장 많은 계정 선택 | 아니요 | `ccam next --auto` |
| 한도에 도달한 `-p` 작업을 다른 계정으로 다시 실행 | 아니요 | `ccam --fallback` |
| 프로젝트 폴더마다 계정 지정 | 아니요 | `ccam folder` |
| 상태 표시줄에 계정 이름 표시 | 아니요 | `ccam statusline install` |
| 깨진 링크 복구 | 직접 | `ccam doctor --fix` |
| 추가 소프트웨어 | 없음 | 정적 바이너리 하나 |

계정이 두 개이고 서로 아무것도 공유할 필요가 없다면 별칭을 선택하세요. 여러 계정을 하나의 환경처럼 쓰고 싶다면 ccam을 선택하세요. 별칭으로 시작했다가 나중에 ccam으로 옮겨도 됩니다. `ccam import`로 기존 `~/.claude-*` 폴더를 원래 위치에서 그대로 등록할 수 있습니다.

[컴퓨터 한 대에서 여러 계정을 쓰는 방법](/guides/multiple-claude-code-accounts)에서 두 방식을 단계별로 설명하고, [ccam과 다른 계정 전환 도구](/compare/account-switchers)에서는 로그인 하나를 바꿔 끼우는 도구와 ccam을 비교합니다.
