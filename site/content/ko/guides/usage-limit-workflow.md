---
title: Claude Code 사용량 한도에 도달했을 때 할 일
description: Claude Code 계정마다 남은 사용량을 확인하고, 한 계정이 한도에 도달하면 다른 계정에서 같은 대화를 이어 가는 방법입니다.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

한 계정이 한도에 도달하면 본인 소유의 다른 계정에서 같은 대화를 이어 가세요. 대화 내역 공유를 한 번만 켜 두면 `ccam next --continue`로 바로 이어 갈 수 있습니다. 계정이 하나뿐이라면 Claude Code에 표시된 초기화 시각까지 기다려야 합니다.

## 한도에 도달하기 전에

명령 하나로 계정을 전환할 수 있도록 미리 설정해 둡니다.

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

이제 작업하는 동안 상태 표시줄에 계정과 5시간 사용량이 표시되므로, 한도가 다가오는 것을 미리 알 수 있습니다.

## 한도에 도달했을 때

Claude Code를 종료(`/exit`)한 뒤 다음 명령을 실행합니다.

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

여유가 가장 많은 계정에서 대화를 엽니다.

```sh
ccam next --auto --continue
```

Claude Code가 `work` 계정으로 시작되고, 현재 폴더의 가장 최근 대화가 지금까지 나눈 내용 그대로 열립니다.

## 알아 둘 점

- `ccam next --continue`는 작업하던 폴더에서 실행하세요. 대화는 프로젝트 폴더별로 저장됩니다.
- 같은 대화를 두 계정에서 동시에 실행하지 마세요. 그러려고 하면 ccam이 경고합니다.
- 다른 계정도 같은 그룹으로 대화 내역을 공유하고 있어야 합니다.
- 이렇게 계정을 전환해 써도 되는지는 각 계정의 약관에 따라 다릅니다.

`ccam usage`와 `ccam next`의 전체 옵션은 [사용량 한도](/docs/usage-limits)를 참고하세요.
