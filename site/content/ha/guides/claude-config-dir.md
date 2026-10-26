---
title: Bayanin CLAUDE_CONFIG_DIR
description: Abin da CLAUDE_CONFIG_DIR ke yi a Claude Code, abin da ke cikin foldar, inda macOS ke ajiye bayanan shiga, da fayilolin da za a iya raba tsakanin asusai.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: dc2fe42e11e1
---

`CLAUDE_CONFIG_DIR` yana gaya wa Claude Code wace folda zai yi amfani da ita maimakon `~/.claude`. Duk abin da ya shafi asusu yana nan a ciki: bayanan shiga, saituna, sabobin MCP, plugins da tarihin tattaunawa. Nuna shi zuwa foldoji daban-daban yana ba ka asusai daban-daban.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Shin CLAUDE_CONFIG_DIR yana da goyon bayan hukuma?

Ee. Claude Code ya lissafa shi cikin [environment variables ɗinsa](https://code.claude.com/docs/en/env-vars) a matsayin hanyar sauya foldar saituna, "useful for running multiple accounts side by side", kuma [jagorar tantancewarsa](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) tana amfani da shi don shiga da asusai da dama. Za ka iya saita shi a shell ɗinka ko a saitunan mai amfani, amma ba a saitunan project ba.

## Me ke cikin foldar?

An duba da Claude Code 2.1.309:

| Hanya | Abin da ke ciki | Ana iya rabawa |
| --- | --- | --- |
| `.credentials.json` | Token na shiga (Linux, Windows), token na MCP OAuth | A'a |
| `.claude.json` | Bayanan asusu, sabobin MCP na mai amfani, yanayin kowane project | A'a (sabobin MCP dabam) |
| `settings.json` | Izini, hooks, model, layin matsayi, plugins | Ee |
| `CLAUDE.md` | Umarni na gaba ɗaya (global) | Ee |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensions | Ee |
| `projects/` | Tattaunawa a matsayin fayilolin `.jsonl`, ƙwaƙwalwar kowane project | Ee |
| `history.jsonl` | Tarihin prompts | Ee |
| `sessions/`, `shell-snapshots/`, `cache/` da sauransu | Yanayin aiki (runtime state) | A'a |

## Me ke ruɗar mutane?

**Ba tare da `CLAUDE_CONFIG_DIR` ba, `.claude.json` ba ya cikin `~/.claude`.** A wannan yanayin Claude Code yana ajiye shi a `~/.claude.json`, kusa da foldar. Idan ka fara amfani da `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code yana neman `~/.claude/.claude.json` maimakon haka kuma ba ya samun asusunka. `ccam import` yana magance wannan ta hanyar haɗa su biyun da mahaɗi. Idan `~/.claude` ya riga yana da `.claude.json` nasa, daga gudu na baya da `CLAUDE_CONFIG_DIR=~/.claude`, ccam yana barin fayilolin biyu kuma yana faɗin haka, tunda kai kaɗai ka san wanne ne daidai.

**A macOS bayanan shiga suna cikin Keychain, ba a cikin foldar ba.** Claude Code yana ajiye su ƙarƙashin sunan `Claude Code-credentials`. Idan an saita `CLAUDE_CONFIG_DIR`, sunan yana kuma ƙarewa da `-` da lambobin hex 8 na farko na SHA-256 na hanyar foldar. Sunan asusun abin Keychain ɗin shi ne `$USER`. Saboda haka foldoji daban-daban suna samun shigarwar Keychain daban-daban, kuma foldar da aka matsar ko aka sake mata suna tana rasa shigarta; shi ya sa `ccam rename` yake riƙe sunan foldar a macOS. Idan Keychain ya ƙi rubutun, Claude Code yana komawa ga `.credentials.json` a cikin foldar, kamar yadda [takardunsa](https://code.claude.com/docs/en/authentication#credential-management) suka bayyana.

**Ba kowace shiga ce ke cikin foldar ba.** Shiga ta Claude Console ba tare da API key ba ana ajiye ta a wajen foldar, don haka foldoji daban-daban ba za su iya ware irin waɗannan shiga biyu ba. Shigar claude.ai da API keys ana ajiye su a kowace folda.

## Ina foldar take a Windows?

A tsoho (default) ita ce `%USERPROFILE%\.claude`. Bayanan shiga suna a `.credentials.json` a cikinta, ko a Windows Credential Manager idan an saita Claude Code ya yi amfani da shi, ƙarƙashin sunan da aka samo daga hanyar foldar kamar a macOS. Saita `CLAUDE_CONFIG_DIR` zuwa wata folda ga kowane ƙarin asusu, kamar yadda ake yi a Linux da macOS.

## Shin foldoji za su iya raba fayiloli?

Symbolic links suna aiki: Claude Code yana karantawa da rubutu ta `settings.json` ko foldar `projects/` da aka haɗa da mahaɗi, kuma yana barin mahaɗin yadda yake. Sabobin MCP ne kawai suka bambanta, domin suna raba `.claude.json` da bayanan shiga; miƙa su da `--mcp-config` maimakon haka.

ccam yana gina aikinsa daidai a kan waɗannan ƙa'idodin. Duba [rabawa](/docs/sharing) don yadda yake haɗa foldoji, da [jagorar asusai da dama](/guides/multiple-claude-code-accounts) don cikakken hoto.
