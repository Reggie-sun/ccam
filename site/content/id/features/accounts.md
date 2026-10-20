---
title: Mengelola akun
description: Tambahkan, login, tampilkan, ganti nama, hapus, dan impor akun Claude Code dengan ccam add, ccam login, ccam list, ccam rename, ccam remove, dan ccam import.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

Setiap akun ccam adalah folder konfigurasi Claude Code tersendiri dengan login-nya sendiri. Perintah-perintah ini membuat, memasukkan login, menampilkan, mengganti nama, dan menghapus akun, serta mendaftarkan folder yang sudah Anda miliki.

## Penggunaan

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Di mana pun sebuah perintah menerima `<name|number>`, Anda dapat memberikan nama akun atau nomor slotnya, yaitu kolom `#` pada `ccam list`.

`ccam add` di terminal tanpa opsi akan membuka wizard. Wizard menanyakan nama dan cara akun itu login (login Claude sekarang, kunci API, atau login Claude nanti), lalu menampilkan grup yang akan diikuti dengan `default` sudah dicentang (sebelum ada grup, yang ditampilkan adalah item yang akan dibagikan grup `default`), lalu bertanya `Create a shortcut command for work? [y/N]`, dan diakhiri dengan ringkasan untuk dikonfirmasi. `Esc` kembali satu langkah. Salah satu dari `--group`, `--no-group`, `--shortcut`, atau `--no-login`, atau pipe sebagai ganti terminal, akan melewati wizard; dalam hal itu nama wajib diberikan dan tidak ada pintasan yang dibuat kecuali `--shortcut` menyebutkannya.

Tanpa `--group`, akun baru bergabung ke grup `default`. Akun pertama membuat grup itu, yang membagikan server MCP, pengaturan, `CLAUDE.md`, dan ekstensi; lihat [berbagi dengan grup](/features/share).

`ccam list` menampilkan slot, nama, pintasan, email, grup, dan waktu terakhir digunakan, dengan `*` di depan akun default:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

Dengan `--json`, setiap akun adalah objek dengan `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (dihilangkan jika akun belum pernah dijalankan), `groups`, dan `share`. `share` memetakan setiap item (`mcp`, `settings`, `memory`, `extensions`, `history`) ke grup asalnya, atau ke `"local"`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `ccam add --no-login` | Membuat akun tanpa login; login belakangan dengan `ccam login <name>` | nonaktif |
| `ccam add --group <group>` | Bergabung ke grup ini alih-alih `default`. Bisa diulang; `default` dibuat jika belum ada | `default` |
| `ccam add --no-group` | Tidak bergabung ke grup mana pun: akun tidak membagikan apa pun | nonaktif |
| `ccam add --shortcut <command>` | Juga membuat perintah pintasan ini di `~/.local/bin` | tidak ada |
| `ccam add --api --key-from <source>` | Akun yang memakai kunci API alih-alih login; lihat [akun kunci API](/features/api) | nonaktif |
| `ccam list --json` | Menampilkan akun sebagai JSON | nonaktif |
| `ccam remove -y`, `--yes` | Menghapus tanpa bertanya | nonaktif |
| `ccam import --name <name>` | Nama akun untuk folder yang Anda berikan | diambil dari nama folder |
| `ccam import -y`, `--yes` | Mengimpor setiap folder yang ditemukan tanpa bertanya | nonaktif |
| `--dry-run` | Dihormati oleh `ccam list` dan `ccam completion`, yang tidak mengubah apa pun. Ditolak oleh `add`, `login`, `logout`, `use`, `rename`, `remove`, dan `import`, dengan `` --dry-run is not supported by `ccam add`; nothing was changed `` (menyebutkan perintahnya) | nonaktif |

`ccam login`, `ccam logout`, `ccam use`, dan `ccam rename` tidak memiliki opsi sendiri. `ccam completion` menerima `bash`, `zsh`, atau `powershell`; selain itu gagal dengan `only bash, zsh and powershell are supported`.

Nama akun terdiri dari 1 sampai 32 karakter `a-z`, `0-9`, `-`, dan `_`, diawali huruf atau angka. Nama tidak boleh seluruhnya angka, karena `ccam 2` berarti slot 2. Nama juga tidak boleh berupa perintah ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`), sehingga `ccam <name>` tidak pernah ambigu. Daftarnya adalah `config.Reserved` di `internal/config`.

## Cara kerjanya

`ccam add work --shortcut ccwork` menjalankan langkah-langkah ini di `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` memeriksa nama, dan `shim.Check` memeriksa pintasan sebelum apa pun dibuat.
2. `profile.Add` membuat `~/.local/share/ccam/profiles/work/` dengan mode `0700` dan memberi akun nomor slot kosong terkecil (`Config.NextSlot`), sehingga nomor akun yang dihapus dipakai ulang. Akun pertama menjadi akun default.
3. Akun bergabung ke grup-grupnya (`default` kecuali `--group` atau `--no-group` menentukan lain; akun pertama membuat `default`), dan `applyShare` menautkan folder baru yang masih kosong ke grup-grup itu dengan `share.Build` lalu `Apply`, misalnya `settings.json` ke `groups/default/settings.json`. Jika langkah ini gagal, folder baru dihapus, karena folder yang tertinggal akan menggagalkan setiap percobaan ulang dengan `folder already exists`.
4. Akun disimpan ke `~/.config/ccam/config.toml` melalui `fsutil.AtomicWrite`: file sementara, `fsync`, lalu `rename`, sehingga crash tidak pernah meninggalkan konfigurasi yang setengah jadi.
5. `runLogin` menjalankan `claude auth login` dengan `CLAUDE_CONFIG_DIR` diatur ke folder akun, dan `--email` dengan email terakhir akun itu agar halaman login sudah terisi. Jika login gagal, akun tetap ada, dan ccam menampilkan galatnya serta meminta Anda menjalankan `ccam login work` nanti.
6. `profile.ReadIdentity` mengambil `oauthAccount` dari `profiles/work/.claude.json`: email, id akun dan organisasi, nama organisasi, dan paketnya (`claude_max`, `claude_team` …). File itu bisa berukuran beberapa megabyte, sehingga `TopLevelValue` memindai satu kunci itu saja alih-alih mem-parsing seluruh file. ccam menyimpannya di `[profile.login]` dalam `config.toml`.
7. `shim.Install` menautkan `~/.local/bin/ccwork` ke biner ccam, yang mencari `argv[0]` di `config.toml` saat dijalankan.

Login bisa memakan waktu beberapa menit, dan ccam lain mungkin menyimpan konfigurasi pada saat itu. Karena itu, perintah yang menunggu menyimpan melalui `env.update`, yang memuat ulang `config.toml`, menerapkan satu perubahan itu (di sini, email), lalu menyimpan. `ccam remove` memuat ulang setelah pertanyaannya karena alasan yang sama.

`ccam login` menjalankan langkah 5 dan 6 untuk akun yang sudah ada. Claude Code menganggap email itu hanya sebagai petunjuk untuk halaman login, sehingga sebuah login bisa berakhir di akun Claude yang berbeda dari akun ccam sebelumnya; ccam membandingkan login baru dengan `[profile.login]` dan memberitahukannya (`ccam: this login is other@company.com, but work was dev@company.com before.`). ccam mempertahankan login baru dan memberi tahu cara membatalkannya: `ccam logout work`, lalu `ccam login work` dan login dengan email yang benar. Jika akun ccam lain sudah login ke akun Claude yang sama, ccam memperingatkan bahwa keduanya berbagi batas penggunaan. `ccam login --expired` melakukan login, satu per satu, untuk setiap akun yang belum pernah login atau yang login-nya (refresh token-nya) sudah kedaluwarsa.

`ccam logout` menjalankan `claude auth logout` di akun tersebut dan menghapus email yang tersimpan; catatan login tetap ada, sehingga login berikutnya diperiksa terhadapnya. `ccam use` hanya menulis ulang `default` di `config.toml`.

`ccam list` menambahkan kolom `Organization` jika akun-akun berasal dari lebih dari satu organisasi (paket pribadi ditampilkan sebagai `personal`), dan `ccam list --json` memiliki `plan` dan `org` untuk setiap akun.

### Login di mesin jarak jauh

Melalui SSH, bagian browser dari `claude auth login` bisa gagal. Claude Code menunggu login di server callback pada `http://localhost:<port>/callback` di server, tetapi browser yang menyelesaikan login berjalan di komputer Anda, tempat alamat itu tidak mengarah ke mana pun dan browser menampilkan galat seperti `ERR_CONNECTION_REFUSED`. Alamat itu tetap membawa kode sekali pakai, sehingga ccam dapat menyampaikannya:

<!-- motion:accounts-relay -->

`runLogin` tidak menyerahkan terminal Anda ke Claude Code. Fungsi ini menghubungkan input standar Claude Code ke sebuah pipe, dan `relayLogin` menyalin setiap baris yang Anda ketik ke dalamnya, sehingga kode yang ditempel di prompt `Paste code here` milik Claude Code tetap sampai. Jika sebuah baris berupa alamat callback, `callbackURL` hanya menerimanya jika host-nya `localhost`, `127.0.0.1`, atau `[::1]`, port disebutkan, path-nya `/callback`, serta `code` dan `state` keduanya ada; `http://` boleh dihilangkan. `deliverCallback` lalu meminta alamat itu di server dengan batas 10 detik dan tanpa mengikuti pengalihan, seperti yang akan dilakukan browser. Permintaan itu sampai ke server callback milik `claude auth login` yang sedang menunggu (garis putus-putus di diagram). Claude Code menerima kode dan menyimpan login, di Linux dalam `.credentials.json` di folder akun. Baris yang ditempel tidak akan pernah membuat ccam menghubungi host lain.

ccam menampilkan ini sebelum login dimulai:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Jika berhasil, ccam menampilkan `→ Passed the sign-in to Claude Code on this machine.`

Selama login menunggu input di terminal, ccam mematikan echo dan mode baris terminal serta menampilkan `*` untuk setiap karakter yang Anda ketik, sehingga kode atau alamat callback yang ditempel tidak muncul di layar. `readLine` menangani penyuntingan baris sendiri: `Enter` mengirim baris, `Backspace` menghapus satu karakter, dan `Ctrl+D` di baris kosong mengakhiri input. Terminal dipulihkan saat login berakhir, dan juga saat `Ctrl+C`, `SIGTERM`, atau `SIGHUP`, yang tetap mengakhiri ccam. Jika input standar dan error standar tidak terhubung ke terminal, baris disalin apa adanya saat masuk.

### Mengganti nama, menghapus, dan mengimpor

`ccam rename` dan `ccam remove` bekerja pada folder akun dan pada apa yang disimpan ccam atas nama akun itu di tempat lain. `nameState` mendaftar path tersebut: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>`, serta cache penggunaan `~/.local/state/ccam/usage/<name>.json` dan `<name>.statusline.json`. `ccam remove work` menanganinya dengan urutan ini:

<!-- motion:accounts-remove -->

`ccam remove` pertama-tama menolak jika akun sedang berjalan, lalu bertanya. Jika akun membagikan ekstensi, `RewritePluginPaths` menulis ulang `installed_plugins.json` dan `known_marketplaces.json` milik grup, karena Claude Code menyimpan path plugin absolut yang mungkin melewati folder akun ini. `profile.Trash` memindahkan folder ke `~/.local/share/ccam/trash/<name>-<time>/`, akun dikeluarkan dari `config.toml`, path `nameState` dihapus agar akun berikutnya dengan nama yang sama tidak mewarisinya, dan tautan pintasan dihapus. Untuk memulihkannya, pindahkan kembali folder itu dari tempat sampah.

`ccam rename` mengganti nama folder, menulis ulang path plugin di `plugins/` miliknya dari folder lama ke folder baru, ikut memindahkan status default, dan memindahkan path `nameState` ke nama baru. Perintah pintasan tetap memakai namanya. Jika penulisan ulang plugin gagal, nama folder dikembalikan, karena konfigurasi tidak disimpan.

`ccam import` mendaftarkan folder di tempatnya, sebagai `dir` di `config.toml`. Perintah ini tidak pernah memindahkan atau menyalinnya: refresh token yang tersalin akan membuat salinan-salinan itu saling mengeluarkan login satu sama lain. Tanpa argumen folder, perintah ini mencari folder `~/.claude` dan `~/.claude-*` yang belum terdaftar dan menanyakan masing-masing; nama diambil dari folder (`~/.claude` menjadi `main`, `~/.claude-work` menjadi `work`, `~/.claude-2` menjadi `acct2`), dan folder yang namanya tidak dapat diturunkan ditampilkan beserta perintah `--name` yang perlu dipakai. Tanpa `CLAUDE_CONFIG_DIR`, Claude Code menyimpan `.claude.json` di folder home Anda di samping `~/.claude`, sehingga saat Anda mengimpor `~/.claude` dan folder itu tidak memiliki `.claude.json` sendiri, ccam membuat tautan `~/.claude/.claude.json` ke `~/.claude.json`. `claude` biasa dan `ccam main` kemudian melihat akun yang sama. Akun yang diimpor tidak berada di grup mana pun sampai Anda menambahkannya ke salah satu grup (`ccam group join <group> <name>`), dan tidak ada pintasan yang dibuat.

`ccam completion bash|zsh` menampilkan skrip pelengkapan yang juga melengkapi nama akun.

## Keamanan dan kasus khusus

- `ccam rename` dan `ccam remove` menolak saat akun sedang berjalan (`refuseRunning`), dengan `work is running; quit its Claude Code session first`, karena sesi terus menulis ke folder tersebut.
- `ccam remove` pada akun yang diimpor hanya membatalkan pendaftarannya: `✓ main unregistered (folder ... left in place)`. `~/.claude` Anda tidak dipindahkan ke tempat sampah.
- `ccam add` menolak nama yang foldernya sudah ada di `profiles/`, dan mengarahkan ke `ccam import`.
- `ccam import` menolak folder yang sudah terdaftar, path yang bukan folder, dan nama yang sudah dipakai.
- `ccam logout` dan `ccam login` tidak memeriksa apakah akun sedang berjalan.
- Jika penerusan tidak dapat menjangkau Claude Code, misalnya setelah login berakhir, ccam menampilkan `ccam: could not pass the sign-in to Claude Code:` beserta alasannya. Galat HTTP biasanya berarti kode sudah kedaluwarsa; mulai login lagi.

## Terkait

- [Menjalankan Claude Code](/features/launch)
- [Perintah pintasan](/features/shortcuts)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Beberapa akun Claude Code](/guides/multiple-claude-code-accounts)
