---
title: "Le tableau de bord"
description: "ccam dans un terminal ouvre un tableau de bord plein écran : comptes, groupes, dossiers, extensions et vérifications ; chaque modification montre son plan."
updated: 2026-10-21
order: 15
translation: machine
sourceHash: "27bd153b48b2"
---

`ccam` sans argument ouvre un tableau de bord plein écran quand l'entrée et la sortie standard sont un terminal ; sinon, il affiche `ccam list`. Le tableau de bord montre chaque compte avec son utilisation et ses groupes, et fait ce que font les commandes : lancer, ajouter, connecter, changer de groupes, définir des règles de dossier, vérifier et réparer. Un lancement le ferme ; ajouter un compte ou se connecter le quitte le temps de passer au terminal, puis y revient.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

La ligne sous les noms des onglets indique en une phrase à quoi sert l'onglet courant. Dans un terminal plus large, l'onglet Accounts affiche aussi l'adresse e-mail de chaque compte.

La première image n'est dessinée qu'à partir de `config.toml`, en bien moins de 100 ms ; l'utilisation s'affiche au fur et à mesure que chaque compte répond. Chaque limite a sa propre barre : la fenêtre de 5 heures, la fenêtre hebdomadaire de tous les modèles et la fenêtre hebdomadaire propre à Fable ; `-` signifie que le compte n'a pas cette fenêtre. Un pourcentage passe au jaune à partir de 80 % et au rouge à 100 %, et les valeurs tirées d'un cache après un échec de consultation sont affichées en atténué.

## Onglets

| Onglet | Ce qu'il montre | Touches |
| --- | --- | --- |
| Accounts | numéro, nom, adresse e-mail, barres sur 5 heures, hebdomadaire et Fable, dernière utilisation, groupes ; `✔` le compte par défaut, `▸` le compte de ce dossier, `•` en cours d'exécution | ci-dessous |
| Groups | chaque groupe, ce qu'il partage et ses comptes, ainsi que les comptes sans groupe | `enter` modifier, `n` nouveau, `e` renommer, `x` supprimer |
| Folders | les règles [`ccam folder`](/features/folder) et celle qui s'applique ici | `n` utiliser un compte ici, `enter` changer, `x` supprimer |
| Extensions | ce que montre [`ccam ext`](/features/ext), chaque extension une seule fois avec les groupes dont elle provient | `r` recharger |
| Health | ce que signale [`ccam doctor`](/features/doctor), sous forme d'un encadré par compte (ci-dessous) | `f` réparer l'encadré sélectionné, `F` tout réparer (`doctor --fix`), `↑`/`↓` passer d'un encadré à l'autre |

`←`/`→` ou `Tab` changent d'onglet, `?` liste les touches de l'onglet courant, `r` consulte à nouveau l'utilisation sans le cache court, et `q` ou `Esc` quitte.

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

## L'onglet Health

L'onglet Health lance les vérifications de [`ccam doctor`](/features/doctor) à son ouverture et dessine la même vue : une première ligne avec un symbole par compte, puis un encadré par compte. Chaque ligne d'un encadré correspond à un domaine vérifié, `Folder`, `Login` et `Share links` en premier, marqué `✔` OK, `ℹ` information, `⚠` avertissement ou `✘` erreur, pour que vous voyiez d'un coup d'œil quel compte a un problème et où. Les constats qui n'appartiennent à aucun compte vont dans un encadré `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` passent d'un encadré à l'autre, et l'encadré sélectionné a une bordure de la couleur d'accentuation. `f` ne répare que l'encadré sélectionné, comme le ferait `ccam doctor --fix` : sur l'encadré d'un compte, il ne répare que ce compte, et sur l'encadré `General`, les constats qui n'appartiennent à aucun compte (chemins des plugins et raccourcis) ; tous les autres comptes sont seulement vérifiés à nouveau. Il demande d'abord, et un encadré sans rien à réparer affiche simplement `nothing to repair in <name>`. `F` répare tout d'un coup, comme `ccam doctor --fix`, et `r` vérifie à nouveau.

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

Un fichier qui diffère de la copie du groupe fait d'abord l'objet d'une question (garder celui du groupe, utiliser le mien, ou fusionner). La modification s'exécute ensuite en arrière-plan avec le même code que [`ccam group`](/features/share) : la configuration est enregistrée d'abord, les données de chaque compte sont déplacées dans l'ordre, et les originaux vont dans `.ccam-backup/`. Une modification qui déplacerait les données d'un compte en cours d'exécution est refusée, avec le nom du compte.

## Ajouter un compte

`a` ouvre l'assistant qu'utilise aussi `ccam add` dans un terminal :

1. Le nom du compte.
2. Sa façon de se connecter : un abonnement Claude maintenant (dans le navigateur), une [clé d'API](/features/api), ou un abonnement Claude plus tard.
3. Pour une clé d'API : où ccam lit la clé à chaque lancement (`env:NAME` ou `command:COMMAND`), le point de terminaison, la façon dont la clé est envoyée et le modèle.
4. Les groupes à rejoindre, avec `default` coché ; tant qu'il n'existe aucun groupe, ce que partage le groupe `default`.
5. Une commande raccourci, si vous en voulez une.
6. Un récapitulatif à confirmer.

`Esc` revient d'une étape en arrière et, à la première étape, ferme l'assistant.

## Couleurs

Les couleurs suivent les rôles de l'interface terminal de Claude Code, pour que ccam s'y fonde naturellement ; les noms de groupe reçoivent une couleur dérivée de leur nom. ccam demande au terminal sa couleur de fond et choisit un thème clair ou sombre. Pour en imposer un, définissez-le dans `config.toml` :

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Les thèmes `-daltonized` sont prévus pour le daltonisme rouge-vert, et les thèmes `-ansi` utilisent les 16 couleurs propres au terminal. `NO_COLOR` désactive les couleurs.

## Sécurité et cas particuliers

- Le tableau de bord appelle les mêmes fonctions que les commandes : toutes les règles s'appliquent donc, à savoir des plans avant tout déplacement de données, des refus pour les comptes en cours d'exécution, des sauvegardes, et des jetons de connexion qui ne quittent jamais un compte.
- Le tableau de bord n'a pas de simulation : `ccam --dry-run` sans nom affiche comment le compte de ce dossier démarrerait, et `o` sur un compte propose une simulation de ce seul lancement.
- Le terminal est restauré avant le démarrage de Claude Code ou l'ouverture de la connexion.
- Sur un écran trop petit pour une boîte de dialogue, les lignes de l'onglet disparaissent d'abord ; les lignes sont coupées avec `…` plutôt que renvoyées à la ligne.

## Voir aussi

- [Partager avec des groupes](/features/share)
- [Lancer Claude Code](/features/launch)
- [Gérer les comptes](/features/accounts)
- [Limites d'utilisation](/features/usage)
