---
title: Início rápido
metaTitle: "Início rápido: crie e troque contas do Claude Code"
description: Adicione duas contas do Claude Code com ccam add, dê a elas atalhos como ccwork, alterne entre elas e escolha o que elas compartilham.
updated: 2026-10-24
order: 2
translation: machine
sourceHash: "e98d0637356a"
---

Adicione cada conta com `ccam add` e depois inicie o Claude Code com `ccam <name>` ou com um atalho que você mesmo nomeia, como `ccwork`. Isso é tudo de que a maioria das pessoas precisa.

## 1. Adicione sua primeira conta

```sh
ccam add
```

O assistente pede um nome, como `personal` ou `work`, quais itens a conta deve compartilhar com suas outras contas e se deve criar um atalho (veja o passo 4). Depois ele executa `claude auth login`, que abre o navegador para você entrar. Em uma máquina remota, o navegador pode parar em uma página `http://localhost:…/callback` que não carrega; cole esse endereço no terminal e o ccam conclui o login. Quando ele termina, você vê:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Traga suas configurações atuais

Se você já usava o Claude Code, suas configurações, o `CLAUDE.md`, as extensões e os servidores MCP estão em `~/.claude`. Copie-os para as configurações compartilhadas, para que todas as contas os recebam:

```sh
ccam migrate
```

O ccam mostra o que vai copiar e pergunta antes. Ele só lê `~/.claude`, nunca copia um login e deixa o histórico de conversas de fora, a menos que você adicione `--history`. Consulte [copiar uma configuração existente](/docs/sharing#copy-an-existing-setup) para ver os detalhes.

## 3. Adicione a segunda conta

```sh
ccam add work
```

Entre com a outra conta no navegador.

## 4. Troque de conta

`ccam work` inicia o Claude Code como a conta work. Para algo mais curto, dê à conta um atalho próprio, como `ccwork`. Depois dos itens compartilhados, o assistente pergunta:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Pressione `y` e digite um nome, ou pressione Enter para usar o sugerido:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Depois do login, o ccam confirma o atalho:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Se você pulou a pergunta, pode adicionar um atalho a qualquer momento:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Agora qualquer um destes comandos inicia o Claude Code como a conta work:

```sh
ccwork
ccam work
ccam 2
```

Executar `ccam` sozinho abre o [painel](/features/dashboard): vá até uma conta com as setas ou com o número dela e pressione Enter.

Tudo o que vem depois do atalho ou do nome da conta vai para o Claude Code sem alterações, então você continua usando as opções que já conhece:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Escolha o que compartilhar

Por padrão, os servidores MCP, o `settings.json`, o `CLAUDE.md` global e os agentes, comandos, skills e plugins são compartilhados entre as contas por meio do grupo `default`. O histórico de conversas não é. Adicione-o ao grupo se você quiser continuar conversas de uma conta na outra:

```sh
ccam group set default --add history
```

O ccam mostra o que vai mover em cada conta antes de alterar qualquer coisa. Consulte [compartilhamento com grupos](/features/share) para ver os detalhes.

## 6. Quando uma conta atinge o limite

```sh
ccam usage
ccam next --continue
```

`ccam usage` mostra o uso de 5 horas e o uso semanal de todas as contas. `ccam next --continue` abre a mesma conversa na próxima conta. Em scripts, `ccam --fallback work -p "…"` faz a troca sozinho. Consulte [limites de uso](/docs/usage-limits).

**Assentos dos planos Team e Enterprise.** A escolha automática (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel) também alterna entre assentos da mesma organização Claude Team ou Enterprise. A central de ajuda da Anthropic diz "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (cada assento de um plano Team é atribuído a uma pessoa específica, e só essa pessoa pode usá-lo) e "Users can't route their work through another user's account to use that user's remaining usage" (os usuários não podem encaminhar o próprio trabalho pela conta de outro usuário para aproveitar o uso restante desse usuário) ([planos Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); o mesmo vale para os [planos Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Por isso, usar a escolha automática entre assentos pode violar os termos da Anthropic. Consulte a sua organização e a Anthropic antes de contar com ela. Os planos pessoais (Pro, Max) são tratados na [nota sobre os termos nas perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Já usa ~/.claude?

Registre essa pasta em vez de entrar novamente:

```sh
ccam import
```

A pasta permanece onde está, e o `claude` simples continua funcionando como antes.
