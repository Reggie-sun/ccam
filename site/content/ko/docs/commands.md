---
title: 명령어 레퍼런스
description: ccam add와 ccs <name>부터 ccam usage, ccam next, ccam doctor까지 ccam의 모든 명령과 옵션입니다.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam`과 `ccs`는 같은 프로그램이므로 편한 쪽을 쓰면 됩니다. 어느 명령이든 앞에 `--dry-run`을 붙이면 실제로 실행하지 않고 무엇을 할지만 보여 줍니다. 예: `ccs --dry-run work`.

<!-- commands -->

## Claude Code 시작하기

`ccs <name|number> [arguments]`는 그 계정으로 Claude Code를 시작합니다. ccam은 이름 뒤의 인자를 읽지 않고 같은 순서로 `claude`에 넘기며, 표준 입력, 종료 코드, 시그널도 `claude`를 직접 실행할 때와 똑같이 동작합니다. 이름 뒤의 `--`는 붙여도 되고 생략해도 됩니다.

`ccs -p "summarise"`처럼 계정 이름 없이 인자만 주고 `ccs`를 실행하면 기본 계정을 씁니다. 기본 계정은 `ccam use <name>`으로 지정합니다.

단축 명령 `claude1`, `claude2` 등은 `~/.local/bin`에 있는, ccam을 가리키는 링크입니다. `ccam shims install --names`를 실행하면 계정마다 `claude-<name>`도 만듭니다.

## 계정 추가와 삭제

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

폴더를 지정하지 않고 `ccam import`를 실행하면 `~/.claude`와 `~/.claude-*`를 찾아 하나씩 물어봅니다. 가져온 폴더는 제자리에서 등록되며, ccam은 로그인을 옮기거나 복사하지 않습니다.

## 공유

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

각 항목에는 `on`(`default` 그룹), `off`, 또는 그룹 이름을 지정합니다. [공유](/docs/sharing)를 참고하세요.

## MCP 서버

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## 사용량과 계정 전환

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## 유지 관리

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## 파일

| 경로 | 내용 |
| --- | --- |
| `~/.config/ccam/config.toml` | 계정 목록과 각 계정의 공유 항목 |
| `~/.local/share/ccam/profiles/<name>/` | 계정별 `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | 공유하는 설정, MCP 서버, 확장 기능, 대화 내역 |
| `~/.local/state/ccam/` | 사용량 캐시와 마지막 사용 시각 |
| `~/.local/bin/claude1`, `ccs` | ccam을 가리키는 링크 |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME` 변수를 설정하면 이 폴더들의 위치가 바뀝니다.
