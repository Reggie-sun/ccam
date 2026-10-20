---
title: "Passer au compte suivant"
description: "ccam next ouvre Claude Code avec le compte suivant le dernier utilisé ; --auto prend le plus de marge, --continue reprend et --fallback relance les tâches -p."
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next` ouvre Claude Code avec un autre compte sans que vous ayez à le nommer. Utilisez-le quand un compte atteint sa limite d'utilisation : `ccam next` seul prend le compte qui suit celui utilisé en dernier, `--auto` prend celui qui a le plus de marge dans sa fenêtre de 5 heures, et `--continue` poursuit la dernière conversation dans le nouveau compte. Pour les scripts qui exécutent `claude -p`, `ccam --fallback` fait ce changement de lui-même (voir [Relancer à la limite](#run-again-at-the-limit)).

## Utilisation

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
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
| `--auto` | Choisir le compte connecté selon l'utilisation (voir plus bas), en ignorant les comptes à 100 % | désactivé |
| `--max-usage <1-100>` | Avec `--auto` : préférer les comptes sous ce niveau d'utilisation | `auto_max_usage`, sinon 100 |
| `--strategy headroom\|expiring` | Avec `--auto` : l'utilisation sur 5 heures la plus basse, ou la réinitialisation sur 5 heures la plus proche | `auto_strategy`, sinon `headroom` |
| `--wait[=<duration>]` | Avec `--auto` : quand aucun compte n'a de marge, attendre la réinitialisation du premier, au plus cette durée | désactivé ; `5h` s'il est donné seul |
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

`--auto` décide selon l'utilisation plutôt que selon l'ordre des emplacements. `cli.pickAuto` garde les comptes connectés, demande leur utilisation à tous en même temps avec le même `Fetcher.FetchAll` que `ccam usage`, et laisse `usage.Policy.Choose()` choisir :

<!-- motion:next-auto -->

Un compte est jugé sur ses fenêtres de 5 heures et hebdomadaire et, quand `--model` figure parmi les arguments de Claude Code, aussi sur la fenêtre hebdomadaire de ce modèle (`opus`, `sonnet`, `haiku` ou `fable` dans le nom). `Policy.Choose()` répartit les comptes en trois groupes :

1. **Marge** : chaque fenêtre sous `auto_max_usage` (`config.toml`, 100 par défaut). La stratégie choisit parmi eux : `headroom` (par défaut) prend l'utilisation sur 5 heures la plus basse, `expiring` prend la fenêtre de 5 heures qui se réinitialise le plus tôt, afin que son quota restant serve avant d'être perdu.
2. **Au-dessus** : sous 100 %, mais à `auto_max_usage` ou au-delà. Utilisés seulement quand aucun compte n'a de marge, l'utilisation sur 5 heures la plus basse d'abord.
3. **Plein** : une fenêtre à 100 %, ou une utilisation sur 5 heures inconnue. Jamais choisis.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

Quand rien ne peut être choisi, ccam quitte avec le code 75 (`EX_TEMPFAIL`), pour qu'un script distingue « réessayer plus tard » d'un échec. Avec `--wait`, il attend plutôt que la fenêtre bloquante du premier compte plein se réinitialise (30 secondes après l'heure de réinitialisation, quand le point de terminaison montre la nouvelle fenêtre), affiche `ccam: no account has room; waiting until 15:10` et choisit à nouveau ; quand aucun compte n'indique quand il se réinitialise, il regarde de nouveau au bout de 10 minutes. Si rien ne se libère avant la fin de l'attente, il quitte avec 75. La consultation suit les règles de `ccam usage` : les chiffres de moins de 60 secondes sont réutilisés, chaque compte dispose de 3 secondes, et un compte qui ne répond pas, ou dont la connexion a expiré, est jugé sur ses chiffres enregistrés. Une fenêtre enregistrée dont l'heure de réinitialisation est passée compte pour 0 %, si bien qu'un compte dont la fenêtre a été réinitialisée entre-temps n'est pas écarté. `--auto` ne renouvelle jamais les jetons. Contrairement à la rotation simple, `--auto` peut choisir le compte utilisé en dernier s'il a encore le plus de marge.

### --continue

Le `--continue` de Claude Code ouvre la conversation la plus récente du dossier courant, ce qui ne fonctionne que si le nouveau compte peut voir cette conversation. Avec `--continue`, ccam ne considère donc que les comptes du même groupe `history` que le compte utilisé en dernier, puis place `--continue` devant vos arguments. Les comptes dont l'historique est `local` ne partagent avec personne, donc `--continue` depuis un tel compte échoue toujours. Activez d'abord le partage avec `ccam share <name> --history on` (voir [partage](/docs/sharing)).

## Relancer à la limite

`ccam next` sert à une session interactive qui s'est arrêtée. Une exécution non interactive, `claude -p` dans un script ou une tâche planifiée, ne peut pas vous attendre : `--fallback` avant le nom du compte laisse donc ccam changer de compte pour elle :

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` ne fonctionne qu'avec `-p` ou `--print` ; sinon ccam s'arrête et renvoie vers `ccam next --auto --continue`.
2. L'entrée standard est lue une seule fois et donnée à chaque tentative.
3. ccam démarre `claude` comme processus enfant au lieu d'`exec`, retient sa sortie standard et laisse passer la sortie d'erreur standard au fur et à mesure.
4. Quand `claude` se termine avec 0, ou échoue pour une autre raison, la sortie retenue et le code de sortie sont transmis sans modification.
5. Un échec compte comme la limite d'utilisation quand le message de Claude Code le dit (`You've hit your limit`, `usage limit reached`), ou quand la consultation de l'utilisation du compte montre sa fenêtre de 5 heures ou hebdomadaire à 100 %.
6. La sortie de cette tentative est alors abandonnée et la même commande s'exécute avec un autre compte connecté : celui dont l'utilisation sur 5 heures est la plus basse, puis les comptes dont l'utilisation ne peut pas être consultée, dans l'ordre des numéros. Les comptes à 100 % sont ignorés, et chaque compte est essayé au plus une fois.
7. Avec `-c`, `--continue`, `-r` ou `--resume` parmi les arguments, seuls les comptes du même groupe d'historique sont essayés, afin que la conversation puisse être trouvée.

Le compte suivant est choisi selon les mêmes règles que `--auto` (`auto_max_usage`, `auto_strategy`, la fenêtre `--model`). Quand il ne reste aucun compte, ccam affiche la sortie de la dernière tentative, indique qu'aucun compte connecté n'a plus de marge et quitte avec le code 75. Avec `--wait` avant le nom du compte (`ccam --fallback --wait=2h work -p …`), il attend plutôt la réinitialisation du premier compte, puis peut réessayer chaque compte ; il attend toujours au moins une fois d'abord, afin qu'un compte qui vient d'atteindre sa limite ne soit pas relancé aussitôt. `ccam --dry-run --fallback work -p …` affiche la commande de la première tentative.

Les sessions interactives ne changent pas de compte d'elles-mêmes. À la place, avec `ccam statusline install`, la ligne d'état ajoute `→ ccam next --auto --continue` dès que l'utilisation sur 5 heures ou hebdomadaire atteint 90 % ([Limites d'utilisation](/features/usage)).

## Sécurité et cas particuliers

- Les comptes sans connexion ne sont jamais choisis. Sans autre compte connecté, `ccam next` seul s'arrête avec `no logged-in account other than <name>`.
- Avec `--continue` et aucun compte qui partage l'historique : `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Avec `--auto` et aucun compte connecté : `no logged-in account to pick`. Quand chaque compte est plein ou que son utilisation est inconnue : `` no account with known usage and headroom (see `ccam usage`) ``. Les deux se terminent avec le code de sortie 75.
- `--max-usage`, `--strategy` et `--wait` ne fonctionnent qu'avec `--auto`, et une valeur hors de 1-100 ou une autre stratégie est refusée.
- `--auto` fait des requêtes réseau même avec `--dry-run`. Il ne fait que lire et n'écrit jamais dans les fichiers de jetons.
- Le compte choisi démarre par le même chemin que `ccam <name>`, si bien que les mêmes vérifications s'appliquent, y compris l'avertissement quand un autre Claude Code qui partage le même historique s'exécute déjà dans ce dossier.
- `last-used` n'est écrit que lorsque Claude Code démarre réellement. `--dry-run` n'y touche pas, si bien qu'une simulation ne fait pas avancer la rotation.
- `--fallback` retient la sortie standard jusqu'à la fin d'une tentative : une sortie en flux comme `--output-format stream-json` arrive donc d'un seul bloc à la fin. Les fichiers que la première tentative a modifiés avant d'atteindre la limite restent modifiés ; la tentative suivante reprend l'invite depuis le début.

## Voir aussi

- [Limites d'utilisation](/features/usage)
- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Gérer les limites d'utilisation au quotidien](/guides/usage-limit-workflow)
