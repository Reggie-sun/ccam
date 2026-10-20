---
title: Compartilhamento com grupos
description: Um grupo do ccam compartilha itens escolhidos entre as suas contas por uma pasta; o ccam group move os dados com plano, backups e escolhas nos conflitos.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

Cada conta tem a própria pasta, mas a maior parte do que há nela não precisa ficar separada. Um **grupo** diz o que é compartilhado e por quem: cada conta do grupo cria links desses itens para uma pasta, `~/.local/share/ccam/groups/<group>/`, então todas leem e gravam a mesma cópia. Os tokens de login nunca são compartilhados.

Uma conta pode estar em vários grupos, mas cada item vem de um único grupo. Uma configuração típica é um grupo para as configurações e as extensões e um segundo para o histórico de conversas entre duas contas:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

A primeira conta que você adiciona cria o `default`, e as contas seguintes entram nele, a menos que você diga o contrário.

## Uso

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Todo comando que move dados imprime um plano, pergunta sobre cada arquivo que difere da cópia do grupo e pergunta `Proceed? [y/N]` antes de alterar qualquer coisa:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

As mesmas alterações estão na aba **Groups** do [painel](/features/dashboard), e `s` em uma conta escolhe os grupos dela.

## Itens

| Item | O que abrange |
| --- | --- |
| `mcp` | servidores MCP, pelo `mcp.json` do grupo ([Servidores MCP compartilhados](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | o `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` e `history.jsonl`, para que outra conta possa usar `--continue` |

Um item que nenhum grupo da conta compartilha fica na própria pasta da conta.

## Comandos e opções

| Comando | O que faz |
| --- | --- |
| `ccam group` ou `ccam group list` | lista os grupos |
| `ccam group show <group>` | o que o grupo compartilha, as contas, a pasta e os servidores MCP dele |
| `ccam group add <group> --items <items> [--members <accounts>]` | cria um grupo; os itens das contas dele são movidos para ele |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | altera o que o grupo compartilha, para todas as contas dele |
| `ccam group join <group> <account>...` | coloca contas no grupo |
| `ccam group leave <group> <account>...` | tira contas do grupo; cada uma fica com uma cópia |
| `ccam group rename <group> <new-name>` | renomeia apenas na configuração |
| `ccam group remove <group> [--leave]` | remove um grupo |

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: a conta cuja cópia entra primeiro no grupo | a conta padrão, depois por número |
| `--replace` | `join`: sai antes dos grupos que compartilham os mesmos itens | desativado |
| `--fresh` | as contas que deixam de compartilhar um item começam sem ele, em vez de com uma cópia | desativado |
| `--leave` | `remove`: tira antes as contas do grupo | necessário quando ele tem contas |
| `--prefer <group\|mine\|merge>` | como decidir os conflitos | pergunta; com `--yes`, o padrão de cada arquivo |
| `-y`, `--yes` | executa sem perguntar; os conflitos seguem `--prefer` | desativado |
| `--dry-run` (global) | imprime o plano e não altera nada | desativado |

As contas são indicadas pelo nome ou pelo número; os itens, como uma lista separada por vírgulas. O `ccam add` aceita `--group <group>` (pode ser repetido) ou `--no-group`.

## Como funciona

Os arquivos envolvidos:

- `~/.config/ccam/config.toml`: as entradas `[[group]]`: `name`, `items`, `members` e, quando o nome da pasta é diferente, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: o `CLAUDE_CONFIG_DIR` da conta.
- `~/.local/share/ccam/groups/<folder>/`: as cópias compartilhadas.

Quando o `config.toml` é lido, os grupos de cada conta são resolvidos em uma pasta por item. Isso leva cerca de 25 µs para cinco contas, então uma inicialização continua lendo apenas o `config.toml`. Uma configuração em que uma conta receberia um item de dois grupos é um erro, e nada inicia até que ela seja corrigida.

Uma alteração trabalha sobre uma cópia da configuração:

1. O comando edita a cópia (`AddGroup`, `SetGroupItems`, `JoinGroup` …), que verifica as regras.
2. `internal/group` compara o antigo e o novo e, para cada conta cujos itens mudam, monta um plano com `share.Build`. Os planos são montados um depois do outro, e cada um vê o que os anteriores vão colocar nas pastas dos grupos, então só a primeira conta move a sua cópia para lá. A primeira cópia vem da pasta do grupo, se ela já tiver uma, depois de `--from`, depois da conta padrão e depois das outras, por número.
3. O ccam se recusa se uma dessas contas estiver em execução e, no Windows sem o Modo de Desenvolvedor, se um arquivo precisar de um link.
4. Os conflitos são decididos, pela pergunta ou por `--prefer`, o plano é impresso, e o ccam pergunta.
5. A nova configuração é salva antes que qualquer coisa seja movida. Se uma etapa parar no meio, `ccam doctor --fix` e uma nova execução montam o plano em direção aos grupos salvos e concluem a alteração.
6. Os planos são executados em ordem. Cada etapa move dados para o grupo ou para a pasta de backup; nada é excluído.

`share.Layout` é uma lista de permissões dos caminhos que cada item abrange. Qualquer coisa fora dela nunca chega a um grupo. `share.Never` indica `.credentials.json` e `.claude.json`, e `Build` para com um erro interno se algum deles aparecer no layout. O item `mcp` não tem caminhos: a lista de servidores fica dentro do `.claude.json`, então o [compartilhamento de MCP](/features/mcp) funciona pelo `mcp.json` do grupo.

Isto é o que acontece com uma pasta, `projects/`, quando `side` começa a compartilhar o histórico e o grupo já tem conteúdo:

<!-- motion:share-on -->

Se o grupo ainda não tem nada nesse caminho, a conta é a primeira a compartilhá-lo, e a pasta dela é movida para lá como está. Caso contrário, cada arquivo que o grupo não tem é movido para o grupo, e os arquivos que existem dos dois lados ficam na pasta da conta, então a versão do grupo é mantida. Depois, a pasta da conta vai para `.ccam-backup/<time>/` dentro da pasta da conta, e um link toma o lugar dela. Os arquivos de conversa recebem o nome do ID da sessão, então os históricos de várias contas se combinam sem conflitos. Os registros de `plugins/` armazenam caminhos absolutos que passam pela pasta da conta; o ccam os aponta para a pasta do grupo e mescla as entradas da conta em um registro que o grupo já tenha.

Sair inverte o link:

<!-- motion:share-off -->

O ccam só remove o link se ele apontar para dentro da pasta de grupos e, em seguida, copia o conteúdo atual do grupo para a conta. Com `--fresh`, a cópia é pulada. O grupo mantém o conteúdo para as outras contas. Uma conta que move um item para um grupo cuja pasta ainda está vazia leva uma cópia do que via, então mantém os seus arquivos; os servidores MCP vêm junto da mesma forma. Uma conta que deixa de compartilhar MCP perde os servidores que o ccam tinha copiado para o `.claude.json` dela (um backup é mantido), já que eles vieram do grupo.

Um único arquivo que existe dos dois lados com conteúdo diferente é um conflito:

<!-- motion:share-conflict -->

`history.jsonl` oferece `merge by line`, `keep group` e `use mine`, com a mesclagem primeiro. Outros arquivos `.json` oferecem `keep group`, `use mine` e `merge JSON`. Qualquer outro arquivo oferece `keep group` e `use mine`. `merge JSON` mescla objetos recursivamente, junta arrays sem duplicatas e mantém o valor do grupo quando os dois lados definem a mesma chave de forma diferente, para que as outras contas não vejam mudança. Toda escolha, exceto `keep group`, salva antes a cópia do grupo em `.ccam-backup/<time>/group/` da conta. O seu próprio arquivo sempre vai para o backup.

A cada inicialização, `share.Check` compara cada caminho do layout com os grupos da conta usando apenas `lstat` e `readlink`. Um link ausente é criado na hora. Um arquivo comum onde deveria haver um link, um link para o lugar errado ou um link obsoleto interrompe a inicialização e indica `ccam doctor --fix`.

### Renomear, remover e o formato antigo

- `ccam group rename` altera apenas a configuração. A pasta mantém o nome, e o grupo o registra como `dir`, então nada é movido e nada pode se perder.
- `ccam group remove` se recusa a remover um grupo com contas, a menos que você adicione `--leave`. Cada conta fica então com uma cópia, e a pasta vai para `~/.local/share/ccam/trash/`, a menos que outro grupo a use.
- Antes dos grupos, o `config.toml` guardava uma tabela `[profile.share]` por conta. O ccam lê esse formato como grupos sem mover nada. No primeiro salvamento, ele grava o formato novo e mantém o arquivo antigo como `config.toml.bak`. Os itens compartilhados pelas mesmas contas viram um grupo com o nome da pasta; os demais recebem nomes como `default-history` e continuam usando a mesma pasta.

### No Windows

O Windows só cria links simbólicos no Modo de Desenvolvedor (Configurações > Sistema > Para desenvolvedores) ou para um administrador. Sem essa permissão, uma pasta compartilhada vira uma junção de diretório, que funciona do mesmo jeito para o Claude Code. Um arquivo compartilhado (`settings.json`, `CLAUDE.md`, `history.jsonl`) só pode ser um link simbólico, então uma alteração de grupo que precise de um é recusada antes que qualquer coisa seja movida.

## Segurança e casos especiais

- **Contas em execução.** Uma alteração é recusada enquanto o Claude Code estiver em execução com uma das contas que ela move, porque a sessão continua gravando nesses arquivos. Sair e entrar em um grupo dizem respeito apenas à conta indicada. `--dry-run` ainda imprime o plano.
- **Backups.** Cada plano recebe a própria pasta `.ccam-backup/<time>/` na pasta da conta. O ccam nunca as exclui.
- **Seus próprios links.** Se um caminho já é um link que você criou, por exemplo para um repositório de dotfiles, e o grupo está vazio, o ccam copia o destino do link para o grupo e move o próprio link para o backup. O destino não é alterado.
- **Pastas importadas.** Uma conta registrada no lugar com `ccam import` (como `~/.claude`) é alterada no lugar, então o `claude` simples também vê a mudança; o plano avisa isso.
- **Duas contas na mesma pasta.** Contas que compartilham o histórico gravam nos mesmos arquivos. Quando uma segunda conta inicia em uma pasta onde outra já está em execução, o ccam avisa e pergunta.
- **Edições manuais.** Um membro que não é uma conta é ignorado e informado pelo `ccam doctor`. Os nomes `local`, `on`, `off` e `none` não podem ser nomes de grupos.

## Relacionados

- [O painel](/features/dashboard)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Servidores MCP compartilhados](/features/mcp)
- [Copiar uma configuração existente](/features/migrate)
- [Todos os comandos](/docs/commands)
