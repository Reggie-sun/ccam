---
title: Claude Code kullanım limitine ulaştığında ne yapmalı
description: Her Claude Code hesabında ne kadar pay kaldığını kontrol edin ve bir hesap limitine ulaştığında aynı konuşmaya başka bir hesapta devam edin.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: db7ed9b7bcfc
---

Bir hesap Claude Code kullanım limitine ulaştığında, aynı konuşmaya sahip olduğunuz başka bir hesapta devam edin: konuşma geçmişini bir kez paylaşın, ardından `ccam next --continue` komutunu çalıştırın. Yalnızca bir hesabınız varsa Claude Code'un gösterdiği sıfırlanma zamanını bekleyin.

## Limit dolmadan önce

Hesapları, geçiş tek bir komut sürecek şekilde ayarlayın:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Artık siz çalışırken durum satırı hesabı ve 5 saatlik kullanımını gösterir; böylece limit sizi hazırlıksız yakalamaz.

## Limite ulaşıldığında

Claude Code'dan çıkın (`/exit`), ardından:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

Konuşmayı en çok pay kalan hesapta açın:

```sh
ccam next --auto --continue
```

Claude Code, geçerli klasörün en son konuşmasıyla, o ana kadar söylenen her şey dahil, `work` olarak başlar.

## Akılda tutulması gerekenler

- `ccam next --continue` komutunu çalıştığınız klasörde çalıştırın; konuşmalar proje klasörü bazında saklanır.
- Aynı konuşmada aynı anda iki hesap çalıştırmaktan kaçının. Bunu denerseniz ccam sizi uyarır.
- Diğer hesap da konuşma geçmişini aynı grupta paylaşmalıdır.
- Hesaplar arasında bu şekilde geçiş yapıp yapamayacağınız her hesabın koşullarına bağlıdır.

`ccam usage` ve `ccam next` komutlarının tüm seçenekleri için [kullanım limitleri](/docs/usage-limits) sayfasına bakın.
