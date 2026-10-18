---
title: Apresentando o ccam
description: Por que o ccam existe, o que ele faz por quem tem várias contas do Claude Code e como ele fica fora do seu caminho.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 5e4121938d8c
---

ccam é uma pequena ferramenta de linha de comando para quem usa mais de uma conta do Claude Code no mesmo computador: uma conta pessoal e uma de trabalho, ou várias contas para um projeto longo.

## O problema

O Claude Code oferece suporte a várias contas por meio do `CLAUDE_CONFIG_DIR`: uma pasta por conta. Na prática, cada pasta é uma ilha. Configurações, `CLAUDE.md`, servidores MCP, agentes e skills precisam ser configurados de novo em cada pasta e mantidos em sincronia manualmente. As conversas ficam onde começaram. E quando uma conta atinge o limite de uso, você só descobre tentando.

## O que o ccam faz

- `ccam add` cria a pasta, executa o login e dá um atalho à conta: `claude1`, `claude2` e assim por diante.
- `ccs work` e `claude2` iniciam o Claude Code como essa conta e repassam todos os argumentos sem alterações.
- Servidores MCP, configurações, `CLAUDE.md` e extensões são compartilhados por meio de links, item por item e conta por conta. Logins nunca são compartilhados.
- `ccs usage` mostra o uso de 5 horas e o uso semanal de todas as contas; `ccs next --continue` retoma a conversa na próxima conta.

## Fora do seu caminho

Uma ferramenta que você executa antes de cada sessão não deve ter um custo perceptível. O ccam lê um arquivo de configuração, verifica cada link compartilhado com uma única chamada de sistema e depois se substitui pelo `claude`. Ele acrescenta cerca de 4 milissegundos e não fica em execução enquanto você trabalha.

Ele também é cuidadoso com seus arquivos. Toda alteração no compartilhamento imprime um plano antes de começar e mantém os originais em uma pasta de backup. `ccam doctor` verifica links, permissões e logins, e corrige o que puder corrigir com segurança.

## Experimente

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

O ccam é de código aberto e não oficial: ele não é afiliado à Anthropic nem endossado por ela. Se você pode usar várias contas depende dos termos de cada conta. Feedback e traduções são bem-vindos no [GitHub](https://github.com/jhyunwoo/ccam).
