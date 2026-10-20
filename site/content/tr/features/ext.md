---
title: Uzantıları inceleme
description: ccam ext her hesabın yüklediği beceri, komut, eklenti, MCP sunucusu ve hook'ları, her birinin nereden geldiğini ve neyin bozuk olduğunu gösterir.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` şu soruyu yanıtlar: "Bu beceri (veya MCP sunucusu) neden bir hesapta çalışıyor da diğerinde çalışmıyor?". Her hesap için Claude Code'un yüklediği uzantıları ve her birinin nereden geldiğini listeler: bir paylaşım grubu, hesabın kendi klasörü ya da içinde bulunduğunuz proje. Tablonun altında, bir uzantının beklediğiniz gibi yüklenmesini engelleyen şeyleri listeler. Yalnızca okur.

## Kullanım

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

Bir hücre, o öğenin o hesaptaki kaynağını adlandırır: paylaşım grubu (`default`), öğe hesabın kendi klasöründeyse `account`, hesapta yoksa `-`. `!` sorunlu bir öğeyi işaretler; sorun tablonun altında açıklanır.

## Neleri okur

| Tür | Nerede | Gösterilen kaynak |
| --- | --- | --- |
| Beceriler | `skills/` içindeki, `SKILL.md` dosyası olan her klasör | `skills/` bağlantısının gösterdiği grup ya da `account` |
| Komutlar, ajanlar | `commands/` ve `agents/` altındaki `.md` dosyaları, alt klasörler `folder/name` olarak | grup ya da `account` |
| Eklentiler | `plugins/installed_plugins.json` | grup ya da `account` |
| MCP sunucuları | grubun `mcp.json` dosyası (MCP paylaşılıyorsa) ve hesabın `.claude.json` dosyası | grup ya da `account` |
| Hook'lar | `hooks` (`settings.json` içinde), olay başına sayılarak | grup ya da `account` |
| Bu proje | geçerli klasördeki `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` ve `.mcp.json` | tüm hesaplar için bir kez listelenir |

## Bildirdiği sorunlar

- **Bozuk bağlantılar**: `skills/`, `commands/` veya `agents/` içinde hedefi kaybolmuş bir bağlantı ya da var olmayan bir yeri gösteren paylaşılan bir klasör (`skills`, `commands`, `agents`, `plugins`) (bunu `ccam doctor --fix` ile onarın).
- **`SKILL.md` olmayan bir beceri klasörü**: Claude Code onu atlar.
- **İki yerde farklı tanımlanmış bir MCP sunucusu**: grubun `mcp.json` dosyası ile hesabın `.claude.json` dosyası ya da bunlardan biri ile projenin `.mcp.json` dosyası. İki yerde de aynı ayarların olması (`mcp_strategy = "sync"` bunu yazar) sorun değildir. Bir hesabın kendi sunucusunu `ccam mcp import <account> <server>` ile gruba taşıyın.
- **Geçerli JSON olmayan bir `installed_plugins.json`.**

## Güvenlik ve uç durumlar

- `ccam ext` hiçbir şeyi değiştirmez ve ağ isteği yapmaz. Sunucu adlarını ve hook sayılarını okur; MCP tanımlarındaki `env` veya başlık değerlerini asla okumaz.
- Diskte olanı gösterir. Kurulu ama ayarlarda kapatılmış bir eklenti ya da yönetilen ayarlarla engellenmiş bir sunucu yine görünür.
- O projenin `.claude` klasörünü ve `.mcp.json` dosyasını da dahil etmek için komutu proje klasöründe çalıştırın.

## İlgili

- [Hesaplar arasında paylaşım](/features/share)
- [Paylaşılan MCP sunucuları](/features/mcp)
- [Kontrol etme ve onarma](/features/doctor)
