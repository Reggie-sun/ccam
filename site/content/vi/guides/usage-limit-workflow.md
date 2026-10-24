---
title: Làm gì khi Claude Code chạm giới hạn sử dụng
description: "Chờ giới hạn được đặt lại hoặc tiếp tục ngay cùng cuộc hội thoại trên một tài khoản Claude Code khác. Giới hạn 5 giờ và hàng tuần, từng bước một."
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "049a46684c7b"
---

Khi một tài khoản chạm giới hạn sử dụng (thường gọi là rate limit), bạn có hai lựa chọn: chờ đến khi giới hạn được đặt lại, hoặc tiếp tục ngay cùng cuộc hội thoại trên một tài khoản khác của bạn. Với ccam, cách thứ hai chỉ cần một lệnh, `ccam next --auto --continue`.

## Claude Code có những giới hạn nào?

Với gói Pro và Max, Claude Code có [giới hạn phiên 5 giờ và giới hạn hàng tuần](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Cả hai đều dùng chung với Claude trên web và trong ứng dụng, và mức sử dụng được tính vào cả hai cùng lúc. Khi một trong hai hết, Claude Code [hiển thị thông báo](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) như "You've hit your session limit · resets 3:45pm". Đổi model không giúp được gì, vì giới hạn áp dụng cho mọi model.

## Cách 1: chờ giới hạn được đặt lại

Nếu bạn chỉ có một tài khoản, hãy chờ. Các phiên bản Claude Code gần đây [chờ ngay trong phiên đang mở](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) và tự tiếp tục công việc khi giới hạn được đặt lại. `/usage` trong Claude Code cho biết bạn đã dùng bao nhiêu trong giới hạn của gói.

## Cách 2: tiếp tục ngay trên tài khoản khác

Nếu có tài khoản thứ hai, bạn có thể tiếp tục ngay với cùng cuộc hội thoại. Thiết lập một lần:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Giờ đây dòng trạng thái hiển thị tài khoản và mức sử dụng 5 giờ của nó trong lúc bạn làm việc, nên việc chạm giới hạn sẽ không còn bất ngờ.

Khi chạm giới hạn, hãy thoát Claude Code (`/exit`), rồi chạy:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Mở cuộc hội thoại trên tài khoản còn nhiều hạn mức nhất:

```sh
ccam next --auto --continue
```

Claude Code khởi động với tài khoản `work` và cuộc hội thoại gần nhất của thư mục hiện tại, bao gồm mọi nội dung đã trao đổi. Dùng chung lịch sử (`ccam group set default --add history` ở trên) giữ nó là một cuộc hội thoại duy nhất; nếu không, Claude Code sao chép cuộc hội thoại vào `work`, và ccam cho bạn biết câu trả lời đầu tiên đó phải đọc lại bao nhiêu token.

Với [plugin ccam](/docs/install#claude-code-plugin) trong Claude Code, một phiên dừng lại vì chạm giới hạn sẽ hiện một thông báo và ghi lại khoảng đã đầy, nên `ccam next --auto` không chọn lại tài khoản đó trước khi đặt lại, và `/ccam:switch` (hoặc `/ccam:switch work`) ghi chú cuộc hội thoại cho `ccam next --continue` mà không để lại điều gì phải đoán.

**Seat của gói Team và Enterprise.** Lựa chọn tự động (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Vì vậy, dùng lựa chọn tự động giữa các seat có thể vi phạm điều khoản của Anthropic. Hãy xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó. Các gói cá nhân (Pro, Max) được đề cập trong [ghi chú về điều khoản ở phần Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## ccam có tự động đổi tài khoản không?

Không. ccam không chạy trong lúc bạn làm việc: nó khởi động Claude Code rồi tự thay thế mình bằng Claude Code. Bạn quyết định khi nào đổi, và `ccam next --auto` chọn tài khoản có mức sử dụng 5 giờ thấp nhất, bỏ qua các tài khoản đã đạt 100%, chưa đăng nhập, hoặc thuộc tổ chức khác. Với `relay = true` trong `config.toml`, thay vào đó ccam chờ bên cạnh phiên, nên `/ccam:switch` tự chuyển cuộc hội thoại đi tiếp, giữ lại các tùy chọn, chế độ quyền và mức effort của phiên. Một số công cụ khác thì hoán đổi đăng nhập của phiên đang chạy; [ccam và các công cụ đổi tài khoản khác](/compare/account-switchers) giải thích sự khác biệt.

## Những điều cần lưu ý

- Chạy `ccam next --continue` trong đúng thư mục bạn đã làm việc; hội thoại được lưu theo từng thư mục dự án.
- Tránh chạy hai tài khoản trên cùng một cuộc hội thoại cùng lúc. ccam sẽ cảnh báo nếu bạn thử.
- Không dùng chung lịch sử, cuộc hội thoại được sao chép sang tài khoản kia và hai bản sao đi theo hướng riêng; các checkpoint mà `/rewind` dùng để khôi phục mã vẫn ở lại với tài khoản đầu tiên.
- Tin nhắn đầu tiên trên tài khoản mới có thể tốn nhiều hạn mức hơn bình thường, vì tài khoản đó chưa lưu cache cuộc hội thoại.
- Việc bạn có được đổi tài khoản như vậy hay không phụ thuộc vào điều khoản của từng tài khoản.

Xem [giới hạn sử dụng](/docs/usage-limits) để biết mọi tùy chọn của `ccam usage` và `ccam next`.
