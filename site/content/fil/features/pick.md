---
title: Pagpili ng account para sa mga script
description: Ipi-print ng ccam pick ang account na may pinakamaraming natitira para sa script, iba't ibang account sa sabay na worker, at exit 75 kapag wala.
updated: 2026-10-24
order: 12
translation: machine
sourceHash: a8baeb64de99
---

Ipi-print ng `ccam pick` ang pangalan ng naka-log in na account na may pinakamaraming natitira, para masimulan ng isang script ang Claude Code gamit ito. Ang mga worker na sabay na sinimulan ay may kanya-kanyang account habang may sapat pang account na may natitira, kaya kumakalat sa mga account mo ang isang batch ng mga `claude -p` job sa halip na maubos ang iisa.

**Mga seat ng Team at Enterprise.** Lumilipat din ang awtomatikong pagpili (`ccam next`, `--fallback`, `ccam pick`, ang `n` ng dashboard) sa pagitan ng mga seat ng iisang Claude Team o Enterprise na organisasyon. Sinasabi ng help center ng Anthropic na "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (nakatalaga ang bawat seat ng Team plan sa iisang partikular na tao, at siya lang ang puwedeng gumamit nito) at "Users can't route their work through another user's account to use that user's remaining usage" (hindi puwedeng idaan ng mga user ang trabaho nila sa account ng ibang user para magamit ang natitirang usage ng user na iyon) ([mga Team plan](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ganoon din sa [mga Enterprise plan](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kaya maaaring lumabag sa terms ng Anthropic ang paggamit ng awtomatikong pagpili sa pagitan ng mga seat. Magtanong muna sa iyong organisasyon at sa Anthropic bago ka umasa rito. Saklaw ng [paalala tungkol sa terms sa FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) ang mga personal na plan (Pro, Max).

## Paggamit

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Mga exit code: `0` na may pangalan sa standard output, `75` (`EX_TEMPFAIL`) kapag walang account na may natitira o walang naka-log in, `1` para sa ibang error gaya ng maling flag.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--json` | I-print ang account, ang 5-oras at lingguhang usage nito, kung saan galing ang mga numero at kung bakit ito napili | off |
| `--model <model>` | Bilangin din ang lingguhang window ng model na iyon, gaya ng `opus` | wala |
| `--max-usage <1-100>` | Unahin ang mga account na mas mababa sa usage na ito | `auto_max_usage`, kung wala ay 100 |
| `--strategy headroom\|expiring\|balance` | Pinakamababang 5-oras na usage, ang pinakamalapit na 5-oras na reset, o ang lingguhang quota na unang masasayang kapag hindi nagamit | `auto_strategy`, kung wala ay `headroom` |
| `--cross-org` | Pumili rin ng mga account ng ibang organisasyon | `cross_org`, kung wala ay off |
| `--wait[=<duration>]` | Kapag walang account na may natitira, hintayin ang unang mag-reset, nang hindi hihigit sa tagal na ito. Inilalagay ang value pagkatapos ng `=`: tinatanggihan ang `--wait 30m` kasama ang pahiwatig na `--wait takes its value after '=': --wait=30m` | off; `5h` kapag ibinigay nang mag-isa |
| `--dry-run` | Pumili nang hindi itinatala ang pagpili | off |

## Paano ito gumagana

1. Inililista ng ccam ang mga naka-log in na account (walang limit ang mga API key account at Console login at hindi isinasama) sa organisasyon ng account na sisimulan ng folder na ito (tingnan ang [mga organisasyon](/features/next#auto)); hindi rin isinasama ang folder na naka-log in sa parehong Claude account niyon.
2. Kumukuha ito ng lock (`~/.local/state/ccam/picks.lock`), kaya isa-isang nangyayari ang mga sabay na pagpili.
3. Sabay-sabay nitong tinitingnan ang usage ng bawat account, sa parehong paraan ng `ccam usage`, gamit ang 60-segundong cache. Itinuturing na 100% ang window na iniulat ng Claude Code na puno na hanggang sa reset nito ([ang tala ng limit](/features/usage#the-limit-record)), kaya nananatiling hindi kasama ang account na kaaabot lang sa limit nito kahit nagpapakita pa ng puwang ang mga naka-save nitong numero.
4. Inaalis muna ang mga account na napili sa loob ng huling 2 minuto (`~/.local/state/ccam/picks/<name>`; hindi binibilang ang pagpiling may petsang lampas sa oras ng orasan), at ang [mga patakaran ng pagpili](/features/next#auto) ng `ccam next --auto` ang pumipili sa natitira: inuuna ang mas mababa sa `auto_max_usage`, ayon sa strategy, at hindi kailanman ang account na nasa 100%.
5. Kung kapipili lang ng bawat account na may natitira, pipili ulit ang mga patakaran mula sa lahat ng iyon, kaya magkakahati sa isang account ang ikaapat na worker sa tatlong account.
6. Itinatala ang pagpili at binibitawan ang lock, saka ipi-print ang pangalan.

Mahalaga ang 2-minutong window dahil natatagalan bago lumabas sa mga numero ang usage ng bagong job: kung wala ito, makikita ng limang sabay na sinimulang worker ang iisang account bilang may pinakamaraming natitira.

## Kaligtasan at mga espesyal na kaso

- Nagbabasa lang ng usage at sumusulat ng oras ng pagpili ang `ccam pick`; hindi ito kailanman nagre-renew o nagbabasa ng token nang higit sa ginagawa ng `ccam usage`.
- Hindi nirereserba ng pagpili ang account. Puwede pa ring magsimula ang dalawang script sa iisang account pagkalipas ng 2 minuto; gamitin ang `--fallback` para makabawi kapag naabot ng account na iyon ang limit nito.
- Walang hawak na lock ang paghihintay gamit ang `--wait`, kaya tuloy ang ibang pagpili habang naghihintay.
- Hindi pinapansin ng `ccam pick` ang [`ccam folder`](/features/folder): usage lang ang batayan nito sa pagpili.

## Kaugnay

- [Paglipat sa susunod na account](/features/next)
- [Patakbuhin ulit sa limit](/features/next#run-again-at-the-limit)
- [Mga usage limit](/features/usage)
- [Mga API key account](/features/api)
