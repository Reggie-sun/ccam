---
title: Gunakan beberapa akun Claude Code di satu komputer
description: Jalankan beberapa akun Claude Code berdampingan dengan CLAUDE_CONFIG_DIR, secara manual lewat alias shell atau otomatis dengan ccam, lalu ganti akun kapan saja.
updated: 2026-10-26
order: 1
translation: machine
sourceHash: "381de7a473b2"
---

Beri setiap akun folder konfigurasinya sendiri. Claude Code membaca folder itu dari variabel lingkungan `CLAUDE_CONFIG_DIR`, dan setiap folder menyimpan login, pengaturan, dan riwayatnya sendiri, sehingga setiap akun tetap login dan Anda bisa menjalankannya bersamaan. Anda bisa menyiapkannya secara manual dengan alias shell, atau membiarkan ccam melakukannya.

## Bagaimana cara menyiapkannya secara manual?

Ini adalah pendekatan yang [dijelaskan dalam dokumentasi Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Tambahkan satu baris untuk setiap akun tambahan ke `~/.bashrc` atau `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Buka terminal baru dan jalankan `claude-work`. Claude Code meminta Anda login, karena foldernya masih baru. Sejak saat itu, `claude` adalah akun pertama Anda dan `claude-work` akun kedua.

Di Windows tanpa WSL, fungsi PowerShell di profil Anda (`notepad $PROFILE`) melakukan hal yang sama:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Apa yang tidak dicakup oleh pendekatan alias?

Cara ini berfungsi baik untuk dua akun. Hal-hal yang perlu diketahui:

- Setiap folder dimulai dalam keadaan kosong: tanpa pengaturan, tanpa `CLAUDE.md`, tanpa server MCP, tanpa agen atau skill. Anda menyalin atau menautkan apa yang Anda perlukan, dan menjaga salinannya tetap sinkron sendiri.
- Alias hanya ada di shell interaktif. Skrip, Makefile, dan program lain tidak melihatnya.
- Server MCP berada di `.claude.json`, bersama dengan login, sehingga tidak bisa begitu saja ditautkan antarfolder.
- Percakapan tetap berada di akun tempat percakapan itu dimulai, sehingga Anda tidak bisa melanjutkannya di akun lain.
- Tidak ada yang menampilkan penggunaan semua akun sekaligus.

## Bisakah saya ganti akun dengan /login saja?

Bisa, jika Anda hanya sesekali berganti akun. Mengetik `/login` di dalam Claude Code [mengganti akun](https://code.claude.com/docs/en/quickstart) pada folder yang sedang Anda pakai. Hanya satu akun yang login pada satu waktu, jadi Anda harus login lagi setiap kali berganti, dan Anda tidak bisa menjalankan dua akun berdampingan. Folder terpisah menghindari kedua masalah itu.

## Bagaimana cara menyiapkannya dengan ccam?

ccam mengotomatiskan ide yang sama dan menambahkan bagian yang kurang:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Sekarang `ccam personal` dan `ccam work` menjalankan kedua akun tersebut, begitu juga `ccwork`, pintasan yang Anda berikan untuk akun kerja. Karena keduanya adalah program sungguhan, bukan alias, keduanya juga berfungsi di skrip. Jalankan `ccam` saja untuk memilih akun dari daftar, dan `ccam use work` untuk menjadikan satu akun sebagai default.

Pengaturan, `CLAUDE.md`, server MCP, dan ekstensi dibagikan sejak awal. `ccam next --continue` melanjutkan percakapan di akun lain. Jika Anda mengaktifkan riwayat percakapan bersama, kedua akun mengerjakan percakapan yang sama, bukan salinannya. `ccam usage` menampilkan penggunaan 5 jam dan mingguan semua akun sekaligus.

Jika Anda sudah menggunakan `~/.claude`, jalankan `ccam import` alih-alih menambahkan akun itu lagi, atau `ccam migrate` untuk menyalin pengaturan, server MCP, dan ekstensinya ke grup bersama `default`.

## Bagaimana saya tahu akun mana yang sedang berjalan?

Di dalam Claude Code, `/status` menampilkan akun yang sedang Anda pakai untuk login. Dengan ccam, jalankan `ccam statusline install` sekali, lalu status line akan selalu menampilkan nama akun dan penggunaannya, misalnya `[work] · Opus · 5h 24% · 7d 41%`.

## Apakah ini berfungsi di Windows dan WSL?

Claude Code berjalan [secara native di Windows dan di dalam WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` berfungsi di keduanya, jadi alias atau fungsi PowerShell di atas sudah cukup di Windows native. ccam juga berjalan secara native di Windows: instal di PowerShell dengan `irm https://ccam.dev/install.ps1 | iex`, atau di dalam WSL dengan perintah Linux. Pintasan menjadi `ccwork.exe`, dan membagikan file seperti `settings.json` memerlukan Mode Pengembang ([Instal ccam](/docs/install)).

## Pemecahan masalah

**Claude Code selalu meminta saya login.** Saat pertama kali dijalankan di folder baru, Claude Code memang selalu meminta login. Jika terus meminta, login Anda sudah kedaluwarsa atau terhapus: jalankan `ccam login work` (atau jalankan alias sekali dan login).

**Server MCP hanya muncul di satu akun.** Server yang ditambahkan di dalam sesi dengan `claude mcp add` menjadi milik akun tersebut, dan tanpa `-s user` hanya untuk folder itu. Pindahkan server ke semua akun dengan `ccam mcp import work <server>`.

**Pengaturan berbeda antarakun.** `ccam share work` menampilkan apa saja yang dibagikan akun itu dan melalui grup mana; `ccam group join default work` membuatnya membagikan apa yang dibagikan `default`, setelah menunjukkan apa yang akan berubah. `ccam doctor --fix` memperbaiki tautan yang rusak dan menyimpan cadangan setiap file yang diubahnya.

**Dua akun membuka percakapan yang sama.** Hindari mengerjakan satu percakapan di dua akun secara bersamaan. ccam memperingatkan Anda jika Anda mencobanya.

## Mana yang sebaiknya dipilih?

Untuk dua akun tanpa ada yang dibagikan, alias sudah cukup. Saat Anda menginginkan pengaturan dan server MCP bersama, lebih banyak akun, penggunaan yang terlihat sekilas, atau melanjutkan percakapan, ccam menghemat pekerjaan administrasinya untuk Anda. [Perbandingan dengan alias](/compare/aliases) mencantumkan perbedaannya secara berdampingan, dan [ccam dan pengalih akun lainnya](/compare/account-switchers) membahas alat yang justru menukar satu login.

Boleh tidaknya Anda menggunakan lebih dari satu akun bergantung pada ketentuan yang berlaku untuk masing-masing akun. Baca ketentuan tersebut sebelum Anda mulai.

**Kursi Team dan Enterprise.** Pemilihan akun otomatis (`ccam next`, `--fallback`, `ccam pick`, tombol `n` di dasbor, dan, setelah Anda mengaktifkannya, [`auto_switch`](/features/next#auto-switch)) juga berganti di antara kursi dalam organisasi Claude Team atau Enterprise yang sama, yang dapat melanggar ketentuan Anthropic. Baca [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) dan tanyakan kepada organisasi Anda dan Anthropic sebelum Anda mengandalkannya.
