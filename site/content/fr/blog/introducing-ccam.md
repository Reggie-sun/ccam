---
title: "Présentation de ccam"
description: "Pourquoi ccam existe, ce qu'il apporte aux personnes qui ont plusieurs comptes Claude Code et comment il se fait oublier."
updated: 2026-10-26
order: 1
translation: machine
sourceHash: e4370f29b818
---

ccam est un petit outil en ligne de commande pour les personnes qui utilisent plus d'un compte Claude Code sur le même ordinateur : un compte personnel et un compte professionnel, ou plusieurs comptes pour un même projet de longue durée.

## Le problème

Claude Code prend en charge plusieurs comptes grâce à [`CLAUDE_CONFIG_DIR`](/guides/claude-config-dir) : un dossier par compte. En pratique, chaque dossier est une île. Les paramètres, `CLAUDE.md`, les serveurs MCP, les agents et les skills doivent être reconfigurés dans chaque dossier et synchronisés à la main. Les conversations restent là où elles ont commencé. Et quand un compte atteint sa limite d'utilisation, vous le découvrez en essayant.

## Ce que fait ccam

- `ccam add` crée le dossier, lance la connexion et, si vous le souhaitez, attribue au compte un raccourci que vous nommez vous-même, par exemple `ccwork`.
- `ccam work` et `ccwork` lancent Claude Code avec ce compte et transmettent tous les arguments tels quels.
- Les serveurs MCP, les paramètres, `CLAUDE.md` et les extensions sont partagés au moyen de liens, élément par élément et compte par compte. Les connexions ne sont jamais partagées.
- `ccam usage` affiche l'utilisation sur 5 heures et sur la semaine de chaque compte ; `ccam next --continue` reprend la conversation dans le compte suivant.

## Se faire oublier

Un outil que l'on exécute avant chaque session ne devrait rien coûter de perceptible. ccam lit un fichier de configuration, vérifie les liens partagés à raison d'un ou deux appels système par lien, puis se remplace par `claude`. Il ajoute environ 4 millisecondes et, par défaut, ne tourne pas pendant que vous travaillez.

Il prend aussi soin de vos fichiers. Chaque modification du partage affiche un plan avant de commencer et conserve les originaux dans un dossier de sauvegarde. `ccam doctor` vérifie les liens, les permissions et les connexions, et répare ce qui peut l'être sans risque.

## Essayer ccam

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

Le [démarrage rapide](/docs/quick-start) montre comment ajouter deux comptes, et [comment utiliser plusieurs comptes Claude Code](/guides/multiple-claude-code-accounts) explique l'approche avec et sans ccam.

ccam est open source et non officiel : il n'est ni affilié à Anthropic ni approuvé par Anthropic. Le droit d'utiliser plus d'un compte dépend des conditions qui s'appliquent à chacun. Les retours et les traductions sont les bienvenus sur [GitHub](https://github.com/jhyunwoo/ccam).
