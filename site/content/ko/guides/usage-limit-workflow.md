---
title: Claude Code 사용량 한도에 도달했을 때 할 일
description: 한도가 초기화될 때까지 기다리거나 다른 Claude Code 계정에서 같은 대화를 바로 이어 가는 방법입니다. 5시간 한도와 주간 한도도 설명합니다.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "049a46684c7b"
---

계정이 사용량 한도(흔히 rate limit이라고 부릅니다)에 도달하면 선택지는 두 가지입니다. 한도가 초기화될 때까지 기다리거나, 가지고 있는 다른 계정에서 같은 대화를 지금 이어 가는 것입니다. ccam에서는 두 번째 방법이 명령 하나 `ccam next --auto --continue`로 끝납니다.

## Claude Code에는 어떤 한도가 있나요?

Pro와 Max 요금제의 Claude Code에는 [5시간 세션 한도와 주간 한도](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)가 있습니다. 두 한도 모두 웹과 앱의 Claude와 함께 쓰며, 사용량은 두 한도에 동시에 반영됩니다. 둘 중 하나가 소진되면 Claude Code는 "You've hit your session limit · resets 3:45pm" 같은 [메시지를 표시합니다](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit). 한도는 모든 모델에 적용되므로 모델을 바꿔도 소용없습니다.

## 방법 1: 초기화 기다리기

계정이 하나뿐이라면 기다리세요. 최신 버전의 Claude Code는 [열린 세션에서 기다렸다가](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) 한도가 초기화되면 작업을 스스로 이어 갑니다. Claude Code 안에서 `/usage`를 입력하면 요금제 한도를 얼마나 썼는지 볼 수 있습니다.

## 방법 2: 지금 다른 계정에서 이어 가기

두 번째 계정이 있다면 같은 대화로 바로 작업을 이어 갈 수 있습니다. 한 번만 설정해 두세요.

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

이제 작업하는 동안 상태 표시줄에 계정과 5시간 사용량이 표시되므로, 갑자기 한도에 걸려 당황할 일이 없습니다.

한도에 도달하면 Claude Code를 종료한 뒤(`/exit`) 다음을 실행합니다.

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

여유가 가장 많은 계정에서 대화를 엽니다.

```sh
ccam next --auto --continue
```

Claude Code가 `work` 계정으로 시작되고, 현재 폴더의 가장 최근 대화를 지금까지 나눈 내용 그대로 불러옵니다. 대화 내역을 공유하면(위의 `ccam group set default --add history`) 대화가 하나로 유지됩니다. 공유하지 않으면 Claude Code가 대화를 `work`로 복사하고, ccam은 첫 응답이 토큰을 얼마나 다시 읽는지 알려 줍니다.

Claude Code에 [ccam 플러그인](/docs/install#claude-code-plugin)을 설치해 두면, 한도에 걸려 멈춘 세션에서 알림이 뜨고 가득 찬 구간을 기록하므로 `ccam next --auto`가 초기화 전에 그 계정을 다시 고르지 않으며, `/ccam:switch`(또는 `/ccam:switch work`)가 `ccam next --continue`를 위해 대화를 기록해 두므로 어느 대화를 이어 갈지 헷갈릴 일이 없습니다.

**Team·Enterprise 시트.** 자동 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`)은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꿉니다. Anthropic 고객센터는 "Each seat on a Team plan is assigned to one named individual, and only that person may use it"(팀 플랜의 각 시트는 한 명의 지정된 개인에게 할당되며, 그 사람만 사용할 수 있습니다), "Users can't route their work through another user's account to use that user's remaining usage"(사용자는 다른 사용자의 계정을 통해 작업을 라우팅하여 해당 사용자의 남은 사용량을 사용할 수 없습니다)라고 안내합니다([Team 요금제](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129), [Enterprise 요금제](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)도 마찬가지). 따라서 시트 사이에서 자동 선택을 쓰면 Anthropic의 약관을 위반할 수 있습니다. 이 기능에 의존하기 전에 소속 조직과 Anthropic에 확인하세요. 개인 요금제(Pro, Max)에 대해서는 [FAQ의 약관 안내](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 참고하세요.

## ccam이 자동으로 계정을 전환하나요?

아니요. 작업하는 동안 ccam은 실행되고 있지 않습니다. ccam은 Claude Code를 시작하면서 자기 자신을 Claude Code로 대체합니다. 언제 전환할지는 사용자가 정하며, `ccam next --auto`는 5시간 사용량이 가장 낮은 계정을 고르고 100%에 도달했거나, 로그인하지 않았거나, 다른 조직에 속한 계정은 건너뜁니다. `config.toml`에 `relay = true`를 설정하면 ccam이 대신 세션 옆에서 기다리므로, `/ccam:switch`가 세션의 옵션, 권한 모드, 추론 수준을 유지한 채 스스로 대화를 이어 갑니다. 다른 도구 중에는 실행 중인 세션의 로그인을 바꿔 끼우는 것도 있습니다. 그 차이는 [ccam과 다른 계정 전환 도구](/compare/account-switchers)에서 설명합니다.

## 알아 둘 점

- `ccam next --continue`는 작업하던 폴더에서 실행하세요. 대화는 프로젝트 폴더별로 저장됩니다.
- 같은 대화를 두 계정에서 동시에 실행하지 마세요. 그렇게 하려고 하면 ccam이 경고합니다.
- 대화 내역을 공유하지 않으면 대화가 다른 계정으로 복사되고, 그 뒤로 두 사본은 따로 이어집니다. `/rewind`가 코드를 되돌릴 때 쓰는 체크포인트는 처음 계정에 남습니다.
- 새 계정에서 보내는 첫 메시지는 평소보다 한도를 더 많이 쓸 수 있습니다. 그 계정에는 아직 대화가 캐시되지 않았기 때문입니다.
- 이렇게 계정을 전환해도 되는지는 각 계정의 약관에 따라 다릅니다.

`ccam usage`와 `ccam next`의 모든 옵션은 [사용량 한도](/docs/usage-limits)를 참고하세요.
