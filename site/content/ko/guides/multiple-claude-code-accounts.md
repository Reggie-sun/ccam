---
title: 컴퓨터 한 대에서 Claude Code 계정 두 개를 쓰는 방법
description: CLAUDE_CONFIG_DIR로 Claude Code 계정 여러 개를 나란히 씁니다. 셸 별칭으로 직접 설정하거나 ccam으로 자동화합니다.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2fbee9f4f91
---

계정마다 설정 폴더를 따로 두면 됩니다. Claude Code는 `CLAUDE_CONFIG_DIR` 환경 변수에서 폴더를 읽고, 각 폴더에는 그 계정의 로그인, 설정, 대화 내역이 따로 들어 있습니다. 이 설정은 셸 별칭으로 직접 할 수도 있고, ccam에 맡길 수도 있습니다.

## 별칭으로 직접 설정하기

Claude Code 문서에서 설명하는 방법입니다. 추가 계정마다 `~/.bashrc`나 `~/.zshrc`에 한 줄씩 넣습니다.

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

새 터미널을 열고 `claude-work`를 실행합니다. 폴더가 새로 만들어졌으므로 Claude Code가 로그인을 요청합니다. 그 뒤로는 `claude`가 첫 번째 계정, `claude-work`가 두 번째 계정입니다.

계정이 두 개라면 이 방법으로도 잘 됩니다. 알아 둘 점은 다음과 같습니다.

- 각 폴더는 빈 상태로 시작합니다. 설정도, `CLAUDE.md`도, MCP 서버도, 에이전트나 스킬도 없습니다. 필요한 것을 복사하거나 링크하고, 복사본은 직접 맞춰 두어야 합니다.
- 별칭은 대화형 셸에서만 존재합니다. 스크립트, Makefile, 다른 프로그램에서는 보이지 않습니다.
- MCP 서버는 로그인과 같은 `.claude.json`에 들어 있으므로, 폴더끼리 간단히 링크할 수 없습니다.
- 대화는 시작한 계정에 남으므로, 다른 계정에서 이어 갈 수 없습니다.
- 모든 계정의 사용량을 한 번에 보여 주는 수단이 없습니다.

## ccam으로 설정하기

ccam은 같은 방식을 자동화하고, 빠진 부분을 채웁니다.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

이제 `claude1`과 `claude2`, 또는 `ccam personal`과 `ccam work`로 두 계정을 시작합니다. 이것들은 별칭이 아니라 실제 프로그램이므로 스크립트에서도 동작합니다. 설정, `CLAUDE.md`, MCP 서버, 확장 기능은 처음부터 공유되며, 대화 내역 공유를 켜면 `ccam next --continue`로 다른 계정에서 대화를 이어 갈 수 있습니다.

이미 `~/.claude`를 쓰고 있다면 그 계정을 다시 추가하는 대신 `ccam import`를 실행하세요.

## 어느 쪽을 고를까

계정이 두 개이고 아무것도 공유하지 않는다면 별칭으로 충분합니다. 설정과 MCP 서버 공유, 더 많은 계정, 사용량 한눈에 보기, 대화 이어 가기가 필요하다면 ccam이 그 관리 작업을 덜어 줍니다. [비교](/compare/aliases)에서 차이를 나란히 정리했습니다.
