---
title: Hesapları yönetme
description: Claude Code hesaplarını ccam add, login, list, rename, remove ve import ile ekleyin, giriş yapın, listeleyin, yeniden adlandırın, kaldırın ve içe aktarın.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

Her ccam hesabı, kendi girişine sahip ayrı bir Claude Code yapılandırma klasörüdür. Bu komutlar hesap oluşturur, hesaplarda oturum açar, hesapları listeler, yeniden adlandırır ve kaldırır; ayrıca zaten sahip olduğunuz klasörleri kaydeder.

## Kullanım

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Bir komut `<name|number>` aldığı her yerde hesap adını veya yuva numarasını, yani `ccam list` çıktısındaki `#` sütununu verebilirsiniz.

Terminalde seçeneksiz çalıştırılan `ccam add` bir sihirbaz açar. Sihirbaz önce adı sorar, ardından paylaşılan öğelerin bir kontrol listesini gösterir (konuşma geçmişi dışındaki her şey işaretlidir), sonra `Create a shortcut command for work? [y/N]` sorusunu, `y` yanıtında da `Shortcut name [ccwork]:` sorusunu sorar. `--share`, `--shortcut` veya `--no-login` seçeneklerinden herhangi biri ya da terminal yerine bir boru (pipe) sihirbazı atlar; bu durumda bir ad zorunludur ve `--shortcut` bir ad vermedikçe kısayol oluşturulmaz.

`ccam list` yuvayı, adı, kısayolu, e-postayı, paylaşılan öğeleri ve son kullanımı gösterir; varsayılan hesabın önünde `*` bulunur:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json` ile her hesap `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (hesap hiç başlatılmadıysa yer almaz) ve `share` alanlarına sahip bir nesnedir. `share` her öğeyi (`mcp`, `settings`, `memory`, `extensions`, `history`) grup adıyla veya `"local"` ile eşler.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `ccam add --no-login` | Hesabı giriş yapmadan oluşturur; daha sonra `ccam login <name>` ile giriş yapın | kapalı |
| `ccam add --share <item=group>` | Bir öğeyi bir grupla paylaşır ya da `item=local` ile özel tutar. Tekrarlanabilir | `default` içinde `mcp`, `settings`, `memory`, `extensions`; `history` yerel |
| `ccam add --shortcut <command>` | Bu kısayol komutunu `~/.local/bin` içinde de oluşturur | yok |
| `ccam list --json` | Hesapları JSON olarak yazdırır | kapalı |
| `ccam remove -y`, `--yes` | Sormadan kaldırır | kapalı |
| `ccam import --name <name>` | Verdiğiniz klasör için hesap adı | klasör adından türetilir |
| `ccam import -y`, `--yes` | Bulunan her klasörü sormadan içe aktarır | kapalı |
| `--dry-run` | Hiçbir şeyi değiştirmeyen `ccam list` ve `ccam completion` bunu dikkate alır. `add`, `login`, `logout`, `use`, `rename`, `remove` ve `import` ise `` --dry-run is not supported by `ccam add`; nothing was changed `` (komutun adıyla) mesajıyla reddeder | kapalı |

`ccam login`, `ccam logout`, `ccam use` ve `ccam rename` komutlarının kendi seçenekleri yoktur. `ccam completion` `bash` veya `zsh` alır; başka bir değer `only bash and zsh are supported` hatasıyla başarısız olur.

Hesap adları `a-z`, `0-9`, `-` ve `_` karakterlerinden oluşan, harf veya rakamla başlayan 1 ile 32 karakter arası adlardır. Bir ad yalnızca rakamlardan oluşamaz, çünkü `ccam 2` 2 numaralı yuva anlamına gelir. Bir ad bir ccam komutu da olamaz (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`); böylece `ccam <name>` hiçbir zaman belirsiz olmaz. Liste, `internal/config` içindeki `config.Reserved` değişkenidir.

## Nasıl çalışır

`ccam add work --shortcut ccwork`, `addProfile` (`internal/cli/profiles.go`) içinde şu adımları çalıştırır:

<!-- motion:accounts-add -->

1. Herhangi bir şey oluşturulmadan önce `config.ValidateName` adı, `shim.Check` ise kısayolu kontrol eder.
2. `profile.Add`, `~/.local/share/ccam/profiles/work/` klasörünü `0700` moduyla oluşturur ve hesaba en küçük boş yuva numarasını (`Config.NextSlot`) verir; böylece kaldırılmış bir hesabın numarası yeniden kullanılır. İlk hesap varsayılan olur.
3. `applyShare`, yeni ve boş klasörü `share.Build` ve ardından `Apply` ile paylaşım gruplarına bağlar; örneğin `settings.json` dosyasını `groups/default/settings.json` dosyasına. Bu başarısız olursa yeni klasör kaldırılır, çünkü geride kalan bir klasör her yeniden denemeyi `folder already exists` ile engeller.
4. Hesap, `fsutil.AtomicWrite` aracılığıyla `~/.config/ccam/config.toml` dosyasına kaydedilir: önce geçici bir dosya, `fsync`, ardından `rename`; böylece bir çökme hiçbir zaman yarım bir yapılandırma bırakmaz.
5. `runLogin`, `CLAUDE_CONFIG_DIR` hesap klasörüne ayarlanmış olarak `claude auth login` komutunu başlatır. Giriş başarısız olursa hesap yerinde kalır; ccam hatayı yazdırır ve daha sonra `ccam login work` çalıştırmanızı söyler.
6. `profile.ReadEmail`, `oauthAccount.emailAddress` değerini `profiles/work/.claude.json` dosyasından alır. Bu dosya birkaç megabayt olabilir; bu yüzden `TopLevelValue` tüm dosyayı ayrıştırmak yerine yalnızca o anahtarı tarar.
7. `shim.Install`, `~/.local/bin/ccwork` bağlantısını ccam ikili dosyasına yönlendirir; ccam başladığında `argv[0]` değerini `config.toml` içinde arar.

Bir giriş dakikalar sürebilir ve bu arada başka bir ccam yapılandırmayı kaydedebilir. Bu yüzden bekleyen komutlar, `config.toml` dosyasını yeniden yükleyen, tek değişikliği (burada e-postayı) uygulayan ve kaydeden `env.update` üzerinden kaydeder. `ccam remove` da aynı nedenle sorusundan sonra yeniden yükler.

`ccam login`, mevcut bir hesap için 5. ve 6. adımları çalıştırır. `ccam logout`, hesapta `claude auth logout` komutunu çalıştırır ve kayıtlı e-postayı temizler. `ccam use` yalnızca `config.toml` içindeki `default` değerini yeniden yazar.

### Uzak bir makinede oturum açma

SSH üzerinden `claude auth login` komutunun tarayıcı tarafı bozulabilir. Claude Code, oturum açmayı sunucudaki `http://localhost:<port>/callback` adresinde bir geri çağırma sunucusunda bekler; ancak oturum açmayı tamamlayan tarayıcı sizin bilgisayarınızda çalışır. Orada bu adres hiçbir yere gitmez ve tarayıcı `ERR_CONNECTION_REFUSED` gibi bir hata gösterir. Adres yine de tek kullanımlık kodu taşıdığı için ccam onu iletebilir:

<!-- motion:accounts-relay -->

`runLogin` terminalinizi Claude Code'a vermez. Claude Code'un standart girdisini bir boruya bağlar ve `relayLogin` yazdığınız her satırı bu boruya kopyalar; böylece Claude Code'un `Paste code here` istemine yapıştırılan bir kod yine ulaşır. Bir satır geri çağırma adresi olduğunda `callbackURL` onu yalnızca ana makine `localhost`, `127.0.0.1` veya `[::1]` ise, bir port verilmişse, yol `/callback` ise ve hem `code` hem `state` mevcutsa kabul eder; `http://` yazılmayabilir. Ardından `deliverCallback`, tarayıcının yapacağı gibi, bu adresi sunucuda 10 saniyelik bir sınırla ve yönlendirmeleri izlemeden ister. İstek, bekleyen `claude auth login` komutunun geri çağırma sunucusuna ulaşır (diyagramdaki kesikli çizgi). Claude Code kodu alır ve girişi kaydeder; Linux'ta hesap klasöründeki `.credentials.json` dosyasına. Yapıştırılan bir satır ccam'in başka bir ana makineyle bağlantı kurmasına asla yol açamaz.

ccam, giriş başlamadan önce şunu yazdırır:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Başarılı olduğunda `→ Passed the sign-in to Claude Code on this machine.` yazdırır.

Oturum açma bir terminalde girdi beklerken ccam terminalin yankısını ve satır modunu kapatır ve yazdığınız her karakter için bir `*` gösterir; böylece yapıştırılan bir kod veya geri çağırma adresi ekranda görünmez. `readLine` satır düzenlemeyi kendisi yapar: `Enter` satırı gönderir, `Backspace` bir karakter siler ve boş bir satırda `Ctrl+D` girdiyi bitirir. Terminal, oturum açma bittiğinde ve ayrıca ccam'i yine sonlandıran `Ctrl+C`, `SIGTERM` veya `SIGHUP` durumunda geri yüklenir. Standart girdi ve standart hata bir terminale bağlı değilse satırlar geldikleri gibi kopyalanır.

### Yeniden adlandırma, kaldırma ve içe aktarma

`ccam rename` ve `ccam remove`, hesabın klasörü ve ccam'in hesabın adı altında başka yerlerde tuttukları üzerinde işlem yapar. `nameState` bu yolları listeler: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` ve kullanım önbellekleri `~/.local/state/ccam/usage/<name>.json` ile `<name>.statusline.json`. `ccam remove work` bunları şu sırayla işler:

<!-- motion:accounts-remove -->

`ccam remove`, hesap çalışıyorsa önce reddeder, ardından sorar. Hesap uzantıları paylaşıyorsa `RewritePluginPaths`, grubun `installed_plugins.json` ve `known_marketplaces.json` dosyalarını yeniden yazar; çünkü Claude Code, bu hesabın klasöründen geçebilecek mutlak eklenti yollarını saklar. `profile.Trash` klasörü `~/.local/share/ccam/trash/<name>-<time>/` konumuna taşır, hesap `config.toml` dosyasından çıkarılır, aynı adlı sonraki bir hesap bunları devralmasın diye `nameState` yolları silinir ve kısayol bağlantısı kaldırılır. Kurtarmak için klasörü çöp kutusundan geri taşıyın.

`ccam rename` klasörü yeniden adlandırır, `plugins/` içindeki eklenti yollarını eski klasörden yenisine göre yeniden yazar, varsayılan ayarını taşır ve `nameState` yollarını yeni ada taşır. Kısayol komutu adını korur. Eklenti yollarının yeniden yazılması başarısız olursa yapılandırma kaydedilmediği için klasör eski adına geri döndürülür.

`ccam import` bir klasörü bulunduğu yerde, `config.toml` içinde `dir` olarak kaydeder. Onu hiçbir zaman taşımaz veya kopyalamaz: kopyalanmış bir yenileme tokenı, kopyaların birbirinin oturumunu kapatmasına neden olur. Klasör argümanı olmadan kaydedilmemiş `~/.claude` ve `~/.claude-*` klasörlerini arar ve her biri için sorar. Adlar klasörden gelir (`~/.claude` için `main`, `~/.claude-work` için `work`, `~/.claude-2` için `acct2`); adı türetilemeyen bir klasör, kullanılacak `--name` komutuyla birlikte listelenir. `CLAUDE_CONFIG_DIR` olmadan Claude Code, `.claude.json` dosyasını ev klasörünüzde `~/.claude` klasörünün yanında tutar; bu yüzden `~/.claude` klasörünü içe aktardığınızda ve kendine ait bir `.claude.json` dosyası yoksa ccam, `~/.claude/.claude.json` bağlantısını `~/.claude.json` dosyasına yönlendirerek oluşturur. Böylece düz `claude` ve `ccam main` aynı hesabı görür. İçe aktarılan hesaplar siz `ccam share` çalıştırana kadar hiçbir şey paylaşmaz ve kısayol oluşturulmaz.

`ccam completion bash|zsh`, hesap adlarını da tamamlayan bir tamamlama betiği yazdırır.

## Güvenlik ve uç durumlar

- `ccam rename` ve `ccam remove`, hesap çalışırken `work is running; quit its Claude Code session first` mesajıyla reddeder (`refuseRunning`), çünkü oturum klasöre yazmaya devam eder.
- İçe aktarılmış bir hesapta `ccam remove` yalnızca kaydı siler: `✓ main unregistered (folder ... left in place)`. `~/.claude` klasörünüz çöp kutusuna taşınmaz.
- `ccam add`, klasörü `profiles/` içinde zaten bulunan bir adı reddeder ve `ccam import` komutunu önerir.
- `ccam import`, zaten kayıtlı bir klasörü, klasör olmayan bir yolu ve kullanılmakta olan bir adı reddeder.
- `ccam logout` ve `ccam login`, hesabın çalışıp çalışmadığını kontrol etmez.
- Aktarma Claude Code'a ulaşamazsa, örneğin giriş sona erdikten sonra, nedeniyle birlikte `ccam: could not pass the sign-in to Claude Code:` yazdırır. Bir HTTP hatası genellikle kodun süresinin dolduğu anlamına gelir; girişi yeniden başlatın.

## İlgili

- [Claude Code'u başlatma](/features/launch)
- [Kısayol komutları](/features/shortcuts)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Birden fazla Claude Code hesabı](/guides/multiple-claude-code-accounts)
