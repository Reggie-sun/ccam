---
title: "Limites d'utilisation"
description: "ccam usage affiche d'un coup les limites sur 5 heures et hebdomadaires de chaque compte, et ccam statusline met le compte et l'utilisation dans Claude Code."
updated: 2026-10-25
order: 7
translation: machine
sourceHash: "c98993fcdc7d"
---

`ccam usage` indique quelle part de la limite sur 5 heures, de la limite hebdomadaire et de la limite hebdomadaire propre à chaque modèle (comme Fable) chaque compte a utilisée, avec une barre par limite et le temps restant avant sa réinitialisation. Utilisez-le avant de commencer une longue tâche, ou quand Claude Code vous signale qu'une limite approche, pour voir quel compte a encore de la marge. `ccam statusline` est la commande de ligne d'état pour Claude Code : elle affiche le nom du compte et l'utilisation dans Claude Code et enregistre les chiffres pour `ccam usage`. Pour les jetons et les coûts estimés derrière cette utilisation, `ccam ccusage` exécute [ccusage](https://ccusage.com) sur les conversations de tous les comptes ([voir plus bas](#tokens-and-costs-ccam-ccusage)).

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
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
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
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Chaque compte a une ligne avec son numéro, son nom et l'ancienneté de ses chiffres ; `(cached)` signifie qu'ils proviennent d'un fichier enregistré et non du point de terminaison. En dessous, une ligne par limite affiche une barre de 20 cases, le pourcentage et le temps restant avant la réinitialisation de la fenêtre : `5-hour`, `Weekly` (tous les modèles), puis la fenêtre hebdomadaire propre à chaque modèle, comme `Fable`. Une fenêtre que le compte n'a pas affiche `-`. Un pourcentage passe au jaune à partir de 80 % et au rouge à 100 %. Une fenêtre que Claude Code a signalée pleine, quand il a refusé une requête à une limite d'utilisation, affiche 100 % jusqu'à sa réinitialisation, quels que soient les chiffres enregistrés (voir [le relevé de limite](#the-limit-record)). Dans un terminal, tous les comptes apparaissent immédiatement sous la forme `fetching…`, et leurs lignes sont dessinées à mesure qu'ils répondent. La dernière ligne nomme le compte que [`ccam next --auto`](/features/next) choisirait selon les mêmes règles (`auto_max_usage`, `auto_strategy`, les sessions en cours, uniquement les comptes de l'organisation du dernier compte utilisé sauf si `cross_org` est défini, toute fenêtre à 100 % écartée : un compte dont la semaine est pleine n'est donc jamais nommé), et n'apparaît que si vous avez plus d'un compte. Avec la stratégie `headroom` par défaut, elle indique `→ Most headroom: <name>` ; avec `expiring` ou `balance`, elle donne la raison, par exemple `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Sièges Team et Enterprise.** Le choix automatique (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord, et [`auto_switch`](/features/next#auto-switch) lorsque vous l'activez) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Utiliser le choix automatique entre plusieurs sièges peut donc enfreindre les conditions d'Anthropic. C'est à vous de décider si vous l'utilisez dans ce cas : renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier. Les offres personnelles (Pro, Max) sont traitées dans la [note de la FAQ sur les conditions](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Un compte non connecté affiche `not logged in → ccam login <name>`, et un compte à clé d'API affiche `no limits (API)`. Un jeton d'accès expiré affiche ``access token expired; `ccam <name>` renews it when it starts``, au-dessus des chiffres enregistrés s'il y en a : Claude Code le renouvelle au prochain démarrage. Une connexion qui a expiré dans son ensemble, parce que son jeton de renouvellement a expiré ou que Claude Code l'a effacée après un refus du serveur, affiche `` login expired → `ccam login <name>` ``.

Après `ccam statusline install`, Claude Code affiche une ligne comme `[work] · Opus · 5h 24% · 7d 41%`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--json` | Afficher les rapports en JSON : `name`, `slot`, `status`, `limits` (avec `credits` quand le compte a des crédits d'utilisation activés), `source`, `at`, `error`, et `limited` et `limited_windows` quand un relevé de limite met des fenêtres à 100 %. Sans aucun compte, il affiche `[]` | désactivé |
| `--watch[=N]` | Redessiner toutes les `N` secondes jusqu'à ce que vous appuyiez sur Ctrl+C. `--watch` seul, ou une valeur inférieure ou égale à 0, signifie 60. La valeur se place après `=` : `--watch 10` est refusé avec une indication | désactivé |
| `-v`, `--verbose` | Ajouter sous chaque compte une ligne `source:` : la provenance de ses chiffres (le point de terminaison d'utilisation, Claude Code, la ligne d'état, un fichier enregistré, ou la limite d'utilisation signalée par Claude Code). Quand un relevé de limite passe par-dessus les chiffres, une ligne comme `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` suit | désactivé |
| `--no-cache` | Ignorer le cache de 60 secondes et interroger à nouveau le point de terminaison | désactivé |
| `--source live\|claude\|cache` | `live` interroge le point de terminaison avec le jeton du compte ; `claude` fait faire la consultation par Claude Code (voir ci-dessous) ; `cache` n'affiche que les chiffres enregistrés | `usage_source`, sinon `live` |
| `--refresh` | Renouveler les jetons d'accès expirés, uniquement pour les comptes qui ne sont pas en cours d'exécution. Un renouvellement peut prendre jusqu'à 30 secondes | désactivé |
| `--dry-run` | Accepté : `ccam usage` ne modifie rien d'autre que son propre cache. Refusé avec `--refresh` : `--dry-run cannot be combined with --refresh, which rewrites login tokens` | désactivé |

Toute autre valeur de `--source` échoue avec `--source must be live, claude or cache`.

### L'utilisation vue par Claude Code lui-même

Avec `usage_source = "claude"` en haut de `config.toml` (au-dessus du premier `[[group]]`, `[[profile]]` ou `[[folder]]`), ou `--source claude`, ccam ne lit pas du tout le jeton de connexion. Pour chaque compte, il lance `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` dans le dossier du compte et envoie la requête `get_usage` de Claude Code sur l'entrée standard (Claude Code 2.1.270 ou plus récent) ; Claude Code consulte l'utilisation avec sa propre connexion et répond avec les fenêtres sur 5 heures, hebdomadaire et par modèle. `--safe-mode` empêche les hooks, les serveurs MCP et les plugins du compte de démarrer, et aucune session n'est enregistrée. Chaque compte demande un démarrage de Claude Code, environ une seconde : cette source est donc plus lente que `live`. Les comptes sont interrogés en même temps, chacun avec 20 secondes au plus, et le cache de 60 secondes s'applique comme d'habitude. Une connexion dont le jeton d'accès a expiré, ou expire dans les 5 minutes, dispose de 50 secondes, puisque Claude Code le renouvelle au démarrage et qu'un renouvellement interrompu peut faire perdre la connexion. Une consultation abandonnée demande à Claude Code de se terminer avec SIGTERM, pour qu'il achève un renouvellement commencé (Windows n'a pas ce signal ; il y est arrêté). Une connexion qui a expiré dans son ensemble affiche `` login expired → `ccam login <name>` `` sans lancer Claude Code. Claude Code marque cette requête comme expérimentale : s'il cesse d'y répondre, ccam affiche les chiffres enregistrés et `ccam doctor` le signale. Quand Claude Code répond que sa propre consultation a échoué, ccam affiche les chiffres enregistrés avec `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, et non une offre sans limites. Un `--dry-run` ne lance pas Claude Code (qui écrit dans le dossier du compte) et affiche les chiffres enregistrés. `usage_source` accepte `live` ou `claude` ; `cache` n'est qu'un choix de `--source`. `ccam usage` n'accepte aucun argument.

`ccam statusline` lit sur l'entrée standard le JSON de ligne d'état de Claude Code. Sa seule option est `--wrap <command>` : une commande de ligne d'état existante dont la première ligne de sortie reçoit le nom du compte en tête. `ccam statusline install` n'accepte aucune option. Les deux refusent `--dry-run` avec `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Jetons et coûts : `ccam ccusage`

`ccam usage` indique quelle part de chaque limite un compte a utilisée. [ccusage](https://ccusage.com), lui, compte les jetons et les coûts estimés dans les journaux de conversation de Claude Code. Seul, il ne cherche les journaux de Claude Code que dans `~/.claude` et `~/.config/claude` (ou dans les dossiers que désigne `CLAUDE_CONFIG_DIR`) : parmi les comptes de ccam, il en voit donc au plus un, importé sur place dans `~/.claude`. `ccam ccusage` l'exécute avec `CLAUDE_CONFIG_DIR` pointant vers les dossiers de tous les comptes qui ont des conversations, de sorte qu'un seul rapport les couvre tous :

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam lit d'abord ses propres options : `-a` (`--account`) prend des noms ou des numéros de compte, séparés par des virgules ou répétés, et `--dry-run` affiche la commande au lieu de l'exécuter. À partir du premier argument qui n'est pas une option de ccam, tout est transmis tel quel à ccusage ; `ccam ccusage -- --help` affiche l'aide de ccusage lui-même. ccam exécute un `ccusage` installé, sinon `bun x ccusage`, sinon `npx ccusage@latest` (ces deux derniers le récupèrent depuis npm) ; si aucun n'est présent, il indique d'installer Node.js, Bun ou ccusage. Sous Linux et macOS, ccusage prend la place de ccam ; sous Windows, ccam l'attend et se termine avec le code de sortie de ccusage. Sous WSL, ccam ignore un `ccusage`, `bun` ou `npx` situé sur un lecteur Windows, qui s'exécuterait comme un programme Windows sans les dossiers des comptes.

- Le rapport couvre les comptes dont le dossier contient des conversations (un dossier `projects`). Une conversation que des comptes partagent via un groupe n'est comptée qu'une fois.
- Avec `-a`, un compte qui partage son historique via un groupe apporte les conversations de tout le groupe, et ccam indique à quels comptes elles appartiennent : `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` affiche des fenêtres de 5 heures, et chaque compte a les siennes : nommez donc un seul compte pour ce rapport, comme ci-dessus.
- ccusage inclut aussi les autres outils de code qu'il trouve, comme Codex, sauf si vous en nommez un : `ccam ccusage claude daily` ne porte que sur Claude Code.
- Quand ccam le lance, ccusage ne lit ses paramètres (`ccusage.json`) que dans `./.ccusage` et dans les dossiers des comptes. `--config` désigne un autre fichier : `ccam ccusage daily --config ~/.config/claude/ccusage.json`.
- ccusage lirait un dossier dont le chemin contient une virgule comme deux dossiers : ccam le refuse donc et nomme le compte ; laissez ce compte de côté avec `-a`.
- `ccusage` est une commande de ccam : un compte nommé ccusage se lance donc avec `ccam run ccusage` ; `ccam doctor` le signale.

## Fonctionnement

Les chiffres de `ccam usage` viennent de deux sources : le point de terminaison d'utilisation non documenté d'Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, appelé avec la propre connexion de chaque compte, et les valeurs que Claude Code transmet à sa commande de ligne d'état. Tout se trouve dans `internal/usage`.

1. `cli.usageCmd` charge `config.toml` et construit un `usage.Account` par compte : son dossier `~/.local/share/ccam/profiles/<name>/`, ses marqueurs de session (avec les propres enregistrements de Claude Code dans le dossier du compte, ils indiquent s'il est en cours d'exécution) et s'il a été importé sur place.
2. `Fetcher.FetchAll` lance une goroutine par compte. Toutes partagent un même `http.Client`, si bien que les connexions au point de terminaison sont réutilisées. La connexion peut prendre 2 s, et chaque compte dispose de 3 s au total. Le client ne suit pas les redirections, puisqu'une redirection renverrait le jeton ailleurs, et lit au plus 1 Mo de la réponse.
3. Dès qu'un compte a terminé, sa ligne est redessinée. Un compte qui échoue ou ne répond pas à temps se rabat sur les chiffres enregistrés, si bien qu'un compte lent ne retarde jamais les autres.

Chaque compte interroge le point de terminaison dans sa propre goroutine ; celui qui dépasse le délai affiche ses chiffres enregistrés à la place.

<!-- motion:usage-fetch -->

`personal` et `work` répondent à temps et affichent des chiffres en direct. `side` met plus de 3 s, donc `fallback()` affiche à la place ses chiffres enregistrés. `usage.Table()` dessine un bloc par compte et `usage.Best()` choisit le compte nommé sur la dernière ligne.

Pour un compte, `Fetcher.Fetch` effectue ces vérifications dans l'ordre :

<!-- motion:usage-cache -->

1. Avec `--source cache`, il passe directement à `fallback()`.
2. Sauf si `--no-cache` est défini, un résultat en direct dans `~/.local/state/ccam/usage/<name>.json` mesuré il y a moins de 60 secondes (`CacheTTL`) est utilisé, évalué à l'instant présent : une fenêtre réinitialisée depuis affiche 0 %, et un fichier daté en avance sur l'horloge ne compte pas comme récent. Lancer `ccam usage` deux fois de suite coûte une seule requête par compte.
3. `readCredentials()` lit `.credentials.json` (ou le trousseau macOS) sans le modifier, et seule la clé `claudeAiOauth` est décodée. L'absence de connexion par abonnement donne `not logged in`, ou `no limits (API)` quand `.claude.json` contient une clé d'API. Une connexion dont le jeton de renouvellement a expiré, ou que Claude Code a effacée, donne `` login expired → `ccam login <name>` `` et n'est pas consultée.
4. Un jeton d'accès expiré n'est pas utilisé. Sans `--refresh`, le compte affiche ses chiffres enregistrés avec ``access token expired; `ccam <name>` renews it when it starts``.
5. Sinon, ccam appelle le point de terminaison. Un 401 ou un 403 compte comme une connexion expirée (`login expired: authentication failed (401)`), tout autre statut ou un corps illisible compte comme une erreur, et les deux se rabattent sur les chiffres enregistrés. Une consultation interrompue par Ctrl+C indique `cancelled`. Un 429 signifie que le point de terminaison est sollicité trop souvent : ccam écrit `<name>.pause` avec l'heure à laquelle redemander (`Retry-After`, ou 5 minutes, une heure au plus), affiche les chiffres enregistrés avec `rate limited by the usage endpoint; asking again after 15:10` et n'appelle plus le point de terminaison pour ce compte d'ici là, même avec `--no-cache`. Réessayer tout de suite ne ferait que prolonger la limite.
6. Une réponse valide est écrite dans `<name>.json` avec l'heure de la mesure.

`fallback()` prend la consultation enregistrée la plus récente (`<name>.json`, ou `<name>.claude.json` pour les consultations de Claude Code) et, quand `<name>.statusline.json` a été écrit après elle, y superpose les fenêtres sur 5 heures et hebdomadaire de la ligne d'état ; les fenêtres par modèle, que la ligne d'état n'a pas, restent telles que la consultation les a enregistrées. Avant d'afficher des chiffres enregistrés, il vérifie l'heure de réinitialisation de chaque fenêtre : une fenêtre réinitialisée depuis l'enregistrement des chiffres affiche 0 %, car l'ancien pourcentage n'est plus valable. C'est pourquoi `personal` ci-dessus affiche `0%` sans heure de réinitialisation.

Le point de terminaison a deux formats de réponse. Le plus récent a un tableau `limits` avec des entrées `session` (5 heures), `weekly_all` et `weekly_scoped` (par modèle). L'ancien a des objets `five_hour`, `seven_day`, `seven_day_opus` et `seven_day_sonnet`. `ParseLive` lit les deux : le tableau `limits` l'emporte, et les anciens objets complètent ce qu'il ne couvre pas. Un corps sans aucune limite est signalé comme un changement de format, dont `ccam doctor` vous informe.

### Le relevé de limite

Quand Claude Code arrête un tour à une limite d'utilisation, le hook du [plugin ccam](/docs/install#claude-code-plugin) lit la requête refusée dans la conversation : quelle fenêtre était pleine (5 heures, hebdomadaire, Opus ou Sonnet) et quand elle se réinitialise. Il l'écrit dans `~/.local/state/ccam/usage/<name>.limit.json`. `Fetch` superpose ce relevé à tout ce qu'il sait par ailleurs, quelle qu'en soit la source : chaque fenêtre qu'il contient compte pour 100 % jusqu'à sa réinitialisation, aussi anciens ou bas que soient les chiffres enregistrés, si bien que `ccam next --auto`, `ccam pick`, `--fallback` et [`auto_switch`](/features/next#auto-switch) ne choisissent pas de nouveau ce compte d'ici là. Quand rien d'autre n'est connu, le relevé seul est ce que le compte affiche (`source: the usage limit Claude Code reported` avec `-v`). Les sessions devant lesquelles personne n'est assis (`-p`, sessions en arrière-plan) n'écrivent pas de relevé, pas plus qu'une suspension par l'organisation (plus d'utilisation disponible, un groupe réglé à 0 $), que ses administrateurs peuvent lever à tout moment. `ccam login` vers un autre compte Claude supprime le relevé ; supprimer le fichier l'efface aussi.

### Renouveler les jetons

Les jetons de renouvellement tournent : en renouveler un vous donne un nouveau jeton de renouvellement et invalide l'ancien. Si ccam renouvelait la connexion pendant que Claude Code s'exécute avec le même compte, cette session serait déconnectée. `ccam usage` lit donc les jetons et ne les écrit jamais, sauf si vous passez `--refresh`. `Fetcher.refresh` alors :

1. refuse les comptes enregistrés sur place avec `ccam import`, puisque `claude` seul peut utiliser cette connexion à l'insu de ccam ;
2. refuse les connexions stockées dans le trousseau macOS ;
3. prend un fichier de verrou, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, afin que deux processus ccam ne renouvellent jamais deux fois le même jeton ; tant qu'il est détenu, le compte affiche `token not renewed: another ccam is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` et `logout` le prennent aussi, et refusent pendant un renouvellement ;
4. sous le verrou, refuse si une session du compte est en cours, quelle que soit la façon dont elle a été lancée (via ccam, `claude` seul, un IDE, ou une session en arrière-plan et le démon derrière elle, qui renouvelle lui-même la connexion), et saute la requête si un autre processus a déjà renouvelé le jeton ;
5. n'envoie rien tant que Claude Code détient ses propres verrous de renouvellement dans le dossier du compte (`.oauth_refresh.lock`, `<folder>.lock`) : Claude Code est alors en train de renouveler cette connexion ;
6. envoie le jeton de renouvellement au point de terminaison des jetons et attend la réponse jusqu'à 30 secondes, le temps que Claude Code accorde à ses propres renouvellements. Une requête envoyée n'est jamais interrompue, pas même par Ctrl+C, car le serveur a peut-être déjà fait tourner le jeton de renouvellement et seule sa réponse conserve encore la connexion : Ctrl+C ou un raccrochage met fin aux autres consultations, et ccam quitte (130 après Ctrl+C) une fois la réponse enregistrée ;
7. enregistre sous le verrou d'écriture de Claude Code (`.storage-write.lock`) et ne modifie que les champs de jeton de `.credentials.json` tel qu'il est à ce moment : ce que Claude Code a écrit entre-temps (une connexion MCP) reste, une connexion que Claude Code a renouvelée entre-temps est conservée, et une déconnexion n'est jamais annulée. Quand la réponse donne la durée de vie du nouveau jeton de renouvellement, elle est enregistrée comme `refreshTokenExpiresAt`, comme le fait Claude Code.

Un renouvellement qui n'aboutit pas indique pourquoi, sur la ligne du compte ou au-dessus de ses chiffres enregistrés : `token not renewed: a session is running`, par exemple, ou `login may be lost: …` quand une requête a été envoyée mais que sa réponse n'a pas pu être utilisée ou enregistrée, avec le conseil d'exécuter `ccam login <name>` si Claude Code vous demande ensuite de vous connecter.

### Le collecteur de la ligne d'état

Claude Code exécute la commande de ligne d'état et lui transmet un document JSON par un tube. ccam y lit le nom du modèle et `rate_limits` :

<!-- motion:usage-statusline -->

1. `runStatusline` trouve le compte à partir de `CCAM_PROFILE`, que ccam définit quand il lance Claude Code, ou sinon en comparant `CLAUDE_CONFIG_DIR` aux dossiers des comptes.
2. `ParseStatusLine` lit `rate_limits.five_hour` et `rate_limits.seven_day`. Le JSON ne contient pas de limites pour les comptes à clé d'API, ni avant la première réponse de la session ; rien n'est alors enregistré.
3. `SaveStatusLine` écrit `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` l'utilise quand le point de terminaison ne répond pas.
4. ccam affiche `[name] · model · 5h N% · 7d N%`. Avec `--wrap`, il exécute votre commande via `sh -c` avec la même entrée et affiche sa sortie comme Claude Code affiche une ligne d'état : chaque ligne non vide, avec `[name]` devant la première. Une commande qui échoue n'affiche que le nom du compte. Sous Windows, il exécute la commande comme le fait Claude Code : avec Git Bash s'il est installé, sinon avec PowerShell.
5. Quand l'utilisation sur 5 heures ou hebdomadaire atteint 90 % ou plus et qu'il existe un autre compte, la ligne se termine par `→ ccam next --auto --continue`, par exemple `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

La commande enveloppée n'a pas de limite de temps propre : Claude Code accorde 10 minutes à une commande de ligne d'état et arrête une exécution quand une mise à jour plus récente commence, si bien qu'une commande lente s'affiche comme elle le ferait sans ccam. Quand Claude Code met fin à ccam, la commande se termine aussi (sous Linux et macOS, avec chaque processus qu'elle a lancé).

`ccam statusline install` écrit `statusLine` dans le `settings.json` de chaque compte et conserve les autres clés de l'objet, comme `padding`, `refreshInterval` et `hideVimModeIndicator`. Quand un compte a déjà une commande de ligne d'état, ccam la conserve sous la forme `--wrap '<old command>'` ; sous Windows sous la forme `--wrap-base64 <base64>`, qui n'a besoin de guillemets dans aucun des deux shells. Une ligne d'état qui lance déjà un ccam, depuis un autre chemin ou entre guillemets à la façon d'un ancien ccam, est remplacée au lieu d'être enveloppée, en conservant ce qu'elle enveloppe. Sous Windows, le chemin de ccam est écrit avec des barres obliques et `~` pour le dossier personnel (`~/.local/bin/ccam.exe statusline`), ce que Git Bash et PowerShell comprennent tous les deux. Un chemin qui demande des guillemets est écrit en entier pour le shell qu'utilise Claude Code : `'…'` pour Git Bash, `& '…'` pour PowerShell. Un `settings.json` partagé par un groupe est écrit une seule fois, via son lien, si bien que chaque compte du groupe le reçoit. Relancer la commande affiche `(already set)` pour les fichiers qui utilisent déjà ccam.

## Sécurité et cas particuliers

- `ccam usage` ne modifie jamais un fichier de jetons sauf si vous passez `--refresh`. Quand un jeton d'accès a expiré, lancez le compte une fois avec `ccam <name>` et Claude Code le renouvellera. Une connexion qui a expiré dans son ensemble demande `ccam login <name>`.
- `--refresh` ne touche jamais aux comptes en cours d'exécution (y compris quand Claude Code a été lancé sans ccam), aux comptes importés ni aux connexions du trousseau. Un renouvellement refusé affiche les chiffres enregistrés et la raison (`token not renewed: …`).
- Le point de terminaison d'utilisation n'est pas documenté publiquement et peut changer. Dans ce cas, vous voyez toujours les chiffres enregistrés, et `ccam doctor` signale le changement de format.
- Sans chiffres enregistrés ni réponse, la ligne du compte affiche `lookup failed:` et la raison, par exemple `no saved usage` avec `--source cache`.
- Les fichiers de cache ne contiennent que des pourcentages et des heures de réinitialisation, jamais de jetons.

## Voir aussi

- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Passer au compte suivant](/features/next)
- [Gérer les limites d'utilisation au quotidien](/guides/usage-limit-workflow)
- [Commandes](/docs/commands)
