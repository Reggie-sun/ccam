---
title: 사용량 한도와 계정 전환
description: ccam usage로 모든 Claude Code 계정의 5시간·주간 사용량을 확인하고, ccam next --auto로 여유가 가장 많은 계정을 엽니다.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 9aebb54d40e2
---

`ccam usage`는 각 계정이 5시간 한도와 주간 한도를 얼마나 썼는지 보여 줍니다. `ccam next --auto`는 여유가 가장 많은 계정을 열고, `--continue`는 지금 대화를 그 계정에서 이어 갑니다.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

각 열은 순서대로 계정, 5시간 사용량, 초기화까지 남은 시간, 주간 사용량, 초기화까지 남은 시간, 수치를 가져온 시점을 나타냅니다.

## 수치의 출처

ccam은 각 계정의 로그인 정보로 Anthropic의 사용량 엔드포인트에 모든 계정의 사용량을 동시에 요청합니다. 계정당 최대 3초까지 기다리며, 응답이 온 계정부터 바로 표시합니다.

이 엔드포인트는 공개 문서가 없고, 바뀔 수도 있습니다. 요청이 실패하면 ccam은 마지막으로 받은 수치를 보여 주고 캐시된 값임을 표시합니다. 그 수치를 저장한 뒤 초기화된 구간은 0%로 표시합니다. 응답 형식이 바뀌면 `ccam doctor`가 알려 줍니다.

가져온 지 60초가 안 된 수치는 다시 요청하지 않고 그대로 씁니다. 사용할 수 있는 옵션은 다음과 같습니다.

| 옵션 | 동작 |
| --- | --- |
| `--watch [seconds]` | 계속 새로 고침(기본 60초 간격) |
| `--json` | 기계가 읽기 쉬운 형식으로 출력 |
| `-v` | 모델별 주간 한도와 각 수치의 출처도 표시 |
| `--no-cache` | 항상 새로 요청 |
| `--source cache` | 새로 요청하지 않고 저장된 수치만 표시 |

## 만료된 로그인

계정의 액세스 토큰이 만료되어도 ccam은 토큰을 갱신하지 않습니다. 같은 로그인을 쓰는 Claude Code 세션이 실행 중일 때 갱신하면 그 세션이 로그아웃될 수 있기 때문입니다. `ccam <name>`으로 그 계정을 한 번 시작하면 Claude Code가 알아서 갱신합니다.

`ccam usage --refresh`는 만료된 토큰을 직접 갱신하지만, 실행 중이 아닌 계정만 대상으로 합니다. macOS 키체인이나 가져온 `~/.claude`에 저장된 로그인은 절대 갱신하지 않습니다.

## 상태 표시줄에 사용량 표시하기

```sh
ccam statusline install
```

그러면 Claude Code 상태 표시줄에 계정 이름과 사용량이 표시됩니다(예: `[work] · Opus · 5h 24% · 7d 41%`). 이미 상태 표시줄 명령을 쓰고 있다면 ccam은 그 명령을 그대로 두고 앞에 계정 이름을 붙입니다. 상태 표시줄은 수치도 저장해 두며, 엔드포인트가 응답하지 않으면 `ccam usage`가 이 값을 사용합니다.

## 계정 전환

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next`는 로그인되지 않은 계정을 건너뛰고, `--auto`는 사용량이 100%인 계정을 건너뜁니다. 대화를 이어 가려면 두 계정 모두 대화 내역을 공유해야 합니다. 자세한 내용은 [공유](/docs/sharing)를 참고하세요.

[사용량 한도 대응 가이드](/guides/usage-limit-workflow)에서 전체 과정을 예시와 함께 살펴볼 수 있습니다.
