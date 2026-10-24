---
title: Menyalin penyiapan yang sudah ada
description: Salin pengaturan, CLAUDE.md, ekstensi, dan server MCP dari ~/.claude yang sudah ada ke grup ccam dengan ccam migrate.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: "7bca9efb6eb0"
---

Jika Anda sudah memakai Claude Code sebelum ccam, penyiapan Anda ada di `~/.claude` dan `~/.claude.json`. `ccam migrate` menyalinnya ke sebuah grup pembagian, sehingga setiap akun yang membagikan item tersebut melalui grup itu dimulai dengan penyiapan yang sudah Anda miliki. Folder sumber hanya dibaca, dan `claude` biasa tetap berfungsi seperti sebelumnya. Jalankan sekali setelah `ccam add` pertama Anda; jika `~/.claude.json` memiliki server MCP, `ccam add` pertama menyebutkannya di tipnya.

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
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Grup hanya menjangkau akun-akunnya sendiri, dan hanya untuk item yang dibagikannya, sehingga ringkasan mendaftar akun-akun dalam grup itu. Jika grup tidak membagikan salah satu item yang disalin, atau ada akun yang tidak berada di dalamnya, ccam menampilkan perintah `ccam group set <group> --add <items>` dan `ccam group join <group> <accounts>` untuk mengubahnya.

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

1. **Menemukan sumber.** `checkSource` menolak path yang tidak ada, bukan folder, atau tidak berisi entri Claude Code yang dikenal. Untuk server MCP saja, `~/.claude` yang kosong atau tidak ada tidak masalah selama `~/.claude.json` memiliki server. `migrate.ClaudeJSON` menemukan `.claude.json` milik folder tersebut: di dalamnya, atau untuk `~/.claude` yang tidak memilikinya di dalam, `~/.claude.json` di sampingnya. Di situlah Claude Code menyimpannya jika `CLAUDE_CONFIG_DIR` tidak diatur.
2. **Memetakan item ke path.** Setiap path berasal dari `share.Layout`, daftar izin yang sama dengan yang dipakai grup. Server MCP dibaca dari `.claude.json` dengan `mcp.LocalServers`, yang memakai pemindai `profile.TopLevelValue` alih-alih mem-parse seluruh file.
3. **Menyalin, tidak pernah memindahkan.** File masuk ke `~/.local/share/ccam/groups/<group>/` dengan `fsutil.AtomicWrite` atau `fsutil.CopyFile`, yang menulis file sementara di samping target lalu mengganti namanya, sehingga salinan yang terhenti tidak pernah meninggalkan file yang terpotong. Server MCP masuk ke `mcp.json` milik grup di bawah lock-nya: server yang ditambahkan ccam lain selama itu tetap ada, dan salinan yang dibuat untuk [Remote Control](/guides/remote-control) diperbarui setelahnya.

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

Server MCP yang masuk dengan OAuth menyimpan tokennya di file login, `.credentials.json` di bawah `mcpOAuth`, bukan di daftar server. ccam tidak menyalinnya. Info masuk yang disalin memegang refresh token yang sama dengan aslinya, dan server merotasi token itu pada pembaruan pertama: semua pemegang lainnya, termasuk `~/.claude` sumber, lalu ter-logout dari server itu. Karena itu baris rencana `sign-ins` menyebutkan, untuk setiap server yang disalin dan sudah dimasuki sumber, akun-akun dalam grup yang perlu masuk sendiri dengan `/mcp`, dan ringkasan mengulanginya. Entri tanpa token bukanlah info masuk dan tidak dicantumkan, begitu juga akun yang sudah memiliki info masuknya sendiri, atau server yang definisinya dipertahankan grup. Tidak ada file login yang ditulis.

## Keamanan dan kasus khusus

- Sumber hanya dibaca. Tidak ada file di dalamnya yang dipindahkan, ditautkan, atau diubah, dan login sumber tidak pernah ditulis, termasuk di macOS, tempat login itu mungkin dibaca dari Rantai Kunci (Keychain).
- Sebelum file grup diganti, oleh `use mine`, penggabungan, atau server MCP baru di `mcp.json` yang sudah ada, salinan milik grup disimpan di `<group>/.ccam-backup/<time>/`. Nama folder itu unik per proses, sehingga dua proses dalam detik yang sama tidak saling menimpa cadangan. Ringkasan menampilkan path-nya.
- Dengan `--yes`, setiap konflik mempertahankan salinan milik grup, dan ccam memberi tahu cara memilih sendiri: jalankan perintah yang sama tanpa `--yes`.
- Jika tidak ada yang berubah, ccam menampilkan `Nothing to copy.`; jika setiap konflik mempertahankan salinan milik grup, `Nothing to copy: the group keeps its own copies.`
- Jika penyalinan gagal di tengah jalan, ccam berhenti dengan `The source folder was not changed. Group files replaced so far are saved in` diikuti folder cadangannya. Menjalankan ulang akan menyalin file yang belum sampai dengan utuh.
- Jika beberapa grup membagikan MCP, perintah yang disarankan `ccam doctor`, `ccam ext`, dan `ccam mcp list` berbunyi `ccam migrate --mcp --group <a|b>`.
- Sumber yang tidak ada berakhir dengan `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Terkait

- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Memeriksa dan memperbaiki](/features/doctor)
- [Beberapa akun Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
