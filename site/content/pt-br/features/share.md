---
title: Compartilhamento com grupos
description: Um grupo do ccam compartilha itens escolhidos entre as suas contas por uma pasta; o ccam group move os dados com plano, backups e escolhas nos conflitos.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "e4f05b0fd572"
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
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
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
| `ccam group leave <group> <account>...` | tira contas do grupo; cada uma fica com uma cópia, inclusive dos servidores MCP |
| `ccam group rename <group> <new-name>` | renomeia apenas na configuração |
| `ccam group remove <group> [--leave]` | remove um grupo |

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: a conta cuja cópia entra primeiro no grupo | a conta padrão, depois por número |
| `--replace` | `join`: sai antes dos grupos que compartilham os mesmos itens | desativado |
| `--fresh` | as contas que deixam de compartilhar um item começam sem ele, em vez de com uma cópia (uma pasta importada ainda fica com uma) | desativado |
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
3. O ccam se recusa se uma dessas contas estiver em execução, seja qual for a forma como o Claude Code foi iniciado (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), e, no Windows sem o Modo de Desenvolvedor, se um arquivo precisar de um link.
4. Os conflitos são decididos, pela pergunta ou por `--prefer`, o plano é impresso, e o ccam pergunta.
5. Logo antes de salvar, o ccam verifica de novo: uma conta que foi iniciada enquanto as perguntas esperavam interrompe a alteração, assim como a alteração dos mesmos grupos feita por outro ccam (`the groups changed while this plan waited; nothing was moved, make the change again`). A nova configuração é salva antes que qualquer coisa seja movida, e a alteração é anotada em `~/.local/state/ccam/group-change.json`. Se uma etapa parar no meio, a anotação fica: `ccam doctor --fix` conclui a alteração na ordem dela, com as respostas que você deu, inclusive as etapas de MCP, e até lá nenhuma outra troca de grupos começa (``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``).
6. As etapas de MCP rodam primeiro. Quando contas tiram as extensões de um grupo, as listas de plugins do grupo são apontadas para a pasta do grupo onde quer que os caminhos delas passem pela pasta de uma conta (`group plugin lists:` no plano). Depois, os planos são executados em ordem. Cada etapa move dados para o grupo ou para a pasta de backup; só os arquivos que o grupo já tem idênticos, byte a byte, são descartados em vez de ir para o backup.

`share.Layout` é uma lista de permissões dos caminhos que cada item abrange. Qualquer coisa fora dela nunca chega a um grupo. `share.Never` indica `.credentials.json` e `.claude.json`, e `Build` para com um erro interno se algum deles aparecer no layout. O item `mcp` não tem caminhos: a lista de servidores fica dentro do `.claude.json`, então o [compartilhamento de MCP](/features/mcp) funciona pelo `mcp.json` do grupo.

Isto é o que acontece com uma pasta, `projects/`, quando `side` começa a compartilhar o histórico e o grupo já tem conteúdo:

<!-- motion:share-on -->

Se o grupo ainda não tem nada nesse caminho, a conta é a primeira a compartilhá-lo, e a pasta dela é movida para lá como está. Caso contrário, cada arquivo que o grupo não tem é movido para o grupo, e os arquivos que existem dos dois lados ficam na pasta da conta, então a versão do grupo é mantida. Os arquivos idênticos aos do grupo são descartados, e o plano indica os que diferem (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`). Depois, a pasta da conta vai para `.ccam-backup/<time>/` dentro da pasta da conta, e um link toma o lugar dela. Se a entrada do grupo for ela mesma um link para a entrada da própria conta que está entrando, o conteúdo da conta vai para o grupo no lugar dela. Os arquivos de conversa recebem o nome do ID da sessão, então os históricos de várias contas se combinam sem conflitos. Os registros de `plugins/` armazenam caminhos absolutos que passam pela pasta da conta; o ccam os aponta para a pasta do grupo e mescla as entradas da conta em um registro que o grupo já tenha.

Sair inverte o link:

<!-- motion:share-off -->

Primeiro, o ccam copia o conteúdo atual do grupo para `.ccam-backup/incomplete-<name>` na pasta da conta e só substitui o link quando a cópia está completa, então uma cópia interrompida deixa o link no lugar e `ccam doctor --fix` copia de novo (o plano dessa conclusão mostra `restore <entry>: … (copied back from the group before a stop)`). Ele só remove o link se ele ainda apontar para dentro da pasta de grupos como apontava quando o plano foi feito; um caminho que mudou desde então, como um arquivo comum onde estava o link, fica como está, e a alteração para. Com `--fresh`, a cópia é pulada. O grupo mantém o conteúdo para as outras contas. Uma conta que move um item para um grupo cuja pasta ainda está vazia leva uma cópia do que via, então mantém os seus arquivos; essa cópia é preparada em `<group>/.ccam-incomplete-<name>` até estar completa. Uma conta que deixa de compartilhar MCP mantém os servidores do grupo: eles vão para o próprio `.claude.json` dela (um backup é mantido), ao lado dos servidores próprios dela com os mesmos nomes, que permanecem. Com `--fresh`, o plano lista os servidores que deixam de carregar em vez disso. Entrar de novo remove as cópias que ninguém alterou. Sair de `extensions` também aponta as listas de plugins do grupo para a pasta do grupo onde os caminhos delas passavam pela pasta da conta que sai; para as outras contas, elas significam o mesmo.

Um único arquivo que existe dos dois lados com conteúdo diferente é um conflito:

<!-- motion:share-conflict -->

Você escolhe para arquivos avulsos como `settings.json`; dentro de pastas, a cópia do grupo é mantida, e o plano indica os arquivos que diferem. `history.jsonl` oferece `merge by line`, `keep group` e `use mine`, com a mesclagem primeiro. Outros arquivos `.json` oferecem `keep group`, `use mine` e `merge JSON`. Qualquer outro arquivo oferece `keep group` e `use mine`. `merge JSON` mescla objetos recursivamente, junta arrays sem duplicatas e mantém o valor do grupo quando os dois lados definem a mesma chave de forma diferente, para que as outras contas não vejam mudança. Toda escolha, exceto `keep group`, salva antes a cópia do grupo em `.ccam-backup/<time>/group/` da conta. O seu próprio arquivo sempre vai para o backup.

A cada inicialização, `share.Check` compara cada caminho do layout com os grupos da conta usando apenas `lstat` e `readlink`. Um link ausente é criado na hora. Um arquivo comum onde deveria haver um link, um link para o lugar errado ou um link obsoleto interrompe a inicialização e indica `ccam doctor --fix`.

### Renomear, remover e o formato antigo

- `ccam group rename` altera apenas a configuração. A pasta mantém o nome, e o grupo o registra como `dir`, então nada é movido e nada pode se perder.
- `ccam group remove` se recusa a remover um grupo com contas, a menos que você adicione `--leave`. Cada conta fica então com uma cópia, e a pasta vai para `~/.local/share/ccam/trash/`, a menos que outro grupo a use.
- Antes dos grupos, o `config.toml` guardava uma tabela `[profile.share]` por conta. O ccam lê esse formato como grupos sem mover nada. No primeiro salvamento, ele grava o formato novo e mantém o arquivo antigo como `config.toml.bak`. Os itens compartilhados pelas mesmas contas viram um grupo com o nome da pasta; os demais recebem nomes como `default-history` e continuam usando a mesma pasta.

### No Windows

O Windows só cria links simbólicos no Modo de Desenvolvedor (Configurações > Sistema > Para desenvolvedores) ou para um administrador. Sem essa permissão, uma pasta compartilhada vira uma junção de diretório, que funciona do mesmo jeito para o Claude Code. Um arquivo compartilhado (`settings.json`, `CLAUDE.md`, `history.jsonl`) só pode ser um link simbólico, então uma alteração de grupo que precise de um é recusada antes que qualquer coisa seja movida. Pelo mesmo motivo, o primeiro `ccam add` ali cria o grupo `default` só com `mcp` e `extensions` e avisa, e `ccam add --group` para um grupo que compartilha configurações, memória ou histórico é recusado.

## Segurança e casos especiais

- **Contas em execução.** Uma alteração é recusada enquanto o Claude Code estiver em execução com uma das contas que ela move, porque a sessão continua gravando nesses arquivos, tenha sido ela iniciada pelo ccam ou não (`claude` simples, uma IDE, sessões em segundo plano). A verificação é repetida logo antes de salvar. Sair e entrar em um grupo dizem respeito apenas à conta indicada; sair de `extensions` também reescreve caminhos nas listas de plugins do grupo, sem mudar o que eles significam para as outras. `--dry-run` ainda imprime o plano.
- **Backups.** Cada plano recebe a própria pasta `.ccam-backup/<time>/` na pasta da conta. O ccam nunca as exclui.
- **Seus próprios links.** Se um caminho já é um link que você criou, por exemplo para um repositório de dotfiles, e o grupo está vazio, o ccam copia o destino do link para o grupo e move o próprio link para o backup. O destino não é alterado.
- **Pastas importadas.** Uma conta registrada no lugar com `ccam import` (como `~/.claude`) é alterada no lugar, então o `claude` simples também vê a mudança; o plano avisa isso. Quando ela sai de um grupo, sempre fica com uma cópia, também com `--fresh`, e o `ccam remove` a tira dos grupos dessa forma antes de esquecer a conta. Os arquivos compartilhados dela ficam na pasta do grupo enquanto ela está em um grupo, então remova a conta (ou saia dos grupos) antes de excluir `~/.local/share/ccam`.
- **Duas contas na mesma pasta.** Contas que compartilham o histórico gravam nos mesmos arquivos. Quando uma segunda conta inicia em uma pasta onde outra já está em execução, o ccam avisa e pergunta.
- **Edições manuais.** Um membro que não é uma conta é ignorado e informado pelo `ccam doctor`. Os nomes `local`, `on`, `off` e `none` não podem ser nomes de grupos.

## Relacionados

- [O painel](/features/dashboard)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Servidores MCP compartilhados](/features/mcp)
- [Copiar uma configuração existente](/features/migrate)
- [Todos os comandos](/docs/commands)
