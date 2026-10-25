---
title: Amfani da asusun Claude Code da dama a kwamfuta ɗaya
description: Gudanar da asusun Claude Code da dama gefe da gefe da CLAUDE_CONFIG_DIR, da hannu ta shell alias ko ta atomatik da ccam, kuma ka sauya tsakaninsu.
updated: 2026-10-25
order: 1
translation: machine
sourceHash: 5c5dbd9fcd23
---

Ba kowane asusu foldar saitunsa ta kansa. Claude Code yana karanta foldar daga environment variable ɗin `CLAUDE_CONFIG_DIR`, kuma kowace folda tana riƙe bayanan shigarta, saitunanta da tarihinta, don haka kowane asusu yana nan a shiga kuma za ka iya gudanar da su a lokaci guda. Za ka iya saita wannan da hannu da shell alias, ko ka bar ccam ya yi maka.

## Yaya zan saita shi da hannu?

Wannan ita ce hanyar da [takardun Claude Code suka bayyana](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Ƙara layi ɗaya ga kowane ƙarin asusu a `~/.bashrc` ko `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Buɗe sabon terminal kuma ka gudanar da `claude-work`. Claude Code zai ce ka shiga, domin foldar sabuwa ce. Daga nan, `claude` shi ne asusunka na farko, `claude-work` kuma na biyu. Ƙara alias ɗaya ga kowane asusu, kowanne da foldarsa.

A Windows ba tare da WSL ba, function na PowerShell a cikin profile ɗinka (`notepad $PROFILE`) yana yin haka:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Me hanyar alias ba ta yi ba?

Tana aiki da kyau ga asusai biyu. Abubuwan da ya kamata ka sani:

- Kowace folda tana farawa babu komai: babu saituna, babu `CLAUDE.md`, babu sabobin MCP, babu agents ko skills. Kai ne za ka kwafa ko haɗa abin da kake buƙata, kuma ka riƙa daidaita kwafin da kanka.
- Alias yana samuwa ne kawai a shell na mu'amala (interactive shell). Scripts, Makefiles da sauran shirye-shirye ba sa ganinsu.
- Sabobin MCP suna cikin `.claude.json`, kusa da bayanan shiga, don haka ba za a iya kawai haɗa su tsakanin foldoji da mahaɗi ba.
- Tattaunawa suna zama a asusun da aka fara su, don haka ba za ka iya ci gaba da ɗaya a wani asusu ba.
- Babu wani abu da ke nuna amfanin duk asusai a lokaci ɗaya.

## Zan iya sauya asusu da /login maimakon haka?

Ee, idan kana sauyawa lokaci-lokaci ne kawai. Rubuta `/login` a cikin Claude Code [yana sauya asusun](https://code.claude.com/docs/en/quickstart) foldar da kake ciki. Asusu ɗaya ne kawai yake a shiga a lokaci guda, don haka sai ka sake shiga a kowace sauyawa, kuma ba za ka iya gudanar da asusai biyu gefe da gefe ba. Foldoji daban-daban suna kauce wa duka biyun.

## Yaya zan saita shi da ccam?

ccam yana sarrafa wannan ra'ayin ta atomatik kuma yana ƙara ɓangarorin da suka ɓace:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Yanzu `ccam personal` da `ccam work` suna fara asusun biyu, haka ma `ccwork`, gajeren umarnin da ka ba asusun work. Saboda shirye-shirye ne na gaske ba alias ba, suna aiki a scripts ma. Gudanar da `ccam` shi kaɗai don zaɓar asusu daga jeri, da `ccam use work` don mayar da ɗaya tsoho (default).

Ana raba saituna, `CLAUDE.md`, sabobin MCP da extensions tun daga farko. `ccam next --continue` yana kai tattaunawa zuwa ɗayan asusun, kuma raba tarihin tattaunawa, wanda za ka iya kunnawa, yana riƙe ta a matsayin tattaunawa ɗaya. `ccam usage` yana nuna amfanin awa 5 da na mako na kowane asusu a lokaci ɗaya.

Idan kana amfani da `~/.claude` tuni, gudanar da `ccam import` maimakon sake ƙara wannan asusun, ko `ccam migrate` don kwafar saitunansa, sabobin MCP da extensions zuwa saitunan da aka raba.

## Ta yaya zan san wane asusu ke aiki?

A cikin Claude Code, `/status` yana nuna asusun da ka shiga da shi. Da ccam, gudanar da `ccam statusline install` sau ɗaya kuma layin matsayi (status line) zai riƙa nuna sunan asusun da amfaninsa a kowane lokaci, misali `[work] · Opus · 5h 24% · 7d 41%`.

## Shin wannan yana aiki a Windows da WSL?

Claude Code yana aiki [kai tsaye a Windows da kuma cikin WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` yana aiki a duka biyun, don haka alias ko function na PowerShell da ke sama ya isa a Windows ba tare da WSL ba. ccam ma yana aiki kai tsaye a Windows: girka shi a PowerShell da `irm https://ccam.dev/install.ps1 | iex`, ko a cikin WSL da umarnin Linux. Gajerun umarni suna zama `ccwork.exe`, kuma raba fayiloli kamar `settings.json` yana buƙatar Developer Mode ([Girka ccam](/docs/install)).

## Magance matsaloli

**Claude Code yana tambaya na shiga kowane lokaci.** Farkon farawa a sabuwar folda koyaushe yana tambaya. Idan ya ci gaba da tambaya, shigar ta ƙare ko an cire ta: gudanar da `ccam login work` (ko ka fara alias ɗin sau ɗaya ka shiga).

**Sabar MCP tana bayyana a asusu ɗaya kawai.** Sabobin da aka ƙara a cikin zama da `claude mcp add` na wannan asusun ne, kuma ba tare da `-s user` ba na wannan foldar kawai. Matsar da ɗaya zuwa duk asusai da `ccam mcp import work <server>`.

**Saituna sun bambanta tsakanin asusai.** `ccam share work` yana nuna abin da asusun ke rabawa da kuma ta wane rukuni; `ccam group join default work` yana sa shi ya raba abin da `default` ke rabawa, bayan ya nuna abin da zai canza. `ccam doctor --fix` yana gyara mahaɗan da suka lalace kuma yana ajiye kwafin kowane fayil da ya canza.

**Asusai biyu suna buɗe tattaunawa ɗaya.** Ka guji yin aiki a kan tattaunawa ɗaya a asusai biyu a lokaci guda. ccam yana gargaɗe ka idan ka gwada.

## Wanne zan zaɓa?

Ga asusai biyu ba tare da raba komai ba, alias ya isa. Idan kana son saituna da sabobin MCP da aka raba, ƙarin asusai, ganin amfani a kallo ɗaya ko ci gaba da tattaunawa, ccam yana ɗauke maka aikin bin diddigi. [Kwatance da alias](/compare/aliases) yana jera bambance-bambancen gefe da gefe, kuma [ccam da sauran masu sauya asusu](/compare/account-switchers) yana bayani kan kayan aikin da ke musanya bayanan shiga ɗaya maimakon haka.

Ko an yarda ka yi amfani da asusai da dama ya dogara da sharuɗɗan kowane asusu. Duba su kafin ka yi haka.

**Kujeru (seats) na Team da Enterprise.** Zaɓin atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard, da [`auto_switch`](/features/next#auto-switch) idan ka kunna shi) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ko za ka yi amfani da shi a can ya rage naka: ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba bayanin da ke sama.
