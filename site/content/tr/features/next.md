---
title: Sonraki hesaba geçme
description: ccam next, Claude Code'u başka bir hesapla açar; --auto en çok payı kalanı seçer, --continue konuşmayı taşır, --fallback -p işlerini yeniden çalıştırır.
updated: 2026-10-24
order: 8
translation: machine
sourceHash: "a38daed5dd97"
---

`ccam next`, Claude Code'u adını vermenize gerek kalmadan başka bir hesapla açar. Bir hesap kullanım limitine ulaştığında kullanın: düz `ccam next` en son kullandığınız hesaptan sonraki hesabı alır, `--auto` en çok pay kalan hesabı alır, `--continue` ise bu klasörün son konuşmasını, iki hesap geçmişi paylaşsın ya da paylaşmasın, yeni hesaba taşır. Claude Code'un içinde [ccam eklentisi](#inside-claude-code) bir limit oturumu durdurduğunda size haber verir ve `/ccam:switch` konuşmayı devreder. `claude -p` çalıştıran betikler için `ccam --fallback` geçişi kendiliğinden yapar (bkz. [Limitte yeniden çalıştırma](#run-again-at-the-limit)).

**Team ve Enterprise koltukları.** Otomatik seçim (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar. Anthropic'in yardım merkezinde şöyle belirtilir: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team planındaki her koltuk, adı belirli tek bir kişiye atanır ve onu yalnızca o kişi kullanabilir) ve "Users can't route their work through another user's account to use that user's remaining usage" (kullanıcılar, başka bir kullanıcının kalan kullanım hakkından yararlanmak için işlerini o kullanıcının hesabı üzerinden yürütemez) ([Team planları](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise planları](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) için de durum aynıdır). Bu nedenle otomatik seçimi koltuklar arasında kullanmak Anthropic'in koşullarını ihlal edebilir. Buna güvenmeden önce kuruluşunuza ve Anthropic'e danışın. Kişisel planlar (Pro, Max), [SSS'deki koşullar notunda](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ele alınır.

## Kullanım

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam geçişi, bir `--auto` seçiminde nedeniyle birlikte, standart hataya `→ <from> → <to>` olarak yazdırır, ardından Claude Code'u `ccam <name>` ile aynı şekilde yeni hesapla başlatır. `<from>`, en son kullandığınız hesaptır; `--continue` ile ise konuşmanın ait olduğu hesaptır. `--` sonrasındaki argümanlar `claude` komutuna değiştirilmeden iletilir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--auto` | %100'deki hesapları atlayarak girişli hesabı kullanıma göre seçer (aşağıya bakın) | kapalı |
| `--max-usage <1-100>` | `--auto` ile: bu kullanımın altındaki hesapları tercih eder | `auto_max_usage`, yoksa 100 |
| `--strategy headroom\|expiring\|balance` | `--auto` ile: en düşük 5 saatlik kullanım, en yakın 5 saatlik sıfırlanma ya da kullanılmadan kalacak ilk haftalık kota | `auto_strategy`, yoksa `headroom` |
| `--cross-org` | `--auto` ile: başka bir kuruluşun hesaplarını da seçer (aşağıya bakın) | `cross_org`, yoksa kapalı |
| `--wait[=<duration>]` | `--auto` ile: hiçbir hesapta pay yoksa ilk hesabın sıfırlanmasını en fazla bu süre kadar bekler. Değer `=` işaretinden sonra yazılır: `--wait 2h`, `--wait takes its value after '=': --wait=2h` ipucuyla reddedilir | kapalı; tek başına verilirse `5h` |
| `-c`, `--continue` | Bu klasörün son konuşmasını yeni hesaba taşır (aşağıya bakın) | kapalı |
| `-- <claude args>` | ccam'in kendi argümanlarından sonra `claude` komutuna iletilen argümanlar | yok |
| `--dry-run` | Dikkate alınır: Claude Code'u başlatmadan veya bir başlatma kaydetmeden `→ <from> → <to>` ve çalıştıracağı komut satırını yazdırır | kapalı |

`--auto` ve `--continue` birlikte çalışır. `--continue` kullanıldığında `--resume` argümanları sizin argümanlarınızdan önce yerleştirilir.

Claude Code'un kendi seçenekleri, örneğin `--dangerously-skip-permissions` veya `--model`, `--` sonrasına yazılır. Adın ardından gelen her şeyi olduğu gibi ileten `ccam <name>` komutunun aksine `ccam next` argümanlarını kendisi okur; bu yüzden `--` olmadan `ccam: unknown flag: --dangerously-skip-permissions` hatasıyla durur.

## Nasıl çalışır

ccam üzerinden yapılan her başlatma, zamanı `~/.local/state/ccam/last-used/<name>` dosyasına yazar. `ccam next` bu dosyaları okur, hangi hesaplarda giriş yapıldığını kontrol eder ve seçer.

1. `cli.candidates`, `config.toml` dosyasından her hesap için bir `next.Candidate` oluşturur: hesap, son başlatma zamanı (`launcher.LastUsed`) ve çalışan bir girişi olup olmadığı (`creds.Has`: `.credentials.json` içinde ya da macOS'ta Anahtar Zinciri'nde, tümüyle süresi dolmamış ve Claude Code tarafından silinmemiş bir Claude girişi). Bunun yerine bir API anahtarı kullanan bir klasörün (bir Console girişi ya da `ccam add --api` ile eklenmiş bir hesap) plan limitleri yoktur ve hiçbir zaman seçilmez.
2. `--auto` olmadan `next.Rotate`, hesapları yuva numarasına göre sıralar ve en son başlatılanı bulur. Bu, kendisinden geçiş yaptığınız hesaptır. Henüz hiçbir hesap başlatılmadıysa bu, en küçük yuvadır.
3. Oradan yuva sırasıyla ileri doğru ilerler, son yuvadan ilk yuvaya geri sarar ve girişi olan ilk hesabı döndürür.
4. `launcher.LaunchWith`, seçilen hesabı tam olarak `ccam <name>` gibi hazırlar ve başlatır: hesabın bağlantılarını kontrol eder, MCP paylaşımını uygular, oturumu kaydeder ve `syscall.Exec` ile ccam'i `claude` ile değiştirir.

İlerleme en son kullanılan hesaptan başlar ve ondan sonraki ilk girişli hesapta durur:

<!-- motion:next-rotate -->

En son `work` başlatıldığı için ilerleme ondan sonra başlar. `side` hesabında giriş yoktur ve atlanır, `lab` seçilir. `lab` hesabından da çıkış yapılmış olsaydı ilerleme başa sararak `personal` hesabına giderdi (diyagramdaki `% len`). Döndürme hiçbir zaman başladığınız hesabı döndürmez: girişli başka bir hesap yoksa bir hatayla durur.

### --auto

`--auto`, yuva sırası yerine kullanıma göre karar verir. `cli.pickAuto` girişli hesapları tutar, `ccam usage` komutunun kullandığı `Fetcher.FetchAll` ile hepsinin kullanımını aynı anda sorar ve seçimi `usage.Policy.Choose()` fonksiyonuna bırakır:

<!-- motion:next-auto -->

Bir hesap 5 saatlik ve haftalık pencerelerine göre, Claude Code argümanları arasında `--model` varsa (`--` sonrasında) o modelin haftalık penceresine göre de değerlendirilir (adında `opus`, `sonnet`, `haiku` veya `fable` geçen). `Policy.Choose()` hesapları üç gruba ayırır:

1. **Payı olan**: her pencere `auto_max_usage` (`config.toml`, varsayılan 100) altında. Aralarından strateji seçer: `headroom` (varsayılan) en düşük 5 saatlik kullanımı alır; `expiring` en erken sıfırlanacak 5 saatlik pencereyi alır, böylece kalan kotası kaybolmadan kullanılır; `balance` haftası sıfırlanana kadar gün başına en çok haftalık kotası kalan hesabı alır ve çalışan her ccam oturumunu bu kotadan sayar; böylece haftalık kota kullanılmadan kalmaz ve paralel oturumlar dağılır.
2. **Eşiği aşan**: %100'ün altında ama `auto_max_usage` değerinde veya üstünde. Yalnızca hiçbir hesapta pay yoksa, en düşük 5 saatlik kullanımdan başlayarak kullanılır.
3. **Dolu**: %100'de bir pencere ya da bilinmeyen 5 saatlik kullanım. Hiçbir zaman seçilmez.

Claude Code orada bir isteği reddettiğinde bir pencere, sıfırlanana kadar %100 sayılır: [ccam eklentisinin](#inside-claude-code) hook'u bunu, kayıtlı sayılar ne derse desin, kaydeder (`~/.local/state/ccam/usage/<name>.limit.json`). Böylece limitine yeni ulaşmış bir hesap, sayıları eski olsa bile, sıfırlanmadan önce yeniden seçilmez.

Tüm hesaplar için ayarlar `config.toml` dosyasının en üstüne, ilk `[[group]]`, `[[profile]]` veya `[[folder]]` satırının üstüne yazılır. Bunlardan birinin altındaki bir ayar o tabloya aittir: ccam onu her başlatmada adıyla belirtir ve yapılandırmayı kaydeden komutlar ayar taşınana kadar çalışmayı reddeder.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Kuruluşlar.** Otomatik seçimler, bıraktığınız hesaba benzeyen hesaplarda kalır: aynı ekip veya kurumsal kuruluş ya da birlikte kişisel planlar (Pro, Max); böylece bir şirket hesabındaki iş kişisel bir hesapta sürmez, tersi de olmaz. ccam her hesabın kuruluşunu girişinden bilir (`ccam login` onu kaydeder; daha eski girişler `.claude.json` dosyasından okunur); kuruluşunu belirleyemediği bir hesap her yere uyar. Aynı Claude hesabına giriş yapmış iki klasör onun limitlerini paylaşır; bu yüzden biri hiçbir zaman diğerinin yerini almaz. `--cross-org` ya da `config.toml` içinde `cross_org = true`, kuruluş kuralını kaldırır.

Hiçbir şey seçilemediğinde ccam 75 koduyla (`EX_TEMPFAIL`) çıkar; böylece bir betik "daha sonra deneyin" durumunu bir hatadan ayırt edebilir. `--wait` ile bunun yerine ilk dolu hesabın engelleyen penceresi sıfırlanana kadar bekler (uç nokta yeni pencereyi gösterdiğinde, sıfırlanma zamanından 30 saniye sonra), `ccam: no account has room; waiting until 15:10` yazdırır ve yeniden seçer; hiçbir hesap ne zaman sıfırlanacağını söylemiyorsa 10 dakika sonra yeniden bakar. Bekleme bitmeden hiçbir şey boşalmazsa 75 ile çıkar. Sorgu `ccam usage` kurallarını izler: 60 saniyeden yeni sayılar yeniden kullanılır, her hesabın 3 saniyesi vardır, yanıt vermeyen veya erişim tokenının süresi dolmuş bir hesap kayıtlı sayılarına göre değerlendirilir ve bir limit kaydı bunların üzerine yerleşir. Sıfırlanma zamanı geçmiş kayıtlı bir pencere %0 sayılır; böylece penceresi o zamandan beri sıfırlanmış bir hesap atlanmaz. `--auto` hiçbir zaman token yenilemez. Düz döndürmenin aksine `--auto`, hâlâ en çok payı varsa en son kullandığınız hesabı seçebilir.

### --continue

1. ccam bu klasörün son konuşmasını her hesapta arar: `projects/<folder>/` içindeki `<session>.jsonl` dosyalarından (klasör adı, Claude Code'un adlandırdığı gibi, bağlantıları çözülmüş ve harf veya rakam dışındaki her karakteri `-` yapılmış çalışma klasörüdür) en yeni istemi veya yanıtı olan. Yalnızca `/mcp`, `/usage` veya `/status` çalıştırmış bir oturum ya da boş bir dosya gibi istemi veya yanıtı olmayan oturumlar atlanır. Bir oturum bir not bıraktıysa (bir kullanım limiti veya `/ccam:switch`, aşağıya bakın) ve hiçbir konuşmada nottan 2 dakikadan daha sonra bir istem veya yanıt yoksa bunun yerine tam olarak o oturum kullanılır; Claude Code'un `/exit` sırasında yazdığı satırlar sayılmaz.
2. Konuşmanın ait olduğu hesap, kendisinden geçtiğiniz hesaptır; yeni hesap ise yukarıdaki gibi, aynı kuruluş içinden seçilir. Birkaç hesap o geçmişi paylaşıyorsa konuşma, `.claude.json` dosyası konuşmayı en son kendisinin çalıştırdığını söyleyen (`lastSessionId`) hesaba, yoksa en son kullanılan hesaba aittir.
3. Yeni hesap o geçmişi paylaşıyorsa (aynı grup klasörü) ccam `claude --resume <session id>` başlatır, yani tam olarak o konuşmayı (`--continue` ise klasördeki en yeni konuşma hangisiyse onu açardı).
4. Aksi halde `claude --resume <the conversation's file> --fork-session` başlatır: Claude Code tüm konuşmayı yeni bir oturum kimliğiyle yeni hesabın kendi geçmişine kopyalar ve orijinal olduğu gibi kalır. Bundan sonra yeni hesabın `--continue` seçeneği onu bulur.

ccam konuşmanın ne kadar büyük olduğunu söyler: yeni hesap ilk yanıtıyla konuşmanın tamamını yeniden okur. İstem önbellekleri kuruluşlar arasında hiçbir zaman paylaşılmaz ve her kişisel plan kendi başına bir kuruluştur; bu yüzden iki hesap farklı kuruluşlara aitse ccam `prompt caches are not shared between organizations, so the first reply reads it all again` ekler. Eski hesap, yeni hesabın giriş yapmadığı MCP sunucularına giriş yapmışsa ccam onları adlandırır: `work is not signed in to MCP servers personal is: linear (sign in with /mcp there; sign-ins are not shared)`. Bir grupta `history` paylaşmak kopyalamayı önler, tek bir konuşma tutar ve `/rewind` komutunun kodu geri yüklediği kontrol noktalarını (`file-history/`) da paylaşır.

## Claude Code'un içinde

[ccam eklentisi](/docs/install#claude-code-plugin) iki hook ekler (Claude Code 2.1.78 veya sonrası):

- Bir tur kullanım limitinde durduğunda konuşmayı not eder ve destekleyen terminallerde (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) bir masaüstü bildirimi gösterir: ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` carries on as another account.`` Claude Code'un kullanım limiti saydığı her mesajı tanır: 5 saatlik, haftalık, düz ve model limitleri (`You've hit your session limit`, `You've hit your limit`, bir modelin kotası için `You've reached your …`), kullanım kredileri ve harcama limitleri ve hesabı durdurmuş bir kuruluş. Claude Code reddedilen isteği konuşmaya yazdıysa hook, dolu pencereyi sıfırlanana kadar da kaydeder ([limit kaydı](/features/usage#the-limit-record)); yöneticilerinin her an kaldırabileceği bir kuruluş durdurması bunun dışındadır. Kısa bir sunucu kısıtlaması veya bir API anahtarının hız sınırı kullanım limiti değildir ve kendi hâline bırakılır; bir alt ajanın durması ve başında kimsenin oturmadığı oturumlar (`-p`, arka plan oturumları) da öyle. Claude Code 2.1.234 ve sonrası ayrıca oturumda bekler ve 5 saatlik bir sıfırlanmadan sonra kendiliğinden devam eder; bildirim size bunun yerine geçiş yapma seçeneği verir.
- `/ccam:switch [account]` konuşmayı not eder ve herhangi bir istek yapılmadan önce turu bitirir: `/exit` yazın, ardından `ccam next --continue` konuşmayı adını verdiğiniz hesapla ya da sıradaki hesapla sürdürür. Henüz konuşması olmayan bir oturumda sürdürülecek bir şey olmadığını söyler.

Not, oturum başına bir dosya olarak `~/.local/state/ccam/handoff/<account>/<session>.json` içinde durur; böylece bir hesabın iki oturumu birbirinin notunun üzerine asla yazmaz. 12 saat saklanır ve bir sonraki oturum başlamak üzereyken bir kez kullanılır; başarısız olan bir başlatma onu bir sonraki deneme için bırakır. Not, Claude Code'un kabuğunda bir `cd` yapıldıktan sonra bile oturumun başladığı klasöre aittir. `ccam rename` içindeki adları günceller; notta adı geçen hesap o zamandan beri kaldırılmışsa ccam bunu söyler ve her zamanki gibi seçim yapar.

### Aktarıcı

`config.toml` dosyasının en üstünde `relay = true` olduğunda ccam, kendini Claude Code ile değiştirmek yerine etkileşimli oturumların yanında kalır ve `/ccam:switch` sonrasında o oturumun konuşmasını kendiliğinden sürdürür (başka bir oturumun bıraktığı nota göre işlem yapılmaz): hook Claude Code'u sonlandırır ve ccam, o oturumun notuyla `ccam next --continue` çalıştırır. Bir hesap adı verilmezse konuşmanın modelini de sayarak en çok payı kalanı alır. Windows'ta hook Claude Code'u sonlandıramaz: `/exit` yazın, ccam devam eder. Bunun bedeli oturum başına bekleyen bir ccam sürecidir; bu yüzden varsayılan olarak kapalıdır.

Aktarılan oturum, bir konuşmayı sürdürmenin geri getirmediklerini korur: `--add-dir`, `--mcp-config`, araç listeleri, `--settings`, `--plugin-dir` ve sistem istemi seçenekleri gibi seçenekleri, konuşmanın izin modunu ve oturumun efor düzeyini. `--model`, `--name` ve `--agent` seçeneklerini Claude Code konuşmadan kendisi geri getirir. Başka her şey (`--worktree`, `--debug`, `--remote-control`, ccam'in bilmediği seçenekler) `ccam: not carried over: …` içinde adlandırılır. Geçiş hiçbir zaman paylaşılan geçmiş sorusunu sormaz ve geçiş sırasında yazdıklarınız atılır.

Aktarıcı yalnızca Claude Code `/ccam:switch` isteğinin istediği şekilde sona erdiğinde devam eder: önce `/exit` yazdığınızda, sona ererken Ctrl+C'ye bastığınızda ya da ccam bir kapanma (hangup) veya başka bir sinyal aldığında devam etmez; terminal kapanmışsa da etmez. Not bu durumda kalır ve `ccam next --continue` konuşmayı daha sonra sürdürür. Daha eski bir ccam'in aktarıcısı altında başlatılmış bir oturumu ya da başında kimsenin oturmadığı bir oturumu (`-p`, arka plan oturumları) hook sonlandırmaz; bunun yerine `/exit` mesajını alır. `claude remote-control` aktarıcı olmadan çalışır.

## Limitte yeniden çalıştırma

`ccam next`, durmuş bir etkileşimli oturum içindir. Etkileşimsiz bir çalıştırma, yani bir betikte veya zamanlanmış bir işte `claude -p`, sizi bekleyemez; bu yüzden hesap adından önce gelen `--fallback`, ccam'in onun yerine geçiş yapmasını sağlar:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` yalnızca `-p` veya `--print` ile birlikte çalışır; aksi halde ccam durur ve `ccam next --auto --continue` komutunu gösterir.
2. Borudan gelen standart girdi en fazla 10 MiB'a kadar bir kez okunur ve her denemeye verilir. ccam, Claude Code gibi, girdinin başlamasını en fazla 3 saniye bekler (`--input-format stream-json` ile beklemez); açık kalan ve hiçbir şey göndermeyen bir boru, standart hataya bir satırla birlikte çalıştırmaya hiç girdi vermez.
3. ccam, `claude` komutunu `exec` yerine bir alt süreç olarak başlatır, standart çıktısını bekletir ve standart hatayı geldikçe iletir.
4. `claude` 0 ile çıktığında ya da başka bir nedenle başarısız olduğunda bekletilen çıktı ve çıkış kodu değiştirilmeden iletilir. Bir sinyalle sonlandırılan çalıştırma da (128'in üzerinde bir çıkış kodu) öyle: ona durması söylenmiştir ve yeniden çalıştırılmaz.
5. Claude Code'un mesajı bunu söylediğinde, yani [hook'un](#inside-claude-code) tanıdığı kullanım limiti mesajlarından herhangi biri, `--output-format json` veya `stream-json` çıktısının içinde de, ya da hesabın kullanımı 5 saatlik, haftalık veya `--model` penceresini %100'de gösterdiğinde (bir limit kaydı dahil), hata kullanım limiti sayılır.
6. Bu durumda o denemenin çıktısı atılır ve aynı komut başka bir girişli hesapla çalışır: önce 5 saatlik kullanımı en düşük olan, ardından kullanımı sorgulanamayan hesaplar, numara sırasıyla. %100'deki hesaplar atlanır, API anahtarıyla giriş yapan hesaplar hiçbir zaman denenmez ve her hesap en fazla bir kez denenir.
7. Argümanlar arasında `-c`, `--continue`, `-r` veya `--resume` varsa yalnızca aynı geçmiş grubundaki hesaplar denenir, böylece konuşma bulunabilir.
8. `--auto` ile olduğu gibi yalnızca aynı kuruluşun hesapları denenir (`cross_org = true` bunu kaldırır).

Sonraki hesap `--auto` ile aynı kurallarla seçilir (`auto_max_usage`, `auto_strategy`, `--model` penceresi). Hiç hesap kalmadığında ccam son denemenin çıktısını yazdırır, hiçbir girişli hesapta pay kalmadığını söyler ve 75 koduyla çıkar. Hesap adından önce `--wait` ile (`ccam --fallback --wait=2h work -p …`) bunun yerine ilk hesabın sıfırlanmasını bekler, ardından tüm hesapları yeniden deneyebilir; limite yeni ulaşmış bir hesap hemen yeniden çalıştırılmasın diye önce her zaman en az bir kez bekler. Süre, her beklemeyi ayrı ayrı değil, çalıştırmanın tüm beklemelerini birlikte sınırlar. Boşlukla yazıldığında (`--wait 30m`) `ccam next` ile aynı ipucuyla reddedilir. Claude Code'un argümanlarından önce gelen bir `--`, her başlatmada olduğu gibi çalışır (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` ilk denemenin komutunu yazdırır.

Etkileşimli oturumlar kendiliğinden geçiş yapmaz. Bunun yerine `ccam statusline install` ile durum satırı, 5 saatlik veya haftalık kullanım %90'a ulaştığında `→ ccam next --auto --continue` ekler ([Kullanım limitleri](/features/usage)).

## Güvenlik ve uç durumlar

- Girişi olmayan hesaplar hiçbir zaman seçilmez. Girişli başka bir hesap yoksa düz `ccam next`, `no logged-in account other than <name>` ile durur.
- `--continue` ile ve klasörde konuşma yoksa: `no account has a conversation in <folder> to continue`.
- Kuruluş kuralı hiç hesap bırakmadığında: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- `--auto` ile ve girişli hesap yoksa: `no logged-in account to pick`. Her hesap doluysa veya kullanımı bilinmiyorsa: `` no account with known usage and headroom (see `ccam usage`) ``. İkisi de 75 çıkış koduyla biter.
- `--max-usage`, `--strategy`, `--wait` ve `--cross-org` yalnızca `--auto` ile çalışır; 1-100 dışındaki bir değer veya başka bir strateji reddedilir.
- Bir çatallama konuşmayı kopyalar, hesabın diğer durumunu değil: geçmiş paylaşılmadıkça orijinal hesabın `file-history/` içindeki kontrol noktaları gelmez.
- `--auto`, `--dry-run` ile bile ağ istekleri yapar. Yalnızca okur ve token dosyalarına hiçbir zaman yazmaz.
- Seçilen hesap `ccam <name>` ile aynı yoldan başlar; bu yüzden aynı geçmişi paylaşan başka bir Claude Code bu klasörde zaten çalışırken verilen uyarı dahil aynı kontroller uygulanır.
- `last-used` yalnızca Claude Code gerçekten başladığında yazılır. `--dry-run` ona dokunmaz; bu yüzden kuru çalıştırma döndürmeyi ilerletmez.
- `--fallback` standart çıktıyı bir deneme bitene kadar beklettiği için `--output-format stream-json` gibi akış çıktısı sonda tek parça halinde gelir. `stream-json` girdisi akıtan ve daha fazlasını göndermeden önce yanıtları bekleyen bir istemci `--fallback` kullanamaz. İlk denemenin limite ulaşmadan önce değiştirdiği dosyalar değişmiş olarak kalır; sonraki deneme istemi baştan başlatır.

## İlgili

- [Kullanım limitleri](/features/usage)
- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Kullanım limiti iş akışı](/guides/usage-limit-workflow)
