# 수동 점검표

실제 계정 2개로 Phase마다 확인합니다. 자동 테스트가 대신할 수 없는 부분(브라우저 로그인, 실제 Claude Code 화면)만 남겼습니다. 확인 날짜와 Claude Code 버전을 함께 적습니다.

## 설치와 추가

- [ ] 새 사용자(빈 HOME)에서 `install.sh` 한 줄로 설치, `ccam version` 출력
- [ ] `ccam add` 마법사: 이름 → 공유 항목 → `Create a shortcut command for <name>? [y/N]` → (y) `Shortcut name [cc<name>]:` (예: `ccwork`) → 브라우저 로그인 → `✓ <name> (<email>) added`와 `  shortcut: <cmd>` (N이면 `  no shortcut …`)
- [ ] 마법사에서 PATH의 다른 프로그램 이름(예: `ls`)이나 `list`를 단축 명령으로 주면 이유와 함께 다시 물음
- [ ] `ccam shortcut <name> <cmd>`로 만들기·이름 바꾸기, `--remove`로 지우기 → `~/.local/bin`의 링크가 따라 바뀌고 새 이름으로 실행됨
- [ ] 예전 버전에서 만든 `claude1`·`claude2`가 업데이트 후에도 그대로 실행되고 `ccam list`의 Shortcut 칸에 보임
- [ ] 두 번째 계정도 같은 방식으로 추가, `ccam list`에 서로 다른 이메일
- [ ] `ccam import`로 기존 `~/.claude` 등록 후 `ccam main`과 그냥 `claude`가 같은 계정

## 전환

- [ ] 단축 명령을 만든 두 계정이 각자의 단축 명령으로 서로 다른 이메일로 열림 (`/status`)
- [ ] `ccam work --resume`, `ccam 2 -p "hi"` 인자 전달
- [ ] `ccam` 선택 화면: 화살표·Enter, 숫자 키, 80열 터미널에서 줄바꿈 없음
- [ ] 터미널 제목이 `claude · <이름>`

## 공유

- [ ] settings 공유: 한 계정에서 `/config` 변경 → 다른 계정에 반영, 링크 유지
- [ ] `ccam mcp add …` 후 두 계정 모두 `/mcp`에 보임
- [ ] 세션 안 `claude mcp add` → 다음 실행 때 안내 → `ccam mcp import`
- [ ] 대화 내역 공유 켜기 → 1번에서 대화 → 종료 → `ccam next --continue`가 2번에서 이어감
- [ ] 같은 폴더에서 두 계정 동시 실행 시 경고
- [ ] 기존 `~/.claude`(설정·CLAUDE.md·스킬·`~/.claude.json`의 MCP)가 있는 상태에서 첫 `ccam add` → `ccam migrate` 안내 한 줄. `ccam --dry-run migrate`로 계획 확인 → `ccam migrate` → 두 계정의 `/config`·`/mcp`·스킬 목록에 기존 설정이 보이고, `~/.claude`와 그냥 `claude`는 그대로

## 사용량

- [ ] `ccam usage`가 claude.ai 설정 화면의 사용률과 일치
- [ ] `ccam statusline install` 후 상태줄에 `[account] · model · 5h N%`
- [ ] 네트워크를 끊고 `ccam usage` → 캐시 값과 "(cached)" 표시

## 점검

- [ ] 공유 항목 자리에 실제 파일을 만들고 실행 → 멈춤 → `ccam doctor --fix` → 백업 후 정상
- [ ] `ccam remove`가 폴더를 휴지통으로 옮기고 단축 명령 정리

## 기록

| 날짜 | Claude Code | 결과 | 메모 |
| --- | --- | --- | --- |
| 2026-10-18 | 2.1.302 | 부분 | 실제 Claude Code로 MCP 공유(flag·sync), `next --continue`, 상태줄 수집, 사용량 조회, tmux에서 선택 화면·마법사 확인. 브라우저 로그인과 두 번째 실제 계정은 미확인 |
