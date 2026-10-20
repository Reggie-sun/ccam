---
title: Referensi perintah
metaTitle: Referensi perintah ccam
description: Semua perintah dan opsi ccam, dari ccam add dan ccam <name> hingga ccam usage, ccam next, dan ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 116c5b34ac9a
---

Letakkan `--dry-run` sebelum perintah untuk melihat apa yang akan dilakukannya tanpa benar-benar menjalankannya, misalnya `ccam --dry-run work`. Opsi ini berfungsi saat menjalankan Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) dan dengan `group`, `migrate`, `shortcut`, `update`, serta `doctor --fix`. Perintah lain yang mengubah sesuatu akan menolaknya alih-alih berjalan.

<!-- commands -->

## Menjalankan Claude Code

`ccam <name|number> [arguments]` menjalankan Claude Code sebagai akun tersebut. ccam tidak membaca argumen setelah nama akun; ccam meneruskannya ke `claude` dengan urutan yang sama, dan input standar, kode keluar, serta sinyal berperilaku persis seperti pada `claude` itu sendiri. `--` setelah nama akun bersifat opsional.

`ccam` dengan argumen tetapi tanpa nama akun, seperti `ccam -p "summarise"`, menggunakan akun yang ditetapkan untuk folder saat ini, atau jika tidak ada, akun default. Tetapkan keduanya dengan `ccam folder <name> [folder]` dan `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` sebelum nama akun ditujukan untuk proses `-p`/`--print`. Saat akun mencapai batas penggunaannya, ccam menjalankan lagi perintah yang sama sebagai akun yang sudah login dengan sisa kuota terbanyak dan hanya mencetak output dari proses yang selesai. Lihat [Akun per folder](/features/folder) dan [Jalankan ulang saat batas tercapai](/features/next#run-again-at-the-limit).

Pintasan seperti `ccwork` menjalankan akunnya dengan cara yang sama: `ccwork --resume` melakukan hal yang sama dengan `ccam work --resume`. Lihat [pintasan](/docs/commands#shortcuts).

## Menambah dan menghapus akun

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` tanpa folder akan mencari `~/.claude` dan `~/.claude-*`, lalu menanyakan masing-masing. Folder yang diimpor didaftarkan di tempatnya; ccam tidak pernah memindahkan atau menyalin login.

`ccam add --api` menambahkan akun yang memakai kunci API alih-alih login Claude; ccam membaca kunci itu setiap kali akun dijalankan dan tidak pernah menyimpannya. Lihat [akun kunci API](/features/api).

`ccam rename`, `ccam remove`, dan perintah `ccam group` menolak mengubah akun saat Claude Code sedang berjalan sebagai akun tersebut. Keluar dari sesi terlebih dahulu.

`ccam list --json` menampilkan setiap akun beserta slot, pintasan, email, dan foldernya. Array `groups` di dalamnya mencantumkan grup-grup akun tersebut, dan objek `share` memetakan setiap item (`mcp`, `settings`, `memory`, `extensions`, `history`) ke grup asal item itu, atau `local`.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Sebuah grup membagikan item di antara akun-akunnya; satu akun bisa masuk ke beberapa grup, tetapi setiap item hanya berasal dari satu grup. Lihat [berbagi dengan grup](/features/share).

## Menyalin penyiapan yang sudah ada

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` menyalin `settings.json`, `CLAUDE.md`, ekstensi, dan server MCP pengguna dari folder Claude Code ke sebuah grup, sehingga setiap akun yang membagikan item tersebut mendapatkannya. Folder itu hanya dibaca. Opsi `--settings`, `--memory`, `--extensions`, `--mcp`, dan `--history` hanya menyalin item yang Anda sebutkan. Tanpa opsi itu, di terminal sebuah daftar centang memungkinkan Anda memilih, dengan semua item kecuali riwayat percakapan sudah dicentang; di luar terminal, semuanya disalin kecuali riwayat percakapan. `-g` atau `--group` memilih grup, dan `-y` atau `--yes` menyalin tanpa bertanya serta mempertahankan salinan grup jika ada file yang berbeda. ccam juga menyalin info masuk server yang tersimpan ke file login milik setiap akun. Lihat [menyalin penyiapan yang sudah ada](/docs/sharing#copy-an-existing-setup).

## Server MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` memindahkan server yang hanya dimiliki akun itu ke grupnya. Server yang namanya sudah dipakai grup untuk definisi yang berbeda akan dilewati; tambahkan `--force` untuk mengganti salinan milik grup, yang dicadangkan ccam terlebih dahulu.

## Penggunaan dan perpindahan akun

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback`, dan `ccam pick` memilih dengan cara yang sama: akun di bawah `auto_max_usage` (`config.toml`, bawaan 100) lebih dulu, menurut `auto_strategy` (`headroom`, `expiring`, atau `balance`), dengan menghitung jendela mingguan `--model`, tidak pernah akun yang sudah 100%, dan di dalam organisasi akun tempat mereka mulai (`--cross-org` atau `cross_org = true` untuk mencari lebih jauh). Jika tidak ada yang bisa dipilih, semuanya keluar dengan kode 75, atau dengan `--wait` menunggu akun pertama di-reset. `ccam next --continue` membawa percakapan bahkan di antara akun yang tidak berbagi riwayat: Claude Code lalu menyalinnya ke akun baru. Di dalam Claude Code, `/ccam:switch [account]` dari plugin meneruskan percakapan, dan `ccam hook` adalah yang dijalankan oleh hook-nya. Lihat [Beralih ke akun berikutnya](/features/next), [Memilih akun untuk skrip](/features/pick), dan [Sesi latar belakang semua akun](/features/agents).

## Pemeliharaan

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` memeriksa checksum SHA-256 rilis sebelum mengganti ccam, dan `--version v0.2.0` menginstal rilis tertentu. ccam yang diinstal dari paket `.deb` diperbarui dengan apt. Lihat [pembaruan](/docs/install#update).

## File

| Path | Isi |
| --- | --- |
| `~/.config/ccam/config.toml` | Akun, grup, dan pengaturan (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` milik setiap akun |
| `~/.local/share/ccam/groups/<group>/` | Pengaturan, server MCP, ekstensi, dan riwayat bersama |
| `~/.local/state/ccam/` | Cache penggunaan, waktu terakhir digunakan, dan catatan yang ditinggalkan `/ccam:switch` serta hook batas penggunaan (`handoff/`) |
| `~/.local/bin/<shortcut>` | Perintah pintasan seperti `ccwork`, berupa tautan ke ccam |

Variabel `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, dan `XDG_STATE_HOME` memindahkan folder-folder ini. Di Windows, `~` adalah `%USERPROFILE%` dan pintasan berupa `<shortcut>.exe`.
