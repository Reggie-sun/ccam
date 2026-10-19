---
title: Quick start
description: Magdagdag ng dalawang Claude Code account gamit ang ccam add, lumipat sa pagitan nila gamit ang claude1, claude2 o ccam <name>, at piliin kung ano ang ishe-share nila.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: c6b1c1c3ebaf
---

Idagdag ang bawat account gamit ang `ccam add`, saka simulan ang Claude Code gamit ang `claude1`, `claude2` o `ccam <name>`. Iyan lang ang kailangan ng karamihan.

## 1. Idagdag ang una mong account

```sh
ccam add
```

Hihingin ng wizard ang isang pangalan gaya ng `personal` o `work`, at kung aling mga item ang ishe-share ng account sa iba mo pang account. Pagkatapos ay patatakbuhin nito ang `claude auth login`, na magbubukas ng browser para makapag-sign in ka. Kapag tapos na, makikita mo:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

Sa ngayon, sa Korean ipini-print ng ccam ang sarili nitong mga mensahe; pareho ang mga command at option sa bawat wika.

## 2. Idagdag ang pangalawang account

```sh
ccam add work
```

Mag-sign in gamit ang isa pang account sa browser. Ito ang magiging `claude2`.

## 3. Lumipat

Alinman sa mga ito ay magsisimula ng Claude Code bilang work account:

```sh
claude2
ccam work
ccam 2
```

Ipinapasa nang walang pagbabago sa Claude Code ang lahat ng kasunod ng pangalan ng account, kaya patuloy mong magagamit ang mga option na alam mo na:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Kapag `ccam` lang ang pinatakbo mo, magbubukas ang isang picker: gamitin ang mga arrow key at Enter, o pindutin ang numero ng account.

## 4. Piliin kung ano ang ishe-share

Bilang default, shared sa pagitan ng mga account ang mga MCP server, `settings.json`, ang global na `CLAUDE.md`, at ang mga agent, command, skill at plugin. Hindi kasama ang kasaysayan ng usapan. I-on ito para sa parehong account kung gusto mong ituloy ang mga usapan sa pagitan nila:

```sh
ccam share personal --history on
ccam share work --history on
```

Ipinapakita ng ccam kung ano ang ililipat nito bago ito magbago ng anuman. Tingnan ang [sharing](/docs/sharing) para sa mga detalye.

## 5. Kapag naabot ng isang account ang limit nito

```sh
ccam usage
ccam next --continue
```

Ipinapakita ng `ccam usage` ang 5-oras at lingguhang usage ng bawat account. Binubuksan ng `ccam next --continue` ang parehong usapan sa susunod na account. Tingnan ang [mga usage limit](/docs/usage-limits).

## Gumagamit ka na ba ng ~/.claude?

Irehistro ito sa halip na mag-sign in ulit:

```sh
ccam import
```

Mananatili ang folder kung nasaan ito, at gagana pa rin gaya ng dati ang karaniwang `claude`.
