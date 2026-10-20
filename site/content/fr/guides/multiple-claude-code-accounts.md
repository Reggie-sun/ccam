---
title: "Plusieurs comptes Claude Code sur un même ordinateur"
description: "Utilisez plusieurs comptes Claude Code côte à côte avec CLAUDE_CONFIG_DIR, à la main avec des alias shell ou avec ccam, et changez de compte facilement."
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e2f4deb40c81
---

Donnez à chaque compte son propre dossier de configuration. Claude Code lit ce dossier dans la variable d'environnement `CLAUDE_CONFIG_DIR`, et chaque dossier contient sa propre connexion, ses paramètres et son historique : chaque compte reste connecté et vous pouvez les utiliser en même temps. Vous pouvez configurer cela à la main avec un alias shell, ou laisser ccam s'en charger.

## Comment le configurer à la main ?

C'est l'approche que [décrit la documentation de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Ajoutez une ligne par compte supplémentaire dans `~/.bashrc` ou `~/.zshrc` :

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Ouvrez un nouveau terminal et exécutez `claude-work`. Claude Code vous demande de vous connecter, car le dossier est nouveau. À partir de là, `claude` est votre premier compte et `claude-work` le second. Ajoutez un alias par compte, chacun avec son propre dossier.

Sous Windows sans WSL, une fonction PowerShell dans votre profil (`notepad $PROFILE`) fait la même chose :

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Que manque-t-il à l'approche par alias ?

Elle fonctionne bien pour deux comptes. À savoir :

- Chaque dossier part de zéro : pas de paramètres, pas de `CLAUDE.md`, pas de serveurs MCP, pas d'agents ni de skills. Vous copiez ou reliez ce dont vous avez besoin, et vous gardez les copies synchronisées vous-même.
- Les alias n'existent que dans les shells interactifs. Les scripts, les Makefiles et les autres programmes ne les voient pas.
- Les serveurs MCP se trouvent dans `.claude.json`, à côté de la connexion : on ne peut donc pas simplement les relier d'un dossier à l'autre.
- Les conversations restent dans le compte où elles ont commencé : vous ne pouvez donc pas en poursuivre une dans un autre compte.
- Rien ne montre l'utilisation de tous les comptes à la fois.

## Peut-on plutôt changer de compte avec /login ?

Oui, si vous ne changez que de temps en temps. Taper `/login` dans Claude Code [change le compte](https://code.claude.com/docs/en/quickstart) du dossier où vous vous trouvez. Un seul compte est connecté à la fois : vous vous reconnectez donc à chaque changement, et vous ne pouvez pas utiliser deux comptes côte à côte. Des dossiers séparés évitent ces deux problèmes.

## Comment le configurer avec ccam ?

ccam automatise la même idée et ajoute ce qui manque :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Désormais, `ccam personal` et `ccam work` lancent les deux comptes, tout comme `ccwork`, le raccourci que vous avez donné au compte work. Comme ce sont de vrais programmes et non des alias, ils fonctionnent aussi dans les scripts. Exécutez `ccam` seul pour choisir un compte dans une liste, et `ccam use work` pour en faire le compte par défaut.

Les paramètres, `CLAUDE.md`, les serveurs MCP et les extensions sont partagés dès le départ, et vous pouvez activer le partage de l'historique des conversations pour poursuivre une conversation sur l'autre compte avec `ccam next --continue`. `ccam usage` affiche l'utilisation sur 5 heures et sur la semaine de tous les comptes à la fois.

Si vous utilisez déjà `~/.claude`, exécutez `ccam import` au lieu d'ajouter de nouveau ce compte, ou `ccam migrate` pour copier ses paramètres, ses serveurs MCP et ses extensions dans les paramètres partagés.

## Comment savoir quel compte est actif ?

Dans Claude Code, `/status` affiche le compte auquel vous êtes connecté. Avec ccam, exécutez une fois `ccam statusline install` : la ligne d'état affiche en permanence le nom du compte et son utilisation, par exemple `[work] · Opus · 5h 24% · 7d 41%`.

## Est-ce que cela fonctionne sous Windows et WSL ?

Claude Code fonctionne [nativement sous Windows et dans WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` fonctionne dans les deux cas : sous Windows natif, l'alias ou la fonction PowerShell ci-dessus suffit. ccam fonctionne sous Linux et macOS ; sous Windows, installez-le et utilisez-le dans WSL, où il se comporte comme sous Linux.

## Dépannage

**Claude Code me demande de me connecter à chaque fois.** Le premier démarrage dans un nouveau dossier le demande toujours. S'il continue à le demander, la connexion a expiré ou a été supprimée : exécutez `ccam login work` (ou lancez l'alias une fois et connectez-vous).

**Un serveur MCP n'apparaît que dans un seul compte.** Les serveurs ajoutés pendant une session avec `claude mcp add` appartiennent à ce compte. Pour en rendre un disponible dans tous les comptes, utilisez `ccam mcp import work <server>`.

**Les paramètres diffèrent d'un compte à l'autre.** `ccam share work` affiche ce que le compte partage ; `ccam share work --settings on` partage ses paramètres, après avoir montré ce qui va changer. `ccam doctor --fix` répare les liens cassés et garde une sauvegarde de chaque fichier qu'il modifie.

**Deux comptes ouvrent la même conversation.** Évitez de travailler sur une même conversation dans deux comptes en même temps. ccam vous avertit si vous essayez.

## Lequel choisir ?

Pour deux comptes sans rien de partagé, un alias suffit. Si vous voulez partager paramètres et serveurs MCP, avoir plus de comptes, voir l'utilisation d'un coup d'œil ou poursuivre des conversations, ccam vous épargne cette gestion. La [comparaison avec les alias](/compare/aliases) présente les différences côte à côte, et [ccam et les autres sélecteurs de compte](/compare/account-switchers) couvre les outils qui échangent plutôt une seule connexion.

Le droit d'utiliser plusieurs comptes dépend des conditions de chaque compte. Vérifiez-les avant de vous lancer.
