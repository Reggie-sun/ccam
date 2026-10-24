# 수동 점검표

실제 계정 2개로 Phase마다 확인합니다. 자동 테스트가 대신할 수 없는 부분(브라우저 로그인, 실제 Claude Code 화면)만 남겼습니다. 확인 날짜와 Claude Code 버전을 함께 적습니다.

## 설치와 추가

- [ ] 새 사용자(빈 HOME)에서 `install.sh` 한 줄로 설치, `ccam version` 출력
- [ ] `ccam add` 마법사: 이름 → 공유 항목 → `Create a shortcut command for <name>? [y/N]` → (y) `Shortcut name [cc<name>]:` (예: `ccwork`) → 브라우저 로그인 → `✓ <name> (<email>) added`와 `  shortcut: <cmd>` (N이면 `  no shortcut …`)
- [ ] 마법사에서 PATH의 다른 프로그램 이름(예: `ls`)이나 `list`를 단축 명령으로 주면 이유와 함께 다시 물음
- [ ] `ccam shortcut <name> <cmd>`로 만들기·이름 바꾸기, `--remove`로 지우기 → `~/.local/bin`의 링크가 따라 바뀌고 새 이름으로 실행됨
- [ ] 두 번째 계정도 같은 방식으로 추가, `ccam list`에 서로 다른 이메일
- [ ] `ccam import`로 기존 `~/.claude` 등록 후 `ccam main`과 그냥 `claude`가 같은 계정

## 업데이트

- [ ] 이전 릴리스를 `install.sh`로 설치(`CCAM_VERSION=v<이전>`)한 뒤 `ccam update --check` → ``ccam <새 버전> is available (you have <이전>). Run `ccam update`.``
- [ ] `ccam --dry-run update`가 받을 주소와 바꿀 경로만 보여주고 `ccam version`은 그대로
- [ ] `ccam update` → `✓ Updated ccam <이전> → <새 버전> (~/.local/bin/ccam)`, `ccam version`과 단축 명령 모두 새 버전, 다시 실행하면 `is up to date.`
- [ ] `.deb`로 설치한 `/usr/bin/ccam`에서 `ccam update` → 바꾸지 않고 `sudo apt install ./ccam_<버전>_linux_<아키텍처>.deb` 안내
- [ ] 네트워크를 끊고 `ccam update` → `could not reach github.com …` 오류, 바이너리 그대로

## 전환

- [ ] 단축 명령을 만든 두 계정이 각자의 단축 명령으로 서로 다른 이메일로 열림 (`/status`)
- [ ] `ccam work --resume`, `ccam 2 -p "hi"` 인자 전달
- [ ] `ccam` 선택 화면: 화살표·Enter, 숫자 키, 80열 터미널에서 줄바꿈 없음
- [ ] 터미널 제목이 `claude · <이름>`

## relay(`relay = true`)

- [ ] Linux·macOS: `ccam p1 --permission-mode plan --add-dir /tmp`로 연 세션에서 Shift+Tab으로 acceptEdits로 바꾸고 `/ccam:switch p2` → Claude Code가 끝나고 `ccam: carrying the conversation on (ccam next --continue)` 뒤 p2로 같은 대화가 열림. 권한 모드는 acceptEdits, `/tmp`가 추가 폴더로 보이고, 끝나는 동안 친 글자는 입력란에 없음(`ccam: what you typed during the switch was discarded.`)
- [ ] 한 계정으로 두 터미널(서로 다른 폴더)에서 대화를 열고 둘 다 `/ccam:switch` → 각 터미널이 자기 대화를 자기 폴더에서 이어 감
- [ ] 대화 내역을 공유하는 다른 계정이 같은 폴더에서 열려 있어도 hop은 `Run anyway?`를 묻지 않음
- [ ] `/ccam:switch` 직후 Claude Code가 끝나는 동안 Ctrl+C → `ccam: switch cancelled. …`로 끝나고 `ccam next --continue`가 그 대화를 이어 감. 그 사이 창을 닫으면 남는 ccam·claude 프로세스 없음(`ps`)
- [ ] 업데이트 전에 시작한 relay 세션(옛 ccam)에서 `/ccam:switch` → 세션이 닫히지 않고 /exit 안내
- [ ] Windows 11 호스트: `/ccam:switch p2` → `type /exit; ccam then carries this conversation on as p2.` → `/exit` → p2로 이어 감
- [ ] Windows 11 호스트: 위 hop 뒤 p2의 Claude Code에서 Ctrl+C(입력 지우기) → 셸 프롬프트가 돌아오지 않고 Claude Code가 그대로 받음. `/exit` 뒤 셸이 claude의 종료 코드를 받음(`$LASTEXITCODE`)
- [ ] Windows 11 호스트: 바로가기(`ccwork`, 심볼릭 링크와 `mklink /H` 하드 링크 둘 다)로 연 세션에서 `/ccam:switch` → `/exit` → 다른 계정으로 같은 대화가 이어 감
- [ ] Windows: npm으로 설치한 Claude Code(`claude.cmd`)로 `ccam p1 -p "a & echo X> x.txt"` → claude가 그 문장을 그대로 받고 `x.txt`가 생기지 않음. 작업 폴더의 `claude.bat`은 실행되지 않음

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
- [ ] `ccam remove <계정>` → `logged out first`, 휴지통 폴더에 `.credentials.json` 없음, 그 계정의 refresh token이 폐기됨(다른 계정은 그대로)

## 로그인 보존 (테스트 계정으로만)

- [ ] `ccam usage --refresh`(만료된 테스트 계정): 갱신 뒤 `.credentials.json`의 mcpOAuth가 그대로, 바로 `ccam <계정>`이 로그인 없이 열림. 갱신 중 Ctrl+C도 같음
- [ ] `ccam logout <계정>`: MCP 로그인이 있으면 서버 이름과 함께 묻고, n이면 그대로
- [ ] macOS: `ccam add t1` → 로그인 → `ccam rename t1 t2` → `its folder stays …` 출력, `ccam t2`가 로그인된 채 열리고 MCP 로그인도 남음
- [ ] macOS: `ccam remove t2` → 키체인에서 `Claude Code-credentials-<해시>` 항목이 사라짐(`security find-generic-password -s …`가 44로 끝남), `ccam add t1 --no-login`이 로그아웃 상태
- [ ] macOS: 예전 ccam으로 지운 계정의 항목이 남은 상태에서 `ccam add <그 이름>` → `security delete-generic-password …` 안내로 멈춤
- [ ] macOS: `USER`를 지운 셸(`env -u USER`)에서 로그인한 계정을 `ccam usage`가 같은 셸에서 읽음
- [ ] Windows(Developer Mode 꺼짐): `ccam shortcut p1 ccp1` → install.ps1을 다시 실행 → `OK 1 shortcuts now run the new ccam`, `ccam doctor`의 Shortcuts가 OK. 세션이 열린 채로 install.ps1을 두 번 더 실행해도 성공(`ccam.exe.old-<숫자>`)
- [ ] Windows: `CLAUDE_CODE_FORCE_WINDOWS_CREDMAN=1`로 로그인한 계정이 `ccam doctor`에서 Info(Credential Manager), `ccam next`의 후보에 남고 `ccam login --expired`에 들지 않음

## 기록

| 날짜 | Claude Code | 결과 | 메모 |
| --- | --- | --- | --- |
| 2026-10-18 | 2.1.302 | 부분 | 실제 Claude Code로 MCP 공유(flag·sync), `next --continue`, 상태줄 수집, 사용량 조회, tmux에서 선택 화면·마법사 확인. 브라우저 로그인과 두 번째 실제 계정은 미확인 |
