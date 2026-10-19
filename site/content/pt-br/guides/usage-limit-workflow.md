---
title: O que fazer quando o Claude Code atinge o limite de uso
description: Veja quanto espaço resta em cada conta do Claude Code e continue a mesma conversa em outra conta quando uma delas atingir o limite.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: db7ed9b7bcfc
---

Quando uma conta atingir o limite, continue a mesma conversa em outra conta sua: compartilhe o histórico de conversas uma vez e depois execute `ccam next --continue`. Se você tem apenas uma conta, espere o horário de reinício que o Claude Code mostra.

## Antes que aconteça

Configure as contas para que a troca leve um único comando:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Agora a linha de status mostra a conta e o uso de 5 horas dela enquanto você trabalha, então o limite não pega você de surpresa.

## Quando o limite é atingido

Saia do Claude Code (`/exit`) e depois execute:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

Abra a conversa na conta com mais espaço:

```sh
ccam next --auto --continue
```

O Claude Code inicia como `work` com a conversa mais recente da pasta atual, incluindo tudo o que foi dito até agora.

## O que ter em mente

- Execute `ccam next --continue` na mesma pasta em que você estava trabalhando; as conversas são armazenadas por pasta de projeto.
- Evite usar duas contas na mesma conversa ao mesmo tempo. O ccam avisa se você tentar.
- A outra conta também precisa compartilhar o histórico de conversas, no mesmo grupo.
- Se você pode trocar de conta dessa forma depende dos termos de cada conta.

Consulte [limites de uso](/docs/usage-limits) para ver todas as opções de `ccam usage` e `ccam next`.
