---
title: CLAUDE_CONFIG_DIR explicado
description: O que o CLAUDE_CONFIG_DIR faz no Claude Code, o que fica na pasta, onde o macOS guarda o login e quais arquivos podem ser compartilhados entre contas.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` diz ao Claude Code qual pasta usar no lugar de `~/.claude`. Tudo o que é específico da conta fica nela: o login, as configurações, os servidores MCP, os plugins e o histórico de conversas. Apontar a variável para pastas diferentes dá a você contas separadas.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## O que há dentro da pasta

Verificado com o Claude Code 2.1.302:

| Caminho | Conteúdo | Compartilhável |
| --- | --- | --- |
| `.credentials.json` | Tokens de login (Linux), tokens OAuth de MCP | Não |
| `.claude.json` | Dados da conta, servidores MCP do usuário, estado por projeto | Não (servidores MCP separadamente) |
| `settings.json` | Permissões, hooks, modelo, linha de status, plugins | Sim |
| `CLAUDE.md` | Instruções globais | Sim |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensões | Sim |
| `projects/` | Conversas como arquivos `.jsonl`, memória por projeto | Sim |
| `history.jsonl` | Histórico de prompts | Sim |
| `sessions/`, `shell-snapshots/`, `cache/` e outros | Estado de execução | Não |

## Dois detalhes que pegam as pessoas de surpresa

**Sem `CLAUDE_CONFIG_DIR`, o `.claude.json` não fica dentro de `~/.claude`.** Nesse caso, o Claude Code o guarda em `~/.claude.json`, ao lado da pasta. Se você começar a usar `CLAUDE_CONFIG_DIR=~/.claude`, o Claude Code passa a procurar `~/.claude/.claude.json` e não encontra sua conta. `ccam import` resolve isso criando um link entre os dois.

**No macOS, o login fica no Acesso às Chaves, não na pasta.** O Claude Code o armazena com o nome `Claude Code-credentials`, seguido de um hash curto do caminho da pasta quando `CLAUDE_CONFIG_DIR` está definido. Por isso, pastas separadas recebem itens separados no Acesso às Chaves. Se o Acesso às Chaves recusar a gravação, o Claude Code passa a usar o `.credentials.json` na pasta.

## Compartilhar entre pastas

Links simbólicos funcionam: o Claude Code lê e grava por meio de um `settings.json` ou de uma pasta `projects/` vinculados e mantém o link. Os servidores MCP são a exceção, porque compartilham o `.claude.json` com o login; passe-os com `--mcp-config`.

O ccam se baseia exatamente nessas regras. Consulte [compartilhamento](/docs/sharing) para ver como ele cria links entre pastas, e [o guia de várias contas](/guides/multiple-claude-code-accounts) para ter uma visão geral.
