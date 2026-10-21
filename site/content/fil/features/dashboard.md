---
title: Ang dashboard
description: Sa terminal, nagbubukas ang ccam ng full-screen dashboard ng mga account, group, folder rule, extension at pagsusuri; may plano muna ang bawat pagbabago.
updated: 2026-10-21
order: 15
translation: machine
sourceHash: "19a99adef661"
---

Binubuksan ng `ccam` na walang argument ang isang full-screen dashboard kapag terminal ang standard input at output; kung hindi, ipi-print nito ang `ccam list`. Ipinapakita ng dashboard ang bawat account kasama ang usage at mga group nito, at ginagawa nito ang ginagawa ng mga command: magsimula, magdagdag, mag-log in, magbago ng mga group, magtakda ng mga folder rule, magsuri at mag-ayos. Nagsasara ito kapag nagsimula ang isang account; kapag nagdagdag ng account o nag-log in, iniiwan nito ang screen para sa terminal at bumabalik pagkatapos.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Sinasabi ng linya sa ilalim ng mga pangalan ng tab, sa isang pangungusap, kung para saan ang kasalukuyang tab. Sa mas malapad na terminal, ipinapakita rin ng tab na Accounts ang email ng bawat account.

Tuwing nagbubukas ito, tinatanong ng dashboard sa GitHub sa background kung may mas bagong ccam release; hindi kailanman naghihintay rito ang unang frame, at walang ipinapakita kapag offline. Kapag mayroon, inirerekomenda ito ng isang linya sa ilalim ng pamagat:

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

Nagtatanong muna ang `U`, saka nag-a-update gaya ng ginagawa ng [`ccam update`](/features/update): sinusuri ang release laban sa mga SHA-256 checksum nito at pinapalitan nito ang ccam sa isang hakbang. Lumang bersyon pa rin ang dashboard na kinaroroonan mo, kaya lumabas at simulan muli ang `ccam`. Kapag hindi magagawa rito ang update, halimbawa para sa ccam na in-install mula sa `.deb` package, ipinapakita ng dashboard ang parehong payo gaya ng `ccam update` at pinananatili ang alok.

Iginuguhit ang unang frame mula sa `config.toml` lang, sa mas mababa pa nang husto sa 100 ms; napupunan ang usage habang sumasagot ang bawat account. May sariling bar ang bawat limit: ang 5-oras na window, ang lingguhang window para sa lahat ng model, at ang sariling lingguhang window ng Fable; ang ibig sabihin ng `-` ay walang ganoong window ang account. Nagiging dilaw ang porsiyento mula 80% at pula sa 100%, at malabong iginuguhit ang mga value na galing sa cache pagkatapos ng nabigong lookup.

## Mga tab

| Tab | Ano ang ipinapakita | Mga key |
| --- | --- | --- |
| Accounts | numero, pangalan, email, mga bar ng 5-oras, lingguhan at Fable, huling paggamit, mga group; `✔` ang default, `▸` ang account ng folder na ito, `•` tumatakbo | sa ibaba |
| Groups | bawat group, kung ano ang shine-share nito at ang mga account nito, at ang mga account na walang group | `enter` baguhin, `n` bago, `e` palitan ang pangalan, `x` alisin |
| Folders | ang mga patakaran ng [`ccam folder`](/features/folder) at ang nalalapat dito | `n` gumamit ng account dito, `enter` baguhin, `x` alisin |
| Extensions | ang ipinapakita ng [`ccam ext`](/features/ext), bawat extension nang isang beses kasama ang mga group na pinanggagalingan nito | `r` i-reload |
| Health | ang iniuulat ng [`ccam doctor`](/features/doctor), bilang isang kahon bawat account (sa ibaba) | `f` ayusin ang napiling kahon, `F` ayusin lahat (`doctor --fix`), `↑`/`↓` lumipat sa pagitan ng mga kahon |

Nagpapalit ng tab ang `←`/`→` o `Tab`, inililista ng `?` ang mga key ng kasalukuyang tab, hinahanap ulit ng `r` ang usage nang hindi ginagamit ang maikling cache, at lumalabas ang `q` o `Esc`.

## Mga key sa tab na Accounts

| Key | Ano ang ginagawa |
| --- | --- |
| `Enter` | Simulan ang Claude Code bilang ang napiling account |
| `1`–`9` | Pumunta sa account na may numerong iyon; i-type nang mabilis ang `1` `2` para sa 12. `Enter` lang ang nagsisimula ng account |
| `↑`/`↓`, `k`/`j` | Lumipat |
| `n` | Simulan ang account na pipiliin ng [`ccam next --auto`](/features/next) |
| `o` | Simulan nang may mga option: `--continue`, `--resume`, ibang argument, o dry run |
| `u` | Mga detalye ng usage: ang 5-oras, lingguhan at model na mga window, kasama ang mga reset nila at kung saan galing ang mga numero |
| `a` | Magdagdag ng account (ang wizard sa ibaba) |
| `l` | Mag-log in muli |
| `s` | Piliin ang mga group ng account |
| `d` | Gawin itong default na account |
| `f` | Gamitin ito sa kasalukuyang folder |
| `e` | Palitan ang pangalan nito |
| `c` | Itakda o alisin ang shortcut command nito |
| `x` | Alisin ito (nagtatanong muna, nakatapat sa No) |

## Ang tab na Health

Pinapatakbo ng tab na Health ang mga pagsusuri ng [`ccam doctor`](/features/doctor) kapag binuksan mo ito at iginuguhit ang parehong view: isang unang linya na may marka bawat account, saka isang kahon bawat account. Ang bawat hilera sa isang kahon ay isang bahaging sinuri, `Folder`, `Login` at `Share links` muna, na may markang `✔` OK, `ℹ` impormasyon, `⚠` babala o `✘` error, kaya makikita mo agad kung aling account ang may problema at saan. Napupunta sa kahong `General` ang mga natuklasang hindi kabilang sa anumang account.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

Lumilipat ang `↑`/`↓` mula kahon papunta sa kahon, at may border na kulay accent ang napiling kahon. Kinukumpuni lang ng `f` ang napiling kahon, sa paraang gagawin ng `ccam doctor --fix`: sa kahon ng isang account, ang account na iyon lang ang kinukumpuni nito, at sa kahong `General`, ang mga natuklasang hindi kabilang sa anumang account (mga path ng plugin at mga shortcut); sinusuri lang ulit ang bawat iba pang account. Nagtatanong muna ito, at ang kahong walang kailangang kumpunihin ay nagsasabi lang ng `nothing to repair in <name>`. Kinukumpuni ng `F` ang lahat nang sabay-sabay, gaya ng `ccam doctor --fix`, at sinusuri ulit ng `r`.

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
