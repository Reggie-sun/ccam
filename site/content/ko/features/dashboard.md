---
title: 대시보드
description: 터미널에서 ccam을 실행하면 계정, 그룹, 폴더 규칙, 확장, 점검을 보여 주는 전체 화면 대시보드가 열리며, 모든 변경은 계획을 먼저 보여 줍니다.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

표준 입력과 출력이 터미널이면 인자 없는 `ccam`은 전체 화면 대시보드를 엽니다. 터미널이 아니면 `ccam list`를 출력합니다. 대시보드는 모든 계정을 사용량, 그룹과 함께 보여 주며, 명령으로 하는 일을 그대로 할 수 있습니다. 실행, 추가, 로그인, 그룹 변경, 폴더 규칙 지정, 점검과 복구가 모두 됩니다. 실행하면 대시보드가 닫히고, 계정을 추가하거나 로그인할 때는 잠시 터미널로 나갔다가 돌아옵니다.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

첫 화면은 `config.toml`만 읽어 100 ms보다 훨씬 짧은 시간에 그리며, 사용량은 계정마다 응답이 오는 대로 채워집니다. 비율은 80%부터 노란색, 100%에서 빨간색으로 바뀌고, 조회에 실패해 캐시에서 가져온 값은 흐리게 표시됩니다.

## 탭

| 탭 | 보여 주는 내용 | 키 |
| --- | --- | --- |
| Accounts | 번호, 이름, 이메일, 5시간 막대, 7일 사용량, 마지막 사용, 그룹. `✔`는 기본 계정, `▸`는 이 폴더의 계정, `•`는 실행 중 | 아래 참고 |
| Groups | 각 그룹, 그룹이 공유하는 항목과 계정, 어느 그룹에도 속하지 않은 계정 | `enter` 편집, `n` 새로 만들기, `e` 이름 바꾸기, `x` 삭제 |
| Folders | [`ccam folder`](/features/folder) 규칙과 여기에 적용되는 규칙 | `n` 여기서 쓸 계정 지정, `enter` 변경, `x` 삭제 |
| Extensions | [`ccam ext`](/features/ext)가 보여 주는 내용. 확장마다 한 번씩, 출처 그룹과 함께 표시 | `r` 다시 불러오기 |
| Health | [`ccam doctor`](/features/doctor)가 보고하는 내용 | `f` 복구(`doctor --fix`) |

`←`/`→`나 `Tab`으로 탭을 바꾸고, `?`는 현재 탭의 키를 보여 주며, `r`은 짧은 캐시를 쓰지 않고 사용량을 다시 조회하고, `q`나 `Esc`는 종료합니다.

## Accounts 탭의 키

| 키 | 동작 |
| --- | --- |
| `Enter` | 선택한 계정으로 Claude Code 시작 |
| `1`–`9` | 그 번호의 계정으로 이동. 12번은 `1` `2`를 빠르게 입력합니다. 실행은 `Enter`로만 합니다 |
| `↑`/`↓`, `k`/`j` | 이동 |
| `n` | [`ccam next --auto`](/features/next)가 고를 계정으로 시작 |
| `o` | 옵션을 주고 시작: `--continue`, `--resume`, 그 밖의 인자, 또는 시험 실행 |
| `u` | 사용량 자세히 보기: 5시간, 7일, 모델별 구간과 초기화 시각, 수치의 출처 |
| `a` | 계정 추가(아래 마법사) |
| `l` | 다시 로그인 |
| `s` | 계정의 그룹 고르기 |
| `d` | 기본 계정으로 지정 |
| `f` | 현재 폴더에서 쓸 계정으로 지정 |
| `e` | 이름 바꾸기 |
| `c` | 단축 명령 지정 또는 삭제 |
| `x` | 삭제(먼저 묻고, 선택은 No에서 시작) |

## 그룹 바꾸기

Groups 탭에서 `Enter`를 누르면 편집기가 열립니다. 왼쪽에는 항목, 오른쪽에는 계정이 있으며, `Tab`으로 오가고 `Space`로 켜고 끕니다. 계정 중 하나가 이미 다른 그룹에서 받는 항목은 이유와 함께 흐리게 표시되고, 같은 항목을 두 번 받게 되는 계정도 마찬가지입니다. 항목은 한 그룹에서만 받을 수 있습니다. `Enter`를 누르면 계획이 표시됩니다.

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

그룹의 사본과 내용이 다른 파일이 있으면 먼저 어떻게 할지 묻습니다(그룹의 것을 유지, 내 것을 사용, 병합). 그다음 [`ccam group`](/features/share)과 같은 코드로 백그라운드에서 변경을 실행합니다. 설정을 먼저 저장하고, 계정마다 데이터를 차례로 옮기며, 원본은 `.ccam-backup/`으로 갑니다. 실행 중인 계정의 데이터를 옮기게 되는 변경은 그 계정 이름과 함께 거부됩니다.

## 계정 추가

`a`를 누르면 터미널에서 `ccam add`도 쓰는 마법사가 열립니다.

1. 계정 이름.
2. 로그인 방식: 지금 Claude 구독으로 로그인(브라우저에서), [API 키](/features/api), 또는 나중에 Claude 구독으로 로그인.
3. API 키를 고른 경우: 실행할 때마다 ccam이 키를 읽을 곳(`env:NAME` 또는 `command:COMMAND`), 엔드포인트, 키를 보내는 방식, 모델.
4. 들어갈 그룹. `default`가 선택되어 있습니다. 아직 그룹이 없으면 `default`가 공유할 항목을 고릅니다.
5. 원하면 단축 명령.
6. 확인할 요약.

`Esc`는 한 단계 뒤로 가며, 첫 단계에서는 마법사를 닫습니다.

## 색

색은 Claude Code 터미널 UI의 역할을 따르므로, ccam이 Claude Code 옆에서도 어색하지 않습니다. 그룹 이름은 이름에 따라 색이 정해집니다. ccam은 터미널에 배경색을 물어 밝은 테마나 어두운 테마를 고릅니다. 직접 정하려면 `config.toml`에 지정하세요.

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` 테마는 적록 색각 이상을 위한 것이고, `-ansi` 테마는 터미널 자체의 16색을 씁니다. `NO_COLOR`를 설정하면 색을 끕니다.

## 안전장치와 예외 상황

- 대시보드는 명령과 같은 함수를 호출하므로 모든 규칙이 그대로 적용됩니다. 데이터를 옮기기 전에 계획을 보여 주고, 실행 중인 계정은 거부하고, 백업을 남기며, 로그인 토큰은 계정 밖으로 나가지 않습니다.
- 대시보드에는 시험 실행이 없습니다. 이름 없이 `ccam --dry-run`을 실행하면 이 폴더의 계정이 어떻게 시작될지 출력하고, 계정에서 `o`를 누르면 그 한 번의 실행만 시험 실행할 수 있습니다.
- Claude Code가 시작되거나 로그인 화면이 열리기 전에 터미널을 원래 상태로 되돌립니다.
- 화면이 대화 상자를 띄우기에 너무 작으면 탭의 행부터 줄입니다. 줄은 바꿈하지 않고 `…`로 자릅니다.

## 관련 문서

- [그룹으로 공유하기](/features/share)
- [Claude Code 시작하기](/features/launch)
- [계정 관리](/features/accounts)
- [사용량 한도](/features/usage)
