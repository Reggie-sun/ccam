---
title: ccam dibandingkan dengan alias shell
description: Perbandingan alias shell dengan CLAUDE_CONFIG_DIR dan ccam untuk menjalankan beberapa akun Claude Code, fitur demi fitur.
updated: 2026-10-26
order: 1
translation: machine
sourceHash: "c30d84bceba0"
---

Satu alias shell per akun adalah cara paling sederhana untuk menjalankan beberapa akun Claude Code. ccam melakukan hal yang sama di baliknya, lalu menambahkan pembagian, perpindahan akun, dan pemantauan penggunaan.

## Penyiapan alias

Tambahkan satu baris per akun tambahan ke `~/.bashrc` atau `~/.zshrc`, buka terminal baru, lalu jalankan alias itu sekali untuk login:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` tetap menjadi akun pertama Anda dan `claude-work` menjalankan akun kedua. Inilah metode yang [dijelaskan dalam dokumentasi Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## Penyiapan ccam

Instal ccam dan tambahkan akunnya, di sini dengan perintah pintasan tersendiri:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` dan `ccwork` menjalankan akun tersebut, dan keduanya juga berfungsi di skrip.

## Berdampingan

| | Alias shell | ccam |
| --- | --- | --- |
| Login terpisah | Ya | Ya |
| Penyiapan per akun | Edit `~/.bashrc`, login | `ccam add` |
| Berfungsi di skrip dan program lain | Tidak, alias hanya ada di shell | Ya, `ccam work` dan pintasan seperti `ccwork` adalah perintah sungguhan |
| Pengaturan, `CLAUDE.md`, ekstensi bersama | Salin atau tautkan secara manual | Ya, dipilih per item, dengan cadangan |
| Server MCP bersama | Manual, per akun | `ccam mcp add` sekali saja |
| Melanjutkan percakapan di akun lain | Tidak | `ccam next --continue` |
| Penggunaan semua akun | Tidak | `ccam usage` |
| Memilih akun dengan sisa kuota terbanyak | Tidak | `ccam next --auto` |
| Menjalankan ulang tugas `-p` sebagai akun lain saat batas tercapai | Tidak | `ccam --fallback` |
| Satu akun untuk setiap folder proyek | Tidak | `ccam folder` |
| Nama akun di status line | Tidak | `ccam statusline install` |
| Memperbaiki tautan yang rusak | Manual | `ccam doctor --fix` |
| Perangkat lunak tambahan | Tidak ada | Satu binary statis |

Pilih alias jika Anda punya dua akun dan tidak ingin ada yang dibagikan di antaranya. Pilih ccam jika Anda ingin akun-akun itu terasa seperti satu penyiapan. Anda bisa mulai dengan alias dan beralih ke ccam nanti: `ccam import` mendaftarkan folder `~/.claude-*` yang sudah ada di tempatnya.

[Cara menggunakan beberapa akun di satu komputer](/guides/multiple-claude-code-accounts) menjelaskan kedua pendekatan langkah demi langkah, dan [ccam dan pengalih akun lainnya](/compare/account-switchers) membandingkan ccam dengan alat yang menukar satu login.
