---
title: Comprobar y reparar
description: Comprueba enlaces compartidos, permisos de archivos, inicios de sesión y la versión de Claude Code con ccam doctor, y repara lo seguro con ccam doctor --fix.
updated: 2026-10-21
order: 9
translation: machine
sourceHash: 9be63350cb69
---

`ccam doctor` comprueba tu instalación de ccam: el binario de Claude Code, la carpeta, el inicio de sesión y los enlaces compartidos de cada cuenta, los servidores MCP compartidos, el endpoint de uso y tus atajos. `ccam doctor --fix` repara lo que puede sin perder datos y antes mueve a una carpeta de copia de seguridad cualquier archivo que reemplace. Ejecútalo cuando `ccam <name>` se detenga por un problema con un enlace compartido, después de restaurar archivos a mano o después de una actualización de Claude Code.

## Uso

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

La primera línea marca todas las cuentas a la vez, para que veas dónde están los problemas antes de leer nada más. Debajo, cada cuenta tiene un recuadro, y los resultados que no pertenecen a ninguna cuenta (Claude Code, variables de entorno, plugins, servidores MCP, atajos, reglas de carpeta, grupos y ajustes) van en un recuadro `General`. El título del recuadro muestra la peor marca de la cuenta y cuántos problemas tiene. Cada fila del recuadro es un área comprobada, siempre `Folder`, `Login` y `Share links` en ese orden, y luego las demás, con una marca: `✔` correcto, `ℹ` información, `⚠` una advertencia, `✘` un error. Una carpeta sin ningún problema muestra `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp               │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Con `--fix`, cada resultado que se reparó termina en `→ fixed` y deja de contar como problema:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ────────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed       │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ───────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed         │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →    │
 │                fixed                                                       │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link → fixed        │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯
```

En una terminal, los recuadros se ajustan a su ancho y los mensajes largos se parten en varias líneas, con los colores de la propia terminal (`NO_COLOR` los desactiva). Enviada a una tubería o a un archivo, la salida no tiene colores y ningún mensaje se parte ni se corta, así que `ccam doctor | grep settings.json` encuentra la línea completa; cada recuadro tiene entonces el ancho de la línea más larga. La pestaña Health del panel muestra los mismos recuadros.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--fix` | Repara lo que es seguro reparar (se hace una copia de seguridad de los archivos modificados) | desactivada |
| `--dry-run` (global) | Ejecuta todas las comprobaciones pero no repara nada, ni siquiera con `--fix` | desactivada |

`ccam doctor` no acepta argumentos. Respeta `--dry-run`: `ccam --dry-run doctor --fix` imprime los mismos resultados que `ccam doctor` y no cambia nada.

El código de salida es `1` cuando queda al menos un resultado `✘` sin reparar, y `0` en caso contrario. Las advertencias y la información nunca hacen fallar la ejecución. Cuando la ejecución falla sin `--fix`, ccam termina con una línea en blanco y `` `ccam doctor --fix` repairs what it can. ``. Una configuración que no se puede cargar también termina con `1` y un error `ccam:`.

## Cómo funciona

`internal/cli/doctor.go` carga `~/.config/ccam/config.toml` y llama a `doctor.Run` en `internal/doctor/doctor.go`, que ejecuta las comprobaciones en un orden fijo y devuelve los resultados, cada uno con la cuenta a la que pertenece y el área que comprobó. `tui.HealthReport` los agrupa en los recuadros. `--fix` se pasa como `Options.Fix`, que vale false en una ejecución de prueba.

<!-- motion:doctor-checks -->

El bloque central se ejecuta una vez por cuenta. Las comprobaciones, y lo que hace `--fix` en cada una:

1. **Versión de Claude Code** (`checkClaude`). Resuelve el binario `claude` igual que el iniciador y ejecuta `claude --version`. Si la versión no empieza con `2.1.302`, la versión con la que se verificó la estructura de carpetas de ccam (`doctor.LayoutVersion`), agrega la información `ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Un binario que falta o un `--version` que falla es un error. Sin reparación.
2. **Credenciales en el entorno** (`checkAuthEnv`). Advierte cuando está definida `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`: todas las cuentas que se inician desde esa shell la heredan, y Claude Code la prefiere al inicio de sesión propio de la cuenta. Sin reparación; quítala tú.
3. **Carpeta** (`checkProfile`). Una carpeta de cuenta que falta es un error y termina las comprobaciones de esa cuenta. Los permisos con algún bit de grupo u otros activado son una advertencia; `--fix` ejecuta `chmod 0700`. Windows no tiene bits de permisos, así que allí se omiten esta comprobación y la de `0600` de más abajo; la carpeta de tu perfil de usuario está protegida por su lista de acceso.
4. **Archivo de inicio de sesión** (`checkCredentials`). `.credentials.json` se comprueba con `lstat`. Un enlace simbólico es un error, porque cada cuenta debe tener su propio token, y no se repara. Un modo distinto de `0600` es una advertencia; `--fix` ejecuta `chmod 0600`. Después se lee el inicio de sesión (en macOS, primero desde el Llavero). La falta de inicio de sesión es una advertencia que remite a `ccam login <name>`. Un token de actualización vencido es el error `login expired`; un token de acceso vencido es solo información, ya que Claude Code lo renueva en el siguiente inicio. Un inicio de sesión válido muestra el correo electrónico y el plan (`login valid (dev@company.com, Claude Team · Acme)`). Nada de esto se repara; `ccam login --expired` inicia sesión en todas las cuentas que lo necesitan.
5. **Enlaces compartidos** (`share.Check`). Compara cada ruta de `share.Layout` con los grupos de la cuenta, usando solo `lstat` y `readlink`. Informa de cuatro tipos: un enlace que falta, un enlace que apunta a otro sitio, un archivo normal donde debería haber un enlace y un enlace que sobra para un elemento que ahora es local. Un enlace que falta es una advertencia; los demás son errores. Las rutas que una versión posterior de ccam agregó a un elemento (`file-history/`, `tasks/` y `paste-cache/` para el historial, `keybindings.json` para la configuración) son solo una advertencia mientras la cuenta todavía tenga su propia copia: `is not shared with the group yet (ccam doctor --fix moves it in)`. Un inicio nunca se detiene por ellas.
6. **Entradas desconocidas** (`checkUnknownEntries`). Lista los nombres de primer nivel de la carpeta de la cuenta que ccam no conoce (`share.Known`), como indicio de que la estructura de Claude Code cambió. Solo información.
7. **Rutas de plugins** (`checkPlugins`). Los registros de plugins (`installed_plugins.json`, `known_marketplaces.json`) guardan rutas absolutas. En el `plugins/` de cada grupo, las rutas que pasan por la carpeta de una cuenta se rompen cuando esa cuenta se quita o cambia de nombre. `--fix` las reescribe hacia la carpeta del grupo con `profile.RewritePluginPaths`, una escritura atómica del archivo de registro.
8. **Servidores MCP** (`checkMCP`). Carga una vez el `mcp.json` de cada grupo. Un archivo que no se puede leer es un error; un servidor que no pasa `mcp.Validate` es una advertencia. Sin reparación.
9. **Formato de la respuesta de uso** (`checkUsageFormat`). Consulta el endpoint de uso, sin la caché y con el límite habitual de 3 segundos, para una cuenta tras otra hasta que una responde. Si la respuesta ya no se puede analizar, advierte que `ccam usage` recurre a la caché de la línea de estado. Es la única comprobación que usa la red.
10. **Atajos** (`checkShims`). Solo cuando una cuenta tiene un atajo. Los enlaces que faltan en `~/.local/bin` son una advertencia; `--fix` llama a `shim.Install`, que solo crea los enlaces que faltan y nunca toca otros archivos. Una carpeta bin que no está en `PATH` es una advertencia y no se repara.
11. **Carpetas** (`checkFolders`). Solo cuando existen entradas de [`ccam folder`](/features/folder). Una entrada cuya cuenta no existe es un error, y una entrada cuya carpeta ya no está es una advertencia. Ninguna se repara, porque quitar una entrada cambiaría qué cuenta se inicia allí.
12. **Grupos y ajustes** (`checkGroups`). Un miembro de grupo que no es una cuenta, una cuenta cuyo nombre se ha convertido desde entonces en un comando de ccam (iníciala con `ccam run <name>` o cámbiale el nombre) y un `theme` o `usage_source` desconocido en `config.toml` son advertencias. Sin reparación.
13. **Inicios de sesión** (`checkIdentities`). Lee con qué cuenta tiene iniciada la sesión cada carpeta a partir de su `.claude.json`. Una carpeta que ahora tiene la sesión iniciada en una cuenta distinta de la que registró `ccam login` (por ejemplo, tras `/login` dentro de una sesión) es una advertencia, y también lo es una carpeta con la sesión iniciada en la misma cuenta de Claude que otra: las dos comparten sus límites de uso. Sin reparación.

### Reparar enlaces compartidos

La ruta de inicio crea en el momento un enlace que falta, pero se detiene ante cualquier cosa que podría perder datos y te remite aquí:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Para cualquier problema con un enlace compartido, `--fix` llama a `share.Build` con la configuración de uso compartido de la cuenta como estado actual y como estado de destino, y después a `Plan.Apply`. `Build` vuelve a derivar cada paso a partir de lo que hay en el disco, el mismo plan que usa un cambio con `ccam group`, así que la reparación sigue las mismas reglas que unirse a un grupo.

<!-- motion:doctor-repair -->

Aquí `settings.json` es un archivo normal. Como el grupo ya tiene un `settings.json` distinto, el conflicto toma su opción predeterminada, `KeepGroup`, así que la copia del grupo se conserva. El archivo de la cuenta se mueve a `.ccam-backup/<time>/settings.json` dentro de la carpeta de la cuenta, y `os.Symlink` enlaza la ruta con el grupo. Otros casos:

- Si el grupo todavía no tiene una copia, el archivo de la cuenta se mueve al grupo y se convierte en su primer contenido.
- En una carpeta, los archivos que el grupo no tiene se mueven al grupo; el resto va a la copia de seguridad.
- En `history.jsonl`, la opción predeterminada es `merge by line`: las líneas del grupo y luego las tuyas nuevas. La copia anterior del grupo se guarda en `.ccam-backup/<time>/group/`.
- Un enlace a otro grupo se reemplaza. Un enlace a una carpeta fuera de ccam va a la copia de seguridad, y su contenido sirve de base para el grupo si el grupo está vacío.
- Un enlace que sobra para un elemento que ahora es local se quita, y la cuenta recibe una copia del contenido del grupo.

Aparte de quitar enlaces, cada paso mueve datos al grupo o a la carpeta de copia de seguridad. No se elimina nada.

## Seguridad y casos límite

- `ccam doctor` a secas nunca cambia nada. En scripts, usa `--dry-run` con `--fix` para estar seguro.
- Para una cuenta en ejecución, `doctor --fix` solo crea los enlaces que faltan y agrega `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, porque la sesión sigue escribiendo en los archivos que una reparación movería.
- Si falla la reparación de los enlaces compartidos, ccam informa `could not repair share links: <error>` para esa cuenta y omite su comprobación de entradas desconocidas. Las demás cuentas se siguen comprobando.
- Si un cambio con `ccam group` falla a medias, los grupos nuevos ya están guardados, y ccam te indica que termines con `ccam doctor --fix`.
- Un `.credentials.json` enlazado simbólicamente, un inicio de sesión vencido y las credenciales del entorno se informan, pero nunca se cambian.
- Las reescrituras del registro de plugins se escriben de forma atómica, pero no se copian a `.ccam-backup/`.

## Relacionado

- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Copiar una configuración existente](/features/migrate)
- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Comandos](/docs/commands)
