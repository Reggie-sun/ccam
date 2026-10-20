---
title: ccam ile Remote Control kullanma
description: "Remote Control'ü ccam <name> rc ile her hesapla başlatın: geçerli seçenekler, ccam'in paylaşılan MCP sunucularını neden kopyalamayı önerdiği ve sorun giderme."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "5f0b90ea53d6"
---

`ccam work rc`, Claude Code'un Remote Control özelliğini `work` hesabıyla başlatır; böylece bu bilgisayardaki oturumları claude.ai/code üzerinden veya Claude mobil uygulamasından yönetebilirsiniz. `claude rc` gibi çalışır, tek bir farkla: hesap MCP sunucularını paylaşıyorsa ccam, bunları hesabın `.claude.json` dosyasına kopyalayıp kopyalamayacağını bir kez sorar, çünkü Remote Control oturumları sunucuları başka hiçbir yerden okumaz.

## Başlatma

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

Bir [kısayol komutu](/features/shortcuts) da aynı şekilde çalışır: `ccwork rc`, `ccam work rc` ile aynıdır. Komutu çalışmak istediğiniz klasörde çalıştırın, ardından `work` ile aynı Claude hesabıyla oturum açılmış claude.ai/code'u veya Claude uygulamasını açın.

## Seçenekler

Burada iki tür seçenek bir araya gelir: hesap adından önce yazılan veya ccam ayarlarında bulunan ccam'in kendi seçenekleri ve `rc` sonrasına yazılıp Claude Code'a değiştirilmeden ulaşan Remote Control seçenekleri.

### ccam seçenekleri ve ayarları

| Seçenek veya ayar | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--dry-run` (hesaptan önce) | Ortamı ve komutu, varsa MCP bildirimini de yazdırır. Ne soru sorar ne de kopyalar | kapalı |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Hesap MCP sunucularını paylaştığında ve `.claude.json` dosyasında grubun güncel listesi olmadığında sorulur. `y` bunları kopyalar, `N` veya Enter onlar olmadan başlatır | `N` |
| `~/.config/ccam/config.toml` içinde `mcp_strategy = "sync"` | Remote Control olsun olmasın, her başlatmadan önce grubun sunucularını `.claude.json` dosyasına kopyalar; böylece hiçbir zaman soru sorulmaz | `"flag"` |
| `ccam group leave <group> work` | Hesabı, onunla MCP paylaşan gruptan çıkarır; Remote Control bu durumda yalnızca hesabın kendi sunucularıyla başlar ve soru sorulmaz | `default` içinde |

### Remote Control seçenekleri

`rc` sonrasındaki her şey Claude Code'a aittir. Bunlar Claude Code 2.1.304 sürümündeki `claude remote-control` seçenekleridir; `ccam work rc --help` sizin sürümünüzün listesini yazdırır.

| Seçenek | Ne yapar |
| --- | --- |
| `--name <name>` | claude.ai/code'da gösterilen oturum adı |
| `--remote-control-session-name-prefix <prefix>` | Otomatik oluşturulan oturum adları için önek (varsayılan: ana makine adı) |
| `-c`, `--continue` | Remote Control'ün bu klasör için en son kaydettiği oturuma, yaklaşık 4 saat içinde yeniden bağlanır |
| `--session-id <id>` | Belirli bir oturuma yeniden bağlanır |
| `--permission-mode <mode>` | Yeni oturumlar için izin modu: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` veya `plan` |
| `--spawn <mode>` | `same-dir` (varsayılan), `worktree` (oturum başına bir git worktree) veya `session` (tek oturum, oturum bitince çıkar) |
| `--capacity <N>` | `same-dir` ve `worktree` modunda aynı anda en fazla oturum sayısı (varsayılan 32) |
| `--[no-]create-session-in-dir` | Başlangıçta geçerli klasörde ilk oturumu oluşturur (varsayılan açık) |
| `--[no-]chrome` | Yeni oturumlar için Claude in Chrome (varsayılan kapalı) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Hata ayıklama ve ayrıntılı çıktı |

`rc` komutunu hesap adının hemen arkasına yazın. ccam, `-` ile başlamayan ilk argümanı Claude Code'un alt komutu olarak alır; bu yüzden `ccam work rc --name laptop` içinde `rc` görür, ama `ccam work --model opus rc` içinde `opus` görür.

## ccam ne yapar

Başlatma her zamanki yolu izler ([Claude Code'u başlatma](/features/launch)), yalnızca MCP adımında bir değişiklik vardır:

1. ccam, her başlatmada olduğu gibi hesabı bulur, paylaşım bağlantılarını denetler ve `claude` yolunu çözer.
2. Alt komutun `rc` veya `remote-control` olduğunu görür ve `--mcp-config` eklemez.
3. Hesap MCP sunucularını paylaşıyorsa grubun `mcp.json` dosyasını `~/.local/state/ccam/mcp/work/` içindeki `sync-stamp` ile karşılaştırır. `.claude.json` güncel sunucuları zaten içeriyorsa veya grubun hiç sunucusu yoksa doğrudan 6. adıma geçer.
4. Aksi halde oturumlarda eksik kalacak sunucuları adlarıyla belirtir:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. `y` yanıtında `.claude.json` dosyasının bir yedeğini alır, Claude Code'un kendi kilidini alır, grubun sunucularını `mcpServers` içine ekler ve hesabın diğer sunucularını korur. `N` yanıtında, terminal yokken veya hesap çalışırken hiçbir şeyi değiştirmez ve bunu belirtir:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Oturumu kaydeder, `CLAUDE_CONFIG_DIR` değişkenini hesap klasörüne ayarlar ve `exec` ile kendini `claude rc` ile değiştirir.

Remote Control bundan sonra bir sunucu olarak çalışmaya devam eder. Sizin için başlattığı her oturum bu sürecin bir alt sürecidir, dolayısıyla `work` olarak çalışır: `work` hesabının girişi, ayarları ve konuşma geçmişiyle çalışır ve `work` hesabının kullanım sınırından düşülür. Sunucu çalıştığı sürece ccam hesabı çalışıyor sayar.

Deneme çalıştırması farkı gösterir. Normal bir başlatma grubun sunucularını bir argüman olarak taşır, Remote Control taşımaz:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Neden böyle çalışır

**Remote Control `--mcp-config` almaz.** Normal başlatmalarda ccam, grubun sunucularını Claude Code'a `--mcp-config=<path>` olarak verir ve Claude Code'un dosyalarına asla yazmaz. `claude remote-control` bu seçeneği alt komut adından önce reddeder, sonrasında da hiç almaz; bu yüzden ccam 0.6.0 öncesinde `ccam work rc`, MCP sunucularını paylaşan her hesapta başarısız oluyordu. ccam artık bu seçeneği eklemiyor.

**Oturumları sunucuları yalnızca `.claude.json` dosyasından okur.** Remote Control'ün başlattığı oturumlar komut satırındaki seçenekleri devralmaz. Kullanıcı MCP sunucularını buldukları tek yer hesabın `.claude.json` dosyasıdır. Bu dosya hesabın giriş bilgilerini ve klasör başına durumu da tuttuğu için ccam onu `settings.json` ile yaptığı gibi paylaşılan bir kopyaya bağlayamaz. Sunucuları bu dosyaya kopyalamak, oturumlara paylaşılan listeyi vermenin tek yoludur.

**ccam yazmadan önce sorar.** Varsayılan bayrak stratejisinde ccam `.claude.json` dosyasına dokunmamayı taahhüt eder. Remote Control oraya yazmak zorunda olduğu tek durumdur; bu yüzden önce sorar, bir yedek tutar ve `mcpServers` dışında hiçbir şeye dokunmaz.

**Hesap çalışırken yazmaz.** Claude Code oturum sırasında `.claude.json` dosyasını yeniden yazar ve dışarıdan yapılan bir değişiklik kaybolur ya da üzerine yazılır. Bu yüzden ccam, hesabın herhangi bir oturumu (başka bir Remote Control sunucusu dahil) çalışırken kopyalamayı atlar.

**Her değişiklikte bir kez sorar.** Kopyalamadan sonra `sync-stamp`, grubun `mcp.json` dosyasının değiştirilme zamanını tutar. ccam yalnızca bu dosya değiştiğinde, örneğin `ccam mcp add` sonrasında yeniden sorar. `N` yanıtını verirseniz henüz bir kopya olmadığından her Remote Control başlatmasında sorar.

**`ccam mcp` kopyaları güncel tutar.** Adlar çakıştığında `.claude.json` içindeki bir girdi `--mcp-config` seçeneğine üstün gelir; bu yüzden eski bir kopya normal başlatmalarda da daha yeni bir tanımı gizler. Bu nedenle `ccam mcp add`, `add-json`, `remove` ve `import`, çalışmayan her hesaptaki kopyaları günceller:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Terminal yoksa yanıt hayırdır.** Remote Control çoğu zaman bir betikten, bir hizmetten veya `nohup` ile, kimsenin yanıt veremeyeceği bir yerden çalıştırılır. ccam bu durumda beklemek yerine sunucular olmadan başlar.

## Sorunlar ve çözümleri

### Remote Control oturumlarında paylaşılan MCP sunucuları eksik

Soruya `N` yanıtı verildi, terminal yoktu veya hesap çalışıyordu. Remote Control'ü durdurun, hesabın diğer oturumlarını kapatın ve `ccam work rc` komutunu bir terminalde yeniden çalıştırın; `y` yanıtını verin. Hiç sorulmaması için `~/.config/ccam/config.toml` içinde `mcp_strategy = "sync"` ayarlayın. Bir oturumdaki sunucular oturum çalışırken değişmez; bu yüzden kopyalamadan sonra yeni bir oturum başlatın.

### "Account work is running, so its .claude.json is left alone"

`work` hesabının bir oturumu hâlâ açık: normal bir `ccam work`, bir kısayol komutu veya başka bir Remote Control sunucusu. Ondan çıkın ve `ccam work rc` komutunu yeniden çalıştırın. Başka bir hesabın çalıştırdığı oturumlar sayılmaz.

### ccam her seferinde soruyor

Sunucuları henüz kopyalamadınız. Bir kez `y` yanıtını verin veya `mcp_strategy = "sync"` ayarlayın. Remote Control'de sunucuları hiç istemiyorsanız hesapla MCP paylaşımını bırakın: hesabı o gruptan çıkarın (`ccam group leave default work`) ya da MCP'yi hesabın içinde olmadığı ayrı bir grupta tutun.

### `ccam work rc` hemen `--mcp-config` ile ilgili bir hatayla başarısız oluyor

ccam sürümünüz 0.6.0'dan eski ve Remote Control'e `--mcp-config` iletiyor. `ccam update` komutunu çalıştırın. Güncel bir ccam ile aynı hata, `ccam work --model opus rc` örneğindeki gibi `rc` öncesinde değer alan bir seçenek olduğu ve ccam'in `opus` değerini alt komut sandığı anlamına gelir; `rc` komutunu hesap adının hemen arkasına yazın.

### `ccam mcp` ile yapılan bir değişiklik Remote Control'e ulaşmıyor

`ccam mcp add` ve `remove` kopyaları yalnızca çalışmayan hesaplarda günceller ve bunu belirtir: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control'ü durdurup yeniden başlatın; ccam değişen `mcp.json` dosyasını görür ve yeniden sorar.

### Kopya benim kendi sunucumun yerini aldı

Hesabın bir grup sunucusuyla aynı adı taşıyan kendi sunucusu varsa kopya onun yerini almıştır. Önceki `.claude.json`, en yeni beş yedeği tutan `~/.local/state/ccam/mcp/work/backups/` içindedir. Kendi sürümünüzü tüm hesaplar için korumak istiyorsanız bir sonraki Remote Control başlatmasından önce onu `ccam mcp import work <server> --force` ile gruba taşıyın.

### Kopyalanan sunucuları kaldırmak istiyorum

`ccam mcp remove <server>` bir sunucuyu gruptan ve tüm kopyalardan kaldırır. Kopyalamadan önceki duruma dönmek için hesabın hiçbir oturumu çalışmıyorken `~/.local/state/ccam/mcp/work/backups/` içindeki en yeni yedeği geri yükleyin.

### Remote Control klasöre güvenilip güvenilmeyeceğini soruyor veya terminal olmadan çıkıyor

Klasör güveni her hesap için ayrı olarak `.claude.json` içinde saklanır; bu yüzden `personal` olarak güvendiğiniz bir klasör `work` için yenidir. O klasörde bir kez `ccam work` çalıştırıp onay verin, ardından Remote Control'ü orada başlatın.

### "You must be logged in with a Claude account that has a subscription"

Remote Control, abonelikli bir Claude hesabıyla oturum açılmasını gerektirir. Hesabın e-postasını `ccam list` ile kontrol edin ve gerekirse `ccam login work` ile yeniden oturum açın.

### Oturum telefonda görünmüyor

Claude uygulamasında veya claude.ai/code'da farklı bir Claude hesabıyla oturum açılmış. Bunu `ccam list` komutunun `work` için gösterdiği e-postayla karşılaştırın.

### ccam, bu klasörde başka bir Claude Code'un konuşma geçmişini paylaştığı konusunda uyarıyor

Aynı geçmiş grubundaki başka bir hesap aynı klasörde çalışıyor. Her iki hesaptan aynı konuşmayı sürdüren oturumlar geçmişi birbirine karıştırabilir. Yalnızca farklı konuşmalar üzerinde çalışacaklarsa `y` yanıtını verin veya Remote Control'ü başka bir klasörde başlatın.

## İlgili

- [Claude Code'u başlatma](/features/launch)
- [Paylaşılan MCP sunucuları](/features/mcp)
- [Hesaplar arasında paylaşım](/features/share)
- [Hesapları yönetme](/features/accounts)
