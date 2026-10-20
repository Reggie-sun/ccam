---
title: Gruplarla paylaşım
description: Bir ccam grubu seçilen öğeleri hesapları arasında tek bir klasör üzerinden paylaşır; ccam group veriyi plan, yedek ve çakışma seçenekleriyle taşır.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

Her hesabın kendi klasörü vardır, ancak içindekilerin çoğunun ayrı olması gerekmez. Bir **grup** neyin, kimler tarafından paylaşılacağını söyler: gruptaki her hesap bu öğeleri tek bir klasöre, `~/.local/share/ccam/groups/<group>/` klasörüne bağlar; böylece hepsi aynı kopyayı okur ve yazar. Giriş token'ları hiçbir zaman paylaşılmaz.

Bir hesap birkaç grupta olabilir, ancak her öğe yalnızca tek bir gruptan gelir. Tipik bir kurulumda ayarlar ve uzantılar için bir grup, iki hesap arasındaki konuşma geçmişi için de ikinci bir grup bulunur:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Eklediğiniz ilk hesap `default` grubunu oluşturur; aksini söylemediğiniz sürece sonraki hesaplar bu gruba katılır.

## Kullanım

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Veri taşıyan her komut bir plan yazdırır, grubun kopyasından farklı olan her dosya için sorar ve herhangi bir şeyi değiştirmeden önce `Proceed? [y/N]` diye sorar:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

Aynı değişiklikler [kontrol panelinin](/features/dashboard) **Groups** sekmesinde de bulunur; bir hesap üzerinde `s` ise o hesabın gruplarını seçer.

## Öğeler

| Öğe | Kapsadığı |
| --- | --- |
| `mcp` | grubun `mcp.json` dosyası üzerinden MCP sunucuları ([Paylaşılan MCP sunucuları](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | genel `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` ve `history.jsonl`; böylece başka bir hesap `--continue` yapabilir |

Hesabın hiçbir grubunun paylaşmadığı bir öğe hesabın kendi klasöründe kalır.

## Komutlar ve seçenekler

| Komut | Ne yapar |
| --- | --- |
| `ccam group` veya `ccam group list` | grupları listeler |
| `ccam group show <group>` | grubun neyi paylaştığını, hesaplarını, klasörünü ve MCP sunucularını gösterir |
| `ccam group add <group> --items <items> [--members <accounts>]` | bir grup oluşturur; hesaplarının öğeleri gruba taşınır |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | grubun paylaştıklarını, içindeki her hesap için değiştirir |
| `ccam group join <group> <account>...` | hesapları gruba koyar |
| `ccam group leave <group> <account>...` | hesapları gruptan çıkarır; her biri bir kopya tutar |
| `ccam group rename <group> <new-name>` | yalnızca yapılandırmada yeniden adlandırır |
| `ccam group remove <group> [--leave]` | bir grubu kaldırır |

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: kopyası gruba ilk giren hesap | varsayılan hesap, ardından numara sırası |
| `--replace` | `join`: önce aynı öğeleri paylaşan gruplardan ayrılır | kapalı |
| `--fresh` | bir öğeyi paylaşmayı bırakan hesaplar bir kopyayla değil, o öğe olmadan başlar | kapalı |
| `--leave` | `remove`: önce grubun hesaplarını gruptan çıkarır | grubun hesapları varsa gerekir |
| `--prefer <group\|mine\|merge>` | çakışmalara nasıl karar verileceği | sorar; `--yes` ile her dosyanın varsayılanı |
| `-y`, `--yes` | sormadan çalışır; çakışmalar `--prefer` seçeneğini izler | kapalı |
| `--dry-run` (genel) | planı yazdırır, hiçbir şeyi değiştirmez | kapalı |

Hesaplar adla veya numarayla, öğeler virgülle ayrılmış bir liste olarak verilir. `ccam add`, `--group <group>` (tekrarlanabilir) veya `--no-group` alır.

## Nasıl çalışır

İlgili dosyalar:

- `~/.config/ccam/config.toml`: `[[group]]` girdileri: `name`, `items`, `members` ve klasörün adı farklıysa `dir`.
- `~/.local/share/ccam/profiles/<name>/`: hesabın `CLAUDE_CONFIG_DIR` klasörü.
- `~/.local/share/ccam/groups/<folder>/`: paylaşılan kopyalar.

`config.toml` okunduğunda her hesabın grupları öğe başına tek bir klasöre çözümlenir. Bu, beş hesap için yaklaşık 25 µs sürer; bu yüzden bir başlatma yine yalnızca `config.toml` dosyasını okur. Bir hesabın bir öğeyi iki gruptan alacağı bir yapılandırma hatadır ve düzeltilene kadar hiçbir şey başlamaz.

Bir değişiklik yapılandırmanın bir kopyası üzerinde çalışır:

1. Komut kopyayı düzenler (`AddGroup`, `SetGroupItems`, `JoinGroup` …); bu işlevler kuralları kontrol eder.
2. `internal/group` eskiyle yeniyi karşılaştırır ve öğeleri değişen her hesap için `share.Build` ile bir plan oluşturur. Planlar art arda oluşturulur ve her biri öncekilerin grup klasörlerine ne koyacağını görür; böylece kopyasını gruba yalnızca ilk hesap taşır. İlk kopya, varsa grup klasörünün kendisinden, sonra `--from` hesabından, sonra varsayılan hesaptan, sonra numara sırasıyla diğer hesaplardan gelir.
3. Bu hesaplardan biri çalışıyorsa ya da Geliştirici Modu olmayan Windows'ta bir dosyanın bağlantı gerektirmesi durumunda ccam reddeder.
4. Çakışmalara soruyla veya `--prefer` ile karar verilir, plan yazdırılır ve ccam onay ister.
5. Hiçbir şey taşınmadan önce yeni yapılandırma kaydedilir. Bir adım yarıda durursa `ccam doctor --fix` ve yeniden çalıştırma, planlarını kaydedilmiş gruplara doğru oluşturur ve değişikliği tamamlar.
6. Planlar sırayla çalışır. Her adım veriyi gruba veya yedek klasörüne taşır; hiçbir şey silinmez.

`share.Layout`, her öğenin kapsadığı yolların bir izin listesidir. Listede olmayan hiçbir şey bir gruba ulaşamaz. `share.Never`, `.credentials.json` ve `.claude.json` dosyalarını adlandırır; bunlardan biri düzende görünürse `Build` dahili bir hatayla durur. `mcp` öğesinin yolu yoktur: sunucu listesi `.claude.json` içinde bulunur; bu yüzden [MCP paylaşımı](/features/mcp) grubun `mcp.json` dosyası üzerinden çalışır.

`side` geçmişi paylaşmaya başladığında ve grupta zaten içerik bulunduğunda tek bir klasöre, `projects/` klasörüne olanlar şunlardır:

<!-- motion:share-on -->

Grupta o yolda henüz hiçbir şey yoksa hesap onu paylaşan ilk hesaptır ve klasörü olduğu gibi oraya taşınır. Aksi halde grupta olmayan her dosya gruba taşınır; her iki tarafta da bulunan dosyalar hesabın klasöründe kalır, böylece grubun sürümü korunur. Ardından hesabın klasörü hesap klasörünün içindeki `.ccam-backup/<time>/` konumuna gider ve yerine bir bağlantı konur. Konuşma dosyaları oturum kimliklerine göre adlandırılır; bu nedenle birden fazla hesabın geçmişi çakışma olmadan birleşir. `plugins/` kayıtları hesap klasöründen geçen mutlak yollar saklar; ccam bunları grubun klasörünü gösterecek şekilde değiştirir ve grupta zaten bir kayıt varsa hesabın girdilerini onunla birleştirir.

Ayrılmak bağlantıyı tersine çevirir:

<!-- motion:share-off -->

ccam bağlantıyı yalnızca gruplar klasörünü gösteriyorsa kaldırır, ardından grubun mevcut içeriğini hesaba kopyalar. `--fresh` ile kopyalama atlanır. Grup, içeriğini diğer hesaplar için korur. Bir öğeyi klasörü henüz boş olan bir gruba taşıyan hesap gördüğünün bir kopyasını getirir; böylece dosyalarını korur. MCP sunucuları da aynı şekilde gelir. MCP paylaşmayı bırakan bir hesap, ccam'in onun `.claude.json` dosyasına kopyaladığı sunucuları kaybeder (bir yedek tutulur), çünkü bunlar gruptan gelmiştir.

Her iki tarafta da bulunan ve içeriği farklı olan tek bir dosya bir çakışmadır:

<!-- motion:share-conflict -->

`history.jsonl`, önce birleştirme olmak üzere `merge by line`, `keep group` ve `use mine` sunar. Diğer `.json` dosyaları `keep group`, `use mine` ve `merge JSON` sunar. Geri kalan her şey `keep group` ve `use mine` sunar. `merge JSON` nesneleri özyinelemeli olarak birleştirir, dizileri yinelenenler olmadan birleştirir ve iki taraf aynı anahtara farklı değerler verdiğinde grubun değerini korur; böylece diğer hesaplar hiçbir değişiklik görmez. `keep group` dışındaki her seçenek önce grubun kopyasını hesabın `.ccam-backup/<time>/group/` klasörüne kaydeder. Kendi dosyanız her zaman yedeğe gider.

Her başlatmada `share.Check`, düzendeki her yolu yalnızca `lstat` ve `readlink` kullanarak hesabın gruplarıyla karşılaştırır. Eksik bir bağlantı anında oluşturulur. Bağlantı olması gereken yerde gerçek bir dosya, yanlış yere yönelen bir bağlantı veya bayat bir bağlantı başlatmayı durdurur ve `ccam doctor --fix` komutunu gösterir.

### Yeniden adlandırma, kaldırma ve eski biçim

- `ccam group rename` yalnızca yapılandırmayı değiştirir. Klasör adını korur ve grup onu `dir` olarak kaydeder; böylece hiçbir şey taşınmaz ve hiçbir şey kaybolamaz.
- `ccam group remove`, `--leave` eklemedikçe hesapları olan bir grubu reddeder. Bu durumda her hesap bir kopya tutar ve klasör, başka bir grup onu kullanmıyorsa `~/.local/share/ccam/trash/` klasörüne gider.
- Gruplardan önce `config.toml`, her hesap için bir `[profile.share]` tablosu saklıyordu. ccam bu biçimi hiçbir şeyi taşımadan gruplar olarak okur. İlk kaydetmede yeni biçimi yazar ve eski dosyayı `config.toml.bak` olarak tutar. Aynı hesapların paylaştığı öğeler, klasörün adını taşıyan tek bir grup olur; geri kalanlar `default-history` gibi adlar alır ve aynı klasörü kullanmaya devam eder.

### Windows'ta

Windows sembolik bağlantıları yalnızca Geliştirici Modu'nda (Ayarlar > Sistem > Geliştiriciler için) veya bir yönetici için oluşturur. Bu hak olmadan paylaşılan bir klasör, Claude Code için aynı şekilde çalışan bir dizin bağlantı noktasına (junction) dönüşür. Paylaşılan bir dosya (`settings.json`, `CLAUDE.md`, `history.jsonl`) yalnızca sembolik bağlantı olabilir; bu yüzden bağlantı gerektiren bir grup değişikliği hiçbir şey taşınmadan önce reddedilir.

## Güvenlik ve uç durumlar

- **Çalışan hesaplar.** Claude Code, değişikliğin taşıdığı hesaplardan biriyle çalışırken değişiklik reddedilir, çünkü oturum bu dosyalara yazmaya devam eder. Ayrılma ve katılma yalnızca adı verilen hesabı ilgilendirir. `--dry-run` yine de planı yazdırır.
- **Yedekler.** Her plan, hesap klasöründe kendi `.ccam-backup/<time>/` klasörünü alır. ccam bunları hiçbir zaman silmez.
- **Kendi bağlantılarınız.** Bir yol zaten sizin oluşturduğunuz bir bağlantıysa, örneğin bir dotfiles deposuna, ve grup boşsa ccam bağlantının hedefini gruba kopyalar ve bağlantının kendisini yedeğe taşır. Hedefe dokunulmaz.
- **İçe aktarılmış klasörler.** `ccam import` ile yerinde kaydedilmiş bir hesap (örneğin `~/.claude`) yerinde değişir; böylece düz `claude` da değişikliği görür. Plan bunu belirtir.
- **Bir klasörde iki hesap.** Geçmiş paylaşan hesaplar aynı dosyalara yazar. Biri zaten çalışırken aynı klasörde ikinci bir hesap başladığında ccam uyarır ve sorar.
- **Elle düzenlemeler.** Hesap olmayan bir üye yok sayılır ve `ccam doctor` tarafından bildirilir. `local`, `on`, `off` ve `none` adları grup adı olamaz.

## İlgili

- [Kontrol paneli](/features/dashboard)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Paylaşılan MCP sunucuları](/features/mcp)
- [Mevcut bir kurulumu kopyalama](/features/migrate)
- [Tüm komutlar](/docs/commands)
