---
title: Câu hỏi thường gặp
metaTitle: "ccam FAQ: dùng nhiều tài khoản Claude Code"
description: "Giải đáp về ccam, tài khoản Claude Code, giới hạn sử dụng, việc dùng chung, quyền riêng tư và các hệ điều hành được hỗ trợ."
updated: 2026-10-20
order: 6
translation: machine
sourceHash: 80a5475ccfee
---

## ccam có phải là công cụ chính thức của Anthropic không?

Không. ccam là dự án mã nguồn mở không chính thức, không liên kết với Anthropic và không được Anthropic xác nhận. Nó dùng thiết lập `CLAUDE_CONFIG_DIR` mà Claude Code [mô tả trong tài liệu](https://code.claude.com/docs/en/env-vars) để chạy song song nhiều tài khoản.

## Tôi có được phép dùng nhiều tài khoản Claude Code không?

Điều đó phụ thuộc vào điều khoản của từng tài khoản và, với tài khoản công việc, vào chính sách của tổ chức bạn. Hãy kiểm tra trước khi dùng nhiều tài khoản. ccam chỉ giúp việc chuyển đổi giữa các tài khoản bạn đã có trở nên dễ dàng hơn.

## ccam có giúp vượt qua giới hạn sử dụng không?

Không. Mỗi tài khoản giữ giới hạn riêng của nó. ccam hiển thị các giới hạn đó và cho phép bạn tiếp tục trên một tài khoản khác của bạn khi một tài khoản hết hạn mức.

## Làm sao để dùng hai tài khoản Claude Code trên một máy tính?

Hãy cho mỗi tài khoản một thư mục cấu hình riêng bằng `CLAUDE_CONFIG_DIR`. ccam làm việc này thay bạn: chạy `ccam add` cho từng tài khoản, rồi khởi động chúng bằng `ccam personal` và `ccam work`, hoặc bằng các lệnh tắt do bạn tự đặt tên, chẳng hạn `ccwork`. [Hướng dẫn](/guides/multiple-claude-code-accounts) cũng chỉ cách làm thủ công.

## Làm sao để tiếp tục một cuộc hội thoại trên tài khoản khác?

Bật dùng chung lịch sử hội thoại giữa hai tài khoản bằng `ccam share <name> --history on`, rồi chạy `ccam next --continue`. Claude Code sẽ mở cuộc hội thoại gần nhất của thư mục hiện tại trong tài khoản kế tiếp.

## ccam có xem hoặc gửi thông tin đăng nhập của tôi đi không?

ccam không bao giờ dùng chung thông tin đăng nhập Claude của bạn giữa các tài khoản. Nó chỉ đọc access token của một tài khoản để hỏi endpoint mức sử dụng của Anthropic về mức sử dụng của tài khoản đó, khi bạn chạy `ccam usage`, `ccam next --auto` hoặc `ccam doctor`. `ccam migrate` sao chép thông tin đăng nhập đã lưu của các máy chủ MCP từ `~/.claude` vào tệp đăng nhập riêng của các tài khoản của bạn trên cùng một máy tính. ccam không có telemetry và không gửi gì đến bất kỳ nơi nào khác.

## Làm sao để đăng nhập trên máy chủ từ xa?

Chạy `ccam add` hoặc `ccam login <name>` trong phiên SSH như bình thường. Claude Code chờ việc đăng nhập trên máy chủ, nên trình duyệt trên máy tính của bạn dừng ở một trang `http://localhost:<port>/callback?code=…` không tải được. Hãy sao chép địa chỉ đó từ trình duyệt và dán vào terminal: ccam chuyển nó cho Claude Code trên máy chủ, và Claude Code hoàn tất việc đăng nhập. Địa chỉ này chứa một mã dùng một lần, vì vậy đừng chia sẻ nó.

## Các tài khoản của tôi được lưu ở đâu?

Trong `~/.local/share/ccam/profiles/<name>/`, mỗi tài khoản một thư mục. Trên macOS, Claude Code giữ chính thông tin đăng nhập trong Chuỗi khóa (Keychain), dưới một tên được tạo từ thư mục đó.

## Thư mục ~/.claude hiện có của tôi sẽ ra sao?

Không có gì thay đổi, trừ khi bạn chạy `ccam import` hoặc `ccam migrate`. Lệnh import đăng ký thư mục ngay tại chỗ, nên lệnh `claude` thông thường và `ccam <name>` tiếp tục dùng cùng một tài khoản và cùng cài đặt. `ccam migrate` sao chép cài đặt, `CLAUDE.md`, tiện ích mở rộng và máy chủ MCP của thư mục đó vào phần cài đặt dùng chung, và chỉ đọc thư mục.

## ccam hỗ trợ những hệ điều hành nào?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine và các bản phân phối khác), macOS và Windows 10 và 11, chạy trực tiếp hoặc bên trong WSL. Trên Windows, việc dùng chung các tệp như `settings.json` cần bật Developer Mode; thư mục dùng chung thì hoạt động mà không cần chế độ này ([dùng chung trên Windows](/features/share#on-windows)).

## Vì sao một máy chủ MCP chỉ xuất hiện ở một tài khoản?

Máy chủ được thêm trong một phiên bằng `claude mcp add` chỉ được lưu cho tài khoản đó. ccam sẽ báo cho bạn vào lần tới bạn khởi động tài khoản; hãy chuyển máy chủ sang nhóm dùng chung bằng `ccam mcp import <account> <server>`.

## ccam có làm Claude Code chậm đi không?

Khởi động Claude Code thông qua ccam tốn thêm khoảng 4 mili giây, đo bằng `scripts/perf-check.sh` trên một máy tính xách tay. ccam đọc một tệp nhỏ, kiểm tra các liên kết dùng chung rồi tự thay thế mình bằng `claude`, nên nó không chạy trong lúc bạn làm việc. Windows không thể thay thế một tiến trình, nên ở đó ccam khởi động Claude Code và chờ nó, ngoài ra không làm gì khác.

## Làm sao để cập nhật ccam?

Chạy `ccam update`. Lệnh này tải về bản phát hành mới nhất, kiểm tra mã băm rồi thay thế tệp `ccam`; `ccam update --check` chỉ cho biết có bản phát hành mới hơn hay không. ccam không bao giờ tự cập nhật. Nếu bạn đã cài gói `.deb`, hãy dùng apt để cài gói `.deb` mới. Xem [cập nhật](/docs/install#update).

## Làm sao để gỡ cài đặt ccam?

Xóa lệnh tắt của từng tài khoản bằng `ccam shortcut <account> --remove`, xóa `ccam` khỏi `~/.local/bin`, sau đó xóa `~/.config/ccam`, `~/.local/share/ccam` và `~/.local/state/ccam` nếu bạn không còn cần các tài khoản. Xem [cài đặt](/docs/install#uninstall).
