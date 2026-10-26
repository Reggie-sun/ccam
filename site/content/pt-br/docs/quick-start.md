---
title: Início rápido
metaTitle: "Início rápido: crie e troque contas do Claude Code"
description: Adicione duas contas do Claude Code com ccam add, dê a elas atalhos como ccwork, alterne entre elas e escolha o que elas compartilham.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "bd893bbb83ab"
---

Adicione cada conta com `ccam add` e depois inicie o Claude Code com `ccam <name>` ou com um atalho que você mesmo nomeia, como `ccwork`. Isso é tudo de que a maioria das pessoas precisa.

Já usa o Claude Code? Leia [Já usa ~/.claude?](#already-using-claude) no fim desta página antes de começar.

## 1. Adicione sua primeira conta

```sh
ccam add
```

O assistente pede um nome, como `personal` ou `work`, quais itens a conta deve compartilhar com suas outras contas e se deve criar um atalho (veja o passo 3). Depois ele executa `claude auth login`, que abre o navegador para você entrar. Em uma máquina remota, o seu navegador pode acabar em uma página `http://localhost:…/callback` que não carrega; cole esse endereço no terminal e o ccam conclui o login. Quando o login termina, você vê:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Traga suas configurações atuais

Se você já usava o Claude Code, suas configurações, o `CLAUDE.md` e as extensões estão em `~/.claude`, e os seus servidores MCP estão em `~/.claude.json`. Copie-os para o grupo compartilhado `default`, para que todas as contas os recebam:

```sh
ccam migrate
```

O ccam mostra o que vai copiar e pergunta antes. Ele só lê `~/.claude`, nunca copia um login e deixa o histórico de conversas de fora, a menos que você adicione `--history`. Consulte [copiar uma configuração existente](/docs/sharing#copy-an-existing-setup) para ver os detalhes.

## 3. Adicione a segunda conta

```sh
ccam add work
```

O assistente pula o nome, que você já informou na linha de comando, e mostra os grupos dos quais participar, com `default` marcado. Desta vez, dê à conta um atalho próprio, como `ccwork`, para iniciá-la com uma só palavra. Depois dos grupos, o assistente pergunta:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Pressione `y` e digite um nome, ou pressione Enter para usar o sugerido:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Depois, entre com a outra conta no navegador. O ccam confirma a conta e o atalho dela:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Troque de conta

`ccam work` inicia o Claude Code como a conta work. Se você pulou a pergunta do atalho, pode adicionar um a qualquer momento:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Agora cada um destes comandos inicia o Claude Code como a conta work (`2` é o número dela em `ccam list`):

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

Por padrão, os servidores MCP, o `settings.json`, o `CLAUDE.md` global e as extensões (agentes, comandos, skills e plugins) são compartilhados entre as contas por meio do grupo `default`. O histórico de conversas não é. Adicione-o ao grupo se você quiser continuar conversas de uma conta na outra:

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

**Assentos dos planos Team e Enterprise.** A escolha automática de conta (`ccam next`, `--fallback`, `ccam pick`, a tecla `n` do painel e, depois que você o ativa, o [`auto_switch`](/features/next#auto-switch)) também alterna entre assentos da mesma organização Claude Team ou Enterprise, o que pode violar os termos da Anthropic. Leia as [perguntas frequentes](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e consulte a sua organização e a Anthropic antes de contar com ela.

## Já usa ~/.claude?

Registre essa pasta em vez de entrar novamente:

```sh
ccam import
```

A pasta permanece onde está, e o `claude` simples continua funcionando como antes.
