---
title: Kuchagua akaunti kwa ajili ya script
description: ccam pick huchapisha akaunti yenye nafasi kubwa zaidi kwa script, huwapa workers walioanzishwa pamoja akaunti tofauti, na hutoka na 75 isipokuwa na nafasi.
updated: 2026-10-24
order: 12
translation: machine
sourceHash: "a8baeb64de99"
---

`ccam pick` huchapisha jina la akaunti iliyoingia yenye nafasi kubwa zaidi, ili script ianzishe Claude Code nayo. Workers walioanzishwa kwa wakati mmoja kila mmoja hupata akaunti tofauti maadamu kuna akaunti za kutosha zenye nafasi, kwa hiyo kundi la kazi za `claude -p` husambaa kwenye akaunti zako badala ya kumaliza moja.

**Viti (seats) vya Team na Enterprise.** Chaguo za kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise. Kituo cha usaidizi cha Anthropic kinasema "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (kila kiti katika mpango wa Team hutolewa kwa mtu mmoja aliyetajwa kwa jina, na ni mtu huyo pekee anayeruhusiwa kukitumia) na "Users can't route their work through another user's account to use that user's remaining usage" (watumiaji hawaruhusiwi kuelekeza kazi yao kupitia akaunti ya mtumiaji mwingine ili kutumia kiasi cha matumizi kilichobaki cha mtumiaji huyo) ([mipango ya Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); vivyo hivyo kwa [mipango ya Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kwa hiyo, kutumia chaguo za kiotomatiki kati ya viti kunaweza kukiuka masharti ya Anthropic. Wasiliana na shirika lako na Anthropic kabla ya kuzitegemea. Mipango ya binafsi (Pro, Max) inaelezwa katika [dokezo la masharti kwenye Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Matumizi

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

Exit codes: `0` pamoja na jina kwenye standard output, `75` (`EX_TEMPFAIL`) wakati hakuna akaunti yenye nafasi au hakuna iliyoingia, `1` kwa makosa mengine kama flag isiyo sahihi.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--json` | Huchapisha akaunti, matumizi yake ya saa 5 na ya wiki, chanzo cha namba na kwa nini ilichaguliwa | imezimwa |
| `--model <model>` | Huhesabu pia kipindi cha wiki cha modeli hiyo, kama `opus` | hakuna |
| `--max-usage <1-100>` | Hupendelea akaunti zilizo chini ya matumizi haya | `auto_max_usage`, vinginevyo 100 |
| `--strategy headroom\|expiring\|balance` | Matumizi ya chini zaidi ya saa 5, kipindi cha saa 5 kinachoanza upya mapema zaidi, au kiasi cha wiki ambacho kingebaki bila kutumika kwanza | `auto_strategy`, vinginevyo `headroom` |
| `--cross-org` | Huchagua pia akaunti za mashirika mengine | `cross_org`, vinginevyo imezimwa |
| `--wait[=<duration>]` | Hakuna akaunti yenye nafasi, husubiri ya kwanza ianze upya, si zaidi ya muda huu. Thamani huwekwa baada ya `=`: `--wait 30m` hukataliwa pamoja na kidokezo `--wait takes its value after '=': --wait=30m` | imezimwa; `5h` ikitolewa peke yake |
| `--dry-run` | Huchagua bila kurekodi uchaguzi | imezimwa |

## Jinsi inavyofanya kazi

1. ccam huorodhesha akaunti zilizoingia (akaunti za API key na kuingia kwa Console hazina vikomo na huachwa nje) katika shirika la akaunti ambayo folda hii ingeanzisha (tazama [mashirika](/features/next#auto)); folda iliyoingia kwenye akaunti ile ile ya Claude kama hiyo huachwa nje pia.
2. Huchukua kufuli (`~/.local/state/ccam/picks.lock`), ili uchaguzi unaoendeshwa kwa wakati mmoja ufanyike mmoja baada ya mwingine.
3. Huomba matumizi ya kila akaunti kwa wakati mmoja, kwa njia ile ile kama `ccam usage`, pamoja na akiba ya sekunde 60. Kipindi ambacho Claude Code iliripoti kuwa kimejaa huhesabiwa kama 100% hadi kianze upya ([rekodi ya kikomo](/features/usage#the-limit-record)), kwa hiyo akaunti iliyofikia kikomo chake sasa hivi hubaki nje hata wakati namba zake zilizohifadhiwa zinaonyesha nafasi.
4. Akaunti zilizochaguliwa ndani ya dakika 2 zilizopita (`~/.local/state/ccam/picks/<name>`; uchaguzi wenye tarehe iliyo baada ya muda wa saa ya kompyuta hauhesabiwi) huachwa nje kwanza, na [kanuni za uchaguzi](/features/next#auto) za `ccam next --auto` huchagua kati ya zilizobaki: zilizo chini ya `auto_max_usage` kwanza, kwa mkakati, na kamwe si akaunti iliyo kwenye 100%.
5. Ikiwa kila akaunti yenye nafasi ilichaguliwa hivi karibuni, kanuni huchagua tena kati ya zote, kwa hiyo worker wa nne kwenye akaunti tatu hushiriki moja.
6. Uchaguzi hurekodiwa na kufuli huachiliwa, kisha jina huchapishwa.

Kipindi cha dakika 2 ni muhimu kwa sababu matumizi ya kazi mpya huchukua muda kuonekana katika namba: bila hicho, workers watano walioanzishwa pamoja wote wangeona akaunti ile ile kama yenye nafasi kubwa zaidi.

## Usalama na hali maalum

- `ccam pick` husoma matumizi tu na huandika muda wa uchaguzi; kamwe haifanyi upya wala kusoma tokeni zaidi ya kile `ccam usage` hufanya.
- Uchaguzi hauihifadhi akaunti. Script mbili bado zinaweza kuanza kwenye akaunti ile ile baada ya dakika 2; tumia `--fallback` kujinasua akaunti hiyo inapofikia kikomo chake.
- Kusubiri kwa `--wait` hakushiki kufuli, kwa hiyo uchaguzi mwingine huendelea wakati huo.
- `ccam pick` hupuuza [`ccam folder`](/features/folder): huchagua kwa matumizi pekee.

## Zinazohusiana

- [Kubadili kwenda akaunti inayofuata](/features/next)
- [Endesha tena kwenye kikomo](/features/next#run-again-at-the-limit)
- [Vikomo vya matumizi](/features/usage)
- [Akaunti za API key](/features/api)
