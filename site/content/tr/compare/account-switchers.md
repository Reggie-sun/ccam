---
title: ccam ve diğer Claude Code hesap değiştirme araçları
description: ccam'ın claude-swap, CCS, clausona, clauth, CAAM, cux ve diğer Claude Code hesap değiştiricileriyle karşılaştırması ve hangi yaklaşımın neye uyduğu.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: b6143bcf3dc7
---

Claude Code hesap değiştirme araçları iki yaklaşımdan birini izler. Çoğu, tek `~/.claude` klasörünün içindeki girişi değiştirir; bu yüzden aynı anda yalnızca bir hesap etkindir. ccam'ın da aralarında olduğu diğerleri, `CLAUDE_CONFIG_DIR` ile her hesaba kendi klasörünü verir; böylece her hesap oturum açık kalır ve birkaçı aynı anda çalışabilir. ccam buna paylaşılan ayarları ve MCP sunucularını ekler ve bir konuşmaya başka bir hesabın klasöründe devam eder.

## Tek bir girişi değiştirmek mi, ayrı klasörler mi?

**Girişi değiştirmek**, tek bir klasör olduğu için tek bir ayar, MCP sunucusu ve geçmiş kümesini korur. Geçiş hızlıdır ve çalışan bir oturum yeni girişi alabilir. Ancak aynı anda yalnızca bir hesap etkindir ve bu, tüm terminaller ve düzenleyiciler için aynı anda geçerlidir.

Girişi değiştirmek için araç, diğer her hesabın girişinin bir kopyasını tutar ve onu yerine yazar. Claude Code bir girişi yenilerken içindeki yenileme belirtecini değiştirir; bu yüzden kaydedilen her kopyanın çalışan oturumlarla eşit tutulması gerekir. Geçişten sonra oturumu kapanan ya da yeniden `/login` gerektiren hesaplar, bu projelerin en sık düzelttiği sorunlar arasındadır.

**Ayrı klasörler**, [Claude Code belgelerinde anlatılan](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) yöntemdir. Her hesap oturum açık kalır ve bir iş oturumuyla kişisel bir oturumu yan yana çalıştırabilirsiniz. Her giriş kendi klasöründe kalır ve Claude Code onu orada her zamanki gibi yeniler. Bedeli, bir araç bunları paylaştırmadıkça her klasörün kendi ayarları, MCP sunucuları ve geçmişiyle başlamasıdır.

## Araçların karşılaştırması

Her projenin kendi belgelerinden ve kaynak kodundan 20 Ekim 2026'da alınmıştır (clausona: v0.12.0-beta). "Belgelenmemiş", özelliğin olmadığı değil, projenin bu konuda bir şey söylemediği anlamına gelir.

| Araç | Yaklaşım | Aynı anda birden fazla hesap | 5 saatlik ve haftalık kullanım | Limite ulaşıldığında | Platformlar |
| --- | --- | --- | --- | --- | --- |
| ccam | Ayrı klasörler | Evet | Evet, `ccam usage` | Oturumda bir bildirim; `ccam next --auto --continue` veya `/ccam:switch` konuşmayı, paylaşılan geçmiş olsun ya da olmasın, en çok payı kalan hesaba taşır; isteğe bağlı `auto_switch` (Linux, macOS) onu kendiliğinden taşır; `ccam --fallback` `-p` işlerini yeniden çalıştırır; `ccam pick` çalışan süreçleri hesaplara dağıtır | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Girişi değiştirir; `cswap run` ile ayrı klasörler | Evet, `cswap run` ile | Evet | `cswap auto` bir eşikte girişi değiştirir; çalışan oturum devam eder | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Ayrı klasörler | Evet | Kısmen | Hesap havuzları istekleri yerel bir proxy üzerinden sıradaki hesaba iletir | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Ayrı klasörler | Evet | Evet | Rotalar her çalıştırmadan önce kotası kalan bir hesap seçer, bir sıfırlanmayı bekleyebilir ve çalışan süreçlere farklı hesaplar dağıtır | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Girişi değiştirir; `clauth start` ile ayrı klasörler | Evet, `clauth start` ile | Evet | Bir yedek hesap zinciri boyunca geçiş yapar | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Girişi değiştirir veya yalıtılmış profiller | Evet, yalıtılmış profillerle | Evet | `caam run` etkileşimsiz komutları başka bir hesapta yeniden dener | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Girişi değiştirir ve `claude` komutunu sarar | Hayır, tek etkin hesap | Evet | Girişi değiştirir ve aynı konuşmayı sürdürür | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Girişi değiştirir | Hayır | Evet | İsteğe bağlı olarak bir eşikte geçiş | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Girişi değiştirir | Hayır | Evet | İsteğe bağlı olarak %100'de geçiş | macOS |

ccswap, Codex girişlerini de yöneten bir claude-swap çatalıdır. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) ve [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) da her hesaba kendi klasörünü verir. Bu araçların birçoğu Codex veya Gemini gibi başka kodlama ajanlarını da yönetir.

## Nasıl yapılmışlar?

| Araç | Yazıldığı dil | Claude Code çalışırken | Girişleri yerler arasında kopyalar | Hesaplar neleri paylaşabilir |
| --- | --- | --- | --- | --- |
| ccam | Go, tek statik ikili dosya | Varsayılan olarak hiçbir şey: ccam kendini `claude` ile değiştirir (bunun mümkün olmadığı Windows'ta Claude Code'un çıkmasını bekler); `relay` veya `auto_switch` açıkken ccam etkileşimli oturumların yanında çalışmaya devam eder | Hayır | Gruplar: her grup paylaştığı öğeleri ve hesaplarını belirtir; bir hesap birkaç grupta olabilir |
| claude-swap | Python 3.12+, uv veya pipx ile kurulur | `cswap run` ile hiçbir şey; `cswap auto` ve menü çubuğu uygulaması çalışmaya devam eder | Evet, girişi değiştirdiğinde | `cswap run` ile: ayarlar, `CLAUDE.md` ve uzantılar birlikte, geçmiş isteğe bağlı; MCP sunucuları her başlatmada kopyalanır |
| CCS | Node.js 18+ veya Bun | CCS üst süreç olarak kalır; havuzlar yerel bir proxy ekler | Kendi hesap klasörleri için hayır; havuzlar kendi girişlerini tutar | Hesap başına paylaşılan veya yerel, geçmiş bağlam gruplarında; MCP sunucuları kopyalanır |
| clausona | TypeScript, Node.js 20+ | Düz `claude` bir kabuk fonksiyonu üzerinden çalışır; `clausona run` Node'u üst süreç olarak tutar; her çıkıştan sonra arka planda bir Node süreci çalışır | Hayır, ama boştaki hesapların kullanımını okumak için onların token'larını kendisi yeniler | Hesaba özgü dosyalar dışında her şey, tek bir ana hesaptan; geçmiş isteğe bağlı; `claude mcp add` ile eklenen sunucular paylaşılmaz |
| clauth | Rust | clauth üst süreç olarak kalır | Evet, girişi değiştirdiğinde | `~/.claude` klasörünün tamamı ya da hiçbir şey |
| CAAM | Go | Moda bağlı; isteğe bağlı bir arka plan hizmeti | Evet, kasa modunda | Moda göre sabittir; MCP sunucuları her başlatmada kopyalanır |
| cux | Go | `claude` etrafında bir sarmalayıcı kalır | Evet | Her şey: tek klasör |
| CCSwitcher, Claude Switcher | macOS uygulamaları (Swift, Python) | Menü çubuğu uygulaması | Evet | Her şey: tek klasör |

## ccam'ın güçlü yanları

### Girişler, Claude Code'un koyduğu yerde kalır

ccam bir Claude girişini asla bir yerden başka bir yere kopyalamaz. Her hesabın girişi kendi klasöründe, macOS'ta ise o klasöre ait Anahtar Zinciri öğesinde kalır ve onu yenileyen program Claude Code'dur. Eşitliği bozulabilecek kaydedilmiş bir kopya yoktur. `ccam import` ile kaydedilen mevcut bir `~/.claude`, düz `claude` komutunun kullandığı Anahtar Zinciri öğesini korur; böylece iki şekilde de oturum açık kalır.

ccam süresi dolmuş bir belirteci yalnızca siz `ccam usage --refresh` çalıştırdığınızda kendisi yeniler. Önce bir kilit alır ve o hesap çalışıyorsa, Claude Code nasıl başlatılmış olursa olsun (arka plan oturumları ve onların daemon'ı da dahil), ya da Claude Code o girişi kendisi yeniliyorsa işlemi reddeder; böylece canlı bir oturumun altındaki belirteci asla değiştirmez.

### Varsayılan olarak sizinle Claude Code arasında hiçbir şey çalışmaz

ccam bir dosya okur, bağlantılarını denetler ve kendini `claude` ile değiştirir (`exec`). Varsayılan olarak siz çalışırken ccam süreci, sarmalayıcı, proxy veya arka plan hizmeti yoktur: standart girdi, sinyaller ve çıkış kodu Claude Code'a aittir ve istekleri doğrudan Anthropic'e gider. Windows bir süreci değiştiremez; bu yüzden orada ccam, Claude Code'u aynı konsolla başlatır ve başka hiçbir şey yapmadan çıkış kodunu bekler. `relay = true` veya `auto_switch = true` olduğunda ccam, bir konuşmaya başka bir hesapta devam edebilmek için etkileşimli oturumların yanında aynı şekilde çalışmaya devam eder; Claude Code'un istekleri yine doğrudan Anthropic'e gider.

Bir hesabı bu şekilde başlatmak yaklaşık 4 ms ekler; bu süre 95. yüzdelikte 6 ms'nin altındadır ve en yüksek bellek kullanımı 8 MB'tır ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) ve hyperfine ile, 300 çalıştırma). Aynı dizüstü bilgisayarda, bir program henüz hiçbir şey yapmadan, yalnızca Python yorumlayıcısını başlatmak 13 ms, Node.js'i başlatmak 31 ms sürer.

### Neyin paylaşılacağına gruplar karar verir

MCP sunucuları, `settings.json`, `CLAUDE.md`, uzantılar (ajanlar, komutlar, beceriler, eklentiler) ve konuşma geçmişi beş ayrı öğedir. Bir grup paylaştığı öğeleri ve içindeki hesapları belirtir; böylece iki iş hesabı geçmişi bir `work` grubunda paylaşırken kişisel bir hesap kendi geçmişini tutar ve üçü de `default` üzerinden aynı MCP sunucularını paylaşır. `ccam group set work --add memory` gruptaki tüm hesapları tek bir planla günceller; [kontrol paneli](/features/dashboard) aynı şeyi onay kutularıyla yapar ve bir hesabın zaten başka bir gruptan aldığı öğeyi gri gösterir.

- MCP sunucuları başlatmada `--mcp-config` ile verilir; bu yüzden bir başlatma varsayılan olarak, Claude Code'un çalışırken yeniden yazdığı `.claude.json` dosyasına asla dokunmaz. ccam oraya yalnızca bunu açıkça söyleyen adımlarda, Claude Code'un kilidi altında ve bir yedekle yazar: sorduktan sonra Remote Control için, bir hesap MCP paylaşımından ayrıldığında ya da yeniden katıldığında ve `ccam mcp import` için. `ccam mcp add` bir sunucuyu tüm grup için bir kez ekler.
- Diğer öğeler grup klasörüne giden bağlantılardır; bu yüzden bir hesapta yapılan değişiklik diğerlerinde hemen görünür. Başlatmada hiçbir şey kopyalanmaz.
- Girişler, `.claude.json` ve çalışma durumu asla paylaşılmaz.

### Her değişiklik önce gösterilir ve yedeklenir

`ccam group` ve `ccam migrate` bir şeyi taşımadan önce bir plan yazdırır ve onay ister. Bir dosya hesapla grubu arasında farklıysa seçim sizindir: grubunkini koruyun, sizinkini kullanın ya da JSON'u birleştirin; klasörlerin içinde grubun kopyası kalır ve plan farklı olan dosyaları adlarıyla belirtir. Orijinaller `.ccam-backup/<time>/` klasörüne gider, kaldırılan bir hesabın oturumu kapatılır ve klasörü bir çöp klasörüne gider; dosyalar tek bir yeniden adlandırmayla değiştirilir, böylece yarıda kalan bir yazma eski dosyayı yerinde bırakır. `--dry-run`, başlatma, paylaşma, taşıma, kısayollar, güncellemeler ve onarımların ne yapacağını gösterir. `ccam doctor --fix` bozuk bağlantıları ve izinleri onarır.

### Sıradaki hesapta aynı konuşma

`ccam next --auto --continue` tüm hesaplara kullanımlarını aynı anda sorar, en çok payı kalanı açar ve konuşmaya devam eder. Paylaşılan geçmişle bu aynı konuşmadır; paylaşılan geçmiş yoksa Claude Code konuşmayı yeni hesaba kopyalar (`--resume <file> --fork-session`) ve ccam ilk yanıtın kaç token'ı yeniden okuyacağını söyler.

Bir limit oturumu durdurduğunda ccam eklentisinin hook'u bir bildirim gösterir ve `/ccam:switch` konuşmayı devreder. Varsayılan olarak ne zaman geçileceğine siz karar verirsiniz. `auto_switch` açıkken (Linux, macOS) ccam oturumu kendisi taşır, ancak yalnızca Claude Code onu bir limitte durdurduktan sonra ve yalnızca ayrıldığı hesapla aynı türden bir hesaba: aynı kuruluşun başka bir koltuğuna ya da başka bir kişisel plana ([ayrıntılar](/features/next#auto-switch)).

Kullanımınızı arka planda izleyen hiçbir şey yoktur ve varsayılan olarak otomatik seçimler işi asla bir Team veya Enterprise kuruluşunun dışına ya da içine taşımaz; böylece bir şirket konuşması yanlışlıkla kişisel bir hesaba düşmez. Bir hesap %90'a ulaştığında durum satırı bu komutu önerir.

Betikler sizi bekleyemez; bu yüzden `ccam --fallback work -p "…"` onların yerine geçiş yapar: hesap limitine ulaştığında aynı komut en çok payı kalan hesapla yeniden çalışır ve yalnızca tamamlanan çalıştırmanın çıktısı yazdırılır. Kullanım limiti olmayan bir hata olduğu gibi iletilir. Paralel çalışan süreçler için `ccam pick` her birine payı kalan farklı bir hesap verir ve hiçbir hesapta pay kalmadığında ikisi de 75 koduyla çıkar.

### Her klasörde doğru hesap

`ccam folder work ~/code/company`, o klasörün içinde hesap adı verilmeden yapılan her başlatmanın `work` kullanmasını sağlar; diğer projeleriniz varsayılan hesabı kullanmaya devam eder. Ayar, her başlatmanın okuduğu aynı `config.toml` dosyasında durur, bu yüzden zaman kaybettirmez.

### Tek ikili dosya, denetlenen indirmeler

ccam, Linux, macOS ve Windows için, ayrıca bir çalışma zamanı kurmayı gerektirmeyen, yaklaşık 10 MB'lık statik bir ikili dosyadır. Kurulum betikleri (`install.sh`, `install.ps1`) ve `ccam update` her indirmeyi sürümün SHA-256 sağlama toplamlarıyla denetler ve ccam size sormadan asla güncellenmez. Telemetrisi yoktur: yaptığı tek istekler, hesaplarınız için Anthropic'e yapılan kullanım sorguları (paralel, 3 saniyelik zaman aşımı ve 60 saniyelik önbellekle), kontrol paneli açıldığında GitHub'da daha yeni bir sürüm olup olmadığını denetlemesi ve sizin istediğiniz indirmelerdir. `usage_source = "claude"` ile bunlar bile Claude Code üzerinden geçer ve ccam hiçbir token göndermez.

## Diğer araçlar nerede daha fazlasını yapar?

- claude-swap, clauth, cux, CAAM ve CCSwitcher bir hesap kullanım eşiğini aştığında çalışan etkileşimli bir oturumu kendiliğinden başka hesaba geçirebilir. ccam `-p` çalıştırmalarını `--fallback` ile geçirir ve etkileşimli bir oturumu yalnızca `auto_switch` açıkken (Linux, macOS) ve yalnızca Claude Code onu bir limitte durdurduktan sonra kendiliğinden taşır; aksi halde oturumu siz, limitteki bir bildirim ve `/ccam:switch` ile taşırsınız.
- Geliştirici Modu olmayan Windows'ta ccam klasörleri paylaşır ama `settings.json` gibi tek dosyaları paylaşamaz (ilk `ccam add` o zaman yalnızca MCP sunucularını ve uzantıları paylaşır); claude-swap ve CCS orada bunları kopyalar.
- CCSwitcher, Claude Switcher ve claude-swap kullanımı macOS menü çubuğunda gösterir.
- CCS, clausona, CAAM, clauth, ccswap ve Claude Switcher Codex'i ve diğer kodlama ajanlarını da yönetir.
- clausona ayrıca uzantıları hesap başına açıp kapatır (ccam'in `ccam ext` komutu bunları yalnızca gösterir), yerel token maliyetini kendisi izler (ccam bu işi [ccusage](https://ccusage.com) aracına bırakır: `ccam ccusage`) ve tmux ile diğer terminallerde çalışan süreç filoları yürüten Claude Code eklentileriyle gelir.
- CAAM, CCSwitcher ve Claude Switcher Homebrew ile kurulur; clauth ve CAAM sürüm sağlama toplamlarını imzalar. ccam kendi betiğiyle, bir `.deb` paketiyle veya `go install` ile kurulur ve sağlama toplamlarını imza olmadan denetler.

## Hangisi size uyar?

Aynı anda tek bir hesap kullanıyorsanız ya da çalışan bir oturumun, limite ulaşmadan bile, yeniden başlatılmadan hesap değiştirmesini istiyorsanız giriş değiştiren bir araç seçin. Linux, macOS veya Windows üzerinde hesapları yan yana ve tek bir ortak kurulumla kullanmak, aynı konuşmaya sıradaki hesapta devam etmek ve varsayılan olarak Claude Code'un yanında hiçbir şeyin çalışmamasını istiyorsanız ccam'ı seçin. Hiçbir şeyin paylaşılmasına ihtiyacınız yoksa [düz kabuk takma adları](/compare/aliases) yeterlidir.

Birden fazla hesap kullanıp kullanamayacağınız, her birine uygulanan koşullara bağlıdır. Başlamadan önce bu koşulları okuyun.

**Team ve Enterprise koltukları.** Otomatik hesap seçimi (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar ve bu, Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce [SSS'yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) okuyun ve kuruluşunuza ve Anthropic'e danışın.
