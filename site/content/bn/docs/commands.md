---
title: "কমান্ড রেফারেন্স"
metaTitle: "ccam কমান্ড রেফারেন্স"
description: "ccam add ও ccam <name> থেকে ccam usage, ccam next ও ccam doctor পর্যন্ত ccam-এর প্রতিটি কমান্ড ও অপশন।"
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 116c5b34ac9a
---

কোনো কমান্ড আসলে না চালিয়ে সেটি কী করবে তা দেখতে কমান্ডের আগে `--dry-run` দিন, যেমন `ccam --dry-run work`। এটি Claude Code চালু করার সময় (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) এবং `group`, `migrate`, `shortcut`, `update` ও `doctor --fix`-এর সঙ্গে কাজ করে। কিছু বদলায় এমন অন্য কমান্ডগুলো চলার বদলে এটি প্রত্যাখ্যান করে।

<!-- commands -->

## Claude Code চালু করা

`ccam <name|number> [arguments]` ওই অ্যাকাউন্ট হিসেবে Claude Code চালু করে। নামের পরের আর্গুমেন্টগুলো ccam পড়ে না; সেগুলো একই ক্রমে `claude`-এ পাঠিয়ে দেয়, এবং স্ট্যান্ডার্ড ইনপুট, এক্সিট কোড ও সিগন্যাল ঠিক `claude`-এর মতোই আচরণ করে। নামের পরে `--` দেওয়া ঐচ্ছিক।

অ্যাকাউন্টের নাম ছাড়া কিন্তু আর্গুমেন্টসহ `ccam`, যেমন `ccam -p "summarise"`, বর্তমান ফোল্ডারের জন্য নির্ধারিত অ্যাকাউন্ট ব্যবহার করে, না থাকলে ডিফল্ট অ্যাকাউন্ট। এগুলো `ccam folder <name> [folder]` ও `ccam use <name>` দিয়ে নির্ধারণ করুন।

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

অ্যাকাউন্টের নামের আগে `--fallback` হলো `-p`/`--print` রানের জন্য। অ্যাকাউন্টটি ব্যবহারের সীমায় পৌঁছালে ccam একই কমান্ড সবচেয়ে বেশি জায়গা থাকা লগইন করা অ্যাকাউন্ট হিসেবে আবার চালায় এবং শুধু যে রানটি শেষ হয়েছে তার আউটপুট প্রিন্ট করে। [ফোল্ডার অনুযায়ী অ্যাকাউন্ট](/features/folder) ও [সীমায় আবার চালানো](/features/next#run-again-at-the-limit) দেখুন।

`ccwork`-এর মতো একটি শর্টকাট তার অ্যাকাউন্টকে একইভাবে চালু করে: `ccwork --resume` ঠিক তাই করে যা `ccam work --resume` করে। [শর্টকাট](/docs/commands#shortcuts) দেখুন।

## অ্যাকাউন্ট যোগ করা ও সরানো

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

ফোল্ডার না দিয়ে `ccam import` চালালে এটি `~/.claude` ও `~/.claude-*` খোঁজে এবং প্রতিটির বিষয়ে আপনাকে জিজ্ঞেস করে। ইমপোর্ট করা ফোল্ডার যেখানে আছে সেখানেই নিবন্ধিত হয়; ccam কখনও কোনো লগইন সরায় বা কপি করে না।

`ccam add --api` এমন একটি অ্যাকাউন্ট যোগ করে যা Claude লগইনের বদলে API key ব্যবহার করে; ccam প্রতিবার চালু করার সময় key পড়ে এবং কখনও সংরক্ষণ করে না। দেখুন [API key অ্যাকাউন্ট](/features/api)।

কোনো অ্যাকাউন্ট হিসেবে Claude Code চলতে থাকলে `ccam rename`, `ccam remove` ও `ccam group` কমান্ডগুলো সেই অ্যাকাউন্ট বদলাতে রাজি হয় না। আগে সেশনটি বন্ধ করুন।

`ccam list --json` প্রতিটি অ্যাকাউন্ট তার স্লট নম্বর, শর্টকাট, ইমেইল ও ফোল্ডারসহ দেখায়। এর `groups` অ্যারে অ্যাকাউন্টটির গ্রুপগুলোর তালিকা দেয়, আর এর `share` অবজেক্ট প্রতিটি আইটেমকে (`mcp`, `settings`, `memory`, `extensions`, `history`) সেই গ্রুপের সঙ্গে মেলায় যেখান থেকে আইটেমটি আসে, অথবা `local`-এর সঙ্গে।

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

একটি গ্রুপ তার অ্যাকাউন্টগুলোর মধ্যে আইটেম শেয়ার করে; একটি অ্যাকাউন্ট একাধিক গ্রুপে থাকতে পারে, তবে প্রতিটি আইটেম আসে কেবল একটি গ্রুপ থেকে। দেখুন [গ্রুপ দিয়ে শেয়ার করা](/features/share)।

## আগের সেটআপ কপি করা

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` একটি Claude Code ফোল্ডার থেকে `settings.json`, `CLAUDE.md`, এক্সটেনশন ও ইউজার MCP সার্ভার একটি গ্রুপে কপি করে, যাতে ওই আইটেমগুলো শেয়ার করা প্রতিটি অ্যাকাউন্ট সেগুলো পায়। ফোল্ডারটি শুধু পড়া হয়। `--settings`, `--memory`, `--extensions`, `--mcp` ও `--history` অপশন শুধু আপনার উল্লেখ করা আইটেমগুলো কপি করে। এগুলো ছাড়া টার্মিনালে একটি চেকলিস্ট থেকে আপনি বেছে নিতে পারেন, যেখানে কথোপকথনের ইতিহাস বাদে সবকিছু চেক করা থাকে; অন্য ক্ষেত্রে কথোপকথনের ইতিহাস বাদে সবকিছু কপি হয়। `-g` বা `--group` গ্রুপ বেছে নেয়, আর `-y` বা `--yes` জিজ্ঞেস না করেই কপি করে এবং কোনো ফাইল আলাদা হলে গ্রুপের কপি রাখে। এটি সার্ভারগুলোর সংরক্ষিত সাইন-ইনও প্রতিটি অ্যাকাউন্টের নিজস্ব লগইন ফাইলে কপি করে। [আগের সেটআপ কপি করা](/docs/sharing#copy-an-existing-setup) দেখুন।

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
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` ও `ccam pick` একইভাবে বেছে নেয়: প্রথমে `auto_max_usage`-এর (`config.toml`, ডিফল্ট 100) নিচে থাকা অ্যাকাউন্ট, `auto_strategy` (`headroom`, `expiring` বা `balance`) অনুযায়ী, `--model`-এর সাপ্তাহিক উইন্ডো ধরে, ১০০%-এ থাকা অ্যাকাউন্ট কখনও নয়, এবং যে অ্যাকাউন্ট থেকে শুরু হয় তার প্রতিষ্ঠানের ভেতরে (আরও দূরে খুঁজতে `--cross-org` বা `cross_org = true`)। কিছুই বেছে নেওয়া না গেলে এগুলো কোড 75 দিয়ে বেরিয়ে যায়, অথবা `--wait` দিলে প্রথম অ্যাকাউন্টটি রিসেট হওয়া পর্যন্ত অপেক্ষা করে। `ccam next --continue` ইতিহাস শেয়ার করে না এমন অ্যাকাউন্টগুলোর মধ্যেও কথোপকথন নিয়ে যায়: তখন Claude Code সেটি নতুন অ্যাকাউন্টে কপি করে। Claude Code-এর ভেতরে প্লাগইনের `/ccam:switch [account]` কথোপকথন অন্য অ্যাকাউন্টের হাতে তুলে দেয়, আর এর হুকগুলো যা চালায় তা হলো `ccam hook`। দেখুন [পরের অ্যাকাউন্টে যাওয়া](/features/next), [স্ক্রিপ্টের জন্য অ্যাকাউন্ট বেছে নেওয়া](/features/pick) ও [প্রতিটি অ্যাকাউন্টের ব্যাকগ্রাউন্ড সেশন](/features/agents)।

## রক্ষণাবেক্ষণ

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

ccam বদলানোর আগে `ccam update` রিলিজের SHA-256 চেকসাম যাচাই করে, আর `--version v0.2.0` একটি নির্দিষ্ট রিলিজ ইনস্টল করে। `.deb` প্যাকেজ থেকে ইনস্টল করা ccam এর বদলে apt দিয়ে আপডেট করুন। [আপডেট](/docs/install#update) দেখুন।

## ফাইল

| পাথ | বিষয়বস্তু |
| --- | --- |
| `~/.config/ccam/config.toml` | অ্যাকাউন্ট, গ্রুপ ও সেটিংস (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | প্রতিটি অ্যাকাউন্টের `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | শেয়ার করা সেটিংস, MCP সার্ভার, এক্সটেনশন ও ইতিহাস |
| `~/.local/state/ccam/` | ব্যবহারের ক্যাশ, সর্বশেষ ব্যবহারের সময়, এবং `/ccam:switch` ও সীমার হুক যে নোট রেখে যায় (`handoff/`) |
| `~/.local/bin/<shortcut>` | `ccwork`-এর মতো শর্টকাট কমান্ড, ccam-এর লিংক |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ও `XDG_STATE_HOME` ভেরিয়েবল এই ফোল্ডারগুলোর অবস্থান বদলায়। Windows-এ `~` হলো `%USERPROFILE%` এবং শর্টকাটগুলো `<shortcut>.exe`।
