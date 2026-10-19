---
title: "অ্যাকাউন্ট পরিচালনা"
description: "ccam add, ccam login, ccam list, ccam rename, ccam remove ও ccam import দিয়ে Claude Code অ্যাকাউন্ট যোগ করুন, লগইন করুন, তালিকা দেখুন, নাম বদলান, সরান ও ইমপোর্ট করুন।"
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 97cef50f48f4
---

প্রতিটি ccam অ্যাকাউন্ট নিজস্ব লগইনসহ আলাদা একটি Claude Code কনফিগ ফোল্ডার। এই কমান্ডগুলো অ্যাকাউন্ট তৈরি করে, সাইন ইন করায়, তালিকা দেখায়, নাম বদলায় ও সরায়, এবং আপনার আগে থেকে থাকা ফোল্ডার নিবন্ধন করে।

## ব্যবহার

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

যেখানে কোনো কমান্ড `<name|number>` নেয়, সেখানে অ্যাকাউন্টের নাম বা তার স্লট নম্বর দিতে পারেন, অর্থাৎ `ccam list`-এর `#` কলাম।

টার্মিনালে কোনো অপশন ছাড়া `ccam add` একটি উইজার্ড খোলে। এটি নাম জিজ্ঞেস করে, তারপর শেয়ার করা আইটেমগুলোর একটি চেকলিস্ট দেখায় (কথোপকথনের ইতিহাস বাদে সবকিছু চেক করা থাকে), তারপর জিজ্ঞেস করে `Create a shortcut command for work? [y/N]` এবং `y` দিলে `Shortcut name [ccwork]:`। `--share`, `--shortcut` বা `--no-login`-এর যেকোনোটি দিলে, অথবা টার্মিনালের বদলে পাইপ থাকলে, উইজার্ড বাদ যায়; তখন নাম দেওয়া বাধ্যতামূলক, এবং `--shortcut` কোনো নাম না দিলে শর্টকাট তৈরি হয় না।

`ccam list` স্লট, নাম, শর্টকাট, ইমেইল, শেয়ার করা আইটেম ও সর্বশেষ ব্যবহার দেখায়, ডিফল্ট অ্যাকাউন্টের আগে `*` থাকে:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

`--json` দিলে প্রতিটি অ্যাকাউন্ট একটি অবজেক্ট, যাতে থাকে `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (অ্যাকাউন্টটি কখনও চালু না হলে বাদ থাকে) ও `share`। `share` প্রতিটি আইটেমকে (`mcp`, `settings`, `memory`, `extensions`, `history`) তার গ্রুপের নাম বা `"local"`-এর সঙ্গে মেলায়।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `ccam add --no-login` | লগইন না করে অ্যাকাউন্ট তৈরি করে; পরে `ccam login <name>` দিয়ে লগইন করুন | বন্ধ |
| `ccam add --share <item=group>` | একটি আইটেম কোনো গ্রুপের সঙ্গে শেয়ার করে, অথবা `item=local` দিয়ে সেটি নিজস্ব রাখে। একাধিকবার দেওয়া যায় | `mcp`, `settings`, `memory`, `extensions` থাকে `default`-এ; `history` লোকাল |
| `ccam add --shortcut <command>` | `~/.local/bin`-এ এই শর্টকাট কমান্ডটিও তৈরি করে | নেই |
| `ccam list --json` | অ্যাকাউন্টগুলো JSON হিসেবে প্রিন্ট করে | বন্ধ |
| `ccam remove -y`, `--yes` | জিজ্ঞেস না করেই সরায় | বন্ধ |
| `ccam import --name <name>` | আপনার দেওয়া ফোল্ডারের জন্য অ্যাকাউন্টের নাম | ফোল্ডারের নাম থেকে নেওয়া |
| `ccam import -y`, `--yes` | পাওয়া প্রতিটি ফোল্ডার জিজ্ঞেস না করেই ইমপোর্ট করে | বন্ধ |
| `--dry-run` | `ccam list` ও `ccam completion` এটি মানে, কারণ এগুলো কিছু বদলায় না। `add`, `login`, `logout`, `use`, `rename`, `remove` ও `import` এটি প্রত্যাখ্যান করে, `` --dry-run is not supported by `ccam add`; nothing was changed `` বার্তাসহ (কমান্ডের নাম উল্লেখ করে) | বন্ধ |

`ccam login`, `ccam logout`, `ccam use` ও `ccam rename`-এর নিজস্ব কোনো অপশন নেই। `ccam completion` নেয় `bash` বা `zsh`; অন্য কিছু দিলে `only bash and zsh are supported` দিয়ে ব্যর্থ হয়।

অ্যাকাউন্টের নাম ১ থেকে ৩২ অক্ষরের, `a-z`, `0-9`, `-` ও `_` দিয়ে গঠিত, এবং কোনো অক্ষর বা অঙ্ক দিয়ে শুরু হয়। নাম শুধু অঙ্ক দিয়ে হতে পারে না, কারণ `ccam 2` মানে স্লট 2। নামটি কোনো ccam কমান্ডও হতে পারে না (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), তাই `ccam <name>` কখনও দ্ব্যর্থক হয় না। তালিকাটি `internal/config`-এর `config.Reserved`।

## এটি কীভাবে কাজ করে

`ccam add work --shortcut ccwork` `addProfile`-এ (`internal/cli/profiles.go`) এই ধাপগুলো চালায়:

<!-- motion:accounts-add -->

1. কিছু তৈরি করার আগে `config.ValidateName` নাম যাচাই করে, আর `shim.Check` শর্টকাট যাচাই করে।
2. `profile.Add` মোড `0700` দিয়ে `~/.local/share/ccam/profiles/work/` তৈরি করে এবং অ্যাকাউন্টটিকে সবচেয়ে ছোট খালি স্লট নম্বর দেয় (`Config.NextSlot`), তাই সরানো অ্যাকাউন্টের নম্বর আবার ব্যবহার হয়। প্রথম অ্যাকাউন্টটি ডিফল্ট হয়।
3. `applyShare` নতুন, খালি ফোল্ডারটিকে `share.Build` ও তারপর `Apply` দিয়ে তার শেয়ার গ্রুপগুলোর সঙ্গে লিংক করে, যেমন `settings.json`-কে `groups/default/settings.json`-এর সঙ্গে। এটি ব্যর্থ হলে নতুন ফোল্ডারটি সরিয়ে দেওয়া হয়, কারণ রয়ে যাওয়া ফোল্ডার প্রতিটি পুনরায় চেষ্টাকে `folder already exists` দিয়ে আটকে দিত।
4. অ্যাকাউন্টটি `fsutil.AtomicWrite`-এর মাধ্যমে `~/.config/ccam/config.toml`-এ সংরক্ষিত হয়: একটি অস্থায়ী ফাইল, `fsync`, তারপর `rename`, তাই ক্র্যাশ হলেও কখনও অর্ধেক কনফিগ থেকে যায় না।
5. `runLogin` অ্যাকাউন্ট ফোল্ডারে `CLAUDE_CONFIG_DIR` নির্ধারণ করে `claude auth login` চালু করে। লগইন ব্যর্থ হলে অ্যাকাউন্টটি থেকে যায়, এবং ccam ত্রুটিটি প্রিন্ট করে আপনাকে পরে `ccam login work` চালাতে বলে।
6. `profile.ReadEmail` `profiles/work/.claude.json` থেকে `oauthAccount.emailAddress` বের করে। ওই ফাইল কয়েক মেগাবাইট হতে পারে, তাই পুরো ফাইল পার্স করার বদলে `TopLevelValue` শুধু ওই একটি কী খোঁজে।
7. `shim.Install` `~/.local/bin/ccwork`-কে ccam বাইনারির সঙ্গে লিংক করে, যা চালু হওয়ার সময় `config.toml`-এ `argv[0]` খোঁজে।

একটি লগইনে কয়েক মিনিট লাগতে পারে, আর এর মধ্যে অন্য কোনো ccam কনফিগ সংরক্ষণ করতে পারে। তাই যে কমান্ডগুলো অপেক্ষা করে, সেগুলো `env.update` দিয়ে সংরক্ষণ করে, যা `config.toml` আবার লোড করে, শুধু ওই একটি পরিবর্তন (এখানে ইমেইল) প্রয়োগ করে এবং সংরক্ষণ করে। একই কারণে `ccam remove` প্রশ্নের পরে আবার লোড করে।

`ccam login` একটি বিদ্যমান অ্যাকাউন্টের জন্য ধাপ ৫ ও ৬ চালায়। `ccam logout` অ্যাকাউন্টে `claude auth logout` চালায় এবং সংরক্ষিত ইমেইল মুছে দেয়। `ccam use` শুধু `config.toml`-এ `default` আবার লেখে।

### দূরের মেশিনে সাইন ইন

SSH-এর মাধ্যমে `claude auth login`-এর ব্রাউজার অংশটি ভেঙে যেতে পারে। Claude Code সার্ভারে `http://localhost:<port>/callback`-এ একটি কলব্যাক সার্ভারে সাইন-ইনের জন্য অপেক্ষা করে, কিন্তু যে ব্রাউজার সাইন-ইন শেষ করে সেটি চলে আপনার কম্পিউটারে, যেখানে ওই ঠিকানা কোথাও পৌঁছায় না এবং ব্রাউজার `ERR_CONNECTION_REFUSED`-এর মতো ত্রুটি দেখায়। ঠিকানাটিতে তবুও এককালীন কোডটি থাকে, তাই ccam সেটি পৌঁছে দিতে পারে:

<!-- motion:accounts-relay -->

`runLogin` আপনার টার্মিনাল Claude Code-কে দিয়ে দেয় না। এটি Claude Code-এর স্ট্যান্ডার্ড ইনপুট একটি পাইপের সঙ্গে যুক্ত করে, এবং `relayLogin` আপনার টাইপ করা প্রতিটি লাইন তাতে কপি করে, তাই Claude Code-এর `Paste code here` প্রম্পটে পেস্ট করা কোডও পৌঁছায়। কোনো লাইন কলব্যাক ঠিকানা হলে `callbackURL` সেটি গ্রহণ করে শুধু তখনই, যখন হোস্ট `localhost`, `127.0.0.1` বা `[::1]`, একটি পোর্ট দেওয়া আছে, পাথ `/callback`, এবং `code` ও `state` দুটোই আছে; `http://` বাদ দেওয়া যায়। তারপর `deliverCallback` সার্ভারে ১০ সেকেন্ডের সীমা দিয়ে এবং রিডাইরেক্ট অনুসরণ না করে ওই ঠিকানায় অনুরোধ পাঠায়, ঠিক যেভাবে ব্রাউজার পাঠাত। অনুরোধটি অপেক্ষারত `claude auth login`-এর কলব্যাক সার্ভারে পৌঁছায় (ডায়াগ্রামের ড্যাশ দেওয়া রেখা)। Claude Code কোডটি পায় এবং লগইন সংরক্ষণ করে, Linux-এ অ্যাকাউন্ট ফোল্ডারের `.credentials.json`-এ। পেস্ট করা কোনো লাইন কখনও ccam-কে অন্য কোনো হোস্টে যোগাযোগ করাতে পারে না।

লগইন শুরুর আগে ccam এটি প্রিন্ট করে:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

সফল হলে এটি প্রিন্ট করে `→ Passed the sign-in to Claude Code on this machine.`

### নাম বদল, সরানো ও ইমপোর্ট

`ccam rename` ও `ccam remove` অ্যাকাউন্টের ফোল্ডারের ওপর এবং অ্যাকাউন্টের নামে ccam অন্য জায়গায় যা রাখে তার ওপর কাজ করে। `nameState` ওই পাথগুলোর তালিকা দেয়: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` এবং ব্যবহারের ক্যাশ `~/.local/state/ccam/usage/<name>.json` ও `<name>.statusline.json`। `ccam remove work` এগুলো এই ক্রমে সামলায়:

<!-- motion:accounts-remove -->

`ccam remove` প্রথমে অ্যাকাউন্টটি চলতে থাকলে প্রত্যাখ্যান করে, তারপর জিজ্ঞেস করে। অ্যাকাউন্টটি এক্সটেনশন শেয়ার করলে `RewritePluginPaths` গ্রুপের `installed_plugins.json` ও `known_marketplaces.json` আবার লেখে, কারণ Claude Code প্লাগইনের পূর্ণ পাথ রাখে যা এই অ্যাকাউন্টের ফোল্ডারের মধ্য দিয়ে যেতে পারে। `profile.Trash` ফোল্ডারটি `~/.local/share/ccam/trash/<name>-<time>/`-এ সরায়, অ্যাকাউন্টটি `config.toml` থেকে বাদ যায়, `nameState` পাথগুলো মুছে ফেলা হয় যাতে পরে একই নামের অ্যাকাউন্ট সেগুলো উত্তরাধিকারে না পায়, এবং শর্টকাট লিংক সরানো হয়। ফিরিয়ে আনতে ফোল্ডারটি ট্র্যাশ থেকে আবার সরিয়ে আনুন।

`ccam rename` ফোল্ডারের নাম বদলায়, এর `plugins/`-এ থাকা প্লাগইন পাথগুলো পুরোনো ফোল্ডার থেকে নতুনটিতে আবার লেখে, ডিফল্ট অবস্থা সঙ্গে নিয়ে যায়, এবং `nameState` পাথগুলো নতুন নামে সরায়। শর্টকাট কমান্ডের নাম একই থাকে। প্লাগইন আবার লেখা ব্যর্থ হলে ফোল্ডারের নাম আগের মতো ফিরিয়ে দেওয়া হয়, কারণ কনফিগ সংরক্ষিত হয়নি।

`ccam import` একটি ফোল্ডার যেখানে আছে সেখানেই নিবন্ধন করে, `config.toml`-এ `dir` হিসেবে। এটি কখনও ফোল্ডার সরায় বা কপি করে না: কপি করা রিফ্রেশ টোকেন কপিগুলোকে একে অপরকে লগ আউট করিয়ে দিত। ফোল্ডার আর্গুমেন্ট ছাড়া এটি অনিবন্ধিত `~/.claude` ও `~/.claude-*` ফোল্ডার খোঁজে এবং প্রতিটির বিষয়ে জিজ্ঞেস করে; নাম আসে ফোল্ডার থেকে (`~/.claude` হলো `main`, `~/.claude-work` হলো `work`, `~/.claude-2` হলো `acct2`), আর যে ফোল্ডারের নাম বের করা যায় না সেটি ব্যবহার করার মতো `--name` কমান্ডসহ তালিকায় দেখানো হয়। `CLAUDE_CONFIG_DIR` ছাড়া Claude Code `.claude.json` আপনার হোম ফোল্ডারে `~/.claude`-এর পাশে রাখে, তাই আপনি `~/.claude` ইমপোর্ট করলে এবং তার নিজের কোনো `.claude.json` না থাকলে ccam `~/.claude.json`-এর দিকে `~/.claude/.claude.json` লিংকটি তৈরি করে। তখন সাধারণ `claude` ও `ccam main` একই অ্যাকাউন্ট দেখে। ইমপোর্ট করা অ্যাকাউন্ট আপনি `ccam share` না চালানো পর্যন্ত কিছু শেয়ার করে না, এবং কোনো শর্টকাট তৈরি হয় না।

`ccam completion bash|zsh` একটি কমপ্লিশন স্ক্রিপ্ট প্রিন্ট করে, যা অ্যাকাউন্টের নামও পূরণ করে।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- অ্যাকাউন্ট চলতে থাকলে `ccam rename` ও `ccam remove` প্রত্যাখ্যান করে (`refuseRunning`), `work is running; quit its Claude Code session first` বার্তাসহ, কারণ সেশনটি ফোল্ডারে লিখতে থাকে।
- ইমপোর্ট করা অ্যাকাউন্টে `ccam remove` শুধু নিবন্ধন বাতিল করে: `✓ main unregistered (folder ... left in place)`। আপনার `~/.claude` ট্র্যাশে সরানো হয় না।
- `profiles/`-এ যে নামের ফোল্ডার আগে থেকেই আছে, `ccam add` সেই নাম প্রত্যাখ্যান করে এবং `ccam import` দেখিয়ে দেয়।
- `ccam import` আগে থেকে নিবন্ধিত ফোল্ডার, ফোল্ডার নয় এমন পাথ, এবং নেওয়া হয়ে গেছে এমন নাম প্রত্যাখ্যান করে।
- `ccam logout` ও `ccam login` অ্যাকাউন্টটি চলছে কি না যাচাই করে না।
- রিলে Claude Code-এ পৌঁছাতে না পারলে, যেমন লগইন শেষ হয়ে যাওয়ার পরে, এটি কারণসহ `ccam: could not pass the sign-in to Claude Code:` প্রিন্ট করে। HTTP ত্রুটির মানে সাধারণত কোডের মেয়াদ শেষ; আবার লগইন শুরু করুন।

## সম্পর্কিত

- [Claude Code চালু করা](/features/launch)
- [শর্টকাট কমান্ড](/features/shortcuts)
- [MCP সার্ভার, সেটিংস ও ইতিহাস শেয়ার করুন](/docs/sharing)
- [একাধিক Claude Code অ্যাকাউন্ট](/guides/multiple-claude-code-accounts)
