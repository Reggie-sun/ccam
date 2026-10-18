---
title: "Installer ccam"
description: "Installez ccam sous Linux ou macOS avec une seule commande, depuis un paquet .deb ou depuis les sources. Claude Code doit être installé avant."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

Exécutez ceci dans un terminal sous Linux ou macOS :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Le script télécharge la version adaptée à votre système, vérifie sa somme de contrôle SHA-256 et place deux commandes dans `~/.local/bin` : `ccam` et son alias court `ccs`. Si ce dossier n'est pas dans votre `PATH`, le script affiche la ligne à ajouter à `~/.bashrc` ou `~/.zshrc`.

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

Le paquet installe `/usr/bin/ccam` et le lien `/usr/bin/ccs`.

**Depuis les sources.** Avec Go 1.26 ou une version plus récente :

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
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

La complétion couvre à la fois `ccam` et `ccs`, y compris les noms de comptes.

## Vérifier l'installation

```sh
ccam version
ccam doctor
```

`ccam doctor` vérifie que Claude Code est trouvé et que chaque dossier de compte, chaque lien et chaque connexion sont en ordre. Poursuivez avec le [démarrage rapide](/docs/quick-start).

## Désinstaller

Supprimez les deux commandes et les raccourcis des comptes, puis les dossiers de données si vous n'avez plus besoin des comptes :

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Supprimer `~/.local/share/ccam` supprime aussi les connexions et l'historique des conversations qui y sont stockés. Un dossier `~/.claude` que vous avez [importé](/docs/commands) reste là où il se trouve.
