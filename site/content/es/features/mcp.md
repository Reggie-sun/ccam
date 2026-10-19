---
title: Servidores MCP compartidos
description: Mantén una sola lista de servidores MCP por grupo con ccam mcp add, add-json, list, remove e import, que se pasa con --mcp-config o se sincroniza en cada cuenta.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 69fd3233fd09
---

Claude Code guarda los servidores MCP de usuario en `.claude.json`, el mismo archivo que contiene el inicio de sesión de la cuenta, así que ese archivo no se puede compartir entre cuentas. ccam guarda los servidores compartidos en un archivo de grupo, `mcp.json`, y se lo entrega a cada cuenta del grupo al iniciar. Usa `ccam mcp` para agregar un servidor una sola vez para todas tus cuentas, y `ccam mcp import` para llevar a las demás un servidor que agregaste dentro de una sesión.

## Uso

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Un cambio se aplica a partir del siguiente inicio de cada cuenta; las sesiones en ejecución conservan los servidores con los que se iniciaron:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
```

Agregar un nombre que el grupo ya tiene reemplaza su definición e imprime `updated` en lugar de `added`. Los servidores llegan a una cuenta solo si comparte MCP con ese grupo (`ccam share <name> --mcp on`); cuando ninguna cuenta comparte MCP todavía, `add` lo indica.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `-g`, `--group <group>` | grupo de destino para `add`, `add-json` y `remove` | el único grupo a través del cual alguna cuenta comparte MCP; si no, `default` |
| `-- <command> [args…]` (`add`) | ejecuta este comando como servidor stdio | |
| `--url <url>` (`add`) | URL de un servidor remoto (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transporte de un servidor con `--url` | `http` |
| `--header <'Name: value'>` (`add`) | encabezado HTTP; se puede repetir | ninguno |
| `-e`, `--env <KEY=VALUE>` (`add`) | variable de entorno para un servidor stdio; se puede repetir | ninguno |
| `--force` (`import`) | reemplaza un servidor del grupo con el mismo nombre que es distinto; se hace una copia de seguridad del `mcp.json` anterior | desactivada |
| `--keep-local` (`import`) | conserva las entradas en el `.claude.json` de la cuenta | desactivada |
| `-h`, `--help` | ayuda de cada comando | |
| `--dry-run` (global) | solo la respeta `mcp list`; `add`, `add-json`, `remove` e `import` la rechazan y no cambian nada | desactivada |

`add` necesita un comando después de `--` o bien `--url`; si no, falla con `a command to run (after --) or --url is required`. `--header` se aplica a los servidores con `--url` y `--env` a los servidores con comando. `add-json` acepta el mismo JSON que `claude mcp add-json`, y la definición debe ser un objeto con `command` o `url`. `import` siempre usa el grupo a través del cual la cuenta comparte MCP, y `list` muestra todos los grupos, así que `-g` no los cambia.

## Cómo funciona

Los archivos que intervienen:

- `~/.local/share/ccam/groups/<group>/mcp.json`: los servidores compartidos, `{"mcpServers": {…}}`, el formato que lee `--mcp-config`. Se escribe con `fsutil.AtomicWrite` y el modo `0600`, porque las entradas de servidores suelen llevar claves de API.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: el archivo propio de Claude Code con los servidores de la cuenta. ccam lee solo su valor `mcpServers` con el escáner `profile.TopLevelValue` y nunca analiza el archivo completo en la ruta de inicio.
- `~/.local/state/ccam/mcp/<name>/`: el archivo `seen`, el `sync-stamp` y las copias de seguridad de `.claude.json`.

Hay dos estrategias, que se eligen con `mcp_strategy` en `~/.config/ccam/config.toml`. La predeterminada es la estrategia de opción:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` busca el grupo de MCP de la cuenta. `mcp.FlagArgs` comprueba con un solo `stat` que el `mcp.json` del grupo existe, sin analizarlo.
2. `Plan.InsertArgs` coloca `--mcp-config=<path>` justo después de `claude`, antes de tus propios argumentos, así que tus propias opciones `--mcp-config` o `--strict-mcp-config` se siguen aplicando encima.
3. `mcp.NewLocal` comprueba si el `.claude.json` de la cuenta cambió desde el último inicio, y solo entonces analiza sus servidores.
4. `syscall.Exec` reemplaza ccam por Claude Code.

La ruta se pasa como un solo argumento, `--mcp-config=<path>`, no como dos. `--mcp-config` de Claude Code acepta varios valores, así que `claude --mcp-config a.json "prompt"` leería la instrucción como un segundo archivo de configuración y fallaría con `MCP config file not found`. La forma de un solo argumento termina ahí, y una instrucción o un subcomando que venga después, como `mcp list`, funciona. Claude Code carga los servidores de `--mcp-config` además de los propios de la cuenta; si un nombre coincide, gana la entrada de `.claude.json`. La estrategia de opción solo escribe en `.claude.json` para Remote Control.

`claude remote-control` (`rc`) no acepta `--mcp-config`, y sus sesiones leen los servidores solo de `.claude.json`. Para `ccam work rc`, ccam omite la opción y, en una terminal, pregunta si quieres copiar los servidores del grupo en `.claude.json`, con una copia de seguridad, como hace la estrategia `sync`. Pregunta hasta que existe una copia y, después, solo cuando cambia el `mcp.json` del grupo ([Remote Control](/features/launch#remote-control)). Como esa copia gana sobre `--mcp-config` si un nombre coincide, `ccam mcp add`, `add-json`, `remove` e `import` actualizan las copias que tienen otras cuentas del grupo:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Un servidor agregado dentro de una sesión con `claude mcp add` queda solo en el `.claude.json` de esa cuenta. ccam no puede verlo en ese momento, porque Claude Code lo ha reemplazado, así que lo comprueba en el siguiente inicio. El archivo `seen` guarda la fecha de modificación de `.claude.json` y los nombres de servidores ya informados. Cuando la fecha es distinta, ccam compara los nombres con el grupo e imprime una sola vez un aviso por cada servidor nuevo, sin esperar ninguna respuesta:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

Con `mcp_strategy = "sync"`, los servidores aparecen en `claude mcp list` y `/mcp` como servidores de usuario normales, y ccam los escribe en el `.claude.json` de cada cuenta antes de iniciar:

<!-- motion:mcp-sync -->

1. `mcp.Sync` compara la fecha de modificación del `mcp.json` del grupo con `sync-stamp`. Si no ha cambiado, el `.claude.json` de la cuenta no se abre.
2. Toma el bloqueo propio de Claude Code, la carpeta `.claude.json.lock` junto al archivo, creada con `mkdir`. Espera hasta 2 segundos y considera obsoleto un bloqueo con más de 10 segundos, igual que la biblioteca de bloqueo de Claude Code.
3. Guarda el `.claude.json` actual en `backups/` dentro de la carpeta de estado y conserva los 5 más recientes.
4. Establece cada servidor del grupo en `mcpServers`, quita los servidores que escribió antes y que el grupo ya no tiene, conserva los servidores propios de la cuenta y escribe el archivo de forma atómica con su modo original.
5. Registra la nueva marca y los nombres que administra, y Claude Code se inicia sin `--mcp-config`.

La sincronización se omite, con `ccam: account work is already running, skipping MCP sync`, cuando hay una sesión de la cuenta en ejecución, porque Claude Code reescribe `.claude.json` mientras se ejecuta. Una sincronización fallida imprime `ccam: MCP sync failed (launching anyway): …` e inicia de todos modos. El aviso de servidores agregados dentro de una sesión es exclusivo de la estrategia de opción.

`ccam mcp import` mueve los servidores propios de una cuenta a su grupo:

<!-- motion:mcp-import -->

ccam lee los servidores de la cuenta con `mcp.LocalServers` y compara cada uno con el servidor del grupo que tiene el mismo nombre mediante `mcp.SameDef`, que ignora el orden de las claves y los espacios. Un servidor que el grupo no tiene, o que tiene con la misma definición, se copia. Uno distinto se omite, porque la definición del grupo es la que usan todas las demás cuentas:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Con `--force`, el servidor del grupo se reemplaza, y el `mcp.json` anterior se guarda primero en `.ccam-backup/<time>/mcp.json` dentro de la carpeta del grupo. Después, `mcp.RemoveLocal` quita las entradas copiadas del `.claude.json` de la cuenta, con el mismo bloqueo y después de hacer una copia de seguridad en la carpeta de estado. Esto es necesario porque, si un nombre coincide, la entrada de la cuenta ganaría sobre la del grupo. `--keep-local` las conserva.

## Seguridad y casos límite

- **Inicios de sesión.** Los servidores que inician sesión con OAuth guardan sus tokens en el `.credentials.json` de cada cuenta (`mcpOAuth`), no en la lista de servidores, así que compartir un servidor no comparte su inicio de sesión. Inicia sesión una vez por cuenta con `/mcp`, o copia los inicios de sesión de `~/.claude` con [`ccam migrate`](/features/migrate). El grupo nunca recibe un archivo de inicio de sesión.
- **Cuentas en ejecución.** `import` copia al grupo, pero conserva las entradas propias de una cuenta en ejecución y lo indica; vuelve a ejecutarlo cuando termine la sesión. La sincronización omite las cuentas en ejecución.
- **Copias para Remote Control.** La copia de una cuenta en ejecución no se toca: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Con `mcp_strategy = "sync"` no hace falta actualizar nada, porque cada inicio sincroniza.
- **Archivos dañados.** Si `.claude.json` no es JSON válido, la sincronización no lo toca e informa el error. Si no se puede tomar el bloqueo en 2 segundos, ccam informa `could not lock … (Claude Code is writing the file)`.
- **Cuentas que no comparten MCP.** `import` se niega y te indica que primero ejecutes `ccam share side --mcp on`.
- **Nombres de grupo.** `-g` acepta solo un nombre de grupo; se convierte en una carpeta dentro de `groups/`, así que nombres como `../x` se rechazan.

## Relacionado

- [Compartir entre cuentas](/features/share)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Copiar una configuración existente](/features/migrate)
- [Todos los comandos](/docs/commands)
