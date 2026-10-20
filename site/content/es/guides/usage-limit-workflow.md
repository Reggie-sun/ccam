---
title: Qué hacer cuando Claude Code llega al límite de uso
description: Espera al reinicio o continúa al instante la misma conversación en otra cuenta de Claude Code. Cómo funcionan los límites de 5 horas y semanal, paso a paso.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 38b80d780bba
---

Cuando una cuenta llega a su límite de uso (a menudo llamado rate limit), tienes dos opciones: esperar a que el límite se reinicie o continuar ahora la misma conversación en otra cuenta tuya. Con ccam, lo segundo es un solo comando, `ccam next --auto --continue`, una vez que el historial de conversaciones se comparte entre las cuentas.

## ¿Qué límites tiene Claude Code?

En los planes Pro y Max, Claude Code tiene un [límite de sesión de cinco horas y un límite semanal](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Ambos se comparten con Claude en la web y en las apps, y el uso cuenta para los dos a la vez. Cuando uno se agota, Claude Code [muestra un mensaje](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) como "You've hit your session limit · resets 3:45pm". Cambiar de modelo no ayuda, porque los límites cubren todos los modelos.

## Opción 1: esperar al reinicio

Si tienes una sola cuenta, espera. Las versiones recientes de Claude Code [esperan en la sesión abierta](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) y continúan la tarea por su cuenta cuando el límite se reinicia. `/usage` dentro de Claude Code muestra cuánto has usado de los límites de tu plan.

## Opción 2: continuar ahora en otra cuenta

Si tienes una segunda cuenta, puedes seguir de inmediato con la misma conversación. Configúralo una vez:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Ahora la línea de estado muestra la cuenta y su uso de 5 horas mientras trabajas, así que el límite no te toma por sorpresa.

Cuando llegues al límite, sal de Claude Code (`/exit`) y luego:

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

Claude Code se inicia como `work` con la conversación más reciente de la carpeta actual, incluido todo lo dicho hasta ahora. Compartir el historial (`ccam group set default --add history` arriba) hace que siga siendo una sola conversación; sin él, Claude Code copia la conversación en `work`, y ccam te indica cuántos tokens vuelve a leer esa primera respuesta.

Con el [plugin de ccam](/docs/install#claude-code-plugin) en Claude Code, una sesión que se detiene en su límite muestra una notificación, y `/ccam:switch` (o `/ccam:switch work`) anota la conversación para `ccam next --continue` sin dejar nada a la adivinación.

## ¿ccam cambia de cuenta automáticamente?

No. ccam no se está ejecutando mientras trabajas: inicia Claude Code y se reemplaza a sí mismo por él. Tú decides cuándo cambiar, y `ccam next --auto` elige la cuenta con el menor uso de 5 horas y omite las cuentas que están al 100%, sin sesión iniciada o en otra organización. Con `relay = true` en `config.toml`, ccam espera en cambio al lado de la sesión, así que `/ccam:switch` continúa por sí solo. Otras herramientas cambian el inicio de sesión de una sesión en curso; [ccam y otros selectores de cuentas](/compare/account-switchers) explica la diferencia.

## Qué tener en cuenta

- Ejecuta `ccam next --continue` en la misma carpeta en la que trabajabas; las conversaciones se guardan por carpeta de proyecto.
- Evita ejecutar dos cuentas sobre la misma conversación al mismo tiempo. ccam te avisa si lo intentas.
- Sin historial compartido, la conversación se copia a la otra cuenta y las dos copias siguen su propio camino; los puntos de control desde los que `/rewind` restaura el código se quedan con la primera cuenta.
- El primer mensaje en la cuenta nueva puede consumir más de su límite de lo habitual, porque esa cuenta todavía no tiene la conversación en caché.
- Si puedes cambiar de cuenta de esta forma depende de los términos de cada cuenta.

Consulta [límites de uso](/docs/usage-limits) para ver todas las opciones de `ccam usage` y `ccam next`.
