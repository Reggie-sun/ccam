---
title: Iniciar o Claude Code
description: Inicie o Claude Code com qualquer conta usando ccam <name|number>, ccam run ou o seletor do ccam; os argumentos vão para o claude sem alterações, e --dry-run mostra o comando.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d033da86f44b
---

`ccam work` inicia o Claude Code com a conta `work`, e tudo o que vem depois do nome vai para o `claude` como está. Use-o onde você digitaria `claude`: em um terminal, em um script ou com uma entrada enviada por pipe.

## Uso

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Informe o nome de uma conta ou o número dela em `ccam list`. O ccam não lê os argumentos depois do nome, então `ccam work --help` mostra a ajuda do Claude Code. Um `--` depois do nome é opcional. A entrada padrão, o código de saída e os sinais se comportam como no `claude`.

Quando o primeiro argumento começa com `-` e não é uma das opções do próprio ccam, o ccam inicia a conta padrão com todos os argumentos. A padrão é a conta definida com `ccam use` ou, se você não definiu nenhuma, a conta com o menor número.

### O seletor

`ccam` sem argumentos abre o seletor quando a entrada e a saída padrão são um terminal e, nos outros casos, imprime `ccam list`. Ele mostra o número, o e-mail, o uso de 5 horas, o último uso e os itens compartilhados de cada conta; o uso é preenchido à medida que cada conta responde.

| Tecla | O que faz |
| --- | --- |
| `Enter` | Inicia o Claude Code com a conta selecionada |
| `1` a `9` | Inicia a conta com esse número |
| `Up`/`Down`, `k`/`j` | Move a seleção |
| `a` | Fecha o seletor e executa o assistente do `ccam add` |
| `s` | Edita os itens compartilhados da conta selecionada (`Space` ou `x` alterna, `Enter` aplica, `Esc` volta) |
| `d` | Torna a conta selecionada a padrão |
| `r` | Consulta o uso de novo, sem o cache curto |
| `q`, `Esc`, `Ctrl+C` | Sai |

O seletor fecha antes de agir, então o terminal é restaurado antes de o Claude Code iniciar. Depois de `d` ou de uma alteração de compartilhamento, ele abre de novo. Uma alteração de compartilhamento mostra o plano e pergunta antes, como o `ccam share` faz.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--dry-run` | Imprime o ambiente e o comando que o ccam executaria e sai | desativado |
| `-v`, `--version` | Imprime `ccam <version>` (só como primeiro argumento) | |
| `-h`, `--help` | A ajuda do ccam, como `ccam --help` ou `ccam run --help` | |

Regras dos argumentos:

- O ccam só lê as próprias opções antes do nome da conta. Depois do nome, todos os argumentos pertencem ao Claude Code, incluindo `--dry-run` e `--help`.
- `--dry-run` funciona como `ccam --dry-run work`, `ccam run --dry-run work` e `ccam --dry-run` (a conta padrão). Ela é respeitada, não recusada.
- `ccam run <name|number> [claude args]` é a forma longa de `ccam <name>`. `ccam run` sem nome usa a conta padrão.
- Os nomes de conta não podem ser nomes de comandos do ccam, como `add`, `list`, `run` ou `usage`, então `ccam <word>` nunca é ambíguo.

A simulação imprime uma linha que você poderia colar em um shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Como funciona

Cada inicialização segue as mesmas etapas, da linha de comando até o `exec`:

<!-- motion:launch -->

Cada caixa acende quando a etapa dela é executada, e as linhas abaixo são o ambiente e os argumentos com que o Claude Code inicia.

1. `cmd/ccam/main.go` chama `app.Main`. Se o ccam foi iniciado com outro nome, é um [atalho](/features/shortcuts) e vai direto para o inicializador.
2. `app.parse` examina `os.Args` antes de o cobra entrar em ação. Só o `ccam` sozinho (o seletor), um nome de comando do ccam (`config.Reserved`), `-h`/`--help` e `-v`/`--version` vão para `cli.Execute`. Todo o resto é uma inicialização.
3. `launcher.Launch` lê `~/.config/ccam/config.toml`, o único arquivo que ele analisa.
4. `launcher.Prepare` encontra a conta pelo nome ou número, descarta um `--` inicial e localiza o executável `claude` com `ResolveClaude`. A busca no PATH fica em cache em `~/.local/state/ccam/claude-path`, junto com o horário de modificação do arquivo, para que o ccam não percorra o `PATH` a cada inicialização. Um arquivo alterado ou ausente provoca uma nova busca. Um caminho `claude_bin` no `config.toml` é usado como está. O ccam recusa um `claude` que seja o próprio ccam, comparado pelo inode, porque isso entraria em loop.
5. `checkLinks` compara a pasta da conta `~/.local/share/ccam/profiles/<name>/` com as configurações de compartilhamento dela usando `share.Check`, que só chama `lstat` e `readlink`, uma ou duas chamadas de sistema por item. Um link ausente é criado na hora. Qualquer coisa que possa perder dados interrompe a inicialização.
6. `applyMCP` adiciona os servidores MCP compartilhados (veja abaixo).
7. `session.Register` registra a sessão. O ccam grava um marcador com o nome do PID dele, contendo o horário de início do processo, em `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (ou em `~/.local/state/ccam/locks/` sem `XDG_RUNTIME_DIR`). Para uma conta que compartilha o histórico de conversas, `historyCheck` primeiro procura uma sessão em execução em `history-<group>-<hash>`, em que o hash é o da pasta atual, e também se registra ali. O `exec` mantém o PID, então o marcador passa a indicar o próprio processo do Claude Code. É um marcador, e não um bloqueio mantido, porque os processos iniciados pelo Claude Code herdariam o bloqueio e o manteriam depois do fim da sessão.
8. `Plan.Exec` define o título do terminal como `claude · <name>` (quando a saída de erro padrão é um terminal), grava o horário da inicialização em `~/.local/state/ccam/last-used/<name>` sem `fsync` e chama `syscall.Exec`.

O Claude Code é executado com `CLAUDE_CONFIG_DIR` apontando para a pasta da conta e `CCAM_PROFILE` com o nome da conta. O ccam remove qualquer valor dessas variáveis que tenha herdado. Para uma pasta `~/.claude` registrada com `ccam import` no macOS, ele também define `CLAUDE_SECURESTORAGE_CONFIG_DIR` com um valor vazio, para que a conta e o `claude` simples usem o mesmo login do Acesso às Chaves.

### Por que exec

`syscall.Exec` substitui o processo do ccam pelo `claude`, então o ccam não fica na memória, as teclas, os sinais e o código de saída são os do próprio Claude Code, e os argumentos continuam sendo um array, nunca uma string de shell. O trabalho posterior a uma sessão, como notar servidores adicionados com `claude mcp add`, acontece na próxima inicialização: o ccam compara o horário de modificação do `.claude.json` com `~/.local/state/ccam/mcp/<name>/seen` e só analisa o arquivo quando ele mudou.

### Como uma linha de comando é encaminhada

Três linhas de comando, três caminhos por `internal/app`:

<!-- motion:launch-route -->

Um atalho é reconhecido pelo `argv[0]`, `ccam <name>` pelo primeiro argumento, e os dois terminam em `syscall.Exec`. Só um comando do ccam, como `ccam usage`, chega ao cobra.

### Servidores MCP compartilhados

Com a estratégia padrão `flag`, o ccam passa os servidores MCP do grupo ao Claude Code como um argumento, colocado antes dos seus:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` só retorna um argumento se o `mcp.json` do grupo existir, e `Plan.InsertArgs` o coloca logo depois de `claude`. Suas próprias opções `--mcp-config` ou `--strict-mcp-config` continuam valendo por cima. É um único argumento, `--mcp-config=<path>`, porque o `--mcp-config` do Claude Code aceita vários valores: escrito como dois argumentos, ele também leria um prompt seguinte como arquivo de configuração.

Com `mcp_strategy = "sync"` no `config.toml`, o ccam grava os servidores no `.claude.json` da conta e pula essa etapa se a conta já estiver em execução.

### O limite de 10 ms

O ccam adiciona no máximo 10 ms (p95) antes de o Claude Code iniciar. O caminho de inicialização nunca monta a árvore de comandos do cobra, não inicia o Bubble Tea nem abre conexões de rede, e `TestPackageInitIsCheap` falha se as funções `init()` dos pacotes levarem mais de 5 ms no total. `scripts/perf-check.sh` mede `ccam --dry-run work`, `ccam work` e um atalho com `hyperfine` e falha acima de 10 ms, e `BenchmarkLaunchPath` em `internal/launcher` executa as etapas em loop. Para ver o tempo de cada etapa, defina `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Segurança e casos especiais

- Um nome ou número desconhecido para com `no account "<name>"` e indica `ccam list`. Sem nenhuma conta, o ccam pede que você adicione uma com `ccam add <name>`.
- Se um caminho compartilhado virou um arquivo comum ou aponta para o lugar errado, o ccam para com `stopped because of share link problems:`, lista os problemas e pede que você execute `ccam doctor --fix`. Ele não move dados no caminho de inicialização.
- Se outro Claude Code que compartilha o histórico de conversas com esta conta (mesmo grupo) estiver em execução na mesma pasta, o ccam avisa que o histórico pode se misturar e pergunta `Run anyway? [y/N]`. Sem um terminal, ele mostra o aviso e continua.
- `--dry-run` pula a sincronização de MCP, o aviso de servidor novo, a verificação do histórico, o marcador de sessão, o título e o horário do último uso, e não faz o exec. Ele ainda verifica os links de compartilhamento e cria um link ausente, como uma inicialização normal faria.
- Com `mcp_strategy = "sync"`, uma sincronização que falha imprime `ccam: MCP sync failed (launching anyway): …`, e o Claude Code inicia com os servidores que já tinha.

## Relacionados

- [Atalhos](/features/shortcuts)
- [Gerenciar contas](/features/accounts)
- [Referência de comandos](/docs/commands)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Como o CLAUDE_CONFIG_DIR funciona](/guides/claude-config-dir)
