---
title: Usar varias cuentas de Claude Code en una computadora
description: Ejecuta varias cuentas de Claude Code en paralelo con CLAUDE_CONFIG_DIR, a mano con alias del shell o automáticamente con ccam, y cambia entre ellas.
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "5c5dbd9fcd23"
---

Dale a cada cuenta su propia carpeta de configuración. Claude Code lee la carpeta de la variable de entorno `CLAUDE_CONFIG_DIR`, y cada carpeta guarda su propio inicio de sesión, configuración e historial, así que todas las cuentas siguen con la sesión iniciada y puedes ejecutarlas a la vez. Puedes configurarlo a mano con un alias del shell o dejar que ccam lo haga.

## ¿Cómo lo configuro a mano?

Este es el método que [describe la documentación de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Agrega una línea por cada cuenta adicional a `~/.bashrc` o `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Abre una terminal nueva y ejecuta `claude-work`. Claude Code te pide que inicies sesión, porque la carpeta es nueva. A partir de ahí, `claude` es tu primera cuenta y `claude-work` la segunda. Agrega un alias por cuenta, cada uno con su propia carpeta.

En Windows sin WSL, una función de PowerShell en tu perfil (`notepad $PROFILE`) hace lo mismo:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## ¿Qué le falta al método de los alias?

Funciona bien para dos cuentas. Lo que debes saber:

- Cada carpeta empieza vacía: sin configuración, sin `CLAUDE.md`, sin servidores MCP, sin agentes ni skills. Copias o enlazas lo que necesitas y mantienes las copias sincronizadas tú mismo.
- Los alias solo existen en shells interactivos. Los scripts, los Makefiles y otros programas no los ven.
- Los servidores MCP viven en `.claude.json`, junto al inicio de sesión, así que no se pueden simplemente enlazar entre carpetas.
- Las conversaciones se quedan en la cuenta donde empezaron, así que no puedes continuar una en otra cuenta.
- No hay nada que muestre el uso de todas las cuentas a la vez.

## ¿Puedo cambiar de cuenta con /login?

Sí, si solo cambias de vez en cuando. Escribir `/login` dentro de Claude Code [cambia la cuenta](https://code.claude.com/docs/en/quickstart) de la carpeta en la que estás. Solo hay una cuenta con la sesión iniciada a la vez, así que vuelves a iniciar sesión en cada cambio y no puedes ejecutar dos cuentas en paralelo. Las carpetas separadas evitan ambas cosas.

## ¿Cómo lo configuro con ccam?

ccam automatiza la misma idea y agrega las partes que faltan:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Ahora `ccam personal` y `ccam work` inician las dos cuentas, y también `ccwork`, el atajo que le diste a la cuenta work. Como son programas reales y no alias, también funcionan en scripts. Ejecuta `ccam` solo para elegir una cuenta de una lista, y `ccam use work` para que una sea la predeterminada.

La configuración, `CLAUDE.md`, los servidores MCP y las extensiones se comparten desde el principio. `ccam next --continue` lleva una conversación a la otra cuenta, y el historial de conversaciones compartido, que puedes activar, la mantiene como una sola conversación. `ccam usage` muestra el uso de 5 horas y semanal de todas las cuentas a la vez.

Si ya usas `~/.claude`, ejecuta `ccam import` en lugar de volver a agregar esa cuenta, o `ccam migrate` para copiar su configuración, sus servidores MCP y sus extensiones a la configuración compartida.

## ¿Cómo sé qué cuenta se está ejecutando?

Dentro de Claude Code, `/status` muestra la cuenta con la que iniciaste sesión. Con ccam, ejecuta `ccam statusline install` una vez y la línea de estado mostrará siempre el nombre de la cuenta y su uso, por ejemplo `[work] · Opus · 5h 24% · 7d 41%`.

## ¿Funciona en Windows y WSL?

Claude Code se ejecuta [de forma nativa en Windows y dentro de WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` funciona en ambos, así que en Windows nativo basta con el alias o la función de PowerShell de arriba. ccam también funciona de forma nativa en Windows: instálalo en PowerShell con `irm https://ccam.dev/install.ps1 | iex`, o dentro de WSL con el comando de Linux. Los atajos pasan a ser `ccwork.exe`, y compartir archivos como `settings.json` requiere el modo de desarrollador ([Instalar ccam](/docs/install)).

## Solución de problemas

**Claude Code me pide iniciar sesión cada vez.** El primer inicio en una carpeta nueva siempre lo pide. Si lo sigue pidiendo, el inicio de sesión caducó o se eliminó: ejecuta `ccam login work` (o inicia el alias una vez e inicia sesión).

**Un servidor MCP aparece solo en una cuenta.** Los servidores agregados dentro de una sesión con `claude mcp add` pertenecen a esa cuenta y, sin `-s user`, solo a esa carpeta. Pasa uno a todas las cuentas con `ccam mcp import work <server>`.

**La configuración es distinta entre cuentas.** `ccam share work` muestra lo que comparte la cuenta y a través de qué grupo; `ccam group join default work` hace que comparta lo que comparte `default`, después de mostrar lo que va a cambiar. `ccam doctor --fix` repara los enlaces rotos y guarda una copia de seguridad de cada archivo que modifica.

**Dos cuentas abren la misma conversación.** Evita trabajar en una conversación desde dos cuentas al mismo tiempo. ccam te avisa si lo intentas.

## ¿Cuál debo elegir?

Para dos cuentas sin nada compartido, basta con un alias. Si quieres configuración y servidores MCP compartidos, más cuentas, ver el uso de un vistazo o continuar conversaciones, ccam te ahorra el trabajo de mantenimiento. La [comparación con los alias](/compare/aliases) muestra las diferencias lado a lado, y [ccam y otros selectores de cuentas](/compare/account-switchers) trata las herramientas que intercambian un único inicio de sesión.

Si puedes usar varias cuentas depende de los términos de cada una. Revísalos antes de hacerlo.

**Asientos de Team y Enterprise.** La elección automática (`ccam next`, `--fallback`, `ccam pick`, la tecla `n` del panel y [`auto_switch`](/features/next#auto-switch) cuando lo activas) también pasa de un asiento a otro dentro de una misma organización de Claude Team o Enterprise. El centro de ayuda de Anthropic indica: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada asiento de un plan Team se asigna a una persona específica, y solo esa persona puede usarlo) y "Users can't route their work through another user's account to use that user's remaining usage" (los usuarios no pueden hacer pasar su trabajo por la cuenta de otro usuario para consumir el uso restante de ese usuario) ([planes Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo mismo para los [planes Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por lo tanto, usar la elección automática entre asientos puede infringir los términos de Anthropic. Si la usas ahí es decisión tuya: consulta con tu organización y con Anthropic antes de depender de ella. Los planes personales (Pro, Max) se tratan en la nota anterior.
