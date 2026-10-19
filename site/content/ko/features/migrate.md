---
title: 기존 설정 복사하기
description: ccam migrate로 기존 ~/.claude의 설정, CLAUDE.md, 확장 기능, MCP 서버, MCP 로그인 정보를 ccam 공유 그룹으로 복사합니다.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

ccam을 쓰기 전에 Claude Code를 써 왔다면 설정이 `~/.claude`와 `~/.claude.json`에 있습니다. `ccam migrate`는 이 설정을 공유 그룹으로 복사하므로, 그룹을 통해 해당 항목을 공유하는 모든 계정이 기존 설정 그대로 시작합니다. 원본 폴더는 읽기만 하므로, `claude`를 그냥 실행해도 이전과 똑같이 동작합니다. 처음 `ccam add`를 실행한 뒤 한 번 실행하고, 어떤 계정이 MCP 서버 공유를 시작할 때마다 다시 실행하세요.

## 사용법

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

터미널에서 항목 옵션과 `--yes` 없이 실행하면 ccam이 먼저 체크리스트를 보여 줍니다.

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space(또는 `x`)로 커서 위치의 항목 선택을 바꾸고, 화살표 키나 `j`/`k`로 이동하며, Enter를 누르면 계획을 볼 수 있습니다. Esc, `q`, Ctrl+C를 누르면 `Cancelled.`를 출력합니다. 원본에 해당 내용이 없는 항목에는 `not found`가 표시됩니다. 모든 선택을 해제하면 ccam은 `Nothing selected.`를 출력하고 멈춥니다.

그다음 ccam은 계획을 출력하고, 복사하기 전에 확인을 받습니다.

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

그룹의 내용은 그 그룹을 통해 항목을 공유하는 계정에만 전달되므로, 요약에는 항목마다 그런 계정이 나열됩니다. 공유하지 않는 계정이 있으면 ccam은 공유를 켜는 `ccam share <name> --<item> on` 명령을 출력합니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `[folder]` | 복사할 원본 Claude Code 폴더 | `~/.claude` |
| `--settings` | `settings.json`을 복사합니다 | 꺼짐 |
| `--memory` | 전역 `CLAUDE.md`를 복사합니다 | 꺼짐 |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/`를 복사합니다 | 꺼짐 |
| `--mcp` | 사용자 MCP 서버(`.claude.json`의 `mcpServers`)를 복사합니다 | 꺼짐 |
| `--history` | 대화 내역(`projects/`, `history.jsonl`)을 복사합니다 | 꺼짐 |
| `-g`, `--group string` | 복사할 대상 공유 그룹 | `default` |
| `-y`, `--yes` | 묻지 않고 복사합니다. 충돌하면 그룹의 사본을 유지합니다 | 꺼짐 |
| `--dry-run`(전역) | 계획을 출력하고, 질문하거나 복사하기 전에 멈춥니다 | 꺼짐 |

항목 옵션이 없으면 터미널에서는 체크리스트로 정합니다. 터미널이 아니거나 `--yes`를 주면 대화 내역을 제외한 모든 항목을 복사합니다. `projects/`는 클 수 있으므로 대화 내역은 기본적으로 제외합니다.

`--group on`은 `default`를 뜻합니다. `--group off`와 `--group local`은 `--group takes a group name, not "off"`로 거부합니다. `ccam migrate`는 `--dry-run`을 지원하며, 시험 실행에서도 체크리스트가 나타나므로 선택 결과를 미리 볼 수 있습니다.

## 동작 방식

`internal/cli/migrate.go`가 항목을 고르면 `migrate.Build`가 원본과 그룹을 읽어 계획을 반환합니다. 확인을 받은 뒤 `Plan.Apply`가 실행되기 전에는 아무것도 기록하지 않습니다.

1. **원본 찾기.** `checkSource`는 없는 경로, 폴더가 아닌 경로, 알려진 Claude Code 항목이 하나도 없는 경로를 거부합니다. `migrate.ClaudeJSON`은 폴더에 해당하는 `.claude.json`을 찾습니다. 폴더 안에 있으면 그것을 쓰고, 안에 없는 `~/.claude`라면 그 옆의 `~/.claude.json`을 씁니다. `CLAUDE_CONFIG_DIR`가 설정되지 않았을 때 Claude Code가 이 파일을 두는 곳입니다.
2. **항목을 경로에 대응시키기.** 모든 경로는 `ccam share`가 쓰는 것과 같은 허용 목록인 `share.Layout`에서 나옵니다. MCP 서버는 `mcp.LocalServers`로 `.claude.json`에서 읽으며, 이 함수는 파일 전체를 해석하지 않고 `profile.TopLevelValue` 스캐너를 씁니다.
3. **옮기지 않고 복사만 하기.** 파일은 `fsutil.AtomicWrite`나 `fsutil.CopyFile`로 `~/.local/share/ccam/groups/<group>/`에 복사하고, MCP 서버는 그룹의 `mcp.json`에 넣습니다.

각 항목은 원본 경로 하나를 그룹 경로 하나에 대응시킵니다.

<!-- motion:migrate-items -->

대화 내역 두 행은 `--history`를 줄 때만 복사되므로 점선으로 표시됩니다. `.credentials.json`과 `.claude.json`은 `share.Never`에 들어 있어 그룹에 절대 들어가지 않습니다. 레이아웃에 이 파일이 나타나면 `Build`가 내부 오류를 반환하고, 폴더를 탐색할 때도 이름으로 건너뜁니다.

`buildEntry`는 모든 경로에 상태를 부여하며, 계획의 각 줄은 이 상태를 설명합니다.

<!-- motion:migrate-states -->

- **없음**: 원본에 없거나, 폴더여야 할 자리에 파일이 있는 경우입니다. 항목마다 한 번 `not found, skipped`로 표시됩니다.
- **링크됨**: 원본 경로가 이미 그룹의 사본을 가리키는 경우입니다. 예를 들어 원본이 ccam 프로필 폴더일 때입니다.
- **같음**: 그룹에 이미 같은 내용이 있는 경우입니다.
- **새 항목**: 그룹에 없으므로 통째로 복사합니다.
- **병합**: 양쪽에 모두 있는 폴더입니다. `copyMissing`은 그룹에 없는 파일만 추가하며, 내용이 다른 파일은 그룹에 있는 그대로 두고 계획에 `group copies kept`로 표시합니다. `plugins/`에서는 레지스트리 `installed_plugins.json`과 `known_marketplaces.json`을 병합하고, 플러그인의 절대 경로를 그룹 폴더를 가리키도록 다시 씁니다.
- **충돌**: 양쪽에 모두 있지만 내용이 다른 파일입니다. 파일마다 `keep group`과 `use mine` 중에서 고르며, `.json` 파일은 `merge JSON`(값이 겹치면 그룹이 우선하는 깊은 병합), `history.jsonl`은 `merge by line`도 고를 수 있습니다. 정의가 다른 MCP 서버도 충돌로 처리하며, `keep group`이나 `use mine`을 고릅니다. 맨 처음 나오는 선택지가 Enter를 눌렀을 때의 기본값이고, `--yes`는 항상 그룹의 사본을 유지합니다.

체크리스트 자체는 `ccam add`가 쓰는 `checklist` 컴포넌트로 만든 `tui.MigratePicker`입니다. `migrate.Found`는 lstat 호출과 MCP 서버 목록으로 없는 항목을 표시합니다.

### MCP 로그인 정보

OAuth로 로그인하는 MCP 서버는 토큰을 서버 목록이 아니라 로그인 파일인 `.credentials.json`의 `mcpOAuth` 아래에 보관합니다. Claude Code는 각 항목의 키를 `<name>|`에 서버의 `type`, `url`, `headers`로 만든 `sha256` 해시의 앞 16자를 붙여 만듭니다. 키가 폴더와 무관하므로, 서버 정의만 같다면 복사한 항목이 다른 계정에서도 동작합니다. 그래서 `internal/migrate/logins.go`의 `buildLogins`는 최종적으로 그룹이 원본의 정의를 갖게 되는 서버의 항목만 복사합니다. 새 항목이거나, 이미 같거나, 충돌에서 `use mine`을 고른 서버입니다. 플러그인 서버의 로그인 정보(`plugin:`으로 시작하는 키)는 확장 기능을 복사하고 원본에 `plugins/`가 있을 때 함께 복사됩니다.

<!-- motion:migrate-signins -->

이 항목들은 받는 계정 자체의 `.credentials.json`에 `fsutil.AtomicWrite`로 `0600` 권한을 주어 기록합니다. 그룹에는 절대 들어가지 않습니다. 계정이 대상 그룹을 통해 `mcp`(플러그인 서버라면 `extensions`)를 공유할 때만 받으며, 계정에 없는 키만 추가합니다. 이미 있는 항목은 더 최신일 수 있으므로 그대로 둡니다. 계획의 `sign-ins` 줄에는 서버, 받는 계정, 건너뛴 계정이 나열됩니다. ccam은 다음과 같은 계정을 건너뛰고, 그 이름과 이유를 알려 줍니다.

- `ccam import`로 제자리에 등록한 계정. 일반 `claude`가 그 로그인 파일을 쓰고 있을 수 있기 때문입니다.
- 실행 중인 계정. Claude Code가 로그인 파일을 다시 쓰기 때문입니다(기록 직전에 다시 확인하며, `started meanwhile`로 보고합니다).
- 아직 로그인하지 않은 계정.
- 로그인 정보를 macOS 키체인에 보관하는 계정. ccam은 키체인에 쓰지 않습니다.

건너뛴 계정에서는 `/mcp`로 서버에 로그인하세요. 나중에 MCP 서버를 공유하기 시작한 계정은 `ccam migrate`를 다시 실행하면 로그인 정보를 받습니다. 요약에도 ``Then run `ccam migrate` again to give them the MCP sign-ins too.``로 안내합니다.

## 안전장치와 예외 상황

- 원본은 읽기만 합니다. 그 안의 어떤 파일도 옮기거나 링크하거나 바꾸지 않으며, 원본의 로그인 정보에는 절대 쓰지 않습니다. 키체인에서 로그인 정보를 읽을 수 있는 macOS에서도 마찬가지입니다.
- `use mine`, 병합, 기존 `mcp.json`에 새 MCP 서버를 추가하는 경우처럼 그룹의 파일을 바꾸기 전에는 그룹의 사본을 `<group>/.ccam-backup/<time>/`에 저장합니다. 폴더 이름은 실행할 때마다 고유하므로, 같은 초에 두 번 실행해도 서로의 백업을 덮어쓰지 않습니다. 요약에 그 경로가 출력됩니다.
- `--yes`를 주면 모든 충돌에서 그룹의 사본을 유지하며, ccam은 직접 고르려면 같은 명령을 `--yes` 없이 실행하라고 안내합니다.
- 바뀔 것이 없으면 `Nothing to copy.`를, 모든 충돌에서 그룹의 사본을 유지하면 `Nothing to copy: the group keeps its own copies.`를 출력합니다.
- 복사가 중간에 실패하면 `The source folder was not changed. Group files replaced so far are saved in`과 백업 폴더를 출력하고 멈춥니다.
- 원본이 없으면 `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`로 끝납니다.

## 관련 문서

- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [점검과 복구](/features/doctor)
- [Claude Code 여러 계정 쓰기](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
