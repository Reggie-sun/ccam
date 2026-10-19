---
title: "Questions fréquentes"
description: "Réponses sur ccam, les comptes Claude Code, les limites d'utilisation, le partage, la confidentialité et les systèmes pris en charge."
updated: 2026-10-19
order: 6
translation: machine
sourceHash: e07d468d453b
---

## ccam est-il un outil officiel d'Anthropic ?

Non. ccam est un projet open source non officiel : il n'est ni affilié à Anthropic ni approuvé par Anthropic. Il utilise le paramètre `CLAUDE_CONFIG_DIR`, que Claude Code documente pour utiliser plusieurs comptes côte à côte.

## Ai-je le droit d'utiliser plusieurs comptes Claude Code ?

Cela dépend des conditions de chaque compte et, pour les comptes professionnels, de la politique de votre organisation. Vérifiez-les avant d'utiliser plusieurs comptes. ccam se contente de faciliter le passage d'un compte à l'autre parmi ceux que vous possédez déjà.

## ccam permet-il de contourner les limites d'utilisation ?

Non. Chaque compte conserve ses propres limites. ccam les affiche et vous permet de continuer sur un autre compte qui vous appartient quand l'un d'eux est épuisé.

## Comment utiliser deux comptes Claude Code sur un seul ordinateur ?

Donnez à chaque compte son propre dossier de configuration avec `CLAUDE_CONFIG_DIR`. ccam le fait pour vous : exécutez `ccam add` pour chaque compte, puis lancez-les avec `ccam personal` et `ccam work`, ou avec des raccourcis que vous nommez vous-même, par exemple `ccwork`. Le [guide](/guides/multiple-claude-code-accounts) montre aussi comment procéder à la main.

## Comment poursuivre une conversation sur un autre compte ?

Partagez l'historique des conversations entre les deux comptes avec `ccam share <name> --history on`, puis exécutez `ccam next --continue`. Claude Code ouvre la conversation la plus récente du dossier courant dans le compte suivant.

## ccam voit-il ou envoie-t-il mes identifiants de connexion ?

ccam ne copie ni ne partage jamais les identifiants de connexion. Il lit le jeton d'accès d'un compte uniquement pour demander l'utilisation de ce compte au point de terminaison d'utilisation d'Anthropic, quand vous exécutez `ccam usage`, `ccam next --auto` ou `ccam doctor`. ccam n'a aucune télémétrie et n'envoie rien nulle part ailleurs.

## Où mes comptes sont-ils stockés ?

Dans `~/.local/share/ccam/profiles/<name>/`, un dossier par compte. Sous macOS, Claude Code conserve la connexion elle-même dans le trousseau, sous un nom dérivé de ce dossier.

## Qu'advient-il de mon dossier ~/.claude existant ?

Rien, sauf si vous exécutez `ccam import` ou `ccam migrate`. L'import enregistre le dossier là où il se trouve, si bien que `claude` seul et `ccam <name>` continuent d'utiliser le même compte et les mêmes paramètres. `ccam migrate` copie ses paramètres, `CLAUDE.md`, ses extensions et ses serveurs MCP dans les paramètres partagés, et ne fait que lire le dossier.

## Quels systèmes ccam prend-il en charge ?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine et d'autres) et macOS. Sous Windows, utilisez-le dans WSL.

## Pourquoi un serveur MCP n'apparaît-il que dans un seul compte ?

Les serveurs ajoutés pendant une session avec `claude mcp add` sont enregistrés pour ce compte uniquement. ccam vous le signale au prochain démarrage du compte ; déplacez le serveur vers le groupe partagé avec `ccam mcp import <account> <server>`.

## ccam ralentit-il Claude Code ?

Lancer Claude Code via ccam ajoute environ 4 millisecondes, mesurées avec `scripts/perf-check.sh` sur un ordinateur portable. ccam lit un petit fichier, vérifie les liens partagés puis se remplace par `claude` : il ne tourne donc pas pendant que vous travaillez.

## Comment désinstaller ccam ?

Supprimez le raccourci de chaque compte avec `ccam shortcut <account> --remove`, supprimez `ccam` de `~/.local/bin`, puis `~/.config/ccam`, `~/.local/share/ccam` et `~/.local/state/ccam` si vous n'avez plus besoin des comptes. Voir [Installer ccam](/docs/install#uninstall).
