---
title: Komut başvurusu
metaTitle: ccam komut başvurusu
description: ccam add ve ccam <name> komutlarından ccam usage, ccam next ve ccam doctor komutlarına kadar tüm ccam komutları ve seçenekleri.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 5d3275878535
---

Bir komutun ne yapacağını, yapmadan görmek için komuttan önce `--dry-run` yazın, örneğin `ccam --dry-run work`. Bu, Claude Code başlatılırken (`ccam <name>`, `ccam run`, `ccam next`) ve `share`, `migrate`, `shortcut`, `update` ile `doctor --fix` komutlarında çalışır. Bir şeyi değiştiren diğer komutlar çalışmak yerine bunu reddeder.

<!-- commands -->

## Claude Code'u başlatma

`ccam <name|number> [arguments]`, Claude Code'u o hesapla başlatır. ccam, addan sonraki argümanları okumaz; bunları aynı sırayla `claude` komutuna iletir. Standart girdi, çıkış kodu ve sinyaller de tam olarak `claude` komutunun kendisinde olduğu gibi davranır. Addan sonra `--` yazmak isteğe bağlıdır.

Hesap adı olmadan argümanlarla çalıştırılan `ccam`, örneğin `ccam -p "summarise"`, geçerli klasör için ayarlanan hesabı, yoksa varsayılan hesabı kullanır. Bunları `ccam folder <name> [folder]` ve `ccam use <name>` ile ayarlayın.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Hesap adından önce gelen `--fallback`, `-p`/`--print` çalıştırmaları içindir. Hesap kullanım limitine ulaştığında ccam aynı komutu en çok payı kalan girişli hesapla yeniden çalıştırır ve yalnızca tamamlanan çalıştırmanın çıktısını yazdırır. [Klasöre göre hesaplar](/features/folder) ve [Limitte yeniden çalıştırma](/features/next#run-again-at-the-limit) sayfalarına bakın.

`ccwork` gibi bir kısayol, kendi hesabını aynı şekilde başlatır: `ccwork --resume`, `ccam work --resume` ile aynı işi yapar. [Kısayollar](/docs/commands#shortcuts) bölümüne bakın.

## Hesap ekleme ve kaldırma

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

Klasör belirtilmeden çalıştırılan `ccam import`, `~/.claude` ve `~/.claude-*` klasörlerini arar ve her biri için sorar. İçe aktarılan klasörler bulundukları yerde kaydedilir; ccam hiçbir zaman bir girişi taşımaz veya kopyalamaz.

`ccam rename`, `ccam remove` ve `ccam share`, Claude Code bir hesapla çalışırken o hesabı değiştirmeyi reddeder. Önce oturumdan çıkın.

`ccam list --json` her hesabı yuvası, kısayolu, e-postası ve klasörüyle birlikte yazdırır. Çıktıdaki `share` nesnesi her öğeyi (`mcp`, `settings`, `memory`, `extensions`, `history`) bir grup adıyla veya `local` ile eşler.

## Kısayollar

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Kısayol, adını kendiniz verdiğiniz ve Claude Code'u tek bir hesapla başlatan bir komuttur. `ccam work` gibi her argümanı Claude Code'a değiştirmeden iletir; bu yüzden `ccwork --resume` ve `ccwork -p "hi"` beklediğiniz gibi çalışır. `~/.local/bin` içinde ccam'e giden bir bağlantıdır.

Kısayollar isteğe bağlıdır. `ccam add` sihirbazı bir kısayol oluşturulup oluşturulmayacağını sorar, `ccam add <name> --shortcut <command>` sormadan oluşturur ve `ccam list` her hesabın kısayolunu, yoksa `-` gösterir.

Adlarda a-z, 0-9, `.`, `-` ve `_` karakterleri kullanılır. Bir ad bir ccam komutu, başka bir hesabın kısayolu veya `PATH` içindeki başka bir programın adı olamaz.

`ccam rename` kısayolun adını korur, `ccam remove` ise kısayolu hesapla birlikte siler. `ccam shims install`, örneğin ccam'i taşıdıktan sonra, tüm hesapların bağlantılarını yeniden oluşturur.

## Paylaşım

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Her öğe `on` (`default` grubu), `off` veya bir grup adı alır. [Paylaşım](/docs/sharing) sayfasına bakın.

## Mevcut bir kurulumu kopyalama

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate`, `settings.json`, `CLAUDE.md`, uzantıları ve kullanıcı MCP sunucularını bir Claude Code klasöründen bir gruba kopyalar; böylece bu öğeleri paylaşan her hesap onlara sahip olur. Klasör yalnızca okunur. `--settings`, `--memory`, `--extensions`, `--mcp` ve `--history` seçenekleri yalnızca belirttiğiniz öğeleri kopyalar. Bunlar olmadan, terminalde bir kontrol listesinden seçim yaparsınız ve konuşma geçmişi dışındaki her şey işaretlidir; diğer durumlarda konuşma geçmişi dışındaki her şey kopyalanır. `-g` veya `--group` grubu seçer; `-y` veya `--yes` sormadan kopyalar ve bir dosya farklıysa grubun kopyasını korur. Ayrıca sunucuların kaydedilmiş girişlerini her hesabın kendi giriş dosyasına kopyalar. [Mevcut bir kurulumu kopyalama](/docs/sharing#copy-an-existing-setup) bölümüne bakın.

## MCP sunucuları

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>`, yalnızca o hesapta bulunan sunucuları hesabın grubuna taşır. Adı grupta zaten farklı bir tanım için kullanılan sunucular atlanır; grubun kopyasını değiştirmek için `--force` ekleyin, ccam önce o kopyayı yedekler.

## Kullanım ve geçiş

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Bakım

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update`, ccam'ı değiştirmeden önce sürümün SHA-256 sağlama toplamını kontrol eder; `--version v0.2.0` ise belirli bir sürümü kurar. `.deb` paketinden kurulan ccam bunun yerine apt ile güncellenir. [Güncelleme](/docs/install#update) bölümüne bakın.

## Dosyalar

| Yol | İçerik |
| --- | --- |
| `~/.config/ccam/config.toml` | Hesaplar ve neyi paylaştıkları |
| `~/.local/share/ccam/profiles/<name>/` | Her hesabın `CLAUDE_CONFIG_DIR` klasörü |
| `~/.local/share/ccam/groups/<group>/` | Paylaşılan ayarlar, MCP sunucuları, uzantılar ve geçmiş |
| `~/.local/state/ccam/` | Kullanım önbelleği ve son kullanım zamanları |
| `~/.local/bin/<shortcut>` | `ccwork` gibi kısayol komutları, ccam'e giden bağlantılar |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ve `XDG_STATE_HOME` değişkenleri bu klasörlerin yerini değiştirir. Windows'ta `~`, `%USERPROFILE%` anlamına gelir ve kısayollar `<shortcut>.exe` olur.
