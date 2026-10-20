---
title: Maraming Claude Code account sa iisang computer
description: Magpatakbo ng maraming Claude Code account nang magkatabi sa CLAUDE_CONFIG_DIR, mano-mano sa shell alias o awtomatiko sa ccam, at lumipat sa pagitan nila.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2f4deb40c81
---

Bigyan ang bawat account ng sariling configuration folder. Binabasa ng Claude Code ang folder mula sa environment variable na `CLAUDE_CONFIG_DIR`, at may sariling login, settings at kasaysayan ang bawat folder, kaya nananatiling naka-log in ang bawat account at puwede mo silang patakbuhin nang sabay. Puwede mo itong i-set up nang mano-mano gamit ang shell alias, o hayaang gawin ito ng ccam.

## Paano ito i-set up nang mano-mano?

Ito ang paraang [inilalarawan sa dokumentasyon ng Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Magdagdag ng isang linya bawat karagdagang account sa `~/.bashrc` o `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Magbukas ng bagong terminal at patakbuhin ang `claude-work`. Hihilingin ng Claude Code na mag-log in ka, dahil bago ang folder. Mula noon, `claude` ang una mong account at `claude-work` ang pangalawa. Magdagdag ng isang alias bawat account, bawat isa may sariling folder.

Sa Windows na walang WSL, ganoon din ang ginagawa ng isang PowerShell function sa iyong profile (`notepad $PROFILE`):

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Ano ang kulang sa paraang alias?

Gumagana ito nang maayos para sa dalawang account. Mga dapat malaman:

- Walang laman ang bawat folder sa simula: walang settings, walang `CLAUDE.md`, walang MCP server, walang agent o skill. Ikaw ang kokopya o magli-link ng kailangan mo, at ikaw rin ang magpapanatiling magkakatugma ng mga kopya.
- Umiiral lang ang mga alias sa mga interactive shell. Hindi sila nakikita ng mga script, Makefile at iba pang program.
- Nasa `.claude.json` ang mga MCP server, katabi ng login, kaya hindi sila basta mali-link sa pagitan ng mga folder.
- Nananatili ang mga usapan sa account kung saan sila nagsimula, kaya hindi mo maitutuloy ang isa sa ibang account.
- Walang anumang nagpapakita ng usage ng lahat ng account nang sabay-sabay.

## Puwede ba akong mag-switch ng account gamit ang /login?

Oo, kung paminsan-minsan ka lang lumilipat. Kapag nag-type ka ng `/login` sa loob ng Claude Code, [napapalitan ang account](https://code.claude.com/docs/en/quickstart) ng folder na ginagamit mo. Isang account lang ang naka-log in sa bawat pagkakataon, kaya kailangan mong mag-log in ulit sa bawat paglipat, at hindi mo mapapatakbo nang magkatabi ang dalawang account. Naiiwasan ng magkakahiwalay na folder ang dalawang problemang ito.

## Paano ito i-set up gamit ang ccam?

Ina-automate ng ccam ang parehong ideya at idinadagdag ang mga kulang na bahagi:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Ngayon, sinisimulan ng `ccam personal` at `ccam work` ang dalawang account, at ganoon din ang `ccwork`, ang shortcut na ibinigay mo sa work account. Dahil totoong mga program sila at hindi mga alias, gumagana rin sila sa mga script. Patakbuhin ang `ccam` nang mag-isa para pumili ng account mula sa listahan, at ang `ccam use work` para gawing default ang isang account.

Shared na mula sa simula ang settings, `CLAUDE.md`, mga MCP server at extension, at puwede mong i-on ang shared na kasaysayan ng usapan para maituloy ang isang usapan sa kabilang account gamit ang `ccam next --continue`. Ipinapakita ng `ccam usage` ang 5-oras at lingguhang usage ng bawat account nang sabay-sabay.

Kung gumagamit ka na ng `~/.claude`, patakbuhin ang `ccam import` sa halip na idagdag ulit ang account na iyon, o ang `ccam migrate` para kopyahin ang settings, mga MCP server at extension nito sa shared na settings.

## Paano ko malalaman kung aling account ang tumatakbo?

Sa loob ng Claude Code, ipinapakita ng `/status` ang account na naka-log in. Sa ccam, patakbuhin nang isang beses ang `ccam statusline install` at palaging ipapakita ng status line ang pangalan ng account at ang usage nito, halimbawa `[work] · Opus · 5h 24% · 7d 41%`.

## Gumagana ba ito sa Windows at WSL?

Tumatakbo ang Claude Code [nang native sa Windows at sa loob ng WSL](https://code.claude.com/docs/en/setup). Gumagana ang `CLAUDE_CONFIG_DIR` sa pareho, kaya sapat na ang alias o ang PowerShell function sa itaas sa native na Windows. Tumatakbo ang ccam sa Linux at macOS; sa Windows, i-install at gamitin ito sa loob ng WSL, kung saan gumagana ito gaya ng sa Linux.

## Pag-troubleshoot

**Palaging hinihiling ng Claude Code na mag-log in ako.** Laging humihingi ito ng login sa unang pagbukas sa bagong folder. Kung paulit-ulit itong humihingi, nag-expire o naalis ang login: patakbuhin ang `ccam login work` (o patakbuhin nang isang beses ang alias at mag-log in).

**Sa iisang account lang lumalabas ang isang MCP server.** Pag-aari ng account na iyon ang mga server na idinagdag sa loob ng session gamit ang `claude mcp add`. Ilipat ang isa sa lahat ng account gamit ang `ccam mcp import work <server>`.

**Magkaiba ang settings ng mga account.** Ipinapakita ng `ccam share work` kung ano ang shine-share ng account; shine-share naman ng `ccam share work --settings on` ang settings nito, matapos ipakita kung ano ang magbabago. Inaayos ng `ccam doctor --fix` ang mga sirang link at nagtatabi ito ng backup ng bawat file na binabago nito.

**Iisang usapan ang binubuksan ng dalawang account.** Iwasang gumawa sa iisang usapan sa dalawang account nang sabay. Babalaan ka ng ccam kapag sinubukan mo.

## Alin ang dapat kong piliin?

Para sa dalawang account na walang shine-share, sapat na ang alias. Kapag gusto mo ng shared na settings at MCP server, mas maraming account, mabilisang tingin sa usage o pagtutuloy ng mga usapan, inaalis ng ccam ang abala ng pag-aasikaso. Inililista ng [paghahambing sa mga alias](/compare/aliases) ang mga pagkakaiba nang magkatabi, at tinatalakay ng [ccam at iba pang account switcher](/compare/account-switchers) ang mga tool na nagpapalit ng iisang login.

Nakadepende sa terms ng bawat account kung puwede kang gumamit ng maraming account. Tingnan muna ang mga ito bago mo gawin.
