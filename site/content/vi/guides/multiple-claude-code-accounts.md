---
title: Dùng nhiều tài khoản Claude Code trên một máy tính
description: "Chạy song song nhiều tài khoản Claude Code với CLAUDE_CONFIG_DIR, thủ công bằng alias của shell hoặc tự động bằng ccam, và đổi tài khoản khi cần."
updated: 2026-10-24
order: 1
translation: machine
sourceHash: "5463dec2f6eb"
---

Hãy cho mỗi tài khoản một thư mục cấu hình riêng. Claude Code đọc thư mục này từ biến môi trường `CLAUDE_CONFIG_DIR`, và mỗi thư mục chứa thông tin đăng nhập, cài đặt và lịch sử riêng, nên mọi tài khoản đều giữ trạng thái đăng nhập và bạn có thể chạy chúng cùng lúc. Bạn có thể tự thiết lập bằng một alias của shell, hoặc để ccam làm việc đó.

## Làm sao để tự thiết lập thủ công?

Đây là cách mà [tài liệu của Claude Code mô tả](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Thêm một dòng cho mỗi tài khoản bổ sung vào `~/.bashrc` hoặc `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Mở một terminal mới và chạy `claude-work`. Claude Code sẽ yêu cầu bạn đăng nhập, vì thư mục này còn mới. Từ đó trở đi, `claude` là tài khoản thứ nhất và `claude-work` là tài khoản thứ hai. Mỗi tài khoản cần một alias riêng, với thư mục riêng.

Trên Windows không có WSL, một hàm PowerShell trong profile của bạn (`notepad $PROFILE`) cũng làm được điều tương tự:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Cách dùng alias còn thiếu những gì?

Cách này hoạt động tốt với hai tài khoản. Những điều cần biết:

- Mỗi thư mục bắt đầu trống: không có cài đặt, không có `CLAUDE.md`, không có máy chủ MCP, không có agent hay skill. Bạn tự sao chép hoặc liên kết những gì cần, và tự giữ các bản sao đồng bộ với nhau.
- Alias chỉ tồn tại trong shell tương tác. Script, Makefile và các chương trình khác không thấy chúng.
- Máy chủ MCP nằm trong `.claude.json`, cùng với thông tin đăng nhập, nên không thể đơn giản liên kết chúng giữa các thư mục.
- Các cuộc hội thoại nằm lại ở tài khoản nơi chúng bắt đầu, nên bạn không thể tiếp tục một cuộc hội thoại ở tài khoản khác.
- Không có gì hiển thị mức sử dụng của tất cả tài khoản cùng lúc.

## Có thể đổi tài khoản bằng /login không?

Có, nếu bạn chỉ thỉnh thoảng mới đổi. Gõ `/login` trong Claude Code sẽ [đổi tài khoản](https://code.claude.com/docs/en/quickstart) của thư mục bạn đang dùng. Mỗi lúc chỉ có một tài khoản được đăng nhập, nên bạn phải đăng nhập lại mỗi lần đổi, và không thể chạy song song hai tài khoản. Dùng các thư mục riêng sẽ tránh được cả hai vấn đề này.

## Làm sao để thiết lập bằng ccam?

ccam tự động hóa cùng ý tưởng đó và bổ sung những phần còn thiếu:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Giờ đây `ccam personal` và `ccam work` sẽ khởi động hai tài khoản, và `ccwork`, lệnh tắt bạn đặt cho tài khoản công việc, cũng vậy. Vì chúng là chương trình thực sự chứ không phải alias, chúng cũng hoạt động trong script. Chạy riêng `ccam` để chọn tài khoản từ danh sách, và `ccam use work` để đặt một tài khoản làm mặc định.

Cài đặt, `CLAUDE.md`, máy chủ MCP và tiện ích mở rộng được dùng chung ngay từ đầu. `ccam next --continue` chuyển một cuộc hội thoại sang tài khoản kia, và lịch sử hội thoại dùng chung, thứ bạn có thể bật, giữ nó là một cuộc hội thoại duy nhất. `ccam usage` hiển thị mức sử dụng 5 giờ và hàng tuần của mọi tài khoản cùng lúc.

Nếu bạn đã dùng `~/.claude`, hãy chạy `ccam import` thay vì thêm lại tài khoản đó, hoặc chạy `ccam migrate` để sao chép cài đặt, máy chủ MCP và tiện ích mở rộng của nó vào phần cài đặt dùng chung.

## Làm sao biết tài khoản nào đang chạy?

Trong Claude Code, `/status` cho biết tài khoản bạn đang đăng nhập. Với ccam, hãy chạy `ccam statusline install` một lần, và dòng trạng thái sẽ luôn hiển thị tên tài khoản cùng mức sử dụng của nó, ví dụ `[work] · Opus · 5h 24% · 7d 41%`.

## Cách này có dùng được trên Windows và WSL không?

Claude Code chạy [trực tiếp trên Windows và bên trong WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` hoạt động ở cả hai, nên trên Windows gốc chỉ cần alias hoặc hàm PowerShell ở trên là đủ. ccam cũng chạy trực tiếp trên Windows: cài nó trong PowerShell bằng `irm https://ccam.dev/install.ps1 | iex`, hoặc bên trong WSL bằng lệnh dành cho Linux. Lệnh tắt trở thành `ccwork.exe`, và việc dùng chung các tệp như `settings.json` cần bật Developer Mode ([Cài đặt ccam](/docs/install)).

## Khắc phục sự cố

**Claude Code bắt tôi đăng nhập mỗi lần.** Lần khởi động đầu tiên trong một thư mục mới luôn yêu cầu đăng nhập. Nếu nó cứ tiếp tục yêu cầu, phiên đăng nhập đã hết hạn hoặc bị xóa: hãy chạy `ccam login work` (hoặc chạy alias một lần và đăng nhập).

**Một máy chủ MCP chỉ xuất hiện ở một tài khoản.** Máy chủ được thêm trong một phiên bằng `claude mcp add` thuộc về tài khoản đó, và nếu không có `-s user` thì chỉ thuộc về thư mục đó. Chuyển nó sang mọi tài khoản bằng `ccam mcp import work <server>`.

**Cài đặt khác nhau giữa các tài khoản.** `ccam share work` in ra những gì tài khoản đó dùng chung và qua nhóm nào; `ccam group join default work` khiến nó dùng chung những gì `default` dùng chung, sau khi cho bạn xem những gì sẽ thay đổi. `ccam doctor --fix` sửa các liên kết bị hỏng và giữ bản sao lưu của mọi tệp mà nó thay đổi.

**Hai tài khoản mở cùng một cuộc hội thoại.** Tránh làm việc trên một cuộc hội thoại bằng hai tài khoản cùng lúc. ccam sẽ cảnh báo nếu bạn thử.

## Nên chọn cách nào?

Với hai tài khoản và không dùng chung gì, một alias là đủ. Khi bạn muốn dùng chung cài đặt và máy chủ MCP, có nhiều tài khoản hơn, xem nhanh mức sử dụng hoặc tiếp tục hội thoại, ccam giúp bạn khỏi phải tự quản lý những việc đó. Phần [so sánh với alias](/compare/aliases) liệt kê song song các điểm khác biệt, còn [ccam và các công cụ đổi tài khoản khác](/compare/account-switchers) nói về những công cụ hoán đổi một lần đăng nhập.

Việc bạn có được dùng nhiều tài khoản hay không phụ thuộc vào điều khoản của từng tài khoản. Hãy kiểm tra trước khi dùng.

**Seat của gói Team và Enterprise.** Lựa chọn tự động (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Vì vậy, dùng lựa chọn tự động giữa các seat có thể vi phạm điều khoản của Anthropic. Hãy xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó. Các gói cá nhân (Pro, Max) được đề cập trong ghi chú ở trên.
