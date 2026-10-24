---
title: Copiar una configuración existente
description: Copia la configuración, CLAUDE.md, las extensiones y los servidores MCP de un ~/.claude existente a un grupo de ccam con ccam migrate.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: 7bca9efb6eb0
---

Si usabas Claude Code antes de ccam, tu configuración está en `~/.claude` y `~/.claude.json`. `ccam migrate` la copia a un grupo de uso compartido, para que todas las cuentas que comparten esos elementos a través del grupo empiecen con la configuración que ya tienes. La carpeta de origen solo se lee, y `claude` a secas sigue funcionando como antes. Ejecútalo una vez después de tu primer `ccam add`; cuando `~/.claude.json` tiene servidores MCP, el primer `ccam add` los nombra en su consejo.

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
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Un grupo solo llega a sus propias cuentas, y solo con los elementos que comparte, así que el resumen lista las cuentas del grupo. Cuando el grupo no comparte alguno de los elementos copiados, o algunas cuentas no están en él, ccam imprime los comandos `ccam group set <group> --add <items>` y `ccam group join <group> <accounts>` que lo cambian.

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

1. **Encontrar el origen.** `checkSource` rechaza una ruta que no existe, que no es una carpeta o que no contiene ninguna entrada conocida de Claude Code. Solo para los servidores MCP, un `~/.claude` vacío o inexistente sirve siempre que `~/.claude.json` tenga servidores. `migrate.ClaudeJSON` busca el `.claude.json` que corresponde a la carpeta: dentro de ella o, para `~/.claude` sin uno dentro, `~/.claude.json` junto a ella. Ahí es donde lo guarda Claude Code cuando `CLAUDE_CONFIG_DIR` no está definida.
2. **Asignar elementos a rutas.** Cada ruta sale de `share.Layout`, la misma lista de rutas permitidas que usan los grupos. Los servidores MCP se leen de `.claude.json` con `mcp.LocalServers`, que usa el escáner `profile.TopLevelValue` en lugar de analizar el archivo completo.
3. **Copiar, nunca mover.** Los archivos van a `~/.local/share/ccam/groups/<group>/` con `fsutil.AtomicWrite` o `fsutil.CopyFile`, que escribe un archivo temporal junto al destino y le cambia el nombre, así que una copia que se detiene nunca deja un archivo cortado. Los servidores MCP van al `mcp.json` del grupo con su bloqueo: un servidor que otro ccam agrega mientras tanto se conserva, y las copias hechas para [Remote Control](/guides/remote-control) se actualizan después.

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

Los servidores MCP que inician sesión con OAuth guardan sus tokens en el archivo de inicio de sesión, `.credentials.json` bajo `mcpOAuth`, no en la lista de servidores. ccam no los copia. Un inicio de sesión copiado tiene el mismo token de actualización que el original, y el servidor rota ese token en la primera renovación: todos los demás que lo tienen, incluido el `~/.claude` de origen, quedan entonces con la sesión cerrada en ese servidor. Por eso la línea `sign-ins` del plan nombra, para cada servidor copiado en el que el origen tiene la sesión iniciada, las cuentas del grupo que inician sesión por su cuenta con `/mcp`, y el resumen lo repite. Una entrada sin token no es un inicio de sesión y no se lista, ni tampoco una cuenta que ya tiene su propio inicio de sesión, ni un servidor cuya definición conserva el grupo. No se escribe ningún archivo de inicio de sesión.

## Seguridad y casos límite

- El origen solo se lee. Ningún archivo en él se mueve, se enlaza ni se modifica, y el inicio de sesión del origen nunca se escribe, ni siquiera en macOS, donde puede leerse desde el Llavero.
- Antes de reemplazar un archivo del grupo, ya sea por `use mine`, una combinación o servidores MCP nuevos en un `mcp.json` existente, la copia del grupo se guarda en `<group>/.ccam-backup/<time>/`. El nombre de la carpeta es único en cada ejecución, así que dos ejecuciones en el mismo segundo no sobrescriben las copias de seguridad de la otra. El resumen imprime la ruta.
- Con `--yes`, cada conflicto conserva la copia del grupo, y ccam explica cómo elegir en su lugar: ejecuta el mismo comando sin `--yes`.
- Si no cambiaría nada, ccam imprime `Nothing to copy.`; si cada conflicto conserva la copia del grupo, `Nothing to copy: the group keeps its own copies.`
- Si una copia falla a medias, ccam se detiene con `The source folder was not changed. Group files replaced so far are saved in` seguido de la carpeta de copia de seguridad. Una nueva ejecución copia los archivos que no llegaron completos.
- Cuando varios grupos comparten MCP, los comandos que sugieren `ccam doctor`, `ccam ext` y `ccam mcp list` son `ccam migrate --mcp --group <a|b>`.
- Un origen que no existe termina con `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Relacionado

- [Compartir servidores MCP, configuración e historial](/docs/sharing)
- [Comprobar y reparar](/features/doctor)
- [Varias cuentas de Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
