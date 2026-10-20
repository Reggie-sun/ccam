---
title: "শর্টকাট কমান্ড"
description: "ccam shortcut দিয়ে একটি অ্যাকাউন্টকে ccwork-এর মতো নিজস্ব কমান্ড দিন, আর ccam shims install দিয়ে ~/.local/bin-এ লিংকগুলো আবার তৈরি করুন।"
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

শর্টকাট হলো এমন একটি কমান্ড যার নাম আপনি নিজে দেন এবং যা একটি অ্যাকাউন্ট হিসেবে Claude Code চালু করে: `ccwork --resume` ঠিক তাই করে যা `ccam work --resume` করে। এটি `~/.local/bin`-এ ccam-এর একটি লিংক, তাই এটি bash, zsh, স্ক্রিপ্ট এবং `PATH` থেকে প্রোগ্রাম চালায় এমন যেকোনো কিছুতে কাজ করে। শর্টকাট ঐচ্ছিক; যে অ্যাকাউন্টগুলো আপনি প্রায়ই চালু করেন সেগুলোর জন্য একটি তৈরি করুন।

## ব্যবহার

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam প্রতিটি পরিবর্তন নিশ্চিত করে:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

একটি অ্যাকাউন্টের সর্বোচ্চ একটি শর্টকাট থাকে, যা `~/.config/ccam/config.toml`-এ `shortcut` হিসেবে সংরক্ষিত থাকে। `ccam add` উইজার্ড জিজ্ঞেস করে `Create a shortcut command for <name>? [y/N]` এবং তারপর `Shortcut name [cc<name>]:`। উইজার্ডের বাইরে `ccam add` শুধু `--shortcut` দিলে একটি তৈরি করে। `ccam list` প্রতিটি অ্যাকাউন্টের শর্টকাট দেখায়, অথবা `-`।

শর্টকাট প্রতিটি আর্গুমেন্ট Claude Code-এ পাঠায়, ঠিক যেভাবে `ccam <name>` নামের পরে পাঠায়। এর মধ্যে `--dry-run` ও `--help`-ও আছে: চালু করার পূর্বরূপ দেখতে `ccam --dry-run work` ব্যবহার করুন।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--remove` | অ্যাকাউন্টের শর্টকাট ও তার লিংক মুছে দেয় | বন্ধ |
| `-h`, `--help` | কমান্ডের সাহায্য | |
| `--dry-run` | `ccam shortcut` এটি মানে: নামটি যাচাই করে এবং সংরক্ষণ বা লিংক না করে `Would set the shortcut of work to ccwork.` বা `Would remove ccwork, the shortcut of work.` প্রিন্ট করে | বন্ধ |

`ccam shims install` কোনো আর্গুমেন্ট নেয় না এবং `--dry-run` প্রত্যাখ্যান করে `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed `` দিয়ে।

নামের নিয়ম, এই ক্রমে যাচাই করা হয়:

- ১ থেকে ৬৪ অক্ষরের, `a-z`, `0-9`, `.`, `-` ও `_` দিয়ে গঠিত, কোনো অক্ষর বা অঙ্ক দিয়ে শুরু;
- `ccam`, `claude` বা `run` বা `usage`-এর মতো কোনো ccam কমান্ড নয়;
- অন্য কোনো অ্যাকাউন্টের শর্টকাট নয়;
- `~/.local/bin`-এ এমন কোনো ফাইল নয় যা আগে থেকেই ccam-এর লিংক নয়;
- আপনার `PATH`-এ থাকা অন্য কোনো এক্সিকিউটেবলের নাম নয়।

একসঙ্গে একটি কমান্ড ও `--remove` দিলে `give a command or --remove, not both` দিয়ে থেমে যায়। অ্যাকাউন্টের আগে থেকে থাকা নামটিই আবার দিলে `work: ccwork (unchanged)` প্রিন্ট হয় এবং লিংক না থাকলে আবার তৈরি হয়।

## এটি কীভাবে কাজ করে

নাম নির্ধারণ করলে `checkShortcut` চলে, যা bin ফোল্ডার ও ccam বাইনারির আসল পাথ দিয়ে `shim.Check` ডাকে:

<!-- motion:shortcuts-check -->

নামটি ক্রমানুসারে যাচাইগুলোর মধ্য দিয়ে নিচে নামে। প্রতিটি যাচাইয়ের পাশে এমন একটি নাম আছে যা সেটি প্রত্যাখ্যান করে: `Work` প্যাটার্নে ব্যর্থ হয়, `run` একটি ccam কমান্ড, `ccside` অন্য অ্যাকাউন্টের, `ccfile` `~/.local/bin`-এ একটি বিদ্যমান ফাইল, আর `git` `/usr/bin/git`-কে আড়াল করত। শুধু প্রতিটি যাচাই পেরোনো নামই সংরক্ষিত ও লিংক করা হয়।

1. `shim.Check` নামটি মেলায় `^[a-z0-9][a-z0-9._-]{0,63}$`-এর সঙ্গে, তারপর `config.Reserved`, তারপর `cfg.FindShortcut`। এটি `~/.local/bin/<name>`-এ `os.Lstat` ডাকে, এবং `PATH`-এর অন্য প্রতিটি ফোল্ডারের জন্য `os.Stat` ডাকে এবং ccam-এর লিংক নয় এমন এক্সিকিউটেবল ফাইল প্রত্যাখ্যান করে।
2. ccam `fsutil.AtomicWrite` দিয়ে `config.toml` সংরক্ষণ করে (অস্থায়ী ফাইল, `fsync`, rename)।
3. `shim.Remove` পুরোনো শর্টকাটের লিংক মুছে দেয়, তবে শুধু যদি সেটি ccam-এর লিংক হয়।
4. `shim.Install` কনফিগার করা প্রতিটি অনুপস্থিত শর্টকাটের লিংক তৈরি করে, `os.Executable` ও `filepath.EvalSymlinks` দিয়ে খুঁজে পাওয়া ccam বাইনারির দিকে `os.Symlink` দিয়ে। এটি কখনও কিছু সরায় না।

একটি লিংক ccam-এর নিজের বলে গণ্য হয় যখন `os.Readlink` ccam পাথ ফেরত দেয়, অথবা যখন লিংক ও ccam পাথ একই ফাইলে পৌঁছায়। প্রতিটি মুছে ফেলা ও প্রতিটি ওভাররাইট যাচাই এই পরীক্ষার মধ্য দিয়ে যায়, তাই ccam নিজের তৈরি না করা ফাইল স্পর্শ করে না।

শর্টকাট চালালে আপনার টাইপ করা নামটিই ccam-এর একমাত্র প্রয়োজন:

<!-- motion:shortcuts-link -->

শেল `PATH`-এ `ccwork` খুঁজে পায়, কার্নেল লিংক অনুসরণ করে ccam বাইনারি চালু করে, এবং `argv[0]` তখনও `ccwork`। `app.Main`-এ `ccam` ছাড়া যেকোনো নাম একবার `config.toml` লোড করে এবং `shim.Lookup` ডাকে, যা সেই অ্যাকাউন্ট খুঁজে বের করে যার `shortcut` `filepath.Base(argv[0])`-এর সমান। তারপর `launcher.LaunchWith` আগে থেকে লোড হওয়া কনফিগ দিয়ে [`ccam <name>`](/features/launch)-এর একই ধাপগুলো নেয়, এবং `argv[1:]` অপরিবর্তিত অবস্থায় `claude`-এ দেয়।

কেন শেল অ্যালিয়াস বা র‍্যাপার স্ক্রিপ্ট নয়, সিমলিংক: সিমলিংক প্রতিটি শেলে ও স্ক্রিপ্টে কাজ করে, এটি অতিরিক্ত কোনো শেল চালু যোগ করে না (প্রায় ১ থেকে ২ ms), এবং চালু করা তখনও একটি ফাইল পড়ে। `scripts/perf-check.sh` `hyperfine` দিয়ে একটি শর্টকাট মাপে, `ccam <name>`-এর একই ১০ ms (p95) সীমার বিপরীতে।

### Windows-এ

শর্টকাট হলো `%USERPROFILE%\.local\bin`-এ `ccwork.exe`, এবং খোঁজার আগে `shim.Name` `.exe` বাদ দেয় (যেকোনো কেসে), তাই `ccwork`, `ccwork.exe` ও `CCWORK.EXE` সবই `work` চালু করে। ccam প্রথমে একটি সিম্বলিক লিংক তৈরি করে, যা Windows Developer Mode-এ বা একজন অ্যাডমিনিস্ট্রেটরের জন্য অনুমতি দেয়। সেই অধিকার না থাকলে এটি একটি হার্ড লিংক তৈরি করে: একই `ccam.exe` ফাইলের দ্বিতীয় একটি নাম, যার জন্য কোনো বিশেষ অধিকার লাগে না। একটি হার্ড লিংক ccam-এর নিজের বলে গণ্য হয় যখন সেটি `ccam.exe`-এর সঙ্গে একই ফাইল। যেহেতু `ccam update` `ccam.exe`-এর জায়গায় একটি নতুন ফাইল বসায়, তারপর এটি `shim.Relink` দিয়ে প্রতিটি হার্ড লিংক করা শর্টকাট নবায়ন করে এবং `✓ shortcut ccwork renewed` প্রিন্ট করে। নাম যাচাইয়ের সময় অন্য `PATH` ফোল্ডারগুলোতে `PATHEXT`-এর প্রতিটি এক্সটেনশন (`.exe`, `.cmd` ইত্যাদি) সহ `<name>` খোঁজা হয়।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- **যে লিংক কোনো অ্যাকাউন্টের নয়।** `~/.local/bin`-এ ccam-এর কোনো লিংকের নাম যদি কোনো অ্যাকাউন্ট ব্যবহার না করে, যেমন আপনি `config.toml`-এ হাতে `shortcut` সম্পাদনা করার পরে, ccam ডিফল্ট অ্যাকাউন্টে ফিরে যায় না। এটি `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) `` দিয়ে থেমে যায়। `ccam shortcut` দিয়ে আবার নামটি নির্ধারণ করুন, অথবা লিংকটি মুছে দিন।
- **বিদ্যমান ফাইল।** ccam কখনও নিজের তৈরি না করা ফাইল প্রতিস্থাপন করে না। `ccam shims install` থেমে যায় `<path> already exists and is not a ccam link, so it was left alone` দিয়ে।
- **সরানো ccam বাইনারি।** তখন পুরোনো লিংকগুলো এমন পাথে নির্দেশ করে যা আর নেই, এবং ccam সেগুলোকে নিজের বলে গণ্য করে না। সেগুলো মুছে দিন, তারপর `ccam shims install` চালান।
- **নাম বদল ও সরানো।** `ccam rename` শর্টকাটের নাম একই রাখে, এবং লিংকটি কাজ করতে থাকে, কারণ এটি ccam-এ নির্দেশ করে আর ccam `config.toml`-এ নামটি খোঁজে। `ccam remove` অ্যাকাউন্টের লিংক মুছে দেয়। `ccam shortcut <account> <new>` নতুনটি তৈরির আগে পুরোনো লিংক মুছে দেয়।
- **PATH।** ccam লিংক তৈরি করার সময় `~/.local/bin` আপনার `PATH`-এ না থাকলে এটি প্রিন্ট করে `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` ও তারপর যোগ করার মতো `export PATH=...` লাইনটি। Windows-এ এর বদলে এটি সেই PowerShell লাইনটি প্রিন্ট করে যা ফোল্ডারটি আপনার ইউজার `PATH`-এ যোগ করে।
- **পরে যাচাই।** `ccam doctor` জানায় `missing shortcuts in /home/me/.local/bin: ccwork` এবং ফোল্ডারটি `PATH`-এ না থাকলে সতর্ক করে। `ccam doctor --fix` অনুপস্থিত লিংকগুলো তৈরি করে।

## সম্পর্কিত

- [Claude Code চালু করা](/features/launch)
- [অ্যাকাউন্ট পরিচালনা](/features/accounts)
- [কমান্ড রেফারেন্স](/docs/commands#shortcuts)
- [ccam ইনস্টল করুন](/docs/install)
