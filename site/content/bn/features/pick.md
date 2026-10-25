---
title: "স্ক্রিপ্টের জন্য অ্যাকাউন্ট বেছে নেওয়া"
description: "ccam pick স্ক্রিপ্টের জন্য সবচেয়ে বেশি জায়গা থাকা অ্যাকাউন্ট প্রিন্ট করে, একসঙ্গে চালু worker-দের আলাদা অ্যাকাউন্ট দেয়, আর জায়গা না থাকলে 75 দিয়ে বেরোয়।"
updated: 2026-10-25
order: 12
translation: machine
sourceHash: 5c4f74db9714
---

`ccam pick` সবচেয়ে বেশি জায়গা থাকা লগইন করা অ্যাকাউন্টের নাম প্রিন্ট করে, যাতে কোনো স্ক্রিপ্ট সেটি দিয়ে Claude Code চালু করতে পারে। জায়গা থাকা অ্যাকাউন্ট যথেষ্ট থাকলে একই সময়ে চালু হওয়া প্রতিটি worker আলাদা অ্যাকাউন্ট পায়, তাই একগুচ্ছ `claude -p` কাজ একটি অ্যাকাউন্ট নিঃশেষ না করে আপনার অ্যাকাউন্টগুলোতে ছড়িয়ে পড়ে।

**Team ও Enterprise সিট।** স্বয়ংক্রিয় বাছাই (`ccam next`, `--fallback`, `ccam pick`, ড্যাশবোর্ডের `n`, এবং আপনি চালু করলে [`auto_switch`](/features/next#auto-switch)) একই Claude Team বা Enterprise প্রতিষ্ঠানের সিটগুলোর মধ্যেও বদল করে। Anthropic-এর সহায়তা কেন্দ্রে বলা হয়েছে, "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team প্ল্যানের প্রতিটি সিট একজন নির্দিষ্ট ব্যক্তির নামে বরাদ্দ করা হয়, এবং কেবল সেই ব্যক্তিই সেটি ব্যবহার করতে পারেন) এবং "Users can't route their work through another user's account to use that user's remaining usage" (ব্যবহারকারীরা অন্য কোনো ব্যবহারকারীর অ্যাকাউন্টের মধ্য দিয়ে নিজেদের কাজ চালিয়ে তাঁর ব্যবহারের সীমার বাকি অংশ কাজে লাগাতে পারেন না) ([Team প্ল্যান](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise প্ল্যানের](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) ক্ষেত্রেও একই কথা)। তাই সিটগুলোর মধ্যে স্বয়ংক্রিয় বাছাই ব্যবহার করলে Anthropic-এর শর্তাবলি লঙ্ঘন হতে পারে। সেখানে এটি ব্যবহার করবেন কি না, সে সিদ্ধান্ত আপনার: এর ওপর নির্ভর করার আগে আপনার প্রতিষ্ঠান ও Anthropic-এর কাছ থেকে জেনে নিন। ব্যক্তিগত প্ল্যানের (Pro, Max) কথা [প্রশ্নোত্তরের শর্তাবলি-বিষয়ক নোটে](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) বলা আছে।

## ব্যবহার

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

এক্সিট কোড: স্ট্যান্ডার্ড আউটপুটে নামসহ `0`, কোনো অ্যাকাউন্টে জায়গা না থাকলে বা কোনোটিতে লগইন না থাকলে `75` (`EX_TEMPFAIL`), এবং ভুল ফ্ল্যাগের মতো অন্য ত্রুটির জন্য `1`।

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--json` | অ্যাকাউন্ট, তার ৫ ঘণ্টা ও সাপ্তাহিক ব্যবহার, সংখ্যাগুলো কোথা থেকে এসেছে এবং কেন সেটি বেছে নেওয়া হয়েছে তা প্রিন্ট করে | বন্ধ |
| `--model <model>` | ওই মডেলের সাপ্তাহিক উইন্ডোও ধরে, যেমন `opus` | নেই |
| `--max-usage <1-100>` | এই ব্যবহারের নিচে থাকা অ্যাকাউন্টগুলোকে অগ্রাধিকার দেয় | `auto_max_usage`, নইলে 100 |
| `--strategy headroom\|expiring\|balance` | সবচেয়ে কম ৫ ঘণ্টার ব্যবহার, সবচেয়ে আগে আসা ৫ ঘণ্টার রিসেট, অথবা যে সাপ্তাহিক কোটা সবার আগে অব্যবহৃত থেকে যেত | `auto_strategy`, নইলে `headroom` |
| `--cross-org` | অন্য প্রতিষ্ঠানের অ্যাকাউন্টও বেছে নেয় | `cross_org`, নইলে বন্ধ |
| `--wait[=<duration>]` | কোনো অ্যাকাউন্টে জায়গা না থাকলে প্রথমটি রিসেট হওয়া পর্যন্ত অপেক্ষা করে, সর্বোচ্চ এতক্ষণ। মানটি `=`-এর পরে বসে: `--wait 30m` প্রত্যাখ্যাত হয়, এই ইঙ্গিতসহ: `--wait takes its value after '=': --wait=30m` | বন্ধ; একা দিলে `5h` |
| `--dry-run` | বাছাইটি নথিভুক্ত না করে বেছে নেয় | বন্ধ |

## এটি কীভাবে কাজ করে

1. ccam লগইন করা অ্যাকাউন্টগুলোর তালিকা করে (API key অ্যাকাউন্ট ও Console লগইনের কোনো সীমা নেই, তাই সেগুলো বাদ থাকে), শুধু সেই প্রতিষ্ঠানের, যার অ্যাকাউন্ট এই ফোল্ডারে চালু হতো ([প্রতিষ্ঠান](/features/next#auto) দেখুন); যে ফোল্ডার ওই অ্যাকাউন্টের মতো একই Claude অ্যাকাউন্টে লগইন করা, সেটিও বাদ যায়।
2. এটি একটি লক (`~/.local/state/ccam/picks.lock`) নেয়, যাতে একই সময়ে চলা বাছাইগুলো একটির পর একটি হয়।
3. এটি একসঙ্গে প্রতিটি অ্যাকাউন্টের ব্যবহার জিজ্ঞেস করে, ঠিক `ccam usage`-এর মতো, ৬০ সেকেন্ডের ক্যাশসহ। Claude Code যে উইন্ডোকে পূর্ণ বলে জানিয়েছে, সেটি রিসেট পর্যন্ত ১০০% ধরা হয় ([সীমার রেকর্ড](/features/usage#the-limit-record)), তাই সবে সীমায় পৌঁছানো একটি অ্যাকাউন্ট বাদই থাকে, এমনকি তার সংরক্ষিত সংখ্যা জায়গা দেখালেও।
4. গত ২ মিনিটে `ccam pick` বা কোনো সেশন সেখানে সরিয়ে নেওয়া [`auto_switch`](/features/next#auto-switch) যে অ্যাকাউন্টগুলো বেছে নিয়েছে (`~/.local/state/ccam/picks/<name>`; ঘড়ির সময়ের পরের তারিখের কোনো বাছাই গোনা হয় না), সেগুলো প্রথমে বাদ দেওয়া হয়, এবং বাকিগুলোর মধ্যে `ccam next --auto`-এর [বাছাইয়ের নিয়ম](/features/next#auto) বেছে নেয়: প্রথমে `auto_max_usage`-এর নিচে থাকাগুলো, কৌশল অনুযায়ী, এবং ১০০%-এ থাকা অ্যাকাউন্ট কখনও নয়। এতে নিয়ম অনুযায়ী প্রথমে যে অ্যাকাউন্ট আসত সেটি বাদ পড়লে `--json`-এর `reason` তা জানায়: `30% of its 5 hours used; work was picked a moment ago`।
5. জায়গা থাকা প্রতিটি অ্যাকাউন্ট সম্প্রতি বেছে নেওয়া হয়ে থাকলে নিয়মগুলো আবার সবগুলোর মধ্যে থেকে বেছে নেয়, তাই তিনটি অ্যাকাউন্টে চতুর্থ worker একটি অ্যাকাউন্ট ভাগ করে নেয়।
6. বাছাইটি নথিভুক্ত হয় ও লক ছেড়ে দেওয়া হয়, তারপর নামটি প্রিন্ট হয়।

২ মিনিটের উইন্ডোটি জরুরি, কারণ নতুন কাজের ব্যবহার সংখ্যায় দেখা দিতে কিছুটা সময় লাগে: এটি না থাকলে একসঙ্গে চালু হওয়া পাঁচটি worker সবাই একই অ্যাকাউন্টকে সবচেয়ে বেশি জায়গা থাকা অ্যাকাউন্ট হিসেবে দেখত।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- `ccam pick` শুধু ব্যবহার পড়ে এবং বাছাইয়ের সময় লেখে; `ccam usage` যতটুকু করে তার বাইরে এটি কখনও টোকেন নবায়ন করে না বা পড়ে না।
- বাছাই অ্যাকাউন্টটি সংরক্ষণ করে রাখে না। ২ মিনিট পরে দুটি স্ক্রিপ্ট তবুও একই অ্যাকাউন্টে চালু হতে পারে; ওই অ্যাকাউন্ট সীমায় পৌঁছালে সামলে নিতে `--fallback` ব্যবহার করুন।
- `--wait` দিয়ে অপেক্ষার সময় কোনো লক ধরে রাখা হয় না, তাই এর মধ্যে অন্য বাছাইগুলো চলতে থাকে।
- `ccam pick` [`ccam folder`](/features/folder) উপেক্ষা করে: এটি শুধু ব্যবহার দেখে বেছে নেয়।

## সম্পর্কিত

- [পরের অ্যাকাউন্টে যাওয়া](/features/next)
- [সীমায় আবার চালানো](/features/next#run-again-at-the-limit)
- [ব্যবহারের সীমা](/features/usage)
- [API key অ্যাকাউন্ট](/features/api)
