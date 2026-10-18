---
title: Claude Code가 사용량 한도에 도달했을 때 할 일
description: Claude Code 계정마다 남은 여유를 확인하고, 한 계정이 한도에 도달하면 다른 계정에서 같은 대화를 이어 갑니다.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 64fe0b561bd1
---

한 계정이 한도에 도달하면 본인이 가진 다른 계정에서 같은 대화를 이어 가세요. 대화 내역 공유를 한 번 켜 둔 뒤 `ccs next --continue`를 실행하면 됩니다. 계정이 하나뿐이라면 Claude Code가 보여 주는 초기화 시각까지 기다려야 합니다.

## 한도에 도달하기 전에

명령 하나로 전환할 수 있도록 계정을 설정해 둡니다.

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

이제 작업하는 동안 상태줄에 계정과 5시간 사용량이 표시되므로, 한도에 갑자기 걸리는 일이 없습니다.

## 한도에 도달했을 때

Claude Code를 종료(`/exit`)한 뒤 다음을 실행합니다.

```sh
ccs usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

여유가 가장 많은 계정에서 대화를 엽니다.

```sh
ccs next --auto --continue
```

Claude Code가 `work` 계정으로 시작되고, 지금까지 나눈 내용을 모두 담은 현재 폴더의 가장 최근 대화가 열립니다.

## 알아 둘 점

- `ccs next --continue`는 작업하던 폴더에서 실행하세요. 대화는 프로젝트 폴더별로 저장됩니다.
- 두 계정으로 같은 대화를 동시에 실행하지 마세요. 그렇게 하려고 하면 ccam이 경고합니다.
- 다른 계정도 같은 그룹에서 대화 내역을 공유하고 있어야 합니다.
- 이렇게 계정을 바꿔 써도 되는지는 각 계정의 약관에 따릅니다.

`ccs usage`와 `ccs next`의 모든 옵션은 [사용량 한도](/docs/usage-limits)를 참고하세요.
