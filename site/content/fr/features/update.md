---
title: "Mettre à jour ccam"
description: "ccam update remplace ccam par la dernière version de GitHub après vérification SHA-256 ; --check informe seulement et --version installe une version précise."
updated: 2026-10-23
order: 10
translation: machine
sourceHash: 000d37de5f29
---

`ccam update` télécharge une version de ccam depuis GitHub, la vérifie à l'aide du `checksums.txt` de la version et remplace le binaire en cours d'exécution en un seul renommage. Utilisez-le si vous avez installé ccam avec le script d'installation ou depuis une archive. Le [tableau de bord](/features/dashboard) demande aussi la dernière version à GitHub à chaque ouverture et la propose avec `U`, qui exécute les mêmes étapes après confirmation. Aucune autre commande ne contacte GitHub, si bien qu'elles fonctionnent de la même façon hors ligne.

## Utilisation

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Quand il n'existe pas de version plus récente, ccam affiche `ccam 0.4.0 is up to date.` et ne modifie rien. Avec `--version` et la version que vous utilisez déjà, il affiche `ccam 0.4.0 is already installed.`

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--check` | Indiquer seulement si une version plus récente existe | désactivé |
| `--version <tag>` | Installer cette étiquette de version, comme `v0.3.0`, plus récente ou plus ancienne que l'actuelle | la dernière |
| `--dry-run` | Respecté : recherche la dernière étiquette (sauf si `--version` en nomme une), effectue les mêmes vérifications sur le binaire, et affiche les chemins de téléchargement, de somme de contrôle et de remplacement sans rien télécharger | désactivé |

`--check` et `--version` ne peuvent pas être utilisés ensemble. L'étiquette doit avoir la forme `vMAJOR.MINOR.PATCH`, avec un suffixe facultatif de préversion ou de build ; le `v` initial peut être omis. Tout le reste échoue, par exemple `--version takes a release tag such as v0.3.0, not "1.2"`.

La variable d'environnement `CCAM_RELEASES_URL` remplace `https://github.com/jhyunwoo/ccam/releases`, pour un miroir ou pour des tests. Elle doit être organisée de la même façon : `<url>/latest` redirige vers `<url>/tag/<tag>`, et les fichiers se trouvent à `<url>/download/<tag>/<asset>`.

## Fonctionnement

Le code se trouve dans `internal/update`, que seul `ccam update` importe. Le chemin de lancement ne l'appelle jamais.

1. `cli.updateCmd` trouve le binaire en cours d'exécution avec `os.Executable` et résout les liens symboliques avec `filepath.EvalSymlinks`. C'est ce chemin réel qui est remplacé.
2. Sans `--version`, `Updater.Latest` demande `<releases>/latest` sans suivre la redirection, et lit l'étiquette dans l'en-tête `Location` (`…/tag/v0.5.0`). Cela passe par la page des versions plutôt que par l'API GitHub, si bien que la limite de 60 requêtes par heure de l'API ne s'applique pas. La recherche dispose de 10 secondes.
3. `parseVersion` vérifie l'étiquette avec un motif strict. Une étiquette se retrouve dans des URL de téléchargement et dans une commande shell affichée, donc rien de ce qui contient `/`, `%`, des espaces ou des caractères spéciaux du shell ne passe. `compareVersions` la compare à la version en cours d'exécution ; une préversion est classée avant sa version finale.
4. `managed()` et `writable()` décident si ccam a le droit de remplacer le binaire (voir plus bas).
5. `Updater.Fetch` télécharge `ccam_<os>_<arch>.tar.gz` (sous Windows `ccam_windows_<arch>.zip`) et `checksums.txt` en même temps, avec 2 minutes pour les deux. Chaque téléchargement est plafonné à 256 Mo.
6. La somme SHA-256 de l'archive doit correspondre à sa ligne dans `checksums.txt` avant toute autre chose. Une ligne absente ou une somme de contrôle différente arrête la mise à jour avec `nothing was changed`.
7. `extract()` lit le tar et prend le premier fichier ordinaire nommé `ccam` (`extractZip` prend `ccam.exe` dans le zip). Les dossiers, les liens et les autres entrées sont ignorés, et le binaire est lui aussi plafonné à 256 Mo.
8. `Updater.Install` écrit le nouveau binaire avec `fsutil.AtomicWrite` : un fichier temporaire `.ccam.tmp*` à côté du binaire, avec les bits de permission de l'ancien fichier, écrit sur le disque puis renommé par-dessus l'ancien.

L'ensemble du processus, de la recherche de l'étiquette jusqu'au renommage :

<!-- motion:update-pipeline -->

L'archive et `checksums.txt` sont téléchargés côte à côte, et rien n'est écrit tant que la somme de contrôle ne correspond pas. Le renommage est atomique : un plantage ou un disque plein laisse soit l'ancien binaire, soit le nouveau, jamais un binaire à moitié écrit. Les commandes raccourcis comme `~/.local/bin/ccwork` sont des liens symboliques vers le binaire ccam ; elles exécutent donc la nouvelle version sans aucune modification.

Sous Windows, le fichier d'un programme en cours d'exécution ne peut être ni remplacé ni supprimé, seulement renommé. `Updater.Install` écrit le nouveau binaire dans `ccam.exe.new`, renomme le `ccam.exe` en cours d'exécution en `ccam.exe.old`, puis renomme le nouveau fichier en `ccam.exe`. Si la dernière étape échoue, l'ancien binaire reprend son nom. La mise à jour suivante supprime les fichiers `ccam.exe.old` qu'elle peut ; tant qu'une session tourne encore depuis un `ccam.exe.old` antérieur, le binaire en cours d'exécution devient plutôt `ccam.exe.old-<time>`. Les raccourcis qui sont des liens physiques plutôt que des liens symboliques ([Commandes raccourcis](/features/shortcuts)) sont renouvelés juste après, pour qu'ils exécutent eux aussi la nouvelle version ; cela couvre les liens physiques de tout ancien ccam que ccam a créés, pas seulement ceux du binaire qui vient d'être remplacé.

### Installations gérées par un paquet

Un ccam installé depuis le `.deb` appartient à apt, et le remplacer à l'insu d'apt fausserait la base de données des paquets. Avant de télécharger, ccam vérifie à qui appartient le binaire :

<!-- motion:update-managed -->

1. Un binaire dans `/usr/bin`, `/bin`, `/usr/sbin` ou `/sbin` est considéré comme géré par un paquet.
2. Ailleurs, si `dpkg` est installé et que `dpkg -S <path>` réussit, il est aussi considéré comme géré par un paquet.
3. Sinon, ccam crée puis supprime un fichier temporaire `.ccam-update-*` à côté du binaire pour s'assurer qu'il peut y écrire.

Pour un binaire géré par un paquet, ccam s'arrête et affiche les étapes apt :

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Sur les systèmes autres que Linux, seule la ligne du script d'installation est affichée : `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex` sous Windows. Quand le dossier n'est pas accessible en écriture, ccam propose de relancer le script d'installation, qui installe ccam dans `~/.local/bin`.

## Sécurité et cas particuliers

- Une version de développement indique `dev` comme version et n'a pas de version publiée à laquelle se comparer ; `ccam update` et `ccam update --check` la refusent donc. `ccam update --version vX.Y.Z` installe une version publiée par-dessus. Une compilation faite avec `go install github.com/jhyunwoo/ccam/cmd/ccam@vX.Y.Z` (ou `@latest`) indique cette version publiée, si bien que `ccam update` et la proposition du tableau de bord fonctionnent pour elle ; `@main` ou une compilation locale postérieure à une étiquette reste `dev`.
- Un téléchargement échoué, un fichier absent pour votre plateforme ou une somme de contrôle qui ne correspond pas laissent le binaire installé tel quel.
- Une version qui n'existe pas échoue avec `there is no release <tag>`. Quand GitHub refuse avec un HTTP 403 ou 429, ccam indique qu'il s'agit probablement d'une limite de débit et vous demande de réessayer plus tard.
- `--check` ne fait que lire. Il ne vérifie pas si le binaire est géré par un paquet ou accessible en écriture.
- `--version` avec une étiquette plus ancienne permet de revenir en arrière.

## Voir aussi

- [Installation](/docs/install)
- [Commandes](/docs/commands)
- [FAQ](/docs/faq)
