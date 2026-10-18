---
title: MCP sunucularını, ayarları ve geçmişi paylaşma
description: Girişler ayrı kalırken MCP sunucuları, settings.json, CLAUDE.md, ajanlar, beceriler ve konuşma geçmişi Claude Code hesapları arasında ccam ile nasıl paylaşılır.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
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

## Paylaşımı kapatma

```sh
ccam share work --history off
```

Bağlantı kaldırılır ve hesap, grubun mevcut içeriğinin kendine ait bir kopyasını alır. Bunun yerine boş başlamak için `--fresh` ekleyin. Grup, diğer hesaplar için içeriğini korur.

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
