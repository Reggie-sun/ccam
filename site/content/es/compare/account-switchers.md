---
title: "Cambiar cuenta de Claude Code: ccam y alternativas"
description: Cómo se compara ccam con claude-swap, CCS, clauth, CAAM, cux y otras herramientas para cambiar de cuenta en Claude Code, y qué enfoque conviene en cada caso.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Las herramientas para cambiar de cuenta en Claude Code siguen uno de dos enfoques. La mayoría intercambia el inicio de sesión dentro de la única carpeta `~/.claude`, así que solo hay una cuenta activa a la vez. Otras, entre ellas ccam, le dan a cada cuenta su propia carpeta con `CLAUDE_CONFIG_DIR`, así que todas las cuentas siguen con la sesión iniciada y varias pueden ejecutarse a la vez. ccam agrega configuración y servidores MCP compartidos, y continúa una conversación en la carpeta de otra cuenta.

## ¿Intercambiar un inicio de sesión o usar carpetas separadas?

**Intercambiar un inicio de sesión** mantiene un solo conjunto de configuración, servidores MCP e historial, porque solo hay una carpeta. El cambio es rápido, y una sesión en curso puede tomar el nuevo inicio de sesión. Pero solo hay una cuenta activa a la vez, para todas las terminales y editores a la vez.

**Las carpetas separadas** son el método que [describe la documentación de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Todas las cuentas siguen con la sesión iniciada, y puedes ejecutar una sesión de trabajo y una personal en paralelo. El precio es que cada carpeta empieza con su propia configuración, sus servidores MCP y su historial, a menos que una herramienta los comparta.

## Comparación de las herramientas

Tomado de la documentación de cada proyecto el 2026-10-19. "No documentado" significa que el proyecto no lo dice, no que la función falte.

| Herramienta | Enfoque | Varias cuentas a la vez | Uso de 5 horas y semanal | Al llegar a un límite | Plataformas |
| --- | --- | --- | --- | --- | --- |
| ccam | Carpetas separadas | Sí | Sí, `ccam usage` | Ejecutas `ccam next --auto --continue` para continuar la conversación en la cuenta con más margen | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Intercambia el inicio de sesión; carpetas separadas con `cswap run` | Sí, con `cswap run` | Sí | `cswap auto` cambia el inicio de sesión al superar un umbral; la sesión en curso continúa | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Carpetas separadas | Sí | En parte, en una app de la barra de menús de macOS | Los grupos de cuentas pasan las solicitudes a la siguiente cuenta mediante un proxy local | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Intercambia el inicio de sesión; carpetas separadas con `clauth start` | Sí, con `clauth start` | Sí | Cambia siguiendo una cadena de respaldo | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Intercambia el inicio de sesión o usa perfiles aislados | Sí, con perfiles aislados | Sí | `caam run` reintenta los comandos no interactivos en otra cuenta | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Intercambia el inicio de sesión y envuelve `claude` | No documentado | Sí | Intercambia el inicio de sesión y reanuda la misma conversación | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Intercambia el inicio de sesión | No | Sí | No documentado | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Intercambia el inicio de sesión | No | Sí | Cambio opcional al 100% | macOS |

ccswap es un fork de claude-swap que también gestiona inicios de sesión de Codex. Varias de estas herramientas también manejan otros agentes de programación, como Codex o Gemini.

## En qué se centra ccam

- Todas las cuentas siguen con la sesión iniciada. Cambiar nunca cierra la sesión de una, y `ccam work` y `ccam personal` pueden ejecutarse en dos terminales a la vez.
- Los servidores MCP, `settings.json`, `CLAUDE.md`, los agentes, las skills y los plugins se comparten entre carpetas mediante enlaces, elemento por elemento y cuenta por cuenta. Los inicios de sesión nunca se comparten.
- Con el historial compartido, un solo comando, `ccam next --auto --continue`, abre la misma conversación en la carpeta de la cuenta con más margen.
- ccam inicia Claude Code y se reemplaza a sí mismo por él, así que no queda nada ejecutándose en segundo plano. No vigila tu uso ni cambia de cuenta por su cuenta; tú decides cuándo.

## Dónde otras herramientas hacen más

- claude-swap, CCS, clauth, CAAM y cux funcionan en Windows sin WSL. ccam no.
- claude-swap, clauth y cux cambian por su cuenta cuando una cuenta supera un umbral de uso. ccam deja el cambio en tus manos.
- CCSwitcher, Claude Switcher y claude-swap muestran el uso en la barra de menús de macOS.
- CCS, CAAM, ccswap y Claude Switcher también gestionan Codex y otros agentes de programación.

## ¿Cuál te conviene?

Elige una herramienta que intercambie el inicio de sesión si usas una cuenta a la vez, quieres que el cambio sea automático o estás en Windows sin WSL. Elige ccam si quieres cuentas en paralelo, una sola configuración compartida y la misma conversación en la siguiente cuenta. Si no necesitas compartir nada, bastan [simples alias del shell](/compare/aliases).

Si puedes usar varias cuentas depende de los términos de cada una. Revísalos antes de hacerlo.
