---
title: "ccam انسٹال کریں"
description: "ccam کو Linux یا macOS پر ایک کمانڈ سے، .deb پیکیج سے یا سورس سے انسٹال کریں۔ اس سے پہلے Claude Code انسٹال ہونا ضروری ہے۔"
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Linux یا macOS پر ٹرمینل میں یہ کمانڈ چلائیں:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

اسکرپٹ آپ کے سسٹم کے لیے ریلیز ڈاؤن لوڈ کرتا ہے، اس کا SHA-256 چیک سم جانچتا ہے اور `~/.local/bin` میں دو کمانڈز رکھتا ہے: `ccam` اور اس کا مختصر نام `ccs`۔ اگر یہ فولڈر آپ کے `PATH` میں شامل نہیں ہے تو اسکرپٹ وہ لائن پرنٹ کرتا ہے جو `~/.bashrc` یا `~/.zshrc` میں شامل کرنی ہے۔

ccam آپ کے لیے Claude Code شروع کرتا ہے، اس لیے پہلے [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) انسٹال کریں اور یقینی بنائیں کہ `claude --version` کام کرتا ہے۔

## معاون سسٹمز

- x86-64 اور ARM64 پر Linux: Ubuntu، Debian، Fedora، Arch، Alpine اور دیگر ڈسٹری بیوشنز۔ بائنری اسٹیٹک ہے اور اسے کسی اور چیز کی ضرورت نہیں۔
- Apple silicon اور Intel پر macOS۔
- Windows براہِ راست سپورٹ نہیں ہے۔ WSL کے اندر ccam کسی بھی Linux سسٹم کی طرح کام کرتا ہے۔

## انسٹال کرنے کے دوسرے طریقے

**Debian اور Ubuntu پیکیج۔** [ریلیزز کے صفحے](https://github.com/jhyunwoo/ccam/releases) سے `ccam_<version>_linux_amd64.deb` ڈاؤن لوڈ کریں اور اسے انسٹال کریں:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

پیکیج `/usr/bin/ccam` اور `/usr/bin/ccs` لنک انسٹال کرتا ہے۔

**سورس سے۔** Go 1.26 یا اس سے نئے ورژن کے ساتھ:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**کوئی مخصوص ورژن۔** `CCAM_VERSION` سیٹ کریں:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## شیل کمپلیشن

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

کمپلیشن `ccam` اور `ccs` دونوں کے لیے کام کرتی ہے، اکاؤنٹس کے ناموں سمیت۔

## انسٹالیشن چیک کریں

```sh
ccam version
ccam doctor
```

`ccam doctor` چیک کرتا ہے کہ Claude Code مل رہا ہے اور ہر اکاؤنٹ فولڈر، لنک اور لاگ اِن درست حالت میں ہے۔ اس کے بعد [فوری آغاز](/docs/quick-start) پر جائیں۔

## اَن انسٹال کرنا

دونوں کمانڈز اور اکاؤنٹس کے شارٹ کٹس ہٹائیں، پھر اگر آپ کو اکاؤنٹس کی مزید ضرورت نہیں تو ڈیٹا فولڈرز بھی ہٹا دیں:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` حذف کرنے سے وہاں محفوظ لاگ اِن اور گفتگو کی ہسٹری بھی حذف ہو جاتی ہے۔ جو `~/.claude` فولڈر آپ نے [امپورٹ کیا تھا](/docs/commands) وہ اپنی جگہ پر رہتا ہے۔
