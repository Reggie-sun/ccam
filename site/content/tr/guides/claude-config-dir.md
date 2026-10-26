---
title: CLAUDE_CONFIG_DIR nedir?
description: CLAUDE_CONFIG_DIR'ın Claude Code'da ne yaptığı, klasörde neler saklandığı, macOS üzerinde girişin yeri ve hesaplar arasında hangi dosyaların paylaşılabileceği.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: dc2fe42e11e1
---

`CLAUDE_CONFIG_DIR`, Claude Code'a `~/.claude` yerine hangi klasörü kullanacağını söyler. Hesaba özgü her şey orada bulunur: giriş bilgileri, ayarlar, MCP sunucuları, eklentiler ve konuşma geçmişi. Bu değişkeni farklı klasörlere yönlendirmek size ayrı hesaplar sağlar.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR resmi olarak destekleniyor mu?

Evet. Claude Code onu [ortam değişkenleri](https://code.claude.com/docs/en/env-vars) arasında, yapılandırma klasörünü değiştirmenin yolu olarak listeler ve "birden fazla hesabı yan yana çalıştırmak için kullanışlı" olduğunu belirtir; [kimlik doğrulama kılavuzu](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) da birden fazla hesapla giriş yapmak için onu kullanır. Bu değişkeni kabuğunuzda veya kullanıcı ayarlarında ayarlayabilirsiniz, ancak bir projenin ayarlarında ayarlayamazsınız.

## Klasörün içinde neler var?

Claude Code 2.1.309 ile kontrol edilmiştir:

| Yol | İçerik | Paylaşılabilir |
| --- | --- | --- |
| `.credentials.json` | Giriş tokenları (Linux, Windows), MCP OAuth tokenları | Hayır |
| `.claude.json` | Hesap bilgileri, kullanıcı MCP sunucuları, proje bazında durum | Hayır (MCP sunucuları ayrı olarak) |
| `settings.json` | İzinler, hook'lar, model, durum satırı, eklentiler | Evet |
| `CLAUDE.md` | Genel talimatlar | Evet |
| `agents/`, `commands/`, `skills/`, `plugins/` | Uzantılar | Evet |
| `projects/` | `.jsonl` dosyaları olarak konuşmalar, proje bazında bellek | Evet |
| `history.jsonl` | İstem geçmişi | Evet |
| `sessions/`, `shell-snapshots/`, `cache/` ve diğerleri | Çalışma zamanı durumu | Hayır |

## İnsanlar en çok nerede takılıyor?

**`CLAUDE_CONFIG_DIR` olmadan `.claude.json`, `~/.claude` içinde değildir.** Claude Code bu durumda onu klasörün yanında, `~/.claude.json` konumunda tutar. `CLAUDE_CONFIG_DIR=~/.claude` kullanmaya başlarsanız Claude Code bunun yerine `~/.claude/.claude.json` dosyasını arar ve hesabınızı bulamaz. `ccam import` bu durumu ikisini birbirine bağlayarak çözer. `~/.claude` klasöründe, `CLAUDE_CONFIG_DIR=~/.claude` ile yapılmış önceki bir çalıştırmadan kalma, kendine ait bir `.claude.json` zaten varsa ccam iki dosyayı da bırakır ve bunu söyler, çünkü hangisinin doğru olduğunu yalnızca siz bilirsiniz.

**macOS üzerinde giriş bilgileri klasörde değil, Anahtar Zinciri'ndedir.** Claude Code bunları `Claude Code-credentials` adı altında saklar. `CLAUDE_CONFIG_DIR` ayarlandığında ad ayrıca `-` ve klasör yolunun SHA-256 değerinin ilk 8 onaltılık basamağıyla biter. Anahtar Zinciri öğesinin hesap adı `$USER` olur. Bu nedenle ayrı klasörler ayrı Anahtar Zinciri kayıtları alır ve taşınan veya yeniden adlandırılan bir klasör girişini kaybeder; `ccam rename` macOS'ta klasörün adını bu yüzden korur. Anahtar Zinciri yazmayı reddederse Claude Code, [belgelerinde](https://code.claude.com/docs/en/authentication#credential-management) anlatıldığı gibi klasördeki `.credentials.json` dosyasına geri döner.

**Her giriş klasörde saklanmaz.** API anahtarı olmadan yapılan bir Claude Console girişi klasörün dışında saklanır; bu yüzden ayrı klasörler bu türden iki girişi birbirinden ayıramaz. claude.ai girişleri ve API anahtarları klasör bazında tutulur.

## Windows'ta klasör nerede?

Varsayılan olarak `%USERPROFILE%\.claude` konumundadır. Giriş bilgileri bunun içindeki `.credentials.json` dosyasında ya da Claude Code bunu kullanacak şekilde ayarlandıysa, macOS'ta olduğu gibi klasör yolundan türetilen bir ad altında, Windows Kimlik Bilgileri Yöneticisi'ndedir. Linux ve macOS'ta olduğu gibi, her ek hesap için `CLAUDE_CONFIG_DIR` değişkenini başka bir klasöre ayarlayın.

## Klasörler dosya paylaşabilir mi?

Sembolik bağlantılar çalışır: Claude Code, bağlanmış bir `settings.json` dosyası veya `projects/` klasörü üzerinden okur ve yazar, bağlantıyı da korur. MCP sunucuları bunun istisnasıdır, çünkü giriş bilgileriyle aynı `.claude.json` dosyasını paylaşırlar; bunun yerine onları `--mcp-config` ile iletin.

ccam tam olarak bu kuralların üzerine kuruludur. Klasörleri nasıl bağladığı için [paylaşım](/docs/sharing) sayfasına, genel tablo için de [birden fazla hesap kılavuzuna](/guides/multiple-claude-code-accounts) bakın.
