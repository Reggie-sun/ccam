---
title: Cara menggunakan dua akun Claude Code di satu komputer
description: Jalankan beberapa akun Claude Code secara berdampingan dengan CLAUDE_CONFIG_DIR, secara manual dengan alias shell atau otomatis dengan ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d82129ddcb4c
---

Untuk menggunakan dua akun Claude Code di satu komputer, beri setiap akun folder konfigurasinya sendiri. Claude Code membaca folder tersebut dari variabel lingkungan `CLAUDE_CONFIG_DIR`, dan setiap folder menyimpan login, pengaturan, dan riwayatnya sendiri. Anda bisa menyiapkannya secara manual dengan alias shell, atau membiarkan ccam melakukannya.

## Secara manual, dengan alias

Ini adalah pendekatan yang dijelaskan dalam dokumentasi Claude Code. Tambahkan satu baris per akun tambahan ke `~/.bashrc` atau `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Buka terminal baru dan jalankan `claude-work`. Claude Code meminta Anda login, karena foldernya masih baru. Sejak saat itu, `claude` adalah akun pertama Anda dan `claude-work` akun kedua.

Cara ini berfungsi baik untuk dua akun. Hal-hal yang perlu diketahui:

- Setiap folder dimulai dalam keadaan kosong: tanpa pengaturan, tanpa `CLAUDE.md`, tanpa server MCP, tanpa agen atau skill. Anda menyalin atau menautkan apa yang Anda perlukan, dan menyelaraskan salinannya sendiri.
- Alias hanya ada di shell interaktif. Skrip, Makefile, dan program lain tidak melihatnya.
- Server MCP berada di `.claude.json`, bersama dengan login, sehingga tidak bisa begitu saja ditautkan antarfolder.
- Percakapan tetap berada di akun tempat percakapan itu dimulai, sehingga Anda tidak bisa melanjutkannya di akun lain.
- Tidak ada yang menampilkan penggunaan semua akun sekaligus.

## Dengan ccam

ccam mengotomatiskan ide yang sama dan menambahkan bagian yang kurang:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Sekarang `claude1` dan `claude2`, atau `ccs personal` dan `ccs work`, menjalankan kedua akun tersebut. Karena keduanya adalah program sungguhan, bukan alias, keduanya juga berfungsi di skrip. Pengaturan, `CLAUDE.md`, server MCP, dan ekstensi dibagikan sejak awal, dan Anda dapat mengaktifkan riwayat percakapan bersama untuk melanjutkan percakapan di akun lain dengan `ccs next --continue`.

Jika Anda sudah menggunakan `~/.claude`, jalankan `ccam import` alih-alih menambahkan akun itu lagi.

## Mana yang dipilih

Untuk dua akun tanpa ada yang dibagikan, alias sudah cukup. Saat Anda menginginkan pengaturan dan server MCP bersama, lebih banyak akun, penggunaan yang terlihat sekilas, atau melanjutkan percakapan, ccam menghemat pekerjaan administrasinya. [Perbandingan](/compare/aliases) mencantumkan perbedaannya secara berdampingan.
