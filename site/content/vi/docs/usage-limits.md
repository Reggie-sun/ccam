---
title: Giới hạn sử dụng và chuyển đổi tài khoản
description: "Xem mức sử dụng 5 giờ và hàng tuần của mọi tài khoản Claude Code bằng ccam usage, và mở tài khoản còn nhiều hạn mức nhất bằng ccam next --auto."
updated: 2026-10-20
order: 5
translation: machine
sourceHash: "40aef0708c81"
---

`ccam usage` cho biết mỗi tài khoản đã dùng bao nhiêu trong giới hạn 5 giờ và giới hạn hàng tuần. `ccam next --auto` mở tài khoản còn nhiều hạn mức nhất, và `--continue` mang cuộc hội thoại hiện tại sang tài khoản đó.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

Các cột lần lượt là: tài khoản, mức sử dụng 5 giờ, thời gian đến lần đặt lại, mức sử dụng hàng tuần, thời gian đến lần đặt lại, và độ cũ của các con số.

## Các con số lấy từ đâu

ccam hỏi endpoint mức sử dụng của Anthropic cho mọi tài khoản cùng lúc, dùng thông tin đăng nhập riêng của từng tài khoản. Nó chờ tối đa 3 giây cho mỗi tài khoản và hiển thị từng tài khoản ngay khi có phản hồi.

Endpoint này không được công bố tài liệu chính thức và có thể thay đổi. Khi một yêu cầu thất bại, ccam hiển thị các con số đã biết gần nhất và đánh dấu chúng là dữ liệu từ bộ nhớ đệm. Một khoảng giới hạn đã được đặt lại kể từ khi các con số đó được lưu sẽ hiển thị 0%. `ccam doctor` sẽ báo cho bạn nếu định dạng phản hồi đã thay đổi.

Các con số được lấy cách đây chưa đến 60 giây sẽ được dùng lại. Các tùy chọn:

| Tùy chọn | Tác dụng |
| --- | --- |
| `--watch [seconds]` | Làm mới liên tục, mặc định mỗi 60 giây |
| `--json` | Đầu ra máy đọc được |
| `-v` | Hiển thị thêm giới hạn hàng tuần theo từng mô hình và nguồn gốc của từng con số |
| `--no-cache` | Luôn hỏi lại |
| `--source cache` | Không bao giờ hỏi; chỉ hiển thị các con số đã lưu |

## Đăng nhập hết hạn

Khi access token của một tài khoản đã hết hạn, ccam không gia hạn nó, vì việc gia hạn trong lúc một phiên Claude Code đang dùng cùng thông tin đăng nhập có thể khiến phiên đó bị đăng xuất. Hãy khởi động tài khoản một lần bằng `ccam <name>` và Claude Code sẽ gia hạn nó.

`ccam usage --refresh` tự gia hạn các token đã hết hạn, nhưng chỉ cho những tài khoản không đang chạy. Lệnh này không bao giờ gia hạn thông tin đăng nhập được lưu trong Chuỗi khóa (Keychain) của macOS hoặc trong một `~/.claude` đã nhập.

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

`ccam next` bỏ qua các tài khoản chưa đăng nhập, và `--auto` bỏ qua các tài khoản đã ở mức 100% và các tài khoản thuộc tổ chức khác. Khi hai tài khoản dùng chung lịch sử hội thoại, `--continue` mở cùng cuộc hội thoại đó; nếu không, Claude Code sao chép nó vào tài khoản mới. Xem [dùng chung](/docs/sharing) và [chuyển sang tài khoản tiếp theo](/features/next#continue).

Với các script và tác vụ định kỳ chạy `claude -p`, hãy đặt `--fallback` trước tên tài khoản. Khi tài khoản đó chạm giới hạn, ccam chạy lại cùng lệnh đó bằng tài khoản còn nhiều hạn mức nhất, và chỉ đầu ra của lần chạy đã hoàn tất được in ra:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Xem [Chạy lại khi chạm giới hạn](/features/next#run-again-at-the-limit) để biết ccam quyết định như thế nào.

[Hướng dẫn quy trình khi chạm giới hạn sử dụng](/guides/usage-limit-workflow) đi qua một ví dụ đầy đủ.
