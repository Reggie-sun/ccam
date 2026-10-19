---
title: Kullanım limitleri ve hesap değiştirme
description: Tüm Claude Code hesaplarının 5 saatlik ve haftalık kullanımını ccam usage ile görün, en çok pay kalan hesabı ccam next --auto ile açın.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a64c2c341db3
---

`ccam usage`, her hesabın 5 saatlik ve haftalık limitlerinin ne kadarını kullandığını gösterir. `ccam next --auto` en çok pay kalan hesabı açar, `--continue` ise mevcut konuşmayı o hesaba taşır.

```console
$ ccam usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccam next --auto 로 바로 열기)
```

Sütunlar şunlardır: hesap, 5 saatlik kullanım, sıfırlanmaya kalan süre, haftalık kullanım, sıfırlanmaya kalan süre ve sayıların ne kadar eski olduğu.

## Sayılar nereden geliyor

ccam, her hesabın kendi giriş bilgilerini kullanarak tüm hesaplar için aynı anda Anthropic'in kullanım uç noktasına sorgu gönderir. Hesap başına en fazla 3 saniye bekler ve her hesabı yanıt verir vermez gösterir.

Bu uç nokta herkese açık olarak belgelenmemiştir ve değişebilir. Bir istek başarısız olduğunda ccam bilinen son sayıları gösterir ve bunları önbellekten alınmış olarak işaretler. Yanıt biçimi değiştiyse `ccam doctor` bunu size bildirir.

60 saniyeden daha kısa süre önce alınan sayılar yeniden kullanılır. Seçenekler:

| Seçenek | Etki |
| --- | --- |
| `--watch [seconds]` | Sürekli yeniler, varsayılan olarak 60 saniyede bir |
| `--json` | Makinece okunabilir çıktı |
| `-v` | Modele özgü haftalık limitleri ve her sayının nereden geldiğini de gösterir |
| `--no-cache` | Her zaman yeniden sorar |
| `--source cache` | Hiç sormaz; yalnızca kayıtlı sayıları gösterir |

## Süresi dolmuş girişler

Bir hesabın erişim tokenının süresi dolduğunda ccam onu yenilemez; çünkü aynı giriş bilgilerini bir Claude Code oturumu kullanırken token yenilenirse o oturumdan çıkış yapılabilir. Hesabı bir kez `ccam <name>` ile başlatın, Claude Code tokenı yeniler.

`ccam usage --refresh` süresi dolmuş tokenları kendisi yeniler, ancak yalnızca çalışmayan hesaplar için. macOS Anahtar Zinciri'nde veya içe aktarılmış bir `~/.claude` klasöründe saklanan girişleri hiçbir zaman yenilemez.

## Durum satırında kullanım

```sh
ccam statusline install
```

Bundan sonra Claude Code, durum satırında hesap adını ve kullanımını gösterir, örneğin `[work] · Opus · 5h 24% · 7d 41%`. Zaten bir durum satırı komutunuz varsa ccam onu korur ve hesap adını başına ekler. Durum satırı sayılarını da kaydeder; uç nokta yanıt vermediğinde `ccam usage` bu sayıları kullanır.

## Geçiş

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` giriş yapılmamış hesapları atlar, `--auto` ise %100'deki hesapları atlar. Bir konuşmaya devam etmek için iki hesabın da konuşma geçmişini paylaşması gerekir; bkz. [paylaşım](/docs/sharing).

[Kullanım limiti iş akışı kılavuzu](/guides/usage-limit-workflow) eksiksiz bir örneği adım adım anlatır.
