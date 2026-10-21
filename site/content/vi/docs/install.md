---
title: Cài đặt ccam
metaTitle: "Cài đặt ccam trên Linux, macOS và Windows"
description: "Cài đặt ccam trên Linux, macOS hoặc Windows bằng một lệnh, từ gói .deb hoặc từ mã nguồn, và cập nhật bằng ccam update. Cần cài Claude Code trước."
updated: 2026-10-21
order: 1
translation: machine
sourceHash: "1e3a53623de0"
---

Chạy lệnh này trong terminal trên Linux hoặc macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Script tải về bản phát hành cho hệ thống của bạn, kiểm tra mã băm SHA-256 và đặt `ccam` vào `~/.local/bin`, nơi ccam cũng tạo các lệnh tắt mà bạn chọn sau này, chẳng hạn `ccwork`. Nếu thư mục đó chưa nằm trong `PATH`, script sẽ in ra dòng cần thêm vào `~/.bashrc` hoặc `~/.zshrc`.

Trên Windows, hãy chạy lệnh này trong PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Script tải về `ccam_windows_<arch>.zip`, kiểm tra mã băm SHA-256, đặt `ccam.exe` vào `%USERPROFILE%\.local\bin`, thư mục mà trình cài đặt của Claude Code đặt `claude.exe`, và thêm thư mục đó vào `PATH` của người dùng nếu chưa có. Sau đó hãy mở một terminal mới.

ccam khởi động Claude Code cho bạn, vì vậy hãy cài [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) trước và đảm bảo `claude --version` chạy được.

## Hệ điều hành được hỗ trợ

- Linux trên x86-64 và ARM64: Ubuntu, Debian, Fedora, Arch, Alpine và các bản phân phối khác. Tệp thực thi là tệp tĩnh và không có phụ thuộc.
- macOS trên Apple silicon và Intel.
- Windows 10 và 11 trên x86-64 và ARM64, trong PowerShell, Windows Terminal hoặc Git Bash. Thư mục dùng chung hoạt động với mọi người dùng. Tệp dùng chung (`settings.json`, `CLAUDE.md`, `history.jsonl`) cần bật Developer Mode (Settings > System > For developers), vì Windows chỉ cho phép chế độ này, hoặc quản trị viên, tạo liên kết tới tệp; xem [dùng chung trên Windows](/features/share#on-windows).
- Bên trong WSL, ccam hoạt động như trên bất kỳ hệ thống Linux nào.

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

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Tự động hoàn thành trong shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Tính năng tự động hoàn thành áp dụng cho mọi lệnh và tùy chọn của `ccam`, bao gồm cả tên tài khoản.

## Plugin Claude Code

Plugin ccam dạy chính Claude Code cách dùng ccam: kiểm tra tài khoản nào còn hạn mức, khởi động các worker `claude -p` song song trên các tài khoản khác nhau bằng `ccam pick`, và giải thích vì sao thiếu một skill hay máy chủ MCP bằng `ccam ext`. Cài nó ngay trong Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Plugin thêm một skill, các lệnh `/ccam:usage`, `/ccam:ext` và `/ccam:switch`, cùng hai hook: khi một phiên dừng lại vì giới hạn sử dụng, bạn nhận được một thông báo trên terminal, và `/ccam:switch [account]` chuyển cuộc hội thoại sang một tài khoản khác ([chi tiết](/features/next#inside-claude-code)). Các hook chạy `ccam`, nên ccam phải nằm trong `PATH` mà Claude Code thấy. Skill này không bao giờ chạy các lệnh đăng nhập hay di chuyển tệp; thay vào đó nó gợi ý các lệnh đó cho bạn.

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

Bảng điều khiển (`ccam` trên terminal) cũng tìm bản phát hành mới mỗi lần mở và, khi có, sẽ đề nghị cập nhật: nhấn `U` rồi xác nhận. Không lệnh nào khác tìm, và ccam không bao giờ cập nhật mà không hỏi. Để cài một bản phát hành cụ thể, kể cả bản cũ hơn, hãy truyền `--version`. `ccam --dry-run update` cho biết những gì sẽ được tải về và thay thế mà không thay đổi gì.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Chạy lại script cài đặt cũng sẽ cập nhật bản sao trong `~/.local/bin`. Nếu bạn đã build ccam bằng `go install`, hãy chạy lại chính lệnh đó.

**Windows.** `ccam update` tải về tệp `.zip` cho Windows. Trên Windows không thể thay thế một chương trình đang chạy, nên tệp thực thi cũ được đổi tên thành `ccam.exe.old`, và lần cập nhật tiếp theo sẽ xóa nó.

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

Trên Windows, trong PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Xóa `~/.local/share/ccam` cũng xóa luôn thông tin đăng nhập và lịch sử hội thoại được lưu ở đó. Thư mục `~/.claude` mà bạn đã [nhập](/docs/commands) vẫn nằm nguyên tại chỗ.
