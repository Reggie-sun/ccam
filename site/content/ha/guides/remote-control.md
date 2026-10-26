---
title: Amfani da Remote Control tare da ccam
description: "Remote Control na Claude Code ga kowane asusu da ccam <name> rc: zaɓuɓɓukan da ke aiki, dalilin da ccam ke tayin kwafar sabobin MCP da aka raba, da gyara."
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 4d46f49c5dea
---

`ccam work rc` yana fara Remote Control na Claude Code a matsayin asusun `work`, don ka iya sarrafa zaman da ke kan wannan kwamfuta daga claude.ai/code ko manhajar Claude ta wayar hannu. Yana aiki kamar `claude rc` da bambanci guda ɗaya: idan asusun yana raba sabobin MCP, ccam yana tambaya sau ɗaya ko zai kwafe su cikin `.claude.json` na asusun, domin wannan fayil shi ne wuri ɗaya tilo da zaman Remote Control ke karanta sabobi daga gare shi.

## Fara shi

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[Gajeren umarni](/features/shortcuts) yana aiki haka nan: `ccwork rc` shi ne `ccam work rc`. Gudanar da umarnin a foldar da kake son yin aiki a ciki, sannan ka buɗe claude.ai/code ko manhajar Claude, an shiga da asusun Claude ɗaya da na `work`.

## Zaɓuɓɓuka

Iri biyu na zaɓuɓɓuka suna haɗuwa a nan: na ccam kansa, waɗanda ke zuwa kafin sunan asusun ko suke cikin saitunan ccam, da na Remote Control, waɗanda ke zuwa bayan `rc` kuma suna isa ga Claude Code ba tare da canji ba.

### Zaɓuɓɓuka da saitunan ccam

| Zaɓi ko saiti | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--dry-run` (kafin asusun) | Yana buga muhalli da umarnin, da sanarwar MCP idan akwai. Ba ya tambaya kuma ba ya kwafi | a kashe |
| `mcp_strategy = "sync"` a saman `~/.config/ccam/config.toml` | Yana kwafar sabobin rukunin cikin `.claude.json` kafin kowace farawa, Remote Control ko ba haka ba, don haka ba a taɓa tambaya ba | `"flag"` |
| `ccam group leave <group> work` | Yana cire asusun daga rukunin da ke raba MCP da shi. Ana kwafar sabobin rukunin cikin `.claude.json` nasa a matsayin nasa na kansa, don haka Remote Control yana da su kuma ba a tambayar komai; da `--fresh` yana farawa da sabobin asusun kansa kawai | a cikin `default` |

Idan asusun yana raba sabobin MCP kuma `.claude.json` ɗinsa ba shi da jerin rukunin na yanzu, ko sabar da aka kwafa ta ɓace daga cikinsa, ccam yana tambaya `Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]`. `y` yana kwafar su; `N` ko Enter yana farawa ba tare da su ba.

### Zaɓuɓɓukan Remote Control

Duk abin da ke bayan `rc` na Claude Code ne. Waɗannan su ne zaɓuɓɓukan `claude remote-control` a Claude Code 2.1.304; `ccam work rc --help` yana buga jerin na sigarka.

| Zaɓi | Abin da yake yi |
| --- | --- |
| `--name <name>` | Sunan zaman da ake nunawa a claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Kari a gaban sunayen zaman da ake ƙirƙira (tsoho: sunan kwamfutar mai masaukin) |
| `-c`, `--continue` | Sake haɗawa da zaman da Remote Control ya rubuta na ƙarshe don wannan foldar, cikin kusan awa 4 |
| `--session-id <id>` | Sake haɗawa da wani takamaiman zama |
| `--permission-mode <mode>` | Yanayin izini don sabbin zaman: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` ko `plan` |
| `--spawn <mode>` | `same-dir` (tsoho), `worktree` (git worktree ɗaya ga kowane zama) ko `session` (zama ɗaya, yana fita idan ya ƙare) |
| `--capacity <N>` | Matsakaicin adadin zaman da ke gudana a lokaci ɗaya a yanayin `same-dir` da `worktree` (tsoho 32) |
| `--[no-]create-session-in-dir` | Ƙirƙiri zama na farko a foldar yanzu lokacin farawa (tsoho a kunne) |
| `--[no-]chrome` | Claude in Chrome don sabbin zaman (tsoho a kashe) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Fitarwa ta gyaran kurakurai da ta cikakken bayani |

ccam yana karanta umarnin Claude Code kamar yadda Claude Code ke yi, yana tsallake zaɓuɓɓuka tare da ƙimominsu, don haka `ccam work rc --name laptop` da `ccam work --model opus rc` duka suna fara Remote Control. Zaɓin da ccam bai sani ba ana ƙirga shi a matsayin wanda ba shi da ƙima; don tsira, saka `rc` kai tsaye bayan sunan asusun.

## Abin da ccam ke yi

Farawar tana bin hanyar da aka saba ([Fara Claude Code](/features/launch)) da canji ɗaya a matakin MCP:

1. ccam yana nemo asusun, yana duba mahaɗan rabawarsa kuma yana nemo `claude`, kamar a kowace farawa.
2. Yana ganin cewa umarnin shi ne `rc` ko `remote-control` kuma yana barin `--mcp-config` a waje.
3. Idan asusun yana raba sabobin MCP, yana kwatanta `mcp.json` na rukunin da `sync-stamp` da ke `~/.local/state/ccam/mcp/work/`, kuma yana duba cewa kowace sabar rukuni har yanzu tana cikin `.claude.json`. Idan `.claude.json` yana ɗauke da sabobin na yanzu tuni, ko rukunin ba shi da ko ɗaya, yana tafiya kai tsaye zuwa mataki na 6.
4. In ba haka ba, yana faɗin sunayen sabobin da zaman za su rasa:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

5. Idan `y`, yana ajiye kwafin `.claude.json`, yana ɗaukar kullen Claude Code kansa, yana ƙara sabobin rukunin cikin `mcpServers` kuma yana riƙe sauran sabobin asusun; idan asusun yana da sabarsa mai suna iri ɗaya da sabar rukuni, ccam yana riƙe sabar asusun kuma yana faɗin haka. Idan `N`, ba tare da terminal ba, ko yayin da asusun ke aiki, ba ya canza komai kuma yana faɗin haka:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Yana rubuta zaman, yana saita `CLAUDE_CONFIG_DIR` zuwa foldar asusun kuma yana maye gurbin kansa da `claude rc` ta hanyar `exec`.

Daga nan Remote Control yana ci gaba da aiki a matsayin sabar. Kowane zama da ya fara maka ɗan aiki ne (child process) na wannan aikin, don haka yana aiki a matsayin `work`: da bayanan shiga, saituna da tarihin tattaunawa na `work`, kuma ana lissafa shi a kan iyakar amfani na `work`. ccam yana ɗaukar asusun a matsayin mai aiki yayin da sabar Remote Control da ya fara ke aiki. Sabar da aka fara da `claude remote-control` na yau da kullum tana ƙirguwa ne kawai yayin da take da zaman; ccam ba zai iya ganin wadda ba ta da ko ɗaya tukuna ba.

Gwajin gudu (dry run) yana nuna bambancin. Farawa ta yau da kullum tana ɗauke da sabobin rukunin a matsayin argument, Remote Control ba ya ɗauka:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## Dalilin da ya sa yake aiki haka

**Remote Control ba ya karɓar `--mcp-config`.** A farawa ta yau da kullum ccam yana miƙa sabobin rukunin ga Claude Code a matsayin `--mcp-config=<path>` kuma ba ya taɓa rubuta fayilolin Claude Code. `claude remote-control` yana ƙin wannan zaɓin kafin sunan umarninsa kuma ba ya karɓar ko ɗaya bayansa, don haka kafin ccam 0.6.0, `ccam work rc` yana gazawa ga kowane asusu da ke raba sabobin MCP. Yanzu ccam yana barin zaɓin a waje.

**Zamansa suna karanta sabobi ne kawai daga `.claude.json`.** Zaman da Remote Control ke farawa ba sa gadon zaɓuɓɓuka daga layin umarni. Wuri ɗaya tilo da suke samun sabobin MCP na mai amfani shi ne `.claude.json` na asusun. Wannan fayil ɗin yana ɗauke da bayanan shiga na asusun da yanayin kowace folda, don haka ccam ba zai iya haɗa shi da kwafin da aka raba kamar yadda yake yi da `settings.json` ba. Kwafar sabobin cikinsa ita ce hanya ɗaya tilo ta ba zaman jerin da aka raba.

**ccam yana tambaya kafin ya rubuta.** Da dabarar flag ta tsoho, ccam yana barin `.claude.json` yadda yake a lokacin farawa. Remote Control ita ce farawa ɗaya tilo da dole ne ya rubuta a can, don haka yana fara tambaya, yana ajiye kwafi, kuma ba ya taɓa komai sai `mcpServers`. (A wajen farawa, barin ko sake shiga rabon MCP da `ccam mcp import` su ma suna rubutu a can, kamar yadda shirye-shiryensu suka faɗa.)

**Ba ya rubutu yayin da asusun ke aiki.** Claude Code yana sake rubuta `.claude.json` yayin zama, kuma canjin da aka yi daga waje zai ɓace ko a rubuta a kansa. Shi ya sa ccam yake tsallake kwafin idan kowane zama na asusun yana aiki, har da wata sabar Remote Control da Claude Code da aka fara ba ta ccam ba.

**Yana tambaya sau ɗaya ga kowane canji.** Bayan kwafi, `sync-stamp` yana riƙe lokacin gyaran `mcp.json` na rukunin. ccam yana sake tambaya ne kawai idan wannan fayil ɗin ya canza, misali bayan `ccam mcp add`, ko idan sabar da aka kwafa ta ɓace daga `.claude.json`. Idan ka amsa `N`, babu kwafi tukuna, don haka yana tambaya a kowace farawar Remote Control.

**`ccam mcp` yana kiyaye kwafin a halin yanzu.** Shigarwa a `.claude.json` tana nasara a kan `--mcp-config` idan sunaye sun yi karo, don haka tsohon kwafi zai ɓoye sabon ma'ana ko a farawa ta yau da kullum. Saboda haka `ccam mcp add`, `add-json`, `remove` da `import`, da `ccam migrate`, suna sabunta kwafin kowane asusu da ba ya aiki, kuma suna faɗin haka idan kwafi ya canza:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Ba tare da terminal ba amsar ita ce a'a.** Remote Control yakan gudana daga rubutun (script), sabis ko `nohup`, inda babu wanda zai iya amsawa. Sai ccam ya fara ba tare da sabobin ba maimakon jira.

## Matsaloli da mafita

### Sabobin MCP da aka raba ba su cikin zaman Remote Control

An amsa tambayar da `N`, babu terminal, ko asusun yana aiki. Dakatar da Remote Control, rufe sauran zaman asusun, kuma sake gudanar da `ccam work rc` a terminal; amsa `y`. Don daina tambayarka, saita `mcp_strategy = "sync"` a `~/.config/ccam/config.toml`. Sabobin da ke cikin zama tuni ba sa canzawa yayin da yake aiki, don haka fara sabon zama bayan kwafin.

### "Account work is running, so its .claude.json is left alone"

Akwai zaman `work` da ke buɗe har yanzu: `ccam work` na yau da kullum, gajeren umarni, ko wata sabar Remote Control. Fita daga ciki kuma sake gudanar da `ccam work rc`. Zaman da wani asusu ke gudanarwa ba sa ƙirguwa.

### ccam yana tambaya kowane lokaci

Ba ka kwafi sabobin ba tukuna. Amsa `y` sau ɗaya, ko saita `mcp_strategy = "sync"`. Idan ba ka taɓa son sabobin a Remote Control ba, daina raba MCP da asusun ba tare da riƙe su ba: `ccam group leave default work --fresh`, ko ka ajiye MCP a rukuni nasa na kansa wanda asusun ba ya ciki. `ccam group leave` kawai yana kwafar sabobin rukunin cikin `.claude.json` na asusun, don haka Remote Control yana da su.

### `ccam work rc` yana gazawa nan take da kuskure game da `--mcp-config`

ccam ɗinka ya tsufa fiye da 0.6.0 kuma yana miƙa `--mcp-config` ga Remote Control. Gudanar da `ccam update`. Da ccam na yanzu, kuskure iri ɗaya yana nufin cewa wani zaɓi da ba a sani ba mai ƙima ya zo kafin `rc`, kuma ccam ya yi kuskuren ɗaukar ƙimar a matsayin umarnin; saka `rc` kai tsaye bayan sunan asusun.

### Canjin da aka yi da `ccam mcp` ba ya isa Remote Control

`ccam mcp add` da `remove` suna sabunta kwafi ne kawai a asusan da ba sa aiki, kuma suna faɗin haka: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Dakatar da Remote Control kuma sake fara shi; ccam yana ganin `mcp.json` da ya canza kuma yana sake tambaya.

### Ana amfani da sabata ta kaina maimakon ta rukunin

Idan asusun yana da sabarsa mai suna ɗaya da sabar rukuni, kwafin yana barin ta a wurinta kuma yana faɗin haka: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`. Don yin amfani da ta rukunin, cire ta asusun (`claude mcp remove -s user <name>`); don ba kowane asusu sigarka, matsar da ita cikin rukunin da `ccam mcp import work <server> --force`. Kwafin ajiya na `.claude.json` suna a `~/.local/state/ccam/mcp/work/backups/`, wanda ke riƙe sabbi biyar na ƙarshe.

### Ina son cire sabobin da aka kwafa

`ccam mcp remove <server>` yana cire sabar daga rukunin da kuma daga kowane kwafi. Don komawa yanayin kafin kwafin, maido da sabon kwafin ajiya daga `~/.local/state/ccam/mcp/work/backups/` yayin da babu zaman asusun da ke aiki.

### Remote Control yana tambaya ko a amince da foldar, ko yana fita ba tare da terminal ba

Ana ajiye amincewa da folda ga kowane asusu a `.claude.json`, don haka foldar da ka amince da ita a matsayin `personal` sabuwa ce ga `work`. Gudanar da `ccam work` sau ɗaya a waccan foldar kuma ka amince, sannan ka fara Remote Control a can.

### "You must be logged in with a Claude account that has a subscription"

Remote Control yana buƙatar shiga da asusun Claude mai biyan kuɗi (subscription). Duba imel ɗin asusun a `ccam list`, kuma sake shiga da `ccam login work` idan ya cancanta.

### Zaman ba ya bayyana a waya

Manhajar Claude ko claude.ai/code an shiga da wani asusun Claude daban. Kwatanta shi da imel ɗin da `ccam list` ke nunawa don `work`.

### ccam yana gargaɗin cewa wani zaman Claude Code yana raba tarihin tattaunawa a wannan foldar

Wani asusu a rukunin tarihi ɗaya yana aiki a folda ɗaya. Zaman da ke ci gaba da tattaunawa ɗaya daga asusun biyu na iya cakuɗa tarihin. Amsa `y` kawai idan za su yi aiki a kan tattaunawa daban-daban, ko fara Remote Control a wata folda.

## Masu alaƙa

- [Fara Claude Code](/features/launch)
- [Sabobin MCP da aka raba](/features/mcp)
- [Rabawa da rukunoni](/features/share)
- [Sarrafa asusai](/features/accounts)
