---
title: Bắt đầu nhanh
metaTitle: "Bắt đầu nhanh: thêm và đổi tài khoản Claude Code"
description: Thêm hai tài khoản Claude Code bằng ccam add, đặt cho chúng lệnh tắt như ccwork, chuyển đổi giữa chúng và chọn những gì dùng chung.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

Thêm từng tài khoản bằng `ccam add`, rồi khởi động Claude Code bằng `ccam <name>` hoặc bằng một lệnh tắt do bạn tự đặt tên, chẳng hạn `ccwork`. Hầu hết mọi người chỉ cần như vậy.

Bạn đã dùng Claude Code? Hãy đọc [Bạn đã dùng ~/.claude?](#already-using-claude) ở cuối trang này trước khi bắt đầu.

## 1. Thêm tài khoản đầu tiên

```sh
ccam add
```

Trình hướng dẫn hỏi một cái tên, chẳng hạn `personal` hoặc `work`, những mục nào tài khoản nên dùng chung với các tài khoản khác của bạn, và có tạo lệnh tắt hay không (xem bước 3). Sau đó nó chạy `claude auth login`, lệnh này mở trình duyệt để bạn đăng nhập. Trên máy từ xa, trình duyệt của bạn có thể dừng ở một trang `http://localhost:…/callback` không tải được; hãy dán địa chỉ đó vào terminal và ccam sẽ hoàn tất việc đăng nhập. Khi đăng nhập xong, bạn sẽ thấy:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Mang theo cài đặt hiện có

Nếu bạn đã từng dùng Claude Code, cài đặt, `CLAUDE.md` và tiện ích mở rộng của bạn nằm trong `~/.claude`, còn máy chủ MCP của bạn nằm trong `~/.claude.json`. Hãy sao chép chúng vào nhóm `default` dùng chung để mọi tài khoản đều có:

```sh
ccam migrate
```

ccam cho bạn xem những gì nó sẽ sao chép và hỏi trước. Nó chỉ đọc `~/.claude`, không bao giờ sao chép thông tin đăng nhập, và không đưa lịch sử hội thoại vào trừ khi bạn thêm `--history`. Xem [sao chép thiết lập hiện có](/docs/sharing#copy-an-existing-setup) để biết chi tiết.

## 3. Thêm tài khoản thứ hai

```sh
ccam add work
```

Trình hướng dẫn bỏ qua bước hỏi tên, vì bạn đã nhập tên trên dòng lệnh, và hiển thị các nhóm để tham gia, với `default` đã được đánh dấu. Lần này, hãy đặt cho tài khoản một lệnh tắt riêng, chẳng hạn `ccwork`, để bạn có thể khởi động nó chỉ bằng một từ. Sau phần chọn nhóm, trình hướng dẫn hỏi:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Nhấn `y`, rồi gõ một cái tên, hoặc nhấn Enter để dùng tên được gợi ý:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Sau đó đăng nhập bằng tài khoản còn lại trong trình duyệt. ccam xác nhận tài khoản và lệnh tắt của nó:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Chuyển đổi tài khoản

`ccam work` khởi động Claude Code bằng tài khoản công việc. Nếu bạn đã bỏ qua câu hỏi về lệnh tắt, bạn có thể thêm lệnh tắt bất cứ lúc nào:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Giờ đây, mỗi lệnh dưới đây đều khởi động Claude Code bằng tài khoản công việc (`2` là số của nó trong `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Chạy `ccam` một mình sẽ mở [bảng điều khiển](/features/dashboard): di chuyển tới một tài khoản bằng các phím mũi tên hoặc số của nó, rồi nhấn Enter.

Mọi thứ sau lệnh tắt hoặc tên tài khoản được chuyển nguyên vẹn cho Claude Code, nên bạn vẫn dùng các tùy chọn quen thuộc:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Chọn những gì dùng chung

Theo mặc định, máy chủ MCP, `settings.json`, `CLAUDE.md` toàn cục và tiện ích mở rộng (agent, lệnh, skill và plugin) được dùng chung giữa các tài khoản qua nhóm `default`. Lịch sử hội thoại thì không. Hãy thêm nó vào nhóm nếu bạn muốn tiếp tục hội thoại qua lại giữa các tài khoản:

```sh
ccam group set default --add history
```

ccam cho bạn xem những gì nó sẽ di chuyển cho từng tài khoản trước khi thay đổi bất cứ thứ gì. Xem [dùng chung bằng nhóm](/features/share) để biết chi tiết.

## 6. Khi một tài khoản chạm giới hạn

```sh
ccam usage
ccam next --continue
```

`ccam usage` hiển thị mức sử dụng 5 giờ và hàng tuần của mọi tài khoản. `ccam next --continue` mở cùng cuộc hội thoại ở tài khoản kế tiếp. Trong script, `ccam --fallback work -p "…"` tự đổi tài khoản. Xem [giới hạn sử dụng](/docs/usage-limits).

**Seat của gói Team và Enterprise.** Việc tự động chọn tài khoản (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển và, khi bạn bật nó, [`auto_switch`](/features/next#auto-switch)) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise, điều này có thể vi phạm điều khoản của Anthropic. Hãy đọc [Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) và xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó.

## Bạn đã dùng ~/.claude?

Hãy đăng ký thư mục đó thay vì đăng nhập lại:

```sh
ccam import
```

Thư mục vẫn nằm nguyên tại chỗ, và lệnh `claude` thông thường vẫn hoạt động như trước.
