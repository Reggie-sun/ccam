---
title: ccam'i güncelleme
description: ccam update, bir SHA-256 kontrolünden sonra ccam'i en son GitHub sürümüyle değiştirir; --check yalnızca bildirir, --version belirli bir sürüm etiketini kurar.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: 598931e4801f
---

`ccam update`, GitHub'dan bir ccam sürümü indirir, onu sürümün `checksums.txt` dosyasıyla doğrular ve çalışan ikili dosyayı tek bir yeniden adlandırmayla değiştirir. ccam'i kurulum betiğiyle veya bir arşivden kurduysanız kullanın. [Kontrol paneli](/features/dashboard) de her açıldığında GitHub'a en son sürümü sorar ve onu `U` ile önerir; bu tuş sorduktan sonra aynı adımları çalıştırır. Diğer komutlar hiçbir zaman GitHub ile bağlantı kurmaz; bu yüzden çevrimdışı da aynı şekilde çalışırlar.

## Kullanım

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Daha yeni bir sürüm yoksa ccam `ccam 0.4.0 is up to date.` yazdırır ve hiçbir şeyi değiştirmez. `--version` zaten sahip olduğunuz sürümü belirtirse `ccam 0.4.0 is already installed.` yazdırır.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--check` | Yalnızca daha yeni bir sürüm olup olmadığını bildirir | kapalı |
| `--version <tag>` | `v0.3.0` gibi, mevcut sürümden yeni veya eski bu sürüm etiketini kurar | en son |
| `--dry-run` | Desteklenir: en son etiketi arar (`--version` bir etiket vermedikçe), ikili dosya üzerinde aynı kontrolleri çalıştırır ve indirmeden indirme, sağlama toplamı ve değiştirme yollarını yazdırır | kapalı |

`--check` ve `--version` birlikte kullanılamaz. Etiket, isteğe bağlı bir ön sürüm veya derleme sonekiyle `vMAJOR.MINOR.PATCH` biçiminde olmalıdır; baştaki `v` yazılmayabilir. Başka her şey başarısız olur, örneğin `--version takes a release tag such as v0.3.0, not "1.2"`.

`CCAM_RELEASES_URL` ortam değişkeni, bir yansı veya testler için `https://github.com/jhyunwoo/ccam/releases` adresinin yerini alır. Aynı şekilde düzenlenmiş olmalıdır: `<url>/latest`, `<url>/tag/<tag>` adresine yönlendirir ve dosyalar `<url>/download/<tag>/<asset>` konumundadır.

## Nasıl çalışır

ccam, kurulu ikili dosyaya dokunmadan önce indirmeyi denetler ve ikili dosyayı tek bir yeniden adlandırmayla değiştirir. Kod, yalnızca `ccam update` komutunun içe aktardığı `internal/update` içindedir. Başlatma yolu onu hiçbir zaman çağırmaz.

1. ccam çalışan ikili dosyayı `os.Executable` ile bulur ve sembolik bağlantıları `filepath.EvalSymlinks` ile çözümler (`cli.updateCmd`). Değiştirilen, bu gerçek yoldur.
2. `--version` olmadan ccam yönlendirmeyi izlemeden `<releases>/latest` adresini ister (`Updater.Latest`) ve etiketi `Location` başlığından okur (`…/tag/v0.5.0`). Bu, GitHub API'si yerine sürümler sayfasını kullanır; böylece API'nin saatte 60 istek sınırı uygulanmaz. Arama 10 saniye sonra zaman aşımına uğrar.
3. ccam etiketi katı bir kalıba göre kontrol eder (`parseVersion`). Etiket, indirme URL'lerinde ve yazdırılan bir kabuk komutunda yer alır; bu yüzden `/`, `%`, boşluk veya kabuk karakterleri içeren hiçbir şey geçemez. ccam ardından onu çalışan sürümle karşılaştırır (`compareVersions`); bir ön sürüm, kendi sürümünden önce sıralanır.
4. ccam ikili dosyayı değiştirip değiştiremeyeceğine karar verir (`managed()`, `writable()`; aşağıya bakın).
5. ccam `ccam_<os>_<arch>.tar.gz` (Windows'ta `ccam_windows_<arch>.zip`) ve `checksums.txt` dosyalarını aynı anda, ikisi için toplam 2 dakikayla indirir (`Updater.Fetch`). Her indirme 256 MB ile sınırlıdır.
6. Başka herhangi bir şey olmadan önce arşivin SHA-256 değeri `checksums.txt` içindeki satırıyla eşleşmelidir. Eksik bir satır veya farklı bir sağlama toplamı güncellemeyi `nothing was changed` ile durdurur.
7. ccam tar dosyasını okur ve `ccam` adlı ilk normal dosyayı alır (`extract()`; `extractZip` zip içinden `ccam.exe` dosyasını alır). Klasörler, bağlantılar ve diğer girdiler yok sayılır; ikili dosya da 256 MB ile sınırlıdır.
8. ccam yeni ikili dosyayı `fsutil.AtomicWrite` ile yazar (`Updater.Install`): ikili dosyanın yanında, eski dosyanın izin bitlerine sahip geçici bir `.ccam.tmp*` dosyası diske yazılır ve eskisinin üzerine yeniden adlandırılır.

Etiket aramasından yeniden adlandırmaya kadar tüm akış:

<!-- motion:update-pipeline -->

Arşiv ve `checksums.txt` paralel olarak indirilir ve sağlama toplamı eşleşene kadar hiçbir şey yazılmaz. Yeniden adlandırma atomiktir; bu yüzden bir çökme veya dolu bir disk ya eski ikili dosyayı ya da yenisini bırakır, hiçbir zaman yarım bir dosya bırakmaz. `~/.local/bin/ccwork` gibi kısayol komutları ccam ikili dosyasına giden sembolik bağlantılardır; bu yüzden hiçbir değişiklik olmadan yeni sürümü çalıştırırlar.

Windows'ta çalışan bir programın dosyası değiştirilemez veya silinemez, yalnızca yeniden adlandırılabilir. ccam yeni ikili dosyayı `ccam.exe.new` olarak yazar, çalışan `ccam.exe` dosyasını `ccam.exe.old` olarak, yeni dosyayı da `ccam.exe` olarak yeniden adlandırır (`Updater.Install`). Son adım başarısız olursa eski ikili dosya geri adlandırılır. Bir sonraki güncelleme silebildiği `ccam.exe.old` dosyalarını siler; daha önceki bir `ccam.exe.old` dosyasından hâlâ bir oturum çalışıyorsa çalışan ikili dosya bunun yerine `ccam.exe.old-<time>` olur. Sembolik bağlantı yerine sabit bağlantı olan kısayollar ([Kısayol komutları](/features/shortcuts)) hemen ardından yenilenir, böylece onlar da yeni sürümü çalıştırır; bu, yalnızca az önce değiştirilen ikili dosyanın değil, ccam'in oluşturduğu, ccam'in daha eski herhangi bir sürümünün sabit bağlantılarını da kapsar.

### Paket yöneticisiyle yapılan kurulumlar

`.deb` paketinden kurulan bir ccam kopyası apt'ye aittir ve onu apt'nin haberi olmadan değiştirmek paket veritabanını dosyayla uyumsuz bırakır. İndirmeden önce ccam ikili dosyanın kime ait olduğunu kontrol eder:

<!-- motion:update-managed -->

1. `/usr/bin`, `/bin`, `/usr/sbin` veya `/sbin` içindeki bir ikili dosya paket yöneticisine ait sayılır.
2. Başka bir yerdeyse, `dpkg` kuruluysa ve `dpkg -S <path>` başarılı olursa yine paket yöneticisine ait sayılır.
3. Aksi halde ccam, oraya yazabildiğinden emin olmak için ikili dosyanın yanında geçici bir `.ccam-update-*` dosyası oluşturur ve kaldırır.

Paket yöneticisine ait bir ikili dosyada ccam durur ve apt adımlarını yazdırır:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

Linux dışındaki sistemlerde yalnızca kurulum betiği satırı yazdırılır: Windows'ta `irm https://ccam.dev/install.ps1 | iex`. Klasör yazılabilir değilse ccam, ccam'i `~/.local/bin` içine kuran kurulum betiğini yeniden çalıştırmayı önerir.

## Güvenlik ve uç durumlar

- Bir geliştirme derlemesi sürümünü `dev` olarak bildirir ve karşılaştırılacak bir sürümü yoktur; bu yüzden `ccam update` ve `ccam update --check` onu reddeder. `ccam update --version vX.Y.Z` onun üzerine bir sürüm kurar. `go install` ile yapılmış derlemeler için bkz. [ccam kurulumu](/docs/install#other-ways-to-install).
- Başarısız bir indirme, platformunuz için eksik bir dosya veya sağlama toplamı uyuşmazlığı kurulu ikili dosyayı olduğu gibi bırakır.
- Var olmayan bir sürüm `there is no release <tag>` ile başarısız olur. GitHub HTTP 403 veya 429 ile reddettiğinde ccam bunun muhtemelen bir istek sınırı olduğunu söyler ve daha sonra yeniden denemenizi ister.
- `--check` yalnızca okur. İkili dosyanın paket yöneticisine ait veya yazılabilir olup olmadığını kontrol etmez.
- Daha eski bir etiketle `--version`, geri dönmenin yoludur.

## İlgili

- [ccam kurulumu](/docs/install)
- [Komut başvurusu](/docs/commands)
- [Sık sorulan sorular](/docs/faq)
