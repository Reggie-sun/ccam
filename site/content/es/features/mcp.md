---
title: Servidores MCP compartidos
description: Mantén una lista de servidores MCP por grupo con ccam mcp add, add-json, list, remove e import, que se pasa con --mcp-config o se sincroniza en cada cuenta.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: "27a06ad35295"
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
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
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
only in account side, for /home/me/code/app only: notes
```

Debajo de los grupos, `ccam mcp list` nombra cada servidor que no llega a las demás cuentas:

- `only in account work: pg` para un servidor que una cuenta agregó para sí misma;
- `only in account side, for /home/me/code/app only: notes` para uno agregado para una carpeta (`claude mcp add` a secas, consulta más abajo);
- `account work has its own …, used in its sessions in place of the group's` cuando un servidor propio de una cuenta tiene el nombre de un servidor del grupo;
- el `.mcp.json` de esta carpeta, con las cuentas en las que está `waiting for approval in`;
- los servidores que están `in no account, only in ~/.claude.json (plain claude's)`, con el comando [`ccam migrate --mcp`](/features/migrate) que los copia.

Agregar un nombre que el grupo ya tiene reemplaza su definición e imprime `updated` en lugar de `added`. Los servidores llegan a una cuenta solo si está en ese grupo y el grupo comparte MCP (`ccam group set <group> --add mcp`); cuando el grupo no lo hace, `add` lo indica.

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
| `--dry-run` (global) | solo la admite `mcp list`; `add`, `add-json`, `remove` e `import` la rechazan y no cambian nada | desactivada |

`add` necesita un comando después de `--` o bien `--url`; si no, falla con `a command to run (after --) or --url is required`. `--header` se aplica a los servidores con `--url` y `--env` a los servidores con comando. `add-json` acepta el mismo JSON que `claude mcp add-json`, y la definición debe ser un objeto con `command` o `url`. `import` siempre usa el grupo a través del cual la cuenta comparte MCP, y `list` muestra todos los grupos, así que `-g` no los cambia.

## Cómo funciona

Los archivos que intervienen:

- `~/.local/share/ccam/groups/<group>/mcp.json`: los servidores compartidos, `{"mcpServers": {…}}`, el formato que lee `--mcp-config`. Se escribe con `fsutil.AtomicWrite` y el modo `0600`, porque las entradas de servidores suelen llevar claves de API, y solo con el bloqueo `<group>/.mcp.json.lock`, así que dos comandos `ccam mcp add` ejecutados al mismo tiempo conservan ambos sus servidores.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: el archivo propio de Claude Code con los servidores de la cuenta. ccam lee solo su valor `mcpServers` con el escáner `profile.TopLevelValue` y nunca analiza el archivo completo en la ruta de inicio.
- `~/.local/state/ccam/mcp/<name>/`: el archivo `seen`, el `sync-stamp` (con un hash de cada copia que escribió ccam, para distinguir una copia que nadie cambió de una que editaste), `sync-checked` (lo que vio la estrategia de sincronización por última vez), `kept` (los servidores que la cuenta conservó al dejar de compartir MCP) y las copias de seguridad de `.claude.json`.

Hay dos estrategias, que se eligen con `mcp_strategy` en `~/.config/ccam/config.toml`. La predeterminada es la estrategia de opción:

<!-- motion:mcp-flag -->

1. ccam busca el grupo de MCP de la cuenta (`launcher.applyMCP`) y comprueba con un solo `stat` que el `mcp.json` del grupo existe, sin analizarlo (`mcp.FlagArgs`).
2. ccam coloca `--mcp-config=<path>` justo después de `claude`, antes de tus propios argumentos, así que tus propias opciones `--mcp-config` o `--strict-mcp-config` se siguen aplicando encima (`Plan.InsertArgs`).
3. ccam comprueba si el `.claude.json` de la cuenta cambió desde el último inicio, y solo entonces analiza sus servidores (`mcp.NewLocal`).
4. ccam se reemplaza a sí mismo por Claude Code (`syscall.Exec`).

La ruta se pasa como un solo argumento, `--mcp-config=<path>`, no como dos. `--mcp-config` de Claude Code acepta varios valores, así que `claude --mcp-config a.json "prompt"` leería la instrucción como un segundo archivo de configuración y fallaría con `MCP config file not found`. La forma de un solo argumento evita esto, así que una instrucción o un subcomando que venga después funciona. `ccam work mcp list`, `get` y `remove` no reciben la opción, ya que Claude Code la rechaza allí; en su lugar, ccam nombra los servidores del grupo en la salida de error estándar ([detalles](/features/launch#shared-mcp-servers)). Claude Code carga los servidores de `--mcp-config` además de los propios de la cuenta; si un nombre coincide, gana la entrada de `.claude.json`. La estrategia de opción solo escribe en `.claude.json` para Remote Control, cuando una cuenta deja de compartir MCP o vuelve a hacerlo ([más abajo](#leaving-and-joining)), y en `ccam mcp import`.

`claude remote-control` (`rc`) no acepta `--mcp-config`, y sus sesiones leen los servidores solo de `.claude.json`. Para `ccam work rc`, ccam omite la opción y, en una terminal, pregunta si quieres copiar los servidores del grupo en `.claude.json`, con una copia de seguridad, como hace la estrategia `sync`. Pregunta hasta que existe una copia y, después, solo cuando cambia el `mcp.json` del grupo o falta un servidor copiado ([Remote Control](/features/launch#remote-control)). Como esa copia gana sobre `--mcp-config` si un nombre coincide, `ccam mcp add`, `add-json`, `remove` e `import` actualizan las copias que tienen otras cuentas del grupo:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

La línea `✓ Updated` aparece solo cuando una copia cambió. Una copia que creó una versión anterior de ccam se actualiza de la misma forma, y también una posterior a `ccam migrate`.

Un servidor agregado dentro de una sesión con `claude mcp add` queda solo en el `.claude.json` de esa cuenta. `claude mcp add` a secas usa el ámbito local: el servidor se carga para esa única cuenta en esa única carpeta (`-s user` lo deja disponible para la cuenta en todas las carpetas). ccam no puede verlo cuando ocurre, porque Claude Code ya ha reemplazado el proceso de ccam, así que lo comprueba en el siguiente inicio. El archivo `seen` guarda la fecha de modificación de `.claude.json` y los nombres de servidores ya informados. Cuando la fecha es distinta, ccam compara los nombres con el grupo e imprime una sola vez un aviso por cada servidor nuevo, sin esperar ninguna respuesta, con ambas estrategias:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

El aviso también nombra un servidor propio de una cuenta que tiene el nombre de un servidor del grupo, ya que el de la cuenta gana en sus sesiones.

Con `mcp_strategy = "sync"`, los servidores aparecen en `claude mcp list` y `/mcp` como servidores de usuario normales, y ccam los escribe en el `.claude.json` de cada cuenta antes de iniciar:

<!-- motion:mcp-sync -->

1. ccam compara la fecha de modificación del `mcp.json` del grupo con `sync-stamp`, y el `.claude.json` de la cuenta con lo que vio por última vez (`mcp.Sync`). Si ninguno cambió, no pasa nada más. Si no, también comprueba que cada servidor del grupo siga en `.claude.json`, y vuelve a poner una copia que falta (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Toma el bloqueo propio de Claude Code, la carpeta `.claude.json.lock` junto al archivo, creada con `mkdir`. Espera hasta 2 segundos y considera obsoleto un bloqueo con más de 10 segundos, igual que la biblioteca de bloqueo de Claude Code.
3. Guarda el `.claude.json` actual en `backups/` dentro de la carpeta de estado y conserva los 5 más recientes.
4. Establece cada servidor del grupo en `mcpServers`, quita los servidores que escribió antes y que el grupo ya no tiene, y escribe el archivo de forma atómica con su modo original, solo cuando algo cambió. Solo cambia o quita las copias que escribió y que nadie cambió desde entonces. Si la cuenta tiene un servidor propio con el mismo nombre que un servidor del grupo, ccam conserva el servidor de la cuenta y lo indica una vez: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Registra la nueva marca con un hash de cada copia, y Claude Code se inicia sin `--mcp-config`.

La sincronización se omite, con `ccam: account work is already running, skipping MCP sync`, cuando hay una sesión de la cuenta en ejecución, sin importar cómo se haya iniciado Claude Code, porque Claude Code reescribe `.claude.json` mientras se ejecuta. Una sincronización fallida imprime `ccam: MCP sync failed (launching anyway): …` e inicia de todos modos.

### Salir y volver a unirse

Una cuenta que deja de compartir MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) conserva los servidores del grupo: el plan los copia en su `.claude.json`, con una copia de seguridad, y sus propios servidores con los mismos nombres se mantienen. Con `--fresh` no recibe ninguno, y el plan enumera los servidores que dejan de cargarse. Cuando se une otra vez a un grupo que tiene esos servidores, las copias conservadas que nadie cambió se quitan, para que le lleguen las ediciones posteriores del grupo.

`ccam mcp import` mueve los servidores propios de una cuenta a su grupo:

<!-- motion:mcp-import -->

ccam lee los servidores de la cuenta, los propios y los que tiene solo para una carpeta, y compara cada uno con el servidor del grupo que tiene el mismo nombre mediante `mcp.SameDef`, que ignora el orden de las claves y los espacios. Un servidor de carpeta se convierte en un servidor del grupo para todas las carpetas (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); un nombre con dos definiciones de carpeta distintas se rechaza. Los servidores de plugins no se ofrecen, ya que el plugin los aporta. Un servidor que el grupo no tiene, o que tiene con la misma definición, se copia. Uno distinto se omite, porque la definición del grupo es la que usan todas las demás cuentas:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Con `--force`, el servidor del grupo se reemplaza, y el `mcp.json` anterior se guarda primero en `.ccam-backup/<time>/mcp.json` dentro de la carpeta del grupo. Después, ccam quita las entradas copiadas del `.claude.json` de la cuenta, las propias y las de la carpeta en una sola escritura, con el mismo bloqueo y después de hacer una copia de seguridad en la carpeta de estado (`mcp.Remove`). ccam mantiene bloqueado el `mcp.json` del grupo desde que lee el archivo hasta que se quitan esas entradas, así que nadie quita el servidor del grupo entretanto. Esto es necesario porque, si un nombre coincide, la entrada de la cuenta ganaría sobre la del grupo. Con `--keep-local`, las entradas se quedan en el `.claude.json` de la cuenta.

## Seguridad y casos límite

- **Inicios de sesión.** Los servidores que inician sesión con OAuth guardan sus tokens en el `.credentials.json` de cada cuenta (`mcpOAuth`), no en la lista de servidores, así que compartir un servidor no comparte su inicio de sesión. Inicia sesión una vez por cuenta con `/mcp`. ccam nunca copia inicios de sesión, y [`ccam migrate`](/features/migrate) tampoco: un token de actualización copiado cierra la sesión de todos los demás que lo tienen. [`ccam ext`](/features/ext) y `ccam doctor` nombran los servidores en los que una cuenta no ha iniciado sesión mientras otra cuenta sí, y `ccam next --continue` indica cuáles le faltan a la nueva cuenta. El grupo nunca recibe un archivo de inicio de sesión.
- **Cuentas en ejecución.** `import` copia al grupo, pero conserva las entradas propias de una cuenta en ejecución y lo indica; vuelve a ejecutarlo cuando termine la sesión. La sincronización omite las cuentas en ejecución. Ambos detectan también las sesiones de Claude Code iniciadas sin ccam.
- **Copias para Remote Control.** La copia de una cuenta en ejecución no se toca: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Con `mcp_strategy = "sync"` no hace falta actualizar nada, porque cada inicio sincroniza.
- **Archivos dañados.** Si `.claude.json` no es JSON válido, la sincronización no lo toca e informa el error. Si no se puede tomar el bloqueo en 2 segundos, ccam informa `could not lock … (Claude Code is writing the file)`.
- **Cuentas que no comparten MCP.** `import` se niega y te indica que primero pongas la cuenta en un grupo que comparta MCP (`ccam group join <group> side`).
- **Nombres de grupo.** `-g` acepta solo un nombre de grupo; se convierte en una carpeta dentro de `groups/`, así que nombres como `../x` se rechazan.

## Relacionado

- [Compartir con grupos](/features/share)
- [Usar Remote Control con ccam](/guides/remote-control)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Copiar una configuración existente](/features/migrate)
- [Referencia de comandos](/docs/commands)
