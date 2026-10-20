---
title: إدارة الحسابات
description: أضف حسابات Claude Code وسجّل الدخول إليها واسردها وأعد تسميتها وأزلها واستوردها عبر ccam add وccam login وccam list وccam rename وccam remove وccam import.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: 9eb992214132
---

كل حساب في ccam هو مجلد إعداد منفصل لبرنامج Claude Code له تسجيل دخوله الخاص. تنشئ هذه الأوامر الحسابات، وتسجّل الدخول إليها، وتسردها، وتعيد تسميتها، وتزيلها، وتسجّل المجلدات التي لديك بالفعل.

## الاستخدام

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

حيثما يأخذ أمر ما `<name|number>`، يمكنك تحديد اسم الحساب أو رقم خانته، أي العمود `#` في `ccam list`.

عند تشغيل `ccam add` في الطرفية دون خيارات، يفتح معالجًا. يسأل عن الاسم، ثم يعرض قائمة اختيار بالعناصر المشتركة (كل شيء محدد ما عدا سجل المحادثات)، ثم يسأل `Create a shortcut command for work? [y/N]`، وعند الإجابة بـ `y` يسأل `Shortcut name [ccwork]:`. أي من `--share` أو `--shortcut` أو `--no-login`، أو أنبوب بدلًا من طرفية، يتخطى المعالج؛ وعندها يلزم تحديد اسم، ولا يُنشأ أمر مختصر ما لم يحدد `--shortcut` واحدًا.

يعرض `ccam list` الخانة والاسم والأمر المختصر والبريد الإلكتروني والعناصر المشتركة وآخر استخدام، مع `*` قبل الحساب الافتراضي:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

مع `--json`، يكون كل حساب كائنًا فيه `slot` و`shortcut` و`name` و`email` و`default` و`dir` و`last_used` (يُحذف إذا لم يُشغَّل الحساب قط) و`share`. ويربط `share` كل عنصر (`mcp` و`settings` و`memory` و`extensions` و`history`) باسم مجموعته أو بالقيمة `"local"`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `ccam add --no-login` | ينشئ الحساب دون تسجيل الدخول؛ سجّل الدخول لاحقًا باستخدام `ccam login <name>` | معطّل |
| `ccam add --share <item=group>` | يشارك عنصرًا مع مجموعة، أو `item=local` لإبقائه خاصًا. يمكن تكراره | `mcp` و`settings` و`memory` و`extensions` في `default`؛ و`history` محلي |
| `ccam add --shortcut <command>` | ينشئ أيضًا هذا الأمر المختصر في `~/.local/bin` | لا شيء |
| `ccam add --api --key-from <source>` | حساب يستخدم مفتاح API بدلًا من تسجيل الدخول؛ راجع [حسابات مفاتيح API](/features/api) | معطّل |
| `ccam list --json` | يطبع الحسابات بصيغة JSON | معطّل |
| `ccam remove -y`، `--yes` | يزيل دون سؤال | معطّل |
| `ccam import --name <name>` | اسم الحساب للمجلد الذي تحدده | مشتق من اسم المجلد |
| `ccam import -y`، `--yes` | يستورد كل مجلد يجده دون سؤال | معطّل |
| `--dry-run` | يدعمه `ccam list` و`ccam completion`، وهما لا يغيّران شيئًا. وترفضه `add` و`login` و`logout` و`use` و`rename` و`remove` و`import` مع `` --dry-run is not supported by `ccam add`; nothing was changed `` (مع ذكر اسم الأمر) | معطّل |

ليس لـ `ccam login` و`ccam logout` و`ccam use` و`ccam rename` خيارات خاصة بها. ويأخذ `ccam completion` القيمة `bash` أو `zsh` أو `powershell`؛ وأي شيء آخر يفشل مع `only bash, zsh and powershell are supported`.

أسماء الحسابات من 1 إلى 32 حرفًا من `a-z` و`0-9` و`-` و`_`، وتبدأ بحرف أو رقم. ولا يجوز أن يتكون الاسم من أرقام فقط، لأن `ccam 2` يعني الخانة 2. ولا يجوز أن يكون الاسم أمرًا من أوامر ccam أيضًا (`add` و`list` و`ls` و`run` و`use` و`next` و`usage` و`u` و`share` و`mcp` و`login` و`logout` و`rename` و`remove` و`rm` و`import` و`migrate` و`doctor` و`shims` و`shortcut` و`completion` و`help` و`version` و`statusline` و`update` و`folder` و`pick` و`ext`)، فلا يكون `ccam <name>` ملتبسًا أبدًا. والقائمة هي `config.Reserved` في `internal/config`.

## كيف يعمل

يشغّل `ccam add work --shortcut ccwork` هذه الخطوات في `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. يتحقق `config.ValidateName` من الاسم، ويتحقق `shim.Check` من الأمر المختصر قبل إنشاء أي شيء.
2. ينشئ `profile.Add` المجلد `~/.local/share/ccam/profiles/work/` بالوضع `0700`، ويمنح الحساب أدنى رقم خانة متاح (`Config.NextSlot`)، فيُعاد استخدام رقم الحساب المُزال. ويصبح الحساب الأول هو الافتراضي.
3. يربط `applyShare` المجلد الجديد الفارغ بمجموعات المشاركة الخاصة به باستخدام `share.Build` ثم `Apply`، مثلًا `settings.json` إلى `groups/default/settings.json`. وإذا فشل ذلك، يُزال المجلد الجديد، لأن بقاءه سيمنع كل محاولة لاحقة مع `folder already exists`.
4. يُحفظ الحساب في `~/.config/ccam/config.toml` عبر `fsutil.AtomicWrite`: ملف مؤقت، ثم `fsync`، ثم `rename`، فلا يترك الانهيار نصف إعداد أبدًا.
5. يشغّل `runLogin` الأمر `claude auth login` مع تعيين `CLAUDE_CONFIG_DIR` إلى مجلد الحساب. وإذا فشل تسجيل الدخول، يبقى الحساب، ويطبع ccam الخطأ ويطلب منك تشغيل `ccam login work` لاحقًا.
6. يأخذ `profile.ReadEmail` القيمة `oauthAccount.emailAddress` من `profiles/work/.claude.json`. قد يبلغ حجم هذا الملف عدة ميغابايتات، لذلك يبحث `TopLevelValue` عن المفتاح الواحد بدلًا من تحليل الملف كله.
7. يربط `shim.Install` المسار `~/.local/bin/ccwork` بملف ccam التنفيذي، الذي يبحث عن `argv[0]` في `config.toml` عندما يبدأ.

قد يستغرق تسجيل الدخول دقائق، وقد يحفظ ccam آخر الإعداد في أثناء ذلك. لذلك تحفظ الأوامر التي تنتظر عبر `env.update`، الذي يعيد تحميل `config.toml`، ويطبّق التغيير الواحد (هنا البريد الإلكتروني) ثم يحفظ. ويعيد `ccam remove` التحميل بعد سؤاله للسبب نفسه.

يشغّل `ccam login` الخطوتين 5 و6 لحساب موجود. ويشغّل `ccam logout` الأمر `claude auth logout` في الحساب ويمسح البريد الإلكتروني المحفوظ. ولا يعيد `ccam use` إلا كتابة `default` في `config.toml`.

### تسجيل الدخول على جهاز بعيد

عبر SSH، قد يتعطل الجزء الخاص بالمتصفح في `claude auth login`. ينتظر Claude Code تسجيل الدخول على خادم استدعاء في `http://localhost:<port>/callback` على الخادم، لكن المتصفح الذي يُكمل تسجيل الدخول يعمل على حاسوبك، حيث لا يؤدي ذلك العنوان إلى أي مكان، فيعرض المتصفح خطأ مثل `ERR_CONNECTION_REFUSED`. ومع ذلك يظل العنوان يحمل الرمز الصالح لمرة واحدة، لذلك يستطيع ccam إيصاله:

<!-- motion:accounts-relay -->

لا يسلّم `runLogin` طرفيتك إلى Claude Code. بل يوصل الإدخال القياسي لـ Claude Code بأنبوب، وينسخ `relayLogin` إليه كل سطر تكتبه، فيصل الرمز الذي تلصقه عند مطالبة `Paste code here` في Claude Code كما هو. وعندما يكون السطر عنوان استدعاء، لا يقبله `callbackURL` إلا إذا كان المضيف `localhost` أو `127.0.0.1` أو `[::1]`، وكان المنفذ محددًا، والمسار `/callback`، وكلا `code` و`state` موجودين؛ ويمكن حذف `http://`. ثم يطلب `deliverCallback` ذلك العنوان على الخادم مع حد قدره 10 ثوانٍ ودون اتباع إعادة التوجيه، كما كان المتصفح سيفعل. يصل الطلب إلى خادم الاستدعاء الخاص بـ `claude auth login` المنتظر (الخط المتقطع في المخطط). يتلقى Claude Code الرمز ويحفظ تسجيل الدخول، وعلى Linux يكون ذلك في `.credentials.json` داخل مجلد الحساب. ولا يمكن لسطر ملصوق أن يجعل ccam يتصل بمضيف آخر أبدًا.

يطبع ccam هذا قبل أن يبدأ تسجيل الدخول:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

وعند النجاح يطبع `→ Passed the sign-in to Claude Code on this machine.`

وبينما ينتظر تسجيل الدخول إدخالًا في طرفية، يوقف ccam صدى الطرفية ووضع السطر، ويعرض `*` لكل حرف تكتبه، فلا يظهر على الشاشة رمز ملصوق أو عنوان استدعاء. ويتولى `readLine` تحرير السطر بنفسه: يرسل `Enter` السطر، ويمحو `Backspace` حرفًا واحدًا، وينهي `Ctrl+D` على سطر فارغ الإدخال. وتُستعاد الطرفية عند انتهاء تسجيل الدخول، وكذلك عند `Ctrl+C` أو `SIGTERM` أو `SIGHUP`، التي تظل تُنهي ccam. وإذا لم تكن هناك طرفية على الإدخال القياسي والخطأ القياسي، تُنسخ الأسطر كما تصل.

### إعادة التسمية والإزالة والاستيراد

يعمل `ccam rename` و`ccam remove` على مجلد الحساب وعلى ما يحتفظ به ccam باسم الحساب في أماكن أخرى. يسرد `nameState` هذه المسارات: `~/.local/state/ccam/last-used/<name>` و`~/.local/state/ccam/mcp/<name>` وملفات ذاكرة التخزين المؤقت للاستخدام `~/.local/state/ccam/usage/<name>.json` و`<name>.statusline.json`. ويتعامل `ccam remove work` معها بهذا الترتيب:

<!-- motion:accounts-remove -->

يرفض `ccam remove` أولًا إذا كان الحساب يعمل، ثم يسأل. وإذا كان الحساب يشارك الإضافات، يعيد `RewritePluginPaths` كتابة `installed_plugins.json` و`known_marketplaces.json` الخاصين بالمجموعة، لأن Claude Code يخزّن مسارات مطلقة للإضافات قد تمر عبر مجلد هذا الحساب. ينقل `profile.Trash` المجلد إلى `~/.local/share/ccam/trash/<name>-<time>/`، ويُحذف الحساب من `config.toml`، وتُحذف مسارات `nameState` فلا يرثها حساب لاحق بالاسم نفسه، ويُزال رابط الأمر المختصر. وللاستعادة، أعد المجلد من سلة المحذوفات.

يعيد `ccam rename` تسمية المجلد، ويعيد كتابة مسارات الإضافات في `plugins/` الخاص به من المجلد القديم إلى الجديد، وينقل صفة الحساب الافتراضي معه، وينقل مسارات `nameState` إلى الاسم الجديد. ويحتفظ الأمر المختصر باسمه. وإذا فشلت إعادة كتابة مسارات الإضافات، يُعاد المجلد إلى اسمه القديم، لأن الإعداد لا يُحفظ.

يسجّل `ccam import` مجلدًا في مكانه، بوصفه `dir` في `config.toml`. ولا ينقله أو ينسخه أبدًا: فرمز التحديث المنسوخ سيجعل النسخ تسجّل خروج بعضها بعضًا. ودون تحديد مجلد، يبحث عن مجلدات `~/.claude` و`~/.claude-*` غير المسجلة ويسألك عن كل منها؛ وتأتي الأسماء من المجلد (`~/.claude` يصبح `main`، و`~/.claude-work` يصبح `work`، و`~/.claude-2` يصبح `acct2`)، والمجلد الذي لا يمكن اشتقاق اسمه يُسرد مع أمر `--name` الذي تستخدمه. ودون `CLAUDE_CONFIG_DIR`، يحتفظ Claude Code بالملف `.claude.json` في مجلد المنزل بجانب `~/.claude`، لذلك عندما تستورد `~/.claude` وليس فيه ملف `.claude.json` خاص به، ينشئ ccam الرابط `~/.claude/.claude.json` إلى `~/.claude.json`. وعندها يرى `claude` العادي و`ccam main` الحساب نفسه. ولا تشارك الحسابات المستوردة أي شيء حتى تشغّل `ccam share`، ولا يُنشأ لها أمر مختصر.

يطبع `ccam completion bash|zsh` سكربت إكمال تلقائي يُكمل أسماء الحسابات أيضًا.

## الأمان والحالات الخاصة

- يرفض `ccam rename` و`ccam remove` العمل ما دام الحساب يعمل (`refuseRunning`)، مع `work is running; quit its Claude Code session first`، لأن الجلسة تظل تكتب في المجلد.
- تشغيل `ccam remove` على حساب مستورد يلغي تسجيله فقط: `✓ main unregistered (folder ... left in place)`. ولا يُنقل `~/.claude` الخاص بك إلى سلة المحذوفات.
- يرفض `ccam add` الاسم الذي يوجد مجلده بالفعل في `profiles/`، ويوجّهك إلى `ccam import`.
- يرفض `ccam import` المجلد المسجل بالفعل، والمسار الذي ليس مجلدًا، والاسم المستخدم.
- لا يتحقق `ccam logout` و`ccam login` هل الحساب يعمل.
- إذا لم يستطع الترحيل الوصول إلى Claude Code، مثلًا بعد انتهاء تسجيل الدخول، يطبع `ccam: could not pass the sign-in to Claude Code:` مع السبب. وخطأ HTTP يعني عادة أن الرمز انتهت صلاحيته؛ ابدأ تسجيل الدخول من جديد.

## صفحات ذات صلة

- [تشغيل Claude Code](/features/launch)
- [الأوامر المختصرة](/features/shortcuts)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [عدة حسابات Claude Code](/guides/multiple-claude-code-accounts)
