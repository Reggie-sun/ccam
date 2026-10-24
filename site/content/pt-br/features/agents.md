---
title: Sessões em segundo plano de todas as contas
description: ccam agents lista em uma só tabela os agentes em segundo plano e as sessões abertas do Claude Code de todas as suas contas, e ccam attach abre uma pelo id.
updated: 2026-10-23
order: 16
translation: machine
sourceHash: "092368f63363"
---

O Claude Code mantém as sessões em segundo plano (`claude --bg`, visão de agentes) por conta: `claude agents` só mostra as sessões da conta com que ele está rodando. `ccam agents` consulta todas as contas ao mesmo tempo e lista tudo junto, e `ccam attach <id|name>` abre uma sessão em segundo plano como a conta que a iniciou, seja ela qual for.

## Uso

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` é o estado da sessão em segundo plano (`working`, `blocked`, `done`, `failed`, `stopped`) e o que o processo dela está fazendo (`busy`, `waiting`, `idle`); uma sessão interativa aberta aparece como `open`.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--all` | Lista também as sessões em segundo plano concluídas, com falha ou interrompidas | desativado |
| `--json` | Imprime as sessões em JSON, cada uma com a `account` a que pertence | desativado |

`ccam attach <id|name> [claude args]` repassa ao `claude attach` tudo o que vem depois do id sem alterar, inclusive opções: `ccam attach 49ec --dry-run` entrega `--dry-run` ao Claude Code. Para imprimir o comando em vez disso, coloque `--dry-run` antes do id (`ccam --dry-run attach 49ec` ou `ccam attach --dry-run 49ec`).

## Como funciona

1. Para cada conta, o ccam executa `claude agents --json` (com `--all`, quando informado) com o `CLAUDE_CONFIG_DIR` da conta, todas ao mesmo tempo, cada uma com limite de 15 segundos. Isso lê a lista de sessões da própria conta e não inicia o serviço em segundo plano do Claude Code.
2. Cada entrada recebe o nome da conta, e a lista é ordenada pela hora de início, da mais recente para a mais antiga. Uma conta que não responde é indicada na saída de erro padrão e fica de fora, e então o ccam sai com código 1, para que um script não confunda uma falha com uma lista vazia.
3. `ccam attach` procura `jobs/<id>` na pasta de cada conta, onde o Claude Code guarda uma sessão em segundo plano, e executa `claude attach <id>` como essa conta, do mesmo jeito que `ccam <name>` inicia o Claude Code. Uma palavra que não é o começo de um id é procurada nos nomes das sessões, sem diferenciar maiúsculas e minúsculas; uma palavra que parece um id nunca é tomada por um nome.

## Segurança e casos especiais

- Os ids das sessões em segundo plano pertencem a uma conta: `claude attach` executado como outra conta não os encontra, e é isso que o `ccam attach` resolve.
- Duas contas podem ter sessões cujos ids começam igual, ou cujos nomes têm a mesma palavra; nesse caso, o `ccam attach` mostra essas sessões e pede uma parte maior do id ou do nome.
- O ccam passa seus servidores MCP compartilhados ao Claude Code com `--mcp-config=<file>` em um único argumento, que `claude agents` e `claude attach` aceitam; as sessões em segundo plano continuam usando esse arquivo quando reiniciam.
- O serviço em segundo plano de uma conta com chave de API pega a chave do shell que o iniciou. Inicie a primeira sessão em segundo plano dessa conta pelo ccam (`ccam gw --bg …`), para que o serviço tenha a chave.

## Relacionados

- [Iniciar o Claude Code](/features/launch)
- [Escolher uma conta para scripts](/features/pick)
- [O painel](/features/dashboard)
