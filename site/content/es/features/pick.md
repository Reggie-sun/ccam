---
title: Elegir una cuenta para scripts
description: "ccam pick imprime la cuenta con más margen para un script, da cuentas distintas a los workers iniciados a la vez y sale con 75 cuando ninguna tiene margen."
updated: 2026-10-25
order: 12
translation: machine
sourceHash: "5c4f74db9714"
---

`ccam pick` imprime el nombre de la cuenta con sesión iniciada que tiene más margen, para que un script inicie Claude Code con ella. Los workers iniciados a la vez reciben cada uno una cuenta distinta mientras haya suficientes cuentas con margen, así que un lote de trabajos `claude -p` se reparte entre tus cuentas en lugar de agotar una sola.

**Asientos de Team y Enterprise.** La elección automática (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y [`auto_switch`](/features/next#auto-switch) cuando lo activas) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por lo tanto, usar la elección automática entre asientos puede infringir los términos de Anthropic. Si la usas ahí es decisión tuya: consulta con tu organización y con Anthropic antes de depender de ella. Los planes personales (Pro, Max) se tratan en la [nota de las preguntas frecuentes sobre los términos](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Uso

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Códigos de salida: `0` con un nombre en la salida estándar, `75` (`EX_TEMPFAIL`) cuando ninguna cuenta tiene margen o ninguna tiene sesión iniciada, `1` para otros errores, como una opción incorrecta.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--json` | Imprime la cuenta, su uso de 5 horas y semanal, de dónde salen las cifras y por qué se eligió | desactivada |
| `--model <model>` | Cuenta también la ventana semanal de ese modelo, como `opus` | ninguno |
| `--max-usage <1-100>` | Prefiere las cuentas por debajo de este uso | `auto_max_usage`, si no 100 |
| `--strategy headroom\|expiring\|balance` | El menor uso de 5 horas, el restablecimiento de 5 horas más próximo o la cuota semanal que antes quedaría sin usar | `auto_strategy`, si no `headroom` |
| `--cross-org` | Elige también cuentas de otras organizaciones | `cross_org`, si no desactivada |
| `--wait[=<duration>]` | Cuando ninguna cuenta tiene margen, espera a que se restablezca la primera, como máximo este tiempo. El valor va después de `=`: `--wait 30m` se rechaza con la sugerencia `--wait takes its value after '=': --wait=30m` | desactivada; `5h` si se da sola |
| `--dry-run` | Elige sin registrar la elección | desactivada |

## Cómo funciona

1. ccam lista las cuentas con sesión iniciada (las cuentas con clave de API y los inicios de sesión de Console no tienen límites y quedan fuera) en la organización de la cuenta que iniciaría esta carpeta (consulta [organizaciones](/features/next#auto)); una carpeta con la sesión iniciada en la misma cuenta de Claude que esa también queda fuera.
2. Toma un bloqueo (`~/.local/state/ccam/picks.lock`), así que las elecciones que se ejecutan a la vez ocurren una tras otra.
3. Consulta el uso de todas las cuentas a la vez, igual que `ccam usage`, con la caché de 60 segundos. Una ventana que Claude Code informó llena cuenta como 100 % hasta su restablecimiento ([el registro de límite](/features/usage#the-limit-record)), así que una cuenta que acaba de llegar a su límite se queda fuera aunque sus cifras guardadas muestren margen.
4. Las cuentas elegidas en los últimos 2 minutos, por `ccam pick` o por [`auto_switch`](/features/next#auto-switch) al mover allí una sesión (`~/.local/state/ccam/picks/<name>`; una elección con fecha posterior a la hora del reloj no cuenta), se dejan fuera primero, y las [reglas de elección](/features/next#auto) de `ccam next --auto` eligen entre las demás: primero por debajo de `auto_max_usage`, según la estrategia, nunca una cuenta al 100 %. Cuando eso pasa por alto la cuenta que las reglas tomarían primero, el `reason` de `--json` lo indica: `30% of its 5 hours used; work was picked a moment ago`.
5. Si todas las cuentas con margen se eligieron hace poco, las reglas vuelven a elegir entre todas, así que un cuarto worker con tres cuentas comparte una.
6. La elección se registra y se libera el bloqueo; después se imprime el nombre.

La ventana de 2 minutos importa porque el uso de un trabajo nuevo tarda un rato en aparecer en las cifras: sin ella, cinco workers iniciados a la vez verían todos la misma cuenta como la de más margen.

## Seguridad y casos límite

- `ccam pick` solo lee el uso y escribe la hora de la elección; nunca renueva ni lee tokens más allá de lo que hace `ccam usage`.
- Una elección no reserva la cuenta. Dos scripts pueden seguir iniciándose en la misma cuenta pasados 2 minutos; usa `--fallback` para recuperarte cuando esa cuenta llegue a su límite.
- Esperar con `--wait` no mantiene ningún bloqueo, así que otras elecciones siguen mientras tanto.
- `ccam pick` ignora [`ccam folder`](/features/folder): elige solo por uso.

## Relacionado

- [Cambiar a la siguiente cuenta](/features/next)
- [Volver a ejecutar al llegar al límite](/features/next#run-again-at-the-limit)
- [Límites de uso](/features/usage)
- [Cuentas con clave de API](/features/api)
