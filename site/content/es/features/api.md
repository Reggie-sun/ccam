---
title: Cuentas con clave de API
description: "ccam add --api agrega una cuenta con clave de API (API de Anthropic o un gateway); la clave se lee al iniciar de una variable o un comando y nunca se guarda."
updated: 2026-10-20
order: 13
translation: machine
sourceHash: "a903c92f2c50"
---

Una cuenta con clave de API inicia Claude Code con una clave de API en lugar de un inicio de sesión de una suscripción de Claude: la API de Anthropic, o un endpoint que hable el mismo formato Messages, como un gateway de LiteLLM u OpenRouter. ccam no guarda la clave. La lee cada vez que se inicia la cuenta, desde una variable de entorno o desde un comando como un gestor de contraseñas, y solo se la entrega a ese proceso de Claude Code.

## Uso

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

Una cuenta con clave de API se une a grupos como cualquier otra cuenta (por defecto a `default`: la configuración, `CLAUDE.md`, las extensiones y los servidores MCP), y funciona con atajos, `ccam folder` y `ccam ext`.

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--api` | Agrega una cuenta con clave de API en lugar de iniciar sesión | desactivada |
| `--key-from env:NAME` | Lee la clave de la variable de entorno `NAME` | |
| `--key-from command:COMMAND` | Ejecuta `COMMAND` (con `sh -c`, o `cmd /C` en Windows) y usa la primera línea que imprime | |
| `--base-url <url>` | El endpoint (`ANTHROPIC_BASE_URL`) | La API de Anthropic |
| `--auth api-key\|bearer` | Envía la clave como `ANTHROPIC_API_KEY` (`x-api-key`) o como `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, que suelen pedir los gateways) | `api-key` |
| `--model <model>` | El modelo (`ANTHROPIC_MODEL`) | El predeterminado de Claude Code |

La configuración se guarda en `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Cómo funciona

En cada inicio de la cuenta, después de las comprobaciones habituales y antes de que se inicie Claude Code:

1. ccam lee la clave: `env:` del entorno, `command:` ejecutando el comando con un límite de 30 segundos. Un comando puede pedir una contraseña en la terminal; su salida de error también va a la terminal. Una clave vacía detiene el inicio con un mensaje que indica dónde se esperaba la clave.
2. Quita todas las variables heredadas que deciden cómo inicia sesión Claude Code o qué proveedor usa: las claves y los tokens (`ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`), el endpoint (`ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`) y los selectores de proveedor (`CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`). Una clave definida en tu shell para otra cosa no puede llegar a esta cuenta, y un inicio de sesión de suscripción no puede imponerse a su clave.
3. Define `ANTHROPIC_BASE_URL` y `ANTHROPIC_MODEL` cuando se indican, y la clave como `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`.

Las cuentas con inicio de sesión de suscripción conservan el entorno que heredan, como antes. En una sesión interactiva, Claude Code puede preguntar una vez si debe usar la clave que encuentra en el entorno.

## Seguridad y casos límite

- La clave nunca se escribe en `config.toml`, en una caché ni en un registro, y `ccam --dry-run gw` no la lee: una ejecución de prueba no ejecuta una fuente `command:`.
- `ccam add --api` comprueba una fuente `env:` enseguida y avisa si la variable está vacía; una fuente `command:` se ejecuta por primera vez al iniciar.
- `ccam login` y `ccam logout` rechazan una cuenta con clave de API, ya que no tiene inicio de sesión. `ccam doctor` muestra su endpoint y la fuente de su clave, y avisa cuando una variable `env:` está vacía.
- `ccam usage` muestra `API` para ella. `ccam next --auto`, `--fallback` y `ccam pick` solo eligen entre cuentas de suscripción, porque una clave de API no tiene un límite de 5 horas ni semanal que comparar.
- Un `settings.json` compartido con un bloque `env` o `apiKeyHelper` se aplica a todas las cuentas del grupo, incluidas las cuentas con clave de API. Deja la configuración del proveedor fuera de la configuración compartida, o pon la cuenta con clave de API en un grupo sin `settings`.

## Relacionado

- [Administrar cuentas](/features/accounts)
- [Iniciar Claude Code](/features/launch)
- [Referencia de comandos](/docs/commands)
