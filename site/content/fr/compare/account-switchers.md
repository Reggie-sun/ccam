---
title: "ccam et les autres sélecteurs de compte Claude Code"
description: "Comparaison de ccam avec claude-swap, CCS, clauth, CAAM, cux et d'autres outils pour changer de compte Claude Code, et quelle approche convient à quel usage."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Les outils pour changer de compte Claude Code suivent l'une de deux approches. La plupart échangent la connexion à l'intérieur de l'unique dossier `~/.claude` : un seul compte est donc actif à la fois. D'autres, dont ccam, donnent à chaque compte son propre dossier avec `CLAUDE_CONFIG_DIR` : chaque compte reste connecté et plusieurs peuvent tourner en même temps. ccam y ajoute le partage des paramètres et des serveurs MCP, et poursuit une conversation dans le dossier d'un autre compte.

## Échanger une connexion ou séparer les dossiers ?

**Échanger la connexion** conserve un seul jeu de paramètres, de serveurs MCP et d'historique, puisqu'il n'y a qu'un dossier. Le changement est rapide, et une session en cours peut reprendre la nouvelle connexion. Mais un seul compte est actif à la fois, pour tous les terminaux et éditeurs en même temps.

**Des dossiers séparés**, c'est la méthode que [décrit la documentation de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Chaque compte reste connecté, et vous pouvez lancer côte à côte une session professionnelle et une session personnelle. La contrepartie : chaque dossier démarre avec ses propres paramètres, serveurs MCP et historique, sauf si un outil les partage.

## Comparaison des outils

D'après la documentation de chaque projet au 2026-10-19. « Non documenté » signifie que le projet n'en parle pas, pas que la fonctionnalité manque.

| Outil | Approche | Plusieurs comptes à la fois | Utilisation sur 5 heures et hebdomadaire | Quand une limite est atteinte | Plateformes |
| --- | --- | --- | --- | --- | --- |
| ccam | Dossiers séparés | Oui | Oui, `ccam usage` | Vous exécutez `ccam next --auto --continue` pour poursuivre la conversation dans le compte qui a le plus de marge | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Échange la connexion ; dossiers séparés avec `cswap run` | Oui, avec `cswap run` | Oui | `cswap auto` change de connexion à un seuil ; la session en cours continue | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Dossiers séparés | Oui | En partie, dans une app de barre des menus macOS | Des groupes de comptes transmettent les requêtes au compte suivant via un proxy local | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Échange la connexion ; dossiers séparés avec `clauth start` | Oui, avec `clauth start` | Oui | Passe au compte suivant d'une chaîne de secours | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Échange la connexion, ou profils isolés | Oui, avec des profils isolés | Oui | `caam run` relance les commandes non interactives sur un autre compte | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Échange la connexion et enveloppe `claude` | Non documenté | Oui | Échange la connexion et reprend la même conversation | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Échange la connexion | Non | Oui | Non documenté | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Échange la connexion | Non | Oui | Changement facultatif à 100 % | macOS |

ccswap est un fork de claude-swap qui gère aussi les connexions Codex. Plusieurs de ces outils prennent aussi en charge d'autres agents de code, comme Codex ou Gemini.

## Les priorités de ccam

- Chaque compte reste connecté. Changer de compte ne déconnecte jamais l'un d'eux, et `ccam work` et `ccam personal` peuvent tourner dans deux terminaux à la fois.
- Les serveurs MCP, `settings.json`, `CLAUDE.md`, les agents, les skills et les plugins sont partagés entre dossiers par des liens, élément par élément et compte par compte. Les connexions ne sont jamais partagées.
- Avec l'historique partagé, une seule commande, `ccam next --auto --continue`, ouvre la même conversation dans le dossier du compte qui a le plus de marge.
- ccam lance Claude Code et se remplace par lui : rien ne tourne en arrière-plan. Il ne surveille pas votre utilisation et ne change pas de compte de lui-même ; c'est vous qui décidez quand.

## Ce que les autres outils font de plus

- claude-swap, CCS, clauth, CAAM et cux fonctionnent sous Windows sans WSL. ccam, non.
- claude-swap, clauth et cux changent de compte d'eux-mêmes quand un compte dépasse un seuil d'utilisation. ccam vous laisse décider.
- CCSwitcher, Claude Switcher et claude-swap affichent l'utilisation dans la barre des menus de macOS.
- CCS, CAAM, ccswap et Claude Switcher gèrent aussi Codex et d'autres agents de code.

## Lequel choisir ?

Choisissez un outil qui échange la connexion si vous utilisez un compte à la fois, si vous voulez que le changement se fasse automatiquement ou si vous êtes sous Windows sans WSL. Choisissez ccam si vous voulez des comptes côte à côte, une seule configuration partagée et la même conversation dans le compte suivant. Si vous n'avez besoin de rien partager, de [simples alias shell](/compare/aliases) suffisent.

Le droit d'utiliser plusieurs comptes dépend des conditions de chaque compte. Vérifiez-les avant de vous lancer.
