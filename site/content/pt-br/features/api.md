---
title: Contas com chave de API
description: ccam add --api adiciona uma conta com chave de API para a API da Anthropic ou um gateway; a chave vem de uma variável ou comando a cada início e não é salva.
updated: 2026-10-23
order: 13
translation: machine
sourceHash: "232a71525e42"
---

Uma conta com chave de API inicia o Claude Code com uma chave de API em vez do login de uma assinatura do Claude: a API da Anthropic, ou um endpoint que fala o mesmo formato Messages, como um gateway LiteLLM ou OpenRouter. O ccam não guarda a chave. Ele a lê toda vez que a conta inicia, de uma variável de ambiente ou de um comando como um gerenciador de senhas, e a entrega só a esse processo do Claude Code.

## Uso

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

Uma conta com chave de API entra em grupos como qualquer outra conta (por padrão, no `default`: configurações, `CLAUDE.md`, extensões e servidores MCP), e funciona com atalhos, `ccam folder` e `ccam ext`.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--api` | Adiciona uma conta com chave de API em vez de fazer login | desativado |
| `--key-from env:NAME` | Lê a chave da variável de ambiente `NAME` | |
| `--key-from command:COMMAND` | Executa `COMMAND` (com `sh -c`, ou no Windows como se fosse digitado em `cmd /d /s /c "…"`, mantendo as aspas) e usa a primeira linha que ele imprime | |
| `--base-url <url>` | O endpoint (`ANTHROPIC_BASE_URL`) | a API da Anthropic |
| `--auth api-key\|bearer` | Envia a chave como `ANTHROPIC_API_KEY` (`x-api-key`) ou como `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, que os gateways costumam pedir) | `api-key` |
| `--model <model>` | O modelo (`ANTHROPIC_MODEL`) | o padrão do Claude Code |

As configurações ficam em `config.toml`:

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

## Como funciona

A cada inicialização da conta, depois das verificações de sempre e antes de o Claude Code iniciar:

1. O ccam lê a chave: `env:` do ambiente, `command:` executando o comando com um limite de 30 segundos. Um comando pode pedir uma senha no terminal; a saída de erro dele também vai para o terminal. Uma chave vazia interrompe a inicialização com uma mensagem dizendo onde a chave era esperada.
2. Ele remove toda variável herdada que define como o Claude Code faz login ou qual provedor usa: as chaves e tokens (`ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`), o endpoint (`ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`) e as opções de provedor (`CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`). Uma chave definida no seu shell para outra coisa não chega a esta conta, e o login de uma assinatura não pode passar na frente da chave dela.
3. Ele define `ANTHROPIC_BASE_URL` e `ANTHROPIC_MODEL` quando informados, e a chave como `ANTHROPIC_API_KEY` ou `ANTHROPIC_AUTH_TOKEN`.

Contas com login de assinatura mantêm o ambiente que herdam, como antes. Em uma sessão interativa, o Claude Code pode perguntar uma vez se deve usar a chave que encontra no ambiente.

## Segurança e casos especiais

- A chave nunca é gravada em `config.toml`, em um cache ou em um log, e `ccam --dry-run gw` não a lê: uma simulação não executa uma origem `command:`. Em vez disso, a simulação mostra o endpoint, o modelo e de onde vem a chave:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  Uma variável herdada que a inicialização remove aparece como `env -u <name>` no começo. Com uma origem `env:` cuja variável está vazia, a simulação falha como a inicialização falharia.
- `ccam add --api` verifica uma origem `env:` na hora e avisa se a variável estiver vazia; uma origem `command:` só é executada pela primeira vez na inicialização.
- `ccam login` e `ccam logout` recusam uma conta com chave de API, já que ela não tem login. `ccam doctor` mostra o endpoint e a origem da chave dela, e avisa quando uma variável `env:` está vazia.
- `ccam usage` mostra `no limits (API)` para ela. `ccam next --auto`, `--fallback` e `ccam pick` escolhem só entre contas de assinatura, porque uma chave de API não tem limite de 5 horas ou semanal para comparar.
- Um `settings.json` compartilhado com um bloco `env` ou `apiKeyHelper` vale para todas as contas do grupo, incluindo as contas com chave de API. Mantenha as configurações de provedor fora das configurações compartilhadas, ou coloque a conta com chave de API em um grupo sem `settings`.

## Relacionados

- [Gerenciar contas](/features/accounts)
- [Iniciar o Claude Code](/features/launch)
- [Referência de comandos](/docs/commands)
