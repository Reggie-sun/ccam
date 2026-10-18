---
title: 사용량 한도와 계정 전환
description: ccs usage로 모든 Claude Code 계정의 5시간·주간 사용량을 보고, ccs next --auto로 여유가 가장 많은 계정을 엽니다.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a9ee77936668
---

`ccs usage`는 각 계정이 5시간 한도와 주간 한도를 얼마나 썼는지 보여 줍니다. `ccs next --auto`는 여유가 가장 많은 계정을 열고, `--continue`는 현재 대화를 그 계정으로 이어 갑니다.

```console
$ ccs usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccs next --auto 로 바로 열기)
```

각 열은 차례로 계정, 5시간 사용량, 초기화까지 남은 시간, 주간 사용량, 초기화까지 남은 시간, 그리고 수치가 얼마나 오래되었는지(기준)를 나타냅니다.

## 수치는 어디에서 오나요

ccam은 계정마다 그 계정의 로그인으로 Anthropic의 사용량 엔드포인트에 동시에 요청합니다. 계정당 최대 3초까지 기다리며, 응답이 오는 대로 각 계정을 표시합니다.

이 엔드포인트는 공개 문서가 없으며 바뀔 수 있습니다. 요청이 실패하면 ccam은 마지막으로 알던 수치를 보여 주고 캐시된 값이라고 표시합니다. 응답 형식이 바뀌었다면 `ccam doctor`가 알려 줍니다.

가져온 지 60초가 지나지 않은 수치는 다시 사용합니다. 옵션은 다음과 같습니다.

| 옵션 | 효과 |
| --- | --- |
| `--watch [seconds]` | 계속 새로 고침. 기본값은 60초마다 |
| `--json` | 기계가 읽을 수 있는 형식으로 출력 |
| `-v` | 모델별 주간 한도와 각 수치의 출처도 표시 |
| `--no-cache` | 항상 새로 요청 |
| `--source cache` | 요청하지 않고 저장된 수치만 표시 |

## 만료된 로그인

계정의 액세스 토큰이 만료되어도 ccam은 갱신하지 않습니다. 같은 로그인을 쓰는 Claude Code 세션이 실행 중일 때 갱신하면 그 세션이 로그아웃될 수 있기 때문입니다. `ccs <name>`으로 그 계정을 한 번 시작하면 Claude Code가 갱신합니다.

`ccs usage --refresh`는 만료된 토큰을 직접 갱신하지만, 실행 중이 아닌 계정에만 그렇게 합니다. macOS 키체인이나 가져온 `~/.claude`에 저장된 로그인은 절대 갱신하지 않습니다.

## 상태줄의 사용량

```sh
ccam statusline install
```

그러면 Claude Code가 상태줄에 계정 이름과 사용량을 보여 줍니다. 예: `[work] · Opus · 5h 24% · 7d 41%`. 이미 상태줄 명령을 쓰고 있다면 ccam은 그 명령을 유지하고 앞에 계정 이름을 붙입니다. 상태줄은 수치도 저장해 두며, 엔드포인트가 응답하지 않을 때 `ccs usage`가 이 값을 씁니다.

## 계정 전환

```sh
ccs next              # the account after the one you used last
ccs next --auto       # the account with the lowest 5-hour usage
ccs next --continue   # the next account that shares your conversation history
```

`ccs next`는 로그인되지 않은 계정을 건너뛰고, `--auto`는 100%에 도달한 계정을 건너뜁니다. 대화를 이어 가려면 두 계정이 모두 대화 내역을 공유해야 합니다. [공유](/docs/sharing)를 참고하세요.

[사용량 한도 대처 가이드](/guides/usage-limit-workflow)에서 전체 예시를 단계별로 살펴봅니다.
