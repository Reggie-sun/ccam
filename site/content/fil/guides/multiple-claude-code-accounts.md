---
title: Paano gumamit ng dalawang Claude Code account sa iisang computer
description: Magpatakbo ng maraming Claude Code account nang magkatabi gamit ang CLAUDE_CONFIG_DIR, mano-mano gamit ang shell alias o awtomatiko gamit ang ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2fbee9f4f91
---

Bigyan ang bawat account ng sariling configuration folder. Binabasa ng Claude Code ang folder mula sa environment variable na `CLAUDE_CONFIG_DIR`, at may sariling login, settings at kasaysayan ang bawat folder. Puwede mo itong i-set up nang mano-mano gamit ang shell alias, o hayaang gawin ito ng ccam.

## Mano-mano, gamit ang alias

Ito ang paraang inilalarawan sa dokumentasyon ng Claude Code. Magdagdag ng isang linya bawat karagdagang account sa `~/.bashrc` o `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Magbukas ng bagong terminal at patakbuhin ang `claude-work`. Hihilingin ng Claude Code na mag-log in ka, dahil bago ang folder. Mula noon, `claude` ang una mong account at `claude-work` ang pangalawa.

Gumagana ito nang maayos para sa dalawang account. Mga dapat malaman:

- Walang laman ang bawat folder sa simula: walang settings, walang `CLAUDE.md`, walang MCP server, walang agent o skill. Ikaw ang kokopya o magli-link ng kailangan mo, at ikaw rin ang magpapanatiling magkakatugma ng mga kopya.
- Umiiral lang ang mga alias sa mga interactive shell. Hindi sila nakikita ng mga script, Makefile at iba pang program.
- Nasa `.claude.json` ang mga MCP server, katabi ng login, kaya hindi sila basta mali-link sa pagitan ng mga folder.
- Nananatili ang mga usapan sa account kung saan sila nagsimula, kaya hindi mo maitutuloy ang isa sa ibang account.
- Walang anumang nagpapakita ng usage ng lahat ng account nang sabay-sabay.

## Gamit ang ccam

Ina-automate ng ccam ang parehong ideya at idinadagdag ang mga kulang na bahagi:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Ngayon, sinisimulan ng `claude1` at `claude2`, o ng `ccam personal` at `ccam work`, ang dalawang account. Dahil totoong mga program sila at hindi mga alias, gumagana rin sila sa mga script. Shared na mula sa simula ang settings, `CLAUDE.md`, mga MCP server at extension, at puwede mong i-on ang shared na kasaysayan ng usapan para maituloy ang isang usapan sa kabilang account gamit ang `ccam next --continue`.

Kung gumagamit ka na ng `~/.claude`, patakbuhin ang `ccam import` sa halip na idagdag ulit ang account na iyon.

## Alin ang pipiliin

Para sa dalawang account na walang shine-share, sapat na ang alias. Kapag gusto mo ng shared na settings at MCP server, mas maraming account, mabilisang tingin sa usage o pagtutuloy ng mga usapan, inaalis ng ccam ang abala ng pag-aasikaso. Inililista ng [paghahambing](/compare/aliases) ang mga pagkakaiba nang magkatabi.
