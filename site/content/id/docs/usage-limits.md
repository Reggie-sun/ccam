---
title: Batas penggunaan dan perpindahan akun
description: Lihat penggunaan 5 jam dan mingguan setiap akun Claude Code dengan ccam usage, lalu buka akun dengan sisa kuota terbanyak lewat ccam next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a64c2c341db3
---

`ccam usage` menampilkan seberapa banyak batas 5 jam dan mingguan yang sudah digunakan setiap akun. `ccam next --auto` membuka akun dengan sisa kuota terbanyak, dan `--continue` membawa percakapan saat ini ke akun tersebut.

```console
$ ccam usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccam next --auto 로 바로 열기)
```

Kolomnya adalah: akun, penggunaan 5 jam, waktu hingga direset, penggunaan mingguan, waktu hingga direset, dan sudah berapa lama sejak angka-angka itu diambil.

## Dari mana angkanya berasal

ccam menanyakan endpoint penggunaan Anthropic untuk semua akun secara bersamaan, menggunakan login masing-masing akun. ccam menunggu paling lama 3 detik per akun dan menampilkan setiap akun begitu akun itu menjawab.

Endpoint ini tidak didokumentasikan secara publik dan dapat berubah. Jika permintaan gagal, ccam menampilkan angka terakhir yang diketahui dan menandainya sebagai cache. `ccam doctor` memberi tahu Anda jika format responsnya berubah.

Angka yang diambil kurang dari 60 detik yang lalu akan dipakai ulang. Opsi:

| Opsi | Efek |
| --- | --- |
| `--watch [seconds]` | Memperbarui terus-menerus, setiap 60 detik secara default |
| `--json` | Keluaran yang dapat dibaca mesin |
| `-v` | Juga menampilkan batas mingguan per model dan asal setiap angka |
| `--no-cache` | Selalu bertanya ulang |
| `--source cache` | Tidak pernah bertanya; hanya menampilkan angka yang tersimpan |

## Login yang kedaluwarsa

Saat token akses sebuah akun sudah kedaluwarsa, ccam tidak memperbaruinya, karena memperbaruinya saat sesi Claude Code sedang memakai login yang sama dapat membuat sesi itu keluar. Jalankan akun tersebut sekali dengan `ccam <name>` dan Claude Code akan memperbaruinya.

`ccam usage --refresh` memperbarui token yang kedaluwarsa sendiri, tetapi hanya untuk akun yang tidak sedang berjalan. Perintah ini tidak pernah memperbarui login yang tersimpan di Rantai Kunci (Keychain) macOS atau di `~/.claude` yang diimpor.

## Penggunaan di status line

```sh
ccam statusline install
```

Claude Code kemudian menampilkan nama akun dan penggunaannya di status line, misalnya `[work] · Opus · 5h 24% · 7d 41%`. Jika Anda sudah memiliki perintah status line, ccam mempertahankannya dan menaruh nama akun di depannya. Status line juga menyimpan angkanya, yang dipakai `ccam usage` saat endpoint tidak menjawab.

## Berpindah akun

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` melewati akun yang belum login, dan `--auto` melewati akun yang sudah mencapai 100%. Untuk melanjutkan percakapan, kedua akun harus berbagi riwayat percakapan; lihat [pembagian](/docs/sharing).

[Panduan alur kerja batas penggunaan](/guides/usage-limit-workflow) menjelaskan contoh lengkapnya.
