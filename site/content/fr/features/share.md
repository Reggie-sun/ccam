---
title: "Partager avec des groupes"
description: "Un groupe ccam partage des éléments choisis entre ses comptes via un dossier ; ccam group déplace les données avec plan, sauvegardes et arbitrage des conflits."
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 3f94dd13dbf1
---

Chaque compte a son propre dossier, mais la plupart de ce qu'il contient n'a pas besoin d'être séparé. Un **groupe** indique ce qui est partagé et par qui : chaque compte du groupe relie ces éléments à un même dossier, `~/.local/share/ccam/groups/<group>/`, si bien qu'ils lisent et écrivent tous la même copie. Les jetons de connexion ne sont jamais partagés.

Un compte peut appartenir à plusieurs groupes, mais chaque élément ne provient que d'un seul groupe. Une configuration typique comprend un groupe pour les paramètres et les extensions, et un second pour l'historique des conversations entre deux comptes :

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Le premier compte que vous ajoutez crée `default`, et les comptes suivants le rejoignent sauf indication contraire.

## Utilisation

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Chaque commande qui déplace des données affiche un plan, pose une question pour chaque fichier qui diffère de la copie du groupe, et demande `Proceed? [y/N]` avant de modifier quoi que ce soit :

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

Vous pouvez faire les mêmes modifications dans l'onglet **Groups** du [tableau de bord](/features/dashboard) ; appuyez sur `s` sur un compte pour choisir ses groupes.

## Éléments

| Élément | Ce qu'il couvre |
| --- | --- |
| `mcp` | les serveurs MCP, via le `mcp.json` du groupe ([Serveurs MCP partagés](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | le `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` et `history.jsonl`, pour qu'un autre compte puisse faire `--continue` |

Si aucun des groupes d'un compte ne partage un élément, cet élément reste dans le dossier propre au compte.

## Commandes et options

| Commande | Effet |
| --- | --- |
| `ccam group` ou `ccam group list` | lister les groupes |
| `ccam group show <group>` | ce que partage le groupe, ses comptes, son dossier et ses serveurs MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | créer un groupe ; les éléments de ses comptes y sont déplacés |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | changer ce que partage le groupe, pour chacun de ses comptes |
| `ccam group join <group> <account>...` | placer des comptes dans le groupe |
| `ccam group leave <group> <account>...` | retirer des comptes ; chacun garde une copie, serveurs MCP compris |
| `ccam group rename <group> <new-name>` | renommer, dans la configuration seulement |
| `ccam group remove <group> [--leave]` | supprimer un groupe |

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--from <account>` | `add`, `set` : le compte dont la copie entre la première dans le groupe | le compte par défaut, puis par numéro |
| `--replace` | `join` : quitter d'abord les groupes qui partagent les mêmes éléments | désactivé |
| `--fresh` | les comptes qui cessent de partager un élément repartent sans lui au lieu d'en garder une copie (un dossier importé en garde quand même une) | désactivé |
| `--leave` | `remove` : retirer d'abord les comptes du groupe | nécessaire quand il a des comptes |
| `--prefer <group\|mine\|merge>` | comment trancher les conflits | demander ; avec `--yes`, le choix par défaut de chaque fichier |
| `-y`, `--yes` | exécuter sans demander ; les conflits suivent `--prefer` | désactivé |
| `--dry-run` (globale) | afficher le plan, ne rien modifier | désactivé |

Les comptes sont indiqués par leur nom ou leur numéro, les éléments par une liste séparée par des virgules. `ccam add` accepte `--group <group>` (répétable) ou `--no-group`.

## Fonctionnement

Les fichiers concernés :

- `~/.config/ccam/config.toml` : les entrées `[[group]]` : `name`, `items`, `members` et, quand le nom du dossier diffère, `dir`.
- `~/.local/share/ccam/profiles/<name>/` : le `CLAUDE_CONFIG_DIR` du compte.
- `~/.local/share/ccam/groups/<folder>/` : les copies partagées.

À la lecture de `config.toml`, les groupes de chaque compte sont résolus en un dossier par élément. Cela prend environ 25 µs pour cinq comptes, si bien qu'un lancement ne lit toujours que `config.toml`. Une configuration où un compte recevrait un élément de deux groupes est une erreur, et rien ne démarre tant qu'elle n'est pas corrigée.

Une modification travaille sur une copie de la configuration :

1. La commande modifie la copie (`AddGroup`, `SetGroupItems`, `JoinGroup` …), qui vérifie les règles.
2. ccam compare l'ancienne et la nouvelle configuration et, pour chaque compte dont les éléments changent, construit un plan (`internal/group`, `share.Build`). Les plans sont construits l'un après l'autre, et chacun voit ce que les précédents vont placer dans les dossiers de groupe : seul le premier compte y déplace donc sa copie. La première copie du groupe provient, dans cet ordre : du dossier du groupe (s'il en a déjà une), du compte `--from`, du compte par défaut, puis des autres comptes par numéro.
3. ccam refuse si l'un de ces comptes est en cours d'exécution, quelle que soit la façon dont Claude Code a été lancé (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), et, sous Windows sans le Mode développeur, si un fichier devait être relié par un lien.
4. ccam tranche les conflits (en posant la question, ou avec `--prefer`), affiche le plan et demande confirmation.
5. Juste avant d'enregistrer, ccam vérifie de nouveau. Si un compte a démarré pendant que ccam attendait vos réponses, la modification s'arrête ; elle s'arrête aussi si un autre processus ccam a modifié les mêmes groupes entre-temps (`the groups changed while this plan waited; nothing was moved, make the change again`). La nouvelle configuration est enregistrée avant tout déplacement, et la modification est notée dans `~/.local/state/ccam/group-change.json`. Si une étape s'arrête à mi-parcours, la note reste : `ccam doctor --fix` termine la modification dans son ordre, avec les réponses que vous avez données, étapes MCP comprises, et d'ici là aucun autre changement de groupe ne démarre (``the last group change stopped halfway; finish it with `ccam doctor --fix` first``).
6. Les étapes MCP s'exécutent en premier. Quand des comptes retirent les extensions d'un groupe, les listes de plugins du groupe sont pointées vers le dossier du groupe partout où leurs chemins passent par le dossier d'un compte (`group plugin lists:` dans le plan). Ensuite, les plans s'exécutent dans l'ordre. Chaque étape déplace des données vers le groupe ou vers le dossier de sauvegarde ; seuls les fichiers identiques à la copie du groupe sont supprimés au lieu d'être sauvegardés.

ccam tient une liste d'autorisation des chemins couverts par chaque élément (`share.Layout`). Rien de ce qui n'y figure ne peut atteindre un groupe. ccam désigne aussi `.credentials.json` et `.claude.json` comme jamais partagés (`share.Never`), et s'arrête avec une erreur interne si l'un d'eux apparaît un jour dans la structure (`Build`). L'élément `mcp` n'a pas de chemins : la liste des serveurs se trouve dans `.claude.json`, donc le [partage MCP](/features/mcp) passe par le `mcp.json` du groupe.

Voici ce qui arrive à un dossier, `projects/`, quand `side` commence à partager l'historique et que le groupe a déjà du contenu :

<!-- motion:share-on -->

Si le groupe n'a encore rien à ce chemin, le compte est le premier à le partager et son dossier y est déplacé tel quel. Sinon, chaque fichier qui manque au groupe y est déplacé, et les fichiers présents des deux côtés restent dans le dossier du compte, si bien que la version du groupe est conservée. Les fichiers identiques à ceux du groupe sont supprimés, et le plan nomme ceux qui diffèrent (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

Le dossier du compte va ensuite dans `.ccam-backup/<time>/` à l'intérieur du dossier du compte, et un lien prend sa place. Si l'entrée du groupe est elle-même un lien vers l'entrée propre du compte qui rejoint, le contenu du compte va dans le groupe à sa place. Les fichiers de conversation portent le nom de leur ID de session : les historiques de plusieurs comptes se combinent donc sans conflit. Les registres de `plugins/` enregistrent des chemins absolus qui passent par le dossier du compte ; ccam les fait pointer vers le dossier du groupe et fusionne les entrées du compte dans un registre que le groupe possède déjà.

Quitter le partage inverse le lien :

<!-- motion:share-off -->

ccam copie d'abord le contenu actuel du groupe dans `.ccam-backup/incomplete-<name>` dans le dossier du compte, et ne remplace le lien qu'une fois la copie complète. Si la copie s'interrompt, le lien reste en place et `ccam doctor --fix` copie de nouveau (le plan d'une telle reprise affiche `restore <entry>: … (copied back from the group before a stop)`). Le groupe conserve son contenu pour les autres comptes.

- ccam ne supprime le lien que s'il pointe toujours dans le dossier des groupes comme au moment où le plan a été créé. Un chemin qui a changé depuis, comme un vrai fichier à la place du lien, est laissé tel quel et la modification s'arrête.
- Avec `--fresh`, la copie est sautée.
- Un compte qui fait passer un élément à un groupe dont le dossier est encore vide apporte une copie de ce qu'il voyait, et garde donc ses fichiers. Cette copie est préparée dans `<group>/.ccam-incomplete-<name>` jusqu'à ce qu'elle soit complète.
- Un compte qui cesse de partager MCP garde les serveurs du groupe : ils vont dans son propre `.claude.json` (une sauvegarde est conservée), sauf quand le compte a déjà un serveur du même nom, qui reste tel quel. Avec `--fresh`, le plan liste plutôt les serveurs qui cessent de se charger. Rejoindre de nouveau retire les copies que personne n'a modifiées.
- Quitter `extensions` fait aussi pointer les listes de plugins du groupe vers le dossier du groupe là où leurs chemins passaient par le dossier du compte qui part ; pour les autres comptes, les chemins réécrits mènent aux mêmes plugins.

Un fichier présent des deux côtés avec un contenu différent est un conflit :

<!-- motion:share-conflict -->

Vous choisissez pour les fichiers isolés comme `settings.json` ; dans les dossiers, la copie du groupe reste, et le plan nomme les fichiers qui diffèrent. Les choix dépendent du fichier ; le premier est le choix par défaut :

| Fichier | Choix |
| --- | --- |
| `history.jsonl` | `merge by line`, `keep group`, `use mine` |
| Autres fichiers `.json` | `keep group`, `use mine`, `merge JSON` |
| Tout le reste | `keep group`, `use mine` |

`merge JSON` fusionne les objets récursivement, joint les tableaux sans doublons et garde la valeur du groupe quand les deux côtés définissent la même clé différemment, si bien que les autres comptes ne voient aucun changement. Chaque choix sauf `keep group` enregistre d'abord la copie du groupe dans `.ccam-backup/<time>/group/` du compte. Votre propre fichier va toujours dans la sauvegarde.

À chaque lancement, ccam compare chaque chemin de la structure avec les groupes du compte en n'utilisant que `lstat` et `readlink` (`share.Check`). Un lien manquant est créé sur-le-champ. Un vrai fichier là où devrait se trouver un lien, un lien vers le mauvais endroit ou un lien périmé arrête le lancement et renvoie vers `ccam doctor --fix`.

### Renommer, supprimer et l'ancien format

- `ccam group rename` ne modifie que la configuration. Le dossier garde son nom et le groupe l'enregistre comme `dir` : rien ne bouge et rien ne peut être perdu.
- `ccam group remove` refuse un groupe qui a des comptes, sauf si vous ajoutez `--leave`. Chaque compte garde alors une copie, et le dossier part dans `~/.local/share/ccam/trash/`, sauf si un autre groupe l'utilise.
- Avant les groupes, `config.toml` enregistrait une table `[profile.share]` par compte. ccam lit ce format comme des groupes sans rien déplacer. Au premier enregistrement, il écrit le nouveau format et garde l'ancien fichier sous le nom `config.toml.bak`. Les éléments partagés par les mêmes comptes deviennent un seul groupe portant le nom du dossier ; les autres reçoivent des noms comme `default-history` et continuent d'utiliser le même dossier.

### Sous Windows

Windows ne crée de liens symboliques qu'en Mode développeur (Paramètres > Système > Pour les développeurs) ou pour un administrateur. Sans ce droit, un dossier partagé devient une jonction de répertoire, qui fonctionne de la même façon pour Claude Code. Un fichier partagé (`settings.json`, `CLAUDE.md`, `history.jsonl`) ne peut être qu'un lien symbolique : une modification de groupe qui en exige un est donc refusée avant tout déplacement. Pour la même raison, le premier `ccam add` y crée le groupe `default` avec seulement `mcp` et `extensions` et le signale, et `ccam add --group` vers un groupe qui partage les paramètres, la mémoire ou l'historique est refusé.

## Sécurité et cas particuliers

- **Comptes en cours d'exécution.** Une modification est refusée pendant que Claude Code s'exécute avec l'un des comptes qu'elle déplace, car la session continue d'écrire dans ces fichiers, que ccam l'ait lancée ou non (`claude` seul, un IDE, les sessions en arrière-plan). La vérification est répétée juste avant l'enregistrement. Quitter un groupe ou en rejoindre un n'affecte que le compte nommé ; quitter `extensions` réécrit aussi des chemins dans les listes de plugins du groupe, sans changer ce qu'elles signifient pour les autres. `--dry-run` affiche tout de même le plan.
- **Sauvegardes.** Chaque plan reçoit son propre dossier `.ccam-backup/<time>/` dans le dossier du compte. ccam ne les supprime jamais.
- **Vos propres liens.** Si un chemin est déjà un lien que vous avez créé, par exemple vers un dépôt de dotfiles, et que le groupe est vide, ccam copie la cible du lien dans le groupe et déplace le lien lui-même dans la sauvegarde. La cible n'est pas modifiée.
- **Dossiers importés.** Un compte enregistré sur place avec `ccam import` (comme `~/.claude`) est modifié sur place, si bien que `claude` seul voit aussi le changement ; le plan l'indique. Quand il quitte un groupe, il garde toujours une copie, même avec `--fresh`, et `ccam remove` le retire ainsi de ses groupes avant d'oublier le compte. Ses fichiers partagés se trouvent dans le dossier du groupe tant qu'il fait partie d'un groupe : supprimez donc le compte (ou quittez les groupes) avant de supprimer `~/.local/share/ccam`.
- **Deux comptes dans un même dossier.** Les comptes qui partagent l'historique écrivent dans les mêmes fichiers. Quand un second compte démarre dans un dossier où un autre s'exécute déjà, ccam avertit et demande confirmation.
- **Modifications manuelles.** Si `members` liste un nom qui n'est pas un compte, ccam l'ignore et `ccam doctor` le signale. Les noms `local`, `on`, `off` et `none` ne peuvent pas être des noms de groupe.

## Voir aussi

- [Le tableau de bord](/features/dashboard)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Serveurs MCP partagés](/features/mcp)
- [Copier une configuration existante](/features/migrate)
- [Référence des commandes](/docs/commands)
