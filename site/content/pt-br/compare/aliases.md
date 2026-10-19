---
title: ccam comparado com aliases de shell
description: Aliases de shell com CLAUDE_CONFIG_DIR e o ccam lado a lado, para executar várias contas do Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6fc4200a8a65
---

Um alias de shell por conta é a forma mais simples de executar várias contas do Claude Code. O ccam faz a mesma coisa por baixo e acrescenta compartilhamento, troca de conta e acompanhamento do uso.

| | Alias de shell | ccam |
| --- | --- | --- |
| Logins separados | Sim | Sim |
| Configuração por conta | Editar `~/.bashrc`, fazer login | `ccam add` |
| Funciona em scripts e outros programas | Não, aliases só existem no shell | Sim, `ccam work` e atalhos como `ccwork` são comandos de verdade |
| Configurações, `CLAUDE.md` e extensões compartilhados | Copiar ou criar links manualmente | Por item, com backups |
| Servidores MCP compartilhados | Manualmente, por conta | `ccam mcp add` uma vez |
| Continuar uma conversa em outra conta | Não | `ccam next --continue` |
| Uso de todas as contas | Não | `ccam usage` |
| Escolher a conta com mais espaço | Não | `ccam next --auto` |
| Nome da conta na linha de status | Não | `ccam statusline install` |
| Corrigir links quebrados | Manualmente | `ccam doctor --fix` |
| Software adicional | Nenhum | Um único binário estático |

Escolha aliases se você tem duas contas e não quer compartilhar nada entre elas. Escolha o ccam se você quer que as contas pareçam uma única configuração. Você pode começar com aliases e passar para o ccam depois: `ccam import` registra as pastas `~/.claude-*` existentes no lugar onde estão.

[Como usar várias contas no mesmo computador](/guides/multiple-claude-code-accounts) explica as duas abordagens passo a passo, e [o ccam e outras ferramentas para trocar de conta](/compare/account-switchers) compara o ccam com ferramentas que trocam um único login.
