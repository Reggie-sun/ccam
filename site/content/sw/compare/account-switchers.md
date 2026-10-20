---
title: ccam na vibadilisha akaunti vingine vya Claude Code
description: Jinsi ccam inavyolinganishwa na claude-swap, CCS, clauth, CAAM, cux na zana nyingine za kubadili akaunti za Claude Code, na njia ipi inafaa matumizi yapi.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: d058da937218
---

Zana za kubadili akaunti za Claude Code hufuata mojawapo ya njia mbili. Nyingi hubadilisha taarifa za kuingia ndani ya folda moja ya `~/.claude`, kwa hiyo akaunti moja tu hutumika kwa wakati mmoja. Nyingine, ikiwemo ccam, huipa kila akaunti folda yake kwa `CLAUDE_CONFIG_DIR`, kwa hiyo kila akaunti hubaki imeingia na kadhaa zinaweza kuendeshwa kwa pamoja. ccam huongeza mipangilio na seva za MCP zinazoshirikiwa, na huendeleza mazungumzo katika folda ya akaunti nyingine.

## Kubadilisha kuingia kumoja au folda tofauti?

**Kubadilisha kuingia** huweka seti moja ya mipangilio, seva za MCP na historia, kwa sababu kuna folda moja tu. Kubadili ni haraka, na kipindi kinachoendelea kinaweza kuchukua kuingia mpya. Lakini akaunti moja tu hutumika kwa wakati mmoja, kwa kila terminal na kihariri kwa pamoja.

Ili kubadilisha, zana huhifadhi nakala ya taarifa za kuingia za kila akaunti nyingine na kuiandika mahali pake. Claude Code hubadilisha refresh token ya kuingia kila inapoifanya upya, kwa hiyo kila nakala iliyohifadhiwa lazima ilingane na vipindi vinavyoendelea. Akaunti zinazotolewa au kuhitaji `/login` tena baada ya kubadili ni miongoni mwa matatizo ambayo miradi hii hurekebisha mara nyingi zaidi.

**Folda tofauti** ndiyo njia ambayo [nyaraka za Claude Code zinaeleza](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Kila akaunti hubaki imeingia, na unaweza kuendesha kipindi cha kazi na cha binafsi sambamba. Taarifa za kuingia za kila akaunti hubaki katika folda yake, ambapo Claude Code huzifanya upya kama kawaida. Gharama yake ni kwamba kila folda huanza na mipangilio, seva za MCP na historia yake yenyewe, isipokuwa zana fulani izishiriki.

## Jinsi zana zinavyolinganishwa

Imechukuliwa kutoka nyaraka na msimbo chanzo wa kila mradi tarehe 2026-10-20. "Haijaelezwa" inamaanisha mradi hausemi, si kwamba kipengele hakipo.

| Zana | Njia | Akaunti kwa pamoja | Matumizi ya saa 5 na ya wiki | Kikomo kinapofikiwa | Mifumo |
| --- | --- | --- | --- | --- | --- |
| ccam | Folda tofauti | Ndiyo | Ndiyo, `ccam usage` | Unaendesha `ccam next --auto --continue` ili kuendeleza mazungumzo katika akaunti yenye nafasi kubwa zaidi; `ccam --fallback` huendesha tena kazi za `-p` humo yenyewe | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Hubadilisha kuingia; folda tofauti kwa `cswap run` | Ndiyo, kwa `cswap run` | Ndiyo | `cswap auto` hubadilisha kuingia kwenye kiwango fulani; kipindi kinachoendelea huendelea | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Folda tofauti | Ndiyo | Kwa sehemu | Makundi ya akaunti hupeleka maombi kwa akaunti inayofuata kupitia proxy ya ndani | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Hubadilisha kuingia; folda tofauti kwa `clauth start` | Ndiyo, kwa `clauth start` | Ndiyo | Hubadili kufuata mnyororo wa akaunti mbadala | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Hubadilisha kuingia, au profile zilizotengwa | Ndiyo, kwa profile zilizotengwa | Ndiyo | `caam run` hujaribu tena amri zisizo shirikishi kwenye akaunti nyingine | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Hubadilisha kuingia na kuzunguka `claude` | Hapana, akaunti moja hai | Ndiyo | Hubadilisha kuingia na kuendeleza mazungumzo yale yale | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Hubadilisha kuingia | Hapana | Ndiyo | Kubadili kwa hiari kwenye kiwango fulani | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Hubadilisha kuingia | Hapana | Ndiyo | Kubadili kwa hiari kwenye 100% | macOS |

ccswap ni fork ya claude-swap ambayo pia husimamia kuingia kwa Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) na [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) pia huipa kila akaunti folda yake. Zana kadhaa kati ya hizi pia hushughulikia mawakala wengine wa kuandika msimbo, kama Codex au Gemini.

## Jinsi zilivyojengwa

| Zana | Imeandikwa kwa | Wakati Claude Code inaendeshwa | Hunakili taarifa za kuingia kati ya maeneo | Kile akaunti zinaweza kushiriki |
| --- | --- | --- | --- | --- |
| ccam | Go, binary moja tuli | Hakuna: ccam hujibadilisha na `claude` (kwenye Windows, ambayo haiwezi, huisubiri) | Hapana | Kila kipengele peke yake: cha ndani au kundi lenye jina, kwa kila akaunti |
| claude-swap | Python 3.12 au mpya zaidi, husakinishwa kwa uv au pipx | Hakuna kwa `cswap run`; `cswap auto` na programu ya menu bar huendelea kuendeshwa | Ndiyo, inapobadilisha kuingia | Kwa `cswap run`: mipangilio, `CLAUDE.md` na viendelezi kwa pamoja, historia kwa hiari; seva za MCP hunakiliwa kila inapoanzishwa |
| CCS | Node.js 18 au mpya zaidi, au Bun | CCS hubaki kama mchakato mzazi; makundi huongeza proxy ya ndani | Hapana kwa folda zake za akaunti; makundi huhifadhi taarifa zao za kuingia | Inashirikiwa au ya ndani kwa kila akaunti, historia katika makundi ya muktadha; seva za MCP hunakiliwa |
| clauth | Rust | clauth hubaki kama mchakato mzazi | Ndiyo, inapobadilisha kuingia | `~/.claude` yote, au hakuna |
| CAAM | Go | Inategemea hali; daemon ya hiari | Ndiyo, katika hali yake ya vault | Huamuliwa na hali; seva za MCP hunakiliwa kila inapoanzishwa |
| cux | Go | Kifuniko (wrapper) hubaki kuzunguka `claude` | Ndiyo | Kila kitu: folda moja |
| CCSwitcher, Claude Switcher | Programu za macOS (Swift, Python) | Programu ya menu bar | Ndiyo | Kila kitu: folda moja |

## Mambo ambayo ccam hufanya vizuri

### Taarifa za kuingia hubaki pale Claude Code ilipoziweka

ccam kamwe hainakili taarifa za kuingia za Claude kutoka mahali pamoja hadi pengine. Taarifa za kuingia za kila akaunti hubaki katika folda yake, au kwenye macOS katika kipengee cha Keychain cha folda hiyo, na Claude Code ndiyo programu inayozifanya upya. Hakuna nakala iliyohifadhiwa inayoweza kupoteza ulinganifu. `~/.claude` iliyopo iliyosajiliwa kwa `ccam import` huhifadhi kipengee cha Keychain kinachotumiwa na `claude` ya kawaida, kwa hiyo hubaki imeingia kwa njia zote mbili.

ccam hufanya upya tokeni iliyokwisha muda yenyewe tu unapoendesha `ccam usage --refresh`. Huchukua kufuli kwanza na hukataa wakati akaunti hiyo inaendeshwa, kwa hiyo kamwe haibadilishi tokeni chini ya kipindi kinachoendelea.

### Hakuna kinachoendeshwa kati yako na Claude Code

ccam husoma faili moja, hukagua viungo vyake na hujibadilisha na `claude` (`exec`). Unapofanya kazi hakuna mchakato wa ccam, kifuniko, proxy wala daemon: ingizo la kawaida, ishara na msimbo wa kutoka ni vya Claude Code, na maombi yake huenda moja kwa moja kwa Anthropic. Windows haiwezi kubadilisha mchakato, kwa hiyo huko ccam huanzisha Claude Code kwa console ile ile na kusubiri msimbo wake wa kutoka, bila kufanya kitu kingine.

Kuanzisha akaunti kwa njia hii huongeza takriban 4 ms, chini ya 6 ms kwenye asilimia ya 95, na kumbukumbu ya juu kabisa ya 8 MB ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) kwa hyperfine, mizunguko 300). Kwenye kompyuta mpakato ile ile, kuanzisha mfasiri wa Python peke yake huchukua 13 ms na Node.js 31 ms, kabla programu haijafanya chochote.

### Shiriki kipengele kwa kipengele, akaunti kwa akaunti

Seva za MCP, `settings.json`, `CLAUDE.md`, viendelezi (agents, commands, skills, plugins) na historia ya mazungumzo ni vipengele vitano tofauti. Kila akaunti huweka kila kipengele kuwa `local` au kundi, kwa hiyo akaunti mbili za kazi zinaweza kushiriki historia katika kundi la `work` wakati akaunti ya binafsi inabaki na yake, na zote tatu hushiriki seva zile zile za MCP.

- Seva za MCP hupitishwa kwa `--mcp-config` wakati wa kuanzisha, kwa hiyo kwa chaguo-msingi ccam kamwe haihariri `.claude.json`, faili ambalo Claude Code huliandika upya inapoendeshwa. `ccam mcp add` huongeza seva mara moja kwa kundi zima.
- Vipengele vingine ni viungo vinavyoelekea folda ya kundi, kwa hiyo badiliko lililofanywa katika akaunti moja huonekana katika nyingine mara moja. Hakuna kinachonakiliwa wakati wa kuanzisha.
- Taarifa za kuingia, `.claude.json` na hali ya uendeshaji hazishirikiwi kamwe.

### Kila badiliko huonyeshwa kwanza na kuhifadhiwa nakala

`ccam share` na `ccam migrate` huchapisha mpango na kuuliza kabla ya kuhamisha chochote. Faili linapotofautiana kati ya akaunti na kundi lake, unachagua: kubaki na la kundi, kutumia lako, au kuunganisha JSON. Faili za asili huenda `.ccam-backup/<time>/`, folda ya akaunti iliyoondolewa huenda kwenye folda ya taka, na faili hubadilishwa kwa rename moja, kwa hiyo uandishi uliokatizwa huacha faili la zamani. `--dry-run` huonyesha kile ambacho kuanzisha, kushiriki, kuhamisha, njia za mkato, masasisho na marekebisho yangefanya. `ccam doctor --fix` hurekebisha viungo na ruhusa zilizoharibika.

### Mazungumzo yale yale katika akaunti inayofuata

Historia ikishirikiwa, `ccam next --auto --continue` huuliza kila akaunti matumizi yake kwa wakati mmoja, hufungua ile yenye nafasi kubwa zaidi na kuendeleza mazungumzo. Wewe ndiye unaamua lini kubadili; ccam haifuatilii matumizi yako chinichini. Mstari wa hali hupendekeza amri hiyo akaunti inapofikia 90%.

Script haziwezi kukusubiri, kwa hiyo `ccam --fallback work -p "…"` hubadili kwa ajili yao: akaunti inapofikia kikomo chake, amri ile ile huendeshwa tena kama akaunti yenye nafasi kubwa zaidi, na matokeo ya uendeshaji uliokamilika pekee ndiyo huchapishwa. Kushindwa kusikotokana na kikomo cha matumizi hupitishwa kama kulivyo.

### Akaunti sahihi katika kila folda

`ccam folder work ~/code/company` hufanya kila uanzishaji usio na jina la akaunti ndani ya folda hiyo utumie `work`, wakati miradi yako mingine inabaki na akaunti chaguo-msingi. Mpangilio huu uko katika `config.toml` ile ile ambayo kila uanzishaji husoma, kwa hiyo hauongezi muda.

### Binary moja, vipakuliwa vilivyokaguliwa

ccam ni binary tuli ya takriban 10 MB bila runtime ya kusakinisha, kwa Linux, macOS na Windows. Hati za usakinishaji (`install.sh`, `install.ps1`) na `ccam update` hukagua kila kipakuliwa dhidi ya checksum za SHA-256 za toleo, na ccam kamwe haitafuti masasisho yenyewe. Haina telemetry: maombi pekee inayotuma ni maswali ya matumizi kwa Anthropic kwa ajili ya akaunti zako, kwa sambamba na kikomo cha sekunde 3 na akiba ya sekunde 60, pamoja na vipakuliwa unavyoomba.

## Pale zana nyingine zinapofanya zaidi

- claude-swap, clauth, cux, CAAM na CCSwitcher zinaweza kubadili zenyewe kipindi shirikishi kinachoendelea akaunti inapovuka kiwango fulani cha matumizi. ccam hubadili uendeshaji wa `-p` kwa `--fallback` na hukuachia vipindi shirikishi.
- Kwenye Windows bila Developer Mode, ccam hushiriki folda lakini haiwezi kushiriki faili moja moja kama `settings.json`; claude-swap na CCS huyanakili huko badala yake.
- CCSwitcher, Claude Switcher na claude-swap huonyesha matumizi kwenye menu bar ya macOS.
- CCS, CAAM, clauth, ccswap na Claude Switcher pia husimamia Codex na mawakala wengine wa kuandika msimbo.
- CAAM, CCSwitcher na Claude Switcher husakinishwa kwa Homebrew, na clauth na CAAM hutia sahihi checksum za matoleo yao. ccam husakinishwa kwa hati yake, `.deb` au `go install`, na hukagua checksum bila sahihi.

## Ipi inafaa?

Chagua zana ya kubadilisha kuingia ikiwa unatumia akaunti moja kwa wakati mmoja au unataka vipindi vinavyoendelea vibadili kiotomatiki. Chagua ccam ikiwa unataka akaunti sambamba, usanidi mmoja unaoshirikiwa, mazungumzo yale yale katika akaunti inayofuata na kutokuwa na kitu kinachoendeshwa kando ya Claude Code, kwenye Linux, macOS au Windows. Ikiwa huhitaji kushiriki chochote, [alias za kawaida za shell](/compare/aliases) zinatosha.

Iwapo unaruhusiwa kutumia akaunti kadhaa inategemea masharti ya kila akaunti. Yakague kabla ya kufanya hivyo.
