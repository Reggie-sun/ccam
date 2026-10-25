---
title: ક્વિક સ્ટાર્ટ
metaTitle: "ક્વિક સ્ટાર્ટ: Claude Code એકાઉન્ટ ઉમેરો, સ્વિચ કરો"
description: "ccam add વડે બે Claude Code એકાઉન્ટ ઉમેરો, તેમને ccwork જેવા શૉર્ટકટ આપો, તેમની વચ્ચે સ્વિચ કરો, અને તેઓ શું શેર કરે તે પસંદ કરો."
updated: 2026-10-25
order: 2
translation: machine
sourceHash: 0570ef27e4ce
---

દરેક એકાઉન્ટ `ccam add` વડે ઉમેરો, પછી `ccam <name>` વડે અથવા તમે જાતે નામ આપેલા શૉર્ટકટ, જેમ કે `ccwork`, વડે Claude Code શરૂ કરો. મોટા ભાગના લોકોને આટલું જ જોઈએ.

## 1. તમારું પહેલું એકાઉન્ટ ઉમેરો

```sh
ccam add
```

વિઝાર્ડ `personal` કે `work` જેવું નામ પૂછે છે, એકાઉન્ટ તમારાં બીજાં એકાઉન્ટ સાથે કઈ આઇટમ શેર કરે તે પૂછે છે, અને શૉર્ટકટ કમાન્ડ બનાવવો કે નહીં તે પણ પૂછે છે (પગલું 4 જુઓ). પછી તે `claude auth login` ચલાવે છે, જે સાઇન ઇન કરવા માટે બ્રાઉઝર ખોલે છે. રિમોટ મશીન પર, બ્રાઉઝર છેલ્લે લોડ ન થતા `http://localhost:…/callback` પેજ પર પહોંચી શકે છે; તે સરનામું ટર્મિનલમાં પેસ્ટ કરો અને ccam સાઇન-ઇન પૂરું કરશે. પૂરું થાય ત્યારે તમને આ દેખાય છે:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. તમારાં હાલનાં સેટિંગ્સ લઈ આવો

જો તમે પહેલાં Claude Code વાપર્યું હોય, તો તમારાં સેટિંગ્સ, `CLAUDE.md`, એક્સટેન્શન્સ અને MCP સર્વર `~/.claude` માં છે. દરેક એકાઉન્ટને તે મળે એ માટે તેમને શેર કરેલાં સેટિંગ્સમાં કૉપિ કરો:

```sh
ccam migrate
```

ccam પહેલાં બતાવે છે કે તે શું કૉપિ કરશે, અને પછી પૂછે છે. તે `~/.claude` ને ફક્ત વાંચે છે, ક્યારેય લૉગિન કૉપિ કરતું નથી, અને તમે `--history` ન ઉમેરો ત્યાં સુધી વાતચીતનો ઇતિહાસ છોડી દે છે. વિગતો માટે [હાલનું સેટઅપ કૉપિ કરવું](/docs/sharing#copy-an-existing-setup) જુઓ.

## 3. બીજું એકાઉન્ટ ઉમેરો

```sh
ccam add work
```

બ્રાઉઝરમાં બીજા એકાઉન્ટથી સાઇન ઇન કરો.

## 4. એકાઉન્ટ સ્વિચ કરો

`ccam work` Claude Code ને work એકાઉન્ટ તરીકે શરૂ કરે છે. એનાથી ટૂંકું કંઈક જોઈએ, તો એકાઉન્ટને તેનો પોતાનો શૉર્ટકટ કમાન્ડ આપો, જેમ કે `ccwork`. શેર કરવાની આઇટમ પછી વિઝાર્ડ પૂછે છે:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

`y` દબાવો, પછી નામ ટાઇપ કરો, અથવા સૂચવેલું નામ વાપરવા Enter દબાવો:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

લૉગિન પછી ccam તેની પુષ્ટિ કરે છે:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

જો તમે આ પ્રશ્ન છોડી દીધો હોય, તો ગમે ત્યારે શૉર્ટકટ ઉમેરી શકો છો:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

હવે આમાંથી કોઈ પણ કમાન્ડ Claude Code ને work એકાઉન્ટ તરીકે શરૂ કરે છે:

```sh
ccwork
ccam work
ccam 2
```

ફક્ત `ccam` ચલાવવાથી [ડેશબોર્ડ](/features/dashboard) ખૂલે છે: એરો કી કે એકાઉન્ટના નંબર વડે એકાઉન્ટ પર જાઓ, પછી Enter દબાવો.

શૉર્ટકટ કે એકાઉન્ટના નામ પછીનું બધું કોઈ ફેરફાર વગર Claude Code ને જાય છે, તેથી તમે જાણીતા ઓપ્શન્સ જ વાપરતા રહો છો:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. શું શેર કરવું તે પસંદ કરો

ડિફૉલ્ટ રીતે MCP સર્વર, `settings.json`, ગ્લોબલ `CLAUDE.md`, અને agents, commands, skills અને plugins `default` ગ્રુપ દ્વારા એકાઉન્ટ વચ્ચે શેર થાય છે. વાતચીતનો ઇતિહાસ શેર થતો નથી. જો તમે એકાઉન્ટ વચ્ચે વાતચીત ચાલુ રાખવા માગતા હો, તો તેને ગ્રુપમાં ઉમેરો:

```sh
ccam group set default --add history
```

કંઈ પણ બદલતાં પહેલાં ccam બતાવે છે કે તે દરેક એકાઉન્ટ માટે શું ખસેડશે. વિગતો માટે [ગ્રુપ વડે શેરિંગ](/features/share) જુઓ.

## 6. જ્યારે એકાઉન્ટ તેની મર્યાદા સુધી પહોંચે

```sh
ccam usage
ccam next --continue
```

`ccam usage` દરેક એકાઉન્ટનો 5 કલાકનો અને સાપ્તાહિક વપરાશ બતાવે છે. `ccam next --continue` એ જ વાતચીત આગલા એકાઉન્ટમાં ખોલે છે. scripts માં `ccam --fallback work -p "…"` જાતે સ્વિચ કરે છે. [વપરાશ મર્યાદા](/docs/usage-limits) જુઓ.

**Team અને Enterprise સીટ.** આપમેળે થતી પસંદગી (`ccam next`, `--fallback`, `ccam pick`, ડેશબોર્ડની `n` કી, અને તમે ચાલુ કરો ત્યારે [`auto_switch`](/features/next#auto-switch)) એક જ Claude Team કે Enterprise સંસ્થાની સીટો વચ્ચે પણ સ્વિચ કરે છે. Anthropic નું સહાય કેન્દ્ર કહે છે કે "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team પ્લાનની દરેક સીટ નામથી નક્કી કરેલી એક વ્યક્તિને ફાળવવામાં આવે છે, અને ફક્ત એ વ્યક્તિ જ તેનો ઉપયોગ કરી શકે છે) અને "Users can't route their work through another user's account to use that user's remaining usage" (વપરાશકર્તાઓ બીજા વપરાશકર્તાના એકાઉન્ટ મારફતે પોતાનું કામ ચલાવીને એ વપરાશકર્તાનો બાકી રહેલો વપરાશ વાપરી શકતા નથી) ([Team પ્લાન](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise પ્લાન](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) માટે પણ એવું જ). તેથી સીટો વચ્ચે આપમેળે થતી પસંદગી વાપરવાથી Anthropic ની શરતોનું ઉલ્લંઘન થઈ શકે છે. સીટો વચ્ચે તેને વાપરવું કે નહીં એ તમારો નિર્ણય છે: તેના પર આધાર રાખતાં પહેલાં તમારી સંસ્થા અને Anthropic પાસેથી ખાતરી કરી લો. અંગત પ્લાન (Pro, Max) [વારંવાર પૂછાતા પ્રશ્નોમાંની શરતો અંગેની નોંધ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) માં આવરી લેવાયા છે.

## પહેલેથી ~/.claude વાપરો છો?

ફરીથી સાઇન ઇન કરવાને બદલે તેને રજિસ્ટર કરો:

```sh
ccam import
```

ફોલ્ડર જ્યાં છે ત્યાં જ રહે છે, અને સાદું `claude` પહેલાંની જેમ કામ કરતું રહે છે.
