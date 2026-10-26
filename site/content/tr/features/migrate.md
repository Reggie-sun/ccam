---
title: Mevcut bir kurulumu kopyalama
description: Mevcut bir ~/.claude klasörünün ayarlarını, CLAUDE.md dosyasını, uzantılarını ve MCP sunucularını ccam migrate ile bir ccam grubuna kopyalayın.
updated: 2026-10-26
order: 6
translation: machine
sourceHash: f835fbe84ade
---

Claude Code'u ccam'den önce kullandıysanız kurulumunuz `~/.claude` ve `~/.claude.json` içindedir. `ccam migrate` bunu bir paylaşım grubuna kopyalar; böylece bu öğeleri grup üzerinden paylaşan her hesap, zaten sahip olduğunuz kurulumla başlar. Kaynak klasör yalnızca okunur ve düz `claude` eskisi gibi çalışmaya devam eder. İlk `ccam add` komutunuzdan sonra bir kez çalıştırın; `~/.claude.json` içinde MCP sunucuları varsa ilk `ccam add` ipucunda onları adlarıyla belirtir.

## Kullanım

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

Terminalde, öğe seçenekleri ve `--yes` olmadan ccam önce bir kontrol listesi gösterir:

```console
 Copy into group default

 From /home/me/.claude
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Space (veya `x`) imlecin altındaki öğeyi değiştirir, ok tuşları veya `j`/`k` hareket ettirir, Enter planı gösterir; Esc, `q` veya Ctrl+C ise `Cancelled.` yazdırır. Kaynakta karşılığı olmayan öğeler `not found` olarak işaretlenir. Her şeyin işaretini kaldırırsanız ccam `Nothing selected.` yazdırır ve durur.

Ardından ccam planı yazdırır ve kopyalamadan önce onay ister:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Bir grup yalnızca kendi hesaplarını ve yalnızca paylaştığı öğeleri etkiler; bu yüzden özet, grubun hesaplarını listeler. Grup kopyalanan öğelerden birini paylaşmıyorsa ya da bazı hesaplar grupta değilse ccam, bunu değiştiren `ccam group set <group> --add <items>` ve `ccam group join <group> <accounts>` komutlarını yazdırır.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `[folder]` | Kopyalanacak Claude Code klasörü | `~/.claude` |
| `--settings` | `settings.json` dosyasını kopyalar | kapalı |
| `--memory` | Genel `CLAUDE.md` dosyasını kopyalar | kapalı |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` klasörlerini kopyalar | kapalı |
| `--mcp` | Kullanıcı MCP sunucularını (`.claude.json` içindeki `mcpServers`) kopyalar | kapalı |
| `--history` | Konuşma geçmişini (`projects/`, `history.jsonl`) kopyalar | kapalı |
| `-g`, `--group <group>` | Kopyalanacak paylaşım grubu | `default` |
| `-y`, `--yes` | Sormadan kopyalar; çakışmada grubun kopyası korunur | kapalı |
| `--dry-run` (genel) | Planı yazdırır ve herhangi bir sorudan veya kopyalamadan önce durur | kapalı |

Terminalde hiçbir öğe seçeneği vermezseniz öğeleri kontrol listesinden seçersiniz. Terminal dışında veya `--yes` ile ccam, konuşma geçmişi dışındaki her şeyi kopyalar. `projects/` büyük olabileceği için geçmiş varsayılan olarak dışarıda bırakılır.

`--group on`, `default` anlamına gelir. `--group off` ve `--group local`, `--group takes a group name, not "off"` ile reddedilir. `ccam migrate`, `--dry-run` seçeneğini destekler; kuru çalıştırmada da kontrol listesi görünür, böylece bir seçimi önizleyebilirsiniz.

## Nasıl çalışır

ccam herhangi bir şeyi kopyalamadan önce eksiksiz bir plan oluşturur. Öğeleri seçer (`internal/cli/migrate.go`), ardından kaynağı ve grubu okur ve planı oluşturur (`migrate.Build`). Siz onaylayıp `Plan.Apply` çalışana kadar hiçbir şey yazılmaz.

1. **Kaynağı bulma.** ccam eksik olan, klasör olmayan veya bilinen hiçbir Claude Code girdisi içermeyen bir yolu reddeder (`checkSource`). Yalnızca MCP sunucuları için, `~/.claude.json` içinde sunucular olduğu sürece boş ya da eksik bir `~/.claude` sorun değildir. ccam klasöre ait `.claude.json` dosyasını bulur (`migrate.ClaudeJSON`): klasörün içinde ya da içinde olmayan `~/.claude` için yanındaki `~/.claude.json`. `CLAUDE_CONFIG_DIR` ayarlı olmadığında Claude Code onu orada tutar.
2. **Öğeleri yollara eşleme.** Her yol, grupların kullandığı izin listesinin aynısı olan `share.Layout` içinden gelir. MCP sunucuları, dosyanın tamamını ayrıştırmak yerine `profile.TopLevelValue` tarayıcısını kullanan `mcp.LocalServers` ile `.claude.json` dosyasından okunur.
3. **Taşımadan kopyalama.** Dosyalar `fsutil.AtomicWrite` veya `fsutil.CopyFile` ile `~/.local/share/ccam/groups/<group>/` içine gider; bu ikisi hedefin yanına geçici bir dosya yazıp onu yeniden adlandırır, böylece kesintiye uğrayan bir kopyalama asla yarım bir dosya bırakmaz. MCP sunucuları, kilidi altında grubun `mcp.json` dosyasına gider: başka bir ccam sürecinin bu arada eklediği bir sunucu korunur ve [Remote Control](/guides/remote-control) için yapılmış kopyalar ardından güncellenir.

Her öğe bir kaynak yolu bir grup yoluna eşler:

<!-- motion:migrate-items -->

İki geçmiş satırı kesikli çizgiyle gösterilir, çünkü yalnızca `--history` ile kopyalanırlar. `.credentials.json` ve `.claude.json`, `share.Never` içinde listelenir ve hiçbir zaman gruba ulaşmaz: bunlardan biri düzende görünürse ccam dahili bir hatayla durur (`Build`) ve klasör taraması bunları adlarına göre atlar.

ccam her yola bir durum verir; her plan satırı bunu açıklar (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: kaynakta yok veya klasör beklenen yerde bir dosya var. Öğe başına bir kez `not found, skipped` olarak listelenir.
- **Linked**: kaynak yol zaten grubun kopyasına çözümlenir; örneğin kaynak bir ccam profil klasörü olduğunda.
- **Same**: grupta zaten aynı içerik var.
- **New**: grupta yok, bu yüzden bütünüyle kopyalanır.
- **Merge**: her iki tarafta da bir klasör. ccam yalnızca grupta olmayan dosyaları ekler (`copyMissing`); farklı olan dosyalar gruptaki haliyle kalır ve plan `group copies kept` der. `plugins/` için `installed_plugins.json` ve `known_marketplaces.json` kayıtları birleştirilir ve mutlak eklenti yolları grubun klasörünü gösterecek şekilde yeniden yazılır.
- **Conflict**: her iki tarafta da farklı içerikli bir dosya. Dosya başına seçersiniz: `keep group`, `use mine`, ayrıca `.json` dosyaları için `merge JSON` (çakışan değerlerde grubun kazandığı derin bir birleştirme) veya `history.jsonl` için `merge by line`. Tanımı farklı olan bir MCP sunucusu da `keep group` veya `use mine` seçenekleriyle bir çakışmadır. Listelenen ilk seçenek, Enter'a bastığınızda varsayılandır; `--yes` her zaman grubun kopyasını korur.

Kontrol listesinin kendisi, `ccam add` komutunun kullandığı `checklist` bileşeni üzerine kurulmuş `tui.MigratePicker` bileşenidir; ccam eksik öğeleri lstat çağrıları ve MCP sunucu listesiyle işaretler (`migrate.Found`).

### MCP oturum açmaları

OAuth ile oturum açan MCP sunucuları tokenlarını sunucu listesinde değil, giriş dosyasında, `.credentials.json` içinde `mcpOAuth` altında tutar. ccam bunları kopyalamaz. Kopyalanmış bir oturum açma, orijinaliyle aynı yenileme tokenını taşır ve sunucu bu tokenı ilk yenilemede değiştirir: kaynak `~/.claude` dahil, onu tutan diğer herkesin o sunucudaki oturumu kapanır. Bu yüzden kaynağın oturum açtığı her kopyalanan sunucu için planın `sign-ins` satırı, kendileri `/mcp` ile oturum açması gereken grup hesaplarını listeler ve özet bunu yineler. Tokenı olmayan bir girdi oturum açma değildir ve listelenmez; zaten kendi oturum açması olan bir hesap ya da tanımını grubun koruduğu bir sunucu da listelenmez. Hiçbir giriş dosyası yazılmaz.

## Güvenlik ve uç durumlar

- Kaynak yalnızca okunur. İçindeki hiçbir dosya taşınmaz, bağlanmaz veya değiştirilmez; kaynak girişi, Anahtar Zinciri'nden okunabildiği macOS dahil, hiçbir zaman yazılmaz.
- Bir grup dosyası `use mine`, bir birleştirme veya mevcut bir `mcp.json` içindeki yeni MCP sunucuları ile değiştirilmeden önce grubun kopyası `<group>/.ccam-backup/<time>/` içine kaydedilir. Klasör adı her çalıştırmada benzersizdir; böylece aynı saniyedeki iki çalıştırma birbirinin yedeklerinin üzerine yazmaz. Özet yolu yazdırır.
- `--yes` ile her çakışmada grubun kopyası korunur ve ccam bunun yerine nasıl seçim yapılacağını söyler: aynı komutu `--yes` olmadan çalıştırın.
- Hiçbir şey değişmeyecekse ccam `Nothing to copy.` yazdırır; her çakışmada grubun kopyası korunuyorsa `Nothing to copy: the group keeps its own copies.` yazdırır.
- Bir kopyalama yarıda başarısız olursa ccam `The source folder was not changed. Group files replaced so far are saved in` ve ardından yedek klasörüyle durur. Yeniden çalıştırmak, tam olarak kopyalanmamış dosyaları kopyalar.
- Birkaç grup MCP paylaştığında `ccam doctor`, `ccam ext` ve `ccam mcp list` komutlarının önerdiği komutlar `ccam migrate --mcp --group <a|b>` biçimindedir.
- Eksik bir kaynak `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` ile biter.

## İlgili

- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Kontrol etme ve onarma](/features/doctor)
- [Tek bilgisayarda birden fazla Claude Code hesabı kullanma](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR nedir?](/guides/claude-config-dir)
