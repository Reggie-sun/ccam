---
title: Tek bilgisayarda iki Claude Code hesabı nasıl kullanılır
description: Birden fazla Claude Code hesabını CLAUDE_CONFIG_DIR ile yan yana çalıştırın; kabuk takma adlarıyla elle ya da ccam ile otomatik olarak.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d82129ddcb4c
---

Her hesaba kendi yapılandırma klasörünü verin. Claude Code klasörü `CLAUDE_CONFIG_DIR` ortam değişkeninden okur ve her klasör kendi giriş bilgilerini, ayarlarını ve geçmişini tutar. Tek bilgisayarda iki Claude Code hesabı için bunu bir kabuk takma adıyla elle kurabilir ya da işi ccam'a bırakabilirsiniz.

## Elle, bir takma adla

Bu, Claude Code belgelerinin anlattığı yaklaşımdır. Her ek hesap için `~/.bashrc` veya `~/.zshrc` dosyasına bir satır ekleyin:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Yeni bir terminal açın ve `claude-work` komutunu çalıştırın. Klasör yeni olduğu için Claude Code sizden giriş yapmanızı ister. Bundan sonra `claude` ilk hesabınız, `claude-work` ise ikinci hesabınızdır.

Bu yöntem iki hesap için iyi çalışır. Bilmeniz gerekenler:

- Her klasör boş başlar: ayar yok, `CLAUDE.md` yok, MCP sunucusu yok, ajan veya beceri yok. İhtiyacınız olanı kopyalar ya da bağlarsınız ve kopyaları kendiniz eşit tutarsınız.
- Takma adlar yalnızca etkileşimli kabuklarda vardır. Betikler, Makefile'lar ve diğer programlar onları görmez.
- MCP sunucuları giriş bilgilerinin yanında, `.claude.json` içinde bulunur; bu yüzden klasörler arasında doğrudan bağlanamazlar.
- Konuşmalar başladıkları hesapta kalır; bu yüzden bir konuşmaya başka bir hesapta devam edemezsiniz.
- Tüm hesapların kullanımını bir arada gösteren hiçbir şey yoktur.

## ccam ile

ccam aynı fikri otomatikleştirir ve eksik parçaları ekler:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Artık `claude1` ve `claude2` ya da `ccs personal` ve `ccs work` iki hesabı başlatır. Bunlar takma ad değil gerçek programlar olduğu için betiklerde de çalışır. Ayarlar, `CLAUDE.md`, MCP sunucuları ve uzantılar baştan paylaşılır; ayrıca paylaşılan konuşma geçmişini açarak `ccs next --continue` ile bir konuşmaya diğer hesapta devam edebilirsiniz.

Zaten `~/.claude` kullanıyorsanız o hesabı yeniden eklemek yerine `ccam import` komutunu çalıştırın.

## Hangisini seçmeli

İki hesap için ve hiçbir şey paylaşılmayacaksa bir takma ad yeterlidir. Paylaşılan ayarlar ve MCP sunucuları, daha fazla hesap, tek bakışta kullanım ya da konuşmalara devam etme istediğinizde ccam sizi bu takip işinden kurtarır. [Karşılaştırma](/compare/aliases) sayfası farkları yan yana listeler.
