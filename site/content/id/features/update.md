---
title: Memperbarui ccam
description: ccam update mengganti ccam dengan rilis GitHub terbaru setelah pemeriksaan SHA-256; --check hanya melaporkan dan --version menginstal tag rilis tertentu.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: "598931e4801f"
---

`ccam update` mengunduh rilis ccam dari GitHub, memeriksanya terhadap `checksums.txt` milik rilis tersebut, dan mengganti biner yang sedang berjalan dengan satu rename. Gunakan jika Anda menginstal ccam dengan skrip instalasi atau dari arsip. [Dasbor](/features/dashboard) juga menanyakan rilis terbaru ke GitHub setiap kali dibuka dan menawarkannya lewat `U`, yang menjalankan langkah yang sama setelah bertanya. Perintah lain tidak pernah menghubungi GitHub, sehingga perintah-perintah itu bekerja sama saja secara offline.

## Penggunaan

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Jika tidak ada rilis yang lebih baru, ccam menampilkan `ccam 0.4.0 is up to date.` dan tidak mengubah apa pun. Jika `--version` menyebutkan versi yang sudah Anda miliki, ccam menampilkan `ccam 0.4.0 is already installed.`

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--check` | Hanya melaporkan apakah ada rilis yang lebih baru | nonaktif |
| `--version <tag>` | Menginstal tag rilis ini, seperti `v0.3.0`, yang lebih baru atau lebih lama dari versi saat ini | terbaru |
| `--dry-run` | Didukung: mencari tag terbaru (kecuali `--version` menyebutkan satu), menjalankan pemeriksaan yang sama pada biner, dan menampilkan path unduhan, checksum, dan penggantian tanpa mengunduh | nonaktif |

`--check` dan `--version` tidak dapat dipakai bersama. Tag harus berbentuk `vMAJOR.MINOR.PATCH`, dengan akhiran pre-release atau build yang opsional; `v` di depan boleh dihilangkan. Selain itu akan gagal, misalnya `--version takes a release tag such as v0.3.0, not "1.2"`.

Variabel environment `CCAM_RELEASES_URL` menggantikan `https://github.com/jhyunwoo/ccam/releases`, untuk mirror atau untuk pengujian. Strukturnya harus sama: `<url>/latest` dialihkan ke `<url>/tag/<tag>`, dan aset berada di `<url>/download/<tag>/<asset>`.

## Cara kerjanya

ccam memeriksa unduhan sebelum menyentuh biner yang terinstal, dan mengganti biner dengan satu rename. Kodenya ada di `internal/update`, yang hanya diimpor oleh `ccam update`. Jalur peluncuran tidak pernah memanggilnya.

1. ccam menemukan biner yang sedang berjalan dengan `os.Executable` dan menyelesaikan symlink dengan `filepath.EvalSymlinks` (`cli.updateCmd`). Path asli itulah yang diganti.
2. Tanpa `--version`, ccam meminta `<releases>/latest` tanpa mengikuti pengalihan (`Updater.Latest`), dan membaca tag dari header `Location` (`…/tag/v0.5.0`). Cara ini memakai halaman rilis alih-alih API GitHub, sehingga batas API sebesar 60 permintaan per jam tidak berlaku. Pencarian ini habis waktunya setelah 10 detik.
3. ccam memeriksa tag terhadap pola yang ketat (`parseVersion`). Tag akan masuk ke URL unduhan dan ke perintah shell yang ditampilkan, sehingga tidak ada yang mengandung `/`, `%`, spasi, atau karakter shell yang lolos. ccam lalu membandingkannya dengan versi yang sedang berjalan (`compareVersions`); pre-release diurutkan sebelum rilisnya.
4. ccam memutuskan apakah ccam boleh mengganti biner sama sekali (`managed()`, `writable()`; lihat di bawah).
5. ccam mengunduh `ccam_<os>_<arch>.tar.gz` (di Windows `ccam_windows_<arch>.zip`) dan `checksums.txt` secara bersamaan, dengan waktu 2 menit untuk keduanya (`Updater.Fetch`). Setiap unduhan dibatasi 256 MB.
6. SHA-256 arsip harus cocok dengan barisnya di `checksums.txt` sebelum hal lain terjadi. Baris yang tidak ada atau checksum yang berbeda menghentikan pembaruan dengan `nothing was changed`.
7. ccam membaca tar dan mengambil file biasa pertama bernama `ccam` (`extract()`; `extractZip` mengambil `ccam.exe` dari zip). Direktori, tautan, dan entri lain diabaikan, dan biner juga dibatasi 256 MB.
8. ccam menulis biner baru dengan `fsutil.AtomicWrite` (`Updater.Install`): file sementara `.ccam.tmp*` di samping biner, dengan bit izin file lama, di-flush ke disk dan di-rename menimpa file lama.

Seluruh alurnya, dari pencarian tag sampai rename:

<!-- motion:update-pipeline -->

Arsip dan `checksums.txt` diunduh secara paralel, dan tidak ada yang ditulis sampai checksum cocok. Rename bersifat atomik, sehingga crash atau disk penuh hanya meninggalkan biner lama atau biner baru, tidak pernah setengahnya. Perintah pintasan seperti `~/.local/bin/ccwork` adalah symlink ke biner ccam, sehingga langsung menjalankan versi baru tanpa perubahan apa pun.

Di Windows, file program yang sedang berjalan tidak bisa diganti atau dihapus, hanya diganti namanya. ccam menulis biner baru ke `ccam.exe.new`, mengganti nama `ccam.exe` yang sedang berjalan menjadi `ccam.exe.old`, lalu mengganti nama file baru menjadi `ccam.exe` (`Updater.Install`). Jika langkah terakhir gagal, nama biner lama dikembalikan. Pembaruan berikutnya menghapus file `ccam.exe.old` yang bisa dihapusnya; selama sebuah sesi masih berjalan dari `ccam.exe.old` sebelumnya, biner yang sedang berjalan menjadi `ccam.exe.old-<time>` sebagai gantinya. Pintasan yang berupa hard link alih-alih tautan simbolik ([Perintah pintasan](/features/shortcuts)) langsung diperbarui setelahnya, sehingga juga menjalankan versi baru; itu mencakup hard link dari versi ccam lama mana pun yang dibuat ccam, bukan hanya dari biner yang baru saja diganti.

### Instalasi yang dikelola paket

Salinan ccam yang diinstal dari `.deb` dikelola oleh apt, dan menggantinya diam-diam di luar apt akan membuat basis data paket tidak sinkron dengan file tersebut. Sebelum mengunduh, ccam memeriksa siapa pemilik biner:

<!-- motion:update-managed -->

1. Biner di `/usr/bin`, `/bin`, `/usr/sbin`, atau `/sbin` dianggap dikelola paket.
2. Di tempat lain, jika `dpkg` terinstal dan `dpkg -S <path>` berhasil, biner itu juga dianggap dikelola paket.
3. Jika tidak, ccam membuat dan menghapus file sementara `.ccam-update-*` di samping biner untuk memastikan ccam dapat menulis di sana.

Untuk biner yang dikelola paket, ccam berhenti dan menampilkan langkah-langkah apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

Di sistem selain Linux, hanya baris skrip instalasi yang ditampilkan: `irm https://ccam.dev/install.ps1 | iex` di Windows. Jika foldernya tidak dapat ditulisi, ccam menyarankan menjalankan ulang skrip instalasi, yang menginstal ccam di `~/.local/bin`.

## Keamanan dan kasus khusus

- Build pengembangan melaporkan versinya sebagai `dev` dan tidak memiliki rilis untuk dibandingkan, sehingga `ccam update` dan `ccam update --check` menolaknya. `ccam update --version vX.Y.Z` menginstal rilis di atasnya. Untuk build yang dibuat dengan `go install`, lihat [Instal ccam](/docs/install#other-ways-to-install).
- Unduhan yang gagal, aset yang tidak tersedia untuk platform Anda, atau checksum yang tidak cocok membiarkan biner terinstal seperti semula.
- Rilis yang tidak ada gagal dengan `there is no release <tag>`. Jika GitHub menolak dengan HTTP 403 atau 429, ccam menyebutkan bahwa itu kemungkinan pembatasan laju dan meminta Anda mencoba lagi nanti.
- `--check` hanya membaca. Opsi ini tidak memeriksa apakah biner dikelola paket atau dapat ditulisi.
- `--version` dengan tag yang lebih lama adalah cara untuk kembali ke versi sebelumnya.

## Terkait

- [Instal ccam](/docs/install)
- [Referensi perintah](/docs/commands)
- [Pertanyaan yang sering diajukan](/docs/faq)
