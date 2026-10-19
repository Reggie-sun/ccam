---
title: Làm gì khi Claude Code chạm giới hạn sử dụng
description: "Chờ giới hạn được đặt lại hoặc tiếp tục ngay cùng cuộc hội thoại trên một tài khoản Claude Code khác. Giới hạn 5 giờ và hàng tuần, từng bước một."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: ff765a189332
---

Khi một tài khoản chạm giới hạn sử dụng (thường gọi là rate limit), bạn có hai lựa chọn: chờ đến khi giới hạn được đặt lại, hoặc tiếp tục ngay cùng cuộc hội thoại trên một tài khoản khác của bạn. Với ccam, cách thứ hai chỉ cần một lệnh, `ccam next --auto --continue`, sau khi lịch sử hội thoại đã được dùng chung giữa các tài khoản.

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
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Giờ đây dòng trạng thái hiển thị tài khoản và mức sử dụng 5 giờ của nó trong lúc bạn làm việc, nên việc chạm giới hạn sẽ không còn bất ngờ.

Khi chạm giới hạn, hãy thoát Claude Code (`/exit`), rồi chạy:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Mở cuộc hội thoại trên tài khoản còn nhiều hạn mức nhất:

```sh
ccam next --auto --continue
```

Claude Code khởi động với tài khoản `work` và cuộc hội thoại gần nhất của thư mục hiện tại, bao gồm mọi nội dung đã trao đổi.

## ccam có tự động đổi tài khoản không?

Không. ccam không chạy trong lúc bạn làm việc: nó khởi động Claude Code rồi tự thay thế mình bằng Claude Code. Bạn quyết định khi nào đổi, và `ccam next --auto` chọn tài khoản có mức sử dụng 5 giờ thấp nhất, bỏ qua các tài khoản đã đạt 100% hoặc chưa đăng nhập. Một số công cụ khác thì hoán đổi đăng nhập của phiên đang chạy; [ccam và các công cụ đổi tài khoản khác](/compare/account-switchers) giải thích sự khác biệt.

## Những điều cần lưu ý

- Chạy `ccam next --continue` trong đúng thư mục bạn đã làm việc; hội thoại được lưu theo từng thư mục dự án.
- Tránh chạy hai tài khoản trên cùng một cuộc hội thoại cùng lúc. ccam sẽ cảnh báo nếu bạn thử.
- Tài khoản kia cũng phải dùng chung lịch sử hội thoại, trong cùng một nhóm.
- Tin nhắn đầu tiên trên tài khoản mới có thể tốn nhiều hạn mức hơn bình thường, vì tài khoản đó chưa lưu cache cuộc hội thoại.
- Việc bạn có được đổi tài khoản như vậy hay không phụ thuộc vào điều khoản của từng tài khoản.

Xem [giới hạn sử dụng](/docs/usage-limits) để biết mọi tùy chọn của `ccam usage` và `ccam next`.
