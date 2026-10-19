---
title: Menggunakan Remote Control dengan ccam
description: Jalankan Remote Control milik Claude Code sebagai akun mana pun dengan ccam <name> rc, lihat opsi mana yang berlaku, mengapa ccam menawarkan untuk menyalin server MCP bersama, dan cara mengatasi masalah umum.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: b7b09f9bcb9e
---

`ccam work rc` menjalankan Remote Control milik Claude Code sebagai akun `work`, sehingga Anda dapat mengendalikan sesi di komputer ini dari claude.ai/code atau aplikasi seluler Claude. Perintah ini bekerja seperti `claude rc` dengan satu perbedaan: jika akun berbagi server MCP, ccam bertanya sekali apakah server tersebut perlu disalin ke `.claude.json` milik akun, karena sesi Remote Control tidak membaca server dari tempat lain.

## Memulai

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[Perintah pintasan](/features/shortcuts) bekerja dengan cara yang sama: `ccwork rc` sama dengan `ccam work rc`. Jalankan perintah di folder tempat Anda ingin bekerja, lalu buka claude.ai/code atau aplikasi Claude yang sudah login dengan akun Claude yang sama seperti `work`.

## Opsi

Ada dua jenis opsi di sini: opsi milik ccam, yang ditulis sebelum nama akun atau berada di pengaturan ccam, dan opsi Remote Control, yang ditulis setelah `rc` dan diteruskan ke Claude Code tanpa diubah.

### Opsi dan pengaturan ccam

| Opsi atau pengaturan | Fungsinya | Default |
| --- | --- | --- |
| `--dry-run` (sebelum akun) | Menampilkan environment dan perintah, serta pemberitahuan MCP jika ada. Tidak bertanya dan tidak menyalin | nonaktif |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Ditanyakan jika akun berbagi server MCP dan `.claude.json` miliknya belum berisi daftar grup saat ini. `y` menyalinnya, `N` atau Enter memulai tanpa server tersebut | `N` |
| `mcp_strategy = "sync"` di `~/.config/ccam/config.toml` | Menyalin server grup ke `.claude.json` sebelum setiap peluncuran, baik Remote Control maupun bukan, sehingga tidak pernah ada pertanyaan | `"flag"` |
| `ccam share work --mcp off` | Berhenti berbagi server MCP dengan akun; Remote Control kemudian dimulai hanya dengan server milik akun sendiri dan tidak ada pertanyaan | aktif |

### Opsi Remote Control

Semua yang ada setelah `rc` adalah milik Claude Code. Berikut opsi `claude remote-control` di Claude Code 2.1.304; `ccam work rc --help` menampilkan daftar untuk versi Anda.

| Opsi | Fungsinya |
| --- | --- |
| `--name <name>` | Nama sesi yang ditampilkan di claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Awalan untuk nama sesi yang dibuat otomatis (default: nama host) |
| `-c`, `--continue` | Menyambung kembali ke sesi yang terakhir dicatat Remote Control untuk folder ini, dalam kurun sekitar 4 jam |
| `--session-id <id>` | Menyambung kembali ke sesi tertentu |
| `--permission-mode <mode>` | Mode izin untuk sesi baru: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk`, atau `plan` |
| `--spawn <mode>` | `same-dir` (default), `worktree` (satu git worktree per sesi), atau `session` (satu sesi, keluar saat sesi berakhir) |
| `--capacity <N>` | Jumlah sesi maksimum sekaligus dalam mode `same-dir` dan `worktree` (default 32) |
| `--[no-]create-session-in-dir` | Membuat sesi pertama di folder saat ini ketika dimulai (default aktif) |
| `--[no-]chrome` | Claude in Chrome untuk sesi baru (default nonaktif) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Output debug dan output rinci |

Tulis `rc` tepat setelah nama akun. ccam menganggap argumen pertama yang tidak diawali `-` sebagai subperintah Claude Code, jadi dalam `ccam work rc --name laptop` ccam melihat `rc`, tetapi dalam `ccam work --model opus rc` ccam melihat `opus`.

## Yang dilakukan ccam

Peluncuran mengikuti jalur biasa ([Menjalankan Claude Code](/features/launch)) dengan satu perubahan pada langkah MCP:

1. ccam menemukan akun, memeriksa tautan berbaginya, dan mencari `claude`, seperti pada peluncuran lainnya.
2. ccam melihat bahwa subperintahnya adalah `rc` atau `remote-control`, lalu tidak menyertakan `--mcp-config`.
3. Jika akun berbagi server MCP, ccam membandingkan `mcp.json` milik grup dengan `sync-stamp` di `~/.local/state/ccam/mcp/work/`. Jika `.claude.json` sudah berisi server saat ini, atau grup tidak punya server, ccam langsung ke langkah 6.
4. Jika tidak, ccam menyebutkan server yang tidak akan didapat sesi:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. Jika Anda menjawab `y`, ccam menyimpan cadangan `.claude.json`, mengambil lock milik Claude Code, menambahkan server grup ke `mcpServers`, dan mempertahankan server lain milik akun. Jika menjawab `N`, tanpa terminal, atau saat akun sedang berjalan, ccam tidak mengubah apa pun dan memberitahukannya:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. ccam mencatat sesi, mengatur `CLAUDE_CONFIG_DIR` ke folder akun, lalu menggantikan dirinya dengan `claude rc` melalui `exec`.

Setelah itu Remote Control terus berjalan sebagai server. Setiap sesi yang dimulainya untuk Anda adalah proses anak dari proses tersebut, sehingga berjalan sebagai `work`: dengan login, pengaturan, dan riwayat percakapan milik `work`, serta dihitung terhadap batas penggunaan `work`. Selama server berjalan, ccam menganggap akun sedang berjalan.

Uji coba (dry run) menunjukkan perbedaannya. Peluncuran normal membawa server grup sebagai argumen, sedangkan Remote Control tidak:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Mengapa cara kerjanya begini

**Remote Control tidak menerima `--mcp-config`.** Pada peluncuran normal, ccam menyerahkan server grup ke Claude Code sebagai `--mcp-config=<path>` dan tidak pernah menulis file milik Claude Code. `claude remote-control` menolak opsi itu sebelum nama subperintahnya dan tidak menerima opsi apa pun setelahnya, sehingga sebelum ccam 0.6.0, `ccam work rc` gagal untuk setiap akun yang berbagi server MCP. Kini ccam tidak menyertakan opsi tersebut.

**Sesinya hanya membaca server dari `.claude.json`.** Sesi yang dimulai Remote Control tidak mewarisi opsi dari baris perintah. Satu-satunya tempat sesi menemukan server MCP pengguna adalah `.claude.json` milik akun. File itu juga menyimpan detail login akun dan state per folder, sehingga ccam tidak dapat menautkannya ke salinan bersama seperti yang dilakukannya pada `settings.json`. Menyalin server ke file tersebut adalah satu-satunya cara memberi sesi daftar bersama.

**ccam bertanya sebelum menulis.** Dengan strategi flag default, ccam berjanji tidak menyentuh `.claude.json`. Remote Control adalah satu-satunya kasus ccam harus menulis ke sana, jadi ccam bertanya dulu, menyimpan cadangan, dan hanya menyentuh `mcpServers`.

**ccam tidak menulis saat akun berjalan.** Claude Code menulis ulang `.claude.json` selama sesi, sehingga perubahan dari luar akan hilang atau tertimpa. Itulah sebabnya ccam melewati penyalinan saat ada sesi akun yang berjalan, termasuk server Remote Control lain.

**ccam bertanya sekali untuk setiap perubahan.** Setelah penyalinan, `sync-stamp` menyimpan waktu modifikasi `mcp.json` milik grup. ccam bertanya lagi hanya jika file itu berubah, misalnya setelah `ccam mcp add`. Jika Anda menjawab `N`, belum ada salinan, jadi ccam bertanya pada setiap peluncuran Remote Control.

**`ccam mcp` menjaga salinan tetap mutakhir.** Entri di `.claude.json` menang atas `--mcp-config` jika namanya bentrok, sehingga salinan lama akan menyembunyikan definisi yang lebih baru pada peluncuran normal juga. Karena itu `ccam mcp add`, `add-json`, `remove`, dan `import` memperbarui salinan di setiap akun yang tidak sedang berjalan:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Tanpa terminal, jawabannya tidak.** Remote Control sering dijalankan dari skrip, layanan, atau `nohup`, tempat tidak ada yang bisa menjawab. ccam lalu memulai tanpa server tersebut alih-alih menunggu.

## Masalah dan solusinya

### Server MCP bersama tidak ada di sesi Remote Control

Pertanyaan dijawab `N`, tidak ada terminal, atau akun sedang berjalan. Hentikan Remote Control, tutup sesi lain milik akun, lalu jalankan lagi `ccam work rc` di terminal; jawab `y`. Agar tidak pernah ditanya, atur `mcp_strategy = "sync"` di `~/.config/ccam/config.toml`. Server dalam sesi yang sedang berjalan tidak berubah, jadi mulai sesi baru setelah penyalinan.

### "Account work is running, so its .claude.json is left alone"

Sesi `work` masih terbuka: `ccam work` biasa, perintah pintasan, atau server Remote Control lain. Keluar dari sesi itu dan jalankan lagi `ccam work rc`. Sesi yang dijalankan akun lain tidak dihitung.

### ccam selalu bertanya

Anda belum menyalin server. Jawab `y` sekali, atau atur `mcp_strategy = "sync"`. Jika Anda tidak pernah ingin server tersebut ada di Remote Control, berhenti berbagi MCP dengan akun itu: `ccam share work --mcp off`.

### `ccam work rc` langsung gagal dengan error tentang `--mcp-config`

ccam Anda lebih lama dari 0.6.0 dan meneruskan `--mcp-config` ke Remote Control. Jalankan `ccam update`. Error yang sama pada ccam versi terbaru berarti ada opsi bernilai sebelum `rc`, seperti dalam `ccam work --model opus rc`, sehingga ccam mengira `opus` adalah subperintahnya; tulis `rc` tepat setelah nama akun.

### Perubahan lewat `ccam mcp` tidak sampai ke Remote Control

`ccam mcp add` dan `remove` hanya memperbarui salinan di akun yang tidak sedang berjalan, dan memberitahukannya: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Hentikan Remote Control lalu mulai lagi; ccam melihat `mcp.json` yang berubah dan bertanya lagi.

### Salinan menggantikan server milik saya sendiri

Jika akun punya server sendiri dengan nama yang sama seperti server grup, salinan menggantikannya. `.claude.json` sebelumnya ada di `~/.local/state/ccam/mcp/work/backups/`, yang menyimpan lima cadangan terbaru. Untuk mempertahankan versi Anda bagi semua akun, pindahkan ke grup dengan `ccam mcp import work <server> --force` sebelum peluncuran Remote Control berikutnya.

### Saya ingin menghapus server yang disalin

`ccam mcp remove <server>` menghapus server dari grup dan dari setiap salinan. Untuk kembali ke keadaan sebelum penyalinan, pulihkan cadangan terbaru dari `~/.local/state/ccam/mcp/work/backups/` saat tidak ada sesi akun yang berjalan.

### Remote Control bertanya apakah folder dipercaya, atau keluar tanpa terminal

Kepercayaan folder disimpan per akun di `.claude.json`, sehingga folder yang Anda percayai sebagai `personal` masih baru bagi `work`. Jalankan `ccam work` sekali di folder itu dan setujui, lalu mulai Remote Control di sana.

### "You must be logged in with a Claude account that has a subscription"

Remote Control memerlukan login Claude dengan langganan. Periksa email akun di `ccam list`, dan login lagi dengan `ccam login work` jika perlu.

### Sesi tidak muncul di ponsel

Aplikasi Claude atau claude.ai/code login dengan akun Claude yang berbeda. Bandingkan dengan email yang ditampilkan `ccam list` untuk `work`.

### ccam memperingatkan bahwa Claude Code lain berbagi riwayat percakapan di folder ini

Akun lain dalam grup riwayat yang sama sedang berjalan di folder yang sama. Sesi yang melanjutkan percakapan yang sama dari kedua akun dapat membuat riwayat tercampur. Jawab `y` hanya jika keduanya akan mengerjakan percakapan yang berbeda, atau mulai Remote Control di folder lain.

## Terkait

- [Menjalankan Claude Code](/features/launch)
- [Server MCP bersama](/features/mcp)
- [Berbagi antarakun](/features/share)
- [Mengelola akun](/features/accounts)
