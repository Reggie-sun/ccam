---
title: ccam dibandingkan dengan alias shell
description: Perbandingan alias shell dengan CLAUDE_CONFIG_DIR dan ccam untuk menjalankan beberapa akun Claude Code, berdampingan.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e9902676e0ed
---

Satu alias shell per akun adalah cara paling sederhana untuk menjalankan beberapa akun Claude Code. ccam melakukan hal yang sama di baliknya, lalu menambahkan pembagian, perpindahan akun, dan pemantauan penggunaan.

| | Alias shell | ccam |
| --- | --- | --- |
| Login terpisah | Ya | Ya |
| Penyiapan per akun | Edit `~/.bashrc`, login | `ccam add` |
| Berfungsi di skrip dan program lain | Tidak, alias hanya ada di shell | Ya, `claude2` adalah perintah sungguhan |
| Pengaturan, `CLAUDE.md`, ekstensi bersama | Salin atau tautkan secara manual | Per item, dengan cadangan |
| Server MCP bersama | Manual, per akun | `ccam mcp add` sekali saja |
| Melanjutkan percakapan di akun lain | Tidak | `ccs next --continue` |
| Penggunaan semua akun | Tidak | `ccs usage` |
| Memilih akun dengan sisa kuota terbanyak | Tidak | `ccs next --auto` |
| Nama akun di status line | Tidak | `ccam statusline install` |
| Memperbaiki tautan yang rusak | Manual | `ccam doctor --fix` |
| Perangkat lunak tambahan | Tidak ada | Satu binary statis |

Pilih alias jika Anda punya dua akun dan tidak ingin ada yang dibagikan di antaranya. Pilih ccam jika Anda ingin akun-akun itu terasa seperti satu penyiapan. Anda bisa mulai dengan alias dan beralih ke ccam nanti: `ccam import` mendaftarkan folder `~/.claude-*` yang sudah ada di tempatnya.

[Cara menggunakan dua akun di satu komputer](/guides/multiple-claude-code-accounts) menjelaskan kedua pendekatan langkah demi langkah.
