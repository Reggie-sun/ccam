---
title: Compartilhar servidores MCP, configuração e histórico
description: Como os grupos do ccam compartilham servidores MCP, settings.json, CLAUDE.md, skills e histórico de conversas entre contas do Claude Code, sem misturar logins.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "27e0a0e5abbd"
---

Cada conta mantém o próprio login. Todo o resto pode ser compartilhado por meio de **grupos**: um grupo indica os itens que compartilha e as contas que fazem parte dele, e o ccam cria links dos arquivos de cada uma dessas contas para a pasta do grupo, então uma alteração feita em uma conta aparece em todas elas.

## O que pode ser compartilhado

| Item | O que é compartilhado |
| --- | --- |
| `mcp` | Servidores MCP: o `mcp.json` do grupo, passado ao Claude Code na inicialização |
| `settings` | `settings.json`: permissões, hooks, modelo, linha de status |
| `memory` | o `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | histórico de conversas: `projects/` e `history.jsonl` |

A sua primeira conta cria o grupo `default`, que compartilha tudo exceto o histórico de conversas, e as contas seguintes entram nele.

Os arquivos de login `.credentials.json` e `.claude.json` nunca são compartilhados. O ccam só cria links para caminhos de uma lista fixa, e esses dois não estão nela.

Esta página cobre as tarefas do dia a dia. Todas as opções estão descritas nas páginas de recursos: [compartilhamento com grupos](/features/share), [servidores MCP compartilhados](/features/mcp) e [cópia de uma configuração existente](/features/migrate).

## Grupos

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Uma conta pode estar em vários grupos, mas cada item vem de um único grupo. Então, para compartilhar o histórico entre `work` e `side`, mas não com as suas outras contas, deixe o histórico fora do `default` e dê às duas um grupo só delas, como acima. `ccam share work` mostra o que uma conta compartilha e por meio de qual grupo. O [painel](/features/dashboard) faz o mesmo na aba Groups.

## Quando uma conta começa a compartilhar

Quando uma conta começa a compartilhar um item (ela entra em um grupo, ou o grupo dela passa a compartilhar o item), o ccam:

1. imprime um plano com cada arquivo que vai mover, copiar ou vincular para cada conta, e pede confirmação antes de começar;
2. move os arquivos da conta para a pasta do grupo se o grupo ainda estiver vazio, ou adiciona os arquivos que o grupo ainda não tem;
3. quando o mesmo arquivo é diferente, pergunta se você quer manter a versão do grupo, usar a sua ou mesclar o JSON;
4. mantém o seu original em `.ccam-backup/<time>/` dentro da pasta da conta;
5. substitui o original por um link para o grupo.

Os arquivos de conversa recebem o nome do ID da sessão, então os históricos de várias contas se combinam sem conflitos.

O ccam não altera uma conta enquanto o Claude Code estiver em execução com ela, seja qual for a forma como ele foi iniciado (inclusive sessões em segundo plano), porque a sessão continua gravando nos arquivos que estão sendo movidos. Encerre-a primeiro. O ccam confere de novo logo antes de salvar, então uma sessão iniciada enquanto o plano esperava também é detectada.

## Quando uma conta deixa de compartilhar

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

O link é removido e a conta recebe a própria cópia do conteúdo atual do grupo; os servidores MCP do grupo vão para o próprio `.claude.json` dela, ao lado dos servidores que já são dela. Adicione `--fresh` para começar do zero; o plano então lista os servidores MCP que deixam de carregar. Uma pasta importada como `~/.claude` sempre fica com uma cópia. O grupo mantém o conteúdo para as outras contas.

## Copiar uma configuração existente

Se você usava o Claude Code antes do ccam, suas configurações estão em `~/.claude`. O `ccam migrate` as copia para um grupo, para que todas as contas que compartilham esses itens as recebam:

```sh
ccam migrate
```

Em um terminal, o ccam mostra primeiro uma lista de seleção dos itens. Tudo vem marcado, exceto o histórico de conversas; pressione Espaço para alterar um item e Enter para ver o plano. Itens para os quais a pasta não tem nada aparecem como `not found`.

O que é copiado:

- o `settings.json` e o `CLAUDE.md` global;
- as extensões `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` e `rules/`;
- seus servidores MCP de usuário, ou seja, a entrada `mcpServers` do `~/.claude.json`, para o `mcp.json` do grupo;
- o histórico de conversas (`projects/` e `history.jsonl`), mas só quando você adiciona `--history`.

A pasta de origem é apenas lida: nada nela é movido, vinculado ou alterado, e o `claude` simples continua funcionando como antes. Os arquivos de login `.credentials.json` e `.claude.json` nunca são copiados para o grupo.

Servidores MCP que fazem login com OAuth guardam os tokens no login, não na lista de servidores. O ccam não copia esses logins: um login copiado tem o mesmo refresh token, e a primeira renovação desconecta daquele servidor todos os outros que o têm, inclusive o `~/.claude`. O plano e o resumo indicam as contas que precisam entrar nos servidores por conta própria com `/mcp`.

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
Group default: personal
To share them:
  ccam group join default work
```

No final, o ccam lista as contas do grupo. Se o grupo ainda não compartilha um dos itens, ou se algumas contas não estão nele, ele imprime os comandos `ccam group` que mudam isso.

Quando o grupo já tem uma cópia diferente de um arquivo, o ccam pergunta, como faz ao ativar o compartilhamento, se você quer manter a cópia do grupo ou usar a sua. O `settings.json` e outros arquivos JSON também podem ser mesclados, e o `history.jsonl` linha por linha. As cópias anteriores do grupo são salvas em `.ccam-backup/<time>/` dentro da pasta do grupo.

| Opção | Efeito |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copia só os itens que você indicar, sem a lista de seleção; sem elas e fora de um terminal, tudo exceto o histórico de conversas |
| `-g`, `--group <group>` | Copia para este grupo em vez de `default` |
| `-y`, `--yes` | Copia sem perguntar; quando um arquivo é diferente, mantém a cópia do grupo |

Coloque `--dry-run` antes do comando, como em `ccam --dry-run migrate`, para ver o plano sem copiar nada. Para copiar de outra pasta, indique-a: `ccam migrate ~/.claude-work --group work`.

## Servidores MCP

O Claude Code armazena os servidores MCP no `.claude.json`, junto com o login, então esse arquivo não pode ser compartilhado. O ccam mantém os servidores compartilhados no `mcp.json` do grupo e inicia o Claude Code com `--mcp-config=<group>/mcp.json`. Suas próprias opções `--mcp-config` continuam valendo por cima.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Um servidor que você adiciona dentro de uma sessão com `claude mcp add` é salvo apenas para essa conta e para a pasta em que você estava (o escopo local do Claude Code; `-s user` o torna da conta em todas as pastas). Na próxima vez que você iniciar essa conta, o ccam informa o nome do servidor, com a pasta, e o comando que o move para o grupo:

```sh
ccam mcp import work gh
```

Se o grupo já tiver um servidor diferente com esse nome, o ccam o ignora e avisa. `--force` substitui o servidor do grupo e guarda o `mcp.json` antigo na pasta `.ccam-backup/` do grupo.

Servidores MCP que fazem login com OAuth armazenam os tokens junto com o login da conta, então autorize-os uma vez em cada conta com `/mcp`. `ccam ext` e `ccam doctor` indicam os servidores em que uma conta não fez login enquanto outra conta fez.

Se você prefere que os servidores apareçam em `claude mcp list` e em `/mcp` como servidores comuns de usuário, defina `mcp_strategy = "sync"` no início do `~/.config/ccam/config.toml`, acima do primeiro `[[group]]`, `[[profile]]` ou `[[folder]]`. Assim, o ccam grava os servidores no `.claude.json` de cada conta antes da inicialização, usando o bloqueio de arquivo do próprio Claude Code e com um backup. Um servidor que a conta tem por conta própria com o mesmo nome continua sendo da conta, e o ccam avisa isso.

## Duas contas na mesma pasta

Contas que compartilham o histórico de conversas gravam nos mesmos arquivos. Se você iniciar um segundo Claude Code em uma pasta onde outra conta do mesmo grupo já está em execução, o ccam avisa e pergunta se você quer continuar.

## Quando algo dá errado

```sh
ccam doctor
ccam doctor --fix
```

Se um caminho compartilhado virou um arquivo comum, por exemplo depois que você restaurou um backup manualmente, o ccam para antes de iniciar o Claude Code e indica `ccam doctor --fix`, que o liga ao grupo de novo: o arquivo vira a cópia do grupo quando o grupo não tem uma e, caso contrário, o grupo mantém a sua e o seu vai para um backup, que o `--fix` indica. Uma troca de grupos que parou no meio também é concluída pelo `ccam doctor --fix`.
