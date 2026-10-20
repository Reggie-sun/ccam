---
title: CLAUDE_CONFIG_DIR explicado
description: O que o CLAUDE_CONFIG_DIR faz no Claude Code, o que fica na pasta, onde o macOS guarda o login e quais arquivos podem ser compartilhados entre contas.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: ef9317030e42
---

`CLAUDE_CONFIG_DIR` diz ao Claude Code qual pasta usar no lugar de `~/.claude`. Tudo o que é específico da conta fica lá: o login, as configurações, os servidores MCP, os plugins e o histórico de conversas. Apontar a variável para pastas diferentes dá a você contas separadas.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## O CLAUDE_CONFIG_DIR tem suporte oficial?

Sim. O Claude Code a lista entre as suas [variáveis de ambiente](https://code.claude.com/docs/en/env-vars) como a forma de substituir a pasta de configuração, "útil para executar várias contas lado a lado", e o [guia de autenticação](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) dele a usa para fazer login com várias contas. Você pode defini-la no shell ou nas configurações de usuário, mas não nas configurações de um projeto.

## O que tem dentro da pasta?

Verificado com o Claude Code 2.1.302:

| Caminho | Conteúdo | Compartilhável |
| --- | --- | --- |
| `.credentials.json` | Tokens de login (Linux, Windows), tokens OAuth de MCP | Não |
| `.claude.json` | Dados da conta, servidores MCP do usuário, estado por projeto | Não (servidores MCP à parte) |
| `settings.json` | Permissões, hooks, modelo, linha de status, plugins | Sim |
| `CLAUDE.md` | Instruções globais | Sim |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensões | Sim |
| `projects/` | Conversas como arquivos `.jsonl`, memória por projeto | Sim |
| `history.jsonl` | Histórico de prompts | Sim |
| `sessions/`, `shell-snapshots/`, `cache/` e outros | Estado de execução | Não |

## O que costuma pegar as pessoas de surpresa?

**Sem `CLAUDE_CONFIG_DIR`, o `.claude.json` não fica dentro de `~/.claude`.** Nesse caso, o Claude Code o guarda em `~/.claude.json`, ao lado da pasta. Se você começar a usar `CLAUDE_CONFIG_DIR=~/.claude`, o Claude Code procura `~/.claude/.claude.json` e não encontra sua conta. `ccam import` resolve isso criando um link entre os dois.

**No macOS, o login fica no Acesso às Chaves, não na pasta.** O Claude Code o guarda com o nome `Claude Code-credentials`, seguido de um hash curto do caminho da pasta quando `CLAUDE_CONFIG_DIR` está definido. Por isso, pastas separadas recebem entradas separadas no Acesso às Chaves. Se o Acesso às Chaves recusar a gravação, o Claude Code recorre ao `.credentials.json` na pasta, como [a documentação dele](https://code.claude.com/docs/en/authentication#credential-management) descreve.

**Nem todo login fica na pasta.** Um login do Claude Console sem chave de API é guardado fora dela, então pastas separadas não mantêm dois desses logins separados. Logins do claude.ai e chaves de API são guardados por pasta.

## Onde fica a pasta no Windows?

Por padrão em `%USERPROFILE%\.claude`, com o login no `.credentials.json` dentro dela. Defina `CLAUDE_CONFIG_DIR` com outra pasta para cada conta extra, do mesmo jeito que no Linux e no macOS.

## As pastas podem compartilhar arquivos?

Links simbólicos funcionam: o Claude Code lê e grava por meio de um `settings.json` ou de uma pasta `projects/` com link e mantém o link. Os servidores MCP são a exceção, porque dividem o `.claude.json` com o login; passe-os com `--mcp-config` em vez disso.

O ccam se baseia exatamente nessas regras. Veja [compartilhamento](/docs/sharing) para saber como ele cria links entre as pastas, e [o guia de várias contas](/guides/multiple-claude-code-accounts) para a visão geral.
