---
title: "ccam et les autres sélecteurs de compte Claude Code"
description: "ccam comparé à claude-swap, CCS, clausona, clauth, CAAM, cux et d'autres outils pour changer de compte Claude Code, et quelle approche pour quel usage."
updated: 2026-10-25
order: 2
translation: machine
sourceHash: a6aa0409a8df
---

Les outils pour changer de compte Claude Code suivent l'une de deux approches. La plupart échangent la connexion à l'intérieur de l'unique dossier `~/.claude` : un seul compte est donc actif à la fois. D'autres, dont ccam, donnent à chaque compte son propre dossier avec `CLAUDE_CONFIG_DIR` : chaque compte reste connecté et plusieurs peuvent tourner en même temps. ccam y ajoute le partage des paramètres et des serveurs MCP, et poursuit une conversation dans le dossier d'un autre compte.

## Échanger une connexion ou séparer les dossiers ?

**Échanger la connexion** conserve un seul jeu de paramètres, de serveurs MCP et d'historique, puisqu'il n'y a qu'un dossier. Le changement est rapide, et une session en cours peut reprendre la nouvelle connexion. Mais un seul compte est actif à la fois, pour tous les terminaux et éditeurs en même temps.

Pour faire l'échange, l'outil garde une copie de la connexion de chaque autre compte et la remet en place. Claude Code remplace le jeton d'actualisation d'une connexion quand il la renouvelle : chaque copie enregistrée doit donc rester synchronisée avec les sessions en cours. Des comptes déconnectés, ou qui demandent à nouveau `/login` après un changement, font partie des problèmes que ces projets corrigent le plus souvent.

**Des dossiers séparés**, c'est la méthode que [décrit la documentation de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Chaque compte reste connecté, et vous pouvez lancer côte à côte une session professionnelle et une session personnelle. Chaque connexion reste dans son propre dossier, où Claude Code la renouvelle comme d'habitude. La contrepartie : chaque dossier démarre avec ses propres paramètres, serveurs MCP et historique, sauf si un outil les partage.

## Comparaison des outils

D'après la documentation et le code source de chaque projet au 2026-10-20 (clausona : v0.12.0-beta). « Non documenté » signifie que le projet n'en parle pas, pas que la fonctionnalité manque.

| Outil | Approche | Plusieurs comptes à la fois | Utilisation sur 5 heures et hebdomadaire | Quand une limite est atteinte | Plateformes |
| --- | --- | --- | --- | --- | --- |
| ccam | Dossiers séparés | Oui | Oui, `ccam usage` | Une notification dans la session ; `ccam next --auto --continue` ou `/ccam:switch` transporte la conversation vers le compte qui a le plus de marge, avec ou sans historique partagé, et avec `auto_switch` activé (Linux, macOS) la session y passe d'elle-même ; `ccam --fallback` y relance les tâches `-p` de lui-même, et `ccam pick` répartit les workers | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Échange la connexion ; dossiers séparés avec `cswap run` | Oui, avec `cswap run` | Oui | `cswap auto` change de connexion à un seuil ; la session en cours continue | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Dossiers séparés | Oui | En partie | Des groupes de comptes transmettent les requêtes au compte suivant via un proxy local | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Dossiers séparés | Oui | Oui | Des routes choisissent un compte qui a du quota avant chaque exécution, peuvent attendre une réinitialisation et attribuent des comptes distincts aux workers | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Échange la connexion ; dossiers séparés avec `clauth start` | Oui, avec `clauth start` | Oui | Passe au compte suivant d'une chaîne de secours | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Échange la connexion, ou profils isolés | Oui, avec des profils isolés | Oui | `caam run` relance les commandes non interactives sur un autre compte | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Échange la connexion et enveloppe `claude` | Non, un seul compte actif | Oui | Échange la connexion et reprend la même conversation | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Échange la connexion | Non | Oui | Changement facultatif à un seuil | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Échange la connexion | Non | Oui | Changement facultatif à 100 % | macOS |

ccswap est un fork de claude-swap qui gère aussi les connexions Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) et [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) donnent eux aussi à chaque compte son propre dossier. Plusieurs de ces outils prennent aussi en charge d'autres agents de code, comme Codex ou Gemini.

## Comment ils sont construits

| Outil | Écrit en | Pendant que Claude Code tourne | Copie les connexions d'un endroit à l'autre | Ce que les comptes peuvent partager |
| --- | --- | --- | --- | --- |
| ccam | Go, un seul binaire statique | Rien par défaut : ccam se remplace par `claude` (sous Windows, qui ne le permet pas, il l'attend) ; avec `relay` ou `auto_switch` activé, ccam attend à côté des sessions interactives | Non | Des groupes : chacun nomme les éléments qu'il partage et ses comptes ; un compte peut appartenir à plusieurs |
| claude-swap | Python 3.12+, installé avec uv ou pipx | Rien avec `cswap run` ; `cswap auto` et l'app de barre des menus continuent de tourner | Oui, quand il échange la connexion | Avec `cswap run` : paramètres, `CLAUDE.md` et extensions ensemble, historique en option ; serveurs MCP copiés à chaque lancement |
| CCS | Node.js 18+ ou Bun | CCS reste le processus parent ; les groupes de comptes ajoutent un proxy local | Non pour ses dossiers de compte ; les groupes de comptes gardent leurs propres connexions | Partagé ou local par compte, historique dans des groupes de contexte ; serveurs MCP copiés |
| clausona | TypeScript, Node.js 20+ | `claude` seul passe par une fonction shell ; `clausona run` garde Node comme processus parent ; un processus Node en arrière-plan après chaque sortie | Non, mais il renouvelle lui-même les jetons des comptes inactifs pour lire leur utilisation | Tout sauf les fichiers propres à chaque compte, depuis un compte principal ; historique facultatif ; les serveurs ajoutés avec `claude mcp add` ne sont pas partagés |
| clauth | Rust | clauth reste le processus parent | Oui, quand il échange la connexion | Tout `~/.claude`, ou rien |
| CAAM | Go | Selon le mode ; un démon facultatif | Oui, en mode coffre | Fixé par le mode ; serveurs MCP copiés à chaque lancement |
| cux | Go | Un wrapper reste autour de `claude` | Oui | Tout : un seul dossier |
| CCSwitcher, Claude Switcher | Apps macOS (Swift, Python) | L'app de barre des menus | Oui | Tout : un seul dossier |

## Les points forts de ccam

### Les connexions restent là où Claude Code les a mises

ccam ne copie jamais une connexion Claude d'un endroit à un autre. La connexion de chaque compte reste dans son propre dossier, ou sous macOS dans l'élément du trousseau propre à ce dossier, et c'est Claude Code qui la renouvelle. Il n'existe aucune copie enregistrée qui puisse se désynchroniser. Un `~/.claude` existant enregistré avec `ccam import` garde l'élément du trousseau qu'utilise `claude` seul : il reste connecté dans les deux cas.

ccam ne renouvelle lui-même un jeton expiré que lorsque vous lancez `ccam usage --refresh`. Il prend d'abord un verrou et refuse tant que ce compte est en cours d'exécution, quelle que soit la façon dont Claude Code a été lancé (sessions en arrière-plan et leur démon compris), et tant que Claude Code renouvelle lui-même cette connexion : il ne fait donc jamais tourner un jeton sous une session active.

### Par défaut, rien ne s'intercale entre vous et Claude Code

ccam lit un fichier, vérifie ses liens et se remplace par `claude` (`exec`). Par défaut, pendant que vous travaillez, il n'y a ni processus ccam, ni wrapper, ni proxy, ni démon : l'entrée standard, les signaux et le code de sortie appartiennent à Claude Code, et ses requêtes vont directement chez Anthropic. Windows ne peut pas remplacer un processus : ccam y lance Claude Code dans la même console et attend son code de sortie, sans rien faire d'autre. Avec `relay = true` ou `auto_switch = true`, ccam attend de la même façon à côté des sessions interactives, pour pouvoir poursuivre une conversation sous un autre compte ; les requêtes de Claude Code vont toujours directement chez Anthropic.

Lancer un compte de cette façon ajoute environ 4 ms, moins de 6 ms au 95e centile, avec 8 Mo de mémoire au pic ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) avec hyperfine, 300 exécutions). Sur le même portable, le simple démarrage de l'interpréteur Python prend 13 ms et celui de Node.js 31 ms, avant qu'un programme ne fasse quoi que ce soit.

### Les groupes décident de ce qui est partagé

Les serveurs MCP, `settings.json`, `CLAUDE.md`, les extensions (agents, commandes, skills, plugins) et l'historique des conversations sont cinq éléments distincts. Un groupe nomme les éléments qu'il partage et les comptes qu'il contient : deux comptes professionnels peuvent ainsi partager l'historique dans un groupe `work` pendant qu'un compte personnel garde le sien, et tous les trois partagent les mêmes serveurs MCP via `default`. `ccam group set work --add memory` déplace tous les comptes du groupe en un seul plan ; le [tableau de bord](/features/dashboard) fait de même avec des cases à cocher et grise un élément qu'un compte reçoit déjà d'un autre groupe.

- Les serveurs MCP sont passés avec `--mcp-config` au lancement : par défaut, un lancement ne touche donc pas à `.claude.json`, le fichier que Claude Code réécrit pendant qu'il tourne. ccam n'y écrit que lors d'étapes qui l'annoncent, sous le verrou de Claude Code et avec une sauvegarde : pour Remote Control après avoir demandé, quand un compte quitte ou rejoint le partage MCP, et pour `ccam mcp import`. `ccam mcp add` ajoute un serveur une seule fois pour tout le groupe.
- Les autres éléments sont des liens vers le dossier du groupe : une modification faite dans un compte apparaît aussitôt dans les autres. Rien n'est copié au lancement.
- Les connexions, `.claude.json` et l'état d'exécution ne sont jamais partagés.

### Chaque modification est montrée d'abord et sauvegardée

`ccam group` et `ccam migrate` affichent un plan et demandent confirmation avant de déplacer quoi que ce soit. Quand un fichier diffère entre un compte et son groupe, vous choisissez : garder celui du groupe, utiliser le vôtre ou fusionner le JSON ; dans les dossiers, la copie du groupe reste et le plan nomme les fichiers qui diffèrent. Les originaux vont dans `.ccam-backup/<time>/`, un compte supprimé est déconnecté et son dossier part dans un dossier de corbeille, et les fichiers sont remplacés en un seul renommage : une écriture interrompue laisse l'ancien fichier en place. `--dry-run` montre ce que feraient le lancement, le partage, la migration, les raccourcis, les mises à jour et les réparations. `ccam doctor --fix` répare les liens cassés et les permissions.

### La même conversation dans le compte suivant

`ccam next --auto --continue` interroge tous les comptes en même temps sur leur utilisation, ouvre celui qui a le plus de marge et poursuit la conversation. Avec l'historique partagé, c'est la même conversation ; sans lui, Claude Code la copie dans le nouveau compte (`--resume <file> --fork-session`), et ccam indique combien de jetons la première réponse relit. Quand une limite arrête une session, le hook du plugin ccam affiche une notification, et `/ccam:switch` passe la conversation au compte suivant. Par défaut, c'est vous qui décidez quand changer. Avec `auto_switch` activé (Linux, macOS), ccam déplace la session lui-même, mais seulement après que Claude Code l'a arrêtée à une limite, et seulement vers un compte du même type que celui qu'elle quitte : un autre siège de la même organisation, ou une autre offre personnelle ([détails](/features/next#auto-switch)). Rien ne surveille votre utilisation en arrière-plan, et par défaut les choix automatiques ne font jamais sortir le travail d'une organisation Team ou Enterprise ni ne l'y font entrer : une conversation d'entreprise n'atterrit donc jamais par accident sur un compte privé. La ligne d'état suggère la commande dès qu'un compte atteint 90 %.

Les scripts ne peuvent pas vous attendre : `ccam --fallback work -p "…"` change donc de compte pour eux. Quand le compte atteint sa limite, la même commande est relancée avec le compte qui a le plus de marge, et seule la sortie de l'exécution qui a abouti est affichée. Un échec qui n'est pas dû à une limite d'utilisation est transmis tel quel. Pour des workers en parallèle, `ccam pick` donne à chacun un compte différent qui a de la marge, et les deux quittent avec le code 75 quand plus aucun compte n'en a.

**Sièges Team et Enterprise.** Le choix automatique (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord, et [`auto_switch`](/features/next#auto-switch) lorsque vous l'activez) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Utiliser le choix automatique entre plusieurs sièges peut donc enfreindre les conditions d'Anthropic. C'est à vous de décider si vous l'utilisez dans ce cas : renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier. Les offres personnelles (Pro, Max) sont traitées dans la [note de la FAQ sur les conditions](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

### Le bon compte dans chaque dossier

`ccam folder work ~/code/company` fait que tout lancement sans nom de compte dans ce dossier utilise `work`, tandis que vos autres projets gardent le compte par défaut. Le réglage se trouve dans le même `config.toml` que lit chaque lancement : il ne coûte donc aucun temps.

### Un seul binaire, des téléchargements vérifiés

ccam est un binaire statique d'environ 10 Mo, sans environnement d'exécution à installer, pour Linux, macOS et Windows. Les scripts d'installation (`install.sh`, `install.ps1`) et `ccam update` vérifient chaque téléchargement avec les sommes de contrôle SHA-256 de la version, et ccam ne se met jamais à jour sans demander. Il n'y a aucune télémétrie : les seules requêtes qu'il envoie sont les consultations d'utilisation de vos comptes auprès d'Anthropic, en parallèle avec une limite de 3 secondes et un cache de 60 secondes, la question que le tableau de bord pose à GitHub sur la dernière version quand il s'ouvre, et les téléchargements que vous demandez. Avec `usage_source = "claude"`, même celles-ci passent par Claude Code, et ccam n'envoie aucun jeton.

## Ce que les autres outils font de plus

- claude-swap, clauth, cux, CAAM et CCSwitcher peuvent changer de compte d'eux-mêmes dans une session interactive en cours quand un compte dépasse un seuil d'utilisation. ccam change de compte pour les exécutions `-p` avec `--fallback`, et ne déplace une session interactive de lui-même qu'avec `auto_switch` activé (Linux, macOS) et seulement après que Claude Code l'a arrêtée à une limite ; sinon, c'est vous qui la déplacez, avec une notification à la limite et `/ccam:switch`.
- Sous Windows sans le Mode développeur, ccam partage les dossiers mais ne peut pas partager de fichiers isolés comme `settings.json` (le premier `ccam add` ne partage alors que les serveurs MCP et les extensions) ; claude-swap et CCS les y copient à la place.
- CCSwitcher, Claude Switcher et claude-swap affichent l'utilisation dans la barre des menus de macOS.
- CCS, clausona, CAAM, clauth, ccswap et Claude Switcher gèrent aussi Codex et d'autres agents de code.
- clausona active et désactive aussi les extensions compte par compte (`ccam ext` de ccam se contente de les afficher), suit lui-même le coût local en jetons (ccam confie cela à [ccusage](https://ccusage.com) : `ccam ccusage`) et fournit des plugins Claude Code qui font tourner des flottes de workers dans tmux et d'autres terminaux.
- CAAM, CCSwitcher et Claude Switcher s'installent avec Homebrew, et clauth et CAAM signent les sommes de contrôle de leurs versions. ccam s'installe avec son script, un `.deb` ou `go install`, et vérifie les sommes de contrôle sans signature.

## Lequel choisir ?

Choisissez un outil qui échange la connexion si vous utilisez un compte à la fois ou si vous voulez que les sessions en cours changent de compte avant la limite ou en échangeant la connexion. Choisissez ccam si vous voulez des comptes côte à côte, une seule configuration partagée, la même conversation dans le compte suivant et, par défaut, rien qui tourne à côté de Claude Code, sous Linux, macOS ou Windows. Si vous n'avez besoin de rien partager, de [simples alias shell](/compare/aliases) suffisent.

Le droit d'utiliser plusieurs comptes dépend des conditions de chaque compte. Vérifiez-les avant de vous lancer.
