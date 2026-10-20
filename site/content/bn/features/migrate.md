---
title: "আগের সেটআপ কপি করা"
description: "ccam migrate দিয়ে বিদ্যমান ~/.claude-এর সেটিংস, CLAUDE.md, এক্সটেনশন, MCP সার্ভার ও MCP সাইন-ইন একটি ccam গ্রুপে কপি করুন।"
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 7cee52e75749
---

ccam-এর আগে Claude Code ব্যবহার করে থাকলে আপনার সেটআপ আছে `~/.claude` ও `~/.claude.json`-এ। `ccam migrate` সেটি একটি শেয়ার গ্রুপে কপি করে, যাতে গ্রুপের মাধ্যমে ওই আইটেমগুলো শেয়ার করা প্রতিটি অ্যাকাউন্ট আপনার আগে থেকে থাকা সেটআপ নিয়েই শুরু হয়। উৎস ফোল্ডারটি শুধু পড়া হয়, এবং সাধারণ `claude` আগের মতোই কাজ করে। আপনার প্রথম `ccam add`-এর পরে এটি একবার চালান, এবং যখনই কোনো অ্যাকাউন্ট MCP সার্ভার শেয়ার করা শুরু করে তখন আবার।

## ব্যবহার

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

টার্মিনালে, আইটেম ফ্ল্যাগ ছাড়া এবং `--yes` ছাড়া, ccam প্রথমে একটি চেকলিস্ট দেখায়:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (বা `x`) কার্সরের নিচের আইটেম টগল করে, অ্যারো কী বা `j`/`k` সরায়, Enter পরিকল্পনা দেখায়, আর Esc, `q` বা Ctrl+C প্রিন্ট করে `Cancelled.`। যে আইটেমের জন্য উৎসে কিছু নেই, সেগুলো `not found` হিসেবে চিহ্নিত থাকে। সবকিছু আনচেক করলে ccam `Nothing selected.` প্রিন্ট করে থেমে যায়।

তারপর ccam পরিকল্পনাটি প্রিন্ট করে এবং কপি করার আগে জিজ্ঞেস করে:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

একটি গ্রুপ শুধু নিজের অ্যাকাউন্টগুলোতে পৌঁছায়, এবং শুধু যে আইটেমগুলো সে শেয়ার করে সেগুলোর জন্য, তাই সারাংশটি গ্রুপের অ্যাকাউন্টগুলোর তালিকা দেয়। কপি করা কোনো আইটেম গ্রুপটি শেয়ার না করলে, বা কিছু অ্যাকাউন্ট তাতে না থাকলে, ccam সেটি বদলানোর `ccam group set <group> --add <items>` ও `ccam group join <group> <accounts>` কমান্ডগুলো প্রিন্ট করে।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `[folder]` | যে Claude Code ফোল্ডার থেকে কপি করা হবে | `~/.claude` |
| `--settings` | `settings.json` কপি করে | বন্ধ |
| `--memory` | গ্লোবাল `CLAUDE.md` কপি করে | বন্ধ |
| `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` কপি করে | বন্ধ |
| `--mcp` | ইউজার MCP সার্ভার কপি করে (`.claude.json`-এর `mcpServers`) | বন্ধ |
| `--history` | কথোপকথনের ইতিহাস কপি করে (`projects/`, `history.jsonl`) | বন্ধ |
| `-g`, `--group string` | যে শেয়ার গ্রুপে কপি করা হবে | `default` |
| `-y`, `--yes` | জিজ্ঞেস না করেই কপি করে; দ্বন্দ্ব হলে গ্রুপের কপি রাখা হয় | বন্ধ |
| `--dry-run` (গ্লোবাল) | পরিকল্পনা প্রিন্ট করে এবং কোনো প্রশ্ন বা কপির আগে থেমে যায় | বন্ধ |

কোনো আইটেম ফ্ল্যাগ না থাকলে টার্মিনালে চেকলিস্ট ঠিক করে। টার্মিনালের বাইরে, বা `--yes` দিলে, ccam কথোপকথনের ইতিহাস বাদে সবকিছু কপি করে। ইতিহাস ডিফল্টভাবে বাদ থাকে, কারণ `projects/` বড় হতে পারে।

`--group on` মানে `default`। `--group off` ও `--group local` প্রত্যাখ্যান করা হয় `--group takes a group name, not "off"` দিয়ে। `ccam migrate` `--dry-run` মানে, এবং ড্রাই রানেও চেকলিস্ট দেখা যায়, যাতে আপনি একটি নির্বাচনের পূর্বরূপ দেখতে পারেন।

## এটি কীভাবে কাজ করে

`internal/cli/migrate.go` আইটেমগুলো বেছে নেয়, তারপর `migrate.Build` উৎস ও গ্রুপ পড়ে একটি পরিকল্পনা ফেরত দেয়। আপনি নিশ্চিত করার পরে `Plan.Apply` পর্যন্ত কিছুই লেখা হয় না।

1. **উৎস খুঁজে বের করা।** যে পাথ নেই, ফোল্ডার নয়, বা যাতে Claude Code-এর কোনো পরিচিত এন্ট্রি নেই, `checkSource` সেটি প্রত্যাখ্যান করে। `migrate.ClaudeJSON` ফোল্ডারটির `.claude.json` খুঁজে বের করে: এর ভেতরে, অথবা ভেতরে না থাকলে `~/.claude`-এর জন্য তার পাশের `~/.claude.json`। `CLAUDE_CONFIG_DIR` নির্ধারিত না থাকলে Claude Code এটি সেখানেই রাখে।
2. **আইটেমকে পাথের সঙ্গে মেলানো।** প্রতিটি পাথ আসে `share.Layout` থেকে, যা গ্রুপগুলোর ব্যবহৃত একই অনুমোদিত তালিকা। MCP সার্ভারগুলো `.claude.json` থেকে `mcp.LocalServers` দিয়ে পড়া হয়, যা পুরো ফাইল পার্স না করে `profile.TopLevelValue` স্ক্যানার ব্যবহার করে।
3. **কপি, কখনও সরানো নয়।** ফাইলগুলো `fsutil.AtomicWrite` বা `fsutil.CopyFile` দিয়ে `~/.local/share/ccam/groups/<group>/`-এ যায়; MCP সার্ভারগুলো যায় গ্রুপের `mcp.json`-এ।

প্রতিটি আইটেম একটি উৎস পাথকে একটি গ্রুপ পাথের সঙ্গে মেলায়:

<!-- motion:migrate-items -->

ইতিহাসের দুটি সারি ড্যাশ দেওয়া, কারণ সেগুলো শুধু `--history` দিলে কপি হয়। `.credentials.json` ও `.claude.json` `share.Never`-এ তালিকাভুক্ত এবং কখনও গ্রুপে পৌঁছায় না: কোনোটি কখনও কাঠামোতে দেখা দিলে `Build` একটি অভ্যন্তরীণ ত্রুটি ফেরত দেয়, এবং ফোল্ডার ঘোরার সময় নাম দেখে সেগুলো বাদ দেওয়া হয়।

`buildEntry` প্রতিটি পাথকে একটি অবস্থা দেয়, যা প্রতিটি পরিকল্পনা লাইন বর্ণনা করে:

<!-- motion:migrate-states -->

- **Missing**: উৎসে নেই, অথবা যেখানে ফোল্ডার থাকার কথা সেখানে একটি ফাইল। প্রতি আইটেমে একবার `not found, skipped` হিসেবে তালিকাভুক্ত হয়।
- **Linked**: উৎস পাথ আগে থেকেই গ্রুপের কপিতে নির্দেশ করে, যেমন উৎসটি যখন একটি ccam প্রোফাইল ফোল্ডার।
- **Same**: গ্রুপে আগে থেকেই হুবহু একই বিষয়বস্তু আছে।
- **New**: গ্রুপে এটি নেই, তাই পুরোটা কপি হয়।
- **Merge**: দুই দিকেই একটি ফোল্ডার। `copyMissing` শুধু সেই ফাইলগুলো যোগ করে যা গ্রুপে নেই; আলাদা ফাইলগুলো গ্রুপে যেমন আছে তেমনই থাকে এবং পরিকল্পনায় লেখা থাকে `group copies kept`। `plugins/`-এর জন্য রেজিস্ট্রি `installed_plugins.json` ও `known_marketplaces.json` মার্জ করা হয়, এবং প্লাগইনের পূর্ণ পাথগুলো গ্রুপের ফোল্ডার নির্দেশ করার জন্য আবার লেখা হয়।
- **Conflict**: দুই দিকেই আলাদা বিষয়বস্তুর একটি ফাইল। আপনি প্রতিটি ফাইলের জন্য বেছে নেন: `keep group`, `use mine`, সঙ্গে `.json` ফাইলের জন্য `merge JSON` (একটি গভীর মার্জ, যেখানে সংঘাতপূর্ণ মানে গ্রুপ জেতে) অথবা `history.jsonl`-এর জন্য `merge by line`। যে MCP সার্ভারের সংজ্ঞা আলাদা, সেটিও একটি দ্বন্দ্ব, `keep group` বা `use mine` সহ। তালিকার প্রথম পছন্দটি Enter চাপলে ডিফল্ট; `--yes` সবসময় গ্রুপের কপি রাখে।

চেকলিস্টটি নিজেই `tui.MigratePicker`, যা `ccam add`-এর ব্যবহৃত `checklist` কম্পোনেন্টের ওপর তৈরি; `migrate.Found` lstat কল ও MCP সার্ভার তালিকা দিয়ে অনুপস্থিত আইটেমগুলো চিহ্নিত করে।

### MCP সাইন-ইন

যে MCP সার্ভারগুলো OAuth দিয়ে সাইন ইন করে, সেগুলো তাদের টোকেন সার্ভার তালিকায় নয়, লগইন ফাইলে রাখে, `.credentials.json`-এ `mcpOAuth`-এর নিচে। Claude Code প্রতিটি এন্ট্রির কী রাখে `<name>|` সঙ্গে সার্ভারের `type`, `url` ও `headers`-এর একটি `sha256` হ্যাশের প্রথম ১৬টি অক্ষর। কী-টি ফোল্ডারের ওপর নির্ভর করে না, তাই সার্ভারের সংজ্ঞা হুবহু এক থাকলে কপি করা এন্ট্রি অন্য অ্যাকাউন্টেও কাজ করে। তাই `internal/migrate/logins.go`-এর `buildLogins` শুধু সেই সার্ভারগুলোর এন্ট্রি কপি করে যেগুলো শেষ পর্যন্ত গ্রুপে উৎসের সংজ্ঞায় থাকে: নতুন, আগে থেকেই একই, অথবা এমন দ্বন্দ্ব যেখানে আপনি `use mine` বেছেছেন। প্লাগইন সার্ভারের সাইন-ইন (`plugin:` দিয়ে শুরু হওয়া কী) সঙ্গে যায় যখন আপনি এক্সটেনশন কপি করেন এবং উৎসে `plugins/` থাকে।

<!-- motion:migrate-signins -->

এন্ট্রিগুলো প্রতিটি প্রাপক অ্যাকাউন্টের নিজস্ব `.credentials.json`-এ যায়, `fsutil.AtomicWrite` দিয়ে মোড `0600`-এ লেখা হয়। এগুলো কখনও গ্রুপে যায় না। কোনো অ্যাকাউন্ট এগুলো পায় শুধু যদি সেটি লক্ষ্য গ্রুপের মাধ্যমে `mcp` (অথবা প্লাগইন সার্ভারের জন্য `extensions`) শেয়ার করে, এবং শুধু তার না থাকা কী-গুলো যোগ হয়: তার আগে থেকে থাকা, হয়তো নতুনতর, এন্ট্রিগুলো রেখে দেওয়া হয়। পরিকল্পনার `sign-ins` লাইন সার্ভারগুলো, প্রাপক অ্যাকাউন্ট ও বাদ পড়াগুলোর তালিকা দেয়। ccam একটি অ্যাকাউন্ট বাদ দেয়, এবং কারণসহ তার নাম জানায়, যখন সেটি:

- `ccam import` দিয়ে নিজের জায়গায় নিবন্ধিত, কারণ সাধারণ `claude` হয়তো তার লগইন ফাইল ব্যবহার করছে;
- চলছে, কারণ Claude Code তার লগইন ফাইল আবার লেখে (লেখার ঠিক আগে আবার যাচাই করা হয়, `started meanwhile` হিসেবে জানানো হয়);
- এখনও লগইন করা হয়নি;
- লগইন macOS Keychain-এ রাখে, যেখানে ccam লেখে না।

বাদ পড়া অ্যাকাউন্টে `/mcp` দিয়ে সার্ভারগুলোতে সাইন ইন করুন। যে অ্যাকাউন্ট পরে MCP সার্ভার শেয়ার করা শুরু করে, আপনি আবার `ccam migrate` চালালে সেটি সাইন-ইনগুলো পায়; সারাংশে তা লেখা থাকে ``Then run `ccam migrate` again to give them the MCP sign-ins too.`` দিয়ে।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- উৎস শুধু পড়া হয়। এর কোনো ফাইল সরানো, লিংক করা বা বদলানো হয় না, এবং উৎসের লগইন কখনও লেখা হয় না, macOS-সহ, যেখানে সেটি Keychain থেকে পড়া হতে পারে।
- কোনো গ্রুপ ফাইল প্রতিস্থাপনের আগে, `use mine`, মার্জ, বা বিদ্যমান `mcp.json`-এ নতুন MCP সার্ভারের কারণে, গ্রুপের কপি `<group>/.ccam-backup/<time>/`-এ সংরক্ষিত হয়। ফোল্ডারের নাম প্রতিটি রানে অনন্য, তাই একই সেকেন্ডে দুটি রান একে অপরের ব্যাকআপ মুছে লেখে না। সারাংশে পাথটি প্রিন্ট হয়।
- `--yes` দিলে প্রতিটি দ্বন্দ্বে গ্রুপের কপি থাকে, এবং ccam জানায় এর বদলে কীভাবে বেছে নেবেন: `--yes` ছাড়া একই কমান্ড চালান।
- কিছুই না বদলালে ccam প্রিন্ট করে `Nothing to copy.`; প্রতিটি দ্বন্দ্বে গ্রুপের কপি থাকলে `Nothing to copy: the group keeps its own copies.`
- কপি মাঝপথে ব্যর্থ হলে ccam থেমে যায় `The source folder was not changed. Group files replaced so far are saved in` ও তারপর ব্যাকআপ ফোল্ডার দিয়ে।
- উৎস না থাকলে শেষ হয় `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>` দিয়ে।

## সম্পর্কিত

- [MCP সার্ভার, সেটিংস ও ইতিহাস শেয়ার করুন](/docs/sharing)
- [যাচাই ও মেরামত](/features/doctor)
- [একাধিক Claude Code অ্যাকাউন্ট](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
