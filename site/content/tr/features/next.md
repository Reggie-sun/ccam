---
title: Sonraki hesaba geçme
description: ccam next, Claude Code'u en son kullandığınız hesaptan sonraki hesapla açar; --auto en düşük 5 saatlik kullanımı seçer, --continue konuşmayı sürdürür.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next`, Claude Code'u adını vermenize gerek kalmadan başka bir hesapla açar. Bir hesap kullanım limitine ulaştığında kullanın: düz `ccam next` en son kullandığınız hesaptan sonraki hesabı alır, `--auto` 5 saatlik penceresinde en çok pay kalan hesabı alır, `--continue` ise son konuşmayı yeni hesapta sürdürür.

## Kullanım

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
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
| `--auto` | %100'deki hesapları atlayarak, girişli hesaplar arasından 5 saatlik kullanımı en düşük olanı seçer | kapalı |
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

`--auto`, yuva sırası yerine kullanıma göre karar verir. `cli.pickAuto` girişli hesapları tutar, `ccam usage` komutunun kullandığı `Fetcher.FetchAll` ile hepsinin kullanımını aynı anda sorar ve seçimi `usage.Best()` fonksiyonuna bırakır:

<!-- motion:next-auto -->

`usage.Best()` en düşük 5 saatlik yüzdeyi alır ve %100'deki veya 5 saatlik kullanımı bilinmeyen hesapları atlar. Sorgu `ccam usage` kurallarını izler: 60 saniyeden yeni sayılar yeniden kullanılır, her hesabın 3 saniyesi vardır ve yanıt vermeyen veya girişinin süresi dolmuş bir hesap kayıtlı sayılarına göre değerlendirilir. Sıfırlanma zamanı geçmiş kayıtlı bir pencere %0 sayılır; böylece penceresi o zamandan beri sıfırlanmış bir hesap atlanmaz. `--auto` hiçbir zaman token yenilemez. Düz döndürmenin aksine `--auto`, hâlâ en çok payı varsa en son kullandığınız hesabı seçebilir.

### --continue

Claude Code'un `--continue` seçeneği geçerli klasördeki en son konuşmayı açar; bu yalnızca yeni hesap o konuşmayı görebiliyorsa çalışır. Bu yüzden `--continue` ile ccam yalnızca en son kullandığınız hesapla aynı `history` grubundaki hesapları dikkate alır, ardından `--continue` değerini argümanlarınızın önüne koyar. Geçmişi `local` olan hesaplar kimseyle paylaşmaz; bu yüzden böyle bir hesaptan `--continue` her zaman başarısız olur. Önce `ccam share <name> --history on` ile paylaşımı açın ([paylaşım](/docs/sharing) sayfasına bakın).

## Güvenlik ve uç durumlar

- Girişi olmayan hesaplar hiçbir zaman seçilmez. Girişli başka bir hesap yoksa düz `ccam next`, `no logged-in account other than <name>` ile durur.
- `--continue` ile ve geçmiş paylaşan hesap yoksa: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- `--auto` ile ve girişli hesap yoksa: `no logged-in account to pick`. Her hesap doluysa veya kullanımı bilinmiyorsa: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto`, `--dry-run` ile bile ağ istekleri yapar. Yalnızca okur ve token dosyalarına hiçbir zaman yazmaz.
- Seçilen hesap `ccam <name>` ile aynı yoldan başlar; bu yüzden aynı geçmişi paylaşan başka bir Claude Code bu klasörde zaten çalışırken verilen uyarı dahil aynı kontroller uygulanır.
- `last-used` yalnızca Claude Code gerçekten başladığında yazılır. `--dry-run` ona dokunmaz; bu yüzden kuru çalıştırma döndürmeyi ilerletmez.

## İlgili

- [Kullanım limitleri](/features/usage)
- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Kullanım limiti iş akışı](/guides/usage-limit-workflow)
