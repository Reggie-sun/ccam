---
title: Câu hỏi thường gặp
metaTitle: "ccam FAQ: dùng nhiều tài khoản Claude Code"
description: "Giải đáp về ccam, tài khoản Claude Code, giới hạn sử dụng, việc dùng chung, quyền riêng tư và các hệ điều hành được hỗ trợ."
updated: 2026-10-24
order: 6
translation: machine
sourceHash: "c4e3b89bb4ea"
---

## ccam có phải là công cụ chính thức của Anthropic không?

Không. ccam là dự án mã nguồn mở không chính thức, không liên kết với Anthropic và không được Anthropic xác nhận. Nó dùng thiết lập `CLAUDE_CONFIG_DIR` mà Claude Code [mô tả trong tài liệu](https://code.claude.com/docs/en/env-vars) để chạy song song nhiều tài khoản.

## Tôi có được phép dùng nhiều tài khoản Claude Code không?

Điều đó phụ thuộc vào điều khoản của từng tài khoản và, với tài khoản công việc, vào chính sách của tổ chức bạn. Hãy kiểm tra trước khi dùng nhiều tài khoản. ccam chỉ giúp việc chuyển đổi giữa các tài khoản bạn đã có trở nên dễ dàng hơn.

**Seat của gói Team và Enterprise.** Lựa chọn tự động (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Vì vậy, dùng lựa chọn tự động giữa các seat có thể vi phạm điều khoản của Anthropic. Hãy xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó. Các gói cá nhân (Pro, Max) được đề cập trong đoạn ở trên.

## ccam có giúp vượt qua giới hạn sử dụng không?

Không. Mỗi tài khoản giữ giới hạn riêng của nó. ccam hiển thị các giới hạn đó và cho phép bạn tiếp tục trên một tài khoản khác của bạn khi một tài khoản hết hạn mức.

## Làm sao để dùng hai tài khoản Claude Code trên một máy tính?

Hãy cho mỗi tài khoản một thư mục cấu hình riêng bằng `CLAUDE_CONFIG_DIR`. ccam làm việc này thay bạn: chạy `ccam add` cho từng tài khoản, rồi khởi động chúng bằng `ccam personal` và `ccam work`, hoặc bằng các lệnh tắt do bạn tự đặt tên, chẳng hạn `ccwork`. [Hướng dẫn](/guides/multiple-claude-code-accounts) cũng chỉ cách làm thủ công.

## Làm sao để tiếp tục một cuộc hội thoại trên tài khoản khác?

Chạy `ccam next --continue` trong thư mục: lệnh này chuyển cuộc hội thoại gần nhất của thư mục sang tài khoản kế tiếp (`--auto` để chọn tài khoản còn nhiều hạn mức nhất). Bên trong Claude Code, `/ccam:switch` từ [plugin ccam](/docs/install#claude-code-plugin) làm điều tương tự cho phiên bạn đang ở. Khi hai tài khoản dùng chung lịch sử hội thoại (`ccam group set default --add history`, hoặc một nhóm riêng của chúng: `ccam group add pair --items history --members work,side`), tài khoản mới mở chính cuộc hội thoại đó; nếu không, Claude Code sao chép nó vào lịch sử của tài khoản mới.

## ccam có xem hoặc gửi thông tin đăng nhập của tôi đi không?

ccam không bao giờ dùng chung hay sao chép thông tin đăng nhập Claude của bạn giữa các tài khoản, và không bao giờ sao chép các lần đăng nhập của máy chủ MCP. Nó chỉ đọc access token của một tài khoản để hỏi endpoint mức sử dụng của Anthropic về mức sử dụng của tài khoản đó, khi bạn chạy `ccam usage`, `ccam next --auto` hoặc `ccam doctor`, và chỉ gia hạn token đã hết hạn bằng `ccam usage --refresh`; với `usage_source = "claude"` nó hoàn toàn không đọc token. `ccam ext` và `ccam doctor` đọc tên các lần đăng nhập MCP của một tài khoản, không bao giờ đọc token của chúng. ccam không có telemetry. Ngoài các lượt tra cứu đó, yêu cầu duy nhất nó tự gửi là câu hỏi của bảng điều khiển tới GitHub về bản phát hành ccam mới nhất, và yêu cầu này không mang thông tin gì về bạn hay các tài khoản của bạn.

## Làm sao để đăng nhập trên máy chủ từ xa?

Chạy `ccam add` hoặc `ccam login <name>` trong phiên SSH như bình thường. Claude Code chờ việc đăng nhập trên máy chủ, nên trình duyệt trên máy tính của bạn dừng ở một trang `http://localhost:<port>/callback?code=…` không tải được. Hãy sao chép địa chỉ đó từ trình duyệt và dán vào terminal: ccam chuyển nó cho Claude Code trên máy chủ, và Claude Code hoàn tất việc đăng nhập. Địa chỉ này chứa một mã dùng một lần, vì vậy đừng chia sẻ nó.

## Các tài khoản của tôi được lưu ở đâu?

Trong `~/.local/share/ccam/profiles/<name>/`, mỗi tài khoản một thư mục. Trên macOS, Claude Code giữ chính thông tin đăng nhập trong Chuỗi khóa (Keychain), dưới một tên được tạo từ thư mục đó, nên `ccam rename` giữ lại tên cũ của thư mục ở đó.

## Thư mục ~/.claude hiện có của tôi sẽ ra sao?

Không có gì thay đổi, trừ khi bạn chạy `ccam import` hoặc `ccam migrate`. Lệnh import đăng ký thư mục ngay tại chỗ, nên lệnh `claude` thông thường và `ccam <name>` tiếp tục dùng cùng một tài khoản và cùng cài đặt. `ccam migrate` sao chép cài đặt, `CLAUDE.md`, tiện ích mở rộng và máy chủ MCP của thư mục đó vào phần cài đặt dùng chung, và chỉ đọc thư mục.

## ccam hỗ trợ những hệ điều hành nào?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine và các bản phân phối khác), macOS và Windows 10 và 11, chạy trực tiếp hoặc bên trong WSL. Trên Windows, việc dùng chung các tệp như `settings.json` cần bật Developer Mode; thư mục dùng chung thì hoạt động mà không cần chế độ này, và khi không có nó, lần `ccam add` đầu tiên tạo nhóm `default` chỉ với máy chủ MCP và tiện ích mở rộng ([dùng chung trên Windows](/features/share#on-windows)).

## Vì sao một máy chủ MCP chỉ xuất hiện ở một tài khoản?

Máy chủ được thêm trong một phiên bằng `claude mcp add` chỉ được lưu cho tài khoản đó, và nếu không có `-s user` thì chỉ cho thư mục hiện tại. ccam sẽ báo cho bạn vào lần tới bạn khởi động tài khoản, và `ccam mcp list` cùng `ccam ext` hiển thị chúng; hãy chuyển máy chủ sang nhóm dùng chung bằng `ccam mcp import <account> <server>`.

## ccam có làm Claude Code chậm đi không?

Khởi động Claude Code thông qua ccam tốn thêm khoảng 4 mili giây, đo bằng `scripts/perf-check.sh` trên một máy tính xách tay. ccam đọc một tệp nhỏ, kiểm tra các liên kết dùng chung rồi tự thay thế mình bằng `claude`, nên nó không chạy trong lúc bạn làm việc. Windows không thể thay thế một tiến trình, nên ở đó ccam khởi động Claude Code và chờ nó, ngoài ra không làm gì khác.

## Làm sao để cập nhật ccam?

Chạy `ccam update`. Lệnh này tải về bản phát hành mới nhất, kiểm tra mã băm rồi thay thế tệp `ccam`; `ccam update --check` chỉ cho biết có bản phát hành mới hơn hay không. Bảng điều khiển cũng báo cho bạn khi có bản phát hành mới hơn, và `U` sẽ cài nó sau khi hỏi; ccam không bao giờ tự cập nhật. Nếu bạn đã cài gói `.deb`, hãy dùng apt để cài gói `.deb` mới. Xem [cập nhật](/docs/install#update).

## Làm sao để gỡ cài đặt ccam?

Xóa lệnh tắt của từng tài khoản bằng `ccam shortcut <account> --remove`. Nếu bạn đã nhập `~/.claude` và nó nằm trong một nhóm, hãy chạy `ccam remove <name>` cho nó trước, để nó nhận lại một bản sao của những gì nó dùng chung. Sau đó xóa `ccam` khỏi `~/.local/bin`, và xóa `~/.config/ccam`, `~/.local/share/ccam` và `~/.local/state/ccam` nếu bạn không còn cần các tài khoản. Xem [cài đặt](/docs/install#uninstall).
