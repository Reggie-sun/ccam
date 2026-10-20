---
title: Giải thích về CLAUDE_CONFIG_DIR
description: "CLAUDE_CONFIG_DIR làm gì trong Claude Code, thư mục chứa những gì, macOS giữ đăng nhập ở đâu, và tệp nào có thể dùng chung giữa các tài khoản."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: ef9317030e42
---

`CLAUDE_CONFIG_DIR` cho Claude Code biết dùng thư mục nào thay cho `~/.claude`. Mọi thứ riêng của tài khoản đều nằm ở đó: thông tin đăng nhập, cài đặt, máy chủ MCP, plugin và lịch sử hội thoại. Trỏ biến này tới các thư mục khác nhau sẽ cho bạn các tài khoản tách biệt.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR có được hỗ trợ chính thức không?

Có. Claude Code liệt kê nó trong danh sách [biến môi trường](https://code.claude.com/docs/en/env-vars) như cách để thay đổi thư mục cấu hình, "hữu ích để chạy song song nhiều tài khoản", và [hướng dẫn xác thực](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) của Claude Code dùng nó để đăng nhập bằng nhiều tài khoản. Bạn có thể đặt nó trong shell hoặc trong cài đặt người dùng, nhưng không thể đặt trong cài đặt của dự án.

## Bên trong thư mục có gì?

Đã kiểm tra với Claude Code 2.1.302:

| Đường dẫn | Nội dung | Dùng chung được |
| --- | --- | --- |
| `.credentials.json` | Token đăng nhập (Linux, Windows), token OAuth của MCP | Không |
| `.claude.json` | Thông tin tài khoản, máy chủ MCP của người dùng, trạng thái theo dự án | Không (máy chủ MCP thì xử lý riêng) |
| `settings.json` | Quyền, hook, model, dòng trạng thái, plugin | Có |
| `CLAUDE.md` | Chỉ dẫn toàn cục | Có |
| `agents/`, `commands/`, `skills/`, `plugins/` | Tiện ích mở rộng | Có |
| `projects/` | Hội thoại dưới dạng tệp `.jsonl`, bộ nhớ theo dự án | Có |
| `history.jsonl` | Lịch sử prompt | Có |
| `sessions/`, `shell-snapshots/`, `cache/` và các thư mục khác | Trạng thái khi chạy | Không |

## Những điểm nào dễ gây nhầm lẫn?

**Khi không có `CLAUDE_CONFIG_DIR`, `.claude.json` không nằm trong `~/.claude`.** Khi đó Claude Code giữ tệp này ở `~/.claude.json`, cạnh thư mục. Nếu bạn bắt đầu dùng `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code sẽ tìm `~/.claude/.claude.json` và không thấy tài khoản của bạn. `ccam import` xử lý việc này bằng cách liên kết hai tệp với nhau.

**Trên macOS, thông tin đăng nhập nằm trong Chuỗi khóa (Keychain), không nằm trong thư mục.** Claude Code lưu nó dưới tên `Claude Code-credentials`, kèm theo một mã băm ngắn của đường dẫn thư mục khi `CLAUDE_CONFIG_DIR` được đặt. Do đó các thư mục riêng biệt sẽ có các mục Chuỗi khóa riêng biệt. Nếu Chuỗi khóa từ chối ghi, Claude Code sẽ dùng `.credentials.json` trong thư mục để thay thế, như [tài liệu của Claude Code](https://code.claude.com/docs/en/authentication#credential-management) mô tả.

**Không phải mọi kiểu đăng nhập đều nằm trong thư mục.** Đăng nhập Claude Console không dùng API key được lưu bên ngoài thư mục, nên các thư mục riêng không tách được hai đăng nhập như vậy. Đăng nhập claude.ai và API key được giữ theo từng thư mục.

## Trên Windows, thư mục nằm ở đâu?

Mặc định là `%USERPROFILE%\.claude`, với thông tin đăng nhập nằm trong `.credentials.json` bên trong. Đặt `CLAUDE_CONFIG_DIR` thành một thư mục khác cho mỗi tài khoản bổ sung, giống như trên Linux và macOS.

## Các thư mục có thể dùng chung tệp không?

Liên kết tượng trưng (symbolic link) dùng được: Claude Code đọc và ghi thông qua `settings.json` hoặc thư mục `projects/` được liên kết, và giữ nguyên liên kết. Máy chủ MCP là ngoại lệ, vì chúng nằm chung `.claude.json` với thông tin đăng nhập; thay vào đó hãy truyền chúng bằng `--mcp-config`.

ccam được xây dựng đúng trên những quy tắc này. Xem [dùng chung](/docs/sharing) để biết cách ccam liên kết thư mục, và [hướng dẫn dùng nhiều tài khoản](/guides/multiple-claude-code-accounts) để có cái nhìn tổng thể.
