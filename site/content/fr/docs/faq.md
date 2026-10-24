---
title: "Questions fréquentes"
metaTitle: "FAQ ccam : plusieurs comptes Claude Code"
description: "Réponses sur ccam, les comptes Claude Code, les limites d'utilisation, le partage, la confidentialité et les systèmes pris en charge."
updated: 2026-10-24
order: 6
translation: machine
sourceHash: c4e3b89bb4ea
---

## ccam est-il un outil officiel d'Anthropic ?

Non. ccam est un projet open source non officiel : il n'est ni affilié à Anthropic ni approuvé par Anthropic. Il utilise le paramètre `CLAUDE_CONFIG_DIR`, que [Claude Code documente](https://code.claude.com/docs/en/env-vars) pour utiliser plusieurs comptes côte à côte.

## Ai-je le droit d'utiliser plusieurs comptes Claude Code ?

Cela dépend des conditions de chaque compte et, pour les comptes professionnels, de la politique de votre organisation. Vérifiez-les avant d'utiliser plusieurs comptes. ccam se contente de faciliter le passage d'un compte à l'autre parmi ceux que vous possédez déjà.

**Sièges Team et Enterprise.** Le choix automatique (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Utiliser le choix automatique entre plusieurs sièges peut donc enfreindre les conditions d'Anthropic. Renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier. Les offres personnelles (Pro, Max) sont traitées dans le paragraphe ci-dessus.

## ccam permet-il de contourner les limites d'utilisation ?

Non. Chaque compte conserve ses propres limites. ccam les affiche et vous permet de continuer sur un autre compte qui vous appartient quand l'un d'eux est épuisé.

## Comment utiliser deux comptes Claude Code sur un seul ordinateur ?

Donnez à chaque compte son propre dossier de configuration avec `CLAUDE_CONFIG_DIR`. ccam le fait pour vous : exécutez `ccam add` pour chaque compte, puis lancez-les avec `ccam personal` et `ccam work`, ou avec des raccourcis que vous nommez vous-même, par exemple `ccwork`. Le [guide](/guides/multiple-claude-code-accounts) montre aussi comment procéder à la main.

## Comment poursuivre une conversation sur un autre compte ?

Exécutez `ccam next --continue` dans le dossier : il transporte la dernière conversation du dossier vers le compte suivant (`--auto` pour celui qui a le plus de marge). Dans Claude Code, `/ccam:switch` du [plugin ccam](/docs/install#claude-code-plugin) fait de même pour la session en cours. Quand les deux comptes partagent l'historique des conversations (`ccam group set default --add history`, ou un groupe à eux : `ccam group add pair --items history --members work,side`), le nouveau compte ouvre cette conversation même ; sinon, Claude Code la copie dans l'historique du nouveau compte.

## ccam voit-il ou envoie-t-il mes identifiants de connexion ?

ccam ne partage ni ne copie jamais votre connexion Claude entre les comptes, et ne copie jamais les connexions des serveurs MCP. Il lit le jeton d'accès d'un compte uniquement pour demander l'utilisation de ce compte au point de terminaison d'utilisation d'Anthropic, quand vous exécutez `ccam usage`, `ccam next --auto` ou `ccam doctor`, et n'en renouvelle un expiré qu'avec `ccam usage --refresh` ; avec `usage_source = "claude"`, il ne lit pas du tout le jeton. `ccam ext` et `ccam doctor` lisent les noms des connexions MCP d'un compte, jamais leurs jetons. ccam n'a aucune télémétrie. En dehors de ces consultations, la seule requête qu'il envoie de lui-même est la question du tableau de bord à GitHub sur la dernière version de ccam, qui ne transmet rien sur vous ni sur vos comptes.

## Comment se connecter sur un serveur distant ?

Exécutez `ccam add` ou `ccam login <name>` dans la session SSH, comme d'habitude. Claude Code attend la connexion sur le serveur : le navigateur de votre ordinateur aboutit donc sur une page `http://localhost:<port>/callback?code=…` qui ne se charge pas. Copiez cette adresse depuis le navigateur et collez-la dans le terminal : ccam la transmet à Claude Code sur le serveur, qui termine la connexion. L'adresse contient un code à usage unique : ne la partagez pas.

## Où mes comptes sont-ils stockés ?

Dans `~/.local/share/ccam/profiles/<name>/`, un dossier par compte. Sous macOS, Claude Code conserve la connexion elle-même dans le trousseau, sous un nom dérivé de ce dossier : `ccam rename` y conserve donc l'ancien nom du dossier.

## Qu'advient-il de mon dossier ~/.claude existant ?

Rien, sauf si vous exécutez `ccam import` ou `ccam migrate`. L'import enregistre le dossier là où il se trouve, si bien que `claude` seul et `ccam <name>` continuent d'utiliser le même compte et les mêmes paramètres. `ccam migrate` copie ses paramètres, `CLAUDE.md`, ses extensions et ses serveurs MCP dans les paramètres partagés, et ne fait que lire le dossier.

## Quels systèmes ccam prend-il en charge ?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine et d'autres), macOS et Windows 10 et 11, en natif ou dans WSL. Sous Windows, partager des fichiers comme `settings.json` nécessite le Mode développeur ; les dossiers partagés fonctionnent sans lui, et sans lui le premier `ccam add` crée le groupe `default` avec seulement les serveurs MCP et les extensions ([partage sous Windows](/features/share#on-windows)).

## Pourquoi un serveur MCP n'apparaît-il que dans un seul compte ?

Les serveurs ajoutés pendant une session avec `claude mcp add` sont enregistrés pour ce compte uniquement, et sans `-s user` pour le dossier courant uniquement. ccam vous le signale au prochain démarrage du compte, et `ccam mcp list` et `ccam ext` les affichent ; déplacez le serveur vers le groupe partagé avec `ccam mcp import <account> <server>`.

## ccam ralentit-il Claude Code ?

Lancer Claude Code via ccam ajoute environ 4 millisecondes, mesurées avec `scripts/perf-check.sh` sur un ordinateur portable. ccam lit un petit fichier, vérifie les liens partagés puis se remplace par `claude` : il ne tourne donc pas pendant que vous travaillez. Windows ne peut pas remplacer un processus : ccam y lance Claude Code et l'attend, sans rien faire d'autre.

## Comment mettre à jour ccam ?

Exécutez `ccam update`. La commande télécharge la dernière version, vérifie sa somme de contrôle et remplace le fichier `ccam` ; `ccam update --check` indique seulement si une version plus récente existe. Le tableau de bord vous signale aussi quand une version plus récente existe, et `U` l'installe après confirmation ; ccam ne se met jamais à jour de lui-même. Si vous avez installé le paquet `.deb`, installez plutôt le nouveau `.deb` avec apt. Voir [Mettre à jour](/docs/install#update).

## Comment désinstaller ccam ?

Supprimez le raccourci de chaque compte avec `ccam shortcut <account> --remove`. Si vous avez importé `~/.claude` et qu'il fait partie d'un groupe, exécutez d'abord `ccam remove <name>` pour lui, afin qu'il récupère une copie de ce qu'il partage. Supprimez ensuite `ccam` de `~/.local/bin`, puis `~/.config/ccam`, `~/.local/share/ccam` et `~/.local/state/ccam` si vous n'avez plus besoin des comptes. Voir [Installer ccam](/docs/install#uninstall).
