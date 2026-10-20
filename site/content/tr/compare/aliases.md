---
title: ccam ile kabuk takma adlarının karşılaştırması
description: Birden fazla Claude Code hesabı çalıştırmak için CLAUDE_CONFIG_DIR ile kabuk takma adları ve ccam, yan yana karşılaştırma.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: b9fbcb88f111
---

Her hesap için bir kabuk takma adı (alias), birden fazla Claude Code hesabı çalıştırmanın en basit yoludur. ccam arka planda aynı şeyi yapar ve buna paylaşım, hesap değiştirme ve kullanım takibi ekler.

## Takma adlarla kurulum

Her ek hesap için `~/.bashrc` veya `~/.zshrc` dosyasına bir satır ekleyin, yeni bir terminal açın ve giriş yapmak için takma adı bir kez çalıştırın:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` ilk hesabınız olarak kalır, `claude-work` ise ikincisini başlatır. Bu, [Claude Code belgelerinin anlattığı](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) yöntemdir.

## ccam ile kurulum

ccam'i kurun ve hesabı ekleyin; bu örnekte hesaba kendi kısayol komutu da veriliyor:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` ve `ccwork` hesabı başlatır ve ikisi de betiklerde çalışır.

## Yan yana

| | Kabuk takma adı | ccam |
| --- | --- | --- |
| Ayrı girişler | Evet | Evet |
| Hesap başına kurulum | `~/.bashrc` dosyasını düzenleme, giriş yapma | `ccam add` |
| Betiklerde ve diğer programlarda çalışır | Hayır, takma adlar yalnızca kabukta çalışır | Evet, `ccam work` ve `ccwork` gibi kısayollar gerçek komutlardır |
| Paylaşılan ayarlar, `CLAUDE.md`, uzantılar | Elle kopyalama veya bağlama | Öğe bazında, yedeklerle |
| Paylaşılan MCP sunucuları | Elle, hesap başına | Bir kez `ccam mcp add` |
| Bir konuşmaya başka hesapta devam etme | Hayır | `ccam next --continue` |
| Tüm hesapların kullanımı | Hayır | `ccam usage` |
| En çok pay kalan hesabı seçme | Hayır | `ccam next --auto` |
| Limitte bir `-p` işini başka bir hesapla yeniden çalıştırma | Hayır | `ccam --fallback` |
| Her proje klasörü için bir hesap | Hayır | `ccam folder` |
| Durum satırında hesap adı | Hayır | `ccam statusline install` |
| Bozuk bağlantıları onarma | Elle | `ccam doctor --fix` |
| Ek yazılım | Yok | Tek bir statik ikili dosya |

İki hesabınız varsa ve aralarında hiçbir şeyin paylaşılmasını istemiyorsanız takma adları seçin. Hesapların tek bir kurulum gibi hissettirmesini istiyorsanız ccam'ı seçin. Takma adlarla başlayıp daha sonra ccam'a geçebilirsiniz: `ccam import`, mevcut `~/.claude-*` klasörlerini bulundukları yerde kaydeder.

[Tek bilgisayarda birden fazla hesap nasıl kullanılır](/guides/multiple-claude-code-accounts) kılavuzu iki yaklaşımı da adım adım açıklar; [ccam ve diğer hesap değiştirme araçları](/compare/account-switchers) ise ccam'ı tek bir girişi değiştiren araçlarla karşılaştırır.
