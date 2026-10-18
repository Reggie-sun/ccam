---
title: "ccam ইনস্টল করুন"
description: "Linux বা macOS-এ একটি কমান্ডে, .deb প্যাকেজ থেকে, অথবা সোর্স থেকে ccam ইনস্টল করুন। আগে Claude Code ইনস্টল থাকতে হবে।"
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Linux বা macOS-এ একটি টার্মিনালে এটি চালান:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

স্ক্রিপ্টটি আপনার সিস্টেমের জন্য রিলিজ ডাউনলোড করে, এর SHA-256 চেকসাম যাচাই করে এবং `~/.local/bin`-এ দুটি কমান্ড রাখে: `ccam` ও এর ছোট বিকল্প নাম `ccs`। ওই ফোল্ডার আপনার `PATH`-এ না থাকলে, `~/.bashrc` বা `~/.zshrc`-এ কোন লাইন যোগ করতে হবে তা স্ক্রিপ্টটি প্রিন্ট করে।

ccam আপনার হয়ে Claude Code চালু করে, তাই আগে [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ইনস্টল করুন এবং নিশ্চিত হন যে `claude --version` কাজ করে।

## সমর্থিত সিস্টেম

- x86-64 ও ARM64-এ Linux: Ubuntu, Debian, Fedora, Arch, Alpine এবং অন্যান্য ডিস্ট্রিবিউশন। বাইনারিটি স্ট্যাটিক এবং এর কোনো নির্ভরতা নেই।
- Apple silicon ও Intel-এ macOS।
- Windows সরাসরি সমর্থিত নয়। WSL-এর ভেতরে ccam অন্য যেকোনো Linux সিস্টেমের মতোই কাজ করে।

## ইনস্টলের অন্যান্য উপায়

**Debian ও Ubuntu প্যাকেজ।** [রিলিজ পাতা](https://github.com/jhyunwoo/ccam/releases) থেকে `ccam_<version>_linux_amd64.deb` ডাউনলোড করে ইনস্টল করুন:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

প্যাকেজটি `/usr/bin/ccam` এবং `/usr/bin/ccs` লিংক ইনস্টল করে।

**সোর্স থেকে।** Go 1.26 বা নতুন সংস্করণ দিয়ে:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
```

**একটি নির্দিষ্ট সংস্করণ।** `CCAM_VERSION` নির্ধারণ করুন:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## শেল কমপ্লিশন

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

কমপ্লিশন `ccam` ও `ccs` দুটিতেই কাজ করে, অ্যাকাউন্টের নামসহ।

## ইনস্টলেশন পরীক্ষা করুন

```sh
ccam version
ccam doctor
```

`ccam doctor` পরীক্ষা করে যে Claude Code পাওয়া যাচ্ছে এবং প্রতিটি অ্যাকাউন্ট ফোল্ডার, লিংক ও লগইন ঠিক আছে। এরপর [কুইক স্টার্ট](/docs/quick-start) দেখুন।

## আনইনস্টল

দুটি কমান্ড ও অ্যাকাউন্টের শর্টকাটগুলো সরিয়ে দিন, তারপর অ্যাকাউন্টগুলোর আর দরকার না থাকলে ডেটা ফোল্ডারগুলোও:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

`~/.local/share/ccam` মুছলে সেখানে সংরক্ষিত লগইন ও কথোপকথনের ইতিহাসও মুছে যায়। আপনি যে `~/.claude` ফোল্ডার [ইমপোর্ট করেছেন](/docs/commands) সেটি যেখানে আছে সেখানেই থাকে।
