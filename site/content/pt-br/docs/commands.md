---
title: Referência de comandos
description: Todos os comandos e opções do ccam, de ccam add e ccam <name> a ccam usage, ccam next e ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 54bdfc489174
---

Coloque `--dry-run` antes de qualquer comando para ver o que ele faria sem fazê-lo, por exemplo `ccam --dry-run work`.

<!-- commands -->

## Iniciar o Claude Code

`ccam <name|number> [arguments]` inicia o Claude Code como essa conta. O ccam não lê os argumentos depois do nome; ele os repassa ao `claude` na mesma ordem, e a entrada padrão, o código de saída e os sinais se comportam exatamente como no próprio `claude`. Um `--` depois do nome é opcional.

`ccam` com argumentos, mas sem nome de conta, como `ccam -p "summarise"`, usa a conta padrão. Defina-a com `ccam use <name>`.

Os atalhos `claude1`, `claude2` e assim por diante são links para o ccam em `~/.local/bin`. `ccam shims install --names` também cria `claude-<name>` para cada conta.

## Adicionar e remover contas

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` sem uma pasta procura `~/.claude` e `~/.claude-*` e pergunta sobre cada uma. As pastas importadas são registradas no lugar onde estão; o ccam nunca move nem copia um login.

## Compartilhamento

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Cada item aceita `on` (o grupo `default`), `off` ou o nome de um grupo. Consulte [compartilhamento](/docs/sharing).

## Servidores MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Uso e troca de conta

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Manutenção

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Arquivos

| Caminho | Conteúdo |
| --- | --- |
| `~/.config/ccam/config.toml` | As contas e o que elas compartilham |
| `~/.local/share/ccam/profiles/<name>/` | O `CLAUDE_CONFIG_DIR` de cada conta |
| `~/.local/share/ccam/groups/<group>/` | Configurações, servidores MCP, extensões e histórico compartilhados |
| `~/.local/state/ccam/` | Cache de uso e horários do último uso |
| `~/.local/bin/claude1`, `claude2`… | Links para o ccam |

As variáveis `XDG_CONFIG_HOME`, `XDG_DATA_HOME` e `XDG_STATE_HOME` mudam o local dessas pastas.
