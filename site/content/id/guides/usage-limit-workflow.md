---
title: Saat Claude Code mencapai batas penggunaan
description: Tunggu reset atau langsung lanjutkan percakapan yang sama di akun Claude Code lain. Cara kerja batas 5 jam dan mingguan, langkah demi langkah.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: "38b80d780bba"
---

Saat sebuah akun mencapai batas penggunaannya (sering disebut rate limit), Anda punya dua pilihan: menunggu sampai batasnya direset, atau melanjutkan percakapan yang sama sekarang juga di akun lain milik Anda. Dengan ccam, pilihan kedua cukup satu perintah, `ccam next --auto --continue`, setelah riwayat percakapan dibagikan antarakun.

## Batas apa saja yang dimiliki Claude Code?

Pada paket Pro dan Max, Claude Code memiliki [batas sesi lima jam dan batas mingguan](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Keduanya dipakai bersama dengan Claude di web dan di aplikasi, dan penggunaan dihitung terhadap keduanya sekaligus. Saat salah satunya habis, Claude Code [menampilkan pesan](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) seperti "You've hit your session limit · resets 3:45pm". Berganti model tidak membantu, karena batasnya berlaku untuk semua model.

## Opsi 1: tunggu reset

Jika Anda hanya punya satu akun, tunggu. Versi terbaru Claude Code [menunggu di sesi yang terbuka](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) dan melanjutkan tugas dengan sendirinya setelah batasnya direset. `/usage` di dalam Claude Code menampilkan seberapa banyak batas paket Anda yang sudah terpakai.

## Opsi 2: lanjutkan di akun lain sekarang

Jika Anda punya akun kedua, Anda bisa langsung melanjutkan percakapan yang sama. Siapkan sekali:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Status line kini menampilkan akun dan penggunaan 5 jamnya selama Anda bekerja, sehingga batas tidak datang tiba-tiba.

Saat batas tercapai, keluar dari Claude Code (`/exit`), lalu:

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

Claude Code dimulai sebagai `work` dengan percakapan terbaru dari folder saat ini, termasuk semua yang sudah dibicarakan. Berbagi riwayat (`ccam group set default --add history` di atas) menjaganya tetap satu percakapan; tanpa itu, Claude Code menyalin percakapan ke `work`, dan ccam memberi tahu berapa token yang dibaca ulang oleh balasan pertama.

Dengan [plugin ccam](/docs/install#claude-code-plugin) di Claude Code, sesi yang berhenti karena batasnya menampilkan notifikasi, dan `/ccam:switch` (atau `/ccam:switch work`) mencatat percakapan itu untuk `ccam next --continue` sehingga tidak ada yang perlu ditebak.

## Apakah ccam berganti akun secara otomatis?

Tidak. ccam tidak berjalan selama Anda bekerja: ccam menjalankan Claude Code lalu menggantikan dirinya dengan Claude Code. Anda yang menentukan kapan berganti, dan `ccam next --auto` memilih akun dengan penggunaan 5 jam terendah serta melewati akun yang sudah 100%, belum login, atau berasal dari organisasi lain. Dengan `relay = true` di `config.toml`, ccam justru menunggu di samping sesi, sehingga `/ccam:switch` melanjutkan percakapan dengan sendirinya. Beberapa alat lain justru menukar login sesi yang sedang berjalan; [ccam dan pengalih akun lainnya](/compare/account-switchers) menjelaskan perbedaannya.

## Hal yang perlu diperhatikan

- Jalankan `ccam next --continue` di folder yang sama tempat Anda bekerja; percakapan disimpan per folder proyek.
- Hindari menjalankan dua akun pada percakapan yang sama secara bersamaan. ccam memperingatkan Anda jika Anda mencobanya.
- Tanpa riwayat bersama, percakapan disalin ke akun lain dan kedua salinan itu berjalan sendiri-sendiri; checkpoint yang dipakai `/rewind` untuk memulihkan kode tetap berada di akun pertama.
- Pesan pertama di akun baru bisa memakai lebih banyak batasnya daripada biasanya, karena akun itu belum menyimpan percakapan tersebut di cache.
- Boleh tidaknya Anda berganti akun seperti ini bergantung pada ketentuan setiap akun.

Lihat [batas penggunaan](/docs/usage-limits) untuk setiap opsi `ccam usage` dan `ccam next`.
