---
title: Memperkenalkan ccam
description: Mengapa ccam dibuat, apa gunanya bagi pengguna beberapa akun Claude Code, dan bagaimana ccam tidak mengganggu pekerjaan Anda.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: b116c91bb8d4
---

ccam adalah alat baris perintah kecil untuk orang yang menggunakan lebih dari satu akun Claude Code di komputer yang sama: akun pribadi dan akun kerja, atau beberapa akun untuk satu proyek panjang.

## Masalahnya

Claude Code mendukung beberapa akun melalui `CLAUDE_CONFIG_DIR`: satu folder per akun. Dalam praktiknya, setiap folder berdiri sendiri. Pengaturan, `CLAUDE.md`, server MCP, agen, dan skill harus disiapkan ulang di setiap folder dan diselaraskan secara manual. Percakapan tetap berada di tempat percakapan itu dimulai. Dan saat satu akun mencapai batas penggunaannya, Anda baru tahu setelah mencobanya.

## Apa yang dilakukan ccam

- `ccam add` membuat folder, menjalankan login, dan jika Anda mau, memberi akun itu pintasan yang namanya Anda pilih sendiri, seperti `ccwork`.
- `ccam work` dan `ccwork` menjalankan Claude Code sebagai akun tersebut dan meneruskan semua argumen tanpa diubah.
- Server MCP, pengaturan, `CLAUDE.md`, dan ekstensi dibagikan melalui tautan, per item dan per akun. Login tidak pernah dibagikan.
- `ccam usage` menampilkan penggunaan 5 jam dan mingguan semua akun; `ccam next --continue` melanjutkan percakapan di akun berikutnya.

## Tidak mengganggu

Alat yang Anda jalankan sebelum setiap sesi seharusnya tidak memakan biaya yang terasa. ccam membaca satu file konfigurasi, memeriksa setiap tautan bersama dengan satu system call, lalu menggantikan dirinya dengan `claude`. ccam hanya menambah sekitar 4 milidetik, dan tidak berjalan saat Anda bekerja.

ccam juga berhati-hati dengan file Anda. Setiap perubahan pada pembagian menampilkan rencana sebelum dimulai dan menyimpan file asli di folder cadangan. `ccam doctor` memeriksa tautan, izin, dan login, lalu memperbaiki apa yang aman untuk diperbaiki.

## Cobalah

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam bersifat open source dan tidak resmi: ccam tidak berafiliasi dengan atau didukung oleh Anthropic. Boleh tidaknya Anda menggunakan beberapa akun bergantung pada ketentuan setiap akun. Masukan dan terjemahan sangat diterima di [GitHub](https://github.com/jhyunwoo/ccam).
