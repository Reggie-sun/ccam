---
title: تثبيت ccam
metaTitle: تثبيت ccam على Linux وmacOS وWindows
description: ثبّت ccam على Linux أو macOS أو Windows بأمر واحد أو من حزمة .deb أو عبر go install، وحدّثه باستخدام ccam update. يجب تثبيت Claude Code أولًا.
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

يشغّل ccam برنامج Claude Code نيابةً عنك، لذا ثبّت [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) أولًا وتأكد من أن `claude --version` يعمل.

شغّل هذا الأمر في طرفية على Linux أو macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

ينزّل السكربت الإصدار المناسب لنظامك، ويتحقق من مجموع SHA-256 الخاص به، ويضع `ccam` في `~/.local/bin`، حيث ينشئ ccam أيضًا الأوامر المختصرة التي تختارها لاحقًا، مثل `ccwork`. إذا لم يكن ذلك المجلد ضمن `PATH` لديك، يطبع السكربت السطر الذي يجب إضافته إلى `~/.bashrc` أو `~/.zshrc`.

على Windows، شغّل هذا في PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

ينزّل الملف `ccam_windows_<arch>.zip`، ويتحقق من مجموع SHA-256 الخاص به، ويضع `ccam.exe` في `%USERPROFILE%\.local\bin`، وهو المجلد الذي يضع فيه مثبّت Claude Code الملف `claude.exe`. وإذا لم يكن ذلك المجلد ضمن `PATH` الخاص بمستخدمك بعد، يضيفه السكربت ويُبقي مدخلات مثل `%USERPROFILE%\bin` كما هي. افتح طرفية جديدة بعد ذلك.

## الأنظمة المدعومة

- Linux على x86-64 وARM64: Ubuntu وDebian وFedora وArch وAlpine وتوزيعات أخرى. الملف التنفيذي ثابت (static) وليست له أي تبعيات.
- macOS على معالجات Apple silicon وIntel.
- Windows 10 و11 على x86-64 وARM64، في PowerShell أو Windows Terminal أو Git Bash:
  - تعمل المجلدات المشتركة دون وضع المطور (Developer Mode).
  - تحتاج الملفات المشتركة (`settings.json` و`CLAUDE.md` و`history.jsonl`) إلى وضع المطور (في الإعدادات > النظام > للمطورين)، لأن Windows لا يسمح بإنشاء روابط إلى الملفات إلا في هذا الوضع أو للمسؤول؛ راجع [المشاركة على Windows](/features/share#on-windows).
  - دون وضع المطور، ينشئ أول `ccam add` المجموعة `default` بخوادم MCP والإضافات فقط، ويخبرك كيف تضيف البقية بعد تشغيل وضع المطور.
- داخل WSL، يعمل ccam كما يعمل على أي نظام Linux آخر.

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

إذا ثبّت من وسم إصدار (`@latest` أو `@vX.Y.Z`)، يُبلغ ccam عن ذلك الإصدار، ويعمل معه `ccam update` وعرض التحديث في لوحة التحكم. أما البناء من `@main`، أو من نسخة محلية فيها إيداعات (commits) بعد وسم، فيُبلغ عن `dev`، ويرفضه `ccam update` ما لم تمرّر `--version`.

**إصدار محدد.** عيّن `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## الإكمال التلقائي في الصدفة

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

يغطي الإكمال التلقائي كل أوامر `ccam` وخياراته، بما في ذلك أسماء الحسابات.

## المكوّن الإضافي لـ Claude Code

يعلّم المكوّن الإضافي لـ ccam برنامج Claude Code نفسه أن يستخدم ccam: أن يتحقق من الحساب الذي لديه هامش، وأن يشغّل عمّال `claude -p` متوازيين على حسابات مختلفة باستخدام `ccam pick`، وأن يشرح غياب مهارة أو خادم MCP باستخدام `ccam ext`. ثبّته من داخل Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

يضيف:

- مهارة تقترح الأوامر التي تسجّل الدخول أو تنقل الملفات لكنها لا تشغّلها أبدًا؛
- الأوامر `/ccam:usage` و`/ccam:ext` و`/ccam:switch`؛
- خطافين ([التفاصيل](/features/next#inside-claude-code)). يعرض أحدهما إشعارًا في الطرفية عندما تتوقف جلسة عند حد الاستخدام؛ وعلى Linux وmacOS، مع تشغيل [`auto_switch`](/features/next#auto-switch)، ينقل ccam بعدها الجلسة إلى حساب آخر لديه متسع. ويتيح الآخر لـ `/ccam:switch [account]` أن يسلّم المحادثة إلى حساب آخر.

تشغّل الخطافات `ccam`، لذا يجب أن يكون ضمن `PATH` الذي يراه Claude Code.

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

تبحث لوحة التحكم (`ccam` دون وسائط، في الطرفية) أيضًا عن إصدار جديد في كل مرة تُفتح فيها، وإذا وُجد إصدار عرضته عليك: اضغط `U` ثم أكّد. لا يبحث أي أمر آخر عن تحديثات، ولا يحدّث ccam نفسه أبدًا دون أن يسألك. لتثبيت إصدار محدد، حتى لو كان أقدم، مرّر الخيار `--version`. ويعرض `ccam --dry-run update` ما سيُنزَّل وما سيُستبدَل دون تغيير أي شيء.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

تشغيل سكربت التثبيت مرة أخرى يحدّث أيضًا النسخة الموجودة في `~/.local/bin`. وإذا بنيت ccam باستخدام `go install` من وسم إصدار، فيمكنك تشغيل `ccam update` أو أمر `go install` نفسه مرة أخرى.

**Windows.** ينزّل `ccam update` ملف `.zip` الخاص بـ Windows. لا يمكن استبدال برنامج قيد التشغيل على Windows، لذا يُعاد تسمية الملف التنفيذي القديم إلى `ccam.exe.old`، ويحذفه التحديث التالي؛ وما دامت جلسة لا تزال تعمل من `ccam.exe.old` سابق، يكون الاسم `ccam.exe.old-<time>`. وتُجدَّد الأوامر المختصرة التي هي روابط صلبة بعد ذلك مباشرة. وتشغيل `install.ps1` مرة أخرى يعمل بالطريقة نفسها، حتى أثناء عمل جلسة Claude Code من الملف التنفيذي القديم، ويطبع `OK 2 shortcuts now run the new ccam` عندما يجدّد الأوامر المختصرة التي هي روابط صلبة.

**حزمة Debian وUbuntu.** لا يستبدل `ccam update` الملف `/usr/bin/ccam` الذي ثبّته apt، بل يطبع هذه الخطوات بدلًا من ذلك. نزّل الملف الجديد `ccam_<version>_linux_amd64.deb` من [صفحة الإصدارات](https://github.com/jhyunwoo/ccam/releases) وثبّته فوق القديم:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## إلغاء التثبيت

يزيل `ccam uninstall` برنامج ccam وكل ما أنشأه. ويعرض أولًا ما سيُزال، ثم يسأل:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- يُزال كل حساب كما يفعل `ccam remove`. فالحساب الموجود في مجلدات ccam يُسجَّل خروجه أولًا ثم يُحذف مع مجلده، ومعه سجل المحادثات والإعدادات أيضًا.
- أما مجلد `~/.claude` الذي [استوردته](/docs/commands) فيبقى في مكانه ليعمل مع `claude` العادي. وإذا كان في مجموعة، يستعيد نسخة من كل ما كان يشاركه. ويحلّ محلَّ سطر الحالة الخاص بـ ccam في ملف `settings.json` الخاص به سطرُ الحالة الذي كان ccam يغلّفه، إن وُجد. وإذا كان المكوّن الإضافي ccam مثبّتًا فيه، يطبع ccam الأمر `/plugin uninstall` لتشغّله في Claude Code.
- تُحذف الأوامر المختصرة و`~/.config/ccam` و`~/.local/share/ccam` (مع سلة محذوفات ccam) و`~/.local/state/ccam`، ثم الأمر `ccam` نفسه.

يرفض ccam العمل ما دام أحد الحسابات قيد التشغيل؛ فأغلق جلسات Claude Code الخاصة به أولًا. أما `/usr/bin/ccam` الذي ثبّته apt فيبقى؛ أزله بعد ذلك باستخدام `sudo apt remove ccam`. وعلى Windows لا يستطيع البرنامج قيد التشغيل حذف نفسه، لذا ينتهي `ccam uninstall` بطباعة الأمر `Remove-Item` الخاص بـ `ccam.exe`.
