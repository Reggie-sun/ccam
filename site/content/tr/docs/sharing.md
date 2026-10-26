---
title: MCP sunucularını, ayarları ve geçmişi paylaşma
description: ccam grupları MCP sunucularını, settings.json, CLAUDE.md, becerileri ve konuşma geçmişini Claude Code hesapları arasında, girişleri ayrı tutarak paylaşır.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 2d8029ad6ec8
---

Her hesap kendi giriş bilgilerini korur. Geri kalan her şey **gruplar** üzerinden paylaşılabilir: bir grup paylaştığı öğeleri ve içindeki hesapları belirtir; ccam da bu hesapların her birinin dosyalarını grubun klasörüne bağlar, böylece bir hesapta yapılan değişiklik hepsinde görünür.

## Neler paylaşılabilir

| Öğe | Paylaşılan |
| --- | --- |
| `mcp` | MCP sunucuları: grubun `mcp.json` dosyası, başlangıçta Claude Code'a iletilir |
| `settings` | `settings.json`: izinler, hook'lar, model, durum satırı |
| `memory` | genel `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | konuşma geçmişi: `projects/` ve `history.jsonl` |

İlk hesabınız, konuşma geçmişi dışında her şeyi paylaşan `default` grubunu oluşturur; sonraki hesaplar bu gruba katılır.

Giriş dosyaları `.credentials.json` ve `.claude.json` hiçbir zaman paylaşılmaz. ccam yalnızca sabit bir listedeki yolları bağlar ve bu iki dosya o listede yoktur.

Bu sayfa günlük işleri anlatır. Tüm seçenekler özellik sayfalarında açıklanır: [gruplarla paylaşım](/features/share), [paylaşılan MCP sunucuları](/features/mcp) ve [mevcut bir kurulumu kopyalama](/features/migrate).

## Gruplar

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Bir hesap birkaç grupta olabilir, ancak her öğe yalnızca tek bir gruptan gelir. Bu yüzden geçmişi `work` ile `side` arasında paylaşıp diğer hesaplarınızla paylaşmamak için geçmişi `default` dışında tutun ve yukarıdaki gibi bu ikisine kendilerine ait bir grup verin. `ccam share work` bir hesabın neyi, hangi grup üzerinden paylaştığını gösterir. [Kontrol paneli](/features/dashboard) aynı şeyi Groups sekmesinde yapar.

## Bir hesap paylaşmaya başladığında

Bir hesap bir öğeyi paylaşmaya başladığında (bir gruba katıldığında ya da grubu o öğeyi eklediğinde) ccam:

1. her hesap için taşıyacağı, kopyalayacağı veya bağlayacağı her dosyanın planını yazdırır ve başlamadan önce onay ister;
2. grup henüz boşsa hesabın dosyalarını grup klasörüne taşır, değilse grupta henüz olmayan dosyaları ekler;
3. aynı dosya farklıysa grubun sürümünü mü koruyacağınızı, sizinkini mi kullanacağınızı yoksa JSON'u mu birleştireceğinizi sorar;
4. orijinalinizi hesap klasörünün içindeki `.ccam-backup/<time>/` altında saklar;
5. orijinali gruba yönelen bir bağlantıyla değiştirir.

Konuşma dosyaları oturum kimliklerine göre adlandırılır; bu nedenle birden fazla hesabın geçmişi çakışma olmadan birleşir.

Claude Code o hesapla çalışırken, nasıl başlatılmış olursa olsun (arka plan oturumları da dahil), ccam hesabı değiştirmez, çünkü oturum taşınan dosyalara yazmaya devam eder. Önce oturumdan çıkın. ccam kaydetmeden hemen önce yeniden denetler; böylece ccam yanıtınızı beklerken başlatılan bir oturumu da yakalar.

## Bir hesap paylaşmayı bıraktığında

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Bağlantı kaldırılır ve hesap, grubun mevcut içeriğinin kendine ait bir kopyasını alır; grubun MCP sunucuları, hesabın kendi sunucularının yanına, kendi `.claude.json` dosyasına girer. Bunun yerine boş başlamak için `--fresh` ekleyin; plan o zaman yüklenmesi duracak MCP sunucularını listeler. `~/.claude` gibi içe aktarılmış bir klasör her zaman bir kopya tutar. Grup, diğer hesaplar için içeriğini korur.

## Mevcut bir kurulumu kopyalama

Claude Code'u ccam'den önce kullandıysanız ayarlarınız `~/.claude` içindedir. `ccam migrate` bunları bir gruba kopyalar; böylece bu öğeleri paylaşan her hesap onlara sahip olur:

```sh
ccam migrate
```

Terminalde ccam önce öğelerin bir kontrol listesini gösterir. Konuşma geçmişi dışındaki her şey işaretlidir; bir öğeyi değiştirmek için Space'e, planı görmek için Enter'a basın. Klasörde karşılığı olmayan öğeler `not found` olarak işaretlenir.

Kopyaladıkları:

- `settings.json` ve genel `CLAUDE.md`;
- `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` ve `rules/` uzantıları;
- kullanıcı MCP sunucularınız, yani `~/.claude.json` içindeki `mcpServers` girdisi, grubun `mcp.json` dosyasına;

- konuşma geçmişi (`projects/` ve `history.jsonl`), ancak yalnızca `--history` eklediğinizde.

Kaynak klasör yalnızca okunur: içindeki hiçbir şey taşınmaz, bağlanmaz veya değiştirilmez ve düz `claude` eskisi gibi çalışmaya devam eder. Giriş dosyaları `.credentials.json` ve `.claude.json` asla gruba kopyalanmaz.

OAuth ile giriş yapan MCP sunucuları tokenlarını sunucu listesinde değil, giriş bilgilerinde tutar. ccam bu girişleri kopyalamaz: kopyalanan bir giriş aynı yenileme tokenını taşır ve ilk yenileme, `~/.claude` dahil, onu tutan diğer herkesin o sunucudaki oturumunu kapatır. Plan ve özet, bu sunuculara kendileri `/mcp` ile giriş yapması gereken hesapları adlarıyla belirtir.

ccam herhangi bir şeyi kopyalamadan önce planı yazdırır ve onay ister:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Group default: personal
To share them:
  ccam group join default work
```

Sonunda ccam grubun hesaplarını listeler. Grup öğelerden birini henüz paylaşmıyorsa ya da bazı hesaplar grupta değilse bunu değiştiren `ccam group` komutlarını yazdırır.

Grupta bir dosyanın farklı bir kopyası zaten varsa ccam, paylaşımı açarken olduğu gibi sorar: grubun kopyasını korumak mı, sizinkini kullanmak mı. `settings.json` ve diğer JSON dosyaları birleştirilebilir de, `history.jsonl` ise satır satır birleştirilebilir. Grubun önceki kopyaları grup klasörünün içindeki `.ccam-backup/<time>/` altında saklanır.

| Seçenek | Etki |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Yalnızca belirtilen öğeleri, kontrol listesi olmadan kopyalar. Terminal dışında ve hiçbir öğe seçeneği verilmeden ccam, konuşma geçmişi dışındaki her şeyi kopyalar. |
| `-g`, `--group <group>` | `default` yerine bu gruba kopyalar |
| `-y`, `--yes` | Sormadan kopyalar; bir dosya farklıysa grubun kopyasını korur |

Hiçbir şey kopyalamadan planı görmek için `--dry-run` seçeneğini komuttan önce yazın: `ccam --dry-run migrate`. Başka bir klasörden kopyalamak için onu belirtin: `ccam migrate ~/.claude-work --group work`.

## MCP sunucuları

Claude Code, MCP sunucularını giriş bilgileriyle birlikte `.claude.json` içinde saklar; bu yüzden bu dosya paylaşılamaz. ccam, paylaşılan sunucuları grubun `mcp.json` dosyasında tutar ve Claude Code'u `--mcp-config=<group>/mcp.json` ile başlatır. Kendi `--mcp-config` seçenekleriniz bunun üzerine uygulanmaya devam eder.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Bir oturum içinde `claude mcp add` ile eklediğiniz sunucu yalnızca o hesap için ve içinde bulunduğunuz klasör için kaydedilir (Claude Code'un yerel kapsamı; `-s user` onu hesap için her klasörde kullanılabilir yapar). O hesabı bir sonraki başlatışınızda ccam, sunucunun adını klasörüyle birlikte ve onu gruba taşıyan komutu gösterir:

```sh
ccam mcp import work gh
```

Grupta aynı adla farklı bir sunucu zaten varsa ccam onu atlar ve bunu bildirir. `--force` grubun sunucusunu değiştirir ve eski `mcp.json` dosyasını grubun `.ccam-backup/` klasöründe saklar.

OAuth ile giriş yapan MCP sunucuları, tokenlarını hesabın giriş bilgileriyle birlikte saklar; bu yüzden her hesapta `/mcp` ile onlara bir kez giriş yapın. `ccam ext` ve `ccam doctor`, başka bir hesap giriş yapmışken bir hesabın giriş yapmadığı sunucuları adlarıyla belirtir.

Sunucuların `claude mcp list` ve `/mcp` içinde sıradan kullanıcı sunucuları olarak görünmesini tercih ediyorsanız `~/.config/ccam/config.toml` dosyasının en üstünde, ilk `[[group]]`, `[[profile]]` veya `[[folder]]` satırının üstünde `mcp_strategy = "sync"` ayarını yapın. ccam bu durumda onları başlatmadan önce, Claude Code'un kendi dosya kilidi altında ve bir yedekle birlikte, her hesabın `.claude.json` dosyasına yazar. Hesabın bir grup sunucusuyla aynı adı taşıyan kendi sunucusu varsa ccam hesabın sunucusunu korur ve bunu söyler.

## Aynı klasörde iki hesap

Konuşma geçmişini paylaşan hesaplar aynı dosyalara yazar. Aynı gruptan başka bir hesabın zaten çalıştığı bir klasörde ikinci bir Claude Code oturumu başlatırsanız ccam sizi uyarır ve devam edilip edilmeyeceğini sorar.

## Bir şey bozulduğunda

```sh
ccam doctor
ccam doctor --fix
```

Paylaşılan bir yol sıradan bir dosyaya dönüşmüşse (örneğin bir yedeği elle geri yükledikten sonra), ccam Claude Code'u başlatmadan önce durur ve sizi `ccam doctor --fix` komutuna yönlendirir. Bu düzeltme yolu yeniden gruba bağlar: grubun bir kopyası yoksa dosya grubun kopyası olur; aksi halde grup kendi kopyasını korur ve sizinki, `--fix` komutunun adını verdiği bir yedeğe gider. Yarıda kalmış bir grup değişikliğini de `ccam doctor --fix` tamamlar.
