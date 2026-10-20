---
title: O painel
description: O ccam em um terminal abre um painel em tela cheia com contas, grupos, regras de pasta, extensões e verificações; cada mudança mostra o plano antes.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` sem argumentos abre um painel em tela cheia quando a entrada e a saída padrão são um terminal; caso contrário, imprime `ccam list`. O painel mostra todas as contas com o uso e os grupos de cada uma e faz o que os comandos fazem: iniciar, adicionar, fazer login, alterar grupos, definir regras de pasta, verificar e reparar. Iniciar uma conta fecha o painel; adicionar uma conta ou fazer login sai dele para o terminal e depois volta.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

O primeiro quadro é desenhado só a partir do `config.toml`, em bem menos de 100 ms; o uso aparece à medida que cada conta responde. Uma porcentagem fica amarela a partir de 80% e vermelha em 100%, e os valores vindos de um cache depois de uma consulta que falhou aparecem esmaecidos.

## Abas

| Aba | O que mostra | Teclas |
| --- | --- | --- |
| Accounts | número, nome, e-mail, barra de 5 horas, uso de 7 dias, último uso, grupos; `✔` a conta padrão, `▸` a conta desta pasta, `•` em execução | abaixo |
| Groups | cada grupo, o que ele compartilha e suas contas, e as contas que não estão em nenhum grupo | `enter` editar, `n` novo, `e` renomear, `x` remover |
| Folders | as regras do [`ccam folder`](/features/folder) e a que vale aqui | `n` usar uma conta aqui, `enter` alterar, `x` remover |
| Extensions | o que o [`ccam ext`](/features/ext) mostra, cada extensão uma vez, com os grupos de onde ela vem | `r` recarregar |
| Health | o que o [`ccam doctor`](/features/doctor) relata | `f` reparar (`doctor --fix`) |

`←`/`→` ou `Tab` trocam de aba, `?` lista as teclas da aba atual, `r` consulta o uso de novo sem o cache curto, e `q` ou `Esc` sai.

## Teclas na aba Accounts

| Tecla | O que faz |
| --- | --- |
| `Enter` | Inicia o Claude Code como a conta selecionada |
| `1`–`9` | Vai para a conta com esse número; digite `1` `2` rapidamente para ir à 12. Só `Enter` inicia |
| `↑`/`↓`, `k`/`j` | Move o cursor |
| `n` | Inicia a conta que o [`ccam next --auto`](/features/next) escolheria |
| `o` | Inicia com opções: `--continue`, `--resume`, outros argumentos ou uma simulação |
| `u` | Detalhes do uso: as janelas de 5 horas, de 7 dias e por modelo, com os reinícios e a origem dos números |
| `a` | Adiciona uma conta (o assistente abaixo) |
| `l` | Faz login novamente |
| `s` | Escolhe os grupos da conta |
| `d` | Torna a conta padrão |
| `f` | Usa a conta na pasta atual |
| `e` | Renomeia a conta |
| `c` | Define ou remove o comando de atalho dela |
| `x` | Remove a conta (pergunta antes, começando em No) |

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
