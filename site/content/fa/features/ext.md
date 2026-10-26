---
title: بررسی افزونه‌ها
description: ccam ext مهارت‌ها، فرمان‌ها، پلاگین‌ها، سرورهای MCP و هوک‌هایی را که هر حساب بارگذاری می‌کند، منبع هر کدام و موارد خراب را نشان می‌دهد.
updated: 2026-10-26
order: 14
translation: machine
sourceHash: 7ca357a7f38a
---

`ccam ext` به این پرسش پاسخ می‌دهد: «چرا این مهارت (یا سرور MCP) در یک حساب کار می‌کند اما در دیگری نه؟» برای هر حساب، افزونه‌هایی را که Claude Code بارگذاری می‌کند و منبع هر کدام را فهرست می‌کند: یک گروه اشتراک، پوشه خود حساب، یا پروژه‌ای که در آن هستید. زیر جدول، چیزهایی را فهرست می‌کند که مانع بارگذاری یک افزونه به شکلی که انتظار دارید می‌شوند. چیزی را تغییر نمی‌دهد.

## نحوه استفاده

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

هر خانه منبع آن مورد در آن حساب را نشان می‌دهد: گروه اشتراک (`default`)، `account` وقتی در پوشه خود حساب است، `local` برای سرور MCP که حساب فقط برای این پوشه دارد، یا `-` وقتی حساب آن را ندارد. پلاگینی که نصب شده اما روشن نیست، پس از منبعش `(off)` نشان می‌دهد. `!` موردی را که مشکلی دارد علامت می‌زند که زیر جدول توضیح داده می‌شود؛ مشکلاتی که به هیچ حسابی تعلق ندارند، مانند سرورهایی که فقط `claude` معمولی دارد، در آخر و بدون نام می‌آیند. `ccam ext --json` همین گزارش را چاپ می‌کند (`accounts`، پروژه، `problems`)، با `remote` روی سرورهای MCP که Claude Code از طریق شبکه به آن‌ها می‌رسد؛ وقتی حسابی وجود نداشته باشد، گزارش را با `accounts: []` چاپ می‌کند.

## آنچه می‌خواند

| نوع | کجا | منبع نمایش‌داده‌شده |
| --- | --- | --- |
| مهارت‌ها | هر پوشه در `skills/` که `SKILL.md` دارد | گروهی که `skills/` به آن پیوند دارد، یا `account` |
| فرمان‌ها، agentها | فایل‌های `.md` زیر `commands/` و `agents/`، زیرپوشه‌ها به شکل `folder/name` | گروه، یا `account` |
| پلاگین‌ها | `plugins/installed_plugins.json`، که با `enabledPlugins` در تنظیمات حساب و پروژه روشن یا خاموش می‌شوند؛ پلاگینی که فقط برای یک پروژه نصب شده فقط در همان پروژه دیده می‌شود | گروه، یا `account`؛ `(off)` وقتی روشن نیست |
| سرورهای MCP | `mcp.json` گروه (وقتی MCP مشترک است)، `.claude.json` حساب (سرورهای خودش، و آن‌هایی که فقط برای این پوشه دارد)، و سرورهای پلاگین‌هایش، به شکل `plugin:<plugin>:<server>` | گروه، `account` یا `local` |
| هوک‌ها | `hooks` در `settings.json`، شمارش‌شده برای هر رویداد | گروه، یا `account` |
| این پروژه | `.claude/skills`، `.claude/commands`، `.claude/agents`، `.claude/settings.json`، `.claude/settings.local.json` و `.mcp.json` در پوشه فعلی | یک بار فهرست می‌شود؛ حسابی که به پوشه اعتماد نکرده است برای سرورهای `.mcp.json` آن یک مشکل می‌گیرد |

## مشکلاتی که گزارش می‌کند

- **پیوندهای خراب**: پیوندی در `skills/`، `commands/` یا `agents/` که مقصدش از بین رفته، یا پوشه‌ای مشترک (`skills`، `commands`، `agents`، `plugins`) که به جایی ناموجود پیوند دارد (با `ccam doctor --fix` تعمیرش کنید که پوشه گروه گم‌شده را از نو می‌سازد).
- **پوشه مهارت بدون `SKILL.md`**: Claude Code از آن می‌گذرد.
- **سرور MCP که در دو جا متفاوت تعریف شده**: `mcp.json` گروه و `.claude.json` حساب، یا یکی از آن‌ها و `.mcp.json` پروژه. تعریف‌های یکسان در هر دو جا (که `mcp_strategy = "sync"` می‌نویسد) مشکل نیست. سرور خود یک حساب را با `ccam mcp import <account> <server>` به گروه منتقل کنید.
- **`installed_plugins.json` که JSON معتبر نیست**، پلاگینی که `enabledPlugins` نامش را نمی‌برد، و پلاگینی که پوشه‌اش وجود ندارد.
- **سرورهایی که فقط یک حساب بارگذاری می‌کند**: سروری که فقط در `.claude.json` این حساب است و حساب‌های دیگر گروهش آن را بارگذاری نمی‌کنند، و سروری در دامنه local حساب برای این پوشه (کاری که `claude mcp add` انجام می‌دهد، مگر با `-s user`). هر کدام همراه با فرمان `ccam mcp import` می‌آید که آن را به اشتراک می‌گذارد. سرورهایی در `~/.claude.json` مربوط به `claude` معمولی که هیچ حسابی بارگذاری‌شان نمی‌کند همراه با فرمان `ccam migrate --mcp` می‌آیند.
- **سرور پروژه‌ای که منتظر تأیید است**: حساب هنوز به پوشه اعتماد نکرده است، پس Claude Code یک بار می‌پرسد و تا آن زمان سرورهای `.mcp.json` پروژه منتظر می‌مانند. مشکل می‌گوید آن‌ها را کجا تأیید کنید: `"enabledMcpjsonServers"` در `settings.json` حساب‌های گروه تنظیمات، یا این حساب. اعتمادی که به worktree اصلی یک مخزن git داده شده برای worktreeهای دیگر آن هم اعمال می‌شود.
- **ورودها و connectorها**: سرور MCP که در این حساب به ورود نیاز دارد در حالی که حساب دیگری به آن وارد شده است (هر ورود فقط برای یک حساب معتبر است؛ با `/mcp` وارد شوید)، و connectorهای claude.ai که حساب دیگری به کار برده و این حساب نه.

## ایمنی و موارد خاص

- `ccam ext` چیزی را تغییر نمی‌دهد و هیچ درخواست شبکه‌ای نمی‌فرستد. نام سرورها و شمار هوک‌ها را می‌خواند و هرگز مقدارهای `env` یا سرآیندها را در تعریف‌های MCP نمی‌خواند. برای بررسی ورودها، نام ورودهای MCP را در ورود هر حساب می‌خواند (در macOS از Keychain)، هرگز توکن‌ها را، و فقط برای حساب‌هایی که سروری را بارگذاری می‌کنند که حساب دیگری هم بارگذاری می‌کند.
- آنچه را روی دیسک هست نشان می‌دهد. سروری که تنظیمات مدیریت‌شده مسدودش کرده‌اند باز هم نمایش داده می‌شود؛ پلاگینی که در تنظیمات خاموش است `(off)` نشان می‌دهد.
- آن را در پوشه پروژه اجرا کنید تا پوشه `.claude` و `.mcp.json` همان پروژه هم در نظر گرفته شوند.

## مطالب مرتبط

- [اشتراک‌گذاری با گروه‌ها](/features/share)
- [سرورهای MCP مشترک](/features/mcp)
- [بررسی و تعمیر](/features/doctor)
