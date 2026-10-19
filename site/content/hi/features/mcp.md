---
title: साझा MCP सर्वर
description: ccam mcp add, add-json, list, remove और import से हर ग्रुप के लिए MCP सर्वरों की एक सूची रखें, जो --mcp-config से दी जाती है या हर अकाउंट में sync होती है।
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
---

Claude Code user MCP सर्वर `.claude.json` में रखता है, वही फ़ाइल जिसमें अकाउंट का लॉगिन होता है, इसलिए यह फ़ाइल अकाउंट्स के बीच साझा नहीं हो सकती। ccam साझा सर्वरों को एक ग्रुप फ़ाइल `mcp.json` में रखता है, और शुरू होते समय उसे ग्रुप के हर अकाउंट को देता है। किसी सर्वर को एक बार में अपने सभी अकाउंट्स के लिए जोड़ने के लिए `ccam mcp` इस्तेमाल करें, और किसी एक session में जोड़े गए सर्वर को दूसरों तक ले जाने के लिए `ccam mcp import`।

## उपयोग

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

बदलाव हर अकाउंट के अगली बार शुरू होने से लागू होता है; चल रहे session उन्हीं सर्वरों के साथ चलते रहते हैं जिनके साथ वे शुरू हुए थे:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

ऐसा नाम जोड़ना जो ग्रुप के पास पहले से है, उसकी definition बदल देता है और `added` की जगह `updated` प्रिंट करता है। सर्वर किसी अकाउंट तक केवल तभी पहुँचते हैं जब वह उस ग्रुप के साथ MCP साझा करता हो (`ccam share <name> --mcp on`); जब अभी कोई अकाउंट MCP साझा नहीं करता, तो `add` यह बताता है।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `-g`, `--group <group>` | `add`, `add-json` और `remove` के लिए लक्ष्य ग्रुप | वह इकलौता ग्रुप जिसके ज़रिए कोई अकाउंट MCP साझा करता है, वरना `default` |
| `-- <command> [args…]` (`add`) | इस कमांड को stdio सर्वर के रूप में चलाएँ | |
| `--url <url>` (`add`) | रिमोट सर्वर का URL (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | `--url` सर्वर का transport | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, कई बार दिया जा सकता है | कोई नहीं |
| `-e`, `--env <KEY=VALUE>` (`add`) | stdio सर्वर के लिए environment variable, कई बार दिया जा सकता है | कोई नहीं |
| `--force` (`import`) | उसी नाम का अलग ग्रुप सर्वर बदल दें; पुरानी `mcp.json` का बैकअप लिया जाता है | बंद |
| `--keep-local` (`import`) | अकाउंट की `.claude.json` में entries रहने दें | बंद |
| `-h`, `--help` | हर कमांड की help | |
| `--dry-run` (ग्लोबल) | केवल `mcp list` इसे मानता है; `add`, `add-json`, `remove` और `import` इसे अस्वीकार करते हैं और कुछ नहीं बदलते | बंद |

`add` को या तो `--` के बाद एक कमांड चाहिए या `--url`; वरना यह `a command to run (after --) or --url is required` के साथ विफल होता है। `--header`, `--url` सर्वरों पर लागू होता है और `--env` कमांड सर्वरों पर। `add-json` वही JSON लेता है जो `claude mcp add-json` लेता है, और definition एक ऑब्जेक्ट होनी चाहिए जिसमें `command` या `url` हो। `import` हमेशा वह ग्रुप इस्तेमाल करता है जिसके ज़रिए अकाउंट MCP साझा करता है, और `list` हर ग्रुप दिखाता है, इसलिए `-g` इन्हें नहीं बदलता।

## यह कैसे काम करता है

शामिल फ़ाइलें:

- `~/.local/share/ccam/groups/<group>/mcp.json`: साझा सर्वर, `{"mcpServers": {…}}`, वह format जिसे `--mcp-config` पढ़ता है। इसे `fsutil.AtomicWrite` से और मोड `0600` के साथ लिखा जाता है, क्योंकि सर्वर entries में अक्सर API keys होती हैं।
- `~/.local/share/ccam/profiles/<name>/.claude.json`: अकाउंट के सर्वरों वाली Claude Code की अपनी फ़ाइल। ccam `profile.TopLevelValue` scanner से केवल इसका `mcpServers` मान पढ़ता है और शुरू करने के रास्ते पर कभी पूरी फ़ाइल parse नहीं करता।
- `~/.local/state/ccam/mcp/<name>/`: `seen` फ़ाइल, `sync-stamp` और `.claude.json` के बैकअप।

दो strategies हैं, जो `~/.config/ccam/config.toml` में `mcp_strategy` से सेट होती हैं। डिफ़ॉल्ट flag strategy है:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` अकाउंट का MCP ग्रुप खोजता है। `mcp.FlagArgs` एक `stat` से जाँचता है कि ग्रुप की `mcp.json` मौजूद है, उसे parse किए बिना।
2. `Plan.InsertArgs`, `--mcp-config=<path>` को `claude` के ठीक बाद, आपके अपने आर्ग्युमेंट से पहले रखता है, ताकि आपका अपना `--mcp-config` या `--strict-mcp-config` उसके ऊपर फिर भी लागू हो।
3. `mcp.NewLocal` जाँचता है कि पिछली बार शुरू होने के बाद अकाउंट की `.claude.json` बदली है या नहीं, और केवल तभी उसके सर्वर parse करता है।
4. `syscall.Exec`, ccam की जगह Claude Code को रख देता है।

पाथ एक आर्ग्युमेंट `--mcp-config=<path>` के रूप में दिया जाता है, दो के रूप में नहीं। Claude Code का `--mcp-config` कई मान लेता है, इसलिए `claude --mcp-config a.json "prompt"` prompt को दूसरी config फ़ाइल समझकर `MCP config file not found` के साथ विफल हो जाता। एक आर्ग्युमेंट वाला रूप वहीं रुक जाता है, और उसके बाद आने वाला prompt या `mcp list` जैसी subcommand काम करती है। Claude Code `--mcp-config` के सर्वरों को अकाउंट के अपने सर्वरों के साथ लोड करता है; नाम टकराने पर `.claude.json` की entry जीतती है। Flag strategy केवल Remote Control के लिए `.claude.json` में लिखती है।

`claude remote-control` (`rc`) कोई `--mcp-config` नहीं लेता, और उसके sessions सर्वर केवल `.claude.json` से पढ़ते हैं। `ccam work rc` के लिए, ccam flag छोड़ देता है और टर्मिनल पर पूछता है कि क्या ग्रुप के सर्वर बैकअप के साथ `.claude.json` में कॉपी किए जाएँ, जैसे `sync` strategy करती है। यह तब तक पूछता है जब तक कोई कॉपी न बन जाए, फिर केवल ग्रुप की `mcp.json` बदलने के बाद ([Remote Control](/features/launch#remote-control))। चूँकि नाम टकराने पर ऐसी कॉपी `--mcp-config` पर जीतती है, `ccam mcp add`, `add-json`, `remove` और `import` ग्रुप के दूसरे अकाउंट्स के पास रखी कॉपियाँ अपडेट करते हैं:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

किसी session के अंदर `claude mcp add` से जोड़ा गया सर्वर केवल उस अकाउंट की `.claude.json` में जाता है। ccam उस समय उसे नहीं देख सकता, क्योंकि Claude Code ने उसकी जगह ले ली है, इसलिए यह अगली बार शुरू होने पर जाँचता है। `seen` फ़ाइल में `.claude.json` का modification time और पहले से बताए गए सर्वरों के नाम होते हैं। जब समय अलग होता है, तो ccam नामों की तुलना ग्रुप से करता है और हर नए सर्वर के लिए एक बार, input का इंतज़ार किए बिना, एक संकेत प्रिंट करता है:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

`mcp_strategy = "sync"` के साथ, सर्वर `claude mcp list` और `/mcp` में सामान्य user सर्वरों की तरह दिखते हैं, और ccam शुरू होने से पहले उन्हें हर अकाउंट की `.claude.json` में लिखता है:

<!-- motion:mcp-sync -->

1. `mcp.Sync` ग्रुप की `mcp.json` के modification time की तुलना `sync-stamp` से करता है। अगर वह नहीं बदला है, तो अकाउंट की `.claude.json` खोली नहीं जाती।
2. यह Claude Code का अपना lock लेता है, यानी फ़ाइल के बगल में `.claude.json.lock` फ़ोल्डर, जो `mkdir` से बनता है। यह 2 सेकंड तक इंतज़ार करता है और 10 सेकंड से पुराने lock को stale मानता है, जैसे Claude Code की lock library करती है।
3. यह मौजूदा `.claude.json` को state फ़ोल्डर के `backups/` में सेव करता है और सबसे नई 5 रखता है।
4. यह `mcpServers` में हर ग्रुप सर्वर सेट करता है, पहले लिखे गए वे सर्वर हटाता है जो अब ग्रुप में नहीं हैं, अकाउंट के अपने सर्वर रहने देता है, और फ़ाइल को उसके मूल मोड के साथ atomic रूप से लिखता है।
5. यह नया stamp और अपने प्रबंधित नाम दर्ज करता है, और Claude Code `--mcp-config` के बिना शुरू होता है।

जब अकाउंट का कोई session चल रहा हो, तो sync `ccam: account work is already running, skipping MCP sync` के साथ छोड़ दिया जाता है, क्योंकि Claude Code चलते समय `.claude.json` फिर से लिखता है। विफल sync `ccam: MCP sync failed (launching anyway): …` प्रिंट करता है और फिर भी शुरू करता है। Session के अंदर जोड़े गए सर्वरों का संकेत केवल flag strategy का हिस्सा है।

`ccam mcp import` अकाउंट के अपने सर्वरों को उसके ग्रुप में move करता है:

<!-- motion:mcp-import -->

ccam `mcp.LocalServers` से अकाउंट के सर्वर पढ़ता है और `mcp.SameDef` से हर एक की तुलना उसी नाम के ग्रुप सर्वर से करता है, जो keys का क्रम और spacing अनदेखा करता है। जो सर्वर ग्रुप में नहीं है, या उसी definition के साथ है, वह कॉपी होता है। अलग सर्वर छोड़ दिया जाता है, क्योंकि ग्रुप की definition ही वह है जिसे हर दूसरा अकाउंट इस्तेमाल करता है:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` के साथ ग्रुप का सर्वर बदल दिया जाता है, और पिछली `mcp.json` पहले ग्रुप फ़ोल्डर के अंदर `.ccam-backup/<time>/mcp.json` में सेव होती है। फिर `mcp.RemoveLocal`, उसी lock के तहत और state फ़ोल्डर में बैकअप के बाद, कॉपी की गई entries को अकाउंट की `.claude.json` से हटाता है। यह ज़रूरी है क्योंकि नाम टकराने पर अकाउंट की entry ग्रुप की entry से जीत जाती। `--keep-local` उन्हें रहने देता है।

## सुरक्षा और विशेष स्थितियाँ

- **साइन-इन।** OAuth से साइन इन करने वाले सर्वर अपने tokens हर अकाउंट की `.credentials.json` (`mcpOAuth`) में रखते हैं, सर्वर सूची में नहीं, इसलिए सर्वर साझा करने से उसका साइन-इन साझा नहीं होता। हर अकाउंट में `/mcp` से एक बार साइन इन करें, या [`ccam migrate`](/features/migrate) से `~/.claude` के साइन-इन कॉपी करें। ग्रुप को कभी कोई लॉगिन फ़ाइल नहीं मिलती।
- **चल रहे अकाउंट।** `import` ग्रुप में कॉपी करता है लेकिन चल रहे अकाउंट की अपनी entries रहने देता है और यह बताता है; session ख़त्म होने के बाद इसे फिर से चलाएँ। Sync चल रहे अकाउंट्स को छोड़ देता है।
- **Remote Control के लिए कॉपियाँ।** चल रहे अकाउंट की कॉपी को छुआ नहीं जाता: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` के साथ कुछ भी अपडेट करने की ज़रूरत नहीं होती, क्योंकि हर बार शुरू करने पर sync होता है।
- **टूटी फ़ाइलें।** अगर `.claude.json` वैध JSON नहीं है, तो sync उसे बिना छुए छोड़ देता है और त्रुटि बताता है। अगर 2 सेकंड में lock न लिया जा सके, तो ccam `could not lock … (Claude Code is writing the file)` बताता है।
- **MCP साझा न करने वाले अकाउंट।** `import` मना कर देता है और पहले `ccam share side --mcp on` चलाने को कहता है।
- **ग्रुप के नाम।** `-g` केवल ग्रुप का नाम लेता है; वह `groups/` के नीचे एक फ़ोल्डर बनता है, इसलिए `../x` जैसे नाम अस्वीकार किए जाते हैं।

## संबंधित

- [अकाउंट्स के बीच साझा करना](/features/share)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [मौजूदा सेटअप कॉपी करना](/features/migrate)
- [सभी कमांड](/docs/commands)
