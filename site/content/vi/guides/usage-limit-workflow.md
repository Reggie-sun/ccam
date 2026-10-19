---
title: Làm gì khi Claude Code chạm giới hạn sử dụng
description: "Kiểm tra mỗi tài khoản Claude Code còn bao nhiêu hạn mức và đổi tài khoản Claude Code để tiếp tục cùng cuộc hội thoại khi một tài khoản chạm giới hạn."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: db7ed9b7bcfc
---

Khi một tài khoản chạm giới hạn, hãy tiếp tục cùng cuộc hội thoại ở một tài khoản khác của bạn: bật dùng chung lịch sử hội thoại một lần, rồi chạy `ccam next --continue`. Nếu bạn chỉ có một tài khoản, hãy chờ đến thời điểm đặt lại mà Claude Code hiển thị.

## Trước khi điều đó xảy ra

Thiết lập các tài khoản để việc chuyển đổi chỉ cần một lệnh:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Giờ đây dòng trạng thái hiển thị tài khoản và mức sử dụng 5 giờ của nó trong lúc bạn làm việc, nên việc chạm giới hạn sẽ không làm bạn bất ngờ.

## Khi đã chạm giới hạn

Thoát Claude Code (`/exit`), sau đó:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

Mở cuộc hội thoại ở tài khoản còn nhiều hạn mức nhất:

```sh
ccam next --auto --continue
```

Claude Code khởi động dưới tài khoản `work` với cuộc hội thoại gần nhất của thư mục hiện tại, bao gồm mọi thứ đã trao đổi cho đến lúc đó.

## Những điều cần lưu ý

- Chạy `ccam next --continue` trong cùng thư mục bạn đã làm việc; hội thoại được lưu theo từng thư mục dự án.
- Tránh chạy hai tài khoản trên cùng một cuộc hội thoại cùng lúc. ccam sẽ cảnh báo nếu bạn thử làm vậy.
- Tài khoản kia cũng phải dùng chung lịch sử hội thoại, trong cùng một nhóm.
- Việc bạn có được chuyển đổi giữa các tài khoản theo cách này hay không phụ thuộc vào điều khoản của từng tài khoản.

Xem [giới hạn sử dụng](/docs/usage-limits) để biết mọi tùy chọn của `ccam usage` và `ccam next`.
