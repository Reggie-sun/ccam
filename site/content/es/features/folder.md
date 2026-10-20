---
title: Cuentas por carpeta
description: ccam folder hace que una carpeta y todo su contenido usen una cuenta al iniciar Claude Code sin nombre de cuenta, como work para el código de tu empresa.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` hace que todo inicio sin nombre de cuenta dentro de `~/code/company`, y en cualquier carpeta debajo de ella, inicie Claude Code como `work`. Fuera de ella, se inicia la cuenta predeterminada como antes. Si nombras una cuenta, con `ccam personal` o un atajo como `ccpersonal`, siempre se inicia esa cuenta.

## Uso

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Estos inicios usan la cuenta de la carpeta:

- `ccam -p "…"` y cualquier otro inicio cuyo primer argumento sea una opción para Claude Code
- `ccam run` sin nombre, y `ccam --dry-run`
- `ccam --fallback -p "…"` en su primer intento
- el selector, `ccam` a secas, que empieza con el cursor en la cuenta de la carpeta

## Opciones

| Opción | Qué hace | Valor predeterminado |
| --- | --- | --- |
| `<account>` | Nombre o número de la cuenta que se usará | |
| `[folder]` | La carpeta, relativa o absoluta | la carpeta actual |
| `--remove` | Quita la entrada de exactamente esa carpeta | desactivada |
| `--dry-run` | Se respeta: imprime el cambio sin guardarlo | desactivada |

## Cómo funciona

Cada carpeta es una entrada `[[folder]]` en `~/.config/ccam/config.toml`, el archivo que cada inicio lee de todos modos:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` convierte la carpeta en una ruta absoluta con los enlaces simbólicos resueltos y la guarda. Volver a asignar una carpeta reemplaza su cuenta.
2. Un inicio sin nombre llama a `launcher.HereProfile`. Sin entradas, toma la cuenta predeterminada de inmediato.
3. Si no, compara el directorio de trabajo con cada `path` como cadenas: la propia carpeta o una debajo de ella, nunca una carpeta que solo empieza con las mismas letras (`/code/company` no abarca `/code/companyx`). Gana la entrada más profunda, así que `~/code` puede ser `personal` mientras `~/code/company` es `work`.
4. Solo si el directorio de trabajo tal cual no coincide con nada, ccam resuelve sus enlaces simbólicos y compara una vez más.
5. En Windows, la comparación no distingue mayúsculas de minúsculas, igual que Windows.

No se lee ningún archivo aparte de `config.toml`, así que un inicio dentro de una carpeta no cuesta un tiempo medible (4,0 ms con o sin entradas, `hyperfine`, 300 ejecuciones).

## Seguridad y casos límite

- `ccam rename` traslada las entradas de carpeta al nuevo nombre, y `ccam remove` elimina las entradas de la cuenta quitada.
- Una entrada cuya cuenta ya no existe detiene un inicio sin nombre en esa carpeta: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` informa de esas entradas y de las entradas cuya carpeta ya no está.
- `ccam folder --remove` dentro de una carpeta que obtiene su cuenta de una carpeta superior indica esa carpeta en lugar de quitar nada.
- `ccam next` elige por rotación o por uso y no tiene en cuenta las carpetas.

## Relacionado

- [Iniciar Claude Code](/features/launch)
- [Atajos](/features/shortcuts)
- [Administrar cuentas](/features/accounts)
- [Referencia de comandos](/docs/commands)
