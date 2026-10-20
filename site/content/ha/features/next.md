---
title: Sauyawa zuwa asusu na gaba
description: ccam next yana buɗe Claude Code a matsayin wani asusu; --auto yana zaɓar mafi sarari, --continue yana kai tattaunawar gaba, --fallback yana maimaita ayyukan -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

`ccam next` yana buɗe Claude Code a matsayin wani asusu ba tare da ka ambaci sunansa ba. Yi amfani da shi idan asusu ya kai iyakar amfaninsa: `ccam next` kawai yana ɗaukar asusun da ke bayan wanda ka yi amfani da shi na ƙarshe, `--auto` yana ɗaukar wanda ya fi sarari, kuma `--continue` yana kai tattaunawar ƙarshe ta wannan foldar zuwa sabon asusun, ko asusun biyu suna raba tarihi ko ba sa raba. A cikin Claude Code, [plugin ɗin ccam](#inside-claude-code) yana sanar da kai idan iyaka ta dakatar da zama, kuma `/ccam:switch` yana miƙa tattaunawar gaba. Don rubutun umarni (scripts) masu gudanar da `claude -p`, `ccam --fallback` yana yin sauyawar da kansa (duba [Sake gudanarwa a iyaka](#run-again-at-the-limit)).

## Amfani

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam yana buga sauyawar a matsayin `→ <from> → <to>` a standard error, tare da dalilin zaɓin `--auto`, sannan yana fara Claude Code a matsayin sabon asusun kamar yadda `ccam <name>` ke yi. `<from>` shi ne asusun da ka yi amfani da shi na ƙarshe, ko da `--continue` asusun da tattaunawar take nasa. Arguments bayan `--` suna zuwa `claude` ba tare da canji ba.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--auto` | Zaɓi asusun da aka shiga da shi bisa amfani (duba ƙasa), ana tsallake asusan da suka kai 100% | a kashe |
| `--max-usage <1-100>` | Da `--auto`: fifita asusan da ke ƙasa da wannan amfani | `auto_max_usage`, in ba haka ba 100 |
| `--strategy headroom\|expiring\|balance` | Da `--auto`: mafi ƙarancin amfanin awa 5, sake saitin awa 5 mafi kusa, ko kason mako da zai fara ɓata ba tare da an yi amfani da shi ba | `auto_strategy`, in ba haka ba `headroom` |
| `--cross-org` | Da `--auto`: zaɓi asusan wata ƙungiya (organization) dabam ma (duba ƙasa) | `cross_org`, in ba haka ba a kashe |
| `--wait[=<duration>]` | Da `--auto`: idan babu asusun da ke da sarari, jira na farko ya sake saitawa, har zuwa iyakar wannan tsawon | a kashe; `5h` idan an ba shi shi kaɗai |
| `-c`, `--continue` | Kai tattaunawar ƙarshe ta wannan foldar zuwa sabon asusun (duba ƙasa) | a kashe |
| `-- <claude args>` | Arguments da ake miƙa wa `claude` bayan na ccam kansa | babu |
| `--dry-run` | Ana karɓar sa: yana buga `→ <from> → <to>` da layin umarnin da zai gudanar, ba tare da fara Claude Code ko rubuta farawa ba | a kashe |

`--auto` da `--continue` suna aiki tare. Da `--continue`, ana saka arguments na `--resume` kafin naka.

## Yadda yake aiki

Kowane farawa ta ccam yana rubuta lokacin zuwa `~/.local/state/ccam/last-used/<name>`. `ccam next` yana karanta waɗannan fayiloli, yana duba asusan da aka shiga da su, kuma yana zaɓa.

1. `cli.candidates` yana gina `next.Candidate` ɗaya ga kowane asusu daga `config.toml`: asusun, lokacin farawarsa na ƙarshe (`launcher.LastUsed`) da ko yana da shiga (`creds.Has`: `.credentials.json` da ake iya karantawa, ko a macOS abu a Keychain).
2. Ba tare da `--auto` ba, `next.Rotate` yana jera asusan bisa lambar rami kuma yana nemo wanda aka fara kwanan nan. Wannan shi ne asusun da kake sauyawa daga gare shi. Idan ba a fara kowane asusu ba tukuna, shi ne mafi ƙanƙantar rami.
3. Daga nan yana tafiya gaba a tsarin ramuka, yana komawa daga rami na ƙarshe zuwa na farko, kuma yana mayar da asusun farko da aka shiga da shi.
4. `launcher.LaunchWith` yana shirya kuma yana fara asusun da aka zaɓa daidai kamar yadda `ccam <name>` zai yi: yana duba mahaɗan asusun, yana aiwatar da rabawar MCP, yana rubuta zaman kuma yana maye gurbin ccam da `claude` ta `syscall.Exec`.

Tafiyar tana farawa daga asusun da aka yi amfani da shi na ƙarshe kuma tana tsayawa a asusun farko da aka shiga da shi bayan sa:

<!-- motion:next-rotate -->

An fara `work` kwanan nan, don haka tafiyar tana farawa bayan sa. `side` ba shi da shiga kuma ana tsallake shi, kuma ana zaɓar `lab`. Da a ce `lab` ma an fita daga gare shi, tafiyar za ta zagaya zuwa `personal` (`% len` a zanen). Juyawa ba ta taɓa mayar da asusun da ka fara daga gare shi ba: idan babu wani asusu da aka shiga da shi, tana tsayawa da kuskure.

### --auto

`--auto` yana yanke shawara bisa amfani maimakon tsarin ramuka. `cli.pickAuto` yana riƙe asusan da aka shiga da su, yana tambayar amfaninsu duka a lokaci ɗaya da `Fetcher.FetchAll` ɗaya da `ccam usage` ke amfani da shi, kuma yana barin `usage.Policy.Choose()` ya zaɓa:

<!-- motion:next-auto -->

Ana auna asusu da tagoginsa na awa 5 da na mako, kuma, idan `--model` yana cikin arguments na Claude Code, da tagar mako ta wannan model ɗin ma (`opus`, `sonnet`, `haiku` ko `fable` a cikin sunan). `Policy.Choose()` yana raba asusan zuwa rukuni uku:

1. **Sarari**: kowace taga tana ƙasa da `auto_max_usage` (`config.toml`, tsoho 100). Dabarar ce ke zaɓa a cikinsu: `headroom` (tsoho) yana ɗaukar mafi ƙarancin amfanin awa 5, `expiring` yana ɗaukar tagar awa 5 da za ta sake saitawa da wuri, don a yi amfani da sauran quota ɗinta kafin ta ɓace; `balance` yana ɗaukar asusun da ke da mafi yawan quota na mako da ya rage ga kowace rana har makonsa ya sake saitawa, yana ƙirga kowane zaman ccam da ke gudana a kansa, don kada a bar quota na mako ba a yi amfani da shi ba kuma zaman da ke gudana a lokaci ɗaya su rarrabu.
2. **Sama**: ƙasa da 100% amma a ko sama da `auto_max_usage`. Ana amfani da su ne kawai idan babu asusun da ke da sarari, mafi ƙarancin amfanin awa 5 da farko.
3. **Cike**: taga a 100%, ko amfanin awa 5 da ba a sani ba. Ba a taɓa zaɓar su ba.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Ƙungiyoyi (organizations).** Zaɓuɓɓukan atomatik suna zama tare da asusan da suka yi kama da wanda kake barin: ƙungiyar team ko enterprise iri ɗaya, ko tsare-tsaren kai (Pro, Max) tare, don kada aikin asusun kamfani ya ci gaba a asusun kai ko akasin haka. ccam yana sanin ƙungiyar kowane asusu daga shigarsa (`ccam login` yana rubuta ta; ana karanta tsofaffin shige-shige daga `.claude.json`); asusun da ba zai iya tantancewa ba yana dacewa a ko'ina. Foldoji biyu da suka shiga asusun Claude iri ɗaya suna raba iyakokinsa, don haka ɗaya ba ya taɓa karɓar aiki daga ɗayan. `--cross-org`, ko `cross_org = true` a `config.toml`, yana ɗage ƙa'idar ƙungiya.

Idan ba za a iya zaɓar komai ba, ccam yana fita da code 75 (`EX_TEMPFAIL`), don rubutun umarni ya bambanta "a sake gwadawa daga baya" da gazawa. Da `--wait`, maimakon haka yana barci har tagar da ke toshe asusun farko da ya cika ta sake saitawa (daƙiƙa 30 bayan lokacin sake saitin, lokacin da endpoint ɗin ke nuna sabuwar tagar), yana buga `ccam: no account has room; waiting until 15:10`, kuma yana sake zaɓa; idan babu asusun da ya faɗi lokacin da zai sake saitawa, yana sake duba bayan minti 10. Idan babu abin da ya samu sarari kafin jiran ya ƙare, yana fita da 75. Neman yana bin ƙa'idodin `ccam usage`: ana sake amfani da lambobin da ba su kai daƙiƙa 60 ba, kowane asusu yana da daƙiƙa 3, kuma asusun da bai amsa ba, ko wanda shigarsa ta ƙare, ana auna shi da lambobinsa da aka ajiye. Tagar da aka ajiye wadda lokacin sake saitinta ya wuce ana ƙirga ta a matsayin 0%, don haka ba a tsallake asusun da tagarsa ta sake saitawa tun daga lokacin ba. `--auto` ba ya taɓa sabunta token. Ba kamar juyawa ta yau da kullum ba, `--auto` zai iya zaɓar asusun da ka yi amfani da shi na ƙarshe idan har yanzu shi ne ya fi sarari.

### --continue

1. ccam yana neman tattaunawar ƙarshe ta wannan foldar a kowane asusu: `<session>.jsonl` mafi sabo a `projects/<folder>/` (sunan foldar shi ne foldar aiki inda aka mai da kowane hali da ba harafi ko lamba ba zuwa `-`, kamar yadda Claude Code ke sanya suna). Idan zaman ya bar bayani (iyakar amfani ko `/ccam:switch`, duba ƙasa) kuma ba a rubuta wani abu sabo ba tun daga lokacin, ana amfani da ainihin wannan zaman maimakon haka.
2. Asusun da take nasa shi ne wanda kake sauyawa daga gare shi, kuma ana zaɓar sabon asusun kamar yadda yake a sama, a cikin ƙungiya ɗaya.
3. Idan sabon asusun yana raba wannan tarihin (foldar rukuni ɗaya), ccam yana fara `claude --resume <session id>`, ainihin wannan tattaunawar (`--continue` zai buɗe duk wadda ta fi sabo a foldar).
4. In ba haka ba yana fara `claude --resume <the conversation's file> --fork-session`: Claude Code yana kwafar dukkan tattaunawar, ƙarƙashin sabon session id, cikin tarihin sabon asusun kansa, kuma ta asali tana zama yadda take. Daga nan `--continue` na sabon asusun yana samun ta.

ccam yana faɗin girman tattaunawar: sabon asusun yana sake karanta dukanta da amsarsa ta farko. Ba a taɓa raba prompt caches tsakanin ƙungiyoyi ba, don haka da `--cross-org` wannan amsa ta farko tana cin cikakken girman. Raba `history` a rukuni yana guje wa kwafin, yana riƙe tattaunawa ɗaya, kuma yana raba checkpoints da `/rewind` ke maido da lamba daga gare su (`file-history/`).

## A cikin Claude Code

[Plugin ɗin ccam](/docs/install#claude-code-plugin) yana ƙara hooks biyu (Claude Code 2.1.78 ko sabo):

- Idan juyi ya tsaya a iyakar amfani (`You've hit your session limit`, iyakar mako ko ta model), yana rubuta bayanin tattaunawar kuma yana nuna sanarwar desktop a terminals da ke goyon bayanta (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Takaita gudu na ɗan lokaci daga sabar (throttle) ko iyakar saurin (rate limit) na API key ba iyakar amfani ba ne, kuma ana barin su. Claude Code 2.1.234 da sabo kuma yana jira a cikin zaman kuma yana ci gaba da kansa bayan sake saitin awa 5; sanarwar tana ba ka zaɓin sauyawa maimakon haka.
- `/ccam:switch [account]` yana rubuta bayanin tattaunawar kuma yana kawo ƙarshen juyin kafin a yi kowace buƙata: rubuta `/exit`, sannan `ccam next --continue` yana ci gaba da ita, a matsayin asusun da ka ambata ko na gaba.

Bayanin yana zama a `~/.local/state/ccam/handoff/<account>.json` na awa 12 kuma ana amfani da shi sau ɗaya, lokacin da zama na gaba ke shirin farawa; farawar da ta gaza tana barin sa don ƙoƙari na gaba. Na foldar da zaman ya fara a cikinta ne, ko bayan `cd` a cikin shell na Claude Code. Idan an sake wa asusun da ya ambata suna ko an cire shi tun daga lokacin, ccam yana faɗin haka kuma yana zaɓa kamar yadda ya saba. Da `relay = true` a `config.toml`, ccam yana zama a gefen zaman mu'amala maimakon maye gurbin kansa da Claude Code, kuma bayan `/ccam:switch` yana kai tattaunawar wannan zaman gaba da kansa (ba ya aiki da bayanin da wani zaman ya bari). Hakan yana cin process ɗaya na ccam mai jira ga kowane zama, shi ya sa yake a kashe ta tsohuwa. Zaman da ke ci gaba da tattaunawar yana farawa da `ccam next --continue`, ba tare da arguments da aka fara zaman farko da su ba, kamar `--model`.

## Sake gudanarwa a iyaka

`ccam next` na zaman mu'amala ne da ya tsaya. Gudu wanda ba na mu'amala ba, `claude -p` a rubutun umarni ko a aikin da aka tsara lokacinsa, ba zai iya jiranka ba, don haka `--fallback` kafin sunan asusun yana bar ccam ya sauya a madadinsa:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` yana aiki ne kawai tare da `-p` ko `--print`; in ba haka ba ccam yana tsayawa kuma yana nuna `ccam next --auto --continue`.
2. Ana karanta standard input sau ɗaya kuma ana ba kowane yunƙuri shi.
3. ccam yana fara `claude` a matsayin child process maimakon `exec`, yana riƙe standard output ɗinsa, kuma yana miƙa standard error yadda yake zuwa.
4. Idan `claude` ya fita da 0, ko ya gaza saboda wani dalili, ana miƙa output ɗin da aka riƙe da exit code ba tare da canji ba.
5. Gazawa tana ƙirguwa a matsayin iyakar amfani idan saƙon Claude Code ya faɗi haka (`You've hit your limit`, `usage limit reached`), ko idan neman amfanin asusun ya nuna tagarsa ta awa 5 ko ta mako a 100%.
6. Sannan ana watsar da output ɗin wannan yunƙurin kuma umarnin iri ɗaya yana gudana a matsayin wani asusu da aka shiga da shi: wanda ke da mafi ƙarancin amfanin awa 5, sannan asusan da ba za a iya neman amfaninsu ba, bisa tsarin lambobi. Ana tsallake asusan da suka kai 100%, kuma ana gwada kowane asusu sau ɗaya kawai.
7. Idan akwai `-c`, `--continue`, `-r` ko `--resume` a cikin arguments, asusan da ke rukunin tarihi ɗaya ne kawai ake gwadawa, don a sami tattaunawar.
8. Asusan ƙungiya ɗaya ne kawai ake gwadawa, kamar da `--auto` (`cross_org = true` yana ɗage wannan).

Ana zaɓar asusu na gaba bisa ƙa'idodi iri ɗaya da na `--auto` (`auto_max_usage`, `auto_strategy`, tagar `--model`). Idan babu sauran asusu, ccam yana buga output na yunƙurin ƙarshe, yana cewa babu asusun da aka shiga da shi da ke da sauran sarari, kuma yana fita da code 75. Da `--wait` kafin sunan asusun (`ccam --fallback --wait=2h work -p …`), maimakon haka yana jiran asusun farko ya sake saitawa, sannan zai iya sake gwada kowane asusu; koyaushe yana fara yin barci aƙalla sau ɗaya, don kada a sake gudanar da asusun da ya kai iyakarsa yanzu nan take. `ccam --dry-run --fallback work -p …` yana buga umarnin yunƙuri na farko.

Ba a sauya zaman mu'amala da kansu. Maimakon haka, da `ccam statusline install`, layin matsayi yana ƙara `→ ccam next --auto --continue` da zarar amfanin awa 5 ko na mako ya kai 90% ([Iyakokin amfani](/features/usage)).

## Tsaro da yanayi na musamman

- Ba a taɓa zaɓar asusan da ba su da shiga. Idan babu wani asusu da aka shiga da shi, `ccam next` kawai yana tsayawa da `no logged-in account other than <name>`.
- Da `--continue` kuma babu tattaunawa a foldar: `no account has a conversation in <folder> to continue`.
- Idan ƙa'idar ƙungiya ba ta bar kowane asusu ba: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Da `--auto` kuma babu asusun da aka shiga da shi: `no logged-in account to pick`. Idan kowane asusu ya cika ko ba a san amfaninsa ba: `` no account with known usage and headroom (see `ccam usage`) ``. Dukansu suna ƙarewa da exit code 75.
- `--max-usage`, `--strategy`, `--wait` da `--cross-org` suna aiki ne kawai da `--auto`, kuma ana ƙin ƙima a wajen 1-100 ko wata dabara dabam.
- Fork yana kwafar tattaunawar, ba sauran yanayin asusun ba: checkpoints na `file-history/` na asusun asali ba sa zuwa tare sai idan an raba tarihi.
- `--auto` yana yin buƙatun network ko da `--dry-run`. Yana karantawa ne kawai, kuma ba ya taɓa rubuta fayilolin token.
- Asusun da aka zaɓa yana farawa ta hanya ɗaya da `ccam <name>`, don haka dubawa iri ɗaya suna aiki, har da gargaɗin idan wani Claude Code da ke raba tarihi ɗaya yana aiki tuni a wannan foldar.
- Ana rubuta `last-used` ne kawai idan Claude Code ya fara da gaske. `--dry-run` yana barin sa, don haka dry run ba ya matsar da juyawar gaba.
- `--fallback` yana riƙe standard output har sai yunƙuri ya ƙare, don haka output mai gudana kamar `--output-format stream-json` yana zuwa gaba ɗaya a ƙarshe. Fayilolin da yunƙurin farko ya canza kafin ya kai iyaka suna nan a canje; yunƙuri na gaba yana sake fara prompt ɗin.

## Masu alaƙa

- [Iyakokin amfani](/features/usage)
- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow)
