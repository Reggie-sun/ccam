---
title: Preguntas frecuentes
metaTitle: "Preguntas frecuentes: varias cuentas de Claude Code"
description: Respuestas sobre ccam, cuentas de Claude Code, límites de uso, uso compartido, privacidad y sistemas compatibles.
updated: 2026-10-20
order: 6
translation: machine
sourceHash: e7d4a2c877ac
---

## ¿ccam es una herramienta oficial de Anthropic?

No. ccam es un proyecto no oficial de código abierto y no está afiliado con Anthropic ni cuenta con su respaldo. Usa la opción `CLAUDE_CONFIG_DIR`, que [Claude Code documenta](https://code.claude.com/docs/en/env-vars) para ejecutar varias cuentas en paralelo.

## ¿Puedo usar varias cuentas de Claude Code?

Depende de los términos de cada cuenta y, en el caso de las cuentas de trabajo, de la política de tu organización. Revísalos antes de usar varias cuentas. ccam solo facilita cambiar entre cuentas que ya tienes.

## ¿ccam evita los límites de uso?

No. Cada cuenta mantiene sus propios límites. ccam los muestra y te permite continuar en otra cuenta tuya cuando una se agota.

## ¿Cómo uso dos cuentas de Claude Code en una computadora?

Dale a cada cuenta su propia carpeta de configuración con `CLAUDE_CONFIG_DIR`. ccam lo hace por ti: ejecuta `ccam add` para cada cuenta y luego inícialas con `ccam personal` y `ccam work`, o con atajos a los que tú les pones nombre, como `ccwork`. La [guía](/guides/multiple-claude-code-accounts) también muestra cómo hacerlo a mano.

## ¿Cómo continúo una conversación en otra cuenta?

Comparte el historial de conversaciones entre las cuentas con `ccam group set default --add history` (o con un grupo propio: `ccam group add pair --items history --members work,side`) y luego ejecuta `ccam next --continue`. Claude Code abre la conversación más reciente de la carpeta actual en la siguiente cuenta.

## ¿ccam ve o envía mi inicio de sesión?

ccam nunca comparte tu inicio de sesión de Claude entre cuentas. Lee el token de acceso de una cuenta solo para consultar al endpoint de uso de Anthropic el uso de esa cuenta, cuando ejecutas `ccam usage`, `ccam next --auto` o `ccam doctor`. `ccam migrate` copia los inicios de sesión guardados de los servidores MCP de `~/.claude` a los archivos de inicio de sesión propios de tus cuentas en la misma computadora. ccam no tiene telemetría y no envía nada a ningún otro lugar.

## ¿Cómo inicio sesión en un servidor remoto?

Ejecuta `ccam add` o `ccam login <name>` en la sesión SSH como de costumbre. Claude Code espera el inicio de sesión en el servidor, así que el navegador de tu computadora termina en una página `http://localhost:<port>/callback?code=…` que no carga. Copia esa dirección del navegador y pégala en la terminal: ccam se la pasa a Claude Code en el servidor, que completa el inicio de sesión. La dirección contiene un código de un solo uso, así que no la compartas.

## ¿Dónde se guardan mis cuentas?

En `~/.local/share/ccam/profiles/<name>/`, una carpeta por cuenta. En macOS, Claude Code guarda el inicio de sesión en sí en el Llavero, con un nombre derivado de esa carpeta.

## ¿Qué pasa con mi ~/.claude actual?

Nada, a menos que ejecutes `ccam import` o `ccam migrate`. La importación registra la carpeta donde está, así que `claude` a secas y `ccam <name>` siguen usando la misma cuenta y la misma configuración. `ccam migrate` copia su configuración, `CLAUDE.md`, sus extensiones y sus servidores MCP a la configuración compartida, y solo lee la carpeta.

## ¿Qué sistemas admite ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine y otros), macOS y Windows 10 y 11, de forma nativa o dentro de WSL. En Windows, compartir archivos como `settings.json` requiere el modo de desarrollador; las carpetas compartidas funcionan sin él ([compartir en Windows](/features/share#on-windows)).

## ¿Por qué un servidor MCP aparece solo en una cuenta?

Los servidores agregados dentro de una sesión con `claude mcp add` se guardan solo para esa cuenta. ccam te avisa la próxima vez que inicies la cuenta; mueve el servidor al grupo compartido con `ccam mcp import <account> <server>`.

## ¿ccam hace más lento a Claude Code?

Iniciar Claude Code a través de ccam agrega unos 4 milisegundos, medidos con `scripts/perf-check.sh` en una laptop. ccam lee un archivo pequeño, comprueba los enlaces compartidos y luego se reemplaza a sí mismo por `claude`, así que no se está ejecutando mientras trabajas. Windows no puede reemplazar un proceso, así que allí ccam inicia Claude Code y lo espera, sin hacer nada más.

## ¿Cómo actualizo ccam?

Ejecuta `ccam update`. Descarga la versión más reciente, comprueba su suma de verificación y reemplaza el archivo `ccam`; `ccam update --check` solo te indica si hay una versión más reciente. ccam nunca se actualiza por su cuenta. Si instalaste el paquete `.deb`, instala el nuevo `.deb` con apt. Consulta [actualizar](/docs/install#update).

## ¿Cómo desinstalo ccam?

Quita el atajo de cada cuenta con `ccam shortcut <account> --remove`, elimina `ccam` de `~/.local/bin` y, luego, `~/.config/ccam`, `~/.local/share/ccam` y `~/.local/state/ccam` si ya no necesitas las cuentas. Consulta [instalar](/docs/install#uninstall).
