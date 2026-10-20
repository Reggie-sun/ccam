---
title: "Vérifier et réparer"
description: "Vérifiez les liens de partage, les permissions, les connexions et la version de Claude Code avec ccam doctor, et réparez ce qui est sûr avec ccam doctor --fix."
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

`ccam doctor` vérifie votre installation de ccam : le binaire Claude Code, le dossier de chaque compte, sa connexion et ses liens de partage, les serveurs MCP partagés, le point de terminaison d'utilisation et vos commandes raccourcis. `ccam doctor --fix` répare ce qu'il peut sans perte de données et déplace d'abord dans un dossier de sauvegarde tout fichier qu'il remplace. Lancez-le quand `ccam <name>` s'arrête sur un problème de lien de partage, après avoir restauré des fichiers à la main ou après une mise à jour de Claude Code.

## Utilisation

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Chaque constat tient sur une ligne avec un symbole : `✓` OK, `·` information, `!` avertissement, `✗` erreur. Les constats propres à un compte commencent par le nom du compte.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Avec `--fix`, chaque constat réparé se termine par `→ fixed` :

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--fix` | Réparer ce qui peut l'être sans risque (les fichiers modifiés sont sauvegardés) | désactivé |
| `--dry-run` (globale) | Exécuter toutes les vérifications sans rien réparer, même avec `--fix` | désactivé |

`ccam doctor` n'accepte aucun argument. Il respecte `--dry-run` : `ccam --dry-run doctor --fix` affiche les mêmes constats que `ccam doctor` et ne modifie rien.

Le code de sortie vaut `1` quand au moins un constat `✗` reste non réparé, et `0` sinon. Les avertissements et les informations ne font jamais échouer l'exécution. Quand l'exécution échoue sans `--fix`, ccam termine par une ligne vide et `` `ccam doctor --fix` repairs what it can. ``. Une configuration impossible à charger se termine aussi par `1` et une erreur `ccam:`.

## Fonctionnement

`internal/cli/doctor.go` charge `~/.config/ccam/config.toml` et appelle `doctor.Run` dans `internal/doctor/doctor.go`, qui exécute les vérifications dans un ordre fixe et renvoie les constats. `--fix` est transmis sous la forme `Options.Fix`, mis à false lors d'une simulation.

<!-- motion:doctor-checks -->

Le bloc central s'exécute une fois par compte. Les vérifications, et ce que fait `--fix` pour chacune :

1. **Version de Claude Code** (`checkClaude`). Résout le binaire `claude` comme le fait le lanceur et exécute `claude --version`. Si la version ne commence pas par `2.1.302`, la version sur laquelle la structure de dossiers de ccam a été vérifiée (`doctor.LayoutVersion`), il ajoute `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Un binaire absent ou un `--version` qui échoue est une erreur. Pas de réparation.
2. **Identifiants dans l'environnement** (`checkAuthEnv`). Avertit quand `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN` est défini : chaque compte lancé depuis ce shell en hérite, et Claude Code le préfère à la connexion propre du compte. Pas de réparation ; supprimez la variable vous-même.
3. **Dossier** (`checkProfile`). Un dossier de compte absent est une erreur et met fin aux vérifications de ce compte. Des permissions avec des bits de groupe ou d'autres utilisateurs donnent un avertissement ; `--fix` exécute `chmod 0700`. Windows n'a pas de bits de permission : cette vérification et celle de `0600` ci-dessous y sont donc ignorées ; le dossier situé dans votre profil utilisateur est protégé par sa liste de contrôle d'accès.
4. **Fichier de connexion** (`checkCredentials`). `.credentials.json` est vérifié avec `lstat`. Un lien symbolique est une erreur, car chaque compte doit avoir son propre jeton, et n'est pas réparé. Un mode autre que `0600` donne un avertissement ; `--fix` exécute `chmod 0600`. Ensuite la connexion est lue (sous macOS, d'abord depuis le trousseau). Une absence de connexion donne un avertissement qui renvoie vers `ccam login <name>`. Un jeton de renouvellement expiré donne l'erreur `login expired` ; un jeton d'accès expiré n'est qu'une information, puisque Claude Code le renouvelle au lancement suivant. Aucun de ces cas n'est réparé.
5. **Liens de partage** (`share.Check`). Compare chaque chemin de `share.Layout` avec les paramètres de partage du compte, en utilisant uniquement `lstat` et `readlink`. Il signale quatre cas : un lien absent, un lien qui pointe ailleurs, un vrai fichier là où devrait se trouver un lien, et un lien restant pour un élément désormais local. Un lien absent est un avertissement, les autres cas sont des erreurs.
6. **Entrées inconnues** (`checkUnknownEntries`). Liste les noms de premier niveau du dossier du compte que ccam ne connaît pas (`share.Known`), comme indice d'un changement de structure de Claude Code. Information uniquement.
7. **Chemins des plugins** (`checkPlugins`). Les registres de plugins (`installed_plugins.json`, `known_marketplaces.json`) enregistrent des chemins absolus. Dans le `plugins/` de chaque groupe, les chemins qui passent par le dossier d'un compte cessent de fonctionner quand ce compte est supprimé ou renommé. `--fix` les réécrit vers le dossier du groupe avec `profile.RewritePluginPaths`, une écriture atomique du fichier de registre.
8. **Serveurs MCP** (`checkMCP`). Charge une fois le `mcp.json` de chaque groupe. Un fichier illisible est une erreur ; un serveur qui ne passe pas `mcp.Validate` donne un avertissement. Pas de réparation.
9. **Format de la réponse d'utilisation** (`checkUsageFormat`). Interroge le point de terminaison d'utilisation, sans le cache et avec la limite habituelle de 3 secondes, pour un compte après l'autre jusqu'à ce que l'un réponde. Si la réponse ne peut plus être analysée, il avertit que `ccam usage` se rabat sur le cache de la ligne d'état. C'est la seule vérification qui utilise le réseau.
10. **Raccourcis** (`checkShims`). Uniquement quand un compte a un raccourci. Des liens absents dans `~/.local/bin` donnent un avertissement ; `--fix` appelle `shim.Install`, qui ne crée que les liens manquants et ne touche jamais aux autres fichiers. Un dossier bin absent du `PATH` donne un avertissement et n'est pas réparé.
11. **Dossiers** (`checkFolders`). Uniquement quand des entrées [`ccam folder`](/features/folder) existent. Une entrée dont le compte n'existe pas est une erreur, et une entrée dont le dossier a disparu donne un avertissement. Aucune n'est réparée, car supprimer une entrée changerait le compte qui démarre à cet endroit.

### Réparer les liens de partage

Le chemin de lancement crée sur-le-champ un lien manquant, mais s'arrête pour tout ce qui pourrait faire perdre des données et renvoie ici :

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Pour tout problème de lien de partage, `--fix` appelle `share.Build` avec les paramètres de partage du compte comme état actuel et comme état cible, puis `Plan.Apply`. `Build` redéduit chaque étape de ce qui se trouve sur le disque, le même plan que celui de `ccam share`, si bien que la réparation suit les mêmes règles que l'activation du partage.

<!-- motion:doctor-repair -->

Ici, `settings.json` est un vrai fichier. Comme le groupe a déjà un `settings.json` différent, le conflit prend son choix par défaut, `KeepGroup`, et la copie du groupe est conservée. Le fichier du compte est déplacé vers `.ccam-backup/<time>/settings.json` dans le dossier du compte, et `os.Symlink` relie le chemin au groupe. Autres cas :

- Si le groupe n'a pas encore de copie, le fichier du compte est déplacé dans le groupe et devient son premier contenu.
- Pour un dossier, les fichiers absents du groupe sont déplacés dans le groupe ; le reste va dans la sauvegarde.
- Pour `history.jsonl`, le choix par défaut est `merge by line` : les lignes du groupe, puis vos nouvelles lignes. La copie précédente du groupe est enregistrée dans `.ccam-backup/<time>/group/`.
- Un lien vers un autre groupe est remplacé. Un lien vers un dossier extérieur à ccam va dans la sauvegarde, et son contenu alimente le groupe si celui-ci est vide.
- Un lien restant pour un élément désormais local est supprimé, et le compte reçoit une copie du contenu du groupe.

À part la suppression de liens, chaque étape déplace des données vers le groupe ou le dossier de sauvegarde. Rien n'est supprimé.

## Sécurité et cas particuliers

- Un simple `ccam doctor` ne modifie jamais rien. Dans les scripts, utilisez `--dry-run` avec `--fix` pour en être sûr.
- `doctor --fix` ne vérifie pas si un compte est en cours d'exécution. Quittez Claude Code pour un compte avant de réparer ses liens de partage, car la session continue d'écrire dans ces fichiers.
- Si une réparation de partage échoue, ccam signale `could not repair share links: <error>` pour ce compte et saute sa vérification des entrées inconnues. Les autres comptes sont tout de même vérifiés.
- Si `ccam share` échoue à mi-parcours, il enregistre d'abord les paramètres cibles et vous indique de terminer avec `ccam doctor --fix`.
- Un `.credentials.json` en lien symbolique, une connexion expirée et des identifiants dans l'environnement sont signalés mais jamais modifiés.
- Les réécritures des registres de plugins sont écrites de façon atomique mais ne sont pas copiées dans `.ccam-backup/`.

## Voir aussi

- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Copier une configuration existante](/features/migrate)
- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Commandes](/docs/commands)
