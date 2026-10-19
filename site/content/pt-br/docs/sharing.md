---
title: Compartilhar servidores MCP, configuração e histórico
description: Como o ccam compartilha servidores MCP, settings.json, CLAUDE.md, agentes, skills e histórico de conversas entre contas do Claude Code, sem compartilhar logins.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 1e2f0d6809fc
---

Cada conta mantém o próprio login. Todo o resto pode ser compartilhado item por item: o ccam cria um link do arquivo ou da pasta da conta para uma pasta de grupo compartilhada, então uma alteração feita em uma conta aparece em todas as contas desse grupo.

## O que pode ser compartilhado

| Item | Opção | O que é compartilhado |
| --- | --- | --- |
| Servidores MCP | `--mcp` | O `mcp.json` do grupo, passado ao Claude Code na inicialização |
| Configurações | `--settings` | `settings.json`: permissões, hooks, modelo, linha de status |
| Memória global | `--memory` | `CLAUDE.md` |
| Extensões | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Histórico de conversas | `--history` | `projects/` e `history.jsonl` |

Contas novas compartilham tudo, exceto o histórico de conversas.

Os arquivos de login `.credentials.json` e `.claude.json` nunca são compartilhados. O ccam só cria links para caminhos de uma lista fixa, e esses dois não estão nela.

## Grupos

`on` significa o grupo `default`. Você também pode informar um nome, para que apenas algumas contas compartilhem um item:

```sh
ccam share work --history team
ccam share side --history team
```

Agora `work` e `side` compartilham conversas entre si, mas não com suas outras contas.

## Ativar o compartilhamento

Quando uma conta começa a compartilhar um item, o ccam:

1. imprime um plano com cada arquivo que vai mover, copiar ou vincular, e pede confirmação antes de começar;
2. move os arquivos da conta para a pasta do grupo se o grupo ainda estiver vazio, ou adiciona os arquivos que o grupo ainda não tem;
3. quando o mesmo arquivo é diferente, pergunta se você quer manter a versão do grupo, usar a sua ou mesclar o JSON;
4. mantém o seu original em `.ccam-backup/<time>/` dentro da pasta da conta;
5. substitui o original por um link para o grupo.

Os arquivos de conversa recebem o nome do ID da sessão, então os históricos de várias contas se combinam sem conflitos.

## Desativar o compartilhamento

```sh
ccam share work --history off
```

O link é removido e a conta recebe a própria cópia do conteúdo atual do grupo. Adicione `--fresh` para começar do zero. O grupo mantém o conteúdo para as outras contas.

## Copiar uma configuração existente

Se você usava o Claude Code antes do ccam, suas configurações estão em `~/.claude`. O `ccam migrate` as copia para um grupo, para que todas as contas que compartilham esses itens as recebam:

```sh
ccam migrate
```

O que é copiado:

- o `settings.json` e o `CLAUDE.md` global;
- as extensões `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` e `rules/`;
- seus servidores MCP de usuário, ou seja, a entrada `mcpServers` do `~/.claude.json`, para o `mcp.json` do grupo;
- o histórico de conversas (`projects/` e `history.jsonl`), mas só quando você adiciona `--history`.

A pasta de origem é apenas lida: nada nela é movido, vinculado ou alterado, e o `claude` simples continua funcionando como antes. Os arquivos de login `.credentials.json` e `.claude.json` nunca são copiados.

O ccam imprime o plano e pergunta antes de copiar qualquer coisa:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

No final, o ccam lista as contas que compartilham cada item por meio do grupo. Se algumas contas não compartilham, ele imprime os comandos `ccam share` que ativam o compartilhamento para elas.

Quando o grupo já tem uma cópia diferente de um arquivo, o ccam pergunta, como faz ao ativar o compartilhamento, se você quer manter a cópia do grupo ou usar a sua. O `settings.json` e outros arquivos JSON também podem ser mesclados, e o `history.jsonl` linha por linha. As cópias anteriores do grupo são salvas em `.ccam-backup/<time>/` dentro da pasta do grupo.

| Opção | Efeito |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copia só os itens que você indicar; sem elas, tudo exceto o histórico de conversas |
| `-g`, `--group <group>` | Copia para este grupo em vez de `default` |
| `-y`, `--yes` | Copia sem perguntar; quando um arquivo é diferente, mantém a cópia do grupo |

Coloque `--dry-run` antes do comando, como em `ccam --dry-run migrate`, para ver o plano sem copiar nada. Para copiar de outra pasta, indique-a: `ccam migrate ~/.claude-work --group work`.

## Servidores MCP

O Claude Code armazena os servidores MCP no `.claude.json`, junto com o login, então esse arquivo não pode ser compartilhado. O ccam mantém os servidores compartilhados no `mcp.json` do grupo e inicia o Claude Code com `--mcp-config=<group>/mcp.json`. Suas próprias opções `--mcp-config` continuam valendo por cima.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Um servidor que você adiciona dentro de uma sessão com `claude mcp add` é salvo apenas para essa conta. Na próxima vez que você iniciar essa conta, o ccam informa o nome do servidor e o comando que o move para o grupo:

```sh
ccam mcp import work gh
```

Servidores MCP que fazem login com OAuth armazenam os tokens junto com o login da conta, então autorize-os uma vez em cada conta.

Se você prefere que os servidores apareçam em `claude mcp list` e em `/mcp` como servidores comuns de usuário, defina `mcp_strategy = "sync"` em `~/.config/ccam/config.toml`. Assim, o ccam grava os servidores no `.claude.json` de cada conta antes da inicialização, usando o bloqueio de arquivo do próprio Claude Code e com um backup.

## Duas contas na mesma pasta

Contas que compartilham o histórico de conversas gravam nos mesmos arquivos. Se você iniciar um segundo Claude Code em uma pasta onde outra conta do mesmo grupo já está em execução, o ccam avisa e pergunta se você quer continuar.

## Quando algo dá errado

```sh
ccam doctor
ccam doctor --fix
```

Se um caminho compartilhado virou um arquivo comum, por exemplo depois que você restaurou um backup manualmente, o ccam para antes de iniciar o Claude Code e indica `ccam doctor --fix`, que mescla o arquivo no grupo e mantém um backup.
