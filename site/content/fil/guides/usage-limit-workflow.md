---
title: Kapag naabot mo ang Claude Code usage limit
description: Maghintay ng reset o ituloy agad ang parehong usapan sa ibang Claude Code account. Paano gumagana ang 5-oras at lingguhang limit, hakbang-hakbang.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 38b80d780bba
---

Kapag naabot ng isang account ang usage limit nito (madalas tawaging rate limit), may dalawa kang pagpipilian: maghintay hanggang mag-reset ang limit, o ituloy ngayon din ang parehong usapan sa isa pang account na pag-aari mo. Sa ccam, isang command lang ang pangalawa, `ccam next --auto --continue`, kapag shared na ang kasaysayan ng usapan sa pagitan ng mga account.

## Anong mga limit mayroon ang Claude Code?

Sa mga Pro at Max plan, may [limang-oras na session limit at lingguhang limit](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan) ang Claude Code. Parehong kahati ng mga ito ang Claude sa web at sa mga app, at sabay na ibinabawas ang usage sa dalawa. Kapag naubos ang isa, [nagpapakita ng mensahe](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) ang Claude Code gaya ng "You've hit your session limit · resets 3:45pm". Hindi makakatulong ang paglipat ng model, dahil saklaw ng mga limit ang lahat ng model.

## Opsyon 1: hintayin ang reset

Kung iisa lang ang account mo, maghintay. Ang mga bagong bersyon ng Claude Code ay [naghihintay sa bukas na session](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) at kusang itinutuloy ang gawain kapag nag-reset na ang limit. Ipinapakita ng `/usage` sa loob ng Claude Code kung gaano na karami ang nagamit mo sa mga limit ng iyong plan.

## Opsyon 2: ituloy ngayon sa ibang account

Kung may pangalawa kang account, puwede mong ituloy agad ang parehong usapan. I-set up ito nang isang beses:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Ipinapakita na ngayon ng status line ang account at ang 5-oras na usage nito habang nagtatrabaho ka, kaya hindi ka mabibigla sa limit.

Kapag naabot ang limit, lumabas sa Claude Code (`/exit`), pagkatapos:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Buksan ang usapan sa account na may pinakamaraming natitira:

```sh
ccam next --auto --continue
```

Magsisimula ang Claude Code bilang `work` kasama ang pinakahuling usapan ng kasalukuyang folder, kasama ang lahat ng nasabi na. Kapag shared ang kasaysayan (ang `ccam group set default --add history` sa itaas), nananatili itong iisang usapan; kung hindi, kinokopya ng Claude Code ang usapan sa `work`, at sinasabi ng ccam kung ilang token ang babasahin ulit ng unang sagot na iyon.

Kapag nasa Claude Code ang [ccam plugin](/docs/install#claude-code-plugin), nagpapakita ng notification ang session na huminto sa limit nito, at itinatala ng `/ccam:switch` (o `/ccam:switch work`) ang usapan para sa `ccam next --continue` nang walang kailangang hulaan.

## Kusa bang nagpapalit ng account ang ccam?

Hindi. Hindi tumatakbo ang ccam habang nagtatrabaho ka: sinisimulan nito ang Claude Code at pinapalitan ang sarili nito ng Claude Code. Ikaw ang nagpapasya kung kailan lilipat, at pinipili ng `ccam next --auto` ang account na may pinakamababang 5-oras na usage at nilalaktawan ang mga account na nasa 100%, hindi naka-log in, o kabilang sa ibang organisasyon. Kapag may `relay = true` sa `config.toml`, naghihintay naman ang ccam sa tabi ng session, kaya kusang nagpapatuloy ang `/ccam:switch`. Ang ilang ibang tool ay nagpapalit ng login ng tumatakbong session; ipinapaliwanag ng [ccam at iba pang account switcher](/compare/account-switchers) ang pagkakaiba.

## Mga dapat tandaan

- Patakbuhin ang `ccam next --continue` sa parehong folder kung saan ka nagtrabaho; naka-store ang mga usapan ayon sa project folder.
- Iwasang magpatakbo ng dalawang account sa parehong usapan nang sabay. Babalaan ka ng ccam kung susubukan mo.
- Kapag hindi shared ang kasaysayan, kinokopya ang usapan sa kabilang account at magkakahiwalay na ang landas ng dalawang kopya; nananatili sa unang account ang mga checkpoint na pinagbabalikan ng code ng `/rewind`.
- Puwedeng mas malaki ang mabawas sa limit ng bagong account sa unang mensahe kaysa karaniwan, dahil hindi pa naka-cache sa account na iyon ang usapan.
- Nakadepende sa terms ng bawat account kung puwede kang lumipat sa pagitan ng mga account nang ganito.

Tingnan ang [mga usage limit](/docs/usage-limits) para sa bawat option ng `ccam usage` at `ccam next`.
