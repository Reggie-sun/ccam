---
title: ccam ile tanışın
description: ccam'ın neden var olduğu, birden fazla Claude Code hesabı olan kişiler için ne yaptığı ve işinize nasıl karışmadığı.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: a555dc969142
---

ccam, aynı bilgisayarda birden fazla Claude Code hesabı kullanan kişiler için küçük bir komut satırı aracıdır: bir kişisel ve bir iş hesabı ya da uzun bir proje için birkaç hesap.

## Sorun

Claude Code, [`CLAUDE_CONFIG_DIR`](/guides/claude-config-dir) aracılığıyla birden fazla hesabı destekler: her hesap için bir klasör. Pratikte her klasör bir adadır. Ayarlar, `CLAUDE.md`, MCP sunucuları, ajanlar ve beceriler her klasörde yeniden kurulmalı ve elle eşit tutulmalıdır. Konuşmalar başladıkları yerde kalır. Bir hesap kullanım limitine ulaştığında da bunu ancak deneyerek öğrenirsiniz.

## ccam ne yapar

- `ccam add` klasörü oluşturur, girişi çalıştırır ve isterseniz hesaba adını kendiniz seçtiğiniz `ccwork` gibi bir kısayol verir.
- `ccam work` ve `ccwork`, Claude Code'u o hesapla başlatır ve her argümanı değiştirmeden iletir.
- MCP sunucuları, ayarlar, `CLAUDE.md` ve uzantılar bağlantılar aracılığıyla, öğe öğe ve hesap hesap paylaşılır. Giriş bilgileri hiçbir zaman paylaşılmaz.
- `ccam usage` tüm hesapların 5 saatlik ve haftalık kullanımını gösterir; `ccam next --continue` konuşmayı bir sonraki hesapta sürdürür.

## İşinize karışmamak

Her oturumdan önce çalıştırdığınız bir aracın fark edeceğiniz bir maliyeti olmamalıdır. ccam bir yapılandırma dosyası okur, paylaşılan bağlantıların her birini tek bir sistem çağrısıyla kontrol eder ve ardından kendini `claude` ile değiştirir. Yaklaşık 4 milisaniye ekler ve siz çalışırken çalışır durumda değildir.

Dosyalarınıza da özen gösterir. Paylaşımdaki her değişiklik, başlamadan önce bir plan yazdırır ve orijinalleri bir yedek klasöründe saklar. `ccam doctor` bağlantıları, izinleri ve girişleri kontrol eder ve güvenle onarabildiklerini onarır.

## Deneyin

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

[Hızlı başlangıç](/docs/quick-start) iki hesap eklemeyi adım adım gösterir, [birden fazla Claude Code hesabı nasıl kullanılır](/guides/multiple-claude-code-accounts) kılavuzu ise yaklaşımı ccam ile ve ccam olmadan açıklar.

ccam açık kaynaklıdır ve resmi değildir: Anthropic ile bağlantılı değildir ve Anthropic tarafından onaylanmamıştır. Birden fazla hesap kullanıp kullanamayacağınız her hesabın koşullarına bağlıdır. Geri bildirimlere ve çevirilere [GitHub](https://github.com/jhyunwoo/ccam) üzerinde açığız.
