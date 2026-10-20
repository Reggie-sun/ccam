---
title: Sarrafa asusai
description: Ƙara, shiga, lissafa, sake suna, cire da shigo da asusun Claude Code da ccam add, ccam login, ccam list, ccam rename, ccam remove da ccam import.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

Kowane asusun ccam foldar saitin Claude Code ce ta daban mai nata shiga. Waɗannan umarnai suna ƙirƙira, shiga, lissafa, sake suna da cire asusai, kuma suna yin rajistar foldojin da kake da su tuni.

## Amfani

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Duk inda umarni ke karɓar `<name|number>`, za ka iya ba da sunan asusun ko lambar ramin (slot) sa, wato ginshiƙin `#` na `ccam list`.

`ccam add` a terminal ba tare da zaɓuɓɓuka ba yana buɗe mataimaki (wizard). Yana tambayar suna da yadda asusun zai shiga (shigar Claude yanzu, API key, ko shigar Claude daga baya), sannan yana nuna rukunonin da za a shiga tare da alama a kan `default` (kafin a sami rukunoni, abubuwan da rukunin `default` zai raba), sannan yana tambaya `Create a shortcut command for work? [y/N]`, kuma yana ƙarewa da taƙaitawa don ka tabbatar. `Esc` yana komawa mataki ɗaya baya. Duk wani daga `--group`, `--no-group`, `--shortcut` ko `--no-login`, ko bututu (pipe) maimakon terminal, yana tsallake mataimakin; to sai dole a ba da suna, kuma ba a ƙirƙirar gajeren umarni sai idan `--shortcut` ya ambaci ɗaya.

Ba tare da `--group` ba, sabon asusu yana shiga rukunin `default`. Asusun farko ne yake ƙirƙirar sa, yana raba sabobin MCP, saituna, `CLAUDE.md` da extensions; duba [rabawa da rukunoni](/features/share).

`ccam list` yana nuna ramin, suna, gajeren umarni, imel, rukunoni da amfani na ƙarshe, tare da `*` a gaban asusun tsoho:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

Da `--json`, kowane asusu abu ne mai `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (ba ya nan idan ba a taɓa fara asusun ba), `groups` da `share`. `share` yana danganta kowane abu (`mcp`, `settings`, `memory`, `extensions`, `history`) da rukunin da yake fitowa daga gare shi, ko da `"local"`.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `ccam add --no-login` | Ƙirƙiri asusun ba tare da shiga ba; shiga daga baya da `ccam login <name>` | a kashe |
| `ccam add --group <group>` | Shiga wannan rukunin maimakon `default`. Ana iya maimaitawa; ana ƙirƙirar `default` idan babu shi | `default` |
| `ccam add --no-group` | Kada a shiga kowane rukuni: asusun ba ya raba komai | a kashe |
| `ccam add --shortcut <command>` | Kuma ƙirƙiri wannan gajeren umarni a `~/.local/bin` | babu |
| `ccam add --api --key-from <source>` | Asusun da ke amfani da API key maimakon shiga; duba [Asusan API key](/features/api) | a kashe |
| `ccam list --json` | Buga asusan a matsayin JSON | a kashe |
| `ccam remove -y`, `--yes` | Cire ba tare da tambaya ba | a kashe |
| `ccam import --name <name>` | Sunan asusu na foldar da ka bayar | daga sunan foldar |
| `ccam import -y`, `--yes` | Shigo da kowace foldar da aka samo ba tare da tambaya ba | a kashe |
| `--dry-run` | `ccam list` da `ccam completion` suna karɓar shi, domin ba sa canza komai. `add`, `login`, `logout`, `use`, `rename`, `remove` da `import` suna ƙin shi, da `` --dry-run is not supported by `ccam add`; nothing was changed `` (tare da sunan umarnin) | a kashe |

`ccam login`, `ccam logout`, `ccam use` da `ccam rename` ba su da zaɓuɓɓukan kansu. `ccam completion` yana karɓar `bash`, `zsh` ko `powershell`; duk wani abu dabam yana gazawa da `only bash, zsh and powershell are supported`.

Sunayen asusai haruffa 1 zuwa 32 ne na `a-z`, `0-9`, `-` da `_`, suna farawa da harafi ko lamba. Suna ba zai iya zama lambobi kawai ba, domin `ccam 2` yana nufin ramin 2. Haka kuma suna ba zai iya zama umarnin ccam ba (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`), don haka `ccam <name>` ba ya taɓa yin shubuha. Jerin shi ne `config.Reserved` a `internal/config`.

## Yadda yake aiki

`ccam add work --shortcut ccwork` yana gudanar da waɗannan matakai a `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` yana duba sunan, kuma `shim.Check` yana duba gajeren umarnin kafin a ƙirƙiri komai.
2. `profile.Add` yana ƙirƙirar `~/.local/share/ccam/profiles/work/` da mode `0700` kuma yana ba asusun mafi ƙanƙantar lambar rami da ke a kwance (`Config.NextSlot`), don haka ana sake amfani da lambar asusun da aka cire. Asusun farko yana zama tsoho.
3. Asusun yana shiga rukunoninsa (`default` sai dai idan `--group` ko `--no-group` ya faɗi wani abu dabam; asusun farko yana ƙirƙirar `default`), kuma `applyShare` yana haɗa sabuwar foldar mara komai da su ta `share.Build` sannan `Apply`, misali `settings.json` zuwa `groups/default/settings.json`. Idan wannan ya gaza, ana cire sabuwar foldar, domin foldar da ta rage za ta hana kowace sake gwadawa da `folder already exists`.
4. Ana ajiye asusun a `~/.config/ccam/config.toml` ta `fsutil.AtomicWrite`: fayil na wucin gadi, `fsync`, sannan `rename`, don haka faɗuwar shiri ba ta taɓa barin rabin saiti ba.
5. `runLogin` yana fara `claude auth login` tare da `CLAUDE_CONFIG_DIR` da aka saita zuwa foldar asusun, da `--email` tare da imel na ƙarshe na asusun don a cike shafin shiga. Idan shigar ta gaza, asusun yana nan, kuma ccam yana buga kuskuren kuma yana gaya maka ka gudanar da `ccam login work` daga baya.
6. `profile.ReadIdentity` yana fitar da `oauthAccount` daga `profiles/work/.claude.json`: imel, id na asusu da na ƙungiya (organization), sunan ƙungiyar da tsarin biyan kuɗi (`claude_max`, `claude_team` …). Wannan fayil zai iya kai megabytes da dama, don haka `TopLevelValue` yana neman maɓalli ɗaya kawai maimakon ya karanta dukkan fayil ɗin. ccam yana ajiye su a `[profile.login]` a cikin `config.toml`.
7. `shim.Install` yana haɗa `~/.local/bin/ccwork` da binary na ccam, wanda ke neman `argv[0]` a `config.toml` lokacin da ya fara.

Shiga na iya ɗaukar mintuna, kuma wani ccam na iya ajiye saitin a wannan lokacin. Don haka umarnin da ke jira suna ajiyewa ta `env.update`, wanda ke sake loda `config.toml`, yana amfani da canji ɗaya (a nan, imel ɗin) kuma yana ajiyewa. `ccam remove` yana sake lodawa bayan tambayarsa saboda dalili ɗaya.

`ccam login` yana gudanar da mataki na 5 da na 6 ga asusun da ke akwai. Claude Code yana ɗaukar imel ɗin ne kawai a matsayin alama ga shafin shiga, don haka shiga zai iya ƙarewa a wani asusun Claude dabam da wanda asusun ccam yake a da; ccam yana kwatanta sabuwar shigar da `[profile.login]` kuma yana faɗin haka (`ccam: this login is other@company.com, but work was dev@company.com before.`). Yana riƙe sabuwar shigar kuma yana gaya maka yadda za ka soke ta: `ccam logout work`, sannan `ccam login work`, kuma ka shiga da imel ɗin da ya dace. Idan wani asusun ccam ya shiga asusun Claude iri ɗaya, ccam yana gargaɗi cewa su biyun suna raba iyakokin amfaninsu. `ccam login --expired` yana shiga, ɗaya bayan ɗaya, kowane asusu da ba a taɓa shiga da shi ba ko wanda shigarsa (refresh token ɗinsa) ta ƙare.

`ccam logout` yana gudanar da `claude auth logout` a cikin asusun kuma yana share imel ɗin da aka ajiye; bayanin shigar yana nan, don haka ana duba shiga ta gaba da shi. `ccam use` yana sake rubuta `default` a `config.toml` kawai.

`ccam list` yana ƙara ginshiƙin `Organization` idan asusan na ƙungiyoyi fiye da ɗaya ne (tsare-tsaren kai suna nuna `personal`), kuma `ccam list --json` yana da `plan` da `org` ga kowane asusu.

### Shiga a kwamfutar nesa

Ta SSH, ɓangaren burauza na `claude auth login` zai iya lalacewa. Claude Code yana jiran shigar a sabar callback a `http://localhost:<port>/callback` a kan sabar, amma burauzar da ke kammala shigar tana aiki a kwamfutarka, inda wannan adireshin ba ya kai ko'ina kuma burauzar tana nuna kuskure kamar `ERR_CONNECTION_REFUSED`. Adireshin har yanzu yana ɗauke da lambar sau ɗaya, don haka ccam zai iya isar da ita:

<!-- motion:accounts-relay -->

`runLogin` ba ya miƙa terminal ɗinka ga Claude Code. Yana haɗa standard input na Claude Code da bututu (pipe), kuma `relayLogin` yana kwafar kowane layi da ka rubuta zuwa cikinsa, don haka lambar da aka liƙa a tambayar `Paste code here` ta Claude Code har yanzu tana isa. Idan layi adireshin callback ne, `callbackURL` yana karɓar sa ne kawai idan host ɗin `localhost`, `127.0.0.1` ko `[::1]` ne, an ba da port, hanyar `/callback` ce kuma duka `code` da `state` suna nan; ana iya barin `http://`. Sannan `deliverCallback` yana neman wannan adireshin a kan sabar da iyakar daƙiƙa 10 kuma ba tare da bin redirects ba, kamar yadda burauzar za ta yi. Buƙatar tana isa sabar callback na `claude auth login` da ke jira (layin da aka yanke a zanen). Claude Code yana karɓar lambar kuma yana ajiye shigar, a Linux a `.credentials.json` cikin foldar asusun. Layin da aka liƙa ba zai taɓa sa ccam ya tuntuɓi wani host ba.

ccam yana buga wannan kafin shigar ta fara:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Idan ya yi nasara yana buga `→ Passed the sign-in to Claude Code on this machine.`

Yayin da shigar ke jiran shigarwa a terminal, ccam yana kashe echo da yanayin layi (line mode) na terminal kuma yana nuna `*` ga kowane harafi da ka rubuta, don haka lambar da aka liƙa ko adireshin callback ba ya bayyana a allo. `readLine` yana gyara layin da kansa: `Enter` yana aika layin, `Backspace` yana goge harafi ɗaya, kuma `Ctrl+D` a layi mara komai yana kawo ƙarshen shigarwa. Ana mayar da terminal yadda yake idan shigar ta ƙare, haka kuma a `Ctrl+C`, `SIGTERM` ko `SIGHUP`, waɗanda har yanzu suke ƙare ccam. Idan babu terminal a standard input da standard error, ana kwafar layuka yadda suke zuwa.

### Sake suna, cirewa da shigowa

`ccam rename` da `ccam remove` suna aiki a kan foldar asusun da kuma abin da ccam ke ajiyewa da sunan asusun a wasu wurare. `nameState` yana lissafa waɗannan hanyoyi: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` da cache na amfani `~/.local/state/ccam/usage/<name>.json` da `<name>.statusline.json`. `ccam remove work` yana sarrafa su a wannan tsari:

<!-- motion:accounts-remove -->

`ccam remove` da farko yana ƙi idan asusun yana aiki, sannan yana tambaya. Idan asusun yana raba extensions, `RewritePluginPaths` yana sake rubuta `installed_plugins.json` da `known_marketplaces.json` na rukunin, domin Claude Code yana ajiye cikakkun hanyoyin plugin waɗanda za su iya ratsa foldar wannan asusun. `profile.Trash` yana matsar da foldar zuwa `~/.local/share/ccam/trash/<name>-<time>/`, asusun yana fita daga `config.toml`, ana goge hanyoyin `nameState` don kada asusun da za a yi daga baya mai suna iri ɗaya ya gaje su, kuma ana cire mahaɗin gajeren umarnin. Don maidowa, mayar da foldar daga shara.

`ccam rename` yana sake sunan foldar, yana sake rubuta hanyoyin plugin a cikin `plugins/` ɗinta daga tsohuwar foldar zuwa sabuwa, yana ɗauke da matsayin tsoho tare, kuma yana matsar da hanyoyin `nameState` zuwa sabon sunan. Gajeren umarnin yana riƙe sunansa. Idan sake rubuta plugin ɗin ya gaza, ana mayar da sunan foldar, tunda ba a ajiye saitin ba.

`ccam import` yana yin rajistar folda a inda take, a matsayin `dir` a `config.toml`. Ba ya taɓa matsar ko kwafar ta: refresh token da aka kwafa zai sa kwafin su fitar da juna. Ba tare da argument na folda ba yana neman foldojin `~/.claude` da `~/.claude-*` da ba a yi wa rajista ba kuma yana tambaya game da kowace; sunaye suna fitowa daga foldar (`~/.claude` shi ne `main`, `~/.claude-work` shi ne `work`, `~/.claude-2` shi ne `acct2`), kuma ana lissafa foldar da ba za a iya samun suna daga gare ta ba tare da umarnin `--name` da za a yi amfani da shi. Ba tare da `CLAUDE_CONFIG_DIR` ba, Claude Code yana ajiye `.claude.json` a foldar gidanka kusa da `~/.claude`, don haka idan ka shigo da `~/.claude` kuma ba shi da `.claude.json` nasa, ccam yana ƙirƙirar mahaɗin `~/.claude/.claude.json` zuwa `~/.claude.json`. Sannan `claude` na yau da kullum da `ccam main` suna ganin asusu ɗaya. Asusan da aka shigo da su ba sa cikin kowane rukuni har sai ka ƙara su a ɗaya (`ccam group join <group> <name>`), kuma ba a ƙirƙirar gajeren umarni.

`ccam completion bash|zsh` yana buga rubutun kammalawa (completion) wanda kuma ke kammala sunayen asusai.

## Tsaro da yanayi na musamman

- `ccam rename` da `ccam remove` suna ƙi yayin da asusun ke aiki (`refuseRunning`), da `work is running; quit its Claude Code session first`, domin zaman yana ci gaba da rubutu a cikin foldar.
- `ccam remove` a kan asusun da aka shigo da shi yana soke rajistarsa kawai: `✓ main unregistered (folder ... left in place)`. Ba a matsar da `~/.claude` naka zuwa shara ba.
- `ccam add` yana ƙin sunan da foldarsa ta riga ta kasance a `profiles/`, kuma yana nuna `ccam import`.
- `ccam import` yana ƙin foldar da aka riga aka yi wa rajista, hanyar da ba folda ba ce, da sunan da aka riga aka ɗauka.
- `ccam logout` da `ccam login` ba sa duba ko asusun yana aiki.
- Idan relay ɗin ba zai iya kai wa Claude Code ba, misali bayan shigar ta ƙare, yana buga `ccam: could not pass the sign-in to Claude Code:` tare da dalilin. Kuskuren HTTP yawanci yana nufin lambar ta ƙare; sake fara shigar.

## Masu alaƙa

- [Fara Claude Code](/features/launch)
- [Gajerun umarni](/features/shortcuts)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Asusun Claude Code da dama](/guides/multiple-claude-code-accounts)
