---
title: Memulai cepat
metaTitle: "Memulai cepat: tambah dan ganti akun Claude Code"
description: Tambahkan dua akun Claude Code dengan ccam add, beri pintasan seperti ccwork, berpindah di antaranya, dan pilih apa yang dibagikan.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "bd893bbb83ab"
---

Tambahkan setiap akun dengan `ccam add`, lalu jalankan Claude Code dengan `ccam <name>` atau dengan pintasan yang namanya Anda pilih sendiri, seperti `ccwork`. Bagi kebanyakan orang, itu sudah cukup.

Sudah memakai Claude Code? Baca [Sudah menggunakan ~/.claude?](#already-using-claude) di akhir halaman ini sebelum Anda mulai.

## 1. Tambahkan akun pertama Anda

```sh
ccam add
```

Wizard akan menanyakan nama, seperti `personal` atau `work`, item apa saja yang perlu dibagikan akun ini dengan akun Anda yang lain, dan apakah perlu membuat perintah pintasan (lihat langkah 3). Lalu wizard menjalankan `claude auth login`, yang membuka browser agar Anda bisa masuk. Di mesin jarak jauh, browser Anda mungkin berakhir di halaman `http://localhost:…/callback` yang tidak bisa dimuat; tempelkan alamat itu ke terminal dan ccam akan menyelesaikan proses masuk. Setelah proses masuk selesai, Anda akan melihat:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Bawa pengaturan yang sudah ada

Jika Anda sudah pernah memakai Claude Code, pengaturan, `CLAUDE.md`, dan ekstensi Anda ada di `~/.claude`, dan server MCP Anda ada di `~/.claude.json`. Salin semuanya ke grup bersama `default` agar setiap akun mendapatkannya:

```sh
ccam migrate
```

ccam menampilkan apa yang akan disalin dan meminta persetujuan terlebih dahulu. ccam hanya membaca `~/.claude`, tidak pernah menyalin login, dan tidak menyertakan riwayat percakapan kecuali Anda menambahkan `--history`. Lihat [menyalin penyiapan yang sudah ada](/docs/sharing#copy-an-existing-setup) untuk detailnya.

## 3. Tambahkan akun kedua

```sh
ccam add work
```

Wizard melewati pertanyaan nama, karena Anda sudah memberikannya di baris perintah, dan menampilkan grup yang bisa diikuti dengan `default` sudah dicentang. Kali ini, beri akun itu perintah pintasannya sendiri, seperti `ccwork`, agar Anda bisa menjalankannya dengan satu kata. Setelah grup, wizard akan bertanya:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Tekan `y`, lalu ketik nama, atau tekan Enter untuk memakai nama yang disarankan:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Lalu masuk dengan akun yang lain di browser. ccam mengonfirmasi akun itu beserta pintasannya:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Berpindah akun

`ccam work` menjalankan Claude Code sebagai akun kerja. Jika Anda melewati pertanyaan pintasan, tambahkan pintasan kapan saja:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Sekarang setiap perintah di bawah ini menjalankan Claude Code sebagai akun kerja (`2` adalah nomornya di `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Menjalankan `ccam` saja akan membuka [dasbor](/features/dashboard): pindah ke sebuah akun dengan tombol panah atau nomornya, lalu tekan Enter.

Semua yang ada setelah pintasan atau nama akun diteruskan ke Claude Code tanpa diubah, jadi Anda tetap bisa memakai opsi yang sudah Anda kenal:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Pilih apa yang dibagikan

Secara default, server MCP, `settings.json`, `CLAUDE.md` global, dan ekstensi (agen, perintah, skill, dan plugin) dibagikan antarakun melalui grup `default`. Riwayat percakapan tidak dibagikan. Tambahkan riwayat ke grup itu jika Anda ingin melanjutkan percakapan antarakun:

```sh
ccam group set default --add history
```

ccam menampilkan apa yang akan dipindahkannya untuk setiap akun sebelum mengubah apa pun. Lihat [berbagi dengan grup](/features/share) untuk detailnya.

## 6. Saat sebuah akun mencapai batasnya

```sh
ccam usage
ccam next --continue
```

`ccam usage` menampilkan penggunaan 5 jam dan mingguan semua akun. `ccam next --continue` membuka percakapan yang sama di akun berikutnya. Di skrip, `ccam --fallback work -p "…"` berganti akun dengan sendirinya. Lihat [batas penggunaan](/docs/usage-limits).

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.

## Sudah menggunakan ~/.claude?

Daftarkan folder itu alih-alih masuk lagi:

```sh
ccam import
```

Folder tersebut tetap berada di tempatnya, dan `claude` biasa tetap berfungsi seperti sebelumnya.
