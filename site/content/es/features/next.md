---
title: Cambiar a la siguiente cuenta
description: ccam next abre Claude Code con otra cuenta; --auto elige la de más margen, --continue lleva consigo la conversación y --fallback reejecuta trabajos -p.
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` abre Claude Code con otra cuenta sin que tengas que nombrarla. Úsalo cuando una cuenta llegue a su límite de uso: `ccam next` a secas toma la cuenta que sigue a la que usaste por última vez, `--auto` toma la que tiene más margen y `--continue` lleva la última conversación de esta carpeta a la cuenta nueva, compartan o no el historial las dos. Dentro de Claude Code, el [plugin de ccam](#inside-claude-code) te avisa cuando un límite detiene una sesión y `/ccam:switch` pasa la conversación a otra cuenta. Para scripts que ejecutan `claude -p`, `ccam --fallback` hace el cambio por sí solo (consulta [Volver a ejecutar al llegar al límite](#run-again-at-the-limit)).

## Uso

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam imprime el cambio como `→ <from> → <to>` en la salida de error estándar, con el motivo de una elección de `--auto`, y luego inicia Claude Code con la cuenta nueva igual que `ccam <name>`. `<from>` es la cuenta que usaste por última vez o, con `--continue`, la cuenta a la que pertenece la conversación. Los argumentos que van después de `--` llegan a `claude` sin cambios.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--auto` | Elige la cuenta con sesión iniciada según el uso (ver más abajo), y omite las cuentas al 100 % | desactivada |
| `--max-usage <1-100>` | Con `--auto`: prefiere las cuentas por debajo de este uso | `auto_max_usage`, si no 100 |
| `--strategy headroom\|expiring\|balance` | Con `--auto`: el menor uso de 5 horas, el restablecimiento de 5 horas más próximo o la cuota semanal que antes quedaría sin usar | `auto_strategy`, si no `headroom` |
| `--cross-org` | Con `--auto`: elige también cuentas de otra organización (ver más abajo) | `cross_org`, si no desactivada |
| `--wait[=<duration>]` | Con `--auto`: cuando ninguna cuenta tiene margen, espera a que se restablezca la primera, como máximo este tiempo | desactivada; `5h` si se da sola |
| `-c`, `--continue` | Lleva la última conversación de esta carpeta a la cuenta nueva (ver más abajo) | desactivada |
| `-- <claude args>` | Argumentos que se pasan a `claude` después de los de ccam | ninguno |
| `--dry-run` | Se respeta: imprime `→ <from> → <to>` y la línea de comandos que ejecutaría, sin iniciar Claude Code ni registrar un inicio | desactivada |

`--auto` y `--continue` se pueden combinar. Con `--continue`, los argumentos de `--resume` se colocan antes de los tuyos.

Las opciones propias de Claude Code, como `--dangerously-skip-permissions` o `--model`, van después de `--`. A diferencia de `ccam <name>`, que pasa todo lo que va después del nombre, `ccam next` lee sus argumentos por sí mismo, así que sin el `--` se detiene con `ccam: unknown flag: --dangerously-skip-permissions`.

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

`--auto` decide por el uso en lugar del orden de las posiciones. `cli.pickAuto` se queda con las cuentas que tienen sesión iniciada, consulta el uso de todas a la vez con el mismo `Fetcher.FetchAll` que usa `ccam usage` y deja que elija `usage.Policy.Choose()`:

<!-- motion:next-auto -->

Una cuenta se juzga por sus ventanas de 5 horas y semanal y, cuando `--model` está entre los argumentos para Claude Code (después de `--`), también por la ventana semanal de ese modelo (`opus`, `sonnet`, `haiku` o `fable` en el nombre). `Policy.Choose()` reparte las cuentas en tres grupos:

1. **Con margen**: todas las ventanas por debajo de `auto_max_usage` (`config.toml`, 100 por defecto). La estrategia elige entre ellas: `headroom` (predeterminada) toma el menor uso de 5 horas; `expiring` toma la ventana de 5 horas que se restablece antes, para que su cuota restante se use antes de perderse; `balance` toma la cuenta con más cuota semanal restante por día hasta que se restablezca su semana, y cuenta en su contra cada sesión de ccam en ejecución, para que la cuota semanal no quede sin usar y las sesiones en paralelo se repartan.
2. **Por encima**: por debajo del 100 %, pero en `auto_max_usage` o por encima. Solo se usan cuando ninguna cuenta tiene margen, primero la de menor uso de 5 horas.
3. **Llenas**: una ventana al 100 %, o un uso de 5 horas desconocido. Nunca se eligen.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizaciones.** Las elecciones automáticas se quedan con cuentas como la que dejas: la misma organización de equipo o de empresa, o los planes personales (Pro, Max) entre sí, para que el trabajo de una cuenta de la empresa no continúe en una privada ni al revés. ccam conoce la organización de cada cuenta por su inicio de sesión (`ccam login` la registra; en los inicios de sesión más antiguos se lee de `.claude.json`); una cuenta cuya organización no puede determinar encaja en cualquier sitio. Dos carpetas con la sesión iniciada en la misma cuenta de Claude comparten sus límites, así que una nunca releva a la otra. `--cross-org`, o `cross_org = true` en `config.toml`, quita la regla de la organización.

Cuando no se puede elegir nada, ccam sale con el código 75 (`EX_TEMPFAIL`), para que un script distinga "inténtalo más tarde" de un fallo. Con `--wait`, en cambio, espera hasta que se restablezca la ventana que bloquea a la primera cuenta llena (30 segundos después de la hora de restablecimiento, cuando el endpoint muestra la nueva ventana), imprime `ccam: no account has room; waiting until 15:10` y vuelve a elegir; cuando ninguna cuenta dice cuándo se restablece, vuelve a mirar a los 10 minutos. Si nada se libera antes de que acabe la espera, sale con 75. La consulta sigue las reglas de `ccam usage`: las cifras con menos de 60 segundos se reutilizan, cada cuenta tiene 3 segundos, y una cuenta que no responde, o cuyo inicio de sesión venció, se juzga por sus cifras guardadas. Una ventana guardada cuya hora de restablecimiento ya pasó cuenta como 0 %, así que no se descarta una cuenta cuya ventana se restableció desde entonces. `--auto` nunca renueva tokens. A diferencia de la rotación simple, `--auto` puede elegir la cuenta que usaste por última vez si sigue teniendo el mayor margen.

### --continue

1. ccam busca la última conversación de esta carpeta en todas las cuentas: el `<session>.jsonl` más reciente de `projects/<folder>/` (el nombre de la carpeta es la carpeta de trabajo con cada carácter que no sea una letra o un dígito convertido en `-`, como la nombra Claude Code). Si la sesión dejó una nota (un límite de uso o `/ccam:switch`, ver más abajo) y desde entonces no se escribió nada más reciente, se usa exactamente esa sesión.
2. La cuenta a la que pertenece es la cuenta de la que te cambias, y la cuenta nueva se elige como se describe arriba, dentro de la misma organización.
3. Si la cuenta nueva comparte ese historial (la misma carpeta de grupo), ccam inicia `claude --resume <session id>`, justo esa conversación (`--continue` abriría la más reciente de la carpeta, sea cual sea).
4. Si no, inicia `claude --resume <the conversation's file> --fork-session`: Claude Code copia la conversación completa, con un nuevo id de sesión, en el historial propio de la cuenta nueva, y el original queda como estaba. A partir de entonces, `--continue` de la cuenta nueva la encuentra.

ccam indica el tamaño de la conversación: la cuenta nueva la vuelve a leer entera con su primera respuesta. Las cachés de prompts nunca se comparten entre organizaciones, así que con `--cross-org` esa primera respuesta cuesta el tamaño completo. Compartir `history` en un grupo evita la copia, mantiene una sola conversación y comparte también los puntos de control desde los que `/rewind` restaura el código (`file-history/`).

## Dentro de Claude Code

El [plugin de ccam](/docs/install#claude-code-plugin) agrega dos hooks (Claude Code 2.1.78 o posterior):

- Cuando un turno se detiene en un límite de uso (`You've hit your session limit`, el límite semanal o el de un modelo), anota la conversación y muestra una notificación de escritorio en las terminales que la admiten (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Una limitación breve del servidor o el límite de solicitudes de una clave de API no es un límite de uso y no se toca. Claude Code 2.1.234 y posteriores también esperan en la sesión y continúan por sí solos tras un restablecimiento de 5 horas; la notificación te da la opción de cambiar de cuenta en su lugar.
- `/ccam:switch [account]` anota la conversación y termina el turno antes de que se haga ninguna solicitud: escribe `/exit`, y luego `ccam next --continue` la continúa, con la cuenta que nombraste o con la siguiente.

La nota vive en `~/.local/state/ccam/handoff/<account>.json` durante 12 horas y se usa una sola vez, cuando la siguiente sesión está a punto de iniciarse; un inicio que falla la deja para el siguiente intento. Pertenece a la carpeta en la que se inició la sesión, incluso después de un `cd` en el shell de Claude Code. Si la cuenta que nombra se ha renombrado o eliminado desde entonces, ccam lo indica y elige como de costumbre. Con `relay = true` en `config.toml`, ccam se queda al lado de las sesiones interactivas en lugar de reemplazarse a sí mismo por Claude Code, y después de `/ccam:switch` continúa por sí solo la conversación de esa sesión (no actúa sobre una nota que haya dejado otra sesión). Eso cuesta un proceso de ccam en espera por sesión, y por eso está desactivado por defecto. La sesión que toma el relevo se inicia con `ccam next --continue`, sin los argumentos con los que se inició la primera, como `--model`.

## Volver a ejecutar al llegar al límite

`ccam next` es para una sesión interactiva que se ha detenido. Una ejecución no interactiva, `claude -p` en un script o una tarea programada, no puede esperarte, así que `--fallback` antes del nombre de la cuenta permite que ccam cambie por ella:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` solo funciona junto con `-p` o `--print`; si no, ccam se detiene y remite a `ccam next --auto --continue`.
2. La entrada estándar se lee una vez y se entrega a cada intento.
3. ccam inicia `claude` como proceso hijo en lugar de con `exec`, retiene su salida estándar y deja pasar la salida de error estándar a medida que llega.
4. Cuando `claude` sale con 0, o falla por otro motivo, la salida retenida y el código de salida se transmiten sin cambios.
5. Un fallo cuenta como límite de uso cuando el mensaje de Claude Code lo indica (`You've hit your limit`, `usage limit reached`), o cuando la consulta de uso de la cuenta muestra su ventana de 5 horas o semanal al 100 %.
6. Entonces se descarta la salida de ese intento y el mismo comando se ejecuta con otra cuenta con sesión iniciada: la de menor uso de 5 horas, y después las cuentas cuyo uso no se puede consultar, por orden de número. Las cuentas al 100 % se omiten, y cada cuenta se prueba como máximo una vez.
7. Con `-c`, `--continue`, `-r` o `--resume` entre los argumentos, solo se prueban las cuentas del mismo grupo de historial, para que se pueda encontrar la conversación.
8. Solo se prueban cuentas de la misma organización, igual que con `--auto` (`cross_org = true` quita esta restricción).

La siguiente cuenta se elige con las mismas reglas que `--auto` (`auto_max_usage`, `auto_strategy`, la ventana de `--model`). Cuando no queda ninguna cuenta, ccam imprime la salida del último intento, indica que ninguna cuenta con sesión iniciada tiene margen y sale con el código 75. Con `--wait` antes del nombre de la cuenta (`ccam --fallback --wait=2h work -p …`), en cambio, espera a que se restablezca la primera cuenta y después puede volver a probar todas las cuentas; siempre espera al menos una vez antes, para que una cuenta que acaba de llegar a su límite no se vuelva a ejecutar de inmediato. `ccam --dry-run --fallback work -p …` imprime el comando del primer intento.

Las sesiones interactivas no cambian por sí solas. En su lugar, con `ccam statusline install`, la línea de estado añade `→ ccam next --auto --continue` cuando el uso de 5 horas o semanal llega al 90 % ([Límites de uso](/features/usage)).

## Seguridad y casos límite

- Las cuentas sin inicio de sesión nunca se eligen. Si no hay otra cuenta con sesión iniciada, `ccam next` a secas se detiene con `no logged-in account other than <name>`.
- Con `--continue` y ninguna conversación en la carpeta: `no account has a conversation in <folder> to continue`.
- Cuando la regla de la organización no deja ninguna cuenta: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Con `--auto` y ninguna cuenta con sesión iniciada: `no logged-in account to pick`. Cuando todas las cuentas están llenas o su uso es desconocido: `` no account with known usage and headroom (see `ccam usage`) ``. Ambos casos terminan con el código de salida 75.
- `--max-usage`, `--strategy`, `--wait` y `--cross-org` solo funcionan con `--auto`, y se rechaza un valor fuera de 1-100 u otra estrategia.
- Una bifurcación copia la conversación, no el resto del estado de la cuenta: los puntos de control de `file-history/` de la cuenta original no la acompañan, salvo que el historial se comparta.
- `--auto` hace solicitudes de red incluso con `--dry-run`. Solo lee, y nunca escribe archivos de tokens.
- La cuenta elegida se inicia por la misma ruta que `ccam <name>`, así que se aplican las mismas comprobaciones, incluida la advertencia cuando otro Claude Code que comparte el mismo historial ya se está ejecutando en esta carpeta.
- `last-used` solo se escribe cuando Claude Code se inicia de verdad. `--dry-run` no lo toca, así que una ejecución de prueba no hace avanzar la rotación.
- `--fallback` retiene la salida estándar hasta que termina un intento, así que una salida continua como `--output-format stream-json` llega de una sola vez al final. Los archivos que el primer intento cambió antes de llegar al límite siguen cambiados; el siguiente intento empieza el prompt de nuevo.

## Relacionado

- [Límites de uso](/features/usage)
- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Flujo de trabajo con límites de uso](/guides/usage-limit-workflow)
