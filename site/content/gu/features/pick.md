---
title: scripts માટે એકાઉન્ટ પસંદ કરવું
description: ccam pick script માટે સૌથી વધુ જગ્યાવાળું એકાઉન્ટ છાપે છે, એકસાથે શરૂ થતા workers ને જુદાં એકાઉન્ટ આપે છે, અને જગ્યા ન હોય ત્યારે 75 સાથે બહાર નીકળે છે.
updated: 2026-10-24
order: 12
translation: machine
sourceHash: a8baeb64de99
---

`ccam pick` સૌથી વધુ જગ્યાવાળા લૉગિન થયેલા એકાઉન્ટનું નામ છાપે છે, જેથી script તેનાથી Claude Code શરૂ કરી શકે. જ્યાં સુધી જગ્યાવાળાં પૂરતાં એકાઉન્ટ હોય, ત્યાં સુધી એકસાથે શરૂ થતા દરેક worker ને જુદું એકાઉન્ટ મળે છે, તેથી `claude -p` jobs નો સમૂહ એક જ એકાઉન્ટ ખાલી કરવાને બદલે તમારાં એકાઉન્ટમાં વહેંચાય છે.

**Team અને Enterprise સીટ.** આપમેળે થતી પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે. Anthropic નું સહાય કેન્દ્ર કહે છે કે "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team પ્લાનની દરેક સીટ નામથી નક્કી કરેલી એક વ્યક્તિને ફાળવવામાં આવે છે, અને ફક્ત એ વ્યક્તિ જ તેનો ઉપયોગ કરી શકે છે) અને "Users can't route their work through another user's account to use that user's remaining usage" (વપરાશકર્તાઓ બીજા વપરાશકર્તાના એકાઉન્ટ મારફતે પોતાનું કામ ચલાવીને એ વપરાશકર્તાનો બાકી રહેલો વપરાશ વાપરી શકતા નથી) ([Team પ્લાન](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise પ્લાન](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) માટે પણ એવું જ). તેથી સીટો વચ્ચે આપમેળે થતી પસંદગી વાપરવાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો. અંગત પ્લાન (Pro, Max) [વારંવાર પૂછાતા પ્રશ્નોમાંની શરતો અંગેની નોંધ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) માં આવરી લેવાયા છે.

## ઉપયોગ

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

Exit code: standard output પર નામ સાથે `0`, કોઈ એકાઉન્ટમાં જગ્યા ન હોય કે કોઈ લૉગિન ન હોય ત્યારે `75` (`EX_TEMPFAIL`), અને ખોટા flag જેવી બીજી ભૂલો માટે `1`.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--json` | એકાઉન્ટ, તેનો 5 કલાકનો અને સાપ્તાહિક વપરાશ, આંકડા ક્યાંથી આવે છે અને તે શા માટે પસંદ થયું તે છાપે છે | off |
| `--model <model>` | એ મૉડલની સાપ્તાહિક વિન્ડો પણ ગણે છે, જેમ કે `opus` | કોઈ નહીં |
| `--max-usage <1-100>` | આ વપરાશથી નીચેનાં એકાઉન્ટને પ્રાથમિકતા આપે છે | `auto_max_usage`, નહીં તો 100 |
| `--strategy headroom\|expiring\|balance` | સૌથી ઓછો 5 કલાકનો વપરાશ, સૌથી વહેલો 5 કલાકનો રીસેટ, અથવા જે સાપ્તાહિક quota સૌથી પહેલાં વણવપરાયેલો રહી જાય | `auto_strategy`, નહીં તો `headroom` |
| `--cross-org` | બીજી સંસ્થાનાં એકાઉન્ટ પણ પસંદ કરે છે | `cross_org`, નહીં તો off |
| `--wait[=<duration>]` | કોઈ એકાઉન્ટમાં જગ્યા ન હોય ત્યારે, પહેલું એકાઉન્ટ રીસેટ થાય તેની રાહ જુએ છે, વધુમાં વધુ આટલો સમય. મૂલ્ય `=` પછી આવે છે: `--wait 30m` નો `--wait takes its value after '=': --wait=30m` સૂચના સાથે ઇનકાર થાય છે | off; એકલું આપો તો `5h` |
| `--dry-run` | પસંદગી નોંધ્યા વગર પસંદ કરે છે | off |

## તે કેવી રીતે કામ કરે છે

1. આ ફોલ્ડર જે એકાઉન્ટ શરૂ કરે તેની સંસ્થામાં ccam લૉગિન થયેલાં એકાઉન્ટની યાદી બનાવે છે (API key એકાઉન્ટ અને Console લૉગિનને કોઈ મર્યાદા નથી, તેથી તે બહાર રહે છે; જુઓ [સંસ્થાઓ](/features/next#auto)); એ જ Claude એકાઉન્ટમાં લૉગિન થયેલું ફોલ્ડર પણ બહાર રહે છે.
2. તે એક lock (`~/.local/state/ccam/picks.lock`) લે છે, તેથી એકસાથે ચાલતી પસંદગીઓ એક પછી એક થાય છે.
3. તે `ccam usage` ની જેમ જ, 60 સેકન્ડના કૅશ સાથે, બધાં એકાઉન્ટનો વપરાશ એકસાથે જુએ છે. જે વિન્ડો ભરાઈ ગઈ હોવાનું Claude Code એ જણાવ્યું હોય તે તેના રીસેટ સુધી 100% ગણાય છે ([મર્યાદા રેકોર્ડ](/features/usage#the-limit-record)), તેથી હમણાં જ મર્યાદાએ પહોંચેલું એકાઉન્ટ, તેના સેવ કરેલા આંકડા જગ્યા બતાવતા હોય ત્યારે પણ, બહાર રહે છે.
4. છેલ્લી 2 મિનિટમાં પસંદ થયેલાં એકાઉન્ટ (`~/.local/state/ccam/picks/<name>`; ઘડિયાળના સમય પછીની તારીખવાળી પસંદગી ગણાતી નથી) પહેલાં બહાર રખાય છે, અને બાકીનામાંથી `ccam next --auto` ના [પસંદગીના નિયમો](/features/next#auto) પસંદ કરે છે: પહેલાં `auto_max_usage` થી નીચેનાં, strategy પ્રમાણે, અને 100% પરનું એકાઉન્ટ ક્યારેય નહીં.
5. જો જગ્યાવાળું દરેક એકાઉન્ટ તાજેતરમાં પસંદ થયું હોય, તો નિયમો ફરી એ બધાંમાંથી પસંદ કરે છે, તેથી ત્રણ એકાઉન્ટ પર ચોથો worker કોઈ એક એકાઉન્ટ વહેંચે છે.
6. પસંદગી નોંધાય છે અને lock છૂટે છે, પછી નામ છપાય છે.

2 મિનિટની વિન્ડો મહત્ત્વની છે, કારણ કે નવા job નો વપરાશ આંકડામાં દેખાતાં થોડો સમય લાગે છે: તેના વગર, એકસાથે શરૂ થયેલા પાંચ workers બધા એક જ એકાઉન્ટને સૌથી વધુ જગ્યાવાળું ગણત.

## સલામતી અને ખાસ કિસ્સા

- `ccam pick` ફક્ત વપરાશ વાંચે છે અને પસંદગીનો સમય લખે છે; `ccam usage` જે કરે તેથી વધુ તે ક્યારેય tokens રિન્યૂ કરતું કે વાંચતું નથી.
- પસંદગી એકાઉન્ટને અનામત રાખતી નથી. 2 મિનિટ પછી બે scripts હજી પણ એક જ એકાઉન્ટ પર શરૂ થઈ શકે છે; એ એકાઉન્ટ મર્યાદાએ પહોંચે ત્યારે સંભાળવા `--fallback` વાપરો.
- `--wait` સાથે રાહ જોતી વખતે કોઈ lock રહેતું નથી, તેથી એ દરમિયાન બીજી પસંદગીઓ ચાલુ રહે છે.
- `ccam pick` [`ccam folder`](/features/folder) ને અવગણે છે: તે ફક્ત વપરાશ પરથી પસંદ કરે છે.

## સંબંધિત

- [આગલા એકાઉન્ટ પર સ્વિચ કરવું](/features/next)
- [મર્યાદાએ ફરી ચલાવવું](/features/next#run-again-at-the-limit)
- [વપરાશ મર્યાદા](/features/usage)
- [API key એકાઉન્ટ](/features/api)
