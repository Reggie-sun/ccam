---
title: Pertanyaan yang sering diajukan
metaTitle: "FAQ ccam: beberapa akun Claude Code"
description: Jawaban seputar ccam, akun Claude Code, batas penggunaan, pembagian, privasi, dan sistem yang didukung.
updated: 2026-10-24
order: 6
translation: machine
sourceHash: c4e3b89bb4ea
---

## Apakah ccam alat resmi dari Anthropic?

Tidak. ccam adalah proyek open source tidak resmi dan tidak berafiliasi dengan atau didukung oleh Anthropic. ccam menggunakan pengaturan `CLAUDE_CONFIG_DIR` yang [didokumentasikan Claude Code](https://code.claude.com/docs/en/env-vars) untuk menjalankan beberapa akun secara berdampingan.

## Apakah saya boleh menggunakan beberapa akun Claude Code?

Itu bergantung pada ketentuan setiap akun dan, untuk akun kerja, kebijakan organisasi Anda. Periksa hal tersebut sebelum Anda menggunakan beberapa akun. ccam hanya mempermudah perpindahan di antara akun yang sudah Anda miliki.

**Kursi Team dan Enterprise.** Pilihan otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama. Pusat bantuan Anthropic menyatakan "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (setiap kursi pada paket Team diberikan kepada satu orang tertentu, dan hanya orang itu yang boleh menggunakannya) dan "Users can't route their work through another user's account to use that user's remaining usage" (pengguna tidak boleh menyalurkan pekerjaannya melalui akun pengguna lain untuk memakai sisa kuota pengguna tersebut) ([paket Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); begitu pula untuk [paket Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Karena itu, memakai pilihan otomatis di antara kursi dapat melanggar ketentuan Anthropic. Tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya. Paket pribadi (Pro, Max) dibahas dalam paragraf di atas.

## Apakah ccam mengakali batas penggunaan?

Tidak. Setiap akun tetap memiliki batasnya sendiri. ccam menampilkan batas tersebut dan memungkinkan Anda melanjutkan di akun lain milik Anda saat satu akun habis kuotanya.

## Bagaimana cara menggunakan dua akun Claude Code di satu komputer?

Beri setiap akun folder konfigurasinya sendiri dengan `CLAUDE_CONFIG_DIR`. ccam melakukannya untuk Anda: jalankan `ccam add` untuk setiap akun, lalu jalankan akun-akun itu dengan `ccam personal` dan `ccam work`, atau dengan pintasan yang namanya Anda pilih sendiri, seperti `ccwork`. [Panduan ini](/guides/multiple-claude-code-accounts) juga menunjukkan cara melakukannya secara manual.

## Bagaimana cara melanjutkan percakapan di akun lain?

Jalankan `ccam next --continue` di folder itu: perintah ini membawa percakapan terakhir folder tersebut ke akun berikutnya (`--auto` untuk akun dengan sisa kuota terbanyak). Di dalam Claude Code, `/ccam:switch` dari [plugin ccam](/docs/install#claude-code-plugin) melakukan hal yang sama untuk sesi yang sedang Anda jalani. Jika kedua akun berbagi riwayat percakapan (`ccam group set default --add history`, atau dengan grup tersendiri: `ccam group add pair --items history --members work,side`), akun baru membuka percakapan yang sama persis; jika tidak, Claude Code menyalinnya ke riwayat akun baru.

## Apakah ccam melihat atau mengirim login saya?

ccam tidak pernah membagikan atau menyalin login Claude Anda antarakun, dan tidak pernah menyalin info masuk server MCP. ccam membaca token akses sebuah akun hanya untuk menanyakan penggunaan akun itu ke endpoint penggunaan Anthropic, saat Anda menjalankan `ccam usage`, `ccam next --auto`, atau `ccam doctor`, dan memperbarui token yang kedaluwarsa hanya dengan `ccam usage --refresh`; dengan `usage_source = "claude"`, ccam sama sekali tidak membaca token. `ccam ext` dan `ccam doctor` membaca nama-nama info masuk MCP sebuah akun, tidak pernah tokennya. ccam tidak memiliki telemetri. Selain pengecekan itu, satu-satunya permintaan yang dibuatnya sendiri adalah pertanyaan dasbor ke GitHub tentang rilis ccam terbaru, yang tidak membawa apa pun tentang Anda atau akun-akun Anda.

## Bagaimana cara masuk di server jarak jauh?

Jalankan `ccam add` atau `ccam login <name>` di sesi SSH seperti biasa. Claude Code menunggu proses masuk di server, sehingga browser di komputer Anda berakhir di halaman `http://localhost:<port>/callback?code=…` yang tidak bisa dimuat. Salin alamat itu dari browser lalu tempelkan ke terminal: ccam meneruskannya ke Claude Code di server, yang kemudian menyelesaikan proses masuk. Alamat itu berisi kode sekali pakai, jadi jangan bagikan kepada siapa pun.

## Di mana akun saya disimpan?

Di `~/.local/share/ccam/profiles/<name>/`, satu folder per akun. Di macOS, Claude Code menyimpan login itu sendiri di Rantai Kunci (Keychain), dengan nama yang diturunkan dari folder tersebut, sehingga `ccam rename` mempertahankan nama lama folder itu di sana.

## Apa yang terjadi pada ~/.claude saya yang sudah ada?

Tidak ada, kecuali Anda menjalankan `ccam import` atau `ccam migrate`. Impor mendaftarkan folder itu di tempatnya, sehingga `claude` biasa dan `ccam <name>` tetap menggunakan akun dan pengaturan yang sama. `ccam migrate` menyalin pengaturan, `CLAUDE.md`, ekstensi, dan server MCP dari folder itu ke pengaturan bersama dan hanya membaca folder tersebut.

## Sistem apa saja yang didukung ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine, dan lainnya), macOS, serta Windows 10 dan 11, secara native atau di dalam WSL. Di Windows, membagikan file seperti `settings.json` memerlukan Mode Pengembang (Developer Mode); folder bersama berfungsi tanpanya, dan tanpanya `ccam add` pertama membuat grup `default` yang hanya berisi server MCP dan ekstensi ([berbagi di Windows](/features/share#on-windows)).

## Mengapa server MCP hanya muncul di satu akun?

Server yang ditambahkan di dalam sesi dengan `claude mcp add` hanya disimpan untuk akun tersebut, dan tanpa `-s user` hanya untuk folder saat ini. ccam memberi tahu Anda saat berikutnya Anda menjalankan akun itu, dan `ccam mcp list` serta `ccam ext` menampilkannya; pindahkan server ke grup bersama dengan `ccam mcp import <account> <server>`.

## Apakah ccam memperlambat Claude Code?

Menjalankan Claude Code melalui ccam menambah sekitar 4 milidetik, diukur dengan `scripts/perf-check.sh` di sebuah laptop. ccam membaca satu file kecil, memeriksa tautan bersama, lalu menggantikan dirinya dengan `claude`, sehingga ccam tidak berjalan saat Anda bekerja. Windows tidak bisa mengganti sebuah proses, jadi di sana ccam menjalankan Claude Code dan menunggunya, tanpa melakukan hal lain.

## Bagaimana cara memperbarui ccam?

Jalankan `ccam update`. Perintah ini mengunduh rilis terbaru, memeriksa checksum-nya, lalu mengganti file `ccam`; `ccam update --check` hanya memberi tahu apakah ada rilis yang lebih baru. Dasbor juga memberi tahu jika ada rilis yang lebih baru, dan `U` menginstalnya setelah bertanya; ccam tidak pernah memperbarui dirinya sendiri. Jika Anda menginstal paket `.deb`, instal `.deb` yang baru dengan apt. Lihat [pembaruan](/docs/install#update).

## Bagaimana cara menghapus instalasi ccam?

Hapus pintasan setiap akun dengan `ccam shortcut <account> --remove`. Jika Anda mengimpor `~/.claude` dan folder itu ada di sebuah grup, jalankan `ccam remove <name>` untuknya terlebih dahulu, agar folder itu mendapatkan kembali salinan dari apa yang dibagikannya. Lalu hapus `ccam` dari `~/.local/bin`, dan hapus `~/.config/ccam`, `~/.local/share/ccam`, dan `~/.local/state/ccam` jika Anda tidak lagi memerlukan akun-akun tersebut. Lihat [instalasi](/docs/install#uninstall).
