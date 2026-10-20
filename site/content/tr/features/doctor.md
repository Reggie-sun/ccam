---
title: Kontrol etme ve onarma
description: ccam doctor ile paylaşım bağlantılarını, dosya izinlerini, girişleri ve Claude Code sürümünü kontrol edin, güvenli olanı ccam doctor --fix ile onarın.
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` ccam kurulumunuzu kontrol eder: Claude Code ikili dosyası, her hesabın klasörü, girişi ve paylaşım bağlantıları, paylaşılan MCP sunucuları, kullanım uç noktası ve kısayol komutlarınız. `ccam doctor --fix`, veri kaybetmeden yapabileceklerini onarır ve değiştirdiği her dosyayı önce bir yedek klasörüne taşır. `ccam <name>` bir paylaşım bağlantısı sorunuyla durduğunda, dosyaları elle geri yükledikten sonra veya bir Claude Code güncellemesinden sonra çalıştırın.

## Kullanım

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Her bulgu, işaretli bir satırdır: `✓` sorun yok, `·` bilgi, `!` uyarı, `✗` hata. Hesap bulguları hesap adıyla başlar.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

`--fix` ile onarılan her bulgu `→ fixed` ile biter:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--fix` | Onarılması güvenli olanı onarır (değişen dosyalar yedeklenir) | kapalı |
| `--dry-run` (genel) | Tüm kontrolleri çalıştırır ama `--fix` ile bile hiçbir şeyi onarmaz | kapalı |

`ccam doctor` argüman almaz. `--dry-run` seçeneğini dikkate alır: `ccam --dry-run doctor --fix`, `ccam doctor` ile aynı bulguları yazdırır ve hiçbir şeyi değiştirmez.

En az bir `✗` bulgusu onarılmadan kaldığında çıkış kodu `1`, aksi halde `0` olur. Uyarılar ve bilgiler çalıştırmayı hiçbir zaman başarısız kılmaz. Çalıştırma `--fix` olmadan başarısız olduğunda ccam boş bir satır ve `` `ccam doctor --fix` repairs what it can. `` ile biter. Yüklenemeyen bir yapılandırma da `1` ve bir `ccam:` hatasıyla biter.

## Nasıl çalışır

`internal/cli/doctor.go`, `~/.config/ccam/config.toml` dosyasını yükler ve `internal/doctor/doctor.go` içindeki `doctor.Run` fonksiyonunu çağırır; bu fonksiyon kontrolleri sabit bir sırayla çalıştırır ve bulguları döndürür. `--fix`, `Options.Fix` olarak iletilir ve kuru çalıştırmada false olarak ayarlanır.

<!-- motion:doctor-checks -->

Ortadaki blok her hesap için bir kez çalışır. Kontroller ve `--fix` seçeneğinin her biri için yaptığı:

1. **Claude Code sürümü** (`checkClaude`). `claude` ikili dosyasını başlatıcının yaptığı gibi bulur ve `claude --version` komutunu çalıştırır. Sürüm, ccam'in klasör düzeninin doğrulandığı sürüm olan `2.1.302` (`doctor.LayoutVersion`) ile başlamıyorsa `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` satırını ekler. Eksik bir ikili dosya veya başarısız olan bir `--version` hatadır. Onarım yok.
2. **Ortamdaki kimlik bilgileri** (`checkAuthEnv`). `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` veya `ANTHROPIC_AUTH_TOKEN` ayarlı olduğunda uyarır: o kabuktan başlatılan her hesap bunu devralır ve Claude Code bunu hesabın kendi girişine tercih eder. Onarım yok; değişkeni kendiniz kaldırın.
3. **Klasör** (`checkProfile`). Eksik bir hesap klasörü hatadır ve o hesabın kontrollerini sonlandırır. Grup veya diğerleri için herhangi bir bit ayarlanmış izinler uyarıdır; `--fix`, `chmod 0700` çalıştırır.
4. **Giriş dosyası** (`checkCredentials`). `.credentials.json`, `lstat` ile kontrol edilir. Sembolik bağlantı bir hatadır, çünkü her hesabın kendi tokenı olmalıdır, ve onarılmaz. `0600` dışındaki bir mod uyarıdır; `--fix`, `chmod 0600` çalıştırır. Ardından giriş okunur (macOS'ta önce Anahtar Zinciri'nden). Giriş yoksa `ccam login <name>` komutunu gösteren bir uyarı verilir. Süresi dolmuş bir yenileme tokenı `login expired` hatasıdır; süresi dolmuş bir erişim tokenı ise yalnızca bilgidir, çünkü Claude Code onu bir sonraki başlatmada yeniler. Bunların hiçbiri onarılmaz.
5. **Paylaşım bağlantıları** (`share.Check`). `share.Layout` içindeki her yolu, yalnızca `lstat` ve `readlink` kullanarak hesabın paylaşım ayarlarıyla karşılaştırır. Dört tür sorun bildirir: eksik bağlantı, başka bir yere yönelen bağlantı, bağlantı olması gereken yerde gerçek bir dosya ve artık yerel olan bir öğe için geride kalmış bir bağlantı. Eksik bağlantı uyarıdır, diğerleri hatadır.
6. **Bilinmeyen girdiler** (`checkUnknownEntries`). Hesap klasöründe ccam'in bilmediği (`share.Known`) üst düzey adları, Claude Code düzeninin değiştiğine dair bir ipucu olarak listeler. Yalnızca bilgi.
7. **Eklenti yolları** (`checkPlugins`). Eklenti kayıtları (`installed_plugins.json`, `known_marketplaces.json`) mutlak yollar saklar. Her grubun `plugins/` klasöründe, bir hesabın klasöründen geçen yollar o hesap kaldırıldığında veya yeniden adlandırıldığında bozulur. `--fix`, bunları kayıt dosyasını atomik olarak yazan `profile.RewritePluginPaths` ile grup klasörüne göre yeniden yazar.
8. **MCP sunucuları** (`checkMCP`). Her grubun `mcp.json` dosyasını bir kez yükler. Okunamayan bir dosya hatadır; `mcp.Validate` kontrolünden geçemeyen bir sunucu uyarıdır. Onarım yok.
9. **Kullanım yanıtı biçimi** (`checkUsageFormat`). Kullanım uç noktasına önbellek olmadan ve olağan 3 saniyelik sınırla, biri yanıt verene kadar hesapları sırayla sorar. Yanıt artık ayrıştırılamıyorsa `ccam usage` komutunun durum satırı önbelleğine geri döndüğü konusunda uyarır. Ağı kullanan tek kontrol budur.
10. **Kısayollar** (`checkShims`). Yalnızca bir hesabın kısayolu olduğunda. `~/.local/bin` içindeki eksik bağlantılar uyarıdır; `--fix`, yalnızca eksik bağlantıları oluşturan ve diğer dosyalara hiç dokunmayan `shim.Install` fonksiyonunu çağırır. `PATH` üzerinde olmayan bir bin klasörü uyarıdır ve onarılmaz.

### Paylaşım bağlantılarını onarma

Başlatma yolu eksik bir bağlantıyı anında oluşturur, ancak veri kaybettirebilecek her durumda durur ve buraya yönlendirir:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Herhangi bir paylaşım bağlantısı sorununda `--fix`, hesabın paylaşım ayarlarını hem mevcut hem de hedef durum olarak vererek `share.Build` fonksiyonunu, ardından `Plan.Apply` fonksiyonunu çağırır. `Build`, her adımı diskte bulunandan yeniden türetir; bu, `ccam share` komutunun kullandığı planın aynısıdır. Böylece onarım, paylaşımı açmakla aynı kurallara uyar.

<!-- motion:doctor-repair -->

Burada `settings.json` gerçek bir dosyadır. Grupta zaten farklı bir `settings.json` bulunduğu için çakışma varsayılan seçimi olan `KeepGroup` ile çözülür ve grubun kopyası kalır. Hesabın dosyası, hesap klasörünün içindeki `.ccam-backup/<time>/settings.json` konumuna taşınır ve `os.Symlink` yolu gruba bağlar. Diğer durumlar:

- Grupta henüz kopya yoksa hesabın dosyası gruba taşınır ve grubun ilk içeriği olur.
- Bir klasörde, grupta olmayan dosyalar gruba taşınır; geri kalanı yedeğe gider.
- `history.jsonl` için varsayılan `merge by line` seçeneğidir: önce grubun satırları, ardından sizin yeni satırlarınız. Grubun önceki kopyası `.ccam-backup/<time>/group/` içine kaydedilir.
- Başka bir gruba yönelen bağlantı değiştirilir. ccam dışındaki bir klasöre yönelen bağlantı yedeğe gider ve grup boşsa içeriği grubun ilk içeriği olur.
- Artık yerel olan bir öğe için bırakılmış bağlantı kaldırılır ve hesap, grubun içeriğinin bir kopyasını alır.

Bağlantıların kaldırılması dışında her adım veriyi gruba veya yedek klasörüne taşır. Hiçbir şey silinmez.

## Güvenlik ve uç durumlar

- Düz bir `ccam doctor` hiçbir zaman hiçbir şeyi değiştirmez. Betiklerde emin olmak için `--fix` ile birlikte `--dry-run` kullanın.
- `doctor --fix`, bir hesabın çalışıp çalışmadığını kontrol etmez. Bir hesabın paylaşım bağlantılarını onarmadan önce o hesaptaki Claude Code'dan çıkın, çünkü oturum bu dosyalara yazmaya devam eder.
- Bir paylaşım onarımı başarısız olursa ccam o hesap için `could not repair share links: <error>` bildirir ve bilinmeyen girdi kontrolünü atlar. Diğer hesaplar yine kontrol edilir.
- `ccam share` yarıda başarısız olursa önce hedef ayarları kaydeder ve işlemi `ccam doctor --fix` ile tamamlamanızı söyler.
- Sembolik bağlantılı bir `.credentials.json`, süresi dolmuş bir giriş ve ortamdaki kimlik bilgileri bildirilir ama hiçbir zaman değiştirilmez.
- Eklenti kaydı yeniden yazımları atomik olarak yazılır, ancak `.ccam-backup/` içine kopyalanmaz.

## İlgili

- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Mevcut bir kurulumu kopyalama](/features/migrate)
- [Kullanım limitleri ve hesap değiştirme](/docs/usage-limits)
- [Komutlar](/docs/commands)
