---
title: Administrar cuentas
description: Agrega, inicia sesión, lista, renombra, quita e importa cuentas de Claude Code con ccam add, ccam login, ccam list, ccam rename, ccam remove y ccam import.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

Cada cuenta de ccam es una carpeta de configuración de Claude Code independiente con su propio inicio de sesión. Estos comandos crean cuentas, inician sesión en ellas, las listan, les cambian el nombre y las quitan, y registran carpetas que ya tienes.

## Uso

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Siempre que un comando acepta `<name|number>`, puedes indicar el nombre de la cuenta o su número de posición, la columna `#` de `ccam list`.

`ccam add` sin opciones en una terminal abre un asistente. Pide el nombre, luego muestra una lista de selección de los elementos compartidos (todo está marcado excepto el historial de conversaciones), después pregunta `Create a shortcut command for work? [y/N]` y, si respondes `y`, `Shortcut name [ccwork]:`. Cualquiera de `--share`, `--shortcut` o `--no-login`, o una tubería en lugar de una terminal, omite el asistente; en ese caso el nombre es obligatorio y no se crea ningún atajo salvo que `--shortcut` indique uno.

`ccam list` muestra la posición, el nombre, el atajo, el correo electrónico, los elementos compartidos y el último uso, con `*` delante de la cuenta predeterminada:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

Con `--json`, cada cuenta es un objeto con `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (se omite si la cuenta nunca se inició) y `share`. `share` asigna a cada elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) su nombre de grupo o `"local"`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `ccam add --no-login` | Crea la cuenta sin iniciar sesión; inicia sesión más tarde con `ccam login <name>` | desactivada |
| `ccam add --share <item=group>` | Comparte un elemento con un grupo, o `item=local` para mantenerlo privado. Se puede repetir | `mcp`, `settings`, `memory`, `extensions` en `default`; `history` local |
| `ccam add --shortcut <command>` | Crea también este atajo en `~/.local/bin` | ninguno |
| `ccam list --json` | Imprime las cuentas en JSON | desactivada |
| `ccam remove -y`, `--yes` | Quita sin preguntar | desactivada |
| `ccam import --name <name>` | Nombre de cuenta para la carpeta que indicas | se deriva del nombre de la carpeta |
| `ccam import -y`, `--yes` | Importa todas las carpetas encontradas sin preguntar | desactivada |
| `--dry-run` | La respetan `ccam list` y `ccam completion`, que no cambian nada. La rechazan `add`, `login`, `logout`, `use`, `rename`, `remove` e `import`, con `` --dry-run is not supported by `ccam add`; nothing was changed `` (con el nombre del comando) | desactivada |

`ccam login`, `ccam logout`, `ccam use` y `ccam rename` no tienen opciones propias. `ccam completion` acepta `bash` o `zsh`; cualquier otro valor falla con `only bash and zsh are supported`.

Los nombres de cuenta tienen de 1 a 32 caracteres entre `a-z`, `0-9`, `-` y `_`, y empiezan con una letra o un dígito. Un nombre no puede tener solo dígitos, porque `ccam 2` significa la posición 2. Tampoco puede ser un comando de ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), así que `ccam <name>` nunca es ambiguo. La lista es `config.Reserved` en `internal/config`.

## Cómo funciona

`ccam add work --shortcut ccwork` ejecuta estos pasos en `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` comprueba el nombre y `shim.Check` comprueba el atajo antes de crear nada.
2. `profile.Add` crea `~/.local/share/ccam/profiles/work/` con el modo `0700` y le asigna a la cuenta el número de posición libre más bajo (`Config.NextSlot`), de modo que el número de una cuenta quitada se reutiliza. La primera cuenta pasa a ser la predeterminada.
3. `applyShare` enlaza la carpeta nueva y vacía con sus grupos de uso compartido mediante `share.Build` y después `Apply`, por ejemplo `settings.json` con `groups/default/settings.json`. Si esto falla, se elimina la carpeta nueva, porque una carpeta sobrante bloquearía cada reintento con `folder already exists`.
4. La cuenta se guarda en `~/.config/ccam/config.toml` mediante `fsutil.AtomicWrite`: un archivo temporal, `fsync` y luego `rename`, así que un fallo nunca deja la configuración a medias.
5. `runLogin` inicia `claude auth login` con `CLAUDE_CONFIG_DIR` apuntando a la carpeta de la cuenta. Si el inicio de sesión falla, la cuenta se conserva, y ccam imprime el error y te indica que ejecutes `ccam login work` más tarde.
6. `profile.ReadEmail` extrae `oauthAccount.emailAddress` de `profiles/work/.claude.json`. Ese archivo puede ocupar varios megabytes, así que `TopLevelValue` busca esa única clave en lugar de analizar el archivo completo.
7. `shim.Install` enlaza `~/.local/bin/ccwork` con el binario de ccam, que busca `argv[0]` en `config.toml` al iniciar.

Un inicio de sesión puede tardar minutos, y mientras tanto otro ccam puede guardar la configuración. Por eso los comandos que esperan guardan mediante `env.update`, que vuelve a cargar `config.toml`, aplica el único cambio (aquí, el correo electrónico) y guarda. `ccam remove` vuelve a cargar después de su pregunta por el mismo motivo.

`ccam login` ejecuta los pasos 5 y 6 para una cuenta existente. `ccam logout` ejecuta `claude auth logout` en la cuenta y borra el correo electrónico guardado. `ccam use` solo reescribe `default` en `config.toml`.

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

`ccam rename` y `ccam remove` actúan sobre la carpeta de la cuenta y sobre lo que ccam guarda en otros lugares con el nombre de la cuenta. `nameState` lista esas rutas: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` y las cachés de uso `~/.local/state/ccam/usage/<name>.json` y `<name>.statusline.json`. `ccam remove work` las procesa en este orden:

<!-- motion:accounts-remove -->

`ccam remove` primero se niega si la cuenta está en ejecución y luego pregunta. Si la cuenta comparte extensiones, `RewritePluginPaths` reescribe los archivos `installed_plugins.json` y `known_marketplaces.json` del grupo, porque Claude Code guarda rutas absolutas de plugins que pueden pasar por la carpeta de esta cuenta. `profile.Trash` mueve la carpeta a `~/.local/share/ccam/trash/<name>-<time>/`, la cuenta sale de `config.toml`, las rutas de `nameState` se eliminan para que una cuenta posterior con el mismo nombre no las herede, y se quita el enlace del atajo. Para recuperarla, mueve la carpeta de vuelta desde la papelera.

`ccam rename` cambia el nombre de la carpeta, reescribe las rutas de plugins de su `plugins/` de la carpeta anterior a la nueva, conserva la cuenta predeterminada y mueve las rutas de `nameState` al nombre nuevo. El atajo conserva su nombre. Si falla la reescritura de los plugins, la carpeta recupera su nombre anterior, ya que la configuración no se guarda.

`ccam import` registra una carpeta donde está, como `dir` en `config.toml`. Nunca la mueve ni la copia: un token de actualización copiado haría que las copias se cerraran la sesión entre sí. Sin una carpeta como argumento, busca carpetas `~/.claude` y `~/.claude-*` no registradas y pregunta por cada una; los nombres salen de la carpeta (`~/.claude` es `main`, `~/.claude-work` es `work`, `~/.claude-2` es `acct2`), y una carpeta de la que no se puede derivar un nombre aparece con el comando `--name` que debes usar. Sin `CLAUDE_CONFIG_DIR`, Claude Code guarda `.claude.json` en tu carpeta personal junto a `~/.claude`, así que cuando importas `~/.claude` y no tiene su propio `.claude.json`, ccam crea el enlace `~/.claude/.claude.json` a `~/.claude.json`. Así, `claude` a secas y `ccam main` ven la misma cuenta. Las cuentas importadas no comparten nada hasta que ejecutas `ccam share`, y no se crea ningún atajo.

`ccam completion bash|zsh` imprime un script de autocompletado que también completa los nombres de cuenta.

## Seguridad y casos límite

- `ccam rename` y `ccam remove` se niegan mientras la cuenta está en ejecución (`refuseRunning`), con `work is running; quit its Claude Code session first`, porque la sesión sigue escribiendo en la carpeta.
- `ccam remove` en una cuenta importada solo anula su registro: `✓ main unregistered (folder ... left in place)`. Tu `~/.claude` no se mueve a la papelera.
- `ccam add` rechaza un nombre cuya carpeta ya existe en `profiles/` y te remite a `ccam import`.
- `ccam import` rechaza una carpeta que ya está registrada, una ruta que no es una carpeta y un nombre que ya está en uso.
- `ccam logout` y `ccam login` no comprueban si la cuenta está en ejecución.
- Si el reenvío no puede llegar a Claude Code, por ejemplo porque el inicio de sesión ya terminó, imprime `ccam: could not pass the sign-in to Claude Code:` con el motivo. Un error HTTP suele significar que el código venció; vuelve a iniciar el inicio de sesión.

## Relacionado

- [Iniciar Claude Code](/features/launch)
- [Atajos](/features/shortcuts)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Varias cuentas de Claude Code](/guides/multiple-claude-code-accounts)
