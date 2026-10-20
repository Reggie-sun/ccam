---
title: Claude Code 시작하기
description: ccam <name|number>, ccam run, ccam 선택 화면으로 원하는 계정의 Claude Code를 시작합니다. 인자는 그대로 claude에 전달되며, --dry-run은 실행할 명령을 보여 줍니다.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: e6b5a8d4df86
---

`ccam work`는 `work` 계정으로 Claude Code를 시작하며, 이름 뒤의 모든 것은 그대로 `claude`에 전달됩니다. 터미널, 스크립트, 파이프로 입력을 넘길 때 등 `claude`를 입력하던 곳이면 어디서나 쓸 수 있습니다.

## 사용법

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

계정 이름이나 `ccam list`의 번호를 지정하세요. ccam은 이름 뒤의 인자를 해석하지 않으므로, `ccam work --help`는 Claude Code의 도움말을 보여 줍니다. 이름 뒤의 `--`는 넣어도 되고 생략해도 됩니다. 표준 입력, 종료 코드, 시그널은 `claude`를 직접 실행할 때와 똑같이 동작합니다.

첫 번째 인자가 `-`로 시작하고 ccam 자체 옵션이 아니면, ccam은 모든 인자를 넘겨 이름 없이 계정을 시작합니다. 이때는 [`ccam folder`](/features/folder)로 현재 폴더에 지정한 계정을, 없으면 기본 계정을 씁니다. 기본 계정은 `ccam use`로 지정한 계정이며, 지정하지 않았다면 번호가 가장 낮은 계정입니다.

### 선택 화면

표준 입력과 출력이 터미널이면 인자 없이 `ccam`을 실행할 때 선택 화면이 열리고, 터미널이 아니면 `ccam list`를 출력합니다. 선택 화면은 계정마다 번호, 이메일, 5시간 사용량, 마지막 사용 시각, 공유 항목을 보여 주며, 사용량은 계정이 응답하는 대로 채워집니다.

| 키 | 동작 |
| --- | --- |
| `Enter` | 선택한 계정으로 Claude Code를 시작합니다 |
| `1`~`9` | 해당 번호의 계정을 시작합니다 |
| `Up`/`Down`, `k`/`j` | 선택을 옮깁니다 |
| `a` | 선택 화면을 닫고 `ccam add` 설정 마법사를 실행합니다 |
| `s` | 선택한 계정의 공유 항목을 편집합니다(`Space`나 `x`로 전환, `Enter`로 적용, `Esc`로 돌아가기) |
| `d` | 선택한 계정을 기본 계정으로 지정합니다 |
| `r` | 짧은 캐시를 무시하고 사용량을 다시 가져옵니다 |
| `q`, `Esc`, `Ctrl+C` | 종료합니다 |

선택 화면은 동작하기 전에 닫히므로, Claude Code가 시작되기 전에 터미널이 원래대로 돌아옵니다. `d`를 누르거나 공유 설정을 바꾼 뒤에는 다시 열립니다. 공유 설정을 바꿀 때는 `ccam share`와 마찬가지로 계획을 보여 주고 먼저 확인을 받습니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--dry-run` | ccam이 실행할 환경 변수와 명령을 출력하고 종료합니다 | 꺼짐 |
| `--fallback` | `-p`/`--print` 실행용: 계정이 사용량 한도에 도달하면 로그인된 계정 중 여유가 가장 많은 계정으로 같은 명령을 다시 실행합니다([자세히](/features/next#run-again-at-the-limit)) | 꺼짐 |
| `-v`, `--version` | `ccam <version>`을 출력합니다(첫 번째 인자일 때만) | |
| `-h`, `--help` | ccam 도움말. `ccam --help` 또는 `ccam run --help`로 씁니다 | |

인자 규칙은 다음과 같습니다.

- ccam은 계정 이름 앞의 자체 옵션만 읽습니다. 이름 뒤의 모든 인자는 `--dry-run`과 `--help`를 포함해 Claude Code의 것입니다.
- `--dry-run`은 `ccam --dry-run work`, `ccam run --dry-run work`, `ccam --dry-run`(이름 없이 시작하는 계정)으로 쓸 수 있습니다. 거부하지 않고 지원합니다. `--dry-run`과 `--fallback`은 순서를 바꿔 써도 됩니다.
- `ccam run <name|number> [claude args]`는 `ccam <name>`의 긴 형식입니다. 이름 없이 `ccam run`을 실행하면 폴더에 지정된 계정이나 기본 계정을 씁니다.
- 계정 이름에는 `add`, `list`, `run`, `usage` 같은 ccam 명령 이름을 쓸 수 없으므로, `ccam <word>`가 모호해지는 일은 없습니다.

시험 실행은 셸에, Windows에서는 PowerShell에 그대로 붙여 넣을 수 있는 한 줄을 출력합니다.

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## 동작 방식

실행할 때마다 명령줄부터 `exec`까지 같은 단계를 거칩니다.

<!-- motion:launch -->

각 상자는 해당 단계가 실행될 때 강조되며, 아래 줄은 Claude Code가 시작될 때 받는 환경 변수와 인자입니다.

1. `cmd/ccam/main.go`가 `app.Main`을 호출합니다. ccam이 다른 이름으로 시작되었다면 [단축 명령](/features/shortcuts)이므로 바로 실행기로 넘어갑니다.
2. `app.parse`는 cobra가 관여하기 전에 `os.Args`를 살펴봅니다. `ccam`만 입력한 경우(선택 화면), ccam 명령 이름(`config.Reserved`), `-h`/`--help`, `-v`/`--version`만 `cli.Execute`로 갑니다. 나머지는 모두 실행입니다.
3. `launcher.Launch`가 `~/.config/ccam/config.toml`을 읽습니다. 해석하는 파일은 이것 하나뿐입니다.
4. `launcher.Prepare`가 이름이나 번호로 계정을 찾고(이름이 없으면 `HereProfile`이 작업 폴더를 `config.toml`의 `[[folder]]` 항목과 문자열로 비교해 폴더의 계정을 고릅니다), 맨 앞의 `--` 하나를 제거한 뒤, `ResolveClaude`로 `claude` 실행 파일을 찾습니다. PATH 검색 결과는 파일의 수정 시각과 함께 `~/.local/state/ccam/claude-path`에 캐시되므로, ccam은 실행할 때마다 `PATH`를 검색하지 않습니다. 파일이 바뀌었거나 없어졌으면 다시 검색합니다. `config.toml`에 `claude_bin` 경로가 있으면 그대로 씁니다. `claude`가 ccam 자신이면 무한 반복이 되므로, inode로 비교해 거부합니다.
5. `checkLinks`가 계정 폴더 `~/.local/share/ccam/profiles/<name>/`를 공유 설정과 비교합니다. 이때 쓰는 `share.Check`는 `lstat`과 `readlink`만 호출하므로 항목당 시스템 호출은 한두 번입니다. 없는 링크는 그 자리에서 만듭니다. 데이터를 잃을 수 있는 문제가 있으면 실행을 멈춥니다.
6. `applyMCP`가 공유 MCP 서버를 추가합니다(아래 참고).
7. `session.Register`가 세션을 기록합니다. ccam은 PID를 이름으로 하고 프로세스 시작 시각을 담은 표시 파일을 `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>`(`XDG_RUNTIME_DIR`가 없으면 `~/.local/state/ccam/locks/`)에 기록합니다. 대화 내역을 공유하는 계정이라면 `historyCheck`가 먼저 `history-<group>-<hash>`에서 실행 중인 세션을 찾고(해시는 현재 폴더의 해시), 그곳에도 등록합니다. `exec`는 PID를 유지하므로, 이후 표시 파일은 Claude Code 프로세스 자체를 가리킵니다. 잠금을 잡아 두지 않고 표시 파일을 쓰는 이유는, Claude Code가 시작하는 프로세스가 잠금을 물려받아 세션이 끝난 뒤에도 쥐고 있게 되기 때문입니다.
8. `Plan.Exec`가 터미널 제목을 `claude · <name>`으로 설정하고(표준 오류가 터미널일 때), `fsync` 없이 실행 시각을 `~/.local/state/ccam/last-used/<name>`에 기록한 뒤 `syscall.Exec`를 호출합니다.

Claude Code는 `CLAUDE_CONFIG_DIR`가 계정 폴더로, `CCAM_PROFILE`이 계정 이름으로 설정된 상태에서 실행됩니다. ccam은 물려받은 이 변수들의 값을 먼저 지웁니다. macOS에서 `ccam import`로 등록한 `~/.claude` 폴더라면 `CLAUDE_SECURESTORAGE_CONFIG_DIR`도 빈 값으로 설정해, 그 계정과 일반 `claude`가 같은 키체인 로그인을 쓰게 합니다.

### exec를 쓰는 이유

`syscall.Exec`는 ccam 프로세스를 `claude`로 바꿉니다. 그래서 ccam이 메모리에 남지 않고, 키 입력, 시그널, 종료 코드가 모두 Claude Code 자체의 것이 되며, 인자는 셸 문자열이 아닌 배열로 유지됩니다.

Windows에는 `exec`가 없습니다. 그래서 Windows에서는 `Plan.Exec`가 같은 환경 변수, 인자, 콘솔로 `claude.exe`를 자식 프로세스로 시작하고, 끝날 때까지 기다린 뒤 그 종료 코드로 종료합니다. Ctrl+C는 콘솔에서 Claude Code로 바로 전달되므로 ccam은 이를 무시합니다. 세션 표시 파일에는 ccam의 PID가 들어가는데, ccam은 Claude Code와 정확히 같은 시간 동안 살아 있으며, 표시 파일은 Claude Code가 끝나면 지워집니다. 기다리는 동안 ccam은 다른 일을 하지 않습니다. `claude mcp add`로 추가한 서버를 알아차리는 것처럼 세션 이후의 작업은 다음 실행 때 처리합니다. ccam은 `.claude.json`의 수정 시각을 `~/.local/state/ccam/mcp/<name>/seen`과 비교하고, 바뀌었을 때만 파일을 해석합니다.

### 명령줄이 처리되는 경로

세 가지 명령줄이 `internal/app`에서 세 가지 경로로 나뉩니다.

<!-- motion:launch-route -->

단축 명령은 `argv[0]`로, `ccam <name>`은 첫 번째 인자로 구분하며, 둘 다 `syscall.Exec`로 끝납니다. `ccam usage` 같은 ccam 명령만 cobra까지 갑니다.

### 공유 MCP 서버

기본값인 `flag` 방식에서는 ccam이 그룹의 MCP 서버를 인자로 Claude Code에 넘기며, 이 인자는 직접 넘긴 인자보다 앞에 놓입니다.

<!-- motion:launch-mcp -->

`mcp.FlagArgs`는 그룹의 `mcp.json`이 있을 때만 인자를 반환하고, `Plan.InsertArgs`가 그 인자를 `claude` 바로 뒤에 넣습니다. 그래서 직접 지정한 `--mcp-config`나 `--strict-mcp-config`도 그 위에 그대로 적용됩니다. Claude Code의 `--mcp-config`는 값을 여러 개 받으므로 `--mcp-config=<path>` 한 인자로 넘깁니다. 두 인자로 쓰면 뒤따르는 프롬프트까지 설정 파일로 읽게 됩니다.

`config.toml`에 `mcp_strategy = "sync"`를 설정하면 ccam은 대신 계정의 `.claude.json`에 서버를 기록하며, 그 계정이 이미 실행 중이면 이 단계를 건너뜁니다.

### Remote Control

`ccam work rc`(또는 `ccam work remote-control`)는 `work` 계정으로 Claude Code의 Remote Control을 시작합니다. `claude remote-control`은 동사 앞의 `--mcp-config`를 거부하고 동사 뒤에서도 받지 않으며, 그 세션은 계정의 `.claude.json`에서만 MCP 서버를 읽습니다. 그래서 여기서는 flag 방식이 다르게 동작합니다.

<!-- motion:launch-rc -->

1. `remoteControl`은 계정 이름 뒤에서 `-`로 시작하지 않는 첫 번째 인자를 Claude Code의 동사로 봅니다. 그 값이 `remote-control`이나 `rc`이면 `mcp.FlagArgs`는 아무것도 추가하지 않습니다.
2. `mcp.Synced`는 그룹 `mcp.json`의 수정 시각을 계정의 `sync-stamp`와 비교합니다. `.claude.json`에 이미 그룹의 현재 서버가 들어 있거나 그룹에 서버가 없으면, 묻지 않고 Remote Control을 시작합니다.
3. 그렇지 않으면 ccam은 세션에서 빠지게 될 서버 이름을 알려 주고 묻습니다.

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y`를 고르면 `mcp.Sync`가 `sync` 방식과 같은 방법으로 서버를 복사합니다. Claude Code의 잠금을 잡고, 백업을 만든 뒤, 계정 자체 서버는 그대로 둡니다. `N`을 고르거나 터미널이 아니면 이 서버 없이 Remote Control을 시작하고, 실행할 때마다 서버를 복사하는 `mcp_strategy = "sync"`를 안내합니다.

ccam은 서버가 복사될 때까지 Remote Control을 실행할 때마다 묻고, 복사한 뒤에는 그룹의 `mcp.json`이 다시 바뀌었을 때만 한 번 묻습니다. 이후 일반 실행에서는 `.claude.json`의 복사본이 `--mcp-config`보다 우선하므로, `ccam mcp add`, `remove`, `import`가 다른 계정이 가진 복사본도 갱신합니다([MCP 서버 공유](/features/mcp)).

### 10 ms 예산

ccam이 Claude Code 시작 전에 더하는 시간은 최대 10 ms(p95)입니다. 실행 경로에서는 cobra 명령 트리를 만들지 않고, Bubble Tea를 시작하지 않으며, 네트워크 연결도 열지 않습니다. 패키지의 `init()` 함수가 모두 합쳐 5 ms보다 오래 걸리면 `TestPackageInitIsCheap`가 실패합니다. `scripts/perf-check.sh`는 `hyperfine`으로 `ccam --dry-run work`, `ccam work`, 단축 명령을 측정해 10 ms를 넘으면 실패하고, `internal/launcher`의 `BenchmarkLaunchPath`는 이 단계들을 반복 실행합니다. 단계별 시간을 보려면 `CCAM_TRACE=1`을 설정하세요.

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## 안전장치와 예외 상황

- 모르는 이름이나 번호를 주면 `no account "<name>"` 오류를 내고 멈추며, `ccam list`를 안내합니다. 계정이 하나도 없으면 `ccam add <name>`으로 추가하라고 안내합니다.
- 공유 경로가 일반 파일로 바뀌었거나 잘못된 곳을 가리키면 ccam은 `stopped because of share link problems:`와 함께 문제를 나열하고 멈추며, `ccam doctor --fix`를 실행하라고 안내합니다. 실행 경로에서는 데이터를 옮기지 않습니다.
- 이 계정과 대화 내역을 공유하는(같은 그룹의) 다른 Claude Code가 같은 폴더에서 실행 중이면, ccam은 대화 내역이 섞일 수 있다고 경고하고 `Run anyway? [y/N]`를 묻습니다. 터미널이 아니면 경고만 보여 주고 계속합니다.
- `--dry-run`은 MCP 동기화, 새 서버 알림, 대화 내역 확인, 세션 표시 파일, 터미널 제목, 마지막 사용 시각을 건너뛰고 exec도 하지 않습니다. 다만 공유 링크는 여전히 확인하며, 없는 링크는 일반 실행처럼 만듭니다.
- `work`가 이미 실행 중일 때 `ccam work rc`를 실행하면 `.claude.json`을 건드리지 않고 그 사실을 알린 뒤, 그룹의 서버 없이 Remote Control을 시작합니다. 서버를 복사하려면 세션이 끝난 뒤 다시 실행하세요. `--dry-run`에서는 메시지만 출력하고 묻거나 복사하지 않습니다.
- `mcp_strategy = "sync"`에서 동기화가 실패하면 `ccam: MCP sync failed (launching anyway): …`를 출력하고, Claude Code는 기존 서버로 시작합니다.

## 관련 문서

- [단축 명령](/features/shortcuts)
- [ccam으로 Remote Control 사용하기](/guides/remote-control)
- [계정 관리](/features/accounts)
- [명령어 레퍼런스](/docs/commands)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [CLAUDE_CONFIG_DIR의 동작 방식](/guides/claude-config-dir)
