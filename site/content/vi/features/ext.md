---
title: Kiểm tra tiện ích mở rộng
description: ccam ext hiện các skill, lệnh, plugin, máy chủ MCP và hook mà mỗi tài khoản nạp, nguồn gốc của từng mục, và những gì bị hỏng.
updated: 2026-10-26
order: 14
translation: machine
sourceHash: 7ca357a7f38a
---

`ccam ext` trả lời câu hỏi "Vì sao skill này (hay máy chủ MCP này) chạy ở một tài khoản nhưng không chạy ở tài khoản kia?" Với mỗi tài khoản, nó liệt kê các tiện ích mở rộng mà Claude Code nạp và nguồn gốc của từng mục: một nhóm dùng chung, thư mục riêng của tài khoản, hoặc dự án bạn đang ở. Bên dưới bảng, nó liệt kê những gì khiến một tiện ích mở rộng không được nạp như bạn mong đợi. Nó không thay đổi gì.

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
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

Mỗi ô nêu nguồn của mục đó trong tài khoản đó: nhóm dùng chung (`default`), `account` khi nó nằm trong thư mục riêng của tài khoản, `local` cho một máy chủ MCP mà tài khoản chỉ có cho thư mục này, hoặc `-` khi tài khoản không có nó. Một plugin đã cài nhưng chưa được bật hiển thị `(off)` sau nguồn của nó. `!` đánh dấu một mục có vấn đề, được giải thích bên dưới bảng; các vấn đề không thuộc tài khoản nào, chẳng hạn các máy chủ mà chỉ `claude` thông thường có, nằm cuối cùng và không kèm tên. `ccam ext --json` in ra cùng báo cáo (`accounts`, dự án, `problems`), với `remote` trên các máy chủ MCP mà Claude Code kết nối qua mạng; khi không có tài khoản nào, nó in báo cáo với `accounts: []`.

## Những gì được đọc

| Loại | Vị trí | Nguồn hiển thị |
| --- | --- | --- |
| Skill | mỗi thư mục trong `skills/` có `SKILL.md` | nhóm mà `skills/` liên kết tới, hoặc `account` |
| Lệnh, agent | các tệp `.md` trong `commands/` và `agents/`, thư mục con hiển thị dạng `folder/name` | nhóm, hoặc `account` |
| Plugin | `plugins/installed_plugins.json`, được bật hoặc tắt bởi `enabledPlugins` trong cài đặt của tài khoản và của dự án; một plugin chỉ được cài cho một dự án thì chỉ xuất hiện trong dự án đó | nhóm, hoặc `account`; `(off)` khi chưa được bật |
| Máy chủ MCP | `mcp.json` của nhóm (khi MCP được dùng chung), `.claude.json` của tài khoản (các máy chủ riêng của nó, và những máy chủ nó chỉ có cho thư mục này), và các máy chủ của plugin của nó, dưới dạng `plugin:<plugin>:<server>` | nhóm, `account` hoặc `local` |
| Hook | `hooks` trong `settings.json`, đếm theo từng sự kiện | nhóm, hoặc `account` |
| Dự án này | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` và `.mcp.json` trong thư mục hiện tại | liệt kê một lần; một tài khoản chưa tin cậy thư mục sẽ nhận một vấn đề cho các máy chủ trong `.mcp.json` của nó |

## Các vấn đề được báo cáo

- **Liên kết hỏng**: một liên kết trong `skills/`, `commands/` hoặc `agents/` có đích không còn nữa, hoặc một thư mục dùng chung (`skills`, `commands`, `agents`, `plugins`) liên kết tới một nơi không tồn tại (sửa bằng `ccam doctor --fix`, lệnh này tạo lại thư mục nhóm bị thiếu).
- **Một thư mục skill không có `SKILL.md`**: Claude Code bỏ qua nó.
- **Một máy chủ MCP được định nghĩa khác nhau ở hai nơi**: `mcp.json` của nhóm và `.claude.json` của tài khoản, hoặc một trong hai và `.mcp.json` của dự án. Các định nghĩa giống hệt nhau ở cả hai nơi (thứ mà `mcp_strategy = "sync"` ghi ra) không phải là vấn đề. Chuyển máy chủ riêng của một tài khoản vào nhóm bằng `ccam mcp import <account> <server>`.
- **`installed_plugins.json` không phải JSON hợp lệ**, một plugin mà `enabledPlugins` không nêu tên, và một plugin bị thiếu thư mục.
- **Các máy chủ chỉ một tài khoản tải**: một máy chủ chỉ có trong `.claude.json` của tài khoản này mà các tài khoản khác trong nhóm không tải, và một máy chủ trong phạm vi local của tài khoản cho thư mục này (điều `claude mcp add` làm nếu không có `-s user`). Mỗi trường hợp đi kèm lệnh `ccam mcp import` để dùng chung nó. Các máy chủ trong `~/.claude.json` của `claude` thông thường mà không tài khoản nào tải đi kèm lệnh `ccam migrate --mcp`.
- **Một máy chủ của dự án đang chờ phê duyệt**: tài khoản chưa tin cậy thư mục, nên Claude Code hỏi một lần, và cho đến lúc đó các máy chủ trong `.mcp.json` của dự án phải chờ. Vấn đề cho biết nơi phê duyệt chúng: `"enabledMcpjsonServers"` trong `settings.json` của các tài khoản thuộc nhóm settings, hoặc của tài khoản này. Sự tin cậy dành cho worktree chính của một kho git cũng có hiệu lực trong các worktree khác của nó.
- **Đăng nhập và connector**: một máy chủ MCP cần đăng nhập trong tài khoản này trong khi một tài khoản khác đã đăng nhập vào nó (một lần đăng nhập chỉ có hiệu lực cho một tài khoản; hãy đăng nhập bằng `/mcp`), và các connector claude.ai mà một tài khoản khác đã dùng còn tài khoản này thì chưa.

## An toàn và các trường hợp đặc biệt

- `ccam ext` không thay đổi gì và không gửi yêu cầu mạng nào. Nó đọc tên máy chủ và số lượng hook, không bao giờ đọc giá trị của `env` hay header trong định nghĩa MCP. Để kiểm tra đăng nhập, nó đọc tên các lần đăng nhập MCP trong thông tin đăng nhập của từng tài khoản (trên macOS là từ Keychain), không bao giờ đọc token, và chỉ với các tài khoản tải một máy chủ mà một tài khoản khác cũng tải.
- Nó hiển thị những gì có trên đĩa. Một máy chủ bị cài đặt được quản lý chặn vẫn xuất hiện; một plugin bị tắt trong cài đặt hiển thị `(off)`.
- Chạy nó trong thư mục dự án để bao gồm thư mục `.claude` và `.mcp.json` của dự án đó.

## Liên quan

- [Dùng chung bằng nhóm](/features/share)
- [Máy chủ MCP dùng chung](/features/mcp)
- [Kiểm tra và sửa chữa](/features/doctor)
