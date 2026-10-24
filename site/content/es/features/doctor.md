---
title: Comprobar y reparar
description: Comprueba enlaces compartidos, permisos de archivos, inicios de sesión y la versión de Claude Code con ccam doctor, y repara lo seguro con ccam doctor --fix.
updated: 2026-10-23
order: 9
translation: machine
sourceHash: 11ad34630dcf
---

`ccam doctor` comprueba tu instalación de ccam: el binario de Claude Code, la carpeta, el inicio de sesión y los enlaces compartidos de cada cuenta, los servidores MCP compartidos, el endpoint de uso y tus atajos. `ccam doctor --fix` repara lo que puede sin perder datos y antes mueve a una carpeta de copia de seguridad cualquier archivo que reemplace. Ejecútalo cuando `ccam <name>` se detenga por un problema con un enlace compartido, después de restaurar archivos a mano o después de una actualización de Claude Code.

## Uso

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

La primera línea marca todas las cuentas a la vez, para que veas dónde están los problemas antes de leer nada más. Debajo, cada cuenta tiene un recuadro, y los resultados que no pertenecen a ninguna cuenta (Claude Code, variables de entorno, plugins, servidores MCP, atajos, reglas de carpeta, grupos y ajustes) van en un recuadro `General`. El título del recuadro muestra la peor marca de la cuenta y cuántos problemas tiene. Cada fila del recuadro es un área comprobada, siempre `Folder`, `Login` y `Share links` en ese orden, y luego las demás, con una marca: `✔` correcto, `ℹ` información, `⚠` una advertencia, `✘` un error. Una carpeta sin ningún problema muestra `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Con `--fix`, cada resultado que se reparó termina en `→ fixed` y deja de contar como problema:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

En una terminal, los recuadros se ajustan a su ancho y los mensajes largos se parten en varias líneas, con los colores de la propia terminal (`NO_COLOR` los desactiva). Un recuadro nunca es más estrecho que 29 columnas, así que ninguna palabra se corta en una terminal estrecha. Enviada a una tubería o a un archivo, la salida no tiene colores y ningún mensaje se parte ni se corta, así que `ccam doctor | grep settings.json` encuentra la línea completa; cada recuadro tiene entonces el ancho de la línea más larga. La pestaña Health del panel muestra los mismos recuadros.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--fix` | Repara lo que es seguro reparar (se hace una copia de seguridad de los archivos modificados) | desactivada |
| `--dry-run` (global) | Ejecuta todas las comprobaciones pero no repara nada, ni siquiera con `--fix` | desactivada |

`ccam doctor` no acepta argumentos. Respeta `--dry-run`: `ccam --dry-run doctor --fix` imprime los resultados de `ccam doctor`, agrega lo que `--fix` apartaría y no cambia nada:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

El código de salida es `1` cuando queda al menos un resultado `✘` sin reparar, y `0` en caso contrario. Las advertencias y la información nunca hacen fallar la ejecución. Cuando la ejecución falla sin `--fix`, ccam termina con una línea en blanco y `` `ccam doctor --fix` repairs what it can. ``. Una configuración que no se puede cargar también termina con `1` y un error `ccam:`.

## Cómo funciona

`internal/cli/doctor.go` carga `~/.config/ccam/config.toml` y llama a `doctor.Run` en `internal/doctor/doctor.go`, que ejecuta las comprobaciones en un orden fijo y devuelve los resultados, cada uno con la cuenta a la que pertenece y el área que comprobó. `tui.HealthReport` los agrupa en los recuadros. `--fix` se pasa como `Options.Fix`, que vale false en una ejecución de prueba.

<!-- motion:doctor-checks -->

El bloque central se ejecuta una vez por cuenta. Las comprobaciones, y lo que hace `--fix` en cada una:

1. **Versión de Claude Code** (`checkClaude`). Resuelve el binario `claude` igual que el iniciador y ejecuta `claude --version`. Un binario que falta, un `claude_bin` que no es un archivo ejecutable o un `--version` que falla es un error. Sin reparación.
2. **Credenciales en el entorno** (`checkAuthEnv`). Advierte cuando está definida `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`: todas las cuentas que se inician desde esa shell la heredan, y Claude Code la prefiere al inicio de sesión propio de la cuenta. Sin reparación; quítala tú.
3. **Carpeta** (`checkProfile`). Una carpeta de cuenta que falta es un error y termina las comprobaciones de esa cuenta. Los permisos con algún bit de grupo u otros activado son una advertencia; `--fix` ejecuta `chmod 0700`. Un `~/.claude` importado con su propio `.claude.json` junto a `~/.claude.json` es información: `ccam main` y `claude` a secas mantienen entonces datos de inicio de sesión y servidores MCP separados. Windows no tiene bits de permisos, así que allí se omiten esta comprobación y la de `0600` de más abajo; la carpeta de tu perfil de usuario está protegida por su lista de acceso.
4. **Archivo de inicio de sesión** (`checkCredentials`). `.credentials.json` se comprueba con `lstat`. Un enlace simbólico es un error, porque cada cuenta debe tener su propio token, y no se repara. Un modo distinto de `0600` es una advertencia; `--fix` ejecuta `chmod 0600`. Después se lee el inicio de sesión (en macOS, primero desde el Llavero). La falta de inicio de sesión, o un archivo de inicio de sesión sin ninguno, es una advertencia que remite a `ccam login <name>`. Un token de actualización vencido, o un inicio de sesión que Claude Code borró después de que el servidor lo rechazara, es el error `` login expired → `ccam login <name>` ``; un token de acceso vencido es solo información, ya que Claude Code lo renueva en el siguiente inicio. Un inicio de sesión guardado en el Administrador de credenciales de Windows es información, ya que ccam no puede leerlo. Un inicio de sesión válido muestra el correo electrónico con el que la carpeta tiene la sesión iniciada ahora y el plan (`login valid (dev@company.com, Claude Team · Acme)`). Nada de esto se repara; `ccam login --expired` inicia sesión en todas las cuentas que lo necesitan.
5. **Enlaces compartidos** (`share.Check`). Compara cada ruta de `share.Layout` con los grupos de la cuenta, usando solo `lstat` y `readlink`. Informa de cinco tipos: un enlace que falta, un enlace que apunta a otro sitio, un archivo normal donde debería haber un enlace, un enlace que sobra para un elemento que ahora es local y un enlace correcto cuya carpeta de grupo falta (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; un inicio no comprueba esto). Un enlace que falta y una carpeta de grupo que falta son advertencias; los demás son errores. También se informa de una copia de vuelta desde el grupo que se quedó a medias: `--fix` pone en su lugar una copia preparada completa, pero cuando la ruta se volvió a crear desde entonces (por un inicio), nombra la carpeta con la copia completa para que la resuelvas tú. Las rutas que una versión posterior de ccam agregó a un elemento (`file-history/`, `tasks/` y `paste-cache/` para el historial, `keybindings.json` para la configuración) son solo una advertencia mientras la cuenta todavía tenga su propia copia: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Un inicio nunca se detiene por ellas.
6. **Entradas desconocidas** (`checkUnknownEntries`). Lista los nombres de primer nivel de la carpeta de la cuenta que ccam no conoce (`share.Known`), como indicio de que la estructura de Claude Code cambió. Solo información.
7. **Rutas de plugins** (`checkPlugins`). Los registros de plugins (`installed_plugins.json`, `known_marketplaces.json`) guardan rutas absolutas. En el `plugins/` de cada grupo, las rutas que pasan por la carpeta de una cuenta se rompen cuando esa cuenta se quita o cambia de nombre. `--fix` las reescribe hacia la carpeta del grupo con `profile.RewritePluginPaths`, una escritura atómica del archivo de registro.
8. **Servidores MCP** (`checkMCP`). Carga una vez el `mcp.json` de cada grupo. Un archivo que no se puede leer es un error; un servidor que no pasa `mcp.Validate` es una advertencia. Los servidores del `~/.claude.json` de `claude` a secas que ninguna cuenta carga son una advertencia que nombra el comando `ccam migrate --mcp` que los copia (con `--group <a|b>` cuando varios grupos comparten MCP). Los servidores que solo tiene una cuenta de un grupo son información con el comando `ccam mcp import` que los comparte, y una cuenta que no ha iniciado sesión en servidores MCP en los que otras cuentas sí la iniciaron recibe una advertencia para iniciar sesión con `/mcp` (ccam nunca copia inicios de sesión). Sin reparación.
9. **Formato de la respuesta de uso** (`checkUsageFormat`). Consulta el endpoint de uso, sin la caché, para una cuenta tras otra hasta que una responde, todo dentro de un único límite de 4 segundos. Con `usage_source = "claude"`, en su lugar le pregunta a Claude Code por una cuenta. Si la respuesta ya no se puede analizar, advierte que `ccam usage` recurre a la caché de la línea de estado; cuando ninguna cuenta respondió, lo indica como información (`could not check the usage response format: …`). Es la única comprobación que usa la red.
10. **Atajos** (`checkShims`). Solo cuando una cuenta tiene un atajo. Los enlaces que faltan en `~/.local/bin` son una advertencia, y también lo son los atajos que creó ccam y que todavía ejecutan un ccam anterior, uno que se reemplazó o se movió (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` llama a `shim.Install`, que crea y renueva solo los enlaces que creó ccam y registró en `~/.local/state/ccam/shortcuts.json`, y nunca toca otros archivos. La advertencia termina en `→ fixed` solo cuando se creó cada uno de esos enlaces; si alguno no se pudo crear, se agrega el motivo. Un atajo que ejecuta otro ccam que sigue instalado funciona, y es información; `--fix` lo apunta a este ccam. Una carpeta bin que no está en `PATH` es una advertencia y no se repara.
11. **Carpetas** (`checkFolders`). Solo cuando existen entradas de [`ccam folder`](/features/folder). Una entrada cuya cuenta no existe es un error, y una entrada cuya carpeta ya no está es una advertencia. Ninguna se repara, porque quitar una entrada cambiaría qué cuenta se inicia allí.
12. **Grupos y ajustes** (`checkGroups`). Un miembro de grupo que no es una cuenta, una cuenta cuyo nombre se ha convertido desde entonces en un comando de ccam (iníciala con `ccam run <name>` o cámbiale el nombre), un `theme` o `usage_source` desconocido y las líneas de `config.toml` que ccam no lee son advertencias. Estas últimas suelen ser un ajuste escrito debajo de una tabla `[[group]]`, `[[profile]]` o `[[folder]]`, donde pertenece a esa tabla (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); los comandos que cambian la configuración se niegan hasta que se corrija. Sin reparación.
13. **Inicios de sesión** (`checkIdentities`). Lee con qué cuenta tiene iniciada la sesión cada carpeta a partir de su `.claude.json`. Una carpeta que ahora tiene la sesión iniciada en una cuenta distinta de la que registró `ccam login` (por ejemplo, tras `/login` dentro de una sesión) es una advertencia; una cuenta sin registro de inicio de sesión, como una importada o creada por un ccam anterior, se compara por su correo electrónico registrado. Una carpeta con la sesión iniciada en la misma cuenta de Claude que otra también es una advertencia: las dos comparten sus límites de uso. Una carpeta que tiene el mismo token de inicio de sesión que otra, una carpeta copiada (incluido el `~/.claude` de `claude` a secas), es un error, ya que la que lo renueve primero cierra la sesión de la otra; ejecuta `ccam login <name>` para darle el suyo. Sin reparación.

Antes de las cuentas, un **cambio de grupos** que se quedó a medias es un error en el recuadro `General`: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` lo termina en el orden en que se planificó, con las respuestas dadas a sus conflictos; hasta entonces no empieza ningún otro cambio de grupos, cambio de nombre ni eliminación. A las cuentas incluidas solo se les crean los enlaces que faltan hasta que termine. Si no se puede terminar, el error dice que quitar ese archivo renuncia al resto del cambio. `ccam --dry-run doctor --fix` muestra el orden (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) y lo que apartaría la parte de cada cuenta. Un cambio que otro ccam está realizando en este momento es información, y mientras se ejecuta `--fix` no repara enlaces compartidos (`its share links are not repaired while another ccam changes the groups; …`). `--fix` mantiene `~/.local/state/ccam/group-change.lock` mientras repara, y vuelve a leer `config.toml` después de terminar un cambio.

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

Aquí `settings.json` es un archivo normal. Como el grupo ya tiene un `settings.json` distinto, el conflicto toma su opción predeterminada, `KeepGroup`, así que la copia del grupo se conserva. El archivo de la cuenta se mueve a `.ccam-backup/<time>/settings.json` dentro de la carpeta de la cuenta, y `os.Symlink` enlaza la ruta con el grupo. `--fix` lo indica: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Otros casos:

- Si el grupo todavía no tiene una copia, el archivo de la cuenta se mueve al grupo y se convierte en su primer contenido.
- En una carpeta, los archivos que el grupo no tiene se mueven al grupo. Los archivos que el grupo tiene idénticos byte a byte se descartan, y los que difieren van a la copia de seguridad, nombrados en una nota (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Una carpeta de grupo que falta se crea primero.
- En `history.jsonl`, la opción predeterminada es `merge by line`: las líneas del grupo y luego las tuyas nuevas. La copia anterior del grupo se guarda en `.ccam-backup/<time>/group/`.
- Un enlace a otro grupo se reemplaza. Un enlace a una carpeta fuera de ccam va a la copia de seguridad, y su contenido sirve de base para el grupo si el grupo está vacío.
- Un enlace que sobra para un elemento que ahora es local se quita, y la cuenta recibe una copia del contenido del grupo.

Aparte de quitar enlaces y descartar archivos que el grupo ya tiene idénticos byte a byte, cada paso mueve datos al grupo o a la carpeta de copia de seguridad. No se pierde nada.

## Seguridad y casos límite

- `ccam doctor` a secas nunca cambia nada. En scripts, usa `--dry-run` con `--fix` para estar seguro.
- Para una cuenta en ejecución, sin importar cómo se haya iniciado Claude Code, `doctor --fix` solo crea los enlaces y las carpetas de grupo que faltan, con lo que no se pierde nada. Cuando queda algo más, agrega `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, porque la sesión sigue escribiendo en los archivos que una reparación movería.
- Si falla la reparación de los enlaces compartidos, ccam informa `could not repair share links: <error>` para esa cuenta y omite su comprobación de entradas desconocidas. Las demás cuentas se siguen comprobando.
- Si un cambio con `ccam group` falla a medias, los grupos nuevos ya están guardados y el cambio queda anotado en `~/.local/state/ccam/group-change.json`; ccam te indica que termines con `ccam doctor --fix`, que realiza el resto, incluidos los pasos de MCP.
- Un `.credentials.json` enlazado simbólicamente, un inicio de sesión vencido y las credenciales del entorno se informan, pero nunca se cambian.
- Las reescrituras del registro de plugins se escriben de forma atómica, pero no se copian a `.ccam-backup/`.

## Relacionado

- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Copiar una configuración existente](/features/migrate)
- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Comandos](/docs/commands)
