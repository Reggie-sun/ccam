---
title: "Limites d'utilisation et changement de compte"
description: "Consultez l'utilisation sur 5 heures et hebdomadaire de chaque compte Claude Code avec ccam usage, et ouvrez celui qui a le plus de marge avec ccam next --auto."
updated: 2026-10-21
order: 5
translation: machine
sourceHash: fadb998482da
---

`ccam usage` indique quelle part de la limite sur 5 heures, de la limite hebdomadaire et de la limite hebdomadaire propre à chaque modèle (comme Fable) chaque compte a utilisée. `ccam next --auto` ouvre le compte qui a le plus de marge, et `--continue` reprend la conversation en cours.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Chaque compte a une ligne qui indique l'ancienneté de ses chiffres, puis une barre par limite avec le temps restant avant sa réinitialisation : la fenêtre de 5 heures, la fenêtre hebdomadaire de tous les modèles et la fenêtre hebdomadaire propre à chaque modèle, comme Fable. À partir de 80 % une barre passe au jaune, et à 100 % au rouge. La dernière ligne nomme le compte que `ccam next --auto` ouvrirait.

## D'où viennent les chiffres

ccam interroge le point de terminaison d'utilisation d'Anthropic pour tous les comptes en même temps, avec la propre connexion de chaque compte. Il attend au plus 3 secondes par compte et affiche chaque compte dès qu'il répond.

Ce point de terminaison n'est pas documenté publiquement et peut changer. Quand une requête échoue, ccam affiche les derniers chiffres connus et les signale comme provenant du cache. Une fenêtre réinitialisée depuis l'enregistrement de ces chiffres affiche 0 %. `ccam doctor` vous prévient si le format de la réponse a changé.

Les chiffres récupérés il y a moins de 60 secondes sont réutilisés. Options :

| Option | Effet |
| --- | --- |
| `--watch [seconds]` | Actualiser en continu, toutes les 60 secondes par défaut |
| `--json` | Sortie lisible par une machine |
| `-v` | Afficher aussi la provenance des chiffres de chaque compte |
| `--no-cache` | Toujours redemander |
| `--source cache` | Ne jamais demander ; afficher uniquement les chiffres enregistrés |

## Connexions expirées

Quand le jeton d'accès d'un compte a expiré, ccam ne le renouvelle pas, car le renouveler pendant qu'une session Claude Code utilise la même connexion pourrait déconnecter cette session. Lancez le compte une fois avec `ccam <name>` et Claude Code le renouvellera.

`ccam usage --refresh` renouvelle lui-même les jetons expirés, mais uniquement pour les comptes qui ne sont pas en cours d'exécution. Il ne renouvelle jamais les connexions stockées dans le trousseau macOS ni dans un `~/.claude` importé.

## L'utilisation dans la ligne d'état

```sh
ccam statusline install
```

Claude Code affiche alors le nom du compte et son utilisation dans la ligne d'état, par exemple `[work] · Opus · 5h 24% · 7d 41%`. À partir de 90 %, la ligne suggère aussi `ccam next --auto --continue`. Si vous avez déjà une commande de ligne d'état, ccam la conserve et place le nom du compte devant. La ligne d'état enregistre aussi ses chiffres, que `ccam usage` utilise quand le point de terminaison ne répond pas.

## Changer de compte

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` ignore les comptes qui ne sont pas connectés, et `--auto` ignore les comptes à 100 % et ceux d'une autre organisation. Quand les deux comptes partagent l'historique des conversations, `--continue` ouvre la même conversation ; sinon, Claude Code la copie dans le nouveau compte. Voir [Partage](/docs/sharing) et [Passer au compte suivant](/features/next#continue).

Pour les scripts et les tâches planifiées qui exécutent `claude -p`, placez `--fallback` avant le nom du compte. Quand ce compte atteint sa limite, ccam relance la même commande avec le compte qui a le plus de marge, et seule la sortie de l'exécution qui a abouti est affichée :

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Voir [Relancer à la limite](/features/next#run-again-at-the-limit) pour savoir comment ccam décide.

Le [guide sur la limite d'utilisation](/guides/usage-limit-workflow) présente un exemple complet.
