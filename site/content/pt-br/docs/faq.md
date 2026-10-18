---
title: Perguntas frequentes
description: Respostas sobre o ccam, contas do Claude Code, limites de uso, compartilhamento, privacidade e sistemas suportados.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 12ad0faef2ce
---

## O ccam é uma ferramenta oficial da Anthropic?

Não. O ccam é um projeto de código aberto não oficial e não é afiliado à Anthropic nem endossado por ela. Ele usa a configuração `CLAUDE_CONFIG_DIR`, que o Claude Code documenta para executar várias contas lado a lado.

## Posso usar várias contas do Claude Code?

Isso depende dos termos de cada conta e, no caso de contas de trabalho, da política da sua organização. Verifique-os antes de usar várias contas. O ccam apenas facilita a troca entre contas que você já tem.

## O ccam contorna os limites de uso?

Não. Cada conta mantém os próprios limites. O ccam mostra esses limites e permite que você continue em outra conta sua quando uma delas se esgota.

## Como usar duas contas do Claude Code no mesmo computador?

Dê a cada conta a própria pasta de configuração com `CLAUDE_CONFIG_DIR`. O ccam faz isso por você: execute `ccam add` para cada conta e depois inicie-as com `claude1` e `claude2`. O [guia](/guides/multiple-claude-code-accounts) também mostra como fazer isso manualmente.

## Como continuar uma conversa em outra conta?

Compartilhe o histórico de conversas entre as duas contas com `ccam share <name> --history on` e depois execute `ccs next --continue`. O Claude Code abre a conversa mais recente da pasta atual na próxima conta.

## O ccam vê ou envia meu login?

O ccam nunca copia nem compartilha um login. Ele lê o token de acesso de uma conta apenas para consultar o uso dessa conta no endpoint de uso da Anthropic, quando você executa `ccs usage`, `ccs next --auto` ou `ccam doctor`. O ccam não tem telemetria e não envia nada para nenhum outro lugar.

## Onde minhas contas ficam armazenadas?

Em `~/.local/share/ccam/profiles/<name>/`, uma pasta por conta. No macOS, o Claude Code guarda o próprio login no Acesso às Chaves, com um nome derivado dessa pasta.

## O que acontece com a minha pasta ~/.claude existente?

Nada, a menos que você execute `ccam import`. A importação registra a pasta no lugar onde ela está, então o `claude` simples e o `ccs <name>` continuam usando a mesma conta e as mesmas configurações.

## Quais sistemas o ccam suporta?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine e outras distribuições) e macOS. No Windows, use-o dentro do WSL.

## Por que um servidor MCP aparece em apenas uma conta?

Servidores adicionados dentro de uma sessão com `claude mcp add` ficam armazenados apenas para essa conta. O ccam avisa você na próxima vez que você iniciar a conta; mova o servidor para o grupo compartilhado com `ccam mcp import <account> <server>`.

## O ccam deixa o Claude Code mais lento?

Iniciar o Claude Code pelo ccam acrescenta cerca de 4 milissegundos, medidos com `scripts/perf-check.sh` em um notebook. O ccam lê um arquivo pequeno, verifica os links compartilhados e depois se substitui pelo `claude`, então ele não fica em execução enquanto você trabalha.

## Como desinstalar o ccam?

Exclua `ccam`, `ccs` e os links `claude1`… de `~/.local/bin` e depois `~/.config/ccam`, `~/.local/share/ccam` e `~/.local/state/ccam` se você não precisar mais das contas. Consulte [instalação](/docs/install#uninstall).
