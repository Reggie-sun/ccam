---
title: Dasbor
description: ccam di terminal membuka dasbor layar penuh berisi akun, grup, aturan folder, ekstensi, dan pemeriksaan; setiap perubahan menampilkan rencananya dulu.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` tanpa argumen membuka dasbor layar penuh jika input dan output standar adalah terminal; jika bukan, ccam mencetak `ccam list`. Dasbor menampilkan setiap akun beserta penggunaan dan grupnya, dan melakukan apa yang dilakukan perintah-perintahnya: menjalankan, menambahkan, login, mengubah grup, menetapkan aturan folder, memeriksa, dan memperbaiki. Menjalankan akun akan menutupnya; menambahkan akun atau login meninggalkannya sebentar untuk terminal lalu kembali lagi.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Bingkai pertama digambar hanya dari `config.toml`, jauh di bawah 100 ms; penggunaan terisi begitu setiap akun menjawab. Persentase berubah kuning mulai 80% dan merah pada 100%, dan nilai dari cache setelah pencarian yang gagal digambar redup.

## Tab

| Tab | Yang ditampilkan | Tombol |
| --- | --- | --- |
| Accounts | nomor, nama, email, bilah 5 jam, penggunaan 7 hari, pemakaian terakhir, grup; `✔` akun default, `▸` akun folder ini, `•` sedang berjalan | lihat di bawah |
| Groups | setiap grup, apa yang dibagikannya dan akun-akunnya, serta akun yang tidak ada di grup mana pun | `enter` sunting, `n` baru, `e` ganti nama, `x` hapus |
| Folders | aturan [`ccam folder`](/features/folder) dan aturan yang berlaku di sini | `n` pakai sebuah akun di sini, `enter` ubah, `x` hapus |
| Extensions | apa yang ditampilkan [`ccam ext`](/features/ext), setiap ekstensi satu kali beserta grup asalnya | `r` muat ulang |
| Health | apa yang dilaporkan [`ccam doctor`](/features/doctor) | `f` perbaiki (`doctor --fix`) |

`←`/`→` atau `Tab` berpindah tab, `?` menampilkan tombol-tombol tab yang sedang dibuka, `r` mencari penggunaan lagi tanpa cache singkat, dan `q` atau `Esc` keluar.

## Tombol di tab Accounts

| Tombol | Fungsinya |
| --- | --- |
| `Enter` | Menjalankan Claude Code sebagai akun yang dipilih |
| `1`–`9` | Pindah ke akun dengan nomor itu; ketik `1` `2` dengan cepat untuk 12. Hanya `Enter` yang menjalankan |
| `↑`/`↓`, `k`/`j` | Berpindah |
| `n` | Menjalankan akun yang akan dipilih [`ccam next --auto`](/features/next) |
| `o` | Menjalankan dengan opsi: `--continue`, `--resume`, argumen lain, atau uji coba (dry run) |
| `u` | Rincian penggunaan: jendela 5 jam, 7 hari, dan per model beserta waktu reset-nya dan asal angka-angkanya |
| `a` | Menambahkan akun (wizard di bawah) |
| `l` | Login ulang |
| `s` | Memilih grup akun tersebut |
| `d` | Menjadikannya akun default |
| `f` | Memakainya di folder saat ini |
| `e` | Mengganti namanya |
| `c` | Menetapkan atau menghapus perintah pintasannya |
| `x` | Menghapusnya (bertanya dulu, dimulai dari No) |

## Mengubah grup

Di tab Groups, `Enter` membuka penyunting: item di kiri, akun di kanan, `Tab` untuk berpindah sisi, `Space` untuk mencentang atau menghapus centang. Item yang sudah didapat salah satu akun dari grup lain diredupkan beserta alasannya, begitu pula akun yang akan mendapat sebuah item dua kali; satu item hanya boleh berasal dari satu grup. `Enter` menampilkan rencananya:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

File yang berbeda dari salinan milik grup ditanyakan lebih dulu (pertahankan milik grup, pakai milik saya, atau gabungkan). Lalu perubahan berjalan di latar belakang dengan kode yang sama seperti [`ccam group`](/features/share): konfigurasi disimpan lebih dulu, data setiap akun dipindahkan secara berurutan, dan file aslinya masuk ke `.ccam-backup/`. Perubahan yang akan memindahkan data akun yang sedang berjalan ditolak dengan menyebut nama akunnya.

## Menambahkan akun

`a` membuka wizard yang juga dipakai `ccam add` di terminal:

1. Nama akun.
2. Cara login: langganan Claude sekarang (di browser), [kunci API](/features/api), atau langganan Claude nanti.
3. Untuk kunci API: dari mana ccam membaca kunci itu setiap kali dijalankan (`env:NAME` atau `command:COMMAND`), endpoint, cara kunci dikirim, dan modelnya.
4. Grup yang akan diikuti, dengan `default` sudah dicentang; jika belum ada grup, apa yang dibagikan grup `default`.
5. Perintah pintasan, jika Anda menginginkannya.
6. Ringkasan untuk dikonfirmasi.

`Esc` kembali satu langkah, dan pada langkah pertama menutup wizard.

## Warna

Warna mengikuti peran warna di UI terminal Claude Code, sehingga ccam terlihat serasi di sampingnya; nama grup mendapat warna berdasarkan namanya. ccam menanyakan warna latar belakang ke terminal lalu memilih tema terang atau gelap. Tetapkan sendiri di `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Tema `-daltonized` ditujukan untuk buta warna merah-hijau, dan tema `-ansi` memakai 16 warna milik terminal itu sendiri. `NO_COLOR` mematikan warna.

## Keamanan dan kasus khusus

- Dasbor memanggil fungsi yang sama dengan perintah-perintahnya, jadi setiap aturan tetap berlaku: rencana sebelum data dipindahkan, penolakan untuk akun yang sedang berjalan, cadangan, dan token login yang tidak pernah meninggalkan akunnya.
- Dasbor tidak punya uji coba (dry run): `ccam --dry-run` tanpa nama mencetak bagaimana akun folder ini akan dijalankan, dan `o` pada sebuah akun menawarkan uji coba untuk satu kali menjalankan akun itu saja.
- Terminal dipulihkan sebelum Claude Code dimulai atau login dibuka.
- Layar yang terlalu kecil untuk sebuah dialog lebih dulu membuang baris-baris tab; baris dipotong dengan `…` alih-alih dibungkus ke baris berikutnya.

## Terkait

- [Berbagi dengan grup](/features/share)
- [Menjalankan Claude Code](/features/launch)
- [Mengelola akun](/features/accounts)
- [Batas penggunaan](/features/usage)
