---
title: ccam ile tanışın
description: ccam'ın neden var olduğu, birden fazla Claude Code hesabı olan kişiler için ne yaptığı ve işinize nasıl karışmadığı.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 5e4121938d8c
---

ccam, aynı bilgisayarda birden fazla Claude Code hesabı kullanan kişiler için küçük bir komut satırı aracıdır: bir kişisel ve bir iş hesabı ya da uzun bir proje için birkaç hesap.

## Sorun

Claude Code, `CLAUDE_CONFIG_DIR` aracılığıyla birden fazla hesabı destekler: her hesap için bir klasör. Pratikte her klasör bir adadır. Ayarlar, `CLAUDE.md`, MCP sunucuları, ajanlar ve beceriler her klasörde yeniden kurulmalı ve elle eşit tutulmalıdır. Konuşmalar başladıkları yerde kalır. Bir hesap kullanım limitine ulaştığında da bunu ancak deneyerek öğrenirsiniz.

## ccam ne yapar

- `ccam add` klasörü oluşturur, girişi çalıştırır ve hesaba bir kısayol verir: `claude1`, `claude2` ve devamı.
- `ccs work` ve `claude2`, Claude Code'u o hesapla başlatır ve her argümanı değiştirmeden iletir.
- MCP sunucuları, ayarlar, `CLAUDE.md` ve uzantılar bağlantılar aracılığıyla, öğe öğe ve hesap hesap paylaşılır. Giriş bilgileri hiçbir zaman paylaşılmaz.
- `ccs usage` tüm hesapların 5 saatlik ve haftalık kullanımını gösterir; `ccs next --continue` konuşmayı bir sonraki hesapta sürdürür.

## İşinize karışmamak

Her oturumdan önce çalıştırdığınız bir aracın fark edeceğiniz bir maliyeti olmamalıdır. ccam bir yapılandırma dosyası okur, paylaşılan bağlantıların her birini tek bir sistem çağrısıyla kontrol eder ve ardından kendini `claude` ile değiştirir. Yaklaşık 4 milisaniye ekler ve siz çalışırken çalışır durumda değildir.

Dosyalarınıza da özen gösterir. Paylaşımdaki her değişiklik, başlamadan önce bir plan yazdırır ve orijinalleri bir yedek klasöründe saklar. `ccam doctor` bağlantıları, izinleri ve girişleri kontrol eder ve güvenle onarabildiklerini onarır.

## Deneyin

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam açık kaynaklıdır ve resmi değildir: Anthropic ile bağlantılı değildir ve Anthropic tarafından onaylanmamıştır. Birden fazla hesap kullanıp kullanamayacağınız her hesabın koşullarına bağlıdır. Geri bildirimlere ve çevirilere [GitHub](https://github.com/jhyunwoo/ccam) üzerinde açığız.
