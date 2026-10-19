---
title: Presentamos ccam
description: Por qué existe ccam, qué hace por quienes tienen varias cuentas de Claude Code y cómo evita estorbar.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 985e2d369c34
---

ccam es una pequeña herramienta de línea de comandos para quienes usan más de una cuenta de Claude Code en la misma computadora: una cuenta personal y una de trabajo, o varias cuentas para un proyecto largo.

## El problema

Claude Code admite varias cuentas mediante `CLAUDE_CONFIG_DIR`: una carpeta por cuenta. En la práctica, cada carpeta es una isla. La configuración, `CLAUDE.md`, los servidores MCP, los agentes y las skills se tienen que configurar de nuevo en cada carpeta y mantener sincronizados a mano. Las conversaciones se quedan donde empezaron. Y cuando una cuenta llega a su límite de uso, te enteras al intentar usarla.

## Qué hace ccam

- `ccam add` crea la carpeta, ejecuta el inicio de sesión y le da a la cuenta un atajo: `claude1`, `claude2`, etc.
- `ccam work` y `claude2` inician Claude Code con esa cuenta y pasan todos los argumentos sin cambios.
- Los servidores MCP, la configuración, `CLAUDE.md` y las extensiones se comparten mediante enlaces, elemento por elemento y cuenta por cuenta. Los inicios de sesión nunca se comparten.
- `ccam usage` muestra el uso de 5 horas y semanal de todas las cuentas; `ccam next --continue` retoma la conversación en la siguiente cuenta.

## Sin estorbar

Una herramienta que ejecutas antes de cada sesión no debería tener un costo que notes. ccam lee un archivo de configuración, comprueba los enlaces compartidos con una sola llamada al sistema cada uno y luego se reemplaza a sí mismo por `claude`. Agrega unos 4 milisegundos y no se está ejecutando mientras trabajas.

También es cuidadoso con tus archivos. Cada cambio en el uso compartido imprime un plan antes de empezar y guarda los originales en una carpeta de copia de seguridad. `ccam doctor` revisa enlaces, permisos e inicios de sesión, y repara lo que puede reparar de forma segura.

## Pruébalo

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam es de código abierto y no oficial: no está afiliado con Anthropic ni cuenta con su respaldo. Si puedes usar varias cuentas depende de los términos de cada cuenta. Los comentarios y las traducciones son bienvenidos en [GitHub](https://github.com/jhyunwoo/ccam).
