---
title: Usar o Remote Control com o ccam
description: "Inicie o Remote Control do Claude Code com qualquer conta usando ccam <name> rc: opções válidas, por que o ccam oferece copiar os servidores MCP e soluções."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "5f0b90ea53d6"
---

`ccam work rc` inicia o Remote Control do Claude Code com a conta `work`, para que você possa conduzir sessões neste computador pelo claude.ai/code ou pelo app móvel do Claude. Ele funciona como `claude rc`, com uma diferença: se a conta compartilha servidores MCP, o ccam pergunta uma vez se deve copiá-los para o `.claude.json` da conta, porque as sessões do Remote Control não leem servidores de nenhum outro lugar.

## Como iniciar

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

Um [atalho](/features/shortcuts) funciona do mesmo jeito: `ccwork rc` equivale a `ccam work rc`. Execute o comando na pasta em que você quer trabalhar e depois abra o claude.ai/code ou o app do Claude, conectado com a mesma conta Claude de `work`.

## Opções

Dois tipos de opção se encontram aqui: as do próprio ccam, que vêm antes do nome da conta ou ficam nas configurações do ccam, e as do Remote Control, que vêm depois de `rc` e chegam ao Claude Code sem alterações.

### Opções e configurações do ccam

| Opção ou configuração | O que faz | Padrão |
| --- | --- | --- |
| `--dry-run` (antes da conta) | Imprime o ambiente e o comando, além do aviso de MCP, se houver. Não pergunta nem copia nada | desativado |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Perguntado quando a conta compartilha servidores MCP e o `.claude.json` dela não tem a lista atual do grupo. `y` copia os servidores; `N` ou Enter inicia sem eles | `N` |
| `mcp_strategy = "sync"` em `~/.config/ccam/config.toml` | Copia os servidores do grupo para o `.claude.json` antes de toda inicialização, seja do Remote Control ou não, então nada é perguntado | `"flag"` |
| `ccam group leave <group> work` | Tira a conta do grupo que compartilha MCP com ela; o Remote Control passa a iniciar só com os servidores próprios da conta, sem perguntar nada | no `default` |

### Opções do Remote Control

Tudo o que vem depois de `rc` pertence ao Claude Code. Estas são as opções de `claude remote-control` no Claude Code 2.1.304; `ccam work rc --help` imprime a lista da sua versão.

| Opção | O que faz |
| --- | --- |
| `--name <name>` | Nome da sessão mostrado no claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Prefixo dos nomes de sessão gerados (padrão: o nome do host) |
| `-c`, `--continue` | Reconecta à última sessão que o Remote Control registrou para esta pasta, dentro de cerca de 4 horas |
| `--session-id <id>` | Reconecta a uma sessão específica |
| `--permission-mode <mode>` | Modo de permissão das novas sessões: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` ou `plan` |
| `--spawn <mode>` | `same-dir` (padrão), `worktree` (um git worktree por sessão) ou `session` (uma única sessão, que encerra quando ela termina) |
| `--capacity <N>` | Máximo de sessões simultâneas nos modos `same-dir` e `worktree` (padrão: 32) |
| `--[no-]create-session-in-dir` | Cria uma primeira sessão na pasta atual ao iniciar (ativado por padrão) |
| `--[no-]chrome` | Claude in Chrome nas novas sessões (desativado por padrão) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Saída de depuração e detalhada |

Coloque `rc` logo depois do nome da conta. O ccam toma como comando do Claude Code o primeiro argumento que não começa com `-`, então em `ccam work rc --name laptop` ele vê `rc`, mas em `ccam work --model opus rc` ele vê `opus`.

## O que o ccam faz

A inicialização segue o caminho de sempre ([Iniciar o Claude Code](/features/launch)), com uma mudança na etapa de MCP:

1. O ccam encontra a conta, verifica os links de compartilhamento dela e resolve o `claude`, como em qualquer inicialização.
2. Ele vê que o comando é `rc` ou `remote-control` e deixa `--mcp-config` de fora.
3. Se a conta compartilha servidores MCP, ele compara o `mcp.json` do grupo com o `sync-stamp` em `~/.local/state/ccam/mcp/work/`. Se o `.claude.json` já tiver os servidores atuais, ou se o grupo não tiver nenhum, ele vai direto para a etapa 6.
4. Caso contrário, ele lista os servidores que faltariam às sessões:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. Com `y`, ele salva um backup do `.claude.json`, obtém o bloqueio do próprio Claude Code, acrescenta os servidores do grupo a `mcpServers` e mantém os outros servidores da conta. Com `N`, sem um terminal ou enquanto a conta está em execução, ele não altera nada e avisa:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Ele registra a sessão, define `CLAUDE_CONFIG_DIR` como a pasta da conta e se substitui por `claude rc` com `exec`.

A partir daí, o Remote Control continua rodando como servidor. Cada sessão que ele inicia para você é um processo filho desse processo, então roda como `work`: com o login, as configurações e o histórico de conversas de `work`, e conta para o limite de uso de `work`. Enquanto o servidor estiver rodando, o ccam considera a conta em execução.

A simulação mostra a diferença. Uma inicialização normal leva os servidores do grupo como argumento; o Remote Control, não:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Por que funciona assim

**O Remote Control não aceita `--mcp-config`.** Nas inicializações normais, o ccam entrega os servidores do grupo ao Claude Code como `--mcp-config=<path>` e nunca grava nos arquivos do Claude Code. `claude remote-control` recusa essa opção antes do nome do comando e não aceita nenhuma depois dele, então, antes do ccam 0.6.0, `ccam work rc` falhava em toda conta que compartilhava servidores MCP. Agora o ccam deixa a opção de fora.

**As sessões dele leem servidores só do `.claude.json`.** As sessões que o Remote Control inicia não herdam opções da linha de comando. O único lugar onde elas encontram servidores MCP de usuário é o `.claude.json` da conta. Esse arquivo também guarda os dados de login da conta e o estado por pasta, então o ccam não pode vinculá-lo a uma cópia compartilhada como faz com o `settings.json`. Copiar os servidores para ele é a única forma de dar às sessões a lista compartilhada.

**O ccam pergunta antes de gravar.** Com a estratégia flag padrão, o ccam se compromete a não mexer no `.claude.json`. O Remote Control é o único caso em que ele precisa gravar ali, então pergunta antes, guarda um backup e não toca em nada além de `mcpServers`.

**Ele não grava enquanto a conta está em execução.** O Claude Code reescreve o `.claude.json` durante uma sessão, e uma alteração feita de fora seria perdida ou sobrescrita. Por isso o ccam pula a cópia quando qualquer sessão da conta está em execução, inclusive outro servidor do Remote Control.

**Ele pergunta uma vez a cada mudança.** Depois de uma cópia, `sync-stamp` guarda o horário de modificação do `mcp.json` do grupo. O ccam só pergunta de novo quando esse arquivo muda, por exemplo depois de `ccam mcp add`. Se você responder `N`, ainda não há cópia, então ele pergunta a cada inicialização do Remote Control.

**O `ccam mcp` mantém as cópias atualizadas.** Uma entrada no `.claude.json` prevalece sobre `--mcp-config` em um conflito de nomes, então uma cópia antiga esconderia uma definição mais nova também nas inicializações normais. Por isso `ccam mcp add`, `add-json`, `remove` e `import` atualizam as cópias de todas as contas que não estão em execução:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Sem um terminal, a resposta é não.** O Remote Control costuma rodar a partir de um script, de um serviço ou com `nohup`, onde ninguém pode responder. Nesse caso, o ccam inicia sem os servidores em vez de ficar esperando.

## Problemas e soluções

### Os servidores MCP compartilhados não aparecem nas sessões do Remote Control

A pergunta foi respondida com `N`, não havia terminal ou a conta estava em execução. Pare o Remote Control, feche as outras sessões da conta e execute `ccam work rc` de novo em um terminal; responda `y`. Para nunca ser perguntado, defina `mcp_strategy = "sync"` em `~/.config/ccam/config.toml`. Os servidores de uma sessão não mudam enquanto ela roda, então inicie uma nova sessão depois da cópia.

### "Account work is running, so its .claude.json is left alone"

Ainda há uma sessão de `work` aberta: um `ccam work` normal, um atalho ou outro servidor do Remote Control. Encerre-a e execute `ccam work rc` de novo. Sessões de outra conta não contam.

### O ccam pergunta toda vez

Você ainda não copiou os servidores. Responda `y` uma vez ou defina `mcp_strategy = "sync"`. Se você nunca quiser os servidores no Remote Control, pare de compartilhar MCP com a conta: tire-a desse grupo (`ccam group leave default work`) ou mantenha o MCP em um grupo próprio do qual a conta não faça parte.

### `ccam work rc` falha na hora com um erro sobre `--mcp-config`

Seu ccam é anterior à versão 0.6.0 e passa `--mcp-config` para o Remote Control. Execute `ccam update`. O mesmo erro com um ccam atual significa que uma opção com valor veio antes de `rc`, como em `ccam work --model opus rc`, e o ccam tomou `opus` como o comando; coloque `rc` logo depois do nome da conta.

### Uma alteração feita com `ccam mcp` não chega ao Remote Control

`ccam mcp add` e `remove` só atualizam as cópias nas contas que não estão em execução, e avisam: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Pare o Remote Control e inicie-o de novo; o ccam vê que o `mcp.json` mudou e pergunta outra vez.

### A cópia substituiu um servidor meu

Se a conta tinha um servidor próprio com o mesmo nome de um servidor do grupo, a cópia o substituiu. O `.claude.json` anterior está em `~/.local/state/ccam/mcp/work/backups/`, que mantém os cinco mais recentes. Para manter a sua versão em todas as contas, mova-a para o grupo com `ccam mcp import work <server> --force` antes da próxima inicialização do Remote Control.

### Quero remover os servidores copiados

`ccam mcp remove <server>` remove um servidor do grupo e de todas as cópias. Para voltar ao estado anterior à cópia, restaure o backup mais recente de `~/.local/state/ccam/mcp/work/backups/` enquanto nenhuma sessão da conta estiver rodando.

### O Remote Control pergunta se deve confiar na pasta, ou encerra sem um terminal

A confiança nas pastas fica guardada por conta no `.claude.json`, então uma pasta em que você confiou como `personal` é nova para `work`. Execute `ccam work` uma vez nessa pasta e aceite; depois inicie o Remote Control ali.

### "You must be logged in with a Claude account that has a subscription"

O Remote Control precisa de um login com assinatura do Claude. Confira o e-mail da conta em `ccam list` e, se necessário, faça login de novo com `ccam login work`.

### A sessão não aparece no celular

O app do Claude ou o claude.ai/code está conectado com outra conta Claude. Compare-a com o e-mail que `ccam list` mostra para `work`.

### O ccam avisa que outro Claude Code compartilha o histórico de conversas nesta pasta

Outra conta do mesmo grupo de histórico está em execução na mesma pasta. Sessões que continuam a mesma conversa a partir das duas contas podem misturar o histórico. Responda `y` só se elas forem trabalhar em conversas diferentes, ou inicie o Remote Control em outra pasta.

## Relacionados

- [Iniciar o Claude Code](/features/launch)
- [Servidores MCP compartilhados](/features/mcp)
- [Compartilhamento entre contas](/features/share)
- [Gerenciar contas](/features/accounts)
