---
title: "Démarrage rapide"
metaTitle: "Démarrage : ajouter et changer de compte Claude Code"
description: Ajoutez deux comptes Claude Code avec ccam add, donnez-leur des raccourcis comme ccwork, passez de l'un à l'autre et choisissez ce qu'ils partagent.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 745a9904e32a
---

Ajoutez chaque compte avec `ccam add`, puis lancez Claude Code avec `ccam <name>` ou avec un raccourci que vous nommez vous-même, par exemple `ccwork`. C'est tout ce dont la plupart des gens ont besoin.

## 1. Ajouter votre premier compte

```sh
ccam add
```

L'assistant vous demande un nom, par exemple `personal` ou `work`, les éléments que le compte doit partager avec vos autres comptes, et s'il faut créer un raccourci (voir l'étape 4). Il exécute ensuite `claude auth login`, qui ouvre le navigateur pour que vous puissiez vous connecter. À la fin, vous voyez :

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Reprendre vos paramètres existants

Si vous utilisiez déjà Claude Code, vos paramètres, `CLAUDE.md`, vos extensions et vos serveurs MCP se trouvent dans `~/.claude`. Copiez-les dans les paramètres partagés pour que chaque compte en dispose :

```sh
ccam migrate
```

ccam montre ce qu'il va copier et demande d'abord confirmation. Il se contente de lire `~/.claude`, ne copie jamais une connexion et laisse l'historique des conversations de côté, sauf si vous ajoutez `--history`. Voir [Copier une configuration existante](/docs/sharing#copy-an-existing-setup) pour les détails.

## 3. Ajouter le deuxième compte

```sh
ccam add work
```

Connectez-vous avec l'autre compte dans le navigateur.

## 4. Changer de compte

`ccam work` lance Claude Code avec le compte work. Pour plus court, donnez au compte son propre raccourci, par exemple `ccwork`. Après les éléments partagés, l'assistant demande :

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Appuyez sur `y`, puis tapez un nom, ou appuyez sur Entrée pour garder celui qui est proposé :

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Après la connexion, ccam le confirme :

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Si vous avez passé la question, vous pouvez ajouter un raccourci à tout moment :

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Désormais, chacune de ces commandes lance Claude Code avec le compte work :

```sh
ccwork
ccam work
ccam 2
```

Exécuter `ccam` seul ouvre un sélecteur : utilisez les touches fléchées et Entrée, ou appuyez sur le numéro du compte.

Tout ce qui suit le raccourci ou le nom du compte est transmis tel quel à Claude Code, si bien que vous continuez à utiliser les options que vous connaissez :

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Choisir ce qui est partagé

Par défaut, les serveurs MCP, `settings.json`, le `CLAUDE.md` global ainsi que les agents, commandes, skills et plugins sont partagés entre les comptes. L'historique des conversations ne l'est pas. Activez-le pour les deux comptes si vous voulez poursuivre des conversations de l'un à l'autre :

```sh
ccam share personal --history on
ccam share work --history on
```

ccam montre ce qu'il va déplacer avant de modifier quoi que ce soit. Voir [Partage](/docs/sharing) pour les détails.

## 6. Quand un compte atteint sa limite

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
