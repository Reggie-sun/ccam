---
title: Komut başvurusu
description: ccam add ve ccam <name> komutlarından ccam usage, ccam next ve ccam doctor komutlarına kadar tüm ccam komutları ve seçenekleri.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

Bir komutun ne yapacağını, yapmadan görmek için komuttan önce `--dry-run` yazın, örneğin `ccam --dry-run work`.

<!-- commands -->

## Claude Code'u başlatma

`ccam <name|number> [arguments]`, Claude Code'u o hesapla başlatır. ccam, addan sonraki argümanları okumaz; bunları aynı sırayla `claude` komutuna iletir. Standart girdi, çıkış kodu ve sinyaller de tam olarak `claude` komutunun kendisinde olduğu gibi davranır. Addan sonra `--` yazmak isteğe bağlıdır.

Hesap adı olmadan argümanlarla çalıştırılan `ccam`, örneğin `ccam -p "summarise"`, varsayılan hesabı kullanır. Varsayılan hesabı `ccam use <name>` ile ayarlayın.

`claude1`, `claude2` ve devamındaki kısayollar, `~/.local/bin` içindeki ccam'a yönelen bağlantılardır. `ccam shims install --names` ayrıca her hesap için `claude-<name>` oluşturur.

## Hesap ekleme ve kaldırma

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

Klasör belirtilmeden çalıştırılan `ccam import`, `~/.claude` ve `~/.claude-*` klasörlerini arar ve her biri için sorar. İçe aktarılan klasörler bulundukları yerde kaydedilir; ccam hiçbir zaman bir girişi taşımaz veya kopyalamaz.

## Paylaşım

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Her öğe `on` (`default` grubu), `off` veya bir grup adı alır. [Paylaşım](/docs/sharing) sayfasına bakın.

## MCP sunucuları

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Kullanım ve geçiş

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Bakım

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Dosyalar

| Yol | İçerik |
| --- | --- |
| `~/.config/ccam/config.toml` | Hesaplar ve neyi paylaştıkları |
| `~/.local/share/ccam/profiles/<name>/` | Her hesabın `CLAUDE_CONFIG_DIR` klasörü |
| `~/.local/share/ccam/groups/<group>/` | Paylaşılan ayarlar, MCP sunucuları, uzantılar ve geçmiş |
| `~/.local/state/ccam/` | Kullanım önbelleği ve son kullanım zamanları |
| `~/.local/bin/claude1`, `claude2`… | ccam'a yönelen bağlantılar |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ve `XDG_STATE_HOME` değişkenleri bu klasörlerin yerini değiştirir.
