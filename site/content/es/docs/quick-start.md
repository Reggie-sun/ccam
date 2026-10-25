---
title: Inicio rápido
metaTitle: "Inicio rápido: agrega y cambia cuentas de Claude Code"
description: Agrega dos cuentas de Claude Code con ccam add, dales atajos como ccwork, cambia entre ellas y elige qué comparten.
updated: 2026-10-25
order: 2
translation: machine
sourceHash: "0570ef27e4ce"
---

Agrega cada cuenta con `ccam add` y luego inicia Claude Code con `ccam <name>` o con un atajo al que tú le pones nombre, como `ccwork`. Es todo lo que la mayoría necesita.

## 1. Agrega tu primera cuenta

```sh
ccam add
```

El asistente te pide un nombre, como `personal` o `work`, qué elementos debe compartir la cuenta con tus otras cuentas y si quieres crear un atajo (consulta el paso 4). Luego ejecuta `claude auth login`, que abre el navegador para que inicies sesión. En una máquina remota, el navegador puede terminar en una página `http://localhost:…/callback` que no carga; pega esa dirección en la terminal y ccam completa el inicio de sesión. Cuando termina, verás:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Trae tu configuración actual

Si ya usabas Claude Code, tu configuración, `CLAUDE.md`, tus extensiones y tus servidores MCP están en `~/.claude`. Cópialos a la configuración compartida para que todas las cuentas los tengan:

```sh
ccam migrate
```

ccam muestra lo que va a copiar y pregunta antes. Solo lee `~/.claude`, nunca copia un inicio de sesión y deja fuera el historial de conversaciones a menos que agregues `--history`. Consulta [copiar una configuración existente](/docs/sharing#copy-an-existing-setup) para ver los detalles.

## 3. Agrega la segunda cuenta

```sh
ccam add work
```

Inicia sesión con la otra cuenta en el navegador.

## 4. Cambia de cuenta

`ccam work` inicia Claude Code con la cuenta work. Para algo más corto, dale a la cuenta su propio atajo, como `ccwork`. Después de los elementos compartidos, el asistente pregunta:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Presiona `y` y escribe un nombre, o presiona Enter para usar el sugerido:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Después del inicio de sesión, ccam lo confirma:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Si te saltaste la pregunta, puedes agregar un atajo en cualquier momento:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Ahora cualquiera de estos comandos inicia Claude Code con la cuenta work:

```sh
ccwork
ccam work
ccam 2
```

Si ejecutas `ccam` solo, se abre el [panel](/features/dashboard): ve a una cuenta con las flechas o con su número y presiona Enter.

Todo lo que va después del atajo o del nombre de la cuenta pasa a Claude Code sin cambios, así que sigues usando las opciones que ya conoces:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Elige qué compartir

De forma predeterminada, los servidores MCP, `settings.json`, el `CLAUDE.md` global, y los agentes, comandos, skills y plugins se comparten entre cuentas a través del grupo `default`. El historial de conversaciones no. Agrégalo al grupo si quieres continuar conversaciones entre cuentas:

```sh
ccam group set default --add history
```

ccam muestra lo que va a mover en cada cuenta antes de cambiar nada. Consulta [compartir con grupos](/features/share) para ver los detalles.

## 6. Cuando una cuenta llega a su límite

```sh
ccam usage
ccam next --continue
```

`ccam usage` muestra el uso de 5 horas y semanal de todas las cuentas. `ccam next --continue` abre la misma conversación en la siguiente cuenta. En scripts, `ccam --fallback work -p "…"` cambia por sí solo. Consulta [límites de uso](/docs/usage-limits).

**Asientos de Team y Enterprise.** La elección automática (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y [`auto_switch`](/features/next#auto-switch) cuando lo activas) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por lo tanto, usar la elección automática entre asientos puede infringir los términos de Anthropic. Si la usas ahí es decisión tuya: consulta con tu organización y con Anthropic antes de depender de ella. Los planes personales (Pro, Max) se tratan en la [nota de las preguntas frecuentes sobre los términos](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## ¿Ya usas ~/.claude?

Regístrala en lugar de volver a iniciar sesión:

```sh
ccam import
```

La carpeta se queda donde está, y `claude` a secas sigue funcionando como antes.
