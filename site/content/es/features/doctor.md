---
title: Comprobar y reparar
description: Comprueba los enlaces compartidos, los permisos de archivos, los inicios de sesión y la versión de Claude Code con ccam doctor, y repara lo que es seguro reparar con ccam doctor --fix.
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` comprueba tu instalación de ccam: el binario de Claude Code, la carpeta, el inicio de sesión y los enlaces compartidos de cada cuenta, los servidores MCP compartidos, el endpoint de uso y tus atajos. `ccam doctor --fix` repara lo que puede sin perder datos y antes mueve a una carpeta de copia de seguridad cualquier archivo que reemplace. Ejecútalo cuando `ccam <name>` se detenga por un problema con un enlace compartido, después de restaurar archivos a mano o después de una actualización de Claude Code.

## Uso

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Cada resultado es una línea con una marca: `✓` correcto, `·` información, `!` una advertencia, `✗` un error. Los resultados de una cuenta empiezan con el nombre de la cuenta.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Con `--fix`, cada resultado que se reparó termina en `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--fix` | Repara lo que es seguro reparar (se hace una copia de seguridad de los archivos modificados) | desactivada |
| `--dry-run` (global) | Ejecuta todas las comprobaciones pero no repara nada, ni siquiera con `--fix` | desactivada |

`ccam doctor` no acepta argumentos. Respeta `--dry-run`: `ccam --dry-run doctor --fix` imprime los mismos resultados que `ccam doctor` y no cambia nada.

El código de salida es `1` cuando queda al menos un resultado `✗` sin reparar, y `0` en caso contrario. Las advertencias y la información nunca hacen fallar la ejecución. Cuando la ejecución falla sin `--fix`, ccam termina con una línea en blanco y `` `ccam doctor --fix` repairs what it can. ``. Una configuración que no se puede cargar también termina con `1` y un error `ccam:`.

## Cómo funciona

`internal/cli/doctor.go` carga `~/.config/ccam/config.toml` y llama a `doctor.Run` en `internal/doctor/doctor.go`, que ejecuta las comprobaciones en un orden fijo y devuelve los resultados. `--fix` se pasa como `Options.Fix`, que vale false en una ejecución de prueba.

<!-- motion:doctor-checks -->

El bloque central se ejecuta una vez por cuenta. Las comprobaciones, y lo que hace `--fix` en cada una:

1. **Versión de Claude Code** (`checkClaude`). Resuelve el binario `claude` igual que el iniciador y ejecuta `claude --version`. Si la versión no empieza con `2.1.302`, la versión con la que se verificó la estructura de carpetas de ccam (`doctor.LayoutVersion`), agrega `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Un binario que falta o un `--version` que falla es un error. Sin reparación.
2. **Credenciales en el entorno** (`checkAuthEnv`). Advierte cuando está definida `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`: todas las cuentas que se inician desde esa shell la heredan, y Claude Code la prefiere al inicio de sesión propio de la cuenta. Sin reparación; quítala tú.
3. **Carpeta** (`checkProfile`). Una carpeta de cuenta que falta es un error y termina las comprobaciones de esa cuenta. Los permisos con algún bit de grupo u otros activado son una advertencia; `--fix` ejecuta `chmod 0700`.
4. **Archivo de inicio de sesión** (`checkCredentials`). `.credentials.json` se comprueba con `lstat`. Un enlace simbólico es un error, porque cada cuenta debe tener su propio token, y no se repara. Un modo distinto de `0600` es una advertencia; `--fix` ejecuta `chmod 0600`. Después se lee el inicio de sesión (en macOS, primero desde el Llavero). La falta de inicio de sesión es una advertencia que remite a `ccam login <name>`. Un token de actualización vencido es el error `login expired`; un token de acceso vencido es solo información, ya que Claude Code lo renueva en el siguiente inicio. Nada de esto se repara.
5. **Enlaces compartidos** (`share.Check`). Compara cada ruta de `share.Layout` con la configuración de uso compartido de la cuenta, usando solo `lstat` y `readlink`. Informa de cuatro tipos: un enlace que falta, un enlace que apunta a otro sitio, un archivo normal donde debería haber un enlace y un enlace que sobra para un elemento que ahora es local. Un enlace que falta es una advertencia; los demás son errores.
6. **Entradas desconocidas** (`checkUnknownEntries`). Lista los nombres de primer nivel de la carpeta de la cuenta que ccam no conoce (`share.Known`), como indicio de que la estructura de Claude Code cambió. Solo información.
7. **Rutas de plugins** (`checkPlugins`). Los registros de plugins (`installed_plugins.json`, `known_marketplaces.json`) guardan rutas absolutas. En el `plugins/` de cada grupo, las rutas que pasan por la carpeta de una cuenta se rompen cuando esa cuenta se quita o cambia de nombre. `--fix` las reescribe hacia la carpeta del grupo con `profile.RewritePluginPaths`, una escritura atómica del archivo de registro.
8. **Servidores MCP** (`checkMCP`). Carga una vez el `mcp.json` de cada grupo. Un archivo que no se puede leer es un error; un servidor que no pasa `mcp.Validate` es una advertencia. Sin reparación.
9. **Formato de la respuesta de uso** (`checkUsageFormat`). Consulta el endpoint de uso, sin la caché y con el límite habitual de 3 segundos, para una cuenta tras otra hasta que una responde. Si la respuesta ya no se puede analizar, advierte que `ccam usage` recurre a la caché de la línea de estado. Es la única comprobación que usa la red.
10. **Atajos** (`checkShims`). Solo cuando una cuenta tiene un atajo. Los enlaces que faltan en `~/.local/bin` son una advertencia; `--fix` llama a `shim.Install`, que solo crea los enlaces que faltan y nunca toca otros archivos. Una carpeta bin que no está en `PATH` es una advertencia y no se repara.

### Reparar enlaces compartidos

La ruta de inicio crea en el momento un enlace que falta, pero se detiene ante cualquier cosa que podría perder datos y te remite aquí:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Para cualquier problema con un enlace compartido, `--fix` llama a `share.Build` con la configuración de uso compartido de la cuenta como estado actual y como estado de destino, y después a `Plan.Apply`. `Build` vuelve a derivar cada paso a partir de lo que hay en el disco, el mismo plan que usa `ccam share`, así que la reparación sigue las mismas reglas que activar el uso compartido.

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
- `doctor --fix` no comprueba si una cuenta está en ejecución. Cierra Claude Code de una cuenta antes de reparar sus enlaces compartidos, porque la sesión sigue escribiendo en esos archivos.
- Si falla la reparación de los enlaces compartidos, ccam informa `could not repair share links: <error>` para esa cuenta y omite su comprobación de entradas desconocidas. Las demás cuentas se siguen comprobando.
- Si `ccam share` falla a medias, primero guarda la configuración de destino y te indica que termines con `ccam doctor --fix`.
- Un `.credentials.json` enlazado simbólicamente, un inicio de sesión vencido y las credenciales del entorno se informan, pero nunca se cambian.
- Las reescrituras del registro de plugins se escriben de forma atómica, pero no se copian a `.ccam-backup/`.

## Relacionado

- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Copiar una configuración existente](/features/migrate)
- [Límites de uso y cambio de cuenta](/docs/usage-limits)
- [Comandos](/docs/commands)
