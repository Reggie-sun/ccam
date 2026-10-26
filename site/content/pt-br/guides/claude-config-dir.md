---
title: CLAUDE_CONFIG_DIR explicado
description: O que o CLAUDE_CONFIG_DIR faz no Claude Code, o que fica na pasta, onde o macOS guarda o login e quais arquivos podem ser compartilhados entre contas.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "dc2fe42e11e1"
---

`CLAUDE_CONFIG_DIR` diz ao Claude Code qual pasta usar no lugar de `~/.claude`. Tudo o que é específico da conta fica lá: o login, as configurações, os servidores MCP, os plugins e o histórico de conversas. Apontar a variável para pastas diferentes dá a você contas separadas.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## O CLAUDE_CONFIG_DIR tem suporte oficial?

Sim. O Claude Code a lista entre as suas [variáveis de ambiente](https://code.claude.com/docs/en/env-vars) como a forma de substituir a pasta de configuração, "útil para executar várias contas lado a lado", e o [guia de autenticação](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) dele a usa para fazer login com várias contas. Você pode defini-la no shell ou nas configurações de usuário, mas não nas configurações de um projeto.

## O que tem dentro da pasta?

Verificado com o Claude Code 2.1.309:

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

**Sem `CLAUDE_CONFIG_DIR`, o `.claude.json` não fica dentro de `~/.claude`.** Nesse caso, o Claude Code o guarda em `~/.claude.json`, ao lado da pasta. Se você começar a usar `CLAUDE_CONFIG_DIR=~/.claude`, o Claude Code procura `~/.claude/.claude.json` e não encontra sua conta. `ccam import` resolve isso criando um link entre os dois. Se `~/.claude` já tem um `.claude.json` próprio, de uma execução anterior com `CLAUDE_CONFIG_DIR=~/.claude`, o ccam deixa os dois arquivos e avisa, já que só você sabe qual é o certo.

**No macOS, o login fica no Acesso às Chaves, não na pasta.** O Claude Code o guarda com o nome `Claude Code-credentials`. Quando `CLAUDE_CONFIG_DIR` está definido, o nome também termina com `-` e os 8 primeiros dígitos hexadecimais do SHA-256 do caminho da pasta. O nome da conta do item no Acesso às Chaves é `$USER`. Por isso, pastas separadas recebem entradas separadas no Acesso às Chaves, e uma pasta movida ou renomeada perde o login; é por isso que o `ccam rename` mantém o nome da pasta no macOS. Se o Acesso às Chaves recusar a gravação, o Claude Code recorre ao `.credentials.json` na pasta, como [a documentação dele](https://code.claude.com/docs/en/authentication#credential-management) descreve.

**Nem todo login fica na pasta.** Um login do Claude Console sem chave de API é guardado fora dela, então pastas separadas não conseguem manter dois desses logins separados. Logins do claude.ai e chaves de API são guardados por pasta.

## Onde fica a pasta no Windows?

Por padrão, ela fica em `%USERPROFILE%\.claude`. O login fica no `.credentials.json` dentro dela, ou no Gerenciador de Credenciais do Windows se o Claude Code estiver configurado para usá-lo, com um nome derivado do caminho da pasta, como no macOS. Defina `CLAUDE_CONFIG_DIR` com outra pasta para cada conta extra, do mesmo jeito que no Linux e no macOS.

## As pastas podem compartilhar arquivos?

Links simbólicos funcionam: o Claude Code lê e grava por meio de um `settings.json` ou de uma pasta `projects/` com link e mantém o link. Os servidores MCP são a exceção, porque dividem o `.claude.json` com o login; passe-os com `--mcp-config` em vez disso.

O ccam se baseia exatamente nessas regras. Veja [compartilhamento](/docs/sharing) para saber como ele cria links entre as pastas, e [o guia de várias contas](/guides/multiple-claude-code-accounts) para a visão geral.
