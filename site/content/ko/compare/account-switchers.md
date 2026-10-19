---
title: ccam과 다른 Claude Code 계정 전환 도구 비교
description: claude-swap, CCS, clauth, CAAM, cux 등 Claude Code 계정 전환 도구와 ccam을 비교하고, 어떤 방식이 어떤 용도에 맞는지 설명합니다.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Claude Code 계정을 전환하는 도구는 두 가지 방식 중 하나를 씁니다. 대부분은 하나뿐인 `~/.claude` 폴더 안에서 로그인을 바꿔 끼우므로, 한 번에 한 계정만 활성화됩니다. ccam을 비롯한 나머지는 `CLAUDE_CONFIG_DIR`로 계정마다 폴더를 따로 두므로, 모든 계정이 로그인 상태를 유지하고 여러 계정을 동시에 실행할 수 있습니다. ccam은 여기에 설정과 MCP 서버 공유를 더하고, 다른 계정의 폴더에서 대화를 이어 갈 수 있게 합니다.

## 로그인을 바꿔 끼울까, 폴더를 나눌까?

**로그인 바꿔 끼우기**는 폴더가 하나뿐이므로 설정, MCP 서버, 대화 내역도 한 벌만 유지합니다. 전환이 빠르고, 실행 중인 세션이 새 로그인을 이어받을 수도 있습니다. 하지만 한 번에 한 계정만 활성화되며, 이는 모든 터미널과 편집기에 똑같이 적용됩니다.

**폴더 나누기**는 [Claude Code 문서에서 설명하는](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) 방법입니다. 모든 계정이 로그인 상태를 유지하고, 업무용 세션과 개인용 세션을 나란히 실행할 수 있습니다. 대신 도구가 공유해 주지 않는 한, 폴더마다 설정, MCP 서버, 대화 내역을 따로 가진 채 시작합니다.

## 도구 비교

2026-10-19 기준으로 각 프로젝트의 문서를 바탕으로 정리했습니다. "문서에 없음"은 프로젝트가 언급하지 않았다는 뜻이며, 기능이 없다는 뜻은 아닙니다.

| 도구 | 방식 | 동시에 여러 계정 | 5시간·주간 사용량 | 한도에 도달했을 때 | 플랫폼 |
| --- | --- | --- | --- | --- | --- |
| ccam | 폴더 분리 | 예 | 예, `ccam usage` | 사용자가 `ccam next --auto --continue`를 실행해 여유가 가장 많은 계정에서 대화를 이어 감 | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | 로그인 교체, `cswap run`으로 폴더 분리 | 예, `cswap run` 사용 시 | 예 | `cswap auto`가 임계값에서 로그인을 바꾸며, 실행 중인 세션은 계속됨 | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | 폴더 분리 | 예 | 일부, macOS 메뉴 막대 앱에서 | 계정 풀이 로컬 프록시를 거쳐 요청을 다음 계정으로 넘김 | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | 로그인 교체, `clauth start`로 폴더 분리 | 예, `clauth start` 사용 시 | 예 | 대체 계정 순서를 따라 전환 | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | 로그인 교체 또는 격리된 프로필 | 예, 격리된 프로필 사용 시 | 예 | `caam run`이 비대화형 명령을 다른 계정으로 다시 시도 | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | 로그인 교체, `claude`를 감싸서 실행 | 문서에 없음 | 예 | 로그인을 바꾸고 같은 대화를 재개 | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | 로그인 교체 | 아니요 | 예 | 문서에 없음 | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | 로그인 교체 | 아니요 | 예 | 100%에서 전환(선택 사항) | macOS |

ccswap은 Codex 로그인도 관리하는 claude-swap의 포크입니다. 이 도구 중 여러 개는 Codex나 Gemini 같은 다른 코딩 에이전트도 다룹니다.

## ccam이 중점을 두는 것

- 모든 계정이 로그인 상태를 유지합니다. 전환해도 로그아웃되지 않으며, `ccam work`와 `ccam personal`을 두 터미널에서 동시에 실행할 수 있습니다.
- MCP 서버, `settings.json`, `CLAUDE.md`, 에이전트, 스킬, 플러그인은 링크를 통해 폴더 간에 공유되며, 항목별·계정별로 정할 수 있습니다. 로그인 정보는 절대 공유하지 않습니다.
- 대화 내역을 공유하면 명령 하나 `ccam next --auto --continue`로 여유가 가장 많은 계정의 폴더에서 같은 대화를 엽니다.
- ccam은 Claude Code를 시작하면서 자기 자신을 Claude Code로 대체하므로 백그라운드에서 아무것도 실행되지 않습니다. 사용량을 지켜보거나 스스로 전환하지 않으며, 언제 전환할지는 사용자가 정합니다.

## 다른 도구가 더 잘하는 점

- claude-swap, CCS, clauth, CAAM, cux는 WSL 없이 Windows에서 실행됩니다. ccam은 그렇지 않습니다.
- claude-swap, clauth, cux는 계정이 사용량 임계값을 넘으면 스스로 전환합니다. ccam은 전환을 사용자에게 맡깁니다.
- CCSwitcher, Claude Switcher, claude-swap은 macOS 메뉴 막대에 사용량을 표시합니다.
- CCS, CAAM, ccswap, Claude Switcher는 Codex 등 다른 코딩 에이전트도 관리합니다.

## 어떤 도구가 맞을까요?

한 번에 한 계정만 쓰거나, 전환이 자동으로 이루어지길 원하거나, WSL 없이 Windows를 쓴다면 로그인을 바꿔 끼우는 도구를 고르세요. 계정을 나란히 쓰고, 하나의 설정을 공유하고, 다음 계정에서 같은 대화를 이어 가고 싶다면 ccam을 고르세요. 아무것도 공유할 필요가 없다면 [평범한 셸 별칭](/compare/aliases)으로 충분합니다.

여러 계정을 써도 되는지는 각 계정의 약관에 따라 다릅니다. 사용하기 전에 확인하세요.
