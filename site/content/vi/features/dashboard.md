---
title: Bảng điều khiển
description: "ccam trên terminal mở bảng điều khiển toàn màn hình gồm tài khoản, nhóm, quy tắc thư mục, tiện ích mở rộng và kiểm tra; mọi thay đổi đều hiện kế hoạch trước."
updated: 2026-10-21
order: 15
translation: machine
sourceHash: "27bd153b48b2"
---

`ccam` không kèm đối số sẽ mở một bảng điều khiển toàn màn hình khi đầu vào và đầu ra chuẩn là terminal; nếu không, nó in ra `ccam list`. Bảng điều khiển hiển thị mọi tài khoản cùng mức sử dụng và các nhóm của chúng, và làm được những gì các lệnh làm: khởi động, thêm, đăng nhập, thay đổi nhóm, đặt quy tắc thư mục, kiểm tra và sửa chữa. Khởi động một tài khoản sẽ đóng bảng điều khiển; thêm tài khoản hoặc đăng nhập thì tạm rời nó để dùng terminal rồi quay lại.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Dòng dưới tên các thẻ cho biết trong một câu thẻ hiện tại dùng để làm gì. Trên terminal rộng hơn, thẻ Accounts còn hiển thị email của từng tài khoản.

Khung hình đầu tiên được vẽ chỉ từ `config.toml`, trong chưa đến 100 ms; mức sử dụng được điền vào khi từng tài khoản phản hồi. Mỗi giới hạn có thanh riêng: khoảng 5 giờ, khoảng hàng tuần cho mọi mô hình, và khoảng hàng tuần riêng của Fable; `-` nghĩa là tài khoản không có khoảng đó. Một tỷ lệ phần trăm chuyển sang màu vàng từ 80% và màu đỏ ở 100%, và các giá trị lấy từ bộ nhớ đệm sau một lần tra cứu thất bại được vẽ mờ.

## Các thẻ

| Thẻ | Nội dung hiển thị | Phím |
| --- | --- | --- |
| Accounts | số, tên, email, thanh 5 giờ, hàng tuần và Fable, lần dùng gần nhất, nhóm; `✔` tài khoản mặc định, `▸` tài khoản của thư mục này, `•` đang chạy | bên dưới |
| Groups | từng nhóm, những gì nhóm dùng chung và các tài khoản của nhóm, cùng các tài khoản không thuộc nhóm nào | `enter` sửa, `n` tạo mới, `e` đổi tên, `x` xóa |
| Folders | các quy tắc [`ccam folder`](/features/folder) và quy tắc áp dụng ở đây | `n` dùng một tài khoản ở đây, `enter` thay đổi, `x` xóa |
| Extensions | những gì [`ccam ext`](/features/ext) hiển thị, mỗi tiện ích mở rộng một lần cùng các nhóm mà nó đến từ | `r` tải lại |
| Health | những gì [`ccam doctor`](/features/doctor) báo cáo, dưới dạng một hộp cho mỗi tài khoản (bên dưới) | `f` sửa chữa hộp đang chọn, `F` sửa chữa tất cả (`doctor --fix`), `↑`/`↓` di chuyển giữa các hộp |

`←`/`→` hoặc `Tab` chuyển thẻ, `?` liệt kê các phím của thẻ hiện tại, `r` tra cứu lại mức sử dụng mà không dùng bộ nhớ đệm ngắn hạn, và `q` hoặc `Esc` để thoát.

## Phím trên thẻ Accounts

| Phím | Tác dụng |
| --- | --- |
| `Enter` | Khởi động Claude Code bằng tài khoản đang chọn |
| `1`–`9` | Di chuyển tới tài khoản có số đó; gõ nhanh `1` `2` để tới 12. Chỉ `Enter` mới khởi động |
| `↑`/`↓`, `k`/`j` | Di chuyển |
| `n` | Khởi động tài khoản mà [`ccam next --auto`](/features/next) sẽ chọn |
| `o` | Khởi động với tùy chọn: `--continue`, `--resume`, các đối số khác, hoặc chạy thử |
| `u` | Chi tiết mức sử dụng: khoảng 5 giờ, hàng tuần và theo mô hình, cùng thời điểm đặt lại và nguồn của các con số |
| `a` | Thêm một tài khoản (trình hướng dẫn bên dưới) |
| `l` | Đăng nhập lại |
| `s` | Chọn các nhóm của tài khoản |
| `d` | Đặt nó làm tài khoản mặc định |
| `f` | Dùng nó trong thư mục hiện tại |
| `e` | Đổi tên nó |
| `c` | Đặt hoặc xóa lệnh tắt của nó |
| `x` | Xóa nó (hỏi trước, con trỏ bắt đầu ở No) |

## Thẻ Health

Thẻ Health chạy các kiểm tra của [`ccam doctor`](/features/doctor) khi bạn mở nó và vẽ cùng một giao diện: dòng đầu có một ký hiệu cho mỗi tài khoản, rồi đến một hộp cho mỗi tài khoản. Mỗi hàng trong hộp là một mục đã được kiểm tra, `Folder`, `Login` và `Share links` đứng đầu, được đánh dấu `✔` ổn, `ℹ` thông tin, `⚠` cảnh báo hoặc `✘` lỗi, nên bạn thấy ngay tài khoản nào có vấn đề và ở đâu. Các kết quả không thuộc tài khoản nào nằm trong hộp `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` di chuyển từ hộp này sang hộp khác, và hộp đang chọn có viền màu nhấn. `f` chỉ sửa chữa hộp đang chọn, theo cách `ccam doctor --fix` sẽ làm: trên hộp của một tài khoản, nó chỉ sửa chữa tài khoản đó, còn trên hộp `General` là các kết quả không thuộc tài khoản nào (đường dẫn plugin và lệnh tắt); mọi tài khoản khác chỉ được kiểm tra lại. Nó hỏi trước, và một hộp không có gì để sửa chữa chỉ hiển thị `nothing to repair in <name>`. `F` sửa chữa mọi thứ cùng lúc, như `ccam doctor --fix`, và `r` kiểm tra lại.

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
