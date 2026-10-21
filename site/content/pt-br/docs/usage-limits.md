---
title: Limites de uso e troca de contas
description: Veja o uso de 5 horas e o uso semanal de cada conta do Claude Code com ccam usage e abra a conta com mais espaço com ccam next --auto.
updated: 2026-10-21
order: 5
translation: machine
sourceHash: "fadb998482da"
---

`ccam usage` mostra quanto cada conta já usou do limite de 5 horas, do limite semanal e do limite semanal próprio de cada modelo (como o Fable). `ccam next --auto` abre a conta com mais espaço, e `--continue` leva junto a conversa atual.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Cada conta ganha uma linha com a idade dos seus números e, abaixo, uma barra por limite com o tempo até reiniciar: a janela de 5 horas, a janela semanal de todos os modelos e a janela semanal própria de cada modelo, como o Fable. A partir de 80% a barra fica amarela, e em 100%, vermelha. A última linha mostra a conta que `ccam next --auto` abriria.

## De onde vêm os números

O ccam consulta o endpoint de uso da Anthropic para todas as contas ao mesmo tempo, usando o login de cada conta. Ele espera no máximo 3 segundos por conta e mostra cada conta assim que ela responde.

Esse endpoint não é documentado publicamente e pode mudar. Quando uma solicitação falha, o ccam mostra os últimos números conhecidos e os marca como cache. Uma janela que reiniciou depois que esses números foram salvos mostra 0%. `ccam doctor` avisa se o formato da resposta mudou.

Números obtidos há menos de 60 segundos são reaproveitados. Opções:

| Opção | Efeito |
| --- | --- |
| `--watch [seconds]` | Atualiza continuamente, a cada 60 segundos por padrão |
| `--json` | Saída legível por máquina |
| `-v` | Também mostra a origem dos números de cada conta |
| `--no-cache` | Sempre consulta de novo |
| `--source cache` | Nunca consulta; mostra apenas os números salvos |

## Logins expirados

Quando o token de acesso de uma conta expira, o ccam não o renova, porque renová-lo enquanto uma sessão do Claude Code usa o mesmo login poderia desconectar essa sessão. Inicie a conta uma vez com `ccam <name>` e o Claude Code renova o token.

`ccam usage --refresh` renova os tokens expirados por conta própria, mas apenas das contas que não estão em execução. Ele nunca renova logins armazenados no Acesso às Chaves do macOS ou em um `~/.claude` importado.

## Uso na linha de status

```sh
ccam statusline install
```

Depois disso, o Claude Code mostra o nome da conta e o uso dela na linha de status, por exemplo `[work] · Opus · 5h 24% · 7d 41%`. A partir de 90%, a linha também sugere `ccam next --auto --continue`. Se você já tem um comando de linha de status, o ccam o mantém e coloca o nome da conta na frente. A linha de status também salva os números, que o `ccam usage` usa quando o endpoint não responde.

## Troca de conta

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` pula as contas que não estão com login feito, e `--auto` pula as contas em 100% e as contas de outra organização. Quando as duas contas compartilham o histórico de conversas, `--continue` abre a mesma conversa; caso contrário, o Claude Code a copia para a nova conta. Consulte [compartilhamento](/docs/sharing) e [trocar para a próxima conta](/features/next#continue).

Para scripts e jobs agendados que executam `claude -p`, coloque `--fallback` antes do nome da conta. Quando essa conta atinge o limite, o ccam executa o mesmo comando de novo com a conta que tem mais espaço, e só a saída da execução que terminou é impressa:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Consulte [Executar de novo no limite](/features/next#run-again-at-the-limit) para ver como o ccam decide.

O [guia do fluxo de trabalho para limites de uso](/guides/usage-limit-workflow) mostra um exemplo completo.
