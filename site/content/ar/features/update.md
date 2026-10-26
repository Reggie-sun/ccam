---
title: تحديث ccam
description: يستبدل ccam update برنامج ccam بأحدث إصدار على GitHub بعد التحقق من SHA-256؛ ويكتفي --check بالإبلاغ، ويثبّت --version وسم إصدار محددًا.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: 598931e4801f
---

ينزّل `ccam update` إصدارًا من ccam من GitHub، ويتحقق منه مقابل ملف `checksums.txt` الخاص بالإصدار، ويستبدل الملف التنفيذي الذي يعمل بعملية إعادة تسمية واحدة. استخدمه إذا ثبّت ccam بسكربت التثبيت أو من أرشيف. وتسأل [لوحة التحكم](/features/dashboard) أيضًا GitHub عن أحدث إصدار في كل مرة تُفتح فيها، وتعرضه عبر `U` الذي ينفّذ الخطوات نفسها بعد أن يسألك. ولا تتصل الأوامر الأخرى بـ GitHub أبدًا، لذلك تعمل بالطريقة نفسها دون اتصال.

## الاستخدام

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

عندما لا يوجد إصدار أحدث، يطبع ccam الرسالة `ccam 0.4.0 is up to date.` ولا يغيّر شيئًا. وإذا حدّد `--version` الإصدار الذي لديك بالفعل، يطبع `ccam 0.4.0 is already installed.`

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `--check` | يكتفي بالإبلاغ عن وجود إصدار أحدث | معطّل |
| `--version <tag>` | يثبّت وسم الإصدار هذا، مثل `v0.3.0`، سواء كان أحدث من الحالي أو أقدم | الأحدث |
| `--dry-run` | مدعوم: يبحث عن أحدث وسم (ما لم يحدد `--version` وسمًا)، ويجري الفحوص نفسها على الملف التنفيذي، ويطبع مسارات التنزيل والمجموع الاختباري والاستبدال دون تنزيل | معطّل |

لا يمكن استخدام `--check` و`--version` معًا. يجب أن يكون الوسم على شكل `vMAJOR.MINOR.PATCH`، مع لاحقة اختيارية لإصدار تجريبي أو بناء؛ ويمكن حذف `v` في البداية. وأي شيء آخر يفشل، مثلًا `--version takes a release tag such as v0.3.0, not "1.2"`.

يستبدل متغير البيئة `CCAM_RELEASES_URL` العنوان `https://github.com/jhyunwoo/ccam/releases`، لاستخدام مرآة أو للاختبارات. ويجب أن يكون منظّمًا بالطريقة نفسها: يعيد `<url>/latest` التوجيه إلى `<url>/tag/<tag>`، والملفات موجودة في `<url>/download/<tag>/<asset>`.

## كيف يعمل

يتحقق ccam من التنزيل قبل أن يلمس الملف التنفيذي المثبّت، ويستبدل الملف التنفيذي بعملية إعادة تسمية واحدة. الشيفرة موجودة في `internal/update`، ولا يستورده إلا `ccam update`. ولا يستدعيه مسار التشغيل أبدًا.

1. يجد ccam الملف التنفيذي الذي يعمل باستخدام `os.Executable`، ويحلّ الروابط الرمزية باستخدام `filepath.EvalSymlinks` (`cli.updateCmd`). وهذا المسار الحقيقي هو ما يُستبدل.
2. دون `--version`، يطلب ccam العنوان `<releases>/latest` دون اتباع إعادة التوجيه (`Updater.Latest`)، ويقرأ الوسم من الترويسة `Location` (`…/tag/v0.5.0`). يستخدم هذا صفحة الإصدارات بدلًا من واجهة GitHub البرمجية، فلا ينطبق حد الواجهة البالغ 60 طلبًا في الساعة. وتنتهي مهلة البحث بعد 10 ثوانٍ.
3. يتحقق ccam من الوسم مقابل نمط صارم (`parseVersion`). فالوسم ينتهي في عناوين التنزيل وفي أمر صدفة مطبوع، لذلك لا يمر أي شيء فيه `/` أو `%` أو مسافات أو رموز الصدفة. ثم يقارنه ccam بالإصدار الذي يعمل (`compareVersions`)؛ ويأتي الإصدار التجريبي قبل إصداره النهائي في الترتيب.
4. يقرر ccam هل يجوز له أن يستبدل الملف التنفيذي أصلًا (`managed()` و`writable()`؛ انظر أدناه).
5. ينزّل ccam الملفين `ccam_<os>_<arch>.tar.gz` (وعلى Windows `ccam_windows_<arch>.zip`) و`checksums.txt` في الوقت نفسه، مع دقيقتين لكليهما (`Updater.Fetch`). ولكل تنزيل حد أقصى قدره 256 MB.
6. يجب أن يطابق مجموع SHA-256 للأرشيف سطره في `checksums.txt` قبل أن يحدث أي شيء آخر. والسطر المفقود أو المجموع الاختباري المختلف يوقف التحديث مع `nothing was changed`.
7. يقرأ ccam ملف tar ويأخذ أول ملف عادي اسمه `ccam` (`extract()`؛ ويأخذ `extractZip` الملف `ccam.exe` من ملف zip). وتُتجاهل المجلدات والروابط والمدخلات الأخرى، وللملف التنفيذي أيضًا حد أقصى قدره 256 MB.
8. يكتب ccam الملف التنفيذي الجديد باستخدام `fsutil.AtomicWrite` (`Updater.Install`): ملف مؤقت `.ccam.tmp*` بجانب الملف التنفيذي، بأذونات الملف القديم نفسها، يُكتب إلى القرص ثم يُعاد تسميته فوق القديم.

خط المعالجة كاملًا، من البحث عن الوسم إلى إعادة التسمية:

<!-- motion:update-pipeline -->

يُنزَّل الأرشيف و`checksums.txt` بالتوازي، ولا يُكتب شيء حتى يتطابق المجموع الاختباري. وإعادة التسمية ذرّية، لذلك يترك الانهيار أو امتلاء القرص إما الملف التنفيذي القديم وإما الجديد، ولا يترك نصف ملف أبدًا. والأوامر المختصرة مثل `~/.local/bin/ccwork` روابط رمزية إلى ملف ccam التنفيذي، فتشغّل الإصدار الجديد دون أي تغيير.

على Windows لا يمكن استبدال ملف برنامج قيد التشغيل أو حذفه، بل يمكن فقط إعادة تسميته. يكتب ccam الملف التنفيذي الجديد إلى `ccam.exe.new`، ويعيد تسمية `ccam.exe` قيد التشغيل إلى `ccam.exe.old`، ثم يعيد تسمية الملف الجديد إلى `ccam.exe` (`Updater.Install`). وإذا فشلت الخطوة الأخيرة، يُعاد الملف التنفيذي القديم إلى اسمه. ويحذف التحديث التالي ملفات `ccam.exe.old` التي يستطيع حذفها؛ وما دامت جلسة لا تزال تعمل من `ccam.exe.old` سابق، يصبح الملف التنفيذي قيد التشغيل `ccam.exe.old-<time>` بدلًا من ذلك. أما الأوامر المختصرة التي هي روابط صلبة بدلًا من روابط رمزية ([الأوامر المختصرة](/features/shortcuts)) فتُجدَّد بعد ذلك مباشرة، فتشغّل هي أيضًا الإصدار الجديد؛ ويشمل ذلك الروابط الصلبة التي أنشأها ccam لأي إصدار أقدم من ccam، لا لملف التنفيذ الذي استُبدل للتو فحسب.

### التثبيتات التي يديرها مدير حزم

نسخة ccam المثبّتة من حزمة `.deb` تخص apt، واستبدالها من وراء apt سيجعل قاعدة بيانات الحزم غير متوافقة مع الملف. قبل التنزيل، يتحقق ccam من مالك الملف التنفيذي:

<!-- motion:update-managed -->

1. الملف التنفيذي الموجود في `/usr/bin` أو `/bin` أو `/usr/sbin` أو `/sbin` يُعدّ مُدارًا بمدير حزم.
2. وفي غير ذلك، إذا كان `dpkg` مثبّتًا ونجح `dpkg -S <path>`، فيُعدّ أيضًا مُدارًا بمدير حزم.
3. وإلا ينشئ ccam ملفًا مؤقتًا `.ccam-update-*` بجانب الملف التنفيذي ثم يحذفه، ليتأكد من أنه يستطيع الكتابة هناك.

مع الملف التنفيذي المُدار بمدير حزم، يتوقف ccam ويطبع خطوات apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

على الأنظمة غير Linux، لا يُطبع إلا سطر سكربت التثبيت: `irm https://ccam.dev/install.ps1 | iex` على Windows. وعندما لا يكون المجلد قابلًا للكتابة، يقترح ccam إعادة تشغيل سكربت التثبيت، الذي يثبّت ccam في `~/.local/bin`.

## الأمان والحالات الخاصة

- بناء التطوير يبلّغ عن إصداره بالقيمة `dev` وليس له إصدار يُقارن به، لذلك يرفضه `ccam update` و`ccam update --check`. ويثبّت `ccam update --version vX.Y.Z` إصدارًا فوقه. وللبناءات المصنوعة باستخدام `go install`، راجع [تثبيت ccam](/docs/install#other-ways-to-install).
- التنزيل الفاشل، أو غياب ملف لنظامك، أو عدم تطابق المجموع الاختباري، يترك الملف التنفيذي المثبّت كما كان.
- الإصدار غير الموجود يفشل مع `there is no release <tag>`. وعندما يرفض GitHub الطلب برمز HTTP 403 أو 429، يقول ccam إن السبب على الأرجح حد للطلبات ويطلب منك المحاولة لاحقًا.
- يقرأ `--check` فقط. ولا يتحقق هل الملف التنفيذي مُدار بمدير حزم أو قابل للكتابة.
- استخدام `--version` مع وسم أقدم هو طريقة الرجوع إلى إصدار سابق.

## صفحات ذات صلة

- [تثبيت ccam](/docs/install)
- [مرجع الأوامر](/docs/commands)
- [الأسئلة الشائعة](/docs/faq)
