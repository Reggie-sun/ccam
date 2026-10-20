---
title: بررسی افزونه‌ها
description: ccam ext مهارت‌ها، فرمان‌ها، پلاگین‌ها، سرورهای MCP و هوک‌هایی را که هر حساب بارگذاری می‌کند، منبع هر کدام و موارد خراب را نشان می‌دهد.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` به این پرسش پاسخ می‌دهد: «چرا این مهارت (یا سرور MCP) در یک حساب کار می‌کند و در دیگری نه؟». برای هر حساب، افزونه‌هایی را که Claude Code بارگذاری می‌کند و منبع هر کدام را فهرست می‌کند: یک گروه اشتراک، پوشه خود حساب، یا پروژه‌ای که در آن هستید. زیر جدول، چیزهایی را فهرست می‌کند که مانع بارگذاری یک افزونه به شکلی که انتظار دارید می‌شوند. فقط می‌خواند.

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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

هر خانه منبع آن مورد در آن حساب را نشان می‌دهد: گروه اشتراک (`default`)، `account` وقتی در پوشه خود حساب است، یا `-` وقتی حساب آن را ندارد. `!` موردی را که مشکلی دارد علامت می‌زند که زیر جدول توضیح داده می‌شود.

## آنچه می‌خواند

| نوع | کجا | منبع نمایش‌داده‌شده |
| --- | --- | --- |
| مهارت‌ها | هر پوشه در `skills/` که `SKILL.md` دارد | گروهی که `skills/` به آن پیوند دارد، یا `account` |
| فرمان‌ها، agentها | فایل‌های `.md` زیر `commands/` و `agents/`، زیرپوشه‌ها به شکل `folder/name` | گروه، یا `account` |
| پلاگین‌ها | `plugins/installed_plugins.json` | گروه، یا `account` |
| سرورهای MCP | `mcp.json` گروه (وقتی MCP مشترک است) و `.claude.json` حساب | گروه، یا `account` |
| هوک‌ها | `hooks` در `settings.json`، شمارش‌شده برای هر رویداد | گروه، یا `account` |
| این پروژه | `.claude/skills`، `.claude/commands`، `.claude/agents`، `.claude/settings.json`، `.claude/settings.local.json` و `.mcp.json` در پوشه فعلی | یک بار، برای همه حساب‌ها فهرست می‌شود |

## مشکلاتی که گزارش می‌کند

- **پیوندهای خراب**: پیوندی در `skills/`، `commands/` یا `agents/` که مقصدش از بین رفته، یا پوشه‌ای مشترک (`skills`، `commands`، `agents`، `plugins`) که به جایی ناموجود پیوند دارد (با `ccam doctor --fix` تعمیرش کنید).
- **پوشه مهارت بدون `SKILL.md`**: Claude Code از آن می‌گذرد.
- **سرور MCP که در دو جا متفاوت تعریف شده**: `mcp.json` گروه و `.claude.json` حساب، یا یکی از آن‌ها و `.mcp.json` پروژه. تنظیمات یکسان در هر دو جا (که `mcp_strategy = "sync"` می‌نویسد) مشکل نیست. سرور خود یک حساب را با `ccam mcp import <account> <server>` به گروه منتقل کنید.
- **`installed_plugins.json` که JSON معتبر نیست.**

## ایمنی و موارد خاص

- `ccam ext` چیزی را تغییر نمی‌دهد و هیچ درخواست شبکه‌ای نمی‌فرستد. نام سرورها و شمار هوک‌ها را می‌خواند و هرگز مقدارهای `env` یا سرآیندها را در تعریف‌های MCP نمی‌خواند.
- آنچه را روی دیسک هست نشان می‌دهد. پلاگینی که نصب شده اما در تنظیمات خاموش است، یا سروری که تنظیمات مدیریت‌شده مسدودش کرده‌اند، باز هم نمایش داده می‌شود.
- آن را در پوشه پروژه اجرا کنید تا پوشه `.claude` و `.mcp.json` همان پروژه هم در نظر گرفته شوند.

## مطالب مرتبط

- [اشتراک‌گذاری بین حساب‌ها](/features/share)
- [سرورهای MCP مشترک](/features/mcp)
- [بررسی و تعمیر](/features/doctor)
