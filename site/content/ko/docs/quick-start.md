---
title: 빠른 시작
description: ccam add로 Claude Code 계정 두 개를 추가하고, claude1, claude2, ccam <name>으로 전환하고, 공유할 항목을 고릅니다.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

계정마다 `ccam add`로 추가한 뒤 `claude1`, `claude2` 또는 `ccam <name>`으로 Claude Code를 시작합니다. 대부분은 이것만으로 충분합니다.

## 1. 첫 계정 추가하기

```sh
ccam add
```

마법사가 `personal`이나 `work` 같은 이름과, 이 계정이 다른 계정과 공유할 항목을 묻습니다. 그다음 `claude auth login`을 실행해 브라우저를 열고 로그인하게 합니다. 끝나면 다음과 같이 표시됩니다.

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam은 자체 메시지를 영어로 출력하므로 위와 같은 메시지가 그대로 표시됩니다. 명령과 옵션은 어느 언어에서나 같습니다.

## 2. 두 번째 계정 추가하기

```sh
ccam add work
```

브라우저에서 다른 계정으로 로그인합니다. 이 계정은 `claude2`가 됩니다.

## 3. 전환하기

다음 중 어느 것을 실행해도 work 계정으로 Claude Code가 시작됩니다.

```sh
claude2
ccam work
ccam 2
```

계정 이름 뒤의 내용은 모두 그대로 Claude Code에 전달되므로, 익숙한 옵션을 그대로 쓰면 됩니다.

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

`ccam`만 실행하면 선택 화면이 열립니다. 화살표 키와 Enter를 쓰거나 계정 번호를 누르세요.

## 4. 공유할 항목 고르기

기본적으로 MCP 서버, `settings.json`, 전역 `CLAUDE.md`, 그리고 에이전트, 명령, 스킬, 플러그인이 계정끼리 공유됩니다. 대화 내역은 공유되지 않습니다. 계정을 오가며 대화를 이어 가려면 두 계정 모두에서 켜세요.

```sh
ccam share personal --history on
ccam share work --history on
```

ccam은 무언가를 바꾸기 전에 무엇을 옮길지 먼저 보여 줍니다. 자세한 내용은 [공유](/docs/sharing)를 참고하세요.

## 5. 계정이 한도에 도달하면

```sh
ccam usage
ccam next --continue
```

`ccam usage`는 모든 계정의 5시간 사용량과 주간 사용량을 보여 줍니다. `ccam next --continue`는 다음 계정에서 같은 대화를 엽니다. [사용량 한도](/docs/usage-limits)를 참고하세요.

## 이미 ~/.claude를 쓰고 있나요?

다시 로그인하는 대신 그 폴더를 등록합니다.

```sh
ccam import
```

폴더는 제자리에 그대로 있고, 그냥 `claude`를 실행해도 전과 똑같이 동작합니다.
