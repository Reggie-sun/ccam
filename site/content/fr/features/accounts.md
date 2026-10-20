---
title: "Gérer les comptes"
description: "Ajoutez, connectez, listez, renommez, retirez et importez des comptes Claude Code : ccam add, ccam login, ccam list, ccam rename, ccam remove et ccam import."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

Chaque compte ccam est un dossier de configuration Claude Code distinct, avec sa propre connexion. Ces commandes créent, connectent, listent, renomment et suppriment des comptes, et enregistrent les dossiers que vous avez déjà.

## Utilisation

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Partout où une commande accepte `<name|number>`, vous pouvez indiquer le nom du compte ou son numéro d'emplacement, la colonne `#` de `ccam list`.

`ccam add` sans option dans un terminal ouvre un assistant. Il demande le nom, puis affiche une liste à cocher des éléments partagés (tout est coché sauf l'historique des conversations), puis demande `Create a shortcut command for work? [y/N]` et, si vous répondez `y`, `Shortcut name [ccwork]:`. L'une des options `--share`, `--shortcut` ou `--no-login`, ou un tube à la place d'un terminal, fait sauter l'assistant ; un nom est alors obligatoire et aucun raccourci n'est créé sauf si `--shortcut` en nomme un.

`ccam list` affiche l'emplacement, le nom, le raccourci, l'adresse e-mail, les éléments partagés et la dernière utilisation, avec `*` devant le compte par défaut :

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

Avec `--json`, chaque compte est un objet avec `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (omis quand le compte n'a jamais été lancé) et `share`. `share` associe chaque élément (`mcp`, `settings`, `memory`, `extensions`, `history`) à son nom de groupe ou à `"local"`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `ccam add --no-login` | Créer le compte sans se connecter ; connectez-vous plus tard avec `ccam login <name>` | désactivé |
| `ccam add --share <item=group>` | Partager un élément avec un groupe, ou `item=local` pour le garder privé. Répétable | `mcp`, `settings`, `memory`, `extensions` dans `default` ; `history` local |
| `ccam add --shortcut <command>` | Créer aussi cette commande raccourci dans `~/.local/bin` | aucun |
| `ccam list --json` | Afficher les comptes en JSON | désactivé |
| `ccam remove -y`, `--yes` | Supprimer sans demander | désactivé |
| `ccam import --name <name>` | Nom du compte pour le dossier indiqué | déduit du nom du dossier |
| `ccam import -y`, `--yes` | Importer tous les dossiers trouvés sans demander | désactivé |
| `--dry-run` | Respecté par `ccam list` et `ccam completion`, qui ne modifient rien. Refusé par `add`, `login`, `logout`, `use`, `rename`, `remove` et `import`, avec `` --dry-run is not supported by `ccam add`; nothing was changed `` (en nommant la commande) | désactivé |

`ccam login`, `ccam logout`, `ccam use` et `ccam rename` n'ont pas d'options propres. `ccam completion` accepte `bash` ou `zsh` ; toute autre valeur échoue avec `only bash and zsh are supported`.

Les noms de compte comptent de 1 à 32 caractères parmi `a-z`, `0-9`, `-` et `_`, et commencent par une lettre ou un chiffre. Un nom ne peut pas être composé uniquement de chiffres, car `ccam 2` désigne l'emplacement 2. Un nom ne peut pas non plus être une commande ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), si bien que `ccam <name>` n'est jamais ambigu. La liste est `config.Reserved` dans `internal/config`.

## Fonctionnement

`ccam add work --shortcut ccwork` exécute ces étapes dans `addProfile` (`internal/cli/profiles.go`) :

<!-- motion:accounts-add -->

1. `config.ValidateName` vérifie le nom, et `shim.Check` vérifie le raccourci avant que quoi que ce soit ne soit créé.
2. `profile.Add` crée `~/.local/share/ccam/profiles/work/` avec le mode `0700` et attribue au compte le plus petit numéro d'emplacement libre (`Config.NextSlot`), de sorte que le numéro d'un compte supprimé est réutilisé. Le premier compte devient le compte par défaut.
3. `applyShare` relie le nouveau dossier vide à ses groupes de partage avec `share.Build` puis `Apply`, par exemple `settings.json` vers `groups/default/settings.json`. En cas d'échec, le nouveau dossier est supprimé, car un dossier restant bloquerait chaque nouvelle tentative avec `folder already exists`.
4. Le compte est enregistré dans `~/.config/ccam/config.toml` via `fsutil.AtomicWrite` : un fichier temporaire, `fsync`, puis `rename`, si bien qu'un plantage ne laisse jamais une configuration à moitié écrite.
5. `runLogin` lance `claude auth login` avec `CLAUDE_CONFIG_DIR` défini sur le dossier du compte. Si la connexion échoue, le compte est conservé, et ccam affiche l'erreur et vous indique de lancer `ccam login work` plus tard.
6. `profile.ReadEmail` extrait `oauthAccount.emailAddress` de `profiles/work/.claude.json`. Ce fichier peut peser plusieurs mégaoctets, donc `TopLevelValue` recherche cette seule clé au lieu d'analyser tout le fichier.
7. `shim.Install` relie `~/.local/bin/ccwork` au binaire ccam, qui recherche `argv[0]` dans `config.toml` au démarrage.

Une connexion peut prendre plusieurs minutes, et un autre ccam peut enregistrer la configuration entre-temps. Les commandes qui attendent enregistrent donc via `env.update`, qui recharge `config.toml`, applique la seule modification (ici, l'adresse e-mail) et enregistre. `ccam remove` recharge après sa question pour la même raison.

`ccam login` exécute les étapes 5 et 6 pour un compte existant. `ccam logout` exécute `claude auth logout` dans le compte et efface l'adresse e-mail enregistrée. `ccam use` ne fait que réécrire `default` dans `config.toml`.

### Se connecter sur une machine distante

Via SSH, la partie navigateur de `claude auth login` peut échouer. Claude Code attend la connexion sur un serveur de rappel à l'adresse `http://localhost:<port>/callback` sur le serveur, mais le navigateur qui termine la connexion s'exécute sur votre ordinateur, où cette adresse ne mène nulle part, et le navigateur affiche une erreur comme `ERR_CONNECTION_REFUSED`. L'adresse contient malgré tout le code à usage unique, que ccam peut donc transmettre :

<!-- motion:accounts-relay -->

`runLogin` ne confie pas votre terminal à Claude Code. Il relie l'entrée standard de Claude Code à un tube, et `relayLogin` y copie chaque ligne que vous tapez, si bien qu'un code collé à l'invite `Paste code here` de Claude Code arrive toujours. Quand une ligne est une adresse de rappel, `callbackURL` ne l'accepte que si l'hôte est `localhost`, `127.0.0.1` ou `[::1]`, qu'un port est indiqué, que le chemin est `/callback` et que `code` et `state` sont tous deux présents ; `http://` peut être omis. `deliverCallback` demande alors cette adresse sur le serveur avec une limite de 10 secondes et sans suivre les redirections, comme l'aurait fait le navigateur. La requête atteint le serveur de rappel du `claude auth login` en attente (la ligne en pointillés du diagramme). Claude Code reçoit le code et enregistre la connexion, sous Linux dans `.credentials.json` dans le dossier du compte. Une ligne collée ne peut jamais amener ccam à contacter un autre hôte.

ccam affiche ceci avant le début de la connexion :

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

En cas de réussite, il affiche `→ Passed the sign-in to Claude Code on this machine.`

Pendant que la connexion attend une saisie dans un terminal, ccam désactive l'écho et le mode ligne du terminal et affiche un `*` pour chaque caractère que vous tapez, si bien qu'un code ou une adresse de rappel collés n'apparaissent pas à l'écran. `readLine` gère lui-même l'édition de la ligne : `Enter` envoie la ligne, `Backspace` efface un caractère, et `Ctrl+D` sur une ligne vide termine la saisie. Le terminal est restauré à la fin de la connexion, ainsi qu'à `Ctrl+C`, `SIGTERM` ou `SIGHUP`, qui terminent toujours ccam. Sans terminal sur l'entrée standard et la sortie d'erreur standard, les lignes sont copiées telles qu'elles arrivent.

### Renommer, supprimer et importer

`ccam rename` et `ccam remove` agissent sur le dossier du compte et sur ce que ccam conserve ailleurs sous le nom du compte. `nameState` liste ces chemins : `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` et les caches d'utilisation `~/.local/state/ccam/usage/<name>.json` et `<name>.statusline.json`. `ccam remove work` les traite dans cet ordre :

<!-- motion:accounts-remove -->

`ccam remove` refuse d'abord si le compte est en cours d'exécution, puis pose la question. Si le compte partage les extensions, `RewritePluginPaths` réécrit les fichiers `installed_plugins.json` et `known_marketplaces.json` du groupe, car Claude Code enregistre des chemins de plugins absolus qui peuvent passer par le dossier de ce compte. `profile.Trash` déplace le dossier vers `~/.local/share/ccam/trash/<name>-<time>/`, le compte est retiré de `config.toml`, les chemins de `nameState` sont supprimés pour qu'un futur compte du même nom n'en hérite pas, et le lien du raccourci est supprimé. Pour récupérer le compte, remettez le dossier en place depuis la corbeille.

`ccam rename` renomme le dossier, réécrit les chemins de plugins dans son `plugins/` de l'ancien dossier vers le nouveau, conserve le statut de compte par défaut et déplace les chemins de `nameState` vers le nouveau nom. La commande raccourci garde son nom. Si la réécriture des plugins échoue, le dossier reprend son ancien nom, puisque la configuration n'est pas enregistrée.

`ccam import` enregistre un dossier là où il se trouve, comme `dir` dans `config.toml`. Il ne le déplace ni ne le copie jamais : un jeton de renouvellement copié ferait que les copies se déconnectent mutuellement. Sans dossier en argument, il cherche les dossiers `~/.claude` et `~/.claude-*` non enregistrés et pose la question pour chacun ; les noms viennent du dossier (`~/.claude` devient `main`, `~/.claude-work` devient `work`, `~/.claude-2` devient `acct2`), et un dossier dont le nom ne peut pas être déduit est listé avec la commande `--name` à utiliser. Sans `CLAUDE_CONFIG_DIR`, Claude Code conserve `.claude.json` dans votre dossier personnel à côté de `~/.claude` ; quand vous importez `~/.claude` et qu'il n'a pas de `.claude.json` à lui, ccam crée donc le lien `~/.claude/.claude.json` vers `~/.claude.json`. `claude` seul et `ccam main` voient alors le même compte. Les comptes importés ne partagent rien tant que vous n'avez pas lancé `ccam share`, et aucun raccourci n'est créé.

`ccam completion bash|zsh` affiche un script de complétion qui complète aussi les noms de compte.

## Sécurité et cas particuliers

- `ccam rename` et `ccam remove` refusent d'agir pendant que le compte s'exécute (`refuseRunning`), avec `work is running; quit its Claude Code session first`, car la session continue d'écrire dans le dossier.
- `ccam remove` sur un compte importé ne fait que le désenregistrer : `✓ main unregistered (folder ... left in place)`. Votre `~/.claude` n'est pas déplacé dans la corbeille.
- `ccam add` refuse un nom dont le dossier existe déjà dans `profiles/`, et renvoie vers `ccam import`.
- `ccam import` refuse un dossier déjà enregistré, un chemin qui n'est pas un dossier et un nom déjà pris.
- `ccam logout` et `ccam login` ne vérifient pas si le compte est en cours d'exécution.
- Si le relais ne peut pas joindre Claude Code, par exemple une fois la connexion terminée, il affiche `ccam: could not pass the sign-in to Claude Code:` avec la raison. Une erreur HTTP signifie généralement que le code a expiré ; relancez la connexion.

## Voir aussi

- [Lancer Claude Code](/features/launch)
- [Commandes raccourcis](/features/shortcuts)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Plusieurs comptes Claude Code](/guides/multiple-claude-code-accounts)
