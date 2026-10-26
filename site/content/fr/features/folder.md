---
title: "Comptes par dossier"
description: "ccam folder fait utiliser un compte à un dossier et à son contenu quand vous lancez Claude Code sans nom de compte, comme work pour le code de l'entreprise."
updated: 2026-10-26
order: 11
translation: machine
sourceHash: 25f82da64580
---

Après `ccam folder work ~/code/company`, tout lancement sans nom de compte dans `~/code/company` ou ses sous-dossiers démarre Claude Code avec `work`. En dehors, le compte par défaut démarre comme avant. Nommer un compte, avec `ccam personal` ou un raccourci comme `ccpersonal`, lance toujours ce compte.

## Utilisation

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Ces lancements utilisent le compte du dossier :

- `ccam -p "…"` et tout autre lancement dont le premier argument est une option destinée à Claude Code
- `ccam run` sans nom, et `ccam --dry-run`
- `ccam --fallback -p "…"` pour sa première tentative
- le tableau de bord, `ccam` seul, qui place d'abord son curseur sur le compte du dossier, y compris quand vous êtes arrivé dans le dossier par un lien symbolique

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `<account>` | Nom ou numéro du compte à utiliser | |
| `[folder]` | Le dossier, relatif ou absolu | le dossier courant |
| `--remove` | Supprimer l'entrée de ce dossier précis | désactivé |
| `--dry-run` | Pris en charge : affiche la modification sans l'enregistrer | désactivé |

## Fonctionnement

Chaque dossier est une entrée `[[folder]]` dans `~/.config/ccam/config.toml`, le fichier que chaque lancement lit de toute façon :

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` convertit le dossier en chemin absolu, liens symboliques résolus, et l'enregistre. Définir de nouveau un dossier remplace son compte.
2. Un lancement sans nom appelle `launcher.HereProfile`. Sans aucune entrée, il prend immédiatement le compte par défaut.
3. Sinon, il compare le répertoire de travail à chaque `path` sous forme de chaînes : le dossier lui-même ou un dossier en dessous, jamais un dossier qui commence seulement par les mêmes lettres (`/code/company` ne couvre pas `/code/companyx`). L'entrée la plus profonde l'emporte : `~/code` peut ainsi être `personal` pendant que `~/code/company` est `work`.
4. Ce n'est que si le répertoire de travail, tel quel, ne correspond à rien que ccam résout ses liens symboliques et compare de nouveau.
5. Sous Windows et macOS, la comparaison ignore les majuscules et les minuscules, comme le font généralement leurs systèmes de fichiers.

Aucun fichier autre que `config.toml` n'est lu : un lancement dans un dossier n'ajoute donc aucun temps mesurable (4,0 ms avec ou sans entrées, `hyperfine`, 300 exécutions).

## Sécurité et cas particuliers

- `ccam rename` reporte les entrées de dossier sur le nouveau nom, et `ccam remove` supprime les entrées du compte supprimé.
- Une entrée dont le compte n'existe plus arrête un lancement sans nom dans ce dossier : `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` signale ces entrées, ainsi que celles dont le dossier a disparu.
- Si vous exécutez `ccam folder --remove` dans un sous-dossier qui hérite son compte d'un dossier parent, ccam nomme le dossier parent et ne supprime rien.
- `ccam next` choisit par rotation ou selon l'utilisation et ne tient pas compte des dossiers.

## Voir aussi

- [Lancer Claude Code](/features/launch)
- [Commandes raccourcis](/features/shortcuts)
- [Gérer les comptes](/features/accounts)
- [Référence des commandes](/docs/commands)
