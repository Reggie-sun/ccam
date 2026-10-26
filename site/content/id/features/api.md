---
title: Akun kunci API
description: ccam add --api menambahkan akun kunci API untuk API Anthropic atau gateway; kunci dibaca saat dijalankan dari variabel atau perintah, tak pernah disimpan.
updated: 2026-10-26
order: 13
translation: machine
sourceHash: "771d749f1f90"
---

Akun kunci API menjalankan Claude Code dengan kunci API alih-alih login langganan Claude. Kunci itu bisa untuk API Anthropic atau untuk endpoint apa pun yang memakai format Messages yang sama, seperti gateway LiteLLM atau OpenRouter. ccam tidak menyimpan kunci itu. ccam membacanya setiap kali akun dijalankan, dari variabel lingkungan atau dari perintah seperti pengelola kata sandi, dan hanya memberikannya ke proses Claude Code tersebut.

## Penggunaan

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

Akun kunci API bergabung ke grup seperti akun lainnya (secara bawaan grup `default`, yang membagikan pengaturan, `CLAUDE.md`, ekstensi, dan server MCP), dan berfungsi dengan pintasan, `ccam folder`, dan `ccam ext`.

## Opsi

| Opsi | Fungsinya | Default |
| --- | --- | --- |
| `--api` | Menambahkan akun kunci API alih-alih login | nonaktif |
| `--key-from env:NAME` | Membaca kunci dari variabel lingkungan `NAME` | |
| `--key-from command:COMMAND` | Menjalankan `COMMAND` (dengan `sh -c`, atau di Windows seperti yang diketik di `cmd /d /s /c "…"`, dengan tanda kutip dipertahankan) dan memakai baris pertama yang dicetaknya | |
| `--base-url <url>` | Endpoint (`ANTHROPIC_BASE_URL`) | API Anthropic |
| `--auth api-key\|bearer` | Mengirim kunci sebagai `ANTHROPIC_API_KEY` (`x-api-key`) atau sebagai `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, yang biasanya diminta gateway) | `api-key` |
| `--model <model>` | Model (`ANTHROPIC_MODEL`) | Bawaan Claude Code |

Pengaturan ini disimpan di `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Cara kerjanya

Setiap kali akun dijalankan, setelah pemeriksaan biasa dan sebelum Claude Code dimulai:

1. ccam membaca kunci: `env:` dari lingkungan, `command:` dengan menjalankan perintah itu dengan batas waktu 30 detik. Perintah bisa meminta kata sandi di terminal; output galatnya juga masuk ke terminal. Kunci yang kosong menghentikan peluncuran dengan pesan yang menyebutkan di mana kunci itu diharapkan.
2. ccam menghapus setiap variabel warisan yang menentukan cara Claude Code masuk atau penyedia mana yang dipakainya:
   - kunci dan token: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`;
   - endpoint: `ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`;
   - pengalih penyedia: `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`.

   Kunci yang diatur di shell Anda untuk keperluan lain tidak bisa sampai ke akun ini, dan login langganan tidak bisa didahulukan di atas kuncinya.
3. ccam mengatur `ANTHROPIC_BASE_URL` dan `ANTHROPIC_MODEL` jika diberikan, dan kunci sebagai `ANTHROPIC_API_KEY` atau `ANTHROPIC_AUTH_TOKEN`.

Akun dengan login langganan tetap memakai lingkungan yang diwarisinya, seperti sebelumnya. Dalam sesi interaktif, Claude Code mungkin bertanya sekali apakah akan memakai kunci yang ditemukannya di lingkungan.

## Keamanan dan kasus khusus

- Kunci tidak pernah ditulis ke `config.toml`, cache, atau log, dan `ccam --dry-run gw` tidak membacanya: sumber `command:` tidak dijalankan oleh dry run. Sebagai gantinya, dry run menampilkan endpoint, model, dan asal kuncinya:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  Variabel warisan yang dihapus oleh peluncuran ditampilkan sebagai `env -u <name>` di awal. Dengan sumber `env:` yang variabelnya kosong, dry run gagal seperti halnya peluncuran.
- `ccam add --api` langsung memeriksa sumber `env:` dan memperingatkan jika variabelnya kosong; sumber `command:` baru dijalankan saat peluncuran.
- `ccam login` dan `ccam logout` menolak akun kunci API, karena akun itu tidak punya login. `ccam doctor` menampilkan endpoint dan sumber kuncinya, dan memperingatkan jika variabel `env:` kosong.
- `ccam usage` menampilkan `no limits (API)` untuknya. `ccam next --auto`, `--fallback`, dan `ccam pick` hanya memilih di antara akun langganan, karena kunci API tidak punya batas 5 jam atau mingguan untuk dibandingkan.
- `settings.json` bersama dengan blok `env` atau `apiKeyHelper` berlaku untuk setiap akun dalam grup, termasuk akun kunci API. Jauhkan pengaturan penyedia dari pengaturan bersama, atau masukkan akun kunci API ke grup tanpa `settings`.

## Terkait

- [Mengelola akun](/features/accounts)
- [Menjalankan Claude Code](/features/launch)
- [Referensi perintah](/docs/commands)
