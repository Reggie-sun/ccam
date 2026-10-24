---
title: "Claude Code ব্যবহারের সীমায় পৌঁছালে কী করবেন"
description: "রিসেটের জন্য অপেক্ষা করুন, অথবা একই কথোপকথন এখনই অন্য Claude Code অ্যাকাউন্টে চালিয়ে যান। ৫ ঘণ্টা ও সাপ্তাহিক সীমা কীভাবে কাজ করে, ধাপে ধাপে।"
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 049a46684c7b
---

কোনো অ্যাকাউন্ট ব্যবহারের সীমায় (প্রায়ই rate limit বলা হয়) পৌঁছালে আপনার সামনে দুটি পথ: সীমা রিসেট হওয়া পর্যন্ত অপেক্ষা করুন, অথবা একই কথোপকথন এখনই আপনার অন্য একটি অ্যাকাউন্টে চালিয়ে যান। ccam-এ দ্বিতীয় পথটি একটি কমান্ড: `ccam next --auto --continue`।

## Claude Code-এর কী কী সীমা আছে?

Pro ও Max প্ল্যানে Claude Code-এর [পাঁচ ঘণ্টার সেশন সীমা ও একটি সাপ্তাহিক সীমা](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) আছে। দুটিই ওয়েব ও অ্যাপের Claude-এর সঙ্গে শেয়ার করা, এবং ব্যবহার একই সময়ে দুটি থেকেই কাটে। একটি শেষ হলে Claude Code "You've hit your session limit · resets 3:45pm"-এর মতো একটি [বার্তা দেখায়](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit)। মডেল বদলে লাভ হয় না, কারণ সীমা সব মডেলের ক্ষেত্রেই প্রযোজ্য।

## বিকল্প ১: রিসেটের জন্য অপেক্ষা করুন

আপনার একটিই অ্যাকাউন্ট থাকলে অপেক্ষা করুন। Claude Code-এর সাম্প্রতিক সংস্করণ [খোলা সেশনেই অপেক্ষা করে](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) এবং সীমা রিসেট হলে নিজে থেকেই কাজ চালিয়ে যায়। Claude Code-এর ভেতরে `/usage` দেখায় আপনার প্ল্যানের সীমার কতটা ব্যবহার হয়েছে।

## বিকল্প ২: এখনই অন্য অ্যাকাউন্টে চালিয়ে যান

আপনার দ্বিতীয় একটি অ্যাকাউন্ট থাকলে একই কথোপকথন নিয়ে সঙ্গে সঙ্গে এগিয়ে যেতে পারেন। একবার সেট আপ করুন:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

এখন কাজ করার সময় স্ট্যাটাস লাইন অ্যাকাউন্ট ও তার ৫ ঘণ্টার ব্যবহার দেখায়, তাই সীমায় পৌঁছানো আর অপ্রত্যাশিত থাকে না।

সীমায় পৌঁছালে Claude Code থেকে বেরিয়ে আসুন (`/exit`), তারপর:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

সবচেয়ে বেশি জায়গা থাকা অ্যাকাউন্টে কথোপকথনটি খুলুন:

```sh
ccam next --auto --continue
```

Claude Code `work` হিসেবে চালু হয়, বর্তমান ফোল্ডারের সবচেয়ে সাম্প্রতিক কথোপকথনসহ, এ পর্যন্ত বলা সবকিছু নিয়ে। ইতিহাস শেয়ার করা থাকলে (ওপরের `ccam group set default --add history`) এটি একটিই কথোপকথন থাকে; তা না থাকলে Claude Code কথোপকথনটি `work`-এ কপি করে, এবং ccam আপনাকে জানায় সেই প্রথম উত্তরে কতগুলো টোকেন আবার পড়া হবে।

Claude Code-এ [ccam প্লাগইন](/docs/install#claude-code-plugin) থাকলে সীমায় থেমে যাওয়া সেশন একটি নোটিফিকেশন দেখায় এবং পূর্ণ উইন্ডোটি রেকর্ড করে, তাই `ccam next --auto` রিসেটের আগে ওই অ্যাকাউন্টটি আবার বেছে নেয় না, আর `/ccam:switch` (বা `/ccam:switch work`) কথোপকথনটি `ccam next --continue`-এর জন্য টুকে রাখে, যাতে কিছুই আন্দাজ করতে না হয়।

**Team ও Enterprise সিট।** স্বয়ংক্রিয় বাছাই (`ccam next`, `--fallback`, `ccam pick`, ড্যাশবোর্ডের `n`) একই Claude Team বা Enterprise প্রতিষ্ঠানের সিটগুলোর মধ্যেও বদল করে। Anthropic-এর সহায়তা কেন্দ্রে বলা হয়েছে, "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team প্ল্যানের প্রতিটি সিট একজন নির্দিষ্ট ব্যক্তির নামে বরাদ্দ করা হয়, এবং কেবল সেই ব্যক্তিই সেটি ব্যবহার করতে পারেন) এবং "Users can't route their work through another user's account to use that user's remaining usage" (ব্যবহারকারীরা অন্য কোনো ব্যবহারকারীর অ্যাকাউন্টের মধ্য দিয়ে নিজেদের কাজ চালিয়ে তাঁর ব্যবহারের সীমার বাকি অংশ কাজে লাগাতে পারেন না) ([Team প্ল্যান](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise প্ল্যানের](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) ক্ষেত্রেও একই কথা)। তাই সিটগুলোর মধ্যে স্বয়ংক্রিয় বাছাই ব্যবহার করলে Anthropic-এর শর্তাবলি লঙ্ঘন হতে পারে। এর ওপর নির্ভর করার আগে আপনার প্রতিষ্ঠান ও Anthropic-এর কাছ থেকে জেনে নিন। ব্যক্তিগত প্ল্যানের (Pro, Max) কথা [প্রশ্নোত্তরের শর্তাবলি-বিষয়ক নোটে](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) বলা আছে।

## ccam কি নিজে থেকে অ্যাকাউন্ট বদল করে?

না। আপনি কাজ করার সময় ccam চলে না: এটি Claude Code চালু করে এবং নিজের জায়গায় সেটিকে বসিয়ে দেয়। কখন বদলাবেন তা আপনি ঠিক করেন, আর `ccam next --auto` সবচেয়ে কম ৫ ঘণ্টার ব্যবহারের অ্যাকাউন্ট বেছে নেয় এবং যেগুলো 100%-এ আছে, লগইন করা নেই, বা অন্য প্রতিষ্ঠানের, সেগুলো বাদ দেয়। `config.toml`-এ `relay = true` থাকলে ccam এর বদলে সেশনের পাশে অপেক্ষা করে, তাই `/ccam:switch` নিজে থেকেই কাজ চালিয়ে যায়, সেশনের অপশন, permission mode ও effort level বজায় রেখে। অন্য কিছু টুল এর বদলে চলমান সেশনের লগইন অদলবদল করে; [ccam ও অন্যান্য অ্যাকাউন্ট সুইচার](/compare/account-switchers) পাতায় পার্থক্যটি ব্যাখ্যা করা আছে।

## যা মনে রাখবেন

- যে ফোল্ডারে কাজ করছিলেন সেখানেই `ccam next --continue` চালান; কথোপকথন প্রকল্প ফোল্ডার অনুযায়ী সংরক্ষিত থাকে।
- একই কথোপকথনে একই সময়ে দুটি অ্যাকাউন্ট চালাবেন না। এমন চেষ্টা করলে ccam সতর্ক করে।
- শেয়ার করা ইতিহাস না থাকলে কথোপকথনটি অন্য অ্যাকাউন্টে কপি হয় এবং দুটি কপি যার যার পথে চলে; `/rewind` যে চেকপয়েন্ট থেকে কোড ফিরিয়ে আনে, সেগুলো প্রথম অ্যাকাউন্টেই থেকে যায়।
- নতুন অ্যাকাউন্টে প্রথম বার্তা স্বাভাবিকের চেয়ে বেশি সীমা খরচ করতে পারে, কারণ সেই অ্যাকাউন্ট এখনো কথোপকথনটি ক্যাশ করেনি।
- এভাবে অ্যাকাউন্ট বদল করা যাবে কি না, তা প্রতিটি অ্যাকাউন্টের শর্তের ওপর নির্ভর করে।

`ccam usage` ও `ccam next`-এর সব অপশনের জন্য [ব্যবহারের সীমা](/docs/usage-limits) দেখুন।
