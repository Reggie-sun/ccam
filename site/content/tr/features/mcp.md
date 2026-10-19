---
title: Paylaşılan MCP sunucuları
description: ccam mcp add, add-json, list, remove ve import ile grup başına tek bir MCP sunucu listesi tutun; liste --mcp-config ile iletilir veya her hesaba eşitlenir.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 46a7aa716f99
---

Claude Code, kullanıcı MCP sunucularını hesabın girişini de tutan dosya olan `.claude.json` içinde saklar; bu yüzden bu dosya hesaplar arasında paylaşılamaz. ccam, paylaşılan sunucuları bir grup dosyasında, `mcp.json` içinde tutar ve başlatmada bunu grubun her hesabına verir. Bir sunucuyu tüm hesaplarınız için bir kez eklemek üzere `ccam mcp` komutunu, bir oturumda eklediğiniz bir sunucuyu diğer hesaplara taşımak için `ccam mcp import` komutunu kullanın.

## Kullanım

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Bir değişiklik her hesabın bir sonraki başlatılmasından itibaren geçerli olur; çalışan oturumlar başladıkları sunucuları korur:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

Grupta zaten bulunan bir adı eklemek, onun tanımını değiştirir ve `added` yerine `updated` yazdırır. Sunucular bir hesaba yalnızca hesap MCP'yi o grupla paylaşıyorsa ulaşır (`ccam share <name> --mcp on`); henüz hiçbir hesap MCP paylaşmıyorsa `add` bunu belirtir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json` ve `remove` için hedef grup | bir hesabın MCP'yi üzerinden paylaştığı tek grup, yoksa `default` |
| `-- <command> [args…]` (`add`) | bu komutu bir stdio sunucusu olarak çalıştırır | |
| `--url <url>` (`add`) | uzak sunucu URL'si (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | bir `--url` sunucusunun aktarım türü | `http` |
| `--header <'Name: value'>` (`add`) | HTTP başlığı, tekrarlanabilir | yok |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio sunucusu için ortam değişkeni, tekrarlanabilir | yok |
| `--force` (`import`) | aynı adlı ama farklı bir grup sunucusunu değiştirir; eski `mcp.json` yedeklenir | kapalı |
| `--keep-local` (`import`) | hesabın `.claude.json` dosyasındaki girdileri korur | kapalı |
| `-h`, `--help` | her komutun yardımı | |
| `--dry-run` (genel) | yalnızca `mcp list` dikkate alır; `add`, `add-json`, `remove` ve `import` reddeder ve hiçbir şeyi değiştirmez | kapalı |

`add`, ya `--` sonrasında bir komut ya da `--url` gerektirir; aksi halde `a command to run (after --) or --url is required` hatasıyla başarısız olur. `--header`, `--url` sunucularına, `--env` ise komut sunucularına uygulanır. `add-json`, `claude mcp add-json` ile aynı JSON'u alır ve tanım `command` veya `url` içeren bir nesne olmalıdır. `import` her zaman hesabın MCP'yi üzerinden paylaştığı grubu kullanır, `list` ise her grubu gösterir; bu yüzden `-g` bunları değiştirmez.

## Nasıl çalışır

İlgili dosyalar:

- `~/.local/share/ccam/groups/<group>/mcp.json`: paylaşılan sunucular, `{"mcpServers": {…}}`, yani `--mcp-config` seçeneğinin okuduğu biçim. Sunucu girdileri sıklıkla API anahtarları taşıdığı için `fsutil.AtomicWrite` ile ve `0600` moduyla yazılır.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: Claude Code'un hesabın sunucularını tutan kendi dosyası. ccam yalnızca `mcpServers` değerini `profile.TopLevelValue` tarayıcısıyla okur ve başlatma yolunda dosyanın tamamını hiçbir zaman ayrıştırmaz.
- `~/.local/state/ccam/mcp/<name>/`: `seen` dosyası, `sync-stamp` ve `.claude.json` yedekleri.

`~/.config/ccam/config.toml` içindeki `mcp_strategy` ile ayarlanan iki strateji vardır. Varsayılan, bayrak stratejisidir:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` hesabın MCP grubunu bulur. `mcp.FlagArgs`, grubun `mcp.json` dosyasının var olduğunu onu ayrıştırmadan tek bir `stat` ile kontrol eder.
2. `Plan.InsertArgs`, `--mcp-config=<path>` değerini `claude` komutunun hemen arkasına, sizin argümanlarınızdan önce koyar; böylece kendi `--mcp-config` veya `--strict-mcp-config` seçeneğiniz bunun üzerine yine uygulanır.
3. `mcp.NewLocal`, hesabın `.claude.json` dosyasının son başlatmadan bu yana değişip değişmediğini kontrol eder ve yalnızca değiştiyse sunucularını ayrıştırır.
4. `syscall.Exec`, ccam'i Claude Code ile değiştirir.

Yol iki değil, tek bir argüman olarak iletilir: `--mcp-config=<path>`. Claude Code'un `--mcp-config` seçeneği birden fazla değer alır; bu yüzden `claude --mcp-config a.json "prompt"` istemi ikinci bir yapılandırma dosyası olarak okur ve `MCP config file not found` ile başarısız olur. Tek argümanlı biçim orada durur ve ardından gelen bir istem veya `mcp list` gibi bir alt komut çalışır. Claude Code, `--mcp-config` sunucularını hesabın kendi sunucularına ek olarak yükler; ad çakışmasında `.claude.json` içindeki girdi kazanır. Bayrak stratejisi hiçbir zaman `.claude.json` dosyasına yazmaz.

Bir oturumda `claude mcp add` ile eklenen sunucu yalnızca o hesabın `.claude.json` dosyasına yazılır. ccam o anda bunu göremez, çünkü Claude Code onun yerini almıştır; bu yüzden bir sonraki başlatmada kontrol eder. `seen` dosyası, `.claude.json` dosyasının değiştirilme zamanını ve daha önce bildirilmiş sunucu adlarını tutar. Zaman farklıysa ccam adları grupla karşılaştırır ve her yeni sunucu için girdi beklemeden bir kez ipucu yazdırır:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` ile sunucular `claude mcp list` ve `/mcp` içinde sıradan kullanıcı sunucuları olarak görünür ve ccam bunları başlatmadan önce her hesabın `.claude.json` dosyasına yazar:

<!-- motion:mcp-sync -->

1. `mcp.Sync`, grubun `mcp.json` dosyasının değiştirilme zamanını `sync-stamp` ile karşılaştırır. Değişmediyse hesabın `.claude.json` dosyası açılmaz.
2. Claude Code'un kendi kilidini, yani dosyanın yanında `mkdir` ile oluşturulan `.claude.json.lock` klasörünü alır. En fazla 2 saniye bekler ve Claude Code'un kilit kütüphanesinin yaptığı gibi 10 saniyeden eski bir kilidi bayat sayar.
3. Mevcut `.claude.json` dosyasını durum klasöründeki `backups/` içine kaydeder ve en yeni 5 yedeği tutar.
4. `mcpServers` içinde her grup sunucusunu ayarlar, daha önce yazdığı ama grupta artık bulunmayan sunucuları kaldırır, hesabın kendi sunucularını korur ve dosyayı özgün moduyla atomik olarak yazar.
5. Yeni damgayı ve yönettiği adları kaydeder; Claude Code `--mcp-config` olmadan başlar.

Hesabın bir oturumu çalışırken eşitleme `ccam: account work is already running, skipping MCP sync` ile atlanır, çünkü Claude Code çalışırken `.claude.json` dosyasını yeniden yazar. Başarısız bir eşitleme `ccam: MCP sync failed (launching anyway): …` yazdırır ve yine de başlatır. Oturum içinde eklenen sunucular için verilen ipucu yalnızca bayrak stratejisine aittir.

`ccam mcp import`, bir hesabın kendi sunucularını grubuna taşır:

<!-- motion:mcp-import -->

ccam hesabın sunucularını `mcp.LocalServers` ile okur ve her birini, anahtar sırasını ve boşlukları yok sayan `mcp.SameDef` ile grubun aynı adlı sunucusuyla karşılaştırır. Grupta bulunmayan veya aynı tanımla bulunan bir sunucu kopyalanır. Farklı olan atlanır, çünkü diğer tüm hesapların kullandığı tanım grubun tanımıdır:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` ile grubun sunucusu değiştirilir ve önceki `mcp.json` önce grup klasörünün içindeki `.ccam-backup/<time>/mcp.json` konumuna kaydedilir. Ardından `mcp.RemoveLocal`, aynı kilit altında ve durum klasörüne bir yedek aldıktan sonra kopyalanan girdileri hesabın `.claude.json` dosyasından kaldırır. Bu gereklidir, çünkü ad çakışmasında hesabın girdisi grubunkine üstün gelir. `--keep-local` bunları yerinde bırakır.

## Güvenlik ve uç durumlar

- **Oturum açmalar.** OAuth ile oturum açan sunucular tokenlarını sunucu listesinde değil, her hesabın `.credentials.json` dosyasında (`mcpOAuth`) tutar; bu yüzden bir sunucuyu paylaşmak oturum açmasını paylaşmaz. Her hesapta `/mcp` ile bir kez oturum açın veya oturum açmaları [`ccam migrate`](/features/migrate) ile `~/.claude` klasöründen kopyalayın. Grup hiçbir zaman bir giriş dosyası almaz.
- **Çalışan hesaplar.** `import` gruba kopyalar ama çalışan bir hesabın kendi girdilerini yerinde bırakır ve bunu belirtir; oturum bittikten sonra yeniden çalıştırın. Eşitleme çalışan hesapları atlar.
- **Bozuk dosyalar.** `.claude.json` geçerli bir JSON değilse eşitleme ona dokunmaz ve hatayı bildirir. Kilit 2 saniye içinde alınamazsa ccam `could not lock … (Claude Code is writing the file)` bildirir.
- **MCP paylaşmayan hesaplar.** `import` reddeder ve önce `ccam share side --mcp on` çalıştırmanızı söyler.
- **Grup adları.** `-g` yalnızca bir grup adı alır; bu ad `groups/` altında bir klasör olur, bu yüzden `../x` gibi adlar reddedilir.

## İlgili

- [Hesaplar arasında paylaşım](/features/share)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Mevcut bir kurulumu kopyalama](/features/migrate)
- [Tüm komutlar](/docs/commands)
