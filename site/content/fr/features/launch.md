---
title: "Lancer Claude Code"
description: "Lancez Claude Code avec n'importe quel compte via ccam <name|number>, ccam run ou le sélecteur ccam ; les arguments sont transmis tels quels à claude, et --dry-run affiche la commande."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6c04bb48df22
---

`ccam work` lance Claude Code avec le compte `work`, et tout ce qui suit le nom est transmis tel quel à `claude`. Utilisez-le partout où vous taperiez `claude` : dans un terminal, dans un script ou avec une entrée redirigée.

## Utilisation

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Indiquez le nom d'un compte ou son numéro dans `ccam list`. ccam ne lit pas les arguments qui suivent le nom, donc `ccam work --help` affiche l'aide de Claude Code. Un `--` après le nom est facultatif. L'entrée standard, le code de sortie et les signaux se comportent comme avec `claude`.

Quand le premier argument commence par `-` et n'est pas l'une des options propres à ccam, ccam lance le compte par défaut avec tous les arguments. Le compte par défaut est celui défini avec `ccam use`, ou le compte au plus petit numéro si vous n'en avez pas défini.

### Le sélecteur

`ccam` sans argument ouvre le sélecteur quand l'entrée et la sortie standard sont un terminal, et affiche `ccam list` dans le cas contraire. Il montre le numéro de chaque compte, son adresse e-mail, son utilisation sur 5 heures, sa dernière utilisation et ses éléments partagés ; l'utilisation s'affiche au fur et à mesure que chaque compte répond.

| Touche | Effet |
| --- | --- |
| `Enter` | Lancer Claude Code avec le compte sélectionné |
| `1` à `9` | Lancer le compte portant ce numéro |
| `Up`/`Down`, `k`/`j` | Déplacer la sélection |
| `a` | Fermer le sélecteur et lancer l'assistant `ccam add` |
| `s` | Modifier les éléments partagés du compte sélectionné (`Space` ou `x` bascule, `Enter` applique, `Esc` revient en arrière) |
| `d` | Faire du compte sélectionné le compte par défaut |
| `r` | Récupérer à nouveau l'utilisation, sans le cache court |
| `q`, `Esc`, `Ctrl+C` | Quitter |

Le sélecteur se ferme avant d'agir, si bien que le terminal est restauré avant le démarrage de Claude Code. Après `d` ou une modification du partage, il se rouvre. Une modification du partage affiche d'abord le plan et demande confirmation, comme le fait `ccam share`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--dry-run` | Afficher l'environnement et la commande que ccam exécuterait, puis quitter | désactivé |
| `-v`, `--version` | Afficher `ccam <version>` (uniquement en premier argument) | |
| `-h`, `--help` | L'aide de ccam, sous la forme `ccam --help` ou `ccam run --help` | |

Règles des arguments :

- ccam ne lit ses propres options qu'avant le nom du compte. Après le nom, chaque argument appartient à Claude Code, y compris `--dry-run` et `--help`.
- `--dry-run` fonctionne sous la forme `ccam --dry-run work`, `ccam run --dry-run work` et `ccam --dry-run` (le compte par défaut). Il est respecté, pas refusé.
- `ccam run <name|number> [claude args]` est la forme longue de `ccam <name>`. `ccam run` sans nom utilise le compte par défaut.
- Les noms de compte ne peuvent pas être des noms de commandes ccam comme `add`, `list`, `run` ou `usage`, si bien que `ccam <word>` n'est jamais ambigu.

La simulation affiche une ligne que vous pourriez coller dans un shell :

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Fonctionnement

Chaque lancement suit les mêmes étapes, de la ligne de commande jusqu'à `exec` :

<!-- motion:launch -->

Chaque case s'allume pendant l'exécution de son étape, et les lignes en dessous sont l'environnement et les arguments avec lesquels Claude Code démarre.

1. `cmd/ccam/main.go` appelle `app.Main`. Si ccam a été lancé sous un autre nom, il s'agit d'un [raccourci](/features/shortcuts), qui va directement au lanceur.
2. `app.parse` examine `os.Args` avant toute intervention de cobra. Seuls `ccam` seul (le sélecteur), un nom de commande ccam (`config.Reserved`), `-h`/`--help` et `-v`/`--version` vont à `cli.Execute`. Tout le reste est un lancement.
3. `launcher.Launch` lit `~/.config/ccam/config.toml`, le seul fichier qu'il analyse.
4. `launcher.Prepare` trouve le compte par son nom ou son numéro, retire un éventuel `--` initial et résout l'exécutable `claude` avec `ResolveClaude`. La recherche dans le PATH est mise en cache dans `~/.local/state/ccam/claude-path` avec la date de modification du fichier, si bien que ccam ne parcourt pas le `PATH` à chaque lancement. Un fichier modifié ou absent déclenche une nouvelle recherche. Un chemin `claude_bin` dans `config.toml` est utilisé tel quel. ccam refuse un `claude` qui serait ccam lui-même, comparé par inode, car cela tournerait en boucle.
5. `checkLinks` compare le dossier du compte `~/.local/share/ccam/profiles/<name>/` avec ses paramètres de partage à l'aide de `share.Check`, qui n'appelle que `lstat` et `readlink`, un ou deux appels système par élément. Un lien manquant est créé sur-le-champ. Tout ce qui pourrait faire perdre des données arrête le lancement.
6. `applyMCP` ajoute les serveurs MCP partagés (voir plus bas).
7. `session.Register` enregistre la session. ccam écrit un marqueur nommé d'après son PID, contenant l'heure de démarrage du processus, dans `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (ou `~/.local/state/ccam/locks/` sans `XDG_RUNTIME_DIR`). Pour un compte qui partage l'historique des conversations, `historyCheck` cherche d'abord une session en cours dans `history-<group>-<hash>`, où le hachage est celui du dossier courant, et s'y enregistre aussi. `exec` conserve le PID, si bien que le marqueur désigne ensuite le processus Claude Code lui-même. Un marqueur plutôt qu'un verrou tenu, car les processus lancés par Claude Code hériteraient du verrou et le garderaient après la fin de la session.
8. `Plan.Exec` définit le titre du terminal à `claude · <name>` (quand la sortie d'erreur standard est un terminal), écrit l'heure de lancement dans `~/.local/state/ccam/last-used/<name>` sans `fsync`, et appelle `syscall.Exec`.

Claude Code s'exécute avec `CLAUDE_CONFIG_DIR` défini sur le dossier du compte et `CCAM_PROFILE` défini sur le nom du compte. ccam supprime toute valeur héritée de ces variables. Pour un dossier `~/.claude` enregistré avec `ccam import` sous macOS, il définit aussi `CLAUDE_SECURESTORAGE_CONFIG_DIR` à une valeur vide, afin que le compte et `claude` seul utilisent la même connexion du trousseau.

### Pourquoi exec

`syscall.Exec` remplace le processus ccam par `claude` : ccam ne reste pas en mémoire, les frappes au clavier, les signaux et le code de sortie sont ceux de Claude Code, et les arguments restent un tableau, jamais une chaîne de shell. Le travail après une session, comme repérer les serveurs ajoutés avec `claude mcp add`, a lieu au lancement suivant : ccam compare la date de modification de `.claude.json` avec `~/.local/state/ccam/mcp/<name>/seen` et n'analyse le fichier que s'il a changé.

### Comment une ligne de commande est aiguillée

Trois lignes de commande, trois chemins à travers `internal/app` :

<!-- motion:launch-route -->

Un raccourci est reconnu sur `argv[0]`, `ccam <name>` sur le premier argument, et les deux aboutissent à `syscall.Exec`. Seule une commande ccam comme `ccam usage` atteint cobra.

### Serveurs MCP partagés

Avec la stratégie `flag` par défaut, ccam transmet les serveurs MCP du groupe à Claude Code sous forme d'argument, placé avant les vôtres :

<!-- motion:launch-mcp -->

`mcp.FlagArgs` ne renvoie un argument que si le `mcp.json` du groupe existe, et `Plan.InsertArgs` le place juste après `claude`. Vos propres `--mcp-config` ou `--strict-mcp-config` s'appliquent alors toujours par-dessus. C'est un seul argument, `--mcp-config=<path>`, car le `--mcp-config` de Claude Code accepte plusieurs valeurs : écrit en deux arguments, il lirait aussi une invite qui suit comme fichier de configuration.

Avec `mcp_strategy = "sync"` dans `config.toml`, ccam écrit plutôt les serveurs dans le `.claude.json` du compte, et ne le fait pas si le compte est déjà en cours d'exécution.

### Remote Control

`ccam work rc` (ou `ccam work remote-control`) démarre le Remote Control de Claude Code en tant que `work`. `claude remote-control` refuse `--mcp-config` avant son verbe et n'en accepte aucun après, et ses sessions ne lisent les serveurs MCP que dans le `.claude.json` du compte. La stratégie de l'option fonctionne donc différemment ici :

<!-- motion:launch-rc -->

1. `remoteControl` considère le premier argument après le compte qui ne commence pas par `-` comme le verbe de Claude Code. S'il vaut `remote-control` ou `rc`, `mcp.FlagArgs` n'ajoute rien.
2. `mcp.Synced` compare l'heure de modification du `mcp.json` du groupe avec le `sync-stamp` du compte. Si `.claude.json` contient déjà les serveurs actuels du groupe, ou si le groupe n'en a aucun, Remote Control démarre sans poser de question.
3. Sinon, ccam nomme les serveurs qui manqueraient aux sessions et demande :

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. Avec `y`, `mcp.Sync` les copie comme le fait la stratégie `sync` : sous le verrou de Claude Code, après une sauvegarde, en conservant les serveurs propres du compte. Avec `N`, ou sans terminal, Remote Control démarre sans eux, et ccam indique `mcp_strategy = "sync"`, qui les copie avant chaque lancement.

ccam pose la question à chaque lancement de Remote Control tant que les serveurs ne sont pas copiés, et après une copie seulement quand le `mcp.json` du groupe change de nouveau. Une copie dans `.claude.json` l'emporte sur `--mcp-config` lors des lancements normaux suivants, donc `ccam mcp add`, `remove` et `import` mettent à jour les copies que détiennent les autres comptes ([Serveurs MCP partagés](/features/mcp)).

### Le budget de 10 ms

ccam ajoute au plus 10 ms (p95) avant le démarrage de Claude Code. Le chemin de lancement ne construit jamais l'arborescence de commandes cobra, ne démarre pas Bubble Tea et n'ouvre aucune connexion réseau, et `TestPackageInitIsCheap` échoue si les fonctions `init()` des paquets prennent plus de 5 ms au total. `scripts/perf-check.sh` mesure `ccam --dry-run work`, `ccam work` et un raccourci avec `hyperfine` et échoue au-delà de 10 ms, et `BenchmarkLaunchPath` dans `internal/launcher` exécute les étapes en boucle. Pour voir le temps de chaque étape, définissez `CCAM_TRACE=1` :

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Sécurité et cas particuliers

- Un nom ou un numéro inconnu s'arrête avec `no account "<name>"` et renvoie vers `ccam list`. Sans aucun compte, ccam vous demande d'en ajouter un avec `ccam add <name>`.
- Si un chemin partagé est devenu un fichier ordinaire ou pointe au mauvais endroit, ccam s'arrête avec `stopped because of share link problems:`, liste les problèmes et vous demande de lancer `ccam doctor --fix`. Il ne déplace pas de données sur le chemin de lancement.
- Si un autre Claude Code qui partage l'historique des conversations avec ce compte (même groupe) s'exécute dans le même dossier, ccam avertit que l'historique peut se mélanger et demande `Run anyway? [y/N]`. Sans terminal, il affiche l'avertissement et continue.
- `--dry-run` saute la synchronisation MCP, l'avis de nouveau serveur, la vérification de l'historique, le marqueur de session, le titre et l'heure de dernière utilisation, et n'exécute pas exec. Il vérifie tout de même les liens de partage, et crée un lien manquant comme le ferait un lancement normal.
- `ccam work rc` alors que `work` est déjà en cours d'exécution laisse `.claude.json` tel quel, le signale, et démarre Remote Control sans les serveurs du groupe ; relancez-le une fois la session terminée pour les copier. Avec `--dry-run`, il affiche le message mais ne pose pas de question et ne copie rien.
- Avec `mcp_strategy = "sync"`, une synchronisation échouée affiche `ccam: MCP sync failed (launching anyway): …` et Claude Code démarre avec les serveurs qu'il avait déjà.

## Voir aussi

- [Commandes raccourcis](/features/shortcuts)
- [Gérer les comptes](/features/accounts)
- [Référence des commandes](/docs/commands)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Fonctionnement de CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
