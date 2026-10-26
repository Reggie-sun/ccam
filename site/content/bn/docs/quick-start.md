---
title: "কুইক স্টার্ট"
metaTitle: "কুইক স্টার্ট: Claude Code অ্যাকাউন্ট যোগ ও বদল করুন"
description: "ccam add দিয়ে দুটি Claude Code অ্যাকাউন্ট যোগ করুন, সেগুলোকে ccwork-এর মতো শর্টকাট দিন, অ্যাকাউন্ট বদল করুন, এবং কী শেয়ার হবে তা বেছে নিন।"
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

`ccam add` দিয়ে প্রতিটি অ্যাকাউন্ট যোগ করুন, তারপর `ccam <name>` দিয়ে বা নিজের দেওয়া নামের একটি শর্টকাট, যেমন `ccwork`, দিয়ে Claude Code চালু করুন। বেশিরভাগ মানুষের এটুকুই দরকার।

আগে থেকেই Claude Code ব্যবহার করেন? শুরু করার আগে এই পাতার শেষে [আগে থেকেই ~/.claude ব্যবহার করছেন?](#already-using-claude) পড়ুন।

## ১. আপনার প্রথম অ্যাকাউন্ট যোগ করুন

```sh
ccam add
```

উইজার্ডটি `personal` বা `work`-এর মতো একটি নাম চায়, জিজ্ঞেস করে অ্যাকাউন্টটি আপনার অন্য অ্যাকাউন্টগুলোর সঙ্গে কোন কোন আইটেম শেয়ার করবে, এবং একটি শর্টকাট কমান্ড তৈরি করবে কি না (ধাপ ৩ দেখুন)। তারপর এটি `claude auth login` চালায়, যা সাইন ইন করার জন্য ব্রাউজার খোলে। রিমোট মেশিনে আপনার ব্রাউজার শেষে একটি `http://localhost:…/callback` পেজে যেতে পারে, যা লোড হয় না; ওই ঠিকানাটি টার্মিনালে পেস্ট করুন, ccam সাইন-ইন শেষ করবে। সাইন-ইন সম্পূর্ণ হলে আপনি দেখবেন:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## ২. আপনার আগের সেটিংস নিয়ে আসুন

আগে Claude Code ব্যবহার করে থাকলে আপনার সেটিংস, `CLAUDE.md` ও এক্সটেনশন `~/.claude`-এ আছে, আর আপনার MCP সার্ভার `~/.claude.json`-এ। সেগুলো শেয়ার করা `default` গ্রুপে কপি করুন, যাতে প্রতিটি অ্যাকাউন্ট সেগুলো পায়:

```sh
ccam migrate
```

ccam আগে দেখায় কী কপি করবে, তারপর জিজ্ঞেস করে। এটি `~/.claude` শুধু পড়ে, কখনও কোনো লগইন কপি করে না, এবং আপনি `--history` যোগ না করলে কথোপকথনের ইতিহাস বাদ রাখে। বিস্তারিত জানতে [আগের সেটআপ কপি করা](/docs/sharing#copy-an-existing-setup) দেখুন।

## ৩. দ্বিতীয় অ্যাকাউন্ট যোগ করুন

```sh
ccam add work
```

উইজার্ডটি নাম জিজ্ঞেস করে না, কারণ আপনি কমান্ড লাইনেই তা দিয়েছেন, এবং `default` চেক করা অবস্থায় যোগ দেওয়ার গ্রুপগুলো দেখায়। এবার অ্যাকাউন্টটিকে তার নিজের একটি শর্টকাট কমান্ড দিন, যেমন `ccwork`, যাতে একটি শব্দেই এটি চালু করতে পারেন। গ্রুপগুলোর পরে উইজার্ড জিজ্ঞেস করে:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` চাপুন, তারপর একটি নাম টাইপ করুন, অথবা প্রস্তাবিত নামটি ব্যবহার করতে Enter চাপুন:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

তারপর ব্রাউজারে অন্য অ্যাকাউন্টটি দিয়ে সাইন ইন করুন। ccam অ্যাকাউন্ট ও তার শর্টকাট নিশ্চিত করে:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## ৪. অ্যাকাউন্ট বদল করুন

`ccam work` work অ্যাকাউন্ট হিসেবে Claude Code চালু করে। শর্টকাটের প্রশ্নটি এড়িয়ে গিয়ে থাকলে যেকোনো সময় একটি শর্টকাট যোগ করতে পারেন:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

এখন এই কমান্ডগুলোর প্রতিটি work অ্যাকাউন্ট হিসেবে Claude Code চালু করে (`2` হলো `ccam list`-এ এর নম্বর):

```sh
ccwork
ccam work
ccam 2
```

শুধু `ccam` চালালে [ড্যাশবোর্ড](/features/dashboard) খোলে: অ্যারো কি বা অ্যাকাউন্টের নম্বর দিয়ে একটি অ্যাকাউন্টে যান, তারপর Enter চাপুন।

শর্টকাট বা অ্যাকাউন্টের নামের পরে যা কিছু লেখেন তা অপরিবর্তিত অবস্থায় Claude Code-এ যায়, তাই আপনার চেনা অপশনগুলোই ব্যবহার করতে থাকুন:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## ৫. কী শেয়ার হবে তা বেছে নিন

ডিফল্টভাবে MCP সার্ভার, `settings.json`, গ্লোবাল `CLAUDE.md` ও এক্সটেনশন (agents, commands, skills ও plugins) `default` গ্রুপের মাধ্যমে অ্যাকাউন্টগুলোর মধ্যে শেয়ার হয়। কথোপকথনের ইতিহাস শেয়ার হয় না। অ্যাকাউন্টগুলোর মধ্যে কথোপকথন চালিয়ে যেতে চাইলে এটি গ্রুপে যোগ করুন:

```sh
ccam group set default --add history
```

কিছু বদলানোর আগে ccam দেখায় প্রতিটি অ্যাকাউন্টের জন্য এটি কী সরাবে। বিস্তারিত জানতে [গ্রুপ দিয়ে শেয়ার করা](/features/share) দেখুন।

## ৬. যখন একটি অ্যাকাউন্ট সীমায় পৌঁছায়

```sh
ccam usage
ccam next --continue
```

`ccam usage` সব অ্যাকাউন্টের ৫ ঘণ্টা ও সাপ্তাহিক ব্যবহার দেখায়। `ccam next --continue` একই কথোপকথন পরের অ্যাকাউন্টে খোলে। স্ক্রিপ্টে `ccam --fallback work -p "…"` নিজে থেকেই বদল করে। দেখুন [ব্যবহারের সীমা](/docs/usage-limits)।

**Team ও Enterprise সিট।** স্বয়ংক্রিয় অ্যাকাউন্ট বাছাই (`ccam next`, `--fallback`, `ccam pick`, ড্যাশবোর্ডের `n` এবং আপনি চালু করলে [`auto_switch`](/features/next#auto-switch)) একই Claude Team বা Enterprise প্রতিষ্ঠানের সিটগুলোর মধ্যেও বদল করে, যা Anthropic-এর শর্তাবলি লঙ্ঘন করতে পারে। এর ওপর নির্ভর করার আগে [প্রশ্নোত্তর](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) পড়ুন এবং আপনার প্রতিষ্ঠান ও Anthropic-এর কাছ থেকে জেনে নিন।

## আগে থেকেই ~/.claude ব্যবহার করছেন?

আবার সাইন ইন না করে সেটিকে নিবন্ধন করুন:

```sh
ccam import
```

ফোল্ডারটি যেখানে আছে সেখানেই থাকে, এবং সাধারণ `claude` আগের মতোই কাজ করে।
