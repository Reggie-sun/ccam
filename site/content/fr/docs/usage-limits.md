---
title: "Limites d'utilisation et changement de compte"
description: "Consultez l'utilisation sur 5 heures et hebdomadaire de chaque compte Claude Code avec ccs usage, et ouvrez celui qui a le plus de marge avec ccs next --auto."
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a9ee77936668
---

`ccs usage` indique quelle part des limites sur 5 heures et sur la semaine chaque compte a utilisée. `ccs next --auto` ouvre le compte qui a le plus de marge, et `--continue` reprend la conversation en cours.

```console
$ ccs usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccs next --auto 로 바로 열기)
```

Les colonnes sont : le compte, l'utilisation sur 5 heures, le temps restant avant sa réinitialisation, l'utilisation sur la semaine, le temps restant avant sa réinitialisation et l'ancienneté des chiffres.

## D'où viennent les chiffres

ccam interroge le point de terminaison d'utilisation d'Anthropic pour tous les comptes en même temps, avec la propre connexion de chaque compte. Il attend au plus 3 secondes par compte et affiche chaque compte dès qu'il répond.

Ce point de terminaison n'est pas documenté publiquement et peut changer. Quand une requête échoue, ccam affiche les derniers chiffres connus et les signale comme provenant du cache. `ccam doctor` vous prévient si le format de la réponse a changé.

Les chiffres récupérés il y a moins de 60 secondes sont réutilisés. Options :

| Option | Effet |
| --- | --- |
| `--watch [seconds]` | Actualiser en continu, toutes les 60 secondes par défaut |
| `--json` | Sortie lisible par une machine |
| `-v` | Afficher aussi les limites hebdomadaires propres à chaque modèle et la provenance de chaque chiffre |
| `--no-cache` | Toujours redemander |
| `--source cache` | Ne jamais demander ; afficher uniquement les chiffres enregistrés |

## Connexions expirées

Quand le jeton d'accès d'un compte a expiré, ccam ne le renouvelle pas, car le renouveler pendant qu'une session Claude Code utilise la même connexion pourrait déconnecter cette session. Lancez le compte une fois avec `ccs <name>` et Claude Code le renouvellera.

`ccs usage --refresh` renouvelle lui-même les jetons expirés, mais uniquement pour les comptes qui ne sont pas en cours d'exécution. Il ne renouvelle jamais les connexions stockées dans le trousseau macOS ni dans un `~/.claude` importé.

## L'utilisation dans la ligne d'état

```sh
ccam statusline install
```

Claude Code affiche alors le nom du compte et son utilisation dans la ligne d'état, par exemple `[work] · Opus · 5h 24% · 7d 41%`. Si vous avez déjà une commande de ligne d'état, ccam la conserve et place le nom du compte devant. La ligne d'état enregistre aussi ses chiffres, que `ccs usage` utilise quand le point de terminaison ne répond pas.

## Changer de compte

```sh
ccs next              # the account after the one you used last
ccs next --auto       # the account with the lowest 5-hour usage
ccs next --continue   # the next account that shares your conversation history
```

`ccs next` ignore les comptes qui ne sont pas connectés, et `--auto` ignore les comptes à 100 %. Pour poursuivre une conversation, les deux comptes doivent partager l'historique des conversations ; voir [Partage](/docs/sharing).

Le [guide sur la limite d'utilisation](/guides/usage-limit-workflow) présente un exemple complet.
