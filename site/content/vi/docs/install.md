---
title: Cài đặt ccam
description: "Cài đặt ccam trên Linux hoặc macOS bằng một lệnh, từ gói .deb hoặc từ mã nguồn. Cần cài Claude Code trước."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Chạy lệnh này trong terminal trên Linux hoặc macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Script tải về bản phát hành cho hệ thống của bạn, kiểm tra mã băm SHA-256 và đặt `ccam` vào `~/.local/bin`, nơi `ccam add` sau này tạo các lệnh tắt `claude1`, `claude2`, v.v. Nếu thư mục đó chưa nằm trong `PATH`, script sẽ in ra dòng cần thêm vào `~/.bashrc` hoặc `~/.zshrc`.

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

## Gỡ cài đặt

Xóa lệnh ccam và các lệnh tắt của tài khoản, sau đó xóa các thư mục dữ liệu nếu bạn không còn cần các tài khoản:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Xóa `~/.local/share/ccam` cũng xóa luôn thông tin đăng nhập và lịch sử hội thoại được lưu ở đó. Thư mục `~/.claude` mà bạn đã [nhập](/docs/commands) vẫn nằm nguyên tại chỗ.
