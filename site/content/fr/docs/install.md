---
title: "Installer ccam"
metaTitle: "Installer ccam sous Linux, macOS et Windows"
description: "Installez ccam sous Linux, macOS ou Windows en une commande, via .deb ou go install, et mettez-le à jour avec ccam update. Installez d'abord Claude Code."
updated: 2026-10-28
order: 1
translation: machine
sourceHash: 193df30c9d9a
---

ccam lance Claude Code pour vous : installez donc d'abord [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) et vérifiez que `claude --version` fonctionne.

Exécutez ceci dans un terminal sous Linux ou macOS :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Le script télécharge la version adaptée à votre système, vérifie sa somme de contrôle SHA-256 et place `ccam` dans `~/.local/bin`, où ccam crée aussi les raccourcis que vous choisirez plus tard, par exemple `ccwork`. Si ce dossier n'est pas dans votre `PATH`, le script affiche la ligne à ajouter à `~/.bashrc` ou `~/.zshrc`.

Sous Windows, exécutez ceci dans PowerShell :

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Il télécharge `ccam_windows_<arch>.zip`, vérifie sa somme de contrôle SHA-256 et place `ccam.exe` dans `%USERPROFILE%\.local\bin`, le dossier où le programme d'installation de Claude Code place `claude.exe`. Si ce dossier n'est pas encore dans votre `PATH` utilisateur, le script l'y ajoute en laissant telles quelles les entrées comme `%USERPROFILE%\bin`. Ouvrez ensuite un nouveau terminal.

## Systèmes pris en charge

- Linux sur x86-64 et ARM64 : Ubuntu, Debian, Fedora, Arch, Alpine et d'autres distributions. Le binaire est statique et n'a aucune dépendance.
- macOS sur puce Apple et Intel.
- Windows 10 et 11 sur x86-64 et ARM64, dans PowerShell, Windows Terminal ou Git Bash :
  - Les dossiers partagés fonctionnent sans le Mode développeur.
  - Les fichiers partagés (`settings.json`, `CLAUDE.md`, `history.jsonl`) nécessitent le Mode développeur (Paramètres > Système > Pour les développeurs), car Windows ne permet qu'à ce mode, ou à un administrateur, de créer des liens vers des fichiers ; voir [partage sous Windows](/features/share#on-windows).
  - Sans le Mode développeur, le premier `ccam add` crée le groupe `default` avec seulement les serveurs MCP et les extensions, et indique comment ajouter le reste une fois le Mode développeur activé.
- Dans WSL, ccam fonctionne comme sur n'importe quel autre système Linux.

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

Si vous installez une étiquette de version (`@latest` ou `@vX.Y.Z`), ccam indique cette version, et `ccam update` ainsi que la proposition de mise à jour du tableau de bord fonctionnent. Une compilation de `@main`, ou d'une copie locale qui contient des commits postérieurs à une étiquette, indique `dev`, et `ccam update` la refuse sauf si vous passez `--version`.

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

Il ajoute :

- un skill qui propose les commandes qui connectent un compte ou déplacent des fichiers, mais ne les exécute jamais ;
- les commandes `/ccam:usage`, `/ccam:ext` et `/ccam:switch` ;
- deux hooks ([détails](/features/next#inside-claude-code)). Le premier affiche une notification dans le terminal quand une session s'arrête à une limite d'utilisation ; sous Linux et macOS, avec [`auto_switch`](/features/next#auto-switch) activé, ccam déplace alors la session vers un autre compte qui a de la marge. Le second permet à `/ccam:switch [account]` de passer la conversation à un autre compte.

Les hooks exécutent `ccam` : il doit donc se trouver dans le `PATH` que voit Claude Code.

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

Le tableau de bord (`ccam` sans argument, dans un terminal) vérifie aussi s'il existe une nouvelle version à chaque ouverture et, s'il y en a une, la propose : appuyez sur `U` et confirmez. Aucune autre commande ne vérifie les mises à jour, et ccam ne se met jamais à jour sans demander. Pour installer une version précise, y compris une version plus ancienne, passez `--version`. `ccam --dry-run update` affiche ce qui serait téléchargé et remplacé, sans rien modifier.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Relancer le script d'installation met aussi à jour la copie présente dans `~/.local/bin`. Si vous avez compilé ccam avec `go install` depuis une étiquette de version, vous pouvez exécuter `ccam update` ou relancer la même commande `go install`.

**Windows.** `ccam update` télécharge le `.zip` pour Windows. Un programme en cours d'exécution ne peut pas être remplacé sous Windows : l'ancien binaire est donc renommé en `ccam.exe.old`, que la mise à jour suivante supprime ; tant qu'une session tourne encore depuis un `ccam.exe.old` antérieur, le nom est `ccam.exe.old-<time>`. Les raccourcis qui sont des liens physiques sont renouvelés juste après. Relancer `install.ps1` fonctionne de la même façon, même pendant qu'une session Claude Code tourne encore depuis l'ancien binaire, et affiche `OK 2 shortcuts now run the new ccam` quand il renouvelle les raccourcis qui sont des liens physiques.

**Paquet Debian et Ubuntu.** `ccam update` ne remplace pas un `/usr/bin/ccam` installé par apt et affiche plutôt ces étapes. Téléchargez le nouveau `ccam_<version>_linux_amd64.deb` depuis la [page des versions](https://github.com/jhyunwoo/ccam/releases) et installez-le par-dessus l'ancien :

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Désinstaller

`ccam uninstall` supprime ccam et tout ce qu'il a créé. Il liste d'abord ce qui sera supprimé, puis demande confirmation :

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- Chaque compte est supprimé comme le fait `ccam remove`. Un compte situé dans les dossiers de ccam est d'abord déconnecté, puis supprimé avec son dossier, y compris son historique des conversations et ses réglages.
- Un dossier `~/.claude` que vous avez [importé](/docs/commands) reste là où il se trouve, pour `claude` tout court. S'il fait partie d'un groupe, il récupère une copie de tout ce qu'il partageait. La ligne d'état de ccam dans son `settings.json` cède la place à la ligne d'état que ccam enveloppait, s'il y en avait une. S'il a le plugin ccam, ccam affiche la commande `/plugin uninstall` à exécuter dans Claude Code.
- Les raccourcis, `~/.config/ccam`, `~/.local/share/ccam` (avec la corbeille de ccam) et `~/.local/state/ccam` sont supprimés, puis la commande `ccam` elle-même.

ccam refuse tant qu'un compte est en cours d'exécution ; quittez d'abord ses sessions Claude Code. Un `/usr/bin/ccam` installé par apt reste en place ; supprimez-le ensuite avec `sudo apt remove ccam`. Sous Windows, un programme en cours d'exécution ne peut pas se supprimer lui-même : `ccam uninstall` se termine donc en affichant la commande `Remove-Item` pour `ccam.exe`.
