---
title: 계정 관리
description: ccam add, ccam login, ccam list, ccam rename, ccam remove, ccam import로 Claude Code 계정을 추가, 로그인, 확인, 이름 변경, 삭제하고 가져옵니다.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

ccam 계정은 각각 자체 로그인 정보를 가진 별도의 Claude Code 설정 폴더입니다. 다음 명령으로 계정을 만들고, 로그인하고, 목록을 보고, 이름을 바꾸고, 삭제하며, 이미 있는 폴더를 등록할 수 있습니다.

## 사용법

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

명령이 `<name|number>`를 받는 곳에는 계정 이름이나 슬롯 번호(`ccam list`의 `#` 열)를 줄 수 있습니다.

터미널에서 옵션 없이 `ccam add`를 실행하면 설정 마법사가 열립니다. 이름과 로그인 방식(지금 Claude 로그인, API 키, 나중에 Claude 로그인)을 묻고, 들어갈 그룹을 `default`가 선택된 상태로 보여 준 다음(아직 그룹이 없으면 `default` 그룹이 공유할 항목), `Create a shortcut command for work? [y/N]`를 묻고, 마지막으로 확인할 요약을 보여 줍니다. `Esc`를 누르면 한 단계 뒤로 갑니다. `--group`, `--no-group`, `--shortcut`, `--no-login` 중 하나라도 주거나 터미널 대신 파이프로 실행하면 마법사를 건너뜁니다. 이때는 이름이 필요하며, `--shortcut`으로 이름을 지정하지 않으면 단축 명령을 만들지 않습니다.

`--group`을 주지 않으면 새 계정은 `default` 그룹에 들어갑니다. 첫 번째 계정이 이 그룹을 만들며, 그룹은 MCP 서버, 설정, `CLAUDE.md`, 확장 기능을 공유합니다. [그룹으로 공유하기](/features/share)를 참고하세요.

`ccam list`는 슬롯, 이름, 단축 명령, 이메일, 그룹, 마지막 사용 시각을 보여 주며, 기본 계정 앞에는 `*`가 붙습니다.

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

`--json`을 붙이면 각 계정이 `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used`(한 번도 시작하지 않은 계정에서는 생략), `groups`, `share`를 가진 객체로 출력됩니다. `share`는 각 항목(`mcp`, `settings`, `memory`, `extensions`, `history`)을 그 항목을 받는 그룹이나 `"local"`에 대응시킵니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `ccam add --no-login` | 로그인하지 않고 계정을 만듭니다. 나중에 `ccam login <name>`으로 로그인합니다 | 꺼짐 |
| `ccam add --group <group>` | `default` 대신 이 그룹에 들어갑니다. 여러 번 지정할 수 있으며, `default`가 없으면 만듭니다 | `default` |
| `ccam add --no-group` | 어느 그룹에도 들어가지 않습니다. 계정은 아무것도 공유하지 않습니다 | 꺼짐 |
| `ccam add --shortcut <command>` | `~/.local/bin`에 이 단축 명령도 만듭니다 | 없음 |
| `ccam add --api --key-from <source>` | 로그인 대신 API 키를 쓰는 계정. [API 키 계정](/features/api) 참고 | 꺼짐 |
| `ccam list --json` | 계정을 JSON으로 출력합니다 | 꺼짐 |
| `ccam remove -y`, `--yes` | 묻지 않고 삭제합니다 | 꺼짐 |
| `ccam import --name <name>` | 지정한 폴더의 계정 이름 | 폴더 이름에서 만듦 |
| `ccam import -y`, `--yes` | 찾은 폴더를 모두 묻지 않고 가져옵니다 | 꺼짐 |
| `--dry-run` | 아무것도 바꾸지 않는 `ccam list`와 `ccam completion`에서 지원합니다. `add`, `login`, `logout`, `use`, `rename`, `remove`, `import`는 `` --dry-run is not supported by `ccam add`; nothing was changed ``(해당 명령 이름이 들어감)를 출력하고 거부합니다 | 꺼짐 |

`ccam login`, `ccam logout`, `ccam use`, `ccam rename`에는 자체 옵션이 없습니다. `ccam completion`은 `bash`, `zsh`, `powershell`을 받으며, 그 밖의 값은 `only bash, zsh and powershell are supported`로 실패합니다.

계정 이름은 `a-z`, `0-9`, `-`, `_`로 이루어진 1~32자이며, 문자나 숫자로 시작해야 합니다. `ccam 2`는 슬롯 2를 뜻하므로 숫자로만 된 이름은 쓸 수 없습니다. ccam 명령 이름(`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`)도 쓸 수 없으므로, `ccam <name>`이 모호해지는 일은 없습니다. 이 목록은 `internal/config`의 `config.Reserved`입니다.

## 동작 방식

`ccam add work --shortcut ccwork`는 `addProfile`(`internal/cli/profiles.go`)에서 다음 단계를 실행합니다.

<!-- motion:accounts-add -->

1. 무엇이든 만들기 전에 `config.ValidateName`이 이름을, `shim.Check`가 단축 명령을 검사합니다.
2. `profile.Add`가 `0700` 권한으로 `~/.local/share/ccam/profiles/work/`를 만들고, 비어 있는 가장 낮은 슬롯 번호(`Config.NextSlot`)를 계정에 줍니다. 그래서 삭제한 계정의 번호가 다시 쓰입니다. 첫 번째 계정은 기본 계정이 됩니다.
3. 계정이 자신의 그룹에 들어가고(`--group`이나 `--no-group`이 따로 정하지 않으면 `default`, 첫 번째 계정은 `default`를 만듭니다), `applyShare`가 `share.Build`와 `Apply`로 새 빈 폴더를 그 그룹에 링크합니다(예: `settings.json`을 `groups/default/settings.json`에). 이 단계가 실패하면 새 폴더를 삭제합니다. 폴더가 남아 있으면 이후의 모든 재시도가 `folder already exists`로 막히기 때문입니다.
4. 계정을 `fsutil.AtomicWrite`로 `~/.config/ccam/config.toml`에 저장합니다. 임시 파일에 쓰고 `fsync`한 뒤 `rename`하므로, 도중에 비정상 종료되어도 설정이 절반만 남지 않습니다.
5. `runLogin`이 `CLAUDE_CONFIG_DIR`를 계정 폴더로 설정하고 `claude auth login`을 시작합니다. 로그인 페이지가 미리 채워지도록 계정의 마지막 이메일을 `--email`로 함께 넘깁니다. 로그인이 실패해도 계정은 남으며, ccam은 오류를 출력하고 나중에 `ccam login work`를 실행하라고 안내합니다.
6. `profile.ReadIdentity`가 `profiles/work/.claude.json`에서 `oauthAccount`를 꺼냅니다. 이메일, 계정 ID와 조직 ID, 조직 이름, 요금제(`claude_max`, `claude_team` …)가 들어 있습니다. 이 파일은 수 MB가 될 수 있으므로, `TopLevelValue`는 파일 전체를 해석하지 않고 그 키 하나만 찾습니다. ccam은 이 값들을 `config.toml`의 `[profile.login]`에 보관합니다.
7. `shim.Install`이 `~/.local/bin/ccwork`를 ccam 실행 파일에 링크합니다. ccam은 시작할 때 `config.toml`에서 `argv[0]`을 찾습니다.

로그인에는 몇 분이 걸릴 수 있고, 그동안 다른 ccam이 설정을 저장할 수도 있습니다. 그래서 기다리는 명령은 `env.update`로 저장합니다. 이 함수는 `config.toml`을 다시 읽고, 바꿀 내용 하나(여기서는 이메일)만 적용해 저장합니다. 같은 이유로 `ccam remove`도 확인 질문 뒤에 설정을 다시 읽습니다.

`ccam login`은 기존 계정에 대해 5단계와 6단계를 실행합니다. Claude Code는 이메일을 로그인 페이지의 힌트로만 쓰므로, 로그인 결과가 그 ccam 계정이 원래 쓰던 Claude 계정과 다를 수 있습니다. ccam은 새 로그인을 `[profile.login]`과 비교해 그 사실을 알려 줍니다(`ccam: this login is other@company.com, but work was dev@company.com before.`). 새 로그인은 그대로 두고, 되돌리는 방법을 안내합니다. `ccam logout work`를 실행한 뒤 `ccam login work`로 올바른 이메일로 로그인하면 됩니다. 다른 ccam 계정이 같은 Claude 계정에 로그인되어 있으면, 두 계정이 사용량 한도를 함께 쓴다고 경고합니다. `ccam login --expired`는 한 번도 로그인하지 않았거나 로그인(리프레시 토큰)이 만료된 모든 계정에 차례로 로그인합니다.

`ccam logout`은 계정 안에서 `claude auth logout`을 실행하고 저장된 이메일을 지웁니다. 로그인 기록은 남으므로, 다음 로그인도 이 기록과 비교합니다. `ccam use`는 `config.toml`의 `default`만 다시 씁니다.

계정들이 둘 이상의 조직에 속하면 `ccam list`에 `Organization` 열이 추가되며(개인 요금제는 `personal`로 표시), `ccam list --json`에는 계정마다 `plan`과 `org`가 들어 있습니다.

### 원격 머신에서 로그인하기

SSH로 접속했을 때는 `claude auth login`의 브라우저 쪽 과정이 실패할 수 있습니다. Claude Code는 서버의 `http://localhost:<port>/callback`에 콜백 서버를 띄워 로그인을 기다리지만, 로그인을 마치는 브라우저는 내 컴퓨터에서 실행되므로 그 주소로는 아무 데도 연결되지 않고 `ERR_CONNECTION_REFUSED` 같은 오류가 표시됩니다. 그래도 주소에는 일회용 코드가 들어 있으므로 ccam이 대신 전달할 수 있습니다.

<!-- motion:accounts-relay -->

`runLogin`은 터미널을 Claude Code에 그대로 넘기지 않습니다. Claude Code의 표준 입력을 파이프에 연결하고, `relayLogin`이 입력한 줄을 하나씩 그 파이프로 복사합니다. 그래서 Claude Code의 `Paste code here` 프롬프트에 붙여 넣은 코드도 그대로 전달됩니다. 입력한 줄이 콜백 주소이면 `callbackURL`이 검사합니다. 호스트가 `localhost`, `127.0.0.1`, `[::1]` 중 하나이고, 포트가 있고, 경로가 `/callback`이며, `code`와 `state`가 모두 있을 때만 받아들입니다. `http://`는 생략해도 됩니다. 그다음 `deliverCallback`이 브라우저가 했을 것처럼 서버에서 그 주소로 요청을 보냅니다. 제한 시간은 10초이고 리디렉션은 따라가지 않습니다. 요청은 기다리고 있던 `claude auth login`의 콜백 서버에 도착합니다(그림의 점선). Claude Code가 코드를 받아 로그인 정보를 저장하며, Linux에서는 계정 폴더의 `.credentials.json`에 저장됩니다. 붙여 넣은 줄 때문에 ccam이 다른 호스트에 접속하는 일은 절대 없습니다.

ccam은 로그인을 시작하기 전에 다음을 출력합니다.

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

성공하면 `→ Passed the sign-in to Claude Code on this machine.`를 출력합니다.

로그인이 터미널에서 입력을 기다리는 동안 ccam은 터미널의 에코와 줄 단위 입력 모드를 끄고, 입력한 글자마다 `*`를 하나씩 보여 줍니다. 그래서 붙여 넣은 코드나 콜백 주소가 화면에 나타나지 않습니다. 줄 편집은 `readLine`이 직접 처리합니다. `Enter`는 줄을 보내고, `Backspace`는 한 글자를 지우며, 빈 줄에서 `Ctrl+D`를 누르면 입력이 끝납니다. 터미널 설정은 로그인이 끝날 때 복원되며, `Ctrl+C`, `SIGTERM`, `SIGHUP`을 받았을 때도 복원됩니다. 이 신호들은 여전히 ccam을 종료합니다. 표준 입력과 표준 오류가 터미널이 아니면 줄을 들어오는 그대로 복사합니다.

### 이름 변경, 삭제, 가져오기

`ccam rename`과 `ccam remove`는 계정 폴더와, ccam이 다른 곳에 계정 이름으로 보관하는 파일을 함께 다룹니다. `nameState`가 그 경로들을 나열합니다. `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>`, 사용량 캐시 `~/.local/state/ccam/usage/<name>.json`과 `<name>.statusline.json`입니다. `ccam remove work`는 이들을 다음 순서로 처리합니다.

<!-- motion:accounts-remove -->

`ccam remove`는 먼저 계정이 실행 중이면 거부하고, 그다음 확인을 받습니다. 계정이 확장 기능을 공유하면 `RewritePluginPaths`가 그룹의 `installed_plugins.json`과 `known_marketplaces.json`을 다시 씁니다. Claude Code는 플러그인의 절대 경로를 저장하는데, 그 경로가 이 계정의 폴더를 거칠 수 있기 때문입니다. `profile.Trash`가 폴더를 `~/.local/share/ccam/trash/<name>-<time>/`으로 옮기고, 계정을 `config.toml`에서 빼고, 나중에 같은 이름으로 만든 계정이 물려받지 않도록 `nameState` 경로를 삭제한 뒤, 단축 명령 링크를 제거합니다. 되살리려면 휴지통에서 폴더를 다시 옮기세요.

`ccam rename`은 폴더 이름을 바꾸고, 폴더의 `plugins/`에 있는 플러그인 경로를 이전 폴더에서 새 폴더로 다시 쓰고, 기본 계정 지정을 그대로 옮기며, `nameState` 경로를 새 이름으로 옮깁니다. 단축 명령 이름은 그대로 유지됩니다. 플러그인 경로를 다시 쓰다 실패하면 설정을 저장하지 않으므로 폴더 이름도 원래대로 되돌립니다.

`ccam import`는 폴더를 원래 위치에 그대로 두고 `config.toml`의 `dir`로 등록합니다. 폴더를 옮기거나 복사하지 않습니다. 리프레시 토큰을 복사하면 사본끼리 서로를 로그아웃시키기 때문입니다. 폴더를 지정하지 않으면 등록되지 않은 `~/.claude`와 `~/.claude-*` 폴더를 찾아 폴더마다 하나씩 묻습니다. 이름은 폴더에서 정해지며(`~/.claude`는 `main`, `~/.claude-work`는 `work`, `~/.claude-2`는 `acct2`), 이름을 정할 수 없는 폴더는 써야 할 `--name` 명령과 함께 나열됩니다. `CLAUDE_CONFIG_DIR`가 없으면 Claude Code는 `.claude.json`을 `~/.claude` 옆, 즉 홈 폴더에 둡니다. 그래서 `~/.claude`를 가져올 때 그 안에 자체 `.claude.json`이 없으면, ccam은 `~/.claude.json`을 가리키는 `~/.claude/.claude.json` 링크를 만듭니다. 그러면 일반 `claude`와 `ccam main`이 같은 계정을 봅니다. 가져온 계정은 직접 그룹에 넣기 전까지(`ccam group join <group> <name>`) 어느 그룹에도 속하지 않으며, 단축 명령도 만들지 않습니다.

`ccam completion bash|zsh`는 계정 이름까지 완성해 주는 자동 완성 스크립트를 출력합니다.

## 안전장치와 예외 상황

- 세션이 폴더에 계속 쓰므로, 계정이 실행 중이면 `ccam rename`과 `ccam remove`는 `work is running; quit its Claude Code session first`를 출력하고 거부합니다(`refuseRunning`).
- 가져온 계정에 `ccam remove`를 실행하면 등록만 해제합니다: `✓ main unregistered (folder ... left in place)`. `~/.claude`는 휴지통으로 옮겨지지 않습니다.
- `ccam add`는 `profiles/`에 폴더가 이미 있는 이름을 거부하고 `ccam import`를 안내합니다.
- `ccam import`는 이미 등록된 폴더, 폴더가 아닌 경로, 이미 쓰이는 이름을 거부합니다.
- `ccam logout`과 `ccam login`은 계정이 실행 중인지 확인하지 않습니다.
- 로그인이 끝난 뒤처럼 전달 과정이 Claude Code에 닿지 못하면 `ccam: could not pass the sign-in to Claude Code:`와 이유를 출력합니다. HTTP 오류는 대개 코드가 만료되었다는 뜻이니 로그인을 다시 시작하세요.

## 관련 문서

- [Claude Code 시작하기](/features/launch)
- [단축 명령](/features/shortcuts)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [Claude Code 여러 계정 쓰기](/guides/multiple-claude-code-accounts)
