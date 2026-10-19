---
title: Referensi perintah
description: Semua perintah dan opsi ccam, dari ccam add dan ccam <name> hingga ccam usage, ccam next, dan ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

Letakkan `--dry-run` sebelum perintah apa pun untuk melihat apa yang akan dilakukannya tanpa benar-benar menjalankannya, misalnya `ccam --dry-run work`.

<!-- commands -->

## Menjalankan Claude Code

`ccam <name|number> [arguments]` menjalankan Claude Code sebagai akun tersebut. ccam tidak membaca argumen setelah nama akun; ccam meneruskannya ke `claude` dengan urutan yang sama, dan input standar, kode keluar, serta sinyal berperilaku persis seperti pada `claude` itu sendiri. `--` setelah nama akun bersifat opsional.

`ccam` dengan argumen tetapi tanpa nama akun, seperti `ccam -p "summarise"`, menggunakan akun default. Tetapkan akun default dengan `ccam use <name>`.

Pintasan `claude1`, `claude2`, dan seterusnya adalah tautan ke ccam di `~/.local/bin`. `ccam shims install --names` juga membuat `claude-<name>` untuk setiap akun.

## Menambah dan menghapus akun

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` tanpa folder akan mencari `~/.claude` dan `~/.claude-*`, lalu menanyakan masing-masing. Folder yang diimpor didaftarkan di tempatnya; ccam tidak pernah memindahkan atau menyalin login.

## Pembagian

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Setiap item menerima `on` (grup `default`), `off`, atau nama grup. Lihat [pembagian](/docs/sharing).

## Server MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Penggunaan dan perpindahan akun

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Pemeliharaan

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## File

| Path | Isi |
| --- | --- |
| `~/.config/ccam/config.toml` | Akun dan apa yang dibagikannya |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` milik setiap akun |
| `~/.local/share/ccam/groups/<group>/` | Pengaturan, server MCP, ekstensi, dan riwayat bersama |
| `~/.local/state/ccam/` | Cache penggunaan dan waktu terakhir digunakan |
| `~/.local/bin/claude1`, `claude2`… | Tautan ke ccam |

Variabel `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, dan `XDG_STATE_HOME` memindahkan folder-folder ini.
