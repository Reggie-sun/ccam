---
title: Referensi perintah
metaTitle: Referensi perintah ccam
description: Semua perintah dan opsi ccam, dari ccam add dan ccam <name> hingga ccam usage, ccam next, dan ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 65fb2b8f6617
---

Letakkan `--dry-run` sebelum perintah apa pun untuk melihat apa yang akan dilakukannya tanpa benar-benar menjalankannya, misalnya `ccam --dry-run work`.

<!-- commands -->

## Menjalankan Claude Code

`ccam <name|number> [arguments]` menjalankan Claude Code sebagai akun tersebut. ccam tidak membaca argumen setelah nama akun; ccam meneruskannya ke `claude` dengan urutan yang sama, dan input standar, kode keluar, serta sinyal berperilaku persis seperti pada `claude` itu sendiri. `--` setelah nama akun bersifat opsional.

`ccam` dengan argumen tetapi tanpa nama akun, seperti `ccam -p "summarise"`, menggunakan akun default. Tetapkan akun default dengan `ccam use <name>`.

Pintasan seperti `ccwork` menjalankan akunnya dengan cara yang sama: `ccwork --resume` melakukan hal yang sama dengan `ccam work --resume`. Lihat [pintasan](/docs/commands#shortcuts).

## Menambah dan menghapus akun

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` tanpa folder akan mencari `~/.claude` dan `~/.claude-*`, lalu menanyakan masing-masing. Folder yang diimpor didaftarkan di tempatnya; ccam tidak pernah memindahkan atau menyalin login.

## Pintasan

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Pintasan adalah perintah yang namanya Anda pilih sendiri dan menjalankan Claude Code sebagai satu akun. Seperti `ccam work`, pintasan meneruskan semua argumen ke Claude Code tanpa diubah, sehingga `ccwork --resume` dan `ccwork -p "hi"` bekerja seperti yang Anda harapkan. Pintasan adalah tautan ke ccam di `~/.local/bin`.

Pintasan bersifat opsional. Wizard `ccam add` menanyakan apakah perlu membuatnya, `ccam add <name> --shortcut <command>` membuatnya tanpa bertanya, dan `ccam list` menampilkan pintasan setiap akun, atau `-` jika akun tidak punya pintasan.

Nama memakai karakter a-z, 0-9, `.`, `-`, dan `_`. Nama tidak boleh berupa perintah ccam, pintasan akun lain, atau nama program lain di `PATH` Anda.

`ccam rename` mempertahankan nama pintasan, dan `ccam remove` menghapus pintasan bersama akunnya. `ccam shims install` membuat ulang tautan semua akun, misalnya setelah Anda memindahkan ccam.

## Pembagian

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Setiap item menerima `on` (grup `default`), `off`, atau nama grup. Lihat [pembagian](/docs/sharing).

## Menyalin penyiapan yang sudah ada

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` menyalin `settings.json`, `CLAUDE.md`, ekstensi, dan server MCP pengguna dari folder Claude Code ke sebuah grup, sehingga setiap akun yang membagikan item tersebut mendapatkannya. Folder itu hanya dibaca. Opsi `--settings`, `--memory`, `--extensions`, `--mcp`, dan `--history` hanya menyalin item yang Anda sebutkan; tanpa opsi itu, semuanya disalin kecuali riwayat percakapan. `-g` atau `--group` memilih grup, dan `-y` atau `--yes` menyalin tanpa bertanya serta mempertahankan salinan grup jika ada file yang berbeda. Lihat [menyalin penyiapan yang sudah ada](/docs/sharing#copy-an-existing-setup).

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
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
```

`ccam update` memeriksa checksum SHA-256 rilis sebelum mengganti ccam, dan `--version v0.2.0` menginstal rilis tertentu. ccam yang diinstal dari paket `.deb` diperbarui dengan apt. Lihat [pembaruan](/docs/install#update).

## File

| Path | Isi |
| --- | --- |
| `~/.config/ccam/config.toml` | Akun dan apa yang dibagikannya |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` milik setiap akun |
| `~/.local/share/ccam/groups/<group>/` | Pengaturan, server MCP, ekstensi, dan riwayat bersama |
| `~/.local/state/ccam/` | Cache penggunaan dan waktu terakhir digunakan |
| `~/.local/bin/<shortcut>` | Perintah pintasan seperti `ccwork`, berupa tautan ke ccam |

Variabel `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, dan `XDG_STATE_HOME` memindahkan folder-folder ini.
