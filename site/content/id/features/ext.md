---
title: Memeriksa ekstensi
description: ccam ext menampilkan skill, perintah, plugin, server MCP, dan hook yang dimuat setiap akun, asal masing-masing, dan apa yang rusak.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: 45740e4efbef
---

`ccam ext` menjawab "mengapa skill (atau server MCP) ini berfungsi di satu akun tetapi tidak di akun lain?". Untuk setiap akun, perintah ini mendaftar ekstensi yang dimuat Claude Code dan asal masing-masing: grup berbagi, folder akun itu sendiri, atau proyek tempat Anda berada. Di bawah tabel, ia mendaftar hal-hal yang membuat sebuah ekstensi tidak dimuat seperti yang Anda harapkan. Perintah ini hanya membaca.

## Penggunaan

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

Setiap sel menyebut asal item itu di akun itu: grup berbagi (`default`), `account` jika item itu ada di folder akun itu sendiri, `local` untuk server MCP yang dimiliki akun hanya untuk folder ini, atau `-` jika akun tidak memilikinya. Plugin yang terpasang tetapi tidak diaktifkan menampilkan `(off)` setelah asalnya. `!` menandai item yang bermasalah, dijelaskan di bawah tabel; masalah yang bukan milik akun mana pun, seperti server yang hanya dimiliki `claude` biasa, muncul paling akhir tanpa nama. `ccam ext --json` mencetak laporan yang sama (`accounts`, proyeknya, `problems`), dengan `remote` pada server MCP yang dijangkau Claude Code melalui jaringan; tanpa akun, perintah ini mencetak laporan dengan `accounts: []`.

## Yang dibacanya

| Jenis | Lokasi | Asal yang ditampilkan |
| --- | --- | --- |
| Skill | setiap folder di `skills/` yang berisi `SKILL.md` | grup yang ditautkan `skills/`, atau `account` |
| Perintah, agen | file `.md` di bawah `commands/` dan `agents/`, subfolder sebagai `folder/name` | grup, atau `account` |
| Plugin | `plugins/installed_plugins.json`, diaktifkan atau dimatikan oleh `enabledPlugins` di pengaturan akun dan proyek; plugin yang dipasang hanya untuk satu proyek hanya muncul di proyek itu | grup, atau `account`; `(off)` jika tidak diaktifkan |
| Server MCP | `mcp.json` milik grup (jika MCP dibagikan), `.claude.json` milik akun (server miliknya sendiri, dan yang dimilikinya hanya untuk folder ini), dan server dari plugin-pluginnya, sebagai `plugin:<plugin>:<server>` | grup, `account`, atau `local` |
| Hook | `hooks` di `settings.json`, dihitung per event | grup, atau `account` |
| Proyek ini | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json`, dan `.mcp.json` di folder saat ini | ditampilkan sekali; akun yang belum memercayai folder itu mendapat masalah untuk server `.mcp.json`-nya |

## Masalah yang dilaporkannya

- **Tautan rusak**: tautan di `skills/`, `commands/`, atau `agents/` yang targetnya sudah tidak ada, atau folder bersama (`skills`, `commands`, `agents`, `plugins`) yang menaut ke tempat yang hilang (perbaiki dengan `ccam doctor --fix`, yang membuat ulang folder grup yang hilang).
- **Folder skill tanpa `SKILL.md`**: Claude Code melewatinya.
- **Server MCP yang didefinisikan berbeda di dua tempat**: `mcp.json` milik grup dan `.claude.json` milik akun, atau salah satunya dan `.mcp.json` milik proyek. Pengaturan yang sama di kedua tempat (yang ditulis oleh `mcp_strategy = "sync"`) bukan masalah. Pindahkan server milik akun sendiri ke grup dengan `ccam mcp import <account> <server>`.
- **`installed_plugins.json` yang bukan JSON yang valid**, plugin yang tidak disebutkan oleh `enabledPlugins`, dan plugin yang foldernya hilang.
- **Server yang hanya dimuat satu akun**: server yang hanya ada di `.claude.json` akun ini, yang tidak dimuat akun lain dalam grupnya, dan server di cakupan lokal akun untuk folder ini (yang dilakukan `claude mcp add` kecuali diberi `-s user`). Masing-masing disertai perintah `ccam mcp import` untuk membagikannya. Server di `~/.claude.json` milik `claude` biasa yang tidak dimuat akun mana pun disertai perintah `ccam migrate --mcp`.
- **Server proyek yang menunggu persetujuan**: akun belum memercayai folder itu, sehingga Claude Code bertanya sekali, dan sampai saat itu server `.mcp.json` proyek menunggu. Masalah ini menyebutkan di mana menyetujuinya: `"enabledMcpjsonServers"` di `settings.json` milik akun-akun grup pengaturan, atau milik akun ini. Kepercayaan yang diberikan ke worktree utama sebuah repositori git berlaku juga di worktree lainnya.
- **Info masuk dan konektor**: server MCP yang memerlukan info masuk di akun ini sementara akun lain sudah masuk ke server itu (info masuk hanya berlaku untuk satu akun; masuk dengan `/mcp`), dan konektor claude.ai yang sudah dipakai akun lain tetapi belum oleh akun ini.

## Keamanan dan kasus khusus

- `ccam ext` tidak mengubah apa pun dan tidak membuat permintaan jaringan. Perintah ini membaca nama server dan jumlah hook, tidak pernah nilai `env` atau header dalam definisi MCP. Untuk pemeriksaan info masuk, perintah ini membaca nama info masuk MCP di login setiap akun (di macOS dari Rantai Kunci), tidak pernah tokennya, dan hanya untuk akun yang memuat server yang juga dimuat akun lain.
- Perintah ini menampilkan apa yang ada di disk. Server yang diblokir oleh pengaturan terkelola tetap muncul; plugin yang dimatikan di pengaturan menampilkan `(off)`.
- Jalankan di folder proyek untuk menyertakan folder `.claude` dan `.mcp.json` proyek itu.

## Terkait

- [Berbagi antarakun](/features/share)
- [Server MCP bersama](/features/mcp)
- [Memeriksa dan memperbaiki](/features/doctor)
