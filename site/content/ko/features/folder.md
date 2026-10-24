---
title: 폴더별 계정
description: ccam folder를 쓰면 계정 이름 없이 Claude Code를 시작할 때 폴더와 그 아래 모든 폴더에서 한 계정을 씁니다. 예를 들어 회사 코드에는 work 계정을 쓸 수 있습니다.
updated: 2026-10-23
order: 11
translation: machine
sourceHash: "48de98d42b99"
---

`ccam folder work ~/code/company`를 실행하면 `~/code/company`와 그 아래 모든 폴더에서 계정 이름 없이 실행할 때 Claude Code가 `work`로 시작합니다. 그 밖의 폴더에서는 전처럼 기본 계정이 시작됩니다. `ccam personal`이나 `ccpersonal` 같은 단축 명령으로 계정 이름을 지정하면 언제나 그 계정이 시작됩니다.

## 사용법

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

다음과 같이 실행하면 폴더에 지정된 계정을 씁니다.

- `ccam -p "…"`처럼 첫 번째 인자가 Claude Code 옵션인 모든 실행
- 이름 없는 `ccam run`과 `ccam --dry-run`
- `ccam --fallback -p "…"`의 첫 번째 시도
- 대시보드(`ccam`만 실행). 커서가 폴더에 지정된 계정에 놓인 상태로 열립니다. 심볼릭 링크를 거쳐 그 폴더로 들어온 경우에도 마찬가지입니다

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `<account>` | 사용할 계정의 이름이나 번호 | |
| `[folder]` | 대상 폴더. 상대 경로나 절대 경로 | 현재 폴더 |
| `--remove` | 정확히 그 폴더의 항목을 지웁니다 | 꺼짐 |
| `--dry-run` | 지원합니다. 저장하지 않고 바뀔 내용만 출력합니다 | 꺼짐 |

## 동작 방식

폴더마다 `~/.config/ccam/config.toml`에 `[[folder]]` 항목이 하나씩 생깁니다. 어차피 모든 실행이 읽는 파일입니다.

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder`는 폴더를 심볼릭 링크까지 풀어낸 절대 경로로 바꿔 저장합니다. 같은 폴더를 다시 지정하면 계정이 바뀝니다.
2. 이름 없이 실행하면 `launcher.HereProfile`을 호출합니다. 항목이 없으면 곧바로 기본 계정을 고릅니다.
3. 항목이 있으면 작업 폴더를 각 `path`와 문자열로 비교합니다. 그 폴더 자체이거나 그 아래 폴더여야 하며, 이름의 앞부분만 같은 폴더는 해당하지 않습니다(`/code/company`는 `/code/companyx`를 포함하지 않습니다). 가장 깊은 항목이 이기므로, `~/code`는 `personal`, `~/code/company`는 `work`로 둘 수 있습니다.
4. 작업 폴더를 있는 그대로 비교해 일치하는 항목이 없을 때만, ccam은 작업 폴더의 심볼릭 링크를 풀어 한 번 더 비교합니다.
5. Windows와 macOS에서는 두 운영체제의 파일 시스템이 보통 그렇듯 대소문자를 구분하지 않고 비교합니다.

`config.toml` 외에는 어떤 파일도 읽지 않으므로, 폴더 안에서 실행해도 측정할 만한 시간이 더 들지 않습니다(항목이 있든 없든 4.0 ms, `hyperfine`, 300회 실행).

## 안전장치와 예외 상황

- `ccam rename`은 폴더 항목을 새 이름으로 옮기고, `ccam remove`는 삭제한 계정의 항목을 지웁니다.
- 항목의 계정이 더 이상 존재하지 않으면 그 폴더에서 이름 없이 실행할 때 `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`를 출력하고 멈춥니다. `ccam doctor`는 이런 항목과, 폴더가 사라진 항목을 보고합니다.
- 위쪽 폴더에서 계정을 물려받는 폴더 안에서 `ccam folder --remove`를 실행하면, 아무것도 지우지 않고 그 위쪽 폴더를 알려 줍니다.
- `ccam next`는 순환 순서나 사용량으로 계정을 고르며, 폴더는 보지 않습니다.

## 관련 문서

- [Claude Code 시작하기](/features/launch)
- [단축 명령](/features/shortcuts)
- [계정 관리](/features/accounts)
- [명령어 레퍼런스](/docs/commands)
