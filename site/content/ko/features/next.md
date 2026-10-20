---
title: 다음 계정으로 전환
description: ccam next는 마지막에 쓴 다음 계정으로 Claude Code를 엽니다. --auto는 가장 여유 있는 계정을 쓰고, --continue는 이어 가며, --fallback은 -p 작업을 재실행합니다.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 7bc55c2dd520
---

`ccam next`는 계정 이름을 지정하지 않아도 다른 계정으로 Claude Code를 엽니다. 계정이 사용량 한도에 도달했을 때 쓰면 됩니다. 그냥 `ccam next`를 실행하면 마지막으로 사용한 계정의 다음 계정을, `--auto`를 붙이면 5시간 구간의 여유가 가장 많은 계정을 고릅니다. `--continue`를 붙이면 마지막 대화를 새 계정에서 이어 갑니다. `claude -p`를 실행하는 스크립트에서는 `ccam --fallback`이 이 전환을 알아서 합니다([한도에서 다시 실행](#run-again-at-the-limit) 참고).

## 사용법

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam은 표준 오류에 전환 내용을 `→ <from> → <to>` 형식으로 출력한 다음, `ccam <name>`과 같은 방식으로 새 계정의 Claude Code를 시작합니다. `<from>`은 마지막으로 사용한 계정입니다. `--` 뒤의 인자는 그대로 `claude`에 전달됩니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--auto` | 로그인된 계정 중 5시간 사용량이 가장 적은 계정을 고릅니다. 사용량이 100%인 계정은 건너뜁니다 | 꺼짐 |
| `-c`, `--continue` | 마지막으로 사용한 계정과 대화 내역을 공유하는 계정만 고려하고, `claude --continue`로 시작합니다 | 꺼짐 |
| `-- <claude args>` | ccam 자체 인자 뒤에 `claude`로 넘길 인자 | 없음 |
| `--dry-run` | 지원합니다. Claude Code를 시작하거나 실행 기록을 남기지 않고, `→ <from> → <to>`와 실행할 명령줄을 출력합니다 | 꺼짐 |

`--auto`와 `--continue`는 함께 쓸 수 있습니다. `--continue`를 쓰면 `--continue`가 직접 넘긴 인자보다 앞에 놓입니다.

## 동작 방식

ccam으로 실행할 때마다 그 시각이 `~/.local/state/ccam/last-used/<name>`에 기록됩니다. `ccam next`는 이 파일들을 읽고, 어떤 계정이 로그인되어 있는지 확인한 뒤 계정을 고릅니다.

1. `cli.candidates`가 `config.toml`의 계정마다 `next.Candidate`를 하나씩 만듭니다. 여기에는 계정, 마지막 실행 시각(`launcher.LastUsed`), 로그인 여부(`creds.Has`: 읽을 수 있는 `.credentials.json`, macOS에서는 키체인 항목)가 들어갑니다.
2. `--auto`가 없으면 `next.Rotate`가 계정을 슬롯 번호순으로 정렬하고 가장 최근에 실행된 계정을 찾습니다. 그 계정이 지금 전환해 떠나는 계정입니다. 아직 실행된 계정이 없으면 가장 낮은 슬롯이 기준이 됩니다.
3. 그 지점부터 슬롯 순서대로 앞으로 진행하며, 마지막 슬롯 다음에는 첫 슬롯으로 돌아갑니다. 처음 만나는 로그인된 계정을 반환합니다.
4. `launcher.LaunchWith`가 고른 계정을 `ccam <name>`과 똑같이 준비하고 시작합니다. 계정의 링크를 확인하고, MCP 공유를 적용하고, 세션을 등록한 뒤 `syscall.Exec`로 ccam을 `claude`로 바꿉니다.

탐색은 마지막으로 사용한 계정에서 시작해, 그 뒤에 처음 나오는 로그인된 계정에서 멈춥니다.

<!-- motion:next-rotate -->

`work`가 가장 최근에 실행되었으므로 탐색은 그 다음부터 시작합니다. `side`는 로그인되어 있지 않아 건너뛰고, `lab`이 선택됩니다. `lab`도 로그아웃되어 있었다면 탐색은 처음으로 돌아가 `personal`로 갑니다(그림의 `% len`). 순환 방식은 출발한 계정을 절대 반환하지 않습니다. 로그인된 다른 계정이 없으면 오류를 내고 멈춥니다.

### --auto

`--auto`는 슬롯 순서가 아니라 사용량으로 결정합니다. `cli.pickAuto`는 로그인된 계정만 남기고, `ccam usage`가 쓰는 것과 같은 `Fetcher.FetchAll`로 모든 계정의 사용량을 한꺼번에 요청한 뒤 `usage.Best()`에게 선택을 맡깁니다.

<!-- motion:next-auto -->

`usage.Best()`는 5시간 사용률이 가장 낮은 계정을 고르며, 100%이거나 5시간 사용량을 알 수 없는 계정은 건너뜁니다. 조회는 `ccam usage`의 규칙을 따릅니다. 가져온 지 60초가 안 된 수치는 다시 쓰고, 계정마다 3초까지 기다리며, 응답하지 않거나 로그인이 만료된 계정은 저장된 수치로 판단합니다. 저장된 구간의 초기화 시각이 지났으면 0%로 계산하므로, 그 사이에 구간이 초기화된 계정이 밀려나지 않습니다. `--auto`는 토큰을 절대 갱신하지 않습니다. 일반 순환과 달리 `--auto`는 마지막으로 사용한 계정의 여유가 여전히 가장 많으면 그 계정을 고를 수도 있습니다.

### --continue

Claude Code의 `--continue`는 현재 폴더의 가장 최근 대화를 여는데, 새 계정이 그 대화를 볼 수 있어야만 동작합니다. 그래서 `--continue`를 쓰면 ccam은 마지막으로 사용한 계정과 같은 `history` 그룹에 있는 계정만 고려하고, 직접 넘긴 인자 앞에 `--continue`를 붙입니다. 대화 내역이 `local`인 계정은 어떤 계정과도 공유하지 않으므로, 그런 계정에서 `--continue`를 쓰면 항상 실패합니다. 먼저 `ccam share <name> --history on`으로 공유를 켜세요([공유](/docs/sharing) 참고).

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
2. 표준 입력은 한 번 읽어 두고 모든 시도에 똑같이 넘깁니다.
3. ccam은 `exec` 대신 `claude`를 자식 프로세스로 시작하고, 표준 출력은 붙잡아 두며, 표준 오류는 나오는 대로 그대로 전달합니다.
4. `claude`가 0으로 종료하거나 다른 이유로 실패하면, 붙잡아 둔 출력과 종료 코드를 바꾸지 않고 그대로 전달합니다.
5. Claude Code의 메시지가 사용량 한도를 알릴 때(`You've hit your limit`, `usage limit reached`), 또는 계정의 사용량 조회에서 5시간이나 주간 구간이 100%로 나올 때 그 실패를 사용량 한도로 봅니다.
6. 그러면 그 시도의 출력은 버리고, 로그인된 다른 계정으로 같은 명령을 실행합니다. 5시간 사용량이 가장 낮은 계정부터, 그다음 사용량을 조회할 수 없는 계정을 번호 순서대로 시도합니다. 100%인 계정은 건너뛰고, 계정마다 최대 한 번만 시도합니다.
7. 인자에 `-c`, `--continue`, `-r`, `--resume`이 있으면 대화를 찾을 수 있도록 같은 대화 내역 그룹의 계정만 시도합니다.

남은 계정이 없으면 ccam은 마지막 시도의 출력을 표시하고, 여유가 남은 로그인 계정이 없다고 알린 뒤, 그 시도의 종료 코드로 종료합니다. `ccam --dry-run --fallback work -p …`는 첫 번째 시도의 명령을 출력합니다.

대화형 세션은 자동으로 전환하지 않습니다. 대신 `ccam statusline install`을 해 두면, 5시간이나 주간 사용량이 90%에 도달했을 때 상태 표시줄 끝에 `→ ccam next --auto --continue`가 붙습니다([사용량 한도](/features/usage)).

## 안전장치와 예외 상황

- 로그인되지 않은 계정은 절대 선택되지 않습니다. 로그인된 다른 계정이 없으면 일반 `ccam next`는 `no logged-in account other than <name>` 오류를 내고 멈춥니다.
- `--continue`를 썼는데 대화 내역을 공유하는 계정이 없으면 `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``가 출력됩니다.
- `--auto`를 썼는데 로그인된 계정이 없으면 `no logged-in account to pick`가 출력됩니다. 모든 계정이 가득 찼거나 사용량을 알 수 없으면 `` no account with known usage and headroom (see `ccam usage`) ``가 출력됩니다.
- `--auto`는 `--dry-run`을 붙여도 네트워크 요청을 보냅니다. 읽기만 하며 토큰 파일에는 절대 쓰지 않습니다.
- 선택된 계정은 `ccam <name>`과 같은 경로로 시작하므로 같은 검사가 적용됩니다. 같은 대화 내역을 공유하는 다른 Claude Code가 이 폴더에서 이미 실행 중일 때 나오는 경고도 마찬가지입니다.
- `last-used`는 Claude Code가 실제로 시작될 때만 기록됩니다. `--dry-run`은 이 파일을 건드리지 않으므로, 시험 실행으로 순환 순서가 넘어가지 않습니다.
- `--fallback`은 시도가 끝날 때까지 표준 출력을 붙잡아 두므로, `--output-format stream-json` 같은 스트리밍 출력도 마지막에 한꺼번에 나옵니다. 첫 번째 시도가 한도에 도달하기 전에 바꾼 파일은 그대로 남고, 다음 시도는 프롬프트를 처음부터 다시 시작합니다.

## 관련 문서

- [사용량 한도](/features/usage)
- [사용량 한도와 계정 전환](/docs/usage-limits)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [사용량 한도 대응 가이드](/guides/usage-limit-workflow)
