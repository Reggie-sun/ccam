---
title: Cambiar a la siguiente cuenta
description: ccam next abre Claude Code con la cuenta que sigue a la que usaste por última vez; --auto elige el menor uso de 5 horas y --continue retoma la conversación.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` abre Claude Code con otra cuenta sin que tengas que nombrarla. Úsalo cuando una cuenta llegue a su límite de uso: `ccam next` a secas toma la cuenta que sigue a la que usaste por última vez, `--auto` toma la que tiene más margen en su ventana de 5 horas y `--continue` sigue la última conversación en la cuenta nueva.

## Uso

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam imprime el cambio como `→ <from> → <to>` en la salida de error estándar y luego inicia Claude Code con la cuenta nueva igual que `ccam <name>`. `<from>` es la cuenta que usaste por última vez. Los argumentos que van después de `--` llegan a `claude` sin cambios.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--auto` | Elige la cuenta con sesión iniciada que tiene el menor uso de 5 horas, y omite las cuentas al 100 % | desactivada |
| `-c`, `--continue` | Solo tiene en cuenta las cuentas que comparten el historial de conversaciones con la cuenta que usaste por última vez, e inicia `claude --continue` | desactivada |
| `-- <claude args>` | Argumentos que se pasan a `claude` después de los de ccam | ninguno |
| `--dry-run` | Se respeta: imprime `→ <from> → <to>` y la línea de comandos que ejecutaría, sin iniciar Claude Code ni registrar un inicio | desactivada |

`--auto` y `--continue` se pueden combinar. Con `--continue`, `--continue` se coloca antes de tus propios argumentos.

## Cómo funciona

Cada inicio a través de ccam escribe la hora en `~/.local/state/ccam/last-used/<name>`. `ccam next` lee estos archivos, comprueba qué cuentas tienen sesión iniciada y elige.

1. `cli.candidates` crea un `next.Candidate` por cada cuenta de `config.toml`: la cuenta, la hora de su último inicio (`launcher.LastUsed`) y si tiene un inicio de sesión (`creds.Has`: un `.credentials.json` legible o, en macOS, un elemento del Llavero).
2. Sin `--auto`, `next.Rotate` ordena las cuentas por número de posición y busca la que se inició más recientemente. Esa es la cuenta de la que te cambias. Si todavía no se ha iniciado ninguna cuenta, es la posición más baja.
3. A partir de ahí avanza en el orden de las posiciones, volviendo de la última posición a la primera, y devuelve la primera cuenta que tiene sesión iniciada.
4. `launcher.LaunchWith` prepara e inicia la cuenta elegida exactamente como lo haría `ccam <name>`: comprueba los enlaces de la cuenta, aplica el uso compartido de MCP, registra la sesión y reemplaza ccam por `claude` mediante `syscall.Exec`.

El recorrido empieza en la última cuenta usada y se detiene en la primera cuenta con sesión iniciada que viene después:

<!-- motion:next-rotate -->

`work` fue la que se inició más recientemente, así que el recorrido empieza después de ella. `side` no tiene sesión iniciada y se omite, y se elige `lab`. Si `lab` tampoco tuviera sesión iniciada, el recorrido volvería al principio hasta `personal` (`% len` en el diagrama). La rotación nunca devuelve la cuenta desde la que empezaste: si no hay otra cuenta con sesión iniciada, se detiene con un error.

### --auto

`--auto` decide por el uso en lugar del orden de las posiciones. `cli.pickAuto` se queda con las cuentas que tienen sesión iniciada, consulta el uso de todas a la vez con el mismo `Fetcher.FetchAll` que usa `ccam usage` y deja que elija `usage.Best()`:

<!-- motion:next-auto -->

`usage.Best()` toma el porcentaje de 5 horas más bajo y omite las cuentas al 100 % o con un uso de 5 horas desconocido. La consulta sigue las reglas de `ccam usage`: las cifras con menos de 60 segundos se reutilizan, cada cuenta tiene 3 segundos, y una cuenta que no responde, o cuyo inicio de sesión venció, se juzga por sus cifras guardadas. Una ventana guardada cuya hora de restablecimiento ya pasó cuenta como 0 %, así que no se descarta una cuenta cuya ventana se restableció desde entonces. `--auto` nunca renueva tokens. A diferencia de la rotación simple, `--auto` puede elegir la cuenta que usaste por última vez si sigue teniendo el mayor margen.

### --continue

`--continue` de Claude Code abre la conversación más reciente de la carpeta actual, lo que solo funciona si la cuenta nueva puede ver esa conversación. Por eso, con `--continue`, ccam solo tiene en cuenta las cuentas del mismo grupo de `history` que la cuenta que usaste por última vez, y luego coloca `--continue` delante de tus argumentos. Las cuentas cuyo historial es `local` no comparten con nadie, así que `--continue` desde una de ellas siempre falla. Activa primero el uso compartido con `ccam share <name> --history on` (consulta [compartir](/docs/sharing)).

## Seguridad y casos límite

- Las cuentas sin inicio de sesión nunca se eligen. Si no hay otra cuenta con sesión iniciada, `ccam next` a secas se detiene con `no logged-in account other than <name>`.
- Con `--continue` y ninguna cuenta que comparta el historial: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Con `--auto` y ninguna cuenta con sesión iniciada: `no logged-in account to pick`. Cuando todas las cuentas están llenas o su uso es desconocido: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` hace solicitudes de red incluso con `--dry-run`. Solo lee, y nunca escribe archivos de tokens.
- La cuenta elegida se inicia por la misma ruta que `ccam <name>`, así que se aplican las mismas comprobaciones, incluida la advertencia cuando otro Claude Code que comparte el mismo historial ya se está ejecutando en esta carpeta.
- `last-used` solo se escribe cuando Claude Code se inicia de verdad. `--dry-run` no lo toca, así que una ejecución de prueba no hace avanzar la rotación.

## Relacionado

- [Límites de uso](/features/usage)
- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Flujo de trabajo con límites de uso](/guides/usage-limit-workflow)
