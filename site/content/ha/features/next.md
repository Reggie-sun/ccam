---
title: Sauyawa zuwa asusu na gaba
description: ccam next yana buɗe Claude Code a matsayin wani asusu; --auto yana zaɓar mafi sarari, --continue yana kai tattaunawar gaba, --fallback yana maimaita ayyukan -p.
updated: 2026-10-24
order: 8
translation: machine
sourceHash: a38daed5dd97
---

`ccam next` yana buɗe Claude Code a matsayin wani asusu ba tare da ka ambaci sunansa ba. Yi amfani da shi idan asusu ya kai iyakar amfaninsa: `ccam next` kawai yana ɗaukar asusun da ke bayan wanda ka yi amfani da shi na ƙarshe, `--auto` yana ɗaukar wanda ya fi sarari, kuma `--continue` yana kai tattaunawar ƙarshe ta wannan foldar zuwa sabon asusun, ko asusun biyu suna raba tarihi ko ba sa raba. A cikin Claude Code, [plugin ɗin ccam](#inside-claude-code) yana sanar da kai idan iyaka ta dakatar da zama, kuma `/ccam:switch` yana miƙa tattaunawar gaba. Don rubutun umarni (scripts) masu gudanar da `claude -p`, `ccam --fallback` yana yin sauyawar da kansa (duba [Sake gudanarwa a iyaka](#run-again-at-the-limit)).

**Kujeru (seats) na Team da Enterprise.** Zaɓin atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba [bayanin sharuɗɗa a cikin Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

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
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam yana buga sauyawar a matsayin `→ <from> → <to>` a standard error, tare da dalilin zaɓin `--auto`, sannan yana fara Claude Code a matsayin sabon asusun kamar yadda `ccam <name>` ke yi. `<from>` shi ne asusun da ka yi amfani da shi na ƙarshe, ko da `--continue` asusun da tattaunawar take nasa. Arguments bayan `--` suna zuwa `claude` ba tare da canji ba.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--auto` | Zaɓi asusun da aka shiga da shi bisa amfani (duba ƙasa), ana tsallake asusan da suka kai 100% | a kashe |
| `--max-usage <1-100>` | Da `--auto`: fifita asusan da ke ƙasa da wannan amfani | `auto_max_usage`, in ba haka ba 100 |
| `--strategy headroom\|expiring\|balance` | Da `--auto`: mafi ƙarancin amfanin awa 5, sake saitin awa 5 mafi kusa, ko kason mako da zai fara ɓata ba tare da an yi amfani da shi ba | `auto_strategy`, in ba haka ba `headroom` |
| `--cross-org` | Da `--auto`: zaɓi asusan wata ƙungiya (organization) dabam ma (duba ƙasa) | `cross_org`, in ba haka ba a kashe |
| `--wait[=<duration>]` | Da `--auto`: idan babu asusun da ke da sarari, jira na farko ya sake saitawa, har zuwa iyakar wannan tsawon. Ƙimar tana zuwa bayan `=`: ana ƙin `--wait 2h` da shawarar `--wait takes its value after '=': --wait=2h` | a kashe; `5h` idan an ba shi shi kaɗai |
| `-c`, `--continue` | Kai tattaunawar ƙarshe ta wannan foldar zuwa sabon asusun (duba ƙasa) | a kashe |
| `-- <claude args>` | Arguments da ake miƙa wa `claude` bayan na ccam kansa | babu |
| `--dry-run` | Ana karɓar sa: yana buga `→ <from> → <to>` da layin umarnin da zai gudanar, ba tare da fara Claude Code ko rubuta farawa ba | a kashe |

`--auto` da `--continue` suna aiki tare. Da `--continue`, ana saka arguments na `--resume` kafin naka.

Zaɓuɓɓukan Claude Code na kansa, kamar `--dangerously-skip-permissions` ko `--model`, suna zuwa bayan `--`. Ba kamar `ccam <name>` ba, wanda ke miƙa duk abin da ke bayan sunan gaba, `ccam next` yana karanta arguments ɗinsa da kansa, don haka ba tare da `--` ba yana tsayawa da `ccam: unknown flag: --dangerously-skip-permissions`.

## Yadda yake aiki

Kowane farawa ta ccam yana rubuta lokacin zuwa `~/.local/state/ccam/last-used/<name>`. `ccam next` yana karanta waɗannan fayiloli, yana duba asusan da aka shiga da su, kuma yana zaɓa.

1. `cli.candidates` yana gina `next.Candidate` ɗaya ga kowane asusu daga `config.toml`: asusun, lokacin farawarsa na ƙarshe (`launcher.LastUsed`) da ko yana da shiga mai aiki (`creds.Has`: shigar Claude a `.credentials.json`, ko a macOS a Keychain, wadda ba ta ƙare gaba ɗaya ba ko Claude Code bai share ta ba). Foldar da ke amfani da API key maimakon haka (shigar Console, ko asusun da aka ƙara da `ccam add --api`) ba ta da iyakokin tsari kuma ba a taɓa zaɓar ta ba.
2. Ba tare da `--auto` ba, `next.Rotate` yana jera asusan bisa lambar rami kuma yana nemo wanda aka fara kwanan nan. Wannan shi ne asusun da kake sauyawa daga gare shi. Idan ba a fara kowane asusu ba tukuna, shi ne mafi ƙanƙantar rami.
3. Daga nan yana tafiya gaba a tsarin ramuka, yana komawa daga rami na ƙarshe zuwa na farko, kuma yana mayar da asusun farko da aka shiga da shi.
4. `launcher.LaunchWith` yana shirya kuma yana fara asusun da aka zaɓa daidai kamar yadda `ccam <name>` zai yi: yana duba mahaɗan asusun, yana aiwatar da rabawar MCP, yana rubuta zaman kuma yana maye gurbin ccam da `claude` ta `syscall.Exec`.

Tafiyar tana farawa daga asusun da aka yi amfani da shi na ƙarshe kuma tana tsayawa a asusun farko da aka shiga da shi bayan sa:

<!-- motion:next-rotate -->

An fara `work` kwanan nan, don haka tafiyar tana farawa bayan sa. `side` ba shi da shiga kuma ana tsallake shi, kuma ana zaɓar `lab`. Da a ce `lab` ma an fita daga gare shi, tafiyar za ta zagaya zuwa `personal` (`% len` a zanen). Juyawa ba ta taɓa mayar da asusun da ka fara daga gare shi ba: idan babu wani asusu da aka shiga da shi, tana tsayawa da kuskure.

### --auto

`--auto` yana yanke shawara bisa amfani maimakon tsarin ramuka. `cli.pickAuto` yana riƙe asusan da aka shiga da su, yana tambayar amfaninsu duka a lokaci ɗaya da `Fetcher.FetchAll` ɗaya da `ccam usage` ke amfani da shi, kuma yana barin `usage.Policy.Choose()` ya zaɓa:

<!-- motion:next-auto -->

Ana auna asusu da tagoginsa na awa 5 da na mako, kuma, idan `--model` yana cikin arguments na Claude Code (bayan `--`), da tagar mako ta wannan model ɗin ma (`opus`, `sonnet`, `haiku` ko `fable` a cikin sunan). `Policy.Choose()` yana raba asusan zuwa rukuni uku:

1. **Sarari**: kowace taga tana ƙasa da `auto_max_usage` (`config.toml`, tsoho 100). Dabarar ce ke zaɓa a cikinsu: `headroom` (tsoho) yana ɗaukar mafi ƙarancin amfanin awa 5, `expiring` yana ɗaukar tagar awa 5 da za ta sake saitawa da wuri, don a yi amfani da sauran quota ɗinta kafin ta ɓace; `balance` yana ɗaukar asusun da ke da mafi yawan quota na mako da ya rage ga kowace rana har makonsa ya sake saitawa, yana ƙirga kowane zaman ccam da ke gudana a kansa, don kada a bar quota na mako ba a yi amfani da shi ba kuma zaman da ke gudana a lokaci ɗaya su rarrabu.
2. **Sama**: ƙasa da 100% amma a ko sama da `auto_max_usage`. Ana amfani da su ne kawai idan babu asusun da ke da sarari, mafi ƙarancin amfanin awa 5 da farko.
3. **Cike**: taga a 100%, ko amfanin awa 5 da ba a sani ba. Ba a taɓa zaɓar su ba.

Ana ƙirga taga a matsayin 100% har zuwa sake saitinta idan Claude Code ya ƙi buƙata a can: hook na [plugin ɗin ccam](#inside-claude-code) yana rubuta hakan (`~/.local/state/ccam/usage/<name>.limit.json`), ko me lambobin da aka ajiye suka ce. Don haka asusun da ya kai iyakarsa yanzu ba a sake zaɓar sa kafin sake saitinsa, ko da lambobinsa sun tsufa.

Saituna na dukkan asusai suna zuwa a saman `config.toml`, kafin `[[group]]`, `[[profile]]` ko `[[folder]]` na farko. Saitin da ke ƙarƙashin ɗaya daga cikinsu na wannan tebur (table) ne: ccam yana ambatar sa a kowane farawa, kuma umarnan da ke ajiye config suna ƙi har sai an matsar da shi.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Ƙungiyoyi (organizations).** Zaɓuɓɓukan atomatik suna zama tare da asusan da suka yi kama da wanda kake barin: ƙungiyar team ko enterprise iri ɗaya, ko tsare-tsaren kai (Pro, Max) tare, don kada aikin asusun kamfani ya ci gaba a asusun kai ko akasin haka. ccam yana sanin ƙungiyar kowane asusu daga shigarsa (`ccam login` yana rubuta ta; ana karanta tsofaffin shige-shige daga `.claude.json`); asusun da ba zai iya tantancewa ba yana dacewa a ko'ina. Foldoji biyu da suka shiga asusun Claude iri ɗaya suna raba iyakokinsa, don haka ɗaya ba ya taɓa karɓar aiki daga ɗayan. `--cross-org`, ko `cross_org = true` a `config.toml`, yana ɗage ƙa'idar ƙungiya.

Idan ba za a iya zaɓar komai ba, ccam yana fita da code 75 (`EX_TEMPFAIL`), don rubutun umarni ya bambanta "a sake gwadawa daga baya" da gazawa. Da `--wait`, maimakon haka yana barci har tagar da ke toshe asusun farko da ya cika ta sake saitawa (daƙiƙa 30 bayan lokacin sake saitin, lokacin da endpoint ɗin ke nuna sabuwar tagar), yana buga `ccam: no account has room; waiting until 15:10`, kuma yana sake zaɓa; idan babu asusun da ya faɗi lokacin da zai sake saitawa, yana sake duba bayan minti 10. Idan babu abin da ya samu sarari kafin jiran ya ƙare, yana fita da 75. Neman yana bin ƙa'idodin `ccam usage`: ana sake amfani da lambobin da ba su kai daƙiƙa 60 ba, kowane asusu yana da daƙiƙa 3, asusun da bai amsa ba, ko wanda access token ɗinsa ya ƙare, ana auna shi da lambobinsa da aka ajiye, kuma rikodin iyaka yana rufe su. Tagar da aka ajiye wadda lokacin sake saitinta ya wuce ana ƙirga ta a matsayin 0%, don haka ba a tsallake asusun da tagarsa ta sake saitawa tun daga lokacin ba. `--auto` ba ya taɓa sabunta token. Ba kamar juyawa ta yau da kullum ba, `--auto` zai iya zaɓar asusun da ka yi amfani da shi na ƙarshe idan har yanzu shi ne ya fi sarari.

### --continue

1. ccam yana neman tattaunawar ƙarshe ta wannan foldar a kowane asusu: a cikin fayilolin `<session>.jsonl` a `projects/<folder>/` (sunan foldar shi ne foldar aiki, bayan warware mahaɗai, inda aka mai da kowane hali da ba harafi ko lamba ba zuwa `-`, kamar yadda Claude Code ke sanya suna), wanda ke da prompt ko amsa mafi sabo. Ana tsallake zaman da ba su da ko ɗaya, kamar zaman da ya gudanar da `/mcp`, `/usage` ko `/status` kawai, ko fayil mara komai. Idan wani zama ya bar bayani (iyakar amfani ko `/ccam:switch`, duba ƙasa) kuma babu tattaunawar da ke da prompt ko amsa fiye da minti 2 bayan sa, ana amfani da ainihin wannan zaman maimakon haka; layukan da Claude Code ke rubutawa a `/exit` ba sa ƙirguwa.
2. Asusun da take nasa shi ne wanda kake sauyawa daga gare shi, kuma ana zaɓar sabon asusun kamar yadda yake a sama, a cikin ƙungiya ɗaya. Idan asusai da dama suna raba wannan tarihin, ta wanda `.claude.json` nasa ke cewa shi ne ya gudanar da tattaunawar na ƙarshe ce (`lastSessionId`), in ba haka ba ta wanda aka yi amfani da shi na ƙarshe.
3. Idan sabon asusun yana raba wannan tarihin (foldar rukuni ɗaya), ccam yana fara `claude --resume <session id>`, ainihin wannan tattaunawar (`--continue` zai buɗe duk wadda ta fi sabo a foldar).
4. In ba haka ba yana fara `claude --resume <the conversation's file> --fork-session`: Claude Code yana kwafar dukkan tattaunawar, ƙarƙashin sabon session id, cikin tarihin sabon asusun kansa, kuma ta asali tana zama yadda take. Daga nan `--continue` na sabon asusun yana samun ta.

ccam yana faɗin girman tattaunawar: sabon asusun yana sake karanta dukanta da amsarsa ta farko. Ba a taɓa raba prompt caches tsakanin ƙungiyoyi ba, kuma kowane tsarin kai ƙungiya ce ta kansa, don haka idan asusan biyu na ƙungiyoyi daban-daban ne ccam yana ƙara `prompt caches are not shared between organizations, so the first reply reads it all again`. Idan tsohon asusun ya shiga sabobin MCP da sabon bai shiga ba, ccam yana ambaton su: `work is not signed in to MCP servers personal is: linear (sign in with /mcp there; sign-ins are not shared)`. Raba `history` a rukuni yana guje wa kwafin, yana riƙe tattaunawa ɗaya, kuma yana raba checkpoints da `/rewind` ke maido da lamba daga gare su (`file-history/`).

## A cikin Claude Code

[Plugin ɗin ccam](/docs/install#claude-code-plugin) yana ƙara hooks biyu (Claude Code 2.1.78 ko sabo):

- Idan juyi ya tsaya a iyakar amfani, yana rubuta bayanin tattaunawar kuma yana nuna sanarwar desktop a terminals da ke goyon bayanta (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` carries on as another account.`` Yana sanin kowane saƙo da Claude Code ke ƙirgawa a matsayin iyakar amfani: iyakokin awa 5, na mako, na yau da kullum da na model (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` ga kason model), credits na amfani da iyakokin kashe kuɗi, da ƙungiyar da ta dakatar da asusun. Idan Claude Code ya rubuta buƙatar da aka ƙi a cikin tattaunawar, hook ɗin yana kuma rubuta cikakkiyar tagar har zuwa sake saitinta ([rikodin iyaka](/features/usage#the-limit-record)), sai dai dakatarwar ƙungiya, wadda admins ɗinta za su iya ɗagewa a kowane lokaci. Takaita gudu na ɗan lokaci daga sabar (throttle) ko iyakar saurin (rate limit) na API key ba iyakar amfani ba ne, kuma ana barin su, haka ma dakatarwar subagent da zaman da babu wanda ke zaune a gabansu (`-p`, zaman bayan fage). Claude Code 2.1.234 da sabo kuma yana jira a cikin zaman kuma yana ci gaba da kansa bayan sake saitin awa 5; sanarwar tana ba ka zaɓin sauyawa maimakon haka.
- `/ccam:switch [account]` yana rubuta bayanin tattaunawar kuma yana kawo ƙarshen juyin kafin a yi kowace buƙata: rubuta `/exit`, sannan `ccam next --continue` yana ci gaba da ita, a matsayin asusun da ka ambata ko na gaba. A zaman da ba shi da tattaunawa tukuna, yana cewa babu abin da za a ci gaba da shi.

Bayanin yana zama a `~/.local/state/ccam/handoff/<account>/<session>.json`, fayil ɗaya ga kowane zama, don haka zama biyu na asusu ɗaya ba sa taɓa rubuta kan na juna. Ana riƙe shi na awa 12 kuma ana amfani da shi sau ɗaya, lokacin da zama na gaba ke shirin farawa; farawar da ta gaza tana barin sa don ƙoƙari na gaba. Na foldar da zaman ya fara a cikinta ne, ko bayan `cd` a cikin shell na Claude Code. `ccam rename` yana sabunta sunayen da ke cikinsa; idan an cire asusun da ya ambata tun daga lokacin, ccam yana faɗin haka kuma yana zaɓa kamar yadda ya saba.

### Relay

Da `relay = true` a saman `config.toml`, ccam yana zama a gefen zaman mu'amala maimakon maye gurbin kansa da Claude Code, kuma bayan `/ccam:switch` yana kai tattaunawar wannan zaman gaba da kansa (ba ya aiki da bayanin da wani zaman ya bari): hook ɗin yana rufe Claude Code, kuma ccam yana gudanar da `ccam next --continue` da bayanin wannan zaman. Idan ba a ambaci asusu ba, yana ɗaukar wanda ya fi sarari, yana ƙirga model na tattaunawar. A Windows hook ɗin ba zai iya rufe Claude Code ba: rubuta `/exit`, kuma ccam yana ci gaba. Hakan yana cin process ɗaya na ccam mai jira ga kowane zama, shi ya sa yake a kashe ta tsohuwa.

Zaman da aka miƙa yana riƙe abin da ci gaba da tattaunawa (resume) ba ya maidowa: zaɓuɓɓuka kamar `--add-dir`, `--mcp-config`, jerin kayan aiki (tools), `--settings`, `--plugin-dir` da zaɓuɓɓukan system prompt, yanayin izini (permission mode) na tattaunawar da matakin ƙoƙari (effort level) na zaman. Claude Code da kansa yana maido da `--model`, `--name` da `--agent` daga tattaunawar. Duk wani abu dabam (`--worktree`, `--debug`, `--remote-control`, zaɓuɓɓukan da ccam bai sani ba) ana ambaton sa a `ccam: not carried over: …`. Sauyawar ba ta taɓa yin tambayar tarihin da aka raba ba, kuma abin da ka rubuta yayin da take gudana ana watsar da shi.

Relay yana ci gaba ne kawai idan Claude Code ya ƙare yadda `/ccam:switch` ya buƙace shi: ba idan ka rubuta `/exit` da farko ba, ko ka danna Ctrl+C yayin da yake ƙarewa, ko ccam ya sami hangup ko wani signal, kuma ba idan terminal ɗin ya tafi ba. Sannan bayanin yana zama, kuma `ccam next --continue` yana kai tattaunawar gaba daga baya. Zaman da aka fara ƙarƙashin relay na tsohon ccam, ko wanda babu wanda ke zaune a gabansa (`-p`, zaman bayan fage), hook ɗin ba ya rufe shi; yana samun saƙon `/exit` maimakon haka. `claude remote-control` yana gudana ba tare da relay ba.

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
2. Ana karanta standard input da aka tura ta bututu (pipe) sau ɗaya, har zuwa 10 MiB, kuma ana ba kowane yunƙuri shi. Kamar Claude Code, ccam yana jira daƙiƙa 3 a mafi yawa don ya fara (ba da `--input-format stream-json` ba); bututun da ya kasance a buɗe kuma bai aika komai ba yana ba gudun rashin input, tare da layi a standard error.
3. ccam yana fara `claude` a matsayin child process maimakon `exec`, yana riƙe standard output ɗinsa, kuma yana miƙa standard error yadda yake zuwa.
4. Idan `claude` ya fita da 0, ko ya gaza saboda wani dalili, ana miƙa output ɗin da aka riƙe da exit code ba tare da canji ba. Haka ma gudun da signal ya ƙare (exit code sama da 128): an nemi ya tsaya, kuma ba a sake gudanar da shi.
5. Gazawa tana ƙirguwa a matsayin iyakar amfani idan saƙon Claude Code ya faɗi haka, kowane daga cikin saƙonnin iyakar amfani da [hook ɗin](#inside-claude-code) ya sani, har a cikin output na `--output-format json` ko `stream-json`, ko idan amfanin asusun ya nuna tagarsa ta awa 5, ta mako ko ta `--model` a 100% (har da rikodin iyaka).
6. Sannan ana watsar da output ɗin wannan yunƙurin kuma umarnin iri ɗaya yana gudana a matsayin wani asusu da aka shiga da shi: wanda ke da mafi ƙarancin amfanin awa 5, sannan asusan da ba za a iya neman amfaninsu ba, bisa tsarin lambobi. Ana tsallake asusan da suka kai 100%, ba a taɓa gwada asusan da ke da shigar API key ba, kuma ana gwada kowane asusu sau ɗaya kawai.
7. Idan akwai `-c`, `--continue`, `-r` ko `--resume` a cikin arguments, asusan da ke rukunin tarihi ɗaya ne kawai ake gwadawa, don a sami tattaunawar.
8. Asusan ƙungiya ɗaya ne kawai ake gwadawa, kamar da `--auto` (`cross_org = true` yana ɗage wannan).

Ana zaɓar asusu na gaba bisa ƙa'idodi iri ɗaya da na `--auto` (`auto_max_usage`, `auto_strategy`, tagar `--model`). Idan babu sauran asusu, ccam yana buga output na yunƙurin ƙarshe, yana cewa babu asusun da aka shiga da shi da ke da sauran sarari, kuma yana fita da code 75. Da `--wait` kafin sunan asusun (`ccam --fallback --wait=2h work -p …`), maimakon haka yana jiran asusun farko ya sake saitawa, sannan zai iya sake gwada kowane asusu; koyaushe yana fara yin barci aƙalla sau ɗaya, don kada a sake gudanar da asusun da ya kai iyakarsa yanzu nan take. Tsawon yana iyakance dukkan jiran gudun tare, ba kowane jira ɗaya ba. Idan an rubuta shi da sarari (`--wait 30m`), ana ƙin sa da shawara iri ɗaya da ta `ccam next`. `--` na farko kafin arguments na Claude Code yana aiki kamar a kowane farawa (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` yana buga umarnin yunƙuri na farko.

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
- `--fallback` yana riƙe standard output har sai yunƙuri ya ƙare, don haka output mai gudana kamar `--output-format stream-json` yana zuwa gaba ɗaya a ƙarshe. Abokin aiki (client) da ke tura input na `stream-json` kuma yana jiran amsoshi kafin ya aika ƙari ba zai iya amfani da `--fallback` ba. Fayilolin da yunƙurin farko ya canza kafin ya kai iyaka suna nan a canje; yunƙuri na gaba yana sake fara prompt ɗin.

## Masu alaƙa

- [Iyakokin amfani](/features/usage)
- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow)
