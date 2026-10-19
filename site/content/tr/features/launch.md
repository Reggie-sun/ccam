---
title: Claude Code'u başlatma
description: Claude Code'u herhangi bir hesapla ccam <name|number>, ccam run veya ccam seçicisiyle başlatın; argümanlar claude'a değiştirilmeden iletilir, --dry-run komutu gösterir.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d033da86f44b
---

`ccam work`, Claude Code'u `work` hesabıyla başlatır ve addan sonraki her şey `claude` komutuna olduğu gibi iletilir. `claude` yazacağınız her yerde kullanın: terminalde, bir betikte veya girdiyi boruyla aktararak.

## Kullanım

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Bir hesabın adını veya `ccam list` çıktısındaki numarasını verin. ccam addan sonraki argümanları okumaz; bu yüzden `ccam work --help` Claude Code'un yardımını gösterir. Addan sonra `--` yazmak isteğe bağlıdır. Standart girdi, çıkış kodu ve sinyaller `claude` komutundaki gibi davranır.

İlk argüman `-` ile başlıyorsa ve ccam'in kendi seçeneklerinden biri değilse ccam, varsayılan hesabı tüm argümanlarla başlatır. Varsayılan, `ccam use` ile ayarlanan hesaptır; ayarlamadıysanız en küçük numaralı hesaptır.

### Seçici

Argümansız `ccam`, standart girdi ve çıktı bir terminal olduğunda seçiciyi açar, aksi halde `ccam list` çıktısını yazdırır. Seçici her hesabın numarasını, e-postasını, 5 saatlik kullanımını, son kullanımını ve paylaşılan öğelerini gösterir; kullanım, her hesap yanıt verdikçe dolar.

| Tuş | Ne yapar |
| --- | --- |
| `Enter` | Claude Code'u seçili hesapla başlatır |
| `1` ile `9` arası | O numaralı hesabı başlatır |
| `Up`/`Down`, `k`/`j` | Seçimi taşır |
| `a` | Seçiciyi kapatır ve `ccam add` sihirbazını çalıştırır |
| `s` | Seçili hesabın paylaşılan öğelerini düzenler (`Space` veya `x` değiştirir, `Enter` uygular, `Esc` geri döner) |
| `d` | Seçili hesabı varsayılan yapar |
| `r` | Kısa önbelleği kullanmadan kullanımı yeniden alır |
| `q`, `Esc`, `Ctrl+C` | Çıkar |

Seçici bir işlem yapmadan önce kapanır; böylece terminal Claude Code başlamadan önce eski haline gelir. `d` veya bir paylaşım değişikliğinden sonra yeniden açılır. Paylaşım değişikliği, `ccam share` gibi önce planı gösterir ve onay ister.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--dry-run` | ccam'in çalıştıracağı ortamı ve komutu yazdırır, ardından çıkar | kapalı |
| `-v`, `--version` | `ccam <version>` yazdırır (yalnızca ilk argüman olarak) | |
| `-h`, `--help` | `ccam --help` veya `ccam run --help` olarak ccam'in yardımı | |

Argüman kuralları:

- ccam kendi seçeneklerini yalnızca hesap adından önce okur. Addan sonra `--dry-run` ve `--help` dahil her argüman Claude Code'a aittir.
- `--dry-run`, `ccam --dry-run work`, `ccam run --dry-run work` ve `ccam --dry-run` (varsayılan hesap) olarak çalışır. Reddedilmez, dikkate alınır.
- `ccam run <name|number> [claude args]`, `ccam <name>` komutunun uzun biçimidir. Adsız `ccam run` varsayılan hesabı kullanır.
- Hesap adları `add`, `list`, `run` veya `usage` gibi ccam komut adları olamaz; böylece `ccam <word>` hiçbir zaman belirsiz olmaz.

Kuru çalıştırma, bir kabuğa yapıştırabileceğiniz tek bir satır yazdırır:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Nasıl çalışır

Her başlatma, komut satırından `exec` çağrısına kadar aynı adımları izler:

<!-- motion:launch -->

Her kutu, adımı çalışırken yanar; alttaki satırlar ise Claude Code'un başladığı ortam ve argümanlardır.

1. `cmd/ccam/main.go`, `app.Main` fonksiyonunu çağırır. ccam başka bir adla başlatıldıysa bu bir [kısayoldur](/features/shortcuts) ve doğrudan başlatıcıya gider.
2. `app.parse`, cobra devreye girmeden önce `os.Args` değerine bakar. Yalnızca tek başına `ccam` (seçici), bir ccam komut adı (`config.Reserved`), `-h`/`--help` ve `-v`/`--version`, `cli.Execute` fonksiyonuna gider. Geri kalan her şey bir başlatmadır.
3. `launcher.Launch`, ayrıştırdığı tek dosya olan `~/.config/ccam/config.toml` dosyasını okur.
4. `launcher.Prepare` hesabı ada veya numaraya göre bulur, baştaki bir `--` değerini atar ve `claude` yürütülebilir dosyasını `ResolveClaude` ile bulur. PATH araması, dosyanın değiştirilme zamanıyla birlikte `~/.local/state/ccam/claude-path` içinde önbelleğe alınır; böylece ccam her başlatmada `PATH` içinde arama yapmaz. Değişmiş veya eksik bir dosya yeni bir aramayı tetikler. `config.toml` içindeki bir `claude_bin` yolu olduğu gibi kullanılır. ccam, inode ile karşılaştırarak ccam'in kendisi olan bir `claude` dosyasını reddeder, çünkü bu bir döngüye yol açar.
5. `checkLinks`, hesap klasörü `~/.local/share/ccam/profiles/<name>/` ile paylaşım ayarlarını `share.Check` ile karşılaştırır; bu fonksiyon yalnızca `lstat` ve `readlink` çağırır, öğe başına bir veya iki sistem çağrısı. Eksik bir bağlantı anında oluşturulur. Veri kaybettirebilecek her şey başlatmayı durdurur.
6. `applyMCP`, paylaşılan MCP sunucularını ekler (aşağıya bakın).
7. `session.Register` oturumu kaydeder. ccam, PID'sinin adını taşıyan ve süreç başlangıç zamanını içeren bir işaret dosyasını `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` konumuna (`XDG_RUNTIME_DIR` yoksa `~/.local/state/ccam/locks/` içine) yazar. Konuşma geçmişini paylaşan bir hesap için `historyCheck`, önce `history-<group>-<hash>` içinde çalışan bir oturum arar (hash, geçerli klasörün hash'idir) ve oraya da kaydolur. `exec` PID'yi korur; böylece işaret dosyası artık Claude Code sürecinin kendisini adlandırır. Tutulan bir kilit değil, bir işaret dosyası kullanılır; çünkü Claude Code'un başlattığı süreçler bir kilidi devralır ve oturum bittikten sonra da tutar.
8. `Plan.Exec`, terminal başlığını `claude · <name>` olarak ayarlar (standart hata bir terminal olduğunda), başlatma zamanını `fsync` olmadan `~/.local/state/ccam/last-used/<name>` dosyasına yazar ve `syscall.Exec` çağırır.

Claude Code, `CLAUDE_CONFIG_DIR` hesap klasörüne ve `CCAM_PROFILE` hesap adına ayarlanmış olarak çalışır. ccam, bunların devraldığı değerlerini kaldırır. macOS'ta `ccam import` ile kaydedilmiş bir `~/.claude` klasörü için ayrıca `CLAUDE_SECURESTORAGE_CONFIG_DIR` değişkenini boş bir değere ayarlar; böylece hesap ve düz `claude` aynı Anahtar Zinciri girişini kullanır.

### Neden exec

`syscall.Exec`, ccam sürecini `claude` ile değiştirir; böylece ccam bellekte kalmaz, tuş vuruşları, sinyaller ve çıkış kodu doğrudan Claude Code'a aittir ve argümanlar hiçbir zaman bir kabuk dizesi olmadan dizi olarak kalır. Bir oturumdan sonraki işler, örneğin `claude mcp add` ile eklenen sunucuları fark etmek, bir sonraki başlatmada yapılır: ccam, `.claude.json` dosyasının değiştirilme zamanını `~/.local/state/ccam/mcp/<name>/seen` ile karşılaştırır ve dosyayı yalnızca değiştiyse ayrıştırır.

### Bir komut satırı nasıl yönlendirilir

Üç komut satırı, `internal/app` içinden üç yol:

<!-- motion:launch-route -->

Bir kısayol `argv[0]` ile, `ccam <name>` ise ilk argümanla eşleşir ve ikisi de `syscall.Exec` ile biter. Yalnızca `ccam usage` gibi bir ccam komutu cobra'ya ulaşır.

### Paylaşılan MCP sunucuları

Varsayılan `flag` stratejisiyle ccam, grubun MCP sunucularını Claude Code'a sizinkilerden önce yerleştirilen bir argüman olarak iletir:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` yalnızca grubun `mcp.json` dosyası varsa bir argüman döndürür ve `Plan.InsertArgs` bunu `claude` komutunun hemen arkasına koyar. Kendi `--mcp-config` veya `--strict-mcp-config` seçeneğiniz bunun üzerine yine uygulanır. Bu tek bir argümandır, `--mcp-config=<path>`, çünkü Claude Code'un `--mcp-config` seçeneği birden fazla değer alır: iki argüman olarak yazılsaydı ardından gelen bir istemi de yapılandırma dosyası olarak okurdu.

`config.toml` içinde `mcp_strategy = "sync"` ile ccam bunun yerine sunucuları hesabın `.claude.json` dosyasına yazar ve hesap zaten çalışıyorsa bunu atlar.

### 10 ms bütçesi

ccam, Claude Code başlamadan önce en fazla 10 ms (p95) ekler. Başlatma yolu hiçbir zaman cobra komut ağacını oluşturmaz, Bubble Tea başlatmaz veya bir ağ bağlantısı açmaz; paket `init()` fonksiyonları toplamda 5 ms'den uzun sürerse `TestPackageInitIsCheap` başarısız olur. `scripts/perf-check.sh`, `ccam --dry-run work`, `ccam work` ve bir kısayolu `hyperfine` ile ölçer ve 10 ms'nin üzerinde başarısız olur; `internal/launcher` içindeki `BenchmarkLaunchPath` ise adımları bir döngüde çalıştırır. Adım başına süreyi görmek için `CCAM_TRACE=1` ayarlayın:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Güvenlik ve uç durumlar

- Bilinmeyen bir ad veya numara `no account "<name>"` ile durur ve `ccam list` komutunu gösterir. Hiç hesap yoksa ccam, `ccam add <name>` ile bir hesap eklemenizi ister.
- Paylaşılan bir yol sıradan bir dosyaya dönüştüyse veya yanlış yere yöneliyorsa ccam `stopped because of share link problems:` ile durur, sorunları listeler ve `ccam doctor --fix` çalıştırmanızı ister. Başlatma yolunda veri taşımaz.
- Bu hesapla konuşma geçmişini paylaşan (aynı gruptaki) başka bir Claude Code aynı klasörde çalışıyorsa ccam, geçmişin karışabileceği konusunda uyarır ve `Run anyway? [y/N]` diye sorar. Terminal yoksa uyarıyı gösterir ve devam eder.
- `--dry-run`; MCP eşitlemesini, yeni sunucu bildirimini, geçmiş kontrolünü, oturum işaretini, başlığı ve son kullanım zamanını atlar ve exec yapmaz. Paylaşım bağlantılarını yine kontrol eder ve eksik bir bağlantıyı normal bir başlatmanın yapacağı gibi oluşturur.
- `mcp_strategy = "sync"` ile başarısız bir eşitleme `ccam: MCP sync failed (launching anyway): …` yazdırır ve Claude Code zaten sahip olduğu sunucularla başlar.

## İlgili

- [Kısayol komutları](/features/shortcuts)
- [Hesapları yönetme](/features/accounts)
- [Komut başvurusu](/docs/commands)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [CLAUDE_CONFIG_DIR nasıl çalışır](/guides/claude-config-dir)
