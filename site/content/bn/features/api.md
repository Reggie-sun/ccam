---
title: "API key অ্যাকাউন্ট"
description: "ccam add --api Anthropic-এর API বা কোনো গেটওয়ের জন্য API key অ্যাকাউন্ট যোগ করে; key চালুর সময় একটি ভেরিয়েবল বা কমান্ড থেকে পড়া হয়, কখনও সংরক্ষিত হয় না।"
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

একটি API key অ্যাকাউন্ট Claude সাবস্ক্রিপশন লগইনের বদলে একটি API key দিয়ে Claude Code চালু করে: Anthropic-এর API, অথবা একই Messages ফরম্যাট বোঝে এমন কোনো এন্ডপয়েন্ট, যেমন LiteLLM বা OpenRouter গেটওয়ে। ccam key সংরক্ষণ করে না। প্রতিবার অ্যাকাউন্ট চালু হলে এটি একটি এনভায়রনমেন্ট ভেরিয়েবল থেকে বা পাসওয়ার্ড ম্যানেজারের মতো কোনো কমান্ড থেকে key পড়ে, এবং শুধু ওই Claude Code প্রসেসকেই দেয়।

## ব্যবহার

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

একটি API key অ্যাকাউন্ট অন্য যেকোনো অ্যাকাউন্টের মতোই সেটিংস, `CLAUDE.md`, এক্সটেনশন ও MCP সার্ভার শেয়ার করে, এবং শর্টকাট, `ccam folder` ও `ccam ext`-এর সঙ্গে কাজ করে।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--api` | লগইনের বদলে একটি API key অ্যাকাউন্ট যোগ করে | বন্ধ |
| `--key-from env:NAME` | এনভায়রনমেন্ট ভেরিয়েবল `NAME` থেকে key পড়ে | |
| `--key-from command:COMMAND` | `COMMAND` চালায় (`sh -c` দিয়ে, বা Windows-এ `cmd /C` দিয়ে) এবং এটি যে প্রথম লাইন প্রিন্ট করে সেটি ব্যবহার করে | |
| `--base-url <url>` | এন্ডপয়েন্ট (`ANTHROPIC_BASE_URL`) | Anthropic-এর API |
| `--auth api-key\|bearer` | key-টি `ANTHROPIC_API_KEY` (`x-api-key`) হিসেবে, অথবা `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, যা গেটওয়ে সাধারণত চায়) হিসেবে পাঠায় | `api-key` |
| `--model <model>` | মডেল (`ANTHROPIC_MODEL`) | Claude Code-এর ডিফল্ট |

সেটিংসগুলো `config.toml`-এ সংরক্ষিত হয়:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## এটি কীভাবে কাজ করে

অ্যাকাউন্টটি প্রতিবার চালু হওয়ার সময়, সাধারণ যাচাইয়ের পরে এবং Claude Code শুরু হওয়ার আগে:

1. ccam key পড়ে: `env:` এনভায়রনমেন্ট থেকে, `command:` ৩০ সেকেন্ডের সীমার মধ্যে কমান্ডটি চালিয়ে। কোনো কমান্ড টার্মিনালে পাসওয়ার্ড চাইতে পারে; তার ত্রুটির আউটপুটও টার্মিনালে যায়। খালি key চালু হওয়া থামিয়ে দেয়, এমন একটি বার্তাসহ যা বলে key কোথায় থাকার কথা ছিল।
2. Claude Code কীভাবে সাইন ইন করবে বা কোন প্রোভাইডার ব্যবহার করবে তা ঠিক করে এমন প্রতিটি উত্তরাধিকারসূত্রে পাওয়া ভেরিয়েবল এটি সরিয়ে দেয়: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` ও `CLAUDE_CODE_USE_FOUNDRY`। অন্য কোনো কাজের জন্য আপনার শেলে সেট করা key এই অ্যাকাউন্টে পৌঁছাতে পারে না, এবং কোনো সাবস্ক্রিপশন লগইন এর key-কে ছাপিয়ে যেতে পারে না।
3. দেওয়া থাকলে এটি `ANTHROPIC_BASE_URL` ও `ANTHROPIC_MODEL` সেট করে, এবং key-টিকে `ANTHROPIC_API_KEY` বা `ANTHROPIC_AUTH_TOKEN` হিসেবে।

সাবস্ক্রিপশন লগইন থাকা অ্যাকাউন্টগুলো আগের মতোই উত্তরাধিকারসূত্রে পাওয়া এনভায়রনমেন্ট রাখে। ইন্টারঅ্যাকটিভ সেশনে Claude Code একবার জিজ্ঞেস করতে পারে এনভায়রনমেন্টে পাওয়া key ব্যবহার করবে কি না।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- key কখনও `config.toml`, কোনো ক্যাশ বা লগে লেখা হয় না, এবং `ccam --dry-run gw` এটি পড়ে না: ড্রাই রানে `command:` সোর্স চালানো হয় না।
- `ccam add --api` একটি `env:` সোর্স সঙ্গে সঙ্গে যাচাই করে এবং ভেরিয়েবলটি খালি হলে সতর্ক করে; `command:` সোর্স প্রথমবার চালুর সময় চলে।
- `ccam login` ও `ccam logout` API key অ্যাকাউন্ট প্রত্যাখ্যান করে, কারণ এর কোনো লগইন নেই। `ccam doctor` এর এন্ডপয়েন্ট ও key সোর্স দেখায়, এবং একটি `env:` ভেরিয়েবল খালি হলে সতর্ক করে।
- `ccam usage` এর জন্য `API` দেখায়। `ccam next --auto`, `--fallback` ও `ccam pick` শুধু সাবস্ক্রিপশন অ্যাকাউন্টগুলোর মধ্যে থেকে বেছে নেয়, কারণ API key-এর তুলনা করার মতো কোনো ৫ ঘণ্টা বা সাপ্তাহিক সীমা নেই।
- `env` ব্লক বা `apiKeyHelper` থাকা একটি শেয়ার করা `settings.json` গ্রুপের প্রতিটি অ্যাকাউন্টে প্রযোজ্য, API key অ্যাকাউন্টসহ। প্রোভাইডারের সেটিংস শেয়ার করা সেটিংসের বাইরে রাখুন।

## সম্পর্কিত

- [অ্যাকাউন্ট পরিচালনা](/features/accounts)
- [Claude Code চালু করা](/features/launch)
- [কমান্ড রেফারেন্স](/docs/commands)
