---
title: "Serveurs MCP partagés"
description: "Gardez une seule liste de serveurs MCP par groupe avec ccam mcp add, add-json, list, remove et import, transmise via --mcp-config ou synchronisée dans chaque compte."
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 46a7aa716f99
---

Claude Code conserve les serveurs MCP utilisateur dans `.claude.json`, le même fichier que celui qui contient la connexion du compte ; ce fichier ne peut donc pas être partagé entre comptes. ccam conserve les serveurs partagés dans un fichier de groupe, `mcp.json`, et le transmet à chaque compte du groupe au lancement. Utilisez `ccam mcp` pour ajouter un serveur une seule fois pour tous vos comptes, et `ccam mcp import` pour transférer aux autres un serveur que vous avez ajouté pendant une session.

## Utilisation

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Une modification s'applique à partir du lancement suivant de chaque compte ; les sessions en cours gardent les serveurs avec lesquels elles ont démarré :

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

Ajouter un nom que le groupe possède déjà remplace sa définition et affiche `updated` au lieu de `added`. Les serveurs n'atteignent un compte que s'il partage MCP avec ce groupe (`ccam share <name> --mcp on`) ; quand aucun compte ne partage encore MCP, `add` le signale.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `-g`, `--group <group>` | groupe cible pour `add`, `add-json` et `remove` | le seul groupe par lequel un compte partage MCP, sinon `default` |
| `-- <command> [args…]` (`add`) | exécuter cette commande comme serveur stdio | |
| `--url <url>` (`add`) | URL du serveur distant (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transport d'un serveur `--url` | `http` |
| `--header <'Name: value'>` (`add`) | en-tête HTTP, répétable | aucun |
| `-e`, `--env <KEY=VALUE>` (`add`) | variable d'environnement pour un serveur stdio, répétable | aucune |
| `--force` (`import`) | remplacer un serveur du groupe de même nom qui diffère ; l'ancien `mcp.json` est sauvegardé | désactivé |
| `--keep-local` (`import`) | conserver les entrées dans le `.claude.json` du compte | désactivé |
| `-h`, `--help` | aide de chaque commande | |
| `--dry-run` (globale) | respecté uniquement par `mcp list` ; `add`, `add-json`, `remove` et `import` le refusent et ne modifient rien | désactivé |

`add` exige soit une commande après `--`, soit `--url` ; sinon il échoue avec `a command to run (after --) or --url is required`. `--header` s'applique aux serveurs `--url` et `--env` aux serveurs à commande. `add-json` accepte le même JSON que `claude mcp add-json`, et la définition doit être un objet avec un `command` ou un `url`. `import` utilise toujours le groupe par lequel le compte partage MCP, et `list` affiche tous les groupes ; `-g` ne les modifie donc pas.

## Fonctionnement

Les fichiers concernés :

- `~/.local/share/ccam/groups/<group>/mcp.json` : les serveurs partagés, `{"mcpServers": {…}}`, le format que lit `--mcp-config`. Il est écrit avec `fsutil.AtomicWrite` et le mode `0600`, car les entrées de serveur contiennent souvent des clés d'API.
- `~/.local/share/ccam/profiles/<name>/.claude.json` : le fichier propre à Claude Code, avec les serveurs du compte. ccam ne lit que sa valeur `mcpServers` avec le scanner `profile.TopLevelValue` et n'analyse jamais le fichier entier sur le chemin de lancement.
- `~/.local/state/ccam/mcp/<name>/` : le fichier `seen`, le `sync-stamp` et les sauvegardes de `.claude.json`.

Il existe deux stratégies, définies par `mcp_strategy` dans `~/.config/ccam/config.toml`. La stratégie par défaut est celle de l'option :

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` trouve le groupe MCP du compte. `mcp.FlagArgs` vérifie avec un seul `stat` que le `mcp.json` du groupe existe, sans l'analyser.
2. `Plan.InsertArgs` place `--mcp-config=<path>` juste après `claude`, avant vos propres arguments, si bien que vos propres `--mcp-config` ou `--strict-mcp-config` s'appliquent toujours par-dessus.
3. `mcp.NewLocal` vérifie si le `.claude.json` du compte a changé depuis le dernier lancement, et n'analyse ses serveurs que dans ce cas.
4. `syscall.Exec` remplace ccam par Claude Code.

Le chemin est transmis en un seul argument, `--mcp-config=<path>`, et non en deux. Le `--mcp-config` de Claude Code accepte plusieurs valeurs, donc `claude --mcp-config a.json "prompt"` lirait l'invite comme un second fichier de configuration et échouerait avec `MCP config file not found`. La forme en un seul argument s'arrête là, et une invite ou une sous-commande comme `mcp list` qui suit fonctionne. Claude Code charge les serveurs de `--mcp-config` en plus de ceux du compte ; en cas de conflit de noms, l'entrée de `.claude.json` l'emporte. La stratégie de l'option n'écrit jamais dans `.claude.json`.

Un serveur ajouté pendant une session avec `claude mcp add` n'arrive que dans le `.claude.json` de ce compte. ccam ne peut pas le voir sur le moment, puisque Claude Code l'a remplacé ; il vérifie donc au lancement suivant. Le fichier `seen` contient la date de modification de `.claude.json` et les noms de serveurs déjà signalés. Quand la date diffère, ccam compare les noms avec le groupe et affiche une seule fois une indication pour chaque nouveau serveur, sans attendre de saisie :

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

Avec `mcp_strategy = "sync"`, les serveurs apparaissent dans `claude mcp list` et `/mcp` comme des serveurs utilisateur ordinaires, et ccam les écrit dans le `.claude.json` de chaque compte avant le lancement :

<!-- motion:mcp-sync -->

1. `mcp.Sync` compare la date de modification du `mcp.json` du groupe avec `sync-stamp`. Si elle n'a pas changé, le `.claude.json` du compte n'est pas ouvert.
2. Il prend le verrou propre à Claude Code, le dossier `.claude.json.lock` à côté du fichier, créé avec `mkdir`. Il attend jusqu'à 2 secondes et considère comme périmé un verrou de plus de 10 secondes, comme le fait la bibliothèque de verrous de Claude Code.
3. Il enregistre le `.claude.json` actuel dans `backups/` du dossier d'état et conserve les 5 plus récents.
4. Il définit chaque serveur du groupe dans `mcpServers`, retire les serveurs qu'il avait écrits et que le groupe n'a plus, conserve les serveurs propres au compte, et écrit le fichier de façon atomique avec son mode d'origine.
5. Il enregistre le nouveau tampon et les noms qu'il gère, et Claude Code démarre sans `--mcp-config`.

La synchronisation est sautée, avec `ccam: account work is already running, skipping MCP sync`, quand une session du compte est en cours, car Claude Code réécrit `.claude.json` pendant son exécution. Une synchronisation échouée affiche `ccam: MCP sync failed (launching anyway): …` et lance quand même. L'indication pour les serveurs ajoutés pendant une session n'existe qu'avec la stratégie de l'option.

`ccam mcp import` déplace les serveurs propres à un compte dans son groupe :

<!-- motion:mcp-import -->

ccam lit les serveurs du compte avec `mcp.LocalServers` et compare chacun avec le serveur du groupe de même nom à l'aide de `mcp.SameDef`, qui ignore l'ordre des clés et les espaces. Un serveur absent du groupe, ou présent avec la même définition, est copié. Un serveur différent est ignoré, car la définition du groupe est celle qu'utilisent tous les autres comptes :

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Avec `--force`, le serveur du groupe est remplacé, et le `mcp.json` précédent est d'abord enregistré dans `.ccam-backup/<time>/mcp.json` dans le dossier du groupe. Ensuite `mcp.RemoveLocal` retire les entrées copiées du `.claude.json` du compte, sous le même verrou et après une sauvegarde dans le dossier d'état. C'est nécessaire, car l'entrée du compte l'emporterait sur celle du groupe en cas de conflit de noms. `--keep-local` les laisse en place.

## Sécurité et cas particuliers

- **Connexions.** Les serveurs qui se connectent avec OAuth conservent leurs jetons dans le `.credentials.json` de chaque compte (`mcpOAuth`), pas dans la liste des serveurs ; partager un serveur ne partage donc pas sa connexion. Connectez-vous une fois par compte avec `/mcp`, ou copiez les connexions depuis `~/.claude` avec [`ccam migrate`](/features/migrate). Le groupe ne reçoit jamais de fichier de connexion.
- **Comptes en cours d'exécution.** `import` copie vers le groupe mais laisse les entrées propres d'un compte en cours d'exécution, et le signale ; relancez-le une fois la session terminée. La synchronisation saute les comptes en cours d'exécution.
- **Fichiers endommagés.** Si `.claude.json` n'est pas un JSON valide, la synchronisation n'y touche pas et signale l'erreur. Si le verrou ne peut pas être pris en 2 secondes, ccam signale `could not lock … (Claude Code is writing the file)`.
- **Comptes qui ne partagent pas MCP.** `import` refuse et vous indique de lancer d'abord `ccam share side --mcp on`.
- **Noms de groupe.** `-g` n'accepte qu'un nom de groupe ; il devient un dossier sous `groups/`, si bien que des noms comme `../x` sont refusés.

## Voir aussi

- [Partage entre comptes](/features/share)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Copier une configuration existante](/features/migrate)
- [Toutes les commandes](/docs/commands)
