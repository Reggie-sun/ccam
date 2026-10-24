---
title: Bayanin umarnai
metaTitle: Bayanin umarnan ccam
description: Kowane umarni da zaɓi na ccam, daga ccam add da ccam <name> zuwa ccam usage, ccam next da ccam doctor.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 36af7e209cd9
---

Saka `--dry-run` kafin umarni don ganin abin da zai yi ba tare da ya yi shi ba, misali `ccam --dry-run work`. Yana aiki lokacin fara Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) da kuma tare da `group`, `migrate`, `shortcut`, `update` da `doctor --fix`. Sauran umarnin da ke canza wani abu suna ƙin shi maimakon su gudana.

<!-- commands -->

## Fara Claude Code

`ccam <name|number> [arguments]` yana fara Claude Code a matsayin wannan asusun. ccam ba ya karanta arguments da ke bayan sunan; yana miƙa su ga `claude` a tsari ɗaya, kuma standard input, exit code da signals suna aiki daidai kamar da `claude` kansa. `--` bayan sunan ba dole ba ne.

`ccam` tare da arguments amma ba tare da sunan asusu ba, kamar `ccam -p "summarise"`, yana amfani da asusun da aka saita don foldar yanzu, in ba haka ba asusun tsoho. Saita su da `ccam folder <name> [folder]` da `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` kafin sunan asusun na gudun `-p`/`--print` ne. Idan asusun ya kai iyakar amfaninsa, ccam yana sake gudanar da umarnin iri ɗaya a matsayin asusun da aka shiga da shi mafi sarari, kuma yana buga sakamakon gudun da ya kammala kaɗai. Duba [Asusai bisa folda](/features/folder) da [Sake gudanarwa a iyaka](/features/next#run-again-at-the-limit).

**Kujeru (seats) na Team da Enterprise.** Zaɓin atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba [bayanin sharuɗɗa a cikin Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Gajeren umarni kamar `ccwork` yana fara asusunsa ta hanya ɗaya: `ccwork --resume` yana yin abin da `ccam work --resume` ke yi. Duba [gajerun umarni](/docs/commands#shortcuts).

## Ƙarawa da cire asusai

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` ba tare da folda ba yana neman `~/.claude` da `~/.claude-*` kuma yana tambaya game da kowanne; ana ba da rahoton foldar da ta gaza a matsayin `Skipped <path>: <reason>` kuma ana yin rajistar sauran, kuma foldar da sunanta ya riga ya zama na wani asusu tana buƙatar `--name`, wanda ake bayarwa tare da foldar. Ana yin rajistar foldojin da aka shigo da su a inda suke; ccam ba ya taɓa motsawa ko kwafar bayanan shiga.

`ccam login --expired` yana kuma shiga asusan da Claude Code ya share shigarsu. `ccam remove` yana fitar da asusun da ccam ke sarrafawa daga shiga kafin foldarsa ta tafi shara (ba idan wani asusu yana riƙe da shiga iri ɗaya ba), kuma `ccam remove` na asusun da aka shigo da shi yana fitar da shi daga rukunoninsa tare da kwafin abin da suke rabawa. A macOS, `ccam rename` yana riƙe sunan foldar, domin Keychain yana ajiye shigar ƙarƙashin hanyar foldar.

`ccam add --api` yana ƙara asusun da ke amfani da API key maimakon shigar Claude; ccam yana karanta key ɗin a kowane ƙaddamarwa kuma ba ya taɓa ajiye shi. Duba [Asusan API key](/features/api).

`ccam rename`, `ccam remove` da umarnan `ccam group` suna ƙin canza asusu yayin da Claude Code ke gudana a matsayin wannan asusun, ko ta yaya aka fara shi (har da zaman bayan fage). Ka rufe zaman (session) tukuna. Suna kuma ƙi idan canjin rukunoni bai kammala ba: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` yana buga kowane asusu tare da ramin (slot) sa, gajeren umarninsa, imel da foldarsa. Array ɗinsa na `groups` yana lissafa rukunonin asusun, kuma abin `share` ɗinsa yana danganta kowane abu (`mcp`, `settings`, `memory`, `extensions`, `history`) da rukunin da yake fitowa daga gare shi, ko `local`.

## Gajerun umarni

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Gajeren umarni umarni ne da kai da kanka ka sa masa suna, wanda ke fara Claude Code a matsayin asusu ɗaya. Kamar `ccam work`, yana miƙa kowane argument ga Claude Code ba tare da canji ba, don haka `ccwork --resume` da `ccwork -p "hi"` suna aiki kamar yadda kake tsammani. Mahaɗi ne zuwa ccam a cikin `~/.local/bin`.

Gajerun umarni zaɓi ne. Mataimakin `ccam add` yana tambaya ko a ƙirƙiri ɗaya, `ccam add <name> --shortcut <command>` yana ƙirƙirar ɗaya ba tare da tambaya ba, kuma `ccam list` yana nuna gajeren umarnin kowane asusu, ko `-` idan ba shi da shi.

Sunaye suna amfani da haruffa a-z, 0-9, `.`, `-` da `_`. Suna ba zai iya zama umarnin ccam ba, ko gajeren umarnin wani asusu, ko sunan wani shiri a `PATH` naka, ko kuma fayil a `~/.local/bin` da ba ccam ne ya yi shi ba, kamar kwafinka ko mahaɗinka na ccam. Gajeren umarnin da aka gyara da hannu a `config.toml` wanda ya karya waɗannan ƙa'idodi, ko shi ne `claude` ko `ccam`, ko an ba asusu biyu shi, yana sa `config.toml` ya kasa lodawa, tare da ambaton asusun.

`ccam rename` yana riƙe sunan gajeren umarnin, kuma `ccam remove` yana goge gajeren umarnin tare da asusun. `ccam shims install` yana sake ƙirƙirar mahaɗan dukkan asusai, misali bayan ka matsar da ccam; fayil da ke tare hanyar gajeren umarni ɗaya ba ya hana sauran.

## Rabawa

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Rukuni yana raba abubuwa tsakanin asusansa; asusu zai iya kasancewa a rukunoni da dama, amma kowane abu yana fitowa daga rukuni ɗaya kawai. Duba [rabawa da rukunoni](/features/share).

## Kwafar saitin da ke akwai

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` yana kwafar `settings.json`, `CLAUDE.md`, extensions da sabobin MCP na mai amfani daga foldar Claude Code zuwa rukuni, don kowane asusu da ke raba waɗannan abubuwa ya same su. Ana karanta foldar ne kawai. Zaɓuɓɓukan `--settings`, `--memory`, `--extensions`, `--mcp` da `--history` suna kwafar abubuwan da ka ambata kawai. Ba tare da su ba, a terminal jerin zaɓi (checklist) yana ba ka damar zaɓa, inda aka yi wa komai alama sai tarihin tattaunawa; a wasu wurare kuma ana kwafar komai sai tarihin tattaunawa. `-g` ko `--group` yana zaɓar rukunin, kuma `-y` ko `--yes` yana kwafa ba tare da tambaya ba kuma yana riƙe kwafin rukunin idan fayil ya bambanta. Ba a kwafar shigar sabobin MCP, domin shigar da aka kwafa tana fitar da sauran masu riƙe da ita; tsarin yana ambaton asusan da za su shiga da `/mcp`. `--mcp` yana aiki ko da `~/.claude` babu komai a ciki ko babu shi, amma `~/.claude.json` yana da sabobi. Duba [kwafar saitin da ke akwai](/docs/sharing#copy-an-existing-setup).

## Sabobin MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` yana motsa sabobin da wannan asusun kaɗai ke da su zuwa rukuninsa, har da waɗanda yake da su don folda ɗaya kawai (`claude mcp add` na yau da kullum), waɗanda daga nan suke lodawa a kowace folda. Ana tsallake sabar da rukunin ya riga ya yi amfani da sunanta don wata ma'ana daban; ƙara `--force` don maye gurbin kwafin rukunin, wanda ccam ke fara ajiye masa madadi.

## Amfani da sauyawa

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
```

`next --auto`, `--fallback` da `ccam pick` suna zaɓa ta hanya ɗaya: na farko waɗanda ke ƙasa da `auto_max_usage` (`config.toml`, tsoho 100), bisa `auto_strategy` (`headroom`, `expiring` ko `balance`), suna ƙirga tagar mako ta `--model`, ba sa taɓa zaɓar asusun da ya kai 100% (tagar da Claude Code ya ba da rahoton ta cika ana ƙirga ta a matsayin 100% har zuwa sake saitinta), ba sa taɓa zaɓar shigar API key, kuma suna zama a cikin ƙungiya (organization) ta asusun da suka fara daga gare shi (`--cross-org` ko `cross_org = true` don duba gaba da haka). `--wait` yana ɗaukar ƙimarsa bayan `=` (`--wait=2h`); idan an rubuta shi da sarari ana ƙin shi tare da shawara. Idan ba za a iya zaɓar kowanne ba, suna fita da code 75, ko da `--wait` suna jiran asusun farko ya sake saitawa. `ccam next --continue` yana kai tattaunawa gaba har ma tsakanin asusan da ba sa raba tarihi: a lokacin Claude Code yana kwafar ta cikin sabon asusun. Ba kamar `ccam <name>` ba, `ccam next` yana karanta zaɓuɓɓukansa da kansa, don haka zaɓuɓɓukan Claude Code suna zuwa bayan `--`; ba tare da shi ba ccam yana tsayawa da `unknown flag`. A cikin Claude Code, `/ccam:switch [account]` na plugin ɗin yana miƙa tattaunawa gaba, kuma `ccam hook` shi ne abin da hooks ɗinsa ke gudanarwa. Duba [Sauyawa zuwa asusu na gaba](/features/next), [Zaɓar asusu don rubutun umarni](/features/pick) da [Zaman bayan fage na kowane asusu](/features/agents).

## Kulawa

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` yana duba checksum na SHA-256 na sakin kafin ya maye gurbin ccam, kuma `--version v0.2.0` yana girka takamaiman saki. ccam da aka girka daga kunshin `.deb` ana sabunta shi da apt a maimakon haka. Duba [sabuntawa](/docs/install#update).

## Fayiloli

| Hanya | Abin da ke ciki |
| --- | --- |
| `~/.config/ccam/config.toml` | Asusai, rukunoni da saituna (`theme`, `usage_source`, `relay`, `cross_org` …). Saituna suna zuwa a sama, kafin `[[group]]`, `[[profile]]` ko `[[folder]]` na farko: saitin da ke ƙarƙashin tebur (table) na wannan tebur ne, don haka ccam yana ambatar sa a kowane ƙaddamarwa da a `ccam doctor`, kuma umarnan da ke ajiye config suna ƙi har sai an matsar da shi. Ba a riƙe sharhi (comments) idan ccam ya ajiye fayil ɗin |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` na kowane asusu |
| `~/.local/share/ccam/groups/<group>/` | Saituna, sabobin MCP, extensions da tarihi da aka raba |
| `~/.local/state/ccam/` | Cache na amfani da rikodin iyaka (`usage/<name>.json`, `usage/<name>.limit.json`), lokutan amfani na ƙarshe, alamomin zama da kulle-kulle (`locks/`), gajerun umarnin da ccam ya yi (`shortcuts.json`), canjin rukunoni da bai kammala ba (`group-change.json`), da bayanan da `/ccam:switch` da hook na iyaka ke barwa (`handoff/<account>/<session>.json`) |
| `~/.local/bin/<shortcut>` | Gajerun umarni kamar `ccwork`, mahaɗai zuwa ccam |

Masu canji (variables) `XDG_CONFIG_HOME`, `XDG_DATA_HOME` da `XDG_STATE_HOME` suna canza wurin waɗannan foldojin. A Windows, `~` shi ne `%USERPROFILE%`, kuma gajerun umarni su ne `<shortcut>.exe`.
