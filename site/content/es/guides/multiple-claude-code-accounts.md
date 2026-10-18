---
title: Cómo usar dos cuentas de Claude Code en una computadora
description: Usa varias cuentas de Claude Code en paralelo con CLAUDE_CONFIG_DIR, a mano con alias del shell o automáticamente con ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d82129ddcb4c
---

Dale a cada cuenta su propia carpeta de configuración. Claude Code lee la carpeta de la variable de entorno `CLAUDE_CONFIG_DIR`, y cada carpeta guarda su propio inicio de sesión, configuración e historial. Puedes configurar varias cuentas de Claude Code a mano con un alias del shell o dejar que ccam lo haga.

## A mano, con un alias

Este es el método que describe la documentación de Claude Code. Agrega una línea por cada cuenta adicional a `~/.bashrc` o `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Abre una terminal nueva y ejecuta `claude-work`. Claude Code te pide que inicies sesión, porque la carpeta es nueva. A partir de ahí, `claude` es tu primera cuenta y `claude-work` la segunda.

Esto funciona bien para dos cuentas. Lo que debes saber:

- Cada carpeta empieza vacía: sin configuración, sin `CLAUDE.md`, sin servidores MCP, sin agentes ni skills. Copias o enlazas lo que necesitas y mantienes las copias sincronizadas tú mismo.
- Los alias solo existen en shells interactivos. Los scripts, los Makefiles y otros programas no los ven.
- Los servidores MCP viven en `.claude.json`, junto al inicio de sesión, así que no se pueden simplemente enlazar entre carpetas.
- Las conversaciones se quedan en la cuenta donde empezaron, así que no puedes continuar una en otra cuenta.
- No hay nada que muestre el uso de todas las cuentas a la vez.

## Con ccam

ccam automatiza la misma idea y agrega las partes que faltan:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Ahora `claude1` y `claude2`, o `ccs personal` y `ccs work`, inician las dos cuentas. Como son programas reales y no alias, también funcionan en scripts. La configuración, `CLAUDE.md`, los servidores MCP y las extensiones se comparten desde el principio, y puedes activar el historial de conversaciones compartido para continuar una conversación en la otra cuenta con `ccs next --continue`.

Si ya usas `~/.claude`, ejecuta `ccam import` en lugar de volver a agregar esa cuenta.

## Cuál elegir

Para dos cuentas sin nada compartido, basta con un alias. Si quieres configuración y servidores MCP compartidos, más cuentas, ver el uso de un vistazo o continuar conversaciones, ccam te ahorra el trabajo de mantenimiento. La [comparación](/compare/aliases) muestra las diferencias lado a lado.
