---
title: Berbagi antarakun
description: Cara ccam share menautkan pengaturan, CLAUDE.md, ekstensi, dan riwayat ke folder grup, dengan rencana, cadangan, dan pilihan konflik; off menyalinnya kembali.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Setiap akun memiliki foldernya sendiri, tetapi sebagian besar isinya tidak perlu terpisah. `ccam share` mengaktifkan atau menonaktifkan pembagian untuk satu akun dan satu item sekaligus: file atau folder milik akun dipindahkan ke folder grup dan diganti dengan tautan simbolik, sehingga setiap akun dalam grup membaca dan menulis salinan yang sama. Gunakan setelah `ccam add` untuk mengubah apa yang dibagikan sebuah akun, atau untuk menempatkan beberapa akun dalam grup tersendiri.

## Penggunaan

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Tanpa opsi item, `ccam share` menampilkan keadaan saat ini:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Dengan opsi, ccam menampilkan rencana, menanyakan setiap file yang berbeda dari salinan milik grup, dan bertanya `Proceed? [y/N]` sebelum mengubah apa pun. Akun baru membagikan semuanya kecuali riwayat percakapan; [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing) menjelaskan setiap item dari sisi pengguna.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Server MCP | tidak berubah |
| `--settings <on\|off\|group>` | `settings.json` | tidak berubah |
| `--memory <on\|off\|group>` | `CLAUDE.md` global | tidak berubah |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | tidak berubah |
| `--history <on\|off\|group>` | `projects/` dan `history.jsonl` | tidak berubah |
| `--fresh` | saat menonaktifkan pembagian, mulai dari kosong alih-alih menyalin isi grup | nonaktif |
| `--prefer <group\|mine\|merge>` | cara memutuskan konflik | bertanya; dengan `--yes`, default setiap item |
| `-y`, `--yes` | berjalan tanpa bertanya; konflik mengikuti `--prefer` | nonaktif |
| `-h`, `--help` | bantuan untuk `share` | |
| `--dry-run` (global) | dihormati: menanyakan konflik, menampilkan rencana, tidak mengubah apa pun | nonaktif |

Nilai `on` berarti grup `default`, `off` berarti lokal (tidak dibagikan), dan nilai lain adalah nama grup. Nama grup menjadi nama folder di bawah `~/.local/share/ccam/groups/`, sehingga `ccam share side --history ../x` gagal dengan `invalid group name: "../x"`. Akun diberikan berdasarkan nama atau nomor, dan hanya item yang Anda sebutkan yang berubah.

## Cara kerjanya

File yang terlibat:

- `~/.config/ccam/config.toml`: tabel `[profile.share]` setiap akun, satu nilai per item.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` milik akun.
- `~/.local/share/ccam/groups/<group>/`: salinan bersama.

1. `internal/cli/share.go` membaca konfigurasi, menerapkan setiap nilai `--<item>` ke salinan pengaturan akun dengan `config.Share.Set`, dan menolak jika akun sedang berjalan.
2. `share.Build` membandingkan pengaturan saat ini dengan tujuan dan mendaftar langkah-langkahnya. Fungsi ini hanya membaca disk: `os.Lstat` dan `os.Readlink` untuk setiap path di `share.Layout`. Rencananya berupa daftar operasi: memindahkan ke slot grup yang kosong (adopt), memindahkan file yang belum dimiliki grup (merge dir), menulis konflik yang sudah diputuskan ke grup, mencadangkan, melepas tautan, menautkan, menyalin kembali, menulis ulang path plugin, dan menggabungkan registri plugin.
3. Konflik diputuskan, lewat pertanyaan atau lewat `--prefer`.
4. `Plan.Describe` menampilkan rencana, dan ccam bertanya.
5. Pengaturan tujuan ditulis ke `config.toml` sebelum apa pun dipindahkan. Jika `Plan.Apply` berhenti di tengah jalan, `ccam doctor --fix` dan menjalankan ulang perintah akan menyusun rencananya menuju tujuan yang tersimpan dan menyelesaikan perubahan.
6. `Plan.Apply` menjalankan langkah-langkah secara berurutan. Setiap langkah memindahkan data ke grup atau ke folder cadangan; tidak ada yang dihapus.

`share.Layout` adalah daftar izin berisi path yang dicakup setiap item. Apa pun yang tidak ada di dalamnya tidak akan pernah sampai ke grup. `share.Never` menyebut `.credentials.json` dan `.claude.json`, dan `Build` berhenti dengan galat internal jika salah satunya muncul di tata letak. Item `mcp` tidak memiliki path: daftar server berada di dalam `.claude.json`, sehingga [pembagian MCP](/features/mcp) bekerja melalui `mcp.json` milik grup, dan `--mcp` hanya mengubah pengaturannya.

Inilah yang terjadi pada satu folder, `projects/`, saat `side` mengaktifkan riwayat dan grup sudah memiliki isi:

<!-- motion:share-on -->

Jika grup belum memiliki apa pun di path itu, akun tersebut adalah yang pertama membagikannya dan foldernya dipindahkan ke sana apa adanya, tanpa cadangan. Jika tidak, `mergeDir` memindahkan setiap file yang belum dimiliki grup ke grup, dan file yang ada di kedua sisi tetap di folder akun, sehingga versi milik grup dipertahankan. Folder akun kemudian masuk ke `.ccam-backup/<time>/` di dalam folder akun, dan `os.Symlink` menempatkan tautan di posisinya. File percakapan diberi nama sesuai ID sesinya, sehingga riwayat dari beberapa akun dapat digabungkan tanpa bentrok.

Pemindahan memakai `rename(2)` dan beralih ke salin lalu hapus jika melintasi sistem file. Tautan relatif di dalam pohon yang dipindahkan yang mengarah ke luar pohon itu ditulis ulang menjadi tautan absolut (`fsutil.PortableLink`), sehingga tetap terhubung dari tempat barunya. Registri `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) menyimpan path absolut melalui folder akun; ccam mengarahkannya ke folder grup dan, jika grup sudah memiliki registri, menggabungkan entri akun ke dalamnya.

Menonaktifkan sebuah item membalik tautannya:

<!-- motion:share-off -->

ccam hanya menghapus tautan jika tautan itu mengarah ke dalam folder grup, lalu menyalin isi grup saat ini ke akun dengan `CopyTree`. Dengan `--fresh`, penyalinan dilewati dan akun dimulai tanpa apa pun di path itu. Grup tetap menyimpan isinya untuk akun-akun lain. Memindahkan item dari satu grup ke grup lain melepas tautannya dan menautkannya ke grup baru, dengan membiarkan grup lama apa adanya; rencana menyebutkannya dalam baris `Note:`.

Satu file yang ada di kedua sisi dengan isi berbeda adalah konflik:

<!-- motion:share-conflict -->

`share.NewConflict` mendaftar pilihan yang sesuai dengan file tersebut. `history.jsonl` menawarkan `merge by line`, `keep group`, dan `use mine`, dengan penggabungan di urutan pertama. File `.json` lain menawarkan `keep group`, `use mine`, dan `merge JSON`. Selain itu menawarkan `keep group` dan `use mine`. Pilihan pertama adalah default:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) menggabungkan objek secara rekursif, menyatukan array tanpa duplikat, dan mempertahankan nilai milik grup jika kedua sisi mengatur kunci yang sama secara berbeda, sehingga akun lain tidak melihat perubahan perilaku. `merge by line` mempertahankan baris-baris milik grup, lalu menambahkan baris Anda yang belum dimilikinya. Setiap pilihan kecuali `keep group` terlebih dahulu menyimpan salinan milik grup ke `.ccam-backup/<time>/group/` di folder akun, lalu menulis hasilnya dengan `fsutil.AtomicWrite`. File milik Anda sendiri selalu masuk ke cadangan.

`--prefer group`, `mine`, atau `merge` memutuskan setiap konflik tanpa bertanya; `merge` berarti `merge JSON` atau `merge by line`, mana pun yang ditawarkan file tersebut, dan file yang tidak menawarkan keduanya tetap memakai default. `--yes` saja mengambil default setiap file.

Pada setiap peluncuran, `share.Check` membandingkan setiap path di tata letak dengan pengaturannya hanya menggunakan `lstat` dan `readlink`, sehingga cukup ringan untuk jalur peluncuran. Tautan yang hilang langsung dibuat, karena tidak ada yang bisa hilang. File asli di tempat yang seharusnya berisi tautan, tautan ke tempat yang salah, atau tautan yang sudah usang akan menghentikan peluncuran:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` menjalankan `share.Build` dengan pengaturan tersimpan di kedua sisi, yang menurunkan setiap langkah dari apa yang ada di disk, mempertahankan isi grup saat terjadi konflik, dan mencadangkan apa pun yang digantinya.

## Keamanan dan kasus khusus

- **Akun yang sedang berjalan.** ccam menolak mengubah pembagian saat Claude Code sedang berjalan sebagai akun tersebut (`account work is running; quit it first, then change its sharing`), karena sesi terus menulis ke file yang sedang dipindahkan dan baris yang ditambahkannya pada saat itu akan hilang. `--dry-run` tetap menampilkan rencana.
- **Cadangan.** Setiap rencana mendapatkan folder `.ccam-backup/<time>/` sendiri di folder akun, dengan akhiran bernomor jika dua rencana dimulai pada detik yang sama. ccam tidak pernah menghapusnya.
- **Kegagalan sebagian.** Jika sebuah langkah gagal, ccam menampilkan galatnya dan lokasi file aslinya, lalu meminta Anda menyelesaikannya dengan `ccam doctor --fix`. Tujuannya sudah tersimpan, sehingga perbaikan melanjutkan ke arah itu.
- **Tautan buatan Anda sendiri.** Jika sebuah path sudah berupa tautan yang Anda buat, misalnya ke repositori dotfiles, dan grup masih kosong, ccam menyalin target tautan itu ke grup dan memindahkan tautannya sendiri ke cadangan. Targetnya tidak disentuh.
- **File login.** `.credentials.json` dan `.claude.json` tidak pernah ditautkan, dipindahkan, atau disalin ke grup.
- **Dua akun di satu folder.** Akun yang berbagi riwayat menulis ke file yang sama. Saat akun kedua dari grup riwayat yang sama dimulai di folder tempat satu akun sudah berjalan, ccam memperingatkan dan bertanya sebelum melanjutkan.

## Terkait

- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Server MCP bersama](/features/mcp)
- [Menyalin penyiapan yang sudah ada](/features/migrate)
- [Semua perintah](/docs/commands)
