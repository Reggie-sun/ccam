---
title: Límites de uso y cambio de cuenta
description: Consulta el uso de 5 horas y semanal de cada cuenta de Claude Code con ccam usage, y abre la cuenta con más margen con ccam next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 8ab41501c548
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

Este endpoint no está documentado públicamente y puede cambiar. Cuando una solicitud falla, ccam muestra las últimas cifras conocidas y las marca como almacenadas en caché. `ccam doctor` te avisa si el formato de la respuesta cambió.

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

Entonces Claude Code muestra el nombre de la cuenta y su uso en la línea de estado, por ejemplo `[work] · Opus · 5h 24% · 7d 41%`. Si ya tienes un comando de línea de estado, ccam lo conserva y pone el nombre de la cuenta delante. La línea de estado también guarda sus cifras, que `ccam usage` usa cuando el endpoint no responde.

## Cambiar de cuenta

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` omite las cuentas que no tienen sesión iniciada, y `--auto` omite las cuentas al 100 %. Para continuar una conversación, ambas cuentas deben compartir el historial de conversaciones; consulta [compartir](/docs/sharing).

La [guía del flujo de trabajo con límites de uso](/guides/usage-limit-workflow) recorre un ejemplo completo.
