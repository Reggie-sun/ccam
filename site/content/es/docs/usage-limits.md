---
title: Límites de uso y cambio de cuenta
description: Consulta el uso de 5 horas y semanal de cada cuenta de Claude Code con ccam usage, y abre la cuenta con más margen con ccam next --auto.
updated: 2026-10-20
order: 5
translation: machine
sourceHash: 40aef0708c81
---

`ccam usage` muestra cuánto ha usado cada cuenta de sus límites de 5 horas y semanal. `ccam next --auto` abre la cuenta con más margen, y `--continue` lleva consigo la conversación actual.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

Las columnas son: cuenta, uso de 5 horas, tiempo hasta que se restablece, uso semanal, tiempo hasta que se restablece y la antigüedad de las cifras.

## De dónde salen las cifras

ccam consulta el endpoint de uso de Anthropic para todas las cuentas al mismo tiempo, con el inicio de sesión propio de cada cuenta. Espera como máximo 3 segundos por cuenta y muestra cada cuenta en cuanto responde.

Este endpoint no está documentado públicamente y puede cambiar. Cuando una solicitud falla, ccam muestra las últimas cifras conocidas y las marca como almacenadas en caché. Una ventana que se haya reiniciado desde que se guardaron esas cifras muestra 0%. `ccam doctor` te avisa si el formato de la respuesta cambió.

Las cifras obtenidas hace menos de 60 segundos se reutilizan. Opciones:

| Opción | Efecto |
| --- | --- |
| `--watch [seconds]` | Actualiza de forma continua, cada 60 segundos de forma predeterminada |
| `--json` | Salida legible por máquinas |
| `-v` | Muestra también los límites semanales por modelo y de dónde salió cada cifra |
| `--no-cache` | Siempre vuelve a consultar |
| `--source cache` | Nunca consulta; solo muestra las cifras guardadas |

## Inicios de sesión vencidos

Cuando el token de acceso de una cuenta venció, ccam no lo renueva, porque renovarlo mientras una sesión de Claude Code usa el mismo inicio de sesión podría cerrar esa sesión. Inicia la cuenta una vez con `ccam <name>` y Claude Code lo renueva.

`ccam usage --refresh` renueva por sí mismo los tokens vencidos, pero solo de las cuentas que no se están ejecutando. Nunca renueva inicios de sesión guardados en el Llavero de macOS ni en un `~/.claude` importado.

## Uso en la línea de estado

```sh
ccam statusline install
```

Entonces Claude Code muestra el nombre de la cuenta y su uso en la línea de estado, por ejemplo `[work] · Opus · 5h 24% · 7d 41%`. A partir del 90 %, la línea también sugiere `ccam next --auto --continue`. Si ya tienes un comando de línea de estado, ccam lo conserva y pone el nombre de la cuenta delante. La línea de estado también guarda sus cifras, que `ccam usage` usa cuando el endpoint no responde.

## Cambiar de cuenta

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` omite las cuentas que no tienen sesión iniciada, y `--auto` omite las cuentas al 100 % y las cuentas de otra organización. Cuando las dos cuentas comparten el historial de conversaciones, `--continue` abre la misma conversación; si no, Claude Code la copia en la nueva cuenta. Consulta [compartir](/docs/sharing) y [cambiar a la siguiente cuenta](/features/next#continue).

Para scripts y tareas programadas que ejecutan `claude -p`, pon `--fallback` antes del nombre de la cuenta. Cuando esa cuenta llega a su límite, ccam vuelve a ejecutar el mismo comando con la cuenta con más margen, y solo se imprime la salida de la ejecución que terminó:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Consulta [Volver a ejecutar al llegar al límite](/features/next#run-again-at-the-limit) para ver cómo decide ccam.

La [guía del flujo de trabajo con límites de uso](/guides/usage-limit-workflow) recorre un ejemplo completo.
