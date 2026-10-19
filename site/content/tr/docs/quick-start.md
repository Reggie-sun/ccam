---
title: Hızlı başlangıç
description: ccam add ile iki Claude Code hesabı ekleyin, claude1, claude2 veya ccam <name> ile aralarında geçiş yapın ve neyi paylaşacaklarını seçin.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

Her hesabı `ccam add` ile ekleyin, ardından Claude Code'u `claude1`, `claude2` veya `ccam <name>` ile başlatın. Çoğu kişinin ihtiyacı bu kadarıdır.

## 1. İlk hesabınızı ekleyin

```sh
ccam add
```

Sihirbaz, `personal` veya `work` gibi bir ad ve hesabın diğer hesaplarınızla hangi öğeleri paylaşacağını sorar. Ardından `claude auth login` komutunu çalıştırır; bu komut giriş yapabilmeniz için tarayıcıyı açar. İşlem bittiğinde şunu görürsünüz:

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam kendi mesajlarını İngilizce yazdırır; komutlar ve seçenekler her dilde aynıdır.

## 2. İkinci hesabı ekleyin

```sh
ccam add work
```

Tarayıcıda diğer hesapla giriş yapın. Bu hesap `claude2` olur.

## 3. Geçiş yapın

Bunlardan herhangi biri Claude Code'u work hesabıyla başlatır:

```sh
claude2
ccam work
ccam 2
```

Hesap adından sonra gelen her şey Claude Code'a değiştirilmeden iletilir; böylece bildiğiniz seçenekleri kullanmaya devam edersiniz:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

`ccam` tek başına çalıştırıldığında bir seçici açar: ok tuşlarını ve Enter tuşunu kullanın ya da hesabın numarasına basın.

## 4. Neyin paylaşılacağını seçin

Varsayılan olarak MCP sunucuları, `settings.json`, genel `CLAUDE.md` ile ajanlar, komutlar, beceriler ve eklentiler hesaplar arasında paylaşılır. Konuşma geçmişi paylaşılmaz. Konuşmalara hesaplar arasında devam etmek istiyorsanız bunu iki hesap için de açın:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam, herhangi bir şeyi değiştirmeden önce neyi taşıyacağını gösterir. Ayrıntılar için [paylaşım](/docs/sharing) sayfasına bakın.

## 5. Bir hesap limitine ulaştığında

```sh
ccam usage
ccam next --continue
```

`ccam usage`, tüm hesapların 5 saatlik ve haftalık kullanımını gösterir. `ccam next --continue`, aynı konuşmayı bir sonraki hesapta açar. [Kullanım limitleri](/docs/usage-limits) sayfasına bakın.

## Zaten ~/.claude kullanıyor musunuz?

Yeniden giriş yapmak yerine onu kaydedin:

```sh
ccam import
```

Klasör olduğu yerde kalır ve düz `claude` komutu eskisi gibi çalışmaya devam eder.
