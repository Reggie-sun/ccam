---
title: Kısayol komutları
description: ccam shortcut ile bir hesaba ccwork gibi kendi komutunu verin ve ~/.local/bin içindeki bağlantıları ccam shims install ile yeniden oluşturun.
updated: 2026-10-23
order: 3
translation: machine
sourceHash: "38e2d1aba397"
---

Kısayol, adını kendiniz verdiğiniz ve Claude Code'u tek bir hesapla başlatan bir komuttur: `ccwork --resume`, `ccam work --resume` ile aynı işi yapar. `~/.local/bin` içinde ccam'e giden bir bağlantıdır; bu yüzden bash, zsh, betikler ve programları `PATH` üzerinden çalıştıran her şeyde çalışır. Kısayollar isteğe bağlıdır; sık başlattığınız hesaplar için bir tane oluşturun.

## Kullanım

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

ccam her değişikliği onaylar:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Bir hesabın en fazla bir kısayolu olur ve bu kısayol `~/.config/ccam/config.toml` içinde `shortcut` olarak saklanır. `ccam add` sihirbazı `Create a shortcut command for <name>? [y/N]` ve ardından `Shortcut name [cc<name>]:` diye sorar. Sihirbaz dışında `ccam add` yalnızca `--shortcut` ile kısayol oluşturur. `ccam list` her hesabın kısayolunu, yoksa `-` gösterir.

Kısayol, `ccam <name>` komutunun addan sonra yaptığı gibi her argümanı Claude Code'a iletir. Buna `--dry-run` ve `--help` de dahildir: bir başlatmayı önizlemek için `ccam --dry-run work` kullanın.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--remove` | Hesabın kısayolunu ve bağlantısını siler | kapalı |
| `-h`, `--help` | Komutun yardımı | |
| `--dry-run` | `ccam shortcut` bunu dikkate alır: adı kontrol eder ve kaydetmeden veya bağlamadan `Would set the shortcut of work to ccwork.` ya da `Would remove ccwork, the shortcut of work.` yazdırır | kapalı |

`ccam shims install` argüman almaz ve `--dry-run` seçeneğini `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` ile reddeder.

Ad kuralları, bu sırayla kontrol edilir:

- `a-z`, `0-9`, `.`, `-` ve `_` karakterlerinden oluşan, harf veya rakamla başlayan 1 ile 64 karakter;
- `ccam`, `claude` veya `run` ya da `usage` gibi bir ccam komutu olamaz;
- başka bir hesabın kısayolu olamaz;
- `~/.local/bin` içinde zaten ccam'in bağlantısı (ya da ccam'in oluşturduğu, bu hesabın kendi kısayolu) olmayan bir dosya olamaz;
- `PATH` üzerindeki başka bir yürütülebilir dosyanın adı olamaz.

`config.toml` yüklenirken aynı kurallara göre denetlenir: kalıba uymayan, `claude` ya da `ccam` olan veya iki hesaba verilmiş, elle düzenlenmiş bir kısayol, düzeltilene kadar her komutun hesabın adını veren bir mesajla başarısız olmasına yol açar.

Bir komutu ve `--remove` seçeneğini birlikte vermek `give a command or --remove, not both` ile durur. Bir hesaba zaten sahip olduğu adı vermek `work: ccwork (unchanged)` yazdırır ve eksikse bağlantısını yeniden oluşturur.

## Nasıl çalışır

Bir ad ayarlamak, bin klasörü ve ccam ikili dosyasının gerçek yoluyla `shim.Check` fonksiyonunu çağıran `checkShortcut` fonksiyonunu çalıştırır:

<!-- motion:shortcuts-check -->

Ad, kontrollerden sırayla geçer. Her kontrolün yanında reddettiği bir ad bulunur: `Work` kalıba uymaz, `run` bir ccam komutudur, `ccside` başka bir hesaba aittir, `ccfile` `~/.local/bin` içinde mevcut bir dosyadır ve `git`, `/usr/bin/git` dosyasını gizlerdi. Yalnızca her kontrolden geçen bir ad kaydedilir ve bağlanır.

1. `shim.Check` adı önce `^[a-z0-9][a-z0-9._-]{0,63}$`, ardından `config.Reserved`, sonra `cfg.FindShortcut` ile karşılaştırır. `~/.local/bin/<name>` üzerinde `os.Lstat` çağırır; `PATH` üzerindeki diğer her klasör için `os.Stat` çağırır ve ccam'e giden bir bağlantı olmayan yürütülebilir bir dosyayı reddeder.
2. ccam, `config.toml` dosyasını `fsutil.AtomicWrite` ile kaydeder (geçici dosya, `fsync`, yeniden adlandırma).
3. `shim.Remove` eski kısayolun bağlantısını siler, ancak yalnızca ccam'e giden bir bağlantıysa.
4. `shim.Install`, yapılandırılmış her kısayolu ad sırasıyla dolaşır ve eksik olan her birinin bağlantısını `os.Executable` ve `filepath.EvalSymlinks` ile çözümlenen ccam ikili dosyasına `os.Symlink` ile oluşturur. Bir kısayolun yolundaki bir dosya diğerlerini durdurmaz; hata her birini adlandırır.

ccam oluşturduğu her bağlantıyı ve onu çalıştırdığı ikili dosyayı `~/.local/state/ccam/shortcuts.json` dosyasına kaydeder. `os.Readlink` ccam yolunu döndürdüğünde veya bağlantı ile ccam yolu aynı dosyaya çözümlendiğinde bir bağlantı ccam'in kendisine ait sayılır. İkili dosyası o zamandan beri taşınmış veya değiştirilmiş (başka bir klasöre taşınmış bir ccam, yeniden kurulum) kayıtlı bir bağlantı eskimiştir ve `shim.Install` onu bu ccam'e yönlendirir. Hâlâ kurulu olan başka bir ccam'e giden kayıtlı bir bağlantı olduğu gibi bırakılır. Her kaldırma ve her üzerine yazma kontrolü bu testlerden geçer; böylece ccam, ccam'in kendi kopyanız veya kendiniz oluşturduğunuz bir bağlantı gibi, kendi oluşturmadığı dosyalara dokunmaz.

Bir kısayolu çalıştırdığınızda ccam'in ihtiyaç duyduğu tek şey yazdığınız addır:

<!-- motion:shortcuts-link -->

Kabuk `ccwork` komutunu `PATH` üzerinde bulur, çekirdek bağlantıyı izleyip ccam ikili dosyasını başlatır ve `argv[0]` hâlâ `ccwork` olur. `app.Main` içinde `ccam` dışındaki her ad `config.toml` dosyasını bir kez yükler ve `shortcut` değeri `filepath.Base(argv[0])` olan hesabı bulan `shim.Lookup` fonksiyonunu çağırır. Ardından `launcher.LaunchWith`, yapılandırma zaten yüklenmiş olarak [`ccam <name>`](/features/launch) ile aynı adımları izler ve `argv[1:]` değerini `claude` komutuna dokunmadan iletir.

Neden kabuk takma adı veya sarmalayıcı betik değil de sembolik bağlantı: sembolik bağlantı her kabukta ve betiklerde çalışır, fazladan bir kabuk başlatması eklemez (yaklaşık 1 ile 2 ms) ve başlatma yine tek bir dosya okur. `scripts/perf-check.sh`, bir kısayolu `hyperfine` ile `ccam <name>` ile aynı 10 ms (p95) sınırına göre ölçer.

### Windows'ta

Kısayol, `%USERPROFILE%\.local\bin` içindeki `ccwork.exe` dosyasıdır ve `shim.Name` aramadan önce `.exe` uzantısını (büyük veya küçük harf fark etmeksizin) atar; böylece `ccwork`, `ccwork.exe` ve `CCWORK.EXE` hepsi `work` hesabını başlatır. ccam önce bir sembolik bağlantı oluşturur; Windows buna Geliştirici Modu'nda veya bir yöneticiye izin verir. Bu hak olmadan bir sabit bağlantı oluşturur: aynı `ccam.exe` dosyası için ikinci bir ad; bunun için ayrıcalık gerekmez. Bir sabit bağlantı, `ccam.exe` ile aynı dosyaysa ccam'in kendi bağlantısı sayılır. Bir sabit bağlantı oluşturulduğu dosyayı çalıştırmaya devam eder; bu yüzden `ccam update`, `ccam.exe` yerine yeni bir dosya koyduğunda sabit bağlantılı her kısayolu `shim.Relink` ile yeniler ve `✓ shortcut ccwork renewed` yazdırır. `install.ps1` betiğini yeniden çalıştırmak da `ccam shims install` üzerinden aynısını yapar (`OK 2 shortcuts now run the new ccam`); daha eski bir ccam'in kopyası olan her kısayol için `ccam shims install` ve `ccam doctor --fix` de öyle. Ad kontrolü, diğer `PATH` klasörlerinde `PATHEXT` içindeki her uzantıyla (`.exe`, `.cmd` vb.) `<name>` arar.

Bir kısayolu bir kabukta (PowerShell, cmd.exe, Git Bash) yazın. Dosya Gezgini'nden, Çalıştır iletişim kutusundan veya `Start-Process` ile açılan bir sembolik bağlantı düz `ccam` olarak başlar ve hesap yerine kontrol panelini açar.

## Güvenlik ve uç durumlar

- **Hiçbir hesabın sahiplenmediği bir bağlantı.** `~/.local/bin` içindeki ccam'e giden bir bağlantının adı hiçbir hesap tarafından kullanılmıyorsa, örneğin `config.toml` içindeki `shortcut` değerini elle düzenledikten sonra, ccam varsayılan hesaba geri dönmez. `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` ile durur. Adı `ccam shortcut` ile yeniden ayarlayın veya bağlantıyı silin.
- **Mevcut dosyalar.** ccam kendi oluşturmadığı bir dosyayı hiçbir zaman değiştirmez. `ccam shims install`, `<path> already exists and is not a ccam link, so it was left alone` bildirir ve diğer bağlantıları yine de oluşturur. `ccam add` ve `ccam import` bu durumda `shortcut: ccwork (not created; see above)` yazdırır.
- **Taşınmış bir ccam ikili dosyası.** ccam'in oluşturduğu bağlantılar bu durumda artık var olmayan bir yolu gösterir. `ccam shims install` (veya `ccam doctor --fix`) onları, komutu çalıştırdığınız ccam'e yönlendirir.
- **Yeniden adlandırma ve kaldırma.** `ccam rename` kısayol adını korur ve bağlantı çalışmaya devam eder, çünkü ccam'i gösterir ve ccam adı `config.toml` içinde arar. `ccam remove` hesabın bağlantısını siler. `ccam shortcut <account> <new>` yenisini oluşturmadan önce eski bağlantıyı siler.
- **PATH.** ccam bir bağlantı oluşturduğunda `~/.local/bin` `PATH` üzerinde değilse, `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` ve ardından eklenecek `export PATH=...` satırını yazdırır. Windows'ta bunun yerine klasörü kullanıcı `PATH` değişkeninize ekleyen PowerShell satırını yazdırır.
- **Sonradan kontrol.** `ccam doctor`, `missing shortcuts in /home/me/.local/bin: ccwork` ve daha eski bir ccam'i çalıştıran kısayolları bildirir, klasör `PATH` üzerinde değilse uyarır. `ccam doctor --fix` eksik bağlantıları oluşturur, eskimiş olanları yeniler ve yalnızca oluşturduğu bağlantılar için `→ fixed` der. Hâlâ kurulu olan başka bir ccam'i çalıştıran bir kısayol çalışır ve yalnızca bilgidir; `ccam doctor --fix` onu komutu çalıştırdığınız ccam'e yönlendirir, `ccam add` ve `ccam import` ise onu ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)`` olarak gösterir.

## İlgili

- [Claude Code'u başlatma](/features/launch)
- [Hesapları yönetme](/features/accounts)
- [Komut başvurusu](/docs/commands#shortcuts)
- [ccam'i kurma](/docs/install)
