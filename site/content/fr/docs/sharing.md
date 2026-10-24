---
title: "Partager serveurs MCP, paramètres et historique"
description: "Comment les groupes ccam partagent serveurs MCP, settings.json, CLAUDE.md, skills et historique entre comptes Claude Code, chaque connexion restant séparée."
updated: 2026-10-23
order: 4
translation: machine
sourceHash: 27e0a0e5abbd
---

Chaque compte conserve sa propre connexion. Tout le reste peut être partagé par des **groupes** : un groupe nomme les éléments qu'il partage et les comptes qu'il contient, et ccam relie les fichiers de chacun de ces comptes au dossier du groupe, si bien qu'une modification faite dans un compte apparaît dans tous.

## Ce qui peut être partagé

| Élément | Ce qui est partagé |
| --- | --- |
| `mcp` | les serveurs MCP : le `mcp.json` du groupe, transmis à Claude Code au démarrage |
| `settings` | `settings.json` : permissions, hooks, modèle, ligne d'état |
| `memory` | le `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | l'historique des conversations : `projects/` et `history.jsonl` |

Votre premier compte crée le groupe `default`, qui partage tout sauf l'historique des conversations, et les comptes suivants le rejoignent.

Les fichiers de connexion `.credentials.json` et `.claude.json` ne sont jamais partagés. ccam ne relie que des chemins figurant dans une liste fixe, et ces deux fichiers n'y figurent pas.

Cette page couvre les tâches courantes. Toutes les options sont décrites dans les pages des fonctionnalités : [partager avec des groupes](/features/share), [serveurs MCP partagés](/features/mcp) et [copie d'une configuration existante](/features/migrate).

## Groupes

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Un compte peut appartenir à plusieurs groupes, mais chaque élément ne provient que d'un seul groupe. Pour partager l'historique entre `work` et `side` mais pas avec vos autres comptes, laissez donc l'historique hors de `default` et donnez à ces deux comptes un groupe à eux, comme ci-dessus. `ccam share work` montre ce que partage un compte et par quel groupe. Le [tableau de bord](/features/dashboard) fait de même dans son onglet Groups.

## Quand un compte commence à partager

Quand un compte commence à partager un élément (il rejoint un groupe, ou son groupe ajoute l'élément), ccam :

1. affiche un plan de chaque fichier qu'il va déplacer, copier ou relier pour chaque compte, et demande confirmation avant de commencer ;
2. déplace les fichiers du compte dans le dossier du groupe si le groupe est encore vide, ou ajoute les fichiers que le groupe n'a pas encore ;
3. quand un même fichier diffère, demande s'il faut garder la version du groupe, utiliser la vôtre ou fusionner le JSON ;
4. conserve votre original dans `.ccam-backup/<time>/` à l'intérieur du dossier du compte ;
5. remplace l'original par un lien vers le groupe.

Les fichiers de conversation portent le nom de leur ID de session : les historiques de plusieurs comptes se combinent donc sans conflit.

ccam ne modifie pas un compte tant que Claude Code s'exécute avec lui, quelle que soit la façon dont Claude Code a été lancé (sessions en arrière-plan comprises), car la session continue d'écrire dans les fichiers en cours de déplacement. Quittez-la d'abord. ccam vérifie de nouveau juste avant d'enregistrer : une session lancée pendant que le plan attendait est donc détectée elle aussi.

## Quand un compte cesse de partager

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Le lien est supprimé et le compte reçoit sa propre copie du contenu actuel du groupe ; les serveurs MCP du groupe vont dans son propre `.claude.json`, à côté de ses propres serveurs. Ajoutez `--fresh` pour repartir d'un contenu vide à la place ; le plan liste alors les serveurs MCP qui cessent de se charger. Un dossier importé comme `~/.claude` garde toujours une copie. Le groupe conserve son contenu pour les autres comptes.

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
- l'historique des conversations (`projects/` et `history.jsonl`), mais seulement si vous ajoutez `--history`.

Le dossier source est seulement lu : rien n'y est déplacé, lié ni modifié, et `claude` seul continue de fonctionner comme avant. Les fichiers de connexion `.credentials.json` et `.claude.json` ne sont jamais copiés dans le groupe.

Les serveurs MCP qui se connectent avec OAuth gardent leurs jetons dans la connexion, pas dans la liste des serveurs. ccam ne copie pas ces connexions : une connexion copiée contient le même jeton d'actualisation, et le premier renouvellement déconnecte de ce serveur tous les autres détenteurs, `~/.claude` compris. Le plan et le résumé nomment les comptes qui se connectent eux-mêmes aux serveurs avec `/mcp`.

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
Group default: personal
To share them:
  ccam group join default work
```

À la fin, ccam liste les comptes du groupe. Si le groupe ne partage pas encore l'un des éléments, ou que certains comptes n'en font pas partie, il affiche les commandes `ccam group` qui y remédient.

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

Un serveur que vous ajoutez pendant une session avec `claude mcp add` est enregistré pour ce compte uniquement, et pour le dossier où vous vous trouviez (la portée locale de Claude Code ; `-s user` le rend propre au compte dans tous les dossiers). Au prochain démarrage de ce compte, ccam indique le nom du serveur, avec son dossier, et la commande qui le déplace vers le groupe :

```sh
ccam mcp import work gh
```

Si le groupe possède déjà un autre serveur portant ce nom, ccam l'ignore et vous le signale. `--force` remplace le serveur du groupe et conserve l'ancien `mcp.json` dans le dossier `.ccam-backup/` du groupe.

Les serveurs MCP qui se connectent avec OAuth stockent leurs jetons avec la connexion du compte : autorisez-les donc une fois dans chaque compte avec `/mcp`. `ccam ext` et `ccam doctor` nomment les serveurs auxquels un compte n'est pas connecté alors qu'un autre compte l'est.

Si vous préférez que les serveurs apparaissent dans `claude mcp list` et `/mcp` comme des serveurs utilisateur ordinaires, définissez `mcp_strategy = "sync"` en haut de `~/.config/ccam/config.toml`, au-dessus du premier `[[group]]`, `[[profile]]` ou `[[folder]]`. ccam les écrit alors dans le `.claude.json` de chaque compte avant le démarrage, sous le verrou de fichier propre à Claude Code et avec une sauvegarde. Un serveur que le compte possède lui-même sous le même nom reste celui du compte, et ccam le signale.

## Deux comptes dans le même dossier

Les comptes qui partagent l'historique des conversations écrivent dans les mêmes fichiers. Si vous lancez un second Claude Code dans un dossier où un autre compte du même groupe est déjà en cours d'exécution, ccam vous avertit et vous demande si vous voulez continuer.

## En cas de problème

```sh
ccam doctor
ccam doctor --fix
```

Si un chemin partagé est devenu un fichier ordinaire, par exemple après que vous avez restauré une sauvegarde à la main, ccam s'arrête avant de lancer Claude Code et vous renvoie vers `ccam doctor --fix`, qui le relie de nouveau au groupe : le fichier devient la copie du groupe quand le groupe n'en a pas, et sinon le groupe garde la sienne et la vôtre part dans une sauvegarde, que `--fix` nomme. Un changement de groupes arrêté à mi-chemin est lui aussi terminé par `ccam doctor --fix`.
