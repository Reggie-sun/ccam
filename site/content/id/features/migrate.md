---
title: Menyalin penyiapan yang sudah ada
description: Salin pengaturan, CLAUDE.md, ekstensi, server MCP, dan info masuk MCP dari ~/.claude yang sudah ada ke grup pembagian ccam dengan ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

Jika Anda sudah memakai Claude Code sebelum ccam, penyiapan Anda ada di `~/.claude` dan `~/.claude.json`. `ccam migrate` menyalinnya ke sebuah grup pembagian, sehingga setiap akun yang membagikan item tersebut melalui grup itu dimulai dengan penyiapan yang sudah Anda miliki. Folder sumber hanya dibaca, dan `claude` biasa tetap berfungsi seperti sebelumnya. Jalankan sekali setelah `ccam add` pertama Anda, dan lagi setiap kali sebuah akun mulai membagikan server MCP.

## Penggunaan

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

Di terminal, tanpa opsi item dan tanpa `--yes`, ccam terlebih dahulu menampilkan daftar centang:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (atau `x`) mengubah centang item di bawah kursor, tombol panah atau `j`/`k` memindahkan kursor, Enter menampilkan rencana, dan Esc, `q`, atau Ctrl+C menampilkan `Cancelled.`. Item yang tidak ada isinya di sumber ditandai `not found`. Jika Anda menghapus semua centang, ccam menampilkan `Nothing selected.` dan berhenti.

Lalu ccam menampilkan rencana dan bertanya sebelum menyalin:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Grup hanya menjangkau akun yang membagikan suatu item melaluinya, sehingga ringkasan mendaftar akun-akun itu per item. Untuk akun yang tidak membagikannya, ccam menampilkan perintah `ccam share <name> --<item> on` untuk mengaktifkan pembagian.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `[folder]` | Folder Claude Code yang menjadi sumber salinan | `~/.claude` |
| `--settings` | Menyalin `settings.json` | nonaktif |
| `--memory` | Menyalin `CLAUDE.md` global | nonaktif |
| `--extensions` | Menyalin `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | nonaktif |
| `--mcp` | Menyalin server MCP pengguna (`mcpServers` di `.claude.json`) | nonaktif |
| `--history` | Menyalin riwayat percakapan (`projects/`, `history.jsonl`) | nonaktif |
| `-g`, `--group string` | Grup pembagian tujuan salinan | `default` |
| `-y`, `--yes` | Menyalin tanpa bertanya; jika ada konflik, salinan milik grup dipertahankan | nonaktif |
| `--dry-run` (global) | Menampilkan rencana dan berhenti sebelum pertanyaan atau penyalinan apa pun | nonaktif |

Tanpa opsi item, daftar centang yang menentukan di terminal. Di luar terminal, atau dengan `--yes`, ccam menyalin semuanya kecuali riwayat percakapan. Riwayat tidak disertakan secara default karena `projects/` bisa berukuran besar.

`--group on` berarti `default`. `--group off` dan `--group local` ditolak dengan `--group takes a group name, not "off"`. `ccam migrate` menghormati `--dry-run`, dan daftar centang tetap muncul saat dry run sehingga Anda dapat meninjau pilihan terlebih dahulu.

## Cara kerjanya

`internal/cli/migrate.go` memilih item, lalu `migrate.Build` membaca sumber dan grup serta mengembalikan rencana. Tidak ada yang ditulis sampai `Plan.Apply`, setelah Anda mengonfirmasi.

1. **Menemukan sumber.** `checkSource` menolak path yang tidak ada, bukan folder, atau tidak berisi entri Claude Code yang dikenal. `migrate.ClaudeJSON` menemukan `.claude.json` milik folder tersebut: di dalamnya, atau untuk `~/.claude` yang tidak memilikinya di dalam, `~/.claude.json` di sampingnya. Di situlah Claude Code menyimpannya jika `CLAUDE_CONFIG_DIR` tidak diatur.
2. **Memetakan item ke path.** Setiap path berasal dari `share.Layout`, daftar izin yang sama dengan yang dipakai `ccam share`. Server MCP dibaca dari `.claude.json` dengan `mcp.LocalServers`, yang memakai pemindai `profile.TopLevelValue` alih-alih mem-parse seluruh file.
3. **Menyalin, tidak pernah memindahkan.** File masuk ke `~/.local/share/ccam/groups/<group>/` dengan `fsutil.AtomicWrite` atau `fsutil.CopyFile`; server MCP masuk ke `mcp.json` milik grup.

Setiap item memetakan satu path sumber ke satu path grup:

<!-- motion:migrate-items -->

Dua baris riwayat digambar putus-putus karena hanya disalin dengan `--history`. `.credentials.json` dan `.claude.json` tercantum di `share.Never` dan tidak pernah sampai ke grup: `Build` mengembalikan galat internal jika salah satunya muncul di tata letak, dan penelusuran folder melewatinya berdasarkan nama.

`buildEntry` memberi setiap path sebuah keadaan, yang dijelaskan oleh setiap baris rencana:

<!-- motion:migrate-states -->

- **Missing**: tidak ada di sumber, atau berupa file padahal yang diharapkan folder. Dicantumkan sekali per item sebagai `not found, skipped`.
- **Linked**: path sumber sudah mengarah ke salinan milik grup, misalnya jika sumbernya adalah folder profil ccam.
- **Same**: grup sudah memiliki isi yang identik.
- **New**: grup belum memilikinya, sehingga disalin utuh.
- **Merge**: folder di kedua sisi. `copyMissing` hanya menambahkan file yang belum dimiliki grup; file yang berbeda tetap seperti milik grup dan rencana menyebut `group copies kept`. Untuk `plugins/`, registri `installed_plugins.json` dan `known_marketplaces.json` digabungkan, dan path plugin absolut ditulis ulang agar mengarah ke folder grup.
- **Conflict**: file di kedua sisi dengan isi berbeda. Anda memilih per file: `keep group`, `use mine`, ditambah `merge JSON` untuk file `.json` (penggabungan mendalam yang memenangkan grup jika nilainya bentrok) atau `merge by line` untuk `history.jsonl`. Server MCP yang definisinya berbeda juga merupakan konflik, dengan pilihan `keep group` atau `use mine`. Pilihan pertama yang tercantum adalah default saat Anda menekan Enter; `--yes` selalu mempertahankan salinan milik grup.

Daftar centang itu sendiri adalah `tui.MigratePicker` yang dibangun di atas komponen `checklist` yang juga dipakai `ccam add`; `migrate.Found` menandai item yang tidak ada dengan panggilan lstat dan daftar server MCP.

### Info masuk MCP

Server MCP yang masuk dengan OAuth menyimpan tokennya di file login, `.credentials.json` di bawah `mcpOAuth`, bukan di daftar server. Claude Code memberi kunci setiap entri berupa `<name>|` ditambah 16 karakter pertama hash `sha256` dari `type`, `url`, dan `headers` server. Kunci itu tidak bergantung pada folder, sehingga entri yang disalin berfungsi di akun lain selama definisi servernya identik. Karena itu `buildLogins` di `internal/migrate/logins.go` hanya menyalin entri server yang akhirnya dimiliki grup dengan definisi dari sumber: baru, sudah sama, atau konflik yang Anda putuskan dengan `use mine`. Info masuk server plugin (kunci yang diawali `plugin:`) ikut disalin saat Anda menyalin ekstensi dan sumbernya memiliki `plugins/`.

<!-- motion:migrate-signins -->

Entri itu masuk ke `.credentials.json` milik setiap akun penerima, ditulis dengan `fsutil.AtomicWrite` dengan mode `0600`. Entri itu tidak pernah masuk ke grup. Akun hanya menerimanya jika akun itu membagikan `mcp` (atau, untuk server plugin, `extensions`) melalui grup tujuan, dan hanya kunci yang belum dimilikinya yang ditambahkan: entri yang sudah dimilikinya, yang mungkin lebih baru, dipertahankan. Baris rencana `sign-ins` mendaftar server, akun penerima, dan akun yang dilewati. ccam melewati sebuah akun, dan menyebutkan namanya beserta alasannya, jika akun itu:

- didaftarkan di tempatnya dengan `ccam import`, karena `claude` biasa mungkin sedang memakai file login-nya;
- sedang berjalan, karena Claude Code menulis ulang file login-nya (diperiksa lagi tepat sebelum menulis, dilaporkan sebagai `started meanwhile`);
- belum login;
- menyimpan login-nya di Rantai Kunci (Keychain) macOS, yang tidak ditulis oleh ccam.

Masuklah ke server di akun yang dilewati dengan `/mcp`. Akun yang baru mulai membagikan server MCP belakangan akan mendapatkan info masuk itu saat Anda menjalankan `ccam migrate` lagi; ringkasan menyebutkannya dengan ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Keamanan dan kasus khusus

- Sumber hanya dibaca. Tidak ada file di dalamnya yang dipindahkan, ditautkan, atau diubah, dan login sumber tidak pernah ditulis, termasuk di macOS, tempat login itu mungkin dibaca dari Rantai Kunci (Keychain).
- Sebelum file grup diganti, oleh `use mine`, penggabungan, atau server MCP baru di `mcp.json` yang sudah ada, salinan milik grup disimpan di `<group>/.ccam-backup/<time>/`. Nama folder itu unik per proses, sehingga dua proses dalam detik yang sama tidak saling menimpa cadangan. Ringkasan menampilkan path-nya.
- Dengan `--yes`, setiap konflik mempertahankan salinan milik grup, dan ccam memberi tahu cara memilih sendiri: jalankan perintah yang sama tanpa `--yes`.
- Jika tidak ada yang berubah, ccam menampilkan `Nothing to copy.`; jika setiap konflik mempertahankan salinan milik grup, `Nothing to copy: the group keeps its own copies.`
- Jika penyalinan gagal di tengah jalan, ccam berhenti dengan `The source folder was not changed. Group files replaced so far are saved in` diikuti folder cadangannya.
- Sumber yang tidak ada berakhir dengan `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Terkait

- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Memeriksa dan memperbaiki](/features/doctor)
- [Beberapa akun Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
