---
title: Inspeccionar extensiones
description: "ccam ext muestra los skills, comandos, plugins, servidores MCP y hooks que carga cada cuenta, de dónde viene cada uno y qué está roto."
updated: 2026-10-26
order: 14
translation: machine
sourceHash: "7ca357a7f38a"
---

`ccam ext` responde a la pregunta "¿Por qué este skill (o servidor MCP) funciona en una cuenta pero no en la otra?". Lista, para cada cuenta, las extensiones que carga Claude Code y de dónde viene cada una: un grupo de uso compartido, la carpeta propia de la cuenta o el proyecto en el que estás. Debajo de la tabla lista lo que impide que una extensión se cargue como esperas. No cambia nada.

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

Una celda indica el origen de ese elemento en esa cuenta: el grupo de uso compartido (`default`), `account` cuando está en la carpeta propia de la cuenta, `local` para un servidor MCP que la cuenta tiene solo para esta carpeta, o `-` cuando la cuenta no lo tiene. Un plugin instalado pero no activado muestra `(off)` después de su origen. `!` marca un elemento con un problema, que se explica debajo de la tabla; los problemas que no pertenecen a ninguna cuenta, como los servidores que solo tiene `claude` a secas, van al final sin nombre. `ccam ext --json` imprime el mismo informe (`accounts`, el proyecto, `problems`), con `remote` en los servidores MCP a los que Claude Code llega por la red; sin cuentas, imprime el informe con `accounts: []`.

## Qué lee

| Tipo | Dónde | Origen mostrado |
| --- | --- | --- |
| Skills | cada carpeta de `skills/` con un `SKILL.md` | el grupo al que enlaza `skills/`, o `account` |
| Comandos, agentes | archivos `.md` dentro de `commands/` y `agents/`, las subcarpetas como `folder/name` | el grupo, o `account` |
| Plugins | `plugins/installed_plugins.json`, activados o desactivados por `enabledPlugins` en la configuración de la cuenta y del proyecto; un plugin instalado solo para un proyecto aparece solo en ese proyecto | el grupo, o `account`; `(off)` cuando no está activado |
| Servidores MCP | el `mcp.json` del grupo (cuando MCP se comparte), el `.claude.json` de la cuenta (sus propios servidores y los que tiene solo para esta carpeta) y los servidores de sus plugins, como `plugin:<plugin>:<server>` | el grupo, `account` o `local` |
| Hooks | `hooks` en `settings.json`, contados por evento | el grupo, o `account` |
| Este proyecto | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` y `.mcp.json` en la carpeta actual | se lista una vez; una cuenta que no ha confiado en la carpeta recibe un problema por los servidores de su `.mcp.json` |

## Problemas que informa

- **Enlaces rotos**: un enlace en `skills/`, `commands/` o `agents/` cuyo destino ya no existe, o una carpeta compartida (`skills`, `commands`, `agents`, `plugins`) que enlaza a un lugar que falta (repárala con `ccam doctor --fix`, que vuelve a crear una carpeta de grupo que falta).
- **Una carpeta de skill sin `SKILL.md`**: Claude Code la omite.
- **Un servidor MCP definido de forma distinta en dos lugares**: el `mcp.json` del grupo y el `.claude.json` de la cuenta, o uno de ellos y el `.mcp.json` del proyecto. Las definiciones idénticas en ambos lugares (lo que escribe `mcp_strategy = "sync"`) no son un problema. Mueve un servidor propio de una cuenta al grupo con `ccam mcp import <account> <server>`.
- **Un `installed_plugins.json` que no es JSON válido**, un plugin que `enabledPlugins` no nombra y un plugin cuya carpeta falta.
- **Servidores que solo carga una cuenta**: un servidor que solo está en el `.claude.json` de esta cuenta, que las demás cuentas de su grupo no cargan, y uno en el ámbito local de la cuenta para esta carpeta (lo que hace `claude mcp add` salvo que se le pase `-s user`). Cada uno viene con el comando `ccam mcp import` que lo comparte. Los servidores del `~/.claude.json` de `claude` a secas que ninguna cuenta carga vienen con el comando `ccam migrate --mcp`.
- **Un servidor del proyecto pendiente de aprobación**: la cuenta todavía no ha confiado en la carpeta, así que Claude Code pregunta una vez, y hasta entonces los servidores del `.mcp.json` del proyecto esperan. El problema indica dónde aprobarlos: `"enabledMcpjsonServers"` en el `settings.json` de las cuentas del grupo de configuración, o de esta cuenta. La confianza dada al worktree principal de un repositorio git también se aplica a sus otros worktrees.
- **Inicios de sesión y conectores**: un servidor MCP que necesita un inicio de sesión en esta cuenta mientras otra cuenta sí lo inició (un inicio de sesión vale solo para una cuenta; inicia sesión con `/mcp`), y los conectores de claude.ai que otra cuenta ha usado y esta no.

## Seguridad y casos límite

- `ccam ext` no cambia nada ni hace solicitudes de red. Lee los nombres de los servidores y el número de hooks, nunca los valores de `env` ni de las cabeceras en las definiciones MCP. Para la comprobación de inicios de sesión lee los nombres de los inicios de sesión MCP del inicio de sesión de cada cuenta (en macOS, desde el Llavero), nunca los tokens, y solo en las cuentas que cargan un servidor que también carga otra cuenta.
- Muestra lo que hay en el disco. Un servidor bloqueado por una configuración administrada sigue apareciendo; un plugin desactivado en la configuración muestra `(off)`.
- Ejecútalo en la carpeta del proyecto para incluir la carpeta `.claude` y el `.mcp.json` de ese proyecto.

## Relacionado

- [Compartir con grupos](/features/share)
- [Servidores MCP compartidos](/features/mcp)
- [Comprobar y reparar](/features/doctor)
