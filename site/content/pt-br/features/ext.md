---
title: Inspecionar extensões
description: ccam ext mostra as skills, comandos, plugins, servidores MCP e hooks que cada conta carrega, de onde vem cada um e o que está quebrado.
updated: 2026-10-23
order: 14
translation: machine
sourceHash: "45740e4efbef"
---

`ccam ext` responde "por que esta skill (ou este servidor MCP) funciona em uma conta e não na outra?". Ele lista, para cada conta, as extensões que o Claude Code carrega e de onde vem cada uma: um grupo de compartilhamento, a pasta da própria conta ou o projeto em que você está. Abaixo da tabela, lista o que impede uma extensão de carregar como você espera. Ele só lê.

## Uso

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

Cada célula indica a origem daquele item naquela conta: o grupo de compartilhamento (`default`), `account` quando ele fica na pasta da própria conta, `local` para um servidor MCP que a conta tem só para esta pasta, ou `-` quando a conta não o tem. Um plugin instalado mas não ativado mostra `(off)` depois da origem. `!` marca um item com problema, explicado abaixo da tabela; problemas que não pertencem a nenhuma conta, como servidores que só o `claude` simples tem, vêm por último, sem nome. `ccam ext --json` imprime o mesmo relatório (`accounts`, o projeto, `problems`), com `remote` nos servidores MCP que o Claude Code acessa pela rede; sem contas, ele imprime o relatório com `accounts: []`.

## O que ele lê

| Tipo | Onde | Origem mostrada |
| --- | --- | --- |
| Skills | cada pasta em `skills/` com um `SKILL.md` | o grupo para o qual `skills/` aponta, ou `account` |
| Comandos, agentes | arquivos `.md` em `commands/` e `agents/`, subpastas como `folder/name` | o grupo, ou `account` |
| Plugins | `plugins/installed_plugins.json`, ativados ou desativados por `enabledPlugins` nas configurações da conta e do projeto; um plugin instalado só para um projeto aparece só nesse projeto | o grupo, ou `account`; `(off)` quando não está ativado |
| Servidores MCP | o `mcp.json` do grupo (quando o MCP é compartilhado), o `.claude.json` da conta (os servidores próprios dela e os que ela tem só para esta pasta) e os servidores dos plugins dela, como `plugin:<plugin>:<server>` | o grupo, `account` ou `local` |
| Hooks | `hooks` em `settings.json`, contados por evento | o grupo, ou `account` |
| Este projeto | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` e `.mcp.json` na pasta atual | listados uma vez; uma conta que não confiou na pasta recebe um problema pelos servidores do `.mcp.json` |

## Problemas que ele aponta

- **Links quebrados**: um link em `skills/`, `commands/` ou `agents/` cujo destino sumiu, ou uma pasta compartilhada (`skills`, `commands`, `agents`, `plugins`) que aponta para um lugar que não existe (repare com `ccam doctor --fix`, que cria de novo uma pasta de grupo ausente).
- **Uma pasta de skill sem `SKILL.md`**: o Claude Code a ignora.
- **Um servidor MCP definido de forma diferente em dois lugares**: o `mcp.json` do grupo e o `.claude.json` da conta, ou um deles e o `.mcp.json` do projeto. As mesmas configurações nos dois lugares (o que `mcp_strategy = "sync"` grava) não são um problema. Mova um servidor próprio da conta para o grupo com `ccam mcp import <account> <server>`.
- **Um `installed_plugins.json` que não é JSON válido**, um plugin que `enabledPlugins` não menciona e um plugin cuja pasta está ausente.
- **Servidores que só uma conta carrega**: um servidor só no `.claude.json` desta conta, que as outras contas do grupo dela não carregam, e um no escopo local da conta para esta pasta (o que `claude mcp add` faz, a menos que receba `-s user`). Cada um vem com o comando `ccam mcp import` que o compartilha. Servidores no `~/.claude.json` do `claude` simples que nenhuma conta carrega vêm com o comando `ccam migrate --mcp`.
- **Um servidor de projeto esperando aprovação**: a conta ainda não confiou na pasta, então o Claude Code pergunta uma vez e, até lá, os servidores do `.mcp.json` do projeto ficam esperando. O problema diz onde aprová-los: `"enabledMcpjsonServers"` no `settings.json` das contas do grupo de configurações, ou desta conta. A confiança dada à worktree principal de um repositório git vale nas outras worktrees dele.
- **Logins e conectores**: um servidor MCP que precisa de login nesta conta enquanto outra conta tem login nele (um login vale só para uma conta; entre com `/mcp`), e conectores do claude.ai que outra conta usou e esta não.

## Segurança e casos especiais

- `ccam ext` não altera nada e não faz solicitações de rede. Ele lê nomes de servidores e contagens de hooks, nunca os valores de `env` ou de cabeçalhos nas definições de MCP. Para a verificação de login, ele lê os nomes dos logins de MCP no login de cada conta (no macOS, do Acesso às Chaves), nunca os tokens, e só para contas que carregam um servidor que outra conta também carrega.
- Ele mostra o que está no disco. Um servidor bloqueado por configurações gerenciadas continua aparecendo; um plugin desativado nas configurações mostra `(off)`.
- Execute-o na pasta do projeto para incluir a pasta `.claude` e o `.mcp.json` desse projeto.

## Relacionados

- [Compartilhamento entre contas](/features/share)
- [Servidores MCP compartilhados](/features/mcp)
- [Verificar e reparar](/features/doctor)
