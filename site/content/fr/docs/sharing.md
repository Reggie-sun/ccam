---
title: "Partager les serveurs MCP, les paramètres et l'historique"
description: "Comment ccam partage serveurs MCP, settings.json, CLAUDE.md, agents, skills et historique entre comptes Claude Code, chaque connexion restant séparée."
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

Chaque compte conserve sa propre connexion. Tout le reste peut être partagé élément par élément : ccam relie le fichier ou le dossier du compte à un dossier de groupe partagé, si bien qu'une modification faite dans un compte apparaît dans tous les comptes de ce groupe.

## Ce qui peut être partagé

| Élément | Option | Ce qui est partagé |
| --- | --- | --- |
| Serveurs MCP | `--mcp` | Le `mcp.json` du groupe, transmis à Claude Code au démarrage |
| Paramètres | `--settings` | `settings.json` : permissions, hooks, modèle, ligne d'état |
| Mémoire globale | `--memory` | `CLAUDE.md` |
| Extensions | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Historique des conversations | `--history` | `projects/` et `history.jsonl` |

Les nouveaux comptes partagent tout sauf l'historique des conversations.

Les fichiers de connexion `.credentials.json` et `.claude.json` ne sont jamais partagés. ccam ne relie que des chemins figurant dans une liste fixe, et ces deux fichiers n'y figurent pas.

## Groupes

`on` désigne le groupe `default`. Vous pouvez aussi indiquer un nom, pour que seuls certains comptes partagent un élément :

```sh
ccam share work --history team
ccam share side --history team
```

`work` et `side` partagent désormais leurs conversations entre eux, mais pas avec vos autres comptes.

## Activer le partage

Quand un compte commence à partager un élément, ccam :

1. affiche un plan de chaque fichier qu'il va déplacer, copier ou relier, et demande confirmation avant de commencer ;
2. déplace les fichiers du compte dans le dossier du groupe si le groupe est encore vide, ou ajoute les fichiers que le groupe n'a pas encore ;
3. quand un même fichier diffère, demande s'il faut garder la version du groupe, utiliser la vôtre ou fusionner le JSON ;
4. conserve votre original dans `.ccam-backup/<time>/` à l'intérieur du dossier du compte ;
5. remplace l'original par un lien vers le groupe.

Les fichiers de conversation portent le nom de leur ID de session : les historiques de plusieurs comptes se combinent donc sans conflit.

## Désactiver le partage

```sh
ccam share work --history off
```

Le lien est supprimé et le compte reçoit sa propre copie du contenu actuel du groupe. Ajoutez `--fresh` pour repartir d'un contenu vide à la place. Le groupe conserve son contenu pour les autres comptes.

## Serveurs MCP

Claude Code stocke les serveurs MCP dans `.claude.json`, avec la connexion : ce fichier ne peut donc pas être partagé. ccam conserve les serveurs partagés dans le `mcp.json` du groupe et lance Claude Code avec `--mcp-config=<group>/mcp.json`. Vos propres options `--mcp-config` s'appliquent toujours par-dessus.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Un serveur que vous ajoutez pendant une session avec `claude mcp add` est enregistré pour ce compte uniquement. Au prochain démarrage de ce compte, ccam indique le nom du serveur et la commande qui le déplace vers le groupe :

```sh
ccam mcp import work gh
```

Les serveurs MCP qui se connectent avec OAuth stockent leurs jetons avec la connexion du compte : autorisez-les donc une fois dans chaque compte.

Si vous préférez que les serveurs apparaissent dans `claude mcp list` et `/mcp` comme des serveurs utilisateur ordinaires, définissez `mcp_strategy = "sync"` dans `~/.config/ccam/config.toml`. ccam les écrit alors dans le `.claude.json` de chaque compte avant le démarrage, sous le verrou de fichier propre à Claude Code et avec une sauvegarde.

## Deux comptes dans le même dossier

Les comptes qui partagent l'historique des conversations écrivent dans les mêmes fichiers. Si vous lancez un second Claude Code dans un dossier où un autre compte du même groupe est déjà en cours d'exécution, ccam vous avertit et vous demande si vous voulez continuer.

## En cas de problème

```sh
ccam doctor
ccam doctor --fix
```

Si un chemin partagé est devenu un fichier ordinaire, par exemple après que vous avez restauré une sauvegarde à la main, ccam s'arrête avant de lancer Claude Code et vous renvoie vers `ccam doctor --fix`, qui fusionne le fichier dans le groupe et conserve une sauvegarde.
