---
title: "ccam comparé aux alias shell"
description: "Alias shell avec CLAUDE_CONFIG_DIR ou ccam pour utiliser plusieurs comptes Claude Code : comparaison côte à côte."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e9902676e0ed
---

Un alias shell par compte est le moyen le plus simple d'utiliser plusieurs comptes Claude Code. ccam fait la même chose en interne et ajoute le partage, le changement de compte et le suivi de l'utilisation.

| | Alias shell | ccam |
| --- | --- | --- |
| Connexions séparées | Oui | Oui |
| Configuration par compte | Modifier `~/.bashrc`, se connecter | `ccam add` |
| Fonctionne dans les scripts et d'autres programmes | Non, les alias n'existent que dans le shell | Oui, `claude2` est une vraie commande |
| Paramètres, `CLAUDE.md` et extensions partagés | Copie ou lien à la main | Élément par élément, avec sauvegardes |
| Serveurs MCP partagés | À la main, compte par compte | `ccam mcp add` une seule fois |
| Poursuivre une conversation dans un autre compte | Non | `ccs next --continue` |
| Utilisation de tous les comptes | Non | `ccs usage` |
| Choisir le compte qui a le plus de marge | Non | `ccs next --auto` |
| Nom du compte dans la ligne d'état | Non | `ccam statusline install` |
| Réparer les liens cassés | À la main | `ccam doctor --fix` |
| Logiciel supplémentaire | Aucun | Un seul binaire statique |

Choisissez les alias si vous avez deux comptes et ne voulez rien partager entre eux. Choisissez ccam si vous voulez que les comptes se comportent comme une seule configuration. Vous pouvez commencer avec des alias et passer à ccam plus tard : `ccam import` enregistre les dossiers `~/.claude-*` existants là où ils se trouvent.

[Utiliser deux comptes sur un seul ordinateur](/guides/multiple-claude-code-accounts) explique les deux approches étape par étape.
