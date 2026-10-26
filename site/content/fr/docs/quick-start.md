---
title: "Démarrage rapide"
metaTitle: "Démarrage : ajouter et changer de compte Claude Code"
description: Ajoutez deux comptes Claude Code avec ccam add, donnez-leur des raccourcis comme ccwork, passez de l'un à l'autre et choisissez ce qu'ils partagent.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

Ajoutez chaque compte avec `ccam add`, puis lancez Claude Code avec `ccam <name>` ou avec un raccourci que vous nommez vous-même, par exemple `ccwork`. C'est tout ce dont la plupart des gens ont besoin.

Vous utilisez déjà Claude Code ? Lisez [Vous utilisez déjà ~/.claude ?](#already-using-claude) à la fin de cette page avant de commencer.

## 1. Ajouter votre premier compte

```sh
ccam add
```

L'assistant vous demande un nom, par exemple `personal` ou `work`, les éléments que le compte doit partager avec vos autres comptes, et s'il faut créer un raccourci (voir l'étape 3). Il exécute ensuite `claude auth login`, qui ouvre le navigateur pour que vous puissiez vous connecter. Sur une machine distante, votre navigateur peut aboutir sur une page `http://localhost:…/callback` qui ne se charge pas ; collez cette adresse dans le terminal et ccam termine la connexion. Une fois la connexion terminée, vous voyez :

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Reprendre vos paramètres existants

Si vous avez déjà utilisé Claude Code, vos paramètres, `CLAUDE.md` et vos extensions se trouvent dans `~/.claude`, et vos serveurs MCP dans `~/.claude.json`. Copiez-les dans le groupe partagé `default` pour que chaque compte en dispose :

```sh
ccam migrate
```

ccam montre ce qu'il va copier et demande d'abord confirmation. Il se contente de lire `~/.claude`, ne copie jamais une connexion et laisse l'historique des conversations de côté, sauf si vous ajoutez `--history`. Voir [Copier une configuration existante](/docs/sharing#copy-an-existing-setup) pour plus de détails.

## 3. Ajouter le deuxième compte

```sh
ccam add work
```

L'assistant passe la question du nom, que vous avez donné sur la ligne de commande, et affiche les groupes à rejoindre, avec `default` coché. Cette fois, donnez au compte son propre raccourci, par exemple `ccwork`, pour pouvoir le lancer d'un seul mot. Après les groupes, l'assistant demande :

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Appuyez sur `y`, puis tapez un nom, ou appuyez sur Entrée pour garder celui qui est proposé :

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Connectez-vous ensuite avec l'autre compte dans le navigateur. ccam confirme le compte et son raccourci :

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Changer de compte

`ccam work` lance Claude Code avec le compte work. Si vous avez passé la question du raccourci, vous pouvez en ajouter un à tout moment :

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Désormais, chacune de ces commandes lance Claude Code avec le compte work (`2` est son numéro dans `ccam list`) :

```sh
ccwork
ccam work
ccam 2
```

Exécuter `ccam` seul ouvre le [tableau de bord](/features/dashboard) : placez-vous sur un compte avec les touches fléchées ou son numéro, puis appuyez sur Entrée.

Tout ce qui suit le raccourci ou le nom du compte est transmis tel quel à Claude Code, si bien que vous continuez à utiliser les options que vous connaissez :

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Choisir ce qui est partagé

Par défaut, les serveurs MCP, `settings.json`, le `CLAUDE.md` global et les extensions (agents, commandes, skills et plugins) sont partagés entre les comptes par le groupe `default`. L'historique des conversations ne l'est pas. Ajoutez-le au groupe si vous voulez poursuivre des conversations d'un compte à l'autre :

```sh
ccam group set default --add history
```

ccam montre ce qu'il va déplacer pour chaque compte avant de modifier quoi que ce soit. Voir [Partager avec des groupes](/features/share) pour plus de détails.

## 6. Quand un compte atteint sa limite

```sh
ccam usage
ccam next --continue
```

`ccam usage` affiche l'utilisation sur 5 heures et sur la semaine de chaque compte. `ccam next --continue` ouvre la même conversation dans le compte suivant. Dans les scripts, `ccam --fallback work -p "…"` change de compte de lui-même. Voir [Limites d'utilisation](/docs/usage-limits).

**Sièges Team et Enterprise.** La sélection automatique de compte (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord et, une fois que vous l'activez, [`auto_switch`](/features/next#auto-switch)) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise, ce qui peut enfreindre les conditions d'Anthropic. Lisez la [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) et renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier.

## Vous utilisez déjà ~/.claude ?

Enregistrez-le au lieu de vous reconnecter :

```sh
ccam import
```

Le dossier reste là où il se trouve, et `claude` seul continue de fonctionner comme avant.
