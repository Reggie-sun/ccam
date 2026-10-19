---
title: Kullanım limitleri
description: ccam usage tüm hesapların 5 saatlik ve haftalık limitlerini aynı anda gösterir, ccam statusline ise hesap adını ve kullanımı Claude Code içine yerleştirir.
updated: 2026-10-19
order: 7
translation: machine
sourceHash: 8306a8f6d950
---

`ccam usage`, her hesabın 5 saatlik ve haftalık limitlerinin ne kadarını kullandığını, her pencerenin sıfırlanmasına kalan süreyle birlikte gösterir. Uzun bir işe başlamadan önce veya Claude Code bir limitin yaklaştığını söylediğinde, hangi hesapta pay kaldığını görmek için kullanın. `ccam statusline`, Claude Code için durum satırı komutudur: hesap adını ve kullanımı Claude Code içinde gösterir ve sayıları `ccam usage` için kaydeder.

## Kullanım

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Terminalde tablo, her satırda `fetching…` ile hemen görünür ve her satır, hesabı yanıt verdikçe yeniden çizilir. `As of` sayıların ne kadar eski olduğunu söyler; `(cached)` sayıların uç noktadan değil, kayıtlı bir dosyadan geldiği anlamına gelir. Son satır, 5 saatlik kullanımı %100'ün altında olan hesaplar arasında en düşük olanı adlandırır ve yalnızca birden fazla hesabınız olduğunda görünür. Giriş yapılmamış bir hesap `not logged in → ccam login <name>`, API anahtarlı bir hesap ise `no limits (API)` gösterir.

`ccam statusline install` sonrasında Claude Code, `[work] · Opus · 5h 24% · 7d 41%` gibi bir satır gösterir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--json` | Raporları JSON olarak yazdırır: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | kapalı |
| `--watch[=N]` | Ctrl+C'ye basana kadar her `N` saniyede bir yeniden çizer. Tek başına `--watch` veya 0 ya da daha küçük bir değer 60 anlamına gelir | kapalı |
| `-v`, `--verbose` | Modele özgü haftalık limitleri ve bir `source:` satırını (`live`, `cache` veya `statusline`) gösterir | kapalı |
| `--no-cache` | 60 saniyelik önbelleği yok sayar ve uç noktaya yeniden sorar | kapalı |
| `--source live\|cache` | `cache` uç noktaya hiç sormaz ve yalnızca kayıtlı sayıları gösterir; `live` varsayılan davranıştır | live |
| `--refresh` | Süresi dolmuş erişim tokenlarını yalnızca çalışmayan hesaplar için yeniler | kapalı |
| `--dry-run` | Kabul edilir: `ccam usage` kendi önbelleği dışında hiçbir şeyi değiştirmez. `--refresh` ile birlikte reddedilir: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | kapalı |

Başka herhangi bir `--source` değeri `--source must be live or cache` ile başarısız olur. `ccam usage` argüman almaz.

`ccam statusline`, Claude Code'un durum satırı JSON'unu standart girdiden okur. Tek seçeneği `--wrap <command>` seçeneğidir: ilk çıktı satırının başına hesap adı eklenen mevcut bir durum satırı komutu. `ccam statusline install` seçenek almaz. İkisi de `--dry-run` seçeneğini `` --dry-run is not supported by `ccam statusline`; nothing was changed `` ile reddeder.

## Nasıl çalışır

Sayılar iki yerden gelir: her hesabın kendi girişiyle çağrılan, Anthropic'in belgelenmemiş kullanım uç noktası `GET https://api.anthropic.com/api/oauth/usage` ve Claude Code'un durum satırı komutuna verdiği değerler. Bunların hepsi `internal/usage` içindedir.

1. `cli.usageCmd`, `config.toml` dosyasını yükler ve her hesap için bir `usage.Account` oluşturur: klasörü `~/.local/share/ccam/profiles/<name>/`, oturum klasörü (çalışıp çalışmadığını görmek için) ve yerinde içe aktarılıp aktarılmadığı.
2. `Fetcher.FetchAll` her hesap için bir goroutine başlatır. Hepsi tek bir `http.Client` paylaşır; böylece uç noktaya yapılan bağlantılar yeniden kullanılır. Bağlanma 2 sn sürebilir ve her hesaba toplam 3 sn verilir. İstemci yönlendirmeleri izlemez, çünkü bir yönlendirme tokenı başka bir yere yeniden gönderir; yanıtın da en fazla 1 MB'ını okur.
3. Her hesap bittikçe satırı yeniden çizilir. Başarısız olan veya zamanında yanıt vermeyen bir hesap kayıtlı sayılara geri döner; böylece yavaş bir hesap diğerlerini hiçbir zaman bekletmez.

Her hesap uç noktaya karşı kendi goroutine'inde çalışır; zaman aşımına uğrayan hesap bunun yerine kayıtlı sayılarını gösterir.

<!-- motion:usage-fetch -->

`personal` ve `work` zamanında yanıt verir ve canlı sayılar gösterir. `side` 3 sn'den uzun sürer; bu yüzden `fallback()` iki önbellek dosyasını okur ve en son yazılanı kullanır. `usage.Table()` satırları çizer ve `usage.Best()` son satırda adı geçen hesabı seçer.

Tek bir hesap için `Fetcher.Fetch` şu kontrolleri sırayla yapar:

<!-- motion:usage-cache -->

1. `--source cache` ile doğrudan `fallback()` fonksiyonuna gider.
2. `--no-cache` ayarlı değilse, `~/.local/state/ccam/usage/<name>.json` içinde 60 saniyeden (`CacheTTL`) yeni olan canlı bir sonuç olduğu gibi döndürülür. `ccam usage` komutunu art arda iki kez çalıştırmak hesap başına bir isteğe mal olur.
3. `readCredentials()`, `.credentials.json` dosyasını (veya macOS Anahtar Zinciri'ni) değiştirmeden okur ve yalnızca `claudeAiOauth` anahtarının kodu çözülür. Abonelik girişi yoksa `not logged in`, `.claude.json` bir API anahtarı tutuyorsa `no limits (API)` gösterilir.
4. Süresi dolmuş bir erişim tokenı kullanılmaz. `--refresh` olmadan hesap, süresi dolmuş olarak işaretlenmiş kayıtlı sayılarını gösterir.
5. Aksi halde ccam uç noktayı çağırır. 401 veya 403 süresi dolmuş bir giriş sayılır, diğer her durum kodu veya okunamayan bir gövde hata sayılır ve her ikisi de kayıtlı sayılara geri döner.
6. İyi bir yanıt, ölçüldüğü zamanla birlikte `<name>.json` dosyasına yazılır.

`fallback()` hem `<name>.json` hem de `<name>.statusline.json` dosyasını okur ve daha yeni olanı tutar. Kayıtlı sayıları göstermeden önce her pencerenin sıfırlanma zamanını kontrol eder: sayılar kaydedildikten sonra sıfırlanmış bir pencere %0 gösterir, çünkü eski yüzde artık geçerli değildir. Yukarıdaki `personal` hesabının sıfırlanma zamanı olmadan `0%` göstermesinin nedeni budur.

Uç noktanın iki yanıt biçimi vardır. Yenisinde `session` (5 saatlik), `weekly_all` ve `weekly_scoped` (modele özgü) girdileri içeren bir `limits` dizisi bulunur. Eskisinde `five_hour`, `seven_day`, `seven_day_opus` ve `seven_day_sonnet` nesneleri bulunur. `ParseLive` ikisini de okur: `limits` dizisi önceliklidir, eski nesneler ise onun kapsamadığını tamamlar. Hiç limit içermeyen bir gövde, değişmiş bir biçim olarak bildirilir ve `ccam doctor` bunu size söyler.

### Tokenları yenileme

Yenileme tokenları döner: birini yenilemek size yeni bir yenileme tokenı verir ve eskisini geçersiz kılar. Claude Code aynı hesapla çalışırken ccam girişi yenileseydi o oturumdan çıkış yapılırdı. Bu yüzden `ccam usage` tokenları okur ve siz `--refresh` vermedikçe hiçbir zaman yazmaz. O durumda `Fetcher.refresh`:

1. `ccam import` ile yerinde kaydedilmiş hesapları reddeder, çünkü düz `claude` o girişi ccam'in haberi olmadan kullanıyor olabilir;
2. macOS Anahtar Zinciri'nde saklanan girişleri reddeder;
3. iki ccam sürecinin aynı tokenı hiçbir zaman iki kez yenilememesi için bir kilit dosyası alır: `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (o değişken ayarlı değilse `~/.local/state/ccam/locks/` altında);
4. kilit altındayken hesabın bir oturumu çalışıyorsa reddeder ve token başka bir süreç tarafından zaten yenilenmişse isteği atlar;
5. yenileme tokenını token uç noktasına gönderir ve yeni tokenları, dosyadaki diğer tüm anahtarları koruyarak `fsutil.AtomicWrite` ile `.credentials.json` dosyasına geri yazar.

### Durum satırı toplayıcısı

Claude Code durum satırı komutunu çalıştırır ve ona boruyla bir JSON belgesi gönderir. ccam bundan model adını ve `rate_limits` değerini okur:

<!-- motion:usage-statusline -->

1. `runStatusline` hesabı, ccam'in Claude Code'u başlatırken ayarladığı `CCAM_PROFILE` değişkeninden, yoksa `CLAUDE_CONFIG_DIR` değerini hesap klasörleriyle eşleştirerek bulur.
2. `ParseStatusLine`, `rate_limits.five_hour` ve `rate_limits.seven_day` değerlerini okur. API anahtarlı hesaplarda veya oturumun ilk yanıtından önce JSON'da hız limitleri yoktur; bu durumda hiçbir şey kaydedilmez.
3. `SaveStatusLine`, `~/.local/state/ccam/usage/<name>.statusline.json` dosyasını yazar. `ccam usage`, uç nokta yanıt vermediğinde bunu kullanır.
4. ccam `[name] · model · 5h N% · 7d N%` yazdırır. `--wrap` ile komutunuzu aynı girdiyle `sh -c` üzerinden çalıştırır ve bunun yerine `[name]` ve ardından komutun ilk satırını yazdırır.

Sarmalanan komuta 2 saniye verilir. Sonrasında sonlandırılır ve satır yalnızca hesap adını gösterir; böylece yavaş bir komut Claude Code'un durum satırını engelleyemez.

`ccam statusline install`, her hesabın `settings.json` dosyasına `statusLine` yazar. Bir hesabın zaten bir durum satırı komutu varsa ccam onu `--wrap '<old command>'` olarak korur. Bir grup üzerinden paylaşılan `settings.json` bağlantısı aracılığıyla bir kez yazılır; böylece gruptaki her hesap onu alır. Yeniden çalıştırmak, zaten ccam kullanan dosyalar için `(already set)` yazdırır.

## Güvenlik ve uç durumlar

- `ccam usage`, siz `--refresh` vermedikçe hiçbir zaman bir token dosyasını değiştirmez. Bir girişin süresi dolduğunda hesabı bir kez `ccam <name>` ile çalıştırın, Claude Code onu yeniler.
- `--refresh` çalışan hesaplara, içe aktarılmış hesaplara veya Anahtar Zinciri girişlerine hiçbir zaman dokunmaz. Reddedilen bir yenileme, kayıtlı sayıları ve nedeni gösterir.
- Kullanım uç noktası herkese açık olarak belgelenmemiştir ve değişebilir. Değiştiğinde kayıtlı sayıları yine görürsünüz ve `ccam doctor` biçim değişikliğini bildirir.
- Kayıtlı sayı ve yanıt olmadığında satır `lookup failed:` ve nedeni gösterir, örneğin `--source cache` ile `no saved usage`.
- Önbellek dosyaları yalnızca yüzdeleri ve sıfırlanma zamanlarını tutar, hiçbir zaman token tutmaz.

## İlgili

- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [Sonraki hesaba geçme](/features/next)
- [Kullanım limiti iş akışı](/guides/usage-limit-workflow)
- [Komutlar](/docs/commands)
