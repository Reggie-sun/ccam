---
title: "Comptes à clé d'API"
description: "ccam add --api ajoute un compte à clé d'API (API d'Anthropic ou passerelle) ; la clé est lue au lancement dans une variable ou une commande, jamais stockée."
updated: 2026-10-23
order: 13
translation: machine
sourceHash: "232a71525e42"
---

Un compte à clé d'API lance Claude Code avec une clé d'API au lieu d'une connexion à un abonnement Claude : l'API d'Anthropic, ou un point de terminaison qui parle le même format Messages, comme une passerelle LiteLLM ou OpenRouter. ccam ne stocke pas la clé. Il la lit à chaque lancement du compte, depuis une variable d'environnement ou une commande comme un gestionnaire de mots de passe, et ne la transmet qu'à ce processus Claude Code.

## Utilisation

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

Un compte à clé d'API rejoint des groupes comme n'importe quel autre compte (par défaut `default` : paramètres, `CLAUDE.md`, extensions et serveurs MCP), et fonctionne avec les raccourcis, `ccam folder` et `ccam ext`.

## Options

| Option | Effet | Par défaut |
| --- | --- | --- |
| `--api` | Ajouter un compte à clé d'API au lieu de se connecter | désactivé |
| `--key-from env:NAME` | Lire la clé dans la variable d'environnement `NAME` | |
| `--key-from command:COMMAND` | Exécuter `COMMAND` (avec `sh -c`, ou sous Windows tel que tapé dans `cmd /d /s /c "…"`, guillemets conservés) et utiliser la première ligne qu'elle affiche | |
| `--base-url <url>` | Le point de terminaison (`ANTHROPIC_BASE_URL`) | L'API d'Anthropic |
| `--auth api-key\|bearer` | Envoyer la clé comme `ANTHROPIC_API_KEY` (`x-api-key`) ou comme `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, ce que les passerelles attendent en général) | `api-key` |
| `--model <model>` | Le modèle (`ANTHROPIC_MODEL`) | Celui de Claude Code par défaut |

Les paramètres sont enregistrés dans `config.toml` :

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Fonctionnement

À chaque lancement du compte, après les vérifications habituelles et avant le démarrage de Claude Code :

1. ccam lit la clé : `env:` depuis l'environnement, `command:` en exécutant la commande avec une limite de 30 secondes. Une commande peut demander un mot de passe dans le terminal ; sa sortie d'erreur va aussi au terminal. Une clé vide arrête le lancement avec un message qui indique où la clé était attendue.
2. Il retire toute variable héritée qui détermine comment Claude Code se connecte ou quel fournisseur il utilise : les clés et jetons (`ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`), le point de terminaison (`ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`) et les sélecteurs de fournisseur (`CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`). Une clé définie dans votre shell pour autre chose ne peut pas atteindre ce compte, et une connexion à un abonnement ne peut pas l'emporter sur sa clé.
3. Il définit `ANTHROPIC_BASE_URL` et `ANTHROPIC_MODEL` s'ils sont donnés, et la clé comme `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN`.

Les comptes connectés à un abonnement gardent l'environnement dont ils héritent, comme avant. Dans une session interactive, Claude Code peut demander une fois s'il doit utiliser la clé qu'il trouve dans l'environnement.

## Sécurité et cas particuliers

- La clé n'est jamais écrite dans `config.toml`, un cache ou un journal, et `ccam --dry-run gw` ne la lit pas : une source `command:` n'est pas exécutée par une simulation. La simulation affiche à la place le point de terminaison, le modèle et la provenance de la clé :

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  Une variable héritée que le lancement retire apparaît au début sous la forme `env -u <name>`. Avec une source `env:` dont la variable est vide, la simulation échoue comme le ferait le lancement.
- `ccam add --api` vérifie une source `env:` tout de suite et avertit si la variable est vide ; une source `command:` est exécutée pour la première fois au lancement.
- `ccam login` et `ccam logout` refusent un compte à clé d'API, puisqu'il n'a pas de connexion. `ccam doctor` affiche son point de terminaison et la source de sa clé, et avertit quand une variable `env:` est vide.
- `ccam usage` affiche `no limits (API)` pour lui. `ccam next --auto`, `--fallback` et `ccam pick` ne choisissent que parmi les comptes à abonnement, car une clé d'API n'a pas de limite sur 5 heures ou hebdomadaire à comparer.
- Un `settings.json` partagé avec un bloc `env` ou `apiKeyHelper` s'applique à chaque compte du groupe, comptes à clé d'API compris. Gardez les paramètres de fournisseur hors des paramètres partagés, ou placez le compte à clé d'API dans un groupe sans `settings`.

## Voir aussi

- [Gérer les comptes](/features/accounts)
- [Lancer Claude Code](/features/launch)
- [Référence des commandes](/docs/commands)
