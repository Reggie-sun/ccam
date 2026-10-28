---
title: Instalar ccam
metaTitle: Instalar ccam en Linux, macOS y Windows
description: Instala ccam en Linux, macOS o Windows con un comando, un paquete .deb o go install, y actualízalo con ccam update. Antes necesitas Claude Code.
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

ccam inicia Claude Code por ti, así que instala primero [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) y asegúrate de que `claude --version` funcione.

Ejecuta esto en una terminal de Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

El script descarga la versión para tu sistema, comprueba su suma de verificación SHA-256 y coloca `ccam` en `~/.local/bin`, donde ccam también crea los atajos que elijas más adelante, como `ccwork`. Si esa carpeta no está en tu `PATH`, el script imprime la línea que debes agregar a `~/.bashrc` o `~/.zshrc`.

En Windows, ejecuta esto en PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Descarga `ccam_windows_<arch>.zip`, comprueba su suma de verificación SHA-256 y coloca `ccam.exe` en `%USERPROFILE%\.local\bin`, la carpeta donde el instalador de Claude Code coloca `claude.exe`. Si esa carpeta aún no está en tu `PATH` de usuario, el script la agrega y deja como están entradas como `%USERPROFILE%\bin`. Después, abre una terminal nueva.

## Sistemas compatibles

- Linux en x86-64 y ARM64: Ubuntu, Debian, Fedora, Arch, Alpine y otras distribuciones. El binario es estático y no tiene dependencias.
- macOS en Apple silicon e Intel.
- Windows 10 y 11 en x86-64 y ARM64, en PowerShell, Windows Terminal o Git Bash:
  - Las carpetas compartidas funcionan sin el modo de desarrollador.
  - Los archivos compartidos (`settings.json`, `CLAUDE.md`, `history.jsonl`) requieren el modo de desarrollador (Configuración > Sistema > Para programadores), porque Windows solo permite crear enlaces a archivos en ese modo o a un administrador; consulta [compartir en Windows](/features/share#on-windows).
  - Sin el modo de desarrollador, el primer `ccam add` crea el grupo `default` solo con servidores MCP y extensiones, e indica cómo agregar el resto una vez activado el modo de desarrollador.
- Dentro de WSL, ccam funciona igual que en cualquier otro sistema Linux.

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

Si instalas una etiqueta de versión (`@latest` o `@vX.Y.Z`), ccam informa esa versión, y `ccam update` y la oferta de actualización del panel funcionan. Una compilación de `@main`, o de una copia local con commits posteriores a una etiqueta, informa `dev`, y `ccam update` la rechaza a menos que pases `--version`.

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

## Plugin de Claude Code

El plugin de ccam le enseña al propio Claude Code a usar ccam: a comprobar qué cuenta tiene margen, a iniciar workers `claude -p` en paralelo en cuentas distintas con `ccam pick` y a explicar con `ccam ext` por qué falta un skill o un servidor MCP. Instálalo dentro de Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Añade:

- un skill que sugiere comandos que inician sesión o mueven archivos, pero nunca los ejecuta;
- los comandos `/ccam:usage`, `/ccam:ext` y `/ccam:switch`;
- dos hooks ([detalles](/features/next#inside-claude-code)). Uno muestra una notificación en la terminal cuando una sesión se detiene en un límite de uso; en Linux y macOS, con [`auto_switch`](/features/next#auto-switch) activado, ccam mueve entonces la sesión a otra cuenta con margen. El otro permite que `/ccam:switch [account]` pase la conversación a otra cuenta.

Los hooks ejecutan `ccam`, así que debe estar en el `PATH` que ve Claude Code.

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

El panel (`ccam` sin argumentos, en una terminal) también busca una versión nueva cada vez que se abre y, si la hay, la ofrece: pulsa `U` y confirma. Ningún otro comando busca actualizaciones, y ccam nunca se actualiza sin preguntar. Para instalar una versión concreta, incluso una anterior, usa `--version`. `ccam --dry-run update` muestra qué se descargaría y qué se reemplazaría, sin cambiar nada.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Volver a ejecutar el script de instalación también actualiza la copia de `~/.local/bin`. Si compilaste ccam con `go install` desde una etiqueta de versión, puedes ejecutar `ccam update` o volver a ejecutar el mismo comando `go install`.

**Windows.** `ccam update` descarga el `.zip` para Windows. En Windows no se puede reemplazar un programa en ejecución, así que el binario antiguo se renombra a `ccam.exe.old`, que la siguiente actualización elimina; mientras una sesión sigue ejecutándose desde un `ccam.exe.old` anterior, el nombre es `ccam.exe.old-<time>`. Los atajos que son enlaces duros se renuevan justo después. Volver a ejecutar `install.ps1` funciona igual, incluso mientras una sesión de Claude Code sigue ejecutándose desde el binario antiguo, e imprime `OK 2 shortcuts now run the new ccam` cuando renueva los atajos que son enlaces duros.

**Paquete para Debian y Ubuntu.** `ccam update` no reemplaza un `/usr/bin/ccam` instalado con apt; en su lugar, muestra estos pasos. Descarga el nuevo `ccam_<version>_linux_amd64.deb` desde la [página de versiones](https://github.com/jhyunwoo/ccam/releases) e instálalo sobre el anterior:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Desinstalar

`ccam uninstall` quita ccam y todo lo que creó. Primero lista lo que se va y pregunta:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- Cada cuenta se quita como lo hace `ccam remove`. A una cuenta que está en las carpetas de ccam primero se le cierra la sesión y luego se elimina con su carpeta, incluidos su historial de conversaciones y su configuración.
- Una carpeta `~/.claude` que hayas [importado](/docs/commands) se queda donde está para el `claude` normal. Si está en un grupo, recupera una copia de todo lo que compartía. La línea de estado de ccam en su `settings.json` cede el sitio a la línea de estado que ccam envolvía, si la había. Si tiene el plugin de ccam, ccam imprime el comando `/plugin uninstall` para que lo ejecutes en Claude Code.
- Se eliminan los atajos, `~/.config/ccam`, `~/.local/share/ccam` (con la papelera de ccam) y `~/.local/state/ccam`, y después el propio comando `ccam`.

ccam se niega mientras una cuenta está en ejecución; cierra antes sus sesiones de Claude Code. Un `/usr/bin/ccam` que instaló apt se queda; quítalo después con `sudo apt remove ccam`. En Windows un programa en ejecución no puede eliminarse a sí mismo, así que `ccam uninstall` termina imprimiendo el comando `Remove-Item` para `ccam.exe`.
