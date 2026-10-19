---
title: Sauyawa zuwa asusu na gaba
description: ccam next yana buɗe Claude Code a matsayin asusun da ke bayan wanda ka yi amfani da shi na ƙarshe; --auto yana zaɓar mafi ƙarancin amfanin awa 5, kuma --continue yana ci gaba da tattaunawar.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` yana buɗe Claude Code a matsayin wani asusu ba tare da ka ambaci sunansa ba. Yi amfani da shi idan asusu ya kai iyakar amfaninsa: `ccam next` kawai yana ɗaukar asusun da ke bayan wanda ka yi amfani da shi na ƙarshe, `--auto` yana ɗaukar wanda ya fi sarari a tagar awa 5 ɗinsa, kuma `--continue` yana ci gaba da tattaunawar ƙarshe a sabon asusun.

## Amfani

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
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
| `--auto` | Zaɓi asusun da aka shiga da shi mai mafi ƙarancin amfanin awa 5, ana tsallake asusan da suka kai 100% | a kashe |
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

`--auto` yana yanke shawara bisa amfani maimakon tsarin ramuka. `cli.pickAuto` yana riƙe asusan da aka shiga da su, yana tambayar amfaninsu duka a lokaci ɗaya da `Fetcher.FetchAll` ɗaya da `ccam usage` ke amfani da shi, kuma yana barin `usage.Best()` ya zaɓa:

<!-- motion:next-auto -->

`usage.Best()` yana ɗaukar mafi ƙarancin kashi na awa 5 kuma yana tsallake asusan da suka kai 100% ko waɗanda ba a san amfaninsu na awa 5 ba. Neman yana bin ƙa'idodin `ccam usage`: ana sake amfani da lambobin da ba su kai daƙiƙa 60 ba, kowane asusu yana da daƙiƙa 3, kuma asusun da bai amsa ba, ko wanda shigarsa ta ƙare, ana auna shi da lambobinsa da aka ajiye. Tagar da aka ajiye wadda lokacin sake saitinta ya wuce ana ƙirga ta a matsayin 0%, don haka ba a tsallake asusun da tagarsa ta sake saitawa tun daga lokacin ba. `--auto` ba ya taɓa sabunta token. Ba kamar juyawa ta yau da kullum ba, `--auto` zai iya zaɓar asusun da ka yi amfani da shi na ƙarshe idan har yanzu shi ne ya fi sarari.

### --continue

`--continue` na Claude Code yana buɗe tattaunawa ta kwanan nan a foldar yanzu, wanda ke aiki ne kawai idan sabon asusun zai iya ganin wannan tattaunawar. Don haka da `--continue`, ccam yana duba asusan da ke rukunin `history` ɗaya da asusun da ka yi amfani da shi na ƙarshe kawai, sannan yana saka `--continue` a gaban arguments naka. Asusan da tarihinsu `local` ne ba sa raba wa kowa, don haka `--continue` daga irin wannan asusu koyaushe yana gazawa. Fara kunna rabawa da `ccam share <name> --history on` (duba [rabawa](/docs/sharing)).

## Tsaro da yanayi na musamman

- Ba a taɓa zaɓar asusan da ba su da shiga. Idan babu wani asusu da aka shiga da shi, `ccam next` kawai yana tsayawa da `no logged-in account other than <name>`.
- Da `--continue` kuma babu asusun da ke raba tarihi: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Da `--auto` kuma babu asusun da aka shiga da shi: `no logged-in account to pick`. Idan kowane asusu ya cika ko ba a san amfaninsa ba: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` yana yin buƙatun network ko da `--dry-run`. Yana karantawa ne kawai, kuma ba ya taɓa rubuta fayilolin token.
- Asusun da aka zaɓa yana farawa ta hanya ɗaya da `ccam <name>`, don haka dubawa iri ɗaya suna aiki, har da gargaɗin idan wani Claude Code da ke raba tarihi ɗaya yana aiki tuni a wannan foldar.
- Ana rubuta `last-used` ne kawai idan Claude Code ya fara da gaske. `--dry-run` yana barin sa, don haka dry run ba ya matsar da juyawar gaba.

## Masu alaƙa

- [Iyakokin amfani](/features/usage)
- [Iyakokin amfani da sauya asusai](/docs/usage-limits)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow)
