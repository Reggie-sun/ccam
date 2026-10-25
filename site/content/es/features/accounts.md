---
title: Administrar cuentas
description: Agrega, inicia sesión, lista, renombra, quita e importa cuentas de Claude Code con ccam add, ccam login, ccam list, ccam rename, ccam remove y ccam import.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "1d40ce0280e0"
---

Cada cuenta de ccam es una carpeta de configuración de Claude Code independiente con su propio inicio de sesión. Estos comandos crean cuentas, inician sesión en ellas, las listan, les cambian el nombre y las quitan, y registran carpetas que ya tienes.

## Uso

```sh
ccam add                                   # wizard on a terminal: name, sign-in, groups, shortcut
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --group default --group team
ccam login side                            # run claude auth login inside the account
ccam login --expired                       # every account never logged in, or whose login expired
ccam logout side                           # log out: the login and its MCP sign-ins
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, logs out, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Siempre que un comando acepta `<name|number>`, puedes indicar el nombre de la cuenta o su número de posición, la columna `#` de `ccam list`.

`ccam add` sin opciones en una terminal abre un asistente. Pide el nombre y cómo inicia sesión la cuenta (un inicio de sesión de Claude ahora, una clave de API o un inicio de sesión de Claude más tarde), luego muestra los grupos a los que unirse con `default` marcado (si todavía no hay grupos, los elementos que compartirá el grupo `default`), después pregunta `Create a shortcut command for work? [y/N]` y termina con un resumen para confirmar. `Esc` vuelve un paso atrás. Cualquiera de `--group`, `--no-group`, `--shortcut` o `--no-login`, o una tubería en lugar de una terminal, omite el asistente; en ese caso el nombre es obligatorio y no se crea ningún atajo salvo que `--shortcut` indique uno.

Sin `--group`, una cuenta nueva se une al grupo `default`. La primera cuenta lo crea, y comparte los servidores MCP, la configuración, `CLAUDE.md` y las extensiones; consulta [compartir con grupos](/features/share). En Windows sin el modo de desarrollador, sin el cual no se pueden enlazar archivos, la primera cuenta crea `default` solo con `mcp` y `extensions` y lo indica (`Note: files can be linked on Windows only in Developer Mode (Settings > System > For developers), so group default shares only mcp, extensions. …`), y se rechaza `--group` hacia un grupo que comparte la configuración, la memoria o el historial.

`ccam list` muestra la posición, el nombre, el atajo, el correo electrónico, los grupos y el último uso, con `*` delante de la cuenta predeterminada. El correo electrónico es aquel con el que la carpeta tiene la sesión iniciada ahora (de su `.claude.json`); cuando difiere del que registró ccam, el registrado aparece después como `(was …)`:

```console
$ ccam list
   #  Name      Shortcut  Email             Groups        Last used
*  1  work      ccwork    user@example.com  default team  -
   2  personal  -         me@side.dev       default       -
```

```console
$ ccam list
   #  Name      Shortcut  Email                                     Groups   Last used
*  1  work      ccwork    other@company.com (was work@example.com)  default  2m ago
```

Con `--json`, cada cuenta es un objeto con `slot`, `shortcut`, `name`, `email` (el del inicio de sesión actual, como en la tabla), `default`, `dir`, `last_used` (se omite si la cuenta nunca se inició), `groups` y `share`. `share` asigna a cada elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) el grupo del que proviene, o `"local"`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `ccam add --no-login` | Crea la cuenta sin iniciar sesión; inicia sesión más tarde con `ccam login <name>` | desactivada |
| `ccam add --group <group>` | Se une a este grupo en lugar de `default`. Se puede repetir; `default` se crea si no existe | `default` |
| `ccam add --no-group` | No se une a ningún grupo: la cuenta no comparte nada | desactivada |
| `ccam add --shortcut <command>` | Crea también este atajo en `~/.local/bin` | ninguno |
| `ccam add --api --key-from <source>` | Una cuenta que usa una clave de API en lugar de un inicio de sesión; consulta [Cuentas con clave de API](/features/api) | desactivada |
| `ccam list --json` | Imprime las cuentas en JSON | desactivada |
| `ccam remove -y`, `--yes` | Quita sin preguntar | desactivada |
| `ccam logout -y`, `--yes` | Cierra la sesión sin preguntar | desactivada |
| `ccam import --name <name>` | Nombre de cuenta para la carpeta que indicas; se rechaza sin una carpeta (`--name needs a folder: ccam import <folder> --name ops`) | se deriva del nombre de la carpeta |
| `ccam import -y`, `--yes` | Importa todas las carpetas encontradas sin preguntar | desactivada |
| `--dry-run` | La respetan `ccam list` y `ccam completion`, que no cambian nada. La rechazan `add`, `login`, `logout`, `use`, `rename`, `remove` e `import`, con `` --dry-run is not supported by `ccam add`; nothing was changed `` (con el nombre del comando) | desactivada |

`ccam login` (aparte de `--expired`), `ccam use` y `ccam rename` no tienen opciones propias. `ccam completion` acepta `bash`, `zsh` o `powershell`; cualquier otro valor falla con `only bash, zsh and powershell are supported`.

Los nombres de cuenta tienen de 1 a 32 caracteres entre `a-z`, `0-9`, `-` y `_`, y empiezan con una letra o un dígito. Un nombre no puede tener solo dígitos, porque `ccam 2` significa la posición 2. Tampoco puede ser un comando de ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`, `group`, `groups`, `hook`, `agents`, `attach`, `ccusage`), así que `ccam <name>` nunca es ambiguo. La lista es `config.Reserved` en `internal/config`.

## Cómo funciona

`ccam add work --shortcut ccwork` ejecuta estos pasos en `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` comprueba el nombre y `shim.Check` comprueba el atajo antes de crear nada.
2. `profile.Add` crea `~/.local/share/ccam/profiles/work/` con el modo `0700` y le asigna a la cuenta el número de posición libre más bajo (`Config.NextSlot`), de modo que el número de una cuenta quitada se reutiliza. La primera cuenta pasa a ser la predeterminada.
3. La cuenta se une a sus grupos (`default` salvo que `--group` o `--no-group` indiquen otra cosa; la primera cuenta crea `default`), y `applyShare` enlaza la carpeta nueva y vacía con ellos mediante `share.Build` y después `Apply`, por ejemplo `settings.json` con `groups/default/settings.json`. Si esto falla, se elimina la carpeta nueva, porque una carpeta sobrante bloquearía cada reintento con `folder already exists`.
4. La cuenta se guarda en `~/.config/ccam/config.toml` mediante `fsutil.AtomicWrite`: un archivo temporal, `fsync` y luego `rename`, así que un fallo nunca deja la configuración a medias.
5. `runLogin` inicia `claude auth login` con `CLAUDE_CONFIG_DIR` apuntando a la carpeta de la cuenta, y `--email` con el último correo electrónico de la cuenta para que la página de inicio de sesión aparezca rellenada. Si el inicio de sesión falla, la cuenta se conserva, y ccam imprime el error y te indica que ejecutes `ccam login work` más tarde.
6. `profile.ReadIdentity` extrae `oauthAccount` de `profiles/work/.claude.json`: el correo electrónico, los ids de la cuenta y de la organización, el nombre de la organización y el plan (`claude_max`, `claude_team` …). Ese archivo puede ocupar varios megabytes, así que `TopLevelValue` busca esa única clave en lugar de analizar el archivo completo. ccam los guarda en `[profile.login]` en `config.toml`.
7. `shim.Install` enlaza `~/.local/bin/ccwork` con el binario de ccam, que busca `argv[0]` en `config.toml` al iniciar, y registra el enlace en `~/.local/state/ccam/shortcuts.json`. Si no se puede crear el enlace, la cuenta se agrega igualmente, ccam indica el motivo (`Could not create shortcuts: …`) y el resumen muestra `shortcut: ccwork (not created; see above)`.

Un inicio de sesión puede tardar minutos, y mientras tanto otro ccam puede guardar la configuración. Por eso los comandos que esperan guardan mediante `env.update`, que vuelve a cargar `config.toml`, aplica el único cambio (aquí, el correo electrónico) y guarda. `ccam remove` vuelve a cargar después de su pregunta por el mismo motivo.

`ccam login` ejecuta los pasos 5 y 6 para una cuenta existente. Claude Code toma el correo electrónico solo como sugerencia para la página de inicio de sesión, así que un inicio de sesión puede acabar en otra cuenta de Claude distinta de la que tenía la cuenta de ccam; ccam compara el nuevo inicio de sesión con `[profile.login]` y lo indica (`ccam: this login is other@company.com, but work was dev@company.com before.`). Conserva el nuevo inicio de sesión y te indica cómo deshacerlo: ``To undo it, run `ccam login work` again and sign in as dev@company.com.`` Un inicio de sesión en otra cuenta de Claude también descarta el [registro de límite](/features/usage#the-limit-record) de la cuenta, ya que esos límites eran de la otra cuenta. Cuando otra cuenta de ccam tiene la sesión iniciada en la misma cuenta de Claude, ccam avisa de que las dos comparten sus límites de uso. `ccam login --expired` inicia sesión, una tras otra, en todas las cuentas que nunca iniciaron sesión o cuyo inicio de sesión caducó: su token de actualización caducó, o Claude Code borró el inicio de sesión después de que el servidor lo rechazara. Antes de que se abra el primer navegador, comprueba que `config.toml` se pueda guardar.

`ccam logout` ejecuta `claude auth logout` en la cuenta y borra el correo electrónico guardado; el registro del inicio de sesión se conserva, así que el siguiente inicio de sesión se compara con él. Claude Code revoca el inicio de sesión y elimina lo que guarda para la cuenta, también los inicios de sesión de los servidores MCP. Por eso ccam primero dice qué se va con él y pregunta (`--yes` omite la pregunta): cuando el inicio de sesión tiene inicios de sesión MCP, cuando la cuenta está en ejecución (sus sesiones también se cierran) y cuando la carpeta de otra cuenta tiene el mismo inicio de sesión (una carpeta copiada), que el cierre de sesión también revoca:

```console
$ ccam logout side
Logging out also removes side's MCP server sign-ins (linear); sign in to them again with /mcp.
Log out side? [y/N]
```

Para iniciar sesión con otra cuenta de Claude basta con `ccam login <name>`, y conserva los inicios de sesión MCP. `ccam use` solo reescribe `default` en `config.toml`.

`ccam list` agrega una columna `Organization` cuando las cuentas pertenecen a más de una organización (los planes personales muestran `personal`), y `ccam list --json` incluye `plan` y `org` para cada cuenta.

### Iniciar sesión en una máquina remota

Por SSH, la parte del navegador de `claude auth login` puede fallar. Claude Code espera el inicio de sesión en un servidor de retorno en `http://localhost:<port>/callback` del servidor, pero el navegador que completa el inicio de sesión se ejecuta en tu computadora, donde esa dirección no lleva a ninguna parte, y el navegador muestra un error como `ERR_CONNECTION_REFUSED`. La dirección sigue llevando el código de un solo uso, así que ccam puede entregarlo:

<!-- motion:accounts-relay -->

`runLogin` no le cede tu terminal a Claude Code. Conecta la entrada estándar de Claude Code a una tubería y `relayLogin` copia en ella cada línea que escribes, de modo que un código pegado en la solicitud `Paste code here` de Claude Code sigue llegando. Cuando una línea es una dirección de retorno, `callbackURL` la acepta solo si el host es `localhost`, `127.0.0.1` o `[::1]`, se indica un puerto, la ruta es `/callback` y están presentes tanto `code` como `state`; `http://` se puede omitir. Después, `deliverCallback` solicita esa dirección en el servidor con un límite de 10 segundos y sin seguir redirecciones, como lo habría hecho el navegador. La solicitud llega al servidor de retorno del `claude auth login` que está esperando (la línea discontinua del diagrama). Claude Code recibe el código y guarda el inicio de sesión, en Linux en `.credentials.json` dentro de la carpeta de la cuenta. Una línea pegada nunca puede hacer que ccam se conecte a otro host.

ccam imprime esto antes de que empiece el inicio de sesión:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Si todo sale bien, imprime `→ Passed the sign-in to Claude Code on this machine.`

Mientras el inicio de sesión espera una entrada en una terminal, ccam desactiva el eco y el modo de línea de la terminal y muestra un `*` por cada carácter que escribes, para que un código o una dirección de retorno pegados no aparezcan en pantalla. `readLine` se encarga de editar la línea: `Enter` envía la línea, `Backspace` borra un carácter y `Ctrl+D` en una línea vacía termina la entrada. La terminal se restaura cuando termina el inicio de sesión, y también con `Ctrl+C`, `SIGTERM` o `SIGHUP`, que siguen cerrando ccam. Sin una terminal en la entrada estándar y la salida de error estándar, las líneas se copian tal como llegan.

### Cambiar el nombre, quitar e importar

`ccam rename` y `ccam remove` actúan sobre la carpeta de la cuenta y sobre lo que ccam guarda en otros lugares con el nombre de la cuenta. `nameState` lista esas rutas: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>`, las cachés de uso `~/.local/state/ccam/usage/<name>.json` y `<name>.statusline.json`, el registro de límite `<name>.limit.json` y las notas de sesión de `~/.local/state/ccam/handoff/<name>/`. `ccam remove work` las procesa en este orden:

<!-- motion:accounts-remove -->

`ccam remove` primero se niega si la cuenta está en ejecución, luego pregunta y vuelve a comprobarlo después de la respuesta (detecta una sesión iniciada mientras la pregunta esperaba). Comprueba que `config.toml` se pueda guardar antes de cambiar nada. Si la cuenta comparte extensiones, `RewritePluginPaths` reescribe los archivos `installed_plugins.json` y `known_marketplaces.json` del grupo, porque Claude Code guarda rutas absolutas de plugins que pueden pasar por la carpeta de esta cuenta. A continuación se cierra la sesión de una cuenta gestionada: `claude auth logout` se ejecuta en su carpeta, así que ningún inicio de sesión válido espera en la papelera (ni, en macOS, en un ítem del Llavero con el nombre de la carpeta, que heredaría una cuenta posterior con ese nombre). Esto se omite en una cuenta con clave de API y cuando la carpeta de otra cuenta tiene el mismo inicio de sesión, que el cierre de sesión también revocaría; un cierre de sesión que falla se informa y la eliminación continúa. `profile.Trash` mueve la carpeta a `~/.local/share/ccam/trash/<name>-<time>/`, la cuenta sale de `config.toml`, las rutas de `nameState` se eliminan para que una cuenta posterior con el mismo nombre no las herede, y se quita el enlace del atajo: `✓ side removed (folder moved to /home/you/.local/share/ccam/trash/side-20261023-212306; logged out first)`. Si al final `config.toml` no se puede guardar, la carpeta vuelve de la papelera. Para recuperar una cuenta quitada, mueve la carpeta de vuelta desde la papelera y vuelve a iniciar sesión.

`ccam rename` cambia el nombre de la carpeta, reescribe las rutas de plugins de su `plugins/` de la carpeta anterior a la nueva, conserva la cuenta predeterminada y mueve las rutas de `nameState` al nombre nuevo; las notas de sesión se actualizan para que `ccam next --continue` las siga encontrando. El atajo conserva su nombre. Si falla la reescritura de los plugins, o no se puede guardar `config.toml`, la carpeta recupera su nombre anterior. En macOS la carpeta conserva su nombre anterior, ya que el Llavero guarda el inicio de sesión y los inicios de sesión MCP con la ruta de la carpeta y una carpeta movida quedaría sin sesión; lo mismo ocurre en Windows cuando Claude Code guarda los inicios de sesión en el Administrador de credenciales. La configuración registra entonces la carpeta como `folder` junto al nombre nuevo, y ccam dice `its folder stays …: the macOS Keychain keeps the login under the folder's path`.

`ccam import` registra una carpeta donde está, como `dir` en `config.toml`. Nunca la mueve ni la copia: un token de actualización copiado haría que las copias se cerraran la sesión entre sí. Sin una carpeta como argumento, busca carpetas `~/.claude` y `~/.claude-*` no registradas y pregunta por cada una; los nombres salen de la carpeta (`~/.claude` es `main`, `~/.claude-work` es `work`, `~/.claude-2` es `acct2`), y una carpeta de la que no se puede derivar un nombre, o cuyo nombre ya es el de una cuenta, aparece con el comando `--name` que debes usar. Una carpeta encontrada que falla se informa como `Skipped <path>: <reason>` y las demás se importan. Primero se recogen las respuestas y luego las carpetas se registran en `config.toml` tal como esté en ese momento, así que se conserva lo que otros comandos de ccam guardaron mientras tanto; las líneas `✓` aparecen después de guardar. Sin `CLAUDE_CONFIG_DIR`, Claude Code guarda `.claude.json` en tu carpeta personal junto a `~/.claude`, así que cuando importas `~/.claude` y no tiene su propio `.claude.json`, ccam crea el enlace `~/.claude/.claude.json` a `~/.claude.json`. Así, `claude` a secas y `ccam main` ven la misma cuenta. Si `~/.claude` tiene su propio `.claude.json` junto a `~/.claude.json`, ccam deja ambos y lo indica: `ccam main` usaría el interior y `claude` a secas el exterior, con sus datos de inicio de sesión y sus servidores MCP separados, así que reemplaza el archivo interior por un enlace para quedarte con uno ([doctor](/features/doctor) también lo informa). Cuando la carpeta tiene el mismo inicio de sesión que una cuenta registrada (una carpeta copiada), ccam lo indica: la que lo renueve primero cierra la sesión de la otra, así que ejecuta `ccam login <name>` antes de usarla. Las cuentas importadas no están en ningún grupo hasta que las agregas a uno (`ccam group join <group> <name>`), y no se crea ningún atajo.

`ccam completion bash|zsh` imprime un script de autocompletado que también completa los nombres de cuenta.

## Seguridad y casos límite

- `ccam rename` y `ccam remove` se niegan mientras la cuenta está en ejecución (`refuseRunning`), con ``work is running; quit its Claude Code sessions first (`ccam agents` lists them, `ccam work stop <id>` ends a background one)``, porque la sesión sigue escribiendo en la carpeta. La comprobación abarca Claude Code sin importar cómo se haya iniciado: mediante ccam, `claude` a secas, un IDE, sesiones en segundo plano y su daemon.
- `ccam rename`, `ccam remove` y `ccam logout` se niegan mientras otro ccam renueva o cambia ese inicio de sesión: ``another ccam is using the login of work right now (a refresh, rename, remove or logout); try again in a moment``.
- Mientras un cambio de grupos se ejecuta en otro ccam, `ccam rename` y `ccam remove` se niegan con `another ccam is changing the groups right now; try again when it is done`; después de uno que se quedó a medias, con ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.
- `ccam remove` en una cuenta importada conserva la carpeta y su inicio de sesión, para `claude` a secas. Si la cuenta está en grupos, primero sale de ellos como hace `ccam group leave`: cada elemento compartido se convierte en una copia del del grupo en la carpeta, y los servidores MCP del grupo van a su `.claude.json`. Después se anula su registro: `✓ main unregistered (folder ... left in place; it left its groups with a copy of what they share)`. Las copias de seguridad de su `.claude.json` que hizo ccam pasan a `<folder>/.ccam-backup/mcp-<time>`; si no se pueden mover, se quedan en `~/.local/state/ccam/mcp/<name>/backups` y el mensaje lo indica. Tu `~/.claude` no se mueve a la papelera.
- Dos cuentas no pueden tener la misma carpeta: un `config.toml` que nombra una carpeta dos veces no se puede cargar.
- `ccam add` rechaza un nombre cuya carpeta ya existe en `profiles/` y te remite a `ccam import`. En macOS también se detiene cuando el Llavero todavía tiene un inicio de sesión o inicios de sesión MCP para la ruta de la nueva carpeta, dejados por un ccam anterior, y nombra el comando `security delete-generic-password` que los quita. Si `config.toml` no se puede guardar, la carpeta nueva se vuelve a quitar.
- `ccam import` rechaza una carpeta que ya está registrada, una ruta que no es una carpeta y un nombre que ya está en uso. También detecta una copia del `~/.claude` de `claude` a secas cuando `~/.claude` no está registrado.
- `ccam logout` pregunta primero cuando la cuenta está en ejecución; `ccam login` no lo comprueba.
- Si el reenvío no puede llegar a Claude Code, por ejemplo porque el inicio de sesión ya terminó, imprime `ccam: could not pass the sign-in to Claude Code:` con el motivo. Un error HTTP suele significar que el código venció; vuelve a iniciar el inicio de sesión.

## Relacionado

- [Iniciar Claude Code](/features/launch)
- [Atajos](/features/shortcuts)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Varias cuentas de Claude Code](/guides/multiple-claude-code-accounts)
