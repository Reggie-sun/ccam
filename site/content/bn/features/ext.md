---
title: "এক্সটেনশন পরীক্ষা"
description: "ccam ext দেখায় প্রতিটি অ্যাকাউন্ট কোন skills, commands, plugins, MCP সার্ভার ও hooks লোড করে, প্রতিটি কোথা থেকে আসে, এবং কী ভাঙা আছে।"
updated: 2026-10-26
order: 14
translation: machine
sourceHash: 7ca357a7f38a
---

`ccam ext` এই প্রশ্নের উত্তর দেয়: "এই skill (বা MCP সার্ভার) একটি অ্যাকাউন্টে কাজ করে কিন্তু অন্যটিতে করে না কেন?"। এটি প্রতিটি অ্যাকাউন্টের জন্য Claude Code যে এক্সটেনশনগুলো লোড করে তার তালিকা দেয়, এবং প্রতিটি কোথা থেকে আসে: একটি শেয়ার গ্রুপ, অ্যাকাউন্টের নিজস্ব ফোল্ডার, বা আপনি যে প্রকল্পে আছেন। টেবিলের নিচে এটি সেই সমস্যাগুলোর তালিকা দেয় যা কোনো এক্সটেনশনকে আপনার প্রত্যাশামতো লোড হতে বাধা দেয়। এটি কিছুই বদলায় না।

## ব্যবহার

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

প্রতিটি ঘর ওই অ্যাকাউন্টে ওই আইটেমের সোর্স জানায়: শেয়ার গ্রুপ (`default`), অ্যাকাউন্টের নিজস্ব ফোল্ডারে থাকলে `account`, অ্যাকাউন্টের শুধু এই ফোল্ডারের জন্য থাকা MCP সার্ভারের ক্ষেত্রে `local`, অথবা অ্যাকাউন্টে না থাকলে `-`। ইনস্টল করা কিন্তু চালু না করা plugin তার সোর্সের পরে `(off)` দেখায়। `!` সমস্যাযুক্ত আইটেম চিহ্নিত করে, যার ব্যাখ্যা টেবিলের নিচে থাকে; যে সমস্যাগুলো কোনো অ্যাকাউন্টের নয়, যেমন শুধু সাধারণ `claude`-এর কাছে থাকা সার্ভার, সেগুলো কোনো নাম ছাড়া সবশেষে আসে। `ccam ext --json` একই রিপোর্ট প্রিন্ট করে (`accounts`, প্রকল্প, `problems`), যেখানে Claude Code নেটওয়ার্কের মাধ্যমে যে MCP সার্ভারগুলোতে পৌঁছায় সেগুলোতে `remote` থাকে; কোনো অ্যাকাউন্ট না থাকলে এটি `accounts: []`-সহ রিপোর্টটি প্রিন্ট করে।

## এটি কী পড়ে

| ধরন | কোথায় | দেখানো সোর্স |
| --- | --- | --- |
| Skills | `skills/`-এর প্রতিটি ফোল্ডার যাতে `SKILL.md` আছে | যে গ্রুপের দিকে `skills/` লিংক করে, অথবা `account` |
| Commands, agents | `commands/` ও `agents/`-এর ভেতরের `.md` ফাইল, সাবফোল্ডার `folder/name` হিসেবে | গ্রুপ, অথবা `account` |
| Plugins | `plugins/installed_plugins.json`, অ্যাকাউন্টের ও প্রকল্পের সেটিংসের `enabledPlugins` দিয়ে চালু বা বন্ধ; শুধু একটি প্রকল্পের জন্য ইনস্টল করা plugin কেবল সেই প্রকল্পেই দেখা যায় | গ্রুপ, অথবা `account`; চালু না থাকলে `(off)` |
| MCP সার্ভার | গ্রুপের `mcp.json` (MCP শেয়ার করা থাকলে), অ্যাকাউন্টের `.claude.json` (এর নিজস্ব সার্ভার, এবং যেগুলো এর শুধু এই ফোল্ডারের জন্য আছে), এবং এর plugin-গুলোর সার্ভার, `plugin:<plugin>:<server>` হিসেবে | গ্রুপ, `account` অথবা `local` |
| Hooks | `settings.json`-এর `hooks`, প্রতি ইভেন্টে গোনা | গ্রুপ, অথবা `account` |
| এই প্রকল্প | বর্তমান ফোল্ডারের `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` ও `.mcp.json` | একবার দেখানো হয়; যে অ্যাকাউন্ট ফোল্ডারটিকে বিশ্বাস করেনি, সেটি এর `.mcp.json` সার্ভারগুলোর জন্য একটি সমস্যা পায় |

## এটি কোন সমস্যাগুলো জানায়

- **ভাঙা লিংক**: `skills/`, `commands/` বা `agents/`-এ এমন লিংক যার লক্ষ্য আর নেই, অথবা এমন শেয়ার করা ফোল্ডার (`skills`, `commands`, `agents`, `plugins`) যা অনুপস্থিত জায়গায় লিংক করে (`ccam doctor --fix` দিয়ে মেরামত করুন, যা অনুপস্থিত গ্রুপ ফোল্ডার আবার তৈরি করে)।
- **`SKILL.md` ছাড়া skill ফোল্ডার**: Claude Code এটি এড়িয়ে যায়।
- **দুই জায়গায় আলাদাভাবে সংজ্ঞায়িত MCP সার্ভার**: গ্রুপের `mcp.json` ও অ্যাকাউন্টের `.claude.json`, অথবা এদের একটি ও প্রকল্পের `.mcp.json`। দুই জায়গায় হুবহু এক সংজ্ঞা (যা `mcp_strategy = "sync"` লেখে) কোনো সমস্যা নয়। কোনো অ্যাকাউন্টের নিজস্ব সার্ভার `ccam mcp import <account> <server>` দিয়ে গ্রুপে সরান।
- **বৈধ JSON নয় এমন `installed_plugins.json`**, এমন plugin যার নাম `enabledPlugins`-এ নেই, এবং এমন plugin যার ফোল্ডার নেই।
- **শুধু একটি অ্যাকাউন্ট যে সার্ভার লোড করে**: শুধু এই অ্যাকাউন্টের `.claude.json`-এ থাকা কোনো সার্ভার, যা তার গ্রুপের অন্য অ্যাকাউন্টগুলো লোড করে না, এবং এই ফোল্ডারের জন্য অ্যাকাউন্টের local scope-এ থাকা কোনো সার্ভার (`-s user` না দিলে `claude mcp add` যা করে)। প্রতিটির সঙ্গে সেটি শেয়ার করার `ccam mcp import` কমান্ড থাকে। সাধারণ `claude`-এর `~/.claude.json`-এ থাকা যে সার্ভারগুলো কোনো অ্যাকাউন্ট লোড করে না, সেগুলোর সঙ্গে `ccam migrate --mcp` কমান্ড থাকে।
- **অনুমোদনের অপেক্ষায় থাকা প্রকল্প সার্ভার**: অ্যাকাউন্টটি এখনও ফোল্ডারটিকে বিশ্বাস করেনি, তাই Claude Code একবার জিজ্ঞেস করে, এবং ততক্ষণ প্রকল্পের `.mcp.json` সার্ভারগুলো অপেক্ষা করে। সমস্যাটি জানায় সেগুলো কোথায় অনুমোদন করতে হবে: সেটিংস গ্রুপের অ্যাকাউন্টগুলোর, অথবা এই অ্যাকাউন্টের `settings.json`-এ `"enabledMcpjsonServers"`। কোনো git রিপোজিটরির মূল worktree-কে দেওয়া বিশ্বাস তার অন্য worktree-গুলোতেও প্রযোজ্য।
- **সাইন-ইন ও connector**: এমন MCP সার্ভার যার জন্য এই অ্যাকাউন্টে সাইন-ইন দরকার, অথচ অন্য একটি অ্যাকাউন্ট তাতে সাইন ইন করা (একটি সাইন-ইন শুধু একটি অ্যাকাউন্টের জন্যই খাটে; `/mcp` দিয়ে সাইন ইন করুন), এবং claude.ai-এর যে connector-গুলো অন্য একটি অ্যাকাউন্ট ব্যবহার করেছে কিন্তু এটি করেনি।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- `ccam ext` কিছুই বদলায় না এবং কোনো নেটওয়ার্ক অনুরোধ পাঠায় না। এটি সার্ভারের নাম ও hooks-এর সংখ্যা পড়ে, MCP সংজ্ঞায় থাকা `env` বা হেডারের মান কখনও নয়। সাইন-ইন যাচাইয়ের জন্য এটি প্রতিটি অ্যাকাউন্টের লগইনে থাকা MCP সাইন-ইনগুলোর নাম পড়ে (macOS-এ Keychain থেকে), টোকেন কখনও নয়, এবং শুধু সেই অ্যাকাউন্টগুলোর জন্য যেগুলো এমন কোনো সার্ভার লোড করে যা অন্য কোনো অ্যাকাউন্টও লোড করে।
- এটি ডিস্কে যা আছে তা-ই দেখায়। managed settings দিয়ে আটকানো সার্ভারও দেখা যায়; সেটিংসে বন্ধ রাখা plugin `(off)` দেখায়।
- কোনো প্রকল্পের `.claude` ফোল্ডার ও `.mcp.json` অন্তর্ভুক্ত করতে প্রকল্পের ফোল্ডারে এটি চালান।

## সম্পর্কিত

- [গ্রুপ দিয়ে শেয়ার করা](/features/share)
- [শেয়ার করা MCP সার্ভার](/features/mcp)
- [যাচাই ও মেরামত](/features/doctor)
