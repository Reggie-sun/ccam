---
title: 컴퓨터 한 대에서 Claude Code 여러 계정 사용하기
description: CLAUDE_CONFIG_DIR로 여러 Claude Code 계정을 나란히 실행하고 전환하는 방법입니다. 셸 별칭으로 직접 설정하거나 ccam으로 자동화할 수 있습니다.
updated: 2026-10-24
order: 1
translation: machine
sourceHash: "5463dec2f6eb"
---

계정마다 설정 폴더를 따로 두면 됩니다. Claude Code는 `CLAUDE_CONFIG_DIR` 환경 변수로 사용할 폴더를 정하고, 폴더마다 로그인 정보, 설정, 대화 내역이 따로 저장됩니다. 그래서 모든 계정이 로그인 상태를 유지하고, 여러 계정을 동시에 실행할 수 있습니다. 셸 별칭으로 직접 설정할 수도 있고, ccam에 맡길 수도 있습니다.

## 직접 설정하려면 어떻게 하나요?

[Claude Code 문서에서 설명하는](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) 방법입니다. 추가할 계정마다 `~/.bashrc`나 `~/.zshrc`에 한 줄씩 넣습니다.

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

새 터미널을 열고 `claude-work`를 실행하세요. 새 폴더이므로 Claude Code가 로그인을 요청합니다. 이후로는 `claude`가 첫 번째 계정, `claude-work`가 두 번째 계정입니다. 계정마다 별칭을 하나씩, 각각 다른 폴더로 추가하세요.

WSL 없이 Windows를 쓴다면 프로필(`notepad $PROFILE`)에 PowerShell 함수를 넣어 같은 효과를 얻을 수 있습니다.

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## 별칭 방식에서 빠지는 것은 무엇인가요?

계정이 두 개라면 잘 동작합니다. 다만 다음 사항을 알아 두세요.

- 각 폴더는 빈 상태에서 시작합니다. 설정, `CLAUDE.md`, MCP 서버, 에이전트, 스킬이 하나도 없습니다. 필요한 것은 직접 복사하거나 링크해야 하고, 복사본끼리 내용을 맞추는 일도 직접 해야 합니다.
- 별칭은 대화형 셸에서만 동작합니다. 스크립트, Makefile, 다른 프로그램에서는 쓸 수 없습니다.
- MCP 서버는 로그인 정보와 같은 `.claude.json`에 저장되므로, 폴더 간에 간단히 링크할 수 없습니다.
- 대화는 시작한 계정에만 남기 때문에 다른 계정에서 이어 갈 수 없습니다.
- 모든 계정의 사용량을 한 번에 확인할 방법이 없습니다.

## /login으로 계정을 전환하면 안 되나요?

가끔만 전환한다면 그래도 됩니다. Claude Code 안에서 `/login`을 입력하면 현재 폴더의 [계정이 전환됩니다](https://code.claude.com/docs/en/quickstart). 한 번에 한 계정만 로그인되므로 전환할 때마다 다시 로그인해야 하고, 두 계정을 나란히 실행할 수도 없습니다. 폴더를 나누면 두 문제가 모두 사라집니다.

## ccam으로 설정하려면 어떻게 하나요?

ccam은 같은 원리를 자동화하고, 부족한 부분을 채워 줍니다.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

이제 `ccam personal`과 `ccam work`로 두 계정을 실행할 수 있고, work 계정에 붙인 단축 명령 `ccwork`로도 실행할 수 있습니다. 별칭이 아니라 실제 프로그램이므로 스크립트에서도 동작합니다. `ccam`만 실행하면 목록에서 계정을 고를 수 있고, `ccam use work`로 기본 계정을 정할 수 있습니다.

설정, `CLAUDE.md`, MCP 서버, 확장 기능은 처음부터 공유됩니다. `ccam next --continue`는 대화를 다른 계정으로 넘겨 이어 가고, 직접 켤 수 있는 대화 내역 공유를 쓰면 그 대화가 하나로 유지됩니다. `ccam usage`는 모든 계정의 5시간·주간 사용량을 한 번에 보여 줍니다.

이미 `~/.claude`를 사용하고 있다면 그 계정을 다시 추가하지 말고 `ccam import`를 실행하세요. 또는 `ccam migrate`로 그 폴더의 설정, MCP 서버, 확장 기능을 공유 설정으로 복사할 수 있습니다.

## 지금 어떤 계정으로 실행 중인지 어떻게 알 수 있나요?

Claude Code 안에서 `/status`를 입력하면 로그인한 계정이 표시됩니다. ccam을 쓴다면 `ccam statusline install`을 한 번 실행해 두세요. 상태 표시줄에 계정 이름과 사용량이 항상 표시됩니다. 예를 들면 `[work] · Opus · 5h 24% · 7d 41%`와 같습니다.

## Windows와 WSL에서도 되나요?

Claude Code는 [Windows에서 네이티브로, 그리고 WSL 안에서](https://code.claude.com/docs/en/setup) 실행됩니다. `CLAUDE_CONFIG_DIR`은 두 환경 모두에서 동작하므로, 네이티브 Windows에서는 위의 별칭이나 PowerShell 함수로 충분합니다. ccam도 Windows에서 네이티브로 실행됩니다. PowerShell에서 `irm https://ccam.dev/install.ps1 | iex`로 설치하거나, WSL 안에서 Linux용 명령으로 설치하세요. 단축 명령은 `ccwork.exe`가 되고, `settings.json` 같은 파일을 공유하려면 개발자 모드가 필요합니다([ccam 설치](/docs/install)).

## 문제 해결

**Claude Code가 매번 로그인을 요청합니다.** 새 폴더에서 처음 실행할 때는 항상 로그인을 요청합니다. 그 뒤에도 계속 요청한다면 로그인이 만료되었거나 삭제된 것입니다. `ccam login work`를 실행하세요(또는 별칭을 한 번 실행해 로그인하세요).

**MCP 서버가 한 계정에만 보입니다.** 세션 안에서 `claude mcp add`로 추가한 서버는 그 계정에만 속하고, `-s user` 없이 추가했다면 그 폴더에만 속합니다. `ccam mcp import work <server>`로 모든 계정에 옮길 수 있습니다.

**계정마다 설정이 다릅니다.** `ccam share work`는 그 계정이 무엇을 어느 그룹으로 공유하는지 보여 줍니다. `ccam group join default work`는 바뀔 내용을 먼저 보여 준 뒤 이 계정이 `default`가 공유하는 항목을 함께 공유하게 합니다. `ccam doctor --fix`는 깨진 링크를 복구하고, 바꾸는 파일은 모두 백업해 둡니다.

**두 계정이 같은 대화를 엽니다.** 한 대화를 두 계정에서 동시에 작업하지 마세요. 그렇게 하려고 하면 ccam이 경고합니다.

## 어느 쪽을 선택해야 하나요?

계정이 두 개이고 공유할 것이 없다면 별칭으로 충분합니다. 설정과 MCP 서버를 공유하거나, 계정을 더 늘리거나, 사용량을 한눈에 보거나, 대화를 이어 가고 싶다면 ccam이 번거로운 관리 작업을 덜어 줍니다. [별칭과의 비교](/compare/aliases)에 차이를 나란히 정리했고, [ccam과 다른 계정 전환 도구](/compare/account-switchers)에서는 로그인 하나를 바꿔 끼우는 도구들을 다룹니다.

여러 계정을 써도 되는지는 각 계정의 약관에 따라 다릅니다. 사용하기 전에 확인하세요.

**Team·Enterprise 시트.** 자동 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`)은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꿉니다. Anthropic 고객센터는 "Each seat on a Team plan is assigned to one named individual, and only that person may use it"(팀 플랜의 각 시트는 한 명의 지정된 개인에게 할당되며, 그 사람만 사용할 수 있습니다), "Users can't route their work through another user's account to use that user's remaining usage"(사용자는 다른 사용자의 계정을 통해 작업을 라우팅하여 해당 사용자의 남은 사용량을 사용할 수 없습니다)라고 안내합니다([Team 요금제](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129), [Enterprise 요금제](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)도 마찬가지). 따라서 시트 사이에서 자동 선택을 쓰면 Anthropic의 약관을 위반할 수 있습니다. 이 기능에 의존하기 전에 소속 조직과 Anthropic에 확인하세요. 개인 요금제(Pro, Max)에 대해서는 위의 안내를 참고하세요.
