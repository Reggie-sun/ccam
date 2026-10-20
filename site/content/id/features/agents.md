---
title: Sesi latar belakang semua akun
description: ccam agents menampilkan agen latar belakang dan sesi terbuka Claude Code dari semua akun dalam satu tabel, dan ccam attach membuka salah satunya lewat id-nya.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code menyimpan sesi latar belakang (`claude --bg`, tampilan agen) per akun: `claude agents` hanya menampilkan sesi milik akun yang menjalankannya. `ccam agents` bertanya ke semua akun sekaligus dan menampilkannya bersama-sama, dan `ccam attach <id>` membuka sesi latar belakang sebagai akun mana pun yang memulainya.

## Penggunaan

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` adalah status sesi latar belakang (`working`, `blocked`, `done`, `failed`, `stopped`) dan apa yang sedang dikerjakan prosesnya (`busy`, `waiting`, `idle`); sesi interaktif yang terbuka menampilkan `open`.

## Opsi

| Opsi | Fungsinya | Bawaan |
| --- | --- | --- |
| `--all` | Juga menampilkan sesi latar belakang yang sudah selesai, gagal, atau dihentikan | nonaktif |
| `--json` | Mencetak sesi sebagai JSON, masing-masing dengan `account` pemiliknya | nonaktif |

`ccam attach <id> [claude args]` meneruskan argumen setelah id ke `claude attach`. Dengan `--dry-run`, perintah itu mencetak perintahnya.

## Cara kerjanya

1. Untuk setiap akun, ccam menjalankan `claude agents --json` (dengan `--all` jika diberikan) dengan `CLAUDE_CONFIG_DIR` milik akun tersebut, semuanya sekaligus, masing-masing dalam batas 15 detik. Ini membaca daftar sesi milik akun itu sendiri dan tidak memulai layanan latar belakang Claude Code.
2. Setiap entri diberi nama akunnya, dan daftar diurutkan menurut waktu mulai, yang terbaru lebih dulu. Akun yang tidak menjawab disebutkan di error standar dan dilewati, lalu ccam keluar dengan kode 1, sehingga skrip tidak mengira kegagalan sebagai daftar kosong.
3. `ccam attach` mencari `jobs/<id>` di folder setiap akun, tempat Claude Code menyimpan sesi latar belakang, lalu menjalankan `claude attach <id>` sebagai akun tersebut, dengan cara yang sama seperti `ccam <name>` memulai Claude Code.

## Keamanan dan kasus khusus

- Id sesi latar belakang milik satu akun: `claude attach` yang dijalankan sebagai akun lain tidak dapat menemukannya, dan itulah yang diselesaikan oleh `ccam attach`.
- Dua akun bisa memiliki sesi yang id-nya berawalan sama; `ccam attach` lalu menyebutkan keduanya dan meminta bagian id yang lebih panjang.
- ccam meneruskan server MCP bersamanya ke Claude Code dengan `--mcp-config=<file>` dalam satu argumen, yang diterima oleh `claude agents` dan `claude attach`; sesi latar belakang tetap memakai file itu saat dimulai ulang.
- Layanan latar belakang akun kunci API mengambil kuncinya dari shell yang memulainya. Mulailah sesi latar belakang pertama akun seperti itu melalui ccam (`ccam gw --bg …`), agar layanan tersebut memiliki kuncinya.

## Terkait

- [Menjalankan Claude Code](/features/launch)
- [Memilih akun untuk skrip](/features/pick)
- [Dasbor](/features/dashboard)
