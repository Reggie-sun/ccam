---
title: Pertanyaan yang sering diajukan
description: Jawaban seputar ccam, akun Claude Code, batas penggunaan, pembagian, privasi, dan sistem yang didukung.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: e07d468d453b
---

## Apakah ccam alat resmi dari Anthropic?

Tidak. ccam adalah proyek open source tidak resmi dan tidak berafiliasi dengan atau didukung oleh Anthropic. ccam menggunakan pengaturan `CLAUDE_CONFIG_DIR` yang didokumentasikan Claude Code untuk menjalankan beberapa akun secara berdampingan.

## Apakah saya boleh menggunakan beberapa akun Claude Code?

Itu bergantung pada ketentuan setiap akun dan, untuk akun kerja, kebijakan organisasi Anda. Periksa hal tersebut sebelum Anda menggunakan beberapa akun. ccam hanya mempermudah perpindahan di antara akun yang sudah Anda miliki.

## Apakah ccam mengakali batas penggunaan?

Tidak. Setiap akun tetap memiliki batasnya sendiri. ccam menampilkan batas tersebut dan memungkinkan Anda melanjutkan di akun lain milik Anda saat satu akun habis kuotanya.

## Bagaimana cara menggunakan dua akun Claude Code di satu komputer?

Beri setiap akun folder konfigurasinya sendiri dengan `CLAUDE_CONFIG_DIR`. ccam melakukannya untuk Anda: jalankan `ccam add` untuk setiap akun, lalu jalankan akun-akun itu dengan `ccam personal` dan `ccam work`, atau dengan pintasan yang namanya Anda pilih sendiri, seperti `ccwork`. [Panduan ini](/guides/multiple-claude-code-accounts) juga menunjukkan cara melakukannya secara manual.

## Bagaimana cara melanjutkan percakapan di akun lain?

Bagikan riwayat percakapan antara kedua akun dengan `ccam share <name> --history on`, lalu jalankan `ccam next --continue`. Claude Code membuka percakapan terbaru dari folder saat ini di akun berikutnya.

## Apakah ccam melihat atau mengirim login saya?

ccam tidak pernah menyalin atau membagikan login. ccam membaca token akses sebuah akun hanya untuk menanyakan penggunaan akun itu ke endpoint penggunaan Anthropic, saat Anda menjalankan `ccam usage`, `ccam next --auto`, atau `ccam doctor`. ccam tidak memiliki telemetri dan tidak mengirim apa pun ke tempat lain.

## Di mana akun saya disimpan?

Di `~/.local/share/ccam/profiles/<name>/`, satu folder per akun. Di macOS, Claude Code menyimpan login itu sendiri di Rantai Kunci (Keychain), dengan nama yang diturunkan dari folder tersebut.

## Apa yang terjadi pada ~/.claude saya yang sudah ada?

Tidak ada, kecuali Anda menjalankan `ccam import` atau `ccam migrate`. Impor mendaftarkan folder itu di tempatnya, sehingga `claude` biasa dan `ccam <name>` tetap menggunakan akun dan pengaturan yang sama. `ccam migrate` menyalin pengaturan, `CLAUDE.md`, ekstensi, dan server MCP dari folder itu ke pengaturan bersama dan hanya membaca folder tersebut.

## Sistem apa saja yang didukung ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine, dan lainnya) serta macOS. Di Windows, gunakan ccam di dalam WSL.

## Mengapa server MCP hanya muncul di satu akun?

Server yang ditambahkan di dalam sesi dengan `claude mcp add` hanya disimpan untuk akun tersebut. ccam memberi tahu Anda saat berikutnya Anda menjalankan akun itu; pindahkan server ke grup bersama dengan `ccam mcp import <account> <server>`.

## Apakah ccam memperlambat Claude Code?

Menjalankan Claude Code melalui ccam menambah sekitar 4 milidetik, diukur dengan `scripts/perf-check.sh` di sebuah laptop. ccam membaca satu file kecil, memeriksa tautan bersama, lalu menggantikan dirinya dengan `claude`, sehingga ccam tidak berjalan saat Anda bekerja.

## Bagaimana cara menghapus instalasi ccam?

Hapus pintasan setiap akun dengan `ccam shortcut <account> --remove`, hapus `ccam` dari `~/.local/bin`, lalu hapus `~/.config/ccam`, `~/.local/share/ccam`, dan `~/.local/state/ccam` jika Anda tidak lagi memerlukan akun-akun tersebut. Lihat [instalasi](/docs/install#uninstall).
