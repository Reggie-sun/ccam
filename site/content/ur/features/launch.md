---
title: "Claude Code شروع کرنا"
description: "ccam <name|number>، ccam run یا ccam کے انتخابی مینو سے Claude Code کو کسی بھی اکاؤنٹ کے طور پر شروع کریں؛ آرگیومنٹس بغیر تبدیلی کے claude کو جاتے ہیں، اور --dry-run کمانڈ دکھاتا ہے۔"
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 401eed57e533
---

`ccam work`، Claude Code کو اکاؤنٹ `work` کے طور پر شروع کرتا ہے، اور نام کے بعد سب کچھ جوں کا توں `claude` کو جاتا ہے۔ اسے ہر اس جگہ استعمال کریں جہاں آپ `claude` ٹائپ کرتے: ٹرمینل میں، اسکرپٹ میں، یا pipe سے آنے والے اِن پٹ کے ساتھ۔

## استعمال

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

اکاؤنٹ کا نام یا `ccam list` سے اس کا نمبر دیں۔ ccam نام کے بعد والے آرگیومنٹس نہیں پڑھتا، اس لیے `ccam work --help`، Claude Code کی مدد دکھاتا ہے۔ نام کے بعد `--` اختیاری ہے۔ اسٹینڈرڈ اِن پٹ، ایگزٹ کوڈ اور سگنلز ویسے ہی کام کرتے ہیں جیسے `claude` کے ساتھ۔

جب پہلا آرگیومنٹ `-` سے شروع ہو اور ccam کا اپنا کوئی فلیگ نہ ہو، تو ccam تمام آرگیومنٹس کے ساتھ ڈیفالٹ اکاؤنٹ شروع کرتا ہے۔ ڈیفالٹ وہ اکاؤنٹ ہے جو `ccam use` سے مقرر کیا گیا ہو، یا اگر آپ نے کوئی مقرر نہ کیا ہو تو سب سے کم نمبر والا اکاؤنٹ۔

### انتخابی مینو

بغیر آرگیومنٹس کے `ccam` انتخابی مینو کھولتا ہے جب اسٹینڈرڈ اِن پٹ اور آؤٹ پٹ ٹرمینل ہوں، ورنہ `ccam list` پرنٹ کرتا ہے۔ یہ ہر اکاؤنٹ کا نمبر، ای میل، 5 گھنٹے کا استعمال، آخری استعمال اور شیئر کیے گئے آئٹمز دکھاتا ہے؛ استعمال ہر اکاؤنٹ کے جواب دینے پر بھرتا جاتا ہے۔

| کی | کیا کرتی ہے |
| --- | --- |
| `Enter` | منتخب اکاؤنٹ کے طور پر Claude Code شروع کرتی ہے |
| `1` سے `9` | اس نمبر والا اکاؤنٹ شروع کرتی ہے |
| `Up`/`Down`، `k`/`j` | انتخاب کو اوپر نیچے کرتی ہے |
| `a` | مینو بند کر کے `ccam add` کا وزرڈ چلاتی ہے |
| `s` | منتخب اکاؤنٹ کے شیئر کیے گئے آئٹمز میں ترمیم (`Space` یا `x` بدلتی ہے، `Enter` لاگو کرتی ہے، `Esc` واپس جاتی ہے) |
| `d` | منتخب اکاؤنٹ کو ڈیفالٹ بناتی ہے |
| `r` | مختصر کیش کے بغیر استعمال دوبارہ حاصل کرتی ہے |
| `q`، `Esc`، `Ctrl+C` | باہر نکلتی ہے |

مینو کوئی عمل کرنے سے پہلے بند ہو جاتا ہے، اس لیے Claude Code شروع ہونے سے پہلے ٹرمینل بحال ہو جاتا ہے۔ `d` یا شیئرنگ کی تبدیلی کے بعد یہ دوبارہ کھلتا ہے۔ شیئرنگ کی تبدیلی پہلے منصوبہ دکھاتی ہے اور پوچھتی ہے، جیسے `ccam share` کرتا ہے۔

## آپشنز

| آپشن | کیا کرتا ہے | ڈیفالٹ |
| --- | --- | --- |
| `--dry-run` | وہ انوائرنمنٹ اور کمانڈ پرنٹ کرتا ہے جو ccam چلاتا، پھر باہر نکل جاتا ہے | بند |
| `-v`، `--version` | `ccam <version>` پرنٹ کرتا ہے (صرف پہلے آرگیومنٹ کے طور پر) | |
| `-h`، `--help` | ccam کی مدد، `ccam --help` یا `ccam run --help` کے طور پر | |

آرگیومنٹس کے اصول:

- ccam اپنے فلیگز صرف اکاؤنٹ کے نام سے پہلے پڑھتا ہے۔ نام کے بعد ہر آرگیومنٹ Claude Code کا ہے، بشمول `--dry-run` اور `--help`۔
- `--dry-run`، `ccam --dry-run work`، `ccam run --dry-run work` اور `ccam --dry-run` (ڈیفالٹ اکاؤنٹ) کے طور پر کام کرتا ہے۔ یہ لاگو ہوتا ہے، مسترد نہیں ہوتا۔
- `ccam run <name|number> [claude args]`، `ccam <name>` کی لمبی شکل ہے۔ نام کے بغیر `ccam run` ڈیفالٹ اکاؤنٹ استعمال کرتا ہے۔
- اکاؤنٹس کے نام `add`، `list`، `run` یا `usage` جیسے ccam کمانڈز کے نام نہیں ہو سکتے، اس لیے `ccam <word>` کبھی مبہم نہیں ہوتا۔

ڈرائی رن ایک لائن پرنٹ کرتا ہے جسے آپ شیل میں چسپاں کر سکتے ہیں:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## یہ کیسے کام کرتا ہے

ہر لانچ کمانڈ لائن سے `exec` تک وہی مراحل اٹھاتا ہے:

<!-- motion:launch -->

ہر خانہ اپنا مرحلہ چلنے پر روشن ہوتا ہے، اور نیچے کی لائنیں وہ انوائرنمنٹ اور آرگیومنٹس ہیں جن کے ساتھ Claude Code شروع ہوتا ہے۔

1. `cmd/ccam/main.go`، `app.Main` کو کال کرتا ہے۔ اگر ccam کسی اور نام سے شروع ہوا ہو تو یہ ایک [شارٹ کٹ](/features/shortcuts) ہے اور سیدھا لانچر پر جاتا ہے۔
2. `app.parse`، cobra کے شامل ہونے سے پہلے `os.Args` دیکھتا ہے۔ صرف اکیلا `ccam` (انتخابی مینو)، ccam کمانڈ کا نام (`config.Reserved`)، `-h`/`--help` اور `-v`/`--version`، `cli.Execute` پر جاتے ہیں۔ باقی سب کچھ لانچ ہے۔
3. `launcher.Launch`، `~/.config/ccam/config.toml` پڑھتا ہے، جو واحد فائل ہے جسے یہ پارس کرتا ہے۔
4. `launcher.Prepare` اکاؤنٹ کو نام یا نمبر سے تلاش کرتا ہے، شروع کا ایک `--` ہٹاتا ہے، اور `ResolveClaude` سے `claude` قابلِ عمل فائل تلاش کرتا ہے۔ PATH کی تلاش فائل کے ترمیم کے وقت کے ساتھ `~/.local/state/ccam/claude-path` میں کیش ہوتی ہے، اس لیے ccam ہر لانچ پر `PATH` تلاش نہیں کرتا۔ بدلی ہوئی یا غائب فائل نئی تلاش شروع کرتی ہے۔ `config.toml` میں `claude_bin` کا راستہ جوں کا توں استعمال ہوتا ہے۔ ccam ایسا `claude` مسترد کرتا ہے جو خود ccam ہو، inode سے موازنہ کر کے، کیونکہ اس سے لوپ بن جائے گا۔
5. `checkLinks` اکاؤنٹ فولڈر `~/.local/share/ccam/profiles/<name>/` کو اس کی شیئرنگ سیٹنگز سے `share.Check` کے ذریعے موازنہ کرتا ہے، جو صرف `lstat` اور `readlink` کال کرتا ہے، ہر آئٹم کے لیے ایک یا دو سسٹم کالز۔ غائب لنک فوراً بنا دیا جاتا ہے۔ جس چیز سے ڈیٹا کھو سکتا ہو وہ لانچ روک دیتی ہے۔
6. `applyMCP` شیئر کیے گئے MCP سرورز شامل کرتا ہے (نیچے دیکھیں)۔
7. `session.Register` سیشن ریکارڈ کرتا ہے۔ ccam اپنے PID کے نام سے ایک نشان، جس میں پروسیس کے شروع ہونے کا وقت ہوتا ہے، `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` میں لکھتا ہے (یا `XDG_RUNTIME_DIR` کے بغیر `~/.local/state/ccam/locks/` میں)۔ گفتگو کی ہسٹری شیئر کرنے والے اکاؤنٹ کے لیے، `historyCheck` پہلے `history-<group>-<hash>` میں کوئی چلتا ہوا سیشن تلاش کرتا ہے، جہاں hash موجودہ فولڈر کا ہوتا ہے، اور وہاں بھی رجسٹر کرتا ہے۔ `exec`، PID برقرار رکھتا ہے، اس لیے نشان پھر خود Claude Code پروسیس کا نام بتاتا ہے۔ پکڑا ہوا لاک نہیں بلکہ نشان، کیونکہ Claude Code جو پروسیس شروع کرتا ہے وہ لاک وراثت میں پا کر سیشن ختم ہونے کے بعد بھی اسے رکھے رہیں گے۔
8. `Plan.Exec` ٹرمینل کا عنوان `claude · <name>` مقرر کرتا ہے (جب اسٹینڈرڈ ایرر ٹرمینل ہو)، لانچ کا وقت `fsync` کے بغیر `~/.local/state/ccam/last-used/<name>` میں لکھتا ہے، اور `syscall.Exec` کال کرتا ہے۔

Claude Code اس طرح چلتا ہے کہ `CLAUDE_CONFIG_DIR` اکاؤنٹ فولڈر پر اور `CCAM_PROFILE` اکاؤنٹ کے نام پر مقرر ہوتا ہے۔ ccam ان کی وہ قدریں ہٹا دیتا ہے جو اسے وراثت میں ملی ہوں۔ macOS پر `ccam import` سے رجسٹر کیے گئے `~/.claude` فولڈر کے لیے یہ `CLAUDE_SECURESTORAGE_CONFIG_DIR` کو بھی خالی قدر پر مقرر کرتا ہے، تاکہ اکاؤنٹ اور سادہ `claude` ایک ہی Keychain لاگ اِن استعمال کریں۔

### exec کیوں

`syscall.Exec`، ccam پروسیس کی جگہ `claude` لے آتا ہے، اس لیے ccam میموری میں نہیں رہتا، کی اسٹروکس، سگنلز اور ایگزٹ کوڈ Claude Code کے اپنے ہوتے ہیں، اور آرگیومنٹس ایک array رہتے ہیں، کبھی شیل اسٹرنگ نہیں بنتے۔ سیشن کے بعد کا کام، جیسے `claude mcp add` سے شامل کیے گئے سرورز کو نوٹ کرنا، اگلے لانچ پر ہوتا ہے: ccam، `.claude.json` کے ترمیم کے وقت کو `~/.local/state/ccam/mcp/<name>/seen` سے موازنہ کرتا ہے اور فائل صرف تب پارس کرتا ہے جب وہ بدلی ہو۔

### کمانڈ لائن کس راستے پر جاتی ہے

تین کمانڈ لائنیں، `internal/app` سے گزرنے کے تین راستے:

<!-- motion:launch-route -->

شارٹ کٹ `argv[0]` پر، اور `ccam <name>` پہلے آرگیومنٹ پر پہچانا جاتا ہے، اور دونوں `syscall.Exec` پر ختم ہوتے ہیں۔ صرف `ccam usage` جیسی ccam کمانڈ cobra تک پہنچتی ہے۔

### شیئر کیے گئے MCP سرورز

ڈیفالٹ `flag` حکمتِ عملی کے ساتھ، ccam گروپ کے MCP سرورز Claude Code کو ایک آرگیومنٹ کے طور پر دیتا ہے، جو آپ کے آرگیومنٹس سے پہلے رکھا جاتا ہے:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` آرگیومنٹ صرف تب لوٹاتا ہے جب گروپ کی `mcp.json` موجود ہو، اور `Plan.InsertArgs` اسے `claude` کے فوراً بعد رکھتا ہے۔ آپ کا اپنا `--mcp-config` یا `--strict-mcp-config` پھر بھی اس کے اوپر لاگو ہوتا ہے۔ یہ ایک ہی آرگیومنٹ ہے، `--mcp-config=<path>`، کیونکہ Claude Code کا `--mcp-config` کئی قدریں لیتا ہے: دو آرگیومنٹس کے طور پر لکھا جائے تو یہ بعد میں آنے والے پرامپٹ کو بھی کنفیگ فائل سمجھ کر پڑھے گا۔

`config.toml` میں `mcp_strategy = "sync"` کے ساتھ، ccam اس کے بجائے سرورز اکاؤنٹ کی `.claude.json` میں لکھتا ہے، اور اگر اکاؤنٹ پہلے سے چل رہا ہو تو یہ چھوڑ دیتا ہے۔

### Remote Control

`ccam work rc` (یا `ccam work remote-control`)، Claude Code کا Remote Control اکاؤنٹ `work` کے طور پر شروع کرتا ہے۔ `claude remote-control` اپنے فعل (verb) سے پہلے `--mcp-config` کو مسترد کرتا ہے اور اس کے بعد کوئی نہیں لیتا، اور اس کے سیشنز MCP سرورز صرف اکاؤنٹ کی `.claude.json` سے پڑھتے ہیں۔ اس لیے یہاں فلیگ والی حکمتِ عملی مختلف طریقے سے کام کرتی ہے:

<!-- motion:launch-rc -->

1. `remoteControl` اکاؤنٹ کے بعد پہلے ایسے آرگیومنٹ کو جو `-` سے شروع نہ ہو، Claude Code کا فعل مانتا ہے۔ جب یہ `remote-control` یا `rc` ہو تو `mcp.FlagArgs` کچھ شامل نہیں کرتا۔
2. `mcp.Synced` گروپ کی `mcp.json` کے ترمیم کے وقت کا موازنہ اکاؤنٹ کے `sync-stamp` سے کرتا ہے۔ اگر `.claude.json` میں پہلے سے گروپ کے موجودہ سرورز ہوں، یا گروپ کے پاس کوئی سرور نہ ہو، تو Remote Control بغیر سوال کے شروع ہو جاتا ہے۔
3. ورنہ ccam ان سرورز کے نام بتاتا ہے جو سیشنز کو نہیں ملیں گے اور پوچھتا ہے:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. `y` پر، `mcp.Sync` انہیں اسی طرح کاپی کرتا ہے جیسے `sync` حکمتِ عملی کرتی ہے: Claude Code کے لاک کے تحت، بیک اپ کے بعد، اور اکاؤنٹ کے اپنے سرورز برقرار رکھتے ہوئے۔ `N` پر، یا ٹرمینل کے بغیر، Remote Control ان کے بغیر شروع ہوتا ہے، اور ccam، `mcp_strategy = "sync"` کی طرف اشارہ کرتا ہے، جو ہر لانچ سے پہلے انہیں کاپی کرتی ہے۔

ccam ہر Remote Control لانچ پر پوچھتا ہے جب تک سرورز کاپی نہ ہو جائیں، اور کاپی کے بعد صرف تب جب گروپ کی `mcp.json` دوبارہ بدلے۔ بعد کے عام لانچز میں `.claude.json` کی کاپی کو `--mcp-config` پر ترجیح ملتی ہے، اس لیے `ccam mcp add`، `remove` اور `import` وہ کاپیاں اپ ڈیٹ کرتے ہیں جو دوسرے اکاؤنٹس کے پاس ہیں ([شیئر کیے گئے MCP سرورز](/features/mcp))۔

### 10 ms کا بجٹ

ccam، Claude Code شروع ہونے سے پہلے زیادہ سے زیادہ 10 ms (p95) شامل کرتا ہے۔ لانچ کا راستہ کبھی cobra کمانڈ ٹری نہیں بناتا، Bubble Tea شروع نہیں کرتا اور نیٹ ورک کنکشن نہیں کھولتا، اور اگر پیکیجز کے `init()` فنکشنز کل ملا کر 5 ms سے زیادہ لیں تو `TestPackageInitIsCheap` ناکام ہوتا ہے۔ `scripts/perf-check.sh`، `hyperfine` سے `ccam --dry-run work`، `ccam work` اور ایک شارٹ کٹ کو ناپتا ہے اور 10 ms سے اوپر ناکام ہوتا ہے، اور `internal/launcher` میں `BenchmarkLaunchPath` مراحل کو لوپ میں چلاتا ہے۔ ہر مرحلے کا وقت دیکھنے کے لیے `CCAM_TRACE=1` مقرر کریں:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## حفاظت اور خاص صورتیں

- نامعلوم نام یا نمبر `no account "<name>"` کے ساتھ رک جاتا ہے اور `ccam list` کی طرف رہنمائی کرتا ہے۔ اگر کوئی اکاؤنٹ ہی نہ ہو تو ccam کہتا ہے کہ `ccam add <name>` سے ایک شامل کریں۔
- اگر کوئی شیئر کیا گیا راستہ عام فائل بن گیا ہو یا غلط جگہ اشارہ کرتا ہو، تو ccam، `stopped because of share link problems:` کے ساتھ رک جاتا ہے، ان کی فہرست دیتا ہے، اور `ccam doctor --fix` چلانے کو کہتا ہے۔ یہ لانچ کے راستے پر ڈیٹا منتقل نہیں کرتا۔
- اگر کوئی دوسرا Claude Code جو اس اکاؤنٹ کے ساتھ گفتگو کی ہسٹری شیئر کرتا ہو (ایک ہی گروپ) اسی فولڈر میں چل رہا ہو، تو ccam انتباہ دیتا ہے کہ ہسٹری آپس میں مل سکتی ہے اور `Run anyway? [y/N]` پوچھتا ہے۔ ٹرمینل کے بغیر یہ انتباہ دکھا کر جاری رکھتا ہے۔
- `--dry-run`، MCP sync، نئے سرور کی اطلاع، ہسٹری کی جانچ، سیشن کا نشان، عنوان اور آخری استعمال کا وقت چھوڑ دیتا ہے، اور exec نہیں کرتا۔ یہ پھر بھی شیئر لنکس جانچتا ہے، اور عام لانچ کی طرح غائب لنک بنا دیتا ہے۔
- `ccam work rc`، جب `work` پہلے سے چل رہا ہو، `.claude.json` کو نہیں چھیڑتا، یہ بتا دیتا ہے، اور Remote Control کو گروپ کے سرورز کے بغیر شروع کرتا ہے؛ انہیں کاپی کرنے کے لیے سیشن ختم ہونے کے بعد اسے دوبارہ چلائیں۔ `--dry-run` کے ساتھ یہ پیغام پرنٹ کرتا ہے لیکن نہ پوچھتا ہے نہ کاپی کرتا ہے۔
- `mcp_strategy = "sync"` کے ساتھ، ناکام sync، `ccam: MCP sync failed (launching anyway): …` پرنٹ کرتا ہے اور Claude Code ان سرورز کے ساتھ شروع ہوتا ہے جو اس کے پاس پہلے سے تھے۔

## متعلقہ

- [شارٹ کٹ کمانڈز](/features/shortcuts)
- [ccam کے ساتھ Remote Control استعمال کرنا](/guides/remote-control)
- [اکاؤنٹس کا انتظام](/features/accounts)
- [کمانڈز کا حوالہ](/docs/commands)
- [MCP سرورز، سیٹنگز اور ہسٹری شیئر کریں](/docs/sharing)
- [CLAUDE_CONFIG_DIR کیسے کام کرتا ہے](/guides/claude-config-dir)
