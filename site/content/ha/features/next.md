---
title: Sauyawa zuwa asusu na gaba
description: ccam next yana buɗe Claude Code da asusun da ke bin na ƙarshe; --auto yana zaɓar mafi sarari, --continue yana ci gaba, --fallback yana maimaita ayyukan -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next` yana buɗe Claude Code a matsayin wani asusu ba tare da ka ambaci sunansa ba. Yi amfani da shi idan asusu ya kai iyakar amfaninsa: `ccam next` kawai yana ɗaukar asusun da ke bayan wanda ka yi amfani da shi na ƙarshe, `--auto` yana ɗaukar wanda ya fi sarari a tagar awa 5 ɗinsa, kuma `--continue` yana ci gaba da tattaunawar ƙarshe a sabon asusun. Don rubutun umarni (scripts) masu gudanar da `claude -p`, `ccam --fallback` yana yin wannan sauyawar da kansa (duba [Sake gudanarwa a iyaka](#run-again-at-the-limit)).

## Amfani

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam yana buga sauyawar a matsayin `→ <from> → <to>` a standard error, sannan yana fara Claude Code a matsayin sabon asusun kamar yadda `ccam <name>` ke yi. `<from>` shi ne asusun da ka yi amfani da shi na ƙarshe. Arguments bayan `--` suna zuwa `claude` ba tare da canji ba.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--auto` | Zaɓi asusun da aka shiga da shi bisa amfani (duba ƙasa), ana tsallake asusan da suka kai 100% | a kashe |
| `--max-usage <1-100>` | Da `--auto`: fifita asusan da ke ƙasa da wannan amfani | `auto_max_usage`, in ba haka ba 100 |
| `--strategy headroom\|expiring` | Da `--auto`: mafi ƙarancin amfanin awa 5, ko sake saitin awa 5 mafi kusa | `auto_strategy`, in ba haka ba `headroom` |
| `--wait[=<duration>]` | Da `--auto`: idan babu asusun da ke da sarari, jira na farko ya sake saitawa, har zuwa iyakar wannan tsawon | a kashe; `5h` idan an ba shi shi kaɗai |
| `-c`, `--continue` | Duba asusan da ke raba tarihin tattaunawa da asusun da ka yi amfani da shi na ƙarshe kawai, kuma fara `claude --continue` | a kashe |
| `-- <claude args>` | Arguments da ake miƙa wa `claude` bayan na ccam kansa | babu |
| `--dry-run` | Ana karɓar sa: yana buga `→ <from> → <to>` da layin umarnin da zai gudanar, ba tare da fara Claude Code ko rubuta farawa ba | a kashe |

`--auto` da `--continue` suna aiki tare. Da `--continue`, ana saka `--continue` kafin arguments naka.

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

1. **Sarari**: kowace taga tana ƙasa da `auto_max_usage` (`config.toml`, tsoho 100). Dabarar ce ke zaɓa a cikinsu: `headroom` (tsoho) yana ɗaukar mafi ƙarancin amfanin awa 5, `expiring` yana ɗaukar tagar awa 5 da za ta sake saitawa da wuri, don a yi amfani da sauran quota ɗinta kafin ta ɓace.
2. **Sama**: ƙasa da 100% amma a ko sama da `auto_max_usage`. Ana amfani da su ne kawai idan babu asusun da ke da sarari, mafi ƙarancin amfanin awa 5 da farko.
3. **Cike**: taga a 100%, ko amfanin awa 5 da ba a sani ba. Ba a taɓa zaɓar su ba.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

Idan ba za a iya zaɓar komai ba, ccam yana fita da code 75 (`EX_TEMPFAIL`), don rubutun umarni ya bambanta "a sake gwadawa daga baya" da gazawa. Da `--wait`, maimakon haka yana barci har tagar da ke toshe asusun farko da ya cika ta sake saitawa (daƙiƙa 30 bayan lokacin sake saitin, lokacin da endpoint ɗin ke nuna sabuwar tagar), yana buga `ccam: no account has room; waiting until 15:10`, kuma yana sake zaɓa; idan babu asusun da ya faɗi lokacin da zai sake saitawa, yana sake duba bayan minti 10. Idan babu abin da ya samu sarari kafin jiran ya ƙare, yana fita da 75. Neman yana bin ƙa'idodin `ccam usage`: ana sake amfani da lambobin da ba su kai daƙiƙa 60 ba, kowane asusu yana da daƙiƙa 3, kuma asusun da bai amsa ba, ko wanda shigarsa ta ƙare, ana auna shi da lambobinsa da aka ajiye. Tagar da aka ajiye wadda lokacin sake saitinta ya wuce ana ƙirga ta a matsayin 0%, don haka ba a tsallake asusun da tagarsa ta sake saitawa tun daga lokacin ba. `--auto` ba ya taɓa sabunta token. Ba kamar juyawa ta yau da kullum ba, `--auto` zai iya zaɓar asusun da ka yi amfani da shi na ƙarshe idan har yanzu shi ne ya fi sarari.

### --continue

`--continue` na Claude Code yana buɗe tattaunawa ta kwanan nan a foldar yanzu, wanda ke aiki ne kawai idan sabon asusun zai iya ganin wannan tattaunawar. Don haka da `--continue`, ccam yana duba asusan da ke rukunin `history` ɗaya da asusun da ka yi amfani da shi na ƙarshe kawai, sannan yana saka `--continue` a gaban arguments naka. Asusan da tarihinsu `local` ne ba sa raba wa kowa, don haka `--continue` daga irin wannan asusu koyaushe yana gazawa. Fara kunna rabawa da `ccam share <name> --history on` (duba [rabawa](/docs/sharing)).

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

Ana zaɓar asusu na gaba bisa ƙa'idodi iri ɗaya da na `--auto` (`auto_max_usage`, `auto_strategy`, tagar `--model`). Idan babu sauran asusu, ccam yana buga output na yunƙurin ƙarshe, yana cewa babu asusun da aka shiga da shi da ke da sauran sarari, kuma yana fita da code 75. Da `--wait` kafin sunan asusun (`ccam --fallback --wait=2h work -p …`), maimakon haka yana jiran asusun farko ya sake saitawa, sannan zai iya sake gwada kowane asusu; koyaushe yana fara yin barci aƙalla sau ɗaya, don kada a sake gudanar da asusun da ya kai iyakarsa yanzu nan take. `ccam --dry-run --fallback work -p …` yana buga umarnin yunƙuri na farko.

Ba a sauya zaman mu'amala da kansu. Maimakon haka, da `ccam statusline install`, layin matsayi yana ƙara `→ ccam next --auto --continue` da zarar amfanin awa 5 ko na mako ya kai 90% ([Iyakokin amfani](/features/usage)).

## Tsaro da yanayi na musamman

- Ba a taɓa zaɓar asusan da ba su da shiga. Idan babu wani asusu da aka shiga da shi, `ccam next` kawai yana tsayawa da `no logged-in account other than <name>`.
- Da `--continue` kuma babu asusun da ke raba tarihi: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Da `--auto` kuma babu asusun da aka shiga da shi: `no logged-in account to pick`. Idan kowane asusu ya cika ko ba a san amfaninsa ba: `` no account with known usage and headroom (see `ccam usage`) ``. Dukansu suna ƙarewa da exit code 75.
- `--max-usage`, `--strategy` da `--wait` suna aiki ne kawai da `--auto`, kuma ana ƙin ƙima a wajen 1-100 ko wata dabara dabam.
- `--auto` yana yin buƙatun network ko da `--dry-run`. Yana karantawa ne kawai, kuma ba ya taɓa rubuta fayilolin token.
- Asusun da aka zaɓa yana farawa ta hanya ɗaya da `ccam <name>`, don haka dubawa iri ɗaya suna aiki, har da gargaɗin idan wani Claude Code da ke raba tarihi ɗaya yana aiki tuni a wannan foldar.
- Ana rubuta `last-used` ne kawai idan Claude Code ya fara da gaske. `--dry-run` yana barin sa, don haka dry run ba ya matsar da juyawar gaba.
- `--fallback` yana riƙe standard output har sai yunƙuri ya ƙare, don haka output mai gudana kamar `--output-format stream-json` yana zuwa gaba ɗaya a ƙarshe. Fayilolin da yunƙurin farko ya canza kafin ya kai iyaka suna nan a canje; yunƙuri na gaba yana sake fara prompt ɗin.

## Masu alaƙa

- [Iyakokin amfani](/features/usage)
- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow)
