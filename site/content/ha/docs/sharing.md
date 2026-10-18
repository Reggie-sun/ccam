---
title: Raba sabobin MCP, saituna da tarihi
description: Yadda ccam ke raba sabobin MCP, settings.json, CLAUDE.md, agents, skills da tarihin tattaunawa tsakanin asusun Claude Code, yayin da kowane shiga yake daban.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

Kowane asusu yana riƙe bayanan shigarsa. Duk sauran abubuwa za a iya raba su ɗaya bayan ɗaya: ccam yana mayar da fayil ko foldar asusun mahaɗi zuwa foldar rukuni da aka raba, don haka canjin da aka yi a asusu ɗaya yana bayyana a duk asusan wannan rukunin.

## Abin da za a iya rabawa

| Abu | Zaɓi | Abin da ake rabawa |
| --- | --- | --- |
| Sabobin MCP | `--mcp` | `mcp.json` na rukunin, wanda ake miƙa wa Claude Code lokacin farawa |
| Saituna | `--settings` | `settings.json`: izini, hooks, model, layin matsayi |
| Ƙwaƙwalwa ta gaba ɗaya | `--memory` | `CLAUDE.md` |
| Extensions | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Tarihin tattaunawa | `--history` | `projects/` da `history.jsonl` |

Sababbin asusai suna raba komai banda tarihin tattaunawa.

Ba a taɓa raba fayilolin shiga `.credentials.json` da `.claude.json` ba. ccam yana haɗa hanyoyi ne kawai daga ƙayyadadden jeri, kuma waɗannan biyun ba sa cikinsa.

## Rukunoni

`on` yana nufin rukunin `default`. Za ka iya kuma ba da suna, don wasu asusai kawai su raba wani abu:

```sh
ccam share work --history team
ccam share side --history team
```

Yanzu `work` da `side` suna raba tattaunawa da juna, amma ba da sauran asusunka ba.

## Kunna rabawa

Idan asusu ya fara raba wani abu, ccam:

1. yana buga shirin kowane fayil da zai motsa, kwafa ko haɗa da mahaɗi, kuma yana tambaya kafin ya fara;
2. yana motsa fayilolin asusun zuwa foldar rukunin idan rukunin har yanzu babu komai a ciki, ko yana ƙara fayilolin da rukunin bai da su tukuna;
3. idan fayil ɗaya ya bambanta, yana tambaya ko a riƙe sigar rukunin, a yi amfani da taka, ko a haɗa JSON ɗin;
4. yana ajiye ainihin naka a `.ccam-backup/<time>/` a cikin foldar asusun;
5. yana maye gurbin ainihin da mahaɗi zuwa rukunin.

Ana sanya wa fayilolin tattaunawa suna bisa session ID ɗinsu, don haka tarihi daga asusai da dama yana haɗuwa ba tare da rikici ba.

## Kashe rabawa

```sh
ccam share work --history off
```

Ana cire mahaɗin, kuma asusun yana samun kwafinsa na abin da ke cikin rukunin a yanzu. Ƙara `--fresh` don farawa babu komai a maimakon haka. Rukunin yana riƙe abin da ke cikinsa don sauran asusai.

## Sabobin MCP

Claude Code yana ajiye sabobin MCP a `.claude.json`, tare da bayanan shiga, don haka ba za a iya raba wannan fayil ba. ccam yana ajiye sabobin da aka raba a `mcp.json` na rukunin kuma yana fara Claude Code da `--mcp-config=<group>/mcp.json`. Zaɓuɓɓukan `--mcp-config` naka ma suna ci gaba da aiki, a saman waɗannan.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Sabar da ka ƙara a cikin zama (session) da `claude mcp add` ana ajiye ta ga wannan asusun kawai. A gaba da ka fara wannan asusun, ccam yana faɗin sunan sabar da umarnin da zai motsa ta zuwa rukunin:

```sh
ccam mcp import work gh
```

Sabobin MCP da ke shiga da OAuth suna ajiye token ɗinsu tare da bayanan shigar asusun, don haka ka ba su izini sau ɗaya a kowane asusu.

Idan ka fi so sabobin su bayyana a `claude mcp list` da `/mcp` a matsayin sabobin mai amfani na yau da kullum, saita `mcp_strategy = "sync"` a `~/.config/ccam/config.toml`. Sannan ccam yana rubuta su a cikin `.claude.json` na kowane asusu kafin farawa, ƙarƙashin kullen fayil na Claude Code kansa kuma tare da kwafin ajiya.

## Asusai biyu a folda ɗaya

Asusan da ke raba tarihin tattaunawa suna rubutu a fayiloli iri ɗaya. Idan ka fara Claude Code na biyu a foldar da wani asusu na rukuni ɗaya yake aiki a ciki tuni, ccam yana gargaɗe ka kuma yana tambaya ko a ci gaba.

## Idan wani abu ya lalace

```sh
ccam doctor
ccam doctor --fix
```

Idan hanyar da aka raba ta zama fayil na yau da kullum, misali bayan ka maido da kwafin ajiya da hannu, ccam yana tsayawa kafin ya fara Claude Code kuma yana nuna maka `ccam doctor --fix`, wanda ke haɗa fayil ɗin cikin rukunin kuma yana ajiye kwafin ajiya.
