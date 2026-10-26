---
title: extensions तपासणे
description: ccam ext प्रत्येक खाते कोणते skills, commands, plugins, MCP सर्व्हर आणि hooks लोड करते, प्रत्येक कुठून येते आणि काय बिघडले आहे ते दाखवते.
updated: 2026-10-26
order: 14
translation: machine
sourceHash: "7ca357a7f38a"
---

`ccam ext` या प्रश्नाचे उत्तर देते: "हे skill (किंवा MCP सर्व्हर) एका खात्यात चालते पण दुसऱ्यात का नाही?" ते प्रत्येक खात्यासाठी Claude Code लोड करत असलेली extensions आणि प्रत्येक कुठून येते ते दाखवते: एखादा share group, खात्याचा स्वतःचा फोल्डर, किंवा तुम्ही ज्या प्रकल्पात आहात तो. table खाली ते अशा गोष्टींची यादी देते ज्या एखादे extension तुमच्या अपेक्षेप्रमाणे लोड होण्यापासून रोखतात. ते काहीही बदलत नाही.

## वापर

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

प्रत्येक cell त्या खात्यातील त्या आयटमचा source सांगतो: share group (`default`), खात्याच्या स्वतःच्या फोल्डरमध्ये असल्यास `account`, खात्याकडे फक्त या फोल्डरसाठी असलेल्या MCP सर्व्हरसाठी `local`, किंवा खात्याकडे तो नसल्यास `-`. इन्स्टॉल केलेले पण चालू न केलेले plugin त्याच्या source नंतर `(off)` दाखवते. `!` समस्या असलेला आयटम दर्शवते, ज्याचे स्पष्टीकरण table खाली असते; कोणत्याही खात्याशी संबंधित नसलेल्या समस्या, जसे फक्त साध्या `claude` कडे असलेले सर्व्हर, शेवटी नावाशिवाय येतात. `ccam ext --json` तोच अहवाल छापते (`accounts`, प्रकल्प, `problems`), Claude Code network वरून ज्यांच्यापर्यंत पोहोचते अशा MCP सर्व्हरवर `remote` सह; एकही खाते नसल्यास ते `accounts: []` सह अहवाल छापते.

## ते काय वाचते

| प्रकार | कुठे | दाखवलेला source |
| --- | --- | --- |
| Skills | `skills/` मधील `SKILL.md` असलेला प्रत्येक फोल्डर | `skills/` ज्या group कडे link करतो तो, किंवा `account` |
| Commands, agents | `commands/` आणि `agents/` खालील `.md` फाइल्स, उप-फोल्डर `folder/name` म्हणून | group, किंवा `account` |
| Plugins | `plugins/installed_plugins.json`, खात्याच्या आणि प्रकल्पाच्या सेटिंग्जमधील `enabledPlugins` ने चालू किंवा बंद; फक्त एका प्रकल्पासाठी इन्स्टॉल केलेले plugin फक्त त्या प्रकल्पात दिसते | group, किंवा `account`; चालू नसल्यास `(off)` |
| MCP सर्व्हर | group चा `mcp.json` (MCP शेअर केलेले असताना), खात्याचा `.claude.json` (त्याचे स्वतःचे सर्व्हर, आणि फक्त या फोल्डरसाठी असलेले सर्व्हर), आणि त्याच्या plugins चे सर्व्हर, `plugin:<plugin>:<server>` म्हणून | group, `account` किंवा `local` |
| Hooks | `settings.json` मधील `hooks`, प्रत्येक event नुसार मोजलेले | group, किंवा `account` |
| हा प्रकल्प | सध्याच्या फोल्डरमधील `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` आणि `.mcp.json` | एकदाच दाखवले जाते; ज्या खात्याने फोल्डरवर विश्वास दाखवलेला नाही त्याला फोल्डरच्या `.mcp.json` सर्व्हरसाठी समस्या दाखवली जाते |

## ते कोणत्या समस्या सांगते

- **तुटलेल्या links**: `skills/`, `commands/` किंवा `agents/` मधील अशी link जिचे target आता नाही, किंवा अस्तित्वात नसलेल्या जागेकडे link करणारा शेअर केलेला फोल्डर (`skills`, `commands`, `agents`, `plugins`) (`ccam doctor --fix` ने दुरुस्त करा, जे नसलेला group फोल्डर पुन्हा तयार करते).
- **`SKILL.md` नसलेला skill फोल्डर**: Claude Code तो वगळते.
- **दोन ठिकाणी वेगवेगळ्या प्रकारे परिभाषित केलेला MCP सर्व्हर**: group चा `mcp.json` आणि खात्याचा `.claude.json`, किंवा त्यांपैकी एक आणि प्रकल्पाचा `.mcp.json`. दोन्ही ठिकाणी एकसारख्या व्याख्या (ज्या `mcp_strategy = "sync"` लिहिते) ही समस्या नाही. खात्याचा स्वतःचा सर्व्हर `ccam mcp import <account> <server>` ने group मध्ये हलवा.
- **वैध JSON नसलेली `installed_plugins.json`**, `enabledPlugins` ज्याचे नाव घेत नाही असे plugin, आणि ज्याचा फोल्डर नाही असे plugin.
- **फक्त एकच खाते लोड करते असे सर्व्हर**: फक्त या खात्याच्या `.claude.json` मधील सर्व्हर, जो त्याच्या group ची इतर खाती लोड करत नाहीत, आणि या फोल्डरसाठी खात्याच्या local scope मधील सर्व्हर (`-s user` दिल्याशिवाय `claude mcp add` हेच करते). प्रत्येकासोबत तो शेअर करणारी `ccam mcp import` कमांड येते. साध्या `claude` च्या `~/.claude.json` मधील, कोणतेही खाते लोड करत नाही असे सर्व्हर `ccam migrate --mcp` कमांडसह येतात.
- **मंजुरीची वाट पाहणारा प्रकल्प सर्व्हर**: खात्याने अजून फोल्डरवर विश्वास दाखवलेला नाही, म्हणून Claude Code एकदा विचारते, आणि तोपर्यंत प्रकल्पाचे `.mcp.json` सर्व्हर थांबून राहतात. त्यांना कुठे मंजुरी द्यायची ते समस्या सांगते: settings group मधील खात्यांच्या, किंवा या खात्याच्या, `settings.json` मध्ये `"enabledMcpjsonServers"`. git repository च्या main worktree ला दिलेला विश्वास त्याच्या इतर worktrees मध्येही लागू होतो.
- **साइन-इन्स आणि connectors**: दुसरे खाते साइन इन असताना या खात्यात साइन-इन लागणारा MCP सर्व्हर (साइन-इन फक्त एकाच खात्यासाठी चालते; `/mcp` ने साइन इन करा), आणि दुसऱ्या खात्याने वापरलेले पण या खात्याने न वापरलेले claude.ai connectors.

## सुरक्षितता आणि अपवादात्मक स्थिती

- `ccam ext` काहीही बदलत नाही आणि कोणतीही network request करत नाही. ते सर्व्हरची नावे आणि hooks ची संख्या वाचते, MCP definitions मधील `env` किंवा headers ची मूल्ये कधीच नाही. साइन-इन तपासणीसाठी ते प्रत्येक खात्याच्या लॉगिनमधील MCP साइन-इन्सची नावे वाचते (macOS वर Keychain मधून), tokens कधीच नाही, आणि फक्त अशा खात्यांसाठी जी दुसरे खातेही लोड करत असलेला सर्व्हर लोड करतात.
- disk वर जे आहे तेच ते दाखवते. managed settings ने अडवलेला सर्व्हर तरीही दिसतो; सेटिंग्जमध्ये बंद केलेले plugin `(off)` दाखवते.
- त्या प्रकल्पाचा `.claude` फोल्डर आणि `.mcp.json` समाविष्ट करण्यासाठी ते प्रकल्पाच्या फोल्डरमध्ये चालवा.

## संबंधित

- [गटांसह शेअरिंग](/features/share)
- [शेअर केलेले MCP सर्व्हर](/features/mcp)
- [तपासणी आणि दुरुस्ती](/features/doctor)
