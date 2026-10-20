---
title: Dashibodi
description: ccam kwenye terminal hufungua dashibodi ya skrini nzima ya akaunti, vikundi, kanuni za folda, viendelezi na ukaguzi; kila badiliko huonyesha mpango wake kwanza.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` bila hoja hufungua dashibodi ya skrini nzima wakati ingizo na matokeo ya kawaida ni terminal; vinginevyo huchapisha `ccam list`. Dashibodi huonyesha kila akaunti pamoja na matumizi yake na vikundi vyake, na hufanya kile ambacho amri hufanya: kuanzisha, kuongeza, kuingia, kubadilisha vikundi, kuweka kanuni za folda, kukagua na kurekebisha. Kuanzisha huifunga; kuongeza akaunti au kuingia huiacha kwa ajili ya terminal kisha hurudi.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Fremu ya kwanza huchorwa kutoka `config.toml` pekee, kwa muda mfupi sana chini ya 100 ms; matumizi hujazwa kadiri kila akaunti inavyojibu. Asilimia hugeuka manjano kuanzia 80% na nyekundu kwenye 100%, na thamani kutoka kwenye akiba (cache) baada ya utafutaji kushindwa huchorwa hafifu.

## Vichupo

| Kichupo | Kinachoonyesha | Vitufe |
| --- | --- | --- |
| Accounts | namba, jina, barua pepe, upau wa saa 5, matumizi ya siku 7, matumizi ya mwisho, vikundi; `✔` akaunti chaguo-msingi, `▸` akaunti ya folda hii, `•` inaendeshwa | hapa chini |
| Groups | kila kikundi, kile kinachoshiriki na akaunti zake, pamoja na akaunti zisizo katika kikundi chochote | `enter` hariri, `n` kipya, `e` badilisha jina, `x` ondoa |
| Folders | kanuni za [`ccam folder`](/features/folder) na ile inayotumika hapa | `n` tumia akaunti hapa, `enter` badilisha, `x` ondoa |
| Extensions | kile ambacho [`ccam ext`](/features/ext) huonyesha, kila kiendelezi mara moja pamoja na vikundi kinakotoka | `r` pakia upya |
| Health | kile ambacho [`ccam doctor`](/features/doctor) huripoti | `f` rekebisha (`doctor --fix`) |

`←`/`→` au `Tab` hubadilisha vichupo, `?` huorodhesha vitufe vya kichupo cha sasa, `r` hutafuta matumizi tena bila akiba fupi, na `q` au `Esc` hufunga.

## Vitufe kwenye kichupo cha Accounts

| Kitufe | Inachofanya |
| --- | --- |
| `Enter` | Anzisha Claude Code kama akaunti iliyochaguliwa |
| `1`–`9` | Nenda kwenye akaunti yenye namba hiyo; andika `1` `2` haraka kwa 12. `Enter` pekee ndicho huanzisha |
| `↑`/`↓`, `k`/`j` | Sogea |
| `n` | Anzisha akaunti ambayo [`ccam next --auto`](/features/next) ingechagua |
| `o` | Anzisha kwa chaguo: `--continue`, `--resume`, hoja nyingine, au jaribio (dry run) |
| `u` | Maelezo ya matumizi: madirisha ya saa 5, siku 7 na ya model pamoja na nyakati zake za kuwekwa upya na mahali namba zinapotoka |
| `a` | Ongeza akaunti (mchawi ulio hapa chini) |
| `l` | Ingia tena |
| `s` | Chagua vikundi vya akaunti |
| `d` | Ifanye akaunti chaguo-msingi |
| `f` | Itumie katika folda ya sasa |
| `e` | Badilisha jina lake |
| `c` | Weka au ondoa amri yake ya njia fupi |
| `x` | Iondoe (huuliza kwanza, ikianzia kwenye No) |

## Kubadilisha vikundi

Kwenye kichupo cha Groups, `Enter` hufungua kihariri: vipengele upande wa kushoto, akaunti upande wa kulia, `Tab` kubadili, `Space` kuwasha au kuzima. Kipengele ambacho mojawapo ya akaunti tayari inakipata kutoka kikundi kingine hufifishwa pamoja na sababu, vivyo hivyo akaunti ambayo ingepata kipengele mara mbili; kipengele kinaweza kutoka kwenye kikundi kimoja tu. `Enter` huonyesha mpango:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Faili inayotofautiana na nakala ya kikundi huulizwa kwanza (kuhifadhi ya kikundi, kutumia yangu, au kuunganisha). Kisha badiliko huendeshwa chinichini kwa msimbo ule ule wa [`ccam group`](/features/share): usanidi huhifadhiwa kwanza, data ya kila akaunti huhamishwa kwa mpangilio, na faili asili huenda `.ccam-backup/`. Badiliko ambalo lingehamisha data ya akaunti inayoendeshwa hukataliwa pamoja na jina la akaunti hiyo.

## Kuongeza akaunti

`a` hufungua mchawi ambao `ccam add` pia hutumia kwenye terminal:

1. Jina la akaunti.
2. Jinsi inavyoingia: usajili wa Claude sasa (kwenye kivinjari), [API key](/features/api), au usajili wa Claude baadaye.
3. Kwa API key: mahali ccam husoma ufunguo katika kila uanzishaji (`env:NAME` au `command:COMMAND`), endpoint, jinsi ufunguo unavyotumwa na model.
4. Vikundi vya kujiunga, `default` ikiwa imeteuliwa; kabla ya kuwepo kwa vikundi, kile ambacho kikundi `default` hushiriki.
5. Amri ya njia fupi, ukiitaka.
6. Muhtasari wa kuthibitisha.

`Esc` hurudi hatua moja nyuma, na kwenye hatua ya kwanza hufunga mchawi.

## Rangi

Rangi hufuata majukumu ya kiolesura cha terminal cha Claude Code, ili ccam ionekane ikiendana nayo; majina ya vikundi hupata rangi kutokana na jina lake. ccam huuliza terminal kuhusu mandharinyuma yake na huchagua mandhari angavu au meusi. Weka moja mwenyewe katika `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Mandhari ya `-daltonized` ni kwa ajili ya upofu wa rangi nyekundu na kijani, na mandhari ya `-ansi` hutumia rangi 16 za terminal yenyewe. `NO_COLOR` huzima rangi.

## Usalama na hali maalum

- Dashibodi huita kazi zile zile ambazo amri huita, kwa hiyo kila kanuni hubaki: mipango kabla data haijahamishwa, kukataa kwa akaunti zinazoendeshwa, nakala rudufu, na tokeni za kuingia ambazo kamwe haziondoki kwenye akaunti.
- Dashibodi haina jaribio (dry run): `ccam --dry-run` bila jina huchapisha jinsi akaunti ya folda hii ingeanza, na `o` kwenye akaunti hutoa jaribio la uanzishaji huo mmoja tu.
- Terminal hurejeshwa kabla Claude Code haijaanza au kuingia hakujafunguka.
- Skrini iliyo ndogo mno kwa kidirisha huondoa kwanza safu za kichupo; mistari hukatwa kwa `…` badala ya kufungwa kwenye mstari unaofuata.

## Zinazohusiana

- [Kushiriki kwa vikundi](/features/share)
- [Kuanzisha Claude Code](/features/launch)
- [Kusimamia akaunti](/features/accounts)
- [Vikomo vya matumizi](/features/usage)
