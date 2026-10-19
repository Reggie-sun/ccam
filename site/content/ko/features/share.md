---
title: 계정 간 공유
description: ccam share가 계획, 백업, 충돌 선택을 거쳐 설정, CLAUDE.md, 확장 기능, 대화 내역을 그룹 폴더에 링크하는 방법과, off로 끌 때 내용을 다시 복사해 오는 방법을 설명합니다.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

계정마다 자체 폴더가 있지만, 그 안의 내용 대부분은 따로 둘 필요가 없습니다. `ccam share`는 한 계정의 한 항목씩 공유를 켜거나 끕니다. 계정의 파일이나 폴더를 그룹 폴더로 옮기고 그 자리를 심볼릭 링크로 바꾸므로, 그룹의 모든 계정이 같은 사본을 읽고 씁니다. `ccam add` 이후에 계정이 공유하는 항목을 바꾸거나, 일부 계정만 따로 그룹으로 묶을 때 쓰세요.

## 사용법

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

항목 옵션 없이 실행하면 `ccam share`는 현재 상태를 출력합니다.

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

옵션을 주면 ccam은 계획을 출력하고, 그룹의 사본과 내용이 다른 파일마다 어떻게 할지 묻고, 무엇이든 바꾸기 전에 `Proceed? [y/N]`를 묻습니다. 새 계정은 대화 내역을 제외한 모든 항목을 공유합니다. 각 항목을 사용자 입장에서 설명한 내용은 [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)를 참고하세요.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP 서버 | 변경 없음 |
| `--settings <on\|off\|group>` | `settings.json` | 변경 없음 |
| `--memory <on\|off\|group>` | 전역 `CLAUDE.md` | 변경 없음 |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | 변경 없음 |
| `--history <on\|off\|group>` | `projects/`와 `history.jsonl` | 변경 없음 |
| `--fresh` | 공유를 끌 때 그룹의 내용을 복사하지 않고 빈 상태로 시작합니다 | 꺼짐 |
| `--prefer <group\|mine\|merge>` | 충돌을 해결하는 방법 | 물어봄. `--yes`이면 항목별 기본값 |
| `-y`, `--yes` | 묻지 않고 실행합니다. 충돌은 `--prefer`를 따릅니다 | 꺼짐 |
| `-h`, `--help` | `share` 도움말 | |
| `--dry-run`(전역) | 지원합니다. 충돌을 묻고 계획을 출력하지만 아무것도 바꾸지 않습니다 | 꺼짐 |

값 `on`은 `default` 그룹, `off`는 로컬(공유 안 함)을 뜻하며, 그 밖의 값은 그룹 이름입니다. 그룹 이름은 `~/.local/share/ccam/groups/` 아래 폴더 이름이 되므로, `ccam share side --history ../x`는 `invalid group name: "../x"`로 실패합니다. 계정은 이름이나 번호로 지정하며, 지정한 항목만 바뀝니다.

## 동작 방식

관련 파일은 다음과 같습니다.

- `~/.config/ccam/config.toml`: 계정마다 있는 `[profile.share]` 테이블. 항목마다 값이 하나씩 있습니다.
- `~/.local/share/ccam/profiles/<name>/`: 계정의 `CLAUDE_CONFIG_DIR`.
- `~/.local/share/ccam/groups/<group>/`: 공유된 사본.

1. `internal/cli/share.go`가 설정을 읽고, `config.Share.Set`으로 계정 설정의 사본에 각 `--<item>` 값을 적용하며, 계정이 실행 중이면 거부합니다.
2. `share.Build`가 현재 설정을 목표 설정과 비교해 단계를 나열합니다. 디스크는 읽기만 하며, `share.Layout`의 모든 경로에 `os.Lstat`과 `os.Readlink`를 호출합니다. 계획은 작업 목록입니다. 비어 있는 그룹 자리로 옮기기(adopt), 그룹에 없는 파일 옮기기(merge dir), 해결한 충돌을 그룹에 기록하기, 백업, 링크 해제, 링크, 다시 복사해 오기, 플러그인 경로 다시 쓰기, 플러그인 레지스트리 병합이 있습니다.
3. 프롬프트나 `--prefer`로 충돌을 해결합니다.
4. `Plan.Describe`가 계획을 출력하고 ccam이 확인을 받습니다.
5. 무엇이든 옮기기 전에 목표 설정을 `config.toml`에 기록합니다. `Plan.Apply`가 중간에 멈추면 `ccam doctor --fix`나 다시 실행한 명령이 저장된 목표를 향해 계획을 만들어 변경을 마무리합니다.
6. `Plan.Apply`가 단계를 순서대로 실행합니다. 모든 단계는 데이터를 그룹이나 백업 폴더로 옮기며, 삭제되는 것은 없습니다.

`share.Layout`은 각 항목이 다루는 경로의 허용 목록입니다. 목록에 없는 것은 절대 그룹에 들어갈 수 없습니다. `share.Never`에는 `.credentials.json`과 `.claude.json`이 들어 있으며, 둘 중 하나라도 레이아웃에 나타나면 `Build`가 내부 오류를 내고 멈춥니다. `mcp` 항목에는 경로가 없습니다. 서버 목록이 `.claude.json` 안에 있기 때문에 [MCP 공유](/features/mcp)는 그 대신 그룹의 `mcp.json`으로 동작하며, `--mcp`는 설정만 바꿉니다.

`side`가 대화 내역 공유를 켜고 그룹에 이미 내용이 있을 때, `projects/` 폴더 하나에 일어나는 일은 다음과 같습니다.

<!-- motion:share-on -->

그룹의 그 경로에 아직 아무것도 없으면 이 계정이 처음 공유하는 것이므로, 계정의 폴더를 백업 없이 그대로 옮깁니다. 그렇지 않으면 `mergeDir`가 그룹에 없는 파일을 하나씩 그룹으로 옮기고, 양쪽에 모두 있는 파일은 계정 폴더에 남겨 그룹의 버전을 유지합니다. 그다음 계정의 폴더를 계정 폴더 안의 `.ccam-backup/<time>/`으로 옮기고, `os.Symlink`가 그 자리에 링크를 둡니다. 대화 파일 이름은 세션 ID를 따르므로, 여러 계정의 대화 내역을 합쳐도 충돌하지 않습니다.

옮길 때는 `rename(2)`를 쓰고, 파일 시스템이 다르면 복사 후 삭제로 대신합니다. 옮긴 트리 안에서 트리 밖을 가리키던 상대 링크는 절대 링크로 다시 쓰므로(`fsutil.PortableLink`), 새 위치에서도 올바르게 해석됩니다. `plugins/` 레지스트리(`installed_plugins.json`, `known_marketplaces.json`)는 계정 폴더를 거치는 절대 경로를 저장합니다. ccam은 이 경로를 그룹 폴더를 가리키도록 바꾸고, 그룹에 이미 레지스트리가 있으면 계정의 항목을 그 안에 병합합니다.

항목의 공유를 끄면 링크를 되돌립니다.

<!-- motion:share-off -->

ccam은 링크가 그룹 폴더 안을 가리킬 때만 링크를 제거하고, `CopyTree`로 그룹의 현재 내용을 계정에 복사합니다. `--fresh`를 주면 복사를 건너뛰어, 계정은 그 경로에 아무것도 없는 상태로 시작합니다. 그룹의 내용은 다른 계정을 위해 그대로 남습니다. 항목을 한 그룹에서 다른 그룹으로 옮기면 링크를 해제하고 새 그룹에 링크하며, 이전 그룹은 그대로 둡니다. 계획에 `Note:` 줄로 이 사실이 표시됩니다.

양쪽에 모두 있지만 내용이 다른 단일 파일은 충돌입니다.

<!-- motion:share-conflict -->

`share.NewConflict`는 파일에 맞는 선택지를 나열합니다. `history.jsonl`은 `merge by line`, `keep group`, `use mine`을 제시하며 병합이 맨 처음입니다. 그 밖의 `.json` 파일은 `keep group`, `use mine`, `merge JSON`을, 나머지 파일은 `keep group`과 `use mine`을 제시합니다. 맨 처음 선택지가 기본값입니다.

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON`(`MergeJSONBytes`)은 객체를 재귀적으로 병합하고, 배열은 중복 없이 합치며, 양쪽이 같은 키에 다른 값을 설정했으면 그룹의 값을 유지합니다. 그래서 다른 계정의 동작은 바뀌지 않습니다. `merge by line`은 그룹의 줄을 유지하고, 그 뒤에 그룹에 없는 내 줄을 추가합니다. `keep group`을 제외한 모든 선택지는 먼저 그룹의 사본을 계정 폴더의 `.ccam-backup/<time>/group/`에 저장한 뒤 결과를 `fsutil.AtomicWrite`로 기록합니다. 내 파일은 항상 백업으로 보관됩니다.

`--prefer group`, `mine`, `merge`를 주면 묻지 않고 모든 충돌을 해결합니다. `merge`는 파일이 제시하는 쪽에 따라 `merge JSON`이나 `merge by line`을 뜻하며, 둘 다 제시하지 않는 파일은 기본값을 따릅니다. `--yes`만 주면 파일마다 기본값을 씁니다.

실행할 때마다 `share.Check`가 레이아웃의 각 경로를 설정과 비교하는데, `lstat`과 `readlink`만 쓰므로 실행 경로에서도 부담이 적습니다. 없는 링크는 잃을 것이 없으므로 그 자리에서 만듭니다. 링크가 있어야 할 자리에 일반 파일이 있거나, 링크가 잘못된 곳을 가리키거나, 더 이상 쓰지 않는 링크가 남아 있으면 실행을 멈춥니다.

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix`는 저장된 설정을 현재와 목표 양쪽에 넘겨 `share.Build`를 실행합니다. 디스크의 실제 상태로부터 모든 단계를 만들어 내고, 충돌하면 그룹의 내용을 유지하며, 바꾸는 것은 모두 백업합니다.

## 안전장치와 예외 상황

- **실행 중인 계정.** 그 계정으로 Claude Code가 실행 중이면 ccam은 공유 설정 변경을 거부합니다(`account work is running; quit it first, then change its sharing`). 세션이 옮기려는 파일에 계속 쓰므로, 그 사이에 추가된 줄을 잃게 되기 때문입니다. `--dry-run`은 그래도 계획을 출력합니다.
- **백업.** 계획마다 계정 폴더 안에 자체 `.ccam-backup/<time>/` 폴더가 생기며, 같은 초에 두 계획이 시작되면 번호 접미사가 붙습니다. ccam은 이 폴더를 절대 삭제하지 않습니다.
- **부분 실패.** 단계가 실패하면 ccam은 오류와 원본이 있는 위치를 출력하고, `ccam doctor --fix`로 마무리하라고 안내합니다. 목표 설정은 이미 저장되어 있으므로 복구는 그 목표를 향해 이어집니다.
- **직접 만든 링크.** 어떤 경로가 이미 직접 만든 링크이고(예: dotfiles 저장소를 가리키는 링크) 그룹이 비어 있으면, ccam은 링크 대상의 내용을 그룹에 복사하고 링크 자체는 백업으로 옮깁니다. 링크 대상은 건드리지 않습니다.
- **로그인 파일.** `.credentials.json`과 `.claude.json`은 절대 그룹에 링크하거나 옮기거나 복사하지 않습니다.
- **같은 폴더에서 두 계정 쓰기.** 대화 내역을 공유하는 계정들은 같은 파일에 기록합니다. 같은 대화 내역 그룹의 다른 계정이 이미 실행 중인 폴더에서 두 번째 계정을 시작하면, ccam이 경고하고 계속할지 묻습니다.

## 관련 문서

- [MCP 서버, 설정, 대화 내역 공유](/docs/sharing)
- [MCP 서버 공유](/features/mcp)
- [기존 설정 복사하기](/features/migrate)
- [전체 명령](/docs/commands)
