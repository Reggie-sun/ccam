---
title: إدارة الحسابات
description: أضف حسابات Claude Code وسجّل الدخول إليها واسردها وأعد تسميتها وأزلها واستوردها عبر ccam add وccam login وccam list وccam rename وccam remove وccam import.
updated: 2026-10-23
order: 2
translation: machine
sourceHash: f4c09f92ce7d
---

كل حساب في ccam هو مجلد إعداد منفصل لبرنامج Claude Code له تسجيل دخوله الخاص. تنشئ هذه الأوامر الحسابات، وتسجّل الدخول إليها، وتسردها، وتعيد تسميتها، وتزيلها، وتسجّل المجلدات التي لديك بالفعل.

## الاستخدام

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

حيثما يأخذ أمر ما `<name|number>`، يمكنك تحديد اسم الحساب أو رقم خانته، أي العمود `#` في `ccam list`.

عند تشغيل `ccam add` في الطرفية دون خيارات، يفتح معالجًا. يسأل عن الاسم وعن طريقة تسجيل دخول الحساب (تسجيل دخول إلى Claude الآن، أو مفتاح API، أو تسجيل دخول إلى Claude لاحقًا)، ثم يعرض المجموعات التي يمكن الانضمام إليها مع تحديد `default` (وقبل وجود أي مجموعات، يعرض العناصر التي ستشاركها المجموعة `default`)، ثم يسأل `Create a shortcut command for work? [y/N]`، وينتهي بملخص لتأكيده. ويعود `Esc` خطوة إلى الوراء. أي من `--group` أو `--no-group` أو `--shortcut` أو `--no-login`، أو أنبوب بدلًا من طرفية، يتخطى المعالج؛ وعندها يلزم تحديد اسم، ولا يُنشأ أمر مختصر ما لم يحدد `--shortcut` واحدًا.

دون `--group`، ينضم الحساب الجديد إلى المجموعة `default`. وينشئها الحساب الأول، فتشارك خوادم MCP والإعدادات و`CLAUDE.md` والإضافات؛ راجع [المشاركة عبر المجموعات](/features/share). وعلى Windows دون وضع المطور (Developer Mode)، الذي لا يمكن ربط الملفات من دونه، ينشئ الحساب الأول `default` بالعنصرين `mcp` و`extensions` فقط ويخبرك بذلك (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`)، ويُرفض `--group` إلى مجموعة تشارك الإعدادات أو الذاكرة أو السجل.

يعرض `ccam list` الخانة والاسم والأمر المختصر والبريد الإلكتروني والمجموعات وآخر استخدام، مع `*` قبل الحساب الافتراضي. والبريد الإلكتروني هو الذي سُجّل دخول المجلد به الآن (من `.claude.json` الخاص به)؛ وعندما يختلف عن البريد الذي سجّله ccam، يتبعه البريد المسجَّل بصيغة `(was …)`:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

مع `--json`، يكون كل حساب كائنًا فيه `slot` و`shortcut` و`name` و`email` (لتسجيل الدخول الحالي، كما في الجدول) و`default` و`dir` و`last_used` (يُحذف إذا لم يُشغَّل الحساب قط) و`groups` و`share`. ويربط `share` كل عنصر (`mcp` و`settings` و`memory` و`extensions` و`history`) بالمجموعة التي يأتي منها، أو بالقيمة `"local"`.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `ccam add --no-login` | ينشئ الحساب دون تسجيل الدخول؛ سجّل الدخول لاحقًا باستخدام `ccam login <name>` | معطّل |
| `ccam add --group <group>` | ينضم إلى هذه المجموعة بدلًا من `default`. يمكن تكراره؛ وتُنشأ `default` إذا لم تكن موجودة | `default` |
| `ccam add --no-group` | لا ينضم إلى أي مجموعة: لا يشارك الحساب شيئًا | معطّل |
| `ccam add --shortcut <command>` | ينشئ أيضًا هذا الأمر المختصر في `~/.local/bin` | لا شيء |
| `ccam add --api --key-from <source>` | حساب يستخدم مفتاح API بدلًا من تسجيل الدخول؛ راجع [حسابات مفاتيح API](/features/api) | معطّل |
| `ccam list --json` | يطبع الحسابات بصيغة JSON | معطّل |
| `ccam remove -y`، `--yes` | يزيل دون سؤال | معطّل |
| `ccam logout -y`، `--yes` | يسجّل الخروج دون سؤال | معطّل |
| `ccam import --name <name>` | اسم الحساب للمجلد الذي تحدده؛ ويُرفض دون مجلد (`--name needs a folder: ccam import <folder> --name ops`) | مشتق من اسم المجلد |
| `ccam import -y`، `--yes` | يستورد كل مجلد يجده دون سؤال | معطّل |
| `--dry-run` | يدعمه `ccam list` و`ccam completion`، وهما لا يغيّران شيئًا. وترفضه `add` و`login` و`logout` و`use` و`rename` و`remove` و`import` مع `` --dry-run is not supported by `ccam add`; nothing was changed `` (مع ذكر اسم الأمر) | معطّل |

ليس لـ `ccam login` (عدا `--expired`) و`ccam use` و`ccam rename` خيارات خاصة بها. ويأخذ `ccam completion` القيمة `bash` أو `zsh` أو `powershell`؛ وأي شيء آخر يفشل مع `only bash, zsh and powershell are supported`.

أسماء الحسابات من 1 إلى 32 حرفًا من `a-z` و`0-9` و`-` و`_`، وتبدأ بحرف أو رقم. ولا يجوز أن يتكون الاسم من أرقام فقط، لأن `ccam 2` يعني الخانة 2. ولا يجوز أن يكون الاسم أمرًا من أوامر ccam أيضًا (`add` و`list` و`ls` و`run` و`use` و`next` و`usage` و`u` و`share` و`mcp` و`login` و`logout` و`rename` و`remove` و`rm` و`import` و`migrate` و`doctor` و`shims` و`shortcut` و`completion` و`help` و`version` و`statusline` و`update` و`folder` و`pick` و`ext` و`group` و`groups`)، فلا يكون `ccam <name>` ملتبسًا أبدًا. والقائمة هي `config.Reserved` في `internal/config`.

## كيف يعمل

يشغّل `ccam add work --shortcut ccwork` هذه الخطوات في `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. يتحقق `config.ValidateName` من الاسم، ويتحقق `shim.Check` من الأمر المختصر قبل إنشاء أي شيء.
2. ينشئ `profile.Add` المجلد `~/.local/share/ccam/profiles/work/` بالوضع `0700`، ويمنح الحساب أدنى رقم خانة متاح (`Config.NextSlot`)، فيُعاد استخدام رقم الحساب المُزال. ويصبح الحساب الأول هو الافتراضي.
3. ينضم الحساب إلى مجموعاته (`default` ما لم يحدد `--group` أو `--no-group` غير ذلك؛ وينشئ الحساب الأول `default`)، ويربط `applyShare` المجلد الجديد الفارغ بها باستخدام `share.Build` ثم `Apply`، مثلًا `settings.json` إلى `groups/default/settings.json`. وإذا فشل ذلك، يُزال المجلد الجديد، لأن بقاءه سيمنع كل محاولة لاحقة مع `folder already exists`.
4. يُحفظ الحساب في `~/.config/ccam/config.toml` عبر `fsutil.AtomicWrite`: ملف مؤقت، ثم `fsync`، ثم `rename`، فلا يترك الانهيار نصف إعداد أبدًا.
5. يشغّل `runLogin` الأمر `claude auth login` مع تعيين `CLAUDE_CONFIG_DIR` إلى مجلد الحساب، ومع `--email` بآخر بريد إلكتروني للحساب فتُملأ صفحة تسجيل الدخول مسبقًا. وإذا فشل تسجيل الدخول، يبقى الحساب، ويطبع ccam الخطأ ويطلب منك تشغيل `ccam login work` لاحقًا.
6. يأخذ `profile.ReadIdentity` القيمة `oauthAccount` من `profiles/work/.claude.json`: البريد الإلكتروني، ومعرّفي الحساب والمؤسسة، واسم المؤسسة، والخطة (`claude_max` و`claude_team` …). قد يبلغ حجم هذا الملف عدة ميغابايتات، لذلك يبحث `TopLevelValue` عن المفتاح الواحد بدلًا من تحليل الملف كله. ويحفظها ccam في `[profile.login]` داخل `config.toml`.
7. يربط `shim.Install` المسار `~/.local/bin/ccwork` بملف ccam التنفيذي، الذي يبحث عن `argv[0]` في `config.toml` عندما يبدأ، ويسجّل الرابط في `~/.local/state/ccam/shortcuts.json`. وإذا تعذّر إنشاء الرابط، يُضاف الحساب مع ذلك، ويخبرك ccam بالسبب (`Could not create shortcuts: …`) ويعرض الملخص `shortcut: ccwork (not created; see above)`.

قد يستغرق تسجيل الدخول دقائق، وقد يحفظ ccam آخر الإعداد في أثناء ذلك. لذلك تحفظ الأوامر التي تنتظر عبر `env.update`، الذي يعيد تحميل `config.toml`، ويطبّق التغيير الواحد (هنا البريد الإلكتروني) ثم يحفظ. ويعيد `ccam remove` التحميل بعد سؤاله للسبب نفسه.

يشغّل `ccam login` الخطوتين 5 و6 لحساب موجود. ولا يأخذ Claude Code البريد الإلكتروني إلا تلميحًا لصفحة تسجيل الدخول، لذلك قد ينتهي تسجيل الدخول في حساب Claude غير الذي كان عليه حساب ccam؛ فيقارن ccam تسجيل الدخول الجديد بـ `[profile.login]` ويخبرك بذلك (`ccam: this login is other@company.com, but work was dev@company.com before.`). ويحتفظ بتسجيل الدخول الجديد ويخبرك كيف تتراجع عنه: ``To undo it, run `ccam login work` again and sign in as dev@company.com.`` وتسجيل الدخول إلى حساب Claude آخر يُسقط أيضًا [سجل الحد](/features/usage#the-limit-record) الخاص بالحساب، لأن تلك الحدود كانت للحساب الآخر. وعندما يكون حساب ccam آخر مسجّل الدخول بحساب Claude نفسه، يحذّر ccam من أن الاثنين يتشاركان حدود الاستخدام. ويسجّل `ccam login --expired`، واحدًا تلو الآخر، دخول كل حساب لم يُسجَّل دخوله قط أو انتهى تسجيل دخوله: انتهت صلاحية رمز التحديث الخاص به، أو مسح Claude Code تسجيل الدخول بعد أن رفضه الخادم. وقبل أن يُفتح المتصفح الأول، يتحقق من إمكانية حفظ `config.toml`.

يشغّل `ccam logout` الأمر `claude auth logout` في الحساب ويمسح البريد الإلكتروني المحفوظ؛ ويبقى سجل تسجيل الدخول، فيُقارَن به تسجيل الدخول التالي. ويُبطل Claude Code تسجيل الدخول ويحذف ما يخزّنه للحساب، بما في ذلك تسجيلات دخول خوادم MCP. لذا يخبرك ccam أولًا بما سيذهب معه ويسألك (ويتخطى `--yes` السؤال): عندما يحمل تسجيل الدخول تسجيلات دخول MCP، وعندما يكون الحساب قيد التشغيل (فيُسجَّل خروج جلساته أيضًا)، وعندما يحمل مجلد حساب آخر تسجيل الدخول نفسه (مجلد منسوخ)، فيُبطله تسجيل الخروج كذلك:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

ولتسجيل الدخول بحساب Claude آخر، يكفي `ccam login <name>` وحده، وهو يحتفظ بتسجيلات دخول MCP. ولا يعيد `ccam use` إلا كتابة `default` في `config.toml`.

يضيف `ccam list` العمود `Organization` عندما تنتمي الحسابات إلى أكثر من مؤسسة واحدة (وتظهر الخطط الشخصية بالقيمة `personal`)، ويتضمن `ccam list --json` الحقلين `plan` و`org` لكل حساب.

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

يعمل `ccam rename` و`ccam remove` على مجلد الحساب وعلى ما يحتفظ به ccam باسم الحساب في أماكن أخرى. يسرد `nameState` هذه المسارات: `~/.local/state/ccam/last-used/<name>` و`~/.local/state/ccam/mcp/<name>` وملفات ذاكرة التخزين المؤقت للاستخدام `~/.local/state/ccam/usage/<name>.json` و`<name>.statusline.json`، وسجل الحد `<name>.limit.json`، وملاحظات الجلسات في `~/.local/state/ccam/handoff/<name>/`. ويتعامل `ccam remove work` معها بهذا الترتيب:

<!-- motion:accounts-remove -->

يرفض `ccam remove` أولًا إذا كان الحساب يعمل، ثم يسأل، ثم يتحقق مرة أخرى بعد الإجابة (فتُكتشف الجلسة التي بدأت أثناء انتظار السؤال). ويتحقق من إمكانية حفظ `config.toml` قبل أن يغيّر أي شيء. وإذا كان الحساب يشارك الإضافات، يعيد `RewritePluginPaths` كتابة `installed_plugins.json` و`known_marketplaces.json` الخاصين بالمجموعة، لأن Claude Code يخزّن مسارات مطلقة للإضافات قد تمر عبر مجلد هذا الحساب. ثم يُسجَّل خروج الحساب المُدار: يُشغَّل `claude auth logout` في مجلده، فلا ينتظر تسجيل دخول صالح في سلة المحذوفات (أو، على macOS، في عنصر سلسلة مفاتيح مسمّى باسم المجلد، قد يرثه حساب لاحق بذلك الاسم). ويُتخطّى ذلك لحساب مفتاح API وعندما يحمل مجلد حساب آخر تسجيل الدخول نفسه، الذي سيُبطله تسجيل الخروج أيضًا؛ ويُبلَّغ عن تسجيل الخروج الذي يفشل، وتستمر الإزالة. ينقل `profile.Trash` المجلد إلى `~/.local/share/ccam/trash/<name>-<time>/`، ويُحذف الحساب من `config.toml`، وتُحذف مسارات `nameState` فلا يرثها حساب لاحق بالاسم نفسه، ويُزال رابط الأمر المختصر: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`. وإذا تعذّر حفظ `config.toml` في النهاية، يعود المجلد من سلة المحذوفات. ولاستعادة حساب مُزال، أعد المجلد من سلة المحذوفات وسجّل الدخول مرة أخرى.

يعيد `ccam rename` تسمية المجلد، ويعيد كتابة مسارات الإضافات في `plugins/` الخاص به من المجلد القديم إلى الجديد، وينقل صفة الحساب الافتراضي معه، وينقل مسارات `nameState` إلى الاسم الجديد؛ وتُحدَّث ملاحظات الجلسات فيظل `ccam next --continue` يجدها. ويحتفظ الأمر المختصر باسمه. وإذا فشلت إعادة كتابة مسارات الإضافات، أو تعذّر حفظ `config.toml`، يُعاد المجلد إلى اسمه القديم. وعلى macOS يحتفظ المجلد باسمه القديم، لأن سلسلة المفاتيح تحفظ تسجيل الدخول وتسجيلات دخول MCP تحت مسار المجلد، والمجلد المنقول سيُسجَّل خروجه؛ والأمر نفسه على Windows عندما يحفظ Claude Code تسجيلات الدخول في Credential Manager. فيسجّل الإعداد حينها المجلد بوصفه `folder` بجانب الاسم الجديد، ويقول ccam `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

يسجّل `ccam import` مجلدًا في مكانه، بوصفه `dir` في `config.toml`. ولا ينقله أو ينسخه أبدًا: فرمز التحديث المنسوخ سيجعل النسخ تسجّل خروج بعضها بعضًا. ودون تحديد مجلد، يبحث عن مجلدات `~/.claude` و`~/.claude-*` غير المسجلة ويسألك عن كل منها؛ وتأتي الأسماء من المجلد (`~/.claude` يصبح `main`، و`~/.claude-work` يصبح `work`، و`~/.claude-2` يصبح `acct2`)، والمجلد الذي لا يمكن اشتقاق اسمه، أو الذي يحمل اسمًا لحساب موجود، يُسرد مع أمر `--name` الذي تستخدمه. والمجلد الذي يُعثر عليه ويفشل يُبلَّغ عنه بالرسالة `Skipped <path>: <reason>` وتُستورد المجلدات الأخرى. وتُجمع الإجابات أولًا ثم تُسجَّل المجلدات في `config.toml` بحالته عندئذٍ، فيبقى ما حفظته أوامر ccam الأخرى في أثناء ذلك؛ وتأتي أسطر `✓` بعد الحفظ. ودون `CLAUDE_CONFIG_DIR`، يحتفظ Claude Code بالملف `.claude.json` في مجلد المنزل بجانب `~/.claude`، لذلك عندما تستورد `~/.claude` وليس فيه ملف `.claude.json` خاص به، ينشئ ccam الرابط `~/.claude/.claude.json` إلى `~/.claude.json`. وعندها يرى `claude` العادي و`ccam main` الحساب نفسه. وإذا كان في `~/.claude` ملف `.claude.json` خاص به بجانب `~/.claude.json`، يترك ccam الاثنين ويخبرك بذلك: سيستخدم `ccam main` الملف الداخلي و`claude` العادي الملف الخارجي، فتبقى تفاصيل تسجيل الدخول وخوادم MCP لكل منهما منفصلة، لذا استبدل الملف الداخلي برابط للإبقاء على ملف واحد (ويُبلغ عنه [doctor](/features/doctor) أيضًا). وعندما يحمل المجلد تسجيل الدخول نفسه الذي يحمله حساب مسجَّل (مجلد منسوخ)، يخبرك ccam بذلك: أيهما يجدّده أولًا يُخرج الآخر، لذا شغّل `ccam login <name>` قبل استخدامه. ولا تنتمي الحسابات المستوردة إلى أي مجموعة حتى تضيفها إلى واحدة (`ccam group join <group> <name>`)، ولا يُنشأ لها أمر مختصر.

يطبع `ccam completion bash|zsh` سكربت إكمال تلقائي يُكمل أسماء الحسابات أيضًا.

## الأمان والحالات الخاصة

- يرفض `ccam rename` و`ccam remove` العمل ما دام الحساب يعمل (`refuseRunning`)، مع ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``، لأن الجلسة تظل تكتب في المجلد. ويشمل الفحص Claude Code أيًّا كانت طريقة تشغيله: عبر ccam، أو `claude` العادي، أو بيئة تطوير، أو الجلسات الخلفية وخدمتها الخلفية.
- يرفض `ccam rename` و`ccam remove` و`ccam logout` العمل ما دام ccam آخر يجدّد تسجيل الدخول ذاك أو يغيّره: ``another ccam is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- بينما يجري تغيير في المجموعات في ccam آخر، يرفض `ccam rename` و`ccam remove` مع `another ccam is changing the groups right now; try again when it is done`؛ وبعد تغيير توقف في منتصفه، مع ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.
- تشغيل `ccam remove` على حساب مستورد يحتفظ بالمجلد وتسجيل دخوله، لأجل `claude` العادي. وإذا كان الحساب في مجموعات، يغادرها أولًا كما يفعل `ccam group leave`: يصبح كل عنصر مشترك نسخة من نسخة المجموعة داخل المجلد، وتذهب خوادم MCP الخاصة بالمجموعة إلى `.claude.json` الخاص به. ثم يُلغى تسجيله: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. وتنتقل النسخ الاحتياطية من `.claude.json` الخاص به التي أنشأها ccam إلى `<folder>/.ccam-backup/mcp-<time>`؛ وإذا تعذّر نقلها، تبقى في `~/.local/state/ccam/mcp/<name>/backups` وتذكر الرسالة ذلك. ولا يُنقل `~/.claude` الخاص بك إلى سلة المحذوفات.
- لا يمكن أن يكون لحسابين المجلد نفسه: فملف `config.toml` الذي يذكر مجلدًا واحدًا مرتين يفشل تحميله.
- يرفض `ccam add` الاسم الذي يوجد مجلده بالفعل في `profiles/`، ويوجّهك إلى `ccam import`. وعلى macOS يتوقف أيضًا عندما تظل سلسلة المفاتيح تحمل تسجيل دخول أو تسجيلات دخول MCP لمسار المجلد الجديد، تركها ccam أقدم، ويذكر أمر `security delete-generic-password` الذي يزيلها. وإذا تعذّر حفظ `config.toml`، يُزال المجلد الجديد مرة أخرى.
- يرفض `ccam import` المجلد المسجل بالفعل، والمسار الذي ليس مجلدًا، والاسم المستخدم. وتُلاحَظ أيضًا نسخة من `~/.claude` الخاص بـ `claude` العادي عندما لا يكون `~/.claude` مسجَّلًا.
- يسأل `ccam logout` أولًا عندما يكون الحساب قيد التشغيل؛ أما `ccam login` فلا يتحقق.
- إذا لم يستطع الترحيل الوصول إلى Claude Code، مثلًا بعد انتهاء تسجيل الدخول، يطبع `ccam: could not pass the sign-in to Claude Code:` مع السبب. وخطأ HTTP يعني عادة أن الرمز انتهت صلاحيته؛ ابدأ تسجيل الدخول من جديد.

## صفحات ذات صلة

- [تشغيل Claude Code](/features/launch)
- [الأوامر المختصرة](/features/shortcuts)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [عدة حسابات Claude Code](/guides/multiple-claude-code-accounts)
