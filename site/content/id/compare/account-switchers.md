---
title: ccam dan pengalih akun Claude Code lainnya
description: Perbandingan ccam dengan claude-swap, CCS, clauth, CAAM, cux, dan alat lain untuk ganti akun Claude Code, serta pendekatan mana yang cocok untuk kebutuhan apa.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Alat untuk berganti akun Claude Code memakai salah satu dari dua pendekatan. Sebagian besar menukar login di dalam satu folder `~/.claude`, sehingga hanya satu akun yang aktif pada satu waktu. Yang lain, termasuk ccam, memberi setiap akun foldernya sendiri dengan `CLAUDE_CONFIG_DIR`, sehingga setiap akun tetap login dan beberapa akun bisa berjalan bersamaan. ccam menambahkan pengaturan dan server MCP bersama, serta bisa melanjutkan percakapan di folder akun lain.

## Menukar satu login atau folder terpisah?

**Menukar login** mempertahankan satu set pengaturan, server MCP, dan riwayat, karena hanya ada satu folder. Pergantiannya cepat, dan sesi yang sedang berjalan bisa memakai login yang baru. Namun, hanya satu akun yang aktif pada satu waktu, untuk semua terminal dan editor sekaligus.

**Folder terpisah** adalah metode yang [dijelaskan dalam dokumentasi Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Setiap akun tetap login, dan Anda bisa menjalankan sesi kerja dan sesi pribadi berdampingan. Konsekuensinya, setiap folder dimulai dengan pengaturan, server MCP, dan riwayatnya sendiri, kecuali ada alat yang membagikannya.

## Perbandingan alat-alatnya

Diambil dari dokumentasi masing-masing proyek pada 2026-10-19. "Tidak didokumentasikan" berarti proyek tersebut tidak menyebutkannya, bukan berarti fiturnya tidak ada.

| Alat | Pendekatan | Beberapa akun sekaligus | Penggunaan 5 jam dan mingguan | Saat batas tercapai | Platform |
| --- | --- | --- | --- | --- | --- |
| ccam | Folder terpisah | Ya | Ya, `ccam usage` | Anda menjalankan `ccam next --auto --continue` untuk melanjutkan percakapan di akun dengan sisa kuota terbanyak | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Menukar login; folder terpisah dengan `cswap run` | Ya, dengan `cswap run` | Ya | `cswap auto` menukar login pada ambang tertentu; sesi yang berjalan tetap berlanjut | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Folder terpisah | Ya | Sebagian, di aplikasi bilah menu macOS | Kumpulan akun meneruskan permintaan ke akun berikutnya melalui proxy lokal | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Menukar login; folder terpisah dengan `clauth start` | Ya, dengan `clauth start` | Ya | Berpindah mengikuti rantai cadangan | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Menukar login, atau profil terisolasi | Ya, dengan profil terisolasi | Ya | `caam run` mengulang perintah non-interaktif di akun lain | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Menukar login dan membungkus `claude` | Tidak didokumentasikan | Ya | Menukar login dan melanjutkan percakapan yang sama | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Menukar login | Tidak | Ya | Tidak didokumentasikan | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Menukar login | Tidak | Ya | Peralihan opsional pada 100% | macOS |

ccswap adalah fork dari claude-swap yang juga mengelola login Codex. Beberapa alat ini juga menangani agen coding lain, seperti Codex atau Gemini.

## Fokus ccam

- Setiap akun tetap login. Berganti akun tidak pernah membuat akun lain logout, dan `ccam work` serta `ccam personal` bisa berjalan di dua terminal sekaligus.
- Server MCP, `settings.json`, `CLAUDE.md`, agen, skill, dan plugin dibagikan antarfolder melalui tautan, per item dan per akun. Login tidak pernah dibagikan.
- Dengan riwayat bersama, satu perintah, `ccam next --auto --continue`, membuka percakapan yang sama di folder akun dengan sisa kuota terbanyak.
- ccam menjalankan Claude Code lalu menggantikan dirinya dengan Claude Code, jadi tidak ada yang berjalan di latar belakang. ccam tidak memantau penggunaan Anda dan tidak berganti akun sendiri; Anda yang menentukan kapan.

## Di mana alat lain lebih unggul

- claude-swap, CCS, clauth, CAAM, dan cux berjalan di Windows tanpa WSL. ccam tidak.
- claude-swap, clauth, dan cux berganti akun sendiri saat sebuah akun melewati ambang penggunaan. ccam menyerahkan pergantian kepada Anda.
- CCSwitcher, Claude Switcher, dan claude-swap menampilkan penggunaan di bilah menu macOS.
- CCS, CAAM, ccswap, dan Claude Switcher juga mengelola Codex dan agen coding lainnya.

## Mana yang cocok?

Pilih alat yang menukar login jika Anda memakai satu akun pada satu waktu, ingin pergantian terjadi otomatis, atau memakai Windows tanpa WSL. Pilih ccam jika Anda ingin akun berjalan berdampingan, satu penyiapan bersama, dan percakapan yang sama di akun berikutnya. Jika tidak ada yang perlu dibagikan, [alias shell biasa](/compare/aliases) sudah cukup.

Boleh tidaknya Anda menggunakan beberapa akun bergantung pada ketentuan setiap akun. Periksa ketentuan itu sebelum melakukannya.
