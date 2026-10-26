---
title: MCP 서버 공유
description: ccam mcp add, add-json, list, remove, import로 그룹마다 MCP 서버 목록을 하나씩 관리하고, --mcp-config로 넘기거나 각 계정에 동기화합니다.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: "27a06ad35295"
---

Claude Code는 사용자 MCP 서버를 계정의 로그인 정보가 들어 있는 `.claude.json`에 함께 저장하므로, 이 파일은 계정 간에 공유할 수 없습니다. ccam은 공유 서버를 그룹 파일인 `mcp.json`에 두고, 실행할 때 그룹의 모든 계정에 이 파일을 넘깁니다. `ccam mcp`로 서버를 한 번만 추가해 모든 계정에서 쓰고, 한 세션 안에서 추가한 서버는 `ccam mcp import`로 다른 계정에 옮기세요.

## 사용법

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
```

변경 사항은 각 계정을 다음에 실행할 때부터 적용됩니다. 실행 중인 세션은 시작할 때의 서버를 그대로 씁니다.

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

그룹 아래에 `ccam mcp list`는 다른 계정에 닿지 않는 모든 서버를 알려 줍니다.

- 계정이 자기만을 위해 추가한 서버는 `only in account work: pg`
- 한 폴더에만 추가한 서버(그냥 `claude mcp add`, 아래 참고)는 `only in account side, for /home/me/code/app only: notes`
- 계정 자체 서버가 그룹 서버와 이름이 같으면 `account work has its own …, used in its sessions in place of the group's`
- 이 폴더의 `.mcp.json`은 승인을 기다리는 계정과 함께 `waiting for approval in`
- `in no account, only in ~/.claude.json (plain claude's)`인 서버는 이를 복사하는 [`ccam migrate --mcp`](/features/migrate) 명령과 함께 표시

그룹에 이미 있는 이름을 추가하면 기존 정의를 바꾸고, `added` 대신 `updated`를 출력합니다. 서버는 그 그룹에 속해 있고 그룹이 MCP를 공유할 때(`ccam group set <group> --add mcp`)만 계정에 전달됩니다. 그룹이 MCP를 공유하지 않으면 `add`가 그 사실을 알려 줍니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json`, `remove`의 대상 그룹 | 어떤 계정이 MCP를 공유하는 그룹이 하나뿐이면 그 그룹, 아니면 `default` |
| `-- <command> [args…]` (`add`) | 이 명령을 stdio 서버로 실행 | |
| `--url <url>` (`add`) | 원격 서버 URL(http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` 서버의 전송 방식 | `http` |
| `--header <'Name: value'>` (`add`) | HTTP 헤더. 여러 번 지정 가능 | 없음 |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio 서버의 환경 변수. 여러 번 지정 가능 | 없음 |
| `--force` (`import`) | 그룹에 같은 이름으로 정의가 다른 서버가 있으면 바꿈. 이전 `mcp.json`은 백업됨 | 꺼짐 |
| `--keep-local` (`import`) | 계정의 `.claude.json`에 있는 항목을 남겨 둠 | 꺼짐 |
| `-h`, `--help` | 각 명령의 도움말 | |
| `--dry-run` (전역) | `mcp list`만 지원. `add`, `add-json`, `remove`, `import`는 이 옵션을 거부하고 아무것도 바꾸지 않음 | 꺼짐 |

`add`에는 `--` 뒤의 명령이나 `--url` 중 하나가 필요하며, 둘 다 없으면 `a command to run (after --) or --url is required`로 실패합니다. `--header`는 `--url` 서버에, `--env`는 명령 서버에 적용됩니다. `add-json`은 `claude mcp add-json`과 같은 JSON을 받으며, 정의는 `command`나 `url`이 있는 객체여야 합니다. `import`는 항상 계정이 MCP를 공유하는 그룹을 쓰고 `list`는 모든 그룹을 보여 주므로, `-g`는 이 두 명령에 영향을 주지 않습니다.

## 동작 방식

관련 파일은 다음과 같습니다.

- `~/.local/share/ccam/groups/<group>/mcp.json`: 공유 서버. `--mcp-config`가 읽는 형식인 `{"mcpServers": {…}}`로 저장됩니다. 서버 항목에 API 키가 들어 있는 경우가 많으므로 `fsutil.AtomicWrite`로 `0600` 권한을 주어 기록합니다. 또한 `<group>/.mcp.json.lock` 잠금을 잡은 상태에서만 기록하므로, `ccam mcp add` 명령 두 개를 동시에 실행해도 두 서버가 모두 남습니다.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: 계정의 서버가 들어 있는 Claude Code 자체 파일. ccam은 `profile.TopLevelValue` 스캐너로 `mcpServers` 값만 읽으며, 실행 경로에서 파일 전체를 해석하지 않습니다.
- `~/.local/state/ccam/mcp/<name>/`: `seen` 파일, `sync-stamp`(ccam이 쓴 복사본마다 해시가 들어 있어, 아무도 바꾸지 않은 복사본과 직접 고친 복사본을 구분할 수 있습니다), `sync-checked`(sync 방식이 마지막으로 본 상태), `kept`(MCP 공유에서 빠질 때 계정이 남겨 둔 서버), `.claude.json` 백업.

방식은 두 가지이며 `~/.config/ccam/config.toml`의 `mcp_strategy`로 정합니다. 기본값은 옵션 방식입니다.

<!-- motion:mcp-flag -->

1. ccam이 계정의 MCP 그룹을 찾고(`launcher.applyMCP`), `stat` 한 번으로 그룹의 `mcp.json`이 있는지만 확인할 뿐 파일을 해석하지는 않습니다(`mcp.FlagArgs`).
2. ccam이 `claude` 바로 뒤, 직접 넘긴 인자보다 앞에 `--mcp-config=<path>`를 넣습니다. 그래서 직접 지정한 `--mcp-config`나 `--strict-mcp-config`도 그 위에 그대로 적용됩니다(`Plan.InsertArgs`).
3. ccam은 계정의 `.claude.json`이 지난 실행 이후 바뀌었는지 확인하고, 바뀐 경우에만 서버 목록을 해석합니다(`mcp.NewLocal`).
4. ccam이 자신을 Claude Code로 바꿉니다(`syscall.Exec`).

경로는 두 인자가 아니라 `--mcp-config=<path>` 한 인자로 넘깁니다. Claude Code의 `--mcp-config`는 값을 여러 개 받으므로, `claude --mcp-config a.json "prompt"`라고 하면 프롬프트를 두 번째 설정 파일로 읽어 `MCP config file not found`로 실패합니다. 한 인자 형식은 이 문제를 피하므로, 뒤따르는 프롬프트나 하위 명령이 정상 동작합니다. `ccam work mcp list`, `get`, `remove`에는 Claude Code가 거기서 이 옵션을 거부하므로 넣지 않으며, 대신 ccam이 그룹의 서버를 표준 오류로 알려 줍니다([자세히](/features/launch#shared-mcp-servers)). Claude Code는 `--mcp-config`의 서버를 계정 자체 서버에 더해 불러오며, 이름이 겹치면 `.claude.json`의 항목이 우선합니다. 옵션 방식이 `.claude.json`에 쓰는 것은 Remote Control, 계정이 MCP 공유에서 빠지거나 다시 들어올 때([아래](#leaving-and-joining)), 그리고 `ccam mcp import`뿐입니다.

`claude remote-control`(`rc`)은 `--mcp-config`를 받지 않으며, 그 세션은 `.claude.json`에서만 서버를 읽습니다. `ccam work rc`에서는 ccam이 이 옵션을 빼고, 터미널이면 `sync` 방식처럼 백업을 남기고 그룹의 서버를 `.claude.json`에 복사할지 묻습니다. 복사본이 생길 때까지 묻고, 그 뒤에는 그룹의 `mcp.json`이 바뀌었거나 복사한 서버가 사라진 뒤에만 다시 묻습니다([Remote Control](/features/launch#remote-control)). 이 복사본은 이름이 겹칠 때 `--mcp-config`보다 우선하므로, `ccam mcp add`, `add-json`, `remove`, `import`는 그룹의 다른 계정이 가진 복사본도 갱신합니다.

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

`✓ Updated` 줄은 복사본이 바뀌었을 때만 나옵니다. 이전 버전의 ccam이 만든 복사본도, `ccam migrate` 뒤의 복사본도 같은 방식으로 최신 상태가 됩니다.

세션 안에서 `claude mcp add`로 추가한 서버는 해당 계정의 `.claude.json`에만 저장됩니다. 그냥 `claude mcp add`는 local 범위를 쓰므로, 그 서버는 그 한 계정의 그 한 폴더에서만 불러옵니다(`-s user`를 주면 그 계정이 모든 폴더에서 쓸 수 있습니다). 그 시점에는 Claude Code가 이미 ccam 프로세스를 대체한 상태라 ccam이 이를 볼 수 없으므로, 다음 실행 때 확인합니다. `seen` 파일에는 `.claude.json`의 수정 시각과 이미 알린 서버 이름이 들어 있습니다. 시각이 다르면 ccam은 서버 이름을 그룹과 비교하고, 새 서버마다 입력을 기다리지 않고 한 번씩 안내를 출력합니다. 두 방식 모두 그렇습니다.

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

이 안내는 그룹 서버와 이름이 같은 계정 자체 서버도 알려 줍니다. 그 계정의 세션에서는 계정 쪽 서버가 우선하기 때문입니다.

`mcp_strategy = "sync"`로 설정하면 서버가 `claude mcp list`와 `/mcp`에 일반 사용자 서버로 나타나며, ccam은 실행 전에 각 계정의 `.claude.json`에 서버를 기록합니다.

<!-- motion:mcp-sync -->

1. ccam이 그룹 `mcp.json`의 수정 시각을 `sync-stamp`와, 계정의 `.claude.json`을 마지막으로 본 상태와 비교합니다(`mcp.Sync`). 둘 다 바뀌지 않았으면 더 하는 일이 없습니다. 바뀌었으면 그룹의 모든 서버가 아직 `.claude.json`에 있는지도 확인하고, 사라진 복사본은 다시 넣습니다(`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Claude Code 자체 잠금, 즉 파일 옆에 `mkdir`로 만드는 `.claude.json.lock` 폴더를 잡습니다. 최대 2초까지 기다리며, Claude Code의 잠금 라이브러리와 마찬가지로 10초보다 오래된 잠금은 유효하지 않은 것으로 봅니다.
3. 현재 `.claude.json`을 상태 폴더의 `backups/`에 저장하고, 최신 5개만 남깁니다.
4. 그룹의 모든 서버를 `mcpServers`에 설정하고, 전에 기록했지만 이제 그룹에 없는 서버는 제거하며, 무언가 바뀌었을 때만 파일을 원래 권한으로 원자적으로 기록합니다. 자신이 쓴 복사본 중 그 뒤로 아무도 바꾸지 않은 것만 바꾸거나 지웁니다. 계정에 그룹 서버와 이름이 같은 자체 서버가 있으면 ccam은 계정의 서버를 유지하고 한 번 알려 줍니다: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. 복사본마다 해시를 담은 새 스탬프를 기록하고, `--mcp-config` 없이 Claude Code를 시작합니다.

Claude Code는 실행 중에 `.claude.json`을 다시 쓰므로, Claude Code를 어떻게 시작했든 그 계정의 세션이 실행 중이면 `ccam: account work is already running, skipping MCP sync`를 출력하고 동기화를 건너뜁니다. 동기화가 실패하면 `ccam: MCP sync failed (launching anyway): …`를 출력하고 그대로 실행합니다.

### 빠질 때와 다시 들어올 때

MCP 공유를 멈추는 계정(`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`)은 그룹의 서버를 남겨 둡니다. 계획이 백업을 남기고 그 서버를 계정의 `.claude.json`에 복사하며, 같은 이름의 계정 자체 서버는 그대로 남습니다. `--fresh`를 주면 아무것도 받지 않고, 계획에 더 이상 불러오지 않을 서버가 표시됩니다. 그 서버가 있는 그룹에 다시 들어가면 아무도 바꾸지 않은 남겨 둔 복사본을 지우므로, 이후 그룹에서 고친 내용이 계정에 닿습니다.

`ccam mcp import`는 계정 자체의 서버를 그 계정의 그룹으로 옮깁니다.

<!-- motion:mcp-import -->

ccam은 계정 자체 서버와 한 폴더에서만 가진 서버를 읽고, 각 서버를 그룹에 있는 같은 이름의 서버와 `mcp.SameDef`로 비교합니다. 이 함수는 키 순서와 공백을 무시합니다. 폴더 서버는 모든 폴더에서 쓰는 그룹 서버가 됩니다(`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`). 한 이름에 서로 다른 폴더 정의가 둘 있으면 거부합니다. 플러그인 서버는 플러그인이 가져오므로 옮기자고 하지 않습니다. 그룹에 없거나 정의가 같은 서버는 복사합니다. 정의가 다른 서버는 건너뜁니다. 다른 모든 계정이 쓰는 것은 그룹의 정의이기 때문입니다.

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force`를 붙이면 그룹의 서버를 바꾸며, 이전 `mcp.json`은 먼저 그룹 폴더 안의 `.ccam-backup/<time>/mcp.json`에 저장합니다. 그다음 ccam이 같은 잠금을 잡고 상태 폴더에 백업을 남긴 뒤, 복사한 항목을 계정 자체 항목과 폴더 항목 모두 한 번에 계정의 `.claude.json`에서 제거합니다(`mcp.Remove`). ccam은 그룹의 `mcp.json`을 읽은 순간부터 그 항목이 제거될 때까지 잠가 두므로, 그 사이에 누가 그룹에서 서버를 지우는 일은 없습니다. 이름이 겹치면 계정의 항목이 그룹의 항목보다 우선하므로 이 단계가 필요합니다. `--keep-local`을 주면 항목이 계정의 `.claude.json`에 남습니다.

## 안전장치와 예외 상황

- **로그인.** OAuth로 로그인하는 서버는 토큰을 서버 목록이 아니라 각 계정의 `.credentials.json`(`mcpOAuth`)에 보관하므로, 서버를 공유해도 로그인은 공유되지 않습니다. 계정마다 `/mcp`로 한 번씩 로그인하세요. ccam은 로그인을 절대 복사하지 않으며, [`ccam migrate`](/features/migrate)도 마찬가지입니다. 복사한 리프레시 토큰은 같은 토큰을 가진 다른 모든 쪽을 로그아웃시키기 때문입니다. [`ccam ext`](/features/ext)와 `ccam doctor`는 다른 계정은 로그인했지만 이 계정은 로그인하지 않은 서버를 알려 주고, `ccam next --continue`는 새 계정에 빠진 로그인을 알려 줍니다. 그룹에는 로그인 파일이 절대 들어가지 않습니다.
- **실행 중인 계정.** `import`는 그룹에 복사하지만 실행 중인 계정 자체의 항목은 남겨 두고 그 사실을 알려 줍니다. 세션이 끝난 뒤 다시 실행하세요. 동기화는 실행 중인 계정을 건너뜁니다. 둘 다 ccam 없이 시작한 Claude Code 세션도 감지합니다.
- **Remote Control용 복사본.** 실행 중인 계정의 복사본은 건드리지 않고 다음과 같이 알립니다. `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"`에서는 실행할 때마다 동기화하므로 갱신할 것이 없습니다.
- **손상된 파일.** `.claude.json`이 올바른 JSON이 아니면 동기화는 파일을 건드리지 않고 오류를 보고합니다. 2초 안에 잠금을 잡지 못하면 `could not lock … (Claude Code is writing the file)`를 보고합니다.
- **MCP를 공유하지 않는 계정.** `import`는 거부하고, 먼저 계정을 MCP를 공유하는 그룹에 넣으라고 안내합니다(`ccam group join <group> side`).
- **그룹 이름.** `-g`에는 그룹 이름만 줄 수 있습니다. 이름이 `groups/` 아래 폴더가 되므로 `../x` 같은 이름은 거부합니다.

## 관련 문서

- [그룹으로 공유하기](/features/share)
- [ccam으로 Remote Control 사용하기](/guides/remote-control)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [기존 설정 복사하기](/features/migrate)
- [명령어 레퍼런스](/docs/commands)
