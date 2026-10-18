---
title: CLAUDE_CONFIG_DIR nedir?
description: CLAUDE_CONFIG_DIR'ın Claude Code'da ne yaptığı, klasörde neler saklandığı, macOS üzerinde girişin yeri ve hesaplar arasında hangi dosyaların paylaşılabileceği.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR`, Claude Code'a `~/.claude` yerine hangi klasörü kullanacağını söyler. Hesaba özgü her şey orada bulunur: giriş bilgileri, ayarlar, MCP sunucuları, eklentiler ve konuşma geçmişi. Bu değişkeni farklı klasörlere yönlendirmek size ayrı hesaplar sağlar.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Klasörün içinde neler var

Claude Code 2.1.302 ile kontrol edilmiştir:

| Yol | İçerik | Paylaşılabilir |
| --- | --- | --- |
| `.credentials.json` | Giriş tokenları (Linux), MCP OAuth tokenları | Hayır |
| `.claude.json` | Hesap bilgileri, kullanıcı MCP sunucuları, proje bazında durum | Hayır (MCP sunucuları ayrı olarak) |
| `settings.json` | İzinler, hook'lar, model, durum satırı, eklentiler | Evet |
| `CLAUDE.md` | Genel talimatlar | Evet |
| `agents/`, `commands/`, `skills/`, `plugins/` | Uzantılar | Evet |
| `projects/` | `.jsonl` dosyaları olarak konuşmalar, proje bazında bellek | Evet |
| `history.jsonl` | İstem geçmişi | Evet |
| `sessions/`, `shell-snapshots/`, `cache/` ve diğerleri | Çalışma zamanı durumu | Hayır |

## İnsanları şaşırtan iki ayrıntı

**`CLAUDE_CONFIG_DIR` olmadan `.claude.json`, `~/.claude` içinde değildir.** Claude Code bu durumda onu klasörün yanında, `~/.claude.json` konumunda tutar. `CLAUDE_CONFIG_DIR=~/.claude` kullanmaya başlarsanız Claude Code bunun yerine `~/.claude/.claude.json` dosyasını arar ve hesabınızı bulamaz. `ccam import` bu durumu ikisini birbirine bağlayarak çözer.

**macOS üzerinde giriş bilgileri klasörde değil, Anahtar Zinciri'ndedir.** Claude Code bunları `Claude Code-credentials` adı altında saklar; `CLAUDE_CONFIG_DIR` ayarlandığında bu adın ardından klasör yolunun kısa bir hash'i gelir. Bu nedenle ayrı klasörler ayrı Anahtar Zinciri kayıtları alır. Anahtar Zinciri yazmayı reddederse Claude Code, klasördeki `.credentials.json` dosyasına geri döner.

## Klasörler arasında paylaşım

Sembolik bağlantılar çalışır: Claude Code, bağlanmış bir `settings.json` dosyası veya `projects/` klasörü üzerinden okur ve yazar, bağlantıyı da korur. MCP sunucuları bunun istisnasıdır, çünkü giriş bilgileriyle aynı `.claude.json` dosyasını paylaşırlar; bunun yerine onları `--mcp-config` ile iletin.

ccam tam olarak bu kuralların üzerine kuruludur. Klasörleri nasıl bağladığı için [paylaşım](/docs/sharing) sayfasına, genel tablo için de [birden fazla hesap kılavuzuna](/guides/multiple-claude-code-accounts) bakın.
