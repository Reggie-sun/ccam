---
title: 그룹으로 공유하기
description: ccam 그룹은 고른 항목을 폴더 하나로 계정끼리 공유하며, ccam group은 계획, 백업, 충돌 선택을 거쳐 데이터를 옮깁니다.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

계정마다 자체 폴더가 있지만, 그 안의 내용 대부분은 따로 둘 필요가 없습니다. **그룹**은 무엇을 누가 공유하는지 정합니다. 그룹의 모든 계정은 그 항목들을 폴더 하나(`~/.local/share/ccam/groups/<group>/`)에 링크하므로, 모두 같은 사본을 읽고 씁니다. 로그인 토큰은 절대 공유하지 않습니다.

한 계정이 여러 그룹에 속할 수 있지만, 각 항목은 한 그룹에서만 받습니다. 흔한 구성은 설정과 확장 기능을 위한 그룹 하나, 그리고 두 계정이 대화 내역을 나누는 그룹 하나입니다.

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

처음 추가한 계정이 `default`를 만들고, 그 뒤에 추가하는 계정은 따로 지정하지 않으면 이 그룹에 들어갑니다.

## 사용법

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

데이터를 옮기는 명령은 모두 계획을 출력하고, 그룹의 사본과 내용이 다른 파일마다 어떻게 할지 묻고, 무엇이든 바꾸기 전에 `Proceed? [y/N]`를 묻습니다.

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

[대시보드](/features/dashboard)의 **Groups** 탭에서도 같은 변경을 할 수 있으며, 계정에서 `s`를 누르면 그 계정의 그룹을 고릅니다.

## 항목

| 항목 | 다루는 범위 |
| --- | --- |
| `mcp` | MCP 서버. 그룹의 `mcp.json`으로 공유합니다([MCP 서버 공유](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | 전역 `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/`와 `history.jsonl`. 다른 계정에서 `--continue`할 수 있게 합니다 |

계정이 속한 어느 그룹도 공유하지 않는 항목은 계정 자체 폴더에 남습니다.

## 명령과 옵션

| 명령 | 동작 |
| --- | --- |
| `ccam group` 또는 `ccam group list` | 그룹 목록 보기 |
| `ccam group show <group>` | 그룹이 공유하는 항목, 계정, 폴더, MCP 서버 보기 |
| `ccam group add <group> --items <items> [--members <accounts>]` | 그룹 만들기. 그 계정들의 항목이 그룹으로 옮겨집니다 |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | 그룹이 공유하는 항목을 그룹의 모든 계정에 대해 바꾸기 |
| `ccam group join <group> <account>...` | 계정을 그룹에 넣기 |
| `ccam group leave <group> <account>...` | 계정을 그룹에서 빼기. 각 계정은 사본을 유지합니다 |
| `ccam group rename <group> <new-name>` | 설정에서만 이름 바꾸기 |
| `ccam group remove <group> [--leave]` | 그룹 삭제 |

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: 그룹에 먼저 사본을 넣을 계정 | 기본 계정, 그다음 번호순 |
| `--replace` | `join`: 같은 항목을 공유하는 그룹에서 먼저 나갑니다 | 꺼짐 |
| `--fresh` | 항목 공유를 멈추는 계정이 사본 없이 빈 상태로 시작합니다 | 꺼짐 |
| `--leave` | `remove`: 그룹의 계정을 먼저 뺍니다 | 계정이 있는 그룹이면 필요 |
| `--prefer <group\|mine\|merge>` | 충돌을 해결하는 방법 | 물어봄. `--yes`이면 파일별 기본값 |
| `-y`, `--yes` | 묻지 않고 실행합니다. 충돌은 `--prefer`를 따릅니다 | 꺼짐 |
| `--dry-run`(전역) | 계획을 출력하고 아무것도 바꾸지 않습니다 | 꺼짐 |

계정은 이름이나 번호로, 항목은 쉼표로 구분한 목록으로 지정합니다. `ccam add`는 `--group <group>`(여러 번 지정 가능)이나 `--no-group`을 받습니다.

## 동작 방식

관련 파일은 다음과 같습니다.

- `~/.config/ccam/config.toml`: `[[group]]` 항목. `name`, `items`, `members`와, 폴더 이름이 그룹 이름과 다를 때 `dir`이 들어 있습니다.
- `~/.local/share/ccam/profiles/<name>/`: 계정의 `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<folder>/`: 공유된 사본.

`config.toml`을 읽을 때 각 계정의 그룹을 항목마다 폴더 하나로 정리합니다. 계정 다섯 개 기준으로 약 25 µs가 걸리므로, 실행할 때는 여전히 `config.toml`만 읽습니다. 한 계정이 같은 항목을 두 그룹에서 받게 되는 설정은 오류이며, 고칠 때까지 아무것도 시작하지 않습니다.

변경은 설정의 사본에서 이루어집니다.

1. 명령이 사본을 고치고(`AddGroup`, `SetGroupItems`, `JoinGroup` …), 이 함수들이 규칙을 확인합니다.
2. `internal/group`이 이전 설정과 새 설정을 비교해, 항목이 바뀌는 계정마다 `share.Build`로 계획을 만듭니다. 계획은 하나씩 차례로 만들고, 각 계획은 앞선 계획이 그룹 폴더에 넣을 내용을 보므로, 자기 사본을 옮겨 넣는 것은 첫 번째 계정뿐입니다. 첫 사본은 그룹 폴더에 이미 있으면 그것을 쓰고, 없으면 `--from`, 기본 계정, 나머지 계정의 번호순으로 정합니다.
3. 그 계정 중 하나가 실행 중이면 ccam은 거부합니다. 개발자 모드가 꺼진 Windows에서 파일에 링크가 필요할 때도 거부합니다.
4. 프롬프트나 `--prefer`로 충돌을 해결하고, 계획을 출력한 뒤 확인을 받습니다.
5. 무엇이든 옮기기 전에 새 설정을 저장합니다. 단계가 중간에 멈추면 `ccam doctor --fix`나 다시 실행한 명령이 저장된 그룹을 향해 계획을 만들어 변경을 마무리합니다.
6. 계획을 순서대로 실행합니다. 모든 단계는 데이터를 그룹이나 백업 폴더로 옮기며, 삭제되는 것은 없습니다.

`share.Layout`은 각 항목이 다루는 경로의 허용 목록입니다. 목록에 없는 것은 절대 그룹에 들어갈 수 없습니다. `share.Never`에는 `.credentials.json`과 `.claude.json`이 들어 있으며, 둘 중 하나라도 레이아웃에 나타나면 `Build`가 내부 오류를 내고 멈춥니다. `mcp` 항목에는 경로가 없습니다. 서버 목록이 `.claude.json` 안에 있기 때문에 [MCP 공유](/features/mcp)는 그룹의 `mcp.json`으로 동작합니다.

`side`가 대화 내역을 공유하기 시작하고 그룹에 이미 내용이 있을 때, `projects/` 폴더 하나에 일어나는 일은 다음과 같습니다.

<!-- motion:share-on -->

그룹의 그 경로에 아직 아무것도 없으면 이 계정이 처음 공유하는 것이므로, 계정의 폴더를 그대로 그곳으로 옮깁니다. 그렇지 않으면 그룹에 없는 파일을 하나씩 그룹으로 옮기고, 양쪽에 모두 있는 파일은 계정 폴더에 남겨 그룹의 버전을 유지합니다. 그다음 계정의 폴더를 계정 폴더 안의 `.ccam-backup/<time>/`으로 옮기고, 그 자리에 링크를 둡니다. 대화 파일 이름은 세션 ID를 따르므로, 여러 계정의 대화 내역을 합쳐도 충돌하지 않습니다. `plugins/` 레지스트리는 계정 폴더를 거치는 절대 경로를 저장합니다. ccam은 이 경로를 그룹 폴더를 가리키도록 바꾸고, 그룹에 이미 레지스트리가 있으면 계정의 항목을 그 안에 병합합니다.

그룹에서 나가면 링크를 되돌립니다.

<!-- motion:share-off -->

ccam은 링크가 그룹 폴더 안을 가리킬 때만 링크를 제거하고, 그룹의 현재 내용을 계정에 복사합니다. `--fresh`를 주면 복사를 건너뜁니다. 그룹의 내용은 다른 계정을 위해 그대로 남습니다. 폴더가 아직 비어 있는 그룹으로 항목을 옮기는 계정은 자신이 보던 내용의 사본을 가져가므로 파일을 잃지 않습니다. MCP 서버도 같은 방식으로 따라갑니다. MCP 공유를 멈춘 계정은 ccam이 `.claude.json`에 복사해 둔 서버를 잃습니다(백업은 남습니다). 그 서버는 그룹에서 온 것이기 때문입니다.

양쪽에 모두 있지만 내용이 다른 단일 파일은 충돌입니다.

<!-- motion:share-conflict -->

`history.jsonl`은 `merge by line`, `keep group`, `use mine`을 제시하며 병합이 맨 처음입니다. 그 밖의 `.json` 파일은 `keep group`, `use mine`, `merge JSON`을, 나머지 파일은 `keep group`과 `use mine`을 제시합니다. `merge JSON`은 객체를 재귀적으로 병합하고, 배열은 중복 없이 합치며, 양쪽이 같은 키에 다른 값을 설정했으면 그룹의 값을 유지합니다. 그래서 다른 계정에는 아무 변화가 없습니다. `keep group`을 제외한 모든 선택지는 먼저 그룹의 사본을 계정 폴더의 `.ccam-backup/<time>/group/`에 저장합니다. 내 파일은 항상 백업으로 보관됩니다.

실행할 때마다 `share.Check`가 레이아웃의 각 경로를 계정의 그룹과 비교하며, 이때 `lstat`과 `readlink`만 씁니다. 없는 링크는 그 자리에서 만듭니다. 링크가 있어야 할 자리에 일반 파일이 있거나, 링크가 잘못된 곳을 가리키거나, 더 이상 쓰지 않는 링크가 남아 있으면 실행을 멈추고 `ccam doctor --fix`를 안내합니다.

### 이름 바꾸기, 삭제, 이전 형식

- `ccam group rename`은 설정만 바꿉니다. 폴더는 이름을 그대로 유지하고 그룹이 그 이름을 `dir`로 기록하므로, 옮겨지는 것도 잃는 것도 없습니다.
- `ccam group remove`는 `--leave`를 주지 않으면 계정이 있는 그룹을 삭제하지 않습니다. 주면 각 계정이 사본을 유지하고, 다른 그룹이 쓰지 않는 한 폴더는 `~/.local/share/ccam/trash/`로 옮겨집니다.
- 그룹이 생기기 전에는 `config.toml`이 계정마다 `[profile.share]` 테이블을 저장했습니다. ccam은 이 형식을 아무것도 옮기지 않고 그룹으로 읽습니다. 처음 저장할 때 새 형식으로 기록하고 이전 파일은 `config.toml.bak`으로 남깁니다. 같은 계정들이 공유하던 항목은 폴더 이름을 딴 그룹 하나가 되고, 나머지는 `default-history` 같은 이름을 받아 같은 폴더를 계속 씁니다.

### Windows에서

Windows는 개발자 모드(설정 > 시스템 > 개발자용)이거나 관리자일 때만 심볼릭 링크를 만들 수 있습니다. 그 권한이 없으면 공유 폴더는 디렉터리 정션이 되며, Claude Code에서는 똑같이 동작합니다. 공유 파일(`settings.json`, `CLAUDE.md`, `history.jsonl`)은 심볼릭 링크로만 만들 수 있으므로, 그런 링크가 필요한 그룹 변경은 무엇이든 옮기기 전에 거부됩니다.

## 안전장치와 예외 상황

- **실행 중인 계정.** 옮기려는 계정 중 하나로 Claude Code가 실행 중이면 변경을 거부합니다. 세션이 그 파일에 계속 쓰기 때문입니다. 그룹에서 나가거나 들어갈 때는 지정한 계정만 해당합니다. `--dry-run`은 그래도 계획을 출력합니다.
- **백업.** 계획마다 계정 폴더 안에 자체 `.ccam-backup/<time>/` 폴더가 생깁니다. ccam은 이 폴더를 절대 삭제하지 않습니다.
- **직접 만든 링크.** 어떤 경로가 이미 직접 만든 링크이고(예: dotfiles 저장소를 가리키는 링크) 그룹이 비어 있으면, ccam은 링크 대상의 내용을 그룹에 복사하고 링크 자체는 백업으로 옮깁니다. 링크 대상은 건드리지 않습니다.
- **가져온 폴더.** `ccam import`로 제자리에 등록한 계정(`~/.claude` 등)은 그 자리에서 바뀌므로, 그냥 `claude`를 실행해도 변경이 보입니다. 계획에 그 사실이 표시됩니다.
- **같은 폴더에서 두 계정 쓰기.** 대화 내역을 공유하는 계정들은 같은 파일에 기록합니다. 한 계정이 이미 실행 중인 폴더에서 두 번째 계정을 시작하면 ccam이 경고하고 계속할지 묻습니다.
- **직접 고친 설정.** 계정이 아닌 멤버는 무시되며 `ccam doctor`가 알려 줍니다. `local`, `on`, `off`, `none`은 그룹 이름으로 쓸 수 없습니다.

## 관련 문서

- [대시보드](/features/dashboard)
- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [MCP 서버 공유](/features/mcp)
- [기존 설정 복사하기](/features/migrate)
- [전체 명령](/docs/commands)
