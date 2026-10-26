---
title: Atalhos
description: Dê a uma conta o próprio comando, como ccwork, com ccam shortcut, e recrie os links em ~/.local/bin com ccam shims install.
updated: 2026-10-26
order: 3
translation: machine
sourceHash: "066e7e622517"
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
ccam shims install            # recreate missing or outdated links of all accounts
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
| `--dry-run` | O `ccam shortcut` a aceita: verifica o nome e imprime `Would set the shortcut of work to ccwork.` ou `Would remove ccwork, the shortcut of work.` sem salvar nem criar links | desativado |

`ccam shims install` não aceita argumentos e recusa `--dry-run` com `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Regras de nome, verificadas nesta ordem:

- de 1 a 64 caracteres entre `a-z`, `0-9`, `.`, `-` e `_`, começando por uma letra ou um dígito;
- não pode ser `ccam`, `claude` nem um comando do ccam, como `run` ou `usage`;
- não pode ser o atalho de outra conta;
- não pode ser um arquivo existente em `~/.local/bin`, a menos que seja um link que o ccam criou (como o próprio atalho desta conta);
- não pode ser o nome de outro executável no seu `PATH`.

O `config.toml` é verificado com as mesmas regras quando é carregado. Se você editar um atalho à mão de forma que ele quebre o padrão, seja `claude` ou `ccam`, ou seja usado por duas contas, todos os comandos falham com uma mensagem que indica a conta, até que você o corrija.

Informar um comando e `--remove` juntos para com `give a command or --remove, not both`. Definir o nome que a conta já tem imprime `work: ccwork (unchanged)` e recria o link dela se ele estiver faltando.

## Como funciona

Definir um nome executa `checkShortcut`, que chama `shim.Check` com a pasta bin e o caminho real do binário do ccam:

<!-- motion:shortcuts-check -->

O nome passa pelas verificações em ordem. Ao lado de cada verificação há um nome que ela recusa: `Work` não segue o padrão, `run` é um comando do ccam, `ccside` pertence a outra conta, `ccfile` é um arquivo existente em `~/.local/bin` e `git` esconderia `/usr/bin/git`. Só um nome que passa em todas as verificações é salvo e vinculado.

1. O ccam compara o nome com `^[a-z0-9][a-z0-9._-]{0,63}$`, depois com `config.Reserved` e depois com `cfg.FindShortcut` (`shim.Check`). Ele chama `os.Lstat` em `~/.local/bin/<name>` e, para cada outra pasta do `PATH`, chama `os.Stat` e recusa um arquivo executável que não seja um link para o ccam.
2. O ccam salva o `config.toml` com `fsutil.AtomicWrite` (arquivo temporário, `fsync`, rename).
3. O ccam exclui o link do atalho antigo, mas só se ele for um link para o ccam (`shim.Remove`).
4. O ccam percorre cada atalho configurado em ordem de nome e cria o link de cada um que esteja faltando (`shim.Install`), com `os.Symlink` para o binário do ccam, localizado com `os.Executable` e `filepath.EvalSymlinks`. Um arquivo no caminho de um atalho não impede os outros; o erro indica cada um.

O ccam registra cada link que cria, e o binário que cada link executa, em `~/.local/state/ccam/shortcuts.json`. Um link conta como do próprio ccam quando `os.Readlink` retorna o caminho do ccam, ou quando o link e o caminho do ccam levam ao mesmo arquivo. Um link registrado cujo binário foi movido ou substituído desde então (o binário do ccam movido para outra pasta, uma reinstalação) está desatualizado, e o ccam o aponta para este binário do ccam (`shim.Install`). Um link registrado para outra cópia instalada do ccam fica como está. Toda remoção e toda verificação de sobrescrita passa por esses testes, então o ccam não mexe em arquivos que ele não criou, como uma cópia sua do ccam ou um link que você mesmo criou.

Quando você executa um atalho, o nome que você digitou é tudo de que o ccam precisa:

<!-- motion:shortcuts-link -->

O shell encontra `ccwork` no `PATH`, o kernel segue o link e inicia o binário do ccam, e `argv[0]` continua sendo `ccwork`. Em `app.Main`, qualquer nome diferente de `ccam` carrega o `config.toml` uma vez e chama `shim.Lookup`, que encontra a conta cujo `shortcut` é igual a `filepath.Base(argv[0])`. Em seguida, o ccam segue as mesmas etapas que [`ccam <name>`](/features/launch), com a configuração já carregada, e passa `argv[1:]` ao `claude` sem alterações (`launcher.LaunchWith`).

Por que um link simbólico e não um alias de shell ou um script wrapper: um link simbólico funciona em todos os shells e em scripts, evita iniciar mais um shell, o que custaria cerca de 1 a 2 ms, e a inicialização continua lendo um único arquivo. `scripts/perf-check.sh` mede um atalho com `hyperfine` contra o mesmo limite de 10 ms (p95) de `ccam <name>`.

### No Windows

Um atalho é `ccwork.exe` em `%USERPROFILE%\.local\bin`, e o ccam remove o `.exe` (sem diferenciar maiúsculas e minúsculas) antes da busca (`shim.Name`), então `ccwork`, `ccwork.exe` e `CCWORK.EXE` iniciam todos `work`. O ccam primeiro cria um link simbólico, que o Windows permite no Modo de Desenvolvedor ou a um administrador. Sem essa permissão, ele cria um hard link: um segundo nome para o mesmo arquivo `ccam.exe`, que não exige privilégio. Um hard link conta como do ccam quando é o mesmo arquivo que `ccam.exe`.

Um hard link sempre executa o arquivo a partir do qual foi criado, então, quando `ccam update` coloca um arquivo novo no lugar de `ccam.exe`, ele renova cada atalho criado como hard link com `shim.Relink` e imprime `✓ shortcut ccwork renewed`. Executar o `install.ps1` de novo faz o mesmo por meio de `ccam shims install` (`OK 2 shortcuts now run the new ccam`), e o mesmo fazem `ccam shims install` e `ccam doctor --fix` para qualquer atalho que seja uma cópia de uma versão mais antiga do ccam. A verificação de nome procura `<name>` com cada extensão de `PATHEXT` (`.exe`, `.cmd` e assim por diante) nas outras pastas do `PATH`.

Digite um atalho em um shell (PowerShell, cmd.exe, Git Bash). Um link simbólico aberto pelo Explorador de Arquivos, pela caixa de diálogo Executar ou por `Start-Process` inicia como o `ccam` simples, que abre o painel em vez da conta.

## Segurança e casos especiais

- **Um link que nenhuma conta reivindica.** Se um link para o ccam em `~/.local/bin` tem um nome que nenhuma conta usa, por exemplo depois que você editou `shortcut` no `config.toml` manualmente, o ccam não recorre à conta padrão. Ele para com `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Defina o nome de novo com `ccam shortcut` ou exclua o link.
- **Arquivos existentes.** O ccam nunca substitui um arquivo que ele não criou. `ccam shims install` informa `<path> already exists and is not a ccam link, so it was left alone` e ainda cria os outros links. `ccam add` e `ccam import` então imprimem `shortcut: ccwork (not created; see above)`.
- **Um binário do ccam movido.** Se você mover o binário do ccam, os links que o ccam criou passam a apontar para um caminho que não existe mais. `ccam shims install` (ou `ccam doctor --fix`) os aponta para o binário do ccam com que você o executa.
- **Renomear e remover.** `ccam rename` mantém o nome do atalho, e o link continua funcionando, porque aponta para o ccam e o ccam procura o nome no `config.toml`. `ccam remove` exclui o link da conta. `ccam shortcut <account> <new>` exclui o link antigo antes de criar o novo.
- **PATH.** Quando o ccam cria um link e `~/.local/bin` não está no seu `PATH`, ele imprime `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` seguido da linha `export PATH=...` a adicionar. No Windows, ele imprime em vez disso a linha do PowerShell que adiciona a pasta ao `PATH` do seu usuário.
- **Verificar depois.** `ccam doctor` informa `missing shortcuts in /home/me/.local/bin: ccwork`, atalhos que executam uma versão mais antiga do ccam, e avisa quando a pasta não está no `PATH`. `ccam doctor --fix` cria os links que faltam e renova os desatualizados, e diz `→ fixed` só para os que criou. Um atalho que executa outro ccam que ainda está instalado funciona, então o `ccam doctor` o lista só como informação; `ccam doctor --fix` o aponta para o binário do ccam com que você o executa, e `ccam add` e `ccam import` o mostram como ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``.

## Relacionados

- [Iniciar o Claude Code](/features/launch)
- [Gerenciar contas](/features/accounts)
- [Referência de comandos](/docs/commands#shortcuts)
- [Instalar o ccam](/docs/install)
