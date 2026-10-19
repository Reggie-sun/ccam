---
title: Hızlı başlangıç
description: ccam add ile iki Claude Code hesabı ekleyin, onlara ccwork gibi kısayollar verin, aralarında geçiş yapın ve neyi paylaşacaklarını seçin.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: eb078efa2ee2
---

Her hesabı `ccam add` ile ekleyin, ardından Claude Code'u `ccam <name>` ile veya adını kendiniz seçtiğiniz `ccwork` gibi bir kısayolla başlatın. Çoğu kişi için bu kadarı yeterlidir.

## 1. İlk hesabınızı ekleyin

```sh
ccam add
```

Sihirbaz, `personal` veya `work` gibi bir ad, hesabın diğer hesaplarınızla hangi öğeleri paylaşacağını ve bir kısayol komutu oluşturulup oluşturulmayacağını sorar (4. adıma bakın). Ardından `claude auth login` komutunu çalıştırır; bu komut giriş yapabilmeniz için tarayıcıyı açar. İşlem bittiğinde şunu görürsünüz:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Mevcut ayarlarınızı taşıyın

Claude Code'u daha önce kullandıysanız ayarlarınız, `CLAUDE.md`, uzantılarınız ve MCP sunucularınız `~/.claude` içindedir. Her hesap bunlara sahip olsun diye onları paylaşılan ayarlara kopyalayın:

```sh
ccam migrate
```

ccam neyi kopyalayacağını gösterir ve önce onay ister. `~/.claude` klasörünü yalnızca okur, giriş bilgilerini asla kopyalamaz ve `--history` eklemediğiniz sürece konuşma geçmişini dışarıda bırakır. Ayrıntılar için [mevcut bir kurulumu kopyalama](/docs/sharing#copy-an-existing-setup) bölümüne bakın.

## 3. İkinci hesabı ekleyin

```sh
ccam add work
```

Tarayıcıda diğer hesapla giriş yapın.

## 4. Hesaplar arasında geçiş

`ccam work`, Claude Code'u iş hesabıyla başlatır. Daha kısası için hesaba `ccwork` gibi kendine ait bir kısayol komutu verin. Paylaşılan öğelerden sonra sihirbaz şunu sorar:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` tuşuna basın, ardından bir ad yazın veya önerilen adı kullanmak için Enter'a basın:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Girişten sonra ccam bunu onaylar:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Bu soruyu atladıysanız istediğiniz zaman bir kısayol ekleyebilirsiniz:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Artık bunlardan herhangi biri Claude Code'u iş hesabıyla başlatır:

```sh
ccwork
ccam work
ccam 2
```

Yalnızca `ccam` çalıştırıldığında bir seçici açılır: ok tuşlarını ve Enter'ı kullanın ya da hesabın numarasına basın.

Kısayoldan veya hesap adından sonra gelen her şey Claude Code'a değiştirilmeden iletilir, böylece bildiğiniz seçenekleri kullanmaya devam edersiniz:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Neyin paylaşılacağını seçin

Varsayılan olarak MCP sunucuları, `settings.json`, genel `CLAUDE.md` ile ajanlar, komutlar, beceriler ve eklentiler hesaplar arasında paylaşılır. Konuşma geçmişi paylaşılmaz. Konuşmalara hesaplar arasında devam etmek istiyorsanız bunu iki hesap için de açın:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam, herhangi bir şeyi değiştirmeden önce neyi taşıyacağını gösterir. Ayrıntılar için [paylaşım](/docs/sharing) sayfasına bakın.

## 6. Bir hesap limitine ulaştığında

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
