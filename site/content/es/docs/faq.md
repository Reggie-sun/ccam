---
title: Preguntas frecuentes
metaTitle: "Preguntas frecuentes: varias cuentas de Claude Code"
description: Respuestas sobre ccam, cuentas de Claude Code, límites de uso, uso compartido, privacidad y sistemas compatibles.
updated: 2026-10-24
order: 6
translation: machine
sourceHash: c4e3b89bb4ea
---

## ¿ccam es una herramienta oficial de Anthropic?

No. ccam es un proyecto no oficial de código abierto y no está afiliado con Anthropic ni cuenta con su respaldo. Usa la opción `CLAUDE_CONFIG_DIR`, que [Claude Code documenta](https://code.claude.com/docs/en/env-vars) para ejecutar varias cuentas en paralelo.

## ¿Puedo usar varias cuentas de Claude Code?

Depende de los términos de cada cuenta y, en el caso de las cuentas de trabajo, de la política de tu organización. Revísalos antes de usar varias cuentas. ccam solo facilita cambiar entre cuentas que ya tienes.

**Asientos de Team y Enterprise.** La elección automática (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por lo tanto, usar la elección automática entre asientos puede infringir los términos de Anthropic. Consulta con tu organización y con Anthropic antes de depender de ella. Los planes personales (Pro, Max) se tratan en el párrafo anterior.

## ¿ccam evita los límites de uso?

No. Cada cuenta mantiene sus propios límites. ccam los muestra y te permite continuar en otra cuenta tuya cuando una se agota.

## ¿Cómo uso dos cuentas de Claude Code en una computadora?

Dale a cada cuenta su propia carpeta de configuración con `CLAUDE_CONFIG_DIR`. ccam lo hace por ti: ejecuta `ccam add` para cada cuenta y luego inícialas con `ccam personal` y `ccam work`, o con atajos a los que tú les pones nombre, como `ccwork`. La [guía](/guides/multiple-claude-code-accounts) también muestra cómo hacerlo a mano.

## ¿Cómo continúo una conversación en otra cuenta?

Ejecuta `ccam next --continue` en la carpeta: lleva la última conversación de la carpeta a la siguiente cuenta (`--auto` para la que tiene más margen). Dentro de Claude Code, `/ccam:switch` del [plugin de ccam](/docs/install#claude-code-plugin) hace lo mismo con la sesión en la que estás. Cuando las dos cuentas comparten el historial de conversaciones (`ccam group set default --add history`, o un grupo propio: `ccam group add pair --items history --members work,side`), la nueva cuenta abre esa misma conversación; si no, Claude Code la copia en el historial de la nueva cuenta.

## ¿ccam ve o envía mi inicio de sesión?

ccam nunca comparte ni copia tu inicio de sesión de Claude entre cuentas, y nunca copia los inicios de sesión de los servidores MCP. Lee el token de acceso de una cuenta solo para consultar al endpoint de uso de Anthropic el uso de esa cuenta, cuando ejecutas `ccam usage`, `ccam next --auto` o `ccam doctor`, y solo renueva uno caducado con `ccam usage --refresh`; con `usage_source = "claude"` no lee el token en absoluto. `ccam ext` y `ccam doctor` leen los nombres de los inicios de sesión MCP de una cuenta, nunca sus tokens. ccam no tiene telemetría. Aparte de esas consultas, la única solicitud que hace por su cuenta es la pregunta del panel a GitHub por la última versión de ccam, que no lleva nada sobre ti ni sobre tus cuentas.

## ¿Cómo inicio sesión en un servidor remoto?

Ejecuta `ccam add` o `ccam login <name>` en la sesión SSH como de costumbre. Claude Code espera el inicio de sesión en el servidor, así que el navegador de tu computadora termina en una página `http://localhost:<port>/callback?code=…` que no carga. Copia esa dirección del navegador y pégala en la terminal: ccam se la pasa a Claude Code en el servidor, que completa el inicio de sesión. La dirección contiene un código de un solo uso, así que no la compartas.

## ¿Dónde se guardan mis cuentas?

En `~/.local/share/ccam/profiles/<name>/`, una carpeta por cuenta. En macOS, Claude Code guarda el inicio de sesión en sí en el Llavero, con un nombre derivado de esa carpeta, así que `ccam rename` conserva allí el nombre anterior de la carpeta.

## ¿Qué pasa con mi ~/.claude actual?

Nada, a menos que ejecutes `ccam import` o `ccam migrate`. La importación registra la carpeta donde está, así que `claude` a secas y `ccam <name>` siguen usando la misma cuenta y la misma configuración. `ccam migrate` copia su configuración, `CLAUDE.md`, sus extensiones y sus servidores MCP a la configuración compartida, y solo lee la carpeta.

## ¿Qué sistemas admite ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine y otros), macOS y Windows 10 y 11, de forma nativa o dentro de WSL. En Windows, compartir archivos como `settings.json` requiere el modo de desarrollador; las carpetas compartidas funcionan sin él, y sin él el primer `ccam add` crea el grupo `default` solo con servidores MCP y extensiones ([compartir en Windows](/features/share#on-windows)).

## ¿Por qué un servidor MCP aparece solo en una cuenta?

Los servidores agregados dentro de una sesión con `claude mcp add` se guardan solo para esa cuenta y, sin `-s user`, solo para la carpeta actual. ccam te avisa la próxima vez que inicies la cuenta, y `ccam mcp list` y `ccam ext` los muestran; mueve el servidor al grupo compartido con `ccam mcp import <account> <server>`.

## ¿ccam hace más lento a Claude Code?

Iniciar Claude Code a través de ccam agrega unos 4 milisegundos, medidos con `scripts/perf-check.sh` en una laptop. ccam lee un archivo pequeño, comprueba los enlaces compartidos y luego se reemplaza a sí mismo por `claude`, así que no se está ejecutando mientras trabajas. Windows no puede reemplazar un proceso, así que allí ccam inicia Claude Code y lo espera, sin hacer nada más.

## ¿Cómo actualizo ccam?

Ejecuta `ccam update`. Descarga la versión más reciente, comprueba su suma de verificación y reemplaza el archivo `ccam`; `ccam update --check` solo te indica si hay una versión más reciente. El panel también te avisa cuando hay una versión más reciente, y `U` la instala después de preguntar; ccam nunca se actualiza por su cuenta. Si instalaste el paquete `.deb`, instala el nuevo `.deb` con apt. Consulta [actualizar](/docs/install#update).

## ¿Cómo desinstalo ccam?

Quita el atajo de cada cuenta con `ccam shortcut <account> --remove`. Si importaste `~/.claude` y está en un grupo, ejecuta primero `ccam remove <name>` para ella, así recupera una copia de lo que comparte. Luego elimina `ccam` de `~/.local/bin` y quita `~/.config/ccam`, `~/.local/share/ccam` y `~/.local/state/ccam` si ya no necesitas las cuentas. Consulta [instalar](/docs/install#uninstall).
