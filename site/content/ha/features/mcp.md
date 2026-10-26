---
title: Sabobin MCP da aka raba
description: Riƙe jerin sabobin MCP ɗaya ga kowane rukuni da ccam mcp add, add-json, list, remove da import, ana miƙa su da --mcp-config ko a daidaita su cikin kowane asusu.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: 27a06ad35295
---

Claude Code yana ajiye sabobin MCP na mai amfani a `.claude.json`, fayil ɗaya da ke riƙe shigar asusun, don haka ba za a iya raba fayil ɗin tsakanin asusai ba. ccam yana ajiye sabobin da aka raba a fayil na rukuni, `mcp.json`, kuma yana miƙa shi ga kowane asusun rukunin lokacin farawa. Yi amfani da `ccam mcp` don ƙara sabar sau ɗaya ga dukkan asusanka, da `ccam mcp import` don matsar da sabar da ka ƙara a cikin zama ɗaya zuwa sauran.

## Amfani

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
```

Canji yana aiki daga farawa ta gaba ta kowane asusu; zaman da ke gudana suna riƙe sabobin da suka fara da su:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

A ƙarƙashin rukunonin, `ccam mcp list` yana ambaton kowace sabar da ba ta kai wa sauran asusai:

- `only in account work: pg` ga sabar da asusu ya ƙara wa kansa;
- `only in account side, for /home/me/code/app only: notes` ga wadda aka ƙara don folda ɗaya (`claude mcp add` na yau da kullum, duba ƙasa);
- `account work has its own …, used in its sessions in place of the group's` idan sabar asusu ta kansa tana da sunan sabar rukuni;
- `.mcp.json` na wannan foldar, tare da asusan da take `waiting for approval in` amincewarsu;
- sabobin da suke `in no account, only in ~/.claude.json (plain claude's)`, tare da umarnin [`ccam migrate --mcp`](/features/migrate) da ke kwafar su.

Ƙara sunan da rukunin ya riga yana da shi yana maye gurbin ma'anarsa kuma yana buga `updated` maimakon `added`. Sabobin suna kai wa asusu ne kawai idan yana cikin wannan rukunin kuma rukunin yana raba MCP (`ccam group set <group> --add mcp`); idan rukunin ba ya yi, `add` yana faɗin haka.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `-g`, `--group <group>` | rukunin da ake nufi don `add`, `add-json` da `remove` | rukuni ɗaya tilo da wani asusu ke raba MCP ta cikinsa, in ba haka ba `default` |
| `-- <command> [args…]` (`add`) | gudanar da wannan umarni a matsayin sabar stdio | |
| `--url <url>` (`add`) | URL na sabar nesa (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | hanyar sadarwa (transport) ta sabar `--url` | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, ana iya maimaitawa | babu |
| `-e`, `--env <KEY=VALUE>` (`add`) | mai canjin muhalli ga sabar stdio, ana iya maimaitawa | babu |
| `--force` (`import`) | maye gurbin sabar rukuni mai suna ɗaya da ta bambanta; ana ajiye kwafin tsohon `mcp.json` | a kashe |
| `--keep-local` (`import`) | riƙe shigarwar a `.claude.json` na asusun | a kashe |
| `-h`, `--help` | taimako ga kowane umarni | |
| `--dry-run` (na gaba ɗaya) | `mcp list` kaɗai ke tallafa masa; `add`, `add-json`, `remove` da `import` suna ƙin sa kuma ba sa canza komai | a kashe |

`add` yana buƙatar ko dai umarni bayan `--` ko `--url`; in ba haka ba yana gazawa da `a command to run (after --) or --url is required`. `--header` yana aiki ga sabobin `--url`, kuma `--env` ga sabobin umarni. `add-json` yana karɓar JSON iri ɗaya da `claude mcp add-json`, kuma dole ma'anar ta zama abu (object) mai `command` ko `url`. `import` koyaushe yana amfani da rukunin da asusun ke raba MCP ta cikinsa, kuma `list` yana nuna kowane rukuni, don haka `-g` ba ya canza su.

## Yadda yake aiki

Fayilolin da abin ya shafa:

- `~/.local/share/ccam/groups/<group>/mcp.json`: sabobin da aka raba, `{"mcpServers": {…}}`, tsarin da `--mcp-config` ke karantawa. Ana rubuta shi da `fsutil.AtomicWrite` da mode `0600`, domin shigarwar sabobi sau da yawa suna ɗauke da API keys, kuma ƙarƙashin kullen `<group>/.mcp.json.lock` kawai, don haka umarnan `ccam mcp add` biyu da aka gudanar a lokaci ɗaya dukansu suna riƙe sabobinsu.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: fayil na Claude Code kansa mai sabobin asusun. ccam yana karanta ƙimar `mcpServers` ɗinsa kawai da scanner na `profile.TopLevelValue` kuma ba ya taɓa karanta dukkan fayil ɗin a hanyar farawa.
- `~/.local/state/ccam/mcp/<name>/`: fayil ɗin `seen`, `sync-stamp` (tare da hash na kowane kwafi da ccam ya rubuta, don ya bambanta kwafin da ba wanda ya canza da wanda ka gyara), `sync-checked` (abin da dabarar sync ta gani na ƙarshe), `kept` (sabobin da asusun ya riƙe lokacin da ya bar rabon MCP) da kwafin ajiya na `.claude.json`.

Akwai dabaru biyu, ana saita su ta `mcp_strategy` a `~/.config/ccam/config.toml`. Tsohuwar ita ce dabarar flag:

<!-- motion:mcp-flag -->

1. ccam yana nemo rukunin MCP na asusun (`launcher.applyMCP`) kuma yana dubawa da `stat` ɗaya cewa `mcp.json` na rukunin yana nan, ba tare da karanta shi ba (`mcp.FlagArgs`).
2. ccam yana saka `--mcp-config=<path>` kai tsaye bayan `claude`, kafin arguments naka, don haka `--mcp-config` ko `--strict-mcp-config` naka har yanzu suna aiki a saman sa (`Plan.InsertArgs`).
3. ccam yana duba ko `.claude.json` na asusun ya canza tun farawa ta ƙarshe, kuma sai a lokacin ne yake karanta sabobinsa (`mcp.NewLocal`).
4. ccam yana maye gurbin kansa da Claude Code (`syscall.Exec`).

Ana miƙa hanyar a matsayin argument ɗaya, `--mcp-config=<path>`, ba biyu ba. `--mcp-config` na Claude Code yana karɓar ƙima da dama, don haka `claude --mcp-config a.json "prompt"` zai karanta prompt ɗin a matsayin fayil ɗin saiti na biyu kuma ya gaza da `MCP config file not found`. Nau'in argument ɗaya yana hana wannan, don haka prompt ko ƙaramin umarni da ke biye da shi yana aiki. `ccam work mcp list`, `get` da `remove` ba sa samun flag, tunda Claude Code yana ƙin shi a can; maimakon haka ccam yana ambaton sabobin rukunin a standard error ([cikakken bayani](/features/launch#shared-mcp-servers)). Claude Code yana loda sabobin `--mcp-config` ƙari ga na asusun kansa; idan sunaye sun yi karo, shigarwar da ke `.claude.json` ce ke nasara. Dabarar flag tana rubuta `.claude.json` ne kawai don Remote Control, idan asusu ya bar rabon MCP ko ya sake shiga ([a ƙasa](#leaving-and-joining)), da kuma don `ccam mcp import`.

`claude remote-control` (`rc`) ba ya karɓar `--mcp-config`, kuma zamansa suna karanta sabobi ne kawai daga `.claude.json`. Don `ccam work rc`, ccam yana barin flag ɗin kuma, a terminal, yana tambaya ko zai kwafi sabobin rukunin cikin `.claude.json`, tare da ajiya, kamar yadda dabarar `sync` ke yi. Yana tambaya har sai kwafi ya kasance, sannan sai kawai bayan `mcp.json` na rukunin ya canza ko sabar da aka kwafa ta ɓace ([Remote Control](/features/launch#remote-control)). Saboda irin wannan kwafi yana nasara a kan `--mcp-config` idan sunaye sun yi karo, `ccam mcp add`, `add-json`, `remove` da `import` suna sabunta kwafin da sauran asusan rukunin ke riƙe da su:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Layin `✓ Updated` yana bayyana ne kawai idan kwafi ya canza. Kwafin da tsohuwar sigar ccam ta ƙirƙira ana sabunta shi ta hanya ɗaya, haka ma wanda ke bayan `ccam migrate`.

Sabar da aka ƙara a cikin zama da `claude mcp add` tana sauka a `.claude.json` na wannan asusun kawai. `claude mcp add` na yau da kullum yana amfani da local scope: sabar tana lodawa ga wannan asusu ɗaya a wannan folda ɗaya (`-s user` yana sa ta samu ga asusun a kowace folda). ccam ba zai iya ganin wannan lokacin da ya faru ba, domin Claude Code ya riga ya maye gurbin process na ccam, don haka yana dubawa a farawa ta gaba. Fayil ɗin `seen` yana riƙe lokacin gyaran `.claude.json` da sunayen sabobin da aka riga aka sanar. Idan lokacin ya bambanta, ccam yana kwatanta sunayen da rukunin kuma yana buga shawara ga kowace sabuwar sabar sau ɗaya, ba tare da jiran shigarwa ba, a dabarun biyu:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Shawarar tana kuma ambaton sabar asusu ta kansa da ke da sunan sabar rukuni, tunda ta asusun ce ke nasara a zamansa.

Da `mcp_strategy = "sync"`, sabobin suna bayyana a `claude mcp list` da `/mcp` a matsayin sabobin mai amfani na yau da kullum, kuma ccam yana rubuta su cikin `.claude.json` na kowane asusu kafin farawa:

<!-- motion:mcp-sync -->

1. ccam yana kwatanta lokacin gyaran `mcp.json` na rukunin da `sync-stamp`, da `.claude.json` na asusun da abin da ya gani na ƙarshe (`mcp.Sync`). Idan babu wanda ya canza, babu wani abu da ke faruwa. In ba haka ba yana kuma duba cewa kowace sabar rukuni har yanzu tana cikin `.claude.json`, kuma yana mayar da kwafin da ya ɓace (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Yana ɗaukar kullen Claude Code kansa, foldar `.claude.json.lock` kusa da fayil ɗin, wadda aka ƙirƙira da `mkdir`. Yana jira har zuwa daƙiƙa 2 kuma yana ɗaukar kullen da ya wuce daƙiƙa 10 a matsayin tsoho (stale), kamar yadda laburaren kulle na Claude Code ke yi.
3. Yana ajiye `.claude.json` na yanzu a `backups/` a foldar state kuma yana riƙe sababbi 5.
4. Yana saita kowace sabar rukuni a `mcpServers`, yana cire sabobin da ya rubuta a baya waɗanda rukunin ba shi da su yanzu, kuma yana rubuta fayil ɗin a atomic da mode ɗinsa na asali, sai idan wani abu ya canza. Yana canzawa ko cire kwafin da shi ya rubuta waɗanda ba wanda ya canza tun daga lokacin kawai. Idan asusun yana da sabarsa mai suna iri ɗaya da sabar rukuni, ccam yana riƙe sabar asusun kuma yana faɗin haka sau ɗaya: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Yana rubuta sabon stamp tare da hash na kowane kwafi, kuma Claude Code yana farawa ba tare da `--mcp-config` ba.

Ana tsallake sync, da `ccam: account work is already running, skipping MCP sync`, idan zaman asusun yana gudana, ko ta yaya aka fara Claude Code, domin Claude Code yana sake rubuta `.claude.json` yayin da yake aiki. Sync da ya gaza yana buga `ccam: MCP sync failed (launching anyway): …` kuma yana farawa duk da haka.

### Barin rabon MCP da komawa

Asusun da ya daina raba MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) yana riƙe sabobin rukunin: shirin yana kwafar su cikin `.claude.json` nasa, tare da ajiya, kuma sabobinsa na kansa masu suna iri ɗaya suna nan. Da `--fresh` ba ya samun ko ɗaya, kuma shirin yana lissafa sabobin da za su daina lodawa. Idan ya sake shiga rukunin da ke da waɗannan sabobin, ana cire kwafin da aka riƙe waɗanda ba wanda ya canza, don gyare-gyaren rukunin na gaba su kai gare shi.

`ccam mcp import` yana matsar da sabobin asusu na kansa cikin rukuninsa:

<!-- motion:mcp-import -->

ccam yana karanta sabobin asusun, nasa na kansa da waɗanda yake da su don folda ɗaya kawai, kuma yana kwatanta kowacce da sabar rukunin mai suna ɗaya ta amfani da `mcp.SameDef`, wanda ke watsi da tsarin maɓallai da sarari. Sabar folda tana zama sabar rukuni ga kowace folda (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); ana ƙin suna ɗaya mai ma'anoni biyu daban-daban na folda. Ba a ba da sabobin plugin ba, tunda plugin ɗin ne ke kawo su. Ana kwafar sabar da rukunin ba shi da ita ko yake da ita da ma'ana ɗaya. Ana tsallake wadda ta bambanta, domin ma'anar rukunin ce kowane sauran asusu ke amfani da ita:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Da `--force` ana maye gurbin sabar rukunin, kuma ana fara ajiye `mcp.json` na baya a `.ccam-backup/<time>/mcp.json` a cikin foldar rukunin. Sannan ccam yana cire shigarwar da aka kwafa daga `.claude.json` na asusun, nasa na kansa da na foldar a rubutu ɗaya, ƙarƙashin kulle ɗaya kuma bayan ajiya a foldar state (`mcp.Remove`). ccam yana riƙe `mcp.json` na rukunin a kulle tun daga lokacin da ya karanta fayil ɗin har sai an cire waɗannan shigarwar, don kada wani ya cire sabar daga rukunin a tsakani. Ana buƙatar wannan domin shigarwar asusun za ta yi nasara a kan ta rukunin idan sunaye sun yi karo. Da `--keep-local`, shigarwar suna zama a `.claude.json` na asusun.

## Tsaro da yanayi na musamman

- **Shiga.** Sabobin da ke shiga da OAuth suna ajiye token ɗinsu a `.credentials.json` na kowane asusu (`mcpOAuth`), ba a jerin sabobin ba, don haka raba sabar ba ya raba shigarta. Shiga sau ɗaya a kowane asusu da `/mcp`. ccam ba ya taɓa kwafar shiga, haka ma [`ccam migrate`](/features/migrate): refresh token da aka kwafa yana fitar da kowane mai riƙe da shi. [`ccam ext`](/features/ext) da `ccam doctor` suna ambaton sabobin da asusu bai shiga ba alhali wani asusu ya shiga, kuma `ccam next --continue` yana faɗin waɗanda sabon asusun ba shi da su. Rukunin ba ya taɓa samun fayil ɗin shiga.
- **Asusan da ke aiki.** `import` yana kwafa zuwa rukunin amma yana barin shigarwar asusun da ke aiki kuma yana faɗin haka; sake gudanar da shi bayan zaman ya ƙare. Sync yana tsallake asusan da ke aiki. Dukansu suna kuma gano zaman Claude Code da aka fara ba ta ccam ba.
- **Kwafi don Remote Control.** Ba a taɓa kwafin asusun da ke aiki ba: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Da `mcp_strategy = "sync"` babu abin da ke buƙatar sabuntawa, domin kowace farawa tana yin sync.
- **Fayilolin da suka lalace.** Idan `.claude.json` ba ingantaccen JSON ba ne, sync yana barin sa yadda yake kuma yana bayar da rahoton kuskuren. Idan ba za a iya ɗaukar kullen cikin daƙiƙa 2 ba, ccam yana bayar da rahoton `could not lock … (Claude Code is writing the file)`.
- **Asusan da ba sa raba MCP.** `import` yana ƙi kuma yana gaya maka ka fara saka asusun a rukunin da ke raba MCP (`ccam group join <group> side`).
- **Sunayen rukunoni.** `-g` yana karɓar sunan rukuni kawai; yana zama folda ƙarƙashin `groups/`, don haka ana ƙin sunaye kamar `../x`.

## Masu alaƙa

- [Rabawa da rukunoni](/features/share)
- [Amfani da Remote Control tare da ccam](/guides/remote-control)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Kwafar saitin da ke akwai](/features/migrate)
- [Bayanin umarnai](/docs/commands)
