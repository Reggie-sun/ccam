---
title: Yadda ake amfani da asusun Claude Code biyu a kwamfuta ɗaya
description: Gudanar da asusun Claude Code da dama gefe da gefe da CLAUDE_CONFIG_DIR, da hannu ta shell alias ko ta atomatik da ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d82129ddcb4c
---

Ba kowane asusu foldar saitunsa ta kansa. Claude Code yana karanta foldar daga environment variable ɗin `CLAUDE_CONFIG_DIR`, kuma kowace folda tana riƙe bayanan shigarta, saitunanta da tarihinta. Za ka iya saita wannan da hannu da shell alias, ko ka bar ccam ya yi maka.

## Da hannu, da alias

Wannan ita ce hanyar da takardun bayanin Claude Code suka bayyana. Ƙara layi ɗaya ga kowane ƙarin asusu a `~/.bashrc` ko `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Buɗe sabon terminal kuma ka gudanar da `claude-work`. Claude Code zai ce ka shiga, domin foldar sabuwa ce. Daga nan, `claude` shi ne asusunka na farko, `claude-work` kuma na biyu.

Wannan yana aiki da kyau ga asusai biyu. Abubuwan da ya kamata ka sani:

- Kowace folda tana farawa babu komai: babu saituna, babu `CLAUDE.md`, babu sabobin MCP, babu agents ko skills. Kai ne za ka kwafa ko haɗa abin da kake buƙata, kuma ka riƙa daidaita kwafin da kanka.
- Alias yana samuwa ne kawai a shell na mu'amala (interactive shell). Scripts, Makefiles da sauran shirye-shirye ba sa ganinsu.
- Sabobin MCP suna cikin `.claude.json`, kusa da bayanan shiga, don haka ba za a iya kawai haɗa su tsakanin foldoji da mahaɗi ba.
- Tattaunawa suna zama a asusun da aka fara su, don haka ba za ka iya ci gaba da ɗaya a wani asusu ba.
- Babu wani abu da ke nuna amfanin duk asusai a lokaci ɗaya.

## Da ccam

ccam yana sarrafa wannan ra'ayin ta atomatik kuma yana ƙara ɓangarorin da suka ɓace:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Yanzu `claude1` da `claude2`, ko `ccs personal` da `ccs work`, suna fara asusun biyu. Saboda shirye-shirye ne na gaske ba alias ba, suna aiki a scripts ma. Ana raba saituna, `CLAUDE.md`, sabobin MCP da extensions tun daga farko, kuma za ka iya kunna rabon tarihin tattaunawa don ci gaba da tattaunawa a ɗayan asusun da `ccs next --continue`.

Idan kana amfani da `~/.claude` tuni, gudanar da `ccam import` maimakon sake ƙara wannan asusun.

## Wanne za a zaɓa

Ga asusai biyu ba tare da raba komai ba, alias ya isa. Idan kana son saituna da sabobin MCP da aka raba, ƙarin asusai, ganin amfani a kallo ɗaya ko ci gaba da tattaunawa, ccam yana ɗauke maka aikin bin diddigi. [Kwatancen](/compare/aliases) yana jera bambance-bambancen gefe da gefe.
