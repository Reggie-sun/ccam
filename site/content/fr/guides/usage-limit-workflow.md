---
title: "Que faire quand Claude Code atteint sa limite"
description: "Attendez la réinitialisation ou poursuivez la même conversation sur un autre compte Claude Code. Les limites sur 5 heures et hebdomadaires, étape par étape."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: ff765a189332
---

Quand un compte atteint sa limite d'utilisation (souvent appelée « rate limit »), vous avez deux options : attendre que la limite se réinitialise, ou poursuivre dès maintenant la même conversation sur un autre compte qui vous appartient. Avec ccam, la seconde ne demande qu'une commande, `ccam next --auto --continue`, une fois l'historique des conversations partagé entre les comptes.

## Quelles sont les limites de Claude Code ?

Avec les forfaits Pro et Max, Claude Code a une [limite de session sur cinq heures et une limite hebdomadaire](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Les deux sont partagées avec Claude sur le web et dans les applications, et l'utilisation compte pour les deux à la fois. Quand l'une est épuisée, Claude Code [affiche un message](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) comme « You've hit your session limit · resets 3:45pm ». Changer de modèle n'aide pas, car les limites couvrent tous les modèles.

## Option 1 : attendre la réinitialisation

Si vous n'avez qu'un compte, attendez. Les versions récentes de Claude Code [attendent dans la session ouverte](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) et reprennent la tâche d'elles-mêmes une fois la limite réinitialisée. `/usage` dans Claude Code indique quelle part des limites de votre forfait vous avez utilisée.

## Option 2 : poursuivre tout de suite sur un autre compte

Si vous avez un second compte, vous pouvez continuer immédiatement la même conversation. Configurez-le une fois :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

La ligne d'état affiche désormais le compte et son utilisation sur 5 heures pendant que vous travaillez : la limite ne vous prend donc pas par surprise.

Quand la limite est atteinte, quittez Claude Code (`/exit`), puis :

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Ouvrez la conversation dans le compte qui a le plus de marge :

```sh
ccam next --auto --continue
```

Claude Code démarre en tant que `work` avec la conversation la plus récente du dossier courant, y compris tout ce qui a été dit jusque-là.

## ccam change-t-il de compte automatiquement ?

Non. ccam ne tourne pas pendant que vous travaillez : il lance Claude Code et se remplace par lui. C'est vous qui décidez quand changer, et `ccam next --auto` choisit le compte dont l'utilisation sur 5 heures est la plus faible, en ignorant les comptes à 100 % ou non connectés. Certains autres outils échangent plutôt la connexion d'une session en cours ; [ccam et les autres sélecteurs de compte](/compare/account-switchers) explique la différence.

## Points à retenir

- Exécutez `ccam next --continue` dans le dossier où vous travailliez ; les conversations sont stockées par dossier de projet.
- Évitez d'utiliser deux comptes sur la même conversation en même temps. ccam vous avertit si vous essayez.
- L'autre compte doit aussi partager l'historique des conversations, dans le même groupe.
- Le premier message sur le nouveau compte peut consommer plus de sa limite que d'habitude, car ce compte n'a pas encore la conversation en cache.
- Le droit de passer ainsi d'un compte à l'autre dépend des conditions de chaque compte.

Voir [Limites d'utilisation](/docs/usage-limits) pour toutes les options de `ccam usage` et `ccam next`.
