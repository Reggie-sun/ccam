---
title: Bayanin umarnai
metaTitle: Bayanin umarnan ccam
description: Kowane umarni da zaɓi na ccam, daga ccam add da ccam <name> zuwa ccam usage, ccam next da ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 3342453e38d3
---

Saka `--dry-run` kafin umarni don ganin abin da zai yi ba tare da ya yi shi ba, misali `ccam --dry-run work`. Yana aiki lokacin fara Claude Code (`ccam <name>`, `ccam run`, `ccam next`) da kuma tare da `share`, `migrate`, `shortcut`, `update` da `doctor --fix`. Sauran umarnin da ke canza wani abu suna ƙin shi maimakon su gudana.

<!-- commands -->

## Fara Claude Code

`ccam <name|number> [arguments]` yana fara Claude Code a matsayin wannan asusun. ccam ba ya karanta arguments da ke bayan sunan; yana miƙa su ga `claude` a tsari ɗaya, kuma standard input, exit code da signals suna aiki daidai kamar da `claude` kansa. `--` bayan sunan ba dole ba ne.

`ccam` tare da arguments amma ba tare da sunan asusu ba, kamar `ccam -p "summarise"`, yana amfani da asusun tsoho. Saita shi da `ccam use <name>`.

Gajeren umarni kamar `ccwork` yana fara asusunsa ta hanya ɗaya: `ccwork --resume` yana yin abin da `ccam work --resume` ke yi. Duba [gajerun umarni](/docs/commands#shortcuts).

## Ƙarawa da cire asusai

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` ba tare da folda ba yana neman `~/.claude` da `~/.claude-*` kuma yana tambaya game da kowanne. Ana yin rajistar foldojin da aka shigo da su a inda suke; ccam ba ya taɓa motsawa ko kwafar bayanan shiga.

`ccam rename`, `ccam remove` da `ccam share` suna ƙin canza asusu yayin da Claude Code ke gudana a matsayin wannan asusun. Ka rufe zaman (session) tukuna.

`ccam list --json` yana buga kowane asusu tare da ramin (slot) sa, gajeren umarninsa, imel da foldarsa. Abin `share` ɗinsa yana danganta kowane abu (`mcp`, `settings`, `memory`, `extensions`, `history`) da sunan rukuni ko `local`.

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
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Kowane abu yana karɓar `on` (rukunin `default`), `off`, ko sunan rukuni. Duba [rabawa](/docs/sharing).

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
ccam next --continue
ccam next --auto
```

## Kulawa

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
```

`ccam update` yana duba checksum na SHA-256 na sakin kafin ya maye gurbin ccam, kuma `--version v0.2.0` yana girka takamaiman saki. ccam da aka girka daga kunshin `.deb` ana sabunta shi da apt a maimakon haka. Duba [sabuntawa](/docs/install#update).

## Fayiloli

| Hanya | Abin da ke ciki |
| --- | --- |
| `~/.config/ccam/config.toml` | Asusai da abin da suke rabawa |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` na kowane asusu |
| `~/.local/share/ccam/groups/<group>/` | Saituna, sabobin MCP, extensions da tarihi da aka raba |
| `~/.local/state/ccam/` | Cache na amfani da lokutan amfani na ƙarshe |
| `~/.local/bin/<shortcut>` | Gajerun umarni kamar `ccwork`, mahaɗai zuwa ccam |

Masu canji (variables) `XDG_CONFIG_HOME`, `XDG_DATA_HOME` da `XDG_STATE_HOME` suna canza wurin waɗannan foldojin.
