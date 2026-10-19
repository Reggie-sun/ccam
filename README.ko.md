# ccam — Claude Code Account Manager (Claude Code 다중 계정 관리)

[English](README.md) · **한국어** · 웹사이트: https://ccam.dev

- ccam은 컴퓨터 한 대에서 Claude Code 계정 여러 개를 쓰고 명령 하나로 바꿀 수 있게 해 주는 무료 오픈소스 명령줄 도구입니다.
- 계정마다 별도의 CLAUDE_CONFIG_DIR에 자기 로그인을 따로 두고, MCP 서버, 설정, CLAUDE.md, 확장 기능, 대화 내역은 계정끼리 공유할 수 있습니다. 그래서 한 계정이 사용량 한도에 도달하면 다른 계정에서 대화를 이어 갈 수 있습니다.
- Linux나 macOS에서 명령 하나로 설치한 뒤, 계정마다 ccam add를 실행합니다.
- ccam은 비공식 프로젝트이며 Anthropic과 제휴하거나 Anthropic의 보증을 받은 것이 아닙니다.

`ccam`은 우분투 터미널에서 Claude Code 계정 여러 개를 추가하고 바꿔 쓰게 해 주는 단일 바이너리 도구입니다. 계정마다 `CLAUDE_CONFIG_DIR`을 따로 두어 로그인을 분리하고, MCP 서버·설정·전역 CLAUDE.md·확장 기능·대화 내역은 항목별로 공유할지 고릅니다.

> **이용약관을 먼저 확인하세요.** 계정을 여러 개 쓰는 것이 허용되는지는 각 계정의 Anthropic 이용약관과 소속 조직의 정책에 따릅니다. ccam은 계정 전환을 편하게 할 뿐, 사용량 한도를 우회하지 않습니다.

프로그램 출력(메시지, 질문, 도움말)은 모두 영어입니다.

## 설치

```bash
curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

리눅스(우분투·데비안·페도라·아치·알파인 등)와 macOS를 지원합니다. `~/.local/bin/ccam`이 생깁니다. 체크섬을 확인한 뒤 설치하고, `~/.local/bin`이 `PATH`에 없으면 알려 줍니다. 그 밖의 방법:

- `.deb`: [릴리스](https://github.com/jhyunwoo/ccam/releases)에서 받아 `sudo apt install ./ccam_*.deb`
- 소스: `go install github.com/jhyunwoo/ccam/cmd/ccam@latest` (Go 1.26 이상)

Claude Code(`claude`)가 먼저 설치되어 있어야 합니다.

## 업데이트

```bash
ccam update            # 최신 릴리스 설치
ccam update --check    # 새 릴리스가 있는지 알려 주기만 함
```

`ccam update`는 내 시스템용 릴리스를 받아 릴리스의 `checksums.txt`(SHA-256)로 확인한 뒤, `ccam` 바이너리를 파일 이름 바꾸기 한 번으로 교체합니다. 단축 명령은 이 파일을 가리키는 링크라 함께 새 버전이 됩니다. `--version v0.2.0`은 그 릴리스를 설치하고(이전 버전도 가능), `ccam --dry-run update`는 받을 파일과 바꿀 파일만 보여줍니다. ccam이 스스로 업데이트를 확인하는 일은 없습니다. `.deb`로 설치했다면 apt에 맡깁니다. [릴리스](https://github.com/jhyunwoo/ccam/releases)에서 새 `.deb`를 받아 `sudo apt install ./ccam_<버전>_linux_<아키텍처>.deb`를 실행하세요. 설치 스크립트를 다시 실행해도 `~/.local/bin`의 ccam이 업데이트됩니다.

## 빠른 시작

```bash
ccam add             # 이름 → 공유 항목 → 단축 명령 만들기? (y/N, 이름은 직접) → 브라우저 로그인
ccam add work --shortcut ccwork    # 두 번째 계정, 단축 명령 ccwork

ccwork               # 단축 명령으로 work 계정의 Claude Code
ccam work --resume   # 이름으로 열기, 뒤의 인자는 그대로 claude에 전달
ccam                 # 선택 화면 (Enter 실행, 1-9 바로 실행, s 공유, a 추가)
ccam usage           # 모든 계정의 5시간·주간 사용률
ccam next --continue # 한도에 걸리면 다음 계정에서 같은 대화 이어가기
```

이미 쓰던 `~/.claude`가 있으면 `ccam import`로 제자리에 등록합니다 (옮기거나 복사하지 않음).

쓰던 설정을 공유 설정의 출발점으로 삼으려면 `ccam migrate`를 실행합니다. `~/.claude`의 settings.json, CLAUDE.md, 확장 기능, MCP 서버를 공유 그룹으로 복사합니다. 원래 폴더는 읽기만 하고, 대화 내역은 `--history`를 줄 때만 복사합니다. 터미널에서는 먼저 복사할 항목을 체크 목록에서 고르고, 그다음 복사할 내용을 보여준 뒤 확인을 받습니다 (`ccam --dry-run migrate`는 보여주기만 함). 복사한 MCP 서버의 저장된 로그인(OAuth)은 그룹이 아니라 각 계정의 로그인 파일에 넣습니다.

SSH로 접속한 원격 서버에서 로그인할 때는 내 컴퓨터의 브라우저가 열리지 않는 `http://localhost:<포트>/callback?code=…` 페이지에서 멈춥니다. 그 주소를 `ccam add`나 `ccam login`이 실행 중인 터미널에 붙여 넣으면 ccam이 서버의 Claude Code에 전달해 로그인을 마칩니다.

## 명령

| 명령 | 하는 일 |
| --- | --- |
| `ccam` | 선택 화면 (터미널이 아니면 목록 출력) |
| `ccam <이름\|번호> [claude 인자]` | 그 계정으로 `claude` 실행. 인자는 해석하지 않고 그대로 전달 |
| `<단축 명령> [claude 인자]` | 계정에 단축 명령이 있으면 그 이름으로 바로 실행 |
| `ccam add [이름]` | 계정 추가. 터미널에서는 마법사, `--no-login`·`--share 항목=그룹`·`--shortcut <명령>`도 가능 |
| `ccam shortcut <이름> [명령] [--remove]` | 계정의 단축 명령 보기·만들기·이름 바꾸기·지우기 |
| `ccam list [--json]` | 계정 목록 |
| `ccam use <이름>` | 기본 계정 (`ccam -p "…"`처럼 이름 없이 실행할 때) |
| `ccam next [--auto] [--continue]` | 다음 계정으로. `--auto`는 5시간 사용률이 가장 낮은 계정 |
| `ccam usage [--watch] [--json] [-v]` | 사용량 (줄임 `ccam u`) |
| `ccam share <이름> --history on …` | 공유 변경. 이전 계획을 보여주고 확인 후 실행 |
| `ccam mcp add\|add-json\|list\|remove\|import` | 공유 MCP 서버 |
| `ccam login\|logout <이름>` | 다시 로그인 / 토큰 삭제 |
| `ccam rename` / `ccam remove` | 이름 변경 / 삭제 (폴더는 휴지통으로) |
| `ccam import [폴더]` | 기존 폴더 등록 |
| `ccam migrate [폴더]` | 기존 설정(기본 `~/.claude`)을 공유 그룹으로 복사. 원본은 바꾸지 않음. `--group`, 항목 옵션, `--history`, `--yes` |
| `ccam doctor [--fix]` | 링크·권한·로그인·버전·사용량 형식 점검과 수리 |
| `ccam statusline install` | Claude Code 상태줄에 `[account] · model · 5h N%` 표시 |
| `ccam update [--check] [--version vX.Y.Z]` | 최신(또는 지정한) 릴리스로 업데이트, 체크섬 확인 |
| `ccam completion bash\|zsh` | 자동완성 |

앞에 `--dry-run`을 붙이면 실행하지 않고 할 일만 보여줍니다 (`ccam --dry-run work`). 실행(`ccam <이름>`, `run`, `next`), `share`, `migrate`, `shortcut`, `update`, `doctor --fix`에서 쓸 수 있고, 무언가를 바꾸는 다른 명령은 실제로 실행하는 대신 거절합니다.

## 공유

항목마다 `local`(분리) 또는 그룹 이름(기본 `default`)을 고릅니다. 같은 그룹의 계정끼리만 내용을 함께 씁니다.

| 항목 | 방식 |
| --- | --- |
| MCP 서버 | 그룹의 `mcp.json`을 실행할 때 `--mcp-config`로 전달 (`.claude.json`은 건드리지 않음) |
| settings.json, 전역 CLAUDE.md | 그룹 파일로 심볼릭 링크 |
| 에이전트·명령·스킬·플러그인 | `agents/ commands/ skills/ plugins/ output-styles/ rules/` 폴더 링크 |
| 대화 내역 | `projects/`와 `history.jsonl` 링크 → 다른 계정에서 `--continue`/`--resume` |

로그인 파일(`.credentials.json`, `.claude.json`)은 어떤 경우에도 공유되지 않습니다. 공유를 켜고 끌 때는 원본을 `.ccam-backup/<시각>/`에 남기고, 파일이 서로 다르면 "keep group / use mine / merge JSON" 중에서 고르게 합니다.

세션 안에서 `claude mcp add`로 넣은 서버는 그 계정에만 저장됩니다. 다음 실행 때 ccam이 알려 주면 `ccam mcp import <계정> <서버>`로 그룹에 옮기세요. `.claude.json`을 직접 고치는 방식이 필요하면 `config.toml`에 `mcp_strategy = "sync"`를 둡니다.

대화 내역을 공유하는 두 계정이 같은 폴더에서 동시에 돌면 경고합니다.

## 사용량

```
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

각 계정의 토큰으로 사용량 API를 동시에 조회합니다 (계정당 3초 제한, 60초 캐시). 이 API는 공개 문서가 없어 바뀔 수 있습니다. 실패하면 마지막 값이나 상태줄이 모아 둔 값을 보여주고, `ccam doctor`가 형식 변경을 알려 줍니다. 만료된 토큰은 직접 갱신하지 않습니다 (`--refresh`는 그 계정이 실행 중이 아닐 때만).

## 파일 위치

```
~/.config/ccam/config.toml           계정 목록과 공유 설정
~/.local/share/ccam/profiles/<이름>/  계정별 CLAUDE_CONFIG_DIR
~/.local/share/ccam/groups/<그룹>/    공유 저장소
~/.local/state/ccam/                 사용량 캐시, 마지막 사용 시각
$XDG_RUNTIME_DIR/ccam/               잠금 파일
~/.local/bin/<단축 명령>              ccam을 가리키는 링크 (만든 계정만)
```

`CCAM_HOME`을 지정하면 모든 경로가 그 아래로 바뀝니다 (테스트용).

## 성능

`ccam <이름>`은 `config.toml` 하나만 읽고 `lstat`으로 링크를 확인한 뒤 `exec`로 `claude`가 됩니다. 일반 노트북에서 ccam이 더하는 시간은 p95 약 4 ms, 메모리는 8 MB입니다. `CCAM_TRACE=1`이면 단계별 시간을, `CCAM_PPROF=<파일>`이면 CPU 프로파일을 남깁니다.

## 문제 해결

- 무엇이든 이상하면 `ccam doctor`, 고칠 수 있는 것은 `ccam doctor --fix` (바뀌는 파일은 백업).
- "공유 링크에 문제가 있어 실행을 멈췄습니다" → 공유 항목 자리에 실제 파일이 생긴 경우입니다. `doctor --fix`가 그룹으로 합치고 원본을 백업합니다.
- 로그인 만료 → `ccam login <이름>`.

## 개발

```bash
go build ./cmd/ccam
go test ./...
golangci-lint run
scripts/perf-check.sh            # hyperfine으로 실행 경로 목표 확인
scripts/bench-compare.sh main    # benchstat으로 20% 넘는 회귀 검사
```

설계와 결정 기록은 [docs/PLAN.md](docs/PLAN.md), Claude Code 폴더 조사는 [docs/claude-layout.md](docs/claude-layout.md), 수동 점검표는 [docs/manual-qa.md](docs/manual-qa.md)에 있습니다.

## 라이선스

[MIT](LICENSE). ccam은 비공식 프로젝트이며 Anthropic과 제휴하거나 Anthropic의 보증을 받은 것이 아닙니다.
