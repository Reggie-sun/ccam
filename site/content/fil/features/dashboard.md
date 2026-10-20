---
title: Ang dashboard
description: Sa terminal, nagbubukas ang ccam ng full-screen dashboard ng mga account, group, folder rule, extension at pagsusuri; may plano muna ang bawat pagbabago.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

Binubuksan ng `ccam` na walang argument ang isang full-screen dashboard kapag terminal ang standard input at output; kung hindi, ipi-print nito ang `ccam list`. Ipinapakita ng dashboard ang bawat account kasama ang usage at mga group nito, at ginagawa nito ang ginagawa ng mga command: magsimula, magdagdag, mag-log in, magbago ng mga group, magtakda ng mga folder rule, magsuri at mag-ayos. Nagsasara ito kapag nagsimula ang isang account; kapag nagdagdag ng account o nag-log in, iniiwan nito ang screen para sa terminal at bumabalik pagkatapos.

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

Iginuguhit ang unang frame mula sa `config.toml` lang, sa mas mababa pa nang husto sa 100 ms; napupunan ang usage habang sumasagot ang bawat account. Nagiging dilaw ang porsiyento mula 80% at pula sa 100%, at malabong iginuguhit ang mga value na galing sa cache pagkatapos ng nabigong lookup.

## Mga tab

| Tab | Ano ang ipinapakita | Mga key |
| --- | --- | --- |
| Accounts | numero, pangalan, email, 5-oras na bar, 7-araw na usage, huling paggamit, mga group; `✔` ang default, `▸` ang account ng folder na ito, `•` tumatakbo | sa ibaba |
| Groups | bawat group, kung ano ang shine-share nito at ang mga account nito, at ang mga account na walang group | `enter` baguhin, `n` bago, `e` palitan ang pangalan, `x` alisin |
| Folders | ang mga patakaran ng [`ccam folder`](/features/folder) at ang nalalapat dito | `n` gumamit ng account dito, `enter` baguhin, `x` alisin |
| Extensions | ang ipinapakita ng [`ccam ext`](/features/ext), bawat extension nang isang beses kasama ang mga group na pinanggagalingan nito | `r` i-reload |
| Health | ang iniuulat ng [`ccam doctor`](/features/doctor) | `f` ayusin (`doctor --fix`) |

Nagpapalit ng tab ang `←`/`→` o `Tab`, inililista ng `?` ang mga key ng kasalukuyang tab, hinahanap ulit ng `r` ang usage nang hindi ginagamit ang maikling cache, at lumalabas ang `q` o `Esc`.

## Mga key sa tab na Accounts

| Key | Ano ang ginagawa |
| --- | --- |
| `Enter` | Simulan ang Claude Code bilang ang napiling account |
| `1`–`9` | Pumunta sa account na may numerong iyon; i-type nang mabilis ang `1` `2` para sa 12. `Enter` lang ang nagsisimula ng account |
| `↑`/`↓`, `k`/`j` | Lumipat |
| `n` | Simulan ang account na pipiliin ng [`ccam next --auto`](/features/next) |
| `o` | Simulan nang may mga option: `--continue`, `--resume`, ibang argument, o dry run |
| `u` | Mga detalye ng usage: ang 5-oras, 7-araw at model na mga window, kasama ang mga reset nila at kung saan galing ang mga numero |
| `a` | Magdagdag ng account (ang wizard sa ibaba) |
| `l` | Mag-log in muli |
| `s` | Piliin ang mga group ng account |
| `d` | Gawin itong default na account |
| `f` | Gamitin ito sa kasalukuyang folder |
| `e` | Palitan ang pangalan nito |
| `c` | Itakda o alisin ang shortcut command nito |
| `x` | Alisin ito (nagtatanong muna, nakatapat sa No) |

## Pagbabago ng mga group

Sa tab na Groups, binubuksan ng `Enter` ang editor: ang mga item sa kaliwa, ang mga account sa kanan, `Tab` para magpalit, `Space` para mag-toggle. Kulay-abo, kasama ang dahilan, ang item na nakukuha na ng isa sa mga account mula sa ibang group, gayundin ang account na makakakuha ng isang item nang dalawang beses; iisang group lang ang puwedeng pinanggagalingan ng isang item. Ipinapakita ng `Enter` ang plano:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Itinatanong muna ang file na naiiba sa kopya ng group (panatilihin ang sa group, gamitin ang akin, o pagsamahin). Pagkatapos ay tumatakbo ang pagbabago sa background gamit ang parehong code ng [`ccam group`](/features/share): sine-save muna ang config, inililipat nang sunod-sunod ang data ng bawat account, at napupunta ang mga orihinal sa `.ccam-backup/`. Tinatanggihan, kasama ang pangalan ng account, ang pagbabagong maglilipat ng data ng tumatakbong account.

## Pagdagdag ng account

Binubuksan ng `a` ang wizard na ginagamit din ng `ccam add` sa terminal:

1. Ang pangalan ng account.
2. Kung paano ito magsa-sign in: Claude subscription ngayon (sa browser), [API key](/features/api), o Claude subscription mamaya.
3. Para sa API key: kung saan binabasa ng ccam ang key sa bawat pagsisimula (`env:NAME` o `command:COMMAND`), ang endpoint, kung paano ipinapadala ang key at ang model.
4. Ang mga group na sasalihan, na naka-check ang `default`; kapag wala pang group, kung ano ang shine-share ng group na `default`.
5. Isang shortcut command, kung gusto mo.
6. Isang buod na kukumpirmahin.

Bumabalik nang isang hakbang ang `Esc`, at sa unang hakbang, isinasara nito ang wizard.

## Mga kulay

Sinusunod ng mga kulay ang mga papel ng kulay sa terminal UI ng Claude Code, kaya bagay na bagay ang ccam sa tabi nito; nakakakuha ng kulay ang mga pangalan ng group batay sa pangalan nila. Tinatanong ng ccam ang terminal tungkol sa background nito at pumipili ng maliwanag o madilim na theme. Puwede mo itong itakda mismo sa `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Para sa red-green color blindness ang mga theme na `-daltonized`, at ginagamit ng mga theme na `-ansi` ang sariling 16 na kulay ng terminal. Pinapatay ng `NO_COLOR` ang mga kulay.

## Kaligtasan at mga espesyal na kaso

- Tinatawag ng dashboard ang parehong mga function ng mga command, kaya nalalapat ang bawat patakaran: mga plano bago maglipat ng data, pagtanggi para sa mga tumatakbong account, mga backup, at mga login token na hindi kailanman umaalis sa isang account.
- Walang dry run ang dashboard: ipi-print ng `ccam --dry-run` na walang pangalan kung paano magsisimula ang account ng folder na ito, at nag-aalok ang `o` sa isang account ng dry run ng iisang pagsisimulang iyon.
- Naibabalik ang terminal bago magsimula ang Claude Code o bumukas ang login.
- Sa screen na masyadong maliit para sa isang dialog, unang tinatanggal ang mga hilera ng tab; pinuputol ang mga linya gamit ang `…` sa halip na i-wrap.

## Kaugnay

- [Pag-share gamit ang mga group](/features/share)
- [Pagsisimula ng Claude Code](/features/launch)
- [Pamamahala ng mga account](/features/accounts)
- [Mga usage limit](/features/usage)
