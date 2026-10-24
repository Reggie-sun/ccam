---
title: "Les sessions en arrière-plan de tous les comptes"
description: "ccam agents liste en un tableau les agents en arrière-plan et les sessions ouvertes de Claude Code de tous vos comptes ; ccam attach en ouvre un par son id."
updated: 2026-10-23
order: 16
translation: machine
sourceHash: "092368f63363"
---

Claude Code conserve les sessions en arrière-plan (`claude --bg`, la vue des agents) compte par compte : `claude agents` ne montre que les sessions du compte avec lequel il s'exécute. `ccam agents` interroge tous les comptes en même temps et les liste ensemble, et `ccam attach <id|name>` ouvre une session en arrière-plan avec le compte qui l'a lancée, quel qu'il soit.

## Utilisation

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` est l'état de la session en arrière-plan (`working`, `blocked`, `done`, `failed`, `stopped`) et ce que fait son processus (`busy`, `waiting`, `idle`) ; une session interactive ouverte affiche `open`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--all` | Lister aussi les sessions en arrière-plan terminées, en échec ou arrêtées | désactivé |
| `--json` | Afficher les sessions en JSON, chacune avec le compte (`account`) auquel elle appartient | désactivé |

`ccam attach <id|name> [claude args]` transmet tout ce qui suit l'id à `claude attach` sans le modifier, options comprises : `ccam attach 49ec --dry-run` passe `--dry-run` à Claude Code. Pour afficher plutôt la commande, placez `--dry-run` avant l'id (`ccam --dry-run attach 49ec` ou `ccam attach --dry-run 49ec`).

## Fonctionnement

1. Pour chaque compte, ccam exécute `claude agents --json` (avec `--all` s'il est indiqué) avec le `CLAUDE_CONFIG_DIR` du compte, pour tous en même temps, chacun avec 15 secondes au plus. Cela lit la liste de sessions propre au compte sans démarrer le service d'arrière-plan de Claude Code.
2. Chaque entrée reçoit le nom du compte, et la liste est triée par heure de démarrage, la plus récente d'abord. Un compte qui ne répond pas est nommé sur la sortie d'erreur standard et laissé de côté, puis ccam quitte avec le code 1, pour qu'un script ne prenne pas un échec pour une liste vide.
3. `ccam attach` cherche `jobs/<id>` dans le dossier de chaque compte, là où Claude Code conserve une session en arrière-plan, et exécute `claude attach <id>` avec ce compte, de la même façon que `ccam <name>` lance Claude Code. Un mot qui n'est le début d'aucun id est recherché dans les noms des sessions, sans tenir compte de la casse ; un mot qui ressemble à un id n'est jamais pris pour un nom.

## Sécurité et cas particuliers

- Les identifiants de session en arrière-plan appartiennent à un seul compte : `claude attach` exécuté avec un autre compte ne peut pas les trouver, et c'est ce que règle `ccam attach`.
- Deux comptes pourraient avoir des sessions dont les identifiants commencent de la même façon, ou dont les noms contiennent le même mot ; `ccam attach` les nomme alors et demande une plus grande partie de l'id ou du nom.
- ccam transmet ses serveurs MCP partagés à Claude Code avec `--mcp-config=<file>` en un seul argument, ce que `claude agents` et `claude attach` acceptent ; les sessions en arrière-plan continuent d'utiliser ce fichier quand elles redémarrent.
- Le service d'arrière-plan d'un compte à clé d'API prend sa clé dans le shell qui l'a lancé. Lancez la première session en arrière-plan d'un tel compte via ccam (`ccam gw --bg …`), pour que le service dispose de la clé.

## Voir aussi

- [Lancer Claude Code](/features/launch)
- [Choisir un compte pour les scripts](/features/pick)
- [Le tableau de bord](/features/dashboard)
