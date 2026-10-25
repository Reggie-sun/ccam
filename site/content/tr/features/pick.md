---
title: Betikler için hesap seçme
description: ccam pick bir betik için en çok payı kalan hesabı yazdırır, birlikte başlatılan süreçlere farklı hesaplar verir ve hiçbirinde pay yoksa 75 ile çıkar.
updated: 2026-10-25
order: 12
translation: machine
sourceHash: "5c4f74db9714"
---

`ccam pick`, bir betiğin Claude Code'u başlatması için en çok payı kalan girişli hesabın adını yazdırır. Aynı anda başlatılan çalışan süreçler, payı kalan yeterli hesap olduğu sürece her biri farklı bir hesap alır; böylece bir grup `claude -p` işi tek bir hesabı tüketmek yerine hesaplarınıza dağılır.

**Team ve Enterprise koltukları.** Otomatik seçim (`ccam next`, `--fallback`, `ccam pick`, kontrol panelindeki `n` tuşu ve açtığınızda [`auto_switch`](/features/next#auto-switch)), aynı Claude Team veya Enterprise kuruluşunun koltukları arasında da geçiş yapar. Anthropic'in yardım merkezinde şöyle belirtilir: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team planındaki her koltuk, adı belirli tek bir kişiye atanır ve onu yalnızca o kişi kullanabilir) ve "Users can't route their work through another user's account to use that user's remaining usage" (kullanıcılar, başka bir kullanıcının kalan kullanım hakkından yararlanmak için işlerini o kullanıcının hesabı üzerinden yürütemez) ([Team planları](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise planları](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) için de durum aynıdır). Bu nedenle otomatik seçimi koltuklar arasında kullanmak Anthropic'in koşullarını ihlal edebilir. Orada kullanıp kullanmamak sizin kararınızdır: buna güvenmeden önce kuruluşunuza ve Anthropic'e danışın. Kişisel planlar (Pro, Max), [SSS'deki koşullar notunda](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ele alınır.

## Kullanım

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Çıkış kodları: standart çıktıda bir adla `0`, hiçbir hesapta pay kalmadığında veya hiçbir hesapta giriş yapılmadığında `75` (`EX_TEMPFAIL`), hatalı bir bayrak gibi diğer hatalar için `1`.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--json` | Hesabı, 5 saatlik ve haftalık kullanımını, sayıların nereden geldiğini ve neden seçildiğini yazdırır | kapalı |
| `--model <model>` | O modelin haftalık penceresini de sayar, örneğin `opus` | yok |
| `--max-usage <1-100>` | Bu kullanımın altındaki hesapları tercih eder | `auto_max_usage`, yoksa 100 |
| `--strategy headroom\|expiring\|balance` | En düşük 5 saatlik kullanım, en yakın 5 saatlik sıfırlanma ya da kullanılmadan kalacak ilk haftalık kota | `auto_strategy`, yoksa `headroom` |
| `--cross-org` | Başka kuruluşların hesaplarını da seçer | `cross_org`, yoksa kapalı |
| `--wait[=<duration>]` | Hiçbir hesapta pay yoksa ilk hesabın sıfırlanmasını en fazla bu süre kadar bekler. Değer `=` işaretinden sonra yazılır: `--wait 30m`, `--wait takes its value after '=': --wait=30m` ipucuyla reddedilir | kapalı; tek başına verilirse `5h` |
| `--dry-run` | Seçimi kaydetmeden seçer | kapalı |

## Nasıl çalışır

1. ccam, bu klasörün başlatacağı hesabın kuruluşundaki girişli hesapları listeler (API anahtarlı hesapların ve Console girişlerinin limiti yoktur ve dışarıda bırakılır; bkz. [kuruluşlar](/features/next#auto)); o hesapla aynı Claude hesabına giriş yapmış bir klasör de dışarıda bırakılır.
2. Bir kilit alır (`~/.local/state/ccam/picks.lock`); böylece aynı anda çalışan seçimler birbiri ardına gerçekleşir.
3. Her hesabın kullanımını, `ccam usage` ile aynı şekilde ve 60 saniyelik önbellekle, aynı anda sorar. Claude Code'un dolu bildirdiği bir pencere sıfırlanana kadar %100 sayılır ([limit kaydı](/features/usage#the-limit-record)); böylece limitine yeni ulaşmış bir hesap, kayıtlı sayıları pay gösterse bile dışarıda kalır.
4. Son 2 dakika içinde `ccam pick` tarafından ya da [`auto_switch`](/features/next#auto-switch) bir oturumu oraya taşıyarak seçilen hesaplar (`~/.local/state/ccam/picks/<name>`; saatin gösterdiği zamandan sonraya tarihli bir seçim sayılmaz) önce dışarıda bırakılır ve kalanlar arasından `ccam next --auto` komutunun [seçim kuralları](/features/next#auto) seçer: önce `auto_max_usage` altındakiler, stratejiye göre, %100'deki bir hesap asla seçilmeden. Bu, kuralların ilk seçeceği hesabı atladığında `--json` çıktısının `reason` alanı bunu söyler: `30% of its 5 hours used; work was picked a moment ago`.
5. Payı kalan her hesap yakın zamanda seçildiyse kurallar yeniden hepsi arasından seçer; böylece üç hesaptaki dördüncü çalışan süreç bir hesabı paylaşır.
6. Seçim kaydedilir, kilit bırakılır, ardından ad yazdırılır.

2 dakikalık pencere önemlidir, çünkü yeni bir işin kullanımının sayılara yansıması biraz zaman alır: o olmasa birlikte başlatılan beş çalışan süreç, en çok payı kalan hesap olarak hep aynı hesabı görürdü.

## Güvenlik ve uç durumlar

- `ccam pick` yalnızca kullanımı okur ve seçim zamanını yazar; `ccam usage` komutunun yaptığının ötesinde hiçbir zaman token yenilemez veya okumaz.
- Bir seçim hesabı ayırmaz. İki betik 2 dakika sonra yine aynı hesapta başlayabilir; o hesap limitine ulaştığında toparlanmak için `--fallback` kullanın.
- `--wait` ile beklemek kilit tutmaz; bu sırada diğer seçimler devam eder.
- `ccam pick`, [`ccam folder`](/features/folder) ayarını yok sayar: yalnızca kullanıma göre seçer.

## İlgili

- [Sonraki hesaba geçme](/features/next)
- [Limitte yeniden çalıştırma](/features/next#run-again-at-the-limit)
- [Kullanım limitleri](/features/usage)
- [API anahtarlı hesaplar](/features/api)
