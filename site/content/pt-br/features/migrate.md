---
title: Copiar uma configuração existente
description: Copie configurações, CLAUDE.md, extensões, servidores MCP e logins de MCP de um ~/.claude existente para um grupo de compartilhamento do ccam com ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

Se você usava o Claude Code antes do ccam, sua configuração está em `~/.claude` e `~/.claude.json`. O `ccam migrate` a copia para um grupo de compartilhamento, para que todas as contas que compartilham esses itens pelo grupo comecem com a configuração que você já tem. A pasta de origem é apenas lida, e o `claude` simples continua funcionando como antes. Execute-o uma vez depois do seu primeiro `ccam add` e de novo sempre que uma conta passar a compartilhar servidores MCP.

## Uso

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

Em um terminal, sem opções de itens e sem `--yes`, o ccam mostra primeiro uma lista de seleção:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Espaço (ou `x`) alterna o item sob o cursor, as setas ou `j`/`k` movem a seleção, Enter mostra o plano, e Esc, `q` ou Ctrl+C imprime `Cancelled.`. Itens para os quais a origem não tem nada aparecem como `not found`. Se você desmarcar tudo, o ccam imprime `Nothing selected.` e para.

Em seguida, o ccam imprime o plano e pergunta antes de copiar:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Um grupo só chega às contas que compartilham um item por meio dele, então o resumo as lista por item. Para as contas que não compartilham, o ccam imprime os comandos `ccam share <name> --<item> on` que ativam o compartilhamento.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `[folder]` | A pasta do Claude Code de onde copiar | `~/.claude` |
| `--settings` | Copia o `settings.json` | desativado |
| `--memory` | Copia o `CLAUDE.md` global | desativado |
| `--extensions` | Copia `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | desativado |
| `--mcp` | Copia os servidores MCP de usuário (`mcpServers` no `.claude.json`) | desativado |
| `--history` | Copia o histórico de conversas (`projects/`, `history.jsonl`) | desativado |
| `-g`, `--group string` | O grupo de compartilhamento para onde copiar | `default` |
| `-y`, `--yes` | Copia sem perguntar; em um conflito, a cópia do grupo é mantida | desativado |
| `--dry-run` (global) | Imprime o plano e para antes de qualquer pergunta ou cópia | desativado |

Sem opções de itens, a lista de seleção decide em um terminal. Fora de um terminal, ou com `--yes`, o ccam copia tudo, exceto o histórico de conversas. O histórico fica de fora por padrão porque `projects/` pode ser grande.

`--group on` significa `default`. `--group off` e `--group local` são recusados com `--group takes a group name, not "off"`. O `ccam migrate` respeita `--dry-run`, e a lista de seleção continua aparecendo em uma simulação, para que você possa ver uma seleção antes.

## Como funciona

`internal/cli/migrate.go` escolhe os itens, e então `migrate.Build` lê a origem e o grupo e retorna um plano. Nada é gravado até `Plan.Apply`, depois que você confirma.

1. **Encontrar a origem.** `checkSource` recusa um caminho que não existe, que não é uma pasta ou que não contém nenhuma entrada conhecida do Claude Code. `migrate.ClaudeJSON` encontra o `.claude.json` que pertence à pasta: dentro dela ou, para um `~/.claude` sem um interno, o `~/.claude.json` ao lado dele. É ali que o Claude Code o mantém quando `CLAUDE_CONFIG_DIR` não está definida.
2. **Associar itens a caminhos.** Todos os caminhos vêm de `share.Layout`, a mesma lista de permissões que o `ccam share` usa. Os servidores MCP são lidos do `.claude.json` com `mcp.LocalServers`, que usa o scanner `profile.TopLevelValue` em vez de analisar o arquivo inteiro.
3. **Copiar, nunca mover.** Os arquivos vão para `~/.local/share/ccam/groups/<group>/` com `fsutil.AtomicWrite` ou `fsutil.CopyFile`; os servidores MCP vão para o `mcp.json` do grupo.

Cada item associa um caminho de origem a um caminho do grupo:

<!-- motion:migrate-items -->

As duas linhas do histórico são tracejadas porque só são copiadas com `--history`. `.credentials.json` e `.claude.json` estão em `share.Never` e nunca chegam ao grupo: `Build` retorna um erro interno se algum deles aparecer no layout, e a varredura de pastas os ignora pelo nome.

`buildEntry` atribui um estado a cada caminho, e é isso que cada linha do plano descreve:

<!-- motion:migrate-states -->

- **Ausente**: não está na origem, ou é um arquivo onde se espera uma pasta. Listado uma vez por item como `not found, skipped`.
- **Vinculado**: o caminho de origem já aponta para a cópia do grupo, por exemplo quando a origem é uma pasta de perfil do ccam.
- **Igual**: o grupo já tem conteúdo idêntico.
- **Novo**: o grupo não o tem, então ele é copiado inteiro.
- **Mesclar**: uma pasta dos dois lados. `copyMissing` adiciona só os arquivos que o grupo não tem; arquivos diferentes ficam como estão no grupo, e o plano diz `group copies kept`. Em `plugins/`, os registros `installed_plugins.json` e `known_marketplaces.json` são mesclados, e os caminhos absolutos de plugins são reescritos para apontar para a pasta do grupo.
- **Conflito**: um arquivo dos dois lados com conteúdo diferente. Você escolhe por arquivo: `keep group`, `use mine` e também `merge JSON` para arquivos `.json` (uma mesclagem profunda em que o grupo prevalece nos valores conflitantes) ou `merge by line` para `history.jsonl`. Um servidor MCP com definição diferente também é um conflito, com `keep group` ou `use mine`. A primeira opção listada é a padrão quando você pressiona Enter; `--yes` sempre mantém a cópia do grupo.

A lista de seleção em si é `tui.MigratePicker`, construída sobre o componente `checklist` que o `ccam add` usa; `migrate.Found` marca os itens ausentes com chamadas lstat e com a lista de servidores MCP.

### Logins de MCP

Servidores MCP que fazem login com OAuth guardam os tokens no arquivo de login, `.credentials.json` em `mcpOAuth`, e não na lista de servidores. O Claude Code identifica cada entrada por `<name>|` mais os 16 primeiros caracteres de um hash `sha256` de `type`, `url` e `headers` do servidor. A chave não depende da pasta, então uma entrada copiada funciona em outra conta desde que a definição do servidor seja idêntica. Por isso, `buildLogins` em `internal/migrate/logins.go` copia só as entradas dos servidores que o grupo acaba tendo com a definição da origem: novos, já iguais ou um conflito em que você escolheu `use mine`. Os logins de servidores de plugins (chaves que começam com `plugin:`) vão junto quando você copia extensões e a origem tem `plugins/`.

<!-- motion:migrate-signins -->

As entradas vão para o próprio `.credentials.json` de cada conta que as recebe, gravado com `fsutil.AtomicWrite` no modo `0600`. Elas nunca vão para o grupo. Uma conta só as recebe se compartilha `mcp` (ou, para servidores de plugins, `extensions`) pelo grupo de destino, e só as chaves que ela não tem são adicionadas: as entradas que ela já tem, talvez mais recentes, são mantidas. A linha `sign-ins` do plano lista os servidores, as contas que os recebem e as que foram ignoradas. O ccam ignora uma conta, e a indica com o motivo, quando ela:

- foi registrada no lugar com `ccam import`, já que o `claude` simples pode estar usando o arquivo de login dela;
- está em execução, já que o Claude Code reescreve o arquivo de login dela (verificado de novo logo antes da gravação, informado como `started meanwhile`);
- ainda não fez login;
- guarda o login no Acesso às Chaves do macOS, onde o ccam não grava.

Em uma conta ignorada, entre nos servidores com `/mcp`. Uma conta que passar a compartilhar servidores MCP depois recebe os logins quando você executar `ccam migrate` de novo; o resumo avisa com ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Segurança e casos especiais

- A origem é apenas lida. Nenhum arquivo nela é movido, vinculado ou alterado, e o login da origem nunca é gravado, inclusive no macOS, onde ele pode ser lido do Acesso às Chaves.
- Antes que um arquivo do grupo seja substituído, por `use mine`, por uma mesclagem ou por novos servidores MCP em um `mcp.json` existente, a cópia do grupo é salva em `<group>/.ccam-backup/<time>/`. O nome da pasta é único por execução, então duas execuções no mesmo segundo não sobrescrevem os backups uma da outra. O resumo imprime o caminho.
- Com `--yes`, todo conflito mantém a cópia do grupo, e o ccam explica como escolher: execute o mesmo comando sem `--yes`.
- Se nada mudaria, o ccam imprime `Nothing to copy.`; se todos os conflitos mantêm a cópia do grupo, `Nothing to copy: the group keeps its own copies.`
- Se uma cópia falhar no meio, o ccam para com `The source folder was not changed. Group files replaced so far are saved in` seguido da pasta de backup.
- Uma origem ausente termina com `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Relacionados

- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Verificar e reparar](/features/doctor)
- [Várias contas do Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
