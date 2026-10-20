---
title: Referencia de comandos
metaTitle: Referencia de comandos de ccam
description: Todos los comandos y opciones de ccam, desde ccam add y ccam <name> hasta ccam usage, ccam next y ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

Pon `--dry-run` antes de un comando para ver lo que haría sin hacerlo, por ejemplo `ccam --dry-run work`. Funciona al iniciar Claude Code (`ccam <name>`, `ccam run`, `ccam next`) y con `share`, `migrate`, `shortcut`, `update` y `doctor --fix`. Los demás comandos que cambian algo lo rechazan en lugar de ejecutarse.

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

Un atajo como `ccwork` inicia su cuenta de la misma forma: `ccwork --resume` hace lo mismo que `ccam work --resume`. Consulta [atajos](/docs/commands#shortcuts).

## Agregar y quitar cuentas

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` sin una carpeta busca `~/.claude` y `~/.claude-*` y pregunta por cada una. Las carpetas importadas se registran donde están; ccam nunca mueve ni copia un inicio de sesión.

`ccam add --api` agrega una cuenta que usa una clave de API en lugar de un inicio de sesión de Claude; ccam lee la clave en cada inicio y nunca la guarda. Consulta [Cuentas con clave de API](/features/api).

`ccam rename`, `ccam remove` y `ccam share` se niegan a cambiar una cuenta mientras Claude Code se está ejecutando con esa cuenta. Cierra primero la sesión.

`ccam list --json` muestra cada cuenta con su número de posición, atajo, correo electrónico y carpeta. Su objeto `share` asigna a cada elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) un nombre de grupo o `local`.

## Atajos

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Un atajo es un comando al que tú le pones nombre y que inicia Claude Code con una cuenta concreta. Igual que `ccam work`, pasa todos los argumentos a Claude Code sin cambios, así que `ccwork --resume` y `ccwork -p "hi"` funcionan como esperas. Es un enlace a ccam en `~/.local/bin`.

Los atajos son opcionales. El asistente de `ccam add` pregunta si quieres crear uno, `ccam add <name> --shortcut <command>` crea uno sin preguntar y `ccam list` muestra el atajo de cada cuenta, o `-` si no tiene.

Los nombres usan los caracteres a-z, 0-9, `.`, `-` y `_`. Un nombre no puede ser un comando de ccam, el atajo de otra cuenta ni el nombre de otro programa que esté en tu `PATH`.

`ccam rename` conserva el nombre del atajo, y `ccam remove` elimina el atajo junto con la cuenta. `ccam shims install` vuelve a crear los enlaces de todas las cuentas, por ejemplo después de mover ccam.

## Compartir

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Cada elemento acepta `on` (el grupo `default`), `off` o el nombre de un grupo. Consulta [compartir](/docs/sharing).

## Copiar una configuración existente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, las extensiones y los servidores MCP de usuario de una carpeta de Claude Code a un grupo, para que todas las cuentas que comparten esos elementos los tengan. La carpeta solo se lee. Las opciones `--settings`, `--memory`, `--extensions`, `--mcp` y `--history` copian solo los elementos que indicas. Sin ellas, en una terminal una lista de selección te permite elegir, con todo marcado excepto el historial de conversaciones; en otros casos se copia todo excepto el historial de conversaciones. `-g` o `--group` elige el grupo, y `-y` o `--yes` copia sin preguntar y conserva la copia del grupo cuando un archivo es distinto. También copia los inicios de sesión guardados de los servidores al archivo de inicio de sesión propio de cada cuenta. Consulta [copiar una configuración existente](/docs/sharing#copy-an-existing-setup).

## Servidores MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` mueve a su grupo los servidores que solo tiene esa cuenta. Si el grupo ya usa el nombre de un servidor para una definición distinta, ese servidor se omite; agrega `--force` para reemplazar la copia del grupo, de la que ccam hace antes una copia de seguridad.

## Uso y cambio de cuenta

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

`next --auto`, `--fallback` y `ccam pick` eligen de la misma forma: primero por debajo de `auto_max_usage` (`config.toml`, 100 por defecto), según `auto_strategy` (`headroom` o `expiring`), contando la ventana semanal de `--model`, nunca una cuenta al 100 %. Cuando no se puede elegir ninguna, salen con el código 75 o, con `--wait`, esperan a que se restablezca la primera cuenta. Consulta [Cambiar a la siguiente cuenta](/features/next) y [Elegir una cuenta para scripts](/features/pick).

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
| `~/.config/ccam/config.toml` | Las cuentas y lo que comparten |
| `~/.local/share/ccam/profiles/<name>/` | El `CLAUDE_CONFIG_DIR` de cada cuenta |
| `~/.local/share/ccam/groups/<group>/` | Configuración, servidores MCP, extensiones e historial compartidos |
| `~/.local/state/ccam/` | Caché de uso y horas del último uso |
| `~/.local/bin/<shortcut>` | Atajos como `ccwork`, enlaces a ccam |

Las variables `XDG_CONFIG_HOME`, `XDG_DATA_HOME` y `XDG_STATE_HOME` cambian la ubicación de estas carpetas. En Windows, `~` es `%USERPROFILE%` y los atajos son `<shortcut>.exe`.
