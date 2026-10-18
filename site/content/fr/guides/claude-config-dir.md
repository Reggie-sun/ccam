---
title: "CLAUDE_CONFIG_DIR expliqué"
description: "Rôle de CLAUDE_CONFIG_DIR dans Claude Code, contenu du dossier, emplacement de la connexion sous macOS et fichiers partageables entre comptes."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` indique à Claude Code quel dossier utiliser à la place de `~/.claude`. Tout ce qui est propre au compte s'y trouve : la connexion, les paramètres, les serveurs MCP, les plugins et l'historique des conversations. Le faire pointer vers des dossiers différents vous donne des comptes séparés.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Contenu du dossier

Vérifié avec Claude Code 2.1.302 :

| Chemin | Contenu | Partageable |
| --- | --- | --- |
| `.credentials.json` | Jetons de connexion (Linux), jetons OAuth MCP | Non |
| `.claude.json` | Informations du compte, serveurs MCP utilisateur, état par projet | Non (serveurs MCP à part) |
| `settings.json` | Permissions, hooks, modèle, ligne d'état, plugins | Oui |
| `CLAUDE.md` | Instructions globales | Oui |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensions | Oui |
| `projects/` | Conversations sous forme de fichiers `.jsonl`, mémoire par projet | Oui |
| `history.jsonl` | Historique des prompts | Oui |
| `sessions/`, `shell-snapshots/`, `cache/` et autres | État d'exécution | Non |

## Deux détails qui piègent souvent

**Sans `CLAUDE_CONFIG_DIR`, `.claude.json` ne se trouve pas dans `~/.claude`.** Claude Code le conserve alors dans `~/.claude.json`, à côté du dossier. Si vous commencez à utiliser `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code cherche `~/.claude/.claude.json` à la place et ne trouve pas votre compte. `ccam import` gère ce cas en reliant les deux.

**Sous macOS, la connexion se trouve dans le trousseau, pas dans le dossier.** Claude Code la stocke sous le nom `Claude Code-credentials`, suivi d'un court hachage du chemin du dossier quand `CLAUDE_CONFIG_DIR` est défini. Des dossiers distincts ont donc des entrées distinctes dans le trousseau. Si le trousseau refuse l'écriture, Claude Code se rabat sur `.credentials.json` dans le dossier.

## Partager entre dossiers

Les liens symboliques fonctionnent : Claude Code lit et écrit à travers un `settings.json` ou un dossier `projects/` lié, et conserve le lien. Les serveurs MCP font exception, car ils partagent `.claude.json` avec la connexion ; transmettez-les plutôt avec `--mcp-config`.

ccam repose exactement sur ces règles. Voir [Partage](/docs/sharing) pour savoir comment il relie les dossiers, et [le guide sur plusieurs comptes](/guides/multiple-claude-code-accounts) pour une vue d'ensemble.
