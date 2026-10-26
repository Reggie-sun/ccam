---
title: Her hesabın arka plan oturumları
description: ccam agents, tüm hesaplarınızın Claude Code arka plan ajanlarını ve açık oturumlarını tek bir tabloda listeler; ccam attach birini kimliğiyle açar.
updated: 2026-10-26
order: 16
translation: machine
sourceHash: 06281fb13119
---

Claude Code arka plan oturumlarını (`claude --bg`, ajan görünümü) hesap başına tutar: `claude agents` yalnızca kendisini çalıştıran hesabın oturumlarını gösterir. `ccam agents` tüm hesaplara aynı anda sorar ve oturumları birlikte listeler; `ccam attach <id|name>` ise bir arka plan oturumunu, onu hangi hesap başlattıysa o hesapla açar.

## Kullanım

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State`, arka plan oturumunun durumunu (`working`, `blocked`, `done`, `failed`, `stopped`) ve sürecinin ne yaptığını (`busy`, `waiting`, `idle`) gösterir; açık bir etkileşimli oturum `open` gösterir.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--all` | Bitmiş, başarısız olmuş veya durdurulmuş arka plan oturumlarını da listeler | kapalı |
| `--json` | Oturumları, her biri ait olduğu `account` ile birlikte JSON olarak yazdırır | kapalı |

`ccam attach <id|name> [claude args]`, kimlikten sonraki her şeyi, seçenekler dahil, değiştirmeden `claude attach` komutuna iletir: `ccam attach 49ec --dry-run`, `--dry-run` seçeneğini Claude Code'a verir. Bunun yerine komutu yazdırmak için `--dry-run` seçeneğini kimlikten önce yazın (`ccam --dry-run attach 49ec` veya `ccam attach --dry-run 49ec`).

## Nasıl çalışır

1. ccam her hesap için, hesabın `CLAUDE_CONFIG_DIR` değeriyle `claude agents --json` komutunu (verildiyse `--all` ile) çalıştırır; hepsini paralel olarak, her birini 15 saniyelik bir zaman aşımıyla. Bu, hesabın kendi oturum listesini okur ve Claude Code'un arka plan hizmetini başlatmaz.
2. Her girdi hesabın adını alır ve liste başlangıç zamanına göre, en yeni önce olacak şekilde sıralanır. Yanıt vermeyen bir hesap standart hataya adıyla yazılır ve listeye alınmaz; ardından ccam 1 çıkış koduyla çıkar, böylece bir betik bir hatayı boş liste sanmaz.
3. `ccam attach`, Claude Code'un bir arka plan oturumunu tuttuğu `jobs/<id>` yolunu her hesabın klasöründe arar ve `claude attach <id>` komutunu, `ccam <name>` komutunun Claude Code'u başlattığı şekilde, o hesapla çalıştırır. Bir kimlik öneki olmayan bir sözcük, büyük/küçük harf fark etmeksizin oturum adlarıyla eşleştirilir; kimliğe benzeyen bir sözcük asla ad olarak alınmaz.

## Güvenlik ve uç durumlar

- Arka plan oturumu kimlikleri tek bir hesaba aittir: başka bir hesapla çalıştırılan `claude attach` onları bulamaz; `ccam attach` tam olarak bu sorunu çözer.
- İki hesaptaki oturumlar aynı kimlik önekini ya da ad sözcüğünü paylaşıyorsa `ccam attach` onları listeler ve daha uzun bir kimlik ya da ad ister.
- ccam paylaşılan MCP sunucularını Claude Code'a tek bir argümanda `--mcp-config=<file>` ile verir; `claude agents` ve `claude attach` bunu kabul eder ve arka plan oturumları yeniden başladıklarında o dosyayı kullanmaya devam eder.
- API anahtarlı bir hesabın arka plan hizmeti anahtarını kendisini başlatan kabuktan alır. Böyle bir hesabın ilk arka plan oturumunu ccam üzerinden (`ccam gw --bg …`) başlatın; böylece hizmet anahtara sahip olur.

## İlgili

- [Claude Code'u başlatma](/features/launch)
- [Betikler için hesap seçme](/features/pick)
- [Kontrol paneli](/features/dashboard)
