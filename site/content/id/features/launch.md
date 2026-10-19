---
title: Menjalankan Claude Code
description: Jalankan Claude Code sebagai akun mana pun dengan ccam <name|number>, ccam run, atau pemilih ccam; argumen diteruskan ke claude tanpa diubah, dan --dry-run menampilkan perintahnya.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 401eed57e533
---

`ccam work` menjalankan Claude Code sebagai akun `work`, dan semua yang ada setelah nama diteruskan ke `claude` apa adanya. Gunakan di mana pun Anda biasa mengetik `claude`: di terminal, di skrip, atau dengan input dari pipe.

## Penggunaan

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Berikan nama akun atau nomornya dari `ccam list`. ccam tidak membaca argumen setelah nama, sehingga `ccam work --help` menampilkan bantuan Claude Code. `--` setelah nama bersifat opsional. Input standar, kode keluar, dan sinyal berperilaku seperti pada `claude`.

Jika argumen pertama diawali `-` dan bukan salah satu flag milik ccam, ccam menjalankan akun default dengan semua argumen. Akun default adalah akun yang ditetapkan dengan `ccam use`, atau akun dengan nomor terkecil jika Anda belum menetapkannya.

### Pemilih

`ccam` tanpa argumen membuka pemilih jika input dan output standar adalah terminal, dan menampilkan `ccam list` jika tidak. Pemilih menampilkan nomor, email, penggunaan 5 jam, waktu terakhir digunakan, dan item yang dibagikan untuk setiap akun; angka penggunaan terisi begitu setiap akun menjawab.

| Tombol | Fungsinya |
| --- | --- |
| `Enter` | Menjalankan Claude Code sebagai akun yang dipilih |
| `1` sampai `9` | Menjalankan akun dengan nomor tersebut |
| `Up`/`Down`, `k`/`j` | Memindahkan pilihan |
| `a` | Menutup pemilih dan menjalankan wizard `ccam add` |
| `s` | Mengubah item yang dibagikan akun yang dipilih (`Space` atau `x` mengubah centang, `Enter` menerapkan, `Esc` kembali) |
| `d` | Menjadikan akun yang dipilih sebagai default |
| `r` | Mengambil ulang penggunaan, tanpa cache singkat |
| `q`, `Esc`, `Ctrl+C` | Keluar |

Pemilih ditutup sebelum bertindak, sehingga terminal dipulihkan sebelum Claude Code dimulai. Setelah `d` atau perubahan pembagian, pemilih terbuka lagi. Perubahan pembagian menampilkan rencana dan bertanya terlebih dahulu, seperti `ccam share`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--dry-run` | Menampilkan environment dan perintah yang akan dijalankan ccam, lalu keluar | nonaktif |
| `-v`, `--version` | Menampilkan `ccam <version>` (hanya sebagai argumen pertama) | |
| `-h`, `--help` | Bantuan ccam, sebagai `ccam --help` atau `ccam run --help` | |

Aturan argumen:

- ccam membaca flag miliknya hanya sebelum nama akun. Setelah nama, setiap argumen milik Claude Code, termasuk `--dry-run` dan `--help`.
- `--dry-run` berfungsi sebagai `ccam --dry-run work`, `ccam run --dry-run work`, dan `ccam --dry-run` (akun default). Opsi ini dihormati, bukan ditolak.
- `ccam run <name|number> [claude args]` adalah bentuk panjang dari `ccam <name>`. `ccam run` tanpa nama memakai akun default.
- Nama akun tidak boleh berupa nama perintah ccam seperti `add`, `list`, `run`, atau `usage`, sehingga `ccam <word>` tidak pernah ambigu.

Dry run menampilkan satu baris yang bisa Anda tempel ke shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Cara kerjanya

Setiap peluncuran melewati langkah yang sama, dari baris perintah sampai `exec`:

<!-- motion:launch -->

Setiap kotak menyala saat langkahnya berjalan, dan baris di bawahnya adalah environment dan argumen yang dipakai Claude Code saat dimulai.

1. `cmd/ccam/main.go` memanggil `app.Main`. Jika ccam dijalankan dengan nama lain, itu adalah [pintasan](/features/shortcuts) dan langsung menuju peluncur.
2. `app.parse` memeriksa `os.Args` sebelum cobra terlibat. Hanya `ccam` saja (pemilih), nama perintah ccam (`config.Reserved`), `-h`/`--help`, dan `-v`/`--version` yang diteruskan ke `cli.Execute`. Selain itu adalah peluncuran.
3. `launcher.Launch` membaca `~/.config/ccam/config.toml`, satu-satunya file yang di-parse-nya.
4. `launcher.Prepare` mencari akun berdasarkan nama atau nomor, membuang satu `--` di depan, dan menemukan executable `claude` dengan `ResolveClaude`. Hasil pencarian PATH disimpan dalam cache di `~/.local/state/ccam/claude-path` bersama waktu modifikasi file, sehingga ccam tidak mencari di `PATH` setiap kali peluncuran. File yang berubah atau hilang memicu pencarian baru. Path `claude_bin` di `config.toml` dipakai apa adanya. ccam menolak `claude` yang ternyata ccam sendiri, dibandingkan berdasarkan inode, karena itu akan berulang tanpa akhir.
5. `checkLinks` membandingkan folder akun `~/.local/share/ccam/profiles/<name>/` dengan pengaturan pembagiannya menggunakan `share.Check`, yang hanya memanggil `lstat` dan `readlink`, satu atau dua system call per item. Tautan yang hilang langsung dibuat. Apa pun yang bisa menyebabkan kehilangan data akan menghentikan peluncuran.
6. `applyMCP` menambahkan server MCP bersama (lihat di bawah).
7. `session.Register` mencatat sesi. ccam menulis penanda yang dinamai sesuai PID-nya, berisi waktu mulai proses, ke `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (atau `~/.local/state/ccam/locks/` tanpa `XDG_RUNTIME_DIR`). Untuk akun yang membagikan riwayat percakapan, `historyCheck` terlebih dahulu mencari sesi yang sedang berjalan di `history-<group>-<hash>`, dengan hash dari folder saat ini, lalu juga mendaftar di sana. `exec` mempertahankan PID, sehingga penanda itu kemudian menunjuk proses Claude Code itu sendiri. Penanda, bukan lock yang dipegang, karena proses yang dijalankan Claude Code akan mewarisi lock dan tetap memegangnya setelah sesi berakhir.
8. `Plan.Exec` mengatur judul terminal menjadi `claude · <name>` (jika error standar adalah terminal), menulis waktu peluncuran ke `~/.local/state/ccam/last-used/<name>` tanpa `fsync`, lalu memanggil `syscall.Exec`.

Claude Code berjalan dengan `CLAUDE_CONFIG_DIR` diatur ke folder akun dan `CCAM_PROFILE` diatur ke nama akun. ccam menghapus nilai-nilai ini yang diwarisinya. Untuk folder `~/.claude` yang didaftarkan dengan `ccam import` di macOS, ccam juga mengatur `CLAUDE_SECURESTORAGE_CONFIG_DIR` ke nilai kosong, sehingga akun tersebut dan `claude` biasa memakai login Rantai Kunci (Keychain) yang sama.

### Mengapa exec

`syscall.Exec` mengganti proses ccam dengan `claude`, sehingga ccam tidak tetap berada di memori, penekanan tombol, sinyal, dan kode keluar sepenuhnya milik Claude Code, dan argumen tetap berupa array, bukan string shell. Pekerjaan setelah sesi, seperti mendeteksi server yang ditambahkan dengan `claude mcp add`, dilakukan pada peluncuran berikutnya: ccam membandingkan waktu modifikasi `.claude.json` dengan `~/.local/state/ccam/mcp/<name>/seen` dan mem-parse file itu hanya jika berubah.

### Bagaimana baris perintah diarahkan

Tiga baris perintah, tiga rute melalui `internal/app`:

<!-- motion:launch-route -->

Pintasan dicocokkan pada `argv[0]`, `ccam <name>` pada argumen pertama, dan keduanya berakhir di `syscall.Exec`. Hanya perintah ccam seperti `ccam usage` yang sampai ke cobra.

### Server MCP bersama

Dengan strategi default `flag`, ccam meneruskan server MCP milik grup ke Claude Code sebagai argumen, ditempatkan sebelum argumen Anda:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` hanya mengembalikan argumen jika `mcp.json` milik grup ada, dan `Plan.InsertArgs` menempatkannya tepat setelah `claude`. `--mcp-config` atau `--strict-mcp-config` milik Anda sendiri tetap berlaku di atasnya. Ini adalah satu argumen, `--mcp-config=<path>`, karena `--mcp-config` milik Claude Code menerima beberapa nilai: jika ditulis sebagai dua argumen, prompt yang mengikutinya juga akan dibaca sebagai file konfigurasi.

Dengan `mcp_strategy = "sync"` di `config.toml`, ccam menulis server ke `.claude.json` milik akun sebagai gantinya, dan melewatinya jika akun sudah berjalan.

### Remote Control

`ccam work rc` (atau `ccam work remote-control`) menjalankan Remote Control milik Claude Code sebagai akun `work`. `claude remote-control` menolak `--mcp-config` sebelum subperintahnya dan tidak menerimanya setelahnya, dan sesinya hanya membaca server MCP dari `.claude.json` milik akun. Jadi strategi flag bekerja secara berbeda di sini:

<!-- motion:launch-rc -->

1. `remoteControl` menganggap argumen pertama setelah akun yang tidak diawali `-` sebagai subperintah Claude Code. Jika argumen itu `remote-control` atau `rc`, `mcp.FlagArgs` tidak menambahkan apa pun.
2. `mcp.Synced` membandingkan waktu modifikasi `mcp.json` milik grup dengan `sync-stamp` milik akun. Jika `.claude.json` sudah berisi server grup saat ini, atau grup tidak punya server, Remote Control dimulai tanpa pertanyaan.
3. Jika tidak, ccam menyebutkan server yang tidak akan didapat sesi dan bertanya:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. Jika Anda menjawab `y`, `mcp.Sync` menyalinnya seperti yang dilakukan strategi `sync`: di bawah lock milik Claude Code, setelah membuat cadangan, dan dengan mempertahankan server milik akun sendiri. Jika menjawab `N`, atau tanpa terminal, Remote Control dimulai tanpa server tersebut, dan ccam menyarankan `mcp_strategy = "sync"`, yang menyalinnya sebelum setiap peluncuran.

ccam bertanya pada setiap peluncuran Remote Control sampai server disalin, dan setelah disalin hanya bertanya lagi saat `mcp.json` milik grup berubah lagi. Salinan di `.claude.json` menang atas `--mcp-config` pada peluncuran normal berikutnya, sehingga `ccam mcp add`, `remove`, dan `import` memperbarui salinan yang dimiliki akun lain ([Server MCP bersama](/features/mcp)).

### Anggaran 10 ms

ccam menambahkan paling banyak 10 ms (p95) sebelum Claude Code dimulai. Jalur peluncuran tidak pernah membangun pohon perintah cobra, menjalankan Bubble Tea, atau membuka koneksi jaringan, dan `TestPackageInitIsCheap` gagal jika fungsi `init()` paket memakan waktu lebih dari 5 ms secara total. `scripts/perf-check.sh` mengukur `ccam --dry-run work`, `ccam work`, dan sebuah pintasan dengan `hyperfine` dan gagal jika melebihi 10 ms, dan `BenchmarkLaunchPath` di `internal/launcher` menjalankan langkah-langkahnya dalam loop. Untuk melihat waktu per langkah, atur `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Keamanan dan kasus khusus

- Nama atau nomor yang tidak dikenal berhenti dengan `no account "<name>"` dan mengarahkan ke `ccam list`. Jika belum ada akun sama sekali, ccam meminta Anda menambahkannya dengan `ccam add <name>`.
- Jika path yang dibagikan berubah menjadi file biasa atau mengarah ke tempat yang salah, ccam berhenti dengan `stopped because of share link problems:`, mendaftar masalahnya, dan meminta Anda menjalankan `ccam doctor --fix`. ccam tidak memindahkan data di jalur peluncuran.
- Jika Claude Code lain yang berbagi riwayat percakapan dengan akun ini (grup yang sama) berjalan di folder yang sama, ccam memperingatkan bahwa riwayat bisa tercampur dan bertanya `Run anyway? [y/N]`. Tanpa terminal, ccam menampilkan peringatan dan melanjutkan.
- `--dry-run` melewati sinkronisasi MCP, pemberitahuan server baru, pemeriksaan riwayat, penanda sesi, judul, dan waktu terakhir digunakan, serta tidak melakukan exec. Opsi ini tetap memeriksa tautan pembagian, dan membuat tautan yang hilang seperti peluncuran biasa.
- `ccam work rc` saat `work` sudah berjalan tidak menyentuh `.claude.json`, memberitahukannya, dan menjalankan Remote Control tanpa server grup; jalankan lagi setelah sesi berakhir untuk menyalinnya. Dengan `--dry-run`, ccam menampilkan pesannya tetapi tidak bertanya maupun menyalin.
- Dengan `mcp_strategy = "sync"`, sinkronisasi yang gagal menampilkan `ccam: MCP sync failed (launching anyway): …` dan Claude Code dimulai dengan server yang sudah dimilikinya.

## Terkait

- [Perintah pintasan](/features/shortcuts)
- [Menggunakan Remote Control dengan ccam](/guides/remote-control)
- [Mengelola akun](/features/accounts)
- [Referensi perintah](/docs/commands)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Cara kerja CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
