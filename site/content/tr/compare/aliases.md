---
title: ccam ile kabuk takma adlarının karşılaştırması
description: Birden fazla Claude Code hesabı çalıştırmak için CLAUDE_CONFIG_DIR ile kabuk takma adları ve ccam, yan yana karşılaştırma.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: f751716e46c9
---

Her hesap için bir kabuk takma adı (alias), birden fazla Claude Code hesabı çalıştırmanın en basit yoludur. ccam arka planda aynı şeyi yapar ve buna paylaşım, hesap değiştirme ve kullanım takibi ekler.

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
| Durum satırında hesap adı | Hayır | `ccam statusline install` |
| Bozuk bağlantıları onarma | Elle | `ccam doctor --fix` |
| Ek yazılım | Yok | Tek bir statik ikili dosya |

İki hesabınız varsa ve aralarında hiçbir şeyin paylaşılmasını istemiyorsanız takma adları seçin. Hesapların tek bir kurulum gibi hissettirmesini istiyorsanız ccam'ı seçin. Takma adlarla başlayıp daha sonra ccam'a geçebilirsiniz: `ccam import`, mevcut `~/.claude-*` klasörlerini bulundukları yerde kaydeder.

[Tek bilgisayarda iki hesap nasıl kullanılır](/guides/multiple-claude-code-accounts) kılavuzu, iki yaklaşımı da adım adım açıklar.
