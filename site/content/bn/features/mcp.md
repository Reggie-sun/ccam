---
title: "শেয়ার করা MCP সার্ভার"
description: "ccam mcp add, add-json, list, remove ও import দিয়ে প্রতি গ্রুপে একটি MCP সার্ভার তালিকা রাখুন, যা --mcp-config দিয়ে দেওয়া বা প্রতিটি অ্যাকাউন্টে সিঙ্ক হয়।"
updated: 2026-10-20
order: 5
translation: machine
sourceHash: 6d3b3c38c6cb
---

Claude Code ইউজার MCP সার্ভারগুলো `.claude.json`-এ রাখে, যে ফাইলে অ্যাকাউন্টের লগইনও থাকে, তাই ফাইলটি অ্যাকাউন্টগুলোর মধ্যে শেয়ার করা যায় না। ccam শেয়ার করা সার্ভারগুলো একটি গ্রুপ ফাইল `mcp.json`-এ রাখে এবং চালু করার সময় গ্রুপের প্রতিটি অ্যাকাউন্টকে সেটি দেয়। আপনার সব অ্যাকাউন্টের জন্য একবারে একটি সার্ভার যোগ করতে `ccam mcp` ব্যবহার করুন, আর একটি সেশনের ভেতরে যোগ করা সার্ভার অন্যগুলোতে সরাতে `ccam mcp import`।

## ব্যবহার

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

পরিবর্তন প্রতিটি অ্যাকাউন্টের পরের বার চালু হওয়া থেকে প্রযোজ্য হয়; চলমান সেশনগুলো যে সার্ভার নিয়ে শুরু হয়েছিল সেগুলোই রাখে:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

গ্রুপে আগে থেকে থাকা কোনো নাম যোগ করলে তার সংজ্ঞা প্রতিস্থাপিত হয় এবং `added`-এর বদলে `updated` প্রিন্ট হয়। সার্ভার কোনো অ্যাকাউন্টে পৌঁছায় শুধু যদি অ্যাকাউন্টটি ওই গ্রুপে থাকে এবং গ্রুপটি MCP শেয়ার করে (`ccam group set <group> --add mcp`); গ্রুপটি তা না করলে `add` সেটি জানায়।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json` ও `remove`-এর লক্ষ্য গ্রুপ | যে একমাত্র গ্রুপের মাধ্যমে কোনো অ্যাকাউন্ট MCP শেয়ার করে, অন্যথায় `default` |
| `-- <command> [args…]` (`add`) | এই কমান্ডটি stdio সার্ভার হিসেবে চালায় | |
| `--url <url>` (`add`) | রিমোট সার্ভারের URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` সার্ভারের ট্রান্সপোর্ট | `http` |
| `--header <'Name: value'>` (`add`) | HTTP হেডার, একাধিকবার দেওয়া যায় | নেই |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio সার্ভারের জন্য এনভায়রনমেন্ট ভেরিয়েবল, একাধিকবার দেওয়া যায় | নেই |
| `--force` (`import`) | একই নামের আলাদা গ্রুপ সার্ভার প্রতিস্থাপন করে; পুরোনো `mcp.json`-এর ব্যাকআপ রাখা হয় | বন্ধ |
| `--keep-local` (`import`) | অ্যাকাউন্টের `.claude.json`-এ এন্ট্রিগুলো রেখে দেয় | বন্ধ |
| `-h`, `--help` | প্রতিটি কমান্ডের সাহায্য | |
| `--dry-run` (গ্লোবাল) | শুধু `mcp list` এটি মানে; `add`, `add-json`, `remove` ও `import` এটি প্রত্যাখ্যান করে এবং কিছু বদলায় না | বন্ধ |

`add`-এর জন্য `--`-এর পরে একটি কমান্ড অথবা `--url` লাগে; না হলে এটি `a command to run (after --) or --url is required` দিয়ে ব্যর্থ হয়। `--header` `--url` সার্ভারে প্রযোজ্য, আর `--env` কমান্ড সার্ভারে। `add-json` `claude mcp add-json`-এর মতো একই JSON নেয়, এবং সংজ্ঞাটি অবশ্যই `command` বা `url`-সহ একটি অবজেক্ট হতে হবে। `import` সবসময় সেই গ্রুপ ব্যবহার করে যার মাধ্যমে অ্যাকাউন্টটি MCP শেয়ার করে, আর `list` প্রতিটি গ্রুপ দেখায়, তাই `-g` এগুলো বদলায় না।

## এটি কীভাবে কাজ করে

সংশ্লিষ্ট ফাইলগুলো:

- `~/.local/share/ccam/groups/<group>/mcp.json`: শেয়ার করা সার্ভার, `{"mcpServers": {…}}`, যে ফরম্যাট `--mcp-config` পড়ে। এটি `fsutil.AtomicWrite` ও মোড `0600` দিয়ে লেখা হয়, কারণ সার্ভার এন্ট্রিতে প্রায়ই API কী থাকে।
- `~/.local/share/ccam/profiles/<name>/.claude.json`: অ্যাকাউন্টের সার্ভারসহ Claude Code-এর নিজস্ব ফাইল। ccam `profile.TopLevelValue` স্ক্যানার দিয়ে শুধু এর `mcpServers` মান পড়ে এবং চালু করার পথে কখনও পুরো ফাইল পার্স করে না।
- `~/.local/state/ccam/mcp/<name>/`: `seen` ফাইল, `sync-stamp` ও `.claude.json`-এর ব্যাকআপ।

দুটি কৌশল আছে, যা `~/.config/ccam/config.toml`-এ `mcp_strategy` দিয়ে নির্ধারিত হয়। ডিফল্ট হলো ফ্ল্যাগ কৌশল:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` অ্যাকাউন্টের MCP গ্রুপ খুঁজে বের করে। `mcp.FlagArgs` পার্স না করে একটি `stat` দিয়ে যাচাই করে যে গ্রুপের `mcp.json` আছে।
2. `Plan.InsertArgs` `--mcp-config=<path>` `claude`-এর ঠিক পরে, আপনার নিজের আর্গুমেন্টের আগে বসায়, তাই আপনার নিজের `--mcp-config` বা `--strict-mcp-config` তখনও এর ওপরে প্রযোজ্য থাকে।
3. `mcp.NewLocal` যাচাই করে শেষবার চালু করার পর থেকে অ্যাকাউন্টের `.claude.json` বদলেছে কি না, এবং শুধু তখনই এর সার্ভারগুলো পার্স করে।
4. `syscall.Exec` ccam-কে Claude Code দিয়ে প্রতিস্থাপন করে।

পাথটি দুটি নয়, একটি আর্গুমেন্ট হিসেবে দেওয়া হয়, `--mcp-config=<path>`। Claude Code-এর `--mcp-config` একাধিক মান নেয়, তাই `claude --mcp-config a.json "prompt"` প্রম্পটটিকে দ্বিতীয় কনফিগ ফাইল হিসেবে পড়ত এবং `MCP config file not found` দিয়ে ব্যর্থ হতো। এক-আর্গুমেন্টের রূপ সেখানেই থামে, এবং পরে আসা প্রম্পট বা `mcp list`-এর মতো সাবকমান্ড কাজ করে। Claude Code অ্যাকাউন্টের নিজস্ব সার্ভারের পাশাপাশি `--mcp-config` সার্ভারগুলোও লোড করে; নামের সংঘাত হলে `.claude.json`-এর এন্ট্রি জেতে। ফ্ল্যাগ কৌশল শুধু Remote Control-এর জন্য `.claude.json` লেখে।

`claude remote-control` (`rc`) কোনো `--mcp-config` নেয় না, এবং এর সেশনগুলো সার্ভার শুধু `.claude.json` থেকে পড়ে। `ccam work rc`-এর জন্য ccam ফ্ল্যাগটি বাদ দেয় এবং টার্মিনালে জিজ্ঞেস করে, গ্রুপের সার্ভারগুলো ব্যাকআপসহ `.claude.json`-এ কপি করবে কি না, যেমন `sync` কৌশল করে। কোনো কপি তৈরি না হওয়া পর্যন্ত এটি জিজ্ঞেস করে, তারপর শুধু গ্রুপের `mcp.json` বদলালে আবার ([Remote Control](/features/launch#remote-control))। নামের সংঘাত হলে এমন কপি `--mcp-config`-এর ওপর জেতে বলে `ccam mcp add`, `add-json`, `remove` ও `import` গ্রুপের অন্য অ্যাকাউন্টগুলোর কাছে থাকা কপিগুলো আপডেট করে:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

সেশনের ভেতরে `claude mcp add` দিয়ে যোগ করা সার্ভার শুধু ওই অ্যাকাউন্টের `.claude.json`-এ যায়। ccam সেই সময় এটি দেখতে পায় না, কারণ Claude Code তাকে প্রতিস্থাপন করেছে, তাই এটি পরের বার চালু করার সময় যাচাই করে। `seen` ফাইলে থাকে `.claude.json`-এর পরিবর্তনের সময় এবং আগে জানানো সার্ভারের নাম। সময় আলাদা হলে ccam নামগুলো গ্রুপের সঙ্গে তুলনা করে এবং প্রতিটি নতুন সার্ভারের জন্য একবার একটি ইঙ্গিত প্রিন্ট করে, ইনপুটের জন্য অপেক্ষা না করে:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` থাকলে সার্ভারগুলো `claude mcp list` ও `/mcp`-এ সাধারণ ইউজার সার্ভার হিসেবে দেখা যায়, এবং ccam চালু করার আগে সেগুলো প্রতিটি অ্যাকাউন্টের `.claude.json`-এ লেখে:

<!-- motion:mcp-sync -->

1. `mcp.Sync` গ্রুপের `mcp.json`-এর পরিবর্তনের সময় `sync-stamp`-এর সঙ্গে তুলনা করে। না বদলালে অ্যাকাউন্টের `.claude.json` খোলা হয় না।
2. এটি Claude Code-এর নিজস্ব লক নেয়, ফাইলের পাশে `mkdir` দিয়ে তৈরি ফোল্ডার `.claude.json.lock`। এটি সর্বোচ্চ ২ সেকেন্ড অপেক্ষা করে এবং ১০ সেকেন্ডের পুরোনো লককে বাসি বলে গণ্য করে, ঠিক যেমন Claude Code-এর লক লাইব্রেরি করে।
3. এটি বর্তমান `.claude.json` স্টেট ফোল্ডারের `backups/`-এ সংরক্ষণ করে এবং সবচেয়ে নতুন ৫টি রাখে।
4. এটি `mcpServers`-এ প্রতিটি গ্রুপ সার্ভার বসায়, আগে নিজের লেখা যে সার্ভারগুলো গ্রুপে আর নেই সেগুলো সরায়, অ্যাকাউন্টের নিজস্ব সার্ভারগুলো রাখে, এবং ফাইলটি তার মূল মোডসহ অ্যাটমিকভাবে লেখে।
5. এটি নতুন স্ট্যাম্প ও তার পরিচালিত নামগুলো নথিভুক্ত করে, এবং Claude Code `--mcp-config` ছাড়াই চালু হয়।

অ্যাকাউন্টের কোনো সেশন চলতে থাকলে সিঙ্ক বাদ দেওয়া হয়, `ccam: account work is already running, skipping MCP sync` বার্তাসহ, কারণ Claude Code চলার সময় `.claude.json` আবার লেখে। ব্যর্থ সিঙ্ক `ccam: MCP sync failed (launching anyway): …` প্রিন্ট করে এবং তবুও চালু করে। সেশনের ভেতরে যোগ করা সার্ভারের ইঙ্গিত শুধু ফ্ল্যাগ কৌশলের অংশ।

`ccam mcp import` একটি অ্যাকাউন্টের নিজস্ব সার্ভারগুলো তার গ্রুপে সরায়:

<!-- motion:mcp-import -->

ccam `mcp.LocalServers` দিয়ে অ্যাকাউন্টের সার্ভারগুলো পড়ে এবং প্রতিটিকে গ্রুপের একই নামের সার্ভারের সঙ্গে `mcp.SameDef` দিয়ে তুলনা করে, যা কী-এর ক্রম ও স্পেস উপেক্ষা করে। যে সার্ভার গ্রুপে নেই বা একই সংজ্ঞায় আছে, সেটি কপি করা হয়। আলাদা হলে বাদ দেওয়া হয়, কারণ গ্রুপের সংজ্ঞাই বাকি সব অ্যাকাউন্ট ব্যবহার করে:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` দিলে গ্রুপের সার্ভার প্রতিস্থাপিত হয়, এবং আগের `mcp.json` প্রথমে গ্রুপ ফোল্ডারের ভেতরে `.ccam-backup/<time>/mcp.json`-এ সংরক্ষিত হয়। তারপর `mcp.RemoveLocal` একই লকের অধীনে এবং স্টেট ফোল্ডারে একটি ব্যাকআপের পরে অ্যাকাউন্টের `.claude.json` থেকে কপি করা এন্ট্রিগুলো সরায়। এটি দরকার, কারণ নামের সংঘাত হলে অ্যাকাউন্টের এন্ট্রি গ্রুপেরটির ওপর জিতত। `--keep-local` সেগুলো রেখে দেয়।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- **সাইন-ইন।** যে সার্ভারগুলো OAuth দিয়ে সাইন ইন করে, সেগুলো তাদের টোকেন সার্ভার তালিকায় নয়, প্রতিটি অ্যাকাউন্টের `.credentials.json`-এ (`mcpOAuth`) রাখে, তাই সার্ভার শেয়ার করলে তার সাইন-ইন শেয়ার হয় না। প্রতিটি অ্যাকাউন্টে একবার `/mcp` দিয়ে সাইন ইন করুন, অথবা [`ccam migrate`](/features/migrate) দিয়ে `~/.claude` থেকে সাইন-ইনগুলো কপি করুন। গ্রুপ কখনও কোনো লগইন ফাইল পায় না।
- **চলমান অ্যাকাউন্ট।** `import` গ্রুপে কপি করে, কিন্তু চলমান অ্যাকাউন্টের নিজস্ব এন্ট্রিগুলো রেখে দেয় এবং তা জানায়; সেশন শেষ হলে আবার চালান। সিঙ্ক চলমান অ্যাকাউন্টগুলো বাদ দেয়।
- **Remote Control-এর জন্য কপি।** চলমান অ্যাকাউন্টের কপিতে হাত দেওয়া হয় না: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` থাকলে কিছুই আপডেট করতে হয় না, কারণ প্রতিবার চালু করার সময় সিঙ্ক হয়।
- **ভাঙা ফাইল।** `.claude.json` বৈধ JSON না হলে সিঙ্ক সেটি অস্পর্শিত রাখে এবং ত্রুটি জানায়। ২ সেকেন্ডে লক নেওয়া না গেলে ccam জানায় `could not lock … (Claude Code is writing the file)`।
- **যে অ্যাকাউন্টগুলো MCP শেয়ার করে না।** `import` প্রত্যাখ্যান করে এবং আপনাকে আগে অ্যাকাউন্টটিকে এমন একটি গ্রুপে রাখতে বলে যা MCP শেয়ার করে (`ccam group join <group> side`)।
- **গ্রুপের নাম।** `-g` শুধু একটি গ্রুপের নাম নেয়; এটি `groups/`-এর নিচে একটি ফোল্ডার হয়, তাই `../x`-এর মতো নাম প্রত্যাখ্যান করা হয়।

## সম্পর্কিত

- [অ্যাকাউন্টগুলোর মধ্যে শেয়ার করা](/features/share)
- [ccam দিয়ে Remote Control ব্যবহার](/guides/remote-control)
- [MCP সার্ভার, সেটিংস ও ইতিহাস শেয়ার করুন](/docs/sharing)
- [আগের সেটআপ কপি করুন](/features/migrate)
- [সব কমান্ড](/docs/commands)
