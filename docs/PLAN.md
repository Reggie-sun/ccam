# ccam 개발 계획 — Claude Code 다중 계정 관리 도구

2026-10-18 작성

## 개요와 목표

`ccam`(Claude Code Account Manager, 짧은 별칭 `ccs`)은 우분투 터미널에서 Claude Code 계정 여러 개를 추가하고 전환하는 Go 단일 바이너리 CLI/TUI 도구입니다. 계정마다 `CLAUDE_CONFIG_DIR`을 따로 두어 로그인 정보를 분리하고, MCP·설정·대화 내역은 계정별로 공유할지 고를 수 있게 합니다.

핵심 요구사항은 세 가지입니다.

1. **쉬운 추가**: `ccam add` 한 번으로 프로필 생성, 로그인, 단축 명령(`claude1`…) 생성까지 끝납니다.
2. **쉬운 전환**: `claude2`, `ccs work`, 또는 `ccam` 한 단어로 뜨는 선택 화면에서 계정을 고릅니다. `ccam next --continue`로 한도에 걸린 대화를 다음 계정에서 이어갑니다.
3. **공유 선택**: MCP 서버, 설정, 전역 CLAUDE.md, 에이전트·명령·스킬, 대화 내역을 항목별·계정별로 공유 또는 분리합니다.

범위 밖: Claude Code 자체 수정, 사용량 한도 우회 자동화, macOS/Windows 지원(v1은 우분투 22.04/24.04 기준).

## 전제: Claude Code가 저장하는 것

`CLAUDE_CONFIG_DIR`을 지정하면 기본 `~/.claude` 대신 그 폴더에 설정과 로그인 정보가 저장됩니다([Maestro 문서](https://docs.runmaestro.ai/multi-claude)). 이 폴더 안의 항목은 "분리 필수"와 "공유 가능"으로 나뉘고, ccam 설계는 이 구분을 따릅니다.

| 경로 (프로필 폴더 기준) | 내용 | ccam 처리 |
| --- | --- | --- |
| `.credentials.json` | OAuth 토큰(리눅스) | 항상 분리, 권한 0600 |
| `.claude.json` | 계정 정보(`oauthAccount`), 사용자 범위 MCP(`mcpServers`), 프로젝트별 상태 | 파일은 분리, MCP 키만 선택 공유 |
| `settings.json` | 권한, 훅, 모델, env | 선택 공유 (심볼릭 링크) |
| `CLAUDE.md` | 전역 메모리 | 선택 공유 (심볼릭 링크) |
| `agents/` `commands/` `skills/` `plugins/` | 확장 기능 | 선택 공유 (심볼릭 링크) |
| `projects/` | 프로젝트별 대화 기록(`.jsonl`) | 선택 공유 (심볼릭 링크) |
| `history.jsonl` | 입력 프롬프트 기록 | 대화 내역과 함께 선택 공유 |
| `todos/` `shell-snapshots/` `statsig/` 등 | 실행 중 상태·캐시 | 항상 분리 |

핵심 제약: MCP 서버 목록은 계정 정보와 같은 `.claude.json`에 들어 있어서 파일째 공유할 수 없습니다. 그래서 MCP만 별도 방식(공유 설계 섹션)이 필요합니다.

버전마다 파일 구성이 바뀔 수 있으므로 위 표는 Phase 0에서 설치된 Claude Code로 직접 확인하고, 확인 결과를 `docs/claude-layout.md`에 기록합니다. 과거에 `CLAUDE_CONFIG_DIR` 사용 시 MCP가 로드되지 않던 버그가 보고된 적도 있습니다([이슈](https://claudeissues.com/issue/42217-mcp-servers-from-mcp-json-not-loaded-when-claude-config-dir-is-set), 2026년 6월 해결로 표시).

## 사용자 경험과 명령어

사용자는 처음에 `ccam add`만 알면 되고, 이후에는 `claude1`·`claude2` 같은 단축 명령이나 `ccam` 선택 화면으로 전환합니다. 단축 명령은 셸 별칭이 아니라 `~/.local/bin`에 ccam 바이너리를 가리키는 심볼릭 링크로 만들어서 bash·zsh·스크립트 어디서나 동작하게 합니다.

`ccs`는 `ccam`과 똑같이 동작하는 짧은 별칭입니다. `ccs add`, `ccs next`처럼 어느 명령에서나 바꿔 쓸 수 있고, 설치할 때 `~/.local/bin/ccs` 심볼릭 링크로 함께 만들어집니다.

가장 자주 쓰는 형태는 `ccs <이름>`입니다. 하위 명령 없이 이름이나 번호만 주면 그 계정으로 Claude Code가 바로 열리고, 뒤에 붙인 인자는 그대로 `claude`에 전달됩니다(예: `ccs work --continue`). 이름과 하위 명령이 헷갈리지 않도록 `add`, `list`, `run` 같은 하위 명령 이름은 프로필 이름으로 쓸 수 없게 `ccs add`에서 막습니다.

**인자 전달 규칙**: 프로필 이름(또는 번호) 뒤의 모든 인자는 ccam이 해석하지 않고 순서 그대로 `claude`에 넘깁니다. 그래서 `claude`로 하던 일은 앞부분만 바꾸면 똑같이 됩니다.

| 입력 | 실제 실행 |
| --- | --- |
| `ccs work --resume` | work 계정으로 `claude --resume` |
| `ccs work -p "테스트 고쳐줘" --output-format json` | work 계정으로 `claude -p "테스트 고쳐줘" --output-format json` |
| `claude2 --model opus --continue` | 2번 계정으로 `claude --model opus --continue` |
| `ccs work --help` | `claude --help` (ccam 도움말이 아님) |
| `ccs -p "요약해줘"` | 기본 계정으로 `claude -p "요약해줘"` |
| `ccs work -- --version` | `claude --version` (`--`는 생략 가능, 넣어도 동작) |
| `cat log.txt \| ccs work -p "원인 찾아줘"` | 표준 입력도 그대로 전달 |

세부 규칙:

- 첫 인자가 `-`로 시작하면 ccam 자체 옵션(`--help`, `--version`), 실행할 명령만 출력하고 멈추는 `--dry-run`인지만 확인하고, 아니면 기본 계정으로 `claude`에 넘깁니다.
- 인자 배열은 셸 문자열로 합치지 않고 배열 그대로 `exec`합니다. 따옴표, 공백, 한글 인자가 깨지지 않고, 종료 코드와 시그널도 `claude`의 것이 그대로 남습니다.
- MCP flag 전략이 붙이는 `--mcp-config`는 사용자 인자 앞에 넣어서, 사용자가 직접 준 `--mcp-config`·`--strict-mcp-config`가 함께 적용되게 합니다. Phase 0에서 여러 개를 줬을 때의 동작을 확인합니다.
- `ccs run <이름> [인자]`도 같은 규칙을 따릅니다.

| 명령 | 동작 |
| --- | --- |
| `ccam` | TUI 선택 화면: 계정 목록, 이메일, 공유 상태, 마지막 사용 시각. Enter로 실행 |
| `ccam <이름\|번호> [claude 인자]` | 그 계정으로 Claude Code를 바로 엽니다 (예: `ccs work`, `ccs 2`). `run`의 줄임 |
| `ccam add [이름]` | 대화형 마법사: 이름 → 공유 항목 선택 → Claude Code 실행해 로그인 → 이메일 확인 → 단축 명령 생성 |
| `ccam list` | 표 출력 (`--json` 지원) |
| `ccam run <이름\|번호> [-- claude 인자]` | 해당 계정으로 `claude` 실행 (`exec`로 교체) |
| `ccam use <이름>` | 기본 계정 지정. `ccam run`에 이름이 없으면 이 계정 사용 |
| `ccam next [--auto] [--continue]` | 다음 계정으로 실행. `--continue`면 대화 내역 공유 시 마지막 대화를 이어감 |
| `ccam usage [--watch] [--json]` | 모든 계정의 5시간·주간 사용률과 초기화 시각을 한 번에 조회 (줄임: `ccs u`) |
| `ccam share <이름> --mcp on --history off …` | 공유 항목 변경 (데이터 이전 포함) |
| `ccam mcp add\|list\|remove` | 공유 MCP 저장소 관리 |
| `ccam login <이름>` / `ccam logout <이름>` | 재로그인 / 토큰 삭제 |
| `ccam rename` / `ccam remove` | 이름 변경 (단축 명령 갱신) / 삭제 (확인 후 휴지통 폴더로 이동) |
| `ccam import` | 기존 `~/.claude`나 `~/.claude-N` 폴더를 프로필로 가져오기 |
| `ccam doctor` | 링크 깨짐, 권한, 로그인 만료, Claude Code 버전 점검 및 자동 수리 |
| `ccam shims install` | `claude1`…`claudeN` 재생성. 번호는 추가 순서, 이름 단축 명령(`claude-work`)도 선택 |
| `ccam completion bash\|zsh` | 자동완성 스크립트 출력 |

`ccam add` 화면 예시:

```
$ ccam add work
? 공유할 항목을 고르세요 (스페이스로 선택)
  [x] MCP 서버
  [x] settings.json
  [x] 전역 CLAUDE.md
  [x] 에이전트·명령·스킬
  [ ] 대화 내역
→ Claude Code를 엽니다. /login 으로 로그인한 뒤 /exit 하세요.
✓ work (dev@company.com) 추가됨 — 단축 명령: claude3
```

실행할 때는 Claude Code 화면 위에 계정 이름이 보이도록 터미널 제목(`\e]0;claude · work\a`)을 바꾸고, 선택적으로 Claude Code 상태줄(statusline)에 계정명을 표시하는 스크립트를 함께 설치합니다.

## 사용량 조회

`ccs usage`(줄여서 `ccs u`)는 모든 프로필의 5시간·주간 사용률과 초기화 시각을 한 화면에 보여줍니다. 각 계정에 동시에 요청하고, 응답이 늦거나 실패한 계정은 마지막으로 저장해 둔 값을 기준 시각과 함께 표시합니다.

```
$ ccs usage
 계정            5시간               초기화     주간                초기화       기준
 1 personal     ███████░░░  72%    1시간 40분  ████░░░░░░  41%    3일 6시간    방금
 2 work         █░░░░░░░░░   8%    4시간 12분  ██░░░░░░░░  19%    5일 2시간    방금
 3 side         ██████████ 100%    22분        ███████░░░  68%    1일 9시간    12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccs next --auto 로 바로 열기)
```

데이터 출처는 두 가지이고, 실시간 조회가 기본입니다.

| 출처 | 방식 | 장점 | 한계 |
| --- | --- | --- | --- |
| 실시간 조회 | 각 프로필의 `.credentials.json` 토큰으로 `GET https://api.anthropic.com/api/oauth/usage` 호출 ([ccusage](https://pypi.org/project/ccusage/)가 쓰는 방식) | Claude Code를 켜지 않아도 지금 값을 받음 | 공개 문서가 없는 엔드포인트라 예고 없이 바뀔 수 있음 |
| 상태줄 캐시 | ccam이 설치하는 상태줄 스크립트가 Claude Code가 넘겨주는 `rate_limits.five_hour`·`seven_day` 값을 프로필별 파일에 저장 ([설명](https://wmedia.es/en/tips/claude-code-usage-limit-status-line)) | Claude Code가 직접 주는 값이라 안정적 | 그 계정을 마지막으로 쓴 시점의 값, 세션의 첫 응답 이후에만 생김 |

표시 규칙:

- 실시간 응답의 새 형식(`limits` 배열의 `session`·`weekly_all`·`weekly_scoped`)과 예전 형식(`five_hour`·`seven_day`)을 모두 읽습니다. 모델별 주간 한도가 있으면 `ccs usage -v`에서 따로 보여줍니다.
- 토큰이 만료됐으면 ccam이 직접 갱신하지 않고 "로그인 만료, `ccs <이름>`을 한 번 실행하면 갱신"이라고 표시합니다. 실행 중인 Claude Code와 동시에 토큰 파일을 고쳐 쓰는 일을 막기 위해서입니다. 갱신을 원하면 `--refresh`를 붙입니다(파일 잠금 후 갱신).
- API 키로 로그인한 계정처럼 구독 한도가 없는 계정은 "한도 없음(API)"으로 표시합니다.
- 같은 계정을 60초 안에 다시 조회하면 저장된 값을 씁니다. `--no-cache`로 무시할 수 있습니다.
- 80% 이상은 노란색, 100%는 빨간색으로 표시합니다. 색만으로 구분하지 않도록 퍼센트 숫자도 항상 함께 씁니다.

옵션: `--json`(스크립트용), `--watch [초]`(주기적으로 새로 고침, 기본 60초), `-v`(모델별 한도와 출처 표시), `--source live|cache`.

다른 기능과의 연결:

- `ccs` 선택 화면의 계정 목록에 5시간 사용률을 함께 표시합니다.
- `ccs next --auto [--continue]`는 5시간 사용률이 가장 낮은 계정을 골라 엽니다. 100%인 계정은 건너뜁니다.

## 공유 설계

공유는 "공유 그룹" 단위로 동작합니다. 항목마다 프로필 값이 `local`(분리)이거나 그룹 이름(기본 `default`)이고, 같은 그룹의 프로필끼리만 내용을 함께 씁니다. 예를 들어 1~3번은 대화 내역을 `team` 그룹으로 공유하고 4~5번은 분리할 수 있습니다.

```
토큰은 계정마다 따로, 나머지는 그룹 폴더로 링크

  claude1 · personal ──┐
  claude2 · work ──────┼──▶  default 그룹
  claude3 · side ──────┘      settings.json · CLAUDE.md
                              mcp.json (MCP 서버)
                              agents · commands · skills

  claude2 · work ──────┬──▶  team 그룹
  claude3 · side ──────┘      projects/ (대화 기록)
                              history.jsonl

  각 계정의 .credentials.json · .claude.json 은 항상 분리
```

세 계정 모두 설정·MCP·확장 기능은 default 그룹을 쓰고, 대화 내역은 2·3번만 team 그룹으로 공유합니다. 로그인 파일은 어떤 그룹에도 연결되지 않습니다.

| 항목 | 공유 방식 | 비고 |
| --- | --- | --- |
| MCP 서버 | 기본: 실행 시 `--mcp-config <그룹>/mcp.json` 전달. 대안: 실행 전 `.claude.json`의 `mcpServers` 키 동기화 | 둘 중 무엇을 기본으로 할지 Phase 0에서 결정 |
| settings.json | 그룹 폴더 파일로 심볼릭 링크 | |
| 전역 CLAUDE.md | 심볼릭 링크 | |
| 에이전트·명령·스킬·플러그인 | 폴더 단위 심볼릭 링크 | 항목 하나로 묶어 토글 |
| 대화 내역 | `projects/`와 `history.jsonl`을 심볼릭 링크 | 공유하면 다른 계정에서 `--continue`/`--resume` 가능 |

MCP 두 방식의 차이:

- **실행 플래그 방식**은 `.claude.json`을 전혀 건드리지 않아 안전합니다. 대신 세션 안에서 `claude mcp add`로 추가한 서버는 그 계정에만 저장되므로, ccam이 다음 실행 때 새 서버를 감지해 "공유 저장소로 옮길까요?"를 묻습니다.
- **동기화 방식**은 `/mcp` 화면에 서버가 원래처럼 보이지만, Claude Code가 수시로 다시 쓰는 파일을 수정하므로 파일 잠금, 백업, 같은 프로필이 실행 중인지 확인이 필요합니다.

공유를 켜고 끌 때의 데이터 이전 규칙:

1. **켜기 (local → 그룹)**: 프로필의 기존 내용을 그룹 폴더로 복사합니다. 대화 파일은 세션 ID가 파일명이라 그냥 합치고, `settings.json`·`CLAUDE.md`가 다르면 차이를 보여주고 "그룹 것 유지 / 내 것으로 덮기 / JSON 병합" 중 고르게 합니다. 그다음 원본을 `.ccam-backup/<시각>/`로 옮기고 링크를 겁니다.
2. **끄기 (그룹 → local)**: 링크를 지우고 그룹 내용을 복사해 둡니다. `--fresh`를 주면 빈 상태로 시작합니다.
3. 모든 이전 작업은 먼저 계획을 출력하고(`--dry-run`과 같음) 확인을 받은 뒤 실행합니다.

동시 실행 주의: 대화 내역을 공유하는 두 계정이 같은 프로젝트 폴더에서 동시에 돌면 같은 세션 파일을 놓고 충돌할 수 있습니다. ccam은 `(그룹, 작업 폴더)`마다 잠금 파일을 두고, 이미 실행 중이면 경고 후 계속할지 묻습니다.

## 아키텍처

Go로 작성한 단일 바이너리입니다. 단축 명령이 매번 ccam을 거쳐 `claude`를 띄우므로 시작 속도가 중요하고, 의존성 없는 바이너리 하나라 설치도 쉽습니다.

| 영역 | 선택 |
| --- | --- |
| 언어 | Go 1.23 이상 |
| CLI | `spf13/cobra` |
| TUI | `charmbracelet/bubbletea`, `bubbles`, `lipgloss`, 입력 폼은 `charmbracelet/huh` |
| 설정 파일 | TOML (`pelletier/go-toml/v2`) |
| 파일 잠금 | `gofrs/flock` |
| 배포 | `goreleaser`로 `.deb`와 `tar.gz`, `install.sh` |

파일 위치 (XDG 규칙):

```
~/.config/ccam/config.toml            # 프로필 목록과 공유 설정
~/.local/share/ccam/profiles/<이름>/   # 각 계정의 CLAUDE_CONFIG_DIR
~/.local/share/ccam/groups/<그룹>/     # 공유 저장소
    settings.json  CLAUDE.md  mcp.json
    agents/ commands/ skills/ plugins/
    projects/  history.jsonl
~/.local/state/ccam/locks/  logs/  usage/   # 잠금, 로그, 사용량 캐시
~/.local/bin/claude1 … claudeN       # 단축 명령 (ccam을 가리키는 링크)
~/.local/bin/ccs -> ccam             # 짧은 별칭
```

`config.toml` 예시:

```toml
version = 1
default = "personal"
claude_bin = "claude"
mcp_strategy = "flag"        # flag | sync

[[profile]]
name = "personal"
slot = 1                      # claude1
email = "me@example.com"     # .claude.json에서 읽어 캐시
# dir = "/home/me/.claude"   # import한 기존 폴더는 제자리 등록
[profile.share]
mcp = "default"
settings = "default"
memory = "default"
extensions = "default"
history = "local"
```

모듈 구성 (`internal/` 아래):

- `config`: TOML 읽기·쓰기, 스키마 버전, 원자적 저장
- `profile`: 생성·이름 변경·삭제·가져오기, `.claude.json`에서 이메일 읽기
- `share`: 심볼릭 링크 관리, 이전 계획 생성(plan)과 실행(apply) 분리
- `mcp`: 그룹 `mcp.json` 관리, flag/sync 두 전략을 같은 인터페이스로
- `launcher`: 환경변수 구성, 잠금, 터미널 제목, `syscall.Exec`로 `claude` 실행
- `shim`: `~/.local/bin` 심볼릭 링크(claude1…, ccs) 생성·정리
- `tui`: 선택 화면과 추가 마법사
- `doctor`: 점검과 자동 수리
- `usage`: 계정별 사용량 병렬 조회, 60초 캐시, 상태줄 수집기
- `fsutil`: 원자적 쓰기, 백업, 경로 유틸

`ccam run` 실행 순서:

1. 프로필 결정: 첫 인자가 하위 명령이 아니면 프로필 이름이나 슬롯 번호로 해석 (`ccs work`, `ccs 2`), `run`에 이름이 없으면 기본 계정, 인자가 아예 없으면 TUI
2. 공유 링크 상태 확인, 깨졌으면 수리하거나 중단
3. MCP 전략 적용 (flag면 인자 추가, sync면 키 병합)
4. 대화 내역 공유 시 `(그룹, 작업 폴더)` 잠금 확인
5. `CLAUDE_CONFIG_DIR` 설정, 터미널 제목 변경, `last_used` 기록
6. `claude`로 `exec` (항상 프로세스 교체, 끝난 뒤 할 일은 다음 실행 때 처리)

단축 명령은 셸 스크립트 없이 ccam 바이너리를 가리키는 심볼릭 링크이고, 실행된 이름으로 계정을 고릅니다(성능 설계 2번).

## 성능 설계

가장 자주 실행되는 `ccs <이름>`과 `claude2`는 ccam이 붙이는 지연을 10 ms 이하로 유지합니다. 사용자는 이 시간을 매번 기다리므로, 이 경로를 가장 가볍게 만드는 것이 성능 목표의 중심입니다.

| 작업 | 목표 (p95, 일반 노트북) | 측정 방법 |
| --- | --- | --- |
| `ccs <이름>` → `claude` 실행 직전까지 | 10 ms 이하 | `hyperfine 'ccs --dry-run work'` |
| `ccs list` | 30 ms 이하 | hyperfine |
| `ccs usage` 실시간 조회, 계정 5개 | 1.5 s 이하 (계정당 제한 3 s) | 통합 테스트에서 지연을 넣은 가짜 서버로 측정 |
| `ccs usage` 캐시만 사용 | 30 ms 이하 | hyperfine |
| `ccs` 선택 화면 첫 표시 | 100 ms 이하 | teatest 타이머 |
| 메모리 (실행 직전) | 15 MB 이하 | `/usr/bin/time -v` |

**실행 경로 최적화**

1. **cobra를 거치지 않는 빠른 경로**: `main()`이 먼저 `os.Args`를 보고 `ccs <이름>`, `ccs run`, 단축 명령이면 바로 실행 경로로 갑니다. 이 경로에서는 cobra, bubbletea, lipgloss, HTTP 패키지를 초기화하지 않습니다.
2. **단축 명령은 심볼릭 링크**: `claude2`는 셸 스크립트가 아니라 ccam 바이너리를 가리키는 링크이고, 실행된 이름(`argv[0]`)으로 계정을 고릅니다. 셸을 한 번 더 띄우는 비용(약 1~2 ms)이 없어집니다.
3. **`exec`로 교체**: 항상 `syscall.Exec`로 `claude` 프로세스로 바뀝니다. ccam이 메모리에 남지 않고, 키 입력·시그널·종료 코드에 끼어들지 않습니다. MCP 새 서버 감지처럼 "끝난 뒤" 할 일은 다음 ccam 실행 때 `.claude.json` 수정 시각을 비교해서 처리합니다.
4. **최소한의 파일 접근**: 실행 경로에서 읽는 파일은 `config.toml` 하나입니다. 공유 링크는 `lstat`만으로 확인하고(항목당 시스템 호출 1번), 폴더 전체를 훑는 점검은 `doctor`에서만 합니다.
5. **`claude` 경로 캐시**: `claude` 실행 파일의 절대 경로와 수정 시각을 상태 파일에 저장해 매번 `PATH`를 뒤지지 않습니다. 파일이 바뀌거나 없어지면 다시 찾습니다.
6. **큰 파일은 통째로 읽지 않음**: `.claude.json`은 프로젝트 기록이 쌓여 수 MB가 될 수 있습니다. 이메일 같은 값은 `config.toml`에 캐시하고, 꼭 읽어야 할 때는 스트리밍 JSON 디코더로 필요한 키만 꺼냅니다.
7. **쓰기 최소화**: `last_used` 기록은 작은 상태 파일 하나에 `fsync` 없이 씁니다. 설정 변경처럼 잃으면 안 되는 쓰기만 임시 파일 + `fsync` + `rename`을 거칩니다.
8. **잠금 파일은 메모리 파일시스템에**: 잠금은 `$XDG_RUNTIME_DIR`(`/run/user/<uid>`, tmpfs)에 두고 기다리지 않는 `flock`을 씁니다.

**그 밖의 최적화**

- **사용량 조회**: 계정마다 goroutine을 하나씩 띄우고 하나의 HTTP 클라이언트(keep-alive, HTTP/2)를 공유합니다. 접속 2 s, 전체 3 s 제한을 두고, 도착한 계정부터 화면에 표시합니다.
- **데이터 이전**: 같은 파일시스템이면 복사 대신 `rename`으로 옮깁니다. 복사가 필요하면 Go의 `io.Copy`가 리눅스에서 쓰는 `copy_file_range`를 그대로 활용합니다. 폴더 순회는 `fs.WalkDir`로 해서 항목마다 `stat`을 부르지 않습니다.
- **TUI**: 첫 화면은 `config.toml`만으로 그리고, 사용량은 뒤이어 비동기로 채웁니다. 스타일 객체는 한 번만 만들어 재사용합니다.
- **빌드**: `CGO_ENABLED=0`, `-trimpath -ldflags "-s -w"`로 정적·경량 바이너리를 만듭니다(목표 10 MB 이하). Phase 6에서 실행 경로 프로파일로 PGO(`default.pgo`)를 적용합니다.

**측정과 회귀 방지**

- `CCAM_TRACE=1`이면 실행 경로 단계별 소요 시간을 표준 에러에 출력합니다.
- `CCAM_PPROF=<파일>`이면 CPU 프로파일을 남깁니다.
- 실행 경로, 설정 파싱, 링크 확인, 사용량 파서에 `go test -bench` 벤치마크를 둡니다.
- CI에서 `benchstat`으로 이전 커밋과 비교해 20% 넘게 느려지거나 위 목표를 넘으면 실패시킵니다.

## 단계별 개발 로드맵

여덟 단계로 나누고, Phase 1이 끝나면 `claude1`·`claude2`가 실제로 동작하는 상태가 됩니다. 각 단계는 Claude Code에 프롬프트 하나로 맡길 수 있는 크기이고, 완료 기준을 통과해야 다음 단계로 넘어갑니다.

**Phase 0 — 조사 스파이크 (코드 없음)**

- 확인: 사용량 엔드포인트와 상태줄 `rate_limits` 응답 형식, 프로필 폴더의 실제 파일 구성, `--mcp-config` 동작, 심볼릭 링크한 `projects/`로 다른 계정에서 `--continue`가 되는지, 로그인을 CLI 하위 명령으로 할 수 있는지(`claude --help`)
- 완료 기준: `docs/claude-layout.md` 작성, `mcp_strategy` 기본값 결정
- 프롬프트: "임시 폴더 두 개를 CLAUDE_CONFIG_DIR로 써서 Claude Code의 설정 폴더 구성을 조사해줘. docs/PLAN.md의 Phase 0 확인 항목을 하나씩 실험하고 결과를 docs/claude-layout.md에 표로 정리해줘. 내 실제 ~/.claude는 건드리지 마."

**Phase 1 — 뼈대와 프로필 (MVP)**

- 구현: `go mod init`, cobra 골격, `config`, `profile`, `launcher`, `shim` / 명령 `add`(대화형 아님), `list`, `run`, `use`, `remove`, `shims install`
- 완료 기준: 두 계정으로 각각 로그인 후 `claude1`·`claude2`가 서로 다른 이메일로 뜸, `ccs <이름>`으로도 같은 계정이 열림, 실행 경로 10 ms 이하(hyperfine), 단위 테스트 통과
- 프롬프트: "docs/PLAN.md의 Phase 1을 구현해줘. 먼저 plan 모드로 파일 구조와 인터페이스를 제안하고, 승인하면 테스트부터 작성해. 모든 파일 경로는 CCAM_HOME 환경변수로 바꿀 수 있게 해서 테스트가 실제 홈 폴더를 건드리지 않게 해줘. 실행 경로는 PLAN.md의 성능 설계대로 cobra를 거치지 않게 만들고 벤치마크도 함께 작성해줘."

**Phase 2 — 공유 엔진**

- 구현: `share` 패키지(그룹, 링크, 이전 plan/apply, 백업), `ccam share`, `doctor` 기본 점검
- 완료 기준: settings·CLAUDE.md·확장 기능을 켜고 끄기를 반복해도 데이터 손실 없음 (통합 테스트)
- 프롬프트: "Phase 2를 구현해줘. 이전 작업은 반드시 Plan 구조체를 먼저 만들고 출력한 뒤 Apply하도록 분리해. 켜기→끄기→켜기 왕복 테스트와 충돌 케이스 테스트를 넣어줘."

**Phase 3 — MCP 공유**

- 구현: 그룹 `mcp.json`, flag 전략(과 Phase 0 결과에 따라 sync 전략), `ccam mcp add|list|remove`, 다음 실행 때 새 서버 감지
- 완료 기준: 그룹에 추가한 MCP 서버가 같은 그룹의 모든 계정 `/mcp`에 보임
- 프롬프트: "Phase 3을 구현해줘. docs/claude-layout.md에서 정한 MCP 전략을 따르고, .claude.json을 수정해야 하면 flock과 백업을 반드시 거쳐."

**Phase 4 — 대화 내역 공유와 계정 이어달리기**

- 구현: `projects/`·`history.jsonl` 공유, `(그룹, 작업 폴더)` 잠금, `ccam next [--auto] [--continue]`
- 완료 기준: 1번 계정 대화를 종료하고 `ccam next --continue`로 2번 계정에서 같은 대화가 이어짐

**Phase 4.5 — 사용량 조회**

- 구현: `usage` 패키지(실시간 조회, 상태줄 캐시, 두 응답 형식 파서), `ccs usage`, `ccs next --auto`, 상태줄 수집 스크립트
- 완료 기준: 계정 5개 조회가 3초 안에 끝나고, 한 계정이 실패해도 나머지는 표시됨. 녹화해 둔 응답 JSON으로 파서 테스트 통과
- 프롬프트: "Phase 4.5를 구현해줘. 사용량 조회는 docs/claude-layout.md에 기록한 응답 형식을 따르고, 실제 API 대신 httptest 서버와 녹화한 JSON으로 테스트해. 토큰 파일은 읽기만 하고 --refresh일 때만 잠금 후 갱신해."

**Phase 5 — TUI**

- 구현: `ccam` 선택 화면, `ccam add` 마법사, 공유 항목 토글 화면
- 완료 기준: 키보드만으로 추가·전환·공유 변경 가능, 80열 터미널에서 깨지지 않음
- 프롬프트: "Phase 5를 bubbletea와 huh로 구현해줘. 화면 로직은 모델 함수로 분리해서 teatest로 테스트해줘."

**Phase 6 — 마감과 배포**

- 구현: `import`, `doctor` 자동 수리, 자동완성, 상태줄 스크립트, PGO 적용, CI 성능 회귀 검사, goreleaser(.deb), `install.sh`, README(한국어/영어)
- 완료 기준: 새 우분투 24.04 컨테이너에서 `install.sh` 한 줄로 설치하고 계정 2개 추가까지 5분 이내

## 테스트와 위험 요소

테스트는 실제 홈 폴더와 실제 Claude 계정 없이 돌아가야 합니다. 그래서 모든 경로를 `CCAM_HOME`으로 바꿀 수 있게 하고, `claude` 대신 받은 인자와 환경변수를 기록만 하는 가짜 실행 파일(`testdata/fake-claude`)을 씁니다.

- **단위 테스트**: 인자 전달(위 표의 모든 입력을 fake-claude가 받은 배열과 비교), config 읽기·쓰기, 슬롯 번호 배정, 이전 계획 생성, MCP JSON 병합
- **통합 테스트**: 임시 `CCAM_HOME`에서 add → share on/off 왕복 → remove 후 파일 상태 비교
- **TUI 테스트**: `teatest`로 키 입력 시나리오 재생
- **수동 점검**: Phase마다 실제 계정 2개로 확인 (체크리스트를 `docs/manual-qa.md`에 유지)
- **CI**: GitHub Actions에서 `go test ./...`, `golangci-lint`, Ubuntu 22.04/24.04, 벤치마크 비교(benchstat)

| 위험 | 대응 |
| --- | --- |
| Claude Code 업데이트로 파일 구성 변경 | `doctor`가 알려진 구성과 비교해 경고, 레이아웃 정의를 코드 한 곳(`layout.go`)에 모음 |
| `.claude.json` 동시 수정으로 손상 | 기본은 flag 전략, sync는 flock + 백업 + 실행 중 확인 |
| 공유 대화 내역 동시 쓰기 충돌 | `(그룹, 작업 폴더)` 잠금과 경고 |
| 토큰 파일 노출 | 프로필 폴더 0700, 토큰 0600, 공유 그룹에 토큰이 링크되지 않도록 허용 목록 방식 |
| 이전 중 데이터 손실 | plan/apply 분리, 항상 백업 후 링크, 삭제는 휴지통 폴더로 이동 |
| 계정 여러 개 사용 조건 | 각 계정의 Anthropic 이용약관을 사용자가 확인하도록 README에 명시 |
| 사용량 엔드포인트 변경·중단 | 상태줄 캐시로 자동 대체, 응답 파서를 한 파일에 두고 형식이 바뀌면 doctor가 알림 |

## Claude Code로 개발 시작하기

이 파일을 `docs/PLAN.md`로 저장하고, 아래 CLAUDE.md를 둔 뒤 Phase 0 프롬프트로 시작합니다.

1. 준비

   ```bash
   sudo snap install go --classic        # 또는 공식 tarball
   go version                            # 1.23 이상
   mkdir -p ~/dev/ccam && cd ~/dev/ccam
   git init && mkdir docs
   # 이 파일을 docs/PLAN.md로 저장
   ```

2. 프로젝트 루트에 `CLAUDE.md` 작성

   ```markdown
   # ccam — Claude Code 다중 계정 관리 도구

   전체 계획은 docs/PLAN.md, Claude Code 폴더 조사 결과는 docs/claude-layout.md를 먼저 읽는다.

   ## 규칙
   - 실행 파일은 ccam, 짧은 별칭 ccs(심볼릭 링크)로 실행해도 똑같이 동작한다
   - ccam <이름|번호>는 run의 줄임이다. 프로필 이름은 하위 명령 이름과 겹칠 수 없다
   - 프로필 이름 뒤의 인자는 파싱하지 않고 배열 그대로 claude에 넘긴다 (cobra의 DisableFlagParsing 또는 cobra 이전에 os.Args 직접 처리)
   - Go 1.23, 모듈 경로 github.com/<내아이디>/ccam
   - 코드는 internal/ 아래 패키지로 나누고 cmd/ccam/main.go는 얇게 유지
   - 모든 경로는 CCAM_HOME(테스트용)과 XDG 변수를 따른다
   - 사용자의 실제 ~/.claude, ~/.local/share/ccam은 절대 수정하지 않는다
   - 파일 쓰기는 fsutil.AtomicWrite, 데이터 이전은 Plan → Apply 순서
   - 토큰 파일(.credentials.json)은 어떤 경우에도 공유 그룹에 링크하지 않는다
   - 사용자에게 보이는 메시지는 한국어, 코드·주석·커밋은 영어

   ## 명령
   - 빌드: go build ./cmd/ccam
   - 테스트: go test ./...
   - 린트: golangci-lint run
   - 벤치마크: go test -bench=. -benchmem ./...
   - 실행 경로 측정: hyperfine 'ccs --dry-run work'

   ## 작업 방식
   - 한 번에 한 Phase만 진행하고, 끝나면 완료 기준을 직접 확인해 보고한다
   - 새 기능은 테스트를 먼저 작성한다

   ## 성능 규칙 (docs/PLAN.md의 성능 설계)
   - ccs <이름> 실행 경로는 10 ms 이하. 이 경로에서는 cobra·bubbletea·lipgloss·net/http를 쓰지 않는다
   - 패키지 init()에서 파일·터미널·네트워크에 접근하지 않는다. 무거운 초기화는 처음 쓸 때 한다
   - 실행 경로에서 읽는 파일은 config.toml 하나, 링크 확인은 lstat만
   - claude 실행은 항상 syscall.Exec, 후처리는 다음 실행 때
   - .claude.json은 통째로 파싱하지 않는다 (캐시 또는 스트리밍 디코더)
   - 네트워크 호출은 병렬 + 제한 시간 + 캐시
   - 실행 경로나 사용량 조회를 바꾸면 benchstat 결과를 커밋 메시지에 남긴다
   ```

3. Claude Code 실행 후 첫 프롬프트

   ```
   docs/PLAN.md와 CLAUDE.md를 읽고 Phase 0부터 시작해줘.
   실험은 /tmp 아래 임시 CLAUDE_CONFIG_DIR로만 하고,
   로그인이 필요한 단계에서는 멈추고 나에게 직접 로그인하라고 알려줘.
   ```

4. Phase가 끝날 때마다 커밋하고, 다음 Phase는 새 세션에서 "docs/PLAN.md의 Phase N을 진행해줘"로 이어갑니다. 진행 중 바뀐 결정은 PLAN.md에 바로 반영하게 하면 세션이 바뀌어도 맥락이 유지됩니다.


## 진행 기록과 결정

구현 중 계획과 달라진 점이나 계획에 없던 결정을 여기에 기록합니다.

- Phase 0: `ccam add`/`login`은 Claude Code 화면을 여는 대신 `claude auth login`을 실행합니다 (CLI 하위 명령으로 로그인 가능 확인).
- Phase 0: `mcp_strategy` 기본값은 `flag`. `--mcp-config`는 여러 번 줘도 모두 적용되고 `.claude.json`의 사용자 MCP와도 함께 로드됩니다.
- Phase 0: 기본 위치에서는 `.claude.json`이 `~/.claude` 안이 아니라 `~/.claude.json`에 있으므로, `import`는 이 파일을 프로필 폴더로 복사해야 합니다.
- Phase 1: 슬롯 번호는 비어 있는 가장 작은 번호를 씁니다 (삭제한 번호를 재사용).
- Phase 1: Phase 5 전까지 `ccam`만 입력하면 선택 화면 대신 목록을 보여줍니다.
- Phase 1: 실행 경로 측정 결과 프로세스 전체 2.5 ms (p95 3.0 ms), 바이너리 4.1 MB.
- Phase 2: Claude Code는 심볼릭 링크된 `settings.json`을 링크를 유지한 채 고쳐 씁니다 (플러그인 설치·해제로 확인). 링크 방식 공유가 안전합니다.
- Phase 2: 확장 기능 항목에 `output-styles/`, `rules/`도 포함합니다 (2.1.302에 존재).
- Phase 2: 공유를 켤 때 그룹에 아직 내용이 없으면 백업 없이 원본을 그룹으로 옮깁니다(첫 공유자). 이미 있으면 새 파일만 옮기고 겹치는 파일은 그룹 것을 유지, 내 것은 `.ccam-backup/`에 남깁니다.
- Phase 2: JSON 병합은 객체는 재귀 병합, 배열은 합집합, 값이 충돌하면 그룹 값을 유지합니다. `history.jsonl` 충돌은 줄 단위 병합이 기본입니다.
- Phase 2: 플러그인 `installPath`는 설치한 계정 폴더 경로로 저장됩니다. 그 계정을 지우면 다른 계정의 플러그인 경로가 깨지므로 `doctor --fix`가 그룹 폴더 경로로 바꿉니다.
- Phase 2: 실행 시 링크가 없기만 하면 바로 만들고, 실제 파일이 있는 등 데이터가 걸린 문제는 실행을 멈추고 `ccam doctor --fix`를 안내합니다. 실행 경로 벤치마크 63 µs/op (링크 확인 포함).
- Phase 3: 파일 잠금은 `gofrs/flock` 대신 `syscall.Flock`을 씁니다 (리눅스 전용이라 의존성 불필요). 잠금 fd의 close-on-exec를 풀어 `exec`된 claude 프로세스가 잠금을 계속 쥐게 하므로, 잠금이 잡혀 있으면 "그 계정이 실행 중"이라는 뜻입니다.
- Phase 3: `.claude.json` 잠금은 Claude Code와 같은 방식(`.claude.json.lock` 디렉터리, 10초 지나면 오래된 잠금)으로 잡습니다.
- Phase 3: 세션 안에서 `claude mcp add`로 추가된 서버는 실행 때 묻지 않고 한 줄 안내(`ccam mcp import <계정> <서버>`)만 보여줍니다. 실행 경로에서 입력을 기다리지 않기 위해서입니다.
- Phase 3: 이름이 겹치면 `.claude.json`의 계정 서버가 `--mcp-config`보다 우선하므로, `mcp import`는 그룹으로 복사한 뒤 계정 쪽 항목을 지웁니다 (백업 후, 실행 중이면 건너뜀).
- Phase 3: `.claude.json`의 키 하나만 꺼내는 바이트 스캐너를 만들었습니다. 1.6 MB 파일에서 json.Decoder보다 6.4배 빠르고(0.77 ms), 퍼징으로 encoding/json과 결과를 대조했습니다. `.claude.json`이 바뀐 직후 실행도 p95 5.5 ms.
- Phase 3: 완료 기준 확인 — 실제 Claude Code 2.1.302로 두 계정 모두 그룹 MCP 서버가 보임 (flag, sync 둘 다).
- Phase 4: `ccam next`는 마지막으로 쓴 계정의 다음 번호(순환)로 갑니다. 로그인 안 된 계정(`.credentials.json` 없음)은 건너뛰고, `--continue`면 같은 대화 내역 그룹의 계정만 고른 뒤 `claude --continue`를 붙입니다.
- Phase 4: `(그룹, 작업 폴더)` 잠금도 claude 프로세스가 쥡니다. 이미 잡혀 있으면 경고하고, 터미널이면 계속할지 묻고, 파이프 입력(`-p` 스크립트 등)이면 경고만 하고 진행합니다.
- Phase 4: 완료 기준 확인 — 실제 Claude Code로 1번 계정에서 시작한 대화를 `ccs next --continue`로 2번 계정에서 이어감.
- Phase 4.5: 사용량 표는 80열에 맞추려고 막대를 8칸으로 줄였습니다. 블록 문자(█░)는 한국어 로캘에서도 1칸으로 셉니다.
- Phase 4.5: 실제 계정 응답에는 `weekly_all`도 `seven_day`도 없고 모델별 `weekly_scoped`만 있었습니다. 이때 주간 칸은 `-`로 두고 `-v`에서 모델별 값을 보여줍니다.
- Phase 4.5: `--refresh`의 토큰 갱신은 실제 서비스로 확인하지 않았습니다 (refresh token이 바뀌어 실행 중인 세션이 로그아웃될 수 있어서). 해당 계정이 실행 중이면 갱신하지 않습니다.
- Phase 4.5: 확인 — 실제 계정으로 `ccs usage` 조회 성공, 실제 Claude Code 상태줄에서 `[계정] · 모델 · 5h N%` 표시와 캐시 저장 확인.
- Phase 5: Charm 라이브러리는 v2(`charm.land/bubbletea/v2`, `lipgloss/v2`)를 씁니다. v2가 Go 1.26 이상을 요구해서 `go.mod`의 Go 버전이 1.26이 됐습니다 (계획의 1.23에서 변경).
- Phase 5: `huh`는 쓰지 않습니다. huh가 끌어오는 클립보드 패키지의 `init()`이 PATH를 뒤져 WSL에서 모든 실행에 약 110 ms를 더했습니다. 추가 마법사는 Bubble Tea로 직접 만들었고, 패키지 init 비용이 5 ms를 넘으면 실패하는 테스트를 넣었습니다.
- Phase 5: 선택 화면에서 고른 동작(실행·추가·공유 변경)은 화면을 닫은 뒤 실행합니다. 공유 변경은 일반 터미널에서 이전 계획을 보여주고 확인을 받습니다.
- Phase 5: 확인 — 80×20 tmux에서 선택 화면·공유 화면·추가 마법사·번호 키 실행을 키보드만으로 확인. 첫 화면 100 ms 이내(teatest).
- Phase 6: `import`는 폴더를 옮기거나 복사하지 않고 제자리에 등록합니다. 토큰을 복사하면 refresh token이 바뀔 때 서로 로그아웃시키기 때문입니다. 기본 `~/.claude`를 가져오면 `~/.claude/.claude.json` → `~/.claude.json` 링크를 만들어 그냥 `claude`와 `ccs <이름>`이 같은 계정을 보게 합니다 (Claude Code가 이 링크를 유지하며 쓰는 것을 확인).
- Phase 6: PGO는 실행 경로 벤치마크 프로파일(`cmd/ccam/default.pgo`)로 적용했고 3% 빨라졌습니다. `CCAM_PPROF`는 `exec` 직전에 프로파일을 닫습니다.
- Phase 6: hyperfine 2.0은 JSON 형식이 바뀌어 `scripts/perf-check.sh`가 두 형식을 모두 읽습니다. 측정 결과 p95 `ccs --dry-run work` 4.2 ms, `ccs work` 4.4 ms, `ccs list` 3.9 ms, 최대 메모리 8 MB.
- Phase 6: Docker가 없어 우분투 24.04 컨테이너 대신 빈 HOME과 로컬 릴리스 서버로 `install.sh`를 확인했습니다 (설치 + 계정 2개 추가 0.2초, 체크섬 위조 시 설치 거부). 실제 컨테이너 확인은 남은 일입니다.
- Phase 6: `doctor`가 사용량 응답 형식 변경과, 확인된 Claude Code 버전(2.1.302)과 다른 버전을 알려 줍니다.
- 최종 검토: ccam이 붙이는 `--mcp-config <경로>`가 뒤따르는 프롬프트·하위 명령을 설정 파일로 삼키는 문제를 실제 Claude Code로 재현하고, `--mcp-config=<경로>` 한 인자로 바꿔 해결했습니다. 함께 고친 것: "내 것으로 덮기·병합" 전에 그룹 파일 백업, `AtomicWrite`가 심볼릭 링크를 따라가 씀(가져온 `~/.claude/.claude.json` 링크 유지), 가져온 계정은 `--refresh` 거부, `rename`·`remove` 때 플러그인 경로 갱신, dotfiles 같은 외부 링크 내용을 그룹으로 복사, JSON을 쓸 때 `&&`·`>`를 이스케이프하지 않음.
- 잔여 수정: 실행 중 판단을 잠금 fd 상속에서 "PID + 프로세스 시작 시각" 표시 파일로 바꿨습니다(claude가 띄운 백그라운드 작업이 잠금을 계속 쥐던 문제). 공유 변경이 중간에 실패해도 목표 설정을 먼저 저장해 `doctor --fix`로 마저 진행됩니다. 두 번째 계정이 확장 기능 공유에 들어오면 플러그인 목록을 합칩니다.
- 플랫폼: 계획의 "v1은 우분투만"을 넓혀 macOS와 다른 리눅스 배포판을 지원합니다. macOS는 로그인이 키체인에 있어 `security`로 읽기만 하고(이름 규칙은 docs/claude-layout.md), 세션 판단은 `kern.proc.pid` 시작 시각을 씁니다. 가져온 `~/.claude`에는 `CLAUDE_SECURESTORAGE_CONFIG_DIR=""`를 넘겨 그냥 `claude`와 같은 키체인 항목을 쓰게 합니다. 확인은 GitHub Actions의 macOS 러너와 데비안·알파인(musl)·페도라·아치 컨테이너에서 전체 테스트로 합니다. 실제 Mac에서 실제 로그인으로는 아직 확인하지 않았습니다.
- 2026-10-19: 소유자 요청으로 짧은 별칭 `ccs`를 없앴습니다. 모든 명령은 `ccam`입니다 (`ccam <이름>`, `ccam usage` 등). `ccam`이 더는 `~/.local/bin/ccs` 링크를 만들지 않고, shim 동기화(`add`·`import`·`rename`·`shims install`)와 `install.sh`가 ccam을 가리키는 예전 `ccs` 링크만 지웁니다. 일반 파일이거나 다른 곳을 가리키는 `ccs`는 건드리지 않습니다. 위 계획 본문의 `ccs` 표기는 당시 기록으로 그대로 둡니다.
- 2026-10-19: 소유자 요청으로 프로그램 출력을 모두 영어로 바꿨습니다 (도움말, 메시지, 오류, 질문, TUI, `install.sh`). 한글은 출력하지 않습니다. `internal/e2e/english_test.go`의 `TestSourceHasNoHangul`이 `cmd/`·`internal/`의 테스트 외 Go 소스(문자열·룬 리터럴은 이스케이프까지, 주석 포함)와 `install.sh`·`scripts/*.sh`에서 한글을 찾으면 실패하고, `TestOutputHasNoHangul`이 빌드한 바이너리로 모든 하위 명령의 도움말과 `list`·`usage`·`doctor`·`--dry-run`·오류 경로 출력을 검사합니다. 시간 표기는 `1h 40m`, `3d 6h`, `12m ago (cached)`이고, 예전 버전이 사용량 캐시에 남긴 한글 모델 이름은 읽을 때 `per-model`로 바꿉니다. 위 계획 본문의 한국어 출력 예시는 당시 기록으로 그대로 둡니다.
- 2026-10-19: 소유자 요청으로 `ccam migrate [폴더]`를 추가했습니다. 기존 Claude Code 설정(기본 `~/.claude`, MCP 서버는 `~/.claude.json`의 `mcpServers`)을 공유 그룹(기본 `default`, `--group`으로 변경)에 복사합니다. 옮기거나 링크하지 않고 복사만 하며 원본 폴더는 읽기만 합니다 (테스트가 원본의 해시·수정 시각을 전후로 비교). 기본 항목은 settings·memory·extensions·mcp이고, 대화 내역은 클 수 있어 `--history`를 줄 때만 복사합니다. 경로는 `share.Layout`만 보므로 `.credentials.json`·`.claude.json`은 어떤 경우에도 복사하지 않습니다. 그룹에 다른 내용이 있으면 `share`와 같은 선택지(keep group / use mine / merge JSON / merge by line)를 묻고, `--yes`면 그룹 것을 유지한 뒤 다시 고르는 방법을 알려 줍니다. 폴더 항목은 `share`처럼 그룹에 없는 파일만 더하고 겹치는 파일은 그룹 것을 유지합니다. 바뀌는 그룹 파일은 `<그룹>/.ccam-backup/<시각>/`에 남깁니다. 그룹에 복사해도 그 항목을 공유하는 계정에만 보이므로, 끝에 항목별로 공유 중인 계정과 `ccam share <이름> --<항목> on` 명령을 보여 줍니다. 첫 `ccam add` 때 `~/.claude`에 설정이 있으면 `ccam migrate` 안내 한 줄만 출력하고 자동으로 옮기지는 않습니다.
- 2026-10-19: 소유자 요청으로 단축 명령을 자동으로 만들지 않습니다. `ccam add` 마법사가 공유 항목 다음에 `Create a shortcut command for <name>? [y/N]`(기본 No)를 묻고, y면 `Shortcut name [cc<이름>]:`(예: work → `ccwork`, 쓸 수 없는 이름이면 기본값 없이 `Shortcut name:`)으로 이름을 직접 정하게 합니다. 사용자에게 보이는 예시는 모두 `ccwork` 형식이고 `claude<번호>`는 예전 링크 호환에만 남습니다. 터미널이 아니면 `--shortcut <명령>`을 줄 때만 만듭니다. 이미 있는 계정은 전용 명령 `ccam shortcut <계정> [명령] [--remove]`로 보기·만들기·이름 바꾸기·지우기를 합니다. 이름은 `[a-z0-9][a-z0-9._-]*`(64자까지)이고 `ccam`·`claude`·ccam 하위 명령, 다른 계정의 단축 명령, PATH의 다른 실행 파일을 가리는 이름은 이유와 함께 거절합니다. 단축 명령은 계정마다 하나로 `config.toml`의 `shortcut`에 저장하고, 설정 형식 버전을 2로 올렸습니다. 버전 1 설정은 CLI가 처음 읽을 때 디스크에 남아 ccam을 가리키는 `claude<번호>`(없으면 `claude-<이름>`) 링크를 그 계정의 단축 명령으로 기록하므로 기존 단축 명령은 그대로 동작하고, 그 전에도 실행 경로는 버전 1이면 예전 이름 규칙으로 계정을 찾습니다. 실행 경로는 `argv[0]`의 이름을 원래 읽던 `config.toml` 한 번으로 찾습니다(측정: `ccam` 경유 단축 명령 p95 4.7 ms, 예전 claude1 대비 중앙값 약 0.2 ms 차이). 설정에 없는 `claude<번호>` 이름으로 실행되면 다른 계정을 열지 않고 오류를 냅니다. `ccam shims install --names`(`claude-<이름>`)는 단축 명령 이름을 직접 정할 수 있게 되어 없앴고, `shims install`은 설정된 단축 명령 링크를 다시 만드는 일만 합니다. 링크 정리는 ccam을 가리키는 예전 `claude<번호>`·`claude-<이름>`·`ccs` 링크 중 설정에 없는 것만 지우고, 다른 파일과 사용자가 만든 링크는 건드리지 않습니다. `ccam import`도 단축 명령을 만들지 않습니다.
- 2026-10-19: 소유자 요청으로 예전 버전(v0.1.0, v0.2.0)과의 호환을 두지 않습니다. 현재 설계만 있습니다. 지운 것: 버전 1 설정의 `claude<번호>`·`claude-<이름>` 링크를 단축 명령으로 기록하던 `shim.Adopt`와 그 저장(`shortcut`의 "unchanged" 경로, `shims install`), 실행 경로가 버전 1 설정에서 예전 이름 규칙으로 계정을 찾던 부분, 설정 형식 버전(`version` 키와 "newer ccam" 검사; 이제 읽지도 쓰지도 않음), shim 동기화가 ccam을 가리키는 예전 `claude<번호>`·`claude-<이름>`·`ccs` 링크를 지우던 부분, `install.sh`의 예전 `ccs` 링크 삭제, 사용량 캐시의 한글 모델 이름을 `per-model`로 바꾸던 부분, 그리고 이것들의 테스트. 남긴 것: 단축 명령을 바꾸거나 지우거나 계정을 지우면 ccam이 만든 링크만 지우고(`shim.Remove`), `shims install`·`doctor --fix`는 빠진 링크를 만들기만 하며 다른 파일은 건드리지 않습니다. 설정에 없는 이름으로 실행됐는데 그 이름이 bin 폴더에서 ccam을 가리키는 링크면(손으로 `config.toml`의 `shortcut`을 바꾼 경우) 기본 계정을 열지 않고 `ccam: <이름> is not the shortcut of any account (see \`ccam list\`, \`ccam shortcut\`)` 오류를 냅니다. 사용량 응답의 예전 형식(`five_hour`·`seven_day`)은 Anthropic API 형식이라 계속 읽습니다.
- 2026-10-19: 소유자 요청으로 `ccam update [--check] [--version vX.Y.Z]`를 추가했습니다. 최신 릴리스는 GitHub API 대신 `https://github.com/jhyunwoo/ccam/releases/latest`의 리다이렉트(`/tag/<태그>`)로 알아냅니다 (API의 시간당 60회 제한을 피함, 10초 제한). 테스트와 미러용으로 `CCAM_RELEASES_URL`이 이 주소를 바꿉니다. `install.sh`의 `CCAM_BASE_URL`은 다운로드 폴더라 "최신" 위치가 없어 이름을 따르지 않았습니다. 실행 중인 버전(`cli.Version`)보다 새 릴리스가 없으면 `ccam <버전> is up to date.`, `--check`는 알려 주기만 합니다. 설치할 때는 `ccam_<os>_<arch>.tar.gz`와 `checksums.txt`를 동시에 받아(2분 제한) SHA-256이 맞을 때만 tar에서 `ccam`을 꺼내, 실행 파일(`os.Executable` + `EvalSymlinks`) 옆에 같은 권한의 임시 파일로 쓰고 이름 바꾸기 한 번으로 교체합니다(`fsutil.AtomicWrite`). 체크섬이 다르거나 파일이 없으면 아무것도 바꾸지 않습니다. 단축 명령은 이 경로를 가리키는 링크라 그대로 새 버전을 실행합니다. `--version`은 그 태그를 설치하며 이전 버전으로 돌아가는 데도 씁니다. 개발 빌드(`dev`)는 비교할 릴리스가 없어 `--version` 없이는 거절합니다. `/usr/bin` 같은 시스템 폴더에 있거나 `dpkg -S`가 소유를 확인한 바이너리는 apt에 맡기고 새 `.deb`를 받아 `sudo apt install ./ccam_<버전>_linux_<아키텍처>.deb`를 하라고 안내합니다. 폴더에 쓸 수 없으면 설치 스크립트를 다시 실행하라고 안내합니다. `--dry-run`은 최신 버전만 조회하고 받을 주소와 바꿀 경로를 보여줍니다. 업데이트는 이 명령으로만 하고, 다른 명령은 업데이트를 확인하지 않습니다 (실행 경로를 빠르게, 오프라인에서도 똑같이 두기 위해). 새 패키지 `internal/update`는 `cli`만 가져오고 실행 경로(`launcher`)는 쓰지 않습니다. `update`는 하위 명령이 되어 계정 이름으로 쓸 수 없습니다.
- 2026-10-19: 전체 점검과 정리. 모든 명령을 가짜 claude로 돌려 보고(82개 경우), 리뷰에서 나온 문제를 고쳤습니다. 동작이 바뀐 것: `--dry-run`을 지원하지 않는 명령(`add`·`remove`·`rename`·`use`·`login`·`logout`·`import`·`mcp add` 등)은 예전에 플래그를 무시하고 실제로 실행했는데, 이제 `--dry-run is not supported by …`로 거절합니다(허용 목록은 `cli.dryRunCommands`, `doctor --fix`는 점검만). `ccam run --dry-run <이름>`과 `ccam run --help`는 ccam이 처리합니다. 계정이 실행 중이면 `rename`·`remove`·`share`를 거절합니다. 브라우저 로그인이나 확인 질문처럼 오래 기다리는 명령은 저장 직전에 설정을 다시 읽어, 그사이 다른 ccam이 쓴 내용을 덮지 않습니다. `list --json`의 `share`는 항목 이름(`mcp`…`history`) → 그룹 또는 `local`입니다. `mcp import`는 그룹에 같은 이름의 다른 서버가 있으면 건너뛰고, `--force`면 `mcp.json`을 백업한 뒤 덮습니다. 초기화 시각이 지난 캐시 사용률은 0%로 봅니다(`next --auto`가 이미 초기화된 계정을 놓치던 문제). `config.toml`의 계정·그룹 이름을 읽을 때 검사해 폴더 밖을 가리키지 못하게 했습니다. 공유·이전 때 밖을 가리키는 상대 링크는 절대 경로로 바꾸고, 새 폴더는 0700으로 만듭니다. 실행 경로 벤치마크는 69.8 µs → 76.6 µs(+9.7%, 명시한 `claude_bin`이 ccam 자신인지 inode로 확인하는 비용), hyperfine p95 4.3~4.8 ms. cobra 템플릿(`SetVersionTemplate` 등)은 text/template을 끌어와 바이너리를 2 MB 늘리므로 쓰지 않습니다(`TestNoTextTemplate`).
- 2026-10-19: 소유자 요청으로 `ccam migrate`가 터미널에서 항목 옵션(`--settings`·`--memory`·`--extensions`·`--mcp`·`--history`)과 `--yes` 없이 실행되면 복사할 항목을 체크 목록으로 고르게 했습니다. 기본 선택은 예전과 같이 대화 내역만 빼고 전부이고, 원본 폴더에 그 항목이 없으면 `not found`로 표시합니다(`migrate.Found`: 파일은 lstat, MCP는 `.claude.json`의 `mcpServers`). Space로 바꾸고 Enter로 계획을 보며, Esc는 `Cancelled.`, 아무것도 고르지 않으면 `Nothing selected.`로 끝납니다. `--dry-run`에서도 고를 수 있습니다. 터미널이 아니거나 항목 옵션·`--yes`를 주면 예전과 같습니다. 체크 목록은 `ccam add` 마법사와 같은 컴포넌트(`tui.checklist`)를 씁니다.
- 2026-10-19: 소유자 요청으로 원격 서버(SSH)에서의 로그인을 돕습니다. `claude auth login`은 서버의 `http://localhost:<포트>/callback`에서 브라우저의 리다이렉트를 기다리는데, 브라우저는 사용자의 컴퓨터에서 돌아 그 주소가 열리지 않습니다. 그래서 `ccam add`·`ccam login`이 claude의 표준 입력을 직접 넘기지 않고 한 줄씩 전달하다가, `http://localhost|127.0.0.1|[::1]:<포트>/callback?code=…&state=…` 형태의 줄(스킴 생략 가능)이 오면 claude에 넘기지 않고 ccam이 서버에서 그 주소를 요청합니다(10초 제한, 리다이렉트는 따라가지 않음, 루프백 외 호스트는 받지 않음). 다른 입력(claude의 "Paste code here" 코드 등)은 그대로 claude에 갑니다. 확인: Claude Code 2.1.303의 `auth login`은 표준 입력이 파이프여도 같은 프롬프트로 기다립니다.
- 2026-10-19: 소유자 요청으로 `ccam migrate`가 복사한 MCP 서버의 저장된 OAuth 로그인(원본 `.credentials.json`의 `mcpOAuth`)도 복사합니다. 그룹에는 넣지 않고(로그인 파일은 계정마다 따로), 그 서버를 그룹으로 받는 각 계정의 `.credentials.json`에 계정에 없는 항목만 더합니다. Claude Code 2.1.303은 이 항목의 키를 `<서버 이름>|sha256({type, url, headers})의 앞 16자`로 만들어(바이너리에서 확인) 폴더와 무관하므로, 원본과 정의가 같은 서버(새로 복사, 이미 같음, use mine)에만 적용됩니다. 플러그인 서버(`plugin:` 키)는 원본에 `plugins/`가 있고 확장 기능을 복사할 때만 따라갑니다. 실행 중인 계정, 로그인 전 계정, 제자리 등록(import)한 계정(그냥 `claude`도 그 파일을 쓸 수 있음), macOS 키체인 로그인은 건너뛰고 이유를 보여 줍니다. 나중에 MCP 공유를 켠 계정은 `ccam migrate`를 다시 실행하면 받습니다.

## 2026-10-20 계획: Windows, 토큰 갱신 주소, 한도 재시도, 폴더별 계정

소유자 요청으로 비교 조사에서 나온 단점 네 가지를 구현합니다. 순서는 작은 것부터이고, 단계마다 테스트를 먼저 쓰고 통과시킨 뒤 다음으로 넘어갑니다.

| 단계 | 내용 | 확인 방법 |
| --- | --- | --- |
| A. 토큰 갱신 주소 | `usage.TokenURL`을 Claude Code 2.1.304가 쓰는 `https://platform.claude.com/v1/oauth/token`으로 바꿉니다(바이너리의 `TOKEN_URL`에서 확인). 429는 "로그인 다시"가 아니라 "잠시 뒤 다시"로 안내합니다. | 단위 테스트(가짜 서버 429·200), 실제 주소에 잘못된 토큰으로 요청해 엔드포인트가 응답하는지 확인(실제 토큰은 쓰지 않음) |
| B. 폴더별 계정 | `config.toml`에 `[[folder]] path, account`를 두고, 이름 없이 실행할 때(`ccam -p …`, `ccam run`, `ccam --dry-run`) 현재 폴더와 가장 길게 겹치는 폴더의 계정을 기본 계정보다 먼저 씁니다. 명령은 `ccam folder [이름] [폴더] [--remove]`. 실행 경로에서 읽는 파일은 여전히 `config.toml` 하나이고 `os.Getwd` 한 번만 더합니다. 선택 화면은 그 계정에 커서를 둡니다. 이름을 준 실행과 단축 명령은 그대로 그 계정입니다. | config·launcher 단위 테스트(하위 폴더, 겹치는 폴더, 경로 경계 `/a/b` vs `/a/bc`), e2e(`ccam folder`, `-p` 실행 시 고른 계정), 벤치마크 비교 |
| C. 한도 재시도 | `ccam --fallback <이름> -p …`(또는 `ccam run --fallback`): `-p`/`--print` 실행에서만 동작합니다. claude를 자식 프로세스로 돌려 표준 출력을 모아 두고, 0이 아닌 코드로 끝나면 그 계정의 사용량을 조회해 5시간·주간 한도가 찼을 때만 기록을 버리고 여유가 가장 많은 다른 계정에서 같은 명령을 다시 실행합니다(계정마다 한 번). 표준 입력은 먼저 다 읽어 두고 매번 다시 넣습니다. 대화형 실행은 지금처럼 `exec`이고, 대신 상태줄이 5시간이나 주간 사용률 90% 이상일 때 `→ ccam next --auto --continue`를 덧붙입니다. 백그라운드 감시는 하지 않습니다. | 가짜 claude가 계정별로 실패·성공하도록 하고 가짜 사용량 서버로 e2e(재시도, 한도가 아닌 실패는 재시도 안 함, 표준 입력 재전달, 종료 코드 전달, `-p` 없으면 거절), 상태줄 단위 테스트 |
| D. Windows | 아래 표. Linux·macOS 동작은 바꾸지 않습니다. | 이 컴퓨터(WSL2)의 Windows 11 호스트에서 교차 컴파일한 테스트 바이너리와 e2e를 실행, 실제 `claude.exe`로 `CLAUDE_CONFIG_DIR` 확인, CI에 `windows-latest` 추가 |

Windows에서 바뀌는 부분:

| 영역 | Linux·macOS | Windows |
| --- | --- | --- |
| claude 실행 | `syscall.Exec` | `exec`이 없어 자식 프로세스로 실행하고 기다린 뒤 종료 코드를 그대로 돌려줍니다. Ctrl+C는 콘솔이 claude에도 보내므로 ccam은 무시합니다. |
| 공유 링크 | 심볼릭 링크 | 심볼릭 링크를 먼저 시도하고, 권한이 없으면(개발자 모드 꺼짐) 폴더는 junction으로 만듭니다. 파일(`settings.json`, `CLAUDE.md`, `history.jsonl`)은 심볼릭 링크만 되므로 개발자 모드를 켜거나 그 항목을 `local`로 두라고 안내합니다. 점검은 junction도 링크로 봅니다. |
| 단축 명령 | `~/.local/bin/<이름>` 링크 | `~/.local/bin/<이름>.exe`. 심볼릭 링크, 안 되면 하드 링크. 하드 링크는 `ccam update` 뒤에 다시 만듭니다. |
| 잠금 | `flock` | `LockFileEx` |
| 실행 중 판단 | PID + 시작 시각 | 같음(`GetProcessTimes`). ccam이 claude를 기다리는 동안 살아 있으므로 표시 파일은 ccam PID를 씁니다. |
| 로그인 입력 가리기 | termios | 콘솔 모드(`ENABLE_ECHO_INPUT`, `ENABLE_LINE_INPUT` 끄기) |
| 상태줄 `--wrap` | `sh -c` | Claude Code처럼 Git Bash, 없으면 PowerShell |
| 상태줄 명령 경로 | 셸 따옴표 | `/` 구분자 경로(bash·PowerShell·cmd에서 모두 동작) |
| `--dry-run` 출력 | sh 형식 | PowerShell 형식 |
| 업데이트 | 이름 바꾸기로 교체 | 실행 중인 `ccam.exe`를 `ccam.exe.old`로 옮기고 새 파일을 놓습니다. `.old`는 다음 실행 때 지웁니다. 자산은 `ccam_windows_<arch>.zip` |
| 설치 | `install.sh` | `install.ps1`(SHA-256 확인, `~/.local/bin`, 사용자 PATH 추가) |
| 폴더 위치 | XDG | 같은 구조(`~/.config/ccam` 등, `%USERPROFILE%` 기준). Claude Code도 `%USERPROFILE%\.claude`, `%USERPROFILE%\.local\bin`을 씁니다. |
| 권한 점검(0600) | doctor가 확인 | Windows에는 해당 없음, 건너뜁니다 |
| 셸 자동 완성 | bash, zsh | powershell 추가 |

문서: README 두 개, 사이트의 설치·명령·FAQ·비교 페이지와 기능 페이지를 고치고 30개 언어로 번역합니다. 비교 페이지의 "Windows는 지원하지 않음"과 "Nothing runs between you and Claude Code"는 Windows 예외를 적어 고칩니다.

진행 기록 (2026-10-20):

- A. 토큰 갱신 주소: `usage.TokenURL`을 `https://platform.claude.com/v1/oauth/token`으로 바꿨습니다. 잘못된 토큰으로 두 주소에 요청하면 둘 다 JSON `rate_limit_error`(429)로 답해, 엔드포인트는 둘 다 살아 있습니다. 갱신 요청이 429를 받으면 `ccam login` 대신 "잠시 뒤 다시"를 안내하고 로그인 파일은 건드리지 않습니다(`TestRefreshRateLimitedSaysWait`).
- B. 폴더별 계정: `ccam folder`. 항목은 심볼릭 링크를 푼 절대 경로로 저장하고, 실행 때는 작업 폴더를 그대로 비교한 뒤 맞는 항목이 없을 때만 링크를 풀어 한 번 더 비교합니다(macOS의 `/var` → `/private/var`). 이름 변경은 항목을 따라가고, 계정 삭제는 그 계정의 항목을 지웁니다. 없는 계정을 가리키는 항목은 그 폴더에서 이름 없는 실행을 멈추고 `doctor`가 오류로 보고합니다. 측정: 항목 없음·불일치·일치 모두 hyperfine 평균 4.0~4.2 ms로 차이 없음.
- C. 한도 재시도: `ccam --fallback [이름] -p …`. 한도 판단은 Claude Code의 메시지(`hit your limit`, `usage limit reached`)나 사용량 조회 결과(5시간·주간 100%) 중 하나입니다. 다음 계정은 5시간 사용률이 가장 낮은 계정, 그다음 사용량을 알 수 없는 계정을 번호순으로 고르고, `-c`·`--continue`·`-r`·`--resume`이 있으면 같은 대화 내역 그룹만 봅니다. 실패한 시도의 출력은 버리고, 한도가 아닌 실패는 출력과 종료 코드를 그대로 넘깁니다. e2e는 가짜 claude(`FAKE_LIMITED`)로 재시도·입력 재전달·전부 한도·대화형 거절을 확인합니다. 상태줄은 5시간이나 주간 90% 이상이고 다른 계정이 있으면 `→ ccam next --auto --continue`를 덧붙입니다.
- 가짜 claude를 셸 스크립트(`testdata/fake-claude`)에서 Go 프로그램(`testdata/fakeclaude`)으로 바꿨습니다. 같은 e2e를 Windows에서도 돌리기 위해서입니다. e2e는 `CCAM_E2E_BIN`·`CCAM_E2E_FAKE`·`CCAM_E2E_ROOT`·`CCAM_E2E_RELEASES`로 미리 빌드한 바이너리를 받을 수 있어, Go가 없는 Windows 호스트에서도 돌아갑니다.
- D. Windows: 이 컴퓨터(WSL2)의 Windows 11 호스트에서 교차 컴파일한 모든 패키지 테스트와 e2e가 통과합니다. 실제 `claude.exe` 2.1.302로 확인한 것: `ccam work --version`과 `ccwork.exe --version`이 실행되고, `ccam work auth status`의 `configDirectory`가 ccam 프로필 폴더이며, `--dry-run`이 출력한 줄을 PowerShell에 그대로 넣으면 실행됩니다. 공유 항목은 실제 심볼릭 링크로 만들어지고 `doctor`가 통과하며, `ccam folder`가 동작합니다. 사용자의 실제 `C:\Users\biost\.claude`와 `.claude.json`의 수정 시각은 전후로 같았습니다. 상태줄 명령(`C:/…/ccam.exe statusline`, `~/…`)은 Git Bash와 PowerShell 양쪽에서 실행됩니다. `install.ps1`은 로컬 릴리스 서버로 설치·덮어쓰기(실행 중 파일은 `.old`로)·체크섬 위조 거부를 확인했습니다. junction 대체 경로는 심볼릭 링크 실패를 흉내 낸 테스트로 확인했고, `RemoveAll`이 junction 너머를 지우지 않는 것도 테스트합니다.
- Windows에서 확인하지 못한 것: 개발자 모드가 꺼진 실제 환경(이 컴퓨터는 심볼릭 링크가 허용됨, 테스트로 대체), 브라우저 로그인과 입력 가리기(대화형 콘솔 필요), 대화형 Claude Code 안의 상태줄 표시, `install.ps1`의 사용자 `PATH` 추가(실제 사용자 환경 변수를 바꾸므로 생략), ARM64.
- 성능: `BenchmarkLaunchPath` main 76.78 µs → 78.18 µs (+1.82%). 처음 구현에서 `share.Check`가 항목마다 `lstat`을 두 번 불러 +24%였고, 첫 `lstat` 결과를 재사용해 고쳤습니다. hyperfine p95 4.5~4.9 ms, 최대 메모리 8.3 MB.

## 2026-10-20 계획 2: clausona 비교에서 나온 보완

clausona(larcane97/clausona v0.12.0-beta) 조사에서 나온 단점을 보완합니다. 원칙은 그대로입니다: 로그인 토큰을 스스로 갱신하지 않고, 실행 경로는 10 ms 안, 셸 함수·백그라운드 프로세스 없음.

| 단계 | 내용 |
| --- | --- |
| 1. 고르는 규칙 | `config.toml`의 `auto_max_usage`(기본 100)와 `auto_strategy`(`headroom` 기본, `expiring`). 계정은 판단하는 모든 창(5시간, 주간, `--model`이 있으면 그 모델의 주간 창)이 임계값 아래면 "여유", 100% 미만이면 "넘침", 100%면 "가득". 여유 계정 중 전략으로 고르고, 없으면 넘침 계정 중 사용률이 가장 낮은 계정. `next --auto`·`--fallback`·`pick`이 같은 규칙(`usage.Policy`)을 씁니다. `next`와 `pick`은 `--max-usage`·`--strategy`로 덮어씁니다. |
| 1. 종료 코드 75 | 고를 계정이 없으면 EX_TEMPFAIL(75)로 끝냅니다(`next --auto`, `pick`, 모든 계정이 한도인 `--fallback`). |
| 1. `--wait[=시간]` | `--fallback`·`next --auto`·`pick`에서 계정이 없으면 가장 빨리 여유가 생기는 시각(막고 있는 창의 초기화 시각)까지 기다렸다가 다시 고릅니다. 기본 최대 5시간, 그 안에 생기지 않으면 75. |
| 1. 429 | 사용량 조회가 429를 받으면 `Retry-After`(없으면 5분) 동안 그 계정을 다시 조회하지 않고 저장된 값을 씁니다. |
| 2. `ccam pick` | 스크립트·병렬 작업자용. 규칙대로 고른 계정 이름을 출력(`--json`도). 잠금 아래에서 고르고, 최근 2분 안에 고른 계정은 뒤로 미뤄 동시에 시작한 작업자들이 서로 다른 계정을 받게 합니다. |
| 3. API 키 계정 | `ccam add <이름> --api --base-url … --key-from env:NAME|command:"…"`. 키는 실행할 때 그 출처에서 읽어 환경 변수로만 넘기고 파일에 저장하지 않습니다. 물려받은 인증·공급자 변수는 지웁니다. |
| 4. `ccam ext` | 읽기 전용 진단: 계정마다 로드되는 스킬·명령·에이전트·플러그인·MCP 서버·훅과 그 출처(그룹, 계정, 프로젝트), 깨진 링크, 여러 곳에 정의된 MCP 이름, 계정 사이 차이. |
| 5. Claude Code 플러그인 | 저장소의 마켓플레이스와 `ccam` 플러그인: 사용량 확인, 다음 계정, 작업자에게 계정 나누기를 알려 주는 스킬과 명령. |
| 6. 문서 | 기능 페이지(pick, API 키 계정, ext)와 바뀐 페이지, 30개 언어, 비교 페이지에 clausona. |

보류: Codex 지원(범위), 비용 집계(ccusage 안내), 쉬는 계정 토큰 자동 갱신(원칙).

진행 기록 (계획 2, 2026-10-20):

- 고르는 규칙: `usage.Policy`(`internal/usage/choose.go`)가 `next --auto`·`--fallback`·`pick`의 공통 규칙입니다. 여유(모든 창이 임계값 아래) → 넘침(100% 미만) → 가득(고르지 않음) 순이고, 전략 `expiring`은 5시간 창이 가장 빨리 초기화되는 계정을 고릅니다. 모델 창은 `--model` 값에 `opus`·`sonnet`·`haiku`·`fable`이 들어 있으면 이름에 같은 낱말이 있는 모델별 주간 창을 봅니다. `FreeAt`은 가득 찬 계정만 봅니다(넘침 계정은 어차피 고를 수 있으므로).
- `--wait`: 막고 있는 창의 초기화 시각 30초 뒤까지 자고 다시 고릅니다. 초기화 시각을 모르면 10분 뒤. `--fallback --wait`은 방금 한도에 걸린 계정이 60초 캐시에서 여유 있어 보일 수 있으므로 처음 한 번은 반드시 잡니다(바로 같은 계정을 다시 돌리는 반복 방지). 시계와 잠은 테스트에서 바꿉니다(`clock`, `sleep`).
- 종료 코드 75: `next --auto`·`pick`에서 고를 계정이 없을 때, `--fallback`에서 모든 계정이 한도일 때. 이전에는 `--fallback`이 마지막 시도의 종료 코드를 돌려줬습니다(동작 변경).
- 429: `<이름>.pause`에 다시 물을 시각(`Retry-After`, 없으면 5분, 최대 1시간)을 적고 그때까지는 `--no-cache`여도 묻지 않습니다. rename·remove가 함께 옮기고 지웁니다.
- `ccam pick`: 잠금(`picks.lock`) 아래에서 고르고 `picks/<이름>`에 시각을 적습니다. 2분 안에 고른 계정은 먼저 뺍니다. 같은 프로세스 안 동시 호출 3개가 서로 다른 계정을 받는 것을 테스트합니다.
- API 키 계정: `config.toml`의 `[profile.api]`(`base_url`, `key_from`, `auth`, `model`). 키는 실행 직전(`launcher.Ready`, dry run이 아닐 때)에 읽어 자식 환경 변수로만 넘기고, 물려받은 인증·공급자 변수 8개를 지웁니다. `command:` 출처는 dry run에서 실행하지 않습니다(e2e가 표시 파일로 확인). login·logout은 거절, usage는 `API`, 자동 선택 대상에서 빠짐.
- `ccam ext`: `internal/ext`가 읽기만 합니다. 공유 폴더 자체가 링크이면 `WalkDir`가 따라가지 않아 명령·에이전트가 빠지던 것을 실제 출력에서 발견해 고쳤습니다(링크를 먼저 풂).
- 플러그인: 저장소 루트 `.claude-plugin/marketplace.json`과 `plugins/ccam`(스킬 하나, `/ccam:usage`·`/ccam:ext`). `claude plugin validate --strict`로 둘 다 통과. 주의: `claude plugin init`은 현재 `CLAUDE_CONFIG_DIR`의 skills 폴더에 만들어서, 이 작업 중 공유 그룹에 생긴 `demo` 폴더를 지웠습니다. 플러그인은 손으로 쓰고 validate로만 확인합니다.
- 확인: 리눅스 전체 테스트, lint(리눅스·Windows), Windows 호스트 전체 테스트 통과. `BenchmarkLaunchPath` main 82.36 µs → 80.90 µs (~, p=0.161), hyperfine p95 5.1~5.7 ms.

## 2026-10-20 계획 3: 그룹, 전체 화면 대시보드, 기능 묶음

주인이 정한 것: 그룹 모델 B(레이어), 전체 화면 대시보드(탭), 숫자 키는 커서 이동·Enter로 실행, 기능 묶음 넷(이어가기 1·2·3·5, 계정 안전 6·7·환경 변수 경고, 사용량 4·9, 실행 8·10). 보고서의 나머지 추천(D2–D12, theme 자동 감지와 `theme` 키, 위저드의 API 단계, `group` 예약, 이름에서 정하는 그룹 색, 같은 변경에서 문서, 버그 먼저)은 추천대로 갑니다.

| 결정 | 내용 |
| --- | --- |
| 그룹 저장 | `[[group]] name, items, members, dir?`. 계정은 여러 그룹에 들 수 있지만 한 항목은 한 그룹에서만 받습니다(겹치면 Load 오류). `Profile.Share`는 메모리에서만 그룹으로부터 풀어 낸 값(값은 그룹 폴더 이름)이라 share.Check/Build/launcher/next는 그대로입니다. Save는 Share가 그룹과 어긋나면 내부 오류를 냅니다(그룹을 거치지 않은 변경을 잡는 안전장치). |
| 옛 형식 | `[profile.share]`만 있는 설정은 Load가 손실 없이 그룹으로 바꿉니다(폴더마다 같은 계정 집합의 항목을 한 그룹으로, 가장 큰 쪽이 폴더 이름, 나머지는 `<폴더>-<항목>`/`<폴더>-<n>`과 `dir`). 쓰는 것은 다음 저장 때이고, 그때 `config.toml.bak`을 한 번 남깁니다. 두 형식이 같이 있으면 오류. property test(2000회)로 모든 계정의 폴더가 그대로인지 확인합니다. |
| 그룹 엔진 | `internal/group`: 옛 설정과 새 설정을 비교해 항목이 바뀐 계정마다 share.Build 계획을 만들고, `share.Options.Pending`으로 앞 계정의 계획이 그룹에 넣을 것을 뒤 계정이 보게 합니다(첫 내용은 그룹 폴더 → `--from` → 기본 계정 → 번호 순). 실행 중인 계정이 있으면 거부, Windows에서 파일 링크가 안 되면 미리 거부(D11). 설정을 먼저 저장하고 순서대로 Apply(중간에 멈추면 doctor --fix가 마무리). |
| 그룹 사이 이동 | 다른 그룹으로 옮길 때 새 그룹 폴더가 비어 있으면 원래 그룹의 내용을 복사해 시작합니다(파일 항목, 그리고 MCP의 mcp.json). MCP 공유를 그만두면 ccam이 `.claude.json`에 복사해 둔 서버를 백업 후 지웁니다(D10). sync stamp에 그룹 폴더가 들어갑니다. |
| 명령 | `ccam group [list|show|add|set|join|leave|rename|remove]`. `set`은 `--add`/`--remove`/`--items`(pflag가 `-history`를 `-h`로 읽기 때문에 `+item -item` 대신 플래그). `join/leave <group> <account>...`로 그룹을 앞에 둡니다. rename은 설정만(폴더는 `dir`로 남음), remove는 계정이 있으면 `--leave` 필요, 폴더는 다른 그룹이 안 쓰면 휴지통으로. |
| 바뀐 명령 | `ccam share <계정>`은 읽기 전용(옛 플래그는 숨기고 group 안내). `ccam add`는 `--group`(반복)·`--no-group`, 기본은 default(없으면 처음 계정이 mcp·settings·memory·extensions로 만듦; 그룹은 있는데 default가 없으면 이름을 바꾼 그룹이 그 폴더를 쓰면 그 그룹, 아니면 아무 그룹도). `--share`는 숨기고 안내. migrate·mcp는 그룹 이름을 폴더로 풉니다. list와 TUI는 그룹 이름을 보여 줍니다. doctor는 계정이 아닌 멤버, 예약어가 된 계정 이름, 모르는 theme을 알립니다. ext는 그룹 이름으로 표시하고 그룹 항목의 문제는 한 번만 씁니다. |
| 계획 출력 | 경로를 데이터 폴더 기준 상대 경로로 쓰고, 같은 종류·같은 두 폴더 사이의 단계는 한 줄로 묶습니다(3계정 extensions 해제가 42줄 → 9줄). |
| 대시보드 | `ccam`(TTY)은 alt screen 탭 화면: Accounts·Groups·Folders·Extensions·Health. `internal/tui`의 `Backend` 인터페이스를 cli가 구현하고(명령과 같은 코드), 실행·추가(로그인)·로그인은 프로그램을 끝낸 뒤 cli가 하고 다시 대시보드로 돌아옵니다(실행만 끝). 테마는 Claude Code의 색 역할(suggestion #b1b9f9/#5769f7 등)을 따르고 오렌지와 ✻는 쓰지 않습니다. `theme` 키(auto 기본, dark, light, *-daltonized, *-ansi), auto는 `RequestBackgroundColor`. 그룹 색은 이름의 해시. |
| 키 | 숫자는 커서 이동(1초 안 연속 입력은 한 숫자, 맞는 번호가 없으면 새로 시작), Enter 실행, n 여유가 가장 많은 계정(next --auto 규칙), o 실행 옵션, u 사용량 상세, a 추가 위저드, l 로그인, s 그룹, d 기본, f 이 폴더, e 이름, c shortcut, x 제거(파괴적 확인, 처음 선택은 No), ? 도움말. 그룹 탭은 항목·계정 두 칸 편집기에서 겹치는 항목·계정을 이유와 함께 막습니다. |
| 위저드 | 이름 → 로그인 방식(지금 로그인 / API 키 / 나중에) → (API: 키 출처, 주소, 인증 방식, 모델) → 그룹(없으면 default가 공유할 항목) → shortcut → 요약. Esc는 한 단계 뒤로. `ccam add`(대화형)와 대시보드가 같은 위저드를 씁니다. |

진행 기록 (계획 3):

- 버그 1–5(보고서 1.2)는 그룹 작업 전에 고쳤습니다: usage의 "Most headroom"이 Policy를 따름, 좁은 폭 redraw, cache 모드의 로그아웃 표시, -v의 출처 줄, migrate의 없는 항목 체크. 6(공유 편집기 정렬)은 새 TUI로, 7(ext 반복)은 그룹 문제를 한 번만 씀으로 해결.
- `BenchmarkLaunchPath`(그룹 해석 추가 후) HEAD 82.18 µs → 82.16 µs (~, p=0.959). `config.Load` 5계정: 그룹 형식 23.5 µs, 옛 형식(변환) 31.1 µs.
- 환경 변수 경고: 로그인 계정을 실행할 때 `CLAUDE_CODE_OAUTH_TOKEN`·`ANTHROPIC_API_KEY`·`ANTHROPIC_AUTH_TOKEN`에 값이 있으면 stderr에 알립니다(getenv만, 파일 읽기 없음).
- 9 균형 전략: `auto_strategy = "balance"`(또는 `--strategy balance`)는 주간 창의 남은 양을 초기화까지 남은 날수로 나눈 값이 큰 계정, 실행 중인 ccam 세션마다 10씩 뺌. `next --auto`는 stderr에, `pick --json`은 `reason`에 고른 이유를 씁니다(모든 전략).
- Claude Code 확인(2.1.304, 읽기 전용 조사 6개, 실험은 scratch 폴더): `oauthAccount`의 accountUuid·organizationUuid·organizationName·organizationType(claude_pro/max/team/enterprise), `claude auth login --email`은 로그인 페이지 미리 채우기뿐, StopFailure(2.1.78+)의 `error`는 모든 한도에서 `rate_limit`이라 `last_assistant_message`("You've hit your session limit · resets …")로 구분, StopFailure 출력은 `terminalSequence`(OSC 9/99/777)만 쓰임, UserPromptExpansion(matcher `^ccam:switch$`)은 `continue:false`로 요청 전에 멈춤, `--resume <절대 경로.jsonl> --fork-session`은 현재 CLAUDE_CONFIG_DIR의 projects/<slug>/에 새 id로 복사, slug 규칙(UTF-16 단위마다 영숫자 외 '-', 200자 넘으면 앞 200자 + '-' + base36(abs(java hash))), `claude agents --json`은 계정 폴더마다 따로이고 supervisor를 띄우지 않음, `get_usage` 제어 요청(2.1.270+)은 Claude Code가 자기 로그인으로 사용량을 조회, file-history/·tasks/·paste-cache/는 세션 id 단위라 공유해도 안전하고 file-history가 없으면 다른 계정에서 /rewind 코드 복원이 실패.
- 6 조직 경계·7 로그인 가드: 로그인할 때 `[profile.login]`(account_uuid, email, org_uuid, org_name, org_type)을 기록. 자동 선택(next --auto, --fallback, pick, 대시보드 n, next --continue)은 같은 팀/엔터프라이즈 조직, 또는 개인 요금제끼리만. 같은 Claude 계정에 로그인한 두 폴더는 한도를 함께 쓰므로 서로 대신하지 않음. `--cross-org`/`cross_org = true`로 풂. 알 수 없는 계정은 어디에나 맞음. `ccam login`은 `--email`을 넘기고 로그인 뒤 기록과 비교(다르면 경고하고 되돌리는 방법 `ccam logout`→`ccam login`을 안내; 묻지 않는 이유는 로그인 중계 goroutine이 아직 stdin을 읽고 있기 때문), 같은 계정의 다른 폴더를 경고. `login --expired`. doctor가 기록과 다른 로그인, 같은 계정 두 폴더를 알림. list는 조직이 둘 이상일 때 Organization 열.
- 3 history 공유 없이 이어가기·2 크기 안내: `next --continue`는 모든 계정에서 이 폴더의 가장 최근 대화(또는 남은 메모)를 찾고, 새 계정이 그 폴더를 보면 `--continue`, 아니면 `--resume <파일> --fork-session`. 대화 크기(마지막 응답의 입력+캐시+출력 토큰)와 제목을 알리고, 조직이 다르면 캐시가 공유되지 않는다고 덧붙임. `internal/transcript`.
- 1 한도 훅·8 /ccam:switch: 플러그인 hooks/hooks.json이 `ccam hook stop-failure`·`ccam hook switch`를 실행. 한도면 `state/handoff/<계정>.json` 메모와 데스크톱 알림, 짧은 throttle·API 키 429는 무시. switch는 메모를 남기고 턴을 멈춤(실제 2.1.304에서 API 호출 0회로 확인). `next --continue`가 12시간 안의 메모를 한 번 쓰고(dry run은 남김), 지정한 계정으로 가며, 떠나는 계정은 --auto 후보에서 뺌. `relay = true`면 대화형 세션을 자식으로 실행해 switch 뒤 `ccam next --continue`로 이어 감(Unix는 훅이 CLAUDE_PID에 TERM, Windows는 /exit 뒤). 기본은 끔(옆에서 아무것도 돌지 않는다는 원칙).
- 4 토큰을 읽지 않는 사용량: `usage_source = "claude"`(또는 `--source claude`)면 계정마다 `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence`에 get_usage를 보내 응답을 읽음. 인증 변수는 지우고, 20초 제한, 캐시 키 `claude`. 기본은 계속 live.
- 5 이어가기 완성: history에 file-history/·tasks/·paste-cache/, settings에 keybindings.json. 나중에 더한 경로(`Late`)는 계정이 자기 사본을 가지고 있어도 실행을 막지 않고 doctor가 경고, `doctor --fix`나 다음 그룹 변경이 옮김. known에 daemon·jobs 등 2.1.304 항목 추가.
- 10 에이전트: `ccam agents [--all] [--json]`(계정마다 `claude agents --json` 병렬, 15초), `ccam attach <id>`(jobs/<id>로 계정을 찾아 `claude attach`).
- 검토 반영(리뷰 워크플로 16건 확인 + 나머지 직접 확인): 지운 계정이 그룹 멤버로 남던 문제, `dir = "local"` 거부, 같은 폴더·같은 항목을 두 그룹이 쓰는 설정 거부, 비어 있는 그룹 폴더도 원래 그룹 내용으로 시작, Late 항목의 다른 링크는 실행을 막지 않음, 그룹→그룹 MCP 이동 때 옛 그룹 사본 제거, 실행 중 계정의 doctor --fix는 QuickFix만, 조직 판단은 지금의 .claude.json 우선(세션 안 /login 반영), 대시보드: dry run이면 추가·로그인도 거부, 백엔드 출력은 버림(alt screen 보호), n이 next --auto처럼 조직·세션 규칙을 따름, 위저드 theme auto가 밝은 배경을 따름, 폴더 탭 열 맞춤, 빈 목록에서 space. API 키 출처 `command:`/`env:` 빈 값 거부. 사용량: 프로브에 WaitDelay 2초, `usage` 끝줄도 조직·세션 규칙, 설정의 `usage_source`는 claude일 때만 따름(cache는 --source 전용, doctor 안내와 일치), FormatChanged는 형식 오류에만, dry run은 Claude Code를 띄우지 않고 저장된 값만, WithoutAuth·API 계정이 지우는 변수에 `*_FILE_DESCRIPTOR`·`ANTHROPIC_UNIX_SOCKET`·`CLAUDE_CODE_USE_ANTHROPIC_AWS`/`_GOOGLE_CLOUD`/`MANTLE`/`GATEWAY` 추가.
- 이어가기 검토 반영: relay는 `CCAM_RELAY=<자기 PID>`를 넘기고 훅이 메모에 적은 값이 같을 때만 이어 가며, 그 메모를 숨은 플래그 `next --note <계정>`으로 정확히 넘김(같은 계정의 다른 세션·같은 폴더의 다른 계정 메모를 잘못 쓰던 문제). 상속된 CCAM_RELAY는 buildEnv가 지움. 훅은 `CLAUDE_PROJECT_DIR`(시작 폴더)를 적고, 메모는 대화 파일의 projects/<slug>로도 폴더를 맞춤(세션 안 cd). 공유 history에서는 `--continue` 대신 `--resume <id>`(폴더의 더 새 대화를 열지 않게). 메모의 대상 계정이 없어지면 경고하고 평소대로 고름. 메모는 실행 준비가 끝난 뒤에 지움(`launcher.LaunchThen`). `ccam agents`는 계정 조회가 실패하면 종료 코드 1. relay된 세션은 처음 인자(`--model` 등)를 넘기지 않음(문서화).
- 보류: 2단계 `--brief`(모델 없는 요약), 대시보드의 에이전트 탭. 실제 로그인이 필요한 확인(실제 get_usage 값, 실제 한도 메시지로 훅 분류, 실제 fork 뒤 cache_read)은 테스트 프로필 로그인 뒤.

## 2026-10-21: `ccam usage` 막대 표시

- 소유자 요청으로 `ccam usage`의 표를 계정별 블록으로 바꿨습니다. 계정 줄(번호, 이름, 기준 시각) 아래에 한도마다 한 줄씩 `5-hour`·`Weekly`·모델별 주간 창(예: `Fable`)을 20칸 막대, 퍼센트, `resets in …`으로 보여 줍니다(대시보드의 사용량 창과 같은 모양). 실제 계정들은 `weekly_all` 없이 Fable 주간 창만 있어 예전 표에서는 주간 칸이 `-`이고 Fable은 `-v`에서만 보였습니다. 이제 모델별 창은 항상 보이고, `-v`는 출처 줄만 더합니다. 이름표는 8칸에서 자릅니다. 80칸 안에 들어가며, 조회 중 다시 그리기는 줄 수가 늘기만 하므로 그대로 동작합니다. `--json` 출력은 바뀌지 않습니다. 사이트(30개 언어)의 출력 예시와 열 설명은 아직 예전 표입니다.

## 2026-10-21: 대시보드 Health 상자, 탭 설명, Accounts의 한도 막대

- 소유자 요청으로 대시보드(`ccam`)의 Health 탭을 계정별 상자로 바꿨습니다. 맨 위 한 줄에 `✔ General   ⚠ c1   ✔ c2 …`로 전체 상태를 요약하고, 상자마다 제목에 가장 나쁜 상태와 문제 수, 줄마다 점검 영역(Folder·Login·Share links 등)과 ✔/ℹ/⚠/✘를 보여 줍니다. 계정 상자는 Folder → Login → Share links 순서이고, 폴더 점검은 문제가 있을 때만 결과를 남기므로 문제가 없으면 `Folder OK` 줄을 넣습니다. 계정과 무관한 결과는 General 상자에 모읍니다. ↑/↓는 상자 사이를 움직이고 선택한 상자의 테두리를 강조색으로 그립니다. 이를 위해 `doctor.Finding`에 `Area`를 더했습니다. `ccam doctor`도 같은 요약 줄과 상자를 출력합니다(`tui.HealthReport`, 예전 `✓ 이름: 내용` 줄과 `Finding.String`은 없앰). 터미널이면 그 폭에 맞춰 줄을 바꾸고, 파이프로 보내면 색 없이 줄을 바꾸지도 자르지도 않으며 모든 상자를 가장 넓은 줄에 맞춥니다(grep으로 메시지를 찾을 수 있도록). 색은 터미널 자체 팔레트(`light-ansi`)이고 `NO_COLOR`를 따릅니다. 종료 코드와 `--fix` 안내 줄은 그대로입니다.
- 모든 탭의 이름 줄 아래에 그 탭이 무엇을 하는지 한 문장(`tabAbout`, 79칸 이하)을 흐리게 보여 줍니다.
- Accounts 탭의 사용량 열을 `5-hour`·`Weekly`·`Fable` 세 칸(6칸 막대와 퍼센트)으로 바꿨습니다. Fable 칸은 이름에 `fable`이 들어간 모델별 주간 창입니다. 사용량 상세 창의 `7-day`도 `Weekly`로 맞췄습니다.

## 2026-10-21: Health 탭에서 고른 상자만 수리

- 소유자 요청으로 대시보드 Health 탭의 `f`가 선택한 상자만 수리합니다. 계정 상자면 그 계정의 결과만, General 상자면 계정과 무관한 결과(플러그인 경로, 단축 명령)만 `ccam doctor --fix`처럼 고치고, 나머지 계정은 점검만 합니다. 문제가 없는 상자에서는 `nothing to repair in <이름>`만 보여 줍니다. 예전처럼 모두 고치는 키는 `F`입니다. 이를 위해 `doctor.Options.FixOnly`(계정 이름, 또는 `doctor.General`)를 더했고 `Backend.Doctor`가 그 값을 받습니다. `ccam doctor --fix` 명령은 그대로 전부 고칩니다.

## 2026-10-21: 대시보드의 업데이트 권유

- 소유자 요청으로 "ccam은 스스로 업데이트를 확인하지 않는다"는 2026-10-19 결정을 바꿨습니다. 대시보드(터미널에서 `ccam`)가 열릴 때마다 백그라운드에서 GitHub의 최신 릴리스를 묻습니다(`update.Updater.Newer`, `ccam update`와 같은 `/releases/latest` 리다이렉트, 10초 제한). 소유자가 고른 대로 하루 캐시나 끄는 설정은 두지 않았습니다. 첫 화면은 기다리지 않고, 오프라인이거나 개발 빌드면 아무것도 표시하지 않습니다.
- 새 버전이 있으면 제목 줄 아래에 `● ccam <새 버전> is available (you have <지금 버전>) · press U to update`를 보여 주고, `U`는 확인 창을 거쳐 `ccam update`와 같은 단계(`env.runUpdate`: 체크섬 확인, 한 번의 이름 바꾸기, Windows 단축 명령 다시 연결)로 설치합니다. 실행 중인 대시보드는 옛 바이너리이므로 "quit and start ccam again"을 안내합니다. apt로 설치한 경우 등 실패하면 `ccam update`와 같은 안내를 오류로 보여 주고 권유는 남깁니다. dry run에서는 설치를 거절합니다.
- 다른 명령(`ccam <이름>` 실행 경로 포함)은 여전히 네트워크에 접근하지 않습니다.
