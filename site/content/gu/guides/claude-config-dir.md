---
title: CLAUDE_CONFIG_DIR સમજૂતી
description: Claude Code માં CLAUDE_CONFIG_DIR શું કરે છે, ફોલ્ડરમાં શું સંગ્રહાય છે, macOS લૉગિન ક્યાં રાખે છે અને એકાઉન્ટ વચ્ચે કઈ ફાઇલો શેર થઈ શકે.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` Claude Code ને જણાવે છે કે `~/.claude` ને બદલે કયું ફોલ્ડર વાપરવું. એકાઉન્ટને લગતું બધું ત્યાં રહે છે: લૉગિન, સેટિંગ્સ, MCP સર્વર, plugins અને વાતચીતનો ઇતિહાસ. તેને જુદાં ફોલ્ડર તરફ નિર્દેશિત કરવાથી તમને અલગ એકાઉન્ટ મળે છે.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## ફોલ્ડરની અંદર શું હોય છે

Claude Code 2.1.302 સાથે તપાસેલું:

| પાથ | સામગ્રી | શેર કરી શકાય |
| --- | --- | --- |
| `.credentials.json` | લૉગિન ટોકન (Linux), MCP OAuth ટોકન | ના |
| `.claude.json` | એકાઉન્ટની વિગતો, user MCP સર્વર, દરેક પ્રોજેક્ટની સ્થિતિ | ના (MCP સર્વર અલગથી) |
| `settings.json` | પરવાનગીઓ, hooks, મૉડલ, સ્ટેટસ લાઇન, plugins | હા |
| `CLAUDE.md` | ગ્લોબલ સૂચનાઓ | હા |
| `agents/`, `commands/`, `skills/`, `plugins/` | એક્સટેન્શન્સ | હા |
| `projects/` | `.jsonl` ફાઇલો તરીકે વાતચીત, દરેક પ્રોજેક્ટની મેમરી | હા |
| `history.jsonl` | પ્રૉમ્પ્ટનો ઇતિહાસ | હા |
| `sessions/`, `shell-snapshots/`, `cache/` અને બીજાં | રનટાઇમ સ્થિતિ | ના |

## બે વિગતો જેમાં લોકો ભૂલ કરે છે

**`CLAUDE_CONFIG_DIR` વગર, `.claude.json` એ `~/.claude` ની અંદર હોતી નથી.** ત્યારે Claude Code તેને ફોલ્ડરની બાજુમાં `~/.claude.json` પર રાખે છે. જો તમે `CLAUDE_CONFIG_DIR=~/.claude` વાપરવાનું શરૂ કરો, તો Claude Code તેના બદલે `~/.claude/.claude.json` શોધે છે અને તમારું એકાઉન્ટ શોધી શકતું નથી. `ccam import` આ બંનેને લિંક કરીને આ સંભાળી લે છે.

**macOS પર લૉગિન ફોલ્ડરમાં નહીં, Keychain માં હોય છે.** Claude Code તેને `Claude Code-credentials` નામ હેઠળ રાખે છે, અને જ્યારે `CLAUDE_CONFIG_DIR` સેટ હોય ત્યારે તેની પાછળ ફોલ્ડર પાથનો ટૂંકો hash જોડાય છે. તેથી અલગ ફોલ્ડરને અલગ Keychain એન્ટ્રી મળે છે. જો Keychain લખવાની ના પાડે, તો Claude Code ફોલ્ડરમાંની `.credentials.json` વાપરે છે.

## ફોલ્ડર વચ્ચે શેરિંગ

Symbolic links કામ કરે છે: Claude Code લિંક કરેલી `settings.json` કે `projects/` ફોલ્ડર દ્વારા વાંચે અને લખે છે, અને લિંક જાળવી રાખે છે. MCP સર્વર અપવાદ છે, કારણ કે તેઓ લૉગિન સાથે `.claude.json` શેર કરે છે; તેના બદલે તેમને `--mcp-config` વડે આપો.

ccam બરાબર આ જ નિયમો પર બનેલું છે. તે ફોલ્ડર કેવી રીતે લિંક કરે છે તે માટે [શેરિંગ](/docs/sharing) જુઓ, અને મોટા ચિત્ર માટે [ઘણાં એકાઉન્ટ વિશેની માર્ગદર્શિકા](/guides/multiple-claude-code-accounts) જુઓ.
