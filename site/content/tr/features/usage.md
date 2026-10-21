---
title: Kullanım limitleri
description: ccam usage tüm hesapların 5 saatlik ve haftalık limitlerini aynı anda gösterir, ccam statusline ise hesap adını ve kullanımı Claude Code içine yerleştirir.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

`ccam usage`, her hesabın 5 saatlik limitin, haftalık limitin ve her modelin kendi haftalık limitinin (Fable gibi) ne kadarını kullandığını, her limit için bir çubuk ve sıfırlanmasına kalan süreyle gösterir. Uzun bir işe başlamadan önce veya Claude Code bir limitin yaklaştığını söylediğinde, hangi hesapta pay kaldığını görmek için kullanın. `ccam statusline`, Claude Code için durum satırı komutudur: hesap adını ve kullanımı Claude Code içinde gösterir ve sayıları `ccam usage` için kaydeder.

## Kullanım

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Her hesap; numarası, adı ve sayılarının ne kadar eski olduğunu gösteren bir satır alır; `(cached)` sayıların uç noktadan değil, kayıtlı bir dosyadan geldiği anlamına gelir. Altında her limit için bir satır, 20 hücrelik bir çubuk, yüzdeyi ve pencerenin sıfırlanmasına kalan süreyi gösterir: `5-hour`, `Weekly` (tüm modeller), ardından her modelin kendi haftalık penceresi, örneğin `Fable`. Hesabın sahip olmadığı bir pencere `-` gösterir. Yüzde %80'den itibaren sarıya, %100'de kırmızıya döner. Terminalde her hesap hemen `fetching…` olarak görünür ve satırları hesap yanıt verdikçe çizilir. Son satır, [`ccam next --auto`](/features/next) komutunun aynı kurallarla seçeceği hesabı adlandırır (`auto_max_usage`, `auto_strategy`, çalışan oturumlar, `cross_org` ayarlanmadıkça yalnızca son kullanılan hesabın kuruluşundaki hesaplar ve %100'deki her pencere atlanarak; böylece haftası dolmuş bir hesap hiçbir zaman adlandırılmaz) ve yalnızca birden fazla hesabınız olduğunda görünür. Giriş yapılmamış bir hesap `not logged in → ccam login <name>`, API anahtarlı bir hesap ise `no limits (API)` gösterir.

`ccam statusline install` sonrasında Claude Code, `[work] · Opus · 5h 24% · 7d 41%` gibi bir satır gösterir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--json` | Raporları JSON olarak yazdırır: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | kapalı |
| `--watch[=N]` | Ctrl+C'ye basana kadar her `N` saniyede bir yeniden çizer. Tek başına `--watch` veya 0 ya da daha küçük bir değer 60 anlamına gelir | kapalı |
| `-v`, `--verbose` | Her hesabın altına bir `source:` satırı ekler: sayılarının nereden geldiği (kullanım uç noktası, Claude Code, durum satırı veya kayıtlı bir dosya) | kapalı |
| `--no-cache` | 60 saniyelik önbelleği yok sayar ve uç noktaya yeniden sorar | kapalı |
| `--source live\|claude\|cache` | `live` uç noktaya hesabın token'ıyla sorar; `claude` kullanımı Claude Code'a sorgulatır (aşağıda); `cache` yalnızca kayıtlı sayıları gösterir | `usage_source`, yoksa `live` |
| `--refresh` | Süresi dolmuş erişim tokenlarını yalnızca çalışmayan hesaplar için yeniler | kapalı |
| `--dry-run` | Kabul edilir: `ccam usage` kendi önbelleği dışında hiçbir şeyi değiştirmez. `--refresh` ile birlikte reddedilir: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | kapalı |

Başka herhangi bir `--source` değeri `--source must be live, claude or cache` ile başarısız olur.

### Claude Code'un kendisinden kullanım

`config.toml` içinde `usage_source = "claude"` ile (ya da `--source claude` ile) ccam giriş token'ını hiç okumaz. Her hesap için hesabın klasöründe `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` başlatır ve Claude Code'un `get_usage` isteğini standart girdiden gönderir (Claude Code 2.1.270 veya sonrası); Claude Code kullanımı kendi girişiyle sorgular ve 5 saatlik, haftalık ve model pencereleriyle yanıt verir. `--safe-mode`, hesabın hook'larının, MCP sunucularının ve eklentilerinin başlamasını engeller ve hiçbir oturum kaydedilmez. Her hesap bir Claude Code başlatması, yaklaşık bir saniye alır; bu yüzden bu kaynak `live` kaynağından yavaştır. Hesaplara aynı anda, her birine 20 saniye içinde sorulur ve 60 saniyelik önbellek her zamanki gibi geçerlidir. Claude Code bu isteği deneysel olarak işaretler; bu yüzden yanıt vermeyi bırakırsa ccam kayıtlı sayıları gösterir ve `ccam doctor` bunu bildirir. Bir `--dry-run` Claude Code'u başlatmaz (Claude Code hesabın klasörüne yazar) ve kayıtlı sayıları gösterir. `usage_source`, `live` veya `claude` alır; `cache` yalnızca bir `--source` seçeneğidir. `ccam usage` argüman almaz.

`ccam statusline`, Claude Code'un durum satırı JSON'unu standart girdiden okur. Tek seçeneği `--wrap <command>` seçeneğidir: ilk çıktı satırının başına hesap adı eklenen mevcut bir durum satırı komutu. `ccam statusline install` seçenek almaz. İkisi de `--dry-run` seçeneğini `` --dry-run is not supported by `ccam statusline`; nothing was changed `` ile reddeder.

## Nasıl çalışır

Sayılar iki yerden gelir: her hesabın kendi girişiyle çağrılan, Anthropic'in belgelenmemiş kullanım uç noktası `GET https://api.anthropic.com/api/oauth/usage` ve Claude Code'un durum satırı komutuna verdiği değerler. Bunların hepsi `internal/usage` içindedir.

1. `cli.usageCmd`, `config.toml` dosyasını yükler ve her hesap için bir `usage.Account` oluşturur: klasörü `~/.local/share/ccam/profiles/<name>/`, oturum klasörü (çalışıp çalışmadığını görmek için) ve yerinde içe aktarılıp aktarılmadığı.
2. `Fetcher.FetchAll` her hesap için bir goroutine başlatır. Hepsi tek bir `http.Client` paylaşır; böylece uç noktaya yapılan bağlantılar yeniden kullanılır. Bağlanma 2 sn sürebilir ve her hesaba toplam 3 sn verilir. İstemci yönlendirmeleri izlemez, çünkü bir yönlendirme tokenı başka bir yere yeniden gönderir; yanıtın da en fazla 1 MB'ını okur.
3. Her hesap bittikçe satırı yeniden çizilir. Başarısız olan veya zamanında yanıt vermeyen bir hesap kayıtlı sayılara geri döner; böylece yavaş bir hesap diğerlerini hiçbir zaman bekletmez.

Her hesap uç noktaya karşı kendi goroutine'inde çalışır; zaman aşımına uğrayan hesap bunun yerine kayıtlı sayılarını gösterir.

<!-- motion:usage-fetch -->

`personal` ve `work` zamanında yanıt verir ve canlı sayılar gösterir. `side` 3 sn'den uzun sürer; bu yüzden `fallback()` iki önbellek dosyasını okur ve en son yazılanı kullanır. `usage.Table()` her hesap için bir blok çizer ve `usage.Best()` son satırda adı geçen hesabı seçer.

Tek bir hesap için `Fetcher.Fetch` şu kontrolleri sırayla yapar:

<!-- motion:usage-cache -->

1. `--source cache` ile doğrudan `fallback()` fonksiyonuna gider.
2. `--no-cache` ayarlı değilse, `~/.local/state/ccam/usage/<name>.json` içinde 60 saniyeden (`CacheTTL`) yeni olan canlı bir sonuç olduğu gibi döndürülür. `ccam usage` komutunu art arda iki kez çalıştırmak hesap başına bir isteğe mal olur.
3. `readCredentials()`, `.credentials.json` dosyasını (veya macOS Anahtar Zinciri'ni) değiştirmeden okur ve yalnızca `claudeAiOauth` anahtarının kodu çözülür. Abonelik girişi yoksa `not logged in`, `.claude.json` bir API anahtarı tutuyorsa `no limits (API)` gösterilir.
4. Süresi dolmuş bir erişim tokenı kullanılmaz. `--refresh` olmadan hesap, süresi dolmuş olarak işaretlenmiş kayıtlı sayılarını gösterir.
5. Aksi halde ccam uç noktayı çağırır. 401 veya 403 süresi dolmuş bir giriş sayılır, diğer her durum kodu veya okunamayan bir gövde hata sayılır ve her ikisi de kayıtlı sayılara geri döner. 429, uç noktaya çok sık sorulduğu anlamına gelir: ccam `<name>.pause` dosyasına yeniden sorma zamanını (`Retry-After` ya da 5 dakika, en fazla bir saat) yazar, kayıtlı sayıları `rate limited by the usage endpoint; asking again after 15:10` ile gösterir ve o zamana kadar, `--no-cache` ile bile, o hesap için uç noktayı çağırmaz. Hemen yeniden denemek yalnızca limitin sürmesine yol açar.
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
4. ccam `[name] · model · 5h N% · 7d N%` yazdırır. `--wrap` ile komutunuzu aynı girdiyle `sh -c` üzerinden çalıştırır ve bunun yerine `[name]` ve ardından komutun ilk satırını yazdırır. Windows'ta komutu Claude Code'un yaptığı gibi çalıştırır: kuruluysa Git Bash ile, değilse PowerShell ile.
5. 5 saatlik veya haftalık kullanım %90 ya da üzerindeyse ve başka bir hesap varsa satır `→ ccam next --auto --continue` ile biter, örneğin `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Sarmalanan komuta 2 saniye verilir. Sonrasında sonlandırılır ve satır yalnızca hesap adını gösterir; böylece yavaş bir komut Claude Code'un durum satırını engelleyemez.

`ccam statusline install`, her hesabın `settings.json` dosyasına `statusLine` yazar. Bir hesabın zaten bir durum satırı komutu varsa ccam onu `--wrap '<old command>'` olarak korur. Windows'ta ccam'in yolu eğik çizgilerle ve ev klasörü için `~` ile yazılır (`~/.local/bin/ccam.exe statusline`); bunu hem Git Bash hem PowerShell anlar. Bir grup üzerinden paylaşılan `settings.json` bağlantısı aracılığıyla bir kez yazılır; böylece gruptaki her hesap onu alır. Yeniden çalıştırmak, zaten ccam kullanan dosyalar için `(already set)` yazdırır.

## Güvenlik ve uç durumlar

- `ccam usage`, siz `--refresh` vermedikçe hiçbir zaman bir token dosyasını değiştirmez. Bir girişin süresi dolduğunda hesabı bir kez `ccam <name>` ile çalıştırın, Claude Code onu yeniler.
- `--refresh` çalışan hesaplara, içe aktarılmış hesaplara veya Anahtar Zinciri girişlerine hiçbir zaman dokunmaz. Reddedilen bir yenileme, kayıtlı sayıları ve nedeni gösterir.
- Kullanım uç noktası herkese açık olarak belgelenmemiştir ve değişebilir. Değiştiğinde kayıtlı sayıları yine görürsünüz ve `ccam doctor` biçim değişikliğini bildirir.
- Kayıtlı sayı ve yanıt olmadığında hesap satırı `lookup failed:` ve nedeni gösterir, örneğin `--source cache` ile `no saved usage`.
- Önbellek dosyaları yalnızca yüzdeleri ve sıfırlanma zamanlarını tutar, hiçbir zaman token tutmaz.

## İlgili

- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [Sonraki hesaba geçme](/features/next)
- [Kullanım limiti iş akışı](/guides/usage-limit-workflow)
- [Komutlar](/docs/commands)
