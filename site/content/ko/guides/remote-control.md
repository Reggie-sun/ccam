---
title: ccam으로 Remote Control 사용하기
description: ccam <name> rc로 원하는 계정의 Claude Code Remote Control을 시작하는 방법, 적용되는 옵션, 공유 MCP 서버 복사를 제안하는 이유, 문제 해결 방법을 설명합니다.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: "4d46f49c5dea"
---

`ccam work rc`는 `work` 계정으로 Claude Code의 Remote Control을 시작하므로, 이 컴퓨터의 세션을 claude.ai/code나 Claude 모바일 앱에서 조작할 수 있습니다. 동작은 `claude rc`와 같고 차이는 하나뿐입니다. 계정이 MCP 서버를 공유하고 있으면 ccam은 그 서버를 계정의 `.claude.json`에 복사할지 한 번 묻습니다. Remote Control 세션이 서버를 읽는 곳은 그 파일뿐이기 때문입니다.

## 시작하기

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[단축 명령](/features/shortcuts)도 똑같이 동작합니다. `ccwork rc`는 `ccam work rc`와 같습니다. 작업할 폴더에서 명령을 실행한 다음, `work`와 같은 Claude 계정으로 로그인한 claude.ai/code나 Claude 앱을 여세요.

## 옵션

여기서는 두 종류의 옵션이 만납니다. 계정 이름 앞에 쓰거나 ccam 설정에 두는 ccam 자체 옵션, 그리고 `rc` 뒤에 써서 Claude Code에 그대로 전달되는 Remote Control 옵션입니다.

### ccam 옵션과 설정

| 옵션 또는 설정 | 동작 | 기본값 |
| --- | --- | --- |
| `--dry-run` (계정 앞) | 환경 변수와 명령, 그리고 MCP 안내가 있으면 그것도 출력합니다. 묻지도 복사하지도 않습니다 | 꺼짐 |
| `~/.config/ccam/config.toml` 맨 위의 `mcp_strategy = "sync"` | Remote Control 여부와 관계없이 실행할 때마다 그룹의 서버를 `.claude.json`에 복사하므로, 아무것도 묻지 않습니다 | `"flag"` |
| `ccam group leave <group> work` | MCP를 공유하는 그룹에서 이 계정을 뺍니다. 그룹의 서버가 계정 자체의 서버로 `.claude.json`에 복사되므로, Remote Control이 그 서버를 갖고 아무것도 묻지 않습니다. `--fresh`를 주면 계정 자체의 서버만으로 시작합니다 | `default`에 속함 |

계정이 MCP 서버를 공유하고, 그 `.claude.json`에 그룹의 현재 목록이 없거나 복사한 서버가 사라졌으면 ccam이 `Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]`라고 묻습니다. `y`는 복사하고, `N`이나 Enter는 복사 없이 시작합니다.

### Remote Control 옵션

`rc` 뒤의 모든 것은 Claude Code가 받습니다. 아래는 Claude Code 2.1.304 기준 `claude remote-control`의 옵션이며, 사용 중인 버전의 목록은 `ccam work rc --help`로 볼 수 있습니다.

| 옵션 | 동작 |
| --- | --- |
| `--name <name>` | claude.ai/code에 표시되는 세션 이름 |
| `--remote-control-session-name-prefix <prefix>` | 자동으로 만드는 세션 이름의 접두사(기본값: 호스트 이름) |
| `-c`, `--continue` | Remote Control이 이 폴더에서 마지막으로 기록한 세션에 다시 연결(약 4시간 이내) |
| `--session-id <id>` | 특정 세션에 다시 연결 |
| `--permission-mode <mode>` | 새 세션의 권한 모드: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk`, `plan` 중 하나 |
| `--spawn <mode>` | `same-dir`(기본값), `worktree`(세션마다 git worktree 하나), `session`(세션 하나만 열고, 끝나면 종료) |
| `--capacity <N>` | `same-dir`와 `worktree` 모드에서 동시에 열 수 있는 최대 세션 수(기본값 32) |
| `--[no-]create-session-in-dir` | 시작할 때 현재 폴더에 첫 세션을 만듦(기본값 켜짐) |
| `--[no-]chrome` | 새 세션에서 Claude in Chrome 사용(기본값 꺼짐) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | 디버그 출력과 자세한 출력 |

ccam은 Claude Code와 같은 방식으로, 옵션과 그 값을 건너뛰며 Claude Code의 명령을 읽습니다. 그래서 `ccam work rc --name laptop`과 `ccam work --model opus rc` 모두 Remote Control을 시작합니다. ccam이 모르는 옵션은 값이 없는 옵션으로 치므로, 안전하게 하려면 `rc`를 계정 이름 바로 뒤에 쓰세요.

## ccam이 하는 일

실행은 평소와 같은 경로([Claude Code 시작하기](/features/launch))를 따르고, MCP 단계 하나만 다릅니다.

1. ccam은 여느 실행과 마찬가지로 계정을 찾고, 공유 링크를 확인하고, `claude`의 위치를 찾습니다.
2. 명령이 `rc`나 `remote-control`임을 확인하고 `--mcp-config`를 넣지 않습니다.
3. 계정이 MCP 서버를 공유하면 그룹의 `mcp.json`을 `~/.local/state/ccam/mcp/work/`의 `sync-stamp`와 비교하고, 그룹의 모든 서버가 아직 `.claude.json`에 있는지 확인합니다. `.claude.json`에 이미 현재 서버가 들어 있거나 그룹에 서버가 없으면 바로 6단계로 넘어갑니다.
4. 그렇지 않으면 세션에서 빠지게 될 서버 이름을 알려 줍니다.

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

5. `y`를 고르면 `.claude.json`을 백업하고, Claude Code 자체 잠금을 잡은 뒤, 그룹의 서버를 `mcpServers`에 추가하고 계정의 다른 서버는 그대로 둡니다. 계정에 그룹 서버와 이름이 같은 자체 서버가 있으면 ccam은 계정의 서버를 유지하고 그 사실을 알립니다. `N`을 고르거나, 터미널이 아니거나, 계정이 실행 중이면 아무것도 바꾸지 않고 그 사실을 알립니다.

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. 세션을 기록하고, `CLAUDE_CONFIG_DIR`을 계정 폴더로 설정한 뒤, `exec`로 자신을 `claude rc`로 바꿉니다.

그 뒤 Remote Control은 서버로 계속 실행됩니다. 여러분을 위해 시작하는 모든 세션은 이 프로세스의 자식이므로 `work`로 실행됩니다. 즉 `work`의 로그인, 설정, 대화 내역을 쓰고, 사용량은 `work`의 사용량 한도에 집계됩니다. ccam은 자신이 시작한 Remote Control 서버가 실행되는 동안 이 계정을 실행 중으로 봅니다. 그냥 `claude remote-control`로 시작한 서버는 세션이 있는 동안에만 실행 중으로 치며, 아직 세션이 없는 서버는 ccam이 알아보지 못합니다.

차이는 dry run으로 확인할 수 있습니다. 일반 실행은 그룹의 서버를 인자로 넘기지만, Remote Control은 넘기지 않습니다.

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## 이렇게 동작하는 이유

**Remote Control은 `--mcp-config`를 받지 않습니다.** 일반 실행에서 ccam은 그룹의 서버를 `--mcp-config=<path>`로 Claude Code에 넘기고, Claude Code의 파일에는 절대 쓰지 않습니다. `claude remote-control`은 명령 이름 앞의 이 옵션을 거부하고 뒤에서도 받지 않기 때문에, ccam 0.6.0 이전에는 MCP 서버를 공유하는 모든 계정에서 `ccam work rc`가 실패했습니다. 이제 ccam은 이 옵션을 넣지 않습니다.

**그 세션은 `.claude.json`에서만 서버를 읽습니다.** Remote Control이 시작하는 세션은 명령줄 옵션을 물려받지 않습니다. 이 세션이 사용자 MCP 서버를 찾는 곳은 계정의 `.claude.json` 하나뿐입니다. 이 파일에는 계정의 로그인 정보와 폴더별 상태도 들어 있으므로, ccam은 `settings.json`처럼 공유 사본에 링크할 수 없습니다. 세션에 공유 목록을 주는 방법은 서버를 이 파일에 복사하는 것뿐입니다.

**ccam은 쓰기 전에 묻습니다.** 기본 flag 방식에서 ccam은 실행할 때 `.claude.json`을 건드리지 않습니다. Remote Control은 그 파일에 써야 하는 유일한 실행이므로, 먼저 묻고, 백업을 남기고, `mcpServers` 외에는 아무것도 건드리지 않습니다. (실행 밖에서는 MCP 공유에서 빠지거나 다시 들어갈 때와 `ccam mcp import`도 계획에 밝힌 대로 이 파일에 씁니다.)

**계정이 실행 중일 때는 쓰지 않습니다.** Claude Code는 세션 중에 `.claude.json`을 다시 쓰므로, 외부에서 바꾼 내용은 사라지거나 덮어써집니다. 그래서 ccam은 계정의 세션이 하나라도 실행 중이면 복사를 건너뜁니다. 다른 Remote Control 서버와 ccam 없이 시작한 Claude Code도 여기에 포함됩니다.

**변경될 때마다 한 번만 묻습니다.** 복사한 뒤에는 `sync-stamp`에 그룹 `mcp.json`의 수정 시각이 들어갑니다. ccam은 `ccam mcp add` 뒤처럼 그 파일이 바뀌었거나 복사한 서버가 `.claude.json`에서 사라졌을 때만 다시 묻습니다. `N`으로 답하면 아직 복사본이 없으므로 Remote Control을 실행할 때마다 묻습니다.

**`ccam mcp`가 복사본을 최신으로 유지합니다.** 이름이 겹치면 `.claude.json`의 항목이 `--mcp-config`보다 우선하므로, 오래된 복사본은 일반 실행에서도 새 정의를 가리게 됩니다. 그래서 `ccam mcp add`, `add-json`, `remove`, `import`와 `ccam migrate`는 실행 중이 아닌 모든 계정의 복사본을 갱신하고, 복사본이 바뀌었으면 그 사실을 알립니다.

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**터미널이 없으면 답은 '아니요'입니다.** Remote Control은 스크립트, 서비스, `nohup`처럼 아무도 답할 수 없는 곳에서 실행되는 경우가 많습니다. 이때 ccam은 기다리지 않고 서버 없이 시작합니다.

## 문제와 해결 방법

### Remote Control 세션에 공유 MCP 서버가 없습니다

질문에 `N`으로 답했거나, 터미널이 없었거나, 계정이 실행 중이었던 경우입니다. Remote Control을 멈추고 계정의 다른 세션을 닫은 뒤, 터미널에서 `ccam work rc`를 다시 실행해 `y`로 답하세요. 더 이상 묻지 않게 하려면 `~/.config/ccam/config.toml`에 `mcp_strategy = "sync"`를 설정하세요. 이미 열린 세션의 서버는 실행 중에 바뀌지 않으므로, 복사한 뒤에는 새 세션을 시작하세요.

### "Account work is running, so its .claude.json is left alone"

`work`의 세션이 아직 열려 있습니다. 일반 `ccam work`, 단축 명령, 또는 다른 Remote Control 서버일 수 있습니다. 그 세션을 종료하고 `ccam work rc`를 다시 실행하세요. 다른 계정이 실행하는 세션은 상관없습니다.

### ccam이 매번 묻습니다

아직 서버를 복사하지 않았기 때문입니다. `y`로 한 번 답하거나 `mcp_strategy = "sync"`를 설정하세요. Remote Control에서 이 서버를 쓸 생각이 없다면 서버를 남기지 않고 이 계정과 MCP 공유를 중단하세요. `ccam group leave default work --fresh`를 실행하거나, MCP를 이 계정이 속하지 않은 별도 그룹에 두면 됩니다. 그냥 `ccam group leave`를 실행하면 그룹의 서버가 계정의 `.claude.json`에 복사되므로 Remote Control이 그 서버를 갖게 됩니다.

### `ccam work rc`가 `--mcp-config` 관련 오류와 함께 바로 실패합니다

사용 중인 ccam이 0.6.0보다 오래되어 Remote Control에 `--mcp-config`를 넘기고 있습니다. `ccam update`를 실행하세요. 최신 ccam에서 같은 오류가 나면, 값을 받는 알 수 없는 옵션이 `rc` 앞에 와서 ccam이 그 값을 명령으로 잘못 본 것입니다. `rc`를 계정 이름 바로 뒤에 쓰세요.

### `ccam mcp`로 바꾼 내용이 Remote Control에 반영되지 않습니다

`ccam mcp add`와 `remove`는 실행 중이 아닌 계정의 복사본만 갱신하며, 그럴 때는 다음과 같이 알립니다: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control을 멈췄다가 다시 시작하세요. ccam이 바뀐 `mcp.json`을 알아채고 다시 묻습니다.

### 그룹 서버 대신 내 서버가 쓰입니다

계정 자체에 그룹 서버와 이름이 같은 서버가 있으면, 복사는 그 서버를 그대로 두고 그 사실을 알립니다: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`. 그룹의 서버를 쓰려면 계정의 서버를 지우세요(`claude mcp remove -s user <name>`). 내 정의를 모든 계정에 주려면 `ccam mcp import work <server> --force`로 그룹에 옮기세요. `.claude.json` 백업은 `~/.local/state/ccam/mcp/work/backups/`에 있으며, 최신 5개가 보관됩니다.

### 복사한 서버를 지우고 싶습니다

`ccam mcp remove <server>`는 그룹과 모든 복사본에서 서버를 지웁니다. 복사하기 전 상태로 돌아가려면, 계정의 세션이 하나도 실행되지 않을 때 `~/.local/state/ccam/mcp/work/backups/`에서 가장 최근 백업을 복원하세요.

### Remote Control이 폴더를 신뢰할지 묻거나, 터미널 없이 종료됩니다

폴더 신뢰 여부는 계정마다 `.claude.json`에 저장되므로, `personal`로 신뢰한 폴더도 `work`에는 처음 보는 폴더입니다. 그 폴더에서 `ccam work`를 한 번 실행해 신뢰를 승인한 다음, 거기서 Remote Control을 시작하세요.

### "You must be logged in with a Claude account that has a subscription"

Remote Control에는 Claude 구독 계정 로그인이 필요합니다. `ccam list`에서 계정의 이메일을 확인하고, 필요하면 `ccam login work`로 다시 로그인하세요.

### 휴대폰에 세션이 보이지 않습니다

Claude 앱이나 claude.ai/code가 다른 Claude 계정으로 로그인되어 있습니다. `ccam list`가 `work`에 대해 보여 주는 이메일과 비교해 보세요.

### 다른 Claude Code 세션이 이 폴더에서 대화 내역을 공유한다고 ccam이 경고합니다

같은 대화 내역 그룹의 다른 계정이 같은 폴더에서 실행 중입니다. 두 계정에서 같은 대화를 이어 가는 세션은 대화 내역을 섞을 수 있습니다. 서로 다른 대화로 작업할 때만 `y`로 답하거나, 다른 폴더에서 Remote Control을 시작하세요.

## 관련 문서

- [Claude Code 시작하기](/features/launch)
- [MCP 서버 공유](/features/mcp)
- [그룹으로 공유하기](/features/share)
- [계정 관리](/features/accounts)
