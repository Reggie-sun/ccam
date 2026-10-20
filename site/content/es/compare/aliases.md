---
title: ccam comparado con los alias del shell
description: Alias del shell con CLAUDE_CONFIG_DIR frente a ccam para usar varias cuentas de Claude Code, lado a lado.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: ea24bb091caf
---

Un alias del shell por cuenta es la forma más sencilla de usar varias cuentas de Claude Code. ccam hace lo mismo por debajo y agrega uso compartido, cambio de cuenta y seguimiento del uso.

## Configuración con alias

Agrega una línea por cada cuenta adicional a `~/.bashrc` o `~/.zshrc`, abre una terminal nueva y ejecuta el alias una vez para iniciar sesión:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` sigue siendo tu primera cuenta y `claude-work` inicia la segunda. Este es el método que [describe la documentación de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## Configuración con ccam

Instala ccam y agrega la cuenta, en este caso con su propio comando de atajo:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` y `ccwork` inician la cuenta, y ambos funcionan también en scripts.

## Lado a lado

| | Alias del shell | ccam |
| --- | --- | --- |
| Inicios de sesión separados | Sí | Sí |
| Configuración por cuenta | Editar `~/.bashrc`, iniciar sesión | `ccam add` |
| Funciona en scripts y otros programas | No, los alias solo existen en el shell | Sí, `ccam work` y los atajos como `ccwork` son comandos reales |
| Configuración, `CLAUDE.md` y extensiones compartidos | Copiar o enlazar a mano | Por elemento, con copias de seguridad |
| Servidores MCP compartidos | A mano, en cada cuenta | `ccam mcp add` una vez |
| Continuar una conversación en otra cuenta | No | `ccam next --continue` |
| Uso de todas las cuentas | No | `ccam usage` |
| Elegir la cuenta con más margen | No | `ccam next --auto` |
| Nombre de la cuenta en la línea de estado | No | `ccam statusline install` |
| Reparar enlaces rotos | A mano | `ccam doctor --fix` |
| Software adicional | Ninguno | Un único binario estático |

Elige los alias si tienes dos cuentas y no quieres compartir nada entre ellas. Elige ccam si quieres que las cuentas funcionen como una sola configuración. Puedes empezar con alias y pasarte a ccam más adelante: `ccam import` registra las carpetas `~/.claude-*` existentes donde están.

[Cómo usar varias cuentas en una computadora](/guides/multiple-claude-code-accounts) explica ambos métodos paso a paso, y [ccam y otros selectores de cuentas](/compare/account-switchers) compara ccam con las herramientas que intercambian un único inicio de sesión.
