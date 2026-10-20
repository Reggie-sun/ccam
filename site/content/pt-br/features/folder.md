---
title: Contas por pasta
description: ccam folder faz uma pasta e tudo o que está dentro dela usarem uma conta quando você inicia o Claude Code sem nome de conta, como work para o código da empresa.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: "46d1410a3d81"
---

`ccam folder work ~/code/company` faz com que toda inicialização sem nome de conta dentro de `~/code/company`, e em qualquer pasta abaixo dela, inicie o Claude Code como `work`. Fora dela, a conta padrão inicia como antes. Nomear uma conta, com `ccam personal` ou com um atalho como `ccpersonal`, sempre inicia essa conta.

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

Estas inicializações usam a conta da pasta:

- `ccam -p "…"` e qualquer outra inicialização cujo primeiro argumento seja uma opção para o Claude Code
- `ccam run` sem nome, e `ccam --dry-run`
- `ccam --fallback -p "…"` na primeira tentativa
- o painel, `ccam` sozinho, que começa com o cursor na conta da pasta

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `<account>` | Nome ou número da conta a usar | |
| `[folder]` | A pasta, relativa ou absoluta | a pasta atual |
| `--remove` | Remove a entrada exatamente dessa pasta | desativado |
| `--dry-run` | Respeitada: imprime a mudança sem salvá-la | desativado |

## Como funciona

Cada pasta é uma entrada `[[folder]]` em `~/.config/ccam/config.toml`, o arquivo que toda inicialização já lê:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` resolve a pasta para um caminho absoluto, com os links simbólicos resolvidos, e o salva. Definir uma pasta de novo substitui a conta dela.
2. Uma inicialização sem nome chama `launcher.HereProfile`. Sem entradas, ela pega a conta padrão na hora.
3. Caso contrário, compara o diretório de trabalho com cada `path` como strings: a própria pasta ou uma abaixo dela, nunca uma pasta que só começa com as mesmas letras (`/code/company` não cobre `/code/companyx`). A entrada mais profunda vence, então `~/code` pode ser `personal` enquanto `~/code/company` é `work`.
4. Só se o diretório de trabalho, como está, não corresponder a nada, o ccam resolve os links simbólicos dele e compara mais uma vez.
5. No Windows, a comparação ignora maiúsculas e minúsculas, como o Windows faz.

Nenhum arquivo além do `config.toml` é lido, então uma inicialização dentro de uma pasta não custa tempo mensurável (4,0 ms com ou sem entradas, `hyperfine`, 300 execuções).

## Segurança e casos especiais

- `ccam rename` leva as entradas de pasta para o novo nome, e `ccam remove` exclui as entradas da conta removida.
- Uma entrada cuja conta não existe mais interrompe uma inicialização sem nome nessa pasta: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` informa essas entradas e as entradas cuja pasta sumiu.
- `ccam folder --remove` dentro de uma pasta que recebe a conta de uma pasta acima indica essa pasta em vez de remover alguma coisa.
- `ccam next` escolhe por rotação ou por uso e não olha as pastas.

## Relacionados

- [Iniciar o Claude Code](/features/launch)
- [Atalhos](/features/shortcuts)
- [Gerenciar contas](/features/accounts)
- [Referência de comandos](/docs/commands)
