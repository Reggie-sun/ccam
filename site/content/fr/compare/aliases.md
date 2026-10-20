---
title: "ccam comparé aux alias shell"
description: "Alias shell avec CLAUDE_CONFIG_DIR ou ccam pour utiliser plusieurs comptes Claude Code : comparaison côte à côte."
updated: 2026-10-20
order: 1
translation: machine
sourceHash: ea24bb091caf
---

Un alias shell par compte est le moyen le plus simple d'utiliser plusieurs comptes Claude Code. ccam fait la même chose en interne et ajoute le partage, le changement de compte et le suivi de l'utilisation.

## La configuration avec des alias

Ajoutez une ligne par compte supplémentaire dans `~/.bashrc` ou `~/.zshrc`, ouvrez un nouveau terminal et lancez l'alias une fois pour vous connecter :

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` reste votre premier compte et `claude-work` lance le second. C'est la méthode [décrite dans la documentation de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## La configuration avec ccam

Installez ccam et ajoutez le compte, ici avec sa propre commande raccourcie :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` et `ccwork` lancent le compte, et tous deux fonctionnent aussi dans les scripts.

## Côte à côte

| | Alias shell | ccam |
| --- | --- | --- |
| Connexions séparées | Oui | Oui |
| Configuration par compte | Modifier `~/.bashrc`, se connecter | `ccam add` |
| Fonctionne dans les scripts et d'autres programmes | Non, les alias n'existent que dans le shell | Oui, `ccam work` et les raccourcis comme `ccwork` sont de vraies commandes |
| Paramètres, `CLAUDE.md` et extensions partagés | Copie ou lien à la main | Élément par élément, avec sauvegardes |
| Serveurs MCP partagés | À la main, compte par compte | `ccam mcp add` une seule fois |
| Poursuivre une conversation dans un autre compte | Non | `ccam next --continue` |
| Utilisation de tous les comptes | Non | `ccam usage` |
| Choisir le compte qui a le plus de marge | Non | `ccam next --auto` |
| Nom du compte dans la ligne d'état | Non | `ccam statusline install` |
| Réparer les liens cassés | À la main | `ccam doctor --fix` |
| Logiciel supplémentaire | Aucun | Un seul binaire statique |

Choisissez les alias si vous avez deux comptes et ne voulez rien partager entre eux. Choisissez ccam si vous voulez que les comptes se comportent comme une seule configuration. Vous pouvez commencer avec des alias et passer à ccam plus tard : `ccam import` enregistre les dossiers `~/.claude-*` existants là où ils se trouvent.

[Utiliser plusieurs comptes sur un même ordinateur](/guides/multiple-claude-code-accounts) explique les deux approches étape par étape, et [ccam et les autres sélecteurs de compte](/compare/account-switchers) compare ccam aux outils qui échangent une seule connexion.
