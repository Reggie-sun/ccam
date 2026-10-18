---
title: ccam comparado con los alias del shell
description: Alias del shell con CLAUDE_CONFIG_DIR frente a ccam para usar varias cuentas de Claude Code, lado a lado.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: e9902676e0ed
---

Un alias del shell por cuenta es la forma más sencilla de usar varias cuentas de Claude Code. ccam hace lo mismo por debajo y agrega uso compartido, cambio de cuenta y seguimiento del uso.

| | Alias del shell | ccam |
| --- | --- | --- |
| Inicios de sesión separados | Sí | Sí |
| Configuración por cuenta | Editar `~/.bashrc`, iniciar sesión | `ccam add` |
| Funciona en scripts y otros programas | No, los alias solo existen en el shell | Sí, `claude2` es un comando real |
| Configuración, `CLAUDE.md` y extensiones compartidos | Copiar o enlazar a mano | Por elemento, con copias de seguridad |
| Servidores MCP compartidos | A mano, en cada cuenta | `ccam mcp add` una vez |
| Continuar una conversación en otra cuenta | No | `ccs next --continue` |
| Uso de todas las cuentas | No | `ccs usage` |
| Elegir la cuenta con más margen | No | `ccs next --auto` |
| Nombre de la cuenta en la línea de estado | No | `ccam statusline install` |
| Reparar enlaces rotos | A mano | `ccam doctor --fix` |
| Software adicional | Ninguno | Un único binario estático |

Elige los alias si tienes dos cuentas y no quieres compartir nada entre ellas. Elige ccam si quieres que las cuentas funcionen como una sola configuración. Puedes empezar con alias y pasarte a ccam más adelante: `ccam import` registra las carpetas `~/.claude-*` existentes donde están.

[Cómo usar dos cuentas en una computadora](/guides/multiple-claude-code-accounts) explica ambos métodos paso a paso.
