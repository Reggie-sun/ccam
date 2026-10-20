---
title: Copiar una configuración existente
description: Copia la configuración, CLAUDE.md, las extensiones, los servidores MCP y sus inicios de sesión de un ~/.claude existente a un grupo de ccam con ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

Si usabas Claude Code antes de ccam, tu configuración está en `~/.claude` y `~/.claude.json`. `ccam migrate` la copia a un grupo de uso compartido, para que todas las cuentas que comparten esos elementos a través del grupo empiecen con la configuración que ya tienes. La carpeta de origen solo se lee, y `claude` a secas sigue funcionando como antes. Ejecútalo una vez después de tu primer `ccam add`, y de nuevo cada vez que una cuenta empiece a compartir servidores MCP.

## Uso

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

En una terminal, sin opciones de elementos y sin `--yes`, ccam primero muestra una lista de selección:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Espacio (o `x`) cambia el elemento bajo el cursor, las flechas o `j`/`k` mueven la selección, Enter muestra el plan, y Esc, `q` o Ctrl+C imprime `Cancelled.`. Los elementos para los que el origen no tiene nada aparecen marcados como `not found`. Si desmarcas todo, ccam imprime `Nothing selected.` y se detiene.

Después, ccam imprime el plan y pregunta antes de copiar:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Un grupo solo llega a las cuentas que comparten un elemento a través de él, así que el resumen las lista por elemento. Para las cuentas que no lo hacen, ccam imprime los comandos `ccam share <name> --<item> on` que activan el uso compartido.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `[folder]` | La carpeta de Claude Code desde la que se copia | `~/.claude` |
| `--settings` | Copia `settings.json` | desactivada |
| `--memory` | Copia el `CLAUDE.md` global | desactivada |
| `--extensions` | Copia `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | desactivada |
| `--mcp` | Copia los servidores MCP de usuario (`mcpServers` en `.claude.json`) | desactivada |
| `--history` | Copia el historial de conversaciones (`projects/`, `history.jsonl`) | desactivada |
| `-g`, `--group string` | El grupo de uso compartido al que se copia | `default` |
| `-y`, `--yes` | Copia sin preguntar; en un conflicto se conserva la copia del grupo | desactivada |
| `--dry-run` (global) | Imprime el plan y se detiene antes de cualquier pregunta o copia | desactivada |

Sin ninguna opción de elemento, en una terminal decide la lista de selección. Fuera de una terminal, o con `--yes`, ccam copia todo excepto el historial de conversaciones. El historial queda fuera de forma predeterminada porque `projects/` puede ser grande.

`--group on` significa `default`. `--group off` y `--group local` se rechazan con `--group takes a group name, not "off"`. `ccam migrate` respeta `--dry-run`, y la lista de selección sigue apareciendo en una ejecución de prueba para que puedas previsualizar una selección.

## Cómo funciona

`internal/cli/migrate.go` elige los elementos, y luego `migrate.Build` lee el origen y el grupo y devuelve un plan. No se escribe nada hasta `Plan.Apply`, después de que confirmas.

1. **Encontrar el origen.** `checkSource` rechaza una ruta que no existe, que no es una carpeta o que no contiene ninguna entrada conocida de Claude Code. `migrate.ClaudeJSON` busca el `.claude.json` que corresponde a la carpeta: dentro de ella o, para `~/.claude` sin uno dentro, `~/.claude.json` junto a ella. Ahí es donde lo guarda Claude Code cuando `CLAUDE_CONFIG_DIR` no está definida.
2. **Asignar elementos a rutas.** Cada ruta sale de `share.Layout`, la misma lista de rutas permitidas que usa `ccam share`. Los servidores MCP se leen de `.claude.json` con `mcp.LocalServers`, que usa el escáner `profile.TopLevelValue` en lugar de analizar el archivo completo.
3. **Copiar, nunca mover.** Los archivos van a `~/.local/share/ccam/groups/<group>/` con `fsutil.AtomicWrite` o `fsutil.CopyFile`; los servidores MCP van al `mcp.json` del grupo.

Cada elemento asigna una ruta de origen a una ruta del grupo:

<!-- motion:migrate-items -->

Las dos filas del historial son discontinuas porque solo se copian con `--history`. `.credentials.json` y `.claude.json` están en `share.Never` y nunca llegan al grupo: `Build` devuelve un error interno si alguno aparece en la estructura, y el recorrido de carpetas los omite por nombre.

`buildEntry` asigna un estado a cada ruta, que es lo que describe cada línea del plan:

<!-- motion:migrate-states -->

- **Missing**: no está en el origen, o es un archivo donde se espera una carpeta. Aparece una vez por elemento como `not found, skipped`.
- **Linked**: la ruta de origen ya apunta a la copia del grupo, por ejemplo cuando el origen es una carpeta de perfil de ccam.
- **Same**: el grupo ya tiene un contenido idéntico.
- **New**: el grupo no lo tiene, así que se copia completo.
- **Merge**: una carpeta en ambos lados. `copyMissing` agrega solo los archivos que el grupo no tiene; los archivos distintos se quedan como los tiene el grupo y el plan dice `group copies kept`. En `plugins/`, los registros `installed_plugins.json` y `known_marketplaces.json` se combinan, y las rutas absolutas de los plugins se reescriben para apuntar a la carpeta del grupo.
- **Conflict**: un archivo en ambos lados con contenido distinto. Eliges por archivo: `keep group`, `use mine`, y además `merge JSON` para archivos `.json` (una combinación profunda en la que gana el grupo cuando los valores chocan) o `merge by line` para `history.jsonl`. Un servidor MCP cuya definición es distinta también es un conflicto, con `keep group` o `use mine`. La primera opción de la lista es la predeterminada al pulsar Enter; `--yes` siempre conserva la copia del grupo.

La lista de selección es `tui.MigratePicker`, construida sobre el componente `checklist` que usa `ccam add`; `migrate.Found` marca los elementos que faltan con llamadas a lstat y la lista de servidores MCP.

### Inicios de sesión de MCP

Los servidores MCP que inician sesión con OAuth guardan sus tokens en el archivo de inicio de sesión, `.credentials.json` bajo `mcpOAuth`, no en la lista de servidores. Claude Code identifica cada entrada con `<name>|` más los primeros 16 caracteres de un hash `sha256` del `type`, la `url` y los `headers` del servidor. La clave no depende de la carpeta, así que una entrada copiada funciona en otra cuenta siempre que la definición del servidor sea idéntica. Por eso `buildLogins` en `internal/migrate/logins.go` solo copia las entradas de los servidores que el grupo termina teniendo con la definición del origen: nuevos, ya iguales o un conflicto en el que elegiste `use mine`. Los inicios de sesión de los servidores de plugins (claves que empiezan con `plugin:`) se copian cuando copias las extensiones y el origen tiene `plugins/`.

<!-- motion:migrate-signins -->

Las entradas van al `.credentials.json` propio de cada cuenta que las recibe, escrito con `fsutil.AtomicWrite` en el modo `0600`. Nunca van al grupo. Una cuenta las recibe solo si comparte `mcp` (o, para los servidores de plugins, `extensions`) a través del grupo de destino, y solo se agregan las claves que no tiene: las entradas que ya tiene, quizá más recientes, se conservan. La línea `sign-ins` del plan lista los servidores, las cuentas que los reciben y las omitidas. ccam omite una cuenta, y la nombra con el motivo, cuando:

- está registrada en su ubicación con `ccam import`, ya que `claude` a secas puede estar usando su archivo de inicio de sesión;
- está en ejecución, ya que Claude Code reescribe su archivo de inicio de sesión (se vuelve a comprobar justo antes de escribir, y se informa como `started meanwhile`);
- todavía no ha iniciado sesión;
- guarda su inicio de sesión en el Llavero de macOS, en el que ccam no escribe.

En una cuenta omitida, inicia sesión en los servidores con `/mcp`. Una cuenta que empieza a compartir servidores MCP más adelante recibe los inicios de sesión cuando vuelves a ejecutar `ccam migrate`; el resumen lo indica con ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Seguridad y casos límite

- El origen solo se lee. Ningún archivo en él se mueve, se enlaza ni se modifica, y el inicio de sesión del origen nunca se escribe, ni siquiera en macOS, donde puede leerse desde el Llavero.
- Antes de reemplazar un archivo del grupo, ya sea por `use mine`, una combinación o servidores MCP nuevos en un `mcp.json` existente, la copia del grupo se guarda en `<group>/.ccam-backup/<time>/`. El nombre de la carpeta es único en cada ejecución, así que dos ejecuciones en el mismo segundo no sobrescriben las copias de seguridad de la otra. El resumen imprime la ruta.
- Con `--yes`, cada conflicto conserva la copia del grupo, y ccam explica cómo elegir en su lugar: ejecuta el mismo comando sin `--yes`.
- Si no cambiaría nada, ccam imprime `Nothing to copy.`; si cada conflicto conserva la copia del grupo, `Nothing to copy: the group keeps its own copies.`
- Si una copia falla a medias, ccam se detiene con `The source folder was not changed. Group files replaced so far are saved in` seguido de la carpeta de copia de seguridad.
- Un origen que no existe termina con `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Relacionado

- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Comprobar y reparar](/features/doctor)
- [Varias cuentas de Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
