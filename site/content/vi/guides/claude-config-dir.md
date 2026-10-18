---
title: Giải thích về CLAUDE_CONFIG_DIR
description: "CLAUDE_CONFIG_DIR làm gì trong Claude Code, thư mục chứa những gì, macOS giữ đăng nhập ở đâu, và tệp nào có thể dùng chung giữa các tài khoản."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` cho Claude Code biết thư mục nào cần dùng thay cho `~/.claude`. Mọi thứ riêng của tài khoản đều nằm ở đó: thông tin đăng nhập, cài đặt, máy chủ MCP, plugin và lịch sử hội thoại. Trỏ biến này đến các thư mục khác nhau sẽ cho bạn các tài khoản riêng biệt.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Bên trong thư mục có gì

Đã kiểm tra với Claude Code 2.1.302:

| Đường dẫn | Nội dung | Dùng chung được |
| --- | --- | --- |
| `.credentials.json` | Token đăng nhập (Linux), token OAuth của MCP | Không |
| `.claude.json` | Thông tin tài khoản, máy chủ MCP của người dùng, trạng thái theo từng dự án | Không (máy chủ MCP được xử lý riêng) |
| `settings.json` | Quyền, hook, mô hình, dòng trạng thái, plugin | Có |
| `CLAUDE.md` | Chỉ dẫn toàn cục | Có |
| `agents/`, `commands/`, `skills/`, `plugins/` | Tiện ích mở rộng | Có |
| `projects/` | Hội thoại dưới dạng tệp `.jsonl`, bộ nhớ theo từng dự án | Có |
| `history.jsonl` | Lịch sử prompt | Có |
| `sessions/`, `shell-snapshots/`, `cache/` và các mục khác | Trạng thái khi chạy | Không |

## Hai chi tiết dễ gây nhầm lẫn

**Khi không có `CLAUDE_CONFIG_DIR`, `.claude.json` không nằm trong `~/.claude`.** Khi đó Claude Code giữ tệp này ở `~/.claude.json`, bên cạnh thư mục. Nếu bạn bắt đầu dùng `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code sẽ tìm `~/.claude/.claude.json` thay vào đó và không tìm thấy tài khoản của bạn. `ccam import` xử lý việc này bằng cách liên kết hai tệp với nhau.

**Trên macOS, thông tin đăng nhập nằm trong Chuỗi khóa (Keychain), không nằm trong thư mục.** Claude Code lưu nó dưới tên `Claude Code-credentials`, kèm theo một mã băm ngắn của đường dẫn thư mục khi `CLAUDE_CONFIG_DIR` được đặt. Do đó các thư mục riêng biệt sẽ có các mục Chuỗi khóa riêng biệt. Nếu Chuỗi khóa từ chối ghi, Claude Code sẽ dùng `.credentials.json` trong thư mục để thay thế.

## Dùng chung giữa các thư mục

Liên kết tượng trưng (symbolic link) hoạt động được: Claude Code đọc và ghi qua một `settings.json` hoặc thư mục `projects/` được liên kết và giữ nguyên liên kết đó. Máy chủ MCP là ngoại lệ, vì chúng nằm chung `.claude.json` với thông tin đăng nhập; thay vào đó, hãy truyền chúng bằng `--mcp-config`.

ccam được xây dựng dựa trên chính những quy tắc này. Xem [dùng chung](/docs/sharing) để biết cách ccam liên kết thư mục, và [hướng dẫn dùng nhiều tài khoản](/guides/multiple-claude-code-accounts) để có bức tranh toàn cảnh.
