---
title: Kusimamia akaunti
description: Ongeza, ingia, orodhesha, badilisha jina, ondoa na ingiza akaunti za Claude Code kwa ccam add, ccam login, ccam list, ccam rename, ccam remove na ccam import.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "8564517d0d71"
---

Kila akaunti ya ccam ni folda tofauti ya usanidi wa Claude Code yenye taarifa zake za kuingia. Amri hizi huunda, huingiza, huorodhesha, hubadilisha majina na huondoa akaunti, na husajili folda ambazo tayari unazo.

## Matumizi

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

Popote amri inapopokea `<name|number>`, unaweza kutoa jina la akaunti au namba yake (slot), yaani safu ya `#` ya `ccam list`.

`ccam add` kwenye terminal bila chaguo hufungua mchawi. Huuliza jina na jinsi akaunti inavyoingia (kuingia kwa Claude sasa, API key, au kuingia kwa Claude baadaye), kisha huonyesha vikundi vya kujiunga huku `default` ikiwa imeteuliwa (kabla ya kuwepo kwa vikundi, vipengele ambavyo kikundi `default` kitashiriki), kisha huuliza `Create a shortcut command for work? [y/N]`, na humalizia kwa muhtasari wa kuthibitisha. `Esc` hurudi hatua moja nyuma. Chaguo lolote kati ya `--group`, `--no-group`, `--shortcut` au `--no-login`, au pipe badala ya terminal, huruka mchawi; hapo jina linahitajika na hakuna njia fupi inayoundwa isipokuwa `--shortcut` itaje moja.

Bila `--group`, akaunti mpya hujiunga na kikundi `default`. Akaunti ya kwanza hukiunda, kikishiriki seva za MCP, mipangilio, `CLAUDE.md` na viendelezi; tazama [kushiriki kwa vikundi](/features/share).

`ccam list` huonyesha namba, jina, njia fupi, barua pepe, vikundi na matumizi ya mwisho, na `*` mbele ya akaunti chaguo-msingi:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

Kwa `--json`, kila akaunti ni kitu chenye `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (huachwa ikiwa akaunti haijawahi kuanzishwa), `groups` na `share`. `share` huhusisha kila kipengele (`mcp`, `settings`, `memory`, `extensions`, `history`) na kikundi kinachotoka, au `"local"`.

## Chaguo

| Chaguo | Kinachofanya | Chaguo-msingi |
| --- | --- | --- |
| `ccam add --no-login` | Unda akaunti bila kuingia; ingia baadaye kwa `ccam login <name>` | imezimwa |
| `ccam add --group <group>` | Jiunge na kikundi hiki badala ya `default`. Kinaweza kurudiwa; `default` huundwa ikiwa hakipo | `default` |
| `ccam add --no-group` | Usijiunge na kikundi chochote: akaunti haishiriki chochote | imezimwa |
| `ccam add --shortcut <command>` | Pia unda amri hii ya njia fupi katika `~/.local/bin` | hakuna |
| `ccam add --api --key-from <source>` | Akaunti inayotumia API key badala ya kuingia; tazama [Akaunti za API key](/features/api) | imezimwa |
| `ccam list --json` | Chapisha akaunti kama JSON | imezimwa |
| `ccam remove -y`, `--yes` | Ondoa bila kuuliza | imezimwa |
| `ccam import --name <name>` | Jina la akaunti kwa folda unayotoa | hutokana na jina la folda |
| `ccam import -y`, `--yes` | Ingiza kila folda iliyopatikana bila kuuliza | imezimwa |
| `--dry-run` | Hukubaliwa na `ccam list` na `ccam completion`, ambazo hazibadilishi kitu. Hukataliwa na `add`, `login`, `logout`, `use`, `rename`, `remove` na `import`, kwa ujumbe `` --dry-run is not supported by `ccam add`; nothing was changed `` (ukitaja amri husika) | imezimwa |

`ccam login`, `ccam logout`, `ccam use` na `ccam rename` hazina chaguo zao wenyewe. `ccam completion` hupokea `bash`, `zsh` au `powershell`; kitu kingine chochote hushindwa kwa `only bash, zsh and powershell are supported`.

Majina ya akaunti yana herufi 1 hadi 32 za `a-z`, `0-9`, `-` na `_`, yakianza na herufi au tarakimu. Jina haliwezi kuwa tarakimu tu, kwa sababu `ccam 2` inamaanisha namba 2. Jina pia haliwezi kuwa amri ya ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`), ili `ccam <name>` isiwe na utata kamwe. Orodha hiyo ni `config.Reserved` katika `internal/config`.

## Jinsi inavyofanya kazi

`ccam add work --shortcut ccwork` hutekeleza hatua hizi katika `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` hukagua jina, na `shim.Check` hukagua njia fupi kabla ya kitu chochote kuundwa.
2. `profile.Add` huunda `~/.local/share/ccam/profiles/work/` kwa mode `0700` na huipa akaunti namba ndogo kabisa iliyo wazi (`Config.NextSlot`), kwa hiyo namba ya akaunti iliyoondolewa hutumiwa tena. Akaunti ya kwanza huwa chaguo-msingi.
3. Akaunti hujiunga na vikundi vyake (`default` isipokuwa `--group` au `--no-group` iseme vinginevyo; akaunti ya kwanza huunda `default`), na `applyShare` huunganisha folda mpya iliyo tupu na vikundi hivyo kwa `share.Build` na kisha `Apply`, kwa mfano `settings.json` kwenda `groups/default/settings.json`. Hili likishindwa, folda mpya huondolewa, kwa sababu folda iliyobaki ingezuia kila jaribio jipya kwa `folder already exists`.
4. Akaunti huhifadhiwa katika `~/.config/ccam/config.toml` kupitia `fsutil.AtomicWrite`: faili ya muda, `fsync`, kisha `rename`, ili hitilafu ya ghafla isiache usanidi nusu.
5. `runLogin` huanzisha `claude auth login` huku `CLAUDE_CONFIG_DIR` ikiwa imewekwa kwenye folda ya akaunti, pamoja na `--email` yenye barua pepe ya mwisho ya akaunti ili ukurasa wa kuingia ujazwe tayari. Kuingia kukishindwa, akaunti hubaki, na ccam huchapisha hitilafu na kukuambia uendeshe `ccam login work` baadaye.
6. `profile.ReadIdentity` huchukua `oauthAccount` kutoka `profiles/work/.claude.json`: barua pepe, id za akaunti na shirika, jina la shirika na mpango (`claude_max`, `claude_team` …). Faili hiyo inaweza kuwa megabaiti kadhaa, kwa hiyo `TopLevelValue` hutafuta ufunguo huo mmoja badala ya kuchanganua faili nzima. ccam huvihifadhi katika `[profile.login]` ndani ya `config.toml`.
7. `shim.Install` huunganisha `~/.local/bin/ccwork` na binary ya ccam, ambayo hutafuta `argv[0]` katika `config.toml` inapoanza.

Kuingia kunaweza kuchukua dakika kadhaa, na ccam nyingine inaweza kuhifadhi usanidi wakati huo. Kwa hiyo amri zinazosubiri huhifadhi kupitia `env.update`, ambayo hupakia upya `config.toml`, hutumia badiliko hilo moja (hapa, barua pepe) na kuhifadhi. `ccam remove` hupakia upya baada ya swali lake kwa sababu hiyo hiyo.

`ccam login` hutekeleza hatua ya 5 na 6 kwa akaunti iliyopo. Claude Code huichukua barua pepe kama kidokezo tu kwa ukurasa wa kuingia, kwa hiyo kuingia kunaweza kuishia kwenye akaunti nyingine ya Claude kuliko ile ambayo akaunti ya ccam ilikuwa; ccam hulinganisha kuingia kupya na `[profile.login]` na husema hivyo (`ccam: this login is other@company.com, but work was dev@company.com before.`). Huhifadhi kuingia kupya na hukuambia jinsi ya kulirudisha: `ccam logout work`, kisha `ccam login work` na uingie kwa barua pepe sahihi. Akaunti nyingine ya ccam ikiwa imeingia kwenye akaunti ile ile ya Claude, ccam huonya kwamba hizo mbili hushiriki vikomo vyao vya matumizi. `ccam login --expired` huingiza, moja baada ya nyingine, kila akaunti ambayo haijawahi kuingia au ambayo kuingia kwake (refresh token yake) kumeisha muda.

`ccam logout` huendesha `claude auth logout` ndani ya akaunti na hufuta barua pepe iliyohifadhiwa; rekodi ya kuingia hubaki, kwa hiyo kuingia kunakofuata hukaguliwa dhidi yake. `ccam use` huandika upya `default` tu katika `config.toml`.

`ccam list` huongeza safu ya `Organization` akaunti zinapokuwa za mashirika zaidi ya moja (mipango ya binafsi huonyesha `personal`), na `ccam list --json` ina `plan` na `org` kwa kila akaunti.

### Kuingia kwenye mashine ya mbali

Kupitia SSH, sehemu ya kivinjari ya `claude auth login` inaweza kushindwa. Claude Code husubiri kuingia kwenye seva ya callback katika `http://localhost:<port>/callback` kwenye seva, lakini kivinjari kinachokamilisha kuingia kinaendeshwa kwenye kompyuta yako, ambako anwani hiyo haielekei popote na kivinjari huonyesha hitilafu kama `ERR_CONNECTION_REFUSED`. Anwani hiyo bado ina msimbo wa mara moja, kwa hiyo ccam inaweza kuuwasilisha:

<!-- motion:accounts-relay -->

`runLogin` haikabidhi terminal yako kwa Claude Code. Huunganisha standard input ya Claude Code na pipe, na `relayLogin` hunakili ndani yake kila mstari unaouandika, kwa hiyo msimbo unaobandikwa kwenye swali la `Paste code here` la Claude Code bado hufika. Mstari unapokuwa anwani ya callback, `callbackURL` huikubali tu ikiwa host ni `localhost`, `127.0.0.1` au `[::1]`, port imetolewa, njia ni `/callback` na `code` pamoja na `state` zote zipo; `http://` inaweza kuachwa. Kisha `deliverCallback` huomba anwani hiyo kwenye seva kwa kikomo cha sekunde 10 na bila kufuata redirects, kama kivinjari kingefanya. Ombi hufika kwenye seva ya callback ya `claude auth login` inayosubiri (mstari wa vistari kwenye mchoro). Claude Code hupokea msimbo na huhifadhi taarifa za kuingia, kwenye Linux katika `.credentials.json` ndani ya folda ya akaunti. Mstari uliobandikwa hauwezi kamwe kuifanya ccam iwasiliane na host nyingine.

ccam huchapisha hiki kabla kuingia hakujaanza:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Ikifaulu huchapisha `→ Passed the sign-in to Claude Code on this machine.`

Wakati kuingia kunasubiri ingizo kwenye terminal, ccam huzima mwangwi (echo) na hali ya mstari ya terminal na huonyesha `*` kwa kila herufi unayoandika, kwa hiyo msimbo au anwani ya callback iliyobandikwa haionekani kwenye skrini. `readLine` hufanya uhariri wa mstari yenyewe: `Enter` hutuma mstari, `Backspace` hufuta herufi moja, na `Ctrl+D` kwenye mstari mtupu humaliza ingizo. Terminal hurejeshwa kuingia kunapoisha, na pia kwa `Ctrl+C`, `SIGTERM` au `SIGHUP`, ambazo bado humaliza ccam. Bila terminal kwenye standard input na standard error, mistari hunakiliwa kama inavyokuja.

### Kubadilisha jina, kuondoa na kuingiza

`ccam rename` na `ccam remove` hufanya kazi kwenye folda ya akaunti na kwenye kile ambacho ccam huhifadhi mahali pengine kwa jina la akaunti. `nameState` huorodhesha njia hizo: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` na cache za matumizi `~/.local/state/ccam/usage/<name>.json` na `<name>.statusline.json`. `ccam remove work` huzishughulikia kwa mpangilio huu:

<!-- motion:accounts-remove -->

`ccam remove` kwanza hukataa ikiwa akaunti inaendeshwa, kisha huuliza. Ikiwa akaunti inashiriki viendelezi, `RewritePluginPaths` huandika upya `installed_plugins.json` na `known_marketplaces.json` za kikundi, kwa sababu Claude Code huhifadhi njia kamili za plugins ambazo zinaweza kupita kwenye folda ya akaunti hii. `profile.Trash` huhamisha folda kwenda `~/.local/share/ccam/trash/<name>-<time>/`, akaunti huondoka kwenye `config.toml`, njia za `nameState` hufutwa ili akaunti ya baadaye yenye jina lile lile isizirithi, na kiungo cha njia fupi huondolewa. Ili kurejesha, hamisha folda kutoka kwenye trash.

`ccam rename` hubadilisha jina la folda, huandika upya njia za plugins katika `plugins/` yake kutoka folda ya zamani kwenda mpya, huhamisha hali ya chaguo-msingi pamoja nayo, na huhamisha njia za `nameState` kwenda jina jipya. Amri ya njia fupi huhifadhi jina lake. Uandikaji upya wa plugins ukishindwa, folda hurudishiwa jina lake la zamani, kwa kuwa usanidi haujahifadhiwa.

`ccam import` husajili folda pale ilipo, kama `dir` katika `config.toml`. Kamwe haiihamishi wala kuinakili: refresh token iliyonakiliwa ingesababisha nakala hizo zitoane nje ya akaunti. Bila hoja ya folda, hutafuta folda za `~/.claude` na `~/.claude-*` ambazo hazijasajiliwa na huuliza kuhusu kila moja; majina hutokana na folda (`~/.claude` ni `main`, `~/.claude-work` ni `work`, `~/.claude-2` ni `acct2`), na folda ambayo jina lake haliwezi kupatikana huorodheshwa pamoja na amri ya `--name` ya kutumia. Bila `CLAUDE_CONFIG_DIR`, Claude Code huhifadhi `.claude.json` katika folda yako ya nyumbani kando ya `~/.claude`, kwa hiyo unapoingiza `~/.claude` na haina `.claude.json` yake yenyewe, ccam huunda kiungo `~/.claude/.claude.json` kwenda `~/.claude.json`. Kisha `claude` ya kawaida na `ccam main` huona akaunti ile ile. Akaunti zilizoingizwa haziko katika kikundi chochote hadi uziongeze kwenye kimoja (`ccam group join <group> <name>`), na hakuna njia fupi inayoundwa.

`ccam completion bash|zsh` huchapisha script ya kukamilisha amri ambayo pia hukamilisha majina ya akaunti.

## Usalama na hali maalum

- `ccam rename` na `ccam remove` hukataa wakati akaunti inaendeshwa (`refuseRunning`), kwa `work is running; quit its Claude Code session first`, kwa sababu kipindi kinaendelea kuandika ndani ya folda.
- `ccam remove` kwenye akaunti iliyoingizwa huifuta tu kwenye usajili: `✓ main unregistered (folder ... left in place)`. `~/.claude` yako haihamishwi kwenda kwenye trash.
- `ccam add` hukataa jina ambalo folda yake tayari ipo katika `profiles/`, na huelekeza kwa `ccam import`.
- `ccam import` hukataa folda ambayo tayari imesajiliwa, njia ambayo si folda, na jina lililokwisha kuchukuliwa.
- `ccam logout` na `ccam login` hazikagui kama akaunti inaendeshwa.
- Ikiwa relay haiwezi kufikia Claude Code, kwa mfano baada ya kuingia kumalizika, huchapisha `ccam: could not pass the sign-in to Claude Code:` pamoja na sababu. Hitilafu ya HTTP kwa kawaida inamaanisha msimbo umeisha muda; anza kuingia tena.

## Zinazohusiana

- [Kuanzisha Claude Code](/features/launch)
- [Amri za njia fupi](/features/shortcuts)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Akaunti kadhaa za Claude Code](/guides/multiple-claude-code-accounts)
