---
title: "কমান্ড রেফারেন্স"
description: "ccam add ও ccam <name> থেকে ccam usage, ccam next ও ccam doctor পর্যন্ত ccam-এর প্রতিটি কমান্ড ও অপশন।"
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

কোনো কমান্ড আসলে না চালিয়ে সেটি কী করবে তা দেখতে কমান্ডের আগে `--dry-run` দিন, যেমন `ccam --dry-run work`।

<!-- commands -->

## Claude Code চালু করা

`ccam <name|number> [arguments]` ওই অ্যাকাউন্ট হিসেবে Claude Code চালু করে। নামের পরের আর্গুমেন্টগুলো ccam পড়ে না; সেগুলো একই ক্রমে `claude`-এ পাঠিয়ে দেয়, এবং স্ট্যান্ডার্ড ইনপুট, এক্সিট কোড ও সিগন্যাল ঠিক `claude`-এর মতোই আচরণ করে। নামের পরে `--` দেওয়া ঐচ্ছিক।

অ্যাকাউন্টের নাম ছাড়া কিন্তু আর্গুমেন্টসহ `ccam`, যেমন `ccam -p "summarise"`, ডিফল্ট অ্যাকাউন্ট ব্যবহার করে। এটি `ccam use <name>` দিয়ে নির্ধারণ করুন।

`claude1`, `claude2` ইত্যাদি শর্টকাট হলো `~/.local/bin`-এ থাকা ccam-এর লিংক। `ccam shims install --names` প্রতিটি অ্যাকাউন্টের জন্য `claude-<name>`-ও তৈরি করে।

## অ্যাকাউন্ট যোগ করা ও সরানো

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

ফোল্ডার না দিয়ে `ccam import` চালালে এটি `~/.claude` ও `~/.claude-*` খোঁজে এবং প্রতিটির বিষয়ে আপনাকে জিজ্ঞেস করে। ইমপোর্ট করা ফোল্ডার যেখানে আছে সেখানেই নিবন্ধিত হয়; ccam কখনও কোনো লগইন সরায় বা কপি করে না।

## শেয়ার করা

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

প্রতিটি আইটেম `on` (`default` গ্রুপ), `off`, অথবা একটি গ্রুপের নাম নেয়। দেখুন [শেয়ার করা](/docs/sharing)।

## MCP সার্ভার

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## ব্যবহার ও অ্যাকাউন্ট বদল

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## রক্ষণাবেক্ষণ

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## ফাইল

| পাথ | বিষয়বস্তু |
| --- | --- |
| `~/.config/ccam/config.toml` | অ্যাকাউন্ট এবং সেগুলো কী শেয়ার করে |
| `~/.local/share/ccam/profiles/<name>/` | প্রতিটি অ্যাকাউন্টের `CLAUDE_CONFIG_DIR` |
| `~/.local/share/ccam/groups/<group>/` | শেয়ার করা সেটিংস, MCP সার্ভার, এক্সটেনশন ও ইতিহাস |
| `~/.local/state/ccam/` | ব্যবহারের ক্যাশ ও সর্বশেষ ব্যবহারের সময় |
| `~/.local/bin/claude1`, `claude2`… | ccam-এর লিংক |

`XDG_CONFIG_HOME`, `XDG_DATA_HOME` ও `XDG_STATE_HOME` ভেরিয়েবল এই ফোল্ডারগুলোর অবস্থান বদলায়।
