---
title: Inicio rápido
description: Agrega dos cuentas de Claude Code con ccam add, cambia entre ellas con claude1, claude2 o ccam <name>, y elige qué comparten.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: c6b1c1c3ebaf
---

Agrega cada cuenta con `ccam add` y luego inicia Claude Code con `claude1`, `claude2` o `ccam <name>`. Eso es todo lo que la mayoría de las personas necesita.

## 1. Agrega tu primera cuenta

```sh
ccam add
```

El asistente te pide un nombre, como `personal` o `work`, y qué elementos debe compartir la cuenta con tus otras cuentas. Luego ejecuta `claude auth login`, que abre el navegador para que inicies sesión. Cuando termina, verás:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

Por ahora, ccam muestra sus propios mensajes en coreano; los comandos y las opciones son los mismos en todos los idiomas.

## 2. Agrega la segunda cuenta

```sh
ccam add work
```

Inicia sesión con la otra cuenta en el navegador. Esta se convierte en `claude2`.

## 3. Cambia de cuenta

Cualquiera de estos comandos inicia Claude Code con la cuenta de trabajo:

```sh
claude2
ccam work
ccam 2
```

Todo lo que va después del nombre de la cuenta pasa a Claude Code sin cambios, así que sigues usando las opciones que ya conoces:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Si ejecutas `ccam` solo, se abre un selector: usa las flechas y Enter, o presiona el número de la cuenta.

## 4. Elige qué compartir

De forma predeterminada, los servidores MCP, `settings.json`, el `CLAUDE.md` global, y los agentes, comandos, skills y plugins se comparten entre cuentas. El historial de conversaciones no. Actívalo en ambas cuentas si quieres continuar conversaciones entre ellas:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam muestra lo que va a mover antes de cambiar nada. Consulta [compartir](/docs/sharing) para ver los detalles.

## 5. Cuando una cuenta llega a su límite

```sh
ccam usage
ccam next --continue
```

`ccam usage` muestra el uso de 5 horas y semanal de todas las cuentas. `ccam next --continue` abre la misma conversación en la siguiente cuenta. Consulta [límites de uso](/docs/usage-limits).

## ¿Ya usas ~/.claude?

Regístrala en lugar de volver a iniciar sesión:

```sh
ccam import
```

La carpeta se queda donde está, y `claude` a secas sigue funcionando como antes.
