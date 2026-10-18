---
title: Preguntas frecuentes
description: Respuestas sobre ccam, cuentas de Claude Code, límites de uso, uso compartido, privacidad y sistemas compatibles.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 12ad0faef2ce
---

## ¿ccam es una herramienta oficial de Anthropic?

No. ccam es un proyecto no oficial de código abierto y no está afiliado con Anthropic ni cuenta con su respaldo. Usa la opción `CLAUDE_CONFIG_DIR`, que Claude Code documenta para ejecutar varias cuentas en paralelo.

## ¿Puedo usar varias cuentas de Claude Code?

Depende de los términos de cada cuenta y, en el caso de las cuentas de trabajo, de la política de tu organización. Revísalos antes de usar varias cuentas. ccam solo facilita cambiar entre cuentas que ya tienes.

## ¿ccam evita los límites de uso?

No. Cada cuenta mantiene sus propios límites. ccam los muestra y te permite continuar en otra cuenta tuya cuando una se agota.

## ¿Cómo uso dos cuentas de Claude Code en una computadora?

Dale a cada cuenta su propia carpeta de configuración con `CLAUDE_CONFIG_DIR`. ccam lo hace por ti: ejecuta `ccam add` para cada cuenta y luego inícialas con `claude1` y `claude2`. La [guía](/guides/multiple-claude-code-accounts) también muestra cómo hacerlo a mano.

## ¿Cómo continúo una conversación en otra cuenta?

Comparte el historial de conversaciones entre las dos cuentas con `ccam share <name> --history on` y luego ejecuta `ccs next --continue`. Claude Code abre la conversación más reciente de la carpeta actual en la siguiente cuenta.

## ¿ccam ve o envía mi inicio de sesión?

ccam nunca copia ni comparte un inicio de sesión. Lee el token de acceso de una cuenta solo para consultar al endpoint de uso de Anthropic el uso de esa cuenta, cuando ejecutas `ccs usage`, `ccs next --auto` o `ccam doctor`. ccam no tiene telemetría y no envía nada a ningún otro lugar.

## ¿Dónde se guardan mis cuentas?

En `~/.local/share/ccam/profiles/<name>/`, una carpeta por cuenta. En macOS, Claude Code guarda el inicio de sesión en sí en el Llavero, con un nombre derivado de esa carpeta.

## ¿Qué pasa con mi ~/.claude actual?

Nada, a menos que ejecutes `ccam import`. La importación registra la carpeta donde está, así que `claude` a secas y `ccs <name>` siguen usando la misma cuenta y la misma configuración.

## ¿Qué sistemas admite ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine y otros) y macOS. En Windows, úsalo dentro de WSL.

## ¿Por qué un servidor MCP aparece solo en una cuenta?

Los servidores agregados dentro de una sesión con `claude mcp add` se guardan solo para esa cuenta. ccam te avisa la próxima vez que inicies la cuenta; mueve el servidor al grupo compartido con `ccam mcp import <account> <server>`.

## ¿ccam hace más lento a Claude Code?

Iniciar Claude Code a través de ccam agrega unos 4 milisegundos, medidos con `scripts/perf-check.sh` en una laptop. ccam lee un archivo pequeño, comprueba los enlaces compartidos y luego se reemplaza a sí mismo por `claude`, así que no se está ejecutando mientras trabajas.

## ¿Cómo desinstalo ccam?

Elimina `ccam`, `ccs` y los enlaces `claude1`… de `~/.local/bin` y, luego, `~/.config/ccam`, `~/.local/share/ccam` y `~/.local/state/ccam` si ya no necesitas las cuentas. Consulta [instalar](/docs/install#uninstall).
