---
title: Como usar duas contas do Claude Code no mesmo computador
description: Execute várias contas do Claude Code lado a lado com CLAUDE_CONFIG_DIR, manualmente com aliases de shell ou automaticamente com o ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6a4b7d961e10
---

Dê a cada conta a própria pasta de configuração. O Claude Code lê a pasta da variável de ambiente `CLAUDE_CONFIG_DIR`, e cada pasta guarda o próprio login, as próprias configurações e o próprio histórico. Você pode configurar isso manualmente com um alias de shell ou deixar o ccam fazer isso.

## Manualmente, com um alias

Esta é a abordagem descrita na documentação do Claude Code. Adicione uma linha para cada conta extra ao `~/.bashrc` ou ao `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Abra um novo terminal e execute `claude-work`. O Claude Code pede que você faça login, porque a pasta é nova. A partir daí, `claude` é sua primeira conta e `claude-work` a segunda.

Isso funciona bem para duas contas. O que você precisa saber:

- Cada pasta começa vazia: sem configurações, sem `CLAUDE.md`, sem servidores MCP, sem agentes nem skills. Você copia ou cria links para o que precisa e mantém as cópias em sincronia por conta própria.
- Aliases existem apenas em shells interativos. Scripts, Makefiles e outros programas não os enxergam.
- Os servidores MCP ficam no `.claude.json`, ao lado do login, então não é possível simplesmente criar links para eles entre pastas.
- As conversas ficam na conta em que começaram, então você não pode continuar uma delas em outra conta.
- Não há nada que mostre o uso de todas as contas de uma vez.

## Com o ccam

O ccam automatiza a mesma ideia e acrescenta as partes que faltam:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Agora `ccam personal` e `ccam work` iniciam as duas contas, assim como `ccwork`, o atalho que você deu à conta work. Como são programas de verdade e não aliases, eles também funcionam em scripts. Configurações, `CLAUDE.md`, servidores MCP e extensões são compartilhados desde o início, e você pode ativar o histórico de conversas compartilhado para continuar uma conversa na outra conta com `ccam next --continue`.

Se você já usa `~/.claude`, execute `ccam import` em vez de adicionar essa conta de novo, ou `ccam migrate` para copiar as configurações, os servidores MCP e as extensões dela para as configurações compartilhadas.

## Qual escolher

Para duas contas sem nada compartilhado, um alias basta. Quando você quer configurações e servidores MCP compartilhados, mais contas, o uso visível de relance ou continuar conversas, o ccam poupa o trabalho de manter tudo em ordem. A [comparação](/compare/aliases) lista as diferenças lado a lado.
