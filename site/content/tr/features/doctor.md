---
title: Kontrol etme ve onarma
description: ccam doctor ile paylaşım bağlantılarını, dosya izinlerini, girişleri ve Claude Code sürümünü kontrol edin, güvenli olanı ccam doctor --fix ile onarın.
updated: 2026-10-23
order: 9
translation: machine
sourceHash: "11ad34630dcf"
---

`ccam doctor` ccam kurulumunuzu kontrol eder: Claude Code ikili dosyası, her hesabın klasörü, girişi ve paylaşım bağlantıları, paylaşılan MCP sunucuları, kullanım uç noktası ve kısayol komutlarınız. `ccam doctor --fix`, veri kaybetmeden yapabileceklerini onarır ve değiştirdiği her dosyayı önce bir yedek klasörüne taşır. `ccam <name>` bir paylaşım bağlantısı sorunuyla durduğunda, dosyaları elle geri yükledikten sonra veya bir Claude Code güncellemesinden sonra çalıştırın.

## Kullanım

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

İlk satır tüm hesapları bir kerede işaretler; böylece başka bir şey okumadan önce sorunların nerede olduğunu görürsünüz. Altında her hesap bir kutu alır, hiçbir hesaba ait olmayan bulgular (Claude Code, ortam değişkenleri, eklentiler, MCP sunucuları, kısayollar, klasör kuralları, gruplar ve ayarlar) ise bir `General` kutusuna girer. Kutu başlığı, hesabın en kötü işaretini ve kaç sorunu olduğunu gösterir. İçerideki her satır kontrol edilen bir alandır; her zaman bu sırayla `Folder`, `Login` ve `Share links`, ardından diğerleri gelir ve bir işaret taşır: `✔` sorun yok, `ℹ` bilgi, `⚠` uyarı, `✘` hata. Hiçbir sorunu olmayan bir klasör `✔ Folder  OK` gösterir.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

`--fix` ile onarılan her bulgu `→ fixed` ile biter ve artık sorun olarak sayılmaz:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

Terminalde kutular terminalin genişliğine sığar, uzun mesajlar sarılır ve terminalin kendi renkleri kullanılır (`NO_COLOR` renkleri kapatır). Bir kutu hiçbir zaman 29 sütundan dar olmaz; böylece dar bir terminalde hiçbir sözcük kesilmez. Bir boruya veya dosyaya gönderildiğinde çıktıda renk olmaz ve hiçbir mesaj sarılmaz veya kesilmez; böylece `ccam doctor | grep settings.json` satırın tamamını bulur. Bu durumda her kutu en uzun satır kadar geniştir. Kontrol panelinin Health sekmesi aynı kutuları gösterir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--fix` | Onarılması güvenli olanı onarır (değişen dosyalar yedeklenir) | kapalı |
| `--dry-run` (genel) | Tüm kontrolleri çalıştırır ama `--fix` ile bile hiçbir şeyi onarmaz | kapalı |

`ccam doctor` argüman almaz. `--dry-run` seçeneğini dikkate alır: `ccam --dry-run doctor --fix`, `ccam doctor` bulgularını yazdırır, `--fix` seçeneğinin kenara ayıracaklarını ekler ve hiçbir şeyi değiştirmez:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

En az bir `✘` bulgusu onarılmadan kaldığında çıkış kodu `1`, aksi halde `0` olur. Uyarılar ve bilgiler çalıştırmayı hiçbir zaman başarısız kılmaz. Çalıştırma `--fix` olmadan başarısız olduğunda ccam boş bir satır ve `` `ccam doctor --fix` repairs what it can. `` ile biter. Yüklenemeyen bir yapılandırma da `1` ve bir `ccam:` hatasıyla biter.

## Nasıl çalışır

`internal/cli/doctor.go`, `~/.config/ccam/config.toml` dosyasını yükler ve `internal/doctor/doctor.go` içindeki `doctor.Run` fonksiyonunu çağırır; bu fonksiyon kontrolleri sabit bir sırayla çalıştırır ve her biri ait olduğu hesap ve kontrol ettiği alanla birlikte bulguları döndürür. `tui.HealthReport` bunları kutulara gruplar. `--fix`, `Options.Fix` olarak iletilir ve kuru çalıştırmada false olarak ayarlanır.

<!-- motion:doctor-checks -->

Ortadaki blok her hesap için bir kez çalışır. Kontroller ve `--fix` seçeneğinin her biri için yaptığı:

1. **Claude Code sürümü** (`checkClaude`). `claude` ikili dosyasını başlatıcının yaptığı gibi bulur ve `claude --version` komutunu çalıştırır. Eksik bir ikili dosya, çalıştırılabilir bir dosya olmayan bir `claude_bin` veya başarısız olan bir `--version` hatadır. Onarım yok.
2. **Ortamdaki kimlik bilgileri** (`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` veya `ANTHROPIC_AUTH_TOKEN` ayarlı olduğunda uyarır: o kabuktan başlatılan her hesap bunu devralır ve Claude Code bunu hesabın kendi girişine tercih eder. Onarım yok; değişkeni kendiniz kaldırın.
3. **Klasör** (`checkProfile`). Eksik bir hesap klasörü hatadır ve o hesabın kontrollerini sonlandırır. Grup veya diğerleri için herhangi bir bit ayarlanmış izinler uyarıdır; `--fix`, `chmod 0700` çalıştırır. `~/.claude.json` dosyasının yanında kendine ait bir `.claude.json` dosyası bulunan içe aktarılmış bir `~/.claude` bilgidir: `ccam main` ve düz `claude` bu durumda giriş ayrıntılarını ve MCP sunucularını ayrı tutar. Windows'ta izin bitleri yoktur; bu yüzden orada bu kontrol ve aşağıdaki `0600` kontrolü atlanır, kullanıcı profilinizdeki klasörü erişim listesi korur.
4. **Giriş dosyası** (`checkCredentials`). `.credentials.json`, `lstat` ile kontrol edilir. Sembolik bağlantı bir hatadır, çünkü her hesabın kendi tokenı olmalıdır, ve onarılmaz. `0600` dışındaki bir mod uyarıdır; `--fix`, `chmod 0600` çalıştırır. Ardından giriş okunur (macOS'ta önce Anahtar Zinciri'nden). Giriş yoksa ya da giriş dosyasında giriş bulunmuyorsa `ccam login <name>` komutunu gösteren bir uyarı verilir. Süresi dolmuş bir yenileme tokenı ya da sunucu reddettikten sonra Claude Code'un sildiği bir giriş `` login expired → `ccam login <name>` `` hatasıdır; süresi dolmuş bir erişim tokenı ise yalnızca bilgidir, çünkü Claude Code onu bir sonraki başlatmada yeniler. Windows Kimlik Bilgileri Yöneticisi'nde tutulan bir giriş bilgidir, çünkü ccam onu okuyamaz. Geçerli bir giriş, klasörün şu anda giriş yaptığı e-postayı ve planı gösterir (`login valid (dev@company.com, Claude Team · Acme)`). Bunların hiçbiri onarılmaz; `ccam login --expired` gereken her hesapta giriş yapar.
5. **Paylaşım bağlantıları** (`share.Check`). `share.Layout` içindeki her yolu, yalnızca `lstat` ve `readlink` kullanarak hesabın gruplarıyla karşılaştırır. Beş tür sorun bildirir: eksik bağlantı, başka bir yere yönelen bağlantı, bağlantı olması gereken yerde gerçek bir dosya, artık yerel olan bir öğe için geride kalmış bir bağlantı ve grup klasörü eksik olan doğru bir bağlantı (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; bir başlatma bunu denetlemez). Eksik bağlantı ve eksik grup klasörü uyarıdır, diğerleri hatadır. Gruptan geri kopyalamanın yarıda kalması da bildirilir: `--fix` hazırlanmış eksiksiz bir kopyayı yerine koyar, ancak yol o zamandan beri (bir başlatma tarafından) yeniden oluşturulmuşsa, sizin düzenlemeniz için eksiksiz kopyanın bulunduğu klasörü adlandırır. ccam'in sonraki bir sürümünde bir öğeye eklenen yollar (geçmiş için `file-history/`, `tasks/` ve `paste-cache/`, ayarlar için `keybindings.json`), hesabın hâlâ kendi kopyası olduğu sürece yalnızca uyarıdır: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Bir başlatma bunlar yüzünden asla durmaz.
6. **Bilinmeyen girdiler** (`checkUnknownEntries`). Hesap klasöründe ccam'in bilmediği (`share.Known`) üst düzey adları, Claude Code düzeninin değiştiğine dair bir ipucu olarak listeler. Yalnızca bilgi.
7. **Eklenti yolları** (`checkPlugins`). Eklenti kayıtları (`installed_plugins.json`, `known_marketplaces.json`) mutlak yollar saklar. Her grubun `plugins/` klasöründe, bir hesabın klasöründen geçen yollar o hesap kaldırıldığında veya yeniden adlandırıldığında bozulur. `--fix`, bunları kayıt dosyasını atomik olarak yazan `profile.RewritePluginPaths` ile grup klasörüne göre yeniden yazar.
8. **MCP sunucuları** (`checkMCP`). Her grubun `mcp.json` dosyasını bir kez yükler. Okunamayan bir dosya hatadır; `mcp.Validate` kontrolünden geçemeyen bir sunucu uyarıdır. Düz `claude` komutunun `~/.claude.json` dosyasında bulunup hiçbir hesabın yüklemediği sunucular, onları kopyalayan `ccam migrate --mcp` komutunu (birkaç grup MCP paylaşıyorsa `--group <a|b>` ile) belirten bir uyarıdır. Bir grubun yalnızca tek bir hesabında bulunan sunucular, onları paylaşan `ccam mcp import` komutuyla birlikte bilgidir; diğer hesapların giriş yaptığı MCP sunucularına giriş yapmamış bir hesap ise `/mcp` ile giriş yapması için bir uyarı alır (ccam girişleri asla kopyalamaz). Onarım yok.
9. **Kullanım yanıtı biçimi** (`checkUsageFormat`). Kullanım uç noktasına önbellek olmadan, biri yanıt verene kadar hesapları sırayla, hepsini tek bir 4 saniyelik sınır içinde sorar. `usage_source = "claude"` ile bunun yerine tek bir hesap için Claude Code'a sorar. Yanıt artık ayrıştırılamıyorsa `ccam usage` komutunun durum satırı önbelleğine geri döndüğü konusunda uyarır; hiçbir hesap yanıt vermediyse bunu bilgi olarak söyler (`could not check the usage response format: …`). Ağı kullanan tek kontrol budur.
10. **Kısayollar** (`checkShims`). Yalnızca bir hesabın kısayolu olduğunda. `~/.local/bin` içindeki eksik bağlantılar uyarıdır; ccam'in oluşturduğu ve hâlâ daha eski, değiştirilmiş veya taşınmış bir ccam'i çalıştıran kısayollar da öyledir (`shortcuts in <bin> that run another ccam (one replaced or moved): …`). `--fix`, yalnızca ccam'in oluşturduğu ve `~/.local/state/ccam/shortcuts.json` dosyasına kaydettiği bağlantıları oluşturan ve yenileyen, diğer dosyalara hiç dokunmayan `shim.Install` fonksiyonunu çağırır. Uyarı yalnızca bu bağlantıların her biri oluşturulduğunda `→ fixed` ile biter; biri oluşturulamadıysa nedeni eklenir. Hâlâ kurulu olan başka bir ccam'i çalıştıran bir kısayol çalışır ve bilgidir; `--fix` onu bu ccam'e yönlendirir. `PATH` üzerinde olmayan bir bin klasörü uyarıdır ve onarılmaz.
11. **Klasörler** (`checkFolders`). Yalnızca [`ccam folder`](/features/folder) girdileri varsa. Hesabı bulunmayan bir girdi hatadır, klasörü artık olmayan bir girdi ise uyarıdır. Hiçbiri düzeltilmez, çünkü bir girdiyi kaldırmak orada hangi hesabın başlayacağını değiştirir.
12. **Gruplar ve ayarlar** (`checkGroups`). Hesap olmayan bir grup üyesi, adı sonradan bir ccam komutu hâline gelmiş bir hesap (onu `ccam run <name>` ile başlatın ya da yeniden adlandırın), bilinmeyen bir `theme` veya `usage_source` ve `config.toml` içinde ccam'in okumadığı satırlar uyarıdır. Sonuncular genellikle bir `[[group]]`, `[[profile]]` veya `[[folder]]` tablosunun altına yazılmış ve bu yüzden o tabloya ait olan bir ayardır (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); yapılandırmayı değiştiren komutlar bu düzeltilene kadar çalışmayı reddeder. Onarılmaz.
13. **Girişler** (`checkIdentities`). Her klasörün kimle giriş yaptığını kendi `.claude.json` dosyasından okur. Artık `ccam login` ile kaydedilenden başka bir hesaba giriş yapmış bir klasör (örneğin bir oturumun içinde `/login` sonrasında) uyarıdır; içe aktarılmış ya da daha eski bir ccam tarafından oluşturulmuş gibi giriş kaydı olmayan bir hesap, kayıtlı e-postasıyla karşılaştırılır. Başka bir klasörle aynı Claude hesabına giriş yapmış bir klasör de uyarıdır: ikisi kullanım limitlerini paylaşır. Başka bir klasörle aynı giriş tokenını tutan bir klasör, yani kopyalanmış bir klasör (düz `claude` komutunun `~/.claude` klasörü dahil), hatadır, çünkü onu ilk yenileyen diğerinin oturumunu kapatır; ona kendi girişini vermek için `ccam login <name>` çalıştırın. Onarılmaz.

Hesaplardan önce, yarıda kalmış bir **grup değişikliği** `General` kutusunda bir hatadır: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` onu planlandığı sırayla, çakışmalarına verilen yanıtlarla tamamlar; o zamana kadar başka hiçbir grup değişikliği, yeniden adlandırma veya kaldırma başlamaz. İçindeki hesaplar, tamamlanana kadar yalnızca eksik bağlantılarını alır. Tamamlanamazsa hata, o dosyayı kaldırmanın değişikliğin geri kalanından vazgeçmek anlamına geldiğini söyler. `ccam --dry-run doctor --fix` sırayı (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) ve her hesabın payının neyi kenara ayıracağını gösterir. Başka bir ccam'in şu anda yürüttüğü bir değişiklik bilgidir ve o çalışırken `--fix` hiçbir paylaşım bağlantısını onarmaz (`its share links are not repaired while another ccam changes the groups; …`). `--fix` onarım yaparken `~/.local/state/ccam/group-change.lock` kilidini tutar ve bir değişikliği tamamladıktan sonra `config.toml` dosyasını yeniden okur.

### Paylaşım bağlantılarını onarma

Başlatma yolu eksik bir bağlantıyı anında oluşturur, ancak veri kaybettirebilecek her durumda durur ve buraya yönlendirir:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Herhangi bir paylaşım bağlantısı sorununda `--fix`, hesabın paylaşım ayarlarını hem mevcut hem de hedef durum olarak vererek `share.Build` fonksiyonunu, ardından `Plan.Apply` fonksiyonunu çağırır. `Build`, her adımı diskte bulunandan yeniden türetir; bu, bir `ccam group` değişikliğinin kullandığı planın aynısıdır. Böylece onarım, bir gruba katılmakla aynı kurallara uyar.

<!-- motion:doctor-repair -->

Burada `settings.json` gerçek bir dosyadır. Grupta zaten farklı bir `settings.json` bulunduğu için çakışma varsayılan seçimi olan `KeepGroup` ile çözülür ve grubun kopyası kalır. Hesabın dosyası, hesap klasörünün içindeki `.ccam-backup/<time>/settings.json` konumuna taşınır ve `os.Symlink` yolu gruba bağlar. `--fix` bunu söyler: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Diğer durumlar:

- Grupta henüz kopya yoksa hesabın dosyası gruba taşınır ve grubun ilk içeriği olur.
- Bir klasörde, grupta olmayan dosyalar gruba taşınır. Grupta bayt bayt aynısı bulunan dosyalar bırakılır, farklı olanlar ise bir notta adlarıyla belirtilerek yedeğe gider (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Eksik bir grup klasörü önce oluşturulur.
- `history.jsonl` için varsayılan `merge by line` seçeneğidir: önce grubun satırları, ardından sizin yeni satırlarınız. Grubun önceki kopyası `.ccam-backup/<time>/group/` içine kaydedilir.
- Başka bir gruba yönelen bağlantı değiştirilir. ccam dışındaki bir klasöre yönelen bağlantı yedeğe gider ve grup boşsa içeriği grubun ilk içeriği olur.
- Artık yerel olan bir öğe için bırakılmış bağlantı kaldırılır ve hesap, grubun içeriğinin bir kopyasını alır.

Bağlantıların kaldırılması ve grupta bayt bayt aynısı bulunan dosyaların bırakılması dışında her adım veriyi gruba veya yedek klasörüne taşır. Hiçbir şey kaybolmaz.

## Güvenlik ve uç durumlar

- Düz bir `ccam doctor` hiçbir zaman hiçbir şeyi değiştirmez. Betiklerde emin olmak için `--fix` ile birlikte `--dry-run` kullanın.
- Claude Code nasıl başlatılmış olursa olsun, çalışan bir hesap için `doctor --fix` yalnızca eksik bağlantıları ve eksik grup klasörlerini oluşturur; bu hiçbir şey kaybettirmez. Geriye başka bir şey kalırsa `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest `` ekler, çünkü oturum, bir onarımın taşıyacağı dosyalara yazmaya devam eder.
- Bir paylaşım onarımı başarısız olursa ccam o hesap için `could not repair share links: <error>` bildirir ve bilinmeyen girdi kontrolünü atlar. Diğer hesaplar yine kontrol edilir.
- Bir `ccam group` değişikliği yarıda başarısız olursa yeni gruplar zaten kaydedilmiştir ve değişiklik `~/.local/state/ccam/group-change.json` dosyasına not edilir; ccam işlemi, MCP adımları dahil geri kalanını yürüten `ccam doctor --fix` ile tamamlamanızı söyler.
- Sembolik bağlantılı bir `.credentials.json`, süresi dolmuş bir giriş ve ortamdaki kimlik bilgileri bildirilir ama hiçbir zaman değiştirilmez.
- Eklenti kaydı yeniden yazımları atomik olarak yazılır, ancak `.ccam-backup/` içine kopyalanmaz.

## İlgili

- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Mevcut bir kurulumu kopyalama](/features/migrate)
- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [Komutlar](/docs/commands)
