---
title: Hızlı başlangıç
metaTitle: "Hızlı başlangıç: Claude Code hesabı ekle ve değiştir"
description: ccam add ile iki Claude Code hesabı ekleyin, onlara ccwork gibi kısayollar verin, aralarında geçiş yapın ve neyi paylaşacaklarını seçin.
updated: 2026-10-24
order: 2
translation: machine
sourceHash: "e98d0637356a"
---

Her hesabı `ccam add` ile ekleyin, ardından Claude Code'u `ccam <name>` ile veya adını kendiniz seçtiğiniz `ccwork` gibi bir kısayolla başlatın. Çoğu kişi için bu kadarı yeterlidir.

## 1. İlk hesabınızı ekleyin

```sh
ccam add
```

Sihirbaz, `personal` veya `work` gibi bir ad, hesabın diğer hesaplarınızla hangi öğeleri paylaşacağını ve bir kısayol komutu oluşturulup oluşturulmayacağını sorar (4. adıma bakın). Ardından `claude auth login` komutunu çalıştırır; bu komut giriş yapabilmeniz için tarayıcıyı açar. Uzak bir makinede tarayıcı, yüklenmeyen bir `http://localhost:…/callback` sayfasında kalabilir; bu adresi terminale yapıştırın, ccam girişi tamamlar. İşlem bittiğinde şunu görürsünüz:

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

Yalnızca `ccam` çalıştırıldığında [kontrol paneli](/features/dashboard) açılır: ok tuşlarıyla veya numarasıyla bir hesaba gidin, ardından Enter'a basın.

Kısayoldan veya hesap adından sonra gelen her şey Claude Code'a değiştirilmeden iletilir, böylece bildiğiniz seçenekleri kullanmaya devam edersiniz:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Neyin paylaşılacağını seçin

Varsayılan olarak MCP sunucuları, `settings.json`, genel `CLAUDE.md` ile ajanlar, komutlar, beceriler ve eklentiler `default` grubu üzerinden hesaplar arasında paylaşılır. Konuşma geçmişi paylaşılmaz. Konuşmalara hesaplar arasında devam etmek istiyorsanız onu gruba ekleyin:

```sh
ccam group set default --add history
```

ccam, herhangi bir şeyi değiştirmeden önce her hesap için neyi taşıyacağını gösterir. Ayrıntılar için [gruplarla paylaşım](/features/share) sayfasına bakın.

## 6. Bir hesap limitine ulaştığında

```sh
ccam usage
ccam next --continue
```

`ccam usage`, tüm hesapların 5 saatlik ve haftalık kullanımını gösterir. `ccam next --continue`, aynı konuşmayı bir sonraki hesapta açar. Betiklerde `ccam --fallback work -p "…"` kendiliğinden geçiş yapar. [Kullanım limitleri](/docs/usage-limits) sayfasına bakın.

**Team ve Enterprise koltukları.** Otomatik seçim (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar. Anthropic'in yardım merkezinde şöyle belirtilir: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team planındaki her koltuk, adı belirli tek bir kişiye atanır ve onu yalnızca o kişi kullanabilir) ve "Users can't route their work through another user's account to use that user's remaining usage" (kullanıcılar, başka bir kullanıcının kalan kullanım hakkından yararlanmak için işlerini o kullanıcının hesabı üzerinden yürütemez) ([Team planları](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise planları](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) için de durum aynıdır). Bu nedenle otomatik seçimi koltuklar arasında kullanmak Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce kuruluşunuza ve Anthropic'e danışın. Kişisel planlar (Pro, Max), [SSS'deki koşullar notunda](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ele alınır.

## Zaten ~/.claude kullanıyor musunuz?

Yeniden giriş yapmak yerine onu kaydedin:

```sh
ccam import
```

Klasör olduğu yerde kalır ve düz `claude` komutu eskisi gibi çalışmaya devam eder.
