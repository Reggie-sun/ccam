---
title: Lệnh tắt
description: "Đặt cho một tài khoản một lệnh riêng như ccwork bằng ccam shortcut, và tạo lại các liên kết trong ~/.local/bin bằng ccam shims install."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

Lệnh tắt là một lệnh do bạn tự đặt tên, dùng để khởi động Claude Code bằng một tài khoản: `ccwork --resume` làm đúng những gì `ccam work --resume` làm. Nó là một liên kết đến ccam trong `~/.local/bin`, nên nó hoạt động trong bash, zsh, script và mọi thứ khác chạy chương trình từ `PATH`. Lệnh tắt là tùy chọn; hãy tạo cho các tài khoản bạn hay khởi động.

## Cách dùng

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam xác nhận từng thay đổi:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Mỗi tài khoản có tối đa một lệnh tắt, được lưu dưới dạng `shortcut` trong `~/.config/ccam/config.toml`. Trình hướng dẫn `ccam add` hỏi `Create a shortcut command for <name>? [y/N]` rồi `Shortcut name [cc<name>]:`. Ngoài trình hướng dẫn, `ccam add` chỉ tạo lệnh tắt khi có `--shortcut`. `ccam list` hiển thị lệnh tắt của từng tài khoản, hoặc `-`.

Lệnh tắt chuyển mọi đối số cho Claude Code, giống như `ccam <name>` làm với các đối số sau tên. Điều đó bao gồm cả `--dry-run` và `--help`: để xem trước một lần khởi động, hãy dùng `ccam --dry-run work`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--remove` | Xóa lệnh tắt của tài khoản cùng liên kết của nó | tắt |
| `-h`, `--help` | Trợ giúp cho lệnh | |
| `--dry-run` | `ccam shortcut` hỗ trợ nó: kiểm tra tên và in ra `Would set the shortcut of work to ccwork.` hoặc `Would remove ccwork, the shortcut of work.` mà không lưu hay tạo liên kết | tắt |

`ccam shims install` không nhận đối số và từ chối `--dry-run` với `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Quy tắc đặt tên, được kiểm tra theo thứ tự sau:

- từ 1 đến 64 ký tự gồm `a-z`, `0-9`, `.`, `-` và `_`, bắt đầu bằng một chữ cái hoặc chữ số;
- không phải `ccam`, `claude` hay một lệnh của ccam như `run` hoặc `usage`;
- không phải lệnh tắt của tài khoản khác;
- không phải một tệp trong `~/.local/bin` mà chưa phải là liên kết đến ccam;
- không phải tên của một chương trình khác trong `PATH` của bạn.

Đưa cả tên lệnh lẫn `--remove` cùng lúc sẽ dừng với `give a command or --remove, not both`. Đặt đúng tên mà tài khoản đã có sẽ in ra `work: ccwork (unchanged)` và tạo lại liên kết nếu nó bị thiếu.

## Cách hoạt động

Việc đặt tên chạy `checkShortcut`, hàm này gọi `shim.Check` với thư mục bin và đường dẫn thật của tệp thực thi ccam:

<!-- motion:shortcuts-check -->

Tên đi lần lượt qua các bước kiểm tra. Bên cạnh mỗi bước là một tên bị từ chối ở đó: `Work` không khớp mẫu, `run` là một lệnh của ccam, `ccside` thuộc về tài khoản khác, `ccfile` là một tệp đã có trong `~/.local/bin`, và `git` sẽ che mất `/usr/bin/git`. Chỉ tên vượt qua mọi bước kiểm tra mới được lưu và liên kết.

1. `shim.Check` so khớp tên với `^[a-z0-9][a-z0-9._-]{0,63}$`, rồi với `config.Reserved`, rồi với `cfg.FindShortcut`. Nó gọi `os.Lstat` trên `~/.local/bin/<name>`, và với mọi thư mục khác trong `PATH`, nó gọi `os.Stat` và từ chối một tệp thực thi không phải là liên kết đến ccam.
2. ccam lưu `config.toml` bằng `fsutil.AtomicWrite` (tệp tạm, `fsync`, đổi tên).
3. `shim.Remove` xóa liên kết của lệnh tắt cũ, nhưng chỉ khi đó là liên kết đến ccam.
4. `shim.Install` tạo liên kết cho mọi lệnh tắt đã cấu hình mà còn thiếu, bằng `os.Symlink` tới tệp thực thi ccam, được xác định bằng `os.Executable` và `filepath.EvalSymlinks`. Nó không bao giờ xóa thứ gì.

Một liên kết được coi là của ccam khi `os.Readlink` trả về đường dẫn ccam, hoặc khi liên kết và đường dẫn ccam trỏ tới cùng một tệp. Mọi thao tác xóa và mọi kiểm tra ghi đè đều đi qua phép thử này, nên ccam không động đến các tệp nó không tạo ra.

Khi bạn chạy một lệnh tắt, tên bạn gõ là thứ duy nhất ccam cần:

<!-- motion:shortcuts-link -->

Shell tìm thấy `ccwork` trong `PATH`, nhân hệ điều hành đi theo liên kết và khởi động tệp thực thi ccam, và `argv[0]` vẫn là `ccwork`. Trong `app.Main`, mọi tên khác `ccam` đều tải `config.toml` một lần và gọi `shim.Lookup`, hàm này tìm tài khoản có `shortcut` bằng `filepath.Base(argv[0])`. Sau đó `launcher.LaunchWith` đi qua cùng các bước như [`ccam <name>`](/features/launch) với cấu hình đã được tải, và chuyển nguyên vẹn `argv[1:]` cho `claude`.

Vì sao dùng liên kết tượng trưng mà không dùng alias của shell hay script bao bọc: liên kết tượng trưng hoạt động trong mọi shell và trong script, không phải khởi động thêm một shell (khoảng 1 đến 2 ms), và lần khởi động vẫn chỉ đọc một tệp. `scripts/perf-check.sh` đo một lệnh tắt bằng `hyperfine` với cùng giới hạn 10 ms (p95) như `ccam <name>`.

## An toàn và các trường hợp đặc biệt

- **Liên kết không thuộc tài khoản nào.** Nếu một liên kết đến ccam trong `~/.local/bin` có tên mà không tài khoản nào dùng, chẳng hạn sau khi bạn sửa `shortcut` trong `config.toml` bằng tay, ccam không quay về tài khoản mặc định. Nó dừng với `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Hãy đặt lại tên bằng `ccam shortcut`, hoặc xóa liên kết.
- **Tệp đã có.** ccam không bao giờ thay thế một tệp nó không tạo ra. `ccam shims install` dừng với `<path> already exists and is not a ccam link, so it was left alone`.
- **Tệp thực thi ccam bị di chuyển.** Khi đó các liên kết cũ trỏ tới một đường dẫn không còn tồn tại, và ccam không coi chúng là của mình. Hãy xóa chúng, rồi chạy `ccam shims install`.
- **Đổi tên và xóa.** `ccam rename` giữ nguyên tên lệnh tắt, và liên kết vẫn hoạt động, vì nó trỏ tới ccam và ccam tra tên trong `config.toml`. `ccam remove` xóa liên kết của tài khoản. `ccam shortcut <account> <new>` xóa liên kết cũ trước khi tạo liên kết mới.
- **PATH.** Khi ccam tạo một liên kết mà `~/.local/bin` không nằm trong `PATH` của bạn, nó in ra `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` theo sau là dòng `export PATH=...` cần thêm.
- **Kiểm tra về sau.** `ccam doctor` báo `missing shortcuts in /home/me/.local/bin: ccwork` và cảnh báo khi thư mục không nằm trong `PATH`. `ccam doctor --fix` tạo các liên kết còn thiếu.

## Liên quan

- [Khởi động Claude Code](/features/launch)
- [Quản lý tài khoản](/features/accounts)
- [Tham chiếu lệnh](/docs/commands#shortcuts)
- [Cài đặt ccam](/docs/install)
