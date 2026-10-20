---
title: Raba sabobin MCP, saituna da tarihi
description: Yadda ccam ke raba sabobin MCP, settings.json, CLAUDE.md, agents, skills da tarihin tattaunawa tsakanin asusun Claude Code, yayin da kowane shiga yake daban.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 888d466b29e9
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

Wannan shafin yana bayani kan ayyukan yau da kullum. An bayyana kowane zaɓi a shafukan fasali: [rabawa tsakanin asusai](/features/share), [sabobin MCP da aka raba](/features/mcp) da [kwafar saitin da ke akwai](/features/migrate).

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

ccam ba ya canza rabawa yayin da Claude Code ke gudana a matsayin asusun, saboda zaman yana ci gaba da rubutu a cikin fayilolin da ake motsawa. Ka rufe shi tukuna.

## Kashe rabawa

```sh
ccam share work --history off
```

Ana cire mahaɗin, kuma asusun yana samun kwafinsa na abin da ke cikin rukunin a yanzu. Ƙara `--fresh` don farawa babu komai a maimakon haka. Rukunin yana riƙe abin da ke cikinsa don sauran asusai.

## Kwafi saitin da ke akwai

Idan ka yi amfani da Claude Code kafin ccam, saitunanka suna cikin `~/.claude`. `ccam migrate` yana kwafar su zuwa rukuni, don kowane asusu da ke raba waɗannan abubuwa ya same su:

```sh
ccam migrate
```

A terminal, ccam yana fara nuna jerin zaɓin abubuwan (checklist). An yi wa komai alama sai tarihin tattaunawa; danna Space don canza abu, da Enter don ganin shirin. Abubuwan da foldar ba ta da komai a kansu ana yi musu alamar `not found`.

Abin da yake kwafa:

- `settings.json` da `CLAUDE.md` na gaba ɗaya;
- extensions `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` da `rules/`;
- sabobin MCP naka na mai amfani, wato shigarwar `mcpServers` ta `~/.claude.json`, zuwa `mcp.json` na rukunin;
- bayanan shiga (OAuth) da aka ajiye na sabobin MCP da yake kwafa, zuwa fayil ɗin shiga na kowane asusu;
- tarihin tattaunawa (`projects/` da `history.jsonl`), amma sai idan ka ƙara `--history`.

Ana karanta foldar tushe ne kawai: ba a matsar, haɗa ko canza komai a cikinta, kuma `claude` na yau da kullum yana ci gaba da aiki kamar da. Ba a taɓa kwafar fayilolin shiga `.credentials.json` da `.claude.json` zuwa rukunin ba.

Sabobin MCP da ke shiga da OAuth suna ajiye token ɗinsu a cikin bayanan shiga, ba a cikin jerin sabobin ba. Don kada ka sake shiga cikinsu, ccam yana kwafar waɗannan bayanan shiga zuwa fayil ɗin shiga na kowane asusu da ke samun sabobin ta rukunin, sai dai idan asusun yana da nasa tuni. Yana tsallake asusun da ke aiki, waɗanda ba a shiga ba tukuna, waɗanda aka yi wa rajista a wurinsu da `ccam import`, ko waɗanda ke ajiye bayanan shigarsu a Keychain na macOS, kuma yana ambatar sunayensu; a cikinsu, shiga sabobin da `/mcp`. Asusun da ya fara raba sabobin MCP daga baya zai sami bayanan shigar idan ka sake gudanar da `ccam migrate`.

ccam yana buga shirin kuma yana tambaya kafin ya kwafi komai:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

A ƙarshe, ccam yana lissafa asusan da ke raba kowane abu ta rukunin. Idan wasu asusai ba sa yi, yana buga umarnin `ccam share` da ke kunna rabawa a gare su.

Idan rukunin yana da wani kwafi na fayil da ya bambanta, ccam yana tambaya, kamar yadda yake yi lokacin da ka kunna rabawa: riƙe kwafin rukunin ko amfani da naka. Ana iya kuma haɗa `settings.json` da sauran fayilolin JSON, da `history.jsonl` layi bayan layi. Ana ajiye kwafin rukunin na baya a `.ccam-backup/<time>/` cikin foldar rukunin.

| Zaɓi | Tasiri |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Kwafi abubuwan da ka ambata kawai, ba tare da jerin zaɓin ba; ba tare da su ba kuma a wajen terminal, komai sai tarihin tattaunawa |
| `-g`, `--group <group>` | Kwafa zuwa wannan rukunin maimakon `default` |
| `-y`, `--yes` | Kwafa ba tare da tambaya ba; idan fayil ya bambanta, riƙe kwafin rukunin |

Sa `--dry-run` kafin umarnin, kamar a `ccam --dry-run migrate`, don ganin shirin ba tare da kwafar komai ba. Don kwafa daga wata foldar, ambaci ta: `ccam migrate ~/.claude-work --group work`.

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

Idan rukunin ya riga yana da wata sabar daban mai wannan sunan, ccam yana tsallake ta kuma yana faɗa. `--force` yana maye gurbin sabar rukunin kuma yana ajiye tsohon `mcp.json` a cikin foldar `.ccam-backup/` na rukunin.

Sabobin MCP da ke shiga da OAuth suna ajiye token ɗinsu tare da bayanan shigar asusun, don haka ka ba su izini sau ɗaya a kowane asusu, ko ka kwafi bayanan shigar daga `~/.claude` da `ccam migrate`.

Idan ka fi so sabobin su bayyana a `claude mcp list` da `/mcp` a matsayin sabobin mai amfani na yau da kullum, saita `mcp_strategy = "sync"` a `~/.config/ccam/config.toml`. Sannan ccam yana rubuta su a cikin `.claude.json` na kowane asusu kafin farawa, ƙarƙashin kullen fayil na Claude Code kansa kuma tare da kwafin ajiya.

## Asusai biyu a folda ɗaya

Asusan da ke raba tarihin tattaunawa suna rubutu a fayiloli iri ɗaya. Idan ka fara Claude Code na biyu a foldar da wani asusu na rukuni ɗaya yake aiki a ciki tuni, ccam yana gargaɗe ka kuma yana tambaya ko a ci gaba.

## Idan wani abu ya lalace

```sh
ccam doctor
ccam doctor --fix
```

Idan hanyar da aka raba ta zama fayil na yau da kullum, misali bayan ka maido da kwafin ajiya da hannu, ccam yana tsayawa kafin ya fara Claude Code kuma yana nuna maka `ccam doctor --fix`, wanda ke haɗa fayil ɗin cikin rukunin kuma yana ajiye kwafin ajiya.
