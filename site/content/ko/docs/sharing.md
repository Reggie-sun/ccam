---
title: MCP 서버, 설정, 대화 내역 공유
description: ccam이 로그인은 분리한 채 MCP 서버, settings.json, CLAUDE.md, 에이전트, 스킬, 대화 내역을 계정끼리 공유하는 방법.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

계정마다 로그인은 따로 둡니다. 그 밖의 모든 것은 항목별로 공유할 수 있습니다. ccam은 계정의 파일이나 폴더를 공유 그룹 폴더에 링크하므로, 한 계정에서 바꾼 내용이 같은 그룹의 모든 계정에 나타납니다.

## 공유할 수 있는 항목

| 항목 | 옵션 | 공유되는 것 |
| --- | --- | --- |
| MCP 서버 | `--mcp` | 그룹의 `mcp.json`. 시작할 때 Claude Code에 전달 |
| 설정 | `--settings` | `settings.json`: 권한, 훅, 모델, 상태줄 |
| 전역 메모리 | `--memory` | `CLAUDE.md` |
| 확장 기능 | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| 대화 내역 | `--history` | `projects/`와 `history.jsonl` |

새 계정은 대화 내역을 뺀 모든 항목을 공유합니다.

로그인 파일인 `.credentials.json`과 `.claude.json`은 절대 공유되지 않습니다. ccam은 정해진 목록에 있는 경로만 링크하며, 이 두 파일은 그 목록에 없습니다.

## 그룹

`on`은 `default` 그룹을 뜻합니다. 그룹 이름을 직접 지정해 일부 계정끼리만 항목을 공유하게 할 수도 있습니다.

```sh
ccam share work --history team
ccam share side --history team
```

이제 `work`와 `side`는 서로 대화를 공유하지만, 다른 계정과는 공유하지 않습니다.

## 공유 켜기

계정이 어떤 항목을 공유하기 시작하면 ccam은 다음과 같이 합니다.

1. 옮기거나, 복사하거나, 링크할 파일을 모두 담은 계획을 출력하고, 시작하기 전에 확인을 받습니다.
2. 그룹이 아직 비어 있으면 계정의 파일을 그룹 폴더로 옮기고, 그렇지 않으면 그룹에 아직 없는 파일만 추가합니다.
3. 같은 파일의 내용이 다르면 그룹 것을 유지할지, 내 것을 쓸지, JSON을 병합할지 묻습니다.
4. 원본은 계정 폴더 안의 `.ccam-backup/<time>/`에 남깁니다.
5. 원본을 그룹을 가리키는 링크로 바꿉니다.

대화 파일은 세션 ID로 이름이 붙으므로, 여러 계정의 대화 내역이 충돌 없이 합쳐집니다.

## 공유 끄기

```sh
ccam share work --history off
```

링크가 제거되고, 계정은 그룹의 현재 내용을 복사한 자기 사본을 갖게 됩니다. 대신 빈 상태로 시작하려면 `--fresh`를 붙이세요. 그룹은 다른 계정을 위해 내용을 그대로 유지합니다.

## MCP 서버

Claude Code는 MCP 서버를 로그인과 함께 `.claude.json`에 저장하므로 이 파일은 공유할 수 없습니다. ccam은 공유 서버를 그룹의 `mcp.json`에 두고 `--mcp-config=<group>/mcp.json`으로 Claude Code를 시작합니다. 직접 지정한 `--mcp-config` 옵션도 그 위에 그대로 적용됩니다.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

세션 안에서 `claude mcp add`로 추가한 서버는 그 계정에만 저장됩니다. 다음에 그 계정을 시작하면 ccam이 그 서버의 이름과, 서버를 그룹으로 옮기는 명령을 알려 줍니다.

```sh
ccam mcp import work gh
```

OAuth로 로그인하는 MCP 서버는 토큰을 계정의 로그인과 함께 저장하므로, 계정마다 한 번씩 인증해야 합니다.

서버가 `claude mcp list`와 `/mcp`에 일반 사용자 서버로 나타나는 쪽이 좋다면 `~/.config/ccam/config.toml`에 `mcp_strategy = "sync"`를 설정하세요. 그러면 ccam은 시작하기 전에 각 계정의 `.claude.json`에 서버를 써 넣으며, 이때 Claude Code 자체의 파일 잠금을 쓰고 백업도 남깁니다.

## 같은 폴더의 두 계정

대화 내역을 공유하는 계정들은 같은 파일에 기록합니다. 같은 그룹의 다른 계정이 이미 실행 중인 폴더에서 Claude Code를 하나 더 시작하면, ccam이 경고하고 계속할지 묻습니다.

## 문제가 생겼을 때

```sh
ccam doctor
ccam doctor --fix
```

예를 들어 백업을 손으로 복원한 뒤 공유 경로가 일반 파일이 되었다면, ccam은 Claude Code를 시작하기 전에 멈추고 `ccam doctor --fix`를 안내합니다. 이 명령은 그 파일을 그룹에 합치고 백업을 남깁니다.
