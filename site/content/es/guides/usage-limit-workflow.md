---
title: Qué hacer cuando Claude Code llega a su límite de uso
description: Revisa cuánto margen le queda a cada cuenta de Claude Code y continúa la misma conversación en otra cuenta cuando una llega a su límite.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

Cuando una cuenta llega a su límite de uso, continúa la misma conversación en otra cuenta tuya: comparte el historial de conversaciones una vez y luego ejecuta `ccam next --continue`. Si solo tienes una cuenta, espera a la hora de restablecimiento que muestra Claude Code.

## Antes de que pase

Configura las cuentas para que el cambio requiera un solo comando:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Ahora la línea de estado muestra la cuenta y su uso de 5 horas mientras trabajas, así que el límite no te toma por sorpresa.

## Cuando se alcanza el límite

Sal de Claude Code (`/exit`) y luego ejecuta:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Abre la conversación en la cuenta con más margen:

```sh
ccam next --auto --continue
```

Claude Code inicia como `work` con la conversación más reciente de la carpeta actual, incluido todo lo dicho hasta ahora.

## Ten en cuenta

- Ejecuta `ccam next --continue` en la misma carpeta en la que trabajabas; las conversaciones se guardan por carpeta de proyecto.
- Evita usar dos cuentas en la misma conversación al mismo tiempo. ccam te avisa si lo intentas.
- La otra cuenta también debe compartir el historial de conversaciones, en el mismo grupo.
- Si puedes cambiar entre cuentas de esta manera depende de los términos de cada cuenta.

Consulta [límites de uso](/docs/usage-limits) para ver todas las opciones de `ccam usage` y `ccam next`.
