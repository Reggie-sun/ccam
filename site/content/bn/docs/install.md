---
title: "ccam ইনস্টল করুন"
metaTitle: "Linux, macOS ও Windows-এ ccam ইনস্টল করুন"
description: "Linux, macOS বা Windows-এ একটি কমান্ডে, .deb প্যাকেজ বা সোর্স থেকে ccam ইনস্টল করুন, আর ccam update দিয়ে আপডেট করুন। আগে Claude Code ইনস্টল থাকতে হবে।"
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 07fab3165030
---

Linux বা macOS-এ একটি টার্মিনালে এটি চালান:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

স্ক্রিপ্টটি আপনার সিস্টেমের জন্য রিলিজ ডাউনলোড করে, তার SHA-256 চেকসাম যাচাই করে এবং `ccam`-কে `~/.local/bin`-এ রাখে; ccam পরে আপনার বেছে নেওয়া শর্টকাট কমান্ড, যেমন `ccwork`, এই ফোল্ডারেই তৈরি করে। ফোল্ডারটি আপনার `PATH`-এ না থাকলে স্ক্রিপ্টটি `~/.bashrc` বা `~/.zshrc`-এ যোগ করার লাইনটি প্রিন্ট করে।

Windows-এ PowerShell-এ এটি চালান:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

এটি `ccam_windows_<arch>.zip` ডাউনলোড করে, তার SHA-256 চেকসাম যাচাই করে, `ccam.exe`-কে `%USERPROFILE%\.local\bin`-এ রাখে, যে ফোল্ডারে Claude Code-এর ইনস্টলার `claude.exe` রাখে, এবং ফোল্ডারটি আপনার ইউজার `PATH`-এ না থাকলে সেখানে যোগ করে। এরপর একটি নতুন টার্মিনাল খুলুন।

ccam আপনার হয়ে Claude Code চালু করে, তাই আগে [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ইনস্টল করুন এবং নিশ্চিত হন যে `claude --version` কাজ করে।

## সমর্থিত সিস্টেম

- x86-64 ও ARM64-এ Linux: Ubuntu, Debian, Fedora, Arch, Alpine এবং অন্যান্য ডিস্ট্রিবিউশন। বাইনারিটি স্ট্যাটিক এবং এর কোনো নির্ভরতা নেই।
- Apple silicon ও Intel-এ macOS।
- x86-64 ও ARM64-এ Windows 10 ও 11, PowerShell, Windows Terminal বা Git Bash-এ। শেয়ার করা ফোল্ডার প্রত্যেক ব্যবহারকারীর জন্য কাজ করে। শেয়ার করা ফাইলের (`settings.json`, `CLAUDE.md`, `history.jsonl`) জন্য Developer Mode (Settings > System > For developers) লাগে, কারণ Windows শুধু এই মোডে, বা একজন অ্যাডমিনিস্ট্রেটরকে, ফাইলের লিংক তৈরি করতে দেয়; [Windows-এ শেয়ার করা](/features/share#on-windows) দেখুন।
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

এটি একটি skill, `/ccam:usage`, `/ccam:ext` ও `/ccam:switch` কমান্ড, এবং দুটি হুক যোগ করে: কোনো সেশন ব্যবহারের সীমায় থেমে গেলে আপনি টার্মিনালে একটি নোটিফিকেশন পান, আর `/ccam:switch [account]` কথোপকথনটি অন্য একটি অ্যাকাউন্টের হাতে তুলে দেয় ([বিস্তারিত](/features/next#inside-claude-code))। হুকগুলো `ccam` চালায়, তাই Claude Code যে `PATH` দেখে, তাতে এটি থাকতে হবে। skill-টি কখনও লগইন করা বা ফাইল সরানোর কমান্ড চালায় না; বরং সেগুলো প্রস্তাব করে।

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

ccam কখনো নিজে থেকে আপডেট খোঁজে না; আপনি এই কমান্ড চালালেই কেবল নতুন রিলিজ খোঁজে। কোনো নির্দিষ্ট রিলিজ, এমনকি পুরোনো রিলিজও, ইনস্টল করতে `--version` দিন। `ccam --dry-run update` কিছু না বদলে দেখায় কী ডাউনলোড হবে আর কী বদলানো হবে।

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

ইনস্টল স্ক্রিপ্টটি আবার চালালেও `~/.local/bin`-এর কপিটি আপডেট হয়। ccam যদি `go install` দিয়ে বিল্ড করে থাকেন, তাহলে একই কমান্ড আবার চালান।

**Windows।** `ccam update` Windows-এর জন্য `.zip` ডাউনলোড করে। Windows-এ চলমান প্রোগ্রাম বদলানো যায় না, তাই পুরোনো বাইনারির নাম বদলে `ccam.exe.old` করা হয়, যা পরের আপডেট মুছে দেয়।

**Debian ও Ubuntu প্যাকেজ।** apt দিয়ে ইনস্টল করা `/usr/bin/ccam` ফাইলটি `ccam update` বদলায় না, বরং এর বদলে এই ধাপগুলো প্রিন্ট করে। [রিলিজ পাতা](https://github.com/jhyunwoo/ccam/releases) থেকে নতুন `ccam_<version>_linux_amd64.deb` ডাউনলোড করে পুরোনোটির ওপর ইনস্টল করুন:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## আনইনস্টল

ccam ইনস্টল থাকতে থাকতেই আগে অ্যাকাউন্টের শর্টকাটগুলো সরান। `ccam list` সেগুলো Shortcut কলামে দেখায়; প্রতিটি `ccam shortcut <account> --remove` দিয়ে সরান, অথবা ওই লিংকগুলো `~/.local/bin` থেকে নিজে মুছে দিন। তারপর কমান্ডটি সরান, আর অ্যাকাউন্টগুলোর আর দরকার না থাকলে ডেটা ফোল্ডারগুলোও:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Windows-এ, PowerShell-এ:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

`~/.local/share/ccam` মুছলে সেখানে সংরক্ষিত লগইন ও কথোপকথনের ইতিহাসও মুছে যায়। আপনি যে `~/.claude` ফোল্ডার [ইমপোর্ট করেছেন](/docs/commands) সেটি যেখানে আছে সেখানেই থাকে।
