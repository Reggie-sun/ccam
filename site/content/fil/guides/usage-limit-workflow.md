---
title: Ano ang gagawin kapag naabot ng Claude Code ang usage limit nito
description: Alamin kung gaano pa karami ang natitira sa bawat Claude Code account at ituloy ang parehong usapan sa ibang account kapag naabot ng isa ang limit nito.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 64fe0b561bd1
---

Kapag naabot ng isang account ang limit nito, ituloy ang parehong usapan sa ibang account na pag-aari mo: i-share nang isang beses ang kasaysayan ng usapan, saka patakbuhin ang `ccs next --continue`. Kung iisa lang ang account mo, hintayin ang oras ng reset na ipinapakita ng Claude Code.

## Bago ito mangyari

I-set up ang mga account para isang command lang ang kailangan sa paglipat:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Ipinapakita na ngayon ng status line ang account at ang 5-oras na usage nito habang nagtatrabaho ka, kaya hindi ka mabibigla sa limit.

## Kapag naabot na ang limit

Lumabas sa Claude Code (`/exit`), pagkatapos:

```sh
ccs usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

Buksan ang usapan sa account na may pinakamaraming natitira:

```sh
ccs next --auto --continue
```

Magsisimula ang Claude Code bilang `work` kasama ang pinakahuling usapan ng kasalukuyang folder, kasama ang lahat ng nasabi na.

## Mga dapat tandaan

- Patakbuhin ang `ccs next --continue` sa parehong folder kung saan ka nagtrabaho; naka-store ang mga usapan ayon sa project folder.
- Iwasang magpatakbo ng dalawang account sa parehong usapan nang sabay. Babalaan ka ng ccam kung susubukan mo.
- Dapat ding nagshe-share ng kasaysayan ng usapan ang kabilang account, sa parehong group.
- Nakadepende sa terms ng bawat account kung puwede kang lumipat sa pagitan ng mga account nang ganito.

Tingnan ang [mga usage limit](/docs/usage-limits) para sa bawat option ng `ccs usage` at `ccs next`.
