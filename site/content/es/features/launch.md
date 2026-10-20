---
title: Iniciar Claude Code
description: Inicia Claude Code como cualquier cuenta con ccam <name|number>, ccam run o el selector; los argumentos pasan intactos a claude y --dry-run muestra el comando.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 401eed57e533
---

`ccam work` inicia Claude Code con la cuenta `work`, y todo lo que va después del nombre llega a `claude` tal cual. Úsalo donde escribirías `claude`: en una terminal, en un script o con la entrada por una tubería.

## Uso

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Indica el nombre de una cuenta o su número de `ccam list`. ccam no lee los argumentos que van después del nombre, así que `ccam work --help` muestra la ayuda de Claude Code. Un `--` después del nombre es opcional. La entrada estándar, el código de salida y las señales se comportan igual que con `claude`.

Cuando el primer argumento empieza con `-` y no es una de las opciones propias de ccam, ccam inicia la cuenta predeterminada con todos los argumentos. La predeterminada es la cuenta establecida con `ccam use`, o la cuenta con el número más bajo si no has establecido ninguna.

### El selector

`ccam` sin argumentos abre el selector cuando la entrada y la salida estándar son una terminal, y en otro caso imprime `ccam list`. Muestra el número, el correo electrónico, el uso de 5 horas, el último uso y los elementos compartidos de cada cuenta; el uso se completa a medida que responde cada cuenta.

| Tecla | Qué hace |
| --- | --- |
| `Enter` | Inicia Claude Code con la cuenta seleccionada |
| `1` a `9` | Inicia la cuenta con ese número |
| `Up`/`Down`, `k`/`j` | Mueve la selección |
| `a` | Cierra el selector y ejecuta el asistente de `ccam add` |
| `s` | Edita los elementos compartidos de la cuenta seleccionada (`Space` o `x` cambia un elemento, `Enter` aplica, `Esc` vuelve atrás) |
| `d` | Convierte la cuenta seleccionada en la predeterminada |
| `r` | Vuelve a obtener el uso, sin la caché breve |
| `q`, `Esc`, `Ctrl+C` | Sale |

El selector se cierra antes de actuar, así que la terminal se restaura antes de que se inicie Claude Code. Después de `d` o de un cambio en el uso compartido, se vuelve a abrir. Un cambio en el uso compartido muestra el plan y pregunta primero, igual que `ccam share`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--dry-run` | Imprime el entorno y el comando que ccam ejecutaría, y sale | desactivada |
| `-v`, `--version` | Imprime `ccam <version>` (solo como primer argumento) | |
| `-h`, `--help` | La ayuda de ccam, como `ccam --help` o `ccam run --help` | |

Reglas de los argumentos:

- ccam lee sus propias opciones solo antes del nombre de la cuenta. Después del nombre, todos los argumentos son de Claude Code, incluidos `--dry-run` y `--help`.
- `--dry-run` funciona como `ccam --dry-run work`, `ccam run --dry-run work` y `ccam --dry-run` (la cuenta predeterminada). Se respeta, no se rechaza.
- `ccam run <name|number> [claude args]` es la forma larga de `ccam <name>`. `ccam run` sin un nombre usa la cuenta predeterminada.
- Los nombres de cuenta no pueden ser nombres de comandos de ccam como `add`, `list`, `run` o `usage`, así que `ccam <word>` nunca es ambiguo.

La ejecución de prueba imprime una línea que podrías pegar en una shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Cómo funciona

Cada inicio sigue los mismos pasos, desde la línea de comandos hasta `exec`:

<!-- motion:launch -->

Cada recuadro se ilumina mientras se ejecuta su paso, y las líneas de abajo son el entorno y los argumentos con los que se inicia Claude Code.

1. `cmd/ccam/main.go` llama a `app.Main`. Si ccam se inició con otro nombre, es un [atajo](/features/shortcuts) y va directamente al iniciador.
2. `app.parse` examina `os.Args` antes de que intervenga cobra. Solo `ccam` a secas (el selector), un nombre de comando de ccam (`config.Reserved`), `-h`/`--help` y `-v`/`--version` van a `cli.Execute`. Todo lo demás es un inicio.
3. `launcher.Launch` lee `~/.config/ccam/config.toml`, el único archivo que analiza.
4. `launcher.Prepare` busca la cuenta por nombre o número, descarta un `--` inicial y resuelve el ejecutable `claude` con `ResolveClaude`. La búsqueda en PATH se guarda en caché en `~/.local/state/ccam/claude-path` junto con la fecha de modificación del archivo, así que ccam no busca en `PATH` en cada inicio. Un archivo modificado o que falta provoca una búsqueda nueva. Una ruta `claude_bin` en `config.toml` se usa tal cual. ccam rechaza un `claude` que sea el propio ccam, comparando el inodo, porque eso crearía un bucle.
5. `checkLinks` compara la carpeta de la cuenta `~/.local/share/ccam/profiles/<name>/` con su configuración de uso compartido mediante `share.Check`, que solo llama a `lstat` y `readlink`, una o dos llamadas al sistema por elemento. Un enlace que falta se crea en el momento. Cualquier cosa que podría perder datos detiene el inicio.
6. `applyMCP` agrega los servidores MCP compartidos (consulta más abajo).
7. `session.Register` registra la sesión. ccam escribe un marcador con el nombre de su PID, que contiene la hora de inicio del proceso, en `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (o en `~/.local/state/ccam/locks/` si no hay `XDG_RUNTIME_DIR`). En una cuenta que comparte el historial de conversaciones, `historyCheck` primero busca una sesión en ejecución en `history-<group>-<hash>`, donde el hash es el de la carpeta actual, y también se registra ahí. `exec` conserva el PID, así que el marcador pasa a nombrar al propio proceso de Claude Code. Es un marcador y no un bloqueo retenido, porque los procesos que inicia Claude Code heredarían el bloqueo y lo conservarían después de que termine la sesión.
8. `Plan.Exec` establece el título de la terminal en `claude · <name>` (cuando la salida de error estándar es una terminal), escribe la hora de inicio en `~/.local/state/ccam/last-used/<name>` sin `fsync` y llama a `syscall.Exec`.

Claude Code se ejecuta con `CLAUDE_CONFIG_DIR` apuntando a la carpeta de la cuenta y `CCAM_PROFILE` con el nombre de la cuenta. ccam quita cualquier valor heredado de estas variables. Para una carpeta `~/.claude` registrada con `ccam import` en macOS, también establece `CLAUDE_SECURESTORAGE_CONFIG_DIR` con un valor vacío, para que la cuenta y `claude` a secas usen el mismo inicio de sesión del Llavero.

### Por qué exec

`syscall.Exec` reemplaza el proceso de ccam por `claude`, así que ccam no se queda en memoria, las pulsaciones de teclas, las señales y el código de salida son los de Claude Code, y los argumentos siguen siendo un arreglo, nunca una cadena de shell. El trabajo posterior a una sesión, como detectar servidores agregados con `claude mcp add`, ocurre en el siguiente inicio: ccam compara la fecha de modificación de `.claude.json` con `~/.local/state/ccam/mcp/<name>/seen` y analiza el archivo solo cuando cambió.

### Cómo se enruta una línea de comandos

Tres líneas de comandos, tres rutas por `internal/app`:

<!-- motion:launch-route -->

Un atajo se reconoce por `argv[0]` y `ccam <name>` por el primer argumento, y ambos terminan en `syscall.Exec`. Solo un comando de ccam como `ccam usage` llega a cobra.

### Servidores MCP compartidos

Con la estrategia predeterminada `flag`, ccam pasa los servidores MCP del grupo a Claude Code como un argumento, colocado antes de los tuyos:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` devuelve un argumento solo si existe el `mcp.json` del grupo, y `Plan.InsertArgs` lo coloca justo después de `claude`. Tus propias opciones `--mcp-config` o `--strict-mcp-config` se siguen aplicando encima. Es un solo argumento, `--mcp-config=<path>`, porque `--mcp-config` de Claude Code acepta varios valores: escrito como dos argumentos, también leería una instrucción posterior como un archivo de configuración.

Con `mcp_strategy = "sync"` en `config.toml`, ccam escribe los servidores en el `.claude.json` de la cuenta, y lo omite si la cuenta ya está en ejecución.

### Remote Control

`ccam work rc` (o `ccam work remote-control`) inicia el Remote Control de Claude Code como `work`. `claude remote-control` rechaza `--mcp-config` antes de su verbo y no acepta ninguno después, y sus sesiones leen los servidores MCP solo del `.claude.json` de la cuenta. Por eso la estrategia de opción funciona de otra manera aquí:

<!-- motion:launch-rc -->

1. `remoteControl` toma como verbo de Claude Code el primer argumento después de la cuenta que no empieza por `-`. Si es `remote-control` o `rc`, `mcp.FlagArgs` no añade nada.
2. `mcp.Synced` compara la hora de modificación del `mcp.json` del grupo con el `sync-stamp` de la cuenta. Si `.claude.json` ya contiene los servidores actuales del grupo, o el grupo no tiene ninguno, Remote Control se inicia sin preguntar.
3. Si no, ccam nombra los servidores que les faltarían a las sesiones y pregunta:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. Con `y`, `mcp.Sync` los copia como lo hace la estrategia `sync`: con el bloqueo de Claude Code, después de una copia de seguridad y conservando los servidores propios de la cuenta. Con `N`, o sin una terminal, Remote Control se inicia sin ellos, y ccam sugiere `mcp_strategy = "sync"`, que los copia antes de cada inicio.

ccam pregunta en cada inicio de Remote Control hasta que los servidores se copian y, después de una copia, solo cuando el `mcp.json` del grupo vuelve a cambiar. Una copia en `.claude.json` gana sobre `--mcp-config` en los inicios normales posteriores, así que `ccam mcp add`, `remove` e `import` actualizan las copias que tienen otras cuentas ([Servidores MCP compartidos](/features/mcp)).

### El presupuesto de 10 ms

ccam agrega como máximo 10 ms (p95) antes de que se inicie Claude Code. La ruta de inicio nunca construye el árbol de comandos de cobra, no inicia Bubble Tea ni abre una conexión de red, y `TestPackageInitIsCheap` falla si las funciones `init()` de los paquetes tardan más de 5 ms en total. `scripts/perf-check.sh` mide `ccam --dry-run work`, `ccam work` y un atajo con `hyperfine` y falla por encima de 10 ms, y `BenchmarkLaunchPath` en `internal/launcher` ejecuta los pasos en un bucle. Para ver el tiempo de cada paso, define `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Seguridad y casos límite

- Un nombre o número desconocido se detiene con `no account "<name>"` y remite a `ccam list`. Si no hay ninguna cuenta, ccam te pide que agregues una con `ccam add <name>`.
- Si una ruta compartida se convirtió en un archivo normal o apunta al lugar equivocado, ccam se detiene con `stopped because of share link problems:`, las lista y te pide que ejecutes `ccam doctor --fix`. No mueve datos en la ruta de inicio.
- Si otro Claude Code que comparte el historial de conversaciones con esta cuenta (el mismo grupo) se está ejecutando en la misma carpeta, ccam advierte que el historial se puede mezclar y pregunta `Run anyway? [y/N]`. Sin una terminal, muestra la advertencia y continúa.
- `--dry-run` omite la sincronización de MCP, el aviso de servidores nuevos, la comprobación del historial, el marcador de sesión, el título y la hora del último uso, y no hace exec. Sí comprueba los enlaces compartidos y crea un enlace que falta, como lo haría un inicio normal.
- `ccam work rc` cuando `work` ya está en ejecución no toca `.claude.json`, lo indica e inicia Remote Control sin los servidores del grupo; vuelve a ejecutarlo cuando termine la sesión para copiarlos. Con `--dry-run` imprime el mensaje, pero no pregunta ni copia.
- Con `mcp_strategy = "sync"`, una sincronización fallida imprime `ccam: MCP sync failed (launching anyway): …` y Claude Code se inicia con los servidores que ya tenía.

## Relacionado

- [Atajos](/features/shortcuts)
- [Usar Remote Control con ccam](/guides/remote-control)
- [Administrar cuentas](/features/accounts)
- [Referencia de comandos](/docs/commands)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Cómo funciona CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
