---
title: Memulai cepat
description: Tambahkan dua akun Claude Code dengan ccam add, berpindah di antaranya dengan claude1, claude2, atau ccam <name>, dan pilih apa yang dibagikan.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

Tambahkan setiap akun dengan `ccam add`, lalu jalankan Claude Code dengan `claude1`, `claude2`, atau `ccam <name>`. Hanya itu yang dibutuhkan sebagian besar orang.

## 1. Tambahkan akun pertama Anda

```sh
ccam add
```

Wizard akan menanyakan nama, seperti `personal` atau `work`, dan item apa saja yang perlu dibagikan akun ini dengan akun Anda yang lain. Lalu wizard menjalankan `claude auth login`, yang membuka browser agar Anda bisa masuk. Setelah selesai, Anda akan melihat:

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam menampilkan pesannya sendiri dalam bahasa Inggris; perintah dan opsinya sama dalam semua bahasa.

## 2. Tambahkan akun kedua

```sh
ccam add work
```

Masuk dengan akun lainnya di browser. Akun ini menjadi `claude2`.

## 3. Berpindah akun

Semua perintah berikut menjalankan Claude Code sebagai akun kerja:

```sh
claude2
ccam work
ccam 2
```

Semua yang ada setelah nama akun diteruskan ke Claude Code tanpa diubah, jadi Anda tetap memakai opsi yang sudah Anda kenal:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Menjalankan `ccam` saja akan membuka pemilih akun: gunakan tombol panah dan Enter, atau tekan nomor akun.

## 4. Pilih apa yang dibagikan

Secara default, server MCP, `settings.json`, `CLAUDE.md` global, serta agen, perintah, skill, dan plugin dibagikan antarakun. Riwayat percakapan tidak dibagikan. Aktifkan untuk kedua akun jika Anda ingin melanjutkan percakapan di antara keduanya:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam menampilkan apa yang akan dipindahkannya sebelum mengubah apa pun. Lihat [pembagian](/docs/sharing) untuk detailnya.

## 5. Saat sebuah akun mencapai batasnya

```sh
ccam usage
ccam next --continue
```

`ccam usage` menampilkan penggunaan 5 jam dan mingguan semua akun. `ccam next --continue` membuka percakapan yang sama di akun berikutnya. Lihat [batas penggunaan](/docs/usage-limits).

## Sudah menggunakan ~/.claude?

Daftarkan folder itu alih-alih masuk lagi:

```sh
ccam import
```

Folder tersebut tetap berada di tempatnya, dan `claude` biasa tetap berfungsi seperti sebelumnya.
