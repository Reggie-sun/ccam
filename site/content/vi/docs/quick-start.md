---
title: Bắt đầu nhanh
description: "Thêm hai tài khoản Claude Code bằng ccam add, chuyển đổi giữa chúng bằng claude1, claude2 hoặc ccs, và chọn những gì chúng dùng chung."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 9462b80d61bc
---

Thêm từng tài khoản bằng `ccam add`, rồi khởi động Claude Code bằng `claude1`, `claude2` hoặc `ccs <name>`. Hầu hết mọi người chỉ cần như vậy.

## 1. Thêm tài khoản đầu tiên

```sh
ccam add
```

Trình hướng dẫn hỏi một cái tên, chẳng hạn `personal` hoặc `work`, và những mục nào tài khoản nên dùng chung với các tài khoản khác của bạn. Sau đó nó chạy `claude auth login`, lệnh này mở trình duyệt để bạn đăng nhập. Khi hoàn tất, bạn sẽ thấy:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

Hiện tại ccam in thông báo của chính nó bằng tiếng Hàn; các lệnh và tùy chọn thì giống nhau ở mọi ngôn ngữ.

## 2. Thêm tài khoản thứ hai

```sh
ccam add work
```

Đăng nhập bằng tài khoản còn lại trong trình duyệt. Tài khoản này sẽ trở thành `claude2`.

## 3. Chuyển đổi

Bất kỳ lệnh nào dưới đây cũng khởi động Claude Code bằng tài khoản công việc:

```sh
claude2
ccs work
ccs 2
```

Mọi thứ sau tên tài khoản được chuyển nguyên vẹn cho Claude Code, nên bạn vẫn dùng các tùy chọn quen thuộc:

```sh
ccs work --resume
ccs work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Chạy `ccs` một mình sẽ mở trình chọn: dùng các phím mũi tên và Enter, hoặc nhấn số của tài khoản.

## 4. Chọn những gì dùng chung

Theo mặc định, máy chủ MCP, `settings.json`, `CLAUDE.md` toàn cục, cùng agent, lệnh, skill và plugin được dùng chung giữa các tài khoản. Lịch sử hội thoại thì không. Hãy bật nó cho cả hai tài khoản nếu bạn muốn tiếp tục hội thoại qua lại giữa chúng:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam cho bạn xem những gì nó sẽ di chuyển trước khi thay đổi bất cứ thứ gì. Xem [dùng chung](/docs/sharing) để biết chi tiết.

## 5. Khi một tài khoản chạm giới hạn

```sh
ccs usage
ccs next --continue
```

`ccs usage` hiển thị mức sử dụng 5 giờ và hàng tuần của mọi tài khoản. `ccs next --continue` mở cùng cuộc hội thoại ở tài khoản kế tiếp. Xem [giới hạn sử dụng](/docs/usage-limits).

## Bạn đã dùng ~/.claude?

Hãy đăng ký thư mục đó thay vì đăng nhập lại:

```sh
ccam import
```

Thư mục vẫn nằm nguyên tại chỗ, và lệnh `claude` thông thường vẫn hoạt động như trước.
