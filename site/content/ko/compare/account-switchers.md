---
title: ccam과 다른 Claude Code 계정 전환 도구 비교
description: claude-swap, CCS, clauth, CAAM, cux 등 Claude Code 계정 전환 도구와 ccam을 비교하고, 어떤 방식이 어떤 용도에 맞는지 설명합니다.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: d058da937218
---

Claude Code 계정을 전환하는 도구는 두 가지 방식 중 하나를 씁니다. 대부분은 하나뿐인 `~/.claude` 폴더 안에서 로그인을 바꿔 끼우므로, 한 번에 한 계정만 활성화됩니다. ccam을 비롯한 나머지는 `CLAUDE_CONFIG_DIR`로 계정마다 폴더를 따로 두므로, 모든 계정이 로그인 상태를 유지하고 여러 계정을 동시에 실행할 수 있습니다. ccam은 여기에 설정과 MCP 서버 공유를 더하고, 다른 계정의 폴더에서 대화를 이어 갈 수 있게 합니다.

## 로그인을 바꿔 끼울까, 폴더를 나눌까?

**로그인 바꿔 끼우기**는 폴더가 하나뿐이므로 설정, MCP 서버, 대화 내역도 한 벌만 유지합니다. 전환이 빠르고, 실행 중인 세션이 새 로그인을 이어받을 수도 있습니다. 하지만 한 번에 한 계정만 활성화되며, 이는 모든 터미널과 편집기에 똑같이 적용됩니다.

로그인을 바꿔 끼우려면 도구가 다른 모든 계정의 로그인 사본을 보관해 두었다가 제자리에 써 넣어야 합니다. Claude Code는 로그인을 갱신할 때마다 그 안의 리프레시 토큰을 새것으로 바꾸므로, 저장해 둔 사본도 실행 중인 세션과 계속 맞춰 두어야 합니다. 전환한 뒤 계정이 로그아웃되거나 `/login`을 다시 해야 하는 문제는 이런 프로젝트들이 가장 자주 고치는 문제 중 하나입니다.

**폴더 나누기**는 [Claude Code 문서에서 설명하는](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) 방법입니다. 모든 계정이 로그인 상태를 유지하고, 업무용 세션과 개인용 세션을 나란히 실행할 수 있습니다. 각 로그인은 자기 폴더에 그대로 있고, Claude Code가 평소처럼 그 자리에서 갱신합니다. 대신 도구가 공유해 주지 않는 한, 폴더마다 설정, MCP 서버, 대화 내역을 따로 가진 채 시작합니다.

## 도구 비교

2026-10-20 기준으로 각 프로젝트의 문서와 소스 코드를 바탕으로 정리했습니다. "문서에 없음"은 프로젝트가 언급하지 않았다는 뜻이며, 기능이 없다는 뜻은 아닙니다.

| 도구 | 방식 | 동시에 여러 계정 | 5시간·주간 사용량 | 한도에 도달했을 때 | 플랫폼 |
| --- | --- | --- | --- | --- | --- |
| ccam | 폴더 분리 | 예 | 예, `ccam usage` | 사용자가 `ccam next --auto --continue`를 실행해 여유가 가장 많은 계정에서 대화를 이어 감. `-p` 작업은 `ccam --fallback`이 그 계정에서 알아서 다시 실행 | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | 로그인 교체, `cswap run`으로 폴더 분리 | 예, `cswap run` 사용 시 | 예 | `cswap auto`가 임계값에서 로그인을 바꾸며, 실행 중인 세션은 계속됨 | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | 폴더 분리 | 예 | 일부 | 계정 풀이 로컬 프록시를 거쳐 요청을 다음 계정으로 넘김 | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | 로그인 교체, `clauth start`로 폴더 분리 | 예, `clauth start` 사용 시 | 예 | 대체 계정 순서를 따라 전환 | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | 로그인 교체 또는 격리된 프로필 | 예, 격리된 프로필 사용 시 | 예 | `caam run`이 비대화형 명령을 다른 계정으로 다시 시도 | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | 로그인 교체, `claude`를 감싸서 실행 | 아니요, 활성 계정은 하나 | 예 | 로그인을 바꾸고 같은 대화를 재개 | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | 로그인 교체 | 아니요 | 예 | 임계값에서 전환(선택 사항) | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | 로그인 교체 | 아니요 | 예 | 100%에서 전환(선택 사항) | macOS |

ccswap은 Codex 로그인도 관리하는 claude-swap의 포크입니다. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile), [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles)도 계정마다 폴더를 따로 둡니다. 이 도구 중 여러 개는 Codex나 Gemini 같은 다른 코딩 에이전트도 다룹니다.

## 구현 방식

| 도구 | 작성 언어 | Claude Code 실행 중 | 로그인을 다른 곳으로 복사 | 계정끼리 공유할 수 있는 것 |
| --- | --- | --- | --- | --- |
| ccam | Go, 정적 바이너리 하나 | 없음: ccam이 자신을 `claude`로 대체(대체할 수 없는 Windows에서는 종료를 기다림) | 아니요 | 항목마다 따로: 계정별로 로컬 또는 이름 붙인 그룹 |
| claude-swap | Python 3.12+, uv 또는 pipx로 설치 | `cswap run`에서는 없음, `cswap auto`와 메뉴 막대 앱은 계속 실행 | 예, 로그인을 바꿔 끼울 때 | `cswap run` 사용 시: 설정, `CLAUDE.md`, 확장을 한꺼번에, 대화 내역은 선택; MCP 서버는 실행할 때마다 복사 |
| CCS | Node.js 18+ 또는 Bun | CCS가 부모 프로세스로 남음, 풀은 로컬 프록시를 추가 | 계정 폴더는 아니요, 풀은 자체 로그인을 보관 | 계정별로 공유 또는 로컬, 대화 내역은 컨텍스트 그룹으로; MCP 서버는 복사 |
| clauth | Rust | clauth가 부모 프로세스로 남음 | 예, 로그인을 바꿔 끼울 때 | `~/.claude` 전체 또는 아무것도 없음 |
| CAAM | Go | 모드에 따라 다름, 선택적 데몬 | 예, 볼트 모드에서 | 모드에 따라 고정; MCP 서버는 실행할 때마다 복사 |
| cux | Go | 래퍼가 `claude`를 감싼 채 남음 | 예 | 전부: 폴더 하나 |
| CCSwitcher, Claude Switcher | macOS 앱(Swift, Python) | 메뉴 막대 앱 | 예 | 전부: 폴더 하나 |

## ccam이 잘하는 것

### 로그인은 Claude Code가 둔 자리에 그대로 있습니다

ccam은 Claude 로그인을 한 곳에서 다른 곳으로 복사하지 않습니다. 각 계정의 로그인은 자기 폴더에, macOS에서는 그 폴더 전용 키체인 항목에 그대로 있으며, 갱신은 Claude Code가 합니다. 어긋날 수 있는 저장 사본이 없습니다. `ccam import`로 등록한 기존 `~/.claude`는 일반 `claude`가 쓰는 키체인 항목을 그대로 쓰므로, 어느 쪽으로 실행해도 로그인 상태가 유지됩니다.

ccam이 만료된 토큰을 직접 갱신하는 것은 `ccam usage --refresh`를 실행했을 때뿐입니다. 먼저 잠금을 걸고, 그 계정이 실행 중이면 갱신을 거부하므로, 실행 중인 세션 아래에서 토큰을 바꾸는 일은 없습니다.

### 사용자와 Claude Code 사이에서 아무것도 실행되지 않습니다

ccam은 파일 하나를 읽고, 링크를 확인한 뒤, 자신을 `claude`로 대체합니다(`exec`). 작업하는 동안에는 ccam 프로세스도, 래퍼도, 프록시도, 데몬도 없습니다. 표준 입력, 시그널, 종료 코드는 Claude Code의 것이고, 요청은 곧바로 Anthropic으로 갑니다. Windows에서는 프로세스를 대체할 수 없으므로, ccam은 같은 콘솔로 Claude Code를 시작하고 종료 코드를 기다릴 뿐 다른 일은 하지 않습니다.

이렇게 계정을 시작할 때 늘어나는 시간은 약 4 ms, 95번째 백분위수에서도 6 ms 미만이며, 최대 메모리는 8 MB입니다(hyperfine으로 [`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) 300회 실행). 같은 노트북에서 프로그램이 아무 일도 하기 전, 인터프리터가 시작하는 데만 Python은 13 ms, Node.js는 31 ms가 걸립니다.

### 항목별로, 계정별로 공유

MCP 서버, `settings.json`, `CLAUDE.md`, 확장(에이전트, 명령, 스킬, 플러그인), 대화 내역은 서로 별개인 다섯 항목입니다. 계정마다 각 항목을 `local` 또는 그룹으로 정하므로, 업무 계정 두 개는 `work` 그룹에서 대화 내역을 공유하고 개인 계정은 자기 내역을 따로 두면서도, 세 계정 모두 같은 MCP 서버를 쓸 수 있습니다.

- MCP 서버는 실행할 때 `--mcp-config`로 넘기므로, 기본적으로 ccam은 Claude Code가 실행 중에 다시 쓰는 파일인 `.claude.json`을 수정하지 않습니다. `ccam mcp add`는 서버를 그룹 전체에 한 번에 추가합니다.
- 나머지 항목은 그룹 폴더를 가리키는 링크이므로, 한 계정에서 바꾼 내용이 다른 계정에도 바로 나타납니다. 실행할 때 복사하는 것은 없습니다.
- 로그인, `.claude.json`, 실행 상태는 절대 공유하지 않습니다.

### 모든 변경은 먼저 보여 주고 백업합니다

`ccam share`와 `ccam migrate`는 무언가를 옮기기 전에 계획을 출력하고 확인을 받습니다. 계정과 그룹의 파일 내용이 다르면 그룹 것을 유지할지, 내 것을 쓸지, JSON을 병합할지 직접 고릅니다. 원본은 `.ccam-backup/<time>/`에, 제거한 계정의 폴더는 휴지통 폴더에 들어가고, 파일은 이름 바꾸기 한 번으로 교체되므로 쓰기가 중간에 끊겨도 이전 파일이 남습니다. `--dry-run`은 실행, 공유, 이전, 단축 명령, 업데이트, 복구가 무엇을 할지 보여 줍니다. `ccam doctor --fix`는 깨진 링크와 권한을 고칩니다.

### 다음 계정에서 같은 대화를

대화 내역을 공유하면 `ccam next --auto --continue`가 모든 계정의 사용량을 한꺼번에 조회하고, 여유가 가장 많은 계정을 열어 대화를 이어 갑니다. 언제 전환할지는 사용자가 정합니다. ccam은 백그라운드에서 사용량을 지켜보지 않습니다. 계정이 90%에 도달하면 상태 표시줄이 이 명령을 제안합니다.

스크립트는 사용자를 기다릴 수 없으므로, `ccam --fallback work -p "…"`가 대신 전환합니다. 계정이 한도에 도달하면 여유가 가장 많은 계정으로 같은 명령을 다시 실행하고, 끝까지 실행된 쪽의 출력만 표시합니다. 사용량 한도가 아닌 실패는 그대로 전달합니다.

### 폴더마다 맞는 계정

`ccam folder work ~/code/company`를 실행하면 그 폴더 안에서 계정 이름 없이 실행할 때마다 `work`를 쓰고, 나머지 프로젝트는 계속 기본 계정을 씁니다. 이 설정은 모든 실행이 읽는 같은 `config.toml`에 들어 있으므로 시간이 들지 않습니다.

### 바이너리 하나, 검증된 다운로드

ccam은 약 10 MB의 정적 바이너리이며, 따로 설치할 런타임이 없고, Linux, macOS, Windows를 지원합니다. 설치 스크립트(`install.sh`, `install.ps1`)와 `ccam update`는 모든 다운로드를 릴리스의 SHA-256 체크섬과 대조하며, ccam이 스스로 업데이트를 확인하는 일은 없습니다. 텔레메트리도 없습니다. ccam이 보내는 요청은 내 계정의 사용량을 Anthropic에 조회하는 것(병렬, 3초 제한, 60초 캐시)과 사용자가 요청한 다운로드뿐입니다.

## 다른 도구가 더 잘하는 점

- claude-swap, clauth, cux, CAAM, CCSwitcher는 계정이 사용량 임계값을 넘으면 실행 중인 대화형 세션을 스스로 전환할 수 있습니다. ccam은 `-p` 실행만 `--fallback`으로 전환하고, 대화형 세션은 사용자에게 맡깁니다.
- 개발자 모드가 꺼진 Windows에서 ccam은 폴더는 공유하지만 `settings.json` 같은 단일 파일은 공유하지 못합니다. claude-swap과 CCS는 이런 경우 파일을 복사합니다.
- CCSwitcher, Claude Switcher, claude-swap은 macOS 메뉴 막대에 사용량을 표시합니다.
- CCS, CAAM, clauth, ccswap, Claude Switcher는 Codex 등 다른 코딩 에이전트도 관리합니다.
- CAAM, CCSwitcher, Claude Switcher는 Homebrew로 설치할 수 있고, clauth와 CAAM은 릴리스 체크섬에 서명합니다. ccam은 설치 스크립트, `.deb`, `go install`로 설치하며, 서명 없이 체크섬만 확인합니다.

## 어떤 도구가 맞을까요?

한 번에 한 계정만 쓰거나, 실행 중인 세션이 자동으로 전환되길 원한다면 로그인을 바꿔 끼우는 도구를 고르세요. Linux, macOS, Windows에서 계정을 나란히 쓰고, 하나의 설정을 공유하고, 다음 계정에서 같은 대화를 이어 가면서, Claude Code 옆에서 아무것도 실행되지 않기를 원한다면 ccam을 고르세요. 아무것도 공유할 필요가 없다면 [평범한 셸 별칭](/compare/aliases)으로 충분합니다.

여러 계정을 써도 되는지는 각 계정의 약관에 따라 다릅니다. 사용하기 전에 확인하세요.
