---
title: Compartir servidores MCP, configuración e historial
description: Cómo los grupos de ccam comparten servidores MCP, settings.json, CLAUDE.md, skills e historial entre cuentas de Claude Code, con cada inicio de sesión aparte.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: 27e0a0e5abbd
---

Cada cuenta mantiene su propio inicio de sesión. Todo lo demás se puede compartir mediante **grupos**: un grupo nombra los elementos que comparte y las cuentas que contiene, y ccam enlaza los archivos de cada una de esas cuentas con la carpeta del grupo, de modo que un cambio hecho en una cuenta aparece en todas.

## Qué se puede compartir

| Elemento | Qué se comparte |
| --- | --- |
| `mcp` | Servidores MCP: el `mcp.json` del grupo, que se pasa a Claude Code al iniciar |
| `settings` | `settings.json`: permisos, hooks, modelo, línea de estado |
| `memory` | el `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | historial de conversaciones: `projects/` y `history.jsonl` |

Tu primera cuenta crea el grupo `default`, que comparte todo excepto el historial de conversaciones, y las cuentas posteriores se unen a él.

Los archivos de inicio de sesión `.credentials.json` y `.claude.json` nunca se comparten. ccam solo enlaza rutas de una lista fija, y estos dos no están en ella.

Esta página cubre las tareas cotidianas. Todas las opciones se describen en las páginas de funciones: [compartir con grupos](/features/share), [servidores MCP compartidos](/features/mcp) y [copiar una configuración existente](/features/migrate).

## Grupos

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Una cuenta puede estar en varios grupos, pero cada elemento proviene de un solo grupo. Así que, para compartir el historial entre `work` y `side` pero no con tus otras cuentas, deja el historial fuera de `default` y dales a las dos un grupo propio, como arriba. `ccam share work` muestra lo que comparte una cuenta y a través de qué grupo. El [panel](/features/dashboard) hace lo mismo en su pestaña Groups.

## Cuando una cuenta empieza a compartir

Cuando una cuenta empieza a compartir un elemento (se une a un grupo, o su grupo agrega el elemento), ccam:

1. imprime un plan con cada archivo que va a mover, copiar o enlazar en cada cuenta, y te pregunta antes de empezar;
2. mueve los archivos de la cuenta a la carpeta del grupo si el grupo todavía está vacío, o agrega los archivos que el grupo aún no tiene;
3. cuando el mismo archivo es distinto, te pregunta si quieres conservar la versión del grupo, usar la tuya o combinar el JSON;
4. guarda tu original en `.ccam-backup/<time>/` dentro de la carpeta de la cuenta;
5. reemplaza el original por un enlace al grupo.

Los archivos de conversación llevan el nombre de su ID de sesión, así que los historiales de varias cuentas se combinan sin conflictos.

ccam no cambia una cuenta mientras Claude Code se está ejecutando con ella, sin importar cómo se haya iniciado Claude Code (también las sesiones en segundo plano), porque la sesión sigue escribiendo en los archivos que se mueven. Ciérrala primero. ccam vuelve a comprobarlo justo antes de guardar, así que también detecta una sesión iniciada mientras el plan esperaba.

## Cuando una cuenta deja de compartir

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Se quita el enlace y la cuenta recibe su propia copia del contenido actual del grupo; los servidores MCP del grupo van a su propio `.claude.json`, junto a sus propios servidores. Agrega `--fresh` para empezar vacía; el plan entonces enumera los servidores MCP que dejan de cargarse. Una carpeta importada como `~/.claude` siempre conserva una copia. El grupo conserva su contenido para las demás cuentas.

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
- el historial de conversaciones (`projects/` y `history.jsonl`), pero solo si agregas `--history`.

La carpeta de origen solo se lee: nada en ella se mueve, se enlaza ni se modifica, y `claude` a secas sigue funcionando como antes. Los archivos de inicio de sesión `.credentials.json` y `.claude.json` nunca se copian al grupo.

Los servidores MCP que inician sesión con OAuth guardan sus tokens en el inicio de sesión, no en la lista de servidores. ccam no copia esos inicios de sesión: un inicio de sesión copiado lleva el mismo token de actualización, y la primera renovación cierra la sesión de ese servidor para todos los demás que lo tienen, `~/.claude` incluido. El plan y el resumen nombran las cuentas que inician sesión en los servidores por sí mismas con `/mcp`.

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
Group default: personal
To share them:
  ccam group join default work
```

Al final, ccam lista las cuentas del grupo. Si el grupo todavía no comparte alguno de los elementos, o algunas cuentas no están en él, imprime los comandos `ccam group` que lo cambian.

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

Un servidor que agregas dentro de una sesión con `claude mcp add` se guarda solo para esa cuenta y para la carpeta en la que estabas (el ámbito local de Claude Code; `-s user` lo hace de la cuenta en todas las carpetas). La próxima vez que inicies esa cuenta, ccam te indica el nombre del servidor, con su carpeta, y el comando que lo mueve al grupo:

```sh
ccam mcp import work gh
```

Si el grupo ya tiene un servidor distinto con ese nombre, ccam lo omite y te lo indica. `--force` reemplaza el servidor del grupo y guarda el `mcp.json` anterior en la carpeta `.ccam-backup/` del grupo.

Los servidores MCP que inician sesión con OAuth guardan sus tokens junto con el inicio de sesión de la cuenta, así que autorízalos una vez en cada cuenta con `/mcp`. `ccam ext` y `ccam doctor` nombran los servidores en los que una cuenta no ha iniciado sesión mientras otra cuenta sí.

Si prefieres que los servidores aparezcan en `claude mcp list` y `/mcp` como servidores de usuario normales, define `mcp_strategy = "sync"` al principio de `~/.config/ccam/config.toml`, antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`. Entonces ccam los escribe en el `.claude.json` de cada cuenta antes de iniciar, con el bloqueo de archivo propio de Claude Code y con una copia de seguridad. Un servidor propio de la cuenta con el mismo nombre sigue siendo de la cuenta, y ccam lo indica.

## Dos cuentas en la misma carpeta

Las cuentas que comparten el historial de conversaciones escriben en los mismos archivos. Si inicias un segundo Claude Code en una carpeta donde ya se está ejecutando otra cuenta del mismo grupo, ccam te avisa y te pregunta si quieres continuar.

## Cuando algo falla

```sh
ccam doctor
ccam doctor --fix
```

Si una ruta compartida se convirtió en un archivo normal, por ejemplo después de restaurar una copia de seguridad a mano, ccam se detiene antes de iniciar Claude Code y te remite a `ccam doctor --fix`, que vuelve a enlazarlo al grupo: el archivo pasa a ser la copia del grupo cuando el grupo no tiene ninguna; si no, el grupo conserva la suya y la tuya va a una copia de seguridad, que `--fix` nombra. Un cambio de grupos que se quedó a medias también lo termina `ccam doctor --fix`.
