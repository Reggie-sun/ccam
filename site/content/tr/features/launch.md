---
title: Claude Code'u başlatma
description: Claude Code'u herhangi bir hesapla ccam <name|number>, ccam run veya kontrol paneliyle başlatın; argümanlar claude'a aynen gider, --dry-run komutu gösterir.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: "4812b197f143"
---

`ccam work`, Claude Code'u `work` hesabıyla başlatır ve addan sonraki her şey `claude` komutuna olduğu gibi iletilir. `claude` yazacağınız her yerde kullanın: terminalde, bir betikte veya girdiyi boruyla aktararak.

## Kullanım

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Bir hesabın adını veya `ccam list` çıktısındaki numarasını verin. ccam addan sonraki argümanları okumaz; bu yüzden `ccam work --help` Claude Code'un yardımını gösterir. Addan sonra `--` yazmak isteğe bağlıdır. Standart girdi, çıkış kodu ve sinyaller `claude` komutundaki gibi davranır.

İlk argüman `-` ile başlıyorsa ve ccam'in kendi seçeneklerinden biri değilse ccam, adsız bir hesabı tüm argümanlarla başlatır. Bu, geçerli klasör için [`ccam folder`](/features/folder) ile ayarlanan hesaptır; yoksa varsayılan hesaptır: `ccam use` ile ayarlanan hesap ya da ayarlamadıysanız en küçük numaralı hesap.

### Kontrol paneli

Argümansız `ccam`, standart girdi ve çıktı bir terminal olduğunda [kontrol panelini](/features/dashboard) açar, aksi halde `ccam list` çıktısını yazdırır. Bir hesaba gidin (ok tuşlarıyla ya da numarasını yazarak) ve başlatmak için `Enter` tuşuna basın; `n` en çok payı kalan hesabı başlatır, `o` ise `--continue`, `--resume` veya başka argümanlar sunar. Kontrol paneli Claude Code başlamadan önce kapanır; böylece terminal önce eski haline gelir.

Kabukta `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` veya `ANTHROPIC_AUTH_TOKEN` ayarlıysa Claude Code, hesabın kendi girişi yerine o anahtarı ya da token'ı kullanabilir; bu yüzden ccam başlamadan önce bunu standart hataya yazar (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Değişkenleri yalnızca okur; bir [API anahtarlı hesap](/features/api) ise onları temizler.

`config.toml` içinde `relay = true` olduğunda ccam etkileşimli oturumları bir alt süreç olarak başlatır ve yanlarında kalır; böylece `/ccam:switch` bir konuşmayı kendiliğinden sürdürebilir ([ayrıntılar](/features/next#inside-claude-code)). Aksi halde ccam kendini Claude Code ile değiştirir ve yanında hiçbir şey çalışmaz.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--dry-run` | ccam'in çalıştıracağı ortamı ve komutu yazdırır, ardından çıkar | kapalı |
| `--fallback` | `-p`/`--print` çalıştırmaları için: hesap kullanım limitine ulaştığında aynı komutu en çok payı kalan girişli hesapla yeniden çalıştırır ([ayrıntılar](/features/next#run-again-at-the-limit)) | kapalı |
| `-v`, `--version` | `ccam <version>` yazdırır (yalnızca ilk argüman olarak) | |
| `-h`, `--help` | `ccam --help` veya `ccam run --help` olarak ccam'in yardımı | |

Argüman kuralları:

- ccam kendi seçeneklerini yalnızca hesap adından önce okur. Addan sonra `--dry-run` ve `--help` dahil her argüman Claude Code'a aittir.
- `--dry-run`, `ccam --dry-run work`, `ccam run --dry-run work` ve `ccam --dry-run` (adsız hesap) olarak çalışır. Reddedilmez, dikkate alınır. `--dry-run` ve `--fallback` herhangi bir sırayla yazılabilir.
- `ccam run <name|number> [claude args]`, `ccam <name>` komutunun uzun biçimidir. Adsız `ccam run` klasörün hesabını ya da varsayılan hesabı kullanır.
- Hesap adları `add`, `list`, `run` veya `usage` gibi ccam komut adları olamaz; böylece `ccam <word>` hiçbir zaman belirsiz olmaz.

Kuru çalıştırma, bir kabuğa ya da Windows'ta PowerShell'e yapıştırabileceğiniz tek bir satır yazdırır:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## Nasıl çalışır

Her başlatma, komut satırından `exec` çağrısına kadar aynı adımları izler:

<!-- motion:launch -->

Her kutu, adımı çalışırken yanar; alttaki satırlar ise Claude Code'un başladığı ortam ve argümanlardır.

1. `cmd/ccam/main.go`, `app.Main` fonksiyonunu çağırır. ccam başka bir adla başlatıldıysa bu bir [kısayoldur](/features/shortcuts) ve doğrudan başlatıcıya gider.
2. `app.parse`, cobra devreye girmeden önce `os.Args` değerine bakar. Yalnızca tek başına `ccam` (kontrol paneli), bir ccam komut adı (`config.Reserved`), `-h`/`--help` ve `-v`/`--version`, `cli.Execute` fonksiyonuna gider. Geri kalan her şey bir başlatmadır.
3. `launcher.Launch`, ayrıştırdığı tek dosya olan `~/.config/ccam/config.toml` dosyasını okur.
4. `launcher.Prepare` hesabı ada veya numaraya göre bulur (ad yoksa `HereProfile`, çalışma dizinini `config.toml` içindeki `[[folder]]` girdileriyle dize olarak karşılaştırıp klasörün hesabını alır), baştaki bir `--` değerini atar ve `claude` yürütülebilir dosyasını `ResolveClaude` ile bulur. PATH araması, dosyanın değiştirilme zamanıyla birlikte `~/.local/state/ccam/claude-path` içinde önbelleğe alınır; böylece ccam her başlatmada `PATH` içinde arama yapmaz. Değişmiş veya eksik bir dosya yeni bir aramayı tetikler. `config.toml` içindeki bir `claude_bin` yolu olduğu gibi kullanılır. ccam, inode ile karşılaştırarak ccam'in kendisi olan bir `claude` dosyasını reddeder, çünkü bu bir döngüye yol açar.
5. `checkLinks`, hesap klasörü `~/.local/share/ccam/profiles/<name>/` ile paylaşım ayarlarını `share.Check` ile karşılaştırır; bu fonksiyon yalnızca `lstat` ve `readlink` çağırır, öğe başına bir veya iki sistem çağrısı. Eksik bir bağlantı anında oluşturulur. Veri kaybettirebilecek her şey başlatmayı durdurur.
6. `applyMCP`, paylaşılan MCP sunucularını ekler (aşağıya bakın).
7. `session.Register` oturumu kaydeder. ccam, PID'sinin adını taşıyan ve süreç başlangıç zamanını içeren bir işaret dosyasını `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` konumuna (`XDG_RUNTIME_DIR` yoksa `~/.local/state/ccam/locks/` içine) yazar. Konuşma geçmişini paylaşan bir hesap için `historyCheck`, önce `history-<group>-<hash>` içinde çalışan bir oturum arar (hash, geçerli klasörün hash'idir) ve oraya da kaydolur. `exec` PID'yi korur; böylece işaret dosyası artık Claude Code sürecinin kendisini adlandırır. Tutulan bir kilit değil, bir işaret dosyası kullanılır; çünkü Claude Code'un başlattığı süreçler bir kilidi devralır ve oturum bittikten sonra da tutar.
8. `Plan.Exec`, terminal başlığını `claude · <name>` olarak ayarlar (standart hata bir terminal olduğunda), başlatma zamanını `fsync` olmadan `~/.local/state/ccam/last-used/<name>` dosyasına yazar ve `syscall.Exec` çağırır.

Claude Code, `CLAUDE_CONFIG_DIR` hesap klasörüne ve `CCAM_PROFILE` hesap adına ayarlanmış olarak çalışır. ccam, bunların devraldığı değerlerini kaldırır. macOS'ta `ccam import` ile kaydedilmiş bir `~/.claude` klasörü için ayrıca `CLAUDE_SECURESTORAGE_CONFIG_DIR` değişkenini boş bir değere ayarlar; böylece hesap ve düz `claude` aynı Anahtar Zinciri girişini kullanır.

### Neden exec

`syscall.Exec`, ccam sürecini `claude` ile değiştirir; böylece ccam bellekte kalmaz, tuş vuruşları, sinyaller ve çıkış kodu doğrudan Claude Code'a aittir ve argümanlar hiçbir zaman bir kabuk dizesi olmadan dizi olarak kalır.

Windows'ta `exec` yoktur. Orada `Plan.Exec`, `claude.exe` dosyasını aynı ortam, argümanlar ve konsolla bir alt süreç olarak başlatır, onu bekler ve onun çıkış koduyla çıkar. Ctrl+C, Claude Code'a doğrudan konsoldan ulaşır; bu yüzden ccam onu yok sayar. Oturum işaretçisi, tam olarak Claude Code kadar yaşayan ccam'in PID'sini taşır ve Claude Code çıktığında kaldırılır. Beklerken ccam başka hiçbir şey yapmaz. Bir oturumdan sonraki işler, örneğin `claude mcp add` ile eklenen sunucuları fark etmek, bir sonraki başlatmada yapılır: ccam, `.claude.json` dosyasının değiştirilme zamanını `~/.local/state/ccam/mcp/<name>/seen` ile karşılaştırır ve dosyayı yalnızca değiştiyse ayrıştırır.

### Bir komut satırı nasıl yönlendirilir

Üç komut satırı, `internal/app` içinden üç yol:

<!-- motion:launch-route -->

Bir kısayol `argv[0]` ile, `ccam <name>` ise ilk argümanla eşleşir ve ikisi de `syscall.Exec` ile biter. Yalnızca `ccam usage` gibi bir ccam komutu cobra'ya ulaşır.

### Paylaşılan MCP sunucuları

Varsayılan `flag` stratejisiyle ccam, grubun MCP sunucularını Claude Code'a sizinkilerden önce yerleştirilen bir argüman olarak iletir:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` yalnızca grubun `mcp.json` dosyası varsa bir argüman döndürür ve `Plan.InsertArgs` bunu `claude` komutunun hemen arkasına koyar. Kendi `--mcp-config` veya `--strict-mcp-config` seçeneğiniz bunun üzerine yine uygulanır. Bu tek bir argümandır, `--mcp-config=<path>`, çünkü Claude Code'un `--mcp-config` seçeneği birden fazla değer alır: iki argüman olarak yazılsaydı ardından gelen bir istemi de yapılandırma dosyası olarak okurdu.

`config.toml` içinde `mcp_strategy = "sync"` ile ccam bunun yerine sunucuları hesabın `.claude.json` dosyasına yazar ve hesap zaten çalışıyorsa bunu atlar.

### Remote Control

`ccam work rc` (veya `ccam work remote-control`), Claude Code'un Remote Control özelliğini `work` hesabıyla başlatır. `claude remote-control`, `--mcp-config` seçeneğini alt komutundan önce reddeder, sonrasında da hiç almaz; oturumları ise MCP sunucularını yalnızca hesabın `.claude.json` dosyasından okur. Bu yüzden bayrak stratejisi burada farklı çalışır:

<!-- motion:launch-rc -->

1. `remoteControl`, hesaptan sonra gelen ve `-` ile başlamayan ilk argümanı Claude Code'un alt komutu olarak alır. Bu `remote-control` veya `rc` olduğunda `mcp.FlagArgs` hiçbir şey eklemez.
2. `mcp.Synced`, grubun `mcp.json` dosyasının değiştirilme zamanını hesabın `sync-stamp` dosyasıyla karşılaştırır. `.claude.json` grubun güncel sunucularını zaten içeriyorsa veya grubun hiç sunucusu yoksa Remote Control soru sormadan başlar.
3. Aksi halde ccam, oturumlarda eksik kalacak sunucuları adlarıyla belirtir ve sorar:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y` yanıtında `mcp.Sync` bunları `sync` stratejisinin yaptığı gibi kopyalar: Claude Code'un kilidi altında, bir yedek aldıktan sonra ve hesabın kendi sunucularını koruyarak. `N` yanıtında veya terminal yoksa Remote Control onlar olmadan başlar ve ccam, bunları her başlatmadan önce kopyalayan `mcp_strategy = "sync"` ayarını önerir.

ccam, sunucular kopyalanana kadar her Remote Control başlatmasında sorar; kopyalamadan sonra ise yalnızca grubun `mcp.json` dosyası yeniden değiştiğinde sorar. `.claude.json` içindeki bir kopya, sonraki normal başlatmalarda `--mcp-config` seçeneğine üstün gelir; bu yüzden `ccam mcp add`, `remove` ve `import` diğer hesapların tuttuğu kopyaları günceller ([Paylaşılan MCP sunucuları](/features/mcp)).

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
- `ccam work rc`, `work` zaten çalışırken `.claude.json` dosyasına dokunmaz, bunu belirtir ve Remote Control'ü grubun sunucuları olmadan başlatır; bunları kopyalamak için oturum bittikten sonra yeniden çalıştırın. `--dry-run` ile mesajı yazdırır ama ne sorar ne de kopyalar.
- `mcp_strategy = "sync"` ile başarısız bir eşitleme `ccam: MCP sync failed (launching anyway): …` yazdırır ve Claude Code zaten sahip olduğu sunucularla başlar.

## İlgili

- [Kısayol komutları](/features/shortcuts)
- [ccam ile Remote Control kullanma](/guides/remote-control)
- [Hesapları yönetme](/features/accounts)
- [Komut başvurusu](/docs/commands)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [CLAUDE_CONFIG_DIR nasıl çalışır](/guides/claude-config-dir)
