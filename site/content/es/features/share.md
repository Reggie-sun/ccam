---
title: Compartir con grupos
description: Un grupo de ccam comparte lo elegido entre sus cuentas con una carpeta; ccam group mueve los datos con un plan, copias de seguridad y opciones ante conflictos.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: e4f05b0fd572
---

Cada cuenta tiene su propia carpeta, pero la mayor parte de lo que contiene no necesita estar separado. Un **grupo** dice qué se comparte y entre quiénes: cada cuenta del grupo enlaza esos elementos con una carpeta, `~/.local/share/ccam/groups/<group>/`, así que todas leen y escriben la misma copia. Los tokens de inicio de sesión nunca se comparten.

Una cuenta puede estar en varios grupos, pero cada elemento proviene de un solo grupo. Una configuración típica es un grupo para la configuración y las extensiones, y un segundo grupo para el historial de conversaciones entre dos cuentas:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

La primera cuenta que agregas crea `default`, y las cuentas posteriores se unen a él salvo que indiques otra cosa.

## Uso

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Todo comando que mueve datos imprime un plan, pregunta por cada archivo que es distinto de la copia del grupo y pregunta `Proceed? [y/N]` antes de cambiar nada:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

Los mismos cambios están en la pestaña **Groups** del [panel](/features/dashboard), y `s` sobre una cuenta elige sus grupos.

## Elementos

| Elemento | Qué cubre |
| --- | --- |
| `mcp` | Servidores MCP, mediante el `mcp.json` del grupo ([Servidores MCP compartidos](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | el `CLAUDE.md` global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` y `history.jsonl`, para que otra cuenta pueda usar `--continue` |

Un elemento que ningún grupo de la cuenta comparte se queda en la carpeta propia de la cuenta.

## Comandos y opciones

| Comando | Qué hace |
| --- | --- |
| `ccam group` o `ccam group list` | lista los grupos |
| `ccam group show <group>` | lo que comparte el grupo, sus cuentas, su carpeta y sus servidores MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | crea un grupo; los elementos de sus cuentas se mueven a él |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | cambia lo que comparte el grupo, para todas sus cuentas |
| `ccam group join <group> <account>...` | pone cuentas en el grupo |
| `ccam group leave <group> <account>...` | saca cuentas del grupo; cada una conserva una copia, incluidos los servidores MCP |
| `ccam group rename <group> <new-name>` | cambia el nombre solo en la configuración |
| `ccam group remove <group> [--leave]` | quita un grupo |

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: la cuenta cuya copia entra primero en el grupo | la cuenta predeterminada, luego por número |
| `--replace` | `join`: sale antes de los grupos que comparten los mismos elementos | desactivada |
| `--fresh` | las cuentas que dejan de compartir un elemento empiezan sin él en lugar de con una copia (una carpeta importada conserva una de todos modos) | desactivada |
| `--leave` | `remove`: saca antes las cuentas del grupo | necesaria si tiene cuentas |
| `--prefer <group\|mine\|merge>` | cómo resolver los conflictos | pregunta; con `--yes`, la opción predeterminada de cada archivo |
| `-y`, `--yes` | se ejecuta sin preguntar; los conflictos siguen `--prefer` | desactivada |
| `--dry-run` (global) | imprime el plan y no cambia nada | desactivada |

Las cuentas se indican por nombre o por número; los elementos, como una lista separada por comas. `ccam add` acepta `--group <group>` (se puede repetir) o `--no-group`.

## Cómo funciona

Los archivos que intervienen:

- `~/.config/ccam/config.toml`: las entradas `[[group]]`: `name`, `items`, `members` y, cuando el nombre de la carpeta es distinto, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: el `CLAUDE_CONFIG_DIR` de la cuenta.
- `~/.local/share/ccam/groups/<folder>/`: las copias compartidas.

Al leer `config.toml`, los grupos de cada cuenta se resuelven en una carpeta por elemento. Esto tarda unos 25 µs con cinco cuentas, así que un inicio sigue leyendo solo `config.toml`. Una configuración en la que una cuenta recibiría un elemento de dos grupos es un error, y no se inicia nada hasta que se corrija.

Un cambio trabaja sobre una copia de la configuración:

1. El comando edita la copia (`AddGroup`, `SetGroupItems`, `JoinGroup` …), que comprueba las reglas.
2. `internal/group` compara la configuración anterior con la nueva y, para cada cuenta cuyos elementos cambian, construye un plan con `share.Build`. Los planes se construyen uno tras otro, y cada uno ve lo que los anteriores pondrán en las carpetas de grupo, así que solo la primera cuenta mueve su copia al grupo. La primera copia viene de la carpeta del grupo si ya tiene una, después de `--from`, después de la cuenta predeterminada y después de las demás por número.
3. ccam se niega si una de esas cuentas está en ejecución, sin importar cómo se haya iniciado Claude Code (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), y en Windows sin el modo de desarrollador si un archivo necesitaría un enlace.
4. Los conflictos se resuelven, con la pregunta o con `--prefer`, se imprime el plan y ccam pregunta.
5. Justo antes de guardar, ccam vuelve a comprobar: una cuenta que se inició mientras esperaban las preguntas detiene el cambio, y también lo hace un cambio de los mismos grupos por parte de otro ccam (`the groups changed while this plan waited; nothing was moved, make the change again`). La nueva configuración se guarda antes de mover nada, y el cambio queda anotado en `~/.local/state/ccam/group-change.json`. Si un paso se detiene a medias, la nota se queda: `ccam doctor --fix` termina el cambio en su orden, con las respuestas que diste, incluidos los pasos de MCP, y hasta entonces no empieza ningún otro cambio de grupos (``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``).
6. Los pasos de MCP se ejecutan primero. Cuando las cuentas sacan las extensiones de un grupo, las listas de plugins del grupo se apuntan a la carpeta del grupo allí donde sus rutas pasan por la carpeta de una cuenta (`group plugin lists:` en el plan). Después los planes se ejecutan en orden. Cada paso mueve datos al grupo o a la carpeta de copia de seguridad; solo los archivos que el grupo ya tiene idénticos byte a byte se descartan en lugar de copiarse a la copia de seguridad.

`share.Layout` es una lista de rutas permitidas que cubre cada elemento. Nada que no esté en ella puede llegar nunca a un grupo. `share.Never` nombra `.credentials.json` y `.claude.json`, y `Build` se detiene con un error interno si alguno aparece alguna vez en la estructura. El elemento `mcp` no tiene rutas: la lista de servidores está dentro de `.claude.json`, así que [el uso compartido de MCP](/features/mcp) funciona mediante el `mcp.json` del grupo.

Esto es lo que le pasa a una carpeta, `projects/`, cuando `side` empieza a compartir el historial y el grupo ya tiene contenido:

<!-- motion:share-on -->

Si el grupo todavía no tiene nada en esa ruta, la cuenta es la primera en compartirla y su carpeta se mueve allí tal cual. En caso contrario, cada archivo que el grupo no tiene se mueve al grupo, y los archivos que existen en ambos lados se quedan en la carpeta de la cuenta, así que se conserva la versión del grupo. Los archivos idénticos a los del grupo se descartan, y el plan nombra los que difieren (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`). Después, la carpeta de la cuenta va a `.ccam-backup/<time>/` dentro de la carpeta de la cuenta, y un enlace ocupa su lugar. Si la entrada del grupo es a su vez un enlace a la entrada propia de la cuenta que se une, el contenido de la cuenta va al grupo en su lugar. Los archivos de conversación llevan el nombre de su ID de sesión, así que los historiales de varias cuentas se combinan sin choques. Los registros de `plugins/` guardan rutas absolutas que pasan por la carpeta de la cuenta; ccam las apunta a la carpeta del grupo y combina las entradas de la cuenta en un registro que el grupo ya tenga.

Salir de un grupo invierte el enlace:

<!-- motion:share-off -->

Primero ccam copia el contenido actual del grupo en `.ccam-backup/incomplete-<name>` dentro de la carpeta de la cuenta, y reemplaza el enlace solo cuando la copia está completa, así que una copia que se detiene deja el enlace en su sitio y `ccam doctor --fix` vuelve a copiar (el plan de esa finalización muestra `restore <entry>: … (copied back from the group before a stop)`). Quita el enlace solo si sigue apuntando dentro de la carpeta de grupos como cuando se hizo el plan; una ruta que cambió desde entonces, como un archivo normal donde estaba el enlace, se deja como está y el cambio se detiene. Con `--fresh` se omite la copia. El grupo conserva su contenido para las demás cuentas. Una cuenta que mueve un elemento a un grupo cuya carpeta todavía está vacía lleva consigo una copia de lo que veía, así que conserva sus archivos; esa copia se prepara en `<group>/.ccam-incomplete-<name>` hasta que está completa. Una cuenta que deja de compartir MCP conserva los servidores del grupo: van a su propio `.claude.json` (se guarda una copia de seguridad), junto a sus propios servidores con los mismos nombres, que se mantienen. Con `--fresh`, en su lugar, el plan enumera los servidores que dejan de cargarse. Volver a unirse quita las copias que nadie cambió. Dejar de compartir `extensions` también apunta las listas de plugins del grupo a la carpeta del grupo donde sus rutas pasaban por la carpeta de la cuenta que sale; para las demás cuentas significan lo mismo.

Un solo archivo que existe en ambos lados con contenido distinto es un conflicto:

<!-- motion:share-conflict -->

Tú eliges en los archivos sueltos como `settings.json`; dentro de las carpetas se queda la copia del grupo, y el plan nombra los archivos que difieren. `history.jsonl` ofrece `merge by line`, `keep group` y `use mine`, con la combinación primero. Los demás archivos `.json` ofrecen `keep group`, `use mine` y `merge JSON`. Cualquier otro archivo ofrece `keep group` y `use mine`. `merge JSON` combina los objetos de forma recursiva, une los arreglos sin duplicados y conserva el valor del grupo cuando ambos lados asignan valores distintos a la misma clave, de modo que las demás cuentas no notan ningún cambio. Todas las opciones excepto `keep group` guardan primero la copia del grupo en `.ccam-backup/<time>/group/` de la cuenta. Tu propio archivo siempre va a la copia de seguridad.

En cada inicio, `share.Check` compara cada ruta de la estructura con los grupos de la cuenta usando solo `lstat` y `readlink`. Un enlace que falta se crea en el momento. Un archivo normal donde debería haber un enlace, un enlace al lugar equivocado o un enlace obsoleto detienen el inicio y remiten a `ccam doctor --fix`.

### Cambiar el nombre, quitar y el formato anterior

- `ccam group rename` solo cambia la configuración. La carpeta conserva su nombre y el grupo lo registra como `dir`, así que no se mueve nada y no se puede perder nada.
- `ccam group remove` se niega a quitar un grupo con cuentas salvo que agregues `--leave`. Entonces cada cuenta conserva una copia, y la carpeta va a `~/.local/share/ccam/trash/` salvo que otro grupo la use.
- Antes de los grupos, `config.toml` guardaba una tabla `[profile.share]` por cuenta. ccam lee ese formato como grupos sin mover nada. Al guardar por primera vez escribe el formato nuevo y conserva el archivo anterior como `config.toml.bak`. Los elementos compartidos por las mismas cuentas se convierten en un grupo con el nombre de la carpeta; los demás reciben nombres como `default-history` y siguen usando la misma carpeta.

### En Windows

Windows crea enlaces simbólicos solo en el modo de desarrollador (Configuración > Sistema > Para programadores) o para un administrador. Sin ese permiso, una carpeta compartida se convierte en una unión de directorio (junction), que funciona igual para Claude Code. Un archivo compartido (`settings.json`, `CLAUDE.md`, `history.jsonl`) solo puede ser un enlace simbólico, así que un cambio de grupo que necesite uno se rechaza antes de mover nada. Por el mismo motivo, allí el primer `ccam add` crea el grupo `default` solo con `mcp` y `extensions` y lo indica, y se rechaza `ccam add --group` hacia un grupo que comparte la configuración, la memoria o el historial.

## Seguridad y casos límite

- **Cuentas en ejecución.** Un cambio se rechaza mientras Claude Code se está ejecutando con una de las cuentas que mueve, porque la sesión sigue escribiendo en esos archivos, la haya iniciado ccam o no (`claude` a secas, un IDE, sesiones en segundo plano). La comprobación se repite justo antes de guardar. Salir de un grupo o unirse a él solo afecta a la cuenta indicada; dejar de compartir `extensions` también reescribe rutas en las listas de plugins del grupo, sin cambiar lo que significan para las demás. `--dry-run` sigue imprimiendo el plan.
- **Copias de seguridad.** Cada plan tiene su propia carpeta `.ccam-backup/<time>/` dentro de la carpeta de la cuenta. ccam nunca las elimina.
- **Tus propios enlaces.** Si una ruta ya es un enlace que creaste, por ejemplo a un repositorio de dotfiles, y el grupo está vacío, ccam copia al grupo el destino del enlace y mueve el propio enlace a la copia de seguridad. El destino no se toca.
- **Carpetas importadas.** Una cuenta registrada en su ubicación con `ccam import` (como `~/.claude`) cambia en su sitio, así que `claude` a secas también ve el cambio; el plan lo indica. Cuando sale de un grupo siempre conserva una copia, también con `--fresh`, y `ccam remove` la saca así de sus grupos antes de olvidar la cuenta. Sus archivos compartidos viven en la carpeta del grupo mientras está en un grupo, así que quita la cuenta (o sal de los grupos) antes de eliminar `~/.local/share/ccam`.
- **Dos cuentas en una carpeta.** Las cuentas que comparten el historial escriben en los mismos archivos. Cuando una segunda cuenta se inicia en una carpeta donde ya se está ejecutando otra, ccam avisa y pregunta.
- **Ediciones a mano.** Un miembro que no es una cuenta se ignora y `ccam doctor` lo informa. Los nombres `local`, `on`, `off` y `none` no pueden ser nombres de grupo.

## Relacionado

- [El panel](/features/dashboard)
- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Servidores MCP compartidos](/features/mcp)
- [Copiar una configuración existente](/features/migrate)
- [Todos los comandos](/docs/commands)
