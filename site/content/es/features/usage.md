---
title: Límites de uso
description: ccam usage muestra a la vez los límites de 5 horas y semanal de todas las cuentas, y ccam statusline pone el nombre de la cuenta y su uso en Claude Code.
updated: 2026-10-19
order: 7
translation: machine
sourceHash: 8306a8f6d950
---

`ccam usage` muestra cuánto ha usado cada cuenta de sus límites de 5 horas y semanal, junto con el tiempo que falta para que se restablezca cada ventana. Úsalo antes de empezar una tarea larga, o cuando Claude Code te avise de que un límite está cerca, para ver qué cuenta tiene margen. `ccam statusline` es el comando de línea de estado para Claude Code: muestra el nombre de la cuenta y su uso dentro de Claude Code y guarda las cifras para `ccam usage`.

## Uso

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

En una terminal, la tabla aparece de inmediato con `fetching…` en cada fila, y cada fila se vuelve a dibujar cuando responde su cuenta. `As of` indica la antigüedad de las cifras; `(cached)` significa que vienen de un archivo guardado, no del endpoint. La última línea nombra la cuenta con el menor uso de 5 horas por debajo del 100 %, y solo aparece si tienes más de una cuenta. Una cuenta sin sesión iniciada muestra `not logged in → ccam login <name>`, y una cuenta con clave de API muestra `no limits (API)`.

Después de `ccam statusline install`, Claude Code muestra una línea como `[work] · Opus · 5h 24% · 7d 41%`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--json` | Imprime los informes en JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | desactivada |
| `--watch[=N]` | Vuelve a dibujar cada `N` segundos hasta que pulses Ctrl+C. `--watch` solo, o un valor de 0 o menos, significa 60 | desactivada |
| `-v`, `--verbose` | Muestra los límites semanales por modelo y una línea `source:` (`live`, `cache` o `statusline`) | desactivada |
| `--no-cache` | Ignora la caché de 60 segundos y vuelve a consultar el endpoint | desactivada |
| `--source live\|cache` | `cache` nunca consulta el endpoint y solo muestra las cifras guardadas; `live` es el comportamiento predeterminado | live |
| `--refresh` | Renueva los tokens de acceso vencidos, solo de las cuentas que no se están ejecutando | desactivada |
| `--dry-run` | Se acepta: `ccam usage` no cambia nada salvo su propia caché. Se rechaza junto con `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | desactivada |

Cualquier otro valor de `--source` falla con `--source must be live or cache`. `ccam usage` no acepta argumentos.

`ccam statusline` lee de la entrada estándar el JSON de la línea de estado de Claude Code. Su única opción es `--wrap <command>`: un comando de línea de estado existente a cuya primera línea de salida se le antepone el nombre de la cuenta. `ccam statusline install` no acepta opciones. Ambos rechazan `--dry-run` con `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Cómo funciona

Las cifras salen de dos lugares: el endpoint de uso no documentado de Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, al que se llama con el inicio de sesión propio de cada cuenta, y los valores que Claude Code le pasa a su comando de línea de estado. Todo esto está en `internal/usage`.

1. `cli.usageCmd` carga `config.toml` y crea un `usage.Account` por cuenta: su carpeta `~/.local/share/ccam/profiles/<name>/`, su carpeta de sesiones (para saber si está en ejecución) y si se importó en su ubicación.
2. `Fetcher.FetchAll` inicia una goroutine por cuenta. Todas comparten un `http.Client`, así que las conexiones al endpoint se reutilizan. La conexión puede tardar 2 s, y cada cuenta tiene 3 s en total. El cliente no sigue redirecciones, ya que una redirección reenviaría el token a otro sitio, y lee como máximo 1 MB de la respuesta.
3. A medida que termina cada cuenta, su fila se vuelve a dibujar. Una cuenta que falla o no responde a tiempo recurre a las cifras guardadas, así que una cuenta lenta nunca retrasa a las demás.

Cada cuenta consulta el endpoint en su propia goroutine; la que agota el tiempo muestra en su lugar sus cifras guardadas.

<!-- motion:usage-fetch -->

`personal` y `work` responden a tiempo y muestran cifras en vivo. `side` tarda más de 3 s, así que `fallback()` lee sus dos archivos de caché y usa el que se escribió por última vez. `usage.Table()` dibuja las filas y `usage.Best()` elige la cuenta que se nombra en la última línea.

Para una cuenta, `Fetcher.Fetch` hace estas comprobaciones en orden:

<!-- motion:usage-cache -->

1. Con `--source cache`, va directamente a `fallback()`.
2. Salvo que se indique `--no-cache`, un resultado en vivo en `~/.local/state/ccam/usage/<name>.json` con menos de 60 segundos (`CacheTTL`) se devuelve tal cual. Ejecutar `ccam usage` dos veces seguidas cuesta una solicitud por cuenta.
3. `readCredentials()` lee `.credentials.json` (o el Llavero de macOS) sin cambiarlo, y solo se decodifica la clave `claudeAiOauth`. Si no hay un inicio de sesión de suscripción, se muestra `not logged in`, o `no limits (API)` cuando `.claude.json` contiene una clave de API.
4. Un token de acceso vencido no se usa. Sin `--refresh`, la cuenta muestra sus cifras guardadas marcadas como vencidas.
5. En caso contrario, ccam llama al endpoint. Un 401 o un 403 cuenta como inicio de sesión vencido, cualquier otro estado o un cuerpo que no puede leer cuenta como error, y en ambos casos se recurre a las cifras guardadas.
6. Una respuesta correcta se escribe en `<name>.json` con la hora en que se midió.

`fallback()` lee tanto `<name>.json` como `<name>.statusline.json` y se queda con el más reciente. Antes de mostrar cifras guardadas, comprueba la hora de restablecimiento de cada ventana: una ventana que se restableció después de guardar las cifras muestra 0 %, porque el porcentaje anterior ya no es válido. Por eso `personal`, más arriba, muestra `0%` sin hora de restablecimiento.

El endpoint tiene dos formatos de respuesta. El más nuevo tiene un arreglo `limits` con entradas `session` (5 horas), `weekly_all` y `weekly_scoped` (por modelo). El más antiguo tiene los objetos `five_hour`, `seven_day`, `seven_day_opus` y `seven_day_sonnet`. `ParseLive` lee ambos: el arreglo `limits` tiene prioridad, y los objetos antiguos completan lo que no cubre. Un cuerpo sin ningún límite se informa como un cambio de formato, del que te avisa `ccam doctor`.

### Renovar tokens

Los tokens de actualización rotan: al renovar uno obtienes un nuevo token de actualización y el anterior queda invalidado. Si ccam renovara el inicio de sesión mientras Claude Code se está ejecutando con la misma cuenta, esa sesión quedaría cerrada. Por eso `ccam usage` lee los tokens y nunca los escribe, salvo que pases `--refresh`. Entonces `Fetcher.refresh`:

1. rechaza las cuentas registradas en su ubicación con `ccam import`, ya que `claude` a secas puede usar ese inicio de sesión sin que ccam lo sepa;
2. rechaza los inicios de sesión guardados en el Llavero de macOS;
3. toma un archivo de bloqueo, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (dentro de `~/.local/state/ccam/locks/` cuando esa variable no está definida), para que dos procesos de ccam nunca renueven dos veces el mismo token;
4. con el bloqueo tomado, se niega si hay una sesión de la cuenta en ejecución, y omite la solicitud si otro proceso ya renovó el token;
5. envía el token de actualización al endpoint de tokens y vuelve a escribir los tokens nuevos en `.credentials.json` con `fsutil.AtomicWrite`, conservando todas las demás claves del archivo.

### El recolector de la línea de estado

Claude Code ejecuta el comando de línea de estado y le pasa un documento JSON por una tubería. ccam lee de él el nombre del modelo y `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` identifica la cuenta a partir de `CCAM_PROFILE`, que ccam define al iniciar Claude Code, o si no, comparando `CLAUDE_CONFIG_DIR` con las carpetas de las cuentas.
2. `ParseStatusLine` lee `rate_limits.five_hour` y `rate_limits.seven_day`. El JSON no tiene límites de uso en las cuentas con clave de API, ni antes de la primera respuesta de la sesión; en esos casos no se guarda nada.
3. `SaveStatusLine` escribe `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` lo usa cuando el endpoint no responde.
4. ccam imprime `[name] · model · 5h N% · 7d N%`. Con `--wrap`, ejecuta tu comando mediante `sh -c` con la misma entrada e imprime en su lugar `[name]` seguido de la primera línea del comando.

El comando envuelto tiene 2 segundos. Pasado ese tiempo se termina y la línea muestra solo el nombre de la cuenta, así que un comando lento no puede bloquear la línea de estado de Claude Code.

`ccam statusline install` escribe `statusLine` en el `settings.json` de cada cuenta. Cuando una cuenta ya tiene un comando de línea de estado, ccam lo conserva como `--wrap '<old command>'`. Un `settings.json` compartido a través de un grupo se escribe una sola vez, a través de su enlace, así que todas las cuentas del grupo lo reciben. Al volver a ejecutarlo, imprime `(already set)` para los archivos que ya usan ccam.

## Seguridad y casos límite

- `ccam usage` nunca cambia un archivo de tokens salvo que pases `--refresh`. Cuando un inicio de sesión venció, inicia la cuenta una vez con `ccam <name>` y Claude Code lo renueva.
- `--refresh` nunca toca las cuentas en ejecución, las cuentas importadas ni los inicios de sesión del Llavero. Una renovación rechazada muestra las cifras guardadas y el motivo.
- El endpoint de uso no está documentado públicamente y puede cambiar. Si cambia, sigues viendo las cifras guardadas, y `ccam doctor` informa el cambio de formato.
- Sin cifras guardadas y sin respuesta, la fila muestra `lookup failed:` y el motivo, por ejemplo `no saved usage` con `--source cache`.
- Los archivos de caché solo contienen porcentajes y horas de restablecimiento, nunca tokens.

## Relacionado

- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Cambiar a la siguiente cuenta](/features/next)
- [Flujo de trabajo con límites de uso](/guides/usage-limit-workflow)
- [Comandos](/docs/commands)
