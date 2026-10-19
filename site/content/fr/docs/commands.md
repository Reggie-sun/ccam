---
title: "Référence des commandes"
description: "Toutes les commandes et options de ccam, de ccam add et ccam <name> à ccam usage, ccam next et ccam doctor."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5a82175c30d3
---

Placez `--dry-run` avant une commande pour voir ce qu'elle ferait sans l'exécuter, par exemple `ccam --dry-run work`. Cela fonctionne au démarrage de Claude Code (`ccam <name>`, `ccam run`, `ccam next`) et avec `share`, `migrate`, `shortcut`, `update` et `doctor --fix`. Les autres commandes qui modifient quelque chose la refusent au lieu de s'exécuter.

<!-- commands -->

## Lancer Claude Code

`ccam <name|number> [arguments]` lance Claude Code avec ce compte. ccam ne lit pas les arguments qui suivent le nom ; il les transmet à `claude` dans le même ordre, et l'entrée standard, le code de sortie et les signaux se comportent exactement comme avec `claude` lui-même. Un `--` après le nom est facultatif.

`ccam` avec des arguments mais sans nom de compte, comme `ccam -p "summarise"`, utilise le compte par défaut. Définissez-le avec `ccam use <name>`.

Un raccourci comme `ccwork` lance son compte de la même façon : `ccwork --resume` fait la même chose que `ccam work --resume`. Voir [Raccourcis](/docs/commands#shortcuts).

## Ajouter et supprimer des comptes

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

Sans dossier, `ccam import` cherche `~/.claude` et `~/.claude-*` et pose la question pour chacun. Les dossiers importés sont enregistrés là où ils se trouvent ; ccam ne déplace ni ne copie jamais une connexion.

`ccam rename`, `ccam remove` et `ccam share` refusent de modifier un compte tant que Claude Code s'exécute avec ce compte. Quittez d'abord la session.

`ccam list --json` affiche chaque compte avec son numéro d'emplacement, son raccourci, son adresse e-mail et son dossier. Son objet `share` associe chaque élément (`mcp`, `settings`, `memory`, `extensions`, `history`) à un nom de groupe ou à `local`.

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
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Chaque élément accepte `on` (le groupe `default`), `off` ou un nom de groupe. Voir [Partage](/docs/sharing).

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
ccam next --continue
ccam next --auto
```

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
```

`ccam update` vérifie la somme de contrôle SHA-256 de la version avant de remplacer ccam, et `--version v0.2.0` installe une version précise. Un ccam installé depuis le paquet `.deb` se met à jour avec apt. Voir [Mettre à jour](/docs/install#update).

## Fichiers

| Chemin | Contenu |
| --- | --- |
| `~/.config/ccam/config.toml` | Les comptes et ce qu'ils partagent |
| `~/.local/share/ccam/profiles/<name>/` | Le `CLAUDE_CONFIG_DIR` de chaque compte |
| `~/.local/share/ccam/groups/<group>/` | Paramètres, serveurs MCP, extensions et historique partagés |
| `~/.local/state/ccam/` | Cache d'utilisation et heures de dernière utilisation |
| `~/.local/bin/<shortcut>` | Raccourcis comme `ccwork`, des liens vers ccam |

Les variables `XDG_CONFIG_HOME`, `XDG_DATA_HOME` et `XDG_STATE_HOME` changent l'emplacement de ces dossiers.
