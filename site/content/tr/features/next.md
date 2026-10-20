---
title: Sonraki hesaba geçme
description: ccam next, Claude Code'u son hesaptan sonrakiyle açar; --auto en çok payı kalanı seçer, --continue sürdürür, --fallback -p işlerini yeniden çalıştırır.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next`, Claude Code'u adını vermenize gerek kalmadan başka bir hesapla açar. Bir hesap kullanım limitine ulaştığında kullanın: düz `ccam next` en son kullandığınız hesaptan sonraki hesabı alır, `--auto` 5 saatlik penceresinde en çok pay kalan hesabı alır, `--continue` ise son konuşmayı yeni hesapta sürdürür. `claude -p` çalıştıran betikler için `ccam --fallback` bu geçişi kendiliğinden yapar (bkz. [Limitte yeniden çalıştırma](#run-again-at-the-limit)).

## Kullanım

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam geçişi standart hataya `→ <from> → <to>` olarak yazdırır, ardından Claude Code'u `ccam <name>` ile aynı şekilde yeni hesapla başlatır. `<from>`, en son kullandığınız hesaptır. `--` sonrasındaki argümanlar `claude` komutuna değiştirilmeden iletilir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--auto` | %100'deki hesapları atlayarak girişli hesabı kullanıma göre seçer (aşağıya bakın) | kapalı |
| `--max-usage <1-100>` | `--auto` ile: bu kullanımın altındaki hesapları tercih eder | `auto_max_usage`, yoksa 100 |
| `--strategy headroom\|expiring` | `--auto` ile: en düşük 5 saatlik kullanım ya da en yakın 5 saatlik sıfırlanma | `auto_strategy`, yoksa `headroom` |
| `--wait[=<duration>]` | `--auto` ile: hiçbir hesapta pay yoksa ilk hesabın sıfırlanmasını en fazla bu süre kadar bekler | kapalı; tek başına verilirse `5h` |
| `-c`, `--continue` | Yalnızca en son kullandığınız hesapla konuşma geçmişini paylaşan hesapları dikkate alır ve `claude --continue` başlatır | kapalı |
| `-- <claude args>` | ccam'in kendi argümanlarından sonra `claude` komutuna iletilen argümanlar | yok |
| `--dry-run` | Dikkate alınır: Claude Code'u başlatmadan veya bir başlatma kaydetmeden `→ <from> → <to>` ve çalıştıracağı komut satırını yazdırır | kapalı |

`--auto` ve `--continue` birlikte çalışır. `--continue` ile `--continue`, sizin argümanlarınızdan önce yerleştirilir.

## Nasıl çalışır

ccam üzerinden yapılan her başlatma, zamanı `~/.local/state/ccam/last-used/<name>` dosyasına yazar. `ccam next` bu dosyaları okur, hangi hesaplarda giriş yapıldığını kontrol eder ve seçer.

1. `cli.candidates`, `config.toml` dosyasından her hesap için bir `next.Candidate` oluşturur: hesap, son başlatma zamanı (`launcher.LastUsed`) ve girişi olup olmadığı (`creds.Has`: okunabilir bir `.credentials.json` ya da macOS'ta bir Anahtar Zinciri öğesi).
2. `--auto` olmadan `next.Rotate`, hesapları yuva numarasına göre sıralar ve en son başlatılanı bulur. Bu, kendisinden geçiş yaptığınız hesaptır. Henüz hiçbir hesap başlatılmadıysa bu, en küçük yuvadır.
3. Oradan yuva sırasıyla ileri doğru ilerler, son yuvadan ilk yuvaya geri sarar ve girişi olan ilk hesabı döndürür.
4. `launcher.LaunchWith`, seçilen hesabı tam olarak `ccam <name>` gibi hazırlar ve başlatır: hesabın bağlantılarını kontrol eder, MCP paylaşımını uygular, oturumu kaydeder ve `syscall.Exec` ile ccam'i `claude` ile değiştirir.

İlerleme en son kullanılan hesaptan başlar ve ondan sonraki ilk girişli hesapta durur:

<!-- motion:next-rotate -->

En son `work` başlatıldığı için ilerleme ondan sonra başlar. `side` hesabında giriş yoktur ve atlanır, `lab` seçilir. `lab` hesabından da çıkış yapılmış olsaydı ilerleme başa sararak `personal` hesabına giderdi (diyagramdaki `% len`). Döndürme hiçbir zaman başladığınız hesabı döndürmez: girişli başka bir hesap yoksa bir hatayla durur.

### --auto

`--auto`, yuva sırası yerine kullanıma göre karar verir. `cli.pickAuto` girişli hesapları tutar, `ccam usage` komutunun kullandığı `Fetcher.FetchAll` ile hepsinin kullanımını aynı anda sorar ve seçimi `usage.Policy.Choose()` fonksiyonuna bırakır:

<!-- motion:next-auto -->

Bir hesap 5 saatlik ve haftalık pencerelerine göre, Claude Code argümanları arasında `--model` varsa o modelin haftalık penceresine göre de değerlendirilir (adında `opus`, `sonnet`, `haiku` veya `fable` geçen). `Policy.Choose()` hesapları üç gruba ayırır:

1. **Payı olan**: her pencere `auto_max_usage` (`config.toml`, varsayılan 100) altında. Aralarından strateji seçer: `headroom` (varsayılan) en düşük 5 saatlik kullanımı alır, `expiring` en erken sıfırlanacak 5 saatlik pencereyi alır; böylece kalan kotası kaybolmadan kullanılır.
2. **Eşiği aşan**: %100'ün altında ama `auto_max_usage` değerinde veya üstünde. Yalnızca hiçbir hesapta pay yoksa, en düşük 5 saatlik kullanımdan başlayarak kullanılır.
3. **Dolu**: %100'de bir pencere ya da bilinmeyen 5 saatlik kullanım. Hiçbir zaman seçilmez.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

Hiçbir şey seçilemediğinde ccam 75 koduyla (`EX_TEMPFAIL`) çıkar; böylece bir betik "daha sonra deneyin" durumunu bir hatadan ayırt edebilir. `--wait` ile bunun yerine ilk dolu hesabın engelleyen penceresi sıfırlanana kadar bekler (uç nokta yeni pencereyi gösterdiğinde, sıfırlanma zamanından 30 saniye sonra), `ccam: no account has room; waiting until 15:10` yazdırır ve yeniden seçer; hiçbir hesap ne zaman sıfırlanacağını söylemiyorsa 10 dakika sonra yeniden bakar. Bekleme bitmeden hiçbir şey boşalmazsa 75 ile çıkar. Sorgu `ccam usage` kurallarını izler: 60 saniyeden yeni sayılar yeniden kullanılır, her hesabın 3 saniyesi vardır ve yanıt vermeyen veya girişinin süresi dolmuş bir hesap kayıtlı sayılarına göre değerlendirilir. Sıfırlanma zamanı geçmiş kayıtlı bir pencere %0 sayılır; böylece penceresi o zamandan beri sıfırlanmış bir hesap atlanmaz. `--auto` hiçbir zaman token yenilemez. Düz döndürmenin aksine `--auto`, hâlâ en çok payı varsa en son kullandığınız hesabı seçebilir.

### --continue

Claude Code'un `--continue` seçeneği geçerli klasördeki en son konuşmayı açar; bu yalnızca yeni hesap o konuşmayı görebiliyorsa çalışır. Bu yüzden `--continue` ile ccam yalnızca en son kullandığınız hesapla aynı `history` grubundaki hesapları dikkate alır, ardından `--continue` değerini argümanlarınızın önüne koyar. Geçmişi `local` olan hesaplar kimseyle paylaşmaz; bu yüzden böyle bir hesaptan `--continue` her zaman başarısız olur. Önce `ccam share <name> --history on` ile paylaşımı açın ([paylaşım](/docs/sharing) sayfasına bakın).

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
2. Standart girdi bir kez okunur ve her denemeye verilir.
3. ccam, `claude` komutunu `exec` yerine bir alt süreç olarak başlatır, standart çıktısını bekletir ve standart hatayı geldikçe iletir.
4. `claude` 0 ile çıktığında ya da başka bir nedenle başarısız olduğunda bekletilen çıktı ve çıkış kodu değiştirilmeden iletilir.
5. Claude Code'un mesajı bunu söylediğinde (`You've hit your limit`, `usage limit reached`) ya da hesabın kullanım sorgusu 5 saatlik veya haftalık penceresini %100'de gösterdiğinde, hata kullanım limiti sayılır.
6. Bu durumda o denemenin çıktısı atılır ve aynı komut başka bir girişli hesapla çalışır: önce 5 saatlik kullanımı en düşük olan, ardından kullanımı sorgulanamayan hesaplar, numara sırasıyla. %100'deki hesaplar atlanır ve her hesap en fazla bir kez denenir.
7. Argümanlar arasında `-c`, `--continue`, `-r` veya `--resume` varsa yalnızca aynı geçmiş grubundaki hesaplar denenir, böylece konuşma bulunabilir.

Sonraki hesap `--auto` ile aynı kurallarla seçilir (`auto_max_usage`, `auto_strategy`, `--model` penceresi). Hiç hesap kalmadığında ccam son denemenin çıktısını yazdırır, hiçbir girişli hesapta pay kalmadığını söyler ve 75 koduyla çıkar. Hesap adından önce `--wait` ile (`ccam --fallback --wait=2h work -p …`) bunun yerine ilk hesabın sıfırlanmasını bekler, ardından tüm hesapları yeniden deneyebilir; limite yeni ulaşmış bir hesap hemen yeniden çalıştırılmasın diye önce her zaman en az bir kez bekler. `ccam --dry-run --fallback work -p …` ilk denemenin komutunu yazdırır.

Etkileşimli oturumlar kendiliğinden geçiş yapmaz. Bunun yerine `ccam statusline install` ile durum satırı, 5 saatlik veya haftalık kullanım %90'a ulaştığında `→ ccam next --auto --continue` ekler ([Kullanım limitleri](/features/usage)).

## Güvenlik ve uç durumlar

- Girişi olmayan hesaplar hiçbir zaman seçilmez. Girişli başka bir hesap yoksa düz `ccam next`, `no logged-in account other than <name>` ile durur.
- `--continue` ile ve geçmiş paylaşan hesap yoksa: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- `--auto` ile ve girişli hesap yoksa: `no logged-in account to pick`. Her hesap doluysa veya kullanımı bilinmiyorsa: `` no account with known usage and headroom (see `ccam usage`) ``. İkisi de 75 çıkış koduyla biter.
- `--max-usage`, `--strategy` ve `--wait` yalnızca `--auto` ile çalışır; 1-100 dışındaki bir değer veya başka bir strateji reddedilir.
- `--auto`, `--dry-run` ile bile ağ istekleri yapar. Yalnızca okur ve token dosyalarına hiçbir zaman yazmaz.
- Seçilen hesap `ccam <name>` ile aynı yoldan başlar; bu yüzden aynı geçmişi paylaşan başka bir Claude Code bu klasörde zaten çalışırken verilen uyarı dahil aynı kontroller uygulanır.
- `last-used` yalnızca Claude Code gerçekten başladığında yazılır. `--dry-run` ona dokunmaz; bu yüzden kuru çalıştırma döndürmeyi ilerletmez.
- `--fallback` standart çıktıyı bir deneme bitene kadar beklettiği için `--output-format stream-json` gibi akış çıktısı sonda tek parça halinde gelir. İlk denemenin limite ulaşmadan önce değiştirdiği dosyalar değişmiş olarak kalır; sonraki deneme istemi baştan başlatır.

## İlgili

- [Kullanım limitleri](/features/usage)
- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Kullanım limiti iş akışı](/guides/usage-limit-workflow)
