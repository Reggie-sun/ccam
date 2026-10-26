---
title: साझा MCP सर्वर
description: ccam mcp add, add-json, list, remove और import से हर ग्रुप के लिए MCP सर्वरों की एक सूची रखें, जो --mcp-config से दी जाती है या हर अकाउंट में sync होती है।
updated: 2026-10-26
order: 5
translation: machine
sourceHash: 27a06ad35295
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
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
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
only in account side, for /home/me/code/app only: notes
```

ग्रुप्स के नीचे, `ccam mcp list` हर उस सर्वर का नाम बताता है जो दूसरे अकाउंट्स तक नहीं पहुँचता:

- किसी अकाउंट के अपने लिए जोड़े गए सर्वर के लिए `only in account work: pg`;
- किसी एक फ़ोल्डर के लिए जोड़े गए सर्वर के लिए `only in account side, for /home/me/code/app only: notes` (सादा `claude mcp add`, नीचे देखें);
- `account work has its own …, used in its sessions in place of the group's`, जब अकाउंट के अपने सर्वर का नाम किसी ग्रुप सर्वर वाला हो;
- इस फ़ोल्डर की `.mcp.json`, उन अकाउंट्स के साथ जिनकी मंज़ूरी का वह `waiting for approval in` है;
- वे सर्वर जो `in no account, only in ~/.claude.json (plain claude's)` हैं, उन्हें कॉपी करने वाली [`ccam migrate --mcp`](/features/migrate) कमांड के साथ।

ऐसा नाम जोड़ना जो ग्रुप के पास पहले से है, उसकी definition बदल देता है और `added` की जगह `updated` प्रिंट करता है। सर्वर किसी अकाउंट तक केवल तभी पहुँचते हैं जब वह उस ग्रुप में हो और ग्रुप MCP साझा करता हो (`ccam group set <group> --add mcp`); जब ग्रुप ऐसा नहीं करता, तो `add` यह बताता है।

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
| `--dry-run` (ग्लोबल) | केवल `mcp list` इसका समर्थन करता है; `add`, `add-json`, `remove` और `import` इसे अस्वीकार करते हैं और कुछ नहीं बदलते | बंद |

`add` को या तो `--` के बाद एक कमांड चाहिए या `--url`; वरना यह `a command to run (after --) or --url is required` के साथ विफल होता है। `--header`, `--url` सर्वरों पर लागू होता है और `--env` कमांड सर्वरों पर। `add-json` वही JSON लेता है जो `claude mcp add-json` लेता है, और definition एक ऑब्जेक्ट होनी चाहिए जिसमें `command` या `url` हो। `import` हमेशा वह ग्रुप इस्तेमाल करता है जिसके ज़रिए अकाउंट MCP साझा करता है, और `list` हर ग्रुप दिखाता है, इसलिए `-g` इन्हें नहीं बदलता।

## यह कैसे काम करता है

शामिल फ़ाइलें:

- `~/.local/share/ccam/groups/<group>/mcp.json`: साझा सर्वर, `{"mcpServers": {…}}`, वह format जिसे `--mcp-config` पढ़ता है। इसे `fsutil.AtomicWrite` से और मोड `0600` के साथ लिखा जाता है, क्योंकि सर्वर entries में अक्सर API keys होती हैं, और केवल lock `<group>/.mcp.json.lock` के तहत, ताकि एक ही समय पर चलाई गई दो `ccam mcp add` कमांड दोनों अपने सर्वर बनाए रखें।
- `~/.local/share/ccam/profiles/<name>/.claude.json`: अकाउंट के सर्वरों वाली Claude Code की अपनी फ़ाइल। ccam `profile.TopLevelValue` scanner से केवल इसका `mcpServers` मान पढ़ता है और शुरू करने के रास्ते पर कभी पूरी फ़ाइल parse नहीं करता।
- `~/.local/state/ccam/mcp/<name>/`: `seen` फ़ाइल, `sync-stamp` (ccam की लिखी हर कॉपी के hash के साथ, ताकि वह ऐसी कॉपी जिसे किसी ने नहीं बदला, आपकी बदली हुई कॉपी से अलग पहचान सके), `sync-checked` (sync strategy ने आख़िरी बार क्या देखा), `kept` (MCP sharing छोड़ते समय अकाउंट ने कौन-से सर्वर रखे) और `.claude.json` के बैकअप।

दो strategies हैं, जो `~/.config/ccam/config.toml` में `mcp_strategy` से सेट होती हैं। डिफ़ॉल्ट flag strategy है:

<!-- motion:mcp-flag -->

1. ccam अकाउंट का MCP ग्रुप खोजता है (`launcher.applyMCP`) और एक `stat` से जाँचता है कि ग्रुप की `mcp.json` मौजूद है, उसे parse किए बिना (`mcp.FlagArgs`)।
2. ccam `--mcp-config=<path>` को `claude` के ठीक बाद, आपके अपने आर्ग्युमेंट से पहले रखता है, ताकि आपका अपना `--mcp-config` या `--strict-mcp-config` उसके ऊपर फिर भी लागू हो (`Plan.InsertArgs`)।
3. ccam जाँचता है कि पिछली बार शुरू होने के बाद अकाउंट की `.claude.json` बदली है या नहीं, और केवल तभी उसके सर्वर parse करता है (`mcp.NewLocal`)।
4. ccam ख़ुद की जगह Claude Code को रख देता है (`syscall.Exec`)।

पाथ एक आर्ग्युमेंट `--mcp-config=<path>` के रूप में दिया जाता है, दो के रूप में नहीं। Claude Code का `--mcp-config` कई मान लेता है, इसलिए `claude --mcp-config a.json "prompt"` prompt को दूसरी config फ़ाइल समझकर `MCP config file not found` के साथ विफल हो जाता। एक आर्ग्युमेंट वाला रूप इससे बचाता है, इसलिए उसके बाद आने वाला prompt या subcommand काम करता है। `ccam work mcp list`, `get` और `remove` को कोई flag नहीं मिलता, क्योंकि Claude Code वहाँ इसे अस्वीकार करता है; इसके बजाय ccam standard error पर ग्रुप के सर्वरों के नाम बताता है ([विवरण](/features/launch#shared-mcp-servers))। Claude Code `--mcp-config` के सर्वरों को अकाउंट के अपने सर्वरों के साथ लोड करता है; नाम टकराने पर `.claude.json` की entry जीतती है। Flag strategy `.claude.json` में केवल Remote Control के लिए, जब कोई अकाउंट MCP sharing छोड़ता है या उसमें फिर से जुड़ता है ([नीचे](#leaving-and-joining)), और `ccam mcp import` के लिए लिखती है।

`claude remote-control` (`rc`) कोई `--mcp-config` नहीं लेता, और उसके sessions सर्वर केवल `.claude.json` से पढ़ते हैं। `ccam work rc` के लिए, ccam flag छोड़ देता है और टर्मिनल पर पूछता है कि क्या ग्रुप के सर्वर बैकअप के साथ `.claude.json` में कॉपी किए जाएँ, जैसे `sync` strategy करती है। यह तब तक पूछता है जब तक कोई कॉपी न बन जाए, फिर केवल ग्रुप की `mcp.json` बदलने या किसी कॉपी किए गए सर्वर के ग़ायब होने के बाद ([Remote Control](/features/launch#remote-control))। चूँकि नाम टकराने पर ऐसी कॉपी `--mcp-config` पर जीतती है, `ccam mcp add`, `add-json`, `remove` और `import` ग्रुप के दूसरे अकाउंट्स के पास रखी कॉपियाँ अपडेट करते हैं:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

`✓ Updated` पंक्ति केवल तब दिखती है जब कोई कॉपी बदली हो। ccam के किसी पुराने वर्ज़न की बनाई कॉपी भी इसी तरह अपडेट होती है, और `ccam migrate` के बाद की कॉपी भी।

किसी session के अंदर `claude mcp add` से जोड़ा गया सर्वर केवल उस अकाउंट की `.claude.json` में जाता है। सादा `claude mcp add` local scope इस्तेमाल करता है: सर्वर उसी एक अकाउंट के लिए उसी एक फ़ोल्डर में लोड होता है (`-s user` उसे हर फ़ोल्डर में अकाउंट के लिए उपलब्ध बना देता है)। ccam इसे होते समय नहीं देख सकता, क्योंकि Claude Code तब तक ccam process की जगह ले चुका होता है, इसलिए यह अगली बार शुरू होने पर जाँचता है। `seen` फ़ाइल में `.claude.json` का modification time और पहले से बताए गए सर्वरों के नाम होते हैं। जब समय अलग होता है, तो ccam नामों की तुलना ग्रुप से करता है और हर नए सर्वर के लिए एक बार, input का इंतज़ार किए बिना, दोनों strategies के साथ, एक संकेत प्रिंट करता है:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

संकेत अकाउंट के ऐसे अपने सर्वर का नाम भी बताता है जिसका नाम किसी ग्रुप सर्वर वाला है, क्योंकि उसके sessions में अकाउंट वाला जीतता है।

`mcp_strategy = "sync"` के साथ, सर्वर `claude mcp list` और `/mcp` में सामान्य user सर्वरों की तरह दिखते हैं, और ccam शुरू होने से पहले उन्हें हर अकाउंट की `.claude.json` में लिखता है:

<!-- motion:mcp-sync -->

1. ccam ग्रुप की `mcp.json` के modification time की तुलना `sync-stamp` से करता है, और अकाउंट की `.claude.json` की उससे जो उसने आख़िरी बार देखा (`mcp.Sync`)। अगर दोनों में से कोई नहीं बदला, तो और कुछ नहीं होता। वरना यह यह भी जाँचता है कि हर ग्रुप सर्वर अभी भी `.claude.json` में है, और ग़ायब हुई कॉपी वापस रखता है (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`)।
2. यह Claude Code का अपना lock लेता है, यानी फ़ाइल के बगल में `.claude.json.lock` फ़ोल्डर, जो `mkdir` से बनता है। यह 2 सेकंड तक इंतज़ार करता है और 10 सेकंड से पुराने lock को stale मानता है, जैसे Claude Code की lock library करती है।
3. यह मौजूदा `.claude.json` को state फ़ोल्डर के `backups/` में सेव करता है और सबसे नई 5 रखता है।
4. यह `mcpServers` में हर ग्रुप सर्वर सेट करता है, पहले लिखे गए वे सर्वर हटाता है जो अब ग्रुप में नहीं हैं, और फ़ाइल को उसके मूल मोड के साथ atomic रूप से लिखता है, केवल तब जब कुछ बदला हो। यह केवल वही कॉपियाँ बदलता या हटाता है जो उसने लिखीं और जिन्हें तब से किसी ने नहीं बदला। अगर अकाउंट का अपना कोई सर्वर किसी ग्रुप सर्वर के नाम वाला है, तो ccam अकाउंट का सर्वर रखता है और यह एक बार बताता है: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`।
5. यह हर कॉपी के hash के साथ नया stamp दर्ज करता है, और Claude Code `--mcp-config` के बिना शुरू होता है।

जब अकाउंट का कोई session चल रहा हो, चाहे Claude Code किसी भी तरह शुरू हुआ हो, तो sync `ccam: account work is already running, skipping MCP sync` के साथ छोड़ दिया जाता है, क्योंकि Claude Code चलते समय `.claude.json` फिर से लिखता है। विफल sync `ccam: MCP sync failed (launching anyway): …` प्रिंट करता है और फिर भी शुरू करता है।

### छोड़ना और जुड़ना

जो अकाउंट MCP साझा करना बंद करता है (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`), वह ग्रुप के सर्वर रखता है: योजना उन्हें बैकअप के साथ उसकी `.claude.json` में कॉपी करती है, और उन्हीं नामों वाले उसके अपने सर्वर बने रहते हैं। `--fresh` के साथ उसे कोई नहीं मिलता, और योजना उन सर्वरों की सूची देती है जो लोड होना बंद कर देते हैं। जब वह फिर से ऐसे ग्रुप में जुड़ता है जिसमें वे सर्वर हैं, तो रखी गई वे कॉपियाँ जिन्हें किसी ने नहीं बदला, हटा दी जाती हैं, ताकि ग्रुप के बाद के बदलाव उस तक पहुँचें।

`ccam mcp import` अकाउंट के अपने सर्वरों को उसके ग्रुप में move करता है:

<!-- motion:mcp-import -->

ccam अकाउंट के सर्वर पढ़ता है, उसके अपने और वे जो उसके पास केवल एक फ़ोल्डर के लिए हैं, और `mcp.SameDef` से हर एक की तुलना उसी नाम के ग्रुप सर्वर से करता है, जो keys का क्रम और spacing अनदेखा करता है। फ़ोल्डर सर्वर हर फ़ोल्डर के लिए ग्रुप सर्वर बन जाता है (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); दो अलग फ़ोल्डर definitions वाला एक नाम अस्वीकार होता है। Plugin सर्वर पेश नहीं किए जाते, क्योंकि plugin उन्हें लाता है। जो सर्वर ग्रुप में नहीं है, या उसी definition के साथ है, वह कॉपी होता है। अलग सर्वर छोड़ दिया जाता है, क्योंकि ग्रुप की definition ही वह है जिसे हर दूसरा अकाउंट इस्तेमाल करता है:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

`--force` के साथ ग्रुप का सर्वर बदल दिया जाता है, और पिछली `mcp.json` पहले ग्रुप फ़ोल्डर के अंदर `.ccam-backup/<time>/mcp.json` में सेव होती है। फिर ccam, उसी lock के तहत और state फ़ोल्डर में बैकअप के बाद, कॉपी की गई entries को अकाउंट की `.claude.json` से हटाता है, उसकी अपनी और फ़ोल्डर वाली एक ही write में (`mcp.Remove`)। ccam ग्रुप की `mcp.json` को फ़ाइल पढ़ने से लेकर उन entries के हटने तक locked रखता है, ताकि बीच में कोई सर्वर को ग्रुप से न हटाए। यह ज़रूरी है क्योंकि नाम टकराने पर अकाउंट की entry ग्रुप की entry से जीत जाती। `--keep-local` के साथ entries अकाउंट की `.claude.json` में रहती हैं।

## सुरक्षा और विशेष स्थितियाँ

- **साइन-इन।** OAuth से साइन इन करने वाले सर्वर अपने tokens हर अकाउंट की `.credentials.json` (`mcpOAuth`) में रखते हैं, सर्वर सूची में नहीं, इसलिए सर्वर साझा करने से उसका साइन-इन साझा नहीं होता। हर अकाउंट में `/mcp` से एक बार साइन इन करें। ccam कभी साइन-इन कॉपी नहीं करता, और [`ccam migrate`](/features/migrate) भी नहीं करता: कॉपी किया गया refresh token हर दूसरे धारक को लॉग आउट कर देता है। [`ccam ext`](/features/ext) और `ccam doctor` उन सर्वरों के नाम बताते हैं जिनमें कोई अकाउंट साइन इन नहीं है जबकि कोई दूसरा अकाउंट है, और `ccam next --continue` बताता है कि नए अकाउंट में कौन-से नहीं हैं। ग्रुप को कभी कोई लॉगिन फ़ाइल नहीं मिलती।
- **चल रहे अकाउंट।** `import` ग्रुप में कॉपी करता है लेकिन चल रहे अकाउंट की अपनी entries रहने देता है और यह बताता है; session ख़त्म होने के बाद इसे फिर से चलाएँ। Sync चल रहे अकाउंट्स को छोड़ देता है। दोनों ccam के बिना शुरू हुए Claude Code sessions को भी पहचानते हैं।
- **Remote Control के लिए कॉपियाँ।** चल रहे अकाउंट की कॉपी को छुआ नहीं जाता: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` `mcp_strategy = "sync"` के साथ कुछ भी अपडेट करने की ज़रूरत नहीं होती, क्योंकि हर बार शुरू करने पर sync होता है।
- **टूटी फ़ाइलें।** अगर `.claude.json` वैध JSON नहीं है, तो sync उसे बिना छुए छोड़ देता है और त्रुटि बताता है। अगर 2 सेकंड में lock न लिया जा सके, तो ccam `could not lock … (Claude Code is writing the file)` बताता है।
- **MCP साझा न करने वाले अकाउंट।** `import` मना कर देता है और कहता है कि पहले अकाउंट को ऐसे ग्रुप में रखें जो MCP साझा करता हो (`ccam group join <group> side`)।
- **ग्रुप के नाम।** `-g` केवल ग्रुप का नाम लेता है; वह `groups/` के नीचे एक फ़ोल्डर बनता है, इसलिए `../x` जैसे नाम अस्वीकार किए जाते हैं।

## संबंधित

- [ग्रुप के साथ साझा करना](/features/share)
- [ccam के साथ Remote Control का उपयोग](/guides/remote-control)
- [MCP सर्वर, सेटिंग्स और इतिहास साझा करें](/docs/sharing)
- [मौजूदा सेटअप कॉपी करना](/features/migrate)
- [कमांड संदर्भ](/docs/commands)
