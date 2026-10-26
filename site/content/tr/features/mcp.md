---
title: Paylaşılan MCP sunucuları
description: ccam mcp add, add-json, list, remove ve import ile grup başına tek bir MCP sunucu listesi tutun; liste --mcp-config ile iletilir veya her hesaba eşitlenir.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: 27a06ad35295
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
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
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
only in account side, for /home/me/code/app only: notes
```

Grupların altında `ccam mcp list` diğer hesaplara ulaşmayan her sunucuyu adlandırır:

- bir hesabın kendisi için eklediği bir sunucu için `only in account work: pg`;
- tek bir klasör için eklenmiş bir sunucu için (düz `claude mcp add`, aşağıya bakın) `only in account side, for /home/me/code/app only: notes`;
- bir hesabın kendi sunucusu bir grup sunucusunun adını taşıdığında `account work has its own …, used in its sessions in place of the group's`;
- bu klasörün `.mcp.json` dosyası, onayını beklediği hesaplarla birlikte (`waiting for approval in`);
- `in no account, only in ~/.claude.json (plain claude's)` olan sunucular, onları kopyalayan [`ccam migrate --mcp`](/features/migrate) komutuyla birlikte.

Grupta zaten bulunan bir adı eklemek, onun tanımını değiştirir ve `added` yerine `updated` yazdırır. Sunucular bir hesaba yalnızca hesap o gruptaysa ve grup MCP'yi paylaşıyorsa ulaşır (`ccam group set <group> --add mcp`); grup paylaşmıyorsa `add` bunu belirtir.

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
| `--dry-run` (genel) | yalnızca `mcp list` destekler; `add`, `add-json`, `remove` ve `import` reddeder ve hiçbir şeyi değiştirmez | kapalı |

`add`, ya `--` sonrasında bir komut ya da `--url` gerektirir; aksi halde `a command to run (after --) or --url is required` hatasıyla başarısız olur. `--header`, `--url` sunucularına, `--env` ise komut sunucularına uygulanır. `add-json`, `claude mcp add-json` ile aynı JSON'u alır ve tanım `command` veya `url` içeren bir nesne olmalıdır. `import` her zaman hesabın MCP'yi üzerinden paylaştığı grubu kullanır, `list` ise her grubu gösterir; bu yüzden `-g` bunları değiştirmez.

## Nasıl çalışır

İlgili dosyalar:

- `~/.local/share/ccam/groups/<group>/mcp.json`: paylaşılan sunucular, `{"mcpServers": {…}}`, yani `--mcp-config` seçeneğinin okuduğu biçim. Sunucu girdileri sıklıkla API anahtarları taşıdığı için `fsutil.AtomicWrite` ile ve `0600` moduyla, yalnızca `<group>/.mcp.json.lock` kilidi altında yazılır; böylece aynı anda çalıştırılan iki `ccam mcp add` komutu da kendi sunucularını korur.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: Claude Code'un hesabın sunucularını tutan kendi dosyası. ccam yalnızca `mcpServers` değerini `profile.TopLevelValue` tarayıcısıyla okur ve başlatma yolunda dosyanın tamamını hiçbir zaman ayrıştırmaz.
- `~/.local/state/ccam/mcp/<name>/`: `seen` dosyası, `sync-stamp` (ccam'in yazdığı her kopyanın bir hash'iyle; böylece kimsenin değiştirmediği bir kopyayı sizin düzenlediğiniz bir kopyadan ayırabilir), `sync-checked` (eşitleme stratejisinin en son gördüğü), `kept` (hesabın MCP paylaşımından ayrılırken koruduğu sunucular) ve `.claude.json` yedekleri.

`~/.config/ccam/config.toml` içindeki `mcp_strategy` ile ayarlanan iki strateji vardır. Varsayılan, bayrak stratejisidir:

<!-- motion:mcp-flag -->

1. ccam hesabın MCP grubunu bulur (`launcher.applyMCP`) ve grubun `mcp.json` dosyasının var olduğunu onu ayrıştırmadan tek bir `stat` ile kontrol eder (`mcp.FlagArgs`).
2. ccam `--mcp-config=<path>` değerini `claude` komutunun hemen arkasına, sizin argümanlarınızdan önce koyar; böylece kendi `--mcp-config` veya `--strict-mcp-config` seçeneğiniz bunun üzerine yine uygulanır (`Plan.InsertArgs`).
3. ccam hesabın `.claude.json` dosyasının son başlatmadan bu yana değişip değişmediğini kontrol eder ve yalnızca değiştiyse sunucularını ayrıştırır (`mcp.NewLocal`).
4. ccam kendini Claude Code ile değiştirir (`syscall.Exec`).

Yol iki değil, tek bir argüman olarak iletilir: `--mcp-config=<path>`. Claude Code'un `--mcp-config` seçeneği birden fazla değer alır; bu yüzden `claude --mcp-config a.json "prompt"` istemi ikinci bir yapılandırma dosyası olarak okur ve `MCP config file not found` ile başarısız olur. Tek argümanlı biçim bunu önler; böylece ardından gelen bir istem veya alt komut çalışır. Claude Code orada bayrağı reddettiği için `ccam work mcp list`, `get` ve `remove` bayrak almaz; ccam bunun yerine grubun sunucularını standart hataya yazar ([ayrıntılar](/features/launch#shared-mcp-servers)). Claude Code, `--mcp-config` sunucularını hesabın kendi sunucularına ek olarak yükler; ad çakışmasında `.claude.json` içindeki girdi kazanır. Bayrak stratejisi `.claude.json` dosyasına yalnızca Remote Control için, bir hesap MCP paylaşımından ayrıldığında veya yeniden katıldığında ([aşağıda](#leaving-and-joining)) ve `ccam mcp import` için yazar.

`claude remote-control` (`rc`) `--mcp-config` almaz ve oturumları sunucuları yalnızca `.claude.json` dosyasından okur. `ccam work rc` için ccam bayrağı eklemez ve bir terminalde, grubun sunucularının `.claude.json` dosyasına kopyalanıp kopyalanmayacağını sorar; kopyalama, `sync` stratejisinde olduğu gibi bir yedekle yapılır. Bir kopya oluşana kadar sorar, sonra yalnızca grubun `mcp.json` dosyası değiştikten veya kopyalanmış bir sunucu kaybolduktan sonra yeniden sorar ([Remote Control](/features/launch#remote-control)). Böyle bir kopya ad çakışmasında `--mcp-config` seçeneğine üstün geldiği için `ccam mcp add`, `add-json`, `remove` ve `import`, grubun diğer hesaplarının tuttuğu kopyaları günceller:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

`✓ Updated` satırı yalnızca bir kopya değiştiğinde görünür. ccam'in daha eski bir sürümünün oluşturduğu bir kopya da, `ccam migrate` sonrasındaki bir kopya da aynı şekilde güncellenir.

Bir oturumda `claude mcp add` ile eklenen sunucu yalnızca o hesabın `.claude.json` dosyasına yazılır. Düz `claude mcp add` yerel kapsamı kullanır: sunucu yalnızca o tek hesap için, o tek klasörde yüklenir (`-s user` onu hesap için her klasörde kullanılabilir yapar). ccam bu olduğu anda bunu göremez, çünkü Claude Code ccam sürecinin yerini çoktan almıştır; bu yüzden bir sonraki başlatmada kontrol eder. `seen` dosyası, `.claude.json` dosyasının değiştirilme zamanını ve daha önce bildirilmiş sunucu adlarını tutar. Zaman farklıysa ccam adları grupla karşılaştırır ve her yeni sunucu için, iki stratejide de, girdi beklemeden bir kez ipucu yazdırır:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

İpucu, bir grup sunucusunun adını taşıyan, hesabın kendine ait bir sunucusunu da adlandırır, çünkü hesabın oturumlarında onunki kazanır.

`mcp_strategy = "sync"` ile sunucular `claude mcp list` ve `/mcp` içinde sıradan kullanıcı sunucuları olarak görünür ve ccam bunları başlatmadan önce her hesabın `.claude.json` dosyasına yazar:

<!-- motion:mcp-sync -->

1. ccam, grubun `mcp.json` dosyasının değiştirilme zamanını `sync-stamp` ile, hesabın `.claude.json` dosyasını da en son gördüğüyle karşılaştırır (`mcp.Sync`). İkisi de değişmediyse başka hiçbir şey olmaz. Aksi halde her grup sunucusunun hâlâ `.claude.json` içinde olduğunu da denetler ve kaybolmuş bir kopyayı geri koyar (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Claude Code'un kendi kilidini, yani dosyanın yanında `mkdir` ile oluşturulan `.claude.json.lock` klasörünü alır. En fazla 2 saniye bekler ve Claude Code'un kilit kütüphanesinin yaptığı gibi 10 saniyeden eski bir kilidi bayat sayar.
3. Mevcut `.claude.json` dosyasını durum klasöründeki `backups/` içine kaydeder ve en yeni 5 yedeği tutar.
4. `mcpServers` içinde her grup sunucusunu ayarlar, daha önce yazdığı ama grupta artık bulunmayan sunucuları kaldırır ve dosyayı, yalnızca bir şey değiştiyse, özgün moduyla atomik olarak yazar. Yalnızca kendi yazdığı ve o zamandan beri kimsenin değiştirmediği kopyaları değiştirir veya kaldırır. Hesabın bir grup sunucusuyla aynı adı taşıyan kendi sunucusu varsa ccam hesabın sunucusunu korur ve bunu bir kez söyler: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Yeni damgayı her kopyanın bir hash'iyle birlikte kaydeder; Claude Code `--mcp-config` olmadan başlar.

Hesabın bir oturumu çalışırken, Claude Code nasıl başlatılmış olursa olsun, eşitleme `ccam: account work is already running, skipping MCP sync` ile atlanır, çünkü Claude Code çalışırken `.claude.json` dosyasını yeniden yazar. Başarısız bir eşitleme `ccam: MCP sync failed (launching anyway): …` yazdırır ve yine de başlatır.

### Ayrılma ve katılma

MCP paylaşmayı bırakan bir hesap (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) grubun sunucularını korur: plan onları bir yedekle birlikte hesabın `.claude.json` dosyasına kopyalar ve hesabın aynı adlı kendi sunucuları kalır. `--fresh` ile hiçbirini almaz ve plan yüklenmesi duracak sunucuları listeler. Bu sunuculara sahip bir gruba yeniden katıldığında, kimsenin değiştirmediği korunan kopyalar kaldırılır; böylece grubun sonraki düzenlemeleri ona ulaşır.

`ccam mcp import`, bir hesabın kendi sunucularını grubuna taşır:

<!-- motion:mcp-import -->

ccam hesabın sunucularını, kendi sunucularını ve yalnızca tek bir klasör için tuttuklarını okur ve her birini, anahtar sırasını ve boşlukları yok sayan `mcp.SameDef` ile grubun aynı adlı sunucusuyla karşılaştırır. Bir klasör sunucusu, her klasör için bir grup sunucusu olur (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); iki farklı klasör tanımına sahip tek bir ad reddedilir. Eklenti sunucuları önerilmez, çünkü onları eklenti getirir. Grupta bulunmayan veya aynı tanımla bulunan bir sunucu kopyalanır. Farklı olan atlanır, çünkü diğer tüm hesapların kullandığı tanım grubun tanımıdır:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` ile grubun sunucusu değiştirilir ve önceki `mcp.json` önce grup klasörünün içindeki `.ccam-backup/<time>/mcp.json` konumuna kaydedilir. Ardından ccam, aynı kilit altında ve durum klasörüne bir yedek aldıktan sonra kopyalanan girdileri, hesabın kendi girdilerini ve klasörün girdilerini tek bir yazmada, hesabın `.claude.json` dosyasından kaldırır (`mcp.Remove`). ccam grubun `mcp.json` dosyasını okuduğu andan bu girdiler kaldırılana kadar kilitli tutar; böylece arada kimse sunucuyu gruptan kaldıramaz. Bu gereklidir, çünkü ad çakışmasında hesabın girdisi grubunkine üstün gelir. `--keep-local` ile girdiler hesabın `.claude.json` dosyasında kalır.

## Güvenlik ve uç durumlar

- **Oturum açmalar.** OAuth ile oturum açan sunucular tokenlarını sunucu listesinde değil, her hesabın `.credentials.json` dosyasında (`mcpOAuth`) tutar; bu yüzden bir sunucuyu paylaşmak oturum açmasını paylaşmaz. Her hesapta `/mcp` ile bir kez oturum açın. ccam oturum açmaları asla kopyalamaz, [`ccam migrate`](/features/migrate) de kopyalamaz: kopyalanan bir yenileme tokenı onu tutan diğer herkesin oturumunu kapatır. [`ccam ext`](/features/ext) ve `ccam doctor`, başka bir hesap oturum açmışken bir hesabın oturum açmadığı sunucuları adlarıyla belirtir; `ccam next --continue` de yeni hesapta hangilerinin eksik olduğunu söyler. Grup hiçbir zaman bir giriş dosyası almaz.
- **Çalışan hesaplar.** `import` gruba kopyalar ama çalışan bir hesabın kendi girdilerini yerinde bırakır ve bunu belirtir; oturum bittikten sonra yeniden çalıştırın. Eşitleme çalışan hesapları atlar. İkisi de ccam olmadan başlatılan Claude Code oturumlarını da algılar.
- **Remote Control için kopyalar.** Çalışan bir hesabın kopyasına dokunulmaz: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` ile hiçbir şeyin güncellenmesi gerekmez, çünkü her başlatma eşitleme yapar.
- **Bozuk dosyalar.** `.claude.json` geçerli bir JSON değilse eşitleme ona dokunmaz ve hatayı bildirir. Kilit 2 saniye içinde alınamazsa ccam `could not lock … (Claude Code is writing the file)` bildirir.
- **MCP paylaşmayan hesaplar.** `import` reddeder ve önce hesabı MCP paylaşan bir gruba koymanızı söyler (`ccam group join <group> side`).
- **Grup adları.** `-g` yalnızca bir grup adı alır; bu ad `groups/` altında bir klasör olur, bu yüzden `../x` gibi adlar reddedilir.

## İlgili

- [Gruplarla paylaşım](/features/share)
- [ccam ile Remote Control kullanma](/guides/remote-control)
- [MCP sunucularını, ayarları ve geçmişi paylaşma](/docs/sharing)
- [Mevcut bir kurulumu kopyalama](/features/migrate)
- [Komut başvurusu](/docs/commands)
