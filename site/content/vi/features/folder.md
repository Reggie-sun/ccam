---
title: Tài khoản theo thư mục
description: "ccam folder cho một thư mục và mọi thư mục con dùng một tài khoản khi bạn chạy Claude Code không nêu tên tài khoản, ví dụ work cho mã công ty."
updated: 2026-10-20
order: 11
translation: machine
sourceHash: "46d1410a3d81"
---

`ccam folder work ~/code/company` khiến mọi lần khởi chạy không nêu tên tài khoản bên trong `~/code/company`, và trong mọi thư mục bên dưới nó, khởi động Claude Code bằng `work`. Bên ngoài thư mục đó, tài khoản mặc định vẫn khởi động như trước. Khi nêu tên một tài khoản, bằng `ccam personal` hoặc một lệnh tắt như `ccpersonal`, tài khoản đó luôn được khởi động.

## Cách dùng

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Các lần khởi chạy sau dùng tài khoản của thư mục:

- `ccam -p "…"` và mọi lần khởi chạy khác có đối số đầu tiên là một tùy chọn cho Claude Code
- `ccam run` không kèm tên, và `ccam --dry-run`
- `ccam --fallback -p "…"` cho lần thử đầu tiên
- bảng điều khiển, `ccam` đứng một mình, mở ra với con trỏ đặt sẵn ở tài khoản của thư mục

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `<account>` | Tên hoặc số của tài khoản cần dùng | |
| `[folder]` | Thư mục, tương đối hoặc tuyệt đối | thư mục hiện tại |
| `--remove` | Xóa mục của đúng thư mục đó | tắt |
| `--dry-run` | Được hỗ trợ: in ra thay đổi mà không lưu | tắt |

## Cách hoạt động

Mỗi thư mục là một mục `[[folder]]` trong `~/.config/ccam/config.toml`, tệp mà mọi lần khởi chạy vốn đã đọc:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` chuyển thư mục thành đường dẫn tuyệt đối, đã phân giải các liên kết tượng trưng, rồi lưu lại. Đặt lại một thư mục sẽ thay tài khoản của nó.
2. Một lần khởi chạy không nêu tên gọi `launcher.HereProfile`. Khi không có mục nào, nó lấy ngay tài khoản mặc định.
3. Nếu không, nó so sánh thư mục làm việc với từng `path` dưới dạng chuỗi: chính thư mục đó hoặc một thư mục bên dưới, không bao giờ là một thư mục chỉ bắt đầu bằng cùng các ký tự (`/code/company` không bao gồm `/code/companyx`). Mục sâu nhất được ưu tiên, nên `~/code` có thể là `personal` trong khi `~/code/company` là `work`.
4. Chỉ khi thư mục làm việc ở dạng nguyên bản không khớp mục nào, ccam mới phân giải các liên kết tượng trưng của nó và so sánh thêm một lần.
5. Trên Windows, phép so sánh không phân biệt chữ hoa và chữ thường, giống như Windows.

Ngoài `config.toml` không có tệp nào khác được đọc, nên một lần khởi chạy bên trong một thư mục không tốn thêm thời gian đo được (4.0 ms dù có hay không có mục nào, `hyperfine`, 300 lần chạy).

## An toàn và các trường hợp đặc biệt

- `ccam rename` chuyển các mục thư mục sang tên mới, và `ccam remove` xóa các mục của tài khoản bị xóa.
- Một mục có tài khoản không còn tồn tại sẽ dừng lần khởi chạy không nêu tên trong thư mục đó: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` báo cáo các mục như vậy, và cả các mục có thư mục đã mất.
- `ccam folder --remove` bên trong một thư mục nhận tài khoản từ một thư mục phía trên sẽ nêu tên thư mục đó thay vì xóa bất cứ thứ gì.
- `ccam next` chọn theo vòng xoay hoặc mức sử dụng và không xét đến thư mục.

## Liên quan

- [Khởi động Claude Code](/features/launch)
- [Lệnh tắt](/features/shortcuts)
- [Quản lý tài khoản](/features/accounts)
- [Tham chiếu lệnh](/docs/commands)
