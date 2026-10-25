---
title: Límites de uso
description: ccam usage muestra a la vez los límites de 5 horas y semanal de todas las cuentas, y ccam statusline pone el nombre de la cuenta y su uso en Claude Code.
updated: 2026-10-25
order: 7
translation: machine
sourceHash: "91ad20327289"
---

`ccam usage` muestra cuánto ha usado cada cuenta de su límite de 5 horas, su límite semanal y el límite semanal propio de cada modelo (como Fable), con una barra por límite y el tiempo que falta para que se restablezca. Úsalo antes de empezar una tarea larga, o cuando Claude Code te avise de que un límite está cerca, para ver qué cuenta tiene margen. `ccam statusline` es el comando de línea de estado para Claude Code: muestra el nombre de la cuenta y su uso dentro de Claude Code y guarda las cifras para `ccam usage`.

## Uso

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Cada cuenta tiene una línea con su número, su nombre y la antigüedad de sus cifras; `(cached)` significa que vienen de un archivo guardado, no del endpoint. Debajo, una línea por límite muestra una barra de 20 celdas, el porcentaje y el tiempo hasta que se restablece la ventana: `5-hour`, `Weekly` (todos los modelos) y luego la ventana semanal propia de cada modelo, como `Fable`. Una ventana que la cuenta no tiene muestra `-`. Un porcentaje se vuelve amarillo a partir del 80 % y rojo al 100 %. Una ventana que Claude Code informó llena, cuando rechazó una solicitud en un límite de uso, muestra 100 % hasta su restablecimiento, digan lo que digan las cifras guardadas (consulta [el registro de límite](#the-limit-record)). En una terminal todas las cuentas aparecen de inmediato como `fetching…`, y sus líneas se dibujan a medida que responden. La última línea nombra la cuenta que elegiría [`ccam next --auto`](/features/next) con las mismas reglas (`auto_max_usage`, `auto_strategy`, las sesiones en ejecución, solo cuentas de la organización de la última cuenta usada salvo que `cross_org` esté activado, se omite toda ventana al 100 %, así que nunca se nombra una cuenta con la semana llena), y solo aparece si tienes más de una cuenta. Con la estrategia predeterminada `headroom` dice `→ Most headroom: <name>`; con `expiring` o `balance` da el motivo, como `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Asientos de Team y Enterprise.** La elección automática (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y [`auto_switch`](/features/next#auto-switch) cuando lo activas) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por lo tanto, usar la elección automática entre asientos puede infringir los términos de Anthropic. Si la usas ahí es decisión tuya: consulta con tu organización y con Anthropic antes de depender de ella. Los planes personales (Pro, Max) se tratan en la [nota de las preguntas frecuentes sobre los términos](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Una cuenta sin sesión iniciada muestra `not logged in → ccam login <name>`, y una cuenta con clave de API muestra `no limits (API)`. Un token de acceso vencido muestra ``access token expired; `ccam <name>` renews it when it starts``, encima de las cifras guardadas cuando las hay: Claude Code lo renueva en el siguiente inicio. Un inicio de sesión que venció por completo, porque su token de actualización venció o Claude Code lo borró después de que el servidor lo rechazara, muestra `` login expired → `ccam login <name>` ``.

Después de `ccam statusline install`, Claude Code muestra una línea como `[work] · Opus · 5h 24% · 7d 41%`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--json` | Imprime los informes en JSON: `name`, `slot`, `status`, `limits` (con `credits` cuando la cuenta tiene activados los créditos de uso), `source`, `at`, `error`, y `limited` y `limited_windows` cuando un registro de límite pone ventanas al 100 %. Sin cuentas, imprime `[]` | desactivada |
| `--watch[=N]` | Vuelve a dibujar cada `N` segundos hasta que pulses Ctrl+C. `--watch` solo, o un valor de 0 o menos, significa 60. El valor va después de `=`: `--watch 10` se rechaza con una sugerencia | desactivada |
| `-v`, `--verbose` | Añade una línea `source:` bajo cada cuenta: de dónde salieron sus cifras (el endpoint de uso, Claude Code, la línea de estado, un archivo guardado o el límite de uso que informó Claude Code). Cuando un registro de límite se impone sobre las cifras, sigue una línea como `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | desactivada |
| `--no-cache` | Ignora la caché de 60 segundos y vuelve a consultar el endpoint | desactivada |
| `--source live\|claude\|cache` | `live` consulta el endpoint con el token de la cuenta; `claude` hace que Claude Code lo consulte (ver más abajo); `cache` solo muestra las cifras guardadas | `usage_source`, si no `live` |
| `--refresh` | Renueva los tokens de acceso vencidos, solo de las cuentas que no se están ejecutando. Una renovación puede tardar hasta 30 segundos | desactivada |
| `--dry-run` | Se acepta: `ccam usage` no cambia nada salvo su propia caché. Se rechaza junto con `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | desactivada |

Cualquier otro valor de `--source` falla con `--source must be live, claude or cache`.

### El uso desde el propio Claude Code

Con `usage_source = "claude"` al principio de `config.toml` (antes del primer `[[group]]`, `[[profile]]` o `[[folder]]`), o `--source claude`, ccam no lee en absoluto el token de inicio de sesión. Para cada cuenta inicia `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` en la carpeta de la cuenta y envía por la entrada estándar la solicitud `get_usage` de Claude Code (Claude Code 2.1.270 o posterior); Claude Code consulta el uso con su propio inicio de sesión y responde con las ventanas de 5 horas, semanal y por modelo. `--safe-mode` impide que se inicien los hooks, los servidores MCP y los plugins de la cuenta, y no se guarda ninguna sesión. Cada cuenta supone un inicio de Claude Code, de alrededor de un segundo, así que esta fuente es más lenta que `live`; las cuentas se consultan a la vez, cada una con un límite de 20 segundos, y la caché de 60 segundos se aplica como siempre. Un inicio de sesión cuyo token de acceso venció, o vence en menos de 5 minutos, recibe 50 segundos, ya que Claude Code lo renueva al iniciar y una renovación cortada puede perder el inicio de sesión. Una consulta abandonada le pide a Claude Code que termine con SIGTERM, para que acabe una renovación que haya empezado (Windows no tiene esa señal; allí se detiene). Un inicio de sesión que venció por completo muestra `` login expired → `ccam login <name>` `` sin iniciar Claude Code. Claude Code marca la solicitud como experimental, así que si deja de responder, ccam muestra las cifras guardadas y `ccam doctor` lo informa. Cuando Claude Code responde que su propia consulta falló, ccam muestra las cifras guardadas con `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, no un plan sin límites. Un `--dry-run` no inicia Claude Code (escribe en la carpeta de la cuenta) y muestra las cifras guardadas. `usage_source` acepta `live` o `claude`; `cache` solo es una opción de `--source`. `ccam usage` no acepta argumentos.

`ccam statusline` lee de la entrada estándar el JSON de la línea de estado de Claude Code. Su única opción es `--wrap <command>`: un comando de línea de estado existente a cuya primera línea de salida se le antepone el nombre de la cuenta. `ccam statusline install` no acepta opciones. Ambos rechazan `--dry-run` con `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Cómo funciona

Las cifras salen de dos lugares: el endpoint de uso no documentado de Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, al que se llama con el inicio de sesión propio de cada cuenta, y los valores que Claude Code le pasa a su comando de línea de estado. Todo esto está en `internal/usage`.

1. `cli.usageCmd` carga `config.toml` y crea un `usage.Account` por cuenta: su carpeta `~/.local/share/ccam/profiles/<name>/`, sus marcadores de sesión (junto con los registros propios de Claude Code en la carpeta de la cuenta, indican si está en ejecución) y si se importó en su ubicación.
2. `Fetcher.FetchAll` inicia una goroutine por cuenta. Todas comparten un `http.Client`, así que las conexiones al endpoint se reutilizan. La conexión puede tardar 2 s, y cada cuenta tiene 3 s en total. El cliente no sigue redirecciones, ya que una redirección reenviaría el token a otro sitio, y lee como máximo 1 MB de la respuesta.
3. A medida que termina cada cuenta, su fila se vuelve a dibujar. Una cuenta que falla o no responde a tiempo recurre a las cifras guardadas, así que una cuenta lenta nunca retrasa a las demás.

Cada cuenta consulta el endpoint en su propia goroutine; la que agota el tiempo muestra en su lugar sus cifras guardadas.

<!-- motion:usage-fetch -->

`personal` y `work` responden a tiempo y muestran cifras en vivo. `side` tarda más de 3 s, así que `fallback()` muestra en su lugar sus cifras guardadas. `usage.Table()` dibuja un bloque por cuenta y `usage.Best()` elige la cuenta que se nombra en la última línea.

Para una cuenta, `Fetcher.Fetch` hace estas comprobaciones en orden:

<!-- motion:usage-cache -->

1. Con `--source cache`, va directamente a `fallback()`.
2. Salvo que se indique `--no-cache`, un resultado en vivo en `~/.local/state/ccam/usage/<name>.json` medido hace menos de 60 segundos (`CacheTTL`) se usa, evaluado a la hora actual: una ventana que se restableció desde entonces muestra 0 %, y un archivo con fecha posterior a la del reloj no cuenta como reciente. Ejecutar `ccam usage` dos veces seguidas cuesta una solicitud por cuenta.
3. `readCredentials()` lee `.credentials.json` (o el Llavero de macOS) sin cambiarlo, y solo se decodifica la clave `claudeAiOauth`. Si no hay un inicio de sesión de suscripción, se muestra `not logged in`, o `no limits (API)` cuando `.claude.json` contiene una clave de API. Un inicio de sesión cuyo token de actualización venció, o que Claude Code borró, es `` login expired → `ccam login <name>` `` y no se consulta.
4. Un token de acceso vencido no se usa. Sin `--refresh`, la cuenta muestra sus cifras guardadas con ``access token expired; `ccam <name>` renews it when it starts``.
5. En caso contrario, ccam llama al endpoint. Un 401 o un 403 cuenta como inicio de sesión vencido (`login expired: authentication failed (401)`), cualquier otro estado o un cuerpo que no puede leer cuenta como error, y en ambos casos se recurre a las cifras guardadas. Una consulta terminada con Ctrl+C dice `cancelled`. Un 429 significa que se está consultando el endpoint con demasiada frecuencia: ccam escribe `<name>.pause` con la hora para volver a preguntar (`Retry-After`, o 5 minutos, como máximo una hora), muestra las cifras guardadas con `rate limited by the usage endpoint; asking again after 15:10` y no llama al endpoint para esa cuenta hasta entonces, ni siquiera con `--no-cache`. Reintentar de inmediato solo mantendría el límite.
6. Una respuesta correcta se escribe en `<name>.json` con la hora en que se midió.

`fallback()` toma la consulta guardada más reciente (`<name>.json`, o `<name>.claude.json` de las consultas de Claude Code) y, cuando `<name>.statusline.json` se escribió después, pone encima las ventanas de 5 horas y semanal de la línea de estado; las ventanas por modelo, que la línea de estado no tiene, se quedan como las guardó la consulta. Antes de mostrar cifras guardadas, comprueba la hora de restablecimiento de cada ventana: una ventana que se restableció después de guardar las cifras muestra 0 %, porque el porcentaje anterior ya no es válido. Por eso `personal`, más arriba, muestra `0%` sin hora de restablecimiento.

El endpoint tiene dos formatos de respuesta. El más nuevo tiene un arreglo `limits` con entradas `session` (5 horas), `weekly_all` y `weekly_scoped` (por modelo). El más antiguo tiene los objetos `five_hour`, `seven_day`, `seven_day_opus` y `seven_day_sonnet`. `ParseLive` lee ambos: el arreglo `limits` tiene prioridad, y los objetos antiguos completan lo que no cubre. Un cuerpo sin ningún límite se informa como un cambio de formato, del que te avisa `ccam doctor`.

### El registro de límite

Cuando Claude Code detiene un turno en un límite de uso, el hook del [plugin de ccam](/docs/install#claude-code-plugin) lee de la conversación la solicitud rechazada: qué ventana estaba llena (5 horas, semanal, Opus o Sonnet) y cuándo se restablece. Lo escribe en `~/.local/state/ccam/usage/<name>.limit.json`. `Fetch` pone el registro encima de todo lo demás que sabe, de cualquier fuente: cada ventana que contiene cuenta como 100 % hasta su restablecimiento, por antiguas o bajas que sean las cifras guardadas, así que `ccam next --auto`, `ccam pick`, `--fallback` y [`auto_switch`](/features/next#auto-switch) no vuelven a elegir la cuenta antes de entonces. Si no se sabe nada más, el registro es lo único que muestra la cuenta (`source: the usage limit Claude Code reported` con `-v`). Las sesiones que nadie atiende (`-p`, sesiones en segundo plano) no escriben ningún registro, ni tampoco la detención por parte de una organización (sin uso disponible, un grupo fijado en $0), que sus administradores pueden levantar en cualquier momento. `ccam login` en otra cuenta de Claude descarta el registro; eliminar el archivo también lo borra.

### Renovar tokens

Los tokens de actualización rotan: al renovar uno obtienes un nuevo token de actualización y el anterior queda invalidado. Si ccam renovara el inicio de sesión mientras Claude Code se está ejecutando con la misma cuenta, esa sesión quedaría cerrada. Por eso `ccam usage` lee los tokens y nunca los escribe, salvo que pases `--refresh`. Entonces `Fetcher.refresh`:

1. rechaza las cuentas registradas en su ubicación con `ccam import`, ya que `claude` a secas puede usar ese inicio de sesión sin que ccam lo sepa;
2. rechaza los inicios de sesión guardados en el Llavero de macOS;
3. toma un archivo de bloqueo, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, para que dos procesos de ccam nunca renueven dos veces el mismo token; mientras está tomado, la cuenta muestra `token not renewed: another ccam is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` y `logout` también lo toman, y se niegan mientras hay una renovación en curso;
4. con el bloqueo tomado, se niega si hay una sesión de la cuenta en ejecución, sin importar cómo se haya iniciado (mediante ccam, `claude` a secas, un IDE, o una sesión en segundo plano y el daemon que la respalda, que renueva el inicio de sesión por sí mismo), y omite la solicitud si otro proceso ya renovó el token;
5. no envía nada mientras Claude Code mantiene sus propios bloqueos de renovación en la carpeta de la cuenta (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code está renovando ese inicio de sesión en ese momento;
6. envía el token de actualización al endpoint de tokens y espera hasta 30 segundos la respuesta, el tiempo que Claude Code da a sus propias renovaciones. Una solicitud ya enviada nunca se corta, ni siquiera con Ctrl+C, ya que puede que el servidor ya haya rotado el token de actualización y solo su respuesta conserve el inicio de sesión: Ctrl+C o un hangup terminan las demás consultas, y ccam sale (130 tras Ctrl+C) una vez guardada la respuesta;
7. guarda con el bloqueo de escritura de Claude Code (`.storage-write.lock`) y cambia solo los campos de token de `.credentials.json` tal como esté en ese momento: lo que Claude Code escribió mientras tanto (un inicio de sesión MCP) se conserva, un inicio de sesión que Claude Code renovó mientras tanto se mantiene, y un cierre de sesión nunca se deshace. Cuando la respuesta indica la duración del nuevo token de actualización, se guarda como `refreshTokenExpiresAt`, igual que hace Claude Code.

Una renovación que no se completa indica el motivo, en la línea de la cuenta o encima de sus cifras guardadas: `token not renewed: a session is running`, por ejemplo, o `login may be lost: …` cuando se envió una solicitud pero su respuesta no se pudo usar ni guardar, con el consejo de ejecutar `ccam login <name>` si Claude Code te pide luego que inicies sesión.

### El recolector de la línea de estado

Claude Code ejecuta el comando de línea de estado y le pasa un documento JSON por una tubería. ccam lee de él el nombre del modelo y `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` identifica la cuenta a partir de `CCAM_PROFILE`, que ccam define al iniciar Claude Code, o si no, comparando `CLAUDE_CONFIG_DIR` con las carpetas de las cuentas.
2. `ParseStatusLine` lee `rate_limits.five_hour` y `rate_limits.seven_day`. El JSON no tiene límites de uso en las cuentas con clave de API, ni antes de la primera respuesta de la sesión; en esos casos no se guarda nada.
3. `SaveStatusLine` escribe `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` lo usa cuando el endpoint no responde.
4. ccam imprime `[name] · model · 5h N% · 7d N%`. Con `--wrap`, ejecuta tu comando mediante `sh -c` con la misma entrada e imprime su salida como Claude Code muestra una línea de estado: cada línea que no está en blanco, con `[name]` delante de la primera. Un comando que falla muestra solo el nombre de la cuenta. En Windows ejecuta el comando como lo hace Claude Code: con Git Bash si está instalado y, si no, con PowerShell.
5. Cuando el uso de 5 horas o semanal está al 90 % o más y hay otra cuenta, la línea termina con `→ ccam next --auto --continue`, por ejemplo `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

El comando envuelto no tiene un límite de tiempo propio: Claude Code da 10 minutos a un comando de línea de estado y detiene una ejecución cuando empieza una actualización más reciente, así que un comando lento se muestra igual que sin ccam. Cuando Claude Code termina ccam, el comando también termina (en Linux y macOS, con todos los procesos que inició).

`ccam statusline install` escribe `statusLine` en el `settings.json` de cada cuenta y conserva las demás claves del objeto, como `padding`, `refreshInterval` y `hideVimModeIndicator`. Cuando una cuenta ya tiene un comando de línea de estado, ccam lo conserva como `--wrap '<old command>'`; en Windows, como `--wrap-base64 <base64>`, que no necesita comillas en ninguna de las dos shells. Una línea de estado que ya ejecuta un ccam, desde otra ruta o entre comillas como lo escribía un ccam anterior, se reemplaza en lugar de envolverse, conservando lo que envuelve. En Windows, la ruta a ccam se escribe con barras normales y `~` para la carpeta personal (`~/.local/bin/ccam.exe statusline`), que entienden tanto Git Bash como PowerShell. Una ruta que necesita comillas se escribe completa para la shell que usa Claude Code: `'…'` para Git Bash, `& '…'` para PowerShell. Un `settings.json` compartido a través de un grupo se escribe una sola vez, a través de su enlace, así que todas las cuentas del grupo lo reciben. Al volver a ejecutarlo, imprime `(already set)` para los archivos que ya usan ccam.

## Seguridad y casos límite

- `ccam usage` nunca cambia un archivo de tokens salvo que pases `--refresh`. Cuando un token de acceso venció, inicia la cuenta una vez con `ccam <name>` y Claude Code lo renueva. Un inicio de sesión que venció por completo necesita `ccam login <name>`.
- `--refresh` nunca toca las cuentas en ejecución (también cuando Claude Code se inició sin ccam), las cuentas importadas ni los inicios de sesión del Llavero. Una renovación rechazada muestra las cifras guardadas y el motivo (`token not renewed: …`).
- El endpoint de uso no está documentado públicamente y puede cambiar. Si cambia, sigues viendo las cifras guardadas, y `ccam doctor` informa el cambio de formato.
- Sin cifras guardadas y sin respuesta, la línea de la cuenta muestra `lookup failed:` y el motivo, por ejemplo `no saved usage` con `--source cache`.
- Los archivos de caché solo contienen porcentajes y horas de restablecimiento, nunca tokens.

## Relacionado

- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Cambiar a la siguiente cuenta](/features/next)
- [Flujo de trabajo con límites de uso](/guides/usage-limit-workflow)
- [Comandos](/docs/commands)
