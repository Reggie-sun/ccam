---
title: "ccam দিয়ে Remote Control ব্যবহার"
description: "ccam <name> rc দিয়ে যেকোনো অ্যাকাউন্টে Claude Code Remote Control চালু করুন; প্রযোজ্য অপশন, শেয়ার করা MCP সার্ভার কপির প্রস্তাব ও সমস্যার সমাধান দেখুন।"
updated: 2026-10-23
order: 4
translation: machine
sourceHash: 3e17d89f98c8
---

`ccam work rc` অ্যাকাউন্ট `work` হিসেবে Claude Code-এর Remote Control চালু করে, যাতে আপনি claude.ai/code বা Claude মোবাইল অ্যাপ থেকে এই কম্পিউটারের সেশনগুলো চালাতে পারেন। এটি `claude rc`-এর মতোই কাজ করে, শুধু একটি পার্থক্য আছে: অ্যাকাউন্টটি MCP সার্ভার শেয়ার করলে ccam একবার জিজ্ঞেস করে সেগুলো অ্যাকাউন্টের `.claude.json`-এ কপি করবে কি না, কারণ Remote Control সেশনগুলো অন্য কোথাও থেকে সার্ভার পড়ে না।

## চালু করা

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

একটি [শর্টকাট](/features/shortcuts) একইভাবে কাজ করে: `ccwork rc` মানে `ccam work rc`। যে ফোল্ডারে কাজ করতে চান সেখানে কমান্ডটি চালান, তারপর `work`-এর মতো একই Claude অ্যাকাউন্টে সাইন ইন করা অবস্থায় claude.ai/code বা Claude অ্যাপ খুলুন।

## অপশন

এখানে দুই ধরনের অপশন মেলে: ccam-এর নিজস্ব, যেগুলো অ্যাকাউন্টের নামের আগে বসে বা ccam-এর সেটিংসে থাকে, আর Remote Control-এর, যেগুলো `rc`-এর পরে বসে এবং অপরিবর্তিত অবস্থায় Claude Code-এ পৌঁছায়।

### ccam-এর অপশন ও সেটিংস

| অপশন বা সেটিং | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--dry-run` (অ্যাকাউন্টের আগে) | এনভায়রনমেন্ট ও কমান্ড প্রিন্ট করে, আর MCP বার্তা থাকলে সেটিও। জিজ্ঞেসও করে না, কপিও করে না | বন্ধ |
| `Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]` | জিজ্ঞেস করা হয় যখন অ্যাকাউন্টটি MCP সার্ভার শেয়ার করে এবং তার `.claude.json`-এ গ্রুপের বর্তমান তালিকা নেই, বা কপি করা কোনো সার্ভার সেখান থেকে হারিয়ে গেছে। `y` সেগুলো কপি করে, `N` বা Enter সেগুলো ছাড়াই চালু করে | `N` |
| `~/.config/ccam/config.toml`-এর একেবারে ওপরে `mcp_strategy = "sync"` | প্রতিবার চালু করার আগে, Remote Control হোক বা না হোক, গ্রুপের সার্ভারগুলো `.claude.json`-এ কপি করে, তাই কখনও কিছু জিজ্ঞেস করা হয় না | `"flag"` |
| `ccam group leave <group> work` | যে গ্রুপ অ্যাকাউন্টটির সঙ্গে MCP শেয়ার করে, সেখান থেকে অ্যাকাউন্টটিকে বের করে নেয়। গ্রুপের সার্ভারগুলো অ্যাকাউন্টের নিজস্ব হিসেবে তার `.claude.json`-এ কপি হয়, তাই Remote Control সেগুলো পায় এবং কিছু জিজ্ঞেস করা হয় না; `--fresh` দিলে এটি শুধু অ্যাকাউন্টের নিজস্ব সার্ভার নিয়ে চালু হয় | `default`-এ আছে |

### Remote Control-এর অপশন

`rc`-এর পরের সবকিছু Claude Code-এর। এগুলো Claude Code 2.1.304-এ `claude remote-control`-এর অপশন; `ccam work rc --help` আপনার ভার্সনের তালিকা প্রিন্ট করে।

| অপশন | কী করে |
| --- | --- |
| `--name <name>` | claude.ai/code-এ দেখানো সেশনের নাম |
| `--remote-control-session-name-prefix <prefix>` | স্বয়ংক্রিয়ভাবে তৈরি সেশন নামের প্রিফিক্স (ডিফল্ট: হোস্টের নাম) |
| `-c`, `--continue` | এই ফোল্ডারের জন্য Remote Control সর্বশেষ যে সেশন রেকর্ড করেছে তাতে আবার যুক্ত হয়, প্রায় ৪ ঘণ্টার মধ্যে |
| `--session-id <id>` | নির্দিষ্ট একটি সেশনে আবার যুক্ত হয় |
| `--permission-mode <mode>` | নতুন সেশনের অনুমতি মোড: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` বা `plan` |
| `--spawn <mode>` | `same-dir` (ডিফল্ট), `worktree` (প্রতি সেশনে একটি git worktree) বা `session` (একটি সেশন, সেটি শেষ হলে বেরিয়ে যায়) |
| `--capacity <N>` | `same-dir` ও `worktree` মোডে একসঙ্গে সর্বোচ্চ সেশন (ডিফল্ট ৩২) |
| `--[no-]create-session-in-dir` | শুরুতে বর্তমান ফোল্ডারে প্রথম সেশন তৈরি করে (ডিফল্ট চালু) |
| `--[no-]chrome` | নতুন সেশনের জন্য Claude in Chrome (ডিফল্ট বন্ধ) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | ডিবাগ ও বিস্তারিত আউটপুট |

ccam Claude Code-এর কমান্ড সেভাবেই পড়ে যেভাবে Claude Code পড়ে, অপশনগুলোকে তাদের মানসহ বাদ দিয়ে, তাই `ccam work rc --name laptop` ও `ccam work --model opus rc` দুটোই Remote Control চালু করে। ccam চেনে না এমন অপশন মান ছাড়া অপশন হিসেবে গণ্য হয়; নিরাপদ থাকতে `rc` অ্যাকাউন্টের নামের ঠিক পরেই বসান।

## ccam কী করে

চালু করা স্বাভাবিক পথেই চলে ([Claude Code চালু করা](/features/launch)), শুধু MCP ধাপে একটি পরিবর্তন থাকে:

1. ccam অ্যাকাউন্টটি খুঁজে বের করে, তার শেয়ার লিংক যাচাই করে এবং `claude` নির্ধারণ করে, যেকোনো চালু করার মতোই।
2. এটি দেখে যে কমান্ডটি `rc` বা `remote-control`, তাই `--mcp-config` বাদ দেয়।
3. অ্যাকাউন্টটি MCP সার্ভার শেয়ার করলে এটি গ্রুপের `mcp.json`-কে `~/.local/state/ccam/mcp/work/`-এর `sync-stamp`-এর সঙ্গে তুলনা করে, এবং যাচাই করে যে প্রতিটি গ্রুপ সার্ভার এখনও `.claude.json`-এ আছে। `.claude.json`-এ আগে থেকেই বর্তমান সার্ভারগুলো থাকলে, বা গ্রুপে কোনো সার্ভার না থাকলে, এটি সরাসরি ৬ নম্বর ধাপে চলে যায়।
4. না হলে এটি সেই সার্ভারগুলোর নাম জানায় যেগুলো সেশনগুলো পাবে না:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

5. `y` দিলে এটি `.claude.json`-এর একটি ব্যাকআপ রাখে, Claude Code-এর নিজস্ব লক নেয়, গ্রুপের সার্ভারগুলো `mcpServers`-এ যোগ করে এবং অ্যাকাউন্টের অন্য সার্ভারগুলো রেখে দেয়; কোনো গ্রুপ সার্ভারের নামে অ্যাকাউন্টের নিজস্ব সার্ভার থাকলে সেটি অ্যাকাউন্টেরই থাকে, এবং ccam তা জানায়। `N` দিলে, টার্মিনাল না থাকলে, বা অ্যাকাউন্টটি চলমান থাকলে এটি কিছুই বদলায় না এবং তা জানিয়ে দেয়:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. এটি সেশনটি রেকর্ড করে, `CLAUDE_CONFIG_DIR` অ্যাকাউন্টের ফোল্ডারে নির্ধারণ করে এবং `exec`-এর মাধ্যমে নিজেকে `claude rc` দিয়ে প্রতিস্থাপন করে।

এরপর Remote Control একটি সার্ভার হিসেবে চলতে থাকে। আপনার জন্য এটি যে সেশনই চালু করে, সেটি ওই প্রসেসের চাইল্ড, তাই সেটি `work` হিসেবে চলে: `work`-এর লগইন, সেটিংস ও কথোপকথনের ইতিহাস নিয়ে, এবং `work`-এর ব্যবহারসীমা থেকে গোনা হয়। ccam-এর চালু করা সার্ভার চলার সময় ccam অ্যাকাউন্টটিকে চলমান ধরে, আর সাধারণ `claude remote-control` দিয়ে চালু করা সার্ভারের ক্ষেত্রেও তাই, যতক্ষণ তার সেশন থাকে; এমন সার্ভারের এখনও কোনো সেশন না থাকলে সেটি দেখা যায় না।

ড্রাই রানে পার্থক্যটি দেখা যায়। স্বাভাবিক চালু করায় গ্রুপের সার্ভারগুলো একটি আর্গুমেন্ট হিসেবে যায়, Remote Control-এ যায় না:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## কেন এভাবে কাজ করে

**Remote Control কোনো `--mcp-config` নেয় না।** স্বাভাবিক চালু করায় ccam গ্রুপের সার্ভারগুলো `--mcp-config=<path>` হিসেবে Claude Code-কে দেয় এবং কখনও Claude Code-এর ফাইলে লেখে না। `claude remote-control` তার কমান্ড নামের আগে এই অপশন প্রত্যাখ্যান করে এবং পরে কোনোটিই নেয় না, তাই ccam 0.6.0-এর আগে MCP সার্ভার শেয়ার করা প্রতিটি অ্যাকাউন্টে `ccam work rc` ব্যর্থ হতো। এখন ccam অপশনটি বাদ দেয়।

**এর সেশনগুলো সার্ভার শুধু `.claude.json` থেকে পড়ে।** Remote Control যে সেশনগুলো চালু করে সেগুলো কমান্ড লাইন থেকে অপশন উত্তরাধিকারসূত্রে পায় না। ব্যবহারকারীর MCP সার্ভার তারা কেবল একটি জায়গায় পায়: অ্যাকাউন্টের `.claude.json`। এই ফাইলে অ্যাকাউন্টের লগইন তথ্য আর ফোল্ডারভিত্তিক অবস্থাও থাকে, তাই ccam এটিকে `settings.json`-এর মতো শেয়ার করা কপির সঙ্গে লিংক করতে পারে না। সেশনগুলোকে শেয়ার করা তালিকা দেওয়ার একমাত্র উপায় হলো সার্ভারগুলো এতে কপি করা।

**লেখার আগে ccam জিজ্ঞেস করে।** ডিফল্ট ফ্ল্যাগ কৌশলে চালু করার সময় ccam `.claude.json`-এ হাত দেয় না। Remote Control-ই একমাত্র চালু করা যেখানে তাকে সেখানে লিখতে হয়, তাই এটি আগে জিজ্ঞেস করে, ব্যাকআপ রাখে, এবং `mcpServers` ছাড়া আর কিছুতে হাত দেয় না। (চালু করার বাইরে, MCP শেয়ারিং ছাড়া বা আবার যোগ দেওয়া এবং `ccam mcp import`-ও সেখানে লেখে, তাদের পরিকল্পনায় যেমন বলা থাকে।)

**অ্যাকাউন্ট চলার সময় এটি লেখে না।** সেশন চলাকালে Claude Code `.claude.json` নতুন করে লেখে, আর বাইরে থেকে করা পরিবর্তন হারিয়ে যাবে বা তার ওপর লেখা হয়ে যাবে। এ কারণেই অ্যাকাউন্টের কোনো সেশন চললে, আরেকটি Remote Control সার্ভার এবং ccam ছাড়া চালু হওয়া Claude Code-সহ, ccam কপি করা বাদ দেয়।

**প্রতিটি পরিবর্তনে এটি একবার জিজ্ঞেস করে।** কপির পরে `sync-stamp` গ্রুপের `mcp.json`-এর পরিবর্তনের সময় ধরে রাখে। ccam আবার জিজ্ঞেস করে শুধু তখনই, যখন ওই ফাইল বদলায়, যেমন `ccam mcp add`-এর পরে, অথবা যখন কপি করা কোনো সার্ভার `.claude.json` থেকে হারিয়ে যায়। আপনি `N` উত্তর দিলে এখনও কোনো কপি নেই, তাই প্রতিবার Remote Control চালু করার সময় এটি জিজ্ঞেস করে।

**`ccam mcp` কপিগুলো হালনাগাদ রাখে।** নামের সংঘাত হলে `.claude.json`-এর এন্ট্রি `--mcp-config`-এর ওপর জেতে, তাই একটি পুরোনো কপি স্বাভাবিক চালু করাতেও নতুন সংজ্ঞাকে আড়াল করবে। সেজন্য `ccam mcp add`, `add-json`, `remove` ও `import`, এবং `ccam migrate`, চলমান নয় এমন প্রতিটি অ্যাকাউন্টের কপি আপডেট করে, এবং কোনো কপি বদলালে তা জানায়:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**টার্মিনাল না থাকলে উত্তর "না"।** Remote Control প্রায়ই কোনো স্ক্রিপ্ট, সার্ভিস বা `nohup` থেকে চলে, যেখানে উত্তর দেওয়ার কেউ থাকে না। তখন ccam অপেক্ষা না করে সার্ভার ছাড়াই চালু করে।

## সমস্যা ও সমাধান

### Remote Control সেশনে শেয়ার করা MCP সার্ভারগুলো নেই

প্রশ্নের উত্তর `N` দেওয়া হয়েছিল, টার্মিনাল ছিল না, বা অ্যাকাউন্টটি চলমান ছিল। Remote Control থামান, অ্যাকাউন্টের অন্য সেশনগুলো বন্ধ করুন, এবং টার্মিনালে আবার `ccam work rc` চালান; উত্তর দিন `y`। কখনও জিজ্ঞেস না করাতে চাইলে `~/.config/ccam/config.toml`-এ `mcp_strategy = "sync"` নির্ধারণ করুন। কোনো সেশনে আগে থেকে থাকা সার্ভার সেটি চলাকালে বদলায় না, তাই কপির পরে একটি নতুন সেশন শুরু করুন।

### "Account work is running, so its .claude.json is left alone"

`work`-এর একটি সেশন এখনও খোলা: একটি সাধারণ `ccam work`, একটি শর্টকাট, বা আরেকটি Remote Control সার্ভার। সেটি থেকে বেরিয়ে আবার `ccam work rc` চালান। অন্য অ্যাকাউন্টের চালানো সেশন গোনা হয় না।

### ccam প্রতিবার জিজ্ঞেস করে

আপনি এখনও সার্ভারগুলো কপি করেননি। একবার `y` উত্তর দিন, বা `mcp_strategy = "sync"` নির্ধারণ করুন। Remote Control-এ কখনও সার্ভার না চাইলে সেগুলো না রেখেই অ্যাকাউন্টের সঙ্গে MCP শেয়ার করা বন্ধ করুন: `ccam group leave default work --fresh`, অথবা MCP-কে এমন একটি আলাদা গ্রুপে রাখুন যাতে অ্যাকাউন্টটি নেই। সাধারণ `ccam group leave` গ্রুপের সার্ভারগুলো অ্যাকাউন্টের `.claude.json`-এ কপি করে, তাই Remote Control সেগুলো পায়।

### `ccam work rc` সঙ্গে সঙ্গে `--mcp-config` সংক্রান্ত ত্রুটি দিয়ে ব্যর্থ হয়

আপনার ccam 0.6.0-এর চেয়ে পুরোনো এবং Remote Control-কে `--mcp-config` দেয়। `ccam update` চালান। বর্তমান ccam-এ একই ত্রুটি মানে ccam চেনে না এমন কোনো মান-সহ অপশন `rc`-এর আগে এসেছে, তাই ccam মানটিকে কমান্ড ধরে নিয়েছে; `rc` অ্যাকাউন্টের নামের ঠিক পরে বসান।

### `ccam mcp` দিয়ে করা পরিবর্তন Remote Control-এ পৌঁছায় না

`ccam mcp add` ও `remove` শুধু চলমান নয় এমন অ্যাকাউন্টের কপি আপডেট করে, এবং তা জানিয়ে দেয়: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Remote Control থামিয়ে আবার চালু করুন; ccam বদলে যাওয়া `mcp.json` দেখে আবার জিজ্ঞেস করবে।

### গ্রুপেরটির বদলে আমার নিজের একটি সার্ভার ব্যবহার হয়

অ্যাকাউন্টে গ্রুপের কোনো সার্ভারের নামে নিজস্ব সার্ভার থাকলে কপি সেটিকে জায়গামতো রেখে দেয় এবং তা জানায়: `account work has its own MCP servers fs, so group default's of those names are not copied there`। গ্রুপেরটি ব্যবহার করতে অ্যাকাউন্টেরটি সরান (`claude mcp remove -s user <name>`); আপনার ভার্সনটি প্রতিটি অ্যাকাউন্টকে দিতে চাইলে `ccam mcp import work <server> --force` দিয়ে সেটি গ্রুপে সরিয়ে নিন। `.claude.json`-এর ব্যাকআপ আছে `~/.local/state/ccam/mcp/work/backups/`-এ, যা সর্বশেষ পাঁচটি রাখে।

### আমি কপি করা সার্ভারগুলো সরাতে চাই

`ccam mcp remove <server>` গ্রুপ থেকে এবং প্রতিটি কপি থেকে একটি সার্ভার সরিয়ে দেয়। কপির আগের অবস্থায় ফিরতে, অ্যাকাউন্টের কোনো সেশন না চলা অবস্থায় `~/.local/state/ccam/mcp/work/backups/` থেকে সর্বশেষ ব্যাকআপ পুনরুদ্ধার করুন।

### Remote Control ফোল্ডারটিকে বিশ্বাস করবেন কি না জিজ্ঞেস করে, বা টার্মিনাল ছাড়া বেরিয়ে যায়

ফোল্ডারের প্রতি বিশ্বাস প্রতিটি অ্যাকাউন্টের জন্য আলাদাভাবে `.claude.json`-এ রাখা হয়, তাই `personal` হিসেবে যে ফোল্ডার আপনি বিশ্বাস করেছেন সেটি `work`-এর কাছে নতুন। ওই ফোল্ডারে একবার `ccam work` চালিয়ে সম্মতি দিন, তারপর সেখানে Remote Control চালু করুন।

### "You must be logged in with a Claude account that has a subscription"

Remote Control-এর জন্য Claude সাবস্ক্রিপশনসহ লগইন দরকার। `ccam list`-এ অ্যাকাউন্টের ইমেইল দেখুন, এবং প্রয়োজনে `ccam login work` দিয়ে আবার সাইন ইন করুন।

### ফোনে সেশনটি দেখা যায় না

Claude অ্যাপ বা claude.ai/code অন্য একটি Claude অ্যাকাউন্টে সাইন ইন করা। `ccam list` `work`-এর জন্য যে ইমেইল দেখায় তার সঙ্গে মিলিয়ে দেখুন।

### ccam সতর্ক করে যে এই ফোল্ডারে আরেকটি Claude Code কথোপকথনের ইতিহাস শেয়ার করছে

একই ইতিহাস গ্রুপের আরেকটি অ্যাকাউন্ট একই ফোল্ডারে চলছে। দুই অ্যাকাউন্ট থেকে একই কথোপকথন চালিয়ে যাওয়া সেশনগুলো ইতিহাস মিশিয়ে ফেলতে পারে। শুধু তখনই `y` উত্তর দিন যখন সেগুলো আলাদা আলাদা কথোপকথনে কাজ করবে, নয়তো অন্য ফোল্ডারে Remote Control চালু করুন।

## সম্পর্কিত

- [Claude Code চালু করা](/features/launch)
- [শেয়ার করা MCP সার্ভার](/features/mcp)
- [অ্যাকাউন্টগুলোর মধ্যে শেয়ার করা](/features/share)
- [অ্যাকাউন্ট পরিচালনা](/features/accounts)
