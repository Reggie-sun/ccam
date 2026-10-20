---
title: Pagpili ng account para sa mga script
description: Ipi-print ng ccam pick ang account na may pinakamaraming natitira para sa script, iba't ibang account sa sabay na worker, at exit 75 kapag wala.
updated: 2026-10-20
order: 12
translation: machine
sourceHash: 65a5a7e179a5
---

Ipi-print ng `ccam pick` ang pangalan ng naka-log in na account na may pinakamaraming natitira, para masimulan ng isang script ang Claude Code gamit ito. Ang mga worker na sabay na sinimulan ay may kanya-kanyang account habang may sapat pang account na may natitira, kaya kumakalat sa mga account mo ang isang batch ng mga `claude -p` job sa halip na maubos ang iisa.

## Paggamit

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
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
  "source": "live"
}
```

Mga exit code: `0` na may pangalan sa standard output, `75` (`EX_TEMPFAIL`) kapag walang account na may natitira o walang naka-log in, `1` para sa ibang error gaya ng maling flag.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--json` | I-print ang account, ang 5-oras at lingguhang usage nito, at kung saan galing ang mga numero | off |
| `--model <model>` | Bilangin din ang lingguhang window ng model na iyon, gaya ng `opus` | wala |
| `--max-usage <1-100>` | Unahin ang mga account na mas mababa sa usage na ito | `auto_max_usage`, kung wala ay 100 |
| `--strategy headroom\|expiring` | Pinakamababang 5-oras na usage, o ang pinakamalapit na 5-oras na reset | `auto_strategy`, kung wala ay `headroom` |
| `--wait[=<duration>]` | Kapag walang account na may natitira, hintayin ang unang mag-reset, nang hindi hihigit sa tagal na ito | off; `5h` kapag ibinigay nang mag-isa |
| `--dry-run` | Pumili nang hindi itinatala ang pagpili | off |

## Paano ito gumagana

1. Inililista ng ccam ang mga naka-log in na account (walang limit ang mga API key account at hindi isinasama).
2. Kumukuha ito ng lock (`~/.local/state/ccam/picks.lock`), kaya isa-isang nangyayari ang mga sabay na pagpili.
3. Sabay-sabay nitong tinitingnan ang usage ng bawat account, sa parehong paraan ng `ccam usage`, gamit ang 60-segundong cache.
4. Inaalis muna ang mga account na napili sa loob ng huling 2 minuto (`~/.local/state/ccam/picks/<name>`), at ang [mga patakaran ng pagpili](/features/next#auto) ng `ccam next --auto` ang pumipili sa natitira: inuuna ang mas mababa sa `auto_max_usage`, ayon sa strategy, at hindi kailanman ang account na nasa 100%.
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
