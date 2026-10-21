---
title: "Limites d'utilisation"
description: "ccam usage affiche d'un coup les limites sur 5 heures et hebdomadaires de chaque compte, et ccam statusline met le compte et l'utilisation dans Claude Code."
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

`ccam usage` indique quelle part de la limite sur 5 heures, de la limite hebdomadaire et de la limite hebdomadaire propre à chaque modèle (comme Fable) chaque compte a utilisée, avec une barre par limite et le temps restant avant sa réinitialisation. Utilisez-le avant de commencer une longue tâche, ou quand Claude Code vous signale qu'une limite approche, pour voir quel compte a encore de la marge. `ccam statusline` est la commande de ligne d'état pour Claude Code : elle affiche le nom du compte et l'utilisation dans Claude Code et enregistre les chiffres pour `ccam usage`.

## Utilisation

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Chaque compte a une ligne avec son numéro, son nom et l'ancienneté de ses chiffres ; `(cached)` signifie qu'ils proviennent d'un fichier enregistré et non du point de terminaison. En dessous, une ligne par limite affiche une barre de 20 cases, le pourcentage et le temps restant avant la réinitialisation de la fenêtre : `5-hour`, `Weekly` (tous les modèles), puis la fenêtre hebdomadaire propre à chaque modèle, comme `Fable`. Une fenêtre que le compte n'a pas affiche `-`. Un pourcentage passe au jaune à partir de 80 % et au rouge à 100 %. Dans un terminal, tous les comptes apparaissent immédiatement sous la forme `fetching…`, et leurs lignes sont dessinées à mesure qu'ils répondent. La dernière ligne nomme le compte que [`ccam next --auto`](/features/next) choisirait selon les mêmes règles (`auto_max_usage`, `auto_strategy`, les sessions en cours, uniquement les comptes de l'organisation du dernier compte utilisé sauf si `cross_org` est défini, toute fenêtre à 100 % écartée : un compte dont la semaine est pleine n'est donc jamais nommé), et n'apparaît que si vous avez plus d'un compte. Un compte non connecté affiche `not logged in → ccam login <name>`, et un compte à clé d'API affiche `no limits (API)`.

Après `ccam statusline install`, Claude Code affiche une ligne comme `[work] · Opus · 5h 24% · 7d 41%`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--json` | Afficher les rapports en JSON : `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | désactivé |
| `--watch[=N]` | Redessiner toutes les `N` secondes jusqu'à ce que vous appuyiez sur Ctrl+C. `--watch` seul, ou une valeur inférieure ou égale à 0, signifie 60 | désactivé |
| `-v`, `--verbose` | Ajouter sous chaque compte une ligne `source:` : la provenance de ses chiffres (le point de terminaison d'utilisation, Claude Code, la ligne d'état ou un fichier enregistré) | désactivé |
| `--no-cache` | Ignorer le cache de 60 secondes et interroger à nouveau le point de terminaison | désactivé |
| `--source live\|claude\|cache` | `live` interroge le point de terminaison avec le jeton du compte ; `claude` fait faire la consultation par Claude Code (voir ci-dessous) ; `cache` n'affiche que les chiffres enregistrés | `usage_source`, sinon `live` |
| `--refresh` | Renouveler les jetons d'accès expirés, uniquement pour les comptes qui ne sont pas en cours d'exécution | désactivé |
| `--dry-run` | Accepté : `ccam usage` ne modifie rien d'autre que son propre cache. Refusé avec `--refresh` : `--dry-run cannot be combined with --refresh, which rewrites login tokens` | désactivé |

Toute autre valeur de `--source` échoue avec `--source must be live, claude or cache`.

### L'utilisation vue par Claude Code lui-même

Avec `usage_source = "claude"` dans `config.toml` (ou `--source claude`), ccam ne lit pas du tout le jeton de connexion. Pour chaque compte, il lance `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` dans le dossier du compte et envoie la requête `get_usage` de Claude Code sur l'entrée standard (Claude Code 2.1.270 ou plus récent) ; Claude Code consulte l'utilisation avec sa propre connexion et répond avec les fenêtres sur 5 heures, hebdomadaire et par modèle. `--safe-mode` empêche les hooks, les serveurs MCP et les plugins du compte de démarrer, et aucune session n'est enregistrée. Chaque compte demande un démarrage de Claude Code, environ une seconde : cette source est donc plus lente que `live`. Les comptes sont interrogés en même temps, chacun avec 20 secondes au plus, et le cache de 60 secondes s'applique comme d'habitude. Claude Code marque cette requête comme expérimentale : s'il cesse d'y répondre, ccam affiche les chiffres enregistrés et `ccam doctor` le signale. Un `--dry-run` ne lance pas Claude Code (qui écrit dans le dossier du compte) et affiche les chiffres enregistrés. `usage_source` accepte `live` ou `claude` ; `cache` n'est qu'un choix de `--source`. `ccam usage` n'accepte aucun argument.

`ccam statusline` lit sur l'entrée standard le JSON de ligne d'état de Claude Code. Sa seule option est `--wrap <command>` : une commande de ligne d'état existante dont la première ligne de sortie reçoit le nom du compte en tête. `ccam statusline install` n'accepte aucune option. Les deux refusent `--dry-run` avec `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Fonctionnement

Les chiffres viennent de deux sources : le point de terminaison d'utilisation non documenté d'Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, appelé avec la propre connexion de chaque compte, et les valeurs que Claude Code transmet à sa commande de ligne d'état. Tout se trouve dans `internal/usage`.

1. `cli.usageCmd` charge `config.toml` et construit un `usage.Account` par compte : son dossier `~/.local/share/ccam/profiles/<name>/`, son dossier de session (pour savoir s'il est en cours d'exécution) et s'il a été importé sur place.
2. `Fetcher.FetchAll` lance une goroutine par compte. Toutes partagent un même `http.Client`, si bien que les connexions au point de terminaison sont réutilisées. La connexion peut prendre 2 s, et chaque compte dispose de 3 s au total. Le client ne suit pas les redirections, puisqu'une redirection renverrait le jeton ailleurs, et lit au plus 1 Mo de la réponse.
3. Dès qu'un compte a terminé, sa ligne est redessinée. Un compte qui échoue ou ne répond pas à temps se rabat sur les chiffres enregistrés, si bien qu'un compte lent ne retarde jamais les autres.

Chaque compte interroge le point de terminaison dans sa propre goroutine ; celui qui dépasse le délai affiche ses chiffres enregistrés à la place.

<!-- motion:usage-fetch -->

`personal` et `work` répondent à temps et affichent des chiffres en direct. `side` met plus de 3 s, donc `fallback()` lit ses deux fichiers de cache et utilise le plus récemment écrit. `usage.Table()` dessine un bloc par compte et `usage.Best()` choisit le compte nommé sur la dernière ligne.

Pour un compte, `Fetcher.Fetch` effectue ces vérifications dans l'ordre :

<!-- motion:usage-cache -->

1. Avec `--source cache`, il passe directement à `fallback()`.
2. Sauf si `--no-cache` est défini, un résultat en direct dans `~/.local/state/ccam/usage/<name>.json` datant de moins de 60 secondes (`CacheTTL`) est renvoyé tel quel. Lancer `ccam usage` deux fois de suite coûte une seule requête par compte.
3. `readCredentials()` lit `.credentials.json` (ou le trousseau macOS) sans le modifier, et seule la clé `claudeAiOauth` est décodée. L'absence de connexion par abonnement donne `not logged in`, ou `no limits (API)` quand `.claude.json` contient une clé d'API.
4. Un jeton d'accès expiré n'est pas utilisé. Sans `--refresh`, le compte affiche ses chiffres enregistrés, marqués comme expirés.
5. Sinon, ccam appelle le point de terminaison. Un 401 ou un 403 compte comme une connexion expirée, tout autre statut ou un corps illisible compte comme une erreur, et les deux se rabattent sur les chiffres enregistrés. Un 429 signifie que le point de terminaison est sollicité trop souvent : ccam écrit `<name>.pause` avec l'heure à laquelle redemander (`Retry-After`, ou 5 minutes, une heure au plus), affiche les chiffres enregistrés avec `rate limited by the usage endpoint; asking again after 15:10` et n'appelle plus le point de terminaison pour ce compte d'ici là, même avec `--no-cache`. Réessayer tout de suite ne ferait que prolonger la limite.
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
4. ccam affiche `[name] · model · 5h N% · 7d N%`. Avec `--wrap`, il exécute votre commande via `sh -c` avec la même entrée et affiche à la place `[name]` suivi de la première ligne de la commande. Sous Windows, il exécute la commande comme le fait Claude Code : avec Git Bash s'il est installé, sinon avec PowerShell.
5. Quand l'utilisation sur 5 heures ou hebdomadaire atteint 90 % ou plus et qu'il existe un autre compte, la ligne se termine par `→ ccam next --auto --continue`, par exemple `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

La commande enveloppée dispose de 2 secondes. Au-delà, elle est tuée et la ligne n'affiche que le nom du compte, si bien qu'une commande lente ne peut pas bloquer la ligne d'état de Claude Code.

`ccam statusline install` écrit `statusLine` dans le `settings.json` de chaque compte. Quand un compte a déjà une commande de ligne d'état, ccam la conserve sous la forme `--wrap '<old command>'`. Sous Windows, le chemin de ccam est écrit avec des barres obliques et `~` pour le dossier personnel (`~/.local/bin/ccam.exe statusline`), ce que Git Bash et PowerShell comprennent tous les deux. Un `settings.json` partagé par un groupe est écrit une seule fois, via son lien, si bien que chaque compte du groupe le reçoit. Relancer la commande affiche `(already set)` pour les fichiers qui utilisent déjà ccam.

## Sécurité et cas particuliers

- `ccam usage` ne modifie jamais un fichier de jetons sauf si vous passez `--refresh`. Quand une connexion a expiré, lancez le compte une fois avec `ccam <name>` et Claude Code la renouvellera.
- `--refresh` ne touche jamais aux comptes en cours d'exécution, aux comptes importés ni aux connexions du trousseau. Un renouvellement refusé affiche les chiffres enregistrés et la raison.
- Le point de terminaison d'utilisation n'est pas documenté publiquement et peut changer. Dans ce cas, vous voyez toujours les chiffres enregistrés, et `ccam doctor` signale le changement de format.
- Sans chiffres enregistrés ni réponse, la ligne du compte affiche `lookup failed:` et la raison, par exemple `no saved usage` avec `--source cache`.
- Les fichiers de cache ne contiennent que des pourcentages et des heures de réinitialisation, jamais de jetons.

## Voir aussi

- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Passer au compte suivant](/features/next)
- [Gérer les limites d'utilisation au quotidien](/guides/usage-limit-workflow)
- [Commandes](/docs/commands)
