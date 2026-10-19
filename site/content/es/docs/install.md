---
title: Instalar ccam
description: Instala ccam en Linux o macOS con un comando, desde un paquete .deb o desde el código fuente. Primero debes instalar Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Ejecuta esto en una terminal de Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

El script descarga la versión para tu sistema, comprueba su suma de verificación SHA-256 y coloca `ccam` en `~/.local/bin`, donde más adelante `ccam add` crea los atajos `claude1`, `claude2`, etc. Si esa carpeta no está en tu `PATH`, el script imprime la línea que debes agregar a `~/.bashrc` o `~/.zshrc`.

ccam inicia Claude Code por ti, así que instala primero [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) y asegúrate de que `claude --version` funcione.

## Sistemas compatibles

- Linux en x86-64 y ARM64: Ubuntu, Debian, Fedora, Arch, Alpine y otras distribuciones. El binario es estático y no tiene dependencias.
- macOS en Apple silicon e Intel.
- Windows no es compatible directamente. Dentro de WSL, ccam funciona como en cualquier sistema Linux.

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

## Autocompletado del shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

El autocompletado cubre todos los comandos y opciones de `ccam`, incluidos los nombres de las cuentas.

## Comprobar la instalación

```sh
ccam version
ccam doctor
```

`ccam doctor` comprueba que se encuentre Claude Code y que cada carpeta de cuenta, enlace e inicio de sesión esté en orden. Continúa con el [inicio rápido](/docs/quick-start).

## Desinstalar

Elimina el comando y los atajos de las cuentas y, luego, las carpetas de datos si ya no necesitas las cuentas:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Al eliminar `~/.local/share/ccam` también se eliminan los inicios de sesión y el historial de conversaciones guardados allí. Una carpeta `~/.claude` que hayas [importado](/docs/commands) se queda donde está.
