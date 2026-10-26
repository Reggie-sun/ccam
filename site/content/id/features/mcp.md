---
title: Server MCP bersama
description: Simpan satu daftar server MCP per grup dengan ccam mcp add, add-json, list, remove, dan import, diteruskan lewat --mcp-config atau disinkronkan ke setiap akun.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: "27a06ad35295"
---

Claude Code menyimpan server MCP pengguna di `.claude.json`, file yang sama yang menyimpan login akun, sehingga file itu tidak dapat dibagikan antarakun. ccam menyimpan server bersama di file grup, `mcp.json`, dan menyerahkannya ke setiap akun dalam grup saat peluncuran. Gunakan `ccam mcp` untuk menambahkan server sekali untuk semua akun Anda, dan `ccam mcp import` untuk memindahkan server yang Anda tambahkan di dalam satu sesi ke akun-akun lain.

## Penggunaan

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
```

Perubahan berlaku mulai peluncuran berikutnya dari setiap akun; sesi yang sedang berjalan tetap memakai server yang dipakainya saat dimulai:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

Di bawah grup-grup, `ccam mcp list` menyebutkan setiap server yang tidak sampai ke akun lain:

- `only in account work: pg` untuk server yang ditambahkan sebuah akun untuk dirinya sendiri;
- `only in account side, for /home/me/code/app only: notes` untuk server yang ditambahkan untuk satu folder (`claude mcp add` biasa, lihat di bawah);
- `account work has its own …, used in its sessions in place of the group's` jika server milik akun sendiri memakai nama server grup;
- `.mcp.json` folder ini, beserta akun-akun yang persetujuannya masih `waiting for approval in`;
- server yang `in no account, only in ~/.claude.json (plain claude's)`, beserta perintah [`ccam migrate --mcp`](/features/migrate) untuk menyalinnya.

Menambahkan nama yang sudah dimiliki grup akan mengganti definisinya dan menampilkan `updated` alih-alih `added`. Server hanya sampai ke akun yang berada dalam grup tersebut, dan hanya jika grup itu membagikan MCP (`ccam group set <group> --add mcp`); jika grup itu tidak membagikannya, `add` akan memberitahukannya.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `-g`, `--group <group>` | grup tujuan untuk `add`, `add-json`, dan `remove` | satu-satunya grup yang dipakai akun untuk membagikan MCP, jika tidak `default` |
| `-- <command> [args…]` (`add`) | menjalankan perintah ini sebagai server stdio | |
| `--url <url>` (`add`) | URL server jarak jauh (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transport untuk server `--url` | `http` |
| `--header <'Name: value'>` (`add`) | header HTTP, bisa diulang | tidak ada |
| `-e`, `--env <KEY=VALUE>` (`add`) | variabel environment untuk server stdio, bisa diulang | tidak ada |
| `--force` (`import`) | mengganti server grup bernama sama yang berbeda; `mcp.json` lama dicadangkan | nonaktif |
| `--keep-local` (`import`) | mempertahankan entri di `.claude.json` milik akun | nonaktif |
| `-h`, `--help` | bantuan untuk setiap perintah | |
| `--dry-run` (global) | hanya didukung oleh `mcp list`; `add`, `add-json`, `remove`, dan `import` menolaknya dan tidak mengubah apa pun | nonaktif |

`add` memerlukan perintah setelah `--` atau `--url`; jika tidak, perintah ini gagal dengan `a command to run (after --) or --url is required`. `--header` berlaku untuk server `--url` dan `--env` untuk server perintah. `add-json` menerima JSON yang sama dengan `claude mcp add-json`, dan definisinya harus berupa objek dengan `command` atau `url`. `import` selalu memakai grup yang dipakai akun untuk membagikan MCP, dan `list` menampilkan setiap grup, sehingga `-g` tidak mengubah keduanya.

## Cara kerjanya

File yang terlibat:

- `~/.local/share/ccam/groups/<group>/mcp.json`: server bersama, `{"mcpServers": {…}}`, format yang dibaca `--mcp-config`. File ini ditulis dengan `fsutil.AtomicWrite` dan mode `0600`, karena entri server sering membawa kunci API, dan hanya di bawah lock `<group>/.mcp.json.lock`, sehingga dua perintah `ccam mcp add` yang dijalankan bersamaan sama-sama mempertahankan servernya.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: file milik Claude Code sendiri yang berisi server akun. ccam hanya membaca nilai `mcpServers`-nya dengan pemindai `profile.TopLevelValue` dan tidak pernah mem-parse seluruh file di jalur peluncuran.
- `~/.local/state/ccam/mcp/<name>/`: file `seen`, `sync-stamp` (dengan hash dari setiap salinan yang ditulis ccam, agar bisa membedakan salinan yang tidak diubah siapa pun dari salinan yang Anda sunting), `sync-checked` (apa yang terakhir dilihat strategi sync), `kept` (server yang dipertahankan akun saat keluar dari berbagi MCP), dan cadangan `.claude.json`.

Ada dua strategi, diatur oleh `mcp_strategy` di `~/.config/ccam/config.toml`. Default-nya adalah strategi flag:

<!-- motion:mcp-flag -->

1. ccam mencari grup MCP milik akun (`launcher.applyMCP`) dan memeriksa dengan satu `stat` bahwa `mcp.json` milik grup ada, tanpa mem-parse-nya (`mcp.FlagArgs`).
2. ccam menempatkan `--mcp-config=<path>` tepat setelah `claude`, sebelum argumen Anda sendiri, sehingga `--mcp-config` atau `--strict-mcp-config` milik Anda tetap berlaku di atasnya (`Plan.InsertArgs`).
3. ccam memeriksa apakah `.claude.json` milik akun berubah sejak peluncuran terakhir, dan hanya dalam hal itu mem-parse servernya (`mcp.NewLocal`).
4. ccam menggantikan dirinya dengan Claude Code (`syscall.Exec`).

Path diteruskan sebagai satu argumen, `--mcp-config=<path>`, bukan dua. `--mcp-config` milik Claude Code menerima beberapa nilai, sehingga `claude --mcp-config a.json "prompt"` akan membaca prompt sebagai file konfigurasi kedua dan gagal dengan `MCP config file not found`. Bentuk satu argumen menghindari hal ini, sehingga prompt atau subperintah setelahnya tetap berfungsi. `ccam work mcp list`, `get`, dan `remove` tidak mendapat flag, karena Claude Code menolaknya di sana; sebagai gantinya ccam menyebutkan server grup di error standar ([detail](/features/launch#shared-mcp-servers)). Claude Code memuat server `--mcp-config` sebagai tambahan dari server milik akun; jika ada nama yang bentrok, entri di `.claude.json` yang menang. Strategi flag hanya menulis `.claude.json` untuk Remote Control, saat sebuah akun keluar dari berbagi MCP atau bergabung lagi ([di bawah](#leaving-and-joining)), dan untuk `ccam mcp import`.

`claude remote-control` (`rc`) tidak menerima `--mcp-config`, dan sesinya hanya membaca server dari `.claude.json`. Untuk `ccam work rc`, ccam tidak menyertakan flag tersebut dan, di terminal, bertanya apakah server grup perlu disalin ke `.claude.json`, dengan cadangan, seperti yang dilakukan strategi `sync`. ccam terus bertanya sampai salinan ada, lalu bertanya lagi hanya setelah `mcp.json` milik grup berubah atau sebuah server yang disalin hilang ([Remote Control](/features/launch#remote-control)). Karena salinan semacam itu menang atas `--mcp-config` jika ada nama yang bentrok, `ccam mcp add`, `add-json`, `remove`, dan `import` memperbarui salinan yang dimiliki akun lain dalam grup:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Baris `✓ Updated` hanya muncul jika sebuah salinan berubah. Salinan yang dibuat versi ccam yang lebih lama diperbarui dengan cara yang sama, begitu pula salinan setelah `ccam migrate`.

Server yang ditambahkan di dalam sesi dengan `claude mcp add` hanya masuk ke `.claude.json` milik akun itu. `claude mcp add` biasa memakai cakupan lokal: server dimuat untuk satu akun itu di satu folder itu saja (`-s user` membuatnya tersedia bagi akun itu di setiap folder). ccam tidak dapat melihat hal ini saat terjadi, karena Claude Code sudah menggantikan proses ccam, sehingga ccam memeriksanya pada peluncuran berikutnya. File `seen` menyimpan waktu modifikasi `.claude.json` dan nama server yang sudah dilaporkan. Jika waktunya berbeda, ccam membandingkan nama-nama itu dengan grup dan menampilkan petunjuk untuk setiap server baru satu kali, tanpa menunggu input, pada kedua strategi:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Petunjuk itu juga menyebutkan server milik akun sendiri yang memakai nama server grup, karena milik akun itulah yang menang di sesinya.

Dengan `mcp_strategy = "sync"`, server muncul di `claude mcp list` dan `/mcp` sebagai server pengguna biasa, dan ccam menulisnya ke `.claude.json` setiap akun sebelum peluncuran:

<!-- motion:mcp-sync -->

1. ccam membandingkan waktu modifikasi `mcp.json` milik grup dengan `sync-stamp`, dan `.claude.json` milik akun dengan apa yang terakhir dilihatnya (`mcp.Sync`). Jika keduanya tidak berubah, tidak ada hal lain yang terjadi. Jika berubah, fungsi ini juga memeriksa bahwa setiap server grup masih ada di `.claude.json`, dan mengembalikan salinan yang hilang (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Fungsi ini mengambil lock milik Claude Code sendiri, yaitu folder `.claude.json.lock` di samping file, yang dibuat dengan `mkdir`. Fungsi ini menunggu paling lama 2 detik dan menganggap lock yang berumur lebih dari 10 detik sudah basi, seperti yang dilakukan pustaka lock milik Claude Code.
3. Fungsi ini menyimpan `.claude.json` saat ini ke `backups/` di folder state dan mempertahankan 5 yang terbaru.
4. Fungsi ini mengatur setiap server grup di `mcpServers`, menghapus server yang pernah ditulisnya tetapi kini tidak dimiliki grup, dan menulis file secara atomik dengan mode aslinya, hanya jika ada yang berubah. Fungsi ini hanya mengubah atau menghapus salinan yang ditulisnya sendiri dan belum diubah siapa pun sejak itu. Jika akun memiliki servernya sendiri dengan nama yang sama dengan server grup, ccam mempertahankan server milik akun dan memberitahukannya sekali: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Fungsi ini mencatat stamp baru beserta hash dari setiap salinan, lalu Claude Code dimulai tanpa `--mcp-config`.

Sinkronisasi dilewati, dengan `ccam: account work is already running, skipping MCP sync`, jika ada sesi akun yang sedang berjalan, bagaimanapun Claude Code dijalankan, karena Claude Code menulis ulang `.claude.json` selama berjalan. Sinkronisasi yang gagal menampilkan `ccam: MCP sync failed (launching anyway): …` dan tetap meluncurkan.

### Keluar dan bergabung

Akun yang berhenti membagikan MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) mempertahankan server grup: rencana menyalinnya ke `.claude.json` milik akun itu, dengan cadangan, dan server miliknya sendiri dengan nama yang sama tetap ada. Dengan `--fresh`, akun tidak mendapat satu pun, dan rencana mencantumkan server yang berhenti dimuat. Saat akun itu bergabung lagi ke grup yang memiliki server tersebut, salinan yang dipertahankan dan belum diubah siapa pun dihapus, sehingga suntingan grup berikutnya sampai ke akun itu.

`ccam mcp import` memindahkan server milik akun sendiri ke grupnya:

<!-- motion:mcp-import -->

ccam membaca server akun, miliknya sendiri dan yang dimilikinya hanya untuk satu folder, dan membandingkan masing-masing dengan server grup bernama sama menggunakan `mcp.SameDef`, yang mengabaikan urutan kunci dan spasi. Server folder menjadi server grup untuk setiap folder (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); satu nama dengan dua definisi folder yang berbeda ditolak. Server plugin tidak ditawarkan, karena plugin itulah yang membawanya. Server yang belum dimiliki grup, atau dimiliki dengan definisi yang sama, disalin. Server yang berbeda dilewati, karena definisi milik grup itulah yang dipakai setiap akun lain:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Dengan `--force`, server grup diganti, dan `mcp.json` sebelumnya terlebih dahulu disimpan ke `.ccam-backup/<time>/mcp.json` di dalam folder grup. Lalu ccam menghapus entri yang sudah disalin dari `.claude.json` milik akun, miliknya sendiri dan milik folder dalam satu penulisan, di bawah lock yang sama dan setelah membuat cadangan di folder state (`mcp.Remove`). ccam tetap mengunci `mcp.json` milik grup sejak membaca file itu sampai entri-entri itu dihapus, sehingga tidak ada yang menghapus server dari grup di sela-selanya. Ini perlu karena entri milik akun akan menang atas entri milik grup jika namanya bentrok. Dengan `--keep-local`, entri-entri itu tetap ada di `.claude.json` milik akun.

## Keamanan dan kasus khusus

- **Info masuk.** Server yang masuk dengan OAuth menyimpan tokennya di `.credentials.json` setiap akun (`mcpOAuth`), bukan di daftar server, sehingga membagikan server tidak ikut membagikan info masuknya. Masuk sekali per akun dengan `/mcp`. ccam tidak pernah menyalin info masuk, begitu pula [`ccam migrate`](/features/migrate): refresh token yang disalin membuat semua pemegang lainnya ter-logout. [`ccam ext`](/features/ext) dan `ccam doctor` menyebutkan server yang belum dimasuki sebuah akun sementara akun lain sudah, dan `ccam next --continue` memberi tahu server mana yang tidak dimiliki akun baru. Grup tidak pernah menerima file login.
- **Akun yang sedang berjalan.** `import` menyalin ke grup tetapi membiarkan entri milik akun yang sedang berjalan dan memberitahukannya; jalankan lagi setelah sesi berakhir. Sinkronisasi melewati akun yang sedang berjalan. Keduanya juga mendeteksi sesi Claude Code yang dijalankan tanpa ccam.
- **Salinan untuk Remote Control.** Salinan milik akun yang sedang berjalan tidak disentuh: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Dengan `mcp_strategy = "sync"`, tidak ada yang perlu diperbarui, karena setiap peluncuran melakukan sinkronisasi.
- **File rusak.** Jika `.claude.json` bukan JSON yang valid, sinkronisasi tidak menyentuhnya dan melaporkan galatnya. Jika lock tidak dapat diambil dalam 2 detik, ccam melaporkan `could not lock … (Claude Code is writing the file)`.
- **Akun yang tidak membagikan MCP.** `import` menolak dan meminta Anda memasukkan akun itu ke grup yang membagikan MCP terlebih dahulu (`ccam group join <group> side`).
- **Nama grup.** `-g` hanya menerima nama grup; nama itu menjadi folder di bawah `groups/`, sehingga nama seperti `../x` ditolak.

## Terkait

- [Berbagi dengan grup](/features/share)
- [Menggunakan Remote Control dengan ccam](/guides/remote-control)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Menyalin penyiapan yang sudah ada](/features/migrate)
- [Referensi perintah](/docs/commands)
