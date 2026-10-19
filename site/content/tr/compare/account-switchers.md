---
title: ccam ve diğer Claude Code hesap değiştirme araçları
description: ccam'ın claude-swap, CCS, clauth, CAAM, cux ve Claude Code hesap değiştirme için kullanılan diğer araçlarla karşılaştırması ve hangi yaklaşımın neye uyduğu.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Claude Code hesap değiştirme araçları iki yaklaşımdan birini izler. Çoğu, tek `~/.claude` klasörünün içindeki girişi değiştirir; bu yüzden aynı anda yalnızca bir hesap etkindir. ccam'ın de aralarında olduğu diğerleri, `CLAUDE_CONFIG_DIR` ile her hesaba kendi klasörünü verir; böylece her hesap oturum açık kalır ve birkaçı aynı anda çalışabilir. ccam buna paylaşılan ayarları ve MCP sunucularını ekler ve bir konuşmaya başka bir hesabın klasöründe devam eder.

## Tek bir girişi değiştirmek mi, ayrı klasörler mi?

**Girişi değiştirmek**, tek bir klasör olduğu için tek bir ayar, MCP sunucusu ve geçmiş kümesini korur. Geçiş hızlıdır ve çalışan bir oturum yeni girişi alabilir. Ancak aynı anda yalnızca bir hesap etkindir ve bu, tüm terminaller ve düzenleyiciler için aynı anda geçerlidir.

**Ayrı klasörler**, [Claude Code belgelerinde anlatılan](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) yöntemdir. Her hesap oturum açık kalır ve bir iş oturumuyla kişisel bir oturumu yan yana çalıştırabilirsiniz. Bedeli, bir araç bunları paylaştırmadıkça her klasörün kendi ayarları, MCP sunucuları ve geçmişiyle başlamasıdır.

## Araçların karşılaştırması

Her projenin kendi belgelerinden 19 Ekim 2026'da alınmıştır. "Belgelenmemiş", özelliğin olmadığı değil, projenin bu konuda bir şey söylemediği anlamına gelir.

| Araç | Yaklaşım | Aynı anda birden fazla hesap | 5 saatlik ve haftalık kullanım | Limite ulaşıldığında | Platformlar |
| --- | --- | --- | --- | --- | --- |
| ccam | Ayrı klasörler | Evet | Evet, `ccam usage` | Konuşmaya en çok payı kalan hesapta devam etmek için `ccam next --auto --continue` komutunu siz çalıştırırsınız | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Girişi değiştirir; `cswap run` ile ayrı klasörler | Evet, `cswap run` ile | Evet | `cswap auto` bir eşikte girişi değiştirir; çalışan oturum devam eder | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Ayrı klasörler | Evet | Kısmen, bir macOS menü çubuğu uygulamasında | Hesap havuzları istekleri yerel bir proxy üzerinden sıradaki hesaba iletir | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Girişi değiştirir; `clauth start` ile ayrı klasörler | Evet, `clauth start` ile | Evet | Bir yedek hesap zinciri boyunca geçiş yapar | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Girişi değiştirir veya yalıtılmış profiller | Evet, yalıtılmış profillerle | Evet | `caam run` etkileşimsiz komutları başka bir hesapta yeniden dener | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Girişi değiştirir ve `claude` komutunu sarar | Belgelenmemiş | Evet | Girişi değiştirir ve aynı konuşmayı sürdürür | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Girişi değiştirir | Hayır | Evet | Belgelenmemiş | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Girişi değiştirir | Hayır | Evet | İsteğe bağlı olarak %100'de geçiş | macOS |

ccswap, Codex girişlerini de yöneten bir claude-swap çatalıdır. Bu araçların birçoğu Codex veya Gemini gibi başka kodlama ajanlarını da yönetir.

## ccam neye odaklanır?

- Her hesap oturum açık kalır. Geçiş hiçbir hesabın oturumunu kapatmaz ve `ccam work` ile `ccam personal` iki terminalde aynı anda çalışabilir.
- MCP sunucuları, `settings.json`, `CLAUDE.md`, ajanlar, beceriler ve eklentiler klasörler arasında bağlantılarla, öğe öğe ve hesap hesap paylaşılır. Giriş bilgileri asla paylaşılmaz.
- Paylaşılan geçmişle tek bir komut, `ccam next --auto --continue`, aynı konuşmayı en çok payı kalan hesabın klasöründe açar.
- ccam Claude Code'u başlatır ve kendi yerini ona bırakır; bu yüzden arka planda hiçbir şey çalışmaz. Kullanımınızı izlemez ve kendiliğinden geçiş yapmaz; ne zaman geçileceğine siz karar verirsiniz.

## Diğer araçlar nerede daha fazlasını yapar?

- claude-swap, CCS, clauth, CAAM ve cux WSL olmadan Windows'ta çalışır. ccam çalışmaz.
- claude-swap, clauth ve cux bir hesap kullanım eşiğini aştığında kendiliğinden geçiş yapar. ccam geçişi size bırakır.
- CCSwitcher, Claude Switcher ve claude-swap kullanımı macOS menü çubuğunda gösterir.
- CCS, CAAM, ccswap ve Claude Switcher Codex'i ve diğer kodlama ajanlarını da yönetir.

## Hangisi size uyar?

Aynı anda tek bir hesap kullanıyorsanız, geçişin otomatik olmasını istiyorsanız ya da WSL olmadan Windows kullanıyorsanız giriş değiştiren bir araç seçin. Hesapları yan yana, tek bir ortak kurulumla kullanmak ve aynı konuşmaya sıradaki hesapta devam etmek istiyorsanız ccam'ı seçin. Hiçbir şeyin paylaşılmasına ihtiyacınız yoksa [düz kabuk takma adları](/compare/aliases) yeterlidir.

Birden fazla hesap kullanıp kullanamayacağınız her hesabın koşullarına bağlıdır. Kullanmadan önce bunları kontrol edin.
