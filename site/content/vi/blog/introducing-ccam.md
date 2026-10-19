---
title: Giới thiệu ccam
description: "Vì sao có ccam, ccam làm gì cho người dùng nhiều tài khoản Claude Code, và cách nó không cản trở công việc của bạn."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: b116c91bb8d4
---

ccam là một công cụ dòng lệnh nhỏ dành cho người dùng nhiều hơn một tài khoản Claude Code trên cùng một máy tính: một tài khoản cá nhân và một tài khoản công việc, hoặc nhiều tài khoản cho một dự án dài.

## Vấn đề

Claude Code hỗ trợ nhiều tài khoản thông qua `CLAUDE_CONFIG_DIR`: mỗi tài khoản một thư mục. Trên thực tế, mỗi thư mục là một hòn đảo riêng. Cài đặt, `CLAUDE.md`, máy chủ MCP, agent và skill phải được thiết lập lại trong từng thư mục và đồng bộ thủ công. Các cuộc hội thoại nằm lại ở nơi chúng bắt đầu. Và khi một tài khoản chạm giới hạn sử dụng, bạn chỉ biết được khi thử dùng.

## ccam làm gì

- `ccam add` tạo thư mục, chạy bước đăng nhập và, nếu bạn muốn, gán cho tài khoản một lệnh tắt do bạn tự đặt tên, chẳng hạn `ccwork`.
- `ccam work` và `ccwork` khởi động Claude Code bằng tài khoản đó và chuyển nguyên vẹn mọi đối số.
- Máy chủ MCP, cài đặt, `CLAUDE.md` và tiện ích mở rộng được dùng chung qua liên kết, theo từng mục và từng tài khoản. Thông tin đăng nhập không bao giờ được dùng chung.
- `ccam usage` hiển thị mức sử dụng 5 giờ và hàng tuần của mọi tài khoản; `ccam next --continue` tiếp tục cuộc hội thoại ở tài khoản kế tiếp.

## Không cản trở công việc

Một công cụ bạn chạy trước mỗi phiên làm việc không nên tốn thêm chi phí nào mà bạn nhận thấy. ccam đọc một tệp cấu hình, kiểm tra mỗi liên kết dùng chung bằng một lời gọi hệ thống duy nhất, rồi tự thay thế mình bằng `claude`. Nó thêm khoảng 4 mili giây và không chạy trong lúc bạn làm việc.

ccam cũng cẩn thận với tệp của bạn. Mọi thay đổi về việc dùng chung đều in ra một kế hoạch trước khi bắt đầu và giữ bản gốc trong một thư mục sao lưu. `ccam doctor` kiểm tra liên kết, quyền truy cập và thông tin đăng nhập, rồi sửa những gì có thể sửa an toàn.

## Dùng thử

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam là phần mềm mã nguồn mở và không chính thức: nó không liên kết với Anthropic và không được Anthropic xác nhận. Việc bạn có được dùng nhiều tài khoản hay không phụ thuộc vào điều khoản của từng tài khoản. Mọi góp ý và bản dịch đều được hoan nghênh trên [GitHub](https://github.com/jhyunwoo/ccam).
