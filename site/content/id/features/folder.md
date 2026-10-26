---
title: Akun per folder
description: ccam folder membuat folder dan semua isinya memakai satu akun saat Anda menjalankan Claude Code tanpa nama akun, seperti work untuk kode kantor Anda.
updated: 2026-10-26
order: 11
translation: machine
sourceHash: "25f82da64580"
---

Setelah `ccam folder work ~/code/company`, setiap peluncuran tanpa nama akun di `~/code/company` atau subfoldernya menjalankan Claude Code sebagai `work`. Di luar folder itu, akun default tetap dijalankan seperti sebelumnya. Menyebutkan nama akun, dengan `ccam personal` atau pintasan seperti `ccpersonal`, selalu menjalankan akun tersebut.

## Penggunaan

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Peluncuran berikut memakai akun folder tersebut:

- `ccam -p "…"` dan peluncuran lain yang argumen pertamanya adalah opsi untuk Claude Code
- `ccam run` tanpa nama, dan `ccam --dry-run`
- `ccam --fallback -p "…"` untuk percobaan pertamanya
- dasbor, `ccam` saja, yang dimulai dengan kursor pada akun folder tersebut, bahkan saat Anda mencapai folder itu melalui tautan simbolik

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `<account>` | Nama atau nomor akun yang dipakai | |
| `[folder]` | Folder, relatif atau absolut | folder saat ini |
| `--remove` | Menghapus entri untuk folder itu persis | nonaktif |
| `--dry-run` | Didukung: menampilkan perubahan tanpa menyimpannya | nonaktif |

## Cara kerjanya

Setiap folder adalah satu entri `[[folder]]` di `~/.config/ccam/config.toml`, file yang memang dibaca setiap peluncuran:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` mengubah folder menjadi path absolut dengan tautan simbolik yang sudah diurai, lalu menyimpannya. Menetapkan folder lagi akan mengganti akunnya.
2. Peluncuran tanpa nama memanggil `launcher.HereProfile`. Jika tidak ada entri, akun default langsung dipakai.
3. Jika ada, direktori kerja dibandingkan dengan setiap `path` sebagai string: folder itu sendiri atau folder di bawahnya, tidak pernah folder yang hanya diawali huruf yang sama (`/code/company` tidak mencakup `/code/companyx`). Entri yang paling dalam menang, sehingga `~/code` bisa memakai `personal` sementara `~/code/company` memakai `work`.
4. Hanya jika direktori kerja apa adanya tidak cocok dengan apa pun, barulah ccam mengurai tautan simboliknya dan membandingkan lagi.
5. Di Windows dan macOS, perbandingan mengabaikan huruf besar dan kecil, seperti yang biasanya dilakukan sistem file keduanya.

Tidak ada file selain `config.toml` yang dibaca, sehingga peluncuran di dalam folder tidak menambah waktu yang terukur (4,0 ms dengan atau tanpa entri, `hyperfine`, 300 kali jalan).

## Keamanan dan kasus khusus

- `ccam rename` membawa entri folder ke nama baru, dan `ccam remove` menghapus entri milik akun yang dihapus.
- Entri yang akunnya sudah tidak ada menghentikan peluncuran tanpa nama di folder itu: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` melaporkan entri seperti itu, serta entri yang foldernya sudah hilang.
- Jika Anda menjalankan `ccam folder --remove` di subfolder yang mewarisi akunnya dari folder induk, ccam menyebutkan folder induk itu dan tidak menghapus apa pun.
- `ccam next` memilih berdasarkan rotasi atau penggunaan dan tidak melihat folder.

## Terkait

- [Menjalankan Claude Code](/features/launch)
- [Perintah pintasan](/features/shortcuts)
- [Mengelola akun](/features/accounts)
- [Referensi perintah](/docs/commands)
