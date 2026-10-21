---
title: Bayanin umarnai
metaTitle: Bayanin umarnan ccam
description: Kowane umarni da zaɓi na ccam, daga ccam add da ccam <name> zuwa ccam usage, ccam next da ccam doctor.
updated: 2026-10-21
order: 3
translation: machine
sourceHash: 338ac11a0b5a
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

Gajeren umarni kamar `ccwork` yana fara asusunsa ta hanya ɗaya: `ccwork --resume` yana yin abin da `ccam work --resume` ke yi. Duba [gajerun umarni](/docs/commands#shortcuts).

## Ƙarawa da cire asusai

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` ba tare da folda ba yana neman `~/.claude` da `~/.claude-*` kuma yana tambaya game da kowanne. Ana yin rajistar foldojin da aka shigo da su a inda suke; ccam ba ya taɓa motsawa ko kwafar bayanan shiga.

`ccam add --api` yana ƙara asusun da ke amfani da API key maimakon shigar Claude; ccam yana karanta key ɗin a kowane ƙaddamarwa kuma ba ya taɓa ajiye shi. Duba [Asusan API key](/features/api).

`ccam rename`, `ccam remove` da umarnan `ccam group` suna ƙin canza asusu yayin da Claude Code ke gudana a matsayin wannan asusun. Ka rufe zaman (session) tukuna.

`ccam list --json` yana buga kowane asusu tare da ramin (slot) sa, gajeren umarninsa, imel da foldarsa. Array ɗinsa na `groups` yana lissafa rukunonin asusun, kuma abin `share` ɗinsa yana danganta kowane abu (`mcp`, `settings`, `memory`, `extensions`, `history`) da rukunin da yake fitowa daga gare shi, ko `local`.

## Gajerun umarni

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Gajeren umarni umarni ne da kai da kanka ka sa masa suna, wanda ke fara Claude Code a matsayin asusu ɗaya. Kamar `ccam work`, yana miƙa kowane argument ga Claude Code ba tare da canji ba, don haka `ccwork --resume` da `ccwork -p "hi"` suna aiki kamar yadda kake tsammani. Mahaɗi ne zuwa ccam a cikin `~/.local/bin`.

Gajerun umarni zaɓi ne. Mataimakin `ccam add` yana tambaya ko a ƙirƙiri ɗaya, `ccam add <name> --shortcut <command>` yana ƙirƙirar ɗaya ba tare da tambaya ba, kuma `ccam list` yana nuna gajeren umarnin kowane asusu, ko `-` idan ba shi da shi.

Sunaye suna amfani da haruffa a-z, 0-9, `.`, `-` da `_`. Suna ba zai iya zama umarnin ccam ba, ko gajeren umarnin wani asusu, ko sunan wani shiri a `PATH` naka.

`ccam rename` yana riƙe sunan gajeren umarnin, kuma `ccam remove` yana goge gajeren umarnin tare da asusun. `ccam shims install` yana sake ƙirƙirar mahaɗan dukkan asusai, misali bayan ka matsar da ccam.

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

`ccam migrate` yana kwafar `settings.json`, `CLAUDE.md`, extensions da sabobin MCP na mai amfani daga foldar Claude Code zuwa rukuni, don kowane asusu da ke raba waɗannan abubuwa ya same su. Ana karanta foldar ne kawai. Zaɓuɓɓukan `--settings`, `--memory`, `--extensions`, `--mcp` da `--history` suna kwafar abubuwan da ka ambata kawai. Ba tare da su ba, a terminal jerin zaɓi (checklist) yana ba ka damar zaɓa, inda aka yi wa komai alama sai tarihin tattaunawa; a wasu wurare kuma ana kwafar komai sai tarihin tattaunawa. `-g` ko `--group` yana zaɓar rukunin, kuma `-y` ko `--yes` yana kwafa ba tare da tambaya ba kuma yana riƙe kwafin rukunin idan fayil ya bambanta. Yana kuma kwafar bayanan shigar sabobin da aka ajiye zuwa fayil ɗin shiga na kowane asusu. Duba [kwafar saitin da ke akwai](/docs/sharing#copy-an-existing-setup).

## Sabobin MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` yana motsa sabobin da wannan asusun kaɗai ke da su zuwa rukuninsa. Ana tsallake sabar da rukunin ya riga ya yi amfani da sunanta don wata ma'ana daban; ƙara `--force` don maye gurbin kwafin rukunin, wanda ccam ke fara ajiye masa madadi.

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
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` da `ccam pick` suna zaɓa ta hanya ɗaya: na farko waɗanda ke ƙasa da `auto_max_usage` (`config.toml`, tsoho 100), bisa `auto_strategy` (`headroom`, `expiring` ko `balance`), suna ƙirga tagar mako ta `--model`, ba sa taɓa zaɓar asusun da ya kai 100%, kuma suna zama a cikin ƙungiya (organization) ta asusun da suka fara daga gare shi (`--cross-org` ko `cross_org = true` don duba gaba da haka). Idan ba za a iya zaɓar kowanne ba, suna fita da code 75, ko da `--wait` suna jiran asusun farko ya sake saitawa. `ccam next --continue` yana kai tattaunawa gaba har ma tsakanin asusan da ba sa raba tarihi: a lokacin Claude Code yana kwafar ta cikin sabon asusun. Ba kamar `ccam <name>` ba, `ccam next` yana karanta zaɓuɓɓukansa da kansa, don haka zaɓuɓɓukan Claude Code suna zuwa bayan `--`; ba tare da shi ba ccam yana tsayawa da `unknown flag`. A cikin Claude Code, `/ccam:switch [account]` na plugin ɗin yana miƙa tattaunawa gaba, kuma `ccam hook` shi ne abin da hooks ɗinsa ke gudanarwa. Duba [Sauyawa zuwa asusu na gaba](/features/next), [Zaɓar asusu don rubutun umarni](/features/pick) da [Zaman bayan fage na kowane asusu](/features/agents).

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
| `~/.config/ccam/config.toml` | Asusai, rukunoni da saituna (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` na kowane asusu |
| `~/.local/share/ccam/groups/<group>/` | Saituna, sabobin MCP, extensions da tarihi da aka raba |
| `~/.local/state/ccam/` | Cache na amfani, lokutan amfani na ƙarshe da bayanan da `/ccam:switch` da hook na iyaka ke barwa (`handoff/`) |
| `~/.local/bin/<shortcut>` | Gajerun umarni kamar `ccwork`, mahaɗai zuwa ccam |

Masu canji (variables) `XDG_CONFIG_HOME`, `XDG_DATA_HOME` da `XDG_STATE_HOME` suna canza wurin waɗannan foldojin. A Windows, `~` shi ne `%USERPROFILE%`, kuma gajerun umarni su ne `<shortcut>.exe`.
