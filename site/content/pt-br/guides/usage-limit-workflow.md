---
title: O que fazer no limite de uso do Claude Code
description: Espere o limite reiniciar ou continue na hora a mesma conversa em outra conta do Claude Code. Como funcionam os limites de 5 horas e semanal, passo a passo.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "049a46684c7b"
---

Quando uma conta atinge o limite de uso (muitas vezes chamado de rate limit), você tem duas opções: esperar até o limite reiniciar ou continuar agora a mesma conversa em outra conta sua. Com o ccam, a segunda opção é um único comando, `ccam next --auto --continue`.

## Quais limites o Claude Code tem?

Nos planos Pro e Max, o Claude Code tem um [limite de sessão de cinco horas e um limite semanal](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Os dois são compartilhados com o Claude na web e nos apps, e o uso conta para os dois ao mesmo tempo. Quando um deles se esgota, o Claude Code [mostra uma mensagem](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) como "You've hit your session limit · resets 3:45pm". Trocar de modelo não ajuda, porque os limites valem para todos os modelos.

## Opção 1: esperar o limite reiniciar

Se você tem uma só conta, espere. As versões recentes do Claude Code [esperam na sessão aberta](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) e continuam a tarefa sozinhas quando o limite reinicia. `/usage` dentro do Claude Code mostra quanto você já usou dos limites do seu plano.

## Opção 2: continuar agora em outra conta

Se você tem uma segunda conta, pode seguir na hora com a mesma conversa. Configure uma vez:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Agora a linha de status mostra a conta e o uso de 5 horas dela enquanto você trabalha, então o limite não pega você de surpresa.

Quando o limite for atingido, saia do Claude Code (`/exit`) e depois:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Abra a conversa na conta com mais espaço:

```sh
ccam next --auto --continue
```

O Claude Code inicia como `work` com a conversa mais recente da pasta atual, incluindo tudo o que foi dito até agora. Compartilhar o histórico (`ccam group set default --add history`, acima) faz com que continue sendo uma única conversa; sem isso, o Claude Code copia a conversa para `work`, e o ccam informa quantos tokens essa primeira resposta lê de novo.

Com o [plugin do ccam](/docs/install#claude-code-plugin) no Claude Code, uma sessão que para no limite mostra uma notificação e registra a janela cheia, para que o `ccam next --auto` não escolha essa conta de novo antes do reinício, e `/ccam:switch` (ou `/ccam:switch work`) registra a conversa para o `ccam next --continue`, sem deixar nada para adivinhar.

**Assentos dos planos Team e Enterprise.** A escolha automática (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel) também alterna entre assentos da mesma organização Claude Team ou Enterprise. A central de ajuda da Anthropic diz "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada assento de um plano Team é atribuído a uma pessoa específica, e só essa pessoa pode usá-lo) e "Users can't route their work through another user's account to use that user's remaining usage" (os usuários não podem encaminhar o próprio trabalho pela conta de outro usuário para aproveitar o uso restante desse usuário) ([planos Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); o mesmo vale para os [planos Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por isso, usar a escolha automática entre assentos pode violar os termos da Anthropic. Consulte a sua organização e a Anthropic antes de contar com ela. Os planos pessoais (Pro, Max) são tratados na [nota sobre os termos nas perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## O ccam troca de conta automaticamente?

Não. O ccam não fica em execução enquanto você trabalha: ele inicia o Claude Code e se substitui por ele. Você decide quando trocar, e `ccam next --auto` escolhe a conta com o menor uso de 5 horas e pula as contas que estão em 100%, sem login ou de outra organização. Com `relay = true` no `config.toml`, o ccam fica, em vez disso, esperando ao lado da sessão, então `/ccam:switch` continua sozinho, mantendo as opções, o modo de permissão e o nível de esforço da sessão. Algumas outras ferramentas trocam o login de uma sessão em andamento; [o ccam e outras ferramentas para trocar de conta](/compare/account-switchers) explica a diferença.

## O que ter em mente

- Execute `ccam next --continue` na mesma pasta em que você estava trabalhando; as conversas são guardadas por pasta de projeto.
- Evite executar duas contas na mesma conversa ao mesmo tempo. O ccam avisa se você tentar.
- Sem histórico compartilhado, a conversa é copiada para a outra conta e as duas cópias seguem caminhos separados; os checkpoints a partir dos quais o `/rewind` restaura o código ficam com a primeira conta.
- A primeira mensagem na nova conta pode consumir mais do limite dela do que o normal, porque essa conta ainda não tem a conversa em cache.
- Se você pode trocar de conta desse jeito depende dos termos de cada conta.

Veja [limites de uso](/docs/usage-limits) para todas as opções de `ccam usage` e `ccam next`.
