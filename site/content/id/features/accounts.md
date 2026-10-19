---
title: Mengelola akun
description: Tambahkan, login, tampilkan, ganti nama, hapus, dan impor akun Claude Code dengan ccam add, ccam login, ccam list, ccam rename, ccam remove, dan ccam import.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 97cef50f48f4
---

Setiap akun ccam adalah folder konfigurasi Claude Code tersendiri dengan login-nya sendiri. Perintah-perintah ini membuat, memasukkan login, menampilkan, mengganti nama, dan menghapus akun, serta mendaftarkan folder yang sudah Anda miliki.

## Penggunaan

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
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

`ccam add` di terminal tanpa opsi akan membuka wizard. Wizard menanyakan nama, lalu menampilkan daftar centang item yang dibagikan (semua item kecuali riwayat percakapan sudah dicentang), lalu bertanya `Create a shortcut command for work? [y/N]` dan, jika Anda menjawab `y`, `Shortcut name [ccwork]:`. Salah satu dari `--share`, `--shortcut`, atau `--no-login`, atau pipe sebagai ganti terminal, akan melewati wizard; dalam hal itu nama wajib diberikan dan tidak ada pintasan yang dibuat kecuali `--shortcut` menyebutkannya.

`ccam list` menampilkan slot, nama, pintasan, email, item yang dibagikan, dan waktu terakhir digunakan, dengan `*` di depan akun default:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

Dengan `--json`, setiap akun adalah objek dengan `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (dihilangkan jika akun belum pernah dijalankan), dan `share`. `share` memetakan setiap item (`mcp`, `settings`, `memory`, `extensions`, `history`) ke nama grupnya atau ke `"local"`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `ccam add --no-login` | Membuat akun tanpa login; login belakangan dengan `ccam login <name>` | nonaktif |
| `ccam add --share <item=group>` | Membagikan item ke sebuah grup, atau `item=local` agar tetap privat. Bisa diulang | `mcp`, `settings`, `memory`, `extensions` di `default`; `history` lokal |
| `ccam add --shortcut <command>` | Juga membuat perintah pintasan ini di `~/.local/bin` | tidak ada |
| `ccam list --json` | Menampilkan akun sebagai JSON | nonaktif |
| `ccam remove -y`, `--yes` | Menghapus tanpa bertanya | nonaktif |
| `ccam import --name <name>` | Nama akun untuk folder yang Anda berikan | diambil dari nama folder |
| `ccam import -y`, `--yes` | Mengimpor setiap folder yang ditemukan tanpa bertanya | nonaktif |
| `--dry-run` | Dihormati oleh `ccam list` dan `ccam completion`, yang tidak mengubah apa pun. Ditolak oleh `add`, `login`, `logout`, `use`, `rename`, `remove`, dan `import`, dengan `` --dry-run is not supported by `ccam add`; nothing was changed `` (menyebutkan perintahnya) | nonaktif |

`ccam login`, `ccam logout`, `ccam use`, dan `ccam rename` tidak memiliki opsi sendiri. `ccam completion` menerima `bash` atau `zsh`; selain itu gagal dengan `only bash and zsh are supported`.

Nama akun terdiri dari 1 sampai 32 karakter `a-z`, `0-9`, `-`, dan `_`, diawali huruf atau angka. Nama tidak boleh seluruhnya angka, karena `ccam 2` berarti slot 2. Nama juga tidak boleh berupa perintah ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), sehingga `ccam <name>` tidak pernah ambigu. Daftarnya adalah `config.Reserved` di `internal/config`.

## Cara kerjanya

`ccam add work --shortcut ccwork` menjalankan langkah-langkah ini di `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` memeriksa nama, dan `shim.Check` memeriksa pintasan sebelum apa pun dibuat.
2. `profile.Add` membuat `~/.local/share/ccam/profiles/work/` dengan mode `0700` dan memberi akun nomor slot kosong terkecil (`Config.NextSlot`), sehingga nomor akun yang dihapus dipakai ulang. Akun pertama menjadi akun default.
3. `applyShare` menautkan folder baru yang masih kosong ke grup pembagiannya dengan `share.Build` lalu `Apply`, misalnya `settings.json` ke `groups/default/settings.json`. Jika langkah ini gagal, folder baru dihapus, karena folder yang tertinggal akan menggagalkan setiap percobaan ulang dengan `folder already exists`.
4. Akun disimpan ke `~/.config/ccam/config.toml` melalui `fsutil.AtomicWrite`: file sementara, `fsync`, lalu `rename`, sehingga crash tidak pernah meninggalkan konfigurasi yang setengah jadi.
5. `runLogin` menjalankan `claude auth login` dengan `CLAUDE_CONFIG_DIR` diatur ke folder akun. Jika login gagal, akun tetap ada, dan ccam menampilkan galatnya serta meminta Anda menjalankan `ccam login work` nanti.
6. `profile.ReadEmail` mengambil `oauthAccount.emailAddress` dari `profiles/work/.claude.json`. File itu bisa berukuran beberapa megabyte, sehingga `TopLevelValue` memindai satu kunci itu saja alih-alih mem-parsing seluruh file.
7. `shim.Install` menautkan `~/.local/bin/ccwork` ke biner ccam, yang mencari `argv[0]` di `config.toml` saat dijalankan.

Login bisa memakan waktu beberapa menit, dan ccam lain mungkin menyimpan konfigurasi pada saat itu. Karena itu, perintah yang menunggu menyimpan melalui `env.update`, yang memuat ulang `config.toml`, menerapkan satu perubahan itu (di sini, email), lalu menyimpan. `ccam remove` memuat ulang setelah pertanyaannya karena alasan yang sama.

`ccam login` menjalankan langkah 5 dan 6 untuk akun yang sudah ada. `ccam logout` menjalankan `claude auth logout` di akun tersebut dan menghapus email yang tersimpan. `ccam use` hanya menulis ulang `default` di `config.toml`.

### Login di mesin jarak jauh

Melalui SSH, bagian browser dari `claude auth login` bisa gagal. Claude Code menunggu login di server callback pada `http://localhost:<port>/callback` di server, tetapi browser yang menyelesaikan login berjalan di komputer Anda, tempat alamat itu tidak mengarah ke mana pun dan browser menampilkan galat seperti `ERR_CONNECTION_REFUSED`. Alamat itu tetap membawa kode sekali pakai, sehingga ccam dapat menyampaikannya:

<!-- motion:accounts-relay -->

`runLogin` tidak menyerahkan terminal Anda ke Claude Code. Fungsi ini menghubungkan input standar Claude Code ke sebuah pipe, dan `relayLogin` menyalin setiap baris yang Anda ketik ke dalamnya, sehingga kode yang ditempel di prompt `Paste code here` milik Claude Code tetap sampai. Jika sebuah baris berupa alamat callback, `callbackURL` hanya menerimanya jika host-nya `localhost`, `127.0.0.1`, atau `[::1]`, port disebutkan, path-nya `/callback`, serta `code` dan `state` keduanya ada; `http://` boleh dihilangkan. `deliverCallback` lalu meminta alamat itu di server dengan batas 10 detik dan tanpa mengikuti pengalihan, seperti yang akan dilakukan browser. Permintaan itu sampai ke server callback milik `claude auth login` yang sedang menunggu (garis putus-putus di diagram). Claude Code menerima kode dan menyimpan login, di Linux dalam `.credentials.json` di folder akun. Baris yang ditempel tidak akan pernah membuat ccam menghubungi host lain.

ccam menampilkan ini sebelum login dimulai:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Jika berhasil, ccam menampilkan `→ Passed the sign-in to Claude Code on this machine.`

### Mengganti nama, menghapus, dan mengimpor

`ccam rename` dan `ccam remove` bekerja pada folder akun dan pada apa yang disimpan ccam atas nama akun itu di tempat lain. `nameState` mendaftar path tersebut: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>`, serta cache penggunaan `~/.local/state/ccam/usage/<name>.json` dan `<name>.statusline.json`. `ccam remove work` menanganinya dengan urutan ini:

<!-- motion:accounts-remove -->

`ccam remove` pertama-tama menolak jika akun sedang berjalan, lalu bertanya. Jika akun membagikan ekstensi, `RewritePluginPaths` menulis ulang `installed_plugins.json` dan `known_marketplaces.json` milik grup, karena Claude Code menyimpan path plugin absolut yang mungkin melewati folder akun ini. `profile.Trash` memindahkan folder ke `~/.local/share/ccam/trash/<name>-<time>/`, akun dikeluarkan dari `config.toml`, path `nameState` dihapus agar akun berikutnya dengan nama yang sama tidak mewarisinya, dan tautan pintasan dihapus. Untuk memulihkannya, pindahkan kembali folder itu dari tempat sampah.

`ccam rename` mengganti nama folder, menulis ulang path plugin di `plugins/` miliknya dari folder lama ke folder baru, ikut memindahkan status default, dan memindahkan path `nameState` ke nama baru. Perintah pintasan tetap memakai namanya. Jika penulisan ulang plugin gagal, nama folder dikembalikan, karena konfigurasi tidak disimpan.

`ccam import` mendaftarkan folder di tempatnya, sebagai `dir` di `config.toml`. Perintah ini tidak pernah memindahkan atau menyalinnya: refresh token yang tersalin akan membuat salinan-salinan itu saling mengeluarkan login satu sama lain. Tanpa argumen folder, perintah ini mencari folder `~/.claude` dan `~/.claude-*` yang belum terdaftar dan menanyakan masing-masing; nama diambil dari folder (`~/.claude` menjadi `main`, `~/.claude-work` menjadi `work`, `~/.claude-2` menjadi `acct2`), dan folder yang namanya tidak dapat diturunkan ditampilkan beserta perintah `--name` yang perlu dipakai. Tanpa `CLAUDE_CONFIG_DIR`, Claude Code menyimpan `.claude.json` di folder home Anda di samping `~/.claude`, sehingga saat Anda mengimpor `~/.claude` dan folder itu tidak memiliki `.claude.json` sendiri, ccam membuat tautan `~/.claude/.claude.json` ke `~/.claude.json`. `claude` biasa dan `ccam main` kemudian melihat akun yang sama. Akun yang diimpor tidak membagikan apa pun sampai Anda menjalankan `ccam share`, dan tidak ada pintasan yang dibuat.

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
