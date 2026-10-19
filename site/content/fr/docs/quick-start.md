---
title: "Démarrage rapide"
description: "Ajoutez deux comptes Claude Code avec ccam add, changez de compte avec claude1, claude2 ou ccam <name>, et choisissez ce qu'ils partagent."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: c6b1c1c3ebaf
---

Ajoutez chaque compte avec `ccam add`, puis lancez Claude Code avec `claude1`, `claude2` ou `ccam <name>`. C'est tout ce dont la plupart des gens ont besoin.

## 1. Ajouter votre premier compte

```sh
ccam add
```

L'assistant vous demande un nom, par exemple `personal` ou `work`, et les éléments que le compte doit partager avec vos autres comptes. Il exécute ensuite `claude auth login`, qui ouvre le navigateur pour que vous puissiez vous connecter. À la fin, vous voyez :

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

Pour l'instant, ccam affiche ses propres messages en coréen ; les commandes et les options sont les mêmes dans toutes les langues.

## 2. Ajouter le deuxième compte

```sh
ccam add work
```

Connectez-vous avec l'autre compte dans le navigateur. Il devient `claude2`.

## 3. Changer de compte

Chacune de ces commandes lance Claude Code avec le compte work :

```sh
claude2
ccam work
ccam 2
```

Tout ce qui suit le nom du compte est transmis tel quel à Claude Code, si bien que vous continuez à utiliser les options que vous connaissez :

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Exécuter `ccam` seul ouvre un sélecteur : utilisez les touches fléchées et Entrée, ou appuyez sur le numéro du compte.

## 4. Choisir ce qui est partagé

Par défaut, les serveurs MCP, `settings.json`, le `CLAUDE.md` global ainsi que les agents, commandes, skills et plugins sont partagés entre les comptes. L'historique des conversations ne l'est pas. Activez-le pour les deux comptes si vous voulez poursuivre des conversations de l'un à l'autre :

```sh
ccam share personal --history on
ccam share work --history on
```

ccam montre ce qu'il va déplacer avant de modifier quoi que ce soit. Voir [Partage](/docs/sharing) pour les détails.

## 5. Quand un compte atteint sa limite

```sh
ccam usage
ccam next --continue
```

`ccam usage` affiche l'utilisation sur 5 heures et sur la semaine de chaque compte. `ccam next --continue` ouvre la même conversation dans le compte suivant. Voir [Limites d'utilisation](/docs/usage-limits).

## Vous utilisez déjà ~/.claude ?

Enregistrez-le au lieu de vous reconnecter :

```sh
ccam import
```

Le dossier reste là où il se trouve, et `claude` seul continue de fonctionner comme avant.
