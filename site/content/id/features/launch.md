---
title: Menjalankan Claude Code
description: Jalankan Claude Code sebagai akun apa pun via ccam <name|number>, ccam run, atau dasbor; argumen diteruskan utuh ke claude, --dry-run menampilkan perintah.
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "afa4942cfed5"
---

`ccam work` menjalankan Claude Code sebagai akun `work`, dan semua yang ada setelah nama diteruskan ke `claude` apa adanya. Gunakan di mana pun Anda biasa mengetik `claude`: di terminal, di skrip, atau dengan input dari pipe.

## Penggunaan

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Berikan nama akun atau nomornya dari `ccam list`. ccam tidak membaca argumen setelah nama, sehingga `ccam work --help` menampilkan bantuan Claude Code. `--` setelah nama bersifat opsional. Input standar, kode keluar, dan sinyal berperilaku seperti pada `claude`.

Jika argumen pertama diawali `-` dan bukan salah satu flag milik ccam, ccam menjalankan akun tanpa nama, dengan semua argumen. Akun itu adalah akun yang ditetapkan untuk folder saat ini dengan [`ccam folder`](/features/folder), atau jika tidak ada, akun default: akun yang ditetapkan dengan `ccam use`, atau akun dengan nomor terkecil jika Anda belum menetapkannya.

### Dasbor

`ccam` tanpa argumen membuka [dasbor](/features/dashboard) jika input dan output standar adalah terminal, dan menampilkan `ccam list` jika tidak. Pindah ke sebuah akun (dengan tombol panah, atau dengan mengetik nomornya) lalu tekan `Enter` untuk menjalankannya; `n` menjalankan akun dengan ruang terbanyak dan `o` menawarkan `--continue`, `--resume`, atau argumen lain. Dasbor ditutup sebelum Claude Code dimulai, sehingga terminal dipulihkan terlebih dahulu.

**Kursi Team dan Enterprise.** Pilihan otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan [`auto_switch`](/features/next#auto-switch) bila Anda mengaktifkannya) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama. Pusat bantuan Anthropic menyatakan "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (setiap kursi pada paket Team diberikan kepada satu orang tertentu, dan hanya orang itu yang boleh menggunakannya) dan "Users can't route their work through another user's account to use that user's remaining usage" (pengguna tidak boleh menyalurkan pekerjaannya melalui akun pengguna lain untuk memakai sisa kuota pengguna tersebut) ([paket Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); begitu pula untuk [paket Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Karena itu, memakai pilihan otomatis di antara kursi dapat melanggar ketentuan Anthropic. Apakah Anda memakainya di sana adalah keputusan Anda: tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya. Paket pribadi (Pro, Max) dibahas dalam [catatan tentang ketentuan di FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Jika shell memiliki `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY`, atau `ANTHROPIC_AUTH_TOKEN` yang ditetapkan, Claude Code akan memakai kunci atau token itu alih-alih login milik akun, sehingga ccam memberitahukannya di error standar sebelum menjalankannya (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). ccam hanya membaca variabel-variabel itu; [akun kunci API](/features/api) justru mengosongkannya.

Dengan `relay = true`, atau [`auto_switch = true`](/features/next#auto-switch) di Linux dan macOS, di bagian atas `config.toml` (di atas `[[group]]`, `[[profile]]`, atau `[[folder]]` pertama), ccam menjalankan sesi interaktif sebagai proses anak dan tetap berjalan di sampingnya, agar `/ccam:switch`, atau `auto_switch` saat batas penggunaan tercapai, bisa melanjutkan percakapan dengan sendirinya, dengan opsi, mode izin, dan tingkat effort sesi tersebut ([detail](/features/next#relay)); jika tidak, ccam menggantikan dirinya dengan Claude Code, dan tidak ada yang berjalan di sampingnya.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--dry-run` | Menampilkan environment dan perintah yang akan dijalankan ccam, lalu keluar | nonaktif |
| `--fallback` | Untuk proses `-p`/`--print`: saat akun mencapai batas penggunaannya, menjalankan lagi perintah yang sama sebagai akun yang sudah login dengan sisa kuota terbanyak ([detail](/features/next#run-again-at-the-limit)) | nonaktif |
| `-v`, `--version` | Menampilkan `ccam <version>` (hanya sebagai argumen pertama) | |
| `-h`, `--help` | Bantuan ccam, sebagai `ccam --help` atau `ccam run --help` | |

Aturan argumen:

- ccam membaca flag miliknya hanya sebelum nama akun. Setelah nama, setiap argumen milik Claude Code, termasuk `--dry-run` dan `--help`.
- `--dry-run` berfungsi sebagai `ccam --dry-run work`, `ccam run --dry-run work`, dan `ccam --dry-run` (akun tanpa nama). Opsi ini dihormati, bukan ditolak. `--dry-run` dan `--fallback` boleh ditulis dalam urutan mana pun.
- `ccam run <name|number> [claude args]` adalah bentuk panjang dari `ccam <name>`. `ccam run` tanpa nama memakai akun folder itu atau akun default.
- Nama akun tidak boleh berupa nama perintah ccam seperti `add`, `list`, `run`, atau `usage`, sehingga `ccam <word>` tidak pernah ambigu.

Dry run menampilkan satu baris yang bisa Anda tempel ke shell, atau ke PowerShell di Windows. Baris itu diawali dengan apa yang diubah peluncuran di environment: selalu `CLAUDE_CONFIG_DIR` dan `CCAM_PROFILE`, `env -u <name>` untuk variabel yang dibuangnya (seperti `CCAM_RELAY` yang diwarisi), dan untuk [akun kunci API](/features/api) endpoint, model, serta kuncinya dalam bentuk sumbernya (`ANTHROPIC_API_KEY="$WORK_KEY"`), tidak pernah kunci itu sendiri:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Baris PowerShell mengembalikan variabel-variabel itu setelahnya, seperti yang dilakukan `VAR=x cmd` di shell, dan menyerahkan baris perintah ke Claude Code melalui `Start-Process`, sehingga Windows PowerShell 5.1 dan PowerShell 7 menjalankannya dengan cara yang sama. Jika peluncuran mengatur sebuah variabel ke nilai kosong, baris `#` di bawahnya memberi tahu bahwa PowerShell sebelum 7.5 justru menghapus variabel seperti itu. Untuk `claude.cmd` yang dipasang npm, baris itu menjalankan perintah `cmd.exe` yang akan dipakai peluncuran (lihat di bawah), dan argumen yang akan ditolak peluncuran di sana juga menggagalkan dry run.

## Cara kerjanya

Setiap peluncuran melewati langkah yang sama, dari baris perintah sampai `exec`:

<!-- motion:launch -->

Setiap kotak menyala saat langkahnya berjalan, dan baris di bawahnya adalah environment dan argumen yang dipakai Claude Code saat dimulai.

1. `cmd/ccam/main.go` memanggil `app.Main`. Jika ccam dijalankan dengan nama lain, itu adalah [pintasan](/features/shortcuts) dan langsung menuju peluncur.
2. `app.parse` memeriksa `os.Args` sebelum cobra terlibat. Hanya `ccam` saja (dasbor), nama perintah ccam (`config.Reserved`), `-h`/`--help`, dan `-v`/`--version` yang diteruskan ke `cli.Execute`. Selain itu adalah peluncuran.
3. `launcher.Launch` membaca `~/.config/ccam/config.toml`, satu-satunya file yang di-parse-nya.
4. `launcher.Prepare` mencari akun berdasarkan nama atau nomor (tanpa nama, `HereProfile` mengambil akun folder itu, dengan membandingkan direktori kerja dengan entri `[[folder]]` di `config.toml` sebagai string, tanpa membedakan huruf besar-kecil di macOS dan Windows), membuang satu `--` di depan, dan menemukan executable `claude` dengan `ResolveClaude`. Hasil pencarian PATH disimpan dalam cache di `~/.local/state/ccam/claude-path` bersama waktu modifikasi file dan hash dari `PATH`, sehingga ccam tidak mencari di `PATH` setiap kali peluncuran. File yang berubah atau hilang, atau `PATH` yang berubah, memicu pencarian baru. Path `claude_bin` di `config.toml` (`~/` diperluas) harus berupa file yang dapat dieksekusi; jika tidak, peluncuran berhenti dengan `claude_bin X is not an executable file. Check claude_bin in config.toml` sebelum apa pun ditulis. ccam menolak `claude` yang ternyata ccam sendiri, dibandingkan berdasarkan inode, karena itu akan berulang tanpa akhir.
5. `checkLinks` membandingkan folder akun `~/.local/share/ccam/profiles/<name>/` dengan pengaturan pembagiannya menggunakan `share.Check`, yang hanya memanggil `lstat` dan `readlink`, satu atau dua system call per item. Tautan yang hilang langsung dibuat. Apa pun yang bisa menyebabkan kehilangan data akan menghentikan peluncuran.
6. `applyMCP` menambahkan server MCP bersama (lihat di bawah).
7. `session.Register` mencatat sesi. ccam menulis penanda yang dinamai sesuai PID-nya, berisi waktu mulai proses, ke `~/.local/state/ccam/locks/profile-<name>/<pid>`, apa pun isi `XDG_RUNTIME_DIR`, sehingga setiap shell melihat penanda yang sama. Untuk akun yang membagikan riwayat percakapan, `historyCheck` terlebih dahulu mencari sesi yang sedang berjalan di `history-<group>-<hash>`, dengan hash dari folder saat ini setelah tautan simbolik di-resolve, lalu juga mendaftar di sana. `exec` mempertahankan PID, sehingga penanda itu kemudian menunjuk proses Claude Code itu sendiri. Penanda, bukan lock yang dipegang, karena proses yang dijalankan Claude Code akan mewarisi lock dan tetap memegangnya setelah sesi berakhir. Penanda yang tidak bisa ditulis tidak menghentikan peluncuran: ccam menampilkan `ccam: could not mark the session as running (…); starting Claude Code anyway`. Pemeriksaan lain ccam tentang apakah sebuah akun sedang berjalan (ganti nama, hapus, perubahan grup, `doctor --fix`, sinkronisasi MCP) juga membaca catatan milik Claude Code sendiri di folder akun (`sessions/<pid>.json`, `daemon.lock`), sehingga juga melihat Claude Code yang dijalankan tanpa ccam: `claude` biasa, IDE, sesi latar belakang, dan daemon-nya.
8. `Plan.Exec` mengatur judul terminal menjadi `claude · <name>` (jika error standar adalah terminal), menulis waktu peluncuran ke `~/.local/state/ccam/last-used/<name>` tanpa `fsync`, lalu memanggil `syscall.Exec`.

Claude Code berjalan dengan `CLAUDE_CONFIG_DIR` diatur ke folder akun dan `CCAM_PROFILE` diatur ke nama akun. ccam menghapus nilai-nilai ini yang diwarisinya, dan juga `CCAM_RELAY` yang diwarisi. Untuk folder `~/.claude` yang didaftarkan dengan `ccam import` di macOS, ccam juga mengatur `CLAUDE_SECURESTORAGE_CONFIG_DIR` ke nilai kosong, sehingga akun tersebut dan `claude` biasa memakai login Rantai Kunci (Keychain) yang sama; di Windows ccam melakukan hal yang sama jika Claude Code menyimpan login di Credential Manager.

### Mengapa exec

`syscall.Exec` mengganti proses ccam dengan `claude`, sehingga ccam tidak tetap berada di memori, penekanan tombol, sinyal, dan kode keluar sepenuhnya milik Claude Code, dan argumen tetap berupa array, bukan string shell.

Windows tidak punya `exec`. Di sana, `Plan.Exec` menjalankan `claude.exe` sebagai proses anak dengan environment, argumen, dan konsol yang sama, menunggunya, lalu keluar dengan kode keluarnya. `claude.cmd` yang dipasang npm berjalan di `cmd.exe` dengan path lengkapnya dan setiap argumen diberi tanda kutip agar `cmd.exe` tidak membaca satu pun sebagai sintaksnya sendiri (`&`, `%PATH%`); argumen yang berisi baris baru tidak bisa diteruskan dengan cara itu dan ditolak, dengan saran untuk memasang Claude Code native. Ctrl+C sampai ke Claude Code langsung dari konsol, jadi ccam mengabaikannya. Penanda sesi berisi PID ccam, yang hidup persis selama Claude Code, dan dihapus saat Claude Code keluar. Selama menunggu, ccam tidak melakukan hal lain. Pekerjaan setelah sesi, seperti mendeteksi server yang ditambahkan dengan `claude mcp add`, dilakukan pada peluncuran berikutnya: ccam membandingkan waktu modifikasi `.claude.json` dengan `~/.local/state/ccam/mcp/<name>/seen` dan mem-parse file itu hanya jika berubah.

### Bagaimana baris perintah diarahkan

Tiga baris perintah, tiga rute melalui `internal/app`:

<!-- motion:launch-route -->

Pintasan dicocokkan pada `argv[0]`, `ccam <name>` pada argumen pertama, dan keduanya berakhir di `syscall.Exec`. Hanya perintah ccam seperti `ccam usage` yang sampai ke cobra.

### Server MCP bersama

Dengan strategi default `flag`, ccam meneruskan server MCP milik grup ke Claude Code sebagai argumen, ditempatkan sebelum argumen Anda:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` hanya mengembalikan argumen jika `mcp.json` milik grup ada, dan `Plan.InsertArgs` menempatkannya tepat setelah `claude`. `--mcp-config` atau `--strict-mcp-config` milik Anda sendiri tetap berlaku di atasnya. Ini adalah satu argumen, `--mcp-config=<path>`, karena `--mcp-config` milik Claude Code menerima beberapa nilai: jika ditulis sebagai dua argumen, prompt yang mengikutinya juga akan dibaca sebagai file konfigurasi.

`claudeVerb` menentukan di mana flag itu tidak boleh ditaruh. Fungsi ini membaca subperintah Claude Code dari argumen seperti yang dilakukan Claude Code, dengan melewati opsi beserta nilainya. `remote-control` atau `rc` tidak mendapat flag (lihat di bawah), begitu juga verb sesi latar belakang yang diberikan pertama (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get`, dan `remove` juga tidak mendapatkannya, karena Claude Code menolaknya di sana; sebagai gantinya ccam memberi tahu di error standar bahwa server grup juga dimuat di sesi-sesi akun itu:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Setiap subperintah lain dan setiap sesi tetap mendapat flag itu.

Dengan `mcp_strategy = "sync"` di `config.toml`, ccam menulis server ke `.claude.json` milik akun sebagai gantinya, dan melewatinya jika akun sudah berjalan. Peluncuran juga melewati pekerjaan itu jika `mcp.json` milik grup maupun `.claude.json` tidak berubah sejak pemeriksaan terakhir. Sinkronisasi tidak pernah mengganti server milik akun sendiri yang memakai nama server grup (ccam memberitahukannya sekali), mengembalikan salinan yang hilang dari `.claude.json` (`put back group default's MCP servers …`), dan hanya mengubah atau menghapus salinan yang ditulisnya sendiri dan belum diubah siapa pun sejak itu.

### Remote Control

`ccam work rc` (atau `ccam work remote-control`) menjalankan Remote Control milik Claude Code sebagai akun `work`. `claude remote-control` menolak `--mcp-config` sebelum subperintahnya dan tidak menerimanya setelahnya, dan sesinya hanya membaca server MCP dari `.claude.json` milik akun. Jadi strategi flag bekerja secara berbeda di sini:

<!-- motion:launch-rc -->

1. `claudeVerb` menemukan `remote-control` atau `rc` sebagai subperintah Claude Code, juga setelah opsi dan nilainya (`ccam work --model opus rc`). Lalu `mcp.FlagArgs` tidak menambahkan apa pun.
2. `mcp.Synced` membandingkan waktu modifikasi `mcp.json` milik grup dengan `sync-stamp` milik akun, dan memeriksa bahwa setiap server grup masih ada di `.claude.json`. Jika `.claude.json` sudah berisi server grup saat ini, atau grup tidak punya server, Remote Control dimulai tanpa pertanyaan.
3. Jika tidak, ccam menyebutkan server yang tidak akan didapat sesi dan bertanya:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

4. Jika Anda menjawab `y`, `mcp.Sync` menyalinnya seperti yang dilakukan strategi `sync`: di bawah lock milik Claude Code, setelah membuat cadangan, dan dengan mempertahankan server milik akun sendiri. Jika menjawab `N`, atau tanpa terminal, Remote Control dimulai tanpa server tersebut, dan ccam menyarankan `mcp_strategy = "sync"`, yang menyalinnya sebelum setiap peluncuran.

ccam bertanya pada setiap peluncuran Remote Control sampai server disalin, dan setelah disalin hanya bertanya lagi saat `mcp.json` milik grup berubah lagi atau sebuah server yang disalin hilang dari `.claude.json`. Salinan di `.claude.json` menang atas `--mcp-config` pada peluncuran normal berikutnya, sehingga `ccam mcp add`, `remove`, dan `import` memperbarui salinan yang dimiliki akun lain ([Server MCP bersama](/features/mcp)).

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
- Jika Claude Code lain yang berbagi riwayat percakapan dengan akun ini (grup yang sama) berjalan di folder yang sama, ccam memperingatkan bahwa riwayat bisa tercampur dan bertanya `Run anyway? [y/N]`. Folder dicocokkan setelah tautan simbolik di-resolve. Tanpa terminal (input standar dari `/dev/null` dihitung sebagai tanpa terminal), atau dari job di latar belakang, ccam menampilkan peringatan dan melanjutkan. Proses `-p` baru memulai percakapannya sendiri dan tidak mendapat peringatan; `--version`, `--help`, dan subperintah yang tidak membuka percakapan, seperti `mcp`, `auth`, atau `plugin`, sama sekali tidak diperiksa.
- `--dry-run` melewati sinkronisasi MCP, pemberitahuan server baru, catatan `mcp list`, pemeriksaan riwayat, penanda sesi, judul, dan waktu terakhir digunakan, serta tidak melakukan exec. Opsi ini tetap memeriksa tautan pembagian, dan membuat tautan yang hilang seperti peluncuran biasa.
- `ccam work rc` saat `work` sudah berjalan tidak menyentuh `.claude.json`, memberitahukannya, dan menjalankan Remote Control tanpa server grup; jalankan lagi setelah sesi berakhir untuk menyalinnya. Dengan `--dry-run`, ccam menampilkan pesannya tetapi tidak bertanya maupun menyalin.
- Dengan `mcp_strategy = "sync"`, sinkronisasi yang gagal menampilkan `ccam: MCP sync failed (launching anyway): …` dan Claude Code dimulai dengan server yang sudah dimilikinya.

## Terkait

- [Perintah pintasan](/features/shortcuts)
- [Menggunakan Remote Control dengan ccam](/guides/remote-control)
- [Mengelola akun](/features/accounts)
- [Referensi perintah](/docs/commands)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Cara kerja CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
