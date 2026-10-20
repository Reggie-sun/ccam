---
title: Beralih ke akun berikutnya
description: ccam next membuka Claude Code sebagai akun setelah akun terakhir; --auto memilih ruang terbanyak, --continue melanjutkan, --fallback mengulang tugas -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 7bc55c2dd520
---

`ccam next` membuka Claude Code sebagai akun lain tanpa Anda perlu menyebutkan namanya. Gunakan saat sebuah akun mencapai batas penggunaannya: `ccam next` biasa mengambil akun setelah akun yang terakhir Anda pakai, `--auto` mengambil akun dengan ruang terbanyak di jendela 5 jamnya, dan `--continue` melanjutkan percakapan terakhir di akun baru. Untuk skrip yang menjalankan `claude -p`, `ccam --fallback` melakukan pergantian itu dengan sendirinya (lihat [Jalankan ulang saat batas tercapai](#run-again-at-the-limit)).

## Penggunaan

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam menampilkan perpindahan sebagai `→ <from> → <to>` di error standar, lalu menjalankan Claude Code sebagai akun baru dengan cara yang sama seperti `ccam <name>`. `<from>` adalah akun yang terakhir Anda pakai. Argumen setelah `--` diteruskan ke `claude` tanpa diubah.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--auto` | Memilih akun yang sudah login dengan penggunaan 5 jam terendah, melewati akun yang sudah 100% | nonaktif |
| `-c`, `--continue` | Hanya mempertimbangkan akun yang berbagi riwayat percakapan dengan akun yang terakhir Anda pakai, dan menjalankan `claude --continue` | nonaktif |
| `-- <claude args>` | Argumen yang diteruskan ke `claude` setelah argumen milik ccam | tidak ada |
| `--dry-run` | Dihormati: menampilkan `→ <from> → <to>` dan baris perintah yang akan dijalankan, tanpa menjalankan Claude Code atau mencatat peluncuran | nonaktif |

`--auto` dan `--continue` dapat dipakai bersama. Dengan `--continue`, `--continue` ditempatkan sebelum argumen Anda sendiri.

## Cara kerjanya

Setiap peluncuran melalui ccam menulis waktunya ke `~/.local/state/ccam/last-used/<name>`. `ccam next` membaca file-file ini, memeriksa akun mana yang sudah login, lalu memilih.

1. `cli.candidates` membuat satu `next.Candidate` per akun dari `config.toml`: akun, waktu peluncuran terakhirnya (`launcher.LastUsed`), dan apakah akun itu memiliki login (`creds.Has`: `.credentials.json` yang dapat dibaca, atau di macOS sebuah item Rantai Kunci (Keychain)).
2. Tanpa `--auto`, `next.Rotate` mengurutkan akun berdasarkan nomor slot dan mencari akun yang paling baru diluncurkan. Itulah akun yang Anda tinggalkan. Jika belum ada akun yang pernah diluncurkan, itu adalah slot terkecil.
3. Dari sana fungsi ini bergerak maju sesuai urutan slot, kembali dari slot terakhir ke slot pertama, dan mengembalikan akun pertama yang sudah login.
4. `launcher.LaunchWith` menyiapkan dan menjalankan akun yang dipilih persis seperti `ccam <name>`: memeriksa tautan akun, menerapkan pembagian MCP, mendaftarkan sesi, dan mengganti ccam dengan `claude` melalui `syscall.Exec`.

Penelusuran dimulai dari akun yang terakhir dipakai dan berhenti di akun pertama yang sudah login setelahnya:

<!-- motion:next-rotate -->

`work` paling baru diluncurkan, sehingga penelusuran dimulai setelahnya. `side` tidak memiliki login dan dilewati, lalu `lab` dipilih. Seandainya `lab` juga belum login, penelusuran akan kembali ke `personal` (`% len` di diagram). Rotasi tidak pernah mengembalikan akun asal: jika tidak ada akun lain yang sudah login, rotasi berhenti dengan galat.

### --auto

`--auto` memutuskan berdasarkan penggunaan alih-alih urutan slot. `cli.pickAuto` mempertahankan akun yang sudah login, menanyakan penggunaan semuanya sekaligus dengan `Fetcher.FetchAll` yang sama dengan yang dipakai `ccam usage`, dan membiarkan `usage.Best()` memilih:

<!-- motion:next-auto -->

`usage.Best()` mengambil persentase 5 jam terendah dan melewati akun yang sudah 100% atau yang penggunaan 5 jamnya tidak diketahui. Pencarian mengikuti aturan `ccam usage`: angka yang berumur kurang dari 60 detik dipakai ulang, setiap akun diberi waktu 3 detik, dan akun yang tidak menjawab, atau yang login-nya sudah kedaluwarsa, dinilai berdasarkan angka yang tersimpan. Jendela tersimpan yang waktu reset-nya sudah lewat dihitung 0%, sehingga akun yang jendelanya sudah direset sejak itu tidak dilewatkan. `--auto` tidak pernah memperbarui token. Tidak seperti rotasi biasa, `--auto` dapat memilih akun yang terakhir Anda pakai jika akun itu masih memiliki ruang terbanyak.

### --continue

`--continue` milik Claude Code membuka percakapan terbaru di folder saat ini, yang hanya berhasil jika akun baru dapat melihat percakapan itu. Jadi dengan `--continue`, ccam hanya mempertimbangkan akun dalam grup `history` yang sama dengan akun yang terakhir Anda pakai, lalu menempatkan `--continue` di depan argumen Anda. Akun yang riwayatnya `local` tidak berbagi dengan siapa pun, sehingga `--continue` dari akun seperti itu selalu gagal. Aktifkan pembagian terlebih dahulu dengan `ccam share <name> --history on` (lihat [pembagian](/docs/sharing)).

## Jalankan ulang saat batas tercapai

`ccam next` ditujukan untuk sesi interaktif yang sudah berhenti. Proses non-interaktif, yaitu `claude -p` di skrip atau tugas terjadwal, tidak bisa menunggu Anda, jadi `--fallback` sebelum nama akun membuat ccam berganti akun untuknya:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` hanya berfungsi bersama `-p` atau `--print`; jika tidak, ccam berhenti dan mengarahkan ke `ccam next --auto --continue`.
2. Input standar dibaca sekali dan diberikan ke setiap percobaan.
3. ccam menjalankan `claude` sebagai proses anak alih-alih `exec`, menahan output standarnya, dan meneruskan error standar begitu muncul.
4. Saat `claude` keluar dengan 0, atau gagal karena alasan lain, output yang ditahan dan kode keluarnya diteruskan tanpa diubah.
5. Kegagalan dianggap sebagai batas penggunaan jika pesan Claude Code menyatakannya (`You've hit your limit`, `usage limit reached`), atau jika pengecekan penggunaan akun menunjukkan jendela 5 jam atau mingguannya di 100%.
6. Kemudian output percobaan itu dibuang dan perintah yang sama dijalankan sebagai akun lain yang sudah login: akun dengan penggunaan 5 jam terendah, lalu akun yang penggunaannya tidak bisa dicek, sesuai urutan nomor. Akun di 100% dilewati, dan setiap akun dicoba paling banyak sekali.
7. Jika ada `-c`, `--continue`, `-r`, atau `--resume` di antara argumen, hanya akun dalam grup riwayat yang sama yang dicoba, sehingga percakapannya bisa ditemukan.

Jika tidak ada akun yang tersisa, ccam mencetak output percobaan terakhir, menyatakan bahwa tidak ada akun yang sudah login yang masih punya ruang, dan keluar dengan kode percobaan itu. `ccam --dry-run --fallback work -p …` mencetak perintah percobaan pertama.

Sesi interaktif tidak berganti akun dengan sendirinya. Sebagai gantinya, dengan `ccam statusline install`, baris status menambahkan `→ ccam next --auto --continue` begitu penggunaan 5 jam atau mingguan mencapai 90% ([Batas penggunaan](/features/usage)).

## Keamanan dan kasus khusus

- Akun tanpa login tidak pernah dipilih. Jika tidak ada akun lain yang sudah login, `ccam next` biasa berhenti dengan `no logged-in account other than <name>`.
- Dengan `--continue` dan tidak ada akun yang berbagi riwayat: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Dengan `--auto` dan tidak ada akun yang sudah login: `no logged-in account to pick`. Jika setiap akun penuh atau penggunaannya tidak diketahui: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` membuat permintaan jaringan bahkan dengan `--dry-run`. Opsi ini hanya membaca, dan tidak pernah menulis file token.
- Akun yang dipilih dijalankan melalui jalur yang sama dengan `ccam <name>`, sehingga pemeriksaan yang sama berlaku, termasuk peringatan jika Claude Code lain yang berbagi riwayat yang sama sudah berjalan di folder ini.
- `last-used` hanya ditulis saat Claude Code benar-benar dimulai. `--dry-run` tidak menyentuhnya, sehingga dry run tidak memajukan rotasi.
- `--fallback` menahan output standar sampai sebuah percobaan selesai, sehingga output streaming seperti `--output-format stream-json` tiba sekaligus di akhir. File yang diubah percobaan pertama sebelum mencapai batas tetap berubah; percobaan berikutnya memulai prompt dari awal.

## Terkait

- [Batas penggunaan](/features/usage)
- [Batas penggunaan dan perpindahan akun](/docs/usage-limits)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Alur kerja batas penggunaan](/guides/usage-limit-workflow)
