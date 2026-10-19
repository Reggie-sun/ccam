---
title: "Utiliser deux comptes Claude Code sur un seul ordinateur"
description: "Utilisez plusieurs comptes Claude Code côte à côte avec CLAUDE_CONFIG_DIR, à la main avec des alias shell ou automatiquement avec ccam."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2fbee9f4f91
---

Donnez à chaque compte son propre dossier de configuration. Claude Code lit ce dossier dans la variable d'environnement `CLAUDE_CONFIG_DIR`, et chaque dossier contient sa propre connexion, ses paramètres et son historique. Vous pouvez configurer cela à la main avec un alias shell, ou laisser ccam s'en charger.

## À la main, avec un alias

C'est l'approche que décrit la documentation de Claude Code. Ajoutez une ligne par compte supplémentaire dans `~/.bashrc` ou `~/.zshrc` :

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Ouvrez un nouveau terminal et exécutez `claude-work`. Claude Code vous demande de vous connecter, car le dossier est nouveau. À partir de là, `claude` est votre premier compte et `claude-work` le second.

Cela fonctionne bien pour deux comptes. À savoir :

- Chaque dossier part de zéro : pas de paramètres, pas de `CLAUDE.md`, pas de serveurs MCP, pas d'agents ni de skills. Vous copiez ou reliez ce dont vous avez besoin, et vous gardez les copies synchronisées vous-même.
- Les alias n'existent que dans les shells interactifs. Les scripts, les Makefiles et les autres programmes ne les voient pas.
- Les serveurs MCP se trouvent dans `.claude.json`, à côté de la connexion : on ne peut donc pas simplement les relier d'un dossier à l'autre.
- Les conversations restent dans le compte où elles ont commencé : vous ne pouvez donc pas en poursuivre une dans un autre compte.
- Rien ne montre l'utilisation de tous les comptes à la fois.

## Avec ccam

ccam automatise la même idée et ajoute ce qui manque :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Désormais, `claude1` et `claude2`, ou `ccam personal` et `ccam work`, lancent les deux comptes. Comme ce sont de vrais programmes et non des alias, ils fonctionnent aussi dans les scripts. Les paramètres, `CLAUDE.md`, les serveurs MCP et les extensions sont partagés dès le départ, et vous pouvez activer le partage de l'historique des conversations pour poursuivre une conversation sur l'autre compte avec `ccam next --continue`.

Si vous utilisez déjà `~/.claude`, exécutez `ccam import` au lieu d'ajouter de nouveau ce compte.

## Lequel choisir

Pour deux comptes sans rien de partagé, un alias suffit. Si vous voulez partager paramètres et serveurs MCP, avoir plus de comptes, voir l'utilisation d'un coup d'œil ou poursuivre des conversations, ccam vous épargne cette gestion. La [comparaison](/compare/aliases) présente les différences côte à côte.
