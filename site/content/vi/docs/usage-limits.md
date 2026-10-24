---
title: Giới hạn sử dụng và chuyển đổi tài khoản
description: "Xem mức sử dụng 5 giờ và hàng tuần của mọi tài khoản Claude Code bằng ccam usage, và mở tài khoản còn nhiều hạn mức nhất bằng ccam next --auto."
updated: 2026-10-24
order: 5
translation: machine
sourceHash: "ccaed084207a"
---

`ccam usage` cho biết mỗi tài khoản đã dùng bao nhiêu trong giới hạn 5 giờ, giới hạn hàng tuần và giới hạn hàng tuần riêng của từng mô hình (như Fable). `ccam next --auto` mở tài khoản còn nhiều hạn mức nhất, và `--continue` mang cuộc hội thoại hiện tại sang tài khoản đó.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Mỗi tài khoản có một dòng cho biết độ cũ của các con số, tiếp theo là một thanh cho mỗi giới hạn cùng thời gian đến lần đặt lại: khoảng 5 giờ, khoảng hàng tuần cho mọi mô hình, và khoảng hàng tuần riêng của từng mô hình, như Fable. Từ 80% thanh chuyển sang màu vàng, và ở 100% thì màu đỏ. Dòng cuối nêu tên tài khoản mà `ccam next --auto` sẽ mở; với `auto_strategy = "expiring"` hoặc `"balance"` thì dòng này là `→ next --auto picks <name>: <reason>`.

**Seat của gói Team và Enterprise.** Lựa chọn tự động (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Vì vậy, dùng lựa chọn tự động giữa các seat có thể vi phạm điều khoản của Anthropic. Hãy xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó. Các gói cá nhân (Pro, Max) được đề cập trong [ghi chú về điều khoản ở phần Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Các con số lấy từ đâu

ccam hỏi endpoint mức sử dụng của Anthropic cho mọi tài khoản cùng lúc, dùng thông tin đăng nhập riêng của từng tài khoản. Nó chờ tối đa 3 giây cho mỗi tài khoản và hiển thị từng tài khoản ngay khi có phản hồi.

Endpoint này không được công bố tài liệu chính thức và có thể thay đổi. Khi một yêu cầu thất bại, ccam hiển thị các con số đã biết gần nhất và đánh dấu chúng là dữ liệu từ bộ nhớ đệm. Một khoảng giới hạn đã được đặt lại kể từ khi các con số đó được lưu sẽ hiển thị 0%. `ccam doctor` sẽ báo cho bạn nếu định dạng phản hồi đã thay đổi.

Khi Claude Code dừng một phiên ở giới hạn sử dụng, [plugin ccam](/docs/install#claude-code-plugin) ghi lại khoảng nào đã đầy và khi nào nó được đặt lại. Cho đến lúc đó, khoảng đó hiển thị 100% bất kể các con số khác nói gì, nên `ccam next --auto` và `--fallback` không chọn lại tài khoản đó quá sớm.

Các con số được lấy cách đây chưa đến 60 giây sẽ được dùng lại. Các tùy chọn:

| Tùy chọn | Tác dụng |
| --- | --- |
| `--watch[=seconds]` | Làm mới liên tục, mặc định mỗi 60 giây |
| `--json` | Đầu ra máy đọc được |
| `-v` | Hiển thị thêm nguồn gốc các con số của từng tài khoản |
| `--no-cache` | Luôn hỏi lại |
| `--source cache` | Không bao giờ hỏi; chỉ hiển thị các con số đã lưu |

## Đăng nhập hết hạn

Khi access token của một tài khoản đã hết hạn, ccam không gia hạn nó, vì việc gia hạn trong lúc một phiên Claude Code đang dùng cùng thông tin đăng nhập có thể khiến phiên đó bị đăng xuất. Tài khoản hiển thị ``access token expired; `ccam <name>` renews it when it starts``: hãy khởi động nó một lần bằng `ccam <name>` và Claude Code sẽ gia hạn nó. Một thông tin đăng nhập đã hết hạn hoàn toàn, hoặc đã bị Claude Code xóa, hiển thị `` login expired → `ccam login <name>` `` và cần đăng nhập lại.

`ccam usage --refresh` tự gia hạn các token đã hết hạn, nhưng chỉ cho những tài khoản không đang chạy, bất kể Claude Code được khởi động theo cách nào. Lệnh này không bao giờ gia hạn thông tin đăng nhập được lưu trong Chuỗi khóa (Keychain) của macOS hoặc trong một `~/.claude` đã nhập, và một lần gia hạn có thể mất tới 30 giây.

## Mức sử dụng trên dòng trạng thái

```sh
ccam statusline install
```

Sau đó Claude Code hiển thị tên tài khoản và mức sử dụng của nó trên dòng trạng thái, ví dụ `[work] · Opus · 5h 24% · 7d 41%`. Từ 90% trở lên, dòng này còn gợi ý `ccam next --auto --continue`. Nếu bạn đã có một lệnh cho dòng trạng thái, ccam giữ nguyên lệnh đó và đặt tên tài khoản lên trước. Dòng trạng thái cũng lưu lại các con số của nó, và `ccam usage` dùng chúng khi endpoint không phản hồi.

## Chuyển đổi tài khoản

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` bỏ qua các tài khoản chưa đăng nhập, và `--auto` bỏ qua các tài khoản đã ở mức 100% (kể cả một khoảng mà Claude Code báo là đã đầy, cho đến khi nó được đặt lại) và các tài khoản thuộc tổ chức khác. Khi hai tài khoản dùng chung lịch sử hội thoại, `--continue` mở cùng cuộc hội thoại đó; nếu không, Claude Code sao chép nó vào tài khoản mới. Xem [dùng chung](/docs/sharing) và [chuyển sang tài khoản tiếp theo](/features/next#continue).

Với các script và tác vụ định kỳ chạy `claude -p`, hãy đặt `--fallback` trước tên tài khoản. Khi tài khoản đó chạm giới hạn, ccam chạy lại cùng lệnh đó bằng tài khoản còn nhiều hạn mức nhất, và chỉ đầu ra của lần chạy đã hoàn tất được in ra:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Xem [Chạy lại khi chạm giới hạn](/features/next#run-again-at-the-limit) để biết ccam quyết định như thế nào.

[Hướng dẫn quy trình khi chạm giới hạn sử dụng](/guides/usage-limit-workflow) đi qua một ví dụ đầy đủ.
