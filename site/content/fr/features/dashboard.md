---
title: "Le tableau de bord"
description: "ccam dans un terminal ouvre un tableau de bord plein écran : comptes, groupes, dossiers, extensions et vérifications ; chaque modification montre son plan."
updated: 2026-10-25
order: 15
translation: machine
sourceHash: "98ccf6cd81e9"
---

`ccam` sans argument ouvre un tableau de bord plein écran quand l'entrée et la sortie standard sont un terminal ; sinon, il affiche `ccam list`. Le tableau de bord montre chaque compte avec son utilisation et ses groupes, et fait ce que font les commandes : lancer, ajouter, connecter, changer de groupes, définir des règles de dossier, lister les sessions en arrière-plan, vérifier et réparer. Un lancement, ou le rattachement à une session en arrière-plan, le ferme ; ajouter un compte ou se connecter le quitte le temps de passer au terminal, puis y revient.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Agents   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

La ligne sous les noms des onglets indique en une phrase à quoi sert l'onglet courant. Dans un terminal plus large, l'onglet Accounts affiche aussi l'adresse e-mail de chaque compte.

À chaque ouverture, le tableau de bord demande en arrière-plan à GitHub si une version plus récente de ccam existe ; la première image ne l'attend jamais, et hors ligne rien n'est affiché. S'il y en a une, une ligne sous le titre la recommande :

```console
 ● ccam 0.10.2 is available (you have 0.10.1) · press U to update
```

`U` demande d'abord confirmation, puis met à jour comme le fait [`ccam update`](/features/update) : la version est vérifiée avec ses sommes de contrôle SHA-256 et remplace ccam en une seule étape. Le tableau de bord dans lequel vous êtes reste l'ancienne version : quittez-le et relancez `ccam`. Quand la mise à jour ne peut pas se faire ici, par exemple pour un ccam installé depuis le paquet `.deb`, le tableau de bord affiche le même conseil que `ccam update` et garde la proposition.

La première image n'est dessinée qu'à partir de `config.toml`, en bien moins de 100 ms ; l'utilisation s'affiche au fur et à mesure que chaque compte répond. Chaque limite a sa propre barre : la fenêtre de 5 heures, la fenêtre hebdomadaire de tous les modèles et la fenêtre hebdomadaire propre à Fable ; `-` signifie que le compte n'a pas cette fenêtre. Un pourcentage passe au jaune à partir de 80 % et au rouge à 100 %, et les valeurs tirées d'un cache après un échec de consultation sont affichées en atténué. Une fenêtre que Claude Code a signalée pleine affiche 100 % jusqu'à sa réinitialisation. Un jeton d'accès expiré affiche un `⚠ token expired` atténué (Claude Code le renouvelle au prochain démarrage), tandis que `⚠ login expired` signifie que le compte a besoin de `ccam login`.

## Onglets

| Onglet | Ce qu'il montre | Touches |
| --- | --- | --- |
| Accounts | numéro, nom, adresse e-mail, barres sur 5 heures, hebdomadaire et Fable, dernière utilisation, groupes ; `✔` le compte par défaut, `▸` le compte de ce dossier, `•` en cours d'exécution (y compris Claude Code lancé sans ccam : `claude` seul, un IDE, les sessions en arrière-plan et leur démon) | ci-dessous |
| Agents | ce que liste [`ccam agents`](/features/agents) : les sessions en arrière-plan et les sessions ouvertes de tous les comptes, avec leur état, leur dossier et leur heure de démarrage (ci-dessous) | `enter` se rattacher (comme `ccam attach`), `x` arrêter, `r` recharger |
| Groups | chaque groupe, ce qu'il partage et ses comptes, ainsi que les comptes sans groupe | `enter` modifier, `n` nouveau, `e` renommer, `x` supprimer |
| Folders | les règles [`ccam folder`](/features/folder) et celle qui s'applique ici, y compris quand vous êtes arrivé dans le dossier par un lien symbolique (les nouvelles règles sont enregistrées avec les liens résolus) | `n` utiliser un compte ici, `enter` changer, `x` supprimer |
| Extensions | ce que montre [`ccam ext`](/features/ext), chaque extension une seule fois avec les groupes dont elle provient, `off:` les comptes où un plugin n'est pas activé, les serveurs qu'un compte a pour ce seul dossier parmi les siens, `This folder's project`, et les problèmes qui n'appartiennent à aucun compte | `r` recharger |
| Health | ce que signale [`ccam doctor`](/features/doctor), sous forme d'un encadré par compte (ci-dessous) | `f` réparer l'encadré sélectionné, `F` tout réparer (`doctor --fix`), `↑`/`↓` passer d'un encadré à l'autre |

`←`/`→` ou `Tab` changent d'onglet, `?` liste les touches de l'onglet courant, `r` consulte à nouveau l'utilisation sans le cache court et recharge les onglets Agents, Extensions et Health, et `q` ou `Esc` quitte.

## Touches de l'onglet Accounts

| Touche | Effet |
| --- | --- |
| `Enter` | Lancer Claude Code avec le compte sélectionné |
| `1`–`9` | Aller au compte portant ce numéro ; tapez `1` `2` rapidement pour 12. Seul `Enter` lance |
| `↑`/`↓`, `k`/`j` | Se déplacer |
| `n` | Lancer le compte que choisirait [`ccam next --auto`](/features/next) |
| `o` | Lancer avec des options : `--continue`, `--resume`, d'autres arguments, ou une simulation |
| `u` | Détails de l'utilisation : fenêtres sur 5 heures, hebdomadaire et par modèle, avec leurs réinitialisations et l'origine des chiffres |
| `a` | Ajouter un compte (l'assistant ci-dessous) |
| `l` | Se reconnecter |
| `s` | Choisir les groupes du compte |
| `d` | En faire le compte par défaut |
| `f` | L'utiliser dans le dossier courant |
| `e` | Le renommer |
| `c` | Définir ou supprimer sa commande raccourci |
| `x` | Le supprimer (demande d'abord confirmation, sur Non par défaut) |

**Sièges Team et Enterprise.** Le choix automatique (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord, et [`auto_switch`](/features/next#auto-switch) lorsque vous l'activez) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Utiliser le choix automatique entre plusieurs sièges peut donc enfreindre les conditions d'Anthropic. C'est à vous de décider si vous l'utilisez dans ce cas : renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier. Les offres personnelles (Pro, Max) sont traitées dans la [note de la FAQ sur les conditions](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## L'onglet Health

L'onglet Health lance les vérifications de [`ccam doctor`](/features/doctor) à son ouverture et dessine la même vue : une première ligne avec un symbole par compte, puis un encadré par compte. Chaque ligne d'un encadré correspond à un domaine vérifié, `Folder`, `Login` et `Share links` en premier, marqué `✔` OK, `ℹ` information, `⚠` avertissement ou `✘` erreur, pour que vous voyiez d'un coup d'œil quel compte a un problème et où. Les constats qui n'appartiennent à aucun compte vont dans un encadré `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` passent d'un encadré à l'autre, et l'encadré sélectionné a une bordure de la couleur d'accentuation. `f` ne répare que l'encadré sélectionné, comme le ferait `ccam doctor --fix` : sur l'encadré d'un compte, il ne répare que ce compte, et sur l'encadré `General`, les constats qui n'appartiennent à aucun compte (chemins des plugins et raccourcis) ; tous les autres comptes sont seulement vérifiés à nouveau. Il demande d'abord (``Like `ccam doctor --fix`: missing links and group folders are made, and a real file where a link belongs is backed up first; where the group has its own copy, the account then sees the group's. Nothing is lost.``), et un encadré sans rien à réparer affiche simplement `nothing to repair in <name>`. `F` répare tout d'un coup, comme `ccam doctor --fix`, et `r` vérifie à nouveau.

## L'onglet Agents

L'onglet Agents montre ce qu'affiche [`ccam agents`](/features/agents) : les sessions en arrière-plan et les sessions ouvertes de tous les comptes, avec le compte auquel chacune appartient. Il interroge le Claude Code de chaque compte en même temps (`claude agents --json`, jusqu'à 15 secondes chacun) à la première ouverture de l'onglet, sur `r`, après un arrêt et après une modification de la configuration, et affiche `Asking every account's Claude Code for its sessions…` jusqu'à ce que tous aient répondu.

```console
  Accounts   Agents   Groups   Folders   Extensions   Health
 Background agents and open sessions of every account; enter attaches to one.

    Account   ID        Name   State          Folder      Started
 ❯  work      -         -      open (busy)    ~/code/web  now
    personal  49ec0028  probe  blocked, idle  ~/code/app  12m

 ✘ lab: could not list its sessions: exit status 1 'claude agents --json' is
   not available in this environment.

 enter to attach · x to stop · r to reload · ? for shortcuts
```

`State` est l'état d'une session en arrière-plan et ce que fait son processus, ou `open` pour une session interactive, comme dans `ccam agents`. Un compte qui n'a pas pu être interrogé reçoit une ligne `✘` sous la liste avec ce qu'a répondu Claude Code, comme `lab` ci-dessus, dont la vue des agents est désactivée. Les sessions en arrière-plan terminées sont laissées de côté (`ccam agents --all` les liste), et quand tous les comptes ont répondu sans aucune session, l'onglet affiche ``No sessions. Start one in the background with `ccam <account> --bg "…"`.`` Sur un terminal étroit, l'onglet retire `Started`, puis `Name`, puis `Folder`, et coupe un dossier trop long par la gauche.

`Enter` sur une session en arrière-plan ferme le tableau de bord et ouvre la session dans ce terminal sous son compte, comme le fait `ccam attach <id>` ; le tableau de bord ne revient pas (dans Claude Code, `←` mène à sa vue des agents et `Ctrl+Z` au shell). Une session ouverte n'a pas d'identifiant, donc `Enter` et `x` indiquent seulement qu'il faut une session en arrière-plan. `x` demande d'abord, en commençant sur `No, keep it` : ``It runs as personal. Like `ccam personal stop 49ec0028`, Claude Code stops it and keeps its conversation, which `ccam attach 49ec0028` opens again.`` Sur oui, ccam exécute `claude stop <id>` sous ce compte, affiche ``49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again`` ou la raison de l'échec, puis relit la liste.

## Changer les groupes

Dans l'onglet Groups, `Enter` ouvre l'éditeur : les éléments à gauche, les comptes à droite, `Tab` pour passer de l'un à l'autre, `Space` pour cocher ou décocher. Un élément que l'un des comptes reçoit déjà d'un autre groupe est grisé avec la raison, tout comme un compte qui recevrait un élément deux fois : un élément ne peut provenir que d'un seul groupe. `Enter` affiche le plan :

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Un fichier qui diffère de la copie du groupe fait d'abord l'objet d'une question (garder celui du groupe, utiliser le mien, ou fusionner). La modification s'exécute ensuite en arrière-plan avec le même code que [`ccam group`](/features/share) : la configuration est enregistrée d'abord, les données de chaque compte sont déplacées dans l'ordre, et les originaux vont dans `.ccam-backup/` ; un compte qui quitte le partage MCP garde les serveurs du groupe dans son `.claude.json`. Une modification qui déplacerait les données d'un compte en cours d'exécution est refusée, avec le nom du compte, y compris quand le compte a démarré pendant que le plan était ouvert, et un plan rendu obsolète par le changement de groupes d'un autre ccam est refusé avec `the groups changed while this plan waited; nothing was moved, make the change again`. Pendant qu'un autre ccam modifie les groupes, une nouvelle modification, un renommage ou une suppression affiche `another ccam is changing the groups right now` ; après une modification arrêtée à mi-chemin, ils attendent que `F` dans l'onglet Health (ou `ccam doctor --fix`) l'ait terminée.

## Ajouter un compte

`a` ouvre l'assistant qu'utilise aussi `ccam add` dans un terminal :

1. Le nom du compte.
2. Sa façon de se connecter : un abonnement Claude maintenant (dans le navigateur), une [clé d'API](/features/api), ou un abonnement Claude plus tard.
3. Pour une clé d'API : où ccam lit la clé à chaque lancement (`env:NAME` ou `command:COMMAND`), le point de terminaison, la façon dont la clé est envoyée et le modèle.
4. Les groupes à rejoindre, avec `default` coché ; tant qu'il n'existe aucun groupe, ce que partage le groupe `default`.
5. Une commande raccourci, si vous en voulez une.
6. Un récapitulatif à confirmer.

`Esc` revient d'une étape en arrière et, à la première étape, ferme l'assistant.

## Supprimer un compte

`x` demande confirmation avant que quoi que ce soit ne se passe, sur Non par défaut, et indique ce que fait la suppression : `It is logged out first (its login is revoked), then its folder goes to ccam's trash; the groups' shared files stay.` Pour un compte à clé d'API, il n'y a pas de connexion à révoquer, et quand un autre compte détient la même connexion, il n'est pas déconnecté. Un compte importé garde son dossier : `It was imported, so its folder /home/me/.claude stays where it is: it leaves its groups first and keeps a copy of what they share, then ccam forgets it.`

## Couleurs

Les couleurs suivent les rôles de l'interface terminal de Claude Code, pour que ccam s'y fonde naturellement ; les noms de groupe reçoivent une couleur dérivée de leur nom. ccam demande au terminal sa couleur de fond et choisit un thème clair ou sombre. Pour en imposer un, définissez-le en haut de `config.toml`, au-dessus du premier `[[group]]`, `[[profile]]` ou `[[folder]]` :

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Les thèmes `-daltonized` sont prévus pour le daltonisme rouge-vert, et les thèmes `-ansi` utilisent les 16 couleurs propres au terminal. `NO_COLOR` désactive les couleurs.

## Sécurité et cas particuliers

- Le tableau de bord appelle les mêmes fonctions que les commandes : toutes les règles s'appliquent donc, à savoir des plans avant tout déplacement de données, des refus pour les comptes en cours d'exécution, des sauvegardes, et des jetons de connexion qui ne quittent jamais un compte.
- Le tableau de bord n'a pas de simulation : `ccam --dry-run` sans nom affiche comment le compte de ce dossier démarrerait, et `o` sur un compte propose une simulation de ce seul lancement.
- Le terminal est restauré avant le démarrage de Claude Code ou l'ouverture de la connexion.
- Sur un écran trop petit pour une boîte de dialogue, les lignes de l'onglet disparaissent d'abord ; les lignes sont coupées avec `…` plutôt que renvoyées à la ligne. Sur un écran peu haut, les lignes vides et la description de l'onglet disparaissent ensuite, pour que la ligne des touches reste. En dessous de 40 colonnes ou de 8 lignes, le tableau de bord n'affiche que `Too small: ccam needs 40x8 (now 36x7)`.
- `Ctrl+C` pendant l'application d'une modification, d'une suppression, d'une mise à jour ou d'une réparation attend qu'elle se termine : la ligne d'état affiche `ccam quits when it is done.`
- Après que `U` a mis à jour ccam, la proposition ne revient pas quand le tableau de bord revient de l'ajout d'un compte ou d'une connexion.

## Voir aussi

- [Partager avec des groupes](/features/share)
- [Lancer Claude Code](/features/launch)
- [Gérer les comptes](/features/accounts)
- [Limites d'utilisation](/features/usage)
