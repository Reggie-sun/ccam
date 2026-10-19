---
title: مرجع فرمان‌ها
description: همه فرمان‌ها و گزینه‌های ccam، از ccam add و ccam <name> تا ccam usage، ccam next و ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

برای اینکه ببینید یک فرمان چه کاری انجام می‌دهد بدون آنکه واقعاً انجامش دهد، `--dry-run` را پیش از آن بگذارید، برای مثال `ccam --dry-run work`.

<!-- commands -->

## اجرای Claude Code

`ccam <name|number> [arguments]` اجرای Claude Code با آن حساب را آغاز می‌کند. ccam آرگومان‌های پس از نام را نمی‌خواند؛ آن‌ها را به همان ترتیب به `claude` می‌دهد، و ورودی استاندارد، کد خروج و سیگنال‌ها دقیقاً مانند خود `claude` رفتار می‌کنند. گذاشتن `--` پس از نام اختیاری است.

`ccam` با آرگومان ولی بدون نام حساب، مانند `ccam -p "summarise"`، از حساب پیش‌فرض استفاده می‌کند. حساب پیش‌فرض را با `ccam use <name>` تعیین کنید.

میان‌برهای `claude1`، `claude2` و مانند آن‌ها پیوندهایی به ccam در `~/.local/bin` هستند. `ccam shims install --names` برای هر حساب `claude-<name>` را هم می‌سازد.

## افزودن و حذف حساب‌ها

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` بدون نام پوشه، به دنبال `~/.claude` و `~/.claude-*` می‌گردد و درباره هر کدام از شما می‌پرسد. پوشه‌های واردشده در همان جایی که هستند ثبت می‌شوند؛ ccam هرگز یک ورود را جابه‌جا یا کپی نمی‌کند.

## اشتراک‌گذاری

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

هر مورد یکی از این مقدارها را می‌پذیرد: `on` (گروه `default`)، `off` یا نام یک گروه. [اشتراک‌گذاری](/docs/sharing) را ببینید.

## سرورهای MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## مصرف و جابه‌جایی

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## نگهداری

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## فایل‌ها

| مسیر | محتوا |
| --- | --- |
| `~/.config/ccam/config.toml` | حساب‌ها و آنچه به اشتراک می‌گذارند |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` هر حساب |
| `~/.local/share/ccam/groups/<group>/` | تنظیمات، سرورهای MCP، افزونه‌ها و تاریخچه مشترک |
| `~/.local/state/ccam/` | حافظه نهان مصرف و زمان‌های آخرین استفاده |
| `~/.local/bin/claude1`، `claude2`… | پیوندهایی به ccam |

متغیرهای `XDG_CONFIG_HOME`، `XDG_DATA_HOME` و `XDG_STATE_HOME` محل این پوشه‌ها را تغییر می‌دهند.
