---
title: "Partager serveurs MCP, paramètres et historique"
description: "Comment ccam partage serveurs MCP, settings.json, CLAUDE.md, agents, skills et historique entre comptes Claude Code, chaque connexion restant séparée."
updated: 2026-10-19
order: 4
translation: machine
sourceHash: f266f8bfd55c
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

ccam ne modifie pas le partage tant que Claude Code s'exécute avec le compte, car la session continue d'écrire dans les fichiers en cours de déplacement. Quittez-la d'abord.

## Désactiver le partage

```sh
ccam share work --history off
```

Le lien est supprimé et le compte reçoit sa propre copie du contenu actuel du groupe. Ajoutez `--fresh` pour repartir d'un contenu vide à la place. Le groupe conserve son contenu pour les autres comptes.

## Copier une configuration existante

Si vous utilisiez Claude Code avant ccam, vos paramètres se trouvent dans `~/.claude`. `ccam migrate` les copie dans un groupe, pour que chaque compte qui partage ces éléments en dispose :

```sh
ccam migrate
```

Dans un terminal, ccam affiche d'abord une liste des éléments à cocher. Tout est coché sauf l'historique des conversations ; appuyez sur Espace pour changer un élément et sur Entrée pour voir le plan. Les éléments pour lesquels le dossier ne contient rien sont marqués `not found`.

Ce qui est copié :

- `settings.json` et le `CLAUDE.md` global ;
- les extensions `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` et `rules/` ;
- vos serveurs MCP utilisateur, c'est-à-dire l'entrée `mcpServers` de `~/.claude.json`, dans le `mcp.json` du groupe ;
- les connexions enregistrées (OAuth) des serveurs MCP qu'il copie, dans le fichier de connexion propre à chaque compte ;
- l'historique des conversations (`projects/` et `history.jsonl`), mais seulement si vous ajoutez `--history`.

Le dossier source est seulement lu : rien n'y est déplacé, lié ni modifié, et `claude` seul continue de fonctionner comme avant. Les fichiers de connexion `.credentials.json` et `.claude.json` ne sont jamais copiés dans le groupe.

Les serveurs MCP qui se connectent avec OAuth gardent leurs jetons dans la connexion, pas dans la liste des serveurs. Pour vous éviter de vous y reconnecter, ccam copie ces connexions dans le fichier de connexion de chaque compte qui reçoit les serveurs par le groupe, sauf si le compte a déjà les siennes. Il ignore les comptes en cours d'exécution, pas encore connectés, enregistrés sur place avec `ccam import` ou qui gardent leur connexion dans le trousseau macOS, et les nomme ; connectez-vous aux serveurs depuis ces comptes avec `/mcp`. Un compte qui commence à partager les serveurs MCP plus tard reçoit les connexions quand vous relancez `ccam migrate`.

ccam affiche le plan et demande confirmation avant de copier quoi que ce soit :

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

À la fin, ccam liste les comptes qui partagent chaque élément via le groupe. Si certains comptes ne le font pas, il affiche les commandes `ccam share` qui activent le partage pour eux.

Quand le groupe possède déjà une copie différente d'un fichier, ccam vous demande, comme lors de l'activation du partage, de garder la copie du groupe ou d'utiliser la vôtre. `settings.json` et les autres fichiers JSON peuvent aussi être fusionnés, et `history.jsonl` ligne par ligne. Les copies précédentes du groupe sont enregistrées dans `.ccam-backup/<time>/` à l'intérieur du dossier du groupe.

| Option | Effet |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copier seulement les éléments indiqués, sans la liste à cocher ; sans elles et hors d'un terminal, tout sauf l'historique des conversations |
| `-g`, `--group <group>` | Copier dans ce groupe au lieu de `default` |
| `-y`, `--yes` | Copier sans rien demander ; quand un fichier diffère, garder la copie du groupe |

Placez `--dry-run` avant la commande, comme dans `ccam --dry-run migrate`, pour voir le plan sans rien copier. Pour copier depuis un autre dossier, indiquez-le : `ccam migrate ~/.claude-work --group work`.

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

Si le groupe possède déjà un autre serveur portant ce nom, ccam l'ignore et vous le signale. `--force` remplace le serveur du groupe et conserve l'ancien `mcp.json` dans le dossier `.ccam-backup/` du groupe.

Les serveurs MCP qui se connectent avec OAuth stockent leurs jetons avec la connexion du compte : autorisez-les donc une fois dans chaque compte, ou copiez les connexions depuis `~/.claude` avec `ccam migrate`.

Si vous préférez que les serveurs apparaissent dans `claude mcp list` et `/mcp` comme des serveurs utilisateur ordinaires, définissez `mcp_strategy = "sync"` dans `~/.config/ccam/config.toml`. ccam les écrit alors dans le `.claude.json` de chaque compte avant le démarrage, sous le verrou de fichier propre à Claude Code et avec une sauvegarde.

## Deux comptes dans le même dossier

Les comptes qui partagent l'historique des conversations écrivent dans les mêmes fichiers. Si vous lancez un second Claude Code dans un dossier où un autre compte du même groupe est déjà en cours d'exécution, ccam vous avertit et vous demande si vous voulez continuer.

## En cas de problème

```sh
ccam doctor
ccam doctor --fix
```

Si un chemin partagé est devenu un fichier ordinaire, par exemple après que vous avez restauré une sauvegarde à la main, ccam s'arrête avant de lancer Claude Code et vous renvoie vers `ccam doctor --fix`, qui fusionne le fichier dans le groupe et conserve une sauvegarde.
