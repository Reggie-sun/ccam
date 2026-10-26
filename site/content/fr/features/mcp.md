---
title: "Serveurs MCP partagés"
description: "Gardez une liste de serveurs MCP par groupe avec ccam mcp add, add-json, list, remove et import, transmise via --mcp-config ou synchronisée dans chaque compte."
updated: 2026-10-26
order: 5
translation: machine
sourceHash: 27a06ad35295
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
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
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
only in account side, for /home/me/code/app only: notes
```

Sous les groupes, `ccam mcp list` nomme chaque serveur qui n'atteint pas les autres comptes :

- `only in account work: pg` pour un serveur qu'un compte a ajouté pour lui-même ;
- `only in account side, for /home/me/code/app only: notes` pour un serveur ajouté pour un seul dossier (`claude mcp add` seul, voir plus bas) ;
- `account work has its own …, used in its sessions in place of the group's` quand un serveur propre à un compte porte le nom d'un serveur du groupe ;
- le `.mcp.json` de ce dossier, avec les comptes dont l'approbation est attendue (`waiting for approval in`) ;
- les serveurs qui sont `in no account, only in ~/.claude.json (plain claude's)`, avec la commande [`ccam migrate --mcp`](/features/migrate) qui les copie.

Ajouter un nom que le groupe possède déjà remplace sa définition et affiche `updated` au lieu de `added`. Les serveurs n'atteignent un compte que s'il appartient à ce groupe et que le groupe partage MCP (`ccam group set <group> --add mcp`) ; sinon, `add` le signale.

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
| `--dry-run` (globale) | pris en charge uniquement par `mcp list` ; `add`, `add-json`, `remove` et `import` le refusent et ne modifient rien | désactivé |

`add` exige soit une commande après `--`, soit `--url` ; sinon il échoue avec `a command to run (after --) or --url is required`. `--header` s'applique aux serveurs `--url` et `--env` aux serveurs à commande. `add-json` accepte le même JSON que `claude mcp add-json`, et la définition doit être un objet avec un `command` ou un `url`. `import` utilise toujours le groupe par lequel le compte partage MCP, et `list` affiche tous les groupes ; `-g` ne les modifie donc pas.

## Fonctionnement

Les fichiers concernés :

- `~/.local/share/ccam/groups/<group>/mcp.json` : les serveurs partagés, `{"mcpServers": {…}}`, le format que lit `--mcp-config`. Il est écrit avec `fsutil.AtomicWrite` et le mode `0600`, car les entrées de serveur contiennent souvent des clés d'API, et uniquement sous le verrou `<group>/.mcp.json.lock`, si bien que deux commandes `ccam mcp add` exécutées en même temps conservent chacune leurs serveurs.
- `~/.local/share/ccam/profiles/<name>/.claude.json` : le fichier propre à Claude Code, avec les serveurs du compte. ccam ne lit que sa valeur `mcpServers` avec le scanner `profile.TopLevelValue` et n'analyse jamais le fichier entier sur le chemin de lancement.
- `~/.local/state/ccam/mcp/<name>/` : le fichier `seen`, le `sync-stamp` (avec un hachage de chaque copie écrite par ccam, pour distinguer une copie que personne n'a modifiée d'une copie que vous avez éditée), `sync-checked` (ce que la stratégie de synchronisation a vu en dernier), `kept` (les serveurs que le compte a gardés en quittant le partage MCP) et les sauvegardes de `.claude.json`.

Il existe deux stratégies, définies par `mcp_strategy` dans `~/.config/ccam/config.toml`. La stratégie par défaut est celle de l'option :

<!-- motion:mcp-flag -->

1. ccam trouve le groupe MCP du compte (`launcher.applyMCP`) et vérifie avec un seul `stat` que le `mcp.json` du groupe existe, sans l'analyser (`mcp.FlagArgs`).
2. ccam place `--mcp-config=<path>` juste après `claude`, avant vos propres arguments, si bien que vos propres `--mcp-config` ou `--strict-mcp-config` s'appliquent toujours par-dessus (`Plan.InsertArgs`).
3. ccam vérifie si le `.claude.json` du compte a changé depuis le dernier lancement, et n'analyse ses serveurs que dans ce cas (`mcp.NewLocal`).
4. ccam se remplace par Claude Code (`syscall.Exec`).

Le chemin est transmis en un seul argument, `--mcp-config=<path>`, et non en deux. Le `--mcp-config` de Claude Code accepte plusieurs valeurs, donc `claude --mcp-config a.json "prompt"` lirait l'invite comme un second fichier de configuration et échouerait avec `MCP config file not found`. La forme en un seul argument évite cela, si bien qu'une invite ou une sous-commande placée après fonctionne. `ccam work mcp list`, `get` et `remove` ne reçoivent pas l'option, car Claude Code la refuse à cet endroit ; ccam nomme plutôt les serveurs du groupe sur la sortie d'erreur standard ([détails](/features/launch#shared-mcp-servers)). Claude Code charge les serveurs de `--mcp-config` en plus de ceux du compte ; en cas de conflit de noms, l'entrée de `.claude.json` l'emporte. La stratégie de l'option n'écrit dans `.claude.json` que pour Remote Control, quand un compte quitte ou rejoint le partage MCP ([plus bas](#leaving-and-joining)), et pour `ccam mcp import`.

`claude remote-control` (`rc`) n'accepte pas `--mcp-config`, et ses sessions ne lisent les serveurs que dans `.claude.json`. Pour `ccam work rc`, ccam omet l'option et, dans un terminal, demande s'il faut copier les serveurs du groupe dans `.claude.json`, avec une sauvegarde, comme le fait la stratégie `sync`. Il pose la question tant qu'aucune copie n'existe, puis seulement après un changement du `mcp.json` du groupe ou la disparition d'un serveur copié ([Remote Control](/features/launch#remote-control)). Comme une telle copie l'emporte sur `--mcp-config` en cas de conflit de noms, `ccam mcp add`, `add-json`, `remove` et `import` mettent à jour les copies que détiennent les autres comptes du groupe :

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

La ligne `✓ Updated` n'apparaît que lorsqu'une copie a changé. Une copie créée par une ancienne version de ccam est mise à jour de la même façon, tout comme une copie après `ccam migrate`.

Un serveur ajouté pendant une session avec `claude mcp add` n'arrive que dans le `.claude.json` de ce compte. `claude mcp add` seul utilise la portée locale : le serveur se charge pour ce seul compte dans ce seul dossier (`-s user` le rend disponible pour le compte dans tous les dossiers). ccam ne peut pas le voir sur le moment, puisque Claude Code a déjà remplacé le processus ccam ; il vérifie donc au lancement suivant. Le fichier `seen` contient la date de modification de `.claude.json` et les noms de serveurs déjà signalés. Quand la date diffère, ccam compare les noms avec le groupe et affiche une seule fois une indication pour chaque nouveau serveur, sans attendre de saisie, avec les deux stratégies :

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

L'indication nomme aussi un serveur propre à un compte qui porte le nom d'un serveur du groupe, puisque celui du compte l'emporte dans ses sessions.

Avec `mcp_strategy = "sync"`, les serveurs apparaissent dans `claude mcp list` et `/mcp` comme des serveurs utilisateur ordinaires, et ccam les écrit dans le `.claude.json` de chaque compte avant le lancement :

<!-- motion:mcp-sync -->

1. ccam compare la date de modification du `mcp.json` du groupe avec `sync-stamp`, et le `.claude.json` du compte avec ce qu'il a vu en dernier (`mcp.Sync`). Si aucun des deux n'a changé, rien d'autre ne se passe. Sinon, il vérifie aussi que chaque serveur du groupe se trouve toujours dans `.claude.json`, et remet en place une copie qui a disparu (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Il prend le verrou propre à Claude Code, le dossier `.claude.json.lock` à côté du fichier, créé avec `mkdir`. Il attend jusqu'à 2 secondes et considère comme périmé un verrou de plus de 10 secondes, comme le fait la bibliothèque de verrous de Claude Code.
3. Il enregistre le `.claude.json` actuel dans `backups/` du dossier d'état et conserve les 5 plus récents.
4. Il définit chaque serveur du groupe dans `mcpServers`, retire les serveurs qu'il avait écrits et que le groupe n'a plus, et écrit le fichier de façon atomique avec son mode d'origine, seulement si quelque chose a changé. Il ne modifie ou ne supprime que les copies qu'il a écrites et que personne n'a modifiées depuis. Si le compte possède son propre serveur portant le même nom qu'un serveur du groupe, ccam garde celui du compte et le signale une fois : `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Il enregistre le nouveau tampon avec un hachage de chaque copie, et Claude Code démarre sans `--mcp-config`.

La synchronisation est sautée, avec `ccam: account work is already running, skipping MCP sync`, quand une session du compte est en cours, quelle que soit la façon dont Claude Code a été lancé, car Claude Code réécrit `.claude.json` pendant son exécution. Une synchronisation échouée affiche `ccam: MCP sync failed (launching anyway): …` et lance quand même.

### Quitter et rejoindre

Un compte qui cesse de partager MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) garde les serveurs du groupe : le plan les copie dans son `.claude.json`, avec une sauvegarde, et ses propres serveurs de mêmes noms restent. Avec `--fresh`, il n'en reçoit aucun, et le plan liste les serveurs qui cessent de se charger. Quand il rejoint un groupe qui a de nouveau ces serveurs, les copies conservées que personne n'a modifiées sont retirées, pour que les modifications ultérieures du groupe l'atteignent.

`ccam mcp import` déplace les serveurs propres à un compte dans son groupe :

<!-- motion:mcp-import -->

ccam lit les serveurs du compte, les siens et ceux qu'il a pour un seul dossier, et compare chacun avec le serveur du groupe de même nom à l'aide de `mcp.SameDef`, qui ignore l'ordre des clés et les espaces. Un serveur de dossier devient un serveur du groupe pour tous les dossiers (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`) ; un même nom avec deux définitions de dossier différentes est refusé. Les serveurs de plugin ne sont pas proposés, puisque le plugin les apporte. Un serveur absent du groupe, ou présent avec la même définition, est copié. Un serveur différent est ignoré, car la définition du groupe est celle qu'utilisent tous les autres comptes :

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Avec `--force`, le serveur du groupe est remplacé, et le `mcp.json` précédent est d'abord enregistré dans `.ccam-backup/<time>/mcp.json` dans le dossier du groupe. Ensuite ccam retire les entrées copiées du `.claude.json` du compte, les siennes et celles du dossier en une seule écriture, sous le même verrou et après une sauvegarde dans le dossier d'état (`mcp.Remove`). ccam garde le `mcp.json` du groupe verrouillé depuis sa lecture jusqu'à ce que ces entrées soient retirées, si bien que personne ne retire le serveur du groupe entre-temps. C'est nécessaire, car l'entrée du compte l'emporterait sur celle du groupe en cas de conflit de noms. Avec `--keep-local`, les entrées restent dans le `.claude.json` du compte.

## Sécurité et cas particuliers

- **Connexions.** Les serveurs qui se connectent avec OAuth conservent leurs jetons dans le `.credentials.json` de chaque compte (`mcpOAuth`), pas dans la liste des serveurs ; partager un serveur ne partage donc pas sa connexion. Connectez-vous une fois par compte avec `/mcp`. ccam ne copie jamais les connexions, et [`ccam migrate`](/features/migrate) non plus : un jeton de renouvellement copié déconnecte tous ses autres détenteurs. [`ccam ext`](/features/ext) et `ccam doctor` nomment les serveurs auxquels un compte n'est pas connecté alors qu'un autre compte l'est, et `ccam next --continue` indique ceux qui manquent au nouveau compte. Le groupe ne reçoit jamais de fichier de connexion.
- **Comptes en cours d'exécution.** `import` copie vers le groupe mais laisse les entrées propres d'un compte en cours d'exécution, et le signale ; relancez-le une fois la session terminée. La synchronisation saute les comptes en cours d'exécution. Les deux détectent aussi les sessions Claude Code lancées sans ccam.
- **Copies pour Remote Control.** La copie d'un compte en cours d'exécution n'est pas modifiée : `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Avec `mcp_strategy = "sync"`, rien n'a besoin d'être mis à jour, car chaque lancement synchronise.
- **Fichiers endommagés.** Si `.claude.json` n'est pas un JSON valide, la synchronisation n'y touche pas et signale l'erreur. Si le verrou ne peut pas être pris en 2 secondes, ccam signale `could not lock … (Claude Code is writing the file)`.
- **Comptes qui ne partagent pas MCP.** `import` refuse et vous indique de placer d'abord le compte dans un groupe qui partage MCP (`ccam group join <group> side`).
- **Noms de groupe.** `-g` n'accepte qu'un nom de groupe ; il devient un dossier sous `groups/`, si bien que des noms comme `../x` sont refusés.

## Voir aussi

- [Partager avec des groupes](/features/share)
- [Utiliser Remote Control avec ccam](/guides/remote-control)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Copier une configuration existante](/features/migrate)
- [Référence des commandes](/docs/commands)
