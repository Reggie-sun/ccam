---
title: خوادم MCP المشتركة
description: احتفظ بقائمة واحدة من خوادم MCP لكل مجموعة باستخدام ccam mcp add وadd-json وlist وremove وimport، وتُمرَّر عبر --mcp-config أو تُزامَن في كل حساب.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 46a7aa716f99
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
ccam mcp list                                                     # every group, and servers only one account has
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
```

إضافة اسم تملكه المجموعة بالفعل تستبدل تعريفه، وتطبع `updated` بدلًا من `added`. ولا تصل الخوادم إلى حساب إلا إذا كان يشارك MCP مع تلك المجموعة (`ccam share <name> --mcp on`)؛ وعندما لا يشارك أي حساب MCP بعد، يخبرك `add` بذلك.

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

- `~/.local/share/ccam/groups/<group>/mcp.json`: الخوادم المشتركة، `{"mcpServers": {…}}`، وهي الصيغة التي يقرؤها `--mcp-config`. يُكتب باستخدام `fsutil.AtomicWrite` وبالوضع `0600`، لأن مدخلات الخوادم تحمل غالبًا مفاتيح API.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: ملف Claude Code نفسه الذي فيه خوادم الحساب. لا يقرأ ccam منه إلا القيمة `mcpServers` باستخدام الماسح `profile.TopLevelValue`، ولا يحلّل الملف كله أبدًا في مسار التشغيل.
- `~/.local/state/ccam/mcp/<name>/`: الملف `seen`، و`sync-stamp`، والنسخ الاحتياطية من `.claude.json`.

هناك استراتيجيتان، يحددهما `mcp_strategy` في `~/.config/ccam/config.toml`. والافتراضية هي استراتيجية الخيار:

<!-- motion:mcp-flag -->

1. يجد `launcher.applyMCP` مجموعة MCP الخاصة بالحساب. ويتحقق `mcp.FlagArgs` باستدعاء `stat` واحد من أن ملف `mcp.json` الخاص بالمجموعة موجود، دون تحليله.
2. يضع `Plan.InsertArgs` الوسيط `--mcp-config=<path>` مباشرة بعد `claude`، قبل وسائطك، فيظل `--mcp-config` أو `--strict-mcp-config` الخاص بك ساريًا فوق ذلك.
3. يتحقق `mcp.NewLocal` هل تغيّر ملف `.claude.json` الخاص بالحساب منذ التشغيل الأخير، ولا يحلّل خوادمه إلا عندئذ.
4. يستبدل `syscall.Exec` برنامج ccam بـ Claude Code.

يُمرَّر المسار وسيطًا واحدًا، `--mcp-config=<path>`، لا وسيطين. فالخيار `--mcp-config` في Claude Code يأخذ عدة قيم، لذلك سيقرأ `claude --mcp-config a.json "prompt"` الطلب على أنه ملف إعداد ثانٍ ويفشل مع `MCP config file not found`. أما صيغة الوسيط الواحد فتتوقف عند ذلك، ويعمل ما يليها من طلب أو أمر فرعي مثل `mcp list`. يحمّل Claude Code خوادم `--mcp-config` إضافة إلى خوادم الحساب نفسه؛ وعند تعارض الأسماء يتقدم المدخل الموجود في `.claude.json`. ولا تكتب استراتيجية الخيار في `.claude.json` أبدًا.

الخادم الذي يُضاف داخل جلسة باستخدام `claude mcp add` لا يصل إلا إلى ملف `.claude.json` الخاص بذلك الحساب. ولا يستطيع ccam رؤيته في حينه، لأن Claude Code قد حلّ محله، لذلك يتحقق عند التشغيل التالي. يحتوي الملف `seen` على وقت تعديل `.claude.json` وأسماء الخوادم التي سبق الإبلاغ عنها. وعندما يختلف الوقت، يقارن ccam الأسماء بالمجموعة ويطبع تلميحًا لكل خادم جديد مرة واحدة، دون انتظار إدخال:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

مع `mcp_strategy = "sync"`، تظهر الخوادم في `claude mcp list` و`/mcp` كخوادم مستخدم عادية، ويكتبها ccam في ملف `.claude.json` الخاص بكل حساب قبل التشغيل:

<!-- motion:mcp-sync -->

1. يقارن `mcp.Sync` وقت تعديل ملف `mcp.json` الخاص بالمجموعة مع `sync-stamp`. وإذا لم يتغيّر، لا يُفتح ملف `.claude.json` الخاص بالحساب.
2. يأخذ قفل Claude Code نفسه، أي المجلد `.claude.json.lock` بجانب الملف، الذي يُنشأ باستخدام `mkdir`. وينتظر حتى ثانيتين، ويعدّ القفل الذي عمره أكثر من 10 ثوانٍ قديمًا، كما تفعل مكتبة القفل في Claude Code.
3. يحفظ ملف `.claude.json` الحالي في `backups/` داخل مجلد الحالة، ويحتفظ بأحدث 5 نسخ.
4. يعيّن كل خادم من المجموعة في `mcpServers`، ويزيل الخوادم التي كتبها سابقًا ولم تعد في المجموعة، ويحتفظ بخوادم الحساب الخاصة، ويكتب الملف ذرّيًا بوضعه الأصلي.
5. يسجّل الختم الجديد والأسماء التي يديرها، ويبدأ Claude Code دون `--mcp-config`.

تُتخطّى المزامنة، مع `ccam: account work is already running, skipping MCP sync`، عندما تعمل جلسة للحساب، لأن Claude Code يعيد كتابة `.claude.json` في أثناء عمله. وتطبع المزامنة الفاشلة `ccam: MCP sync failed (launching anyway): …` ويستمر التشغيل رغم ذلك. أما التلميح الخاص بالخوادم المضافة داخل جلسة فيخص استراتيجية الخيار وحدها.

ينقل `ccam mcp import` خوادم الحساب الخاصة إلى مجموعته:

<!-- motion:mcp-import -->

يقرأ ccam خوادم الحساب باستخدام `mcp.LocalServers` ويقارن كلًا منها بخادم المجموعة الذي يحمل الاسم نفسه باستخدام `mcp.SameDef`، الذي يتجاهل ترتيب المفاتيح والمسافات. يُنسخ الخادم الذي لا تملكه المجموعة أو تملكه بالتعريف نفسه. ويُتخطّى المختلف، لأن تعريف المجموعة هو ما تستخدمه كل الحسابات الأخرى:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

مع `--force` يُستبدل خادم المجموعة، ويُحفظ ملف `mcp.json` السابق أولًا في `.ccam-backup/<time>/mcp.json` داخل مجلد المجموعة. ثم يزيل `mcp.RemoveLocal` المدخلات المنسوخة من ملف `.claude.json` الخاص بالحساب، تحت القفل نفسه وبعد نسخة احتياطية في مجلد الحالة. وهذا ضروري لأن مدخل الحساب سيتقدم على مدخل المجموعة عند تعارض الأسماء. ويُبقيها `--keep-local`.

## الأمان والحالات الخاصة

- **تسجيلات الدخول.** الخوادم التي تسجّل الدخول عبر OAuth تحتفظ برموزها في ملف `.credentials.json` الخاص بكل حساب (`mcpOAuth`)، لا في قائمة الخوادم، لذلك لا تؤدي مشاركة خادم إلى مشاركة تسجيل دخوله. سجّل الدخول مرة واحدة في كل حساب باستخدام `/mcp`، أو انسخ تسجيلات الدخول من `~/.claude` باستخدام [`ccam migrate`](/features/migrate). ولا تتلقى المجموعة أي ملف تسجيل دخول أبدًا.
- **الحسابات التي تعمل.** ينسخ `import` إلى المجموعة لكنه يترك المدخلات الخاصة بالحساب الذي يعمل ويخبرك بذلك؛ شغّله مرة أخرى بعد انتهاء الجلسة. وتتخطى المزامنة الحسابات التي تعمل.
- **الملفات التالفة.** إذا لم يكن `.claude.json` بصيغة JSON صالحة، تتركه المزامنة دون تغيير وتبلّغ عن الخطأ. وإذا تعذّر أخذ القفل خلال ثانيتين، يبلّغ ccam بـ `could not lock … (Claude Code is writing the file)`.
- **الحسابات التي لا تشارك MCP.** يرفض `import` ويطلب منك تشغيل `ccam share side --mcp on` أولًا.
- **أسماء المجموعات.** لا يأخذ `-g` إلا اسم مجموعة؛ فهو يصبح مجلدًا تحت `groups/`، لذلك تُرفض أسماء مثل `../x`.

## صفحات ذات صلة

- [المشاركة بين الحسابات](/features/share)
- [مشاركة خوادم MCP والإعدادات والسجل](/docs/sharing)
- [نسخ إعداد موجود](/features/migrate)
- [كل الأوامر](/docs/commands)
