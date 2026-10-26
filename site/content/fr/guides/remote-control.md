---
title: "Utiliser Remote Control avec ccam"
description: "Lancez le Remote Control de Claude Code avec un compte via ccam <name> rc : options valables, pourquoi ccam propose de copier les serveurs MCP, dépannage."
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 4d46f49c5dea
---

`ccam work rc` démarre le Remote Control de Claude Code avec le compte `work`, pour que vous puissiez contrôler des sessions sur cet ordinateur depuis claude.ai/code ou l'application mobile Claude. Il fonctionne comme `claude rc`, à une différence près : si le compte partage des serveurs MCP, ccam demande une fois s'il faut les copier dans le `.claude.json` du compte, car ce fichier est le seul endroit où les sessions Remote Control lisent les serveurs.

## Démarrer

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

Un [raccourci](/features/shortcuts) fonctionne de la même façon : `ccwork rc` équivaut à `ccam work rc`. Lancez la commande dans le dossier où vous voulez travailler, puis ouvrez claude.ai/code ou l'application Claude, connecté avec le même compte Claude que `work`.

## Options

Deux sortes d'options se rencontrent ici : celles de ccam, qui se placent avant le nom du compte ou se trouvent dans les paramètres de ccam, et celles de Remote Control, qui se placent après `rc` et parviennent telles quelles à Claude Code.

### Options et paramètres de ccam

| Option ou paramètre | Effet | Par défaut |
| --- | --- | --- |
| `--dry-run` (avant le compte) | Affiche l'environnement et la commande, ainsi que l'avis MCP s'il y en a un. Ne pose aucune question et ne copie rien | désactivé |
| `mcp_strategy = "sync"` en haut de `~/.config/ccam/config.toml` | Copie les serveurs du groupe dans `.claude.json` avant chaque lancement, Remote Control ou non, si bien qu'aucune question n'est jamais posée | `"flag"` |
| `ccam group leave <group> work` | Retire le compte du groupe qui partage MCP avec lui. Les serveurs du groupe sont copiés dans son `.claude.json` comme les siens, si bien que Remote Control les a et qu'aucune question n'est posée ; avec `--fresh`, il démarre avec les seuls serveurs propres du compte | dans `default` |

Quand le compte partage des serveurs MCP et que son `.claude.json` ne contient pas la liste actuelle du groupe, ou qu'un serveur copié en a disparu, ccam demande `Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]`. `y` les copie ; `N` ou Entrée démarre sans eux.

### Options de Remote Control

Tout ce qui suit `rc` appartient à Claude Code. Voici les options de `claude remote-control` dans Claude Code 2.1.304 ; `ccam work rc --help` affiche la liste pour votre version.

| Option | Effet |
| --- | --- |
| `--name <name>` | Nom de la session affiché dans claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Préfixe des noms de session générés (par défaut : le nom de la machine) |
| `-c`, `--continue` | Se rattache à la dernière session que Remote Control a enregistrée pour ce dossier, dans un délai d'environ 4 heures |
| `--session-id <id>` | Se rattache à une session précise |
| `--permission-mode <mode>` | Mode d'autorisation des nouvelles sessions : `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` ou `plan` |
| `--spawn <mode>` | `same-dir` (par défaut), `worktree` (un git worktree par session) ou `session` (une seule session, quitte à sa fin) |
| `--capacity <N>` | Nombre maximal de sessions simultanées en mode `same-dir` et `worktree` (32 par défaut) |
| `--[no-]create-session-in-dir` | Crée une première session dans le dossier courant au démarrage (activé par défaut) |
| `--[no-]chrome` | Claude in Chrome pour les nouvelles sessions (désactivé par défaut) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Sortie de débogage et détaillée |

ccam lit la commande de Claude Code comme le fait Claude Code, en sautant les options avec leurs valeurs : `ccam work rc --name laptop` et `ccam work --model opus rc` démarrent donc tous deux Remote Control. Une option que ccam ne connaît pas compte comme une option sans valeur ; par prudence, placez `rc` juste après le nom du compte.

## Ce que fait ccam

Le lancement suit le chemin habituel ([Lancer Claude Code](/features/launch)), avec un seul changement à l'étape MCP :

1. ccam trouve le compte, vérifie ses liens de partage et résout `claude`, comme pour tout lancement.
2. Il voit que la commande est `rc` ou `remote-control` et omet `--mcp-config`.
3. Si le compte partage des serveurs MCP, il compare le `mcp.json` du groupe avec le `sync-stamp` de `~/.local/state/ccam/mcp/work/`, et vérifie que chaque serveur du groupe se trouve toujours dans `.claude.json`. Si `.claude.json` contient déjà les serveurs actuels, ou si le groupe n'en a aucun, il passe directement à l'étape 6.
4. Sinon, il nomme les serveurs qui manqueraient aux sessions :

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; the account's own servers stay)? [y/N]
   ```

5. Avec `y`, il enregistre une sauvegarde de `.claude.json`, prend le verrou propre à Claude Code, ajoute les serveurs du groupe à `mcpServers` et conserve les autres serveurs du compte ; si le compte possède son propre serveur portant le même nom qu'un serveur du groupe, ccam garde celui du compte et le signale. Avec `N`, sans terminal ou pendant que le compte est en cours d'exécution, il ne modifie rien et le signale :

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Il enregistre la session, fait pointer `CLAUDE_CONFIG_DIR` vers le dossier du compte et se remplace par `claude rc` via `exec`.

Remote Control continue ensuite de tourner comme serveur. Chaque session qu'il démarre pour vous est un processus enfant de ce processus et s'exécute donc en tant que `work` : avec la connexion, les paramètres et l'historique des conversations de `work`, et en comptant dans la limite d'utilisation de `work`. ccam considère le compte comme en cours d'exécution tant qu'un serveur Remote Control qu'il a lancé tourne. Un serveur lancé avec `claude remote-control` seul ne compte que tant qu'il a des sessions ; ccam ne peut pas en voir un qui n'en a encore aucune.

La simulation montre la différence. Un lancement normal transmet les serveurs du groupe en argument, Remote Control non :

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## Pourquoi ça fonctionne ainsi

**Remote Control n'accepte pas `--mcp-config`.** Pour les lancements normaux, ccam transmet les serveurs du groupe à Claude Code sous la forme `--mcp-config=<path>` et n'écrit jamais dans les fichiers de Claude Code. `claude remote-control` refuse cette option avant son nom de commande et n'en accepte aucune après ; avant ccam 0.6.0, `ccam work rc` échouait donc pour tout compte qui partageait des serveurs MCP. ccam omet désormais l'option.

**Ses sessions ne lisent les serveurs que dans `.claude.json`.** Les sessions que démarre Remote Control n'héritent pas des options de la ligne de commande. Le seul endroit où elles trouvent des serveurs MCP utilisateur est le `.claude.json` du compte. Ce fichier contient aussi les informations de connexion du compte et l'état par dossier, si bien que ccam ne peut pas le lier à une copie partagée comme il le fait avec `settings.json`. Y copier les serveurs est le seul moyen de donner la liste partagée aux sessions.

**ccam demande avant d'écrire.** Avec la stratégie de l'option par défaut, ccam ne touche pas à `.claude.json` au lancement. Remote Control est le seul lancement où il doit y écrire : il demande donc d'abord, garde une sauvegarde et ne touche qu'à `mcpServers`. (En dehors des lancements, quitter ou rejoindre le partage MCP et `ccam mcp import` y écrivent aussi, comme l'indiquent leurs plans.)

**Il n'écrit pas pendant que le compte est en cours d'exécution.** Claude Code réécrit `.claude.json` pendant une session, et une modification faite de l'extérieur serait perdue ou écrasée. C'est pourquoi ccam saute la copie dès qu'une session du compte est en cours, y compris un autre serveur Remote Control et Claude Code lancé sans ccam.

**Il demande une fois par changement.** Après une copie, `sync-stamp` contient l'heure de modification du `mcp.json` du groupe. ccam ne repose la question que lorsque ce fichier change, par exemple après `ccam mcp add`, ou quand un serveur copié a disparu de `.claude.json`. Si vous répondez `N`, aucune copie n'existe encore, et la question revient donc à chaque lancement de Remote Control.

**`ccam mcp` tient les copies à jour.** Une entrée de `.claude.json` l'emporte sur `--mcp-config` en cas de conflit de noms : une ancienne copie masquerait donc une définition plus récente, y compris lors des lancements normaux. `ccam mcp add`, `add-json`, `remove` et `import`, ainsi que `ccam migrate`, mettent donc à jour les copies de chaque compte qui n'est pas en cours d'exécution, et le signalent quand une copie a changé :

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Sans terminal, la réponse est non.** Remote Control est souvent lancé depuis un script, un service ou `nohup`, où personne ne peut répondre. ccam démarre alors sans les serveurs au lieu d'attendre.

## Problèmes et solutions

### Les serveurs MCP partagés manquent dans les sessions Remote Control

La question a reçu la réponse `N`, il n'y avait pas de terminal, ou le compte était en cours d'exécution. Arrêtez Remote Control, fermez les autres sessions du compte et relancez `ccam work rc` dans un terminal ; répondez `y`. Pour ne plus voir la question, définissez `mcp_strategy = "sync"` dans `~/.config/ccam/config.toml`. Les serveurs d'une session ne changent pas tant qu'elle tourne : démarrez donc une nouvelle session après la copie.

### "Account work is running, so its .claude.json is left alone"

Une session de `work` est encore ouverte : un `ccam work` normal, un raccourci ou un autre serveur Remote Control. Quittez-la et relancez `ccam work rc`. Les sessions d'un autre compte ne comptent pas.

### ccam pose la question à chaque fois

Vous n'avez pas encore copié les serveurs. Répondez `y` une fois, ou définissez `mcp_strategy = "sync"`. Si vous ne voulez jamais des serveurs dans Remote Control, cessez de partager MCP avec le compte sans les garder : `ccam group leave default work --fresh`, ou placez MCP dans un groupe à part dont le compte ne fait pas partie. Un simple `ccam group leave` copie les serveurs du groupe dans le `.claude.json` du compte, si bien que Remote Control les a.

### `ccam work rc` échoue immédiatement avec une erreur concernant `--mcp-config`

Votre ccam est antérieur à la version 0.6.0 et transmet `--mcp-config` à Remote Control. Lancez `ccam update`. Avec un ccam à jour, la même erreur signifie qu'une option inconnue suivie d'une valeur précédait `rc`, et que ccam a pris cette valeur pour la commande ; placez `rc` juste après le nom du compte.

### Une modification faite avec `ccam mcp` n'atteint pas Remote Control

`ccam mcp add` et `remove` ne mettent à jour les copies que dans les comptes qui ne sont pas en cours d'exécution, et le signalent : `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Arrêtez Remote Control puis redémarrez-le ; ccam voit que `mcp.json` a changé et repose la question.

### Un de mes serveurs est utilisé à la place de celui du groupe

Si le compte a son propre serveur portant le même nom qu'un serveur du groupe, la copie le laisse en place et le signale : `account work has its own MCP servers fs, so group default's servers with those names are not copied there`. Pour utiliser celui du groupe, supprimez celui du compte (`claude mcp remove -s user <name>`) ; pour donner votre version à tous les comptes, déplacez-la dans le groupe avec `ccam mcp import work <server> --force`. Les sauvegardes de `.claude.json` se trouvent dans `~/.local/state/ccam/mcp/work/backups/`, qui conserve les cinq plus récentes.

### Je veux retirer les serveurs copiés

`ccam mcp remove <server>` retire un serveur du groupe et de chaque copie. Pour revenir à l'état d'avant la copie, restaurez la sauvegarde la plus récente de `~/.local/state/ccam/mcp/work/backups/` pendant qu'aucune session du compte n'est en cours.

### Remote Control demande s'il faut faire confiance au dossier, ou quitte sans terminal

La confiance accordée aux dossiers est enregistrée par compte dans `.claude.json` : un dossier approuvé en tant que `personal` est donc nouveau pour `work`. Lancez `ccam work` une fois dans ce dossier et acceptez, puis démarrez Remote Control à cet endroit.

### "You must be logged in with a Claude account that has a subscription"

Remote Control exige une connexion avec un abonnement Claude. Vérifiez l'adresse e-mail du compte dans `ccam list`, et reconnectez-vous avec `ccam login work` si nécessaire.

### La session n'apparaît pas sur le téléphone

L'application Claude ou claude.ai/code est connectée avec un autre compte Claude. Comparez-le avec l'adresse e-mail que `ccam list` affiche pour `work`.

### ccam signale qu'une autre session Claude Code partage l'historique des conversations dans ce dossier

Un autre compte du même groupe d'historique tourne dans le même dossier. Des sessions qui poursuivent la même conversation depuis les deux comptes peuvent mélanger l'historique. Répondez `y` seulement si elles travaillent sur des conversations différentes, ou démarrez Remote Control dans un autre dossier.

## Voir aussi

- [Lancer Claude Code](/features/launch)
- [Serveurs MCP partagés](/features/mcp)
- [Partager avec des groupes](/features/share)
- [Gérer les comptes](/features/accounts)
