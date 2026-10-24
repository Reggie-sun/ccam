---
title: Zaɓar asusu don rubutun umarni
description: ccam pick yana buga asusun da ya fi sarari don rubutun umarni, yana ba workers da aka fara tare asusai daban-daban, kuma yana fita da 75 idan babu sarari.
updated: 2026-10-24
order: 12
translation: machine
sourceHash: a8baeb64de99
---

`ccam pick` yana buga sunan asusun da aka shiga da shi wanda ya fi sarari, don rubutun umarni (script) ya fara Claude Code da shi. Workers da aka fara a lokaci ɗaya kowannensu yana samun asusu daban muddin akwai isassun asusai masu sarari, don haka rukunin ayyukan `claude -p` yana rarrabuwa a kan asusanka maimakon ya shanye ɗaya.

**Kujeru (seats) na Team da Enterprise.** Zaɓin atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba [bayanin sharuɗɗa a cikin Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Amfani

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Exit codes: `0` tare da suna a standard output, `75` (`EX_TEMPFAIL`) idan babu asusun da ke da sarari ko babu wanda aka shiga da shi, `1` don wasu kurakurai kamar flag mara kyau.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--json` | Buga asusun, amfaninsa na awa 5 da na mako, inda lambobin suka fito da dalilin da ya sa aka zaɓe shi | a kashe |
| `--model <model>` | Ƙirga tagar mako ta wannan model ɗin ma, kamar `opus` | babu |
| `--max-usage <1-100>` | Fifita asusan da ke ƙasa da wannan amfani | `auto_max_usage`, in ba haka ba 100 |
| `--strategy headroom\|expiring\|balance` | Mafi ƙarancin amfanin awa 5, sake saitin awa 5 mafi kusa, ko kason mako da zai fara ɓata ba tare da an yi amfani da shi ba | `auto_strategy`, in ba haka ba `headroom` |
| `--cross-org` | Zaɓi asusan wasu ƙungiyoyi (organizations) ma | `cross_org`, in ba haka ba a kashe |
| `--wait[=<duration>]` | Idan babu asusun da ke da sarari, jira na farko ya sake saitawa, har zuwa iyakar wannan tsawon. Ƙimar tana zuwa bayan `=`: ana ƙin `--wait 30m` da shawarar `--wait takes its value after '=': --wait=30m` | a kashe; `5h` idan an ba shi shi kaɗai |
| `--dry-run` | Zaɓa ba tare da rubuta zaɓin ba | a kashe |

## Yadda yake aiki

1. ccam yana lissafa asusan da aka shiga da su (asusan API key da shigar Console ba su da iyaka kuma ana barin su a waje) a cikin ƙungiyar asusun da wannan foldar za ta fara (duba [ƙungiyoyi](/features/next#auto)); ana barin foldar da ta shiga asusun Claude iri ɗaya da wannan a waje ita ma.
2. Yana ɗaukar kulle (`~/.local/state/ccam/picks.lock`), don zaɓuɓɓukan da ke gudana a lokaci ɗaya su faru ɗaya bayan ɗaya.
3. Yana neman amfanin kowane asusu a lokaci ɗaya, kamar yadda `ccam usage` ke yi, tare da cache na daƙiƙa 60. Tagar da Claude Code ya ba da rahoton ta cika ana ƙirga ta a matsayin 100% har zuwa sake saitinta ([rikodin iyaka](/features/usage#the-limit-record)), don haka asusun da ya kai iyakarsa yanzu yana zama a waje ko da lambobinsa da aka ajiye suna nuna sarari.
4. Ana fara barin asusan da aka zaɓa a cikin minti 2 da suka gabata (`~/.local/state/ccam/picks/<name>`; zaɓin da kwanan watansa ya wuce lokacin agogo ba ya ƙirguwa) a waje, kuma [ƙa'idodin zaɓi](/features/next#auto) na `ccam next --auto` suna zaɓa a cikin sauran: waɗanda ke ƙasa da `auto_max_usage` da farko, bisa dabarar, kuma ba a taɓa zaɓar asusun da ya kai 100% ba.
5. Idan an zaɓi kowane asusu mai sarari kwanan nan, ƙa'idodin suna sake zaɓa a cikin dukansu, don haka worker na huɗu a kan asusai uku yana raba ɗaya.
6. Ana rubuta zaɓin kuma a saki kullen, sannan a buga sunan.

Tagar minti 2 tana da muhimmanci domin amfanin sabon aiki yana ɗaukar ɗan lokaci kafin ya bayyana a lambobin: in babu ita, workers biyar da aka fara tare duk za su ga asusu ɗaya a matsayin wanda ya fi sarari.

## Tsaro da yanayi na musamman

- `ccam pick` yana karanta amfani ne kawai kuma yana rubuta lokacin zaɓi; ba ya taɓa sabunta ko karanta token fiye da abin da `ccam usage` ke yi.
- Zaɓi ba ya keɓe asusun. Rubutun umarni biyu za su iya fara a asusu ɗaya bayan minti 2; yi amfani da `--fallback` don farfaɗowa idan wannan asusun ya kai iyakarsa.
- Jira da `--wait` ba ya riƙe kulle, don haka sauran zaɓuɓɓuka suna ci gaba a halin yanzu.
- `ccam pick` yana watsi da [`ccam folder`](/features/folder): yana zaɓa bisa amfani kawai.

## Masu alaƙa

- [Sauyawa zuwa asusu na gaba](/features/next)
- [Sake gudanarwa a iyaka](/features/next#run-again-at-the-limit)
- [Iyakokin amfani](/features/usage)
- [Asusan API key](/features/api)
