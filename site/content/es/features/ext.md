---
title: Inspeccionar extensiones
description: "ccam ext muestra los skills, comandos, plugins, servidores MCP y hooks que carga cada cuenta, de dónde viene cada uno y qué está roto."
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` responde a "¿por qué este skill (o servidor MCP) funciona en una cuenta y no en la otra?". Lista, para cada cuenta, las extensiones que carga Claude Code y de dónde viene cada una: un grupo de uso compartido, la carpeta propia de la cuenta o el proyecto en el que estás. Debajo de la tabla lista lo que impide que una extensión se cargue como esperas. Solo lee.

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

Una celda indica el origen de ese elemento en esa cuenta: el grupo de uso compartido (`default`), `account` cuando está en la carpeta propia de la cuenta, o `-` cuando la cuenta no lo tiene. `!` marca un elemento con un problema, que se explica debajo de la tabla.

## Qué lee

| Tipo | Dónde | Origen mostrado |
| --- | --- | --- |
| Skills | cada carpeta de `skills/` con un `SKILL.md` | el grupo al que enlaza `skills/`, o `account` |
| Comandos, agentes | archivos `.md` dentro de `commands/` y `agents/`, las subcarpetas como `folder/name` | el grupo, o `account` |
| Plugins | `plugins/installed_plugins.json` | el grupo, o `account` |
| Servidores MCP | el `mcp.json` del grupo (cuando MCP se comparte) y el `.claude.json` de la cuenta | el grupo, o `account` |
| Hooks | `hooks` en `settings.json`, contados por evento | el grupo, o `account` |
| Este proyecto | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` y `.mcp.json` en la carpeta actual | se lista una vez, para cada cuenta |

## Problemas que informa

- **Enlaces rotos**: un enlace en `skills/`, `commands/` o `agents/` cuyo destino ya no existe, o una carpeta compartida (`skills`, `commands`, `agents`, `plugins`) que enlaza a un lugar que falta (repárala con `ccam doctor --fix`).
- **Una carpeta de skill sin `SKILL.md`**: Claude Code la omite.
- **Un servidor MCP definido de forma distinta en dos lugares**: el `mcp.json` del grupo y el `.claude.json` de la cuenta, o uno de ellos y el `.mcp.json` del proyecto. La misma configuración en ambos lugares (lo que escribe `mcp_strategy = "sync"`) no es un problema. Mueve un servidor propio de una cuenta al grupo con `ccam mcp import <account> <server>`.
- **Un `installed_plugins.json` que no es JSON válido.**

## Seguridad y casos límite

- `ccam ext` no cambia nada ni hace solicitudes de red. Lee los nombres de los servidores y el número de hooks, nunca los valores de `env` ni de las cabeceras en las definiciones MCP.
- Muestra lo que hay en el disco. Un plugin instalado pero desactivado en la configuración, o un servidor bloqueado por una configuración administrada, sigue apareciendo.
- Ejecútalo en la carpeta del proyecto para incluir la carpeta `.claude` y el `.mcp.json` de ese proyecto.

## Relacionado

- [Compartir entre cuentas](/features/share)
- [Servidores MCP compartidos](/features/mcp)
- [Comprobar y reparar](/features/doctor)
