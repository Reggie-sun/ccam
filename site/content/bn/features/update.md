---
title: "ccam আপডেট করা"
description: "ccam update একটি SHA-256 যাচাইয়ের পরে ccam-কে সর্বশেষ GitHub রিলিজ দিয়ে প্রতিস্থাপন করে; --check শুধু জানায় আর --version একটি নির্দিষ্ট রিলিজ ট্যাগ ইনস্টল করে।"
updated: 2026-10-19
order: 10
translation: machine
sourceHash: 83fc20191328
---

`ccam update` GitHub থেকে একটি ccam রিলিজ ডাউনলোড করে, রিলিজের `checksums.txt`-এর সঙ্গে সেটি যাচাই করে এবং একটি rename-এ চলমান বাইনারি প্রতিস্থাপন করে। ইনস্টল স্ক্রিপ্ট দিয়ে বা আর্কাইভ থেকে ccam ইনস্টল করে থাকলে এটি ব্যবহার করুন। ccam কখনও নিজে থেকে আপডেট খোঁজে না: শুধু এই কমান্ডটি GitHub-এর সঙ্গে যোগাযোগ করে, তাই বাকি প্রতিটি কমান্ড অফলাইনেও একইভাবে কাজ করে।

## ব্যবহার

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

নতুনতর কোনো রিলিজ না থাকলে ccam প্রিন্ট করে `ccam 0.4.0 is up to date.` এবং কিছু বদলায় না। `--version`-এ আপনার বর্তমানে চলমান সংস্করণটি দিলে এটি প্রিন্ট করে `ccam 0.4.0 is already installed.`

## অপশন

| অপশন | কী করে | ডিফল্ট |
| --- | --- | --- |
| `--check` | শুধু জানায় নতুনতর রিলিজ আছে কি না | বন্ধ |
| `--version <tag>` | এই রিলিজ ট্যাগ ইনস্টল করে, যেমন `v0.3.0`, বর্তমানটির চেয়ে নতুন বা পুরোনো | সর্বশেষ |
| `--dry-run` | মানা হয়: সর্বশেষ ট্যাগ খোঁজে (যদি না `--version` কোনোটির নাম দেয়), বাইনারিতে একই যাচাইগুলো চালায়, এবং ডাউনলোড না করে ডাউনলোড, চেকসাম ও প্রতিস্থাপনের পাথ প্রিন্ট করে | বন্ধ |

`--check` ও `--version` একসঙ্গে ব্যবহার করা যায় না। ট্যাগটি অবশ্যই `vMAJOR.MINOR.PATCH`-এর মতো হতে হবে, ঐচ্ছিক প্রি-রিলিজ বা বিল্ড প্রত্যয়সহ; সামনের `v` বাদ দেওয়া যায়। অন্য কিছু ব্যর্থ হয়, যেমন `--version takes a release tag such as v0.3.0, not "1.2"`।

এনভায়রনমেন্ট ভেরিয়েবল `CCAM_RELEASES_URL` মিরর বা পরীক্ষার জন্য `https://github.com/jhyunwoo/ccam/releases`-এর জায়গা নেয়। এটি অবশ্যই একইভাবে সাজানো হতে হবে: `<url>/latest` `<url>/tag/<tag>`-এ রিডাইরেক্ট করে, আর অ্যাসেটগুলো থাকে `<url>/download/<tag>/<asset>`-এ।

## এটি কীভাবে কাজ করে

কোডটি আছে `internal/update`-এ, যা শুধু `ccam update` ইমপোর্ট করে। চালু করার পথ কখনও এটি ডাকে না।

1. `cli.updateCmd` `os.Executable` দিয়ে চলমান বাইনারি খুঁজে বের করে এবং `filepath.EvalSymlinks` দিয়ে সিমলিংকগুলো অনুসরণ করে। সেই আসল পাথটিই প্রতিস্থাপিত হয়।
2. `--version` ছাড়া `Updater.Latest` রিডাইরেক্ট অনুসরণ না করে `<releases>/latest`-এ অনুরোধ পাঠায়, এবং `Location` হেডার থেকে ট্যাগটি পড়ে (`…/tag/v0.5.0`)। এটি GitHub API-এর বদলে রিলিজ পেজ ব্যবহার করে, তাই API-এর ঘণ্টায় ৬০টি অনুরোধের সীমা প্রযোজ্য হয় না। খোঁজার জন্য সময় ১০ সেকেন্ড।
3. `parseVersion` ট্যাগটিকে একটি কঠোর প্যাটার্নের সঙ্গে যাচাই করে। ট্যাগটি ডাউনলোড URL-এ এবং প্রিন্ট করা শেল কমান্ডে যায়, তাই `/`, `%`, স্পেস বা শেলের অক্ষরসহ কিছুই পেরোয় না। `compareVersions` এটিকে চলমান সংস্করণের সঙ্গে তুলনা করে; প্রি-রিলিজ তার রিলিজের আগে সাজানো হয়।
4. `managed()` ও `writable()` ঠিক করে ccam আদৌ বাইনারিটি প্রতিস্থাপন করতে পারবে কি না (নিচে দেখুন)।
5. `Updater.Fetch` একই সময়ে `ccam_<os>_<arch>.tar.gz` ও `checksums.txt` ডাউনলোড করে, দুটির জন্য মোট ২ মিনিট। প্রতিটি ডাউনলোড সর্বোচ্চ 256 MB।
6. অন্য কিছু ঘটার আগে আর্কাইভের SHA-256 অবশ্যই `checksums.txt`-এর সংশ্লিষ্ট লাইনের সঙ্গে মিলতে হবে। লাইন না থাকলে বা চেকসাম আলাদা হলে আপডেট `nothing was changed` দিয়ে থেমে যায়।
7. `extract()` tar পড়ে এবং `ccam` নামের প্রথম সাধারণ ফাইলটি নেয়। ডিরেক্টরি, লিংক ও অন্যান্য এন্ট্রি উপেক্ষা করা হয়, এবং বাইনারিও সর্বোচ্চ 256 MB।
8. `Updater.Install` `fsutil.AtomicWrite` দিয়ে নতুন বাইনারি লেখে: বাইনারির পাশে একটি অস্থায়ী ফাইল `.ccam.tmp*`, পুরোনো ফাইলের অনুমতি বিটসহ, ডিস্কে ফ্লাশ করে পুরোনোটির ওপর rename করা হয়।

ট্যাগ খোঁজা থেকে rename পর্যন্ত পুরো পাইপলাইন:

<!-- motion:update-pipeline -->

আর্কাইভ ও `checksums.txt` পাশাপাশি ডাউনলোড হয়, এবং চেকসাম না মেলা পর্যন্ত কিছুই লেখা হয় না। rename অ্যাটমিক, তাই ক্র্যাশ বা ভরা ডিস্কের পরে থাকে হয় পুরোনো বাইনারি নয়তো নতুনটি, কখনও অর্ধেক নয়। `~/.local/bin/ccwork`-এর মতো শর্টকাট কমান্ড ccam বাইনারির সিমলিংক, তাই কোনো পরিবর্তন ছাড়াই সেগুলো নতুন সংস্করণ চালায়।

### প্যাকেজ ম্যানেজারের ইনস্টল

`.deb` থেকে ইনস্টল করা ccam apt-এর অধীনে, এবং apt-এর অজান্তে সেটি প্রতিস্থাপন করলে প্যাকেজ ডেটাবেস ভুল থেকে যেত। ডাউনলোডের আগে ccam যাচাই করে বাইনারিটি কার:

<!-- motion:update-managed -->

1. `/usr/bin`, `/bin`, `/usr/sbin` বা `/sbin`-এ থাকা বাইনারি প্যাকেজ-পরিচালিত বলে গণ্য হয়।
2. অন্য জায়গায়, `dpkg` ইনস্টল থাকলে এবং `dpkg -S <path>` সফল হলে সেটিও প্যাকেজ-পরিচালিত বলে গণ্য হয়।
3. অন্যথায় ccam সেখানে লিখতে পারে কি না নিশ্চিত করতে বাইনারির পাশে একটি অস্থায়ী `.ccam-update-*` ফাইল তৈরি করে আবার সরায়।

প্যাকেজ-পরিচালিত বাইনারির জন্য ccam থেমে যায় এবং apt-এর ধাপগুলো প্রিন্ট করে:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Linux ছাড়া অন্য সিস্টেমে শুধু ইনস্টল স্ক্রিপ্টের লাইনটি প্রিন্ট হয়। ফোল্ডারে লেখা না গেলে ccam ইনস্টল স্ক্রিপ্টটি আবার চালানোর পরামর্শ দেয়, যা `~/.local/bin`-এ ccam ইনস্টল করে।

## নিরাপত্তা ও বিশেষ ক্ষেত্র

- ডেভেলপমেন্ট বিল্ড তার সংস্করণ `dev` হিসেবে জানায় এবং তুলনা করার মতো কোনো রিলিজ নেই, তাই `ccam update` ও `ccam update --check` এটি প্রত্যাখ্যান করে। `ccam update --version vX.Y.Z` এর ওপর একটি রিলিজ ইনস্টল করে।
- ব্যর্থ ডাউনলোড, আপনার প্ল্যাটফর্মের অ্যাসেট না থাকা, বা চেকসাম না মেলা ইনস্টল করা বাইনারি আগের মতোই রাখে।
- যে রিলিজ নেই সেটি `there is no release <tag>` দিয়ে ব্যর্থ হয়। GitHub HTTP 403 বা 429 দিয়ে প্রত্যাখ্যান করলে ccam জানায় এটি সম্ভবত রেট লিমিট এবং পরে আবার চেষ্টা করতে বলে।
- `--check` শুধু পড়ে। বাইনারিটি প্যাকেজ-পরিচালিত বা লেখার যোগ্য কি না তা এটি যাচাই করে না।
- পুরোনো ট্যাগসহ `--version` দিয়েই আগের সংস্করণে ফেরা যায়।

## সম্পর্কিত

- [ইনস্টল](/docs/install)
- [কমান্ড](/docs/commands)
- [সাধারণ প্রশ্ন](/docs/faq)
