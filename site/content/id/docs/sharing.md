---
title: Berbagi server MCP, pengaturan, dan riwayat
description: Cara grup ccam membagikan server MCP, settings.json, CLAUDE.md, skill, dan riwayat percakapan antarakun Claude Code, dengan setiap login tetap terpisah.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: "2d8029ad6ec8"
---

Setiap akun menyimpan login-nya sendiri. Semua hal lain dapat dibagikan melalui **grup**: sebuah grup menyebutkan item yang dibagikannya dan akun-akun di dalamnya, dan ccam menautkan file setiap akun itu ke folder milik grup, sehingga perubahan yang dibuat di satu akun muncul di semuanya.

## Apa saja yang bisa dibagikan

| Item | Yang dibagikan |
| --- | --- |
| `mcp` | server MCP: `mcp.json` milik grup, diteruskan ke Claude Code saat dimulai |
| `settings` | `settings.json`: izin, hook, model, status line |
| `memory` | `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | riwayat percakapan: `projects/` dan `history.jsonl` |

Akun pertama Anda membuat grup `default`, yang membagikan semuanya kecuali riwayat percakapan, dan akun-akun berikutnya bergabung ke grup itu.

File login `.credentials.json` dan `.claude.json` tidak pernah dibagikan. ccam hanya menautkan path dari daftar tetap, dan kedua file ini tidak ada di dalamnya.

Halaman ini membahas tugas sehari-hari. Setiap opsi dijelaskan di halaman fitur: [berbagi dengan grup](/features/share), [server MCP bersama](/features/mcp), dan [menyalin penyiapan yang sudah ada](/features/migrate).

## Grup

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Satu akun bisa masuk ke beberapa grup, tetapi setiap item hanya berasal dari satu grup. Jadi, untuk berbagi riwayat antara `work` dan `side` tetapi tidak dengan akun Anda yang lain, jangan masukkan riwayat ke `default` dan berikan kedua akun itu grup tersendiri, seperti di atas. `ccam share work` menampilkan apa yang dibagikan sebuah akun dan melalui grup mana. [Dasbor](/features/dashboard) melakukan hal yang sama di tab Groups.

## Saat sebuah akun mulai berbagi

Saat sebuah akun mulai membagikan suatu item (akun itu bergabung ke sebuah grup, atau grupnya menambahkan item tersebut), ccam:

1. menampilkan rencana untuk setiap file yang akan dipindahkan, disalin, atau ditautkan bagi setiap akun, dan meminta persetujuan sebelum memulai;
2. memindahkan file akun ke folder grup jika grup masih kosong, atau menambahkan file yang belum dimiliki grup;
3. jika file yang sama berbeda isinya, menanyakan apakah akan mempertahankan versi grup, memakai versi Anda, atau menggabungkan JSON-nya;
4. menyimpan file asli Anda di `.ccam-backup/<time>/` di dalam folder akun;
5. mengganti file asli dengan tautan ke grup.

File percakapan diberi nama sesuai ID sesinya, sehingga riwayat dari beberapa akun dapat digabungkan tanpa konflik.

ccam tidak mengubah sebuah akun saat Claude Code sedang berjalan sebagai akun tersebut, bagaimanapun Claude Code dijalankan (termasuk sesi latar belakang), karena sesi terus menulis ke file yang sedang dipindahkan. Keluar dari sesi itu terlebih dahulu. ccam memeriksa lagi tepat sebelum menyimpan, sehingga ccam juga mendeteksi sesi yang dimulai selagi ccam menunggu jawaban Anda.

## Saat sebuah akun berhenti berbagi

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Tautan dihapus dan akun mendapatkan salinannya sendiri dari isi grup saat ini; server MCP grup masuk ke `.claude.json` milik akun itu sendiri, di samping server miliknya sendiri. Tambahkan `--fresh` untuk memulai dari kosong; rencana kemudian mencantumkan server MCP yang berhenti dimuat. Folder yang diimpor seperti `~/.claude` selalu menyimpan salinan. Grup tetap menyimpan isinya untuk akun-akun lain.

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
- riwayat percakapan (`projects/` dan `history.jsonl`), tetapi hanya jika Anda menambahkan `--history`.

Folder sumber hanya dibaca: tidak ada isinya yang dipindahkan, ditautkan, atau diubah, dan `claude` biasa tetap berfungsi seperti sebelumnya. File login `.credentials.json` dan `.claude.json` tidak pernah disalin ke grup.

Server MCP yang masuk dengan OAuth menyimpan tokennya di login, bukan di daftar server. ccam tidak menyalin info masuk tersebut: info masuk yang disalin memegang refresh token yang sama, dan pembaruan pertama membuat semua pemegang lainnya ter-logout dari server itu, termasuk `~/.claude`. Rencana dan ringkasan menyebutkan akun-akun yang perlu masuk sendiri ke server-server itu dengan `/mcp`.

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
Group default: personal
To share them:
  ccam group join default work
```

Di akhir, ccam mencantumkan akun-akun dalam grup. Jika grup belum membagikan salah satu item, atau ada akun yang belum masuk ke dalamnya, ccam menampilkan perintah `ccam group` untuk mengubahnya.

Jika grup sudah memiliki salinan file yang berbeda, ccam bertanya, sama seperti saat Anda mengaktifkan pembagian: mempertahankan salinan grup atau memakai milik Anda. `settings.json` dan file JSON lainnya juga bisa digabungkan, dan `history.jsonl` bisa digabungkan baris per baris. Salinan grup sebelumnya disimpan di `.ccam-backup/<time>/` di dalam folder grup.

| Opsi | Efek |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Hanya menyalin item yang disebutkan, tanpa daftar centang. Di luar terminal, tanpa opsi item, ccam menyalin semuanya kecuali riwayat percakapan. |
| `-g`, `--group <group>` | Menyalin ke grup ini, bukan ke `default` |
| `-y`, `--yes` | Menyalin tanpa bertanya; jika ada file yang berbeda, mempertahankan salinan grup |

Tulis `--dry-run` sebelum perintahnya, seperti `ccam --dry-run migrate`, untuk melihat rencana tanpa menyalin apa pun. Untuk menyalin dari folder lain, sebutkan foldernya: `ccam migrate ~/.claude-work --group work`.

## Server MCP

Claude Code menyimpan server MCP di `.claude.json`, bersama dengan login, sehingga file tersebut tidak dapat dibagikan. ccam menyimpan server bersama di `mcp.json` milik grup dan menjalankan Claude Code dengan `--mcp-config=<group>/mcp.json`. Opsi `--mcp-config` Anda sendiri tetap berlaku di atasnya.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Server yang Anda tambahkan di dalam sesi dengan `claude mcp add` hanya disimpan untuk akun tersebut, dan untuk folder tempat Anda berada (cakupan lokal Claude Code; `-s user` membuatnya tersedia bagi akun itu di setiap folder). Saat berikutnya Anda menjalankan akun itu, ccam menyebutkan nama server tersebut, beserta foldernya, dan perintah untuk memindahkannya ke grup:

```sh
ccam mcp import work gh
```

Jika grup sudah memiliki server lain dengan nama yang sama, ccam melewatinya dan memberi tahu Anda. `--force` mengganti server milik grup dan menyimpan `mcp.json` yang lama di folder `.ccam-backup/` milik grup.

Server MCP yang masuk dengan OAuth menyimpan tokennya bersama login akun, jadi masuklah ke server tersebut sekali di setiap akun dengan `/mcp`. `ccam ext` dan `ccam doctor` menyebutkan server yang belum dimasuki sebuah akun sementara akun lain sudah.

Jika Anda lebih suka server muncul di `claude mcp list` dan `/mcp` sebagai server pengguna biasa, tetapkan `mcp_strategy = "sync"` di bagian atas `~/.config/ccam/config.toml`, di atas `[[group]]`, `[[profile]]`, atau `[[folder]]` pertama. ccam kemudian menuliskannya ke `.claude.json` setiap akun sebelum memulai, di bawah file lock milik Claude Code sendiri dan dengan cadangan. Jika akun itu memiliki servernya sendiri dengan nama yang sama dengan server grup, ccam mempertahankan server milik akun dan memberi tahu hal itu.

## Dua akun di folder yang sama

Akun yang berbagi riwayat percakapan menulis ke file yang sama. Jika Anda menjalankan sesi Claude Code kedua di folder tempat akun lain dari grup yang sama sedang berjalan, ccam akan memperingatkan Anda dan menanyakan apakah akan melanjutkan.

## Jika ada yang rusak

```sh
ccam doctor
ccam doctor --fix
```

Jika path bersama berubah menjadi file biasa, misalnya setelah Anda memulihkan cadangan secara manual, ccam berhenti sebelum menjalankan Claude Code dan mengarahkan Anda ke `ccam doctor --fix`. Perbaikan itu menautkan path tersebut lagi ke grup: file itu menjadi salinan grup jika grup belum memilikinya; jika sudah, grup mempertahankan salinannya sendiri, dan file Anda masuk ke cadangan yang disebutkan oleh `--fix`. Perubahan grup yang terhenti di tengah jalan juga diselesaikan oleh `ccam doctor --fix`.
