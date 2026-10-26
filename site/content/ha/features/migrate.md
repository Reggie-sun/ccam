---
title: Kwafar saitin da ke akwai
description: Kwafi saituna, CLAUDE.md, extensions da sabobin MCP na ~/.claude da ke akwai zuwa rukuni na ccam da ccam migrate.
updated: 2026-10-26
order: 6
translation: machine
sourceHash: f835fbe84ade
---

Idan ka yi amfani da Claude Code kafin ccam, saitinka yana cikin `~/.claude` da `~/.claude.json`. `ccam migrate` yana kwafar sa zuwa rukunin rabawa, don kowane asusu da ke raba waɗannan abubuwa ta rukunin ya fara da saitin da kake da shi tuni. Ana karanta foldar tushe ne kawai, kuma `claude` na yau da kullum yana ci gaba da aiki kamar da. Gudanar da shi sau ɗaya bayan `ccam add` ɗinka na farko; idan `~/.claude.json` yana da sabobin MCP, `ccam add` na farko yana ambaton su a shawararsa.

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
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Space (ko `x`) yana sauya abin da ke ƙarƙashin cursor, maɓallan kibiya ko `j`/`k` suna matsawa, Enter yana nuna shirin, kuma Esc, `q` ko Ctrl+C yana buga `Cancelled.`. Abubuwan da tushen ba shi da komai a kansu ana yi musu alamar `not found`. Idan ka cire alamar komai, ccam yana buga `Nothing selected.` kuma yana tsayawa.

Sannan ccam yana buga shirin kuma yana tambaya kafin ya kwafa:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Rukuni yana shafar asusansa ne kawai, kuma abubuwan da yake rabawa kawai, don haka taƙaitawar tana lissafa asusan rukunin. Idan rukunin ba ya raba ɗaya daga cikin abubuwan da aka kwafa, ko wasu asusai ba sa cikinsa, ccam yana buga umarnan `ccam group set <group> --add <items>` da `ccam group join <group> <accounts>` da ke canza hakan.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `[folder]` | Foldar Claude Code da za a kwafa daga gare ta | `~/.claude` |
| `--settings` | Kwafi `settings.json` | a kashe |
| `--memory` | Kwafi `CLAUDE.md` na gaba ɗaya | a kashe |
| `--extensions` | Kwafi `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | a kashe |
| `--mcp` | Kwafi sabobin MCP na mai amfani (`mcpServers` a `.claude.json`) | a kashe |
| `--history` | Kwafi tarihin tattaunawa (`projects/`, `history.jsonl`) | a kashe |
| `-g`, `--group <group>` | Rukunin rabawa da za a kwafa zuwa cikinsa | `default` |
| `-y`, `--yes` | Kwafa ba tare da tambaya ba; idan an yi karo ana riƙe kwafin rukunin | a kashe |
| `--dry-run` (na gaba ɗaya) | Buga shirin kuma tsaya kafin kowace tambaya ko kwafi | a kashe |

A terminal, idan ba ka ba da flags na abubuwa ba, kai ne kake zaɓar abubuwan a jerin zaɓin. A wajen terminal, ko da `--yes`, ccam yana kwafar komai sai tarihin tattaunawa. Ana barin tarihi ta tsohuwa domin `projects/` zai iya yin girma.

`--group on` yana nufin `default`. Ana ƙin `--group off` da `--group local` da `--group takes a group name, not "off"`. `ccam migrate` yana tallafa wa `--dry-run`, kuma jerin zaɓin har yanzu yana bayyana a dry run don ka iya duba zaɓi kafin lokaci.

## Yadda yake aiki

ccam yana gina cikakken shiri kafin ya kwafi komai. Yana zaɓar abubuwan (`internal/cli/migrate.go`), sannan yana karanta tushen da rukunin kuma yana gina shirin (`migrate.Build`). Ba a rubuta komai har sai `Plan.Apply`, bayan ka tabbatar.

1. **Nemo tushen.** ccam yana ƙin hanyar da ta ɓace, ba folda ba ce, ko ba ta da shigarwar Claude Code da aka sani (`checkSource`). Don sabobin MCP kaɗai, `~/.claude` da babu komai a ciki ko wanda babu shi ba matsala ba ne muddin `~/.claude.json` yana da sabobi. ccam yana nemo `.claude.json` na foldar (`migrate.ClaudeJSON`): a cikinta, ko ga `~/.claude` da ba shi da ɗaya a ciki, `~/.claude.json` kusa da shi. A nan ne Claude Code ke ajiye shi idan ba a saita `CLAUDE_CONFIG_DIR` ba.
2. **Danganta abubuwa da hanyoyi.** Kowace hanya tana fitowa daga `share.Layout`, jerin izini ɗaya da rukunoni ke amfani da shi. Ana karanta sabobin MCP daga `.claude.json` da `mcp.LocalServers`, wanda ke amfani da scanner na `profile.TopLevelValue` maimakon karanta dukkan fayil ɗin.
3. **Kwafi, ba matsawa ba.** Fayiloli suna zuwa `~/.local/share/ccam/groups/<group>/` da `fsutil.AtomicWrite` ko `fsutil.CopyFile`, wanda ke rubuta fayil na wucin gadi kusa da inda ake nufi sannan ya sauya masa suna, don haka kwafin da aka katse ba ya taɓa barin fayil da bai cika ba. Sabobin MCP suna shiga `mcp.json` na rukunin ƙarƙashin kullensa: sabar da wani process na ccam ya ƙara a wannan lokacin ana riƙe ta, kuma ana sabunta kwafin da aka yi don [Remote Control](/guides/remote-control) bayan haka.

Kowane abu yana danganta hanyar tushe ɗaya da hanyar rukuni ɗaya:

<!-- motion:migrate-items -->

Layukan tarihi biyu an yanke su (dashed) domin ana kwafar su ne kawai da `--history`. An lissafa `.credentials.json` da `.claude.json` a `share.Never` kuma ba sa taɓa kai wa rukunin: ccam yana tsayawa da kuskure na ciki idan ɗaya ya taɓa bayyana a tsarin (`Build`), kuma zagayen foldar yana tsallake su ta suna.

ccam yana ba kowace hanya yanayi, wanda shi ne abin da kowane layin shiri ke bayyanawa (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: ba ya cikin tushen, ko fayil ne inda ake tsammanin folda. Ana lissafa shi sau ɗaya ga kowane abu a matsayin `not found, skipped`.
- **Linked**: hanyar tushen ta riga tana kai wa kwafin rukunin, misali idan tushen foldar profile ce ta ccam.
- **Same**: rukunin ya riga yana da abu iri ɗaya.
- **New**: rukunin ba shi da shi, don haka ana kwafar sa gaba ɗaya.
- **Merge**: folda a ɓangarorin biyu. ccam yana ƙara fayilolin da rukunin ba shi da su kawai (`copyMissing`); fayilolin da suka bambanta suna zama yadda rukunin yake da su kuma shirin yana cewa `group copies kept`. Ga `plugins/`, ana haɗa rajistocin `installed_plugins.json` da `known_marketplaces.json`, kuma ana sake rubuta cikakkun hanyoyin plugin don su nuna foldar rukunin.
- **Conflict**: fayil a ɓangarorin biyu mai abubuwa daban. Kai ne ke zaɓa ga kowane fayil: `keep group`, `use mine`, da kuma `merge JSON` ga fayilolin `.json` (haɗi mai zurfi inda rukunin ke nasara a ƙimomin da suka yi karo) ko `merge by line` ga `history.jsonl`. Sabar MCP da ma'anarta ta bambanta ita ma rikici ce, da `keep group` ko `use mine`. Zaɓi na farko da aka lissafa shi ne tsoho idan ka danna Enter; `--yes` koyaushe yana riƙe kwafin rukunin.

Jerin zaɓin kansa shi ne `tui.MigratePicker`, wanda aka gina a kan ɓangaren `checklist` da `ccam add` ke amfani da shi; ccam yana yi wa abubuwan da suka ɓace alama da kiran lstat da jerin sabobin MCP (`migrate.Found`).

### Shigar sabobin MCP

Sabobin MCP da ke shiga da OAuth suna ajiye token ɗinsu a fayil ɗin shiga, `.credentials.json` ƙarƙashin `mcpOAuth`, ba a jerin sabobin ba. ccam ba ya kwafar su. Shigar da aka kwafa tana riƙe da refresh token iri ɗaya da na asali, kuma sabar tana juya wannan token a sabuntawa ta farko: daga nan ana fitar da kowane sauran mai riƙe da shi daga wannan sabar, har da `~/.claude` na tushe. Don haka ga kowace sabar da aka kwafa wadda tushen ya shiga, layin `sign-ins` na shirin yana lissafa asusan rukunin da dole su shiga da kansu da `/mcp`, kuma taƙaitawar tana maimaita shi. Shigarwar da babu token a cikinta ba shiga ba ce kuma ba a lissafa ta, haka ma asusun da ke da nasa shigar tuni, ko sabar da rukunin ke riƙe ma'anarta. Ba a rubuta fayil ɗin shiga ko ɗaya.

## Tsaro da yanayi na musamman

- Ana karanta tushen ne kawai. Ba a matsar, haɗa ko canza kowane fayil a cikinsa, kuma ba a taɓa rubuta shigar tushen, har a macOS inda za a iya karanta ta daga Keychain.
- Kafin a maye gurbin fayil na rukuni, ta `use mine`, haɗawa, ko sababbin sabobin MCP a `mcp.json` da ke akwai, ana ajiye kwafin rukunin a `<group>/.ccam-backup/<time>/`. Sunan foldar na musamman ne ga kowace gudanarwa, don haka gudanarwa biyu a daƙiƙa ɗaya ba sa share ajiyar juna. Taƙaitawar tana buga hanyar.
- Da `--yes`, kowane rikici yana riƙe kwafin rukunin, kuma ccam yana faɗin yadda za a zaɓa maimakon haka: gudanar da umarni ɗaya ba tare da `--yes` ba.
- Idan babu abin da zai canza, ccam yana buga `Nothing to copy.`; idan kowane rikici yana riƙe kwafin rukunin, `Nothing to copy: the group keeps its own copies.`
- Idan kwafi ya gaza a tsakiya, ccam yana tsayawa da `The source folder was not changed. Group files replaced so far are saved in` sannan foldar ajiyar. Sake gudanar da shi yana kwafar fayilolin da ba a kwafa su gaba ɗaya ba.
- Idan rukunoni da dama suna raba MCP, umarnan da `ccam doctor`, `ccam ext` da `ccam mcp list` ke ba da shawara suna cewa `ccam migrate --mcp --group <a|b>`.
- Tushen da ya ɓace yana ƙarewa da `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Masu alaƙa

- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Dubawa da gyarawa](/features/doctor)
- [Amfani da asusun Claude Code da dama a kwamfuta ɗaya](/guides/multiple-claude-code-accounts)
- [Bayanin CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
