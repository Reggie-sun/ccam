---
title: Beralih ke akun berikutnya
description: ccam next membuka Claude Code sebagai akun lain; --auto memilih ruang terbanyak, --continue membawa percakapan, dan --fallback mengulang tugas -p.
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` membuka Claude Code sebagai akun lain tanpa Anda perlu menyebutkan namanya. Gunakan saat sebuah akun mencapai batas penggunaannya: `ccam next` biasa mengambil akun setelah akun yang terakhir Anda pakai, `--auto` mengambil akun dengan ruang terbanyak, dan `--continue` membawa percakapan terakhir di folder ini ke akun baru, baik kedua akun berbagi riwayat maupun tidak. Di dalam Claude Code, [plugin ccam](#inside-claude-code) memberi tahu Anda saat sebuah batas menghentikan sesi, dan `/ccam:switch` meneruskan percakapan. Untuk skrip yang menjalankan `claude -p`, `ccam --fallback` melakukan pergantian itu dengan sendirinya (lihat [Jalankan ulang saat batas tercapai](#run-again-at-the-limit)).

## Penggunaan

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam menampilkan perpindahan sebagai `→ <from> → <to>` di error standar, beserta alasan pilihan `--auto`, lalu menjalankan Claude Code sebagai akun baru dengan cara yang sama seperti `ccam <name>`. `<from>` adalah akun yang terakhir Anda pakai, atau dengan `--continue`, akun pemilik percakapan itu. Argumen setelah `--` diteruskan ke `claude` tanpa diubah.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--auto` | Memilih akun yang sudah login berdasarkan penggunaan (lihat di bawah), melewati akun yang sudah 100% | nonaktif |
| `--max-usage <1-100>` | Dengan `--auto`: utamakan akun di bawah penggunaan ini | `auto_max_usage`, jika tidak ada 100 |
| `--strategy headroom\|expiring\|balance` | Dengan `--auto`: penggunaan 5 jam terendah, reset 5 jam yang paling cepat, atau kuota mingguan yang paling dulu akan tidak terpakai | `auto_strategy`, jika tidak ada `headroom` |
| `--cross-org` | Dengan `--auto`: juga memilih akun dari organisasi lain (lihat di bawah) | `cross_org`, jika tidak ada nonaktif |
| `--wait[=<duration>]` | Dengan `--auto`: jika tidak ada akun yang punya ruang, tunggu akun pertama di-reset, paling lama selama ini | nonaktif; `5h` jika diberikan tanpa nilai |
| `-c`, `--continue` | Membawa percakapan terakhir di folder ini ke akun baru (lihat di bawah) | nonaktif |
| `-- <claude args>` | Argumen yang diteruskan ke `claude` setelah argumen milik ccam | tidak ada |
| `--dry-run` | Dihormati: menampilkan `→ <from> → <to>` dan baris perintah yang akan dijalankan, tanpa menjalankan Claude Code atau mencatat peluncuran | nonaktif |

`--auto` dan `--continue` dapat dipakai bersama. Dengan `--continue`, argumen `--resume` ditempatkan sebelum argumen Anda sendiri.

Opsi milik Claude Code sendiri, seperti `--dangerously-skip-permissions` atau `--model`, ditaruh setelah `--`. Tidak seperti `ccam <name>`, yang meneruskan semua yang ada setelah nama, `ccam next` membaca argumennya sendiri, jadi tanpa `--` ccam berhenti dengan `ccam: unknown flag: --dangerously-skip-permissions`.

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

`--auto` memutuskan berdasarkan penggunaan alih-alih urutan slot. `cli.pickAuto` mempertahankan akun yang sudah login, menanyakan penggunaan semuanya sekaligus dengan `Fetcher.FetchAll` yang sama dengan yang dipakai `ccam usage`, dan membiarkan `usage.Policy.Choose()` memilih:

<!-- motion:next-auto -->

Akun dinilai berdasarkan jendela 5 jam dan mingguannya, dan, jika `--model` ada di antara argumen untuk Claude Code (setelah `--`), juga berdasarkan jendela mingguan model itu (`opus`, `sonnet`, `haiku`, atau `fable` di namanya). `Policy.Choose()` membagi akun menjadi tiga kelompok:

1. **Ada ruang**: setiap jendela di bawah `auto_max_usage` (`config.toml`, bawaan 100). Strategi memilih di antaranya: `headroom` (bawaan) mengambil penggunaan 5 jam terendah; `expiring` mengambil jendela 5 jam yang paling cepat di-reset, sehingga sisa kuotanya terpakai sebelum hangus; `balance` mengambil akun dengan sisa kuota mingguan terbanyak per hari hingga minggunya di-reset, dengan menghitung setiap sesi ccam yang sedang berjalan sebagai beban akun itu, sehingga kuota mingguan tidak dibiarkan tak terpakai dan sesi paralel tersebar.
2. **Lewat batas**: di bawah 100% tetapi sama dengan atau di atas `auto_max_usage`. Hanya dipakai jika tidak ada akun yang punya ruang, dimulai dari penggunaan 5 jam terendah.
3. **Penuh**: sebuah jendela di 100%, atau penggunaan 5 jam yang tidak diketahui. Tidak pernah dipilih.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organisasi.** Pilihan otomatis tetap berada di akun yang serupa dengan akun yang Anda tinggalkan: organisasi tim atau enterprise yang sama, atau paket pribadi (Pro, Max) bersama-sama, sehingga pekerjaan dari akun perusahaan tidak berlanjut di akun pribadi, atau sebaliknya. ccam mengetahui organisasi setiap akun dari login-nya (`ccam login` mencatatnya; login lama dibaca dari `.claude.json`); akun yang organisasinya tidak bisa dipastikan cocok di mana saja. Dua folder yang login ke akun Claude yang sama berbagi batasnya, sehingga yang satu tidak pernah menggantikan yang lain. `--cross-org`, atau `cross_org = true` di `config.toml`, mencabut aturan organisasi ini.

Jika tidak ada yang bisa dipilih, ccam keluar dengan kode 75 (`EX_TEMPFAIL`), sehingga skrip bisa membedakan "coba lagi nanti" dari kegagalan. Dengan `--wait`, ccam justru tidur sampai jendela penghalang milik akun penuh pertama di-reset (30 detik setelah waktu reset, saat endpoint menampilkan jendela baru), mencetak `ccam: no account has room; waiting until 15:10`, lalu memilih lagi; jika tidak ada akun yang menyebutkan kapan di-reset, ccam memeriksa lagi setelah 10 menit. Jika tidak ada yang kosong sebelum waktu tunggu berakhir, ccam keluar dengan 75. Pencarian mengikuti aturan `ccam usage`: angka yang berumur kurang dari 60 detik dipakai ulang, setiap akun diberi waktu 3 detik, dan akun yang tidak menjawab, atau yang login-nya sudah kedaluwarsa, dinilai berdasarkan angka yang tersimpan. Jendela tersimpan yang waktu reset-nya sudah lewat dihitung 0%, sehingga akun yang jendelanya sudah direset sejak itu tidak dilewatkan. `--auto` tidak pernah memperbarui token. Tidak seperti rotasi biasa, `--auto` dapat memilih akun yang terakhir Anda pakai jika akun itu masih memiliki ruang terbanyak.

### --continue

1. ccam mencari percakapan terakhir di folder ini di setiap akun: `<session>.jsonl` terbaru di `projects/<folder>/` (nama foldernya adalah folder kerja dengan setiap karakter selain huruf atau angka diubah menjadi `-`, seperti cara Claude Code menamainya). Jika sesi itu meninggalkan catatan (batas penggunaan atau `/ccam:switch`, lihat di bawah) dan belum ada yang lebih baru ditulis sejak itu, sesi itulah yang dipakai.
2. Akun pemilik percakapan itu adalah akun yang Anda tinggalkan, dan akun baru dipilih seperti di atas, di dalam organisasi yang sama.
3. Jika akun baru berbagi riwayat itu (folder grup yang sama), ccam menjalankan `claude --resume <session id>`, tepat percakapan itu (`--continue` akan membuka percakapan mana pun yang terbaru di folder tersebut).
4. Jika tidak, ccam menjalankan `claude --resume <the conversation's file> --fork-session`: Claude Code menyalin seluruh percakapan, dengan id sesi baru, ke riwayat milik akun baru, dan aslinya tetap seperti semula. Sejak itu `--continue` di akun baru akan menemukannya.

ccam menyebutkan seberapa besar percakapannya: akun baru membaca seluruhnya lagi pada balasan pertamanya. Cache prompt tidak pernah dibagikan antarorganisasi, sehingga dengan `--cross-org` balasan pertama itu memakan ukuran penuh. Berbagi `history` dalam sebuah grup menghindari penyalinan, mempertahankan satu percakapan, dan juga membagikan checkpoint tempat `/rewind` memulihkan kode (`file-history/`).

## Di dalam Claude Code

[Plugin ccam](/docs/install#claude-code-plugin) menambahkan dua hook (Claude Code 2.1.78 atau yang lebih baru):

- Saat sebuah giliran berhenti di batas penggunaan (`You've hit your session limit`, batas mingguan, atau batas model), hook mencatat percakapan itu dan menampilkan notifikasi desktop di terminal yang mendukungnya (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Pembatasan singkat dari server atau batas laju (rate limit) kunci API bukanlah batas penggunaan dan dibiarkan. Claude Code 2.1.234 dan yang lebih baru juga menunggu di dalam sesi dan melanjutkan dengan sendirinya setelah reset 5 jam; notifikasi itu memberi Anda pilihan untuk berpindah akun saja.
- `/ccam:switch [account]` mencatat percakapan dan mengakhiri giliran sebelum permintaan apa pun dibuat: ketik `/exit`, lalu `ccam next --continue` melanjutkannya, sebagai akun yang Anda sebutkan atau sebagai akun berikutnya.

Catatan itu disimpan di `~/.local/state/ccam/handoff/<account>.json` selama 12 jam dan dipakai sekali, saat sesi berikutnya akan dimulai; jika peluncurannya gagal, catatan itu tetap ada untuk percobaan berikutnya. Catatan itu milik folder tempat sesi dimulai, bahkan setelah `cd` di shell Claude Code. Jika akun yang disebut di dalamnya sejak itu diganti namanya atau dihapus, ccam memberitahukannya dan memilih seperti biasa. Dengan `relay = true` di `config.toml`, ccam tetap berjalan di samping sesi interaktif alih-alih menggantikan dirinya dengan Claude Code, dan setelah `/ccam:switch` ccam melanjutkan percakapan sesi itu dengan sendirinya (catatan yang ditinggalkan sesi lain tidak ditindaklanjuti). Itu memerlukan satu proses ccam yang menunggu per sesi, karena itulah fitur ini nonaktif secara default. Sesi lanjutan itu dimulai dengan `ccam next --continue`, tanpa argumen yang dipakai saat sesi pertama dimulai, seperti `--model`.

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
8. Hanya akun dari organisasi yang sama yang dicoba, seperti pada `--auto` (`cross_org = true` mencabut aturan ini).

Akun berikutnya dipilih dengan aturan yang sama seperti `--auto` (`auto_max_usage`, `auto_strategy`, jendela `--model`). Jika tidak ada akun yang tersisa, ccam mencetak output percobaan terakhir, menyatakan bahwa tidak ada akun yang sudah login yang masih punya ruang, dan keluar dengan kode 75. Dengan `--wait` sebelum nama akun (`ccam --fallback --wait=2h work -p …`), ccam justru menunggu akun pertama di-reset, lalu dapat mencoba setiap akun lagi; ccam selalu tidur setidaknya sekali terlebih dahulu, sehingga akun yang baru saja mencapai batasnya tidak langsung dijalankan lagi. `ccam --dry-run --fallback work -p …` mencetak perintah percobaan pertama.

Sesi interaktif tidak berganti akun dengan sendirinya. Sebagai gantinya, dengan `ccam statusline install`, baris status menambahkan `→ ccam next --auto --continue` begitu penggunaan 5 jam atau mingguan mencapai 90% ([Batas penggunaan](/features/usage)).

## Keamanan dan kasus khusus

- Akun tanpa login tidak pernah dipilih. Jika tidak ada akun lain yang sudah login, `ccam next` biasa berhenti dengan `no logged-in account other than <name>`.
- Dengan `--continue` dan tidak ada percakapan di folder ini: `no account has a conversation in <folder> to continue`.
- Jika aturan organisasi tidak menyisakan akun apa pun: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Dengan `--auto` dan tidak ada akun yang sudah login: `no logged-in account to pick`. Jika setiap akun penuh atau penggunaannya tidak diketahui: `` no account with known usage and headroom (see `ccam usage`) ``. Keduanya berakhir dengan kode keluar 75.
- `--max-usage`, `--strategy`, `--wait`, dan `--cross-org` hanya berfungsi dengan `--auto`, dan nilai di luar 1-100 atau strategi lain ditolak.
- Fork menyalin percakapan, bukan status akun lainnya: checkpoint di `file-history/` milik akun asal tidak ikut terbawa kecuali riwayat dibagikan.
- `--auto` membuat permintaan jaringan bahkan dengan `--dry-run`. Opsi ini hanya membaca, dan tidak pernah menulis file token.
- Akun yang dipilih dijalankan melalui jalur yang sama dengan `ccam <name>`, sehingga pemeriksaan yang sama berlaku, termasuk peringatan jika Claude Code lain yang berbagi riwayat yang sama sudah berjalan di folder ini.
- `last-used` hanya ditulis saat Claude Code benar-benar dimulai. `--dry-run` tidak menyentuhnya, sehingga dry run tidak memajukan rotasi.
- `--fallback` menahan output standar sampai sebuah percobaan selesai, sehingga output streaming seperti `--output-format stream-json` tiba sekaligus di akhir. File yang diubah percobaan pertama sebelum mencapai batas tetap berubah; percobaan berikutnya memulai prompt dari awal.

## Terkait

- [Batas penggunaan](/features/usage)
- [Batas penggunaan dan perpindahan akun](/docs/usage-limits)
- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Alur kerja batas penggunaan](/guides/usage-limit-workflow)
