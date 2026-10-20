---
title: El panel
description: ccam en una terminal abre un panel a pantalla completa con cuentas, grupos, reglas de carpeta, extensiones y comprobaciones; cada cambio muestra antes su plan.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` sin argumentos abre un panel a pantalla completa cuando la entrada y la salida estándar son una terminal; si no, imprime `ccam list`. El panel muestra cada cuenta con su uso y sus grupos, y hace lo mismo que los comandos: iniciar, agregar, iniciar sesión, cambiar grupos, definir reglas de carpeta, comprobar y reparar. Al iniciar una cuenta se cierra; al agregar una cuenta o iniciar sesión, deja la terminal libre y luego vuelve.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

El primer fotograma se dibuja solo a partir de `config.toml`, en mucho menos de 100 ms; el uso se completa a medida que responde cada cuenta. Un porcentaje se vuelve amarillo a partir del 80 % y rojo al 100 %, y los valores tomados de una caché tras una consulta fallida se dibujan atenuados.

## Pestañas

| Pestaña | Qué muestra | Teclas |
| --- | --- | --- |
| Accounts | número, nombre, correo electrónico, barra de 5 horas, uso de 7 días, último uso, grupos; `✔` la predeterminada, `▸` la cuenta de esta carpeta, `•` en ejecución | más abajo |
| Groups | cada grupo, lo que comparte y sus cuentas, y las cuentas que no están en ningún grupo | `enter` editar, `n` nuevo, `e` renombrar, `x` quitar |
| Folders | las reglas de [`ccam folder`](/features/folder) y la que se aplica aquí | `n` usar una cuenta aquí, `enter` cambiar, `x` quitar |
| Extensions | lo que muestra [`ccam ext`](/features/ext), cada extensión una vez con los grupos de los que proviene | `r` volver a cargar |
| Health | lo que informa [`ccam doctor`](/features/doctor) | `f` reparar (`doctor --fix`) |

`←`/`→` o `Tab` cambian de pestaña, `?` lista las teclas de la pestaña actual, `r` vuelve a consultar el uso sin la caché corta, y `q` o `Esc` salen.

## Teclas en la pestaña Accounts

| Tecla | Qué hace |
| --- | --- |
| `Enter` | Inicia Claude Code con la cuenta seleccionada |
| `1`–`9` | Va a la cuenta con ese número; escribe `1` `2` rápido para la 12. Solo `Enter` inicia |
| `↑`/`↓`, `k`/`j` | Mueve la selección |
| `n` | Inicia la cuenta que elegiría [`ccam next --auto`](/features/next) |
| `o` | Inicia con opciones: `--continue`, `--resume`, otros argumentos o una ejecución de prueba |
| `u` | Detalles de uso: las ventanas de 5 horas, de 7 días y por modelo, con sus restablecimientos y de dónde vienen las cifras |
| `a` | Agrega una cuenta (el asistente de más abajo) |
| `l` | Vuelve a iniciar sesión |
| `s` | Elige los grupos de la cuenta |
| `d` | La convierte en la cuenta predeterminada |
| `f` | La usa en la carpeta actual |
| `e` | Le cambia el nombre |
| `c` | Define o quita su atajo |
| `x` | La quita (pregunta antes, con No seleccionado) |

## Cambiar grupos

En la pestaña Groups, `Enter` abre el editor: los elementos a la izquierda, las cuentas a la derecha, `Tab` para cambiar de lado y `Space` para marcar o desmarcar. Un elemento que una de las cuentas ya recibe de otro grupo aparece atenuado junto con el motivo, igual que una cuenta que recibiría un elemento dos veces; un elemento solo puede provenir de un grupo. `Enter` muestra el plan:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Primero se pregunta por cada archivo que difiere de la copia del grupo (conservar el del grupo, usar el mío o combinar). Después el cambio se ejecuta en segundo plano con el mismo código que [`ccam group`](/features/share): primero se guarda la configuración, los datos de cada cuenta se mueven en orden y los originales van a `.ccam-backup/`. Un cambio que movería los datos de una cuenta en ejecución se rechaza con el nombre de la cuenta.

## Agregar una cuenta

`a` abre el asistente que `ccam add` también usa en una terminal:

1. El nombre de la cuenta.
2. Cómo inicia sesión: con una suscripción de Claude ahora (en el navegador), con una [clave de API](/features/api) o con una suscripción de Claude más tarde.
3. Para una clave de API: de dónde lee ccam la clave en cada inicio (`env:NAME` o `command:COMMAND`), el endpoint, cómo se envía la clave y el modelo.
4. Los grupos a los que unirse, con `default` marcado; si todavía no hay grupos, lo que comparte el grupo `default`.
5. Un atajo, si quieres uno.
6. Un resumen para confirmar.

`Esc` vuelve un paso atrás y, en el primer paso, cierra el asistente.

## Colores

Los colores siguen los roles de la interfaz de terminal de Claude Code, así que ccam encaja a su lado; los nombres de grupo reciben un color según su nombre. ccam le pregunta a la terminal por su fondo y elige un tema claro u oscuro. Puedes fijar uno tú mismo en `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Los temas `-daltonized` son para el daltonismo rojo-verde y los temas `-ansi` usan los 16 colores propios de la terminal. `NO_COLOR` desactiva los colores.

## Seguridad y casos límite

- El panel llama a las mismas funciones que los comandos, así que se cumplen todas las reglas: planes antes de mover datos, rechazos para las cuentas en ejecución, copias de seguridad y tokens de inicio de sesión que nunca salen de una cuenta.
- El panel no tiene ejecución de prueba: `ccam --dry-run` sin un nombre imprime cómo se iniciaría la cuenta de esta carpeta, y `o` sobre una cuenta ofrece una ejecución de prueba de ese único inicio.
- La terminal se restaura antes de que se inicie Claude Code o se abra el inicio de sesión.
- En una pantalla demasiado pequeña para un diálogo, primero se ocultan las filas de la pestaña; las líneas se recortan con `…` en lugar de partirse.

## Relacionado

- [Compartir con grupos](/features/share)
- [Iniciar Claude Code](/features/launch)
- [Administrar cuentas](/features/accounts)
- [Límites de uso](/features/usage)
