---
title: "Choisir un compte pour les scripts"
description: "ccam pick affiche le compte qui a le plus de marge pour un script, donne des comptes différents aux workers lancés ensemble et quitte avec 75 sans marge."
updated: 2026-10-25
order: 12
translation: machine
sourceHash: 5c4f74db9714
---

`ccam pick` affiche le nom du compte connecté qui a le plus de marge, pour qu'un script lance Claude Code avec. Les workers lancés en même temps reçoivent chacun un compte différent tant qu'il y a assez de comptes avec de la marge, si bien qu'un lot de tâches `claude -p` se répartit sur vos comptes au lieu d'en épuiser un seul.

**Sièges Team et Enterprise.** Le choix automatique (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord, et [`auto_switch`](/features/next#auto-switch) lorsque vous l'activez) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Utiliser le choix automatique entre plusieurs sièges peut donc enfreindre les conditions d'Anthropic. C'est à vous de décider si vous l'utilisez dans ce cas : renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier. Les offres personnelles (Pro, Max) sont traitées dans la [note de la FAQ sur les conditions](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Utilisation

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

Codes de sortie : `0` avec un nom sur la sortie standard, `75` (`EX_TEMPFAIL`) quand aucun compte n'a de marge ou qu'aucun n'est connecté, `1` pour les autres erreurs, comme une option invalide.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--json` | Afficher le compte, son utilisation sur 5 heures et hebdomadaire, l'origine des chiffres et la raison de ce choix | désactivé |
| `--model <model>` | Compter aussi la fenêtre hebdomadaire de ce modèle, comme `opus` | aucun |
| `--max-usage <1-100>` | Préférer les comptes sous ce niveau d'utilisation | `auto_max_usage`, sinon 100 |
| `--strategy headroom\|expiring\|balance` | L'utilisation sur 5 heures la plus basse, la réinitialisation sur 5 heures la plus proche, ou le quota hebdomadaire qui serait perdu en premier | `auto_strategy`, sinon `headroom` |
| `--cross-org` | Choisir aussi des comptes d'autres organisations | `cross_org`, sinon désactivé |
| `--wait[=<duration>]` | Quand aucun compte n'a de marge, attendre la réinitialisation du premier, au plus cette durée. La valeur se place après `=` : `--wait 30m` est refusé avec l'indication `--wait takes its value after '=': --wait=30m` | désactivé ; `5h` s'il est donné seul |
| `--dry-run` | Choisir sans enregistrer le choix | désactivé |

## Fonctionnement

1. ccam liste les comptes connectés (les comptes à clé d'API et les connexions Console n'ont pas de limites et sont laissés de côté) dans l'organisation du compte que lancerait ce dossier (voir [organisations](/features/next#auto)) ; un dossier connecté au même compte Claude que celui-ci est lui aussi laissé de côté.
2. Il prend un verrou (`~/.local/state/ccam/picks.lock`), si bien que les choix lancés en même temps se font l'un après l'autre.
3. Il consulte l'utilisation de tous les comptes en même temps, de la même façon que `ccam usage`, avec le cache de 60 secondes. Une fenêtre que Claude Code a signalée pleine compte pour 100 % jusqu'à sa réinitialisation ([le relevé de limite](/features/usage#the-limit-record)), si bien qu'un compte qui vient d'atteindre sa limite reste écarté même quand ses chiffres enregistrés montrent de la marge.
4. Les comptes choisis au cours des 2 dernières minutes, par `ccam pick` ou par [`auto_switch`](/features/next#auto-switch) qui y a déplacé une session (`~/.local/state/ccam/picks/<name>` ; un choix daté après l'heure de l'horloge ne compte pas), sont d'abord écartés, et les [règles de choix](/features/next#auto) de `ccam next --auto` choisissent parmi les autres : d'abord sous `auto_max_usage`, selon la stratégie, jamais un compte à 100 %. Quand cela écarte le compte que les règles prendraient en premier, la `reason` de `--json` le signale : `30% of its 5 hours used; work was picked a moment ago`.
5. Si tous les comptes qui ont de la marge ont été choisis récemment, les règles choisissent de nouveau parmi tous, si bien qu'un quatrième worker sur trois comptes en partage un.
6. Le choix est enregistré et le verrou libéré, puis le nom est affiché.

La fenêtre de 2 minutes compte, car l'utilisation d'une nouvelle tâche met un moment à apparaître dans les chiffres : sans elle, cinq workers lancés ensemble verraient tous le même compte comme celui qui a le plus de marge.

## Sécurité et cas particuliers

- `ccam pick` ne fait que lire l'utilisation et écrire l'heure du choix ; il ne renouvelle ni ne lit jamais de jetons au-delà de ce que fait `ccam usage`.
- Un choix ne réserve pas le compte. Deux scripts peuvent encore démarrer sur le même compte après 2 minutes ; utilisez `--fallback` pour reprendre la main quand ce compte atteint sa limite.
- L'attente avec `--wait` ne garde aucun verrou, si bien que les autres choix continuent pendant ce temps.
- `ccam pick` ignore [`ccam folder`](/features/folder) : il choisit uniquement selon l'utilisation.

## Voir aussi

- [Passer au compte suivant](/features/next)
- [Relancer à la limite](/features/next#run-again-at-the-limit)
- [Limites d'utilisation](/features/usage)
- [Comptes à clé d'API](/features/api)
