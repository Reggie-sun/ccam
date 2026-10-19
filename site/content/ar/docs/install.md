---
title: تثبيت ccam
description: ثبّت ccam على Linux أو macOS بأمر واحد، أو من حزمة .deb، أو من الشيفرة المصدرية، وأبقِه محدّثًا باستخدام ccam update. يجب تثبيت Claude Code أولًا.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: bde486b82b2e
---

شغّل هذا الأمر في طرفية على Linux أو macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

ينزّل السكربت الإصدار المناسب لنظامك، ويتحقق من مجموع SHA-256 الخاص به، ويضع `ccam` في `~/.local/bin`، حيث ينشئ ccam أيضًا الأوامر المختصرة التي تختارها لاحقًا، مثل `ccwork`. إذا لم يكن ذلك المجلد ضمن `PATH` لديك، يطبع السكربت السطر الذي يجب إضافته إلى `~/.bashrc` أو `~/.zshrc`.

يشغّل ccam برنامج Claude Code نيابةً عنك، لذا ثبّت [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) أولًا وتأكد من أن `claude --version` يعمل.

## الأنظمة المدعومة

- Linux على x86-64 وARM64: Ubuntu وDebian وFedora وArch وAlpine وتوزيعات أخرى. الملف التنفيذي ثابت (static) وليست له أي تبعيات.
- macOS على معالجات Apple silicon وIntel.
- Windows غير مدعوم مباشرة. داخل WSL، يعمل ccam كما يعمل على أي نظام Linux.

## طرق تثبيت أخرى

**حزمة Debian وUbuntu.** نزّل `ccam_<version>_linux_amd64.deb` من [صفحة الإصدارات](https://github.com/jhyunwoo/ccam/releases) وثبّتها:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

تثبّت الحزمة `/usr/bin/ccam`.

**من الشيفرة المصدرية.** باستخدام Go 1.26 أو أحدث:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**إصدار محدد.** عيّن `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## الإكمال التلقائي في الصدفة

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

يغطي الإكمال التلقائي كل أوامر `ccam` وخياراته، بما في ذلك أسماء الحسابات.

## التحقق من التثبيت

```sh
ccam version
ccam doctor
```

يتحقق `ccam doctor` من العثور على Claude Code ومن سلامة كل مجلد حساب ورابط وتسجيل دخول. تابع مع [البدء السريع](/docs/quick-start).

## التحديث

```sh
ccam update --check
ccam update
```

يكتفي `ccam update --check` بإخبارك بوجود إصدار أحدث أو عدمه. أما `ccam update` فينزّل الإصدار المناسب لنظامك، ويتحقق من مجموع SHA-256 الخاص به بمطابقته مع ملف `checksums.txt` المرفق بالإصدار، ثم يستبدل الملف `ccam` في خطوة واحدة. إذا لم يتطابق المجموع، فلن يتغير أي شيء. الأوامر المختصرة مثل `ccwork` هي روابط إلى ccam، لذا تشغّل الإصدار الجديد فورًا.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

لا يبحث ccam عن التحديثات من تلقاء نفسه أبدًا، بل يبحث عن إصدار جديد فقط عندما تشغّل هذا الأمر. لتثبيت إصدار محدد، حتى لو كان أقدم، مرّر الخيار `--version`. ويعرض `ccam --dry-run update` ما سيُنزَّل وما سيُستبدَل دون تغيير أي شيء.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

تشغيل سكربت التثبيت مرة أخرى يحدّث أيضًا النسخة الموجودة في `~/.local/bin`. وإذا بنيت ccam باستخدام `go install`، فشغّل الأمر نفسه مرة أخرى.

**حزمة Debian وUbuntu.** لا يستبدل `ccam update` الملف `/usr/bin/ccam` الذي ثبّته apt، بل يطبع هذه الخطوات بدلًا من ذلك. نزّل الملف الجديد `ccam_<version>_linux_amd64.deb` من [صفحة الإصدارات](https://github.com/jhyunwoo/ccam/releases) وثبّته فوق القديم:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## إلغاء التثبيت

أزل الأوامر المختصرة للحسابات أولًا، ما دام ccam مثبتًا. يعرضها `ccam list` في عمود Shortcut؛ أزل كلًّا منها باستخدام `ccam shortcut <account> --remove`، أو احذف تلك الروابط من `~/.local/bin` بنفسك. ثم أزل الأمر، ومجلدات البيانات إذا لم تعد بحاجة إلى الحسابات:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

حذف `~/.local/share/ccam` يحذف أيضًا تسجيلات الدخول وسجل المحادثات المخزنة فيه. أما مجلد `~/.claude` الذي [استوردته](/docs/commands) فيبقى في مكانه.
