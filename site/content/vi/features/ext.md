---
title: Kiểm tra tiện ích mở rộng
description: ccam ext hiện các skill, lệnh, plugin, máy chủ MCP và hook mà mỗi tài khoản nạp, nguồn gốc của từng mục, và những gì bị hỏng.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` trả lời câu hỏi "vì sao skill này (hay máy chủ MCP này) chạy ở tài khoản này mà không chạy ở tài khoản kia?". Với mỗi tài khoản, nó liệt kê các tiện ích mở rộng mà Claude Code nạp và nguồn gốc của từng mục: một nhóm dùng chung, thư mục riêng của tài khoản, hoặc dự án bạn đang ở. Bên dưới bảng, nó liệt kê những gì khiến một tiện ích mở rộng không được nạp như bạn mong đợi. Nó chỉ đọc.

## Cách dùng

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

Mỗi ô nêu nguồn của mục đó trong tài khoản đó: nhóm dùng chung (`default`), `account` khi nó nằm trong thư mục riêng của tài khoản, hoặc `-` khi tài khoản không có nó. `!` đánh dấu một mục có vấn đề, được giải thích bên dưới bảng.

## Những gì được đọc

| Loại | Vị trí | Nguồn hiển thị |
| --- | --- | --- |
| Skill | mỗi thư mục trong `skills/` có `SKILL.md` | nhóm mà `skills/` liên kết tới, hoặc `account` |
| Lệnh, agent | các tệp `.md` trong `commands/` và `agents/`, thư mục con hiển thị dạng `folder/name` | nhóm, hoặc `account` |
| Plugin | `plugins/installed_plugins.json` | nhóm, hoặc `account` |
| Máy chủ MCP | `mcp.json` của nhóm (khi MCP được dùng chung) và `.claude.json` của tài khoản | nhóm, hoặc `account` |
| Hook | `hooks` trong `settings.json`, đếm theo từng sự kiện | nhóm, hoặc `account` |
| Dự án này | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` và `.mcp.json` trong thư mục hiện tại | liệt kê một lần, cho mọi tài khoản |

## Các vấn đề được báo cáo

- **Liên kết hỏng**: một liên kết trong `skills/`, `commands/` hoặc `agents/` có đích không còn nữa, hoặc một thư mục dùng chung (`skills`, `commands`, `agents`, `plugins`) liên kết tới một nơi không tồn tại (sửa bằng `ccam doctor --fix`).
- **Một thư mục skill không có `SKILL.md`**: Claude Code bỏ qua nó.
- **Một máy chủ MCP được định nghĩa khác nhau ở hai nơi**: `mcp.json` của nhóm và `.claude.json` của tài khoản, hoặc một trong hai và `.mcp.json` của dự án. Cùng một thiết lập ở cả hai nơi (thứ mà `mcp_strategy = "sync"` ghi ra) không phải là vấn đề. Chuyển máy chủ riêng của một tài khoản vào nhóm bằng `ccam mcp import <account> <server>`.
- **`installed_plugins.json` không phải JSON hợp lệ.**

## An toàn và các trường hợp đặc biệt

- `ccam ext` không thay đổi gì và không gửi yêu cầu mạng nào. Nó đọc tên máy chủ và số lượng hook, không bao giờ đọc giá trị của `env` hay header trong định nghĩa MCP.
- Nó hiển thị những gì có trên đĩa. Một plugin đã cài nhưng bị tắt trong cài đặt, hoặc một máy chủ bị cài đặt được quản lý chặn, vẫn xuất hiện.
- Chạy nó trong thư mục dự án để bao gồm thư mục `.claude` và `.mcp.json` của dự án đó.

## Liên quan

- [Dùng chung giữa các tài khoản](/features/share)
- [Máy chủ MCP dùng chung](/features/mcp)
- [Kiểm tra và sửa chữa](/features/doctor)
