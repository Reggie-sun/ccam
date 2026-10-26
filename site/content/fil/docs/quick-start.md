---
title: Mabilis na pagsisimula
metaTitle: "Mabilis na pagsisimula: magdagdag ng mga account"
description: Magdagdag ng dalawang Claude Code account gamit ang ccam add, bigyan sila ng shortcut gaya ng ccwork, lumipat sa pagitan nila, at piliin ang ishe-share nila.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

Idagdag ang bawat account gamit ang `ccam add`, saka simulan ang Claude Code gamit ang `ccam <name>` o ang shortcut na ikaw mismo ang nagpangalan, gaya ng `ccwork`. Iyan lang ang kailangan ng karamihan.

Gumagamit ka na ba ng Claude Code? Basahin ang [Gumagamit ka na ba ng ~/.claude?](#already-using-claude) sa dulo ng pahinang ito bago ka magsimula.

## 1. Idagdag ang una mong account

```sh
ccam add
```

Hihingin ng wizard ang isang pangalan gaya ng `personal` o `work`, kung aling mga item ang ishe-share ng account sa iba mo pang account, at kung gagawa ng shortcut command (tingnan ang hakbang 3). Pagkatapos ay patatakbuhin nito ang `claude auth login`, na magbubukas ng browser para makapag-sign in ka. Sa isang remote machine, maaaring mapunta ang browser mo sa isang `http://localhost:…/callback` na page na hindi naglo-load; i-paste ang address na iyon sa terminal at tatapusin ng ccam ang sign-in. Kapag tapos na ang sign-in, makikita mo:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Dalhin ang dati mong settings

Kung gumamit ka na ng Claude Code dati, nasa `~/.claude` ang iyong settings, `CLAUDE.md` at mga extension, at nasa `~/.claude.json` ang iyong mga MCP server. Kopyahin ang mga ito sa shared na group na `default` para makuha ng bawat account:

```sh
ccam migrate
```

Ipapakita muna ng ccam kung ano ang kokopyahin nito at magtatanong bago magsimula. Binabasa lang nito ang `~/.claude`, hindi kailanman kinokopya ang login, at hindi isinasama ang kasaysayan ng usapan maliban kung idagdag mo ang `--history`. Tingnan ang [pagkopya ng dati nang setup](/docs/sharing#copy-an-existing-setup) para sa mga detalye.

## 3. Idagdag ang pangalawang account

```sh
ccam add work
```

Nilalaktawan ng wizard ang pangalan, na ibinigay mo na sa command line, at ipinapakita ang mga group na puwedeng salihan, na naka-check ang `default`. Sa pagkakataong ito, bigyan ang account ng sarili nitong shortcut command, gaya ng `ccwork`, para masimulan mo ito sa isang salita. Pagkatapos ng mga group, magtatanong ang wizard:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Pindutin ang `y`, saka mag-type ng pangalan, o pindutin ang Enter para gamitin ang iminungkahing pangalan:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Saka mag-sign in gamit ang isa pang account sa browser. Kinukumpirma ng ccam ang account at ang shortcut nito:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Lumipat ng account

Sinisimulan ng `ccam work` ang Claude Code bilang work account. Kung nilaktawan mo ang tanong tungkol sa shortcut, puwede kang magdagdag ng shortcut anumang oras:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Sinisimulan na ngayon ng bawat isa sa mga command na ito ang Claude Code bilang work account (`2` ang numero nito sa `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Kapag `ccam` lang ang pinatakbo mo, magbubukas ang [dashboard](/features/dashboard): pumunta sa isang account gamit ang mga arrow key o ang numero nito, saka pindutin ang Enter.

Ipinapasa nang walang pagbabago sa Claude Code ang lahat ng kasunod ng shortcut o ng pangalan ng account, kaya patuloy mong magagamit ang mga option na alam mo na:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Piliin kung ano ang ishe-share

Bilang default, shared sa pagitan ng mga account ang mga MCP server, `settings.json`, ang global na `CLAUDE.md` at ang mga extension (agent, command, skill at plugin), sa pamamagitan ng group na `default`. Hindi kasama ang kasaysayan ng usapan. Idagdag ito sa group kung gusto mong ituloy ang mga usapan sa pagitan ng mga account:

```sh
ccam group set default --add history
```

Ipinapakita ng ccam kung ano ang ililipat nito para sa bawat account bago ito magbago ng anuman. Tingnan ang [pag-share gamit ang mga group](/features/share) para sa mga detalye.

## 6. Kapag naabot ng isang account ang limit nito

```sh
ccam usage
ccam next --continue
```

Ipinapakita ng `ccam usage` ang 5-oras at lingguhang usage ng bawat account. Binubuksan ng `ccam next --continue` ang parehong usapan sa susunod na account. Sa mga script, kusang lumilipat ang `ccam --fallback work -p "…"`. Tingnan ang [mga usage limit](/docs/usage-limits).

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili ng account (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard at, kapag binuksan mo ito, ang [`auto_switch`](/features/next#auto-switch)) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon, na maaaring lumabag sa terms ng Anthropic. Basahin ang [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) at magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito.

## Gumagamit ka na ba ng ~/.claude?

Irehistro ito sa halip na mag-sign in ulit:

```sh
ccam import
```

Mananatili ang folder kung nasaan ito, at gagana pa rin gaya ng dati ang karaniwang `claude`.
