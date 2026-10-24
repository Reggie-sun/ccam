---
title: Penjelasan CLAUDE_CONFIG_DIR
description: Fungsi CLAUDE_CONFIG_DIR di Claude Code, isi foldernya, tempat macOS menyimpan login, dan file mana yang bisa dibagikan antarakun.
updated: 2026-10-23
order: 2
translation: machine
sourceHash: 2cf08478085f
---

`CLAUDE_CONFIG_DIR` memberi tahu Claude Code folder mana yang dipakai sebagai pengganti `~/.claude`. Semua yang khusus untuk akun berada di sana: login, pengaturan, server MCP, plugin, dan riwayat percakapan. Mengarahkannya ke folder yang berbeda memberi Anda akun yang terpisah.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Apakah CLAUDE_CONFIG_DIR didukung secara resmi?

Ya. Claude Code mencantumkannya di antara [variabel lingkungannya](https://code.claude.com/docs/en/env-vars) sebagai cara mengganti folder konfigurasi, "berguna untuk menjalankan beberapa akun secara berdampingan", dan [panduan autentikasinya](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) memakainya untuk login dengan beberapa akun. Anda bisa mengaturnya di shell atau di pengaturan pengguna, tetapi tidak di pengaturan proyek.

## Apa isi folder tersebut?

Diperiksa dengan Claude Code 2.1.309:

| Path | Isi | Bisa dibagikan |
| --- | --- | --- |
| `.credentials.json` | Token login (Linux, Windows), token OAuth MCP | Tidak |
| `.claude.json` | Detail akun, server MCP pengguna, status per proyek | Tidak (server MCP secara terpisah) |
| `settings.json` | Izin, hook, model, status line, plugin | Ya |
| `CLAUDE.md` | Instruksi global | Ya |
| `agents/`, `commands/`, `skills/`, `plugins/` | Ekstensi | Ya |
| `projects/` | Percakapan sebagai file `.jsonl`, memori per proyek | Ya |
| `history.jsonl` | Riwayat prompt | Ya |
| `sessions/`, `shell-snapshots/`, `cache/` dan lainnya | Status runtime | Tidak |

## Apa yang sering menjebak?

**Tanpa `CLAUDE_CONFIG_DIR`, `.claude.json` tidak berada di dalam `~/.claude`.** Dalam hal itu, Claude Code menyimpannya di `~/.claude.json`, di sebelah folder tersebut. Jika Anda mulai menggunakan `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code akan mencari `~/.claude/.claude.json` dan tidak menemukan akun Anda. `ccam import` menangani hal ini dengan menautkan keduanya. Jika `~/.claude` sudah memiliki `.claude.json` sendiri, dari proses sebelumnya dengan `CLAUDE_CONFIG_DIR=~/.claude`, ccam membiarkan kedua file itu dan memberitahukannya, karena hanya Anda yang tahu mana yang benar.

**Di macOS, login berada di Rantai Kunci (Keychain), bukan di folder.** Claude Code menyimpannya dengan nama `Claude Code-credentials`, diikuti `-` dan 8 digit heksadesimal pertama dari SHA-256 path folder saat `CLAUDE_CONFIG_DIR` ditetapkan, untuk akun yang disebut oleh `$USER`. Karena itu, folder yang berbeda mendapatkan entri Rantai Kunci yang berbeda, dan folder yang dipindahkan atau diganti namanya kehilangan login-nya; itulah sebabnya `ccam rename` mempertahankan nama folder di macOS. Jika Rantai Kunci menolak penulisan, Claude Code beralih ke `.credentials.json` di dalam folder, seperti yang dijelaskan [dokumentasinya](https://code.claude.com/docs/en/authentication#credential-management).

**Tidak semua login disimpan di folder.** Login Claude Console tanpa kunci API disimpan di luar folder, sehingga folder terpisah tidak bisa memisahkan dua login semacam itu. Login claude.ai dan kunci API disimpan per folder.

## Di mana folder itu di Windows?

Secara default di `%USERPROFILE%\.claude`, dengan login di `.credentials.json` di dalamnya, atau di Windows Credential Manager jika Claude Code diatur untuk memakainya, dengan nama yang dibuat dari path folder seperti di macOS. Atur `CLAUDE_CONFIG_DIR` ke folder lain untuk setiap akun tambahan, sama seperti di Linux dan macOS.

## Bisakah folder berbagi file?

Tautan simbolik berfungsi: Claude Code membaca dan menulis melalui `settings.json` atau folder `projects/` yang ditautkan, dan mempertahankan tautannya. Server MCP adalah pengecualian, karena berbagi `.claude.json` dengan login; teruskan server MCP dengan `--mcp-config` sebagai gantinya.

ccam dibangun tepat di atas aturan-aturan ini. Lihat [pembagian](/docs/sharing) untuk cara ccam menautkan folder, dan [panduan beberapa akun](/guides/multiple-claude-code-accounts) untuk gambaran yang lebih luas.
