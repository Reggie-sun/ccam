---
title: 자주 묻는 질문
description: ccam, Claude Code 계정, 사용량 한도, 공유, 개인정보, 지원 시스템에 관한 질문과 답변입니다.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: e07d468d453b
---

## ccam은 Anthropic 공식 도구인가요?

아닙니다. ccam은 비공식 오픈소스 프로젝트로, Anthropic과 제휴 관계가 없으며 Anthropic의 보증을 받지 않았습니다. ccam은 Claude Code 문서에서 여러 계정을 함께 실행하는 방법으로 안내하는 `CLAUDE_CONFIG_DIR` 설정을 사용합니다.

## Claude Code 계정을 여러 개 써도 되나요?

각 계정의 약관과, 회사 계정이라면 소속 조직의 정책에 따라 다릅니다. 여러 계정을 쓰기 전에 확인하세요. ccam은 이미 가지고 있는 계정 간 전환을 쉽게 해 줄 뿐입니다.

## ccam으로 사용량 한도를 우회할 수 있나요?

아닙니다. 각 계정의 한도는 그대로 적용됩니다. ccam은 한도를 보여 주고, 한 계정의 한도가 소진되면 본인 소유의 다른 계정에서 작업을 이어 가게 해 줍니다.

## 컴퓨터 한 대에서 Claude Code 계정 두 개를 쓰려면 어떻게 하나요?

`CLAUDE_CONFIG_DIR`로 계정마다 설정 폴더를 따로 지정하면 됩니다. 이 작업은 ccam이 대신 해 줍니다. 계정마다 `ccam add`를 실행한 뒤 `ccam personal`, `ccam work`로 시작하거나, `ccwork`처럼 직접 이름을 붙인 단축 명령으로 시작하세요. [가이드](/guides/multiple-claude-code-accounts)에서는 직접 설정하는 방법도 설명합니다.

## 다른 계정에서 대화를 이어 가려면 어떻게 하나요?

`ccam share <name> --history on`으로 두 계정이 대화 내역을 공유하도록 설정한 뒤 `ccam next --continue`를 실행하세요. 그러면 Claude Code가 다음 계정에서 현재 폴더의 가장 최근 대화를 엽니다.

## ccam이 제 로그인 정보를 보거나 외부로 보내나요?

ccam은 로그인 정보를 절대 복사하거나 공유하지 않습니다. 계정의 액세스 토큰은 `ccam usage`, `ccam next --auto`, `ccam doctor`를 실행할 때 Anthropic의 사용량 엔드포인트에 해당 계정의 사용량을 조회하는 용도로만 읽습니다. ccam에는 원격 측정(telemetry) 기능이 없으며, 그 밖의 어디에도 데이터를 보내지 않습니다.

## 계정은 어디에 저장되나요?

`~/.local/share/ccam/profiles/<name>/`에 계정마다 폴더가 하나씩 만들어집니다. macOS에서는 Claude Code가 로그인 정보 자체를 키체인에 보관하며, 항목 이름은 그 폴더를 바탕으로 정해집니다.

## 기존 ~/.claude는 어떻게 되나요?

`ccam import`나 `ccam migrate`를 실행하지 않는 한 아무것도 바뀌지 않습니다. 가져오기는 폴더를 원래 위치에서 그대로 등록하므로, `claude`를 그냥 실행하든 `ccam <name>`으로 실행하든 같은 계정과 설정을 계속 사용합니다. `ccam migrate`는 폴더를 읽기만 하면서 설정, `CLAUDE.md`, 확장 기능, MCP 서버를 공유 설정으로 복사합니다.

## ccam은 어떤 시스템을 지원하나요?

Linux(Ubuntu, Debian, Fedora, Arch, Alpine 등)와 macOS를 지원합니다. Windows에서는 WSL 안에서 사용하세요.

## 왜 MCP 서버가 한 계정에서만 보이나요?

세션 안에서 `claude mcp add`로 추가한 서버는 해당 계정에만 저장됩니다. 다음에 그 계정을 시작하면 ccam이 알려 주니, `ccam mcp import <account> <server>`로 서버를 공유 그룹으로 옮기세요.

## ccam을 쓰면 Claude Code가 느려지나요?

ccam을 거쳐 Claude Code를 시작하면 약 4밀리초가 더 걸립니다. 노트북에서 `scripts/perf-check.sh`로 측정한 값입니다. ccam은 작은 파일 하나를 읽고 공유 링크를 확인한 다음 자기 프로세스를 `claude`로 바꾸므로, 작업하는 동안에는 실행되고 있지 않습니다.

## ccam은 어떻게 제거하나요?

계정마다 `ccam shortcut <account> --remove`로 단축 명령을 지우고 `~/.local/bin`에서 `ccam`을 삭제한 뒤, 계정이 더 이상 필요 없다면 `~/.config/ccam`, `~/.local/share/ccam`, `~/.local/state/ccam`도 삭제하세요. 자세한 내용은 [설치](/docs/install#uninstall)를 참고하세요.
