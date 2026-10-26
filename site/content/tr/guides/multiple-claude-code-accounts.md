---
title: Tek bilgisayarda birden fazla Claude Code hesabı
description: CLAUDE_CONFIG_DIR ile birden fazla Claude Code hesabını yan yana çalıştırın; elle takma adlarla ya da ccam ile otomatik kurun ve hesaplar arasında geçiş yapın.
updated: 2026-10-26
order: 1
translation: machine
sourceHash: 381de7a473b2
---

Her hesaba kendi yapılandırma klasörünü verin. Claude Code klasörü `CLAUDE_CONFIG_DIR` ortam değişkeninden okur ve her klasör kendi giriş bilgilerini, ayarlarını ve geçmişini tutar; böylece her hesap oturum açık kalır ve hesapları aynı anda çalıştırabilirsiniz. Bunu bir kabuk takma adıyla elle kurabilir ya da ccam'a bırakabilirsiniz.

## Elle nasıl kurulur?

Bu, [Claude Code belgelerinde anlatılan](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) yöntemdir. Her ek hesap için `~/.bashrc` veya `~/.zshrc` dosyasına bir satır ekleyin:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Yeni bir terminal açın ve `claude-work` komutunu çalıştırın. Klasör yeni olduğu için Claude Code sizden giriş yapmanızı ister. Bundan sonra `claude` ilk hesabınız, `claude-work` ise ikinci hesabınız olur.

WSL olmadan Windows'ta, profilinizdeki (`notepad $PROFILE`) bir PowerShell fonksiyonu aynı işi görür:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Takma ad yönteminde neler eksik kalır?

İki hesap için iyi çalışır. Bilmeniz gerekenler:

- Her klasör boş başlar: ayar yok, `CLAUDE.md` yok, MCP sunucusu yok, ajan veya beceri yok. İhtiyacınız olanı kopyalar ya da bağlarsınız ve kopyaları kendiniz eşit tutarsınız.
- Takma adlar yalnızca etkileşimli kabuklarda vardır. Betikler, Makefile'lar ve diğer programlar onları görmez.
- MCP sunucuları giriş bilgilerinin yanında, `.claude.json` içinde durur; bu yüzden klasörler arasında basitçe bağlanamazlar.
- Konuşmalar başladıkları hesapta kalır; bu yüzden bir konuşmaya başka bir hesapta devam edemezsiniz.
- Tüm hesapların kullanımını bir arada gösteren bir şey yoktur.

## Bunun yerine /login ile hesap değiştirebilir miyim?

Evet, yalnızca ara sıra geçiş yapıyorsanız. Claude Code içinde `/login` yazmak, bulunduğunuz klasörün [hesabını değiştirir](https://code.claude.com/docs/en/quickstart). Aynı anda yalnızca bir hesapta oturum açıktır; bu yüzden her geçişte yeniden giriş yaparsınız ve iki hesabı yan yana çalıştıramazsınız. Ayrı klasörler bu iki sorunu da ortadan kaldırır.

## ccam ile nasıl kurulur?

ccam aynı fikri otomatikleştirir ve eksik parçaları ekler:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Artık `ccam personal` ve `ccam work` iki hesabı başlatır; iş hesabına verdiğiniz kısayol olan `ccwork` da öyle. Bunlar takma ad değil gerçek programlar olduğu için betiklerde de çalışır. Listeden bir hesap seçmek için yalnızca `ccam` komutunu, bir hesabı varsayılan yapmak için `ccam use work` komutunu çalıştırın.

Ayarlar, `CLAUDE.md`, MCP sunucuları ve uzantılar baştan paylaşılır. `ccam next --continue` bir konuşmaya diğer hesapta devam eder. Paylaşılan konuşma geçmişini açarsanız iki hesap bir kopya yerine aynı konuşma üzerinde çalışır. `ccam usage` her hesabın 5 saatlik ve haftalık kullanımını bir arada gösterir.

Zaten `~/.claude` kullanıyorsanız o hesabı yeniden eklemek yerine `ccam import` komutunu çalıştırın ya da ayarlarını, MCP sunucularını ve uzantılarını paylaşılan `default` grubuna kopyalamak için `ccam migrate` komutunu kullanın.

## Hangi hesabın çalıştığını nasıl anlarım?

Claude Code içinde `/status`, oturum açtığınız hesabı gösterir. ccam ile bir kez `ccam statusline install` komutunu çalıştırın; durum satırı hesap adını ve kullanımını sürekli gösterir, örneğin `[work] · Opus · 5h 24% · 7d 41%`.

## Bu Windows ve WSL'de çalışır mı?

Claude Code [Windows'ta yerel olarak ve WSL içinde](https://code.claude.com/docs/en/setup) çalışır. `CLAUDE_CONFIG_DIR` ikisinde de çalışır; bu yüzden yerel Windows'ta yukarıdaki takma ad veya PowerShell fonksiyonu yeterlidir. ccam de Windows'ta yerel olarak çalışır: PowerShell'de `irm https://ccam.dev/install.ps1 | iex` ile ya da WSL içinde Linux komutuyla kurun. Kısayollar `ccwork.exe` olur ve `settings.json` gibi dosyaları paylaşmak Geliştirici Modu gerektirir ([ccam kurulumu](/docs/install)).

## Sorun giderme

**Claude Code her seferinde giriş yapmamı istiyor.** Yeni bir klasördeki ilk başlatma her zaman giriş ister. İstemeye devam ederse girişin süresi dolmuş ya da giriş bilgileri silinmiştir: `ccam login work` komutunu çalıştırın (ya da takma adı bir kez başlatıp giriş yapın).

**Bir MCP sunucusu yalnızca bir hesapta görünüyor.** Bir oturum içinde `claude mcp add` ile eklenen sunucular o hesaba, `-s user` olmadan da yalnızca o klasöre aittir. Bir sunucuyu tüm hesaplara taşımak için `ccam mcp import work <server>` komutunu kullanın.

**Ayarlar hesaplar arasında farklı.** `ccam share work`, hesabın neyi, hangi grup üzerinden paylaştığını yazdırır; `ccam group join default work`, neyin değişeceğini gösterdikten sonra hesabın `default` grubunun paylaştıklarını paylaşmasını sağlar. `ccam doctor --fix` bozuk bağlantıları onarır ve değiştirdiği her dosyanın yedeğini tutar.

**İki hesap aynı konuşmayı açıyor.** Aynı konuşma üzerinde iki hesapla aynı anda çalışmaktan kaçının. Bunu denediğinizde ccam sizi uyarır.

## Hangisini seçmeliyim?

İki hesap için ve hiçbir şeyin paylaşılmasını istemiyorsanız bir takma ad yeterlidir. Paylaşılan ayarlar ve MCP sunucuları, daha fazla hesap, kullanımı bir bakışta görmek ya da konuşmalara devam etmek istediğinizde ccam sizi bu takipten kurtarır. [Takma adlarla karşılaştırma](/compare/aliases) farkları yan yana listeler; [ccam ve diğer hesap değiştirme araçları](/compare/account-switchers) ise bunun yerine tek bir girişi değiştiren araçları ele alır.

Birden fazla hesap kullanıp kullanamayacağınız, her birine uygulanan koşullara bağlıdır. Başlamadan önce bu koşulları okuyun.

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.
