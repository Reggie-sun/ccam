---
title: Memeriksa dan memperbaiki
description: Periksa tautan pembagian, izin file, login, dan versi Claude Code dengan ccam doctor, lalu perbaiki yang aman diperbaiki dengan ccam doctor --fix.
updated: 2026-10-23
order: 9
translation: machine
sourceHash: "11ad34630dcf"
---

`ccam doctor` memeriksa instalasi ccam Anda: biner Claude Code, folder setiap akun, login dan tautan pembagian, server MCP bersama, endpoint penggunaan, dan perintah pintasan Anda. `ccam doctor --fix` memperbaiki apa yang bisa diperbaiki tanpa kehilangan data dan terlebih dahulu memindahkan setiap file yang digantinya ke folder cadangan. Jalankan perintah ini saat `ccam <name>` berhenti karena masalah tautan pembagian, setelah Anda memulihkan file secara manual, atau setelah Claude Code diperbarui.

## Penggunaan

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

Baris pertama menandai semua akun sekaligus, sehingga Anda melihat letak masalahnya sebelum membaca yang lain. Di bawahnya, setiap akun mendapat satu kotak, dan temuan yang bukan milik akun mana pun (Claude Code, variabel environment, plugin, server MCP, pintasan, aturan folder, grup, dan pengaturan) masuk ke kotak `General`. Judul kotak menampilkan tanda terburuk akun itu dan jumlah masalahnya. Setiap baris di dalamnya adalah satu area yang diperiksa, selalu `Folder`, `Login`, dan `Share links` dengan urutan itu, lalu area lainnya, dengan tanda: `✔` OK, `ℹ` informasi, `⚠` peringatan, `✘` galat. Folder yang tidak bermasalah menampilkan `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Dengan `--fix`, setiap temuan yang berhasil diperbaiki diakhiri dengan `→ fixed`, dan tidak lagi dihitung sebagai masalah:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

Di terminal, kotak-kotak menyesuaikan lebar terminal dan pesan panjang dibungkus ke baris berikutnya, dengan warna terminal itu sendiri (`NO_COLOR` mematikannya). Sebuah kotak tidak pernah lebih sempit dari 29 kolom, sehingga tidak ada kata yang terpotong di terminal yang sempit. Jika dikirim ke pipe atau file, keluarannya tanpa warna dan tidak ada pesan yang dibungkus atau dipotong, sehingga `ccam doctor | grep settings.json` menemukan baris utuhnya; setiap kotak lalu selebar baris terpanjang. Tab Health di dasbor menampilkan kotak-kotak yang sama.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--fix` | Memperbaiki yang aman diperbaiki (file yang diubah dicadangkan) | nonaktif |
| `--dry-run` (global) | Menjalankan semua pemeriksaan tetapi tidak memperbaiki apa pun, bahkan dengan `--fix` | nonaktif |

`ccam doctor` tidak menerima argumen. Perintah ini menghormati `--dry-run`: `ccam --dry-run doctor --fix` menampilkan temuan `ccam doctor`, menambahkan apa yang akan disisihkan oleh `--fix`, dan tidak mengubah apa pun:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

Kode keluarnya `1` jika setidaknya satu temuan `✘` tidak diperbaiki, dan `0` jika tidak. Peringatan dan informasi tidak pernah menggagalkan proses. Jika proses gagal tanpa `--fix`, ccam diakhiri dengan baris kosong dan `` `ccam doctor --fix` repairs what it can. ``. Konfigurasi yang tidak dapat dimuat juga berakhir dengan `1` dan galat `ccam:`.

## Cara kerjanya

`internal/cli/doctor.go` memuat `~/.config/ccam/config.toml` dan memanggil `doctor.Run` di `internal/doctor/doctor.go`, yang menjalankan pemeriksaan dengan urutan tetap dan mengembalikan temuan, masing-masing dengan akun pemiliknya dan area yang diperiksanya. `tui.HealthReport` mengelompokkannya ke dalam kotak-kotak. `--fix` diteruskan sebagai `Options.Fix`, yang diatur ke false saat dry run.

<!-- motion:doctor-checks -->

Blok tengah dijalankan sekali per akun. Pemeriksaannya, dan apa yang dilakukan `--fix` untuk masing-masing:

1. **Versi Claude Code** (`checkClaude`). Mencari biner `claude` dengan cara yang sama seperti peluncur dan menjalankan `claude --version`. Biner yang tidak ada, `claude_bin` yang bukan file yang dapat dieksekusi, atau `--version` yang gagal adalah galat. Tidak ada perbaikan.
2. **Info masuk di environment** (`checkAuthEnv`). Memberi peringatan jika `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY`, atau `ANTHROPIC_AUTH_TOKEN` diatur: setiap akun yang dijalankan dari shell itu mewarisinya, dan Claude Code lebih mengutamakannya daripada login milik akun. Tidak ada perbaikan; hapus sendiri variabel itu.
3. **Folder** (`checkProfile`). Folder akun yang tidak ada adalah galat dan mengakhiri pemeriksaan akun tersebut. Izin dengan bit group atau other yang diatur adalah peringatan; `--fix` menjalankan `chmod 0700`. `~/.claude` yang diimpor dengan `.claude.json` sendiri di samping `~/.claude.json` adalah informasi: `ccam main` dan `claude` biasa lalu menyimpan detail login dan server MCP secara terpisah. Windows tidak punya bit izin, jadi pemeriksaan ini dan pemeriksaan `0600` di bawah dilewati di sana; folder di profil pengguna Anda dilindungi oleh daftar aksesnya.
4. **File login** (`checkCredentials`). `.credentials.json` diperiksa dengan `lstat`. Symlink adalah galat, karena setiap akun harus memiliki tokennya sendiri, dan tidak diperbaiki. Mode selain `0600` adalah peringatan; `--fix` menjalankan `chmod 0600`. Lalu login dibaca (di macOS dari Rantai Kunci (Keychain) terlebih dahulu). Tidak adanya login, atau file login tanpa login di dalamnya, adalah peringatan yang mengarah ke `ccam login <name>`. Refresh token yang kedaluwarsa, atau login yang dihapus Claude Code setelah server menolaknya, adalah galat `` login expired → `ccam login <name>` ``; token akses yang kedaluwarsa hanya informasi, karena Claude Code memperbaruinya saat dijalankan berikutnya. Login yang disimpan di Windows Credential Manager adalah informasi, karena ccam tidak dapat membacanya. Login yang valid menampilkan email yang saat ini dipakai folder itu untuk login dan paketnya (`login valid (dev@company.com, Claude Team · Acme)`). Tidak satu pun dari ini diperbaiki; `ccam login --expired` melakukan login untuk setiap akun yang memerlukannya.
5. **Tautan pembagian** (`share.Check`). Membandingkan setiap path di `share.Layout` dengan grup-grup akun, hanya menggunakan `lstat` dan `readlink`. Pemeriksaan ini melaporkan lima jenis: tautan yang hilang, tautan yang mengarah ke tempat lain, file asli di tempat yang seharusnya berisi tautan, tautan yang tertinggal untuk item yang kini lokal, dan tautan yang benar tetapi folder grupnya hilang (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; peluncuran tidak memeriksa hal ini). Tautan yang hilang dan folder grup yang hilang adalah peringatan, sisanya galat. Penyalinan kembali dari grup yang terhenti di tengah jalan juga dilaporkan: `--fix` menempatkan salinan bertahap yang lengkap, tetapi jika path itu sudah dibuat ulang sejak saat itu (oleh peluncuran), ia menyebutkan folder berisi salinan lengkap itu agar Anda yang membereskannya. Path yang ditambahkan ke sebuah item di versi ccam yang lebih baru (`file-history/`, `tasks/`, dan `paste-cache/` untuk riwayat, `keybindings.json` untuk pengaturan) hanya menjadi peringatan selama akun masih memiliki salinannya sendiri: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Peluncuran tidak pernah berhenti karenanya.
6. **Entri tak dikenal** (`checkUnknownEntries`). Mendaftar nama tingkat atas di folder akun yang tidak dikenal ccam (`share.Known`), sebagai petunjuk bahwa tata letak Claude Code berubah. Hanya informasi.
7. **Path plugin** (`checkPlugins`). Registri plugin (`installed_plugins.json`, `known_marketplaces.json`) menyimpan path absolut. Di `plugins/` setiap grup, path yang melewati folder satu akun akan rusak saat akun itu dihapus atau diganti namanya. `--fix` menulis ulang path itu ke folder grup dengan `profile.RewritePluginPaths`, yaitu penulisan atomik file registri.
8. **Server MCP** (`checkMCP`). Memuat `mcp.json` setiap grup sekali. File yang tidak dapat dibaca adalah galat; server yang gagal `mcp.Validate` adalah peringatan. Server di `~/.claude.json` milik `claude` biasa yang tidak dimuat akun mana pun adalah peringatan yang menyebutkan perintah `ccam migrate --mcp` untuk menyalinnya (dengan `--group <a|b>` jika beberapa grup membagikan MCP). Server yang hanya dimiliki satu akun dalam sebuah grup adalah informasi beserta perintah `ccam mcp import` untuk membagikannya, dan akun yang belum masuk ke server MCP yang sudah dimasuki akun lain mendapat peringatan untuk masuk dengan `/mcp` (ccam tidak pernah menyalin info masuk). Tidak ada perbaikan.
9. **Format respons penggunaan** (`checkUsageFormat`). Menanyakan endpoint penggunaan, tanpa cache, untuk satu akun demi satu akun sampai ada yang menjawab, semuanya dalam satu batas 4 detik. Dengan `usage_source = "claude"`, pemeriksaan ini menanyakan Claude Code untuk satu akun sebagai gantinya. Jika respons tidak lagi dapat di-parse, pemeriksaan ini memperingatkan bahwa `ccam usage` beralih ke cache status line; jika tidak ada akun yang menjawab, hal itu disampaikan sebagai informasi (`could not check the usage response format: …`). Ini satu-satunya pemeriksaan yang memakai jaringan.
10. **Pintasan** (`checkShims`). Hanya jika ada akun yang memiliki pintasan. Tautan yang hilang di `~/.local/bin` adalah peringatan, begitu juga pintasan buatan ccam yang masih menjalankan ccam yang lebih lama, yang sudah diganti atau dipindahkan (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` memanggil `shim.Install`, yang hanya membuat dan memperbarui tautan yang dibuat ccam dan dicatat di `~/.local/state/ccam/shortcuts.json`, dan tidak pernah menyentuh file lain. Peringatan itu diakhiri `→ fixed` hanya jika setiap tautan tersebut berhasil dibuat; jika ada yang gagal, alasannya ditambahkan. Pintasan yang menjalankan ccam lain yang masih terpasang tetap berfungsi, dan merupakan informasi; `--fix` mengarahkannya ke ccam ini. Folder bin yang tidak ada di `PATH` adalah peringatan dan tidak diperbaiki.
11. **Folder proyek** (`checkFolders`). Hanya jika ada entri [`ccam folder`](/features/folder). Entri yang akunnya tidak ada adalah galat, dan entri yang foldernya sudah hilang adalah peringatan. Keduanya tidak diperbaiki, karena menghapus entri akan mengubah akun mana yang dijalankan di sana.
12. **Grup dan pengaturan** (`checkGroups`). Anggota grup yang bukan akun, akun yang namanya kini menjadi perintah ccam (jalankan dengan `ccam run <name>` atau ganti namanya), `theme` atau `usage_source` yang tidak dikenal, serta baris `config.toml` yang tidak dibaca ccam adalah peringatan. Yang terakhir biasanya berupa pengaturan yang ditulis di bawah tabel `[[group]]`, `[[profile]]`, atau `[[folder]]`, tempat pengaturan itu menjadi milik tabel tersebut (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); perintah yang mengubah konfigurasi menolak sampai hal itu diperbaiki. Tidak ada perbaikan.
13. **Login** (`checkIdentities`). Membaca akun mana yang dipakai setiap folder untuk login dari `.claude.json`-nya. Folder yang kini login ke akun lain daripada yang dicatat `ccam login` (misalnya setelah `/login` di dalam sesi) adalah peringatan; akun tanpa catatan login, seperti akun yang diimpor atau dibuat oleh ccam versi lama, dibandingkan berdasarkan email yang dicatatnya. Folder yang login ke akun Claude yang sama dengan folder lain juga merupakan peringatan: keduanya berbagi batas penggunaan. Folder yang memegang token login yang sama dengan folder lain, yaitu folder yang disalin (termasuk `~/.claude` milik `claude` biasa), adalah galat, karena mana pun yang memperbaruinya lebih dulu akan membuat yang lain ter-logout; jalankan `ccam login <name>` untuk memberinya login sendiri. Tidak ada perbaikan.

Sebelum akun-akun, **perubahan grup** yang terhenti di tengah jalan adalah galat di kotak `General`: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` menyelesaikannya sesuai urutan yang direncanakan, dengan jawaban yang sudah diberikan untuk konfliknya; sampai saat itu, tidak ada perubahan grup lain, penggantian nama, atau penghapusan yang dimulai. Akun-akun di dalamnya hanya dibuatkan tautan yang hilang sampai perubahan itu selesai. Jika perubahan itu tidak bisa diselesaikan, galatnya menyebutkan bahwa menghapus file tersebut berarti meninggalkan sisa perubahan. `ccam --dry-run doctor --fix` menampilkan urutannya (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) dan apa yang akan disisihkan oleh bagian setiap akun. Perubahan yang sedang dijalankan ccam lain saat ini adalah informasi, dan selama perubahan itu berjalan, `--fix` tidak memperbaiki tautan pembagian (`its share links are not repaired while another ccam changes the groups; …`). `--fix` memegang `~/.local/state/ccam/group-change.lock` selama memperbaiki, dan membaca `config.toml` lagi setelah menyelesaikan sebuah perubahan.

### Memperbaiki tautan pembagian

Jalur peluncuran langsung membuat tautan yang hilang, tetapi berhenti untuk apa pun yang bisa menyebabkan kehilangan data dan mengarahkan ke sini:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Untuk setiap masalah tautan pembagian, `--fix` memanggil `share.Build` dengan pengaturan pembagian akun sebagai keadaan saat ini sekaligus keadaan tujuan, lalu `Plan.Apply`. `Build` menurunkan ulang setiap langkah dari apa yang ada di disk, rencana yang sama dengan yang dipakai perubahan `ccam group`, sehingga perbaikan mengikuti aturan yang sama seperti saat bergabung ke grup.

<!-- motion:doctor-repair -->

Di sini `settings.json` adalah file asli. Karena grup sudah memiliki `settings.json` yang berbeda, konflik mengambil pilihan default-nya, `KeepGroup`, sehingga salinan milik grup tetap dipakai. File milik akun dipindahkan ke `.ccam-backup/<time>/settings.json` di dalam folder akun, dan `os.Symlink` menautkan path itu ke grup. `--fix` memberitahukannya: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Kasus lainnya:

- Jika grup belum memiliki salinan, file akun dipindahkan ke grup dan menjadi isi pertamanya.
- Untuk folder, file yang belum dimiliki grup dipindahkan ke grup. File yang dimiliki grup persis sama byte demi byte dibuang, dan yang berbeda masuk ke cadangan, disebutkan dalam sebuah catatan (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Folder grup yang hilang dibuat lebih dulu.
- Untuk `history.jsonl`, default-nya adalah `merge by line`: baris-baris milik grup, lalu baris baru Anda. Salinan grup sebelumnya disimpan di `.ccam-backup/<time>/group/`.
- Tautan ke grup lain diganti. Tautan ke folder di luar ccam masuk ke cadangan, dan isinya menjadi isi awal grup jika grup masih kosong.
- Tautan yang tertinggal untuk item yang kini lokal dihapus, dan akun mendapatkan salinan isi grup.

Selain menghapus tautan dan membuang file yang sudah dimiliki grup persis sama byte demi byte, setiap langkah memindahkan data ke grup atau ke folder cadangan. Tidak ada yang hilang.

## Keamanan dan kasus khusus

- `ccam doctor` biasa tidak pernah mengubah apa pun. Gunakan `--dry-run` bersama `--fix` agar pasti aman di skrip.
- Untuk akun yang sedang berjalan, bagaimanapun Claude Code dijalankan, `doctor --fix` hanya membuat tautan yang hilang dan folder grup yang hilang, yang tidak menghilangkan apa pun. Jika masih ada yang tersisa, ia menambahkan `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, karena sesi terus menulis ke file-file yang akan dipindahkan oleh perbaikan.
- Jika perbaikan pembagian gagal, ccam melaporkan `could not repair share links: <error>` untuk akun tersebut dan melewati pemeriksaan entri tak dikenalnya. Akun lain tetap diperiksa.
- Jika perubahan `ccam group` gagal di tengah jalan, grup-grup yang baru sudah tersimpan dan perubahan itu dicatat di `~/.local/state/ccam/group-change.json`; ccam meminta Anda menyelesaikannya dengan `ccam doctor --fix`, yang menjalankan sisanya, termasuk langkah-langkah MCP.
- `.credentials.json` yang berupa symlink, login yang kedaluwarsa, dan info masuk di environment dilaporkan tetapi tidak pernah diubah.
- Penulisan ulang registri plugin ditulis secara atomik tetapi tidak disalin ke `.ccam-backup/`.

## Terkait

- [Berbagi server MCP, pengaturan, dan riwayat](/docs/sharing)
- [Menyalin penyiapan yang sudah ada](/features/migrate)
- [Batas penggunaan dan perpindahan akun](/docs/usage-limits)
- [Perintah](/docs/commands)
