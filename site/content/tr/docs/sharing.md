---
title: MCP sunucularını, ayarları ve geçmişi paylaşma
description: Girişler ayrı kalırken MCP sunucuları, settings.json, CLAUDE.md, ajanlar, beceriler ve konuşma geçmişi Claude Code hesapları arasında ccam ile nasıl paylaşılır.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 0864ea87c3e9
---

Her hesap kendi giriş bilgilerini korur. Geri kalan her şey öğe öğe paylaşılabilir: ccam, hesabın dosyasını veya klasörünü paylaşılan bir grup klasörüne bağlar; böylece bir hesapta yapılan değişiklik o gruptaki tüm hesaplarda görünür.

## Neler paylaşılabilir

| Öğe | Seçenek | Paylaşılan |
| --- | --- | --- |
| MCP sunucuları | `--mcp` | Grubun `mcp.json` dosyası, başlangıçta Claude Code'a iletilir |
| Ayarlar | `--settings` | `settings.json`: izinler, hook'lar, model, durum satırı |
| Genel bellek | `--memory` | `CLAUDE.md` |
| Uzantılar | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Konuşma geçmişi | `--history` | `projects/` ve `history.jsonl` |

Yeni hesaplar konuşma geçmişi dışında her şeyi paylaşır.

Giriş dosyaları `.credentials.json` ve `.claude.json` hiçbir zaman paylaşılmaz. ccam yalnızca sabit bir listedeki yolları bağlar ve bu iki dosya o listede yoktur.

## Gruplar

`on`, `default` grubu anlamına gelir. Bir öğeyi yalnızca bazı hesapların paylaşması için bir ad da verebilirsiniz:

```sh
ccam share work --history team
ccam share side --history team
```

`work` ve `side` artık konuşmaları birbiriyle paylaşır, ancak diğer hesaplarınızla paylaşmaz.

## Paylaşımı açma

Bir hesap bir öğeyi paylaşmaya başladığında ccam:

1. taşıyacağı, kopyalayacağı veya bağlayacağı her dosyanın planını yazdırır ve başlamadan önce onay ister;
2. grup henüz boşsa hesabın dosyalarını grup klasörüne taşır, değilse grupta henüz olmayan dosyaları ekler;
3. aynı dosya farklıysa grubun sürümünü mü koruyacağınızı, sizinkini mi kullanacağınızı yoksa JSON'u mu birleştireceğinizi sorar;
4. orijinalinizi hesap klasörünün içindeki `.ccam-backup/<time>/` altında saklar;
5. orijinali gruba yönelen bir bağlantıyla değiştirir.

Konuşma dosyaları oturum kimliklerine göre adlandırılır; bu nedenle birden fazla hesabın geçmişi çakışma olmadan birleşir.

Claude Code o hesapla çalışırken ccam paylaşımı değiştirmez, çünkü oturum taşınan dosyalara yazmaya devam eder. Önce oturumdan çıkın.

## Paylaşımı kapatma

```sh
ccam share work --history off
```

Bağlantı kaldırılır ve hesap, grubun mevcut içeriğinin kendine ait bir kopyasını alır. Bunun yerine boş başlamak için `--fresh` ekleyin. Grup, diğer hesaplar için içeriğini korur.

## Mevcut bir kurulumu kopyalama

Claude Code'u ccam'den önce kullandıysanız ayarlarınız `~/.claude` içindedir. `ccam migrate` bunları bir gruba kopyalar; böylece bu öğeleri paylaşan her hesap onlara sahip olur:

```sh
ccam migrate
```

Kopyaladıkları:

- `settings.json` ve genel `CLAUDE.md`;
- `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` ve `rules/` uzantıları;
- kullanıcı MCP sunucularınız, yani `~/.claude.json` içindeki `mcpServers` girdisi, grubun `mcp.json` dosyasına;
- konuşma geçmişi (`projects/` ve `history.jsonl`), ancak yalnızca `--history` eklediğinizde.

Kaynak klasör yalnızca okunur: içindeki hiçbir şey taşınmaz, bağlanmaz veya değiştirilmez ve düz `claude` eskisi gibi çalışmaya devam eder. Giriş dosyaları `.credentials.json` ve `.claude.json` asla kopyalanmaz.

ccam herhangi bir şeyi kopyalamadan önce planı yazdırır ve onay ister:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Sonunda ccam, her öğeyi grup üzerinden paylaşan hesapları listeler. Paylaşmayan hesaplar varsa onlar için paylaşımı açan `ccam share` komutlarını yazdırır.

Grupta bir dosyanın farklı bir kopyası zaten varsa ccam, paylaşımı açarken olduğu gibi sorar: grubun kopyasını korumak mı, sizinkini kullanmak mı. `settings.json` ve diğer JSON dosyaları birleştirilebilir de, `history.jsonl` ise satır satır birleştirilebilir. Grubun önceki kopyaları grup klasörünün içindeki `.ccam-backup/<time>/` altında saklanır.

| Seçenek | Etki |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Yalnızca belirttiğiniz öğeleri kopyalar; bunlar olmadan konuşma geçmişi dışındaki her şeyi |
| `-g`, `--group <group>` | `default` yerine bu gruba kopyalar |
| `-y`, `--yes` | Sormadan kopyalar; bir dosya farklıysa grubun kopyasını korur |

Hiçbir şey kopyalamadan planı görmek için `--dry-run` seçeneğini komuttan önce yazın: `ccam --dry-run migrate`. Başka bir klasörden kopyalamak için onu belirtin: `ccam migrate ~/.claude-work --group work`.

## MCP sunucuları

Claude Code, MCP sunucularını giriş bilgileriyle birlikte `.claude.json` içinde saklar; bu yüzden bu dosya paylaşılamaz. ccam, paylaşılan sunucuları grubun `mcp.json` dosyasında tutar ve Claude Code'u `--mcp-config=<group>/mcp.json` ile başlatır. Kendi `--mcp-config` seçenekleriniz bunun üzerine uygulanmaya devam eder.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Bir oturum içinde `claude mcp add` ile eklediğiniz sunucu yalnızca o hesap için kaydedilir. O hesabı bir sonraki başlatışınızda ccam, sunucunun adını ve onu gruba taşıyan komutu gösterir:

```sh
ccam mcp import work gh
```

Grupta aynı adla farklı bir sunucu zaten varsa ccam onu atlar ve bunu bildirir. `--force` grubun sunucusunu değiştirir ve eski `mcp.json` dosyasını grubun `.ccam-backup/` klasöründe saklar.

OAuth ile giriş yapan MCP sunucuları, tokenlarını hesabın giriş bilgileriyle birlikte saklar; bu yüzden onları her hesapta bir kez yetkilendirin.

Sunucuların `claude mcp list` ve `/mcp` içinde sıradan kullanıcı sunucuları olarak görünmesini tercih ediyorsanız `~/.config/ccam/config.toml` dosyasında `mcp_strategy = "sync"` ayarını yapın. ccam bu durumda onları başlatmadan önce, Claude Code'un kendi dosya kilidi altında ve bir yedekle birlikte, her hesabın `.claude.json` dosyasına yazar.

## Aynı klasörde iki hesap

Konuşma geçmişini paylaşan hesaplar aynı dosyalara yazar. Aynı gruptan başka bir hesabın zaten çalıştığı bir klasörde ikinci bir Claude Code başlatırsanız ccam sizi uyarır ve devam edilip edilmeyeceğini sorar.

## Bir şey bozulduğunda

```sh
ccam doctor
ccam doctor --fix
```

Paylaşılan bir yol sıradan bir dosyaya dönüşmüşse (örneğin bir yedeği elle geri yükledikten sonra), ccam Claude Code'u başlatmadan önce durur ve sizi `ccam doctor --fix` komutuna yönlendirir. Bu komut dosyayı grupla birleştirir ve bir yedek tutar.
