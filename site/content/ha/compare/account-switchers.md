---
title: ccam da sauran masu sauya asusun Claude Code
description: Yadda ccam ya kwatanta da claude-swap, CCS, clauth, CAAM, cux da sauran kayan aikin sauya asusun Claude Code, da wace hanya ta dace da wane amfani.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: d058da937218
---

Kayan aikin sauya asusun Claude Code suna bin ɗaya daga cikin hanyoyi biyu. Yawancinsu suna musanya bayanan shiga a cikin foldar `~/.claude` guda ɗaya, don haka asusu ɗaya ne kawai yake aiki a lokaci guda. Wasu, ciki har da ccam, suna ba kowane asusu foldarsa da `CLAUDE_CONFIG_DIR`, don haka kowane asusu yana nan a shiga kuma da dama za su iya aiki tare. ccam yana ƙara saituna da sabobin MCP da aka raba, kuma yana ci gaba da tattaunawa a foldar wani asusun.

## Musanya shiga ɗaya ko foldoji daban-daban?

**Musanya shiga** yana riƙe saiti ɗaya na saituna, sabobin MCP da tarihi, domin folda ɗaya ce kawai. Sauyawa yana da sauri, kuma zaman da ke gudana zai iya ɗaukar sabuwar shigar. Amma asusu ɗaya ne kawai yake aiki a lokaci guda, ga kowane terminal da editor a lokaci ɗaya.

Don musanyawa, kayan aikin yana ajiye kwafin shigar kowane asusun daban, kuma yana rubuta shi a wurinsa. Claude Code yana maye gurbin refresh token na shiga a duk lokacin da yake sabunta shi, don haka dole ne kowane kwafi da aka ajiye ya kasance daidai da zaman da ke gudana. Asusan da aka fitar da su ko suke buƙatar `/login` kuma bayan sauyawa suna cikin matsalolin da waɗannan projects suka fi gyarawa.

**Foldoji daban-daban** ita ce hanyar da [takardun Claude Code suka bayyana](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Kowane asusu yana nan a shiga, kuma za ka iya gudanar da zaman aiki da na kai gefe da gefe. Kowace shiga tana zama a foldarta, inda Claude Code yake sabunta ta kamar yadda aka saba. Farashin shi ne kowace folda tana farawa da saitunanta, sabobin MCP da tarihinta, sai dai idan wani kayan aiki ya raba su.

## Yadda kayan aikin suka kwatanta

An ɗauko daga takardun kowane project da source code ɗinsa a ranar 2026-10-20. "Ba a bayyana ba" yana nufin project ɗin bai faɗa ba, ba wai babu fasalin ba.

| Kayan aiki | Hanya | Asusai a lokaci ɗaya | Amfanin awa 5 da na mako | Idan an kai iyaka | Tsaruka |
| --- | --- | --- | --- | --- | --- |
| ccam | Foldoji daban-daban | Ee | Ee, `ccam usage` | Kai ne za ka gudanar da `ccam next --auto --continue` don ci gaba da tattaunawar a asusun da ya fi sarari; `ccam --fallback` yana sake gudanar da ayyukan `-p` a can da kansa | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Yana musanya shiga; foldoji daban-daban da `cswap run` | Ee, da `cswap run` | Ee | `cswap auto` yana musanya shiga a wani mataki; zaman da ke gudana yana ci gaba | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Foldoji daban-daban | Ee | Wani ɓangare | Rukunin asusai suna miƙa buƙatu ga asusu na gaba ta proxy na cikin gida | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Yana musanya shiga; foldoji daban-daban da `clauth start` | Ee, da `clauth start` | Ee | Yana sauyawa bisa jerin asusan maye | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Yana musanya shiga, ko profiles da aka ware | Ee, da profiles da aka ware | Ee | `caam run` yana sake gwada umarnan da ba na mu'amala ba a wani asusu | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Yana musanya shiga kuma yana naɗe `claude` | A'a, asusu ɗaya mai aiki | Ee | Yana musanya shiga kuma yana ci gaba da tattaunawar | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Yana musanya shiga | A'a | Ee | Sauyawa na zaɓi a wani mataki | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Yana musanya shiga | A'a | Ee | Sauyawa na zaɓi a 100% | macOS |

ccswap fork ne na claude-swap wanda kuma yake sarrafa shigar Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) da [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) suma suna ba kowane asusu foldarsa. Da yawa daga cikin waɗannan kayan aikin suna kuma kula da wasu coding agents, kamar Codex ko Gemini.

## Yadda aka gina su

| Kayan aiki | Harshen da aka rubuta da shi | Yayin da Claude Code ke aiki | Yana kwafa shiga tsakanin wurare | Abin da asusai za su iya raba |
| --- | --- | --- | --- | --- |
| ccam | Go, static binary ɗaya | Babu komai: ccam yana maye gurbin kansa da `claude` (a Windows, wanda ba zai iya ba, yana jiran sa) | A'a | Kowane abu daban: na gida ko rukuni mai suna, ga kowane asusu |
| claude-swap | Python 3.12 ko sama, ana shigar da shi da uv ko pipx | Babu komai da `cswap run`; `cswap auto` da manhajar menu bar suna ci gaba da aiki | Ee, idan yana musanya shiga | Da `cswap run`: saituna, `CLAUDE.md` da extensions tare, tarihi na zaɓi; ana kwafa sabobin MCP a kowane ƙaddamarwa |
| CCS | Node.js 18 ko sama, ko Bun | CCS yana zama parent process; rukunin asusai suna ƙara proxy na cikin gida | A'a ga foldojin asusansa; rukunin asusai suna riƙe shigarsu | Raba ko na gida ga kowane asusu, tarihi a cikin context groups; ana kwafa sabobin MCP |
| clauth | Rust | clauth yana zama parent process | Ee, idan yana musanya shiga | Duk `~/.claude`, ko babu komai |
| CAAM | Go | Ya dogara da yanayin aiki; da daemon na zaɓi | Ee, a yanayin vault ɗinsa | Yanayin aiki ne ke tsayar da shi; ana kwafa sabobin MCP a kowane ƙaddamarwa |
| cux | Go | Wrapper yana zama a kewaye da `claude` | Ee | Komai: folda ɗaya |
| CCSwitcher, Claude Switcher | Manhajojin macOS (Swift, Python) | Manhajar menu bar | Ee | Komai: folda ɗaya |

## Abin da ccam ya kware a kai

### Shiga tana zama inda Claude Code ya ajiye ta

ccam ba ya taɓa kwafa shigar Claude daga wani wuri zuwa wani. Shigar kowane asusu tana zama a foldarsa, ko a macOS a cikin abin Keychain na wannan foldar, kuma Claude Code ne shirin da ke sabunta ta. Babu kwafin da aka ajiye da zai iya rasa daidaito. `~/.claude` da ke akwai wanda aka yi rajistarsa da `ccam import` yana riƙe abin Keychain da `claude` na yau da kullum ke amfani da shi, don haka yana nan a shiga a kowace hanya.

ccam yana sabunta token da ya ƙare da kansa ne kawai idan ka gudanar da `ccam usage --refresh`. Yana fara ɗaukar kulle, kuma yana ƙi idan wannan asusun yana aiki, don haka ba ya taɓa canza token a ƙarƙashin zaman da ke gudana.

### Babu abin da ke gudana tsakaninka da Claude Code

ccam yana karanta fayil ɗaya, yana duba mahaɗansa, sannan yana maye gurbin kansa da `claude` (`exec`). Yayin da kake aiki babu process na ccam, wrapper, proxy ko daemon: standard input, signals da exit code na Claude Code ne, kuma buƙatunsa suna zuwa kai tsaye ga Anthropic. Windows ba zai iya maye gurbin process ba, don haka a can ccam yana fara Claude Code da console ɗaya kuma yana jiran exit code ɗinsa, ba tare da yin wani abu ba.

Fara asusu ta wannan hanya yana ƙara kusan 4 ms, ƙasa da 6 ms a 95th percentile, da ƙwaƙwalwar ƙoli ta 8 MB ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) da hyperfine, gudu 300). A wannan kwamfutar tafi-da-gidanka, fara Python interpreter kaɗai yana ɗaukar 13 ms, Node.js kuma 31 ms, kafin shiri ya yi komai.

### Raba abu bayan abu, asusu bayan asusu

Sabobin MCP, `settings.json`, `CLAUDE.md`, extensions (agents, commands, skills, plugins) da tarihin tattaunawa abubuwa biyar ne daban-daban. Kowane asusu yana saita kowane abu zuwa `local` ko wani rukuni, don haka asusan aiki biyu za su iya raba tarihi a rukunin `work` yayin da asusun kai yake riƙe nasa, kuma dukansu uku suna raba sabobin MCP iri ɗaya.

- Ana ba da sabobin MCP da `--mcp-config` a lokacin ƙaddamarwa, don haka ta tsohuwa ccam ba ya taɓa gyara `.claude.json`, fayil ɗin da Claude Code ke sake rubutawa yayin da yake aiki. `ccam mcp add` yana ƙara saba sau ɗaya ga duk rukunin.
- Sauran abubuwan mahaɗai ne zuwa foldar rukunin, don haka canjin da aka yi a wani asusu yana bayyana a sauran nan take. Ba a kwafa komai a lokacin ƙaddamarwa.
- Ba a taɓa raba shiga, `.claude.json` ko yanayin aiki ba.

### Ana nuna kowane canji da farko kuma ana ajiye kwafinsa

`ccam share` da `ccam migrate` suna buga tsari kuma suna tambaya kafin su motsa komai. Idan fayil ya bambanta tsakanin asusu da rukuninsa, kai ne kake zaɓa: ka riƙe na rukunin, ka yi amfani da naka, ko ka haɗa JSON ɗin. Fayilolin asali suna zuwa `.ccam-backup/<time>/`, foldar asusun da aka cire tana zuwa foldar shara, kuma ana maye gurbin fayiloli da rename ɗaya, don haka rubutun da aka katse yana barin tsohon fayil. `--dry-run` yana nuna abin da ƙaddamarwa, rabawa, ƙaura, gajerun hanyoyi, sabuntawa da gyare-gyare za su yi. `ccam doctor --fix` yana gyara mahaɗai da izini da suka lalace.

### Wannan tattaunawar a asusu na gaba

Idan an raba tarihi, `ccam next --auto --continue` yana tambayar kowane asusu amfaninsa a lokaci ɗaya, yana buɗe wanda ya fi sarari kuma yana ci gaba da tattaunawar. Kai ne kake yanke shawarar lokacin sauyawa; ccam ba ya sa ido kan amfaninka a bayan fage. Layin matsayi yana ba da shawarar umarnin da zarar wani asusu ya kai 90%.

Rubutun umarni (scripts) ba za su iya jiranka ba, don haka `ccam --fallback work -p "…"` yana sauyawa a madadinsu: idan asusun ya kai iyakarsa, umarnin iri ɗaya yana sake gudana a matsayin asusun da ya fi sarari, kuma sakamakon gudun da ya kammala kaɗai ake bugawa. Gazawar da ba ta iyakar amfani ba ana miƙa ta yadda take.

### Asusun da ya dace a kowace folda

`ccam folder work ~/code/company` yana sa kowane ƙaddamarwa ba tare da sunan asusu ba a cikin wannan foldar ya yi amfani da `work`, yayin da sauran projects ɗinka suke ci gaba da amfani da tsohon asusu. Saitin yana cikin `config.toml` ɗin da kowane ƙaddamarwa ke karantawa, don haka ba ya ɗaukar ƙarin lokaci.

### Binary ɗaya, saukarwa da aka tantance

ccam static binary ne mai kusan 10 MB ba tare da runtime da za a shigar ba, don Linux, macOS da Windows. Rubutun shigarwa (`install.sh`, `install.ps1`) da `ccam update` suna tantance kowace saukarwa da SHA-256 checksums na release, kuma ccam ba ya taɓa neman sabuntawa da kansa. Ba shi da telemetry: buƙatun da yake yi kawai su ne tambayoyin amfani ga Anthropic don asusanka, a lokaci ɗaya tare da iyakar daƙiƙa 3 da cache na daƙiƙa 60, da saukarwar da ka nema.

## Inda sauran kayan aikin suka fi

- claude-swap, clauth, cux, CAAM da CCSwitcher za su iya sauya zaman mu'amala da ke gudana da kansu idan wani asusu ya wuce wani matakin amfani. ccam yana sauya gudun `-p` da `--fallback` kuma yana barin zaman mu'amala a hannunka.
- A Windows ba tare da Developer Mode ba, ccam yana raba foldoji amma ba zai iya raba fayiloli guda ɗaya kamar `settings.json` ba; claude-swap da CCS suna kwafa su a can maimakon haka.
- CCSwitcher, Claude Switcher da claude-swap suna nuna amfani a menu bar na macOS.
- CCS, CAAM, clauth, ccswap da Claude Switcher suna kuma sarrafa Codex da wasu coding agents.
- Ana shigar da CAAM, CCSwitcher da Claude Switcher da Homebrew, kuma clauth da CAAM suna sa hannu kan checksums na release ɗinsu. Ana shigar da ccam da rubutunsa, `.deb` ko `go install`, kuma yana tantance checksums ba tare da sa hannu ba.

## Wanne ya dace?

Zaɓi kayan aikin da ke musanya shiga idan kana amfani da asusu ɗaya a lokaci guda ko kana son zaman da ke gudana su sauya ta atomatik. Zaɓi ccam idan kana son asusai gefe da gefe, saiti ɗaya da aka raba, wannan tattaunawar a asusu na gaba, kuma babu abin da ke gudana kusa da Claude Code, a Linux, macOS ko Windows. Idan ba ka buƙatar raba komai, [shell alias kawai](/compare/aliases) ya isa.

Ko an yarda ka yi amfani da asusai da dama ya dogara da sharuɗɗan kowane asusu. Duba su kafin ka yi haka.
