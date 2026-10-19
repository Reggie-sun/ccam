---
title: "কমান্ড রেফারেন্স"
description: "ccam add ও ccam <name> থেকে ccam usage, ccam next ও ccam doctor পর্যন্ত ccam-এর প্রতিটি কমান্ড ও অপশন।"
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 8ac012f7f7cf
---

কোনো কমান্ড আসলে না চালিয়ে সেটি কী করবে তা দেখতে কমান্ডের আগে `--dry-run` দিন, যেমন `ccam --dry-run work`। এটি Claude Code চালু করার সময় (`ccam <name>`, `ccam run`, `ccam next`) এবং `share`, `migrate`, `shortcut`, `update` ও `doctor --fix`-এর সঙ্গে কাজ করে। কিছু বদলায় এমন অন্য কমান্ডগুলো চলার বদলে এটি প্রত্যাখ্যান করে।

<!-- commands -->

## Claude Code চালু করা

`ccam <name|number> [arguments]` ওই অ্যাকাউন্ট হিসেবে Claude Code চালু করে। নামের পরের আর্গুমেন্টগুলো ccam পড়ে না; সেগুলো একই ক্রমে `claude`-এ পাঠিয়ে দেয়, এবং স্ট্যান্ডার্ড ইনপুট, এক্সিট কোড ও সিগন্যাল ঠিক `claude`-এর মতোই আচরণ করে। নামের পরে `--` দেওয়া ঐচ্ছিক।

অ্যাকাউন্টের নাম ছাড়া কিন্তু আর্গুমেন্টসহ `ccam`, যেমন `ccam -p "summarise"`, ডিফল্ট অ্যাকাউন্ট ব্যবহার করে। এটি `ccam use <name>` দিয়ে নির্ধারণ করুন।

`ccwork`-এর মতো একটি শর্টকাট তার অ্যাকাউন্টকে একইভাবে চালু করে: `ccwork --resume` ঠিক তাই করে যা `ccam work --resume` করে। [শর্টকাট](/docs/commands#shortcuts) দেখুন।

## অ্যাকাউন্ট যোগ করা ও সরানো

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

ফোল্ডার না দিয়ে `ccam import` চালালে এটি `~/.claude` ও `~/.claude-*` খোঁজে এবং প্রতিটির বিষয়ে আপনাকে জিজ্ঞেস করে। ইমপোর্ট করা ফোল্ডার যেখানে আছে সেখানেই নিবন্ধিত হয়; ccam কখনও কোনো লগইন সরায় বা কপি করে না।

কোনো অ্যাকাউন্ট হিসেবে Claude Code চলতে থাকলে `ccam rename`, `ccam remove` ও `ccam share` সেই অ্যাকাউন্ট বদলাতে রাজি হয় না। আগে সেশনটি বন্ধ করুন।

`ccam list --json` প্রতিটি অ্যাকাউন্ট তার স্লট নম্বর, শর্টকাট, ইমেইল ও ফোল্ডারসহ দেখায়। এর `share` অবজেক্ট প্রতিটি আইটেমকে (`mcp`, `settings`, `memory`, `extensions`, `history`) একটি গ্রুপের নাম বা `local`-এর সঙ্গে মেলায়।

## শর্টকাট

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

শর্টকাট হলো এমন একটি কমান্ড যার নাম আপনি নিজে দেন এবং যা একটি অ্যাকাউন্ট হিসেবে Claude Code চালু করে। `ccam work`-এর মতোই এটি প্রতিটি আর্গুমেন্ট অপরিবর্তিত অবস্থায় Claude Code-এ পাঠায়, তাই `ccwork --resume` ও `ccwork -p "hi"` প্রত্যাশামতো কাজ করে। এটি `~/.local/bin`-এ ccam-এর একটি লিংক।

শর্টকাট ঐচ্ছিক। `ccam add` উইজার্ড শর্টকাট তৈরি করবে কি না জিজ্ঞেস করে, `ccam add <name> --shortcut <command>` কিছু জিজ্ঞেস না করেই একটি তৈরি করে, আর `ccam list` প্রতিটি অ্যাকাউন্টের শর্টকাট দেখায়, শর্টকাট না থাকলে `-`।

নামে a-z, 0-9, `.`, `-` ও `_` অক্ষর ব্যবহার করা যায়। নামটি কোনো ccam কমান্ড, অন্য অ্যাকাউন্টের শর্টকাট, বা আপনার `PATH`-এ থাকা অন্য কোনো প্রোগ্রামের নাম হতে পারে না।

`ccam rename` শর্টকাটের নাম একই রাখে, আর `ccam remove` অ্যাকাউন্টের সঙ্গে তার শর্টকাটও মুছে দেয়। `ccam shims install` সব অ্যাকাউন্টের লিংক আবার তৈরি করে, যেমন ccam অন্য জায়গায় সরানোর পরে।

## শেয়ার করা

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

প্রতিটি আইটেম `on` (`default` গ্রুপ), `off`, অথবা একটি গ্রুপের নাম নেয়। দেখুন [শেয়ার করা](/docs/sharing)।

## আগের সেটআপ কপি করা

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` একটি Claude Code ফোল্ডার থেকে `settings.json`, `CLAUDE.md`, এক্সটেনশন ও ইউজার MCP সার্ভার একটি গ্রুপে কপি করে, যাতে ওই আইটেমগুলো শেয়ার করা প্রতিটি অ্যাকাউন্ট সেগুলো পায়। ফোল্ডারটি শুধু পড়া হয়। `--settings`, `--memory`, `--extensions`, `--mcp` ও `--history` অপশন শুধু আপনার উল্লেখ করা আইটেমগুলো কপি করে; এগুলো ছাড়া কথোপকথনের ইতিহাস বাদে সবকিছু কপি হয়। `-g` বা `--group` গ্রুপ বেছে নেয়, আর `-y` বা `--yes` জিজ্ঞেস না করেই কপি করে এবং কোনো ফাইল আলাদা হলে গ্রুপের কপি রাখে। [আগের সেটআপ কপি করা](/docs/sharing#copy-an-existing-setup) দেখুন।

## MCP সার্ভার

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` শুধু ওই অ্যাকাউন্টে থাকা সার্ভারগুলোকে তার গ্রুপে সরায়। যে সার্ভারের নাম গ্রুপ ইতিমধ্যে অন্য কোনো সংজ্ঞার জন্য ব্যবহার করছে, সেটি বাদ দেওয়া হয়; গ্রুপের কপি বদলাতে `--force` যোগ করুন, ccam আগে সেটির ব্যাকআপ রাখে।

## ব্যবহার ও অ্যাকাউন্ট বদল

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## রক্ষণাবেক্ষণ

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
```

ccam বদলানোর আগে `ccam update` রিলিজের SHA-256 চেকসাম যাচাই করে, আর `--version v0.2.0` একটি নির্দিষ্ট রিলিজ ইনস্টল করে। `.deb` প্যাকেজ থেকে ইনস্টল করা ccam এর বদলে apt দিয়ে আপডেট করুন। [আপডেট](/docs/install#update) দেখুন।

## ফাইল

| পাথ | বিষয়বস্তু |
| --- | --- |
| `~/.config/ccam/config.toml` | অ্যাকাউন্ট এবং সেগুলো কী শেয়ার করে |
| `~/.local/share/ccam/profiles/<name>/` | প্রতিটি অ্যাকাউন্টের `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | শেয়ার করা সেটিংস, MCP সার্ভার, এক্সটেনশন ও ইতিহাস |
| `~/.local/state/ccam/` | ব্যবহারের ক্যাশ ও সর্বশেষ ব্যবহারের সময় |
| `~/.local/bin/<shortcut>` | `ccwork`-এর মতো শর্টকাট কমান্ড, ccam-এর লিংক |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ও `XDG_STATE_HOME` ভেরিয়েবল এই ফোল্ডারগুলোর অবস্থান বদলায়।
