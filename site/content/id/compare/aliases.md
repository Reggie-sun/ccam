---
title: ccam dibandingkan dengan alias shell
description: Perbandingan alias shell dengan CLAUDE_CONFIG_DIR dan ccam untuk menjalankan beberapa akun Claude Code, berdampingan.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6fc4200a8a65
---

Satu alias shell per akun adalah cara paling sederhana untuk menjalankan beberapa akun Claude Code. ccam melakukan hal yang sama di baliknya, lalu menambahkan pembagian, perpindahan akun, dan pemantauan penggunaan.

| | Alias shell | ccam |
| --- | --- | --- |
| Login terpisah | Ya | Ya |
| Penyiapan per akun | Edit `~/.bashrc`, login | `ccam add` |
| Berfungsi di skrip dan program lain | Tidak, alias hanya ada di shell | Ya, `ccam work` dan pintasan seperti `ccwork` adalah perintah sungguhan |
| Pengaturan, `CLAUDE.md`, ekstensi bersama | Salin atau tautkan secara manual | Per item, dengan cadangan |
| Server MCP bersama | Manual, per akun | `ccam mcp add` sekali saja |
| Melanjutkan percakapan di akun lain | Tidak | `ccam next --continue` |
| Penggunaan semua akun | Tidak | `ccam usage` |
| Memilih akun dengan sisa kuota terbanyak | Tidak | `ccam next --auto` |
| Nama akun di status line | Tidak | `ccam statusline install` |
| Memperbaiki tautan yang rusak | Manual | `ccam doctor --fix` |
| Perangkat lunak tambahan | Tidak ada | Satu binary statis |

Pilih alias jika Anda punya dua akun dan tidak ingin ada yang dibagikan di antaranya. Pilih ccam jika Anda ingin akun-akun itu terasa seperti satu penyiapan. Anda bisa mulai dengan alias dan beralih ke ccam nanti: `ccam import` mendaftarkan folder `~/.claude-*` yang sudah ada di tempatnya.

[Cara menggunakan beberapa akun di satu komputer](/guides/multiple-claude-code-accounts) menjelaskan kedua pendekatan langkah demi langkah, dan [ccam dan pengalih akun lainnya](/compare/account-switchers) membandingkan ccam dengan alat yang menukar satu login.
