---
title: Claude Code kullanım limitine ulaştığında ne yapmalı
description: Sıfırlanmayı bekleyin ya da aynı konuşmaya hemen başka bir Claude Code hesabında devam edin. 5 saatlik ve haftalık limitlerin işleyişi, adım adım.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: ff765a189332
---

Bir hesap kullanım limitine (çoğu zaman rate limit olarak da anılır) ulaştığında iki seçeneğiniz vardır: limitin sıfırlanmasını beklemek ya da aynı konuşmaya şimdi sahip olduğunuz başka bir hesapta devam etmek. Konuşma geçmişi hesaplar arasında paylaşıldığında ccam ile ikincisi tek bir komut sürer: `ccam next --auto --continue`.

## Claude Code'un hangi limitleri var?

Pro ve Max planlarında Claude Code'un [beş saatlik bir oturum limiti ve haftalık bir limiti](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) vardır. İkisi de web'deki ve uygulamalardaki Claude ile paylaşılır ve kullanım ikisinden aynı anda düşülür. Biri dolduğunda Claude Code "You've hit your session limit · resets 3:45pm" gibi [bir mesaj gösterir](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit). Model değiştirmek işe yaramaz, çünkü limitler tüm modelleri kapsar.

## Seçenek 1: sıfırlanmayı bekleyin

Tek hesabınız varsa bekleyin. Claude Code'un son sürümleri [açık oturumda bekler](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) ve limit sıfırlandığında göreve kendiliğinden devam eder. Claude Code içindeki `/usage`, planınızın limitlerinin ne kadarını kullandığınızı gösterir.

## Seçenek 2: hemen başka bir hesapta devam edin

İkinci bir hesabınız varsa aynı konuşmaya hemen devam edebilirsiniz. Bir kez kurun:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Artık siz çalışırken durum satırı hesabı ve 5 saatlik kullanımını gösterir; böylece limit sizi hazırlıksız yakalamaz.

Limite ulaşıldığında Claude Code'dan çıkın (`/exit`), ardından:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Konuşmayı en çok payı kalan hesapta açın:

```sh
ccam next --auto --continue
```

Claude Code, `work` olarak, geçerli klasörün en son konuşmasıyla ve o ana kadar söylenen her şeyle birlikte başlar.

## ccam hesapları otomatik olarak değiştirir mi?

Hayır. Siz çalışırken ccam çalışmaz: Claude Code'u başlatır ve kendi yerini ona bırakır. Ne zaman geçiş yapacağınıza siz karar verirsiniz; `ccam next --auto`, 5 saatlik kullanımı en düşük olan hesabı seçer ve %100'de olan ya da oturum açılmamış hesapları atlar. Bazı diğer araçlar bunun yerine çalışan bir oturumun girişini değiştirir; farkı [ccam ve diğer hesap değiştirme araçları](/compare/account-switchers) sayfası açıklar.

## Akılda tutulması gerekenler

- `ccam next --continue` komutunu çalıştığınız klasörde çalıştırın; konuşmalar proje klasörü bazında saklanır.
- Aynı konuşma üzerinde iki hesabı aynı anda çalıştırmaktan kaçının. Bunu denerseniz ccam sizi uyarır.
- Diğer hesap da konuşma geçmişini aynı grupta paylaşıyor olmalıdır.
- Yeni hesaptaki ilk mesaj limitin normalden fazlasını kullanabilir, çünkü o hesap konuşmayı henüz önbelleğe almamıştır.
- Hesaplar arasında bu şekilde geçiş yapıp yapamayacağınız her hesabın koşullarına bağlıdır.

`ccam usage` ve `ccam next` komutlarının tüm seçenekleri için [kullanım limitleri](/docs/usage-limits) sayfasına bakın.
