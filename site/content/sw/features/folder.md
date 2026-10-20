---
title: Akaunti kwa folda
description: ccam folder hufanya folda na kila kitu chini yake kitumie akaunti moja unapoanzisha Claude Code bila jina la akaunti, kama work kwa msimbo wa kampuni yako.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: "46d1410a3d81"
---

`ccam folder work ~/code/company` hufanya kila uanzishaji usio na jina la akaunti ndani ya `~/code/company`, na katika folda yoyote iliyo chini yake, uanzishe Claude Code kama `work`. Nje yake, akaunti chaguo-msingi huanza kama hapo awali. Kutaja akaunti, kwa `ccam personal` au njia fupi kama `ccpersonal`, huanzisha akaunti hiyo kila mara.

## Matumizi

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Uanzishaji huu hutumia akaunti ya folda:

- `ccam -p "…"` na uanzishaji mwingine wowote ambao hoja yake ya kwanza ni chaguo la Claude Code
- `ccam run` bila jina, na `ccam --dry-run`
- `ccam --fallback -p "…"` kwa jaribio lake la kwanza
- dashibodi, `ccam` peke yake, ambacho huanza kiashiria chake kikiwa kwenye akaunti ya folda

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `<account>` | Jina au namba ya akaunti ya kutumia | |
| `[folder]` | Folda, ya kiasi (relative) au kamili (absolute) | folda ya sasa |
| `--remove` | Huondoa ingizo la folda hiyo hasa | imezimwa |
| `--dry-run` | Inaheshimiwa: huchapisha badiliko bila kulihifadhi | imezimwa |

## Jinsi inavyofanya kazi

Kila folda ni ingizo moja la `[[folder]]` katika `~/.config/ccam/config.toml`, faili ambalo kila uanzishaji husoma hata hivyo:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` hubadilisha folda kuwa njia kamili, viungo vya ishara (symbolic links) vikiwa vimefuatiliwa, na huihifadhi. Kuweka folda tena hubadilisha akaunti yake.
2. Uanzishaji bila jina huita `launcher.HereProfile`. Kukiwa hakuna maingizo, huchukua akaunti chaguo-msingi mara moja.
3. Vinginevyo hulinganisha folda ya kazi na kila `path` kama mifuatano ya maandishi: folda yenyewe au iliyo chini yake, kamwe si folda inayoanza tu kwa herufi zile zile (`/code/company` haihusishi `/code/companyx`). Ingizo lenye kina zaidi hushinda, kwa hiyo `~/code` inaweza kuwa `personal` wakati `~/code/company` ni `work`.
4. Ikiwa tu folda ya kazi kama ilivyo hailingani na chochote, ccam hufuatilia viungo vyake vya ishara na kulinganisha mara moja zaidi.
5. Kwenye Windows ulinganisho hupuuza herufi kubwa na ndogo, kama Windows inavyofanya.

Hakuna faili jingine linalosomwa isipokuwa `config.toml`, kwa hiyo uanzishaji ndani ya folda hauongezi muda unaoweza kupimwa (4.0 ms ikiwa na maingizo au bila, `hyperfine`, mizunguko 300).

## Usalama na hali maalum

- `ccam rename` huhamisha maingizo ya folda kwenda kwenye jina jipya, na `ccam remove` hufuta maingizo ya akaunti iliyoondolewa.
- Ingizo ambalo akaunti yake haipo tena husimamisha uanzishaji bila jina katika folda hiyo: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` huripoti maingizo kama hayo, na maingizo ambayo folda zake zimetoweka.
- `ccam folder --remove` ndani ya folda inayopata akaunti yake kutoka folda iliyo juu yake hutaja folda hiyo badala ya kuondoa chochote.
- `ccam next` huchagua kwa mzunguko au matumizi na haiangalii folda.

## Zinazohusiana

- [Kuanzisha Claude Code](/features/launch)
- [Amri za njia fupi](/features/shortcuts)
- [Kusimamia akaunti](/features/accounts)
- [Marejeo ya amri](/docs/commands)
