---
title: Inspecionar extensões
description: ccam ext mostra as skills, comandos, plugins, servidores MCP e hooks que cada conta carrega, de onde vem cada um e o que está quebrado.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

Cada célula indica a origem daquele item naquela conta: o grupo de compartilhamento (`default`), `account` quando ele fica na pasta da própria conta, ou `-` quando a conta não o tem. `!` marca um item com problema, explicado abaixo da tabela.

## O que ele lê

| Tipo | Onde | Origem mostrada |
| --- | --- | --- |
| Skills | cada pasta em `skills/` com um `SKILL.md` | o grupo para o qual `skills/` aponta, ou `account` |
| Comandos, agentes | arquivos `.md` em `commands/` e `agents/`, subpastas como `folder/name` | o grupo, ou `account` |
| Plugins | `plugins/installed_plugins.json` | o grupo, ou `account` |
| Servidores MCP | o `mcp.json` do grupo (quando o MCP é compartilhado) e o `.claude.json` da conta | o grupo, ou `account` |
| Hooks | `hooks` em `settings.json`, contados por evento | o grupo, ou `account` |
| Este projeto | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` e `.mcp.json` na pasta atual | listados uma vez, para todas as contas |

## Problemas que ele aponta

- **Links quebrados**: um link em `skills/`, `commands/` ou `agents/` cujo destino sumiu, ou uma pasta compartilhada (`skills`, `commands`, `agents`, `plugins`) que aponta para um lugar que não existe (repare com `ccam doctor --fix`).
- **Uma pasta de skill sem `SKILL.md`**: o Claude Code a ignora.
- **Um servidor MCP definido de forma diferente em dois lugares**: o `mcp.json` do grupo e o `.claude.json` da conta, ou um deles e o `.mcp.json` do projeto. As mesmas configurações nos dois lugares (o que `mcp_strategy = "sync"` grava) não são um problema. Mova um servidor próprio da conta para o grupo com `ccam mcp import <account> <server>`.
- **Um `installed_plugins.json` que não é JSON válido.**

## Segurança e casos especiais

- `ccam ext` não altera nada e não faz solicitações de rede. Ele lê nomes de servidores e contagens de hooks, nunca os valores de `env` ou de cabeçalhos nas definições de MCP.
- Ele mostra o que está no disco. Um plugin instalado mas desativado nas configurações, ou um servidor bloqueado por configurações gerenciadas, continua aparecendo.
- Execute-o na pasta do projeto para incluir a pasta `.claude` e o `.mcp.json` desse projeto.

## Relacionados

- [Compartilhamento entre contas](/features/share)
- [Servidores MCP compartilhados](/features/mcp)
- [Verificar e reparar](/features/doctor)
