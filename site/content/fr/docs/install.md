---
title: "Installer ccam"
description: "Installez ccam sous Linux ou macOS avec une seule commande, depuis un paquet .deb ou depuis les sources, et gardez-le à jour avec ccam update."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 4d7291e0f86e
---

Exécutez ceci dans un terminal sous Linux ou macOS :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Le script télécharge la version adaptée à votre système, vérifie sa somme de contrôle SHA-256 et place `ccam` dans `~/.local/bin`, où ccam crée aussi les raccourcis que vous choisirez plus tard, par exemple `ccwork`. Si ce dossier n'est pas dans votre `PATH`, le script affiche la ligne à ajouter à `~/.bashrc` ou `~/.zshrc`.

ccam lance Claude Code pour vous : installez donc d'abord [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) et vérifiez que `claude --version` fonctionne.

## Systèmes pris en charge

- Linux sur x86-64 et ARM64 : Ubuntu, Debian, Fedora, Arch, Alpine et d'autres distributions. Le binaire est statique et n'a aucune dépendance.
- macOS sur puce Apple et Intel.
- Windows n'est pas pris en charge directement. Dans WSL, ccam fonctionne comme sur n'importe quel système Linux.

## Autres méthodes d'installation

**Paquet Debian et Ubuntu.** Téléchargez `ccam_<version>_linux_amd64.deb` depuis la [page des versions](https://github.com/jhyunwoo/ccam/releases) et installez-le :

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Le paquet installe `/usr/bin/ccam`.

**Depuis les sources.** Avec Go 1.26 ou une version plus récente :

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Une version précise.** Définissez `CCAM_VERSION` :

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Complétion du shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

La complétion couvre toutes les commandes et options de `ccam`, y compris les noms de comptes.

## Vérifier l'installation

```sh
ccam version
ccam doctor
```

`ccam doctor` vérifie que Claude Code est trouvé et que chaque dossier de compte, chaque lien et chaque connexion sont en ordre. Poursuivez avec le [démarrage rapide](/docs/quick-start).

## Mettre à jour

```sh
ccam update --check
ccam update
```

`ccam update --check` indique seulement si une version plus récente existe. `ccam update` télécharge la version adaptée à votre système, vérifie sa somme de contrôle SHA-256 à l'aide du fichier `checksums.txt` de la version, puis remplace le fichier `ccam` en une seule étape. Si la somme de contrôle ne correspond pas, rien n'est modifié. Les raccourcis comme `ccwork` sont des liens vers ccam : ils lancent donc immédiatement la nouvelle version.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam ne recherche jamais de mises à jour de lui-même ; il ne cherche une nouvelle version que lorsque vous exécutez cette commande. Pour installer une version précise, y compris une version plus ancienne, passez `--version`. `ccam --dry-run update` affiche ce qui serait téléchargé et remplacé, sans rien modifier.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Relancer le script d'installation met aussi à jour la copie présente dans `~/.local/bin`. Si vous avez compilé ccam avec `go install`, exécutez de nouveau la même commande.

**Paquet Debian et Ubuntu.** `ccam update` ne remplace pas un `/usr/bin/ccam` installé par apt et affiche plutôt ces étapes. Téléchargez le nouveau `ccam_<version>_linux_amd64.deb` depuis la [page des versions](https://github.com/jhyunwoo/ccam/releases) et installez-le par-dessus l'ancien :

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Désinstaller

Supprimez d'abord les raccourcis des comptes, tant que ccam est encore installé. `ccam list` les affiche dans la colonne Shortcut ; supprimez chacun avec `ccam shortcut <account> --remove`, ou supprimez vous-même ces liens de `~/.local/bin`. Supprimez ensuite la commande, puis les dossiers de données si vous n'avez plus besoin des comptes :

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Supprimer `~/.local/share/ccam` supprime aussi les connexions et l'historique des conversations qui y sont stockés. Un dossier `~/.claude` que vous avez [importé](/docs/commands) reste là où il se trouve.
