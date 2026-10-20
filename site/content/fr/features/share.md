---
title: "Partage entre comptes"
description: "Comment ccam share relie paramètres, CLAUDE.md, extensions et historique à un dossier de groupe, avec plan et sauvegardes, et comment off les recopie."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 16de29bbdd57
---

Chaque compte a son propre dossier, mais la plupart de ce qu'il contient n'a pas besoin d'être séparé. `ccam share` active ou désactive le partage pour un compte et un élément à la fois : le fichier ou le dossier du compte est déplacé dans un dossier de groupe et remplacé par un lien symbolique, si bien que chaque compte du groupe lit et écrit la même copie. Utilisez-le après `ccam add` pour changer ce que partage un compte, ou pour placer quelques comptes dans un groupe à part.

## Utilisation

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Sans option d'élément, `ccam share` affiche l'état actuel :

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Avec des options, ccam affiche un plan, pose une question pour chaque fichier qui diffère de la copie du groupe, et demande `Proceed? [y/N]` avant de modifier quoi que ce soit. Un nouveau compte partage tout sauf l'historique des conversations ; [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing) décrit chaque élément du point de vue de l'utilisateur.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Serveurs MCP | inchangé |
| `--settings <on\|off\|group>` | `settings.json` | inchangé |
| `--memory <on\|off\|group>` | le `CLAUDE.md` global | inchangé |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | inchangé |
| `--history <on\|off\|group>` | `projects/` et `history.jsonl` | inchangé |
| `--fresh` | en désactivant le partage, partir de rien au lieu de copier le contenu du groupe | désactivé |
| `--prefer <group\|mine\|merge>` | comment trancher les conflits | demander ; avec `--yes`, le choix par défaut de chaque élément |
| `-y`, `--yes` | exécuter sans demander ; les conflits suivent `--prefer` | désactivé |
| `-h`, `--help` | aide de `share` | |
| `--dry-run` (globale) | respecté : pose les questions sur les conflits, affiche le plan, ne modifie rien | désactivé |

La valeur `on` désigne le groupe `default`, `off` signifie local (non partagé), et toute autre valeur est un nom de groupe. Les noms de groupe deviennent des noms de dossiers sous `~/.local/share/ccam/groups/`, donc `ccam share side --history ../x` échoue avec `invalid group name: "../x"`. Le compte est indiqué par son nom ou son numéro, et seuls les éléments que vous nommez changent.

## Fonctionnement

Les fichiers concernés :

- `~/.config/ccam/config.toml` : la table `[profile.share]` de chaque compte, une valeur par élément.
- `~/.local/share/ccam/profiles/<name>/` : le `CLAUDE_CONFIG_DIR` du compte.
- `~/.local/share/ccam/groups/<group>/` : les copies partagées.

1. `internal/cli/share.go` lit la configuration, applique chaque valeur `--<item>` à une copie des paramètres du compte avec `config.Share.Set`, et refuse si le compte est en cours d'exécution.
2. `share.Build` compare les paramètres actuels avec la cible et liste les étapes. Il ne fait que lire le disque : `os.Lstat` et `os.Readlink` pour chaque chemin de `share.Layout`. Le plan est une liste d'opérations : déplacer dans un emplacement de groupe vide (adoption), déplacer les fichiers absents du groupe (fusion de dossier), écrire un conflit résolu dans le groupe, sauvegarder, supprimer un lien, créer un lien, recopier, réécrire les chemins de plugins et fusionner le registre des plugins.
3. Les conflits sont tranchés, par la question posée ou par `--prefer`.
4. `Plan.Describe` affiche le plan, et ccam demande confirmation.
5. Les paramètres cibles sont écrits dans `config.toml` avant tout déplacement. Si `Plan.Apply` s'arrête à mi-parcours, `ccam doctor --fix` et une nouvelle exécution construisent leur plan vers la cible enregistrée et terminent la modification.
6. `Plan.Apply` exécute les étapes dans l'ordre. Chaque étape déplace des données vers le groupe ou vers le dossier de sauvegarde ; rien n'est supprimé.

`share.Layout` est une liste d'autorisation des chemins couverts par chaque élément. Rien de ce qui n'y figure ne peut atteindre un groupe. `share.Never` nomme `.credentials.json` et `.claude.json`, et `Build` s'arrête avec une erreur interne si l'un d'eux apparaît un jour dans la structure. L'élément `mcp` n'a pas de chemins : la liste des serveurs se trouve dans `.claude.json`, donc le [partage MCP](/features/mcp) passe plutôt par le `mcp.json` du groupe, et `--mcp` ne modifie que le paramètre.

Voici ce qui arrive à un dossier, `projects/`, quand `side` active l'historique et que le groupe a déjà du contenu :

<!-- motion:share-on -->

Si le groupe n'a encore rien à ce chemin, le compte est le premier à le partager et son dossier y est déplacé tel quel, sans sauvegarde. Sinon, `mergeDir` déplace dans le groupe chaque fichier qui lui manque, et les fichiers présents des deux côtés restent dans le dossier du compte, si bien que la version du groupe est conservée. Le dossier du compte va ensuite dans `.ccam-backup/<time>/` à l'intérieur du dossier du compte, et `os.Symlink` met un lien à sa place. Les fichiers de conversation portent le nom de leur ID de session : les historiques de plusieurs comptes se combinent donc sans conflit.

Le déplacement utilise `rename(2)` et se rabat sur une copie suivie d'une suppression entre systèmes de fichiers différents. Un lien relatif dans l'arborescence déplacée qui pointait à l'extérieur est réécrit en lien absolu (`fsutil.PortableLink`), afin qu'il soit toujours résolu depuis le nouvel emplacement. Les registres de `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) enregistrent des chemins absolus qui passent par le dossier du compte ; ccam les fait pointer vers le dossier du groupe et, quand le groupe a déjà un registre, y fusionne les entrées du compte.

Désactiver un élément inverse le lien :

<!-- motion:share-off -->

ccam ne supprime le lien que s'il pointe dans le dossier des groupes, puis copie le contenu actuel du groupe dans le compte avec `CopyTree`. Avec `--fresh`, la copie est sautée et le compte part de rien à ce chemin. Le groupe conserve son contenu pour les autres comptes. Faire passer un élément d'un groupe à un autre supprime le lien et le relie au nouveau groupe, en laissant l'ancien groupe tel quel ; le plan l'indique dans une ligne `Note:`.

Un fichier présent des deux côtés avec un contenu différent est un conflit :

<!-- motion:share-conflict -->

`share.NewConflict` liste les choix adaptés au fichier. `history.jsonl` propose `merge by line`, `keep group` et `use mine`, avec la fusion en premier. Les autres fichiers `.json` proposent `keep group`, `use mine` et `merge JSON`. Tout le reste propose `keep group` et `use mine`. Le premier choix est celui par défaut :

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) fusionne les objets récursivement, joint les tableaux sans doublons et garde la valeur du groupe quand les deux côtés définissent la même clé différemment, si bien que les autres comptes ne voient aucun changement de comportement. `merge by line` garde les lignes du groupe, puis ajoute celles des vôtres qui lui manquent. Chaque choix sauf `keep group` enregistre d'abord la copie du groupe dans `.ccam-backup/<time>/group/` du dossier du compte, puis écrit le résultat avec `fsutil.AtomicWrite`. Votre propre fichier va toujours dans la sauvegarde.

`--prefer group`, `mine` ou `merge` tranche chaque conflit sans demander ; `merge` signifie `merge JSON` ou `merge by line`, selon ce que propose le fichier, et un fichier qui ne propose ni l'un ni l'autre garde le choix par défaut. `--yes` seul prend le choix par défaut de chaque fichier.

À chaque lancement, `share.Check` compare chaque chemin de la structure avec le paramètre en n'utilisant que `lstat` et `readlink`, ce qui le rend assez léger pour le chemin de lancement. Un lien manquant est créé sur-le-champ, car il n'y a rien à perdre. Un vrai fichier là où devrait se trouver un lien, un lien vers le mauvais endroit ou un lien périmé arrête le lancement :

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` exécute `share.Build` avec les paramètres enregistrés des deux côtés, ce qui déduit chaque étape de ce qui se trouve sur le disque, garde le contenu du groupe en cas de conflit et sauvegarde tout ce qu'il remplace.

### Sous Windows

Windows ne crée de liens symboliques qu'en Mode développeur (Paramètres > Système > Pour les développeurs) ou pour un administrateur. `fsutil.Link` essaie d'abord un lien symbolique. Sans ce droit, un dossier partagé (`agents`, `skills`, `projects` et les autres) devient une jonction de répertoire, que tout utilisateur peut créer, et fonctionne de la même façon pour Claude Code. Un fichier partagé (`settings.json`, `CLAUDE.md`, `history.jsonl`) ne peut être qu'un lien symbolique : ccam s'arrête donc avec `Windows links files only in Developer Mode` et propose d'activer ce mode ou de garder cet élément local. `share.Check` lit les jonctions comme des liens (`fsutil.LinkTarget`), et supprimer un dossier de compte supprime la jonction elle-même, jamais le contenu du groupe vers lequel elle pointe.

## Sécurité et cas particuliers

- **Comptes en cours d'exécution.** ccam refuse de modifier le partage pendant que Claude Code s'exécute avec le compte (`account work is running; quit it first, then change its sharing`), car la session continue d'écrire dans les fichiers en cours de déplacement et les lignes qu'elle ajoute entre-temps seraient perdues. `--dry-run` affiche tout de même le plan.
- **Sauvegardes.** Chaque plan reçoit son propre dossier `.ccam-backup/<time>/` dans le dossier du compte, avec un suffixe numéroté quand deux plans démarrent dans la même seconde. ccam ne les supprime jamais.
- **Échecs partiels.** Si une étape échoue, ccam affiche l'erreur et l'emplacement des originaux, et vous indique de terminer avec `ccam doctor --fix`. La cible est déjà enregistrée, donc la réparation continue vers elle.
- **Vos propres liens.** Si un chemin est déjà un lien que vous avez créé, par exemple vers un dépôt de dotfiles, et que le groupe est vide, ccam copie la cible du lien dans le groupe et déplace le lien lui-même dans la sauvegarde. La cible n'est pas modifiée.
- **Fichiers de connexion.** `.credentials.json` et `.claude.json` ne sont jamais reliés, déplacés ni copiés dans un groupe.
- **Deux comptes dans un même dossier.** Les comptes qui partagent l'historique écrivent dans les mêmes fichiers. Quand un second compte du même groupe d'historique démarre dans un dossier où un autre s'exécute déjà, ccam avertit et demande confirmation avant de continuer.

## Voir aussi

- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Serveurs MCP partagés](/features/mcp)
- [Copier une configuration existante](/features/migrate)
- [Toutes les commandes](/docs/commands)
