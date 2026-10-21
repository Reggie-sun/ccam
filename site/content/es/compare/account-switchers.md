---
title: "Cambiar cuenta de Claude Code: ccam y alternativas"
description: ccam frente a claude-swap, CCS, clausona, clauth, CAAM, cux y otras herramientas para cambiar de cuenta en Claude Code, y qué enfoque conviene en cada caso.
updated: 2026-10-21
order: 2
translation: machine
sourceHash: c74be64d2978
---

Las herramientas para cambiar de cuenta en Claude Code siguen uno de dos enfoques. La mayoría intercambia el inicio de sesión dentro de la única carpeta `~/.claude`, así que solo hay una cuenta activa a la vez. Otras, entre ellas ccam, le dan a cada cuenta su propia carpeta con `CLAUDE_CONFIG_DIR`, así que todas las cuentas siguen con la sesión iniciada y varias pueden ejecutarse a la vez. ccam agrega configuración y servidores MCP compartidos, y continúa una conversación en la carpeta de otra cuenta.

## ¿Intercambiar un inicio de sesión o usar carpetas separadas?

**Intercambiar un inicio de sesión** mantiene un solo conjunto de configuración, servidores MCP e historial, porque solo hay una carpeta. El cambio es rápido, y una sesión en curso puede tomar el nuevo inicio de sesión. Pero solo hay una cuenta activa a la vez, para todas las terminales y editores a la vez.

Para hacer el intercambio, la herramienta guarda una copia del inicio de sesión de cada una de las demás cuentas y la vuelve a colocar en su sitio. Claude Code reemplaza el token de actualización de un inicio de sesión cada vez que lo renueva, así que cada copia guardada tiene que mantenerse al día con las sesiones en curso. Las cuentas que quedan con la sesión cerrada o que piden `/login` de nuevo tras un cambio están entre los problemas que estos proyectos corrigen con más frecuencia.

**Las carpetas separadas** son el método que [describe la documentación de Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Todas las cuentas siguen con la sesión iniciada, y puedes ejecutar una sesión de trabajo y una personal en paralelo. Cada inicio de sesión se queda en su propia carpeta, donde Claude Code lo renueva como siempre. El precio es que cada carpeta empieza con su propia configuración, sus servidores MCP y su historial, a menos que una herramienta los comparta.

## Comparación de las herramientas

Tomado de la documentación y el código fuente de cada proyecto el 2026-10-20 (clausona: v0.12.0-beta). "No documentado" significa que el proyecto no lo dice, no que la función falte.

| Herramienta | Enfoque | Varias cuentas a la vez | Uso de 5 horas y semanal | Al llegar a un límite | Plataformas |
| --- | --- | --- | --- | --- | --- |
| ccam | Carpetas separadas | Sí | Sí, `ccam usage` | Una notificación en la sesión; `ccam next --auto --continue` o `/ccam:switch` lleva la conversación a la cuenta con más margen, con o sin historial compartido; `ccam --fallback` vuelve a ejecutar allí los trabajos `-p` por sí solo, y `ccam pick` reparte los workers | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Intercambia el inicio de sesión; carpetas separadas con `cswap run` | Sí, con `cswap run` | Sí | `cswap auto` cambia el inicio de sesión al superar un umbral; la sesión en curso continúa | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Carpetas separadas | Sí | En parte | Los grupos de cuentas pasan las solicitudes a la siguiente cuenta mediante un proxy local | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Carpetas separadas | Sí | Sí | Las rutas eligen una cuenta con cuota antes de cada ejecución, pueden esperar a un restablecimiento y reparten cuentas distintas entre los workers | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Intercambia el inicio de sesión; carpetas separadas con `clauth start` | Sí, con `clauth start` | Sí | Cambia siguiendo una cadena de respaldo | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Intercambia el inicio de sesión o usa perfiles aislados | Sí, con perfiles aislados | Sí | `caam run` reintenta los comandos no interactivos en otra cuenta | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Intercambia el inicio de sesión y envuelve `claude` | No, una cuenta activa | Sí | Intercambia el inicio de sesión y reanuda la misma conversación | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Intercambia el inicio de sesión | No | Sí | Cambio opcional al superar un umbral | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Intercambia el inicio de sesión | No | Sí | Cambio opcional al 100% | macOS |

ccswap es un fork de claude-swap que también gestiona inicios de sesión de Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) y [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) también le dan a cada cuenta su propia carpeta. Varias de estas herramientas también manejan otros agentes de programación, como Codex o Gemini.

## Cómo están construidas

| Herramienta | Escrita en | Mientras Claude Code se ejecuta | Copia inicios de sesión de un lugar a otro | Qué pueden compartir las cuentas |
| --- | --- | --- | --- | --- |
| ccam | Go, un único binario estático | Nada: ccam se reemplaza a sí mismo por `claude` (en Windows, que no puede hacerlo, lo espera; un relé opcional se queda a su lado) | No | Grupos: cada uno nombra los elementos que comparte y sus cuentas; una cuenta puede estar en varios |
| claude-swap | Python 3.12+, instalado con uv o pipx | Nada con `cswap run`; `cswap auto` y la app de la barra de menús siguen ejecutándose | Sí, cuando intercambia el inicio de sesión | Con `cswap run`: configuración, `CLAUDE.md` y extensiones juntas, historial opcional; servidores MCP copiados en cada inicio |
| CCS | Node.js 18+ o Bun | CCS sigue como proceso padre; los grupos de cuentas añaden un proxy local | No en sus carpetas de cuenta; los grupos de cuentas guardan sus propios inicios de sesión | Compartido o local por cuenta, historial en grupos de contexto; servidores MCP copiados |
| clausona | TypeScript, Node.js 20+ | `claude` a secas pasa por una función de shell; `clausona run` mantiene Node como proceso padre; un proceso de Node en segundo plano tras cada salida | No, pero renueva por sí mismo los tokens de las cuentas inactivas para leer su uso | Todo excepto los archivos propios de cada cuenta, desde una cuenta principal; historial opcional; los servidores añadidos con `claude mcp add` no se comparten |
| clauth | Rust | clauth sigue como proceso padre | Sí, cuando intercambia el inicio de sesión | Todo `~/.claude`, o nada |
| CAAM | Go | Depende del modo; un daemon opcional | Sí, en su modo de bóveda | Fijado por el modo; servidores MCP copiados en cada inicio |
| cux | Go | Un wrapper sigue envolviendo `claude` | Sí | Todo: una sola carpeta |
| CCSwitcher, Claude Switcher | Apps de macOS (Swift, Python) | La app de la barra de menús | Sí | Todo: una sola carpeta |

## Lo que ccam hace bien

### Los inicios de sesión se quedan donde Claude Code los dejó

ccam nunca copia un inicio de sesión de Claude de un lugar a otro. El inicio de sesión de cada cuenta se queda en su propia carpeta, o en macOS en el ítem del llavero propio de esa carpeta, y Claude Code es el programa que lo renueva. No hay ninguna copia guardada que pueda desincronizarse. Un `~/.claude` existente registrado con `ccam import` conserva el ítem del llavero que usa `claude` a secas, así que sigue con la sesión iniciada en ambos casos.

ccam solo renueva por sí mismo un token caducado cuando ejecutas `ccam usage --refresh`. Antes toma un bloqueo y se niega mientras esa cuenta esté en ejecución, así que nunca rota un token bajo una sesión activa.

### Nada se ejecuta entre tú y Claude Code

ccam lee un archivo, comprueba sus enlaces y se reemplaza a sí mismo por `claude` (`exec`). Mientras trabajas no hay ningún proceso de ccam, wrapper, proxy ni daemon: la entrada estándar, las señales y el código de salida pertenecen a Claude Code, y sus solicitudes van directamente a Anthropic. Windows no puede reemplazar un proceso, así que allí ccam inicia Claude Code con la misma consola y espera su código de salida, sin hacer nada más.

Iniciar una cuenta así añade unos 4 ms, menos de 6 ms en el percentil 95, con 8 MB de memoria máxima ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) con hyperfine, 300 ejecuciones). En el mismo portátil, solo arrancar el intérprete de Python tarda 13 ms y Node.js 31 ms, antes de que un programa haga nada.

### Los grupos deciden qué se comparte

Los servidores MCP, `settings.json`, `CLAUDE.md`, las extensiones (agentes, comandos, skills, plugins) y el historial de conversaciones son cinco elementos separados. Un grupo nombra los elementos que comparte y las cuentas que contiene, así que dos cuentas de trabajo pueden compartir el historial en un grupo `work` mientras una cuenta personal conserva el suyo, y las tres comparten los mismos servidores MCP a través de `default`. `ccam group set work --add memory` mueve todas las cuentas del grupo en un solo plan; el [panel](/features/dashboard) hace lo mismo con casillas y atenúa un elemento que una cuenta ya recibe de otro grupo.

- Los servidores MCP se pasan con `--mcp-config` al iniciar, así que por defecto ccam nunca edita `.claude.json`, el archivo que Claude Code reescribe mientras se ejecuta. `ccam mcp add` agrega un servidor una sola vez para todo el grupo.
- Los demás elementos son enlaces a la carpeta del grupo, así que un cambio hecho en una cuenta aparece en las otras al instante. No se copia nada al iniciar.
- Los inicios de sesión, `.claude.json` y el estado de ejecución nunca se comparten.

### Cada cambio se muestra primero y se respalda

`ccam group` y `ccam migrate` muestran un plan y preguntan antes de mover nada. Cuando un archivo difiere entre una cuenta y su grupo, tú eliges: conservar el del grupo, usar el tuyo o combinar el JSON. Los originales van a `.ccam-backup/<time>/`, la carpeta de una cuenta eliminada va a una carpeta de papelera, y los archivos se reemplazan con un solo cambio de nombre, así que una escritura interrumpida deja el archivo anterior. `--dry-run` muestra lo que harían el inicio, el uso compartido, la migración, los atajos, las actualizaciones y las reparaciones. `ccam doctor --fix` repara enlaces rotos y permisos.

### La misma conversación en la siguiente cuenta

`ccam next --auto --continue` consulta el uso de todas las cuentas a la vez, abre la que tiene más margen y continúa la conversación. Con el historial compartido es la misma conversación; sin él, Claude Code la copia en la nueva cuenta (`--resume <file> --fork-session`), y ccam indica cuántos tokens vuelve a leer la primera respuesta. Cuando un límite detiene una sesión, el hook del plugin de ccam muestra una notificación, y `/ccam:switch` pasa la conversación a otra cuenta. Tú decides cuándo cambiar; nada vigila tu uso en segundo plano, y las elecciones automáticas se quedan dentro de una misma organización, así que una conversación de la empresa nunca acaba por accidente en una cuenta privada. La línea de estado sugiere el comando cuando una cuenta llega al 90 %.

Los scripts no pueden esperarte, así que `ccam --fallback work -p "…"` cambia por ellos: cuando la cuenta llega a su límite, el mismo comando se ejecuta de nuevo con la cuenta con más margen, y solo se imprime la salida de la ejecución que terminó. Un fallo que no sea un límite de uso se transmite tal cual. Para workers en paralelo, `ccam pick` le da a cada uno una cuenta distinta con margen, y ambos salen con el código 75 cuando ya ninguna cuenta tiene margen.

### La cuenta correcta en cada carpeta

`ccam folder work ~/code/company` hace que todo inicio sin nombre de cuenta dentro de esa carpeta use `work`, mientras el resto de tus proyectos siguen con la cuenta predeterminada. El ajuste vive en el mismo `config.toml` que lee cada inicio, así que no cuesta tiempo.

### Un solo binario, descargas verificadas

ccam es un binario estático de unos 10 MB, sin entorno de ejecución que instalar, para Linux, macOS y Windows. Los scripts de instalación (`install.sh`, `install.ps1`) y `ccam update` verifican cada descarga con las sumas de comprobación SHA-256 de la versión, y ccam nunca se actualiza sin preguntar. No tiene telemetría: las únicas solicitudes que hace son las consultas de uso de tus cuentas a Anthropic, en paralelo con un límite de 3 segundos y una caché de 60 segundos, la pregunta del panel a GitHub por la última versión cuando se abre, y las descargas que tú pides. Con `usage_source = "claude"`, incluso esas pasan por Claude Code, y ccam nunca maneja un token.

## Dónde otras herramientas hacen más

- claude-swap, clauth, cux, CAAM y CCSwitcher pueden cambiar por su cuenta una sesión interactiva en curso cuando una cuenta supera un umbral de uso. ccam cambia las ejecuciones `-p` con `--fallback` y deja las sesiones interactivas en tus manos, con una notificación al llegar al límite y `/ccam:switch`.
- En Windows sin el modo de desarrollador, ccam comparte carpetas pero no puede compartir archivos sueltos como `settings.json`; claude-swap y CCS los copian allí en su lugar.
- CCSwitcher, Claude Switcher y claude-swap muestran el uso en la barra de menús de macOS.
- CCS, clausona, CAAM, clauth, ccswap y Claude Switcher también gestionan Codex y otros agentes de programación.
- clausona también activa y desactiva extensiones por cuenta (`ccam ext` de ccam solo las muestra), registra el coste local en tokens e incluye plugins de Claude Code que ejecutan flotas de workers en tmux y otras terminales.
- CAAM, CCSwitcher y Claude Switcher se instalan con Homebrew, y clauth y CAAM firman las sumas de comprobación de sus versiones. ccam se instala con su script, un `.deb` o `go install`, y verifica las sumas de comprobación sin firma.

## ¿Cuál te conviene?

Elige una herramienta que intercambie el inicio de sesión si usas una cuenta a la vez o quieres que las sesiones en curso cambien automáticamente. Elige ccam si quieres cuentas en paralelo, una sola configuración compartida, la misma conversación en la siguiente cuenta y nada ejecutándose junto a Claude Code, en Linux, macOS o Windows. Si no necesitas compartir nada, bastan [simples alias del shell](/compare/aliases).

Si puedes usar varias cuentas depende de los términos de cada una. Revísalos antes de hacerlo.
