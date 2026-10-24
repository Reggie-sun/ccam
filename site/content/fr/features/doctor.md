---
title: "Vérifier et réparer"
description: "Vérifiez les liens de partage, les permissions, les connexions et la version de Claude Code avec ccam doctor, et réparez ce qui est sûr avec ccam doctor --fix."
updated: 2026-10-23
order: 9
translation: machine
sourceHash: 11ad34630dcf
---

`ccam doctor` vérifie votre installation de ccam : le binaire Claude Code, le dossier de chaque compte, sa connexion et ses liens de partage, les serveurs MCP partagés, le point de terminaison d'utilisation et vos commandes raccourcis. `ccam doctor --fix` répare ce qu'il peut sans perte de données et déplace d'abord dans un dossier de sauvegarde tout fichier qu'il remplace. Lancez-le quand `ccam <name>` s'arrête sur un problème de lien de partage, après avoir restauré des fichiers à la main ou après une mise à jour de Claude Code.

## Utilisation

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

La première ligne marque tous les comptes d'un coup, pour que vous voyiez où sont les problèmes avant de lire quoi que ce soit d'autre. En dessous, chaque compte a son encadré, et les constats qui n'appartiennent à aucun compte (Claude Code, variables d'environnement, plugins, serveurs MCP, raccourcis, règles de dossier, groupes et réglages) vont dans un encadré `General`. Le titre de l'encadré montre le pire symbole du compte et son nombre de problèmes. Chaque ligne à l'intérieur correspond à un domaine vérifié, toujours `Folder`, `Login` et `Share links` dans cet ordre, puis les autres, avec un symbole : `✔` OK, `ℹ` information, `⚠` avertissement, `✘` erreur. Un dossier sans problème affiche `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
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

`ccam doctor --fix` repairs what it can.
```

Avec `--fix`, chaque constat réparé se termine par `→ fixed` et ne compte plus comme un problème :

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

Dans un terminal, les encadrés s'adaptent à sa largeur et les longs messages passent à la ligne, dans les couleurs du terminal (`NO_COLOR` les désactive). Un encadré n'est jamais plus étroit que 29 colonnes, si bien qu'aucun mot n'est coupé dans un terminal étroit. Envoyée vers un pipe ou un fichier, la sortie n'a pas de couleurs et aucun message n'est coupé ni renvoyé à la ligne, si bien que `ccam doctor | grep settings.json` trouve la ligne entière ; chaque encadré a alors la largeur de la plus longue ligne. L'onglet Health du tableau de bord affiche les mêmes encadrés.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--fix` | Réparer ce qui peut l'être sans risque (les fichiers modifiés sont sauvegardés) | désactivé |
| `--dry-run` (globale) | Exécuter toutes les vérifications sans rien réparer, même avec `--fix` | désactivé |

`ccam doctor` n'accepte aucun argument. Il respecte `--dry-run` : `ccam --dry-run doctor --fix` affiche les constats de `ccam doctor`, y ajoute ce que `--fix` mettrait de côté, et ne modifie rien :

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

Le code de sortie vaut `1` quand au moins un constat `✘` reste non réparé, et `0` sinon. Les avertissements et les informations ne font jamais échouer l'exécution. Quand l'exécution échoue sans `--fix`, ccam termine par une ligne vide et `` `ccam doctor --fix` repairs what it can. ``. Une configuration impossible à charger se termine aussi par `1` et une erreur `ccam:`.

## Fonctionnement

`internal/cli/doctor.go` charge `~/.config/ccam/config.toml` et appelle `doctor.Run` dans `internal/doctor/doctor.go`, qui exécute les vérifications dans un ordre fixe et renvoie les constats, chacun avec le compte auquel il appartient et le domaine vérifié. `tui.HealthReport` les regroupe dans les encadrés. `--fix` est transmis sous la forme `Options.Fix`, mis à false lors d'une simulation.

<!-- motion:doctor-checks -->

Le bloc central s'exécute une fois par compte. Les vérifications, et ce que fait `--fix` pour chacune :

1. **Version de Claude Code** (`checkClaude`). Résout le binaire `claude` comme le fait le lanceur et exécute `claude --version`. Un binaire absent, un `claude_bin` qui n'est pas un fichier exécutable ou un `--version` qui échoue est une erreur. Pas de réparation.
2. **Identifiants dans l'environnement** (`checkAuthEnv`). Avertit quand `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN` est défini : chaque compte lancé depuis ce shell en hérite, et Claude Code le préfère à la connexion propre du compte. Pas de réparation ; supprimez la variable vous-même.
3. **Dossier** (`checkProfile`). Un dossier de compte absent est une erreur et met fin aux vérifications de ce compte. Des permissions avec des bits de groupe ou d'autres utilisateurs donnent un avertissement ; `--fix` exécute `chmod 0700`. Un `~/.claude` importé qui a son propre `.claude.json` à côté de `~/.claude.json` donne une information : `ccam main` et `claude` seul gardent alors des informations de connexion et des serveurs MCP séparés. Windows n'a pas de bits de permission : cette vérification et celle de `0600` ci-dessous y sont donc ignorées ; le dossier situé dans votre profil utilisateur est protégé par sa liste de contrôle d'accès.
4. **Fichier de connexion** (`checkCredentials`). `.credentials.json` est vérifié avec `lstat`. Un lien symbolique est une erreur, car chaque compte doit avoir son propre jeton, et n'est pas réparé. Un mode autre que `0600` donne un avertissement ; `--fix` exécute `chmod 0600`. Ensuite la connexion est lue (sous macOS, d'abord depuis le trousseau). Une absence de connexion, ou un fichier de connexion qui n'en contient pas, donne un avertissement qui renvoie vers `ccam login <name>`. Un jeton de renouvellement expiré, ou une connexion que Claude Code a effacée après un refus du serveur, donne l'erreur `` login expired → `ccam login <name>` `` ; un jeton d'accès expiré n'est qu'une information, puisque Claude Code le renouvelle au lancement suivant. Une connexion conservée dans le Gestionnaire d'identification de Windows est une information, puisque ccam ne peut pas la lire. Une connexion valide affiche l'adresse e-mail avec laquelle le dossier est connecté maintenant et l'offre (`login valid (dev@company.com, Claude Team · Acme)`). Aucun de ces cas n'est réparé ; `ccam login --expired` reconnecte chaque compte qui en a besoin.
5. **Liens de partage** (`share.Check`). Compare chaque chemin de `share.Layout` avec les groupes du compte, en utilisant uniquement `lstat` et `readlink`. Il signale cinq cas : un lien absent, un lien qui pointe ailleurs, un vrai fichier là où devrait se trouver un lien, un lien restant pour un élément désormais local, et un lien correct dont le dossier de groupe est absent (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)` ; un lancement ne vérifie pas ce cas). Un lien absent et un dossier de groupe absent sont des avertissements, les autres cas sont des erreurs. Une copie depuis le groupe arrêtée à mi-chemin est aussi signalée : `--fix` met en place une copie préparée complète, mais quand le chemin a été recréé entre-temps (par un lancement), il nomme le dossier qui contient la copie complète pour que vous fassiez le tri. Les chemins ajoutés à un élément dans une version ultérieure de ccam (`file-history/`, `tasks/` et `paste-cache/` pour l'historique, `keybindings.json` pour les paramètres) ne donnent qu'un avertissement tant que le compte en a encore sa propre copie : `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Un lancement ne s'arrête jamais pour eux.
6. **Entrées inconnues** (`checkUnknownEntries`). Liste les noms de premier niveau du dossier du compte que ccam ne connaît pas (`share.Known`), comme indice d'un changement de structure de Claude Code. Information uniquement.
7. **Chemins des plugins** (`checkPlugins`). Les registres de plugins (`installed_plugins.json`, `known_marketplaces.json`) enregistrent des chemins absolus. Dans le `plugins/` de chaque groupe, les chemins qui passent par le dossier d'un compte cessent de fonctionner quand ce compte est supprimé ou renommé. `--fix` les réécrit vers le dossier du groupe avec `profile.RewritePluginPaths`, une écriture atomique du fichier de registre.
8. **Serveurs MCP** (`checkMCP`). Charge une fois le `mcp.json` de chaque groupe. Un fichier illisible est une erreur ; un serveur qui ne passe pas `mcp.Validate` donne un avertissement. Les serveurs du `~/.claude.json` de `claude` seul qu'aucun compte ne charge donnent un avertissement qui nomme la commande `ccam migrate --mcp` qui les copie (avec `--group <a|b>` quand plusieurs groupes partagent les serveurs MCP). Les serveurs qu'un seul compte d'un groupe possède donnent une information avec la commande `ccam mcp import` qui les partage, et un compte non connecté à des serveurs MCP auxquels d'autres comptes sont connectés reçoit un avertissement l'invitant à s'y connecter avec `/mcp` (ccam ne copie jamais les connexions). Pas de réparation.
9. **Format de la réponse d'utilisation** (`checkUsageFormat`). Interroge le point de terminaison d'utilisation, sans le cache, pour un compte après l'autre jusqu'à ce que l'un réponde, le tout dans une limite unique de 4 secondes. Avec `usage_source = "claude"`, il interroge plutôt Claude Code pour un compte. Si la réponse ne peut plus être analysée, il avertit que `ccam usage` se rabat sur le cache de la ligne d'état ; quand aucun compte n'a répondu, il le signale comme information (`could not check the usage response format: …`). C'est la seule vérification qui utilise le réseau.
10. **Raccourcis** (`checkShims`). Uniquement quand un compte a un raccourci. Des liens absents dans `~/.local/bin` donnent un avertissement, tout comme des raccourcis créés par ccam qui lancent encore un ancien ccam, remplacé ou déplacé (`shortcuts in <bin> that run another ccam (one replaced or moved): …`) ; `--fix` appelle `shim.Install`, qui ne crée et ne renouvelle que les liens que ccam a créés et enregistrés dans `~/.local/state/ccam/shortcuts.json`, et ne touche jamais aux autres fichiers. L'avertissement ne se termine par `→ fixed` que lorsque chacun de ces liens a été créé ; si l'un n'a pas pu l'être, la raison est ajoutée. Un raccourci qui lance un autre ccam encore installé fonctionne et donne une information ; `--fix` le fait pointer vers ce ccam. Un dossier bin absent du `PATH` donne un avertissement et n'est pas réparé.
11. **Dossiers** (`checkFolders`). Uniquement quand des entrées [`ccam folder`](/features/folder) existent. Une entrée dont le compte n'existe pas est une erreur, et une entrée dont le dossier a disparu donne un avertissement. Aucune n'est réparée, car supprimer une entrée changerait le compte qui démarre à cet endroit.
12. **Groupes et réglages** (`checkGroups`). Un membre de groupe qui n'est pas un compte, un compte dont le nom est depuis devenu une commande ccam (lancez-le avec `ccam run <name>` ou renommez-le), un `theme` ou un `usage_source` inconnu, et des lignes de `config.toml` que ccam ne lit pas donnent des avertissements. Ces dernières sont généralement un réglage écrit sous une table `[[group]]`, `[[profile]]` ou `[[folder]]`, où il appartient à cette table (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`) ; les commandes qui modifient la configuration refusent tant que ce n'est pas corrigé. Pas de réparation.
13. **Connexions** (`checkIdentities`). Lit dans son `.claude.json` sous quel compte chaque dossier est connecté. Un dossier désormais connecté à un autre compte que celui enregistré par `ccam login` (par exemple après un `/login` dans une session) donne un avertissement ; un compte sans enregistrement de connexion, comme un compte importé ou créé par un ancien ccam, est comparé d'après son adresse e-mail enregistrée. Un dossier connecté au même compte Claude qu'un autre donne aussi un avertissement : les deux partagent leurs limites d'utilisation. Un dossier qui contient le même jeton de connexion qu'un autre, un dossier copié (y compris le `~/.claude` de `claude` seul), est une erreur, car celui qui le renouvelle en premier déconnecte l'autre ; exécutez `ccam login <name>` pour lui donner le sien. Pas de réparation.

Avant les comptes, un **changement de groupes** arrêté à mi-chemin est une erreur dans l'encadré `General` : ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` le termine dans l'ordre prévu, avec les réponses données à ses conflits ; d'ici là, aucun autre changement de groupes, renommage ou suppression ne démarre. Les comptes concernés ne reçoivent que leurs liens manquants tant qu'il n'est pas terminé. S'il ne peut pas être terminé, l'erreur indique que supprimer ce fichier abandonne le reste du changement. `ccam --dry-run doctor --fix` montre l'ordre (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) et ce que la partie de chaque compte mettrait de côté. Un changement qu'un autre ccam effectue en ce moment est une information, et pendant qu'il s'exécute, `--fix` ne répare aucun lien de partage (`its share links are not repaired while another ccam changes the groups; …`). `--fix` détient `~/.local/state/ccam/group-change.lock` pendant qu'il répare, et relit `config.toml` après avoir terminé un changement.

### Réparer les liens de partage

Le chemin de lancement crée sur-le-champ un lien manquant, mais s'arrête pour tout ce qui pourrait faire perdre des données et renvoie ici :

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Pour tout problème de lien de partage, `--fix` appelle `share.Build` avec les paramètres de partage du compte comme état actuel et comme état cible, puis `Plan.Apply`. `Build` redéduit chaque étape de ce qui se trouve sur le disque, le même plan que celui d'une modification `ccam group`, si bien que la réparation suit les mêmes règles que l'entrée dans un groupe.

<!-- motion:doctor-repair -->

Ici, `settings.json` est un vrai fichier. Comme le groupe a déjà un `settings.json` différent, le conflit prend son choix par défaut, `KeepGroup`, et la copie du groupe est conservée. Le fichier du compte est déplacé vers `.ccam-backup/<time>/settings.json` dans le dossier du compte, et `os.Symlink` relie le chemin au groupe. `--fix` le signale : `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Autres cas :

- Si le groupe n'a pas encore de copie, le fichier du compte est déplacé dans le groupe et devient son premier contenu.
- Pour un dossier, les fichiers absents du groupe sont déplacés dans le groupe. Les fichiers identiques octet pour octet à ceux du groupe sont abandonnés, et ceux qui diffèrent vont dans la sauvegarde, nommés dans une note (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Un dossier de groupe absent est d'abord créé.
- Pour `history.jsonl`, le choix par défaut est `merge by line` : les lignes du groupe, puis vos nouvelles lignes. La copie précédente du groupe est enregistrée dans `.ccam-backup/<time>/group/`.
- Un lien vers un autre groupe est remplacé. Un lien vers un dossier extérieur à ccam va dans la sauvegarde, et son contenu alimente le groupe si celui-ci est vide.
- Un lien restant pour un élément désormais local est supprimé, et le compte reçoit une copie du contenu du groupe.

À part la suppression de liens et l'abandon des fichiers que le groupe contient déjà octet pour octet, chaque étape déplace des données vers le groupe ou le dossier de sauvegarde. Rien n'est perdu.

## Sécurité et cas particuliers

- Un simple `ccam doctor` ne modifie jamais rien. Dans les scripts, utilisez `--dry-run` avec `--fix` pour en être sûr.
- Pour un compte en cours d'exécution, quelle que soit la façon dont Claude Code a été lancé, `doctor --fix` crée seulement les liens manquants et les dossiers de groupe manquants, ce qui ne fait rien perdre. Quand il reste autre chose, il ajoute `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, car la session continue d'écrire dans les fichiers qu'une réparation déplacerait.
- Si une réparation de partage échoue, ccam signale `could not repair share links: <error>` pour ce compte et saute sa vérification des entrées inconnues. Les autres comptes sont tout de même vérifiés.
- Si une modification `ccam group` échoue à mi-parcours, les nouveaux groupes sont déjà enregistrés et le changement est noté dans `~/.local/state/ccam/group-change.json` ; ccam vous indique de terminer avec `ccam doctor --fix`, qui en effectue le reste, étapes MCP comprises.
- Un `.credentials.json` en lien symbolique, une connexion expirée et des identifiants dans l'environnement sont signalés mais jamais modifiés.
- Les réécritures des registres de plugins sont écrites de façon atomique mais ne sont pas copiées dans `.ccam-backup/`.

## Voir aussi

- [Partager les serveurs MCP, les paramètres et l'historique](/docs/sharing)
- [Copier une configuration existante](/features/migrate)
- [Limites d'utilisation et changement de compte](/docs/usage-limits)
- [Commandes](/docs/commands)
