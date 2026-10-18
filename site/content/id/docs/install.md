---
title: Instal ccam
description: Instal ccam di Linux atau macOS dengan satu perintah, dari paket .deb, atau dari kode sumber. Claude Code harus diinstal terlebih dahulu.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Jalankan ini di terminal pada Linux atau macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skrip ini mengunduh rilis untuk sistem Anda, memeriksa checksum SHA-256-nya, dan menaruh dua perintah di `~/.local/bin`: `ccam` dan alias pendeknya `ccs`. Jika folder tersebut tidak ada di `PATH` Anda, skrip akan menampilkan baris yang perlu ditambahkan ke `~/.bashrc` atau `~/.zshrc`.

ccam menjalankan Claude Code untuk Anda, jadi instal [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) terlebih dahulu dan pastikan `claude --version` berfungsi.

## Sistem yang didukung

- Linux di x86-64 dan ARM64: Ubuntu, Debian, Fedora, Arch, Alpine, dan distribusi lainnya. Binary-nya statis dan tidak memiliki dependensi.
- macOS di Apple silicon dan Intel.
- Windows tidak didukung secara langsung. Di dalam WSL, ccam berfungsi seperti di sistem Linux mana pun.

## Cara instalasi lainnya

**Paket Debian dan Ubuntu.** Unduh `ccam_<version>_linux_amd64.deb` dari [halaman rilis](https://github.com/jhyunwoo/ccam/releases) lalu instal:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Paket ini menginstal `/usr/bin/ccam` dan tautan `/usr/bin/ccs`.

**Dari kode sumber.** Dengan Go 1.26 atau yang lebih baru:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**Versi tertentu.** Tetapkan `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Pelengkapan otomatis shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Pelengkapan otomatis mencakup `ccam` dan `ccs`, termasuk nama akun.

## Periksa instalasi

```sh
ccam version
ccam doctor
```

`ccam doctor` memeriksa bahwa Claude Code ditemukan dan bahwa setiap folder akun, tautan, dan login dalam kondisi baik. Lanjutkan ke [panduan memulai cepat](/docs/quick-start).

## Menghapus instalasi

Hapus kedua perintah dan pintasan akun, lalu folder data jika Anda tidak lagi memerlukan akun-akun tersebut:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Menghapus `~/.local/share/ccam` juga menghapus login dan riwayat percakapan yang tersimpan di sana. Folder `~/.claude` yang Anda [impor](/docs/commands) tetap berada di tempatnya.
