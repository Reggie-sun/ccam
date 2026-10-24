---
title: "Inspecter les extensions"
description: "ccam ext affiche les skills, commandes, plugins, serveurs MCP et hooks que charge chaque compte, l'origine de chacun et ce qui est cassé."
updated: 2026-10-23
order: 14
translation: machine
sourceHash: 45740e4efbef
---

`ccam ext` répond à la question « pourquoi ce skill (ou ce serveur MCP) fonctionne-t-il dans un compte et pas dans l'autre ? ». Il liste, pour chaque compte, les extensions que Claude Code charge et l'origine de chacune : un groupe de partage, le dossier propre au compte ou le projet dans lequel vous êtes. Sous le tableau, il liste ce qui empêche une extension de se charger comme prévu. Il ne fait que lire.

## Utilisation

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

Une cellule indique la source de cet élément dans ce compte : le groupe de partage (`default`), `account` quand il se trouve dans le dossier propre au compte, `local` pour un serveur MCP que le compte a pour ce seul dossier, ou `-` quand le compte ne l'a pas. Un plugin installé mais non activé affiche `(off)` après sa source. `!` signale un élément qui pose problème, expliqué sous le tableau ; les problèmes qui n'appartiennent à aucun compte, comme les serveurs que seul `claude` seul possède, viennent en dernier, sans nom. `ccam ext --json` affiche le même rapport (`accounts`, le projet, `problems`), avec `remote` sur les serveurs MCP que Claude Code joint par le réseau ; sans aucun compte, il affiche le rapport avec `accounts: []`.

## Ce qu'il lit

| Type | Emplacement | Source affichée |
| --- | --- | --- |
| Skills | chaque dossier de `skills/` qui contient un `SKILL.md` | le groupe vers lequel pointe `skills/`, ou `account` |
| Commandes, agents | les fichiers `.md` sous `commands/` et `agents/`, les sous-dossiers sous la forme `folder/name` | le groupe, ou `account` |
| Plugins | `plugins/installed_plugins.json`, activés ou désactivés par `enabledPlugins` dans les paramètres du compte et du projet ; un plugin installé pour un seul projet n'apparaît que dans ce projet | le groupe, ou `account` ; `(off)` quand il n'est pas activé |
| Serveurs MCP | le `mcp.json` du groupe (quand MCP est partagé), le `.claude.json` du compte (ses propres serveurs, et ceux qu'il a pour ce seul dossier), et les serveurs de ses plugins, sous la forme `plugin:<plugin>:<server>` | le groupe, `account` ou `local` |
| Hooks | `hooks` dans `settings.json`, comptés par événement | le groupe, ou `account` |
| Ce projet | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` et `.mcp.json` dans le dossier courant | listé une seule fois ; un compte qui n'a pas approuvé le dossier reçoit un problème pour les serveurs de son `.mcp.json` |

## Problèmes signalés

- **Liens cassés** : un lien dans `skills/`, `commands/` ou `agents/` dont la cible a disparu, ou un dossier partagé (`skills`, `commands`, `agents`, `plugins`) qui pointe vers un emplacement absent (réparez-le avec `ccam doctor --fix`, qui recrée un dossier de groupe absent).
- **Un dossier de skill sans `SKILL.md`** : Claude Code l'ignore.
- **Un serveur MCP défini différemment à deux endroits** : le `mcp.json` du groupe et le `.claude.json` du compte, ou l'un des deux et le `.mcp.json` du projet. Les mêmes paramètres aux deux endroits (ce qu'écrit `mcp_strategy = "sync"`) ne posent pas problème. Déplacez un serveur propre à un compte dans le groupe avec `ccam mcp import <account> <server>`.
- **Un `installed_plugins.json` qui n'est pas du JSON valide**, un plugin que `enabledPlugins` ne nomme pas, et un plugin dont le dossier est absent.
- **Des serveurs qu'un seul compte charge** : un serveur présent seulement dans le `.claude.json` de ce compte, que les autres comptes de son groupe ne chargent pas, et un serveur dans la portée locale du compte pour ce dossier (ce que fait `claude mcp add` sans `-s user`). Chacun est accompagné de la commande `ccam mcp import` qui le partage. Les serveurs du `~/.claude.json` de `claude` seul qu'aucun compte ne charge sont accompagnés de la commande `ccam migrate --mcp`.
- **Un serveur de projet en attente d'approbation** : le compte n'a pas encore approuvé le dossier, donc Claude Code pose la question une fois, et d'ici là les serveurs du `.mcp.json` du projet attendent. Le problème indique où les approuver : `"enabledMcpjsonServers"` dans le `settings.json` des comptes du groupe des paramètres, ou de ce compte. L'approbation donnée au worktree principal d'un dépôt git vaut aussi dans ses autres worktrees.
- **Connexions et connecteurs** : un serveur MCP qui demande une connexion dans ce compte alors qu'un autre compte y est connecté (une connexion ne vaut que pour un seul compte ; connectez-vous avec `/mcp`), et des connecteurs claude.ai qu'un autre compte a utilisés et pas celui-ci.

## Sécurité et cas particuliers

- `ccam ext` ne modifie rien et ne fait aucune requête réseau. Il lit les noms des serveurs et le nombre de hooks, jamais les valeurs de `env` ni des en-têtes dans les définitions MCP. Pour la vérification des connexions, il lit les noms des connexions MCP dans la connexion de chaque compte (sous macOS, depuis le trousseau), jamais les jetons, et seulement pour les comptes qui chargent un serveur qu'un autre compte charge aussi.
- Il montre ce qui se trouve sur le disque. Un serveur bloqué par des paramètres gérés apparaît quand même ; un plugin désactivé dans les paramètres affiche `(off)`.
- Lancez-le dans le dossier du projet pour inclure le dossier `.claude` et le `.mcp.json` de ce projet.

## Voir aussi

- [Partage entre comptes](/features/share)
- [Serveurs MCP partagés](/features/mcp)
- [Vérifier et réparer](/features/doctor)
