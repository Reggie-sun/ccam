---
title: Berbagi dengan grup
description: Grup ccam membagikan item pilihan di antara akun-akunnya melalui satu folder; ccam group memindahkan data dengan rencana, cadangan, dan pilihan konflik.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: "3f94dd13dbf1"
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
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
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
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

Anda dapat membuat perubahan yang sama di tab **Groups** pada [dasbor](/features/dashboard); tekan `s` pada sebuah akun untuk memilih grup-grupnya.

## Item

| Item | Yang dicakup |
| --- | --- |
| `mcp` | server MCP, melalui `mcp.json` milik grup ([Server MCP bersama](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` dan `history.jsonl`, sehingga akun lain dapat `--continue` |

Jika tidak ada grup akun itu yang membagikan sebuah item, item itu tetap berada di folder akun itu sendiri.

## Perintah dan opsi

| Perintah | Fungsinya |
| --- | --- |
| `ccam group` atau `ccam group list` | menampilkan daftar grup |
| `ccam group show <group>` | apa yang dibagikan grup, akun-akunnya, foldernya, dan server MCP-nya |
| `ccam group add <group> --items <items> [--members <accounts>]` | membuat grup; item milik akun-akunnya dipindahkan ke dalamnya |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | mengubah apa yang dibagikan grup, untuk setiap akun di dalamnya |
| `ccam group join <group> <account>...` | memasukkan akun ke grup |
| `ccam group leave <group> <account>...` | mengeluarkan akun; masing-masing menyimpan salinan, termasuk server MCP |
| `ccam group rename <group> <new-name>` | mengganti nama hanya di konfigurasi |
| `ccam group remove <group> [--leave]` | menghapus grup |

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: akun yang salinannya masuk ke grup lebih dulu | akun default, lalu menurut nomor |
| `--replace` | `join`: keluar lebih dulu dari grup yang membagikan item yang sama | nonaktif |
| `--fresh` | akun yang berhenti membagikan sebuah item dimulai tanpa item itu, bukan dengan salinannya (folder yang diimpor tetap menyimpan salinan) | nonaktif |
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
2. ccam membandingkan yang lama dan yang baru dan, untuk setiap akun yang itemnya berubah, menyusun rencana (`internal/group`, `share.Build`). Rencana disusun satu per satu, dan masing-masing melihat apa yang akan dimasukkan rencana sebelumnya ke folder grup, sehingga hanya akun pertama yang memindahkan salinannya ke dalam. Salinan pertama grup berasal dari, dengan urutan ini: folder grup (jika sudah memilikinya), akun `--from`, akun default, lalu akun-akun lain menurut nomor.
3. ccam menolak jika salah satu akun itu sedang berjalan, bagaimanapun Claude Code dijalankan (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), dan di Windows tanpa Mode Pengembang jika sebuah file memerlukan tautan.
4. ccam menyelesaikan konflik (dengan bertanya, atau dengan `--prefer`), menampilkan rencana, dan meminta konfirmasi.
5. Tepat sebelum menyimpan, ccam memeriksa lagi. Jika sebuah akun dimulai selagi ccam menunggu jawaban Anda, perubahan berhenti; perubahan juga berhenti jika proses ccam lain mengubah grup yang sama pada saat itu (`the groups changed while this plan waited; nothing was moved, make the change again`). Konfigurasi baru disimpan sebelum apa pun dipindahkan, dan perubahan itu dicatat di `~/.local/state/ccam/group-change.json`. Jika sebuah langkah berhenti di tengah jalan, catatan itu tetap ada: `ccam doctor --fix` menyelesaikan perubahan sesuai urutannya, dengan jawaban yang sudah Anda berikan, termasuk langkah-langkah MCP, dan sampai saat itu tidak ada perubahan grup lain yang dimulai (``the last group change stopped halfway; finish it with `ccam doctor --fix` first``).
6. Langkah-langkah MCP dijalankan lebih dulu. Saat akun mengeluarkan ekstensi dari sebuah grup, daftar plugin grup diarahkan ke folder grup di mana pun path-nya melewati folder sebuah akun (`group plugin lists:` di rencana). Lalu rencana dijalankan secara berurutan. Setiap langkah memindahkan data ke grup atau ke folder cadangan; hanya file yang identik dengan salinan milik grup yang dihapus alih-alih dicadangkan.

ccam menyimpan daftar izin berisi path yang dicakup setiap item (`share.Layout`). Apa pun yang tidak ada di dalamnya tidak akan pernah sampai ke grup. ccam juga menetapkan `.credentials.json` dan `.claude.json` sebagai file yang tidak pernah dibagikan (`share.Never`), dan berhenti dengan galat internal jika salah satunya muncul di tata letak (`Build`). Item `mcp` tidak memiliki path: daftar server berada di dalam `.claude.json`, sehingga [pembagian MCP](/features/mcp) bekerja melalui `mcp.json` milik grup.

Inilah yang terjadi pada satu folder, `projects/`, saat `side` mulai membagikan riwayat dan grup sudah memiliki isi:

<!-- motion:share-on -->

Jika grup belum memiliki apa pun di path itu, akun tersebut adalah yang pertama membagikannya dan foldernya dipindahkan ke sana apa adanya. Jika tidak, setiap file yang belum dimiliki grup dipindahkan ke grup, dan file yang ada di kedua sisi tetap di folder akun, sehingga versi milik grup dipertahankan. File yang identik dengan milik grup dihapus, dan rencana menyebutkan file yang berbeda (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

Folder akun kemudian masuk ke `.ccam-backup/<time>/` di dalam folder akun, dan sebuah tautan menggantikan posisinya. Jika entri milik grup itu sendiri berupa tautan ke entri milik akun yang bergabung, isi akun masuk ke grup menggantikannya. File percakapan diberi nama sesuai ID sesinya, sehingga riwayat dari beberapa akun dapat digabungkan tanpa bentrok. Registri `plugins/` menyimpan path absolut melalui folder akun; ccam mengarahkannya ke folder grup dan menggabungkan entri akun ke dalam registri yang sudah dimiliki grup.

Keluar dari grup membalik tautannya:

<!-- motion:share-off -->

ccam lebih dulu menyalin isi grup saat ini ke `.ccam-backup/incomplete-<name>` di folder akun, dan baru mengganti tautan setelah salinannya lengkap. Jika penyalinan terhenti, tautan tetap di tempatnya dan `ccam doctor --fix` menyalin lagi (rencana penyelesaian semacam itu menampilkan `restore <entry>: … (copied back from the group before a stop)`). Grup tetap menyimpan isinya untuk akun-akun lain.

- ccam hanya menghapus tautan jika tautan itu masih mengarah ke dalam folder grup seperti saat rencana dibuat. Path yang berubah sejak itu, seperti file asli di tempat tautan tadinya berada, dibiarkan apa adanya dan perubahan berhenti.
- Dengan `--fresh`, penyalinan dilewati.
- Akun yang memindahkan sebuah item ke grup yang foldernya masih kosong membawa salinan dari apa yang dilihatnya, sehingga file-filenya tetap ada. Salinan itu disiapkan di `<group>/.ccam-incomplete-<name>` sampai lengkap.
- Akun yang berhenti membagikan MCP mempertahankan server grup: server itu masuk ke `.claude.json` miliknya sendiri (cadangannya disimpan), kecuali jika akun itu sudah memiliki server dengan nama yang sama, yang tetap apa adanya. Dengan `--fresh`, rencana justru mencantumkan server yang berhenti dimuat. Bergabung lagi menghapus salinan yang belum diubah siapa pun.
- Keluar dari `extensions` juga mengarahkan daftar plugin grup ke folder grup di tempat path-nya melewati folder akun yang keluar; bagi akun lain, path yang ditulis ulang itu mengarah ke plugin yang sama.

Satu file yang ada di kedua sisi dengan isi berbeda adalah konflik:

<!-- motion:share-conflict -->

Anda memilih untuk file tunggal seperti `settings.json`; di dalam folder, salinan milik grup dipertahankan, dan rencana menyebutkan file yang berbeda. Pilihannya bergantung pada file; pilihan pertama adalah default:

| File | Pilihan |
| --- | --- |
| `history.jsonl` | `merge by line`, `keep group`, `use mine` |
| File `.json` lain | `keep group`, `use mine`, `merge JSON` |
| Lainnya | `keep group`, `use mine` |

`merge JSON` menggabungkan objek secara rekursif, menyatukan array tanpa duplikat, dan mempertahankan nilai milik grup jika kedua sisi mengatur kunci yang sama secara berbeda, sehingga akun lain tidak melihat perubahan. Setiap pilihan kecuali `keep group` terlebih dahulu menyimpan salinan milik grup ke `.ccam-backup/<time>/group/` milik akun. File milik Anda sendiri selalu masuk ke cadangan.

Pada setiap peluncuran, ccam membandingkan setiap path di tata letak dengan grup-grup akun hanya menggunakan `lstat` dan `readlink` (`share.Check`). Tautan yang hilang langsung dibuat. File asli di tempat yang seharusnya berisi tautan, tautan ke tempat yang salah, atau tautan yang sudah usang akan menghentikan peluncuran dan mengarahkan ke `ccam doctor --fix`.

### Mengganti nama, menghapus, dan format lama

- `ccam group rename` hanya mengubah konfigurasi. Foldernya tetap memakai namanya dan grup mencatatnya sebagai `dir`, sehingga tidak ada yang dipindahkan dan tidak ada yang bisa hilang.
- `ccam group remove` menolak grup yang memiliki akun kecuali Anda menambahkan `--leave`. Setiap akun lalu menyimpan salinan, dan foldernya masuk ke `~/.local/share/ccam/trash/` kecuali grup lain memakainya.
- Sebelum ada grup, `config.toml` menyimpan tabel `[profile.share]` per akun. ccam membaca format itu sebagai grup tanpa memindahkan apa pun. Saat penyimpanan pertama, ccam menulis format baru dan menyimpan file lama sebagai `config.toml.bak`. Item yang dibagikan oleh akun-akun yang sama menjadi satu grup yang dinamai sesuai foldernya; sisanya mendapat nama seperti `default-history` dan tetap memakai folder yang sama.

### Di Windows

Windows hanya membuat tautan simbolik dalam Mode Pengembang (Settings > System > For developers) atau untuk administrator. Tanpa hak itu, folder bersama menjadi directory junction, yang berfungsi sama bagi Claude Code. File bersama (`settings.json`, `CLAUDE.md`, `history.jsonl`) hanya bisa berupa tautan simbolik, jadi perubahan grup yang memerlukannya ditolak sebelum apa pun dipindahkan. Karena alasan yang sama, `ccam add` pertama di sana membuat grup `default` hanya dengan `mcp` dan `extensions` serta memberitahukannya, dan `ccam add --group` ke grup yang membagikan pengaturan, memori, atau riwayat ditolak.

## Keamanan dan kasus khusus

- **Akun yang sedang berjalan.** Perubahan ditolak selama Claude Code sedang berjalan sebagai salah satu akun yang dipindahkannya, karena sesi terus menulis ke file-file itu, entah ccam yang menjalankannya atau bukan (`claude` biasa, IDE, sesi latar belakang). Pemeriksaan diulang tepat sebelum menyimpan. Keluar dan bergabung hanya memengaruhi akun yang disebutkan; keluar dari `extensions` juga menulis ulang path di daftar plugin grup, tanpa mengubah artinya bagi akun lain. `--dry-run` tetap menampilkan rencana.
- **Cadangan.** Setiap rencana mendapatkan folder `.ccam-backup/<time>/` sendiri di folder akun. ccam tidak pernah menghapusnya.
- **Tautan buatan Anda sendiri.** Jika sebuah path sudah berupa tautan yang Anda buat, misalnya ke repositori dotfiles, dan grup masih kosong, ccam menyalin target tautan itu ke grup dan memindahkan tautannya sendiri ke cadangan. Targetnya tidak disentuh.
- **Folder yang diimpor.** Akun yang didaftarkan di tempatnya dengan `ccam import` (seperti `~/.claude`) berubah di tempatnya, sehingga `claude` biasa juga melihat perubahan itu; rencana menyebutkannya. Saat keluar dari grup, akun itu selalu menyimpan salinan, bahkan dengan `--fresh`, dan `ccam remove` mengeluarkannya dari grup-grupnya dengan cara itu sebelum melupakan akun tersebut. File yang dibagikannya berada di folder grup selama akun itu ada di grup, jadi hapus akun itu (atau keluar dari grup-grupnya) sebelum Anda menghapus `~/.local/share/ccam`.
- **Dua akun di satu folder.** Akun yang berbagi riwayat menulis ke file yang sama. Saat akun kedua dimulai di folder tempat satu akun sudah berjalan, ccam memperingatkan dan bertanya.
- **Suntingan manual.** Jika `members` mencantumkan nama yang bukan akun, ccam mengabaikannya dan `ccam doctor` melaporkannya. Nama `local`, `on`, `off`, dan `none` tidak bisa dipakai sebagai nama grup.

## Terkait

- [Dasbor](/features/dashboard)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Server MCP bersama](/features/mcp)
- [Menyalin penyiapan yang sudah ada](/features/migrate)
- [Referensi perintah](/docs/commands)
