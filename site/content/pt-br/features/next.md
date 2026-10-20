---
title: Trocar para a próxima conta
description: ccam next abre o Claude Code com a conta seguinte à última usada; --auto escolhe a de mais espaço, --continue retoma e --fallback reexecuta jobs -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 7bc55c2dd520
---

`ccam next` abre o Claude Code com outra conta sem que você precise dizer qual. Use-o quando uma conta atingir o limite de uso: o `ccam next` simples pega a conta seguinte à última que você usou, `--auto` pega a que tem mais espaço na janela de 5 horas e `--continue` continua a última conversa na nova conta. Para scripts que executam `claude -p`, `ccam --fallback` faz essa troca sozinho (consulte [Executar de novo no limite](#run-again-at-the-limit)).

## Uso

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

O ccam imprime a troca como `→ <from> → <to>` na saída de erro padrão e depois inicia o Claude Code com a nova conta, do mesmo jeito que `ccam <name>`. `<from>` é a última conta que você usou. Os argumentos depois de `--` vão para o `claude` sem alterações.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--auto` | Escolhe a conta com login feito que tem o menor uso de 5 horas, pulando as contas em 100% | desativado |
| `-c`, `--continue` | Considera só as contas que compartilham o histórico de conversas com a última conta usada e inicia `claude --continue` | desativado |
| `-- <claude args>` | Argumentos passados ao `claude` depois dos do próprio ccam | nenhum |
| `--dry-run` | Respeitado: imprime `→ <from> → <to>` e a linha de comando que executaria, sem iniciar o Claude Code nem registrar uma inicialização | desativado |

`--auto` e `--continue` funcionam juntos. Com `--continue`, `--continue` é colocado antes dos seus próprios argumentos.

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

`--auto` decide pelo uso em vez da ordem dos slots. `cli.pickAuto` mantém as contas com login feito, consulta o uso de todas ao mesmo tempo com o mesmo `Fetcher.FetchAll` que o `ccam usage` usa e deixa `usage.Best()` escolher:

<!-- motion:next-auto -->

`usage.Best()` pega a menor porcentagem de 5 horas e pula as contas em 100% ou com uso de 5 horas desconhecido. A consulta segue as regras do `ccam usage`: números obtidos há menos de 60 segundos são reaproveitados, cada conta tem 3 segundos e uma conta que não responde, ou cujo login expirou, é avaliada pelos números salvos. Uma janela salva cujo horário de reinício já passou conta como 0%, então uma conta cuja janela reiniciou nesse meio-tempo não é deixada de lado. `--auto` nunca renova tokens. Diferente da rotação simples, `--auto` pode escolher a última conta usada se ela ainda tiver mais espaço.

### --continue

O `--continue` do Claude Code abre a conversa mais recente da pasta atual, o que só funciona se a nova conta enxergar essa conversa. Por isso, com `--continue`, o ccam considera só as contas do mesmo grupo de `history` da última conta usada e coloca `--continue` na frente dos seus argumentos. Contas cujo histórico é `local` não compartilham com ninguém, então `--continue` a partir de uma delas sempre falha. Ative o compartilhamento antes com `ccam share <name> --history on` (consulte [compartilhamento](/docs/sharing)).

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

Quando não sobra nenhuma conta, o ccam imprime a saída da última tentativa, avisa que nenhuma conta logada tem espaço e sai com o código dessa tentativa. `ccam --dry-run --fallback work -p …` imprime o comando da primeira tentativa.

Sessões interativas não trocam de conta sozinhas. Em vez disso, com `ccam statusline install`, a linha de status acrescenta `→ ccam next --auto --continue` quando o uso de 5 horas ou o semanal chega a 90% ([Limites de uso](/features/usage)).

## Segurança e casos especiais

- Contas sem login nunca são escolhidas. Sem outra conta com login feito, o `ccam next` simples para com `no logged-in account other than <name>`.
- Com `--continue` e nenhuma conta compartilhando o histórico: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Com `--auto` e nenhuma conta com login feito: `no logged-in account to pick`. Quando todas as contas estão cheias ou com uso desconhecido: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` faz solicitações de rede mesmo com `--dry-run`. Ele só lê e nunca grava arquivos de token.
- A conta escolhida é iniciada pelo mesmo caminho que `ccam <name>`, então as mesmas verificações se aplicam, incluindo o aviso quando outro Claude Code que compartilha o mesmo histórico já está em execução nesta pasta.
- `last-used` só é gravado quando o Claude Code realmente inicia. `--dry-run` não mexe nele, então uma simulação não avança a rotação.
- `--fallback` retém a saída padrão até uma tentativa terminar, então uma saída em streaming como `--output-format stream-json` chega de uma vez no final. Arquivos que a primeira tentativa alterou antes de atingir o limite continuam alterados; a tentativa seguinte começa o prompt de novo.

## Relacionados

- [Limites de uso](/features/usage)
- [Limites de uso e troca de contas](/docs/usage-limits)
- [Compartilhar servidores MCP, configurações e histórico](/docs/sharing)
- [Fluxo de trabalho para limites de uso](/guides/usage-limit-workflow)
