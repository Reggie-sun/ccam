---
title: CLAUDE_CONFIG_DIR સમજૂતી
description: Claude Code માં CLAUDE_CONFIG_DIR શું કરે છે, ફોલ્ડરમાં શું સંગ્રહાય છે, macOS લૉગિન ક્યાં રાખે છે અને એકાઉન્ટ વચ્ચે કઈ ફાઇલો શેર થઈ શકે.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: dc2fe42e11e1
---

`CLAUDE_CONFIG_DIR` Claude Code ને જણાવે છે કે `~/.claude` ને બદલે કયું ફોલ્ડર વાપરવું. એકાઉન્ટને લગતું બધું ત્યાં રહે છે: લૉગિન, સેટિંગ્સ, MCP સર્વર, plugins અને વાતચીતનો ઇતિહાસ. તેને જુદાં ફોલ્ડર તરફ નિર્દેશિત કરવાથી તમને અલગ એકાઉન્ટ મળે છે.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## શું CLAUDE_CONFIG_DIR સત્તાવાર રીતે સમર્થિત છે?

હા. Claude Code તેને પોતાના [environment variables](https://code.claude.com/docs/en/env-vars) માં કૉન્ફિગરેશન ફોલ્ડર બદલવાની રીત તરીકે ગણાવે છે, જે "ઘણાં એકાઉન્ટ સાથે-સાથે ચલાવવા માટે ઉપયોગી" છે, અને તેની [authentication માર્ગદર્શિકા](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) ઘણાં એકાઉન્ટથી લૉગિન કરવા માટે તેનો ઉપયોગ કરે છે. તમે તેને તમારા શેલમાં કે user settings માં સેટ કરી શકો છો, પણ પ્રોજેક્ટનાં settings માં નહીં.

## ફોલ્ડરની અંદર શું હોય છે?

Claude Code 2.1.309 સાથે તપાસેલું:

| પાથ | સામગ્રી | શેર કરી શકાય |
| --- | --- | --- |
| `.credentials.json` | લૉગિન ટોકન (Linux, Windows), MCP OAuth ટોકન | ના |
| `.claude.json` | એકાઉન્ટની વિગતો, user MCP સર્વર, દરેક પ્રોજેક્ટની સ્થિતિ | ના (MCP સર્વર અલગથી) |
| `settings.json` | પરવાનગીઓ, hooks, મૉડલ, સ્ટેટસ લાઇન, plugins | હા |
| `CLAUDE.md` | ગ્લોબલ સૂચનાઓ | હા |
| `agents/`, `commands/`, `skills/`, `plugins/` | એક્સટેન્શન્સ | હા |
| `projects/` | `.jsonl` ફાઇલો તરીકે વાતચીત, દરેક પ્રોજેક્ટની મેમરી | હા |
| `history.jsonl` | પ્રૉમ્પ્ટનો ઇતિહાસ | હા |
| `sessions/`, `shell-snapshots/`, `cache/` અને બીજાં | રનટાઇમ સ્થિતિ | ના |

## લોકો ક્યાં અટવાય છે?

**`CLAUDE_CONFIG_DIR` વગર, `.claude.json` એ `~/.claude` ની અંદર હોતી નથી.** ત્યારે Claude Code તેને ફોલ્ડરની બાજુમાં `~/.claude.json` પર રાખે છે. જો તમે `CLAUDE_CONFIG_DIR=~/.claude` વાપરવાનું શરૂ કરો, તો Claude Code તેના બદલે `~/.claude/.claude.json` શોધે છે અને તમારું એકાઉન્ટ શોધી શકતું નથી. `ccam import` આ બંનેને લિંક કરીને આ સંભાળી લે છે. જો `~/.claude` પાસે, `CLAUDE_CONFIG_DIR=~/.claude` સાથેના અગાઉના run માંથી, પહેલેથી પોતાની `.claude.json` હોય, તો ccam બંને ફાઇલો રહેવા દે છે અને તેમ જણાવે છે, કારણ કે કઈ સાચી છે તે ફક્ત તમે જ જાણો છો.

**macOS પર લૉગિન ફોલ્ડરમાં નહીં, Keychain માં હોય છે.** Claude Code તેને `Claude Code-credentials` નામ હેઠળ રાખે છે. જ્યારે `CLAUDE_CONFIG_DIR` સેટ હોય, ત્યારે નામના અંતે `-` અને ફોલ્ડર પાથના SHA-256 ના પહેલા 8 hex અંક પણ જોડાય છે. Keychain આઇટમનું એકાઉન્ટ નામ `$USER` છે. તેથી અલગ ફોલ્ડરને અલગ Keychain એન્ટ્રી મળે છે, અને ખસેડેલું કે નામ બદલેલું ફોલ્ડર પોતાનું લૉગિન ગુમાવે છે; આ જ કારણે macOS પર `ccam rename` ફોલ્ડરનું નામ એ જ રાખે છે. જો Keychain લખવાની ના પાડે, તો Claude Code ફોલ્ડરમાંની `.credentials.json` વાપરે છે, જેમ [તેના દસ્તાવેજીકરણમાં](https://code.claude.com/docs/en/authentication#credential-management) વર્ણવેલું છે.

**દરેક સાઇન-ઇન ફોલ્ડરમાં નથી રહેતું.** API key વગરનું Claude Console સાઇન-ઇન ફોલ્ડરની બહાર સંગ્રહાય છે, તેથી અલગ ફોલ્ડર આવાં બે સાઇન-ઇનને અલગ રાખી શકતાં નથી. claude.ai લૉગિન અને API key દરેક ફોલ્ડરમાં અલગ રહે છે.

## Windows પર ફોલ્ડર ક્યાં હોય છે?

ડિફૉલ્ટ રીતે તે `%USERPROFILE%\.claude` છે. લૉગિન તેની અંદર `.credentials.json` માં હોય છે, અથવા જો Claude Code ને Windows Credential Manager વાપરવા સેટ કર્યું હોય તો તેમાં, macOS ની જેમ ફોલ્ડર પાથ પરથી બનેલા નામ હેઠળ. દરેક વધારાના એકાઉન્ટ માટે `CLAUDE_CONFIG_DIR` ને બીજા ફોલ્ડર પર સેટ કરો, Linux અને macOS ની જેમ જ.

## શું ફોલ્ડર ફાઇલો શેર કરી શકે?

Symbolic links કામ કરે છે: Claude Code લિંક કરેલી `settings.json` કે `projects/` ફોલ્ડર દ્વારા વાંચે અને લખે છે, અને લિંક જાળવી રાખે છે. MCP સર્વર અપવાદ છે, કારણ કે તેઓ લૉગિન સાથે `.claude.json` શેર કરે છે; તેના બદલે તેમને `--mcp-config` વડે આપો.

ccam બરાબર આ જ નિયમો પર બનેલું છે. તે ફોલ્ડર કેવી રીતે લિંક કરે છે તે માટે [શેરિંગ](/docs/sharing) જુઓ, અને મોટા ચિત્ર માટે [ઘણાં એકાઉન્ટ વિશેની માર્ગદર્શિકા](/guides/multiple-claude-code-accounts) જુઓ.
