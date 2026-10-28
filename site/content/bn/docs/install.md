---
title: "ccam ইনস্টল করুন"
metaTitle: "Linux, macOS ও Windows-এ ccam ইনস্টল করুন"
description: "Linux, macOS বা Windows-এ একটি কমান্ডে, .deb প্যাকেজ বা go install দিয়ে ccam ইনস্টল করুন, আর ccam update দিয়ে আপডেট করুন। আগে Claude Code ইনস্টল থাকতে হবে।"
updated: 2026-10-28
order: 1
translation: machine
sourceHash: 193df30c9d9a
---

ccam আপনার হয়ে Claude Code চালু করে, তাই আগে [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ইনস্টল করুন এবং নিশ্চিত হন যে `claude --version` কাজ করে।

Linux বা macOS-এ একটি টার্মিনালে এটি চালান:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

স্ক্রিপ্টটি আপনার সিস্টেমের জন্য রিলিজ ডাউনলোড করে, তার SHA-256 চেকসাম যাচাই করে এবং `ccam`-কে `~/.local/bin`-এ রাখে; ccam পরে আপনার বেছে নেওয়া শর্টকাট কমান্ড, যেমন `ccwork`, এই ফোল্ডারেই তৈরি করে। ফোল্ডারটি আপনার `PATH`-এ না থাকলে স্ক্রিপ্টটি `~/.bashrc` বা `~/.zshrc`-এ যোগ করার লাইনটি প্রিন্ট করে।

Windows-এ PowerShell-এ এটি চালান:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

এটি `ccam_windows_<arch>.zip` ডাউনলোড করে, তার SHA-256 চেকসাম যাচাই করে এবং `ccam.exe`-কে `%USERPROFILE%\.local\bin`-এ রাখে, যে ফোল্ডারে Claude Code-এর ইনস্টলার `claude.exe` রাখে। ফোল্ডারটি আপনার ইউজার `PATH`-এ এখনও না থাকলে স্ক্রিপ্টটি সেটি যোগ করে, আর `%USERPROFILE%\bin`-এর মতো এন্ট্রিগুলো যেমন আছে তেমনই রাখে। এরপর একটি নতুন টার্মিনাল খুলুন।

## সমর্থিত সিস্টেম

- x86-64 ও ARM64-এ Linux: Ubuntu, Debian, Fedora, Arch, Alpine এবং অন্যান্য ডিস্ট্রিবিউশন। বাইনারিটি স্ট্যাটিক এবং এর কোনো নির্ভরতা নেই।
- Apple silicon ও Intel-এ macOS।
- x86-64 ও ARM64-এ Windows 10 ও 11, PowerShell, Windows Terminal বা Git Bash-এ:
  - শেয়ার করা ফোল্ডার Developer Mode ছাড়াই কাজ করে।
  - শেয়ার করা ফাইলের (`settings.json`, `CLAUDE.md`, `history.jsonl`) জন্য Developer Mode (Settings > System > For developers) লাগে, কারণ Windows শুধু এই মোডে, বা একজন অ্যাডমিনিস্ট্রেটরকে, ফাইলের লিংক তৈরি করতে দেয়; [Windows-এ শেয়ার করা](/features/share#on-windows) দেখুন।
  - Developer Mode ছাড়া প্রথম `ccam add` শুধু MCP সার্ভার ও এক্সটেনশন নিয়ে `default` গ্রুপ তৈরি করে, এবং Developer Mode চালু হলে বাকিগুলো কীভাবে যোগ করবেন তা জানায়।
- WSL-এর ভেতরে ccam অন্য যেকোনো Linux সিস্টেমের মতোই কাজ করে।

## ইনস্টলের অন্যান্য উপায়

**Debian ও Ubuntu প্যাকেজ।** [রিলিজ পাতা](https://github.com/jhyunwoo/ccam/releases) থেকে `ccam_<version>_linux_amd64.deb` ডাউনলোড করে ইনস্টল করুন:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

প্যাকেজটি `/usr/bin/ccam` ইনস্টল করে।

**সোর্স থেকে।** Go 1.26 বা নতুন সংস্করণ দিয়ে:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

কোনো রিলিজ ট্যাগ (`@latest` বা `@vX.Y.Z`) ইনস্টল করলে ccam সেই সংস্করণ জানায়, এবং `ccam update` ও ড্যাশবোর্ডের আপডেট প্রস্তাব কাজ করে। `@main`-এর, বা কোনো ট্যাগের পরে কমিট থাকা লোকাল চেকআউটের বিল্ড `dev` জানায়, এবং আপনি `--version` না দিলে `ccam update` সেটি প্রত্যাখ্যান করে।

**একটি নির্দিষ্ট সংস্করণ।** `CCAM_VERSION` নির্ধারণ করুন:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## শেল কমপ্লিশন

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

কমপ্লিশন `ccam`-এর সব কমান্ড ও অপশনে কাজ করে, অ্যাকাউন্টের নামসহ।

## Claude Code প্লাগইন

ccam প্লাগইন Claude Code-কে নিজেই ccam ব্যবহার করতে শেখায়: কোন অ্যাকাউন্টে জায়গা আছে তা দেখা, `ccam pick` দিয়ে আলাদা আলাদা অ্যাকাউন্টে সমান্তরাল `claude -p` worker চালু করা, এবং `ccam ext` দিয়ে কোনো skill বা MCP সার্ভার কেন পাওয়া যাচ্ছে না তা ব্যাখ্যা করা। Claude Code-এর ভেতরে এটি ইনস্টল করুন:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

এটি যোগ করে:

- একটি skill, যা লগইন করা বা ফাইল সরানোর কমান্ড প্রস্তাব করে কিন্তু কখনও সেগুলো চালায় না;
- `/ccam:usage`, `/ccam:ext` ও `/ccam:switch` কমান্ড;
- দুটি হুক ([বিস্তারিত](/features/next#inside-claude-code))। একটি হুক কোনো সেশন ব্যবহারের সীমায় থেমে গেলে টার্মিনালে একটি নোটিফিকেশন দেখায়; Linux ও macOS-এ [`auto_switch`](/features/next#auto-switch) চালু থাকলে ccam তারপর সেশনটি জায়গা থাকা অন্য একটি অ্যাকাউন্টে সরিয়ে নেয়। অন্যটি `/ccam:switch [account]`-কে কথোপকথনটি অন্য একটি অ্যাকাউন্টের হাতে তুলে দিতে দেয়।

হুকগুলো `ccam` চালায়, তাই Claude Code যে `PATH` দেখে, তাতে এটি থাকতে হবে।

## ইনস্টলেশন পরীক্ষা করুন

```sh
ccam version
ccam doctor
```

`ccam doctor` পরীক্ষা করে যে Claude Code পাওয়া যাচ্ছে এবং প্রতিটি অ্যাকাউন্ট ফোল্ডার, লিংক ও লগইন ঠিক আছে। এরপর [কুইক স্টার্ট](/docs/quick-start) দেখুন।

## আপডেট

```sh
ccam update --check
ccam update
```

`ccam update --check` শুধু জানায় নতুন কোনো রিলিজ আছে কি না। `ccam update` আপনার সিস্টেমের জন্য রিলিজ ডাউনলোড করে, রিলিজের `checksums.txt` দিয়ে তার SHA-256 চেকসাম যাচাই করে, তারপর এক ধাপে `ccam` ফাইলটি বদলে দেয়। চেকসাম না মিললে কিছুই বদলায় না। `ccwork`-এর মতো শর্টকাটগুলো ccam-এর লিংক, তাই সেগুলো সঙ্গে সঙ্গে নতুন সংস্করণ চালায়।

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ড্যাশবোর্ডও (টার্মিনালে আর্গুমেন্ট ছাড়া `ccam`) প্রতিবার খোলার সময় নতুন রিলিজ খোঁজে এবং থাকলে সেটি প্রস্তাব করে: `U` চাপুন এবং নিশ্চিত করুন। অন্য কোনো কমান্ড আপডেট খোঁজে না, এবং ccam জিজ্ঞেস না করে কখনো আপডেট হয় না। কোনো নির্দিষ্ট রিলিজ, এমনকি পুরোনো রিলিজও, ইনস্টল করতে `--version` দিন। `ccam --dry-run update` কিছু না বদলে দেখায় কী ডাউনলোড হবে আর কী বদলানো হবে।

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

ইনস্টল স্ক্রিপ্টটি আবার চালালেও `~/.local/bin`-এর কপিটি আপডেট হয়। ccam যদি কোনো রিলিজ ট্যাগ থেকে `go install` দিয়ে বিল্ড করে থাকেন, তাহলে `ccam update` বা একই `go install` কমান্ড আবার চালাতে পারেন।

**Windows।** `ccam update` Windows-এর জন্য `.zip` ডাউনলোড করে। Windows-এ চলমান প্রোগ্রাম বদলানো যায় না, তাই পুরোনো বাইনারির নাম বদলে `ccam.exe.old` করা হয়, যা পরের আপডেট মুছে দেয়; আগের কোনো `ccam.exe.old` থেকে কোনো সেশন তখনও চললে নামটি হয় `ccam.exe.old-<time>`। হার্ড লিংক শর্টকাটগুলো ঠিক পরেই নবায়ন করা হয়। `install.ps1` আবার চালালেও একইভাবে কাজ হয়, এমনকি পুরোনো বাইনারি থেকে কোনো Claude Code সেশন তখনও চললেও, এবং হার্ড লিংক শর্টকাটগুলো নবায়ন করার সময় এটি `OK 2 shortcuts now run the new ccam` প্রিন্ট করে।

**Debian ও Ubuntu প্যাকেজ।** apt দিয়ে ইনস্টল করা `/usr/bin/ccam` ফাইলটি `ccam update` বদলায় না, বরং এর বদলে এই ধাপগুলো প্রিন্ট করে। [রিলিজ পাতা](https://github.com/jhyunwoo/ccam/releases) থেকে নতুন `ccam_<version>_linux_amd64.deb` ডাউনলোড করে পুরোনোটির ওপর ইনস্টল করুন:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## আনইনস্টল

`ccam uninstall` ccam এবং এর তৈরি করা সবকিছু সরিয়ে দেয়। এটি আগে দেখায় কী কী যাবে, তারপর জিজ্ঞেস করে:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- প্রতিটি অ্যাকাউন্ট `ccam remove`-এর মতো করেই সরানো হয়। ccam-এর ফোল্ডারে থাকা অ্যাকাউন্ট আগে লগআউট হয়, তারপর তার ফোল্ডারসহ মুছে যায়, সঙ্গে তার কথোপকথনের ইতিহাস ও সেটিংসও।
- আপনি যে `~/.claude` ফোল্ডার [ইমপোর্ট করেছেন](/docs/commands) সেটি সাধারণ `claude`-এর জন্য যেখানে আছে সেখানেই থাকে। সেটি কোনো গ্রুপে থাকলে, যা যা শেয়ার করত তার একটি কপি ফিরে পায়। এর `settings.json`-এ থাকা ccam-এর স্ট্যাটাস লাইনের জায়গায় ফিরে আসে সেই স্ট্যাটাস লাইন, যেটিকে ccam মুড়ে রেখেছিল, যদি থেকে থাকে। এতে ccam প্লাগইন থাকলে, Claude Code-এ চালানোর জন্য ccam `/plugin uninstall` কমান্ডটি প্রিন্ট করে।
- শর্টকাট কমান্ডগুলো, `~/.config/ccam`, `~/.local/share/ccam` (ccam-এর ট্র্যাশসহ) ও `~/.local/state/ccam` মুছে ফেলা হয়, তারপর `ccam` কমান্ডটি নিজেও।

কোনো অ্যাকাউন্ট চলতে থাকলে ccam রাজি হয় না; আগে সেটির Claude Code সেশনগুলো বন্ধ করুন। apt যে `/usr/bin/ccam` ইনস্টল করেছে সেটি থেকে যায়; পরে `sudo apt remove ccam` দিয়ে সেটি সরান। Windows-এ চলমান প্রোগ্রাম নিজেকে মুছতে পারে না, তাই `ccam uninstall` শেষে `ccam.exe`-এর জন্য `Remove-Item` কমান্ডটি প্রিন্ট করে।
