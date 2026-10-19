---
title: ક્વિક સ્ટાર્ટ
description: ccam add વડે બે Claude Code એકાઉન્ટ ઉમેરો, claude1, claude2 કે ccam <name> વડે તેમની વચ્ચે સ્વિચ કરો અને તેઓ શું શેર કરે તે પસંદ કરો.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

દરેક એકાઉન્ટ `ccam add` વડે ઉમેરો, પછી `claude1`, `claude2` કે `ccam <name>` વડે Claude Code શરૂ કરો. મોટા ભાગના લોકોને આટલું જ જોઈએ.

## 1. તમારું પહેલું એકાઉન્ટ ઉમેરો

```sh
ccam add
```

વિઝાર્ડ `personal` કે `work` જેવું નામ પૂછે છે, અને એકાઉન્ટ તમારાં બીજાં એકાઉન્ટ સાથે કઈ વસ્તુઓ શેર કરે તે પૂછે છે. પછી તે `claude auth login` ચલાવે છે, જે બ્રાઉઝર ખોલે છે જેથી તમે સાઇન ઇન કરી શકો. તે પૂરું થાય ત્યારે તમને આ દેખાય છે:

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam પોતાના સંદેશા અંગ્રેજીમાં છાપે છે; કમાન્ડ અને ઓપ્શન દરેક ભાષામાં સરખા છે.

## 2. બીજું એકાઉન્ટ ઉમેરો

```sh
ccam add work
```

બ્રાઉઝરમાં બીજા એકાઉન્ટથી સાઇન ઇન કરો. તે `claude2` બને છે.

## 3. સ્વિચ કરો

આમાંથી કોઈ પણ કમાન્ડ work એકાઉન્ટ તરીકે Claude Code શરૂ કરે છે:

```sh
claude2
ccam work
ccam 2
```

એકાઉન્ટના નામ પછીનું બધું કોઈ ફેરફાર વગર Claude Code ને જાય છે, તેથી તમે જાણીતા ઓપ્શન વાપરવાનું ચાલુ રાખી શકો છો:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

ફક્ત `ccam` ચલાવવાથી એક પસંદગી મેનૂ ખૂલે છે: એરો કી અને Enter વાપરો, અથવા એકાઉન્ટનો નંબર દબાવો.

## 4. શું શેર કરવું તે પસંદ કરો

ડિફૉલ્ટ રીતે MCP સર્વર, `settings.json`, ગ્લોબલ `CLAUDE.md`, અને agents, commands, skills અને plugins એકાઉન્ટ વચ્ચે શેર થાય છે. વાતચીતનો ઇતિહાસ શેર થતો નથી. જો તમે એકાઉન્ટ વચ્ચે વાતચીત ચાલુ રાખવા માગતા હો, તો બંને એકાઉન્ટ માટે તે ચાલુ કરો:

```sh
ccam share personal --history on
ccam share work --history on
```

કંઈ પણ બદલતાં પહેલાં ccam બતાવે છે કે તે શું ખસેડશે. વિગતો માટે [શેરિંગ](/docs/sharing) જુઓ.

## 5. જ્યારે એકાઉન્ટ તેની મર્યાદા સુધી પહોંચે

```sh
ccam usage
ccam next --continue
```

`ccam usage` દરેક એકાઉન્ટનો 5 કલાકનો અને સાપ્તાહિક વપરાશ બતાવે છે. `ccam next --continue` એ જ વાતચીત આગલા એકાઉન્ટમાં ખોલે છે. [વપરાશ મર્યાદા](/docs/usage-limits) જુઓ.

## પહેલેથી ~/.claude વાપરો છો?

ફરીથી સાઇન ઇન કરવાને બદલે તેને રજિસ્ટર કરો:

```sh
ccam import
```

ફોલ્ડર જ્યાં છે ત્યાં જ રહે છે, અને સાદું `claude` પહેલાંની જેમ કામ કરતું રહે છે.
