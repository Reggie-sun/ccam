---
title: Dubawa da gyarawa
description: Duba mahaɗan rabawa, izinin fayiloli, shiga da sigar Claude Code da ccam doctor, kuma gyara abin da ke da aminci a gyara da ccam doctor --fix.
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` yana duba girkin ccam ɗinka: binary na Claude Code, foldar kowane asusu, shiga da mahaɗan rabawa, sabobin MCP da aka raba, endpoint na amfani da gajerun umarninka. `ccam doctor --fix` yana gyara abin da zai iya yi ba tare da rasa bayanai ba, kuma yana fara matsar da duk fayil da zai maye gurbinsa zuwa foldar ajiya. Gudanar da shi idan `ccam <name>` ya tsaya saboda matsalar mahaɗin rabawa, bayan ka maido da fayiloli da hannu, ko bayan sabunta Claude Code.

## Amfani

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Kowane sakamako layi ɗaya ne mai alama: `✓` lafiya, `·` bayani, `!` gargaɗi, `✗` kuskure. Sakamakon asusu yana farawa da sunan asusun.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Da `--fix`, kowane sakamako da aka gyara yana ƙarewa da `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--fix` | Gyara abin da ke da aminci a gyara (ana ajiye kwafin fayilolin da aka canza) | a kashe |
| `--dry-run` (na gaba ɗaya) | Gudanar da kowane dubawa amma kada a gyara komai, ko da `--fix` | a kashe |

`ccam doctor` ba ya karɓar arguments. Yana karɓar `--dry-run`: `ccam --dry-run doctor --fix` yana buga sakamako iri ɗaya da `ccam doctor` kuma ba ya canza komai.

Exit code shi ne `1` idan aƙalla sakamakon `✗` ɗaya ya rage ba a gyara ba, kuma `0` in ba haka ba. Gargaɗi da bayani ba sa taɓa sa gudanarwar ta gaza. Idan gudanarwar ta gaza ba tare da `--fix` ba, ccam yana ƙarewa da layi mara komai da `` `ccam doctor --fix` repairs what it can. ``. Saitin da ba za a iya lodawa ba shi ma yana ƙarewa da `1` da kuskuren `ccam:`.

## Yadda yake aiki

`internal/cli/doctor.go` yana loda `~/.config/ccam/config.toml` kuma yana kiran `doctor.Run` a `internal/doctor/doctor.go`, wanda ke gudanar da dubawar a ƙayyadadden tsari kuma yana mayar da sakamako. Ana miƙa `--fix` a matsayin `Options.Fix`, wanda ake saita shi zuwa false a dry run.

<!-- motion:doctor-checks -->

Toshen tsakiya yana gudana sau ɗaya ga kowane asusu. Dubawar, da abin da `--fix` ke yi ga kowanne:

1. **Sigar Claude Code** (`checkClaude`). Yana nemo binary na `claude` kamar yadda mai farawa (launcher) ke yi kuma yana gudanar da `claude --version`. Idan sigar ba ta fara da `2.1.302` ba, sigar da aka tabbatar da tsarin foldar ccam a kanta (`doctor.LayoutVersion`), yana ƙara `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Rashin binary ko gazawar `--version` kuskure ne. Babu gyara.
2. **Bayanan shiga a muhalli (environment)** (`checkAuthEnv`). Yana gargaɗi idan an saita `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ko `ANTHROPIC_AUTH_TOKEN`: kowane asusu da aka fara daga wannan shell yana gadon sa, kuma Claude Code yana fifita shi a kan shigar asusun kansa. Babu gyara; cire shi da kanka.
3. **Folda** (`checkProfile`). Rashin foldar asusu kuskure ne kuma yana kawo ƙarshen dubawar wannan asusun. Izini da ke da kowane bit na group ko other gargaɗi ne; `--fix` yana gudanar da `chmod 0700`.
4. **Fayil ɗin shiga** (`checkCredentials`). Ana duba `.credentials.json` da `lstat`. Symlink kuskure ne, domin dole kowane asusu ya sami nasa token, kuma ba a gyara shi. Mode da ba `0600` ba gargaɗi ne; `--fix` yana gudanar da `chmod 0600`. Sannan ana karanta shigar (a macOS daga Keychain da farko). Rashin shiga gargaɗi ne da ke nuna `ccam login <name>`. Refresh token da ya ƙare kuskuren `login expired` ne; access token da ya ƙare bayani ne kawai, tunda Claude Code yana sabunta shi a farawa ta gaba. Ba a gyara ko ɗaya daga cikin waɗannan ba.
5. **Mahaɗan rabawa** (`share.Check`). Yana kwatanta kowace hanya a `share.Layout` da saitunan rabawar asusun, ta amfani da `lstat` da `readlink` kawai. Yana bayar da rahoton nau'i huɗu: mahaɗin da ya ɓace, mahaɗin da ke nuna wani wuri, fayil na gaske inda ya kamata mahaɗi ya kasance, da mahaɗin da ya rage ga abin da yanzu na gida ne. Mahaɗin da ya ɓace gargaɗi ne, sauran kurakurai ne.
6. **Shigarwar da ba a sani ba** (`checkUnknownEntries`). Yana lissafa sunayen babban mataki a foldar asusun da ccam bai sani ba (`share.Known`), a matsayin alamar cewa tsarin Claude Code ya canza. Bayani ne kawai.
7. **Hanyoyin plugin** (`checkPlugins`). Rajistocin plugin (`installed_plugins.json`, `known_marketplaces.json`) suna ajiye cikakkun hanyoyi. A cikin `plugins/` na kowane rukuni, hanyoyin da ke ratsa foldar asusu ɗaya suna lalacewa idan aka cire ko aka sake sunan wannan asusun. `--fix` yana sake rubuta su zuwa foldar rukunin da `profile.RewritePluginPaths`, rubutu na atomic na fayil ɗin rajistar.
8. **Sabobin MCP** (`checkMCP`). Yana loda `mcp.json` na kowane rukuni sau ɗaya. Fayil da ba za a iya karantawa ba kuskure ne; sabar da ta gaza `mcp.Validate` gargaɗi ce. Babu gyara.
9. **Tsarin amsar amfani** (`checkUsageFormat`). Yana tambayar endpoint na amfani, ba tare da cache ba kuma da iyakar daƙiƙa 3 da aka saba, ga asusu ɗaya bayan ɗaya har sai ɗaya ya amsa. Idan ba a iya karanta amsar kuma ba, yana gargaɗin cewa `ccam usage` zai koma ga cache na layin matsayi. Wannan ita ce dubawa kaɗai da ke amfani da network.
10. **Gajerun umarni** (`checkShims`). Sai idan asusu yana da gajeren umarni. Mahaɗan da suka ɓace a `~/.local/bin` gargaɗi ne; `--fix` yana kiran `shim.Install`, wanda ke ƙirƙirar mahaɗan da suka ɓace kawai kuma ba ya taɓa wasu fayiloli. Foldar bin da ba ta cikin `PATH` gargaɗi ce kuma ba a gyara ta.

### Gyara mahaɗan rabawa

Hanyar farawa tana ƙirƙirar mahaɗin da ya ɓace nan take, amma tana tsayawa ga duk abin da zai iya rasa bayanai kuma tana nuna nan:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Ga kowace matsalar mahaɗin rabawa, `--fix` yana kiran `share.Build` da saitunan rabawar asusun a matsayin duka yanayin yanzu da yanayin da ake so, sannan `Plan.Apply`. `Build` yana sake samar da kowane mataki daga abin da ke kan faifai, shiri ɗaya da `ccam share` ke amfani da shi, don haka gyaran yana bin ƙa'idodi iri ɗaya da kunna rabawa.

<!-- motion:doctor-repair -->

A nan `settings.json` fayil ne na gaske. Saboda rukunin ya riga yana da `settings.json` daban, rikicin yana ɗaukar zaɓinsa na tsoho, `KeepGroup`, don haka kwafin rukunin yana nan. Fayil ɗin asusun yana matsawa zuwa `.ccam-backup/<time>/settings.json` a cikin foldar asusun, kuma `os.Symlink` yana haɗa hanyar da rukunin. Sauran yanayi:

- Idan rukunin bai da kwafi tukuna, fayil ɗin asusun yana matsawa cikin rukunin kuma ya zama abin da ke cikinsa na farko.
- Ga folda, fayilolin da rukunin bai da su suna matsawa cikin rukunin; sauran suna zuwa ajiya.
- Ga `history.jsonl`, tsohon zaɓi shi ne `merge by line`: layukan rukunin, sannan sababbin naka. Ana ajiye kwafin rukunin na baya a `.ccam-backup/<time>/group/`.
- Ana maye gurbin mahaɗi zuwa wani rukuni. Mahaɗi zuwa foldar da ke wajen ccam yana zuwa ajiya, kuma abin da ke cikinsa yana zama farkon abin rukunin idan rukunin babu komai a ciki.
- Ana cire mahaɗin da ya rage ga abin da yanzu na gida ne, kuma asusun yana samun kwafin abin da ke cikin rukunin.

Banda cire mahaɗai, kowane mataki yana matsar da bayanai cikin rukunin ko foldar ajiya. Ba a goge komai.

## Tsaro da yanayi na musamman

- `ccam doctor` kawai ba ya taɓa canza komai. Yi amfani da `--dry-run` tare da `--fix` don tabbatarwa a cikin rubutun (scripts).
- `doctor --fix` ba ya duba ko asusu yana aiki. Rufe Claude Code na asusu kafin ka gyara mahaɗan rabawarsa, domin zaman yana ci gaba da rubutu a waɗannan fayiloli.
- Idan gyaran rabawa ya gaza, ccam yana bayar da rahoton `could not repair share links: <error>` ga wannan asusun kuma yana tsallake dubawar shigarwar da ba a sani ba. Ana ci gaba da duba sauran asusai.
- Idan `ccam share` ya gaza a tsakiya, yana fara ajiye saitunan da ake so kuma yana gaya maka ka kammala da `ccam doctor --fix`.
- Ana bayar da rahoton `.credentials.json` da aka yi wa symlink, shigar da ta ƙare da bayanan shiga a muhalli, amma ba a taɓa canza su.
- Ana rubuta sake rubutun rajistar plugin a atomic, amma ba a kwafar su zuwa `.ccam-backup/`.

## Masu alaƙa

- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Kwafar saitin da ke akwai](/features/migrate)
- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Umarnai](/docs/commands)
