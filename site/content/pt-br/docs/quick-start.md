---
title: Início rápido
description: Adicione duas contas do Claude Code com ccam add, troque entre elas com claude1, claude2 ou ccs e escolha o que elas compartilham.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 9462b80d61bc
---

Adicione cada conta com `ccam add` e depois inicie o Claude Code com `claude1`, `claude2` ou `ccs <name>`. Isso é tudo de que a maioria das pessoas precisa.

## 1. Adicione sua primeira conta

```sh
ccam add
```

O assistente pede um nome, como `personal` ou `work`, e quais itens a conta deve compartilhar com suas outras contas. Depois ele executa `claude auth login`, que abre o navegador para você entrar. Quando ele termina, você vê:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

Por enquanto, o ccam imprime as próprias mensagens em coreano; os comandos e as opções são os mesmos em todos os idiomas.

## 2. Adicione a segunda conta

```sh
ccam add work
```

Entre com a outra conta no navegador. Ela se torna `claude2`.

## 3. Troque de conta

Qualquer um destes comandos inicia o Claude Code como a conta de trabalho:

```sh
claude2
ccs work
ccs 2
```

Tudo o que vem depois do nome da conta vai para o Claude Code sem alterações, então você continua usando as opções que já conhece:

```sh
ccs work --resume
ccs work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Executar `ccs` sozinho abre um seletor: use as setas e Enter, ou pressione o número da conta.

## 4. Escolha o que compartilhar

Por padrão, os servidores MCP, o `settings.json`, o `CLAUDE.md` global e os agentes, comandos, skills e plugins são compartilhados entre as contas. O histórico de conversas não é. Ative-o nas duas contas se você quiser continuar conversas de uma conta na outra:

```sh
ccam share personal --history on
ccam share work --history on
```

O ccam mostra o que vai mover antes de alterar qualquer coisa. Consulte [compartilhamento](/docs/sharing) para ver os detalhes.

## 5. Quando uma conta atinge o limite

```sh
ccs usage
ccs next --continue
```

`ccs usage` mostra o uso de 5 horas e o uso semanal de todas as contas. `ccs next --continue` abre a mesma conversa na próxima conta. Consulte [limites de uso](/docs/usage-limits).

## Já usa ~/.claude?

Registre essa pasta em vez de entrar novamente:

```sh
ccam import
```

A pasta permanece onde está, e o `claude` simples continua funcionando como antes.
