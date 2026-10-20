---
title: Kwafar saitin da ke akwai
description: Kwafi saituna, CLAUDE.md, extensions, sabobin MCP da shigar sabobin MCP na ~/.claude da ke akwai zuwa rukuni na ccam da ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 7cee52e75749
---

Idan ka yi amfani da Claude Code kafin ccam, saitinka yana cikin `~/.claude` da `~/.claude.json`. `ccam migrate` yana kwafar sa zuwa rukunin rabawa, don kowane asusu da ke raba waɗannan abubuwa ta rukunin ya fara da saitin da kake da shi tuni. Ana karanta foldar tushe ne kawai, kuma `claude` na yau da kullum yana ci gaba da aiki kamar da. Gudanar da shi sau ɗaya bayan `ccam add` ɗinka na farko, kuma ka sake gudanar da shi duk lokacin da wani asusu ya fara raba sabobin MCP.

## Amfani

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

A terminal, ba tare da flags na abubuwa ba kuma ba tare da `--yes` ba, ccam yana fara nuna jerin zaɓi (checklist):

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (ko `x`) yana sauya abin da ke ƙarƙashin cursor, maɓallan kibiya ko `j`/`k` suna matsawa, Enter yana nuna shirin, kuma Esc, `q` ko Ctrl+C yana buga `Cancelled.`. Abubuwan da tushen ba shi da komai a kansu ana yi musu alamar `not found`. Idan ka cire alamar komai, ccam yana buga `Nothing selected.` kuma yana tsayawa.

Sannan ccam yana buga shirin kuma yana tambaya kafin ya kwafa:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Rukuni yana kai wa asusansa ne kawai, kuma ga abubuwan da yake rabawa kawai, don haka taƙaitawar tana lissafa asusan rukunin. Idan rukunin ba ya raba ɗaya daga cikin abubuwan da aka kwafa, ko wasu asusai ba sa cikinsa, ccam yana buga umarnan `ccam group set <group> --add <items>` da `ccam group join <group> <accounts>` da ke canza hakan.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `[folder]` | Foldar Claude Code da za a kwafa daga gare ta | `~/.claude` |
| `--settings` | Kwafi `settings.json` | a kashe |
| `--memory` | Kwafi `CLAUDE.md` na gaba ɗaya | a kashe |
| `--extensions` | Kwafi `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | a kashe |
| `--mcp` | Kwafi sabobin MCP na mai amfani (`mcpServers` a `.claude.json`) | a kashe |
| `--history` | Kwafi tarihin tattaunawa (`projects/`, `history.jsonl`) | a kashe |
| `-g`, `--group string` | Rukunin rabawa da za a kwafa zuwa cikinsa | `default` |
| `-y`, `--yes` | Kwafa ba tare da tambaya ba; idan an yi karo ana riƙe kwafin rukunin | a kashe |
| `--dry-run` (na gaba ɗaya) | Buga shirin kuma tsaya kafin kowace tambaya ko kwafi | a kashe |

Ba tare da flag na abu ba, jerin zaɓin ne ke yanke shawara a terminal. A wajen terminal, ko da `--yes`, ccam yana kwafar komai sai tarihin tattaunawa. Ana barin tarihi ta tsohuwa domin `projects/` zai iya yin girma.

`--group on` yana nufin `default`. Ana ƙin `--group off` da `--group local` da `--group takes a group name, not "off"`. `ccam migrate` yana karɓar `--dry-run`, kuma jerin zaɓin har yanzu yana bayyana a dry run don ka iya duba zaɓi kafin lokaci.

## Yadda yake aiki

`internal/cli/migrate.go` yana zaɓar abubuwan, sannan `migrate.Build` yana karanta tushen da rukunin kuma yana mayar da shiri. Ba a rubuta komai har sai `Plan.Apply`, bayan ka tabbatar.

1. **Nemo tushen.** `checkSource` yana ƙin hanyar da ta ɓace, ba folda ba ce, ko ba ta da shigarwar Claude Code da aka sani. `migrate.ClaudeJSON` yana nemo `.claude.json` na foldar: a cikinta, ko ga `~/.claude` da ba shi da ɗaya a ciki, `~/.claude.json` kusa da shi. A nan ne Claude Code ke ajiye shi idan ba a saita `CLAUDE_CONFIG_DIR` ba.
2. **Danganta abubuwa da hanyoyi.** Kowace hanya tana fitowa daga `share.Layout`, jerin izini ɗaya da rukunoni ke amfani da shi. Ana karanta sabobin MCP daga `.claude.json` da `mcp.LocalServers`, wanda ke amfani da scanner na `profile.TopLevelValue` maimakon karanta dukkan fayil ɗin.
3. **Kwafi, ba matsawa ba.** Fayiloli suna zuwa `~/.local/share/ccam/groups/<group>/` da `fsutil.AtomicWrite` ko `fsutil.CopyFile`; sabobin MCP suna shiga `mcp.json` na rukunin.

Kowane abu yana danganta hanyar tushe ɗaya da hanyar rukuni ɗaya:

<!-- motion:migrate-items -->

Layukan tarihi biyu an yanke su (dashed) domin ana kwafar su ne kawai da `--history`. An lissafa `.credentials.json` da `.claude.json` a `share.Never` kuma ba sa taɓa kai wa rukunin: `Build` yana mayar da kuskure na ciki idan ɗaya ya taɓa bayyana a tsarin, kuma zagayen foldar yana tsallake su ta suna.

`buildEntry` yana ba kowace hanya yanayi, wanda shi ne abin da kowane layin shiri ke bayyanawa:

<!-- motion:migrate-states -->

- **Missing**: ba ya cikin tushen, ko fayil ne inda ake tsammanin folda. Ana lissafa shi sau ɗaya ga kowane abu a matsayin `not found, skipped`.
- **Linked**: hanyar tushen ta riga tana kai wa kwafin rukunin, misali idan tushen foldar profile ce ta ccam.
- **Same**: rukunin ya riga yana da abu iri ɗaya.
- **New**: rukunin ba shi da shi, don haka ana kwafar sa gaba ɗaya.
- **Merge**: folda a ɓangarorin biyu. `copyMissing` yana ƙara fayilolin da rukunin ba shi da su kawai; fayilolin da suka bambanta suna zama yadda rukunin yake da su kuma shirin yana cewa `group copies kept`. Ga `plugins/`, ana haɗa rajistocin `installed_plugins.json` da `known_marketplaces.json`, kuma ana sake rubuta cikakkun hanyoyin plugin don su nuna foldar rukunin.
- **Conflict**: fayil a ɓangarorin biyu mai abubuwa daban. Kai ne ke zaɓa ga kowane fayil: `keep group`, `use mine`, da kuma `merge JSON` ga fayilolin `.json` (haɗi mai zurfi inda rukunin ke nasara a ƙimomin da suka yi karo) ko `merge by line` ga `history.jsonl`. Sabar MCP da ma'anarta ta bambanta ita ma rikici ce, da `keep group` ko `use mine`. Zaɓi na farko da aka lissafa shi ne tsoho idan ka danna Enter; `--yes` koyaushe yana riƙe kwafin rukunin.

Jerin zaɓin kansa shi ne `tui.MigratePicker`, wanda aka gina a kan ɓangaren `checklist` da `ccam add` ke amfani da shi; `migrate.Found` yana yi wa abubuwan da suka ɓace alama da kiran lstat da jerin sabobin MCP.

### Shigar sabobin MCP

Sabobin MCP da ke shiga da OAuth suna ajiye token ɗinsu a fayil ɗin shiga, `.credentials.json` ƙarƙashin `mcpOAuth`, ba a jerin sabobin ba. Claude Code yana sa wa kowace shigarwa maɓalli `<name>|` tare da haruffa 16 na farko na hash na `sha256` na `type`, `url` da `headers` na sabar. Maɓallin bai dogara da foldar ba, don haka shigarwar da aka kwafa tana aiki a wani asusu muddin ma'anar sabar iri ɗaya ce. Saboda haka `buildLogins` a `internal/migrate/logins.go` yana kwafar shigarwar sabobin da rukunin zai ƙare da su a ma'anar tushen kawai: sababbi, waɗanda suka riga sun zama iri ɗaya, ko rikicin da ka zaɓi `use mine`. Shigar sabobin plugin (maɓallai masu farawa da `plugin:`) suna tafiya tare idan ka kwafi extensions kuma tushen yana da `plugins/`.

<!-- motion:migrate-signins -->

Shigarwar suna shiga `.credentials.json` na kowane asusu da ke karɓa, ana rubuta su da `fsutil.AtomicWrite` a mode `0600`. Ba sa taɓa shiga rukunin. Asusu yana karɓar su ne kawai idan yana raba `mcp` (ko, ga sabobin plugin, `extensions`) ta rukunin da ake nufi, kuma maɓallan da ba shi da su kawai ake ƙarawa: shigarwar da yake da su tuni, wataƙila sababbi, ana riƙe su. Layin shiri `sign-ins` yana lissafa sabobin, asusan da ke karɓa da waɗanda aka tsallake. ccam yana tsallake asusu, kuma yana ambatar sa tare da dalili, idan:

- an yi masa rajista a wurinsa da `ccam import`, tunda `claude` na yau da kullum zai iya amfani da fayil ɗin shigarsa;
- yana aiki, tunda Claude Code yana sake rubuta fayil ɗin shigarsa (ana sake dubawa kafin rubutu, ana bayar da rahoto a matsayin `started meanwhile`);
- ba a shiga da shi ba tukuna;
- yana ajiye shigarsa a Keychain na macOS, wanda ccam ba ya rubutawa.

Shiga sabobin a asusun da aka tsallake da `/mcp`. Asusun da ya fara raba sabobin MCP daga baya yana samun shigar idan ka sake gudanar da `ccam migrate`; taƙaitawar tana faɗin haka da ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Tsaro da yanayi na musamman

- Ana karanta tushen ne kawai. Ba a matsar, haɗa ko canza kowane fayil a cikinsa, kuma ba a taɓa rubuta shigar tushen, har a macOS inda za a iya karanta ta daga Keychain.
- Kafin a maye gurbin fayil na rukuni, ta `use mine`, haɗawa, ko sababbin sabobin MCP a `mcp.json` da ke akwai, ana ajiye kwafin rukunin a `<group>/.ccam-backup/<time>/`. Sunan foldar na musamman ne ga kowace gudanarwa, don haka gudanarwa biyu a daƙiƙa ɗaya ba sa share ajiyar juna. Taƙaitawar tana buga hanyar.
- Da `--yes`, kowane rikici yana riƙe kwafin rukunin, kuma ccam yana faɗin yadda za a zaɓa maimakon haka: gudanar da umarni ɗaya ba tare da `--yes` ba.
- Idan babu abin da zai canza, ccam yana buga `Nothing to copy.`; idan kowane rikici yana riƙe kwafin rukunin, `Nothing to copy: the group keeps its own copies.`
- Idan kwafi ya gaza a tsakiya, ccam yana tsayawa da `The source folder was not changed. Group files replaced so far are saved in` sannan foldar ajiyar.
- Tushen da ya ɓace yana ƙarewa da `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Masu alaƙa

- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Dubawa da gyarawa](/features/doctor)
- [Asusun Claude Code da dama](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
