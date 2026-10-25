---
title: "Démarrage rapide"
metaTitle: "Démarrage : ajouter et changer de compte Claude Code"
description: Ajoutez deux comptes Claude Code avec ccam add, donnez-leur des raccourcis comme ccwork, passez de l'un à l'autre et choisissez ce qu'ils partagent.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: 0570ef27e4ce
---

Ajoutez chaque compte avec `ccam add`, puis lancez Claude Code avec `ccam <name>` ou avec un raccourci que vous nommez vous-même, par exemple `ccwork`. C'est tout ce dont la plupart des gens ont besoin.

## 1. Ajouter votre premier compte

```sh
ccam add
```

L'assistant vous demande un nom, par exemple `personal` ou `work`, les éléments que le compte doit partager avec vos autres comptes, et s'il faut créer un raccourci (voir l'étape 4). Il exécute ensuite `claude auth login`, qui ouvre le navigateur pour que vous puissiez vous connecter. Sur une machine distante, le navigateur peut aboutir sur une page `http://localhost:…/callback` qui ne se charge pas ; collez cette adresse dans le terminal et ccam termine la connexion. À la fin, vous voyez :

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

Exécuter `ccam` seul ouvre le [tableau de bord](/features/dashboard) : placez-vous sur un compte avec les touches fléchées ou son numéro, puis appuyez sur Entrée.

Tout ce qui suit le raccourci ou le nom du compte est transmis tel quel à Claude Code, si bien que vous continuez à utiliser les options que vous connaissez :

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Choisir ce qui est partagé

Par défaut, les serveurs MCP, `settings.json`, le `CLAUDE.md` global ainsi que les agents, commandes, skills et plugins sont partagés entre les comptes par le groupe `default`. L'historique des conversations ne l'est pas. Ajoutez-le au groupe si vous voulez poursuivre des conversations d'un compte à l'autre :

```sh
ccam group set default --add history
```

ccam montre ce qu'il va déplacer pour chaque compte avant de modifier quoi que ce soit. Voir [Partager avec des groupes](/features/share) pour les détails.

## 6. Quand un compte atteint sa limite

```sh
ccam usage
ccam next --continue
```

`ccam usage` affiche l'utilisation sur 5 heures et sur la semaine de chaque compte. `ccam next --continue` ouvre la même conversation dans le compte suivant. Dans les scripts, `ccam --fallback work -p "…"` change de compte de lui-même. Voir [Limites d'utilisation](/docs/usage-limits).

**Sièges Team et Enterprise.** Le choix automatique (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord, et [`auto_switch`](/features/next#auto-switch) lorsque vous l'activez) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Utiliser le choix automatique entre plusieurs sièges peut donc enfreindre les conditions d'Anthropic. C'est à vous de décider si vous l'utilisez dans ce cas : renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier. Les offres personnelles (Pro, Max) sont traitées dans la [note de la FAQ sur les conditions](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Vous utilisez déjà ~/.claude ?

Enregistrez-le au lieu de vous reconnecter :

```sh
ccam import
```

Le dossier reste là où il se trouve, et `claude` seul continue de fonctionner comme avant.
