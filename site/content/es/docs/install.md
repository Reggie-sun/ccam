---
title: Instalar ccam
metaTitle: Instalar ccam en Linux, macOS y Windows
description: Instala ccam en Linux, macOS o Windows con un comando, un paquete .deb o desde el código fuente, y actualízalo con ccam update. Antes necesitas Claude Code.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 07eadc0769c8
---

Ejecuta esto en una terminal de Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

El script descarga la versión para tu sistema, comprueba su suma de verificación SHA-256 y coloca `ccam` en `~/.local/bin`, donde ccam también crea los atajos que elijas más adelante, como `ccwork`. Si esa carpeta no está en tu `PATH`, el script imprime la línea que debes agregar a `~/.bashrc` o `~/.zshrc`.

En Windows, ejecuta esto en PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Descarga `ccam_windows_<arch>.zip`, comprueba su suma de verificación SHA-256, coloca `ccam.exe` en `%USERPROFILE%\.local\bin`, la carpeta donde el instalador de Claude Code coloca `claude.exe`, y agrega esa carpeta a tu `PATH` de usuario si aún no está. Después, abre una terminal nueva.

ccam inicia Claude Code por ti, así que instala primero [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) y asegúrate de que `claude --version` funcione.

## Sistemas compatibles

- Linux en x86-64 y ARM64: Ubuntu, Debian, Fedora, Arch, Alpine y otras distribuciones. El binario es estático y no tiene dependencias.
- macOS en Apple silicon e Intel.
- Windows 10 y 11 en x86-64 y ARM64, en PowerShell, Windows Terminal o Git Bash. Las carpetas compartidas funcionan para cualquier usuario. Los archivos compartidos (`settings.json`, `CLAUDE.md`, `history.jsonl`) requieren el modo de desarrollador (Configuración > Sistema > Para programadores), porque Windows solo permite crear enlaces a archivos en ese modo o a un administrador; consulta [compartir en Windows](/features/share#on-windows).
- Dentro de WSL, ccam funciona como en cualquier sistema Linux.

## Otras formas de instalar

**Paquete para Debian y Ubuntu.** Descarga `ccam_<version>_linux_amd64.deb` desde la [página de versiones](https://github.com/jhyunwoo/ccam/releases) e instálalo:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

El paquete instala `/usr/bin/ccam`.

**Desde el código fuente.** Con Go 1.26 o posterior:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

**Una versión específica.** Define `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Autocompletado del shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

El autocompletado cubre todos los comandos y opciones de `ccam`, incluidos los nombres de las cuentas.

## Comprobar la instalación

```sh
ccam version
ccam doctor
```

`ccam doctor` comprueba que se encuentre Claude Code y que cada carpeta de cuenta, enlace e inicio de sesión esté en orden. Continúa con el [inicio rápido](/docs/quick-start).

## Actualizar

```sh
ccam update --check
ccam update
```

`ccam update --check` solo te indica si hay una versión más reciente. `ccam update` descarga la versión para tu sistema, comprueba su suma de verificación SHA-256 con el archivo `checksums.txt` de la versión y luego reemplaza el archivo `ccam` en un solo paso. Si la suma de verificación no coincide, no se cambia nada. Los atajos como `ccwork` son enlaces a ccam, así que ejecutan la nueva versión de inmediato.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

ccam nunca busca actualizaciones por su cuenta; solo busca una versión nueva cuando ejecutas este comando. Para instalar una versión concreta, incluso una anterior, usa `--version`. `ccam --dry-run update` muestra qué se descargaría y qué se reemplazaría, sin cambiar nada.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Volver a ejecutar el script de instalación también actualiza la copia de `~/.local/bin`. Si compilaste ccam con `go install`, vuelve a ejecutar el mismo comando.

**Windows.** `ccam update` descarga el `.zip` para Windows. En Windows no se puede reemplazar un programa en ejecución, así que el binario antiguo se renombra a `ccam.exe.old`, que la siguiente actualización elimina.

**Paquete para Debian y Ubuntu.** `ccam update` no reemplaza un `/usr/bin/ccam` instalado con apt; en su lugar, muestra estos pasos. Descarga el nuevo `ccam_<version>_linux_amd64.deb` desde la [página de versiones](https://github.com/jhyunwoo/ccam/releases) e instálalo sobre el anterior:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Desinstalar

Primero quita los atajos de las cuentas, mientras ccam sigue instalado. `ccam list` los muestra en la columna Shortcut; quita cada uno con `ccam shortcut <account> --remove`, o elimina tú mismo esos enlaces de `~/.local/bin`. Luego elimina el comando y, si ya no necesitas las cuentas, las carpetas de datos:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

En Windows, en PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Al eliminar `~/.local/share/ccam` también se eliminan los inicios de sesión y el historial de conversaciones guardados allí. Una carpeta `~/.claude` que hayas [importado](/docs/commands) se queda donde está.
