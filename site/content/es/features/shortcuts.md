---
title: Atajos
description: Dale a una cuenta su propio comando, como ccwork, con ccam shortcut, y vuelve a crear los enlaces en ~/.local/bin con ccam shims install.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

Un atajo es un comando al que tú le pones nombre y que inicia Claude Code con una cuenta concreta: `ccwork --resume` hace lo mismo que `ccam work --resume`. Es un enlace a ccam en `~/.local/bin`, así que funciona en bash, zsh, scripts y cualquier otra cosa que ejecute programas desde `PATH`. Los atajos son opcionales; crea uno para las cuentas que inicias a menudo.

## Uso

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam confirma cada cambio:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Una cuenta tiene como máximo un atajo, guardado como `shortcut` en `~/.config/ccam/config.toml`. El asistente de `ccam add` pregunta `Create a shortcut command for <name>? [y/N]` y luego `Shortcut name [cc<name>]:`. Fuera del asistente, `ccam add` solo crea uno con `--shortcut`. `ccam list` muestra el atajo de cada cuenta, o `-`.

Un atajo pasa todos los argumentos a Claude Code, igual que `ccam <name>` después del nombre. Eso incluye `--dry-run` y `--help`: para previsualizar un inicio, usa `ccam --dry-run work`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--remove` | Elimina el atajo de la cuenta y su enlace | desactivada |
| `-h`, `--help` | Ayuda del comando | |
| `--dry-run` | `ccam shortcut` la respeta: comprueba el nombre e imprime `Would set the shortcut of work to ccwork.` o `Would remove ccwork, the shortcut of work.` sin guardar ni enlazar | desactivada |

`ccam shims install` no acepta argumentos y rechaza `--dry-run` con `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Reglas de los nombres, comprobadas en este orden:

- de 1 a 64 caracteres entre `a-z`, `0-9`, `.`, `-` y `_`, empezando con una letra o un dígito;
- no puede ser `ccam`, `claude` ni un comando de ccam como `run` o `usage`;
- no puede ser el atajo de otra cuenta;
- no puede ser un archivo de `~/.local/bin` que no sea ya un enlace a ccam;
- no puede ser el nombre de otro ejecutable que esté en tu `PATH`.

Indicar un comando y `--remove` a la vez se detiene con `give a command or --remove, not both`. Establecer el nombre que una cuenta ya tiene imprime `work: ccwork (unchanged)` y vuelve a crear su enlace si falta.

## Cómo funciona

Establecer un nombre ejecuta `checkShortcut`, que llama a `shim.Check` con la carpeta bin y la ruta real del binario de ccam:

<!-- motion:shortcuts-check -->

El nombre pasa por las comprobaciones en orden. Junto a cada comprobación hay un nombre que rechaza: `Work` no cumple el patrón, `run` es un comando de ccam, `ccside` pertenece a otra cuenta, `ccfile` es un archivo existente en `~/.local/bin` y `git` ocultaría `/usr/bin/git`. Solo un nombre que pasa todas las comprobaciones se guarda y se enlaza.

1. `shim.Check` compara el nombre con `^[a-z0-9][a-z0-9._-]{0,63}$`, luego con `config.Reserved` y luego con `cfg.FindShortcut`. Llama a `os.Lstat` sobre `~/.local/bin/<name>` y, para cada otra carpeta de `PATH`, llama a `os.Stat` y rechaza un archivo ejecutable que no sea un enlace a ccam.
2. ccam guarda `config.toml` con `fsutil.AtomicWrite` (archivo temporal, `fsync`, cambio de nombre).
3. `shim.Remove` elimina el enlace del atajo anterior, pero solo si es un enlace a ccam.
4. `shim.Install` crea el enlace de cada atajo configurado que falta, con `os.Symlink` hacia el binario de ccam, resuelto con `os.Executable` y `filepath.EvalSymlinks`. Nunca quita nada.

Un enlace cuenta como propio de ccam cuando `os.Readlink` devuelve la ruta de ccam, o cuando el enlace y la ruta de ccam se resuelven al mismo archivo. Cada eliminación y cada comprobación de sobrescritura pasa por esta prueba, así que ccam no toca archivos que no creó.

Cuando ejecutas un atajo, el nombre que escribiste es lo único que ccam necesita:

<!-- motion:shortcuts-link -->

La shell encuentra `ccwork` en `PATH`, el kernel sigue el enlace e inicia el binario de ccam, y `argv[0]` sigue siendo `ccwork`. En `app.Main`, cualquier nombre distinto de `ccam` carga `config.toml` una vez y llama a `shim.Lookup`, que busca la cuenta cuyo `shortcut` es igual a `filepath.Base(argv[0])`. Después, `launcher.LaunchWith` sigue los mismos pasos que [`ccam <name>`](/features/launch) con la configuración ya cargada, y pasa `argv[1:]` a `claude` sin tocarlo.

Por qué un enlace simbólico y no un alias de shell o un script envoltorio: un enlace simbólico funciona en todas las shells y en scripts, no agrega el inicio de una shell adicional (de 1 a 2 ms aproximadamente) y el inicio sigue leyendo un solo archivo. `scripts/perf-check.sh` mide un atajo con `hyperfine` contra el mismo límite de 10 ms (p95) que `ccam <name>`.

### En Windows

Un atajo es `ccwork.exe` en `%USERPROFILE%\.local\bin`, y `shim.Name` quita el `.exe` (en mayúsculas o minúsculas) antes de la búsqueda, así que `ccwork`, `ccwork.exe` y `CCWORK.EXE` inician `work`. ccam primero crea un enlace simbólico, que Windows permite en el modo de desarrollador o para un administrador. Sin ese permiso crea un enlace duro: un segundo nombre para el mismo archivo `ccam.exe`, que no necesita privilegios. Un enlace duro cuenta como propio de ccam cuando es el mismo archivo que `ccam.exe`. Como `ccam update` pone un archivo nuevo en lugar de `ccam.exe`, después renueva cada atajo con enlace duro mediante `shim.Relink` e imprime `✓ shortcut ccwork renewed`. La comprobación del nombre busca `<name>` con cada extensión de `PATHEXT` (`.exe`, `.cmd`, etc.) en las demás carpetas de `PATH`.

## Seguridad y casos límite

- **Un enlace que ninguna cuenta reclama.** Si un enlace a ccam en `~/.local/bin` tiene un nombre que ninguna cuenta usa, por ejemplo después de editar `shortcut` a mano en `config.toml`, ccam no recurre a la cuenta predeterminada. Se detiene con `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Vuelve a establecer el nombre con `ccam shortcut` o elimina el enlace.
- **Archivos existentes.** ccam nunca reemplaza un archivo que no creó. `ccam shims install` se detiene con `<path> already exists and is not a ccam link, so it was left alone`.
- **Un binario de ccam movido.** En ese caso, los enlaces antiguos apuntan a una ruta que ya no existe, y ccam no los considera propios. Elimínalos y luego ejecuta `ccam shims install`.
- **Cambiar el nombre y quitar.** `ccam rename` conserva el nombre del atajo, y el enlace sigue funcionando, porque apunta a ccam y ccam busca el nombre en `config.toml`. `ccam remove` elimina el enlace de la cuenta. `ccam shortcut <account> <new>` elimina el enlace anterior antes de crear el nuevo.
- **PATH.** Cuando ccam crea un enlace y `~/.local/bin` no está en tu `PATH`, imprime `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` seguido de la línea `export PATH=...` que debes agregar. En Windows, en su lugar imprime la línea de PowerShell que agrega la carpeta a tu `PATH` de usuario.
- **Comprobar más tarde.** `ccam doctor` informa `missing shortcuts in /home/me/.local/bin: ccwork` y avisa cuando la carpeta no está en `PATH`. `ccam doctor --fix` crea los enlaces que faltan.

## Relacionado

- [Iniciar Claude Code](/features/launch)
- [Administrar cuentas](/features/accounts)
- [Referencia de comandos](/docs/commands#shortcuts)
- [Instalar ccam](/docs/install)
