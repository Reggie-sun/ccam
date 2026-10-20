---
title: Sesiones en segundo plano de todas las cuentas
description: ccam agents lista en una tabla los agentes en segundo plano y las sesiones abiertas de Claude Code de todas tus cuentas, y ccam attach abre uno por su id.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code mantiene las sesiones en segundo plano (`claude --bg`, la vista de agentes) por cuenta: `claude agents` solo muestra las sesiones de la cuenta con la que se ejecuta. `ccam agents` consulta todas las cuentas a la vez y las lista juntas, y `ccam attach <id>` abre una sesión en segundo plano con la cuenta que la inició, sea cual sea.

## Uso

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` es el estado de la sesión en segundo plano (`working`, `blocked`, `done`, `failed`, `stopped`) y lo que está haciendo su proceso (`busy`, `waiting`, `idle`); una sesión interactiva abierta muestra `open`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--all` | Lista también las sesiones en segundo plano terminadas, fallidas o detenidas | desactivado |
| `--json` | Imprime las sesiones como JSON, cada una con la cuenta (`account`) a la que pertenece | desactivado |

`ccam attach <id> [claude args]` pasa a `claude attach` los argumentos que van después del id. Con `--dry-run` imprime el comando.

## Cómo funciona

1. Para cada cuenta, ccam ejecuta `claude agents --json` (con `--all` si se indica) con el `CLAUDE_CONFIG_DIR` de la cuenta, todas a la vez, cada una con un límite de 15 segundos. Esto lee la lista de sesiones propia de la cuenta y no inicia el servicio en segundo plano de Claude Code.
2. Cada entrada recibe el nombre de la cuenta, y la lista se ordena por hora de inicio, de la más reciente a la más antigua. Una cuenta que no responde se nombra en la salida de error estándar y se deja fuera, y ccam sale entonces con el código 1, para que un script no confunda un fallo con una lista vacía.
3. `ccam attach` busca `jobs/<id>` en la carpeta de cada cuenta, donde Claude Code guarda una sesión en segundo plano, y ejecuta `claude attach <id>` con esa cuenta, de la misma forma en que `ccam <name>` inicia Claude Code.

## Seguridad y casos límite

- Los ids de las sesiones en segundo plano pertenecen a una sola cuenta: `claude attach` ejecutado con otra cuenta no puede encontrarlas, y eso es lo que resuelve `ccam attach`.
- Dos cuentas podrían tener sesiones cuyos ids empiezan igual; en ese caso `ccam attach` nombra ambas y pide más caracteres del id.
- ccam pasa sus servidores MCP compartidos a Claude Code con `--mcp-config=<file>` en un solo argumento, que `claude agents` y `claude attach` aceptan; las sesiones en segundo plano siguen usando ese archivo cuando se reinician.
- El servicio en segundo plano de una cuenta con clave de API toma la clave del shell que lo inició. Inicia la primera sesión en segundo plano de una cuenta así a través de ccam (`ccam gw --bg …`), para que el servicio tenga la clave.

## Relacionado

- [Iniciar Claude Code](/features/launch)
- [Elegir una cuenta para scripts](/features/pick)
- [El panel](/features/dashboard)
