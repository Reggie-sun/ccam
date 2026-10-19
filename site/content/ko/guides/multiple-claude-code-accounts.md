---
title: 컴퓨터 한 대에서 Claude Code 계정 두 개를 쓰는 방법
description: CLAUDE_CONFIG_DIR로 여러 Claude Code 계정을 함께 사용하는 방법입니다. 셸 별칭으로 직접 설정하거나 ccam으로 자동화할 수 있습니다.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2fbee9f4f91
---

계정마다 설정 폴더를 따로 두면 됩니다. Claude Code는 `CLAUDE_CONFIG_DIR` 환경 변수로 사용할 폴더를 정하며, 폴더마다 로그인 정보, 설정, 대화 내역이 따로 저장됩니다. 셸 별칭으로 직접 설정할 수도 있고, ccam에 맡길 수도 있습니다.

## 별칭으로 직접 설정하기

Claude Code 문서에서 소개하는 방법입니다. 추가할 계정마다 `~/.bashrc`나 `~/.zshrc`에 한 줄씩 넣습니다.

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

새 터미널을 열고 `claude-work`를 실행하세요. 새 폴더이므로 Claude Code가 로그인을 요청합니다. 이후로는 `claude`가 첫 번째 계정, `claude-work`가 두 번째 계정입니다.

계정이 두 개라면 이 방법으로도 충분합니다. 다만 다음 사항을 알아 두세요.

- 각 폴더는 빈 상태에서 시작합니다. 설정, `CLAUDE.md`, MCP 서버, 에이전트, 스킬이 하나도 없습니다. 필요한 것은 직접 복사하거나 링크해야 하고, 복사본끼리 내용을 맞추는 일도 직접 해야 합니다.
- 별칭은 대화형 셸에서만 동작합니다. 스크립트, Makefile, 다른 프로그램에서는 쓸 수 없습니다.
- MCP 서버는 로그인 정보와 같은 `.claude.json`에 저장되므로, 폴더 간에 간단히 링크할 수 없습니다.
- 대화는 시작한 계정에만 남기 때문에 다른 계정에서 이어 갈 수 없습니다.
- 모든 계정의 사용량을 한 번에 확인할 방법이 없습니다.

## ccam으로 설정하기

ccam은 같은 원리를 자동화하고, 부족한 부분을 채워 줍니다.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

이제 `claude1`과 `claude2`, 또는 `ccam personal`과 `ccam work`로 두 계정을 실행할 수 있습니다. 별칭이 아니라 실제 프로그램이므로 스크립트에서도 동작합니다. 설정, `CLAUDE.md`, MCP 서버, 확장 기능은 처음부터 공유되며, 대화 내역 공유를 켜면 `ccam next --continue`로 다른 계정에서 대화를 이어 갈 수 있습니다.

이미 `~/.claude`를 사용하고 있다면 그 계정을 다시 추가하지 말고 `ccam import`를 실행하세요.

## 어느 쪽을 선택할까요

계정이 두 개이고 공유할 것이 없다면 별칭으로 충분합니다. 설정과 MCP 서버를 공유하거나, 계정을 더 늘리거나, 사용량을 한눈에 보거나, 대화를 이어 가고 싶다면 ccam이 번거로운 관리 작업을 덜어 줍니다. 자세한 차이는 [비교](/compare/aliases) 페이지에 정리해 두었습니다.
