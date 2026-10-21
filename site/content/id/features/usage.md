---
title: Batas penggunaan
description: ccam usage menampilkan batas 5 jam dan mingguan semua akun sekaligus, dan ccam statusline menampilkan nama akun dan penggunaannya di Claude Code.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

`ccam usage` menampilkan seberapa banyak batas 5 jam, batas mingguan, dan batas mingguan khusus setiap model (seperti Fable) yang sudah digunakan setiap akun, sebagai satu bilah per batas beserta waktu hingga direset. Gunakan sebelum memulai tugas panjang, atau saat Claude Code memberi tahu bahwa batas sudah dekat, untuk melihat akun mana yang masih punya ruang. `ccam statusline` adalah perintah status line untuk Claude Code: perintah ini menampilkan nama akun dan penggunaan di dalam Claude Code serta menyimpan angkanya untuk `ccam usage`.

## Penggunaan

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Setiap akun mendapat satu baris berisi nomornya, namanya, dan umur angka-angkanya; `(cached)` berarti angkanya berasal dari file yang tersimpan, bukan dari endpoint. Di bawahnya, satu baris per batas menampilkan bilah 20 sel, persentase, dan waktu hingga jendela direset: `5-hour`, `Weekly` (semua model), lalu jendela mingguan khusus setiap model, seperti `Fable`. Jendela yang tidak dimiliki akun menampilkan `-`. Persentase berubah kuning mulai 80% dan merah pada 100%. Di terminal, setiap akun langsung muncul sebagai `fetching…`, dan baris-barisnya digambar begitu akun itu menjawab. Baris terakhir menyebut akun yang akan dipilih [`ccam next --auto`](/features/next) dengan aturan yang sama (`auto_max_usage`, `auto_strategy`, sesi yang sedang berjalan, hanya akun dari organisasi akun yang terakhir dipakai kecuali `cross_org` diatur, setiap jendela yang sudah 100% dilewati, sehingga akun yang kuota mingguannya sudah penuh tidak pernah disebut), dan hanya muncul jika Anda memiliki lebih dari satu akun. Akun yang belum login menampilkan `not logged in → ccam login <name>`, dan akun dengan kunci API menampilkan `no limits (API)`.

Setelah `ccam statusline install`, Claude Code menampilkan baris seperti `[work] · Opus · 5h 24% · 7d 41%`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--json` | Menampilkan laporan sebagai JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | nonaktif |
| `--watch[=N]` | Menggambar ulang setiap `N` detik sampai Anda menekan Ctrl+C. `--watch` saja, atau nilai 0 atau kurang, berarti 60 | nonaktif |
| `-v`, `--verbose` | Menambahkan baris `source:` di bawah setiap akun: asal angka-angkanya (endpoint penggunaan, Claude Code, status line, atau file yang tersimpan) | nonaktif |
| `--no-cache` | Mengabaikan cache 60 detik dan bertanya ulang ke endpoint | nonaktif |
| `--source live\|claude\|cache` | `live` bertanya ke endpoint dengan token akun; `claude` meminta Claude Code mencarinya (lihat di bawah); `cache` hanya menampilkan angka yang tersimpan | `usage_source`, jika tidak ada `live` |
| `--refresh` | Memperbarui token akses yang kedaluwarsa, hanya untuk akun yang tidak sedang berjalan | nonaktif |
| `--dry-run` | Diterima: `ccam usage` tidak mengubah apa pun selain cache-nya sendiri. Ditolak jika dipakai bersama `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | nonaktif |

Nilai `--source` lainnya gagal dengan `--source must be live, claude or cache`.

### Penggunaan dari Claude Code sendiri

Dengan `usage_source = "claude"` di `config.toml` (atau `--source claude`), ccam sama sekali tidak membaca token login. Untuk setiap akun, ccam menjalankan `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` di folder akun itu dan mengirim permintaan `get_usage` milik Claude Code lewat input standar (Claude Code 2.1.270 atau lebih baru); Claude Code mencari penggunaannya dengan login miliknya sendiri dan menjawab dengan jendela 5 jam, mingguan, dan per model. `--safe-mode` mencegah hook, server MCP, dan plugin akun itu dimulai, dan tidak ada sesi yang disimpan. Setiap akun memerlukan satu kali start Claude Code, sekitar satu detik, sehingga sumber ini lebih lambat daripada `live`; semua akun ditanya sekaligus, masing-masing dalam 20 detik, dan cache 60 detik berlaku seperti biasa. Claude Code menandai permintaan ini sebagai eksperimental, jadi jika permintaan itu berhenti dijawab, ccam menampilkan angka yang tersimpan dan `ccam doctor` melaporkannya. `--dry-run` tidak menjalankan Claude Code (karena Claude Code menulis ke folder akun) dan menampilkan angka yang tersimpan. `usage_source` menerima `live` atau `claude`; `cache` hanya bisa dipilih lewat `--source`. `ccam usage` tidak menerima argumen.

`ccam statusline` membaca JSON status line milik Claude Code dari input standar. Satu-satunya opsinya adalah `--wrap <command>`: perintah status line yang sudah ada, yang baris keluaran pertamanya diberi nama akun di depannya. `ccam statusline install` tidak menerima opsi. Keduanya menolak `--dry-run` dengan `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Cara kerjanya

Angkanya berasal dari dua tempat: endpoint penggunaan Anthropic yang tidak didokumentasikan, `GET https://api.anthropic.com/api/oauth/usage`, yang dipanggil dengan login masing-masing akun, dan nilai yang diberikan Claude Code ke perintah status line-nya. Semuanya ada di `internal/usage`.

1. `cli.usageCmd` memuat `config.toml` dan membuat satu `usage.Account` per akun: foldernya `~/.local/share/ccam/profiles/<name>/`, folder sesinya (untuk melihat apakah akun sedang berjalan), dan apakah akun itu diimpor di tempatnya.
2. `Fetcher.FetchAll` menjalankan satu goroutine per akun. Semuanya berbagi satu `http.Client`, sehingga koneksi ke endpoint dipakai ulang. Menyambung boleh memakan waktu 2 detik, dan setiap akun mendapat total 3 detik. Klien tidak mengikuti pengalihan, karena pengalihan akan mengirim ulang token ke tempat lain, dan hanya membaca paling banyak 1 MB jawaban.
3. Begitu setiap akun selesai, barisnya digambar ulang. Akun yang gagal atau tidak menjawab tepat waktu beralih ke angka yang tersimpan, sehingga satu akun yang lambat tidak pernah menahan akun lain.

Setiap akun berjalan di goroutine-nya sendiri terhadap endpoint; akun yang melewati batas waktu menampilkan angka tersimpannya sebagai gantinya.

<!-- motion:usage-fetch -->

`personal` dan `work` menjawab tepat waktu dan menampilkan angka langsung. `side` membutuhkan lebih dari 3 detik, sehingga `fallback()` membaca dua file cache-nya dan memakai yang paling akhir ditulis. `usage.Table()` menggambar satu blok per akun dan `usage.Best()` memilih akun yang disebut di baris terakhir.

Untuk satu akun, `Fetcher.Fetch` melakukan pemeriksaan ini secara berurutan:

<!-- motion:usage-cache -->

1. Dengan `--source cache`, langsung menuju `fallback()`.
2. Kecuali `--no-cache` diatur, hasil langsung di `~/.local/state/ccam/usage/<name>.json` yang berumur kurang dari 60 detik (`CacheTTL`) dikembalikan apa adanya. Menjalankan `ccam usage` dua kali berturut-turut hanya memakan satu permintaan per akun.
3. `readCredentials()` membaca `.credentials.json` (atau Rantai Kunci (Keychain) macOS) tanpa mengubahnya, dan hanya kunci `claudeAiOauth` yang didekode. Tidak adanya login langganan berarti `not logged in`, atau `no limits (API)` jika `.claude.json` berisi kunci API.
4. Token akses yang kedaluwarsa tidak dipakai. Tanpa `--refresh`, akun menampilkan angka tersimpannya dengan tanda kedaluwarsa.
5. Jika tidak, ccam memanggil endpoint. 401 atau 403 dianggap sebagai login yang kedaluwarsa, status lain atau isi yang tidak dapat dibaca dianggap galat, dan keduanya beralih ke angka yang tersimpan. 429 berarti endpoint terlalu sering ditanya: ccam menulis `<name>.pause` berisi waktu untuk bertanya lagi (`Retry-After`, atau 5 menit, paling lama satu jam), menampilkan angka yang tersimpan dengan `rate limited by the usage endpoint; asking again after 15:10`, dan tidak memanggil endpoint untuk akun itu sampai saat tersebut, bahkan dengan `--no-cache`. Mencoba lagi seketika hanya akan membuat batas itu terus berlaku.
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
4. ccam menampilkan `[name] · model · 5h N% · 7d N%`. Dengan `--wrap`, ccam menjalankan perintah Anda melalui `sh -c` dengan input yang sama dan menampilkan `[name]` diikuti baris pertama perintah itu sebagai gantinya. Di Windows, ccam menjalankan perintah itu dengan cara yang sama seperti Claude Code: dengan Git Bash jika terpasang, atau dengan PowerShell jika tidak.
5. Saat penggunaan 5 jam atau mingguan mencapai 90% atau lebih dan ada akun lain, baris itu diakhiri dengan `→ ccam next --auto --continue`, misalnya `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Perintah yang dibungkus diberi waktu 2 detik. Setelah itu perintah dihentikan dan baris hanya menampilkan nama akun, sehingga perintah yang lambat tidak dapat memblokir status line Claude Code.

`ccam statusline install` menulis `statusLine` ke `settings.json` setiap akun. Jika sebuah akun sudah memiliki perintah status line, ccam mempertahankannya sebagai `--wrap '<old command>'`. Di Windows, path ke ccam ditulis dengan garis miring biasa dan `~` untuk folder home (`~/.local/bin/ccam.exe statusline`), yang dipahami oleh Git Bash maupun PowerShell. `settings.json` yang dibagikan melalui grup ditulis sekali, melalui tautannya, sehingga setiap akun dalam grup mendapatkannya. Menjalankannya lagi menampilkan `(already set)` untuk file yang sudah memakai ccam.

## Keamanan dan kasus khusus

- `ccam usage` tidak pernah mengubah file token kecuali Anda memberikan `--refresh`. Jika login sudah kedaluwarsa, jalankan akun tersebut sekali dengan `ccam <name>` dan Claude Code akan memperbaruinya.
- `--refresh` tidak pernah menyentuh akun yang sedang berjalan, akun yang diimpor, atau login Rantai Kunci (Keychain). Pembaruan yang ditolak menampilkan angka tersimpan beserta alasannya.
- Endpoint penggunaan tidak didokumentasikan secara publik dan dapat berubah. Jika berubah, Anda tetap melihat angka yang tersimpan, dan `ccam doctor` melaporkan perubahan format.
- Jika tidak ada angka tersimpan dan tidak ada jawaban, baris akun menampilkan `lookup failed:` beserta alasannya, misalnya `no saved usage` dengan `--source cache`.
- File cache hanya berisi persentase dan waktu reset, tidak pernah token.

## Terkait

- [Batas penggunaan dan perpindahan akun](/docs/usage-limits)
- [Beralih ke akun berikutnya](/features/next)
- [Alur kerja batas penggunaan](/guides/usage-limit-workflow)
- [Perintah](/docs/commands)
