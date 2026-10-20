---
title: "Installer ccam"
metaTitle: "Installer ccam sous Linux, macOS et Windows"
description: "Installez ccam sous Linux, macOS ou Windows en une commande, via .deb ou les sources, et mettez-le à jour avec ccam update. Installez d'abord Claude Code."
updated: 2026-10-20
order: 1
translation: machine
sourceHash: aea4984c1576
---

Exécutez ceci dans un terminal sous Linux ou macOS :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Le script télécharge la version adaptée à votre système, vérifie sa somme de contrôle SHA-256 et place `ccam` dans `~/.local/bin`, où ccam crée aussi les raccourcis que vous choisirez plus tard, par exemple `ccwork`. Si ce dossier n'est pas dans votre `PATH`, le script affiche la ligne à ajouter à `~/.bashrc` ou `~/.zshrc`.

Sous Windows, exécutez ceci dans PowerShell :

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Il télécharge `ccam_windows_<arch>.zip`, vérifie sa somme de contrôle SHA-256, place `ccam.exe` dans `%USERPROFILE%\.local\bin`, le dossier où le programme d'installation de Claude Code place `claude.exe`, et ajoute ce dossier à votre `PATH` utilisateur s'il n'y figure pas encore. Ouvrez ensuite un nouveau terminal.

ccam lance Claude Code pour vous : installez donc d'abord [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) et vérifiez que `claude --version` fonctionne.

## Systèmes pris en charge

- Linux sur x86-64 et ARM64 : Ubuntu, Debian, Fedora, Arch, Alpine et d'autres distributions. Le binaire est statique et n'a aucune dépendance.
- macOS sur puce Apple et Intel.
- Windows 10 et 11 sur x86-64 et ARM64, dans PowerShell, Windows Terminal ou Git Bash. Les dossiers partagés fonctionnent pour tous les utilisateurs. Les fichiers partagés (`settings.json`, `CLAUDE.md`, `history.jsonl`) nécessitent le Mode développeur (Paramètres > Système > Pour les développeurs), car Windows ne permet qu'à ce mode, ou à un administrateur, de créer des liens vers des fichiers ; voir [partage sous Windows](/features/share#on-windows).
- Dans WSL, ccam fonctionne comme sur n'importe quel système Linux.

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

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Complétion du shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

La complétion couvre toutes les commandes et options de `ccam`, y compris les noms de comptes.

## Plugin Claude Code

Le plugin ccam apprend à Claude Code lui-même à utiliser ccam : vérifier quel compte a de la marge, lancer des workers `claude -p` en parallèle sur différents comptes avec `ccam pick`, et expliquer l'absence d'un skill ou d'un serveur MCP avec `ccam ext`. Installez-le dans Claude Code :

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Il ajoute un skill et deux commandes, `/ccam:usage` et `/ccam:ext`. Le skill n'exécute jamais de commandes qui connectent un compte ou déplacent des fichiers ; il les propose à la place.

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

**Windows.** `ccam update` télécharge le `.zip` pour Windows. Un programme en cours d'exécution ne peut pas être remplacé sous Windows : l'ancien binaire est donc renommé en `ccam.exe.old`, que la mise à jour suivante supprime.

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

Sous Windows, dans PowerShell :

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Supprimer `~/.local/share/ccam` supprime aussi les connexions et l'historique des conversations qui y sont stockés. Un dossier `~/.claude` que vous avez [importé](/docs/commands) reste là où il se trouve.
