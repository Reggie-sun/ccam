---
title: ccam kurulumu
description: ccam'ı Linux veya macOS üzerine tek komutla, bir .deb paketinden ya da kaynak koddan kurun. Önce Claude Code kurulu olmalıdır.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Linux veya macOS üzerinde bir terminalde şunu çalıştırın:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Betik, sisteminize uygun sürümü indirir, SHA-256 sağlama toplamını doğrular ve `~/.local/bin` klasörüne iki komut yerleştirir: `ccam` ve onun kısa takma adı `ccs`. Bu klasör `PATH` içinde değilse betik, `~/.bashrc` veya `~/.zshrc` dosyasına eklemeniz gereken satırı yazdırır.

ccam, Claude Code'u sizin için başlatır; bu nedenle önce [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) kurun ve `claude --version` komutunun çalıştığından emin olun.

## Desteklenen sistemler

- x86-64 ve ARM64 üzerinde Linux: Ubuntu, Debian, Fedora, Arch, Alpine ve diğer dağıtımlar. İkili dosya statiktir ve hiçbir bağımlılığı yoktur.
- Apple silicon ve Intel üzerinde macOS.
- Windows doğrudan desteklenmez. WSL içinde ccam, herhangi bir Linux sisteminde olduğu gibi çalışır.

## Diğer kurulum yolları

**Debian ve Ubuntu paketi.** `ccam_<version>_linux_amd64.deb` dosyasını [sürümler sayfasından](https://github.com/jhyunwoo/ccam/releases) indirin ve kurun:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Paket, `/usr/bin/ccam` dosyasını ve `/usr/bin/ccs` bağlantısını kurar.

**Kaynak koddan.** Go 1.26 veya daha yeni bir sürümle:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**Belirli bir sürüm.** `CCAM_VERSION` değişkenini ayarlayın:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Kabuk tamamlama

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Tamamlama, hesap adları da dahil olmak üzere hem `ccam` hem de `ccs` için çalışır.

## Kurulumu kontrol edin

```sh
ccam version
ccam doctor
```

`ccam doctor`, Claude Code'un bulunduğunu ve her hesap klasörünün, bağlantının ve girişin yerinde olduğunu kontrol eder. [Hızlı başlangıç](/docs/quick-start) ile devam edin.

## Kaldırma

İki komutu ve hesap kısayollarını, ardından hesaplara artık ihtiyacınız yoksa veri klasörlerini silin:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` klasörünü silmek, orada saklanan giriş bilgilerini ve konuşma geçmişini de siler. [İçe aktardığınız](/docs/commands) bir `~/.claude` klasörü olduğu yerde kalır.
