---
title: "Que faire quand Claude Code atteint sa limite"
description: "Attendez la réinitialisation ou poursuivez la même conversation sur un autre compte Claude Code. Les limites sur 5 heures et hebdomadaires, étape par étape."
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 049a46684c7b
---

Quand un compte atteint sa limite d'utilisation (souvent appelée « rate limit »), vous avez deux options : attendre que la limite se réinitialise, ou poursuivre dès maintenant la même conversation sur un autre compte qui vous appartient. Avec ccam, la seconde ne demande qu'une commande, `ccam next --auto --continue`.

## Quelles sont les limites de Claude Code ?

Avec les forfaits Pro et Max, Claude Code a une [limite de session sur cinq heures et une limite hebdomadaire](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Les deux sont partagées avec Claude sur le web et dans les applications, et l'utilisation compte pour les deux à la fois. Quand l'une est épuisée, Claude Code [affiche un message](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) comme « You've hit your session limit · resets 3:45pm ». Changer de modèle n'aide pas, car les limites couvrent tous les modèles.

## Option 1 : attendre la réinitialisation

Si vous n'avez qu'un compte, attendez. Les versions récentes de Claude Code [attendent dans la session ouverte](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) et reprennent la tâche d'elles-mêmes une fois la limite réinitialisée. `/usage` dans Claude Code indique quelle part des limites de votre forfait vous avez utilisée.

## Option 2 : poursuivre tout de suite sur un autre compte

Si vous avez un second compte, vous pouvez continuer immédiatement la même conversation. Configurez-le une fois :

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

La ligne d'état affiche désormais le compte et son utilisation sur 5 heures pendant que vous travaillez : la limite ne vous prend donc pas par surprise.

Quand la limite est atteinte, quittez Claude Code (`/exit`), puis :

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Ouvrez la conversation dans le compte qui a le plus de marge :

```sh
ccam next --auto --continue
```

Claude Code démarre en tant que `work` avec la conversation la plus récente du dossier courant, y compris tout ce qui a été dit jusque-là. Partager l'historique (`ccam group set default --add history` ci-dessus) en fait une seule conversation ; sans cela, Claude Code copie la conversation dans `work`, et ccam vous indique combien de jetons cette première réponse relit.

Avec le [plugin ccam](/docs/install#claude-code-plugin) dans Claude Code, une session qui s'arrête à sa limite affiche une notification et enregistre la fenêtre pleine, si bien que `ccam next --auto` ne choisit pas ce compte de nouveau avant sa réinitialisation, et `/ccam:switch` (ou `/ccam:switch work`) note la conversation pour `ccam next --continue`, sans rien laisser à deviner.

**Sièges Team et Enterprise.** Le choix automatique (`ccam next`, `--fallback`, `ccam pick`, la touche `n` du tableau de bord) passe aussi d'un siège à l'autre au sein d'une même organisation Claude Team ou Enterprise. Le centre d'aide d'Anthropic indique : « Each seat on a Team plan is assigned to one named individual, and only that person may use it » (chaque siège d'une offre Team est attribué à une personne nommément désignée, et seule cette personne peut l'utiliser) et « Users can't route their work through another user's account to use that user's remaining usage » (les utilisateurs ne peuvent pas faire passer leur travail par le compte d'un autre utilisateur pour profiter de l'utilisation restante de celui-ci) ([offres Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129) ; de même pour les [offres Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Utiliser le choix automatique entre plusieurs sièges peut donc enfreindre les conditions d'Anthropic. Renseignez-vous auprès de votre organisation et d'Anthropic avant de vous y fier. Les offres personnelles (Pro, Max) sont traitées dans la [note de la FAQ sur les conditions](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## ccam change-t-il de compte automatiquement ?

Non. ccam ne tourne pas pendant que vous travaillez : il lance Claude Code et se remplace par lui. C'est vous qui décidez quand changer, et `ccam next --auto` choisit le compte dont l'utilisation sur 5 heures est la plus faible, en ignorant les comptes à 100 %, non connectés ou d'une autre organisation. Avec `relay = true` dans `config.toml`, ccam attend plutôt à côté de la session, si bien que `/ccam:switch` enchaîne de lui-même, en conservant les options, le mode de permission et le niveau d'effort de la session. Certains autres outils échangent plutôt la connexion d'une session en cours ; [ccam et les autres sélecteurs de compte](/compare/account-switchers) explique la différence.

## Points à retenir

- Exécutez `ccam next --continue` dans le dossier où vous travailliez ; les conversations sont stockées par dossier de projet.
- Évitez d'utiliser deux comptes sur la même conversation en même temps. ccam vous avertit si vous essayez.
- Sans historique partagé, la conversation est copiée dans l'autre compte et les deux copies évoluent chacune de leur côté ; les points de contrôle à partir desquels `/rewind` restaure le code restent dans le premier compte.
- Le premier message sur le nouveau compte peut consommer plus de sa limite que d'habitude, car ce compte n'a pas encore la conversation en cache.
- Le droit de passer ainsi d'un compte à l'autre dépend des conditions de chaque compte.

Voir [Limites d'utilisation](/docs/usage-limits) pour toutes les options de `ccam usage` et `ccam next`.
