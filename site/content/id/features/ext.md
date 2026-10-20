---
title: Memeriksa ekstensi
description: ccam ext menampilkan skill, perintah, plugin, server MCP, dan hook yang dimuat setiap akun, asal masing-masing, dan apa yang rusak.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

Setiap sel menyebut asal item itu di akun itu: grup berbagi (`default`), `account` jika item itu ada di folder akun itu sendiri, atau `-` jika akun tidak memilikinya. `!` menandai item yang bermasalah, dijelaskan di bawah tabel.

## Yang dibacanya

| Jenis | Lokasi | Asal yang ditampilkan |
| --- | --- | --- |
| Skill | setiap folder di `skills/` yang berisi `SKILL.md` | grup yang ditautkan `skills/`, atau `account` |
| Perintah, agen | file `.md` di bawah `commands/` dan `agents/`, subfolder sebagai `folder/name` | grup, atau `account` |
| Plugin | `plugins/installed_plugins.json` | grup, atau `account` |
| Server MCP | `mcp.json` milik grup (jika MCP dibagikan) dan `.claude.json` milik akun | grup, atau `account` |
| Hook | `hooks` di `settings.json`, dihitung per event | grup, atau `account` |
| Proyek ini | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json`, dan `.mcp.json` di folder saat ini | ditampilkan sekali, untuk setiap akun |

## Masalah yang dilaporkannya

- **Tautan rusak**: tautan di `skills/`, `commands/`, atau `agents/` yang targetnya sudah tidak ada, atau folder bersama (`skills`, `commands`, `agents`, `plugins`) yang menaut ke tempat yang hilang (perbaiki dengan `ccam doctor --fix`).
- **Folder skill tanpa `SKILL.md`**: Claude Code melewatinya.
- **Server MCP yang didefinisikan berbeda di dua tempat**: `mcp.json` milik grup dan `.claude.json` milik akun, atau salah satunya dan `.mcp.json` milik proyek. Pengaturan yang sama di kedua tempat (yang ditulis oleh `mcp_strategy = "sync"`) bukan masalah. Pindahkan server milik akun sendiri ke grup dengan `ccam mcp import <account> <server>`.
- **`installed_plugins.json` yang bukan JSON yang valid.**

## Keamanan dan kasus khusus

- `ccam ext` tidak mengubah apa pun dan tidak membuat permintaan jaringan. Perintah ini membaca nama server dan jumlah hook, tidak pernah nilai `env` atau header dalam definisi MCP.
- Perintah ini menampilkan apa yang ada di disk. Plugin yang terpasang tetapi dimatikan di pengaturan, atau server yang diblokir oleh pengaturan terkelola, tetap muncul.
- Jalankan di folder proyek untuk menyertakan folder `.claude` dan `.mcp.json` proyek itu.

## Terkait

- [Berbagi antarakun](/features/share)
- [Server MCP bersama](/features/mcp)
- [Memeriksa dan memperbaiki](/features/doctor)
