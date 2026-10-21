---
title: "Passer au compte suivant"
description: "ccam next ouvre Claude Code avec un autre compte ; --auto prend le plus de marge, --continue y emporte la conversation et --fallback relance les tâches -p."
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` ouvre Claude Code avec un autre compte sans que vous ayez à le nommer. Utilisez-le quand un compte atteint sa limite d'utilisation : `ccam next` seul prend le compte qui suit celui utilisé en dernier, `--auto` prend celui qui a le plus de marge, et `--continue` emporte la dernière conversation de ce dossier dans le nouveau compte, que les deux partagent l'historique ou non. Dans Claude Code, le [plugin ccam](#inside-claude-code) vous prévient quand une limite arrête une session, et `/ccam:switch` passe la conversation à un autre compte. Pour les scripts qui exécutent `claude -p`, `ccam --fallback` fait le changement de lui-même (voir [Relancer à la limite](#run-again-at-the-limit)).

## Utilisation

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam affiche le changement sous la forme `→ <from> → <to>` sur la sortie d'erreur standard, avec la raison d'un choix `--auto`, puis lance Claude Code avec le nouveau compte de la même façon que `ccam <name>`. `<from>` est le compte utilisé en dernier ou, avec `--continue`, le compte auquel appartient la conversation. Les arguments après `--` sont transmis tels quels à `claude`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--auto` | Choisir le compte connecté selon l'utilisation (voir plus bas), en ignorant les comptes à 100 % | désactivé |
| `--max-usage <1-100>` | Avec `--auto` : préférer les comptes sous ce niveau d'utilisation | `auto_max_usage`, sinon 100 |
| `--strategy headroom\|expiring\|balance` | Avec `--auto` : l'utilisation sur 5 heures la plus basse, la réinitialisation sur 5 heures la plus proche, ou le quota hebdomadaire qui serait perdu en premier | `auto_strategy`, sinon `headroom` |
| `--cross-org` | Avec `--auto` : choisir aussi des comptes d'une autre organisation (voir plus bas) | `cross_org`, sinon désactivé |
| `--wait[=<duration>]` | Avec `--auto` : quand aucun compte n'a de marge, attendre la réinitialisation du premier, au plus cette durée | désactivé ; `5h` s'il est donné seul |
| `-c`, `--continue` | Emporter la dernière conversation de ce dossier dans le nouveau compte (voir plus bas) | désactivé |
| `-- <claude args>` | Arguments transmis à `claude` après ceux de ccam | aucun |
| `--dry-run` | Respecté : affiche `→ <from> → <to>` et la ligne de commande qui serait exécutée, sans lancer Claude Code ni enregistrer de lancement | désactivé |

`--auto` et `--continue` fonctionnent ensemble. Avec `--continue`, les arguments `--resume` sont placés avant les vôtres.

Les options propres à Claude Code, comme `--dangerously-skip-permissions` ou `--model`, se placent après `--`. Contrairement à `ccam <name>`, qui transmet tout ce qui suit le nom, `ccam next` lit lui-même ses arguments : sans le `--`, il s'arrête donc avec `ccam: unknown flag: --dangerously-skip-permissions`.

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

Un compte est jugé sur ses fenêtres de 5 heures et hebdomadaire et, quand `--model` figure parmi les arguments de Claude Code (après `--`), aussi sur la fenêtre hebdomadaire de ce modèle (`opus`, `sonnet`, `haiku` ou `fable` dans le nom). `Policy.Choose()` répartit les comptes en trois groupes :

1. **Marge** : chaque fenêtre sous `auto_max_usage` (`config.toml`, 100 par défaut). La stratégie choisit parmi eux : `headroom` (par défaut) prend l'utilisation sur 5 heures la plus basse ; `expiring` prend la fenêtre de 5 heures qui se réinitialise le plus tôt, afin que son quota restant serve avant d'être perdu ; `balance` prend le compte auquel il reste le plus de quota hebdomadaire par jour jusqu'à la réinitialisation de sa semaine, en décomptant chaque session ccam en cours, afin que le quota hebdomadaire ne reste pas inutilisé et que les sessions parallèles se répartissent.
2. **Au-dessus** : sous 100 %, mais à `auto_max_usage` ou au-delà. Utilisés seulement quand aucun compte n'a de marge, l'utilisation sur 5 heures la plus basse d'abord.
3. **Plein** : une fenêtre à 100 %, ou une utilisation sur 5 heures inconnue. Jamais choisis.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organisations.** Les choix automatiques restent parmi les comptes semblables à celui que vous quittez : la même organisation d'équipe ou d'entreprise, ou bien les offres personnelles (Pro, Max) entre elles, afin qu'un travail commencé sur un compte d'entreprise ne se poursuive pas sur un compte privé, ni l'inverse. ccam connaît l'organisation de chaque compte grâce à sa connexion (`ccam login` l'enregistre ; pour les connexions plus anciennes, elle est lue dans `.claude.json`) ; un compte dont il ne peut pas la déterminer convient partout. Deux dossiers connectés au même compte Claude partagent ses limites : l'un ne prend donc jamais le relais de l'autre. `--cross-org`, ou `cross_org = true` dans `config.toml`, lève la règle d'organisation.

Quand rien ne peut être choisi, ccam quitte avec le code 75 (`EX_TEMPFAIL`), pour qu'un script distingue « réessayer plus tard » d'un échec. Avec `--wait`, il attend plutôt que la fenêtre bloquante du premier compte plein se réinitialise (30 secondes après l'heure de réinitialisation, quand le point de terminaison montre la nouvelle fenêtre), affiche `ccam: no account has room; waiting until 15:10` et choisit à nouveau ; quand aucun compte n'indique quand il se réinitialise, il regarde de nouveau au bout de 10 minutes. Si rien ne se libère avant la fin de l'attente, il quitte avec 75. La consultation suit les règles de `ccam usage` : les chiffres de moins de 60 secondes sont réutilisés, chaque compte dispose de 3 secondes, et un compte qui ne répond pas, ou dont la connexion a expiré, est jugé sur ses chiffres enregistrés. Une fenêtre enregistrée dont l'heure de réinitialisation est passée compte pour 0 %, si bien qu'un compte dont la fenêtre a été réinitialisée entre-temps n'est pas écarté. `--auto` ne renouvelle jamais les jetons. Contrairement à la rotation simple, `--auto` peut choisir le compte utilisé en dernier s'il a encore le plus de marge.

### --continue

1. ccam cherche la dernière conversation de ce dossier dans chaque compte : le `<session>.jsonl` le plus récent dans `projects/<folder>/` (le nom du dossier est le dossier de travail dont chaque caractère autre qu'une lettre ou un chiffre devient `-`, comme le nomme Claude Code). Si la session a laissé une note (une limite d'utilisation ou `/ccam:switch`, voir plus bas) et que rien de plus récent n'a été écrit depuis, c'est cette session précise qui est utilisée.
2. Le compte auquel elle appartient est celui que vous quittez, et le nouveau compte est choisi comme ci-dessus, au sein de la même organisation.
3. Si le nouveau compte partage cet historique (le même dossier de groupe), ccam lance `claude --resume <session id>`, précisément cette conversation (`--continue` ouvrirait la plus récente du dossier, quelle qu'elle soit).
4. Sinon, il lance `claude --resume <the conversation's file> --fork-session` : Claude Code copie toute la conversation, sous un nouvel identifiant de session, dans l'historique propre au nouveau compte, et l'original reste tel quel. Dès lors, le `--continue` du nouveau compte la trouve.

ccam indique la taille de la conversation : le nouveau compte la relit entièrement avec sa première réponse. Les caches de prompt ne sont jamais partagés entre organisations : avec `--cross-org`, cette première réponse coûte donc la taille complète. Partager `history` dans un groupe évite la copie, conserve une seule conversation et partage aussi les points de contrôle à partir desquels `/rewind` restaure le code (`file-history/`).

## Dans Claude Code

Le [plugin ccam](/docs/install#claude-code-plugin) ajoute deux hooks (Claude Code 2.1.78 ou plus récent) :

- Quand un tour s'arrête à une limite d'utilisation (`You've hit your session limit`, la limite hebdomadaire ou celle d'un modèle), il note la conversation et affiche une notification de bureau dans les terminaux qui le permettent (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty) : `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Un bref ralentissement imposé par le serveur ou la limite de débit d'une clé d'API n'est pas une limite d'utilisation et n'est pas pris en compte. Claude Code 2.1.234 et les versions ultérieures attendent aussi dans la session et reprennent d'eux-mêmes après une réinitialisation de la fenêtre de 5 heures ; la notification vous laisse le choix de changer de compte à la place.
- `/ccam:switch [account]` note la conversation et termine le tour avant toute requête : tapez `/exit`, puis `ccam next --continue` la poursuit, avec le compte que vous avez nommé ou avec le suivant.

La note est conservée dans `~/.local/state/ccam/handoff/<account>.json` pendant 12 heures et ne sert qu'une fois, au moment où la session suivante va démarrer ; un lancement qui échoue la laisse pour la tentative suivante. Elle appartient au dossier dans lequel la session a démarré, même après un `cd` dans le shell de Claude Code. Si le compte qu'elle nomme a été renommé ou supprimé entre-temps, ccam le signale et choisit comme d'habitude. Avec `relay = true` dans `config.toml`, ccam reste à côté des sessions interactives au lieu de se remplacer par Claude Code, et après `/ccam:switch` il poursuit de lui-même la conversation de cette session (une note laissée par une autre session n'est pas prise en compte). Cela coûte un processus ccam en attente par session, c'est pourquoi c'est désactivé par défaut. La session qui prend le relais démarre avec `ccam next --continue`, sans les arguments avec lesquels la première avait été lancée, comme `--model`.

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
8. Seuls les comptes de la même organisation sont essayés, comme avec `--auto` (`cross_org = true` lève cette règle).

Le compte suivant est choisi selon les mêmes règles que `--auto` (`auto_max_usage`, `auto_strategy`, la fenêtre `--model`). Quand il ne reste aucun compte, ccam affiche la sortie de la dernière tentative, indique qu'aucun compte connecté n'a plus de marge et quitte avec le code 75. Avec `--wait` avant le nom du compte (`ccam --fallback --wait=2h work -p …`), il attend plutôt la réinitialisation du premier compte, puis peut réessayer chaque compte ; il attend toujours au moins une fois d'abord, afin qu'un compte qui vient d'atteindre sa limite ne soit pas relancé aussitôt. `ccam --dry-run --fallback work -p …` affiche la commande de la première tentative.

Les sessions interactives ne changent pas de compte d'elles-mêmes. À la place, avec `ccam statusline install`, la ligne d'état ajoute `→ ccam next --auto --continue` dès que l'utilisation sur 5 heures ou hebdomadaire atteint 90 % ([Limites d'utilisation](/features/usage)).

## Sécurité et cas particuliers

- Les comptes sans connexion ne sont jamais choisis. Sans autre compte connecté, `ccam next` seul s'arrête avec `no logged-in account other than <name>`.
- Avec `--continue` et aucune conversation dans le dossier : `no account has a conversation in <folder> to continue`.
- Quand la règle d'organisation ne laisse aucun compte : `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Avec `--auto` et aucun compte connecté : `no logged-in account to pick`. Quand chaque compte est plein ou que son utilisation est inconnue : `` no account with known usage and headroom (see `ccam usage`) ``. Les deux se terminent avec le code de sortie 75.
- `--max-usage`, `--strategy`, `--wait` et `--cross-org` ne fonctionnent qu'avec `--auto`, et une valeur hors de 1-100 ou une autre stratégie est refusée.
- Une bifurcation copie la conversation, pas le reste de l'état du compte : les points de contrôle du `file-history/` du compte d'origine ne suivent pas, sauf si l'historique est partagé.
- `--auto` fait des requêtes réseau même avec `--dry-run`. Il ne fait que lire et n'écrit jamais dans les fichiers de jetons.
- Le compte choisi démarre par le même chemin que `ccam <name>`, si bien que les mêmes vérifications s'appliquent, y compris l'avertissement quand un autre Claude Code qui partage le même historique s'exécute déjà dans ce dossier.
- `last-used` n'est écrit que lorsque Claude Code démarre réellement. `--dry-run` n'y touche pas, si bien qu'une simulation ne fait pas avancer la rotation.
- `--fallback` retient la sortie standard jusqu'à la fin d'une tentative : une sortie en flux comme `--output-format stream-json` arrive donc d'un seul bloc à la fin. Les fichiers que la première tentative a modifiés avant d'atteindre la limite restent modifiés ; la tentative suivante reprend l'invite depuis le début.

## Voir aussi

- [Limites d'utilisation](/features/usage)
- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Gérer les limites d'utilisation au quotidien](/guides/usage-limit-workflow)
