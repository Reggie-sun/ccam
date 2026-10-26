---
title: Hızlı başlangıç
metaTitle: "Hızlı başlangıç: Claude Code hesabı ekle ve değiştir"
description: ccam add ile iki Claude Code hesabı ekleyin, onlara ccwork gibi kısayollar verin, aralarında geçiş yapın ve neyi paylaşacaklarını seçin.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

Her hesabı `ccam add` ile ekleyin, ardından Claude Code'u `ccam <name>` ile veya adını kendiniz seçtiğiniz `ccwork` gibi bir kısayolla başlatın. Çoğu kişi için bu kadarı yeterlidir.

Claude Code'u zaten kullanıyor musunuz? Başlamadan önce bu sayfanın sonundaki [Zaten ~/.claude kullanıyor musunuz?](#already-using-claude) bölümünü okuyun.

## 1. İlk hesabınızı ekleyin

```sh
ccam add
```

Sihirbaz, `personal` veya `work` gibi bir ad, hesabın diğer hesaplarınızla hangi öğeleri paylaşacağını ve bir kısayol komutu oluşturulup oluşturulmayacağını sorar (3. adıma bakın). Ardından `claude auth login` komutunu çalıştırır; bu komut giriş yapabilmeniz için tarayıcıyı açar. Uzak bir makinede tarayıcınız, yüklenmeyen bir `http://localhost:…/callback` sayfasında kalabilir; bu adresi terminale yapıştırın, ccam girişi tamamlar. Giriş tamamlandığında şunu görürsünüz:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Mevcut ayarlarınızı taşıyın

Claude Code'u daha önce kullandıysanız ayarlarınız, `CLAUDE.md` ve uzantılarınız `~/.claude` içinde, MCP sunucularınız ise `~/.claude.json` içindedir. Her hesap bunlara sahip olsun diye onları paylaşılan `default` grubuna kopyalayın:

```sh
ccam migrate
```

ccam neyi kopyalayacağını gösterir ve önce onay ister. `~/.claude` klasörünü yalnızca okur, giriş bilgilerini asla kopyalamaz ve `--history` eklemediğiniz sürece konuşma geçmişini dışarıda bırakır. Ayrıntılar için [mevcut bir kurulumu kopyalama](/docs/sharing#copy-an-existing-setup) bölümüne bakın.

## 3. İkinci hesabı ekleyin

```sh
ccam add work
```

Adı komut satırında verdiğiniz için sihirbaz adı sormaz ve katılınacak grupları `default` işaretli olarak gösterir. Bu kez hesaba `ccwork` gibi kendine ait bir kısayol komutu verin, böylece onu tek kelimeyle başlatabilirsiniz. Gruplardan sonra sihirbaz şunu sorar:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` tuşuna basın, ardından bir ad yazın veya önerilen adı kullanmak için Enter'a basın:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Ardından tarayıcıda diğer hesapla giriş yapın. ccam hesabı ve kısayolunu onaylar:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Hesaplar arasında geçiş

`ccam work`, Claude Code'u iş hesabıyla başlatır. Kısayol sorusunu atladıysanız istediğiniz zaman bir kısayol ekleyebilirsiniz:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Artık bu komutların her biri Claude Code'u iş hesabıyla başlatır (`2`, hesabın `ccam list` içindeki numarasıdır):

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

Varsayılan olarak MCP sunucuları, `settings.json`, genel `CLAUDE.md` ve uzantılar (ajanlar, komutlar, beceriler ve eklentiler) `default` grubu üzerinden hesaplar arasında paylaşılır. Konuşma geçmişi paylaşılmaz. Konuşmalara hesaplar arasında devam etmek istiyorsanız onu gruba ekleyin:

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

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.

## Zaten ~/.claude kullanıyor musunuz?

Yeniden giriş yapmak yerine onu kaydedin:

```sh
ccam import
```

Klasör olduğu yerde kalır ve düz `claude` komutu eskisi gibi çalışmaya devam eder.
