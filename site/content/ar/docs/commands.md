---
title: مرجع الأوامر
description: كل أوامر ccam وخياراتها، من ccam add وccam <name> إلى ccam usage وccam next وccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

ضع `--dry-run` قبل أي أمر لترى ما سيفعله دون أن يفعله، مثلًا `ccam --dry-run work`.

<!-- commands -->

## تشغيل Claude Code

يشغّل `ccam <name|number> [arguments]` برنامج Claude Code بذلك الحساب. لا يقرأ ccam الوسائط التي تأتي بعد الاسم؛ بل يمرّرها إلى `claude` بالترتيب نفسه، ويتصرف الإدخال القياسي ورمز الخروج والإشارات تمامًا كما مع `claude` نفسه. وضع `--` بعد الاسم اختياري.

أما `ccam` مع وسائط ولكن دون اسم حساب، مثل `ccam -p "summarise"`، فيستخدم الحساب الافتراضي. عيّنه باستخدام `ccam use <name>`.

الأوامر المختصرة `claude1` و`claude2` وهكذا هي روابط إلى ccam في `~/.local/bin`. وينشئ `ccam shims install --names` أيضًا `claude-<name>` لكل حساب.

## إضافة الحسابات وإزالتها

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

عند تشغيل `ccam import` دون تحديد مجلد، يبحث عن `~/.claude` و`~/.claude-*` ويسألك عن كل واحد منها. تُسجَّل المجلدات المستوردة في أماكنها؛ ولا ينقل ccam أي تسجيل دخول أو ينسخه أبدًا.

## المشاركة

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

يقبل كل عنصر القيمة `on` (المجموعة `default`) أو `off` أو اسم مجموعة. راجع [المشاركة](/docs/sharing).

## خوادم MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## الاستخدام والتبديل

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## الصيانة

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## الملفات

| المسار | المحتوى |
| --- | --- |
| `~/.config/ccam/config.toml` | الحسابات وما تشاركه |
| `~/.local/share/ccam/profiles/<name>/` | مجلد `CLAUDE_CONFIG_DIR` لكل حساب |
| `~/.local/share/ccam/groups/<group>/` | الإعدادات وخوادم MCP والإضافات والسجل المشتركة |
| `~/.local/state/ccam/` | ذاكرة التخزين المؤقت للاستخدام وأوقات آخر استخدام |
| `~/.local/bin/claude1`، `claude2`… | روابط إلى ccam |

تغيّر المتغيرات `XDG_CONFIG_HOME` و`XDG_DATA_HOME` و`XDG_STATE_HOME` مواقع هذه المجلدات.
