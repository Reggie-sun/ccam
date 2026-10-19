---
title: "Que faire quand Claude Code atteint sa limite d'utilisation"
description: "Vérifiez la marge restante de chaque compte Claude Code et poursuivez la même conversation sur un autre compte quand l'un d'eux atteint sa limite."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

Quand un compte atteint sa limite, poursuivez la même conversation dans un autre compte qui vous appartient : partagez l'historique des conversations une fois, puis exécutez `ccam next --continue`. Si vous n'avez qu'un seul compte, attendez l'heure de réinitialisation qu'affiche Claude Code.

## Avant que cela arrive

Configurez les comptes pour que le changement ne demande qu'une seule commande :

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

La ligne d'état affiche désormais le compte et son utilisation sur 5 heures pendant que vous travaillez : la limite ne vous prend donc pas par surprise.

## Quand la limite est atteinte

Quittez Claude Code (`/exit`), puis :

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Ouvrez la conversation dans le compte qui a le plus de marge :

```sh
ccam next --auto --continue
```

Claude Code démarre en tant que `work` avec la conversation la plus récente du dossier courant, y compris tout ce qui a été dit jusque-là.

## Points à retenir

- Exécutez `ccam next --continue` dans le dossier où vous travailliez ; les conversations sont stockées par dossier de projet.
- Évitez d'utiliser deux comptes sur la même conversation en même temps. ccam vous avertit si vous essayez.
- L'autre compte doit aussi partager l'historique des conversations, dans le même groupe.
- Le droit de passer ainsi d'un compte à l'autre dépend des conditions de chaque compte.

Voir [Limites d'utilisation](/docs/usage-limits) pour toutes les options de `ccam usage` et `ccam next`.
