---
title: 빠른 시작
metaTitle: "빠른 시작: Claude Code 계정 추가와 전환"
description: ccam add로 Claude Code 계정 두 개를 추가하고, ccwork 같은 단축 명령을 붙여 계정을 전환한 뒤, 공유할 항목을 고릅니다.
updated: 2026-10-24
order: 2
translation: machine
sourceHash: "e98d0637356a"
---

`ccam add`로 계정을 하나씩 추가한 뒤, `ccam <name>` 또는 `ccwork`처럼 직접 이름을 붙인 단축 명령으로 Claude Code를 시작합니다. 대부분은 이것만으로 충분합니다.

## 1. 첫 번째 계정 추가하기

```sh
ccam add
```

설정 마법사가 `personal`이나 `work` 같은 계정 이름, 이 계정이 다른 계정과 공유할 항목, 그리고 단축 명령을 만들지 여부를 묻습니다(4단계 참고). 이어서 `claude auth login`이 실행되어 브라우저가 열리면 로그인하면 됩니다. 원격 컴퓨터에서는 브라우저가 로드되지 않는 `http://localhost:…/callback` 페이지에서 멈출 수 있습니다. 그 주소를 터미널에 붙여 넣으면 ccam이 로그인을 마칩니다. 완료되면 다음과 같이 표시됩니다.

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. 기존 설정 가져오기

이전에 Claude Code를 써 왔다면 설정, `CLAUDE.md`, 확장 기능, MCP 서버가 `~/.claude`에 있습니다. 이것을 공유 설정으로 복사하면 모든 계정이 같은 설정을 쓰게 됩니다.

```sh
ccam migrate
```

ccam은 무엇을 복사할지 먼저 보여 주고 확인을 받습니다. `~/.claude`는 읽기만 하고, 로그인 정보는 복사하지 않으며, 대화 내역은 `--history`를 붙였을 때만 복사합니다. 자세한 내용은 [기존 설정 복사하기](/docs/sharing#copy-an-existing-setup)를 참고하세요.

## 3. 두 번째 계정 추가하기

```sh
ccam add work
```

브라우저에서 다른 계정으로 로그인합니다.

## 4. 계정 전환하기

`ccam work`를 실행하면 work 계정으로 Claude Code가 시작됩니다. 더 짧게 쓰고 싶다면 `ccwork`처럼 계정 전용 단축 명령을 만드세요. 공유할 항목을 고르고 나면 설정 마법사가 다음과 같이 묻습니다.

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y`를 누른 뒤 이름을 입력하거나, Enter를 눌러 제안된 이름을 그대로 씁니다.

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

로그인을 마치면 ccam이 단축 명령을 알려 줍니다.

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

이 질문을 건너뛰었더라도 단축 명령은 언제든 추가할 수 있습니다.

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

이제 다음 중 어느 명령을 실행해도 work 계정으로 Claude Code가 시작됩니다.

```sh
ccwork
ccam work
ccam 2
```

`ccam`만 실행하면 [대시보드](/features/dashboard)가 열립니다. 화살표 키나 계정 번호로 계정으로 이동한 뒤 Enter를 누르세요.

단축 명령이나 계정 이름 뒤에 쓴 내용은 모두 그대로 Claude Code에 전달되므로, 평소 쓰던 옵션을 그대로 쓰면 됩니다.

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. 공유할 항목 고르기

기본적으로 MCP 서버, `settings.json`, 전역 `CLAUDE.md`, 그리고 에이전트, 명령, 스킬, 플러그인은 `default` 그룹을 통해 계정 간에 공유됩니다. 대화 내역은 공유되지 않습니다. 계정을 오가며 대화를 이어 가려면 대화 내역을 그룹에 추가하세요.

```sh
ccam group set default --add history
```

ccam은 변경하기 전에 계정마다 무엇을 옮길지 먼저 보여 줍니다. 자세한 내용은 [그룹으로 공유하기](/features/share)를 참고하세요.

## 6. 계정이 한도에 도달했을 때

```sh
ccam usage
ccam next --continue
```

`ccam usage`는 모든 계정의 5시간 사용량과 주간 사용량을 보여 줍니다. `ccam next --continue`는 다음 계정에서 같은 대화를 엽니다. 스크립트에서는 `ccam --fallback work -p "…"`가 알아서 전환합니다. 자세한 내용은 [사용량 한도](/docs/usage-limits)를 참고하세요.

**Team·Enterprise 시트.** 자동 선택(`ccam next`, `--fallback`, `ccam pick`, 대시보드의 `n`)은 같은 Claude Team이나 Enterprise 조직의 시트 사이에서도 계정을 바꿉니다. Anthropic 고객센터는 "Each seat on a Team plan is assigned to one named individual, and only that person may use it"(팀 플랜의 각 시트는 한 명의 지정된 개인에게 할당되며, 그 사람만 사용할 수 있습니다), "Users can't route their work through another user's account to use that user's remaining usage"(사용자는 다른 사용자의 계정을 통해 작업을 라우팅하여 해당 사용자의 남은 사용량을 사용할 수 없습니다)라고 안내합니다([Team 요금제](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129), [Enterprise 요금제](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)도 마찬가지). 따라서 시트 사이에서 자동 선택을 쓰면 Anthropic의 약관을 위반할 수 있습니다. 이 기능에 의존하기 전에 소속 조직과 Anthropic에 확인하세요. 개인 요금제(Pro, Max)에 대해서는 [FAQ의 약관 안내](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)를 참고하세요.

## 이미 ~/.claude를 쓰고 있다면

다시 로그인하지 말고 기존 폴더를 등록하세요.

```sh
ccam import
```

폴더는 원래 위치에 그대로 있고, `claude`를 그냥 실행해도 이전과 똑같이 동작합니다.
