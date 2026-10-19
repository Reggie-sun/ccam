---
title: Início rápido
description: Adicione duas contas do Claude Code com ccam add, dê a elas atalhos como ccwork, alterne entre elas e escolha o que elas compartilham.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: b73914ff371f
---

Adicione cada conta com `ccam add` e depois inicie o Claude Code com `ccam <name>` ou com um atalho que você mesmo nomeia, como `ccwork`. Isso é tudo de que a maioria das pessoas precisa.

## 1. Adicione sua primeira conta

```sh
ccam add
```

O assistente pede um nome, como `personal` ou `work`, quais itens a conta deve compartilhar com suas outras contas e se deve criar um atalho (veja o passo 4). Depois ele executa `claude auth login`, que abre o navegador para você entrar. Quando ele termina, você vê:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

O ccam imprime as próprias mensagens em inglês; os comandos e as opções são os mesmos em todos os idiomas.

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

Executar `ccam` sozinho abre um seletor: use as setas e Enter, ou pressione o número da conta.

Tudo o que vem depois do atalho ou do nome da conta vai para o Claude Code sem alterações, então você continua usando as opções que já conhece:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Escolha o que compartilhar

Por padrão, os servidores MCP, o `settings.json`, o `CLAUDE.md` global e os agentes, comandos, skills e plugins são compartilhados entre as contas. O histórico de conversas não é. Ative-o nas duas contas se você quiser continuar conversas de uma conta na outra:

```sh
ccam share personal --history on
ccam share work --history on
```

O ccam mostra o que vai mover antes de alterar qualquer coisa. Consulte [compartilhamento](/docs/sharing) para ver os detalhes.

## 6. Quando uma conta atinge o limite

```sh
ccam usage
ccam next --continue
```

`ccam usage` mostra o uso de 5 horas e o uso semanal de todas as contas. `ccam next --continue` abre a mesma conversa na próxima conta. Consulte [limites de uso](/docs/usage-limits).

## Já usa ~/.claude?

Registre essa pasta em vez de entrar novamente:

```sh
ccam import
```

A pasta permanece onde está, e o `claude` simples continua funcionando como antes.
