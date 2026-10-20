---
title: Trocar para a próxima conta
description: ccam next abre o Claude Code com outra conta; --auto escolhe a de mais espaço, --continue leva a conversa junto e --fallback reexecuta jobs -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

`ccam next` abre o Claude Code com outra conta sem que você precise dizer qual. Use-o quando uma conta atingir o limite de uso: o `ccam next` simples pega a conta seguinte à última que você usou, `--auto` pega a que tem mais espaço e `--continue` leva a última conversa desta pasta para a nova conta, compartilhem as duas o histórico ou não. Dentro do Claude Code, o [plugin do ccam](#inside-claude-code) avisa quando um limite interrompe uma sessão, e `/ccam:switch` passa a conversa adiante. Para scripts que executam `claude -p`, `ccam --fallback` faz a troca sozinho (consulte [Executar de novo no limite](#run-again-at-the-limit)).

## Uso

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

O ccam imprime a troca como `→ <from> → <to>` na saída de erro padrão, com o motivo de uma escolha do `--auto`, e depois inicia o Claude Code com a nova conta, do mesmo jeito que `ccam <name>`. `<from>` é a última conta que você usou ou, com `--continue`, a conta à qual a conversa pertence. Os argumentos depois de `--` vão para o `claude` sem alterações.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--auto` | Escolhe a conta com login feito pelo uso (veja abaixo), pulando as contas em 100% | desativado |
| `--max-usage <1-100>` | Com `--auto`: prefere as contas abaixo deste uso | `auto_max_usage`, senão 100 |
| `--strategy headroom\|expiring\|balance` | Com `--auto`: o menor uso de 5 horas, o reinício de 5 horas mais próximo ou a cota semanal que ficaria sem uso primeiro | `auto_strategy`, senão `headroom` |
| `--cross-org` | Com `--auto`: também escolhe contas de outra organização (veja abaixo) | `cross_org`, senão desativado |
| `--wait[=<duration>]` | Com `--auto`: quando nenhuma conta tem espaço, espera o reinício da primeira, no máximo por este tempo | desativado; `5h` quando usado sozinho |
| `-c`, `--continue` | Leva a última conversa desta pasta para a nova conta (veja abaixo) | desativado |
| `-- <claude args>` | Argumentos passados ao `claude` depois dos do próprio ccam | nenhum |
| `--dry-run` | Respeitado: imprime `→ <from> → <to>` e a linha de comando que executaria, sem iniciar o Claude Code nem registrar uma inicialização | desativado |

`--auto` e `--continue` funcionam juntos. Com `--continue`, os argumentos de `--resume` são colocados antes dos seus.

## Como funciona

Cada inicialização pelo ccam grava o horário em `~/.local/state/ccam/last-used/<name>`. `ccam next` lê esses arquivos, verifica quais contas estão com login feito e escolhe.

1. `cli.candidates` monta um `next.Candidate` por conta a partir do `config.toml`: a conta, o horário da última inicialização (`launcher.LastUsed`) e se ela tem login (`creds.Has`: um `.credentials.json` legível ou, no macOS, um item do Acesso às Chaves).
2. Sem `--auto`, `next.Rotate` ordena as contas pelo número do slot e encontra a iniciada mais recentemente. Essa é a conta da qual você está saindo. Se nenhuma conta foi iniciada ainda, é o menor slot.
3. A partir dela, avança na ordem dos slots, voltando do último slot para o primeiro, e retorna a primeira conta com login feito.
4. `launcher.LaunchWith` prepara e inicia a conta escolhida exatamente como `ccam <name>` faria: verifica os links da conta, aplica o compartilhamento de MCP, registra a sessão e substitui o ccam pelo `claude` com `syscall.Exec`.

A busca começa na última conta usada e para na primeira conta com login feito depois dela:

<!-- motion:next-rotate -->

`work` foi a iniciada mais recentemente, então a busca começa depois dela. `side` não tem login e é pulada, e `lab` é escolhida. Se `lab` também estivesse sem login, a busca voltaria ao início, em `personal` (`% len` no diagrama). A rotação nunca retorna a conta de onde você partiu: sem outra conta com login feito, ela para com um erro.

### --auto

`--auto` decide pelo uso em vez da ordem dos slots. `cli.pickAuto` mantém as contas com login feito, consulta o uso de todas ao mesmo tempo com o mesmo `Fetcher.FetchAll` que o `ccam usage` usa e deixa `usage.Policy.Choose()` escolher:

<!-- motion:next-auto -->

Uma conta é avaliada pelas janelas de 5 horas e semanal e, quando `--model` está entre os argumentos para o Claude Code, também pela janela semanal desse modelo (`opus`, `sonnet`, `haiku` ou `fable` no nome). `Policy.Choose()` divide as contas em três grupos:

1. **Com espaço**: todas as janelas abaixo de `auto_max_usage` (`config.toml`, padrão 100). A estratégia escolhe entre elas: `headroom` (padrão) pega o menor uso de 5 horas; `expiring` pega a janela de 5 horas que reinicia antes, para que a cota restante seja usada antes de se perder; `balance` pega a conta com mais cota semanal restante por dia até a semana dela reiniciar, contando cada sessão do ccam em execução contra ela, para que a cota semanal não fique sem uso e as sessões em paralelo se distribuam.
2. **Acima do limite**: abaixo de 100%, mas em `auto_max_usage` ou acima. Usadas só quando nenhuma conta tem espaço, começando pelo menor uso de 5 horas.
3. **Cheias**: uma janela em 100%, ou uso de 5 horas desconhecido. Nunca escolhidas.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizações.** As escolhas automáticas ficam entre contas como a que você deixa: a mesma organização Team ou Enterprise, ou os planos pessoais (Pro, Max) juntos, para que o trabalho de uma conta da empresa não continue em uma conta particular, nem o contrário. O ccam sabe a organização de cada conta pelo login dela (`ccam login` a registra; logins mais antigos são lidos do `.claude.json`); uma conta cuja organização ele não consegue identificar serve em qualquer caso. Duas pastas com login na mesma conta Claude compartilham os limites dela, então uma nunca substitui a outra. `--cross-org`, ou `cross_org = true` no `config.toml`, suspende a regra da organização.

Quando nada pode ser escolhido, o ccam sai com código 75 (`EX_TEMPFAIL`), para que um script consiga distinguir "tente mais tarde" de uma falha. Com `--wait`, ele espera até a janela que bloqueia a primeira conta cheia reiniciar (30 segundos depois do horário de reinício, quando o endpoint mostra a nova janela), imprime `ccam: no account has room; waiting until 15:10` e escolhe de novo; quando nenhuma conta informa quando reinicia, ele verifica de novo depois de 10 minutos. Se nada se liberar antes do fim da espera, sai com 75. A consulta segue as regras do `ccam usage`: números obtidos há menos de 60 segundos são reaproveitados, cada conta tem 3 segundos e uma conta que não responde, ou cujo login expirou, é avaliada pelos números salvos. Uma janela salva cujo horário de reinício já passou conta como 0%, então uma conta cuja janela reiniciou nesse meio-tempo não é deixada de lado. `--auto` nunca renova tokens. Diferente da rotação simples, `--auto` pode escolher a última conta usada se ela ainda tiver mais espaço.

### --continue

1. O ccam procura a última conversa desta pasta em todas as contas: o `<session>.jsonl` mais recente em `projects/<folder>/` (o nome da pasta é a pasta de trabalho com todo caractere que não seja letra ou dígito trocado por `-`, como o Claude Code a nomeia). Se a sessão deixou uma nota (um limite de uso ou `/ccam:switch`, veja abaixo) e nada mais recente foi gravado desde então, é essa sessão exata que é usada.
2. A conta à qual ela pertence é a conta da qual você sai, e a nova conta é escolhida como acima, dentro da mesma organização.
3. Se a nova conta compartilha esse histórico (a mesma pasta de grupo), o ccam inicia `claude --resume <session id>`, exatamente aquela conversa (`--continue` abriria a mais recente da pasta, qualquer que fosse).
4. Caso contrário, ele inicia `claude --resume <the conversation's file> --fork-session`: o Claude Code copia a conversa inteira, com um novo id de sessão, para o histórico da própria nova conta, e a original fica como estava. A partir daí, o `--continue` da nova conta a encontra.

O ccam informa o tamanho da conversa: a nova conta lê tudo de novo na primeira resposta. Os caches de prompt nunca são compartilhados entre organizações, então com `--cross-org` essa primeira resposta custa o tamanho completo. Compartilhar `history` em um grupo evita a cópia, mantém uma só conversa e também compartilha os checkpoints a partir dos quais o `/rewind` restaura o código (`file-history/`).

## Dentro do Claude Code

O [plugin do ccam](/docs/install#claude-code-plugin) adiciona dois hooks (Claude Code 2.1.78 ou posterior):

- Quando um turno para em um limite de uso (`You've hit your session limit`, o limite semanal ou o de um modelo), ele anota a conversa e mostra uma notificação na área de trabalho nos terminais compatíveis (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Uma limitação breve do servidor ou o limite de taxa de uma chave de API não é um limite de uso e é ignorado. O Claude Code 2.1.234 e posteriores também esperam na sessão e continuam sozinhos depois de um reinício de 5 horas; a notificação dá a você a opção de trocar de conta em vez disso.
- `/ccam:switch [account]` anota a conversa e encerra o turno antes que qualquer solicitação seja feita: digite `/exit`, e então `ccam next --continue` a continua, com a conta que você indicou ou com a próxima.

A nota fica em `~/.local/state/ccam/handoff/<account>.json` por 12 horas e é usada uma vez, quando a próxima sessão está para começar; uma inicialização que falha a deixa para a próxima tentativa. Ela pertence à pasta em que a sessão começou, mesmo depois de um `cd` no shell do Claude Code. Se a conta que ela indica foi renomeada ou removida desde então, o ccam avisa e escolhe como de costume. Com `relay = true` no `config.toml`, o ccam fica ao lado das sessões interativas em vez de se substituir pelo Claude Code e, depois de `/ccam:switch`, leva a conversa daquela sessão adiante sozinho (uma nota deixada por outra sessão não é usada). Isso custa um processo do ccam esperando por sessão, e é por isso que vem desativado por padrão. A sessão que dá continuidade começa com `ccam next --continue`, sem os argumentos com que a primeira foi iniciada, como `--model`.

## Executar de novo no limite

`ccam next` serve para uma sessão interativa que parou. Uma execução não interativa, `claude -p` em um script ou em um job agendado, não pode esperar por você, então `--fallback` antes do nome da conta deixa o ccam fazer a troca por ela:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` só funciona junto com `-p` ou `--print`; caso contrário, o ccam para e indica `ccam next --auto --continue`.
2. A entrada padrão é lida uma vez e entregue a cada tentativa.
3. O ccam inicia o `claude` como processo filho em vez de usar `exec`, retém a saída padrão dele e repassa a saída de erro à medida que ela chega.
4. Quando o `claude` sai com 0, ou falha por outro motivo, a saída retida e o código de saída são repassados sem alteração.
5. Uma falha conta como limite de uso quando a mensagem do Claude Code diz isso (`You've hit your limit`, `usage limit reached`), ou quando a consulta de uso da conta mostra a janela de 5 horas ou a semanal em 100%.
6. Então a saída dessa tentativa é descartada e o mesmo comando roda com outra conta logada: a que tem o menor uso de 5 horas, depois as contas cujo uso não pode ser consultado, pela ordem dos números. Contas em 100% são puladas, e cada conta é tentada no máximo uma vez.
7. Com `-c`, `--continue`, `-r` ou `--resume` entre os argumentos, só são tentadas contas do mesmo grupo de histórico, para que a conversa possa ser encontrada.
8. Só são tentadas contas da mesma organização, como com `--auto` (`cross_org = true` suspende isso).

A próxima conta é escolhida pelas mesmas regras do `--auto` (`auto_max_usage`, `auto_strategy`, a janela do `--model`). Quando não sobra nenhuma conta, o ccam imprime a saída da última tentativa, avisa que nenhuma conta logada tem espaço e sai com código 75. Com `--wait` antes do nome da conta (`ccam --fallback --wait=2h work -p …`), ele espera o reinício da primeira conta, depois pode tentar todas as contas de novo; ele sempre espera pelo menos uma vez antes, para que uma conta que acabou de atingir o limite não rode de novo na hora. `ccam --dry-run --fallback work -p …` imprime o comando da primeira tentativa.

Sessões interativas não trocam de conta sozinhas. Em vez disso, com `ccam statusline install`, a linha de status acrescenta `→ ccam next --auto --continue` quando o uso de 5 horas ou o semanal chega a 90% ([Limites de uso](/features/usage)).

## Segurança e casos especiais

- Contas sem login nunca são escolhidas. Sem outra conta com login feito, o `ccam next` simples para com `no logged-in account other than <name>`.
- Com `--continue` e nenhuma conversa na pasta: `no account has a conversation in <folder> to continue`.
- Quando a regra da organização não deixa nenhuma conta: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Com `--auto` e nenhuma conta com login feito: `no logged-in account to pick`. Quando todas as contas estão cheias ou com uso desconhecido: `` no account with known usage and headroom (see `ccam usage`) ``. Os dois casos terminam com código de saída 75.
- `--max-usage`, `--strategy`, `--wait` e `--cross-org` só funcionam com `--auto`, e um valor fora de 1-100 ou outra estratégia é recusado.
- Um fork copia a conversa, não o restante do estado da conta: os checkpoints do `file-history/` da conta original não vêm junto, a menos que o histórico seja compartilhado.
- `--auto` faz solicitações de rede mesmo com `--dry-run`. Ele só lê e nunca grava arquivos de token.
- A conta escolhida é iniciada pelo mesmo caminho que `ccam <name>`, então as mesmas verificações se aplicam, incluindo o aviso quando outro Claude Code que compartilha o mesmo histórico já está em execução nesta pasta.
- `last-used` só é gravado quando o Claude Code realmente inicia. `--dry-run` não mexe nele, então uma simulação não avança a rotação.
- `--fallback` retém a saída padrão até uma tentativa terminar, então uma saída em streaming como `--output-format stream-json` chega de uma vez no final. Arquivos que a primeira tentativa alterou antes de atingir o limite continuam alterados; a tentativa seguinte começa o prompt de novo.

## Relacionados

- [Limites de uso](/features/usage)
- [Limites de uso e troca de contas](/docs/usage-limits)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Fluxo de trabalho para limites de uso](/guides/usage-limit-workflow)
