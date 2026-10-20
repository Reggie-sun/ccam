---
title: API anahtarlı hesaplar
description: ccam add --api, Anthropic API'si veya bir ağ geçidi için API anahtarlı hesap ekler; anahtar başlatmada bir değişkenden ya da komuttan okunur, saklanmaz.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

API anahtarlı bir hesap, Claude Code'u Claude abonelik girişi yerine bir API anahtarıyla başlatır: Anthropic API'si ya da aynı Messages biçimini konuşan bir uç nokta, örneğin bir LiteLLM veya OpenRouter ağ geçidi. ccam anahtarı saklamaz. Hesap her başladığında anahtarı bir ortam değişkeninden ya da parola yöneticisi gibi bir komuttan okur ve yalnızca o Claude Code sürecine verir.

## Kullanım

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

API anahtarlı bir hesap, diğer her hesap gibi ayarları, `CLAUDE.md` dosyasını, uzantıları ve MCP sunucularını paylaşır; kısayollarla, `ccam folder` ve `ccam ext` ile de çalışır.

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `--api` | Giriş yapmak yerine API anahtarlı bir hesap ekler | kapalı |
| `--key-from env:NAME` | Anahtarı `NAME` ortam değişkeninden okur | |
| `--key-from command:COMMAND` | `COMMAND` komutunu çalıştırır (`sh -c` ile, Windows'ta `cmd /C` ile) ve yazdırdığı ilk satırı kullanır | |
| `--base-url <url>` | Uç nokta (`ANTHROPIC_BASE_URL`) | Anthropic API'si |
| `--auth api-key\|bearer` | Anahtarı `ANTHROPIC_API_KEY` (`x-api-key`) olarak ya da `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, ağ geçitlerinin genellikle istediği) olarak gönderir | `api-key` |
| `--model <model>` | Model (`ANTHROPIC_MODEL`) | Claude Code'un varsayılanı |

Ayarlar `config.toml` içinde saklanır:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Nasıl çalışır

Hesabın her başlatılışında, olağan kontrollerden sonra ve Claude Code başlamadan önce:

1. ccam anahtarı okur: `env:` ortamdan, `command:` komutu 30 saniyelik bir sınırla çalıştırarak. Bir komut terminalde parola isteyebilir; hata çıktısı da terminale gider. Boş bir anahtar, anahtarın nerede beklendiğini söyleyen bir mesajla başlatmayı durdurur.
2. Claude Code'un nasıl giriş yaptığını veya hangi sağlayıcıyı kullandığını belirleyen tüm devralınan değişkenleri kaldırır: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` ve `CLAUDE_CODE_USE_FOUNDRY`. Kabuğunuzda başka bir şey için ayarlanmış bir anahtar bu hesaba ulaşamaz ve bir abonelik girişi onun anahtarının önüne geçemez.
3. Verildiyse `ANTHROPIC_BASE_URL` ve `ANTHROPIC_MODEL` değişkenlerini, anahtarı da `ANTHROPIC_API_KEY` veya `ANTHROPIC_AUTH_TOKEN` olarak ayarlar.

Abonelik girişli hesaplar, önceden olduğu gibi devraldıkları ortamı korur. Etkileşimli bir oturumda Claude Code, ortamda bulduğu anahtarı kullanıp kullanmayacağını bir kez sorabilir.

## Güvenlik ve uç durumlar

- Anahtar hiçbir zaman `config.toml` dosyasına, bir önbelleğe veya bir günlüğe yazılmaz ve `ccam --dry-run gw` onu okumaz: bir `command:` kaynağı kuru çalıştırmada çalıştırılmaz.
- `ccam add --api` bir `env:` kaynağını hemen kontrol eder ve değişken boşsa uyarır; bir `command:` kaynağı ilk kez başlatmada çalıştırılır.
- `ccam login` ve `ccam logout`, girişi olmadığı için API anahtarlı bir hesabı reddeder. `ccam doctor` onun uç noktasını ve anahtar kaynağını gösterir ve bir `env:` değişkeni boşsa uyarır.
- `ccam usage` onun için `API` gösterir. `ccam next --auto`, `--fallback` ve `ccam pick` yalnızca abonelik hesapları arasından seçer, çünkü bir API anahtarının karşılaştırılacak 5 saatlik veya haftalık limiti yoktur.
- Paylaşılan bir `settings.json`, içinde `env` bloğu veya `apiKeyHelper` varsa, API anahtarlı hesaplar dahil gruptaki her hesaba uygulanır. Sağlayıcı ayarlarını paylaşılan ayarların dışında tutun.

## İlgili

- [Hesapları yönetme](/features/accounts)
- [Claude Code'u başlatma](/features/launch)
- [Komut başvurusu](/docs/commands)
