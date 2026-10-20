---
title: "Référence des commandes"
metaTitle: "Référence des commandes ccam"
description: "Toutes les commandes et options de ccam, de ccam add et ccam <name> à ccam usage, ccam next et ccam doctor."
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 116c5b34ac9a
---

Placez `--dry-run` avant une commande pour voir ce qu'elle ferait sans l'exécuter, par exemple `ccam --dry-run work`. Cela fonctionne au démarrage de Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) et avec `group`, `migrate`, `shortcut`, `update` et `doctor --fix`. Les autres commandes qui modifient quelque chose la refusent au lieu de s'exécuter.

<!-- commands -->

## Lancer Claude Code

`ccam <name|number> [arguments]` lance Claude Code avec ce compte. ccam ne lit pas les arguments qui suivent le nom ; il les transmet à `claude` dans le même ordre, et l'entrée standard, le code de sortie et les signaux se comportent exactement comme avec `claude` lui-même. Un `--` après le nom est facultatif.

`ccam` avec des arguments mais sans nom de compte, comme `ccam -p "summarise"`, utilise le compte défini pour le dossier courant, sinon le compte par défaut. Définissez-les avec `ccam folder <name> [folder]` et `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` avant le nom du compte s'applique aux exécutions `-p`/`--print`. Quand le compte atteint sa limite d'utilisation, ccam relance la même commande avec le compte connecté qui a le plus de marge et n'affiche que la sortie de l'exécution qui a abouti. Voir [Comptes par dossier](/features/folder) et [Relancer à la limite](/features/next#run-again-at-the-limit).

Un raccourci comme `ccwork` lance son compte de la même façon : `ccwork --resume` fait la même chose que `ccam work --resume`. Voir [Raccourcis](/docs/commands#shortcuts).

## Ajouter et supprimer des comptes

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

Sans dossier, `ccam import` cherche `~/.claude` et `~/.claude-*` et pose la question pour chacun. Les dossiers importés sont enregistrés là où ils se trouvent ; ccam ne déplace ni ne copie jamais une connexion.

`ccam add --api` ajoute un compte qui utilise une clé d'API au lieu d'une connexion Claude ; ccam lit la clé à chaque lancement et ne la stocke jamais. Voir [Comptes à clé d'API](/features/api).

`ccam rename`, `ccam remove` et les commandes `ccam group` refusent de modifier un compte tant que Claude Code s'exécute avec ce compte. Quittez d'abord la session.

`ccam list --json` affiche chaque compte avec son numéro d'emplacement, son raccourci, son adresse e-mail et son dossier. Son tableau `groups` liste les groupes du compte, et son objet `share` associe chaque élément (`mcp`, `settings`, `memory`, `extensions`, `history`) au groupe dont il provient, ou à `local`.

## Raccourcis

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Un raccourci est une commande que vous nommez vous-même et qui lance Claude Code avec un compte donné. Comme `ccam work`, il transmet tous les arguments tels quels à Claude Code, si bien que `ccwork --resume` et `ccwork -p "hi"` fonctionnent comme prévu. C'est un lien vers ccam dans `~/.local/bin`.

Les raccourcis sont facultatifs. L'assistant de `ccam add` demande s'il faut en créer un, `ccam add <name> --shortcut <command>` en crée un sans rien demander, et `ccam list` affiche le raccourci de chaque compte, ou `-` s'il n'en a pas.

Les noms utilisent les caractères a-z, 0-9, `.`, `-` et `_`. Un nom ne peut pas être une commande ccam, le raccourci d'un autre compte, ni le nom d'un autre programme présent dans votre `PATH`.

`ccam rename` conserve le nom du raccourci, et `ccam remove` supprime le raccourci avec le compte. `ccam shims install` recrée les liens de tous les comptes, par exemple après avoir déplacé ccam.

## Partage

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Un groupe partage des éléments entre ses comptes ; un compte peut appartenir à plusieurs groupes, mais chaque élément ne provient que d'un seul groupe. Voir [Partager avec des groupes](/features/share).

## Copier une configuration existante

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copie `settings.json`, `CLAUDE.md`, les extensions et les serveurs MCP utilisateur d'un dossier Claude Code dans un groupe, pour que chaque compte qui partage ces éléments en dispose. Le dossier est seulement lu. Les options `--settings`, `--memory`, `--extensions`, `--mcp` et `--history` ne copient que les éléments indiqués. Sans elles, une liste à cocher vous permet de choisir dans un terminal, avec tout coché sauf l'historique des conversations ; ailleurs, tout est copié sauf l'historique des conversations. `-g` ou `--group` choisit le groupe, et `-y` ou `--yes` copie sans rien demander et garde la copie du groupe quand un fichier diffère. Il copie aussi les connexions enregistrées des serveurs dans le fichier de connexion propre à chaque compte. Voir [Copier une configuration existante](/docs/sharing#copy-an-existing-setup).

## Serveurs MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` déplace dans son groupe les serveurs que seul ce compte possède. Un serveur dont le groupe utilise déjà le nom pour une autre définition est ignoré ; ajoutez `--force` pour remplacer la copie du groupe, que ccam sauvegarde d'abord.

## Utilisation et changement de compte

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` et `ccam pick` choisissent de la même façon : d'abord sous `auto_max_usage` (`config.toml`, 100 par défaut), selon `auto_strategy` (`headroom`, `expiring` ou `balance`), en comptant la fenêtre hebdomadaire de `--model`, jamais un compte à 100 %, et au sein de l'organisation du compte de départ (`--cross-org` ou `cross_org = true` pour chercher plus loin). Quand aucun ne peut être choisi, ils quittent avec le code 75 ou, avec `--wait`, attendent la réinitialisation du premier compte. `ccam next --continue` transporte une conversation même entre des comptes qui ne partagent pas l'historique : Claude Code la copie alors dans le nouveau compte. Dans Claude Code, la commande `/ccam:switch [account]` du plugin passe une conversation à un autre compte, et `ccam hook` est ce que lancent ses hooks. Voir [Passer au compte suivant](/features/next), [Choisir un compte pour les scripts](/features/pick) et [Les sessions en arrière-plan de tous les comptes](/features/agents).

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` vérifie la somme de contrôle SHA-256 de la version avant de remplacer ccam, et `--version v0.2.0` installe une version précise. Un ccam installé depuis le paquet `.deb` se met à jour avec apt. Voir [Mettre à jour](/docs/install#update).

## Fichiers

| Chemin | Contenu |
| --- | --- |
| `~/.config/ccam/config.toml` | Les comptes, les groupes et les réglages (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | Le `CLAUDE_CONFIG_DIR` de chaque compte |
| `~/.local/share/ccam/groups/<group>/` | Paramètres, serveurs MCP, extensions et historique partagés |
| `~/.local/state/ccam/` | Cache d'utilisation, heures de dernière utilisation et les notes que laissent `/ccam:switch` et le hook de limite (`handoff/`) |
| `~/.local/bin/<shortcut>` | Raccourcis comme `ccwork`, des liens vers ccam |

Les variables `XDG_CONFIG_HOME`, `XDG_DATA_HOME` et `XDG_STATE_HOME` changent l'emplacement de ces dossiers. Sous Windows, `~` correspond à `%USERPROFILE%` et les raccourcis sont des `<shortcut>.exe`.
