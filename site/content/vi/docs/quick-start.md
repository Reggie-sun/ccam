---
title: Bắt đầu nhanh
metaTitle: "Bắt đầu nhanh: thêm và đổi tài khoản Claude Code"
description: Thêm hai tài khoản Claude Code bằng ccam add, đặt cho chúng lệnh tắt như ccwork, chuyển đổi giữa chúng và chọn những gì dùng chung.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 745a9904e32a
---

Thêm từng tài khoản bằng `ccam add`, rồi khởi động Claude Code bằng `ccam <name>` hoặc bằng một lệnh tắt do bạn tự đặt tên, chẳng hạn `ccwork`. Hầu hết mọi người chỉ cần như vậy.

## 1. Thêm tài khoản đầu tiên

```sh
ccam add
```

Trình hướng dẫn hỏi một cái tên, chẳng hạn `personal` hoặc `work`, những mục nào tài khoản nên dùng chung với các tài khoản khác của bạn, và có tạo lệnh tắt hay không (xem bước 4). Sau đó nó chạy `claude auth login`, lệnh này mở trình duyệt để bạn đăng nhập. Khi hoàn tất, bạn sẽ thấy:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Mang theo cài đặt hiện có

Nếu bạn đã từng dùng Claude Code, cài đặt, `CLAUDE.md`, tiện ích mở rộng và máy chủ MCP của bạn nằm trong `~/.claude`. Hãy sao chép chúng vào phần cài đặt dùng chung để mọi tài khoản đều có:

```sh
ccam migrate
```

ccam cho bạn xem những gì nó sẽ sao chép và hỏi trước. Nó chỉ đọc `~/.claude`, không bao giờ sao chép thông tin đăng nhập, và không đưa lịch sử hội thoại vào trừ khi bạn thêm `--history`. Xem [sao chép thiết lập hiện có](/docs/sharing#copy-an-existing-setup) để biết chi tiết.

## 3. Thêm tài khoản thứ hai

```sh
ccam add work
```

Đăng nhập bằng tài khoản còn lại trong trình duyệt.

## 4. Chuyển đổi tài khoản

`ccam work` khởi động Claude Code bằng tài khoản công việc. Nếu muốn ngắn hơn, hãy đặt cho tài khoản một lệnh tắt riêng, chẳng hạn `ccwork`. Sau phần chọn mục dùng chung, trình hướng dẫn hỏi:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Nhấn `y`, rồi gõ một cái tên, hoặc nhấn Enter để dùng tên được gợi ý:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Sau khi đăng nhập, ccam xác nhận:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Nếu bạn đã bỏ qua câu hỏi này, bạn có thể thêm lệnh tắt bất cứ lúc nào:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Giờ đây, bất kỳ lệnh nào dưới đây cũng khởi động Claude Code bằng tài khoản công việc:

```sh
ccwork
ccam work
ccam 2
```

Chạy `ccam` một mình sẽ mở trình chọn: dùng các phím mũi tên và Enter, hoặc nhấn số của tài khoản.

Mọi thứ sau lệnh tắt hoặc tên tài khoản được chuyển nguyên vẹn cho Claude Code, nên bạn vẫn dùng các tùy chọn quen thuộc:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Chọn những gì dùng chung

Theo mặc định, máy chủ MCP, `settings.json`, `CLAUDE.md` toàn cục, cùng agent, lệnh, skill và plugin được dùng chung giữa các tài khoản. Lịch sử hội thoại thì không. Hãy bật nó cho cả hai tài khoản nếu bạn muốn tiếp tục hội thoại qua lại giữa chúng:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam cho bạn xem những gì nó sẽ di chuyển trước khi thay đổi bất cứ thứ gì. Xem [dùng chung](/docs/sharing) để biết chi tiết.

## 6. Khi một tài khoản chạm giới hạn

```sh
ccam usage
ccam next --continue
```

`ccam usage` hiển thị mức sử dụng 5 giờ và hàng tuần của mọi tài khoản. `ccam next --continue` mở cùng cuộc hội thoại ở tài khoản kế tiếp. Xem [giới hạn sử dụng](/docs/usage-limits).

## Bạn đã dùng ~/.claude?

Hãy đăng ký thư mục đó thay vì đăng nhập lại:

```sh
ccam import
```

Thư mục vẫn nằm nguyên tại chỗ, và lệnh `claude` thông thường vẫn hoạt động như trước.
