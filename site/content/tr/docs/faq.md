---
title: Sık sorulan sorular
description: ccam, Claude Code hesapları, kullanım limitleri, paylaşım, gizlilik ve desteklenen sistemler hakkında yanıtlar.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: e07d468d453b
---

## ccam resmi bir Anthropic aracı mı?

Hayır. ccam resmi olmayan, açık kaynaklı bir projedir; Anthropic ile bağlantılı değildir ve Anthropic tarafından onaylanmamıştır. Birden fazla hesabı yan yana çalıştırmak için Claude Code'un belgelediği `CLAUDE_CONFIG_DIR` ayarını kullanır.

## Birden fazla Claude Code hesabı kullanmama izin var mı?

Bu, her hesabın koşullarına ve iş hesapları için kurumunuzun politikasına bağlıdır. Birden fazla hesap kullanmadan önce bunları kontrol edin. ccam yalnızca zaten sahip olduğunuz hesaplar arasında geçiş yapmayı kolaylaştırır.

## ccam kullanım limitlerini atlatmanızı sağlar mı?

Hayır. Her hesabın kendi limitleri geçerli olmaya devam eder. ccam bu limitleri gösterir ve bir hesabın limiti dolduğunda sahip olduğunuz başka bir hesapta devam etmenizi sağlar.

## Tek bilgisayarda iki Claude Code hesabını nasıl kullanırım?

Her hesaba `CLAUDE_CONFIG_DIR` ile kendi yapılandırma klasörünü verin. ccam bunu sizin için yapar: her hesap için `ccam add` komutunu çalıştırın, ardından hesapları `ccam personal` ve `ccam work` ile ya da adını kendiniz seçtiğiniz `ccwork` gibi kısayollarla başlatın. [Kılavuz](/guides/multiple-claude-code-accounts), bunun elle nasıl yapılacağını da gösterir.

## Bir konuşmaya başka bir hesapta nasıl devam ederim?

İki hesap arasında konuşma geçmişini `ccam share <name> --history on` ile paylaşın, ardından `ccam next --continue` komutunu çalıştırın. Claude Code, geçerli klasörün en son konuşmasını bir sonraki hesapta açar.

## ccam giriş bilgilerimi görüyor veya gönderiyor mu?

ccam hiçbir zaman bir girişi kopyalamaz veya paylaşmaz. Bir hesabın erişim tokenını yalnızca `ccam usage`, `ccam next --auto` veya `ccam doctor` çalıştırdığınızda, o hesabın kullanımını Anthropic'in kullanım uç noktasına sormak için okur. ccam'da telemetri yoktur ve başka hiçbir yere bir şey göndermez.

## Hesaplarım nerede saklanıyor?

Her hesap için bir klasör olmak üzere `~/.local/share/ccam/profiles/<name>/` içinde. macOS üzerinde Claude Code, girişin kendisini Anahtar Zinciri'nde, o klasörden türetilen bir ad altında saklar.

## Mevcut ~/.claude klasörüme ne olur?

`ccam import` veya `ccam migrate` çalıştırmadığınız sürece hiçbir şey. İçe aktarma, klasörü bulunduğu yerde kaydeder; böylece düz `claude` ve `ccam <name>` aynı hesabı ve ayarları kullanmaya devam eder. `ccam migrate` ise klasörün ayarlarını, `CLAUDE.md` dosyasını, uzantılarını ve MCP sunucularını paylaşılan ayarlara kopyalar ve klasörü yalnızca okur.

## ccam hangi sistemleri destekliyor?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine ve diğerleri) ve macOS. Windows üzerinde WSL içinde kullanın.

## Bir MCP sunucusu neden yalnızca bir hesapta görünüyor?

Bir oturum içinde `claude mcp add` ile eklenen sunucular yalnızca o hesap için saklanır. ccam, hesabı bir sonraki başlatışınızda bunu size bildirir; sunucuyu `ccam mcp import <account> <server>` ile paylaşılan gruba taşıyın.

## ccam, Claude Code'u yavaşlatır mı?

Claude Code'u ccam üzerinden başlatmak yaklaşık 4 milisaniye ekler; bu süre bir dizüstü bilgisayarda `scripts/perf-check.sh` ile ölçülmüştür. ccam küçük bir dosya okur, paylaşılan bağlantıları kontrol eder ve ardından kendini `claude` ile değiştirir; dolayısıyla siz çalışırken çalışır durumda değildir.

## ccam'ı nasıl kaldırırım?

Her hesabın kısayolunu `ccam shortcut <account> --remove` ile kaldırın, `ccam` komutunu `~/.local/bin` içinden silin; hesaplara artık ihtiyacınız yoksa ardından `~/.config/ccam`, `~/.local/share/ccam` ve `~/.local/state/ccam` klasörlerini de silin. [Kurulum](/docs/install#uninstall) sayfasına bakın.
