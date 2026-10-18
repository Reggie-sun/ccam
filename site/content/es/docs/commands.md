---
title: Referencia de comandos
description: Todos los comandos y opciones de ccam, desde ccam add y ccs <name> hasta ccam usage, ccam next y ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` y `ccs` son el mismo programa; usa el que prefieras. Pon `--dry-run` antes de cualquier comando para ver lo que haría sin hacerlo, por ejemplo `ccs --dry-run work`.

<!-- commands -->

## Iniciar Claude Code

`ccs <name|number> [arguments]` inicia Claude Code con esa cuenta. ccam no lee los argumentos que van después del nombre; los pasa a `claude` en el mismo orden, y la entrada estándar, el código de salida y las señales se comportan exactamente igual que con `claude`. Un `--` después del nombre es opcional.

`ccs` con argumentos pero sin nombre de cuenta, como `ccs -p "summarise"`, usa la cuenta predeterminada. Establécela con `ccam use <name>`.

Los atajos `claude1`, `claude2`, etc., son enlaces a ccam en `~/.local/bin`. `ccam shims install --names` también crea `claude-<name>` para cada cuenta.

## Agregar y quitar cuentas

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` sin una carpeta busca `~/.claude` y `~/.claude-*` y pregunta por cada una. Las carpetas importadas se registran donde están; ccam nunca mueve ni copia un inicio de sesión.

## Compartir

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Cada elemento acepta `on` (el grupo `default`), `off` o el nombre de un grupo. Consulta [compartir](/docs/sharing).

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
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Mantenimiento

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Archivos

| Ruta | Contenido |
| --- | --- |
| `~/.config/ccam/config.toml` | Las cuentas y lo que comparten |
| `~/.local/share/ccam/profiles/<name>/` | El `CLAUDE_CONFIG_DIR` de cada cuenta |
| `~/.local/share/ccam/groups/<group>/` | Configuración, servidores MCP, extensiones e historial compartidos |
| `~/.local/state/ccam/` | Caché de uso y horas del último uso |
| `~/.local/bin/claude1`, `ccs` | Enlaces a ccam |

Las variables `XDG_CONFIG_HOME`, `XDG_DATA_HOME` y `XDG_STATE_HOME` cambian la ubicación de estas carpetas.
