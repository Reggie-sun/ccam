---
title: "ফোল্ডার অনুযায়ী অ্যাকাউন্ট"
description: "অ্যাকাউন্টের নাম ছাড়া Claude Code চালু করলে ccam folder একটি ফোল্ডার ও তার নিচের সবকিছুতে একটি অ্যাকাউন্ট ব্যবহার করায়, যেমন আপনার কোম্পানির কোডের জন্য work।"
updated: 2026-10-26
order: 11
translation: machine
sourceHash: 25f82da64580
---

`ccam folder work ~/code/company` চালানোর পর `~/code/company`-এ বা তার সাবফোল্ডারগুলোতে অ্যাকাউন্টের নাম ছাড়া যেকোনো চালু করা Claude Code-কে `work` হিসেবে চালু করে। এর বাইরে আগের মতোই ডিফল্ট অ্যাকাউন্ট চালু হয়। `ccam personal` বা `ccpersonal`-এর মতো শর্টকাট দিয়ে অ্যাকাউন্টের নাম দিলে সবসময় সেই অ্যাকাউন্টই চালু হয়।

## ব্যবহার

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

এই চালু করাগুলো ফোল্ডারের অ্যাকাউন্ট ব্যবহার করে:

- `ccam -p "…"` এবং অন্য যেকোনো চালু করা যার প্রথম আর্গুমেন্ট Claude Code-এর একটি অপশন
- নাম ছাড়া `ccam run`, এবং `ccam --dry-run`
- `ccam --fallback -p "…"`, তার প্রথম চেষ্টার জন্য
- ড্যাশবোর্ড, একা `ccam`, যা ফোল্ডারের অ্যাকাউন্টের ওপর কার্সর রেখে শুরু হয়, এমনকি যখন আপনি একটি সিম্বলিক লিংকের মধ্য দিয়ে ফোল্ডারটিতে পৌঁছেছেন তখনও

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `<account>` | যে অ্যাকাউন্ট ব্যবহার করা হবে তার নাম বা নম্বর | |
| `[folder]` | ফোল্ডারটি, আপেক্ষিক বা সম্পূর্ণ পাথ | বর্তমান ফোল্ডার |
| `--remove` | ঠিক ওই ফোল্ডারের এন্ট্রি সরায় | বন্ধ |
| `--dry-run` | সমর্থিত: সংরক্ষণ না করে পরিবর্তনটি প্রিন্ট করে | বন্ধ |

## এটি কীভাবে কাজ করে

প্রতিটি ফোল্ডার `~/.config/ccam/config.toml`-এ একটি `[[folder]]` এন্ট্রি, যে ফাইলটি প্রতিটি চালু করা এমনিতেই পড়ে:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` ফোল্ডারটিকে সিম্বলিক লিংক মিলিয়ে একটি সম্পূর্ণ পাথে রূপান্তর করে এবং সংরক্ষণ করে। একটি ফোল্ডার আবার নির্ধারণ করলে তার অ্যাকাউন্ট বদলে যায়।
2. নাম ছাড়া একটি চালু করা `launcher.HereProfile` ডাকে। কোনো এন্ট্রি না থাকলে এটি সঙ্গে সঙ্গে ডিফল্ট অ্যাকাউন্ট নেয়।
3. অন্যথায় এটি ওয়ার্কিং ডিরেক্টরিকে প্রতিটি `path`-এর সঙ্গে স্ট্রিং হিসেবে তুলনা করে: ফোল্ডারটি নিজে বা তার নিচের একটি, কখনও এমন ফোল্ডার নয় যা শুধু একই অক্ষর দিয়ে শুরু হয় (`/code/company` `/code/companyx`-কে অন্তর্ভুক্ত করে না)। সবচেয়ে গভীর এন্ট্রি জেতে, তাই `~/code` `personal` হতে পারে আর `~/code/company` `work`।
4. ওয়ার্কিং ডিরেক্টরি যেমন দেওয়া আছে তেমন অবস্থায় কিছুর সঙ্গে না মিললেই শুধু ccam তার সিম্বলিক লিংক মিলিয়ে আবার তুলনা করে।
5. Windows ও macOS-এ তুলনাটি বড় ও ছোট হাতের অক্ষরের পার্থক্য উপেক্ষা করে, যেমন সাধারণত তাদের ফাইল সিস্টেমও করে।

`config.toml` ছাড়া আর কোনো ফাইল পড়া হয় না, তাই একটি ফোল্ডারের ভেতরে চালু করতে পরিমাপযোগ্য কোনো বাড়তি সময় লাগে না (এন্ট্রিসহ বা ছাড়া 4.0 ms, `hyperfine`, 300 বার চালিয়ে)।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- `ccam rename` ফোল্ডার এন্ট্রিগুলো নতুন নামে নিয়ে যায়, আর `ccam remove` সরানো অ্যাকাউন্টের এন্ট্রিগুলো মুছে দেয়।
- যে এন্ট্রির অ্যাকাউন্ট আর নেই, সেটি ওই ফোল্ডারে নাম ছাড়া চালু করা থামিয়ে দেয়: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`। `ccam doctor` এমন এন্ট্রি জানায়, এবং যে এন্ট্রিগুলোর ফোল্ডার আর নেই সেগুলোও।
- যে সাবফোল্ডার তার অ্যাকাউন্ট কোনো প্যারেন্ট ফোল্ডার থেকে পায়, সেখানে `ccam folder --remove` চালালে ccam প্যারেন্ট ফোল্ডারটির নাম জানায় এবং কিছুই সরায় না।
- `ccam next` রোটেশন বা ব্যবহার দেখে বেছে নেয় এবং ফোল্ডার দেখে না।

## সম্পর্কিত

- [Claude Code চালু করা](/features/launch)
- [শর্টকাট কমান্ড](/features/shortcuts)
- [অ্যাকাউন্ট পরিচালনা](/features/accounts)
- [কমান্ড রেফারেন্স](/docs/commands)
