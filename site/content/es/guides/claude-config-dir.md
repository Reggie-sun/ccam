---
title: CLAUDE_CONFIG_DIR explicado
description: Qué hace CLAUDE_CONFIG_DIR en Claude Code, qué se guarda en la carpeta, dónde guarda macOS el inicio de sesión y qué archivos se pueden compartir entre cuentas.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: ef9317030e42
---

`CLAUDE_CONFIG_DIR` le indica a Claude Code qué carpeta usar en lugar de `~/.claude`. Todo lo específico de la cuenta vive ahí: el inicio de sesión, la configuración, los servidores MCP, los plugins y el historial de conversaciones. Si lo apuntas a carpetas distintas, obtienes cuentas separadas.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## ¿CLAUDE_CONFIG_DIR tiene soporte oficial?

Sí. Claude Code la incluye entre sus [variables de entorno](https://code.claude.com/docs/en/env-vars) como la forma de cambiar la carpeta de configuración, "útil para ejecutar varias cuentas en paralelo", y su [guía de autenticación](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) la usa para iniciar sesión con varias cuentas. Puedes definirla en tu shell o en la configuración de usuario, pero no en la configuración de un proyecto.

## ¿Qué hay dentro de la carpeta?

Comprobado con Claude Code 2.1.302:

| Ruta | Contenido | Se puede compartir |
| --- | --- | --- |
| `.credentials.json` | Tokens de inicio de sesión (Linux, Windows), tokens OAuth de MCP | No |
| `.claude.json` | Datos de la cuenta, servidores MCP de usuario, estado por proyecto | No (los servidores MCP, por separado) |
| `settings.json` | Permisos, hooks, modelo, línea de estado, plugins | Sí |
| `CLAUDE.md` | Instrucciones globales | Sí |
| `agents/`, `commands/`, `skills/`, `plugins/` | Extensiones | Sí |
| `projects/` | Conversaciones como archivos `.jsonl`, memoria por proyecto | Sí |
| `history.jsonl` | Historial de prompts | Sí |
| `sessions/`, `shell-snapshots/`, `cache/` y otros | Estado en tiempo de ejecución | No |

## ¿Qué suele causar confusión?

**Sin `CLAUDE_CONFIG_DIR`, `.claude.json` no está dentro de `~/.claude`.** En ese caso, Claude Code lo guarda en `~/.claude.json`, junto a la carpeta. Si empiezas a usar `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code busca `~/.claude/.claude.json` en su lugar y no encuentra tu cuenta. `ccam import` resuelve esto enlazando los dos.

**En macOS, el inicio de sesión está en el Llavero, no en la carpeta.** Claude Code lo guarda con el nombre `Claude Code-credentials`, seguido de un hash corto de la ruta de la carpeta cuando `CLAUDE_CONFIG_DIR` está definido. Por eso, carpetas distintas obtienen entradas distintas en el Llavero. Si el Llavero rechaza la escritura, Claude Code recurre a `.credentials.json` en la carpeta, como describe [su documentación](https://code.claude.com/docs/en/authentication#credential-management).

**No todos los inicios de sesión viven en la carpeta.** Un inicio de sesión de Claude Console sin clave de API se guarda fuera de ella, así que las carpetas separadas no mantienen dos de ellos aparte. Los inicios de sesión de claude.ai y las claves de API se guardan por carpeta.

## ¿Dónde está la carpeta en Windows?

De forma predeterminada en `%USERPROFILE%\.claude`, con el inicio de sesión en `.credentials.json` dentro de ella. Define `CLAUDE_CONFIG_DIR` con otra carpeta para cada cuenta adicional, igual que en Linux y macOS.

## ¿Las carpetas pueden compartir archivos?

Los enlaces simbólicos funcionan: Claude Code lee y escribe a través de un `settings.json` o una carpeta `projects/` enlazados y conserva el enlace. Los servidores MCP son la excepción, porque comparten `.claude.json` con el inicio de sesión; pásalos con `--mcp-config` en su lugar.

ccam se basa exactamente en estas reglas. Consulta [compartir](/docs/sharing) para ver cómo enlaza las carpetas, y [la guía de varias cuentas](/guides/multiple-claude-code-accounts) para ver el panorama general.
