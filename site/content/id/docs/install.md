---
title: Instal ccam
metaTitle: Instal ccam di Linux, macOS, dan Windows
description: Instal ccam di Linux, macOS, atau Windows dengan satu perintah, paket .deb, atau kode sumber, lalu perbarui dengan ccam update. Instal Claude Code dulu.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 07eadc0769c8
---

Jalankan ini di terminal pada Linux atau macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skrip ini mengunduh rilis untuk sistem Anda, memeriksa checksum SHA-256-nya, dan menaruh `ccam` di `~/.local/bin`, tempat ccam nanti juga membuat perintah pintasan yang Anda pilih, seperti `ccwork`. Jika folder tersebut tidak ada di `PATH` Anda, skrip akan menampilkan baris yang perlu ditambahkan ke `~/.bashrc` atau `~/.zshrc`.

Di Windows, jalankan ini di PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Skrip ini mengunduh `ccam_windows_<arch>.zip`, memeriksa checksum SHA-256-nya, menaruh `ccam.exe` di `%USERPROFILE%\.local\bin`, folder tempat installer Claude Code menaruh `claude.exe`, dan menambahkan folder itu ke `PATH` pengguna Anda jika belum ada. Setelah itu, buka terminal baru.

ccam menjalankan Claude Code untuk Anda, jadi instal [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) terlebih dahulu dan pastikan `claude --version` berfungsi.

## Sistem yang didukung

- Linux di x86-64 dan ARM64: Ubuntu, Debian, Fedora, Arch, Alpine, dan distribusi lainnya. Binary-nya statis dan tidak memiliki dependensi.
- macOS di Apple silicon dan Intel.
- Windows 10 dan 11 di x86-64 dan ARM64, di PowerShell, Windows Terminal, atau Git Bash. Folder bersama berfungsi untuk semua pengguna. File bersama (`settings.json`, `CLAUDE.md`, `history.jsonl`) memerlukan Mode Pengembang (Settings > System > For developers), karena Windows hanya mengizinkan mode itu, atau administrator, untuk membuat tautan ke file; lihat [berbagi di Windows](/features/share#on-windows).
- Di dalam WSL, ccam berfungsi seperti di sistem Linux mana pun.

## Cara instalasi lainnya

**Paket Debian dan Ubuntu.** Unduh `ccam_<version>_linux_amd64.deb` dari [halaman rilis](https://github.com/jhyunwoo/ccam/releases) lalu instal:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Paket ini menginstal `/usr/bin/ccam`.

**Dari kode sumber.** Dengan Go 1.26 atau yang lebih baru:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Versi tertentu.** Tetapkan `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Pelengkapan otomatis shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Pelengkapan otomatis mencakup semua perintah dan opsi `ccam`, termasuk nama akun.

## Periksa instalasi

```sh
ccam version
ccam doctor
```

`ccam doctor` memeriksa bahwa Claude Code ditemukan dan bahwa setiap folder akun, tautan, dan login dalam kondisi baik. Lanjutkan ke [panduan memulai cepat](/docs/quick-start).

## Pembaruan

```sh
ccam update --check
ccam update
```

`ccam update --check` hanya memberi tahu apakah ada rilis yang lebih baru. `ccam update` mengunduh rilis untuk sistem Anda, mencocokkan checksum SHA-256-nya dengan `checksums.txt` dari rilis tersebut, lalu mengganti file `ccam` dalam satu langkah. Jika checksum tidak cocok, tidak ada yang diubah. Pintasan seperti `ccwork` adalah tautan ke ccam, jadi pintasan itu langsung menjalankan versi baru.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam tidak pernah memeriksa pembaruan dengan sendirinya; ccam hanya mencari rilis baru saat Anda menjalankan perintah ini. Untuk menginstal rilis tertentu, termasuk rilis yang lebih lama, berikan `--version`. `ccam --dry-run update` menampilkan apa yang akan diunduh dan diganti tanpa mengubah apa pun.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Menjalankan skrip instalasi lagi juga memperbarui salinan di `~/.local/bin`. Jika Anda membangun ccam dengan `go install`, jalankan lagi perintah yang sama.

**Windows.** `ccam update` mengunduh `.zip` untuk Windows. Program yang sedang berjalan tidak bisa diganti di Windows, jadi binary lama diganti namanya menjadi `ccam.exe.old`, yang akan dihapus oleh pembaruan berikutnya.

**Paket Debian dan Ubuntu.** `ccam update` tidak mengganti `/usr/bin/ccam` yang diinstal oleh apt, dan sebagai gantinya menampilkan langkah-langkah berikut. Unduh `ccam_<version>_linux_amd64.deb` yang baru dari [halaman rilis](https://github.com/jhyunwoo/ccam/releases) lalu instal di atas versi yang lama:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Menghapus instalasi

Hapus dulu pintasan akun selagi ccam masih terpasang. `ccam list` menampilkannya di kolom Shortcut; hapus masing-masing dengan `ccam shortcut <account> --remove`, atau hapus sendiri tautan itu dari `~/.local/bin`. Lalu hapus perintahnya, serta folder data jika Anda tidak lagi memerlukan akun-akun tersebut:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Di Windows, di PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Menghapus `~/.local/share/ccam` juga menghapus login dan riwayat percakapan yang tersimpan di sana. Folder `~/.claude` yang Anda [impor](/docs/commands) tetap berada di tempatnya.
