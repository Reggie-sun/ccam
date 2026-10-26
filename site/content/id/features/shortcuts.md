---
title: Perintah pintasan
description: Berikan akun perintahnya sendiri seperti ccwork dengan ccam shortcut, dan buat ulang tautan di ~/.local/bin dengan ccam shims install.
updated: 2026-10-26
order: 3
translation: machine
sourceHash: "066e7e622517"
---

Pintasan adalah perintah yang namanya Anda pilih sendiri dan menjalankan Claude Code sebagai satu akun: `ccwork --resume` melakukan hal yang sama dengan `ccam work --resume`. Pintasan adalah tautan ke ccam di `~/.local/bin`, sehingga berfungsi di bash, zsh, skrip, dan apa pun yang menjalankan program dari `PATH`. Pintasan bersifat opsional; buatlah untuk akun yang sering Anda jalankan.

## Penggunaan

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

ccam mengonfirmasi setiap perubahan:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Sebuah akun memiliki paling banyak satu pintasan, disimpan sebagai `shortcut` di `~/.config/ccam/config.toml`. Wizard `ccam add` bertanya `Create a shortcut command for <name>? [y/N]` lalu `Shortcut name [cc<name>]:`. Di luar wizard, `ccam add` hanya membuat pintasan dengan `--shortcut`. `ccam list` menampilkan pintasan setiap akun, atau `-`.

Pintasan meneruskan setiap argumen ke Claude Code, sama seperti `ccam <name>` setelah nama. Itu termasuk `--dry-run` dan `--help`: untuk meninjau peluncuran, gunakan `ccam --dry-run work`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--remove` | Menghapus pintasan akun beserta tautannya | nonaktif |
| `-h`, `--help` | Bantuan untuk perintah | |
| `--dry-run` | `ccam shortcut` mendukungnya: memeriksa nama dan menampilkan `Would set the shortcut of work to ccwork.` atau `Would remove ccwork, the shortcut of work.` tanpa menyimpan atau menautkan | nonaktif |

`ccam shims install` tidak menerima argumen dan menolak `--dry-run` dengan `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Aturan nama, diperiksa dengan urutan ini:

- 1 sampai 64 karakter `a-z`, `0-9`, `.`, `-`, dan `_`, diawali huruf atau angka;
- bukan `ccam`, `claude`, atau perintah ccam seperti `run` atau `usage`;
- bukan pintasan akun lain;
- bukan file yang sudah ada di `~/.local/bin`, kecuali file itu tautan yang dibuat ccam (seperti pintasan milik akun ini sendiri);
- bukan nama executable lain di `PATH` Anda.

`config.toml` diperiksa dengan aturan yang sama saat dimuat. Jika Anda mengubah pintasan secara manual sehingga melanggar pola, berupa `claude` atau `ccam`, atau dipakai oleh dua akun, setiap perintah gagal dengan pesan yang menyebutkan akunnya sampai Anda memperbaikinya.

Memberikan perintah dan `--remove` sekaligus akan berhenti dengan `give a command or --remove, not both`. Menetapkan nama yang sudah dimiliki akun menampilkan `work: ccwork (unchanged)` dan membuat ulang tautannya jika hilang.

## Cara kerjanya

Menetapkan nama menjalankan `checkShortcut`, yang memanggil `shim.Check` dengan folder bin dan path asli biner ccam:

<!-- motion:shortcuts-check -->

Nama melewati pemeriksaan secara berurutan. Di samping setiap pemeriksaan ada nama yang ditolaknya: `Work` tidak cocok dengan pola, `run` adalah perintah ccam, `ccside` milik akun lain, `ccfile` adalah file yang sudah ada di `~/.local/bin`, dan `git` akan menutupi `/usr/bin/git`. Hanya nama yang lolos semua pemeriksaan yang disimpan dan ditautkan.

1. ccam mencocokkan nama dengan `^[a-z0-9][a-z0-9._-]{0,63}$`, lalu `config.Reserved`, lalu `cfg.FindShortcut` (`shim.Check`). ccam memanggil `os.Lstat` pada `~/.local/bin/<name>`, dan untuk setiap folder lain di `PATH` memanggil `os.Stat` serta menolak file executable yang bukan tautan ke ccam.
2. ccam menyimpan `config.toml` dengan `fsutil.AtomicWrite` (file sementara, `fsync`, rename).
3. ccam menghapus tautan pintasan lama, tetapi hanya jika itu tautan ke ccam (`shim.Remove`).
4. ccam menelusuri setiap pintasan terkonfigurasi menurut urutan nama dan membuat tautan setiap pintasan yang hilang (`shim.Install`), dengan `os.Symlink` ke biner ccam, yang path-nya ditentukan dengan `os.Executable` dan `filepath.EvalSymlinks`. File yang menghalangi satu pintasan tidak menghentikan pintasan lainnya; galatnya menyebutkan masing-masing.

ccam mencatat setiap tautan yang dibuatnya, beserta biner yang dijalankan setiap tautan, di `~/.local/state/ccam/shortcuts.json`. Sebuah tautan dianggap milik ccam jika `os.Readlink` mengembalikan path ccam, atau jika tautan dan path ccam mengarah ke file yang sama. Tautan tercatat yang binernya sejak itu dipindahkan atau diganti (biner ccam yang dipindahkan ke folder lain, instalasi ulang) sudah usang, dan ccam mengarahkannya ke biner ccam ini (`shim.Install`). Tautan tercatat ke salinan ccam lain yang masih terpasang dibiarkan apa adanya. Setiap penghapusan dan setiap pemeriksaan penimpaan melewati uji-uji ini, sehingga ccam tidak menyentuh file yang tidak dibuatnya, seperti salinan ccam milik Anda sendiri atau tautan yang Anda buat sendiri.

Saat Anda menjalankan pintasan, nama yang Anda ketik adalah satu-satunya hal yang diperlukan ccam:

<!-- motion:shortcuts-link -->

Shell menemukan `ccwork` di `PATH`, kernel mengikuti tautan dan menjalankan biner ccam, dan `argv[0]` tetap `ccwork`. Di `app.Main`, nama apa pun selain `ccam` memuat `config.toml` sekali dan memanggil `shim.Lookup`, yang mencari akun yang `shortcut`-nya sama dengan `filepath.Base(argv[0])`. ccam lalu menjalankan langkah yang sama seperti [`ccam <name>`](/features/launch) dengan konfigurasi yang sudah dimuat, dan meneruskan `argv[1:]` ke `claude` tanpa diubah (`launcher.LaunchWith`).

Mengapa symlink dan bukan alias shell atau skrip pembungkus: symlink berfungsi di setiap shell dan di skrip, menghindari menjalankan shell tambahan, yang akan memakan sekitar 1 sampai 2 ms, dan peluncuran tetap hanya membaca satu file. `scripts/perf-check.sh` mengukur pintasan dengan `hyperfine` terhadap batas 10 ms (p95) yang sama dengan `ccam <name>`.

### Di Windows

Pintasan berupa `ccwork.exe` di `%USERPROFILE%\.local\bin`, dan ccam membuang `.exe` (tanpa membedakan huruf besar-kecil) sebelum pencarian (`shim.Name`), sehingga `ccwork`, `ccwork.exe`, dan `CCWORK.EXE` semuanya menjalankan `work`. ccam terlebih dahulu membuat tautan simbolik, yang diizinkan Windows dalam Mode Pengembang atau untuk administrator. Tanpa hak itu, ccam membuat hard link: nama kedua untuk file `ccam.exe` yang sama, yang tidak memerlukan hak khusus. Hard link dianggap milik ccam jika merupakan file yang sama dengan `ccam.exe`.

Hard link selalu menjalankan file asal pembuatannya, jadi saat `ccam update` menaruh file baru menggantikan `ccam.exe`, ccam memperbarui setiap pintasan hard link dengan `shim.Relink` dan menampilkan `✓ shortcut ccwork renewed`. Menjalankan `install.ps1` lagi melakukan hal yang sama melalui `ccam shims install` (`OK 2 shortcuts now run the new ccam`), begitu pula `ccam shims install` dan `ccam doctor --fix` untuk pintasan mana pun yang merupakan salinan versi ccam yang lebih lama. Pemeriksaan nama mencari `<name>` dengan setiap ekstensi di `PATHEXT` (`.exe`, `.cmd`, dan seterusnya) di folder `PATH` lainnya.

Ketik pintasan di shell (PowerShell, cmd.exe, Git Bash). Tautan simbolik yang dibuka dari Explorer, dialog Run, atau `Start-Process` dijalankan sebagai `ccam` biasa, yang membuka dasbor alih-alih akun itu.

## Keamanan dan kasus khusus

- **Tautan yang tidak dimiliki akun mana pun.** Jika tautan ke ccam di `~/.local/bin` memiliki nama yang tidak dipakai akun mana pun, misalnya setelah Anda mengubah `shortcut` di `config.toml` secara manual, ccam tidak beralih ke akun default. ccam berhenti dengan `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Tetapkan nama itu lagi dengan `ccam shortcut`, atau hapus tautannya.
- **File yang sudah ada.** ccam tidak pernah mengganti file yang tidak dibuatnya. `ccam shims install` melaporkan `<path> already exists and is not a ccam link, so it was left alone` dan tetap membuat tautan lainnya. `ccam add` dan `ccam import` lalu menampilkan `shortcut: ccwork (not created; see above)`.
- **Biner ccam yang dipindahkan.** Jika Anda memindahkan biner ccam, tautan yang dibuat ccam mengarah ke path yang sudah tidak ada. `ccam shims install` (atau `ccam doctor --fix`) mengarahkannya ke biner ccam yang Anda pakai untuk menjalankannya.
- **Mengganti nama dan menghapus.** `ccam rename` mempertahankan nama pintasan, dan tautannya tetap berfungsi, karena mengarah ke ccam dan ccam mencari nama itu di `config.toml`. `ccam remove` menghapus tautan akun. `ccam shortcut <account> <new>` menghapus tautan lama sebelum membuat yang baru.
- **PATH.** Saat ccam membuat tautan dan `~/.local/bin` tidak ada di `PATH` Anda, ccam menampilkan `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` diikuti baris `export PATH=...` yang perlu ditambahkan. Di Windows, ccam menampilkan baris PowerShell yang menambahkan folder itu ke `PATH` pengguna Anda sebagai gantinya.
- **Memeriksa belakangan.** `ccam doctor` melaporkan `missing shortcuts in /home/me/.local/bin: ccwork`, pintasan yang menjalankan versi ccam yang lebih lama, dan memperingatkan jika folder itu tidak ada di `PATH`. `ccam doctor --fix` membuat tautan yang hilang dan memperbarui yang usang, dan menampilkan `→ fixed` hanya untuk yang berhasil dibuatnya. Pintasan yang menjalankan ccam lain yang masih terpasang tetap berfungsi, jadi `ccam doctor` hanya mencantumkannya sebagai informasi; `ccam doctor --fix` mengarahkannya ke biner ccam yang Anda pakai untuk menjalankannya, dan `ccam add` serta `ccam import` menampilkannya sebagai ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``.

## Terkait

- [Menjalankan Claude Code](/features/launch)
- [Mengelola akun](/features/accounts)
- [Referensi perintah](/docs/commands#shortcuts)
- [Menginstal ccam](/docs/install)
