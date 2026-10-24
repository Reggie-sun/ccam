---
title: ccam dan pengalih akun Claude Code lainnya
description: Perbandingan ccam dengan claude-swap, CCS, clausona, clauth, CAAM, cux, dan alat lain untuk ganti akun Claude Code, serta pendekatan mana yang cocok.
updated: 2026-10-24
order: 2
translation: machine
sourceHash: 6d08203f1ea6
---

Alat untuk berganti akun Claude Code memakai salah satu dari dua pendekatan. Sebagian besar menukar login di dalam satu folder `~/.claude`, sehingga hanya satu akun yang aktif pada satu waktu. Yang lain, termasuk ccam, memberi setiap akun foldernya sendiri dengan `CLAUDE_CONFIG_DIR`, sehingga setiap akun tetap login dan beberapa akun bisa berjalan bersamaan. ccam menambahkan pengaturan dan server MCP bersama, serta bisa melanjutkan percakapan di folder akun lain.

## Menukar satu login atau folder terpisah?

**Menukar login** mempertahankan satu set pengaturan, server MCP, dan riwayat, karena hanya ada satu folder. Pergantiannya cepat, dan sesi yang sedang berjalan bisa memakai login yang baru. Namun, hanya satu akun yang aktif pada satu waktu, untuk semua terminal dan editor sekaligus.

Untuk menukar login, alat itu menyimpan salinan login setiap akun lain dan menuliskannya ke tempatnya saat dibutuhkan. Claude Code mengganti refresh token di dalam login setiap kali memperbaruinya, sehingga setiap salinan yang tersimpan harus terus diselaraskan dengan sesi yang sedang berjalan. Akun yang ter-logout atau perlu `/login` lagi setelah berganti termasuk masalah yang paling sering diperbaiki oleh proyek-proyek ini.

**Folder terpisah** adalah metode yang [dijelaskan dalam dokumentasi Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Setiap akun tetap login, dan Anda bisa menjalankan sesi kerja dan sesi pribadi berdampingan. Setiap login tetap berada di foldernya sendiri, tempat Claude Code memperbaruinya seperti biasa. Konsekuensinya, setiap folder dimulai dengan pengaturan, server MCP, dan riwayatnya sendiri, kecuali ada alat yang membagikannya.

## Perbandingan alat-alatnya

Diambil dari dokumentasi dan kode sumber masing-masing proyek pada 2026-10-20 (clausona: v0.12.0-beta). "Tidak didokumentasikan" berarti proyek tersebut tidak menyebutkannya, bukan berarti fiturnya tidak ada.

| Alat | Pendekatan | Beberapa akun sekaligus | Penggunaan 5 jam dan mingguan | Saat batas tercapai | Platform |
| --- | --- | --- | --- | --- | --- |
| ccam | Folder terpisah | Ya | Ya, `ccam usage` | Notifikasi di dalam sesi; `ccam next --auto --continue` atau `/ccam:switch` membawa percakapan ke akun dengan sisa kuota terbanyak, dengan atau tanpa riwayat bersama; `ccam --fallback` menjalankan ulang tugas `-p` di sana dengan sendirinya, dan `ccam pick` menyebar worker ke berbagai akun | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Menukar login; folder terpisah dengan `cswap run` | Ya, dengan `cswap run` | Ya | `cswap auto` menukar login pada ambang tertentu; sesi yang berjalan tetap berlanjut | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Folder terpisah | Ya | Sebagian | Kumpulan akun meneruskan permintaan ke akun berikutnya melalui proxy lokal | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Folder terpisah | Ya | Ya | Rute memilih akun yang masih punya kuota sebelum setiap proses, bisa menunggu reset, dan membagikan akun yang berbeda ke setiap worker | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Menukar login; folder terpisah dengan `clauth start` | Ya, dengan `clauth start` | Ya | Berpindah mengikuti rantai cadangan | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Menukar login, atau profil terisolasi | Ya, dengan profil terisolasi | Ya | `caam run` mengulang perintah non-interaktif di akun lain | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Menukar login dan membungkus `claude` | Tidak, satu akun aktif | Ya | Menukar login dan melanjutkan percakapan yang sama | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Menukar login | Tidak | Ya | Peralihan opsional pada ambang tertentu | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Menukar login | Tidak | Ya | Peralihan opsional pada 100% | macOS |

ccswap adalah fork dari claude-swap yang juga mengelola login Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile), dan [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) juga memberi setiap akun foldernya sendiri. Beberapa alat ini juga menangani agen coding lain, seperti Codex atau Gemini.

## Cara alat-alat ini dibuat

| Alat | Ditulis dengan | Selama Claude Code berjalan | Menyalin login antartempat | Yang bisa dibagikan antarakun |
| --- | --- | --- | --- | --- |
| ccam | Go, satu biner statis | Tidak ada: ccam menggantikan dirinya dengan `claude` (di Windows, yang tidak bisa melakukannya, ccam menunggunya; relay opsional tetap berjalan di sampingnya) | Tidak | Grup: masing-masing menyebutkan item yang dibagikannya dan akun-akunnya; satu akun bisa masuk ke beberapa grup |
| claude-swap | Python 3.12+, dipasang dengan uv atau pipx | Tidak ada dengan `cswap run`; `cswap auto` dan aplikasi bilah menu tetap berjalan | Ya, saat menukar login | Dengan `cswap run`: pengaturan, `CLAUDE.md`, dan ekstensi sekaligus, riwayat opsional; server MCP disalin setiap kali dijalankan |
| CCS | Node.js 18+ atau Bun | CCS tetap sebagai proses induk; kumpulan akun menambahkan proxy lokal | Tidak untuk folder akunnya; kumpulan akun menyimpan login sendiri | Bersama atau lokal per akun, riwayat dalam grup konteks; server MCP disalin |
| clausona | TypeScript, Node.js 20+ | `claude` biasa berjalan melalui fungsi shell; `clausona run` mempertahankan Node sebagai proses induk; proses Node di latar belakang setelah setiap keluar | Tidak, tetapi ia memperbarui sendiri token akun yang menganggur untuk membaca penggunaannya | Semuanya kecuali file per akun, dari satu akun utama; riwayat opsional; server yang ditambahkan dengan `claude mcp add` tidak dibagikan |
| clauth | Rust | clauth tetap sebagai proses induk | Ya, saat menukar login | Seluruh `~/.claude`, atau tidak sama sekali |
| CAAM | Go | Tergantung mode; daemon opsional | Ya, dalam mode vault-nya | Ditentukan oleh mode; server MCP disalin setiap kali dijalankan |
| cux | Go | Sebuah wrapper tetap membungkus `claude` | Ya | Semuanya: satu folder |
| CCSwitcher, Claude Switcher | Aplikasi macOS (Swift, Python) | Aplikasi bilah menu | Ya | Semuanya: satu folder |

## Keunggulan ccam

### Login tetap di tempat Claude Code menaruhnya

ccam tidak pernah menyalin login Claude dari satu tempat ke tempat lain. Login setiap akun tetap berada di foldernya sendiri, atau di macOS dalam item Rantai Kunci (Keychain) milik folder itu sendiri, dan Claude Code-lah program yang memperbaruinya. Tidak ada salinan tersimpan yang bisa menjadi tidak selaras. `~/.claude` lama yang didaftarkan dengan `ccam import` tetap memakai item Rantai Kunci yang dipakai `claude` biasa, sehingga tetap login dengan cara mana pun.

ccam memperbarui token kedaluwarsa sendiri hanya saat Anda menjalankan `ccam usage --refresh`. ccam mengambil kunci (lock) terlebih dahulu dan menolak jika akun itu sedang berjalan, bagaimanapun Claude Code dijalankan (termasuk sesi latar belakang dan daemon-nya), dan selama Claude Code sendiri sedang memperbarui login itu, sehingga tidak pernah merotasi token di bawah sesi yang aktif.

### Tidak ada yang berjalan di antara Anda dan Claude Code

ccam membaca satu file, memeriksa tautannya, lalu menggantikan dirinya dengan `claude` (`exec`). Selama Anda bekerja, tidak ada proses ccam, wrapper, proxy, atau daemon: input standar, sinyal, dan kode keluar milik Claude Code, dan permintaannya langsung menuju Anthropic. Windows tidak bisa mengganti sebuah proses, jadi di sana ccam menjalankan Claude Code dengan konsol yang sama dan menunggu kode keluarnya, tanpa melakukan hal lain.

Menjalankan akun dengan cara ini menambah sekitar 4 ms, di bawah 6 ms pada persentil ke-95, dengan memori puncak 8 MB ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) dengan hyperfine, 300 kali jalan). Di laptop yang sama, memulai interpreter Python saja butuh 13 ms dan Node.js 31 ms, sebelum program melakukan apa pun.

### Grup menentukan apa yang dibagikan

Server MCP, `settings.json`, `CLAUDE.md`, ekstensi (agen, perintah, skill, plugin), dan riwayat percakapan adalah lima item terpisah. Sebuah grup menyebutkan item yang dibagikannya dan akun-akun di dalamnya, sehingga dua akun kerja bisa berbagi riwayat dalam grup `work` sementara akun pribadi menyimpan riwayatnya sendiri, dan ketiganya berbagi server MCP yang sama melalui `default`. `ccam group set work --add memory` memindahkan setiap akun dalam grup itu dalam satu rencana; [dasbor](/features/dashboard) melakukan hal yang sama dengan kotak centang dan meredupkan item yang sudah didapat sebuah akun dari grup lain.

- Server MCP diberikan dengan `--mcp-config` saat dijalankan, jadi secara bawaan peluncuran tidak menyentuh `.claude.json`, file yang ditulis ulang Claude Code selama berjalan. ccam hanya menulis ke sana pada langkah yang menyatakannya, di bawah kunci milik Claude Code dan dengan cadangan: untuk Remote Control setelah bertanya, saat sebuah akun keluar dari atau bergabung kembali ke berbagi MCP, dan untuk `ccam mcp import`. `ccam mcp add` menambahkan server sekali untuk seluruh grup.
- Item lainnya berupa tautan ke folder grup, sehingga perubahan di satu akun langsung terlihat di akun lain. Tidak ada yang disalin saat dijalankan.
- Login, `.claude.json`, dan status yang sedang berjalan tidak pernah dibagikan.

### Setiap perubahan ditampilkan dulu dan dicadangkan

`ccam group` dan `ccam migrate` menampilkan rencana dan bertanya sebelum memindahkan apa pun. Jika sebuah file berbeda antara akun dan grupnya, Anda yang memilih: pertahankan milik grup, pakai milik Anda, atau gabungkan JSON-nya; di dalam folder, salinan milik grup dipertahankan dan rencana menyebutkan file yang berbeda. File asli masuk ke `.ccam-backup/<time>/`, akun yang dihapus di-logout dan foldernya masuk ke folder sampah, dan file diganti dengan satu kali rename, sehingga penulisan yang terputus tetap meninggalkan file lama. `--dry-run` menunjukkan apa yang akan dilakukan oleh peluncuran, berbagi, migrasi, pintasan, pembaruan, dan perbaikan. `ccam doctor --fix` memperbaiki tautan dan izin yang rusak.

### Percakapan yang sama di akun berikutnya

`ccam next --auto --continue` menanyakan penggunaan semua akun sekaligus, membuka akun dengan sisa kuota terbanyak, lalu melanjutkan percakapan. Dengan riwayat bersama, percakapannya tetap sama; tanpa riwayat bersama, Claude Code menyalinnya ke akun baru (`--resume <file> --fork-session`), dan ccam memberi tahu berapa token yang dibaca ulang oleh balasan pertama. Saat batas penggunaan menghentikan sebuah sesi, hook dari plugin ccam menampilkan notifikasi, dan `/ccam:switch` meneruskan percakapan itu. Anda yang menentukan kapan berganti; tidak ada yang memantau penggunaan Anda di latar belakang, dan pilihan otomatis tetap berada dalam satu organisasi, sehingga percakapan perusahaan tidak akan tanpa sengaja berakhir di akun pribadi. Baris status menyarankan perintah itu begitu sebuah akun mencapai 90%.

Skrip tidak bisa menunggu Anda, jadi `ccam --fallback work -p "…"` berganti akun untuknya: saat akun mencapai batasnya, perintah yang sama dijalankan lagi sebagai akun dengan sisa kuota terbanyak, dan hanya output dari proses yang selesai yang dicetak. Kegagalan yang bukan karena batas penggunaan diteruskan apa adanya. Untuk worker paralel, `ccam pick` memberi setiap worker akun berbeda yang masih punya kuota, dan keduanya keluar dengan kode 75 jika tidak ada akun yang tersisa kuotanya.

**Kursi Team dan Enterprise.** Pilihan otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama. Pusat bantuan Anthropic menyatakan "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (setiap kursi pada paket Team diberikan kepada satu orang tertentu, dan hanya orang itu yang boleh menggunakannya) dan "Users can't route their work through another user's account to use that user's remaining usage" (pengguna tidak boleh menyalurkan pekerjaannya melalui akun pengguna lain untuk memakai sisa kuota pengguna tersebut) ([paket Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); begitu pula untuk [paket Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Karena itu, memakai pilihan otomatis di antara kursi dapat melanggar ketentuan Anthropic. Tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya. Paket pribadi (Pro, Max) dibahas dalam [catatan tentang ketentuan di FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

### Akun yang tepat di setiap folder

`ccam folder work ~/code/company` membuat setiap peluncuran tanpa nama akun di dalam folder itu memakai `work`, sementara proyek Anda yang lain tetap memakai akun bawaan. Pengaturan ini disimpan di `config.toml` yang sama yang dibaca setiap peluncuran, sehingga tidak menambah waktu.

### Satu biner, unduhan yang diperiksa

ccam adalah biner statis sekitar 10 MB tanpa runtime yang perlu dipasang, untuk Linux, macOS, dan Windows. Skrip instalasi (`install.sh`, `install.ps1`) dan `ccam update` memeriksa setiap unduhan terhadap checksum SHA-256 rilis, dan ccam tidak pernah memperbarui diri tanpa bertanya. Tidak ada telemetri: satu-satunya permintaan yang dibuatnya adalah pengecekan penggunaan akun Anda ke Anthropic, secara paralel dengan batas 3 detik dan cache 60 detik, pertanyaan dasbor ke GitHub tentang rilis terbaru saat dasbor dibuka, serta unduhan yang Anda minta. Dengan `usage_source = "claude"`, permintaan itu pun melewati Claude Code, dan ccam tidak pernah menyentuh token sama sekali.

## Di mana alat lain lebih unggul

- claude-swap, clauth, cux, CAAM, dan CCSwitcher bisa mengganti akun sesi interaktif yang sedang berjalan dengan sendirinya saat sebuah akun melewati ambang penggunaan. ccam mengganti akun untuk proses `-p` dengan `--fallback` dan menyerahkan sesi interaktif kepada Anda, dengan notifikasi saat batas tercapai dan `/ccam:switch`.
- Di Windows tanpa Mode Pengembang (Developer Mode), ccam bisa membagikan folder tetapi tidak bisa membagikan file tunggal seperti `settings.json` (maka `ccam add` pertama hanya membagikan server MCP dan ekstensi); claude-swap dan CCS menyalinnya ke sana sebagai gantinya.
- CCSwitcher, Claude Switcher, dan claude-swap menampilkan penggunaan di bilah menu macOS.
- CCS, clausona, CAAM, clauth, ccswap, dan Claude Switcher juga mengelola Codex dan agen coding lainnya.
- clausona juga menyalakan dan mematikan ekstensi per akun (`ccam ext` milik ccam hanya menampilkannya), mencatat biaya token secara lokal, dan menyertakan plugin Claude Code yang menjalankan sekumpulan worker di tmux dan terminal lain.
- CAAM, CCSwitcher, dan Claude Switcher bisa dipasang dengan Homebrew, dan clauth serta CAAM menandatangani checksum rilisnya. ccam dipasang dengan skripnya, `.deb`, atau `go install`, dan memeriksa checksum tanpa tanda tangan.

## Mana yang cocok?

Pilih alat yang menukar login jika Anda memakai satu akun pada satu waktu atau ingin sesi yang sedang berjalan berganti akun secara otomatis. Pilih ccam jika Anda ingin akun berjalan berdampingan, satu penyiapan bersama, percakapan yang sama di akun berikutnya, dan tidak ada yang berjalan di samping Claude Code, di Linux, macOS, atau Windows. Jika tidak ada yang perlu dibagikan, [alias shell biasa](/compare/aliases) sudah cukup.

Boleh tidaknya Anda menggunakan beberapa akun bergantung pada ketentuan setiap akun. Periksa ketentuan itu sebelum melakukannya.
