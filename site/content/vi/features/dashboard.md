---
title: Bảng điều khiển
description: "ccam trên terminal mở bảng điều khiển toàn màn hình gồm tài khoản, nhóm, quy tắc thư mục, tiện ích mở rộng và kiểm tra; mọi thay đổi đều hiện kế hoạch trước."
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` không kèm đối số sẽ mở một bảng điều khiển toàn màn hình khi đầu vào và đầu ra chuẩn là terminal; nếu không, nó in ra `ccam list`. Bảng điều khiển hiển thị mọi tài khoản cùng mức sử dụng và các nhóm của chúng, và làm được những gì các lệnh làm: khởi động, thêm, đăng nhập, thay đổi nhóm, đặt quy tắc thư mục, kiểm tra và sửa chữa. Khởi động một tài khoản sẽ đóng bảng điều khiển; thêm tài khoản hoặc đăng nhập thì tạm rời nó để dùng terminal rồi quay lại.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Khung hình đầu tiên được vẽ chỉ từ `config.toml`, trong chưa đến 100 ms; mức sử dụng được điền vào khi từng tài khoản phản hồi. Một tỷ lệ phần trăm chuyển sang màu vàng từ 80% và màu đỏ ở 100%, và các giá trị lấy từ bộ nhớ đệm sau một lần tra cứu thất bại được vẽ mờ.

## Các thẻ

| Thẻ | Nội dung hiển thị | Phím |
| --- | --- | --- |
| Accounts | số, tên, email, thanh 5 giờ, mức sử dụng 7 ngày, lần dùng gần nhất, nhóm; `✔` tài khoản mặc định, `▸` tài khoản của thư mục này, `•` đang chạy | bên dưới |
| Groups | từng nhóm, những gì nhóm dùng chung và các tài khoản của nhóm, cùng các tài khoản không thuộc nhóm nào | `enter` sửa, `n` tạo mới, `e` đổi tên, `x` xóa |
| Folders | các quy tắc [`ccam folder`](/features/folder) và quy tắc áp dụng ở đây | `n` dùng một tài khoản ở đây, `enter` thay đổi, `x` xóa |
| Extensions | những gì [`ccam ext`](/features/ext) hiển thị, mỗi tiện ích mở rộng một lần cùng các nhóm mà nó đến từ | `r` tải lại |
| Health | những gì [`ccam doctor`](/features/doctor) báo cáo | `f` sửa chữa (`doctor --fix`) |

`←`/`→` hoặc `Tab` chuyển thẻ, `?` liệt kê các phím của thẻ hiện tại, `r` tra cứu lại mức sử dụng mà không dùng bộ nhớ đệm ngắn hạn, và `q` hoặc `Esc` để thoát.

## Phím trên thẻ Accounts

| Phím | Tác dụng |
| --- | --- |
| `Enter` | Khởi động Claude Code bằng tài khoản đang chọn |
| `1`–`9` | Di chuyển tới tài khoản có số đó; gõ nhanh `1` `2` để tới 12. Chỉ `Enter` mới khởi động |
| `↑`/`↓`, `k`/`j` | Di chuyển |
| `n` | Khởi động tài khoản mà [`ccam next --auto`](/features/next) sẽ chọn |
| `o` | Khởi động với tùy chọn: `--continue`, `--resume`, các đối số khác, hoặc chạy thử |
| `u` | Chi tiết mức sử dụng: khoảng 5 giờ, 7 ngày và theo mô hình, cùng thời điểm đặt lại và nguồn của các con số |
| `a` | Thêm một tài khoản (trình hướng dẫn bên dưới) |
| `l` | Đăng nhập lại |
| `s` | Chọn các nhóm của tài khoản |
| `d` | Đặt nó làm tài khoản mặc định |
| `f` | Dùng nó trong thư mục hiện tại |
| `e` | Đổi tên nó |
| `c` | Đặt hoặc xóa lệnh tắt của nó |
| `x` | Xóa nó (hỏi trước, con trỏ bắt đầu ở No) |

## Thay đổi nhóm

Trên thẻ Groups, `Enter` mở trình chỉnh sửa: các mục ở bên trái, các tài khoản ở bên phải, `Tab` để chuyển qua lại, `Space` để bật/tắt. Một mục mà một trong các tài khoản đã nhận từ một nhóm khác sẽ bị làm mờ kèm lý do, và một tài khoản sẽ nhận một mục hai lần cũng vậy; một mục chỉ có thể đến từ một nhóm. `Enter` hiển thị kế hoạch:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Một tệp khác với bản của nhóm sẽ được hỏi trước (giữ bản của nhóm, dùng bản của tôi, hoặc gộp). Sau đó thay đổi chạy ở chế độ nền với cùng mã như [`ccam group`](/features/share): cấu hình được lưu trước, dữ liệu của từng tài khoản được di chuyển theo thứ tự, và bản gốc được đưa vào `.ccam-backup/`. Một thay đổi sẽ di chuyển dữ liệu của một tài khoản đang chạy bị từ chối kèm tên tài khoản đó.

## Thêm tài khoản

`a` mở trình hướng dẫn mà `ccam add` cũng dùng trên terminal:

1. Tên tài khoản.
2. Cách nó đăng nhập: một gói đăng ký Claude ngay bây giờ (trong trình duyệt), một [API key](/features/api), hoặc một gói đăng ký Claude sau.
3. Với API key: nơi ccam đọc khóa mỗi lần khởi động (`env:NAME` hoặc `command:COMMAND`), endpoint, cách gửi khóa và mô hình.
4. Các nhóm để tham gia, với `default` đã được đánh dấu; khi chưa có nhóm nào, đó là những gì nhóm `default` dùng chung.
5. Một lệnh tắt, nếu bạn muốn.
6. Một bản tóm tắt để xác nhận.

`Esc` quay lại một bước, và ở bước đầu tiên thì đóng trình hướng dẫn.

## Màu sắc

Màu sắc theo các vai trò trong giao diện terminal của Claude Code, nên ccam trông hài hòa khi đặt cạnh nó; tên nhóm nhận một màu dựa trên chính tên đó. ccam hỏi terminal về màu nền và chọn chủ đề sáng hoặc tối. Bạn có thể tự đặt một chủ đề trong `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Các chủ đề `-daltonized` dành cho người mù màu đỏ-lục, còn các chủ đề `-ansi` dùng 16 màu của chính terminal. `NO_COLOR` tắt màu.

## An toàn và các trường hợp đặc biệt

- Bảng điều khiển gọi cùng các hàm như các lệnh, nên mọi quy tắc đều được giữ: kế hoạch trước khi dữ liệu được di chuyển, từ chối với tài khoản đang chạy, bản sao lưu, và token đăng nhập không bao giờ rời khỏi tài khoản.
- Bảng điều khiển không có chế độ chạy thử: `ccam --dry-run` không kèm tên sẽ in ra cách tài khoản của thư mục này sẽ khởi động, và `o` trên một tài khoản cho phép chạy thử riêng lần khởi động đó.
- Terminal được khôi phục trước khi Claude Code khởi động hoặc việc đăng nhập mở ra.
- Một màn hình quá nhỏ cho hộp thoại sẽ bỏ bớt các hàng của thẻ trước; các dòng bị cắt bằng `…` thay vì xuống dòng.

## Liên quan

- [Dùng chung bằng nhóm](/features/share)
- [Khởi động Claude Code](/features/launch)
- [Quản lý tài khoản](/features/accounts)
- [Giới hạn sử dụng](/features/usage)
