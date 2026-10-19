---
title: "Limites d'utilisation"
description: "ccam usage affiche en une fois les limites sur 5 heures et hebdomadaires de chaque compte, et ccam statusline place le nom du compte et l'utilisation dans Claude Code."
updated: 2026-10-19
order: 7
translation: machine
sourceHash: 8306a8f6d950
---

`ccam usage` indique quelle part des limites sur 5 heures et sur la semaine chaque compte a utilisée, avec le temps restant avant la réinitialisation de chaque fenêtre. Utilisez-le avant de commencer une longue tâche, ou quand Claude Code vous signale qu'une limite approche, pour voir quel compte a encore de la marge. `ccam statusline` est la commande de ligne d'état pour Claude Code : elle affiche le nom du compte et l'utilisation dans Claude Code et enregistre les chiffres pour `ccam usage`.

## Utilisation

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Dans un terminal, le tableau apparaît immédiatement avec `fetching…` dans chaque ligne, et chaque ligne est redessinée quand son compte répond. `As of` indique l'ancienneté des chiffres ; `(cached)` signifie qu'ils proviennent d'un fichier enregistré et non du point de terminaison. La dernière ligne nomme le compte dont l'utilisation sur 5 heures est la plus basse en dessous de 100 %, et n'apparaît que si vous avez plus d'un compte. Un compte non connecté affiche `not logged in → ccam login <name>`, et un compte à clé d'API affiche `no limits (API)`.

Après `ccam statusline install`, Claude Code affiche une ligne comme `[work] · Opus · 5h 24% · 7d 41%`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--json` | Afficher les rapports en JSON : `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | désactivé |
| `--watch[=N]` | Redessiner toutes les `N` secondes jusqu'à ce que vous appuyiez sur Ctrl+C. `--watch` seul, ou une valeur inférieure ou égale à 0, signifie 60 | désactivé |
| `-v`, `--verbose` | Afficher les limites hebdomadaires propres à chaque modèle et une ligne `source:` (`live`, `cache` ou `statusline`) | désactivé |
| `--no-cache` | Ignorer le cache de 60 secondes et interroger à nouveau le point de terminaison | désactivé |
| `--source live\|cache` | `cache` n'interroge jamais le point de terminaison et n'affiche que les chiffres enregistrés ; `live` est le comportement par défaut | live |
| `--refresh` | Renouveler les jetons d'accès expirés, uniquement pour les comptes qui ne sont pas en cours d'exécution | désactivé |
| `--dry-run` | Accepté : `ccam usage` ne modifie rien d'autre que son propre cache. Refusé avec `--refresh` : `--dry-run cannot be combined with --refresh, which rewrites login tokens` | désactivé |

Toute autre valeur de `--source` échoue avec `--source must be live or cache`. `ccam usage` n'accepte aucun argument.

`ccam statusline` lit sur l'entrée standard le JSON de ligne d'état de Claude Code. Sa seule option est `--wrap <command>` : une commande de ligne d'état existante dont la première ligne de sortie reçoit le nom du compte en tête. `ccam statusline install` n'accepte aucune option. Les deux refusent `--dry-run` avec `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Fonctionnement

Les chiffres viennent de deux sources : le point de terminaison d'utilisation non documenté d'Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, appelé avec la propre connexion de chaque compte, et les valeurs que Claude Code transmet à sa commande de ligne d'état. Tout se trouve dans `internal/usage`.

1. `cli.usageCmd` charge `config.toml` et construit un `usage.Account` par compte : son dossier `~/.local/share/ccam/profiles/<name>/`, son dossier de session (pour savoir s'il est en cours d'exécution) et s'il a été importé sur place.
2. `Fetcher.FetchAll` lance une goroutine par compte. Toutes partagent un même `http.Client`, si bien que les connexions au point de terminaison sont réutilisées. La connexion peut prendre 2 s, et chaque compte dispose de 3 s au total. Le client ne suit pas les redirections, puisqu'une redirection renverrait le jeton ailleurs, et lit au plus 1 Mo de la réponse.
3. Dès qu'un compte a terminé, sa ligne est redessinée. Un compte qui échoue ou ne répond pas à temps se rabat sur les chiffres enregistrés, si bien qu'un compte lent ne retarde jamais les autres.

Chaque compte interroge le point de terminaison dans sa propre goroutine ; celui qui dépasse le délai affiche ses chiffres enregistrés à la place.

<!-- motion:usage-fetch -->

`personal` et `work` répondent à temps et affichent des chiffres en direct. `side` met plus de 3 s, donc `fallback()` lit ses deux fichiers de cache et utilise le plus récemment écrit. `usage.Table()` dessine les lignes et `usage.Best()` choisit le compte nommé sur la dernière ligne.

Pour un compte, `Fetcher.Fetch` effectue ces vérifications dans l'ordre :

<!-- motion:usage-cache -->

1. Avec `--source cache`, il passe directement à `fallback()`.
2. Sauf si `--no-cache` est défini, un résultat en direct dans `~/.local/state/ccam/usage/<name>.json` datant de moins de 60 secondes (`CacheTTL`) est renvoyé tel quel. Lancer `ccam usage` deux fois de suite coûte une seule requête par compte.
3. `readCredentials()` lit `.credentials.json` (ou le trousseau macOS) sans le modifier, et seule la clé `claudeAiOauth` est décodée. L'absence de connexion par abonnement donne `not logged in`, ou `no limits (API)` quand `.claude.json` contient une clé d'API.
4. Un jeton d'accès expiré n'est pas utilisé. Sans `--refresh`, le compte affiche ses chiffres enregistrés, marqués comme expirés.
5. Sinon, ccam appelle le point de terminaison. Un 401 ou un 403 compte comme une connexion expirée, tout autre statut ou un corps illisible compte comme une erreur, et les deux se rabattent sur les chiffres enregistrés.
6. Une réponse valide est écrite dans `<name>.json` avec l'heure de la mesure.

`fallback()` lit à la fois `<name>.json` et `<name>.statusline.json` et garde le plus récent. Avant d'afficher des chiffres enregistrés, il vérifie l'heure de réinitialisation de chaque fenêtre : une fenêtre réinitialisée depuis l'enregistrement des chiffres affiche 0 %, car l'ancien pourcentage n'est plus valable. C'est pourquoi `personal` ci-dessus affiche `0%` sans heure de réinitialisation.

Le point de terminaison a deux formats de réponse. Le plus récent a un tableau `limits` avec des entrées `session` (5 heures), `weekly_all` et `weekly_scoped` (par modèle). L'ancien a des objets `five_hour`, `seven_day`, `seven_day_opus` et `seven_day_sonnet`. `ParseLive` lit les deux : le tableau `limits` l'emporte, et les anciens objets complètent ce qu'il ne couvre pas. Un corps sans aucune limite est signalé comme un changement de format, dont `ccam doctor` vous informe.

### Renouveler les jetons

Les jetons de renouvellement tournent : en renouveler un vous donne un nouveau jeton de renouvellement et invalide l'ancien. Si ccam renouvelait la connexion pendant que Claude Code s'exécute avec le même compte, cette session serait déconnectée. `ccam usage` lit donc les jetons et ne les écrit jamais, sauf si vous passez `--refresh`. `Fetcher.refresh` alors :

1. refuse les comptes enregistrés sur place avec `ccam import`, puisque `claude` seul peut utiliser cette connexion à l'insu de ccam ;
2. refuse les connexions stockées dans le trousseau macOS ;
3. prend un fichier de verrou, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (sous `~/.local/state/ccam/locks/` quand cette variable n'est pas définie), afin que deux processus ccam ne renouvellent jamais deux fois le même jeton ;
4. sous le verrou, refuse si une session du compte est en cours, et saute la requête si un autre processus a déjà renouvelé le jeton ;
5. envoie le jeton de renouvellement au point de terminaison des jetons et réécrit les nouveaux jetons dans `.credentials.json` avec `fsutil.AtomicWrite`, en conservant toutes les autres clés du fichier.

### Le collecteur de la ligne d'état

Claude Code exécute la commande de ligne d'état et lui transmet un document JSON par un tube. ccam y lit le nom du modèle et `rate_limits` :

<!-- motion:usage-statusline -->

1. `runStatusline` trouve le compte à partir de `CCAM_PROFILE`, que ccam définit quand il lance Claude Code, ou sinon en comparant `CLAUDE_CONFIG_DIR` aux dossiers des comptes.
2. `ParseStatusLine` lit `rate_limits.five_hour` et `rate_limits.seven_day`. Le JSON ne contient pas de limites pour les comptes à clé d'API, ni avant la première réponse de la session ; rien n'est alors enregistré.
3. `SaveStatusLine` écrit `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` l'utilise quand le point de terminaison ne répond pas.
4. ccam affiche `[name] · model · 5h N% · 7d N%`. Avec `--wrap`, il exécute votre commande via `sh -c` avec la même entrée et affiche à la place `[name]` suivi de la première ligne de la commande.

La commande enveloppée dispose de 2 secondes. Au-delà, elle est tuée et la ligne n'affiche que le nom du compte, si bien qu'une commande lente ne peut pas bloquer la ligne d'état de Claude Code.

`ccam statusline install` écrit `statusLine` dans le `settings.json` de chaque compte. Quand un compte a déjà une commande de ligne d'état, ccam la conserve sous la forme `--wrap '<old command>'`. Un `settings.json` partagé par un groupe est écrit une seule fois, via son lien, si bien que chaque compte du groupe le reçoit. Relancer la commande affiche `(already set)` pour les fichiers qui utilisent déjà ccam.

## Sécurité et cas particuliers

- `ccam usage` ne modifie jamais un fichier de jetons sauf si vous passez `--refresh`. Quand une connexion a expiré, lancez le compte une fois avec `ccam <name>` et Claude Code la renouvellera.
- `--refresh` ne touche jamais aux comptes en cours d'exécution, aux comptes importés ni aux connexions du trousseau. Un renouvellement refusé affiche les chiffres enregistrés et la raison.
- Le point de terminaison d'utilisation n'est pas documenté publiquement et peut changer. Dans ce cas, vous voyez toujours les chiffres enregistrés, et `ccam doctor` signale le changement de format.
- Sans chiffres enregistrés ni réponse, la ligne affiche `lookup failed:` et la raison, par exemple `no saved usage` avec `--source cache`.
- Les fichiers de cache ne contiennent que des pourcentages et des heures de réinitialisation, jamais de jetons.

## Voir aussi

- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Passer au compte suivant](/features/next)
- [Gérer les limites d'utilisation au quotidien](/guides/usage-limit-workflow)
- [Commandes](/docs/commands)
