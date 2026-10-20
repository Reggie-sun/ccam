---
title: "Passer au compte suivant"
description: "ccam next ouvre Claude Code avec le compte suivant le dernier utilisé ; --auto choisit le plus bas usage sur 5 heures et --continue reprend la conversation."
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` ouvre Claude Code avec un autre compte sans que vous ayez à le nommer. Utilisez-le quand un compte atteint sa limite d'utilisation : `ccam next` seul prend le compte qui suit celui utilisé en dernier, `--auto` prend celui qui a le plus de marge dans sa fenêtre de 5 heures, et `--continue` poursuit la dernière conversation dans le nouveau compte.

## Utilisation

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam affiche le changement sous la forme `→ <from> → <to>` sur la sortie d'erreur standard, puis lance Claude Code avec le nouveau compte de la même façon que `ccam <name>`. `<from>` est le compte utilisé en dernier. Les arguments après `--` sont transmis tels quels à `claude`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--auto` | Choisir le compte connecté dont l'utilisation sur 5 heures est la plus basse, en ignorant les comptes à 100 % | désactivé |
| `-c`, `--continue` | Ne considérer que les comptes qui partagent l'historique des conversations avec le compte utilisé en dernier, et lancer `claude --continue` | désactivé |
| `-- <claude args>` | Arguments transmis à `claude` après ceux de ccam | aucun |
| `--dry-run` | Respecté : affiche `→ <from> → <to>` et la ligne de commande qui serait exécutée, sans lancer Claude Code ni enregistrer de lancement | désactivé |

`--auto` et `--continue` fonctionnent ensemble. Avec `--continue`, `--continue` est placé avant vos propres arguments.

## Fonctionnement

Chaque lancement via ccam écrit l'heure dans `~/.local/state/ccam/last-used/<name>`. `ccam next` lit ces fichiers, vérifie quels comptes sont connectés, puis choisit.

1. `cli.candidates` construit un `next.Candidate` par compte à partir de `config.toml` : le compte, l'heure de son dernier lancement (`launcher.LastUsed`) et s'il dispose d'une connexion (`creds.Has` : un `.credentials.json` lisible ou, sous macOS, un élément du trousseau).
2. Sans `--auto`, `next.Rotate` trie les comptes par numéro d'emplacement et trouve celui lancé le plus récemment. C'est le compte que vous quittez. Si aucun compte n'a encore été lancé, c'est l'emplacement le plus bas.
3. À partir de là, il avance dans l'ordre des emplacements, en revenant du dernier au premier, et renvoie le premier compte connecté.
4. `launcher.LaunchWith` prépare et lance le compte choisi exactement comme le ferait `ccam <name>` : il vérifie les liens du compte, applique le partage MCP, enregistre la session et remplace ccam par `claude` via `syscall.Exec`.

Le parcours commence au compte utilisé en dernier et s'arrête au premier compte connecté qui le suit :

<!-- motion:next-rotate -->

`work` a été lancé le plus récemment, donc le parcours commence après lui. `side` n'a pas de connexion et est ignoré, et `lab` est choisi. Si `lab` avait lui aussi été déconnecté, le parcours serait revenu à `personal` (`% len` dans le diagramme). La rotation ne renvoie jamais le compte de départ : sans autre compte connecté, elle s'arrête avec une erreur.

### --auto

`--auto` décide selon l'utilisation plutôt que selon l'ordre des emplacements. `cli.pickAuto` garde les comptes connectés, demande leur utilisation à tous en même temps avec le même `Fetcher.FetchAll` que `ccam usage`, et laisse `usage.Best()` choisir :

<!-- motion:next-auto -->

`usage.Best()` prend le pourcentage sur 5 heures le plus bas et ignore les comptes à 100 % ou dont l'utilisation sur 5 heures est inconnue. La consultation suit les règles de `ccam usage` : les chiffres de moins de 60 secondes sont réutilisés, chaque compte dispose de 3 secondes, et un compte qui ne répond pas, ou dont la connexion a expiré, est jugé sur ses chiffres enregistrés. Une fenêtre enregistrée dont l'heure de réinitialisation est passée compte pour 0 %, si bien qu'un compte dont la fenêtre a été réinitialisée entre-temps n'est pas écarté. `--auto` ne renouvelle jamais les jetons. Contrairement à la rotation simple, `--auto` peut choisir le compte utilisé en dernier s'il a encore le plus de marge.

### --continue

Le `--continue` de Claude Code ouvre la conversation la plus récente du dossier courant, ce qui ne fonctionne que si le nouveau compte peut voir cette conversation. Avec `--continue`, ccam ne considère donc que les comptes du même groupe `history` que le compte utilisé en dernier, puis place `--continue` devant vos arguments. Les comptes dont l'historique est `local` ne partagent avec personne, donc `--continue` depuis un tel compte échoue toujours. Activez d'abord le partage avec `ccam share <name> --history on` (voir [partage](/docs/sharing)).

## Sécurité et cas particuliers

- Les comptes sans connexion ne sont jamais choisis. Sans autre compte connecté, `ccam next` seul s'arrête avec `no logged-in account other than <name>`.
- Avec `--continue` et aucun compte qui partage l'historique : `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Avec `--auto` et aucun compte connecté : `no logged-in account to pick`. Quand chaque compte est plein ou que son utilisation est inconnue : `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` fait des requêtes réseau même avec `--dry-run`. Il ne fait que lire et n'écrit jamais dans les fichiers de jetons.
- Le compte choisi démarre par le même chemin que `ccam <name>`, si bien que les mêmes vérifications s'appliquent, y compris l'avertissement quand un autre Claude Code qui partage le même historique s'exécute déjà dans ce dossier.
- `last-used` n'est écrit que lorsque Claude Code démarre réellement. `--dry-run` n'y touche pas, si bien qu'une simulation ne fait pas avancer la rotation.

## Voir aussi

- [Limites d'utilisation](/features/usage)
- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Gérer les limites d'utilisation au quotidien](/guides/usage-limit-workflow)
