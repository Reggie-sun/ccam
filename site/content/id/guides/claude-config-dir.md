---
title: Penjelasan CLAUDE_CONFIG_DIR
description: Fungsi CLAUDE_CONFIG_DIR di Claude Code, isi foldernya, tempat macOS menyimpan login, dan file mana yang bisa dibagikan antarakun.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` memberi tahu Claude Code folder mana yang dipakai sebagai pengganti `~/.claude`. Semua yang khusus untuk akun berada di sana: login, pengaturan, server MCP, plugin, dan riwayat percakapan. Mengarahkannya ke folder yang berbeda memberi Anda akun yang terpisah.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Apa isi folder tersebut

Diperiksa dengan Claude Code 2.1.302:

| Path | Isi | Bisa dibagikan |
| --- | --- | --- |
| `.credentials.json` | Token login (Linux), token OAuth MCP | Tidak |
| `.claude.json` | Detail akun, server MCP pengguna, status per proyek | Tidak (server MCP secara terpisah) |
| `settings.json` | Izin, hook, model, status line, plugin | Ya |
| `CLAUDE.md` | Instruksi global | Ya |
| `agents/`, `commands/`, `skills/`, `plugins/` | Ekstensi | Ya |
| `projects/` | Percakapan sebagai file `.jsonl`, memori per proyek | Ya |
| `history.jsonl` | Riwayat prompt | Ya |
| `sessions/`, `shell-snapshots/`, `cache/` dan lainnya | Status runtime | Tidak |

## Dua detail yang sering menjebak

**Tanpa `CLAUDE_CONFIG_DIR`, `.claude.json` tidak berada di dalam `~/.claude`.** Dalam hal itu, Claude Code menyimpannya di `~/.claude.json`, di sebelah folder tersebut. Jika Anda mulai menggunakan `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code akan mencari `~/.claude/.claude.json` dan tidak menemukan akun Anda. `ccam import` menangani hal ini dengan menautkan keduanya.

**Di macOS, login berada di Rantai Kunci (Keychain), bukan di folder.** Claude Code menyimpannya dengan nama `Claude Code-credentials`, diikuti hash pendek dari path folder saat `CLAUDE_CONFIG_DIR` ditetapkan. Karena itu, folder yang berbeda mendapatkan entri Rantai Kunci yang berbeda. Jika Rantai Kunci menolak penulisan, Claude Code beralih ke `.credentials.json` di dalam folder.

## Berbagi antarfolder

Tautan simbolik berfungsi: Claude Code membaca dan menulis melalui `settings.json` atau folder `projects/` yang ditautkan, dan mempertahankan tautannya. Server MCP adalah pengecualian, karena berbagi `.claude.json` dengan login; teruskan server MCP dengan `--mcp-config` sebagai gantinya.

ccam dibangun tepat di atas aturan-aturan ini. Lihat [pembagian](/docs/sharing) untuk cara ccam menautkan folder, dan [panduan beberapa akun](/guides/multiple-claude-code-accounts) untuk gambaran yang lebih luas.
