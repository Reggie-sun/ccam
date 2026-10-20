---
title: Berbagi dengan grup
description: Grup ccam membagikan item pilihan di antara akun-akunnya melalui satu folder; ccam group memindahkan data dengan rencana, cadangan, dan pilihan konflik.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

Setiap akun memiliki foldernya sendiri, tetapi sebagian besar isinya tidak perlu terpisah. Sebuah **grup** menentukan apa yang dibagikan dan oleh siapa: setiap akun dalam grup menautkan item-item itu ke satu folder, `~/.local/share/ccam/groups/<group>/`, sehingga semuanya membaca dan menulis salinan yang sama. Token login tidak pernah dibagikan.

Satu akun bisa masuk ke beberapa grup, tetapi setiap item hanya berasal dari satu grup. Penyiapan yang umum adalah satu grup untuk pengaturan dan ekstensi, dan grup kedua untuk riwayat percakapan antara dua akun:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Akun pertama yang Anda tambahkan membuat grup `default`, dan akun-akun berikutnya bergabung ke grup itu kecuali Anda menentukan lain.

## Penggunaan

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Setiap perintah yang memindahkan data menampilkan rencana, menanyakan setiap file yang berbeda dari salinan milik grup, dan bertanya `Proceed? [y/N]` sebelum mengubah apa pun:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

Perubahan yang sama tersedia di tab **Groups** pada [dasbor](/features/dashboard), dan `s` pada sebuah akun memilih grup-grupnya.

## Item

| Item | Yang dicakup |
| --- | --- |
| `mcp` | server MCP, melalui `mcp.json` milik grup ([Server MCP bersama](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` dan `history.jsonl`, sehingga akun lain dapat `--continue` |

Item yang tidak dibagikan oleh grup mana pun dari akun tersebut tetap berada di folder akun itu sendiri.

## Perintah dan opsi

| Perintah | Fungsinya |
| --- | --- |
| `ccam group` atau `ccam group list` | menampilkan daftar grup |
| `ccam group show <group>` | apa yang dibagikan grup, akun-akunnya, foldernya, dan server MCP-nya |
| `ccam group add <group> --items <items> [--members <accounts>]` | membuat grup; item milik akun-akunnya dipindahkan ke dalamnya |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | mengubah apa yang dibagikan grup, untuk setiap akun di dalamnya |
| `ccam group join <group> <account>...` | memasukkan akun ke grup |
| `ccam group leave <group> <account>...` | mengeluarkan akun; masing-masing menyimpan salinan |
| `ccam group rename <group> <new-name>` | mengganti nama hanya di konfigurasi |
| `ccam group remove <group> [--leave]` | menghapus grup |

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: akun yang salinannya masuk ke grup lebih dulu | akun default, lalu menurut nomor |
| `--replace` | `join`: keluar lebih dulu dari grup yang membagikan item yang sama | nonaktif |
| `--fresh` | akun yang berhenti membagikan sebuah item dimulai tanpa item itu, bukan dengan salinannya | nonaktif |
| `--leave` | `remove`: mengeluarkan akun-akun grup terlebih dahulu | diperlukan jika grup memiliki akun |
| `--prefer <group\|mine\|merge>` | cara memutuskan konflik | bertanya; dengan `--yes`, default setiap file |
| `-y`, `--yes` | berjalan tanpa bertanya; konflik mengikuti `--prefer` | nonaktif |
| `--dry-run` (global) | menampilkan rencana, tidak mengubah apa pun | nonaktif |

Akun diberikan berdasarkan nama atau nomor; item sebagai daftar yang dipisahkan koma. `ccam add` menerima `--group <group>` (bisa diulang) atau `--no-group`.

## Cara kerjanya

File yang terlibat:

- `~/.config/ccam/config.toml`: entri `[[group]]`: `name`, `items`, `members` dan, jika nama foldernya berbeda, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` milik akun.
- `~/.local/share/ccam/groups/<folder>/`: salinan bersama.

Saat `config.toml` dibaca, grup-grup setiap akun diuraikan menjadi satu folder per item. Ini memakan waktu sekitar 25 µs untuk lima akun, sehingga peluncuran tetap hanya membaca `config.toml`. Konfigurasi yang membuat sebuah akun mendapatkan satu item dari dua grup adalah galat, dan tidak ada yang dimulai sampai hal itu diperbaiki.

Perubahan dikerjakan pada salinan konfigurasi:

1. Perintah mengubah salinan itu (`AddGroup`, `SetGroupItems`, `JoinGroup` …), yang memeriksa aturan-aturannya.
2. `internal/group` membandingkan yang lama dan yang baru dan, untuk setiap akun yang itemnya berubah, menyusun rencana dengan `share.Build`. Rencana disusun satu per satu, dan masing-masing melihat apa yang akan dimasukkan rencana sebelumnya ke folder grup, sehingga hanya akun pertama yang memindahkan salinannya ke dalam. Salinan pertama berasal dari folder grup jika folder itu sudah memilikinya, lalu dari `--from`, lalu dari akun default, lalu dari akun lain menurut nomor.
3. ccam menolak jika salah satu akun itu sedang berjalan, dan di Windows tanpa Mode Pengembang jika sebuah file memerlukan tautan.
4. Konflik diputuskan, lewat pertanyaan atau lewat `--prefer`, rencana ditampilkan, dan ccam bertanya.
5. Konfigurasi baru disimpan sebelum apa pun dipindahkan. Jika sebuah langkah berhenti di tengah jalan, `ccam doctor --fix` dan menjalankan ulang perintah akan menyusun rencananya menuju grup yang tersimpan dan menyelesaikan perubahan.
6. Rencana dijalankan secara berurutan. Setiap langkah memindahkan data ke grup atau ke folder cadangan; tidak ada yang dihapus.

`share.Layout` adalah daftar izin berisi path yang dicakup setiap item. Apa pun yang tidak ada di dalamnya tidak akan pernah sampai ke grup. `share.Never` menyebut `.credentials.json` dan `.claude.json`, dan `Build` berhenti dengan galat internal jika salah satunya muncul di tata letak. Item `mcp` tidak memiliki path: daftar server berada di dalam `.claude.json`, sehingga [pembagian MCP](/features/mcp) bekerja melalui `mcp.json` milik grup.

Inilah yang terjadi pada satu folder, `projects/`, saat `side` mulai membagikan riwayat dan grup sudah memiliki isi:

<!-- motion:share-on -->

Jika grup belum memiliki apa pun di path itu, akun tersebut adalah yang pertama membagikannya dan foldernya dipindahkan ke sana apa adanya. Jika tidak, setiap file yang belum dimiliki grup dipindahkan ke grup, dan file yang ada di kedua sisi tetap di folder akun, sehingga versi milik grup dipertahankan. Folder akun kemudian masuk ke `.ccam-backup/<time>/` di dalam folder akun, dan sebuah tautan menggantikan posisinya. File percakapan diberi nama sesuai ID sesinya, sehingga riwayat dari beberapa akun dapat digabungkan tanpa bentrok. Registri `plugins/` menyimpan path absolut melalui folder akun; ccam mengarahkannya ke folder grup dan menggabungkan entri akun ke dalam registri yang sudah dimiliki grup.

Keluar dari grup membalik tautannya:

<!-- motion:share-off -->

ccam hanya menghapus tautan jika tautan itu mengarah ke dalam folder grup, lalu menyalin isi grup saat ini ke akun. Dengan `--fresh`, penyalinan dilewati. Grup tetap menyimpan isinya untuk akun-akun lain. Akun yang memindahkan sebuah item ke grup yang foldernya masih kosong membawa salinan dari apa yang dilihatnya, sehingga file-filenya tetap ada; server MCP ikut terbawa dengan cara yang sama. Akun yang berhenti membagikan MCP kehilangan server yang sebelumnya disalin ccam ke `.claude.json` miliknya (cadangannya disimpan), karena server itu berasal dari grup.

Satu file yang ada di kedua sisi dengan isi berbeda adalah konflik:

<!-- motion:share-conflict -->

`history.jsonl` menawarkan `merge by line`, `keep group`, dan `use mine`, dengan penggabungan di urutan pertama. File `.json` lain menawarkan `keep group`, `use mine`, dan `merge JSON`. Selain itu menawarkan `keep group` dan `use mine`. `merge JSON` menggabungkan objek secara rekursif, menyatukan array tanpa duplikat, dan mempertahankan nilai milik grup jika kedua sisi mengatur kunci yang sama secara berbeda, sehingga akun lain tidak melihat perubahan. Setiap pilihan kecuali `keep group` terlebih dahulu menyimpan salinan milik grup ke `.ccam-backup/<time>/group/` milik akun. File milik Anda sendiri selalu masuk ke cadangan.

Pada setiap peluncuran, `share.Check` membandingkan setiap path di tata letak dengan grup-grup akun hanya menggunakan `lstat` dan `readlink`. Tautan yang hilang langsung dibuat. File asli di tempat yang seharusnya berisi tautan, tautan ke tempat yang salah, atau tautan yang sudah usang akan menghentikan peluncuran dan mengarahkan ke `ccam doctor --fix`.

### Mengganti nama, menghapus, dan format lama

- `ccam group rename` hanya mengubah konfigurasi. Foldernya tetap memakai namanya dan grup mencatatnya sebagai `dir`, sehingga tidak ada yang dipindahkan dan tidak ada yang bisa hilang.
- `ccam group remove` menolak grup yang memiliki akun kecuali Anda menambahkan `--leave`. Setiap akun lalu menyimpan salinan, dan foldernya masuk ke `~/.local/share/ccam/trash/` kecuali grup lain memakainya.
- Sebelum ada grup, `config.toml` menyimpan tabel `[profile.share]` per akun. ccam membaca format itu sebagai grup tanpa memindahkan apa pun. Saat penyimpanan pertama, ccam menulis format baru dan menyimpan file lama sebagai `config.toml.bak`. Item yang dibagikan oleh akun-akun yang sama menjadi satu grup yang dinamai sesuai foldernya; sisanya mendapat nama seperti `default-history` dan tetap memakai folder yang sama.

### Di Windows

Windows hanya membuat tautan simbolik dalam Mode Pengembang (Settings > System > For developers) atau untuk administrator. Tanpa hak itu, folder bersama menjadi directory junction, yang berfungsi sama bagi Claude Code. File bersama (`settings.json`, `CLAUDE.md`, `history.jsonl`) hanya bisa berupa tautan simbolik, jadi perubahan grup yang memerlukannya ditolak sebelum apa pun dipindahkan.

## Keamanan dan kasus khusus

- **Akun yang sedang berjalan.** Perubahan ditolak selama Claude Code berjalan sebagai salah satu akun yang dipindahkannya, karena sesi terus menulis ke file-file itu. Keluar dan bergabung hanya menyangkut akun yang disebutkan. `--dry-run` tetap menampilkan rencana.
- **Cadangan.** Setiap rencana mendapatkan folder `.ccam-backup/<time>/` sendiri di folder akun. ccam tidak pernah menghapusnya.
- **Tautan buatan Anda sendiri.** Jika sebuah path sudah berupa tautan yang Anda buat, misalnya ke repositori dotfiles, dan grup masih kosong, ccam menyalin target tautan itu ke grup dan memindahkan tautannya sendiri ke cadangan. Targetnya tidak disentuh.
- **Folder yang diimpor.** Akun yang didaftarkan di tempatnya dengan `ccam import` (seperti `~/.claude`) berubah di tempatnya, sehingga `claude` biasa juga melihat perubahan itu; rencana menyebutkannya.
- **Dua akun di satu folder.** Akun yang berbagi riwayat menulis ke file yang sama. Saat akun kedua dimulai di folder tempat satu akun sudah berjalan, ccam memperingatkan dan bertanya.
- **Suntingan manual.** Anggota yang bukan akun diabaikan dan dilaporkan oleh `ccam doctor`. Nama `local`, `on`, `off`, dan `none` tidak bisa dipakai sebagai nama grup.

## Terkait

- [Dasbor](/features/dashboard)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Server MCP bersama](/features/mcp)
- [Menyalin penyiapan yang sudah ada](/features/migrate)
- [Semua perintah](/docs/commands)
