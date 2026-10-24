---
title: CLAUDE_CONFIG_DIR समजून घ्या
description: Claude Code मध्ये CLAUDE_CONFIG_DIR काय करते, फोल्डरमध्ये काय साठवले जाते, macOS लॉगिन कुठे ठेवते, आणि खात्यांमध्ये कोणत्या फाइल्स शेअर करता येतात.
updated: 2026-10-23
order: 2
translation: machine
sourceHash: "2cf08478085f"
---

`CLAUDE_CONFIG_DIR` हे Claude Code ला `~/.claude` ऐवजी कोणता फोल्डर वापरायचा ते सांगते. खात्याशी संबंधित सर्व काही तिथे असते: लॉगिन, सेटिंग्ज, MCP सर्व्हर, plugins आणि संभाषण इतिहास. ते वेगवेगळ्या फोल्डर्सकडे निर्देशित केल्यास तुम्हाला वेगवेगळी खाती मिळतात.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR ला अधिकृत समर्थन आहे का?

हो. Claude Code त्याच्या [environment variables](https://code.claude.com/docs/en/env-vars) मध्ये configuration फोल्डर बदलण्याचा मार्ग म्हणून याचा उल्लेख करते, जो "अनेक खाती शेजारी शेजारी चालवण्यासाठी उपयोगी" आहे, आणि त्याचे [authentication मार्गदर्शक](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) अनेक खात्यांनी लॉगिन करण्यासाठी हेच वापरते. तुम्ही ते तुमच्या शेलमध्ये किंवा user settings मध्ये सेट करू शकता, पण प्रोजेक्टच्या settings मध्ये नाही.

## फोल्डरमध्ये काय असते?

Claude Code 2.1.309 सोबत तपासले:

| पाथ | मजकूर | शेअर करता येते |
| --- | --- | --- |
| `.credentials.json` | लॉगिन tokens (Linux, Windows), MCP OAuth tokens | नाही |
| `.claude.json` | खात्याचा तपशील, user MCP सर्व्हर, प्रत्येक प्रोजेक्टची स्थिती | नाही (MCP सर्व्हर वेगळे) |
| `settings.json` | परवानग्या, hooks, मॉडेल, स्टेटस लाइन, plugins | हो |
| `CLAUDE.md` | global सूचना | हो |
| `agents/`, `commands/`, `skills/`, `plugins/` | एक्स्टेंशन्स | हो |
| `projects/` | `.jsonl` फाइल्सच्या रूपातील संभाषणे, प्रत्येक प्रोजेक्टची memory | हो |
| `history.jsonl` | prompt इतिहास | हो |
| `sessions/`, `shell-snapshots/`, `cache/` आणि इतर | runtime स्थिती | नाही |

## लोक कुठे अडकतात?

**`CLAUDE_CONFIG_DIR` नसताना `.claude.json` हे `~/.claude` च्या आत नसते.** तेव्हा Claude Code ते फोल्डरच्या शेजारी `~/.claude.json` मध्ये ठेवते. तुम्ही `CLAUDE_CONFIG_DIR=~/.claude` वापरायला सुरुवात केली, तर Claude Code त्याऐवजी `~/.claude/.claude.json` शोधते आणि तुमचे खाते त्याला सापडत नाही. `ccam import` या दोघांना लिंक करून हे सांभाळते. `CLAUDE_CONFIG_DIR=~/.claude` सह आधी चालवल्यामुळे `~/.claude` मध्ये आधीच स्वतःची `.claude.json` असेल, तर ccam दोन्ही फाइल्स तशाच ठेवते आणि तसे सांगते, कारण कोणती योग्य आहे ते फक्त तुम्हालाच माहीत असते.

**macOS वर लॉगिन फोल्डरमध्ये नसून Keychain मध्ये असते.** Claude Code ते `Claude Code-credentials` या नावाखाली साठवते, आणि `CLAUDE_CONFIG_DIR` सेट असेल तेव्हा त्या नावानंतर `-` आणि फोल्डर पाथच्या SHA-256 चे पहिले 8 hex अंक येतात; ही entry `$USER` ने दिलेल्या account नावाखाली असते. त्यामुळे वेगळ्या फोल्डर्सना वेगळ्या Keychain entries मिळतात, आणि हलवलेला किंवा नाव बदललेला फोल्डर आपले लॉगिन गमावतो; म्हणूनच macOS वर `ccam rename` फोल्डरचे नाव तेच ठेवते. Keychain ने लिहिणे नाकारले, तर Claude Code फोल्डरमधील `.credentials.json` वापरते, जसे [त्याच्या दस्तऐवजात](https://code.claude.com/docs/en/authentication#credential-management) सांगितले आहे.

**प्रत्येक sign-in फोल्डरमध्ये नसते.** API key शिवाय केलेले Claude Console sign-in फोल्डरच्या बाहेर साठवले जाते, त्यामुळे वेगळे फोल्डर अशी दोन sign-ins वेगळी ठेवू शकत नाहीत. claude.ai लॉगिन आणि API keys प्रत्येक फोल्डरमध्ये वेगळे ठेवले जातात.

## Windows वर फोल्डर कुठे असतो?

डीफॉल्टनुसार `%USERPROFILE%\.claude`, आणि लॉगिन त्याच्या आत `.credentials.json` मध्ये, किंवा Claude Code ला ते वापरायला सेट केले असेल तर Windows Credential Manager मध्ये, macOS प्रमाणेच फोल्डर पाथवरून बनवलेल्या नावाखाली. प्रत्येक अतिरिक्त खात्यासाठी `CLAUDE_CONFIG_DIR` दुसऱ्या फोल्डरवर सेट करा, अगदी Linux आणि macOS प्रमाणे.

## फोल्डर्स फाइल्स शेअर करू शकतात का?

Symbolic links काम करतात: Claude Code लिंक केलेल्या `settings.json` किंवा `projects/` फोल्डरमधून वाचते आणि लिहिते, आणि लिंक कायम ठेवते. MCP सर्व्हर याला अपवाद आहेत, कारण ते लॉगिनसोबत `.claude.json` शेअर करतात; त्याऐवजी ते `--mcp-config` ने द्या.

ccam नेमक्या याच नियमांवर आधारित आहे. ते फोल्डर्स कसे लिंक करते यासाठी [शेअरिंग](/docs/sharing) पाहा, आणि संपूर्ण चित्रासाठी [अनेक खात्यांबद्दलचे मार्गदर्शक](/guides/multiple-claude-code-accounts) पाहा.
