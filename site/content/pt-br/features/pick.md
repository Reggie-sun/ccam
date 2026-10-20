---
title: Escolher uma conta para scripts
description: ccam pick imprime a conta com mais espaço para um script, dá contas diferentes a workers iniciados juntos e sai com 75 quando nenhuma tem espaço.
updated: 2026-10-20
order: 12
translation: machine
sourceHash: 65a5a7e179a5
---

`ccam pick` imprime o nome da conta com login feito que tem mais espaço, para um script iniciar o Claude Code com ela. Workers iniciados ao mesmo tempo recebem cada um uma conta diferente enquanto houver contas com espaço suficientes, então um lote de jobs `claude -p` se espalha pelas suas contas em vez de esgotar uma só.

## Uso

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live"
}
```

Códigos de saída: `0` com um nome na saída padrão, `75` (`EX_TEMPFAIL`) quando nenhuma conta tem espaço ou nenhuma tem login feito, `1` para outros erros, como uma opção inválida.

## Opções

| Opção | O que faz | Padrão |
| --- | --- | --- |
| `--json` | Imprime a conta, o uso de 5 horas e o semanal e de onde vêm os números | desativado |
| `--model <model>` | Também conta a janela semanal desse modelo, como `opus` | nenhum |
| `--max-usage <1-100>` | Prefere as contas abaixo deste uso | `auto_max_usage`, senão 100 |
| `--strategy headroom\|expiring` | O menor uso de 5 horas, ou o reinício de 5 horas mais próximo | `auto_strategy`, senão `headroom` |
| `--wait[=<duration>]` | Quando nenhuma conta tem espaço, espera o reinício da primeira, no máximo por este tempo | desativado; `5h` quando usado sozinho |
| `--dry-run` | Escolhe sem registrar a escolha | desativado |

## Como funciona

1. O ccam lista as contas com login feito (contas com chave de API não têm limites e ficam de fora).
2. Ele pega um lock (`~/.local/state/ccam/picks.lock`), para que escolhas feitas ao mesmo tempo aconteçam uma depois da outra.
3. Consulta o uso de todas as contas ao mesmo tempo, do mesmo jeito que o `ccam usage`, com o cache de 60 segundos.
4. Contas escolhidas nos últimos 2 minutos (`~/.local/state/ccam/picks/<name>`) ficam de fora primeiro, e as [regras de escolha](/features/next#auto) do `ccam next --auto` escolhem entre as restantes: primeiro as abaixo de `auto_max_usage`, pela estratégia, nunca uma conta em 100%.
5. Se todas as contas com espaço foram escolhidas há pouco, as regras escolhem de novo entre todas elas, então um quarto worker em três contas divide uma delas.
6. A escolha é registrada e o lock liberado, e então o nome é impresso.

A janela de 2 minutos importa porque o uso de um job novo demora um pouco para aparecer nos números: sem ela, cinco workers iniciados juntos veriam todos a mesma conta como a de mais espaço.

## Segurança e casos especiais

- `ccam pick` só lê o uso e grava o horário da escolha; ele nunca renova nem lê tokens além do que o `ccam usage` faz.
- Uma escolha não reserva a conta. Dois scripts ainda podem começar na mesma conta depois de 2 minutos; use `--fallback` para se recuperar quando essa conta atingir o limite.
- Esperar com `--wait` não segura nenhum lock, então outras escolhas continuam nesse meio-tempo.
- `ccam pick` ignora o [`ccam folder`](/features/folder): ele escolhe só pelo uso.

## Relacionados

- [Trocar para a próxima conta](/features/next)
- [Executar de novo no limite](/features/next#run-again-at-the-limit)
- [Limites de uso](/features/usage)
- [Contas com chave de API](/features/api)
