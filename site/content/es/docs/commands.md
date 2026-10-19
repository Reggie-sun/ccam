---
title: Referencia de comandos
metaTitle: Referencia de comandos de ccam
description: Todos los comandos y opciones de ccam, desde ccam add y ccam <name> hasta ccam usage, ccam next y ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 65fb2b8f6617
---

Pon `--dry-run` antes de cualquier comando para ver lo que haría sin hacerlo, por ejemplo `ccam --dry-run work`.

<!-- commands -->

## Iniciar Claude Code

`ccam <name|number> [arguments]` inicia Claude Code con esa cuenta. ccam no lee los argumentos que van después del nombre; los pasa a `claude` en el mismo orden, y la entrada estándar, el código de salida y las señales se comportan exactamente igual que con `claude`. Un `--` después del nombre es opcional.

`ccam` con argumentos pero sin nombre de cuenta, como `ccam -p "summarise"`, usa la cuenta predeterminada. Establécela con `ccam use <name>`.

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
```

`ccam import` sin una carpeta busca `~/.claude` y `~/.claude-*` y pregunta por cada una. Las carpetas importadas se registran donde están; ccam nunca mueve ni copia un inicio de sesión.

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

`ccam migrate` copia `settings.json`, `CLAUDE.md`, las extensiones y los servidores MCP de usuario de una carpeta de Claude Code a un grupo, para que todas las cuentas que comparten esos elementos los tengan. La carpeta solo se lee. Las opciones `--settings`, `--memory`, `--extensions`, `--mcp` y `--history` copian solo los elementos que indicas; sin ellas, se copia todo excepto el historial de conversaciones. `-g` o `--group` elige el grupo, y `-y` o `--yes` copia sin preguntar y conserva la copia del grupo cuando un archivo es distinto. Consulta [copiar una configuración existente](/docs/sharing#copy-an-existing-setup).

## Servidores MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Uso y cambio de cuenta

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Mantenimiento

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
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

Las variables `XDG_CONFIG_HOME`, `XDG_DATA_HOME` y `XDG_STATE_HOME` cambian la ubicación de estas carpetas.
