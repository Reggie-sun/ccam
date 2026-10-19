---
title: Instal ccam
description: Instal ccam di Linux atau macOS dengan satu perintah, dari paket .deb, atau dari kode sumber, lalu jaga agar tetap terbaru dengan ccam update.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 4d7291e0f86e
---

Jalankan ini di terminal pada Linux atau macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Skrip ini mengunduh rilis untuk sistem Anda, memeriksa checksum SHA-256-nya, dan menaruh `ccam` di `~/.local/bin`, tempat ccam nanti juga membuat perintah pintasan yang Anda pilih, seperti `ccwork`. Jika folder tersebut tidak ada di `PATH` Anda, skrip akan menampilkan baris yang perlu ditambahkan ke `~/.bashrc` atau `~/.zshrc`.

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

Paket ini menginstal `/usr/bin/ccam`.

**Dari kode sumber.** Dengan Go 1.26 atau yang lebih baru:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
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

Menghapus `~/.local/share/ccam` juga menghapus login dan riwayat percakapan yang tersimpan di sana. Folder `~/.claude` yang Anda [impor](/docs/commands) tetap berada di tempatnya.
