---
title: "CLAUDE_CONFIG_DIR expliqué"
description: "Rôle de CLAUDE_CONFIG_DIR dans Claude Code, contenu du dossier, emplacement de la connexion sous macOS et fichiers partageables entre comptes."
updated: 2026-10-23
order: 2
translation: machine
sourceHash: 2cf08478085f
---

`CLAUDE_CONFIG_DIR` indique à Claude Code quel dossier utiliser à la place de `~/.claude`. Tout ce qui est propre au compte s'y trouve : la connexion, les paramètres, les serveurs MCP, les plugins et l'historique des conversations. Le faire pointer vers des dossiers différents vous donne des comptes séparés.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR est-il officiellement pris en charge ?

Oui. Claude Code le cite parmi ses [variables d'environnement](https://code.claude.com/docs/en/env-vars) comme moyen de remplacer le dossier de configuration, « utile pour exécuter plusieurs comptes côte à côte », et son [guide d'authentification](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) l'utilise pour se connecter avec plusieurs comptes. Vous pouvez le définir dans votre shell ou dans les paramètres utilisateur, mais pas dans les paramètres d'un projet.

## Que contient le dossier ?

Vérifié avec Claude Code 2.1.309 :

| Chemin | Contenu | Partageable |
| --- | --- | --- |
| `.credentials.json` | Jetons de connexion (Linux, Windows), jetons OAuth MCP | Non |
| `.claude.json` | Informations du compte, serveurs MCP utilisateur, état par projet | Non (serveurs MCP à part) |
| `settings.json` | Permissions, hooks, modèle, ligne d'état, plugins | Oui |
| `CLAUDE.md` | Instructions globales | Oui |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensions | Oui |
| `projects/` | Conversations sous forme de fichiers `.jsonl`, mémoire par projet | Oui |
| `history.jsonl` | Historique des prompts | Oui |
| `sessions/`, `shell-snapshots/`, `cache/` et autres | État d'exécution | Non |

## Quels sont les pièges ?

**Sans `CLAUDE_CONFIG_DIR`, `.claude.json` ne se trouve pas dans `~/.claude`.** Claude Code le conserve alors dans `~/.claude.json`, à côté du dossier. Si vous commencez à utiliser `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code cherche `~/.claude/.claude.json` à la place et ne trouve pas votre compte. `ccam import` gère ce cas en reliant les deux. Si `~/.claude` a déjà son propre `.claude.json`, issu d'une exécution antérieure avec `CLAUDE_CONFIG_DIR=~/.claude`, ccam laisse les deux fichiers et le signale, car vous seul savez lequel est le bon.

**Sous macOS, la connexion se trouve dans le trousseau, pas dans le dossier.** Claude Code la stocke sous le nom `Claude Code-credentials`, suivi de `-` et des 8 premiers chiffres hexadécimaux du SHA-256 du chemin du dossier quand `CLAUDE_CONFIG_DIR` est défini, pour le compte nommé par `$USER`. Des dossiers distincts ont donc des entrées distinctes dans le trousseau, et un dossier déplacé ou renommé perd sa connexion ; c'est pourquoi `ccam rename` conserve le nom du dossier sous macOS. Si le trousseau refuse l'écriture, Claude Code se rabat sur `.credentials.json` dans le dossier, comme le décrit [sa documentation](https://code.claude.com/docs/en/authentication#credential-management).

**Toutes les connexions ne sont pas stockées dans le dossier.** Une connexion Claude Console sans clé API est stockée en dehors : des dossiers séparés ne permettent donc pas d'en garder deux distinctes. Les connexions claude.ai et les clés API sont conservées par dossier.

## Où se trouve le dossier sous Windows ?

Par défaut dans `%USERPROFILE%\.claude`, avec la connexion dans le fichier `.credentials.json` qu'il contient, ou dans le Gestionnaire d'identification de Windows quand Claude Code est configuré pour l'utiliser, sous un nom construit à partir du chemin du dossier comme sous macOS. Définissez `CLAUDE_CONFIG_DIR` sur un autre dossier pour chaque compte supplémentaire, comme sous Linux et macOS.

## Les dossiers peuvent-ils partager des fichiers ?

Les liens symboliques fonctionnent : Claude Code lit et écrit à travers un `settings.json` ou un dossier `projects/` lié, et conserve le lien. Les serveurs MCP font exception, car ils partagent `.claude.json` avec la connexion ; transmettez-les plutôt avec `--mcp-config`.

ccam repose exactement sur ces règles. Voir [Partage](/docs/sharing) pour savoir comment il relie les dossiers, et [le guide sur plusieurs comptes](/guides/multiple-claude-code-accounts) pour une vue d'ensemble.
