---
title: Ang ccam at iba pang Claude Code account switcher
description: Paano ikinukumpara ang ccam sa claude-swap, CCS, clauth, CAAM, cux at iba pang tool sa pag-switch ng Claude Code account, at aling paraan ang bagay saan.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

May dalawang paraan ang mga tool sa pag-switch ng Claude Code account. Karamihan ay nagpapalit ng login sa loob ng iisang `~/.claude` folder, kaya isang account lang ang aktibo sa bawat pagkakataon. Ang iba, kasama ang ccam, ay nagbibigay sa bawat account ng sariling folder gamit ang `CLAUDE_CONFIG_DIR`, kaya nananatiling naka-log in ang bawat account at puwedeng tumakbo nang sabay ang ilan. Nagdadagdag ang ccam ng shared na settings at MCP server, at itinutuloy nito ang isang usapan sa folder ng ibang account.

## Pagpapalit ng iisang login o magkakahiwalay na folder?

**Ang pagpapalit ng login** ay nagpapanatili ng iisang set ng settings, MCP server at kasaysayan, dahil iisa lang ang folder. Mabilis ang paglipat, at puwedeng gamitin ng tumatakbong session ang bagong login. Pero isang account lang ang aktibo sa bawat pagkakataon, para sa lahat ng terminal at editor nang sabay.

**Ang magkakahiwalay na folder** ang paraang [inilalarawan sa dokumentasyon ng Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Nananatiling naka-log in ang bawat account, at puwede mong patakbuhin nang magkatabi ang session para sa trabaho at ang personal na session. Ang kapalit nito: may sariling settings, MCP server at kasaysayan ang bawat folder sa simula, maliban kung may tool na nagshe-share ng mga ito.

## Paghahambing ng mga tool

Kinuha mula sa sariling dokumentasyon ng bawat project noong 2026-10-19. Ang "Hindi nakadokumento" ay nangangahulugang hindi ito binabanggit ng project, hindi na wala ang feature.

| Tool | Paraan | Sabay-sabay na account | 5-oras at lingguhang usage | Kapag naabot ang limit | Platform |
| --- | --- | --- | --- | --- | --- |
| ccam | Magkakahiwalay na folder | Oo | Oo, `ccam usage` | Ikaw ang nagpapatakbo ng `ccam next --auto --continue` para ituloy ang usapan sa account na may pinakamaraming natitira | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Nagpapalit ng login; magkakahiwalay na folder gamit ang `cswap run` | Oo, gamit ang `cswap run` | Oo | Nagpapalit ng login ang `cswap auto` sa isang threshold; tuloy ang tumatakbong session | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Magkakahiwalay na folder | Oo | Bahagya, sa isang macOS menu bar app | Ipinapasa ng mga account pool ang mga request sa susunod na account sa pamamagitan ng lokal na proxy | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Nagpapalit ng login; magkakahiwalay na folder gamit ang `clauth start` | Oo, gamit ang `clauth start` | Oo | Lumilipat pababa sa isang fallback chain | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Nagpapalit ng login, o mga isolated na profile | Oo, gamit ang mga isolated na profile | Oo | Inuulit ng `caam run` ang mga non-interactive na command sa ibang account | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Nagpapalit ng login at binabalot ang `claude` | Hindi nakadokumento | Oo | Nagpapalit ng login at itinutuloy ang parehong usapan | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Nagpapalit ng login | Hindi | Oo | Hindi nakadokumento | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Nagpapalit ng login | Hindi | Oo | Opsyonal na paglipat sa 100% | macOS |

Ang ccswap ay fork ng claude-swap na namamahala rin ng mga Codex login. Ilan sa mga tool na ito ay humahawak din ng ibang coding agent, gaya ng Codex o Gemini.

## Ang pinagtutuunan ng ccam

- Nananatiling naka-log in ang bawat account. Hindi kailanman nala-log out ang isang account sa paglipat, at puwedeng tumakbo nang sabay ang `ccam work` at `ccam personal` sa dalawang terminal.
- Shine-share sa pagitan ng mga folder ang mga MCP server, `settings.json`, `CLAUDE.md`, agent, skill at plugin sa pamamagitan ng mga link, bawat item at bawat account. Hindi kailanman shine-share ang mga login.
- Kapag shared ang kasaysayan, binubuksan ng iisang command, `ccam next --auto --continue`, ang parehong usapan sa folder ng account na may pinakamaraming natitira.
- Sinisimulan ng ccam ang Claude Code at pinapalitan ang sarili nito ng Claude Code, kaya walang tumatakbo sa background. Hindi nito binabantayan ang usage mo at hindi ito kusang lumilipat; ikaw ang nagpapasya kung kailan.

## Kung saan mas marami ang nagagawa ng ibang tool

- Tumatakbo ang claude-swap, CCS, clauth, CAAM at cux sa Windows nang walang WSL. Hindi ang ccam.
- Kusang lumilipat ang claude-swap, clauth at cux kapag lumampas ang isang account sa usage threshold. Iniiwan ng ccam sa iyo ang paglipat.
- Ipinapakita ng CCSwitcher, Claude Switcher at claude-swap ang usage sa macOS menu bar.
- Namamahala rin ang CCS, CAAM, ccswap at Claude Switcher ng Codex at iba pang coding agent.

## Alin ang bagay sa iyo?

Pumili ng tool na nagpapalit ng login kung isang account lang ang ginagamit mo sa bawat pagkakataon, gusto mong awtomatiko ang paglipat, o nasa Windows ka nang walang WSL. Piliin ang ccam kung gusto mo ng mga account na magkatabi, iisang shared na setup at ang parehong usapan sa susunod na account. Kung wala kang kailangang i-share, sapat na ang [simpleng mga shell alias](/compare/aliases).

Nakadepende sa terms ng bawat account kung puwede kang gumamit ng maraming account. Tingnan muna ang mga ito bago mo gawin.
