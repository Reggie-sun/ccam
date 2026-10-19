---
title: Atalhos
description: Dê a uma conta o próprio comando, como ccwork, com ccam shortcut, e recrie os links em ~/.local/bin com ccam shims install.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
---

Um atalho é um comando que você mesmo nomeia e que inicia o Claude Code como uma determinada conta: `ccwork --resume` faz o mesmo que `ccam work --resume`. Ele é um link para o ccam em `~/.local/bin`, então funciona no bash, no zsh, em scripts e em qualquer outra coisa que execute programas a partir do `PATH`. Os atalhos são opcionais; crie um para as contas que você inicia com frequência.

## Uso

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

O ccam confirma cada alteração:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Uma conta tem no máximo um atalho, armazenado como `shortcut` em `~/.config/ccam/config.toml`. O assistente do `ccam add` pergunta `Create a shortcut command for <name>? [y/N]` e depois `Shortcut name [cc<name>]:`. Fora do assistente, o `ccam add` só cria um com `--shortcut`. `ccam list` mostra o atalho de cada conta, ou `-`.

Um atalho repassa todos os argumentos ao Claude Code, do mesmo jeito que `ccam <name>` faz depois do nome. Isso inclui `--dry-run` e `--help`: para ver uma inicialização antes, use `ccam --dry-run work`.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--remove` | Exclui o atalho da conta e o link dele | desativado |
| `-h`, `--help` | Ajuda do comando | |
| `--dry-run` | O `ccam shortcut` a respeita: verifica o nome e imprime `Would set the shortcut of work to ccwork.` ou `Would remove ccwork, the shortcut of work.` sem salvar nem criar links | desativado |

`ccam shims install` não aceita argumentos e recusa `--dry-run` com `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Regras de nome, verificadas nesta ordem:

- de 1 a 64 caracteres entre `a-z`, `0-9`, `.`, `-` e `_`, começando por uma letra ou um dígito;
- não pode ser `ccam`, `claude` nem um comando do ccam, como `run` ou `usage`;
- não pode ser o atalho de outra conta;
- não pode ser um arquivo em `~/.local/bin` que já não seja um link para o ccam;
- não pode ser o nome de outro executável no seu `PATH`.

Informar um comando e `--remove` juntos para com `give a command or --remove, not both`. Definir o nome que a conta já tem imprime `work: ccwork (unchanged)` e recria o link dela se ele estiver faltando.

## Como funciona

Definir um nome executa `checkShortcut`, que chama `shim.Check` com a pasta bin e o caminho real do binário do ccam:

<!-- motion:shortcuts-check -->

O nome passa pelas verificações em ordem. Ao lado de cada verificação há um nome que ela recusa: `Work` não segue o padrão, `run` é um comando do ccam, `ccside` pertence a outra conta, `ccfile` é um arquivo existente em `~/.local/bin` e `git` esconderia `/usr/bin/git`. Só um nome que passa em todas as verificações é salvo e vinculado.

1. `shim.Check` compara o nome com `^[a-z0-9][a-z0-9._-]{0,63}$`, depois com `config.Reserved` e depois com `cfg.FindShortcut`. Ele chama `os.Lstat` em `~/.local/bin/<name>` e, para cada outra pasta do `PATH`, chama `os.Stat` e recusa um arquivo executável que não seja um link para o ccam.
2. O ccam salva o `config.toml` com `fsutil.AtomicWrite` (arquivo temporário, `fsync`, rename).
3. `shim.Remove` exclui o link do atalho antigo, mas só se ele for um link para o ccam.
4. `shim.Install` cria o link de cada atalho configurado que esteja faltando, com `os.Symlink` para o binário do ccam, localizado com `os.Executable` e `filepath.EvalSymlinks`. Ele nunca remove nada.

Um link conta como do próprio ccam quando `os.Readlink` retorna o caminho do ccam, ou quando o link e o caminho do ccam levam ao mesmo arquivo. Toda remoção e toda verificação de sobrescrita passa por esse teste, então o ccam não mexe em arquivos que ele não criou.

Quando você executa um atalho, o nome que você digitou é tudo de que o ccam precisa:

<!-- motion:shortcuts-link -->

O shell encontra `ccwork` no `PATH`, o kernel segue o link e inicia o binário do ccam, e `argv[0]` continua sendo `ccwork`. Em `app.Main`, qualquer nome diferente de `ccam` carrega o `config.toml` uma vez e chama `shim.Lookup`, que encontra a conta cujo `shortcut` é igual a `filepath.Base(argv[0])`. Em seguida, `launcher.LaunchWith` segue as mesmas etapas que [`ccam <name>`](/features/launch), com a configuração já carregada, e passa `argv[1:]` ao `claude` sem alterações.

Por que um link simbólico e não um alias de shell ou um script wrapper: um link simbólico funciona em todos os shells e em scripts, não adiciona a inicialização de mais um shell (cerca de 1 a 2 ms), e a inicialização continua lendo um único arquivo. `scripts/perf-check.sh` mede um atalho com `hyperfine` contra o mesmo limite de 10 ms (p95) de `ccam <name>`.

## Segurança e casos especiais

- **Um link que nenhuma conta reivindica.** Se um link para o ccam em `~/.local/bin` tem um nome que nenhuma conta usa, por exemplo depois que você editou `shortcut` no `config.toml` manualmente, o ccam não recorre à conta padrão. Ele para com `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Defina o nome de novo com `ccam shortcut` ou exclua o link.
- **Arquivos existentes.** O ccam nunca substitui um arquivo que ele não criou. `ccam shims install` para com `<path> already exists and is not a ccam link, so it was left alone`.
- **Um binário do ccam movido.** Os links antigos passam a apontar para um caminho que não existe mais, e o ccam não os considera dele. Exclua-os e execute `ccam shims install`.
- **Renomear e remover.** `ccam rename` mantém o nome do atalho, e o link continua funcionando, porque aponta para o ccam e o ccam procura o nome no `config.toml`. `ccam remove` exclui o link da conta. `ccam shortcut <account> <new>` exclui o link antigo antes de criar o novo.
- **PATH.** Quando o ccam cria um link e `~/.local/bin` não está no seu `PATH`, ele imprime `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` seguido da linha `export PATH=...` a adicionar.
- **Verificar depois.** `ccam doctor` informa `missing shortcuts in /home/me/.local/bin: ccwork` e avisa quando a pasta não está no `PATH`. `ccam doctor --fix` cria os links que faltam.

## Relacionados

- [Iniciar o Claude Code](/features/launch)
- [Gerenciar contas](/features/accounts)
- [Referência de comandos](/docs/commands#shortcuts)
- [Instalar o ccam](/docs/install)
