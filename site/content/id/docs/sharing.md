---
title: Berbagi server MCP, pengaturan, dan riwayat
description: Cara ccam membagikan server MCP, settings.json, CLAUDE.md, agen, skill, dan riwayat percakapan antarakun Claude Code dengan login tetap terpisah.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 888d466b29e9
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

Halaman ini membahas tugas sehari-hari. Setiap opsi dijelaskan di halaman fitur: [berbagi antarakun](/features/share), [server MCP bersama](/features/mcp), dan [menyalin penyiapan yang sudah ada](/features/migrate).

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

ccam tidak mengubah pembagian saat Claude Code sedang berjalan sebagai akun tersebut, karena sesi terus menulis ke file yang sedang dipindahkan. Keluar dari sesi itu terlebih dahulu.

## Menonaktifkan pembagian

```sh
ccam share work --history off
```

Tautan dihapus dan akun mendapatkan salinannya sendiri dari isi grup saat ini. Tambahkan `--fresh` untuk memulai dari kosong. Grup tetap menyimpan isinya untuk akun-akun lain.

## Menyalin penyiapan yang sudah ada

Jika Anda sudah memakai Claude Code sebelum ccam, pengaturan Anda ada di `~/.claude`. `ccam migrate` menyalinnya ke sebuah grup, sehingga setiap akun yang membagikan item tersebut mendapatkannya:

```sh
ccam migrate
```

Di terminal, ccam lebih dulu menampilkan daftar centang berisi item-item tersebut. Semua item kecuali riwayat percakapan sudah dicentang; tekan Space untuk mengubah sebuah item dan Enter untuk melihat rencananya. Item yang tidak ada isinya di folder itu ditandai `not found`.

Yang disalin:

- `settings.json` dan `CLAUDE.md` global;
- ekstensi `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, dan `rules/`;
- server MCP pengguna Anda, yaitu entri `mcpServers` di `~/.claude.json`, ke `mcp.json` milik grup;
- info masuk (OAuth) yang tersimpan untuk server MCP yang disalin, ke file login milik setiap akun;
- riwayat percakapan (`projects/` dan `history.jsonl`), tetapi hanya jika Anda menambahkan `--history`.

Folder sumber hanya dibaca: tidak ada isinya yang dipindahkan, ditautkan, atau diubah, dan `claude` biasa tetap berfungsi seperti sebelumnya. File login `.credentials.json` dan `.claude.json` tidak pernah disalin ke grup.

Server MCP yang masuk dengan OAuth menyimpan tokennya di login, bukan di daftar server. Agar Anda tidak perlu masuk lagi ke server itu, ccam menyalin info masuk tersebut ke file login setiap akun yang mendapatkan server melalui grup, kecuali akun itu sudah memiliki info masuknya sendiri. ccam melewati akun yang sedang berjalan, belum login, didaftarkan di tempatnya dengan `ccam import`, atau menyimpan loginnya di Rantai Kunci (Keychain) macOS, lalu menyebutkan nama akun-akun itu; di akun tersebut, masuklah ke server dengan `/mcp`. Akun yang baru mulai membagikan server MCP belakangan akan mendapatkan info masuk itu saat Anda menjalankan `ccam migrate` lagi.

Sebelum menyalin apa pun, ccam menampilkan rencananya dan meminta persetujuan:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Di akhir, ccam mencantumkan akun yang membagikan setiap item melalui grup tersebut. Jika ada akun yang belum, ccam menampilkan perintah `ccam share` untuk mengaktifkan pembagian bagi akun itu.

Jika grup sudah memiliki salinan file yang berbeda, ccam bertanya, sama seperti saat Anda mengaktifkan pembagian: mempertahankan salinan grup atau memakai milik Anda. `settings.json` dan file JSON lainnya juga bisa digabungkan, dan `history.jsonl` bisa digabungkan baris per baris. Salinan grup sebelumnya disimpan di `.ccam-backup/<time>/` di dalam folder grup.

| Opsi | Efek |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Hanya menyalin item yang Anda sebutkan, tanpa daftar centang; tanpa opsi ini dan di luar terminal, semuanya kecuali riwayat percakapan |
| `-g`, `--group <group>` | Menyalin ke grup ini, bukan ke `default` |
| `-y`, `--yes` | Menyalin tanpa bertanya; jika ada file yang berbeda, mempertahankan salinan grup |

Tulis `--dry-run` sebelum perintahnya, seperti `ccam --dry-run migrate`, untuk melihat rencana tanpa menyalin apa pun. Untuk menyalin dari folder lain, sebutkan foldernya: `ccam migrate ~/.claude-work --group work`.

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

Jika grup sudah memiliki server lain dengan nama yang sama, ccam melewatinya dan memberi tahu Anda. `--force` mengganti server milik grup dan menyimpan `mcp.json` yang lama di folder `.ccam-backup/` milik grup.

Server MCP yang masuk dengan OAuth menyimpan tokennya bersama login akun, jadi otorisasi server tersebut sekali di setiap akun, atau salin info masuknya dari `~/.claude` dengan `ccam migrate`.

Jika Anda lebih suka server muncul di `claude mcp list` dan `/mcp` sebagai server pengguna biasa, tetapkan `mcp_strategy = "sync"` di `~/.config/ccam/config.toml`. ccam kemudian menuliskannya ke `.claude.json` setiap akun sebelum memulai, di bawah file lock milik Claude Code sendiri dan dengan cadangan.

## Dua akun di folder yang sama

Akun yang berbagi riwayat percakapan menulis ke file yang sama. Jika Anda menjalankan Claude Code kedua di folder tempat akun lain dari grup yang sama sedang berjalan, ccam akan memperingatkan Anda dan menanyakan apakah akan melanjutkan.

## Jika ada yang rusak

```sh
ccam doctor
ccam doctor --fix
```

Jika path bersama berubah menjadi file biasa, misalnya setelah Anda memulihkan cadangan secara manual, ccam berhenti sebelum menjalankan Claude Code dan mengarahkan Anda ke `ccam doctor --fix`, yang menggabungkan file tersebut ke grup dan menyimpan cadangan.
