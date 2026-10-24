---
title: Claude Code વપરાશ મર્યાદા આવે ત્યારે શું કરવું
description: રીસેટની રાહ જુઓ અથવા એ જ વાતચીત તરત બીજા Claude Code એકાઉન્ટ પર ચાલુ રાખો. 5 કલાકની અને સાપ્તાહિક મર્યાદા કેવી રીતે કામ કરે છે, પગલું દર પગલું.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 049a46684c7b
---

જ્યારે કોઈ એકાઉન્ટ તેની વપરાશ મર્યાદા (જેને ઘણી વાર rate limit કહેવાય છે) સુધી પહોંચે, ત્યારે તમારી પાસે બે વિકલ્પ હોય છે: મર્યાદા રીસેટ થાય ત્યાં સુધી રાહ જુઓ, અથવા એ જ વાતચીત અત્યારે જ તમારા બીજા એકાઉન્ટ પર ચાલુ રાખો. ccam સાથે બીજો વિકલ્પ ફક્ત એક કમાન્ડ છે: `ccam next --auto --continue`.

## Claude Code માં કઈ મર્યાદાઓ હોય છે?

Pro અને Max પ્લાન પર Claude Code માં [પાંચ કલાકની સેશન મર્યાદા અને સાપ્તાહિક મર્યાદા](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) હોય છે. બંને મર્યાદા વેબ અને ઍપમાંના Claude સાથે વહેંચાયેલી છે, અને વપરાશ એકસાથે બંનેમાં ગણાય છે. જ્યારે કોઈ એક પૂરી થાય, ત્યારે Claude Code "You've hit your session limit · resets 3:45pm" જેવો [સંદેશ બતાવે છે](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit). મૉડલ બદલવાથી કંઈ નહીં વળે, કારણ કે મર્યાદા બધાં મૉડલ પર લાગુ પડે છે.

## વિકલ્પ 1: રીસેટની રાહ જુઓ

જો તમારી પાસે એક જ એકાઉન્ટ હોય, તો રાહ જુઓ. Claude Code ની તાજેતરની આવૃત્તિઓ [ખુલ્લા સેશનમાં રાહ જુએ છે](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) અને મર્યાદા રીસેટ થતાં જ જાતે કામ આગળ વધારે છે. Claude Code ની અંદર `/usage` બતાવે છે કે તમે તમારા પ્લાનની મર્યાદામાંથી કેટલું વાપર્યું છે.

## વિકલ્પ 2: અત્યારે જ બીજા એકાઉન્ટ પર ચાલુ રાખો

જો તમારી પાસે બીજું એકાઉન્ટ હોય, તો તમે એ જ વાતચીત તરત આગળ ચલાવી શકો છો. એક વાર આ સેટ કરો:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

હવે તમે કામ કરતા હો ત્યારે સ્ટેટસ લાઇન એકાઉન્ટ અને તેનો 5 કલાકનો વપરાશ બતાવે છે, તેથી મર્યાદા અચાનક આવતી નથી.

જ્યારે મર્યાદા આવી જાય, ત્યારે Claude Code માંથી બહાર નીકળો (`/exit`), પછી:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

સૌથી વધુ જગ્યાવાળા એકાઉન્ટમાં વાતચીત ખોલો:

```sh
ccam next --auto --continue
```

Claude Code હાલના ફોલ્ડરની સૌથી તાજેતરની વાતચીત સાથે `work` તરીકે શરૂ થાય છે, અત્યાર સુધી કહેલી દરેક વાત સહિત. ઇતિહાસ શેર કરવાથી (ઉપરનું `ccam group set default --add history`) તે એક જ વાતચીત રહે છે; તેના વગર Claude Code વાતચીતની નકલ `work` માં કરે છે, અને ccam તમને જણાવે છે કે એ પહેલો જવાબ કેટલા tokens ફરી વાંચે છે.

Claude Code માં [ccam plugin](/docs/install#claude-code-plugin) સાથે, મર્યાદાએ અટકેલું સેશન એક સૂચના બતાવે છે અને ભરાયેલી વિન્ડો નોંધે છે, જેથી `ccam next --auto` તેના રીસેટ પહેલાં એ એકાઉન્ટ ફરી પસંદ ન કરે, અને `/ccam:switch` (અથવા `/ccam:switch work`) `ccam next --continue` માટે વાતચીતની નોંધ કરે છે, જેથી અનુમાન કરવા માટે કંઈ બાકી ન રહે.

**Team અને Enterprise સીટ.** આપમેળે થતી પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે. Anthropic નું સહાય કેન્દ્ર કહે છે કે "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team પ્લાનની દરેક સીટ નામથી નક્કી કરેલી એક વ્યક્તિને ફાળવવામાં આવે છે, અને ફક્ત એ વ્યક્તિ જ તેનો ઉપયોગ કરી શકે છે) અને "Users can't route their work through another user's account to use that user's remaining usage" (વપરાશકર્તાઓ બીજા વપરાશકર્તાના એકાઉન્ટ મારફતે પોતાનું કામ ચલાવીને એ વપરાશકર્તાનો બાકી રહેલો વપરાશ વાપરી શકતા નથી) ([Team પ્લાન](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise પ્લાન](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) માટે પણ એવું જ). તેથી સીટો વચ્ચે આપમેળે થતી પસંદગી વાપરવાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. તેના પર આધાર રાખતાં પહેલાં તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો. અંગત પ્લાન (Pro, Max) [વારંવાર પૂછાતા પ્રશ્નોમાંની શરતો અંગેની નોંધ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) માં આવરી લેવાયા છે.

## શું ccam આપમેળે એકાઉન્ટ સ્વિચ કરે છે?

ના. તમે કામ કરતા હો ત્યારે ccam ચાલુ હોતું નથી: તે Claude Code શરૂ કરે છે અને પોતાની જગ્યા તેને આપી દે છે. ક્યારે સ્વિચ કરવું તે તમે નક્કી કરો છો, અને `ccam next --auto` સૌથી ઓછા 5 કલાકના વપરાશવાળું એકાઉન્ટ પસંદ કરે છે અને 100% પર પહોંચેલાં, લૉગિન ન હોય તેવાં કે બીજી સંસ્થાનાં એકાઉન્ટ છોડી દે છે. `config.toml` માં `relay = true` સાથે, ccam તેના બદલે સેશનની બાજુમાં રાહ જુએ છે, તેથી `/ccam:switch` જાતે આગળ વધે છે, સેશનના ઓપ્શન, permission mode અને effort level જાળવીને. કેટલાંક બીજાં ટૂલ તેના બદલે ચાલુ સેશનનું લૉગિન બદલે છે; [ccam અને બીજાં એકાઉન્ટ સ્વિચર](/compare/account-switchers) આ તફાવત સમજાવે છે.

## ધ્યાનમાં રાખવાની બાબતો

- `ccam next --continue` એ જ ફોલ્ડરમાં ચલાવો જ્યાં તમે કામ કરતા હતા; વાતચીત દરેક પ્રોજેક્ટ ફોલ્ડર પ્રમાણે સંગ્રહાય છે.
- એક જ વાતચીત પર એક જ સમયે બે એકાઉન્ટ ચલાવવાનું ટાળો. જો તમે પ્રયત્ન કરો તો ccam તમને ચેતવે છે.
- શેર કરેલા ઇતિહાસ વગર, વાતચીતની નકલ બીજા એકાઉન્ટમાં થાય છે અને બંને નકલ પોતપોતાના માર્ગે આગળ વધે છે; `/rewind` જેમાંથી કોડ પાછો લાવે છે તે checkpoints પહેલા એકાઉન્ટ પાસે જ રહે છે.
- નવા એકાઉન્ટ પરનો પહેલો સંદેશ સામાન્ય કરતાં તેની મર્યાદામાંથી વધુ વાપરી શકે છે, કારણ કે એ એકાઉન્ટે હજી વાતચીત કૅશ કરી નથી.
- તમે આ રીતે એકાઉન્ટ વચ્ચે સ્વિચ કરી શકો કે નહીં, તે દરેક એકાઉન્ટની શરતો પર આધાર રાખે છે.

`ccam usage` અને `ccam next` ના દરેક ઓપ્શન માટે [વપરાશ મર્યાદા](/docs/usage-limits) જુઓ.
