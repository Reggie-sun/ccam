---
title: Berbagi server MCP, pengaturan, dan riwayat
description: Cara ccam membagikan server MCP, settings.json, CLAUDE.md, agen, skill, dan riwayat percakapan antarakun Claude Code dengan login tetap terpisah.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

Setiap akun menyimpan login-nya sendiri. Semua hal lain dapat dibagikan per item: ccam menautkan file atau folder milik akun ke folder grup bersama, sehingga perubahan yang dibuat di satu akun muncul di semua akun dalam grup tersebut.

## Apa saja yang bisa dibagikan

| Item | Opsi | Yang dibagikan |
| --- | --- | --- |
| Server MCP | `--mcp` | `mcp.json` milik grup, diteruskan ke Claude Code saat dimulai |
| Pengaturan | `--settings` | `settings.json`: izin, hook, model, status line |
| Memori global | `--memory` | `CLAUDE.md` |
| Ekstensi | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Riwayat percakapan | `--history` | `projects/` dan `history.jsonl` |

Akun baru membagikan semuanya kecuali riwayat percakapan.

File login `.credentials.json` dan `.claude.json` tidak pernah dibagikan. ccam hanya menautkan path dari daftar tetap, dan kedua file ini tidak ada di dalamnya.

## Grup

`on` berarti grup `default`. Anda juga bisa memberi nama, sehingga hanya sebagian akun yang berbagi suatu item:

```sh
ccam share work --history team
ccam share side --history team
```

Sekarang `work` dan `side` berbagi percakapan satu sama lain, tetapi tidak dengan akun Anda yang lain.

## Mengaktifkan pembagian

Saat sebuah akun mulai membagikan suatu item, ccam:

1. menampilkan rencana untuk setiap file yang akan dipindahkan, disalin, atau ditautkan, dan meminta persetujuan sebelum memulai;
2. memindahkan file akun ke folder grup jika grup masih kosong, atau menambahkan file yang belum dimiliki grup;
3. jika file yang sama berbeda isinya, menanyakan apakah akan mempertahankan versi grup, memakai versi Anda, atau menggabungkan JSON-nya;
4. menyimpan file asli Anda di `.ccam-backup/<time>/` di dalam folder akun;
5. mengganti file asli dengan tautan ke grup.

File percakapan diberi nama sesuai ID sesinya, sehingga riwayat dari beberapa akun dapat digabungkan tanpa konflik.

## Menonaktifkan pembagian

```sh
ccam share work --history off
```

Tautan dihapus dan akun mendapatkan salinannya sendiri dari isi grup saat ini. Tambahkan `--fresh` untuk memulai dari kosong. Grup tetap menyimpan isinya untuk akun-akun lain.

## Server MCP

Claude Code menyimpan server MCP di `.claude.json`, bersama dengan login, sehingga file tersebut tidak dapat dibagikan. ccam menyimpan server bersama di `mcp.json` milik grup dan menjalankan Claude Code dengan `--mcp-config=<group>/mcp.json`. Opsi `--mcp-config` Anda sendiri tetap berlaku di atasnya.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Server yang Anda tambahkan di dalam sesi dengan `claude mcp add` hanya disimpan untuk akun tersebut. Saat berikutnya Anda menjalankan akun itu, ccam menyebutkan nama server tersebut beserta perintah untuk memindahkannya ke grup:

```sh
ccam mcp import work gh
```

Server MCP yang masuk dengan OAuth menyimpan tokennya bersama login akun, jadi otorisasi server tersebut sekali di setiap akun.

Jika Anda lebih suka server muncul di `claude mcp list` dan `/mcp` sebagai server pengguna biasa, tetapkan `mcp_strategy = "sync"` di `~/.config/ccam/config.toml`. ccam kemudian menuliskannya ke `.claude.json` setiap akun sebelum memulai, di bawah file lock milik Claude Code sendiri dan dengan cadangan.

## Dua akun di folder yang sama

Akun yang berbagi riwayat percakapan menulis ke file yang sama. Jika Anda menjalankan Claude Code kedua di folder tempat akun lain dari grup yang sama sedang berjalan, ccam akan memperingatkan Anda dan menanyakan apakah akan melanjutkan.

## Jika ada yang rusak

```sh
ccam doctor
ccam doctor --fix
```

Jika path bersama berubah menjadi file biasa, misalnya setelah Anda memulihkan cadangan secara manual, ccam berhenti sebelum menjalankan Claude Code dan mengarahkan Anda ke `ccam doctor --fix`, yang menggabungkan file tersebut ke grup dan menyimpan cadangan.
