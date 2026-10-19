---
title: Yang perlu dilakukan saat Claude Code mencapai batas penggunaan
description: Periksa sisa kuota setiap akun Claude Code dan lanjutkan percakapan yang sama di akun lain saat satu akun mencapai batasnya.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

Saat satu akun mencapai batas penggunaan Claude Code, lanjutkan percakapan yang sama di akun lain milik Anda: bagikan riwayat percakapan sekali, lalu jalankan `ccam next --continue`. Jika Anda hanya punya satu akun, tunggu hingga waktu reset yang ditampilkan Claude Code.

## Sebelum terjadi

Siapkan akun-akun Anda agar perpindahan cukup dengan satu perintah:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Sekarang status line menampilkan akun dan penggunaan 5 jamnya saat Anda bekerja, sehingga batas penggunaan tidak datang sebagai kejutan.

## Saat batas tercapai

Keluar dari Claude Code (`/exit`), lalu:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Buka percakapan di akun dengan sisa kuota terbanyak:

```sh
ccam next --auto --continue
```

Claude Code dimulai sebagai `work` dengan percakapan terbaru dari folder saat ini, termasuk semua yang sudah dibicarakan sejauh ini.

## Hal yang perlu diingat

- Jalankan `ccam next --continue` di folder yang sama tempat Anda bekerja; percakapan disimpan per folder proyek.
- Hindari menjalankan dua akun pada percakapan yang sama secara bersamaan. ccam memperingatkan Anda jika Anda mencobanya.
- Akun lainnya juga harus berbagi riwayat percakapan, dalam grup yang sama.
- Boleh tidaknya Anda berpindah akun seperti ini bergantung pada ketentuan setiap akun.

Lihat [batas penggunaan](/docs/usage-limits) untuk semua opsi `ccam usage` dan `ccam next`.
