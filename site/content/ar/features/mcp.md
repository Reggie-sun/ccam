---
title: خوادم MCP المشتركة
description: احتفظ بقائمة واحدة من خوادم MCP لكل مجموعة باستخدام ccam mcp add وadd-json وlist وremove وimport، وتُمرَّر عبر --mcp-config أو تُزامَن في كل حساب.
updated: 2026-10-23
order: 5
translation: machine
sourceHash: fefca7f8d857
---

يحتفظ Claude Code بخوادم MCP الخاصة بالمستخدم في `.claude.json`، وهو الملف نفسه الذي يحتوي على تسجيل دخول الحساب، لذلك لا يمكن مشاركة هذا الملف بين الحسابات. يحتفظ ccam بالخوادم المشتركة في ملف للمجموعة، `mcp.json`، ويسلّمه إلى كل حساب في المجموعة عند التشغيل. استخدم `ccam mcp` لإضافة خادم مرة واحدة لكل حساباتك، و`ccam mcp import` لنقل خادم أضفته داخل جلسة واحدة إلى الحسابات الأخرى.

## الاستخدام

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

يسري التغيير من التشغيل التالي لكل حساب؛ وتحتفظ الجلسات التي تعمل بالخوادم التي بدأت بها:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

تحت المجموعات، يذكر `ccam mcp list` كل خادم لا يصل إلى الحسابات الأخرى: `only in account work: pg` للخادم الذي أضافه حساب لنفسه، و`only in account side, for /home/me/code/app only: notes` للخادم المضاف لمجلد واحد (`claude mcp add` العادي، انظر أدناه)، و`account work has its own …, used in its sessions in place of the group's` عندما يحمل خادم الحساب الخاص اسم خادم من المجموعة، و`.mcp.json` الخاص بهذا المجلد مع الحسابات التي ينتظر موافقتها `waiting for approval in`، والخوادم التي تكون `in no account, only in ~/.claude.json (plain claude's)`، مع أمر [`ccam migrate --mcp`](/features/migrate) الذي ينسخها.

إضافة اسم تملكه المجموعة بالفعل تستبدل تعريفه، وتطبع `updated` بدلًا من `added`. ولا تصل الخوادم إلى حساب إلا إذا كان في تلك المجموعة وكانت المجموعة تشارك MCP (`ccam group set <group> --add mcp`)؛ وعندما لا تشاركه المجموعة، يخبرك `add` بذلك.

## الخيارات

| الخيار | ما يفعله | القيمة الافتراضية |
| --- | --- | --- |
| `-g`، `--group <group>` | المجموعة المستهدفة لـ `add` و`add-json` و`remove` | المجموعة الوحيدة التي يشارك عبرها حساب ما MCP، وإلا `default` |
| `-- <command> [args…]` (`add`) | يشغّل هذا الأمر بوصفه خادم stdio | |
| `--url <url>` (`add`) | عنوان URL لخادم بعيد (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | طريقة النقل لخادم `--url` | `http` |
| `--header <'Name: value'>` (`add`) | ترويسة HTTP، يمكن تكرارها | لا شيء |
| `-e`، `--env <KEY=VALUE>` (`add`) | متغير بيئة لخادم stdio، يمكن تكراره | لا شيء |
| `--force` (`import`) | يستبدل خادم المجموعة الذي يحمل الاسم نفسه ويختلف عنه؛ ويُحتفظ بنسخة احتياطية من `mcp.json` القديم | معطّل |
| `--keep-local` (`import`) | يُبقي المدخلات في ملف `.claude.json` الخاص بالحساب | معطّل |
| `-h`، `--help` | مساعدة كل أمر | |
| `--dry-run` (عام) | يدعمه `mcp list` فقط؛ وترفضه `add` و`add-json` و`remove` و`import` ولا تغيّر شيئًا | معطّل |

يحتاج `add` إما إلى أمر بعد `--` وإما إلى `--url`؛ وإلا يفشل مع `a command to run (after --) or --url is required`. ينطبق `--header` على خوادم `--url`، و`--env` على خوادم الأوامر. ويأخذ `add-json` صيغة JSON نفسها التي يأخذها `claude mcp add-json`، ويجب أن يكون التعريف كائنًا فيه `command` أو `url`. يستخدم `import` دائمًا المجموعة التي يشارك عبرها الحساب MCP، ويعرض `list` كل المجموعات، لذلك لا يغيّرهما `-g`.

## كيف يعمل

الملفات المعنية:

- `~/.local/share/ccam/groups/<group>/mcp.json`: الخوادم المشتركة، `{"mcpServers": {…}}`، وهي الصيغة التي يقرؤها `--mcp-config`. يُكتب باستخدام `fsutil.AtomicWrite` وبالوضع `0600`، لأن مدخلات الخوادم تحمل غالبًا مفاتيح API، ولا يُكتب إلا تحت القفل `<group>/.mcp.json.lock`، فإذا نُفّذ أمرا `ccam mcp add` في آن واحد احتفظ كلاهما بخادمه.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: ملف Claude Code نفسه الذي فيه خوادم الحساب. لا يقرأ ccam منه إلا القيمة `mcpServers` باستخدام الماسح `profile.TopLevelValue`، ولا يحلّل الملف كله أبدًا في مسار التشغيل.
- `~/.local/state/ccam/mcp/<name>/`: الملف `seen`، و`sync-stamp` (مع تجزئة لكل نسخة كتبها ccam، فيستطيع التمييز بين نسخة لم يغيّرها أحد ونسخة عدّلتها أنت)، و`sync-checked` (ما رأته استراتيجية المزامنة آخر مرة)، و`kept` (الخوادم التي احتفظ بها الحساب عندما غادر مشاركة MCP)، والنسخ الاحتياطية من `.claude.json`.

هناك استراتيجيتان، يحددهما `mcp_strategy` في `~/.config/ccam/config.toml`. والافتراضية هي استراتيجية الخيار:

<!-- motion:mcp-flag -->

1. يجد `launcher.applyMCP` مجموعة MCP الخاصة بالحساب. ويتحقق `mcp.FlagArgs` باستدعاء `stat` واحد من أن ملف `mcp.json` الخاص بالمجموعة موجود، دون تحليله.
2. يضع `Plan.InsertArgs` الوسيط `--mcp-config=<path>` مباشرة بعد `claude`، قبل وسائطك، فيظل `--mcp-config` أو `--strict-mcp-config` الخاص بك ساريًا فوق ذلك.
3. يتحقق `mcp.NewLocal` هل تغيّر ملف `.claude.json` الخاص بالحساب منذ التشغيل الأخير، ولا يحلّل خوادمه إلا عندئذ.
4. يستبدل `syscall.Exec` برنامج ccam بـ Claude Code.

يُمرَّر المسار وسيطًا واحدًا، `--mcp-config=<path>`، لا وسيطين. فالخيار `--mcp-config` في Claude Code يأخذ عدة قيم، لذلك سيقرأ `claude --mcp-config a.json "prompt"` الطلب على أنه ملف إعداد ثانٍ ويفشل مع `MCP config file not found`. أما صيغة الوسيط الواحد فتتوقف عند ذلك، ويعمل ما يليها من طلب أو أمر فرعي. ولا يحصل `ccam work mcp list` و`get` و`remove` على الخيار، لأن Claude Code يرفضه هناك؛ وبدلًا من ذلك يذكر ccam خوادم المجموعة في الخطأ القياسي ([التفاصيل](/features/launch#shared-mcp-servers)). يحمّل Claude Code خوادم `--mcp-config` إضافة إلى خوادم الحساب نفسه؛ وعند تعارض الأسماء يتقدم المدخل الموجود في `.claude.json`. ولا تكتب استراتيجية الخيار في `.claude.json` إلا من أجل Remote Control، وعندما يغادر حساب مشاركة MCP أو ينضم إليها من جديد ([أدناه](#leaving-and-joining))، ولـ `ccam mcp import`.

لا يأخذ `claude remote-control` (`rc`) الخيار `--mcp-config`، ولا تقرأ جلساته الخوادم إلا من `.claude.json`. ففي حالة `ccam work rc`، يحذف ccam هذا الخيار، وعلى الطرفية يسأل هل ينسخ خوادم المجموعة إلى `.claude.json`، مع نسخة احتياطية، كما تفعل استراتيجية `sync`. ويظل يسأل إلى أن توجد نسخة، ثم لا يسأل مرة أخرى إلا بعد تغيّر ملف `mcp.json` الخاص بالمجموعة أو اختفاء خادم منسوخ ([Remote Control](/features/launch#remote-control)). ولأن هذه النسخة تتقدم على `--mcp-config` عند تعارض الأسماء، تحدّث الأوامر `ccam mcp add` و`add-json` و`remove` و`import` النسخ التي تحتفظ بها الحسابات الأخرى في المجموعة:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

لا يظهر السطر `✓ Updated` إلا عندما تتغيّر نسخة. والنسخة التي أنشأها ccam أقدم تُحدَّث بالطريقة نفسها، وكذلك النسخة بعد `ccam migrate`.

الخادم الذي يُضاف داخل جلسة باستخدام `claude mcp add` لا يصل إلا إلى ملف `.claude.json` الخاص بذلك الحساب. ويستخدم `claude mcp add` العادي النطاق المحلي: فيُحمَّل الخادم لذلك الحساب وحده في ذلك المجلد وحده (ويجعله `-s user` خاصًا بالحساب في كل مجلد). ولا يستطيع ccam رؤيته في حينه، لأن Claude Code قد حلّ محله، لذلك يتحقق عند التشغيل التالي. يحتوي الملف `seen` على وقت تعديل `.claude.json` وأسماء الخوادم التي سبق الإبلاغ عنها. وعندما يختلف الوقت، يقارن ccam الأسماء بالمجموعة ويطبع تلميحًا لكل خادم جديد مرة واحدة، دون انتظار إدخال، مع الاستراتيجيتين كلتيهما:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

ويذكر التلميح أيضًا خادم الحساب الخاص الذي يحمل اسم خادم من المجموعة، لأن خادم الحساب يتقدم في جلساته.

مع `mcp_strategy = "sync"`، تظهر الخوادم في `claude mcp list` و`/mcp` كخوادم مستخدم عادية، ويكتبها ccam في ملف `.claude.json` الخاص بكل حساب قبل التشغيل:

<!-- motion:mcp-sync -->

1. يقارن `mcp.Sync` وقت تعديل ملف `mcp.json` الخاص بالمجموعة مع `sync-stamp`، وملف `.claude.json` الخاص بالحساب مع ما رآه آخر مرة. وإذا لم يتغيّر أي منهما، لا يحدث شيء آخر. وإلا يتحقق أيضًا من أن كل خادم من المجموعة لا يزال في `.claude.json`، ويعيد النسخة التي اختفت (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. يأخذ قفل Claude Code نفسه، أي المجلد `.claude.json.lock` بجانب الملف، الذي يُنشأ باستخدام `mkdir`. وينتظر حتى ثانيتين، ويعدّ القفل الذي عمره أكثر من 10 ثوانٍ قديمًا، كما تفعل مكتبة القفل في Claude Code.
3. يحفظ ملف `.claude.json` الحالي في `backups/` داخل مجلد الحالة، ويحتفظ بأحدث 5 نسخ.
4. يعيّن كل خادم من المجموعة في `mcpServers`، ويزيل الخوادم التي كتبها سابقًا ولم تعد في المجموعة، ويكتب الملف ذرّيًا بوضعه الأصلي، فقط عندما يتغيّر شيء. ولا يغيّر أو يزيل إلا النسخ التي كتبها ولم يغيّرها أحد منذ ذلك الحين. والخادم الذي يملكه الحساب بنفسه باسم خادم من المجموعة يبقى خادم الحساب، ويقول ccam ذلك مرة واحدة: `account work has its own MCP servers fs, so group default's of those names are not copied there`.
5. يسجّل الختم الجديد مع تجزئة لكل نسخة، ويبدأ Claude Code دون `--mcp-config`.

تُتخطّى المزامنة، مع `ccam: account work is already running, skipping MCP sync`، عندما تعمل جلسة للحساب، أيًّا كانت طريقة تشغيل Claude Code، لأن Claude Code يعيد كتابة `.claude.json` في أثناء عمله. وتطبع المزامنة الفاشلة `ccam: MCP sync failed (launching anyway): …` ويستمر التشغيل رغم ذلك.

### المغادرة والانضمام

الحساب الذي يتوقف عن مشاركة MCP (`ccam group leave` أو `ccam group set --remove mcp` أو `ccam group remove --leave`) يحتفظ بخوادم المجموعة: تنسخها الخطة إلى `.claude.json` الخاص به، مع نسخة احتياطية، وتبقى خوادمه الخاصة التي تحمل الأسماء نفسها. ومع `--fresh` لا يحصل على أي منها، وتسرد الخطة الخوادم التي سيتوقف تحميلها. وعندما ينضم مرة أخرى إلى مجموعة فيها تلك الخوادم، تُزال النسخ المحتفظ بها التي لم يغيّرها أحد، فتصل إليه تعديلات المجموعة اللاحقة.

ينقل `ccam mcp import` خوادم الحساب الخاصة إلى مجموعته:

<!-- motion:mcp-import -->

يقرأ ccam خوادم الحساب، الخاصة به والتي يملكها لمجلد واحد فقط، ويقارن كلًا منها بخادم المجموعة الذي يحمل الاسم نفسه باستخدام `mcp.SameDef`، الذي يتجاهل ترتيب المفاتيح والمسافات. ويصبح خادم المجلد خادمًا للمجموعة في كل مجلد (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`)؛ ويُرفض الاسم الواحد الذي له تعريفان مختلفان في مجلدين. ولا تُعرض خوادم المكوّنات الإضافية، لأن المكوّن الإضافي يأتي بها. يُنسخ الخادم الذي لا تملكه المجموعة أو تملكه بالتعريف نفسه. ويُتخطّى المختلف، لأن تعريف المجموعة هو ما تستخدمه كل الحسابات الأخرى:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

مع `--force` يُستبدل خادم المجموعة، ويُحفظ ملف `mcp.json` السابق أولًا في `.ccam-backup/<time>/mcp.json` داخل مجلد المجموعة. ثم يزيل `mcp.Remove` المدخلات المنسوخة من ملف `.claude.json` الخاص بالحساب، الخاصة به وتلك الخاصة بالمجلد في كتابة واحدة، تحت القفل نفسه وبعد نسخة احتياطية في مجلد الحالة. ويبقى ملف `mcp.json` الخاص بالمجموعة مقفلًا منذ تحميله حتى تختفي تلك المدخلات، فلا يزيل أحد الخادم من المجموعة في الأثناء. وهذا ضروري لأن مدخل الحساب سيتقدم على مدخل المجموعة عند تعارض الأسماء. ويُبقيها `--keep-local`.

## الأمان والحالات الخاصة

- **تسجيلات الدخول.** الخوادم التي تسجّل الدخول عبر OAuth تحتفظ برموزها في ملف `.credentials.json` الخاص بكل حساب (`mcpOAuth`)، لا في قائمة الخوادم، لذلك لا تؤدي مشاركة خادم إلى مشاركة تسجيل دخوله. سجّل الدخول مرة واحدة في كل حساب باستخدام `/mcp`. ولا ينسخ ccam تسجيلات الدخول أبدًا، ولا حتى [`ccam migrate`](/features/migrate): فرمز التحديث المنسوخ يُخرج كل حامل آخر. ويذكر [`ccam ext`](/features/ext) و`ccam doctor` الخوادم التي لم يسجّل حساب الدخول إليها بينما سجّل حساب آخر، ويقول `ccam next --continue` أيها ينقص الحساب الجديد. ولا تتلقى المجموعة أي ملف تسجيل دخول أبدًا.
- **الحسابات التي تعمل.** ينسخ `import` إلى المجموعة لكنه يترك المدخلات الخاصة بالحساب الذي يعمل ويخبرك بذلك؛ شغّله مرة أخرى بعد انتهاء الجلسة. وتتخطى المزامنة الحسابات التي تعمل. وكلاهما يحتسب أيضًا Claude Code الذي بدأ دون ccam.
- **النسخ الخاصة بـ Remote Control.** لا تُمَس نسخة الحساب الذي يعمل: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` ومع `mcp_strategy = "sync"` لا يحتاج شيء إلى تحديث، لأن كل تشغيل يُجري المزامنة.
- **الملفات التالفة.** إذا لم يكن `.claude.json` بصيغة JSON صالحة، تتركه المزامنة دون تغيير وتبلّغ عن الخطأ. وإذا تعذّر أخذ القفل خلال ثانيتين، يبلّغ ccam بـ `could not lock … (Claude Code is writing the file)`.
- **الحسابات التي لا تشارك MCP.** يرفض `import` ويطلب منك أولًا وضع الحساب في مجموعة تشارك MCP (`ccam group join <group> side`).
- **أسماء المجموعات.** لا يأخذ `-g` إلا اسم مجموعة؛ فهو يصبح مجلدًا تحت `groups/`، لذلك تُرفض أسماء مثل `../x`.

## صفحات ذات صلة

- [المشاركة بين الحسابات](/features/share)
- [استخدام Remote Control مع ccam](/guides/remote-control)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [نسخ إعداد موجود](/features/migrate)
- [كل الأوامر](/docs/commands)
