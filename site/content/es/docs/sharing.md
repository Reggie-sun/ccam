---
title: Compartir servidores MCP, configuración e historial
description: Cómo ccam comparte servidores MCP, settings.json, CLAUDE.md, agentes, skills e historial entre cuentas de Claude Code y mantiene separado cada inicio de sesión.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

Cada cuenta mantiene su propio inicio de sesión. Todo lo demás se puede compartir elemento por elemento: ccam enlaza el archivo o la carpeta de la cuenta con una carpeta de grupo compartida, de modo que un cambio hecho en una cuenta aparece en todas las cuentas de ese grupo.

## Qué se puede compartir

| Elemento | Opción | Qué se comparte |
| --- | --- | --- |
| Servidores MCP | `--mcp` | El `mcp.json` del grupo, que se pasa a Claude Code al iniciar |
| Configuración | `--settings` | `settings.json`: permisos, hooks, modelo, línea de estado |
| Memoria global | `--memory` | `CLAUDE.md` |
| Extensiones | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Historial de conversaciones | `--history` | `projects/` y `history.jsonl` |

Las cuentas nuevas comparten todo excepto el historial de conversaciones.

Los archivos de inicio de sesión `.credentials.json` y `.claude.json` nunca se comparten. ccam solo enlaza rutas de una lista fija, y estos dos no están en ella.

## Grupos

`on` significa el grupo `default`. También puedes indicar un nombre, para que solo algunas cuentas compartan un elemento:

```sh
ccam share work --history team
ccam share side --history team
```

Ahora `work` y `side` comparten conversaciones entre sí, pero no con tus otras cuentas.

## Activar el uso compartido

Cuando una cuenta empieza a compartir un elemento, ccam:

1. imprime un plan con cada archivo que va a mover, copiar o enlazar, y te pregunta antes de empezar;
2. mueve los archivos de la cuenta a la carpeta del grupo si el grupo todavía está vacío, o agrega los archivos que el grupo aún no tiene;
3. cuando el mismo archivo es distinto, te pregunta si quieres conservar la versión del grupo, usar la tuya o combinar el JSON;
4. guarda tu original en `.ccam-backup/<time>/` dentro de la carpeta de la cuenta;
5. reemplaza el original por un enlace al grupo.

Los archivos de conversación llevan el nombre de su ID de sesión, así que los historiales de varias cuentas se combinan sin conflictos.

## Desactivar el uso compartido

```sh
ccam share work --history off
```

Se quita el enlace y la cuenta recibe su propia copia del contenido actual del grupo. Agrega `--fresh` para empezar vacía. El grupo conserva su contenido para las demás cuentas.

## Servidores MCP

Claude Code guarda los servidores MCP en `.claude.json`, junto con el inicio de sesión, así que ese archivo no se puede compartir. ccam guarda los servidores compartidos en el `mcp.json` del grupo e inicia Claude Code con `--mcp-config=<group>/mcp.json`. Tus propias opciones `--mcp-config` se siguen aplicando encima.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Un servidor que agregas dentro de una sesión con `claude mcp add` se guarda solo para esa cuenta. La próxima vez que inicies esa cuenta, ccam te indica el nombre del servidor y el comando que lo mueve al grupo:

```sh
ccam mcp import work gh
```

Los servidores MCP que inician sesión con OAuth guardan sus tokens junto con el inicio de sesión de la cuenta, así que autorízalos una vez en cada cuenta.

Si prefieres que los servidores aparezcan en `claude mcp list` y `/mcp` como servidores de usuario normales, define `mcp_strategy = "sync"` en `~/.config/ccam/config.toml`. Entonces ccam los escribe en el `.claude.json` de cada cuenta antes de iniciar, con el bloqueo de archivo propio de Claude Code y con una copia de seguridad.

## Dos cuentas en la misma carpeta

Las cuentas que comparten el historial de conversaciones escriben en los mismos archivos. Si inicias un segundo Claude Code en una carpeta donde ya se está ejecutando otra cuenta del mismo grupo, ccam te avisa y te pregunta si quieres continuar.

## Cuando algo falla

```sh
ccam doctor
ccam doctor --fix
```

Si una ruta compartida se convirtió en un archivo normal, por ejemplo después de restaurar una copia de seguridad a mano, ccam se detiene antes de iniciar Claude Code y te remite a `ccam doctor --fix`, que combina el archivo con el grupo y guarda una copia de seguridad.
