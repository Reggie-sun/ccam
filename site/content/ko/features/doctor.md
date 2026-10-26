---
title: 점검과 복구
description: ccam doctor로 공유 링크, 파일 권한, 로그인, Claude Code 버전을 점검하고, ccam doctor --fix로 안전하게 고칠 수 있는 문제를 복구합니다.
updated: 2026-10-26
order: 9
translation: machine
sourceHash: "282972ac5053"
---

`ccam doctor`는 ccam 설치 상태를 점검합니다. Claude Code 실행 파일, 모든 계정의 폴더, 로그인, 공유 링크, 공유 MCP 서버, 사용량 엔드포인트, 단축 명령을 확인합니다. `ccam doctor --fix`는 데이터를 잃지 않고 고칠 수 있는 문제를 복구하며, 교체하는 파일은 먼저 백업 폴더로 옮깁니다. `ccam <name>`이 공유 링크 문제로 멈췄을 때, 파일을 직접 복원한 뒤, 또는 Claude Code를 업데이트한 뒤에 실행하세요.

## 사용법

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

첫 줄은 모든 계정을 한꺼번에 표시로 보여 주므로, 다른 내용을 읽기 전에 어디에 문제가 있는지 알 수 있습니다. 그 아래에는 계정마다 상자가 하나씩 있고, 어느 계정에도 속하지 않는 결과(Claude Code, 환경 변수, 플러그인, MCP 서버, 단축 명령, 폴더 규칙, 그룹, 설정, `auto_switch`)는 `General` 상자에 들어갑니다. 상자 제목에는 그 계정의 가장 심각한 표시와 문제 개수가 나옵니다. 상자 안의 각 행은 점검한 영역 하나이며, 항상 `Folder`, `Login`, `Share links` 순서로 나오고 그다음에 나머지가 나옵니다. 앞에 붙는 표시는 `✔` 정상, `ℹ` 정보, `⚠` 경고, `✘` 오류입니다. 아무 문제가 없는 폴더는 `✔ Folder  OK`로 표시됩니다.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix`를 붙이면 복구된 항목 끝에 `→ fixed`가 붙고, 그 항목은 더 이상 문제로 세지 않습니다.

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

터미널에서는 상자가 터미널 너비에 맞춰지고 긴 메시지는 줄바꿈되며, 터미널 자체의 색을 씁니다(`NO_COLOR`로 끌 수 있습니다). 상자는 29열보다 좁아지지 않으므로, 좁은 터미널에서도 단어가 잘리지 않습니다. 파이프나 파일로 보내면 색이 없고 어떤 메시지도 줄바꿈되거나 잘리지 않으므로, `ccam doctor | grep settings.json`으로 줄 전체를 찾을 수 있습니다. 이때 모든 상자의 너비는 가장 긴 줄에 맞춰집니다. 대시보드의 Health 탭도 같은 상자를 보여 줍니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--fix` | 안전하게 고칠 수 있는 문제를 복구합니다(바뀌는 파일은 백업) | 꺼짐 |
| `--dry-run`(전역) | 모든 점검을 실행하지만, `--fix`를 붙여도 아무것도 복구하지 않습니다 | 꺼짐 |

`ccam doctor`는 인자를 받지 않습니다. `--dry-run`을 지원하므로, `ccam --dry-run doctor --fix`는 `ccam doctor`의 결과에 `--fix`가 따로 치워 둘 것을 더해 출력하고, 아무것도 바꾸지 않습니다.

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

복구되지 않은 `✘` 항목이 하나라도 남으면 종료 코드는 `1`이고, 그렇지 않으면 `0`입니다. 경고와 정보는 실패로 치지 않습니다. `--fix` 없이 실행해 실패하면 ccam은 빈 줄과 `` `ccam doctor --fix` repairs what it can. ``으로 끝납니다. 설정을 읽을 수 없을 때도 `ccam:` 오류와 함께 `1`로 끝납니다.

## 동작 방식

점검은 정해진 순서로 실행되며, 결과는 계정마다 상자 하나에 나타납니다. ccam은 `~/.config/ccam/config.toml`을 읽고(`internal/cli/doctor.go`) 점검을 실행하며, 각 결과에는 속한 계정과 점검한 영역이 담깁니다(`internal/doctor/doctor.go`의 `doctor.Run`). 그다음 결과를 상자로 묶습니다(`tui.HealthReport`). `--fix`는 `Options.Fix`로 전달되며, 시험 실행에서는 false가 됩니다.

<!-- motion:doctor-checks -->

가운데 블록은 계정마다 한 번씩 실행됩니다. 각 점검과 `--fix`가 하는 일은 다음과 같습니다.

1. **Claude Code 버전**(`checkClaude`). 실행기와 같은 방식으로 `claude` 실행 파일을 찾아 `claude --version`을 실행합니다. 실행 파일이 없거나, `claude_bin`이 실행 가능한 파일이 아니거나, `--version`이 실패하면 오류입니다. 복구하지 않습니다.
2. **환경 변수의 로그인 정보**(`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN` 중 하나가 설정되어 있으면 경고합니다. 그 셸에서 시작하는 모든 계정이 이 값을 물려받고, Claude Code는 계정 자체의 로그인보다 이 값을 우선하기 때문입니다. 복구하지 않으니 직접 해제하세요.
3. **폴더**(`checkProfile`). 계정 폴더가 없으면 오류이며, 그 계정의 점검은 여기서 끝납니다. 그룹이나 기타 사용자 권한 비트가 하나라도 설정되어 있으면 경고이고, `--fix`가 `chmod 0700`을 실행합니다. 가져온 `~/.claude`에 `~/.claude.json`과 별도로 자체 `.claude.json`이 있으면 정보로 표시합니다. 이 경우 `ccam main`과 일반 `claude`가 로그인 정보와 MCP 서버를 따로 갖게 됩니다. Windows에는 권한 비트가 없으므로 이 점검과 아래의 `0600` 점검은 건너뜁니다. 사용자 프로필 안의 폴더는 액세스 목록으로 보호됩니다.
4. **로그인 파일**(`checkCredentials`). `.credentials.json`을 `lstat`으로 확인합니다. 계정마다 자체 토큰이 있어야 하므로 심볼릭 링크이면 오류이며, 복구하지 않습니다. 권한이 `0600`이 아니면 경고이고, `--fix`가 `chmod 0600`을 실행합니다.

   그다음 로그인 정보를 읽습니다(macOS에서는 키체인을 먼저 봅니다). 로그인 정보가 없거나 로그인 파일에 로그인이 들어 있지 않으면 `ccam login <name>`을 안내하는 경고입니다. 리프레시 토큰이 만료되었거나, 서버가 거부한 뒤 Claude Code가 로그인을 지웠으면 `` login expired → `ccam login <name>` `` 오류이고, 액세스 토큰만 만료되었으면 Claude Code가 다음 실행 때 갱신하므로 정보로만 표시합니다. Windows 자격 증명 관리자에 보관된 로그인은 ccam이 읽을 수 없으므로 정보로 표시합니다. 로그인이 유효하면 폴더가 지금 로그인된 이메일과 요금제를 보여 줍니다(`login valid (dev@company.com, Claude Team · Acme)`). 이 중 어느 것도 복구하지 않습니다. 로그인이 필요한 모든 계정에는 `ccam login --expired`로 로그인할 수 있습니다.
5. **공유 링크**(`share.Check`). `share.Layout`의 모든 경로를 계정의 그룹과 비교하며, `lstat`과 `readlink`만 사용합니다. 링크가 없는 경우, 링크가 다른 곳을 가리키는 경우, 링크가 있어야 할 자리에 일반 파일이 있는 경우, 이제 로컬인 항목에 링크가 남아 있는 경우, 링크는 맞지만 링크가 가리키는 그룹 폴더가 없는 경우(`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`. 실행할 때는 이것을 확인하지 않습니다)의 다섯 가지를 보고합니다. 링크가 없는 경우와 그룹 폴더가 없는 경우는 경고이고, 나머지는 오류입니다.

   그룹에서 다시 복사하다 중간에 멈춘 경우도 보고합니다. `--fix`는 준비해 둔 완전한 사본을 제자리에 두지만, 그 사이에 실행이 그 경로를 다시 만들었다면 완전한 사본이 있는 폴더를 알려 주어 사용자가 직접 정리하게 합니다. 이후 버전의 ccam에서 항목에 추가된 경로(대화 내역의 `file-history/`, `tasks/`, `paste-cache/`, 설정의 `keybindings.json`)는 계정에 아직 자체 사본이 있는 동안 경고로만 표시됩니다: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. 이 경로 때문에 실행이 멈추는 일은 없습니다.
6. **알 수 없는 항목**(`checkUnknownEntries`). 계정 폴더 최상위에서 ccam이 모르는 이름(`share.Known`에 없는 이름)을 나열해, Claude Code 폴더 구조가 바뀌었을 수 있다는 단서를 줍니다. 정보로만 표시합니다.
7. **플러그인 경로**(`checkPlugins`). 플러그인 레지스트리(`installed_plugins.json`, `known_marketplaces.json`)는 절대 경로를 저장합니다. 각 그룹의 `plugins/`에서 특정 계정 폴더를 거치는 경로는 그 계정을 삭제하거나 이름을 바꾸면 깨집니다. `--fix`는 `profile.RewritePluginPaths`로 이 경로를 그룹 폴더 기준으로 다시 씁니다. 레지스트리 파일은 원자적으로 기록됩니다.
8. **MCP 서버**(`checkMCP`). 그룹마다 `mcp.json`을 한 번 읽습니다. 파일을 읽을 수 없으면 오류이고, `mcp.Validate`를 통과하지 못한 서버는 경고입니다. 일반 `claude`의 `~/.claude.json`에 있지만 어느 계정도 불러오지 않는 서버는 경고이며, 이를 복사하는 `ccam migrate --mcp` 명령을 알려 줍니다(여러 그룹이 MCP를 공유하면 `--group <a|b>`도 붙입니다). 그룹에서 한 계정에만 있는 서버는 이를 공유하는 `ccam mcp import` 명령과 함께 정보로 표시하고, 다른 계정은 로그인한 MCP 서버에 로그인하지 않은 계정에는 `/mcp`로 로그인하라는 경고를 보여 줍니다(ccam은 로그인을 절대 복사하지 않습니다). 복구하지 않습니다.
9. **사용량 응답 형식**(`checkUsageFormat`). 캐시 없이, 전체 4초 타임아웃 안에서 계정을 하나씩 돌아가며 응답이 올 때까지 사용량 엔드포인트에 요청합니다. `usage_source = "claude"`이면 대신 한 계정에 대해 Claude Code에 묻습니다. 응답을 더 이상 해석할 수 없으면 `ccam usage`가 상태 표시줄 캐시를 대신 쓴다고 경고하고, 응답한 계정이 없으면 그 사실을 정보로 알립니다(`could not check the usage response format: …`). 네트워크를 쓰는 점검은 이것뿐입니다.
10. **단축 명령**(`checkShims`). 단축 명령이 있는 계정이 있을 때만 실행합니다. `~/.local/bin`에 링크가 없으면 경고이고, ccam이 만든 단축 명령이 교체되거나 옮겨진 이전 버전의 ccam을 아직 실행해도 경고입니다(`shortcuts in <bin> that run another ccam (one replaced or moved): …`). `--fix`가 `shim.Install`을 호출하며, 이 함수는 ccam이 만들어 `~/.local/state/ccam/shortcuts.json`에 기록한 링크만 만들고 새로 고칠 뿐 다른 파일은 절대 건드리지 않습니다. 경고 끝에 `→ fixed`가 붙는 것은 그 링크를 모두 만들었을 때뿐이며, 하나라도 만들지 못하면 이유가 덧붙습니다. 아직 설치되어 있는 다른 ccam을 실행하는 단축 명령은 동작하므로 `ccam doctor`는 정보로만 표시하고, `--fix`가 이 ccam 실행 파일을 가리키게 바꿉니다. bin 폴더가 `PATH`에 없으면 경고이며, 복구하지 않습니다.
11. **폴더별 계정**(`checkFolders`). [`ccam folder`](/features/folder) 항목이 있을 때만 실행합니다. 항목의 계정이 존재하지 않으면 오류이고, 항목의 폴더가 사라졌으면 경고입니다. 항목을 지우면 그 폴더에서 시작되는 계정이 바뀌므로, 둘 다 복구하지 않습니다.
12. **그룹과 설정**(`checkGroups`). 계정이 아닌 그룹 멤버, 이름이 그 뒤에 ccam 명령 이름이 되어 버린 계정(`ccam run <name>`으로 시작하거나 이름을 바꾸세요), 알 수 없는 `theme`이나 `usage_source` 값, ccam이 읽지 않는 `config.toml`의 줄은 경고입니다. 마지막 경우는 대개 `[[group]]`, `[[profile]]`, `[[folder]]` 테이블 아래에 적어 그 테이블에 속하게 된 설정입니다(`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`). 이를 고칠 때까지 설정을 바꾸는 명령은 거부합니다. 복구하지 않습니다.
13. **로그인**(`checkIdentities`). 각 폴더의 `.claude.json`에서 어느 계정으로 로그인되어 있는지 읽습니다. 폴더가 `ccam login`이 기록한 것과 다른 계정으로 로그인되어 있으면(예: 세션 안에서 `/login`을 실행한 뒤) 경고입니다. 가져온 계정이나 이전 버전의 ccam이 만든 계정처럼 로그인 기록이 없는 계정은 기록된 이메일로 비교합니다. 다른 폴더와 같은 Claude 계정으로 로그인되어 있어도 경고입니다. 두 폴더가 사용량 한도를 함께 쓰기 때문입니다. 다른 폴더와 같은 로그인 토큰을 가진 폴더, 즉 복사한 폴더(일반 `claude`의 `~/.claude` 포함)는 오류입니다. 먼저 갱신하는 쪽이 다른 쪽을 로그아웃시키기 때문이며, `ccam login <name>`을 실행해 자체 로그인을 주세요. 복구하지 않습니다.
14. **`auto_switch`**(`checkAutoSwitch`). `auto_switch = true`일 때만 점검합니다. `General` 상자의 정보 줄이 `auto_switch`가 하는 일을 설명하고 [그 설명](/features/next#auto-switch)으로 연결합니다. 계정 가운데 시트가 있는 요금제(Claude Team, Claude Enterprise)마다 한 줄이 그 계정들의 이름을 대고, `auto_switch`가 한 조직의 시트 사이에서 세션을 옮긴다는 것을 알리며, 고객센터 링크와 함께 Anthropic의 시트 규칙을 인용하고, 선택은 사용자에게 맡깁니다. 훅에 자기 프로세스(`CLAUDE_PID`)를 알려 주지 않는 2.1.214보다 오래된 Claude Code는 세션을 옮길 수 없으므로 그 상자에서 경고입니다.

    ccam 플러그인이 없거나 플러그인이 꺼진 계정은 자기 상자에 경고가 붙고, Claude Code에서 이를 고치는 명령이 함께 나옵니다: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam`(또는 `/plugin enable ccam@ccam`). API 키 계정은 제외합니다. Windows에서는 `auto_switch is ignored on Windows; with relay = true, /ccam:switch followed by /exit still moves a conversation to another account.` 한 줄만 나옵니다. 복구하지 않습니다.

계정을 점검하기 전에, 중간에 멈춘 **그룹 변경**은 `General` 상자의 오류로 표시됩니다: ``the last group change stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix`는 계획한 순서대로, 충돌에 대해 받은 답을 써서 이를 마무리하며, 그때까지 다른 그룹 변경, 이름 바꾸기, 삭제는 시작하지 않습니다. 변경이 끝날 때까지 ccam은 그 변경에 포함된 계정에 없는 링크만 만듭니다. 마무리할 수 없으면, 그 파일을 지우면 나머지 변경을 포기하게 된다고 오류에서 알려 줍니다.

`ccam --dry-run doctor --fix`는 순서(`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`)와 계정마다 따로 치워 둘 것을 보여 줍니다. 다른 ccam 프로세스가 지금 진행 중인 그룹 변경은 정보로 표시하며, 진행 중인 동안 `--fix`는 공유 링크를 복구하지 않습니다(`its share links are not repaired while another ccam changes the groups; …`). `--fix`는 복구하는 동안 `~/.local/state/ccam/group-change.lock`을 잡고 있으며, 변경을 마무리한 뒤 `config.toml`을 다시 읽습니다.

### 공유 링크 복구

실행 경로는 없는 링크를 그 자리에서 만들지만, 데이터를 잃을 수 있는 문제에서는 멈추고 이 명령을 안내합니다.

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

공유 링크 문제가 있으면 `--fix`는 계정의 공유 설정을 현재 상태이자 목표 상태로 넘겨 `share.Build`를 호출한 다음 `Plan.Apply`를 실행합니다. ccam은 디스크의 실제 상태로부터 모든 단계를 다시 만들어 `ccam group` 변경이 쓰는 것과 같은 계획을 얻으므로(`Build`), 복구는 그룹에 들어갈 때와 같은 규칙을 따릅니다.

<!-- motion:doctor-repair -->

이 예에서 `settings.json`은 일반 파일입니다. 그룹에 이미 내용이 다른 `settings.json`이 있으므로, 충돌은 기본 선택인 `keep group`(그림의 `KeepGroup`)으로 처리되어 그룹의 사본이 유지됩니다. 계정의 파일은 계정 폴더 안의 `.ccam-backup/<time>/settings.json`으로 옮겨지고, ccam이 그 경로를 그룹에 링크합니다(`os.Symlink`). `--fix`는 이를 알려 줍니다: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. 다른 경우는 다음과 같습니다.

- 그룹에 아직 사본이 없으면 계정의 파일이 그룹으로 옮겨져 그룹의 첫 내용이 됩니다.
- 폴더라면 그룹에 없는 파일은 그룹으로 옮깁니다. 그룹의 사본과 똑같은 파일은 지우고, 내용이 다른 파일은 백업으로 보내며 메모로 알려 줍니다(`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- 그룹 폴더가 없으면 먼저 만듭니다.
- `history.jsonl`은 기본값이 `merge by line`입니다. 그룹의 줄 다음에 새 줄을 붙입니다. 그룹의 이전 사본은 `.ccam-backup/<time>/group/`에 보관됩니다.
- 다른 그룹을 가리키는 링크는 교체합니다. ccam 밖의 폴더를 가리키는 링크는 백업으로 보내고, 그룹이 비어 있으면 그 내용으로 그룹을 채웁니다.
- 이제 로컬인 항목에 남은 링크는 제거하고, 계정은 그룹 내용의 사본을 받습니다.

링크를 제거하고 그룹의 사본과 똑같은 파일을 지우는 것 말고는 모든 단계가 데이터를 그룹이나 백업 폴더로 옮깁니다. 잃는 것은 없습니다.

## 안전장치와 예외 상황

- 그냥 `ccam doctor`를 실행하면 아무것도 바꾸지 않습니다. 스크립트에서 확실히 하려면 `--fix`와 함께 `--dry-run`을 쓰세요.
- 실행 중인 계정에 대해서는 Claude Code를 어떻게 시작했든 `doctor --fix`가 없는 링크와 없는 그룹 폴더만 만듭니다. 이것으로는 잃는 것이 없습니다. 그 밖에 남은 것이 있으면 `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``를 덧붙입니다. 복구하면서 옮길 파일에 세션이 계속 쓰기 때문입니다.
- 공유 복구가 실패하면 ccam은 그 계정에 대해 `could not repair share links: <error>`를 보고하고 알 수 없는 항목 점검을 건너뜁니다. 다른 계정은 계속 점검합니다.
- `ccam group` 변경이 중간에 실패해도 새 그룹은 이미 저장되어 있고 그 변경은 `~/.local/state/ccam/group-change.json`에 기록됩니다. ccam은 `ccam doctor --fix`로 마무리하라고 안내하며, 이 명령이 MCP 단계를 포함해 나머지 변경을 실행합니다.
- 심볼릭 링크인 `.credentials.json`, 만료된 로그인, 환경 변수의 로그인 정보는 보고만 하고 절대 바꾸지 않습니다.
- 플러그인 레지스트리를 다시 쓸 때는 원자적으로 기록하지만, `.ccam-backup/`에 복사해 두지는 않습니다.

## 관련 문서

- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [기존 설정 복사하기](/features/migrate)
- [사용량 한도와 계정 전환](/docs/usage-limits)
- [명령어 레퍼런스](/docs/commands)
