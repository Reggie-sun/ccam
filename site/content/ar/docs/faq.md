---
title: الأسئلة الشائعة
metaTitle: "أسئلة شائعة عن ccam: عدة حسابات Claude Code"
description: إجابات عن ccam وحسابات Claude Code وحدود الاستخدام والمشاركة والخصوصية والأنظمة المدعومة.
updated: 2026-10-24
order: 6
translation: machine
sourceHash: c4e3b89bb4ea
---

## هل ccam أداة رسمية من Anthropic؟

لا. ccam مشروع مفتوح المصدر غير رسمي، وليس تابعًا لشركة Anthropic ولا معتمدًا منها. يستخدم إعداد `CLAUDE_CONFIG_DIR` الذي [يوثّقه Claude Code](https://code.claude.com/docs/en/env-vars) لتشغيل عدة حسابات جنبًا إلى جنب.

## هل يُسمح لي باستخدام عدة حسابات Claude Code؟

يعتمد ذلك على شروط كل حساب، وعلى سياسة مؤسستك بالنسبة إلى حسابات العمل. راجعها قبل أن تستخدم عدة حسابات. كل ما يفعله ccam هو تسهيل التبديل بين حسابات لديك بالفعل.

**مقاعد Team وEnterprise.** يبدّل الاختيار التلقائي (`ccam next` و`--fallback` و`ccam pick` والمفتاح `n` في لوحة التحكم) أيضًا بين مقاعد المؤسسة نفسها في Claude Team أو Enterprise. يقول مركز المساعدة لدى Anthropic: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (يُخصَّص كل مقعد في خطة Team لشخص واحد محدد بالاسم، ولا يجوز أن يستخدمه إلا ذلك الشخص) و"Users can't route their work through another user's account to use that user's remaining usage" (لا يجوز للمستخدمين تمرير عملهم عبر حساب مستخدم آخر للاستفادة من الاستخدام المتبقي لذلك المستخدم) ([خطط Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129)؛ وبالمثل في [خطط Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). لذا فإن استخدام الاختيار التلقائي بين المقاعد قد يخالف شروط Anthropic. تحقّق من ذلك مع مؤسستك ومع Anthropic قبل أن تعتمد عليه. أما الخطط الشخصية (Pro وMax) فتتناولها الفقرة أعلاه.

## هل يتجاوز ccam حدود الاستخدام؟

لا. يحتفظ كل حساب بحدوده الخاصة. يعرضها ccam ويتيح لك المتابعة على حساب آخر تملكه عندما ينفد أحدها.

## كيف أستخدم حسابَي Claude Code على حاسوب واحد؟

امنح كل حساب مجلد إعدادات خاصًا به باستخدام `CLAUDE_CONFIG_DIR`. يفعل ccam ذلك نيابةً عنك: شغّل `ccam add` لكل حساب، ثم شغّل الحسابين باستخدام `ccam personal` و`ccam work`، أو بأوامر مختصرة تختار أسماءها بنفسك، مثل `ccwork`. ويوضح [الدليل](/guides/multiple-claude-code-accounts) أيضًا كيفية فعل ذلك يدويًا.

## كيف أتابع محادثة على حساب آخر؟

شغّل `ccam next --continue` في المجلد: فهو ينقل آخر محادثة في المجلد إلى الحساب التالي (`--auto` للحساب الذي لديه أكبر هامش). وداخل Claude Code، يفعل `/ccam:switch` من [المكوّن الإضافي لـ ccam](/docs/install#claude-code-plugin) الشيء نفسه للجلسة التي أنت فيها. وعندما يتشارك الحسابان سجل المحادثات (`ccam group set default --add history`، أو بمجموعة خاصة بهما: `ccam group add pair --items history --members work,side`)، يفتح الحساب الجديد المحادثة نفسها؛ وإلا ينسخها Claude Code إلى سجل الحساب الجديد.

## هل يرى ccam تسجيل دخولي أو يرسله؟

لا يشارك ccam تسجيل دخولك إلى Claude بين الحسابات ولا ينسخه أبدًا، ولا ينسخ تسجيلات دخول خوادم MCP أبدًا. يقرأ رمز الوصول الخاص بالحساب فقط ليطلب من نقطة نهاية الاستخدام لدى Anthropic بيانات استخدام ذلك الحساب، وذلك عندما تشغّل `ccam usage` أو `ccam next --auto` أو `ccam doctor`، ولا يجدّد الرمز المنتهي إلا باستخدام `ccam usage --refresh`؛ ومع `usage_source = "claude"` لا يقرأ الرمز إطلاقًا. ويقرأ `ccam ext` و`ccam doctor` أسماء تسجيلات دخول MCP الخاصة بالحساب، لا رموزها أبدًا. لا يحتوي ccam على أي قياس عن بُعد (telemetry). وبخلاف هذه الاستعلامات، فالطلب الوحيد الذي يرسله من تلقاء نفسه هو سؤال لوحة التحكم لـ GitHub عن أحدث إصدار من ccam، وهو لا يحمل أي شيء عنك أو عن حساباتك.

## كيف أسجّل الدخول على خادم بعيد؟

شغّل `ccam add` أو `ccam login <name>` في جلسة SSH كالمعتاد. ينتظر Claude Code تسجيل الدخول على الخادم، لذا ينتهي المتصفح على حاسوبك عند صفحة `http://localhost:<port>/callback?code=…` لا تُحمَّل. انسخ هذا العنوان من المتصفح والصقه في الطرفية: يمرّره ccam إلى Claude Code على الخادم، فيُكمل تسجيل الدخول. يحتوي العنوان على رمز يُستخدم مرة واحدة، فلا تشاركه.

## أين تُخزَّن حساباتي؟

في `~/.local/share/ccam/profiles/<name>/`، مجلد لكل حساب. على macOS، يحتفظ Claude Code بتسجيل الدخول نفسه في سلسلة المفاتيح، تحت اسم مشتق من ذلك المجلد، لذا يحتفظ `ccam rename` هناك بالاسم القديم للمجلد.

## ماذا يحدث لمجلد ~/.claude الموجود لدي؟

لا شيء، إلا إذا شغّلت `ccam import` أو `ccam migrate`. يسجّل الاستيراد المجلد في مكانه، فيظل `claude` العادي و`ccam <name>` يستخدمان الحساب والإعدادات نفسها. أما `ccam migrate` فينسخ إعداداته و`CLAUDE.md` والإضافات وخوادم MCP إلى الإعدادات المشتركة، ويكتفي بقراءة المجلد.

## ما الأنظمة التي يدعمها ccam؟

Linux (Ubuntu وDebian وFedora وArch وAlpine وغيرها) وmacOS وWindows 10 و11، أصليًا أو داخل WSL. على Windows، تحتاج مشاركة الملفات مثل `settings.json` إلى وضع المطور (Developer Mode)؛ أما المجلدات المشتركة فتعمل دونه، ودونه ينشئ أول `ccam add` المجموعة `default` بخوادم MCP والإضافات فقط ([المشاركة على Windows](/features/share#on-windows)).

## لماذا يظهر خادم MCP في حساب واحد فقط؟

الخوادم التي تُضاف داخل جلسة باستخدام `claude mcp add` تُخزَّن لذلك الحساب فقط، ودون `-s user` للمجلد الحالي فقط. يخبرك ccam بذلك في المرة التالية التي تشغّل فيها الحساب، ويعرضها `ccam mcp list` و`ccam ext`؛ انقل الخادم إلى المجموعة المشتركة باستخدام `ccam mcp import <account> <server>`.

## هل يبطئ ccam برنامج Claude Code؟

تشغيل Claude Code عبر ccam يضيف نحو 4 أجزاء من الألف من الثانية، وقد قيس ذلك باستخدام `scripts/perf-check.sh` على حاسوب محمول. يقرأ ccam ملفًا صغيرًا واحدًا، ويفحص الروابط المشتركة، ثم يستبدل نفسه بـ `claude`، لذلك لا يكون قيد التشغيل أثناء عملك. لا يستطيع Windows استبدال عملية، لذا يشغّل ccam هناك Claude Code وينتظره، دون أن يفعل أي شيء آخر.

## كيف أحدّث ccam؟

شغّل `ccam update`. ينزّل هذا الأمر أحدث إصدار، ويتحقق من مجموع التحقق الخاص به، ويستبدل الملف `ccam`؛ أما `ccam update --check` فيكتفي بإخبارك بوجود إصدار أحدث أو عدمه. وتخبرك لوحة التحكم أيضًا عند وجود إصدار أحدث، ويثبّته `U` بعد أن يسألك؛ ولا يحدّث ccam نفسه تلقائيًا أبدًا. إذا ثبّت حزمة `.deb`، فثبّت حزمة `.deb` الجديدة باستخدام apt بدلًا من ذلك. راجع [التحديث](/docs/install#update).

## كيف ألغي تثبيت ccam؟

أزل الأمر المختصر لكل حساب باستخدام `ccam shortcut <account> --remove`. وإذا كنت قد استوردت `~/.claude` وكان في مجموعة، فشغّل `ccam remove <name>` له أولًا، ليستعيد نسخة مما يشاركه. ثم احذف `ccam` من `~/.local/bin`، واحذف `~/.config/ccam` و`~/.local/share/ccam` و`~/.local/state/ccam` إذا لم تعد بحاجة إلى الحسابات. راجع [التثبيت](/docs/install#uninstall).
