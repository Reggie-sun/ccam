---
title: Cài đặt ccam
metaTitle: "Cài đặt ccam trên Linux và macOS"
description: "Cài đặt ccam trên Linux hoặc macOS bằng một lệnh, từ gói .deb hoặc từ mã nguồn, và cập nhật bằng ccam update. Cần cài Claude Code trước."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 7ecd1ed00e09
---

Chạy lệnh này trong terminal trên Linux hoặc macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Script tải về bản phát hành cho hệ thống của bạn, kiểm tra mã băm SHA-256 và đặt `ccam` vào `~/.local/bin`, nơi ccam cũng tạo các lệnh tắt mà bạn chọn sau này, chẳng hạn `ccwork`. Nếu thư mục đó chưa nằm trong `PATH`, script sẽ in ra dòng cần thêm vào `~/.bashrc` hoặc `~/.zshrc`.

ccam khởi động Claude Code cho bạn, vì vậy hãy cài [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) trước và đảm bảo `claude --version` chạy được.

## Hệ điều hành được hỗ trợ

- Linux trên x86-64 và ARM64: Ubuntu, Debian, Fedora, Arch, Alpine và các bản phân phối khác. Tệp thực thi là tệp tĩnh và không có phụ thuộc.
- macOS trên Apple silicon và Intel.
- Windows không được hỗ trợ trực tiếp. Bên trong WSL, ccam hoạt động như trên bất kỳ hệ thống Linux nào.

## Các cách cài đặt khác

**Gói cho Debian và Ubuntu.** Tải `ccam_<version>_linux_amd64.deb` từ [trang phát hành](https://github.com/jhyunwoo/ccam/releases) và cài đặt:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Gói này cài `/usr/bin/ccam`.

**Từ mã nguồn.** Với Go 1.26 trở lên:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Một phiên bản cụ thể.** Đặt `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Tự động hoàn thành trong shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Tính năng tự động hoàn thành áp dụng cho mọi lệnh và tùy chọn của `ccam`, bao gồm cả tên tài khoản.

## Kiểm tra cài đặt

```sh
ccam version
ccam doctor
```

`ccam doctor` kiểm tra xem có tìm thấy Claude Code không và mọi thư mục tài khoản, liên kết và thông tin đăng nhập có ổn không. Tiếp tục với phần [bắt đầu nhanh](/docs/quick-start).

## Cập nhật

```sh
ccam update --check
ccam update
```

`ccam update --check` chỉ cho bạn biết có bản phát hành mới hơn hay không. `ccam update` tải về bản phát hành cho hệ thống của bạn, đối chiếu mã băm SHA-256 với tệp `checksums.txt` của bản phát hành đó, rồi thay thế tệp `ccam` trong một bước. Nếu mã băm không khớp, sẽ không có gì bị thay đổi. Các lệnh tắt như `ccwork` là liên kết tới ccam, nên chúng chạy phiên bản mới ngay lập tức.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam không bao giờ tự kiểm tra cập nhật; nó chỉ tìm bản phát hành mới khi bạn chạy lệnh này. Để cài một bản phát hành cụ thể, kể cả bản cũ hơn, hãy truyền `--version`. `ccam --dry-run update` cho biết những gì sẽ được tải về và thay thế mà không thay đổi gì.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Chạy lại script cài đặt cũng sẽ cập nhật bản sao trong `~/.local/bin`. Nếu bạn đã build ccam bằng `go install`, hãy chạy lại chính lệnh đó.

**Gói cho Debian và Ubuntu.** `ccam update` không thay thế `/usr/bin/ccam` do apt cài đặt mà in ra các bước sau. Tải `ccam_<version>_linux_amd64.deb` mới từ [trang phát hành](https://github.com/jhyunwoo/ccam/releases) rồi cài đè lên bản cũ:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Gỡ cài đặt

Trước tiên, hãy xóa các lệnh tắt của tài khoản khi ccam vẫn còn được cài. `ccam list` hiển thị chúng trong cột Shortcut; xóa từng lệnh bằng `ccam shortcut <account> --remove`, hoặc tự xóa các liên kết đó khỏi `~/.local/bin`. Sau đó xóa lệnh ccam, và xóa các thư mục dữ liệu nếu bạn không còn cần các tài khoản:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Xóa `~/.local/share/ccam` cũng xóa luôn thông tin đăng nhập và lịch sử hội thoại được lưu ở đó. Thư mục `~/.claude` mà bạn đã [nhập](/docs/commands) vẫn nằm nguyên tại chỗ.
