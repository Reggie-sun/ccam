---
title: Cách dùng hai tài khoản Claude Code trên một máy tính
description: "Chạy song song Claude Code nhiều tài khoản với CLAUDE_CONFIG_DIR, thủ công bằng alias của shell hoặc tự động bằng ccam."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d82129ddcb4c
---

Hãy cho mỗi tài khoản một thư mục cấu hình riêng. Claude Code đọc thư mục này từ biến môi trường `CLAUDE_CONFIG_DIR`, và mỗi thư mục chứa thông tin đăng nhập, cài đặt và lịch sử riêng. Bạn có thể tự thiết lập bằng một alias của shell, hoặc để ccam làm việc đó.

## Thủ công, bằng alias

Đây là cách mà tài liệu của Claude Code mô tả. Thêm một dòng cho mỗi tài khoản bổ sung vào `~/.bashrc` hoặc `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Mở một terminal mới và chạy `claude-work`. Claude Code sẽ yêu cầu bạn đăng nhập, vì thư mục này còn mới. Từ đó trở đi, `claude` là tài khoản thứ nhất và `claude-work` là tài khoản thứ hai.

Cách này hoạt động tốt với hai tài khoản. Những điều cần biết:

- Mỗi thư mục bắt đầu trống: không có cài đặt, không có `CLAUDE.md`, không có máy chủ MCP, không có agent hay skill. Bạn tự sao chép hoặc liên kết những gì cần, và tự giữ các bản sao đồng bộ với nhau.
- Alias chỉ tồn tại trong shell tương tác. Script, Makefile và các chương trình khác không thấy chúng.
- Máy chủ MCP nằm trong `.claude.json`, cùng với thông tin đăng nhập, nên không thể đơn giản liên kết chúng giữa các thư mục.
- Các cuộc hội thoại nằm lại ở tài khoản nơi chúng bắt đầu, nên bạn không thể tiếp tục một cuộc hội thoại ở tài khoản khác.
- Không có gì hiển thị mức sử dụng của tất cả tài khoản cùng lúc.

## Với ccam

ccam tự động hóa cùng ý tưởng đó và bổ sung những phần còn thiếu:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Giờ đây `claude1` và `claude2`, hoặc `ccs personal` và `ccs work`, sẽ khởi động hai tài khoản. Vì chúng là chương trình thực sự chứ không phải alias, chúng cũng hoạt động trong script. Cài đặt, `CLAUDE.md`, máy chủ MCP và tiện ích mở rộng được dùng chung ngay từ đầu, và bạn có thể bật dùng chung lịch sử hội thoại để tiếp tục một cuộc hội thoại trên tài khoản kia bằng `ccs next --continue`.

Nếu bạn đã dùng `~/.claude`, hãy chạy `ccam import` thay vì thêm lại tài khoản đó.

## Nên chọn cách nào

Với hai tài khoản và không dùng chung gì, một alias là đủ. Khi bạn muốn dùng chung cài đặt và máy chủ MCP, có nhiều tài khoản hơn, xem nhanh mức sử dụng hoặc tiếp tục hội thoại, ccam giúp bạn khỏi phải tự quản lý những việc đó. Phần [so sánh](/compare/aliases) liệt kê song song các điểm khác biệt.
