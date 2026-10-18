---
title: Bayanin CLAUDE_CONFIG_DIR
description: Abin da CLAUDE_CONFIG_DIR ke yi a Claude Code, abin da ke cikin foldar, inda macOS ke ajiye bayanan shiga, da fayilolin da za a iya raba tsakanin asusai.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` yana gaya wa Claude Code wace folda zai yi amfani da ita maimakon `~/.claude`. Duk abin da ya shafi asusu yana nan a ciki: bayanan shiga, saituna, sabobin MCP, plugins da tarihin tattaunawa. Nuna shi zuwa foldoji daban-daban yana ba ka asusai daban-daban.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Abin da ke cikin foldar

An duba da Claude Code 2.1.302:

| Hanya | Abin da ke ciki | Ana iya rabawa |
| --- | --- | --- |
| `.credentials.json` | Token na shiga (Linux), token na MCP OAuth | A'a |
| `.claude.json` | Bayanan asusu, sabobin MCP na mai amfani, yanayin kowane project | A'a (sabobin MCP dabam) |
| `settings.json` | Izini, hooks, model, layin matsayi, plugins | Ee |
| `CLAUDE.md` | Umarni na gaba ɗaya (global) | Ee |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensions | Ee |
| `projects/` | Tattaunawa a matsayin fayilolin `.jsonl`, ƙwaƙwalwar kowane project | Ee |
| `history.jsonl` | Tarihin prompts | Ee |
| `sessions/`, `shell-snapshots/`, `cache/` da sauransu | Yanayin aiki (runtime state) | A'a |

## Bayanai biyu da ke ruɗar mutane

**Ba tare da `CLAUDE_CONFIG_DIR` ba, `.claude.json` ba ya cikin `~/.claude`.** A wannan yanayin Claude Code yana ajiye shi a `~/.claude.json`, kusa da foldar. Idan ka fara amfani da `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code yana neman `~/.claude/.claude.json` maimakon haka kuma ba ya samun asusunka. `ccam import` yana magance wannan ta hanyar haɗa su biyun da mahaɗi.

**A macOS bayanan shiga suna cikin Keychain, ba a cikin foldar ba.** Claude Code yana ajiye su ƙarƙashin sunan `Claude Code-credentials`, sannan gajeren hash na hanyar foldar idan an saita `CLAUDE_CONFIG_DIR`. Saboda haka foldoji daban-daban suna samun shigarwar Keychain daban-daban. Idan Keychain ya ƙi rubutun, Claude Code yana komawa ga `.credentials.json` a cikin foldar.

## Rabawa tsakanin foldoji

Symbolic links suna aiki: Claude Code yana karantawa da rubutu ta `settings.json` ko foldar `projects/` da aka haɗa da mahaɗi, kuma yana barin mahaɗin yadda yake. Sabobin MCP ne kawai suka bambanta, domin suna raba `.claude.json` da bayanan shiga; miƙa su da `--mcp-config` maimakon haka.

ccam yana gina aikinsa daidai a kan waɗannan ƙa'idodin. Duba [rabawa](/docs/sharing) don yadda yake haɗa foldoji, da [jagorar asusai da dama](/guides/multiple-claude-code-accounts) don cikakken hoto.
