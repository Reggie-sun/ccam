---
title: "Copier une configuration existante"
description: "Copiez les paramètres, CLAUDE.md, les extensions et les serveurs MCP d'un ~/.claude existant vers un groupe ccam avec ccam migrate."
updated: 2026-10-23
order: 6
translation: machine
sourceHash: 7bca9efb6eb0
---

Si vous utilisiez Claude Code avant ccam, votre configuration se trouve dans `~/.claude` et `~/.claude.json`. `ccam migrate` la copie dans un groupe de partage, si bien que chaque compte qui partage ces éléments par ce groupe démarre avec la configuration que vous avez déjà. Le dossier source est seulement lu, et `claude` seul continue de fonctionner comme avant. Lancez-le une fois après votre premier `ccam add` ; quand `~/.claude.json` contient des serveurs MCP, le premier `ccam add` les nomme dans son conseil.

## Utilisation

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

Dans un terminal, sans option d'élément et sans `--yes`, ccam affiche d'abord une liste à cocher :

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Espace (ou `x`) coche ou décoche l'élément sous le curseur, les flèches ou `j`/`k` déplacent le curseur, Entrée affiche le plan, et Échap, `q` ou Ctrl+C affiche `Cancelled.`. Les éléments pour lesquels la source n'a rien sont marqués `not found`. Si vous décochez tout, ccam affiche `Nothing selected.` et s'arrête.

ccam affiche ensuite le plan et demande confirmation avant de copier :

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Un groupe n'atteint que ses propres comptes, et seulement pour les éléments qu'il partage ; le résumé liste donc les comptes du groupe. Quand le groupe ne partage pas l'un des éléments copiés, ou que certains comptes n'en font pas partie, ccam affiche les commandes `ccam group set <group> --add <items>` et `ccam group join <group> <accounts>` qui y remédient.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `[folder]` | Le dossier Claude Code à copier | `~/.claude` |
| `--settings` | Copier `settings.json` | désactivé |
| `--memory` | Copier le `CLAUDE.md` global | désactivé |
| `--extensions` | Copier `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | désactivé |
| `--mcp` | Copier les serveurs MCP utilisateur (`mcpServers` dans `.claude.json`) | désactivé |
| `--history` | Copier l'historique des conversations (`projects/`, `history.jsonl`) | désactivé |
| `-g`, `--group string` | Le groupe de partage dans lequel copier | `default` |
| `-y`, `--yes` | Copier sans demander ; en cas de conflit, la copie du groupe est conservée | désactivé |
| `--dry-run` (globale) | Afficher le plan et s'arrêter avant toute question ou copie | désactivé |

Sans option d'élément, la liste à cocher décide dans un terminal. Hors terminal, ou avec `--yes`, ccam copie tout sauf l'historique des conversations. L'historique est exclu par défaut car `projects/` peut être volumineux.

`--group on` désigne `default`. `--group off` et `--group local` sont refusés avec `--group takes a group name, not "off"`. `ccam migrate` respecte `--dry-run`, et la liste à cocher apparaît quand même lors d'une simulation, pour que vous puissiez prévisualiser une sélection.

## Fonctionnement

`internal/cli/migrate.go` choisit les éléments, puis `migrate.Build` lit la source et le groupe et renvoie un plan. Rien n'est écrit avant `Plan.Apply`, après votre confirmation.

1. **Trouver la source.** `checkSource` refuse un chemin absent, qui n'est pas un dossier ou qui ne contient aucune entrée Claude Code connue. Pour les seuls serveurs MCP, un `~/.claude` vide ou absent convient tant que `~/.claude.json` contient des serveurs. `migrate.ClaudeJSON` trouve le `.claude.json` qui appartient au dossier : à l'intérieur, ou, pour un `~/.claude` qui n'en contient pas, `~/.claude.json` à côté. C'est là que Claude Code le conserve quand `CLAUDE_CONFIG_DIR` n'est pas défini.
2. **Associer les éléments aux chemins.** Chaque chemin vient de `share.Layout`, la même liste d'autorisation que celle qu'utilisent les groupes. Les serveurs MCP sont lus dans `.claude.json` avec `mcp.LocalServers`, qui utilise le scanner `profile.TopLevelValue` au lieu d'analyser tout le fichier.
3. **Copier, jamais déplacer.** Les fichiers vont dans `~/.local/share/ccam/groups/<group>/` avec `fsutil.AtomicWrite` ou `fsutil.CopyFile`, qui écrit un fichier temporaire à côté de la cible et le renomme, si bien qu'une copie interrompue ne laisse jamais un fichier tronqué. Les serveurs MCP vont dans le `mcp.json` du groupe sous son verrou : un serveur qu'un autre ccam ajoute entre-temps reste, et les copies faites pour [Remote Control](/guides/remote-control) sont mises à jour ensuite.

Chaque élément associe un chemin source à un chemin du groupe :

<!-- motion:migrate-items -->

Les deux lignes de l'historique sont en pointillés car elles ne sont copiées qu'avec `--history`. `.credentials.json` et `.claude.json` figurent dans `share.Never` et n'atteignent jamais le groupe : `Build` renvoie une erreur interne si l'un d'eux apparaît un jour dans la structure, et le parcours des dossiers les ignore par leur nom.

`buildEntry` attribue un état à chaque chemin, et c'est ce que décrit chaque ligne du plan :

<!-- motion:migrate-states -->

- **Missing** : absent de la source, ou un fichier là où un dossier est attendu. Listé une fois par élément sous la forme `not found, skipped`.
- **Linked** : le chemin source pointe déjà vers la copie du groupe, par exemple quand la source est un dossier de profil ccam.
- **Same** : le groupe a déjà un contenu identique.
- **New** : le groupe ne l'a pas, il est donc copié en entier.
- **Merge** : un dossier présent des deux côtés. `copyMissing` n'ajoute que les fichiers absents du groupe ; les fichiers qui diffèrent restent tels que le groupe les a et le plan indique `group copies kept`. Pour `plugins/`, les registres `installed_plugins.json` et `known_marketplaces.json` sont fusionnés, et les chemins de plugins absolus sont réécrits pour pointer vers le dossier du groupe.
- **Conflict** : un fichier présent des deux côtés avec un contenu différent. Vous choisissez fichier par fichier : `keep group`, `use mine`, plus `merge JSON` pour les fichiers `.json` (une fusion profonde où le groupe l'emporte sur les valeurs en conflit) ou `merge by line` pour `history.jsonl`. Un serveur MCP dont la définition diffère est aussi un conflit, avec `keep group` ou `use mine`. Le premier choix listé est celui par défaut quand vous appuyez sur Entrée ; `--yes` conserve toujours la copie du groupe.

La liste à cocher elle-même est `tui.MigratePicker`, construite sur le composant `checklist` qu'utilise `ccam add` ; `migrate.Found` marque les éléments absents à l'aide d'appels lstat et de la liste des serveurs MCP.

### Connexions MCP

Les serveurs MCP qui se connectent avec OAuth conservent leurs jetons dans le fichier de connexion, `.credentials.json` sous `mcpOAuth`, et non dans la liste des serveurs. ccam ne les copie pas. Une connexion copiée contient le même jeton de renouvellement que l'original, et le serveur fait tourner ce jeton au premier renouvellement : tous les autres détenteurs, `~/.claude` source compris, sont alors déconnectés de ce serveur. La ligne `sign-ins` du plan nomme donc, pour chaque serveur copié auquel la source est connectée, les comptes du groupe qui se connectent eux-mêmes avec `/mcp`, et le résumé le répète. Une entrée sans jeton n'est pas une connexion et n'est pas listée, pas plus qu'un compte qui a déjà sa propre connexion, ni un serveur dont le groupe garde la définition. Aucun fichier de connexion n'est écrit.

## Sécurité et cas particuliers

- La source est seulement lue. Aucun de ses fichiers n'est déplacé, relié ni modifié, et la connexion de la source n'est jamais écrite, y compris sous macOS où elle peut être lue depuis le trousseau.
- Avant qu'un fichier du groupe ne soit remplacé, par `use mine`, une fusion ou de nouveaux serveurs MCP dans un `mcp.json` existant, la copie du groupe est enregistrée dans `<group>/.ccam-backup/<time>/`. Le nom du dossier est unique pour chaque exécution, si bien que deux exécutions dans la même seconde n'écrasent pas leurs sauvegardes respectives. Le résumé affiche le chemin.
- Avec `--yes`, chaque conflit conserve la copie du groupe, et ccam indique comment choisir à la place : relancer la même commande sans `--yes`.
- Si rien ne changerait, ccam affiche `Nothing to copy.` ; si chaque conflit conserve la copie du groupe, `Nothing to copy: the group keeps its own copies.`
- Si une copie échoue à mi-parcours, ccam s'arrête avec `The source folder was not changed. Group files replaced so far are saved in` suivi du dossier de sauvegarde. Une nouvelle exécution copie les fichiers qui ne sont pas arrivés en entier.
- Quand plusieurs groupes partagent MCP, les commandes que suggèrent `ccam doctor`, `ccam ext` et `ccam mcp list` s'écrivent `ccam migrate --mcp --group <a|b>`.
- Une source absente se termine par `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Voir aussi

- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Vérifier et réparer](/features/doctor)
- [Plusieurs comptes Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
