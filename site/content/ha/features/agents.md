---
title: Zaman bayan fage na kowane asusu
description: ccam agents yana lissafa background agents da zaman Claude Code da ke buɗe na dukkan asusanka a tebur ɗaya, kuma ccam attach yana buɗe ɗaya ta ID ɗinsa.
updated: 2026-10-26
order: 16
translation: machine
sourceHash: 06281fb13119
---

Claude Code yana ajiye zaman bayan fage (`claude --bg`, agent view) bisa asusu: `claude agents` yana nuna zaman asusun da yake gudana a matsayinsa kawai. `ccam agents` yana tambayar kowane asusu a lokaci ɗaya kuma yana lissafa su tare, kuma `ccam attach <id|name>` yana buɗe zaman bayan fage a matsayin asusun da ya fara shi, ko wane ne.

## Amfani

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` shi ne yanayin zaman bayan fagen (`working`, `blocked`, `done`, `failed`, `stopped`) da abin da process ɗinsa ke yi (`busy`, `waiting`, `idle`); zaman mu'amala da ke buɗe yana nuna `open`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--all` | Lissafa zaman bayan fage da suka gama, suka gaza ko aka dakatar ma | a kashe |
| `--json` | Buga zaman a matsayin JSON, kowanne tare da `account` da yake nasa | a kashe |

`ccam attach <id|name> [claude args]` yana miƙa duk abin da ke bayan ID ɗin ga `claude attach` ba tare da canji ba, har da flags: `ccam attach 49ec --dry-run` yana miƙa `--dry-run` ga Claude Code. Don buga umarnin a maimakon haka, saka `--dry-run` kafin ID ɗin (`ccam --dry-run attach 49ec` ko `ccam attach --dry-run 49ec`).

## Yadda yake aiki

1. Ga kowane asusu, ccam yana gudanar da `claude agents --json` (tare da `--all` idan an ba da shi) da `CLAUDE_CONFIG_DIR` na asusun, a lokaci guda, kowanne da iyakar lokaci ta daƙiƙa 15. Wannan yana karanta jerin zaman asusun kansa kuma ba ya fara sabis na bayan fage na Claude Code.
2. Kowace shigarwa tana samun sunan asusun, kuma ana jera jerin bisa lokacin farawa, mafi sabo da farko. Ana ambaton asusun da bai amsa ba a standard error kuma a bar shi a waje, sannan ccam yana fita da lamba 1, don kada rubutun umarni (script) ya yi kuskuren ɗaukar gazawa a matsayin jeri marar komai.
3. `ccam attach` yana neman `jobs/<id>` a foldar kowane asusu, inda Claude Code ke ajiye zaman bayan fage, kuma yana gudanar da `claude attach <id>` a matsayin wannan asusun, ta hanya ɗaya da `ccam <name>` ke fara Claude Code. Kalmar da ba farkon (prefix) wani ID ba ce ana kwatanta ta da sunayen zaman, ba tare da la'akari da manya ko ƙananan haruffa ba; kalmar da ke kama da ID ba a taɓa ɗaukar ta a matsayin suna ba.

## Tsaro da yanayi na musamman

- ID na zaman bayan fage na asusu ɗaya ne: `claude attach` da aka gudanar a matsayin wani asusu ba zai iya samun su ba, kuma wannan shi ne abin da `ccam attach` ke warwarewa.
- Idan zaman a asusai biyu suna da farkon ID iri ɗaya ko kalmar suna iri ɗaya, `ccam attach` yana lissafa su kuma yana neman ID ko suna mafi tsawo.
- ccam yana miƙa sabobin MCP ɗinsa da aka raba ga Claude Code da `--mcp-config=<file>` a argument ɗaya, wanda `claude agents` da `claude attach` ke karɓa; zaman bayan fage suna ci gaba da amfani da wannan fayil idan sun sake farawa.
- Sabis na bayan fage na asusun API key yana ɗaukar key ɗinsa daga harsashin (shell) da ya fara shi. Fara zaman bayan fage na farko na irin wannan asusu ta ccam (`ccam gw --bg …`), don sabis ɗin ya sami key ɗin.

## Masu alaƙa

- [Fara Claude Code](/features/launch)
- [Zaɓar asusu don rubutun umarni](/features/pick)
- [Dashboard](/features/dashboard)
