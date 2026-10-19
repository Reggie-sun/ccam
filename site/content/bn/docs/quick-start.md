---
title: "কুইক স্টার্ট"
metaTitle: "কুইক স্টার্ট: Claude Code অ্যাকাউন্ট যোগ ও বদল করুন"
description: "ccam add দিয়ে দুটি Claude Code অ্যাকাউন্ট যোগ করুন, সেগুলোকে ccwork-এর মতো শর্টকাট দিন, অ্যাকাউন্ট বদল করুন, এবং কী শেয়ার হবে তা বেছে নিন।"
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 745a9904e32a
---

`ccam add` দিয়ে প্রতিটি অ্যাকাউন্ট যোগ করুন, তারপর `ccam <name>` দিয়ে বা নিজের দেওয়া নামের একটি শর্টকাট, যেমন `ccwork`, দিয়ে Claude Code চালু করুন। বেশিরভাগ মানুষের এটুকুই দরকার।

## ১. আপনার প্রথম অ্যাকাউন্ট যোগ করুন

```sh
ccam add
```

উইজার্ডটি `personal` বা `work`-এর মতো একটি নাম চায়, জিজ্ঞেস করে অ্যাকাউন্টটি আপনার অন্য অ্যাকাউন্টগুলোর সঙ্গে কোন কোন আইটেম শেয়ার করবে, এবং একটি শর্টকাট কমান্ড তৈরি করবে কি না (ধাপ ৪ দেখুন)। তারপর এটি `claude auth login` চালায়, যা সাইন ইন করার জন্য ব্রাউজার খোলে। শেষ হলে আপনি দেখবেন:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## ২. আপনার আগের সেটিংস নিয়ে আসুন

আগে Claude Code ব্যবহার করে থাকলে আপনার সেটিংস, `CLAUDE.md`, এক্সটেনশন ও MCP সার্ভার `~/.claude`-এ আছে। সেগুলো শেয়ার করা সেটিংসে কপি করুন, যাতে প্রতিটি অ্যাকাউন্ট সেগুলো পায়:

```sh
ccam migrate
```

ccam আগে দেখায় কী কপি করবে, তারপর জিজ্ঞেস করে। এটি `~/.claude` শুধু পড়ে, কখনও কোনো লগইন কপি করে না, এবং আপনি `--history` যোগ না করলে কথোপকথনের ইতিহাস বাদ রাখে। বিস্তারিত জানতে [আগের সেটআপ কপি করা](/docs/sharing#copy-an-existing-setup) দেখুন।

## ৩. দ্বিতীয় অ্যাকাউন্ট যোগ করুন

```sh
ccam add work
```

ব্রাউজারে অন্য অ্যাকাউন্টটি দিয়ে সাইন ইন করুন।

## ৪. অ্যাকাউন্ট বদল করুন

`ccam work` work অ্যাকাউন্ট হিসেবে Claude Code চালু করে। আরও ছোট কিছু চাইলে অ্যাকাউন্টটিকে তার নিজের একটি শর্টকাট কমান্ড দিন, যেমন `ccwork`। শেয়ার করা আইটেমগুলোর পরে উইজার্ড জিজ্ঞেস করে:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` চাপুন, তারপর একটি নাম টাইপ করুন, অথবা প্রস্তাবিত নামটি ব্যবহার করতে Enter চাপুন:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

লগইনের পরে ccam এটি নিশ্চিত করে:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

প্রশ্নটি এড়িয়ে গিয়ে থাকলে যেকোনো সময় একটি শর্টকাট যোগ করতে পারেন:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

এখন এগুলোর যেকোনোটি work অ্যাকাউন্ট হিসেবে Claude Code চালু করে:

```sh
ccwork
ccam work
ccam 2
```

শুধু `ccam` চালালে একটি বাছাইয়ের তালিকা খোলে: অ্যারো কি ও Enter ব্যবহার করুন, অথবা অ্যাকাউন্টের নম্বর চাপুন।

শর্টকাট বা অ্যাকাউন্টের নামের পরে যা কিছু লেখেন তা অপরিবর্তিত অবস্থায় Claude Code-এ যায়, তাই আপনার চেনা অপশনগুলোই ব্যবহার করতে থাকুন:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## ৫. কী শেয়ার হবে তা বেছে নিন

ডিফল্টভাবে MCP সার্ভার, `settings.json`, গ্লোবাল `CLAUDE.md`, এবং agents, commands, skills ও plugins অ্যাকাউন্টগুলোর মধ্যে শেয়ার হয়। কথোপকথনের ইতিহাস শেয়ার হয় না। অ্যাকাউন্টগুলোর মধ্যে কথোপকথন চালিয়ে যেতে চাইলে দুটি অ্যাকাউন্টেই এটি চালু করুন:

```sh
ccam share personal --history on
ccam share work --history on
```

কিছু বদলানোর আগে ccam দেখায় এটি কী সরাবে। বিস্তারিত জানতে [শেয়ার করা](/docs/sharing) দেখুন।

## ৬. যখন একটি অ্যাকাউন্ট সীমায় পৌঁছায়

```sh
ccam usage
ccam next --continue
```

`ccam usage` সব অ্যাকাউন্টের ৫ ঘণ্টা ও সাপ্তাহিক ব্যবহার দেখায়। `ccam next --continue` একই কথোপকথন পরের অ্যাকাউন্টে খোলে। দেখুন [ব্যবহারের সীমা](/docs/usage-limits)।

## আগে থেকেই ~/.claude ব্যবহার করছেন?

আবার সাইন ইন না করে সেটিকে নিবন্ধন করুন:

```sh
ccam import
```

ফোল্ডারটি যেখানে আছে সেখানেই থাকে, এবং সাধারণ `claude` আগের মতোই কাজ করে।
