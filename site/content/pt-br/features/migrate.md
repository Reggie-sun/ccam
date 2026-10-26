---
title: Copiar uma configuração existente
description: Copie configurações, CLAUDE.md, extensões e servidores MCP de um ~/.claude existente para um grupo do ccam com ccam migrate.
updated: 2026-10-26
order: 6
translation: machine
sourceHash: "f835fbe84ade"
---

Se você usava o Claude Code antes do ccam, sua configuração está em `~/.claude` e `~/.claude.json`. O `ccam migrate` a copia para um grupo de compartilhamento, para que todas as contas que compartilham esses itens pelo grupo comecem com a configuração que você já tem. A pasta de origem é apenas lida, e o `claude` simples continua funcionando como antes. Execute-o uma vez depois do seu primeiro `ccam add`; quando o `~/.claude.json` tem servidores MCP, o primeiro `ccam add` os indica na dica dele.

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
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Espaço (ou `x`) alterna o item sob o cursor, as setas ou `j`/`k` movem a seleção, Enter mostra o plano, e Esc, `q` ou Ctrl+C imprime `Cancelled.`. Itens para os quais a origem não tem nada aparecem como `not found`. Se você desmarcar tudo, o ccam imprime `Nothing selected.` e para.

Em seguida, o ccam imprime o plano e pergunta antes de copiar:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Um grupo só afeta as próprias contas, e só nos itens que compartilha, então o resumo lista as contas do grupo. Quando o grupo não compartilha um dos itens copiados, ou algumas contas não estão nele, o ccam imprime os comandos `ccam group set <group> --add <items>` e `ccam group join <group> <accounts>` que mudam isso.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `[folder]` | A pasta do Claude Code de onde copiar | `~/.claude` |
| `--settings` | Copia o `settings.json` | desativado |
| `--memory` | Copia o `CLAUDE.md` global | desativado |
| `--extensions` | Copia `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | desativado |
| `--mcp` | Copia os servidores MCP de usuário (`mcpServers` no `.claude.json`) | desativado |
| `--history` | Copia o histórico de conversas (`projects/`, `history.jsonl`) | desativado |
| `-g`, `--group <group>` | O grupo de compartilhamento para onde copiar | `default` |
| `-y`, `--yes` | Copia sem perguntar; em um conflito, a cópia do grupo é mantida | desativado |
| `--dry-run` (global) | Imprime o plano e para antes de qualquer pergunta ou cópia | desativado |

Em um terminal, se você não informar opções de itens, você escolhe os itens na lista de seleção. Fora de um terminal, ou com `--yes`, o ccam copia tudo, exceto o histórico de conversas. O histórico fica de fora por padrão porque `projects/` pode ser grande.

`--group on` significa `default`. `--group off` e `--group local` são recusados com `--group takes a group name, not "off"`. O `ccam migrate` aceita `--dry-run`, e a lista de seleção continua aparecendo em uma simulação, para que você possa ver uma seleção antes.

## Como funciona

O ccam monta um plano completo antes de copiar qualquer coisa. Ele escolhe os itens (`internal/cli/migrate.go`), depois lê a origem e o grupo e monta o plano (`migrate.Build`). Nada é gravado até `Plan.Apply`, depois que você confirma.

1. **Encontrar a origem.** O ccam recusa um caminho que não existe, que não é uma pasta ou que não contém nenhuma entrada conhecida do Claude Code (`checkSource`). Só para servidores MCP, um `~/.claude` vazio ou ausente não é problema, desde que o `~/.claude.json` tenha servidores. O ccam encontra o `.claude.json` que pertence à pasta (`migrate.ClaudeJSON`): dentro dela ou, para um `~/.claude` sem um interno, o `~/.claude.json` ao lado dele. É ali que o Claude Code o mantém quando `CLAUDE_CONFIG_DIR` não está definida.
2. **Associar itens a caminhos.** Todos os caminhos vêm de `share.Layout`, a mesma lista de permissões que os grupos usam. Os servidores MCP são lidos do `.claude.json` com `mcp.LocalServers`, que usa o scanner `profile.TopLevelValue` em vez de analisar o arquivo inteiro.
3. **Copiar, nunca mover.** Os arquivos vão para `~/.local/share/ccam/groups/<group>/` com `fsutil.AtomicWrite` ou `fsutil.CopyFile`, que grava um arquivo temporário ao lado do destino e o renomeia, para que uma cópia interrompida nunca deixe um arquivo parcial. Os servidores MCP vão para o `mcp.json` do grupo com o bloqueio dele: um servidor que outro processo do ccam adicione nesse meio-tempo é mantido, e as cópias feitas para o [Remote Control](/guides/remote-control) são atualizadas depois.

Cada item associa um caminho de origem a um caminho do grupo:

<!-- motion:migrate-items -->

As duas linhas do histórico são tracejadas porque só são copiadas com `--history`. `.credentials.json` e `.claude.json` estão em `share.Never` e nunca chegam ao grupo: o ccam para com um erro interno se algum deles aparecer no layout (`Build`), e a varredura de pastas os ignora pelo nome.

O ccam atribui um estado a cada caminho, e é isso que cada linha do plano descreve (`buildEntry`):

<!-- motion:migrate-states -->

- **Ausente**: não está na origem, ou é um arquivo onde se espera uma pasta. Listado uma vez por item como `not found, skipped`.
- **Vinculado**: o caminho de origem já aponta para a cópia do grupo, por exemplo quando a origem é uma pasta de perfil do ccam.
- **Igual**: o grupo já tem conteúdo idêntico.
- **Novo**: o grupo não o tem, então ele é copiado inteiro.
- **Mesclar**: uma pasta dos dois lados. O ccam adiciona só os arquivos que o grupo não tem (`copyMissing`); arquivos diferentes ficam como estão no grupo, e o plano diz `group copies kept`. Em `plugins/`, os registros `installed_plugins.json` e `known_marketplaces.json` são mesclados, e os caminhos absolutos de plugins são reescritos para apontar para a pasta do grupo.
- **Conflito**: um arquivo dos dois lados com conteúdo diferente. Você escolhe por arquivo: `keep group`, `use mine` e também `merge JSON` para arquivos `.json` (uma mesclagem profunda em que o grupo prevalece nos valores conflitantes) ou `merge by line` para `history.jsonl`. Um servidor MCP com definição diferente também é um conflito, com `keep group` ou `use mine`. A primeira opção listada é a padrão quando você pressiona Enter; `--yes` sempre mantém a cópia do grupo.

A lista de seleção em si é `tui.MigratePicker`, construída sobre o componente `checklist` que o `ccam add` usa; o ccam marca os itens ausentes com chamadas lstat e com a lista de servidores MCP (`migrate.Found`).

### Logins de MCP

Servidores MCP que fazem login com OAuth guardam os tokens no arquivo de login, `.credentials.json` em `mcpOAuth`, e não na lista de servidores. O ccam não os copia. Um login copiado tem o mesmo refresh token que o original, e o servidor troca esse token na primeira renovação: todos os outros que o têm, inclusive o `~/.claude` de origem, são então desconectados daquele servidor. Por isso, para cada servidor copiado em que a origem tem login, a linha `sign-ins` do plano lista as contas do grupo que precisam entrar nele por conta própria com `/mcp`, e o resumo repete isso. Uma entrada sem token não é um login e não é listada, assim como uma conta que já tem o próprio login ou um servidor cuja definição o grupo mantém. Nenhum arquivo de login é gravado.

## Segurança e casos especiais

- A origem é apenas lida. Nenhum arquivo nela é movido, vinculado ou alterado, e o login da origem nunca é gravado, inclusive no macOS, onde ele pode ser lido do Acesso às Chaves.
- Antes que um arquivo do grupo seja substituído, por `use mine`, por uma mesclagem ou por novos servidores MCP em um `mcp.json` existente, a cópia do grupo é salva em `<group>/.ccam-backup/<time>/`. O nome da pasta é único por execução, então duas execuções no mesmo segundo não sobrescrevem os backups uma da outra. O resumo imprime o caminho.
- Com `--yes`, todo conflito mantém a cópia do grupo, e o ccam explica como escolher: execute o mesmo comando sem `--yes`.
- Se nada mudaria, o ccam imprime `Nothing to copy.`; se todos os conflitos mantêm a cópia do grupo, `Nothing to copy: the group keeps its own copies.`
- Se uma cópia falhar no meio, o ccam para com `The source folder was not changed. Group files replaced so far are saved in` seguido da pasta de backup. Executá-lo de novo copia os arquivos que não foram copiados por completo.
- Quando vários grupos compartilham MCP, os comandos que `ccam doctor`, `ccam ext` e `ccam mcp list` sugerem são `ccam migrate --mcp --group <a|b>`.
- Uma origem ausente termina com `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Relacionados

- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Verificar e reparar](/features/doctor)
- [Usar várias contas do Claude Code no mesmo computador](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR explicado](/guides/claude-config-dir)
