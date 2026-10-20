---
title: So sánh ccam với alias của shell
description: "So sánh song song alias của shell dùng CLAUDE_CONFIG_DIR với ccam để chạy nhiều tài khoản Claude Code."
updated: 2026-10-20
order: 1
translation: machine
sourceHash: b9fbcb88f111
---

Mỗi tài khoản một alias của shell là cách đơn giản nhất để chạy Claude Code nhiều tài khoản. ccam bên dưới cũng làm điều tương tự, và bổ sung thêm việc dùng chung, chuyển đổi tài khoản và theo dõi mức sử dụng.

## Thiết lập bằng alias

Thêm một dòng cho mỗi tài khoản bổ sung vào `~/.bashrc` hoặc `~/.zshrc`, mở một terminal mới và chạy alias một lần để đăng nhập:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` vẫn là tài khoản đầu tiên của bạn, còn `claude-work` khởi động tài khoản thứ hai. Đây là cách mà [tài liệu của Claude Code mô tả](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## Thiết lập bằng ccam

Cài đặt ccam và thêm tài khoản, ở đây kèm luôn một lệnh tắt riêng:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` và `ccwork` đều khởi động tài khoản này, và cả hai cũng dùng được trong script.

## So sánh song song

| | Alias của shell | ccam |
| --- | --- | --- |
| Đăng nhập riêng biệt | Có | Có |
| Thiết lập cho mỗi tài khoản | Sửa `~/.bashrc`, đăng nhập | `ccam add` |
| Dùng được trong script và chương trình khác | Không, alias chỉ có trong shell | Có, `ccam work` và các lệnh tắt như `ccwork` là lệnh thực sự |
| Dùng chung cài đặt, `CLAUDE.md`, tiện ích mở rộng | Sao chép hoặc liên kết thủ công | Theo từng mục, có sao lưu |
| Dùng chung máy chủ MCP | Thủ công, cho từng tài khoản | `ccam mcp add` một lần |
| Tiếp tục cuộc hội thoại ở tài khoản khác | Không | `ccam next --continue` |
| Mức sử dụng của mọi tài khoản | Không | `ccam usage` |
| Chọn tài khoản còn nhiều hạn mức nhất | Không | `ccam next --auto` |
| Chạy lại một tác vụ `-p` bằng tài khoản khác khi chạm giới hạn | Không | `ccam --fallback` |
| Một tài khoản cho mỗi thư mục dự án | Không | `ccam folder` |
| Tên tài khoản trên dòng trạng thái | Không | `ccam statusline install` |
| Sửa liên kết bị hỏng | Thủ công | `ccam doctor --fix` |
| Phần mềm bổ sung | Không có | Một tệp thực thi tĩnh |

Hãy chọn alias nếu bạn có hai tài khoản và không muốn dùng chung gì giữa chúng. Hãy chọn ccam nếu bạn muốn các tài khoản hoạt động như một bộ thiết lập duy nhất. Bạn có thể bắt đầu với alias và chuyển sang ccam sau: `ccam import` đăng ký các thư mục `~/.claude-*` có sẵn ngay tại chỗ.

[Cách dùng nhiều tài khoản trên một máy tính](/guides/multiple-claude-code-accounts) giải thích từng bước cả hai cách, còn [ccam và các công cụ đổi tài khoản khác](/compare/account-switchers) so sánh ccam với những công cụ hoán đổi một lần đăng nhập.
