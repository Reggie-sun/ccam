---
title: "যাচাই ও মেরামত"
description: "ccam doctor দিয়ে শেয়ার লিংক, ফাইলের অনুমতি, লগইন ও Claude Code সংস্করণ যাচাই করুন, আর ccam doctor --fix দিয়ে নিরাপদে মেরামতযোগ্য সমস্যা ঠিক করুন।"
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` আপনার ccam ইনস্টলেশন যাচাই করে: Claude Code বাইনারি, প্রতিটি অ্যাকাউন্টের ফোল্ডার, লগইন ও শেয়ার লিংক, শেয়ার করা MCP সার্ভার, ব্যবহার এন্ডপয়েন্ট এবং আপনার শর্টকাট কমান্ড। `ccam doctor --fix` ডেটা না হারিয়ে যা করা যায় তা মেরামত করে, এবং যে ফাইল প্রতিস্থাপন করে সেটি আগে একটি ব্যাকআপ ফোল্ডারে সরায়। `ccam <name>` শেয়ার লিংকের সমস্যায় থেমে গেলে, হাতে ফাইল পুনরুদ্ধার করার পরে, বা Claude Code আপডেটের পরে এটি চালান।

## ব্যবহার

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

প্রতিটি ফলাফল একটি চিহ্নসহ এক লাইনের: `✓` ঠিক আছে, `·` তথ্য, `!` সতর্কতা, `✗` ত্রুটি। অ্যাকাউন্টের ফলাফল অ্যাকাউন্টের নাম দিয়ে শুরু হয়।

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

`--fix` দিলে মেরামত হওয়া প্রতিটি ফলাফলের শেষে থাকে `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--fix` | নিরাপদে মেরামতযোগ্য সমস্যা মেরামত করে (বদলানো ফাইলের ব্যাকআপ রাখা হয়) | বন্ধ |
| `--dry-run` (গ্লোবাল) | প্রতিটি যাচাই চালায় কিন্তু কিছু মেরামত করে না, `--fix` দিলেও | বন্ধ |

`ccam doctor` কোনো আর্গুমেন্ট নেয় না। এটি `--dry-run` মানে: `ccam --dry-run doctor --fix` `ccam doctor`-এর মতোই ফলাফল প্রিন্ট করে এবং কিছু বদলায় না।

অন্তত একটি `✗` ফলাফল মেরামত না হয়ে থাকলে এক্সিট কোড `1`, অন্যথায় `0`। সতর্কতা ও তথ্য কখনও রানকে ব্যর্থ করে না। `--fix` ছাড়া রান ব্যর্থ হলে ccam একটি ফাঁকা লাইন ও `` `ccam doctor --fix` repairs what it can. `` দিয়ে শেষ করে। লোড করা যায় না এমন কনফিগও `1` ও একটি `ccam:` ত্রুটি দিয়ে শেষ হয়।

## এটি কীভাবে কাজ করে

`internal/cli/doctor.go` `~/.config/ccam/config.toml` লোড করে এবং `internal/doctor/doctor.go`-এর `doctor.Run` ডাকে, যা একটি নির্দিষ্ট ক্রমে যাচাইগুলো চালায় এবং ফলাফল ফেরত দেয়। `--fix` `Options.Fix` হিসেবে পাঠানো হয়, ড্রাই রানে যা false থাকে।

<!-- motion:doctor-checks -->

মাঝের অংশটি প্রতিটি অ্যাকাউন্টের জন্য একবার চলে। যাচাইগুলো, এবং প্রতিটির জন্য `--fix` কী করে:

1. **Claude Code সংস্করণ** (`checkClaude`)। লঞ্চার যেভাবে করে সেভাবে `claude` বাইনারি খুঁজে বের করে এবং `claude --version` চালায়। সংস্করণটি `2.1.302` দিয়ে শুরু না হলে, অর্থাৎ যে সংস্করণের সঙ্গে ccam-এর ফোল্ডার কাঠামো যাচাই করা হয়েছে (`doctor.LayoutVersion`), এটি যোগ করে `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`। বাইনারি না থাকা বা `--version` ব্যর্থ হওয়া একটি ত্রুটি। কোনো মেরামত নেই।
2. **এনভায়রনমেন্টে ক্রেডেনশিয়াল** (`checkAuthEnv`)। `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` বা `ANTHROPIC_AUTH_TOKEN` নির্ধারিত থাকলে সতর্ক করে: ওই শেল থেকে চালু হওয়া প্রতিটি অ্যাকাউন্ট সেটি উত্তরাধিকারে পায়, এবং Claude Code অ্যাকাউন্টের নিজের লগইনের চেয়ে সেটিকে অগ্রাধিকার দেয়। কোনো মেরামত নেই; নিজে সেটি আনসেট করুন।
3. **ফোল্ডার** (`checkProfile`)। অ্যাকাউন্ট ফোল্ডার না থাকা একটি ত্রুটি এবং ওই অ্যাকাউন্টের যাচাই সেখানেই শেষ হয়। গ্রুপ বা অন্যদের কোনো বিট নির্ধারিত থাকা অনুমতি একটি সতর্কতা; `--fix` `chmod 0700` চালায়।
4. **লগইন ফাইল** (`checkCredentials`)। `.credentials.json` `lstat` দিয়ে যাচাই করা হয়। সিমলিংক একটি ত্রুটি, কারণ প্রতিটি অ্যাকাউন্টের নিজস্ব টোকেন থাকতে হবে, এবং এটি মেরামত করা হয় না। `0600` ছাড়া অন্য মোড একটি সতর্কতা; `--fix` `chmod 0600` চালায়। তারপর লগইন পড়া হয় (macOS-এ প্রথমে Keychain থেকে)। লগইন না থাকা একটি সতর্কতা, যা `ccam login <name>` দেখিয়ে দেয়। মেয়াদোত্তীর্ণ রিফ্রেশ টোকেন হলো ত্রুটি `login expired`; মেয়াদোত্তীর্ণ অ্যাক্সেস টোকেন শুধু তথ্য, কারণ Claude Code পরের বার চালু হলে সেটি নবায়ন করে। এগুলোর কোনোটিই মেরামত করা হয় না।
5. **শেয়ার লিংক** (`share.Check`)। শুধু `lstat` ও `readlink` ব্যবহার করে `share.Layout`-এর প্রতিটি পাথ অ্যাকাউন্টের শেয়ার সেটিংসের সঙ্গে তুলনা করে। এটি চার ধরনের সমস্যা জানায়: লিংক না থাকা, অন্য কোথাও নির্দেশ করা লিংক, যেখানে লিংক থাকার কথা সেখানে আসল ফাইল, এবং এখন লোকাল হয়ে যাওয়া আইটেমের জন্য রয়ে যাওয়া লিংক। লিংক না থাকা একটি সতর্কতা, বাকিগুলো ত্রুটি।
6. **অজানা এন্ট্রি** (`checkUnknownEntries`)। অ্যাকাউন্ট ফোল্ডারের যে শীর্ষ-স্তরের নামগুলো ccam চেনে না (`share.Known`), সেগুলোর তালিকা দেয়, Claude Code-এর কাঠামো বদলেছে এমন ইঙ্গিত হিসেবে। শুধু তথ্য।
7. **প্লাগইন পাথ** (`checkPlugins`)। প্লাগইন রেজিস্ট্রি (`installed_plugins.json`, `known_marketplaces.json`) পূর্ণ পাথ রাখে। প্রতিটি গ্রুপের `plugins/`-এ যে পাথগুলো কোনো একটি অ্যাকাউন্টের ফোল্ডারের মধ্য দিয়ে যায়, সেই অ্যাকাউন্ট সরানো বা নাম বদলানো হলে সেগুলো ভেঙে যায়। `--fix` `profile.RewritePluginPaths` দিয়ে সেগুলো গ্রুপ ফোল্ডারে আবার লেখে, রেজিস্ট্রি ফাইলের একটি অ্যাটমিক রাইট।
8. **MCP সার্ভার** (`checkMCP`)। প্রতিটি গ্রুপের `mcp.json` একবার লোড করে। পড়া যায় না এমন ফাইল একটি ত্রুটি; `mcp.Validate`-এ ব্যর্থ সার্ভার একটি সতর্কতা। কোনো মেরামত নেই।
9. **ব্যবহারের রেসপন্স ফরম্যাট** (`checkUsageFormat`)। ক্যাশ ছাড়া এবং স্বাভাবিক ৩ সেকেন্ডের সীমা দিয়ে একটির পর একটি অ্যাকাউন্টের জন্য ব্যবহার এন্ডপয়েন্টে জিজ্ঞেস করে, যতক্ষণ না কোনোটি উত্তর দেয়। রেসপন্স আর পার্স না হলে এটি সতর্ক করে যে `ccam usage` স্ট্যাটাস লাইনের ক্যাশে ফিরে যায়। এটিই একমাত্র যাচাই যা নেটওয়ার্ক ব্যবহার করে।
10. **শর্টকাট** (`checkShims`)। শুধু যখন কোনো অ্যাকাউন্টের শর্টকাট থাকে। `~/.local/bin`-এ লিংক না থাকা একটি সতর্কতা; `--fix` `shim.Install` ডাকে, যা শুধু অনুপস্থিত লিংক তৈরি করে এবং অন্য ফাইল কখনও স্পর্শ করে না। `PATH`-এ না থাকা bin ফোল্ডার একটি সতর্কতা এবং মেরামত করা হয় না।

### শেয়ার লিংক মেরামত

চালু করার পথ অনুপস্থিত লিংক সঙ্গে সঙ্গে তৈরি করে, কিন্তু যেকোনো কিছু যাতে ডেটা হারাতে পারে তার জন্য থেমে যায় এবং এখানে দেখিয়ে দেয়:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

যেকোনো শেয়ার লিংক সমস্যার জন্য `--fix` অ্যাকাউন্টের শেয়ার সেটিংসকে বর্তমান ও লক্ষ্য দুই অবস্থা হিসেবে দিয়ে `share.Build` ডাকে, তারপর `Plan.Apply`। `Build` ডিস্কে যা আছে তা থেকে প্রতিটি ধাপ নতুন করে বের করে, যা `ccam share`-এর ব্যবহৃত একই পরিকল্পনা, তাই মেরামত শেয়ার চালু করার মতো একই নিয়ম মেনে চলে।

<!-- motion:doctor-repair -->

এখানে `settings.json` একটি আসল ফাইল। গ্রুপে আগে থেকেই আলাদা একটি `settings.json` থাকায় দ্বন্দ্বটি তার ডিফল্ট পছন্দ `KeepGroup` নেয়, তাই গ্রুপের কপি থেকে যায়। অ্যাকাউন্টের ফাইলটি অ্যাকাউন্ট ফোল্ডারের ভেতরে `.ccam-backup/<time>/settings.json`-এ সরে যায়, এবং `os.Symlink` পাথটিকে গ্রুপের সঙ্গে লিংক করে। অন্যান্য ক্ষেত্র:

- গ্রুপে এখনও কোনো কপি না থাকলে অ্যাকাউন্টের ফাইলটি গ্রুপে সরে যায় এবং তার প্রথম বিষয়বস্তু হয়।
- ফোল্ডারের ক্ষেত্রে গ্রুপে যে ফাইলগুলো নেই সেগুলো গ্রুপে সরে যায়; বাকিগুলো ব্যাকআপে যায়।
- `history.jsonl`-এর জন্য ডিফল্ট হলো `merge by line`: গ্রুপের লাইনগুলো, তারপর আপনার নতুনগুলো। গ্রুপের আগের কপি `.ccam-backup/<time>/group/`-এ সংরক্ষিত হয়।
- অন্য গ্রুপের লিংক প্রতিস্থাপন করা হয়। ccam-এর বাইরের কোনো ফোল্ডারের লিংক ব্যাকআপে যায়, এবং গ্রুপ খালি থাকলে তার বিষয়বস্তু দিয়ে গ্রুপ শুরু হয়।
- এখন লোকাল হয়ে যাওয়া আইটেমের জন্য রয়ে যাওয়া লিংক সরানো হয়, এবং অ্যাকাউন্টটি গ্রুপের বিষয়বস্তুর একটি কপি পায়।

লিংক সরানো ছাড়া প্রতিটি ধাপ ডেটা গ্রুপে বা ব্যাকআপ ফোল্ডারে সরায়। কিছুই মুছে ফেলা হয় না।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- সাধারণ `ccam doctor` কখনও কিছু বদলায় না। স্ক্রিপ্টে নিশ্চিত হতে `--fix`-এর সঙ্গে `--dry-run` ব্যবহার করুন।
- `doctor --fix` কোনো অ্যাকাউন্ট চলছে কি না যাচাই করে না। কোনো অ্যাকাউন্টের শেয়ার লিংক মেরামতের আগে তার Claude Code বন্ধ করুন, কারণ সেশনটি ওই ফাইলগুলোতে লিখতে থাকে।
- শেয়ার মেরামত ব্যর্থ হলে ccam ওই অ্যাকাউন্টের জন্য `could not repair share links: <error>` জানায় এবং তার অজানা এন্ট্রির যাচাই বাদ দেয়। অন্য অ্যাকাউন্টগুলো তবুও যাচাই করা হয়।
- `ccam share` মাঝপথে ব্যর্থ হলে এটি আগে লক্ষ্য সেটিংস সংরক্ষণ করে এবং আপনাকে `ccam doctor --fix` দিয়ে শেষ করতে বলে।
- সিমলিংক করা `.credentials.json`, মেয়াদোত্তীর্ণ লগইন ও এনভায়রনমেন্টের ক্রেডেনশিয়াল জানানো হয়, কিন্তু কখনও বদলানো হয় না।
- প্লাগইন রেজিস্ট্রি অ্যাটমিকভাবে আবার লেখা হয়, কিন্তু `.ccam-backup/`-এ কপি করা হয় না।

## সম্পর্কিত

- [MCP সার্ভার, সেটিংস ও ইতিহাস শেয়ার করুন](/docs/sharing)
- [আগের সেটআপ কপি করা](/features/migrate)
- [ব্যবহারের সীমা ও অ্যাকাউন্ট বদল](/docs/usage-limits)
- [কমান্ড](/docs/commands)
