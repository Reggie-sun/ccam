---
title: Compartir servidores MCP, configuración e historial
description: Cómo ccam comparte servidores MCP, settings.json, CLAUDE.md, agentes, skills e historial entre cuentas de Claude Code y mantiene separado cada inicio de sesión.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 888d466b29e9
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

Esta página cubre las tareas cotidianas. Todas las opciones se describen en las páginas de funciones: [compartir entre cuentas](/features/share), [servidores MCP compartidos](/features/mcp) y [copiar una configuración existente](/features/migrate).

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

ccam no cambia el uso compartido mientras Claude Code se está ejecutando con la cuenta, porque la sesión sigue escribiendo en los archivos que se mueven. Ciérrala primero.

## Desactivar el uso compartido

```sh
ccam share work --history off
```

Se quita el enlace y la cuenta recibe su propia copia del contenido actual del grupo. Agrega `--fresh` para empezar vacía. El grupo conserva su contenido para las demás cuentas.

## Copiar una configuración existente

Si usabas Claude Code antes de ccam, tu configuración está en `~/.claude`. `ccam migrate` la copia a un grupo, para que todas las cuentas que comparten esos elementos la tengan:

```sh
ccam migrate
```

En una terminal, ccam primero muestra una lista de selección con los elementos. Todo está marcado excepto el historial de conversaciones; pulsa Espacio para cambiar un elemento y Enter para ver el plan. Los elementos para los que la carpeta no tiene nada aparecen como `not found`.

Qué copia:

- `settings.json` y el `CLAUDE.md` global;
- las extensiones `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` y `rules/`;
- tus servidores MCP de usuario, es decir, la entrada `mcpServers` de `~/.claude.json`, al `mcp.json` del grupo;
- los inicios de sesión guardados (OAuth) de los servidores MCP que copia, al archivo de inicio de sesión propio de cada cuenta;
- el historial de conversaciones (`projects/` y `history.jsonl`), pero solo si agregas `--history`.

La carpeta de origen solo se lee: nada en ella se mueve, se enlaza ni se modifica, y `claude` a secas sigue funcionando como antes. Los archivos de inicio de sesión `.credentials.json` y `.claude.json` nunca se copian al grupo.

Los servidores MCP que inician sesión con OAuth guardan sus tokens en el inicio de sesión, no en la lista de servidores. Para que no tengas que volver a iniciar sesión en ellos, ccam copia esos inicios de sesión al archivo de inicio de sesión de cada cuenta que recibe los servidores a través del grupo, salvo que la cuenta ya tenga los suyos. Omite las cuentas que están en ejecución, que aún no iniciaron sesión, que se registraron en su ubicación con `ccam import` o que guardan su inicio de sesión en el Llavero de macOS, y las nombra; en ellas, inicia sesión en los servidores con `/mcp`. Una cuenta que empieza a compartir servidores MCP más adelante recibe los inicios de sesión cuando vuelves a ejecutar `ccam migrate`.

ccam imprime el plan y pregunta antes de copiar nada:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Al final, ccam lista las cuentas que comparten cada elemento a través del grupo. Si algunas cuentas no lo hacen, imprime los comandos `ccam share` que activan el uso compartido para ellas.

Cuando el grupo ya tiene una copia distinta de un archivo, ccam pregunta, igual que al activar el uso compartido, si quieres conservar la copia del grupo o usar la tuya. `settings.json` y los demás archivos JSON también se pueden combinar, y `history.jsonl` línea por línea. Las copias anteriores del grupo se guardan en `.ccam-backup/<time>/` dentro de la carpeta del grupo.

| Opción | Efecto |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copiar solo los elementos que indicas, sin la lista de selección; sin ellas y fuera de una terminal, todo excepto el historial de conversaciones |
| `-g`, `--group <group>` | Copiar a este grupo en lugar de `default` |
| `-y`, `--yes` | Copiar sin preguntar; cuando un archivo es distinto, conservar la copia del grupo |

Pon `--dry-run` antes del comando, como en `ccam --dry-run migrate`, para ver el plan sin copiar nada. Para copiar desde otra carpeta, indícala: `ccam migrate ~/.claude-work --group work`.

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

Si el grupo ya tiene un servidor distinto con ese nombre, ccam lo omite y te lo indica. `--force` reemplaza el servidor del grupo y guarda el `mcp.json` anterior en la carpeta `.ccam-backup/` del grupo.

Los servidores MCP que inician sesión con OAuth guardan sus tokens junto con el inicio de sesión de la cuenta, así que autorízalos una vez en cada cuenta o copia los inicios de sesión de `~/.claude` con `ccam migrate`.

Si prefieres que los servidores aparezcan en `claude mcp list` y `/mcp` como servidores de usuario normales, define `mcp_strategy = "sync"` en `~/.config/ccam/config.toml`. Entonces ccam los escribe en el `.claude.json` de cada cuenta antes de iniciar, con el bloqueo de archivo propio de Claude Code y con una copia de seguridad.

## Dos cuentas en la misma carpeta

Las cuentas que comparten el historial de conversaciones escriben en los mismos archivos. Si inicias un segundo Claude Code en una carpeta donde ya se está ejecutando otra cuenta del mismo grupo, ccam te avisa y te pregunta si quieres continuar.

## Cuando algo falla

```sh
ccam doctor
ccam doctor --fix
```

Si una ruta compartida se convirtió en un archivo normal, por ejemplo después de restaurar una copia de seguridad a mano, ccam se detiene antes de iniciar Claude Code y te remite a `ccam doctor --fix`, que combina el archivo con el grupo y guarda una copia de seguridad.
