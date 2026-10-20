---
title: ccam और दूसरे Claude Code account switchers
description: Claude Code अकाउंट स्विच करने के लिए ccam की तुलना claude-swap, CCS, clausona, clauth, CAAM, cux और दूसरे टूल्स से, और कौन-सा तरीका किस काम के लिए ठीक है।
updated: 2026-10-20
order: 2
translation: machine
sourceHash: 68048287f792
---

Claude Code अकाउंट स्विच करने वाले टूल दो में से एक तरीका अपनाते हैं। ज़्यादातर टूल एक ही `~/.claude` फ़ोल्डर के अंदर लॉगिन बदलते हैं, इसलिए एक समय में एक ही अकाउंट सक्रिय रहता है। दूसरे टूल, जिनमें ccam भी है, `CLAUDE_CONFIG_DIR` से हर अकाउंट को उसका अपना फ़ोल्डर देते हैं, इसलिए हर अकाउंट लॉग इन रहता है और कई अकाउंट एक साथ चल सकते हैं। ccam इसमें साझा सेटिंग्स और MCP सर्वर जोड़ता है, और किसी बातचीत को दूसरे अकाउंट के फ़ोल्डर में जारी रखता है।

## एक लॉगिन बदलना या अलग-अलग फ़ोल्डर?

**लॉगिन बदलने** में सेटिंग्स, MCP सर्वर और इतिहास का एक ही सेट रहता है, क्योंकि फ़ोल्डर एक ही है। स्विच करना तेज़ है, और चल रहा session नया लॉगिन ले सकता है। लेकिन एक समय में एक ही अकाउंट सक्रिय रहता है, सभी terminals और editors के लिए एक साथ।

लॉगिन बदलने के लिए टूल बाक़ी हर अकाउंट के लॉगिन की एक कॉपी रखता है और ज़रूरत पड़ने पर उसे सही जगह लिख देता है। Claude Code लॉगिन को renew करते समय उसका refresh token बदल देता है, इसलिए हर सहेजी गई कॉपी को चल रहे sessions के साथ मेल में रखना पड़ता है। स्विच के बाद लॉग आउट हो जाने वाले या दोबारा `/login` माँगने वाले अकाउंट उन समस्याओं में से हैं जिन्हें ये प्रोजेक्ट सबसे ज़्यादा ठीक करते हैं।

**अलग-अलग फ़ोल्डर** वह तरीका है जो [Claude Code के दस्तावेज़ में बताया गया है](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)। हर अकाउंट लॉग इन रहता है, और आप work और personal session साथ-साथ चला सकते हैं। हर लॉगिन अपने ही फ़ोल्डर में रहता है, जहाँ Claude Code उसे हमेशा की तरह renew करता है। इसकी क़ीमत यह है कि हर फ़ोल्डर अपनी अलग सेटिंग्स, MCP सर्वर और इतिहास के साथ शुरू होता है, जब तक कोई टूल उन्हें साझा न करे।

## टूल्स की तुलना

2026-10-20 को हर प्रोजेक्ट के अपने दस्तावेज़ और source code से लिया गया (clausona: v0.12.0-beta)। "दस्तावेज़ में नहीं" का मतलब है कि प्रोजेक्ट इसके बारे में कुछ नहीं कहता, यह नहीं कि सुविधा मौजूद नहीं है।

| टूल | तरीका | एक साथ कई अकाउंट | 5 घंटे और साप्ताहिक उपयोग | सीमा पूरी होने पर | प्लेटफ़ॉर्म |
| --- | --- | --- | --- | --- | --- |
| ccam | अलग-अलग फ़ोल्डर | हाँ | हाँ, `ccam usage` | आप `ccam next --auto --continue` चलाते हैं, और बातचीत सबसे ज़्यादा गुंजाइश वाले अकाउंट में जारी रहती है; `ccam --fallback` `-p` jobs को अपने-आप वहाँ फिर से चलाता है, और `ccam pick` workers को अकाउंट्स में बाँटता है | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | लॉगिन बदलता है; `cswap run` से अलग-अलग फ़ोल्डर | हाँ, `cswap run` से | हाँ | `cswap auto` एक सीमा-बिंदु पर लॉगिन बदलता है; चल रहा session जारी रहता है | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | अलग-अलग फ़ोल्डर | हाँ | आंशिक रूप से | Account pools एक local proxy के ज़रिए requests अगले अकाउंट को भेजते हैं | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | अलग-अलग फ़ोल्डर | हाँ | हाँ | Routes हर run से पहले quota वाला अकाउंट चुनते हैं, reset का इंतज़ार कर सकते हैं, और workers को अलग-अलग अकाउंट देते हैं | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | लॉगिन बदलता है; `clauth start` से अलग-अलग फ़ोल्डर | हाँ, `clauth start` से | हाँ | fallback chain में अगले अकाउंट पर स्विच करता है | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | लॉगिन बदलता है, या isolated profiles | हाँ, isolated profiles से | हाँ | `caam run` non-interactive कमांड को दूसरे अकाउंट पर फिर से चलाता है | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | लॉगिन बदलता है और `claude` को wrap करता है | नहीं, एक ही सक्रिय अकाउंट | हाँ | लॉगिन बदलता है और वही बातचीत फिर से शुरू करता है | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | लॉगिन बदलता है | नहीं | हाँ | एक सीमा-बिंदु पर वैकल्पिक स्विच | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | लॉगिन बदलता है | नहीं | हाँ | 100% पर वैकल्पिक स्विच | macOS |

ccswap, claude-swap का एक fork है जो Codex लॉगिन भी संभालता है। [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) और [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) भी हर अकाउंट को उसका अपना फ़ोल्डर देते हैं। इनमें से कई टूल Codex या Gemini जैसे दूसरे coding agents भी संभालते हैं।

## ये कैसे बने हैं

| टूल | किस भाषा में लिखा गया | Claude Code चलते समय | लॉगिन को एक जगह से दूसरी जगह कॉपी करता है | अकाउंट क्या साझा कर सकते हैं |
| --- | --- | --- | --- | --- |
| ccam | Go, एक static binary | कुछ नहीं: ccam ख़ुद को `claude` से बदल देता है (Windows पर, जहाँ यह संभव नहीं, यह उसके ख़त्म होने का इंतज़ार करता है) | नहीं | हर आइटम अलग से: हर अकाउंट के लिए local या कोई नामित group |
| claude-swap | Python 3.12+, uv या pipx से इंस्टॉल | `cswap run` के साथ कुछ नहीं; `cswap auto` और menu bar ऐप चलते रहते हैं | हाँ, जब यह लॉगिन बदलता है | `cswap run` के साथ: सेटिंग्स, `CLAUDE.md` और extensions एक साथ, इतिहास वैकल्पिक; MCP सर्वर हर launch पर कॉपी होते हैं |
| CCS | Node.js 18+ या Bun | CCS parent process के रूप में बना रहता है; pools एक local proxy जोड़ते हैं | अपने अकाउंट फ़ोल्डरों के लिए नहीं; pools अपने अलग लॉगिन रखते हैं | हर अकाउंट के लिए साझा या local, इतिहास context groups में; MCP सर्वर कॉपी होते हैं |
| clausona | TypeScript, Node.js 20+ | साधारण `claude` एक shell function से चलता है; `clausona run` Node को parent process के रूप में रखता है; हर exit के बाद एक background Node process | नहीं, लेकिन idle अकाउंट्स का उपयोग पढ़ने के लिए यह ख़ुद उनके tokens renew करता है | हर अकाउंट की अपनी फ़ाइलों को छोड़कर सब कुछ, एक primary से; इतिहास वैकल्पिक; `claude mcp add` से जोड़े गए सर्वर साझा नहीं होते |
| clauth | Rust | clauth parent process के रूप में बना रहता है | हाँ, जब यह लॉगिन बदलता है | पूरा `~/.claude`, या कुछ नहीं |
| CAAM | Go | mode पर निर्भर; एक वैकल्पिक daemon | हाँ, अपने vault mode में | mode से तय; MCP सर्वर हर launch पर कॉपी होते हैं |
| cux | Go | `claude` के चारों ओर एक wrapper बना रहता है | हाँ | सब कुछ: एक ही फ़ोल्डर |
| CCSwitcher, Claude Switcher | macOS ऐप (Swift, Python) | menu bar ऐप | हाँ | सब कुछ: एक ही फ़ोल्डर |

## ccam क्या अच्छा करता है

### लॉगिन वहीं रहते हैं जहाँ Claude Code ने उन्हें रखा

ccam कभी किसी Claude लॉगिन को एक जगह से दूसरी जगह कॉपी नहीं करता। हर अकाउंट का लॉगिन उसके अपने फ़ोल्डर में रहता है, या macOS पर उस फ़ोल्डर के अपने कीचेन आइटम में, और उसे renew करने वाला प्रोग्राम Claude Code ही है। ऐसी कोई सहेजी गई कॉपी नहीं है जो मेल से बाहर हो सके। मौजूदा `~/.claude`, जिसे `ccam import` से रजिस्टर किया गया है, वही कीचेन आइटम रखता है जिसे साधारण `claude` इस्तेमाल करता है, इसलिए दोनों तरह से वह लॉग इन रहता है।

ccam किसी expired token को ख़ुद तभी renew करता है जब आप `ccam usage --refresh` चलाते हैं। यह पहले एक lock लेता है और वह अकाउंट चल रहा हो तो मना कर देता है, इसलिए यह किसी चालू session के नीचे से token कभी नहीं बदलता।

### आपके और Claude Code के बीच कुछ नहीं चलता

ccam एक फ़ाइल पढ़ता है, अपने links जाँचता है और ख़ुद को `claude` से बदल देता है (`exec`)। काम करते समय कोई ccam process, wrapper, proxy या daemon नहीं होता: standard input, signals और exit code Claude Code के होते हैं, और उसकी requests सीधे Anthropic तक जाती हैं। Windows किसी process को बदल नहीं सकता, इसलिए वहाँ ccam Claude Code को उसी console के साथ शुरू करता है और उसके exit code का इंतज़ार करता है, इसके अलावा कुछ नहीं करता।

इस तरह किसी अकाउंट को शुरू करने में लगभग 4 ms जुड़ते हैं, 95th percentile पर 6 ms से कम, और peak memory 8 MB रहती है ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) hyperfine के साथ, 300 runs)। उसी laptop पर, कोई प्रोग्राम कुछ करे उससे पहले ही, सिर्फ़ Python interpreter शुरू होने में 13 ms और Node.js में 31 ms लगते हैं।

### हर आइटम और हर अकाउंट के हिसाब से साझा करें

MCP सर्वर, `settings.json`, `CLAUDE.md`, extensions (agents, commands, skills, plugins) और बातचीत का इतिहास पाँच अलग आइटम हैं। हर अकाउंट हर आइटम को `local` या किसी group पर सेट करता है, इसलिए दो work अकाउंट `work` group में इतिहास साझा कर सकते हैं जबकि personal अकाउंट अपना इतिहास अलग रखता है, और तीनों एक ही MCP सर्वर साझा करते हैं।

- MCP सर्वर launch के समय `--mcp-config` से दिए जाते हैं, इसलिए default रूप से ccam कभी `.claude.json` को edit नहीं करता, जिसे Claude Code चलते समय दोबारा लिखता है। `ccam mcp add` पूरे group के लिए एक ही बार में सर्वर जोड़ता है।
- बाक़ी आइटम group फ़ोल्डर की ओर links हैं, इसलिए एक अकाउंट में किया गया बदलाव तुरंत बाक़ी अकाउंट में दिखता है। launch के समय कुछ भी कॉपी नहीं होता।
- लॉगिन, `.claude.json` और running state कभी साझा नहीं होते।

### हर बदलाव पहले दिखाया जाता है और उसका backup बनता है

`ccam share` और `ccam migrate` कुछ भी खिसकाने से पहले एक plan दिखाते हैं और पूछते हैं। जब किसी अकाउंट और उसके group में कोई फ़ाइल अलग हो, तो आप चुनते हैं: group वाली रखें, अपनी इस्तेमाल करें, या JSON merge करें। मूल फ़ाइलें `.ccam-backup/<time>/` में जाती हैं, हटाए गए अकाउंट का फ़ोल्डर एक trash फ़ोल्डर में जाता है, और फ़ाइलें एक ही rename में बदली जाती हैं, इसलिए बीच में रुका write पुरानी फ़ाइल को जस का तस छोड़ देता है। `--dry-run` दिखाता है कि launch, sharing, migration, shortcuts, updates और repairs क्या करेंगे। `ccam doctor --fix` टूटे links और permissions ठीक करता है।

### अगले अकाउंट में वही बातचीत

साझा इतिहास के साथ, `ccam next --auto --continue` हर अकाउंट से एक साथ उसका उपयोग पूछता है, सबसे ज़्यादा गुंजाइश वाला अकाउंट खोलता है और बातचीत जारी रखता है। कब स्विच करना है, यह आप तय करते हैं; ccam background में आपके उपयोग पर नज़र नहीं रखता। किसी अकाउंट के 90% तक पहुँचने पर status line यह कमांड सुझाती है।

Scripts आपका इंतज़ार नहीं कर सकतीं, इसलिए `ccam --fallback work -p "…"` उनके लिए स्विच करता है: जब अकाउंट अपनी सीमा तक पहुँचता है, तो वही कमांड सबसे ज़्यादा गुंजाइश वाले अकाउंट के रूप में फिर से चलती है, और सिर्फ़ उस run का output छपता है जो पूरा हुआ। जो विफलता उपयोग की सीमा नहीं है, उसे जस का तस आगे दे दिया जाता है। समानांतर workers के लिए, `ccam pick` हर एक को गुंजाइश वाला अलग अकाउंट देता है, और जब किसी अकाउंट में गुंजाइश नहीं बचती तो दोनों code 75 के साथ बाहर निकलते हैं।

### हर फ़ोल्डर में सही अकाउंट

`ccam folder work ~/code/company` से उस फ़ोल्डर के अंदर बिना अकाउंट नाम वाला हर launch `work` इस्तेमाल करता है, जबकि आपके बाक़ी प्रोजेक्ट default अकाउंट पर रहते हैं। यह setting उसी `config.toml` में रहती है जिसे हर launch पढ़ता है, इसलिए इसमें कोई समय नहीं लगता।

### एक binary, जाँचे गए downloads

ccam लगभग 10 MB की एक static binary है, जिसके लिए कोई runtime इंस्टॉल नहीं करना पड़ता, Linux, macOS और Windows के लिए। install scripts (`install.sh`, `install.ps1`) और `ccam update` हर download को release के SHA-256 checksums से जाँचते हैं, और ccam अपने-आप कभी updates नहीं खोजता। इसमें कोई telemetry नहीं है: यह सिर्फ़ आपके अकाउंट्स का उपयोग जानने के लिए Anthropic को requests भेजता है, समानांतर में, 3 सेकंड की सीमा और 60 सेकंड के cache के साथ, और वे downloads करता है जो आप माँगते हैं।

## जहाँ दूसरे टूल ज़्यादा करते हैं

- claude-swap, clauth, cux, CAAM और CCSwitcher, किसी अकाउंट के उपयोग की तय सीमा पार होने पर चल रहे interactive session को अपने-आप स्विच कर सकते हैं। ccam `-p` runs को `--fallback` से स्विच करता है और interactive sessions आप पर छोड़ता है।
- Developer Mode के बिना Windows पर ccam फ़ोल्डर साझा करता है, लेकिन `settings.json` जैसी अकेली फ़ाइलें साझा नहीं कर सकता; claude-swap और CCS वहाँ उन्हें कॉपी कर देते हैं।
- CCSwitcher, Claude Switcher और claude-swap उपयोग को macOS menu bar में दिखाते हैं।
- CCS, clausona, CAAM, clauth, ccswap और Claude Switcher, Codex और दूसरे coding agents भी संभालते हैं।
- clausona हर अकाउंट के लिए extensions को चालू और बंद भी करता है (ccam का `ccam ext` उन्हें सिर्फ़ दिखाता है), local token लागत का हिसाब रखता है, और ऐसे Claude Code plugins देता है जो tmux और दूसरे terminals में workers के पूरे समूह चलाते हैं।
- CAAM, CCSwitcher और Claude Switcher, Homebrew से इंस्टॉल होते हैं, और clauth और CAAM अपने release checksums पर signature करते हैं। ccam अपनी script, एक `.deb` या `go install` से इंस्टॉल होता है, और checksums को बिना signature के जाँचता है।

## कौन-सा ठीक रहेगा?

अगर आप एक समय में एक ही अकाउंट इस्तेमाल करते हैं या चाहते हैं कि चल रहे sessions अपने-आप स्विच हों, तो लॉगिन बदलने वाला टूल चुनें। अगर आप अकाउंट्स को साथ-साथ चलाना, एक ही साझा सेटअप रखना, अगले अकाउंट में वही बातचीत जारी रखना चाहते हैं और नहीं चाहते कि Claude Code के साथ कुछ और चले, Linux, macOS या Windows पर, तो ccam चुनें। अगर आपको कुछ भी साझा नहीं करना है, तो [साधारण shell aliases](/compare/aliases) काफ़ी हैं।

कई अकाउंट इस्तेमाल करने की अनुमति है या नहीं, यह हर अकाउंट की शर्तों पर निर्भर करता है। ऐसा करने से पहले उन्हें जाँच लें।
