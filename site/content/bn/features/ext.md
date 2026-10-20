---
title: "এক্সটেনশন পরীক্ষা"
description: "ccam ext দেখায় প্রতিটি অ্যাকাউন্ট কোন skills, commands, plugins, MCP সার্ভার ও hooks লোড করে, প্রতিটি কোথা থেকে আসে, এবং কী ভাঙা আছে।"
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` এই প্রশ্নের উত্তর দেয়: "এই skill (বা MCP সার্ভার) একটি অ্যাকাউন্টে কাজ করে কিন্তু অন্যটিতে করে না কেন?"। এটি প্রতিটি অ্যাকাউন্টের জন্য Claude Code যে এক্সটেনশনগুলো লোড করে তার তালিকা দেয়, এবং প্রতিটি কোথা থেকে আসে: একটি শেয়ার গ্রুপ, অ্যাকাউন্টের নিজস্ব ফোল্ডার, বা আপনি যে প্রকল্পে আছেন। টেবিলের নিচে এটি সেই সমস্যাগুলোর তালিকা দেয় যা কোনো এক্সটেনশনকে আপনার প্রত্যাশামতো লোড হতে বাধা দেয়। এটি শুধু পড়ে।

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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

প্রতিটি ঘর ওই অ্যাকাউন্টে ওই আইটেমের সোর্স জানায়: শেয়ার গ্রুপ (`default`), অ্যাকাউন্টের নিজস্ব ফোল্ডারে থাকলে `account`, অথবা অ্যাকাউন্টে না থাকলে `-`। `!` সমস্যাযুক্ত আইটেম চিহ্নিত করে, যার ব্যাখ্যা টেবিলের নিচে থাকে।

## এটি কী পড়ে

| ধরন | কোথায় | দেখানো সোর্স |
| --- | --- | --- |
| Skills | `skills/`-এর প্রতিটি ফোল্ডার যাতে `SKILL.md` আছে | যে গ্রুপের দিকে `skills/` লিংক করে, অথবা `account` |
| Commands, agents | `commands/` ও `agents/`-এর ভেতরের `.md` ফাইল, সাবফোল্ডার `folder/name` হিসেবে | গ্রুপ, অথবা `account` |
| Plugins | `plugins/installed_plugins.json` | গ্রুপ, অথবা `account` |
| MCP সার্ভার | গ্রুপের `mcp.json` (MCP শেয়ার করা থাকলে) ও অ্যাকাউন্টের `.claude.json` | গ্রুপ, অথবা `account` |
| Hooks | `settings.json`-এর `hooks`, প্রতি ইভেন্টে গোনা | গ্রুপ, অথবা `account` |
| এই প্রকল্প | বর্তমান ফোল্ডারের `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` ও `.mcp.json` | একবার দেখানো হয়, প্রতিটি অ্যাকাউন্টের জন্য |

## এটি কোন সমস্যাগুলো জানায়

- **ভাঙা লিংক**: `skills/`, `commands/` বা `agents/`-এ এমন লিংক যার লক্ষ্য আর নেই, অথবা এমন শেয়ার করা ফোল্ডার (`skills`, `commands`, `agents`, `plugins`) যা অনুপস্থিত জায়গায় লিংক করে (`ccam doctor --fix` দিয়ে মেরামত করুন)।
- **`SKILL.md` ছাড়া skill ফোল্ডার**: Claude Code এটি এড়িয়ে যায়।
- **দুই জায়গায় আলাদাভাবে সংজ্ঞায়িত MCP সার্ভার**: গ্রুপের `mcp.json` ও অ্যাকাউন্টের `.claude.json`, অথবা এদের একটি ও প্রকল্পের `.mcp.json`। দুই জায়গায় একই সেটিংস (যা `mcp_strategy = "sync"` লেখে) কোনো সমস্যা নয়। কোনো অ্যাকাউন্টের নিজস্ব সার্ভার `ccam mcp import <account> <server>` দিয়ে গ্রুপে সরান।
- **বৈধ JSON নয় এমন `installed_plugins.json`।**

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- `ccam ext` কিছুই বদলায় না এবং কোনো নেটওয়ার্ক অনুরোধ পাঠায় না। এটি সার্ভারের নাম ও hooks-এর সংখ্যা পড়ে, MCP সংজ্ঞায় থাকা `env` বা হেডারের মান কখনও নয়।
- এটি ডিস্কে যা আছে তা-ই দেখায়। ইনস্টল করা কিন্তু সেটিংসে বন্ধ রাখা plugin, বা managed settings দিয়ে আটকানো সার্ভারও দেখা যায়।
- কোনো প্রকল্পের `.claude` ফোল্ডার ও `.mcp.json` অন্তর্ভুক্ত করতে প্রকল্পের ফোল্ডারে এটি চালান।

## সম্পর্কিত

- [অ্যাকাউন্টগুলোর মধ্যে শেয়ার করা](/features/share)
- [শেয়ার করা MCP সার্ভার](/features/mcp)
- [যাচাই ও মেরামত](/features/doctor)
