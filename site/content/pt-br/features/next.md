---
title: Trocar para a próxima conta
description: ccam next abre o Claude Code com outra conta; --auto escolhe a de mais espaço, --continue leva a conversa junto e --fallback reexecuta jobs -p.
updated: 2026-10-24
order: 8
translation: machine
sourceHash: "a38daed5dd97"
---

`ccam next` abre o Claude Code com outra conta sem que você precise dizer qual. Use-o quando uma conta atingir o limite de uso: o `ccam next` simples pega a conta seguinte à última que você usou, `--auto` pega a que tem mais espaço e `--continue` leva a última conversa desta pasta para a nova conta, compartilhem as duas o histórico ou não. Dentro do Claude Code, o [plugin do ccam](#inside-claude-code) avisa quando um limite interrompe uma sessão, e `/ccam:switch` passa a conversa adiante. Para scripts que executam `claude -p`, `ccam --fallback` faz a troca sozinho (consulte [Executar de novo no limite](#run-again-at-the-limit)).

**Assentos dos planos Team e Enterprise.** A escolha automática (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel) também alterna entre assentos da mesma organização Claude Team ou Enterprise. A central de ajuda da Anthropic diz "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada assento de um plano Team é atribuído a uma pessoa específica, e só essa pessoa pode usá-lo) e "Users can't route their work through another user's account to use that user's remaining usage" (os usuários não podem encaminhar o próprio trabalho pela conta de outro usuário para aproveitar o uso restante desse usuário) ([planos Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); o mesmo vale para os [planos Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por isso, usar a escolha automática entre assentos pode violar os termos da Anthropic. Consulte a sua organização e a Anthropic antes de contar com ela. Os planos pessoais (Pro, Max) são tratados na [nota sobre os termos nas perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

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
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work CCAM_PROFILE=work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

O ccam imprime a troca como `→ <from> → <to>` na saída de erro padrão, com o motivo de uma escolha do `--auto`, e depois inicia o Claude Code com a nova conta, do mesmo jeito que `ccam <name>`. `<from>` é a última conta que você usou ou, com `--continue`, a conta à qual a conversa pertence. Os argumentos depois de `--` vão para o `claude` sem alterações.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--auto` | Escolhe a conta com login feito pelo uso (veja abaixo), pulando as contas em 100% | desativado |
| `--max-usage <1-100>` | Com `--auto`: prefere as contas abaixo deste uso | `auto_max_usage`, senão 100 |
| `--strategy headroom\|expiring\|balance` | Com `--auto`: o menor uso de 5 horas, o reinício de 5 horas mais próximo ou a cota semanal que ficaria sem uso primeiro | `auto_strategy`, senão `headroom` |
| `--cross-org` | Com `--auto`: também escolhe contas de outra organização (veja abaixo) | `cross_org`, senão desativado |
| `--wait[=<duration>]` | Com `--auto`: quando nenhuma conta tem espaço, espera o reinício da primeira, no máximo por este tempo. O valor vai depois de `=`: `--wait 2h` é recusado com a dica `--wait takes its value after '=': --wait=2h` | desativado; `5h` quando usado sozinho |
| `-c`, `--continue` | Leva a última conversa desta pasta para a nova conta (veja abaixo) | desativado |
| `-- <claude args>` | Argumentos passados ao `claude` depois dos do próprio ccam | nenhum |
| `--dry-run` | Respeitado: imprime `→ <from> → <to>` e a linha de comando que executaria, sem iniciar o Claude Code nem registrar uma inicialização | desativado |

`--auto` e `--continue` funcionam juntos. Com `--continue`, os argumentos de `--resume` são colocados antes dos seus.

As opções do próprio Claude Code, como `--dangerously-skip-permissions` ou `--model`, vão depois de `--`. Diferente de `ccam <name>`, que repassa tudo o que vem depois do nome, o `ccam next` lê os próprios argumentos, então sem o `--` ele para com `ccam: unknown flag: --dangerously-skip-permissions`.

## Como funciona

Cada inicialização pelo ccam grava o horário em `~/.local/state/ccam/last-used/<name>`. `ccam next` lê esses arquivos, verifica quais contas estão com login feito e escolhe.

1. `cli.candidates` monta um `next.Candidate` por conta a partir do `config.toml`: a conta, o horário da última inicialização (`launcher.LastUsed`) e se ela tem um login que funciona (`creds.Has`: um login do Claude no `.credentials.json`, ou no macOS no Acesso às Chaves, que não expirou por inteiro nem foi apagado pelo Claude Code). Uma pasta que usa uma chave de API em vez disso (um login do Console, ou uma conta adicionada com `ccam add --api`) não tem limites de plano e nunca é escolhida.
2. Sem `--auto`, `next.Rotate` ordena as contas pelo número do slot e encontra a iniciada mais recentemente. Essa é a conta da qual você está saindo. Se nenhuma conta foi iniciada ainda, é o menor slot.
3. A partir dela, avança na ordem dos slots, voltando do último slot para o primeiro, e retorna a primeira conta com login feito.
4. `launcher.LaunchWith` prepara e inicia a conta escolhida exatamente como `ccam <name>` faria: verifica os links da conta, aplica o compartilhamento de MCP, registra a sessão e substitui o ccam pelo `claude` com `syscall.Exec`.

A busca começa na última conta usada e para na primeira conta com login feito depois dela:

<!-- motion:next-rotate -->

`work` foi a iniciada mais recentemente, então a busca começa depois dela. `side` não tem login e é pulada, e `lab` é escolhida. Se `lab` também estivesse sem login, a busca voltaria ao início, em `personal` (`% len` no diagrama). A rotação nunca retorna a conta de onde você partiu: sem outra conta com login feito, ela para com um erro.

### --auto

`--auto` decide pelo uso em vez da ordem dos slots. `cli.pickAuto` mantém as contas com login feito, consulta o uso de todas ao mesmo tempo com o mesmo `Fetcher.FetchAll` que o `ccam usage` usa e deixa `usage.Policy.Choose()` escolher:

<!-- motion:next-auto -->

Uma conta é avaliada pelas janelas de 5 horas e semanal e, quando `--model` está entre os argumentos para o Claude Code (depois de `--`), também pela janela semanal desse modelo (`opus`, `sonnet`, `haiku` ou `fable` no nome). `Policy.Choose()` divide as contas em três grupos:

1. **Com espaço**: todas as janelas abaixo de `auto_max_usage` (`config.toml`, padrão 100). A estratégia escolhe entre elas: `headroom` (padrão) pega o menor uso de 5 horas; `expiring` pega a janela de 5 horas que reinicia antes, para que a cota restante seja usada antes de se perder; `balance` pega a conta com mais cota semanal restante por dia até a semana dela reiniciar, contando cada sessão do ccam em execução contra ela, para que a cota semanal não fique sem uso e as sessões em paralelo se distribuam.
2. **Acima do limite**: abaixo de 100%, mas em `auto_max_usage` ou acima. Usadas só quando nenhuma conta tem espaço, começando pelo menor uso de 5 horas.
3. **Cheias**: uma janela em 100%, ou uso de 5 horas desconhecido. Nunca escolhidas.

Uma janela conta como 100% até reiniciar quando o Claude Code recusou uma solicitação ali: o hook do [plugin do ccam](#inside-claude-code) registra isso (`~/.local/state/ccam/usage/<name>.limit.json`), não importa o que digam os números salvos. Assim, uma conta que acabou de atingir o limite não é escolhida de novo antes do reinício, mesmo quando os números dela são antigos.

As configurações para todas as contas ficam no início do `config.toml`, acima do primeiro `[[group]]`, `[[profile]]` ou `[[folder]]`. Uma configuração abaixo de um deles pertence àquela tabela: o ccam a aponta a cada inicialização, e os comandos que salvam a configuração se recusam até que ela seja movida.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizações.** As escolhas automáticas ficam entre contas como a que você deixa: a mesma organização Team ou Enterprise, ou os planos pessoais (Pro, Max) juntos, para que o trabalho de uma conta da empresa não continue em uma conta particular, nem o contrário. O ccam sabe a organização de cada conta pelo login dela (`ccam login` a registra; logins mais antigos são lidos do `.claude.json`); uma conta cuja organização ele não consegue identificar serve em qualquer caso. Duas pastas com login na mesma conta Claude compartilham os limites dela, então uma nunca substitui a outra. `--cross-org`, ou `cross_org = true` no `config.toml`, suspende a regra da organização.

Quando nada pode ser escolhido, o ccam sai com código 75 (`EX_TEMPFAIL`), para que um script consiga distinguir "tente mais tarde" de uma falha. Com `--wait`, ele espera até a janela que bloqueia a primeira conta cheia reiniciar (30 segundos depois do horário de reinício, quando o endpoint mostra a nova janela), imprime `ccam: no account has room; waiting until 15:10` e escolhe de novo; quando nenhuma conta informa quando reinicia, ele verifica de novo depois de 10 minutos. Se nada se liberar antes do fim da espera, sai com 75. A consulta segue as regras do `ccam usage`: números obtidos há menos de 60 segundos são reaproveitados, cada conta tem 3 segundos, uma conta que não responde, ou cujo token de acesso expirou, é avaliada pelos números salvos, e um registro de limite fica por cima deles. Uma janela salva cujo horário de reinício já passou conta como 0%, então uma conta cuja janela reiniciou nesse meio-tempo não é deixada de lado. `--auto` nunca renova tokens. Diferente da rotação simples, `--auto` pode escolher a última conta usada se ela ainda tiver mais espaço.

### --continue

1. O ccam procura a última conversa desta pasta em todas as contas: dos arquivos `<session>.jsonl` em `projects/<folder>/` (o nome da pasta é a pasta de trabalho, com os links resolvidos e todo caractere que não seja letra ou dígito trocado por `-`, como o Claude Code a nomeia), aquele com o prompt ou a resposta mais recente. Sessões sem nenhum dos dois, como uma sessão que só executou `/mcp`, `/usage` ou `/status`, ou um arquivo vazio, são puladas. Se uma sessão deixou uma nota (um limite de uso ou `/ccam:switch`, veja abaixo) e nenhuma conversa tem um prompt ou uma resposta mais de 2 minutos depois dela, é essa sessão exata que é usada; as linhas que o Claude Code grava no `/exit` não contam.
2. A conta à qual ela pertence é a conta da qual você sai, e a nova conta é escolhida como acima, dentro da mesma organização. Quando várias contas compartilham esse histórico, ela pertence àquela cujo `.claude.json` diz que executou a conversa por último (`lastSessionId`); caso contrário, à usada por último.
3. Se a nova conta compartilha esse histórico (a mesma pasta de grupo), o ccam inicia `claude --resume <session id>`, exatamente aquela conversa (`--continue` abriria a mais recente da pasta, qualquer que fosse).
4. Caso contrário, ele inicia `claude --resume <the conversation's file> --fork-session`: o Claude Code copia a conversa inteira, com um novo id de sessão, para o histórico da própria nova conta, e a original fica como estava. A partir daí, o `--continue` da nova conta a encontra.

O ccam informa o tamanho da conversa: a nova conta lê tudo de novo na primeira resposta. Os caches de prompt nunca são compartilhados entre organizações, e cada plano pessoal é uma organização própria, então, quando as duas contas pertencem a organizações diferentes, o ccam acrescenta `prompt caches are not shared between organizations, so the first reply reads it all again`. Quando a conta antiga tem login em servidores MCP em que a nova não tem, o ccam os indica: `work is not signed in to MCP servers personal is: linear (sign in with /mcp there; sign-ins are not shared)`. Compartilhar `history` em um grupo evita a cópia, mantém uma só conversa e também compartilha os checkpoints a partir dos quais o `/rewind` restaura o código (`file-history/`).

## Dentro do Claude Code

O [plugin do ccam](/docs/install#claude-code-plugin) adiciona dois hooks (Claude Code 2.1.78 ou posterior):

- Quando um turno para em um limite de uso, ele anota a conversa e mostra uma notificação na área de trabalho nos terminais compatíveis (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): ``ccam: work hit its 5-hour limit (resets 3:45pm). `ccam next --continue` carries on as another account.`` Ele conhece todas as mensagens que o Claude Code trata como limite de uso: os limites de 5 horas, semanal, simples e de modelo (`You've hit your session limit`, `You've hit your limit`, `You've reached your …` para a cota de um modelo), créditos de uso e limites de gasto, e uma organização que bloqueou a conta. Quando o Claude Code gravou a solicitação recusada na conversa, o hook também registra a janela cheia até o reinício ([o registro de limite](/features/usage#the-limit-record)), exceto no bloqueio de uma organização, que os administradores dela podem suspender a qualquer momento. Uma limitação breve do servidor ou o limite de taxa de uma chave de API não é um limite de uso e é ignorado, assim como a parada de um subagente e as sessões sem ninguém diante delas (`-p`, sessões em segundo plano). O Claude Code 2.1.234 e posteriores também esperam na sessão e continuam sozinhos depois de um reinício de 5 horas; a notificação dá a você a opção de trocar de conta em vez disso.
- `/ccam:switch [account]` anota a conversa e encerra o turno antes que qualquer solicitação seja feita: digite `/exit`, e então `ccam next --continue` a continua, com a conta que você indicou ou com a próxima. Em uma sessão que ainda não tem conversa, ele avisa que não há nada para continuar.

A nota fica em `~/.local/state/ccam/handoff/<account>/<session>.json`, um arquivo por sessão, para que duas sessões de uma conta nunca sobrescrevam a nota uma da outra. Ela é mantida por 12 horas e usada uma vez, quando a próxima sessão está para começar; uma inicialização que falha a deixa para a próxima tentativa. Ela pertence à pasta em que a sessão começou, mesmo depois de um `cd` no shell do Claude Code. `ccam rename` atualiza os nomes nela; se a conta que ela indica foi removida desde então, o ccam avisa e escolhe como de costume.

### Relay

Com `relay = true` no início do `config.toml`, o ccam fica ao lado das sessões interativas em vez de se substituir pelo Claude Code e, depois de `/ccam:switch`, leva a conversa daquela sessão adiante sozinho (uma nota deixada por outra sessão não é usada): o hook encerra o Claude Code, e o ccam executa `ccam next --continue` com a nota daquela sessão. Sem uma conta indicada, ele pega a que tem mais espaço, contando o modelo da conversa. No Windows, o hook não consegue encerrar o Claude Code: digite `/exit`, e o ccam continua. Isso custa um processo do ccam esperando por sessão, e é por isso que vem desativado por padrão.

A sessão que dá continuidade mantém o que retomar uma conversa não restaura: opções como `--add-dir`, `--mcp-config`, as listas de ferramentas, `--settings`, `--plugin-dir` e as opções de prompt de sistema, o modo de permissão da conversa e o nível de esforço da sessão. O próprio Claude Code restaura `--model`, `--name` e `--agent` a partir da conversa. Qualquer outra coisa (`--worktree`, `--debug`, `--remote-control`, opções que o ccam não conhece) é indicada em `ccam: not carried over: …`. A troca nunca faz a pergunta do histórico compartilhado, e o que você digitou enquanto ela rodava é descartado.

O relay só continua quando o Claude Code terminou do jeito que o `/ccam:switch` pediu: não quando você digitou `/exit` antes, pressionou Ctrl+C enquanto ele terminava, ou o ccam recebeu um hangup ou outro sinal, nem quando o terminal sumiu. A nota então fica, e `ccam next --continue` leva a conversa adiante depois. Uma sessão iniciada pelo relay de um ccam mais antigo, ou uma sem ninguém diante dela (`-p`, sessões em segundo plano), não é encerrada pelo hook; ela recebe a mensagem do `/exit` em vez disso. `claude remote-control` roda sem relay.

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
2. A entrada padrão vinda de um pipe é lida uma vez, até 10 MiB, e entregue a cada tentativa. Como o Claude Code, o ccam espera no máximo 3 segundos para ela começar (não com `--input-format stream-json`); um pipe que fica aberto e não envia nada deixa a execução sem entrada, com uma linha na saída de erro padrão.
3. O ccam inicia o `claude` como processo filho em vez de usar `exec`, retém a saída padrão dele e repassa a saída de erro à medida que ela chega.
4. Quando o `claude` sai com 0, ou falha por outro motivo, a saída retida e o código de saída são repassados sem alteração. O mesmo vale para uma execução encerrada por um sinal (um código de saída acima de 128): pediram que ela parasse, e ela não roda de novo.
5. Uma falha conta como limite de uso quando a mensagem do Claude Code diz isso, qualquer uma das mensagens de limite de uso que o [hook](#inside-claude-code) conhece, também dentro da saída de `--output-format json` ou `stream-json`, ou quando o uso da conta mostra a janela de 5 horas, a semanal ou a do `--model` em 100% (incluindo um registro de limite).
6. Então a saída dessa tentativa é descartada e o mesmo comando roda com outra conta logada: a que tem o menor uso de 5 horas, depois as contas cujo uso não pode ser consultado, pela ordem dos números. Contas em 100% são puladas, contas com login por chave de API nunca são tentadas, e cada conta é tentada no máximo uma vez.
7. Com `-c`, `--continue`, `-r` ou `--resume` entre os argumentos, só são tentadas contas do mesmo grupo de histórico, para que a conversa possa ser encontrada.
8. Só são tentadas contas da mesma organização, como com `--auto` (`cross_org = true` suspende isso).

A próxima conta é escolhida pelas mesmas regras do `--auto` (`auto_max_usage`, `auto_strategy`, a janela do `--model`). Quando não sobra nenhuma conta, o ccam imprime a saída da última tentativa, avisa que nenhuma conta logada tem espaço e sai com código 75. Com `--wait` antes do nome da conta (`ccam --fallback --wait=2h work -p …`), ele espera o reinício da primeira conta, depois pode tentar todas as contas de novo; ele sempre espera pelo menos uma vez antes, para que uma conta que acabou de atingir o limite não rode de novo na hora. A duração limita toda a espera da execução em conjunto, não cada espera. Escrita com espaço (`--wait 30m`), ela é recusada com a mesma dica do `ccam next`. Um `--` antes dos argumentos do Claude Code funciona como em qualquer inicialização (`ccam --fallback work -- -p …`). `ccam --dry-run --fallback work -p …` imprime o comando da primeira tentativa.

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
- `--fallback` retém a saída padrão até uma tentativa terminar, então uma saída em streaming como `--output-format stream-json` chega de uma vez no final. Um cliente que envia entrada `stream-json` em streaming e espera respostas antes de enviar mais não pode usar `--fallback`. Arquivos que a primeira tentativa alterou antes de atingir o limite continuam alterados; a tentativa seguinte começa o prompt de novo.

## Relacionados

- [Limites de uso](/features/usage)
- [Limites de uso e troca de contas](/docs/usage-limits)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Fluxo de trabalho para limites de uso](/guides/usage-limit-workflow)
