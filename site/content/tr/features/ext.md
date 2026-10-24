---
title: Uzantıları inceleme
description: ccam ext her hesabın yüklediği beceri, komut, eklenti, MCP sunucusu ve hook'ları, her birinin nereden geldiğini ve neyin bozuk olduğunu gösterir.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: "45740e4efbef"
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
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

Bir hücre, o öğenin o hesaptaki kaynağını adlandırır: paylaşım grubu (`default`), öğe hesabın kendi klasöründeyse `account`, hesabın yalnızca bu klasör için tuttuğu bir MCP sunucusu için `local`, hesapta yoksa `-`. Kurulu ama açık olmayan bir eklenti, kaynağından sonra `(off)` gösterir. `!` sorunlu bir öğeyi işaretler; sorun tablonun altında açıklanır. Yalnızca düz `claude` komutunda bulunan sunucular gibi hiçbir hesaba ait olmayan sorunlar, ad olmadan en sona gelir. `ccam ext --json` aynı raporu (`accounts`, proje, `problems`) yazdırır; Claude Code'un ağ üzerinden ulaştığı MCP sunucularında `remote` bulunur; hiç hesap yoksa raporu `accounts: []` ile yazdırır.

## Neleri okur

| Tür | Nerede | Gösterilen kaynak |
| --- | --- | --- |
| Beceriler | `skills/` içindeki, `SKILL.md` dosyası olan her klasör | `skills/` bağlantısının gösterdiği grup ya da `account` |
| Komutlar, ajanlar | `commands/` ve `agents/` altındaki `.md` dosyaları, alt klasörler `folder/name` olarak | grup ya da `account` |
| Eklentiler | `plugins/installed_plugins.json`; hesabın ve projenin ayarlarındaki `enabledPlugins` ile açılır veya kapatılır; yalnızca tek bir proje için kurulmuş bir eklenti yalnızca o projede görünür | grup ya da `account`; açık değilse `(off)` |
| MCP sunucuları | grubun `mcp.json` dosyası (MCP paylaşılıyorsa), hesabın `.claude.json` dosyası (kendi sunucuları ve yalnızca bu klasör için tuttukları) ve eklentilerinin sunucuları, `plugin:<plugin>:<server>` olarak | grup, `account` ya da `local` |
| Hook'lar | `hooks` (`settings.json` içinde), olay başına sayılarak | grup ya da `account` |
| Bu proje | geçerli klasördeki `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` ve `.mcp.json` | bir kez listelenir; klasöre güvenmemiş bir hesap, `.mcp.json` sunucuları için bir sorun alır |

## Bildirdiği sorunlar

- **Bozuk bağlantılar**: `skills/`, `commands/` veya `agents/` içinde hedefi kaybolmuş bir bağlantı ya da var olmayan bir yeri gösteren paylaşılan bir klasör (`skills`, `commands`, `agents`, `plugins`) (bunu, eksik bir grup klasörünü yeniden oluşturan `ccam doctor --fix` ile onarın).
- **`SKILL.md` olmayan bir beceri klasörü**: Claude Code onu atlar.
- **İki yerde farklı tanımlanmış bir MCP sunucusu**: grubun `mcp.json` dosyası ile hesabın `.claude.json` dosyası ya da bunlardan biri ile projenin `.mcp.json` dosyası. İki yerde de aynı ayarların olması (`mcp_strategy = "sync"` bunu yazar) sorun değildir. Bir hesabın kendi sunucusunu `ccam mcp import <account> <server>` ile gruba taşıyın.
- **Geçerli JSON olmayan bir `installed_plugins.json`**, `enabledPlugins` içinde adı geçmeyen bir eklenti ve klasörü eksik bir eklenti.
- **Yalnızca tek bir hesabın yüklediği sunucular**: grubundaki diğer hesapların yüklemediği, yalnızca bu hesabın `.claude.json` dosyasında bulunan bir sunucu ve hesabın bu klasör için yerel kapsamındaki bir sunucu (`-s user` verilmedikçe `claude mcp add` bunu yapar). Her biri, onu paylaşan `ccam mcp import` komutuyla birlikte gelir. Düz `claude` komutunun `~/.claude.json` dosyasında bulunup hiçbir hesabın yüklemediği sunucular `ccam migrate --mcp` komutuyla birlikte gelir.
- **Onay bekleyen bir proje sunucusu**: hesap klasöre henüz güvenmemiştir, bu yüzden Claude Code bir kez sorar ve o zamana kadar projenin `.mcp.json` sunucuları bekler. Sorun, onları nerede onaylayacağınızı söyler: ayarlar grubunun hesaplarının ya da bu hesabın `settings.json` dosyasındaki `"enabledMcpjsonServers"`. Bir git deposunun ana çalışma ağacına verilen güven, diğer çalışma ağaçlarında da geçerlidir.
- **Girişler ve bağlayıcılar**: başka bir hesap giriş yapmışken bu hesapta giriş gerektiren bir MCP sunucusu (bir giriş yalnızca tek bir hesap için geçerlidir; `/mcp` ile giriş yapın) ve başka bir hesabın kullandığı ama bu hesabın kullanmadığı claude.ai bağlayıcıları.

## Güvenlik ve uç durumlar

- `ccam ext` hiçbir şeyi değiştirmez ve ağ isteği yapmaz. Sunucu adlarını ve hook sayılarını okur; MCP tanımlarındaki `env` veya başlık değerlerini asla okumaz. Giriş denetimi için her hesabın girişindeki MCP girişlerinin adlarını okur (macOS'ta Anahtar Zinciri'nden), tokenları asla okumaz; bunu da yalnızca başka bir hesabın da yüklediği bir sunucuyu yükleyen hesaplar için yapar.
- Diskte olanı gösterir. Yönetilen ayarlarla engellenmiş bir sunucu yine görünür; ayarlarda kapatılmış bir eklenti `(off)` gösterir.
- O projenin `.claude` klasörünü ve `.mcp.json` dosyasını da dahil etmek için komutu proje klasöründe çalıştırın.

## İlgili

- [Hesaplar arasında paylaşım](/features/share)
- [Paylaşılan MCP sunucuları](/features/mcp)
- [Kontrol etme ve onarma](/features/doctor)
