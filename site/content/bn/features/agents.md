---
title: "প্রতিটি অ্যাকাউন্টের ব্যাকগ্রাউন্ড সেশন"
description: "ccam agents আপনার সব অ্যাকাউন্টের Claude Code ব্যাকগ্রাউন্ড এজেন্ট ও খোলা সেশন একটি টেবিলে দেখায়, আর ccam attach id দিয়ে একটি সেশন খোলে।"
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code ব্যাকগ্রাউন্ড সেশন (`claude --bg`, agent view) প্রতি অ্যাকাউন্টে আলাদা রাখে: `claude agents` শুধু সেই অ্যাকাউন্টের সেশন দেখায় যা হিসেবে এটি চলছে। `ccam agents` একসঙ্গে প্রতিটি অ্যাকাউন্টকে জিজ্ঞেস করে এবং সবগুলো একসঙ্গে তালিকায় দেখায়, আর `ccam attach <id>` একটি ব্যাকগ্রাউন্ড সেশন সেই অ্যাকাউন্ট হিসেবে খোলে যেটি সেশনটি শুরু করেছিল।

## ব্যবহার

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` হলো ব্যাকগ্রাউন্ড সেশনের অবস্থা (`working`, `blocked`, `done`, `failed`, `stopped`) এবং তার প্রসেস কী করছে (`busy`, `waiting`, `idle`); খোলা ইন্টারঅ্যাকটিভ সেশন দেখায় `open`।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--all` | যে ব্যাকগ্রাউন্ড সেশনগুলো শেষ, ব্যর্থ বা থামানো, সেগুলোও তালিকায় দেখায় | বন্ধ |
| `--json` | সেশনগুলো JSON হিসেবে প্রিন্ট করে, প্রতিটি যে `account`-এর তা সহ | বন্ধ |

`ccam attach <id> [claude args]` id-এর পরের আর্গুমেন্টগুলো `claude attach`-এ পাঠায়। `--dry-run` দিলে এটি কমান্ডটি প্রিন্ট করে।

## এটি কীভাবে কাজ করে

1. প্রতিটি অ্যাকাউন্টের জন্য ccam অ্যাকাউন্টের `CLAUDE_CONFIG_DIR` দিয়ে `claude agents --json` চালায় (দেওয়া থাকলে `--all`-সহ), সবগুলো একসঙ্গে, প্রত্যেকটি ১৫ সেকেন্ডের মধ্যে। এটি অ্যাকাউন্টের নিজস্ব সেশন তালিকা পড়ে এবং Claude Code-এর ব্যাকগ্রাউন্ড সার্ভিস চালু করে না।
2. প্রতিটি এন্ট্রি অ্যাকাউন্টের নাম পায়, এবং তালিকাটি শুরুর সময় অনুযায়ী সাজানো হয়, সবচেয়ে নতুনটি আগে। যে অ্যাকাউন্ট উত্তর দেয় না, তার নাম স্ট্যান্ডার্ড এররে জানানো হয় এবং সেটি বাদ পড়ে, তারপর ccam কোড 1 দিয়ে বেরিয়ে যায়, যাতে কোনো স্ক্রিপ্ট ব্যর্থতাকে খালি তালিকা বলে ভুল না করে।
3. `ccam attach` প্রতিটি অ্যাকাউন্টের ফোল্ডারে `jobs/<id>` খোঁজে, যেখানে Claude Code একটি ব্যাকগ্রাউন্ড সেশন রাখে, এবং সেই অ্যাকাউন্ট হিসেবে `claude attach <id>` চালায়, ঠিক যেভাবে `ccam <name>` Claude Code চালু করে।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- ব্যাকগ্রাউন্ড সেশনের id একটি অ্যাকাউন্টেরই: অন্য অ্যাকাউন্ট হিসেবে চালানো `claude attach` সেগুলো খুঁজে পায় না, আর `ccam attach` ঠিক এই সমস্যাটিই মেটায়।
- দুটি অ্যাকাউন্টে এমন সেশন থাকতে পারে যাদের id একইভাবে শুরু হয়; তখন `ccam attach` দুটোরই নাম জানায় এবং id-এর আরও অংশ চায়।
- ccam তার শেয়ার করা MCP সার্ভারগুলো একটি আর্গুমেন্টে `--mcp-config=<file>` দিয়ে Claude Code-কে দেয়, যা `claude agents` ও `claude attach` গ্রহণ করে; ব্যাকগ্রাউন্ড সেশন আবার চালু হলেও ওই ফাইলটিই ব্যবহার করতে থাকে।
- API key অ্যাকাউন্টের ব্যাকগ্রাউন্ড সার্ভিস তার key নেয় সেই শেল থেকে যা সার্ভিসটি চালু করেছিল। এমন অ্যাকাউন্টের প্রথম ব্যাকগ্রাউন্ড সেশনটি ccam-এর মাধ্যমে চালু করুন (`ccam gw --bg …`), যাতে সার্ভিসটি key পায়।

## সম্পর্কিত

- [Claude Code চালু করা](/features/launch)
- [স্ক্রিপ্টের জন্য অ্যাকাউন্ট বেছে নেওয়া](/features/pick)
- [ড্যাশবোর্ড](/features/dashboard)
