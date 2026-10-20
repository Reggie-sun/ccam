---
title: ccam kurulumu
metaTitle: ccam'ı Linux, macOS ve Windows üzerine kurun
description: ccam'ı Linux, macOS veya Windows'a tek komutla, bir .deb paketiyle ya da kaynak koddan kurun, ccam update ile güncelleyin. Önce Claude Code kurulu olmalıdır.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: "07fab3165030"
---

Linux veya macOS üzerinde bir terminalde şunu çalıştırın:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Betik, sisteminize uygun sürümü indirir, SHA-256 sağlama toplamını doğrular ve `ccam` komutunu `~/.local/bin` klasörüne yerleştirir; ccam daha sonra seçtiğiniz `ccwork` gibi kısayol komutlarını da aynı klasörde oluşturur. Bu klasör `PATH` içinde değilse betik, `~/.bashrc` veya `~/.zshrc` dosyasına eklemeniz gereken satırı yazdırır.

Windows'ta bunu PowerShell'de çalıştırın:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Betik `ccam_windows_<arch>.zip` dosyasını indirir, SHA-256 sağlama toplamını denetler, `ccam.exe` dosyasını Claude Code yükleyicisinin `claude.exe` dosyasını koyduğu klasöre, yani `%USERPROFILE%\.local\bin` içine koyar ve henüz orada değilse bu klasörü kullanıcı `PATH` değişkeninize ekler. Ardından yeni bir terminal açın.

ccam, Claude Code'u sizin için başlatır; bu nedenle önce [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) kurun ve `claude --version` komutunun çalıştığından emin olun.

## Desteklenen sistemler

- x86-64 ve ARM64 üzerinde Linux: Ubuntu, Debian, Fedora, Arch, Alpine ve diğer dağıtımlar. İkili dosya statiktir ve hiçbir bağımlılığı yoktur.
- Apple silicon ve Intel üzerinde macOS.
- x86-64 ve ARM64 üzerinde Windows 10 ve 11; PowerShell, Windows Terminal veya Git Bash içinde. Paylaşılan klasörler her kullanıcı için çalışır. Paylaşılan dosyalar (`settings.json`, `CLAUDE.md`, `history.jsonl`) Geliştirici Modu (Ayarlar > Sistem > Geliştiriciler için) gerektirir, çünkü Windows dosyalara bağlantı oluşturmaya yalnızca bu modda veya bir yöneticiye izin verir; bkz. [Windows'ta paylaşım](/features/share#on-windows).
- WSL içinde ccam, herhangi bir Linux sisteminde olduğu gibi çalışır.

## Diğer kurulum yolları

**Debian ve Ubuntu paketi.** `ccam_<version>_linux_amd64.deb` dosyasını [sürümler sayfasından](https://github.com/jhyunwoo/ccam/releases) indirin ve kurun:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Paket, `/usr/bin/ccam` dosyasını kurar.

**Kaynak koddan.** Go 1.26 veya daha yeni bir sürümle:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Belirli bir sürüm.** `CCAM_VERSION` değişkenini ayarlayın:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Kabuk tamamlama

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Tamamlama, hesap adları da dahil olmak üzere `ccam` komutunun tüm alt komutları ve seçenekleri için çalışır.

## Claude Code eklentisi

ccam eklentisi Claude Code'un kendisine ccam'i kullanmayı öğretir: hangi hesapta pay kaldığını kontrol etmeyi, farklı hesaplarda paralel `claude -p` çalışan süreçlerini `ccam pick` ile başlatmayı ve eksik bir beceriyi veya MCP sunucusunu `ccam ext` ile açıklamayı. Claude Code içinde kurun:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Bir beceri, `/ccam:usage`, `/ccam:ext` ve `/ccam:switch` komutlarını ve iki hook ekler: bir oturum kullanım limitinde durduğunda bir terminal bildirimi alırsınız ve `/ccam:switch [account]` konuşmayı başka bir hesaba devreder ([ayrıntılar](/features/next#inside-claude-code)). Hook'lar `ccam` komutunu çalıştırır; bu yüzden ccam, Claude Code'un gördüğü `PATH` üzerinde olmalıdır. Beceri, giriş yapan veya dosya taşıyan komutları hiçbir zaman çalıştırmaz; bunun yerine onları önerir.

## Kurulumu kontrol edin

```sh
ccam version
ccam doctor
```

`ccam doctor`, Claude Code'un bulunduğunu ve her hesap klasörünün, bağlantının ve girişin yerinde olduğunu kontrol eder. [Hızlı başlangıç](/docs/quick-start) ile devam edin.

## Güncelleme

```sh
ccam update --check
ccam update
```

`ccam update --check` yalnızca daha yeni bir sürüm olup olmadığını söyler. `ccam update` ise sisteminize uygun sürümü indirir, SHA-256 sağlama toplamını sürümün `checksums.txt` dosyasıyla karşılaştırır ve ardından `ccam` dosyasını tek adımda değiştirir. Sağlama toplamı eşleşmezse hiçbir şey değiştirilmez. `ccwork` gibi kısayollar ccam'a giden bağlantılar olduğundan yeni sürümü hemen çalıştırır.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam güncellemeleri hiçbir zaman kendiliğinden kontrol etmez; yeni bir sürümü yalnızca bu komutu çalıştırdığınızda arar. Daha eski bir sürüm de dahil olmak üzere belirli bir sürümü kurmak için `--version` seçeneğini verin. `ccam --dry-run update`, hiçbir şeyi değiştirmeden neyin indirilip neyin değiştirileceğini gösterir.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Kurulum betiğini yeniden çalıştırmak da `~/.local/bin` içindeki kopyayı günceller. ccam'ı `go install` ile derlediyseniz aynı komutu yeniden çalıştırın.

**Windows.** `ccam update`, Windows için `.zip` dosyasını indirir. Windows'ta çalışan bir program değiştirilemez; bu yüzden eski ikili dosya `ccam.exe.old` olarak yeniden adlandırılır ve bir sonraki güncelleme onu siler.

**Debian ve Ubuntu paketi.** `ccam update`, apt ile kurulmuş bir `/usr/bin/ccam` dosyasını değiştirmez; bunun yerine bu adımları yazdırır. Yeni `ccam_<version>_linux_amd64.deb` dosyasını [sürümler sayfasından](https://github.com/jhyunwoo/ccam/releases) indirin ve eskisinin üzerine kurun:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Kaldırma

Önce, ccam hâlâ kuruluyken hesap kısayollarını kaldırın. `ccam list` bunları Shortcut sütununda gösterir; her birini `ccam shortcut <account> --remove` ile kaldırın ya da bu bağlantıları `~/.local/bin` içinden kendiniz silin. Ardından komutu ve hesaplara artık ihtiyacınız yoksa veri klasörlerini silin:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows'ta, PowerShell'de:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam` klasörünü silmek, orada saklanan giriş bilgilerini ve konuşma geçmişini de siler. [İçe aktardığınız](/docs/commands) bir `~/.claude` klasörü olduğu yerde kalır.
