---
title: "Référence des commandes"
description: "Toutes les commandes et options de ccam, de ccam add et ccam <name> à ccam usage, ccam next et ccam doctor."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

Placez `--dry-run` avant n'importe quelle commande pour voir ce qu'elle ferait sans l'exécuter, par exemple `ccam --dry-run work`.

<!-- commands -->

## Lancer Claude Code

`ccam <name|number> [arguments]` lance Claude Code avec ce compte. ccam ne lit pas les arguments qui suivent le nom ; il les transmet à `claude` dans le même ordre, et l'entrée standard, le code de sortie et les signaux se comportent exactement comme avec `claude` lui-même. Un `--` après le nom est facultatif.

`ccam` avec des arguments mais sans nom de compte, comme `ccam -p "summarise"`, utilise le compte par défaut. Définissez-le avec `ccam use <name>`.

Les raccourcis `claude1`, `claude2`, etc. sont des liens vers ccam dans `~/.local/bin`. `ccam shims install --names` crée aussi `claude-<name>` pour chaque compte.

## Ajouter et supprimer des comptes

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

Sans dossier, `ccam import` cherche `~/.claude` et `~/.claude-*` et pose la question pour chacun. Les dossiers importés sont enregistrés là où ils se trouvent ; ccam ne déplace ni ne copie jamais une connexion.

## Partage

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Chaque élément accepte `on` (le groupe `default`), `off` ou un nom de groupe. Voir [Partage](/docs/sharing).

## Serveurs MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

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
ccam shims install
```

## Fichiers

| Chemin | Contenu |
| --- | --- |
| `~/.config/ccam/config.toml` | Les comptes et ce qu'ils partagent |
| `~/.local/share/ccam/profiles/<name>/` | Le `CLAUDE_CONFIG_DIR` de chaque compte |
| `~/.local/share/ccam/groups/<group>/` | Paramètres, serveurs MCP, extensions et historique partagés |
| `~/.local/state/ccam/` | Cache d'utilisation et heures de dernière utilisation |
| `~/.local/bin/claude1`, `claude2`… | Liens vers ccam |

Les variables `XDG_CONFIG_HOME`, `XDG_DATA_HOME` et `XDG_STATE_HOME` changent l'emplacement de ces dossiers.
