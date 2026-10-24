---
title: Memulai cepat
metaTitle: "Memulai cepat: tambah dan ganti akun Claude Code"
description: Tambahkan dua akun Claude Code dengan ccam add, beri pintasan seperti ccwork, berpindah di antaranya, dan pilih apa yang dibagikan.
updated: 2026-10-24
order: 2
translation: machine
sourceHash: e98d0637356a
---

Tambahkan setiap akun dengan `ccam add`, lalu jalankan Claude Code dengan `ccam <name>` atau dengan pintasan yang namanya Anda pilih sendiri, seperti `ccwork`. Bagi kebanyakan orang, itu sudah cukup.

## 1. Tambahkan akun pertama Anda

```sh
ccam add
```

Wizard akan menanyakan nama, seperti `personal` atau `work`, item apa saja yang perlu dibagikan akun ini dengan akun Anda yang lain, dan apakah perlu membuat perintah pintasan (lihat langkah 4). Lalu wizard menjalankan `claude auth login`, yang membuka browser agar Anda bisa masuk. Di mesin jarak jauh, browser mungkin berakhir di halaman `http://localhost:…/callback` yang tidak bisa dimuat; tempelkan alamat itu ke terminal dan ccam akan menyelesaikan proses masuk. Setelah selesai, Anda akan melihat:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Bawa pengaturan yang sudah ada

Jika Anda sudah pernah memakai Claude Code, pengaturan, `CLAUDE.md`, ekstensi, dan server MCP Anda ada di `~/.claude`. Salin semuanya ke pengaturan bersama agar setiap akun mendapatkannya:

```sh
ccam migrate
```

ccam menampilkan apa yang akan disalin dan meminta persetujuan terlebih dahulu. ccam hanya membaca `~/.claude`, tidak pernah menyalin login, dan tidak menyertakan riwayat percakapan kecuali Anda menambahkan `--history`. Lihat [menyalin penyiapan yang sudah ada](/docs/sharing#copy-an-existing-setup) untuk detailnya.

## 3. Tambahkan akun kedua

```sh
ccam add work
```

Masuk dengan akun yang lain di browser.

## 4. Berpindah akun

`ccam work` menjalankan Claude Code sebagai akun kerja. Agar lebih singkat, beri akun itu perintah pintasannya sendiri, seperti `ccwork`. Setelah item yang dibagikan, wizard akan bertanya:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Tekan `y`, lalu ketik nama, atau tekan Enter untuk memakai nama yang disarankan:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Setelah login, ccam mengonfirmasinya:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Jika Anda melewati pertanyaan itu, tambahkan pintasan kapan saja:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Sekarang perintah mana pun di bawah ini menjalankan Claude Code sebagai akun kerja:

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

Secara default, server MCP, `settings.json`, `CLAUDE.md` global, serta agen, perintah, skill, dan plugin dibagikan antarakun melalui grup `default`. Riwayat percakapan tidak dibagikan. Tambahkan riwayat ke grup itu jika Anda ingin melanjutkan percakapan antarakun:

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

**Kursi Team dan Enterprise.** Pilihan otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama. Pusat bantuan Anthropic menyatakan "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (setiap kursi pada paket Team diberikan kepada satu orang tertentu, dan hanya orang itu yang boleh menggunakannya) dan "Users can't route their work through another user's account to use that user's remaining usage" (pengguna tidak boleh menyalurkan pekerjaannya melalui akun pengguna lain untuk memakai sisa kuota pengguna tersebut) ([paket Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); begitu pula untuk [paket Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Karena itu, memakai pilihan otomatis di antara kursi dapat melanggar ketentuan Anthropic. Tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya. Paket pribadi (Pro, Max) dibahas dalam [catatan tentang ketentuan di FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Sudah menggunakan ~/.claude?

Daftarkan folder itu alih-alih masuk lagi:

```sh
ccam import
```

Folder tersebut tetap berada di tempatnya, dan `claude` biasa tetap berfungsi seperti sebelumnya.
