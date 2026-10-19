---
title: "Commandes raccourcis"
description: "Donnez à un compte sa propre commande, comme ccwork, avec ccam shortcut, et recréez les liens dans ~/.local/bin avec ccam shims install."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

Un raccourci est une commande que vous nommez vous-même et qui lance Claude Code avec un compte donné : `ccwork --resume` fait la même chose que `ccam work --resume`. C'est un lien vers ccam dans `~/.local/bin` ; il fonctionne donc dans bash, zsh, les scripts et tout ce qui lance des programmes depuis le `PATH`. Les raccourcis sont facultatifs ; créez-en un pour les comptes que vous lancez souvent.

## Utilisation

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam confirme chaque modification :

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Un compte a au plus un raccourci, enregistré sous `shortcut` dans `~/.config/ccam/config.toml`. L'assistant de `ccam add` demande `Create a shortcut command for <name>? [y/N]` puis `Shortcut name [cc<name>]:`. En dehors de l'assistant, `ccam add` n'en crée un qu'avec `--shortcut`. `ccam list` affiche le raccourci de chaque compte, ou `-`.

Un raccourci transmet tous les arguments à Claude Code, de la même façon que `ccam <name>` après le nom. Cela inclut `--dry-run` et `--help` : pour prévisualiser un lancement, utilisez `ccam --dry-run work`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--remove` | Supprimer le raccourci du compte et son lien | désactivé |
| `-h`, `--help` | Aide de la commande | |
| `--dry-run` | `ccam shortcut` le respecte : il vérifie le nom et affiche `Would set the shortcut of work to ccwork.` ou `Would remove ccwork, the shortcut of work.` sans rien enregistrer ni relier | désactivé |

`ccam shims install` n'accepte aucun argument et refuse `--dry-run` avec `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Règles de nommage, vérifiées dans cet ordre :

- de 1 à 64 caractères parmi `a-z`, `0-9`, `.`, `-` et `_`, commençant par une lettre ou un chiffre ;
- ni `ccam`, ni `claude`, ni une commande ccam comme `run` ou `usage` ;
- pas le raccourci d'un autre compte ;
- pas un fichier de `~/.local/bin` qui ne serait pas déjà un lien vers ccam ;
- pas le nom d'un autre exécutable présent dans votre `PATH`.

Indiquer à la fois une commande et `--remove` s'arrête avec `give a command or --remove, not both`. Redonner à un compte le nom qu'il a déjà affiche `work: ccwork (unchanged)` et recrée son lien s'il est absent.

## Fonctionnement

Définir un nom exécute `checkShortcut`, qui appelle `shim.Check` avec le dossier bin et le chemin réel du binaire ccam :

<!-- motion:shortcuts-check -->

Le nom descend les vérifications dans l'ordre. À côté de chaque vérification figure un nom qu'elle refuse : `Work` ne respecte pas le motif, `run` est une commande ccam, `ccside` appartient à un autre compte, `ccfile` est un fichier existant dans `~/.local/bin`, et `git` masquerait `/usr/bin/git`. Seul un nom qui passe toutes les vérifications est enregistré et relié.

1. `shim.Check` compare le nom à `^[a-z0-9][a-z0-9._-]{0,63}$`, puis à `config.Reserved`, puis à `cfg.FindShortcut`. Il appelle `os.Lstat` sur `~/.local/bin/<name>` et, pour chaque autre dossier du `PATH`, il appelle `os.Stat` et refuse un fichier exécutable qui n'est pas un lien vers ccam.
2. ccam enregistre `config.toml` avec `fsutil.AtomicWrite` (fichier temporaire, `fsync`, renommage).
3. `shim.Remove` supprime le lien de l'ancien raccourci, mais seulement s'il s'agit d'un lien vers ccam.
4. `shim.Install` crée le lien de chaque raccourci configuré qui manque, avec `os.Symlink` vers le binaire ccam, résolu avec `os.Executable` et `filepath.EvalSymlinks`. Il ne supprime jamais rien.

Un lien est considéré comme appartenant à ccam quand `os.Readlink` renvoie le chemin de ccam, ou quand le lien et le chemin de ccam se résolvent vers le même fichier. Chaque suppression et chaque vérification avant écrasement passe par ce test, si bien que ccam ne touche pas aux fichiers qu'il n'a pas créés.

Quand vous lancez un raccourci, le nom que vous avez tapé est tout ce dont ccam a besoin :

<!-- motion:shortcuts-link -->

Le shell trouve `ccwork` dans le `PATH`, le noyau suit le lien et lance le binaire ccam, et `argv[0]` vaut toujours `ccwork`. Dans `app.Main`, tout nom autre que `ccam` charge une fois `config.toml` et appelle `shim.Lookup`, qui trouve le compte dont le `shortcut` est égal à `filepath.Base(argv[0])`. `launcher.LaunchWith` suit ensuite les mêmes étapes que [`ccam <name>`](/features/launch) avec la configuration déjà chargée, et transmet `argv[1:]` tel quel à `claude`.

Pourquoi un lien symbolique plutôt qu'un alias de shell ou un script enveloppe : un lien symbolique fonctionne dans tous les shells et dans les scripts, il n'ajoute pas de démarrage de shell supplémentaire (environ 1 à 2 ms), et le lancement ne lit toujours qu'un seul fichier. `scripts/perf-check.sh` mesure un raccourci avec `hyperfine` par rapport à la même limite de 10 ms (p95) que `ccam <name>`.

## Sécurité et cas particuliers

- **Un lien qu'aucun compte ne revendique.** Si un lien vers ccam dans `~/.local/bin` porte un nom qu'aucun compte n'utilise, par exemple après avoir modifié `shortcut` à la main dans `config.toml`, ccam ne se rabat pas sur le compte par défaut. Il s'arrête avec `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Redéfinissez le nom avec `ccam shortcut`, ou supprimez le lien.
- **Fichiers existants.** ccam ne remplace jamais un fichier qu'il n'a pas créé. `ccam shims install` s'arrête avec `<path> already exists and is not a ccam link, so it was left alone`.
- **Un binaire ccam déplacé.** Les anciens liens pointent alors vers un chemin qui n'existe plus, et ccam ne les considère pas comme les siens. Supprimez-les, puis lancez `ccam shims install`.
- **Renommer et supprimer.** `ccam rename` conserve le nom du raccourci, et le lien continue de fonctionner, car il pointe vers ccam et ccam recherche le nom dans `config.toml`. `ccam remove` supprime le lien du compte. `ccam shortcut <account> <new>` supprime l'ancien lien avant de créer le nouveau.
- **PATH.** Quand ccam crée un lien et que `~/.local/bin` n'est pas dans votre `PATH`, il affiche `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` suivi de la ligne `export PATH=...` à ajouter.
- **Vérifier plus tard.** `ccam doctor` signale `missing shortcuts in /home/me/.local/bin: ccwork` et avertit quand le dossier n'est pas dans le `PATH`. `ccam doctor --fix` crée les liens manquants.

## Voir aussi

- [Lancer Claude Code](/features/launch)
- [Gérer les comptes](/features/accounts)
- [Référence des commandes](/docs/commands#shortcuts)
- [Installer ccam](/docs/install)
