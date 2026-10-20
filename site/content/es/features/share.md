---
title: Compartir entre cuentas
description: Cómo ccam share enlaza configuración, CLAUDE.md, extensiones e historial a una carpeta de grupo con plan y copias de seguridad, y cómo off los vuelve a copiar.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Cada cuenta tiene su propia carpeta, pero la mayor parte de lo que contiene no necesita estar separado. `ccam share` activa o desactiva el uso compartido para una cuenta y un elemento a la vez: el archivo o la carpeta de la cuenta se mueve a una carpeta de grupo y se reemplaza por un enlace simbólico, de modo que todas las cuentas del grupo leen y escriben la misma copia. Úsalo después de `ccam add` para cambiar lo que comparte una cuenta, o para poner algunas cuentas en un grupo propio.

## Uso

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Sin opciones de elementos, `ccam share` imprime el estado actual:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Con opciones, ccam imprime un plan, pregunta por cada archivo que es distinto de la copia del grupo y pregunta `Proceed? [y/N]` antes de cambiar nada. Una cuenta nueva comparte todo excepto el historial de conversaciones; [Compartir servidores MCP, configuración e historial](/docs/sharing) describe cada elemento desde el punto de vista del usuario.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Servidores MCP | sin cambios |
| `--settings <on\|off\|group>` | `settings.json` | sin cambios |
| `--memory <on\|off\|group>` | el `CLAUDE.md` global | sin cambios |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | sin cambios |
| `--history <on\|off\|group>` | `projects/` y `history.jsonl` | sin cambios |
| `--fresh` | al desactivar el uso compartido, empieza vacío en lugar de copiar el contenido del grupo | desactivada |
| `--prefer <group\|mine\|merge>` | cómo resolver los conflictos | pregunta; con `--yes`, la opción predeterminada de cada elemento |
| `-y`, `--yes` | se ejecuta sin preguntar; los conflictos siguen `--prefer` | desactivada |
| `-h`, `--help` | ayuda de `share` | |
| `--dry-run` (global) | se respeta: pregunta por los conflictos, imprime el plan y no cambia nada | desactivada |

El valor `on` significa el grupo `default`, `off` significa local (no compartido) y cualquier otro valor es un nombre de grupo. Los nombres de grupo se convierten en nombres de carpeta dentro de `~/.local/share/ccam/groups/`, así que `ccam share side --history ../x` falla con `invalid group name: "../x"`. La cuenta se indica por nombre o por número, y solo cambian los elementos que indicas.

## Cómo funciona

Los archivos que intervienen:

- `~/.config/ccam/config.toml`: la tabla `[profile.share]` de cada cuenta, con un valor por elemento.
- `~/.local/share/ccam/profiles/<name>/`: el `CLAUDE_CONFIG_DIR` de la cuenta.
- `~/.local/share/ccam/groups/<group>/`: las copias compartidas.

1. `internal/cli/share.go` lee la configuración, aplica cada valor de `--<item>` a una copia de la configuración de uso compartido de la cuenta con `config.Share.Set` y se niega si la cuenta está en ejecución.
2. `share.Build` compara la configuración actual con la de destino y lista los pasos. Solo lee el disco: `os.Lstat` y `os.Readlink` para cada ruta de `share.Layout`. El plan es una lista de operaciones: mover a un espacio vacío del grupo (adoptar), mover los archivos que el grupo no tiene (combinar carpeta), escribir en el grupo un conflicto resuelto, hacer una copia de seguridad, quitar un enlace, enlazar, volver a copiar, reescribir las rutas de los plugins y combinar el registro de plugins.
3. Los conflictos se resuelven, con la pregunta o con `--prefer`.
4. `Plan.Describe` imprime el plan, y ccam pregunta.
5. La configuración de destino se escribe en `config.toml` antes de mover nada. Si `Plan.Apply` se detiene a medias, `ccam doctor --fix` y una nueva ejecución construyen su plan hacia el destino guardado y terminan el cambio.
6. `Plan.Apply` ejecuta los pasos en orden. Cada paso mueve datos al grupo o a la carpeta de copia de seguridad; no se elimina nada.

`share.Layout` es una lista de rutas permitidas que cubre cada elemento. Nada que no esté en ella puede llegar nunca a un grupo. `share.Never` nombra `.credentials.json` y `.claude.json`, y `Build` se detiene con un error interno si alguno aparece alguna vez en la estructura. El elemento `mcp` no tiene rutas: la lista de servidores está dentro de `.claude.json`, así que [el uso compartido de MCP](/features/mcp) funciona mediante el `mcp.json` del grupo, y `--mcp` solo cambia la configuración.

Esto es lo que le pasa a una carpeta, `projects/`, cuando `side` activa el historial y el grupo ya tiene contenido:

<!-- motion:share-on -->

Si el grupo todavía no tiene nada en esa ruta, la cuenta es la primera en compartirla y su carpeta se mueve allí tal cual, sin copia de seguridad. En caso contrario, `mergeDir` mueve al grupo cada archivo que el grupo no tiene, y los archivos que existen en ambos lados se quedan en la carpeta de la cuenta, así que se conserva la versión del grupo. Después, la carpeta de la cuenta va a `.ccam-backup/<time>/` dentro de la carpeta de la cuenta, y `os.Symlink` pone un enlace en su lugar. Los archivos de conversación llevan el nombre de su ID de sesión, así que los historiales de varias cuentas se combinan sin conflictos.

Para mover se usa `rename(2)`, y entre sistemas de archivos distintos se recurre a copiar y eliminar. Un enlace relativo dentro del árbol movido que apuntaba fuera de él se reescribe como enlace absoluto (`fsutil.PortableLink`), para que siga resolviéndose desde la nueva ubicación. Los registros de `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) guardan rutas absolutas que pasan por la carpeta de la cuenta; ccam las apunta a la carpeta del grupo y, cuando el grupo ya tiene un registro, combina en él las entradas de la cuenta.

Desactivar un elemento invierte el enlace:

<!-- motion:share-off -->

ccam quita el enlace solo si apunta dentro de la carpeta de grupos, y después copia el contenido actual del grupo a la cuenta con `CopyTree`. Con `--fresh` se omite la copia y la cuenta empieza sin nada en esa ruta. El grupo conserva su contenido para las demás cuentas. Mover un elemento de un grupo a otro quita el enlace y lo enlaza con el grupo nuevo, sin tocar el grupo anterior; el plan lo indica en una línea `Note:`.

Un solo archivo que existe en ambos lados con contenido distinto es un conflicto:

<!-- motion:share-conflict -->

`share.NewConflict` lista las opciones adecuadas para el archivo. `history.jsonl` ofrece `merge by line`, `keep group` y `use mine`, con la combinación primero. Los demás archivos `.json` ofrecen `keep group`, `use mine` y `merge JSON`. Cualquier otro archivo ofrece `keep group` y `use mine`. La primera opción es la predeterminada:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) combina los objetos de forma recursiva, une los arreglos sin duplicados y conserva el valor del grupo cuando ambos lados asignan valores distintos a la misma clave, de modo que las demás cuentas no notan ningún cambio de comportamiento. `merge by line` conserva las líneas del grupo y luego agrega las tuyas que no tiene. Todas las opciones excepto `keep group` guardan primero la copia del grupo en `.ccam-backup/<time>/group/` dentro de la carpeta de la cuenta, y después escriben el resultado con `fsutil.AtomicWrite`. Tu propio archivo siempre va a la copia de seguridad.

`--prefer group`, `mine` o `merge` resuelve todos los conflictos sin preguntar; `merge` significa `merge JSON` o `merge by line`, la que ofrezca el archivo, y un archivo que no ofrece ninguna de las dos conserva la opción predeterminada. `--yes` solo toma la opción predeterminada de cada archivo.

En cada inicio, `share.Check` compara cada ruta de la estructura con la configuración usando solo `lstat` y `readlink`, lo que lo mantiene lo bastante ligero para la ruta de inicio. Un enlace que falta se crea en el momento, porque no hay nada que perder. Un archivo normal donde debería haber un enlace, un enlace al lugar equivocado o un enlace obsoleto detienen el inicio:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` ejecuta `share.Build` con la configuración guardada en ambos lados, lo que deriva cada paso de lo que hay en el disco, conserva el contenido del grupo en los conflictos y hace una copia de seguridad de todo lo que reemplaza.

## Seguridad y casos límite

- **Cuentas en ejecución.** ccam se niega a cambiar el uso compartido mientras Claude Code se está ejecutando con la cuenta (`account work is running; quit it first, then change its sharing`), porque la sesión sigue escribiendo en los archivos que se mueven y las líneas que agregue mientras tanto se perderían. `--dry-run` sigue imprimiendo el plan.
- **Copias de seguridad.** Cada plan tiene su propia carpeta `.ccam-backup/<time>/` dentro de la carpeta de la cuenta, con un sufijo numérico cuando dos planes empiezan en el mismo segundo. ccam nunca las elimina.
- **Fallos parciales.** Si un paso falla, ccam imprime el error y dónde están los originales, y te indica que termines con `ccam doctor --fix`. El destino ya está guardado, así que la reparación continúa hacia él.
- **Tus propios enlaces.** Si una ruta ya es un enlace que creaste, por ejemplo a un repositorio de dotfiles, y el grupo está vacío, ccam copia al grupo el destino del enlace y mueve el propio enlace a la copia de seguridad. El destino no se toca.
- **Archivos de inicio de sesión.** `.credentials.json` y `.claude.json` nunca se enlazan, mueven ni copian a un grupo.
- **Dos cuentas en una carpeta.** Las cuentas que comparten el historial escriben en los mismos archivos. Cuando una segunda cuenta del mismo grupo de historial se inicia en una carpeta donde ya se está ejecutando otra, ccam avisa y pregunta antes de continuar.

## Relacionado

- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Servidores MCP compartidos](/features/mcp)
- [Copiar una configuración existente](/features/migrate)
- [Todos los comandos](/docs/commands)
