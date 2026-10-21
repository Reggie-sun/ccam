---
title: Actualizar ccam
description: ccam update reemplaza ccam por la última versión de GitHub tras comprobar su SHA-256; --check solo informa y --version instala una versión concreta.
updated: 2026-10-21
order: 10
translation: machine
sourceHash: a92dbf37399f
---

`ccam update` descarga una versión de ccam desde GitHub, la comprueba con el `checksums.txt` de la versión y reemplaza el binario en ejecución con un solo cambio de nombre. Úsalo si instalaste ccam con el script de instalación o desde un archivo comprimido. El [panel](/features/dashboard) también pregunta a GitHub por la última versión cada vez que se abre y la ofrece con `U`, que ejecuta los mismos pasos después de preguntar. Ningún otro comando se conecta a GitHub, así que funcionan igual sin conexión.

## Uso

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Cuando no hay una versión más reciente, ccam imprime `ccam 0.4.0 is up to date.` y no cambia nada. Con `--version` y la versión que ya ejecutas, imprime `ccam 0.4.0 is already installed.`

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `--check` | Solo informa si hay una versión más reciente | desactivada |
| `--version <tag>` | Instala esta etiqueta de versión, como `v0.3.0`, más reciente o más antigua que la actual | la más reciente |
| `--dry-run` | Se respeta: busca la etiqueta más reciente (salvo que `--version` indique una), ejecuta las mismas comprobaciones sobre el binario e imprime las rutas de descarga, suma de verificación y reemplazo sin descargar | desactivada |

`--check` y `--version` no se pueden usar a la vez. La etiqueta debe tener la forma `vMAJOR.MINOR.PATCH`, con un sufijo opcional de versión preliminar o de compilación; la `v` inicial se puede omitir. Cualquier otra cosa falla, por ejemplo `--version takes a release tag such as v0.3.0, not "1.2"`.

La variable de entorno `CCAM_RELEASES_URL` reemplaza `https://github.com/jhyunwoo/ccam/releases`, para un servidor espejo o para pruebas. Debe tener la misma estructura: `<url>/latest` redirige a `<url>/tag/<tag>`, y los archivos están en `<url>/download/<tag>/<asset>`.

## Cómo funciona

El código está en `internal/update`, que solo importa `ccam update`. La ruta de inicio nunca lo llama.

1. `cli.updateCmd` encuentra el binario en ejecución con `os.Executable` y resuelve los enlaces simbólicos con `filepath.EvalSymlinks`. Esa ruta real es la que se reemplaza.
2. Sin `--version`, `Updater.Latest` solicita `<releases>/latest` sin seguir la redirección y lee la etiqueta del encabezado `Location` (`…/tag/v0.5.0`). Esto usa la página de versiones en lugar de la API de GitHub, así que no se aplica el límite de 60 solicitudes por hora de la API. La consulta tiene 10 segundos.
3. `parseVersion` comprueba la etiqueta con un patrón estricto. Una etiqueta acaba en las URL de descarga y en un comando de shell impreso, así que no pasa nada que tenga `/`, `%`, espacios o caracteres de shell. `compareVersions` la compara con la versión en ejecución; una versión preliminar se ordena antes que su versión final.
4. `managed()` y `writable()` deciden si ccam puede reemplazar el binario (consulta más abajo).
5. `Updater.Fetch` descarga `ccam_<os>_<arch>.tar.gz` (en Windows, `ccam_windows_<arch>.zip`) y `checksums.txt` a la vez, con 2 minutos para ambos. Cada descarga tiene un límite de 256 MB.
6. El SHA-256 del archivo comprimido debe coincidir con su línea en `checksums.txt` antes de que ocurra cualquier otra cosa. Una línea que falta o una suma de verificación distinta detiene la actualización con `nothing was changed`.
7. `extract()` lee el tar y toma el primer archivo normal llamado `ccam` (`extractZip` toma `ccam.exe` del zip). Las carpetas, los enlaces y las demás entradas se ignoran, y el binario también tiene un límite de 256 MB.
8. `Updater.Install` escribe el binario nuevo con `fsutil.AtomicWrite`: un archivo temporal `.ccam.tmp*` junto al binario, con los bits de permiso del archivo anterior, volcado al disco y renombrado sobre el anterior.

El proceso completo, desde la búsqueda de la etiqueta hasta el cambio de nombre:

<!-- motion:update-pipeline -->

El archivo comprimido y `checksums.txt` se descargan en paralelo, y no se escribe nada hasta que la suma de verificación coincide. El cambio de nombre es atómico, así que un fallo o un disco lleno deja el binario anterior o el nuevo, nunca la mitad de uno. Los atajos como `~/.local/bin/ccwork` son enlaces simbólicos al binario de ccam, así que ejecutan la versión nueva sin ningún cambio.

En Windows, el archivo de un programa en ejecución no se puede reemplazar ni eliminar, solo renombrar. `Updater.Install` escribe el binario nuevo en `ccam.exe.new`, renombra el `ccam.exe` en ejecución a `ccam.exe.old` y renombra el archivo nuevo a `ccam.exe`. Si el último paso falla, el binario anterior recupera su nombre. La siguiente actualización elimina `ccam.exe.old`. Los atajos que son enlaces duros en lugar de enlaces simbólicos ([Atajos](/features/shortcuts)) se renuevan justo después, así que también ejecutan la versión nueva.

### Instalaciones administradas por paquetes

Un ccam instalado desde el `.deb` pertenece a apt, y reemplazarlo a espaldas de apt dejaría incorrecta la base de datos de paquetes. Antes de descargar, ccam comprueba a quién pertenece el binario:

<!-- motion:update-managed -->

1. Un binario en `/usr/bin`, `/bin`, `/usr/sbin` o `/sbin` se considera administrado por paquetes.
2. En otro lugar, si `dpkg` está instalado y `dpkg -S <path>` tiene éxito, también se considera administrado por paquetes.
3. En caso contrario, ccam crea y elimina un archivo temporal `.ccam-update-*` junto al binario para asegurarse de que puede escribir ahí.

Con un binario administrado por paquetes, ccam se detiene e imprime los pasos de apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

En sistemas que no son Linux, solo se imprime la línea del script de instalación: `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex` en Windows. Cuando no se puede escribir en la carpeta, ccam sugiere volver a ejecutar el script de instalación, que instala ccam en `~/.local/bin`.

## Seguridad y casos límite

- Una compilación de desarrollo informa su versión como `dev` y no tiene una versión publicada con la cual compararse, así que `ccam update` y `ccam update --check` la rechazan. `ccam update --version vX.Y.Z` instala una versión publicada encima.
- Una descarga fallida, un archivo que falta para tu plataforma o una suma de verificación que no coincide dejan el binario instalado como estaba.
- Una versión que no existe falla con `there is no release <tag>`. Cuando GitHub rechaza la solicitud con HTTP 403 o 429, ccam indica que probablemente sea un límite de solicitudes y te pide que lo intentes de nuevo más tarde.
- `--check` solo lee. No comprueba si el binario está administrado por paquetes ni si se puede escribir en él.
- `--version` con una etiqueta más antigua es la forma de volver a una versión anterior.

## Relacionado

- [Instalar](/docs/install)
- [Comandos](/docs/commands)
- [Preguntas frecuentes](/docs/faq)
