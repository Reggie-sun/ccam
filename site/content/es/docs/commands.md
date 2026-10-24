---
title: Referencia de comandos
metaTitle: Referencia de comandos de ccam
description: Todos los comandos y opciones de ccam, desde ccam add y ccam <name> hasta ccam usage, ccam next y ccam doctor.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 36af7e209cd9
---

Pon `--dry-run` antes de un comando para ver lo que haría sin hacerlo, por ejemplo `ccam --dry-run work`. Funciona al iniciar Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) y con `group`, `migrate`, `shortcut`, `update` y `doctor --fix`. Los demás comandos que cambian algo lo rechazan en lugar de ejecutarse.

<!-- commands -->

## Iniciar Claude Code

`ccam <name|number> [arguments]` inicia Claude Code con esa cuenta. ccam no lee los argumentos que van después del nombre; los pasa a `claude` en el mismo orden, y la entrada estándar, el código de salida y las señales se comportan exactamente igual que con `claude`. Un `--` después del nombre es opcional.

`ccam` con argumentos pero sin nombre de cuenta, como `ccam -p "summarise"`, usa la cuenta asignada a la carpeta actual o, si no hay ninguna, la cuenta predeterminada. Establécelas con `ccam folder <name> [folder]` y `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` antes del nombre de la cuenta es para ejecuciones `-p`/`--print`. Cuando la cuenta llega a su límite de uso, ccam vuelve a ejecutar el mismo comando con la cuenta con sesión iniciada que tenga más margen e imprime solo la salida de la ejecución que terminó. Consulta [Cuentas por carpeta](/features/folder) y [Volver a ejecutar al llegar al límite](/features/next#run-again-at-the-limit).

**Asientos de Team y Enterprise.** La elección automática (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por lo tanto, usar la elección automática entre asientos puede infringir los términos de Anthropic. Consulta con tu organización y con Anthropic antes de depender de ella. Los planes personales (Pro, Max) se tratan en la [nota de las preguntas frecuentes sobre los términos](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Un atajo como `ccwork` inicia su cuenta de la misma forma: `ccwork --resume` hace lo mismo que `ccam work --resume`. Consulta [atajos](/docs/commands#shortcuts).

## Agregar y quitar cuentas

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` sin una carpeta busca `~/.claude` y `~/.claude-*` y pregunta por cada una; una carpeta que falla se informa como `Skipped <path>: <reason>` y las demás se registran, y una cuyo nombre ya es el de una cuenta necesita `--name`, que se indica junto con la carpeta. Las carpetas importadas se registran donde están; ccam nunca mueve ni copia un inicio de sesión.

`ccam login --expired` también inicia sesión en las cuentas cuyo inicio de sesión borró Claude Code. `ccam remove` cierra la sesión de una cuenta gestionada antes de que su carpeta vaya a la papelera (no cuando otra cuenta tiene el mismo inicio de sesión), y `ccam remove` de una cuenta importada sale de sus grupos con una copia de lo que comparten. En macOS, `ccam rename` conserva el nombre de la carpeta, ya que el Llavero guarda el inicio de sesión con la ruta de la carpeta.

`ccam add --api` agrega una cuenta que usa una clave de API en lugar de un inicio de sesión de Claude; ccam lee la clave en cada inicio y nunca la guarda. Consulta [Cuentas con clave de API](/features/api).

`ccam rename`, `ccam remove` y los comandos `ccam group` se niegan a cambiar una cuenta mientras Claude Code se está ejecutando con esa cuenta, sin importar cómo se haya iniciado (también las sesiones en segundo plano). Cierra primero la sesión. También se niegan mientras un cambio de grupos esté sin terminar: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` muestra cada cuenta con su número de posición, atajo, correo electrónico y carpeta. Su arreglo `groups` lista los grupos de la cuenta, y su objeto `share` asigna a cada elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) el grupo del que proviene, o `local`.

## Atajos

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Un atajo es un comando al que tú le pones nombre y que inicia Claude Code con una cuenta concreta. Igual que `ccam work`, pasa todos los argumentos a Claude Code sin cambios, así que `ccwork --resume` y `ccwork -p "hi"` funcionan como esperas. Es un enlace a ccam en `~/.local/bin`.

Los atajos son opcionales. El asistente de `ccam add` pregunta si quieres crear uno, `ccam add <name> --shortcut <command>` crea uno sin preguntar y `ccam list` muestra el atajo de cada cuenta, o `-` si no tiene.

Los nombres usan los caracteres a-z, 0-9, `.`, `-` y `_`. Un nombre no puede ser un comando de ccam, el atajo de otra cuenta ni el nombre de otro programa que esté en tu `PATH`, ni un archivo de `~/.local/bin` que ccam no haya creado, como tu propia copia o enlace de ccam. Un atajo editado a mano en `config.toml` que incumple estas reglas, que es `claude` o `ccam`, o que se da a dos cuentas hace que `config.toml` no se pueda cargar, y se nombra la cuenta.

`ccam rename` conserva el nombre del atajo, y `ccam remove` elimina el atajo junto con la cuenta. `ccam shims install` vuelve a crear los enlaces de todas las cuentas, por ejemplo después de mover ccam; un archivo que estorba a un atajo no detiene los demás.

## Compartir

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Un grupo comparte elementos entre sus cuentas; una cuenta puede estar en varios grupos, pero cada elemento proviene de un solo grupo. Consulta [compartir con grupos](/features/share).

## Copiar una configuración existente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, las extensiones y los servidores MCP de usuario de una carpeta de Claude Code a un grupo, para que todas las cuentas que comparten esos elementos los tengan. La carpeta solo se lee. Las opciones `--settings`, `--memory`, `--extensions`, `--mcp` y `--history` copian solo los elementos que indicas. Sin ellas, en una terminal una lista de selección te permite elegir, con todo marcado excepto el historial de conversaciones; en otros casos se copia todo excepto el historial de conversaciones. `-g` o `--group` elige el grupo, y `-y` o `--yes` copia sin preguntar y conserva la copia del grupo cuando un archivo es distinto. Los inicios de sesión MCP no se copian, ya que un inicio de sesión copiado cierra la sesión de los demás que lo tienen; el plan nombra las cuentas que inician sesión con `/mcp`. `--mcp` también funciona cuando `~/.claude` está vacío o no existe pero `~/.claude.json` tiene servidores. Consulta [copiar una configuración existente](/docs/sharing#copy-an-existing-setup).

## Servidores MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` mueve a su grupo los servidores que solo tiene esa cuenta, también los que tiene solo para una carpeta (`claude mcp add` a secas), que entonces se cargan en todas las carpetas. Si el grupo ya usa el nombre de un servidor para una definición distinta, ese servidor se omite; agrega `--force` para reemplazar la copia del grupo, de la que ccam hace antes una copia de seguridad.

## Uso y cambio de cuenta

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
```

`next --auto`, `--fallback` y `ccam pick` eligen de la misma forma: primero por debajo de `auto_max_usage` (`config.toml`, 100 por defecto), según `auto_strategy` (`headroom`, `expiring` o `balance`), contando la ventana semanal de `--model`, nunca una cuenta al 100 % (una ventana que Claude Code informó llena cuenta como 100 % hasta su restablecimiento), nunca un inicio de sesión con clave de API, y dentro de la organización de la cuenta desde la que parten (`--cross-org` o `cross_org = true` para buscar más allá). `--wait` toma su valor después de `=` (`--wait=2h`); escrito con un espacio se rechaza con una sugerencia. Cuando no se puede elegir ninguna, salen con el código 75 o, con `--wait`, esperan a que se restablezca la primera cuenta. `ccam next --continue` lleva una conversación a otra cuenta incluso entre cuentas que no comparten el historial: Claude Code la copia entonces en la nueva cuenta. A diferencia de `ccam <name>`, `ccam next` lee sus propias opciones, así que las opciones de Claude Code van después de `--`; sin él, ccam se detiene con `unknown flag`. Dentro de Claude Code, `/ccam:switch [account]` del plugin pasa una conversación a otra cuenta, y `ccam hook` es lo que ejecutan sus hooks. Consulta [Cambiar a la siguiente cuenta](/features/next), [Elegir una cuenta para scripts](/features/pick) y [Sesiones en segundo plano de todas las cuentas](/features/agents).

## Mantenimiento

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` comprueba la suma de verificación SHA-256 de la versión antes de reemplazar ccam, y `--version v0.2.0` instala una versión concreta. Si instalaste ccam desde el paquete `.deb`, actualízalo con apt. Consulta [actualizar](/docs/install#update).

## Archivos

| Ruta | Contenido |
| --- | --- |
| `~/.config/ccam/config.toml` | Cuentas, grupos y ajustes (`theme`, `usage_source`, `relay`, `cross_org` …). Los ajustes van al principio, antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`: uno que está debajo de una tabla pertenece a esa tabla, así que ccam lo nombra en cada inicio y en `ccam doctor`, y los comandos que guardan la configuración se niegan hasta que se mueva. Los comentarios no se conservan cuando ccam guarda el archivo |
| `~/.local/share/ccam/profiles/<name>/` | El `CLAUDE_CONFIG_DIR` de cada cuenta |
| `~/.local/share/ccam/groups/<group>/` | Configuración, servidores MCP, extensiones e historial compartidos |
| `~/.local/state/ccam/` | Caché de uso y registros de límite (`usage/<name>.json`, `usage/<name>.limit.json`), horas del último uso, marcadores de sesión y bloqueos (`locks/`), los atajos que creó ccam (`shortcuts.json`), un cambio de grupos sin terminar (`group-change.json`) y las notas que dejan `/ccam:switch` y el hook de límite (`handoff/<account>/<session>.json`) |
| `~/.local/bin/<shortcut>` | Atajos como `ccwork`, enlaces a ccam |

Las variables `XDG_CONFIG_HOME`, `XDG_DATA_HOME` y `XDG_STATE_HOME` cambian la ubicación de estas carpetas. En Windows, `~` es `%USERPROFILE%` y los atajos son `<shortcut>.exe`.
