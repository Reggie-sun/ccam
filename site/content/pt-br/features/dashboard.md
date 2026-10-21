---
title: O painel
description: O ccam em um terminal abre um painel em tela cheia com contas, grupos, regras de pasta, extensões e verificações; cada mudança mostra o plano antes.
updated: 2026-10-21
order: 15
translation: machine
sourceHash: "27bd153b48b2"
---

`ccam` sem argumentos abre um painel em tela cheia quando a entrada e a saída padrão são um terminal; caso contrário, imprime `ccam list`. O painel mostra todas as contas com o uso e os grupos de cada uma e faz o que os comandos fazem: iniciar, adicionar, fazer login, alterar grupos, definir regras de pasta, verificar e reparar. Iniciar uma conta fecha o painel; adicionar uma conta ou fazer login sai dele para o terminal e depois volta.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

A linha abaixo dos nomes das abas diz em uma frase para que serve a aba atual. Em um terminal mais largo, a aba Accounts também mostra o e-mail de cada conta.

O primeiro quadro é desenhado só a partir do `config.toml`, em bem menos de 100 ms; o uso aparece à medida que cada conta responde. Cada limite tem a própria barra: a janela de 5 horas, a janela semanal de todos os modelos e a janela semanal própria do Fable; `-` significa que a conta não tem essa janela. Uma porcentagem fica amarela a partir de 80% e vermelha em 100%, e os valores vindos de um cache depois de uma consulta que falhou aparecem esmaecidos.

## Abas

| Aba | O que mostra | Teclas |
| --- | --- | --- |
| Accounts | número, nome, e-mail, barras de 5 horas, semanal e do Fable, último uso, grupos; `✔` a conta padrão, `▸` a conta desta pasta, `•` em execução | abaixo |
| Groups | cada grupo, o que ele compartilha e suas contas, e as contas que não estão em nenhum grupo | `enter` editar, `n` novo, `e` renomear, `x` remover |
| Folders | as regras do [`ccam folder`](/features/folder) e a que vale aqui | `n` usar uma conta aqui, `enter` alterar, `x` remover |
| Extensions | o que o [`ccam ext`](/features/ext) mostra, cada extensão uma vez, com os grupos de onde ela vem | `r` recarregar |
| Health | o que o [`ccam doctor`](/features/doctor) relata, com uma caixa por conta (abaixo) | `f` reparar a caixa selecionada, `F` reparar tudo (`doctor --fix`), `↑`/`↓` passar de uma caixa para outra |

`←`/`→` ou `Tab` trocam de aba, `?` lista as teclas da aba atual, `r` consulta o uso de novo sem o cache curto, e `q` ou `Esc` sai.

## Teclas na aba Accounts

| Tecla | O que faz |
| --- | --- |
| `Enter` | Inicia o Claude Code como a conta selecionada |
| `1`–`9` | Vai para a conta com esse número; digite `1` `2` rapidamente para ir à 12. Só `Enter` inicia |
| `↑`/`↓`, `k`/`j` | Move o cursor |
| `n` | Inicia a conta que o [`ccam next --auto`](/features/next) escolheria |
| `o` | Inicia com opções: `--continue`, `--resume`, outros argumentos ou uma simulação |
| `u` | Detalhes do uso: as janelas de 5 horas, semanal e por modelo, com os reinícios e a origem dos números |
| `a` | Adiciona uma conta (o assistente abaixo) |
| `l` | Faz login novamente |
| `s` | Escolhe os grupos da conta |
| `d` | Torna a conta padrão |
| `f` | Usa a conta na pasta atual |
| `e` | Renomeia a conta |
| `c` | Define ou remove o comando de atalho dela |
| `x` | Remove a conta (pergunta antes, começando em No) |

## A aba Health

A aba Health executa as verificações do [`ccam doctor`](/features/doctor) quando você a abre e desenha a mesma visão: uma primeira linha com uma marca por conta e depois uma caixa por conta. Cada linha de uma caixa é uma área verificada, primeiro `Folder`, `Login` e `Share links`, marcada com `✔` OK, `ℹ` informação, `⚠` um aviso ou `✘` um erro, para você ver de imediato qual conta tem um problema e onde. Os resultados que não pertencem a nenhuma conta vão para uma caixa `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` passam de uma caixa para outra, e a caixa selecionada tem a borda na cor de destaque. `f` repara só a caixa selecionada, como o `ccam doctor --fix` faria: na caixa de uma conta, repara só essa conta, e na caixa `General`, os resultados que não pertencem a nenhuma conta (caminhos de plugins e atalhos); as demais contas são apenas verificadas de novo. Ele pergunta antes, e uma caixa sem nada para reparar só mostra `nothing to repair in <name>`. `F` repara tudo de uma vez, como o `ccam doctor --fix`, e `r` verifica de novo.

## Alterar grupos

Na aba Groups, `Enter` abre o editor: os itens à esquerda, as contas à direita, `Tab` para alternar entre eles e `Space` para marcar ou desmarcar. Um item que uma das contas já recebe de outro grupo aparece esmaecido, com o motivo, assim como uma conta que receberia um item duas vezes; um item só pode vir de um grupo. `Enter` mostra o plano:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Primeiro, o ccam pergunta sobre cada arquivo que difere da cópia do grupo (manter o do grupo, usar o meu ou mesclar). Depois, a alteração roda em segundo plano com o mesmo código do [`ccam group`](/features/share): a configuração é salva primeiro, os dados de cada conta são movidos em ordem, e os originais vão para `.ccam-backup/`. Uma alteração que moveria os dados de uma conta em execução é recusada, com o nome da conta.

## Adicionar uma conta

`a` abre o assistente que o `ccam add` também usa em um terminal:

1. O nome da conta.
2. Como ela faz login: com uma assinatura do Claude agora (no navegador), com uma [chave de API](/features/api) ou com uma assinatura do Claude depois.
3. Para uma chave de API: de onde o ccam lê a chave a cada inicialização (`env:NAME` ou `command:COMMAND`), o endpoint, como a chave é enviada e o modelo.
4. Os grupos em que a conta entra, com `default` marcado; antes de existirem grupos, o que o grupo `default` compartilha.
5. Um comando de atalho, se você quiser.
6. Um resumo para confirmar.

`Esc` volta uma etapa e, na primeira etapa, fecha o assistente.

## Cores

As cores seguem os papéis da interface de terminal do Claude Code, para que o ccam pareça em casa ao lado dele; os nomes de grupos recebem uma cor derivada do próprio nome. O ccam pergunta ao terminal qual é a cor de fundo e escolhe um tema claro ou escuro. Para escolher um você mesmo, defina-o no `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Os temas `-daltonized` são para daltonismo vermelho-verde, e os temas `-ansi` usam as 16 cores do próprio terminal. `NO_COLOR` desativa as cores.

## Segurança e casos especiais

- O painel chama as mesmas funções que os comandos, então todas as regras valem: planos antes de mover dados, recusas para contas em execução, backups e tokens de login que nunca saem de uma conta.
- O painel não tem simulação: `ccam --dry-run` sem um nome imprime como a conta desta pasta seria iniciada, e `o` em uma conta oferece uma simulação só daquela inicialização.
- O terminal é restaurado antes de o Claude Code iniciar ou de o login abrir.
- Em uma tela pequena demais para um diálogo, as linhas da aba são as primeiras a sair; as linhas são cortadas com `…` em vez de quebradas.

## Relacionados

- [Compartilhamento com grupos](/features/share)
- [Iniciar o Claude Code](/features/launch)
- [Gerenciar contas](/features/accounts)
- [Limites de uso](/features/usage)
