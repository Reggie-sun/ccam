---
title: Limites de uso e troca de contas
description: Veja o uso de 5 horas e o uso semanal de cada conta do Claude Code com ccs usage e abra a conta com mais espaço com ccs next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a9ee77936668
---

`ccs usage` mostra quanto de cada limite, o de 5 horas e o semanal, cada conta já usou. `ccs next --auto` abre a conta com mais espaço, e `--continue` leva junto a conversa atual.

```console
$ ccs usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccs next --auto 로 바로 열기)
```

As colunas são: conta, uso de 5 horas, tempo até reiniciar, uso semanal, tempo até reiniciar e a idade dos números.

## De onde vêm os números

O ccam consulta o endpoint de uso da Anthropic para todas as contas ao mesmo tempo, usando o login de cada conta. Ele espera no máximo 3 segundos por conta e mostra cada conta assim que ela responde.

Esse endpoint não é documentado publicamente e pode mudar. Quando uma solicitação falha, o ccam mostra os últimos números conhecidos e os marca como cache. `ccam doctor` avisa se o formato da resposta mudou.

Números obtidos há menos de 60 segundos são reaproveitados. Opções:

| Opção | Efeito |
| --- | --- |
| `--watch [seconds]` | Atualiza continuamente, a cada 60 segundos por padrão |
| `--json` | Saída legível por máquina |
| `-v` | Também mostra os limites semanais por modelo e a origem de cada número |
| `--no-cache` | Sempre consulta de novo |
| `--source cache` | Nunca consulta; mostra apenas os números salvos |

## Logins expirados

Quando o token de acesso de uma conta expira, o ccam não o renova, porque renová-lo enquanto uma sessão do Claude Code usa o mesmo login poderia desconectar essa sessão. Inicie a conta uma vez com `ccs <name>` e o Claude Code renova o token.

`ccs usage --refresh` renova os tokens expirados por conta própria, mas apenas das contas que não estão em execução. Ele nunca renova logins armazenados no Acesso às Chaves do macOS ou em um `~/.claude` importado.

## Uso na linha de status

```sh
ccam statusline install
```

Depois disso, o Claude Code mostra o nome da conta e o uso dela na linha de status, por exemplo `[work] · Opus · 5h 24% · 7d 41%`. Se você já tem um comando de linha de status, o ccam o mantém e coloca o nome da conta na frente. A linha de status também salva os números, que o `ccs usage` usa quando o endpoint não responde.

## Troca de conta

```sh
ccs next              # the account after the one you used last
ccs next --auto       # the account with the lowest 5-hour usage
ccs next --continue   # the next account that shares your conversation history
```

`ccs next` pula as contas que não estão com login feito, e `--auto` pula as contas em 100%. Para continuar uma conversa, as duas contas precisam compartilhar o histórico de conversas; consulte [compartilhamento](/docs/sharing).

O [guia do fluxo de trabalho para limites de uso](/guides/usage-limit-workflow) mostra um exemplo completo.
