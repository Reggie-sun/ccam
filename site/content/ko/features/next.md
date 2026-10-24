---
title: 다음 계정으로 전환
description: ccam next는 다른 계정으로 Claude Code를 엽니다. --auto는 여유가 가장 많은 계정을 고르고, --continue는 대화를 넘겨받으며, --fallback은 -p 작업을 재실행합니다.
updated: 2026-10-24
order: 8
translation: machine
sourceHash: "a38daed5dd97"
---

`ccam next`는 계정 이름을 지정하지 않아도 다른 계정으로 Claude Code를 엽니다. 계정이 사용량 한도에 도달했을 때 쓰면 됩니다. 그냥 `ccam next`를 실행하면 마지막으로 사용한 계정의 다음 계정을, `--auto`를 붙이면 여유가 가장 많은 계정을 고릅니다. `--continue`를 붙이면 두 계정이 대화 내역을 공유하든 안 하든 이 폴더의 마지막 대화를 새 계정으로 넘겨 이어 갑니다. Claude Code 안에서는 [ccam 플러그인](#inside-claude-code)이 한도 때문에 세션이 멈추면 알려 주고, `/ccam:switch`로 대화를 넘길 수 있습니다. `claude -p`를 실행하는 스크립트에서는 `ccam --fallback`이 이 전환을 알아서 합니다([한도에서 다시 실행](#run-again-at-the-limit) 참고).

**Team·Enterprise 시트.** 자동 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`)은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꿉니다. Anthropic 고객센터는 "Each seat on a Team plan is assigned to one named individual, and only that person may use it"(팀 플랜의 각 시트는 한 명의 지정된 개인에게 할당되며, 그 사람만 사용할 수 있습니다), "Users can't route their work through another user's account to use that user's remaining usage"(사용자는 다른 사용자의 계정을 통해 작업을 라우팅하여 해당 사용자의 남은 사용량을 사용할 수 없습니다)라고 안내합니다([Team 요금제](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129), [Enterprise 요금제](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)도 마찬가지). 따라서 시트 사이에서 자동 선택을 쓰면 Anthropic의 약관을 위반할 수 있습니다. 이 기능에 의존하기 전에 소속 조직과 Anthropic에 확인하세요. 개인 요금제(Pro, Max)에 대해서는 [FAQ의 약관 안내](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 참고하세요.

## 사용법

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam은 표준 오류에 전환 내용을 `→ <from> → <to>` 형식으로 출력하고(`--auto`로 골랐다면 그 이유도 함께), `ccam <name>`과 같은 방식으로 새 계정의 Claude Code를 시작합니다. `<from>`은 마지막으로 사용한 계정이며, `--continue`를 쓰면 대화가 속한 계정입니다. `--` 뒤의 인자는 그대로 `claude`에 전달됩니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--auto` | 로그인된 계정을 사용량으로 고릅니다(아래 참고). 사용량이 100%인 계정은 건너뜁니다 | 꺼짐 |
| `--max-usage <1-100>` | `--auto`와 함께: 사용량이 이 값 미만인 계정을 우선합니다 | `auto_max_usage`, 없으면 100 |
| `--strategy headroom\|expiring\|balance` | `--auto`와 함께: 5시간 사용량이 가장 낮은 계정, 5시간 구간이 가장 먼저 초기화되는 계정, 또는 주간 할당량이 가장 먼저 쓰이지 못하고 남을 계정 | `auto_strategy`, 없으면 `headroom` |
| `--cross-org` | `--auto`와 함께: 다른 조직의 계정도 고릅니다(아래 참고) | `cross_org`, 없으면 꺼짐 |
| `--wait[=<duration>]` | `--auto`와 함께: 여유 있는 계정이 없으면 첫 계정이 초기화될 때까지 최대 이 시간만큼 기다립니다. 값은 `=` 뒤에 씁니다. `--wait 2h`는 `--wait takes its value after '=': --wait=2h`라는 힌트와 함께 거부합니다 | 꺼짐, 값 없이 주면 `5h` |
| `-c`, `--continue` | 이 폴더의 마지막 대화를 새 계정으로 넘겨 이어 갑니다(아래 참고) | 꺼짐 |
| `-- <claude args>` | ccam 자체 인자 뒤에 `claude`로 넘길 인자 | 없음 |
| `--dry-run` | 지원합니다. Claude Code를 시작하거나 실행 기록을 남기지 않고, `→ <from> → <to>`와 실행할 명령줄을 출력합니다 | 꺼짐 |

`--auto`와 `--continue`는 함께 쓸 수 있습니다. `--continue`를 쓰면 `--resume` 인자가 직접 넘긴 인자보다 앞에 놓입니다.

`--dangerously-skip-permissions`나 `--model` 같은 Claude Code 자체 옵션은 `--` 뒤에 둡니다. 이름 뒤의 모든 것을 그대로 넘기는 `ccam <name>`과 달리 `ccam next`는 인자를 직접 읽으므로, `--`가 없으면 `ccam: unknown flag: --dangerously-skip-permissions` 오류를 내고 멈춥니다.

## 동작 방식

ccam으로 실행할 때마다 그 시각이 `~/.local/state/ccam/last-used/<name>`에 기록됩니다. `ccam next`는 이 파일들을 읽고, 어떤 계정이 로그인되어 있는지 확인한 뒤 계정을 고릅니다.

1. `cli.candidates`가 `config.toml`의 계정마다 `next.Candidate`를 하나씩 만듭니다. 여기에는 계정, 마지막 실행 시각(`launcher.LastUsed`), 동작하는 로그인이 있는지(`creds.Has`: `.credentials.json`, macOS에서는 키체인에 있는 Claude 로그인으로, 통째로 만료되지 않았고 Claude Code가 지우지 않은 것)가 들어갑니다. 대신 API 키를 쓰는 폴더(Console 로그인이나 `ccam add --api`로 추가한 계정)에는 요금제 한도가 없으므로 절대 고르지 않습니다.
2. `--auto`가 없으면 `next.Rotate`가 계정을 슬롯 번호순으로 정렬하고 가장 최근에 실행된 계정을 찾습니다. 그 계정이 지금 전환해 떠나는 계정입니다. 아직 실행된 계정이 없으면 가장 낮은 슬롯이 기준이 됩니다.
3. 그 지점부터 슬롯 순서대로 앞으로 진행하며, 마지막 슬롯 다음에는 첫 슬롯으로 돌아갑니다. 처음 만나는 로그인된 계정을 반환합니다.
4. `launcher.LaunchWith`가 고른 계정을 `ccam <name>`과 똑같이 준비하고 시작합니다. 계정의 링크를 확인하고, MCP 공유를 적용하고, 세션을 등록한 뒤 `syscall.Exec`로 ccam을 `claude`로 바꿉니다.

탐색은 마지막으로 사용한 계정에서 시작해, 그 뒤에 처음 나오는 로그인된 계정에서 멈춥니다.

<!-- motion:next-rotate -->

`work`가 가장 최근에 실행되었으므로 탐색은 그 다음부터 시작합니다. `side`는 로그인되어 있지 않아 건너뛰고, `lab`이 선택됩니다. `lab`도 로그아웃되어 있었다면 탐색은 처음으로 돌아가 `personal`로 갑니다(그림의 `% len`). 순환 방식은 출발한 계정을 절대 반환하지 않습니다. 로그인된 다른 계정이 없으면 오류를 내고 멈춥니다.

### --auto

`--auto`는 슬롯 순서가 아니라 사용량으로 결정합니다. `cli.pickAuto`는 로그인된 계정만 남기고, `ccam usage`가 쓰는 것과 같은 `Fetcher.FetchAll`로 모든 계정의 사용량을 한꺼번에 요청한 뒤 `usage.Policy.Choose()`에게 선택을 맡깁니다.

<!-- motion:next-auto -->

계정은 5시간 구간과 주간 구간으로 판단하며, Claude Code에 넘기는 인자(`--` 뒤)에 `--model`이 있으면 그 모델의 주간 구간도 따집니다(이름에 `opus`, `sonnet`, `haiku`, `fable`이 들어간 경우). `Policy.Choose()`는 계정을 세 부류로 나눕니다.

1. **여유**: 모든 구간이 `auto_max_usage`(`config.toml`, 기본값 100) 미만인 계정. 이 중에서는 전략에 따라 고릅니다. `headroom`(기본값)은 5시간 사용량이 가장 낮은 계정을, `expiring`은 5시간 구간이 가장 먼저 초기화되는 계정을 골라, 남은 할당량이 사라지기 전에 쓰게 합니다. `balance`는 주간 구간이 초기화될 때까지 하루당 남은 주간 할당량이 가장 많은 계정을 고르며, 실행 중인 ccam 세션도 그 계정의 몫으로 계산합니다. 그래서 주간 할당량이 쓰이지 못하고 남지 않으며, 병렬 세션이 여러 계정에 고루 나뉩니다.
2. **초과**: 100% 미만이지만 `auto_max_usage` 이상인 계정. 여유 있는 계정이 하나도 없을 때만 쓰며, 5시간 사용량이 가장 낮은 계정부터 고릅니다.
3. **가득 참**: 100%인 구간이 있거나 5시간 사용량을 알 수 없는 계정. 절대 고르지 않습니다.

Claude Code가 어떤 구간에서 요청을 거부했으면, 저장된 수치와 상관없이 그 구간은 초기화될 때까지 100%로 칩니다. [ccam 플러그인](#inside-claude-code)의 훅이 이를 기록합니다(`~/.local/state/ccam/usage/<name>.limit.json`). 그래서 방금 한도에 도달한 계정은 수치가 오래되었더라도 초기화 전에 다시 고르지 않습니다.

모든 계정에 적용되는 설정은 `config.toml`의 맨 위, 첫 `[[group]]`, `[[profile]]`, `[[folder]]`보다 위에 둡니다. 이들 아래에 둔 설정은 그 테이블에 속하므로, ccam은 실행할 때마다 그 설정을 알려 주고, 설정을 저장하는 명령은 옮길 때까지 거부합니다.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**조직.** 자동으로 고를 때는 떠나는 계정과 비슷한 계정 안에서만 고릅니다. 같은 Team이나 Enterprise 조직의 계정끼리, 또는 개인 요금제(Pro, Max) 계정끼리 고르므로, 회사 계정의 작업이 개인 계정에서 이어지거나 그 반대가 되는 일이 없습니다. ccam은 각 계정의 조직을 로그인에서 알아냅니다(`ccam login`이 기록하며, 이전 로그인은 `.claude.json`에서 읽습니다). 조직을 알 수 없는 계정은 어디에나 들어갈 수 있습니다. 같은 Claude 계정으로 로그인한 두 폴더는 한도를 함께 쓰므로, 한쪽이 다른 쪽을 이어받는 일은 없습니다. `--cross-org`나 `config.toml`의 `cross_org = true`를 쓰면 조직 규칙이 풀립니다.

고를 수 있는 계정이 없으면 ccam은 종료 코드 75(`EX_TEMPFAIL`)로 끝나므로, 스크립트는 "나중에 다시 시도"와 실패를 구별할 수 있습니다. `--wait`를 주면 대신, 가득 찬 계정 중 막고 있는 구간이 가장 먼저 초기화될 때까지(엔드포인트에 새 구간이 나타나는 초기화 시각 30초 뒤) 기다리며 `ccam: no account has room; waiting until 15:10`을 출력하고, 다시 고릅니다. 초기화 시각을 알려 주는 계정이 없으면 10분 뒤에 다시 확인합니다. 기다리는 시간이 끝날 때까지 여유가 생기지 않으면 75로 종료합니다. 조회는 `ccam usage`의 규칙을 따릅니다. 가져온 지 60초가 안 된 수치는 다시 쓰고, 계정마다 3초까지 기다리며, 응답하지 않거나 액세스 토큰이 만료된 계정은 저장된 수치로 판단하고, 한도 기록이 그 위에 덮입니다. 저장된 구간의 초기화 시각이 지났으면 0%로 계산하므로, 그 사이에 구간이 초기화된 계정이 밀려나지 않습니다. `--auto`는 토큰을 절대 갱신하지 않습니다. 일반 순환과 달리 `--auto`는 마지막으로 사용한 계정의 여유가 여전히 가장 많으면 그 계정을 고를 수도 있습니다.

### --continue

1. ccam은 모든 계정에서 이 폴더의 마지막 대화를 찾습니다. `projects/<folder>/`의 `<session>.jsonl` 파일 중 가장 최근의 프롬프트나 응답이 있는 파일입니다(폴더 이름은 Claude Code가 짓는 방식대로, 링크를 푼 작업 폴더 경로에서 문자와 숫자가 아닌 모든 글자를 `-`로 바꾼 것입니다). `/mcp`, `/usage`, `/status`만 실행한 세션이나 빈 파일처럼 프롬프트나 응답이 없는 세션은 건너뜁니다. 세션이 메모를 남겼고(사용량 한도나 `/ccam:switch`, 아래 참고) 그보다 2분 넘게 뒤에 프롬프트나 응답이 있는 대화가 없으면 그 세션을 그대로 씁니다. Claude Code가 `/exit`에서 쓰는 줄은 치지 않습니다.
2. 대화가 속한 계정이 전환해 떠나는 계정이 되고, 새 계정은 위와 같은 방식으로 같은 조직 안에서 고릅니다. 여러 계정이 그 대화 내역을 공유하면, `.claude.json`에 그 대화를 마지막으로 실행했다고 적힌(`lastSessionId`) 계정의 것으로, 없으면 마지막으로 사용한 계정의 것으로 봅니다.
3. 새 계정이 그 대화 내역을 공유하면(같은 그룹 폴더) ccam은 `claude --resume <session id>`로 바로 그 대화를 시작합니다(`--continue`로는 그 폴더에서 가장 최근의 대화가 열립니다).
4. 그렇지 않으면 `claude --resume <the conversation's file> --fork-session`으로 시작합니다. Claude Code가 대화 전체를 새 세션 ID로 새 계정 자체의 대화 내역에 복사하며, 원본은 그대로 남습니다. 그 뒤로는 새 계정의 `--continue`가 이 대화를 찾습니다.

ccam은 대화의 크기를 알려 줍니다. 새 계정은 첫 응답에서 대화 전체를 다시 읽기 때문입니다. 프롬프트 캐시는 조직 사이에 절대 공유되지 않고 개인 요금제는 각각 별도의 조직이므로, 두 계정이 서로 다른 조직에 속하면 ccam이 `prompt caches are not shared between organizations, so the first reply reads it all again`을 덧붙입니다. 이전 계정은 로그인했지만 새 계정은 로그인하지 않은 MCP 서버가 있으면 ccam이 이를 알려 줍니다: `work is not signed in to MCP servers personal is: linear (sign in with /mcp there; sign-ins are not shared)`. 그룹에서 `history`를 공유하면 복사할 필요가 없어 대화가 하나로 유지되고, `/rewind`가 코드를 되돌릴 때 쓰는 체크포인트(`file-history/`)도 함께 공유됩니다.

## Claude Code 안에서

[ccam 플러그인](/docs/install#claude-code-plugin)은 훅 두 개를 추가합니다(Claude Code 2.1.78 이상).

- 턴이 사용량 한도에 걸려 멈추면 대화를 기록해 두고, 데스크톱 알림을 지원하는 터미널(iTerm2, WezTerm, Windows Terminal, kitty, Ghostty)에서 알림을 띄웁니다: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` carries on as another account.`` 훅은 Claude Code가 사용량 한도로 치는 모든 메시지를 압니다. 5시간, 주간, 일반, 모델 한도(`You've hit your session limit`, `You've hit your limit`, 모델 할당량의 `You've reached your …`), 사용 크레딧과 지출 한도, 그리고 조직이 계정을 막은 경우입니다. Claude Code가 거부된 요청을 대화에 기록했으면, 훅은 가득 찬 구간도 초기화될 때까지 기록합니다([한도 기록](/features/usage#the-limit-record)). 단, 관리자가 언제든 풀 수 있는 조직의 차단은 기록하지 않습니다. 서버의 짧은 요청 제한이나 API 키의 속도 제한은 사용량 한도가 아니므로 건드리지 않으며, 하위 에이전트가 멈춘 경우와 아무도 앞에 앉아 있지 않은 세션(`-p`, 백그라운드 세션)도 마찬가지입니다. Claude Code 2.1.234 이상은 세션 안에서 기다렸다가 5시간 구간이 초기화되면 스스로 이어 가기도 합니다. 알림은 그 대신 전환할 수 있는 선택지를 줍니다.
- `/ccam:switch [account]`는 대화를 기록해 두고, 요청을 보내기 전에 턴을 끝냅니다. `/exit`를 입력한 뒤 `ccam next --continue`를 실행하면, 지정한 계정이나 다음 계정으로 대화를 이어 갑니다. 아직 대화가 없는 세션에서는 이어 갈 것이 없다고 알려 줍니다.

메모는 세션마다 한 파일씩 `~/.local/state/ccam/handoff/<account>/<session>.json`에 남으므로, 한 계정의 두 세션이 서로의 메모를 덮어쓰지 않습니다. 메모는 12시간 동안 남아 있으며, 다음 세션이 시작되려 할 때 한 번만 쓰입니다. 실행에 실패하면 메모는 다음 시도를 위해 그대로 남습니다. 메모는 세션이 시작된 폴더에 속하며, Claude Code의 셸에서 `cd`로 옮겨 가도 마찬가지입니다. `ccam rename`은 메모 안의 이름을 고쳐 씁니다. 메모에 적힌 계정이 그 사이에 삭제되었으면 ccam은 그 사실을 알리고 평소처럼 고릅니다.

### 릴레이

`config.toml`의 맨 위에 `relay = true`를 설정하면 ccam은 자신을 Claude Code로 대체하지 않고 대화형 세션 옆에 남아 있으며, `/ccam:switch` 뒤에 그 세션의 대화를 스스로 이어 갑니다(다른 세션이 남긴 메모는 처리하지 않습니다). 훅이 Claude Code를 끝내면 ccam이 그 세션의 메모로 `ccam next --continue`를 실행합니다. 계정을 지정하지 않으면 대화의 모델까지 따져 여유가 가장 많은 계정을 고릅니다. Windows에서는 훅이 Claude Code를 끝낼 수 없으니 `/exit`를 입력하세요. 그러면 ccam이 이어 갑니다. 세션마다 기다리는 ccam 프로세스가 하나씩 생기므로 기본값은 꺼짐입니다.

이어 받은 세션은 대화를 재개해도 복원되지 않는 것을 유지합니다. `--add-dir`, `--mcp-config`, 도구 목록, `--settings`, `--plugin-dir`, 시스템 프롬프트 옵션 같은 옵션과 대화의 권한 모드, 세션의 추론 수준입니다. `--model`, `--name`, `--agent`는 Claude Code가 대화에서 직접 복원합니다. 그 밖의 것(`--worktree`, `--debug`, `--remote-control`, ccam이 모르는 옵션)은 `ccam: not carried over: …`로 알려 줍니다. 전환할 때는 대화 내역 공유에 관한 질문을 절대 하지 않으며, 전환하는 동안 입력한 내용은 버려집니다.

릴레이는 Claude Code가 `/ccam:switch`가 요청한 방식으로 끝났을 때만 이어 갑니다. 먼저 `/exit`를 입력했거나, 끝나는 동안 Ctrl+C를 눌렀거나, ccam이 행업이나 다른 시그널을 받았거나, 터미널이 사라졌을 때는 이어 가지 않습니다. 이때 메모는 남아 있으므로, 나중에 `ccam next --continue`로 대화를 이어 갈 수 있습니다. 이전 ccam의 릴레이로 시작한 세션이나 아무도 앞에 앉아 있지 않은 세션(`-p`, 백그라운드 세션)은 훅이 끝내지 않고, 대신 `/exit` 메시지를 보여 줍니다. `claude remote-control`은 릴레이 없이 실행됩니다.

## 한도에서 다시 실행

`ccam next`는 멈춘 대화형 세션을 위한 명령입니다. 스크립트나 예약 작업에서 실행하는 `claude -p` 같은 비대화형 실행은 사용자를 기다릴 수 없으므로, 계정 이름 앞에 `--fallback`을 붙이면 ccam이 대신 전환합니다.

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback`은 `-p`나 `--print`와 함께 쓸 때만 동작합니다. 그렇지 않으면 ccam은 멈추고 `ccam next --auto --continue`를 안내합니다.
2. 파이프로 들어온 표준 입력은 최대 10 MiB까지 한 번 읽어 두고 모든 시도에 똑같이 넘깁니다. Claude Code와 마찬가지로 입력이 시작되기를 최대 3초까지 기다리며(`--input-format stream-json`일 때는 기다리지 않습니다), 열려 있지만 아무것도 보내지 않는 파이프라면 입력 없이 실행하고 표준 오류에 한 줄을 남깁니다.
3. ccam은 `exec` 대신 `claude`를 자식 프로세스로 시작하고, 표준 출력은 붙잡아 두며, 표준 오류는 나오는 대로 그대로 전달합니다.
4. `claude`가 0으로 종료하거나 다른 이유로 실패하면, 붙잡아 둔 출력과 종료 코드를 바꾸지 않고 그대로 전달합니다. 시그널로 끝난 실행(128보다 큰 종료 코드)도 마찬가지입니다. 멈추라는 요청을 받은 것이므로 다시 실행하지 않습니다.
5. Claude Code의 메시지가 사용량 한도를 알릴 때(`--output-format json`이나 `stream-json` 출력 안에 있는 경우를 포함해 [훅](#inside-claude-code)이 아는 사용량 한도 메시지 중 하나), 또는 계정의 사용량에서 5시간, 주간, `--model` 구간이 100%로 나올 때(한도 기록 포함) 그 실패를 사용량 한도로 봅니다.
6. 그러면 그 시도의 출력은 버리고, 로그인된 다른 계정으로 같은 명령을 실행합니다. 5시간 사용량이 가장 낮은 계정부터, 그다음 사용량을 조회할 수 없는 계정을 번호 순서대로 시도합니다. 100%인 계정은 건너뛰고, API 키로 로그인한 계정은 절대 시도하지 않으며, 계정마다 최대 한 번만 시도합니다.
7. 인자에 `-c`, `--continue`, `-r`, `--resume`이 있으면 대화를 찾을 수 있도록 같은 대화 내역 그룹의 계정만 시도합니다.
8. `--auto`와 마찬가지로 같은 조직의 계정만 시도합니다(`cross_org = true`로 이 제한을 풀 수 있습니다).

다음 계정은 `--auto`와 같은 규칙(`auto_max_usage`, `auto_strategy`, `--model` 구간)으로 고릅니다. 남은 계정이 없으면 ccam은 마지막 시도의 출력을 표시하고, 여유가 남은 로그인 계정이 없다고 알린 뒤, 종료 코드 75로 종료합니다. 계정 이름 앞에 `--wait`를 주면(`ccam --fallback --wait=2h work -p …`) 대신 첫 번째 계정이 초기화될 때까지 기다린 뒤 모든 계정을 다시 시도할 수 있습니다. 항상 먼저 한 번은 기다리므로, 방금 한도에 도달한 계정을 곧바로 다시 실행하지는 않습니다. 이 시간은 기다림 하나하나가 아니라 실행 전체의 기다림을 합한 상한입니다. 띄어 쓰면(`--wait 30m`) `ccam next`와 같은 힌트와 함께 거부합니다. 다른 실행과 마찬가지로 Claude Code 인자 앞에 `--`를 붙여도 됩니다(`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …`는 첫 번째 시도의 명령을 출력합니다.

대화형 세션은 자동으로 전환하지 않습니다. 대신 `ccam statusline install`을 해 두면, 5시간이나 주간 사용량이 90%에 도달했을 때 상태 표시줄 끝에 `→ ccam next --auto --continue`가 붙습니다([사용량 한도](/features/usage)).

## 안전장치와 예외 상황

- 로그인되지 않은 계정은 절대 선택되지 않습니다. 로그인된 다른 계정이 없으면 일반 `ccam next`는 `no logged-in account other than <name>` 오류를 내고 멈춥니다.
- `--continue`를 썼는데 이 폴더에 대화가 없으면 `no account has a conversation in <folder> to continue`가 출력됩니다.
- 조직 규칙 때문에 고를 계정이 남지 않으면 `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``가 출력됩니다.
- `--auto`를 썼는데 로그인된 계정이 없으면 `no logged-in account to pick`가 출력됩니다. 모든 계정이 가득 찼거나 사용량을 알 수 없으면 `` no account with known usage and headroom (see `ccam usage`) ``가 출력됩니다. 둘 다 종료 코드 75로 끝납니다.
- `--max-usage`, `--strategy`, `--wait`, `--cross-org`는 `--auto`와 함께만 쓸 수 있으며, 1-100 밖의 값이나 다른 전략은 거부됩니다.
- 포크는 대화만 복사할 뿐 계정의 다른 상태는 복사하지 않습니다. 대화 내역을 공유하지 않으면 원래 계정의 `file-history/`에 있는 체크포인트는 따라오지 않습니다.
- `--auto`는 `--dry-run`을 붙여도 네트워크 요청을 보냅니다. 읽기만 하며 토큰 파일에는 절대 쓰지 않습니다.
- 선택된 계정은 `ccam <name>`과 같은 경로로 시작하므로 같은 검사가 적용됩니다. 같은 대화 내역을 공유하는 다른 Claude Code가 이 폴더에서 이미 실행 중일 때 나오는 경고도 마찬가지입니다.
- `last-used`는 Claude Code가 실제로 시작될 때만 기록됩니다. `--dry-run`은 이 파일을 건드리지 않으므로, 시험 실행으로 순환 순서가 넘어가지 않습니다.
- `--fallback`은 시도가 끝날 때까지 표준 출력을 붙잡아 두므로, `--output-format stream-json` 같은 스트리밍 출력도 마지막에 한꺼번에 나옵니다. `stream-json` 입력을 스트리밍하면서 응답을 기다린 뒤에 더 보내는 클라이언트는 `--fallback`을 쓸 수 없습니다. 첫 번째 시도가 한도에 도달하기 전에 바꾼 파일은 그대로 남고, 다음 시도는 프롬프트를 처음부터 다시 시작합니다.

## 관련 문서

- [사용량 한도](/features/usage)
- [사용량 한도와 계정 전환](/docs/usage-limits)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [사용량 한도 대응 가이드](/guides/usage-limit-workflow)
