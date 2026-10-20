---
title: 사용량 한도
description: ccam usage는 모든 계정의 5시간 한도와 주간 한도를 한 번에 보여 주고, ccam statusline은 Claude Code 안에 계정 이름과 사용량을 표시합니다.
updated: 2026-10-20
order: 7
translation: machine
sourceHash: "ec1d1ee27d9f"
---

`ccam usage`는 각 계정이 5시간 한도와 주간 한도를 얼마나 썼는지, 각 구간이 초기화되기까지 얼마나 남았는지 보여 줍니다. 긴 작업을 시작하기 전이나 Claude Code가 한도가 가까워졌다고 알릴 때 실행해, 어느 계정에 여유가 남았는지 확인하세요. `ccam statusline`은 Claude Code의 상태 표시줄 명령입니다. Claude Code 안에 계정 이름과 사용량을 표시하고, `ccam usage`가 쓸 수 있도록 수치를 저장합니다.

## 사용법

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

터미널에서는 표가 바로 나타나고 각 행에 `fetching…`이 표시되며, 계정이 응답하는 대로 해당 행을 다시 그립니다. `As of`는 수치를 가져온 지 얼마나 되었는지 나타내고, `(cached)`는 엔드포인트가 아니라 저장된 파일에서 가져온 수치라는 뜻입니다. 마지막 줄은 [`ccam next --auto`](/features/next)가 같은 규칙으로 고를 계정을 알려 줍니다(`auto_max_usage`, `auto_strategy`, 실행 중인 세션을 따지고, `cross_org`를 설정하지 않으면 마지막으로 사용한 계정과 같은 조직의 계정만 보며, 100%인 구간이 있는 계정은 건너뛰므로 주간 한도가 찬 계정은 절대 나오지 않습니다). 이 줄은 계정이 둘 이상일 때만 나타납니다. 로그인되지 않은 계정은 `not logged in → ccam login <name>`, API 키 계정은 `no limits (API)`로 표시됩니다.

`ccam statusline install`을 실행하면 Claude Code에 `[work] · Opus · 5h 24% · 7d 41%` 같은 줄이 표시됩니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--json` | 결과를 JSON으로 출력합니다: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | 꺼짐 |
| `--watch[=N]` | Ctrl+C를 누를 때까지 `N`초마다 다시 그립니다. `--watch`만 주거나 0 이하의 값을 주면 60초입니다 | 꺼짐 |
| `-v`, `--verbose` | 모델별 주간 한도와 `source:` 줄(`live`, `cache`, `statusline`)을 표시합니다 | 꺼짐 |
| `--no-cache` | 60초 캐시를 무시하고 엔드포인트에 다시 요청합니다 | 꺼짐 |
| `--source live\|claude\|cache` | `live`는 계정의 토큰으로 엔드포인트에 요청하고, `claude`는 Claude Code가 조회하게 하며(아래 참고), `cache`는 저장된 수치만 보여 줍니다 | `usage_source`, 없으면 `live` |
| `--refresh` | 만료된 액세스 토큰을 갱신합니다. 실행 중이 아닌 계정만 대상입니다 | 꺼짐 |
| `--dry-run` | 받아들입니다. `ccam usage`는 자체 캐시 말고는 아무것도 바꾸지 않습니다. `--refresh`와 함께 쓰면 `--dry-run cannot be combined with --refresh, which rewrites login tokens`로 거부합니다 | 꺼짐 |

`--source`에 다른 값을 주면 `--source must be live, claude or cache`로 실패합니다.

### Claude Code 자체에서 사용량 받기

`config.toml`에 `usage_source = "claude"`를 설정하면(또는 `--source claude`를 주면) ccam은 로그인 토큰을 전혀 읽지 않습니다. 계정마다 계정 폴더에서 `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence`를 시작하고, 표준 입력으로 Claude Code의 `get_usage` 요청을 보냅니다(Claude Code 2.1.270 이상). Claude Code가 자체 로그인으로 사용량을 조회해 5시간, 주간, 모델별 구간을 답합니다. `--safe-mode` 덕분에 계정의 훅, MCP 서버, 플러그인은 시작되지 않으며, 세션도 저장되지 않습니다. 계정마다 Claude Code를 한 번 시작해야 하므로 약 1초가 걸려 `live`보다 느립니다. 모든 계정에 동시에 요청하며 계정마다 제한 시간은 20초이고, 60초 캐시는 평소처럼 적용됩니다. Claude Code는 이 요청을 실험적 기능으로 표시하므로, 응답이 멈추면 ccam은 저장된 수치를 보여 주고 `ccam doctor`가 이를 보고합니다. `--dry-run`은 Claude Code를 시작하지 않고(시작하면 계정 폴더에 쓰기 때문입니다) 저장된 수치를 보여 줍니다. `usage_source`에는 `live`나 `claude`를 줄 수 있으며, `cache`는 `--source`에서만 고를 수 있습니다. `ccam usage`는 인자를 받지 않습니다.

`ccam statusline`은 표준 입력으로 Claude Code의 상태 표시줄 JSON을 읽습니다. 옵션은 `--wrap <command>` 하나뿐이며, 기존 상태 표시줄 명령을 지정하면 그 출력의 첫 줄 앞에 계정 이름을 붙입니다. `ccam statusline install`은 옵션을 받지 않습니다. 두 명령 모두 `--dry-run`을 주면 `` --dry-run is not supported by `ccam statusline`; nothing was changed ``를 출력하고 거부합니다.

## 동작 방식

수치의 출처는 두 곳입니다. 하나는 각 계정 자체의 로그인 정보로 호출하는 Anthropic의 비공개 사용량 엔드포인트 `GET https://api.anthropic.com/api/oauth/usage`이고, 다른 하나는 Claude Code가 상태 표시줄 명령에 넘기는 값입니다. 관련 코드는 모두 `internal/usage`에 있습니다.

1. `cli.usageCmd`가 `config.toml`을 읽고 계정마다 `usage.Account`를 하나씩 만듭니다. 여기에는 계정 폴더 `~/.local/share/ccam/profiles/<name>/`, 세션 폴더(실행 중인지 확인용), 제자리에 가져온 계정인지 여부가 들어갑니다.
2. `Fetcher.FetchAll`이 계정마다 고루틴을 하나씩 시작합니다. 모든 고루틴이 `http.Client` 하나를 공유하므로 엔드포인트 연결이 재사용됩니다. 연결에는 최대 2초, 계정마다 전체 3초가 주어집니다. 리디렉션을 따라가면 토큰이 다른 곳으로 다시 전송되므로 클라이언트는 리디렉션을 따르지 않으며, 응답은 최대 1 MB까지만 읽습니다.
3. 계정마다 처리가 끝나는 대로 그 행을 다시 그립니다. 실패하거나 제때 응답하지 않은 계정은 저장된 수치로 대신하므로, 느린 계정 하나 때문에 다른 계정이 늦어지지 않습니다.

각 계정은 자체 고루틴에서 엔드포인트에 요청하며, 시간 초과된 계정은 대신 저장된 수치를 보여 줍니다.

<!-- motion:usage-fetch -->

`personal`과 `work`는 제때 응답해 실시간 수치를 보여 줍니다. `side`는 3초보다 오래 걸리므로, `fallback()`이 두 캐시 파일을 읽어 가장 나중에 기록된 것을 씁니다. `usage.Table()`이 행을 그리고, `usage.Best()`가 마지막 줄에 표시할 계정을 고릅니다.

계정 하나에 대해 `Fetcher.Fetch`는 다음 순서로 확인합니다.

<!-- motion:usage-cache -->

1. `--source cache`이면 바로 `fallback()`으로 갑니다.
2. `--no-cache`가 없으면, `~/.local/state/ccam/usage/<name>.json`에 있는 실시간 결과가 60초(`CacheTTL`)보다 최근이면 그대로 반환합니다. `ccam usage`를 연달아 두 번 실행해도 요청은 계정당 한 번입니다.
3. `readCredentials()`가 `.credentials.json`(또는 macOS 키체인)을 바꾸지 않고 읽으며, `claudeAiOauth` 키만 해석합니다. 구독 로그인이 없으면 `not logged in`이고, `.claude.json`에 API 키가 있으면 `no limits (API)`입니다.
4. 만료된 액세스 토큰은 쓰지 않습니다. `--refresh`가 없으면 그 계정은 만료 표시와 함께 저장된 수치를 보여 줍니다.
5. 그렇지 않으면 ccam이 엔드포인트를 호출합니다. 401이나 403은 로그인 만료로, 그 밖의 상태 코드나 읽을 수 없는 본문은 오류로 처리하며, 둘 다 저장된 수치로 대신합니다. 429는 엔드포인트에 너무 자주 묻고 있다는 뜻입니다. ccam은 `<name>.pause`에 다시 물어볼 시각(`Retry-After`, 없으면 5분, 최대 1시간)을 기록하고, 저장된 수치를 `rate limited by the usage endpoint; asking again after 15:10`과 함께 보여 주며, 그때까지는 `--no-cache`를 붙여도 그 계정의 엔드포인트를 호출하지 않습니다. 곧바로 다시 시도하면 한도가 계속 이어질 뿐입니다.
6. 정상 응답은 측정 시각과 함께 `<name>.json`에 기록합니다.

`fallback()`은 `<name>.json`과 `<name>.statusline.json`을 모두 읽고 더 최근 것을 씁니다. 저장된 수치를 보여 주기 전에 각 구간의 초기화 시각을 확인해, 수치를 저장한 뒤 초기화된 구간은 0%로 표시합니다. 이전 사용률은 더 이상 맞지 않기 때문입니다. 위 예시에서 `personal`이 초기화 시각 없이 `0%`로 표시되는 이유가 이것입니다.

엔드포인트의 응답 형식은 두 가지입니다. 새 형식에는 `session`(5시간), `weekly_all`, `weekly_scoped`(모델별) 항목이 담긴 `limits` 배열이 있습니다. 이전 형식에는 `five_hour`, `seven_day`, `seven_day_opus`, `seven_day_sonnet` 객체가 있습니다. `ParseLive`는 둘 다 읽습니다. `limits` 배열을 우선하고, 배열이 다루지 않는 부분은 이전 형식의 객체로 채웁니다. 한도가 하나도 없는 본문은 형식이 바뀐 것으로 보고하며, `ccam doctor`가 이를 알려 줍니다.

### 토큰 갱신

리프레시 토큰은 갱신할 때마다 바뀝니다. 갱신하면 새 리프레시 토큰을 받고 이전 토큰은 무효가 됩니다. 같은 계정으로 Claude Code가 실행 중일 때 ccam이 로그인을 갱신하면 그 세션은 로그아웃됩니다. 그래서 `ccam usage`는 `--refresh`를 주지 않는 한 토큰을 읽기만 하고 절대 쓰지 않습니다. 이 옵션을 주면 `Fetcher.refresh`는 다음과 같이 동작합니다.

1. `ccam import`로 제자리에 등록한 계정은 거부합니다. ccam이 모르는 사이에 일반 `claude`가 그 로그인을 쓸 수 있기 때문입니다.
2. macOS 키체인에 저장된 로그인은 거부합니다.
3. 잠금 파일 `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock`(이 변수가 없으면 `~/.local/state/ccam/locks/` 아래)을 잡아, 두 ccam 프로세스가 같은 토큰을 두 번 갱신하지 않게 합니다.
4. 잠금을 잡은 상태에서 그 계정의 세션이 실행 중이면 거부하고, 다른 프로세스가 이미 토큰을 갱신했으면 요청을 건너뜁니다.
5. 리프레시 토큰을 토큰 엔드포인트에 보내고, 새 토큰을 `fsutil.AtomicWrite`로 `.credentials.json`에 다시 기록합니다. 파일의 다른 키는 모두 그대로 둡니다.

### 상태 표시줄 수집기

Claude Code는 상태 표시줄 명령을 실행하면서 JSON 문서를 파이프로 넘깁니다. ccam은 여기서 모델 이름과 `rate_limits`를 읽습니다.

<!-- motion:usage-statusline -->

1. `runStatusline`은 ccam이 Claude Code를 시작할 때 설정하는 `CCAM_PROFILE`로 계정을 찾고, 없으면 `CLAUDE_CONFIG_DIR`를 계정 폴더와 대조해 찾습니다.
2. `ParseStatusLine`이 `rate_limits.five_hour`와 `rate_limits.seven_day`를 읽습니다. API 키 계정이거나 세션의 첫 응답 전이면 JSON에 사용량 한도가 없으며, 이때는 아무것도 저장하지 않습니다.
3. `SaveStatusLine`이 `~/.local/state/ccam/usage/<name>.statusline.json`을 기록합니다. 엔드포인트가 응답하지 않으면 `ccam usage`가 이 파일을 씁니다.
4. ccam은 `[name] · model · 5h N% · 7d N%`를 출력합니다. `--wrap`을 주면 같은 입력으로 `sh -c`를 통해 지정한 명령을 실행하고, 대신 `[name]` 뒤에 그 명령의 첫 줄을 붙여 출력합니다. Windows에서는 Claude Code와 같은 방식으로 명령을 실행합니다. Git Bash가 설치되어 있으면 Git Bash로, 없으면 PowerShell로 실행합니다.
5. 5시간이나 주간 사용량이 90% 이상이고 다른 계정이 있으면, 줄 끝에 `→ ccam next --auto --continue`가 붙습니다. 예: `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

감싼 명령에는 2초가 주어집니다. 그 시간이 지나면 명령을 종료하고 계정 이름만 표시하므로, 느린 명령 때문에 Claude Code 상태 표시줄이 멈추지 않습니다.

`ccam statusline install`은 각 계정의 `settings.json`에 `statusLine`을 기록합니다. 계정에 이미 상태 표시줄 명령이 있으면 ccam은 그 명령을 `--wrap '<old command>'`로 보존합니다. Windows에서는 ccam 경로를 슬래시와 홈 폴더를 뜻하는 `~`로 기록하며(`~/.local/bin/ccam.exe statusline`), Git Bash와 PowerShell 모두 이 경로를 이해합니다. 그룹을 통해 공유되는 `settings.json`은 링크를 거쳐 한 번만 기록되므로, 그룹의 모든 계정에 적용됩니다. 다시 실행하면 이미 ccam을 쓰는 파일에는 `(already set)`를 출력합니다.

## 안전장치와 예외 상황

- `ccam usage`는 `--refresh`를 주지 않는 한 토큰 파일을 절대 바꾸지 않습니다. 로그인이 만료되었으면 `ccam <name>`으로 그 계정을 한 번 실행하세요. Claude Code가 갱신합니다.
- `--refresh`는 실행 중인 계정, 가져온 계정, 키체인 로그인을 절대 건드리지 않습니다. 갱신이 거부되면 저장된 수치와 그 이유를 보여 줍니다.
- 사용량 엔드포인트는 공개 문서가 없고, 바뀔 수도 있습니다. 바뀌더라도 저장된 수치는 계속 볼 수 있으며, `ccam doctor`가 형식 변경을 보고합니다.
- 저장된 수치도 없고 응답도 없으면 그 행에 `lookup failed:`와 이유를 표시합니다(예: `--source cache`일 때 `no saved usage`).
- 캐시 파일에는 사용률과 초기화 시각만 들어 있으며, 토큰은 절대 들어 있지 않습니다.

## 관련 문서

- [사용량 한도와 계정 전환](/docs/usage-limits)
- [다음 계정으로 전환](/features/next)
- [사용량 한도 대응 가이드](/guides/usage-limit-workflow)
- [명령어](/docs/commands)
