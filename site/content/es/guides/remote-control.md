---
title: Usar Remote Control con ccam
description: "Usa Remote Control de Claude Code con cualquier cuenta mediante ccam <name> rc: opciones válidas, por qué ccam ofrece copiar los servidores MCP y soluciones."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 5f0b90ea53d6
---

`ccam work rc` inicia el Remote Control de Claude Code con la cuenta `work`, para que puedas manejar sesiones de este equipo desde claude.ai/code o desde la app móvil de Claude. Funciona como `claude rc` con una diferencia: si la cuenta comparte servidores MCP, ccam pregunta una vez si quieres copiarlos en el `.claude.json` de la cuenta, porque las sesiones de Remote Control no leen los servidores de ningún otro sitio.

## Iniciarlo

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

Un [atajo](/features/shortcuts) funciona igual: `ccwork rc` equivale a `ccam work rc`. Ejecuta el comando en la carpeta en la que quieres trabajar y después abre claude.ai/code o la app de Claude, con la sesión iniciada en la misma cuenta de Claude que `work`.

## Opciones

Aquí se juntan dos tipos de opciones: las de ccam, que van antes del nombre de la cuenta o están en la configuración de ccam, y las de Remote Control, que van después de `rc` y llegan a Claude Code sin cambios.

### Opciones y configuración de ccam

| Opción o ajuste | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--dry-run` (antes de la cuenta) | Imprime el entorno y el comando, y el aviso de MCP si lo hay. No pregunta ni copia nada | desactivada |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Se pregunta cuando la cuenta comparte servidores MCP y a su `.claude.json` le falta la lista actual del grupo. `y` los copia; `N` o Enter inicia sin ellos | `N` |
| `mcp_strategy = "sync"` en `~/.config/ccam/config.toml` | Copia los servidores del grupo en `.claude.json` antes de cada inicio, sea de Remote Control o no, así que nunca se pregunta nada | `"flag"` |
| `ccam group leave <group> work` | Saca la cuenta del grupo que comparte MCP con ella; Remote Control se inicia entonces solo con los servidores propios de la cuenta y no se pregunta nada | en `default` |

### Opciones de Remote Control

Todo lo que va después de `rc` pertenece a Claude Code. Estas son las opciones de `claude remote-control` en Claude Code 2.1.304; `ccam work rc --help` imprime la lista de tu versión.

| Opción | Qué hace |
| --- | --- |
| `--name <name>` | Nombre de la sesión que se muestra en claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Prefijo de los nombres de sesión generados (predeterminado: el nombre del equipo) |
| `-c`, `--continue` | Vuelve a conectarse a la última sesión que Remote Control registró para esta carpeta, dentro de unas 4 horas |
| `--session-id <id>` | Vuelve a conectarse a una sesión concreta |
| `--permission-mode <mode>` | Modo de permisos de las sesiones nuevas: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` o `plan` |
| `--spawn <mode>` | `same-dir` (predeterminado), `worktree` (un git worktree por sesión) o `session` (una sola sesión; sale cuando termina) |
| `--capacity <N>` | Máximo de sesiones simultáneas en los modos `same-dir` y `worktree` (predeterminado: 32) |
| `--[no-]create-session-in-dir` | Crea una primera sesión en la carpeta actual al iniciar (activado de forma predeterminada) |
| `--[no-]chrome` | Claude in Chrome para las sesiones nuevas (desactivado de forma predeterminada) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Salida de depuración y detallada |

Pon `rc` justo después del nombre de la cuenta. ccam toma como comando de Claude Code el primer argumento que no empieza por `-`, así que en `ccam work rc --name laptop` ve `rc`, pero en `ccam work --model opus rc` ve `opus`.

## Qué hace ccam

El inicio sigue la ruta habitual ([Iniciar Claude Code](/features/launch)) con un cambio en el paso de MCP:

1. ccam busca la cuenta, comprueba sus enlaces compartidos y resuelve `claude`, como en cualquier inicio.
2. Ve que el comando es `rc` o `remote-control` y omite `--mcp-config`.
3. Si la cuenta comparte servidores MCP, compara el `mcp.json` del grupo con el `sync-stamp` de `~/.local/state/ccam/mcp/work/`. Si `.claude.json` ya contiene los servidores actuales, o el grupo no tiene ninguno, pasa directamente al paso 6.
4. Si no, nombra los servidores que les faltarían a las sesiones:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. Con `y`, guarda una copia de seguridad de `.claude.json`, toma el bloqueo propio de Claude Code, añade los servidores del grupo a `mcpServers` y conserva los demás servidores de la cuenta. Con `N`, sin una terminal o mientras la cuenta está en ejecución, no cambia nada y lo indica:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Registra la sesión, establece `CLAUDE_CONFIG_DIR` en la carpeta de la cuenta y se reemplaza a sí mismo por `claude rc` mediante `exec`.

A partir de ahí, Remote Control sigue ejecutándose como servidor. Cada sesión que inicia para ti es un proceso hijo de ese proceso, así que se ejecuta como `work`: con el inicio de sesión, la configuración y el historial de conversaciones de `work`, y cuenta para el límite de uso de `work`. Mientras el servidor está en ejecución, ccam considera que la cuenta está en ejecución.

La ejecución de prueba muestra la diferencia. Un inicio normal lleva los servidores del grupo como argumento; Remote Control, no:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Por qué funciona así

**Remote Control no acepta `--mcp-config`.** En los inicios normales, ccam entrega los servidores del grupo a Claude Code como `--mcp-config=<path>` y nunca escribe en los archivos de Claude Code. `claude remote-control` rechaza esa opción antes de su nombre de comando y no acepta ninguna después, así que antes de ccam 0.6.0, `ccam work rc` fallaba en todas las cuentas que compartían servidores MCP. Ahora ccam omite la opción.

**Sus sesiones leen los servidores solo de `.claude.json`.** Las sesiones que inicia Remote Control no heredan las opciones de la línea de comandos. El único lugar donde encuentran servidores MCP de usuario es el `.claude.json` de la cuenta. Ese archivo también contiene los datos de inicio de sesión de la cuenta y el estado por carpeta, así que ccam no puede enlazarlo a una copia compartida como hace con `settings.json`. Copiar los servidores en él es la única forma de dar a las sesiones la lista compartida.

**ccam pregunta antes de escribir.** Con la estrategia de opción predeterminada, ccam se compromete a no tocar `.claude.json`. Remote Control es el único caso en el que tiene que escribir ahí, así que pregunta primero, guarda una copia de seguridad y no toca nada más que `mcpServers`.

**No escribe mientras la cuenta está en ejecución.** Claude Code reescribe `.claude.json` durante una sesión, y un cambio hecho desde fuera se perdería o se sobrescribiría. Por eso ccam omite la copia cuando hay cualquier sesión de la cuenta en ejecución, incluido otro servidor de Remote Control.

**Pregunta una vez por cambio.** Después de una copia, `sync-stamp` guarda la fecha de modificación del `mcp.json` del grupo. ccam vuelve a preguntar solo cuando ese archivo cambia, por ejemplo después de `ccam mcp add`. Si respondes `N`, todavía no hay ninguna copia, así que pregunta en cada inicio de Remote Control.

**`ccam mcp` mantiene las copias al día.** Una entrada de `.claude.json` gana sobre `--mcp-config` cuando los nombres coinciden, así que una copia antigua ocultaría una definición más reciente también en los inicios normales. Por eso `ccam mcp add`, `add-json`, `remove` e `import` actualizan las copias de todas las cuentas que no están en ejecución:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Sin una terminal, la respuesta es no.** Remote Control se ejecuta a menudo desde un script, un servicio o `nohup`, donde nadie puede responder. En ese caso, ccam inicia sin los servidores en lugar de quedarse esperando.

## Problemas y soluciones

### Faltan los servidores MCP compartidos en las sesiones de Remote Control

Se respondió `N` a la pregunta, no había terminal o la cuenta estaba en ejecución. Detén Remote Control, cierra las demás sesiones de la cuenta y vuelve a ejecutar `ccam work rc` en una terminal; responde `y`. Para que no se pregunte nunca, establece `mcp_strategy = "sync"` en `~/.config/ccam/config.toml`. Los servidores de una sesión no cambian mientras se ejecuta, así que inicia una sesión nueva después de la copia.

### "Account work is running, so its .claude.json is left alone"

Sigue abierta una sesión de `work`: un `ccam work` normal, un atajo u otro servidor de Remote Control. Ciérrala y vuelve a ejecutar `ccam work rc`. Las sesiones que ejecuta otra cuenta no cuentan.

### ccam pregunta cada vez

Todavía no has copiado los servidores. Responde `y` una vez o establece `mcp_strategy = "sync"`. Si no quieres nunca los servidores en Remote Control, deja de compartir MCP con la cuenta: sácala de ese grupo (`ccam group leave default work`) o mantén MCP en un grupo propio en el que la cuenta no esté.

### `ccam work rc` falla de inmediato con un error sobre `--mcp-config`

Tu ccam es anterior a la 0.6.0 y pasa `--mcp-config` a Remote Control. Ejecuta `ccam update`. Si ves el mismo error con un ccam actual, una opción con valor iba antes de `rc`, como en `ccam work --model opus rc`, y ccam tomó `opus` por el comando; pon `rc` justo después del nombre de la cuenta.

### Un cambio hecho con `ccam mcp` no llega a Remote Control

`ccam mcp add` y `remove` actualizan las copias solo en las cuentas que no están en ejecución, y lo indican: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Detén Remote Control y vuelve a iniciarlo; ccam ve que `mcp.json` cambió y vuelve a preguntar.

### La copia reemplazó un servidor propio

Si la cuenta tenía un servidor propio con el mismo nombre que un servidor del grupo, la copia lo reemplazó. El `.claude.json` anterior está en `~/.local/state/ccam/mcp/work/backups/`, que conserva los cinco más recientes. Para conservar tu versión en todas las cuentas, muévela al grupo con `ccam mcp import work <server> --force` antes del siguiente inicio de Remote Control.

### Quiero quitar los servidores copiados

`ccam mcp remove <server>` quita un servidor del grupo y de todas las copias. Para volver al estado anterior a la copia, restaura la copia de seguridad más reciente de `~/.local/state/ccam/mcp/work/backups/` mientras no haya ninguna sesión de la cuenta en ejecución.

### Remote Control pregunta si confiar en la carpeta, o sale sin una terminal

La confianza en las carpetas se guarda por cuenta en `.claude.json`, así que una carpeta en la que confiaste como `personal` es nueva para `work`. Ejecuta `ccam work` una vez en esa carpeta y acepta; después inicia Remote Control ahí.

### "You must be logged in with a Claude account that has a subscription"

Remote Control necesita un inicio de sesión con una suscripción de Claude. Comprueba el correo electrónico de la cuenta en `ccam list` y, si hace falta, vuelve a iniciar sesión con `ccam login work`.

### La sesión no aparece en el teléfono

La app de Claude o claude.ai/code tiene la sesión iniciada con otra cuenta de Claude. Compárala con el correo electrónico que `ccam list` muestra para `work`.

### ccam advierte que otro Claude Code comparte el historial de conversaciones en esta carpeta

Otra cuenta del mismo grupo de historial se está ejecutando en la misma carpeta. Las sesiones que continúan la misma conversación desde ambas cuentas pueden mezclar el historial. Responde `y` solo si van a trabajar en conversaciones distintas, o inicia Remote Control en otra carpeta.

## Relacionado

- [Iniciar Claude Code](/features/launch)
- [Servidores MCP compartidos](/features/mcp)
- [Compartir entre cuentas](/features/share)
- [Administrar cuentas](/features/accounts)
