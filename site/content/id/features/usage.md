---
title: Batas penggunaan
description: ccam usage menampilkan batas 5 jam dan mingguan semua akun sekaligus, dan ccam statusline menampilkan nama akun dan penggunaannya di Claude Code.
updated: 2026-10-19
order: 7
translation: machine
sourceHash: 8306a8f6d950
---

`ccam usage` menampilkan seberapa banyak batas 5 jam dan mingguan yang sudah digunakan setiap akun, beserta waktu hingga setiap jendela direset. Gunakan sebelum memulai tugas panjang, atau saat Claude Code memberi tahu bahwa batas sudah dekat, untuk melihat akun mana yang masih punya ruang. `ccam statusline` adalah perintah status line untuk Claude Code: perintah ini menampilkan nama akun dan penggunaan di dalam Claude Code serta menyimpan angkanya untuk `ccam usage`.

## Penggunaan

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Di terminal, tabel langsung muncul dengan `fetching…` di setiap baris, dan setiap baris digambar ulang begitu akunnya menjawab. `As of` menunjukkan umur angka-angka itu; `(cached)` berarti angkanya berasal dari file yang tersimpan, bukan dari endpoint. Baris terakhir menyebut akun dengan penggunaan 5 jam terendah di bawah 100%, dan hanya muncul jika Anda memiliki lebih dari satu akun. Akun yang belum login menampilkan `not logged in → ccam login <name>`, dan akun dengan kunci API menampilkan `no limits (API)`.

Setelah `ccam statusline install`, Claude Code menampilkan baris seperti `[work] · Opus · 5h 24% · 7d 41%`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--json` | Menampilkan laporan sebagai JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | nonaktif |
| `--watch[=N]` | Menggambar ulang setiap `N` detik sampai Anda menekan Ctrl+C. `--watch` saja, atau nilai 0 atau kurang, berarti 60 | nonaktif |
| `-v`, `--verbose` | Menampilkan batas mingguan per model dan baris `source:` (`live`, `cache`, atau `statusline`) | nonaktif |
| `--no-cache` | Mengabaikan cache 60 detik dan bertanya ulang ke endpoint | nonaktif |
| `--source live\|cache` | `cache` tidak pernah bertanya ke endpoint dan hanya menampilkan angka yang tersimpan; `live` adalah perilaku default | live |
| `--refresh` | Memperbarui token akses yang kedaluwarsa, hanya untuk akun yang tidak sedang berjalan | nonaktif |
| `--dry-run` | Diterima: `ccam usage` tidak mengubah apa pun selain cache-nya sendiri. Ditolak jika dipakai bersama `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | nonaktif |

Nilai `--source` lainnya gagal dengan `--source must be live or cache`. `ccam usage` tidak menerima argumen.

`ccam statusline` membaca JSON status line milik Claude Code dari input standar. Satu-satunya opsinya adalah `--wrap <command>`: perintah status line yang sudah ada, yang baris keluaran pertamanya diberi nama akun di depannya. `ccam statusline install` tidak menerima opsi. Keduanya menolak `--dry-run` dengan `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Cara kerjanya

Angkanya berasal dari dua tempat: endpoint penggunaan Anthropic yang tidak didokumentasikan, `GET https://api.anthropic.com/api/oauth/usage`, yang dipanggil dengan login masing-masing akun, dan nilai yang diberikan Claude Code ke perintah status line-nya. Semuanya ada di `internal/usage`.

1. `cli.usageCmd` memuat `config.toml` dan membuat satu `usage.Account` per akun: foldernya `~/.local/share/ccam/profiles/<name>/`, folder sesinya (untuk melihat apakah akun sedang berjalan), dan apakah akun itu diimpor di tempatnya.
2. `Fetcher.FetchAll` menjalankan satu goroutine per akun. Semuanya berbagi satu `http.Client`, sehingga koneksi ke endpoint dipakai ulang. Menyambung boleh memakan waktu 2 detik, dan setiap akun mendapat total 3 detik. Klien tidak mengikuti pengalihan, karena pengalihan akan mengirim ulang token ke tempat lain, dan hanya membaca paling banyak 1 MB jawaban.
3. Begitu setiap akun selesai, barisnya digambar ulang. Akun yang gagal atau tidak menjawab tepat waktu beralih ke angka yang tersimpan, sehingga satu akun yang lambat tidak pernah menahan akun lain.

Setiap akun berjalan di goroutine-nya sendiri terhadap endpoint; akun yang melewati batas waktu menampilkan angka tersimpannya sebagai gantinya.

<!-- motion:usage-fetch -->

`personal` dan `work` menjawab tepat waktu dan menampilkan angka langsung. `side` membutuhkan lebih dari 3 detik, sehingga `fallback()` membaca dua file cache-nya dan memakai yang paling akhir ditulis. `usage.Table()` menggambar baris-barisnya dan `usage.Best()` memilih akun yang disebut di baris terakhir.

Untuk satu akun, `Fetcher.Fetch` melakukan pemeriksaan ini secara berurutan:

<!-- motion:usage-cache -->

1. Dengan `--source cache`, langsung menuju `fallback()`.
2. Kecuali `--no-cache` diatur, hasil langsung di `~/.local/state/ccam/usage/<name>.json` yang berumur kurang dari 60 detik (`CacheTTL`) dikembalikan apa adanya. Menjalankan `ccam usage` dua kali berturut-turut hanya memakan satu permintaan per akun.
3. `readCredentials()` membaca `.credentials.json` (atau Rantai Kunci (Keychain) macOS) tanpa mengubahnya, dan hanya kunci `claudeAiOauth` yang didekode. Tidak adanya login langganan berarti `not logged in`, atau `no limits (API)` jika `.claude.json` berisi kunci API.
4. Token akses yang kedaluwarsa tidak dipakai. Tanpa `--refresh`, akun menampilkan angka tersimpannya dengan tanda kedaluwarsa.
5. Jika tidak, ccam memanggil endpoint. 401 atau 403 dianggap sebagai login yang kedaluwarsa, status lain atau isi yang tidak dapat dibaca dianggap galat, dan keduanya beralih ke angka yang tersimpan.
6. Jawaban yang baik ditulis ke `<name>.json` bersama waktu pengukurannya.

`fallback()` membaca `<name>.json` dan `<name>.statusline.json` sekaligus dan mempertahankan yang lebih baru. Sebelum menampilkan angka tersimpan, fungsi ini memeriksa waktu reset setiap jendela: jendela yang sudah direset sejak angka disimpan menampilkan 0%, karena persentase lama sudah tidak berlaku. Itulah sebabnya `personal` di atas menampilkan `0%` tanpa waktu reset.

Endpoint memiliki dua format respons. Yang lebih baru memiliki array `limits` dengan entri `session` (5 jam), `weekly_all`, dan `weekly_scoped` (per model). Yang lebih lama memiliki objek `five_hour`, `seven_day`, `seven_day_opus`, dan `seven_day_sonnet`. `ParseLive` membaca keduanya: array `limits` yang menang, dan objek lama mengisi apa yang tidak dicakupnya. Isi tanpa batas sama sekali dilaporkan sebagai format yang berubah, yang akan diberitahukan oleh `ccam doctor`.

### Memperbarui token

Refresh token berotasi: memperbaruinya memberi Anda refresh token baru dan membatalkan yang lama. Jika ccam memperbarui login saat Claude Code sedang berjalan sebagai akun yang sama, sesi itu akan keluar. Karena itu `ccam usage` membaca token dan tidak pernah menulisnya, kecuali Anda memberikan `--refresh`. Dalam hal itu `Fetcher.refresh`:

1. menolak akun yang didaftarkan di tempatnya dengan `ccam import`, karena `claude` biasa mungkin memakai login itu tanpa sepengetahuan ccam;
2. menolak login yang tersimpan di Rantai Kunci (Keychain) macOS;
3. mengambil file lock, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (di bawah `~/.local/state/ccam/locks/` jika variabel itu tidak diatur), sehingga dua proses ccam tidak pernah memperbarui token yang sama dua kali;
4. di bawah lock, menolak jika ada sesi akun yang sedang berjalan, dan melewati permintaan jika proses lain sudah memperbarui token;
5. mengirim refresh token ke endpoint token dan menulis token baru kembali ke `.credentials.json` dengan `fsutil.AtomicWrite`, dengan mempertahankan setiap kunci lain di file itu.

### Pengumpul status line

Claude Code menjalankan perintah status line dan mengirimkan dokumen JSON kepadanya melalui pipe. ccam membaca nama model dan `rate_limits` darinya:

<!-- motion:usage-statusline -->

1. `runStatusline` menemukan akun dari `CCAM_PROFILE`, yang diatur ccam saat menjalankan Claude Code, atau jika tidak ada, dengan mencocokkan `CLAUDE_CONFIG_DIR` dengan folder akun.
2. `ParseStatusLine` membaca `rate_limits.five_hour` dan `rate_limits.seven_day`. JSON tidak berisi batas laju untuk akun kunci API, atau sebelum respons pertama sesi; dalam hal itu tidak ada yang disimpan.
3. `SaveStatusLine` menulis `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` memakainya jika endpoint tidak menjawab.
4. ccam menampilkan `[name] · model · 5h N% · 7d N%`. Dengan `--wrap`, ccam menjalankan perintah Anda melalui `sh -c` dengan input yang sama dan menampilkan `[name]` diikuti baris pertama perintah itu sebagai gantinya.

Perintah yang dibungkus diberi waktu 2 detik. Setelah itu perintah dihentikan dan baris hanya menampilkan nama akun, sehingga perintah yang lambat tidak dapat memblokir status line Claude Code.

`ccam statusline install` menulis `statusLine` ke `settings.json` setiap akun. Jika sebuah akun sudah memiliki perintah status line, ccam mempertahankannya sebagai `--wrap '<old command>'`. `settings.json` yang dibagikan melalui grup ditulis sekali, melalui tautannya, sehingga setiap akun dalam grup mendapatkannya. Menjalankannya lagi menampilkan `(already set)` untuk file yang sudah memakai ccam.

## Keamanan dan kasus khusus

- `ccam usage` tidak pernah mengubah file token kecuali Anda memberikan `--refresh`. Jika login sudah kedaluwarsa, jalankan akun tersebut sekali dengan `ccam <name>` dan Claude Code akan memperbaruinya.
- `--refresh` tidak pernah menyentuh akun yang sedang berjalan, akun yang diimpor, atau login Rantai Kunci (Keychain). Pembaruan yang ditolak menampilkan angka tersimpan beserta alasannya.
- Endpoint penggunaan tidak didokumentasikan secara publik dan dapat berubah. Jika berubah, Anda tetap melihat angka yang tersimpan, dan `ccam doctor` melaporkan perubahan format.
- Jika tidak ada angka tersimpan dan tidak ada jawaban, baris menampilkan `lookup failed:` beserta alasannya, misalnya `no saved usage` dengan `--source cache`.
- File cache hanya berisi persentase dan waktu reset, tidak pernah token.

## Terkait

- [Batas penggunaan dan perpindahan akun](/docs/usage-limits)
- [Beralih ke akun berikutnya](/features/next)
- [Alur kerja batas penggunaan](/guides/usage-limit-workflow)
- [Perintah](/docs/commands)
