---
title: Referência de comandos
metaTitle: Referência de comandos do ccam
description: Todos os comandos e opções do ccam, de ccam add e ccam <name> a ccam usage, ccam next e ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

Coloque `--dry-run` antes de um comando para ver o que ele faria sem fazê-lo, por exemplo `ccam --dry-run work`. Funciona ao iniciar o Claude Code (`ccam <name>`, `ccam run`, `ccam next`) e com `share`, `migrate`, `shortcut`, `update` e `doctor --fix`. Outros comandos que alteram algo recusam a opção em vez de executar.

<!-- commands -->

## Iniciar o Claude Code

`ccam <name|number> [arguments]` inicia o Claude Code como essa conta. O ccam não lê os argumentos depois do nome; ele os repassa ao `claude` na mesma ordem, e a entrada padrão, o código de saída e os sinais se comportam exatamente como no próprio `claude`. Um `--` depois do nome é opcional.

`ccam` com argumentos, mas sem nome de conta, como `ccam -p "summarise"`, usa a conta definida para a pasta atual ou, se não houver, a conta padrão. Defina-as com `ccam folder <name> [folder]` e `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` antes do nome da conta serve para execuções `-p`/`--print`. Quando a conta atinge o limite de uso, o ccam executa o mesmo comando de novo com a conta logada que tem mais espaço e imprime só a saída da execução que terminou. Consulte [Contas por pasta](/features/folder) e [Executar de novo no limite](/features/next#run-again-at-the-limit).

Um atalho como `ccwork` inicia a conta dele do mesmo jeito: `ccwork --resume` faz o mesmo que `ccam work --resume`. Consulte [atalhos](/docs/commands#shortcuts).

## Adicionar e remover contas

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` sem uma pasta procura `~/.claude` e `~/.claude-*` e pergunta sobre cada uma. As pastas importadas são registradas no lugar onde estão; o ccam nunca move nem copia um login.

`ccam add --api` adiciona uma conta que usa uma chave de API em vez de um login do Claude; o ccam lê a chave a cada início e nunca a guarda. Veja [Contas com chave de API](/features/api).

`ccam rename`, `ccam remove` e `ccam share` se recusam a alterar uma conta enquanto o Claude Code estiver em execução com essa conta. Encerre a sessão primeiro.

`ccam list --json` mostra todas as contas com o número (slot), atalho, e-mail e pasta de cada uma. O objeto `share` associa cada item (`mcp`, `settings`, `memory`, `extensions`, `history`) a um nome de grupo ou a `local`.

## Atalhos

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Um atalho é um comando que você mesmo nomeia e que inicia o Claude Code como uma determinada conta. Assim como `ccam work`, ele repassa todos os argumentos ao Claude Code sem alterações, então `ccwork --resume` e `ccwork -p "hi"` funcionam como você espera. Ele é um link para o ccam em `~/.local/bin`.

Os atalhos são opcionais. O assistente do `ccam add` pergunta se você quer criar um, `ccam add <name> --shortcut <command>` cria um sem perguntar e `ccam list` mostra o atalho de cada conta, ou `-` quando ela não tem.

Os nomes usam os caracteres a-z, 0-9, `.`, `-` e `_`. Um nome não pode ser um comando do ccam, o atalho de outra conta nem o nome de outro programa no seu `PATH`.

`ccam rename` mantém o nome do atalho, e `ccam remove` exclui o atalho junto com a conta. `ccam shims install` recria os links de todas as contas, por exemplo depois que você move o ccam.

## Compartilhamento

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Cada item aceita `on` (o grupo `default`), `off` ou o nome de um grupo. Consulte [compartilhamento](/docs/sharing).

## Copiar uma configuração existente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, extensões e servidores MCP de usuário de uma pasta do Claude Code para um grupo, para que todas as contas que compartilham esses itens os recebam. A pasta é apenas lida. As opções `--settings`, `--memory`, `--extensions`, `--mcp` e `--history` copiam só os itens que você indicar. Sem elas, em um terminal, uma lista de seleção permite escolher, com tudo marcado exceto o histórico de conversas; fora de um terminal, tudo é copiado, exceto o histórico de conversas. `-g` ou `--group` escolhe o grupo, e `-y` ou `--yes` copia sem perguntar e mantém a cópia do grupo quando um arquivo é diferente. Ele também copia os logins salvos dos servidores para o próprio arquivo de login de cada conta. Consulte [copiar uma configuração existente](/docs/sharing#copy-an-existing-setup).

## Servidores MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` move para o grupo da conta os servidores que só essa conta tem. Um servidor cujo nome o grupo já usa para uma definição diferente é ignorado; adicione `--force` para substituir a cópia do grupo, da qual o ccam faz backup antes.

## Uso e troca de conta

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

`next --auto`, `--fallback` e `ccam pick` escolhem do mesmo jeito: primeiro as contas abaixo de `auto_max_usage` (`config.toml`, padrão 100), por `auto_strategy` (`headroom` ou `expiring`), contando a janela semanal do `--model`, nunca uma conta em 100%. Quando nenhuma pode ser escolhida, saem com código 75, ou com `--wait` esperam o reinício da primeira conta. Veja [Trocar para a próxima conta](/features/next) e [Escolher uma conta para scripts](/features/pick).

## Manutenção

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` verifica o checksum SHA-256 da versão antes de substituir o ccam, e `--version v0.2.0` instala uma versão específica. Um ccam instalado pelo pacote `.deb` é atualizado com o apt. Consulte [atualização](/docs/install#update).

## Arquivos

| Caminho | Conteúdo |
| --- | --- |
| `~/.config/ccam/config.toml` | As contas e o que elas compartilham |
| `~/.local/share/ccam/profiles/<name>/` | O `CLAUDE_CONFIG_DIR` de cada conta |
| `~/.local/share/ccam/groups/<group>/` | Configurações, servidores MCP, extensões e histórico compartilhados |
| `~/.local/state/ccam/` | Cache de uso e horários do último uso |
| `~/.local/bin/<shortcut>` | Atalhos como `ccwork`, links para o ccam |

As variáveis `XDG_CONFIG_HOME`, `XDG_DATA_HOME` e `XDG_STATE_HOME` mudam o local dessas pastas. No Windows, `~` é `%USERPROFILE%` e os atalhos são `<shortcut>.exe`.
