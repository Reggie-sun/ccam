---
title: Farawa da sauri
metaTitle: "Farawa da sauri: ƙara da sauya asusun Claude Code"
description: Ƙara asusun Claude Code biyu da ccam add, ba su gajerun umarni kamar ccwork, sauya tsakaninsu, kuma zaɓi abin da za su raba.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: 0570ef27e4ce
---

Ƙara kowane asusu da `ccam add`, sannan ka fara Claude Code da `ccam <name>` ko da gajeren umarni da kai da kanka ka sa masa suna, kamar `ccwork`. Wannan ne kawai abin da yawancin mutane suke buƙata.

## 1. Ƙara asusunka na farko

```sh
ccam add
```

Mataimakin saitawa (wizard) yana tambayar suna kamar `personal` ko `work`, abubuwan da asusun zai raba da sauran asusunka, da kuma ko a ƙirƙiri gajeren umarni (duba mataki na 4). Sannan yana gudanar da `claude auth login`, wanda ke buɗe burauza don ka shiga. A na'ura mai nisa, burauza na iya ƙarewa a shafin `http://localhost:…/callback` wanda ba ya buɗewa; liƙa wannan adireshin a terminal kuma ccam zai kammala shigarwar. Idan ya gama, za ka ga:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Kawo saitunanka da ke akwai

Idan ka taɓa amfani da Claude Code a baya, saitunanka, `CLAUDE.md`, extensions da sabobin MCP suna cikin `~/.claude`. Kwafe su zuwa saitunan da aka raba, don kowane asusu ya same su:

```sh
ccam migrate
```

ccam yana nuna abin da zai kwafa kuma yana tambaya tukuna. Yana karanta `~/.claude` ne kawai, ba ya taɓa kwafar bayanan shiga, kuma ba ya haɗa tarihin tattaunawa sai ka ƙara `--history`. Duba [kwafar saitin da ke akwai](/docs/sharing#copy-an-existing-setup) don cikakken bayani.

## 3. Ƙara asusu na biyu

```sh
ccam add work
```

Shiga da ɗayan asusun a burauza.

## 4. Sauya asusu

`ccam work` yana fara Claude Code a matsayin asusun work. Don wani abu mafi gajarta, ba asusun nasa gajeren umarni, kamar `ccwork`. Bayan abubuwan da za a raba, mataimakin saitawa yana tambaya:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Danna `y`, sannan ka rubuta suna, ko ka danna Enter don amfani da wanda aka ba da shawara:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Bayan shiga, ccam yana tabbatar da shi:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Idan ka tsallake tambayar, za ka iya ƙara gajeren umarni a kowane lokaci:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Yanzu kowane ɗaya daga cikin waɗannan yana fara Claude Code a matsayin asusun work:

```sh
ccwork
ccam work
ccam 2
```

Gudanar da `ccam` shi kaɗai yana buɗe [dashboard](/features/dashboard): je kan asusu da maɓallan kibiya ko lambarsa, sannan danna Enter.

Duk abin da ke bayan gajeren umarnin ko sunan asusun yana zuwa Claude Code ba tare da canji ba, don haka za ka ci gaba da amfani da zaɓuɓɓukan da ka sani:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Zaɓi abin da za a raba

Ta tsohuwa, ana raba sabobin MCP, `settings.json`, `CLAUDE.md` na gaba ɗaya (global), da agents, commands, skills da plugins tsakanin asusai ta rukunin `default`. Ba a raba tarihin tattaunawa ba. Ƙara shi ga rukunin idan kana son ci gaba da tattaunawa tsakanin asusai:

```sh
ccam group set default --add history
```

ccam yana nuna abin da zai motsa ga kowane asusu kafin ya canza komai. Duba [rabawa da rukunoni](/features/share) don cikakken bayani.

## 6. Idan asusu ya kai iyakarsa

```sh
ccam usage
ccam next --continue
```

`ccam usage` yana nuna amfanin awa 5 da na mako na kowane asusu. `ccam next --continue` yana buɗe wannan tattaunawar a asusu na gaba. A cikin rubutun umarni (scripts), `ccam --fallback work -p "…"` yana sauyawa da kansa. Duba [iyakokin amfani](/docs/usage-limits).

**Kujeru (seats) na Team da Enterprise.** Zaɓin atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard, da [`auto_switch`](/features/next#auto-switch) idan ka kunna shi) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ko za ka yi amfani da shi a can ya rage naka: ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba [bayanin sharuɗɗa a cikin Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Kana amfani da ~/.claude tuni?

Yi masa rajista maimakon sake shiga:

```sh
ccam import
```

Foldar tana nan a inda take, kuma `claude` kai tsaye yana ci gaba da aiki kamar da.
