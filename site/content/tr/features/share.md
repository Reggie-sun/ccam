---
title: Hesaplar arasında paylaşım
description: ccam share ayarları, CLAUDE.md dosyasını, uzantıları ve geçmişi plan, yedek ve çakışma seçenekleriyle grup klasörüne nasıl bağlar, off nasıl geri kopyalar.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Her hesabın kendi klasörü vardır, ancak içindekilerin çoğunun ayrı olması gerekmez. `ccam share`, paylaşımı her seferinde bir hesap ve bir öğe için açar veya kapatır: hesabın dosyası veya klasörü bir grup klasörüne taşınır ve yerine sembolik bir bağlantı konur; böylece gruptaki her hesap aynı kopyayı okur ve yazar. Bir hesabın neyi paylaştığını değiştirmek veya birkaç hesabı kendilerine ait bir gruba koymak için `ccam add` sonrasında kullanın.

## Kullanım

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Öğe seçenekleri olmadan `ccam share` mevcut durumu yazdırır:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Seçeneklerle ccam bir plan yazdırır, grubun kopyasından farklı olan her dosya için sorar ve herhangi bir şeyi değiştirmeden önce `Proceed? [y/N]` diye sorar. Yeni bir hesap konuşma geçmişi dışında her şeyi paylaşır; [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing) sayfası her öğeyi kullanıcı açısından açıklar.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | MCP sunucuları | değişmez |
| `--settings <on\|off\|group>` | `settings.json` | değişmez |
| `--memory <on\|off\|group>` | genel `CLAUDE.md` | değişmez |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | değişmez |
| `--history <on\|off\|group>` | `projects/` ve `history.jsonl` | değişmez |
| `--fresh` | paylaşımı kapatırken grubun içeriğini kopyalamak yerine boş başlar | kapalı |
| `--prefer <group\|mine\|merge>` | çakışmalara nasıl karar verileceği | sorar; `--yes` ile her öğenin varsayılanı |
| `-y`, `--yes` | sormadan çalışır; çakışmalar `--prefer` seçeneğini izler | kapalı |
| `-h`, `--help` | `share` yardımı | |
| `--dry-run` (genel) | dikkate alınır: çakışmaları sorar, planı yazdırır, hiçbir şeyi değiştirmez | kapalı |

`on` değeri `default` grubu, `off` yerel (paylaşılmayan) anlamına gelir; diğer her değer bir grup adıdır. Grup adları `~/.local/share/ccam/groups/` altında klasör adları olur; bu yüzden `ccam share side --history ../x`, `invalid group name: "../x"` ile başarısız olur. Hesap adla veya numarayla verilir ve yalnızca belirttiğiniz öğeler değişir.

## Nasıl çalışır

İlgili dosyalar:

- `~/.config/ccam/config.toml`: her hesabın `[profile.share]` tablosu, öğe başına bir değer.
- `~/.local/share/ccam/profiles/<name>/`: hesabın `CLAUDE_CONFIG_DIR` klasörü.
- `~/.local/share/ccam/groups/<group>/`: paylaşılan kopyalar.

1. `internal/cli/share.go` yapılandırmayı okur, her `--<item>` değerini `config.Share.Set` ile hesabın ayarlarının bir kopyasına uygular ve hesap çalışıyorsa reddeder.
2. `share.Build` mevcut ayarları hedefle karşılaştırır ve adımları listeler. Diski yalnızca okur: `share.Layout` içindeki her yol için `os.Lstat` ve `os.Readlink`. Plan bir işlemler listesidir: boş bir grup yuvasına taşıma (adopt), grupta olmayan dosyaları taşıma (merge dir), çözümlenmiş bir çakışmayı gruba yazma, yedekleme, bağlantıyı kaldırma, bağlama, geri kopyalama, eklenti yollarını yeniden yazma ve eklenti kaydını birleştirme.
3. Çakışmalara soruyla veya `--prefer` ile karar verilir.
4. `Plan.Describe` planı yazdırır ve ccam onay ister.
5. Hiçbir şey taşınmadan önce hedef ayarlar `config.toml` dosyasına yazılır. `Plan.Apply` yarıda durursa `ccam doctor --fix` ve yeniden çalıştırma, planlarını kaydedilmiş hedefe doğru oluşturur ve değişikliği tamamlar.
6. `Plan.Apply` adımları sırayla çalıştırır. Her adım veriyi gruba veya yedek klasörüne taşır; hiçbir şey silinmez.

`share.Layout`, her öğenin kapsadığı yolların bir izin listesidir. Listede olmayan hiçbir şey bir gruba ulaşamaz. `share.Never`, `.credentials.json` ve `.claude.json` dosyalarını adlandırır; bunlardan biri düzende görünürse `Build` dahili bir hatayla durur. `mcp` öğesinin yolu yoktur: sunucu listesi `.claude.json` içinde bulunur; bu yüzden [MCP paylaşımı](/features/mcp) bunun yerine grubun `mcp.json` dosyası üzerinden çalışır ve `--mcp` yalnızca ayarı değiştirir.

`side` geçmişi açtığında ve grupta zaten içerik bulunduğunda tek bir klasöre, `projects/` klasörüne olanlar şunlardır:

<!-- motion:share-on -->

Grupta o yolda henüz hiçbir şey yoksa hesap onu paylaşan ilk hesaptır ve klasörü yedek alınmadan olduğu gibi oraya taşınır. Aksi halde `mergeDir`, grupta olmayan her dosyayı gruba taşır; her iki tarafta da bulunan dosyalar hesabın klasöründe kalır, böylece grubun sürümü korunur. Ardından hesabın klasörü hesap klasörünün içindeki `.ccam-backup/<time>/` konumuna gider ve `os.Symlink` yerine bir bağlantı koyar. Konuşma dosyaları oturum kimliklerine göre adlandırılır; bu nedenle birden fazla hesabın geçmişi çakışma olmadan birleşir.

Taşıma `rename(2)` kullanır ve dosya sistemleri arasında kopyalayıp silmeye geri döner. Taşınan ağacın içinde, ağacın dışını gösteren göreli bir bağlantı mutlak bir bağlantı olarak yeniden yazılır (`fsutil.PortableLink`); böylece yeni yerden de çözümlenir. `plugins/` kayıtları (`installed_plugins.json`, `known_marketplaces.json`) hesap klasöründen geçen mutlak yollar saklar; ccam bunları grubun klasörünü gösterecek şekilde değiştirir ve grupta zaten bir kayıt varsa hesabın girdilerini onunla birleştirir.

Bir öğeyi kapatmak bağlantıyı tersine çevirir:

<!-- motion:share-off -->

ccam bağlantıyı yalnızca gruplar klasörünü gösteriyorsa kaldırır, ardından grubun mevcut içeriğini `CopyTree` ile hesaba kopyalar. `--fresh` ile kopyalama atlanır ve hesap o yolda hiçbir şey olmadan başlar. Grup, içeriğini diğer hesaplar için korur. Bir öğeyi bir gruptan diğerine taşımak, bağlantıyı kaldırır ve yeni gruba bağlar; eski grup olduğu gibi kalır. Plan bunu bir `Note:` satırında belirtir.

Her iki tarafta da bulunan ve içeriği farklı olan tek bir dosya bir çakışmadır:

<!-- motion:share-conflict -->

`share.NewConflict`, dosyaya uyan seçenekleri listeler. `history.jsonl`, önce birleştirme olmak üzere `merge by line`, `keep group` ve `use mine` sunar. Diğer `.json` dosyaları `keep group`, `use mine` ve `merge JSON` sunar. Geri kalan her şey `keep group` ve `use mine` sunar. İlk seçenek varsayılandır:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) nesneleri özyinelemeli olarak birleştirir, dizileri yinelenenler olmadan birleştirir ve iki taraf aynı anahtara farklı değerler verdiğinde grubun değerini korur; böylece diğer hesaplar davranışta hiçbir değişiklik görmez. `merge by line` grubun satırlarını korur, ardından grupta olmayan sizin satırlarınızı ekler. `keep group` dışındaki her seçenek önce grubun kopyasını hesap klasöründeki `.ccam-backup/<time>/group/` içine kaydeder, ardından sonucu `fsutil.AtomicWrite` ile yazar. Kendi dosyanız her zaman yedeğe gider.

`--prefer group`, `mine` veya `merge` her çakışmaya sormadan karar verir; `merge`, dosyanın sunduğuna göre `merge JSON` veya `merge by line` anlamına gelir ve ikisini de sunmayan bir dosya varsayılanı korur. Tek başına `--yes` her dosyanın varsayılanını alır.

Her başlatmada `share.Check`, düzendeki her yolu yalnızca `lstat` ve `readlink` kullanarak ayarla karşılaştırır; bu, onu başlatma yolu için yeterince ucuz tutar. Eksik bir bağlantı anında oluşturulur, çünkü kaybedilecek bir şey yoktur. Bağlantı olması gereken yerde gerçek bir dosya, yanlış yere yönelen bir bağlantı veya bayat bir bağlantı başlatmayı durdurur:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix`, `share.Build` fonksiyonunu kayıtlı ayarları her iki tarafa vererek çalıştırır; bu, her adımı diskte bulunandan türetir, çakışmalarda grubun içeriğini korur ve değiştirdiği her şeyi yedekler.

## Güvenlik ve uç durumlar

- **Çalışan hesaplar.** Claude Code o hesapla çalışırken ccam paylaşımı değiştirmeyi reddeder (`account work is running; quit it first, then change its sharing`), çünkü oturum taşınan dosyalara yazmaya devam eder ve bu arada eklediği satırlar kaybolur. `--dry-run` yine de planı yazdırır.
- **Yedekler.** Her plan, hesap klasöründe kendi `.ccam-backup/<time>/` klasörünü alır; iki plan aynı saniyede başlarsa numaralı bir sonek eklenir. ccam bunları hiçbir zaman silmez.
- **Kısmi hatalar.** Bir adım başarısız olursa ccam hatayı ve orijinallerin nerede olduğunu yazdırır ve işlemi `ccam doctor --fix` ile tamamlamanızı söyler. Hedef zaten kaydedilmiştir; bu yüzden onarım ona doğru devam eder.
- **Kendi bağlantılarınız.** Bir yol zaten sizin oluşturduğunuz bir bağlantıysa, örneğin bir dotfiles deposuna, ve grup boşsa ccam bağlantının hedefini gruba kopyalar ve bağlantının kendisini yedeğe taşır. Hedefe dokunulmaz.
- **Giriş dosyaları.** `.credentials.json` ve `.claude.json` hiçbir zaman bir gruba bağlanmaz, taşınmaz veya kopyalanmaz.
- **Bir klasörde iki hesap.** Geçmiş paylaşan hesaplar aynı dosyalara yazar. Aynı geçmiş grubundan ikinci bir hesap, birinin zaten çalıştığı bir klasörde başladığında ccam uyarır ve devam etmeden önce sorar.

## İlgili

- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Paylaşılan MCP sunucuları](/features/mcp)
- [Mevcut bir kurulumu kopyalama](/features/migrate)
- [Tüm komutlar](/docs/commands)
