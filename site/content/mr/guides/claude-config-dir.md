---
title: CLAUDE_CONFIG_DIR समजून घ्या
description: Claude Code मध्ये CLAUDE_CONFIG_DIR काय करते, फोल्डरमध्ये काय साठवले जाते, macOS लॉगिन कुठे ठेवते, आणि खात्यांमध्ये कोणत्या फाइल्स शेअर करता येतात.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` हे Claude Code ला `~/.claude` ऐवजी कोणता फोल्डर वापरायचा ते सांगते. खात्याशी संबंधित सर्व काही तिथे असते: लॉगिन, सेटिंग्ज, MCP सर्व्हर, plugins आणि संभाषण इतिहास. ते वेगवेगळ्या फोल्डर्सकडे निर्देशित केल्यास तुम्हाला वेगवेगळी खाती मिळतात.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## फोल्डरमध्ये काय असते

Claude Code 2.1.302 सोबत तपासले:

| पाथ | मजकूर | शेअर करता येते |
| --- | --- | --- |
| `.credentials.json` | लॉगिन tokens (Linux), MCP OAuth tokens | नाही |
| `.claude.json` | खात्याचा तपशील, user MCP सर्व्हर, प्रत्येक प्रोजेक्टची स्थिती | नाही (MCP सर्व्हर वेगळे) |
| `settings.json` | परवानग्या, hooks, मॉडेल, स्टेटस लाइन, plugins | हो |
| `CLAUDE.md` | global सूचना | हो |
| `agents/`, `commands/`, `skills/`, `plugins/` | एक्स्टेंशन्स | हो |
| `projects/` | `.jsonl` फाइल्सच्या रूपातील संभाषणे, प्रत्येक प्रोजेक्टची memory | हो |
| `history.jsonl` | prompt इतिहास | हो |
| `sessions/`, `shell-snapshots/`, `cache/` आणि इतर | runtime स्थिती | नाही |

## लोकांना अडचणीत आणणारे दोन तपशील

**`CLAUDE_CONFIG_DIR` नसताना `.claude.json` हे `~/.claude` च्या आत नसते.** तेव्हा Claude Code ते फोल्डरच्या शेजारी `~/.claude.json` मध्ये ठेवते. तुम्ही `CLAUDE_CONFIG_DIR=~/.claude` वापरायला सुरुवात केली, तर Claude Code त्याऐवजी `~/.claude/.claude.json` शोधते आणि तुमचे खाते त्याला सापडत नाही. `ccam import` या दोघांना लिंक करून हे सांभाळते.

**macOS वर लॉगिन फोल्डरमध्ये नसून Keychain मध्ये असते.** Claude Code ते `Claude Code-credentials` या नावाखाली साठवते, आणि `CLAUDE_CONFIG_DIR` सेट असेल तेव्हा त्या नावानंतर फोल्डर पाथचा छोटा hash येतो. त्यामुळे वेगळ्या फोल्डर्सना वेगळ्या Keychain entries मिळतात. Keychain ने लिहिणे नाकारले, तर Claude Code फोल्डरमधील `.credentials.json` वापरते.

## फोल्डर्समध्ये शेअरिंग

Symbolic links काम करतात: Claude Code लिंक केलेल्या `settings.json` किंवा `projects/` फोल्डरमधून वाचते आणि लिहिते, आणि लिंक कायम ठेवते. MCP सर्व्हर याला अपवाद आहेत, कारण ते लॉगिनसोबत `.claude.json` शेअर करतात; त्याऐवजी ते `--mcp-config` ने द्या.

ccam नेमक्या याच नियमांवर आधारित आहे. ते फोल्डर्स कसे लिंक करते यासाठी [शेअरिंग](/docs/sharing) पाहा, आणि संपूर्ण चित्रासाठी [अनेक खात्यांबद्दलचे मार्गदर्शक](/guides/multiple-claude-code-accounts) पाहा.
