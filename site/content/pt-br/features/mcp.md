---
title: Servidores MCP compartilhados
description: Mantenha uma lista de servidores MCP por grupo com ccam mcp add, add-json, list, remove e import, passada com --mcp-config ou sincronizada em cada conta.
updated: 2026-10-20
order: 5
translation: machine
sourceHash: "6d3b3c38c6cb"
---

O Claude Code mantém os servidores MCP de usuário no `.claude.json`, o mesmo arquivo que guarda o login da conta, então esse arquivo não pode ser compartilhado entre contas. O ccam mantém os servidores compartilhados em um arquivo de grupo, `mcp.json`, e o entrega a todas as contas do grupo na inicialização. Use `ccam mcp` para adicionar um servidor uma única vez para todas as suas contas, e `ccam mcp import` para levar às outras um servidor que você adicionou dentro de uma sessão.

## Uso

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Uma alteração vale a partir da próxima inicialização de cada conta; as sessões em execução mantêm os servidores com que começaram:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

Adicionar um nome que o grupo já tem substitui a definição dele e imprime `updated` em vez de `added`. Os servidores só chegam a uma conta se ela está nesse grupo e o grupo compartilha MCP (`ccam group set <group> --add mcp`); quando o grupo não compartilha, o `add` avisa.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `-g`, `--group <group>` | grupo de destino de `add`, `add-json` e `remove` | o único grupo pelo qual alguma conta compartilha MCP; caso contrário, `default` |
| `-- <command> [args…]` (`add`) | executa este comando como um servidor stdio | |
| `--url <url>` (`add`) | URL de um servidor remoto (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transporte de um servidor `--url` | `http` |
| `--header <'Name: value'>` (`add`) | cabeçalho HTTP; pode ser repetida | nenhum |
| `-e`, `--env <KEY=VALUE>` (`add`) | variável de ambiente de um servidor stdio; pode ser repetida | nenhuma |
| `--force` (`import`) | substitui um servidor do grupo com o mesmo nome que seja diferente; o `mcp.json` antigo ganha backup | desativado |
| `--keep-local` (`import`) | mantém as entradas no `.claude.json` da conta | desativado |
| `-h`, `--help` | ajuda de cada comando | |
| `--dry-run` (global) | respeitado apenas por `mcp list`; `add`, `add-json`, `remove` e `import` o recusam e não alteram nada | desativado |

`add` precisa de um comando depois de `--` ou de `--url`; caso contrário, falha com `a command to run (after --) or --url is required`. `--header` vale para servidores `--url`, e `--env` para servidores de comando. `add-json` aceita o mesmo JSON que `claude mcp add-json`, e a definição precisa ser um objeto com `command` ou `url`. `import` sempre usa o grupo pelo qual a conta compartilha MCP, e `list` mostra todos os grupos, então `-g` não os altera.

## Como funciona

Os arquivos envolvidos:

- `~/.local/share/ccam/groups/<group>/mcp.json`: os servidores compartilhados, `{"mcpServers": {…}}`, o formato que `--mcp-config` lê. Ele é gravado com `fsutil.AtomicWrite` e modo `0600`, porque as entradas de servidores costumam conter chaves de API.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: o arquivo do próprio Claude Code com os servidores da conta. O ccam lê só o valor `mcpServers` com o scanner `profile.TopLevelValue` e nunca analisa o arquivo inteiro no caminho de inicialização.
- `~/.local/state/ccam/mcp/<name>/`: o arquivo `seen`, o `sync-stamp` e os backups do `.claude.json`.

Há duas estratégias, definidas por `mcp_strategy` em `~/.config/ccam/config.toml`. A padrão é a estratégia flag:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` encontra o grupo de MCP da conta. `mcp.FlagArgs` verifica com um único `stat` que o `mcp.json` do grupo existe, sem analisá-lo.
2. `Plan.InsertArgs` coloca `--mcp-config=<path>` logo depois de `claude`, antes dos seus próprios argumentos, então as suas opções `--mcp-config` ou `--strict-mcp-config` continuam valendo por cima.
3. `mcp.NewLocal` verifica se o `.claude.json` da conta mudou desde a última inicialização e só então analisa os servidores dele.
4. `syscall.Exec` substitui o ccam pelo Claude Code.

O caminho é passado como um único argumento, `--mcp-config=<path>`, e não como dois. O `--mcp-config` do Claude Code aceita vários valores, então `claude --mcp-config a.json "prompt"` leria o prompt como um segundo arquivo de configuração e falharia com `MCP config file not found`. A forma de um argumento termina ali, e um prompt ou um subcomando como `mcp list` que venha em seguida funciona. O Claude Code carrega os servidores de `--mcp-config` além dos da própria conta; em um conflito de nomes, vale a entrada do `.claude.json`. A estratégia flag só grava no `.claude.json` para o Remote Control.

`claude remote-control` (`rc`) não aceita `--mcp-config`, e as sessões dele leem os servidores só do `.claude.json`. Para `ccam work rc`, o ccam deixa a opção de fora e, em um terminal, pergunta se deve copiar os servidores do grupo para o `.claude.json`, com um backup, como a estratégia `sync` faz. Ele pergunta até existir uma cópia e, depois, só quando o `mcp.json` do grupo mudar ([Remote Control](/features/launch#remote-control)). Como essa cópia prevalece sobre `--mcp-config` em um conflito de nomes, `ccam mcp add`, `add-json`, `remove` e `import` atualizam as cópias que outras contas do grupo têm:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Um servidor adicionado dentro de uma sessão com `claude mcp add` vai apenas para o `.claude.json` dessa conta. O ccam não consegue vê-lo nesse momento, porque o Claude Code o substituiu, então ele verifica na próxima inicialização. O arquivo `seen` guarda o horário de modificação do `.claude.json` e os nomes de servidores já informados. Quando o horário é diferente, o ccam compara os nomes com o grupo e imprime uma dica para cada servidor novo, uma única vez, sem esperar entrada:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

Com `mcp_strategy = "sync"`, os servidores aparecem em `claude mcp list` e em `/mcp` como servidores comuns de usuário, e o ccam os grava no `.claude.json` de cada conta antes da inicialização:

<!-- motion:mcp-sync -->

1. `mcp.Sync` compara o horário de modificação do `mcp.json` do grupo com o `sync-stamp`. Se ele não mudou, o `.claude.json` da conta nem é aberto.
2. Ele obtém o bloqueio do próprio Claude Code, a pasta `.claude.json.lock` ao lado do arquivo, criada com `mkdir`. Ele espera até 2 segundos e considera abandonado um bloqueio com mais de 10 segundos, como faz a biblioteca de bloqueio do Claude Code.
3. Ele salva o `.claude.json` atual em `backups/` na pasta de estado e mantém os 5 mais recentes.
4. Ele define cada servidor do grupo em `mcpServers`, remove os servidores que gravou antes e que o grupo não tem mais, mantém os servidores próprios da conta e grava o arquivo de forma atômica, com o modo original.
5. Ele registra o novo carimbo e os nomes que gerencia, e o Claude Code inicia sem `--mcp-config`.

A sincronização é pulada, com `ccam: account work is already running, skipping MCP sync`, quando uma sessão da conta está em execução, porque o Claude Code reescreve o `.claude.json` enquanto roda. Uma sincronização que falha imprime `ccam: MCP sync failed (launching anyway): …` e inicia mesmo assim. A dica sobre servidores adicionados dentro de uma sessão existe só na estratégia flag.

`ccam mcp import` move os servidores próprios de uma conta para o grupo dela:

<!-- motion:mcp-import -->

O ccam lê os servidores da conta com `mcp.LocalServers` e compara cada um com o servidor de mesmo nome do grupo usando `mcp.SameDef`, que ignora a ordem das chaves e os espaços. Um servidor que o grupo não tem, ou que tem com a mesma definição, é copiado. Um diferente é ignorado, porque a definição do grupo é a que todas as outras contas usam:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Com `--force`, o servidor do grupo é substituído, e o `mcp.json` anterior é salvo antes em `.ccam-backup/<time>/mcp.json` dentro da pasta do grupo. Depois, `mcp.RemoveLocal` remove as entradas copiadas do `.claude.json` da conta, com o mesmo bloqueio e depois de um backup na pasta de estado. Isso é necessário porque, em um conflito de nomes, a entrada da conta prevaleceria sobre a do grupo. `--keep-local` as mantém.

## Segurança e casos especiais

- **Logins.** Servidores que fazem login com OAuth guardam os tokens no `.credentials.json` de cada conta (`mcpOAuth`), não na lista de servidores, então compartilhar um servidor não compartilha o login dele. Entre uma vez em cada conta com `/mcp` ou copie os logins do `~/.claude` com [`ccam migrate`](/features/migrate). O grupo nunca recebe um arquivo de login.
- **Contas em execução.** `import` copia para o grupo, mas mantém as entradas próprias de uma conta em execução e avisa; execute-o de novo depois que a sessão terminar. A sincronização pula as contas em execução.
- **Cópias para o Remote Control.** A cópia de uma conta em execução não é alterada: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Com `mcp_strategy = "sync"`, nada precisa ser atualizado, porque cada inicialização sincroniza.
- **Arquivos corrompidos.** Se o `.claude.json` não for um JSON válido, a sincronização não mexe nele e informa o erro. Se o bloqueio não puder ser obtido em 2 segundos, o ccam informa `could not lock … (Claude Code is writing the file)`.
- **Contas que não compartilham MCP.** `import` se recusa e pede que você coloque antes a conta em um grupo que compartilha MCP (`ccam group join <group> side`).
- **Nomes de grupos.** `-g` aceita apenas um nome de grupo; ele vira uma pasta em `groups/`, então nomes como `../x` são recusados.

## Relacionados

- [Compartilhamento entre contas](/features/share)
- [Usar o Remote Control com o ccam](/guides/remote-control)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Copiar uma configuração existente](/features/migrate)
- [Todos os comandos](/docs/commands)
